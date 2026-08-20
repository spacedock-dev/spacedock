// ABOUTME: AC-2 content-contract parity through structured advance inputs.
// ABOUTME: Narrative wording remains covered only by regression snapshots.
package dispatch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// advanceParityFixture is the fixed set of structured values compared with the
// bounded sections of one generated advance artifact.
type advanceParityFixture struct {
	stage           string
	entityTitle     string
	entityPath      string
	workflowDir     string
	checklist       []string
	feedbackContext string
}

// TestBuildAdvanceContentContractParity compares the structured fetch, checklist,
// and optional feedback sections with their source inputs. It deliberately does
// not infer behavior from surrounding assignment narration.
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
				Fetch []string `json:"fetch_commands"`
			}
			if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
				t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
			}
			body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))
			sections := dispatchArtifactSections(t, body)

			// The entity is worktree-stamped and non-split-root, so the file
			// rewrites entity references into the worktree copy (identical
			// resolution to fresh dispatch).
			worktreeEntityPath := filepath.Join(root, worktreeRel, "thing.md")
			fixture := advanceParityFixture{
				stage:           tc.stage,
				entityTitle:     entityTitle,
				entityPath:      worktreeEntityPath,
				workflowDir:     root,
				checklist:       checklist,
				feedbackContext: tc.feedbackContext,
			}

			wantFetch := testWorkflowLauncher + " dispatch show-stage-def --workflow-dir " + shlexQuote(fixture.workflowDir) + " --stage " + fixture.stage
			if len(out.Fetch) != 1 || out.Fetch[0] != wantFetch || sections["Fetch commands"] != wantFetch {
				t.Fatalf("advance fetch differs from structured source: fetch=%v section=%q", out.Fetch, sections["Fetch commands"])
			}
			if sections["Completion checklist"] != strings.Join(fixture.checklist, "\n") {
				t.Fatalf("advance checklist differs from structured source: %q", sections["Completion checklist"])
			}
			feedback, exists := sections["Feedback from prior review"]
			if tc.feedbackReflow {
				if !exists || feedback != fixture.feedbackContext {
					t.Fatalf("feedback section=%q exists=%v, want exact input %q", feedback, exists, fixture.feedbackContext)
				}
			} else if exists {
				t.Fatalf("plain advance unexpectedly emitted a feedback section: %q", feedback)
			}
		})
	}
}
