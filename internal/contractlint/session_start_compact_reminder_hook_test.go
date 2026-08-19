package contractlint

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// sessionStartCompactCommand reads hooks.json and returns the sole command
// registered under SessionStart/^compact$, failing the test if there isn't
// exactly one. Shared by the two parity tests below.
func sessionStartCompactCommand(t *testing.T, repo string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repo, "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Hooks struct {
			SessionStart []struct {
				Matcher string `json:"matcher"`
				Hooks   []struct {
					Command string `json:"command"`
				} `json:"hooks"`
			} `json:"SessionStart"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse hooks.json: %v", err)
	}
	var matched []string
	for _, entry := range doc.Hooks.SessionStart {
		if entry.Matcher != "^compact$" {
			continue
		}
		for _, h := range entry.Hooks {
			matched = append(matched, h.Command)
		}
	}
	if len(matched) != 1 {
		t.Fatalf("hooks.json SessionStart/^compact$ has %d command(s), want exactly 1 (got %v)", len(matched), matched)
	}
	return matched[0]
}

// TestSessionStartCompactReminderHookIsWired (static) and
// TestSessionStartCompactReminderPluginRootFallbackResolves (dynamic)
// together prove this entity's invariant, HOST SIGNATURE PARITY: one
// script, one SessionStart/^compact$ entry, both manifests activating the
// same hooks.json, resolving identically under either host's plugin-root var.
func TestSessionStartCompactReminderHookIsWired(t *testing.T) {
	repo := repoRoot(t)

	if _, err := os.Stat(filepath.Join(repo, "hooks", "codex_session_start_compact.sh")); !os.IsNotExist(err) {
		t.Errorf("hooks/codex_session_start_compact.sh still exists; it should have been consolidated into session_start_compact_reminder.sh: %v", err)
	}

	const wantCommand = "${CLAUDE_PLUGIN_ROOT:-${PLUGIN_ROOT}}/hooks/session_start_compact_reminder.sh"
	if got := sessionStartCompactCommand(t, repo); got != wantCommand {
		t.Errorf("hooks.json SessionStart/^compact$ command = %q, want %q", got, wantCommand)
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
// open/closed behavior, distinct from the parity invariant above: an unset
// SPACEDOCK_BIN and a non-compact source both stay silent (external
// properties), and the open row's payload is checked by property — valid
// JSON, SessionStart event name, names the boot command — not by
// byte-comparing the wording. Stdin malformation and field-typing
// permutations are dropped: they test sed and test(1), not this invariant.
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
		{"launcher-marked + source=startup -> silent", true, `{"session_id":"s1","hook_event_name":"SessionStart","source":"startup"}`, false},
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
