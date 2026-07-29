// ABOUTME: Unit tests for the sandbox registry and the two render surfaces —
// ABOUTME: Inside matches on value, SessionState reports the session, LaunchState the launch.
package safehouse

import (
	"strings"
	"testing"
)

// fakeStateEnv backs Inside's getenv seam from a map, so no test reads the
// running machine's real environment.
func fakeStateEnv(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

// TestInsideMatchesValueNotPresence pins the registry's matching rule. The
// load-bearing row is the third: APP_SANDBOX_CONTAINER_ID is a generic macOS
// app-sandbox variable, so an implementation that tested mere presence would
// claim some other container as agent-safehouse — that row turns red.
func TestInsideMatchesValueNotPresence(t *testing.T) {
	cases := []struct {
		name     string
		vars     map[string]string
		wantName string
		wantOK   bool
	}{
		{"inside-safehouse", map[string]string{"APP_SANDBOX_CONTAINER_ID": "agent-safehouse"}, "agent-safehouse", true},
		{"unset", map[string]string{}, "", false},
		{"some-other-container-is-not-safehouse", map[string]string{"APP_SANDBOX_CONTAINER_ID": "com.example.other"}, "", false},
		{"empty-value", map[string]string{"APP_SANDBOX_CONTAINER_ID": ""}, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name, ok := Inside(fakeStateEnv(tc.vars))
			if name != tc.wantName || ok != tc.wantOK {
				t.Fatalf("Inside() = (%q, %v), want (%q, %v)", name, ok, tc.wantName, tc.wantOK)
			}
		})
	}
}

// TestSessionStateThreeWay pins the three strings `--version` and `status --boot`
// render. The regression row is `inside-safehouse-absent-from-PATH`: that is the
// live configuration this fix exists for — sandboxed, and safehouse off PATH
// precisely BECAUSE the wrap already happened — where the old renderer said
// `unavailable (safehouse not on PATH)`. The expected strings are test-supplied
// literals, not the production constants, so renaming a production string cannot
// silently take the assertion with it.
func TestSessionStateThreeWay(t *testing.T) {
	cases := []struct {
		name       string
		insideName string
		inside     bool
		available  bool
		want       string
	}{
		{"inside-safehouse-absent-from-PATH", "agent-safehouse", true, false, "inside (agent-safehouse)"},
		{"inside-safehouse-on-PATH", "agent-safehouse", true, true, "inside (agent-safehouse)"},
		{"not-sandboxed-available", "", false, true, "not sandboxed (safehouse available)"},
		{"not-sandboxed-not-installed", "", false, false, "not sandboxed (safehouse not installed)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SessionState(tc.insideName, tc.inside, tc.available)
			if got != tc.want {
				t.Fatalf("SessionState(%q, %v, %v) = %q, want %q", tc.insideName, tc.inside, tc.available, got, tc.want)
			}
			// AC-2b: a profile is a launch fact, and this surface is not about a
			// launch. Collapsing the two renderers back into one would put the
			// profile text here and turn this red.
			if strings.Contains(got, ".safehouse") {
				t.Fatalf("SessionState rendered %q, which mentions a .safehouse profile — a launch fact on a session surface", got)
			}
		})
	}
}

// TestLaunchStateFourWay pins the four strings the pre-launch banner renders.
// The inside arm dominates: launched from within the sandbox, safehouse is off
// PATH and a profile is typically present, which the old renderer called
// `unavailable` — the same inversion, on the banner.
func TestLaunchStateFourWay(t *testing.T) {
	cases := []struct {
		name       string
		insideName string
		inside     bool
		selected   bool
		available  bool
		want       string
	}{
		{
			"inside-dominates-even-with-profile-and-no-binary",
			"agent-safehouse", true, true, false,
			"inside (agent-safehouse) — launching without re-wrapping",
		},
		{
			"wrapping-this-launch",
			"", false, true, true,
			"wrapping this launch (safehouse, .safehouse profile)",
		},
		{
			"no-profile-not-wrapping",
			"", false, false, true,
			"not wrapping this launch (no .safehouse profile)",
		},
		{
			"profile-but-safehouse-not-installed",
			"", false, true, false,
			"not wrapped (safehouse not installed; .safehouse profile present)",
		},
		{
			"no-profile-and-not-installed",
			"", false, false, false,
			"not wrapping this launch (no .safehouse profile)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := LaunchState(tc.insideName, tc.inside, tc.selected, tc.available)
			if got != tc.want {
				t.Fatalf("LaunchState(%q, %v, %v, %v) = %q, want %q",
					tc.insideName, tc.inside, tc.selected, tc.available, got, tc.want)
			}
		})
	}
}
