// ABOUTME: Executes the break-glass worker's shipped commit command on a single-root
// ABOUTME: and a split-root workflow and checks the entity report lands committed.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestBreakGlassCommitLandsInEntityRepo runs each `### Commit` command from the
// recovery skill's worker prompts. On a split-root workflow the entity lives in
// a nested, ignored state checkout, so a commit run from {workflow_dir} fails.
func TestBreakGlassCommitLandsInEntityRepo(t *testing.T) {
	data, err := os.ReadFile("../fo-dispatch-recovery/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	commands := regexp.MustCompile("commit your work path-scoped: `([^`]+)`").FindAllStringSubmatch(string(data), -1)
	if len(commands) != 2 {
		t.Fatalf("want the commit command in both break-glass prompts, found %d", len(commands))
	}
	for _, layout := range []string{"single-root", "split-root"} {
		for i, command := range commands {
			root := t.TempDir()
			workflowDir := filepath.Join(root, "docs", "dev")
			entityDir := workflowDir
			if layout == "split-root" {
				entityDir = filepath.Join(workflowDir, ".spacedock-state")
			}
			if err := os.MkdirAll(entityDir, 0o755); err != nil {
				t.Fatal(err)
			}
			git(t, root, "init", "-q")
			if layout == "split-root" {
				if err := os.WriteFile(filepath.Join(workflowDir, ".gitignore"), []byte(".spacedock-state/\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				git(t, entityDir, "init", "-q")
			}
			entity := filepath.Join(entityDir, "widget-task.md")
			if err := os.WriteFile(entity, []byte("## Stage Report: implementation\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			script := strings.NewReplacer(`\"`, `"`, "{workflow_dir}", workflowDir, "{entity_file_path}", entity,
				"{entity title}", "Widget task", "{stage}", "implementation").Replace(command[1])
			cmd := exec.Command("sh", "-c", script)
			cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("%s prompt %d: commit command failed: %v\n%s\n%s", layout, i, err, script, out)
			}
			if got := git(t, entityDir, "log", "-1", "--format=%s", "--", entity); got != "Widget task: implementation report\n" {
				t.Fatalf("%s prompt %d: entity commit subject = %q", layout, i, got)
			}
		}
	}
}
