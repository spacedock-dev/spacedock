// ABOUTME: AC-4 verdict gate — --set refuses a terminal-status transition whose
// ABOUTME: post-state verdict is empty; verdict-present and --force both pass.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestTerminalSetWithoutVerdictRefused is AC-4: a --set that brings status to a
// terminal stage with an empty post-state verdict is refused (exit 1, entity
// unmutated). seq-workflow has no merge hook and no mod-block, so the verdict
// gate is the ONLY guard that can fire — `--set 002-vendor-script status=done`
// succeeds today (proven baseline) and must now refuse for a missing verdict.
func TestTerminalSetWithoutVerdictRefused(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done")
	out, errOut, code := runNative(t, root, env, args...)

	if code != 1 {
		t.Fatalf("terminal --set with no verdict must refuse (exit 1), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(errOut, "verdict") {
		t.Fatalf("stderr must name the missing verdict, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if strings.Contains(fm, "status: done") {
		t.Fatalf("entity must stay unmutated when the verdict gate refuses:\n%s", fm)
	}
}

// TestTerminalSetWithVerdictAllowed is AC-4: the same terminal transition with a
// verdict set in the SAME --set call is allowed (exit 0, status advances).
func TestTerminalSetWithVerdictAllowed(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "verdict=passed")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("terminal --set WITH a verdict must succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: done") || !strings.Contains(fm, "verdict: passed") {
		t.Fatalf("entity should have advanced to done with a verdict:\n%s", fm)
	}
}

// TestTerminalSetWithoutVerdictForceBypasses is AC-4: --force bypasses the
// verdict gate (the same escape hatch the mod-block / merge-hook guards honor).
func TestTerminalSetWithoutVerdictForceBypasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "--force")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("--force must bypass the verdict gate (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: done") {
		t.Fatalf("entity should have advanced to done under --force:\n%s", fm)
	}
}
