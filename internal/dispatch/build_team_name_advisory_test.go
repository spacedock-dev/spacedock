// ABOUTME: AC-1/AC-2/AC-3 — the --team-name-on-claude stderr advisory fires on the
// ABOUTME: legacy arm, is silent on the merged arm, and the --help docs the shape flags.
package dispatch

import (
	"encoding/json"
	"strings"
	"testing"
)

// teamNameAdvisoryMarker is a stable fragment of the legacy --team-name advisory.
// The gate keys on its presence/absence rather than the full sentence so a future
// wording tweak in claudeteam.LegacyTeamNameAdvisory does not break the test
// (matches the advisoryMarker convention in build_advisory_probe_test.go).
const teamNameAdvisoryMarker = "legacy TeamCreate-registry dispatch shape"

// TestBuildLegacyTeamNameAdvisory is AC-1's behavioral gate: on host=claude, a
// build that passes team_name emits exactly the legacy advisory on stderr while
// the same inputs WITHOUT team_name (the merged shape) emit none. The legacy arm
// also re-asserts the legacy envelope shape on stdout (team_name present,
// run_in_background absent) — AC-3 in the same arm: the advisory is stderr-only
// and does not flip the emitted dispatch envelope.
func TestBuildLegacyTeamNameAdvisory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd := writeGood(t, root)
	ep := writeFlatEntity(t, wd, "backlog", "")

	legacyStdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"host":           "claude",
		"team_name":      "fixture-team",
		"bare_mode":      false,
	}, nil)
	mergedStdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"host":           "claude",
		"bare_mode":      false,
		// team_name intentionally absent: the merged .178+ shape.
	}, nil)

	legacy := runNative(legacyStdin, "build", "--workflow-dir", wd)
	if legacy.exit != 0 {
		t.Fatalf("legacy --team-name build exit=%d, want 0\nstderr:\n%s", legacy.exit, legacy.stderr)
	}
	merged := runNative(mergedStdin, "build", "--workflow-dir", wd)
	if merged.exit != 0 {
		t.Fatalf("merged build exit=%d, want 0\nstderr:\n%s", merged.exit, merged.stderr)
	}

	// AC-1: advisory PRESENT on the --team-name arm, ABSENT on the merged arm.
	if !strings.Contains(legacy.stderr, teamNameAdvisoryMarker) {
		t.Errorf("host=claude --team-name build must emit the legacy advisory; stderr=%q", legacy.stderr)
	}
	if strings.Contains(merged.stderr, teamNameAdvisoryMarker) {
		t.Errorf("merged build (no --team-name) must emit NO legacy advisory; stderr=%q", merged.stderr)
	}
	// Exactly one advisory line on the legacy arm (not one per nothing).
	if got := strings.Count(legacy.stderr, teamNameAdvisoryMarker); got != 1 {
		t.Errorf("legacy advisory must appear exactly once; got %d:\n%s", got, legacy.stderr)
	}

	// AC-3 (same arm): the legacy envelope shape is unchanged by the advisory —
	// team_name present, run_in_background absent on the --team-name arm; and the
	// merged arm keeps run_in_background present, team_name absent.
	var legacyOut mergedBuildOutput
	if err := json.Unmarshal([]byte(legacy.stdout), &legacyOut); err != nil {
		t.Fatalf("legacy stdout is not build JSON: %v\n%s", err, legacy.stdout)
	}
	if legacyOut.TeamName == nil || *legacyOut.TeamName != "fixture-team" {
		t.Errorf("legacy --team-name envelope must keep team_name=%q; got %v", "fixture-team", legacyOut.TeamName)
	}
	if legacyOut.RunInBackground != nil {
		t.Errorf("legacy --team-name envelope must NOT gain run_in_background; got %v", *legacyOut.RunInBackground)
	}
	var mergedOut mergedBuildOutput
	if err := json.Unmarshal([]byte(merged.stdout), &mergedOut); err != nil {
		t.Fatalf("merged stdout is not build JSON: %v\n%s", err, merged.stdout)
	}
	if mergedOut.TeamName != nil {
		t.Errorf("merged envelope must omit team_name; got %q", *mergedOut.TeamName)
	}
	if mergedOut.RunInBackground == nil || !*mergedOut.RunInBackground {
		t.Errorf("merged envelope must keep run_in_background=true; got %v", mergedOut.RunInBackground)
	}
}

// TestBuildHelpDocumentsShapeFlags is AC-2's golden: dispatch build --help stdout
// documents the shape-selecting flags (--host, --team-name, --bare-mode), with the
// --team-name line naming the legacy shape and the auto-team default. Behavioral:
// it runs the command and reads stdout, not a source grep.
func TestBuildHelpDocumentsShapeFlags(t *testing.T) {
	res := runNative("", "build", "--help")
	if res.exit != 0 {
		t.Fatalf("dispatch build --help exit=%d, want 0\nstderr=%q", res.exit, res.stderr)
	}
	if res.stderr != "" {
		t.Fatalf("dispatch build --help stderr=%q, want empty", res.stderr)
	}
	assertContainsAll(t, res.stdout,
		"--host",
		"--team-name",
		"--bare-mode",
		"legacy TeamCreate-registry dispatch shape",
		"auto-team is the default",
	)
}
