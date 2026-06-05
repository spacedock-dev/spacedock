package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// Auto-continue fixtures and assertion for AC-5: a dev-shaped workflow parked at
// an implementation-ready state, exercising whether the FO continues the
// lifecycle (advance to validation + dispatch a fresh validator / present the
// gate) instead of stopping after the implementation report is filed. Like the
// other shared fixtures these live under the DEFAULT build tag so the offline
// negative case (auto_continue_negative_test.go) grades the assertion with no
// model spend, while the //go:build live half (auto_continue_live_test.go) drives
// the same fixture against a real agent.

// validationStatusOrBeyond matches an entity whose status advanced to validation
// (or the terminal done) — i.e. the FO did NOT leave it parked at implementation.
var validationStatusOrBeyond = regexp.MustCompile(`(?im)^status:\s*(validation|done)\s*$`)

// implementationStatusAC5 matches an entity still parked at implementation — the
// failure mode this regression guards against.
var implementationStatusAC5 = regexp.MustCompile(`(?im)^status:\s*implementation\s*$`)

func writeAutoContinueWorkflowNoGit(dir string) (string, error) {
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(autoContinueReadme()), 0o644); err != nil {
		return "", err
	}
	entityPath := filepath.Join(dir, "auto-continue-task.md")
	if err := os.WriteFile(entityPath, []byte(autoContinueEntity()), 0o644); err != nil {
		return "", err
	}
	return entityPath, nil
}

func writeAutoContinueWorkflow(t *testing.T, root string) string {
	t.Helper()
	entityPath, err := writeAutoContinueWorkflowNoGit(root)
	if err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)
	return entityPath
}

// autoContinueReadme is the dev-shaped fixture workflow from AC-4: a non-split-
// root workflow backlog → implementation → validation(worktree, fresh, gate) →
// done. validation is `fresh: true` so the FO must dispatch a FRESH validator,
// and `gate: true` so the FO presents the validation gate after the validator
// reports. The implementation stage carries a concrete validator instruction so a
// dispatched validator has real work that leaves a durable validation report.
func autoContinueReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"      worktree: true\n" +
		"    - name: validation\n" +
		"      worktree: true\n" +
		"      fresh: true\n" +
		"      feedback-to: implementation\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Auto-Continue Fixture\n\n" +
		"### backlog\n\nSeed the task.\n\n- **Outputs:** A seed task.\n\n" +
		"### implementation\n\nProduce the deliverable.\n\n- **Outputs:** The deliverable plus an implementation stage report.\n\n" +
		"### validation\n\n" +
		"Verify the implementation against the acceptance criteria. Append a `## Stage Report: validation` " +
		"section to the entity with one `- DONE:` item and a PASSED or REJECTED recommendation.\n\n" +
		"- **Outputs:** A PASSED or REJECTED validation stage report.\n\n" +
		"### done\n\nTerminal state.\n"
}

// autoContinueEntity is parked at status: implementation with a filed
// implementation stage report — the implementation-ready state from AC-4. A
// correct FO verifies this report and immediately advances to validation +
// dispatches a fresh validator (then presents the validation gate); a broken FO
// stops here with the report filed and the status unchanged.
func autoContinueEntity() string {
	return "---\n" +
		"id: auto-continue-task\n" +
		"title: Auto Continue Task\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Auto Continue Task\n\n" +
		"The implementation is complete and its stage report is filed below. The next lifecycle " +
		"step is independent validation.\n\n" +
		"## Acceptance criteria\n\n" +
		"**AC-1** The deliverable exists and is committed.\n" +
		"Verified by: the implementation stage report below plus a validation pass.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Produce the deliverable\n" +
		"  The deliverable is committed and ready for independent verification.\n" +
		"\n### Summary\n\n" +
		"Implementation is complete. The first officer must advance to validation and dispatch a fresh validator.\n"
}

// autoContinuePrompt is the NEUTRAL runbook from AC-4 — `Use $spacedock:first-
// officer` with no "drive to done" coaching. It points the FO at the workflow and
// the one parked entity and asks it to proceed normally. It deliberately does NOT
// tell the FO to advance, dispatch, or validate — whether it does so is exactly
// the behavior under test. Run non-interactively (`claude -p`), the FO enters
// single-entity mode and drives the parked implementation forward on its own; a
// broken FO stops after the implementation report instead.
func autoContinuePrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"Process the entity `auto-continue-task`. Its implementation worker has just completed and filed its stage report.",
		"Proceed with the workflow as the first-officer contract directs, then give your final response.",
	)
}

// assertAutoContinue is host-neutral: before/after entity-state strings plus the
// FO's observed output. It grades the DURABLE outcome, not transcript phrasing.
// The lifecycle continued when the entity is no longer parked at implementation
// (status advanced to validation or done) AND a validation stage report appears
// in the entity body — the durable footprint of a fresh validator the FO
// dispatched. A run that narrates "advancing to validation" in the transcript but
// leaves the durable state at status: implementation with no validation report
// fails on the state checks, not on transcript shape.
func assertAutoContinue(before, after, observed string) error {
	if implementationStatusAC5.MatchString(after) {
		return fmt.Errorf("FO left the entity parked at status: implementation — it stopped instead of advancing")
	}
	if !validationStatusOrBeyond.MatchString(after) {
		return fmt.Errorf("FO did not advance the entity to status: validation (or beyond)")
	}
	if !regexpValidationReport.MatchString(after) {
		return fmt.Errorf("no `## Stage Report: validation` appeared — the FO did not dispatch/run a validator")
	}
	// The implementation report must still be present (the FO advanced, it did not
	// discard the prior stage's report) — guards against an after-state that simply
	// replaced the body rather than appending the validation report.
	if !regexpImplementationReport.MatchString(before) {
		return fmt.Errorf("fixture invariant broken: before-state lacks the implementation stage report")
	}
	if !regexpImplementationReport.MatchString(after) {
		return fmt.Errorf("the implementation stage report was lost from the entity after the run")
	}
	return nil
}

var (
	regexpValidationReport     = regexp.MustCompile(`(?m)^## Stage Report: validation\b`)
	regexpImplementationReport = regexp.MustCompile(`(?m)^## Stage Report: implementation\b`)
)
