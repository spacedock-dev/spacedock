---
name: feedback-rejection-flow
description: "First-officer feedback-rejection routing — the correction-round record, reuse-or-fresh routing of findings to the `feedback-to` target, the reviewer re-run, gate re-entry. Invoke at the rejection-handling point when a feedback gate recommends REJECTED or the captain rejects at a feedback-to stage."
user-invocable: false
---

# Feedback Rejection Flow

This skill loads at the rejection-handling point; rejection detection, correction-record write scope, reuse conditions, and the budget probe stay always-on in the FO contract.

## Feedback Rejection Flow

When a feedback stage recommends REJECTED, run these five steps in order. Each is one action with one completion condition, and is unfinished until that condition holds: never treat a step's first command as its completion.

1. **Deliver the authorized correction.** Read the rejected stage's `feedback-to` target — the stage that receives the fix request, not the reviewer — and the already-authorized workflow package: rejected snapshot, finding evidence, existing workflow classifications, FO-authorized dispositions, and concrete revise assignment. When the workflow has no review-finding checkpoint, the captain's concrete revise instruction supplies correction authority; do not invent development-specific classifications. If the distinct authorization or assignment is missing, hold; routing is ineligible. Route the package unchanged to the target stage in the same worktree, carrying the concrete assignment, not an acknowledgment request or a new classification request. Reuse the existing handle through `«addressable-worker»` only when it is addressable, reuse conditions pass, and `«context-budget»()` reports it under budget; otherwise shut down and fresh-dispatch. If no probe is declared, proceed to reuse.
   **Done when** the correction is complete in durable workflow state: the target worker's own entries closing this round's review log where the workflow keeps one, otherwise its `«completion-signal»` attributed by mailbox content, task path, or durable state. The immediate routing response is never completion.

2. **Record required evidence, once.** When the active workflow requires a canonical correction round, record it exactly once with the existing recorder. Otherwise continue without creating a round. Apply the workflow's Feedback Cycles projection separately; its absence does not waive a required canonical round. Append the authorized `### Feedback Cycles` line for this round when the active workflow declares that projection, and, when a canonical round is required, invoke the neutral recorder exactly once for the whole rejection cycle: `${SPACEDOCK_BIN:-spacedock} gate record ENTITY --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl`. ENTITY is the entity operand the recorder requires; omitting it is a usage error, not a recording. CYCLE is this rejection cycle's number and the only round id this cycle publishes; there is no second publication after the reviewer re-run. Do not define, normalize, or interpret the Cycle line's category labels, fields, tolerance, estimate, or drift grammar. The recorder retains the canonical two-file room and advances `review-round`, without receiving or interpreting the Cycle line.
   **Done when** every declared projection is written and either no canonical round is required or the recorder exits successfully reporting the complete round summary, counting every entry this round accumulated. A usage error — the invocation itself was malformed — is not a round failure: correct the invocation and run it again, which completes this one publication rather than adding a second. If the recorder instead refuses the round, produces no result, or reports an incomplete round, hold the flow: do not claim that the round was recorded, re-run the reviewer, or prepare the next gate.

3. **Commit the correction and required evidence.** Invoke `«state.commit»(slug)`.
   **Done when** the corrected entity and all required round evidence are durably committed and the tree is clean for that entity. An absent round obligation requires no manufactured evidence. Until it is, the entity is absent from `status --next` and refused by `gate prepare`: an empty scheduler here means this step is unfinished, never that the run is over.

4. **Satisfy the review requirement and cycle limit.** On cycle 3, record escalation to the human and stop here, even when no reviewer is declared; skipping review never bypasses this limit. When the workflow requires independent review, obtain its fresh verdict on the corrected state. Never use the correction worker as its own reviewer. Otherwise, completed and committed correction work satisfies the review prerequisite. Same-stage feedback alone does not waive declared review requirements. Re-run the kept-alive reviewer through the `«addressable-worker»` capability step 1 used, when it remains addressable and reuse conditions pass; fresh-dispatch a reviewer only when it is not.
   **Done when** this cycle has every required fresh reviewer verdict, or no review is required and the committed correction is complete — except on cycle 3, when the escalation must be recorded and the flow must stop here.

5. **Re-enter the gate.** Invoke `Skill(skill="spacedock:fo-gate-lifecycle")` and run `«gate.lifecycle»` for the updated stage. A `needs-preparation` row is work to perform, not a stopping condition. Reach this step only after every applicable requirement in steps 1–4 completes. A `needs-preparation` row observed before then is work pending those requirements, not authority to prepare; the rejected report alone can surface one. When independent review is required, the old rejected report cannot substitute for its fresh verdict.
   **Done when** exactly one fresh open gate has been prepared and presented, and the flow stops — without resolving or applying it, changing terminal state, or dispatching a successor.

The FO owns the shared correction-round section and writes it under `«write.classify»`: worktree-side when `worktree:` is set, main-side otherwise. The generic recorder owns only the immutable room and pointer bytes.

## Workflow-defined correction-round projection

The active workflow owns whether a projection exists and its exact grammar. This skill transports the selected projection without supplying defaults or interpreting its contents.
