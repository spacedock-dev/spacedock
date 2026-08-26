---
id: 16yzsp1zf341v1x4rxfprm8h
title: The dispatch assignment carries the workflow-local entity-body schema the gate reads require
status: backlog
source: "Captain-filed spacedock-dev/spacedock#765, 2026-08-26: a worker wrote complete acceptance criteria under a level-1 heading, all checklist items DONE, gate prepare bound — then the mandatory post-prepare ac-scan failed with 'no ## Acceptance criteria section in this file'; the consumer grammar is valid but absent from the producer contract"
started:
completed:
verdict:
score:
worktree:
issue: "spacedock-dev/spacedock#765"
pr:
mod-block:
---

A generated stage assignment must carry enough workflow-local entity-body schema for the worker to produce a body the mandatory structured gate reads accept. Today the dispatch transports the seed, the stage definition, the checklist, and the report protocol — but not the Task Template's body grammar, and not the fact that `--ac-scan` requires the exact level-2 heading.

## Problem

{Ideation fills this in. Seeded from the issue: the failure arrives AFTER an open binding exists, forcing withdrawal, rebinding, and worker rework. The invariant to design for: producer and consumer share one declared grammar, without hard-coding the dev workflow's template into the universal ensign core — non-dev workflows keep their own schema. Related mechanism family: embed-stage-report-protocol-in-dispatch (t4, in flight) transports the REPORT grammar the same way this task transports the BODY grammar; keep the two composable, not duplicated.}

## Proposed approach

{Ideation fills this in. The issue names candidate mechanisms — transport the Task Template/schema, emit a canonical body skeleton, or another form that preserves workflow-local templates — and an acceptance boundary: a fresh seed completes ideation on the delivered assignment alone and passes both structured reads; an invalid shape fails BEFORE a binding is created, with a diagnostic naming the required shape.}

## Risk evidence

{Backlog: the issue's concrete reproduction (level-1 heading, three DONE items, --checklist green, --ac-scan non-zero) decides design should start.}

## Out of scope

The report-protocol embed (t4). The scanner's citation-source gap (#766, its own task).

## Expected surface and tolerance

{Backlog seed; ideation estimates with the production/proof split.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A fresh seed with no body headings enters ideation through the normal dispatch build path, completes using only the delivered assignment and installed ensign contract, and passes status --read --checklist plus --ac-scan.**
Verified by: {ideation refines — seed: the issue's acceptance boundary as a fixture journey; falsifying perturbation: an invalid acceptance-section shape must fail before an open binding exists, with the diagnostic naming the required workflow-local shape.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
