// ABOUTME: Structural-absence guard keeping host-lane literals (claude-live/codex-live/pi-live/TestLive*)
// ABOUTME: out of the runtime-neutral FO Working Principles, so the self-evidence bar stays host-agnostic.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// hostLaneLiteralRe matches a concrete host CI-lane literal — the per-host live
// lanes (claude-live/codex-live/pi-live) or a host live-test name (TestLive*). The
// FO's runtime-neutral Working Principles must name none of these: the self-evidence
// bar is host-agnostic, and its dev-specific path→lane realization lives in
// docs/dev/README.md, not the shared core. The expected value ("no host-lane
// literal") comes from the EXTERNAL runtime-neutrality requirement, not the file's
// own prose, so a host-neutral paraphrase carries no match and a host-coupled
// wording does — the structural-absence family, not a banned prose-grep.
var hostLaneLiteralRe = regexp.MustCompile(`claude-live|codex-live|pi-live|TestLive\w*`)

// workingPrinciplesBlock returns the `## Working Principles` section body (heading to
// the next `## ` or EOF). Scoping the scan to this block binds the runtime-neutrality
// requirement to the section that carries the self-evidence bar, so a host-lane
// literal introduced anywhere in it reds — while a legitimate lane reference in a
// different, host-specific section of the file would not be in scope.
func workingPrinciplesBlock(t *testing.T, body string) string {
	t.Helper()
	loc := regexp.MustCompile(`(?m)^## Working Principles$`).FindStringIndex(body)
	if loc == nil {
		t.Fatal("shared core has no `## Working Principles` block — the self-evidence bar's home section is missing")
	}
	rest := body[loc[1]:]
	if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return rest
}

// scanForHostLaneLiteral returns the 1-based line numbers within content whose text
// names a concrete host CI lane.
func scanForHostLaneLiteral(content string) []int {
	var hits []int
	for i, line := range strings.Split(content, "\n") {
		if hostLaneLiteralRe.MatchString(line) {
			hits = append(hits, i+1)
		}
	}
	return hits
}

// TestWorkingPrinciplesNamesNoHostLane is the AC-2 structural-absence check: the FO
// Working Principles — where the self-evidence bar binds the FO's own gate/merge/triage
// decisions — may carry no host-lane literal, so the runtime-neutrality the entity's
// title requires is enforced, not asserted. The incident-shaped natural wording bakes
// in `claude-live` / `TestLive*`; this reds on it. Non-vacuity is held by the paired
// discriminator and re-insertion controls below.
func TestWorkingPrinciplesNamesNoHostLane(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, sharedCorePath()))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	block := workingPrinciplesBlock(t, string(data))
	lines := strings.Split(block, "\n")
	for _, ln := range scanForHostLaneLiteral(block) {
		t.Errorf("Working Principles line %d names a host CI lane — the runtime-neutral self-evidence bar must stay host-agnostic (the dev path→lane realization lives in docs/dev/README.md): %q", ln, strings.TrimSpace(lines[ln-1]))
	}
}

// TestHostLaneScannerDiscriminates is the DISCRIMINATOR control: the scanner flags
// the host-coupled incident wording and passes the host-neutral paraphrase the bar
// actually ships — so TestWorkingPrinciplesNamesNoHostLane can never pass vacuously
// (e.g. via a typo'd regex that never matches anything).
func TestHostLaneScannerDiscriminates(t *testing.T) {
	// Host-coupled wordings — each MUST flag. These are the natural incident-shaped
	// phrasings that bake a single host's lane into the bar; the check must bite each.
	mustFlag := []string{
		"a relevant claude-live check that flakes is re-run to green, never skipped",
		"read the failure from TestLiveZeroDiscover, not the inherited label",
		"the codex-live lane is unapproved, so it is not a pass",
		"the pi-live lane must actually run and pass",
	}
	for _, line := range mustFlag {
		if !hostLaneLiteralRe.MatchString(line) {
			t.Errorf("discriminator control: a host-coupled line was NOT flagged (the scanner would pass vacuously): %q", line)
		}
	}

	// Host-neutral paraphrases (and incidental lowercase "test") — each MUST pass.
	mustPass := []string{
		`A result is "green" only when the relevant check actually ran and passed.`,
		"a relevant check that flakes is re-run to green (serial, isolated), never skipped",
		"A failure is read from this run's evidence — the failing test, assertion, or error in front of you.",
	}
	for _, line := range mustPass {
		if hostLaneLiteralRe.MatchString(line) {
			t.Errorf("discriminator control: a host-neutral paraphrase was wrongly flagged: %q", line)
		}
	}
}

// TestWorkingPrinciplesHostLaneGuardRedsOnReinsertion is the RE-INSERTION control: it
// drives the same scanner over the REAL shipped Working Principles after splicing a
// host-lane literal into the self-evidence bar, and asserts the absence scan reds —
// proving the guard catches a regression in the file itself, not only in hand-written
// discriminator strings.
func TestWorkingPrinciplesHostLaneGuardRedsOnReinsertion(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, sharedCorePath()))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	original := string(data)
	block := workingPrinciplesBlock(t, original)
	if hits := scanForHostLaneLiteral(block); len(hits) != 0 {
		t.Fatalf("re-insertion control precondition: shipped Working Principles already names a host lane at line(s) %v — the absence check should have caught this", hits)
	}

	reinsertions := []struct {
		name   string
		anchor string
		leak   string
	}{
		{"claude-live", "a relevant check that flakes is re-run to green (serial, isolated)", "a relevant claude-live check that flakes is re-run to green (serial, isolated)"},
		{"TestLive", "the failing test, assertion, or error in front of you", "the failing TestLiveZeroDiscover, assertion, or error in front of you"},
	}
	for _, rc := range reinsertions {
		if !strings.Contains(block, rc.anchor) {
			t.Fatalf("re-insertion control %q: anchor not found in the Working Principles block, cannot mutate: %q", rc.name, rc.anchor)
		}
		mutated := strings.Replace(block, rc.anchor, rc.leak, 1)
		if mutated == block {
			t.Fatalf("re-insertion control %q: splicing the leak back in was a no-op", rc.name)
		}
		if hits := scanForHostLaneLiteral(mutated); len(hits) == 0 {
			t.Errorf("re-insertion control %q: the guard did NOT red after re-inserting %q (the absence check is vacuous against this regression)", rc.name, rc.leak)
		}
	}
}
