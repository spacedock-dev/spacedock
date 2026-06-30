// ABOUTME: Structural-absence guard keeping concrete workflow stage-name enumerations
// ABOUTME: out of the universal ensign core, which loads per-dispatch and must stay stage-neutral.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ensignSharedCore is the universal ensign core, loaded on EVERY worker dispatch via
// the `Skill(skill="spacedock:ensign")` first-action. Because the ensign also runs
// non-dev workflows (ticket, experiment, survey), the core must name no concrete
// workflow stage: a stage-name parenthetical here is a per-dispatch stage-neutrality
// leak. The dev stage discipline rides the dev-shape scaffolding the dev ensign
// fetches per-dispatch (see internal/dispatch's behavior-loss control), not this core.
var ensignSharedCore = filepath.Join("ensign", "references", "ensign-shared-core.md")

// stageNameAlt is the dev workflow's concrete stage vocabulary. A parenthetical that
// enumerates two or more of these names — joined by a comma, slash, semicolon, or
// whitespace — is a stage-neutrality leak.
const stageNameAlt = `(?:backlog|ideation|implementation|validation|done)`

// stageEnumParenRe matches a parenthetical ENUMERATING workflow stage names — an open
// paren wrapping a list of two or more stage-name tokens joined by any of the common
// separators `,` `/` `;` or whitespace, e.g. `(implementation, validation)`,
// `(implementation/validation)`, `(implementation; validation)`, or
// `(implementation validation)`. The separator class is widened past the comma so the
// natural shorthands cannot re-leak the vocab unflagged. The expected value comes from
// the EXTERNAL rule (the universal core is stage-neutral), not the file's own prose, so
// a stage-neutral paraphrase carries no match and a re-introduced stage-name
// parenthetical does — same structural-absence family as
// TestFOContractCoresHaveNoDeferredTierToken.
//
// By-design tradeoffs (kept narrow to avoid over-flagging incidental English):
//   - A bare word is NOT matched — it would wrongly flag "Signaling done" (core line 51).
//   - A single-name parenthetical (e.g. `(implementation)` alone) is NOT matched; the
//     leak shape this guards is a multi-name ENUMERATION, and the `+` requires ≥2 names.
//   - Non-dev stage vocab (a ticket/experiment/survey workflow's own stage names) is NOT
//     matched — stageNameAlt lists only the dev stages, the host whose names leaked here.
//   - The parenthetical must START with a stage name and be a pure stage-name list, so an
//     incidental "done" buried in English prose inside parens does not engage the scanner.
var stageEnumParenRe = regexp.MustCompile(`\(\s*` + stageNameAlt + `\b(?:[,/;\s]+` + stageNameAlt + `\b)+\s*\)`)

// lineEnumeratesStageNames reports whether a line carries a workflow-stage-enumeration
// parenthetical. This is the single scanner the absence check, its discriminator
// control, and the re-insertion control all drive, so defeating it reds every caller.
func lineEnumeratesStageNames(line string) bool {
	return stageEnumParenRe.MatchString(line)
}

// scanForStageEnumeration returns the 1-based line numbers of content whose text
// carries a stage-enumeration parenthetical.
func scanForStageEnumeration(content string) []int {
	var hits []int
	for i, line := range strings.Split(content, "\n") {
		if lineEnumeratesStageNames(line) {
			hits = append(hits, i+1)
		}
	}
	return hits
}

// TestEnsignCoreEnumeratesNoStageNames is a structural-ABSENCE check: the universal
// ensign core may carry no parenthetical enumerating concrete workflow stage names.
// The leak this guards against was the Split-Root State Contract's worktree-illustration
// parentheticals (`(implementation, validation)`, `(ideation, backlog)`); the whole
// file is scanned so a re-introduction at any site fails. Non-vacuity is held by the
// paired discriminator and re-insertion controls below.
func TestEnsignCoreEnumeratesNoStageNames(t *testing.T) {
	path := filepath.Join(skillsRoot(t), ensignSharedCore)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ensign core %s: %v", path, err)
	}
	lines := strings.Split(string(data), "\n")
	for _, ln := range scanForStageEnumeration(string(data)) {
		t.Errorf("%s:%d enumerates concrete workflow stage names in a parenthetical — the universal ensign core loads per-dispatch and must stay stage-neutral (dev stage discipline rides the dev-shape scaffolding): %q", path, ln, strings.TrimSpace(lines[ln-1]))
	}
}

// TestStageEnumerationScannerDiscriminates is the DISCRIMINATOR control: it proves the
// scanner flags the genuine pre-strip leak shapes and passes the stage-neutral phrasing
// the strip produces (plus the incidental "Signaling done") — so
// TestEnsignCoreEnumeratesNoStageNames can never pass vacuously (e.g. by a typo'd regex
// that never matches anything).
func TestStageEnumerationScannerDiscriminates(t *testing.T) {
	// Genuine pre-strip leak lines — each MUST flag. The comma form is the original
	// leak; the slash, semicolon, and bare-whitespace forms are the natural shorthands
	// an author would reach for, and the widened separator class must bite on each.
	mustFlag := []string{
		"With a worktree (implementation, validation), the worktree isolates the deliverable work product only.",
		"Without one (ideation, backlog), you run from the repo root; entity/report still go to the state checkout.",
		"With a worktree (implementation/validation), the worktree isolates the deliverable work product only.",
		"With a worktree (implementation; validation), the worktree isolates the deliverable work product only.",
		"With a worktree (implementation validation), the worktree isolates the deliverable work product only.",
	}
	for _, line := range mustFlag {
		if !lineEnumeratesStageNames(line) {
			t.Errorf("discriminator control: a genuine stage-enumeration leak line was NOT flagged (the scanner would pass vacuously): %q", line)
		}
	}

	// Stage-neutral phrasings (and incidental stage words) — each MUST pass.
	mustPass := []string{
		"With a worktree, the worktree isolates the deliverable work product only. Without one, you run from the repo root; entity/report still go to the state checkout.",
		"Signaling done without committing forces the FO to re-dispatch just to get a commit — the most common cause of nudge loops.",
	}
	for _, line := range mustPass {
		if lineEnumeratesStageNames(line) {
			t.Errorf("discriminator control: a stage-neutral phrasing was wrongly flagged: %q", line)
		}
	}
}

// TestEnsignCoreStageEnumerationGuardRedsOnReinsertion is the RE-INSERTION control: it
// drives the same scanner over the REAL shipped core after splicing each stage-name
// parenthetical back into the worktree-isolation sentence, and asserts the absence scan
// reds. This proves the guard catches a regression in the file itself, not only in
// hand-written discriminator strings — closing the "the regex never matched the real
// content anyway" vacuity gap.
func TestEnsignCoreStageEnumerationGuardRedsOnReinsertion(t *testing.T) {
	path := filepath.Join(skillsRoot(t), ensignSharedCore)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ensign core %s: %v", path, err)
	}
	original := string(data)
	if hits := scanForStageEnumeration(original); len(hits) != 0 {
		t.Fatalf("re-insertion control precondition: shipped core already enumerates stage names at line(s) %v — the absence check should have caught this", hits)
	}

	reinsertions := []struct {
		name   string
		anchor string
		leak   string
	}{
		{"implementation/validation", "With a worktree, the worktree isolates", "With a worktree (implementation, validation), the worktree isolates"},
		{"ideation/backlog", "Without one, you run from the repo root", "Without one (ideation, backlog), you run from the repo root"},
	}
	for _, rc := range reinsertions {
		if !strings.Contains(original, rc.anchor) {
			t.Fatalf("re-insertion control %q: anchor not found in the core, cannot mutate: %q", rc.name, rc.anchor)
		}
		mutated := strings.Replace(original, rc.anchor, rc.leak, 1)
		if mutated == original {
			t.Fatalf("re-insertion control %q: splicing the leak back in was a no-op", rc.name)
		}
		if hits := scanForStageEnumeration(mutated); len(hits) == 0 {
			t.Errorf("re-insertion control %q: the guard did NOT red after re-inserting %q into the core (the absence check is vacuous against this regression)", rc.name, rc.leak)
		}
	}
}
