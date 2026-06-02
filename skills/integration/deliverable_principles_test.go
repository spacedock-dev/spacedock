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
// self-proof criterion (one whose only proof is review of the entity's own prose,
// which can never fail), and refuses advancing a pure-decision entity to a
// terminal PASSED verdict. It must also POINT AT the 2a backstop (the
// `status --validate` self-proof lint / terminal-PASSED set guard) without
// re-declaring that guard as this entity's own code.
//
// This is the legitimate doc-as-deliverable case: the deliverable is the FO's own
// loaded contract text, so a relationship/polarity check over the real file is
// proof at the claim's own level. The BEHAVIORAL teeth — a self-proof AC
// mechanically cannot reach terminal PASSED — live in entity 2a
// (require-external-proof-guard), explicitly out of scope here; this test asserts
// the contract points at that guard, and the relationship half asserts this entity
// does NOT re-file the guard's code by claiming the prose is self-sufficient.
//
// Pins RELATIONSHIP and POLARITY, not token presence (the #262 M1/M2 lesson this
// entity teaches against, applied to its own tests): each required clause is
// asserted together with a ban on its negated/disclaimer inversion, so a
// meaning-flipping rewrite that keeps the token fails.
//
// Failing-first: the pre-edit paragraph (FO core un-hardened original) says only
// "confirm each AC has at least one evidence citation" — it carries none of the
// external-proof clause, the can-never-fail refusal, or the pure-decision clause.
func TestACCrossCheckRequiresExternalProof(t *testing.T) {
	para := acCrossCheckParagraph(t)
	lower := strings.ToLower(para)

	// Each required clause: the affirmative marker MUST be present, AND a polarity
	// inversion that flips its meaning while keeping the token MUST be absent. The
	// inversion phrases are the disclaimer/negation forms the #262-M2 lesson warns
	// a bare-substring check would green-light.
	requiredClauses := []struct {
		clause     string
		marker     string   // affirmative phrasing that must be present
		inversions []string // negated/disclaimer forms that must NOT appear (lowercased)
	}{
		{
			clause:     "external-proof requirement (evidence comes from outside the entity body)",
			marker:     "OUTSIDE the entity body",
			inversions: []string{"need not come from", "no external", "inside the entity body is sufficient", "self-review is acceptable"},
		},
		{
			clause:     "names the concrete external proof forms (test / command output / on-disk state)",
			marker:     "test, a command's output or exit code, or resulting on-disk state",
			inversions: nil,
		},
		{
			clause:     "self-proof refusal — a prose-only self-proof can never fail, so does not satisfy the check",
			marker:     "can never fail",
			inversions: []string{"no longer applies", "is acceptable", "is now acceptable", "historically", "but that concern", "review of the entity's own prose is sufficient"},
		},
		{
			clause:     "pure-decision entity must not advance to terminal PASSED",
			marker:     "terminal PASSED",
			inversions: []string{"may advance it to a terminal passed", "is allowed to reach terminal passed", "can advance to a terminal passed"},
		},
		{
			clause:     "routes a pure-decision entity to the roadmap",
			marker:     "roadmap",
			inversions: nil,
		},
	}
	for _, rc := range requiredClauses {
		if !strings.Contains(para, rc.marker) {
			t.Errorf("AC-coverage cross-check missing the %q clause (expected marker %q)", rc.clause, rc.marker)
		}
		for _, inv := range rc.inversions {
			if strings.Contains(lower, inv) {
				t.Errorf("AC-coverage cross-check carries an inversion of the %q clause (%q present) — the rule's polarity is flipped while keeping the token", rc.clause, inv)
			}
		}
	}

	// Polarity of the self-proof refusal: "can never fail" must be FOLLOWED by the
	// consequence "does not satisfy" within the same clause, so the rewrite
	// "can never fail, but that concern no longer applies; self-review is
	// acceptable" (which drops/contradicts the consequence) fails. This pins the
	// rule's relationship, not the presence of the phrase alone.
	if i := strings.Index(para, "can never fail"); i >= 0 {
		tail := para[i:]
		if !strings.Contains(tail, "does not satisfy") {
			t.Errorf("AC-coverage cross-check's `can never fail` clause is not followed by its consequence (`does not satisfy this cross-check`) — the refusal's polarity is broken")
		}
	}

	// Backstop pointer (relationship OUT to 2a): the paragraph must point at 2a's
	// gates by name (the self-proof lint and the terminal-PASSED set guard) AND
	// frame them as the real assurance / backstop, so the wording references the
	// behavioral backstop rather than restating or replacing it.
	for _, pointer := range []string{
		"status --validate",
		"set guard",
		"backstop",
	} {
		if !strings.Contains(para, pointer) {
			t.Errorf("AC-coverage cross-check missing the 2a backstop pointer %q — the wording must point OUT at the code guard, not re-assert it", pointer)
		}
	}

	// Relationship half (replaces the old 3-string denylist): this entity ships
	// WORDING only. The clause must NOT claim the cross-check is itself
	// self-sufficient / the binding guarantee with no external code gate needed —
	// that would re-file 2a's guard as this entity's own deliverable (double-file)
	// and contradict the "real assurance is that gate" subordination. Self-
	// sufficiency claims vary in wording; this bans the relationship pattern
	// (some negation of "external gate" / some claim "this … is the guarantee").
	//
	// Residual limit (honestly recorded): a wholly novel self-sufficiency phrasing
	// matching NONE of these patterns could still slip a prose check. That ceiling
	// is the doc-as-deliverable limit this entity documents — the BEHAVIORAL teeth
	// (a self-proof AC mechanically cannot reach terminal PASSED) are 2a's
	// `status --validate` guard, explicitly out of scope here. The pointer asserts
	// above guarantee the prose subordinates to that gate rather than replacing it.
	selfSufficiencyPatterns := []string{
		"no external code gate",
		"no external gate",
		"without an external",
		"is itself the binding guarantee",
		"is itself the guarantee",
		"this cross-check is sufficient on its own",
		"this clause is sufficient on its own",
		"self-sufficient",
	}
	for _, p := range selfSufficiencyPatterns {
		if strings.Contains(lower, p) {
			t.Errorf("AC-coverage cross-check claims self-sufficiency (%q) — the behavioral guard is 2a's; this edit ships only the FO-facing wording that points OUT at it", p)
		}
	}
}

// auditSection isolates the `## Detached Adversarial Audit` FO-contract section.
func auditSection(t *testing.T) string {
	t.Helper()
	section := sectionAfter(foSharedCore(t), "## Detached Adversarial Audit")
	if strings.TrimSpace(section) == "" {
		t.Fatalf("FO contract has no `## Detached Adversarial Audit` section")
	}
	return section
}

// auditTriggerSentence returns the opening trigger-surface sentence of the audit
// section — the WHEN axis — scoped from the section start up to the first bullet
// (`\n- `). The four trigger surfaces are named in this sentence; scoping here
// stops a surface marker from "passing" by leaking into unrelated later prose
// ("front-door change", "contract_gate_test.go"), the #262-M2 bare-substring leak
// the validator flagged.
func auditTriggerSentence(t *testing.T, section string) string {
	t.Helper()
	if i := strings.Index(section, "\n- "); i >= 0 {
		return section[:i]
	}
	return section
}

// TestDetachedAdversarialAuditSection locks AC-2: the FO contract carries a
// concrete `## Detached Adversarial Audit` section pinning the three axes the
// captain named — WHEN it triggers, WHAT it produces, HOW it is recorded.
//
// Pins STRUCTURE, not tokens-anywhere (the #262 lesson applied to this entity's
// own test). The three load-bearing instructions must appear as DISTINCT bullets
// carrying their verbs, and an inversion ("intentionally left blank; do not
// perform any audit") is rejected — so replacing the whole section with a
// token-salad line fails. Trigger surfaces are scoped to the WHEN sentence so a
// surface removed from the trigger line fails even if the token appears elsewhere.
//
// Failing-first is trivial: the section does not exist in the pre-edit file, so
// auditSection's Fatalf fires.
func TestDetachedAdversarialAuditSection(t *testing.T) {
	section := auditSection(t)
	lower := strings.ToLower(section)

	// Inversion guard: a "left blank" / "do not perform an audit" rewrite that
	// keeps the heading and sprays tokens must fail. These phrases negate the
	// section's instruction outright.
	for _, inv := range []string{
		"left blank",
		"do not perform",
		"no audit is required",
		"no audit needed",
	} {
		if strings.Contains(lower, inv) {
			t.Errorf("Detached Adversarial Audit section is inverted to a no-op (%q present) — the section must instruct that the audit IS performed", inv)
		}
	}

	// (a) WHEN — the four high-stakes trigger surfaces, asserted on the trigger
	// SENTENCE (not section-wide) with trigger-line-unique phrasing so a surface
	// removed from the WHEN axis fails even though its token recurs later.
	trigger := auditTriggerSentence(t, section)
	triggerSurfaces := map[string]string{
		"front-door launcher":          "front-door launcher",
		"status mutation/guard paths":  "`status` mutation/guard",
		"shipped contract/scaffolding": "shipped contract/scaffolding",
		"CI/release machinery":         "CI/release machinery",
	}
	for surface, marker := range triggerSurfaces {
		if !strings.Contains(trigger, marker) {
			t.Errorf("Detached Adversarial Audit WHEN-trigger sentence missing the %q surface (expected trigger-line phrasing %q)", surface, marker)
		}
	}

	// (b)+(c) the three load-bearing instruction bullets, each a DISTINCT bullet
	// (opens with `- **`) carrying its verb/clause. Asserting the bullet opener +
	// the distinctive instruction text pins the structure, so a token-salad
	// replacement that lacks the bullet structure fails.
	requiredBullets := []struct {
		axis   string
		opener string // distinctive bullet heading
		clause string // load-bearing instruction within that bullet
	}{
		{
			axis:   "detached + read-only (HOW the audit runs)",
			opener: "- **Detached and read-only.**",
			clause: "adversarial edit",
		},
		{
			axis:   "two-tier finding output (WHAT it produces)",
			opener: "- **What it produces.**",
			clause: "Material",
		},
		{
			axis:   "Feedback-Cycles recording/routing (HOW it is recorded)",
			opener: "- **How it is recorded and routed.**",
			clause: "Feedback Rejection Flow",
		},
	}
	for _, rb := range requiredBullets {
		bi := strings.Index(section, rb.opener)
		if bi < 0 {
			t.Errorf("Detached Adversarial Audit section missing the %q bullet (expected distinct bullet %q)", rb.axis, rb.opener)
			continue
		}
		// Scope to this bullet: from its opener to the next bullet or blank line.
		rest := section[bi+len(rb.opener):]
		end := len(rest)
		if j := strings.Index(rest, "\n- "); j >= 0 && j < end {
			end = j
		}
		if j := strings.Index(rest, "\n\n"); j >= 0 && j < end {
			end = j
		}
		bullet := rest[:end]
		if !strings.Contains(bullet, rb.clause) {
			t.Errorf("Detached Adversarial Audit %q bullet missing its load-bearing instruction %q", rb.axis, rb.clause)
		}
	}

	// Detached + read-only must both be required, and the audit is refutation not
	// re-implementation. These are the load-bearing HOW invariants.
	for _, marker := range []string{"DETACHED checkout", "read-only", "REFUTE"} {
		if !strings.Contains(section, marker) {
			t.Errorf("Detached Adversarial Audit section missing the %q requirement", marker)
		}
	}
}

// archived262Path is the in-repo testdata fixture vendoring the #262
// (binary-absent-fo-bootstrap) Feedback Cycles fragment. The live record lives in
// the state-checkout orphan branch that is NOT present in a code-only worktree, so
// the fixture is vendored here to make AC-3's cross-check ALWAYS run in CI (no
// skip-on-absence hole — the #262-M1 lesson applied to this entity's own test).
func archived262FeedbackCycles(t *testing.T) string {
	t.Helper()
	p := filepath.Join("testdata", "binary-absent-fo-bootstrap-feedback-cycles.md")
	b, err := os.ReadFile(p)
	if err != nil {
		// Hard-fail, never skip: the fixture is committed in-repo and must be
		// present. A missing fixture is a test-integrity failure, not a reason to
		// silently pass.
		t.Fatalf("read vendored #262 fixture %s: %v — fixture must be present in-repo so the cross-check always runs", p, err)
	}
	return string(b)
}

// TestAuditGroundingCitesRealCatch locks AC-3: the audit formalization is grounded
// in #262's real catch, not invented. The audit section must cite #262's two
// concrete test-strength mechanisms — M1 (the `strings.Count(...) > 0` skip-on-zero
// hole) and M2 (the bare `strings.Contains` disclaimer hole) — and the citation
// must match the archived #262 record (vendored as in-repo testdata).
//
// Pins RELATIONSHIP/POLARITY and runs ALWAYS (closes the two #262 holes in this
// entity's own test):
//   - M1 (no skip-on-absence): the #262 record is a committed in-repo fixture, so
//     the cross-check runs in CI / a code-only worktree, never skips.
//   - M2 (reject inversion): a citation that negates the grounding ("unrelated to
//     #262, which found nothing; ignore strings.Count…") is rejected, not just a
//     token-presence pass.
func TestAuditGroundingCitesRealCatch(t *testing.T) {
	section := auditSection(t)
	lower := strings.ToLower(section)

	// Citation-inversion guard (M2 lesson): a negated grounding that keeps the
	// tokens must fail. These phrases flip the citation's meaning.
	for _, inv := range []string{
		"unrelated to #262",
		"found nothing",
		"ignore any mention of strings.count",
		"ignore strings.count",
		"not grounded in",
		"no real catch",
	} {
		if strings.Contains(lower, inv) {
			t.Errorf("audit grounding is inverted (%q present) — the citation must affirm #262's real catch, not negate it", inv)
		}
	}

	// Affirmative M1/M2 mechanism citations, by concrete mechanism, not by label.
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

	// Cross-check against the real prior artifact, ALWAYS (no skip-on-absence). The
	// archived #262 Feedback Cycles fixture must itself record the strings.Count
	// (M1) and strings.Contains (M2) mechanisms the contract cites — proving the
	// citation matches the external record, not an invented catch. Drift between
	// the contract's citation and the record fails here.
	record := archived262FeedbackCycles(t)
	for _, m := range []string{"strings.Count", "strings.Contains"} {
		if !strings.Contains(record, m) {
			t.Errorf("cited mechanism %q is not present in the archived #262 record — the grounding citation does not match the real prior artifact", m)
		}
	}
}
