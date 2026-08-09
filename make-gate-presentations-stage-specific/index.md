---
title: Make gate presentations stage-specific and omit empty result classes
status: backlog
id: krbaeb3resfpbh1qvnb65krf
score: 0.8
source: "Captain feedback on 2026-08-08 after repeated gate reviews rendered FAILED: None. Extended by the 2026-08-09 journey-friction audit: the same stage-specific presentation must show the exact selected artifact and references plus the executable next action, so one gate operation does not require filesystem searches or help retries."
issue:
sprint: durable-decisions
sprint-readiness: ready
started:
completed:
verdict:
worktree:
pr:
mod-block:
---

## Outcome

Each gate review presents the evidence that matters at its current stage. The review does not invent checklist rows for empty result classes.

## Problem

The generic presentation template displays placeholders for `DONE`, `SKIPPED`, and `FAILED`. First Officers repeatedly convert empty classes into rows such as `FAILED: None`.

The same template also assumes that every gate has a stage report. A backlog gate reviews a seed task, not a completed stage report. This mismatch caused a fabricated backlog checklist during task `0y` dispatch.

## Required direction

- A backlog review presents the outcome, scope boundary, and proof readiness of the seed task.
- An ideation review presents the chosen direction, risky-mechanism evidence, expected surface, and acceptance proof.
- A validation review presents real checklist results, findings, executed checks, acceptance evidence, and delivery readiness.
- The presentation omits each empty result class. It never writes `None`, `N/A`, or a zero-result placeholder row.
- Every stage retains one recommendation, the exact bound Briefing identity, a digest, and one concrete decision line.

## Acceptance criteria

**AC-1 (VALUE) - A captain receives a concise review that contains only decision-relevant evidence for the current gate stage.**
Verified by: controlled backlog, ideation, and validation presentations each contain their required stage evidence and exclude evidence that does not exist at that stage.

**AC-2 - Empty checklist result classes produce no presentation row.**
Verified by: a review with zero failed and zero skipped items contains no fabricated `FAILED` or `SKIPPED` row. The assessment still reports the numeric totals.

**AC-3 - The stage-specific forms preserve gate authority.**
Verified by: each presentation names the task and stage, gives one recommendation, identifies the bound Briefing and digest, and ends with one decision effect.

**AC-4 - The change does not duplicate adjacent gate tasks.**
Verified by: `w5` remains the owner of exact-digest reliability. Task `xx` remains the owner of zero-item dispatch checklist acceptance.

## Test plan

Use the existing gate-presentation path to produce one controlled review for backlog, ideation, and validation. Grade the visible review against the stage requirements.

Include a validation case with no failed or skipped items. Make sure that the assessment reports zero without a `None` or `N/A` row.

Run the relevant contract smoke checks and one live First Officer presentation. Do not add a prose-presence test or a second presentation renderer.
