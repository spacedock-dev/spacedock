// ABOUTME: Unit table for the stable-regression gate predicate, including the
// ABOUTME: `>=` boundary that lets a re-run of the current stable release pass.
package release

import "testing"

// TestEvaluateStableRegressionGate pins the decision for each input class that
// release.yml can give the gate. The equality row is the load-bearing one. The
// gate blocks a STRICTLY older tag only, so a re-run of the release that
// `stable` points at stays green. A change of the comparison to `<= 0` reds
// that row. The two error rows keep the gate loud on input that it cannot read.
// A silent pass publishes the regression that the gate must stop.
func TestEvaluateStableRegressionGate(t *testing.T) {
	for _, tc := range []struct {
		name     string
		tag      string
		stable   string
		wantPass bool
		wantErr  bool
	}{
		{"older patch on an older line blocks", "v0.27.1", "0.28.0", false, false},
		{"older minor blocks", "v0.26.9", "0.27.0", false, false},
		{"older tag against a prerelease stable blocks", "v0.27.0", "0.28.0-pre0", false, false},
		{"newer minor passes", "v0.28.0", "0.27.0", true, false},
		{"newer patch on the current line passes", "v0.27.1", "0.27.0", true, false},
		{"equal version passes, because a re-run must stay green", "v0.28.0", "0.28.0", true, false},
		{"prerelease tag errors", "v0.28.0-pre1", "0.27.0", false, true},
		{"unparseable tag errors", "v-broken", "0.27.0", false, true},
		{"unparseable stable version errors", "v0.28.0", "not-a-version", false, true},
		{"empty stable version errors", "v0.28.0", "", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dec, err := EvaluateStableRegressionGate(tc.tag, tc.stable)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("EvaluateStableRegressionGate(%q, %q) = %+v, want an error", tc.tag, tc.stable, dec)
				}
				return
			}
			if err != nil {
				t.Fatalf("EvaluateStableRegressionGate(%q, %q): %v", tc.tag, tc.stable, err)
			}
			if dec.Pass != tc.wantPass {
				t.Fatalf("EvaluateStableRegressionGate(%q, %q).Pass = %v, want %v (%s)",
					tc.tag, tc.stable, dec.Pass, tc.wantPass, dec.Reason)
			}
			if dec.Reason == "" {
				t.Fatal("decision carries no Reason for the step log")
			}
		})
	}
}
