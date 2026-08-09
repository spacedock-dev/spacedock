// ABOUTME: AC-2 content-contract parity — walks every row of the entity's
// ABOUTME: content-contract table against the emitted --advance file body + prompt.
package dispatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// contractRow is one row of the pinned content contract (docs/dev/.spacedock-
// state/reuse-advance-file-pointer-dispatch.md `## Content contract`): the
// current-template element and the check that the new file/pointer shape
// still carries it.
type contractRow struct {
	element string
	check   func(t *testing.T, prompt, body string, f advanceParityFixture)
}

// advanceParityFixture is the fixed set of values a contract-row check needs
// to build its expected string.
type advanceParityFixture struct {
	stage           string
	entityTitle     string
	entityPath      string
	workflowDir     string
	checklist       []string
	feedbackContext string
}

var advanceContractRows = []contractRow{
	{
		element: "Next stage name (`Advancing to next stage: {next_stage_name}`)",
		check: func(t *testing.T, prompt, body string, f advanceParityFixture) {
			t.Helper()
			want := "Advancing to next stage: " + f.stage
			if !strings.Contains(prompt, want) {
				t.Errorf("pointer message missing %q: %q", want, prompt)
			}
			if !strings.Contains(body, "## "+want) {
				t.Errorf("file header missing %q: body=%q", "## "+want, body)
			}
		},
	},
	{
		element: "`### Stage definition:` — full README subsection verbatim (via show-stage-def fetch line)",
		check: func(t *testing.T, prompt, body string, f advanceParityFixture) {
			t.Helper()
			want := fmt.Sprintf("%s dispatch show-stage-def --workflow-dir %s --stage %s", shlexQuote(testWorkflowLauncher), shlexQuote(f.workflowDir), shlexQuote(f.stage))
			if !strings.Contains(body, want) {
				t.Errorf("fetch commands missing exact show-stage-def line %q: body=%q", want, body)
			}
		},
	},
	{
		element: "`### Completion checklist`",
		check: func(t *testing.T, prompt, body string, f advanceParityFixture) {
			t.Helper()
			if !strings.Contains(body, "### Completion checklist") {
				t.Errorf("file missing '### Completion checklist' heading: body=%q", body)
			}
			for _, item := range f.checklist {
				if !strings.Contains(body, item) {
					t.Errorf("file missing checklist item %q: body=%q", item, body)
				}
			}
		},
	},
	{
		element: "`Continue working on {entity title} at {entity_file_path}`",
		check: func(t *testing.T, prompt, body string, f advanceParityFixture) {
			t.Helper()
			if !strings.Contains(body, "You are continuing work on: "+f.entityTitle) {
				t.Errorf("file header missing entity title continuation line: body=%q", body)
			}
			if !strings.Contains(body, "Continue working on the entity at "+f.entityPath+".") {
				t.Errorf("file missing continue-on-entity line for %q: body=%q", f.entityPath, body)
			}
		},
	},
	{
		element: "`Commit before sending your completion message` (now pinning the next-stage Done: wording)",
		check: func(t *testing.T, prompt, body string, f advanceParityFixture) {
			t.Helper()
			if !strings.Contains(body, "after all commits and stage report writes are done") {
				t.Errorf("file missing commit-before-signal clause: body=%q", body)
			}
			wantDone := fmt.Sprintf("Done: %s completed %s.", f.entityTitle, f.stage)
			if !strings.Contains(body, wantDone) {
				t.Errorf("file missing next-stage Done: wording %q: body=%q", wantDone, body)
			}
		},
	},
	{
		element: "Feedback context when re-entering a feedback-to stage (`### Feedback from prior review`)",
		check: func(t *testing.T, prompt, body string, f advanceParityFixture) {
			t.Helper()
			if f.feedbackContext == "" {
				if strings.Contains(body, "### Feedback from prior review") {
					t.Errorf("file carries a Feedback section with no feedback_context supplied: body=%q", body)
				}
				return
			}
			if !strings.Contains(body, "### Feedback from prior review") {
				t.Errorf("file missing '### Feedback from prior review' heading: body=%q", body)
			}
			if !strings.Contains(body, f.feedbackContext) {
				t.Errorf("file missing feedback context text %q: body=%q", f.feedbackContext, body)
			}
		},
	},
}

// TestBuildAdvanceContentContractParity is AC-2: every row of the pinned
// content-contract table is asserted against the emitted advance file body +
// pointer prompt, for both a plain advance and a feedback-reflow advance (the
// only row that varies between the two fixtures is the feedback-context row).
func TestBuildAdvanceContentContractParity(t *testing.T) {
	cases := []struct {
		name            string
		stage           string
		feedbackReflow  bool
		feedbackContext string
	}{
		{name: "plain", stage: "validation"},
		{name: "feedback-reflow", stage: "implementation", feedbackReflow: true, feedbackContext: "REJECTED: fix the widget."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
			worktreeRel := ".worktrees/spacedock-ensign-thing"
			if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
				t.Fatal(err)
			}
			entityPath := filepath.Join(root, "thing.md")
			entityTitle := "Thing"
			writeFile(t, entityPath, entityFM(entityTitle, tc.stage, worktreeRel))
			gitInit(t, root)

			checklist := []string{"- validate the thing", "- confirm no regressions"}
			stdinFields := map[string]any{
				"schema_version": 2,
				"entity_path":    entityPath,
				"workflow_dir":   root,
				"stage":          tc.stage,
				"checklist":      checklist,
				"team_name":      "fixture-team",
				"bare_mode":      false,
				"advance":        true,
			}
			if tc.feedbackReflow {
				stdinFields["is_feedback_reflow"] = true
				stdinFields["feedback_context"] = tc.feedbackContext
			}
			stdin := mergeStdin(stdinFields, nil)

			native := runNative(stdin, "build", "--workflow-dir", root)
			if native.exit != 0 {
				t.Fatalf("advance build exit=%d stderr=%s", native.exit, native.stderr)
			}
			var out struct {
				Prompt string `json:"prompt"`
			}
			if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
				t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
			}
			body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

			// The entity is worktree-stamped and non-split-root, so the file
			// rewrites entity references into the worktree copy (identical
			// resolution to fresh dispatch) — the row checks below expect this
			// path, not the project-root entity_path passed on stdin.
			worktreeEntityPath := filepath.Join(root, worktreeRel, "thing.md")

			fixture := advanceParityFixture{
				stage:           tc.stage,
				entityTitle:     entityTitle,
				entityPath:      worktreeEntityPath,
				workflowDir:     root,
				checklist:       []string{"validate the thing", "confirm no regressions"},
				feedbackContext: tc.feedbackContext,
			}

			for _, row := range advanceContractRows {
				t.Run(row.element, func(t *testing.T) {
					row.check(t, out.Prompt, body, fixture)
				})
			}
		})
	}
}
