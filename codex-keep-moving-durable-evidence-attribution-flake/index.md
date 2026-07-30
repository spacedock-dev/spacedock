---
id: 8bnkrtq4rw46xkbez5zrbmmj
title: Codex keep-moving durable-evidence attribution false-red
status: implementation
source: "PR #513 Runtime Live E2E run 29392675038, codex-live job 87279446937"
started: 2026-07-30T23:05:41Z
completed:
verdict:
score: 0.8
worktree: .worktrees/spacedock-ensign-codex-keep-moving-durable-evidence-attribution-flake
issue:
milestone: 0.25.0
gates:
    version: 1
    current:
        gate: gate:8bnkrtq4rw46xkbez5zrbmmj:ideation
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
                state: consumed
                blockers: []
        - id: gate:8bnkrtq4rw46xkbez5zrbmmj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:8bnkrtq4rw46xkbez5zrbmmj-ideation-1
              briefing:
                id: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-1:revision-1
                digest: sha256:a9e29deee9056d26159a1c6772e5b213e73467dd03ce07cfe9186d9660bbdde2
                digest-domain: canonical-bytes
                request-digest: sha256:1e8eef615d1bd19e2011de8bac8ca2650c19caa690e13e9a85903a430c39ba30
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:8bnkrtq4rw46xkbez5zrbmmj:ideation:1
                briefing: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T23:04:28.063854Z"
                decision: approve
                reason: 'Captain conn accepted the no-transcript, net-negative design: ordered per-task Git history directly tests the keep-moving journey, adversarial controls can falsify it, and the implementation boundary forbids observer dialects or product semantics.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
---

The keep-moving live scenario must credit each completed task from its own durable workflow journey, without making a provider transcript dialect part of the product verdict.

## Problem

PR #513's `codex-live` lane failed `TestLiveCodexSharedScenarios/keep-moving-posture` even though `approved-gate`, `ready-one`, and `ready-two` had each completed the real workflow journey: dispatch entry, worker completion, terminalization, and archive. The observer instead tried to reconstruct that journey from Codex JSONL, report-heading cardinality, command text, and merge-loop output. Harmless changes in those representations produced a false red twice after the same exact-head behavior had passed offline, install, Claude/Opus, and Pi validation.

The failure is architectural, not another dialect bug. Adding a generic-heading exception, recognizing a shell variable, normalizing `item.completed`, or replaying a smaller JSONL sample would retain a second protocol whose grammar changes independently of Spacedock behavior. The captain therefore rejected transcript grammar parsing in full.

Exact evidence: [PR #513](https://github.com/spacedock-dev/spacedock/pull/513), [Runtime Live E2E run 29392675038 / codex-live job 87279446937](https://github.com/spacedock-dev/spacedock/actions/runs/29392675038/job/87279446937).

## Spike result

The cheapest mechanism already exists in `internal/ensigncycle/liveassert_test.go`: locate an entity in its active or archive location and interrogate its path-scoped Git history. The focused spike

`go test ./internal/ensigncycle -run 'TestLocateEntity|TestSomeCommitNamesOnly|TestIntegrationTransitionCommitted|TestCompletedSetAnchor|TestLiveStageReportHeading|TestTerminalFrontmatterAnchors' -count=1 -v`

passed on 2026-07-31. Its falsifying control is load-bearing: a clean `status: done` entity with a path-scoped terminal commit passes the superficial terminal checks but `TestIntegrationTransitionCommitted/skipped_advance_fails` rejects it because the required earlier transition is absent. Therefore final archived state alone is insufficient; ordered, per-path durable history is the minimum proof.

## Proposed approach

Replace the keep-moving transcript trace with one host-neutral durable-task journey oracle. For each expected completed slug, follow that entity's Git history from its archive path and require, in ancestry order:

1. A path-scoped `dispatch: {slug} entering {stage}` commit whose entity blob has the expected stage and non-empty `started`. This is the existing dispatch-entry contract, not a new receipt.
2. A later path-scoped worker commit whose entity blob contains that stage's durable Stage Report. `started` alone never earns dispatch credit; the later worker-owned report is what closes the dispatched-work claim under the existing FO/ensign authority split.
3. A later terminal blob with non-empty `completed` and `verdict`, followed by the entity existing only at its canonical archive location. The archive path and Git ancestry bind all facts to the same slug; timestamps are not used for ordering.

Grade the three independent tasks as a set of three separate journeys. `questioned` remains a negative guard: it must be durably re-shaped and nonterminal, and none of its commits may satisfy another slug. The live runner supplies only the workflow root and expected slugs/stages; it does not supply JSONL, final narration, commands, or provider events.

Reuse the same durable completion oracle for the commissioned-task completion fallback in the smallest-sufficient-mechanism scenario, so deleting the shared Codex evidence parser does not silently leave a second consumer behind. That scenario's unrelated edit/commit scope checks stay unchanged.

The rejected alternatives are: final archive state alone (falsified by the spike), `status --archived` alone (no ordered dispatch/completion proof), an instrumented wrapper log (a new observer schema), and any transcript/event/command parser (the rejected architecture).

## Deletion inventory

- Delete all 433 lines of `internal/ensigncycle/shared_keep_moving_test.go`: provider event decoding, command token/regex inference, narration inference, host-neutral motion trace, and transcript correlation tests.
- Delete all 409 lines of `internal/ensigncycle/shared_keep_moving_negative_test.go`: fabricated Claude/Codex streams, replay dialects, and final-message grammar controls.
- Delete all 202 lines of `internal/ensigncycle/codex_dispatch_evidence_test.go`: `item.completed` decoding, dispatch-build result decoding, Stage Report cardinality, status text matching, merge-output matching, and shell-loop variable parsing.
- Delete all 350 lines of `internal/ensigncycle/codex_dispatch_evidence_regression_test.go`: JSONL constructors and observer compatibility cases. Move only the generic `codexCommandOutput` fixture constructor if its unrelated round-recording consumer still needs it.

No compatibility is retained for these internal observer formats.

## Expected surface and semantic boundary

Expected files: delete the four files above; add `internal/ensigncycle/shared_keep_moving_durable_test.go`; adjust `shared_smallest_mechanism_test.go`, `shared_fixtures_test.go`, `codex_live_runner_test.go`, `claude_live_runner_test.go`, and `docs/runtime-live-ci.md`. A tiny shared fixture-helper move is allowed if compilation requires it.

Budget: at most 10 files plus one helper-only file; at most 300 inserted lines; at least 1,100 deleted lines; cumulative diff must remain at least 700 lines net negative. Tolerance is +1 file and +80 insertions only when needed to keep an unrelated test fixture compiling; the net-negative floor is not waived.

Observable semantics changed: test-oracle runtime behavior only. Keep-moving and commissioned completion are credited from existing durable state and Git ancestry rather than transcript/final-message syntax. Command grammar, CLI output, stored formats, mutation authority, runtime dispatch behavior, retry policy, and provider adapters do not change. There is no new observer schema, receipt, retry controller, transcript grammar, or runtime normalization.

## Out of scope

- The permission-narration oracle removal in #514.
- The separate Codex rejection-flow worker-reuse flake.
- Retry policy or the task-15 wait-watchdog replacement.
- Changes to the keep-moving behavioral contract or its requirement to commission each ready entity.
- Runtime host adapters, dispatch tools, workflow frontmatter, and Stage Report format.

## Acceptance criteria

**AC-1 - Every completed independent task is credited from its own durable journey.**
Verified by: a deterministic real-Git fixture produces three ordered dispatch-entry → worker-report → terminalize → archive journeys and the oracle reports `3/3`; removing any one journey reports `2/3`. This replaces PR #513's false-red baseline, where the completed motion was credited as fewer than `3/3`.

**AC-2 - Missing, stale, reordered, or cross-attributed durable steps remain red per task.**
Verified by: table-driven real-Git controls independently remove the dispatch-entry commit, worker report, terminal fields, or archive; place a report before dispatch; and give one slug another slug's report/commit. Each control names and rejects only the affected slug.

**AC-3 - The observer surface is smaller and provider-independent.**
Verified by: `git diff --numstat` meets the declared deletion and net-negative floors; focused tests accept identical durable journeys with empty/arbitrary transcript and final-message bytes, and no keep-moving completion code reads JSONL, shell/JavaScript text, provider event types, or model narration.

**AC-4 - Repository and live confirmation gates are green.**
Verified by: focused durable-journey tests, `gofmt -l` empty for changed Go files, `go test ./...`, `go test ./... -race`, and one exact-head Codex keep-moving live run after deterministic proof is green.

## Test plan

Implement one deterministic real-Git fixture that uses the existing entity layout, dispatch commit convention, Stage Report contract, terminal fields, and archive locations. The positive case costs four commits per task at most; table-driven negatives mutate one fact at a time and assert the exact slug/reason, including stale order and cross-attribution. No JSONL fixture is added.

Run the focused durable tests first, then `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. After those pass, run exactly one exact-head Codex `keep-moving-posture` live scenario and retain its ordinary runtime artifacts only for diagnosis, not grading.

## Stage Report: ideation

- DONE: Replace the seed’s transcript-dialect replay/parser direction with a behavior-first design that observes actual dispatch, completion, terminalization, and durable per-task state; deletion must be the primary mechanism metric.
  The design uses existing per-entity state plus path-scoped Git ancestry and deletes 1,394 parser/replay lines, with a net-negative floor of 700 lines.
- DONE: Spike the cheapest real mechanism that can distinguish completed keep-moving behavior from false-red without parsing model prose, shell commands, JavaScript wrappers, or provider event dialects; record the falsifying result before finalizing the design.
  Focused durable helper tests passed; the skipped-transition control falsified final-state-only grading while ordered Git history rejected it.
- DONE: Specify value ACs, semantic boundary, exact files/LOC and tolerance, negative controls, and one live confirmation only after deterministic behavior proof; forbid new observer schema, retry controller, transcript grammar, and unrelated runtime changes.
  ACs, a 10-file/+1 tolerance, insertion/deletion floors, per-task controls, semantic exclusions, and the single post-deterministic Codex run are recorded above.

### Summary

Ideation now removes the rejected transcript observer instead of extending its dialect grammar. The replacement proves each completed task from existing dispatch-entry, worker-report, terminal, archive, and Git-history facts, with a durable falsifier and strict net-negative implementation budget.
