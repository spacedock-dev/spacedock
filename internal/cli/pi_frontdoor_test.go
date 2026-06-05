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
)

type fakePiRuntimeOps struct {
	lookPath map[string]string
	statOK   map[string]bool
	launched []string
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

func (f *fakePiRuntimeOps) Launch(argv []string) error {
	f.launched = append([]string(nil), argv...)
	return nil
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
		lookPath: piHealthyPathFixtures(),
		statOK:   statOKForPiResources(repo, pkg),
	}
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"review this", "--plugin-dir", repo, "--", "--model", "google/gemini"}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	wantPrefix := []string{
		"pi",
		"--extension", filepath.Join(pkg, "src", "extension", "index.ts"),
		"--skill", filepath.Join(pkg, "skills", "pi-subagents"),
		"--skill", filepath.Join(repo, "skills", "first-officer"),
		"--skill", filepath.Join(repo, "skills", "ensign"),
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
	prompt := ops.launched[len(ops.launched)-1]
	if !strings.Contains(prompt, "Use $spacedock:first-officer") || !strings.Contains(prompt, "review this") {
		t.Fatalf("pi prompt missing FO skill or task: %q", prompt)
	}
}

func TestPiInstallRejectsPluginDir(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ops := &fakeHost{}
	piOps := &fakePiRuntimeOps{}

	code := runInitWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", "/checkout"}, ops, piOps, nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit=%d want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "--plugin-dir is not supported") {
		t.Fatalf("stderr should clearly reject install --plugin-dir, got %q", stderr.String())
	}
	if len(ops.installCmds) != 0 {
		t.Fatalf("install seam called despite rejected --plugin-dir: %v", ops.installCmds)
	}
}

func TestPiInstallAcceptedAndDoesNotUsePluginCommands(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakeHost{}
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "pi"}, ops, &fakePiRuntimeOps{
		lookPath: piHealthyPathFixtures(),
		statOK:   statOKForPiResources(repo, pkg),
	}, append(piTestEnv(pkg, t.TempDir()), "SPACEDOCK_REPO_ROOT="+repo), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(ops.installCmds) != 0 {
		t.Fatalf("install --host pi called host plugin install seam: %v", ops.installCmds)
	}
	out := stdout.String()
	for _, want := range []string{"Pi runtime ready", "pi-subagents", "pi-intercom", pkg, repo, "necessary supervisor-talkback setup prerequisites only"} {
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
			name: "install codex",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runInitWithPi(context.Background(), []string{"--host", "codex", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
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
}

func TestRuntimeSupportDocsKeepPiDoctorVsLiveTalkbackBoundary(t *testing.T) {
	doc, err := os.ReadFile(filepath.Join("..", "..", "docs", "runtime-support.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(doc)
	for _, want := range []string{
		"pi-intercom",
		"intercom bridge",
		"necessary setup checks but insufficient to prove live supervisor talkback",
		"progress update -> decision request -> supervisor reply -> child resume -> durable marker evidence",
		"pi-intercom-supervisor-talkback",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("runtime-support.md missing %q", want)
		}
	}
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
			lookPath: piHealthyPathFixtures(),
			statOK:   statOKForPiResources(repo, pkg),
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
			lookPath: piHealthyPathFixtures(),
			statOK:   statOK,
		}, piTestEnv(pkg, home), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
		}
		out := stdout.String()
		for _, want := range []string{"OK pi CLI", "OK Pi auth", "OK pi-subagents extension", "OK pi-subagents intercom bridge", "OK pi-intercom package root", "OK pi-intercom skill", "OK Spacedock first-officer skill", "OK Spacedock ensign skill", "live child talkback", "durable marker probe"} {
			if !strings.Contains(out, want) {
				t.Fatalf("healthy doctor output missing %q:\n%s", want, out)
			}
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
		"pi": "/bin/pi",
	}
}

func piTestEnv(pkg, home string) []string {
	return []string{
		"PI_SUBAGENTS_PACKAGE_ROOT=" + pkg,
		"PI_INTERCOM_PACKAGE_ROOT=" + pkg + "-intercom",
		"HOME=" + home,
	}
}
