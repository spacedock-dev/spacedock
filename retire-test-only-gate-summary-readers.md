---
id: 77k7m0dmwm10mz6zrdq86tv3
title: Retire or re-justify the test-only gate summary reader trio
status: backlog
source: "Gate finding from remove-gate-validate-subcommand ideation (2026-08-15): SummaryFile/SummaryFileAt/SummaryFileDiagnosticsAt dead-end in tests after the subcommand removal; captain accepted the follow-up at the gate"
started:
completed:
verdict:
score:
worktree:
issue:
---

After remove-gate-validate-subcommand lands, the `SummaryFile` / `SummaryFileAt` / `SummaryFileDiagnosticsAt` trio has no production caller, and the warnings slice `SummaryFileDiagnosticsAt` returns is observable by no caller. Retire the trio, or record the one concrete justification that keeps it. Removing it requires re-pointing two internal/gates tests.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

`ReadDiagnostics` and its `status --validate` consumer. The read tolerance.

## Expected surface and tolerance

Estimate net LOC change: -NN, across ~3 files.

## Acceptance criteria

**AC-1 - No production code path retains a reader whose output no caller can observe.**
Verified by: grep shows no non-test caller of the trio, or the entity records the justification and its consumer.

**AC-2 - The change removes more lines than it adds.**
Verified by: cumulative line delta against origin/main is negative.

**AC-3 - The suite stays green.**
Verified by: `go test ./...` and `go test ./... -race` pass.

## Test plan

Deletion plus re-pointed tests, existing suite as the net.
