---
title: Make the host-neutral dispatch core genuinely runtime-neutral
status: backlog
source: "captain + FO follow-up after 2y merged as v0.20.3 (2026-06-16) — remaining concern: the extracted host-neutral dispatch core still carries Claude/team-only language, even though Codex and Pi now load it too. Verified against origin/main 9bd1f46a: fo-dispatch-core.md still says first team-mode dispatch, keeps spawn-standing-all --team {team_name} in the core, requires team_name as helper output, and makes reuse depend on not being in bare mode."
started:
completed:
verdict:
score: 0.66
worktree:
issue:
id: ezfwkw33awtqgztgr6v7bb59
sprint: 0204-structured-reads
sprint-readiness: ready
---

The 2y merge fixed the important reachability gap by extracting `fo-dispatch-core.md` and `fo-merge-core.md`, making the merge and dispatch ceremony reachable to Codex and Pi. One follow-up remains: the new `fo-dispatch-core.md` is named as host-neutral, but parts of it still describe Claude team mode as if it were universal.

## Problem

On `origin/main` at `9bd1f46a` (`v0.20.3` plus manifest stamp), `skills/first-officer/references/fo-dispatch-core.md` still carries host-specific assumptions:

- The module says it is loaded at the first "team-mode dispatch" even though Codex/Pi do not necessarily have team mode.
- Standing-teammate injection is in the host-neutral core and calls `spacedock dispatch spawn-standing-all --team {team_name}`, which is Claude-team-shaped.
- Reuse condition 1 says `Not in bare mode (teams available)`, which makes Codex `send_input` reuse look invalid despite Codex having a reusable mailbox handle.
- The dispatch builder section says `team_name` is a key field that MUST come from helper output, while Codex explicitly has no `team_name` lifecycle.
- The core says bare mode dispatch blocks until subagent completion, which is a host behavior, not a host-neutral invariant.

This is not a blocker for the 2y merge because the extracted cores are reachable and tested, but it is a real contract-quality issue: future Codex/Pi FOs load a host-neutral core that still reads partly like Claude.

## Desired direction

Keep the 2y split. Do not collapse the cores back into per-host copies.

Make `fo-dispatch-core.md` describe only host-neutral dispatch lifecycle rules:

- load at first worker dispatch, not first team-mode dispatch;
- dispatch via the runtime adapter's spawn call;
- reuse when the adapter exposes a live reusable handle and all generic reuse conditions pass;
- helper output fields are forwarded when emitted/applicable, with host adapters mapping concrete fields;
- blocking/nonblocking dispatch and bare/team behavior live in runtime adapters.

Move or gate Claude-only pieces in `claude-fo-dispatch.md`:

- `spawn-standing-all --team {team_name}` and standing teammate injection;
- Claude `team_name` / TeamCreate sequencing;
- bare-mode blocking semantics;
- any rule whose source of truth is Claude team behavior.

Codex should remain explicit that `spawn_agent` is initial dispatch, `send_input` is reuse/feedback routing, `wait_agent` is only a foreground idle wait, there is no `team_name` lifecycle, and there is no reconcile sweep. Pi should retain its `subagent(...)` / `pi-agent-teams` substrate split.

## Acceptance criteria

**AC-1 — The dispatch core contains no unconditional Claude/team-only dispatch requirements.**
Verified by a structural contractlint check in `internal/contractlint/` that fails if `fo-dispatch-core.md` contains unconditional `team-mode`, `spawn-standing-all`, or `--team {team_name}` requirements outside an explicitly adapter-delegated sentence. The check should be structural and narrow; do not turn it into a broad prose-grep behavior substitute.

**AC-2 — Codex reuse remains expressible through the shared core.**
Verified by a fixture or existing ensigncycle test showing a Codex continuation/reuse path uses `send_input` against a live worker handle and is not rejected merely because there is no team mode or `team_name`.

**AC-3 — Claude behavior is preserved.**
Verified by existing focused gates: `go test ./internal/contractlint`, `go test ./internal/dispatch -run 'Standing|Build'`, and the Claude standing/teardown focused tests that cover `spawn-standing-all`, TeamCreate sequencing, and the terminal marker.

## Test plan

Start with a small structural or fixture test that reproduces the current leak: the host-neutral core should not require team-mode dispatch, `spawn-standing-all --team`, or `team_name` as universal dispatch facts.

Then refactor the prose only: host-neutral language stays in `fo-dispatch-core.md`; Claude team behavior moves to or remains in `claude-fo-dispatch.md`; Codex/Pi adapters state their concrete runtime mapping.

Run:

```bash
go test ./internal/contractlint
go test ./internal/dispatch -run 'Codex|Pi|Standing|Build'
go test ./internal/ensigncycle -run 'Codex|SharedReviewerReuse|WaitWatchdog|ForceTeamMode|TestGradeMarkerMatchesContract'
go test ./...
```
