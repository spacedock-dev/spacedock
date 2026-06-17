// ABOUTME: AC-2/AC-3 structural checks for the commission templates — the universal
// ABOUTME: doctrine has a single contract owner (not restated per-template) and each template leads with its outcome.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// commissionTemplates is the set of shipped commission-template files the
// defer-to-contract restructure governs, each with the contract file that
// solely owns the universal doctrine markers it used to restate.
var commissionTemplates = []string{
	"skills/commission/references/templates/development.md",
	"skills/commission/references/templates/experiment.md",
	"skills/commission/references/templates/refinement.md",
}

// The AC-2 doctrine single-source guard was RETIRED here per doc_test.go's policy
// (this package's own rule, line 11): "Do not add prose-to-code consistency checks
// as behavior substitutes ... if no behavior test exists yet, delete the read and
// report the owed test." It walked skills/+agents/ for a hardcoded doctrine phrase
// (`strings.Contains(b, "Prefer a code gate over a prose-only rule")`, etc.) and
// claimed to red when "the doctrine drifts back into a template". A doctrine is a
// MEANING, not a literal: a paraphrased restatement ("Favor an enforceable code gate
// rather than a prose-only rule…") carries the doctrine while dropping the bytes, so
// the byte-grep stayed green on exactly the drift it advertised. That is the banned
// prose-grep — a literal-phrase substring used as a proxy for whether a meaning is
// present — and unlike the sibling claudeTeamDispatchTokens absence check (where the
// literal command token IS the thing and a meaning-inverting paraphrase necessarily
// drops it), no token here IS the doctrine; detecting paraphrase-drift requires
// interpreting prose, which a machine cannot. Narrowing it to "verbatim-phrase
// absence" would still be a prose-phrase-absence standing in for the meaning, i.e.
// the same banned grep. OWED PROOF (#388 "templates defer doctrine to the contract"):
// the human gate-review of the commission templates, the same review that owns the
// qualitative lede judgement in TestTemplatesLeadWithOutcome below. The structural
// half of AC-2 that a machine CAN see survives as TestTemplatesCarryWorkflowSpecificRulesSlot.

// TestTemplatesCarryWorkflowSpecificRulesSlot is the AC-2 heading-presence half: each
// commission template carries a `## Workflow-specific rules` section. This is a
// structural section-presence property a machine can see — the slot that holds each
// shape's unique rules after the generic doctrine is deferred to the contract — not a
// doctrine-substring match.
func TestTemplatesCarryWorkflowSpecificRulesSlot(t *testing.T) {
	root := repoRoot(t)
	const heading = "## Workflow-specific rules"
	for _, rel := range commissionTemplates {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Errorf("read %s: %v", rel, err)
			continue
		}
		if !strings.Contains(string(data), heading) {
			t.Errorf("%s carries no %q section — the workflow-specific-rules slot is missing", rel, heading)
		}
	}
}

// TestTemplatesLeadWithOutcome is the AC-3 structural lede-position check: each
// template opens with a non-empty prose paragraph (the workflow's outcome) positioned
// between its frontmatter and the first `## ` mechanics heading. This is a
// positional/ordering assertion over the file — lead-with-the-end shape is
// structurally present — not a value-judgement substring match; the qualitative
// "does the lede actually lead with value" judgement is a human gate-review item.
func TestTemplatesLeadWithOutcome(t *testing.T) {
	root := repoRoot(t)
	for _, rel := range commissionTemplates {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Errorf("read %s: %v", rel, err)
			continue
		}
		doc := string(data)
		fm, ok := frontmatter(doc)
		if !ok {
			t.Errorf("%s has no parseable frontmatter — cannot locate the lede", rel)
			continue
		}
		// The body begins after the closing `---` of the frontmatter block.
		body := doc[strings.Index(doc, fm)+len(fm):]
		if i := strings.Index(body, "\n---"); i >= 0 {
			body = body[i+len("\n---"):]
		}
		firstHeading := strings.Index(body, "\n## ")
		if firstHeading < 0 {
			t.Errorf("%s has no `## ` mechanics heading — cannot bound the lede", rel)
			continue
		}
		// The title `# ...` line is the workflow name, not the lede; the lede is the
		// prose between the title and the first `## ` mechanics heading.
		lede := body[:firstHeading]
		lede = strings.TrimSpace(stripTitleLine(lede))
		if lede == "" {
			t.Errorf("%s has no lede paragraph before its first `## ` mechanics heading", rel)
		}
	}
}

// stripTitleLine removes a leading `# Title` line (and surrounding blank lines) from
// a template body so the lede check measures the prose paragraph, not the H1.
func stripTitleLine(body string) string {
	body = strings.TrimLeft(body, "\n")
	if strings.HasPrefix(body, "# ") {
		if nl := strings.Index(body, "\n"); nl >= 0 {
			return body[nl+1:]
		}
		return ""
	}
	return body
}
