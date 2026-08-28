package cli

import (
	"bytes"
	"context"
	"errors"
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
	if !strings.Contains(prompt, "Use $spacedock:first-officer") || !strings.Contains(prompt, "review this") {
		t.Fatalf("pi prompt missing FO skill or task: %q", prompt)
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

// TestRunPi_DevOverrideWithoutExtensionFallsBackGracefully pins the os.Stat
// guard: when the dev-override is active but <repo>/.pi/extensions/spacedock.ts
// is absent, runPi does NOT add the Spacedock extension/skill flags (no crash).
func TestRunPi_DevOverrideWithoutExtensionFallsBackGracefully(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	// NOTE: .pi/extensions/spacedock.ts intentionally NOT created.
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
	if strings.Contains(joined, wantExt) {
		t.Fatalf("dev-override argv must not reference absent extension %q: %v", wantExt, ops.launched)
	}
	if strings.Contains(joined, "--skill "+filepath.Join(repo, "skills")) {
		t.Fatalf("dev-override argv must not add spacedock skills when extension absent: %v", ops.launched)
	}
	if got := strings.Count(joined, "--skill"); got != 1 {
		t.Fatalf("expected exactly 1 --skill flag (pi-subagents only), got %d: %v", got, ops.launched)
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
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	// No --plugin-dir / SPACEDOCK_REPO_ROOT: installed path (cfg.repoRoot == "").
	env := piTestEnv(pkg, t.TempDir())
	// Provide a healthy registered package status so the launch-ready gate passes.
	status := healthyPiPackageStatus()
	// Point the package root at the repo so the ensign skill Stat check resolves.
	status.packageRoot = repo
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
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
		// The retired repo-path skill checks must not appear.
		for _, notWant := range []string{"OK Spacedock first-officer skill", "OK Spacedock ensign skill"} {
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
	// sandbox boundary. The operator's add-dirs follows.
	wantExtra = append(launcherBinEnvPassFlags(), wantExtra...)
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
	if len(inner) < len(wantPrefix)+1 {
		t.Fatalf("wrapped inner argv too short: %v", inner)
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
	if !equalArgv(piSafehouseExtra(ops.launched), launcherBinEnvPassFlags()) {
		t.Fatalf("extra should be just launcherBinEnvPassFlags for profile-alone wrap, got %v", piSafehouseExtra(ops.launched))
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

// TestPiResumeSuppressesBootstrapPrompt pins AC-1/AC-2: a Pi launch with a resume
// token in the passthrough (--resume, --resume=<id>, -r, --continue, -c) must
// NOT append piBootstrapPrompt to the argv, while a non-resume launch still
// does. The non-resume case is the independent baseline that can move the wrong
// way (if the gate regresses or the non-resume path loses the prompt).
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
				if strings.Contains(tok, piBootstrapPrompt) {
					t.Fatalf("resume token %q: argv contains piBootstrapPrompt: %v", token, ops.launched)
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
			prompt := ops.launched[len(ops.launched)-1]
			if !strings.Contains(prompt, piBootstrapPrompt) {
				t.Fatalf("non-resume passthrough %v: last argv token missing piBootstrapPrompt: %v", tc.passthru, ops.launched)
			}
		})
	}
}
