---
title: Repair entered-stage dispatch evidence and post-gate terminalization
status: ideation
source: "PR #585 exact-head run 30706782428: Codex job 91387287118 and Sonnet job 91387287120"
started: 2026-08-09T18:34:29Z
completed:
verdict:
score: 0.9
worktree:
issue:
milestone: 0.27.0
id: 9adv48yhye5s2vkhwd7ge52d
sprint: test-behavior-completeness
gates:
    version: 1
    records:
        - id: gate:9adv48yhye5s2vkhwd7ge52d:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9adv48yhye5s2vkhwd7ge52d-backlog-1
              briefing:
                id: briefing:9adv48yhye5s2vkhwd7ge52d:backlog:attempt-1:revision-1
                digest: sha256:2cb413910c222d4d8b9a3b47fa2f43705c6a5198c796988425da0956e65f8c5e
                request-digest: sha256:9bcea1bdbd613334d989eaf10920c3e50f40c4e6f6386993352a0bd258b4989e
                room-ref: ./repair-entered-stage-dispatch-and-post-gate-terminalization/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9adv48yhye5s2vkhwd7ge52d:backlog:1
                briefing: briefing:9adv48yhye5s2vkhwd7ge52d:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:17.217951Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; shaping assigns initial-stage dispatch to task 6x and keeps post-gate behavior here.
              application:
                target-stage: ideation
                state: consumed
sprint-readiness: ready
group: common-evidence
---

## Problem

PR #585's runtime-neutral live oracle found two defects in one post-gate journey.

First, a gate consume can advance an entity to `implementation`, then the dispatch evidence can split across two commits. The first commit names the entered stage before `started` exists. The second commit adds `started`. The durable oracle requires one exact path-scoped dispatch commit with both values.

Second, a consumed nonterminal approval is treated as an error during ordinary terminal status writes. A worker can complete its report, but the entity cannot reach `done` without `--force`. That breaks the non-force safety rule.

The initial-stage successor defect is not part of this task. Task `6x50qafc8566zc6p1qpb6y30` owns that behavior and its smoke tests. This task must not duplicate that repair or its live journey.

Staff finding M4 adds an ordering constraint. The two tasks share
`internal/ensigncycle/shared_live_runner_test.go` and
`skills/first-officer/references/fo-dispatch-core.md`. Task
`6x50qafc8566zc6p1qpb6y30` owns every `smallest-sufficient-mechanism` binding.
This task owns no smallest-sufficient-mechanism binding.

## Value slice

This task owns one end-to-end value slice: consume an approved nonterminal gate, dispatch the successor with durable evidence, record the worker result, and complete the entity with ordinary non-forced terminal fields.

The `keep-moving-posture` journey is the live value check. Its other ready entities remain controls. The `smallest-sufficient-mechanism` journey and its initial-stage rows remain with task `6x50qafc8566zc6p1qpb6y30`.

## Proposed approach

1. Run the affected `keep-moving-posture` cells with the strict XFAIL mechanism from `ts7gq0mr9s3chx2w4wppd1kt` before any product edit. Bind one stable semantic failure code to each executable target. Keep a TODO only when that target cannot run. Do not add or change a `smallest-sufficient-mechanism` binding.

2. Keep gate consume and successor dispatch as one documented ceremony boundary. The consume command writes the successor status and `consumed` gate state. One `dispatch build --stamp` call then writes `started` and `worktree`, commits the path-scoped entity, and builds the envelope. Do not add a manual status write, a separate state commit, or a plain dispatch build between these commands.

3. Add a smoke assertion that reads the entity blob from the exact `dispatch: <slug> entering <stage>` commit. The blob must contain the entered stage and a nonempty `started` value. A later commit that adds `started` must fail the assertion.

4. Narrow `pendingTerminalApproval` so a consumed nonterminal application is ordinary stage history. It must allow the atomic non-forced terminal fields after the worker report. It must still fail closed for pending terminal-target approval, unreadable authority, stale authority, and a binding to a terminal target.

5. Run the live value cells again. XPASS must remove the repaired XFAIL binding. The durable Git-history oracle remains unchanged.

The simplest alternatives are insufficient. Status prose cannot prove the exact commit blob. A second stamp commit cannot meet the single-boundary requirement. `--force` hides an authority error and violates the value slice. A broad bypass for all consumed records would weaken terminal-target safety.

## Out of scope

- Initial-stage dispatch when `current=ready,next=done`. Task `6x50qafc8566zc6p1qpb6y30` owns this item.
- All `smallest-sufficient-mechanism` bindings and their live proof. Task `6x50qafc8566zc6p1qpb6y30` owns the Claude Sonnet, Codex, and Pi bindings.
- Changes to the durable Git-history oracle, its path scope, or its failure rules.
- New command flags, new frontmatter fields, or a new gate authority format.
- Agent instructions that use `--force` to finish the journey.
- The generic strict-XFAIL runner work owned by `ts7gq0mr9s3chx2w4wppd1kt`.

## Acceptance criteria

**AC-1 (VALUE) — A consumed nonterminal post-gate journey reaches normal completion.**

The `keep-moving-posture` fixture moves `approved-gate` from `review` to `implementation`, records the worker report, and reaches `done` with non-forced atomic terminal fields. The path-scoped history has the consume, dispatch, report, and terminal evidence. The journey does not use `--force`.

Verified by: the strict XFAIL-enabled live cells for each executable host, plus the durable journey grader. Falsifiers: a refusal at `done`, a `--force` invocation, a missing report, or a missing terminal record.

**AC-2 — The dispatch boundary contains complete entered-stage evidence.**

After nonterminal gate consume, the exact path-scoped commit named `dispatch: <slug> entering implementation` contains `status: implementation` and a nonempty `started` field in the entity blob. The envelope exists after that commit. A second commit is not needed.

Verified by: a focused CLI or dispatch smoke test that reads the exact commit blob. Falsifiers: an absent `started` field, a split commit, a missing envelope, or an extra dispatch stamp commit.

**AC-3 — Terminal authority remains fail closed.**

A consumed nonterminal record allows ordinary atomic terminal fields. Pending terminal-target approval, unreadable authority, digest-stale authority, and a terminal-target binding still refuse the non-forced status write and direct the caller to the merge guard.

Verified by: the focused CLI regression and the existing terminal refusal matrix. Falsifiers: `--force` being required for the consumed nonterminal case, or any protected terminal-target case succeeding without its merge guard.

**AC-4 — The 9a candidate uses the landed 6x candidate.**

Task 6x lands before 9a implementation. The 9a branch rebases onto the exact 6x landing commit before it changes a shared runner or dispatch file. The exact candidate reruns the 9a deterministic tests and the `keep-moving-posture` cells.

Verified by: the 6x landing SHA, the serial rebase record, and test records from the rebased candidate. Falsifiers: a stale-base candidate, a duplicate smallest-sufficient binding, or a shared-file change without an exact-candidate rerun.

## Test plan

Run the smallest proof first:

1. Run the existing `internal/dispatch` stamp tests for commit, worktree, mismatch, and idempotence.
2. Run the existing terminal refusal and merge-guard tests.
3. Add and run the exact dispatch-blob smoke test.
4. Add and run the `review(gate) -> implementation(consumed) -> done` CLI regression.
5. Run the recorded gate lifecycle and durable journey tests.
6. After 6x lands, rebase 9a onto the exact 6x landing commit.
7. Run strict XFAIL for the three `keep-moving-posture` target cells before the repair, then rerun them after the repair.
8. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

Each new test must name its falsifier. The consumed-nonterminal regression must prove that the final status write does not call `--force`. The dispatch smoke must inspect the exact commit, not only the final worktree.

## Expected surface and semantic changes

The expected implementation surface is below. Estimates exclude the generic XFAIL runner and task `6x`.

| File | Gross additions | Gross deletions | Net | Purpose |
| --- | ---: | ---: | ---: | --- |
| `internal/status/merge.go` | 12 | 4 | +8 | Classify consumed nonterminal authority as ordinary history while preserving fail-closed terminal-target checks. |
| `internal/cli/terminal_consume_test.go` | 120 | 0 | +120 | Add the focused consumed-nonterminal completion and refusal matrix. |
| `internal/cli/gate_ceremony_count_test.go` | 45 | 0 | +45 | Inspect the exact dispatch commit entity blob. |
| `internal/ensigncycle/shared_live_runner_test.go` | 28 | 6 | +22 | After 6x lands, bind only the three executable `keep-moving-posture` cells to strict XFAIL, then remove those bindings on XPASS. |
| `skills/fo-gate-lifecycle/SKILL.md` | 12 | 4 | +8 | State the one consume-then-`dispatch build --stamp` boundary. |
| `skills/first-officer/references/fo-dispatch-core.md` | 10 | 2 | +8 | Require complete entered-stage evidence in the dispatch commit. |
| `docs/site/concepts/gates-and-decisions.md` | 6 | 1 | +5 | Document ordinary terminalization after consumed nonterminal approval. |
| `docs/site/reference/command-reference.md` | 6 | 1 | +5 | Document the non-forced terminal status rule and dispatch evidence. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | 8 | 1 | +7 | State the authority distinction in the contract. |
| **Total** | **247** | **19** | **+228** | **Tolerance: ±25% for gross and net estimates.** |

Semantic changes are narrow. Command grammar does not change. Stored frontmatter formats do not change. Authority behavior changes only for a consumed nonterminal application during ordinary terminalization. Pending, unreadable, stale, and terminal-target authority remains protected. Runtime behavior requires `dispatch build --stamp` to provide the single complete dispatch evidence boundary.

## Documentation diff proposal

Add this rule to the gates concept and command reference:

> After a nonterminal approval is consumed and its worker report is present, ordinary atomic terminal fields can complete the entity without `--force`. Pending, unreadable, stale, or terminal-target gate authority still fails closed and uses the merge guard.

Add this rule to the dispatch contract:

> For a consumed nonterminal entry, the exact `dispatch: <slug> entering <stage>` commit must contain the entered stage and nonempty `started` value. Do not split those values across commits.

## Dependencies and XFAIL-first order

Task `6x50qafc8566zc6p1qpb6y30` is a landing dependency. It must land its initial-stage repair and all three `smallest-sufficient-mechanism` bindings for Claude Sonnet, Codex, and Pi before 9a implementation starts. The 9a task owns no smallest-sufficient-mechanism binding.

Staff finding M4 requires serial edits to shared files. After 6x lands, rebase 9a onto the exact 6x landing commit. Run the exact 9a candidate tests on that rebased commit. If either shared file changes again, repeat the rebase and exact-candidate reruns.

The strict-XFAIL runner from `ts7gq0mr9s3chx2w4wppd1kt` must be available before either live binding changes. Neither dependency changes this task's value slice.

Before a product edit, run each executable `keep-moving-posture` target and record one stable failure code. After the smoke and CLI tests pass, implement the narrow status change and contract updates. Then rerun the live cells. XPASS removes the 9a binding. A target that cannot run keeps a TODO with its reason.

## Spike result

The riskiest available path was exercised first. The existing stamp tests pass for commit creation, worktree creation, stage mismatch, and idempotent rerun. The existing terminal tests pass for pending, stale, unreadable, and terminal-target authority. Recorded gate lifecycle and durable journey tests also pass.

The spike found the remaining proof gaps. No current test reads the entity blob from the exact dispatch commit. No current CLI regression covers `review(gate) -> implementation(consumed) -> done`. Source tracing shows that `pendingTerminalApproval` currently rejects consumed nonterminal authority. The plan adds proof before the narrow repair. It introduces no new state or authority primitive.

## Stage Report: ideation

- DONE: Remove initial-stage defect from this task and keep one post-gate value slice.
  Task `6x50qafc8566zc6p1qpb6y30` owns initial-stage dispatch and every `smallest-sufficient-mechanism` binding. This entity now owns only `keep-moving-posture` after gate consume.
- DONE: Exercise gate-consume dispatch evidence and safe terminalization before selecting the repair.
  Existing stamp, terminal-authority, recorded-lifecycle, and durable-journey tests passed. Source tracing found the consumed-nonterminal classifier and two missing proof tests.
- DONE: Give gross and net line estimates with XFAIL-first dependencies.
  The plan estimates +247/-19, net +228, with ±25% tolerance. Strict XFAIL runner `ts7gq0mr9s3chx2w4wppd1kt` is the first dependency.
- DONE: Fold staff finding M4 into this ideation.
  Task 6x lands first. The 9a task then rebases onto its exact candidate and reruns every shared-runner and dispatch-file proof from that candidate.

### Summary

Ideation is complete. The plan keeps one post-gate value slice, requires exact dispatch-commit evidence, and narrows terminal authority only for consumed nonterminal history.

Task `6x50qafc8566zc6p1qpb6y30` owns all smallest-sufficient-mechanism bindings and must land before 9a implementation. Staff finding M4 now requires a serial rebase and exact-candidate reruns for the shared runner and dispatch files.
