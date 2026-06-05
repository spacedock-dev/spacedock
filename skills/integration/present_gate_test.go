// ABOUTME: AC-1/AC-2 oracles for the present-gate extraction — the Gate
// ABOUTME: Presentation template + nine assembly rules live in the skill, not the FO core, with a Skill() seam.
package integration

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// assemblyRuleFingerprints uniquely identifies each of the nine captain-facing
// assembly rules moved into the present-gate skill. Each literal is unique-1 in
// the pre-change first-officer-shared-core.md (verified during ideation). The
// count is the teeth of AC-2: a dropped rule reds the absence of its
// fingerprint in the skill.
var assemblyRuleFingerprints = map[string]string{
	"lede-first/decision-last":       "Lede first, decision last",
	"chosen-direction-required":      "Chosen direction is required as FO prose",
	"cite-the-report":                "Cite the Stage Report; render a one-line gist roll-up",
	"reviewer-findings-in-tiers":     "Reviewer findings render in priority tiers",
	"recommendation-appears-once":    "Recommendation appears exactly once",
	"bounce-back-names-asks":         "Bounce-back recommendations name the concrete asks",
	"no-format-pedantry-asides":      "No format-pedantry asides",
	"one-sentence-worktree-heads-up": "One sentence of worktree heads-up when approval changes worktree state",
	"target-15-25-lines":             "Target length: 15-25 lines",
}

// gatePresentationFingerprints identifies the moved Gate-Presentation content:
// the format template plus a representative subset of the assembly rules. AC-1(a)
// asserts presence in the skill; AC-1(b) asserts absence from the FO core.
var gatePresentationFingerprints = map[string]string{
	"template":                  "Gate review: {entity title}",
	"lede-first/decision-last":  "Lede first, decision last",
	"chosen-direction-required": "Chosen direction is required as FO prose",
	"no-format-pedantry-asides": "No format-pedantry asides",
	"target-15-25-lines":        "Target length: 15-25 lines",
}

// presentGateSkill reads the new skill body under test.
func presentGateSkill(t *testing.T) string {
	t.Helper()
	p := filepath.Join(skillsRoot(t), "present-gate", "SKILL.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read present-gate SKILL.md: %v", err)
	}
	return string(b)
}

// foCore reads the FO shared-core contract under test.
func foCore(t *testing.T) string {
	t.Helper()
	return vendoredSkillFiles(t)["first-officer/references/first-officer-shared-core.md"]
}

// TestGatePresentationPresentInSkill is a non-AC text-consistency lint: it
// asserts the moved Gate-Presentation fingerprints (template + assembly rules) are
// present in skills/present-gate/SKILL.md — that the prose MOVED here, real
// authoring work. Per the proof policy (f8b257cf) this presence check does NOT
// prove the FO actually renders the gate from the skill: a meaning-inverted skill
// body keeps every fingerprint. The behavior — the FO loads present-gate via
// Skill() and presents the gate without self-approving — is proven by the live
// gate-guardrail scenario (internal/ensigncycle, runClaudeGateGuardrailScenario /
// runCodexGateGuardrailScenario, asserted by assertGateHeld) and its offline
// mutation control TestGateGuardrailNegativeBrokenStateTransition. This lint only
// guards against the fingerprints being dropped or the prose being deleted.
func TestGatePresentationPresentInSkill(t *testing.T) {
	markNonAC(t, "live gate-guardrail scenario (assertGateHeld) + TestGateGuardrailNegativeBrokenStateTransition")
	skill := presentGateSkill(t)
	for name, fp := range gatePresentationFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("present-gate SKILL.md missing %s fingerprint %q", name, fp)
		}
	}
}

// TestAllNineAssemblyRulesPresentInSkill is a non-AC text-consistency lint: it
// asserts the skill carries all nine captain-facing assembly-rule fingerprints
// (the count is the teeth — a dropped rule reds the absence of its fingerprint).
// Per the proof policy this is text authoring, not behavioral proof: an inverted
// rule body keeps the fingerprint. The behavior that the FO actually FOLLOWS the
// assembly rules when rendering a gate is proven by the live gate-guardrail
// scenario, not this presence check.
func TestAllNineAssemblyRulesPresentInSkill(t *testing.T) {
	markNonAC(t, "live gate-guardrail scenario (assertGateHeld) + TestGateGuardrailNegativeBrokenStateTransition")
	skill := presentGateSkill(t)
	if len(assemblyRuleFingerprints) != 9 {
		t.Fatalf("expected 9 assembly-rule fingerprints, have %d", len(assemblyRuleFingerprints))
	}
	for name, fp := range assemblyRuleFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("present-gate SKILL.md missing assembly rule %s fingerprint %q", name, fp)
		}
	}
}

// TestGatePresentationAbsentFromFOCore is a non-AC text-consistency lint (dedup):
// it asserts the moved fingerprints are NO LONGER present in
// first-officer-shared-core.md — moved, not duplicated. Whole-file (NOT
// region-scoped): region-scoping an absence check would false-pass content that
// moved elsewhere in the file. This is a structural dedup property, not a
// behavioral claim; the FO's gate behavior is proven by the live gate-guardrail
// scenario. The lint guards against the block being re-inlined (which would
// re-introduce a fingerprint and flip this RED).
func TestGatePresentationAbsentFromFOCore(t *testing.T) {
	markNonAC(t, "dedup lint; behavior via live gate-guardrail scenario (assertGateHeld)")
	fo := foCore(t)
	for name, fp := range gatePresentationFingerprints {
		if strings.Contains(fo, fp) {
			t.Errorf("first-officer-shared-core.md still inlines %s fingerprint %q (moved, not duplicated)", name, fp)
		}
	}
}

// presentGateAtInclude matches an `@`-prefixed path token that resolves toward
// the present-gate skill — `@present-gate`, `@./present-gate`,
// `@../present-gate/SKILL.md`, etc. The leading `@` plus any run of relative-path
// segments (`./`, `../`, bare) ending in a `present-gate` path component is the
// disproven cross-skill include the seam must NOT use. A structural scan, not an
// enumerated literal table — it catches the `@./present-gate/...` family the old
// enum missed.
var presentGateAtInclude = regexp.MustCompile(`@(?:\.{1,2}/)*present-gate\b`)

// TestFOCoreInvokesPresentGateSkill is a non-AC text-consistency lint: it asserts
// the FO core's `## Completion and Gates` section carries the
// Skill(skill="spacedock:present-gate") invocation literal and no disproven
// cross-skill @-include. Per the proof policy this presence check does NOT prove
// the FO invokes the skill: a meaning-inverted clause ("NEVER invoke present-gate;
// self-approve silently") keeps the Skill(...) substring and passes (verified in
// ideation — the mutation harness left this GREEN under inversion). The behavior —
// the FO actually loads present-gate and presents the gate without self-approving
// — is proven only by the live gate-guardrail scenario (assertGateHeld) and its
// offline mutation control. This lint guards the seam STRING (so the skill name in
// the contract and the skill's own frontmatter cannot silently drift apart) and
// bans the @-include mechanism; it is the text half, not the behavioral proof.
func TestFOCoreInvokesPresentGateSkill(t *testing.T) {
	markNonAC(t, "live gate-guardrail scenario (assertGateHeld) + TestGateGuardrailNegativeBrokenStateTransition")
	fo := foCore(t)
	region := sectionAfter(fo, "## Completion and Gates")
	if region == "" {
		t.Fatal("first-officer-shared-core.md has no `## Completion and Gates` section")
	}
	if !strings.Contains(region, `Skill(skill="spacedock:present-gate")`) {
		t.Errorf("`## Completion and Gates` section does not invoke Skill(skill=\"spacedock:present-gate\")")
	}
	if m := presentGateAtInclude.FindString(region); m != "" {
		t.Errorf("`## Completion and Gates` section uses the disproven cross-skill @-include %q", m)
	}
}

// presentGateBannedHelperPrefixes selects, from the code-derived spacedock
// vocabulary, the dispatch/status command PREFIXES the gate-presentation skill
// must not name — its prose is FO judgment/format, not shell wiring. It
// deliberately omits the stage-option keys (the skill legitimately references
// `{feedback-to target}` when describing a bounce-back decision), so it bans only
// the qualified command invocations.
func presentGateBannedHelperPrefixes(t *testing.T) []string {
	t.Helper()
	var out []string
	for _, tok := range spacedockLeakageTokens(t) {
		if tok == "spacedock dispatch" || tok == "spacedock status" {
			out = append(out, tok)
		}
	}
	return out
}

// TestPresentGateSkillFreeOfDispatchHelperLeak is a code-bound invariant: the
// gate-presentation skill is free of the binary's `spacedock dispatch` /
// `spacedock status` command prefixes. The expected token set is DERIVED from the
// binary's registered command verbs (spacedockTopLevelCommands), not a literal
// frozen against the skill — so it diverges when a command verb is renamed in
// cli.go, which is what lets this fail as an invariant. A `spacedock dispatch` /
// `spacedock status` token leaking into the skill reds it.
func TestPresentGateSkillFreeOfDispatchHelperLeak(t *testing.T) {
	markCodeBoundInvariant(t, "spacedockTopLevelCommands (cli.go Use: verbs)")
	skill := presentGateSkill(t)
	banned := presentGateBannedHelperPrefixes(t)
	if len(banned) == 0 {
		t.Fatal("derived zero command prefixes — the cli.go command surface diverged")
	}
	for _, b := range banned {
		if strings.Contains(skill, b) {
			t.Errorf("present-gate SKILL.md leaks spacedock command prefix %q (gate-presentation prose is FO judgment, not shell wiring)", b)
		}
	}
}

// presentGateFrontmatterValue returns the trimmed scalar value of a top-level
// `key:` line in skills/present-gate/SKILL.md's YAML frontmatter, with any
// surrounding quotes stripped. The bool reports whether the key was found.
func presentGateFrontmatterValue(t *testing.T, key string) (string, bool) {
	t.Helper()
	fm, ok := frontmatter(presentGateSkill(t))
	if !ok {
		t.Fatal("present-gate SKILL.md has no YAML frontmatter block")
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

// presentGateSeamName is the seam target name the FO contract actually invokes.
// It is read from a DIFFERENT file than the skill under test — the FO shared core,
// the file that drives the FO — so the skill's frontmatter `name:` and the
// contract's `Skill(skill="spacedock:NAME")` invocation have independent sources
// that can diverge.
const presentGateSeamLiteral = "present-gate"

// TestPresentGateSkillNameMatchesSeam is a code-bound invariant: the skill's
// frontmatter `name:` equals the seam name the FO CONTRACT invokes
// (Skill(skill="spacedock:NAME") in first-officer-shared-core.md). The expected
// value comes from the contract, not from the skill file under test — so renaming
// the skill's frontmatter, or renaming the contract's invocation, makes the two
// diverge and reds this. That is the seam-drift the check exists to catch: a
// renamed skill the FO's invocation no longer reaches.
func TestPresentGateSkillNameMatchesSeam(t *testing.T) {
	markCodeBoundInvariant(t, "FO contract Skill(skill=\"spacedock:present-gate\") invocation (first-officer-shared-core.md)")
	name, ok := presentGateFrontmatterValue(t, "name")
	if !ok {
		t.Fatal("present-gate SKILL.md frontmatter has no name field")
	}
	seam := invokedSeamName(foCore(t), presentGateSeamLiteral)
	if seam == "" {
		t.Fatalf("FO contract does not invoke Skill(skill=\"spacedock:%s\") — the seam the skill name must match is gone", presentGateSeamLiteral)
	}
	if name != seam {
		t.Errorf("present-gate SKILL.md frontmatter name is %q, but the FO contract invokes the seam %q — a renamed skill the FO invocation no longer reaches", name, seam)
	}
}

// TestPresentGateSkillIsFOInternal is a code-bound invariant binding the skill's
// `user-invocable` frontmatter to its ROLE in the FO contract: a skill the FO
// invokes mid-run via Skill(skill="spacedock:NAME") is FO-internal and MUST be
// `user-invocable: false`, never a captain-facing user skill. The expected value
// is not a free literal — it is REQUIRED by the contract invoking the seam: the
// presence of the invocation (an independent source) is what makes
// `user-invocable: true` wrong. Flipping the frontmatter to `true` while the FO
// still invokes the seam reds this.
func TestPresentGateSkillIsFOInternal(t *testing.T) {
	markCodeBoundInvariant(t, "FO contract Skill(skill=\"spacedock:present-gate\") invocation implies FO-internal")
	if invokedSeamName(foCore(t), presentGateSeamLiteral) == "" {
		t.Fatalf("FO contract does not invoke the present-gate seam — the FO-internal premise no longer holds")
	}
	v, ok := presentGateFrontmatterValue(t, "user-invocable")
	if !ok {
		t.Fatal("present-gate SKILL.md frontmatter has no user-invocable field")
	}
	if v != "false" {
		t.Errorf("present-gate SKILL.md frontmatter user-invocable is %q, but the FO contract invokes it as a mid-run seam — an FO-internal skill must be user-invocable: false", v)
	}
}
