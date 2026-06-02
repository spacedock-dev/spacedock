// ABOUTME: Codex dispatch-build host shape — the outer prompt and dispatch-file
// ABOUTME: body use Codex mailbox semantics instead of Claude Skill()/SendMessage.
package dispatch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCodexHostPromptShape(t *testing.T) {
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
		"team_name":      "fixture-team",
		"bare_mode":      false,
		"host":           "codex",
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
	if strings.Contains(out.Prompt, "Skill(skill=") {
		t.Fatalf("codex prompt must not depend on Skill(...): %q", out.Prompt)
	}
	if !strings.Contains(out.Prompt, "Read ") || !strings.Contains(out.Prompt, "treat its content as your assignment") {
		t.Fatalf("codex prompt should be the read-dispatch-file form: %q", out.Prompt)
	}

	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))
	for _, banned := range []string{"Skill(skill=\"spacedock:ensign\")", "SendMessage(to=\"team-lead\""} {
		if strings.Contains(body, banned) {
			t.Fatalf("codex dispatch body must omit %q:\n%s", banned, body)
		}
	}
	for _, want := range []string{
		"Read this dispatch file directly",
		"Codex final-status notification",
		"FO mailbox",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("codex dispatch body missing %q:\n%s", want, body)
		}
	}
}
