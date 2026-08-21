package contractlint

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type sessionStartEntry struct {
	Matcher string `json:"matcher"`
	Hooks   []struct {
		Command string `json:"command"`
	} `json:"hooks"`
}

// sessionStartEntries returns EVERY SessionStart entry in hooks.json. Nothing
// filters on the matcher: keying the lookup off the matcher we expect makes it
// a lookup key rather than a value under test, and hides an added sibling entry.
func sessionStartEntries(t *testing.T, repo string) []sessionStartEntry {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repo, "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Hooks struct {
			SessionStart []sessionStartEntry `json:"SessionStart"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse hooks.json: %v", err)
	}
	return doc.Hooks.SessionStart
}

// sessionStartCompactCommand returns the sole command registered under
// SessionStart across all entries, failing the test if there isn't exactly
// one. Shared by the two parity tests below.
func sessionStartCompactCommand(t *testing.T, repo string) string {
	t.Helper()
	var matched []string
	for _, entry := range sessionStartEntries(t, repo) {
		for _, h := range entry.Hooks {
			matched = append(matched, h.Command)
		}
	}
	if len(matched) != 1 {
		t.Fatalf("hooks.json SessionStart has %d command(s), want exactly 1 (got %v)", len(matched), matched)
	}
	return matched[0]
}

// sessionStartSources reports which of the host's SessionStart source values
// the declared matchers accept — the matcher exercised as a regexp, not
// string-compared. Every matcher compiles or the test fails.
func sessionStartSources(t *testing.T, repo string) string {
	t.Helper()
	var matchers []*regexp.Regexp
	for _, entry := range sessionStartEntries(t, repo) {
		re, err := regexp.Compile(entry.Matcher)
		if err != nil {
			t.Fatalf("SessionStart matcher %q does not compile: %v", entry.Matcher, err)
		}
		matchers = append(matchers, re)
	}
	var accepted []string
	for _, src := range []string{"startup", "resume", "clear", "compact", "fork"} {
		for _, re := range matchers {
			if re.MatchString(src) {
				accepted = append(accepted, src)
				break
			}
		}
	}
	return strings.Join(accepted, ",")
}

// TestSessionStartCompactReminderHookIsWired (static) and
// TestSessionStartCompactReminderPluginRootFallbackResolves (dynamic)
// together prove this entity's invariant, HOST SIGNATURE PARITY: one script,
// one SessionStart entry accepting exactly the two context-loss sources, both
// manifests activating the same hooks.json, resolving identically under
// either host's plugin-root var.
func TestSessionStartCompactReminderHookIsWired(t *testing.T) {
	repo := repoRoot(t)

	if _, err := os.Stat(filepath.Join(repo, "hooks", "codex_session_start_compact.sh")); !os.IsNotExist(err) {
		t.Errorf("hooks/codex_session_start_compact.sh still exists; it should have been consolidated into session_start_compact_reminder.sh: %v", err)
	}

	const wantCommand = "${CLAUDE_PLUGIN_ROOT:-${PLUGIN_ROOT}}/hooks/session_start_compact_reminder.sh"
	if got := sessionStartCompactCommand(t, repo); got != wantCommand {
		t.Errorf("hooks.json SessionStart command = %q, want %q", got, wantCommand)
	}

	// The matcher is behavior, not a label: it must accept exactly the two
	// context-loss sources. Covers a widened matcher (`.*` takes all five), a
	// wrong matcher, and an added permissive sibling entry, in one assertion.
	if got, want := sessionStartSources(t, repo), "clear,compact"; got != want {
		t.Errorf("SessionStart matchers accept [%s], want exactly [%s]", got, want)
	}

	for _, manifest := range []string{
		filepath.Join(".claude-plugin", "plugin.json"),
		filepath.Join(".codex-plugin", "plugin.json"),
	} {
		data, err := os.ReadFile(filepath.Join(repo, manifest))
		if err != nil {
			t.Fatal(err)
		}
		var plugin map[string]json.RawMessage
		if err := json.Unmarshal(data, &plugin); err != nil {
			t.Fatal(err)
		}
		var hooksPath string
		if err := json.Unmarshal(plugin["hooks"], &hooksPath); err != nil {
			t.Fatalf("%s has no usable \"hooks\" field: %v", manifest, err)
		}
		if hooksPath != "./hooks.json" {
			t.Errorf("%s hooks = %q, want %q", manifest, hooksPath, "./hooks.json")
		}
	}
}

// TestSessionStartCompactReminderPluginRootFallbackResolves runs the
// hooks.json command through /bin/sh instead of string-comparing it — what
// caught the shipped bug: Claude Code sets CLAUDE_PLUGIN_ROOT, never a bare
// PLUGIN_ROOT (Codex's token, e143969b8); only execution reveals the broken path.
func TestSessionStartCompactReminderPluginRootFallbackResolves(t *testing.T) {
	repo := repoRoot(t)
	command := sessionStartCompactCommand(t, repo)

	cases := []struct{ name, env, want string }{
		{"Claude Code sets CLAUDE_PLUGIN_ROOT", "CLAUDE_PLUGIN_ROOT=/claude/plugin/root", "/claude/plugin/root/hooks/session_start_compact_reminder.sh"},
		{"Codex sets PLUGIN_ROOT (no CLAUDE_PLUGIN_ROOT)", "PLUGIN_ROOT=/codex/plugin/root", "/codex/plugin/root/hooks/session_start_compact_reminder.sh"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("/bin/sh", "-c", "echo "+command)
			cmd.Env = []string{"PATH=" + os.Getenv("PATH"), tc.env}
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("sh -c echo failed: %v", err)
			}
			if got := strings.TrimSpace(string(out)); got != tc.want {
				t.Errorf("resolved path = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestSessionStartCompactReminderHookGate exercises the script's own
// open/closed behavior, distinct from the parity invariant above: the script
// gates on SPACEDOCK_BIN alone, and an open row's payload is checked by
// property — valid JSON, SessionStart event name, names the boot command —
// not by byte-comparing the wording. Source filtering belongs to the host
// ^compact$ matcher pinned above; when the script re-implemented it with a
// greedy sed, a nested "source" object silently suppressed the reminder.
func TestSessionStartCompactReminderHookGate(t *testing.T) {
	repo := repoRoot(t)
	script := filepath.Join(repo, "hooks", "session_start_compact_reminder.sh")

	info, err := os.Stat(script)
	if err != nil {
		t.Fatalf("stat hook script: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Errorf("%s is not executable (mode %v); the plugin invokes it directly", script, info.Mode())
	}

	cases := []struct {
		name         string
		setSpacedock bool
		stdin        string
		wantReminded bool
	}{
		{"launcher-marked + source=compact emits the reminder", true, `{"session_id":"s1","hook_event_name":"SessionStart","source":"compact"}`, true},
		{"no SPACEDOCK_BIN -> silent even on source=compact", false, `{"session_id":"s1","hook_event_name":"SessionStart","source":"compact"}`, false},
		{"launcher-marked emits whatever the payload says: a nested \"source\" no longer suppresses it", true, `{"session_id":"s1","hook_event_name":"SessionStart","source":"compact","extra":{"source":"startup"}}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Explicit env: this repo's ambient shell often exports a real
			// SPACEDOCK_BIN (this entity's own incident is why).
			cmd := exec.Command("/bin/sh", script)
			cmd.Stdin = strings.NewReader(tc.stdin)
			env := []string{"PATH=" + os.Getenv("PATH")}
			if tc.setSpacedock {
				env = append(env, "SPACEDOCK_BIN=/usr/bin/true")
			}
			cmd.Env = env
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("hook exited non-zero: %v (stderr: %s)", err, stderr.String())
			}
			if !tc.wantReminded {
				if got := stdout.String(); got != "" {
					t.Errorf("stdout = %q, want empty (gate closed)", got)
				}
				return
			}
			var payload struct {
				HookSpecificOutput struct {
					HookEventName     string `json:"hookEventName"`
					AdditionalContext string `json:"additionalContext"`
				} `json:"hookSpecificOutput"`
			}
			if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
				t.Fatalf("stdout did not parse as JSON: %v (stdout: %s)", err, stdout.String())
			}
			if payload.HookSpecificOutput.HookEventName != "SessionStart" {
				t.Errorf("hookEventName = %q, want %q", payload.HookSpecificOutput.HookEventName, "SessionStart")
			}
			if !strings.Contains(payload.HookSpecificOutput.AdditionalContext, "status --boot") {
				t.Errorf("additionalContext does not name the boot command (status --boot): %q", payload.HookSpecificOutput.AdditionalContext)
			}
		})
	}
}
