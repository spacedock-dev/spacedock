---
id: vtgw5b1ar42s6x807nfdzwb9
title: Three shared live-lane failures that are not model faults
status: backlog
source: "Captain, 2026-10-27: Haiku analysis after PR #853 (runs 37703335512, 40448057053, 40447244422) found four of nine Haiku journey failures are grader bugs or cross-model failures. Captain: file the targeted fixes and get them passing locally."
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:vtgw5b1ar42s6x807nfdzwb9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:vtgw5b1ar42s6x807nfdzwb9-backlog-1
              briefing:
                id: briefing:vtgw5b1ar42s6x807nfdzwb9:backlog:attempt-1:revision-1
                digest: sha256:481cb8d061369714414c8a8464c4d84fb6dd9483d2bbb40976d16f250c72edff
                room-ref: '@review/backlog/briefing-1'
---

Fix three live-lane failures that hit Sonnet and Haiku alike, so live results show real model behavior and not harness or contract faults.

## Problem

1. **Filing grader false negative.** `commandFilesViaNew` in `internal/ensigncycle/shared_filing_test.go` accepts only a space, tab, quote, newline, or end of line after the slug. A correct `${SPACEDOCK_BIN:-spacedock} new wire-the-thing; echo "exit=$?"` fails with `filing-command-not-observed`. Both Haiku runs did exactly this (10-07 and 10-27), and the entity landed with id 001. Sonnet passed only because its command ended with `--workflow-dir . ;`.
2. **keep-moving-posture: merge guard refuses the approved-gate entity.** After the captain-approved `approved-gate` entity finishes implementation and reaches the ungated terminal stage, `spacedock merge guard` reports it "ineligible" (no pending terminal approval), so it is never archived. The Sonnet and Haiku streams show the same refusal. This is the failure behind Sonnet's XFAIL (owner `060xp69y61yhrww23g3wvwqy`) and Haiku's blank-owner XFAIL. It is not yet known whether the defect is in the binary's merge guard (an entity that was gated earlier, then reached an ungated terminal) or in the FO contract's terminal path.
3. **Break-glass `selected-team` fails on every model, and the `team-lead` completion target is unreachable.** `TestLiveBreakGlassShimRecovery/selected-team` fails on Sonnet (PR #853 run 40447244422 and the `kent/compass-poc` run 39608837699) and on Haiku. Every worker-dispatching scenario, for both models, logs `No agent named 'team-lead' is reachable`. The completion target is pinned to `team-lead`, but the current Claude Code resolves only `main`. The suspected link between the two is unproven.

## Proposed approach

{Ideation fills this in. Diagnose items 2 and 3 to a root cause before choosing a fix.}

## Risk evidence

Backlog: the three failure signatures and their streams are in the artifacts of runs 37703335512, 40448057053 and 40447244422 (expiry about 90 days). Item 1 is confirmed by applying the matcher regex to the recorded commands. Items 2 and 3 have a known signature but no known cause.

## Out of scope

- Haiku-specific behavior failures: the rejection-flow worker ordering, the named bare `Agent()` call, and the literal FAILED worker report.
- Changing which model or cadence the lane runs, or the blank-owner Haiku bindings, except to remove a binding that a fix makes obsolete.
- The `parseLiveGap` wiring and the two docs that disagree on whether `pi-live` must pass.

## Expected surface and tolerance

Estimate net LOC change: {ideation}, across {ideation} files.

## Acceptance criteria

**AC-1 - A correct filing command with a trailing `;` or `&&` grades as filed.**
Verified by: a focused matcher test over the recorded Haiku command, red before the fix and green after, plus a local live `filing` journey run that passes.

**AC-2 - The keep-moving approved-gate entity reaches done and is archived through the merge guard.**
Verified by: a local live `keep-moving-posture` run that passes on Sonnet with no XFAIL needed. The Sonnet `060xp…` binding and the Haiku binding on this journey are removed or re-scoped, with the reason recorded.

**AC-3 - Break-glass `selected-team` passes, and worker completion messages reach the FO without a `team-lead`-unreachable error.**
Verified by: a local live `TestLiveBreakGlassShimRecovery` run that passes both subtests, and a stream that has no `No agent named 'team-lead' is reachable`.

## Test plan

Focused unit tests for each fix, then local live runs of the affected journeys (`SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -run ...` per `docs/runtime-live-ci.md`), on Sonnet at minimum. Then `go test ./...` and `go test ./... -race`.

### Feedback Cycles
