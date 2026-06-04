---
title: FO should auto-continue from completed non-gated stages to next-stage dispatch
status: ideation
source: captain (2026-06-04) — FO stopped after implementation reporting instead of immediately advancing to validation; AI-engineer review found the current contract implies but does not enforce this lifecycle invariant
score: "0.32"
started: 2026-06-04T15:05:37Z
completed:
verdict:
worktree:
issue:
id: wmn2x3k7j0fjshvdz126ray3
---

The first officer should not stop after an implementation worker completes and the implementation stage report is filed. For a non-gated, non-terminal stage, the FO should continue the workflow lifecycle: verify the stage report, advance to the next stage, and dispatch the next worker. In the dev workflow, that means implementation completion immediately advances to validation and dispatches an independent fresh validator unless a gate, terminal merge ceremony, blocker, or captain decision interrupts.

## Problem

The existing contract implies this behavior but does not make it obvious or enforceable. `docs/dev/README.md` describes implementation output as ready for independent verification and validation as the verifier stage. The shared first-officer core describes completion handling and says fresh dispatch should run `status --next` and dispatch the next stage. However, it does not explicitly forbid a captain-facing stop after a non-gated, non-terminal stage report, and it does not name implementation-to-validation as a mandatory auto-continuation edge.

This allowed the FO to report implementation completion for `pi-stage-dispatch-uses-build-artifact` without immediately advancing to validation and dispatching the validator.

## Desired outcome

Make the lifecycle invariant explicit, testable, and hard to miss:

- After a completed non-gated, non-terminal stage, the FO must not stop with a completion-only report.
- The FO must advance to the next stage and dispatch it before ending the turn, unless a gate, terminal merge/mod flow, blocker, or captain decision is required.
- For `validation` stages marked `fresh: true`, the FO must dispatch a fresh independent validator rather than reuse the implementation worker.
- The ideation stage for this entity must define at least one concrete runtime/test scenario that would fail if the FO stops after implementation completion.

## Acceptance criteria

**AC-1 - The FO shared contract explicitly requires auto-continuation after non-gated, non-terminal stage completion.**
Verified by: a skill/integration invariant test over `skills/first-officer/references/first-officer-shared-core.md` that fails unless the Completion/Gates section forbids stopping after a completion-only implementation report and requires advancing/dispatching the next stage before ending the turn.

**AC-2 - The dev workflow docs state implementation report filing is not a stopping point.**
Verified by: a docs/invariant test or validation check over `docs/dev/README.md` requiring wording that implementation completion routes immediately to independent validation dispatch unless gated/blocked.

**AC-3 - Pi runtime guidance preserves the same lifecycle after `pi-subagents` completion.**
Verified by: a skill/integration invariant test over `skills/first-officer/references/pi-first-officer-runtime.md` requiring the parent to continue shared lifecycle after verifying the Pi subagent result/stage report, and not final-response unless gated, terminal, blocked, or awaiting captain decision.

**AC-4 - Ideation defines a concrete happy-path runtime/test scenario.**
Verified by: this entity's ideation section containing a scenario definition with fixture shape, initial state, FO prompt, expected durable state, and a negative case where stopping at implementation report fails the assertion.

**AC-5 - A failable test scenario catches the regression.**
Verified by: an implementation test, live/shared scenario, or fixture-level assertion that starts from an implementation-ready entity, observes/report-files implementation completion, and fails unless the FO advances to validation and dispatches/runs a fresh validator or presents the validation gate as appropriate.

## Notes

An AI-engineer review found this should not be docs-only. The final solution should likely include skill contract wording, docs/invariant tests, and either a runtime scenario or a product helper/next-action command so the FO has a concrete state-machine action to follow.
