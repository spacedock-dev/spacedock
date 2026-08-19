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

// TestSessionStartCompactReminderHookGate exercises the shipped SessionStart
// hook script directly — the only properties a test can prove mechanically:
// it emits the exact reminder JSON only when SPACEDOCK_BIN is set (a
// launcher-marked session — the plugin being merely installed must not hand
// out First Officer instructions after every compaction) AND the hook JSON
// on stdin names source=compact, and it stays silent (exit 0, no stdout) on
// every other combination, including malformed/empty input. It does NOT and
// cannot prove that Claude Code or Codex injects this stdout into the
// resumed model's context, or that the FO obeys the reminder once injected —
// see docs/dev/.spacedock-state/force-boot-at-compaction-boundary.md's
// "How proven" section for why those are out of reach of a unit test.
//
// Each subprocess runs with an explicit, minimal environment rather than the
// test process's inherited one: the ambient shell in this repo commonly
// exports a real SPACEDOCK_BIN (this entity's own incident is why), which
// would silently make the "gate closed" cases pass for the wrong reason.
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

	const wantReminder = `{"hookSpecificOutput":{"hookEventName":"SessionStart","additionalContext":"COMPACTION BOUNDARY — your bindings are stale, not just your narrative.\n\nBefore any gate, merge, or state mutation, re-run the Startup procedure: resolve the binary, then ` +
		"`${SPACEDOCK_BIN:-spacedock} status --boot --identify --json`" +
		`, and READ the whole boot record — not only the keys you came looking for. It carries the registered mods, ready-gate readiness, the binary version gate, PR state, and live team state. Any of those may contradict what your summary says.\n\nSpecifically distrust: which gates you believe are presented (check readiness, not memory), which workers you believe are alive, the SPACEDOCK_BIN path, and which contract version is installed."}}` + "\n"

	cases := []struct {
		name           string
		setSpacedock   bool   // false: do not put SPACEDOCK_BIN in the env at all
		spacedockValue string // used only when setSpacedock is true
		stdin          string
		wantReminded   bool
	}{
		{
			name:           "launcher-marked + source=compact emits the reminder",
			setSpacedock:   true,
			spacedockValue: "/usr/bin/true",
			stdin:          `{"session_id":"s1","hook_event_name":"SessionStart","source":"compact"}`,
			wantReminded:   true,
		},
		{
			name:  "no SPACEDOCK_BIN at all -> silent even on source=compact",
			stdin: `{"session_id":"s1","hook_event_name":"SessionStart","source":"compact"}`,
		},
		{
			name:           "SPACEDOCK_BIN set but empty -> silent (sh -n semantics)",
			setSpacedock:   true,
			spacedockValue: "",
			stdin:          `{"session_id":"s1","hook_event_name":"SessionStart","source":"compact"}`,
		},
		{
			name:           "launcher-marked + source=startup -> silent",
			setSpacedock:   true,
			spacedockValue: "/usr/bin/true",
			stdin:          `{"session_id":"s1","hook_event_name":"SessionStart","source":"startup"}`,
		},
		{
			name:           "launcher-marked + source=resume -> silent",
			setSpacedock:   true,
			spacedockValue: "/usr/bin/true",
			stdin:          `{"session_id":"s1","hook_event_name":"SessionStart","source":"resume"}`,
		},
		{
			name:           "launcher-marked + source=clear -> silent",
			setSpacedock:   true,
			spacedockValue: "/usr/bin/true",
			stdin:          `{"session_id":"s1","hook_event_name":"SessionStart","source":"clear"}`,
		},
		{
			name:           "launcher-marked + empty stdin -> silent",
			setSpacedock:   true,
			spacedockValue: "/usr/bin/true",
		},
		{
			name:           "launcher-marked + malformed non-JSON stdin -> silent",
			setSpacedock:   true,
			spacedockValue: "/usr/bin/true",
			stdin:          "not json at all",
		},
		{
			name:           "launcher-marked + source present but not a string -> silent",
			setSpacedock:   true,
			spacedockValue: "/usr/bin/true",
			stdin:          `{"session_id":"s1","hook_event_name":"SessionStart","source":123}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("/bin/sh", script)
			cmd.Stdin = strings.NewReader(tc.stdin)
			env := []string{"PATH=" + os.Getenv("PATH")}
			if tc.setSpacedock {
				env = append(env, "SPACEDOCK_BIN="+tc.spacedockValue)
			}
			cmd.Env = env
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("hook exited non-zero: %v (stderr: %s)", err, stderr.String())
			}
			want := ""
			if tc.wantReminded {
				want = wantReminder
			}
			if got := stdout.String(); got != want {
				t.Errorf("stdout = %q, want %q", got, want)
			}
		})
	}
}

// TestSessionStartCompactReminderHookIsWired confirms the hook is installed
// by the declared mechanism: ONE script registered ONCE in the shared
// plugin hooks.json under SessionStart/^compact$, and BOTH plugin manifests
// (.claude-plugin and .codex-plugin) activate that same hooks.json — so
// Claude and Codex sessions each get exactly one reminder block, not two.
func TestSessionStartCompactReminderHookIsWired(t *testing.T) {
	repo := repoRoot(t)

	if _, err := os.Stat(filepath.Join(repo, "hooks", "codex_session_start_compact.sh")); !os.IsNotExist(err) {
		t.Errorf("hooks/codex_session_start_compact.sh still exists; it should have been consolidated into session_start_compact_reminder.sh: %v", err)
	}

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

	const wantCommand = "${CLAUDE_PLUGIN_ROOT:-${PLUGIN_ROOT}}/hooks/session_start_compact_reminder.sh"
	var matched []struct {
		Type    string
		Command string
	}
	for _, entry := range doc.Hooks.SessionStart {
		if entry.Matcher != "^compact$" {
			continue
		}
		for _, h := range entry.Hooks {
			matched = append(matched, struct{ Type, Command string }{h.Type, h.Command})
		}
	}
	if len(matched) != 1 {
		t.Fatalf("hooks.json SessionStart/^compact$ has %d command(s), want exactly 1 (got %+v)", len(matched), matched)
	}
	if matched[0].Command != wantCommand {
		t.Errorf("hooks.json SessionStart/^compact$ command = %q, want %q", matched[0].Command, wantCommand)
	}

	for _, manifest := range []string{
		filepath.Join(".claude-plugin", "plugin.json"),
		filepath.Join(".codex-plugin", "plugin.json"),
	} {
		pluginData, err := os.ReadFile(filepath.Join(repo, manifest))
		if err != nil {
			t.Fatal(err)
		}
		var plugin map[string]json.RawMessage
		if err := json.Unmarshal(pluginData, &plugin); err != nil {
			t.Fatal(err)
		}
		hooksRef, ok := plugin["hooks"]
		if !ok {
			t.Errorf("%s has no \"hooks\" field; the SessionStart hook above is dead for this host", manifest)
			continue
		}
		var hooksPath string
		if err := json.Unmarshal(hooksRef, &hooksPath); err != nil {
			t.Fatalf("parse %s hooks field: %v", manifest, err)
		}
		if hooksPath != "./hooks.json" {
			t.Errorf("%s hooks = %q, want %q", manifest, hooksPath, "./hooks.json")
		}
	}
}

// TestSessionStartCompactReminderPluginRootFallbackResolves proves the
// hooks.json command's variable expansion, not just its literal string, by
// actually running it through /bin/sh — the same interpreter Claude Code and
// Codex invoke it with. This is the regression test for a real bug this
// entity's own live spike caught: Claude Code sets CLAUDE_PLUGIN_ROOT for a
// hook subprocess, never a bare PLUGIN_ROOT (that is Codex's token,
// e143969b8). Before the fix, hooks.json referenced only ${PLUGIN_ROOT},
// which is unset under Claude Code and expands to empty under /bin/sh —
// producing the broken path /hooks/session_start_compact_reminder.sh
// (verified live: "SessionStart:compact hook error ... No such file or
// directory", recorded verbatim in ideation-spike-evidence.md section 6). A
// JSON string-equality check would not have caught this: the literal
// ${PLUGIN_ROOT} template is syntactically well-formed JSON and a
// syntactically valid shell command; only executing it under each host's
// actual environment shape reveals the resolved path is wrong.
func TestSessionStartCompactReminderPluginRootFallbackResolves(t *testing.T) {
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
					Command string `json:"command"`
				} `json:"hooks"`
			} `json:"SessionStart"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse hooks.json: %v", err)
	}
	var command string
	for _, entry := range doc.Hooks.SessionStart {
		if entry.Matcher != "^compact$" {
			continue
		}
		for _, h := range entry.Hooks {
			command = h.Command
		}
	}
	if command == "" {
		t.Fatal("could not find the SessionStart/^compact$ command in hooks.json")
	}

	cases := []struct {
		name string
		env  []string
		want string
	}{
		{
			name: "Claude Code sets CLAUDE_PLUGIN_ROOT",
			env:  []string{"CLAUDE_PLUGIN_ROOT=/claude/plugin/root"},
			want: "/claude/plugin/root/hooks/session_start_compact_reminder.sh",
		},
		{
			name: "Codex sets PLUGIN_ROOT (no CLAUDE_PLUGIN_ROOT)",
			env:  []string{"PLUGIN_ROOT=/codex/plugin/root"},
			want: "/codex/plugin/root/hooks/session_start_compact_reminder.sh",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("/bin/sh", "-c", "echo "+command)
			cmd.Env = append([]string{"PATH=" + os.Getenv("PATH")}, tc.env...)
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
