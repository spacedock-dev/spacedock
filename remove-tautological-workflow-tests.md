---
id: 0mdhjk9jdf10h1qnhvx5tvn8
title: Remove the tautological workflow-file tests
status: backlog
source: "Captain directive, CL 2026-08-17, after measuring the edge-advance case: \"the goal of the task is removal of those tautological tests\". Evidence from this session: 1307 lines guarding a 173-line file were green across three consecutive releases while the mechanism they covered failed silently, and validation later found two material defects in that same area."
started:
completed:
verdict:
score:
worktree:
issue:
---

Delete the tests that assert a workflow file contains the text someone wrote into it. The goal is removal, not replacement.

## Problem

A family of tests in this repository reads a `.github/workflows/*.yml` file and string-matches its contents. The shape is:

```go
strings.Contains(strings.ReplaceAll(ifCond, " ", ""),
    "steps.decision.outputs.advance=='true'")
```

That asserts a string is present in a file the author just edited. It cannot observe a release, and it cannot fail for a real reason.

The cost is measured, not asserted. `edge_reconcile_test.go` (648 lines), `edge_advance_wiring_test.go` (341), and `edge_advance_noregress_test.go` (318) guarded `edge_advance_decision.go` (173 lines) at a ratio of 7.5 to 1. All three were green on 2026-08-15, 08-16 and 08-17. Across those same three releases the `edge-advance` job skipped silently inside green runs, `next` froze at `v0.27.0-pre4` and fell 99 commits behind, and the captain ran a four-week-old first-officer contract against a current binary. The tests could not have caught it: the workflow step was present and merely evaluated false, and reading the YAML never evaluates anything.

Task `2d` already deletes those three as collateral, because the mechanism they covered is gone. At least eight more files follow the same shape and roughly 2000 lines remain, the largest being `internal/release/journey_workflow_test.go` at 827 lines. Nothing removes them.

The class is self-perpetuating: this kind of test is cheap to write and never goes red, so it accumulates until the ratio looks like coverage.

## Proposed approach

{Ideation fills this in. The captain's stated goal is REMOVAL. The task should not become a rewrite programme — a deleted tautology needs no replacement, because it was never proving anything. Where a file turns out to hold a legitimate check mixed in with tautological ones, keep the legitimate check and delete the rest rather than preserving the file wholesale.}

## Out of scope

The three edge-advance files, which task `2d` already removes. Do not touch that worktree or those files.

Adding new tests of any kind. If a real gap is found behind a deleted tautology, record it and let the captain decide; do not fill it in this task.

## Not every file in the sweep is the same — this is the trap

A blanket deletion would be wrong, and the FO's initial survey found likely survivors. `node24_actions_guard_test.go` and `claude_version_float_guard_test.go` appear to check a config value against an independent rule — that an action version is pinned rather than floating. That is a real value which can diverge on its own, and the workflow's own proof policy admits it: "a static check counts only when it tests a real value against an independent source that can diverge from it, not as a spelling check over a file the model reads."

The discriminating question for each file is therefore: **can this fail for a reason other than someone editing the file it reads?** If the expected value is just the text the implementer wrote into the file under test, it is a tautology and goes. If it compares the file against an independent source of truth — a pinned version policy, a published manifest, a released artifact — it stays.

Candidate set (line counts as of 2026-08-17, all under `internal/`):
`release/journey_workflow_test.go` 827, `release/runtime_live_evidence_workflow_test.go` 298, `release/claude_candidate_binary_workflow_test.go` 210, `contractlint/workflow_trunk_test.go` 189, `contractlint/fo_write_core_mutation_gate_test.go` 171, `release/manifest_tag_gate_workflow_test.go` 163, `release/node24_actions_guard_test.go` 152, `release/e2egate_workflow_test.go` 151, `release/claude_version_float_guard_test.go` 144.

## Expected surface and tolerance

Estimate net LOC change: −1500 across roughly 9 files, of which the great majority is deletion. Declare insertions and deletions separately from the net figure at the ideation gate; do not declare a gross tolerance. Semantics changed: none — no production code, no command grammar, no stored format, no release behavior.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - No test in the repository asserts only that a workflow file contains a string the author put there.**
This is the measuring AC: the count of surviving assertions whose expected value is text read from the file under test must be ZERO, counted over the candidate set. The count moves the wrong way the moment such a test is re-added. Verified by: the audit's own per-file record, plus a reviewer reproducing the discriminating question on each surviving file.

**AC-2 - Every surviving check can fail for a reason other than an edit to the file it reads.**
Verified by: mutating the independent source rather than the file under test — for a version-pinning guard, move the pin policy and watch the test go red while the workflow file is untouched. Fails if a check can only be reddened by editing its own input, which is the definition this task exists to remove.

**AC-3 - No behavior lost with the deletions.**
Verified by: `go test ./...` green, and a named record of anything deleted that was guarding a live behavior, with the captain's decision on each. Fails if a deleted test was the only thing standing between a real regression and main — the same failure mode `2d`'s validator found when the reconcile deletion silently dropped a never-force-push guard.

## Test plan

{Ideation fills this in. Note the shape of this task: its deliverable is deletion, so its proof is the surviving set and the green suite, not new coverage. Resist writing a lint that detects tautological tests — that is a new standing enforcement mechanism, it is the last resort under the workflow's own rules, and it would need explicit captain approval as its own task.}
