// ABOUTME: Unit tests for env-only runtime-host detection: the marker matrix,
// ABOUTME: the same-host-two-markers non-ambiguity case, identity, and ShortID.
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
			host, markers, _, ambiguous := Detect(fakeEnv(tc.vars))
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

// TestDetectIdentity pins the identity column against the detection marker. The
// load-bearing row is claude: CLAUDECODE is the literal flag `1`, so an
// implementation that reuses the detection marker's VALUE as the identifier
// would render `session 1` for every claude session — that row asserts the
// identity is not "1". An ambiguous marker set carries no identity, since naming
// one would assert a host had been resolved.
func TestDetectIdentity(t *testing.T) {
	cases := []struct {
		name string
		vars map[string]string
		want string
	}{
		{
			"claude-reads-CLAUDE_CODE_SESSION_ID-not-CLAUDECODE",
			map[string]string{"CLAUDECODE": "1", "CLAUDE_CODE_SESSION_ID": "afd74765-9000-4e63-acf4-3b1f4645a8f3"},
			"afd74765-9000-4e63-acf4-3b1f4645a8f3",
		},
		{"claude-identity-unset", map[string]string{"CLAUDECODE": "1"}, ""},
		{
			"codex-marker-doubles-as-identity",
			map[string]string{"CODEX_THREAD_ID": "01937f2a-bbbb"},
			"01937f2a-bbbb",
		},
		{
			"pi-exposes-no-session-identity",
			map[string]string{"PI_CODING_AGENT": "1", "PI_CODING_AGENT_DIR": "/home/u/.pi/agent"},
			"",
		},
		{
			"ambiguous-carries-no-identity",
			map[string]string{"CLAUDECODE": "1", "CLAUDE_CODE_SESSION_ID": "afd74765", "CODEX_THREAD_ID": "01937f2a"},
			"",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, identity, _ := Detect(fakeEnv(tc.vars))
			if identity != tc.want {
				t.Fatalf("identity = %q, want %q", identity, tc.want)
			}
			if identity == "1" {
				t.Fatalf("identity = %q — the detection marker's value was reused as the session id", identity)
			}
		})
	}
}

// TestShortID pins the truncation rule: longer than the prefix truncates,
// exactly the prefix and shorter render whole (never padded), and empty stays
// empty so the caller omits the segment rather than printing `session `.
func TestShortID(t *testing.T) {
	cases := []struct {
		name string
		id   string
		want string
	}{
		{"longer-than-8", "afd74765-9000-4e63-acf4-3b1f4645a8f3", "afd74765"},
		{"exactly-8", "afd74765", "afd74765"},
		{"shorter-than-8", "ab12", "ab12"},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ShortID(tc.id); got != tc.want {
				t.Fatalf("ShortID(%q) = %q, want %q", tc.id, got, tc.want)
			}
		})
	}
}
