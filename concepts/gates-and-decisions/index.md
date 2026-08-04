---
title: "Gates & decisions"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-04 14:53:17"
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
Reviewed snapshot: Briefing `...` at compact digest `sha256:1a2b3c4d…`.
Findings: material: AC-2 cites a test file that does not exist.
Recommendation: revise to add the retry scenario.
Decision ask: approve to close, revise to bounce back, or hold at review.
```

Material findings are the ones that should move your vote; polish never blocks. The decision ask tells you concretely what each call does. Every acceptance criterion is cross-checked before the review reaches you; delegated authority does not hide the review.

The first officer may act only after an explicit grant in the active conversation, including a grant given later in that conversation. Durable gate state records the first officer as the decision renderer and its evidence reason; it does not quote or authenticate the grant's wording or scope. Keep any required chat provenance in the host's own audit system.

## How the review reaches you

Gate reviews appear in chat by default. Before presenting one, the first officer commits newly authored selected sources and calls `gate prepare` with its question, primary Markdown review, exact concise summary, and References. Spacedock authors and binds a two-file recorder-ready room; the first officer commits that entity-owned room. The selected source payloads remain singular local Git objects rather than room copies.

A workflow or session may select a presentation override. Only after the prepare and bind commit, the selected channel receives exactly the emitted room through its declared interface as an opaque handoff. The generic Spacedock contract neither defines channel execution nor reconstructs authority outside that room. When the room is recorder-ready, the first officer passes that same room to `gate record --room`; the recorder recomputes request, Briefing, Result, inventory, and Git pins, derives the complete association in memory, and writes no `association.json`.

If a prepared room becomes stale before any provider decision, the first officer runs `gate withdraw` with a reason. Withdrawal is not approve, revise, or hold: it preserves the old room without a Resolution, provider evidence, application, or stage change. Cold boot reports `withdrawn-awaiting-prepare`; ordinary `gate prepare` then appends attempt N+1, which alone can receive the later Captain decision.

## The three calls

- **Approve.** The decision is recorded first. A separate application step may then advance eligible work exactly once.
- **Redo with feedback.** You accept the direction but send concrete fixes back. The recorded decision is `revise`, and its reason says the direction is accepted.
- **Reject.** With a configured feedback target, the recorded decision is `revise`, its reason says the direction is rejected, and the work bounces to that owner. Without a feedback target, the decision is `hold` and the first officer stops for routing help.

Only approve creates an application. Revise and hold are complete when their Resolution is recorded; workflow routing, not an application payload, handles feedback.

Redo and reject differ only in whether you accept the direction; both carry your concrete asks so the next worker has something to act on. Nothing closes without its verdict on the record. Your call translates into the existing `approve`, `revise`, and `hold` record; automatic bounce applies only when a reviewer recommends `REJECTED` at a configured feedback gate.

After completion verification, a gate with no selected attempt remains `validating`. The first officer binds and commits the retained Briefing before presenting anything. That bind selects the current-stage gate attempt, letting startup distinguish work still validating, an open attempt awaiting the Captain, a withdrawal awaiting replacement preparation, an approval awaiting nonterminal advance, and an approval awaiting merge. Approval to a terminal target is *held* at consume: `gate consume` spends nothing and writes no status — it leaves the application `pending` and returns the `approved-awaiting-merge` route, and `merge guard` is the sole terminal consumer. `merge guard` spends only with delivery proof: the `mod-block` is cleared in its own step first, then `application.state: consumed`, the terminal status, `verdict`, and `completed` move in one locked write, and the `pr` merge sentinel is retained through archive as durable delivery proof. A failed delivery that needs rework returns through the record stage's declared `feedback-to` as `superseded` (`merge guard --rework`); retryable delivery trouble leaves the approval pending and is safe to retry.

Before the first officer shows a gate, it captures the exact bound Briefing identity, digest, and emitted room in committed machine state, then presents a compact snapshot identity in prose. A run without decision authority stops with that attempt open: it writes no Resolution, consumes nothing, advances nothing, and dispatches nothing. After an authorized decision, the recorder itself commits and syncs the Resolution before every route (`--consume` folds the approve's consume into the same call). Approval then uses `gate consume`, which rechecks the retained request, Briefing, Git sources, and eligibility before atomically writing the successor stage and consumed mark. Until that first-entered working stage has a durable, complete Stage Report, `status --next` and boot name it as both `current` and `next`. Once the same-stage dispatch sets its worktree, every away-status `status --set`—backward or forward, even with `--force`—is refused until the report is durable. The consumed descendant commit therefore lands before one recoverable successor dispatch. Revise routes feedback after its close commit, and hold stays at the gate. `gate validate` and `gate eligibility` remain optional diagnostics, not positive-path lifecycle steps.

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
