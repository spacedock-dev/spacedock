// ABOUTME: Pi dispatch-build host shape — the outer prompt and dispatch body avoid
// ABOUTME: Claude team tool signatures and target pi-subagents style completion.
package dispatch

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/piruntime"
)

func TestBuildPiHostPromptShape(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	var out struct {
		Prompt string `json:"prompt"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	for _, banned := range []string{"Skill(skill=", "Agent(", "SendMessage", "TeamCreate", "TeamDelete"} {
		if strings.Contains(out.Prompt, banned) {
			t.Fatalf("pi prompt must not depend on Claude syntax %q: %q", banned, out.Prompt)
		}
	}
	if !strings.Contains(out.Prompt, "Read ") || !strings.Contains(out.Prompt, "treat its content as your assignment") {
		t.Fatalf("pi prompt should be the read-dispatch-file form: %q", out.Prompt)
	}

	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))
	for _, banned := range []string{"Skill(skill=\"spacedock:ensign\")", "SendMessage(to=\"team-lead\"", "TeamCreate", "TeamDelete", "Agent("} {
		if strings.Contains(body, banned) {
			t.Fatalf("pi dispatch body must omit Claude syntax %q:\n%s", banned, body)
		}
	}
	for _, want := range []string{
		"Read this dispatch file directly",
		"Pi subagent completion result",
		"Do not emit Claude team-tool calls",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("pi dispatch body missing %q:\n%s", want, body)
		}
	}
}

// TestBuildPiHostIgnoresModelWithNote is AC-5: dispatch build --host pi over a
// fable-declaring README exits 0, emits model: null, and prints the
// ignore-with-note stderr line in place of the effective_model line — pi's
// dispatch-settable model space does not (yet) honor claude-enum names, so a
// declared value is dropped rather than rejected.
func TestBuildPiHostIgnoresModelWithNote(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeModelsFable)
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "stagemodel", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "stagemodel",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	assertGolden(t, "build-host-pi-model-ignored", goldenEnvelope{res: normRun(native, root, home)})
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	if !strings.Contains(native.stdout, `"model": null`) {
		t.Errorf("stdout model must be null:\n%s", native.stdout)
	}
	if strings.Contains(native.stderr, "effective_model=") {
		t.Errorf("stderr must not contain the effective_model= line on host=pi:\n%s", native.stderr)
	}
}

func TestBuildPiHostArtifactCarriesCanonicalStageFactsThroughPiWrapper(t *testing.T) {
	root := t.TempDir()
	stateDir := filepath.Join(root, "state-checkout")
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(true))
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	worktreeRel := ".worktrees/spacedock-ensign-canonical-stage"
	worktreePath := filepath.Join(root, worktreeRel)
	if err := os.MkdirAll(worktreePath, 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(stateDir, "canonical-stage", "index.md")
	writeFile(t, entityPath, entityFM("Canonical Stage", "implementation", worktreeRel))
	gitInit(t, root)

	checklist := []string{"- keep the builder assignment", "- write the stage report"}
	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      checklist,
		"bare_mode":      true,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	var out struct {
		Name         *string `json:"name"`
		Description  string  `json:"description"`
		DispatchFile string  `json:"dispatch_file_path"`
		Prompt       string  `json:"prompt"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	if out.Name != nil {
		t.Fatalf("bare Pi dispatch should not require team-style worker name, got %q", *out.Name)
	}
	if out.Description != "Canonical Stage: implementation" {
		t.Fatalf("description = %q", out.Description)
	}
	if !strings.Contains(out.DispatchFile, "spacedock-ensign-canonical-stage-implementation.md") {
		t.Fatalf("dispatch file path does not carry builder-derived slug/stage: %q", out.DispatchFile)
	}

	wrapped := piruntime.SubagentStageDispatch(out.Prompt, "implementation", "canonical-stage implementation")
	if wrapped.Context != "fresh" {
		t.Fatalf("Pi wrapper context = %q, want fresh", wrapped.Context)
	}
	if wrapped.Task != out.Prompt {
		t.Fatalf("Pi wrapper replaced the dispatch-build prompt:\nwant: %s\n got: %s", out.Prompt, wrapped.Task)
	}
	if strings.Contains(wrapped.Task, "acceptance") {
		t.Fatalf("Pi stage wrapper task unexpectedly contains same-agent acceptance contract: %q", wrapped.Task)
	}

	body := readDispatchBody(t, out.DispatchFile)
	for _, want := range []string{
		"You are working on: Canonical Stage",
		"Stage: implementation",
		"Read the entity file at " + entityPath,
		"Your working directory for CODE is " + worktreePath,
		"All CODE reads and writes MUST use paths under " + worktreePath,
		"- keep the builder assignment",
		"- write the stage report",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("Pi dispatch artifact missing builder-derived fact %q:\n%s", want, body)
		}
	}
	assertRenderedStageDefCommand(t, body, root, "implementation")
	if !strings.Contains(wrapped.Task, out.DispatchFile) {
		t.Fatalf("Pi wrapper does not forward dispatch file path %s in task %q", out.DispatchFile, wrapped.Task)
	}
}

func TestPiStageDispatchSmokeUsesBuildArtifactThroughWrapper(t *testing.T) {
	root := t.TempDir()
	stateDir := filepath.Join(root, "state-checkout")
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(true))
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	worktreeRel := ".worktrees/spacedock-ensign-smoke-stage"
	worktreePath := filepath.Join(root, worktreeRel)
	if err := os.MkdirAll(worktreePath, 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(stateDir, "smoke-stage", "index.md")
	writeFile(t, entityPath, entityFM("Smoke Stage", "implementation", worktreeRel))
	gitInit(t, stateDir)

	checklist := []string{"- run the fixture Pi worker", "- append durable stage report"}
	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      checklist,
		"bare_mode":      true,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	var out struct {
		DispatchFile string `json:"dispatch_file_path"`
		Prompt       string `json:"prompt"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	wrapped := piruntime.SubagentStageDispatch(out.Prompt, "implementation", "z6 implementation fixback")
	payload, err := json.Marshal(wrapped)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestPiStageDispatchSmokeFixtureWorker$", "-test.v")
	cmd.Env = append(os.Environ(),
		"SPACEDOCK_PI_FIXTURE_WORKER=1",
		"SPACEDOCK_PI_WRAPPER_JSON="+string(payload),
		"SPACEDOCK_PI_ENTITY_PATH="+entityPath,
		"SPACEDOCK_PI_STATE_DIR="+stateDir,
		"SPACEDOCK_PI_WORKFLOW_DIR="+root,
		"SPACEDOCK_PI_WORKTREE_PATH="+worktreePath,
	)
	workerOut, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fixture Pi worker process failed: %v\n%s", err, workerOut)
	}

	entityBody, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"## Stage Report: implementation fixback smoke",
		"- process: fixture Pi worker exited 0",
		"- dispatch_file: " + out.DispatchFile,
		"- wrapper_context: fresh",
		"- checklist_seen: run the fixture Pi worker; append durable stage report",
	} {
		if !strings.Contains(string(entityBody), want) {
			t.Fatalf("entity report missing durable evidence %q:\n%s", want, entityBody)
		}
	}

	logCmd := exec.Command("git", "-C", stateDir, "log", "--oneline", "--", entityPath)
	logOut, err := logCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("state git log failed: %v\n%s", err, logOut)
	}
	if !strings.Contains(string(logOut), "fixture Pi stage report") {
		t.Fatalf("state git log missing fixture stage-report commit:\n%s", logOut)
	}
	statusCmd := exec.Command("git", "-C", stateDir, "status", "--short", "--", entityPath)
	statusOut, err := statusCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("state git status failed: %v\n%s", err, statusOut)
	}
	if strings.TrimSpace(string(statusOut)) != "" {
		t.Fatalf("state entity should be clean after path-scoped commit:\n%s", statusOut)
	}
}

func TestPiStageDispatchSmokeFixtureWorker(t *testing.T) {
	if os.Getenv("SPACEDOCK_PI_FIXTURE_WORKER") != "1" {
		return
	}
	var wrapped piruntime.SubagentDispatch
	if err := json.Unmarshal([]byte(os.Getenv("SPACEDOCK_PI_WRAPPER_JSON")), &wrapped); err != nil {
		t.Fatalf("wrapper JSON: %v", err)
	}
	if wrapped.Context != "fresh" {
		t.Fatalf("wrapper context = %q, want fresh", wrapped.Context)
	}
	if strings.Contains(os.Getenv("SPACEDOCK_PI_WRAPPER_JSON"), "acceptance") {
		t.Fatalf("wrapper unexpectedly serialized same-agent acceptance: %s", os.Getenv("SPACEDOCK_PI_WRAPPER_JSON"))
	}

	dispatchPath, err := piDispatchFileFromTask(wrapped.Task)
	if err != nil {
		t.Fatal(err)
	}
	dispatchBody, err := os.ReadFile(dispatchPath)
	if err != nil {
		t.Fatalf("read dispatch artifact: %v", err)
	}
	entityPath := os.Getenv("SPACEDOCK_PI_ENTITY_PATH")
	workflowDir := os.Getenv("SPACEDOCK_PI_WORKFLOW_DIR")
	worktreePath := os.Getenv("SPACEDOCK_PI_WORKTREE_PATH")
	dispatchText := string(dispatchBody)
	for _, want := range []string{
		"You are working on: Smoke Stage",
		"Stage: implementation",
		"Read the entity file at " + entityPath,
		"Your working directory for CODE is " + worktreePath,
		"- run the fixture Pi worker",
		"- append durable stage report",
	} {
		if !strings.Contains(dispatchText, want) {
			t.Fatalf("dispatch artifact missing %q:\n%s", want, dispatchBody)
		}
	}
	assertRenderedStageDefCommand(t, dispatchText, workflowDir, "implementation")

	report := fmt.Sprintf(`
## Stage Report: implementation fixback smoke

- process: fixture Pi worker exited 0
- dispatch_file: %s
- wrapper_context: %s
- wrapper_phase: %s
- wrapper_label: %s
- checklist_seen: run the fixture Pi worker; append durable stage report
`, dispatchPath, wrapped.Context, wrapped.Phase, wrapped.Label)
	f, err := os.OpenFile(entityPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open entity report: %v", err)
	}
	if _, err := f.WriteString(report); err != nil {
		_ = f.Close()
		t.Fatalf("append entity report: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close entity report: %v", err)
	}

	stateDir := os.Getenv("SPACEDOCK_PI_STATE_DIR")
	for _, args := range [][]string{
		{"-c", "user.email=t@t", "-c", "user.name=t", "add", entityPath},
		{"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", "fixture Pi stage report", "--", entityPath},
	} {
		cmd := exec.Command("git", append([]string{"-C", stateDir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

func assertRenderedStageDefCommand(t *testing.T, body, workflowDir, stage string) {
	t.Helper()
	want := shlexQuote(testWorkflowLauncher) + " dispatch show-stage-def --workflow-dir " + shlexQuote(workflowDir) + " --stage " + stage
	if !strings.Contains(body, want) {
		t.Fatalf("dispatch artifact missing exact stage-definition command %q:\n%s", want, body)
	}
}

func piDispatchFileFromTask(task string) (string, error) {
	const prefix = "Read "
	if !strings.HasPrefix(task, prefix) {
		return "", fmt.Errorf("Pi task does not start with %q: %q", prefix, task)
	}
	rest := strings.TrimPrefix(task, prefix)
	marker := " and treat its content as your assignment"
	idx := strings.Index(rest, marker)
	if idx < 0 {
		return "", fmt.Errorf("Pi task does not forward read-dispatch-file assignment: %q", task)
	}
	return rest[:idx], nil
}

func TestBuildPiHostPreservesSplitRootEntityPath(t *testing.T) {
	root := t.TempDir()
	stateDir := filepath.Join(root, "state-checkout")
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(true))
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(stateDir, "thing", "index.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- preserve split-root entity path"},
		"bare_mode":      false,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))
	worktreeEntityPath := filepath.Join(root, worktreeRel, "state-checkout", "thing", "index.md")
	if strings.Contains(body, worktreeEntityPath) {
		t.Fatalf("pi split-root dispatch must not rewrite entity path into code worktree: %s\n%s", worktreeEntityPath, body)
	}
	if !strings.Contains(body, "Read the entity file at "+entityPath) {
		t.Fatalf("pi split-root dispatch body does not point at state-checkout entity %s:\n%s", entityPath, body)
	}
	if !strings.Contains(body, "This workflow is split-root") {
		t.Fatalf("pi split-root dispatch body missing split-root state guidance:\n%s", body)
	}
}

// TestBuildPiHostEmitsSpawnAgentAndSkill is AC-2's default-path assertion: a
// pi dispatch with no stage agent: override carries the pi-subagents generic
// write-capable agent ("worker") and the basename skill ("ensign") as spawn
// fields, while subagent_type keeps the host-neutral role identity.
func TestBuildPiHostEmitsSpawnAgentAndSkill(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	var out struct {
		SubagentType string  `json:"subagent_type"`
		Agent        *string `json:"agent"`
		Skill        *string `json:"skill"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	if out.SubagentType != "spacedock:ensign" {
		t.Errorf("subagent_type = %q, want host-neutral spacedock:ensign", out.SubagentType)
	}
	if out.Agent == nil || *out.Agent != "worker" {
		t.Errorf("agent = %v, want %q (pi-subagents generic write-capable agent)", out.Agent, "worker")
	}
	if out.Skill == nil || *out.Skill != "ensign" {
		t.Errorf("skill = %v, want %q (basename, the only form pi's resolver binds)", out.Skill, "ensign")
	}
	// AC-4's banned-string half for the build pi branch: the namespaced agent
	// string must stay out of pi spawn surfaces (it names no pi-subagents agent).
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))
	if strings.Contains(body, "agent: spacedock:ensign") {
		t.Errorf("pi dispatch body must not name a pi-subagents agent spacedock:ensign:\n%s", body)
	}
}

// TestBuildPiHostAgentOverrideOmitsSkill is AC-2's override assertion: a stage
// agent: declaration takes over the spawn agent (an override agent owns its
// own contract), so the skill field is omitted and subagent_type follows the
// override as before.
func TestBuildPiHostAgentOverrideOmitsSkill(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWithAgentOverride())
	worktreeRel := ".worktrees/spacedock-satellite-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(native.stdout), &fields); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	var agent string
	if err := json.Unmarshal(fields["agent"], &agent); err != nil || agent != "custom:agent" {
		t.Errorf("agent = %q (err %v), want the stage's agent: override", agent, err)
	}
	if _, present := fields["skill"]; present {
		t.Errorf("skill key must be omitted when the stage overrides agent (the override owns its contract):\n%s", native.stdout)
	}
	var subagentType string
	if err := json.Unmarshal(fields["subagent_type"], &subagentType); err != nil || subagentType != "custom:agent" {
		t.Errorf("subagent_type = %q (err %v), want the override as before", subagentType, err)
	}
}

// readmeWithAgentOverride declares an implementation stage routed to a custom
// agent.
func readmeWithAgentOverride() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: initialization\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"      worktree: true\n" +
		"      agent: custom:agent\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Agent Override Workflow\n" +
		"\n### initialization\n\nseed.\n\n- **Outputs:** x.\n\n" +
		"### implementation\n\nwork.\n\n- **Outputs:** y.\n\n" +
		"### done\n\nterm.\n"
}

// TestBuildClaudeHostGoldenByteIdentical is AC-2's no-churn guard: the exact
// claude-build fixture of the existing parity case single+flat+nonworktree+bare
// is re-run and byte-compared against its checked-in golden, so the added pi
// spawn fields cannot have disturbed claude output bytes.
func TestBuildClaudeHostGoldenByteIdentical(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "backlog",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      true,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	for _, banned := range []string{`"agent":`, `"skill":`} {
		if strings.Contains(native.stdout, banned) {
			t.Fatalf("claude build stdout must not carry pi spawn field %s:\n%s", banned, native.stdout)
		}
	}
	nativeBody := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))
	env := goldenEnvelope{res: normRun(native, root, home), body: normPaths(nativeBody, root, home)}
	assertGolden(t, "build-crossproduct-single+flat+nonworktree+bare", env)
}
