// ABOUTME: --version reports the session it is actually in — the Runtime line
// ABOUTME: with its session id, the session sandbox state, and no host CLI exec.
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// versionEnv backs printVersion's getenv seam from a map, so no test here reads
// the running machine's real environment.
func versionEnv(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

// renderSession runs the real printVersion against a pinned environment and a
// pinned safehouse lookPath, returning the whole rendered output.
func renderSession(vars map[string]string, lookPath func(string) (string, error)) string {
	var buf bytes.Buffer
	printVersion(&buf, versionEnv(vars), lookPath)
	return buf.String()
}

// claudeSession is the live configuration this task was raised from: a claude
// session running inside the safehouse sandbox.
var claudeSession = map[string]string{
	"CLAUDECODE":               "1",
	"CLAUDE_CODE_SESSION_ID":   "afd74765-9000-4e63-acf4-3b1f4645a8f3",
	"APP_SANDBOX_CONTAINER_ID": "agent-safehouse",
}

// TestVersionSessionRender pins the whole `--version` render per marker set. The
// expected outputs are in-test literals, independent of the production strings.
//
// The organising rule under test: inside a session, report the session; outside
// one, report the version. Both shapes start with the same TWO lines — the
// version line and the `OS: <goos>/<goarch>` line — so the outside case is two
// lines (no Runtime line, no Sandbox line, and no contract token), while every
// in-session case adds exactly three lines below them.
func TestVersionSessionRender(t *testing.T) {
	version := displayVersion()
	osToken := runtime.GOOS + "/" + runtime.GOARCH
	cases := []struct {
		name     string
		vars     map[string]string
		lookPath func(string) (string, error)
		want     string
	}{
		{
			// The live regression case: sandboxed, and safehouse off PATH precisely
			// BECAUSE the wrap already happened. This rendered
			// `Sandbox: unavailable (safehouse not on PATH)` before this change.
			name:     "claude-inside-safehouse",
			vars:     claudeSession,
			lookPath: lookMissing,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Runtime: claude (CLAUDECODE, session afd74765)\n" +
				"Sandbox: inside (agent-safehouse)\n" +
				"contract 3\n",
		},
		{
			name:     "outside-every-runtime-is-two-lines",
			vars:     map[string]string{},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n",
		},
		{
			// A host whose identity variable is unset takes the same path as a host
			// with no identity variable at all: the segment is omitted. A stray
			// `, session ` or `()` fails this exact-match.
			name:     "claude-without-session-id",
			vars:     map[string]string{"CLAUDECODE": "1"},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Runtime: claude (CLAUDECODE)\n" +
				"Sandbox: not sandboxed (safehouse available)\n" +
				"contract 3\n",
		},
		{
			// Pi exposes no per-session identifier at all (PI_CODING_AGENT_SESSION_DIR
			// is the sessions COLLECTION), and two markers of the same host are not
			// ambiguity.
			name: "pi-two-markers-no-identity",
			vars: map[string]string{
				"PI_CODING_AGENT":     "1",
				"PI_CODING_AGENT_DIR": "/home/u/.pi/agent",
			},
			lookPath: lookMissing,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)\n" +
				"Sandbox: not sandboxed (safehouse not installed)\n" +
				"contract 3\n",
		},
		{
			// codex's detection marker doubles as its identity variable — the only
			// host where they coincide.
			name:     "codex-marker-doubles-as-identity",
			vars:     map[string]string{"CODEX_THREAD_ID": "01937f2a-bbbb-cccc"},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Runtime: codex (CODEX_THREAD_ID, session 01937f2a)\n" +
				"Sandbox: not sandboxed (safehouse available)\n" +
				"contract 3\n",
		},
		{
			// AC-3: ambiguity is REPORTED, not guessed at, and carries no session
			// identifier — printing one would assert a host had been resolved.
			name: "ambiguous-markers",
			vars: map[string]string{
				"CODEX_THREAD_ID":        "01937f2a",
				"CLAUDECODE":             "1",
				"CLAUDE_CODE_SESSION_ID": "afd74765-9000",
			},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE) — pass --host\n" +
				"Sandbox: not sandboxed (safehouse available)\n" +
				"contract 3\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderSession(tc.vars, tc.lookPath)
			if got != tc.want {
				t.Fatalf("--version render mismatch:\n got=%q\nwant=%q", got, tc.want)
			}
			// AC-2a: no empty or dangling parenthetical is ever reachable.
			if strings.Contains(got, "()") || strings.Contains(got, "session )") {
				t.Fatalf("--version rendered an empty parenthetical: %q", got)
			}
			// AC-2b: the session sandbox line never mentions a launch fact.
			if strings.Contains(got, ".safehouse") {
				t.Fatalf("--version mentions a .safehouse profile — a launch fact on a session surface: %q", got)
			}
			// AC-4: line 1 is the version token and nothing after it.
			line1 := strings.SplitN(got, "\n", 2)[0]
			if !regexp.MustCompile(`^spacedock \S+$`).MatchString(line1) {
				t.Fatalf("line 1 = %q, want to match ^spacedock \\S+$", line1)
			}
			if line1 != "spacedock "+version {
				t.Fatalf("line 1 = %q, want %q", line1, "spacedock "+version)
			}
			// AC (OS line): line 2 is always the OS/arch line, in both shapes.
			line2 := strings.SplitN(got, "\n", 3)[1]
			if !regexp.MustCompile(`^OS: [a-z0-9]+/[a-z0-9]+$`).MatchString(line2) {
				t.Fatalf("line 2 = %q, want `OS: <goos>/<goarch>`", line2)
			}
			if line2 != "OS: "+osToken {
				t.Fatalf("line 2 = %q, want %q", line2, "OS: "+osToken)
			}
		})
	}
}

// TestVersionContractTokenPlacement (AC-4) pins the moved sentinel: `contract 3`
// appears BELOW the sandbox line inside a session, never on line 1, and does not
// appear at all outside every runtime (every integer-era reader is itself a
// session).
func TestVersionContractTokenPlacement(t *testing.T) {
	inSession := renderSession(claudeSession, lookMissing)
	lines := strings.Split(strings.TrimRight(inSession, "\n"), "\n")
	if lines[len(lines)-1] != "contract 3" {
		t.Fatalf("in-session last line = %q, want %q", lines[len(lines)-1], "contract 3")
	}
	if strings.Contains(lines[0], "contract") {
		t.Fatalf("line 1 = %q, want the contract token moved off it", lines[0])
	}

	outside := renderSession(map[string]string{}, lookFound)
	if strings.Contains(outside, "contract") {
		t.Fatalf("outside-every-runtime render = %q, want no contract token", outside)
	}
}

// TestVersionRuntimeLineDistinguishesConcurrentSessions (AC-2) is the point of
// printing an identifier at all: two concurrent sessions on the same host must
// render DIFFERENT Runtime lines. An implementation that reused the detection
// marker's value would render `session 1` for both and turn this red.
func TestVersionRuntimeLineDistinguishesConcurrentSessions(t *testing.T) {
	runtimeLineOf := func(sessionID string) string {
		out := renderSession(map[string]string{
			"CLAUDECODE":             "1",
			"CLAUDE_CODE_SESSION_ID": sessionID,
		}, lookFound)
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(line, "Runtime: ") {
				return line
			}
		}
		t.Fatalf("no Runtime: line in %q", out)
		return ""
	}

	first := runtimeLineOf("afd74765-9000-4e63-acf4-3b1f4645a8f3")
	second := runtimeLineOf("afd74799-9000-4e63-acf4-3b1f4645a8f3")
	if first == second {
		t.Fatalf("two sessions rendered the same Runtime line %q — the identifier does not distinguish them", first)
	}
	if first != "Runtime: claude (CLAUDECODE, session afd74765)" {
		t.Fatalf("first Runtime line = %q", first)
	}
	if second != "Runtime: claude (CLAUDECODE, session afd74799)" {
		t.Fatalf("second Runtime line = %q", second)
	}
}

// TestVersionExecutesNoHostCLI (AC-2) proves the per-host block's removal is REAL
// rather than cosmetic. Executable `claude`/`codex`/`pi` shims are placed on a
// temp PATH that append to a witness file when run; `--version` is then driven
// through the real Run entry point. The assertion is that the witness file was
// never created. Before this change the shims got three entries.
func TestVersionExecutesNoHostCLI(t *testing.T) {
	shimDir := t.TempDir()
	witness := filepath.Join(t.TempDir(), "witness")
	for _, host := range []string{"claude", "codex", "pi", "safehouse"} {
		script := "#!/bin/sh\necho " + host + " >> " + witness + "\nexit 0\n"
		if err := os.WriteFile(filepath.Join(shimDir, host), []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", shimDir)

	var stdout, stderr bytes.Buffer
	if code := Run([]string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("--version exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}

	if data, err := os.ReadFile(witness); err == nil {
		t.Fatalf("--version executed a host CLI; witness file contains:\n%s", data)
	}

	for _, line := range strings.Split(stdout.String(), "\n") {
		for _, host := range []string{"claude:", "codex:", "pi:"} {
			if strings.HasPrefix(line, host) {
				t.Fatalf("--version emitted a per-host line %q; the per-host block is gone", line)
			}
		}
	}
}

// TestVersionAmbiguousMarkersExitZero (AC-3) pins the one branch that can break
// every boot. The FO version gate runs `--version` before discovery, so refusing
// on ambiguous markers — which the DISPATCH detector correctly does — would abort
// every session where a nested runtime leaks a second marker. Driven through Run,
// so the exit code is the real one.
func TestVersionAmbiguousMarkersExitZero(t *testing.T) {
	t.Setenv("CODEX_THREAD_ID", "01937f2a")
	t.Setenv("CLAUDECODE", "1")

	var stdout, stderr bytes.Buffer
	code := Run([]string{"--version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("--version with ambiguous markers exited %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := "Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE) — pass --host"
	if !lineEquals(stdout.String(), want) {
		t.Fatalf("--version output = %q, want a whole line %q", stdout.String(), want)
	}
}
