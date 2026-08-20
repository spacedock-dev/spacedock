// ABOUTME: AC-2/AC-3/AC-4 coverage — split-root state-commit guidance degrades to
// ABOUTME: local-only when the state checkout has no origin, and keeps push/pull when it does.
package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildSplitRootDispatchBody drives `dispatch build` over a split-root workflow
// whose state checkout lives at workflowDir/state-checkout, returning the emitted
// dispatch body. withOrigin adds a named origin remote to the git repo the state
// checkout resolves against, so stateHasOrigin reports true.
func buildSplitRootDispatchBody(t *testing.T, withOrigin bool) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()

	workflowDir := root
	stateCheckout := filepath.Join(workflowDir, "state-checkout")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(true))

	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}

	entityPath := filepath.Join(stateCheckout, "thing", "index.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))

	gitInit(t, root)
	if withOrigin {
		gitAddOrigin(t, root)
	}

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   workflowDir,
		"stage":          "implementation",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", workflowDir)
	return readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))
}

// TestStateCommitGuidanceNoOriginDropsRemoteSync (AC-2, AC-4) pins the degrade: a
// split-root state checkout with NO origin keeps the path-scoped local commit but
// emits a local-only line and drops both `git push origin` and
// `git pull --rebase origin`.
func TestStateCommitGuidanceNoOriginDropsRemoteSync(t *testing.T) {
	body := buildSplitRootDispatchBody(t, false)

	// Path-scoped local commit instruction is retained.
	if !strings.Contains(body, "git -C ") || !strings.Contains(body, " add ") {
		t.Fatalf("no-origin body dropped the path-scoped commit instruction\n--- body ---\n%s", body)
	}
	if !strings.Contains(body, "never a bare `git add -A`") {
		t.Fatalf("no-origin body dropped the concurrency-safety phrase\n--- body ---\n%s", body)
	}

	// The impossible remote-sync commands are gone. The shipped wording is
	// `git -C <checkout> push origin <branch>` / `pull --rebase origin`, so the
	// remote-sync verbs are the discriminators.
	if strings.Contains(body, "push origin") {
		t.Errorf("no-origin body still instructs `push origin`\n--- body ---\n%s", body)
	}
	if strings.Contains(body, "pull --rebase origin") {
		t.Errorf("no-origin body still instructs `pull --rebase origin`\n--- body ---\n%s", body)
	}

	// AC-4: the local-only mode is named, not silent.
	if !strings.Contains(body, "no `origin` remote") {
		t.Errorf("no-origin body does not name the missing-origin condition\n--- body ---\n%s", body)
	}
	if !strings.Contains(body, "local-only") {
		t.Errorf("no-origin body does not name the local-only state mode\n--- body ---\n%s", body)
	}
}

// TestStateCommitGuidanceWithOriginKeepsRemoteSync (AC-3) pins the unchanged
// origin path: a split-root state checkout WITH an origin keeps the push and
// pull-rebase reminders verbatim and does NOT carry the local-only line.
func TestStateCommitGuidanceWithOriginKeepsRemoteSync(t *testing.T) {
	body := buildSplitRootDispatchBody(t, true)

	if !strings.Contains(body, "push origin") {
		t.Errorf("origin body dropped the `push origin` reminder\n--- body ---\n%s", body)
	}
	if !strings.Contains(body, "pull --rebase origin") {
		t.Errorf("origin body dropped the `pull --rebase origin` reminder\n--- body ---\n%s", body)
	}
	if strings.Contains(body, "local-only") {
		t.Errorf("origin body leaks the no-origin local-only line\n--- body ---\n%s", body)
	}
}
