// ABOUTME: Executes the break-glass worker's shipped commit command on a single-root
// ABOUTME: and a split-root workflow and checks the entity report lands committed.
package contractlint

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestBreakGlassCommitLandsInEntityRepo runs each break-glass `### Commit` command. A split-root
// entity lives in a nested state checkout the outer repo ignores, so {workflow_dir} cannot commit it.
func TestBreakGlassCommitLandsInEntityRepo(t *testing.T) {
	skill := readRepoFile(t, "skills/fo-dispatch-recovery/SKILL.md")
	commands := regexp.MustCompile("commit your work path-scoped: `([^`]+)`").FindAllStringSubmatch(skill, -1)
	if len(commands) != 2 {
		t.Fatalf("want the commit command in both break-glass prompts, found %d", len(commands))
	}
	for _, stateDir := range []string{"", ".spacedock-state"} { // single-root, then split-root
		for i, command := range commands {
			workflowDir := t.TempDir()
			entityDir := filepath.Join(workflowDir, stateDir)
			entity := filepath.Join(entityDir, "widget-task.md")
			// The outer repo ignores the state checkout; on single-root the second init is a no-op.
			setup := "git init -q " + workflowDir + " && mkdir -p " + entityDir + " && git init -q " + entityDir +
				" && echo .spacedock-state/ > " + filepath.Join(workflowDir, ".gitignore") + " && echo report > " + entity + " && "
			shipped := strings.NewReplacer(`\"`, `"`, "{workflow_dir}", workflowDir, "{entity_file_path}", entity,
				"{entity title}", "Widget task", "{stage}", "implementation").Replace(command[1])
			cmd := exec.Command("sh", "-c", setup+shipped+" && git -C "+entityDir+" log -1 --format=%s -- "+entity)
			cmd.Dir = t.TempDir() // never inside this repo, whatever the command does
			cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
			if out, err := cmd.CombinedOutput(); err != nil || !strings.HasSuffix(string(out), "\nWidget task: implementation report\n") {
				t.Fatalf("state dir %q, prompt %d: %v\n%s\n%s", stateDir, i, err, shipped, out)
			}
		}
	}
}
