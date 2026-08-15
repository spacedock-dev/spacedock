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
