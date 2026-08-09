---
title: Repair entered-stage dispatch evidence and post-gate terminalization
status: backlog
source: "PR #585 exact-head run 30706782428: Codex job 91387287118 and Sonnet job 91387287120"
started:
completed:
verdict:
score: 0.9
worktree:
issue:
milestone: 0.27.0
id: 9adv48yhye5s2vkhwd7ge52d
sprint: test-behavior-completeness
gates:
    version: 1
    records:
        - id: gate:9adv48yhye5s2vkhwd7ge52d:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9adv48yhye5s2vkhwd7ge52d-backlog-1
              briefing:
                id: briefing:9adv48yhye5s2vkhwd7ge52d:backlog:attempt-1:revision-1
                digest: sha256:2cb413910c222d4d8b9a3b47fa2f43705c6a5198c796988425da0956e65f8c5e
                request-digest: sha256:9bcea1bdbd613334d989eaf10920c3e50f40c4e6f6386993352a0bd258b4989e
                room-ref: ./repair-entered-stage-dispatch-and-post-gate-terminalization/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9adv48yhye5s2vkhwd7ge52d:backlog:1
                briefing: briefing:9adv48yhye5s2vkhwd7ge52d:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:17.217951Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; shaping assigns initial-stage dispatch to task 6x and keeps post-gate behavior here.
              application:
                target-stage: ideation
                state: pending
---

PR #585's runtime-neutral durable live oracle exposed three shipped behavior defects outside the PR's Codex configuration and launcher-shim scope.

Before a product change, use the strict XFAIL behavior from
`ts7gq0mr9s3chx2w4wppd1kt`. Run each affected target and record one stable
semantic failure code for each executable journey cell. Keep TODO only when a
cell cannot run.

After a product repair, XPASS must force removal of the repaired XFAIL binding.
Do not weaken the durable oracle to obtain PASS or XFAIL.

1. Initial-stage successor dispatch is underspecified. At exact head `d28834249b23df204292149c7581a295e85c10dd`, both Codex/Luna and Claude Sonnet changed `ready-one` and `ready-two` from `ready` directly to terminal `done`, committed `dispatch: <slug> entering done`, and built `--stage done`, while the fixture and worker assignment require a `ready` stage report. Owner: `skills/first-officer/references/fo-dispatch-core.md` plus its mandatory skill smoke tests. Make the initial-stage rule explicit and prove `current=ready,next=done` dispatches `current` with `status=ready started`, an exact path-scoped `dispatch: <slug> entering ready` commit, and `dispatch build --stage ready`.

2. Gate-consume dispatch evidence can be split across commits. Sonnet consumed `approved-gate` from `review` to `implementation` and committed `dispatch: approved-gate entering implementation` before `started` existed, then added `started` in a second commit named `dispatch: approved-gate implementation started`. The durable oracle correctly requires the exact dispatch commit's entity blob to contain both the entered stage and nonempty `started`. Owner: gate-lifecycle/dispatch contract and smoke tests. Prove a consumed nonterminal approval records the dispatch boundary only after `started` is present, without weakening path-scoped durable grading.

3. A consumed nonterminal gate record blocks later ordinary terminalization. Both hosts received `condition=consumed` after the approved `review -> implementation` transition. `internal/status/handlers.go` sends every terminal status write through `pendingTerminalApproval`; `internal/status/merge.go` classifies consumed/nonterminal authority as an error. Codex could not finalize; Sonnet escaped with `--force`, violating the live scenario's non-force safety posture. Owner: `internal/status`/`internal/gates`. Add a focused CLI regression for `review(gate) -> implementation(consumed) -> done`, allowing the atomic non-forced terminal fields while pending, unreadable, stale, or binding terminal-target authority remains fail-closed.

Evidence: archived Sonnet artifact `runtime-live-e2e-claude-live-sonnet`, `claude-shared-scenarios-detail.jsonl`, and scenario streams under `live-artifacts/claude/sonnet/claude-shared-scenarios/`; archived Codex artifact `runtime-live-e2e-codex-live` and the corresponding scenario streams. Restore the TODO-disabled live journeys only after focused smoke/CLI proof lands; do not relax the durable Git-history oracle or teach agents to use `--force`.
