// ABOUTME: Fixture test package proving the CI clean-output + jsonl-archive
// ABOUTME: shape; failures are gated by FIXTURE_FAIL so one binary runs green or red.
package cleanoutputfixture

import (
	"os"
	"testing"
)

// failOn reports whether the planted failures should fire. The AC test runs the
// same compiled binary twice — once with FIXTURE_FAIL unset (green) and once set
// (red) — so a single fixture exercises both the clean-green and failures-only-red
// surfaces without a second fixture package.
func failOn() bool { return os.Getenv("FIXTURE_FAIL") == "1" }

func TestAlphaPasses(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("arithmetic is broken")
	}
}

func TestGammaFails(t *testing.T) {
	if !failOn() {
		return
	}
	got := 7
	if got != 42 {
		t.Errorf("compute() = %d, want 42", got)
	}
}

func TestZetaSubtests(t *testing.T) {
	t.Run("case_good", func(t *testing.T) {})
	t.Run("case_bad", func(t *testing.T) {
		if !failOn() {
			return
		}
		t.Errorf("subtest assertion: got %q want %q", "x", "y")
	})
}
