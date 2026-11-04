// ABOUTME: Path-scoped state commits: entity B's commit must not sweep a concurrent
// ABOUTME: writer's staged entity A, and the break-glass commit lands in the entity's repo.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// git runs a git subcommand in dir and fails the test on error.
func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-C", dir,
		"-c", "user.email=t@t", "-c", "user.name=t"}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// TestPathScopedCommitDoesNotSweepSibling reproduces the cross-attribution
// hazard the split-root contract forbids: in one shared, non-branched state
// checkout, writer A stages its entity file while writer B commits its own. The
// prescribed path-scoped commit (`git commit -- {B}`) must commit ONLY B and
// leave A's staged change intact — a bare `git commit` would sweep A in.
func TestPathScopedCommitDoesNotSweepSibling(t *testing.T) {
	checkout := t.TempDir()
	git(t, checkout, "init", "-q")

	entityA := filepath.Join(checkout, "entity-a", "index.md")
	entityB := filepath.Join(checkout, "entity-b", "index.md")
	for _, p := range []string{entityA, entityB} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("seed\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, checkout, "add", "-A")
	git(t, checkout, "commit", "-q", "-m", "seed both entities")

	// Writer A stages a change to entity A (mid-flight, not yet committed).
	if err := os.WriteFile(entityA, []byte("seed\nA's in-flight change\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, checkout, "add", entityA)

	// Writer B makes and commits ITS OWN change with the prescribed path-scoped
	// commit, naming only entity B's path.
	if err := os.WriteFile(entityB, []byte("seed\nB's change\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, checkout, "add", entityB)
	git(t, checkout, "commit", "-m", "state: entity-b update", "--", entityB)

	// B's commit must contain ONLY entity B, never entity A.
	committed := git(t, checkout, "show", "--name-only", "--format=", "HEAD")
	if strings.Contains(committed, "entity-a/index.md") {
		t.Fatalf("path-scoped commit for B swept in entity A:\n%s", committed)
	}
	if !strings.Contains(committed, "entity-b/index.md") {
		t.Fatalf("path-scoped commit for B did not include entity B:\n%s", committed)
	}

	// A's change must still be staged afterward — uncommitted, not clobbered.
	staged := git(t, checkout, "diff", "--cached", "--name-only")
	if !strings.Contains(staged, "entity-a/index.md") {
		t.Fatalf("entity A's staged change was swept away by B's commit; staged now:\n%q", staged)
	}
}

// TestBreakGlassCommitLandsInEntityRepo runs each break-glass `### Commit` command. A split-root
// entity lives in a nested state checkout the outer repo ignores, so {workflow_dir} cannot commit it.
func TestBreakGlassCommitLandsInEntityRepo(t *testing.T) {
	data, err := os.ReadFile("../fo-dispatch-recovery/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	commands := regexp.MustCompile("commit your work path-scoped: `([^`]+)`").FindAllStringSubmatch(string(data), -1)
	if len(commands) != 2 {
		t.Fatalf("want the commit command in both break-glass prompts, found %d", len(commands))
	}
	for _, stateDir := range []string{"", ".spacedock-state"} { // single-root, then split-root
		for i, command := range commands {
			workflowDir := t.TempDir()
			entityDir := filepath.Join(workflowDir, stateDir)
			entity := filepath.Join(entityDir, "widget-task.md")
			git(t, workflowDir, "init", "-q")
			if err := os.MkdirAll(entityDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if stateDir != "" {
				git(t, entityDir, "init", "-q")
				err = os.WriteFile(filepath.Join(workflowDir, ".gitignore"), []byte(stateDir+"/\n"), 0o644)
			}
			if err != nil || os.WriteFile(entity, []byte("report\n"), 0o644) != nil {
				t.Fatal("write fixture", err)
			}
			script := strings.NewReplacer(`\"`, `"`, "{workflow_dir}", workflowDir, "{entity_file_path}", entity,
				"{entity title}", "Widget task", "{stage}", "implementation").Replace(command[1])
			cmd := exec.Command("sh", "-c", script)
			cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("state dir %q, prompt %d: %v\n%s\n%s", stateDir, i, err, script, out)
			}
			if got := git(t, entityDir, "log", "-1", "--format=%s", "--", entity); got != "Widget task: implementation report\n" {
				t.Fatalf("state dir %q, prompt %d: entity commit subject = %q", stateDir, i, got)
			}
		}
	}
}
