---
title: Schedule live journeys with committed duration hints
status: ideation
source: Captain request 2026-09-15; CI run 34996910090
started: 2026-09-15T18:21:56Z
completed:
verdict:
score: 0.8
worktree:
issue:
pr:
mod-block:
id: pytyzge5v85mcy8c02t1khqv
---

Reduce live CI completion time with committed per-runtime duration hints and one bounded queue for common journeys and Claude substrate proofs.

## Problem

CI run 34996910090 tested 380bd7bf34b3dd1de07a125d75035c593f5ff930 (PR #798). Its tree is identical to main 438053493838dc70c9478b3d991309d566783e85 after PRs #794–#798 merged. Common journeys use three slots without duration ordering; Claude then runs three substrate tests serially. Observed package elapsed: Codex common 945.863s; Claude common 1334.022s plus substrate 520.586s = 1854.608s.

## Selected design

Add one live test entry point, `TestLiveScheduled`, for Claude and Codex. It selects a small literal table of existing test functions with their names and rounded second hints. Claude selects the 17 common journeys plus merged dispatch, bare dispatch and the whole break-glass test; Codex selects the 17 common journeys. Both break-glass variants and all existing sequential journey variants stay inside their original test bodies. Pi keeps its current selector and concurrency policy.

Sort selected jobs by descending committed hint, then exact exported test name ascending. Three worker subtests call `t.Parallel()` once each, take the next job under one mutex, and run its existing function in a synchronous child `t.Run`. This uses Go's test lifecycle and native failure propagation. Workers continue after a failing child; CI must not pass `-failfast`. No separate process launcher, queue service, dependency, or remote history is needed.

The guarantee is deterministic queue admission order, not nanosecond ordering of host process startup across three CPUs. Emit each admitted ordinal and exact job name in the existing test log. A freed slot takes the next longest remaining job. The suite has at most three occupied worker slots even with `-parallel 8`; `-parallel 1` may reduce concurrency. There is no second substrate pool.

`liveJourney` suppresses its existing `t.Parallel()` only when `t.Name()` begins with `TestLiveScheduled/`, so a worker owns the slot until that whole job and its subtests finish. Existing exported direct selectors keep their behavior. This one ancestry check avoids rewriting 17 wrappers or introducing an environment flag that could alter unrelated runs. Register the new suite as orchestration, while retaining every existing journey/proof registration and assertion.

Common Claude scenarios, bare dispatch, and each break-glass variant already put config under `<CLAUDE_CONFIG_DIR>/<scenario-name>`. Merged dispatch currently uses the parent config root directly. Give it `<CLAUDE_CONFIG_DIR>/merged-team-mode`, including its metadata reader and reconcile subprocess. Every test already receives a separate HOME and workflow root; existing scenario artifact paths and Codex setup paths remain in use. Do not mutate process-wide environment or global auth/config in scheduled jobs.

Claude's existing gotestsum common step becomes the single 20-job scheduled run. Remove the later serial substrate step. Keep the 90-minute suite backstop, all existing per-runner liveness budgets, the same clean summary output, and `always()` artifact upload. The separate 20-minute substrate-package timeout disappears because the proofs now belong to the combined suite; this is an explicit runtime change. Its event content moves into `live-e2e-detail.jsonl`; stream files, metadata, metrics, config projects and candidate provenance remain retained. Codex keeps its 40-minute backstop and current detail filename with the new selector.

Hints are positive integer seconds in the literal scheduling table, keyed by runtime and exact exported test name. Round initial whole-test Go elapsed times to the nearest 10 seconds; do not add nested variant elapsed again. Sonnet measurements seed both Claude cadences; these are priorities, not timeout budgets. Missing, duplicate or nonpositive entries fail reconciliation before live execution. A new journey requires an explicit conservative hint in the same change. Update hints manually only when sustained whole-test timings would change useful ordering. Merged dispatch lacks a journey metric today, so use Go JSON terminal events for it. Scheduled whole-test event names are nested under `TestLiveScheduled/slot-N/<exported-name>`; exclude their deeper variant events when measuring.

## Mechanism necessity

- The literal hint table serves AC-1. Source declaration order cannot differ by runtime, and Go's `t.Parallel()` resumption order is not a priority contract. A parsed JSON manifest would add a format and parser without improving this small checked-in table.
- Three workers plus one sorted queue serve AC-1/AC-2. Three static longest-first lanes are smaller conceptually but strand capacity when observed durations differ from hints. One shared dynamic queue reuses a freed slot. A semaphore around arbitrary parallel test starts limits concurrency but does not preserve priority and can deadlock if higher-priority tests cannot obtain Go's parallel slots.
- The ancestry check serves AC-2. Reusing unmodified `t.Parallel()` children would return worker capacity early; rewriting each test into a second callable body expands the surface. Direct tests stay available for diagnosis.
- The merged config child serves AC-2/AC-3. Merely retaining distinct working directories leaves shared Claude settings and metadata under the same config root.
- Registry-to-scheduled-function equality serves AC-3. Hint rows alone cannot prove the scheduled suite still covers all registered tests. Extend the current AST reconciliation rather than create another registry.

## Risk evidence

The isolated spike is committed locally at `72718be3f` in `/tmp/spacedock-schedule-ideation`, based on candidate `380bd7bf`; no product branch changed. `SCHEDULE-SPIKE.md` in that commit records runnable commands. This is a local throwaway commit, not a pushed prerequisite; implementation must copy its behavioral proof into ordinary repository tests.

- Actual Go worker/subtest exercise: peak 3 and all 8 jobs completed; race run repeated 20 times passed. An injected `t.Fatal` in job 2 returned exit 1 while all 8 jobs completed.
- Hint-order mutant admitted `7,6,5,4,3,2,1,0` instead of the fixed independent expectation `0,1,2,3,4,5,6,7`; the admission-order oracle rejected it. Four-worker mutant with `-parallel 8` failed the peak=3 assertion. Restore was completed before the spike commit.
- Actual gotestsum failure run retained 74 JSON events, 8 terminal job events and package failure; exit 1. Forced 200ms suite timeout retained 128 events and package failure; exit 1. This proves retained partial evidence, not completion or descendant cleanup after cancellation.
- Three simultaneous shell subprocesses using the existing environment helpers retained distinct HOME, config/projects and workflow markers under race. This proves filesystem/environment separation only; the test models the existing scenario suffix and proposed merged suffix, so it does not independently prove actual runner selection.
- Existing `TestDecideClaudeEnv`, `TestResolveClaudeConfigDir`, `TestIsolatedClaudeEnv*`, and `TestLiveCIStep*` controls passed.
- Actual `TestLiveScheduled` overlap attempt selected merged dispatch, bare dispatch and common shallow boot. All three SKIPPED before host launch: API-key/OAuth environment absent and the OS denied access to `~/.claude/benchmark-token`. No access bypass, global auth change or provider claim was made. **Actual Claude overlap remains unproven and is required before implementation acceptance.**

Retained Go JSON was analyzed locally from `/tmp/ci-schedule-34996910090`; these paths are evidence locations, not scheduler inputs. Frozen whole-test seconds and independently rounded priorities are recorded below so the deterministic regression can be recreated without downloads. Common rows abbreviate the `TestLiveCommon` prefix; substrate rows use their full names.

| Test | Codex observed / hint | Claude observed / hint |
|---|---:|---:|
| ACValueReanchor | 102.61 / 100 | 361.96 / 360 |
| AutoContinueAfterImplementation | 340.06 / 340 | 439.65 / 440 |
| DefaultHeadlessGateStop | 283.77 / 280 | 244.41 / 240 |
| FeedbackThreeCycleEscalation | 70.43 / 70 | 370.59 / 370 |
| Filing | 43.43 / 40 | 74.55 / 70 |
| FullEnsignCycle | 106.05 / 110 | 169.67 / 170 |
| GateGuardrail | 94.89 / 90 | 86.00 / 90 |
| KeepMovingPosture | 266.03 / 270 | 337.62 / 340 |
| MergeHookGuardrail | 46.22 / 50 | 30.85 / 30 |
| OwnedConflictOwnerHandoff | 134.34 / 130 | 302.56 / 300 |
| RecordedGateLifecycle | 158.71 / 160 | 192.01 / 190 |
| RejectionFlow | 335.78 / 340 | 425.09 / 430 |
| SelfEvidenceMergeTriage | 68.13 / 70 | 118.67 / 120 |
| ShallowBoot | 24.18 / 20 | 27.75 / 30 |
| SmallestSufficientMechanism | 241.45 / 240 | 245.81 / 250 |
| WithdrawnGateRecovery | 88.72 / 90 | 96.46 / 100 |
| ZeroDiscovery | 26.23 / 30 | 17.51 / 20 |
| TestLiveBareReachable | — | 111.12 / 110 |
| TestLiveBreakGlassShimRecovery | — | 274.27 / 270 |
| TestLiveMergedTeamModeDispatch | — | 135.19 / 140 |

Rounded-priority replay gives Codex 817.88s (127.983s below the 945.863s observed baseline) and pooled Claude 1363.72s (490.888s below 1854.608s). Reversing priorities gives 929.68s and 1500.31s respectively. These are retrospective fixed-duration predictions, not live speed guarantees.

## Acceptance criteria

**AC-1 — Committed priorities reduce predicted completion time and release freed capacity without external state.** The frozen independent durations above complete within 830s for Codex and 1400s for pooled Claude at three slots, against observed baselines of 945.863s and 1854.608s. Every selected job is admitted exactly once in descending hint order with lexical ties; a short actual job releases capacity while another long job is still active. Proof: deterministic scheduler/subtest exercise plus independent duration replay. Reversing the comparator fails both thresholds; changing only the tie comparator fails the equal-hint case. No wallclock timing gate or remote lookup enters CI.

**AC-2 — Common journeys and Claude substrate proofs share one cap of three occupied tests with isolated host state.** Proof: barrier-controlled jobs with `-parallel 8` never exceed three, cover mixed common/substrate labels, and hold capacity through sequential variants and cleanup. Separately, one bounded local real Claude overlap runs merged dispatch, break-glass (both sequential variants) and common shallow boot with shared archive root and distinct config/workflow/HOME state; inspect transcript session IDs, durable fixture output, and artifact paths. A fourth worker, early slot release, or removal of the merged suffix must fail its corresponding check. Fixture evidence and auth skips cannot satisfy the real-host claim.

**AC-3 — Scheduled CI retains the complete registered coverage and failure evidence.** Claude schedules 20 existing exported tests, Codex 17; Pi retains its current surface. Original assertions, TODO/XFAIL classification, all variant counts, runner cleanup and metrics stay intact. Proof: existing AST registry reconciliation extended to compare scheduled function references with registered common/proof functions, plus actual gotestsum runs with injected failure and timeout. An omitted/duplicate callable, swallowed failure, skipped later job, or missing detail archive fails its corresponding check. A canceled/timed-out suite may leave jobs unrun but must report failure and retain available evidence; the task must not claim canceled work completed.

## Expected surface and tolerance

Estimate net LOC change: +330, across 8 files. Estimated insertions +370, deletions -40. Tolerance: at most +100 additional net lines and 2 additional files (ceiling +430 net / 10 files); fewer lines/files are welcome. A new scheduler service, process supervisor, runtime flag or production CLI package exceeds the design regardless of line count.

1. `internal/ensigncycle/live_schedule_test.go` (new): small queue helper plus offline behavioral tests and frozen baseline literals.
2. `internal/ensigncycle/scheduled_live_test.go` (new): literal hint/function table and suite entry point.
3. `internal/ensigncycle/shared_live_runner_test.go`: scheduled ancestry guard only.
4. `internal/ensigncycle/merged_team_mode_live_test.go`: merged config child.
5. `internal/contractlint/live_registry_reconciliation_test.go`: scheduled membership and selector reconciliation.
6. `.github/workflows/runtime-live-e2e.yml`: Claude combined run, Codex selector, artifact entry/comment updates.
7. `docs/runtime-live-ci.md`: canonical commands, concurrency, hint maintenance and combined evidence location.
8. `docs/runtime-live-ci-registry.md`: orchestration entry, retaining all journey/proof identities.

Allowed observable changes: internal Go test selector and nested test event names, test admission order, Claude substrate/common overlap, merged config location, combined Claude JSON detail file, and replacement of the separate substrate 20-minute package timeout by the existing 90-minute combined-suite timeout. No spacedock command grammar, durable workflow format, FO/ensign authority, assertion semantics, Pi behavior, auth policy or provider settings change.

Expected staffing: one ideation/implementation owner and one independent validator/auditor; tolerance one extra reviewer only for a distinct uncovered claim. Deliver as a separate PR above merged main or the forthcoming Claude repair; no old stack branch mutation. No push or CI until the captain's stack-ready authorization covers it.

## Test plan

The current proof owners are the 17 exported live journeys and three Claude substrate tests for runtime behavior; `TestRuntimeLiveRegistryReconciliation` for coverage; and `internal/release/cilog_clean_output_test.go` for gotestsum summary/detail/exit semantics. The queue itself has no existing proof owner, so add focused tests before implementation.

1. Write queue/admission and capacity tests first (subsecond deterministic jobs; no provider). A reverse comparator, unstable tie, fourth slot, duplicate pop, early slot release or fixed partition with a slow lane must produce an observable failure. Use barriers instead of elapsed-time races for the cap/continued-admission assertions. The frozen replay is separate from real job timing.
2. Extend existing registry reconciliation before changing selectors. Compare real callable references, exact names, counts, positive hints and per-runtime membership, retaining all existing AST journey/assertion checks. Deleting a row or referencing the wrong function must fail; parsing only a label is insufficient. Keep this within the current test file.
3. Use the existing gotestsum proof owner plus one queue failure/timeout subprocess exercise (seconds). Keep later jobs running after `t.Fatal`, assert every expected terminal event, nonzero exit, and retained failure/partial-timeout detail. Reuse standard Go/runner teardown; do not build a cancellation framework. Existing cleanup assertions remain; if actual overlap exposes descendant leaks, report a scope decision rather than silently adding a supervisor.
4. Finish the required bounded local Claude overlap once authorized readable credentials exist (three occupied tests maximum; 15-minute local cap). Use the real merged, break-glass and shallow-boot bodies with all assertions; no provider/model mutation and no full-suite burn. Record skipped runs distinctly. Any failure must be attributed before treating it as scheduling evidence.
5. Run focused changed-package checks, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` before implementation completion. Run an independent adversarial audit of registry/CI wiring. Reserve full Claude/Codex stack-tip CI for the ready stack.

## Proposed documentation diff

In `docs/runtime-live-ci.md`, replace:

> Runtime-specific substrate proofs stay separate because they verify host boundaries rather than workflow semantics.

with:

> Runtime-specific substrate proofs retain separate assertions. Claude schedules its three substrate tests alongside the 17 common journeys in one queue.

Replace the local paragraph beginning “Run the common journeys by selecting one transport” through the separate Claude substrate command with:

> Claude and Codex use committed whole-test duration hints to admit longer tests first. Each scheduled suite has one three-slot queue and continues after a test failure. Claude runs the 17 common journeys and three substrate tests together; Codex runs the 17 common journeys. The hints are priorities, not timeout budgets. The Claude suite has a 90-minute backstop. Existing exported test names remain available for targeted diagnosis.
>
> `SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveScheduled$' -parallel 3 ./internal/ensigncycle -v`

Replace the canonical Codex command with:

> `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveScheduled$' -parallel 3 ./internal/ensigncycle -v`

Add after those commands:

> Update the rounded second hints in `internal/ensigncycle/scheduled_live_test.go` manually when sustained whole-test timings change useful ordering. Use the terminal event for `TestLiveScheduled/slot-N/<exported-test-name>`; do not add its nested variant times again. Claude CI keeps common and substrate events in `live-e2e-detail.jsonl`. Runtime streams, journey metrics and per-scenario config projects remain in the uploaded artifact. Scheduling does not read remote metrics or rewrite hints.

In `docs/runtime-live-ci-registry.md`, add the orchestration entry without changing the existing 20 identities:

> ### `TestLiveScheduled`
>
> Scheduling entry point for Claude and Codex. It invokes the existing registered tests once each, using committed duration priorities and one three-slot queue. It adds no journey or substrate assertion.

### Feedback Cycles

## Stage Report: ideation

- DONE: Prove the smallest scheduler orders by committed hints and enforces one three-slot cap without remote state.
  Local spike 72718be3f exercises native Go subtests; reversed-order and fourth-worker mutants are rejected; rounded retained-duration replay is 817.88s Codex / 1363.72s Claude.
- FAILED: Establish actual substrate/common isolation and preserve test coverage, failure propagation, and artifacts.
  Environment/subprocess isolation and gotestsum failure/timeout evidence pass, but all three real-host tests skipped for unavailable readable credentials; actual host overlap remains mandatory before implementation acceptance.
- DONE: Record a minimal design, exact surface/tolerance, and targeted local proof plan for a separate stack PR.
  +330 net LOC (+370/-40), 8 files; tolerance +100 net / +2 files; original tests retained, canonical doc edits and bounded proof plan recorded above.

### Summary

Selected a three-worker Go test wrapper with a single duration-ordered queue, reusing all existing test bodies and gotestsum evidence. Deterministic scheduling and failure behavior were exercised in an isolated committed spike; live host overlap remains an explicit unmet proof because the available credential path was denied by the OS. No product branch, global auth/config, push or CI mutation was performed.
