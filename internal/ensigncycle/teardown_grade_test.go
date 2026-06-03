// ABOUTME: AC-2 offline grade — gradeTerminalTeardown greens ONLY on the
// ABOUTME: terminal-status marker + hold tail, and REDs both authentic bug recordings (pre-yy give-up + post-yy retry-loop).
package ensigncycle

import "testing"

// TestGradeGreensOnlyOnMarkerAndHold is AC-2: the offline grade greens ONLY on
// the fix shape (a bounded count of TeamDelete attempts → the contract-mandated
// terminal-status marker → a hold) and REDs BOTH authentic bug recordings. It
// replaces the cycle-1 fakeProc-never-exits tautology: the proof is "greens only
// on the marker+hold," demonstrated by REDDING both real bug streams. Sub-second,
// no model spend.
//
// The three fixtures are real captured stream-json (the live test's t.Log tee
// format), each failing/passing for the RIGHT reason:
//   - GREEN: sonnet_teamdelete_bounded_marker — 5 bounded TeamDelete attempts +
//     the marker + a hold tail (no further teardown tool_use). PASS.
//   - RED #1: sonnet_teamdelete_hang — the pre-yy give-up: ONE TeamDelete, fails
//     "active member(s)", FO ends the turn with NO retry and NO marker. RED
//     because the marker is absent. (This is why an attempt-COUNT grade fails: 1
//     attempt is ≤ cap, so "bounded-then-stopped" alone would wrongly GREEN it.)
//   - RED #2: sonnet_teamdelete_retryloop — the post-yy retry-loop from the real
//     failed CI run 26891717026: 6 TeamDelete attempts past the cap, NO marker,
//     never holds. RED because the marker is absent.
func TestGradeGreensOnlyOnMarkerAndHold(t *testing.T) {
	t.Run("green_bounded_marker_hold", func(t *testing.T) {
		lines := loadStreamFixture(t, "sonnet_teamdelete_bounded_marker.stream.jsonl")
		ok, reason := gradeTerminalTeardown(lines)
		if !ok {
			t.Errorf("the fix shape (bounded attempts + marker + hold) must PASS, got reason: %s", reason)
		}
	})

	t.Run("red1_pre_yy_give_up_no_marker", func(t *testing.T) {
		lines := loadStreamFixture(t, "sonnet_teamdelete_hang.stream.jsonl")
		ok, reason := gradeTerminalTeardown(lines)
		if ok {
			t.Error("the pre-yy give-up (one TeamDelete, no retry, NO marker) must RED — boundedness alone (1 ≤ cap) must not green it")
		}
		if reason == "" {
			t.Error("a RED grade must carry a localizing reason")
		}
	})

	t.Run("red2_post_yy_retry_loop_no_marker", func(t *testing.T) {
		lines := loadStreamFixture(t, "sonnet_teamdelete_retryloop.stream.jsonl")
		ok, reason := gradeTerminalTeardown(lines)
		if ok {
			t.Error("the post-yy retry-loop (6 TeamDelete past the cap, NO marker, never holds) must RED")
		}
		if reason == "" {
			t.Error("a RED grade must carry a localizing reason")
		}
	})
}

// TestGradeRedsWhenTeardownContinuesAfterMarker pins the HOLD half of the grade
// independently of the marker half: even WITH the marker present, a teardown
// tool_use emitted AFTER the marker (the FO kept retrying instead of holding)
// must RED. Without this, a grade keyed on marker-PRESENCE alone would green an
// FO that emitted the marker and then kept hammering TeamDelete — which is not
// the bounded terminus the contract mandates. The fixture is the GREEN stream
// with one extra teardown line appended after the marker.
func TestGradeRedsWhenTeardownContinuesAfterMarker(t *testing.T) {
	lines := loadStreamFixture(t, "sonnet_teamdelete_bounded_marker.stream.jsonl")

	// Append a TeamDelete tool_use AFTER the marker — the FO did NOT hold.
	postMarkerRetry := `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"tu_post","name":"TeamDelete","input":{}}]}}`
	lines = append(lines, postMarkerRetry)

	ok, reason := gradeTerminalTeardown(lines)
	if ok {
		t.Error("a teardown tool_use AFTER the marker must RED — the FO must HOLD, not keep retrying")
	}
	if reason == "" {
		t.Error("a RED grade must carry a localizing reason")
	}
}
