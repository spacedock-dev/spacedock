---
title: "Exception paths (Degraded Mode, Break-Glass, dead-ensign detail) split out of claude-fo-dispatch.md into failure-triggered deferred modules"
status: ideation
group: tooling
source: "fable-token-trim-scout analysis 2026-07-02: Degraded Mode (claude-fo-dispatch.md:87-117, ~780 tok), the Break-Glass manual-dispatch template (:42-52, ~520), and Context-Budget/Dead-Ensign detail (:128-149, ~650) load at first dispatch but fire only on failure (second dispatch failure, helper non-zero exit, budget-fail) — ~1.6-1.9k tok on every dispatching session's happy path. #457 just shipped the same non-user-invocable-skill pattern for adapter-less deferred modules."
id: 41cfak9bgwtpa01m1z4qkprq
started: 2026-07-02T03:02:51Z
---

## Problem
The Claude dispatch module front-loads its exception machinery: every session that dispatches even one worker carries the full Degraded Mode contract, the Break-Glass template, and the dead-ensign/budget-failure detail, though the happy path never fires them.

## Desired direction (for ideation to refine)
Move each exception body into a non-user-invocable skill (the #457 pattern) loaded at its trigger; keep resident only the trigger lines and any pre-invocation guards (the 2-line context-budget probe invocation stays — reuse-condition-0 runs on the happy path; the Awaiting Completion section stays — it guards the common wait path). Design must weigh the named risk: exception time is when the FO is least reliable, so the resident trigger must carry the first action, not just a pointer.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- claude-fo-dispatch.md's resident size drops by a measured amount (~1.5k+ tok) with the exception bodies loading only at their triggers.
- Each trigger line names its skill and the first action; contractlint reference-closure stays green.
- Degraded-mode and break-glass behavior verified equivalent post-split (touches skills/**, so claude-live gates the merge).
