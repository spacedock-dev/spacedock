---
title: Gate-extract verbs — structured extraction of stage report, AC coverage, and reviewer findings (the verdict stays level-3)
status: ideation
source: 'FO shaping 0205 (2026-06-15, this session) — the FO gate prep (checklist DONE/SKIPPED/FAILED review, the AC cross-check, the Material/Polish reviewer tiers) is structured extraction over a stage-report file: deterministic, but done by the model today. A weak FO can do it reliably only as a binary call returning structured data; the verdict (judgment) stays level-3. No 2y dependency. 0205 layered-FO, no-2y parallel track.'
started: 2026-06-16T02:21:21Z
completed:
verdict:
score: 0.4
worktree:
issue:
sprint: 0205-layered-fo
sprint-readiness: ready
id: 6reqad9gff9wk544det3x4fj
---

Three deterministic extraction verbs so a weak FO assembles a gate from structured data instead of re-parsing prose. The verdict is NOT computed here; it routes to level-3.

## Problem

Gate preparation today is the model reading a stage report and a `## Acceptance criteria` section and extracting: the checklist accounting, each AC's evidence citation, and the reviewer findings tiered Material vs Polish. This is deterministic extraction a binary should own, so a Haiku FO gets the same structured input every time instead of re-deriving it (and mis-deriving it).

## Proposed approach

{Ideation designs. Candidate verbs: `spacedock gate checklist-extract` (DONE/SKIPPED/FAILED items + line ranges + a chosen-direction-required flag), `spacedock gate ac-scan` (per-AC evidence citations + a natural-place flag that tells level-3 which ACs this stage was the natural place to satisfy), `spacedock gate reviewer-parse` (Material/Polish tiers). Deterministic structured output over a stage-report file; behavior-first over committed fixtures.}

## Out of scope

{The verdict decision (approve/reject) — that is judgment and routes to level-3 via fo-tier-delegation; building the rest of the FO loop; 2y (no dependency).}

## Acceptance criteria

**AC-1 — Each verb emits correct structured data for a committed stage-report fixture: checklist accounting, per-AC evidence + natural-place flags, and Material/Polish reviewer tiers, asserted against the fixture's known contents.**
Verified by: {Go tests over committed stage-report fixtures (a passing report, a report with a FAILED item, a report with an unevidenced AC), asserting the extracted structure matches the fixture, never a prose match.}

## Test plan

{Ideation fills. Go fixture tests over stage-report samples; the verdict is explicitly NOT asserted (it is not computed).}
