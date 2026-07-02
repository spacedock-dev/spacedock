// ABOUTME: Contract anchor — «gate.assemble-verdict» drains the inbox BEFORE
// ABOUTME: presenting a gate, so a captain decision queued from Bridge for that
// ABOUTME: gate is honored instead of colliding with a redundant terminal prompt.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGateAssembleVerdictDrainsBeforePresenting locks the drain-before-present
// fix: when the FO reaches a gate, it must first fire that entity's idle hooks
// (draining a Bridge-queued `decision` record) and skip the presentation if the
// decision already resolved the gate. Bridge wake is best-effort, while delivery
// is confirmed only by the FO-owned drain and ack; without this a queued captain
// decision can still collide with a redundant terminal prompt.
func TestGateAssembleVerdictDrainsBeforePresenting(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "first-officer-shared-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	content := string(data)

	for _, r := range []string{
		"drain before presenting",
		"Bridge wake is best-effort",
		"FO-owned drain and ack",
		"do NOT present",
	} {
		if !strings.Contains(content, r) {
			t.Errorf("«gate.assemble-verdict» no longer drains the inbox before presenting a gate: missing %q.\n"+
				"Without it, a Bridge-queued captain decision collides with a redundant terminal gate prompt.", r)
		}
	}

	// Ordering: the drain effect must precede the decide/present effect — draining
	// after presenting would not prevent the redundant prompt.
	drainAt := strings.Index(content, "drain before presenting")
	decideAt := strings.Index(content, "effect — decide")
	if drainAt < 0 || decideAt < 0 || drainAt > decideAt {
		t.Errorf("the drain effect must come BEFORE the decide/present effect in «gate.assemble-verdict» (drain@%d, decide@%d)", drainAt, decideAt)
	}
}
