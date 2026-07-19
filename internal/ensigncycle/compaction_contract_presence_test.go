// ABOUTME: Contract-presence guard binding the AC-1/AC-2 offline oracles to the shipped
// ABOUTME: shared-core rule, so the acceptance gate fails if the compaction rule is deleted.
package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCompactionContinuityRuleShipped is a contract-presence regression guard, NOT
// behavioral proof of the value ACs. The AC-1 timing oracle and the AC-2 reload oracle
// (compaction_timing_test.go, compaction_reload_test.go) run over hand-authored fixtures,
// so on their own they stay green even if the "## Compaction continuity" rule were
// removed from first-officer-shared-core.md. This guard closes that gap by asserting the
// shipped rule and its load-bearing clauses are present: deleting the rule turns the
// offline acceptance gate RED. It proves the rule SHIPS with the requirements the oracles
// encode; it does NOT prove the FO obeys the rule at runtime — that behavioral linkage is
// the opt-in live path (test plan item 4), outside this offline gate.
func TestCompactionContinuityRuleShipped(t *testing.T) {
	root := postCompactRepoRoot(t)
	path := filepath.Join(root, "skills", "first-officer", "references", "first-officer-shared-core.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read shared-core contract: %v", err)
	}

	section := extractMarkdownSection(string(body), "## Compaction continuity")
	if section == "" {
		t.Fatalf("shared-core is missing the '## Compaction continuity' rule — the AC-1/AC-2 oracles would be false-green without it")
	}

	// Load-bearing clauses the AC-1 (before-compaction timing) and AC-2 (after-compaction
	// reload) oracles depend on. Each is a requirement the oracle encodes; if the shipped
	// rule drops it, the corresponding oracle no longer reflects the contract.
	clauses := []struct {
		ac, phrase string
	}{
		{"AC-1", "durable, recoverable boundary"},     // safe-to-compact timing condition
		{"AC-1", "safe time to compact"},              // the non-blocking suggestion itself
		{"AC-2", "reread the authoritative contract"}, // reload trigger
		{"AC-2", "first-officer-shared-core.md"},      // one of the three reread targets
		{"AC-2", "active host runtime adapter"},       // the host-adapter reread target
		{"AC-2", "«post-compact-notice»"},             // the per-host delivery binding link
		{"AC-2", "never authoritative"},               // the compacted summary is not the contract
	}
	for _, c := range clauses {
		if !strings.Contains(section, c.phrase) {
			t.Errorf("%s: the '## Compaction continuity' rule is missing load-bearing clause %q; the offline oracle would no longer bind to the shipped contract", c.ac, c.phrase)
		}
	}
}

// extractMarkdownSection returns the body of the section whose heading begins with
// heading, up to the next same-or-higher-level heading, or "" if absent. It bounds the
// clause checks to the rule's own span so a phrase elsewhere in the file cannot satisfy
// the guard.
func extractMarkdownSection(doc, heading string) string {
	lines := strings.Split(doc, "\n")
	start := -1
	for i, ln := range lines {
		if strings.HasPrefix(ln, heading) {
			start = i
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}
