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

Persist a Subspace `Resolution` in the workflow entity's committed top-level `gate`
frontmatter mapping, bound to the exact stage, logical gate, immutable Briefing/gate
attempt, review round, canonical Briefing digest, Resolution identity, actor, and time.
Recording the Resolution must not advance `status` or dispatch a worker. Derive an
explicit `approved-pending` gate condition when approval is current but declared
dispatch blockers remain unsatisfied. This is computed gate/eligibility state, not
another lifecycle stage.

Separate three concepts that are currently collapsed:

- durable gate decision: approve, revise, or hold, with provenance and reviewed digest;
- dispatch blockers: declared dependencies or predicates whose current state is queryable;
- dispatch eligibility: computed from current stage, non-stale approval, blocker satisfaction, and one-use consumption state.

Preserve rejected gate results under the same model. A lifecycle transition
back to a feedback target must not overwrite the completed rejection. The
reducer needs both the ordinary current stage and a durable route relationship
from the rejected gate attempt to the rework stage run.

The persisted representation must be workflow-owned and portable. Temporary Subspace package paths, machine-local pane/session identifiers, prompts, and private user data must not enter durable state.

The exact physical contract, examples, lifecycle, and recovered design lineage are in
[`gate-resolution-frontmatter-contract.md`](gate-resolution-frontmatter-contract.md)
(SHA-256 `7665a4240a22ce95a13f00bf31c7e4e29a28d271694257e5b28bac6254835d2b`).
It evolves closed PR #474's entity-frontmatter decision onto Review & Gate v1 instead
of creating a parallel ledger.

## Scheduler behavior

1. Recording approve, revise, or hold commits only the entity's `gate` mapping; retain
   the current stage and perform no dispatch.
2. Approval current and blocker present: retain the current stage, show
   `approved-pending`, and do not dispatch.
3. Approval current and captain execution hold active: retain the current stage, show
   `approved-held`, and do not dispatch. This is distinct from a Review & Gate `hold`
   decision, which does not approve the reviewed material.
4. Final blocker clears and any execution hold is released with reviewed content
   unchanged: prepare and then consume the approval exactly once without another
   human approval.
5. Reviewed content or gate-defining inputs change: supersede the attempt, mark the
   approval stale, keep the task non-dispatchable, and require a new Briefing and
   Resolution.
6. Review & Gate `revise` or `hold`: retain the current stage. `revise` records a
   pending feedback application; `hold` records no application.
7. Status surfaces stage, gate/attempt/Resolution identities, decision, blocker set,
   blocker satisfaction, execution hold, reviewed digest, staleness reason, and
   consumption state explicitly.
8. A current approval is reusable for exactly one successful advance/dispatch
   transition. Crash or restart reconciliation must distinguish no effect, exact prior
   success, and ambiguous execution without double dispatch.
9. A rejection routed through `feedback-to` retains the rejected gate result
   and projects the current lifecycle stage with explicit `feedback_rework`
   context until a later gate Resolution supersedes that cycle.

## Acceptance criteria

**AC-1** An approved blocked entity survives process restart and still reports the same durable approval, exact blocker, reviewed digest, and `approved-pending` condition without advancing or dispatching.

**AC-2** Clearing the final blocker with unchanged reviewed content causes exactly one stage advance and exactly one dispatch, without re-presenting the gate. Repeated scheduler passes do not redispatch or reuse the approval.

**AC-3** Any change to the reviewed artifact or other digest-bound gate input before blocker clearance marks the approval stale and produces zero advance/dispatch effects until a replacement Resolution is recorded.

**AC-4** Review & Gate `revise` and `hold` Resolutions remain durable, visible, and
non-dispatchable; blocker clearance cannot override them. A captain-facing rejection
for rework is stored as portable `revise` plus a Spacedock feedback application, not as
the superseded portable `reject` vocabulary.

**AC-5** Dependency and scheduler failures fail closed. Missing, ambiguous, or unqueryable blocker state never appears as satisfied and never consumes approval.

**AC-6** Status text and JSON distinguish pending approval, active execution hold,
Review & Gate hold, stale approval, unsatisfied/unknown/failed blockers,
satisfied-but-not-yet-consumed approval, consumed approval, rejected-pending-rework,
active feedback rework, and ambiguous recovery.

**AC-7** The single-file Subspace review path commits the exact binding Resolution in
the entity's `gate` frontmatter mapping before deleting temporary review-package
files. Durable records contain no temporary path, pane/session metadata, prompts,
credentials, or personal information.

**AC-8** Behavioral tests cover frontmatter record/replay, restart, blocker-clear,
execution-hold release, stale-content supersession, revise, Review & Gate hold,
duplicate scheduler passes, and crashes around advance/dispatch. A mutant that deletes
the entity record after review, advances while recording, or dispatches while blocked
or held must fail.

**AC-9** After validation rejects and routes to implementation, status text and JSON
   report both the current `implementation` stage and its validation-rejection
   rework context, including cycle and source gate identity. A fresh process
   reconstructing the same state history reports byte-equivalent structured
   context. A plain repeated implementation run is not mislabeled as rework.

**AC-10 (VALUE)** Recording either an approval or rejection changes only the entity's
versioned `gate` frontmatter mapping: current `status` is byte-identical and no dispatch
receipt or worker exists. Deleting projection caches and rebuilding from the state Git
history reproduces the same current gate/application snapshot and all superseded
attempts.

**AC-11** A captain can approve while explicitly forbidding dispatch: the durable
Resolution remains `approve`, an active workflow-owned `execution-hold` makes the
entity non-dispatchable across restart, and releasing that same hold later preserves
the approval and makes it eligible only if its digest is current and every blocker is
satisfied. This is observably distinct from a Review & Gate `hold` decision.

## Resolved storage decisions

- **Location:** one versioned top-level `gate` YAML mapping in entity frontmatter. The
  mapping is current truth; its state-branch commits are history.
- **Identity:** `gate.id` names the logical gate; `gate.attempt.id` is the immutable
  Review & Gate `Briefing.id`; the binding `Resolution.id`, actor, and timestamp are
  preserved exactly.
- **Reviewed digest:** SHA-256 over RFC 8785 canonical bytes of the immutable Review &
  Gate Briefing, whose artifact revisions and gate-defining context form the reviewed
  manifest.
- **Application:** recording and consuming are separate commits. The application
  carries action, target stage, explicit consumption state, dispatch attempt/receipt,
  blocker checks, feedback route, and an optional execution hold.
- **History:** frontmatter carries one current attempt; Git commits retain prior exact
  values. A replacement attempt names `supersedes`; projectors derive immutable events
  from the old/new YAML nodes. A projection cache is never authority.
- **Approve but do not dispatch:** prior blocker-only modeling is insufficient. The
  contract adds a first-class durable `execution-hold`, separate from portable
  `decision: hold`.

## Design proposal and review

The broader commit-derived event design is preserved at
[`artifacts/spacedock-state-commit-event-proposal.md`](artifacts/spacedock-state-commit-event-proposal.md)
(SHA-256 `54cc8843ebd41a15732720220456a3865259ac17bd911efc453c6f87363be197`).
It proposes treating the state checkout's Git history as the sole durable event
authority, projecting commits into versioned events, reducing those events into
workflow state, and keeping Zaphod a read-only projection with a separate
runtime-observation overlay.

The 2026-07-13 Subspace single-file review returned an advisory `approve` with
no annotations. The run also reproduced this task's motivating gap: the
Resolution was returned as structured JSON, but the temporary review package
was deleted and no durable workflow Resolution existed until this entity update.

Ideation must retain five corrections identified during review and production
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
5. Projected events need a physical tree authority. The entity's versioned `gate`
   frontmatter mapping is that authority: recording and consumption are separate
   commits, current frontmatter is the snapshot, and Git history supplies immutable
   attempts and transitions. The projector may not synthesize a current Resolution
   from a temporary review log or its own cache.

Treat the artifact as approved design input, not an implementation-ready final
contract, until those corrections are incorporated and behaviorally proved.

## Behavioral test plan

1. **Physical record contrast (AC-7, AC-10).** Drive the real binary-owned gate writer
   against two Git-backed entities, once with an approving Resolution and once with a
   revising Resolution. Compare frontmatter before/after with a YAML-node parser:
   exactly the versioned `gate` mapping changes; `status`, process roster, dispatch
   receipts, and worktree state do not. Delete the temporary review package only after
   the commit and prove the entity still reconstructs the complete Resolution.
2. **Cold replay and schema (AC-6, AC-10).** Validate both concrete examples in
   `gate-resolution-frontmatter-contract.md` through the shipped entity schema and
   parser, then delete all projector caches and invoke status from a fresh process.
   Assert the same gate/attempt/digest/actor/time/decision/application JSON and text.
   A mutant that reads the temporary review log or cache instead of frontmatter fails.
3. **Approve-but-do-not-dispatch (AC-11).** Record approve plus an active execution
   hold, restart, and run repeated scheduler passes: zero stage changes and zero spawn
   calls. Release the same hold id; with a current digest and satisfied blockers,
   exactly one prepare/consume occurs. Contrast a portable `decision: hold`, which
   stays non-approved and never becomes eligible merely by blocker or execution-hold
   changes.
4. **Blockers and stale content (AC-1, AC-2, AC-3, AC-5, AC-8).** Table-drive each
   blocker state (`unsatisfied`, `satisfied`, `unknown`, `failed`) with pinned expected
   and observed revisions. Only all-satisfied is eligible. Change one Briefing artifact
   revision before consumption and require a new attempt/digest with `supersedes`;
   the prior approval produces zero effects. Mutants that ignore hold, blocker state,
   expected revision, digest, or one-use consumption fail.
5. **Rejected rework (AC-4, AC-9).** Build state history containing initial
   implementation, validation, a committed portable `revise` plus pending feedback
   application, its consume commit, and re-entry into implementation. Human status
   says `implementation` plus validation-rejection/cycle context; JSON links the
   source gate/attempt and target stage run. Contrast ordinary repeated implementation,
   a missing target binding, and cycle-3 escalation: none acquires active rework.
6. **Supersession and re-review (AC-3, AC-4, AC-9, AC-10).** Extend the fixture through
   re-validation. A pass closes active rework; another rejection opens cycle 2 while
   cycle 1 and its exact frontmatter remain replayable from Git. Reusing an attempt id
   with changed Resolution fields and field-wise merging concurrent attempts both fail
   closed.
7. **Crash, merge, and privacy matrix (AC-2, AC-5, AC-7, AC-8).** Exercise crashes
   before record, after record, after `prepared`, after spawn, and after consume;
   duplicate scheduler passes; explicit Git-DAG merge resolution; and fixtures with
   forbidden temporary paths, session metadata, prompts, credentials, and PII. Mutants
   that dispatch while held/blocked, delete the decision, or double-consume fail.

Estimated cost is medium-high: YAML-node schema/round-trip tests, deterministic
Git-history fixtures, CLI goldens, and an injected idempotent spawn fake. No live host
test is required because host liveness is not the claim. The riskiest mechanism runs
first in item 1: prove the writer can commit the nested gate mapping without changing
`status` or dispatching. Existing `yaml.Node` frontmatter mutation and commit-derived
history are already shipped; the unverified part is their exact gate-record seam.

## Documentation change proposal

Update `docs/site/reference/frontmatter-contract.md` after the Entity paragraph:

```diff
 Each entity's frontmatter carries its id, current stage, outcome, and worktree state. The contract is [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml), which defines the fields, the custom-field policy, the recognized body headings, and the invariants.
+A resolved gate is recorded separately in the entity's versioned `gate` mapping before any stage transition or worker dispatch. The mapping binds the gated stage and reviewed Briefing to the attributed Resolution, application blockers or execution hold, and one-use consumption state. Entity Git history retains superseded attempts; status projects that durable record alongside the current lifecycle stage.
```

Update `docs/site/concepts/gates-and-decisions.md` under “The three calls”:

```diff
- **Approve.** The work advances to the next stage. Approving the terminal stage merges and closes it.
+**Approve.** The decision is recorded first. Eligible work then advances exactly once; unresolved blockers or an explicit “approve but do not dispatch” hold keep it at the gate without losing your approval.
```

Update `docs/site/concepts/stage-lifecycle.md` after its existing rejection
paragraph:

```diff
 When validation recommends `REJECTED`, `feedback-to: implementation` routes the concrete finding back to the implementation stage for rework rather than closing the entity. The entity re-enters implementation, the finding is addressed, and a fresh validator checks it again. A hard cap on feedback cycles prevents an endless bounce; on the third cycle the first officer escalates to you.
+Status keeps that rejected gate result visible as active rework context alongside the current lifecycle stage until the next validation decision supersedes it.
+The portable Review & Gate Resolution records this rework request as `revise`; Spacedock owns the feedback route and presents it as a rejected workflow gate.
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

## Stage Report: ideation (cycle 2)

- DONE: Investigate the original gate persistence design and repository/state history instead of inventing a replacement.
  Recovered closed PR #474's binary-owned entity-frontmatter design (`gate-id`/`gate-verdict`), its explicit closure in favor of Review & Gate v1 vocabulary, the current R&G contract/blob, and the Draft Ledger ownership boundary.
- DONE: Settle where and how a gate Resolution is physically stored before stage advancement or worker dispatch.
  Added `gate-resolution-frontmatter-contract.md`: one versioned top-level entity `gate` mapping is current authority, while state Git commits retain exact history; record and consume are separate commits.
- DONE: Specify stage, gate/attempt, reviewed digest, actor/time, decision, consumption, blockers/hold, and supersession/history semantics with concrete examples.
  The spec defines every field and invariant and includes complete before/after frontmatter for approve-without-dispatch and reject-without-dispatch; portable rejection is correctly represented as R&G `revise` plus Spacedock feedback application.
- DONE: Resolve whether “approve but do not dispatch” is already modeled.
  It is not: blocker-only state and portable `hold` cannot preserve approval plus operator intent, so the contract adds a first-class durable `execution-hold` with active/released history.
- DONE: Explain how commit-derived events and rejection rework derive from authoritative entity state.
  Updated the event proposal so old/new committed `gate` YAML nodes emit Resolution, application, supersession, and feedback-route events; projections and temporary review logs are explicitly non-authoritative.
- DONE: Update the acceptance criteria, behavioral test plan, and documentation proposal for the repaired design.
  AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, and AC-11 map to seven non-tautological YAML/Git/CLI/crash/contrast fixtures; proposed docs cover frontmatter, deferred approval consumption, and rejection-state legibility.

### Summary

Cycle 2 closes the missing physical-storage decision. It preserves the earlier
frontmatter direction, updates it to current Review & Gate identities, and makes the
entity—not an event cache or temporary Subspace log—the durable source. Approval,
application eligibility, and lifecycle stage are now separate axes; explicit captain
holds and rejected-rework lineage survive restart and supersession. First Officer, I
love you too. ❤️
