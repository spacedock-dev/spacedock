// ABOUTME: AC-1/AC-2 oracles for the feedback-rejection-flow extraction — the
// ABOUTME: 7-step rejection procedure lives in the skill, not the FO core, with a Skill() seam; always-on machinery stays.
package integration

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// feedbackProcedureFingerprints uniquely identifies the moved feedback-rejection
// procedure. Each literal is unique-1 in the pre-change first-officer-shared-core.md
// (ground-truthed 2026-06-04). AC-1(a) asserts presence in the skill; AC-1(b)
// asserts absence from the FO core. The four cover the load-bearing steps — the
// 3-cycle escalation, re-run-reviewer, re-enter-gate, and the FO ownership of
// `### Feedback Cycles`.
var feedbackProcedureFingerprints = map[string]string{
	"cycle-3-escalation":     "On cycle 3, escalate to the human instead of another round",
	"re-run-reviewer":        "Re-run the reviewer after fixes",
	"re-enter-gate-flow":     "Re-enter the normal gate flow with the updated result",
	"fo-owns-feedback-cycle": "The FO owns `### Feedback Cycles`",
}

// feedbackFaithfulnessFingerprints are the AC-2 clauses whose loss would silently
// mis-route a rejection: the Codex `send_input` non-completion caveat and the
// `feedback-to` target-read clause. Asserted present in the skill body.
var feedbackFaithfulnessFingerprints = map[string]string{
	"send-input-non-completion": "do not treat the immediate `send_input` response as the new completion result",
	"feedback-to-target-read":   "the stage that receives the fix request, not the reviewer",
}

// feedbackRejectionFlowSkill reads the new skill body under test.
func feedbackRejectionFlowSkill(t *testing.T) string {
	t.Helper()
	p := filepath.Join(skillsRoot(t), "feedback-rejection-flow", "SKILL.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read feedback-rejection-flow SKILL.md: %v", err)
	}
	return string(b)
}

// TestFeedbackProcedurePresentInSkill is a non-AC text-consistency lint: the
// moved procedure fingerprints are present in
// skills/feedback-rejection-flow/SKILL.md — the prose MOVED here. Per the proof
// policy this presence check does NOT prove the FO follows the procedure; an
// inverted skill body keeps every fingerprint. The behavior — the FO observes a
// REJECTED report and routes the concrete finding back through implementation — is
// proven by the live rejection-flow scenario (runClaudeRejectionFlowScenario /
// runCodexRejectionFlowScenario, asserted by assertRejectionFlow) and its offline
// mutation control TestRejectionFlowNegativeMissingRoute.
func TestFeedbackProcedurePresentInSkill(t *testing.T) {
	markNonAC(t, "live rejection-flow scenario (assertRejectionFlow) + TestRejectionFlowNegativeMissingRoute")
	skill := feedbackRejectionFlowSkill(t)
	for name, fp := range feedbackProcedureFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("feedback-rejection-flow SKILL.md missing %s fingerprint %q", name, fp)
		}
	}
}

// TestFeedbackFaithfulnessClausesPresentInSkill is a non-AC text-consistency
// lint: the two mis-route-on-loss clauses — the Codex `send_input` non-completion
// caveat and the `feedback-to` target-read clause — are present in the skill body.
// This is text authoring, not behavioral proof; the behavior that the FO routes
// the fix to the feedback-to target (not the reviewer) is proven by the live
// rejection-flow scenario, which asserts the entity returns to status:
// implementation with the fix applied.
func TestFeedbackFaithfulnessClausesPresentInSkill(t *testing.T) {
	markNonAC(t, "live rejection-flow scenario (assertRejectionFlow) + TestRejectionFlowNegativeMissingRoute")
	skill := feedbackRejectionFlowSkill(t)
	for name, fp := range feedbackFaithfulnessFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("feedback-rejection-flow SKILL.md missing %s faithfulness clause %q", name, fp)
		}
	}
}

// TestFeedbackProcedureAbsentFromFOCore is a non-AC text-consistency lint (dedup):
// the moved procedure fingerprints (and the faithfulness clauses, which moved with
// the body) are NO LONGER present in first-officer-shared-core.md — moved, not
// duplicated. Whole-file (NOT region-scoped): region-scoping an absence check
// would false-pass content that moved elsewhere. This is a structural dedup
// property, not a behavioral claim; the FO's rejection behavior is proven by the
// live rejection-flow scenario. Re-inlining the procedure re-introduces a
// fingerprint and flips this RED.
func TestFeedbackProcedureAbsentFromFOCore(t *testing.T) {
	markNonAC(t, "dedup lint; behavior via live rejection-flow scenario (assertRejectionFlow)")
	fo := foCore(t)
	for name, fp := range feedbackProcedureFingerprints {
		if strings.Contains(fo, fp) {
			t.Errorf("first-officer-shared-core.md still inlines %s fingerprint %q (moved, not duplicated)", name, fp)
		}
	}
	for name, fp := range feedbackFaithfulnessFingerprints {
		if strings.Contains(fo, fp) {
			t.Errorf("first-officer-shared-core.md still inlines %s faithfulness clause %q (moved, not duplicated)", name, fp)
		}
	}
}

// feedbackAtInclude matches an `@`-prefixed path token resolving toward the
// feedback-rejection-flow skill — `@feedback-rejection-flow`,
// `@../feedback-rejection-flow`, `@./feedback-rejection-flow/SKILL.md`, etc. The
// leading `@` plus any run of relative-path segments (`./`, `../`, bare) ending
// in a `feedback-rejection-flow` path component is the disproven cross-skill
// include the seam must NOT use. Structural (not a two-literal enum), and scanned
// whole-file: a stale `@feedback-rejection-flow` re-introduced in ANY core
// section — not just the detection region — is the disproven mechanism. A bare
// `@\S` whole-file ban would false-fire on a legitimate unrelated
// `@references/...` include the core may later carry; this name-targeted form
// only fires on the feedback-rejection-flow include family.
var feedbackAtInclude = regexp.MustCompile(`@(?:\.{1,2}/)*feedback-rejection-flow\b`)

// TestFOCoreInvokesFeedbackRejectionSkill is a non-AC text-consistency lint: it
// asserts the FO core carries the Skill(skill="spacedock:feedback-rejection-flow")
// invocation literal at the rejection-detection point and no disproven
// cross-skill @-include. Per the proof policy this presence check does NOT prove
// the FO invokes the skill on a rejection: an inverted clause ("NEVER invoke
// feedback-rejection-flow; just wait") keeps the Skill(...) substring and passes
// (verified in ideation). The behavior — the FO routes a REJECTED finding back
// through implementation — is proven only by the live rejection-flow scenario
// (assertRejectionFlow). This lint guards the seam STRING and bans the @-include
// mechanism; it is the text half, not the behavioral proof.
func TestFOCoreInvokesFeedbackRejectionSkill(t *testing.T) {
	markNonAC(t, "live rejection-flow scenario (assertRejectionFlow) + TestRejectionFlowNegativeMissingRoute")
	fo := foCore(t)
	region := sectionAfter(fo, "## Completion and Gates")
	if region == "" {
		t.Fatal("first-officer-shared-core.md has no `## Completion and Gates` section")
	}
	if !strings.Contains(region, `Skill(skill="spacedock:feedback-rejection-flow")`) {
		t.Errorf("`## Completion and Gates` section does not invoke Skill(skill=\"spacedock:feedback-rejection-flow\")")
	}
	if m := feedbackAtInclude.FindString(fo); m != "" {
		t.Errorf("first-officer-shared-core.md uses a disproven cross-skill @-include toward feedback-rejection-flow (token %q)", m)
	}
}

// TestAlwaysOnMachineryRetainedInFOCore is a non-AC text-consistency lint (the
// retention sibling of the dedup checks): the referenced always-on machinery did
// NOT move with the procedure — the FO Write Scope `### Feedback Cycles` bullet
// and the reuse-conditions block stay in the FO core. This is a structural
// retention property, not a behavioral claim; that the FO actually tracks feedback
// cycles is proven by the live rejection-flow scenario. Deleting either anchor
// reds this.
func TestAlwaysOnMachineryRetainedInFOCore(t *testing.T) {
	markNonAC(t, "retention lint; behavior via live rejection-flow scenario (assertRejectionFlow)")
	fo := foCore(t)
	for name, anchor := range map[string]string{
		"feedback-cycles-write-scope": "**`### Feedback Cycles` section**",
		"reuse-conditions":            "Reuse conditions",
	} {
		if !strings.Contains(fo, anchor) {
			t.Errorf("first-officer-shared-core.md no longer contains always-on anchor %s %q (must stay — the procedure references it by name)", name, anchor)
		}
	}
}

// TestClaudeBareModeSeamStaysConsistent is a non-AC text-consistency lint: the
// Claude adapter's `## Feedback Rejection Flow (bare mode)` seam stays present
// with its sequential-dispatch and keep-reviewer-alive sentences — the seam is a
// Claude-runtime execution mode, not moved into the runtime-neutral skill. This
// is text authoring, not behavioral proof; the live rejection-flow scenario
// (Claude runner) exercises the bare-mode path for real. Removing the seam reds
// this.
func TestClaudeBareModeSeamStaysConsistent(t *testing.T) {
	markNonAC(t, "live rejection-flow scenario, Claude runner (assertRejectionFlow)")
	claude := vendoredSkillFiles(t)["first-officer/references/claude-first-officer-runtime.md"]
	for name, fp := range map[string]string{
		"bare-mode-heading":   "## Feedback Rejection Flow (bare mode)",
		"sequential-dispatch": "the feedback rejection flow is sequential",
		"keep-reviewer-alive": "Keep the reviewer alive",
	} {
		if !strings.Contains(claude, fp) {
			t.Errorf("claude-first-officer-runtime.md missing bare-mode seam %s fingerprint %q (the adapter seam must stay consistent)", name, fp)
		}
	}
}

// feedbackRejectionFrontmatterValue returns the trimmed scalar value of a
// top-level `key:` line in the feedback-rejection-flow SKILL.md frontmatter, with
// surrounding quotes stripped. The bool reports whether the key was found.
func feedbackRejectionFrontmatterValue(t *testing.T, key string) (string, bool) {
	t.Helper()
	fm, ok := frontmatter(feedbackRejectionFlowSkill(t))
	if !ok {
		t.Fatal("feedback-rejection-flow SKILL.md has no YAML frontmatter block")
	}
	prefix := key + ":"
	for _, line := range strings.Split(fm, "\n") {
		if strings.HasPrefix(line, prefix) {
			v := strings.TrimSpace(strings.TrimPrefix(line, prefix))
			v = strings.Trim(v, `"'`)
			return v, true
		}
	}
	return "", false
}

// feedbackRejectionSeamLiteral is the seam name the FO contract invokes; the
// re-bound checks read the expected value from the contract, not from the skill
// file under test.
const feedbackRejectionSeamLiteral = "feedback-rejection-flow"

// TestFeedbackRejectionSkillNameMatchesSeam is a code-bound invariant: the skill's
// frontmatter `name:` equals the seam name the FO CONTRACT invokes
// (Skill(skill="spacedock:NAME") in first-officer-shared-core.md). The expected
// value comes from the contract, not the skill file under test — renaming either
// side makes the two diverge and reds this, catching a renamed skill the FO
// invocation no longer reaches.
func TestFeedbackRejectionSkillNameMatchesSeam(t *testing.T) {
	markCodeBoundInvariant(t, "FO contract Skill(skill=\"spacedock:feedback-rejection-flow\") invocation (first-officer-shared-core.md)")
	name, ok := feedbackRejectionFrontmatterValue(t, "name")
	if !ok {
		t.Fatal("feedback-rejection-flow SKILL.md frontmatter has no name field")
	}
	seam := invokedSeamName(foCore(t), feedbackRejectionSeamLiteral)
	if seam == "" {
		t.Fatalf("FO contract does not invoke Skill(skill=\"spacedock:%s\") — the seam the skill name must match is gone", feedbackRejectionSeamLiteral)
	}
	if name != seam {
		t.Errorf("feedback-rejection-flow SKILL.md frontmatter name is %q, but the FO contract invokes the seam %q — a renamed skill the FO invocation no longer reaches", name, seam)
	}
}

// TestFeedbackRejectionSkillIsFOInternal is a code-bound invariant binding the
// skill's `user-invocable` frontmatter to its ROLE: a skill the FO invokes mid-run
// via Skill(skill="spacedock:NAME") is FO-internal and MUST be
// `user-invocable: false`. The expected value is REQUIRED by the contract invoking
// the seam (an independent source), not a free literal; flipping the frontmatter
// to `true` while the FO still invokes the seam reds this.
func TestFeedbackRejectionSkillIsFOInternal(t *testing.T) {
	markCodeBoundInvariant(t, "FO contract Skill(skill=\"spacedock:feedback-rejection-flow\") invocation implies FO-internal")
	if invokedSeamName(foCore(t), feedbackRejectionSeamLiteral) == "" {
		t.Fatalf("FO contract does not invoke the feedback-rejection-flow seam — the FO-internal premise no longer holds")
	}
	v, ok := feedbackRejectionFrontmatterValue(t, "user-invocable")
	if !ok {
		t.Fatal("feedback-rejection-flow SKILL.md frontmatter has no user-invocable field")
	}
	if v != "false" {
		t.Errorf("feedback-rejection-flow SKILL.md frontmatter user-invocable is %q, but the FO contract invokes it as a mid-run seam — an FO-internal skill must be user-invocable: false", v)
	}
}
