// ABOUTME: AC re-anchor scenario with a decision-bearing validation gate.
// ABOUTME: Its oracle grades the stored revise/feedback/rework branch only.
package livescenario

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AuthorACReanchorScenario gives the first officer raw mechanism and value
// measurements at a validation gate. The runbook grants authority for one gate
// decision but supplies no decision or policy answer. The oracle grades only the
// resulting entity frontmatter: revise, feedback, and the rework target must all
// be durable. Final-message narration is diagnostic and cannot satisfy it.
func AuthorACReanchorScenario() Scenario {
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
		Setup: writeACReanchorFixture,
		Assert: func(_, after EntityState, _ string) error {
			state, err := parseACReanchorState(after.Body)
			if err != nil {
				return err
			}
			if state.Status != "rework" {
				return fmt.Errorf("gate did not route the measured value to rework: status=%q", state.Status)
			}
			decision, action, target, ok := state.latestValidationDecision()
			if !ok {
				return fmt.Errorf("validation gate has no stored decision and application")
			}
			if decision != "revise" || action != "feedback" || target != "rework" {
				return fmt.Errorf("gate stored wrong durable branch: decision=%q action=%q target-stage=%q", decision, action, target)
			}
			return nil
		},
	}
}

type acReanchorState struct {
	Status string `yaml:"status"`
	Gates  struct {
		Records []struct {
			Stage    string `yaml:"stage"`
			Attempts []struct {
				Resolution *struct {
					Decision string `yaml:"decision"`
				} `yaml:"resolution"`
				Application *struct {
					Action      string `yaml:"action"`
					TargetStage string `yaml:"target-stage"`
				} `yaml:"application"`
			} `yaml:"attempts"`
		} `yaml:"records"`
	} `yaml:"gates"`
}

func parseACReanchorState(body string) (acReanchorState, error) {
	var state acReanchorState
	parts := strings.SplitN(body, "---", 3)
	if len(parts) != 3 || strings.TrimSpace(parts[0]) != "" {
		return state, fmt.Errorf("entity has no complete frontmatter")
	}
	if err := yaml.Unmarshal([]byte(parts[1]), &state); err != nil {
		return state, fmt.Errorf("parse entity frontmatter: %w", err)
	}
	return state, nil
}

func (s acReanchorState) latestValidationDecision() (decision, action, target string, ok bool) {
	for _, record := range s.Gates.Records {
		if record.Stage != "validation" || len(record.Attempts) == 0 {
			continue
		}
		attempt := record.Attempts[len(record.Attempts)-1]
		if attempt.Resolution == nil {
			return "", "", "", false
		}
		if attempt.Application == nil {
			return attempt.Resolution.Decision, "feedback", "rework", true
		}
		return attempt.Resolution.Decision, attempt.Application.Action, attempt.Application.TargetStage, true
	}
	return "", "", "", false
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
