---
name: feedback-rejection-flow
description: "First-officer feedback-rejection routing — the correction-round record, the triage block that rides the routed findings, reuse-or-fresh routing to the `feedback-to` target, the reviewer re-run, gate re-entry. Invoke at the rejection-handling point when a feedback gate recommends REJECTED or the captain rejects at a feedback-to stage."
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
5. Route findings back to the target stage in the same worktree using `«addressable-worker»` when the existing handle is addressable and reuse conditions pass; otherwise shut down and fresh-dispatch. Head the feedback context with the finding-triage block below, then the findings, so the routed worker gets the rule and the findings it governs in one packet. The routed message must carry the concrete next-stage assignment and fix work, not just an acknowledgment request. Do not treat the immediate routing response as the new completion result: if the follow-up is on the entity's critical path, wait for the reused worker's next completion through `«completion-signal»` before advancing or shutting it down. Attribute completion by mailbox content, task path, or durable workflow state, not by narration.
6. Re-run the reviewer after fixes. When the existing reviewer remains addressable and reuse conditions pass, re-run the kept-alive reviewer through the same `«addressable-worker»` capability used for feedback routing; the message must ask that reviewer to re-review the updated entity state, not validate its own fix work. Fresh-dispatch the reviewer only when the existing reviewer is no longer addressable or reuse conditions fail.
7. Re-enter the normal gate flow with the updated result.

The FO owns `### Feedback Cycles` and writes it under the eagerly loaded first-officer `«write.classify»` scope: worktree-side when `worktree:` is set, main-side otherwise.

## Feedback Cycles entry

One line per correction round — a cross-stage bounce or an in-stage review round, one section, one format:

    - Cycle {N}: {verdict} — {reviewer/loop}; surface {actuals} vs estimate {declared} ({P}%); findings {none | {F} fixed, {D} declined: <ref · class · why not material · promotes when>}; AC {unchanged | narrowed: <note>}

Take the actuals from the one-liner the workflow documents for its own surface unit. Deviation is the actuals over the estimate the ideation gate approved, never over the prior round — every round of a runaway passes against the prior round. Record an entry whenever a round produces a disposition, fixes or declines, not only when it triggers another round. `findings none` says no finding arrived; an all-declines round records `0 fixed` and names each decline, so "nothing was found" and "everything found was declined" never read alike.

## Finding-triage block

Head the routed feedback context with this, verbatim, above the findings:

    Review findings are inputs to triage, not a fix list. A **review finding** is the output of a review instance the active stage's definition declares — the panels, audits, or validators it names, carried to you in the dispatch packet — or feedback the first officer explicitly routes to you as a review outcome. A direct instruction from the captain is not a finding: it is an order to follow or a decision to seek, never something to triage away. Before changing anything in response to a finding, classify each against the workflow's declared finding-triage taxonomy, where one is declared, and this entity's own value acceptance criteria:

    - **Material** — breaks a value AC, or a declared non-negotiable boundary (safety, security, data-integrity, compatibility) reachable through the supported workflow. Fix it.
    - **Correct-but-disproportionate** (a deferred risk or polish) — right, but no value AC breaks and its trigger is outside the supported or promised workflow. Record a decline; do not fix it. The decline is your licensed disposition, not a dodge: name the finding, its class, why it is not material, and what would promote it to material.
    - **Needs decision** — a genuine product or compatibility fork. Escalate to the first officer; do not resolve it privately.

    Record the disposition — fixed as material, or declined and why — in this round's `### Feedback Cycles` entry so the gate sees it. A finding you neither fix nor record is not triaged.

    **Narrowing an acceptance criterion to make a finding or a rejection pass is not a licensed disposition.** Declining a disproportionate finding and narrowing the claim it targets are opposite moves under the same pressure: the first leaves the promised value intact and is yours to make; the second weakens that value and is a design-reset event requiring the captain's sign-off, recorded so it is captain-visible — never a task-internal edit.
