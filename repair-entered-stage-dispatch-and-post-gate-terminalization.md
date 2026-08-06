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
---

PR #585's runtime-neutral durable live oracle exposed three shipped behavior defects outside the PR's Codex configuration and launcher-shim scope.

1. Initial-stage successor dispatch is underspecified. At exact head `d28834249b23df204292149c7581a295e85c10dd`, both Codex/Luna and Claude Sonnet changed `ready-one` and `ready-two` from `ready` directly to terminal `done`, committed `dispatch: <slug> entering done`, and built `--stage done`, while the fixture and worker assignment require a `ready` stage report. Owner: `skills/first-officer/references/fo-dispatch-core.md` plus its mandatory skill smoke tests. Make the initial-stage rule explicit and prove `current=ready,next=done` dispatches `current` with `status=ready started`, an exact path-scoped `dispatch: <slug> entering ready` commit, and `dispatch build --stage ready`.

2. Gate-consume dispatch evidence can be split across commits. Sonnet consumed `approved-gate` from `review` to `implementation` and committed `dispatch: approved-gate entering implementation` before `started` existed, then added `started` in a second commit named `dispatch: approved-gate implementation started`. The durable oracle correctly requires the exact dispatch commit's entity blob to contain both the entered stage and nonempty `started`. Owner: gate-lifecycle/dispatch contract and smoke tests. Prove a consumed nonterminal approval records the dispatch boundary only after `started` is present, without weakening path-scoped durable grading.

3. A consumed nonterminal gate record blocks later ordinary terminalization. Both hosts received `condition=consumed` after the approved `review -> implementation` transition. `internal/status/handlers.go` sends every terminal status write through `pendingTerminalApproval`; `internal/status/merge.go` classifies consumed/nonterminal authority as an error. Codex could not finalize; Sonnet escaped with `--force`, violating the live scenario's non-force safety posture. Owner: `internal/status`/`internal/gates`. Add a focused CLI regression for `review(gate) -> implementation(consumed) -> done`, allowing the atomic non-forced terminal fields while pending, unreadable, stale, or binding terminal-target authority remains fail-closed.

Evidence: archived Sonnet artifact `runtime-live-e2e-claude-live-sonnet`, `claude-shared-scenarios-detail.jsonl`, and scenario streams under `live-artifacts/claude/sonnet/claude-shared-scenarios/`; archived Codex artifact `runtime-live-e2e-codex-live` and the corresponding scenario streams. Restore the TODO-disabled live journeys only after focused smoke/CLI proof lands; do not relax the durable Git-history oracle or teach agents to use `--force`.
