// ABOUTME: Structural-absence guard keeping the DEFERRED-72 (fo-tier-delegation) tier
// ABOUTME: vocabulary out of the shipped FO contract cores until member 72 un-defers.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// foContractCores are the shipped FO contract cores the gate-verdict body lives in —
// the boot-resident shared core and the two host-neutral cores it names, plus each
// host runtime adapter. Every host loads the shared core at boot, so a deferred-tier
// token anywhere here reaches a booting FO as an instruction to escalate a verdict to
// a mechanism that does not ship (member 72 — fo-tier-delegation — is DEFERRED).
var foContractCores = []string{
	filepath.Join("first-officer", "references", "first-officer-shared-core.md"),
	filepath.Join("first-officer", "references", "fo-dispatch-core.md"),
	filepath.Join("first-officer", "references", "fo-merge-core.md"),
	filepath.Join("first-officer", "references", "claude-first-officer-runtime.md"),
	filepath.Join("first-officer", "references", "codex-first-officer-runtime.md"),
	filepath.Join("first-officer", "references", "pi-first-officer-runtime.md"),
	filepath.Join("first-officer", "references", "claude-fo-dispatch.md"),
}

// deferredTierLiteralTokens are the LITERAL deferred-72 tier-vocabulary tokens that must
// NOT appear in a shipped FO contract core. They are member 72's (fo-tier-delegation)
// tier/level-3-judge mechanism — DEFERRED — so nothing in the shipped contract defines
// an FO "level" or a "level-3 judge". Their presence IS the defect (structural-absence,
// same family as claudeModelTokens / claudeTeamDispatchTokens), not a paraphrasable
// meaning: a mechanism-neutral statement that the FO renders its own `Recommend` line
// carries none of these tokens and passes. The tier mechanism re-enters WITH member 72
// (and its own `«fo.tier»` source) when 72 un-defers, not before.
var deferredTierLiteralTokens = []string{
	"level-2-only",
	"level-3 judge",
	"«fo.tier»",
}

// bareL3RouteRe matches the bare `L3` route framing as a word (`routes to L3`), not an
// incidental `L3` inside a hash or path. The leak shape was "the decision routes to L3
// or the FO's own present-gate Recommend line"; a word-bounded match flags that framing
// while leaving an unrelated `L3` substring alone.
var bareL3RouteRe = regexp.MustCompile(`\bL3\b`)

// lineLeaksDeferredTierToken reports whether a line carries any deferred-72 tier-vocabulary
// token. This is the single scanner the absence check and its discriminator control both
// drive, so defeating it reds both callers.
func lineLeaksDeferredTierToken(line string) bool {
	for _, tok := range deferredTierLiteralTokens {
		if strings.Contains(line, tok) {
			return true
		}
	}
	return bareL3RouteRe.MatchString(line)
}

// TestFOContractCoresHaveNoDeferredTierToken is a structural-ABSENCE check: no shipped FO
// contract core may carry a deferred-72 tier-vocabulary token (`level-2-only`,
// `level-3 judge`, `«fo.tier»`, the bare `L3` route framing). The leak was in
// `«gate.assemble-verdict»`'s decide-effect body, but the whole file set is scanned so a
// re-introduction at any site fails. The expected value comes from the rule (member 72 is
// DEFERRED, so its tier vocabulary must not ship), not the file's own prose, so a
// mechanism-neutral phrasing (the FO renders its own `Recommend` line) passes and a
// re-introduced tier escalation fails — same family as TestDispatchCoreHasNoClaudeModelToken.
// The paired discriminator control keeps this non-vacuous.
func TestFOContractCoresHaveNoDeferredTierToken(t *testing.T) {
	root := skillsRoot(t)
	for _, rel := range foContractCores {
		path := filepath.Join(root, rel)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read FO contract core %s: %v", path, err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			if lineLeaksDeferredTierToken(line) {
				t.Errorf("%s:%d carries a deferred-72 tier-vocabulary token in a shipped FO contract core (member 72 fo-tier-delegation is DEFERRED — the verdict is the FO's own `Recommend` line): %q", path, i+1, strings.TrimSpace(line))
			}
		}
	}
}

// TestDeferredTierTokenScannerDiscriminates is the DISCRIMINATOR control for the
// deferred-tier absence check: it proves the scanner flags the genuine pre-strip leak
// shapes and passes the mechanism-neutral phrasings the strip produces — so
// TestFOContractCoresHaveNoDeferredTierToken can never pass vacuously (e.g. by a typo'd
// token never matching anything).
func TestDeferredTierTokenScannerDiscriminates(t *testing.T) {
	// Genuine pre-strip leak shapes — each MUST flag.
	mustFlag := []string{
		"a level-2-only FO escalates the decision to a level-3 judge; a capable FO renders its own `Recommend` line",
		"the decision routes to L3 or the FO's own `present-gate` `Recommend` line",
		"the tier mechanism rides member 72 WITH its «fo.tier» source",
	}
	for _, line := range mustFlag {
		if !lineLeaksDeferredTierToken(line) {
			t.Errorf("discriminator control: a genuine deferred-tier leak line was NOT flagged (the scanner would pass vacuously): %q", line)
		}
	}

	// Mechanism-neutral phrasings the strip produces — each MUST pass.
	mustPass := []string{
		"The verdict is irreducible JUDGMENT. The FO renders its own `Recommend` line.",
		"the decision is the FO's own `present-gate` `Recommend` line",
	}
	for _, line := range mustPass {
		if lineLeaksDeferredTierToken(line) {
			t.Errorf("discriminator control: a mechanism-neutral phrasing was wrongly flagged: %q", line)
		}
	}
}
