---
id: ef4p1sbvkk9zxjxm0ez9npdr
title: Claude live harness forbids ScheduleWakeup so haiku runs exit when the work is done
status: done
source: "Captain directive, 2026-11-10: the fix is a test-harness injection on haiku that forbids the tool, not a first-officer instruction."
started: 2026-11-10T18:36:35Z
completed: 2026-11-11T05:15:23Z
verdict: PASSED
score: 0.6
worktree: .worktrees/spacedock-ensign-claude-live-forbid-schedulewakeup
issue:
pr: pr-merge:873
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
        - id: gate:ef4p1sbvkk9zxjxm0ez9npdr:validation
          stage: validation
          attempts:
            - id: gate-attempt:ef4p1sbvkk9zxjxm0ez9npdr-validation-1
              briefing:
                id: briefing:ef4p1sbvkk9zxjxm0ez9npdr:validation:attempt-1:revision-1
                digest: sha256:39434e79be4d88b573ea12f3f2e1433dc7bd27a1efc3e4afa269744969ce3304
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:ef4p1sbvkk9zxjxm0ez9npdr:validation:1
                briefing: briefing:ef4p1sbvkk9zxjxm0ez9npdr:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-10T19:47:30.043181Z"
                decision: revise
                reason: 'Captain adds scope to PR #873 before merge: run claude-haiku as a second claude-live leg on pull requests with continue-on-error: true, so haiku is sampled on every PR without blocking merges or the release gate. The validated restriction at ebf5dccfa is accepted as-is; the revision is the added workflow leg.'
                conn:
                    quote: do it as part of the pr with continue-on-error
                    source: captain chat, 2026-11-10
            - id: gate-attempt:ef4p1sbvkk9zxjxm0ez9npdr-validation-2
              briefing:
                id: briefing:ef4p1sbvkk9zxjxm0ez9npdr:validation:attempt-2:revision-1
                digest: sha256:54abbe4acb6bdabd09e4d7a26d200b99bac34a47f272e6265695243cfa31015a
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:ef4p1sbvkk9zxjxm0ez9npdr:validation:2
                briefing: briefing:ef4p1sbvkk9zxjxm0ez9npdr:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-11-11T05:08:55.243802Z"
                decision: approve
                reason: Cycle 2 re-review PASSED with no material finding; restriction proven in CI (0 haiku kills across three runs); pull-request haiku leg is non-blocking and dispatch paths unchanged. Polish findings declined; journey-ledger haiku records accepted for 0.28.0 with a follow-up.
                conn:
                    quote: Do that in background. Get this to release
                    source: 'captain chat, 2026-11-11, in reply to the FO recommendation to approve the #873 gate and merge the stack'
              application:
                target-stage: done
                state: consumed
archived: 2026-11-11T05:15:23Z
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

## Scope amendment

Captain decision, 2026-11-10: the same pull request also adds a non-blocking haiku leg to pull-request runs of Runtime Live E2E.

Outcome: a pull-request run executes `claude-live` on the selected sonnet model (blocking, check name unchanged) and on the registered `claude.haiku` model with `continue-on-error`, so a haiku failure is visible on its leg but does not fail the run conclusion or the release gate. Every `workflow_dispatch` path stays one blocking leg on the selected model, including `claude_model=haiku`. The haiku ID comes from the offline resolver's `claude_haiku` output, and the resolver refuses a pull-request registry without `claude.haiku`.

Added acceptance criterion:

5. The workflow defines two `claude-live` legs for `pull_request` (selected model blocking, `claude.haiku` continue-on-error) and one blocking leg for `workflow_dispatch`. Tested by `TestClaudeLiveKeepsCadenceBoundaryAndSelectedModel`, which expands both event branches of the `claude-live` model axis into legs and classifies each by the job's `continue-on-error` rule: pull request is `[selected blocking, haiku continue-on-error]`, dispatch is `[selected blocking]`. It fails when `continue-on-error` is dropped or set to `true`, when the dispatch branch gains the haiku leg, and when an include entry names a model. `TestLiveModelResolverEmitsPullRequestLegs` runs the resolver step and fails when the `claude.haiku` guard is removed.

Amended expected surface: 3 more files, net +107 (+120/−13).

| File | Change | +/− |
|---|---|---|
| `.github/workflows/runtime-live-e2e.yml` | model axis per event, `continue-on-error`, explicit job name, `claude_haiku` output, resolver guard, header comment | +14 / −4 |
| `internal/release/livemodels_test.go` | the two tests above | +105 / −8 |
| `docs/runtime-live-ci.md` | the pull-request lane bullet names the haiku leg | +1 / −1 |

## Stage Report: ideation

- DONE: Acceptance criteria name the measured end-value (a haiku CI live run with no no-stream-progress kill following a ScheduleWakeup call) and the deterministic test that pins the injected restriction in the built claude argument list
  AC-1 measures kills and ScheduleWakeup calls in the stacked PR's haiku job against a baseline of 7 kills in 3 attempts. AC-2 is the untagged `TestClaudeFrontDoorArgs`, with four falsifying edits named.
- DONE: The design names exactly which haiku journeys gain a claude-haiku xfail binding, with the run evidence for each, and states what happens to haiku bindings that passed in all four runs
  `rejection-flow` and `TestLiveBareReachable` are bound. `ac-value-reanchor` is not. The per-run table covers attempts 1, 2, 3 and the local run. The three never-red bindings stay; removing them is an open question for the captain.
- DONE: Expected surface (files, net LOC, tolerance) and the stacked-PR base ci/go-126-workflows are declared
  Net +30 (+46/−16) across 5 files, ±15 net. Base `ci/go-126-workflows` @ `eef2d1302` (PR #849). The registry doc diff is included.

### Summary

The design is small. It moves the `claude` argument list into an untagged helper that appends `--disallowedTools ScheduleWakeup` for haiku models only, adds `claude-haiku` xfail bindings to the two journeys still red with the restriction, and adds the two registry rows that contractlint requires. Re-pulling the CI logs showed that the local `-rerun3` evidence was a copy of attempt 2. Attempt 3 actually had 4 more kills (bare-reachable, rejection-flow, default-headless-gate-stop, recorded-gate-lifecycle), each after a 1200 s ScheduleWakeup. The body now uses the corrected 7-of-7 correlation and the corrected binding history: keep-moving-posture was red twice, and default-headless-gate-stop was killed once. The bare-reachable binding needs a non-Test gap function, because contractlint ignores gap literals inside Test functions.

## Stage Report: implementation

- DONE: The Claude live harness appends --disallowedTools ScheduleWakeup for haiku models only, through claudeFrontDoorArgs, and TestClaudeFrontDoorArgs fails under each falsifying edit named in the test plan
  Commit `ebf5dccfa`. `TestClaudeFrontDoorArgs` (untagged) checks that the pair is forwarded after `--`, and that `--model <model>` is present, for the haiku ID, `haiku`, the sonnet ID, `sonnet` and the opus ID. Each falsifying edit was applied and the test failed. Deleting the haiku branch fails both haiku cases. Making the append unconditional fails the sonnet and opus cases. Moving the pair before `--` fails the haiku cases. Renaming `claude.haiku` to `claude-small-5-5` in `live_models.txt` fails the registry case.
- DONE: rejection-flow and bare-reachable carry claude-haiku xfail bindings with matching registry rows, and go test ./..., go test ./... -race, and go vet -tags live ./internal/ensigncycle pass
  The bindings and rows are in place, and `TestRuntimeLiveRegistryMatchesBindings` passes. Deleting the `bareReachableGap` row makes it fail with "has no registry row". At `ebf5dccfa`: `go test ./...` passes, `go vet -tags live ./internal/ensigncycle` passes, and `go test ./... -race` passes (exit 0, 21 packages ok, no FAIL and no `DATA RACE`; cold build).
- DONE: The branch is based on origin/ci/go-126-workflows and the diff against it stays within the declared 5-file, +30 net LOC (±15) surface
  The worktree was fast-forwarded to `eef2d1302`. `git diff --numstat origin/ci/go-126-workflows`: 5 files, +46/−17, net +29. These are the 5 declared files, and `skills/` is unchanged.

### Summary

The `spacedock claude` argument list moved into the untagged `claudeFrontDoorArgs` helper. The helper forwards `--disallowedTools ScheduleWakeup` after `--` when the model name contains `haiku`, and `run` calls it. `rejection-flow` gained `liveXFail("claude-haiku", "")`. `TestLiveBareReachable` now grades through `finishLiveScenario` with the new `bareReachableGap`, and an unbound red still fails it. The registry has the two rows from the design. The first draft of the helper was +54 net, outside tolerance, so the helper and test were compacted to +29 before commit. `gofmt -l` reports an unrelated file on the base, `internal/release/runtime_live_evidence_workflow_test.go`. That file was left as is so that the diff stays within the 5 declared files.

## Stage Report: validation

- DONE: Each acceptance criterion has reproduced evidence: AC-2 and AC-3's deterministic half by rerunning the named tests and one falsifying edit each; AC-4 by the skills diff; AC-1 and AC-3's live half read from the CI haiku run the first officer names, or reported as pending with the exact check to run
  - AC-2, reproduced at `ebf5dccfa`: `TestClaudeFrontDoorArgs` passes. It asserts that `--disallowedTools ScheduleWakeup` follows `--` exactly for `claude-haiku-5-5` and `haiku`, and that `--model <model>` is forwarded in all five cases. In a throwaway export it failed under each of three edits: deleting the haiku branch (both haiku cases), making the append unconditional (`sonnet`, `claude-sonnet-5-5`, `claude-opus-5-5`), and moving the pair before `--` (both haiku cases).
  - AC-3 deterministic half, reproduced: `TestRuntimeLiveRegistryMatchesBindings` passes. It failed with "has no registry row" when either new registry row was deleted, and with "has no binding in internal/ensigncycle" when the `rejection-flow` haiku binding was deleted and its row kept.
  - AC-4, reproduced: `git diff origin/ci/go-126-workflows --stat -- skills/ cmd/` is empty. The sonnet and opus cases of AC-2 show their argument lists carry no restriction.
  - AC-1, measured on run 42882274740 (`ebf5dccfa`, job 130062411033 `claude-live (sonnet, claude-haiku-5-5, max, CI-E2E)`, conclusion success; artifact in `/tmp/sw-live-B3WiMQ`): 0 lines contain `no stream progress within` in the job log and in `live-e2e-detail.jsonl`, against the baseline of 7 kills. The artifact holds 22 `claude-stream.jsonl` files with 34 init events, all on `claude-haiku-5-5`. None lists `ScheduleWakeup` in its 23-tool init list, none contains a `ScheduleWakeup` tool_use, and the string appears in no stream file. This also shows `run` applies the helper.
  - Control from the same run: `TestLiveMergedTeamModeDispatch`, which builds its own argument list and is declared unchanged, still lists `ScheduleWakeup` in its init tools. It made no call and passed.
  - AC-3 live half, same run: `XFAIL claude-haiku/bare-reachable owner= observed=[bare-reachable-observables]` and `XPASS ALERT claude-haiku/rejection-flow owner= observed=[]`. The package passed in 306.7 s with 26 tests passed and 0 failed.
  - Other grade lines: XFAIL for `keep-moving-posture` and `default-headless-gate-stop`; XPASS ALERT for `withdrawn-gate-recovery` and break-glass selected-bare. No FAIL line.
- DONE: A semantic adversarial pass covers the haiku model-name match (aliases, registered ID, non-haiku names) and confirms a red bare-reachable or rejection-flow result still fails on claude-sonnet
  A scratch test in the throwaway export, built with `-tags live`, ran the real `claudeFrontDoorArgs`, `claudeLiveRole`, `bareReachableGap`, `selectLiveGap` and `liveScenarioGrade`.
  - Model names: the restriction is present for `claude-haiku-5-5`, `haiku` and a dated `claude-haiku-5-5-20261001`, and absent for `sonnet`, `claude-sonnet-5-5`, `opus`, `claude-opus-5-5`, `Haiku`, the empty string and `gpt-6-luna`. Every name `claudeLiveRole` accepts gets the restriction exactly when its role is `claude-haiku`. The pair is the last two arguments, so the variadic flag cannot swallow a later argument.
  - Grading: for both journeys a red grades `xfail` on `claude-haiku` and `fail` on `claude-sonnet`, `claude-opus`, `codex` and an empty target. Green grades `xpass` on haiku and `pass` elsewhere. An ungraded (infrastructure) error grades `fail` on every target, haiku included. The `rejection-flow` rows used a copy of its gaps literal, because the literal sits inside the Test function.
- DONE: The diff against origin/ci/go-126-workflows stays within the approved 5-file, +30 net LOC (±15) surface with no skill or product-binary change
  `git diff --numstat origin/ci/go-126-workflows...HEAD`: the 5 declared files, +46/−17, net +29, one commit. Four are `_test.go` files and one is `docs/runtime-live-ci-registry.md`. `go vet -tags live ./internal/ensigncycle` exits 0. `/tmp/race-forbid-schedulewakeup.log` shows 21 packages `ok` and `RACE EXIT 0`; the race suite was not rerun.

### Findings

1. Polish (evidence, no AC affected): with `SPACEDOCK_LIVE_MODEL=haiku`, the alias, the argument list is restricted but `claudeLiveRole` rejects the name. `TestLiveBareReachable` now stops with "want a registered Claude model ID or sonnet" before launch, as the shared journeys and break-glass already do; before this change it ran on any model string. CI passes the registered ID `claude-haiku-5-5`, so the supported path is unaffected. Recommended disposition: decline.

### Summary

Recommendation: PASSED. All four acceptance criteria have reproduced evidence at `ebf5dccfa`, and there is no material finding and no deferred risk. The haiku CI run had no no-progress kill and no `ScheduleWakeup` in any restricted stream, against a baseline of 7 kills, and it is a single run. The deterministic tests each failed under a falsifying edit, and an adversarial matrix over model names and grading targets held. The candidate worktree was not modified; falsifying edits ran in a throwaway export that has been removed. The run's `codex-live` job was still in progress when validation ended and is outside this task.

## Stage Report: implementation (cycle 2)

- DONE: Pull-request runs execute a second claude-live leg on the registered haiku model with continue-on-error, while the sonnet leg and every workflow_dispatch path are unchanged
  Commit `5c862c23c` on top of `ebf5dccfa` (not pushed). The `claude-live` matrix `model` axis is `fromJSON(format(...))` over `claude_selected` plus `claude_haiku` on `pull_request`, and `claude_selected` alone otherwise. A single model-free include entry carries cadence, effort and environment. `continue-on-error: ${{ matrix.model != needs.offline.outputs.claude_selected }}`, so only the extra leg is non-blocking, and an explicit job `name:` reproduces the include-only check name (`claude-live (pull-request, claude-sonnet-5-5, max, CI-E2E)`). Artifact names, config dirs and metrics dirs are per `matrix.model`, so the legs do not collide. `journey-delta-comment` still reads only the sonnet artifact. The release gate reads the run conclusion, which a continue-on-error failure does not turn red. `main` has no required status checks.
- DONE: The workflow-shape test pins the new matrix and fails when continue-on-error is dropped or the haiku leg leaks onto dispatch; go test ./... and go test ./... -race pass
  `TestClaudeLiveKeepsCadenceBoundaryAndSelectedModel` expands both event branches into legs through the offline outputs and classifies each by the continue-on-error rule. Five falsifying edits each failed: dropping `continue-on-error`, setting it to `true`, adding the haiku leg to the dispatch branch, adding a haiku include entry, and removing the resolver guard (`TestLiveModelResolverEmitsPullRequestLegs`, which runs the resolver step on the shipped registry and on one without `claude.haiku`). `go test ./...` and `go test ./... -race` exit 0. `actionlint` v1.7.7 (`-shellcheck= -pyflakes=`) reports no issue, before or after the change.
- DONE: The task body carries a Scope amendment section with the added acceptance criterion and amended surface, and the report states what remains unproven until the stack targets main
  `## Scope amendment` sits before the stage reports, with AC-5 and its tests, and the surface is 3 files, +120/−13, net +107. Unproven: no `pull_request` run has executed the two-leg matrix. The Runtime Live E2E `pull_request` trigger fires only for PRs into `main`, and #873 targets `ci/go-126-workflows`. Until the stack targets `main`, it is unproven that GitHub expands the expression axis into two legs, applies the job name, and reports the run as success when the haiku leg fails.

### Summary

The haiku leg is a base matrix axis chosen per event, not an include entry. GitHub turns an include entry that matches no existing combination into a new leg, so a haiku include entry would also run on dispatch. The resolver now refuses a pull-request registry without `claude.haiku`, which would otherwise produce a leg with an empty model. The test diff (+97 net) is most of the surface. `gofmt -l` still lists `internal/release/runtime_live_evidence_workflow_test.go`, which was already unformatted on the base and which this commit does not touch.

## Stage Report: validation (cycle 2)

- DONE: Re-review the corrected state: commit 5c862c23c on top of ebf5dccfa (the Scope amendment section and AC-5 in the task body); do not re-validate ebf5dccfa itself beyond confirming it is unchanged
  `5c862c23c^` is `ebf5dccfa`, whose diff against `eef2d1302` is still 5 files, +46/−17. `git diff --numstat ebf5dccfa 5c862c23c` is the 3 amended files, +120/−13, matching the amended surface. The worktree is clean at `5c862c23c`.
- DONE: Verify independently that on pull_request the claude-live matrix expands to the selected sonnet leg plus a registered-haiku leg, and on every workflow_dispatch cadence (including claude_model=haiku) to the selected leg only
  A scratch script ran the real "Resolve live lane models" step for each event, cadence and model key, then expanded the axis expression by hand. Pull request: `claude-sonnet-5-5` blocking, `claude-haiku-5-5` continue-on-error (cadence and model inputs are ignored). Dispatch `sonnet` with empty, `sonnet` or `haiku` key, and `opus-pre-release` with empty, `sonnet` or `haiku` key: one leg on the selected model. `sonnet` + `opus` and a `push` event are refused. `pi` selects no model, and the job's `if` skips it as before.
- DONE: Verify the haiku pull-request leg is continue-on-error and the selected leg stays blocking on every event; verify the sonnet job name and artifact names are unchanged and the two legs artifacts do not collide; verify journey-delta-comment still reads the sonnet artifact
  By reading the workflow: `continue-on-error` is true only when `matrix.model` differs from `claude_selected`, so the selected leg is blocking on every row above. The explicit `name:` gives `claude-live (pull-request, claude-sonnet-5-5, max, CI-E2E)`, the same order as the earlier include-only name seen on run 42882274740. The artifact name, artifact dir, metrics dir and `CLAUDE_CONFIG_DIR` all carry `matrix.model`. `journey-delta-comment` selects `runtime-live-e2e-claude-live-<claude_sonnet>` by exact name. No test pins these three (finding 1).
- DONE: Falsify the new tests yourself: make at least two breaking edits to the workflow (not the same ones blindly trusted from the implementation report) and confirm a test fails, then restore
  Ten edits in a throwaway export (removed), full `./internal/release/` each time. Seven failed a test: swapping the event branches, wiring the offline `claude_haiku` output to sonnet, comparing `continue-on-error` against `claude_sonnet`, listing the selected model twice on pull request, dropping the job `name:`, adding haiku to the dispatch branch (all `TestClaudeLiveKeepsCadenceBoundaryAndSelectedModel`), and inverting the resolver guard (`TestLiveModelResolverEmitsPullRequestLegs`). Three passed every test (finding 1).
- DONE: Run go test ./... and go test ./... -race, gofmt -l on touched Go files, and actionlint on the workflow if available
  In the worktree at `5c862c23c`: `go test ./...` exit 0, `go test ./... -race` exit 0, `go vet -tags live ./internal/ensigncycle` exit 0. `gofmt -l` on the five touched Go files prints nothing. actionlint was not run: it is not installed here, and fetching it was denied by the session's permission policy.
- DONE: State plainly what remains unproven until a real pull_request run executes the two-leg matrix, and give a PASSED or REJECTED recommendation with findings
  Unproven, because no `pull_request` run into `main` has used this workflow: that GitHub accepts the `fromJSON(format(...))` axis and combines it with the model-free include entry into exactly two legs; that the rendered sonnet check name is byte-identical to the old one; that a failed haiku leg leaves the run conclusion `success` and lets `journey-delta-comment` run; that the haiku leg waits for its own `CI-E2E` approval. The tests pin the expression text and the resolver's output, not GitHub's evaluation.

### Findings

1. Polish (evidence, AC-5 not affected): three properties the implementation report claims have no test. These edits passed the whole release package: keying the artifact name by `claude_selected` (both legs upload the same name), pointing `journey-delta-comment` at the `claude_haiku` artifact, and dropping `matrix.model` from the metrics dir. The shipped workflow is correct on all three by inspection. Recommended disposition: decline, or add one assertion each if the captain wants them pinned.
2. Needs decision (outcome, no AC affected): the release ledger step in `release.yml` downloads every artifact of every `success` run and folds in all `journey-metrics` JSON. Each pull-request run will now add haiku records to the journey-cost ledger, including records from a haiku leg that failed, since that run still concludes `success`. Record filenames and the ledger sort key include the model, so sonnet records are not overwritten. I did not trace how ledger consumers treat a second model per scenario. Promote to material if a consumer aggregates across models.
3. Polish: the pull-request lane now requests two `CI-E2E` approvals and spends a second roughly 7-minute haiku run per pull request. `docs/runtime-live-ci.md` states the second leg and its approval.

### Summary

Recommendation: PASSED. AC-5's two tests pass and failed under seven independent breaking edits, the resolver run across every event and cadence gives two legs on pull request and one blocking leg on every dispatch, and the full and race suites pass at `5c862c23c`. No finding is material. The real two-leg expansion, the check name and the run conclusion under a failed haiku leg stay unproven until a pull request into `main` runs this workflow, and actionlint was not rerun in this pass.
