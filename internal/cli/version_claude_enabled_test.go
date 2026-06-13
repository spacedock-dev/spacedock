// ABOUTME: AC-4 oracle — claude spacedock-enablement reflects the plugin list
// ABOUTME: `enabled` boolean, not mere presence; a false entry is not "enabled".
package cli

import "testing"

// TestClaudeEnablementFromListJSON (AC-4) feeds the enablement reader a `claude
// plugin list --json` body with the spacedock entry's `enabled` set true vs false
// (and absent) and asserts the resolved enablement. The fixture JSON — shaped from
// the live spike capture (id + installPath + enabled) — is the independent source
// of truth: an `enabled:false` entry must NOT resolve to enabled.
func TestClaudeEnablementFromListJSON(t *testing.T) {
	const enabledTrue = `[{"id":"spacedock@spacedock","installPath":"/p/0.19.9","enabled":true},
		{"id":"other@market","installPath":"/q","enabled":true}]`
	const enabledFalse = `[{"id":"spacedock@spacedock","installPath":"/p/0.19.9","enabled":false}]`
	const noSpacedock = `[{"id":"other@market","installPath":"/q","enabled":true}]`

	cases := []struct {
		name string
		body string
		want enablement
	}{
		{"enabled-true", enabledTrue, enablementEnabled},
		{"enabled-false", enabledFalse, enablementNotEnabled},
		{"no-spacedock-entry", noSpacedock, enablementNotEnabled},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := claudeEnablement([]byte(tc.body))
			if err != nil {
				t.Fatalf("claudeEnablement(%s) err = %v, want nil", tc.name, err)
			}
			if got != tc.want {
				t.Fatalf("claudeEnablement(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// TestClaudeEnablementMalformedJSONErrors (AC-5 feeder) asserts a body that does
// not parse is an error, so the caller renders `enablement unknown` rather than
// silently downgrading to not-enabled.
func TestClaudeEnablementMalformedJSONErrors(t *testing.T) {
	if _, err := claudeEnablement([]byte("not json")); err == nil {
		t.Fatalf("claudeEnablement(malformed) err = nil, want a parse error so the caller renders enablement unknown")
	}
}
