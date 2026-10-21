package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	piTestListing = `{"plugins":[
 {"name":"spacedock","source":{"source":"url","url":"https://github.com/spacedock-dev/spacedock.git","ref":"stable"}},
 {"name":"subspace","source":{"source":"url","url":"https://github.com/spacedock-dev/subspace.git","ref":"v9.9.9"}},
 {"name":"compass","source":{"source":"git-subdir","url":"https://github.com/spacedock-dev/compass.git","path":"compass","ref":"stable"}}]}`
	piTestSubspace = "git:github.com/spacedock-dev/subspace"
	piTestAlpha    = "git:github.com/spacedock-dev/alpha"
)

func installPiForTest(t *testing.T, ops *fakePiRuntimeOps) (code int, stdout string) {
	t.Helper()
	repo, pkg := t.TempDir(), t.TempDir()
	writePiSkillFixtures(t, repo)
	writePiSubagentsFixtures(t, pkg)
	ops.lookPath, ops.statOK, ops.packageStatus = piHealthyPathFixtures(), statOKForPiResources(repo, pkg), healthyPiPackageStatus()
	var out, errOut bytes.Buffer
	code = runInitWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, ops, append(piTestEnv(pkg, t.TempDir()), "SPACEDOCK_REPO_ROOT="+repo), &out, &errOut)
	return code, out.String()
}

func TestPiInstallPinsBothPackagesRemovingEachBeforeItsInstall(t *testing.T) {
	savedBranch := devBranch
	defer func() { devBranch = savedBranch }()
	for _, tc := range []struct{ branch, ref string }{{"main", "@stable"}, {"next", "@main"}} {
		devBranch = tc.branch
		ops := &fakePiRuntimeOps{listing: []byte(piTestListing)}
		if code, out := installPiForTest(t, ops); code != 0 {
			t.Fatalf("%s: exit=%d\n%s", tc.branch, code, out)
		}
		want := "remove " + piSpacedockPackageSource + "\ninstall " + piSpacedockPackageSource + tc.ref +
			"\nremove " + piTestSubspace + "\ninstall " + piTestSubspace + "@v9.9.9"
		if got := strings.Join(ops.piCalls, "\n"); got != want {
			t.Errorf("%s: pi calls =\n%s\nwant\n%s", tc.branch, got, want)
		}
	}
}

func TestPiInstallSkipsAnExtensionWithoutAPiPackage(t *testing.T) {
	saved := extensions
	defer func() { extensions = saved }()
	extensions = []extension{{name: "subspace", program: "spacedock-review"}}
	ops := &fakePiRuntimeOps{listing: []byte(piTestListing)}
	installPiForTest(t, ops)
	if len(ops.piCalls) != 2 {
		t.Errorf("pi calls = %v, want only the Spacedock package's remove and install", ops.piCalls)
	}
}

func TestPiInstallNamesAFailedExtensionPackageAndExitsZero(t *testing.T) {
	for name, tc := range map[string]struct {
		ops  fakePiRuntimeOps
		want string
	}{
		"unreadable listing": {fakePiRuntimeOps{}, "could not read the extension listing"},
		"no subspace entry":  {fakePiRuntimeOps{listing: []byte(`{"plugins":[]}`)}, "has no subspace entry"},
		"entry without ref":  {fakePiRuntimeOps{listing: []byte(`{"plugins":[{"name":"subspace","source":{"source":"url","url":"https://github.com/spacedock-dev/subspace.git"}}]}`)}, "does not give subspace as a git url with a ref"},
		"pi install fails":   {fakePiRuntimeOps{listing: []byte(piTestListing), piFailing: piTestSubspace + "@v9.9.9"}, "pi install \"" + piTestSubspace + "@v9.9.9\" failed"},
	} {
		ops := tc.ops
		code, out := installPiForTest(t, &ops)
		line := "spacedock install: "
		if code != 0 || !strings.Contains(out, line) || !strings.Contains(out, tc.want) {
			t.Errorf("%s: exit=%d, want 0 and a %q line containing %q:\n%s", name, code, line, tc.want, out)
		}
	}
}

func TestPiDoctorListsTheSubspacePackageWithoutChangingTheExit(t *testing.T) {
	pkg, home := t.TempDir(), t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	statOK := statOKForPiResources(t.TempDir(), pkg)
	statOK[filepath.Join(home, ".pi", "agent", "auth.json")] = true
	doctor := func() (int, string) {
		var stdout, stderr bytes.Buffer
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath: piHealthyPathFixtures(), statOK: statOK, packageStatus: healthyPiPackageStatus(),
		}, piTestEnv(pkg, home), &stdout, &stderr)
		return code, stdout.String()
	}

	missingCode, missing := doctor()
	if !strings.Contains(missing, "MISSING Subspace package") || !strings.Contains(missing, "remedy: run `spacedock install --host pi`") {
		t.Errorf("doctor without the package lacks the MISSING row and its remedy:\n%s", missing)
	}
	settings := filepath.Join(home, ".pi", "agent", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settings), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, settings, `{"packages":["`+piTestSubspace+`@v9.9.9"]}`)
	okCode, ok := doctor()
	if want := "OK Subspace package: " + filepath.Join(home, ".pi", "agent", "git", "github.com", "spacedock-dev", "subspace") + " (v9.9.9)"; !strings.Contains(ok, want) {
		t.Errorf("doctor with the package lacks %q:\n%s", want, ok)
	}
	if missingCode != 0 || okCode != 0 {
		t.Errorf("doctor exits = %d without, %d with the package, want 0 both: the row is information", missingCode, okCode)
	}
}

type piUpdateFixture struct {
	deps     updateDeps
	settings string
	ran      []string
}

func newPiUpdateFixture(t *testing.T, status piPackageStatus) *piUpdateFixture {
	t.Helper()
	s := stageInstall(t)
	extensions[0].piTitle = "Alpha"
	script := filepath.Join(t.TempDir(), "install.sh")
	mustWrite(t, script, "#!/bin/sh\n:\n")
	agentDir := t.TempDir()
	f := &piUpdateFixture{deps: updateDepsFor(s, s.host, script), settings: filepath.Join(agentDir, "settings.json")}
	writeFile(t, f.settings, `{"packages":["`+piSpacedockPackageSource+`@stable","`+piTestAlpha+`@v1"]}`)
	f.deps.env = append(f.deps.env, "PI_CODING_AGENT_DIR="+agentDir)
	f.deps.lookPath = hostsOnPath("pi")
	f.deps.piOps = &fakePiRuntimeOps{packageStatus: status}
	f.deps.runSelf = func(args ...string) (string, error) {
		f.ran = append(f.ran, strings.Join(args, " "))
		writeFile(t, f.settings, `{"packages":["`+piSpacedockPackageSource+`@stable","`+piTestAlpha+`@v2"]}`)
		return "noise\nspacedock install: pi install \"" + piTestAlpha + "@v2\" failed: exit status 1\n", nil
	}
	return f
}

func TestUpdateRefreshesThePiPackagesWithTheNewProgramAndShowsTheirRefs(t *testing.T) {
	f := newPiUpdateFixture(t, piPackageStatus{registered: true, source: piSpacedockPackageSource + "@stable"})
	var stdout, stderr bytes.Buffer
	if code := runUpdate(f.deps, &stdout, &stderr); code != 0 {
		t.Fatalf("update exited %d: %s", code, stderr.String())
	}
	if strings.Join(f.ran, ";") != "install --host pi" {
		t.Errorf("runSelf calls = %v, want one install --host pi", f.ran)
	}
	out := stdout.String()
	for _, want := range []string{"alpha package      v1 -> v2", "spacedock package  stable          unchanged", "Open Pi sessions keep the old skills until you run /reload.", "spacedock install: pi install"} {
		if !strings.Contains(out, want) {
			t.Errorf("update output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "noise") {
		t.Errorf("the new program's other output was printed:\n%s", out)
	}
}

func TestUpdateLeavesPiAloneWhenItCannotRefreshThePublishedPackage(t *testing.T) {
	for name, tc := range map[string]struct {
		status piPackageStatus
		onPath bool
		want   string
	}{
		"pi not on PATH":       {piPackageStatus{registered: true, source: piSpacedockPackageSource}, false, "pi is not on PATH"},
		"no Spacedock package": {piPackageStatus{}, true, "pi has no Spacedock package"},
		"local checkout":       {piPackageStatus{registered: true, source: "/work/spacedock"}, true, "is /work/spacedock, not a published source"},
	} {
		f := newPiUpdateFixture(t, tc.status)
		if !tc.onPath {
			f.deps.lookPath = hostsOnPath()
		}
		var stdout, stderr bytes.Buffer
		if code := runUpdate(f.deps, &stdout, &stderr); code != 0 || len(f.ran) != 0 || !strings.Contains(stdout.String(), tc.want) {
			t.Errorf("%s: exit=%d, runSelf calls=%v, want 0, none and %q:\n%s", name, code, f.ran, tc.want, stdout.String())
		}
	}
}

func TestUpdateOnHomebrewRefreshesPiInProcessAndFailsWhenItFails(t *testing.T) {
	f := newPiUpdateFixture(t, piPackageStatus{registered: true, source: piSpacedockPackageSource + "@stable"})
	f.deps.selfPath = "/opt/homebrew/Caskroom/spacedock/0.28.0/spacedock"
	f.deps.runSelf = func(...string) (string, error) {
		t.Fatal("Homebrew has not replaced the program; update must not run it")
		return "", nil
	}
	ops := f.deps.piOps.(*fakePiRuntimeOps)
	ops.piFailing = piSpacedockPackagePin(devBranch)
	var stdout, stderr bytes.Buffer
	code := runUpdate(f.deps, &stdout, &stderr)
	if len(ops.piInstalls) != 1 || code != 1 || !strings.Contains(stderr.String(), "the Pi package update failed") {
		t.Errorf("exit=%d, pi installs=%v, stderr=%q, want one in-process install and exit 1", code, ops.piInstalls, stderr.String())
	}
}

func TestReadExtensionListingReadsADirectoryAndRefusesAnythingElse(t *testing.T) {
	dir := t.TempDir()
	writeFileWithDirs(t, filepath.Join(dir, ".claude-plugin", "marketplace.json"), piTestListing)
	if got, err := (execPiRuntimeOps{}).ReadExtensionListing(dir); err != nil || string(got) != piTestListing {
		t.Errorf("directory listing = %q, %v, want the file's bytes", got, err)
	}
	if _, err := (execPiRuntimeOps{}).ReadExtensionListing(filepath.Join(dir, "missing")); err == nil {
		t.Error("a source that is neither a directory nor owner/repo was accepted")
	}
}
