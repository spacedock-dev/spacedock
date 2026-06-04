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
4. Consult the budget probe (reuse condition 0). If the old ensign is over budget or the source is unavailable, shut down and fresh-dispatch; if no probe is declared, proceed to reuse below.
5. Route findings back to the target stage in the same worktree using the existing handle when addressable and reuse conditions pass (`send_input` on Codex, `SendMessage` on Claude teams); otherwise shut down and fresh-dispatch. The routed message must carry the concrete next-stage assignment and fix work, not just an acknowledgment request. On Codex, do not treat the immediate `send_input` response as the new completion result — if the follow-up is on the entity's critical path, wait for the reused worker's next completion before advancing or shutting it down (entity-scoped wait, not a global scheduling stop).
6. Re-run the reviewer after fixes.
7. Re-enter the normal gate flow with the updated result.

The FO owns `### Feedback Cycles`. Routing follows FO Write Scope: worktree-side when `worktree:` is set, main-side otherwise.
