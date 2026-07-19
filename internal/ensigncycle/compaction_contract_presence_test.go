// ABOUTME: AC-1 contract-SHAPE presence check — asserts first-officer-shared-core.md ships the
// ABOUTME: "## Compaction continuity" rule with both clauses, host-neutral. Makes NO runtime-behavior claim.
package ensigncycle

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestCompactionContinuityRuleShippedInContractShape is a CONTRACT-SHAPE check for AC-1,
// and only that. It reads the shipped first-officer-shared-core.md and asserts the
// "## Compaction continuity" rule is present with both required clauses — a before-compaction
// durable-boundary suggest, and an after-compaction reread-and-reconcile that names the three
// reread targets — and that the rule names no host (host-neutral).
//
// It asserts the CONTRACT SHAPE only. It makes NO runtime-behavior claim: it does not assert
// the FO produces exactly one suggestion, nor that the next workflow effect occurs only after
// the reads. The linkage from this prose to a live FO is an unenforced judgment-rule property
// (the design ships two judgment rules, no controller); its proof is the out-of-scope live
// compaction scenario, not this test. This is a contract-shape check for a contract-shape AC —
// not a prose guard dressed as behavioral proof.
//
// The shared-core byte ratchet is enforced by the existing internal/contractlint startup_collapse
// (TestStartupRecipeCollapsedAndLeaner, ceiling preChangeSharedCoreBytes), which AC-1 names as the
// byte-ceiling ratchet; it is not re-implemented here.
func TestCompactionContinuityRuleShippedInContractShape(t *testing.T) {
	root := postCompactRepoRoot(t)
	path := filepath.Join(root, "skills", "first-officer", "references", "first-officer-shared-core.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read shared-core contract: %v", err)
	}

	section := extractMarkdownSection(string(body), "## Compaction continuity")
	if section == "" {
		t.Fatalf("shared-core is missing the '## Compaction continuity' rule")
	}

	// Both required clauses ship, keyed on their load-bearing phrases. This asserts the
	// rule SHIPS with these clauses (contract shape); it is NOT a claim the FO obeys them.
	clauses := []struct {
		clause, phrase string
	}{
		{"before-compaction durable-boundary suggest", "**Before compaction:**"},
		{"before-compaction durable-boundary suggest", "durable, recoverable boundary"},
		{"before-compaction durable-boundary suggest", "safe time to compact"},
		{"after-compaction reread-and-reconcile", "**After compaction**"},
		{"after-compaction reread-and-reconcile", "reread the authoritative contract"},
		{"after-compaction reread target", "SKILL.md"},
		{"after-compaction reread target", "first-officer-shared-core.md"},
		{"after-compaction reread target", "active host runtime adapter"},
		{"after-compaction reconcile-before-effect", "before the next workflow effect"},
		{"compacted summary is not the contract", "never authoritative"},
		{"per-host delivery binding link", "«post-compact-notice»"},
	}
	for _, c := range clauses {
		if !strings.Contains(section, c.phrase) {
			t.Errorf("the '## Compaction continuity' rule is missing the %s clause phrase %q", c.clause, c.phrase)
		}
	}

	// Host-neutral: the rule names no host inside its own text (0 host-specific terms).
	// Host-specific delivery lives in the fo-dispatch-core «post-compact-notice» binding,
	// not in this rule.
	if m := hostNameRe.FindString(section); m != "" {
		t.Errorf("the '## Compaction continuity' rule names host %q — the rule must be host-neutral", m)
	}
}

// hostNameRe matches a host NAME as a whole word. Naming a host inside the host-neutral
// compaction-continuity rule would couple the rule to a host mechanism.
var hostNameRe = regexp.MustCompile(`\b(Codex|Claude|Pi)\b`)

// extractMarkdownSection returns the body of the section whose heading begins with
// heading, up to the next same-or-higher-level heading, or "" if absent. It bounds the
// clause checks to the rule's own span so a phrase elsewhere in the file cannot satisfy
// the check.
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
