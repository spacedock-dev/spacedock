package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/safehouse"
)

type fakePiRuntimeOps struct {
	lookPath      map[string]string
	statOK        map[string]bool
	launched      []string
	launchedEnv   []string
	launchCode    int      // host exit code Launch returns (default 0)
	piInstalls    []string // sources captured by PiInstall
	piInstallOut  string
	piInstallErr  error
	packageStatus piPackageStatus
	// piVersionOut/Err fake `pi --version`. Empty out with nil err means a
	// healthy current binary (0.85.1, above the 0.83.0 floor) so the many
	// existing fixtures don't each need a version field; sub-floor and
	// garbage-version tests set it explicitly.
	piVersionOut string
	piVersionErr error
}

func (f *fakePiRuntimeOps) LookPath(name string) (string, error) {
	if p, ok := f.lookPath[name]; ok {
		return p, nil
	}
	return "", errors.New("not found")
}

func (f *fakePiRuntimeOps) Stat(path string) error {
	if f.statOK[path] {
		return nil
	}
	return errors.New("missing")
}

func (f *fakePiRuntimeOps) Launch(argv []string, env []string) (int, error) {
	f.launched = append([]string(nil), argv...)
	f.launchedEnv = append([]string(nil), env...)
	return f.launchCode, nil
}

func (f *fakePiRuntimeOps) PiInstall(source string) (string, error) {
	f.piInstalls = append(f.piInstalls, source)
	return f.piInstallOut, f.piInstallErr
}

func (f *fakePiRuntimeOps) SpacedockPackageStatus(agentDir, home string) piPackageStatus {
	return f.packageStatus
}

func (f *fakePiRuntimeOps) PiVersion() (string, error) {
	if f.piVersionOut == "" && f.piVersionErr == nil {
		return "0.85.1\n", nil
	}
	return f.piVersionOut, f.piVersionErr
}

// healthyPiPackageStatus is the canned status for a registered, discoverable
// Spacedock package — the state `spacedock install --host pi` produces.
func healthyPiPackageStatus() piPackageStatus {
	return piPackageStatus{
		registered:               true,
		ensignDiscoverable:       true,
		firstOfficerDiscoverable: true,
		source:                   "git:github.com/spacedock-dev/spacedock",
		packageRoot:              "/pkg-store/spacedock",
	}
}

func TestPiCommandRegisteredInTopLevelHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"pi      [task] [-- pi-flags]",
		"Start Pi as your Spacedock first officer",
		"install  [--host claude|codex|pi]",
		"doctor   [--host claude|codex|pi]",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("top-level help missing %q:\n%s", want, out)
		}
	}
}

func TestPiFrontDoorLaunchesWithNativeResourcePaths(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"review this", "--plugin-dir", repo, "--", "--model", "google/gemini"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	// The retired --skill <repo>/skills/{first-officer,ensign} flags are absent;
	// only the pi-subagents skill is passed. The Spacedock skills are discovered
	// from the installed package's extension (resources_discover), not the flags.
	wantPrefix := []string{
		"pi",
		"--extension", filepath.Join(pkg, "src", "extension", "index.ts"),
		"--skill", filepath.Join(pkg, "skills", "pi-subagents"),
		"--model", "google/gemini",
	}
	if len(ops.launched) < len(wantPrefix)+1 {
		t.Fatalf("launch argv too short: %v", ops.launched)
	}
	for i, want := range wantPrefix {
		if ops.launched[i] != want {
			t.Fatalf("launch argv[%d]=%q want %q\nargv=%v", i, ops.launched[i], want, ops.launched)
		}
	}
	joined := strings.Join(ops.launched, " ")
	for _, banned := range []string{"Agent", "SendMessage", "TeamCreate", "TeamDelete", "--agent", "codex"} {
		if strings.Contains(joined, banned) {
			t.Fatalf("pi launch argv contains banned runtime token %q: %v", banned, ops.launched)
		}
	}
	// Exactly one --skill flag (pi-subagents); the retired first-officer/ensign
	// skill flags must not appear.
	if got := strings.Count(joined, "--skill"); got != 1 {
		t.Fatalf("expected exactly 1 --skill flag (pi-subagents only), got %d: %v", got, ops.launched)
	}
	if strings.Contains(joined, filepath.Join(repo, "skills", "first-officer")) || strings.Contains(joined, filepath.Join(repo, "skills", "ensign")) {
		t.Fatalf("pi launch argv must not pass retired repo skill flags: %v", ops.launched)
	}
	prompt := ops.launched[len(ops.launched)-1]
	// The frontdoor launch prompt is argv-only now: the operator task lands as
	// the launch prompt; the contract sentence is gone (the extension's gated
	// injection owns the contract delivery).
	if prompt != "review this" {
		t.Fatalf("pi launch prompt = %q, want the bare operator task (no contract sentence)", prompt)
	}
	if strings.Contains(strings.Join(ops.launched, " "), "$spacedock:") {
		t.Fatalf("pi launch argv must not carry $spacedock: syntax (unexpandable in a launch prompt): %v", ops.launched)
	}
	// The launch marker is set in the child env on every launch shape (AC-3/AC-4).
	if !hasLaunchMarkerEnv(ops.launchedEnv) {
		t.Fatalf("launched env missing %s=1: %v", piLaunchMarkerEnv, ops.launchedEnv)
	}
}

// TestRunPi_DevOverridePassesSpacedockExtensionAndSkills pins the dev-override
// skill-loading fix (AC-2): when --plugin-dir / SPACEDOCK_REPO_ROOT is set
// (cfg.repoRoot != "") AND the Spacedock extension exists at
// <repo>/.pi/extensions/spacedock.ts, runPi appends --extension <ext> and
// --skill <repo>/skills so pi loads the extension (resources_discover) and the
// parent session announces the Spacedock FO/ensign skills. pi does NOT
// auto-discover .pi/extensions/ from cwd, so without these flags the dev path
// boots skill-less (the eq regression).
func TestRunPi_DevOverridePassesSpacedockExtensionAndSkills(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	// The dev-override extension must exist for the os.Stat guard to add the flags.
	writeFileWithDirs(t, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer
	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	joined := strings.Join(ops.launched, " ")
	wantExt := filepath.Join(repo, ".pi", "extensions", "spacedock.ts")
	wantSkills := filepath.Join(repo, "skills")
	if !strings.Contains(joined, "--extension "+wantExt) {
		t.Fatalf("dev-override argv missing --extension %q: %v", wantExt, ops.launched)
	}
	if !strings.Contains(joined, "--skill "+wantSkills) {
		t.Fatalf("dev-override argv missing --skill %q: %v", wantSkills, ops.launched)
	}
	// The pi-subagents extension + skill are still passed first.
	if got := strings.Count(joined, "--skill"); got != 2 {
		t.Fatalf("expected 2 --skill flags (pi-subagents + spacedock skills), got %d: %v", got, ops.launched)
	}
	// Ordering: pi-subagents extension/skill precede the Spacedock extension/skill.
	piSubExt := strings.Index(joined, "--extension "+filepath.Join(pkg, "src", "extension", "index.ts"))
	sdExt := strings.Index(joined, "--extension "+wantExt)
	if piSubExt < 0 || sdExt < 0 || sdExt < piSubExt {
		t.Fatalf("expected pi-subagents extension before spacedock extension: %v", ops.launched)
	}
}

// TestRunPi_DevOverrideWithoutExtensionFallsBackGracefully pins the ready-gate
// tightening (AC-5a): when the dev-override is active but
// <repo>/.pi/extensions/spacedock.ts is absent, runPi REFUSES (the not-ready
// path) instead of silently launching a contract-less session — the old
// graceful flag-drop launched a skill-less pi, the 2026-09-10 incident's
// failure shape.
func TestRunPi_DevOverrideWithoutExtensionFallsBackGracefully(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	// NOTE: .pi/extensions/spacedock.ts intentionally NOT created.
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	statOK := statOKForPiResources(repo, pkg)
	delete(statOK, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"))
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOK,
		packageStatus: piPackageStatus{}, // not registered: the dev-override fallback arm
	}
	var stdout, stderr bytes.Buffer
	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit=%d want 1 (not-ready refusal when the dev-override checkout has no Spacedock extension): stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "MISSING Spacedock extension") {
		t.Fatalf("doctor report missing MISSING Spacedock extension line:\n%s", stdout.String())
	}
	if len(ops.launched) != 0 {
		t.Fatalf("refused runtime must not launch: %v", ops.launched)
	}
}

// TestRunPi_InstalledPathDoesNotPassSpacedockExtension pins AC-3: when there is
// NO dev override (cfg.repoRoot == ""), runPi does NOT add the Spacedock
// extension/skill flags. The installed path relies on pi auto-loading registered
// extensions from settings.json `packages` at startup (the install-managed
// contract shipped by eq) — verified empirically; runPi does not need to pass
// the extension explicitly in the installed case.
func TestRunPi_InstalledPathDoesNotPassSpacedockExtension(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	// The ready gate Stats the effective package root's extension (AC-5a);
	// this test points packageRoot at the repo.
	writeFileWithDirs(t, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	// No --plugin-dir / SPACEDOCK_REPO_ROOT: installed path (cfg.repoRoot == "").
	env := piTestEnv(pkg, t.TempDir())
	// Provide a healthy registered package status so the launch-ready gate passes.
	status := healthyPiPackageStatus()
	// Point the package root at the repo so the ensign skill Stat check resolves.
	status.packageRoot = repo
	statOK := statOKForPiResources(repo, pkg)
	statOK[filepath.Join(repo, ".pi", "extensions", "spacedock.ts")] = true
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOK,
		packageStatus: status,
	}
	var stdout, stderr bytes.Buffer
	code := runPi(context.Background(), []string{"--", "--version"}, t.TempDir(), env, ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	joined := strings.Join(ops.launched, " ")
	if strings.Contains(joined, filepath.Join(".pi", "extensions", "spacedock.ts")) {
		t.Fatalf("installed-path argv must not pass the Spacedock extension: %v", ops.launched)
	}
	if strings.Contains(joined, "--skill "+filepath.Join(repo, "skills")) {
		t.Fatalf("installed-path argv must not pass the repo skills: %v", ops.launched)
	}
	if got := strings.Count(joined, "--skill"); got != 1 {
		t.Fatalf("installed-path argv must have exactly 1 --skill (pi-subagents), got %d: %v", got, ops.launched)
	}
}

func TestPiInstallAcceptsPluginDirAsDevOverride(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakeHost{}
	piOps := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", "/checkout"}, ops, piOps, piTestEnv(pkg, t.TempDir()), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	// --plugin-dir is the dev-override install source: pi install <path>.
	if len(piOps.piInstalls) != 1 || piOps.piInstalls[0] != "/checkout" {
		t.Fatalf("expected pi install /checkout, got %v", piOps.piInstalls)
	}
	if len(ops.installCmds) != 0 {
		t.Fatalf("install --host pi must not call the host plugin install seam: %v", ops.installCmds)
	}
}

func TestPiInstallRunsPiInstallAndDoesNotUsePluginCommands(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakeHost{}
	piOps := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "pi"}, ops, piOps, append(piTestEnv(pkg, t.TempDir()), "SPACEDOCK_REPO_ROOT="+repo), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	// install --host pi runs `pi install <published source>`, not the host plugin seam.
	if len(piOps.piInstalls) != 1 || piOps.piInstalls[0] != piSpacedockPackageSource {
		t.Fatalf("expected pi install %q, got %v", piSpacedockPackageSource, piOps.piInstalls)
	}
	if len(ops.installCmds) != 0 {
		t.Fatalf("install --host pi called host plugin install seam: %v", ops.installCmds)
	}
	out := stdout.String()
	for _, want := range []string{"Pi runtime ready", "pi-subagents", "pi-intercom", pkg, "Spacedock package", "necessary supervisor-talkback setup prerequisites only"} {
		if !strings.Contains(out, want) {
			t.Fatalf("install --host pi output missing %q:\n%s", want, out)
		}
	}
}

func TestPiInstallMissingSubagentsPrintsActionableInstructions(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, &fakePiRuntimeOps{
		lookPath: map[string]string{"pi": "/bin/pi"},
		statOK: map[string]bool{
			filepath.Join(repo, "skills", "first-officer", "SKILL.md"): true,
			filepath.Join(repo, "skills", "ensign", "SKILL.md"):        true,
		},
	}, append(piTestEnv(pkg, t.TempDir()), "SPACEDOCK_REPO_ROOT="+repo), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("install --host pi should be idempotent/instructive, exit=%d stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"Pi runtime setup incomplete", "pi install npm:pi-subagents", "PI_SUBAGENTS_PACKAGE_ROOT", "pi-intercom", "PI_INTERCOM_PACKAGE_ROOT"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing-subagents output missing %q:\n%s", want, out)
		}
	}
}

func TestNonPiSetupRejectsPluginDir(t *testing.T) {
	for _, tc := range []struct {
		name       string
		run        func(hostOps, io.Writer, io.Writer) int
		wantStderr string
	}{
		{
			name: "install claude",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runInitWithPi(context.Background(), []string{"--host", "claude", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
			},
			wantStderr: "--plugin-dir is not supported",
		},
		{
			name: "doctor claude",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runDoctorWithPi(context.Background(), []string{"--host", "claude", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
			},
			wantStderr: "unknown argument \"--plugin-dir\"",
		},
		{
			name: "doctor codex",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runDoctorWithPi(context.Background(), []string{"--host", "codex", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
			},
			wantStderr: "unknown argument \"--plugin-dir\"",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ops := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := tc.run(ops, &stdout, &stderr)
			if code != 2 {
				t.Fatalf("exit=%d want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stderr.String(), tc.wantStderr) {
				t.Fatalf("stderr should contain %q, got %q", tc.wantStderr, stderr.String())
			}
			if len(ops.installCmds) != 0 {
				t.Fatalf("install seam called despite rejected --plugin-dir: %v", ops.installCmds)
			}
		})
	}
}

func TestPiInstallCheckFailsForMissingSupervisorTalkbackPrerequisites(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	home := t.TempDir()
	statOK := statOKForPiResources(repo, pkg)
	statOK[filepath.Join(home, ".pi", "agent", "auth.json")] = true
	delete(statOK, pkg+"-intercom")
	delete(statOK, filepath.Join(pkg+"-intercom", "skills", "pi-intercom", "SKILL.md"))
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "pi", "--check"}, &fakeHost{}, &fakePiRuntimeOps{
		lookPath: map[string]string{"pi": "/bin/pi"},
		statOK:   statOK,
	}, append(piTestEnv(pkg, home), "SPACEDOCK_REPO_ROOT="+repo), &stdout, &stderr)
	if code == 0 {
		t.Fatalf("install --host pi --check exit=0 want non-zero for missing supervisor-talkback prerequisites; stdout=%q", stdout.String())
	}
	out := stdout.String()
	for _, want := range []string{"OK pi-subagents intercom bridge", "MISSING pi-intercom package root", "MISSING pi-intercom skill", "PI_INTERCOM_PACKAGE_ROOT"} {
		if !strings.Contains(out, want) {
			t.Fatalf("install check output missing %q:\n%s", want, out)
		}
	}
}

func TestPiInstallCheckFailsForMissingAuthLikeDoctor(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	home := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "pi", "--check"}, &fakeHost{}, &fakePiRuntimeOps{
		lookPath: piHealthyPathFixtures(),
		statOK:   statOKForPiResources(repo, pkg),
	}, append(piTestEnv(pkg, home), "SPACEDOCK_REPO_ROOT="+repo), &stdout, &stderr)
	if code == 0 {
		t.Fatalf("install --host pi --check exit=0 want non-zero for missing Pi auth; stdout=%q", stdout.String())
	}
	out := stdout.String()
	for _, want := range []string{"Pi runtime check", "MISSING Pi auth", filepath.Join(home, ".pi", "agent", "auth.json"), "OK pi-intercom package root", "OK pi-intercom skill"} {
		if !strings.Contains(out, want) {
			t.Fatalf("install check missing-auth output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Pi runtime ready") {
		t.Fatalf("install check with missing auth should not print ready:\n%s", out)
	}
}

func TestPiRuntimeConfigResolvesEnvPathsForSubagentsIntercomAuthAndSessions(t *testing.T) {
	repo := filepath.Join(t.TempDir(), "repo")
	subagents := filepath.Join(t.TempDir(), "pi-subagents")
	intercom := filepath.Join(t.TempDir(), "pi-intercom")
	authRoot := filepath.Join(t.TempDir(), "coding-agent")
	sessionDir := filepath.Join(t.TempDir(), "sessions")

	cfg := piRuntimeConfigFromEnv([]string{
		"SPACEDOCK_REPO_ROOT=" + repo,
		"PI_SUBAGENTS_PACKAGE_ROOT=" + subagents,
		"PI_INTERCOM_PACKAGE_ROOT=" + intercom,
		"PI_CODING_AGENT_DIR=" + authRoot,
		"PI_CODING_AGENT_SESSION_DIR=" + sessionDir,
	}, t.TempDir(), "")

	assertEqual(t, cfg.repoRoot, repo)
	assertEqual(t, cfg.packageRoot, subagents)
	assertEqual(t, cfg.intercomPackageRoot, intercom)
	assertEqual(t, cfg.extensionPath, filepath.Join(subagents, "src", "extension", "index.ts"))
	assertEqual(t, cfg.subagentsSkill, filepath.Join(subagents, "skills", "pi-subagents"))
	assertEqual(t, cfg.authPath, filepath.Join(authRoot, "auth.json"))
	assertEqual(t, cfg.sessionDir, sessionDir)
	assertEqual(t, cfg.packageRootSource, "PI_SUBAGENTS_PACKAGE_ROOT")
	assertEqual(t, cfg.intercomPackageSource, "PI_INTERCOM_PACKAGE_ROOT")
	assertEqual(t, cfg.authPathSource, "PI_CODING_AGENT_DIR")
	assertEqual(t, cfg.sessionDirSource, "PI_CODING_AGENT_SESSION_DIR")
}

func TestPiRuntimeConfigDefaultsIntercomAndAuthPathsUnderHome(t *testing.T) {
	home := t.TempDir()
	cfg := piRuntimeConfigFromEnv([]string{"HOME=" + home}, "/checkout", "")

	assertEqual(t, cfg.packageRoot, filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-subagents"))
	assertEqual(t, cfg.intercomPackageRoot, filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-intercom"))
	assertEqual(t, cfg.authPath, filepath.Join(home, ".pi", "agent", "auth.json"))
	assertEqual(t, cfg.sessionDir, filepath.Join(home, ".pi", "agent", "sessions"))
	assertEqual(t, cfg.agentDir, filepath.Join(home, ".pi", "agent"))
}

// TestPiRuntimeConfigRetiresSkillFlagsAndCwdFallback is the AC-3 behavior test: the
// launcher's --skill first-officer/ensign flags and the cwd fallback are retired,
// the doctor gates on spacedockPackageOK (package registered + ensign discoverable
// as user-package), and the doctor reports OK from a non-repo cwd when installed.
func TestPiRuntimeConfigRetiresSkillFlagsAndCwdFallback(t *testing.T) {
	t.Run("cwd fallback removed", func(t *testing.T) {
		// No --plugin-dir, no SPACEDOCK_REPO_ROOT: repoRoot is empty (NOT the cwd),
		// even from a non-repo cwd. The cwd fallback is removed, not demoted.
		cfg := piRuntimeConfigFromEnv([]string{"HOME=/h"}, "/some/non-repo/cwd", "")
		assertEqual(t, cfg.repoRoot, "")
		assertEqual(t, cfg.pluginDirSource, "SPACEDOCK_REPO_ROOT")
	})

	t.Run("dev override retained", func(t *testing.T) {
		cfg := piRuntimeConfigFromEnv([]string{"HOME=/h"}, "/cwd", "/checkout")
		assertEqual(t, cfg.repoRoot, "/checkout")
		assertEqual(t, cfg.pluginDirSource, "--plugin-dir")
		cfg2 := piRuntimeConfigFromEnv([]string{"HOME=/h", "SPACEDOCK_REPO_ROOT=/env-repo"}, "/cwd", "")
		assertEqual(t, cfg2.repoRoot, "/env-repo")
	})

	t.Run("launch args have no retired skill flags", func(t *testing.T) {
		repo := t.TempDir()
		writePiSkillFixtures(t, repo)
		pkg := t.TempDir()
		writePiSubagentsFixtures(t, pkg)
		ops := &fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOKForPiResources(repo, pkg),
			packageStatus: healthyPiPackageStatus(),
		}
		var stdout, stderr bytes.Buffer
		code := runPi(context.Background(), []string{"--plugin-dir", repo}, "/tmp", piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
		}
		joined := strings.Join(ops.launched, " ")
		if got := strings.Count(joined, "--skill"); got != 1 {
			t.Fatalf("expected exactly 1 --skill flag (pi-subagents only), got %d: %v", got, ops.launched)
		}
		for _, banned := range []string{filepath.Join(repo, "skills", "first-officer"), filepath.Join(repo, "skills", "ensign")} {
			if strings.Contains(joined, banned) {
				t.Fatalf("retired repo skill flag present in launch args: %v", ops.launched)
			}
		}
	})

	t.Run("doctor gates on spacedockPackageOK from non-repo cwd", func(t *testing.T) {
		pkg := t.TempDir()
		writePiSubagentsFixtures(t, pkg)
		home := t.TempDir()
		auth := filepath.Join(home, ".pi", "agent", "auth.json")
		statOK := statOKForPiResources(t.TempDir(), pkg)
		statOK[auth] = true

		// Not installed: packageStatus zero -> spacedockPackageOK false -> not ready,
		// reported from a non-repo cwd (/tmp). This is the inverse of the old
		// cwd-fallback false-positive.
		var stdout, stderr bytes.Buffer
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath: piHealthyPathFixtures(),
			statOK:   statOK,
		}, piTestEnv(pkg, home), &stdout, &stderr)
		if code == 0 {
			t.Fatalf("doctor should be non-zero when package not installed; stdout=%q", stdout.String())
		}
		if !strings.Contains(stdout.String(), "MISSING Spacedock package") {
			t.Fatalf("doctor should report MISSING Spacedock package when not installed:\n%s", stdout.String())
		}

		// Installed: packageStatus healthy -> spacedockPackageOK true -> ready,
		// still from a non-repo cwd (/tmp). No --plugin-dir, no SPACEDOCK_REPO_ROOT.
		var stdout2, stderr2 bytes.Buffer
		code2 := runDoctorWithPi(context.Background(), []string{"--host", "pi"}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOK,
			packageStatus: healthyPiPackageStatus(),
		}, piTestEnv(pkg, home), &stdout2, &stderr2)
		if code2 != 0 {
			t.Fatalf("doctor should be zero when package installed from non-repo cwd; exit=%d stdout=%q", code2, stdout2.String())
		}
		if !strings.Contains(stdout2.String(), "OK Spacedock package") {
			t.Fatalf("doctor should report OK Spacedock package when installed:\n%s", stdout2.String())
		}
	})
}

func TestPiDoctorReportsMissingAndHealthyRuntime(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	home := t.TempDir()
	auth := filepath.Join(home, ".pi", "agent", "auth.json")

	t.Run("missing", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{}, piTestEnv(pkg, home), &stdout, &stderr)
		if code == 0 {
			t.Fatalf("exit=0 want non-zero for missing pi runtime")
		}
		out := stdout.String()
		for _, want := range []string{"Pi runtime check", "MISSING pi CLI", "MISSING Pi auth", "MISSING pi-subagents", "Supervisor-talkback setup prerequisites", "MISSING pi-subagents intercom bridge", "MISSING pi-intercom package root", "MISSING pi-intercom skill", "necessary supervisor-talkback setup prerequisites only"} {
			if !strings.Contains(out, want) {
				t.Fatalf("missing doctor output missing %q:\n%s", want, out)
			}
		}
		for _, notWant := range []string{"pi-intercom command", "subagents-doctor bridge-health command"} {
			if strings.Contains(out, notWant) {
				t.Fatalf("doctor output should not require unstable command contract %q:\n%s", notWant, out)
			}
		}
	})

	t.Run("openai-api-key-auth", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOKForPiResources(repo, pkg),
			packageStatus: healthyPiPackageStatus(),
		}, append(piTestEnv(pkg, home), "OPENAI_API_KEY=test-key"), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
		}
		if !strings.Contains(stdout.String(), "OK Pi auth") {
			t.Fatalf("OpenAI-key doctor output should accept env auth:\n%s", stdout.String())
		}
	})

	t.Run("healthy", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		statOK := statOKForPiResources(repo, pkg)
		statOK[auth] = true
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOK,
			packageStatus: healthyPiPackageStatus(),
		}, piTestEnv(pkg, home), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
		}
		out := stdout.String()
		for _, want := range []string{"OK pi CLI", "OK Pi auth", "OK pi-subagents extension", "OK pi-subagents intercom bridge", "OK pi-intercom package root", "OK pi-intercom skill", "OK Spacedock package", "live child talkback", "durable marker probe"} {
			if !strings.Contains(out, want) {
				t.Fatalf("healthy doctor output missing %q:\n%s", want, out)
			}
		}
		// The retired repo-path skill checks must not appear: the retired render
		// was `OK Spacedock <skill> skill: <cwd-derived path>`. The package-
		// discovery first-officer line (AC-7) uses a distinct label and is
		// asserted in TestPiDoctorReportsFirstOfficerVersionDuplicates.
		for _, notWant := range []string{"OK Spacedock first-officer skill: ", "OK Spacedock ensign skill"} {
			if strings.Contains(out, notWant) {
				t.Fatalf("healthy doctor output should not print retired skill check %q:\n%s", notWant, out)
			}
		}
	})
}

// TestPiRuntimeDevOverrideSatisfiesPackageGate verifies the regression fix for
// the --plugin-dir / SPACEDOCK_REPO_ROOT dev-override launch path: when the
// Spacedock package is NOT registered in settings.json (fresh pi-home), a
// dev-override repoRoot that contains skills/ensign/SKILL.md satisfies the
// package gate so the launch path reaches the ensign. The inverse (no
// repoRoot, no package) still fails the gate — the install-managed contract.
func TestPiRuntimeDevOverrideSatisfiesPackageGate(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	// The ready gate now also requires the effective package root's extension
	// (AC-5a); give the dev-override checkout one.
	writeFileWithDirs(t, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)

	t.Run("dev override satisfies gate without installed package", func(t *testing.T) {
		home := t.TempDir()
		cfg := piRuntimeConfigFromEnv(append(piTestEnv(pkg, home), "SPACEDOCK_REPO_ROOT="+repo), "/non-repo-cwd", "")
		if cfg.repoRoot != repo {
			t.Fatalf("cfg.repoRoot=%q want %q", cfg.repoRoot, repo)
		}
		check := checkPiRuntime(&fakePiRuntimeOps{
			lookPath: piHealthyPathFixtures(),
			statOK:   statOKForPiResources(repo, pkg),
			// No package registered in settings.json.
			packageStatus: piPackageStatus{},
		}, cfg)
		if !check.spacedockPackageOK {
			t.Fatalf("dev override should satisfy spacedockPackageOK; packageStatus=%+v", check.packageStatus)
		}
		if !piRuntimeLaunchReady(check) {
			t.Fatalf("dev override should make runtime launch-ready; check=%+v", check)
		}
		if check.packageStatus.source != repo+" (dev override)" {
			t.Fatalf("packageStatus.source=%q want %q", check.packageStatus.source, repo+" (dev override)")
		}
		if check.packageStatus.packageRoot != repo {
			t.Fatalf("packageStatus.packageRoot=%q want %q", check.packageStatus.packageRoot, repo)
		}
	})

	t.Run("runPi launches with dev override and no installed package", func(t *testing.T) {
		home := t.TempDir()
		var stdout, stderr bytes.Buffer
		code := runPi(context.Background(), []string{"do work", "--plugin-dir", repo, "--", "--print"}, "/non-repo-cwd",
			piTestEnv(pkg, home), &fakePiRuntimeOps{
				lookPath:      piHealthyPathFixtures(),
				statOK:        statOKForPiResources(repo, pkg),
				packageStatus: piPackageStatus{}, // not installed
			}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("runPi exit=%d want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
		}
		if strings.Contains(stdout.String(), "MISSING Spacedock package") {
			t.Fatalf("dev override must not report MISSING Spacedock package:\n%s", stdout.String())
		}
	})

	t.Run("no repoRoot and no package fails gate", func(t *testing.T) {
		home := t.TempDir()
		cfg := piRuntimeConfigFromEnv(piTestEnv(pkg, home), "/non-repo-cwd", "")
		if cfg.repoRoot != "" {
			t.Fatalf("cfg.repoRoot=%q want empty", cfg.repoRoot)
		}
		check := checkPiRuntime(&fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOKForPiResources(t.TempDir(), pkg),
			packageStatus: piPackageStatus{},
		}, cfg)
		if check.spacedockPackageOK {
			t.Fatalf("without repoRoot or installed package, spacedockPackageOK should be false")
		}
		if piRuntimeLaunchReady(check) {
			t.Fatalf("without repoRoot or installed package, runtime should not be launch-ready")
		}
	})

	t.Run("dev override without ensign skill does not satisfy gate", func(t *testing.T) {
		bareRepo := t.TempDir() // no skills/ensign/SKILL.md
		home := t.TempDir()
		statOK := map[string]bool{
			filepath.Join(pkg, "src", "extension", "index.ts"):          true,
			filepath.Join(pkg, "skills", "pi-subagents", "SKILL.md"):    true,
			filepath.Join(pkg, "src", "intercom", "intercom-bridge.ts"): true,
			pkg + "-intercom": true,
			filepath.Join(pkg+"-intercom", "skills", "pi-intercom", "SKILL.md"): true,
		}
		cfg := piRuntimeConfigFromEnv(append(piTestEnv(pkg, home), "SPACEDOCK_REPO_ROOT="+bareRepo), "/non-repo-cwd", "")
		check := checkPiRuntime(&fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOK,
			packageStatus: piPackageStatus{},
		}, cfg)
		if check.spacedockPackageOK {
			t.Fatalf("dev override without ensign skill must not satisfy spacedockPackageOK")
		}
	})
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func writePiSkillFixtures(t *testing.T, repo string) {
	t.Helper()
	writeFileWithDirs(t, filepath.Join(repo, "skills", "first-officer", "SKILL.md"), "---\nname: first-officer\ndescription: test\n---\n")
	writeFileWithDirs(t, filepath.Join(repo, "skills", "ensign", "SKILL.md"), "---\nname: ensign\ndescription: test\n---\n")
}

func writePiSubagentsFixtures(t *testing.T, pkg string) {
	t.Helper()
	writeFileWithDirs(t, filepath.Join(pkg, "src", "extension", "index.ts"), "export default function() {}\n")
	writeFileWithDirs(t, filepath.Join(pkg, "skills", "pi-subagents", "SKILL.md"), "---\nname: pi-subagents\ndescription: test\n---\n")
}

func writeFileWithDirs(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, path, content)
}

func statOKForPiResources(repo, pkg string) map[string]bool {
	return map[string]bool{
		filepath.Join(pkg, "src", "extension", "index.ts"):          true,
		filepath.Join(pkg, "skills", "pi-subagents", "SKILL.md"):    true,
		filepath.Join(pkg, "src", "intercom", "intercom-bridge.ts"): true,
		filepath.Join(repo, "skills", "first-officer", "SKILL.md"):  true,
		filepath.Join(repo, "skills", "ensign", "SKILL.md"):         true,
		// The repo (dev-override) extension: the ready gate Stats the
		// effective package root's extension (AC-5a).
		filepath.Join(repo, ".pi", "extensions", "spacedock.ts"): true,
		// healthyPiPackageStatus's package root — same gate for the installed
		// package root.
		filepath.Join("/pkg-store/spacedock", ".pi", "extensions", "spacedock.ts"): true,
		pkg + "-intercom": true,
		filepath.Join(pkg+"-intercom", "skills", "pi-intercom", "SKILL.md"): true,
	}
}

func piHealthyPathFixtures() map[string]string {
	return map[string]string{
		"pi":        "/bin/pi",
		"safehouse": "/bin/safehouse",
	}
}

func piTestEnv(pkg, home string) []string {
	return []string{
		"PI_SUBAGENTS_PACKAGE_ROOT=" + pkg,
		"PI_INTERCOM_PACKAGE_ROOT=" + pkg + "-intercom",
		"HOME=" + home,
	}
}

// piSafehouseReadyOps builds a fakePiRuntimeOps with the healthy path fixtures
// (pi + safehouse resolvable) and the stat set for the repo/pkg resources, so a
// safehouse-wrapped runPi reaches the launch seam.
func piSafehouseReadyOps(repo, pkg string) *fakePiRuntimeOps {
	return &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
}

// piSafehouseInnerArgv returns the inner pi argv (after the safehouse `--`
// separator) from a wrapped launch argv, or argv unchanged when unwrapped.
func piSafehouseInnerArgv(argv []string) []string {
	if len(argv) == 0 || argv[0] != "safehouse" {
		return argv
	}
	for i, tok := range argv {
		if tok == "--" {
			return argv[i+1:]
		}
	}
	return nil
}

// piSafehouseExtra returns the safehouse `extra` slot (tokens between
// --trust-workdir-config and the `--` separator) from a wrapped launch argv.
func piSafehouseExtra(argv []string) []string {
	if len(argv) < 2 || argv[0] != "safehouse" || argv[1] != "--trust-workdir-config" {
		return nil
	}
	for i := 2; i < len(argv); i++ {
		if argv[i] == "--" {
			return argv[2:i]
		}
	}
	return nil
}

// TestPiFrontDoorAcceptsSafehouseFlags pins AC-1: parsePiFrontDoorArgs accepts the
// four --safehouse-* flags (space/= /repeatable) without error and populates
// fd.forceSafehouse + fd.safehouseFlags with the re-prefixed tokens. Today pflag
// rejects these with exit 2.
func TestPiFrontDoorAcceptsSafehouseFlags(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		forceSafehouse bool
		safehouseFlags []string
	}{
		{"bare-safehouse", []string{"--safehouse"}, true, nil},
		{"enable-equals", []string{"--safehouse-enable=ssh"}, false, []string{"enable=ssh"}},
		{"enable-space", []string{"--safehouse-enable", "ssh"}, false, []string{"enable=ssh"}},
		{"enable-comma", []string{"--safehouse-enable=ssh,docker"}, false, []string{"enable=ssh,docker"}},
		{"add-dirs-equals", []string{"--safehouse-add-dirs=~/scratch"}, false, []string{"add-dirs=~/scratch"}},
		{"add-dirs-space", []string{"--safehouse-add-dirs", "~/scratch"}, false, []string{"add-dirs=~/scratch"}},
		{"add-dirs-repeat", []string{"--safehouse-add-dirs", "/a", "--safehouse-add-dirs", "/b"}, false, []string{"add-dirs=/a", "add-dirs=/b"}},
		{"add-dirs-ro-equals", []string{"--safehouse-add-dirs-ro=~/ro"}, false, []string{"add-dirs-ro=~/ro"}},
		{"add-dirs-ro-space", []string{"--safehouse-add-dirs-ro", "~/ro"}, false, []string{"add-dirs-ro=~/ro"}},
		{"all-knobs", []string{"--safehouse", "--safehouse-enable", "ssh", "--safehouse-add-dirs", "~/scratch", "--safehouse-add-dirs-ro", "~/ro"}, true, []string{"enable=ssh", "add-dirs=~/scratch", "add-dirs-ro=~/ro"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fd, _, err := parsePiFrontDoorArgs(tc.args)
			if err != nil {
				t.Fatalf("parsePiFrontDoorArgs(%v) err = %v (today pflag rejects these with exit 2)", tc.args, err)
			}
			if fd.forceSafehouse != tc.forceSafehouse {
				t.Errorf("forceSafehouse = %v, want %v", fd.forceSafehouse, tc.forceSafehouse)
			}
			if !equalArgv(fd.safehouseFlags, tc.safehouseFlags) {
				t.Errorf("safehouseFlags = %v, want %v", fd.safehouseFlags, tc.safehouseFlags)
			}
		})
	}
}

// TestPiFrontDoorWrapsWhenKnobPresent pins AC-2: with a knob present and a
// resolvable safehouse binary, runPi produces safehouse --trust-workdir-config
// <TranslateFlags-extra> -- pi …, where extra matches safehouse.TranslateFlags'
// output (the same function claude/codex use).
func TestPiFrontDoorWrapsWhenKnobPresent(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := piSafehouseReadyOps(repo, pkg)
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--safehouse-add-dirs=/a"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(ops.launched) == 0 || ops.launched[0] != "safehouse" {
		t.Fatalf("expected safehouse-wrapped argv, got %v", ops.launched)
	}
	wantExtra, err := safehouse.TranslateFlags([]string{"add-dirs=/a"})
	if err != nil {
		t.Fatalf("TranslateFlags err = %v", err)
	}
	// Feedback cycle 2: pi now mirrors claude/codex's wrap plumbing —
	// launcherBinEnvPassFlags() prefixes the safehouse extra with --env-pass
	// SPACEDOCK_BIN (AC-4) so the launcher the helper calls resolve survives the
	// sandbox boundary. The operator's add-dirs follows. Pi adds its own
	// --env-pass PI_SPACEDOCK_LAUNCH so the gated extension's marker survives
	// the sandbox boundary too.
	wantExtra = append(launcherBinEnvPassFlags(), wantExtra...)
	wantExtra = append(wantExtra, "--env-pass", piLaunchMarkerEnv, "--env-pass", piLaunchTasklessEnv)
	if !equalArgv(piSafehouseExtra(ops.launched), wantExtra) {
		t.Fatalf("extra = %v, want TranslateFlags output prefixed by launcherBinEnvPassFlags %v\nargv=%v", piSafehouseExtra(ops.launched), wantExtra, ops.launched)
	}
	inner := piSafehouseInnerArgv(ops.launched)
	// The retired --skill <repo>/skills/{first-officer,ensign} flags are absent
	// (eq's install-managed merge: the installed package's extension discovers
	// the Spacedock skills via resources_discover); only the pi-subagents skill is
	// passed, matching TestPiFrontDoorLaunchesWithNativeResourcePaths.
	wantPrefix := []string{
		"pi",
		"--extension", filepath.Join(pkg, "src", "extension", "index.ts"),
		"--skill", filepath.Join(pkg, "skills", "pi-subagents"),
	}
	if len(inner) != len(wantPrefix) {
		t.Fatalf("wrapped inner argv = %v, want exactly the prefix (no launch prompt — argv-only now, no task in this launch)", inner)
	}
	for i, want := range wantPrefix {
		if inner[i] != want {
			t.Fatalf("wrapped inner[%d]=%q want %q\ninner=%v", i, inner[i], want, inner)
		}
	}
}

// TestPiFrontDoorPlainWhenNoTrigger pins AC-2: with no knob, no --safehouse, and
// no .safehouse profile, runPi produces the plain pi argv (no safehouse prefix).
func TestPiFrontDoorPlainWhenNoTrigger(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := piSafehouseReadyOps(repo, pkg)
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"--plugin-dir", repo}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
	if len(ops.launched) == 0 || ops.launched[0] != "pi" {
		t.Fatalf("expected plain pi argv (no safehouse prefix), got %v", ops.launched)
	}
}

// TestPiFrontDoorWrapsWhenSafehouseProfileAlone pins AC-2: a .safehouse profile
// alone (no flags) also triggers the wrap; the extra slot carries only the
// env-pass prefix (launcherBinEnvPassFlags, AC-4) — no operator add-dirs/enables.
func TestPiFrontDoorWrapsWhenSafehouseProfileAlone(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	dir := safehouseFixtureDir(t)
	ops := piSafehouseReadyOps(repo, pkg)
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"--plugin-dir", repo}, dir, piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
	if len(ops.launched) == 0 || ops.launched[0] != "safehouse" {
		t.Fatalf("expected safehouse-wrapped argv for .safehouse profile alone, got %v", ops.launched)
	}
	wantProfileAlone := append(launcherBinEnvPassFlags(), "--env-pass", piLaunchMarkerEnv, "--env-pass", piLaunchTasklessEnv)
	if !equalArgv(piSafehouseExtra(ops.launched), wantProfileAlone) {
		t.Fatalf("extra should be launcherBinEnvPassFlags + the pi launch marker env-pass for profile-alone wrap, got %v", piSafehouseExtra(ops.launched))
	}
}

// TestPiFrontDoorUnknownSafehouseKeyErrors pins AC-2: an unknown --safehouse-* key
// exits non-zero with no Launch, mirroring TestUnknownSafehouseKeyErrors.
func TestPiFrontDoorUnknownSafehouseKeyErrors(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := piSafehouseReadyOps(repo, pkg)
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"--safehouse-bogus=x"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("exit=0 want non-zero for unknown --safehouse-* key")
	}
	if ops.launched != nil {
		t.Fatalf("Launch invoked on unknown --safehouse-* key: %v", ops.launched)
	}
	if !strings.Contains(stderr.String(), "--safehouse-bogus") {
		t.Fatalf("error does not name the knob --safehouse-bogus: %q", stderr.String())
	}
}

// TestPiFrontDoorWrapInnerEqualsUnwrapped pins AC-3 (PI-SPECIFIC DECISION 1): the
// wrapped inner argv (tokens after safehouse --trust-workdir-config [extra] --) is
// byte-identical to the unwrapped pi argv. No --dangerously-skip-permissions-
// equivalent and no --tools/--exclude-tools default is injected by the wrap; the
// wrap prefix is the only difference.
func TestPiFrontDoorWrapInnerEqualsUnwrapped(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)

	// Unwrapped: no knob, no --safehouse, no .safehouse profile.
	plainOps := piSafehouseReadyOps(repo, pkg)
	var stdout, stderr bytes.Buffer
	if code := runPi(context.Background(), []string{"--plugin-dir", repo}, t.TempDir(), piTestEnv(pkg, t.TempDir()), plainOps, &stdout, &stderr); code != 0 {
		t.Fatalf("plain exit=%d stderr=%q", code, stderr.String())
	}
	unwrapped := plainOps.launched

	// Wrapped: a knob triggers the wrap, no .safehouse profile.
	wrapOps := piSafehouseReadyOps(repo, pkg)
	var wout, werr bytes.Buffer
	if code := runPi(context.Background(), []string{"--plugin-dir", repo, "--safehouse-add-dirs=/tmp/probe"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), wrapOps, &wout, &werr); code != 0 {
		t.Fatalf("wrapped exit=%d stderr=%q", code, werr.String())
	}
	wrapped := wrapOps.launched
	if len(wrapped) == 0 || wrapped[0] != "safehouse" {
		t.Fatalf("expected safehouse-wrapped argv, got %v", wrapped)
	}
	inner := piSafehouseInnerArgv(wrapped)
	if !equalArgv(inner, unwrapped) {
		t.Fatalf("wrapped inner argv != unwrapped pi argv:\n wrapped inner=%v\n unwrapped    =%v", inner, unwrapped)
	}
	for _, banned := range []string{"--dangerously-skip-permissions", "--dangerously-bypass-approvals-and-sandbox"} {
		for _, tok := range inner {
			if tok == banned || strings.HasPrefix(tok, banned+"=") {
				t.Fatalf("wrap injected a permission-mode flag %q into the inner argv: %v", banned, inner)
			}
		}
	}
}

// TestPiHelpCarriesSafehouseDetail pins AC-4: `spacedock pi --help` (exit 0)
// carries the four --safehouse-* usages + the --safehouse-add-dirs ~/scratch
// example + the -- forwarding note, and does NOT declare --skip-compat-check / --no-install
// (pi has no version gate).
func TestPiHelpCarriesSafehouseDetail(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"pi", "--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"--safehouse",
		"--safehouse-enable",
		"--safehouse-add-dirs",
		"--safehouse-add-dirs-ro",
		"--plugin-dir",
		"--safehouse-add-dirs ~/scratch",
		"forward verbatim",
		"Examples:",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("pi --help missing %q:\n%s", want, out)
		}
	}
	for _, notWant := range []string{"--skip-compat-check", "--no-install"} {
		if strings.Contains(out, notWant) {
			t.Errorf("pi --help should not declare %q:\n%s", notWant, out)
		}
	}
}

// writePiSettingsJSON writes a settings.json with the given package source
// entries under <agentDir>/settings.json and creates the resolved package
// directory trees (package.json with name + pi.skills) for each npm: entry.
func writePiSettingsJSON(t *testing.T, agentDir, home string, packages []string) {
	t.Helper()
	for _, src := range packages {
		if !strings.HasPrefix(src, "npm:") {
			continue
		}
		name := parseNpmPackageName(strings.TrimPrefix(src, "npm:"))
		root := filepath.Join(agentDir, "npm", "node_modules", name)
		skills := []string{"skills"}
		var piSkills string
		for _, s := range skills {
			if piSkills != "" {
				piSkills += ","
			}
			piSkills += "\"" + s + "\""
		}
		pkgJSON := "{\"name\":\"" + name + "\",\"pi\":{\"skills\":[" + piSkills + "]}}"
		writeFileWithDirs(t, filepath.Join(root, "package.json"), pkgJSON)
		writeFileWithDirs(t, filepath.Join(root, "skills", "ensign", "SKILL.md"), "---\nname: ensign\n---\n")
		writeFileWithDirs(t, filepath.Join(root, "skills", "first-officer", "SKILL.md"), "---\nname: first-officer\n---\n")
	}
	var entries string
	for i, src := range packages {
		if i > 0 {
			entries += ","
		}
		entries += "\"" + src + "\""
	}
	settings := "{\"packages\":[" + entries + "]}"
	writeFileWithDirs(t, filepath.Join(agentDir, "settings.json"), settings)
}

// TestPiSpacedockPackageStatus_SubagentsRegistered pins the new detection: the
// settings.json scan sets subagentsRegistered iff npm:pi-subagents is in packages.
func TestPiSpacedockPackageStatus_SubagentsRegistered(t *testing.T) {
	home := t.TempDir()
	for _, tc := range []struct {
		name     string
		packages []string
		want     bool
	}{
		{"with pi-subagents", []string{"npm:spacedock", "npm:pi-subagents"}, true},
		{"without pi-subagents", []string{"npm:spacedock"}, false},
		{"pi-subagents only", []string{"npm:pi-subagents"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			agentDir := filepath.Join(home, ".pi", "agent-"+strings.ReplaceAll(tc.name, " ", "-"))
			writePiSettingsJSON(t, agentDir, home, tc.packages)
			if got := piSpacedockPackageStatus(agentDir, home).subagentsRegistered; got != tc.want {
				t.Fatalf("subagentsRegistered = %v, want %v (packages=%v)", got, tc.want, tc.packages)
			}
		})
	}
}

// TestPiResumeSuppressesBootstrapPrompt pins AC-3's resume arm: a Pi launch
// with a resume token in the passthrough (--resume, --resume=<id>, -r,
// --continue, -c) appends NO launch prompt, and no launch shape ever carries
// $spacedock: contract syntax (the extension's gated injection owns the
// contract). The non-resume cases are the independent baselines that can move
// the wrong way: a task still lands as the bare launch prompt, and a
// non-task passthrough appends nothing.
func TestPiResumeSuppressesBootstrapPrompt(t *testing.T) {
	resumeTokens := []string{
		"--resume",
		"--resume=abc123",
		"-r",
		"--continue",
		"-c",
	}
	for _, token := range resumeTokens {
		t.Run("resume/"+token, func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			ops := &fakePiRuntimeOps{
				lookPath:      piHealthyPathFixtures(),
				statOK:        statOKForPiResources(repo, pkg),
				packageStatus: healthyPiPackageStatus(),
			}
			var stdout, stderr bytes.Buffer
			args := []string{"--plugin-dir", repo, "--", token}
			code := runPi(context.Background(), args, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			for _, tok := range ops.launched {
				if strings.Contains(tok, "$spacedock:") {
					t.Fatalf("resume token %q: argv contains $spacedock: contract syntax: %v", token, ops.launched)
				}
			}
		})
	}

	nonResumeCases := []struct {
		name     string
		passthru []string
	}{
		{"model_flag", []string{"--model", "google/gemini"}},
		{"task_string", []string{"review this code"}},
	}
	for _, tc := range nonResumeCases {
		t.Run("nonresume/"+tc.name, func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			ops := &fakePiRuntimeOps{
				lookPath:      piHealthyPathFixtures(),
				statOK:        statOKForPiResources(repo, pkg),
				packageStatus: healthyPiPackageStatus(),
			}
			var stdout, stderr bytes.Buffer
			args := append([]string{"--plugin-dir", repo, "--"}, tc.passthru...)
			code := runPi(context.Background(), args, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			for _, tok := range ops.launched {
				if strings.Contains(tok, "$spacedock:") {
					t.Fatalf("non-resume passthrough %v: argv contains $spacedock: contract syntax: %v", tc.passthru, ops.launched)
				}
			}
			last := ops.launched[len(ops.launched)-1]
			switch tc.name {
			case "task_string":
				if last != "review this code" {
					t.Fatalf("non-resume task: last argv = %q, want the bare task (no contract sentence)", last)
				}
			case "model_flag":
				// A flag passthrough appends no prompt: the argv tail stays the flags.
				if last != "google/gemini" {
					t.Fatalf("non-resume flag passthrough: unexpected argv tail %q (a prompt was appended?)", last)
				}
			}
		})
	}
}

// hasLaunchMarkerEnv reports whether env carries PI_SPACEDOCK_LAUNCH=1.
func hasLaunchMarkerEnv(env []string) bool {
	for _, kv := range env {
		if kv == piLaunchMarkerEnv+"=1" {
			return true
		}
	}
	return false
}

// envHasKV reports whether env carries the exact KEY=VALUE entry.
func envHasKV(env []string, kv string) bool {
	for _, e := range env {
		if e == kv {
			return true
		}
	}
	return false
}

// TestRunPi_SetsTasklessMarkerIffNoTaskOrResume pins the taskless-greet env
// gate: PI_SPACEDOCK_LAUNCH_TASKLESS=1 EXACTLY when the launch carries no
// operator task and no resume passthrough (the only shape where nothing else
// triggers a first model request), and "0" otherwise — the deterministic "0"
// also overrides any operator shell value (os/exec keeps the last entry per
// key). Falsifying edits: gating on hasTask alone (resume case fails),
// dropping the append (taskless case fails), or removing the "0" arm (the
// taskful/resume cases fail).
func TestRunPi_SetsTasklessMarkerIffNoTaskOrResume(t *testing.T) {
	cases := []struct {
		name string
		args []string
		wrap bool
		want string
	}{
		{"taskless_installed", []string{"--", "--version"}, false, "1"},
		{"taskless_dev_override", []string{"--plugin-dir", "REPO", "--", "--version"}, false, "1"},
		{"taskful", []string{"review this", "--plugin-dir", "REPO"}, false, "0"},
		{"resume_passthrough", []string{"--plugin-dir", "REPO", "--", "--resume"}, false, "0"},
		{"taskless_wrap", []string{"--plugin-dir", "REPO"}, true, "1"},
		{"taskful_wrap", []string{"review this", "--plugin-dir", "REPO"}, true, "0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			writeFileWithDirs(t, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			var ops *fakePiRuntimeOps
			var dir string
			if tc.wrap {
				ops = piSafehouseReadyOps(repo, pkg)
				dir = safehouseFixtureDir(t)
			} else {
				ops = &fakePiRuntimeOps{
					lookPath:      piHealthyPathFixtures(),
					statOK:        statOKForPiResources(repo, pkg),
					packageStatus: healthyPiPackageStatus(),
				}
				dir = t.TempDir()
			}
			var stdout, stderr bytes.Buffer
			args := make([]string, 0, len(tc.args))
			for _, a := range tc.args {
				if a == "REPO" {
					args = append(args, repo)
					continue
				}
				args = append(args, a)
			}
			code := runPi(context.Background(), args, dir, piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			if !envHasKV(ops.launchedEnv, piLaunchTasklessEnv+"="+tc.want) {
				t.Fatalf("%s: launched env missing %s=%s: %v", tc.name, piLaunchTasklessEnv, tc.want, ops.launchedEnv)
			}
		})
	}
}

// TestPiFrontDoorTasklessArgvCarriesNoGreetText pins AC-3's taskless arm: the
// taskless marker is env-only — the greet/bootstrap text NEVER appears in the
// launch argv, and a taskless launch appends no task positional (argv shape
// unchanged vs pre-taskless-launch), while a taskful launch keeps the bare
// operator task as the last argv token. Falsifying edit: any frontdoor-owned
// greet text appended to argv (either arm fails).
func TestPiFrontDoorTasklessArgvCarriesNoGreetText(t *testing.T) {
	for _, tc := range []struct {
		name    string
		taskful bool
	}{
		{"taskless", false},
		{"taskful", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			writeFileWithDirs(t, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			ops := &fakePiRuntimeOps{
				lookPath:      piHealthyPathFixtures(),
				statOK:        statOKForPiResources(repo, pkg),
				packageStatus: healthyPiPackageStatus(),
			}
			var stdout, stderr bytes.Buffer
			args := []string{"--plugin-dir", repo}
			if tc.taskful {
				args = append(args, "review this")
			}
			code := runPi(context.Background(), args, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			joined := strings.Join(ops.launched, " ")
			for _, banned := range []string{"SPACEDOCK-FO-BOOTSTRAP", "EXTREMELY_IMPORTANT", "greet"} {
				if strings.Contains(joined, banned) {
					t.Fatalf("%s: launch argv carries greet/bootstrap text %q: %v", tc.name, banned, ops.launched)
				}
			}
			if tc.taskful {
				if last := ops.launched[len(ops.launched)-1]; last != "review this" {
					t.Fatalf("taskful launch must keep the bare operator task as the last argv token, got %q (argv=%v)", last, ops.launched)
				}
				if !envHasKV(ops.launchedEnv, piLaunchTasklessEnv+"=0") {
					t.Fatalf("taskful launch env must carry %s=0: %v", piLaunchTasklessEnv, ops.launchedEnv)
				}
			} else {
				// Taskless: no task positional appended — the argv still ends
				// with the dev-override skill flag, exactly as before.
				if last := ops.launched[len(ops.launched)-1]; last != filepath.Join(repo, "skills") {
					t.Fatalf("taskless launch must append no task positional, argv ends %q (argv=%v)", last, ops.launched)
				}
				if !envHasKV(ops.launchedEnv, piLaunchTasklessEnv+"=1") {
					t.Fatalf("taskless launch env must carry %s=1: %v", piLaunchTasklessEnv, ops.launchedEnv)
				}
			}
		})
	}
}

// TestRunPi_SetsLaunchMarkerEnvOnEveryLaunchShape pins AC-3's dev-override arm:
// PI_SPACEDOCK_LAUNCH=1 is in the launched child env with AND without repoRoot
// set. Falsifying edit: drop the env append in runPi.
func TestRunPi_SetsLaunchMarkerEnvOnEveryLaunchShape(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"installed", []string{"--", "--version"}},
		{"dev_override", []string{"--plugin-dir", "REPO", "--", "--version"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			writeFileWithDirs(t, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			ops := &fakePiRuntimeOps{
				lookPath:      piHealthyPathFixtures(),
				statOK:        statOKForPiResources(repo, pkg),
				packageStatus: healthyPiPackageStatus(),
			}
			var stdout, stderr bytes.Buffer
			args := make([]string, 0, len(tc.args))
			for _, a := range tc.args {
				if a == "REPO" {
					args = append(args, repo)
					continue
				}
				args = append(args, a)
			}
			code := runPi(context.Background(), args, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			if !hasLaunchMarkerEnv(ops.launchedEnv) {
				t.Fatalf("%s: launched env missing %s=1: %v", tc.name, piLaunchMarkerEnv, ops.launchedEnv)
			}
		})
	}
}

// TestRunPi_ReadyGateRequiresSpacedockExtension pins AC-5a: the ready gate
// refuses when the package root's .pi/extensions/spacedock.ts is missing.
func TestRunPi_ReadyGateRequiresSpacedockExtension(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	statOK := statOKForPiResources(repo, pkg)
	delete(statOK, filepath.Join("/pkg-store/spacedock", ".pi", "extensions", "spacedock.ts"))
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOK,
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer
	code := runPi(context.Background(), []string{"--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit=%d want 1 (missing installed Spacedock extension must refuse): stdout=%q", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "MISSING Spacedock extension") {
		t.Fatalf("doctor report missing MISSING Spacedock extension line:\n%s", stdout.String())
	}
	if len(ops.launched) != 0 {
		t.Fatalf("refused runtime must not launch: %v", ops.launched)
	}
}

// TestRunPi_ReadyGateRequiresFirstOfficerSkill pins AC-5b: the ready gate
// refuses when first-officer is not discoverable from the package.
func TestRunPi_ReadyGateRequiresFirstOfficerSkill(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	status := healthyPiPackageStatus()
	status.firstOfficerDiscoverable = false
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: status,
	}
	var stdout, stderr bytes.Buffer
	code := runPi(context.Background(), []string{"--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit=%d want 1 (undiscoverable first-officer must refuse): stdout=%q", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "MISSING Spacedock first-officer skill") {
		t.Fatalf("doctor report missing MISSING first-officer line:\n%s", stdout.String())
	}
	if len(ops.launched) != 0 {
		t.Fatalf("refused runtime must not launch: %v", ops.launched)
	}
}

// TestRunPi_ReadyGateRefusesSubFloorPiVersion pins AC-5's pi >= 0.83.0 floor,
// read from the binary's `pi --version` only; unparseable output fails closed.
func TestRunPi_ReadyGateRefusesSubFloorPiVersion(t *testing.T) {
	cases := []struct {
		version string
		want    int
	}{
		{"0.82.9", 1},
		{"0.83.0", 0},
		{"0.85.1", 0},
		{"1.0.0", 0},
		{"dev\n", 1},   // unparseable fails closed
		{"garbage", 1}, // unparseable fails closed
	}
	for _, tc := range cases {
		t.Run(tc.version, func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			ops := &fakePiRuntimeOps{
				lookPath:      piHealthyPathFixtures(),
				statOK:        statOKForPiResources(repo, pkg),
				packageStatus: healthyPiPackageStatus(),
				piVersionOut:  tc.version,
			}
			var stdout, stderr bytes.Buffer
			code := runPi(context.Background(), []string{"--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
			if code != tc.want {
				t.Fatalf("pi version %q: exit=%d want %d (stdout=%q)", tc.version, code, tc.want, stdout.String())
			}
			if tc.want != 0 && !strings.Contains(stdout.String(), "MISSING pi version") {
				t.Fatalf("sub-floor doctor report missing MISSING pi version line:\n%s", stdout.String())
			}
		})
	}
}

// TestPiVersionAtLeast unit-tests the floor parse directly: bare semver in,
// floor compare, fail closed on garbage.
func TestPiVersionAtLeast(t *testing.T) {
	cases := []struct {
		version string
		want    bool
	}{
		{"0.82.9", false},
		{"0.83.0", true},
		{"0.83.1", true},
		{"0.85.1", true},
		{"1.0.0", true},
		{"0.8.9", false},
		{"0.9.0", false},
		{"", false},
		{"dev", false},
		{"garbage 0.90.0 trailing", true}, // first semver triple in the output wins
	}
	for _, tc := range cases {
		if got := piVersionAtLeast(tc.version, piVersionFloor); got != tc.want {
			t.Errorf("piVersionAtLeast(%q, %q) = %v, want %v", tc.version, piVersionFloor, got, tc.want)
		}
	}
}

// TestRunPi_WarnsOnDuplicateSpacedockRegistration pins AC-5's non-fatal
// duplicate warning: all roots, the first-match winner, the route-aware count.
func TestRunPi_WarnsOnDuplicateSpacedockRegistration(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	status := healthyPiPackageStatus()
	status.spacedockEntries = 2
	status.packageRoots = []string{"/pkg-store/spacedock", "/dev-link/spacedock"}
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: status,
	}
	var stdout, stderr bytes.Buffer
	code := runPi(context.Background(), []string{"--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("duplicate registration is non-fatal: exit=%d stderr=%q", code, stderr.String())
	}
	warnings := stderr.String()
	for _, want := range []string{
		"registered 2 times",
		"/pkg-store/spacedock (wins, first match)",
		"/dev-link/spacedock",
		"up to 4 skill registrations",
		"pi remove",
	} {
		if !strings.Contains(warnings, want) {
			t.Fatalf("duplicate-registration warning missing %q:\n%s", want, warnings)
		}
	}
}

// TestRunPi_WarnsOnDevOverrideDoubleExtension pins AC-5's non-fatal warning on
// the static dev-override double-extension condition (repoRoot + registered).
func TestRunPi_WarnsOnDevOverrideDoubleExtension(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	writeFileWithDirs(t, filepath.Join(repo, ".pi", "extensions", "spacedock.ts"), "export default function(){}\n")
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer
	code := runPi(context.Background(), []string{"--plugin-dir", repo, "--", "--version"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("double extension is non-fatal: exit=%d stderr=%q", code, stderr.String())
	}
	warnings := stderr.String()
	for _, want := range []string{"both extension copies load", repo, "pi remove"} {
		if !strings.Contains(warnings, want) {
			t.Fatalf("double-extension warning missing %q:\n%s", want, warnings)
		}
	}
}

// TestPiDoctorReportsFirstOfficerVersionDuplicates pins AC-7: the first-officer
// skill line, the pi version floor line (sub-floor flagged as the stale
// @mariozechner-org signal with an @earendil-works remedy), the route-aware
// duplicate count and the double-extension line, each with `pi remove`.
func TestPiDoctorReportsFirstOfficerVersionDuplicates(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	home := t.TempDir()

	t.Run("healthy floors and duplicates flagged", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		auth := filepath.Join(home, ".pi", "agent", "auth.json")
		statOK := statOKForPiResources(repo, pkg)
		statOK[auth] = true
		status := healthyPiPackageStatus()
		status.spacedockEntries = 2
		status.packageRoots = []string{"/pkg-store/spacedock", "/dev-link/spacedock"}
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOK,
			packageStatus: status,
			// repoRoot is derived from --plugin-dir; the double-extension
			// condition is repoRoot set AND the package registered.
		}, piTestEnv(pkg, home), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
		}
		out := stdout.String()
		for _, want := range []string{
			"OK Spacedock first-officer skill (package discovery)",
			"OK Spacedock extension",
			"OK pi version: 0.85.1 (floor 0.83.0)",
			"WARN duplicate Spacedock registration: 2 package entries",
			"up to 4 skill registrations",
			"remedy: `pi remove` the stale entry",
			"WARN dev-override double extension",
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("doctor output missing %q:\n%s", want, out)
			}
		}
	})

	t.Run("sub-floor version flagged as stale-org signal", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOKForPiResources(repo, pkg),
			packageStatus: healthyPiPackageStatus(),
			piVersionOut:  "0.73.1\n",
		}, piTestEnv(pkg, home), &stdout, &stderr)
		if code == 0 {
			t.Fatalf("sub-floor doctor must exit non-zero:\n%s", stdout.String())
		}
		out := stdout.String()
		for _, want := range []string{
			"MISSING pi version: 0.73.1 (floor 0.83.0)",
			"@earendil-works",
			"@mariozechner",
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("sub-floor doctor output missing %q:\n%s", want, out)
			}
		}
	})

	t.Run("missing first-officer and extension lines", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		status := healthyPiPackageStatus()
		status.firstOfficerDiscoverable = false
		statOK := statOKForPiResources(repo, pkg)
		delete(statOK, filepath.Join("/pkg-store/spacedock", ".pi", "extensions", "spacedock.ts"))
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath:      piHealthyPathFixtures(),
			statOK:        statOK,
			packageStatus: status,
		}, piTestEnv(pkg, home), &stdout, &stderr)
		if code == 0 {
			t.Fatalf("doctor must exit non-zero when the FO contract path is broken:\n%s", stdout.String())
		}
		out := stdout.String()
		for _, want := range []string{"MISSING Spacedock extension", "MISSING Spacedock first-officer skill"} {
			if !strings.Contains(out, want) {
				t.Fatalf("doctor output missing %q:\n%s", want, out)
			}
		}
	})
}

// TestPiSpacedockPackageStatus_DoubleRegistration pins the route-aware
// duplicate detection at the source: two spacedock entries → spacedockEntries=2,
// packageRoots in settings order (first wins), discovery from the first only.
func TestPiSpacedockPackageStatus_DoubleRegistration(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	for _, root := range []string{first, second} {
		writeFileWithDirs(t, filepath.Join(root, "package.json"), `{"name":"spacedock","pi":{"skills":["./skills"]}}`)
		writeFileWithDirs(t, filepath.Join(root, "skills", "ensign", "SKILL.md"), "---\nname: ensign\ndescription: test\n---\n")
	}
	// Only the first root carries first-officer — the winner's discovery is
	// what counts.
	writeFileWithDirs(t, filepath.Join(first, "skills", "first-officer", "SKILL.md"), "---\nname: first-officer\ndescription: test\n---\n")
	agentDir := t.TempDir()
	writeFile(t, filepath.Join(agentDir, "settings.json"), fmt.Sprintf(`{"packages":["%s","%s"]}`, first, second))
	status := piSpacedockPackageStatus(agentDir, t.TempDir())
	if !status.registered {
		t.Fatalf("spacedock package should be registered")
	}
	if status.spacedockEntries != 2 {
		t.Fatalf("spacedockEntries = %d, want 2", status.spacedockEntries)
	}
	if len(status.packageRoots) != 2 || status.packageRoots[0] != first || status.packageRoots[1] != second {
		t.Fatalf("packageRoots = %v, want [%q %q] (settings order, first wins)", status.packageRoots, first, second)
	}
	if status.packageRoot != first {
		t.Fatalf("packageRoot = %q, want the first entry %q (first-match winner)", status.packageRoot, first)
	}
	if !status.firstOfficerDiscoverable {
		t.Fatalf("first-officer should be discoverable from the winning (first) root")
	}
}
