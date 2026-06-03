// ABOUTME: Structural lint over the FO terminal-teardown contract — step 10 and the
// ABOUTME: Claude runtime require team teardown to retry to success, not stop at the first failed TeamDelete.
package integration

import (
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
// Oracle honesty: this lint asserts the contract carries the retry clause; it
// does NOT prove the live FO obeys it. The behavioral oracle is the live-e2e CI
// run (AC-1) — a prose-only rule's ceiling is "the wording is present," and
// TeamDelete is a Claude-runtime tool with no spacedock-binary seam to gate. The
// lint guards against the clause being dropped or reverted in a future edit.
func TestTerminalTeardownRetriesToSuccess(t *testing.T) {
	files := vendoredSkillFiles(t)

	// Shared core, step 10 in Merge and Cleanup. The terminal-teardown step must
	// name a retry-to-success discipline for the team teardown so a single failed
	// teardown call does not strand the FO subprocess.
	core := files["first-officer/references/first-officer-shared-core.md"]
	region := sectionAfter(core, "## Merge and Cleanup")
	if region == "" {
		t.Fatal("FO shared core missing the `## Merge and Cleanup` section")
	}
	if !strings.Contains(region, "retry-to-success") {
		t.Error("Merge and Cleanup step 10 must require a `retry-to-success` team teardown, not stop at the first failed teardown call")
	}
	if !strings.Contains(region, "attempt cap") {
		t.Error("Merge and Cleanup step 10 must bound the retry by a small attempt cap, not loop unboundedly")
	}

	// Claude runtime: the Awaiting-Completion section already bans retrying
	// TeamDelete BEFORE the completion signal (the pre-completion wait phase). The
	// fix must distinguish that from the TERMINAL teardown phase, where TeamDelete
	// MUST be retried to success. Assert the runtime carries a terminal-teardown
	// retry clause naming TeamDelete and the active-member race.
	claude := files["first-officer/references/claude-first-officer-runtime.md"]
	teardown := sectionAfter(claude, "## Terminal Team Teardown")
	if teardown == "" {
		t.Fatal("Claude FO runtime missing the `## Terminal Team Teardown` section (the retry-to-success realization of shared-core step 10)")
	}
	if !strings.Contains(teardown, "TeamDelete") {
		t.Error("Terminal Team Teardown must name the `TeamDelete` call it retries")
	}
	if !strings.Contains(teardown, "active member") {
		t.Error("Terminal Team Teardown must name the `active member(s)` race it retries past")
	}
	if !strings.Contains(teardown, "retry") {
		t.Error("Terminal Team Teardown must require a retry (not end the turn on the first failed TeamDelete)")
	}

	// The retry must be BOUNDED — a small cap, not an unbounded loop — so a
	// genuinely stuck teardown still terminates rather than spinning forever.
	if !strings.Contains(teardown, "cap") {
		t.Error("Terminal Team Teardown retry must be bounded by a small cap")
	}
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
