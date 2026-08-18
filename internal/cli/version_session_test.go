// ABOUTME: --version reports the session it is actually in — the Runtime line,
// ABOUTME: the session sandbox state, and no host CLI exec.
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
// one, report the version. Both shapes start with the same THREE lines — the
// version line, the `OS: <goos>/<goarch>` line, and the `Channel:` line — so the
// outside case is three lines (no Runtime line, no Sandbox line, and no contract
// token), while every in-session case adds exactly three lines below them.
// devBranch is untouched by this file (package default "next"), so every
// Channel line here renders the edge rendering: `edge (spacedock@spacedock-edge)`.
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
			// BECAUSE the wrap already happened. The sandbox line names the sandbox
			// and nothing else. CLAUDE_CODE_SESSION_ID is SET here and appears
			// nowhere in the render — the Runtime line is byte-identical to the
			// `claude-no-session-id-var` case below, which is what the session
			// segment's removal means.
			name:     "claude-inside-safehouse",
			vars:     claudeSession,
			lookPath: lookMissing,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Channel: edge (spacedock@spacedock-edge)\n" +
				"Runtime: claude (CLAUDECODE)\n" +
				"Sandbox: agent-safehouse\n" +
				"contract 3\n",
		},
		{
			name:     "outside-every-runtime-is-three-lines",
			vars:     map[string]string{},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Channel: edge (spacedock@spacedock-edge)\n",
		},
		{
			// The same host with no session variable set at all. Its Runtime line
			// matches the case above character for character; only the sandbox
			// differs. A stray `, session ` or `()` fails this exact-match.
			name:     "claude-no-session-id-var",
			vars:     map[string]string{"CLAUDECODE": "1"},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Channel: edge (spacedock@spacedock-edge)\n" +
				"Runtime: claude (CLAUDECODE)\n" +
				"Sandbox: none (safehouse available)\n" +
				"contract 3\n",
		},
		{
			// Two markers of the SAME host are not ambiguity: pi resolves, and both
			// markers are reported in table order.
			name: "pi-two-markers-same-host",
			vars: map[string]string{
				"PI_CODING_AGENT":     "1",
				"PI_CODING_AGENT_DIR": "/home/u/.pi/agent",
			},
			lookPath: lookMissing,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Channel: edge (spacedock@spacedock-edge)\n" +
				"Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)\n" +
				"Sandbox: none (safehouse not installed)\n" +
				"contract 3\n",
		},
		{
			// codex's marker carries a UUID-shaped value. The Runtime line reports
			// that the marker is SET, never any part of what it holds, so no
			// fragment of `01937f2a-bbbb-cccc` may appear.
			name:     "codex-marker-value-never-rendered",
			vars:     map[string]string{"CODEX_THREAD_ID": "01937f2a-bbbb-cccc"},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Channel: edge (spacedock@spacedock-edge)\n" +
				"Runtime: codex (CODEX_THREAD_ID)\n" +
				"Sandbox: none (safehouse available)\n" +
				"contract 3\n",
		},
		{
			// Ambiguity is REPORTED, not guessed at, and carries no remedy: this
			// surface exits 0 and has no fault to remedy. The one command the
			// ambiguity blocks, `dispatch build`, prints its own complete remedy
			// naming all three valid hosts.
			name: "ambiguous-markers",
			vars: map[string]string{
				"CODEX_THREAD_ID":        "01937f2a",
				"CLAUDECODE":             "1",
				"CLAUDE_CODE_SESSION_ID": "afd74765-9000",
			},
			lookPath: lookFound,
			want: "spacedock " + version + "\n" +
				"OS: " + osToken + "\n" +
				"Channel: edge (spacedock@spacedock-edge)\n" +
				"Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE)\n" +
				"Sandbox: none (safehouse available)\n" +
				"contract 3\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderSession(tc.vars, tc.lookPath)
			// AC-2: the trimmed decorations are unreachable in EVERY shape. These
			// run before the exact match on purpose: re-adding a decoration and
			// regenerating the `want` literals in step would slip past an
			// exact-match-only test, and these name what came back instead.
			for _, banned := range []struct{ token, what string }{
				{", session ", "a session identifier segment"},
				{"pass --host", "a --host remedy suffix"},
				{"inside", "the Sandbox line restating its own label"},
			} {
				if strings.Contains(got, banned.token) {
					t.Fatalf("--version rendered %q — %s, which has no reader:\n%s", banned.token, banned.what, got)
				}
			}
			if got != tc.want {
				t.Fatalf("--version render mismatch:\n got=%q\nwant=%q", got, tc.want)
			}
			// AC-2a: no empty or dangling parenthetical is ever reachable.
			if strings.Contains(got, "()") {
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

// TestVersionSandboxLineNamesTheSandbox (AC-2) is the positive half of the
// Sandbox rename: the absence loop above says what must not appear, and this
// says what must. Sandboxed, the value is the sandbox's NAME and nothing else —
// under a `Sandbox: ` label the name alone already carries both facts. Both
// unsandboxed arms lead with `none`, so a reader classifies on name-or-`none`
// rather than on a relationship word.
//
// It holds independently of the banned literals: restoring SessionState's old
// `inside (...)` / `not sandboxed (...)` shape turns this red by name even with
// every `want` string in TestVersionSessionRender regenerated to match.
func TestVersionSandboxLineNamesTheSandbox(t *testing.T) {
	sandboxLineOf := func(out string) string {
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(line, "Sandbox: ") {
				return line
			}
		}
		t.Fatalf("no Sandbox: line in %q", out)
		return ""
	}

	cases := []struct {
		name     string
		vars     map[string]string
		lookPath func(string) (string, error)
		want     string
	}{
		// Sandboxed with safehouse off PATH — the live configuration, and the one
		// where the name is the only honest answer.
		{"sandboxed-renders-the-bare-name", claudeSession, lookMissing, "Sandbox: agent-safehouse"},
		{"unsandboxed-safehouse-available", map[string]string{"CLAUDECODE": "1"}, lookFound, "Sandbox: none (safehouse available)"},
		{"unsandboxed-safehouse-not-installed", map[string]string{"CLAUDECODE": "1"}, lookMissing, "Sandbox: none (safehouse not installed)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sandboxLineOf(renderSession(tc.vars, tc.lookPath))
			if got != tc.want {
				t.Fatalf("Sandbox line = %q, want %q", got, tc.want)
			}
		})
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
	want := "Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE)"
	if !lineEquals(stdout.String(), want) {
		t.Fatalf("--version output = %q, want a whole line %q", stdout.String(), want)
	}
}
