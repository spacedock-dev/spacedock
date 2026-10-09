// ABOUTME: AC-3/AC-4 front-door + init seam tests — version-gate fail-fast,
// ABOUTME: launch-seam argv on compatible, install-seam host commands, codex prose.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// fakeHost records every seam interaction and returns canned results so the
// front door / init paths run with no real host CLI, no exec, no network.
type fakeHost struct {
	// manifest is the path returned by ResolveManifest; "" means no plugin found.
	manifest string
	// manifestAfterInstall, when non-empty, is what ResolveManifest returns once
	// Install has been called — simulating a real install actually landing a
	// (usually compatible) plugin, so D6's re-gate-once step sees the fix. Left
	// empty (the default), ResolveManifest keeps returning the pre-install
	// manifest even after Install — simulating an install that does NOT change
	// resolution, the "second miss" case.
	manifestAfterInstall  string
	installed             bool
	resolveErr            error
	launchedArg           []string // argv captured by Launch
	launchedEnv           []string // env captured by Launch
	launchCode            int      // host exit code Launch returns (default 0)
	launchErr             error
	installCmds           []string // host commands captured by Install
	installOut            string
	inventory             []pluginInventoryEntry
	inventoryAfterInstall []pluginInventoryEntry
	inventoryErr          error
	inventoryCalls        int
}

func (f *fakeHost) ResolveManifest(host string) (string, error) {
	if f.installed && f.manifestAfterInstall != "" {
		return f.manifestAfterInstall, nil
	}
	return f.manifest, f.resolveErr
}

func (f *fakeHost) PluginInventory(host string) ([]pluginInventoryEntry, error) {
	f.inventoryCalls++
	if f.installed && f.inventoryAfterInstall != nil {
		return f.inventoryAfterInstall, f.inventoryErr
	}
	return f.inventory, f.inventoryErr
}

func (f *fakeHost) Launch(argv []string, env []string) (int, error) {
	f.launchedArg = argv
	f.launchedEnv = env
	return f.launchCode, f.launchErr
}

func (f *fakeHost) Install(host, source, branch string) (string, error) {
	f.installCmds = append(f.installCmds, host, source, branch)
	f.installed = true
	return f.installOut, nil
}

func (f *fakeHost) InstallCodexLocalPluginDir(source string) (string, error) {
	f.installCmds = append(f.installCmds, "codex", source)
	f.installed = true
	return f.installOut, nil
}

func stableInventory(selected, sibling bool) []pluginInventoryEntry {
	return []pluginInventoryEntry{
		{ID: "spacedock@spacedock", Version: "0.27.1", Installed: true, Enabled: selected},
		{ID: "spacedock@spacedock-edge", Version: "0.28.0-pre0", Installed: true, Enabled: sibling},
	}
}

func repeatedStableInventory() []pluginInventoryEntry {
	return append(stableInventory(true, false), pluginInventoryEntry{
		ID: "spacedock@spacedock-edge", Version: "0.28.0-pre0", Installed: true, Enabled: true,
	})
}

func TestStable0271FrontDoorsHealEnabledSiblingBeforeLaunch(t *testing.T) {
	withVersion(t, "0.27.1")
	savedBranch := devBranch
	devBranch = "main"
	t.Cleanup(func() { devBranch = savedBranch })

	frontDoors := []struct {
		name string
		run  func([]string, *fakeHost, *bytes.Buffer) int
	}{
		{name: "claude", run: func(args []string, host *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, t.TempDir(), host, lookFound, &stdout, stderr)
		}},
		{name: "codex", run: func(args []string, host *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, t.TempDir(), host, lookFound, &stdout, stderr)
		}},
	}

	for _, door := range frontDoors {
		t.Run(door.name+"/repeated-scope-auto-heal", func(t *testing.T) {
			host := &fakeHost{manifest: writeVersionedManifest(t, "0.27.1"), inventory: repeatedStableInventory(), inventoryAfterInstall: stableInventory(true, false)[:1]}
			var stderr bytes.Buffer
			if code := door.run(nil, host, &stderr); code != 0 || host.launchedArg == nil {
				t.Fatalf("exit=%d launched=%v stderr=%q", code, host.launchedArg != nil, stderr.String())
			}
			if len(host.installCmds) != 3 || host.inventoryCalls != 2 {
				t.Fatalf("install=%v inventoryCalls=%d, want one heal and verification", host.installCmds, host.inventoryCalls)
			}
		})
		t.Run(door.name+"/auto-heal", func(t *testing.T) {
			host := &fakeHost{manifest: writeVersionedManifest(t, "0.27.1"), inventory: stableInventory(true, true), inventoryAfterInstall: stableInventory(true, false)[:1]}
			var stderr bytes.Buffer
			if code := door.run(nil, host, &stderr); code != 0 || host.launchedArg == nil {
				t.Fatalf("exit=%d launched=%v stderr=%q", code, host.launchedArg != nil, stderr.String())
			}
			if len(host.installCmds) != 3 || host.inventoryCalls != 2 {
				t.Fatalf("install=%v inventoryCalls=%d, want one heal and verification", host.installCmds, host.inventoryCalls)
			}
		})
		t.Run(door.name+"/no-install", func(t *testing.T) {
			host := &fakeHost{manifest: writeVersionedManifest(t, "0.27.1"), inventory: stableInventory(true, true)}
			var stderr bytes.Buffer
			if code := door.run([]string{"--no-install"}, host, &stderr); code == 0 || host.launchedArg != nil || len(host.installCmds) != 0 {
				t.Fatalf("exit=%d launched=%v install=%v", code, host.launchedArg != nil, host.installCmds)
			}
			want := "Run `spacedock install --host " + door.name + "` to keep only the stable channel.\n"
			if stderr.String() != want {
				t.Fatalf("stderr=%q, want %q", stderr.String(), want)
			}
		})
		t.Run(door.name+"/inventory-failure", func(t *testing.T) {
			host := &fakeHost{manifest: writeVersionedManifest(t, "0.27.1"), inventoryErr: errors.New("host unavailable")}
			var stderr bytes.Buffer
			if code := door.run(nil, host, &stderr); code == 0 || host.launchedArg != nil || len(host.installCmds) != 0 {
				t.Fatalf("exit=%d launched=%v install=%v", code, host.launchedArg != nil, host.installCmds)
			}
			want := "Spacedock: could not verify the " + door.name + " plugin enablement state: host unavailable.\n" +
				"Run `spacedock install --host " + door.name + "` before launching.\n"
			if stderr.String() != want {
				t.Fatalf("stderr=%q, want %q", stderr.String(), want)
			}
		})
	}
}

func TestDoctorSiblingInventory(t *testing.T) {
	withVersion(t, "0.27.1")
	savedBranch := devBranch
	devBranch = "main"
	t.Cleanup(func() { devBranch = savedBranch })
	conflicts := []struct {
		name, host string
		inventory  []pluginInventoryEntry
	}{{"claude conflict", "claude", stableInventory(true, true)}, {"codex conflict", "codex", stableInventory(true, true)}, {"repeated sibling scopes", "claude", repeatedStableInventory()}}
	for _, conflict := range conflicts {
		t.Run(conflict.name, func(t *testing.T) {
			host := &fakeHost{manifest: writeVersionedManifest(t, "0.27.1"), inventory: conflict.inventory}
			var stdout, stderr bytes.Buffer
			if code := runDoctor(context.Background(), []string{"--host", conflict.host}, host, &stdout, &stderr); code != 0 {
				t.Fatalf("exit=%d stderr=%q", code, stderr.String())
			}
			want := "OK: spacedock binary 0.27.1 and plugin 0.27.1 are compatible.\n" +
				"CONFLICT: " + conflict.host + " can load a different Spacedock plugin than doctor checked.\n" +
				"  checked: spacedock@spacedock 0.27.1 (installed, enabled)\n" +
				"  sibling: spacedock@spacedock-edge 0.28.0-pre0 (installed, enabled)\n" +
				"Run `spacedock install --host " + conflict.host + "` to keep only the stable channel.\n"
			if stdout.String() != want || len(host.installCmds) != 0 {
				t.Fatalf("stdout=%q, want %q; doctor installs=%v", stdout.String(), want, host.installCmds)
			}
		})
	}
	t.Run("disabled sibling", func(t *testing.T) {
		host := &fakeHost{manifest: writeVersionedManifest(t, "0.27.1"), inventory: stableInventory(true, false)}
		var stdout, stderr bytes.Buffer
		want := "OK: spacedock binary 0.27.1 and plugin 0.27.1 are compatible.\n"
		if code := runDoctor(context.Background(), nil, host, &stdout, &stderr); code != 0 || stdout.String() != want {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	})
	t.Run("inventory failure", func(t *testing.T) {
		host := &fakeHost{manifest: writeVersionedManifest(t, "0.27.1"), inventoryErr: errors.New("host unavailable")}
		var stdout, stderr bytes.Buffer
		want := "OK: spacedock binary 0.27.1 and plugin 0.27.1 are compatible.\n" +
			"INCOMPLETE: doctor checked compatibility but did not read the claude plugin enablement state: host unavailable\n"
		if code := runDoctor(context.Background(), nil, host, &stdout, &stderr); code != 0 || stdout.String() != want {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	})
}

func TestParsePluginInventoryLiveSchemas(t *testing.T) {
	want := stableInventory(false, true)
	claude := `[{"id":"spacedock@spacedock","version":"0.27.1","enabled":false,"installPath":"/stable"},{"id":"spacedock@spacedock-edge","version":"0.28.0-pre0","enabled":true,"installPath":"/edge"}]`
	codex := `{"installed":[{"pluginId":"spacedock@spacedock","version":"0.27.1","installed":true,"enabled":false},{"pluginId":"spacedock@spacedock-edge","version":"0.28.0-pre0","installed":true,"enabled":true}]}`
	for _, fixture := range []struct{ host, data string }{{"claude", claude}, {"codex", codex}} {
		got, err := parsePluginInventory(fixture.host, []byte(fixture.data))
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("%s inventory=%#v err=%v, want %#v", fixture.host, got, err, want)
		}
	}
}

// testBinaryVersion is the deterministic binary version a handful of gating
// tests pin via withVersion when they need an exact, self-chosen relationship to
// a fixture (e.g. the upgrade-hint patch-skew cases). It shares its minor (19)
// with binaryMinor's fallback below.
const testBinaryVersion = "0.19.4"

// binaryMinor returns the (major, minor) pair of the CURRENT effective binary
// version (displayVersion() — Version, or its dev-embed resolution). The fixture
// helpers below derive their versions from it, so a "compatible" / "too-old-*"
// fixture is self-consistent with whatever this test binary reports, whether
// that is an explicit withVersion() stamp or the `dev` default's embedded
// checkout minor — no hardcoded fixture value can drift out of sync with it.
func binaryMinor(t *testing.T) (major, minor int) {
	t.Helper()
	major, minor, ok := contract.ParseMajorMinor(displayVersion())
	if !ok {
		t.Fatalf("displayVersion() %q has no parseable major.minor", displayVersion())
	}
	return major, minor
}

// writeVersionedManifest writes a minimal plugin manifest fixture at the given
// version under a fresh temp dir and returns its path.
func writeVersionedManifest(t *testing.T, version string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "plugin.json")
	body := `{ "name": "spacedock", "version": "` + version + `", "skills": "./skills/" }` + "\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// compatibleManifest returns a fixture path whose version shares the current
// binary's major.minor (patch 8, distinct from the binary's own patch to prove
// patch-skew tolerance).
func compatibleManifest(t *testing.T) string {
	t.Helper()
	major, minor := binaryMinor(t)
	return writeVersionedManifest(t, fmt.Sprintf("%d.%d.8", major, minor))
}

// tooOldBinaryManifest returns a fixture path one minor AHEAD of the current
// binary — the binary is too old for this plugin.
func tooOldBinaryManifest(t *testing.T) string {
	t.Helper()
	major, minor := binaryMinor(t)
	return writeVersionedManifest(t, fmt.Sprintf("%d.%d.0", major, minor+1))
}

// tooOldPluginManifest returns a fixture path one minor BEHIND the current
// binary — the plugin is too old for this binary, the D6 auto-heal verdict.
func tooOldPluginManifest(t *testing.T) string {
	t.Helper()
	major, minor := binaryMinor(t)
	if minor == 0 {
		t.Fatalf("binary minor is 0 — no room for a too-old-plugin fixture one minor behind")
	}
	return writeVersionedManifest(t, fmt.Sprintf("%d.%d.0", major, minor-1))
}

func withExecutablePath(t *testing.T, path string, err error) {
	t.Helper()
	orig := executablePath
	executablePath = func() (string, error) { return path, err }
	t.Cleanup(func() { executablePath = orig })
}

func TestResolvedDispatchLauncherPrefersExplicitSpacedockBin(t *testing.T) {
	dir := t.TempDir()
	candidate := filepath.Join(dir, "candidate spacedock")
	if err := os.WriteFile(candidate, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	withExecutablePath(t, filepath.Join(dir, "missing-current"), nil)

	got, ok := resolvedDispatchLauncher([]string{"SPACEDOCK_BIN=" + candidate})
	want, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || got != want {
		t.Fatalf("resolvedDispatchLauncher() = %q, %t; want explicit candidate %q", got, ok, want)
	}
}

func TestResolvedDispatchLauncherFallsBackFromInvalidSpacedockBin(t *testing.T) {
	dir := t.TempDir()
	current := filepath.Join(dir, "current-spacedock")
	if err := os.WriteFile(current, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	withExecutablePath(t, current, nil)

	got, ok := resolvedDispatchLauncher([]string{"SPACEDOCK_BIN=" + filepath.Join(dir, "missing")})
	want, err := filepath.EvalSymlinks(current)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || got != want {
		t.Fatalf("resolvedDispatchLauncher() = %q, %t; want current executable %q", got, ok, want)
	}
}

// withVersion stamps the package Version (the binary display semver the gate
// feeds to the upgrade-hint compare), restoring it after the test. The package
// default is the `dev` sentinel, which suppresses the hint, so a test that
// exercises the behind-plugin hint must stamp a real semver.
func withVersion(t *testing.T, v string) {
	t.Helper()
	orig := Version
	Version = v
	t.Cleanup(func() { Version = orig })
}

func executableFixture(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "spacedock")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	real, err := filepath.EvalSymlinks(p)
	if err != nil {
		t.Fatal(err)
	}
	return real
}

// TestSubprocessEnvScrubActive covers the truthy/falsy parsing of
// CLAUDE_CODE_SUBPROCESS_ENV_SCRUB that gates the launch warning (Task 2) and
// the doctor note (Task 3).
func TestSubprocessEnvScrubActive(t *testing.T) {
	cases := []struct {
		name string
		env  []string
		want bool
	}{
		{"unset", nil, false},
		{"empty value", []string{"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB="}, false},
		{"explicit zero", []string{"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=0"}, false},
		{"set to 1", []string{"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=1"}, true},
		{"set to true", []string{"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=true"}, true},
		{"unrelated env untouched", []string{"OTHER_VAR=1"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := subprocessEnvScrubActive(tc.env); got != tc.want {
				t.Fatalf("subprocessEnvScrubActive(%v) = %v, want %v", tc.env, got, tc.want)
			}
		})
	}
}

// TestClaudeFrontDoorLaunchesOnCompatible: on a compatible contract the front
// door invokes the launch seam with argv beginning `claude --agent
// spacedock:first-officer` and passes through the operator's trailing args.
func TestClaudeFrontDoorLaunchesOnCompatible(t *testing.T) {
	withVersion(t, testBinaryVersion)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "-p", "do the thing"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "-p", "do the thing", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// TestClaudeFrontDoorWarnsOnSubprocessEnvScrub: when
// CLAUDE_CODE_SUBPROCESS_ENV_SCRUB is set truthy and the operator has not
// declared --allowedTools themselves, spacedock claude prints its own
// attributed warning — on both the unsandboxed and --safehouse launch paths,
// since --dangerously-skip-permissions is not known to be exempt from Claude
// Code's hardening.
func TestClaudeFrontDoorWarnsOnSubprocessEnvScrub(t *testing.T) {
	cases := []struct {
		name string
		args []string // bare --safehouse forces the wrap path, matching TestClaudeForceSafehouseWrapsNoProfile
	}{
		{"unsandboxed", []string{"--", "-p", "do the thing"}},
		{"safehouse", []string{"--safehouse"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withVersion(t, testBinaryVersion)
			t.Setenv(subprocessEnvScrubEnv, "1")
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := runClaude(context.Background(), tc.args, t.TempDir(), fake, lookFound, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if !strings.Contains(stderr.String(), "CLAUDE_CODE_SUBPROCESS_ENV_SCRUB") {
				t.Fatalf("stderr missing env-scrub warning: %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), "--allowedTools") {
				t.Fatalf("stderr missing the --allowedTools remedy: %q", stderr.String())
			}
		})
	}
}

// TestClaudeFrontDoorSuppressesSubprocessEnvScrubWarning covers the two
// suppression cases: the var isn't truthy, or the operator already declared
// --allowedTools themselves (the documented workaround).
func TestClaudeFrontDoorSuppressesSubprocessEnvScrubWarning(t *testing.T) {
	cases := []struct {
		name string
		env  string // value to set subprocessEnvScrubEnv to ("" means leave unset)
		args []string
	}{
		{"var unset", "", []string{"--", "-p", "do the thing"}},
		{"var explicit zero", "0", []string{"--", "-p", "do the thing"}},
		{"operator declared allowedTools", "1", []string{"--", "--allowedTools", "Bash(git *)", "-p", "do the thing"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withVersion(t, testBinaryVersion)
			if tc.env == "" {
				t.Setenv(subprocessEnvScrubEnv, "") // register restore-on-cleanup, then unset for real
				os.Unsetenv(subprocessEnvScrubEnv)
			} else {
				t.Setenv(subprocessEnvScrubEnv, tc.env)
			}
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := runClaude(context.Background(), tc.args, t.TempDir(), fake, lookFound, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if strings.Contains(stderr.String(), "CLAUDE_CODE_SUBPROCESS_ENV_SCRUB") {
				t.Fatalf("stderr should not carry the env-scrub warning: %q", stderr.String())
			}
		})
	}
}

// TestFrontDoorUpgradeHintOnBehindPlugin is AC-4: the front-door gate prints the
// opt-in upgrade hint to stderr when the resolved plugin is contract-compatible
// but behind the binary, then proceeds to launch (the hint never blocks). The
// compatible.json fixture is plugin 0.19.8; stamping the binary Version to a
// strictly-newer semver makes it behind-but-compatible. A paired equal-version
// case (binary == plugin) asserts the gate stays silent. Both arms observe the
// recorded launch + stderr, never a source grep.
func TestFrontDoorUpgradeHintOnBehindPlugin(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" behind-plugin hints + launches", func(t *testing.T) {
			withVersion(t, "0.19.20") // same minor as the fixture plugin 0.19.8, strictly newer patch
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stderr bytes.Buffer

			code := tc.run(nil, t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (the hint must not block launch) (stderr=%q)", code, stderr.String())
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked — the hint must not change the launch path")
			}
			out := stderr.String()
			if !strings.Contains(out, "newer plugin") {
				t.Fatalf("stderr missing the behind-plugin upgrade hint: %q", out)
			}
			if !strings.Contains(out, "spacedock install --host "+tc.host) {
				t.Fatalf("stderr hint missing host install command for %s: %q", tc.host, out)
			}
		})

		t.Run(tc.name+" equal-version stays silent", func(t *testing.T) {
			withVersion(t, "0.19.8") // exactly the fixture plugin version
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stderr bytes.Buffer

			code := tc.run(nil, t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked on an equal-version compatible plugin")
			}
			if strings.Contains(stderr.String(), "newer plugin") || strings.Contains(stderr.String(), "spacedock install") {
				t.Fatalf("equal-version gate must stay silent (no upgrade hint): %q", stderr.String())
			}
		})
	}
}

// TestClaudeFrontDoorFailFastOnMismatch: on a mismatch verdict the launch seam is
// NOT invoked and the process exits non-zero with the pinned remedy on stderr.
func TestClaudeFrontDoorInjectsResolvedLauncherBin(t *testing.T) {
	withVersion(t, testBinaryVersion)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	t.Setenv(spacedockBinEnv, "/old/spacedock")
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	got, ok := envValueOf(fake.launchedEnv, spacedockBinEnv)
	if !ok || got != bin {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", spacedockBinEnv, got, ok, bin, fake.launchedEnv)
	}
}

func TestClaudeFrontDoorOmitsStaleLauncherBinWhenResolutionFails(t *testing.T) {
	withVersion(t, testBinaryVersion)
	withExecutablePath(t, "", errors.New("boom"))
	t.Setenv(spacedockBinEnv, "/old/spacedock")
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if got, ok := envValueOf(fake.launchedEnv, spacedockBinEnv); ok {
		t.Fatalf("%s in launch env = %q, want omitted", spacedockBinEnv, got)
	}
}

func TestClaudeFrontDoorLaunchEnvResolvesSymlink(t *testing.T) {
	withVersion(t, testBinaryVersion)
	real := executableFixture(t)
	link := filepath.Join(t.TempDir(), "spacedock-link")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	withExecutablePath(t, link, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if got, ok := envValueOf(fake.launchedEnv, spacedockBinEnv); !ok || got != real {
		t.Fatalf("%s = %q, %v; want symlink target %q, true", spacedockBinEnv, got, ok, real)
	}
}

// TestClaudeFrontDoorEnablesAgentTeamsWhenParentUnset (AC-1): with no parent
// value for CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS, the launched claude child env
// carries it set to 1 — a best-effort export, not the authoritative enabler (the
// authoritative enabler is ~/.claude/settings.json; see agentTeamsEnv).
func TestClaudeFrontDoorEnablesAgentTeamsWhenParentUnset(t *testing.T) {
	withVersion(t, testBinaryVersion)
	t.Setenv(agentTeamsEnv, "") // register restore-on-cleanup, then unset for real
	os.Unsetenv(agentTeamsEnv)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	got, ok := envValueOf(fake.launchedEnv, agentTeamsEnv)
	if !ok || got != "1" {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", agentTeamsEnv, got, ok, "1", fake.launchedEnv)
	}
}

// TestClaudeFrontDoorPreservesExplicitAgentTeams (AC-2): an explicit parent value
// for CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS is preserved, not overridden — an
// operator who set =0 keeps =0. The launcher respects the override either way.
func TestClaudeFrontDoorPreservesExplicitAgentTeams(t *testing.T) {
	withVersion(t, testBinaryVersion)
	t.Setenv(agentTeamsEnv, "0")
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	got, ok := envValueOf(fake.launchedEnv, agentTeamsEnv)
	if !ok || got != "0" {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", agentTeamsEnv, got, ok, "0", fake.launchedEnv)
	}
}

// TestClaudeFrontDoorFailFastOnMismatch (AC-2): a real version mismatch still
// fails fast even without --no-install — auto-install must NOT paper over an
// incompatibility. The verdict reaches runClaude's mismatch branch, not the
// no-plugin branch, so no install is invoked and no launch is reached.
func TestClaudeFrontDoorFailFastOnMismatch(t *testing.T) {
	withVersion(t, testBinaryVersion)
	fake := &fakeHost{manifest: tooOldBinaryManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero on mismatch")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("auto-install invoked over a mismatch: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("launch seam invoked on mismatch: %v", fake.launchedArg)
	}
	assertGateMismatchMessage(t, stderr.String())
}

// assertGateMismatchMessage is the gate-parity oracle: the live launch gate emits
// the same version-bearing mismatch message the doctor does — it names both the
// too-old-binary fixture's plugin display version (one minor ahead of the
// binary, patch 0 — see tooOldBinaryManifest) and the binary display version
// (Version), and carries no `contract N` token in the user-facing line.
func assertGateMismatchMessage(t *testing.T, out string) {
	t.Helper()
	major, minor := binaryMinor(t)
	wantPluginVersion := fmt.Sprintf("plugin %d.%d.0", major, minor+1)
	if !strings.Contains(out, wantPluginVersion) {
		t.Fatalf("gate mismatch message missing %q: %q", wantPluginVersion, out)
	}
	if !strings.Contains(out, "binary "+Version) {
		t.Fatalf("gate mismatch message missing binary version %q: %q", Version, out)
	}
	if regexp.MustCompile(`contract \d`).MatchString(out) {
		t.Fatalf("gate mismatch message must not carry a contract-N token: %q", out)
	}
}

// TestClaudeFrontDoorNoPluginAutoInstalls (AC-1a): with no installed plugin and
// no flags the front door auto-installs the plugin then launches a working
// session — the single command the user typed yields a working FO session
// rather than refusing. Install-invocation and launch-reached are observed
// behaviors recorded by the stub, not a string match.
func TestClaudeFrontDoorNoPluginAutoInstalls(t *testing.T) {
	fake := &fakeHost{manifest: "", manifestAfterInstall: compatibleManifest(t)} // no plugin found; install lands a compatible one
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 when no plugin → auto-install + launch (stderr=%q)", code, stderr.String())
	}
	if len(fake.installCmds) == 0 {
		t.Fatalf("install seam not invoked: auto-install did not run")
	}
	if fake.launchedArg == nil {
		t.Fatalf("launch seam not invoked after auto-install")
	}
}

// TestClaudeFrontDoorNoPluginNoInstallRefuses (AC-1b): --no-install preserves
// the old refuse-and-instruct behavior — no install, no launch, non-zero exit,
// with the actionable instruct message on stderr.
func TestClaudeFrontDoorNoPluginNoInstallRefuses(t *testing.T) {
	fake := &fakeHost{manifest: ""} // no plugin found
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero with --no-install and no plugin")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), "spacedock install") {
		t.Fatalf("stderr missing actionable instruct message: %q", stderr.String())
	}
}

// TestClaudeFrontDoorNonEmptyMissingManifestAutoInstalls: a host that reports a
// non-empty installPath to a directory LACKING the plugin manifest is the
// NoPluginFound verdict from a phantom installPath — the plugin genuinely is not
// on disk, so the default auto-installs (a deliberate behavior change recorded in
// ideation). --no-install preserves the refuse-and-instruct arm.
func TestClaudeFrontDoorNonEmptyMissingManifestAutoInstalls(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".claude-plugin", "plugin.json")

	t.Run("default auto-installs", func(t *testing.T) {
		fake := &fakeHost{manifest: missing, manifestAfterInstall: compatibleManifest(t)} // non-empty path, but the file is absent; install lands a compatible one
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 when phantom manifest → auto-install + launch (stderr=%q)", code, stderr.String())
		}
		if len(fake.installCmds) == 0 {
			t.Fatalf("install seam not invoked for a phantom (missing) manifest")
		}
		if fake.launchedArg == nil {
			t.Fatalf("launch seam not invoked after auto-install for a phantom manifest")
		}
	})

	t.Run("--no-install refuses", func(t *testing.T) {
		fake := &fakeHost{manifest: missing}
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

		if code == 0 {
			t.Fatalf("exit = 0, want non-zero with --no-install and a phantom manifest")
		}
		if len(fake.installCmds) != 0 {
			t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
		}
		if fake.launchedArg != nil {
			t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
		}
	})
}

// TestGateRemedyNamesLiveInstallCommand: every remedy must point at a command
// the binary actually recognizes. After the init->install rename a user who hits
// the gate and is told to "run spacedock init --host …" runs a command that now
// exits 2 (unknown command). gateHost owns only the always-fail-fast remedies
// (resolve error → MalformedVersion); the NoPluginFound message is the caller's
// (it auto-installs by default, refuses under --no-install), so the no-plugin /
// missing-manifest remedies are asserted on the launcher's --no-install output.
// Each remedy must name `spacedock install` and never `spacedock init`.
func TestGateRemedyNamesLiveInstallCommand(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".claude-plugin", "plugin.json")

	// gateHost owns the resolve-error remedy (a host-CLI failure is a hard fail,
	// not a missing plugin). It MUST print here.
	t.Run("resolve error (gateHost-owned)", func(t *testing.T) {
		var stderr bytes.Buffer
		if res := gateHost(&fakeHost{resolveErr: errors.New("host CLI failed")}, "claude", &stderr); res.Verdict == contract.Compatible {
			t.Fatalf("gateHost = Compatible, want denied")
		}
		assertRemedyNamesInstall(t, stderr.String())
	})

	// The NoPluginFound remedies (no plugin, phantom manifest) are caller-owned:
	// gateHost no longer prints for them, so assert the remedy on the launcher's
	// --no-install refuse output.
	noPluginCases := []struct {
		name     string
		manifest string
	}{
		{"no plugin", ""},
		{"missing manifest", missing},
	}
	for _, tc := range noPluginCases {
		t.Run(tc.name+" (caller-owned via --no-install)", func(t *testing.T) {
			fake := &fakeHost{manifest: tc.manifest}
			var stdout, stderr bytes.Buffer
			if code := runClaude(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr); code == 0 {
				t.Fatalf("exit = 0, want non-zero with --no-install and no plugin")
			}
			assertRemedyNamesInstall(t, stderr.String())
		})
		t.Run(tc.name+" (gateHost stays silent)", func(t *testing.T) {
			var stderr bytes.Buffer
			if res := gateHost(&fakeHost{manifest: tc.manifest}, "claude", &stderr); res.Verdict != contract.NoPluginFound {
				t.Fatalf("gateHost verdict = %v, want NoPluginFound", res.Verdict)
			}
			if stderr.Len() != 0 {
				t.Fatalf("gateHost printed for NoPluginFound (caller owns the message): %q", stderr.String())
			}
		})
	}

	// The command the remedy names must resolve in the live command tree. cobra's
	// Find returns the matched command for a registered name and falls back to the
	// root for an unknown one — so `install` must resolve to a non-root command
	// while the removed `init` must fall back to root (the unknown-command path
	// that exits 2). Resolution is deterministic and touches no host CLI.
	var stdout, stderr bytes.Buffer
	root := newRootCommand(context.Background(), nil, nil, t.TempDir(), nil, &stdout, &stderr, &fakeRunner{}, nil)
	if cmd, _, err := root.Find([]string{"install"}); err != nil || cmd == root {
		t.Fatalf("`install` did not resolve to a registered command (cmd=%v, err=%v)", cmd.Name(), err)
	}
	if cmd, _, _ := root.Find([]string{"init"}); cmd != root {
		t.Fatalf("`init` resolved to a command (%v) — the removed verb must fall back to the unknown-command path", cmd.Name())
	}
}

// TestGateHostStaysSilentForTooOldPlugin mirrors the existing NoPluginFound
// silence test for D6's second healable verdict: gateHost must NOT print the
// mismatch remedy for too-old-plugin (the caller decides between the auto-heal
// announcement and the --no-install remedy; a print here would show the
// operator a scary message right before the silent heal).
func TestGateHostStaysSilentForTooOldPlugin(t *testing.T) {
	var stderr bytes.Buffer
	res := gateHost(&fakeHost{manifest: tooOldPluginManifest(t)}, "claude", &stderr)
	if res.Verdict != contract.TooOldPlugin {
		t.Fatalf("gateHost verdict = %v, want TooOldPlugin", res.Verdict)
	}
	if stderr.Len() != 0 {
		t.Fatalf("gateHost printed for TooOldPlugin (the caller owns the message, D6): %q", stderr.String())
	}
}

// TestFrontDoorTooOldPluginAutoRefreshes is D6's default-arm proof for both
// front doors: a too-old-plugin verdict announces "Refreshing the {host}
// plugin…" (not "Installing…" — no plugin was missing, one is just behind),
// installs, re-gates ONCE, and on success proceeds to launch — mirroring the
// NoPluginFound auto-install arm but for the SECOND healable verdict.
func TestFrontDoorTooOldPluginAutoRefreshes(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: tooOldPluginManifest(t), manifestAfterInstall: compatibleManifest(t)}
			var stderr bytes.Buffer

			code := tc.run(nil, t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (auto-refresh + re-gate + launch) (stderr=%q)", code, stderr.String())
			}
			want := "Refreshing the " + tc.host + " plugin"
			if !strings.Contains(stderr.String(), want) {
				t.Fatalf("stderr missing %q announcement: %q", want, stderr.String())
			}
			if strings.Contains(stderr.String(), "Spacedock version mismatch") {
				t.Fatalf("stderr shows the scary mismatch remedy before the silent heal (D6 forbids it): %q", stderr.String())
			}
			if len(fake.installCmds) == 0 {
				t.Fatalf("install seam not invoked: auto-refresh did not run")
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked after auto-refresh")
			}
		})
	}
}

// TestFrontDoorTooOldPluginNoInstallRefuses is D6's --no-install refuse arm:
// with a too-old-plugin verdict and --no-install, the launcher refuses with the
// gate's own mismatch remedy (not the no-plugin remedy — a stale plugin is a
// different situation), installs nothing, and never launches.
func TestFrontDoorTooOldPluginNoInstallRefuses(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: tooOldPluginManifest(t)}
			var stderr bytes.Buffer

			code := tc.run([]string{"--no-install"}, t.TempDir(), fake, &stderr)

			if code == 0 {
				t.Fatalf("exit = 0, want non-zero with --no-install on a too-old-plugin verdict")
			}
			if !strings.Contains(stderr.String(), "Update the plugin to continue.") {
				t.Fatalf("stderr missing the too-old-plugin remedy: %q", stderr.String())
			}
			if len(fake.installCmds) != 0 {
				t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
			}
			if fake.launchedArg != nil {
				t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
			}
		})
	}
}

// TestFrontDoorTooOldPluginSecondMissRefuses is D6's "one retry, no loop"
// guarantee: when the auto-refresh install does NOT fix the resolution (a
// release-window race, a stale channel — modeled here by a fakeHost whose
// manifest never changes after Install), the launcher refuses rather than
// launching blind, and Install is called exactly ONCE (no retry loop).
func TestFrontDoorTooOldPluginSecondMissRefuses(t *testing.T) {
	fake := &fakeHost{manifest: tooOldPluginManifest(t)} // manifestAfterInstall unset: install does not fix it
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero when the refresh does not fix the mismatch")
	}
	if len(fake.installCmds) != 3 { // one Install call records {host, source, branch}
		t.Fatalf("install seam called %v, want exactly one Install invocation (no retry loop)", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("launch seam invoked despite the persisted mismatch: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), "Update the plugin to continue.") {
		t.Fatalf("stderr missing the second-miss remedy: %q", stderr.String())
	}
}

// TestNoPluginAutoInstallAnnouncesHostCorrectly (AC-A, auto-install arm): with
// no installed plugin the launcher announces `Installing the {host} plugin…` on
// stderr before installing, and a codex run never names `claude` (the old
// gateHost remedy hardcoded a `spacedock claude` hint that was wrong in a codex
// run). Install + launch are observed via the recorded seams, not a string
// match.
func TestNoPluginAutoInstallAnnouncesHostCorrectly(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: "", manifestAfterInstall: compatibleManifest(t)} // no plugin found; install lands a compatible one
			var stderr bytes.Buffer

			code := tc.run(nil, t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (auto-install + launch) (stderr=%q)", code, stderr.String())
			}
			want := "Installing the " + tc.host + " plugin"
			if !strings.Contains(stderr.String(), want) {
				t.Fatalf("stderr missing %q announcement: %q", want, stderr.String())
			}
			if tc.host == "codex" && strings.Contains(stderr.String(), "spacedock claude") {
				t.Fatalf("codex auto-install stderr names claude: %q", stderr.String())
			}
			if len(fake.installCmds) == 0 {
				t.Fatalf("install seam not invoked: auto-install did not run")
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked after auto-install")
			}
		})
	}
}

// TestNoPluginNoInstallRemedyIsHostCorrect (AC-A, refuse arm): with --no-install
// and no plugin the launcher prints the manual remedy naming the host-correct
// install/bootstrap commands (`spacedock install --host {host}`, `spacedock
// {host} --skip-compat-check`) and NEVER a `claude` hint in a codex run, then
// exits non-zero without installing or launching. This is the message the caller
// now owns (gateHost stopped printing for NoPluginFound).
func TestNoPluginNoInstallRemedyIsHostCorrect(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: ""} // no plugin found
			var stderr bytes.Buffer

			code := tc.run([]string{"--no-install"}, t.TempDir(), fake, &stderr)

			if code == 0 {
				t.Fatalf("exit = 0, want non-zero with --no-install and no plugin")
			}
			out := stderr.String()
			if !strings.Contains(out, "spacedock install --host "+tc.host) {
				t.Fatalf("remedy missing host-correct install command for %s: %q", tc.host, out)
			}
			if !strings.Contains(out, "spacedock "+tc.host+" --skip-compat-check") {
				t.Fatalf("remedy missing host-correct bootstrap command for %s: %q", tc.host, out)
			}
			if tc.host == "codex" && strings.Contains(out, "claude") {
				t.Fatalf("codex --no-install remedy names claude: %q", out)
			}
			if len(fake.installCmds) != 0 {
				t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
			}
			if fake.launchedArg != nil {
				t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
			}
		})
	}
}

// assertRemedyNamesInstall: every gate / refuse remedy must name the live
// `spacedock install` command and never the removed `spacedock init` (which now
// exits 2).
func assertRemedyNamesInstall(t *testing.T, remedy string) {
	t.Helper()
	if !strings.Contains(remedy, "spacedock install") {
		t.Fatalf("remedy does not name the live install command: %q", remedy)
	}
	if strings.Contains(remedy, "spacedock init") {
		t.Fatalf("remedy names the removed init command (exits 2): %q", remedy)
	}
}

// TestClaudeFrontDoorSkipCompatCheckBootstrap: the --skip-compat-check
// override launches without resolving the manifest (bootstrap case where the
// plugin is being installed for the first time).
func TestClaudeFrontDoorSkipCompatCheckBootstrap(t *testing.T) {
	fake := &fakeHost{manifest: tooOldBinaryManifest(t)} // would mismatch if checked
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--skip-compat-check"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 with --skip-compat-check (stderr=%q)", code, stderr.String())
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v (skip-check must not pass the flag through)", fake.launchedArg, want)
	}
}

// TestCodexFrontDoorLaunchesOnCompatible: on a compatible contract the Codex
// front door preserves post-fence argv order and retains fresh-launch posture
// when no exact `resume` token is present.
func TestCodexFrontDoorLaunchesOnCompatible(t *testing.T) {
	withVersion(t, testBinaryVersion)
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "-m", "gpt-x"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"},
		wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", "-m", "gpt-x", wantCodexBootstrapPrompt)...)
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// TestCodexFrontDoorFailFastOnMismatch: codex fails fast on a mismatch verdict
// with the pinned remedy and does NOT launch.
func TestCodexFrontDoorInjectsLauncherBinThroughSafehouseResume(t *testing.T) {
	withVersion(t, testBinaryVersion)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	t.Setenv(spacedockBinEnv, "/old/spacedock")
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "resume", "abc123"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	wantArgv := append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"},
		wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", "resume", "abc123")...)
	if !equalArgv(fake.launchedArg, wantArgv) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, wantArgv)
	}
	got, ok := envValueOf(fake.launchedEnv, spacedockBinEnv)
	if !ok || got != bin {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", spacedockBinEnv, got, ok, bin, fake.launchedEnv)
	}
}

func TestCodexFrontDoorForwardsPostFenceExampleThroughSafehouse(t *testing.T) {
	withVersion(t, testBinaryVersion)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "--model", "gpt-x", "resume", "abc123"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	wantArgv := append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"},
		wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", "--model", "gpt-x", "resume", "abc123")...)
	if !equalArgv(fake.launchedArg, wantArgv) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, wantArgv)
	}
	if stderr.Len() != 0 {
		t.Fatalf("resume launch produced Spacedock output: %q", stderr.String())
	}
}

func TestCodexFrontDoorBootstrapsPostFenceWithoutResumeThroughSafehouse(t *testing.T) {
	withVersion(t, testBinaryVersion)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	args := []string{"--", "--future-codex-flag=handoff", "opaque-argument"}
	code := runCodex(context.Background(), args, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	wantArgv := append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"},
		wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", "--future-codex-flag=handoff", "opaque-argument", wantCodexBootstrapPrompt)...)
	if !equalArgv(fake.launchedArg, wantArgv) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, wantArgv)
	}
	if !strings.Contains(stderr.String(), "\u00b7 launching codex") {
		t.Fatalf("fresh post-fence launch omitted the Codex banner: %q", stderr.String())
	}
}

// TestCodexFrontDoorDoesNotEnableAgentTeams: the agent-teams flag is a
// claude-only concern; the codex launch path must NOT inject it (scope guard).
func TestCodexFrontDoorDoesNotEnableAgentTeams(t *testing.T) {
	withVersion(t, testBinaryVersion)
	t.Setenv(agentTeamsEnv, "") // register restore-on-cleanup, then unset for real
	os.Unsetenv(agentTeamsEnv)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if got, ok := envValueOf(fake.launchedEnv, agentTeamsEnv); ok {
		t.Fatalf("%s in codex launch env = %q, want omitted", agentTeamsEnv, got)
	}
}

func TestCodexFrontDoorFailFastOnMismatch(t *testing.T) {
	withVersion(t, testBinaryVersion)
	fake := &fakeHost{manifest: tooOldBinaryManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero on mismatch")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("install seam invoked on codex mismatch: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("launch seam invoked on mismatch: %v", fake.launchedArg)
	}
	assertGateMismatchMessage(t, stderr.String())
}

// TestCodexFrontDoorNoPluginAutoInstalls (AC-2): with no installed codex plugin
// and no flags the front door auto-installs the plugin then launches a working
// session — the single command the user typed yields a working FO session rather
// than refusing. Both no-plugin verdicts (an empty manifest AND a phantom
// installPath to an absent manifest) auto-install, matching the claude branch.
// The recorded install seam ({host, source, branch}) and the captured launch argv
// are the independent sources of truth.
func TestCodexFrontDoorNoPluginAutoInstalls(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".codex-plugin", "plugin.json")
	cases := []struct {
		name     string
		manifest string
	}{
		{name: "empty manifest", manifest: ""},
		{name: "missing manifest", manifest: missing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: tc.manifest, manifestAfterInstall: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 when codex has no plugin → auto-install + launch (stderr=%q)", code, stderr.String())
			}
			wantInstall := []string{"codex", channelMarketplaceSource(devBranch), devBranch}
			if !equalArgv(fake.installCmds, wantInstall) {
				t.Fatalf("install seam = %v, want %v (the {host, source, branch} seam)", fake.installCmds, wantInstall)
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked after codex auto-install")
			}
		})
	}
}

// TestCodexFrontDoorNoPluginNoInstallRefuses (AC-3): --no-install preserves the
// refuse-and-instruct behavior — no install, no launch, non-zero exit, with the
// codex no-plugin remedy on stderr. This is the migrated assertion from the old
// fail-fast test, now gated behind --no-install. Both no-plugin verdicts refuse.
func TestCodexFrontDoorNoPluginNoInstallRefuses(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".codex-plugin", "plugin.json")
	cases := []struct {
		name     string
		manifest string
	}{
		{name: "empty manifest", manifest: ""},
		{name: "missing manifest", manifest: missing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: tc.manifest}
			var stdout, stderr bytes.Buffer

			code := runCodex(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

			if code == 0 {
				t.Fatalf("exit = 0, want non-zero with --no-install and no codex plugin")
			}
			if len(fake.installCmds) != 0 {
				t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
			}
			if fake.launchedArg != nil {
				t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
			}
			if !strings.Contains(stderr.String(), "codex plugin") {
				t.Fatalf("stderr missing codex no-plugin remedy: %q", stderr.String())
			}
		})
	}
}

func equalArgv(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
