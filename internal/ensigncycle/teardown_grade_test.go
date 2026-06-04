// ABOUTME: AC-2 offline grade — gradeTerminalTeardown greens on FO-authored
// ABOUTME: terminal-status marker emission (incl. the real re-invoke-retry shape), REDs both authentic bug recordings.
package ensigncycle

import "testing"

// TestGradeGreensOnMarkerEmission is AC-2: the offline grade greens when the FO
// EMITS the contract-mandated terminal-status marker (reaching the bounded
// terminus) and REDs BOTH authentic bug recordings, which emit no marker. The
// grade keys on the FO AUTHORING the marker — NOT on a clean post-marker hold,
// which the PR #285 live run proved the real sonnet FO cannot deliver (it resumes
// teardown on each harness re-invoke). Sub-second, no model spend.
//
// Fixtures (real captured stream-json, the live test's t.Log tee format):
//   - GREEN: sonnet_teamdelete_bounded_marker — bounded TeamDelete attempts +
//     the FO emits the marker. PASS.
//   - GREEN: sonnet_teamdelete_marker_continues — the REAL failing PR #285 sonnet
//     run: the FO emits the marker (in a thinking block) AND continues teardown
//     across re-invokes. This is the producer reality a hold-grade red'd; marker
//     emission must PASS it.
//   - RED #1: sonnet_teamdelete_hang — the pre-yy give-up: ONE TeamDelete, fails
//     "active member(s)", FO ends the turn with NO retry and NO marker. RED.
//   - RED #2: sonnet_teamdelete_retryloop — the post-yy retry-loop (real run
//     26891717026): 6 TeamDelete past the cap, NO marker. RED.
func TestGradeGreensOnMarkerEmission(t *testing.T) {
	t.Run("green_bounded_marker", func(t *testing.T) {
		lines := loadStreamFixture(t, "sonnet_teamdelete_bounded_marker.stream.jsonl")
		ok, reason := gradeTerminalTeardown(lines)
		if !ok {
			t.Errorf("the fix shape (bounded attempts + the FO emits the marker) must PASS, got reason: %s", reason)
		}
	})

	t.Run("green_real_marker_continues", func(t *testing.T) {
		// The real PR #285 sonnet run: marker emitted, then teardown CONTINUES on
		// re-invoke (no clean hold). The grade keys on emission, so this PASSES —
		// the exact shape the prior hold-grade wrongly red'd against the live FO.
		lines := loadStreamFixture(t, "sonnet_teamdelete_marker_continues.stream.jsonl")
		ok, reason := gradeTerminalTeardown(lines)
		if !ok {
			t.Errorf("the real producer (marker emitted + teardown continues on re-invoke) must PASS — a hold-grade over-fits the recordings; got reason: %s", reason)
		}
	})

	t.Run("red1_pre_yy_give_up_no_marker", func(t *testing.T) {
		lines := loadStreamFixture(t, "sonnet_teamdelete_hang.stream.jsonl")
		ok, reason := gradeTerminalTeardown(lines)
		if ok {
			t.Error("the pre-yy give-up (one TeamDelete, no retry, NO marker) must RED")
		}
		if reason == "" {
			t.Error("a RED grade must carry a localizing reason")
		}
	})

	t.Run("red2_post_yy_retry_loop_no_marker", func(t *testing.T) {
		lines := loadStreamFixture(t, "sonnet_teamdelete_retryloop.stream.jsonl")
		ok, reason := gradeTerminalTeardown(lines)
		if ok {
			t.Error("the post-yy retry-loop (6 TeamDelete past the cap, NO marker) must RED")
		}
		if reason == "" {
			t.Error("a RED grade must carry a localizing reason")
		}
	})
}

// TestGradeIgnoresMarkerInContractRead pins the EMISSION discriminator: the grade
// must key on the FO AUTHORING the marker, not on the marker appearing in a
// contract-Read tool_result. The FO Reads the marker-bearing contract files
// (shared-core step 10 / the Claude runtime) at startup — the real run shows the
// marker in TWO user/tool_result entries (contract reads) plus ONE assistant
// thinking emission. A stream whose ONLY marker is a contract-Read must RED; the
// authentic green fixtures pass because the FO ALSO authors the marker.
func TestGradeIgnoresMarkerInContractRead(t *testing.T) {
	// A user tool_result that carries the marker verbatim (a contract Read) but no
	// assistant-authored marker → RED.
	contractRead := `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"r","content":[{"type":"text","text":"10. ...emit TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher. ..."}]}]}}`
	teamDelete := `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"tu","name":"TeamDelete","input":{}}]}}`
	lines := []string{contractRead, teamDelete}

	ok, reason := gradeTerminalTeardown(lines)
	if ok {
		t.Error("a marker that appears ONLY in a contract-Read tool_result must RED — the FO did not author it")
	}
	if reason == "" {
		t.Error("a RED grade must carry a localizing reason")
	}
}
