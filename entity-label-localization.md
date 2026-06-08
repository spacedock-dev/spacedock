---
id: jm0vqtx3j5vfw106kg0nz87b
title: Localize the operating voice to the workflow's entity-label (keep the shared contract generic)
status: backlog
source: "captain (2026-06-08, this session). Captain-facing + workflow-specific prose says generic \"entity\"; it should read as the workflow's declared label (this workflow: \"task\"). Principle: localize the OPERATING VOICE, keep the SHARED ABSTRACTION generic — do NOT rename entity→task in the shared contract (it serves any workflow: ticket/story/experiment)."
started:
completed:
verdict:
score:
worktree:
issue:
---

Make captain-facing and workflow-specific prose speak the workflow's own noun. A workflow that declares `entity-label: task` should present gates, dispatch packages, sprint docs, and status reports in terms of "task", while the shared contract keeps the generic "entity" that serves any workflow (ticket / story / experiment).

## Problem

The README already declares `entity-label` / `entity-label-plural` (here: `task` / `tasks`), and the FO already reads them at boot — but the operating voice still says generic "entity" everywhere the captain reads it. The fix is not to rename the abstraction; it is to *use* the label that already exists, only at the layer the human sees.

## Proposed approach

**Layer 1 (do first).** The FO/Commander resolves the README `entity-label` / `entity-label-plural` in captain-facing + workflow-specific prose: gate presentations, the Commander dispatch package, sprint docs, status reports. Add `{entity-label}` placeholders to the present-gate and commander-dispatch templates so the discipline is "use the label you already read," not "remember to localize."

**Layer 2 (optional follow-up).** `spacedock dispatch build` / `status` substitute the label so localization is automatic, not discipline-dependent.

**Graduation.** Bake Layer 1 into the `spacedock:roadmap` skill graduation so the commander-dispatch template ships with `{entity-label}` and this never recurs in a new sprint.

**Guardrail (load-bearing).** The shared contract — `first-officer-shared-core.md`, `ensign-shared-core.md`, the present-gate/commander-dispatch *abstraction* — keeps generic "entity". Localization happens only in the human-facing rendering layer.

## Out of scope

- Renaming `entity` → `task` (or any label) in the shared first-officer / ensign contracts. The abstraction stays generic.
- Per-workflow forks of the shared skills.

## Acceptance criteria

Ideation/implementation fills in. Sketch:

- The present-gate and commander-dispatch templates carry `{entity-label}` placeholders that resolve from the README at render time (verified by rendering a gate/dispatch for a workflow whose `entity-label` ≠ "entity" and observing the declared label in the output — not by grepping the template for the placeholder).
- The shared contract files still read generic "entity" (verified by their unchanged abstraction text — the guardrail is that this AC and the one above do not collide).
- The `spacedock:roadmap` graduation emits a commander-dispatch template carrying the placeholder (verified by graduating a sample roadmap and inspecting the produced template).

## Test plan

Ideation/implementation fills in. The behavioral half (a rendered gate/dispatch shows the declared label) must be proven by rendering, not by a substring match over the template (proof-policy). Layer 2, if taken, is a `dispatch build` / `status` golden over a fixture workflow with a non-default `entity-label`.
