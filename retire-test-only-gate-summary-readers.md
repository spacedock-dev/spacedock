---
id: 77k7m0dmwm10mz6zrdq86tv3
title: Retire or re-justify the test-only gate summary reader trio
status: ideation
source: "Gate finding from remove-gate-validate-subcommand ideation (2026-08-15): SummaryFile/SummaryFileAt/SummaryFileDiagnosticsAt dead-end in tests after the subcommand removal; captain accepted the follow-up at the gate"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:77k7m0dmwm10mz6zrdq86tv3:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:77k7m0dmwm10mz6zrdq86tv3-backlog-1
              briefing:
                id: briefing:77k7m0dmwm10mz6zrdq86tv3:backlog:attempt-1:revision-1
                digest: sha256:8bd171d5a202a453728de4bce4cf9760d3044a3a17e4cf28615810d5d311c20d
                request-digest: sha256:db59d64f51fc38a4a61f0950e227498d48e75206c0eedae708aa20dc832c0ece
                room-ref: ./retire-test-only-gate-summary-readers/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:77k7m0dmwm10mz6zrdq86tv3:backlog:1
                briefing: briefing:77k7m0dmwm10mz6zrdq86tv3:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T21:24:39.885375Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: dispatch all five onto the stack tip'
              application:
                target-stage: ideation
                state: consumed
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
