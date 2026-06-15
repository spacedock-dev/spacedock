// ABOUTME: AC-3 `new --help` coverage — the per-command synopsis and flag surface
// ABOUTME: render instead of the inherited top-level grouped menu.
package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestNewHelpRendersCommandUsage pins AC-3: `spacedock new --help` shows the
// per-command synopsis and its full flag surface, and does NOT fall through to
// the top-level grouped menu (the Launch section header is the general-menu
// marker).
func TestNewHelpRendersCommandUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"new", "--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"spacedock new [--folder] SLUG",
		"--workflow-dir",
		"--folder",
		"--id-seed",
		"--id-actor",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("new --help missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Launch") {
		t.Errorf("new --help leaked the top-level grouped menu (Launch header):\n%s", out)
	}
}
