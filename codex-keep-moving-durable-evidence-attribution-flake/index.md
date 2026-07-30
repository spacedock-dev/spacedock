---
id: 8bnkrtq4rw46xkbez5zrbmmj
title: Codex keep-moving durable-evidence attribution false-red
status: backlog
source: "PR #513 Runtime Live E2E run 29392675038, codex-live job 87279446937"
started:
completed:
verdict:
score: 0.8
worktree:
issue:
milestone: 0.25.0
gates:
    version: 1
    current:
        gate: gate:8bnkrtq4rw46xkbez5zrbmmj:backlog
    records:
        - id: gate:8bnkrtq4rw46xkbez5zrbmmj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:8bnkrtq4rw46xkbez5zrbmmj-backlog-1
              briefing:
                id: briefing:8bnkrtq4rw46xkbez5zrbmmj:backlog:attempt-1:revision-1
                digest: sha256:f74c978f7bac349914d3380c445388cf82767c0647d3267768c425649fb719a8
                digest-domain: canonical-bytes
                request-digest: sha256:83ed6b8050afb6ea2e3737bf14a7c57f164095922ad8f37db84546e5cde6f84d
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:8bnkrtq4rw46xkbez5zrbmmj:backlog:1
                briefing: briefing:8bnkrtq4rw46xkbez5zrbmmj:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-30T22:50:16.874174Z"
                decision: approve
                reason: Captain explicitly directed dispatch of 8b after Q3 merge and ruled that transcript-grammar parsing is offrail; ideation must replace the dialect proposal with behavior and durable-state proof.
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

The Codex keep-moving grader must recognize a completed per-entity dispatch and terminalization without turning harmless transcript-shape variation into an unrelated PR blocker.

## Problem

PR #513's `codex-live` lane failed `TestLiveCodexSharedScenarios/keep-moving-posture` even though the transcript and durable state showed `approved-gate`, `ready-one`, and `ready-two` dispatched, reported, terminalized, and archived. This is a new variant of the evidence dialect tracked and closed by [`fo-function-reference-invariant`](../_archive/fo-function-reference-invariant/index.md): Codex emitted no visible `spawn_agent` event, `approved-gate` used `## Stage Report` without the canonical colon/stage suffix, and finalization ran through a shell loop whose command named a variable rather than each entity literally. The existing fallback credited the two canonical report headings but could not attribute the third report or the loop's successful per-entity merge-guard output, so it reported a missing dispatch.

This exact signature is not covered by the unrelated `codex-live-rejection-flow-flake` or by `retire-keep-moving-permission-narration-oracle` (#514). Re-running is reasonable, but leaving the variant untracked makes the same unrelated live-lane false-red likely to recur.

Exact evidence: [PR #513](https://github.com/spacedock-dev/spacedock/pull/513), [Runtime Live E2E run 29392675038 / codex-live job 87279446937](https://github.com/spacedock-dev/spacedock/actions/runs/29392675038/job/87279446937).

## Proposed approach

Start with a minimized committed replay fixture derived from the failing Codex JSONL: successful per-entity dispatch builds, a completed collaboration wait, a post-wait batched durable read with two canonical `## Stage Report:` headings and one generic `## Stage Report` heading, then successful per-entity merge-guard results emitted from a variable-driven loop. Use the fixture to decide the smallest authoritative boundary: either make the worker/report path repair or reject noncanonical headings before the FO claims completion, or let the grader consume structured durable entity/finalization evidence that does not depend on raw Markdown heading cardinality or literal shell arguments.

Do not merely broaden the heading regex, count incidental prose, infer success from the final summary, or add a general shell parser. Preserve the ordered successful build → subsequent completed wait → subsequent durable evidence invariant and its per-entity attribution controls.

## Out of scope

- The permission-narration oracle removal in #514.
- The separate Codex rejection-flow worker-reuse flake.
- Retry policy or the task-15 wait-watchdog replacement.
- Changes to the keep-moving behavioral contract or its requirement to commission each ready entity.

## Acceptance criteria

**AC-1 - The PR #513 completed-motion dialect no longer false-reds.**
Verified by: a committed minimized replay fixture derived from job 87279446937 fails on the pre-fix grader and passes after the repair, with `approved-gate`, `ready-one`, and `ready-two` each independently credited from durable evidence.

**AC-2 - Invalid or cross-attributed streams remain red.**
Verified by: the retained stale-report-before-build, report-before-wait, failed-build, failed-target-in-a-batch, one-report-for-multiple-targets, and incidental-prose controls all remain failing cases; add a negative proving one entity's loop output cannot bless another entity.

**AC-3 - The authority boundary is structural and per entity.**
Verified by: focused tests demonstrate that success comes from a canonicalized report/state transition or structured per-entity command result, not final narration, raw `Stage Report` substring counts, or a general shell-command parser.

**AC-4 - Repository and live confirmation gates are green.**
Verified by: focused keep-moving replay tests, `gofmt -l` empty for changed Go files, `go test ./...`, `go test ./... -race`, and one exact-head Codex keep-moving live run.

## Test plan

Add the minimized PR #513 replay and adversarial per-entity negatives to `internal/ensigncycle`. Run focused keep-moving tests first, then the full and race suites. Because this bug is specific to Codex's live transcript dialect, require one exact-head Codex live keep-moving confirmation after all deterministic tests pass; do not use repeated green retries as the primary proof.
