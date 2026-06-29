// ABOUTME: AC-2 real behavioral proof — livescenario test exercising
// ABOUTME: actual FO agent against means-only + regressed-value fixture.
package livescenario

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AuthorACReanchorScenario creates the scenario for AC-2 real behavioral proof.
// This is authored here (in the livescenario package) and will be invoked by
// the ensigncycle live adapter (claudeRunnerAdapter) to run against a real FO.
func AuthorACReanchorScenario(t *testing.T) Scenario {
	return Scenario{
		Name: "ac2-reanchor-means-only-regressed-reject",
		Runbook: "Inspect the ac2-design-proof-fixture entity. " +
			"It has AC-1 (mechanism-only: prose was updated) and AC-2 (end-value: contract shrink 20% but actually grew 2%). " +
			"Apply the AC cross-check with the re-anchor rule from the updated first-officer-shared-core.md. " +
			"Mechanism-only AC is satisfied only when its end-value AC is satisfied. " +
			"Since AC-2 (the end-value) is regressed, reject the entity at the gate. " +
			"Present your gate review with the rejection reasoning.",
		Setup: func(dir string) (string, error) {
			// Stage the fixture from the workflow state checkout
			sourceFixture := "docs/dev/.spacedock-state/ac2-design-proof-fixture.md"
			sourcePath := filepath.Join(dir, sourceFixture)

			// Read fixture from the state checkout
			fixtureBody, err := os.ReadFile(sourcePath)
			if err != nil {
				// Fallback: fixture might not exist yet in this context
				// Write a minimal fixture for this test run
				return writeMinimalACReanchorFixture(dir)
			}

			// Write fixture to the temp dir for the FO to process
			entityPath := filepath.Join(dir, "ac2-design-proof-fixture.md")
			if err := os.WriteFile(entityPath, fixtureBody, 0o644); err != nil {
				return "", err
			}
			return entityPath, nil
		},
		Assert: func(before, after EntityState, observed string) error {
			// Durable outcome 1: Entity verdict must change to REJECTED
			hasVerdictBefore := strings.Contains(before.Body, "verdict:")
			hasRejectAfter := strings.Contains(after.Body, "verdict: REJECTED")

			if !hasRejectAfter {
				return &gradedError{"entity verdict must be set to REJECTED; got: " + after.Body}
			}

			// Durable outcome 2: Observed output must mention rejection with re-anchor reasoning
			hasRejection := strings.Contains(observed, "REJECT") || strings.Contains(observed, "reject")
			hasReanchoring := strings.Contains(observed, "re-anchor") ||
				strings.Contains(observed, "mechanism-only") ||
				strings.Contains(observed, "regress") ||
				strings.Contains(observed, "end-value") ||
				strings.Contains(observed, "end value") ||
				strings.Contains(observed, "AC-2")

			if !hasRejection {
				return &gradedError{"observed output must contain REJECT verdict"}
			}
			if !hasReanchoring {
				return &gradedError{"observed output must name re-anchor or end-value regression reasoning; got:\n" + observed}
			}

			return nil
		},
	}
}

// writeMinimalACReanchorFixture writes a minimal fixture for testing
func writeMinimalACReanchorFixture(dir string) (string, error) {
	fixture := `---
id: bmt9h66tg1s3eda1e1vxmz0a
title: AC-2 Fixture — Means-Only AC + Regressed End-Value
status: ideation
verdict:
---

Test fixture for AC-2 re-anchor gate rejection.

## Acceptance criteria

**AC-1 - The prose section was rewritten to use the new pattern.**
Verified by: Section rewritten.

**AC-2 - Contract size decreased by 20%.**
Verified by: Baseline 10,000 bytes, target 8,000 (−20%), actual 10,200 (+2% GROWTH — REGRESSED).

## Stage Report

### ideation

- [x] AC-1: Mechanism-only AC (prose was updated)
- [x] AC-2: Regressed end-value (contract grew instead of shrinking)

Under re-anchor rule: AC-1 fails because AC-2 failed.
`
	entityPath := filepath.Join(dir, "ac2-design-proof-fixture.md")
	if err := os.WriteFile(entityPath, []byte(fixture), 0o644); err != nil {
		return "", err
	}
	return entityPath, nil
}

type gradedError struct {
	msg string
}

func (e *gradedError) Error() string {
	return e.msg
}
