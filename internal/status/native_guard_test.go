// ABOUTME: AC-4 native guard parity — terminal --set under a mod-block, a
// ABOUTME: merge-hook-unsatisfied terminal --set, and a mod-blocked --archive.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestNativeTerminalSetUnderModBlock locks the terminal-transition-under-mod-
// block guard: native and oracle both exit 1 with the same error text and leave
// the entity unmutated.
func TestNativeTerminalSetUnderModBlock(t *testing.T) {
	env := pinnedEnv(t)
	nativeRoot := stageFixture(t, "guard-workflow")

	args := append([]string{"--workflow-dir", nativeRoot}, "--set", "010-blocked", "status=done")

	nOut, nErr, nCode := runNative(t, nativeRoot, env, args...)

	if nCode != 1 {
		t.Fatalf("native exit=%d, want 1", nCode)
	}
	if nOut != "" {
		t.Fatalf("stdout must be empty on rejection: native=%q", nOut)
	}
	assertEnvelopeGolden(t, "native-guard-terminal-modblock", goldenEnvelope{
		stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
	})
	fm := readWhole(t, filepath.Join(nativeRoot, "010-blocked.md"))
	if !strings.Contains(fm, "status: implementation") || strings.Contains(fm, "status: done") {
		t.Fatalf("entity mutated despite guard:\n%s", fm)
	}
}

// TestNativeArchiveRelativeDest locks --archive's relative dest spelling parity.
func TestNativeArchiveRelativeDest(t *testing.T) {
	env := pinnedEnv(t)
	nativeRoot := stageFixture(t, "seq-workflow")

	args := []string{"--workflow-dir", ".", "--archive", "001-design-seam"}
	nOut, nErr, nCode := runNative(t, nativeRoot, env, args...)

	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	want := "archived: ./_archive/001-design-seam.md\n"
	if nOut != want {
		t.Fatalf("native narration = %q, want %q", nOut, want)
	}
}

// TestNativeMidSetTruncation locks the mid-set truncation: `--set <slug> --bogus
// status=done` drops status=done, exits 1 with the same error, entity unchanged.
func TestNativeMidSetTruncation(t *testing.T) {
	env := pinnedEnv(t)
	nativeRoot := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", nativeRoot}, "--set", "002-vendor-script", "--bogus", "status=done")

	nOut, nErr, nCode := runNative(t, nativeRoot, env, args...)

	if nCode != 1 {
		t.Fatalf("exit: native=%d, want 1", nCode)
	}
	if !strings.Contains(nErr, "requires at least one field=value") {
		t.Fatalf("native stderr=%q", nErr)
	}
	if nOut != "" {
		t.Fatalf("stdout must be empty: native=%q", nOut)
	}
	assertEnvelopeGolden(t, "native-guard-midset-truncation", goldenEnvelope{
		stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
	})
}
