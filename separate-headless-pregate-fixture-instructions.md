---
title: "Separate pre-gate dispatch instructions from gate-ready fixture instructions"
status: ideation
source: "Captain filing request; pre3 and 0.27.3 live release failure triage"
started: 2026-09-15T04:02:22Z
completed: ""
verdict: ""
score: "0.9"
worktree: ""
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
