// ABOUTME: team_name path-safety — a path-unsafe or over-long team_name fails
// ABOUTME: with a clean error before the dispatch-file path is constructed.
package dispatch

import (
	"path/filepath"
	"strings"
	"testing"
)

// teamNameStdin builds a well-formed build request for the good README's backlog
// stage with the given team_name, leaving every other guard satisfied so the
// team_name path-safety check is the one under test.
func teamNameStdin(t *testing.T, root, teamName string) (workflowDir, stdin string) {
	t.Helper()
	wd := writeGood(t, root)
	ep := writeFlatEntity(t, wd, "backlog", "")
	return wd, mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"team_name":      teamName,
		"bare_mode":      false,
	}, nil)
}

// TestBuildTeamNamePathSafety asserts a path-unsafe team_name (path separator,
// dot-dot traversal, or odd characters) fails with a clean exit-1 error rather
// than building an unsanitized /tmp path. The derived name beside it is already
// validated this way; team_name gets the same treatment before path construction.
func TestBuildTeamNamePathSafety(t *testing.T) {
	cases := []struct {
		name     string
		teamName string
	}{
		{"dot-dot-traversal", "../escape"},
		{"path-separator", "a/b"},
		{"leading-slash", "/abs"},
		{"uppercase", "BadTeam"},
		{"space", "bad team"},
		{"trailing-hyphen", "team-"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()
			wd, stdin := teamNameStdin(t, root, tc.teamName)

			native := runNative(stdin, "build", "--workflow-dir", wd)

			if native.exit != 1 {
				t.Errorf("exit native=%d, want 1 for team_name %q\nstdout:\n%s\nstderr:\n%s",
					native.exit, tc.teamName, native.stdout, native.stderr)
			}
			if !strings.HasPrefix(native.stderr, "error: ") {
				t.Errorf("stderr lacks clean error prefix for team_name %q:\n%q", tc.teamName, native.stderr)
			}
		})
	}
}

// TestBuildTeamNameCombinedLengthCap asserts that a team_name and a derived name
// each short enough to pass their per-name length cap still fail when their
// combination pushes the on-disk filename past the filesystem name limit. The
// derived name is lengthened via a long entity slug.
func TestBuildTeamNameCombinedLengthCap(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd := writeGood(t, root)
	// A long (but per-name-valid) entity slug stretches the derived name so the
	// combined filename overshoots dispatchFileNameMaxLen.
	ep := filepath.Join(wd, strings.Repeat("a", 100)+".md")
	writeFile(t, ep, entityFM("Thing", "backlog", ""))
	teamName := strings.Repeat("b", nameMaxLen) // valid in isolation (== nameMaxLen)
	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"team_name":      teamName,
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", wd)

	if native.exit != 1 {
		t.Errorf("exit native=%d, want 1 for combined-length overflow\nstderr:\n%s", native.exit, native.stderr)
	}
	if !strings.HasPrefix(native.stderr, "error: ") {
		t.Errorf("stderr lacks clean error prefix:\n%q", native.stderr)
	}
}

// TestBuildTeamNameValidPath asserts a path-safe team_name still keys the
// dispatch filename as {teamName}-{derivedName}.md — the happy path is unchanged.
func TestBuildTeamNameValidPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd, stdin := teamNameStdin(t, root, "fixture-team")

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("exit native=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
	}
	got := dispatchFilePathFromStdout(t, native.stdout)
	want := filepath.Join(dispatchFileDir, "fixture-team-spacedock-ensign-thing-backlog.md")
	if got != want {
		t.Errorf("dispatch_file_path = %q, want %q", got, want)
	}
}
