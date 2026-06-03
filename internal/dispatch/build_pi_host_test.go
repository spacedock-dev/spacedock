// ABOUTME: Pi dispatch-build host shape — the outer prompt and dispatch body avoid
// ABOUTME: Claude team tool signatures and target pi-subagents style completion.
package dispatch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
		"team_name":      "fixture-pi-team",
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
		"team_name":      "fixture-pi-team",
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
