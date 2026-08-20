// ABOUTME: AC-2's focused fixture — the retired --team-name flag is a usage
// ABOUTME: error, and a stdin team_name key is a byte-identical no-op.
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

// TestBuildStdinTeamNameKeyIsByteIdenticalNoOp is AC-2's envelope-equality half:
// the SAME request with and without a stdin team_name key emits a byte-identical
// envelope. The retired field degrades to the ignore-unknown-keys path rather
// than selecting any special shape.
func TestBuildStdinTeamNameKeyIsByteIdenticalNoOp(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd := writeGood(t, root)
	ep := writeFlatEntity(t, wd, "backlog", "")

	withoutTeamName := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
	}, nil)
	withTeamName := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		"team_name":      "fixture-team",
	}, nil)

	without := runNative(withoutTeamName, "build", "--workflow-dir", wd)
	with := runNative(withTeamName, "build", "--workflow-dir", wd)

	if without.exit != 0 || with.exit != 0 {
		t.Fatalf("exit without=%d with=%d, want 0/0\nstderr without=%q\nstderr with=%q",
			without.exit, with.exit, without.stderr, with.stderr)
	}
	if without.stdout != with.stdout {
		t.Errorf("stdout differs with a stdin team_name key present:\n--- without ---\n%s\n--- with ---\n%s",
			without.stdout, with.stdout)
	}
	if without.stderr != with.stderr {
		t.Errorf("stderr differs with a stdin team_name key present:\n--- without ---\n%q\n--- with ---\n%q",
			without.stderr, with.stderr)
	}
}
