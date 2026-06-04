// ABOUTME: AC-1/AC-2 oracles for the using-claude-team extraction — the four
// ABOUTME: generic team-lifecycle blocks live in the skill, not the FO runtime, with a clean boundary.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// foRuntimeLineBaseline is the pre-change line count of
// claude-first-officer-runtime.md (ground-truthed 2026-06-04, 305 lines). The
// four extracted blocks span ~90 lines; AC-1(d) asserts a material drop.
const foRuntimeLineBaseline = 305

// foRuntimeLineDropFloor is the AC-1(d) minimum line removal. The four blocks
// span ~90 lines (Team Creation generic ~22, Awaiting Completion ~29, Terminal
// Teardown ~12, Degraded Mode ~31); a ≥70-line drop makes "drops materially"
// checkable rather than vibes.
const foRuntimeLineDropFloor = 70

// genericBlockFingerprints uniquely identifies each of the four extracted
// generic team-lifecycle blocks. Each literal is unique-1 in the pre-change
// claude-first-officer-runtime.md (verified during ideation). AC-1(a) asserts
// presence in the skill; AC-1(b) asserts absence from the FO runtime.
var genericBlockFingerprints = map[string]string{
	"Team Creation":          "TeamCreate failure recovery (priority-ordered ladder)",
	"Awaiting Completion":    "A new `system init` entry in the stream is NOT a completion signal",
	"Terminal Team Teardown": "TERMINAL_TEARDOWN_BOUNDED",
	"Degraded Mode":          "Cooperative Shutdown Sweep",
}

// usingClaudeTeamSkill reads the new skill body under test.
func usingClaudeTeamSkill(t *testing.T) string {
	t.Helper()
	p := filepath.Join(skillsRoot(t), "using-claude-team", "SKILL.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read using-claude-team SKILL.md: %v", err)
	}
	return string(b)
}

// foRuntime reads the FO Claude runtime adapter under test.
func foRuntime(t *testing.T) string {
	t.Helper()
	return vendoredSkillFiles(t)["first-officer/references/claude-first-officer-runtime.md"]
}

// TestGenericBlocksPresentInSkill locks AC-1(a): each of the four generic
// team-lifecycle blocks is present in skills/using-claude-team/SKILL.md, keyed
// by its unique fingerprint literal.
func TestGenericBlocksPresentInSkill(t *testing.T) {
	skill := usingClaudeTeamSkill(t)
	for block, fp := range genericBlockFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("using-claude-team SKILL.md missing %s block fingerprint %q", block, fp)
		}
	}
}

// TestGenericBlocksAbsentFromFORuntime locks AC-1(b): the four generic
// fingerprints are NO LONGER present in claude-first-officer-runtime.md — the
// blocks moved, not duplicated. Whole-file (NOT region-scoped): region-scoping a
// generic-absence check would false-pass content that moved elsewhere in the
// file. Negative-proof: re-inlining any block re-introduces its fingerprint and
// flips this RED.
func TestGenericBlocksAbsentFromFORuntime(t *testing.T) {
	fo := foRuntime(t)
	for block, fp := range genericBlockFingerprints {
		if strings.Contains(fo, fp) {
			t.Errorf("claude-first-officer-runtime.md still inlines %s block fingerprint %q (moved, not duplicated)", block, fp)
		}
	}
}

// TestFORuntimeInvokesSkill locks AC-1(c): the FO runtime's `## Team Creation`
// section invokes the skill via Skill(...) and does NOT use the disproven
// cross-skill @-include. Region-scoped to `## Team Creation` (the positive
// Skill()-present / @-absent assertions only — the region legitimately retains
// the standing-teammate subsections). The Skill(...) literal is the integration
// seam; the @-form is the spike-disproven mechanism.
func TestFORuntimeInvokesSkill(t *testing.T) {
	fo := foRuntime(t)
	region := sectionAfter(fo, "## Team Creation")
	if region == "" {
		t.Fatal("claude-first-officer-runtime.md has no `## Team Creation` section")
	}
	if !strings.Contains(region, `Skill(skill="spacedock:using-claude-team")`) {
		t.Errorf("`## Team Creation` section does not invoke Skill(skill=\"spacedock:using-claude-team\")")
	}
	for _, banned := range []string{"@../using-claude-team", "@using-claude-team"} {
		if strings.Contains(region, banned) {
			t.Errorf("`## Team Creation` section uses the disproven cross-skill @-include %q", banned)
		}
	}
}

// TestFORuntimeDroppedMaterially locks AC-1(d): claude-first-officer-runtime.md
// line count dropped by at least foRuntimeLineDropFloor vs the pre-change
// baseline. Secondary signal to the fingerprint-absence teeth above — the floor
// catches a no-op edit, the fingerprints catch a duplicate-not-moved edit.
func TestFORuntimeDroppedMaterially(t *testing.T) {
	fo := foRuntime(t)
	lines := strings.Count(fo, "\n")
	if fo != "" && !strings.HasSuffix(fo, "\n") {
		lines++
	}
	dropped := foRuntimeLineBaseline - lines
	if dropped < foRuntimeLineDropFloor {
		t.Errorf("claude-first-officer-runtime.md dropped only %d lines (baseline %d, now %d); want ≥ %d removed",
			dropped, foRuntimeLineBaseline, lines, foRuntimeLineDropFloor)
	}
}

// skillLeakageLiterals are spacedock-specific tokens the generic team-harness
// skill must NOT name. The qualified `spacedock dispatch` covers every
// dispatch-helper leak (build / reconcile / context-budget) in one. Mirrors the
// sibling devLeakageLiterals table. Bare `reconcile` is deliberately NOT here —
// it appears in legitimate generic event-loop/backstop prose and would
// false-fire.
var skillLeakageLiterals = []string{
	"spacedock dispatch",
	"spacedock status",
	"feedback-to",
	"context-budget",
}

// TestSkillFreeOfSpacedockTokens locks AC-2 (absence half): the generic skill is
// free of spacedock-specific tokens. Negative-proof: a `spacedock dispatch`
// token leaking into the skill reds this.
func TestSkillFreeOfSpacedockTokens(t *testing.T) {
	skill := usingClaudeTeamSkill(t)
	for _, banned := range skillLeakageLiterals {
		if strings.Contains(skill, banned) {
			t.Errorf("using-claude-team SKILL.md leaks spacedock-specific token %q (must stay team-harness-generic)", banned)
		}
	}
}

// TestSpacedockDecisionsStayInFORuntime locks AC-2 (presence half): the
// spacedock decision points REMAIN in the FO contract. The positive anchors are
// the QUALIFIED dispatch-helper calls — `spacedock dispatch build`,
// `spacedock dispatch context-budget`, and the fully-qualified
// `spacedock dispatch reconcile` (NOT bare `reconcile`, which is ×4 incl.
// generic prose). Negative-proof: the qualified build/context-budget/reconcile
// call wrongly moved out of the FO contract reds this.
func TestSpacedockDecisionsStayInFORuntime(t *testing.T) {
	fo := foRuntime(t)
	for _, anchor := range []string{
		"spacedock dispatch build",
		"spacedock dispatch context-budget",
		"spacedock dispatch reconcile",
	} {
		if !strings.Contains(fo, anchor) {
			t.Errorf("claude-first-officer-runtime.md no longer contains spacedock decision anchor %q (must stay in the FO contract)", anchor)
		}
	}
}
