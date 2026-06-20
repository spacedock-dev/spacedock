package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildFnmMultishellTree builds a synthetic fnm tree in a temp dir mirroring
// the live fnm layout: a per-shell `fnm_multishells/<pid>_<ts>` symlink to a
// stable `node-versions/<ver>/installation` dir, with `…/installation/bin/pi`
// as a relative symlink to `../lib/node_modules/pkg/dist/cli.js` (the target
// file is created). It returns the multishell looked-up pi path and the stable
// install bin path. The stable install bin's parent is a REAL directory (not a
// symlink), matching fnm — fnm only ever unlinks the per-shell `<pid>_<ts>`
// symlink, never the shared installation tree.
func buildFnmMultishellTree(t *testing.T) (multishellPi, stablePi string) {
	t.Helper()
	root := t.TempDir()
	installation := filepath.Join(root, "stable-versions", "v1.2.3", "installation")
	binDir := filepath.Join(installation, "bin")
	libTarget := filepath.Join(installation, "lib", "node_modules", "pkg", "dist", "cli.js")
	writeFileWithDirs(t, libTarget, "#!/usr/bin/env node\n")
	// bin/pi is a relative symlink to ../lib/node_modules/pkg/dist/cli.js, exactly
	// as fnm's installed pi bin is.
	piLink := filepath.Join(binDir, "pi")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../lib/node_modules/pkg/dist/cli.js", piLink); err != nil {
		t.Fatal(err)
	}
	// Per-shell multishell dir: fnm_multishells/<pid>_<ts> is a symlink to the
	// stable installation dir.
	multishellParent := filepath.Join(root, "fnm_multishells")
	if err := os.MkdirAll(multishellParent, 0o755); err != nil {
		t.Fatal(err)
	}
	multishellLink := filepath.Join(multishellParent, "123_456")
	if err := os.Symlink(installation, multishellLink); err != nil {
		t.Fatal(err)
	}
	multishellPi = filepath.Join(multishellLink, "bin", "pi")
	// Compute the expected stable path the same way the helper does:
	// EvalSymlinks of the bin dir (resolves any OS-level /var -> /private/var
	// symlink on macOS) then Join(..., "pi").
	stableBinDir, err := filepath.EvalSymlinks(binDir)
	if err != nil {
		t.Fatal(err)
	}
	stablePi = filepath.Join(stableBinDir, "pi")
	return multishellPi, stablePi
}

// TestResolveFnmMultishellPi_StableChosen pins AC-1: when lookedUp is a path
// under */fnm_multishells/*/bin/pi, the helper resolves through the per-shell
// symlink to the stable node-installation bin (NOT under fnm_multishells/).
func TestResolveFnmMultishellPi_StableChosen(t *testing.T) {
	multishellPi, stablePi := buildFnmMultishellTree(t)

	got, ok := resolveFnmMultishellPi(multishellPi)
	if !ok {
		t.Fatalf("resolveFnmMultishellPi(%q) ok=false, want true", multishellPi)
	}
	if got != stablePi {
		t.Fatalf("resolveFnmMultishellPi(%q) = %q, want stable %q", multishellPi, got, stablePi)
	}
	if strings.Contains(got, "/fnm_multishells/") {
		t.Fatalf("stable path still under fnm_multishells: %q", got)
	}
}

// TestResolveFnmMultishellPi_NonMultishell_Unchanged pins AC-2: a non-multishell
// lookedUp path returns ("", false) so runPi leaves argv[0]="pi" unchanged.
func TestResolveFnmMultishellPi_NonMultishell_Unchanged(t *testing.T) {
	for _, lookedUp := range []string{"/usr/local/bin/pi", "/opt/bin/pi"} {
		if got, ok := resolveFnmMultishellPi(lookedUp); ok || got != "" {
			t.Fatalf("resolveFnmMultishellPi(%q) = (%q, %v), want (\"\", false)", lookedUp, got, ok)
		}
	}
}

// TestResolveFnmMultishellPi_TornDownParent_Fallback pins the fall-back
// guarantee: when the multishell parent symlink has been torn down (a sibling
// shell exited), EvalSymlinks errors and the helper returns ("", false) — runPi
// leaves argv[0]="pi" rather than blocking the launch.
func TestResolveFnmMultishellPi_TornDownParent_Fallback(t *testing.T) {
	multishellPi, _ := buildFnmMultishellTree(t)
	// Simulate fnm tearing down the per-shell multishell dir.
	multishellLink := filepath.Dir(filepath.Dir(multishellPi))
	if err := os.Remove(multishellLink); err != nil {
		t.Fatal(err)
	}
	if got, ok := resolveFnmMultishellPi(multishellPi); ok || got != "" {
		t.Fatalf("resolveFnmMultishellPi after teardown = (%q, %v), want (\"\", false)", got, ok)
	}
}

// TestRunPi_LaunchArgv0_StableForMultishell pins AC-1 end-to-end at the runPi
// seam: when ops.LookPath("pi") returns a multishell path, the launched argv[0]
// is the stable install bin (NOT under fnm_multishells/).
func TestRunPi_LaunchArgv0_StableForMultishell(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	multishellPi, stablePi := buildFnmMultishellTree(t)

	ops := &fakePiRuntimeOps{
		lookPath: map[string]string{
			"pi":        multishellPi,
			"safehouse": "/bin/safehouse",
		},
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runPi exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(ops.launched) == 0 {
		t.Fatalf("Launch was not invoked")
	}
	if got := ops.launched[0]; got != stablePi {
		t.Fatalf("launched argv[0]=%q, want stable %q\nargv=%v", got, stablePi, ops.launched)
	}
	if strings.Contains(ops.launched[0], "/fnm_multishells/") {
		t.Fatalf("launched argv[0] is still a multishell path: %q", ops.launched[0])
	}
}

// TestRunPi_LaunchArgv0_BarePiUnderWrap pins the feedback-cycle-1 REVISED
// regression: when BOTH the safehouse wrap is triggered AND ops.LookPath("pi")
// returns a fnm multishell path, the fnm resolution must be SUPPRESSED under
// wrap. The launched argv[0] must be "safehouse" and the inner token after the
// safehouse "--" separator must be the bare "pi" (NOT the stable install bin
// path). Rationale: safehouse execs the inner command INSIDE the sandbox, whose
// filesystem view does not include ~/.local/share/fnm/node-versions/... (not
// the workdir, not added via --safehouse-add-dirs), so handing safehouse the
// stable ABSOLUTE path as the inner token makes Node's Module._resolveFilename
// fail with MODULE_NOT_FOUND (manual pi-live smoke, feedback cycle 1 revised).
// Under wrap we leave argv[0]="pi" (parity with claude/codex, which pass the
// bare name as inner[0] even when wrapped — frontdoor.go:315,498); the race is
// accepted as a known shared limitation under safehouse. This test REDs on the
// cycle-1 reorder code (which set inner=stable even under wrap) and GREENs after
// the !wrap gate lands.
func TestRunPi_LaunchArgv0_BarePiUnderWrap(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	multishellPi, stablePi := buildFnmMultishellTree(t)

	ops := &fakePiRuntimeOps{
		lookPath: map[string]string{
			"pi":        multishellPi,
			"safehouse": "/bin/safehouse",
		},
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	// A --safehouse-* knob triggers the wrap (no .safehouse profile needed);
	// safehouse.Available passes because the fake resolves "safehouse".
	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--safehouse-add-dirs=/a", "--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runPi exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(ops.launched) == 0 {
		t.Fatalf("Launch was not invoked")
	}
	// argv[0] must be "safehouse" — the fnm resolution must not clobber the wrap.
	if got := ops.launched[0]; got != "safehouse" {
		t.Fatalf("launched argv[0]=%q, want safehouse\nargv=%v", got, ops.launched)
	}
	// The inner pi token (right after the safehouse "--" separator) must be the
	// BARE "pi", NOT the stable install bin path and NOT the multishell path.
	// Handing safehouse the stable absolute path makes Node's
	// Module._resolveFilename fail with MODULE_NOT_FOUND inside the sandbox.
	inner := piSafehouseInnerArgv(ops.launched)
	if len(inner) == 0 {
		t.Fatalf("no inner argv after safehouse -- separator: %v", ops.launched)
	}
	if got := inner[0]; got != "pi" {
		t.Fatalf("inner argv[0]=%q, want bare %q (fnm resolution must be suppressed under wrap; stable path %q is not visible in the sandbox)\ninner=%v", got, "pi", stablePi, inner)
	}
}

// TestRunPi_LaunchArgv0_UnchangedForDirect pins AC-2: when ops.LookPath("pi")
// returns a non-multishell path, the launched argv[0] is the literal "pi"
// exactly as today.
func TestRunPi_LaunchArgv0_UnchangedForDirect(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)

	ops := &fakePiRuntimeOps{
		lookPath: map[string]string{
			"pi":        "/usr/local/bin/pi",
			"safehouse": "/bin/safehouse",
		},
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runPi exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(ops.launched) == 0 {
		t.Fatalf("Launch was not invoked")
	}
	if got := ops.launched[0]; got != "pi" {
		t.Fatalf("launched argv[0]=%q, want \"pi\" (unchanged)\nargv=%v", got, ops.launched)
	}
}
