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
	"strconv"
	"strings"
	"testing"
	"time"
)

//spacedock:live-proof id=pi-front-door-subagent-dispatch lane=pi-live
func TestLivePiFrontDoorSmoke(t *testing.T) {
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)
	workflowRoot, stateRoot, entityPath, artifactDir, env, model := newPiLiveSmokeFixture(t, "pi-frontdoor-smoke", repo, piSubagentsRoot, binary)

	envelope := runPiSmokeDispatchBuild(t, binary, workflowRoot, entityPath)
	prompt := piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath, envelope)
	runPiLiveCommand(t, artifactDir, workflowRoot, env, binary,
		"pi",
		prompt,
		"--plugin-dir", repo,
		"--",
		"--print",
		"--model", model,
		"--session-dir", filepath.Join(artifactDir, "sessions"),
	)
	assertPiLiveSmokeResult(t, stateRoot, entityPath, artifactDir)
	assertPiEnsignBootContract(t, workflowRoot, envelope, artifactDir)
}

func newPiLiveSmokeFixture(t *testing.T, name, repo, piSubagentsRoot, binary string) (workflowRoot, stateRoot, entityPath, artifactDir string, env []string, model string) {
	t.Helper()
	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	decision := seedPiLiveAuth(t, piHome, os.Getenv("HOME"), os.Getenv("CODEX_AUTH_JSON"), os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
	// Patch 3 (validation attempt-1 correction): seed piHome/settings.json with
	// the repo as a path package so pi-subagents' settings-package skill
	// discovery (skills.ts collectSettingsPackageSkillPaths over
	// agentDir/settings.json) resolves the basename skill "ensign"; auth-only
	// piHome boots the child contract-free (skills: []).
	writeFile(t, filepath.Join(piHome, "settings.json"), fmt.Sprintf("{\"packages\":[%q]}\n", "file:"+repo))
	writePiSubagentsProjectArtifactDir(t, piHome)
	workflowRoot, stateRoot, entityPath = writePiSplitRootSmokeWorkflow(t)
	artifactDir = filepath.Join(piLiveArtifactDir(t, name), "run")
	if err := os.MkdirAll(filepath.Join(artifactDir, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	env = piLiveEnvForAuth(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot, os.Getenv("OPENAI_API_KEY"), decision.mode)
	model = piLiveChildModel(decision)
	return workflowRoot, stateRoot, entityPath, artifactDir, env, model
}

func runPiLiveCommand(t *testing.T, artifactDir, workflowRoot string, env []string, argv ...string) {
	t.Helper()
	stdoutPath := filepath.Join(artifactDir, "pi-stdout.txt")
	stderrPath := filepath.Join(artifactDir, "pi-stderr.txt")
	ctx, cancel := context.WithTimeout(context.Background(), piLiveRunTimeout(10*time.Minute))
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
		t.Fatalf("pi live smoke timed out after the per-run cap (SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES, default 10m); artifacts in %s", artifactDir)
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

// runPiSmokeDispatchBuild assembles the real initial-dispatch artifact for the
// smoke entity through `dispatch build --host pi` (the AC-1 spawn source of
// truth) and asserts it carries the default worker/ensign spawn fields.
func runPiSmokeDispatchBuild(t *testing.T, binary, workflowRoot, entityPath string) piSmokeEnvelope {
	t.Helper()
	checklist := []string{
		"- First read " + filepath.Join(repoRoot(t), "skills", "ensign", "SKILL.md") + " and " + filepath.Join(repoRoot(t), "skills", "ensign", "references", "pi-ensign-runtime.md") + "; then append a stage report with the exact heading '## Stage Report: implementation' containing the exact marker " + piLiveSmokeMarker + ", at least one '- DONE:' item, and a '### Summary' subsection",
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

//spacedock:live-fixture id=pi/split-root-smoke
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
	seedPiLiveAuth(t, piHome, realHome, "", os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
}

func seedPiLiveAuth(t *testing.T, piHome, realHome, oauthJSON, openAIAPIKey, required string) piLiveAuthDecision {
	t.Helper()
	decision := decidePiLiveAuth(oauthJSON, openAIAPIKey, required)
	if decision.mode == piAuthOAuth {
		if err := seedPiOAuthAuth(piHome, oauthJSON); err != nil {
			t.Fatal(err)
		}
		return decision
	}
	if decision.mode == piAuthAPIKey {
		return decision
	}
	if realHome != "" {
		b, err := os.ReadFile(filepath.Join(realHome, ".pi", "agent", "auth.json"))
		if err == nil && strings.TrimSpace(string(b)) != "" {
			if err := os.MkdirAll(piHome, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(piHome, "auth.json"), b, 0o600); err != nil {
				t.Fatal(err)
			}
			// Custom providers (e.g. lunaroute) declare their models in
			// models.json, not auth.json. Mirror it alongside auth.json so a
			// custom-provider SPACEDOCK_PI_LIVE_CHILD_MODEL resolves instead
			// of failing with "Model ... not found".
			if models, merr := os.ReadFile(filepath.Join(realHome, ".pi", "agent", "models.json")); merr == nil && strings.TrimSpace(string(models)) != "" {
				if err := os.WriteFile(filepath.Join(piHome, "models.json"), models, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			return piLiveAuthDecision{mode: piAuthOAuth, model: "openai-codex/gpt-5.6-luna:max"}
		}
	}
	if required != "" {
		t.Fatal(decision.message)
	}
	t.Skip(decision.message)
	return decision
}

// writePiSubagentsProjectArtifactDir opts the live test fixture into the
// "project" artifact dir so spawned-worker meta artifacts land in
// workflowRoot/.pi/subagents/artifacts/ (pi-subagents 0.53.0's
// PROJECT_SUBAGENTS_RELATIVE_DIR is ".pi/subagents", not ".pi-subagents"),
// where the FrontDoorSmoke grader globs for them. pi-subagents 0.53.0
// (#1062) flipped the default to "session" to keep worktrees clean; the live
// tests need a stable, inspectable location.
func writePiSubagentsProjectArtifactDir(t *testing.T, piHome string) {
	t.Helper()
	configDir := filepath.Join(piHome, "extensions", "subagent")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte("{\"artifactDir\":\"project\"}\n"), 0o644); err != nil {
		t.Fatal(err)
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

// piLiveChildModel resolves the model the Pi live child runs on. An operator
// sets SPACEDOCK_PI_LIVE_CHILD_MODEL (provider/model:thinking) to re-run
// journeys against a non-default model; the operator-mirrored auth.json and
// models.json from seedPiLiveAuth make custom providers resolve. Unset, the
// auth decision's model is used, so the CI pi-live lane keeps its
// openai-codex default.
func piLiveChildModel(decision piLiveAuthDecision) string {
	if override := strings.TrimSpace(os.Getenv("SPACEDOCK_PI_LIVE_CHILD_MODEL")); override != "" {
		return override
	}
	return decision.model
}

// piLiveRunTimeout returns the per-run cap for a Pi live journey. It reads
// SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES (a positive integer) and falls back to dflt
// when the env var is unset or invalid. Raise it for slow :max-thinking models so
// multi-dispatch journeys complete to a graded result instead of timing out;
// make the outer `go test -timeout` longer than this per-run cap.
func piLiveRunTimeout(dflt time.Duration) time.Duration {
	if v := os.Getenv("SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Minute
		}
	}
	return dflt
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
	artifactsDir := filepath.Join(workflowRoot, ".pi", "subagents", "artifacts")
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
		// pi-subagents 0.53.0+ redacts the task in the meta artifact
		// ("[prompt redacted]", live Prompt Audit #1021), so the dispatch-file
		// pointer is no longer recoverable from meta.Task. Verify it instead from
		// the parent FO transcript: the subagent toolCall's task argument is the
		// unredacted spawn task the FO forwarded to the worker.
		parentSession := onePiSession(t, filepath.Join(artifactDir, "sessions", "*.jsonl"), "parent")
		dispatchForwarded := false
		for _, task := range piTranscriptSubagentTasks(t, parentSession) {
			if strings.Contains(task, envelope.DispatchFile) {
				dispatchForwarded = true
				break
			}
		}
		if !dispatchForwarded {
			t.Fatalf("spawn task does not forward the artifact's dispatch-file pointer %s (meta task redacted; checked parent transcript):", envelope.DispatchFile)
		}
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

	rootSession := onePiSession(t, filepath.Join(artifactDir, "sessions", "*.jsonl"), "root")
	childSession := onePiSession(t, filepath.Join(artifactDir, "sessions", "*", "*", "run-*", "session.jsonl"), "child")
	grade, err := buildPiFrontDoorEvidenceGrade(rootSession, childSession, true, piBootContractEvidence{
		Agent:                 meta.Agent,
		Skills:                meta.Skills,
		DispatchFileForwarded: strings.Contains(meta.Task, envelope.DispatchFile),
		ReadCallCount:         len(reads),
		EnsignSkillReadRank:   ensignRank,
		FirstOfficerReads:     foReads,
		Transcript:            meta.TranscriptPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	gradeJSON, err := json.MarshalIndent(grade, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	gradePath := filepath.Join(artifactDir, "pi-ensign-boot-grade.json")
	if err := os.WriteFile(gradePath, append(gradeJSON, '\n'), 0o644); err != nil {
		t.Fatalf("write graded transcript artifact: %v", err)
	}
	t.Logf("Pi front-door evidence pass: %d reads, ensign SKILL.md at read #%d, root %s/%dms, child %s/%dms; grade artifact %s",
		len(reads), ensignRank, grade.Root.Model, grade.Root.DurationMS, grade.Child.Model, grade.Child.DurationMS, gradePath)
}

func onePiSession(t *testing.T, pattern, role string) string {
	t.Helper()
	paths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("Pi %s sessions=%d, want exactly one: %v", role, len(paths), paths)
	}
	return paths[0]
}

// piTranscriptReadPaths extracts the ordered read-type tool-call paths from a
// pi-subagents child transcript (.jsonl message records whose content blocks
// are read tool calls).
func piTranscriptToolValues(t *testing.T, transcriptPath, tool, alternate string) []string {
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
				Path    string `json:"path"`
				Command string `json:"command"`
			} `json:"arguments"`
		}
		if err := json.Unmarshal(record.Message.Content, &blocks); err != nil {
			continue // string content (plain text messages) carries no tool calls
		}
		for _, b := range blocks {
			if b.Type == "toolCall" && (b.Name == tool || alternate != "" && b.Name == alternate) {
				value := b.Arguments.Path
				if alternate != "" {
					value = b.Arguments.Command
				}
				if alternate == "" || value != "" {
					reads = append(reads, value)
				}
			}
		}
	}
	return reads
}

func piTranscriptReadPaths(t *testing.T, transcriptPath string) []string {
	return piTranscriptToolValues(t, transcriptPath, "read", "")
}

// piTranscriptSubagentTasks extracts the unredacted task arguments from every
// subagent spawn toolCall in a parent FO transcript. pi-subagents 0.53.0+
// redacts the task in the worker meta artifact ("[prompt redacted]", live
// Prompt Audit #1021), so the dispatch-file pointer the FO forwarded is
// recoverable only from the parent's spawn toolCall, not from meta.Task.
func piTranscriptSubagentTasks(t *testing.T, transcriptPath string) []string {
	t.Helper()
	var tasks []string
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
				Task  string `json:"task"`
				Agent string `json:"agent"`
			} `json:"arguments"`
		}
		if err := json.Unmarshal(record.Message.Content, &blocks); err != nil {
			continue
		}
		for _, b := range blocks {
			if b.Type == "toolCall" && b.Name == "subagent" && b.Arguments.Agent != "" {
				tasks = append(tasks, b.Arguments.Task)
			}
		}
	}
	return tasks
}

func headStrings(s []string, n int) []string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
