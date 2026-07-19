// ABOUTME: Offline fixture for the Codex «post-compact-notice» binding — parses the
// ABOUTME: shipped .codex-plugin hook config, drives manual|auto, and proves harmless absence.
package ensigncycle

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// postCompactRepoRoot resolves the plugin root (repo root) from the test package
// directory internal/ensigncycle. The live-tagged repoRoot helper is not compiled
// into the default build, so this locates the shipped .codex-plugin files directly.
func postCompactRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

// codexHookHandler is one command handler inside a hooks.json matcher group.
type codexHookHandler struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// codexHookGroup is one matcher group under an event key in hooks.json.
type codexHookGroup struct {
	Matcher string             `json:"matcher"`
	Hooks   []codexHookHandler `json:"hooks"`
}

type codexHooksConfig struct {
	Hooks map[string][]codexHookGroup `json:"hooks"`
}

// requiredNoticePhrases are the load-bearing substrings the reread-and-reconcile
// systemMessage must carry; asserting the phrases (not a whole-string equality)
// keeps the fixture from breaking on incidental punctuation edits while still
// proving the instruction — reread the authoritative FO contract and reconcile —
// is present.
var requiredNoticePhrases = []string{
	"reread",
	"`spacedock:first-officer`",
	"reconcile durable",
}

// loadShippedPostCompactHook reads the shipped manifest + hooks.json and returns
// the plugin-root path and the single PostCompact matcher group.
func loadShippedPostCompactHook(t *testing.T) (root string, group codexHookGroup) {
	t.Helper()
	root = postCompactRepoRoot(t)

	manifestBytes, err := os.ReadFile(filepath.Join(root, ".codex-plugin", "plugin.json"))
	if err != nil {
		t.Fatalf("read .codex-plugin/plugin.json: %v", err)
	}
	var manifest struct {
		Hooks string `json:"hooks"`
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("parse plugin.json: %v", err)
	}
	if manifest.Hooks != "./hooks.json" {
		t.Fatalf("plugin.json hooks = %q, want %q (the declared hooks key Codex loads)", manifest.Hooks, "./hooks.json")
	}

	hooksBytes, err := os.ReadFile(filepath.Join(root, "hooks.json"))
	if err != nil {
		t.Fatalf("read hooks.json: %v", err)
	}
	var cfg codexHooksConfig
	if err := json.Unmarshal(hooksBytes, &cfg); err != nil {
		t.Fatalf("parse hooks.json: %v", err)
	}
	groups := cfg.Hooks["PostCompact"]
	if len(groups) != 1 {
		t.Fatalf("hooks.json PostCompact groups = %d, want exactly 1", len(groups))
	}
	if len(groups[0].Hooks) != 1 || groups[0].Hooks[0].Type != "command" {
		t.Fatalf("PostCompact group must carry exactly one command handler, got %+v", groups[0].Hooks)
	}
	return root, groups[0]
}

// pluginRootToken is the plugin-root variable Codex substitutes into a hook command
// string before it exec's the command (there is no shell). Verified against a live
// Codex 0.144.x CLI (artifacts/codex-0.144.4-plugin-hooks-spike.md): the brace form
// ${PLUGIN_ROOT} is replaced with the materialized plugin directory, so a bundled
// script referenced as ${PLUGIN_ROOT}/hooks/x.sh resolves to an absolute path
// independent of the session cwd. The bare relative form ./hooks/x.sh is NOT
// plugin-root-relative: Codex resolves it against the session cwd (the operator's
// project), so it fails whenever the FO runs in any repo other than the plugin's own.
const pluginRootToken = "${PLUGIN_ROOT}"

// resolveHookCommand mirrors Codex's command resolution: the shipped command MUST be
// plugin-root-absolute via ${PLUGIN_ROOT}; the token is then substituted with the real
// plugin root. It FAILS on a cwd-relative command — the exact defect where the notice
// would not fire from any session cwd other than the plugin repo. Because it substitutes
// ${PLUGIN_ROOT} rather than joining a bare relative path onto the plugin root, a
// cwd-relative command can no longer masquerade as resolvable in the offline gate.
func resolveHookCommand(t *testing.T, root, command string) string {
	t.Helper()
	if !strings.HasPrefix(command, pluginRootToken+"/") {
		t.Fatalf("hook command %q must be plugin-root-absolute (%s/...); a cwd-relative command resolves against the session cwd, not the plugin root, so the notice would not fire whenever the FO operates outside the plugin repo", command, pluginRootToken)
	}
	rel := strings.TrimPrefix(command, pluginRootToken+"/")
	return filepath.Join(root, filepath.FromSlash(rel))
}

// TestCodexPostCompactHookMatchesManualAndAuto proves AC-3's configuration half:
// the shipped PostCompact matcher fires for both compaction sources.
func TestCodexPostCompactHookMatchesManualAndAuto(t *testing.T) {
	_, group := loadShippedPostCompactHook(t)

	re, err := regexp.Compile(group.Matcher)
	if err != nil {
		t.Fatalf("PostCompact matcher %q is not a valid regexp: %v", group.Matcher, err)
	}
	for _, source := range []string{"manual", "auto"} {
		if !re.MatchString(source) {
			t.Errorf("PostCompact matcher %q does not match compaction source %q", group.Matcher, source)
		}
	}
}

// TestCodexPostCompactHookEmitsOneSystemMessagePerEvent is AC-3's command-level
// fixture: it drives the shipped hook command with both the manual and auto event
// payloads and asserts one valid JSON systemMessage per event carrying the required
// reread-and-reconcile instruction. It does NOT claim the warning enters model
// context — the 0.144.4 probe already bounded that to a captain-facing warning.
func TestCodexPostCompactHookEmitsOneSystemMessagePerEvent(t *testing.T) {
	root, group := loadShippedPostCompactHook(t)

	command := resolveHookCommand(t, root, group.Hooks[0].Command)
	info, err := os.Stat(command)
	if err != nil {
		t.Fatalf("shipped hook command %q not found: %v", command, err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatalf("shipped hook command %q is not executable (mode %v)", command, info.Mode())
	}

	var firstMessage string
	for i, source := range []string{"manual", "auto"} {
		payload := `{"hook_event_name":"PostCompact","trigger":"` + source + `"}`
		out := runShippedHook(t, root, command, payload)

		var decoded map[string]json.RawMessage
		if err := json.Unmarshal([]byte(out), &decoded); err != nil {
			t.Fatalf("%s: hook stdout is not valid JSON: %v\nstdout: %q", source, err, out)
		}
		if len(decoded) != 1 {
			t.Fatalf("%s: hook output has %d keys, want exactly systemMessage: %v", source, len(decoded), decoded)
		}
		var message string
		if err := json.Unmarshal(decoded["systemMessage"], &message); err != nil {
			t.Fatalf("%s: hook output has no string systemMessage: %v", source, err)
		}
		for _, phrase := range requiredNoticePhrases {
			if !strings.Contains(message, phrase) {
				t.Errorf("%s: systemMessage is missing required phrase %q\nmessage: %q", source, phrase, message)
			}
		}
		if i == 0 {
			firstMessage = message
		} else if message != firstMessage {
			t.Errorf("hook emitted a different systemMessage for %q than for manual (the reminder is fixed, source-independent)", source)
		}
	}
}

// TestCodexPostCompactHookFiresFromUnrelatedCwdViaPluginRoot is the M1 regression: the
// shipped command must be plugin-root-absolute so it fires from ANY session cwd. It
// proves the defect concretely — the cwd-relative form (./hooks/<script>) does not
// resolve from an unrelated project directory, while the ${PLUGIN_ROOT}-resolved
// absolute form emits the systemMessage there.
func TestCodexPostCompactHookFiresFromUnrelatedCwdViaPluginRoot(t *testing.T) {
	root, group := loadShippedPostCompactHook(t)

	if !strings.HasPrefix(group.Hooks[0].Command, pluginRootToken+"/") {
		t.Fatalf("shipped hook command %q must begin with %s/ so Codex resolves it to the plugin root, not the session cwd", group.Hooks[0].Command, pluginRootToken)
	}

	absolute := resolveHookCommand(t, root, group.Hooks[0].Command)
	unrelated := t.TempDir() // a project dir that is NOT the plugin repo and has no ./hooks/
	payload := `{"hook_event_name":"PostCompact","trigger":"manual"}`

	// Negative: the cwd-relative form Codex would exec from the session cwd cannot find
	// the script from an unrelated cwd — the exact failure the plugin-root form avoids.
	relative := "./hooks/" + filepath.Base(absolute)
	relCmd := exec.Command(relative)
	relCmd.Dir = unrelated
	relCmd.Stdin = strings.NewReader(payload)
	if _, err := relCmd.Output(); err == nil {
		t.Fatalf("cwd-relative command %q unexpectedly resolved from an unrelated cwd; the plugin-root form would be untested", relative)
	}

	// Positive: the plugin-root-absolute command fires from the same unrelated cwd.
	out := runShippedHook(t, root, absolute, payload)
	var decoded map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("hook stdout is not valid JSON from an unrelated cwd: %v\nstdout: %q", err, out)
	}
	var message string
	if err := json.Unmarshal(decoded["systemMessage"], &message); err != nil {
		t.Fatalf("hook output from an unrelated cwd has no string systemMessage: %v", err)
	}
	for _, phrase := range requiredNoticePhrases {
		if !strings.Contains(message, phrase) {
			t.Errorf("systemMessage from an unrelated cwd is missing required phrase %q\nmessage: %q", phrase, message)
		}
	}
}

// TestCodexPostCompactHookHarmlessAbsenceMatrix is AC-4: across the present-ok,
// absent, disabled, and failing states the binding creates no Spacedock state file,
// spawns no background process, and mutates no workflow. Offline, the provable
// invariant is that the shipped hook is a stdout-only script that invokes no
// spacedock binary and writes nothing, so every degraded state is inert.
func TestCodexPostCompactHookHarmlessAbsenceMatrix(t *testing.T) {
	root, group := loadShippedPostCompactHook(t)
	command := resolveHookCommand(t, root, group.Hooks[0].Command)

	scriptBytes, err := os.ReadFile(command)
	if err != nil {
		t.Fatalf("read shipped hook script: %v", err)
	}
	script := string(scriptBytes)
	// The hook must not reach for the workflow: no mutation command and no file
	// redirection. A captain-facing UI cue only. ("spacedock:first-officer" appears
	// in the reminder text, so the ban is on mutation verbs, not the bare word.)
	for _, banned := range []string{"status --set", "state commit", "dispatch build", "spawn_agent", ">"} {
		if strings.Contains(script, banned) {
			t.Errorf("shipped hook script contains %q — a post-compact notice must not touch workflow state or write files", banned)
		}
	}

	cases := []struct {
		name    string
		handler []string // argv run in place of the hook; nil means "absent/disabled — nothing runs"
	}{
		{name: "present-ok", handler: []string{command}},
		{name: "absent", handler: nil},
		{name: "disabled", handler: nil},
		{name: "failing", handler: []string{"sh", "-c", "exit 3"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			work := t.TempDir()
			if tc.handler != nil {
				cmd := exec.Command(tc.handler[0], tc.handler[1:]...)
				cmd.Dir = work
				cmd.Env = append(os.Environ(), "HOME="+home, "CODEX_HOME="+home)
				cmd.Stdin = strings.NewReader(`{"hook_event_name":"PostCompact","trigger":"auto"}`)
				_ = cmd.Run() // a nonzero exit (failing) must still be harmless — do not fail the test on it
			}
			for label, dir := range map[string]string{"HOME": home, "cwd": work} {
				entries, err := os.ReadDir(dir)
				if err != nil {
					t.Fatalf("read %s: %v", label, err)
				}
				if len(entries) != 0 {
					t.Errorf("%s: the %s hook state left %d entr(ies) in %s — the binding must create no state file or process", tc.name, label, len(entries), label)
				}
			}
		})
	}
}

// runShippedHook runs the resolved hook command from an UNRELATED project directory
// (a fresh temp dir that is not the plugin repo and has no local ./hooks/), with
// PLUGIN_ROOT set as Codex sets it, and returns its stdout — asserting a zero exit and
// that it created nothing on disk. Running from an unrelated cwd is what exposes the
// cwd-relative resolution defect: only a plugin-root-absolute command fires here.
func runShippedHook(t *testing.T, root, command, stdinPayload string) string {
	t.Helper()
	work := t.TempDir()
	cmd := exec.Command(command)
	cmd.Dir = work
	cmd.Env = append(os.Environ(), "PLUGIN_ROOT="+root)
	cmd.Stdin = strings.NewReader(stdinPayload)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("hook command exited non-zero: %v", err)
	}
	entries, err := os.ReadDir(work)
	if err != nil {
		t.Fatalf("read work dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("hook wrote %d file(s) to its cwd — a post-compact notice must not write anything", len(entries))
	}
	return string(out)
}
