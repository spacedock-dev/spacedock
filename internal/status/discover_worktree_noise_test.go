// ABOUTME: discoverWorkflows prunes agent/linked-worktree noise via .gitignore
// ABOUTME: dir patterns AND nested-checkout (.git) skip — not a hardcoded list.
package status

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDiscoverWorkflowsPrunesWorktreeNoise: agent and linked worktrees are full
// repo checkouts that each carry a copy of the repo's commissioned workflow
// (e.g. `docs/dev`). They live under `.claude/worktrees/<agent>/…` (an untracked
// agent dir the repo `.gitignore` lists as `.claude/`) and `.worktrees/<name>/…`
// (each a linked worktree rooted by its own `.git` gitlink). Discovery must
// resolve the ONE real workflow by two composable mechanisms, NOT a hardcoded
// path skip: (a) honoring the `.gitignore` directory patterns so the `.claude`
// subtree drops out, and (b) not descending into a nested git checkout.
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

	// The repo .gitignore lists .claude/ as an ignored directory — the mechanism
	// that drops the agent-worktree copies under it (no .git needed at that copy).
	if err := os.WriteFile(filepath.Join(repo, ".gitignore"), []byte(".claude/\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// (real) the single workflow, part of the main checkout (no own .git).
	write(filepath.Join("docs", "dev"))
	// (a) gitignored .claude subtree — pruned by the .gitignore dir pattern.
	write(filepath.Join(".claude", "worktrees", "agent-a", "docs", "dev"))
	write(filepath.Join(".claude", "worktrees", "agent-b", "docs", "dev"))
	// (b) a linked worktree rooted by a .git gitlink — pruned as a nested checkout.
	gitlink(filepath.Join(".worktrees", "wt-x"))
	write(filepath.Join(".worktrees", "wt-x", "docs", "dev"))

	got := discoverWorkflows(repo)
	want := realpathOf(filepath.Join(repo, "docs", "dev"))
	if len(got) != 1 || got[0] != want {
		t.Fatalf("discoverWorkflows = %v, want exactly [%s] (worktree noise must be pruned)", got, want)
	}
}

// TestDiscoverWorkflowsHonorsGitignoreDirPattern isolates mechanism (a): a noise
// workflow copy under a dir the .gitignore lists is excluded purely by the
// gitignore pattern — no `.git` gitlink at the copy, no hardcoded basename.
func TestDiscoverWorkflowsHonorsGitignoreDirPattern(t *testing.T) {
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
	// An arbitrarily-named ignored dir (not in discoverIgnoreDirs) proves the
	// gitignore pattern — not a hardcoded list — does the exclusion.
	if err := os.WriteFile(filepath.Join(repo, ".gitignore"), []byte("agent-cache/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	write(filepath.Join("docs", "dev"))
	write(filepath.Join("agent-cache", "copy", "docs", "dev"))

	got := discoverWorkflows(repo)
	want := realpathOf(filepath.Join(repo, "docs", "dev"))
	if len(got) != 1 || got[0] != want {
		t.Fatalf("discoverWorkflows = %v, want exactly [%s] (gitignore dir pattern must prune the copy)", got, want)
	}
}

// TestDiscoverWorkflowsSkipsNestedCheckout isolates mechanism (b): a nested git
// checkout's workflow copy is pruned by the `.git`-presence skip ALONE, with no
// help from a basename or gitignore rule. The copy lives under `submods/checkout-z`
// — a dir name that is NEITHER in discoverIgnoreDirs NOR any `.gitignore` — so the
// ONLY thing that can drop it is hasGitEntry seeing the nested `.git` gitlink. This
// is the unmasked check: neuter hasGitEntry → `return false` and this test reds
// (the copy re-appears as a second workflow), which the existing `.worktrees`-sited
// fixtures do NOT (the `.worktrees` basename in discoverIgnoreDirs masks them).
func TestDiscoverWorkflowsSkipsNestedCheckout(t *testing.T) {
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

	// The real workflow, part of the main checkout (no own .git).
	write(filepath.Join("docs", "dev"))
	// A nested checkout under a NON-pruned, NON-gitignored dir name, rooted by a
	// `.git` gitlink — pruned only by hasGitEntry.
	nested := filepath.Join(repo, "submods", "checkout-z")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, ".git"), []byte("gitdir: /elsewhere\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	write(filepath.Join("submods", "checkout-z", "docs", "dev"))

	got := discoverWorkflows(repo)
	want := realpathOf(filepath.Join(repo, "docs", "dev"))
	if len(got) != 1 || got[0] != want {
		t.Fatalf("discoverWorkflows = %v, want exactly [%s] (nested-checkout .git skip must prune the copy)", got, want)
	}
}

// TestDiscoverWorkflowsPrunesTestdata isolates the `testdata` basename prune: a
// commissioned-shape README that exists only as a test fixture lives under a
// package's `testdata/` subtree (the Go-idiomatic, package-adjacent home, e.g.
// skills/integration/testdata/refit-content-propagation/site-workflow). Discovery
// must NOT count it as a real workflow while still finding the one real workflow
// at docs/dev. The `testdata` basename is the ONLY thing that drops the fixture
// row here — its parent dirs (skills, integration) are neither in
// discoverIgnoreDirs nor any .gitignore — so removing `testdata` from
// discoverIgnoreDirs re-surfaces the fixture as a second workflow and reds this.
func TestDiscoverWorkflowsPrunesTestdata(t *testing.T) {
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

	// The one real workflow, part of the main checkout.
	write(filepath.Join("docs", "dev"))
	// A commissioned README carried as a fixture under a package's testdata/ —
	// the shape of the relocated skills/integration/testdata fixtures. Pruned by
	// the `testdata` basename, not by any parent dir name or gitignore rule.
	write(filepath.Join("skills", "integration", "testdata", "refit-content-propagation", "site-workflow"))

	got := discoverWorkflows(repo)
	want := realpathOf(filepath.Join(repo, "docs", "dev"))
	if len(got) != 1 || got[0] != want {
		t.Fatalf("discoverWorkflows = %v, want exactly [%s] (a commissioned README under testdata/ must be pruned)", got, want)
	}
}
