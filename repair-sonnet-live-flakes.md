---
title: Repair Sonnet keep-moving and selected-bare live flakes
status: validation
score: "0.95"
source: "Exact-main Runtime Live E2E run 31894934349 attempts 1 and 2"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: 060xp69y61yhrww23g3wvwqy
gates:
    version: 1
    records:
        - id: gate:060xp69y61yhrww23g3wvwqy:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:060xp69y61yhrww23g3wvwqy-backlog-1
              briefing:
                id: briefing:060xp69y61yhrww23g3wvwqy:backlog:attempt-1:revision-1
                digest: sha256:09c57e5fa0e468308825f688c19ec7539dc7e75e89e655d69e0a5f1416aa6012
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:060xp69y61yhrww23g3wvwqy:backlog:1
                briefing: briefing:060xp69y61yhrww23g3wvwqy:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-03T18:26:25.611996Z"
                decision: approve
                reason: 'Captain directed in chat on 2026-11-03: ''for those 2, do we know how to fix it? if not dispatch local probe on top of the worktree.'''
              application:
                target-stage: ideation
                state: consumed
        - id: gate:060xp69y61yhrww23g3wvwqy:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:060xp69y61yhrww23g3wvwqy-ideation-1
              briefing:
                id: briefing:060xp69y61yhrww23g3wvwqy:ideation:attempt-1:revision-1
                digest: sha256:8d5375fc28e2403f34f93c4a2f7bd822c14c174b27b3c834b95cb5087aaf7c6e
                room-ref: '@review/ideation/briefing-1'
              withdrawal:
                by: agent:first-officer
                at: "2026-11-03T21:36:51.86205Z"
                reason: 'New probe evidence after the briefing: 4 unchanged-byte runs killed by the quiet budget, and a second keep-moving cause (host stopped the background worker 600 s after the FO''s last turn).'
            - id: gate-attempt:060xp69y61yhrww23g3wvwqy-ideation-2
              briefing:
                id: briefing:060xp69y61yhrww23g3wvwqy:ideation:attempt-2:revision-1
                digest: sha256:fa1131626a16070e5ff32a62db21822f57cffdf92afcc1c842f8c749212dff77
                room-ref: '@review/ideation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:060xp69y61yhrww23g3wvwqy:ideation:2
                briefing: briefing:060xp69y61yhrww23g3wvwqy:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-11-03T23:57:36.952597Z"
                decision: approve
                reason: 'Captain approves the step-5 fix with the smaller test: hardcode the advance sequence and exercise only the binary path; the test does not read fo-dispatch-core.md. The 060xp keep-moving binding stays until retirement evidence exists. The background-worker stop and the quiet-budget kill are routed outside this task.'
              application:
                target-stage: implementation
                state: consumed
started: 2026-11-03T18:28:41Z
worktree: .worktrees/spacedock-ensign-repair-sonnet-live-flakes
mod-block:
pr:
---
## Problem

Exact main `c6a336a33f866536d2b1ea9ef313bbbc6f244dd1` gave opposite Sonnet outcomes on unchanged bytes in run `31894934349`: `TestLiveCommonKeepMovingPosture` passed attempt 1 and failed attempt 2 with `keep-moving-violation`. The claude-sonnet binding `liveXFail("claude-sonnet", "060xp69y61yhrww23g3wvwqy")` on that journey stays until the journey passes without it.

Scope (captain, 2026-11-03): this task owns the keep-moving flake only. The selected-bare half moved to `align-claude-break-glass-agent-proof`, which removes the Sonnet 060xp selected-bare binding.

### What the failing attempt did (decision point)

The recorded label "the FO bypasses the formal revise/feedback route for `questioned`" is not the failure. It comes from the attempt-2 final message ("no formal reviewer/feedback package existed yet, so I routed it directly"). `assertDurableKeepMoving` does not grade the revise route: for `questioned` it checks only that the entity changed status, got an entity-scoped stage-report commit, stayed nonterminal, and that its last commit is entity-owned. The passing attempt 1 routed `questioned` the same way (`gate record questioned --decision revise`, then `status --set questioned status=ideation started`).

The failing decision point in the attempt-2 stream (`runtime-live-e2e-claude-live-claude-sonnet-5`, artifact 9250257670) is the dispatch window. The FO stamped three entities and then committed them by hand before the `--stamp` builds:

```
status --set ready-one status=implementation started
status --set ready-two status=implementation started
status --set questioned status=ideation started
"Now let's build and check dispatch commits for each, then dispatch the four workers."
git add -A && git commit -m "dispatch: ready-one, ready-two entering implementation; questioned entering ideation (feedback route)"
  -> 8 files: approved-gate.md, approved-gate/review/review/briefing-1/*, questioned.md, questioned/review/review/briefing-1/*, ready-one.md, ready-two.md
dispatch build ... --stamp   (x4)
```

Because `started` was already committed, the `--stamp` builds for `ready-one` and `ready-two` had nothing to stage. `stampCommitInline` returns without a commit when the index is empty (`internal/dispatch/stamp.go`), so the exact `dispatch: ready-one entering implementation` commit never exists. The batch commit cannot stand in for it: `durableBatchDispatchProof` rejects it because it also carries `approved-gate.md` (consumed, no `started`) and two gate rooms.

### Root cause

This is a defect in the FO contract text, with model variance as the trigger. `fo-dispatch-core.md` step 5 tells the FO to write `started` itself (`status --set {slug} status={dispatch_stage} started`) before `--stamp`. That leaves a window in which a hand commit uses up the only change `--stamp` would commit. The binary works as designed. The oracle is correct: the contract requires the exact path-scoped `dispatch: {slug} entering {dispatch_stage}` commit.

### Spike (deterministic replay, real binary, live oracle)

A throwaway test on the probe checkout builds the head `keep-moving/mixed-events` fixture and replays the attempt-2 command order with the real binary: consume, revise, advance, the stray `git add -A` commit, four `--stamp` builds, worker reports, `merge guard --verdict passed`, and the questioned re-review. It then calls `assertDurableKeepMoving`:

- Control (advance then `--stamp` per entity, no stray commit): `nil`.
- Attempt-2 order with `started` in the advance: `durable keep-moving journeys = 1/3: ready-one: missing path-scoped dispatch entry with stage and started; ready-two: (same)`. This is the CI red, reproduced exactly.
- Attempt-2 order with `started` dropped from the advance: `nil`. `--stamp` writes `started` and makes its own path-scoped `dispatch:` commit, even after the stray commit.

Also checked: `status --set e status=b started` does not overwrite an existing `started` value. `--stamp` fills `started` only when it is empty (`needStarted`). Moving the write into `--stamp` therefore keeps the stored value rule.

### Other Sonnet keep-moving reds (already handled, not this task)

- Run 40447244422 (2026-10-27, `861a6e13a`, `claude-sonnet-5-5`): XFAIL `keep-moving-violation`. The FO terminalized with `status --set {e} status=done verdict=PASSED completed` + `state commit`, not `merge guard`. Commit `0841d0165` on this branch fixes that path: the contract names the guard for non-gated terminal stages, and the guard accepts consumed history. In all four head probes that got that far, the FO used `merge guard`.

## Proposed approach

Make `--stamp` the only writer of `started` on the non-gated dispatch path. Then a stray FO commit cannot use up the dispatch commit. This is one sentence in `skills/first-officer/references/fo-dispatch-core.md` step 5, plus a deterministic test that runs the contract's own advance command.

Before:

> For a non-gated entry, first advance with `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={dispatch_stage} started`, then run the same `--stamp --stage {dispatch_stage}` build.

After:

> For a non-gated entry, first advance with `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={dispatch_stage}` (no `started`: the `--stamp` build writes it inside its own path-scoped `dispatch:` commit), then run the same `--stamp --stage {dispatch_stage}` build with no commit in between.

Alternatives considered:
1. Retire the binding on evidence alone, with no product change. This fails the "smallest sufficient" test only if the gate wants protection against model variance. The deviation reproduced in 1 of 2 CI attempts on Sonnet 5. It did not occur in any of 10 head probes on Sonnet 5.5, 2 of which were graded (see Evidence). This option is acceptable if the captain prefers no prose change.
2. Prose only ("do not commit between the advance and the stamp"). This is weaker: it asks the model to follow a rule, while dropping `started` removes the window mechanically.
3. Make `--stamp` create a dispatch commit when there is nothing to stage, or loosen `durableBatchDispatchProof`. This is rejected: an empty commit or a looser oracle accepts unscoped batch commits that sweep in other entities' gate state, and the contract forbids such commits.

Then retire the claude-sonnet keep-moving binding per AC-3.

Expected surface: `skills/first-officer/references/fo-dispatch-core.md` (+1/-1), `internal/ensigncycle/shared_keep_moving_fixture_test.go` (new test, about +90), `internal/ensigncycle/shared_live_runner_test.go` (binding removal, +1/-1), `docs/runtime-live-ci-registry.md` (-1 row). Estimate net LOC change: +89, across 4 files (about 92 insertions, 3 deletions). Tolerance: ±30 net LOC, ±1 file.

Semantics: runtime behavior changes (the FO no longer writes `started` in the non-gated advance; `--stamp` writes it). Command grammar, stored formats (`started` keeps its fill-when-empty rule), and authority do not change. No user-facing doc changes; at retirement the registry row `keep-moving-posture | claude-sonnet | 060xp69y61yhrww23g3wvwqy | ...` is deleted.

No spike needed beyond the one recorded above. The proven mechanisms are `--stamp`'s fill-when-empty `started` plus its path-scoped commit, and `status --set` not overwriting `started`.

## Acceptance criteria

- **AC-1:** Sonnet keep-moving posture passes without an XFAIL binding: `TestLiveCommonKeepMovingPosture` on `claude-sonnet-5-5` grades PASS in 3 of 3 unbound runs on the candidate bytes, where the baseline was a graded red in 1 of 2 CI attempts on `c6a336a33` and 1 of 1 on `861a6e13a`. Tested by: exact live runs (local `SPACEDOCK_LIVE_MODEL=claude-sonnet-5-5`, or the CI Sonnet lane). An infra kill (quiet budget, auth) is not a graded run.
- **AC-2:** A stray FO commit between the non-gated advance and `--stamp` no longer removes the path-scoped dispatch proof: with the contract's step-5 advance command, then `git add -A && git commit`, then `--stamp`, `assertDurableKeepMoving` returns nil. The same sequence with `started` added to the advance returns `ready-one: missing path-scoped dispatch entry`. Tested by: a deterministic Go test that reads the advance command from `fo-dispatch-core.md` and runs it. Falsifying edit: put `started` back in step 5, or stop `--stamp` from stamping an empty `started`.
- **AC-3:** The claude-sonnet `060xp69y61yhrww23g3wvwqy` binding on keep-moving is removed only after a graded bound XFAIL or XPASS on the candidate, followed by an unbound PASS on otherwise unchanged bytes. The registry row is removed in the same change. Tested by: the run records cited in the validation report, and `TestRuntimeLiveRegistryMatchesBindings`.
- **AC-4:** Registry reconciliation, active-owner join, `go test ./...`, `go test ./... -race` and the exact Sonnet live check pass.

## Test plan

- Primary proof owner: `TestLiveCommonKeepMovingPosture` (live, about $3 and 7 to 18 min per run). It covers AC-1 and AC-3.
- New deterministic test (AC-2): reuse `writeKeepMovingWorkflow` and one built binary. Replay the attempt-2 order with subtests `contract-advance` (must pass) and `advance-with-started` (must fail with the missing-dispatch reason, which shows the test bites). About 60 s, because it builds the binary once. This is executable behavior, not a substring check: the test runs the command the contract tells the FO to run.
- Live retirement cost: at least 1 bound and 3 unbound graded Sonnet runs, about $12 to $15.
- Risk to watch: on unchanged head bytes, 7 of 8 local runs were killed by the 60 s no-progress quiet budget (the 8th by an auth 401). The binding does not absorb those kills; they fail the lane. The quiet budget has its own finding under Evidence. Retirement runs come from the parent task's validator and the CI Sonnet lane, which ran this journey in about 62 s.
- Second risk: the "background worker stopped" finding under Evidence is a keep-moving failure that the step-5 fix does not address. If it recurs, AC-1 cannot reach 3 of 3.

## Evidence

- Run `31894934349` (`claude-sonnet-5`), attempt 1: keep-moving PASS; attempt 2: keep-moving FAIL `keep-moving-violation` (decision point quoted above).
- Run `40447244422` (`claude-sonnet-5-5`, `861a6e13a`): keep-moving XFAIL from merge-guard bypass, fixed by `0841d0165`.
- Local probes, first round, on `64274cef5` + a test-only diagnostic patch (logs XFAIL details and keeps the graded root), Sonnet `claude-sonnet-5-5`, `--effort max` shim: 4 keep-moving runs, 0 graded (3 quiet-budget kills, 1 auth 401). In all 4, the FO chained `status --set … started && dispatch build --stamp` per entity with no hand commit before `--stamp`, and used `merge guard` for terminal stages.
- Local probes, second round (2026-11-03, rotated token, 060xp binding in place), Sonnet `claude-sonnet-5-5`, `--effort max` shim as in CI, unchanged `64274cef5` bytes. Three other workers shared the machine. Total 4 runs, 0 graded, 0 auth errors (`401` and `api_retry` counts were 0 in every stream). Each kill has the same shape: the FO finished `approved-gate`, `ready-one` and `ready-two` through `merge guard`, then ended its turn while the background `Questioned: ideation` worker was in a silent turn.

  | Run | Duration | Grade | Hand commit before `--stamp` | Stalled step |
  |---|---|---|---|---|
  | a1 | 682 s | infra kill (60 s quiet budget) | no | `Questioned: ideation` worker, next model turn after `Read fo-status-viewer/SKILL.md` at 19:49:50Z |
  | a2 | 811 s | infra kill (60 s quiet budget) | no | same worker, last line 20:06:05Z (read `/tmp/spacedock-dispatch/feedback-questioned-ideation-cycle1.md`); FO turn ended 20:06:23Z |
  | a3 | 710 s | infra kill (60 s quiet budget) | no | same worker, last line 20:21:52Z (a scratch git symlink experiment); FO turn ended 20:21:53Z |
  | a4 | 732 s | infra kill (60 s quiet budget) | no | same worker, last line 20:39:22Z (a `status` validate baseline); FO turn ended 20:39:42Z |

- Behavior evidence, not retirement evidence: 2 runs on the same checkout with `quietBudgetDefault` raised from 60 s to 180 s, the binding in place, and `-test.timeout 45m`.

  | Run | Duration | Grade and codes | Hand commit before `--stamp` | Notes |
  |---|---|---|---|---|
  | b1 | 1647 s | XPASS, `observed=[]` | no (the FO made no `git commit` at all) | Exact `dispatch: {slug} entering {stage}` commits for all 4 entities. `merge guard` for the 3 terminal entities. Silent gaps of 128, 104, 100, 97 and 78 s, so the 60 s budget would have killed this run. |
  | b2 | 1216 s | XFAIL, `observed=[keep-moving-violation]` | no (each advance chained `status --set … started && dispatch build --stamp`) | Cause: the "background worker stopped" finding below. Largest silent gap was 59 s, so the unchanged 60 s budget would not have fired. |

- Side observations, out of scope here: `/tmp/spacedock-dispatch/` is shared across sessions. In a2 the worker found an unprefixed `feedback-questioned-ideation-cycle1.md` from an earlier session (dated 11:08 local) with `ls` and read it. The secondary "FO booted the wrong root" and "FO broad-searched the filesystem at boot" diagnostics fired on worker scratch commands in a2, a3 and a4, not on FO boot.

### Finding: quiet budget kills runs while a background worker is in a silent turn

1. Released user and normal workflow: none for released users. The affected workflow is a maintainer running the documented local Claude live lane (`SPACEDOCK_LIVE_RUNTIME=claude`, `claude-sonnet-5-5`, effort max as in CI).
2. Observable harm: 7 of 8 local unchanged-byte keep-moving runs (both rounds) ended with no grade, at about $3 and 12 min each. Local runs cannot supply retirement evidence.
3. Affected value: value-ac[AC-1] names local Sonnet runs as one source of graded evidence, and that source produces no grades under this load.
4. Trigger evidence: `drainToExit` (`internal/ensigncycle/streamwatch_test.go`) resets the 60 s deadline only when a stream line arrives. A background worker's max-effort turn writes no line until the turn ends. In a1 to a4 the FO had ended its turn, so no other source wrote lines. b1 had five gaps of 78 to 128 s. The kill depends on load: the parent task's implementation worker got keep-moving to pass 3 of 3 locally earlier, CI ran it in about 62 s, and three workers shared this machine.

Proposed materiality: Deferred risk. Promote it to Material if the CI Sonnet lane gets the same kill. Proposed owner: a new harness task that owns `quietBudgetDefault`, not this task. Proposed disposition: route for decision. A candidate mechanism is to extend the budget while a background dispatch is open, as the sanctioned `quietBudgetDispatchClose` exemption already does for a similar case.

### Finding: background worker stopped before it reported

1. Released user and normal workflow: a Claude FO in a headless `claude -p` session that dispatches background workers and keeps moving. This is the path the keep-moving journey exercises.
2. Observable harm: in b2 the `questioned` ideation worker was stopped mid-stage with no report and no commit, so `questioned` has no entity-owned stage re-shape. The lane grades `keep-moving-violation`: XFAIL under the binding, FAIL without it.
3. Affected value: value-ac[AC-1] requires 3 of 3 unbound PASS, and the step-5 fix does not prevent this failure.
4. Trigger evidence (b2): the worker started at 21:18:29Z and ran nine background stress experiments on symlink and hardlink transports. In its first message to the FO it said the entity had no design text to re-shape. The FO's last line was at 21:23:35Z ("waiting on that completion signal"). At 21:33:35Z, 600 s later, the host injected `[Request interrupted by user]` into the worker. Its `task_notification` reported `status: stopped`, `background_tasks_changed` reported `tasks: []`, and the process exited. The 600 s value was seen once. It fits a host limit on how long background tasks may run after the main turn ends, but host documentation does not confirm it. This happened in 1 of 2 graded runs. a3 and a4 show the same worker doing similar symlink experiments.

Proposed materiality: Material, because it is observed and blocks AC-1. Proposed owner: this task, because it causes the same binding's red on the same journey. Proposed disposition: route for decision, because the remedy is outside the approved step-5 surface. Two options are FO headless wait posture (do not end the turn idle with an open background worker) and ensign stage proportionality (a deep spike in ideation on a fixture entity).
- Selected-bare (moved to `align-claude-break-glass-agent-proof`), local Sonnet 5.5 for that task's use: default prompt bound 3 of 3 XPASS (FO passed `--bare-mode`); with a `/spacedock bare` cue 3 of 3 XPASS. The Aug attempt-1 red came after the FO probed `ToolSearch select:SendMessage` and chose team ("team mode, since SendMessage/Agent are available").

## Stage Report: ideation

- DONE: With local Sonnet probes on a throwaway checkout of the live-lane-shared-failure-fixes branch head, reproduce and quote the exact decision point for each recorded flake: keep-moving (FO bypasses the formal revise/feedback route for questioned) and break-glass selected-bare (Agent call does not preserve blocking bare mode), or show it no longer reproduces, with run counts.
  Probe base `64274cef5` (detached). Keep-moving: the decision point is the hand `git add -A` batch dispatch commit before `--stamp` (quoted in Problem), not the revise route, which the oracle does not grade. A deterministic replay of that order on head reproduces `ready-one/ready-two: missing path-scoped dispatch entry` (control: nil). Live on head, both rounds: 8 unchanged-byte runs, 0 graded (7 quiet-budget kills, 1 auth 401), plus 2 behavior runs with a 180 s budget (b1 XPASS; b2 XFAIL `keep-moving-violation` from a background worker the host stopped). 0 hand commits before `--stamp` in all 10. Selected-bare moved to `align-claude-break-glass-agent-proof` by captain scope change; my head counts (3/3 and 3/3 XPASS) are recorded under Evidence for that task.
- DONE: Name the root cause of each as an FO/ensign contract, binary, harness, or oracle defect, and propose the smallest fix with a falsifiable AC and a net LOC estimate.
  Keep-moving: FO contract defect. Step 5 has the FO write `started` before `--stamp`, so a stray commit uses up the dispatch commit. Fix: drop `started` from the advance (one sentence) plus a deterministic test that runs the contract's command (AC-2). Net +89 LOC across 4 files. Two more findings under Evidence. The quiet-budget kills are a harness defect, proposed Deferred risk, with a new harness owner. The background worker stopped by the host (b2) is a second keep-moving cause, proposed Material and route for decision, because the step-5 fix does not cover it.
- DONE: State the retirement evidence each 060xp binding needs (bound failure or XPASS, then unchanged-byte unbound PASS) and whether the probes already supply it.
  Keep-moving needs a graded bound XFAIL or XPASS on the candidate, then 3 of 3 unbound PASS on unchanged bytes (AC-1, AC-3). The probes do not supply it. The unchanged-byte runs produced 0 grades. The 2 graded runs used a raised budget and are behavior evidence only, by FO decision. Retirement evidence comes from the parent task's validator and the CI Sonnet lane.

### Summary

The keep-moving flake comes from a hand batch commit in the FO's dispatch window, not from the revise route. A replay with the real binary and the live oracle reproduces the CI red exactly. The smallest fix is to let `--stamp` alone write `started` (one sentence in `fo-dispatch-core.md`), proven by a deterministic replay test. In 10 head probes on Sonnet 5.5 the FO never made that hand commit. The unchanged-byte runs produced no grades, because the 60 s quiet budget kills a run while a background worker is in a silent max-effort turn. That is a separate harness finding. Two behavior runs with a 180 s budget graded XPASS and XFAIL. The XFAIL is a second keep-moving cause: the host stopped the `questioned` ideation worker 600 s after the FO's last turn, before it reported. The step-5 fix does not cover this, so it is routed for decision. Retirement evidence remains open and comes from the parent task's validator and CI. The selected-bare half moved to `align-claude-break-glass-agent-proof`.

## Stage Report: implementation

- DONE: Change step 5 of skills/first-officer/references/fo-dispatch-core.md so the non-gated advance writes only status={dispatch_stage} and --stamp writes started inside its own path-scoped dispatch commit, with no commit in between.
  Commit `0542c94ad` on `repair-sonnet-live-flakes` (on top of `64274cef5`): the non-gated advance is now `status --set {slug} status={dispatch_stage}` (no `started`), then the same `--stamp` build "with no commit in between". The gate-consumed sentence and the parent's step-53 edits are unchanged.
- DONE: Add a small deterministic test that hardcodes the advance sequence (advance without started, a stray git add -A commit, then dispatch build --stamp) and shows the path-scoped dispatch commit exists; it goes red when started is written in the advance. The test does not read fo-dispatch-core.md or any other instruction file.
  `TestKeepMovingStrayCommitKeepsDispatchProof` (`internal/ensigncycle/keep_moving_stamp_window_test.go`, about 2 s) builds the real binary and replays advance, a stray `git add -A` commit that also sweeps `questioned.md`, then `dispatch build --stamp`. It grades with the live oracle's `durableDispatchProof`. Subtest `advance-without-started` asserts the exact `dispatch: ready-one entering implementation` commit exists and `started` is set. Subtest `advance-with-started` is the negative control: it asserts the proof is missing. Falsified: when I added `started` to the first subtest's advance, it failed with `path-scoped dispatch proof = false, want true`. The test would also fail if `--stamp` stopped filling an empty `started`. The test reads no instruction file.
- DONE: Keep the 060xp keep-moving binding and its registry row; measured diff within net +59..+119 across 3 to 5 files; gofmt, go test ./..., and go test ./... -race pass (per package with -timeout 30m if needed).
  The binding `liveXFail("claude-sonnet", "060xp69y61yhrww23g3wvwqy")` (`shared_live_runner_test.go:173`) and registry row 21 are not changed. `git diff --numstat 64274cef5..HEAD` gives +71/-1, net +70, within the LOC band. The diff touches only 2 files, below the 3-file floor. The 4-file estimate counted the binding removal and the registry row, which now stay, and I did not add a file only to meet the count. `go test -count=1 ./...` and `go test -race -count=1 ./...` pass (21/21 packages ok), with `-p 2 -timeout 30m`. An earlier run failed only with `no space left on device` and passed again after disk space came back. `gofmt -l ./cmd ./internal` shows only `internal/release/runtime_live_evidence_workflow_test.go`, which was already unformatted at `64274cef5` (last changed by `e68aa5b7a`). It is not in this diff and I did not change it.

### Summary

Step 5 no longer tells the FO to write `started` on the non-gated advance. `--stamp` is now the only writer, inside its own path-scoped `dispatch:` commit, so a stray hand commit can no longer use up the dispatch proof. A 2-second deterministic replay with the real binary and the live oracle's dispatch check proves this. A negative control shows the proof disappears when the advance writes `started`. The surface is 2 files, net +70, so the file count is one below the 3-to-5 band, because the binding and registry row stay as directed. Note for retirement: the registry row text ("the FO skips the formal revise route") gives a cause that ideation showed is wrong. I did not change it, because the row stays with the binding.
