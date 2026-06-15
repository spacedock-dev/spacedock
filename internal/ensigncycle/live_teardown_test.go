//go:build live

// ABOUTME: Team-FORCED live teardown coverage — a headless FO told to run in team
// ABOUTME: mode must reach the bounded terminal teardown and EMIT the TERMINAL_TEARDOWN_BOUNDED marker.
package ensigncycle

import (
	"testing"
)

// TestLiveEnsignCycleTeamTeardown is the team-FORCED companion to TestLiveEnsignCycle.
// The default cycle is TEAM-AGNOSTIC: it gates on dispatch-close + the on-disk
// terminal end-state, which hold whether the FO drove team or bare, so it must NOT
// gate on the team-only TERMINAL_TEARDOWN_BOUNDED marker (a legitimate bare drive
// never emits it — that is the relocated coin the team-agnostic retarget removes).
// But the bounded best-effort terminal-teardown path (AC-2) is real coverage worth
// keeping, and it ONLY runs when there is a team to tear down. So this test FORCES
// team mode via an explicit drivePrompt cue and keeps the teardown-marker grade —
// the marker coverage lives HERE, where a team is guaranteed, not on the default
// path where the mode is the FO's choice.
//
// The team-forcing cue is a legitimate operator instruction (the captain can tell
// the FO to use team mode for concurrent dispatch — `using-claude-team` step 1
// runs TeamCreate first whenever team mode is in effect). It is GENERIC: it names
// no stage or task, only the dispatch mode. With CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
// in the env (the live job and the operator re-run instructions both set it) and a
// fetchable TeamCreate, the FO creates a team, dispatches the ensign through it,
// terminalizes, then runs the bounded teardown and emits the marker.
func TestLiveEnsignCycleTeamTeardown(t *testing.T) {
	// FORCE team mode: the conn-cue (so the gateless fixture drives to terminal)
	// PLUS an explicit team-mode instruction so the FO creates a team rather than
	// taking the bare single-entity path. antiShutdownOverride still rides along —
	// it fights the per-turn #55297 teardown nag, which is exactly the team path
	// this test exercises.
	drivePrompt := "Run in team mode (create a team for concurrent dispatch). " +
		"Drive the workflow to completion; you have the conn to resolve gates from each stage report's verdict (auto-approve). " +
		antiShutdownOverride
	watcher, root := startRealisticLifecycleDrive(t, drivePrompt)

	// Step 1: the first ensign dispatch OPENS — the cycle progressed past boot into
	// dispatch. The barrier is the dispatch OPEN, not its close: in this Claude Code
	// version the team-mode ensign completion arrives as a `direct` message, not the
	// `task_notification status=completed` anchor expectDispatchClose keys on, so a
	// healthy team run can leave the dispatch "open" by that anchor's reckoning even
	// after the ensign finished. The real proof that the FULL cycle ran is Step 2's
	// terminal-teardown MARKER (emitted only after terminalize+archive+teardown), so
	// the dispatch OPEN is a sufficient and reliable early progress beat here; the
	// marker is the load-bearing one.
	if _, err := watcher.expect(isEnsignDispatch, quietBudgetDispatchClose, "ensign dispatch open"); err != nil {
		if wrongRoot := detectWrongRootBoot(watcher.fullTranscript(), root); wrongRoot != nil {
			t.Fatalf("live team-teardown cycle failed waiting for the ensign dispatch to open due to a wrong-root boot: %v\nUnderlying watcher error: %v", wrongRoot, err)
		}
		t.Fatalf("live team-teardown cycle failed waiting for the ensign dispatch to open: %v", err)
	}

	// Step 2 (the team-only coverage this test exists for): grade the BOUNDED
	// best-effort terminal teardown. The FO terminalizes + archives, then attempts
	// to tear the team down; the harness will not let claude -p exit while the
	// team's members[] is populated (the dead-but-listed member is never cleared —
	// upstream #38116/#57681), so a clean self-exit is impossible. The FIX's unique
	// signal is the FO EMITTING the contract-mandated TERMINAL_TEARDOWN_BOUNDED
	// MARKER (a text/thinking block it authors, not a contract-Read it merely saw).
	// expectTerminalTeardownGrade keys PASS on that marker emission — NOT on the
	// shutdown_request/TeamDelete beats both bug shapes also emit, and NOT on a
	// post-marker HOLD the real sonnet FO cannot deliver (it RESUMES teardown on
	// each harness re-invoke). On marker emission the t.Cleanup(poller.kill) reaps
	// the still-running subprocess and the cycle PASSES. The budget stays ≤60s so
	// the AC-1 timeout guard is unaffected.
	if err := watcher.expectTerminalTeardownGrade(quietBudgetDefault); err != nil {
		t.Fatalf("live team-teardown cycle failed grading the terminal teardown: %v", err)
	}

	// The terminal end-state must still be PRESENT and CORRECT on disk — the marker
	// alone is not proof the cycle completed, only that the FO reached the teardown
	// terminus. A team FO archives BEFORE the teardown hold, so by the time the
	// marker is emitted the entity is terminalized + archived. These are the same
	// team-INDEPENDENT end-state checks the default cycle runs; asserting them here
	// too keeps the team path honest (a team that emits the marker but left the
	// entity un-terminalized is a real regression).
	entity, where, found := locateEntity(root, "make-it-work")
	if !found {
		t.Fatalf("entity make-it-work not found in place or under _archive/ after the team-teardown cycle")
	}
	t.Logf("located entity at %s", where)
	if !liveStageReportHeading.MatchString(entity) {
		t.Errorf("entity missing anchored stage-report heading\n%s", entity)
	}
	if !frontmatterField.MatchString(entity) {
		t.Errorf("entity missing terminal `status: done`\n%s", entity)
	}
	if !verdictSet.MatchString(entity) {
		t.Errorf("entity missing a finalized (non-empty) `verdict:`\n%s", entity)
	}
	if !someCommitNamesOnly(t, root, "make-it-work") {
		t.Errorf("no path-scoped commit named only the entity in the team-teardown cycle history")
	}
}
