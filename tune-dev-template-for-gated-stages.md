---
id: 42chs9dh7nq22f8at4szvbxp
title: Tune the dev task template for gated stages
status: backlog
source: Captain direction, 2026-08-13
sprint-readiness: defer
score: 0.8
---

Tune the reusable dev task template so task authors supply the decision evidence required by each gated stage without making task files verbose.

## Problem

The stage definitions declare authoritative `Gate content`, but the reusable task template does not consistently prompt authors for those inputs. The LOC estimate was also more verbose than needed.

## Proposed approach

Align the template with the backlog, ideation, and validation gate-content contracts. Keep prompts compact and avoid duplicating stage instructions.
Treat chosen direction as ideation-specific, not as a generic task-template field. Backlog prompts for scope and required proof; validation prompts for results, evidence, findings, and readiness.
Treat repeated clarification on the same gate as evidence that its decision presentation is incomplete. Identify the missing decision input: if task evidence changed, withdraw the stale gate and prepare a new snapshot; if only presentation guidance changed, update the stage-specific `Gate content` for future gates. Never silently alter a bound snapshot.

## Out of scope

Do not change gate authority, lifecycle behavior, product code, or live-runtime grading.

## Expected surface and tolerance

Estimate net LOC change: +-30, across 1 file.

## Acceptance criteria

**AC-1 — The template prompts for every decision input required by gated stages without restating their full stage contracts.**
Verified by: compare each gated stage's `Gate content` row with the completed template and identify the corresponding prompt.

**AC-2 — Chosen direction appears only in the ideation-stage prompt.**
Verified by: backlog and validation prompts retain their own stage-specific evidence fields and do not request a chosen direction.

**AC-3 — The template uses the compact LOC estimate form.**
Verified by: the rendered template contains `Estimate net LOC change: +-{NNN}, across {M} files.`

**AC-4 — Existing task files and gate lifecycle semantics remain unchanged.**
Verified by: the diff is confined to `docs/dev/README.md` and workflow validation passes.

**AC-5 — Repeated clarification routes according to snapshot currency.**
Verified by: a behavior exercise shows evidence changes withdraw and reprepare the gate, while presentation-only clarifications preserve the bound snapshot and improve the stage-specific `Gate content` for future gates.

## Test plan

Run workflow validation and a one-off gated-stage coverage comparison. Do not add a standing prose-grep test.
