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

// TestCompactBootReminderHookFiresOnlyOnCompact exercises the shipped
// SessionStart hook script directly — the only property a test can prove
// mechanically: it emits the exact reminder bytes when the hook JSON on
// stdin names source=compact, and stays silent (exit 0, no stdout) on every
// other source and on malformed/empty input. It does NOT and cannot prove
// that Claude Code injects this stdout into the resumed model's context, or
// that the FO obeys the reminder once injected — see
// docs/dev/.spacedock-state/force-boot-at-compaction-boundary.md's "How
// proven" section for why those are out of reach of a unit test.
func TestCompactBootReminderHookFiresOnlyOnCompact(t *testing.T) {
	repo := repoRoot(t)
	script := filepath.Join(repo, "hooks", "compact_boot_reminder.sh")

	info, err := os.Stat(script)
	if err != nil {
		t.Fatalf("stat hook script: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Errorf("%s is not executable (mode %v); the plugin invokes it directly", script, info.Mode())
	}

	// @BT@ stands in for a literal backtick, which cannot appear inside the
	// raw string it delimits.
	wantReminder := strings.ReplaceAll(`COMPACTION BOUNDARY — your bindings are stale, not just your narrative.

Before any gate, merge, or state mutation, re-run the Startup procedure:
resolve the binary, then @BT@${SPACEDOCK_BIN:-spacedock} status --boot --identify
--json@BT@, and READ the whole boot record — not only the keys you came looking
for. It carries the registered mods, ready-gate readiness, the binary version
gate, PR state, and live team state. Any of those may contradict what your
summary says.

Specifically distrust: which gates you believe are presented (check
readiness, not memory), which workers you believe are alive, the
SPACEDOCK_BIN path, and which contract version is installed.
`, "@BT@", "`")

	cases := []struct {
		name       string
		stdin      string
		wantStdout string
	}{
		{
			name:       "source=compact emits the exact reminder",
			stdin:      `{"session_id":"s1","hook_event_name":"SessionStart","source":"compact"}`,
			wantStdout: wantReminder,
		},
		{
			name:  "source=startup is silent",
			stdin: `{"session_id":"s1","hook_event_name":"SessionStart","source":"startup"}`,
		},
		{
			name:  "source=resume is silent",
			stdin: `{"session_id":"s1","hook_event_name":"SessionStart","source":"resume"}`,
		},
		{
			name:  "source=clear is silent",
			stdin: `{"session_id":"s1","hook_event_name":"SessionStart","source":"clear"}`,
		},
		{
			name:  "empty stdin is silent",
			stdin: "",
		},
		{
			name:  "malformed non-JSON stdin is silent",
			stdin: "not json at all",
		},
		{
			name:  "source field present but not a string is silent",
			stdin: `{"session_id":"s1","hook_event_name":"SessionStart","source":123}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("/bin/sh", script)
			cmd.Stdin = strings.NewReader(tc.stdin)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("hook exited non-zero: %v (stderr: %s)", err, stderr.String())
			}
			if got := stdout.String(); got != tc.wantStdout {
				t.Errorf("stdout = %q, want %q", got, tc.wantStdout)
			}
		})
	}
}

// TestCompactBootReminderHookIsWired confirms the hook is installed by the
// declared mechanism: registered in the shared plugin hooks.json under
// SessionStart/^compact$, and that .claude-plugin/plugin.json actually
// activates hooks.json (it did not before this entity — see
// TestDispatchAckMachineryIsAbsent's prior wantHooks:false expectation). It
// does not touch the pre-existing codex_session_start_compact.sh entry,
// which stays registered unchanged.
func TestCompactBootReminderHookIsWired(t *testing.T) {
	repo := repoRoot(t)

	data, err := os.ReadFile(filepath.Join(repo, "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Hooks struct {
			SessionStart []struct {
				Matcher string `json:"matcher"`
				Hooks   []struct {
					Type    string `json:"type"`
					Command string `json:"command"`
				} `json:"hooks"`
			} `json:"SessionStart"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse hooks.json: %v", err)
	}

	const wantCommand = "${PLUGIN_ROOT}/hooks/compact_boot_reminder.sh"
	found := false
	for _, entry := range doc.Hooks.SessionStart {
		if entry.Matcher != "^compact$" {
			continue
		}
		for _, h := range entry.Hooks {
			if h.Type == "command" && h.Command == wantCommand {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("hooks.json has no SessionStart matcher=^compact$ entry running %q", wantCommand)
	}

	pluginData, err := os.ReadFile(filepath.Join(repo, ".claude-plugin", "plugin.json"))
	if err != nil {
		t.Fatal(err)
	}
	var plugin map[string]json.RawMessage
	if err := json.Unmarshal(pluginData, &plugin); err != nil {
		t.Fatal(err)
	}
	hooksRef, ok := plugin["hooks"]
	if !ok {
		t.Fatal(".claude-plugin/plugin.json has no \"hooks\" field; the SessionStart hook above is dead")
	}
	var hooksPath string
	if err := json.Unmarshal(hooksRef, &hooksPath); err != nil {
		t.Fatalf("parse .claude-plugin/plugin.json hooks field: %v", err)
	}
	if hooksPath != "./hooks.json" {
		t.Errorf(".claude-plugin/plugin.json hooks = %q, want %q", hooksPath, "./hooks.json")
	}
}
