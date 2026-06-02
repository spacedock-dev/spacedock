// ABOUTME: Codex runtime adapter contract tests — the skill surface loads
// ABOUTME: dedicated Codex FO/ensign adapters with mailbox wait semantics.
package hostneutrality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexRuntimeAdaptersAreLoadable(t *testing.T) {
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
			if !strings.Contains(skill, "CODEX_THREAD_ID") || !strings.Contains(skill, adapterBase) {
				t.Fatalf("%s SKILL.md must branch on CODEX_THREAD_ID and load %s:\n%s", tc.name, adapterBase, skill)
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

func TestCodexAwaitingCompletionPinsMailboxSemantics(t *testing.T) {
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
