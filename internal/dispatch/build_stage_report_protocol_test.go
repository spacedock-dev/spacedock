// ABOUTME: AC-1 — the Pi First-action block no longer overclaims the full
// ABOUTME: ensign discipline; the stage-report format is attributed to the body
// ABOUTME: and the rest to the ensign skill.
package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildPiFirstActionNarrowedToStageReportFormat (AC-1) asserts the Pi
// First-action claim no longer overclaims the full ensign discipline (polling,
// worktree ownership, completion protocol) and instead attributes the
// stage-report format to the body and the rest to the ensign skill.
func TestBuildPiFirstActionNarrowedToStageReportFormat(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-first-action"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	// The overclaim is gone.
	for _, overclaim := range []string{
		"This file contains the shared ensign discipline entry points",
	} {
		if strings.Contains(body, overclaim) {
			t.Fatalf("pi First-action still overclaims: %q present in body:\n%s", overclaim, body)
		}
	}
	// The narrowed claim attributes the format template to the body and the
	// rest of the discipline to the ensign skill.
	for _, want := range []string{
		"This file carries the stage-report format template",
		"The ensign skill supplies the remaining shared discipline",
		"not auto-loaded",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("pi First-action missing narrowed claim %q:\n%s", want, body)
		}
	}
}
