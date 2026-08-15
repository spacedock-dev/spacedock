---
id: e627de4qbh7xeeff2fxjz9kz
title: Fix the fixture-repo object-read flake in ensigncycle durable tests
status: backlog
source: "CI run 31896322753/31896322757 on PR #698, 2026-08-15: TestDurableQuestionedRejectsTerminalHistory/status failed with git log --follow fatal: unable to read 8354e03b in the test's own temp repo; local -count=3 and CI rerun green"
started:
completed:
verdict:
score:
worktree:
issue:
---

A test's own throwaway git fixture lost an object mid-walk under CI load. The failure is not the code under test - it is a race in fixture construction or an fsync/odb hazard. Rerun-to-green established that run's verdict but left the race in place; the next occurrence reds an unrelated PR and costs the same diagnosis hour.

## Problem

{Ideation fills this in: reproduce under load or by instrumentation, name the race.}

## Proposed approach

{Ideation fills this in: candidate directions - serialize fixture git operations, fsync/pack after fixture setup, or prove the hazard external and bound it with a targeted retry at fixture level only.}

## Out of scope

Blanket test retries; weakening any assertion.

## Expected surface and tolerance

Estimate net LOC change: small, test files only.

## Acceptance criteria

**AC-1 - The failure signature has a named root cause with evidence, and the fix targets it.**
Verified by: the report cites the mechanism and a repro or instrumentation trace, not a guess.

**AC-2 - The durable-history tests survive a stress run.**
Verified by: -count and parallel stress of the affected tests green where the baseline shape flaked.

**AC-3 - The suite stays green.**
Verified by: go test ./internal/ensigncycle/ offline shape.

## Test plan

Instrument, reproduce, fix, stress; the CI run URL is the baseline evidence.
