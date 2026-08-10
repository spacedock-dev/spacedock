---
name: feedback-rejection-flow
description: "First-officer feedback-rejection routing — the correction-round record, reuse-or-fresh routing of findings to the `feedback-to` target, the reviewer re-run, gate re-entry. Invoke at the rejection-handling point when a feedback gate recommends REJECTED or the captain rejects at a feedback-to stage."
user-invocable: false
---

# Feedback Rejection Flow

This skill loads at the rejection-handling point; rejection detection, correction-record write scope, reuse conditions, and the budget probe stay always-on in the FO contract.

## Feedback Rejection Flow

When a feedback stage recommends REJECTED:

1. Read the rejected stage's `feedback-to` target — the stage that receives the fix request, not the reviewer.
2. Read the already-authorized workflow package: rejected snapshot, finding evidence, existing workflow classifications, FO-authorized dispositions, and concrete revise assignment. If the distinct authorization or assignment is missing, hold at the active workflow's review-finding checkpoint; routing is ineligible.
3. Invoke `«context-budget»()`. If the old ensign is over budget or the source is unavailable, shut down and fresh-dispatch; if no probe is declared, proceed to reuse below.
4. Route the authorized package unchanged to the target stage in the same worktree using `«addressable-worker»` when the existing handle is addressable and reuse conditions pass; otherwise shut down and fresh-dispatch. The routed message carries the concrete assignment, not an acknowledgment request or a new classification request. Do not treat the immediate routing response as completion: wait for the reused worker's next `«completion-signal»` when it is on the critical path, and attribute completion by mailbox content, task path, or durable workflow state.
5. After correction completes, wait until the worker entries are complete. If the workflow declares a `### Feedback Cycles` correction-round projection, the First Officer must append its authorized line directly. Do not define, normalize, or interpret its category labels, fields, tolerance, estimate, or drift grammar.
6. Invoke the neutral `${SPACEDOCK_BIN:-spacedock} gate record --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl` once. Require a successful result with the complete round summary. If the command fails, produces no result, or retains an incomplete log, hold the flow. Do not claim that the round was recorded. Do not re-run the reviewer or prepare the next gate. The recorder retains the canonical two-file room and advances `review-round`, without receiving or interpreting the Cycle line. On cycle 3, publish the rejected round, then escalate to the human instead of another reviewer round.
7. Re-run the reviewer after publication. When the existing reviewer remains addressable and reuse conditions pass, re-run the kept-alive reviewer through the same `«addressable-worker»` capability used for feedback routing; the message must ask that reviewer to re-review the updated entity state, not validate its own fix work. Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail.
8. Re-enter the normal gate flow with the updated result.

The FO owns the shared correction-round section and writes it under `«write.classify»`: worktree-side when `worktree:` is set, main-side otherwise. The generic recorder owns only the immutable room and pointer bytes.

## Workflow-defined correction-round projection

The active workflow owns whether a projection exists and its exact grammar. This skill transports the selected projection without supplying defaults or interpreting its contents.
