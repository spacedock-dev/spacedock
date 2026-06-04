// ABOUTME: AC-3 presence lock — the dev template carries live-scenario-for-
// ABOUTME: runtime-claims guidance as a third opt-in recommended practice.
package hostneutrality

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// recommendedPracticesSectionRe scopes the presence check to the dev template's
// "## Recommended practices (opt-in)" section, the same scoping the existing
// External-proof / Detached-audit / Test-first locks use, so a stray mention of
// "live scenario" elsewhere in the file cannot satisfy the lock.
var recommendedPracticesSectionRe = regexp.MustCompile(`(?is)## Recommended practices \(opt-in\).*`)

// TestLiveScenarioRecommendedPracticePresent locks AC-3: the dev template
// (skills/commission/references/templates/development.md) names the live-
// scenario pattern for runtime claims as a THIRD opt-in recommended practice
// beside External-proof and Detached-audit. The required clauses encode the
// pattern's load-bearing distinction (a recording proves the WATCHER not the
// producer; a contract-text check proves the WORDS not the behaviour), so the
// lock rests on the guidance's substance, not merely a heading word. Same kind
// of presence check that guards the existing recommended-practice blocks; the
// claim is about the text itself, so proof at the claim's own level is legit.
func TestLiveScenarioRecommendedPracticePresent(t *testing.T) {
	path := filepath.Join("..", "..", "skills", "commission", "references", "templates", "development.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	section := recommendedPracticesSectionRe.FindString(string(body))
	if section == "" {
		t.Fatalf("%s missing the Recommended-practices section", path)
	}
	lower := strings.ToLower(section)

	// A heading naming the live-scenario practice — the third opt-in block.
	if !strings.Contains(lower, "live scenario") {
		t.Errorf("%s recommended-practices section does not name the live-scenario practice", path)
	}
	// The load-bearing distinction the pattern exists to teach.
	for _, want := range []string{
		"recording proves the watcher",
		"text",
		"durable",
	} {
		if !strings.Contains(lower, want) {
			t.Errorf("%s live-scenario block missing required substance %q", path, want)
		}
	}
}
