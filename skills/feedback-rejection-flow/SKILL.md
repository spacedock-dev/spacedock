---
name: feedback-rejection-flow
description: "First-officer feedback-rejection routing — the correction-round record, reuse-or-fresh routing of findings to the `feedback-to` target, the reviewer re-run, gate re-entry. Invoke at the rejection-handling point when a feedback gate recommends REJECTED or the captain rejects at a feedback-to stage."
user-invocable: false
---

# Feedback Rejection Flow

This skill loads at the rejection-handling point; the rejection DETECTION, the `### Feedback Cycles` write-scope rules, the reuse conditions, and the budget probe stay always-on in the FO contract, and this procedure references them by name.

## Feedback Rejection Flow

When a feedback stage recommends REJECTED:

1. Read the rejected stage's `feedback-to` target — the stage that receives the fix request, not the reviewer.
2. Append this round's entry to `### Feedback Cycles` in the entity body, in the format below.
3. Read the deviation and the AC-drift note from the two most recent entries. Past the declared tolerance — 2× unless the entity declares its own — or on a narrowed AC, record a design-reset decision (reconfirm, re-scope, park, or escalate) before any further round; no automatic re-dispatch while that decision is absent. On cycle 3, escalate to the human instead of another round.
4. Invoke `«context-budget»()`. If the old ensign is over budget or the source is unavailable, shut down and fresh-dispatch; if no probe is declared, proceed to reuse below.
5. Route findings back to the target stage in the same worktree using `«addressable-worker»` when the existing handle is addressable and reuse conditions pass; otherwise shut down and fresh-dispatch. The routed message must carry the concrete next-stage assignment and fix work, not just an acknowledgment request. Do not treat the immediate routing response as the new completion result: if the follow-up is on the entity's critical path, wait for the reused worker's next completion through `«completion-signal»` before advancing or shutting it down. Attribute completion by mailbox content, task path, or durable workflow state, not by narration.

   Findings are inputs to triage, not a fix list; a direct captain instruction is an order, not a finding. Ask the routed worker to record its triage as an advisory resolution on the round's briefing—never a binding one, because a round cannot advance status. Material findings are fixed; each correct-but-disproportionate finding is declined with its class, why it is not material, and what would promote it, using the deferred-risk fields validation already requires; a needs-decision finding is escalated. A finding neither fixed nor recorded is not triaged. Narrowing a value AC to pass a finding is a captain-owned design reset recorded through a binding gate attempt, not a worker disposition.

   After the reviewer and worker entries are complete, invoke `${SPACEDOCK_BIN:-spacedock} gate record <entity> --round <stage>/<cycle> --briefing <path>/briefing.json --log <path>/briefing.review.jsonl --feedback-cycle <file>`; the round records no gate application and does not change status. Omit `--feedback-cycle` for a complete no-findings round.
6. Re-run the reviewer after fixes. When the existing reviewer remains addressable and reuse conditions pass, re-run the kept-alive reviewer through the same `«addressable-worker»` capability used for feedback routing; the message must ask that reviewer to re-review the updated entity state, not validate its own fix work. Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail.
7. Re-enter the normal gate flow with the updated result.

The FO owns `### Feedback Cycles` and writes it under the eagerly loaded first-officer `«write.classify»` scope: worktree-side when `worktree:` is set, main-side otherwise.

## Feedback Cycles entry

One line per correction round — a cross-stage bounce or an in-stage review round, one section, one format:

    - Cycle {N}: {verdict} — {reviewer/loop}; surface {actuals} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}

Take the actuals from the one-liner the workflow documents for its own surface unit. Deviation is the actuals over the estimate the ideation gate approved, never over the prior round — every round of a runaway passes against the prior round.
