// ABOUTME: AC-1 contract-SHAPE presence check — asserts first-officer-shared-core.md ships the
// ABOUTME: "## Compaction continuity" rule and that the rule names no host. Makes NO runtime-behavior claim.
package ensigncycle

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestCompactionContinuityRuleShippedInContractShape is a minimal CONTRACT-SHAPE check for
// AC-1, and only that. It reads the shipped first-officer-shared-core.md and asserts the
// "## Compaction continuity" rule section is present and that the rule names no host
// (host-neutral: the `\b(Codex|Claude|Pi)\b` check). Host-specific delivery lives in the
// fo-dispatch-core «post-compact-notice» binding and its runtime adapters, not in this rule.
//
// It deliberately does NOT grep the rule's sentences: pinning exact phrases copied from the
// shipped .md is brittle prose-grep that breaks on any meaning-preserving reword and pins
// wording, not the rule. Clause content is confirmed by code review; the shared-core byte
// ratchet is enforced by the existing internal/contractlint startup_collapse
// (TestStartupRecipeCollapsedAndLeaner). This makes NO runtime-behavior claim — the linkage
// from the prose to a live FO is an unenforced judgment-rule property (two judgment rules, no
// controller), whose proof is the out-of-scope live compaction scenario, not this test.
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

	// Host-neutral: the rule names no host inside its own text (0 host-specific terms).
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
