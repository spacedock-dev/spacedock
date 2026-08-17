package ensigncycle

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGradeLiveTargetMatrix(t *testing.T) {
	for name, test := range map[string]struct {
		xfail  bool
		errs   []error
		status string
		codes  []string
	}{
		"pass":           {status: "pass"},
		"semantic xfail": {true, []error{&gradedErr{code: "gate-not-held"}}, "xfail", []string{"gate-not-held"}},
		"many semantics": {true, []error{&gradedErr{code: "worker-order"}, &gradedErr{code: "gate-not-held"}}, "xfail", []string{"gate-not-held", "worker-order"}},
		"empty xpass":    {true, nil, "xpass", nil},
		"infrastructure": {true, []error{fmt.Errorf("state read failed")}, "fail", nil},
	} {
		t.Run(name, func(t *testing.T) {
			got := gradeLive(test.xfail, test.errs...)
			if got.status != test.status || !reflect.DeepEqual(got.codes, test.codes) {
				t.Fatalf("grade = %#v, want status=%s codes=%v", got, test.status, test.codes)
			}
		})
	}
}

func TestGradeLiveRunsEveryAssertion(t *testing.T) {
	var calls int
	assert := func(code string) error { calls++; return &gradedErr{code: code} }
	gradeLive(true, assert("first"), assert("second"))
	if calls != 2 {
		t.Fatalf("assertion calls = %d, want 2", calls)
	}
}

func TestLiveGradeLaneResult(t *testing.T) {
	for status, wantFail := range map[string]bool{
		"pass":  false,
		"xfail": false,
		"xpass": false,
		"fail":  true,
	} {
		if got := liveGradeFailsLane(status); got != wantFail {
			t.Errorf("liveGradeFailsLane(%q) = %t, want %t", status, got, wantFail)
		}
	}
}

func TestFilingSemanticFailureUsesTargetXFail(t *testing.T) {
	err := &gradedErr{code: "filing-command-not-observed", msg: "filing command log has no spacedock new invocation"}
	for name, test := range map[string]struct {
		xfail bool
		want  string
	}{
		"bound":   {true, "xfail"},
		"unbound": {false, "fail"},
	} {
		t.Run(name, func(t *testing.T) {
			got := gradeLive(test.xfail, err)
			if got.status != test.want || !reflect.DeepEqual(got.codes, []string{"filing-command-not-observed"}) {
				t.Fatalf("grade = %#v, want status=%s code=filing-command-not-observed", got, test.want)
			}
		})
	}
}

// TestGradeLiveRetainsFindingMessages pins AC-3's second half. Every graded
// finding constructs a message explaining what it saw, and until now gradeLive
// discarded all of them — a CI failure printed a bare code list, and the durable
// end state that produced it lives in a t.TempDir that is gone by the time anyone
// reads the run. Each observed code must carry its own message through, paired
// with that code and ordered with it. Dropping the msg field, or pairing a message
// with the wrong code, fails this.
func TestGradeLiveRetainsFindingMessages(t *testing.T) {
	grade := gradeLive(false,
		&gradedErr{code: "rejection-round-missing", msg: "resolved launcher never invoked `gate record --round validation/1`"},
		&gradedErr{code: "rejection-gate-not-prepared", msg: "FO never prepared the cycle-2 validation gate: entity has no gates record"},
	)
	want := []string{
		"rejection-gate-not-prepared: FO never prepared the cycle-2 validation gate: entity has no gates record",
		"rejection-round-missing: resolved launcher never invoked `gate record --round validation/1`",
	}
	if !reflect.DeepEqual(grade.details, want) {
		t.Fatalf("details = %#v, want %#v", grade.details, want)
	}
	if !reflect.DeepEqual(grade.codes, []string{"rejection-gate-not-prepared", "rejection-round-missing"}) {
		t.Fatalf("codes = %v, want the details' codes in the same order", grade.codes)
	}

	// A repeated code keeps the finding that first reported it; a messageless
	// finding contributes a code but no empty detail line.
	repeated := gradeLive(false,
		&gradedErr{code: "gate-not-held", msg: "first observation"},
		&gradedErr{code: "gate-not-held", msg: "second observation"},
		&gradedErr{code: "worker-order"},
	)
	if !reflect.DeepEqual(repeated.details, []string{"gate-not-held: first observation"}) {
		t.Fatalf("repeated/messageless details = %#v", repeated.details)
	}

	// durableSemantic is the wrapper every scenario assertion goes through; the
	// assertion's own error text must be what survives.
	wrapped := gradeLive(false, durableSemantic("rejection-flow-state", fmt.Errorf("fix marker absent from the corrected candidate")))
	if !reflect.DeepEqual(wrapped.details, []string{"rejection-flow-state: fix marker absent from the corrected candidate"}) {
		t.Fatalf("durableSemantic details = %#v", wrapped.details)
	}
}
