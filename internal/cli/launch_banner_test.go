// ABOUTME: AC-B oracles for the pre-launch info banner — version line + detected
// ABOUTME: workflow (repo-wide discovery, noise-pruned) vs none/multiple.
package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// commissionWorkflowAt writes a README.md whose frontmatter declares
// `commissioned-by: spacedock@…` — the same predicate the workflow discovery
// recognizes — at the given absolute dir, creating it first.
func commissionWorkflowAt(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := "---\ncommissioned-by: spacedock@1.0\nid-style: sequential\n---\n# WF\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
}

// gitRepoFixture returns a temp dir holding a `.git` directory so FindGitRoot
// resolves it as the enclosing repository root.
func gitRepoFixture(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	return repo
}

// TestLaunchBannerNamesDetectedWorkflow (AC-B): the banner finds the repo's REAL
// top-level workflow from ANY launch dir (the enclosing git repo is scanned
// downward, with linked/agent-worktree + VCS noise pruned), and names it relative
// to the repo root (e.g. `docs/dev`). From the repo root, from inside the
// workflow, and from a deep subdir it names the same workflow. A noise copy under
// `.claude/worktrees/...` (a real concern: agent worktrees are full repo checkouts
// each carrying a `docs/dev`) is NOT named and does NOT inflate the list. Outside
// any workflow it reads "none detected"; with more than one top-level workflow it
// lists the detected workflow paths space-separated. The fixture layout is the
// independent expected value.
func TestLaunchBannerNamesDetectedWorkflow(t *testing.T) {
	// repoWithRealWorkflowAndNoise builds a temp git repo with the single real
	// workflow at <repo>/docs/dev plus noise copies that each exercise an exclusion
	// MECHANISM (not a hardcoded skip): the `.claude/worktrees/agent-x/docs/dev`
	// copy is dropped because the repo `.gitignore` lists `.claude/`, and the
	// `.worktrees/wt-y/docs/dev` copy is dropped because its root is a nested git
	// checkout (a `.git` gitlink). Only docs/dev is a real top-level workflow.
	gitlinkAt := func(t *testing.T, dir string) {
		t.Helper()
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /elsewhere\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	repoWithRealWorkflowAndNoise := func(t *testing.T) (repo, realWorkflow string) {
		t.Helper()
		repo = gitRepoFixture(t)
		if err := os.WriteFile(filepath.Join(repo, ".gitignore"), []byte(".claude/\n.safehouse\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		realWorkflow = filepath.Join(repo, "docs", "dev")
		commissionWorkflowAt(t, realWorkflow)
		// (a) gitignored .claude subtree — no .git at the copy; the .gitignore prunes it.
		commissionWorkflowAt(t, filepath.Join(repo, ".claude", "worktrees", "agent-x", "docs", "dev"))
		// (b) a nested linked worktree — its .git gitlink stops the descent.
		gitlinkAt(t, filepath.Join(repo, ".worktrees", "wt-y"))
		commissionWorkflowAt(t, filepath.Join(repo, ".worktrees", "wt-y", "docs", "dev"))
		return repo, realWorkflow
	}

	t.Run("from the repo root names the real workflow, not the .claude/.worktrees noise", func(t *testing.T) {
		repo, _ := repoWithRealWorkflowAndNoise(t)
		var buf bytes.Buffer
		launchBanner("claude", repo, false, bannerEnv(nil), lookMissing, &buf)

		out := buf.String()
		if !strings.Contains(out, "spacedock "+displayVersion()) {
			t.Fatalf("banner missing version line naming %q: %q", "spacedock "+displayVersion(), out)
		}
		if !strings.Contains(out, "Workflow: "+filepath.Join("docs", "dev")+"\n") {
			t.Fatalf("banner from repo root does not name the real workflow docs/dev: %q", out)
		}
		if strings.Contains(out, "none detected") {
			t.Fatalf("banner reads none detected from the repo root of a repo with a real workflow: %q", out)
		}
		if strings.Contains(out, "workflows detected") {
			t.Fatalf("banner counted the .claude/.worktrees noise copies as extra workflows: %q", out)
		}
	})

	t.Run("from a deep subdir of the workflow names the same workflow", func(t *testing.T) {
		repo, realWorkflow := repoWithRealWorkflowAndNoise(t)
		sub := filepath.Join(realWorkflow, "nested", "deep")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		launchBanner("codex", sub, false, bannerEnv(nil), lookMissing, &buf)

		out := buf.String()
		if !strings.Contains(out, "Workflow: "+filepath.Join("docs", "dev")+"\n") {
			t.Fatalf("banner from a deep subdir of repo %s does not name docs/dev: %q", repo, out)
		}
	})

	t.Run("more than one real top-level workflow lists the detected paths", func(t *testing.T) {
		repo := gitRepoFixture(t)
		commissionWorkflowAt(t, filepath.Join(repo, "docs", "dev"))
		commissionWorkflowAt(t, filepath.Join(repo, "ops", "release"))
		var buf bytes.Buffer
		launchBanner("claude", repo, false, bannerEnv(nil), lookMissing, &buf)

		out := buf.String()
		wantLine := "Workflows: " + filepath.Join("docs", "dev") + " " + filepath.Join("ops", "release") + "\n"
		if !strings.Contains(out, wantLine) {
			t.Fatalf("banner does not list the two workflow paths %q: %q", wantLine, out)
		}
		if strings.Contains(out, "none detected") {
			t.Fatalf("banner reads none detected with two real workflows: %q", out)
		}
	})

	t.Run("commissioned workflow with no enclosing git repo falls back to the workflow name", func(t *testing.T) {
		// A workflow dir with no `.git` on the way up: the discovery falls back to
		// the bounded walk-up and renders the workflow dir's own name — never `.`/`..`.
		workflow := t.TempDir()
		commissionWorkflowAt(t, workflow)
		var buf bytes.Buffer
		launchBanner("codex", workflow, false, bannerEnv(nil), lookMissing, &buf)

		out := buf.String()
		if !strings.Contains(out, "Workflow: "+filepath.Base(workflow)+"\n") {
			t.Fatalf("banner fallback does not name the workflow dir base %q: %q", filepath.Base(workflow), out)
		}
		if strings.Contains(out, "Workflow: .\n") || strings.Contains(out, "Workflow: ..") {
			t.Fatalf("banner fallback is the cwd-relative `.`/`..` form: %q", out)
		}
		if strings.Contains(out, "none detected") {
			t.Fatalf("banner reads none detected inside a commissioned workflow: %q", out)
		}
	})

	t.Run("outside any workflow", func(t *testing.T) {
		bare := t.TempDir() // no commissioned README in or above this temp dir
		var buf bytes.Buffer
		launchBanner("codex", bare, false, bannerEnv(nil), lookMissing, &buf)

		out := buf.String()
		if !strings.Contains(out, "spacedock "+displayVersion()) {
			t.Fatalf("banner missing version line naming %q: %q", "spacedock "+displayVersion(), out)
		}
		if !strings.Contains(out, "none detected") {
			t.Fatalf("banner does not read none detected outside a workflow: %q", out)
		}
	})
}

// TestLaunchBannerReachesStderrBeforeLaunch (AC-B launcher half): the banner's
// version line reaches stderr on a real launch (gate compatible, no resume), and
// the launch seam is reached. The fakeHost records Launch, so a banner on stderr
// PLUS a recorded launch proves the banner was emitted before the host launched.
func TestLaunchBannerReachesStderrBeforeLaunch(t *testing.T) {
	cases := []struct {
		name string
		run  func(dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", func(dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), nil, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", func(dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), nil, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stderr bytes.Buffer

			code := tc.run(t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if !strings.Contains(stderr.String(), "spacedock "+displayVersion()) {
				t.Fatalf("banner version line did not reach stderr: %q", stderr.String())
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not reached after banner")
			}
		})
	}
}

// TestLaunchBannerSuppressedOnResume (AC-B polish): a resume continues an
// existing session, not a fresh launch, so the banner is suppressed — its
// version line must NOT appear on stderr. Claude's resume is a flag; Codex uses
// an exact `resume` token in its post-fence argv.
func TestLaunchBannerSuppressedOnResume(t *testing.T) {
	t.Run("claude --resume", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		if code := runClaude(context.Background(), []string{"--", "--resume"}, t.TempDir(), fake, lookFound, &stdout, &stderr); code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if strings.Contains(stderr.String(), "spacedock "+displayVersion()) {
			t.Fatalf("banner emitted on --resume (should be suppressed): %q", stderr.String())
		}
	})

	t.Run("codex exact post-fence resume", func(t *testing.T) {
		dir := safehouseFixtureDir(t)
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		if code := runCodex(context.Background(), []string{"--", "resume", "abc123"}, dir, fake, lookFound, &stdout, &stderr); code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if strings.Contains(stderr.String(), "spacedock "+displayVersion()) {
			t.Fatalf("banner emitted on codex exact resume: %q", stderr.String())
		}
	})
}
