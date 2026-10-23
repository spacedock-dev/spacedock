package release

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type extrasRun struct {
	t               *testing.T
	home            string
	bin             string
	data            string
	dist            string
	extras          string
	log             string
	pathDir         []string
	launcherVersion string
}

func writeTarGzFiles(t *testing.T, path string, files map[string]string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
}

const fakeExtrasRows = `alpha|spacedock-dev/alpha|alpha|alpha-tui alpha-client|alpha-tui:spacedock-alpha|bin
beta|spacedock-dev/beta|beta|beta/beta|beta:spacedock-beta|tree`

func scriptWithExtras(t *testing.T, rows string) string {
	t.Helper()
	src, err := os.ReadFile(filepath.Join("..", "..", "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	table := regexp.MustCompile(`(?s)\nEXTRAS='[^']*'\n`)
	if !table.Match(src) {
		t.Fatal("install.sh declares no EXTRAS table; the fake rows have nowhere to go")
	}
	script := table.ReplaceAllLiteral(src, []byte("\nEXTRAS='"+rows+"'\n"))
	path := filepath.Join(t.TempDir(), "install.sh")
	if err := os.WriteFile(path, script, 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func newExtrasRun(t *testing.T, betaVersion string) *extrasRun {
	t.Helper()
	return newExtrasRunWithLauncher(t, betaVersion, "9.9.9")
}

func newExtrasRunWithLauncher(t *testing.T, betaVersion, launcherVersion string) *extrasRun {
	t.Helper()
	osName, arch := goosArch(t)
	root := t.TempDir()
	r := &extrasRun{t: t, home: filepath.Join(root, "home"), bin: filepath.Join(root, "bin"), data: filepath.Join(root, "data"),
		dist: filepath.Join(root, "dist"), extras: filepath.Join(root, "extras"), log: filepath.Join(root, "calls.log"), launcherVersion: launcherVersion}
	for _, d := range []string{r.home, r.bin, r.data, r.dist, r.extras} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	spacedock := "#!/bin/sh\necho \"installed-spacedock $*\" >> '" + r.log + "'\necho 'spacedock " + launcherVersion + "'\n"
	writeTarGzFiles(t, filepath.Join(r.dist, "spacedock_"+launcherVersion+"_"+osName+"_"+arch+".tar.gz"), map[string]string{"spacedock": spacedock})
	r.writeChecksums(r.dist)

	r.stageExtras(osName, arch, betaVersion)
	return r
}

func (r *extrasRun) stageExtras(osName, arch, betaVersion string) {
	r.t.Helper()
	versionScript := "#!/bin/sh\necho 0.16.0\n"
	writeTarGzFiles(r.t, filepath.Join(r.extras, "alpha_0.16.0_"+osName+"_"+arch+".tar.gz"), map[string]string{
		"alpha-tui": versionScript, "alpha-client": versionScript})
	matches, _ := filepath.Glob(filepath.Join(r.extras, "beta_*"))
	for _, m := range matches {
		os.Remove(m)
	}
	writeTarGzFiles(r.t, filepath.Join(r.extras, "beta_"+betaVersion+"_"+osName+"_"+arch+".tar.gz"), map[string]string{
		"beta/beta": "#!/bin/sh\necho beta\n"})
	r.writeChecksums(r.extras)
}

func (r *extrasRun) writeChecksums(dir string) {
	r.t.Helper()
	archives, _ := filepath.Glob(filepath.Join(dir, "*.tar.gz"))
	var lines []string
	for _, a := range archives {
		lines = append(lines, sha256OfFile(r.t, a)+"  "+filepath.Base(a))
	}
	if err := osWriteFile(filepath.Join(dir, "checksums.txt"), strings.Join(lines, "\n")+"\n"); err != nil {
		r.t.Fatal(err)
	}
}

func (r *extrasRun) addProgram(name, body string) string {
	r.t.Helper()
	dir := r.t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		r.t.Fatal(err)
	}
	r.pathDir = append(r.pathDir, dir)
	return dir
}

func (r *extrasRun) run() (stderr string, code int) {
	r.t.Helper()
	return r.runWith("SPACEDOCK_INSTALL_FROM="+r.dist, "SPACEDOCK_EXTRAS=1", "SPACEDOCK_EXTRAS_FROM="+r.extras)
}

func (r *extrasRun) runWith(env ...string) (stderr string, code int) {
	r.t.Helper()
	cmd := exec.Command("sh", scriptWithExtras(r.t, fakeExtrasRows))
	cmd.Env = append([]string{
		"PATH=" + strings.Join(append(append([]string{}, r.pathDir...), "/usr/bin", "/bin"), ":"),
		"HOME=" + r.home, "XDG_DATA_HOME=" + r.data, "SPACEDOCK_INSTALL_DIR=" + r.bin,
	}, env...)
	var errOut strings.Builder
	cmd.Stderr = &errOut
	err := cmd.Run()
	r.t.Logf("install.sh stderr:\n%s", errOut.String())
	if ee, ok := err.(*exec.ExitError); ok {
		return errOut.String(), ee.ExitCode()
	} else if err != nil {
		r.t.Fatal(err)
	}
	return errOut.String(), 0
}

func (r *extrasRun) calls() string {
	data, _ := os.ReadFile(r.log)
	return string(data)
}

func TestInstallPluginsRunThroughTheJustInstalledBinary(t *testing.T) {
	r := newExtrasRun(t, "1.0.1")
	r.addProgram("spacedock", "[ \"$1\" = --version ] || echo 'older copy ran' >> '"+r.log+"'\necho 'spacedock 0.27.2'")
	r.addProgram("claude", "exit 0")

	stderr, code := r.run()
	if code != 0 {
		t.Fatalf("install.sh exited %d", code)
	}
	if got := r.calls(); !strings.Contains(got, "installed-spacedock install --host claude") || strings.Contains(got, "older copy ran") {
		t.Errorf("calls = %q, want the installed spacedock to run `install --host claude` and the older copy not to run", got)
	}
	if !strings.Contains(stderr, "(spacedock 0.27.2) comes first on PATH") {
		t.Errorf("stderr does not name the older copy that comes first:\n%s", stderr)
	}
}

func TestInstallWithoutClaudeSaysHowToAddThePluginsLater(t *testing.T) {
	r := newExtrasRun(t, "1.0.1")
	stderr, code := r.run()
	if code != 0 {
		t.Fatalf("install.sh exited %d", code)
	}
	if strings.Contains(r.calls(), "install --host claude") {
		t.Errorf("plugin install ran without claude on PATH: %q", r.calls())
	}
	if !strings.Contains(stderr, "claude is not on PATH; run `spacedock install --host claude`") {
		t.Errorf("stderr lacks the instruction for later:\n%s", stderr)
	}
}

// A file added inside the installed beta tree survives only when nothing re-extracts it.
func TestInstallExtrasSecondRunLeavesInstalledVersionsAlone(t *testing.T) {
	r := newExtrasRun(t, "1.0.1")
	if _, code := r.run(); code != 0 {
		t.Fatalf("first install exited %d", code)
	}
	marker := filepath.Join(r.data, "beta", "1.0.1", "marker")
	if err := os.WriteFile(marker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	stderr, code := r.run()
	if code != 0 {
		t.Fatalf("second install exited %d", code)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("second run extracted beta 1.0.1 again: %v", err)
	}
	for _, want := range []string{"alpha 0.16.0 unchanged", "beta 1.0.1 unchanged"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr lacks %q:\n%s", want, stderr)
		}
	}
}

func TestInstallTreeExtraKeepsNewAndPreviousVersionOnly(t *testing.T) {
	r := newExtrasRun(t, "1.0.0")
	if _, code := r.run(); code != 0 {
		t.Fatalf("first install exited %d", code)
	}
	if err := os.MkdirAll(filepath.Join(r.data, "beta", "0.9.0"), 0o755); err != nil {
		t.Fatal(err)
	}
	osName, arch := goosArch(t)
	r.stageExtras(osName, arch, "1.0.1")

	if _, code := r.run(); code != 0 {
		t.Fatalf("update install exited %d", code)
	}
	entries, err := os.ReadDir(filepath.Join(r.data, "beta"))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	if strings.Join(got, ",") != "1.0.0,1.0.1" {
		t.Errorf("beta versions on disk = %v, want [1.0.0 1.0.1]", got)
	}
	target, err := os.Readlink(filepath.Join(r.bin, "beta"))
	if err != nil || !strings.Contains(target, filepath.Join("beta", "1.0.1")) {
		t.Errorf("beta link = %q (%v), want it to point into 1.0.1", target, err)
	}
}

func TestInstallExtrasRepairsAPartialInstall(t *testing.T) {
	for _, missing := range []string{"alpha-client", "spacedock-alpha"} {
		t.Run(missing, func(t *testing.T) {
			r := newExtrasRun(t, "1.0.1")
			if _, code := r.run(); code != 0 {
				t.Fatalf("first install exited %d", code)
			}
			if err := os.Remove(filepath.Join(r.bin, missing)); err != nil {
				t.Fatal(err)
			}

			stderr, code := r.run()
			if code != 0 {
				t.Fatalf("second install exited %d", code)
			}
			if _, err := os.Stat(filepath.Join(r.bin, missing)); err != nil {
				t.Errorf("%s is still missing after the second run: %v", missing, err)
			}
			if strings.Contains(stderr, "alpha 0.16.0 unchanged") {
				t.Errorf("a partial install was reported unchanged:\n%s", stderr)
			}
		})
	}
}

func TestInstallTreeExtraKeepsOtherVersionsWhenThisDirectoryHasNoLink(t *testing.T) {
	r := newExtrasRun(t, "1.0.1")
	other := filepath.Join(r.data, "beta", "0.9.0")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}

	if _, code := r.run(); code != 0 {
		t.Fatalf("install exited %d", code)
	}
	if _, err := os.Stat(other); err != nil {
		t.Errorf("a first install into this directory removed a version it never linked: %v", err)
	}
}

func TestInstallWithNoExtensionTableRowInstallsOnlySpacedock(t *testing.T) {
	r := newExtrasRun(t, "1.0.1")
	cmd := exec.Command("sh", scriptWithExtras(t, ""))
	cmd.Env = []string{
		"PATH=/usr/bin:/bin", "HOME=" + r.home, "XDG_DATA_HOME=" + r.data,
		"SPACEDOCK_INSTALL_DIR=" + r.bin, "SPACEDOCK_INSTALL_FROM=" + r.dist,
		"SPACEDOCK_EXTRAS=1", "SPACEDOCK_EXTRAS_FROM=" + r.extras,
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("install.sh failed: %v\n%s", err, out)
	}
	if strings.Contains(string(out), "failed to install") || strings.Contains(string(out), "installed alpha") {
		t.Errorf("an empty table installed or failed an extension:\n%s", out)
	}
	if entries, _ := os.ReadDir(r.data); len(entries) != 0 {
		t.Errorf("an empty table wrote under the data directory: %v", entries)
	}
}

// stubCurl serves install.sh's release lookups and downloads from the staged
// directories and logs every URL, so a run without SPACEDOCK_INSTALL_FROM stays
// offline. STUB_API=fail and STUB_REDIRECT=fail|norelease break those two lookups,
// STUB_FAIL_REPO breaks every lookup of one repo, and the stub git prints
// STUB_TAGS for `ls-remote --tags` (nothing when unset).
func (r *extrasRun) stubCurl(osName, arch string) (log string) {
	r.t.Helper()
	log = filepath.Join(r.t.TempDir(), "urls.log")
	stable := filepath.Join(r.dist, "spacedock_"+r.launcherVersion+"_"+osName+"_"+arch+".tar.gz")
	src, err := os.ReadFile(stable)
	if err != nil {
		r.t.Fatal(err)
	}
	if err := os.WriteFile(strings.TrimSuffix(stable, ".tar.gz")+"_edge.tar.gz", src, 0o644); err != nil {
		r.t.Fatal(err)
	}
	r.writeChecksums(r.dist)
	body := `#!/bin/sh
out=""; url=""
while [ $# -gt 0 ]; do
	case "$1" in -o) out="$2"; shift ;; -w) shift ;; -*) ;; *) url="$1" ;; esac
	shift
done
echo "$url" >> '` + log + `'
repo="${url#https://api.github.com/repos/}"; repo="${repo#https://github.com/}"; repo="${repo%%/releases*}"
case "$repo" in
	spacedock-dev/spacedock) tag=v` + r.launcherVersion + ` ;; spacedock-dev/alpha) tag=v0.16.0 ;; *) tag=v1.0.1 ;;
esac
case "$url" in
	https://api.github.com/*)
		[ "$STUB_API" = fail ] || [ "$repo" = "$STUB_FAIL_REPO" ] && exit 22
		printf '{"tag_name": "%s"}\n' "$tag" ;;
	https://github.com/*/releases/latest)
		[ "$STUB_REDIRECT" = fail ] || [ "$repo" = "$STUB_FAIL_REPO" ] && exit 22
		if [ "$STUB_REDIRECT" = norelease ]; then printf 'https://github.com/%s/releases' "$repo"; else printf 'https://github.com/%s/releases/tag/%s' "$repo" "$tag"; fi ;;
	https://github.com/*/releases/download/*)
		dir='` + r.extras + `'; [ "$repo" = spacedock-dev/spacedock ] && dir='` + r.dist + `'
		cp "$dir/${url##*/}" "$out" ;;
	*) exit 22 ;;
esac
`
	dir := r.addProgram("curl", "")
	if err := os.WriteFile(filepath.Join(dir, "curl"), []byte(body), 0o755); err != nil {
		r.t.Fatal(err)
	}
	r.addProgram("git", `for t in $STUB_TAGS; do printf 'abc123\trefs/tags/%s\n' "$t"; done`)
	return log
}

func TestInstallExtrasRunOnEveryChannelUnlessNone(t *testing.T) {
	cases := []struct {
		name      string
		env       []string
		extension bool
	}{
		{"stable installs them", []string{"SPACEDOCK_CHANNEL=stable"}, true},
		{"edge installs them", []string{"SPACEDOCK_CHANNEL=edge"}, true},
		{"stable with none installs none", []string{"SPACEDOCK_CHANNEL=stable", "SPACEDOCK_EXTRAS=none"}, false},
		{"edge with none installs none", []string{"SPACEDOCK_CHANNEL=edge", "SPACEDOCK_EXTRAS=none"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newExtrasRun(t, "1.0.1")
			r.addProgram("claude", "exit 0")
			osName, arch := goosArch(t)
			log := r.stubCurl(osName, arch)

			stderr, code := r.runWith(tc.env...)
			if code != 0 {
				t.Fatalf("install.sh exited %d", code)
			}
			urls, _ := os.ReadFile(log)
			if !strings.Contains(stderr, "installed spacedock") {
				t.Errorf("spacedock itself was not installed:\n%s", stderr)
			}
			_, linkErr := os.Stat(filepath.Join(r.bin, "spacedock-alpha"))
			asked := strings.Contains(string(urls), "spacedock-dev/alpha/")
			if tc.extension != (linkErr == nil) || tc.extension != asked {
				t.Errorf("extension wanted=%v, but spacedock-alpha linked=%v and its release asked=%v\nurls:\n%s\nstderr:\n%s",
					tc.extension, linkErr == nil, asked, urls, stderr)
			}
			pluginStep := strings.Contains(r.calls(), "install --host claude")
			if tc.extension != pluginStep || (!tc.extension && strings.Contains(stderr, "Claude Code")) {
				t.Errorf("plugin step wanted=%v, but it ran=%v\ncalls:\n%s\nstderr:\n%s", tc.extension, pluginStep, r.calls(), stderr)
			}
		})
	}
}

// The script on main can be newer than the latest stable launcher. Older stable
// launchers must keep their binary-only install until the combined stack ships.
func TestInstallExtrasRespectStableLauncherBoundary(t *testing.T) {
	cases := []struct {
		name      string
		version   string
		channel   string
		local     bool
		extension bool
	}{
		{name: "default stable 0.27.3", version: "0.27.3"},
		{name: "local stable 0.27.3", version: "0.27.3", local: true},
		{name: "stable 0.28.0", version: "0.28.0", extension: true},
		{name: "local stable 0.28.0", version: "0.28.0", local: true, extension: true},
		{name: "stable 0.29.0", version: "0.29.0", extension: true},
		{name: "stable 1.0.0", version: "1.0.0", extension: true},
		{name: "edge unchanged", version: "0.28.0-pre4", channel: "edge", extension: true},
		{name: "older edge unchanged", version: "0.27.3", channel: "edge", extension: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newExtrasRunWithLauncher(t, "1.0.1", tc.version)
			r.addProgram("claude", "exit 0")
			osName, arch := goosArch(t)
			log := r.stubCurl(osName, arch)
			var env []string
			if tc.channel != "" {
				env = append(env, "SPACEDOCK_CHANNEL="+tc.channel)
			}
			if tc.local {
				env = append(env, "SPACEDOCK_INSTALL_FROM="+r.dist, "SPACEDOCK_EXTRAS=1", "SPACEDOCK_EXTRAS_FROM="+r.extras)
			}
			stderr, code := r.runWith(env...)
			if code != 0 {
				t.Fatalf("install.sh exited %d:\n%s", code, stderr)
			}
			if _, err := os.Stat(filepath.Join(r.bin, "spacedock")); err != nil {
				t.Fatalf("launcher was not installed: %v", err)
			}
			urls, _ := os.ReadFile(log)
			for _, name := range []string{"alpha-tui", "alpha-client", "spacedock-alpha", "beta", "spacedock-beta"} {
				_, err := os.Stat(filepath.Join(r.bin, name))
				if installed := err == nil; installed != tc.extension {
					t.Errorf("%s installed=%v, want %v", name, installed, tc.extension)
				}
			}
			if !tc.extension {
				if strings.Contains(string(urls), "spacedock-dev/alpha/") || strings.Contains(string(urls), "spacedock-dev/beta/") {
					t.Errorf("binary-only install looked up extension releases:\n%s", urls)
				}
				if entries, _ := os.ReadDir(r.data); len(entries) != 0 {
					t.Errorf("binary-only install wrote extension data: %v", entries)
				}
			}
			if pluginStep := strings.Contains(r.calls(), "install --host claude"); pluginStep != tc.extension {
				t.Errorf("plugin bootstrap ran=%v, want %v; calls=%q", pluginStep, tc.extension, r.calls())
			}
		})
	}
}
