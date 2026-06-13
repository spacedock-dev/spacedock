// ABOUTME: Unit tests for the shared three-way sandbox-state render strings the
// ABOUTME: launcher banner, status --boot, and --version all source from State.
package safehouse

import "testing"

// TestStateThreeWay pins the three rendered sandbox-state strings against their
// (selected, available) inputs. `selected` means a launch would be wrapped (a
// .safehouse profile is present, or a --safehouse* flag forced it); `available`
// means the safehouse binary resolves on PATH. The three surfaces (banner, boot,
// --version) all read these strings, so they are the single source of truth.
func TestStateThreeWay(t *testing.T) {
	cases := []struct {
		name      string
		selected  bool
		available bool
		want      string
	}{
		// available and selected → the launch will be wrapped through safehouse.
		{"enabled", true, true, "enabled (safehouse)"},
		// available but nothing selected it → the binary is installed, this launch
		// is not sandboxed.
		{"available-not-enabled", false, true, "available, not enabled (no .safehouse profile)"},
		// the binary is not on PATH → a present profile cannot take effect.
		{"unavailable", false, false, "unavailable (safehouse not on PATH)"},
		// the binary is absent even with a profile selected → still unavailable; the
		// missing binary dominates because nothing can wrap the launch.
		{"unavailable-even-when-selected", true, false, "unavailable (safehouse not on PATH)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := State(tc.selected, tc.available); got != tc.want {
				t.Fatalf("State(%v, %v) = %q, want %q", tc.selected, tc.available, got, tc.want)
			}
		})
	}
}
