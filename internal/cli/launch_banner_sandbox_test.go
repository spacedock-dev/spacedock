// ABOUTME: AC-1 oracle — the launcher banner carries a Sandbox: line whose text
// ABOUTME: matches the three-way state from (profile-present, binary-available).
package cli

import (
	"bytes"
	"testing"
)

// renderBannerSandbox exercises the real launchBanner for host from dir against a
// pinned lookPath stub (so the safehouse-binary availability is controlled, not
// read from the machine PATH) and returns the rendered stderr bytes. `selected`
// mirrors the launcher's wrap decision (a profile present, or a --safehouse* flag).
func renderBannerSandbox(host, dir string, selected bool, lookPath func(string) (string, error)) string {
	var buf bytes.Buffer
	launchBanner(host, dir, selected, lookPath, &buf)
	return buf.String()
}

// TestLaunchBannerSandboxLine (AC-1) drives launchBanner over each of the three
// (selected, available) input combinations and asserts the exact rendered Sandbox:
// line. The expected state strings are written here independently of the source,
// so a reverted-wording edit fails this assertion. `selected` is the launcher's
// own wrap decision; availability is pinned via lookFound / lookMissing.
func TestLaunchBannerSandboxLine(t *testing.T) {
	cases := []struct {
		name     string
		selected bool
		lookPath func(string) (string, error)
		want     string
	}{
		{"enabled", true, lookFound, "Sandbox: enabled (safehouse)"},
		{"available-not-enabled", false, lookFound, "Sandbox: available, not enabled (no .safehouse profile)"},
		{"unavailable", false, lookMissing, "Sandbox: unavailable (safehouse not on PATH)"},
		{"unavailable-even-when-selected", true, lookMissing, "Sandbox: unavailable (safehouse not on PATH)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := renderBannerSandbox("claude", t.TempDir(), tc.selected, tc.lookPath)
			if !lineEquals(out, tc.want) {
				t.Fatalf("banner sandbox line = %q, want a whole line %q", out, tc.want)
			}
		})
	}
}
