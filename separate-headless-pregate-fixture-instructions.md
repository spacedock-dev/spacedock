---
title: "Separate pre-gate dispatch instructions from gate-ready fixture instructions"
status: validation
source: "Captain filing request; pre3 and 0.27.3 live release failure triage"
started: 2026-09-15T04:02:22Z
completed: ""
verdict: ""
score: "0.9"
worktree: .worktrees/spacedock-ensign-separate-headless-pregate-fixture-instructions
issue: ""
pr: ""
mod-block: ""
id: 8k3cmcg2gg1e88w0qe6sq3vw
gates:
    version: 1
    records:
        - id: gate:8k3cmcg2gg1e88w0qe6sq3vw:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:8k3cmcg2gg1e88w0qe6sq3vw-backlog-1
              briefing:
                id: briefing:8k3cmcg2gg1e88w0qe6sq3vw:backlog:attempt-1:revision-1
                digest: sha256:7922fa71be8298f5d1f0c06ae08aa6f4bc4587e1074e8774450b93afc4595f17
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:8k3cmcg2gg1e88w0qe6sq3vw:backlog:1
                briefing: briefing:8k3cmcg2gg1e88w0qe6sq3vw:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T04:02:07.367243Z"
                decision: approve
                reason: Captain requested dispatch of the Codex live failure tasks, targeted local verification, and stacked PRs.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:8k3cmcg2gg1e88w0qe6sq3vw:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:8k3cmcg2gg1e88w0qe6sq3vw-ideation-1
              briefing:
                id: briefing:8k3cmcg2gg1e88w0qe6sq3vw:ideation:attempt-1:revision-1
                digest: sha256:a51f84f66dd4cf42b45627b2d736df37a4f9ba1e0df6fd2347d30ca901a90c71
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:8k3cmcg2gg1e88w0qe6sq3vw:ideation:1
                briefing: briefing:8k3cmcg2gg1e88w0qe6sq3vw:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T04:07:50.357067Z"
                decision: approve
                reason: Current fixture contradiction and coherent control are reproduced; two-file correction preserves worker and human-authority boundaries.
                conn:
                    quote: dispatch codex live test failure tasks, verify locally for targeted failure, and open PR as stack, then trigger codex ci on stack tip.
                    source: active thread goal supplied by captain
              application:
                target-stage: implementation
                state: consumed
---

Separate pre-gate dispatch instructions from gate-ready fixture instructions.

## Problem

Both release Codex runs attempted gate prepare while recorded-gate-task was queued; the binary correctly refused that queued is not an actionable gate. No implementation dispatch occurred.

Evidence: https://github.com/spacedock-dev/spacedock/actions/runs/34923154318 (pre3, job 104235689038, artifact 10379087396) and https://github.com/spacedock-dev/spacedock/actions/runs/34923156550 (0.27.3, job 104235722711, artifact 10379038870), default-headless-gate-stop.

runGateStopScenario shares gatePrompt between pre-gate and gate-ready fixtures. gatePrompt requests a recorder-ready room from a committed review. writePreGateWorkflow resets the entity to queued but retains a prewritten validation gate-review recommending approval. These are contradictory task cues. Existing completed headless-recorded-gate-stop-stage-coherence and select-actionable-codex-default-headless-task addressed earlier state/dispatch defects; this task owns the remaining prompt and review-artifact contradiction. The pre3 fixture additionally commits README; both releases still show the same premature prepare.

## Proposed approach

Use distinct instructions for pre-gate and already-gated starting states. The pre-gate request drives current work to the human decision boundary. Remove completed-review cues from its initial artifacts. Preserve implementation-worker dispatch and all gate authority assertions. If a coherent fixture still reproduces the failure, record that separately as runtime evidence rather than relaxing the guard.

## Risk evidence

The captured release artifacts above establish the failure trigger. Replay them before implementation; a live runtime claim requires the targeted live AC below, not a prose-presence check.

## Out of scope

Unrelated release failures and broader test-harness redesign.

## Expected surface and tolerance

Files: `internal/ensigncycle/shared_fixtures_test.go` (fixture, prompt, existing coherence test) and `internal/ensigncycle/claude_live_runner_test.go` (prompt selection). Estimate net LOC change: +35, across 2 files; approximately 45 insertions and 10 deletions, net tolerance +20 to +60. A third adjacent existing test file is permitted only for a distinct fixture/lifecycle falsifier.

Observable semantics permitted: the test-only pre-gate starting artifacts and request become coherent with queued implementation. Command grammar, stored product formats, approval authority, product guards, and host lifecycle semantics do not change. No user-facing documentation changes are needed; this is internal harness behavior.

## Acceptance criteria

**AC-1 — Pre-gate setup has no completed-review evidence.**
Verified by: fixture construction and on-disk checks show queued work with no completed validation review or report; inserting such a completed artifact makes the check fail.

**AC-2 — Codex dispatches implementation before presenting the human gate.**
Verified by: a targeted local live default-headless-gate-stop journey observes worker spawn and completion before gate preparation, then stops without approval/consume/successor dispatch. Skipping the worker or consuming the gate fails existing lifecycle assertions.

**AC-3 — The already-gated control retains its intended boundary.**
Verified by: the existing gate-ready control remains valid and the relevant fixture/lifecycle tests reject extra dispatch or authority consumption.

## Test plan

Primary deterministic owner: `TestPreGateWorkflowIsStageCoherent` in `shared_fixtures_test.go`; extend it before implementation to check `os.Stat` reports no selected gate-review file, no parsed completed report sections, and real boot JSON still offers queued → implementation with zero ready gates. Restoring the retained review must fail. Keep `TestAssertGateHeld` and `TestAssertGateHeldAcceptsPreparedFixtureBinding` as existing authority/control owners; an approval or consume in the held entity must fail.

Primary live owners: `TestLiveCommonDefaultHeadlessGateStop` and `TestLiveCommonGateGuardrail`. Preserve `assertImplementationWorkerLifecycle`, `assertRecordedGateHoldLog`, and `assertGateHeld`; skipping worker completion, moving to validation early, approving/consuming, or dispatching a successor must remain failures. Any further native lifecycle failure from the fresh live run must be named separately; other stack tasks do not own this observer. Do not substitute prompt substring tests for this behavior.

Focused layer command: `go test ./internal/ensigncycle -run 'Test(PreGateWorkflowIsStageCoherent|AssertGateHeld|AssertGateHeldAcceptsPreparedFixtureBinding)$' -count=1` (seconds, deterministic).

Targeted local Codex command, from the approved stack checkout after building its binary: `go build -o /tmp/8k3cmcg2gg-candidate-spacedock ./cmd/spacedock`, then `SPACEDOCK_BIN=/tmp/8k3cmcg2gg-candidate-spacedock SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_CODEX_LIVE_REQUIRED=1 SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/8k3cmcg2gg-codex-live go test -tags live ./internal/ensigncycle -run '^TestLiveCommon(DefaultHeadlessGateStop|GateGuardrail)$' -count=1 -parallel=1 -timeout=20m -v`. Requires local Codex CLI and existing auth, reused through the existing isolated-home harness; record resulting artifacts and durable state. This paid live check is minutes, not a new framework.

Captain-approved stack validation: run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` once at stack tip; per-layer validation is focused and targeted live. Report unavailable live execution separately. No prose-grep proof.

## Bounded implementation plan

1. Extend the existing fixture test with on-disk absence checks before changing setup. Preserve its real boot/dispatchable assertion and the prepared control.
2. In `writePreGateWorkflow`, remove and commit deletion of `selected/gate-review.md` together with the queued entity/snapshot reset. Keep the neutral recorder contract and implementation-stage instruction. The validation review must be produced from actual completed implementation, not seeded as evidence.
3. Add `preGatePrompt` and select it only for `default-headless-gate-stop` in the existing runner. Keep `gatePrompt` for already-gated/withdrawn controls. No generalized scenario configuration abstraction.
4. New pre-gate wording: “Use $spacedock:first-officer for this whole run. Workflow directory: <root>. Engage only `recorded-gate-task`. Drive its current work through implementation to the human decision boundary and stop there.” Existing gate-ready wording remains “Prepare its recorder-ready room from the committed gate review and supporting References.” Neither prompt grants the conn or approval authority.

The two changes serve AC-2 together: removing only the review leaves a request to prepare nonexistent committed evidence; changing only the prompt leaves contradictory completed-review evidence. Removing this one artifact and selecting one prompt is sufficient; a new fixture framework or product guard change is unnecessary.

## Ideation proof

Source inspected: `origin/main` at `2a7b8719843e40b79545f0bb4def6609cdd9ebbf`, not the root worktree HEAD. Read `docs/runtime-support.md` and apply its assume-it-works policy to later live setup.

Throwaway snapshot `/tmp/8k3cmcg2gg-ideation-spike`, unique test `internal/ensigncycle/ideation_8k3cmcg2gg_test.go`: `go test ./internal/ensigncycle -run 'Test(Ideation8k3FixtureEvidence|PreGateWorkflowIsStageCoherent|AssertGateHeld)$' -count=1 -v` passed in 4.015s. Observed: pre-gate retains a 346-byte review; real `gate prepare` exits 1 with “workflow stage queued is not an actionable gate”; gate-ready control creates `state=open`. Existing coherence test passes despite retained review. This exercises fixture and guard behavior; it does not prove the new prompt succeeds live.

The release runs linked above remain the captured runtime baseline. Captured payloads are `/tmp/spacedock-pre3-codex-live.zip` and `/tmp/spacedock-0273-codex-live.zip` (FO supplied). Direct captured-artifact replay remains an implementation prerequisite; the targeted new live run is the value proof and runs after FO scheduling.

### Feedback Cycles


## Stage Report: ideation

- DONE: Establish the smallest coherent pre-gate setup and distinct prompt while preserving the gate-ready control and authority assertions.
  Two-file plan removes the retained review and selects a pre-gate request; origin/main spike observed queued refusal and a successful gate-ready prepare.
- DONE: Record the bounded implementation plan, test owners, and exact targeted local Codex command for the approved stack.
  Body names existing fixture/live owners, net +35 LOC baseline, falsifiers, two exact focused commands, and stack-tip normal/race/format validation.
- SKIPPED: Fresh targeted Codex execution during ideation.
  No candidate changes are authorized yet; targeted live runs after FO scheduling, and no live success is claimed.

### Summary

The current fixture test passes while the queued entity retains a 346-byte completed review; the real binary correctly refuses premature gate preparation. The plan removes that contradiction without changing product authority or lifecycle assertions, and keeps any newly observed native lifecycle issue separate.

## Stage Report: implementation

- DONE: Remove completed review cues from pre-gate setup and select a distinct current-work prompt, preserving gate-ready controls and strict lifecycle assertions.
  Commit `9e4c43ceff6fb9f6c500cd602f4fb18cc5eaae3b` removes and commits the seeded review deletion; only default-headless-gate-stop selects preGatePrompt. Two files, +27/-3 (net +24), within approved tolerance.
- DONE: Commit the bounded two-file change with focused regression proof and prepare the exact targeted local Codex run for FO scheduling.
  Clean branch rebased onto approved `origin/main` `2a7b8719843e40b79545f0bb4def6609cdd9ebbf`; focused command below passed in 5.767s, and live-tag deterministic compilation/control run passed in 4.353s.
- DONE: Replay captured release evidence before implementation.
  Temporary `TestReplay8k3CapturedPreGateFailure` applied unchanged assertRecordedGateHoldLog to command.log from both supplied release ZIPs; both rejected with no successful gate prepare recorded. Temporary test removed after replay.
- DONE: Prove the fixture regression before the fix.
  TestPreGateWorkflowIsStageCoherent failed with retained gate-review stat err=<nil>; after deletion it verifies absent review, parsed stage-report absence, and real boot queued → implementation with zero ready gates. Restoring seeded review fails its on-disk check.
- DONE: Preserve gate-ready authority controls.
  TestAssertGateHeld, TestAssertGateHeldAcceptsPreparedFixtureBinding, and TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle pass; approval/consume mutants remain rejected. No lifecycle assertion changed.
- SKIPPED: Fresh targeted Codex run and stack-wide normal/race/format checks in this worker turn.
  FO explicitly serializes live execution and schedules the full/race/format checks once at combined tip; no live success is claimed and AC-2 remains pending that run.

### Focused proof and scheduled live command

`go test ./internal/ensigncycle -run 'Test(PreGateWorkflowIsStageCoherent|AssertGateHeld|AssertGateHeldAcceptsPreparedFixtureBinding)$' -count=1`

`go test -tags live ./internal/ensigncycle -run 'Test(PreGateWorkflowIsStageCoherent|AssertGateHeld|AssertGateHeldAcceptsPreparedFixtureBinding|AssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle)$' -count=1`

Both edited files were gofmt formatted; `git diff --check` passed. From the approved stack checkout after FO schedules:

```bash
go build -o /tmp/8k3cmcg2gg-candidate-spacedock ./cmd/spacedock
SPACEDOCK_BIN=/tmp/8k3cmcg2gg-candidate-spacedock SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_CODEX_LIVE_REQUIRED=1 SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/8k3cmcg2gg-codex-live go test -tags live ./internal/ensigncycle -run '^TestLiveCommon(DefaultHeadlessGateStop|GateGuardrail)$' -count=1 -parallel=1 -timeout=20m -v
```

### Summary

The queued fixture now starts without a completed review and asks the first officer to drive current implementation work to the human boundary. Gate-ready prompts and strict authority/lifecycle checks retain their previous behavior; fresh Codex evidence is intentionally pending FO scheduling on the combined stack.


## Stage Report: validation

- DONE: Independently review headless layer da50d61d relative keep-moving predecessor 7ca67fdd; verify fixture absence and preserved authority/lifecycle controls.
  Reviewed exact `da50d61d6df3e8347ab9a8d594e31a8d7cf3cfba` against `7ca67fdda5433654b6f67276560b3eb19c0684f6`: two test files, +27/-3, clean candidate; no product or assertion weakening.
- DONE: Assess AC-1 — Pre-gate setup has no completed-review evidence.
  Reviewed existing green implementation proof without rerunning per FO: TestPreGateWorkflowIsStageCoherent checks disk absence, parsed/raw report absence, and real boot queued → implementation with zero ready gates; restoring the review fails os.Stat, while a completed report fails the heading checks.
- DONE: Perform the bounded semantic adversarial review.
  Pre-gate alone deletes the selected review and receives current-work instructions; ready and withdrawn fixtures retain their review and original prompt. Removal is committed path-scoped, and missing deletion fails setup rather than silently retaining stale evidence.
- DONE: Verify preserved negative lifecycle and authority invariants.
  Unchanged held-state checks reject wrong gate/briefing/attempt/digest, approval/application, verdict, and successor status; command-log checks reject repeated prepare, consume, successor dispatch, withdrawal, or post-prepare status mutation. Native lifecycle checks require a correlated completion before validation and one parsed DONE report.
- SKIPPED: After FO grants the serialized live slot, run targeted Codex default-headless and gate-ready control, preserve artifacts and exact commit, and assess all ACs without editing candidate code.
  FO withheld slot after the preceding keep-moving run hit local Codex 0.154 code-mode negotiation failure before gate commands; minimal isolated host smoke is pending. This validator launched no live run and attributes no candidate defect from that sibling failure.
- SKIPPED: AC-2 — Codex dispatches implementation before presenting the human gate.
  Fresh runtime evidence remains pending; unchanged assertions and prior labels alone do not establish worker dispatch/completion, gate preparation order, and the final unconsumed boundary.
- SKIPPED: AC-3 — The already-gated control retains its intended boundary (fresh live portion).
  Deterministic control proof and unchanged boundaries were reviewed; TestLiveCommonGateGuardrail must still run alongside the default-headless scenario after the host smoke and slot grant.
- DONE: Prepare exact-candidate execution without candidate changes.
  `go build -o /tmp/spacedock-headless-da50d61d ./cmd/spacedock` succeeded; `git diff --check` passed and HEAD stayed exact/clean. Scheduled artifacts: `/tmp/spacedock-stack-headless-live`; mandated CI model shim and isolated copied local OAuth will be used.
- SKIPPED: Repeat deterministic/full/race suites.
  FO explicitly retained existing green focused evidence and assigned full/race/format checks once at the combined tip; no new test concern justified repetition here.

### Summary

Review found no material candidate defect within the approved two-file fixture/prompt correction. Validation remains pending, with no PASSED recommendation until both targeted Codex journeys produce fresh evidence; the worker remains addressable for the serialized slot. No code, frontmatter, PR, or merge was changed.


## Review-finding disposition

### Validation finding H-1 — completed implementation did not enter validation

- Reviewer observation: exact candidate `da50d61d6df3e8347ab9a8d594e31a8d7cf3cfba` dispatched and completed its implementation worker, then attempted prepare at `status: implementation`; the binary correctly refused. No validation stamp/dispatch or prepared Briefing exists.
- Released user and normal workflow: Codex first officer driving the supported default-headless queued task to its human validation gate.
- Observable harm: the committed implementation report never becomes a reviewed, prepared human gate; the run stops early.
- Affected authority: value-ac[AC-2] the live journey must dispatch implementation and reach the prepared human boundary without consuming it.
- Trigger evidence: fresh default-headless FAIL in 122.09s, native `spawns=1 completed=130 validation=-1 report=<nil>`; preserved state commits `019cd54` dispatch then `0993259` implementation report; prepare exits 1.
- Advisory classification: outcome defect, material observed AC-2 failure. Task ownership: separate FO completion-routing investigation, coordinated with explicit stage-completion work; no remaining fixture-review contradiction observed.
- Proposed disposition: route for decision; keep this candidate and lifecycle oracle unchanged. The smallest missing behavior is implementation completion → FO validation transition/review dispatch → completed validation report → prepare; no instruction-only repair is proven sufficient.
- FO authorization: pending separate disposition; validator has made no candidate edit or retry.

## Stage Report: validation (cycle 2)

- DONE: Independently review headless layer da50d61d relative keep-moving predecessor 7ca67fdd; verify fixture absence and preserved authority/lifecycle controls.
  Earlier review stands; candidate remained exact and clean after both live runs, with unchanged strict assertions.
- DONE: After FO grants the serialized live slot, run targeted Codex default-headless and gate-ready control, preserve artifacts and exact commit, and assess all ACs without editing candidate code.
  Both targeted journeys ran once serially at `da50d61d6df3e8347ab9a8d594e31a8d7cf3cfba`, current CLI 0.154.0, CI shim Luna/max; total 248.729s, exit 1; slot released to FO.
- DONE: AC-1 — Pre-gate setup has no completed-review evidence.
  Existing deterministic proof retained per FO; live initial entity was queued with no report, and selected gate-review was absent. Worker later committed exactly one implementation report, showing it was newly produced.
- FAILED: AC-2 — Codex dispatches implementation before presenting the human gate.
  TestLiveCommonDefaultHeadlessGateStop FAIL 122.09s: real correlated spawn/completion succeeded, but validation=-1 and no prepared gate. Removing completion or skipping validation still fails the unchanged lifecycle oracle.
- DONE: AC-3 — The already-gated control retains its intended boundary.
  TestLiveCommonGateGuardrail PASS 126.18s: one successful prepare, state commit/head, matching held gate, no approval/consume/successor. Wrong binding or authority consumption would fail the unchanged state/log checks.
- DONE: Diagnose the new failure without retrying or changing candidate code.
  Loaded fo-dispatch-core explicitly requires FO advancement and gated-stage dispatch after implementation completion; trace instead loads gate lifecycle while status remains implementation. Original queued/no-worker issue is repaired; H-1 records the separate observed routing failure.
- DONE: Preserve exact source, host, state, and command evidence.
  `/tmp/spacedock-stack-headless-live/` contains source.json, test.log, finding.md, commands-summary.json, scenario command/JSONL logs, retained native rollouts, and retained workflow/state Git snapshots; no copied auth is included in retained snapshots.
- SKIPPED: Full/race suites, new candidate edits, and live retries.
  Full/race checks belong to combined-tip FO validation; H-1 requires a distinct FO disposition before any remedy/rerun. Local auth remained unchanged and isolated runtime markers were cleaned as in the proven host smoke.

### Summary

Recommend REJECTED against unchanged AC-2: the fixture correction now reaches genuine worker completion, but the first officer fails to enter validation before preparing its human gate. AC-1 remains supported and fresh AC-3 passed; no lifecycle assertion was relaxed and no candidate defect is inferred merely from the test's broad `implementation-worker-not-dispatched` label. H-1 awaits FO ownership/disposition, with the live slot released and this worker addressable.
