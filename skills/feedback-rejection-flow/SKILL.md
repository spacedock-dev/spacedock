---
name: feedback-rejection-flow
description: "First-officer feedback-rejection routing — read the `feedback-to` target, track `### Feedback Cycles`, escalate on cycle 3, consult the budget probe, route findings back to the target stage in the worktree (else fresh), re-run the reviewer, re-enter the gate flow. Invoke at the rejection-handling point when a feedback gate recommends REJECTED or the captain rejects at a feedback-to stage."
user-invocable: false
---

# Feedback Rejection Flow

This skill carries the first-officer's feedback-rejection routing procedure. The rejection DETECTION stays always-on in the FO contract's gate flow; this skill loads at the rejection-handling point to route the findings back to the target stage and re-run the reviewer. The `### Feedback Cycles` write-scope rules, the reuse conditions, and the budget probe stay always-on in the FO contract — this procedure references them by name.

## Feedback Rejection Flow

When a feedback stage recommends REJECTED:

1. Read the rejected stage's `feedback-to` target — the stage that receives the fix request, not the reviewer.
2. Track cycles in `### Feedback Cycles` in the entity body.
3. On cycle 3, escalate to the human instead of another round.
4. Invoke `«context-budget»()`. If the old ensign is over budget or the source is unavailable, shut down and fresh-dispatch; if no probe is declared, proceed to reuse below.
5. Route findings back to the target stage in the same worktree using `«addressable-worker»` when the existing handle is addressable and reuse conditions pass; otherwise shut down and fresh-dispatch. The routed message must carry the concrete next-stage assignment and fix work, not just an acknowledgment request. Do not treat the immediate routing response as the new completion result: if the follow-up is on the entity's critical path, wait for the reused worker's next completion through `«completion-signal»` before advancing or shutting it down. Attribute completion by mailbox content, task path, or durable workflow state, not by narration.
6. Re-run the reviewer after fixes. When the existing reviewer remains addressable and reuse conditions pass, re-run the kept-alive reviewer through the same `«addressable-worker»` capability used for feedback routing; the message must ask that reviewer to re-review the updated entity state, not validate its own fix work. Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail.
7. Re-enter the normal gate flow with the updated result.

The FO owns `### Feedback Cycles`. Routing follows the eagerly loaded first-officer `«write.classify»` scope: worktree-side when `worktree:` is set, main-side otherwise.
