// ABOUTME: AC-6 oracle over the FO event-loop step-0 reconcile prose — roster
// ABOUTME: reconciliation needs a team identity; bare reconcile is git-only D/E.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// claudeFORuntime reads the vendored Claude first-officer runtime contract text.
func claudeFORuntime(t *testing.T) string {
	t.Helper()
	p := filepath.Join(skillsRoot(t), "first-officer", "references", "claude-first-officer-runtime.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read Claude FO runtime: %v", err)
	}
	return string(b)
}

// reconcileStep0Region returns the body of the FO event-loop step-0 reconcile
// item: from the `0. **Reconcile sweep.**` list line up to (but excluding) the
// `1.` list item that follows. Scopes the AC-6 assertions to the one numbered
// step that documents reconcile, so the check is a structural region oracle, not
// a free-floating grep over the whole 34KB file.
func reconcileStep0Region(t *testing.T, text string) string {
	t.Helper()
	lines := strings.Split(text, "\n")
	start := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "0. **Reconcile sweep.**") {
			start = i
			break
		}
	}
	if start == -1 {
		t.Fatal("FO runtime has no `0. **Reconcile sweep.**` event-loop step")
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "1. ") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}

// TestReconcileStep0RequiresTeamIdentityForRoster locks AC-6: the FO event-loop
// step-0 prose must state that the roster-derived classes (A/B/C) require a team
// identity — an explicit --team-name or a current-session match — and that bare
// reconcile without one is git-only (Class D/E). This is the prose half of the
// guarantee the AC-1/AC-4 code gates enforce; pairing them keeps the contract
// from silently re-describing bare reconcile as a safe default.
//
// Oracle: a structural region scan over the single step-0 list item. It asserts
// the region binds the team-identity precondition to the roster classes and the
// git-only fallback to bare invocation — not the mere presence of a phrase
// anywhere in the file. A paraphrase that drops the precondition (re-inviting the
// unsafe bare-fallback) fails this.
func TestReconcileStep0RequiresTeamIdentityForRoster(t *testing.T) {
	region := reconcileStep0Region(t, claudeFORuntime(t))

	// The region must tie roster reconciliation to a required team identity.
	if !strings.Contains(region, "team identity") {
		t.Errorf("step-0 region does not state roster reconciliation needs a team identity:\n%s", region)
	}
	// It must name both ways to supply that identity: explicit --team-name and
	// a current-session match.
	if !strings.Contains(region, "--team-name") {
		t.Errorf("step-0 region does not reference explicit --team-name:\n%s", region)
	}
	if !strings.Contains(region, "current-session") && !strings.Contains(region, "current session") {
		t.Errorf("step-0 region does not reference the current-session match:\n%s", region)
	}
	// It must state that bare reconcile (no team identity) is git-only D/E.
	if !strings.Contains(region, "git-only") {
		t.Errorf("step-0 region does not state bare reconcile is git-only:\n%s", region)
	}
	if !strings.Contains(region, "D/E") {
		t.Errorf("step-0 region does not name the git/filesystem classes D/E for the bare path:\n%s", region)
	}
}

// TestReconcileStep0DropsOptionalTeamNameFraming locks the second half of AC-6:
// the step-0 invocation line must no longer frame --team-name as a free-floating
// optional flag (`[--team-name {team_name}]`) whose omission silently falls back
// to an unsafe heuristic. The bracketed-optional form is exactly the wording the
// entity flagged as inviting the bare unsafe invocation.
func TestReconcileStep0DropsOptionalTeamNameFraming(t *testing.T) {
	region := reconcileStep0Region(t, claudeFORuntime(t))
	if strings.Contains(region, "[--team-name {team_name}]") {
		t.Errorf("step-0 region still frames --team-name as a bracketed-optional flag, re-inviting the unsafe bare fallback:\n%s", region)
	}
}
