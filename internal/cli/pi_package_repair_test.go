package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

// withPiReleaseIdentity overrides the binary's release identity (linker stamp
// + build-info main version) and restores both derivation seams on cleanup.
// stamp uses the goreleaser form (no "v" prefix) or the "dev" sentinel;
// buildInfoVersion models debug.ReadBuildInfo's main module version ("v0.27.2"
// for a proxy install, "(devel)" for a checkout build).
func withPiReleaseIdentity(t *testing.T, stamp, buildInfoVersion string) {
	t.Helper()
	savedVersion, savedBuildInfo := Version, piBuildInfoMainVersion
	Version = stamp
	piBuildInfoMainVersion = func() string { return buildInfoVersion }
	t.Cleanup(func() {
		Version, piBuildInfoMainVersion = savedVersion, savedBuildInfo
	})
}

// piRepairFixture creates the repo + pi-subagents package roots and their
// stat fixtures, so a runPi launch reaches the launch seam when the package
// gate passes.
func piRepairFixture(t *testing.T) (repo, pkg string) {
	t.Helper()
	repo, pkg = t.TempDir(), t.TempDir()
	writePiSkillFixtures(t, repo)
	writePiSubagentsFixtures(t, pkg)
	return repo, pkg
}

// piRepairOps builds the healthy-otherwise fake ops with a canned package
// status; after, when non-nil, replaces the status after the first PiInstall
// (the model for "the repair ran; here is the recheck result").
func piRepairOps(repo, pkg string, status piPackageStatus, after *piPackageStatus) *fakePiRuntimeOps {
	return &fakePiRuntimeOps{
		lookPath:           piHealthyPathFixtures(),
		statOK:             statOKForPiResources(repo, pkg),
		packageStatus:      status,
		statusAfterInstall: after,
	}
}

// runPiWith runs the front door with the standard temp cwd/env and the given
// extra args; returns (exit, stdout, stderr).
func runPiWith(t *testing.T, ops *fakePiRuntimeOps, pkg string, extraArgs []string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	args := append([]string{"review this"}, extraArgs...)
	code := runPi(context.Background(), args, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestPiReleaseRefDerivation(t *testing.T) {
	const pseudo1 = "v0.28.0-pre0.0.20260828120000-abcdef123456" // @next over a pre tag
	const pseudo2 = "v0.27.1-0.20260828120000-abcdef123456"      // no base tag
	for _, tc := range []struct {
		name, stamp, buildInfo, want string
	}{
		{"linker stamp wins", "0.27.2", "v0.26.0", "v0.27.2"},
		{"v-prefixed stamp normalizes", "v0.27.2", "", "v0.27.2"},
		{"prerelease stamp", "0.28.0-pre1", "(devel)", "v0.28.0-pre1"},
		{"proxy build pins its tag", "dev", "v0.27.2", "v0.27.2"},
		{"checkout build is the dev sentinel", "dev", "(devel)", ""},
		{"empty build info is the dev sentinel", "dev", "", ""},
		{"next pseudo-version floats", "dev", pseudo1, ""},
		{"plain pseudo-version floats", "dev", pseudo2, ""},
		{"non-semver build info floats", "dev", "not-a-version", ""},
		{"empty stamp uses build info", "", "v0.27.2", "v0.27.2"},
	} {
		if got := piReleaseRefFrom(tc.stamp, tc.buildInfo); got != tc.want {
			t.Fatalf("%s: piReleaseRefFrom(%q,%q)=%q want %q", tc.name, tc.stamp, tc.buildInfo, got, tc.want)
		}
	}
}

func TestPiPinnedSourceAndRefParsing(t *testing.T) {
	const repo = "git:github.com/spacedock-dev/spacedock"
	for _, tc := range []struct{ ref, want string }{
		{"v0.27.2", repo + "@v0.27.2"},
		{"v0.28.0-pre0", repo + "@v0.28.0-pre0"},
		{"", repo}, // dev sentinel keeps the bare floating source
	} {
		if got := piPinnedSource(tc.ref); got != tc.want {
			t.Fatalf("piPinnedSource(%q)=%q want %q", tc.ref, got, tc.want)
		}
	}
	for _, tc := range []struct {
		source, ref string
		isGit       bool
	}{
		{"git:github.com/spacedock-dev/spacedock", "", true},
		{"git:github.com/spacedock-dev/spacedock@v0.27.2", "v0.27.2", true},
		{"file:/tmp/repo", "", false},
		{"npm:spacedock", "", false},
		{"/usr/local/spacedock", "", false},
		{"", "", false},
	} {
		if ref, ok := piGitSourceRef(tc.source); ref != tc.ref || ok != tc.isGit {
			t.Fatalf("piGitSourceRef(%q)=(%q,%v) want (%q,%v)", tc.source, ref, ok, tc.ref, tc.isGit)
		}
	}
}

func TestPiPackageNeedsRepairTable(t *testing.T) {
	missing := piPackageStatus{}
	unpinned := piPackageStatus{registered: true, source: "git:github.com/spacedock-dev/spacedock"}
	wrongLine := piPackageStatus{registered: true, source: "git:github.com/spacedock-dev/spacedock@v0.28.0-pre1"}
	ownRef := piPackageStatus{registered: true, source: "git:github.com/spacedock-dev/spacedock@v0.27.2"}
	userEntry := piPackageStatus{registered: true, source: "file:/x"}
	for _, tc := range []struct {
		name       string
		st         piPackageStatus
		relRef     string
		want       bool
		wantReason string
	}{
		{"missing entry, release-shaped", missing, "v0.27.2", true, "missing"},
		{"unpinned git entry", unpinned, "v0.27.2", true, "unpinned"},
		{"wrong release line", wrongLine, "v0.27.2", true, "wrong release line"},
		{"pinned to own ref", ownRef, "v0.27.2", false, ""},
		{"non-git entry is user-managed", userEntry, "v0.27.2", false, ""},
		{"dev sentinel never repairs", missing, "", false, ""},
	} {
		got, reason := piPackageNeedsRepair(tc.st, tc.relRef)
		if got != tc.want || (got && reason != tc.wantReason) {
			t.Fatalf("%s: repair=(%v,%q) want (%v,%q)", tc.name, got, reason, tc.want, tc.wantReason)
		}
	}
}

func TestPiInstallUsesReleasePinnedSourcePerIdentity(t *testing.T) {
	pkg := t.TempDir()
	t.Run("linker-stamped binary installs its own tag", func(t *testing.T) {
		withPiReleaseIdentity(t, "0.27.2", "(devel)")
		piOps := &fakePiRuntimeOps{lookPath: piHealthyPathFixtures()}
		var stdout, stderr bytes.Buffer
		if code := runInitWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, piOps, piTestEnv(pkg, t.TempDir()), &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr.String())
		}
		want := "git:github.com/spacedock-dev/spacedock@v0.27.2"
		if len(piOps.piInstalls) != 1 || piOps.piInstalls[0] != want {
			t.Fatalf("stamped binary must install %q, got %v", want, piOps.piInstalls)
		}
	})

	t.Run("proxy-tagged build installs its tag", func(t *testing.T) {
		withPiReleaseIdentity(t, "dev", "v0.27.2")
		piOps := &fakePiRuntimeOps{lookPath: piHealthyPathFixtures()}
		var stdout, stderr bytes.Buffer
		if code := runInitWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, piOps, piTestEnv(pkg, t.TempDir()), &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr.String())
		}
		want := "git:github.com/spacedock-dev/spacedock@v0.27.2"
		if len(piOps.piInstalls) != 1 || piOps.piInstalls[0] != want {
			t.Fatalf("proxy build must pin its tag: want %q, got %v", want, piOps.piInstalls)
		}
	})

	t.Run("dev sentinel keeps the floating source", func(t *testing.T) {
		withPiReleaseIdentity(t, "dev", "(devel)")
		piOps := &fakePiRuntimeOps{lookPath: piHealthyPathFixtures()}
		var stdout, stderr bytes.Buffer
		if code := runInitWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, piOps, piTestEnv(pkg, t.TempDir()), &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr.String())
		}
		if len(piOps.piInstalls) != 1 || piOps.piInstalls[0] != piSpacedockPackageSource {
			t.Fatalf("dev build must keep the bare floating source, got %v", piOps.piInstalls)
		}
	})

	t.Run("plugin-dir override wins and never uses the release source", func(t *testing.T) {
		withPiReleaseIdentity(t, "0.27.2", "(devel)")
		piOps := &fakePiRuntimeOps{lookPath: piHealthyPathFixtures()}
		var stdout, stderr bytes.Buffer
		if code := runInitWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", "/checkout"}, &fakeHost{}, piOps, piTestEnv(pkg, t.TempDir()), &stdout, &stderr); code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr.String())
		}
		if len(piOps.piInstalls) != 1 || piOps.piInstalls[0] != "/checkout" {
			t.Fatalf("plugin dir must be the install source, got %v", piOps.piInstalls)
		}
	})
}

func TestPiFrontDoorRepairsUnpinnedPackageAndLaunches(t *testing.T) {
	repo, pkg := piRepairFixture(t)
	// The released v0.27.2 incident shape: registered but UNPINNED.
	repaired := healthyPiPackageStatus()
	ops := piRepairOps(repo, pkg, piPackageStatus{registered: true, ensignDiscoverable: false, source: "git:github.com/spacedock-dev/spacedock"}, &repaired)
	withPiReleaseIdentity(t, "0.27.2", "(devel)")
	code, _, stderr := runPiWith(t, ops, pkg, nil)

	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if len(ops.launched) == 0 {
		t.Fatalf("expected launch after repair; stderr=%q", stderr)
	}
	want := "git:github.com/spacedock-dev/spacedock@v0.27.2"
	if len(ops.piInstalls) != 1 || ops.piInstalls[0] != want {
		t.Fatalf("expected exactly one install of %q, got %v", want, ops.piInstalls)
	}
	// Initial check + one post-repair recheck; never a second install.
	if ops.statusCalls != 2 {
		t.Fatalf("expected 2 status reads (initial + recheck), got %d", ops.statusCalls)
	}
}

func TestPiFrontDoorRepairsWrongReleaseLine(t *testing.T) {
	repo, pkg := piRepairFixture(t)
	// The incident: @v0.28.0-pre1 installed under a v0.27.2-identity binary.
	wrongLine := healthyPiPackageStatus()
	wrongLine.source = "git:github.com/spacedock-dev/spacedock@v0.28.0-pre1"
	repaired := healthyPiPackageStatus()
	repaired.source = "git:github.com/spacedock-dev/spacedock@v0.27.2"
	ops := piRepairOps(repo, pkg, wrongLine, &repaired)
	withPiReleaseIdentity(t, "0.27.2", "(devel)")
	code, _, stderr := runPiWith(t, ops, pkg, nil)

	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if want := "git:github.com/spacedock-dev/spacedock@v0.27.2"; len(ops.piInstalls) != 1 || ops.piInstalls[0] != want {
		t.Fatalf("wrong-line package must repair to the binary's own ref %q, got %v", want, ops.piInstalls)
	}
	if len(ops.launched) == 0 {
		t.Fatalf("expected launch after wrong-line repair")
	}
}

func TestPiFrontDoorRefusesWhenRepairInstallFails(t *testing.T) {
	repo, pkg := piRepairFixture(t)
	ops := piRepairOps(repo, pkg, piPackageStatus{}, nil)
	ops.piInstallErr = errors.New("clone failed")
	withPiReleaseIdentity(t, "0.27.2", "(devel)")
	code, _, stderr := runPiWith(t, ops, pkg, nil)

	if code != 1 {
		t.Fatalf("exit=%d want 1", code)
	}
	// AC-4: one install attempt, one recheck, no launch, actionable error.
	if len(ops.piInstalls) != 1 || len(ops.launched) != 0 {
		t.Fatalf("want 1 install and no launch: installs=%v launched=%v", ops.piInstalls, ops.launched)
	}
	if ops.statusCalls != 2 {
		t.Fatalf("expected exactly one recheck after the repair, got %d status reads", ops.statusCalls)
	}
	if !strings.Contains(stderr, "package repair install failed") || !strings.Contains(stderr, "spacedock doctor --host pi") {
		t.Fatalf("actionable error missing; stderr=%q", stderr)
	}
}

func TestPiFrontDoorRefusesAfterIneffectiveRepair(t *testing.T) {
	repo, pkg := piRepairFixture(t)
	stillUnpinned := piPackageStatus{registered: true, ensignDiscoverable: false, source: "git:github.com/spacedock-dev/spacedock"}
	ops := piRepairOps(repo, pkg, piPackageStatus{registered: true, source: "git:github.com/spacedock-dev/spacedock"}, &stillUnpinned)
	withPiReleaseIdentity(t, "0.27.2", "(devel)")
	code, _, _ := runPiWith(t, ops, pkg, nil)

	if code != 1 {
		t.Fatalf("exit=%d want 1", code)
	}
	// AC-4: no second install attempt, no launch after a remaining mismatch.
	if len(ops.piInstalls) != 1 || len(ops.launched) != 0 {
		t.Fatalf("one install, no launch expected: installs=%v launched=%v", ops.piInstalls, ops.launched)
	}
}

func TestPiFrontDoorSuppressedRepairs(t *testing.T) {
	t.Run("non-git healthy entry never repairs", func(t *testing.T) {
		repo, pkg := piRepairFixture(t)
		userEntry := piPackageStatus{registered: true, ensignDiscoverable: true, source: "file:/x", packageRoot: "/x"}
		ops := piRepairOps(repo, pkg, userEntry, nil)
		withPiReleaseIdentity(t, "0.27.2", "(devel)")
		code, _, stderr := runPiWith(t, ops, pkg, nil)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr)
		}
		if len(ops.piInstalls) != 0 {
			t.Fatalf("user-managed entry must not be repaired: %v", ops.piInstalls)
		}
	})

	t.Run("dev-sentinel binary never repairs", func(t *testing.T) {
		repo, pkg := piRepairFixture(t)
		ops := piRepairOps(repo, pkg, piPackageStatus{}, nil)
		withPiReleaseIdentity(t, "dev", "(devel)")
		code, _, stderr := runPiWith(t, ops, pkg, nil)
		// A dev build has no pin target: the pre-existing refusal stands.
		if code != 1 {
			t.Fatalf("exit=%d want 1; stderr=%q", code, stderr)
		}
		if len(ops.piInstalls) != 0 {
			t.Fatalf("dev build must not repair-install: %v", ops.piInstalls)
		}
	})

	t.Run("plugin-dir override suppresses the repair", func(t *testing.T) {
		repo, pkg := piRepairFixture(t)
		ops := piRepairOps(repo, pkg, piPackageStatus{}, nil)
		withPiReleaseIdentity(t, "0.27.2", "(devel)")
		var stdout, stderr bytes.Buffer
		code := runPi(context.Background(), []string{"review this", "--plugin-dir", repo}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr.String())
		}
		if len(ops.piInstalls) != 0 {
			t.Fatalf("dev override suppresses the repair: %v", ops.piInstalls)
		}
	})
}

func TestPiFrontDoorHealthyPinnedPackageDoesNotRepair(t *testing.T) {
	repo, pkg := piRepairFixture(t)
	pinned := healthyPiPackageStatus()
	pinned.source = "git:github.com/spacedock-dev/spacedock@v0.27.2"
	ops := piRepairOps(repo, pkg, pinned, nil)
	withPiReleaseIdentity(t, "0.27.2", "(devel)")
	code, _, stderr := runPiWith(t, ops, pkg, nil)

	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if len(ops.piInstalls) != 0 {
		t.Fatalf("healthy pinned package must not repair: %v", ops.piInstalls)
	}
	if len(ops.launched) == 0 {
		t.Fatalf("expected launch")
	}
}
