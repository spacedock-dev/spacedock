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

Reduce live CI completion time by running a committed duration-sorted list of the existing tests through three slots. This task is the **bottom PR**; the other in-flight PRs stack above it.

## Problem

Run 34996910090 tested 380bd7bf34b3dd1de07a125d75035c593f5ff930, whose tree matches merged main 438053493838dc70c9478b3d991309d566783e85. Codex common journeys took 945.863s. Claude common journeys took 1334.022s, then three serial substrate tests added 520.586s. Committed hints can put longer work first and let those Claude proofs use the same three slots.

## Binding feedback and revision

Captain review `resolution:binding-1789503968390265000`, artifact digest `sha256:1d756f6a0d260bbdabb84bdf915477f6bc2c3b89f5b66ae6b9e430f6f9cab935`, returned **REVISE**:

> this feels a lot of code. are you adding unnecessary tests to the scheduler? just run the tests in sorted order.

Before: +330 net lines across 8 files, with dedicated queue/capacity/lifecycle controls, a frozen timing-replay fixture, extended callable/priority checks, and fresh failure/timeout exercises.

After: **+120 net lines across the same 8 files**, one literal hint/function list, one small ordering check, and a short three-worker loop that calls the existing test bodies. No frozen replay fixture, new failure/timeout/cancellation suite, parallel-ancestor AST guard, generic scheduler API, generic AST interpreter, process launcher or new runtime flag. Existing registry and gotestsum tests retain their jobs. The prior spike evidence remains useful; it does not need a permanent copy of every experiment.

## Smallest supported design

Add the exact live selector `^TestLiveScheduled$` for Claude and Codex. A literal table has the 17 existing common test functions and their rounded Claude/Codex hints. For Claude, append the three existing substrate functions. Keep the whole break-glass test together, including both sequential variants. Pi is unchanged.

Sort the selected list by descending hint and exact test-name ties. Use three native Go worker subtests, each calling `t.Parallel()` once and taking the next item from the shared sorted list before a synchronous `t.Run`. Protect the next index with a mutex. This is a short loop local to this test suite, not a reusable scheduler. The common runner suppresses its nested `t.Parallel()` only beneath `TestLiveScheduled/`; standalone exported tests retain their existing behavior. Scheduled test bodies and their variants must remain synchronous beneath each worker.

Why a loop remains: an isolated actual-Go experiment registered 12 parallel subtests in order `0..11` with `-parallel 3`. The three runs started in orders:

- `0,5,6,4,3,2,1,9,11,10,8,7`
- `0,5,6,1,4,3,9,11,10,8,2,7`
- `0,5,6,4,3,2,1,9,11,10,8,7`

Sorted registration therefore does not produce useful priority order. The retained native experiment is `/tmp/spacedock-native-order-spike/order_test.go`; reproduction is `go test -v -parallel 3 -count=3 /tmp/spacedock-native-order-spike/order_test.go`. It is a throwaway probe, not a new repository test. The already committed worker spike `72718be3f` shows the small worker loop works. No need to invent another scheduling mechanism or serialize all tests. The contract is sorted admission to three slots; CPU-level host startup order can differ.

Give merged dispatch its own `<CLAUDE_CONFIG_DIR>/merged-team-mode` child. Common, bare and break-glass scenarios already use their own scenario children. Keep their separate HOME/workflow roots and existing artifact paths. No global auth/config change.

Claude uses one gotestsum invocation for all 20 tests, writing `live-e2e-detail.jsonl`; remove its later serial substrate invocation. Keep `always()` artifact upload, streams, metrics, config projects and candidate provenance. Its existing 90-minute suite backstop now covers all 20 tests; the separate 20-minute substrate-package backstop disappears. Codex keeps its 40-minute timeout and detail filename. Do not pass `-failfast` to either scheduled suite. Existing runner assertions, liveness and cleanup stay unchanged.

Keep the original exported tests for targeted diagnosis. The three-slot/exactly-once promise applies to the canonical anchored scheduled selector. Broad selectors such as `-run TestLive` would select both the wrapper and original tests, duplicating work and possibly artifact names; document the exact selector instead of adding a runtime flag or suppression layer.

## What remains, and why

- **Literal hints and ordering helper (AC-1):** needed to order differently per runtime. One focused test supplies a few unsorted jobs, including an equal-hint pair, and checks the exact independently expected order. Reversing the comparison or tie order fails it. No captured-duration fixture or timing threshold in CI.
- **Three-worker loop and nested-parallel suppression (AC-2):** needed because native sorted registration does not control admission. Reuse the retained spike's evidence for three occupied slots, synchronous completion and continued execution after a failing child; do not copy its separate mutation/failure/timeout experiments into a new test suite.
- **Merged config child (AC-2):** needed before overlap; distinct working directories alone do not isolate the shared Claude config root. Actual host overlap remains the acceptance proof.
- **Small adaptation of existing registry reconciliation (AC-3):** accept the new canonical selector and orchestration entry, and reconcile the literal function references with the current 17 common/3 Claude proof declarations. Check each row's name matches its callable and detect omission/duplication while retaining current assertion/fixture/TODO/XFAIL checks. Do this inside the existing parser walk; no second registry or generalized parser. Compiler/runtime selection and the literal append make Claude-only membership reviewable. Do not add a separate speculative hint-schema or future parallel-call checker.
- **Existing gotestsum controls (AC-3):** keep their ownership of summary, archive and exit-code behavior. No additional permanent failure/cancellation controls are needed for unchanged gotestsum/native `t.Run` behavior.

## Risk evidence and remaining proof

Retained local spike commit `72718be3f` at `/tmp/spacedock-schedule-ideation` is an experiment, not a product branch. It observed peak 3, all 8 fake jobs complete, 20 race repetitions passing, and all 8 complete with exit 1 after injected `t.Fatal`. The reverse-order and fourth-worker mutations were rejected. Gotestsum failure and timeout runs retained JSON detail with nonzero exits. Existing environment and `TestLiveCIStep*` controls passed. These are already-run evidence, not instructions to build a parallel proof suite.

Three concurrent shell subprocesses retained distinct HOME, config/projects and workflow markers. This models the intended path separation and is **not** real Claude isolation proof. The attempted actual merged/bare/shallow-boot overlap skipped all three tests before launch: no API-key/OAuth environment and the OS denied access to `~/.claude/benchmark-token`. Do not bypass the denial or treat skipped tests as passing evidence. Before implementation acceptance, run one bounded real overlap of merged dispatch, break-glass (both sequential variants) and common shallow boot; require overlapping host intervals, distinct actual config/workflow/session paths, durable assertions and retained artifacts.

The measurements below are design evidence for the hints, not a new test-data file or a future remote dependency. Round whole-test terminal durations to the nearest 10 seconds; do not add nested variant elapsed again. Sonnet measurements seed both Claude cadences. Update hints manually if sustained timings change useful ordering. No automatic updater or performance gate.

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

**AC-1 — Longer hinted tests are admitted first from a committed per-runtime list.** The single ordering check proves descending hints with lexical ties against an independent expected list. The retained measurements predict 817.88s Codex versus 945.863s observed, and 1363.72s pooled Claude versus 1854.608s observed; these quantify the intended value without promising provider wallclock gains or adding a timing fixture. No remote reads or automatic hint writes occur.

**AC-2 — Claude common and substrate tests use one three-slot queue with isolated host state.** The short worker loop retains a slot through synchronous `t.Run` and cleanup, including sequential variants; `-parallel 1` can reduce concurrency. Use the existing bounded spike evidence for queue behavior and complete the one real overlap described above before acceptance. A fixture pass or auth skip cannot close this criterion. Existing exported direct tests remain usable with exact selectors.

**AC-3 — The canonical scheduled invocation preserves all existing tests, assertions and evidence.** Claude runs the existing 20 functions and Codex the existing 17; Pi is unchanged. Existing registry reconciliation checks the literal scheduled callables and canonical selectors. Existing gotestsum tests own clean summary, JSON detail and failure exit behavior; ordinary test failure still allows queued jobs to finish, whereas timeout/cancellation may leave unrun jobs and retain only available detail. No new cancellation or process-cleanup guarantee is introduced.

## Expected surface and tolerance

Estimate net LOC change: **+120, across 8 files**. Estimated insertions +155, deletions -35. Tolerance: **+30 net lines, no additional files**; ceiling **+150 net / 8 files**. This is a 64% reduction from the original +330 estimate. If the concise adaptation cannot fit, report the concrete reason instead of adding scaffolding.

1. `internal/ensigncycle/live_schedule_test.go` (new): small row type/order helper and one focused independent-order test.
2. `internal/ensigncycle/scheduled_live_test.go` (new): literal hint/function table, runtime selection and short three-worker loop.
3. `internal/ensigncycle/shared_live_runner_test.go`: scheduled ancestry guard.
4. `internal/ensigncycle/merged_team_mode_live_test.go`: merged config child.
5. `internal/contractlint/live_registry_reconciliation_test.go`: concise callable/selector adaptation inside existing reconciliation.
6. `.github/workflows/runtime-live-e2e.yml`: combined Claude command and Codex selector; remove obsolete serial step/artifact entry.
7. `docs/runtime-live-ci.md`: canonical commands, hints and combined evidence location.
8. `docs/runtime-live-ci-registry.md`: orchestration entry; existing journey/proof identities retained.

Allowed semantics: test selector/nested event names, admission order, Claude overlap, merged config location, combined Claude JSON detail, and the explicit timeout consolidation. No spacedock CLI grammar, durable state format, authority, assertions, Pi behavior, auth policy or provider settings change. Staffing remains one implementation owner and one independent validator; one additional reviewer only for a distinct uncovered claim. Scheduling is the bottom PR, with other in-flight PRs above it. No implementation before the revised binding gate; no push or CI.

## Minimal proof plan

1. Before implementation, add the one ordering check with a small independent expected sequence and equal-hint pair. No timing sleeps, replay fixture or scheduler test framework.
2. Adapt and run the existing registry reconciliation and gotestsum proof owners. Preserve exact callable coverage and use anchored selectors. The retained worker/failure spike is sufficient design evidence for the few lines of native worker glue; review that glue directly instead of introducing new permanent lifecycle/cancellation controls.
3. Complete one bounded local real Claude overlap (merged, break-glass, shallow boot; maximum three occupied tests, 15-minute local cap) with readable authorized credentials. Record actual intervals/paths and original assertions; skipped execution remains outstanding.
4. Implementation still owes repo-required `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. No full suite is being rerun for this ideation revision. Independent validation reviews the small diff and existing checks; full live CI waits for the ready stack.

## Proposed documentation diff

In `docs/runtime-live-ci.md`, replace “Runtime-specific substrate proofs stay separate because they verify host boundaries rather than workflow semantics” with:

> Claude's three substrate proofs retain their separate assertions and share the common journeys' three-slot queue.

Replace the common scheduling paragraph and separate Claude substrate command with:

> Claude and Codex run longer hinted tests first from a committed list, using three slots and continuing after a test failure. Claude includes its 17 common journeys and three substrate tests. Use the exact scheduled selector below; broad selectors such as `TestLive` also select the original tests and duplicate work. Exact original test names remain available for diagnosis.
>
> `SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveScheduled$' -parallel 3 ./internal/ensigncycle -v`

Replace the canonical Codex command with:

> `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveScheduled$' -parallel 3 ./internal/ensigncycle -v`

Add:

> Update the rounded hints in `internal/ensigncycle/scheduled_live_test.go` manually when sustained whole-test timings change useful ordering. Use the terminal `TestLiveScheduled/slot-N/<exported-name>` event without adding its nested variants. Hints are priorities, not timeout budgets. Claude keeps all test events in `live-e2e-detail.jsonl` under one 90-minute suite backstop; streams, metrics and config projects remain archived. Scheduling reads no remote history.

In `docs/runtime-live-ci-registry.md`, add:

> ### `TestLiveScheduled`
>
> Claude/Codex scheduling wrapper for the existing registered test functions; adds no journey or substrate assertion.

### Feedback Cycles

Cycle 1: binding REVISE recorded above. The prior independent review accepted the architecture but recommended extra proof controls; the captain's subsequent binding feedback supersedes those additions. The selected revision keeps only the small ordering check and existing proof owners. Await another binding review before implementation.

## Stage Report: ideation

- DONE: Prove the smallest scheduler orders by committed hints and enforces one three-slot cap without remote state.
  Local spike 72718be3f exercises native Go subtests; reversed-order and fourth-worker mutants are rejected; rounded retained-duration replay is 817.88s Codex / 1363.72s Claude.
- FAILED: Establish actual substrate/common isolation and preserve test coverage, failure propagation, and artifacts.
  Environment/subprocess isolation and gotestsum failure/timeout evidence pass, but all three real-host tests skipped for unavailable readable credentials; actual host overlap remains mandatory before implementation acceptance.
- DONE: Record a minimal design, exact surface/tolerance, and targeted local proof plan for a separate stack PR.
  +330 net LOC (+370/-40), 8 files; tolerance +100 net / +2 files; original tests retained, canonical doc edits and bounded proof plan recorded above.

### Summary

Selected a three-worker Go test wrapper with a single duration-ordered queue, reusing all existing test bodies and gotestsum evidence. Deterministic scheduling and failure behavior were exercised in an isolated committed spike; live host overlap remains an explicit unmet proof because the available credential path was denied by the OS. No product branch, global auth/config, push or CI mutation was performed.


## Stage Report: ideation (cycle 2)

- DONE: Prove the smallest scheduler orders by committed hints and enforces one three-slot cap without remote state.
  Native sorted registration was exercised in three tiny runs and reordered jobs; retained worker spike 72718be3f proves the minimal three-slot alternative. Only one permanent order check is proposed.
- FAILED: Establish actual substrate/common isolation and preserve test coverage, failure propagation, and artifacts.
  Prior environment/failure/artifact evidence stands; actual Claude overlap remains unproven because readable credentials are unavailable, and remains required before acceptance.
- DONE: Record a minimal design, exact surface/tolerance, and targeted local proof plan for a separate stack PR.
  Binding feedback recorded; revised +120 net (+155/-35), 8 files, ceiling +150 net / 8 files; scheduling is the bottom PR and duplicate proof scaffolding is removed.

### Summary

Reduced the design from +330 to +120 estimated net lines by removing the broad new proof suite and retaining existing registry/gotestsum ownership. Sorted native registration demonstrably reorders work, so the design retains only a short three-worker loop to consume the sorted list. No implementation, full-suite rerun, push or CI action occurred; another binding review is required.
