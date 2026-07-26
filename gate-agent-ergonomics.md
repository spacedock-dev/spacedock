---
title: Make recorded gate operation self-guiding for First Officers
status: backlog
source: "Durable-decisions sprint dogfood: manual 0c and xb gate/round operation, 2026-07-24."
score: 1.0
id: skwchfe30ac6ntr63j1g0txj
sprint: durable-decisions
started: 2026-07-26T12:36:08Z
---

The recorded gate model is trustworthy but still makes the First Officer construct and
remember too much. During dogfood, the FO manually assembled Briefing ids, room paths,
relative artifact URIs, SHA-256 revisions, routing, request metadata, lifecycle order,
and Feedback Cycle/round identity. A valid repository-relative Briefing path also failed
until replaced by an absolute path. The recorder rejected malformed state before
mutation, which is valuable; the ergonomic problem is missing scaffolding and
corrective guidance, not excessive validation.

The minimum valuable slice has two capabilities:

1. **Scaffold a prepared gate package from FO-authored judgment.** The FO supplies the
   readable gate review and selects existing entity/spec/evidence artifacts. Spacedock
   derives the current stage, next attempt/revision identity, room path, canonical
   entity-relative URIs, media types, digests, routing context, `briefing.json`, and
   `request.json`. Existing files remain references by URI plus revision; the scaffold
   does not copy them or create another recording system.
2. **Expose a machine-readable lifecycle next action.** A gate status read distinguishes
   `needs-briefing`, `open-awaiting-decision`, `closed-needs-commit`, `eligible`, and
   `consumed`; reports the exact gate/attempt/Briefing, blockers, and descendant-commit
   state; and emits the next valid command with required arguments. It guides the
   existing explicit authority mutations and commits rather than collapsing them into
   one opaque transaction.

The same design should account for four adjacent frictions without automatically
expanding the minimum implementation:

- normalize every CLI path from the caller's working directory before storing a
  canonical entity-relative reference;
- let advisory-round preparation derive the next global Feedback Cycle and consistent
  round/Briefing ids instead of synchronizing them manually;
- make validation errors name the rejected field/value and show the accepted canonical
  shape, especially triage dispositions and Feedback Cycle projections;
- project gate readiness through `status --boot`/`engage` so a complete validation
  report awaiting Captain action is distinct from validation still in progress.

Preserve the current integrity boundaries: exact package/result validation, separate
Briefing/decision/consume mutations, durable Git commit ancestry, advisory rounds that
cannot apply a gate, and provider-neutral Spacedock ownership.

Do not dispatch this task before the current durable-decisions lane lands and a new
prerelease is dogfooded by other First Officers. Their observed correctness and friction
feedback is an input to ideation; prefer removing demonstrated toil over predicting a
larger workflow engine.
