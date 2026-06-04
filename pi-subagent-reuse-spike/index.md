---
title: Pi subagent reuse and feedback-to resume semantics
status: backlog
source: captain (2026-06-04) — Pi runtime can technically resume subagents, but Spacedock has not proven reusable worker semantics for feedback-to loops or token/cost tradeoffs
score: "0.33"
started:
completed:
verdict:
worktree:
issue:
id: 7qcp0ttfhzq01cetbx85vcjx
---

# Pi subagent reuse and feedback-to resume semantics

Explore whether Spacedock's Pi runtime should support reusing/resuming an existing Pi subagent worker instead of always dispatching a fresh worker, especially for validation `feedback-to: implementation` loops.

## Problem

The current Pi runtime contract treats fresh redispatch as the safe default. `pi-subagents` can technically resume a previous run with `subagent({ action: "resume", id, message })`, and Spacedock's Pi runtime notes already mention minimal reuse metadata and completion epochs. But Spacedock has not proven what is actually achievable or desirable for workflow semantics.

The main unknowns are:

- whether resumed Pi workers retain useful task context without becoming stale or confused;
- whether a validation rejection can route back through a previous implementation worker (`feedback-to`) safely;
- whether reuse meaningfully reduces token/time cost versus fresh dispatch;
- how to prevent stale completions from satisfying a new follow-up assignment;
- what durable state should record run IDs, epochs, and reuse decisions.

## Proposed approach

Ideation should spike the runtime mechanism before designing product behavior. The spike should use the smallest safe Pi subagent experiment that can answer what is achievable without mutating product code:

1. Launch a Pi subagent worker with a tiny temp workflow/entity assignment.
2. Record returned run/session identifiers and approximate token/time cost if available.
3. Resume the same worker with a feedback-style follow-up assignment.
4. Verify the resumed worker can write a second durable marker/report and commit it, or record exactly why not.
5. Compare fresh redispatch for the same follow-up: behavior quality, elapsed time, token/cost observability, and state complexity.

The ideation output should decide whether Spacedock should implement first-class Pi reuse now, keep fresh dispatch as default, or only document resume as an operator/debug tool.

## Out of scope

- Implementing production reuse behavior in this stage.
- Changing existing Pi frontdoor/support code before the spike is understood.
- Reusing Claude/Codex team semantics blindly; Pi reuse must be Pi-native.
- Treating transcript text as proof without durable entity/file evidence.

## Acceptance criteria

**AC-1 - The spike proves what Pi subagent resume can and cannot do for a feedback-to loop.**
Verified by: ideation report containing the exact live/manual commands or subagent calls used, durable marker/entity evidence, and pass/fail classification for resumed follow-up.

**AC-2 - Reuse is compared against fresh redispatch for cost and reliability.**
Verified by: ideation report includes elapsed-time and available token/cost observations for resumed follow-up versus fresh follow-up, or records exact missing telemetry if Pi does not expose it.

**AC-3 - Stale-completion and epoch risks are explicitly modeled.**
Verified by: proposed design names the durable metadata needed to prevent a previous completion from satisfying a new follow-up: run/session handle, worker label, entity slug, stage, epoch, state, and completion evidence.

**AC-4 - The recommendation is concrete.**
Verified by: ideation report chooses one of: implement first-class Pi reuse, keep fresh dispatch default and defer reuse, or support resume only as manual/debug tooling, with tradeoffs and test plan.

## Test plan

- Live/manual Pi spike in ideation when Pi auth, `pi-subagents`, and local package paths are available; otherwise record exact missing prerequisites and do not recommend implementation.
- Use a temp split-root workflow/entity, not product state, for live spike evidence.
- If implementation is recommended, future tests should include registry epoch unit tests, Pi runtime instruction invariant tests, and a live resume smoke that verifies two distinct durable markers/reports across initial and resumed assignments.
