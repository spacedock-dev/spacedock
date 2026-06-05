// ABOUTME: Codex runtime adapter contract tests — the skill surface loads
// ABOUTME: dedicated Codex FO/ensign adapters with mailbox wait semantics.
package hostneutrality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCodexRuntimeAdaptersAreLoadable is a code-bound invariant: each runtime
// SKILL.md branches on the SAME host env var the binary reads
// (CODEX_THREAD_ID, AST-extracted from internal/dispatch/build.go) and loads its
// Codex adapter file. The env-var expectation comes from the binary's
// host-derivation code, not a literal frozen against the skill — if the binary
// stops reading CODEX_THREAD_ID, or the skill stops branching on it, the two
// diverge and this reds. The adapter-content tokens are the remaining
// text-consistency portion.
func TestCodexRuntimeAdaptersAreLoadable(t *testing.T) {
	markCodeBoundInvariant(t, "hostEnvVar CODEX_THREAD_ID (internal/dispatch/build.go host-derivation)")
	envVar := hostEnvVar(t, "CODEX_THREAD_ID")
	if envVar == "" {
		t.Fatal("the binary no longer reads CODEX_THREAD_ID for host derivation — the env var the skill must branch on is gone")
	}
	root := filepath.Join("..", "..")
	cases := []struct {
		name      string
		skillPath string
		adapter   string
	}{
		{
			name:      "first-officer",
			skillPath: filepath.Join(root, "skills", "first-officer", "SKILL.md"),
			adapter:   filepath.Join(root, "skills", "first-officer", "references", "codex-first-officer-runtime.md"),
		},
		{
			name:      "ensign",
			skillPath: filepath.Join(root, "skills", "ensign", "SKILL.md"),
			adapter:   filepath.Join(root, "skills", "ensign", "references", "codex-ensign-runtime.md"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			skill := readText(t, tc.skillPath)
			adapterBase := filepath.Base(tc.adapter)
			if !strings.Contains(skill, envVar) || !strings.Contains(skill, adapterBase) {
				t.Fatalf("%s SKILL.md must branch on %s (the binary's host-derivation env var) and load %s:\n%s", tc.name, envVar, adapterBase, skill)
			}

			body := readText(t, tc.adapter)
			for _, want := range []string{"## Awaiting Completion", "## Dispatch", "send_input", "## Completion Signal", "Codex declares none"} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s adapter missing %q:\n%s", tc.name, want, body)
				}
			}
		})
	}
}

// TestCodexAwaitingCompletionPinsMailboxSemantics is a non-AC text-consistency
// lint: the Codex FO adapter carries the mailbox-wait clauses (async final-status
// notification, wait_agent-timeout-is-normal, do-not-poll). Per the proof policy
// this presence check does NOT prove the FO obeys the mailbox semantics; the
// behavior is exercised by the Codex live runner's awaiting-completion path
// (codex_live_runner_test.go / codex_idle_notification_test.go). This lint guards
// the clauses being dropped from the adapter.
func TestCodexAwaitingCompletionPinsMailboxSemantics(t *testing.T) {
	markNonAC(t, "Codex live runner awaiting-completion path (internal/ensigncycle codex_live_runner + codex_idle_notification)")
	body := readText(t, filepath.Join("..", "..", "skills", "first-officer", "references", "codex-first-officer-runtime.md"))
	for _, want := range []string{
		"async final-status notification in the FO mailbox",
		"wait_agent timeout return is normal",
		"do not poll, re-dispatch, or close the worker",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("Codex FO adapter must pin waiting clause %q:\n%s", want, body)
		}
	}
}

func readText(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
