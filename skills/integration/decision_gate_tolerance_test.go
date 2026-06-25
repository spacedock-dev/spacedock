// ABOUTME: Decision-gate tolerance smoke — a gate's self-described `decision`
// ABOUTME: block (the Bridge↔FO hookup for CLOSE/KEEP gates) is an OPTIONAL stage
// ABOUTME: key Bridge parses directly; this locks that Spacedock's own `status
// ABOUTME: --read` tolerates it (ignores it) rather than failing the parse.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestStatusReadToleratesDecisionBlock locks that a workflow whose gate declares a
// machine-readable `decision` block still parses through `status --read`: the
// command exits 0 and emits every stage. Bridge consumes the block from the README
// directly; Spacedock must not choke on it (the FO's `status --set`/boot reads run
// against these same workflows).
func TestStatusReadToleratesDecisionBlock(t *testing.T) {
	dir := t.TempDir()
	readme := filepath.Join(dir, "README.md")
	content := "---\ncommissioned-by: spacedock@1.0\nentity-type: ticket\nid-style: slug\n" +
		"stages:\n  states:\n    - name: review\n      initial: true\n" +
		"    - name: escalated\n      gate: true\n      feedback-to: review\n" +
		"      decision:\n        field: verdict\n        options:\n" +
		"          - {label: Close, value: CLOSED, handoff: fo}\n" +
		"          - {label: Keep, value: IMPROVED, handoff: fo}\n" +
		"    - name: reviewed\n      terminal: true\n---\n# rev\n"
	if err := os.WriteFile(readme, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(spacedockBinary(t), "status", "--read", readme, "--json")
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("status --read on a decision-block workflow failed: %v\n%s", err, out)
	}
	got := string(out)
	for _, stage := range []string{"review", "escalated", "reviewed"} {
		if !strings.Contains(got, `"name":"`+stage+`"`) {
			t.Errorf("status --read did not emit stage %q (decision block broke the parse):\n%s", stage, got)
		}
	}
}
