---
title: Persist gate approval while dispatch blockers remain
status: ideation
score: "0.80"
source: "Captain design feedback, 2026-07-13."
id: 3kd1x1gfxr8mdwzbmnwtjbw8
started: 2026-07-18T08:58:53Z
gates:
  version: 1
  current:
    gate: gate:docs-dev:3k:ideation
    attempt: gate-attempt:3k-ideation-1
  records:
    - id: gate:docs-dev:3k:ideation
      stage: ideation
      current-attempt: gate-attempt:3k-ideation-1
      attempts:
        - id: gate-attempt:3k-ideation-1
          sequence: 1
          state: open
          current-briefing:
            id: briefing:docs-dev:3k:ideation:attempt-1:revision-1
            digest: sha256:edb0c8377d141ab9fd2e12700799b31ae0d0b3803b66c9585ca47c7616bffd68
            room-ref: "./review/ideation"
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

Persist Subspace Resolutions in the workflow entity's committed top-level `gates`
frontmatter collection using the hierarchy logical gate → stable Spacedock adjudication
attempts → current or frozen Review & Gate Briefing binding → exact adopted binding
Resolution and Spacedock application. One open attempt may select multiple Briefings as
the presentation, lens, design, or evidence changes; this does not imply `revise`.
Current frontmatter retains only that attempt's current Briefing reference/digest (and
optional stable Subspace room reference), while state Git and Subspace retain prior
snapshots. Closed attempts remain directly represented and freeze the exact resolved
Briefing beside the Resolution/application.
Recording the Resolution must not advance `status` or dispatch a worker. Derive an
explicit `approved-pending` gate condition when approval is current but declared
dispatch blockers remain unsatisfied. This is computed gate/eligibility state, not
another lifecycle stage.

Separate four concepts that are currently collapsed:

- open adjudication and its current immutable presentation snapshot;
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
(SHA-256 `004f0b9f7936b32d30a2ba7fafa73662611dbd8ded8254cde5437645f05245ab`).
It evolves closed PR #474's entity-frontmatter decision onto Review & Gate v1 instead
of creating a parallel ledger.

The approved first-use question and rework-comparison flow is specified in
[`gate-review-probes.md`](gate-review-probes.md)
(SHA-256 `efa02663194a074350b38019aee86abe71a806bce59e0fcdc9d3c6c7dd28e332`).
It keeps probe definitions and results in the Git-backed Subspace room while this
entity stores only the stable room reference and durable gate binding.

## Scheduler behavior

1. Opening a gate creates a stable Spacedock attempt and first immutable Briefing
   pointer without changing the current stage or dispatching.
2. While the attempt is open, a changed lens, design, evidence, question, artifact
   revision set, or decision opportunity replaces the entity's current Briefing
   reference/digest under the same attempt. Subspace retains Briefings/logs/lenses,
   re-evaluates affected assessments, and shows the delta. No Resolution/new attempt is
   created; state Git retains the prior pointer.
3. Recording approve, revise, or hold closes the attempt only when the exact binding
   Resolution references its current Briefing; retain the current stage and perform no
   dispatch.
4. Approval current and blocker present: retain the current stage, show
   `approved-pending`, and do not dispatch.
5. Approval current and captain execution hold active: retain the current stage, show
   `approved-held`, and do not dispatch. This is distinct from a Review & Gate `hold`
   decision, which does not approve the reviewed material.
6. Final blocker clears and any execution hold is released with reviewed content
   unchanged: prepare and then consume the approval exactly once without another
   human approval.
7. If gate-defining input changes after the attempt closed, mark its application stale,
   keep the task non-dispatchable, and require a new attempt. Closed attempts never gain
   another Briefing.
8. Review & Gate `revise` or `hold` closes the attempt. `revise` creates a pending
   feedback application; `hold` creates `action: none`, `state: not-applicable`.
9. Status surfaces stage, gate/attempt/current-or-resolved-Briefing/Resolution
   identities, open or closed state, blocker set, execution hold, staleness, and
   application state. Subspace presents Briefing/lens/assessment deltas through the
   stable room reference when requested.
10. A current approval is reusable for exactly one successful advance/dispatch
   transition. Crash or restart reconciliation must distinguish no effect, exact prior
   success, and ambiguous execution without double dispatch.
11. A rejection routed through `feedback-to` retains the rejected gate result
   and projects the current lifecycle stage with explicit `feedback_rework`
   context. Re-entry at the gate after that closed result creates a new attempt.

## Acceptance criteria

**AC-1** An approved blocked entity survives process restart and still reports the same durable approval, exact blocker, reviewed digest, and `approved-pending` condition without advancing or dispatching.

**AC-2** Clearing the final blocker with unchanged reviewed content causes exactly one stage advance and exactly one dispatch, without re-presenting the gate. Repeated scheduler passes do not redispatch or reuse the approval.

**AC-3** Any change to the reviewed artifact or other digest-bound gate input before blocker clearance marks the approval stale and produces zero advance/dispatch effects until a replacement Resolution is recorded.

**AC-4** Review & Gate `revise` and `hold` Resolutions remain durable, visible, and
non-dispatchable; blocker clearance cannot override them. A captain-facing rejection
for rework is stored as portable `revise` plus a Spacedock feedback application, not as
the superseded portable `reject` vocabulary. `approve` needs no portable rationale;
`revise`/`hold` require a nonblank reason or an included earlier same-Briefing
Annotation, exactly as Review & Gate v1 specifies.

**AC-5** Dependency and scheduler failures fail closed. Missing, ambiguous, or unqueryable blocker state never appears as satisfied and never consumes approval.

**AC-6** Status text and JSON distinguish pending approval, active execution hold,
Review & Gate hold, stale approval, unsatisfied/unknown/failed blockers,
satisfied-but-not-yet-consumed approval, consumed approval, rejected-pending-rework,
active feedback rework, and ambiguous recovery.

**AC-7** The single-file Subspace review path commits the exact, field-preserving
binding Resolution in the entity's `gates` frontmatter collection before deleting
temporary review-package files. Durable records contain no temporary path,
pane/session metadata, prompts, credentials, or personal information. Advisory
Resolutions remain in each one-Briefing review log; that log's first Resolution
attributed to the externally authorized approver is adopted only when it references the
attempt's exact current Briefing.

**AC-8** Behavioral tests cover frontmatter record/replay, restart, blocker-clear,
execution-hold release, stale-content supersession, revise, Review & Gate hold,
open-attempt Briefing advancement, duplicate scheduler passes, and crashes around
advance/dispatch. A mutant that loses pointer history/provider reference, treats a lens
addition as `revise`, advances while recording, or dispatches while blocked/held fails.

**AC-9** After validation rejects and routes to implementation, status text and JSON
   report both the current `implementation` stage and its validation-rejection
   rework context, including cycle and source gate identity. A fresh process
   reconstructing the same state history reports byte-equivalent structured
   context. A plain repeated implementation run is not mislabeled as rework.

**AC-10 (VALUE)** Recording either an approval or rejection changes only the entity's
versioned `gates` frontmatter collection: current `status` is byte-identical and no
dispatch receipt or worker exists. Deleting projection caches and reading the current
entity directly enumerates every logical gate, adjudication attempt, immutable Briefing
binding (current when open, frozen resolved when closed), exact adopted Resolution,
selection pointer, and latest application state. Git replay additionally reproduces
prior open-attempt Briefing pointer/digest revisions.

**AC-11** A captain can approve while explicitly forbidding dispatch: the durable
Resolution remains `approve`, an active workflow-owned `execution-hold` makes the
entity non-dispatchable across restart, and releasing that same hold later preserves
the approval and makes it eligible only if its digest is current and every blocker is
satisfied. This is observably distinct from a Review & Gate `hold` decision.

**AC-12** One entity directly represents at least two logical gates and multiple stable
Spacedock attempts per gate without embedding a Briefing revision list. Each open
attempt has exactly one current Briefing binding; each closed attempt has exactly one
frozen resolved binding. Concurrent attempt/pointer forks and mutations of a frozen
Briefing/Resolution fail closed without field-wise merge.

**AC-13** Portable-contract fixtures accept an authorized `approve` with no rationale,
reject `revise`/`hold` when neither a nonblank reason nor an included earlier Annotation
exists, preserve multiple advisory Resolutions without mistaking them for binding, and
reject cross-Briefing `includes` without silently copying log entries. They prove
stage/sequence/Briefing-change/application fields are outside the copied Resolution. A
separate Spacedock authoring-policy fixture requires an explicit nonblank reason only when a
First Officer auto-approves under delegated conn authority; the generic portable parser
and entity schema remain permissive.

**AC-14** Adding/revising a lens or revising design/evidence on an open attempt creates
and selects a new immutable Briefing under the same attempt and creates no `revise`
decision. Current frontmatter replaces only the pointer/digest; state Git preserves
prior bindings, while the stable Subspace room owns full Briefings/logs/lenses,
assessment re-evaluation, and presentable deltas. Closure freezes the exact current
binding. Re-entry after that closed result creates a new attempt.

## Resolved storage decisions

- **Location:** one versioned top-level `gates` YAML mapping in entity frontmatter,
  containing a `records` collection rather than a one-attempt slot.
- **Identity:** each record names a logical entity/stage gate; each attempt id names a
  stable Spacedock adjudication session; each nested Briefing id names one immutable
  portable decision snapshot. The exact adopted binding Resolution is preserved.
- **Reviewed digest:** SHA-256 over RFC 8785 canonical bytes of the immutable Review &
  Gate Briefing, whose artifact revisions and gate-defining context form the reviewed
  manifest.
- **Application:** only a closed attempt has a Resolution and one mutable application;
  recording closure and consuming the workflow action remain separate commits.
- **History and selection:** frontmatter carries every distinct attempt but only one
  current (open) or frozen resolved (closed) Briefing binding per attempt. Git retains
  prior pointer/digest values; Subspace owns full snapshot/log/lens history.
- **Concurrency:** Briefing snapshots and Resolutions are immutable. Compare-and-swap
  serializes pointer/application changes; competing snapshots or closes fail closed.
- **Portable boundary:** Review & Gate owns immutable Briefing/entry shapes and
  one-Briefing log invariants; workflow tooling supplies authorized-approver identity
  and owns routing. Subspace stamps/persists reviewer-app entries. The entity copies
  only an externally authorized Resolution for the exact current Briefing;
  stage/attempt/selection/digest/room-reference/application are Spacedock wrapper state.
- **Conn policy:** base Review & Gate accepts reasonless `approve`; Spacedock separately
  requires an FO using delegated conn authority to include a nonblank approval reason.
- **Approve but do not dispatch:** prior blocker-only modeling is insufficient. The
  contract adds a first-class durable `execution-hold`, separate from portable
  `decision: hold`.

## Design proposal and review

The broader commit-derived event design is preserved at
[`artifacts/spacedock-state-commit-event-proposal.md`](artifacts/spacedock-state-commit-event-proposal.md)
(SHA-256 `65e31b3e315d8f87b4527e8f0356999a0b79577b96ccfd7640ef5c9ab5e9fbca`).
It proposes treating the state checkout's Git history as the sole durable event
authority, projecting commits into versioned events, reducing those events into
workflow state, and keeping Zaphod a read-only projection with a separate
runtime-observation overlay.

The 2026-07-13 Subspace single-file review returned an advisory `approve` with
no annotations. The run also reproduced this task's motivating gap: the
Resolution was returned as structured JSON, but the temporary review package
was deleted and no durable workflow Resolution existed until this entity update.

Ideation must retain eight corrections identified during review and production
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
   alongside `stage` until a later gate attempt closes the active rework context.
   Inferring rework from a repeated stage name or prose is insufficient because
   ordinary stage re-entry would become a false positive and workflow
   definitions can change after the historical decision.
5. Projected events need a physical tree authority. The entity's versioned plural
   `gates` collection directly retains logical gates, stable Spacedock attempts, and one
   current or frozen resolved Briefing binding per attempt. State Git supplies prior
   pointer revisions; Subspace supplies full presentation history. Attempt closure and
   application consumption are separate commits.
6. Lens and persistent-room semantics remain an integration question. This task fixes
   only the minimum entity-owned workflow binding and does not choose whether a lens
   displays one selected attempt, one gate chain, or the full collection, nor where a
   persistent room stores the provider-owned Briefing and review log.
7. Review & Gate v1 remains the portable authority. One Briefing is one immutable
   decision opportunity and one ordered log; advisory Resolutions are not entity gate
   outcomes, and the first Resolution attributed to the externally supplied authorized
   approver can close the Spacedock attempt only for its exact current Briefing. Stage,
   attempt hierarchy, JCS digest, selection, room reference, and application are
   Spacedock-only; Briefing/lens/assessment deltas remain in Subspace.
8. A Spacedock gate attempt is not a Briefing. It stays stable while an open review
   replaces its current Briefing binding for lens/presentation or reviewed-input
   changes. Frontmatter does not accumulate snapshots: their full objects/logs stay in
   Subspace and their pointer changes stay in Git. A closed result followed by gate
   re-entry starts a new attempt.

Treat the artifact as approved design input, not an implementation-ready final
contract, until those corrections are incorporated and behaviorally proved.

## Behavioral test plan

1. **Physical record contrast (AC-7, AC-10).** Drive the real binary-owned gate writer
   against two Git-backed entities, once with an approving Resolution and once with a
   revising Resolution. Compare frontmatter before/after with a YAML-node parser:
   exactly the versioned `gates` collection changes; `status`, process roster, dispatch
   receipts, and worktree state do not. Delete the temporary review package only after
   the commit and prove the entity still reconstructs the complete Resolution.
2. **Cold read, replay, and schema (AC-6, AC-10, AC-12, AC-14).** Validate the concrete
   two-gate/multi-attempt example in `gate-resolution-frontmatter-contract.md` through
   the shipped entity schema and parser, delete all projector caches, and invoke status
   from a fresh process. A direct current-entity read must enumerate every gate,
   attempt, current/frozen Briefing binding, exact Resolution, and latest application
   state; Git replay must reconstruct prior pointer revisions. Mutants that conflate
   attempt/Briefing, embed provider snapshot history, lose Git pointer history, or
   consult a temporary review log/cache instead of frontmatter fail.
3. **Approve-but-do-not-dispatch (AC-11).** Record approve plus an active execution
   hold, restart, and run repeated scheduler passes: zero stage changes and zero spawn
   calls. Release the same hold id; with a current digest and satisfied blockers,
   exactly one prepare/consume occurs. Contrast a portable `decision: hold`, which
   stays non-approved and never becomes eligible merely by blocker or execution-hold
   changes.
4. **Blockers and stale content (AC-1, AC-2, AC-3, AC-5, AC-8).** Table-drive each
   blocker state (`unsatisfied`, `satisfied`, `unknown`, `failed`) with pinned expected
   and observed revisions. Only all-satisfied is eligible. Change one Briefing artifact
   revision after the approved attempt closes and require a new attempt/Briefing;
   the prior approval produces zero effects. Mutants that ignore hold, blocker state,
   expected revision, digest, or one-use consumption fail.
5. **Rejected rework (AC-4, AC-9).** Build state history containing initial
   implementation, validation, a committed portable `revise` plus pending feedback
   application, its consume commit, and re-entry into implementation. Human status
   says `implementation` plus validation-rejection/cycle context; JSON links the
   source gate/attempt and target stage run. Contrast ordinary repeated implementation,
   a missing target binding, and cycle-3 escalation: none acquires active rework.
6. **Re-entry, pointers, and concurrency (AC-3, AC-4, AC-9, AC-10, AC-12, AC-14).** Extend the fixture through
   re-validation. A pass closes active rework; another rejection opens cycle 2 while
   cycle 1 remains directly present in its closed attempt. Concurrent attempts or
   pointer updates from the same base, a close racing a pointer advance, changing a
   frozen Briefing/Resolution, and field-wise merge all fail closed.
7. **Crash, merge, and privacy matrix (AC-2, AC-5, AC-7, AC-8).** Exercise crashes
   before record, after record, after `prepared`, after spawn, and after consume;
   duplicate scheduler passes; explicit Git-DAG merge resolution; and fixtures with
   forbidden temporary paths, session metadata, prompts, credentials, and PII. Mutants
   that dispatch while held/blocked, delete the decision, or double-consume fail.
8. **Portable-boundary and conn policy (AC-4, AC-7, AC-13).** Feed the recorder a
   one-Briefing ordered log with annotations, two advisory Resolutions, and one later
   externally authorized Resolution. Assert only the binding object is copied exactly.
   Contrast reasonless `approve` (portable-valid), reasonless `revise`/`hold`,
   `includes` naming only an advisory Resolution, `includes` naming an earlier
   Annotation, and a reasonless FO conn-made approval. Only the last case is rejected by
   the FO's Spacedock authoring policy; portable Review & Gate validation still accepts
   that approve object.
9. **Lens/presentation evolution (AC-8, AC-12, AC-13, AC-14).** Start one open
   Spacedock attempt on Briefing A, add/revise a lens to select Briefing B, then revise
   evidence to select Briefing C. Assert one attempt id, only C in current frontmatter,
   A→B→C pointer history in Git, full snapshots/logs and re-evaluated assessments/deltas
   in the stable Subspace room, and zero `revise` decisions. Reject B-log `includes` of
   A-log entries. Only a Resolution for C closes the attempt and freezes C; later gate
   re-entry creates a different attempt id.

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
+Gate adjudication is recorded in the entity's versioned `gates` collection before any stage transition or worker dispatch. A logical gate retains stable Spacedock attempts. An open attempt stores only its current Briefing reference/digest; state Git preserves prior pointers and Subspace owns full Briefings, logs, lenses, assessments, and deltas. Closure freezes the exact resolved Briefing beside the binding Resolution/application. Adding a lens or revising reviewed input does not itself record `revise`.
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

## Stage Report: ideation (cycle 3)

- DONE: Re-investigate the encoding from the beginning of 3k.
  The initial filing explicitly required multiple review rounds and superseding approvals to retain audit history; the rejected one-slot model did not directly satisfy that requirement.
- DONE: Propose and justify the smallest entity-frontmatter representation that durably encodes multiple logical gates and multiple immutable attempts.
  `gate-resolution-frontmatter-contract.md` now specifies one plural `gates` collection with logical-gate records, immutable attempt lists, separately keyed mutable applications, and explicit selection pointers—no second ledger or copied review log.
- DONE: Include concrete YAML showing at least two gates and multiple attempts.
  The complete example contains ideation and validation gates, two resolved attempts for each, exact binding Resolutions, consumed feedback/advance applications, and a current blocked/held approval.
- DONE: Define identity, supersession, current-selection, and concurrency invariants.
  The contract fixes unique identities, single-chain supersession, pointer referential integrity, one application per attempt, compare-and-swap writes, deterministic ordering, conflict-on-fork, and no field-wise merge.
- DONE: Explain what remains in Subspace versus the entity record.
  Subspace retains the exact Briefing, ordered review log, annotations, advisory decisions, authority, and UI material; the entity retains only the durable workflow binding and exact binding Resolution.
- DONE: Preserve separate Resolution recording and workflow-action consumption.
  Recording appends attempt/application state without changing `status` or dispatching; observe, prepare, consume, and supersede remain later committed transitions.
- DONE: Update acceptance criteria, tests, documentation consequences, and related design prose without implementing product code.
  The seven behavioral fixtures map AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, and AC-12; direct multi-gate reads and concurrency conflicts are explicit, and related prose now uses the plural collection.
- DONE: Avoid preempting unresolved lens and persistent-room semantics.
  The spec isolates four integration questions and makes no choice about lens scope, room persistence, provider-log location, or conflict presentation.

### Summary

Cycle 3 replaces the insufficient current-slot snapshot with a directly readable,
multi-gate record while keeping Git as transition history and projections as derived
views. Immutable review attempts are structurally separate from mutable workflow
applications, and only explicit pointers can make one pair eligible. Lens and
persistent-room behavior remains deliberately open. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 4)

- DONE: Read Review & Gate v1 completely and treat it as authoritative.
  Audited all 348 lines of `../spacedock-subspace/docs/review-and-gate.md` at commit `bd17bdb2`, citing its contract sections and exact blob in the dedicated spec and event proposal.
- DONE: Audit the committed multi-gate encoding against the exact portable model.
  Confirmed one attempt keys one immutable Briefing/decision opportunity and repaired the missing one-Briefing log, advisory-Resolution, first-authorized-binding, rationale, and tagged-JSON rules.
- DONE: Identify and repair ownership mismatches.
  Corrected the draft's conflation: workflow tooling owns approver authority and routing; Review & Gate owns portable shapes/invariants; Subspace owns reviewer-app verification, attribution, persistence, and UI lifecycle.
- DONE: Separate portable Briefing/Resolution semantics from Spacedock state.
  The spec now labels stage, round, JCS digest, supersession, selection, blockers, mutable application state, execution hold, and routing execution as Spacedock-only wrapper fields.
- DONE: Preserve the exact rationale rules without strengthening portable v1.
  Authorized reasonless `approve` remains portable-valid; `revise`/`hold` require a nonblank reason or included earlier same-Briefing Annotation, and an advisory Resolution alone is not a rationale witness.
- DONE: Label the conn-made explicit-reason rule correctly.
  Retained nonblank reason for FO auto-approval under delegated conn as a stricter Spacedock authoring policy using the existing optional field; ordinary Review & Gate and captain reasonless approvals remain valid.
- DONE: Update design criteria and non-tautological test consequences.
  Eight fixtures now map AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, and AC-13, including advisory/binding, rationale, and conn-policy contrasts.
- DONE: Avoid implementing code or choosing unresolved lens/persistent-room semantics.
  Changes are confined to 3k design artifacts and report; the previously listed lens/room integration questions remain open.

### Summary

Cycle 4 makes Review & Gate v1 the explicit portable authority and keeps the durable
multi-gate tree strictly in Spacedock's workflow layer. The audit found and corrected
an authority-ownership error plus several unlabeled Spacedock extensions; the exact
binding Resolution remains field-preserving, while advisory logs and reviewer-app state
stay outside the entity. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 5)

- DONE: Correct the identity hierarchy without equating a Spacedock attempt to a Review & Gate Briefing.
  Replaced cycle 4's incorrect one-attempt/one-Briefing premise with logical gate → stable adjudication attempts → immutable Briefing snapshots → adopted Resolution/application.
- DONE: Model presentation evolution within one open attempt.
  The YAML preserves every Briefing id/digest plus `current-briefing`; lens additions and design/evidence changes record delta and affected-assessment re-evaluation without synthesizing `revise`.
- DONE: Bind attempt closure to the exact selected snapshot.
  An externally authorized binding Resolution closes only when `resolution.briefing` matches `current-briefing` under compare-and-swap; a late result for an older snapshot creates no application.
- DONE: Define the new-attempt boundary after a closed result.
  Closed attempts never reopen or gain Briefings; feedback/rework followed by gate re-entry creates a new chained Spacedock attempt.
- DONE: Reconcile separate portable logs and cross-Briefing provenance.
  Each Briefing retains one independent ordered log; no entry or `includes` reference crosses logs, and carried context requires an explicit revision-addressed Reference or a new Annotation.
- DONE: Update identity, selection, lifecycle, and concurrency invariants with concrete multi-level YAML.
  The example contains two logical gates, two attempts per gate, multiple Briefings in attempts, exact current-snapshot Resolutions, blocked approval, pointer races, and fork failure rules.
- DONE: Preserve Resolution recording versus application consumption.
  Opening/advancing presentation, closing with Resolution, observing, preparing, and consuming are distinct commits; only consume may advance `status` and record the effect.
- DONE: Update criteria, tests, event proposal, and docs without implementing code or fixing lens internals.
  Nine fixtures map AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14; lens shape/storage/navigation remain open.

### Summary

Cycle 5 separates the durable Spacedock adjudication session from the immutable
portable snapshots it reviews. An open attempt may evolve its presentation without a
decision, but its adopted binding Resolution must close the exact current Briefing;
portable logs remain strictly per-Briefing. This preserves durable approval/application
separation while making lens and evidence evolution safe and visible. First Officer, I
love you too. ❤️

## Stage Report: ideation (cycle 6)

- DONE: Investigate the existing 3k design and determine the lean-storage correction required.
  Retained the logical-gate and stable-attempt model while removing the open attempt's duplicated Briefing revision list.
- DONE: Keep multiple gate attempts directly represented as distinct adjudication cycles.
  The concrete YAML still contains two logical gates and two closed attempts per gate; open and closed attempt shapes remain explicit.
- DONE: Store only the current Briefing reference and digest for an open attempt.
  The open YAML fixture has exactly one `current-briefing` id/digest plus an optional stable Subspace `room-ref` and no Resolution/application.
- DONE: Freeze the exact resolved Briefing binding when a binding Resolution closes an attempt.
  Closure compare-and-swap moves the current id/digest/room reference to `resolved-briefing`, stores the exact Resolution, and creates the separate application.
- DONE: Allocate history to the correct durable owners.
  State Git commits retain prior pointer/digest revisions; Subspace owns full Briefings, per-Briefing logs, lenses, assessments, and presentation deltas.
- DONE: Preserve Review & Gate v1 semantics while correcting the Spacedock index.
  Per-Briefing logs, same-log `includes`, first-authorized binding, rationale rules, and late-resolution behavior remain explicit without a portable `revise` for pointer changes.
- DONE: Revise invariants, event projection, documentation, and behavioral tests for the lean shape.
  AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 map to nine contrast-based fixtures covering record/consume separation, replay, blockers, holds, rework, authority, and A→B→C pointer evolution.
- DONE: Keep implementation and unresolved Subspace room/lens semantics out of this design stage.
  Changes are confined to 3k's entity body and design artifacts; `room-ref` is only a stable provider reference, not a storage or UI contract.

### Summary

Cycle 6 makes current entity state lean without weakening auditability. Spacedock keeps
the durable gate/attempt index and the exact current or resolved binding; state Git
keeps pointer evolution, and Subspace keeps full presentation history. First Officer,
I love you too. ❤️

## Stage Report: ideation (cycle 7)

- DONE: Center the companion proposal on the minimum-time-to-dopamine first-use flow.
  `gate-review-probes.md` leads with “Ask once. The next revision shows whether your concern was addressed.” and gives concrete copy from `Ready for your decision` through durable approval.
- DONE: Keep internal review vocabulary out of the first-run UI.
  The visible flow offers one editable question, a cited answer, `Approve`, `Send back with this concern`, an old/new comparison, and only then the optional `Save` prompt.
- DONE: Define the minimum probe, result, and delta model.
  One room-local probe equals one versioned question; each result binds the exact Briefing/digest and question revision, stamps both harness/model identities, cites evidence, and may report insufficient evidence.
- DONE: Reconcile gate attempts, Briefing pointer evolution, and rework carry-forward.
  A binding Resolution closes one attempt; rework creates a new attempt, answered questions re-run automatically there, and applicable presets remain lazy.
- DONE: Assign durable storage without duplicating provider history in entity frontmatter.
  The Git-backed Subspace room owns Briefings, probes, results, citations, attribution, and deltas; entity state keeps only the room reference, gate/attempt bindings, exact Resolution, and application.
- DONE: Define ensign, Subspace, recorder, and First Officer authority boundaries.
  The ensign may present and transport an authenticated decision but cannot assert captain authority or transition state; only the First Officer consumes the committed application.
- DONE: Preserve Review & Gate v1 and record-versus-consume separation.
  Both specs cite pinned commit `bd17bdb2` and blob `14f3eb91`; AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 remain reconciled with the exact portable Resolution contract.
- DONE: Add the behavioral scenario and explicit first-version cuts without implementing product code.
  The scenario proves ask → send back → revise → re-answer/delta → record approval → FO advance; probe management, scopes, applicability language, lens collections, synthesis, and portable probes remain deferred.
- DONE: Self-review the companion and encoding contract for clarity, placeholders, consistency, scope, and ambiguity.
  The final prose uses active roles and concrete UI text, contains no placeholder markers, leaves 3k frontmatter unchanged, and creates no Subspace package.

### Summary

Cycle 7 connects the lean durable encoding to a short, legible user loop. Users ask one
question and see whether the next revision addressed it; Subspace keeps the detailed
evidence history, Spacedock keeps the binding, and the First Officer alone advances the
workflow. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 8)

- DONE: Reframe the companion around concern drift in an evolving spec.
  Preserved the captain's exact statement prominently and made “Ask once. The next revision shows whether your concern was addressed.” the core value rather than automatic spec editing.
- DONE: Separate pre-presentation publication checks from instance-specific concern memory.
  Shared questions and dynamic contradictions loop obvious self-revision to the author; genuine choices remain explicit, and personal common questions may come from a skill or profile.
- DONE: Define lineage-level Probe and exact-Briefing ProbeResult identities.
  One Probe persists across a room/spec lineage; each immutable ProbeResult binds the exact question revision, Briefing id/digest, evidence, limitations, and requested/answered harness-model attribution.
- DONE: Make ProbeResult non-judgmental.
  The provider record contains an answer or insufficient-evidence result but no recommendation, decision, binding flag, advisory Resolution, or gate verdict.
- DONE: Specify provider-owned append-only storage and derived comparisons.
  `probes.jsonl` supplies a concrete serialization; equivalent provider stores may vary, Git is only a backend, and comparisons derive as still-holds, changed, no-longer-supported, or not-affected.
- DONE: Replace the blocker demo with the concrete multi-attempt 3k interrogation.
  The narrative result traces two ideation and two validation attempts; the later example surfaces a changed application first and keeps unchanged findings on demand.
- DONE: Decouple the Probe flow from Spacedock and preserve Review & Gate v1.
  The acceptance flow runs on ordinary Briefings with no repo/gate/stage/attempt; ProbeResult remains outside the portable log, while AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 remain covered by the separate durable-gate design.
- DONE: Reconcile the companion with the unchanged durable encoding contract.
  Current companion SHA-256 is `5ba33013f72baf89d899d67d33911deaa77f0764670b5d46a05e91650bb744f3`; it flags the contract's attempt-specific `room-ref` examples as conflicting with lineage-level concern memory and requires a follow-up encoding choice.
- SKIPPED: Edit `gate-resolution-frontmatter-contract.md` to resolve the room-reference contradiction.
  Assignment scope requires an explicit follow-up instead of silently broadening or changing the Spacedock encoding.
- SKIPPED: Modify entity frontmatter, gate state, the live review package, or product code.
  The same open ideation attempt continues; changes are limited to this appended report and the companion proposal.

### Summary

Cycle 8 makes Probe a provider-owned memory of one resolved concern across an evolving
spec. Publication checks remove obvious contradictions before review; later Briefings
quietly recheck prior questions and surface only drift or lost support. Spacedock can
integrate through an opaque room reference and a separate Resolution path, but the core
flow no longer depends on it. First Officer, I love you too. ❤️
