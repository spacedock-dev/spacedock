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

// TestFeedbackProcedurePresentInSkill locks AC-1(a): the moved procedure
// fingerprints are present in skills/feedback-rejection-flow/SKILL.md.
func TestFeedbackProcedurePresentInSkill(t *testing.T) {
	skill := feedbackRejectionFlowSkill(t)
	for name, fp := range feedbackProcedureFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("feedback-rejection-flow SKILL.md missing %s fingerprint %q", name, fp)
		}
	}
}

// TestFeedbackFaithfulnessClausesPresentInSkill locks AC-2 (faithfulness): the
// two mis-route-on-loss clauses — the Codex `send_input` non-completion caveat
// and the `feedback-to` target-read clause — are present in the skill body.
func TestFeedbackFaithfulnessClausesPresentInSkill(t *testing.T) {
	skill := feedbackRejectionFlowSkill(t)
	for name, fp := range feedbackFaithfulnessFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("feedback-rejection-flow SKILL.md missing %s faithfulness clause %q", name, fp)
		}
	}
}

// TestFeedbackProcedureAbsentFromFOCore locks AC-1(b): the moved procedure
// fingerprints (and the faithfulness clauses, which moved with the body) are NO
// LONGER present in first-officer-shared-core.md — moved, not duplicated. Whole-
// file (NOT region-scoped): region-scoping an absence check would false-pass
// content that moved elsewhere in the file. Negative-proof: re-inlining the
// procedure re-introduces a fingerprint and flips this RED.
func TestFeedbackProcedureAbsentFromFOCore(t *testing.T) {
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

// TestFOCoreInvokesFeedbackRejectionSkill locks AC-1(c): the FO core invokes the
// skill via Skill(...) at the rejection-detection point and uses NO cross-skill
// @-include toward feedback-rejection-flow ANYWHERE in the file. The positive
// Skill(...) check is region-scoped to `## Completion and Gates` (the seam lives
// at the detection point); the @-include ban is WHOLE-FILE so a stale include
// re-introduced in any other section (e.g. `## Merge and Cleanup`) is caught, not
// just one in the detection region. The Skill(...) literal is the integration
// seam; any `@`-token resolving toward feedback-rejection-flow is the disproven
// mechanism.
func TestFOCoreInvokesFeedbackRejectionSkill(t *testing.T) {
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

// TestAlwaysOnMachineryRetainedInFOCore locks AC-1(d): the referenced always-on
// machinery did NOT move with the procedure. The FO Write Scope `### Feedback
// Cycles` write-scope bullet and the reuse-conditions/budget-probe block stay in
// the FO core. Negative-proof: deleting either anchor reds this.
func TestAlwaysOnMachineryRetainedInFOCore(t *testing.T) {
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

// TestClaudeBareModeSeamStaysConsistent locks AC-2 (seam): the Claude adapter's
// `## Feedback Rejection Flow (bare mode)` seam stays — still present, still the
// sequential-dispatch sentence and the keep-reviewer-alive sentence. The seam is
// a Claude-runtime execution mode, NOT moved into the runtime-neutral skill.
// Negative-proof: removing the seam (or either sentence) reds this.
func TestClaudeBareModeSeamStaysConsistent(t *testing.T) {
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

// TestFeedbackRejectionSkillNameMatchesSeam locks AC-2 (hardened seam): the
// frontmatter `name:` VALUE equals `feedback-rejection-flow` — the directory name
// AND the `Skill(skill="spacedock:feedback-rejection-flow")` invocation seam.
// Token-presence alone (skill_surface_test.go) would pass a renamed skill the
// seam no longer reaches; binding the value to the seam target catches that
// drift. Negative-proof: a bogus name value reds this.
func TestFeedbackRejectionSkillNameMatchesSeam(t *testing.T) {
	name, ok := feedbackRejectionFrontmatterValue(t, "name")
	if !ok {
		t.Fatal("feedback-rejection-flow SKILL.md frontmatter has no name field")
	}
	if name != "feedback-rejection-flow" {
		t.Errorf("feedback-rejection-flow SKILL.md frontmatter name is %q, want %q (the directory name and the Skill(skill=\"spacedock:feedback-rejection-flow\") seam)", name, "feedback-rejection-flow")
	}
}

// TestFeedbackRejectionSkillIsFOInternal locks AC-2 (hardened seam): the
// frontmatter carries `user-invocable: false` — the skill is FO-internal (loaded
// mid-run via Skill()), not a captain-facing user skill. Negative-proof: flipping
// to `true` reds this.
func TestFeedbackRejectionSkillIsFOInternal(t *testing.T) {
	v, ok := feedbackRejectionFrontmatterValue(t, "user-invocable")
	if !ok {
		t.Fatal("feedback-rejection-flow SKILL.md frontmatter has no user-invocable field")
	}
	if v != "false" {
		t.Errorf("feedback-rejection-flow SKILL.md frontmatter user-invocable is %q, want \"false\" (the skill is FO-internal)", v)
	}
}
