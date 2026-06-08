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

// commissionedWorkflowDir returns a temp dir holding a README.md whose
// frontmatter declares `commissioned-by: spacedock@…` — the same predicate
// DiscoverWorkflowDir recognizes. The returned dir is the workflow root; tests
// launch the banner from a SUBDIRECTORY so the discovered workflow is a real
// ancestor and the rendered rel path is non-trivial. The expectation is sourced
// from this fixture, independent of frontdoor.go.
func commissionedWorkflowDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	readme := "---\ncommissioned-by: spacedock@1.0\nid-style: sequential\n---\n# WF\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TestLaunchBannerNamesDetectedWorkflow (AC-B): launched from inside a
// commissioned workflow, the banner names the workflow's path relative to the
// launch dir; launched outside any workflow, it reads "none detected". Both
// carry the version line naming cli.Version.
func TestLaunchBannerNamesDetectedWorkflow(t *testing.T) {
	t.Run("inside a commissioned workflow", func(t *testing.T) {
		root := commissionedWorkflowDir(t)
		sub := filepath.Join(root, "nested", "deep")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		launchBanner("claude", sub, &buf)

		out := buf.String()
		if !strings.Contains(out, "spacedock "+Version) {
			t.Fatalf("banner missing version line naming %q: %q", "spacedock "+Version, out)
		}
		wantRel, err := filepath.Rel(sub, root)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Workflow: "+wantRel) {
			t.Fatalf("banner workflow line does not name the detected rel path %q: %q", wantRel, out)
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
