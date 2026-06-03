---
id: gnanx8t260dyax3x6s841bgc
title: Make the Codex no-wait subagent completion behavior retryable
status: backlog
source: "captain (2026-06-03) — dogfood session found Codex may queue subagent completion notifications without explicit wait_agent; document a repeatable idle-wake probe in case Codex changes behavior"
started:
completed:
verdict:
score: "0.26"
worktree:
issue:
---

Codex first-officer operation depends on knowing when `wait_agent` is required and when a worker can safely run in the background. The current adapter says the Codex completion signal is the async final-status notification in the FO mailbox and `wait_agent` is only an optional accelerator. A 2026-06-03 dogfood session confirmed a useful but subtle behavior: a no-wait probe worker completed and its notification was queued before the captain's next message. It did not fully prove that the FO is automatically re-entered while the user stays idle, because queued notification delivery and actual model wake-up are observably different.

## Problem

The current evidence for Codex no-wait completion is anecdotal and easy to misread. If Codex changes how interactive sessions wake on background subagent completion, the FO runtime contract could become either too blocking or too optimistic:

- too blocking: the FO foreground-waits after every dispatch even when Codex could wake from the mailbox, slowing interactive workflow operation and making steering awkward;
- too optimistic: the FO assumes background completion will wake it, leaves critical-path workers unwatched, and does not continue dispatching until the captain sends another message.

The team needs a documented recipe that separates three cases:

1. `wait_agent` foreground wait returns or is interrupted.
2. No `wait_agent`, but later user/tool activity flushes a queued completion notification.
3. No `wait_agent` and no user/tool activity; Codex independently wakes the FO from the subagent notification.

## Retry recipe

Run this manually in an interactive Codex FO session when checking host behavior:

1. Start from a clean observation point and record the wall-clock time.
2. Spawn a short no-write probe worker with `spawn_agent(fork_context=false)` and an instruction like:

   ```text
   Idle-wake notification probe. Do not edit files. Do not inspect the repo. Wait about 25 seconds, then finish with exactly this final message: Done: idle-wake Codex notification probe completed.
   ```

3. Do not call `wait_agent`.
4. End the FO turn immediately.
5. Do not type, shell out with `!`, run terminal jobs, or trigger any other activity for at least 90 seconds.
6. Record whether the UI/model starts a new assistant turn solely from the subagent completion notification.
7. If no wake occurs, send a timestamped captain message and record whether the queued subagent notification arrives before that message in the delivered event batch.

For comparison, run a second probe with explicit foreground waiting:

1. Spawn the same short worker.
2. Tell the captain that the wait is preemptible and steering input is safe.
3. Call `wait_agent` on the returned handle.
4. Interrupt once with a non-stopping message or `! echo ...`, then resume waiting on the same handle unless the captain says pause or stop.

## Acceptance criteria

Each AC names a property of the finished task and a check outside this task body that can fail.

**AC-1 - The Codex FO runtime contract distinguishes queued notification delivery from autonomous idle wake-up.**
Verified by: an instruction-text test over `skills/first-officer/references/codex-first-officer-runtime.md` asserts the adapter names both behaviors and says `wait_agent` is required when the FO is blocked on a critical-path result unless idle-wake behavior is deliberately being exercised.

**AC-2 - A repeatable manual probe exists for Codex no-wait completion semantics.**
Verified by: a checked-in dev recipe or fixture document includes the no-wait idle-wake steps, the foreground-wait comparison steps, expected observations, minimum wait duration, and interpretation rules for queued-before-user-message versus autonomous FO wake-up.

**AC-3 - The current Codex behavior is captured with timestamped evidence.**
Verified by: a small checked-in evidence file records a run date, Codex CLI version, worker handle, spawn time, expected delay, notification time or next-user-message time, and whether the FO visibly woke before user activity.

**AC-4 - The runtime guidance stays compatible with background notifications without requiring foreground waits for every dispatch.**
Verified by: a host-neutrality or skill-text test rejects wording that says every fresh Codex dispatch must immediately call `wait_agent`, while still requiring `wait_agent` or explicit operator-visible waiting when the next FO state transition depends on the worker result.

## Test plan

- Add narrow instruction-text tests for the Codex runtime adapter wording. Cost: low.
- Add a dev recipe or evidence document that can be rerun manually after Codex CLI upgrades. Cost: low.
- Optionally add a live Codex probe later if the multi-agent tool surface exposes enough timestamped event data to avoid relying on UI observation. Cost: medium to high.

## Notes

- 2026-06-03 dogfood observation: `Dalton` (`019e8e34-ec78-7691-b56d-14528a0a9f71`) was spawned without `wait_agent` and returned `Done: idle-wake Codex notification probe completed.` The notification appeared in the delivered batch before the captain asked "and how long should we wait?", but that only proves queued notification ordering, not necessarily autonomous FO re-entry while the user remains idle.
- The same session found a foreground wait can be interrupted by captain input, and a shell-out like `! echo ...` is a steering path while waiting. That path should be documented as operator interaction, not mistaken for an idle-wake proof.
