---
title: Retire the keep-moving permission narration oracle
status: implementation
source: "PR #512 Runtime Live E2E run 29382760645: Opus jobs 87249808752 and 87252929149 both advanced and dispatched the approved entity, then failed because kmPermissionRe matched a negated quotation in the final summary."
started: 2026-07-15T04:14:25Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-retire-keep-moving-permission-narration-oracle
issue:
milestone: 0.25.0
id: bjdm4tdnk93813p9nj913j2y
---

The keep-moving live grader must prove that an approved entity advances and dispatches without letting free-form summary wording veto structured action evidence.

## Problem

`kmPermissionRe` scans the model's final prose for permission-question phrases. In two exact-head Opus runs, the FO performed the required advance and dispatch, then accurately summarized that it used `no "want me to advance?" pause`; the regex matched the quoted phrase and failed both otherwise-correct journeys. This is an evidence defect and a recurring unrelated release blocker, not an outcome defect.

## Proposed approach

Retire the final-message permission regex and its synthetic negative. Keep the existing structured per-entity advance and dispatch checks as the authoritative proof of the user-visible no-false-stop value. Do not add quote stripping, negation grammar, another narration parser, or a model-facing wording instruction.

## Out of scope

- Changes to the keep-moving contract or its structured advance/dispatch expectations.
- Changes to 1k's shallow-boot deferred-loading evidence.
- Any new transcript dialect, shell parser, controller, lifecycle layer, or retry policy.

## Acceptance criteria

**AC-1 - Correct completed motion cannot fail on a negated quotation.**
Verified by: a committed replay fixture based on the two failing Opus final-message shapes passes when structured evidence shows `approved-gate` advanced and dispatched.

**AC-2 - A real false stop remains red.**
Verified by: existing negative cases with missing approved advance or missing approved dispatch still fail independently and name the missing action.

**AC-3 - The suite has less inferred evidence machinery.**
Verified by: `kmPermissionRe`, `askedPermission`, and their narration-only negative are deleted with no replacement parser; the task has a negative test-infrastructure line delta.

**AC-4 - Repository gates remain green.**
Verified by: focused keep-moving tests, `go test ./...`, `go test ./... -race`, exact-head Roborev after local green, and the next Opus Runtime Live E2E lane.

## Test plan

Use offline replay/negative tests for the grader change, then the required full and race suites. The next already-required PR live lane is confirmation; do not build a new live harness or run local Claude.

## Stage Report: implementation

- DONE: Delete kmPermissionRe, the askedPermission field and every assignment (claudeKeepMovingTrace, codexKeepMovingTrace, gradeKeepMoving), and the narration-only synthetic negative — with NO replacement parser; structured advance/dispatch checks remain the authoritative no-false-stop proof (AC-3).
  Commit 3b11e34d: regex, `askedPermission` field, its `gradeKeepMoving` check, both extractor assignments, `kmPermissionFinal`, and the three permission-question negatives removed; grep of the four identifiers now returns nothing; no quote-strip/negation parser added.
- DONE: Add a committed replay fixture from the two failing Opus final-message shapes (negated quotation): correct completed motion PASSES on structured evidence (AC-1); a real false stop with missing approved advance/dispatch stays RED and names the missing action (AC-2).
  `TestKeepMovingNegatedQuotationReplay`; both finals genuinely matched the retired regex (exercised: `want me to advance`, `should I proceed`) yet pass on structured advance+dispatch; missing-advance error names "advance", missing-dispatch names "dispatch".
- DONE: gofmt clean; go build/go vet; go test ./... and go test ./... -race green; report the negative test-infrastructure line delta (AC-3/value).
  `gofmt -l` empty; build/vet exit 0; full suite + `-race` both `ok` across all packages; grader-machinery file `shared_keep_moving_test.go` net −10 lines (+13/−23) — the inferred-evidence apparatus shrank; the tests file's +43 is the grounded replay fixture, not a parser.

### Summary

Retired the final-message permission-narration oracle: the regex, the `askedPermission` trace field with its grade check and both host-extractor assignments, and the narration-only permission-question negatives (plus `kmPermissionFinal`) are gone, with no replacement parser — the existing structured advance/dispatch checks are now the sole no-false-stop proof. Added `TestKeepMovingNegatedQuotationReplay` distilled from the two failing Opus finals (jobs 87249808752/87252929149), proving a correct completed motion carrying a negated permission quotation passes while a genuine missing-advance/missing-dispatch false stop stays red and names the missing action. gofmt/build/vet clean and the full + race suites are green; Roborev and the next Opus live lane are the downstream confirmation per the test plan.
