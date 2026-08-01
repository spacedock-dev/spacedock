//go:build live

package ensigncycle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const piLiveSmokeMarker = "PI-LIVE-SUBAGENT-ENSIGN-SMOKE"
const defaultPiLiveModel = "openrouter/openai/gpt-5.4"

func TestLivePiSubagentEnsignSmoke(t *testing.T) {
	piBin, err := exec.LookPath("pi")
	if err != nil {
		t.Skip("pi not on PATH; install Pi CLI before running the live Pi smoke")
	}
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)
	workflowRoot, stateRoot, entityPath, artifactDir, env := newPiLiveSmokeFixture(t, "pi-subagent-ensign-smoke", repo, piSubagentsRoot, binary)

	envelope := runPiSmokeDispatchBuild(t, binary, workflowRoot, entityPath)
	prompt := piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath, envelope)
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
	assertPiEnsignBootContract(t, workflowRoot, envelope, artifactDir)
}

func TestLivePiFrontDoorSmoke(t *testing.T) {
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)
	workflowRoot, stateRoot, entityPath, artifactDir, env := newPiLiveSmokeFixture(t, "pi-frontdoor-smoke", repo, piSubagentsRoot, binary)

	envelope := runPiSmokeDispatchBuild(t, binary, workflowRoot, entityPath)
	prompt := piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath, envelope)
	runPiLiveCommand(t, artifactDir, workflowRoot, env, binary,
		"pi",
		prompt,
		"--plugin-dir", repo,
		"--",
		"--print",
		"--model", piLiveModelName(),
		"--session-dir", filepath.Join(artifactDir, "sessions"),
	)
	assertPiLiveSmokeResult(t, stateRoot, entityPath, artifactDir)
}

func newPiLiveSmokeFixture(t *testing.T, name, repo, piSubagentsRoot, binary string) (workflowRoot, stateRoot, entityPath, artifactDir string, env []string) {
	t.Helper()
	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	seedPiLiveAuth(t, piHome, os.Getenv("HOME"), os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
	// Patch 3 (validation attempt-1 correction): seed piHome/settings.json with
	// the repo as a path package so pi-subagents' settings-package skill
	// discovery (skills.ts collectSettingsPackageSkillPaths over
	// agentDir/settings.json) resolves the basename skill "ensign"; auth-only
	// piHome boots the child contract-free (skills: []).
	writeFile(t, filepath.Join(piHome, "settings.json"), fmt.Sprintf("{\"packages\":[%q]}\n", "file:"+repo))
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
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

func TestPiLiveSmokePromptRequiresExactStageReportHeading(t *testing.T) {
	envelope := piSmokeEnvelope{Agent: "worker", Skill: "ensign", Prompt: "Read /tmp/spacedock-dispatch/x.md and treat its content as your assignment."}
	prompt := piLiveSmokePrompt("/repo", "/workflow", "/workflow/.spacedock-state", "/workflow/.spacedock-state/pi-live-smoke/index.md", envelope)
	want := "exact heading '## Stage Report: implementation'"
	if !strings.Contains(prompt, want) {
		t.Fatalf("pi live smoke prompt missing %q:\n%s", want, prompt)
	}
}

// piSmokeEnvelope is the dispatch-build --host pi artifact fields the smoke
// forwards to pi-subagents and grades the spawn against.
type piSmokeEnvelope struct {
	Agent        string `json:"agent"`
	Skill        string `json:"skill"`
	Prompt       string `json:"prompt"`
	DispatchFile string `json:"dispatch_file_path"`
}

// runPiSmokeDispatchBuild assembles the real initial-dispatch artifact for the
// smoke entity through `dispatch build --host pi` (the AC-1 spawn source of
// truth) and asserts it carries the default worker/ensign spawn fields.
func runPiSmokeDispatchBuild(t *testing.T, binary, workflowRoot, entityPath string) piSmokeEnvelope {
	t.Helper()
	checklist := []string{
		"- Append a stage report with the exact heading '## Stage Report: implementation' containing the exact marker " + piLiveSmokeMarker + ", at least one '- DONE:' item, and a '### Summary' subsection",
		"- Commit only the entity path in the state checkout with message 'ensign: pi live smoke' (path-scoped git add/commit for pi-live-smoke/index.md)",
	}
	stdin, err := json.Marshal(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   workflowRoot,
		"stage":          "implementation",
		"checklist":      checklist,
		"bare_mode":      true,
		"host":           "pi",
	})
	if err != nil {
		t.Fatal(err)
	}
	// Patches 1+2 (validation attempt-1 correction): the CLI surface is
	// `dispatch build`, and stderr (e.g. the bare-mode advisory) must not
	// contaminate the stdout JSON envelope parse.
	cmd := exec.Command(binary, "dispatch", "build", "--workflow-dir", workflowRoot)
	cmd.Dir = workflowRoot
	cmd.Stdin = strings.NewReader(string(stdin))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("dispatch build --host pi failed: %v\nstderr:\n%s", err, stderr.String())
	}
	out := stdout.Bytes()
	var envelope piSmokeEnvelope
	if err := json.Unmarshal(out, &envelope); err != nil {
		t.Fatalf("dispatch build stdout is not the build envelope: %v\n%s\nstderr:\n%s", err, out, stderr.String())
	}
	if envelope.Agent != "worker" || envelope.Skill != "ensign" {
		t.Fatalf("pi build envelope = agent %q skill %q, want worker/ensign:\n%s", envelope.Agent, envelope.Skill, out)
	}
	if envelope.Prompt == "" || envelope.DispatchFile == "" {
		t.Fatalf("pi build envelope missing prompt/dispatch_file_path:\n%s", out)
	}
	return envelope
}

func piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath string, envelope piSmokeEnvelope) string {
	return fmt.Sprintf(`You are the Spacedock first officer for a live Pi smoke test.

An initial-dispatch artifact was assembled for the entity with `+"`spacedock dispatch build --host pi`"+`; forward it through pi-subagents exactly as emitted — this smoke exists to prove the build artifact's spawn fields drive the worker boot.

  agent: %[5]s
  skill: %[6]s
  task: %[7]s

Use the pi-subagents subagent(...) tool exactly once with those fields verbatim (context must be "fresh", working directory %[2]s). Do not use or mention Claude Agent, SendMessage, TeamCreate, or TeamDelete tools. Do not paraphrase, re-order, or extend the task string.

After subagent(...) returns, you as first officer must verify the entity file %[4]s contains %[8]s and verify the state checkout %[3]s git log contains 'ensign: pi live smoke' over pi-live-smoke/index.md. The stage report must use the exact heading '## Stage Report: implementation' — confirm that too. Exit successfully only after those durable checks pass; your final message names the agent and skill values you passed to subagent(...) and the child's run id.

Reference paths: ensign contract at %[1]s/skills/ensign/SKILL.md; Pi ensign adapter at %[1]s/skills/ensign/references/pi-ensign-runtime.md (the worker's dispatch artifact already points at them).`,
		repo, workflowRoot, stateRoot, entityPath, envelope.Agent, envelope.Skill, envelope.Prompt, piLiveSmokeMarker)
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
	seedPiLiveAuth(t, piHome, realHome, os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
}

func seedPiLiveAuth(t *testing.T, piHome, realHome, openAIAPIKey, required string) {
	t.Helper()
	if realHome != "" {
		authPath := filepath.Join(realHome, ".pi", "agent", "auth.json")
		b, err := os.ReadFile(authPath)
		if err == nil && strings.TrimSpace(string(b)) != "" {
			if err := os.MkdirAll(piHome, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(piHome, "auth.json"), b, 0o600); err != nil {
				t.Fatal(err)
			}
			return
		}
	}
	if strings.TrimSpace(openAIAPIKey) != "" {
		return
	}
	message := "no live Pi auth available: expected ~/.pi/agent/auth.json or OPENAI_API_KEY"
	if required != "" {
		t.Fatal(message + " for the approval-gated pi-live lane")
	}
	t.Skip(message + "; run pi login or set OPENAI_API_KEY to run the live Pi suite")
}

func piLiveEnv(piHome, sessionDir, cleanHome, binaryDir, piSubagentsRoot string) []string {
	env := cleanEnviron(
		"CODEX_THREAD_ID", "CLAUDECODE", "HOME", "PI_CODING_AGENT_DIR",
		"PI_CODING_AGENT_SESSION_DIR", "PI_INTERCOM_PACKAGE_ROOT",
		"PI_SUBAGENTS_PACKAGE_ROOT", "PI_OFFLINE",
	)
	// Optional-adjacent scrub (validation attempt-1 correction): an ambient
	// PI_SUBAGENT_* family leaks hermeticity when the live lane runs nested
	// inside a pi-subagents session (e.g. PI_SUBAGENT_CHILD=1 would exempt the
	// parent FO from its own extension bootstrap).
	env = dropEnvPrefix(env, "PI_SUBAGENT_")
	env = append(env,
		"HOME="+cleanHome,
		"PI_CODING_AGENT_DIR="+piHome,
		"PI_CODING_AGENT_SESSION_DIR="+sessionDir,
		"PI_INTERCOM_PACKAGE_ROOT="+piIntercomPackageRoot(piSubagentsRoot),
		"PI_SUBAGENTS_PACKAGE_ROOT="+piSubagentsRoot,
		"PI_OFFLINE=1",
	)
	return withBinaryOnPath(env, filepath.Join(binaryDir, "spacedock"))
}

func TestPiLiveEnvDropsForeignRuntimeMarkers(t *testing.T) {
	for key, value := range map[string]string{"CODEX_THREAD_ID": "codex", "CLAUDECODE": "claude", "PI_CODING_AGENT": "pi", "PI_CODING_AGENT_DIR": "/parent/pi",
		"PI_CODING_AGENT_SESSION_DIR": "/parent/sessions", "PI_INTERCOM_PACKAGE_ROOT": "/parent/intercom",
		"PI_SUBAGENTS_PACKAGE_ROOT": "/parent/package",
		"PI_OFFLINE":                "0", "HOME": "/parent/home", "OPENAI_API_KEY": "key", "PATH": "/parent/bin"} {
		t.Setenv(key, value)
	}

	env := piLiveEnv("/target/pi", "/target/sessions", "/target/home", "/spacedock/bin", "/target/package")

	want := map[string]string{"CODEX_THREAD_ID": "", "CLAUDECODE": "",
		"PI_CODING_AGENT": "pi", "PI_CODING_AGENT_DIR": "/target/pi",
		"PI_CODING_AGENT_SESSION_DIR": "/target/sessions", "PI_SUBAGENTS_PACKAGE_ROOT": "/target/package",
		"PI_INTERCOM_PACKAGE_ROOT": "/parent/intercom",
		"PI_OFFLINE":               "1", "HOME": "/target/home", "OPENAI_API_KEY": "key",
		"PATH": "/spacedock/bin" + string(os.PathListSeparator) + "/parent/bin"}
	for key, value := range want {
		assertEnvValue(t, env, key, value)
	}
}

// dropEnvPrefix removes every KEY=VALUE entry whose key carries prefix.
func dropEnvPrefix(env []string, prefix string) []string {
	kept := env[:0]
	for _, kv := range env {
		key := kv
		if i := strings.IndexByte(kv, '='); i >= 0 {
			key = kv[:i]
		}
		if strings.HasPrefix(key, prefix) {
			continue
		}
		kept = append(kept, kv)
	}
	return kept
}

func TestPiLiveEnvScrubsAmbientPiSubagentMarkers(t *testing.T) {
	t.Setenv("PI_SUBAGENT_CHILD", "1")
	t.Setenv("PI_SUBAGENT_RUN_ID", "ambient-run")
	t.Setenv("PI_SUBAGENT_DEPTH", "1")

	env := piLiveEnv("/target/pi", "/target/sessions", "/target/home", "/spacedock/bin", "/target/package")

	for _, kv := range env {
		key := kv
		if i := strings.IndexByte(kv, '='); i >= 0 {
			key = kv[:i]
		}
		if strings.HasPrefix(key, "PI_SUBAGENT_") {
			t.Fatalf("piLiveEnv leaked ambient marker %s", kv)
		}
	}
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

func piIntercomPackageRoot(piSubagentsRoot string) string {
	if p := os.Getenv("PI_INTERCOM_PACKAGE_ROOT"); p != "" {
		return p
	}
	return filepath.Join(filepath.Dir(piSubagentsRoot), "pi-intercom")
}

func TestPiIntercomPackageRootDefaultsBesideSubagents(t *testing.T) {
	t.Setenv("PI_INTERCOM_PACKAGE_ROOT", "")
	if got := piIntercomPackageRoot("/packages/pi-subagents"); got != "/packages/pi-intercom" {
		t.Fatalf("piIntercomPackageRoot() = %q, want sibling package", got)
	}
}

func piLiveModelName() string {
	return envOr("SPACEDOCK_PI_LIVE_CHILD_MODEL", defaultPiLiveModel)
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

// assertPiEnsignBootContract is the AC-1 grader: it walks the run's
// .pi-subagents/artifacts records and proves the spawned worker (a) was
// dispatched with the build artifact's agent/skill fields, (b) read
// skills/ensign/SKILL.md among its first five read-type tool calls, and
// (c) never read any path naming first-officer. A graded summary JSON is
// written next to the run artifacts as the durable acceptance trail.
func assertPiEnsignBootContract(t *testing.T, workflowRoot string, envelope piSmokeEnvelope, artifactDir string) {
	t.Helper()
	artifactsDir := filepath.Join(workflowRoot, ".pi-subagents", "artifacts")
	metaPaths, err := filepath.Glob(filepath.Join(artifactsDir, "*_meta.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(metaPaths) != 1 {
		t.Fatalf("expected exactly one spawned worker, found %d meta artifacts in %s", len(metaPaths), artifactsDir)
	}

	type childMeta struct {
		Agent          string   `json:"agent"`
		Skills         []string `json:"skills"`
		Task           string   `json:"task"`
		TranscriptPath string   `json:"transcriptPath"`
	}
	var meta childMeta
	if err := json.Unmarshal([]byte(readFile(t, metaPaths[0])), &meta); err != nil {
		t.Fatalf("worker meta artifact unreadable: %v", err)
	}
	if meta.Agent != envelope.Agent {
		t.Fatalf("spawn consumed agent %q, want the artifact's %q", meta.Agent, envelope.Agent)
	}
	skillForwarded := false
	for _, s := range meta.Skills {
		if s == envelope.Skill {
			skillForwarded = true
		}
	}
	if !skillForwarded {
		t.Fatalf("spawn skills %v do not include the artifact's skill %q", meta.Skills, envelope.Skill)
	}
	if !strings.Contains(meta.Task, envelope.DispatchFile) {
		t.Fatalf("spawn task does not forward the artifact's dispatch-file pointer %s:\n%s", envelope.DispatchFile, tail(meta.Task, 400))
	}

	reads := piTranscriptReadPaths(t, meta.TranscriptPath)
	if len(reads) == 0 {
		t.Fatalf("child transcript %s records no read-type tool calls", meta.TranscriptPath)
	}
	ensignRank := 0 // 1-based rank of the first ensign SKILL.md read; 0 = absent
	for i, p := range reads {
		if strings.Contains(p, "skills/ensign/SKILL.md") {
			ensignRank = i + 1
			break
		}
	}
	if ensignRank == 0 || ensignRank > 5 {
		t.Fatalf("ensign SKILL.md must be among the first five read calls (rank %d); first reads: %v", ensignRank, headStrings(reads, 5))
	}
	foReads := 0
	for _, p := range reads {
		if strings.Contains(p, "first-officer") {
			foReads++
		}
	}
	if foReads != 0 {
		t.Fatalf("child transcript must contain zero first-officer reads, found %d: %v", foReads, reads)
	}

	grade := map[string]any{
		"verdict":                   "pass",
		"worker_transcripts_graded": 1,
		"read_calls":                len(reads),
		"ensign_skill_read_rank":    ensignRank,
		"first_officer_reads":       0,
		"spawn_agent":               meta.Agent,
		"spawn_skills":              meta.Skills,
		"transcript":                meta.TranscriptPath,
	}
	gradeJSON, err := json.MarshalIndent(grade, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	gradePath := filepath.Join(artifactDir, "pi-ensign-boot-grade.json")
	if err := os.WriteFile(gradePath, append(gradeJSON, '\n'), 0o644); err != nil {
		t.Fatalf("write graded transcript artifact: %v", err)
	}
	t.Logf("AC-1 boot contract pass: %d reads, ensign SKILL.md at read #%d, zero first-officer reads; grade artifact %s", len(reads), ensignRank, gradePath)
}

// piTranscriptReadPaths extracts the ordered read-type tool-call paths from a
// pi-subagents child transcript (.jsonl message records whose content blocks
// are read tool calls).
func piTranscriptReadPaths(t *testing.T, transcriptPath string) []string {
	t.Helper()
	var reads []string
	for lineNo, line := range strings.Split(readFile(t, transcriptPath), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record struct {
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("transcript %s line %d is not JSON: %v", transcriptPath, lineNo+1, err)
		}
		var blocks []struct {
			Type      string `json:"type"`
			Name      string `json:"name"`
			Arguments struct {
				Path string `json:"path"`
			} `json:"arguments"`
		}
		if err := json.Unmarshal(record.Message.Content, &blocks); err != nil {
			continue // string content (plain text messages) carries no tool calls
		}
		for _, b := range blocks {
			if b.Type == "toolCall" && b.Name == "read" {
				reads = append(reads, b.Arguments.Path)
			}
		}
	}
	return reads
}

func headStrings(s []string, n int) []string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
