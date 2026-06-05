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

// TestGenericBlocksPresentInSkill is a non-AC text-consistency lint: each of the
// four generic team-lifecycle blocks is present in
// skills/using-claude-team/SKILL.md (the blocks MOVED here). Per the proof policy
// this is text authoring, not behavioral proof; the behavior these blocks govern
// (the FO creates a team, dispatches workers, awaits completion, tears down) is
// exercised for real by every team-using live scenario (the gate-guardrail and
// rejection-flow runs both launch a real team via the FO runtime). This lint
// guards against a block fingerprint being dropped.
func TestGenericBlocksPresentInSkill(t *testing.T) {
	markNonAC(t, "live team-using scenarios (gate-guardrail/rejection-flow launch a real team)")
	skill := usingClaudeTeamSkill(t)
	for block, fp := range genericBlockFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("using-claude-team SKILL.md missing %s block fingerprint %q", block, fp)
		}
	}
}

// TestGenericBlocksAbsentFromFORuntime is a non-AC text-consistency lint (dedup):
// the four generic fingerprints are NO LONGER present in
// claude-first-officer-runtime.md — the blocks moved, not duplicated. Whole-file
// (NOT region-scoped) so content that moved elsewhere does not false-pass. This is
// a structural dedup property, not a behavioral claim; the team behavior is proven
// by the live team-using scenarios. Re-inlining any block re-introduces its
// fingerprint and flips this RED.
func TestGenericBlocksAbsentFromFORuntime(t *testing.T) {
	markNonAC(t, "dedup lint; behavior via live team-using scenarios")
	fo := foRuntime(t)
	for block, fp := range genericBlockFingerprints {
		if strings.Contains(fo, fp) {
			t.Errorf("claude-first-officer-runtime.md still inlines %s block fingerprint %q (moved, not duplicated)", block, fp)
		}
	}
}

// TestFORuntimeInvokesSkill is a non-AC text-consistency lint: it asserts the FO
// runtime's `## Team Creation` section carries the
// Skill(skill="spacedock:using-claude-team") invocation literal and no disproven
// cross-skill @-include. Per the proof policy this presence check does NOT prove
// the FO invokes the skill: an inverted clause keeps the substring. The behavior —
// the FO loads the team-harness discipline and runs a real team — is exercised by
// every team-using live scenario. This lint guards the seam STRING and bans the
// @-include mechanism; it is the text half, not the behavioral proof.
func TestFORuntimeInvokesSkill(t *testing.T) {
	markNonAC(t, "live team-using scenarios (gate-guardrail/rejection-flow launch a real team)")
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

// TestFORuntimeDroppedMaterially is a non-AC structural lint: it asserts
// claude-first-officer-runtime.md dropped at least foRuntimeLineDropFloor lines vs
// a hardcoded pre-change baseline. This is a drift-prone numeric floor (a
// secondary signal to the fingerprint-absence dedup lints), not a behavioral
// claim; no behavior depends on the exact line count. Kept as a sanity check that
// the extraction actually removed bulk, not as proof of any AC.
func TestFORuntimeDroppedMaterially(t *testing.T) {
	markNonAC(t, "n/a — structural line-count floor, no behavior to drive")
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

// TestSkillFreeOfSpacedockTokens is a code-bound invariant: the generic
// team-harness skill is free of the spacedock-specific vocabulary. The banned set
// is DERIVED from code (spacedockLeakageTokens: the dispatch router's helper
// subcommands qualified with their `spacedock dispatch ` prefix, the
// `spacedock dispatch`/`spacedock status` command prefixes from cli.go, and the
// hyphenated stage-option keys like `feedback-to`), not a literal frozen against
// the skill — so the set shifts when the binary's command surface changes, which
// is what lets this fail as an invariant. A spacedock token leaking into the
// skill reds it. The qualified `spacedock dispatch <sub>` forms mean a bare
// English word (a generic `reconcile` in event-loop prose) never false-fires.
func TestSkillFreeOfSpacedockTokens(t *testing.T) {
	markCodeBoundInvariant(t, "spacedockLeakageTokens (dispatch router + cli.go verbs + status stage keys)")
	skill := usingClaudeTeamSkill(t)
	banned := spacedockLeakageTokens(t)
	if len(banned) == 0 {
		t.Fatal("derived zero leakage tokens — the code-side vocabulary source diverged")
	}
	for _, b := range banned {
		if strings.Contains(skill, b) {
			t.Errorf("using-claude-team SKILL.md leaks spacedock-specific token %q (must stay team-harness-generic)", b)
		}
	}
}

// TestSpacedockDecisionsStayInFORuntime is a code-bound invariant: the spacedock
// decision points REMAIN in the FO contract. The required anchors are the
// QUALIFIED dispatch-helper invocations DERIVED from the dispatch router
// (`spacedock dispatch build`, `spacedock dispatch context-budget`,
// `spacedock dispatch reconcile`) rather than literals frozen against the file —
// the binary defines these subcommands, so the anchor set tracks the real command
// surface and a subcommand renamed in code shifts the expectation. A qualified
// decision call wrongly moved out of the FO contract reds this.
func TestSpacedockDecisionsStayInFORuntime(t *testing.T) {
	markCodeBoundInvariant(t, "spacedockDispatchSubcommands (dispatch.go router)")
	fo := foRuntime(t)
	subs := spacedockDispatchSubcommands(t)
	required := spacedockQualified(subs, "build", "context-budget", "reconcile")
	if len(required) != 3 {
		t.Fatalf("expected build/context-budget/reconcile in the dispatch router, derived %v from %v", required, subs)
	}
	for _, anchor := range required {
		if !strings.Contains(fo, anchor) {
			t.Errorf("claude-first-officer-runtime.md no longer contains spacedock decision anchor %q (must stay in the FO contract)", anchor)
		}
	}
}

// spacedockQualified returns `spacedock dispatch <sub>` for each `want` that the
// router actually exposes in subs — so a decision anchor cannot name a subcommand
// the binary does not route, and a renamed subcommand drops out of the set.
func spacedockQualified(subs []string, want ...string) []string {
	have := map[string]bool{}
	for _, s := range subs {
		have[s] = true
	}
	var out []string
	for _, w := range want {
		if have[w] {
			out = append(out, "spacedock dispatch "+w)
		}
	}
	return out
}
