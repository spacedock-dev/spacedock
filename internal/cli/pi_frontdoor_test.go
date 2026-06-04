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
		lookPath: map[string]string{"pi": "/bin/pi"},
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

func TestPiInstallAcceptedAndDoesNotUsePluginCommands(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakeHost{}
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, ops, &fakePiRuntimeOps{
		lookPath: map[string]string{"pi": "/bin/pi"},
		statOK:   statOKForPiResources(repo, pkg),
	}, piTestEnv(pkg, t.TempDir()), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(ops.installCmds) != 0 {
		t.Fatalf("install --host pi called host plugin install seam: %v", ops.installCmds)
	}
	out := stdout.String()
	for _, want := range []string{"Pi runtime ready", "pi-subagents", pkg, repo} {
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

	code := runInitWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{
		lookPath: map[string]string{"pi": "/bin/pi"},
		statOK: map[string]bool{
			filepath.Join(repo, "skills", "first-officer", "SKILL.md"): true,
			filepath.Join(repo, "skills", "ensign", "SKILL.md"):        true,
		},
	}, piTestEnv(pkg, t.TempDir()), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("install --host pi should be idempotent/instructive, exit=%d stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"Pi runtime setup incomplete", "pi install npm:pi-subagents", "PI_SUBAGENTS_PACKAGE_ROOT"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing-subagents output missing %q:\n%s", want, out)
		}
	}
}

func TestNonPiSetupRejectsPluginDir(t *testing.T) {
	for _, tc := range []struct {
		name string
		run  func(hostOps, io.Writer, io.Writer) int
	}{
		{
			name: "install claude",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runInitWithPi(context.Background(), []string{"--host", "claude", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
			},
		},
		{
			name: "install codex",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runInitWithPi(context.Background(), []string{"--host", "codex", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
			},
		},
		{
			name: "doctor claude",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runDoctorWithPi(context.Background(), []string{"--host", "claude", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
			},
		},
		{
			name: "doctor codex",
			run: func(hostOps hostOps, stdout, stderr io.Writer) int {
				return runDoctorWithPi(context.Background(), []string{"--host", "codex", "--plugin-dir", "/checkout"}, hostOps, &fakePiRuntimeOps{}, nil, stdout, stderr)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ops := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := tc.run(ops, &stdout, &stderr)
			if code != 2 {
				t.Fatalf("exit=%d want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stderr.String(), "unknown argument \"--plugin-dir\"") {
				t.Fatalf("stderr should reject --plugin-dir, got %q", stderr.String())
			}
			if len(ops.installCmds) != 0 {
				t.Fatalf("install seam called despite rejected --plugin-dir: %v", ops.installCmds)
			}
		})
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
		for _, want := range []string{"Pi runtime check", "MISSING pi CLI", "MISSING Pi auth", "MISSING pi-subagents"} {
			if !strings.Contains(out, want) {
				t.Fatalf("missing doctor output missing %q:\n%s", want, out)
			}
		}
	})

	t.Run("healthy", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		statOK := statOKForPiResources(repo, pkg)
		statOK[auth] = true
		code := runDoctorWithPi(context.Background(), []string{"--host", "pi", "--plugin-dir", repo}, &fakeHost{}, &fakePiRuntimeOps{
			lookPath: map[string]string{"pi": "/bin/pi"},
			statOK:   statOK,
		}, piTestEnv(pkg, home), &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
		}
		out := stdout.String()
		for _, want := range []string{"OK pi CLI", "OK Pi auth", "OK pi-subagents extension", "OK Spacedock first-officer skill", "OK Spacedock ensign skill"} {
			if !strings.Contains(out, want) {
				t.Fatalf("healthy doctor output missing %q:\n%s", want, out)
			}
		}
	})
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
		filepath.Join(pkg, "src", "extension", "index.ts"):         true,
		filepath.Join(pkg, "skills", "pi-subagents", "SKILL.md"):   true,
		filepath.Join(repo, "skills", "first-officer", "SKILL.md"): true,
		filepath.Join(repo, "skills", "ensign", "SKILL.md"):        true,
	}
}

func piTestEnv(pkg, home string) []string {
	return []string{
		"PI_SUBAGENTS_PACKAGE_ROOT=" + pkg,
		"HOME=" + home,
	}
}
