---
title: "Gates & decisions"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-07-25 19:04:19"
---

# Gates & decisions

A gate is the decision point at the end of a stage where nothing advances without your vote. When a stage declares a gate, the first officer stops after the worker completes, presents a review, and waits for you. It never self-approves.

Each call you make sharpens the bar, and the destination is delegation. When you are sufficiently confident in the workflow and the bar you set, hand over the conn and let the first officer drive multiple tasks with auto-approval:

```
you have the conn to drive toward your sprint goal, authorized to approve and
merge PR on CI green. use your judgement.
```

Until then, the gates are yours.

## What you see at a gate

A chat gate review has one concise evidence spine. The first officer emits it before recording either your decision or a decision made under delegated conn:

```
Capability/change: replace sleep-based waits with event polling.
Test and evidence: login is stable across 50 runs; AC-2 retry has no covering test.
Reviewed snapshot: Briefing `...` at digest `sha256:...`.
Findings: material: AC-2 cites a test file that does not exist.
Recommendation: revise to add the retry scenario.
Decision ask: approve to close, revise to bounce back, or hold at review.
```

Material findings are the ones that should move your vote; polish never blocks. The decision ask tells you concretely what each call does. Every acceptance criterion is cross-checked before the review reaches you; delegated authority does not hide the review.

The first officer may act only after an explicit grant in the active conversation, including a grant given later in that conversation. Durable gate state records the first officer as the decision renderer and its evidence reason; it does not quote or authenticate the grant's wording or scope. Keep any required chat provenance in the host's own audit system.

## How the review reaches you

Gate reviews appear in chat by default. A workflow or session can opt into a review provider that presents the same canonical Briefing as a blocking review and returns an exact retained Result. The scaffold prepares one gate room that binds the request authority, gate attempt, canonical Briefing, and fixed provider outputs. Binding the attempt freezes both the Briefing digest and the request digest. The first officer passes only that room to the provider. The provider must show the Briefing's question, every Artifact, and every recursively reached Reference at its recorded revision. The recorder derives the complete presentation association from the room; the first officer does not assemble it.

The provider owns its presentation transport and retained files. The `spacedock` binary verifies a direct binding Result through its nested `Resolution.by`; advisory output remains evidence and cannot close the gate through an adoption note. If the provider is missing or has the wrong version, the first officer names the remedy and returns to chat without launching it or creating retention files. After launch, a failure never falls back to chat.

## The three calls

- **Approve.** The decision is recorded first. A separate application step may then advance eligible work exactly once.
- **Redo with feedback.** You accept the direction but send concrete fixes back. The recorded decision is `revise`, and its reason says the direction is accepted.
- **Reject.** With a configured feedback target, the recorded decision is `revise`, its reason says the direction is rejected, and the work bounces to that owner. Without a feedback target, the decision is `hold` and the first officer stops for routing help.

Redo and reject differ only in whether you accept the direction; both carry your concrete asks so the next worker has something to act on. Nothing closes without its verdict on the record. Your call translates into the existing `approve`, `revise`, and `hold` record; automatic bounce applies only when a reviewer recommends `REJECTED` at a configured feedback gate.

After completion verification, a gate with no selected attempt remains `validating`. The first officer binds and commits the retained Briefing before presenting anything. That bind selects the current-stage gate attempt, letting startup distinguish work still validating, an open attempt awaiting the Captain, an approval awaiting nonterminal advance, and an approval awaiting merge. Approval to a terminal target is consumed before the existing merge and terminalization path begins.

Before the first officer shows a gate, it captures the exact bound Briefing identity and digest from committed entity state. A run without decision authority stops with that attempt open: it writes no Resolution, consumes nothing, advances nothing, and dispatches nothing. After an authorized decision, it records and commits the Resolution before every route. Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark; the consumed descendant commit lands before ordinary successor dispatch. Revise routes feedback after its close commit, and hold stays at the gate. `gate validate` and `gate eligibility` remain optional diagnostics, not positive-path lifecycle steps.

The review itself stays concise: capability, evidence, reviewed snapshot, findings, recommendation, and decision ask. The entity, spec, and package remain linked references rather than replacing that review with raw artifacts.

## Rejections

Rejections bounce automatically: the findings route back, the work is redone, and the reviewer re-runs; no stop at your desk. The gate reaches you only when the work passes review, or after **three failed rounds**, when the call returns to you instead of bouncing again. Every round is on the record in the item's file.

A useful rejection to type at a gate: "send it back unless this now needs reframing".

## Reviews beyond validation

A typical validation stage already covers code review: the work is checked against your acceptance criteria, with the rejection loop behind it. Adversarial review is built in as well: a `fresh: true` validation stage is exactly that, and high-stakes changes can also get a detached, out-of-workflow pass.

Adversarial review differs from validation in what it distrusts: validation checks the work; the adversarial pass checks the validation. It is read-only and tries to refute the result by constructing an adversarial edit the deliverable's own tests should catch, then confirming they do. A test that stays green under an edit that breaks the claim is a hole validation cannot see on its own. "Refuted nothing material" is a valid recorded outcome, and material findings route back through the rejection loop.

The workflow is flexible beyond that: you can add conditional, lens-specific reviews of your own, like checking that the documentation was updated when a change affects end users.

## Where to go next

- [Operate a workflow](../../running-workflows/operating/) covers answering gates in the day-to-day loop.

## Sitemap

- [The operating model](../operating-model/index.md)
- [Workflows & entities](../workflows-and-entities/index.md)
- [The stage lifecycle](../stage-lifecycle/index.md)
