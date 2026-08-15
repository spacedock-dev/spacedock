---
id: 42chs9dh7nq22f8at4szvbxp
title: Tune the dev task template for gated stages
status: backlog
source: Captain direction, 2026-08-13
sprint-readiness: ready
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

Estimate net LOC change: +30, across 1 file.

## Acceptance criteria

**AC-1 — The template prompts for every decision input required by gated stages without restating their full stage contracts.**
Verified by: compare each gated stage's `Gate content` row with the completed template and identify the corresponding prompt.

**AC-2 — Chosen direction appears only in the ideation-stage prompt.**
Verified by: backlog and validation prompts retain their own stage-specific evidence fields and do not request a chosen direction.

**AC-3 — The template uses the compact LOC estimate form.**
Verified by: the rendered template requests one explicit signed estimate, such as `+60` or `-25`, with tolerance stated separately.

**AC-4 — Existing task files and gate lifecycle semantics remain unchanged.**
Verified by: the diff is confined to `docs/dev/README.md` and workflow validation passes.

**AC-5 — Repeated clarification routes according to snapshot currency.**
Verified by: a behavior exercise shows evidence changes withdraw and reprepare the gate, while presentation-only clarifications preserve the bound snapshot and improve the stage-specific `Gate content` for future gates.

## Test plan

Run workflow validation and a one-off gated-stage coverage comparison. Do not add a standing prose-grep test.

## Captain-directed scope addition (2026-08-15)

Fold stacked-PR support into the pr-merge mod (docs/dev/_mods/pr-merge.md), from the 0.27 stack experience:

- The mod's back half (MERGED detection, sentinel, merge guard) worked unchanged on stacked PRs and needs no edit - encode that as an explicit statement, not an accident.
- The front half must gain a stacked mode: the FO builds template-conformant PR bodies BEFORE `gh stack submit` (or patches them immediately after via the REST API - the GraphQL path chokes on deprecated projectCards), records the per-layer candidate SHA, and accepts a stack-sibling base as valid where the mod today assumes the trunk.
- Body edits on existing PRs must use `gh api --method PATCH repos/{repo}/pulls/{n}` (observed working) rather than `gh pr edit` (observed failing).
- Live-lane economics note for the mod prose: required checks fire on the bottom and top layers; the top's tree is the composition, so one live approval covers the stack.
- Evidence base: PRs #699-#710, stack #707, and this session's ceremony record.

Design constraint: the mod is blocked-product (dispatched-worker territory); the README template is FO-owned process. Ideation designs both edits; implementation is dispatched on top of the current stack (worktree based on the stack tip) and lands as the next stack layer.
