---
id: hz2ankag6fk379ssabpv4ckc
title: Repair Codex rejection-round recording in the live rejection flow
status: backlog
source: "Captain directive 2026-08-16 after two same-day codex failures (runs 31915540750 and 31922268382, both FAIL /rejection-flow observed=[rejection-round-missing]) on a journey whose XFAIL c6a336a33 retired on one unbound pass; old owner continue-codex-rejection-after-first-validation is archived done and fixed a different mode"
started:
completed:
verdict:
score: "0.90"
worktree:
issue:
---

Codex intermittently completes the live rejection flow without recording the rejection round: the FO-side flow reaches the feedback stage but `gate record --round` never runs, so the journey assertion finds no round record (`rejection-round-missing`). Fail-pass-fail across the last three runs proves the c6a336a33 retirement premature for this journey.

Two-part deliverable:
1. Restore the codex XFAIL binding on the `rejection-flow` journey in internal/ensigncycle/shared_live_runner_test.go, bound to THIS entity's id (the TODO-owner lint requires an active owner). This makes the lane honest immediately.
2. Diagnose and repair the codex behavior: why the round recording step is skipped (contract prose the codex FO under-weights, ordering, or a missing completion condition), fix at the owning surface, prove with a TARGETED LOCAL live run of the rejection-flow journey (SPACEDOCK_LIVE_RUNTIME=codex, -run scoped to the one journey, repeated until the mechanism - not luck - explains the green), then retire the XFAIL again in the same branch if the fix holds.

Captain-directed process: work on top of the current stack (branch from stack #717's tip, spacedock-ensign/rule-superseded-verdict-vocabulary at 4edc82f07); iterate with the targeted local live run; PR onto the stack as the next layer for the full-matrix run.

## Problem

{Ideation fills this in: the precise mechanism of the skip, from the two runs' codex-exec streams.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

Pi's rejection flow (p1 owns it). The registry amendment policy itself.

## Expected surface and tolerance

Estimate net LOC change: small - the binding line, plus the behavioral fix at its owning surface (likely contract prose or the journey runner's completion conditions).

## Acceptance criteria

**AC-1 - The rejection-flow journey passes a targeted local codex live run repeatedly, with the fix - not stochasticity - explaining the green.**
Verified by: N consecutive targeted runs green (ideation sets N) plus the mechanism named from the exec streams; a run WITHOUT the fix still shows the failure mode or its precondition.

**AC-2 - The lane is honest at every commit: the XFAIL is bound to this active entity while the flake exists, and removed in the same branch only if the fix holds.**
Verified by: contractlint TODO-owner join green; the binding's presence or absence matches the fix state.

**AC-3 - The suite stays green and the full matrix runs on the stack PR.**
Verified by: offline suites plain and -race; the stack-tip full run after this layer lands.

## Test plan

Targeted local live runs for the repair loop (cheap, one journey); the stack PR's full matrix as the final proof; no new standing harness.
