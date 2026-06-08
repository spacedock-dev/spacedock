// ABOUTME: discoverWorkflows prunes linked/agent-worktree checkouts so a repo's
// ABOUTME: real workflow is not multiplied by `.claude/worktrees` / `.worktrees` copies.
package status

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDiscoverWorkflowsPrunesWorktreeNoise: agent and linked worktrees are full
// repo checkouts (each rooted by a `.git` gitlink) that carry a copy of the
// repo's commissioned workflow (e.g. `docs/dev`). They live under
// `.claude/worktrees/<agent>/…` and `.worktrees/<name>/…`. Discovery must prune a
// nested checkout wholesale so a repo with ONE real workflow resolves to exactly
// that one — not dozens of duplicate copies. The prune is host-neutral: it skips
// any descended dir carrying its own `.git`, so it catches the `.claude` agent
// worktrees without naming a host-specific path.
func TestDiscoverWorkflowsPrunesWorktreeNoise(t *testing.T) {
	repo := t.TempDir()
	commissioned := "---\ncommissioned-by: spacedock@1.0\nid-style: sequential\n---\n# WF\n"
	write := func(rel string) {
		p := filepath.Join(repo, rel, "README.md")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(commissioned), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// gitlink plants a `.git` regular file at a worktree root, matching how git
	// records a linked worktree (a `gitdir: …` pointer file, not a directory).
	gitlink := func(rel string) {
		p := filepath.Join(repo, rel)
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(p, ".git"), []byte("gitdir: /elsewhere\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// The single real workflow (part of the main checkout, no own .git), plus noise
	// copies under agent/linked worktree roots that each carry a .git gitlink.
	write(filepath.Join("docs", "dev"))
	gitlink(filepath.Join(".claude", "worktrees", "agent-a"))
	write(filepath.Join(".claude", "worktrees", "agent-a", "docs", "dev"))
	gitlink(filepath.Join(".claude", "worktrees", "agent-b"))
	write(filepath.Join(".claude", "worktrees", "agent-b", "docs", "dev"))
	gitlink(filepath.Join(".worktrees", "wt-x"))
	write(filepath.Join(".worktrees", "wt-x", "docs", "dev"))

	got := discoverWorkflows(repo)
	want := realpathOf(filepath.Join(repo, "docs", "dev"))
	if len(got) != 1 || got[0] != want {
		t.Fatalf("discoverWorkflows = %v, want exactly [%s] (worktree noise must be pruned)", got, want)
	}
}
