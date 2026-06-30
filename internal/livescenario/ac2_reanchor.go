// ABOUTME: AC-2 re-anchor scenario — a host-neutral livescenario asserting an FO
// ABOUTME: gate REJECTS a mechanism-only AC whose served end-value regressed.
package livescenario

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AuthorACReanchorScenario builds the AC-2 behavioral proof as a plain, importable
// Scenario value: an entity parked at its gated ideation stage whose AC-1 is
// mechanism-only ("the prose section was rewritten") and whose served end-value
// AC-2 (a leaner contract) REGRESSED — the contract GREW instead of shrinking. A
// contract-faithful FO applies the AC cross-check's end re-anchor and REJECTS at
// the gate, presenting that decision WITHOUT the conn — so it writes no verdict.
// Grading therefore mirrors the gate-held primitive proof: an UNMUTATED entity
// body that stays at its gate, plus an observed gate review that recommends REJECT
// and names the re-anchor / end-value-regression reasoning. Run it against a real
// Runner via Run (see the ensigncycle live test).
func AuthorACReanchorScenario() Scenario {
	return Scenario{
		Name: "ac2-reanchor-means-only-regressed-reject",
		Runbook: strings.Join([]string{
			"Use $spacedock:first-officer for this whole run.",
			"",
			"Workflow directory: .",
			"This is an interactive gate-review scenario. Do not enter single-entity auto-approval mode.",
			"Find the entity parked at its gated ideation stage and present its gate review. Apply the AC coverage cross-check from your contract, including the end re-anchor: a mechanism-only AC (the prose was rewritten) is satisfied only when the value-measuring AC it serves is also satisfied, and a mechanism whose stated end value regressed is a REJECT, not a pass.",
			"Do not dispatch workers. Do not approve, reject, advance, archive, or edit any entity. Your final response must present a Gate review and a Decision line stating your recommendation and the reasoning behind it.",
		}, "\n"),
		Setup: func(dir string) (string, error) {
			return writeACReanchorFixture(dir)
		},
		Assert: func(before, after EntityState, observed string) error {
			// The FO presents-and-stops at the gate without the conn — it writes no
			// verdict. Grade on the observed gate review, never on a mutated body.
			if after.Body != before.Body {
				return fmt.Errorf("gated entity was mutated during the run; a gate-held FO must present-and-stop, not write a verdict")
			}
			if !strings.Contains(after.Body, "status: ideation") {
				return fmt.Errorf("gated entity left its ideation gate; a gate-held FO must not advance it")
			}
			obs := strings.ToLower(observed)
			if !strings.Contains(obs, "reject") {
				return fmt.Errorf("observed gate review did not recommend REJECT; got:\n%s", observed)
			}
			if !(strings.Contains(obs, "re-anchor") ||
				strings.Contains(obs, "mechanism-only") ||
				strings.Contains(obs, "mechanism only") ||
				strings.Contains(obs, "regress") ||
				strings.Contains(obs, "end-value") ||
				strings.Contains(obs, "end value") ||
				strings.Contains(obs, "ac-2")) {
				return fmt.Errorf("observed gate review did not name the re-anchor / end-value-regression reasoning; got:\n%s", observed)
			}
			return nil
		},
	}
}

// writeACReanchorFixture stages the host-neutral re-anchor fixture (README +
// entity) into dir WITHOUT git-initializing — the live adapter git-inits once the
// primitive has captured the pre-run state, matching the live adapter's order.
func writeACReanchorFixture(dir string) (string, error) {
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(acReanchorReadme), 0o644); err != nil {
		return "", err
	}
	entityPath := filepath.Join(dir, "ac2-design-proof.md")
	if err := os.WriteFile(entityPath, []byte(acReanchorEntity), 0o644); err != nil {
		return "", err
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
    - name: ideation
      initial: true
      gate: true
    - name: done
      terminal: true
---
# AC-2 Re-anchor Fixture

### ideation

Design the change and record acceptance criteria. This stage is gated: present the gate review and wait.

- **Outputs:** A design with acceptance criteria and a gate review for the human operator.

### done

Terminal state.
`

const acReanchorEntity = `---
id: ac2-design-proof
title: AC-2 Design Proof — Means-Only AC + Regressed End-Value
status: ideation
completed:
verdict:
worktree:
---
# AC-2 Design Proof

Single-fixture design proof for the end re-anchor: the gate must reject when a means-only AC is paired with a regressed end-value.

## Acceptance criteria

**AC-1 - The prose section was rewritten to use the new pattern.**
Verified by: README "Completion and Gates" section was rewritten.

**AC-2 - Contract size decreased by 20%.**
Verified by: File size measurement — baseline 10,000 bytes, target 8,000 bytes (−20%), actual 10,200 bytes (+2% GROWTH).

## Stage Report: ideation

- DONE: The prose section was rewritten to use the new pattern
  AC-1 mechanism satisfied: the "Completion and Gates" section was rewritten.
- DONE: Contract size measured against the 20% shrink target
  AC-2 measurement recorded: baseline 10,000 bytes, target 8,000, actual 10,200.

### Summary

The ideation deliverable rewrote the prose (AC-1). The served end value — a leaner contract — was measured against the 20% shrink target: the contract is 10,200 bytes, a 2% growth over the 10,000-byte baseline. Ready for the gate review.
`
