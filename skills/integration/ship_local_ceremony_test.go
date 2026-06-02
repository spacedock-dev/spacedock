// ABOUTME: Structural lint over the FO shared core's ship-local ceremony block —
// ABOUTME: the named subsection exists and references the merge-policy branch and sentinel.
package integration

import (
	"strings"
	"testing"
)

// subsectionAfter returns the body of the markdown subsection beginning at the
// line equal to heading, up to (but excluding) the next heading at the same or a
// higher level (`### ` or `## `). It is `sectionAfter`'s `###`-aware sibling,
// needed because a `### ` subsection is otherwise swept past `### ` siblings.
func subsectionAfter(text, heading string) string {
	lines := strings.Split(text, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == heading {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") || strings.HasPrefix(lines[i], "### ") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}

// TestShipLocalCeremonyBlockExists is a structural lint: the FO shared core
// carries a single named ship-local ceremony block that references the
// merge-policy branch and the local-merge sentinel. It scopes to the ceremony
// subsection — a real on-disk anchor whose absence (or a missing policy/sentinel
// reference) breaks the documented ship-local flow. The no-`--force` behavioral
// property of the terminal/local-merge transition is enforced by the real
// status codepath in internal/status (merge_policy_guard_test.go:
// TestMergeLocalNoSentinelTerminalSetSucceeds and siblings), so this lint does
// not grep the ceremony prose for force-related wording.
func TestShipLocalCeremonyBlockExists(t *testing.T) {
	fo := vendoredSkillFiles(t)["first-officer/references/first-officer-shared-core.md"]
	region := subsectionAfter(fo, "### Ship-Local Ceremony")
	if region == "" {
		t.Fatal("FO shared core missing the `### Ship-Local Ceremony` block (AC-5)")
	}
	if !strings.Contains(region, "merge: local") {
		t.Error("ship-local ceremony must reference the `merge: local` policy branch")
	}
	if !strings.Contains(region, "local-merge:") {
		t.Error("ship-local ceremony must name the local-merge sentinel")
	}
}
