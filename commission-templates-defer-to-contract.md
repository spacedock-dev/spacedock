---
id: 48wktz0b51941chr9c8kfask
title: Restructure the shipped commission templates — lead-with-the-end + defer universal rules to the FO/ensign contract + a workflow-specific-rules slot
status: ideation
source: captain (2026-06-14) — apply the dev-README slim ethos (lead-with-the-end + defer-to-contract; the README "move-3" work) to the SHIPPED commission templates. A templates assessment this session found development.md / experiment.md / refinement.md restate universal rules inline (proof discipline in development.md's "Recommended practices" ~L111-137; the Inputs/Outputs/Good/Bad stage-semantics block repeated per stage across all three). This is the template half descoped from rzp. Adjacent to but distinct from ey (which ports the proof-policy RULE into the contract + templates); this task makes the templates DEFER universal rules to the contract and keep only workflow-specific rules. 0.20.4 read-cost theme (every commissioned workflow inherits the scaffolding).
started: 2026-06-15T05:35:19Z
completed:
verdict:
score: 0.33
worktree:
issue:
sprint: 0204-structured-reads
sprint-readiness: ready
---

The three shipped commission templates duplicate universal rules that already live in the FO/ensign operating contract. Restructure them the way the dev README is being slimmed: lead with the workflow's outcome, defer universal rules to the contract, and keep only each shape's workflow-specific rules in a dedicated slot. Every new commissioned workflow then inherits leaner, non-stale scaffolding.

## Problem

`skills/commission/references/templates/development.md`, `experiment.md`, and `refinement.md` each restate rules that the FO/ensign contract already owns:
- proof discipline (development.md's "Recommended practices" ~L111-137: external-proof ACs, live-scenario, detached audit) parallels the contract's "prefer a code gate over a prose-only rule" and the ensign's "prove by exercising, not re-reading";
- stage semantics (the Inputs / Outputs / Good / Bad bullets) repeated per stage across all three;
- AC-must-be-testable and lead-with-the-end, both already universal.

The duplication means every commissioned workflow carries the long restated prose, and the template copies drift independently of the contract they paraphrase.

## Proposed approach

{Ideation fills. Candidate direction from the assessment, to evaluate not prescribe:
(a) lead each template with the workflow's outcome/value, not "what this is";
(b) a "Universal operating contract" section that defers to first-officer-shared-core.md + ensign-shared-core.md for stage semantics + proof discipline, rather than restating them;
(c) a "Workflow-specific rules" slot carrying only the 3-5 rules unique to each shape (development: the repo-mutation worktree layer + pr-merge mod offer; experiment / refinement: their own specifics);
(d) compact per-stage prose to one line plus any workflow-specific stage note.
Sequencing: this composes with ey — ey lands the proof-policy rule IN the contract; this makes the templates DEFER to it. Resolve overlap on development.md with ey (one edits the rule in, the other defers to it).}

## Out of scope

{Ideation fills. Likely: the contract-side proof-policy port (that is ey); codebase-specific rules that legitimately stay per-workflow; the dev README itself (FO-direct move-3).}

## Acceptance criteria

Each AC names a property of the finished scaffolding, not a stage action, and how it is verified — proof is behavioral or binds an independent source, never a prose-grep over the template.

**AC-1 — A workflow commissioned from a restructured template still inherits the universal rules (stage semantics, proof discipline) from the contract.** Behavior is preserved; the rules are inherited, not restated.
Verified by: {a live commission from each restructured template, then a boot/behavior check that the FO/ensign apply the universal rules — not a substring match over the template. Ideation pins the exercise.}

**AC-2 — Each template carries a workflow-specific-rules slot with only that shape's unique rules.** The universal rules are absent from the template body (deferred), the shape-specific ones present.
Verified by: {a structural check binding the template's rule set to an independent source (the contract carries the universal ones; the template carries the shape-specific ones) — two sets that can diverge, not a tautological presence assert. Ideation pins it.}

**AC-3 — Each template leads with the workflow's outcome before its mechanics.**
Verified by: {ideation pins the check (structural opening-section assertion, or part of the commission behavioral exercise).}

## Test plan

{Ideation fills. Note this is prose-scaffolding: the load-bearing proof is a live commission + behavior check that universal rules survive the defer, plus an independent-source structural check for the workflow-specific slot — never a substring match over the template prose.}
