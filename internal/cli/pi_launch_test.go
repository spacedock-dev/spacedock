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

// TestFnmStableSandboxDir_NormalInstall pins the normal `npm i -g` walk
// (feedback cycle 2): the stable `pi` is a symlink to a real script under
// installation/lib/node_modules/... with NO workspace root, so the walk grants
// the highest node_modules-bearing ancestor = installation/lib (covers the pi
// package + all its deps).
func TestFnmStableSandboxDir_NormalInstall(t *testing.T) {
	_, stablePi := buildFnmMultishellTree(t)
	installation := filepath.Dir(filepath.Dir(stablePi)) // .../installation/bin/pi -> installation
	want := filepath.Join(installation, "lib")

	got := fnmStableSandboxDir(stablePi)
	if got != want {
		t.Fatalf("fnmStableSandboxDir(%q) = %q, want %q (installation/lib, the node_modules root)", stablePi, got, want)
	}
}

// TestFnmStableSandboxDir_DevLinkWorkspaces pins the dev-link (npm link) walk
// (feedback cycle 2, the captain's actual env): the stable `pi` symlink chain
// escapes the fnm installation out to a monorepo workspace root whose
// node_modules holds the HOISTED runtime deps. The walk must stop at the
// workspace root (package.json with `workspaces`), NOT the package's own
// node_modules (which holds only dev/test deps). Mirrors the live env where
// cross-spawn is hoisted to pi-mono/node_modules, not coding-agent/node_modules.
func TestFnmStableSandboxDir_DevLinkWorkspaces(t *testing.T) {
	root := t.TempDir()
	// Workspace root: package.json with `workspaces`, node_modules with the
	// hoisted runtime dep (cross-spawn), and a packages/<pkg> workspace.
	ws := filepath.Join(root, "pi-mono")
	writeFileWithDirs(t, filepath.Join(ws, "package.json"), `{"name":"pi-mono","workspaces":["packages/*"]}`+"\n")
	if err := os.MkdirAll(filepath.Join(ws, "node_modules", "cross-spawn"), 0o755); err != nil {
		t.Fatal(err)
	}
	pkg := filepath.Join(ws, "packages", "coding-agent")
	// The package's OWN node_modules holds only dev/test deps (NOT cross-spawn),
	// exactly as the captain verified for the live env.
	if err := os.MkdirAll(filepath.Join(pkg, "node_modules", "chai"), 0o755); err != nil {
		t.Fatal(err)
	}
	realScript := filepath.Join(pkg, "dist", "cli.js")
	writeFileWithDirs(t, realScript, "#!/usr/bin/env node\n")

	// fnm installation: bin/pi -> lib/node_modules/@earendil-works/pi-coding-agent/dist/cli.js,
	// and @earendil-works/pi-coding-agent is a symlink out to the workspace pkg
	// (the npm link). EvalSymlinks(installation/bin/pi) resolves to ws pkg cli.js.
	installation := filepath.Join(root, "node-versions", "v24", "installation")
	binDir := filepath.Join(installation, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	nmPkg := filepath.Join(installation, "lib", "node_modules", "@earendil-works", "pi-coding-agent")
	if err := os.MkdirAll(filepath.Dir(nmPkg), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(pkg, nmPkg); err != nil { // npm link: package -> workspace pkg
		t.Fatal(err)
	}
	stable := filepath.Join(binDir, "pi")
	if err := os.Symlink(filepath.Join("..", "lib", "node_modules", "@earendil-works", "pi-coding-agent", "dist", "cli.js"), stable); err != nil {
		t.Fatal(err)
	}

	got := fnmStableSandboxDir(stable)
	// EvalSymlinks walks the real (resolved) path, which on macOS resolves the
	// temp dir's /var -> /private/var symlink, so compare against the resolved ws.
	wantWS, err := filepath.EvalSymlinks(ws)
	if err != nil {
		wantWS = ws
	}
	if got != wantWS {
		t.Fatalf("fnmStableSandboxDir dev-link = %q, want workspace root %q\n(must stop at the workspaces root whose node_modules holds the hoisted deps, not the package's own node_modules)", got, wantWS)
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

// TestRunPi_LaunchArgv0_StableUnderWrapWithAddDirs pins the feedback-cycle-2 fix:
// when BOTH the safehouse wrap is triggered AND ops.LookPath("pi") returns a fnm
// multishell path, the fnm resolution is NO LONGER suppressed under wrap. The
// launched argv[0] must be "safehouse" and the inner token after the safehouse
// "--" separator must be the STABLE install bin path (NOT bare "pi", NOT the
// multishell path), AND the safehouse extra/argv must contain
// --add-dirs=<stable-bin-dir> so the sandbox can see the stable path. Rationale
// (feedback cycle 2): cycle-1-revised gated the resolution to !wrap and left bare
// "pi" under wrap, but the sandbox PRESERVES the inherited PATH, which carried a
// stale (tearing-down) fnm_multishells/<pid>_<ts> dir → MODULE_NOT_FOUND on the
// multishell path. The coupled fix resolves the stable path ALWAYS and grants
// the stable bin's parent dir to the sandbox via --safehouse-add-dirs (a REAL dir
// fnm never tears down) so Node's Module._resolveFilename sees it. Adversarial:
// (a) drop the add-dirs grant under wrap → this test REDs on the missing
// --add-dirs token (the sandbox would have no visibility). (b) re-gate the
// resolution to !wrap → this test REDs with inner argv[0]=="pi" (stable path
// suppressed under wrap).
func TestRunPi_LaunchArgv0_StableUnderWrapWithAddDirs(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	multishellPi, stablePi := buildFnmMultishellTree(t)
	// The synthetic tree is the normal-install shape (no workspaces), so
	// fnmStableSandboxDir walks up from the resolved real script and lands on the
	// highest node_modules-bearing ancestor = installation/lib (covers the pi
	// package + its deps). Dev-link workspaces path is covered by
	// TestFnmStableSandboxDir_DevLinkWorkspaces below.
	installation := filepath.Dir(filepath.Dir(stablePi)) // .../installation/bin/pi -> installation
	wantGrant := filepath.Join(installation, "lib")

	ops := &fakePiRuntimeOps{
		lookPath: map[string]string{
			"pi":        multishellPi,
			"safehouse": "/bin/safehouse",
		},
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	// --safehouse-enable triggers the wrap (no .safehouse profile, no operator
	// add-dirs that would confound the fnm add-dirs assertion); safehouse.Available
	// passes because the fake resolves "safehouse".
	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--safehouse-enable=ssh", "--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runPi exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(ops.launched) == 0 {
		t.Fatalf("Launch was not invoked")
	}
	// argv[0] must be "safehouse" — the wrap fired.
	if got := ops.launched[0]; got != "safehouse" {
		t.Fatalf("launched argv[0]=%q, want safehouse\nargv=%v", got, ops.launched)
	}
	// The inner pi token (right after the safehouse "--" separator) must be the
	// STABLE install bin path — NOT bare "pi" and NOT the multishell path. The
	// resolution now applies under wrap too (feedback cycle 2).
	inner := piSafehouseInnerArgv(ops.launched)
	if len(inner) == 0 {
		t.Fatalf("no inner argv after safehouse -- separator: %v", ops.launched)
	}
	if got := inner[0]; got != stablePi {
		t.Fatalf("inner argv[0]=%q, want stable %q (resolution must apply under wrap)\ninner=%v", got, stablePi, inner)
	}
	if strings.Contains(inner[0], "/fnm_multishells/") {
		t.Fatalf("inner argv[0] is still a multishell path: %q", inner[0])
	}
	// The safehouse extra (before the "--" separator) must grant the dir whose
	// node_modules holds the stable `pi`'s real script + hoisted deps, computed
	// by fnmStableSandboxDir. For the synthetic normal-install tree that is
	// installation/lib. The token is --add-dirs=<grant>.
	wantAddDirs := "--add-dirs=" + wantGrant
	sep := -1
	for i, tok := range ops.launched {
		if tok == "--" {
			sep = i
			break
		}
	}
	if sep < 0 {
		t.Fatalf("no -- separator in wrapped argv: %v", ops.launched)
	}
	extra := ops.launched[:sep]
	found := false
	for _, tok := range extra {
		if tok == wantAddDirs {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("safehouse extra missing %q (sandbox would not see the stable bin)\nextra=%v", wantAddDirs, extra)
	}
}

// TestRunPi_LaunchEnvForwarding pins AC-4: the Launch(argv, env) signature
// change lands and env is threaded into the launch (mirroring claude/codex's
// launchEnv plumbing, frontdoor.go:348). runPi must call ops.Launch(argv,
// launchEnv(os.Environ())) — the env reaches the fake, and launchEnv sets
// SPACEDOCK_BIN from the resolved launcher binary so the sandbox --env-pass
// forwarding carries it through. This test fails to compile (and REDs) if the
// piRuntimeOps.Launch signature regresses to (argv []string) — the fake's
// method no longer satisfies the interface.
func TestRunPi_LaunchEnvForwarding(t *testing.T) {
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
	if len(ops.launchedEnv) == 0 {
		t.Fatalf("Launch was invoked with no env; env must be threaded (Launch(argv, env))")
	}
	// launchEnv sets SPACEDOCK_BIN from the resolved launcher binary (the test
	// binary under `go test`); it must reach the launched env so the sandbox
	// --env-pass forwarding carries it through, exactly as claude/codex do.
	var spacedockBin string
	found := false
	for _, kv := range ops.launchedEnv {
		if k, v, ok := strings.Cut(kv, "="); ok && k == spacedockBinEnv {
			spacedockBin = v
			found = true
		}
	}
	if !found {
		t.Fatalf("%s not in launched env (launchEnv must set it): %v", spacedockBinEnv, ops.launchedEnv)
	}
	if spacedockBin == "" {
		t.Fatalf("%s in launched env is empty: %v", spacedockBinEnv, ops.launchedEnv)
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
