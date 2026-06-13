// ABOUTME: AC oracle — claude's best-effort enabled/disabled marker reflects the
// ABOUTME: plugin list `enabled` boolean; an absent entry is unknown, not disabled.
package cli

import "testing"

// TestClaudeMarkerFromListJSON feeds the marker reader a `claude plugin list
// --json` body with the spacedock entry's `enabled` set true vs false (and absent)
// and asserts the resolved best-effort marker. The fixture JSON — shaped from the
// live spike capture (id + installPath + enabled) — is the independent source of
// truth: an `enabled:false` entry resolves to disabled, an absent entry to unknown
// (the version still comes from the manifest, so we never claim it is disabled).
func TestClaudeMarkerFromListJSON(t *testing.T) {
	const enabledTrue = `[{"id":"spacedock@spacedock","installPath":"/p/0.19.9","enabled":true},
		{"id":"other@market","installPath":"/q","enabled":true}]`
	const enabledFalse = `[{"id":"spacedock@spacedock","installPath":"/p/0.19.9","enabled":false}]`
	const noSpacedock = `[{"id":"other@market","installPath":"/q","enabled":true}]`

	cases := []struct {
		name string
		body string
		want enabledMarker
	}{
		{"enabled-true", enabledTrue, markerEnabled},
		{"enabled-false", enabledFalse, markerDisabled},
		{"no-spacedock-entry", noSpacedock, markerUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := claudeMarker([]byte(tc.body))
			if err != nil {
				t.Fatalf("claudeMarker(%s) err = %v, want nil", tc.name, err)
			}
			if got != tc.want {
				t.Fatalf("claudeMarker(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// TestClaudeMarkerMalformedJSONErrors asserts a body that does not parse is an
// error, so the caller renders the bare version (markerUnknown) rather than
// silently claiming the plugin is disabled.
func TestClaudeMarkerMalformedJSONErrors(t *testing.T) {
	if _, err := claudeMarker([]byte("not json")); err == nil {
		t.Fatalf("claudeMarker(malformed) err = nil, want a parse error so the caller renders the bare version")
	}
}
