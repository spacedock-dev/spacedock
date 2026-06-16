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

// doctrineMarker pins one canonical contract phrasing to its sole owning file.
// The marker set binds two independent, divergeable sources — the contract's own
// wording and whatever a template carries — so the check reds precisely when a
// template re-restates the doctrine, not merely because we wrote the words.
type doctrineMarker struct {
	phrase string
	owner  string
}

var doctrineMarkers = []doctrineMarker{
	{"Prefer a code gate over a prose-only rule", "skills/first-officer/references/first-officer-shared-core.md"},
	{"wording-present is not behavior", "skills/first-officer/references/first-officer-shared-core.md"},
	{"Prove by exercising, not by re-reading", "skills/ensign/references/ensign-shared-core.md"},
	{"A substring search is not proof of behavior", "skills/ensign/references/ensign-shared-core.md"},
}

// TestUniversalDoctrineHasSingleSource is the AC-2 single-source / dedup check: each
// universal-doctrine marker the templates used to paraphrase appears in exactly ONE
// file under skills/+agents/ — its owning contract — and a second copy (a template
// re-restating it) reds the test. This is not a presence tautology: it binds the
// contract's content against every other shipped file, so it fails when the doctrine
// drifts back into a template, the drift regression the restructure prevents.
func TestUniversalDoctrineHasSingleSource(t *testing.T) {
	root := repoRoot(t)
	for _, marker := range doctrineMarkers {
		var sources []string
		for _, tree := range []string{"skills", "agents"} {
			base := filepath.Join(root, tree)
			err := filepath.WalkDir(base, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() || filepath.Ext(p) != ".md" {
					return nil
				}
				b, readErr := os.ReadFile(p)
				if readErr != nil {
					return readErr
				}
				if strings.Contains(string(b), marker.phrase) {
					rel, _ := filepath.Rel(root, p)
					sources = append(sources, rel)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("walk %s: %v", tree, err)
			}
		}
		if len(sources) != 1 || sources[0] != marker.owner {
			t.Errorf("doctrine marker %q has sources %v, want exactly [%s] (single source of truth — a template re-restating it reds this)", marker.phrase, sources, marker.owner)
		}
	}
}

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
