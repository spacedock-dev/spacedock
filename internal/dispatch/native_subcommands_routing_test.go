// ABOUTME: routing regression — runtime-coupled inspection subcommands stay native.
// ABOUTME: Dispatch build emits exact full-path stage and standing fetch commands.
package dispatch

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// TestRuntimeCoupledSubcommandsRouteNative asserts the runtime-coupled
// subcommands are no longer rejected as deferred (the prerequisite seam returned
// exit 2 for them). Each now routes to a native handler: a recognized subcommand
// invoked with its required flag missing returns the native usage error (exit 2
// naming the subcommand), NOT the old "deferred to the claude-runtime-segregation
// surface" diagnostic. Re-adding a deferral would change this exit/message shape.
func TestRuntimeCoupledSubcommandsRouteNative(t *testing.T) {
	for _, sc := range []string{"context-budget", "list-standing", "show-standing", "spawn-standing-all"} {
		var out, errBuf bytes.Buffer
		code := Run(claudeteam.Probe, []string{sc}, strings.NewReader(""), &out, &errBuf)
		if code != 2 {
			t.Errorf("%s with no flags: exit=%d, want 2 (native usage error)", sc, code)
		}
		if strings.Contains(errBuf.String(), "deferred") {
			t.Errorf("%s still emits the deferral diagnostic: %q", sc, errBuf.String())
		}
		if !strings.Contains(errBuf.String(), sc) {
			t.Errorf("%s usage error does not name the subcommand: %q", sc, errBuf.String())
		}
	}
}

// TestBuildOmitsStandingFetchLineEvenUnderMods asserts dispatch build's fetch
// commands stay just show-stage-def even when the workflow declares a standing
// mod. The auto-injected show-standing fetch line only ever fired on the retired
// legacy team_name path (merged and bare dispatches always omitted it, per
// build.go's documented behavior); the pre-retirement sibling of this test,
// TestBuildEmitsStandingFetchLineUnderMods, drove that path with a team_name.
// The standing-teammate flow stays reachable directly via show-standing /
// spawn-standing-all.
func TestBuildOmitsStandingFetchLineEvenUnderMods(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	if err := os.MkdirAll(filepath.Join(root, "_mods"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "_mods", "helper.md"),
		"---\nstanding: true\nname: helper\n---\n## Hook: startup\n- name: helper\n## Agent Prompt\ny\n")
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "backlog",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
		"host":           "claude",
	}, nil)

	var out, errBuf bytes.Buffer
	if code := RunWithLauncher(claudeteam.Probe, testWorkflowLauncher, []string{"build", "--workflow-dir", root}, strings.NewReader(stdin), &out, &errBuf); code != 0 {
		t.Fatalf("build exit=%d stderr=%q", code, errBuf.String())
	}

	var env struct {
		FetchCommands []string `json:"fetch_commands"`
	}
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("build output not JSON: %v\n%s", err, out.String())
	}
	want := testWorkflowLauncher + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage backlog"
	if len(env.FetchCommands) != 1 || env.FetchCommands[0] != want {
		t.Fatalf("expected exactly one full-path show-stage-def command even with a standing mod declared; got %v", env.FetchCommands)
	}
}

// TestBuildOmitsStandingFetchLineWithoutMods asserts the same one-fetch-command
// shape holds for a workflow with no standing mods at all.
func TestBuildOmitsStandingFetchLineWithoutMods(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "backlog",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
		"host":           "claude",
	}, nil)

	var out, errBuf bytes.Buffer
	if code := RunWithLauncher(claudeteam.Probe, testWorkflowLauncher, []string{"build", "--workflow-dir", root}, strings.NewReader(stdin), &out, &errBuf); code != 0 {
		t.Fatalf("build exit=%d stderr=%q", code, errBuf.String())
	}

	var env struct {
		FetchCommands []string `json:"fetch_commands"`
	}
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("build output not JSON: %v\n%s", err, out.String())
	}
	want := testWorkflowLauncher + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage backlog"
	if len(env.FetchCommands) != 1 || env.FetchCommands[0] != want {
		t.Fatalf("expected exactly one full-path show-stage-def command; got %v", env.FetchCommands)
	}
}
