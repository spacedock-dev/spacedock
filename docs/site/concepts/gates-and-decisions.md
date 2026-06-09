# Gates & decisions

A gate is the decision point at the end of a stage. Nothing advances past a gated stage without a decision — yours, or one you've delegated to an agent review.

## What a gate carries

Each gate carries a stage report: findings, verdicts, artifacts, and anomalies. You decide on the evidence, not the transcript. The record outlives the reviewer, so a bad result can be traced back to the call that caused it.

At a gate, the first officer pauses and presents the stage report. You make one of three calls:

- **Approve** — the work meets the bar; the entity advances.
- **Redo with feedback** — the work needs revision; it bounces back to an earlier stage with your concrete feedback.
- **Reject** — the work does not meet the bar.

Some gates wait on you. Others resolve through a delegated agent review, when you've handed the call to an agent.

## Feedback cycles and the loop cap

Rejected work bounces back to an earlier stage for revision. To prevent loops, a hard cap limits how many times an entity can bounce. On the third consecutive rejection, the first officer escalates to you instead of auto-bouncing a fourth time — the loop never runs unbounded.

## The detached adversarial audit

For high-stakes surfaces — the front-door launcher, the state-mutation guards, the shipped scaffolding, and the CI/release machinery — a passing validation is necessary but not sufficient. Before merging such a change, a read-only adversarial audit runs on a separate throwaway checkout. The auditor tries to *refute* the validation: it constructs an edit that breaks the claim and confirms the deliverable's own tests catch it. A test that stays green under an edit that breaks the claim is a hole.

This catches the class of bug where a test passes but would also pass on a broken future edit — which validation, trusting its own green suite, cannot see. "Refuted nothing material" is itself a valid, recorded outcome.
