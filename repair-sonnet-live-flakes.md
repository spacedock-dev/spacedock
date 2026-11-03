---
title: Repair Sonnet keep-moving and selected-bare live flakes
status: ideation
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
started: 2026-11-03T18:28:41Z
worktree:
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
1. Retire the binding on evidence alone, with no product change. This fails the "smallest sufficient" test only if the gate wants protection against model variance. The deviation reproduced in 1 of 2 CI attempts on Sonnet 5. It did not occur in any of 4 head probes on Sonnet 5.5, but none of those probes was graded (see Evidence). This option is acceptable if the captain prefers no prose change.
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
- Risk to watch: in local probes, 3 of 4 runs were killed by the 60 s no-progress quiet budget while the `questioned` ideation worker was silent (the 4th by an auth 401). The binding does not absorb those kills; they fail the lane. CI Sonnet 5.5 ran this journey in 62 s, so CI may be the practical place for retirement runs. If local kills keep happening, the quiet budget needs its own owner.

## Evidence

- Run `31894934349` (`claude-sonnet-5`), attempt 1: keep-moving PASS; attempt 2: keep-moving FAIL `keep-moving-violation` (decision point quoted above).
- Run `40447244422` (`claude-sonnet-5-5`, `861a6e13a`): keep-moving XFAIL from merge-guard bypass, fixed by `0841d0165`.
- Local probes on `64274cef5` + a test-only diagnostic patch (logs XFAIL details and keeps the graded root), Sonnet `claude-sonnet-5-5`, `--effort max` shim: 4 keep-moving runs, 0 graded (3 quiet-budget kills, 1 auth 401). In all 4, the FO chained `status --set … started && dispatch build --stamp` per entity with no hand commit before `--stamp`, and used `merge guard` for terminal stages.
- Selected-bare (moved to `align-claude-break-glass-agent-proof`), local Sonnet 5.5 for that task's use: default prompt bound 3 of 3 XPASS (FO passed `--bare-mode`); with a `/spacedock bare` cue 3 of 3 XPASS. The Aug attempt-1 red came after the FO probed `ToolSearch select:SendMessage` and chose team ("team mode, since SendMessage/Agent are available").

## Stage Report: ideation

- DONE: With local Sonnet probes on a throwaway checkout of the live-lane-shared-failure-fixes branch head, reproduce and quote the exact decision point for each recorded flake: keep-moving (FO bypasses the formal revise/feedback route for questioned) and break-glass selected-bare (Agent call does not preserve blocking bare mode), or show it no longer reproduces, with run counts.
  Probe base `64274cef5` (detached). Keep-moving: the decision point is the hand `git add -A` batch dispatch commit before `--stamp` (quoted in Problem), not the revise route, which the oracle does not grade. A deterministic replay of that order on head reproduces `ready-one/ready-two: missing path-scoped dispatch entry` (control: nil). Live on head: 4 runs, 0 graded (3 quiet-budget kills, 1 auth 401), 0 hand commits before `--stamp`. Selected-bare moved to `align-claude-break-glass-agent-proof` by captain scope change; my head counts (3/3 and 3/3 XPASS) are recorded under Evidence for that task.
- DONE: Name the root cause of each as an FO/ensign contract, binary, harness, or oracle defect, and propose the smallest fix with a falsifiable AC and a net LOC estimate.
  Keep-moving: FO contract defect. Step 5 has the FO write `started` before `--stamp`, so a stray commit uses up the dispatch commit. Fix: drop `started` from the advance (one sentence) plus a deterministic test that runs the contract's command (AC-2). Net +89 LOC across 4 files. The live quiet-budget kills are flagged as a separate harness risk.
- DONE: State the retirement evidence each 060xp binding needs (bound failure or XPASS, then unchanged-byte unbound PASS) and whether the probes already supply it.
  Keep-moving needs a graded bound XFAIL or XPASS on the candidate, then 3 of 3 unbound PASS on unchanged bytes (AC-1, AC-3). The probes do not supply it: 0 graded runs. Live probing then stopped: `~/.claude/benchmark-token` is revoked, the harness supports only that token or `ANTHROPIC_API_KEY` (no keychain path), and the disk filled (ENOSPC).

### Summary

The keep-moving flake comes from a hand batch commit in the FO's dispatch window, not from the revise route. A replay with the real binary and the live oracle reproduces the CI red exactly. The smallest fix is to let `--stamp` alone write `started` (one sentence in `fo-dispatch-core.md`), proven by a deterministic replay test. Head probes on Sonnet 5.5 showed no deviation but produced no graded runs: 3 of 4 were killed by the 60 s quiet budget and 1 by auth expiry. Retirement evidence is therefore still open, and CI may be the practical lane for it. The selected-bare half moved to `align-claude-break-glass-agent-proof`.
