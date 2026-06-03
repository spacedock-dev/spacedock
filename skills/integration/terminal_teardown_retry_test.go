// ABOUTME: Structural lint over the FO terminal-teardown contract — step 10 and the
// ABOUTME: Claude runtime require team teardown to retry to success, not stop at the first failed TeamDelete.
package integration

import (
	"strconv"
	"strings"
	"testing"
)

// TestTerminalTeardownRetriesToSuccess is a structural lint pinning the
// sonnet-live-ci-flake fix: the FO terminal-teardown contract must require the
// runtime team-teardown call to RETRY TO SUCCESS after the cooperative shutdown,
// not end the turn on the first "active member(s)" failure. The archived sonnet
// CI stream proved the defect — terminal TeamDelete raced the ensign's
// cooperative shutdown, failed with "active member(s)", and the FO never retried,
// so the claude -p subprocess never exited and the live cycle hung at expectExit.
//
// Inversion-resistant (cycle-1 audit fix): a bare substring-presence lint passes
// even an INVERTED mandate ("End the turn on the first failure; do NOT re-send
// shutdown_request; do NOT call TeamDelete again") because the positive grep
// tokens survive. This lint instead asserts BOTH directions — the positive
// mandate IS present AND the negating phrases are ABSENT — scoped to the
// terminal-teardown step only (shared-core step 10; the Claude `## Terminal Team
// Teardown` section), so the audit's two inverting edits turn it RED.
//
// Oracle honesty: this lint asserts the contract carries the retry clause; it
// does NOT prove the live FO obeys it. The behavioral oracle is the live-e2e CI
// run (AC-1) — a prose-only rule's ceiling is "the wording is present," and
// TeamDelete is a Claude-runtime tool with no spacedock-binary seam to gate. The
// lint guards against the clause being dropped or INVERTED in a future edit.
func TestTerminalTeardownRetriesToSuccess(t *testing.T) {
	files := vendoredSkillFiles(t)

	// negatingPhrases are the inversion fingerprints — the exact mandates the #275
	// bug followed. Their PRESENCE in a terminal-teardown region means the contract
	// was inverted (cycle-1) or gutted (cycle-2) back to the hang. Lower-cased; the
	// regions are matched lower.
	negatingPhrases := []string{
		// cycle-1 inversions: give up on the first failure / no retry.
		"end the turn on the first",    // the inverted directive
		"do not call teamdelete again", // the inverted no-retry directive
		"do not re-send",               // the inverted no-cooperative-reshutdown directive
		"do not retry",                 // any flat no-retry directive in this region
		"stop at the first",            // give-up-on-first-failure phrasing
		// cycle-2 gut-edits: the millisecond re-race + the non-interactive hang.
		// NB: do NOT add a bare "immediately re-fire" — the CORRECT prose says
		// "do NOT immediately re-fire", so that fingerprint would false-positive.
		// The gut-edit's directive is the affirmative "immediately call TeamDelete
		// again" (no leading "do NOT"), which correct prose never contains.
		"immediately call teamdelete again",       // re-fire with no settle (re-loses the race, #275)
		"surface to the captain and end the turn", // non-interactive cap-exhaustion hang
	}

	// Shared core, step 10 in Merge and Cleanup. Scope to step 10 ALONE — the
	// other nine steps legitimately contain no retry language, and an inverting
	// edit lands inside step 10, so a whole-section scope would dilute the signal.
	// Then DROP the trailing delegation sentence ("…are the adapter's"): it NAMES
	// settle/cap/non-interactive-exit while delegating them to the adapter, so a
	// required-phrase that matched only there would let a gut-edit of the
	// behavioral mandate pass green (the cycle-2 hole). The load-bearing clauses
	// must appear in the BEHAVIORAL prose, not the delegation tail.
	core := files["first-officer/references/first-officer-shared-core.md"]
	step10 := numberedStep(sectionAfter(core, "## Merge and Cleanup"), 10)
	if step10 == "" {
		t.Fatal("FO shared core missing Merge-and-Cleanup step 10 (terminal teardown)")
	}
	step10Behavioral := stripDelegationTail(step10)
	assertDirectionalMandate(t, "shared-core step 10", step10Behavioral, negatingPhrases,
		[]string{
			"retry the team-teardown call to success",                   // retry-to-success
			"Do NOT end the turn",                                       // no give-up-on-first-failure
			"wait a short settle interval before the next teardown",     // inter-attempt settle (load-bearing, not delegated)
			"keep settling and retrying the teardown until it succeeds", // non-interactive cap-exhaustion exit obligation
			"attempt cap", // bounded fast retries
		})

	// Claude runtime: the Awaiting-Completion section already bans retrying
	// TeamDelete BEFORE the completion signal (the pre-completion wait phase). The
	// fix must distinguish that from the TERMINAL teardown phase, where TeamDelete
	// MUST be retried to success WITH an inter-attempt settle and a non-interactive
	// exit obligation. Scope to the `## Terminal Team Teardown` section. There is no
	// delegation tail here (the Claude runtime IS the adapter realization), so the
	// behavioral phrases are asserted directly.
	claude := files["first-officer/references/claude-first-officer-runtime.md"]
	teardown := sectionAfter(claude, "## Terminal Team Teardown")
	if teardown == "" {
		t.Fatal("Claude FO runtime missing the `## Terminal Team Teardown` section (the retry-to-success realization of shared-core step 10)")
	}
	assertDirectionalMandate(t, "Terminal Team Teardown", teardown, negatingPhrases,
		[]string{
			"TeamDelete",                          // names the call it retries
			"active member",                       // names the race it retries past
			"do NOT end the turn",                 // no give-up-on-first-failure
			"wait for the settle before the next", // inter-attempt settle mandate
			"sleep 2",                             // the concrete settle the gut-edit removes
			"keep the settle-then-",               // non-interactive loop obligation (settle-then-TeamDelete loop)
			"until `TeamDelete` succeeds",         // …kept going until success, not abandoned
			"cap",                                 // bounded fast retries
		})
}

// assertDirectionalMandate fails if any negating (inverted-mandate) phrase is
// present in region, or if any required positive-mandate phrase is missing. The
// two-sided check is what makes the lint inversion-resistant: an inverted edit
// either introduces a negating phrase or drops a positive one (or both), so it
// cannot pass by merely preserving grep tokens.
func assertDirectionalMandate(t *testing.T, label, region string, negating, required []string) {
	t.Helper()
	lower := strings.ToLower(region)
	for _, neg := range negating {
		if strings.Contains(lower, neg) {
			t.Errorf("%s contains the inverted-mandate phrase %q — the terminal teardown must RETRY to success, not give up on the first failure", label, neg)
		}
	}
	for _, req := range required {
		if !strings.Contains(region, req) {
			t.Errorf("%s missing required directional-mandate phrase %q", label, req)
		}
	}
}

// stripDelegationTail removes the trailing "…are the adapter's" delegation
// sentence from a teardown step. That sentence NAMES the load-bearing concepts
// (settle interval, cap value, non-interactive exit obligation) while delegating
// their realization to the adapter — so a required-phrase that matched only there
// would let a gut-edit of the BEHAVIORAL mandate pass green (the cycle-2 hole).
// Dropping it forces the behavioral phrases to be asserted against the actionable
// prose. If the delegation sentence is absent, the region is returned unchanged.
func stripDelegationTail(region string) string {
	// Target the FINAL "are the adapter's" — the trailing delegation sentence.
	// Step 10 also says "(roster and decomposition are the adapter's)" near its
	// start, so the first occurrence is the wrong one; the delegation tail is
	// always the last.
	marker := "are the adapter's"
	idx := strings.LastIndex(region, marker)
	if idx < 0 {
		return region
	}
	// Cut back to the start of the sentence containing the marker. Sentences in
	// this prose end with ". " or ".** "; find the last sentence boundary before
	// the marker and drop everything from there.
	head := region[:idx]
	if cut := strings.LastIndex(head, ". "); cut >= 0 {
		return region[:cut+1]
	}
	return region[:idx]
}

// numberedStep returns the body of the markdown ordered-list item beginning
// "{n}. " in region, up to (but excluding) the next sibling item "{n+1}. " or a
// following blank-line+non-list boundary. It scopes the lint to a single
// numbered step so an inverting edit inside that step is not diluted by the rest
// of the section. Returns "" if the item is not found.
func numberedStep(region string, n int) string {
	lines := strings.Split(region, "\n")
	startPrefix := strconv.Itoa(n) + ". "
	nextPrefix := strconv.Itoa(n+1) + ". "
	start := -1
	for i, line := range lines {
		if strings.HasPrefix(line, startPrefix) {
			start = i
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		// Next sibling numbered item, or a markdown subheading, ends the step.
		if strings.HasPrefix(lines[i], nextPrefix) || strings.HasPrefix(lines[i], "#") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}

// TestAwaitingCompletionStillBansPreCompletionTeamDelete guards the boundary the
// fix must NOT erode: BEFORE the ensign's completion signal, the FO must still
// not emit TeamDelete (retrying it during the wait phase is the original
// premature-teardown bug). The terminal retry clause is a separate phase; this
// lint ensures the Awaiting-Completion ban survives the amendment.
func TestAwaitingCompletionStillBansPreCompletionTeamDelete(t *testing.T) {
	claude := vendoredSkillFiles(t)["first-officer/references/claude-first-officer-runtime.md"]
	region := sectionAfter(claude, "## Awaiting Completion")
	if region == "" {
		t.Fatal("Claude FO runtime missing the `## Awaiting Completion` section")
	}
	if !strings.Contains(region, "emit `TeamDelete`") {
		t.Error("Awaiting Completion must still ban emitting TeamDelete before the completion signal arrives")
	}
}
