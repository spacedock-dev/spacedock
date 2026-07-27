// ABOUTME: resolveBuildHost's refusal policy and error strings, pinned across
// ABOUTME: the extraction of the marker table into internal/runtimehost.
package dispatch

import (
	"strings"
	"testing"
)

func hostEnv(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

// TestResolveBuildHostFromMarkers pins that dispatch still REFUSES on ambiguity
// after the marker table moved to internal/runtimehost, where --version reports
// the same facts and exits 0. The two-pi-markers row is the one that would go red
// under a naive "more than one marker set" ambiguity rule; the codex+claude row
// is the recorded nested-session marker leak, and it must remain an error here
// even though it is a reported state on --version.
func TestResolveBuildHostFromMarkers(t *testing.T) {
	cases := []struct {
		name     string
		vars     map[string]string
		wantHost string
		wantErr  string
	}{
		{"codex", map[string]string{"CODEX_THREAD_ID": "t1"}, "codex", ""},
		{"claude", map[string]string{"CLAUDECODE": "1"}, "claude", ""},
		{
			"pi-two-markers-resolve-not-ambiguous",
			map[string]string{"PI_CODING_AGENT": "1", "PI_CODING_AGENT_DIR": "/home/u/.pi/agent"},
			"pi", "",
		},
		{
			"codex-and-claude-refused",
			map[string]string{"CODEX_THREAD_ID": "t1", "CLAUDECODE": "1"},
			"",
			"ambiguous runtime host sources: multiple runtime markers are set (CODEX_THREAD_ID, CLAUDECODE); pass --host claude, codex, or pi",
		},
		{
			"no-markers",
			map[string]string{},
			"",
			"missing host source: pass --host, set JSON host, or run under CODEX_THREAD_ID, CLAUDECODE, PI_CODING_AGENT, or PI_CODING_AGENT_DIR",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			host, err := resolveBuildHost("", "", hostEnv(tc.vars))
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("resolveBuildHost = error %v, want host %q", err, tc.wantHost)
				}
				if host != tc.wantHost {
					t.Fatalf("host = %q, want %q", host, tc.wantHost)
				}
				return
			}
			if err == nil {
				t.Fatalf("resolveBuildHost = %q, want error %q", host, tc.wantErr)
			}
			if err.Error() != tc.wantErr {
				t.Fatalf("error = %q, want %q", err.Error(), tc.wantErr)
			}
		})
	}
}

// TestResolveBuildHostExplicitWins pins that the explicit sources still short
// out marker detection entirely — the extraction touched only the marker branch.
func TestResolveBuildHostExplicitWins(t *testing.T) {
	env := hostEnv(map[string]string{"CODEX_THREAD_ID": "t1", "CLAUDECODE": "1"})
	host, err := resolveBuildHost("pi", "", env)
	if err != nil || host != "pi" {
		t.Fatalf("resolveBuildHost(--host=pi) = (%q, %v), want (\"pi\", nil) despite ambiguous markers", host, err)
	}
	if _, err := resolveBuildHost("bogus", "", env); err == nil || !strings.Contains(err.Error(), "unsupported host") {
		t.Fatalf("resolveBuildHost(--host=bogus) error = %v, want unsupported host", err)
	}
}
