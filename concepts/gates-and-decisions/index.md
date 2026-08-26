---
title: "Gates & decisions"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-26 04:53:24"
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

Before it asks for a decision, the first officer reads the current stage's `Gate content` instruction and shows the evidence that instruction requests. `Gate content` is the authoritative presentation preference and override. Without it, the first officer selects concise evidence that supports the decision. This evidence comes from the stage definition, selected Artifact and References, current stage report, checklist, acceptance criteria, and findings. The first officer does not invent facts or show every source. It omits missing evidence and empty result or finding groups. The first officer shows the review before recording either your decision or a decision made under delegated conn:

```
Capability/change: replace sleep-based waits with event polling.
Test and evidence: login is stable across 50 runs; AC-2 retry has no covering test.
Reviewed snapshot: Briefing `...` at compact digest `sha256:1a2b3c4d…`.
Findings: material: AC-2 cites a test file that does not exist.
Recommendation: revise to add the retry scenario.
Decision ask: approve to close, revise to bounce back, or hold at review.
```

Material findings are the ones that should move your vote; polish never blocks. The decision ask tells you concretely what each call does. Every acceptance criterion is cross-checked before the review reaches you; delegated authority does not hide the review.

The first officer may act only after an explicit grant in the active conversation, including a grant given later in that conversation. Durable gate state records the first officer as the decision renderer, its evidence reason, and a citation of the grant it acted under — the grant's wording and where it was given. The record attributes; it does not authenticate the grant. Captain decisions carry no citation, so the two shapes stay distinguishable.

## How the review reaches you

Stable v1 gate reviews may appear in chat or Subspace. Before presenting one, the first officer commits newly authored selected sources and calls `gate prepare` with its question, primary Markdown review, exact concise summary, and References. Spacedock authors and binds a two-file recorder-ready room; the first officer commits that entity-owned room. The selected source payloads remain singular local Git objects rather than room copies.

Both presentation interfaces return semantic decision and reason input to the first officer. The first officer records that input through the same `gate record --decision` command; Subspace is not a second recorder and does not return Result or inventory files for Spacedock to ingest.

If a prepared room becomes stale before any Captain decision, the first officer runs `gate withdraw` with a reason. Withdrawal is not approve, revise, or hold: it preserves the old room without a Resolution, provider evidence, application, or stage change. Cold boot reports `withdrawn-awaiting-prepare`; ordinary `gate prepare` then appends attempt N+1, which alone can receive the later Captain decision.

### Explicitly outside v1

Provider-specific room-backed recording, `gate record --room`, Result or inventory ingestion, retained provider evidence, and provider package selection are not stable-v1 surfaces. Chat and Subspace presentation remain supported through the one semantic decision recorder.

## The three calls

- **Approve.** The decision is recorded first. A separate application step may then advance eligible work exactly once.
- **Redo with feedback.** You accept the direction but send concrete fixes back. The recorded decision is `revise`, and its reason says the direction is accepted.
- **Reject.** With a configured feedback target, the recorded decision is `revise`, its reason says the direction is rejected, and the work bounces to that owner. Without a feedback target, the decision is `hold` and the first officer stops for routing help.

Only approve creates an application. Revise and hold are complete when their Resolution is recorded; workflow routing, not an application payload, handles feedback.

Redo and reject differ only in whether you accept the direction; both carry your concrete asks so the next worker has something to act on. Nothing closes without its verdict on the record. Your call translates into the existing `approve`, `revise`, and `hold` record; automatic bounce applies only when a reviewer recommends `REJECTED` at a configured feedback gate.

After completion verification, a gate with no current-stage authority remains `validating` until the mechanical report checks pass. It then appears as `needs-preparation` on boot and every machine scheduler read. Engage performs semantic report review. A concrete `report-incomplete:` veto stops without mutation; otherwise the First Officer calls `gate prepare` exactly once with its question, Artifact, summary, and References, commits the emitted `state=open` binding, performs one structured checklist/AC read, and presents it without another boot projection. Open, withdrawn, stale, closed, and spent attempts retain their existing lifecycle routes. A gated stage the workflow marks `initial: true` has no prior stage to have written a report: there the committed, clean-in-HEAD seed is itself the proof, and the entity appears as `needs-preparation` directly. Engage prepares from the seed and skips the structured checklist/AC read.

Approval to a terminal target is *held* at consume: `gate consume` spends nothing and writes no status — it leaves the application `pending` and returns the `approved-awaiting-merge` route, and `merge guard` is the sole terminal consumer. `merge guard` spends only with delivery proof: the `mod-block` is cleared in its own step first, then `application.state: consumed`, the terminal status, `verdict`, and `completed` move in one locked write, and the `pr` merge sentinel is retained through archive as durable delivery proof. A failed delivery that needs rework returns through the record stage's declared `feedback-to` as `superseded` (`merge guard --rework`); retryable delivery trouble leaves the approval pending and is safe to retry.

Before the first officer shows a gate, it captures the exact bound Briefing identity, digest, and emitted room in committed machine state, then presents a compact snapshot identity in prose. A run without decision authority stops with that attempt open: it writes no Resolution, consumes nothing, advances nothing, and dispatches nothing. After an authorized decision, the recorder itself commits and syncs the Resolution before every route. `gate record --decision approve --consume` is the shortest approval path: close, sync, consume, and sync in one call. The supported standalone `gate consume` path rechecks the retained request, Briefing, Git sources, and authority before atomically writing the successor stage and consumed mark. Until that first-entered working stage has a durable, complete Stage Report, `status --next` and boot name it as both `current` and `next`. Once the same-stage dispatch sets its worktree, every away-status `status --set`—backward or forward, even with `--force`—is refused until the report is durable. The consumed descendant commit therefore lands before one recoverable successor dispatch. No separate state commit follows a successful record or consume write. Revise routes feedback after its self-synced close, and hold stays at the gate. `status` projects the next action from durable facts; the acting command reports any authoritative refusal.

After the nonterminal approval is consumed, it is ordinary stage history. The worker can write atomic terminal fields without `--force` after its report is durable. Pending terminal approval and unreadable or stale authority still fail closed. A consumed application that targets a terminal stage also fails closed.

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
