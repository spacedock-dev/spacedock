// ABOUTME: AC-1 oracle — the launcher banner's Sandbox: line answers the LAUNCH
// ABOUTME: question over (already-inside, profile-present, binary-available).
package cli

import (
	"bytes"
	"testing"
)

// bannerEnv backs launchBanner's getenv seam from a map, so the already-inside-a-
// sandbox signal is pinned rather than read from the running machine.
func bannerEnv(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

// renderBannerSandbox exercises the real launchBanner for host from dir against a
// pinned getenv and lookPath stub (so both the inside-a-sandbox signal and the
// safehouse-binary availability are controlled, not read from the machine) and
// returns the rendered stderr bytes. `selected` mirrors the launcher's wrap
// decision (a profile present, or a --safehouse* flag).
func renderBannerSandbox(host, dir string, selected bool, getenv func(string) string, lookPath func(string) (string, error)) string {
	var buf bytes.Buffer
	launchBanner(host, dir, selected, getenv, lookPath, &buf)
	return buf.String()
}

// TestLaunchBannerSandboxLine (AC-1) drives launchBanner over the
// (inside, selected, available) input combinations and asserts the exact rendered
// Sandbox: line. The expected state strings are written here independently of the
// source, so a reverted-wording edit fails this assertion.
//
// The REGRESSION row is `inside-safehouse-absent-from-PATH`: launching from
// within the sandbox, safehouse is off PATH precisely because the wrap already
// happened, and the banner used to render `unavailable (safehouse not on PATH)`
// there. The banner keeps LAUNCH semantics — it is about to perform a launch, so
// the .safehouse profile is load-bearing here and only here.
func TestLaunchBannerSandboxLine(t *testing.T) {
	insideSafehouse := map[string]string{"APP_SANDBOX_CONTAINER_ID": "agent-safehouse"}
	outside := map[string]string{}

	cases := []struct {
		name     string
		vars     map[string]string
		selected bool
		lookPath func(string) (string, error)
		want     string
	}{
		{
			"inside-safehouse-absent-from-PATH", insideSafehouse, true, lookMissing,
			"Sandbox: inside (agent-safehouse) — launching without re-wrapping",
		},
		{"wrapping-this-launch", outside, true, lookFound, "Sandbox: wrapping this launch (safehouse, .safehouse profile)"},
		{"no-profile", outside, false, lookFound, "Sandbox: not wrapping this launch (no .safehouse profile)"},
		{
			"profile-but-safehouse-not-installed", outside, true, lookMissing,
			"Sandbox: not wrapped (safehouse not installed; .safehouse profile present)",
		},
		{"no-profile-and-not-installed", outside, false, lookMissing, "Sandbox: not wrapping this launch (no .safehouse profile)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := renderBannerSandbox("claude", t.TempDir(), tc.selected, bannerEnv(tc.vars), tc.lookPath)
			if !lineEquals(out, tc.want) {
				t.Fatalf("banner sandbox line = %q, want a whole line %q", out, tc.want)
			}
		})
	}
}
