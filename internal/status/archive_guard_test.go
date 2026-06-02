// ABOUTME: AC-3 --archive dest-spelling parity (relative vs absolute --workflow-dir)
// ABOUTME: and the terminal-transition-under-mod-block guard (exit 1, current text).
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestArchiveRelativeDest locks that --archive's dest tracks the --workflow-dir
// spelling: a relative `--workflow-dir .` (run with cwd=root) yields
// `archived: ./_archive/{slug}.md`. Compared launcher-vs-oracle, both run from
// the same relative spelling, and the moved file lands under _archive.
func TestArchiveRelativeDest(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := []string{"--workflow-dir", ".", "--archive", "001-design-seam"}
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	want := "archived: ./_archive/001-design-seam.md\n"
	if nOut != want {
		t.Fatalf("native narration = %q, want %q (relative dest spelling)", nOut, want)
	}
	// The entity actually moved under _archive and left the active dir.
	if _, err := os.Stat(filepath.Join(root, "_archive", "001-design-seam.md")); err != nil {
		t.Fatalf("archived file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "001-design-seam.md")); !os.IsNotExist(err) {
		t.Fatalf("source file should be gone after archive, stat err=%v", err)
	}
}

// TestArchiveAbsoluteDest locks the absolute-spelling case: an absolute
// --workflow-dir yields an absolute archived: dest. Compared launcher-vs-oracle
// with the root prefix normalized so no machine path is asserted literally.
func TestArchiveAbsoluteDest(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := []string{"--workflow-dir", root, "--archive", "001-design-seam"}
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	// Native emits an absolute dest (not realpath'd) under its own root.
	wantNative := "archived: " + filepath.Join(root, "_archive", "001-design-seam.md") + "\n"
	if nOut != wantNative {
		t.Fatalf("native narration = %q, want %q (absolute dest spelling)", nOut, wantNative)
	}
}

// TestTerminalSetUnderModBlockRejected locks the guard: a terminal --set
// (status -> terminal stage) on an entity with an active mod-block exits 1 with
// the current error text, and the entity is NOT mutated. Compared launcher vs
// oracle for exit code, stderr text, and unchanged frontmatter.
func TestTerminalSetUnderModBlockRejected(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "guard-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "010-blocked", "status=done")
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 1 {
		t.Fatalf("native exit=%d, want 1 (guard must reject)", nCode)
	}
	wantErr := "Error: entity 010-blocked has pending mod-block (merge:pr-merge). Clear mod-block in a separate --set call, or use --force."
	if !strings.Contains(nErr, wantErr) {
		t.Fatalf("native stderr = %q, want it to contain %q", nErr, wantErr)
	}
	if nOut != "" {
		t.Fatalf("stdout must be empty on rejection: native=%q", nOut)
	}
	assertEnvelopeGolden(t, "archive-guard-terminal-modblock", goldenEnvelope{
		stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
	})
	// The entity status must be unchanged (still implementation, not done).
	fm := readFrontmatter(t, filepath.Join(root, "010-blocked.md"))
	if !strings.Contains(fm, "status: implementation") {
		t.Fatalf("entity was mutated despite guard rejection:\n%s", fm)
	}
	if strings.Contains(fm, "status: done") {
		t.Fatalf("entity advanced to done despite guard rejection:\n%s", fm)
	}
}
