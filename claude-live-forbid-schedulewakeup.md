---
id: ef4p1sbvkk9zxjxm0ez9npdr
title: Claude live harness forbids ScheduleWakeup so haiku runs exit when the work is done
status: implementation
source: "Captain directive, 2026-11-10: the fix is a test-harness injection on haiku that forbids the tool, not a first-officer instruction."
started: 2026-11-10T18:36:35Z
completed:
verdict:
score: 0.6
worktree: .worktrees/spacedock-ensign-claude-live-forbid-schedulewakeup
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:ef4p1sbvkk9zxjxm0ez9npdr:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ef4p1sbvkk9zxjxm0ez9npdr-backlog-1
              briefing:
                id: briefing:ef4p1sbvkk9zxjxm0ez9npdr:backlog:attempt-1:revision-1
                digest: sha256:ecbd5ac3e78d627a49fcd0be234df9b0515c0f45e3a56dc8c7d4e50863ec533b
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:ef4p1sbvkk9zxjxm0ez9npdr:backlog:1
                briefing: briefing:ef4p1sbvkk9zxjxm0ez9npdr:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-10T18:34:22.537219Z"
                decision: approve
                reason: Seed states the problem with run evidence, a bounded harness-only scope, and the captain's chosen mechanism; ready for ideation.
                conn:
                    quote: for haiku - dispatch this with the disallowed tool change the remaining failing marked as xpass. create PR onto the current stack and we can run ci
                    source: captain chat, 2026-11-10
              application:
                target-stage: ideation
                state: consumed
        - id: gate:ef4p1sbvkk9zxjxm0ez9npdr:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ef4p1sbvkk9zxjxm0ez9npdr-ideation-1
              briefing:
                id: briefing:ef4p1sbvkk9zxjxm0ez9npdr:ideation:attempt-1:revision-1
                digest: sha256:f286fba584c01092850dbcf7db9842e94668d7a85a04a039d912ce6bcff1c379
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:ef4p1sbvkk9zxjxm0ez9npdr:ideation:1
                briefing: briefing:ef4p1sbvkk9zxjxm0ez9npdr:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-10T18:43:33.306267Z"
                decision: approve
                reason: Design is the captain's harness-only haiku restriction pinned by a deterministic test, binds only journeys graded red with the restriction applied, measures the end-value on a CI haiku run against the 7-kill baseline, and declares a 5-file +30 LOC surface on the ci/go-126-workflows base.
                conn:
                    quote: for haiku - dispatch this with the disallowed tool change the remaining failing marked as xpass. create PR onto the current stack and we can run ci
                    source: captain chat, 2026-11-10
              application:
                target-stage: implementation
                state: consumed
---

On the claude-haiku live lane, journeys that complete their work are killed as hung. The haiku first officer calls the host tool `ScheduleWakeup` as a fallback heartbeat while waiting on background workers; after the work finishes and the final turn ends, the pending wakeup keeps the `claude -p` process alive and silent, and the harness kills it after its 60-second no-progress budget, before grading.

## Problem

- Run 42865474062 (`ci/go-126-workflows` @ `eef2d1302`, model `claude-haiku-5-5`, 2026-11-10) failed in all three attempts. Every failure except one graded `ac-value-reanchor` red was a kill with "made no stream progress within 1m0s (no-progress quiet budget)":
  - attempt 1: `owned-conflict-owner-handoff`, `smallest-sufficient-mechanism`
  - attempt 2: `rejection-flow`
  - attempt 3: `rejection-flow`, `bare-reachable`, `default-headless-gate-stop`, `recorded-gate-lifecycle`
- The three attempts produced 66 haiku streams. Exactly 7 contain a `ScheduleWakeup` call with a delay (1200 s, once 60 s), and those 7 are the 7 killed streams. Each killed stream ends with a final `end_turn`, then 60 s of silence. The same journeys exit about a second after their last event in any attempt where they make no scheduling call.
- An xfail binding does not absorb a kill: the kill fails the test before grading. `default-headless-gate-stop` carries a `claude-haiku` binding and still failed in attempt 3.
- Three green sonnet runs (66 streams: 41838996310, 42848106064, 42849771310) contain no delayed scheduling call; one stream calls `ScheduleWakeup` with `stop: true` only.
- No file under `skills/` mentions `ScheduleWakeup`; it comes from the Claude Code host tool list.
- The kill is `drainToExit` in `internal/ensigncycle/streamwatch_test.go`, run with `quietBudgetDefault` by `claudeLiveRunner.run` in `internal/ensigncycle/claude_live_runner_test.go`, which builds the `claude` argument list.

Evidence locations: attempts 1 and 2 under `/tmp/live-hist/42865474062{,-rerun}` (`-rerun3` is a byte-identical copy of attempt 2); attempt 3 is artifact 13787964968 of the run.

## Spike: the restriction removes the tool and the kills

A local haiku run on `main` @ `5785dbed7` (Claude Code 2.1.317) with `--disallowedTools ScheduleWakeup` appended after `--` in `claudeLiveRunner.run` (`/tmp/haiku-local-run1.jsonl`, `/tmp/haiku-local/run1`):

- `ScheduleWakeup` is absent from the session tool list of all 18 shared-scenario streams, and no stream calls it.
- No journey was killed for stream silence. All six journeys killed in CI completed and were graded; four passed.
- Remaining red: `rejection-flow` (graded `rejection-worker-topology`) and `TestLiveBareReachable` (no bare-shaped `Agent()` call; the first officer dispatched a named worker). `TestLiveMergedTeamModeDispatch` failed on a revoked local OAuth token before any model work; it has its own launch path and is not affected by this change.

The cause is still a correlation: a pending wakeup holding `claude -p` open has not been reproduced in isolation. It is 7 of 7 kills in streams with a scheduling call, 0 of 59 clean exits with one, and 0 kills once the tool is gone. The CI run on the pull request is the measurement that confirms or refutes it.

## Approach

1. Move the `spacedock claude` argument list out of `claudeLiveRunner.run` into an untagged helper in `internal/ensigncycle/claude_runtime_helpers_test.go`, beside `unsetNestedSessionArgs`:

   ```go
   func claudeFrontDoorArgs(pluginDir, prompt, model string) []string
   ```

   It returns the current list unchanged and, when `model` contains `haiku`, appends `--disallowedTools ScheduleWakeup` after `--` so it forwards to `claude`. `run` calls it with `prompt + " " + antiShutdownOverride` and `r.modelName`. The helper is untagged because `claude_live_runner_test.go` builds only under `-tags live`, and the pin must run in `go test ./...`.

2. Haiku only. Sonnet made no delayed scheduling call in 66 streams and its lane is green, so restricting it would change a working lane's tool list with no fault to fix. The captain scoped the fix to haiku.

3. Add `claude-haiku` xfail bindings (owner blank, the convention for the existing haiku bindings) to the two journeys that stayed red with the restriction applied. XPASS does not fail the lane (`liveGradeFailsLane`), so a binding on a journey that sometimes passes logs `XPASS ALERT` and stays green.

The merged-team-mode launch (`merged_team_mode_live_test.go`) and the haiku loop spike build their own argument lists. Neither made a scheduling call in any attempt, and they stay unchanged.

## Haiku xfail bindings

Captain direction: "the remaining failing marked as xpass". This design reads that as: **a haiku journey that is still graded red with the restriction applied gets a `claude-haiku` xfail binding**. Journeys that were only killed, or that were red only without the restriction, get no binding. The captain can correct this reading at the gate.

Per-journey results for attempts 1, 2, 3 and the restricted local run (L):

| Journey | 1 | 2 | 3 | L | Decision |
|---|---|---|---|---|---|
| `rejection-flow` | pass | kill | kill | red `rejection-worker-topology` | **bind** |
| `TestLiveBareReachable` | pass | pass | kill | red: no bare-shaped `Agent()` | **bind** |
| `ac-value-reanchor` | red `ac-value-reanchor-violation` | pass | pass | pass | no binding |
| `owned-conflict-owner-handoff`, `smallest-sufficient-mechanism`, `recorded-gate-lifecycle` | killed once | | | pass | no binding |

- `rejection-flow`: the restriction does not fix it. The kills in attempts 2 and 3 hid the graded topology fault, which is the same fault class the Pi binding on this journey already owns. The binding is `liveXFail("claude-haiku", "")` added to its gaps list.
- `TestLiveBareReachable`: it failed graded with the restriction applied, so it meets the rule, though it passed both graded CI runs. It is a plain test, not a `liveJourney`, so the binding needs a gap function, `bareReachableGap(target)`, shaped like `breakGlassSelectedBareGap` in the same file. contractlint only records a gap built in a non-Test function as a binding. The test then grades through `finishLiveScenario` instead of `t.Fatalf`, with the role from `claudeLiveRole(runner.model())`. For unbound targets the behavior does not change: red still fails and green still passes. Dropping this binding saves about 12 LOC and one registry row, and leaves the lane red about one run in three on this evidence.
- `ac-value-reanchor`: its one failure came before the restriction and is a gate-routing fault, not a kill. It passed in the only restricted run, so it does not meet the rule. If it is red in the pull request's CI run, it gets a binding then.

Existing haiku bindings stay as they are:

- `withdrawn-gate-recovery` and break-glass selected-bare: XPASS in all four runs.
- `default-headless-gate-stop`: XPASS in its three graded runs (killed in attempt 3).
- `keep-moving-posture`: XFAIL in attempts 1 and 3, XPASS in attempt 2 and L, so it is still red.

**Open question for the captain:** should the three bindings that never went red (`withdrawn-gate-recovery`, break-glass selected-bare, `default-headless-gate-stop`) be removed? This task does not remove them.

## Acceptance criteria

1. A haiku live run in CI on the pull request has no "made no stream progress" kill, and no stream calls `ScheduleWakeup`. The baseline is run 42865474062: 7 kills in 3 attempts, every one after a delayed `ScheduleWakeup` call. Tested by the `claude-live (…, claude-haiku-5-5, …)` job on the stacked pull request: grep the job log for `no stream progress within`, and grep the artifact's `claude-stream.jsonl` files for `"name":"ScheduleWakeup"`. Both counts must be 0.
2. The built `claude` argument list carries `--disallowedTools ScheduleWakeup` after `--` for the registered haiku model and the `haiku` alias, and not for sonnet or opus. Tested by `TestClaudeFrontDoorArgs` (untagged, runs in `go test ./...`).
3. `rejection-flow` and `TestLiveBareReachable` grade XFAIL when red on `claude-haiku` and still FAIL when red on `claude-sonnet`. Every binding has a matching row in `docs/runtime-live-ci-registry.md`. Tested by the CI haiku job log (`XFAIL claude-haiku/rejection-flow` or `XPASS ALERT …`, and the same for `bare-reachable`) and by `TestRuntimeLiveRegistryMatchesBindings`.
4. Skill text and sonnet launch arguments are unchanged. Tested by `git diff origin/ci/go-126-workflows --stat -- skills/` being empty and by the sonnet cases of AC-2.

## Test plan

- `TestClaudeFrontDoorArgs` (new, deterministic, cheap) sits beside `TestUnsetNestedSessionArgs`. It is a table over `release.LiveModels["claude.haiku"]`, `haiku`, `release.LiveModels["claude.sonnet"]`, `sonnet` and `release.LiveModels["claude.opus"]`.
  - It asserts the adjacent pair `--disallowedTools`, `ScheduleWakeup` exists, at an index after the `--` separator, exactly for the haiku cases, and that `--model <model>` is present in every case.
  - Falsifying edits: delete the haiku branch (haiku cases fail); make the append unconditional (sonnet and opus fail); put the pair before `--` (ordering check fails); rename the registered haiku ID to one without `haiku` (registry case fails).
  - This test pins the helper, not that `run` calls it. That gap is covered by AC-1, the live run. A stream whose init tool list still contains `ScheduleWakeup` would show it.
- `TestRuntimeLiveRegistryMatchesBindings` (existing primary owner of binding/registry agreement) fails if either new binding lacks its registry row, or a row lacks its binding.
- Live (one CI run, the measurement): the haiku `claude-live` job on the stacked pull request, read for AC-1 and AC-3. Cost: one roughly 7-minute haiku lane run.
- Before claiming done: `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `go vet -tags live ./internal/ensigncycle` (the live-tagged files compile only under the tag).

## Expected surface

Estimate net LOC change: +30, across 5 files (about +46 / −16). Tolerance: ±15 net LOC; a sixth file is a boundary breach.

| File | Change | +/− |
|---|---|---|
| `internal/ensigncycle/claude_runtime_helpers_test.go` | `claudeFrontDoorArgs` and `TestClaudeFrontDoorArgs` | +30 / 0 |
| `internal/ensigncycle/claude_live_runner_test.go` | `run` calls the helper in place of the literal | +1 / −10 |
| `internal/ensigncycle/dispatch_recovery_live_test.go` | `bareReachableGap`; `TestLiveBareReachable` grades through `finishLiveScenario` | +10 / −5 |
| `internal/ensigncycle/shared_live_runner_test.go` | `rejection-flow` gaps gain `liveXFail("claude-haiku", "")` | +1 / −1 |
| `docs/runtime-live-ci-registry.md` | two rows | +2 / 0 |

Semantics this task may change: Claude live harness runtime behavior on haiku model names only (the host tool list lacks `ScheduleWakeup`), and live grading of `rejection-flow` and `bare-reachable` on `claude-haiku` (red becomes XFAIL). No command grammar, stored format, authority, product binary or skill text changes.

## Documentation diff

`docs/runtime-live-ci-registry.md` gets two rows appended to the table. contractlint requires them.

```diff
 | `breakGlassSelectedBareGap` (break-glass selected-bare, `dispatch_recovery_live_test.go`) | claude-haiku | unowned (blank; pending analysis) | claude-haiku red `break-glass-observables` |
+| rejection-flow | claude-haiku | unowned (blank; pending analysis) | claude-haiku red `rejection-worker-topology` |
+| `bareReachableGap` (bare-reachable, `dispatch_recovery_live_test.go`) | claude-haiku | unowned (blank; pending analysis) | claude-haiku red `bare-reachable-observables`: no bare-shaped `Agent()` call |
```

No user-facing documentation changes. The flag lives only in the test harness.

## Branch and delivery

The implementation branch starts from `ci/go-126-workflows` (PR #849, tip `eef2d1302`), not `main`. The pull request targets `ci/go-126-workflows` so that CI runs the haiku lane on the stack.

## Alternatives considered

- **Instruct the first officer not to schedule wakeups (skill text).** The captain ruled this out. It would also change product behavior to fix a harness-only symptom.
- **Raise the no-progress quiet budget.** A 1200 s wakeup outlasts any reasonable budget, and a longer budget weakens hang detection for every journey.

## Stage Report: ideation

- DONE: Acceptance criteria name the measured end-value (a haiku CI live run with no no-stream-progress kill following a ScheduleWakeup call) and the deterministic test that pins the injected restriction in the built claude argument list
  AC-1 measures kills and ScheduleWakeup calls in the stacked PR's haiku job against a baseline of 7 kills in 3 attempts. AC-2 is the untagged `TestClaudeFrontDoorArgs`, with four falsifying edits named.
- DONE: The design names exactly which haiku journeys gain a claude-haiku xfail binding, with the run evidence for each, and states what happens to haiku bindings that passed in all four runs
  `rejection-flow` and `TestLiveBareReachable` are bound. `ac-value-reanchor` is not. The per-run table covers attempts 1, 2, 3 and the local run. The three never-red bindings stay; removing them is an open question for the captain.
- DONE: Expected surface (files, net LOC, tolerance) and the stacked-PR base ci/go-126-workflows are declared
  Net +30 (+46/−16) across 5 files, ±15 net. Base `ci/go-126-workflows` @ `eef2d1302` (PR #849). The registry doc diff is included.

### Summary

The design is small. It moves the `claude` argument list into an untagged helper that appends `--disallowedTools ScheduleWakeup` for haiku models only, adds `claude-haiku` xfail bindings to the two journeys still red with the restriction, and adds the two registry rows that contractlint requires. Re-pulling the CI logs showed that the local `-rerun3` evidence was a copy of attempt 2. Attempt 3 actually had 4 more kills (bare-reachable, rejection-flow, default-headless-gate-stop, recorded-gate-lifecycle), each after a 1200 s ScheduleWakeup. The body now uses the corrected 7-of-7 correlation and the corrected binding history: keep-moving-posture was red twice, and default-headless-gate-stop was killed once. The bare-reachable binding needs a non-Test gap function, because contractlint ignores gap literals inside Test functions.
