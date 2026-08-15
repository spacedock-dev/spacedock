// ABOUTME: Unit tests for env-only runtime-host detection: the marker matrix
// ABOUTME: and the same-host-two-markers non-ambiguity case.
package runtimehost

import (
	"strings"
	"testing"
)

// fakeEnv backs Detect's getenv seam from a map, so no test reads the running
// machine's real environment.
func fakeEnv(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

// TestDetectMarkerMatrix pins (host, markers, ambiguous) over the marker sets
// that matter. The two rows that carry the design: two markers of the SAME host
// (pi) are NOT ambiguous — an ambiguity rule of "more than one marker set" turns
// that row red — and two markers of DIFFERENT hosts are ambiguous with no host
// named.
func TestDetectMarkerMatrix(t *testing.T) {
	cases := []struct {
		name          string
		vars          map[string]string
		wantHost      string
		wantMarkers   string
		wantAmbiguous bool
	}{
		{"none", map[string]string{}, "", "", false},
		{"codex", map[string]string{"CODEX_THREAD_ID": "01937f2a-aaaa"}, "codex", "CODEX_THREAD_ID", false},
		{"claude", map[string]string{"CLAUDECODE": "1"}, "claude", "CLAUDECODE", false},
		{"pi-one-marker", map[string]string{"PI_CODING_AGENT": "1"}, "pi", "PI_CODING_AGENT", false},
		{
			"pi-two-markers-same-host-not-ambiguous",
			map[string]string{"PI_CODING_AGENT": "1", "PI_CODING_AGENT_DIR": "/home/u/.pi/agent"},
			"pi", "PI_CODING_AGENT,PI_CODING_AGENT_DIR", false,
		},
		{
			"codex-and-claude-ambiguous",
			map[string]string{"CODEX_THREAD_ID": "01937f2a", "CLAUDECODE": "1"},
			"", "CODEX_THREAD_ID,CLAUDECODE", true,
		},
		{
			"all-markers-ambiguous",
			map[string]string{
				"CODEX_THREAD_ID": "01937f2a", "CLAUDECODE": "1",
				"PI_CODING_AGENT": "1", "PI_CODING_AGENT_DIR": "/home/u/.pi/agent",
			},
			"", "CODEX_THREAD_ID,CLAUDECODE,PI_CODING_AGENT,PI_CODING_AGENT_DIR", true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			host, markers, ambiguous := Detect(fakeEnv(tc.vars))
			if host != tc.wantHost {
				t.Errorf("host = %q, want %q", host, tc.wantHost)
			}
			if got := strings.Join(markers, ","); got != tc.wantMarkers {
				t.Errorf("markers = %q, want %q", got, tc.wantMarkers)
			}
			if ambiguous != tc.wantAmbiguous {
				t.Errorf("ambiguous = %v, want %v", ambiguous, tc.wantAmbiguous)
			}
		})
	}
}
