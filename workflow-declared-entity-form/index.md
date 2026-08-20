---
id: qpa28ewtx6b9f50b4e3qcbqs
title: Let a workflow declare its entity form instead of deciding per filing
status: backlog
source: "Captain, 2026-08-20, while triaging issue #739: there is no filing-time default for entity form, and the per-filing judgment is made before anyone knows whether the entity will be gated."
started:
completed:
verdict:
score:
worktree:
issue:
---

Let a workflow declare that its entities carry their own artifacts, so filing does not depend on a per-entity guess.

## Problem

Entity form is chosen once, at filing, by whoever runs `spacedock new`. Flat is the default by omission: `--folder` is an explicit opt-in flag, the workflow README frontmatter carries no form key, and no entity-form concept exists in the binary. The only guidance is one README line, "Use folder-form entities when reports or artifacts may accumulate beside the task."

That judgment is made at the point of least information. Whether artifacts accumulate beside an entity is not a property of the entity's subject. It follows from the workflow: a workflow whose stages are gated accumulates review rooms beside every entity it gates, and a workflow that keeps reports beside its entities does so for all of them. The filer is asked to predict a fact the workflow already determines.

The consequence is recorded in `gate-room-form-breaks-on-conversion`: a gated flat entity becomes a hybrid, `<slug>.md` beside a `<slug>/` directory of review rooms. In this checkout every one of the 10 flat gated entities is in that state. That entity owns the gate-side fix. This one removes the guess.

## Proposed approach

Two surfaces, one decision.

A workflow frontmatter key declares the default form for entities filed into that workflow. `spacedock new` follows the declaration, and the existing `--folder` flag stays available for a per-entity override. A workflow that declares nothing keeps today's behavior.

A commission-time parameter sets that key when a workflow is created, phrased as the property the operator actually knows: whether the workflow's entities are likely to carry their own artifacts, including whether its stages use gate preparation for review. Commissioning already knows the stage list and which stages are gated, so it can propose the answer rather than ask blind.

Ideation owns the key's name and spelling, whether the default is derived from the declared stages or stated independently, and what an existing workflow with no key does.

## Out of scope

Do not change the gate ceremony, room layout, or retained-reference resolution — `gate-room-form-breaks-on-conversion` owns those. Do not migrate existing entities here. Do not remove the `--folder` flag.

## Expected surface and tolerance

Ideation sets this. The surfaces are the workflow frontmatter schema, `spacedock new`'s form selection, and the commission flow's parameter.

## Acceptance criteria

Ideation writes these. They must cover a workflow that declares the key, a workflow that declares nothing, and the per-entity override still winning.

## Test plan

Ideation writes this. Filing is cheap to exercise offline against both declarations.
