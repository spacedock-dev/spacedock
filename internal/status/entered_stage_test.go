// ABOUTME: Entered-stage scheduling requires a committed, structurally complete
// ABOUTME: current-stage report before projecting or mutating to another stage.
package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

const enteredStageReadme = `---
commissioned-by: spacedock@1
entity-type: task
id-style: slug
state: .spacedock-state
stages:
  defaults:
    worktree: false
    concurrency: 2
  states:
    - name: backlog
      initial: true
    - name: implementation
      worktree: true
      concurrency: 3
    - name: validation
      gate: true
    - name: handoff
    - name: done
      terminal: true
---

# Entered Stage Workflow
`

const enteredStageEntity = `---
id: entered-task
title: Entered Task
status: implementation
completed:
verdict:
worktree:
---
# Entered Task
`

const completeImplementationReport = `
## Stage Report: implementation

- DONE: Produce the entered implementation.
  Commit abc123 contains the exercised implementation evidence.

### Summary

The entered implementation is complete and ready for independent validation.
`

type enteredStageRow struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Current  string `json:"current"`
	Next     string `json:"next"`
	Worktree string `json:"worktree"`
}

func TestEnteredStageProjectionRequiresCommittedCompleteReport(t *testing.T) {
	cases := []struct {
		name   string
		report string
	}{
		{"no report", ""},
		{"heading only", "\n## Stage Report: implementation\n"},
		{"empty checklist", "\n## Stage Report: implementation\n\n### Summary\n\nNot enough.\n"},
		{"blank checklist item", "\n## Stage Report: implementation\n\n- DONE:\n  Evidence without an obligation.\n\n### Summary\n\nNot enough.\n"},
		{"item without evidence", "\n## Stage Report: implementation\n\n- DONE: Produce the implementation.\n\n### Summary\n\nNot enough.\n"},
		{"failed item", "\n## Stage Report: implementation\n\n- FAILED: Produce the implementation.\n  The required behavior is still broken.\n\n### Summary\n\nNot enough.\n"},
		{"missing summary", "\n## Stage Report: implementation\n\n- DONE: Produce the implementation.\n  Commit abc123 contains evidence.\n"},
		{"empty summary", "\n## Stage Report: implementation\n\n- DONE: Produce the implementation.\n  Commit abc123 contains evidence.\n\n### Summary\n\n"},
		{"wrong stage", strings.ReplaceAll(completeImplementationReport, "Stage Report: implementation", "Stage Report: handoff")},
		{"stage-token near miss", strings.ReplaceAll(completeImplementationReport, "Stage Report: implementation", "Stage Report: implementation-notes")},
		{"later malformed masks older valid", completeImplementationReport + "\n## Stage Report: implementation (cycle 2)\n\n- DONE: Later report has no evidence.\n\n### Summary\n\nLater but malformed.\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def, _, _ := buildEnteredStageFixture(t, enteredStageEntity+tc.report)
			assertEnteredStageRows(t, def, enteredStageRow{
				ID: "entered-task", Slug: "entered-task",
				Current: "implementation", Next: "implementation", Worktree: "yes",
			})
		})
	}
	t.Run("uncommitted complete report", func(t *testing.T) {
		def, _, entity := buildEnteredStageFixture(t, enteredStageEntity)
		writeFile(t, entity, enteredStageEntity+completeImplementationReport)
		assertEnteredStageRows(t, def, enteredStageRow{
			ID: "entered-task", Slug: "entered-task",
			Current: "implementation", Next: "implementation", Worktree: "yes",
		})
	})
	t.Run("dirty committed report", func(t *testing.T) {
		def, _, entity := buildEnteredStageFixture(t, enteredStageEntity+completeImplementationReport)
		f, err := os.OpenFile(entity, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString("\nUncommitted entity dirt.\n"); err != nil {
			f.Close()
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		assertEnteredStageRows(t, def, enteredStageRow{
			ID: "entered-task", Slug: "entered-task",
			Current: "implementation", Next: "implementation", Worktree: "yes",
		})
	})
	t.Run("committed complete report ignores sibling dirt", func(t *testing.T) {
		def, state, _ := buildEnteredStageFixture(t, enteredStageEntity+completeImplementationReport)
		writeFile(t, filepath.Join(state, "unrelated-dirty-sibling.md"), "concurrent sibling dirt\n")
		assertEnteredStageRows(t, def, enteredStageRow{
			ID: "entered-task", Slug: "entered-task",
			Current: "implementation", Next: "validation", Worktree: "no",
		})
	})

	t.Run("skipped item with rationale is structurally complete", func(t *testing.T) {
		report := "\n## Stage Report: implementation\n\n- SKIPPED: Optional benchmark.\n  Not applicable because this change has no performance path.\n\n### Summary\n\nThe required implementation is complete.\n"
		def, _, _ := buildEnteredStageFixture(t, enteredStageEntity+report)
		assertEnteredStageRows(t, def, enteredStageRow{
			ID: "entered-task", Slug: "entered-task",
			Current: "implementation", Next: "validation", Worktree: "no",
		})
	})
	t.Run("spaced heading suffix remains accepted", func(t *testing.T) {
		report := strings.Replace(completeImplementationReport, "Stage Report: implementation", "Stage Report: implementation (cycle 2)", 1)
		def, _, _ := buildEnteredStageFixture(t, enteredStageEntity+report)
		assertEnteredStageRows(t, def, enteredStageRow{
			ID: "entered-task", Slug: "entered-task",
			Current: "implementation", Next: "validation", Worktree: "no",
		})
	})
}

func TestEnteredWorktreeStageFailureDiagnosticsAreDistinctAndByteClean(t *testing.T) {
	annotated := "\n## Stage Report: implementation\n\n- DONE: Produce the implementation.\n  Commit abc123 contains evidence.\n- DONE (annotation): Document the edge case.\n  The annotation must not hide this obligation.\n\n### Summary\n\nThe implementation is complete.\n"
	cases := []struct {
		name   string
		report string
		setup  func(t *testing.T, state, entity string)
		want   func(body, state, entity string) string
	}{
		{"missing report", "", nil, func(_, _, _ string) string {
			return `missing current-stage report: no heading whose first stage token is "implementation"; add "## Stage Report: implementation"`
		}},
		{"incomplete checklist", annotated, nil, func(body, _, _ string) string {
			line := strings.Count(body[:strings.Index(body, "- DONE (annotation):")], "\n") + 1
			return fmt.Sprintf(`incomplete current-stage report: line %d "- DONE (annotation): Document the edge case." is not canonical; use "- DONE: <item text>"`, line)
		}},
		{"untracked entity", completeImplementationReport, func(t *testing.T, state, _ string) {
			gitC(t, state, "rm", "-q", "--cached", "--", "entered-task.md")
		}, func(_, state, entity string) string {
			return fmt.Sprintf("untracked completion artifact: %s is not tracked in local Git root %s; add and commit that path", entity, state)
		}},
		{"dirty entity", completeImplementationReport, func(t *testing.T, _, entity string) {
			writeFile(t, entity, readBytes(t, entity)+"\nUncommitted entity dirt.\n")
		}, func(_, state, entity string) string {
			return fmt.Sprintf("dirty completion artifact: %s differs from local HEAD in Git root %s; commit that path (this guard does not require a remote push)", entity, state)
		}},
	}
	gotMessages := map[string]bool{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.Replace(enteredStageEntity+tc.report, "worktree:", "worktree: .worktrees/entered-task", 1)
			def, state, entity := buildEnteredStageFixture(t, body)
			if tc.setup != nil {
				tc.setup(t, state, entity)
			}
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			args := []string{"--workflow-dir", def, "--set", "entered-task", "status=validation"}
			if tc.name == "missing report" {
				args = append(args, "--force")
			}
			stdout, stderr, code := runNative(t, def, pinnedEnv(t), args...)
			if code != 1 {
				t.Fatalf("away mutation exit=%d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
			}
			if stdout != "" {
				t.Fatalf("away mutation emitted success stdout: %q", stdout)
			}
			want := fmt.Sprintf("Error: entity entered-task cannot change status away from entered stage %q: %s.\n",
				"implementation", tc.want(body, state, entity))
			if stderr != want {
				t.Fatalf("stderr=%q, want %q", stderr, want)
			}
			gotMessages[stderr] = true
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Fatalf("away mutation changed entity bytes\n--- before ---\n%s\n--- after ---\n%s", before, after)
			}
		})
	}
	if len(gotMessages) != 4 {
		t.Fatalf("distinct completion diagnostics=%d, want 4", len(gotMessages))
	}
}

func TestEnteredStageMutationControls(t *testing.T) {
	t.Run("same-stage dispatch mutation remains allowed", func(t *testing.T) {
		def, _, entity := buildEnteredStageFixture(t, enteredStageEntity)
		stdout, stderr, code := runNative(t, def, pinnedEnv(t),
			"--workflow-dir", def, "--set", "entered-task", "status=implementation", "started")
		if code != 0 || !strings.Contains(stdout, "started:") {
			t.Fatalf("same-stage dispatch exit=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		if !strings.Contains(readBytes(t, entity), "status: implementation") {
			t.Fatal("same-stage dispatch changed the current status")
		}
	})

	t.Run("unrelated mutation remains allowed", func(t *testing.T) {
		def, _, _ := buildEnteredStageFixture(t, enteredStageEntity)
		stdout, stderr, code := runNative(t, def, pinnedEnv(t),
			"--workflow-dir", def, "--set", "entered-task", "score=0.7")
		if code != 0 || !strings.Contains(stdout, "score:") {
			t.Fatalf("unrelated mutation exit=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})

	t.Run("committed completion unlocks successor", func(t *testing.T) {
		def, _, entity := buildEnteredStageFixture(t, enteredStageEntity+completeImplementationReport)
		stdout, stderr, code := runNative(t, def, pinnedEnv(t),
			"--workflow-dir", def, "--set", "entered-task", "status=validation")
		if code != 0 || !strings.Contains(stdout, "status: implementation -> validation") {
			t.Fatalf("completed transition exit=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		if !strings.Contains(readBytes(t, entity), "status: validation") {
			t.Fatal("completed transition did not update status")
		}
	})
	t.Run("normal worktree dispatch stays guarded until report is committed", func(t *testing.T) {
		def, state, entity := buildEnteredStageFixture(t, enteredStageEntity)
		stdout, stderr, code := runNative(t, def, pinnedEnv(t),
			"--workflow-dir", def, "--set", "entered-task",
			"status=implementation", "worktree=.worktrees/entered-task", "started")
		if code != 0 || !strings.Contains(stdout, "worktree:") {
			t.Fatalf("worktree dispatch exit=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		gitC(t, state, "add", "--", "entered-task.md")
		gitC(t, state, "commit", "-q", "-m", "dispatch entered stage", "--", "entered-task.md")

		before := readBytes(t, entity)
		stdout, stderr, code = runNative(t, def, pinnedEnv(t),
			"--workflow-dir", def, "--set", "entered-task", "status=validation")
		if code != 1 || stdout != "" || !strings.Contains(stderr, "Stage Report") {
			t.Fatalf("unfinished worktree transition exit=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		if after := readBytes(t, entity); after != before {
			t.Fatalf("unfinished worktree transition changed entity bytes\n--- before ---\n%s\n--- after ---\n%s", before, after)
		}

		writeFile(t, entity, before+completeImplementationReport)
		gitC(t, state, "add", "--", "entered-task.md")
		gitC(t, state, "commit", "-q", "-m", "complete entered stage", "--", "entered-task.md")
		stdout, stderr, code = runNative(t, def, pinnedEnv(t),
			"--workflow-dir", def, "--set", "entered-task", "status=validation")
		if code != 0 || !strings.Contains(stdout, "status: implementation -> validation") {
			t.Fatalf("completed worktree transition exit=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
	})
}

func TestEnteredStageLegacySuppressionControls(t *testing.T) {
	cases := []struct {
		name       string
		entityBody string
		want       []enteredStageRow
	}{
		{
			name:       "initial stage keeps successor projection",
			entityBody: strings.Replace(enteredStageEntity, "status: implementation", "status: backlog", 1),
			want: []enteredStageRow{{
				ID: "entered-task", Slug: "entered-task",
				Current: "backlog", Next: "implementation", Worktree: "yes",
			}},
		},
		{
			name:       "gate stage remains suppressed",
			entityBody: strings.Replace(enteredStageEntity, "status: implementation", "status: validation", 1),
		},
		{
			name:       "terminal stage remains suppressed",
			entityBody: strings.Replace(enteredStageEntity, "status: implementation", "status: done", 1),
		},
		{
			name:       "set worktree remains suppressed before first-entry rule",
			entityBody: strings.Replace(enteredStageEntity, "worktree:", "worktree: .worktrees/entered-task", 1),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def, _, _ := buildEnteredStageFixture(t, tc.entityBody)
			assertEnteredStageRows(t, def, tc.want...)
		})
	}
}

func buildEnteredStageFixture(t *testing.T, body string) (def, state, entity string) {
	t.Helper()
	def, state = buildSplitRoot(t, enteredStageReadme, map[string]string{"entered-task.md": body})
	entity = filepath.Join(state, "entered-task.md")
	testgit.InitRepo(t, state, "-q")
	gitC(t, state, "add", "--", "entered-task.md")
	gitC(t, state, "commit", "-q", "-m", "seed entered stage", "--", "entered-task.md")
	return def, state, entity
}

func assertEnteredStageRows(t *testing.T, def string, want ...enteredStageRow) {
	t.Helper()
	nextOut, nextErr, nextCode := runNative(t, def, pinnedEnv(t),
		"--workflow-dir", def, "--next", "--json")
	if nextCode != 0 {
		t.Fatalf("--next exit=%d stderr=%q", nextCode, nextErr)
	}
	bootOut, bootErr, bootCode := runNative(t, def, pinnedEnv(t),
		"--workflow-dir", def, "--boot", "--identify", "--json")
	if bootCode != 0 {
		t.Fatalf("--boot exit=%d stderr=%q", bootCode, bootErr)
	}
	var next, boot struct {
		Dispatchable []enteredStageRow `json:"dispatchable"`
	}
	if err := json.Unmarshal([]byte(nextOut), &next); err != nil {
		t.Fatalf("parse --next: %v\n%s", err, nextOut)
	}
	if err := json.Unmarshal([]byte(bootOut), &boot); err != nil {
		t.Fatalf("parse --boot: %v\n%s", err, bootOut)
	}
	if len(next.Dispatchable) != len(want) {
		t.Fatalf("--next rows=%+v, want %+v", next.Dispatchable, want)
	}
	for i := range want {
		if next.Dispatchable[i] != want[i] {
			t.Fatalf("--next row[%d]=%+v, want %+v", i, next.Dispatchable[i], want[i])
		}
	}
	nextJSON, err := json.Marshal(next.Dispatchable)
	if err != nil {
		t.Fatal(err)
	}
	bootJSON, err := json.Marshal(boot.Dispatchable)
	if err != nil {
		t.Fatal(err)
	}
	if string(bootJSON) != string(nextJSON) {
		t.Fatalf("boot dispatchable=%s, --next dispatchable=%s", bootJSON, nextJSON)
	}
}
