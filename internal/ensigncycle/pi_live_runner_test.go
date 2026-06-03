//go:build live

package ensigncycle

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const piLiveSmokeMarker = "PI-LIVE-SUBAGENT-ENSIGN-SMOKE"

func TestLivePiSubagentEnsignSmoke(t *testing.T) {
	piBin, err := exec.LookPath("pi")
	if err != nil {
		t.Skip("pi not on PATH; install Pi CLI before running the live Pi smoke")
	}
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)
	workflowRoot, stateRoot, entityPath, artifactDir, env := newPiLiveSmokeFixture(t, "pi-subagent-ensign-smoke", repo, piSubagentsRoot, binary)

	prompt := piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath)
	runPiLiveCommand(t, artifactDir, workflowRoot, env, piBin,
		"--print",
		"--session-dir", filepath.Join(artifactDir, "sessions"),
		"--extension", filepath.Join(piSubagentsRoot, "src", "extension", "index.ts"),
		"--skill", filepath.Join(piSubagentsRoot, "skills", "pi-subagents"),
		"--skill", filepath.Join(repo, "skills", "first-officer"),
		"--skill", filepath.Join(repo, "skills", "ensign"),
		prompt,
	)
	assertPiLiveSmokeResult(t, stateRoot, entityPath, artifactDir)
}

func TestLivePiFrontDoorSmoke(t *testing.T) {
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)
	workflowRoot, stateRoot, entityPath, artifactDir, env := newPiLiveSmokeFixture(t, "pi-frontdoor-smoke", repo, piSubagentsRoot, binary)

	prompt := piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath)
	runPiLiveCommand(t, artifactDir, workflowRoot, env, binary,
		"pi",
		prompt,
		"--plugin-dir", repo,
		"--",
		"--print",
		"--session-dir", filepath.Join(artifactDir, "sessions"),
	)
	assertPiLiveSmokeResult(t, stateRoot, entityPath, artifactDir)
}

func newPiLiveSmokeFixture(t *testing.T, name, repo, piSubagentsRoot, binary string) (workflowRoot, stateRoot, entityPath, artifactDir string, env []string) {
	t.Helper()
	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	seedPiLocalAuth(t, piHome, os.Getenv("HOME"))
	workflowRoot, stateRoot, entityPath = writePiSplitRootSmokeWorkflow(t)
	artifactDir = filepath.Join(piLiveArtifactDir(t, name), "run")
	if err := os.MkdirAll(filepath.Join(artifactDir, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	env = piLiveEnv(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot)
	return workflowRoot, stateRoot, entityPath, artifactDir, env
}

func runPiLiveCommand(t *testing.T, artifactDir, workflowRoot string, env []string, argv ...string) {
	t.Helper()
	stdoutPath := filepath.Join(artifactDir, "pi-stdout.txt")
	stderrPath := filepath.Join(artifactDir, "pi-stderr.txt")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = workflowRoot
	cmd.Env = env
	stdout, err := os.Create(stdoutPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()
	stderr, err := os.Create(stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	runErr := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("pi live smoke timed out; artifacts in %s", artifactDir)
	}
	if runErr != nil {
		t.Fatalf("pi live smoke failed: %v; artifacts in %s\nstderr tail:\n%s", runErr, artifactDir, tail(readFile(t, stderrPath), 4000))
	}
}

func assertPiLiveSmokeResult(t *testing.T, stateRoot, entityPath, artifactDir string) {
	t.Helper()
	entity := readFile(t, entityPath)
	for _, want := range []string{piLiveSmokeMarker, "## Stage Report: implementation", "- DONE:", "### Summary"} {
		if !strings.Contains(entity, want) {
			t.Fatalf("entity missing %q after pi subagent smoke; artifacts in %s\n%s", want, artifactDir, entity)
		}
	}
	log := git(t, stateRoot, "log", "--oneline", "--", "pi-live-smoke", "index.md")
	if !strings.Contains(log, "ensign: pi live smoke") {
		t.Fatalf("state checkout git log missing worker commit; artifacts in %s\n%s", artifactDir, log)
	}
	if strings.TrimSpace(git(t, stateRoot, "status", "--short", "--", "pi-live-smoke", "index.md")) != "" {
		t.Fatalf("state checkout entity has uncommitted changes after worker commit; artifacts in %s\n%s", artifactDir, git(t, stateRoot, "status", "--short"))
	}
}

func piSpacedockBinary(t *testing.T, repo string) string {
	t.Helper()
	if os.Getenv("SPACEDOCK_BIN") != "" {
		return spacedockBinary(t)
	}
	out := filepath.Join(t.TempDir(), "spacedock")
	cmd := exec.Command("go", "build", "-o", out, "./cmd/spacedock")
	cmd.Dir = repo
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build spacedock for Pi live smoke: %v\n%s", err, b)
	}
	return out
}

func piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath string) string {
	return fmt.Sprintf(`You are the Spacedock first officer for a live Pi smoke test.

Use the pi-subagents subagent(...) tool exactly once to dispatch one Pi ensign worker. Do not use or mention Claude Agent, SendMessage, TeamCreate, or TeamDelete tools.

Dispatch a worker with agent "delegate" and this task:

Load and follow the local Spacedock ensign skill at %[1]s/skills/ensign/SKILL.md and the Pi ensign adapter at %[1]s/skills/ensign/references/pi-ensign-runtime.md. This is a split-root Spacedock workflow.

Workflow directory: %[2]s
State checkout: %[3]s
Entity file: %[4]s
Target stage: implementation

Required worker actions:
1. Read the workflow README and entity file.
2. Do not edit YAML frontmatter.
3. Append an implementation stage report to the entity body containing the exact marker %[5]s, at least one '- DONE:' item, and a '### Summary' subsection.
4. Commit only the entity path in the state checkout with message 'ensign: pi live smoke'. Use a path-scoped git add/commit for pi-live-smoke/index.md.
5. Return a concise completion result naming the entity file and commit evidence.

After subagent(...) returns, you as first officer must verify the entity file contains %[5]s and verify the state checkout git log contains 'ensign: pi live smoke'. Exit successfully only after those durable checks pass.`, repo, workflowRoot, stateRoot, entityPath, piLiveSmokeMarker)
}

func writePiSplitRootSmokeWorkflow(t *testing.T) (workflowRoot, stateRoot, entityPath string) {
	t.Helper()
	workflowRoot = t.TempDir()
	stateRoot = filepath.Join(workflowRoot, ".spacedock-state")
	writeFile(t, filepath.Join(workflowRoot, "README.md"), piSplitRootSmokeReadme())
	entityPath = filepath.Join(stateRoot, "pi-live-smoke", "index.md")
	writeFile(t, entityPath, piLiveSmokeEntity())
	gitInit(t, workflowRoot)
	gitInit(t, stateRoot)
	return workflowRoot, stateRoot, entityPath
}

func piSplitRootSmokeReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"state: .spacedock-state\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: implementation\n" +
		"      initial: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Pi Split Root Smoke\n\n" +
		"### implementation\n\n" +
		"Append the live Pi smoke marker to the entity stage report.\n\n" +
		"- **Outputs:** Stage report containing the exact Pi live smoke marker.\n\n" +
		"### done\n\nTerminal state.\n"
}

func piLiveSmokeEntity() string {
	return "---\n" +
		"id: pi-live-smoke\n" +
		"title: Pi Live Smoke\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Pi Live Smoke\n\n" +
		"This entity is mutated only by the Pi subagent live smoke.\n"
}

func seedPiLocalAuth(t *testing.T, piHome, realHome string) {
	t.Helper()
	if realHome == "" {
		t.Skip("no HOME set; cannot locate ~/.pi/agent/auth.json for Pi live smoke")
	}
	authPath := filepath.Join(realHome, ".pi", "agent", "auth.json")
	b, err := os.ReadFile(authPath)
	if err != nil {
		t.Skipf("no live Pi auth available: expected %s; run pi login or provide the auth file", authPath)
	}
	if strings.TrimSpace(string(b)) == "" {
		t.Skipf("live Pi auth file is empty: %s", authPath)
	}
	if err := os.MkdirAll(piHome, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(piHome, "auth.json"), b, 0o600); err != nil {
		t.Fatal(err)
	}
}

func piLiveEnv(piHome, sessionDir, cleanHome, binaryDir, piSubagentsRoot string) []string {
	env := os.Environ()
	env = append(env,
		"HOME="+cleanHome,
		"PI_CODING_AGENT_DIR="+piHome,
		"PI_CODING_AGENT_SESSION_DIR="+sessionDir,
		"PI_SUBAGENTS_PACKAGE_ROOT="+piSubagentsRoot,
		"PI_OFFLINE=1",
	)
	return withBinaryOnPath(env, filepath.Join(binaryDir, "spacedock"))
}

func piSubagentsPackageRoot(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("PI_SUBAGENTS_PACKAGE_ROOT"); p != "" {
		return p
	}
	home := os.Getenv("HOME")
	if home == "" {
		t.Fatal("HOME is empty; set PI_SUBAGENTS_PACKAGE_ROOT to the local pi-subagents package")
	}
	p := filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-subagents")
	if _, err := os.Stat(filepath.Join(p, "src", "extension", "index.ts")); err != nil {
		t.Fatalf("pi-subagents package extension not found at %s: %v; set PI_SUBAGENTS_PACKAGE_ROOT", p, err)
	}
	return p
}

func piLiveArtifactDir(t *testing.T, name string) string {
	t.Helper()
	root := os.Getenv("SPACEDOCK_LIVE_ARTIFACT_DIR")
	if root == "" {
		return t.TempDir()
	}
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}
