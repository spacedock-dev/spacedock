---
id: b8ewpvd65epkckvng0n38809
title: Workflow Fit Gate before FO entity creation
status: backlog
source: "Captain draft and directive, 2026-08-16: the FO tends to add stuff into existing workflows and be ceremonial about things not supposed to be there. Session evidence 2026-08-14/15: the banned doc-only journey entity, its mechanism-without-value reshape, and the release-cut task question."
started:
completed:
verdict:
score:
worktree:
issue:
---

Amend the shipped FO write core (skills/first-officer/references/fo-write-core.md) with an admissibility gate ahead of the "FO may write new entity files" rule. Captain's draft, the seed text:

## Workflow Fit Gate

Before creating or materially reclassifying an entity, the FO must verify the work fits the commissioned workflow's subject and value model. Write authorization is not workflow-fit authorization.

A new entity belongs only when it produces or validates a deliverable the workflow exists to track, using that workflow's own `entity-type`, README purpose, stage outputs, and acceptance/proof policy.

Do not file FO/process maintenance, workflow refits, split-root migration, debriefing, status reporting, cleanup of agent/session state, or operating-ledger work into a product/dev workflow unless the workflow README explicitly names that class as an executable deliverable. Record those in a debrief, reconciliation ledger, runbook, roadmap/planning doc, or a separate workflow/process track instead.

If fit is ambiguous, stop before `spacedock new` or `status --set` and ask the captain where the work should live.

And one line under New entity files: `spacedock new` is only the atomic creation mechanism after the Workflow Fit Gate passes; it does not decide whether the work belongs in this workflow.

Enforcement honesty, binding on ideation: this is contract prose executed by a model. The enforcement points are the FO asking the fit question BEFORE `new`, and the backlog gate where the captain catches misses. Do NOT design a committed prose-grep or a lint that reads this file - the write-core fixture precedent already records that a Go reimplementation proves table content, not FO obedience. A one-off falsifiable exercise at validation (replay a known-banned seed scenario against the amended contract) is legitimate evidence; a standing check is not.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in - including reconciliation with the dev README's existing line "if the task's only output is a decision, record it in the roadmap instead", which this gate generalizes; cross-reference, do not duplicate.}

## Out of scope

The dev workflow's own README (FO-owned process; any local mirror belongs to the deferred refit). New enforcement machinery.

## Expected surface and tolerance

Estimate net LOC change: +NN small, 1 shipped file (fo-write-core.md); contractlint reference closure unchanged.

## Acceptance criteria

**AC-1 - The shipped write core carries the fit gate ahead of the entity-creation rule, and the atomic-mechanism clarification under it.**
Verified by: the amended file plus contractlint green; falsifying edit - removing the gate section breaks the stated ordering.

**AC-2 - The gate can fail in exercise: a known-banned seed class, replayed against the amended contract at validation, is refused or routed to a non-workflow home.**
Verified by: one-off validation exercise recorded in the report, with the 2026-08-14 journey-entity filing as the replay scenario; never a committed test.

**AC-3 - The suite stays green.**
Verified by: go test ./internal/contractlint/ plain and -race.

## Test plan

Prose amendment plus reference-closure lint; the one-off replay exercise as the value evidence.
