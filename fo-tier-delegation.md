---
title: Model-aware tier self-identification + level-3 delegation routing
status: ideation
source: 'FO shaping 0205 (2026-06-15, this session) — for a weak FO to be safe by construction it must self-identify its tier and route every judgment call to a stronger standing teammate, never adjudicate alone. No such mechanism exists. The gate-half (route all gate verdicts to level-3 when the FO is Haiku and the workflow has gate:true stages) is the no-2y startable slice; the full seam visibility composes with 2y. 0205 layered-FO.'
started: 2026-06-16T02:21:21Z
completed:
verdict:
score: 0.5
worktree:
issue:
sprint: 0205-layered-fo
sprint-readiness: ready
id: 72r2x0nnvx9az9x1adf08svq
---

The mechanism that makes a weak FO safe by construction: it knows it is weak and escalates structurally, never by luck.

## Problem

A Haiku FO has no way to know it should not make a gate verdict, a scope call, or a conflict-recovery decision. Without a self-identification + routing mechanism, safety depends on the model happening to escalate. The fix is structural: the FO self-identifies its tier and a named routing table sends every judgment category to a standing level-3 teammate on a stronger model.

## Proposed approach

{Ideation designs. Candidate shape: «fo.tier»() resolves the FO tier at boot (from the runtime model name or a launcher flag); a named routing table maps judgment categories (gate verdicts, design/scope, feedback-cycle-3 escalation, model-mismatch reuse, conflict recovery, teardown health) to a level-3 route; a standing level-3 teammate (stronger model) is spawned on demand; a boot gate routes ALL gate verdicts to level-3 when the FO is Haiku and the workflow has gate:true stages. The gate-half is the no-2y startable slice.}

## Out of scope

{The full mechanical/judgment seam visibility (composes with 2y); building the binary verbs (follows the spike); codex/pi substrates.}

## Acceptance criteria

**AC-1 — A Haiku FO self-identifies its tier and a gate verdict is made by a level-3 teammate, never by the Haiku FO alone, with the level-3-made verdict on the durable record.**
Verified by: {a live drive (coordinate with haiku-loop-spike) where the Haiku FO routes a gate verdict to a standing level-3 teammate; graded on the durable workflow state showing the verdict was made by level-3, never transcript phrasing.}

## Test plan

{Ideation fills; coordinate the live-drive proof with haiku-loop-spike. The riskiest mechanism — does tier self-ID + routing actually hold in a live Haiku drive — is exercised in that spike.}
