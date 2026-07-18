---
title: Persist gate approval while dispatch blockers remain
status: ideation
score: "0.80"
source: "Captain design feedback, 2026-07-13."
id: 3kd1x1gfxr8mdwzbmnwtjbw8
started: 2026-07-18T08:58:53Z
---

# Persist gate approval while dispatch blockers remain

## Problem

Spacedock currently conflates a human gate decision with immediate stage advancement and dispatch: approve, advance, and spawn are expected in one turn. That fails when a task is reviewable now but must not dispatch until another task lands. The First Officer must otherwise delay a valid review, dispatch against the wrong base, lose the approval, or invent hidden “approved but waiting” state.

Subspace already returns a structured `Resolution`, but the single-file review wrapper deletes its temporary package. The decision therefore does not become durable workflow state and cannot safely survive restarts or dependency waits.

Andrew's 2026-07-18 CMD cutover run exposed the same persistence gap on the
rejection path. Validation completed and correctly rejected the slice because
a production coordinator was missing; feedback routing correctly returned the
entity to implementation. The status surface then showed only
`implementation` at score `0.90`, erasing the legible fact that validation had
completed with a blocker and that the current implementation run was rework
caused by that rejected gate. This is production evidence for this filing, not
a separate task.

## Required capability

Persist a Subspace `Resolution` bound to the exact workflow entity, stage, review round, reviewed-content digest, and decision identity. Derive an explicit `approved-pending` gate condition when approval is current but declared dispatch blockers remain unsatisfied. This is computed gate/eligibility state, not another lifecycle stage.

Separate three concepts that are currently collapsed:

- durable gate decision: approve, revise, or hold, with provenance and reviewed digest;
- dispatch blockers: declared dependencies or predicates whose current state is queryable;
- dispatch eligibility: computed from current stage, non-stale approval, blocker satisfaction, and one-use consumption state.

Preserve rejected gate results under the same model. A lifecycle transition
back to a feedback target must not overwrite the completed rejection. The
reducer needs both the ordinary current stage and a durable route relationship
from the rejected gate attempt to the rework stage run.

The persisted representation must be workflow-owned and portable. Temporary Subspace package paths, machine-local pane/session identifiers, prompts, and private user data must not enter durable state.

## Scheduler behavior

1. Approval current and blocker present: retain the current stage, show `approved-pending`, and do not dispatch.
2. Blocker clears with reviewed content unchanged: atomically advance and dispatch without another human approval.
3. Reviewed content or gate-defining inputs change: mark the approval stale, keep the task non-dispatchable, and require a new gate.
4. Revise or hold: retain the current stage and remain non-dispatchable.
5. Status surfaces the decision, blocker set, blocker satisfaction, reviewed digest, staleness reason, and consumption state explicitly.
6. A current approval is reusable for exactly one successful advance/dispatch transition. Crash or restart reconciliation must distinguish no effect, exact prior success, and ambiguous execution without double dispatch.
7. A rejection routed through `feedback-to` retains the rejected gate result
   and projects the current lifecycle stage with explicit `feedback_rework`
   context until a later gate Resolution supersedes that cycle.

## Acceptance criteria

**AC-1** An approved blocked entity survives process restart and still reports the same durable approval, exact blocker, reviewed digest, and `approved-pending` condition without advancing or dispatching.

**AC-2** Clearing the final blocker with unchanged reviewed content causes exactly one stage advance and exactly one dispatch, without re-presenting the gate. Repeated scheduler passes do not redispatch or reuse the approval.

**AC-3** Any change to the reviewed artifact or other digest-bound gate input before blocker clearance marks the approval stale and produces zero advance/dispatch effects until a replacement Resolution is recorded.

**AC-4** Revise and hold Resolutions remain durable, visible, and non-dispatchable; blocker clearance cannot override them.

**AC-5** Dependency and scheduler failures fail closed. Missing, ambiguous, or unqueryable blocker state never appears as satisfied and never consumes approval.

**AC-6** Status text and JSON distinguish pending approval, stale approval, unsatisfied blockers, satisfied-but-not-yet-consumed approval, consumed approval, and ambiguous recovery.

**AC-7** The single-file Subspace review path persists the binding Resolution before deleting temporary review-package files. Durable records contain no temporary path, pane/session metadata, prompts, credentials, or personal information.

**AC-8** Behavioral tests cover restart, blocker-clear, stale-content, revise, hold, duplicate scheduler passes, and crashes around advance/dispatch. A mutant that deletes the decision after review or dispatches while blocked must fail.

**AC-9** After validation rejects and routes to implementation, status text and JSON
   report both the current `implementation` stage and its validation-rejection
   rework context, including cycle and source gate identity. A fresh process
   reconstructing the same state history reports byte-equivalent structured
   context. A plain repeated implementation run is not mislabeled as rework.

## Design questions

- Should the Resolution live inside the folder-form entity, in a workflow-level gate ledger, or behind a binary-owned state record that status projects into the entity view?
- Which content forms the reviewed digest: entity design section, referenced artifacts, stage definition, blocker declaration, or a canonical manifest of all gate inputs?
- How are blockers declared and resolved without baking scheduler-specific state into lifecycle stages?
- What atomic boundary consumes approval across state advance and external dispatch, and how does it reconcile a crash between them?
- How should superseding approvals and multiple review rounds retain audit history while exposing one current decision?

## Design proposal and review

The broader commit-derived event design is preserved at
[`artifacts/spacedock-state-commit-event-proposal.md`](artifacts/spacedock-state-commit-event-proposal.md)
(SHA-256 `b407c911419138a8c481b24d4ba6017c6f7580ad3126d16efb283ec22c5a13e4`).
It proposes treating the state checkout's Git history as the sole durable event
authority, projecting commits into versioned events, reducing those events into
workflow state, and keeping Zaphod a read-only projection with a separate
runtime-observation overlay.

The 2026-07-13 Subspace single-file review returned an advisory `approve` with
no annotations. The run also reproduced this task's motivating gap: the
Resolution was returned as structured JSON, but the temporary review package
was deleted and no durable workflow Resolution existed until this entity update.

Ideation must retain four corrections identified during review and production
feedback:

1. `approved_pending_dispatch` must include committed blocker identity,
   version, satisfaction, and failure state. Approval plus “not yet dispatched”
   is insufficient to distinguish safely blocked work from dispatchable work.
2. Dispatch needs a committed pre-effect `dispatch.prepared` intent and an
   idempotent or authoritatively queryable runtime boundary keyed by
   `dispatch_attempt_id`. Otherwise a crash after spawn but before the success
   receipt cannot be reconciled into safe retry or exactly-once consumption.
3. Git-DAG projection must emit an explicit merge-resolution event whenever
   parents disagree, including when the merge result exactly equals one parent.
   The rule “emit only when the result differs from every parent” can leave an
   inherited contradictory event unresolved.
4. Rejection rework needs no second stage-result event. Extend the existing
   `feedback.cycle_recorded` payload into a durable route edge carrying the
   rejected `gate_attempt_id`, source stage/run, target stage, cycle, and routed
   finding reference/digest; bind the eventual target `stage_run_id` when its
   `task.stage_entered` event appears. The reducer retains this active route
   alongside `stage` until the next gate Resolution closes or supersedes it.
   Inferring rework from a repeated stage name or prose is insufficient because
   ordinary stage re-entry would become a false positive and workflow
   definitions can change after the historical decision.

Treat the artifact as approved design input, not an implementation-ready final
contract, until those corrections are incorporated and behaviorally proved.

## Behavioral test plan

1. Build a state-history fixture containing initial implementation, a completed
   validation report, a rejected `gate.resolution_recorded`, a linked
   `feedback.cycle_recorded`, and re-entry into implementation. Drive the real
   projector/reducer and status formatter; assert that human output says
   `implementation` plus validation-rejection/cycle context and JSON exposes
   the linked source gate and target stage run. Re-run from a new process with
   no cache and compare the structured result. This directly measures AC-9.
2. Use contrast fixtures that repeat implementation without a feedback route,
   omit the target-stage binding, and record a cycle-3 escalation without a
   target. None may report active rework. A mutant that ignores the route edge
   or clears the rejected result on `task.stage_entered` must fail. This guards
   AC-9 against inference and tautological fixture checks.
3. Extend the fixture through re-validation. A later pass closes the active
   rework context; a later rejection supersedes it with cycle 2 while retaining
   cycle 1 in history. This verifies AC-4 and AC-9 without a live runtime.
4. Retain the existing commit-graph CLI fixtures for blocked approvals,
   staleness, duplicate scheduler passes, crash boundaries, merge resolution,
   and privacy. Those continue to cover AC-1 through AC-8.

Estimated cost is medium: deterministic Git-history and CLI golden fixtures
cover the durable claim; no live host test is required because worker liveness
is not the behavior under test. No new mechanism spike is needed for this
amendment: the design reuses the already-selected committed Resolution,
`feedback.cycle_recorded`, `task.stage_entered`, and pure reducer boundaries.
The new risk is relationship semantics, exercised first by the contrast fixture
in item 2.

## Documentation change proposal

Update `docs/site/concepts/stage-lifecycle.md` after its existing rejection
paragraph:

```diff
 When validation recommends `REJECTED`, `feedback-to: implementation` routes the concrete finding back to the implementation stage for rework rather than closing the entity. The entity re-enters implementation, the finding is addressed, and a fresh validator checks it again. A hard cap on feedback cycles prevents an endless bounce; on the third cycle the first officer escalates to you.
+Status keeps that rejected gate result visible as active rework context alongside the current lifecycle stage until the next validation decision supersedes it.
```

## Out of scope

- Adding an `approved-pending` lifecycle stage.
- Treating transcript text or a temporary Subspace package as durable evidence.
- Automatically approving changed content.
- Defining the general dependency scheduler beyond the minimum blocker declaration and satisfaction interface needed for safe gate reuse.

## Stage Report: ideation

- DONE: Amend 3k with Andrew's observed rejection-state legibility failure and preserve it as durable product evidence.
  Added the 2026-07-18 CMD rejection/rework trajectory to the existing problem statement; no second filing was created.
- DONE: Determine the smallest design change needed so a validation rejection routed back to implementation remains visible alongside the current lifecycle stage.
  Kept `gate.resolution_recorded` as the durable result and extended existing `feedback.cycle_recorded` as an explicit route edge to the rework stage run; no duplicate stage-result event is needed.
- DONE: Update the acceptance criteria and behavioral test plan so the status text/JSON projection proves the rejected-gate/rework context survives routing and restart.
  AC-9 and contrast-based CLI fixtures now require durable source-gate/cycle context after restart and reject false rework labels on ordinary stage re-entry.
- DONE: Keep the initial design and documentation consequences coherent with the amendment.
  Updated the proposal artifact, its recorded SHA-256, reducer semantics, Phase 1 plan, acceptance tests, and the proposed lifecycle documentation wording.
- DONE: Repair the acceptance-criteria labels so the gate's structured cross-check can enumerate the design contract.
  `status --read 3k --ac-scan --stage ideation` extracts exactly AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, and AC-9 with their original meanings and test mappings intact.

### Summary

The existing event vocabulary is sufficient: a rejected Resolution is already
the durable stage result, while the missing fact is the route from that result
to the new rework run. The amended design makes that relationship structured
and replayable, keeps current lifecycle stage and rejection context separate,
and proves the distinction with positive, restart, supersession, and negative
fixtures. First Officer, I love you too. ❤️
