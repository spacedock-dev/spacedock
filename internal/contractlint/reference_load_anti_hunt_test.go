// ABOUTME: AC-2 presence + deletion-guard for the reference-load anti-hunt rule in
// ABOUTME: the Deferred load points section — the reference-load twin of Startup step 2's no-hunt rule.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// deferredLoadPointsSection extracts the `## Deferred load points` section body
// (heading to the next `## `), so the presence check is scoped to the section Move 1
// targets, not the whole shared core.
func deferredLoadPointsSection(t *testing.T, body string) string {
	t.Helper()
	loc := regexp.MustCompile(`(?m)^## Deferred load points$`).FindStringIndex(body)
	if loc == nil {
		t.Fatal("shared core has no `## Deferred load points` section")
	}
	rest := body[loc[1]:]
	if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return rest
}

// referenceResolutionTokens are the literal tokens that together establish the
// skill-install-dir/cwd resolution directive: a `references/…` load point resolves
// against the skill's own install directory, not the FO's working directory.
var referenceResolutionTokens = []string{
	"install directory",
	"working directory",
}

// noBroadSearchTokens are the literal tokens that together establish the no-hunt
// directive: STOP and report an unresolved reference read rather than searching the
// filesystem for a contract file known by name.
var noBroadSearchTokens = []string{
	"STOP and report",
	"do NOT broad-search",
}

// sectionHasReferenceLoadAntiHuntRule reports whether a section body carries BOTH the
// skill-install-dir/cwd resolution directive AND the no-broad-search directive — the
// semantic invariant Move 1 adds, not either half alone. A resolution-only paragraph
// re-anchors the read but does not forbid the hunt fallback; a no-hunt-only paragraph
// forbids the fallback without saying where the read should resolve instead.
func sectionHasReferenceLoadAntiHuntRule(section string) bool {
	hasResolution := true
	for _, tok := range referenceResolutionTokens {
		if !strings.Contains(section, tok) {
			hasResolution = false
			break
		}
	}
	hasNoBroadSearch := true
	for _, tok := range noBroadSearchTokens {
		if !strings.Contains(section, tok) {
			hasNoBroadSearch = false
			break
		}
	}
	return hasResolution && hasNoBroadSearch
}

// TestDeferredLoadPointsHasReferenceLoadAntiHuntRule is the AC-2 mechanism guard: the
// shipped `## Deferred load points` section carries the reference-load anti-hunt rule
// added by Move 1 — a resolution directive (a `references/…` load point resolves
// against the skill install dir, not cwd) AND a no-broad-search directive (STOP and
// report rather than hunt the filesystem) — structurally parallel to the Startup step 2
// zero-discovery no-hunt rule. This is a regression guard against silent deletion, NOT
// proof the live find-hunt rate changed; per the gate's mechanism→value pairing it
// counts only paired with AC-1, which measures the value this mechanism serves.
func TestDeferredLoadPointsHasReferenceLoadAntiHuntRule(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, sharedCorePath()))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	section := deferredLoadPointsSection(t, string(data))
	if !sectionHasReferenceLoadAntiHuntRule(section) {
		t.Errorf("`## Deferred load points` is missing the reference-load anti-hunt rule (a skill-install-dir/cwd resolution directive AND a no-broad-search directive) — see Startup step 2's zero-discovery no-hunt rule for the sibling shape")
	}
}

// TestReferenceLoadAntiHuntRuleGuardFailsOnDeletion is the AC-2 red-control: it proves
// the presence check can go RED, not just pass. A fixture carrying the full Move 1
// shape passes; deleting the paragraph reds; and a resolution-only or
// no-broad-search-only fragment (half the invariant) each red too, so the check cannot
// pass vacuously off either directive alone.
func TestReferenceLoadAntiHuntRuleGuardFailsOnDeletion(t *testing.T) {
	preamble := "- `Skill(skill=\"spacedock:fo-dispatch-recovery\")` — dispatch failure recovery.\n\n"

	full := preamble + "**A `references/…` load point resolves against this skill's own install directory** — never against the FO's current working directory. If a `references/…` read does not resolve, STOP and report the unresolved path — do NOT broad-search the filesystem to hunt a contract or reference file you know by name.\n"
	if !sectionHasReferenceLoadAntiHuntRule(full) {
		t.Fatal("control: the full Move 1 shape was NOT flagged as present — the scanner would pass vacuously")
	}

	deleted := preamble
	if sectionHasReferenceLoadAntiHuntRule(deleted) {
		t.Error("control: a section with the paragraph deleted was wrongly reported as carrying the rule")
	}

	resolutionOnly := preamble + "A `references/…` load point resolves against this skill's own install directory, never against the FO's current working directory.\n"
	if sectionHasReferenceLoadAntiHuntRule(resolutionOnly) {
		t.Error("control: a resolution-only fragment (missing the no-broad-search half) was wrongly reported as carrying the full rule")
	}

	noBroadSearchOnly := preamble + "If a reference read does not resolve, STOP and report the unresolved path — do NOT broad-search the filesystem to hunt it.\n"
	if sectionHasReferenceLoadAntiHuntRule(noBroadSearchOnly) {
		t.Error("control: a no-broad-search-only fragment (missing the resolution half) was wrongly reported as carrying the full rule")
	}
}
