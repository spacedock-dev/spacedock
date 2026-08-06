// ABOUTME: AC re-anchor scenario with a decision-bearing validation gate.
// ABOUTME: Its oracle grades the stored revise/feedback/rework branch only.
package livescenario

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/gates"
	statuspkg "github.com/spacedock-dev/spacedock/internal/status"
)

// AuthorACReanchorScenario gives the first officer raw mechanism and value
// measurements at a validation gate. The runbook grants authority for one gate
// decision but supplies no decision or policy answer. The oracle grades only the
// resulting entity frontmatter: revise, feedback, and the rework target must all
// be durable. Final-message narration is diagnostic and cannot satisfy it.
func AuthorACReanchorScenario() Scenario {
	var entityPath, workflowDir string
	return Scenario{
		Name: "ac-reanchor/means-pass-value-regressed",
		Runbook: strings.Join([]string{
			"Use $spacedock:first-officer for this whole run.",
			"",
			"Workflow directory: .",
			"Process only `ac2-design-proof`, which is parked at its validation gate.",
			"You have the conn to resolve exactly this one gate from the existing evidence. Prepare and present the gate review, record the decision, and apply its configured route until the entity reaches the resulting parked stage.",
			"Then stop. Do not dispatch workers, edit the report, archive the entity, or perform work in the resulting stage.",
		}, "\n"),
		Setup: func(dir string) (string, error) {
			workflowDir = dir
			var err error
			entityPath, err = writeACReanchorFixture(dir)
			return entityPath, err
		},
		Assert: func(_, after EntityState, _ string) error {
			doc, _, err := gates.Read(entityPath)
			if err != nil {
				return fmt.Errorf("validation gate is not canonical: %w", err)
			}
			summary := gates.CurrentSummary(doc, "validation")
			if summary.State != "closed" || summary.Decision != "revise" || summary.Application != "" {
				return fmt.Errorf("validation gate stored wrong latest decision: state=%q decision=%q application=%q", summary.State, summary.Decision, summary.Application)
			}
			feedbackTo, err := gates.DeclaredReworkTarget(workflowDir, "validation")
			if err != nil {
				return fmt.Errorf("validation feedback route is invalid: %w", err)
			}
			if feedbackTo != "rework" {
				return fmt.Errorf("validation gate has wrong feedback route: feedback-to=%q", feedbackTo)
			}
			if status := statuspkg.ParseFrontmatterData([]byte(after.Body))["status"]; status != feedbackTo {
				return fmt.Errorf("gate did not apply the feedback route: status=%q feedback-to=%q", status, feedbackTo)
			}
			return nil
		},
	}
}

// writeACReanchorFixture creates and commits a standalone workflow. Gate record
// and consume therefore have a real Git root in which to persist their writes.
func writeACReanchorFixture(dir string) (string, error) {
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(acReanchorReadme), 0o644); err != nil {
		return "", err
	}
	entityPath := filepath.Join(dir, "ac2-design-proof.md")
	if err := os.WriteFile(entityPath, []byte(acReanchorEntity), 0o644); err != nil {
		return "", err
	}
	commands := [][]string{
		{"init", "-q"},
		{"add", "--", "README.md", "ac2-design-proof.md"},
		{"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", "init"},
	}
	for _, args := range commands {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
		}
	}
	return entityPath, nil
}

const acReanchorReadme = `---
entity-type: task
id-style: slug
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: rework
      initial: true
    - name: validation
      gate: true
      feedback-to: rework
    - name: accepted
    - name: done
      terminal: true
---
# Contract Measurement Fixture

### rework

Parked target stage.

### validation

Review the recorded acceptance evidence.

- **Outputs:** A gate decision from the recorded evidence.

### accepted

Parked target stage.

### done

Terminal state.
`

const acReanchorEntity = `---
id: ac2-design-proof
title: Contract Measurement
status: validation
completed:
verdict:
worktree:
---
# Contract Measurement

## Acceptance criteria

**AC-1 - The prose section uses the new pattern.**
Verified by: README "Completion and Gates" section was rewritten.

**AC-2 - Contract size meets its measured target.**
Verified by: baseline 10,000 bytes; target 8,000 bytes; actual 10,200 bytes.

## Stage Report: validation

- DONE: The prose section uses the new pattern
  The "Completion and Gates" section was rewritten.
- DONE: The contract size was measured
  Baseline 10,000 bytes; target 8,000 bytes; actual 10,200 bytes.

### Summary

The recorded evidence is ready for gate review.
`
