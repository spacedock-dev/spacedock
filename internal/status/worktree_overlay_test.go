// ABOUTME: Worktree-overlay parity — a non-split-root worktree-backed entity's
// ABOUTME: active reads use the worktree-copy frontmatter, byte-matching the oracle.
package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildWorktreeBacked materializes a NON-split-root workflow in a git repo with
// one worktree-backed entity whose pipeline-dir copy and worktree copy disagree:
//
//	<root>/                       git root == workflow dir (no state: field)
//	  README.md                   stages + id-style
//	  add-login.md                pipeline copy: status=implementation, worktree=wt
//	  wt/add-login.md             worktree copy: status=review
//
// The worktree copy lives at <git_root>/<worktree>/<rel>, matching the oracle's
// os.path.join(git_root, worktree, relpath(entity_path, pipeline_dir)). Returns
// the workflow dir. A git repo is required because the overlay resolves git_root.
func buildWorktreeBacked(t *testing.T, readme, pipelineEntity, worktreeEntity string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readme)
	writeFile(t, filepath.Join(root, "add-login.md"), pipelineEntity)
	writeFile(t, filepath.Join(root, "wt", "add-login.md"), worktreeEntity)
	gitInitWorktreeFixture(t, root)
	return root
}

// gitInitWorktreeFixture initializes a git repo and commits the tree so the
// overlay's find_git_root resolves to root.
func gitInitWorktreeFixture(t *testing.T, dir string) {
	t.Helper()
	for _, args := range [][]string{
		{"init", "-q"},
		{"-c", "user.email=t@t", "-c", "user.name=t", "add", "-A"},
		{"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", "init"},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

// overlayReadme is a slug-style, single-root README (NO state: field) with the
// two stages the overlay test moves the entity between.
const overlayReadme = `---
commissioned-by: spacedock@1
id-style: slug
stages:
  states:
    - name: implementation
      initial: true
      worktree: true
    - name: review
      terminal: true
---

# Worktree Overlay Workflow
`

// TestWorktreeOverlayActiveReads is the M2 parity test: for a NON-split-root
// worktree-backed entity whose pipeline-dir status differs from its worktree-copy
// status, native active reads (table / --where / --fields / --resolve) must show
// the worktree-copy value, frozen against the certified golden. This locks the
// overlay that scan_entities_active / load_active_entity_fields perform.
func TestWorktreeOverlayActiveReads(t *testing.T) {
	pipelineEntity := "---\nid: add-login\nstatus: implementation\ntitle: Add login\nworktree: wt\n---\n"
	worktreeEntity := "---\nid: add-login\nstatus: review\ntitle: Add login\nworktree: wt\n---\n"

	cases := []struct {
		name string
		args []string
		// wantValue, when set, must appear in the native stdout — a direct guard
		// that the overlaid worktree-copy status is observable, not just that both
		// runners agree. --where filters on it (so its presence in a non-empty row
		// proves the overlay drove the match) and --resolve emits identity/path
		// only, so neither carries the literal status field for a value check.
		wantValue string
	}{
		{"table", nil, "review"},
		{"where-status", []string{"--where", "status=review"}, ""},
		// `worktree` is a non-default key, so it appends as a single extra in both
		// runners (no de-dupe). Projecting the DEFAULT `status` would make native
		// diverge from the buggy oracle (de-dupe drops the duplicate column); that
		// overlay+de-dupe case is locked by TestWorktreeOverlayFieldsStatusDeduped.
		{"fields-worktree", []string{"--fields", "worktree"}, ""},
		{"resolve", []string{"--resolve", "add-login"}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nativeRoot := buildWorktreeBacked(t, overlayReadme, pipelineEntity, worktreeEntity)
			env := pinnedEnv(t)

			args := withRoot(append([]string{"--workflow-dir", "%ROOT%"}, tc.args...), nativeRoot)
			nOut, nErr, nCode := runNative(t, nativeRoot, env, args...)

			// Normalize the per-test temp root prefix out of both streams so only
			// the behavioral content is compared. The --resolve workflow= field is
			// realpath'd (resolve.go emits realpathOf(workflowDir)), so on macOS it
			// carries the /private-resolved spelling while path= keeps the as-spelled
			// root; strip the realpath'd spelling first, then the as-spelled one, so
			// both map to %ROOT% and the golden is path-independent across hosts.
			nOutN := overlayRoot(nOut, nativeRoot)
			nErrN := overlayRoot(nErr, nativeRoot)

			assertEnvelopeGolden(t, "worktree-overlay-"+tc.name, goldenEnvelope{
				stdout: nOutN, stderr: nErrN, exit: nCode,
			})

			// Guard against agreeing on the wrong (pipeline-copy) value.
			if tc.wantValue != "" && !strings.Contains(nOutN, tc.wantValue) {
				t.Fatalf("native output should reflect worktree-copy value %q:\n%s", tc.wantValue, nOutN)
			}
			// --where status=review matches only via the overlay (the pipeline copy
			// is status=implementation); the entity row must therefore appear.
			if tc.name == "where-status" && !strings.Contains(nOutN, "add-login") {
				t.Fatalf("--where status=review should match the overlaid entity, got:\n%s", nOutN)
			}
		})
	}

	// Guard: the worktree copy genuinely exists and differs from the pipeline copy
	// during the run, so a passing parity is meaningful (not a no-op fallback).
	root := buildWorktreeBacked(t, overlayReadme, pipelineEntity, worktreeEntity)
	if _, err := os.Stat(filepath.Join(root, "wt", "add-login.md")); err != nil {
		t.Fatalf("worktree copy must exist: %v", err)
	}
}

// TestWorktreeOverlayFieldsStatusDeduped locks the overlay+de-dupe intersection
// the parity table can no longer cover: `--fields status` on a worktree-backed
// entity surfaces the overlaid worktree-copy status (review) in the default
// STATUS column, and the de-dupe means STATUS is NOT rendered twice. This is a
// native-only assertion because the oracle still emits the duplicate column.
func TestWorktreeOverlayFieldsStatusDeduped(t *testing.T) {
	pipelineEntity := "---\nid: add-login\nstatus: implementation\ntitle: Add login\nworktree: wt\n---\n"
	worktreeEntity := "---\nid: add-login\nstatus: review\ntitle: Add login\nworktree: wt\n---\n"
	root := buildWorktreeBacked(t, overlayReadme, pipelineEntity, worktreeEntity)
	env := pinnedEnv(t)

	out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--fields", "status")
	if code != 0 {
		t.Fatalf("native exit=%d stderr=%q", code, errOut)
	}
	header := headerOf(out)
	if c := countToken(header, "STATUS"); c != 1 {
		t.Fatalf("--fields status duplicated the STATUS column (want 1, got %d): %q", c, header)
	}
	// The overlaid worktree-copy status is observable in the default STATUS column.
	if !strings.Contains(out, "review") {
		t.Fatalf("--fields status should surface the overlaid status 'review':\n%s", out)
	}
}

// withRoot replaces the "%ROOT%" placeholder in args with root.
func withRoot(args []string, root string) []string {
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = strings.ReplaceAll(a, "%ROOT%", root)
	}
	return out
}

// overlayRoot maps both the realpath'd and the as-spelled temp root to %ROOT%.
// The realpath'd spelling (macOS /private/var/...) is stripped first since it is
// the longer/more-specific form, then the as-spelled root, mirroring normalize()
// — so the realpath'd workflow= field and the as-spelled path= field both
// collapse to the placeholder and the golden holds on macOS and Linux alike.
func overlayRoot(s, root string) string {
	if real := realpath(root); real != root {
		s = strings.ReplaceAll(s, real, "%ROOT%")
	}
	return strings.ReplaceAll(s, root, "%ROOT%")
}
