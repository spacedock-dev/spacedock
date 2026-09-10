// ABOUTME: AC-2's focused fixture — the retired --team-name flag is a usage
// ABOUTME: error through the surviving CLI/file interface.
package dispatch

import (
	"strings"
	"testing"
)

// TestBuildTeamNameFlagRefused asserts `dispatch build --team-name NAME` exits 2
// with a usage error, the CLI-flag half of AC-2. --team-name selected the
// retired legacy TeamCreate-registry envelope; the binary now refuses the flag
// outright rather than silently ignoring it as an unrecognized token.
func TestBuildTeamNameFlagRefused(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd := writeGood(t, root)
	ep := writeFlatEntity(t, wd, "backlog", "")
	checklistPath := wd + "/checklist.txt"
	writeFile(t, checklistPath, "- a\n")

	for _, args := range [][]string{
		{"build", "--workflow-dir", wd, "--entity-path", ep, "--stage", "backlog", "--checklist-file", checklistPath, "--team-name", "fixture-team"},
		{"build", "--workflow-dir", wd, "--entity-path", ep, "--stage", "backlog", "--checklist-file", checklistPath, "--team-name=fixture-team"},
	} {
		native := runNative("", args...)
		if native.exit != 2 {
			t.Errorf("--team-name exit=%d, want 2 (usage error)\nargs=%v\nstderr=%q", native.exit, args, native.stderr)
		}
		if !strings.Contains(native.stderr, "--team-name") {
			t.Errorf("stderr does not name the refused flag: %q", native.stderr)
		}
	}
}
