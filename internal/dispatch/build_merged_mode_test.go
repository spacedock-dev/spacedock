// ABOUTME: merged-mode dispatch — Claude .178+ emits name present + team_name
// ABOUTME: absent + run_in_background; the only shape a non-bare claude build emits.
package dispatch

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// mergedBuildOutput is the build stdout JSON, decoded to the fields the merged
// shape constrains. Name/TeamName/RunInBackground are pointers so absent (key
// omitted) is distinguishable from present-with-zero-value.
type mergedBuildOutput struct {
	Name            *string `json:"name"`
	TeamName        *string `json:"team_name"`
	RunInBackground *bool   `json:"run_in_background"`
	DispatchFile    string  `json:"dispatch_file_path"`
}

// mergedStdin builds a well-formed non-bare claude request with NO team_name —
// the merged dispatch shape. Every other guard is satisfied so the Rule-8
// merged-path behavior is the one under test.
func mergedStdin(t *testing.T, root, stage string) (workflowDir, stdin string) {
	t.Helper()
	wd := writeGood(t, root)
	ep := writeFlatEntity(t, wd, stage, "")
	return wd, mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          stage,
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		// team_name intentionally absent: this is the merged .178+ shape.
	}, nil)
}

// TestBuildMergedModeEmission is AC-1's fixture: on host=claude, not bare, no
// team_name, the build emits name present + team_name absent + run_in_background
// true, with NO "team mode requires team_name" error (the new Rule-8 path).
func TestBuildMergedModeEmission(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd, stdin := mergedStdin(t, root, "backlog")

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("merged build exit=%d, want 0\nstdout:\n%s\nstderr:\n%s",
			native.exit, native.stdout, native.stderr)
	}
	if strings.Contains(native.stderr, "team mode requires team_name") {
		t.Fatalf("merged build must not emit the legacy Rule-8 error\nstderr:\n%s", native.stderr)
	}

	var out mergedBuildOutput
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	if out.Name == nil || *out.Name == "" {
		t.Errorf("merged build must emit a non-empty name (the reuse-advance + SendMessage handle); got %v", out.Name)
	}
	if out.TeamName != nil {
		t.Errorf("merged build must omit team_name (.178+ has no TeamCreate name); got %q", *out.TeamName)
	}
	if out.RunInBackground == nil || !*out.RunInBackground {
		t.Errorf("merged build must emit run_in_background=true (the worker→lead background channel); got %v", out.RunInBackground)
	}
}

// TestBuildMergedModeCompletionSignal is AC-6's fixture: the merged dispatch
// body carries the completion-signal block pinned to the single worker→lead
// target "team-lead" (matching the ensign runtime), not "main", and not absent.
func TestBuildMergedModeCompletionSignal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd, stdin := mergedStdin(t, root, "backlog")

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("merged build exit=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	if !strings.Contains(body, "### Completion Signal") {
		t.Fatalf("merged dispatch body must carry a completion-signal block (the worker→lead signal):\n%s", body)
	}
	if !strings.Contains(body, `SendMessage(to="team-lead"`) {
		t.Errorf("merged completion signal must target the pinned single reply name \"team-lead\":\n%s", body)
	}
	if strings.Contains(body, `SendMessage(to="main"`) {
		t.Errorf("merged completion signal must pin ONE target — \"main\" must not appear alongside \"team-lead\":\n%s", body)
	}
}

// TestBuildStdinTeamNameIgnored is AC-2's fixture: a stdin team_name key is an
// unrecognized field on the retired legacy path — the build emits the SAME
// merged shape (name present, team_name absent, run_in_background true) with or
// without it, byte-identical to TestBuildMergedModeEmission's assertions.
func TestBuildStdinTeamNameIgnored(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd := writeGood(t, root)
	ep := writeFlatEntity(t, wd, "backlog", "")
	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"team_name":      "fixture-team",
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("build exit=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
	}
	var out mergedBuildOutput
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	if out.Name == nil || *out.Name == "" {
		t.Errorf("build must emit a non-empty name regardless of the ignored team_name key; got %v", out.Name)
	}
	if out.TeamName != nil {
		t.Errorf("a stdin team_name key must not resurrect the retired field; got %q", *out.TeamName)
	}
	if out.RunInBackground == nil || !*out.RunInBackground {
		t.Errorf("build must still emit run_in_background=true with team_name present in stdin; got %v", out.RunInBackground)
	}
}

// TestBuildBareModeUnchanged guards that bare mode (no name, no team_name, no
// run_in_background — the blocking sequential shape) is untouched by the merged
// path: bare omits all three keys, distinct from merged's name+run_in_background.
func TestBuildBareModeUnchanged(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd := writeGood(t, root)
	ep := writeFlatEntity(t, wd, "backlog", "")
	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"bare_mode":      true,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("bare build exit=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
	}
	var out mergedBuildOutput
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	if out.Name != nil || out.TeamName != nil || out.RunInBackground != nil {
		t.Errorf("bare build must omit name/team_name/run_in_background; got name=%v team_name=%v run_in_background=%v", out.Name, out.TeamName, out.RunInBackground)
	}
}

// TestBuildMergedModeDispatchFileDisambiguator is the scope-finding guard: with
// no team_name to prepend, two concurrent merged dispatches of one slug+stage
// must NOT alias the same dispatch file (the stale-pointer hazard). The merged
// path keys the filename on the session id so concurrent FOs get distinct files.
func TestBuildMergedModeDispatchFileDisambiguator(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd, stdin := mergedStdin(t, root, "backlog")

	t.Setenv("CLAUDE_CODE_SESSION_ID", "sessionaaa")
	a := runNative(stdin, "build", "--workflow-dir", wd)
	if a.exit != 0 {
		t.Fatalf("merged build (session A) exit=%d\nstderr:\n%s", a.exit, a.stderr)
	}
	t.Setenv("CLAUDE_CODE_SESSION_ID", "sessionbbb")
	b := runNative(stdin, "build", "--workflow-dir", wd)
	if b.exit != 0 {
		t.Fatalf("merged build (session B) exit=%d\nstderr:\n%s", b.exit, b.stderr)
	}

	pathA := dispatchFilePathFromStdout(t, a.stdout)
	pathB := dispatchFilePathFromStdout(t, b.stdout)
	if pathA == pathB {
		t.Errorf("two concurrent merged sessions aliased the same dispatch file %q — a stale prior dispatch's pointer can be read", pathA)
	}
	if !strings.Contains(filepath.Base(pathA), "sessionaaa") {
		t.Errorf("merged dispatch filename must embed the session-id disambiguator; got %q", filepath.Base(pathA))
	}
}
