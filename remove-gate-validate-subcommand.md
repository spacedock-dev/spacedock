---
id: 0tmv5bry1wbkww2758y88pay
title: Remove the gate validate subcommand
status: backlog
source: "Captain directive, 2026-08-14: value review found near-zero value; every fault it reports also surfaces elsewhere"
started:
completed:
verdict:
score:
worktree:
issue:
sprint: durable-decisions
gates:
    version: 1
    records:
        - id: gate:0tmv5bry1wbkww2758y88pay:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0tmv5bry1wbkww2758y88pay-backlog-1
              briefing:
                id: briefing:0tmv5bry1wbkww2758y88pay:backlog:attempt-1:revision-1
                digest: sha256:30b8fe0e5bc71046ebb679fe9cc6bab32c2bae8fc294c99e68b2307502b675a3
                request-digest: sha256:c67727e7678ea741b41d589a7b94161862beb12ae8bfa7c2554546a0c4c5157b
                room-ref: ./remove-gate-validate-subcommand/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0tmv5bry1wbkww2758y88pay:backlog:1
                briefing: briefing:0tmv5bry1wbkww2758y88pay:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:25.738011Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: pending
---

Remove the `gate validate` CLI subcommand. It has zero live uses across 424 recorded attempts. Fatal faults surface identically on every gate command through the shared reader. The warning class prints at state publish. The read-only digest sweep and the round check have no recorded consumer.

Scope: the two CLI branches in internal/cli/cli.go (plain and --round), the usage text, tests of the subcommand, and skill prose that names `gate validate`. Keep `ValidateRoundFile` (live consumer: `gate record --round`) and `SummaryFileDiagnosticsAt` (live consumer: `SummaryFileAt`).

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

The rest of the gate decision surface. The shared readers and `ValidateRoundFile`.

## Expected surface and tolerance

Estimate net LOC change: -NNN, across ~4 files. Observable semantics change: `gate validate` becomes an unknown subcommand.

## Acceptance criteria

**AC-1 - The command surface no longer carries the subcommand.**
Verified by: a CLI test asserts that `spacedock gate validate` exits with the usage error.

**AC-2 - The change removes more lines than it adds.**
Verified by: cumulative line delta of the diff against origin/main is negative.

**AC-3 - The suite stays green.**
Verified by: `go test ./...` and `go test ./... -race` pass.

## Test plan

Deletion plus the existing suite, plus one usage-error assertion.
