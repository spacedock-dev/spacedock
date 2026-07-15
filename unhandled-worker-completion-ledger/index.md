---
title: Make unhandled worker completions durable and scheduler-visible
status: backlog
source: "Captain report of s0 completion parked behind unchanged implementation state; relevant to c6 continuation ledger feed, 2026-07-15"
started:
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: 4etz8eaktsc19jmq8e9k76ch
---

The first-officer contract already requires a completed non-gated stage to advance, dispatch, or terminalize before the FO yields. The remaining failure is event-loop execution: a worker completion can arrive while the FO handles another captain interaction, remain unacknowledged behind unchanged entity frontmatter, and disappear from scheduling because `status --next` sees only workflow state.

## Observed s0 failure

1. The implementation worker completed with accepted exceptions but omitted required follow-up links.
2. The FO held advancement, filed the follow-ups, and requested a report-only correction.
3. The corrected completion arrived while the FO handled a captain sprint question.
4. The FO inspected the sprint and ran `status --next`.
5. `status --next` returned empty because `s0` still had `status: implementation`; it could not see that the implementation worker had completed.
6. The FO answered and yielded instead of advancing `s0` and dispatching fresh validation.

The result was an unacknowledged completion event parked behind unchanged workflow state.

## Cause

- **Primary:** FO scheduling allowed a new query to preempt an unhandled completion, then treated an empty state-derived ready set as quiescence.
- **Tooling contributor:** `status --next` cannot represent completed workers awaiting advancement.
- **Codex contributor:** `list_agents` observes worker completion, but no durable workflow record says whether the FO handled that completion.
- **Semantic contributor:** a completion may mean `stage_complete`, `blocked`, or `report_update`; the FO must verify the durable report before acting.
- **Not the cause:** the contract's completion-first rule is already explicit. Compaction amplifies the failure but did not cause this incident.

## Required outcome

Add a persistent `awaiting FO action` completion feed to the launcher-owned continuation ledger. A completion record contains:

```text
entity
completed stage
worker outcome: stage_complete | blocked | report_update
tip
report commit
completion epoch
handled: false
```

`spacedock status --next` or the future `spacedock dispatch next-action` must surface unhandled completions before ordinary dispatchables. The FO acknowledges a record only after it advances the entity, terminalizes it, or records a real block. An empty dispatchable set cannot mean idle while an unhandled completion exists. Records survive compaction and host restart.

The runtime loop drains completion observations after every tool return, wait return, compaction, and user interruption; verifies the durable report; advances, terminalizes, or records a block; dispatches the next stage when required; and only then may declare quiescence or yield a final response.

For Codex, reconcile the durable feed with `list_agents` so a completed worker and its completion epoch cannot remain indefinitely unacknowledged. Do not treat roster state alone as stage success; the report and entity state remain authoritative.

## Relationships and scope

- **Feeds into `c6` `codex-post-compaction-contract-reload`:** extend its typed continuation ledger and priority order rather than create a parallel ledger or lifecycle controller.
- **Coordinates with `jv` `codex-worker-shutdown-and-roster-reconcile`:** reuse its worker identity, roster, and completion-epoch evidence where compatible.
- **Coordinates with roadmap 0222:** the unshipped `spacedock dispatch next-action` is the preferred deterministic consumer of this feed.
- A mandatory turn-boundary check for unhandled completions, completed stages still recorded at the old stage, and newly unblocked dependencies is an acceptable near-term guard only if backed by a failing behavioral test. Prose alone is not the end state.
- Do not add a watcher, daemon, lease, retry controller, or transcript/narration oracle.

## Acceptance criteria

- **AC-1 (VALUE — no parked completions):** A worker completion that arrives during a user query while the entity still names the completed stage is durably visible as unhandled and is selected before the FO may report idle or yield.
  - **Verified by:** a Codex-backed or fixture-backed scheduler journey recreates the `s0` ordering and fails unless the completion is verified, acknowledged through advance/terminal/block, and followed by the required next-stage dispatch.
- **AC-2 (durability):** The completion record carries entity, completed stage, typed worker outcome, tip, report commit, completion epoch, and `handled`; it survives process restart and compaction without duplicating or losing the obligation.
  - **Verified by:** restart/rehydration tests persist an unhandled record, reopen the ledger, reconcile the matching worker/report evidence, and prove the same epoch is handled exactly once.
- **AC-3 (scheduler priority):** `status --next` or `dispatch next-action` reports unhandled completions before ordinary dispatchables, and never reports quiescence while one exists.
  - **Verified by:** a deterministic mixed-queue fixture contains an unhandled completion, an ordinary ready task, and unchanged frontmatter; the returned action names completion verification first, then the ready task only after acknowledgment.
- **AC-4 (truthful acknowledgment):** The FO may set `handled: true` only after durable evidence shows the matching epoch advanced, terminalized, or recorded a real block; a mailbox notification or roster status alone cannot acknowledge it.
  - **Verified by:** negative tests for missing report, mismatched epoch/tip/report commit, `report_update`, and blocked outcomes keep the record unhandled until the corresponding authoritative state transition or block record exists.
- **AC-5 (Codex reconciliation):** Codex roster reconciliation imports a matching completed worker observation into the ledger or matches it to an existing record without allowing stale epochs, report-only updates, or blocked workers to masquerade as stage completion.
  - **Verified by:** `list_agents` fixture cases cover stage completion, report-only correction, blocked completion, stale epoch, duplicate observation, compaction, and user-interruption ordering.
