// ABOUTME: Paseo dispatch-build host shape — host-neutral schema-v2 envelope,
// ABOUTME: bare-mode refusal, PASEO_AGENT_ID file disambiguation, and banned-token checks.
package dispatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var paseoBannedTokens = []string{
	"Skill(skill=", "SendMessage", "$spacedock:ensign", "/skill:ensign",
	"pi-subagents", "TeamCreate", "TeamDelete", "Agent(", "FO mailbox",
}

func TestBuildPaseoHostPromptShapeAndSplitRoot(t *testing.T) {
	t.Setenv("PASEO_AGENT_ID", "")
	root := t.TempDir()
	stateDir := filepath.Join(root, "state-checkout")
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(true))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	for _, dir := range []string{stateDir, filepath.Join(root, worktreeRel)} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	entityPath := filepath.Join(stateDir, "thing", "index.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)
	gitInit(t, stateDir)

	res := runNative("- first item\n- second item",
		"build", "--workflow-dir", root, "--entity-path", entityPath,
		"--stage", "implementation", "--checklist-file", "-", "--host", "paseo",
	)
	if res.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", res.exit, res.stderr)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(res.stdout), &raw); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, res.stdout)
	}
	for _, omitted := range []string{"run_in_background", "agent", "skill"} {
		if _, present := raw[omitted]; present {
			t.Errorf("paseo envelope must omit %q:\n%s", omitted, res.stdout)
		}
	}
	if string(raw["model"]) != "null" {
		t.Errorf("model = %s, want null", raw["model"])
	}
	var out struct {
		SchemaVersion int      `json:"schema_version"`
		SubagentType  string   `json:"subagent_type"`
		Description   string   `json:"description"`
		FetchCommands []string `json:"fetch_commands"`
		DispatchFile  string   `json:"dispatch_file_path"`
		Prompt        string   `json:"prompt"`
		Name          string   `json:"name"`
	}
	if err := json.Unmarshal([]byte(res.stdout), &out); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if out.SchemaVersion != 2 || out.SubagentType != "spacedock:ensign" || out.Description != "Thing: implementation" || out.Name != "thing-implementation" || len(out.FetchCommands) != 1 || filepath.Base(out.DispatchFile) != "thing-implementation.md" {
		t.Fatalf("unexpected envelope: %+v", out)
	}
	wantPrompt := fmt.Sprintf("Load the Spacedock ensign skill (with references/ensign-shared-core.md and references/paseo-ensign-runtime.md); then Read %s and treat its content as your assignment.", out.DispatchFile)
	if out.Prompt != wantPrompt {
		t.Fatalf("prompt = %q, want %q", out.Prompt, wantPrompt)
	}
	body := readDispatchBody(t, out.DispatchFile)
	for _, banned := range paseoBannedTokens {
		if strings.Contains(out.Prompt, banned) || strings.Contains(body, banned) {
			t.Errorf("paseo dispatch contains banned token %q", banned)
		}
	}
	for _, want := range []string{
		"load the Spacedock `ensign` skill (`SKILL.md`, `references/ensign-shared-core.md`, and `references/paseo-ensign-runtime.md`)",
		"This assignment is delivered through Paseo.",
		"Read the entity file at " + entityPath,
		"This workflow is split-root",
		"### Completion Signal",
		"Done: Thing completed implementation. Report written to " + entityPath + ".",
		"verifies completion from committed Git state and the entity stage report after your turn becomes idle",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("paseo dispatch body missing %q:\n%s", want, body)
		}
	}
}

func TestBuildPaseoHostWithStampAndBareRejection(t *testing.T) {
	t.Setenv("PASEO_AGENT_ID", "")
	root := t.TempDir()
	wfDir := filepath.Join(root, "docs", "wf")
	stateDir := filepath.Join(wfDir, ".spacedock-state")
	writeFile(t, filepath.Join(wfDir, "README.md"), strings.ReplaceAll(readmeWorktree(true), "state: state-checkout", "state: .spacedock-state"))
	entityPath := filepath.Join(stateDir, "stamped", "index.md")
	writeFile(t, entityPath, entityFM("Stamped Paseo", "implementation", ""))
	gitInit(t, root)
	gitInit(t, stateDir)
	runGitFatal(t, stateDir, "branch", "-M", "spacedock-state/wf")

	res := runNative("- stamped item", "build", "--host", "paseo", "--stamp",
		"--workflow-dir", wfDir, "--entity-path", entityPath, "--stage", "implementation", "--checklist-file", "-")
	if res.exit != 0 {
		t.Fatalf("build --host paseo --stamp exit=%d stderr=%q", res.exit, res.stderr)
	}
	if info, err := os.Stat(filepath.Join(root, ".worktrees", "spacedock-ensign-stamped")); err != nil || !info.IsDir() {
		t.Fatalf("expected stamped worktree: err=%v", err)
	}

	bareRoot := t.TempDir()
	writeFile(t, filepath.Join(bareRoot, "README.md"), readmeNonWorktreeStages())
	bareEntity := filepath.Join(bareRoot, "paseobare.md")
	writeFile(t, bareEntity, entityFM("Paseo Bare", "implementation", ""))
	gitInit(t, bareRoot)
	artifactPath := filepath.Join(dispatchFileDir, "paseobare-implementation.md")
	_ = os.Remove(artifactPath)
	t.Cleanup(func() { _ = os.Remove(artifactPath) })

	bare := runNative("- item", "build", "--workflow-dir", bareRoot, "--entity-path", bareEntity,
		"--stage", "implementation", "--checklist-file", "-", "--host", "paseo", "--bare-mode")
	const wantErr = "bare_mode is unsupported on host paseo; Paseo worker.spawn requires a named worker label"
	if bare.exit != 2 || !strings.Contains(bare.stderr, wantErr) {
		t.Fatalf("bare paseo exit=%d stderr=%q, want exit 2 + %q", bare.exit, bare.stderr, wantErr)
	}
	if _, err := os.Stat(artifactPath); !os.IsNotExist(err) {
		t.Fatalf("bare paseo must not create dispatch artifact %s; stat err=%v", artifactPath, err)
	}
}

func TestBuildPaseoHostModelDisambiguationAndResolve(t *testing.T) {
	t.Setenv("PASEO_AGENT_ID", "")
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeModelsFable)
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "stagemodel", ""))
	gitInit(t, root)

	res := runNative("- a", "build", "--workflow-dir", root, "--entity-path", entityPath,
		"--stage", "stagemodel", "--checklist-file", "-", "--host", "paseo")
	wantNote := "[build] declared model 'fable' ignored on host paseo: outside paseo's dispatch-settable model space; emitting model=null"
	if res.exit != 0 || !strings.Contains(res.stdout, `"model": null`) || !strings.Contains(res.stderr, wantNote) || strings.Contains(res.stderr, "effective_model=") {
		t.Fatalf("model ignore failed: exit=%d stdout=%s stderr=%s", res.exit, res.stdout, res.stderr)
	}

	hRoot, hEntity := buildHostFixture(t)
	args := []string{"build", "--workflow-dir", hRoot, "--entity-path", hEntity, "--stage", "backlog", "--checklist-file", "-", "--host", "paseo"}
	t.Setenv("PASEO_AGENT_ID", "Agent_Alpha/01")
	a := runNative("- a", args...)
	t.Setenv("PASEO_AGENT_ID", "Agent_Beta/02")
	b := runNative("- a", args...)
	if a.exit != 0 || b.exit != 0 {
		t.Fatalf("disambiguated builds failed: a=%d b=%d", a.exit, b.exit)
	}
	pathA, pathB := dispatchFilePathFromStdout(t, a.stdout), dispatchFilePathFromStdout(t, b.stdout)
	if pathA == pathB || filepath.Base(pathA) != "agent-alpha-01-thing-backlog.md" || filepath.Base(pathB) != "agent-beta-02-thing-backlog.md" {
		t.Fatalf("unexpected disambiguated paths: a=%q b=%q", pathA, pathB)
	}

	if host, err := resolveBuildHost("", hostEnv(map[string]string{"PASEO_AGENT_ID": "agt-1", "CODEX_THREAD_ID": "t1"})); err != nil || host != "codex" {
		t.Fatalf("PASEO_AGENT_ID caused ambiguity: (%q, %v)", host, err)
	}
	if host, err := resolveBuildHost("paseo", hostEnv(map[string]string{"CODEX_THREAD_ID": "t1", "CLAUDECODE": "1"})); err != nil || host != "paseo" {
		t.Fatalf("explicit --host paseo failed: (%q, %v)", host, err)
	}
}
