---
name: feedback-rejection-flow
description: "First-officer feedback-rejection routing — the correction-round record, reuse-or-fresh routing of findings to the `feedback-to` target, the reviewer re-run, gate re-entry. Invoke at the rejection-handling point when a feedback gate recommends REJECTED or the captain rejects at a feedback-to stage."
user-invocable: false
---

# Feedback Rejection Flow

This skill loads at the rejection-handling point; rejection detection, correction-record write scope, reuse conditions, and the budget probe stay always-on in the FO contract.

## Feedback Rejection Flow

When a feedback stage recommends REJECTED, run these five steps in order. Each is one action with one completion condition, and is unfinished until that condition holds: never treat a step's first command as its completion.

1. **Deliver the authorized correction.** Read the rejected stage's `feedback-to` target — the stage that receives the fix request, not the reviewer — and the already-authorized workflow package: rejected snapshot, finding evidence, existing workflow classifications, FO-authorized dispositions, and concrete revise assignment. If the distinct authorization or assignment is missing, hold at the active workflow's review-finding checkpoint; routing is ineligible. Route the package unchanged to the target stage in the same worktree, carrying the concrete assignment, not an acknowledgment request or a new classification request. Reuse the existing handle through `«addressable-worker»` only when it is addressable, reuse conditions pass, and `«context-budget»()` reports it under budget; otherwise shut down and fresh-dispatch. If no probe is declared, proceed to reuse.
   **Done when** the correction is complete in durable workflow state: the target worker's own entries closing this round's review log where the workflow keeps one, otherwise its `«completion-signal»` attributed by mailbox content, task path, or durable state. The immediate routing response is never completion.

2. **Record this round, once.** Append the authorized `### Feedback Cycles` line for this round when the active workflow declares that projection, then invoke the neutral recorder exactly once for the whole rejection cycle: `${SPACEDOCK_BIN:-spacedock} gate record ENTITY --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl`. ENTITY is the entity operand the recorder requires; omitting it is a usage error, not a recording. CYCLE is this rejection cycle's number and the only round id this cycle publishes; there is no second publication after the reviewer re-run. Do not define, normalize, or interpret the Cycle line's category labels, fields, tolerance, estimate, or drift grammar. The recorder retains the canonical two-file room and advances `review-round`, without receiving or interpreting the Cycle line.
   **Done when** the recorder exits successfully reporting the complete round summary, counting every entry this round accumulated. A usage error — the invocation itself was malformed — is not a round failure: correct the invocation and run it again, which completes this one publication rather than adding a second. If the recorder instead refuses the round, produces no result, or reports an incomplete round, hold the flow: do not claim that the round was recorded, re-run the reviewer, or prepare the next gate.

3. **Commit the recorded round.** Invoke `«state.commit»(slug)`.
   **Done when** the recorded round is durably committed and the tree is clean for that entity. Until it is, the entity is absent from `status --next` and refused by `gate prepare`: an empty scheduler here means this step is unfinished, never that the run is over.

4. **Re-run the reviewer on the corrected state.** Ask that reviewer to re-review the updated entity state, never to validate its own fix work. Re-run the kept-alive reviewer through the `«addressable-worker»` capability step 1 used, when it remains addressable and reuse conditions pass; fresh-dispatch a reviewer only when it is not. On cycle 3 run no reviewer: escalate to the human instead of a fourth round.
   **Done when** this cycle has a reviewer verdict — or, on cycle 3, the escalation to the human is recorded and the flow has stopped here.

5. **Re-enter the gate.** Invoke `Skill(skill="spacedock:fo-gate-lifecycle")` and run `«gate.lifecycle»` for the updated stage. A `needs-preparation` row is work to perform, not a stopping condition. This step is reached only from step 4's verdict: a `needs-preparation` row observed before that is step 3's commit landing, not a gate to prepare. The rejected report alone can surface one, and preparing on it would hand the captain a gate over work no reviewer has seen since the rejection.
   **Done when** exactly one fresh open gate has been prepared and presented, and the flow stops — without resolving or applying it, changing terminal state, or dispatching a successor.

The FO owns the shared correction-round section and writes it under `«write.classify»`: worktree-side when `worktree:` is set, main-side otherwise. The generic recorder owns only the immutable room and pointer bytes.

## Workflow-defined correction-round projection

The active workflow owns whether a projection exists and its exact grammar. This skill transports the selected projection without supplying defaults or interpreting its contents.
