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

// TestGatePresentationPresentInSkill locks AC-1(a): the moved Gate-Presentation
// fingerprints (template + assembly rules) are present in
// skills/present-gate/SKILL.md.
func TestGatePresentationPresentInSkill(t *testing.T) {
	skill := presentGateSkill(t)
	for name, fp := range gatePresentationFingerprints {
		if !strings.Contains(skill, fp) {
			t.Errorf("present-gate SKILL.md missing %s fingerprint %q", name, fp)
		}
	}
}

// TestAllNineAssemblyRulesPresentInSkill locks AC-2(a): the skill carries all
// nine captain-facing assembly-rule fingerprints — the count is the teeth, a
// dropped rule reds the absence of its fingerprint.
func TestAllNineAssemblyRulesPresentInSkill(t *testing.T) {
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

// TestGatePresentationAbsentFromFOCore locks AC-1(b): the moved fingerprints are
// NO LONGER present in first-officer-shared-core.md — moved, not duplicated.
// Whole-file (NOT region-scoped): region-scoping an absence check would
// false-pass content that moved elsewhere in the file. Negative-proof:
// re-inlining the block re-introduces a fingerprint and flips this RED.
func TestGatePresentationAbsentFromFOCore(t *testing.T) {
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

// TestFOCoreInvokesPresentGateSkill locks AC-1(c): the FO core's `## Completion
// and Gates` section invokes the skill via Skill(...) at the gate point and does
// NOT use the spike-disproven cross-skill @-include. Region-scoped to
// `## Completion and Gates` (the positive Skill()-present / @-absent assertions
// only). The Skill(...) literal is the integration seam; any `@`-token resolving
// toward present-gate is the disproven mechanism.
func TestFOCoreInvokesPresentGateSkill(t *testing.T) {
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

// presentGateLeakageLiterals are spacedock dispatch-helper tokens the
// gate-presentation skill must NOT name — the prose is FO judgment/format, not
// shell wiring. Mirrors the sibling using-claude-team leakage table.
var presentGateLeakageLiterals = []string{
	"spacedock dispatch",
	"spacedock status",
}

// TestPresentGateSkillFreeOfDispatchHelperLeak locks AC-2 (absence half): the
// gate-presentation skill is free of any spacedock-dispatch-helper token.
// Negative-proof: a `spacedock dispatch`/`spacedock status` token leaking into
// the skill reds this.
func TestPresentGateSkillFreeOfDispatchHelperLeak(t *testing.T) {
	skill := presentGateSkill(t)
	for _, banned := range presentGateLeakageLiterals {
		if strings.Contains(skill, banned) {
			t.Errorf("present-gate SKILL.md leaks spacedock dispatch-helper token %q (gate-presentation prose is FO judgment, not shell wiring)", banned)
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

// TestPresentGateSkillNameMatchesSeam locks AC-2: the frontmatter `name:` VALUE
// equals `present-gate` — the directory name AND the
// `Skill(skill="spacedock:present-gate")` invocation seam. Token-presence alone
// (skill_surface_test.go) would pass a renamed skill that the seam no longer
// reaches; binding the value to the seam target catches that drift. Negative-
// proof: a bogus name value reds this.
func TestPresentGateSkillNameMatchesSeam(t *testing.T) {
	name, ok := presentGateFrontmatterValue(t, "name")
	if !ok {
		t.Fatal("present-gate SKILL.md frontmatter has no name field")
	}
	if name != "present-gate" {
		t.Errorf("present-gate SKILL.md frontmatter name is %q, want %q (the directory name and the Skill(skill=\"spacedock:present-gate\") seam)", name, "present-gate")
	}
}

// TestPresentGateSkillIsFOInternal locks AC-2: the frontmatter carries
// `user-invocable: false` — the skill is FO-internal (loaded mid-run via
// Skill()), not a captain-facing user skill. Negative-proof: flipping to `true`
// reds this.
func TestPresentGateSkillIsFOInternal(t *testing.T) {
	v, ok := presentGateFrontmatterValue(t, "user-invocable")
	if !ok {
		t.Fatal("present-gate SKILL.md frontmatter has no user-invocable field")
	}
	if v != "false" {
		t.Errorf("present-gate SKILL.md frontmatter user-invocable is %q, want \"false\" (the skill is FO-internal)", v)
	}
}
