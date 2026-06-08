// ABOUTME: AC-B oracles for the pre-launch info banner — version line + detected
// ABOUTME: workflow rel-path (commissioned README) vs none-detected (bare dir).
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
// `commissioned-by: spacedock@…` — the same predicate DiscoverWorkflowDir
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

// TestLaunchBannerNamesDetectedWorkflow (AC-B): launched from inside a
// commissioned workflow, the banner names the workflow's path RELATIVE TO THE
// GIT REPO ROOT (so a workflow at <repo>/docs/dev reads `Workflow: docs/dev`, a
// recognizable path that orients the operator to which workflow) — never the
// cwd-relative `.`/`..`. Outside any workflow it reads "none detected"; inside a
// workflow with no enclosing `.git`, it falls back to the workflow dir's name
// (never `.`/`..`). Every case carries the version line naming cli.Version. The
// fixture's repo-relative location is the independent expected value.
func TestLaunchBannerNamesDetectedWorkflow(t *testing.T) {
	t.Run("inside a commissioned workflow under a git repo", func(t *testing.T) {
		repo := gitRepoFixture(t)
		workflow := filepath.Join(repo, "docs", "dev")
		commissionWorkflowAt(t, workflow)
		sub := filepath.Join(workflow, "nested", "deep")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		launchBanner("claude", sub, &buf)

		out := buf.String()
		if !strings.Contains(out, "spacedock "+Version) {
			t.Fatalf("banner missing version line naming %q: %q", "spacedock "+Version, out)
		}
		// Repo-relative: the workflow sits at docs/dev under the repo root, so the
		// banner must name docs/dev regardless of how deep the launch dir is.
		if !strings.Contains(out, "Workflow: "+filepath.Join("docs", "dev")) {
			t.Fatalf("banner workflow line does not name the repo-relative path docs/dev: %q", out)
		}
		if strings.Contains(out, "Workflow: .") {
			t.Fatalf("banner workflow line is the cwd-relative `.`/`..` form, not the repo-relative path: %q", out)
		}
		if strings.Contains(out, "none detected") {
			t.Fatalf("banner reads none detected inside a commissioned workflow: %q", out)
		}
	})

	t.Run("commissioned workflow with no enclosing git repo falls back to the workflow name", func(t *testing.T) {
		// A workflow dir with no `.git` on the way up: FindGitRoot finds nothing, so
		// the banner falls back to the workflow dir's own name — never `.`/`..`.
		workflow := t.TempDir()
		commissionWorkflowAt(t, workflow)
		var buf bytes.Buffer
		launchBanner("codex", workflow, &buf)

		out := buf.String()
		if !strings.Contains(out, "Workflow: "+filepath.Base(workflow)) {
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
		bare := t.TempDir() // no commissioned README on the way up to the temp root
		var buf bytes.Buffer
		launchBanner("codex", bare, &buf)

		out := buf.String()
		if !strings.Contains(out, "spacedock "+Version) {
			t.Fatalf("banner missing version line naming %q: %q", "spacedock "+Version, out)
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
			if !strings.Contains(stderr.String(), "spacedock "+Version) {
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
// version line must NOT appear on stderr. claude's resume is a flag; codex's is
// the leading `resume` subcommand.
func TestLaunchBannerSuppressedOnResume(t *testing.T) {
	t.Run("claude --resume", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		if code := runClaude(context.Background(), []string{"--", "--resume"}, t.TempDir(), fake, lookFound, &stdout, &stderr); code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if strings.Contains(stderr.String(), "spacedock "+Version) {
			t.Fatalf("banner emitted on --resume (should be suppressed): %q", stderr.String())
		}
	})

	t.Run("codex resume subcommand", func(t *testing.T) {
		dir := safehouseFixtureDir(t)
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		if code := runCodex(context.Background(), []string{"--", "resume", "abc123"}, dir, fake, lookFound, &stdout, &stderr); code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if strings.Contains(stderr.String(), "spacedock "+Version) {
			t.Fatalf("banner emitted on codex resume (should be suppressed): %q", stderr.String())
		}
	})
}
