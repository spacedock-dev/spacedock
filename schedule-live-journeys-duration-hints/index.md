---
title: Schedule live journeys with committed duration hints
status: validation
source: Captain request 2026-09-15; CI run 34996910090
started: 2026-09-15T18:21:56Z
completed:
verdict:
score: 0.8
worktree: .worktrees/spacedock-ensign-schedule-live-journeys-duration-hints
issue:
pr: "#800"
mod-block: merge:pr-merge
id: pytyzge5v85mcy8c02t1khqv
gates:
    version: 1
    records:
        - id: gate:pytyzge5v85mcy8c02t1khqv:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:pytyzge5v85mcy8c02t1khqv-ideation-1
              briefing:
                id: briefing:pytyzge5v85mcy8c02t1khqv:ideation:attempt-1:revision-1
                digest: sha256:2e92e016dc3de77e21db5d699185aa479418ee201509cf38df031db3ba74d947
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:pytyzge5v85mcy8c02t1khqv:ideation:1
                briefing: briefing:pytyzge5v85mcy8c02t1khqv:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T22:40:16.37176Z"
                decision: approve
                reason: Binding Subspace resolution binding-1789511305813868000 approved reduced design, cap +150 net/8 files, scheduling as bottom PR; actual Claude overlap remains required before acceptance. Reviewed approval artifact sha256:c5fffcc6f164abd78be095350372d91b0563a628775646f506e7c29c12d339cc.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:pytyzge5v85mcy8c02t1khqv:validation
          stage: validation
          attempts:
            - id: gate-attempt:pytyzge5v85mcy8c02t1khqv-validation-1
              briefing:
                id: briefing:pytyzge5v85mcy8c02t1khqv:validation:attempt-1:revision-1
                digest: sha256:8bab61d8e05fc594a993ff91cfe44702447171af4a99f98c29a242ed827fba98
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:pytyzge5v85mcy8c02t1khqv:validation:1
                briefing: briefing:pytyzge5v85mcy8c02t1khqv:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-16T04:39:58.914631Z"
                decision: approve
                reason: Captain approved validation in Subspace resolution:binding-1789533522567170000; local validation accepted and real Claude overlap/isolation remains required at stack-tip CI before delivery.
              application:
                target-stage: done
                state: pending
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

Sorted registration therefore does not produce useful priority order. The native experiment source and outcomes are committed in `SCHEDULE-SPIKE.md` at local spike commit `58a777779`; its executed copy is `/tmp/spacedock-native-order-spike/order_test.go`. Reproduction is `go test -v -parallel 3 -count=3 /tmp/spacedock-native-order-spike/order_test.go`. It is a throwaway probe, not a new repository test. The already committed worker spike `72718be3f` shows the small worker loop works. No need to invent another scheduling mechanism or serialize all tests. The contract is sorted admission to three slots; CPU-level host startup order can differ.

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

Three concurrent shell subprocesses retained distinct HOME, config/projects and workflow markers. This models the intended path separation and is **not** real Claude isolation proof. The attempted actual merged/bare/shallow-boot overlap skipped all three tests before launch: no API-key/OAuth environment and the OS denied access to `~/.claude/benchmark-token`. Do not bypass the denial or treat skipped tests as passing evidence. Captain amendment of 2026-09-15 ("that can be in ci") moves the required real Claude overlap proof to the single full stack-tip CI run. Local acceptance covers deterministic verification and independent review; final delivery requires the existing full `TestLiveScheduled` invocation across all 20 Claude functions to show actual common/substrate overlap under the three-slot cap, overlapping host intervals, distinct actual config/workflow/session paths, passing original durable assertions (including both sequential break-glass variants) and retained artifacts. No separate bounded CI invocation or simultaneous overlap of three specifically selected tests is required. The outcome is unchanged; no local token rotation or local live rerun is required or authorized by this amendment.

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

**AC-2 — Claude common and substrate tests use one three-slot queue with isolated host state.** The short worker loop retains a slot through synchronous `t.Run` and cleanup, including sequential variants; `-parallel 1` can reduce concurrency. Use the existing bounded spike evidence for queue behavior and independent local review. Per the captain's 2026-09-15 amendment, complete the real overlap described above in the single full stack-tip CI run before final delivery; local validation may proceed while this proof is explicitly deferred. CI must retain actual overlapping host intervals, distinct config/workflow/session paths and the original durable assertions, including both sequential break-glass variants. A fixture pass or auth skip cannot close this criterion. Existing exported direct tests remain usable with exact selectors.

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

Allowed semantics: test selector/nested event names, admission order, Claude overlap, merged config location, combined Claude JSON detail, and the explicit timeout consolidation. No spacedock CLI grammar, durable state format, authority, assertions, Pi behavior, auth policy or provider settings change. Staffing remains one implementation owner and one independent validator; one additional reviewer only for a distinct uncovered claim. Scheduling is the bottom PR, with other in-flight PRs above it. The implementation gate was approved. This state amendment authorizes no candidate edits, local reruns, code push, PR or CI action by the implementation worker; the first officer coordinates the ready stack and its single full tip CI run.

## Minimal proof plan

1. Before implementation, add the one ordering check with a small independent expected sequence and equal-hint pair. No timing sleeps, replay fixture or scheduler test framework.
2. Adapt and run the existing registry reconciliation and gotestsum proof owners. Preserve exact callable coverage and use anchored selectors. The retained worker/failure spike is sufficient design evidence for the few lines of native worker glue; review that glue directly instead of introducing new permanent lifecycle/cancellation controls.
3. Defer real Claude overlap to the single full stack-tip CI run, per captain amendment of 2026-09-15. Use the existing full `TestLiveScheduled` invocation across all 20 Claude functions and require actual common/substrate overlap with a maximum of three occupied tests, overlapping host intervals, distinct config/workflow/session paths, passing original durable assertions including both sequential break-glass variants, and retained artifacts. This does not require another bounded CI invocation or simultaneous overlap of three specifically selected tests. No local token rotation or local live rerun; historical skipped execution remains evidence of a skip, and final delivery remains conditional on this CI proof.
4. Required local `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` were completed as recorded below: both suites retain the sole confirmed resolver baseline failure; they are not all-green. Independent local validation reviews the small diff and deterministic evidence without duplicate broad runs. Full live CI waits for the ready stack and must satisfy the deferred AC-2 proof before delivery.

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

## Stage Report: implementation

- DONE: Implement the approved minimal sorted-list runner and shared three-slot cap using existing test bodies and one focused ordering check.
  Commit `4ce49f1ea` adds the literal 17-common/3-Claude table, descending hint/lexical ordering, three synchronous worker slots, common nested-parallel suppression and merged config child; existing exported test bodies remain intact.
- FAILED: Prove actual Claude overlap and isolation locally; preserve failure propagation, coverage and artifacts without duplicate scaffolding.
  The bounded real-host attempt skipped merged dispatch, break-glass and shallow boot for unavailable live auth before any host launch; actual intervals, session/config isolation and both break-glass variants remain unproven. Existing registry/gotestsum checks pass; no auth denial bypass or duplicate proof framework.
- DONE: Commit the minimal candidate with focused and required checks, exact surface count, and a self-contained report for independent validation.
  Candidate `4ce49f1ea` on `spacedock-ensign/schedule-live-journeys-duration-hints`, based on `438053493`: 8 files, +178/-32, net +146 (21.7% above +120 estimate, within +150/8 ceiling); no code push or CI.

### Proof and limitations

- `go test ./internal/ensigncycle -run '^TestLiveScheduleOrder$'` passed after an initial missing-helper compile failure: independently expects long/a/z/short; reversing hint or equal-hint comparison fails it.
- `go test ./internal/contractlint -run 'TestRuntimeLive(RegistryReconciliation|CommonSuiteTimeouts|CommonFailFastPolicy)$'` passed: omitted/duplicated or misnamed callables and incorrect canonical selectors fail; existing fixture/assertion/TODO/XFAIL reconciliation remains active.
- `go test ./internal/release -run '^TestLiveCIStep' -count=1` passed: existing controls exercise clean summary, retained JSON, one execution and underlying failure exit; losing detail or swallowing a failing process exit fails them. The existing named-evidence workflow control also passed.
- `gofmt -w ./cmd ./internal` ran; unrelated pre-existing whitespace in `internal/release/runtime_live_evidence_workflow_test.go` was restored. `git diff --check` passed.
- `go test ./...` and then `go test ./... -race` each exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`: stable `spacedock@spacedock` absent while resolver returned the installed `spacedock-local/.../0.28.0-pre0/.codex-plugin/plugin.json`. All other packages passed; no race reports. Full logs: worktree `.local/schedule-proof/go-test.log` and `go-test-race.log`.
- Resolver disposition: observed user/workflow is local developer repository verification; harm is a failing installed-host check; authority `none: unchanged resolver behavior does not affect a scheduling value AC`; trigger is the exact installed-local/absent-stable mismatch above. Worker proposed Deferred risk/outside scheduling scope; FO explicitly DECLINED a scheduling fix after matching retained independent baseline evidence. No broader green-suite claim.
- The live attempt used the current candidate binary built with `go build -o .local/schedule-proof/spacedock ./cmd/spacedock`; exported `SPACEDOCK_BIN` and `SPACEDOCK_REPO_ROOT` point to that binary and this worktree, and artifact/config roots are isolated beneath `.local/schedule-proof/`.
- Exact live command: `SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 15m -parallel 3 -run '^TestLiveScheduled$/^slot-[0-2]$/^(TestLiveMergedTeamModeDispatch|TestLiveBreakGlassShimRecovery|TestLiveCommonShallowBoot)$' ./internal/ensigncycle -json`.
- Retained `.local/schedule-proof/overlap.jsonl` records all three selected tests as `skip`, with `no live auth available: set ~/.claude/benchmark-token (operator/OAuth) or ANTHROPIC_API_KEY (CI)`. Exit 0 is not host proof. With an authorized readable credential, rerun that bounded command and retain actual overlapping host intervals, distinct paths and original durable assertions before AC-2 acceptance.

### Summary

The minimal scheduled candidate is committed and ready for independent code/fixture validation, with the canonical Claude and Codex invocations sharing their existing evidence owners and Pi unchanged. Actual Claude overlap remains an explicit acceptance blocker; required broad checks also retain the independently confirmed local resolver failure, which the FO declined to fix in this task. No provider/auth policy, global configuration or runtime assertions were changed.

### Credential availability recheck — 2026-09-15

`claude auth status` exited successfully with `loggedIn: true`, `authMethod: claude.ai`, first-party provider and team subscription. Both `ANTHROPIC_API_KEY` and `CLAUDE_CODE_OAUTH_TOKEN` are unset. The existing `decideClaudeEnv` accepts only a readable `$HOME/.claude/benchmark-token` or `ANTHROPIC_API_KEY`; native login alone cannot pass its isolated-harness gate. The previously denied token path was not retried, and no keychain extraction, global auth mutation, candidate edit or redundant live/broad run occurred.

AC-2 still requires an authorized API key in the test process or an approved readable benchmark-token at a harness-compatible HOME. Candidate remains `4ce49f1ea`; after credentials are available, rerun only the recorded bounded overlap proof and retain original assertions and actual overlapping host/path evidence. This recheck does not close the FAILED checklist item or authorize stage advancement.


## Captain amendment — CI proof location (2026-09-15)

The captain's “that can be in ci” moves real overlap verification to the single full stack-tip CI run. It supersedes the historical local-before-acceptance requirement in the implementation report and credential recheck above; their skipped-host and resolver-failure observations remain unchanged. No local credential repair or rerun is required. The task is ready for independent local validation, and final delivery remains conditional on the unchanged AC-2 outcome being proved in CI.

## Stage Report: implementation (cycle 2)

- DONE: Record the captain-authorized move of real overlap proof to stack-tip CI without weakening the required behavior.
  Risk evidence, AC-2 and minimal proof plan now place actual common/substrate overlap, three-slot bound, distinct config/workflow/session paths and passing original assertions (including both break-glass variants) in the existing full 20-function Claude `TestLiveScheduled` stack-tip CI invocation before delivery.
- DONE: Reconcile the implementation report with completed local checks, explicit deferred CI proof, and retained failure evidence.
  This report supersedes the earlier local-before-acceptance disposition; focused ordering/registry/gotestsum controls passed, while both broad suites retain the confirmed resolver baseline failure and the local live attempt remains three auth skips.
- DONE: Publish the scoped state amendment and return ready for independent local validation.
  Only this entity body/report is amended and committed path-scoped; candidate `4ce49f1ea` remains unchanged at +178/-32, net +146 across 8 files, with no candidate edit, local rerun, code push, PR or CI launch.
- DONE: Implement the approved minimal sorted-list runner and shared three-slot cap using existing test bodies and one focused ordering check.
  Retained candidate `4ce49f1ea` and prior focused evidence stand: wrong hint/tie comparison fails the ordering check; omission/duplication/callable mismatch fails registry reconciliation.
- SKIPPED: Prove actual Claude overlap and isolation locally; preserve failure propagation, coverage and artifacts without duplicate scaffolding.
  Local real-host proof is explicitly deferred by captain to stack-tip CI, not waived or claimed green; original registry/gotestsum controls passed and CI must supply actual host intervals, isolated paths, durable assertions and retained artifacts before final delivery.
- DONE: Commit the minimal candidate with focused and required checks, exact surface count, and a self-contained report for independent validation.
  Candidate and deterministic checks were committed/reported previously; normal/race exit 1 solely at the FO-declined resolver baseline is retained, and no redundant checks ran for this state-only amendment.

### Summary

Ready for independent local validation of the committed minimal candidate and its deterministic evidence. Captain-authorized CI deferral replaces the former local overlap blocker without closing AC-2: final delivery still requires the real overlap/isolation proof at the full stack tip, and neither auth skips nor the local resolver failures are called green.


## Stage Report: validation

- DONE: Assess the minimal scheduling candidate and completed deterministic evidence against AC-1 and AC-3 and the captain-amended CI proof boundary for AC-2.
  Local PASSED for candidate `4ce49f1ea34502380d892c1f96ea5c8f3f22e07e` on parent `438053493`; AC-1/AC-3 supported, AC-2 actual-host proof explicitly remains open until stack-tip CI.
- DONE: Perform the required detached adversarial audit of ordering, shared three-slot queue, isolation wiring and evidence preservation without repeating owned green suites or launching live hosts.
  Detached checkout `/tmp/schedule-validation-detached-4ce49f1`; seven minimal mutants failed existing owners, restored checkout has no diff, and live-tag compilation with `-run '^$'` passed without executing hosts.
- DONE: Report local PASSED or REJECTED with explicit unresolved CI overlap proof, exact eight-file/150-net scope, and the known resolver failure limitation.
  Exact scope is 8 files, +178/-32 = +146 net, within the +150-net/8-file ceiling; no candidate edits, new standing tests, broad reruns, live launch, PR or CI action.

### Acceptance evidence and adversarial pass

- AC-1 — Longer hinted tests are admitted first from a committed per-runtime list: reused green `TestLiveScheduleOrder`; detached reversal of hint comparison failed with `[short a z long]`, and reversal of lexical ties failed with `[long z a short]`, both against independent `[long a z short]`.
- Independently compared all 17 literal common-runtime hint pairs with the entity's measurement table: exact match; Claude-only 110/270/140 hints match the three listed substrate measurements. Sorting is local O(n log n) on at most 20 entries; no external read, automatic update, unbounded allocation or new blocking I/O.
- AC-2 — Claude common and substrate tests use one three-slot queue with isolated host state: source trace confirms one sorted list, mutex-protected exactly-once index admission, three parallel workers and synchronous child `t.Run`; queued jobs continue after ordinary child failure. Reused committed spike `72718be3f` documents peak 3, eight completions, race success and all eight completions after injected Fatal; no duplicate queue suite.
- Queue variant matrix reviewed: empty/exhausted queue returns unlocked; equal hints use lexical ties; Claude selects 20 and Codex 17; other runtimes skip wrapper; `-parallel 1` reduces active workers; child failure returns to next job; timeout/cancellation may leave work unrun. Common nested parallelism is suppressed only under scheduled ancestry; both break-glass variants remain synchronous and retain cleanup before releasing their slot. CPU host-start order need not equal mutex admission order.
- Isolation wiring traced from fresh HOME and workflow temp roots to scenario config children, launch env and session-based project readers: merged now uses `merged-team-mode`; bare/common and both break-glass variants retain distinct scenario children. The same merged config child reaches launch and reconciliation. This is wiring review, not actual Claude isolation proof.
- AC-2 outstanding: existing full 20-function Claude `TestLiveScheduled` at the complete stack tip must show actual common/substrate host interval overlap, no more than three occupied tests, distinct config/workflow/session paths, original durable assertions including both sequential break-glass variants, and retained artifacts. Auth skips and shell-marker fixtures cannot close it; no local credential action is needed under the captain amendment.
- AC-3 — The canonical scheduled invocation preserves all existing tests, assertions and evidence: retained green registry proof reconciles 17 common and three Claude callable identities; detached omission, duplication and name/callable mismatch each failed `TestRuntimeLiveRegistryReconciliation` at the exact defective row.
- Detached unanchored Claude selector failed `TestRuntimeLiveCommonSuiteTimeouts`; inserting `-failfast` failed `TestRuntimeLiveCommonFailFastPolicy`. Canonical anchored invocation runs originals through the wrapper once; broad selectors remain explicitly unsupported for exactly-once scheduling. Original exported bodies, assertions and Pi invocation are unchanged.
- Evidence trace: both runtime commands retain one gotestsum execution and original detail filenames; Claude substrate events move into `live-e2e-detail.jsonl`; `always()` uploads retain binary provenance, child config/projects, streams and metrics. Reused green `TestLiveCIStep*` controls exercise actual clean summary, failing exit and JSON retention; removing JSON output or swallowing the process exit violates those owners. Existing named-evidence control establishes step ownership only, not real host coverage.
- Audit command pattern: `go test ./internal/ensigncycle -run '^TestLiveScheduleOrder$' -count=1` for the two comparison mutants; `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1` for row mutants; analogous exact timeout/failfast owners for the two workflow mutants. Each mutant exited 1 with its expected assertion; `/tmp/schedule-validation-<mutation>.log` retains output. Detached `go test -tags live ./internal/ensigncycle -run '^$' -count=1` exited 0.
- Required implementation normal/race logs were inspected and reused: both exit 1 solely at `TestCodexResolveManifestAgainstInstalledHost` (stable marketplace manifest absent; installed local manifest selected), with no race report. Formatting and diff checks were already completed; this validation does not call those broad suites green or repeat them.

### Findings and limitations

No new material outcome/evidence defect, deferred risk or polish finding is proposed; the detached audit refuted nothing material. The already FO-declined resolver baseline remains outside this scheduling change. Actual-host AC-2 proof is an explicit captain-authorized pending delivery condition, not a waived criterion or a local validation rejection.

### Summary

Local PASSED: the minimal candidate and its existing deterministic evidence satisfy AC-1 and AC-3, and independent detached mutation checks rejected the targeted regressions. Candidate bytes and HEAD remain unchanged; final delivery remains conditional on full stack-tip CI satisfying AC-2, with the known resolver failure recorded separately.
