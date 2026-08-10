package ensigncycle

import (
	"reflect"
	"testing"
)

func TestGradeLiveStrictMatrix(t *testing.T) {
	expected := "implementation-worker-not-dispatched"
	for name, test := range map[string]struct {
		expected string
		errs     []error
		status   string
		codes    []string
	}{
		"pass":       {status: "pass"},
		"xfail":      {expected, []error{&gradedErr{code: expected}}, "xfail", []string{expected}},
		"xpass":      {expected, nil, "xpass", nil},
		"different":  {expected, []error{&gradedErr{code: "gate-not-held"}}, "fail", []string{"gate-not-held"}},
		"additional": {expected, []error{&gradedErr{code: expected}, &gradedErr{code: "gate-not-held"}}, "fail", []string{"gate-not-held", expected}},
	} {
		t.Run(name, func(t *testing.T) {
			got := gradeLive(test.expected, test.errs...)
			if got.status != test.status || !reflect.DeepEqual(got.codes, test.codes) {
				t.Fatalf("grade = %#v, want status=%s codes=%v", got, test.status, test.codes)
			}
		})
	}
}

func TestGradeLiveRunsEveryAssertion(t *testing.T) {
	var calls int
	assert := func(code string) error { calls++; return &gradedErr{code: code} }
	gradeLive("first", assert("first"), assert("second"))
	if calls != 2 {
		t.Fatalf("assertion calls = %d, want 2", calls)
	}
}
