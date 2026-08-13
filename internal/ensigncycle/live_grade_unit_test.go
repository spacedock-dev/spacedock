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
