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

Test fixtures, scenario prompt selection, and adjacent behavioral tests in internal/ensigncycle; estimate net +20 to +60 LOC across 2–3 files. No product guard, host adapter, or workflow state semantics change.

## Acceptance criteria

**AC-1 — Pre-gate setup has no completed-review evidence.**
Verified by: fixture construction and on-disk checks show queued work with no completed validation review or report; inserting such a completed artifact makes the check fail.

**AC-2 — Codex dispatches implementation before presenting the human gate.**
Verified by: a targeted local live default-headless-gate-stop journey observes worker spawn and completion before gate preparation, then stops without approval/consume/successor dispatch. Skipping the worker or consuming the gate fails existing lifecycle assertions.

**AC-3 — The already-gated control retains its intended boundary.**
Verified by: the existing gate-ready control remains valid and the relevant fixture/lifecycle tests reject extra dispatch or authority consumption.

## Test plan

Add focused failing behavioral cases first. Reuse the existing test owner and captured artifacts, then run the targeted live AC. Before completion run go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal. Report unavailable live execution separately from passing evidence. No prose-grep proof.

### Feedback Cycles

