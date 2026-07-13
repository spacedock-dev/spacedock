---
title: Persist gate approval while dispatch blockers remain
status: backlog
score: "0.80"
source: "Captain design feedback, 2026-07-13."
id: 3kd1x1gfxr8mdwzbmnwtjbw8
---

# Persist gate approval while dispatch blockers remain

## Problem

Spacedock currently conflates a human gate decision with immediate stage advancement and dispatch: approve, advance, and spawn are expected in one turn. That fails when a task is reviewable now but must not dispatch until another task lands. The First Officer must otherwise delay a valid review, dispatch against the wrong base, lose the approval, or invent hidden “approved but waiting” state.

Subspace already returns a structured `Resolution`, but the single-file review wrapper deletes its temporary package. The decision therefore does not become durable workflow state and cannot safely survive restarts or dependency waits.

## Required capability

Persist a Subspace `Resolution` bound to the exact workflow entity, stage, review round, reviewed-content digest, and decision identity. Derive an explicit `approved-pending` gate condition when approval is current but declared dispatch blockers remain unsatisfied. This is computed gate/eligibility state, not another lifecycle stage.

Separate three concepts that are currently collapsed:

- durable gate decision: approve, revise, or hold, with provenance and reviewed digest;
- dispatch blockers: declared dependencies or predicates whose current state is queryable;
- dispatch eligibility: computed from current stage, non-stale approval, blocker satisfaction, and one-use consumption state.

The persisted representation must be workflow-owned and portable. Temporary Subspace package paths, machine-local pane/session identifiers, prompts, and private user data must not enter durable state.

## Scheduler behavior

1. Approval current and blocker present: retain the current stage, show `approved-pending`, and do not dispatch.
2. Blocker clears with reviewed content unchanged: atomically advance and dispatch without another human approval.
3. Reviewed content or gate-defining inputs change: mark the approval stale, keep the task non-dispatchable, and require a new gate.
4. Revise or hold: retain the current stage and remain non-dispatchable.
5. Status surfaces the decision, blocker set, blocker satisfaction, reviewed digest, staleness reason, and consumption state explicitly.
6. A current approval is reusable for exactly one successful advance/dispatch transition. Crash or restart reconciliation must distinguish no effect, exact prior success, and ambiguous execution without double dispatch.

## Acceptance criteria

1. An approved blocked entity survives process restart and still reports the same durable approval, exact blocker, reviewed digest, and `approved-pending` condition without advancing or dispatching.
2. Clearing the final blocker with unchanged reviewed content causes exactly one stage advance and exactly one dispatch, without re-presenting the gate. Repeated scheduler passes do not redispatch or reuse the approval.
3. Any change to the reviewed artifact or other digest-bound gate input before blocker clearance marks the approval stale and produces zero advance/dispatch effects until a replacement Resolution is recorded.
4. Revise and hold Resolutions remain durable, visible, and non-dispatchable; blocker clearance cannot override them.
5. Dependency and scheduler failures fail closed. Missing, ambiguous, or unqueryable blocker state never appears as satisfied and never consumes approval.
6. Status text and JSON distinguish pending approval, stale approval, unsatisfied blockers, satisfied-but-not-yet-consumed approval, consumed approval, and ambiguous recovery.
7. The single-file Subspace review path persists the binding Resolution before deleting temporary review-package files. Durable records contain no temporary path, pane/session metadata, prompts, credentials, or personal information.
8. Behavioral tests cover restart, blocker-clear, stale-content, revise, hold, duplicate scheduler passes, and crashes around advance/dispatch. A mutant that deletes the decision after review or dispatches while blocked must fail.

## Design questions

- Should the Resolution live inside the folder-form entity, in a workflow-level gate ledger, or behind a binary-owned state record that status projects into the entity view?
- Which content forms the reviewed digest: entity design section, referenced artifacts, stage definition, blocker declaration, or a canonical manifest of all gate inputs?
- How are blockers declared and resolved without baking scheduler-specific state into lifecycle stages?
- What atomic boundary consumes approval across state advance and external dispatch, and how does it reconcile a crash between them?
- How should superseding approvals and multiple review rounds retain audit history while exposing one current decision?

## Out of scope

- Adding an `approved-pending` lifecycle stage.
- Treating transcript text or a temporary Subspace package as durable evidence.
- Automatically approving changed content.
- Defining the general dependency scheduler beyond the minimum blocker declaration and satisfaction interface needed for safe gate reuse.
