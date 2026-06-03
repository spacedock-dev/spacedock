// ABOUTME: Structural lint over the FO terminal-teardown contract — step 10 and the
// ABOUTME: Claude runtime require BOUNDED best-effort teardown then a verbatim marker + hold; the launcher owns the exit.
package integration

import (
	"strconv"
	"strings"
	"testing"
)

// terminalTeardownMarker is the contract-mandated terminal-status sentinel the FO
// emits on cap-exhaustion of the bounded teardown. It is the load-bearing
// discriminator the AC-2 offline grade and the AC-1 live grade key on: neither
// bug shape emits it (the pre-yy give-up ends the turn silently after one failed
// TeamDelete; the post-yy retry-loop retries TeamDelete past the cap and never
// reaches a marker). The contract MUST mandate this string verbatim so the
// streamwatcher has a fixed substring to grep — drop it from the contract and the
// oracle has nothing to key on. Kept in sync with the same constant in
// internal/ensigncycle (terminalTeardownMarker).
const terminalTeardownMarker = "TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher."

// TestTerminalTeardownIsBoundedBestEffort is a structural lint pinning the
// REVERSAL of the sonnet-live-ci #282 contract: the FO terminal teardown is a
// BOUNDED best-effort — a cooperative shutdown, then a bounded set of TeamDelete
// attempts with an inter-attempt settle, then STOP, emit the verbatim
// terminal-status marker, and HOLD (no further teardown tool calls). The PROCESS
// EXIT is the launcher's responsibility, NOT the FO's.
//
// This INVERTS yy/#282's shipped mandate ("retry the team-teardown call to
// success" / "keep settling and retrying the teardown until it succeeds" / "Only
// a TeamDelete success lets the FO end its turn"). The live AC-1 confirmation
// proved retry-to-success is UNREACHABLE: the dispatched member approves shutdown
// and its session dies, but Claude Code never clears it from the team members[]
// (upstream #38116/#57681), so TeamDelete never succeeds; and claude -p will not
// self-exit while members[] is populated, so an FO-driven process exit is
// impossible. The contract therefore cannot demand a TeamDelete success OR an FO
// self-exit — it grades the FO's CORRECT bounded best-effort and hands the exit
// to the launcher.
//
// Inversion-resistant (two-sided): the lint asserts BOTH the bounded-best-effort
// mandate IS present (including the verbatim marker, the hold, and launcher-owned
// exit) AND the retry-to-success / FO-self-exit phrases are ABSENT, scoped to the
// terminal-teardown step only (shared-core step 10; the Claude `## Terminal Team
// Teardown` section). An edit re-introducing retry-to-success, or an FO self-exit,
// or dropping the marker mandate, turns it RED.
//
// Oracle honesty: this lint asserts the contract carries the bounded-best-effort
// + marker clause; it does NOT prove the live FO obeys it. The behavioral oracle
// is the live-e2e CI run (AC-1). The lint guards against the clause being dropped
// or re-inverted in a future edit.
func TestTerminalTeardownIsBoundedBestEffort(t *testing.T) {
	files := vendoredSkillFiles(t)

	// negatingPhrases are the inversion fingerprints. Two groups:
	//   - the ORIGINAL #275 give-up fingerprints (still forbidden — the terminal
	//     teardown must not silently abandon on the first failure with no marker).
	//   - the yy/#282 retry-to-success + FO-self-exit fingerprints now DISPROVEN
	//     (retry-to-success is unreachable; FO self-exit is impossible). Their
	//     presence in a terminal-teardown region means the contract was re-inverted
	//     back to a demand the runtime cannot satisfy. Lower-cased; regions matched
	//     lower.
	negatingPhrases := []string{
		// #275 give-up inversions: give up on the first failure / no retry / no marker.
		"end the turn on the first",    // the give-up directive
		"do not call teamdelete again", // the no-retry directive
		"do not retry",                 // any flat no-retry directive in this region
		"stop at the first",            // give-up-on-first-failure phrasing
		// yy/#282 retry-to-success + self-exit, now disproven by the live run.
		"retry the team-teardown call to success",            // retry-to-success (unreachable)
		"until `teamdelete` succeeds",                        // retry-to-success (unreachable)
		"only a `teamdelete` success lets the fo end",        // retry-to-success exit gate
		"keep settling and retrying the teardown until it",   // unbounded retry obligation
		"keep the settle-then-`teamdelete` loop going",       // unbounded retry obligation
		"hard-exit at the cap",                               // disproven FO self-exit
		"residual cleanup at process death",                  // disproven FO self-exit
		"the fo exits the process",                           // disproven FO self-exit
		"the fo must keep the settle-then-`teamdelete` loop", // unbounded retry obligation
		// audit-cycle-1 near-synonym re-introductions (M1): an unbounded
		// retry-to-success that swaps the exact tokens for a "rounds-until-clears"
		// paraphrase. These are the common rewordings that re-establish the
		// retry-to-success the runtime CANNOT satisfy.
		"continues issuing",  // "continues issuing settle-then-TeamDelete rounds…"
		"rounds until",       // "…rounds until the registry finally clears"
		"until the registry", // "…until the registry clears/settles/empties"
		"finally clears",     // "…until the registry finally clears"
		// audit-cycle-1 FO-self-exit re-introductions (M2): the FO killing its own
		// process group / runtime and demoting the launcher to a "backstop". Scoped
		// to the self-exit FINGERPRINTS, not over-broad: "terminates its own" (not a
		// bare "the fo terminates", which a legit "the FO terminates the attempts"
		// would false-red) and "launcher is only a" (not a bare "backstop", which
		// any legit launcher-as-backstop phrasing would false-red).
		"kill -9",            // self-kill via Bash
		"$ppid",              // kill the parent process (the harness)
		"terminates its own", // "the FO terminates its own runtime…"
		"launcher is only a", // "…the launcher is only a backstop" (demotes launcher-owned exit)
	}

	// Shared core, step 10 in Merge and Cleanup. Scope to step 10 ALONE — the
	// other nine steps legitimately contain no teardown language, and an inverting
	// edit lands inside step 10, so a whole-section scope would dilute the signal.
	// Then DROP the trailing delegation sentence ("…are the adapter's"): it NAMES
	// the cap/settle/marker while delegating their realization to the adapter, so a
	// required-phrase that matched only there would let a gut-edit of the
	// behavioral mandate pass green. The load-bearing clauses must appear in the
	// BEHAVIORAL prose, not the delegation tail.
	core := files["first-officer/references/first-officer-shared-core.md"]
	step10 := numberedStep(sectionAfter(core, "## Merge and Cleanup"), 10)
	if step10 == "" {
		t.Fatal("FO shared core missing Merge-and-Cleanup step 10 (terminal teardown)")
	}
	step10Behavioral := stripDelegationTail(step10)
	// The verbatim marker must be present (asserted directly — it is NOT masked
	// for its own presence check); the behavioral required phrases below are
	// checked with the marker MASKED so a marker-substring cannot satisfy them.
	if !strings.Contains(step10Behavioral, terminalTeardownMarker) {
		t.Errorf("shared-core step 10 missing the verbatim terminal-status marker %q", terminalTeardownMarker)
	}
	assertDirectionalMandate(t, "shared-core step 10", step10Behavioral, negatingPhrases,
		[]string{
			"bounded best-effort",          // the bounded framing (not retry-to-success)
			"attempt cap",                  // bounded fast retries
			"wait a short settle interval", // inter-attempt settle (load-bearing, not delegated)
			// STOP/hold semantics that the marker substring CANNOT satisfy (M1 fix):
			// the FO STOPS the attempts and makes NO FURTHER teardown calls. A
			// near-synonym retry-to-success ("continues issuing … rounds until …
			// clears") cannot keep these.
			"STOPS the teardown attempts",    // the bound terminates the attempts
			"no further teardown tool calls", // the hold: nothing more is issued
			// Launcher-owned exit, distinct from the marker's "for launcher." (M2 fix):
			// the prose must assert the launcher owns the exit AND it is not the FO's.
			// A `kill -9 $PPID` self-exit that demotes the launcher to a "backstop"
			// cannot keep this.
			"launcher's** responsibility, not the FO's", // launcher owns the exit, not the FO
		})

	// Claude runtime: the Awaiting-Completion section already bans retrying
	// TeamDelete BEFORE the completion signal (the pre-completion wait phase). The
	// terminal teardown is a DIFFERENT phase: a BOUNDED set of TeamDelete attempts
	// with an inter-attempt settle, then the verbatim marker + hold, with the exit
	// owned by the launcher. Scope to the `## Terminal Team Teardown` section. The
	// Claude runtime IS the adapter realization, so the behavioral phrases are
	// asserted directly.
	claude := files["first-officer/references/claude-first-officer-runtime.md"]
	teardown := sectionAfter(claude, "## Terminal Team Teardown")
	if teardown == "" {
		t.Fatal("Claude FO runtime missing the `## Terminal Team Teardown` section (the bounded-best-effort realization of shared-core step 10)")
	}
	if !strings.Contains(teardown, terminalTeardownMarker) {
		t.Errorf("Terminal Team Teardown missing the verbatim terminal-status marker %q", terminalTeardownMarker)
	}
	assertDirectionalMandate(t, "Terminal Team Teardown", teardown, negatingPhrases,
		[]string{
			"TeamDelete",          // names the call it attempts
			"active member",       // names the race it attempts past
			"wait for the settle", // inter-attempt settle mandate
			"sleep 2",             // the concrete settle
			"attempt cap",         // bounded attempts
			// STOP/hold semantics the marker substring CANNOT satisfy (M1 fix): on
			// cap-exhaustion the FO STOPS calling TeamDelete and makes NO FURTHER
			// teardown calls. A "continues issuing … rounds until … clears"
			// near-synonym cannot keep these.
			"STOPS calling",                  // the bound terminates the TeamDelete calls
			"no further teardown tool calls", // the hold: nothing more is issued
			// Launcher-owned exit, distinct from the marker's "for launcher." (M2 fix):
			// the prose must explicitly forbid an FO self-exit. A `kill -9 $PPID`
			// self-exit cannot keep "never an FO self-exit".
			"never an FO self-exit", // the exit is the launcher's, never the FO's
		})
}

// assertDirectionalMandate fails if any negating (inverted-mandate) phrase is
// present in region, or if any required positive-mandate phrase is missing. The
// two-sided check is what makes the lint inversion-resistant: an inverted edit
// either introduces a negating phrase or drops a positive one (or both), so it
// cannot pass by merely preserving grep tokens.
//
// Required phrases are checked against the region with the terminal-status MARKER
// MASKED OUT (markerStripped). The marker string itself contains the substrings
// "holding" and "for launcher.", so a naive `strings.Contains(region, "hold")`
// or `…, "launcher")` is satisfied by the marker MENTION alone — they add ZERO
// discriminating power and let a gut-edit that keeps the verbatim marker but
// guts the STOP/hold/launcher-exit BEHAVIORAL prose pass green (the audit-cycle-1
// M1/M2 holes). Masking the marker forces the behavioral required phrases to be
// satisfied by the actionable prose, not the sentinel. The marker's own presence
// is asserted separately by the caller (it is a required phrase in its own right).
func assertDirectionalMandate(t *testing.T, label, region string, negating, required []string) {
	t.Helper()
	lower := strings.ToLower(region)
	for _, neg := range negating {
		if strings.Contains(lower, neg) {
			t.Errorf("%s contains the inverted-mandate phrase %q — the terminal teardown is BOUNDED best-effort then a marker + hold, NOT retry-to-success or an FO self-exit", label, neg)
		}
	}
	masked := strings.Replace(region, terminalTeardownMarker, " <marker> ", -1)
	for _, req := range required {
		if !strings.Contains(masked, req) {
			t.Errorf("%s missing required directional-mandate phrase %q (checked against prose with the marker masked, so a marker-substring does not count)", label, req)
		}
	}
}

// stripDelegationTail removes the trailing "…are the adapter's" delegation
// sentence from a teardown step. That sentence NAMES the load-bearing concepts
// (settle interval, cap value, marker) while delegating their realization to the
// adapter — so a required-phrase that matched only there would let a gut-edit of
// the BEHAVIORAL mandate pass green. Dropping it forces the behavioral phrases to
// be asserted against the actionable prose. If the delegation sentence is absent,
// the region is returned unchanged.
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
// premature-teardown bug). The terminal bounded-teardown clause is a separate
// phase; this lint ensures the Awaiting-Completion ban survives the reversal.
//
// Structural ban check, with an HONEST ceiling (audit-cycle-2 P2 correction): a
// bare substring-PRESENCE check on "emit `TeamDelete`" passes an INVERTED ban
// that keeps the token, and a CLOSED list of forbidden affirmative spellings
// passes any affirmative phrased OUTSIDE the list (e.g. "it is fine to go ahead
// and call `TeamDelete` right now"). This lint does two structural things that do
// NOT depend on an exact spelling:
//
//  1. Positive framing anchors — the ban must remain a `- emit `TeamDelete“
//     bullet under the `Do not:` prohibition, with the premature-teardown
//     rationale. Flipping the ban into a directive drops the bullet form or the
//     rationale.
//  2. Affirmative-permission CO-OCCURRENCE scan — any line that pairs a
//     permission cue (fine to / go ahead / should / may / ok to / feel free /
//     allowed to) with a `TeamDelete` reference is flagged, regardless of the
//     exact spelling. This catches the open class of "go ahead and call
//     TeamDelete" re-introductions, not just a fixed four.
//
// HONEST CEILING: prose lints are inherently reword-evadable — a sufficiently
// novel paraphrase of "you may tear down early" that uses none of the permission
// cues above will still pass. This lint raises the bar to "the common affirmative
// re-intros red AND the ban's positive framing is pinned"; the BEHAVIORAL oracle
// for the ban surviving is the live-e2e run (AC-1), not this structural lint.
func TestAwaitingCompletionStillBansPreCompletionTeamDelete(t *testing.T) {
	claude := vendoredSkillFiles(t)["first-officer/references/claude-first-officer-runtime.md"]
	region := sectionAfter(claude, "## Awaiting Completion")
	if region == "" {
		t.Fatal("Claude FO runtime missing the `## Awaiting Completion` section")
	}
	if !strings.Contains(region, "emit `TeamDelete`") {
		t.Error("Awaiting Completion must still ban emitting TeamDelete before the completion signal arrives")
	}
	// (1) Positive framing anchors.
	if !strings.Contains(region, "Do not:") {
		t.Error("Awaiting Completion must keep the negative `Do not:` framing that governs the TeamDelete ban")
	}
	if !strings.Contains(region, "- emit `TeamDelete`") {
		t.Error("the TeamDelete ban must remain a `- emit `TeamDelete`` bullet under the `Do not:` list, not a free-standing (potentially affirmative) mention")
	}
	if !strings.Contains(region, "tearing down is premature") {
		t.Error("the TeamDelete ban must keep its premature-teardown rationale (the semantic that pins it as a ban, not a directive)")
	}
	// (2) Affirmative-permission co-occurrence scan: any LINE pairing a permission
	// cue with a TeamDelete reference is an affirmative re-intro, regardless of the
	// exact wording. The legitimate `- emit `TeamDelete`` ban bullet carries NO
	// permission cue, so it is not flagged; an inverted "it is fine to go ahead and
	// call `TeamDelete` right now" line carries both and reds.
	permissionCues := []string{"fine to", "go ahead", "should", "may ", "ok to", "okay to", "feel free", "allowed to"}
	lower := strings.ToLower(region)
	for _, line := range strings.Split(lower, "\n") {
		if !strings.Contains(line, "teamdelete") {
			continue
		}
		for _, cue := range permissionCues {
			if strings.Contains(line, cue) {
				t.Errorf("Awaiting Completion line pairs a permission cue %q with TeamDelete — the pre-completion ban was inverted into an affirmative directive: %q", strings.TrimSpace(cue), strings.TrimSpace(line))
			}
		}
	}
}
