// ABOUTME: Locks the two surviving FO-contract edits — the hardened AC-coverage
// ABOUTME: cross-check and the new Detached Adversarial Audit section — over the real file.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// acCrossCheckParagraph isolates the `**AC coverage cross-check.**` paragraph
// inside the FO contract's `## Completion and Gates` section. The paragraph is
// one line in the source (markdown soft-wraps), bounded by the surrounding blank
// lines, so the assertions scope to exactly the cross-check clause and do not
// leak into the count-summary or reuse-conditions prose around it.
func acCrossCheckParagraph(t *testing.T) string {
	t.Helper()
	section := sectionAfter(foSharedCore(t), "## Completion and Gates")
	const marker = "**AC coverage cross-check.**"
	i := strings.Index(section, marker)
	if i < 0 {
		t.Fatalf("Completion and Gates section has no %q paragraph", marker)
	}
	rest := section[i:]
	// Paragraphs are blank-line delimited.
	if j := strings.Index(rest, "\n\n"); j >= 0 {
		return rest[:j]
	}
	return rest
}

// TestACCrossCheckRequiresExternalProof locks AC-1: the FO contract's AC-coverage
// cross-check requires evidence from a check OUTSIDE the entity body, refuses a
// self-oracle criterion (one whose only proof is review of the entity's own
// prose, which can never fail), and refuses advancing a pure-decision entity to a
// terminal PASSED verdict. It must also POINT AT the 2a backstop (the
// `status --validate` self-oracle lint / terminal-PASSED set guard) without
// re-declaring that guard as this entity's own code.
//
// This is the legitimate doc-as-deliverable case: the deliverable is the FO's own
// loaded contract text, so a presence/banned-relationship check over the real file
// is proof at the claim's own level. The BEHAVIORAL teeth — a self-oracle AC
// mechanically cannot reach terminal PASSED — live in entity 2a
// (require-external-proof-guard), explicitly out of scope here; this test asserts
// the contract points at that guard, and the banned half asserts this entity does
// NOT re-file the guard's code.
//
// Failing-first: the pre-edit paragraph (FO core un-hardened original) says only
// "confirm each AC has at least one evidence citation" — it carries none of the
// external-proof clause, the can-never-fail refusal, or the pure-decision clause.
func TestACCrossCheckRequiresExternalProof(t *testing.T) {
	para := acCrossCheckParagraph(t)

	// The clause must not use the banned insider-jargon token (locked separately
	// by working_principles_test.go); these markers are deliberately plain.
	requiredClauses := map[string]string{
		"external-proof requirement (evidence comes from outside the entity body)":        "OUTSIDE the entity body",
		"names the concrete external proof forms (test / command output / on-disk state)": "test, a command's output or exit code, or resulting on-disk state",
		"self-oracle refusal — a prose-only self-proof can never fail":                    "can never fail",
		"pure-decision entity must not advance to terminal PASSED":                        "terminal PASSED",
		"routes a pure-decision entity to the roadmap":                                    "roadmap",
	}
	for clause, marker := range requiredClauses {
		if !strings.Contains(para, marker) {
			t.Errorf("AC-coverage cross-check missing the %q clause (expected marker %q)", clause, marker)
		}
	}

	// AC-1 backstop pointer: the paragraph must point at 2a's gates by name
	// (the self-oracle lint and the terminal-PASSED set guard) so the contract
	// wording references the behavioral backstop rather than restating it.
	for _, pointer := range []string{
		"status --validate",
		"set guard",
	} {
		if !strings.Contains(para, pointer) {
			t.Errorf("AC-coverage cross-check missing the 2a backstop pointer %q — the wording must point at the code guard, not re-assert it", pointer)
		}
	}

	// Banned half: this entity ships WORDING only. It must not re-declare 2a's
	// guard as its own deliverable. The clause points at the guard ("when
	// present", "backstop") rather than claiming this edit implements it. A phrase
	// asserting this edit itself enforces the guarantee would be double-filing.
	bannedSelfClaims := []string{
		"this cross-check enforces",
		"this clause guarantees",
		"this edit implements the guard",
	}
	lower := strings.ToLower(para)
	for _, banned := range bannedSelfClaims {
		if strings.Contains(lower, strings.ToLower(banned)) {
			t.Errorf("AC-coverage cross-check claims to enforce the guard itself (%q) — the behavioral guard is 2a's; this edit ships only the FO-facing wording that points at it", banned)
		}
	}
}

// TestDetachedAdversarialAuditSection locks AC-2: the FO contract carries a
// concrete `## Detached Adversarial Audit` section pinning the three axes the
// captain named — WHEN it triggers, WHAT it produces, HOW it is recorded.
//
// Trigger surfaces (front-door / status / contract / CI), the detached-and-
// read-only requirement, and the Feedback-Cycles recording/routing path are the
// load-bearing invariants. The check is a presence test over the real section
// (a property-of-the-text check at the claim's own level — the deliverable is the
// contract an FO loads, not a runtime behavior).
//
// Failing-first is trivial: the section does not exist in the pre-edit file, so
// sectionAfter returns "" and every assertion fails.
func TestDetachedAdversarialAuditSection(t *testing.T) {
	section := sectionAfter(foSharedCore(t), "## Detached Adversarial Audit")
	if strings.TrimSpace(section) == "" {
		t.Fatalf("FO contract has no `## Detached Adversarial Audit` section")
	}

	// (a) WHEN — the four high-stakes trigger surfaces.
	triggerSurfaces := map[string]string{
		"front-door launcher":          "front-door",
		"status mutation/guard":        "status",
		"shipped contract/scaffolding": "contract",
		"CI/release machinery":         "CI",
	}
	for surface, marker := range triggerSurfaces {
		if !strings.Contains(section, marker) {
			t.Errorf("Detached Adversarial Audit section missing the %q trigger surface (expected marker %q)", surface, marker)
		}
	}

	// (b) HOW it works — detached AND read-only, on a separate checkout, refuting
	// the validation rather than re-implementing it.
	mechanismMarkers := map[string]string{
		"detached checkout requirement":                               "detached",
		"read-only requirement":                                       "read-only",
		"refutation (not re-implementation), via an adversarial edit": "adversarial edit",
	}
	for mechanism, marker := range mechanismMarkers {
		if !strings.Contains(section, marker) {
			t.Errorf("Detached Adversarial Audit section missing the %q (expected marker %q)", mechanism, marker)
		}
	}

	// (c) HOW it is recorded/routed — Feedback Cycles into the prior stage.
	for _, marker := range []string{
		"Feedback",
		"Material",
	} {
		if !strings.Contains(section, marker) {
			t.Errorf("Detached Adversarial Audit section missing the recording/routing marker %q", marker)
		}
	}
}

// archivedBinaryAbsentFeedbackCycles returns the Feedback Cycles text of the
// archived #262 entity when the state checkout is reachable from the code repo,
// and ("", false) when it is not. The state checkout is an orphan branch that is
// NOT present in a clean code-only worktree, so AC-3's cross-check against the
// real prior artifact is exercised when the artifact is present and skipped
// (never failed) when it is absent — honoring the no-hidden-machine-dependencies
// rule for a clean-room `go test`.
func archivedBinaryAbsentFeedbackCycles(t *testing.T) (string, bool) {
	t.Helper()
	p := filepath.Join(repoRoot(t), "docs", "dev", ".spacedock-state", "_archive", "binary-absent-fo-bootstrap", "index.md")
	b, err := os.ReadFile(p)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// TestAuditGroundingCitesRealCatch locks AC-3: the audit formalization is grounded
// in #262's real catch, not invented. The audit section (or the entity-body
// grounding note) must cite #262's two concrete test-strength mechanisms — M1 (the
// `strings.Count(...) > 0` skip-on-zero hole) and M2 (the bare `strings.Contains`
// disclaimer hole). These are checkable against the archived
// binary-absent-fo-bootstrap Feedback Cycles record (external to THIS entity), so
// the citation can fail if the named mechanism does not match the prior artifact.
//
// The grounding lives in the FO contract section (an agent-read surface). The
// cross-check against the archived record runs only when that record is reachable
// — when it is, the test confirms the cited mechanism strings actually appear in
// #262's record, so a fabricated citation fails against the real text.
func TestAuditGroundingCitesRealCatch(t *testing.T) {
	section := sectionAfter(foSharedCore(t), "## Detached Adversarial Audit")
	if strings.TrimSpace(section) == "" {
		t.Fatalf("FO contract has no `## Detached Adversarial Audit` section to ground")
	}

	// M1 and M2 mechanism citations, by their concrete mechanism, not by label.
	mechanisms := map[string]string{
		"M1 — strings.Count skip-on-zero hole":       "strings.Count",
		"M1 — skip-on-zero phrasing":                 "skip-on-zero",
		"M2 — bare strings.Contains disclaimer hole": "strings.Contains",
		"M2 — disclaimer phrasing":                   "disclaimer",
		"#262 reference":                             "#262",
	}
	for mechanism, marker := range mechanisms {
		if !strings.Contains(section, marker) {
			t.Errorf("audit grounding missing the %q citation (expected mechanism marker %q)", mechanism, marker)
		}
	}

	// Cross-check against the real prior artifact when reachable. The archived
	// #262 Feedback Cycles must itself record the strings.Count and
	// strings.Contains mechanisms the contract cites — proving the citation
	// matches the external record, not an invented catch.
	record, ok := archivedBinaryAbsentFeedbackCycles(t)
	if !ok {
		t.Log("archived #262 record not reachable (state checkout absent); skipping cross-check against prior artifact")
		return
	}
	for _, m := range []string{"strings.Count", "strings.Contains"} {
		if !strings.Contains(record, m) {
			t.Errorf("cited mechanism %q is not present in the archived #262 record — the grounding citation does not match the real prior artifact", m)
		}
	}
}
