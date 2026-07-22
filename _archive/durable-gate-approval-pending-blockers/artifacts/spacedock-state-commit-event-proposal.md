# Proposal: Commit-Derived Spacedock Events and a Read-Only Zaphod Projection

Status: Draft  
Date: 2026-07-13

## Decision

Spacedock should derive its durable event stream from the Git history of the
workflow state checkout. It should not maintain a second append-only ledger.

A deterministic projector will convert state commits into versioned event
envelopes. A pure reducer will replay those envelopes into task timelines and
current state. Zaphod will consume that projection through a native adapter and
render it without changing workflow state.

Runtime facts that do not exist in state commits—live worker handles, pane
presence, or an open review UI—belong to a separate observational overlay. If
one of those facts must survive restart, Spacedock must record it in a state
commit. The projector remains the sole source of durable events.

This design gives Spacedock one durable authority: state Git history.

For gate state, the tree authority inside each commit is the entity's versioned
top-level `gates` frontmatter collection defined by
[`../gate-resolution-frontmatter-contract.md`](../gate-resolution-frontmatter-contract.md).
The event stream is derived from successive committed values of that mapping; it is
not a replacement persistence location.

Portable review semantics remain those of Review & Gate v1 at
`spacedock-subspace` commit `bd17bdb23318f815d17a1d10ea2a6d39ab449520`,
`docs/review-and-gate.md` blob
`14f3eb91ec85bfcc08bb3330c21b94cc77f4529f`. One immutable Briefing is one
decision opportunity and owns one ordered review log. Stage, round, event identity,
Spacedock gate-attempt identity, JCS Briefing digest, selection, and application below
are Spacedock projection/index fields, not additions to the portable contract. A
Spacedock gate attempt may select several immutable Briefings while open.

## Goals

- Reconstruct the same event stream from the same state history.
- Derive current workflow state through a pure, replayable reducer.
- Represent approval before dispatch without inventing another lifecycle
  stage.
- Correlate tasks, stage runs, dispatch attempts, workers, gates, and evidence
  across Spacedock, Subspace, Roborev, and Zaphod.
- Let Zaphod show durable state and live observations without becoming a
  workflow authority.
- Detect stale, missing, duplicate, concurrent, and rewritten history.

## Non-goals

- Zaphod will not advance tasks, resolve gates, dispatch workers, or repair
  state.
- The event stream will not copy prompts, transcripts, review bodies, or other
  sensitive prose. Events will carry identifiers, digests, and references.
- `m3` will not implement the ledger. It will supply stable gate and pane
  correlation metadata.
- This proposal does not make a tab title, CWD, pane geometry, or plugin path
  substring authoritative identity.

## Sources of truth

The system has three distinct sources:

1. **State commits** are the durable workflow authority. They record task
   creation, field changes, stage reports, feedback cycles, the entity's structured
   `gates` frontmatter collection, completion, verdicts, and explicit receipts.
2. **Product commits** prove code integration and merge ancestry. State may
   reference an exact product commit, but state cannot prove that commit exists
   or has the claimed ancestry without inspecting the product repository.
3. **Runtime observations** report current workers, panes, processes, and open
   review surfaces. They can disappear and can disagree with durable state.

The durable reducer consumes only state-derived events. A view reducer may
overlay product verification and runtime observations for display.

## Stable identity

The contract must keep these identities separate:

| Identity | Source | Meaning |
|---|---|---|
| `workflow_id` | Workflow frontmatter | Stable identity independent of path or remote URL |
| `task_id` | Task frontmatter | Existing stable Spacedock task identity |
| `event_id` | Projector | One normalized event derived from one state commit |
| `stage_run_id` | Projector | One entry into a task stage |
| `dispatch_attempt_id` | Dispatch preparation | One concrete attempt to spawn a worker |
| `worker_id` | Runtime receipt | Worker handle returned by a successful spawn |
| `gate_id` | Entity `gates.records[]` | One logical entity/stage gate across adjudication attempts |
| `gate_attempt_id` | Entity `gates.records[].attempts[].id` | One stable Spacedock adjudication session, possibly spanning Briefing snapshots |
| `briefing_id` | Review & Gate `Briefing.id` | One immutable portable decision opportunity and its separate ordered log |

`stage_run_id` should derive from the event that enters the stage:

```text
stage_run_id = hash(workflow_id, task_id, stage_entered_event_id)
```

`gate_attempt_id` is minted by Spacedock and stays stable while its open adjudication
advances across immutable `briefing_id` snapshots. The current entity binds the attempt
to `gate_id` and stage and stores only its current Briefing id/digest; state Git commits
retain prior pointer values and Subspace retains the full snapshots/logs. A changed
question, artifact revision set, decision opportunity, lens presentation, or reviewed
evidence gets a new Briefing id; it does not necessarily get a new gate-attempt id or
imply `revise`. A closed result followed by gate re-entry normally gets a new
gate-attempt id.

`dispatch_attempt_id` cannot derive from a stage alone because one stage run may
need retries or replacement workers. Dispatch preparation must mint it before
the spawn call and pass it into the package, worker environment, pane title,
result record, and any later receipt.

## Canonical event envelope

The projector should emit JSON Lines. Every line uses one versioned envelope:

```json
{
  "schema": "spacedock.state-event/v1",
  "event_id": "se1:8c6b…",
  "kind": "task.stage_entered",
  "workflow_id": "wf1:…",
  "task_id": "kjhq0t2h6drse6b32cqybggv",
  "task_slug": "managed-tab-safety-session-integration",
  "stage": "implementation",
  "stage_run_id": "sr1:4e91…",
  "recorded_at": "2026-07-13T12:30:05Z",
  "source": {
    "ref": "spacedock-state/agent-rail-dev",
    "commit": "b24f470…",
    "parents": ["9a76e74…"],
    "path": "managed-tab-safety-session-integration.md",
    "ordinal": 2
  },
  "data": {
    "from": "ideation",
    "to": "implementation"
  }
}
```

The projector computes `event_id` from the schema version, workflow identity,
source commit, path, ordinal, event kind, and normalized payload. Replaying the
same history therefore produces byte-identical IDs and payloads.

`recorded_at` comes from the state commit's committer timestamp. The envelope
may preserve an authored or externally observed timestamp inside `data`, but
the reducer orders durable events by state history, not wall-clock time.

## Events derived from tree changes

The first contract should recognize these durable events:

- `task.created`
- `task.field_changed`
- `task.stage_entered`
- `task.worktree_assigned`
- `stage.report_recorded`
- `feedback.cycle_recorded`
- `gate.attempt_opened`
- `gate.briefing_selected`
- `gate.attempt_closed`
- `gate.resolution_recorded`
- `task.verdict_recorded`
- `task.completed`
- `task.merge_reference_recorded`
- `task.archived`

The projector should parse structured frontmatter and committed decision-log
records. In particular, it parses the complete old/new entity `gates` YAML trees; a
temporary review log is evidence input, never current Spacedock state. It should treat
narrative prose as an opaque artifact. For example,
`stage.report_recorded` should carry the heading, stage, digest, and source
lines; it should not pretend to understand free-form claims unless the Stage
Report grammar supplies structured checklist results.

Gate resolution becomes durable only when the exact binding Resolution is committed in
an entity `gates.records[].attempts[].resolution` node. Recording it does not change `status` or
dispatch a worker. A temporary Subspace briefing, external review log, or FO-authored
presentation file is not current gate state. The mapping carries logical gate id,
separate Briefing and gate-attempt ids, stage, canonical Briefing digest, complete
portable Resolution, application intent, blockers/hold, and consumption state. The
current tree directly retains all gates and attempts plus one current (open) or frozen
resolved (closed) Briefing binding per attempt; Git replay supplies prior pointer
transitions and Subspace supplies full presentation history.

The adopted Resolution is the first one attributed to the authorized approver identity
that workflow tooling supplied externally, and it closes the attempt only when its
`briefing` equals the attempt's exact current-Briefing pointer. Earlier Resolutions from
other actors remain advisory in that one-Briefing provider log. Entries never carry
across Briefing logs, and cross-Briefing `includes` is invalid. Review & Gate permits
reasonless `approve`; `revise`/`hold` require a nonblank reason or an included earlier
same-Briefing Annotation. Spacedock's stricter reason policy for FO conn-made approvals
is authoring policy, not portable validation.

`feedback.cycle_recorded` is the durable route edge between a rejected gate
result and the stage run that receives its rework. Its structured payload must
carry the feedback cycle, rejected `gate_attempt_id`, source stage and
`stage_run_id`, target stage, and the routed finding reference and digest. When
the matching `task.stage_entered` event is projected, the projector binds that
new target `stage_run_id` to the feedback event. The relationship is explicit;
the reducer must not infer rework merely from a repeated stage name, report
prose, or the current workflow definition's `feedback-to` value.

It is derived when the authoritative frontmatter's pending feedback application is
consumed with the matching target-stage transition. This extends an existing event
rather than adding a second stage-result event.
The rejected result remains `gate.resolution_recorded`; the feedback event says
where that result was routed. A cycle-3 escalation records no target stage run
and therefore cannot masquerade as dispatched rework.

## Explicit receipts for external side effects

Tree changes cannot prove that an external side effect occurred. A transition
commit written before `spawn_agent` proves dispatch intent, not spawn success.
Likewise, a gate artifact proves that a review can be opened, not that its UI
opened.

When a side effect must become durable, Spacedock should create a small state
commit with machine-readable commit trailers. The projector will convert those
trailers into ordinary envelopes. No second ledger file is required.

Example trailers:

```text
Spacedock-Receipt: dispatch.spawned/v1
Spacedock-Workflow-ID: wf1:…
Spacedock-Task-ID: kjhq0t2h6drse6b32cqybggv
Spacedock-Stage-Run-ID: sr1:…
Spacedock-Dispatch-Attempt-ID: da1:…
Spacedock-Worker-ID: codex:/root/spacedock_ensign_…
```

Candidate receipt events are:

- `dispatch.prepared`
- `dispatch.spawned`
- `dispatch.failed`
- `worker.completion_signaled`
- `worker.superseded`
- `gate.opened`
- `gate.open_failed`
- `merge.verified`
- `worker.cancelled`

The recorder command must create a state commit; it must not append to a
parallel journal. Repeated calls with the same receipt identity must be
idempotent. The command should reject conflicting payloads for an existing
receipt ID.

Spacedock need not persist every observation. A product may choose a smaller
durable contract and leave `gate.opened` or worker liveness in the runtime
overlay. The UI must then label those facts as observations.

## Git history and concurrency

The projector must support the state branch's commit graph without double
counting merged work.

Recommended projection rules:

1. Walk every commit reachable from the canonical state ref in deterministic
   topological order. Break ties by commit SHA.
2. For a normal commit, compare its tree with its first parent.
3. For a merge commit, emit tree-derived events only for paths whose result
   differs from every parent. Parent commits already account for inherited
   changes; the merge contributes only conflict resolutions or new merge-time
   edits.
4. Read recognized receipt trailers from every reachable commit, including
   commits introduced through a merge.
5. Sort multiple events from one commit by task ID, path, event precedence,
   and normalized payload before assigning ordinals.
6. Detect concurrent contradictory transitions for one task. Mark the reduced
   task conflicted until a later state commit resolves it; never select a
   winner by timestamp.

A projection cursor is `{canonical_ref_tip, event_id}`. If the canonical ref no
longer descends from the stored tip, the consumer must report history rewrite,
find a common ancestor, and rebuild. It must not silently continue from an
unrelated history.

## Reducer

The reducer is a pure function:

```text
reduce(previous_snapshot, ordered_event) -> next_snapshot
```

Each task snapshot should contain:

- durable stage and `stage_run_id`;
- worktree and referenced product head;
- latest Stage Report reference and digest;
- current gate attempt, decision, and artifact digest;
- active feedback route, including its source gate result, cycle, source and
  target stage runs, routed finding reference, and resolution state;
- dispatch attempts and durable receipts, if recorded;
- completion, verdict, and merge reference;
- conflicts, stale approvals, and invalid transitions.

The reducer must never inspect live workers or mutate workflow files.

### Approved but not dispatched

`approved_pending_dispatch` is a computed gate condition, not a lifecycle
stage. It holds when:

1. the selected closed Spacedock attempt's Resolution approves its exact current
   Briefing and reviewed digest;
2. no later event invalidates that digest or decision; and
3. no later `task.stage_entered` or durable `dispatch.spawned` receipt consumes
   the approval.

Changing the reviewed artifact makes the approval stale. A new gate attempt
must resolve the changed digest.

To preserve this state, Spacedock must stop treating approval as an instruction
to commit the next stage before a worker exists. The revised sequence is:

1. Under the already-open attempt's current-Briefing compare-and-swap, store the exact
   binding Resolution, close the attempt, create `application.state: pending`, and leave
   `status` unchanged.
2. After confirming the stored digest is current, every blocker is satisfied, and no
   execution hold is active, prepare a package and commit a minted
   `dispatch_attempt_id` with `state: prepared` without consuming approval.
3. Call the idempotent or authoritatively queryable runtime spawn boundary.
4. After success, atomically commit the stage transition, consumed state, and spawn
   receipt.
5. If the process dies between steps 3 and 4, reconciliation reports a live
   orphan; it does not fabricate state.

A failed spawn may receive a committed failure receipt while the approval
remains pending for retry.

### Rejected result routed to rework

A rejected Resolution is a durable completed gate result even when workflow
routing changes the task's current lifecycle stage. On
the entity records Review & Gate `decision: revise` plus Spacedock
`application.action: feedback`; the captain-facing word “rejected” is a workflow
projection, not a resurrected portable `reject` decision. On
`feedback.cycle_recorded`, the reducer opens an active feedback route linked to
the rejected `gate_attempt_id`. On the matching `task.stage_entered`, it binds
the target `stage_run_id` and computes a `feedback_rework` condition alongside
the ordinary current stage.

For example, after validation rejects cycle 1 and routes to implementation,
the snapshot still has `stage: implementation`, while its structured context
identifies `condition: feedback_rework`, `feedback_cycle: 1`,
`feedback_from_stage: validation`, and the rejected gate attempt. Text status
must render the same distinction, such as `implementation (rework: validation
rejected, cycle 1)`. It must not collapse the row to bare `implementation` plus
an unchanged score.

The active route survives restart and remains attached through remediation and
the resulting re-review. A later committed gate Resolution supersedes it: a
pass closes the active route; another rejection closes the prior route and
opens the next cycle. Historical gate and feedback events remain replayable in
all cases. Re-entering a stage without a linked feedback event is an ordinary
stage run and must not be labeled rework.

## Runtime overlay and reconciliation

Runtime observations should use a separate schema, such as
`spacedock.runtime-observation/v1`. They may report:

- worker handle and liveness;
- pane and process identity;
- open Subspace review surface;
- AgentsView session identity;
- current Zellij session and stable tab ID.

The overlay must link to durable `workflow_id`, `task_id`, `stage_run_id`,
`dispatch_attempt_id`, or `gate_attempt_id`. It must not infer identity from a
title, CWD, geometry, or path substring.

Reconciliation compares reduced durable state with current observations and
emits diagnostics:

- `unconfirmed`: durable running state has no live worker observation;
- `orphaned`: a live worker has no durable spawn receipt;
- `stale`: a recorded worker or gate observation has expired;
- `superseded`: a later attempt owns the stage run;
- `conflicted`: concurrent state commits disagree.

Reconciliation never changes task state merely to repair the display. An
operator or explicit recovery command must authorize any state commit.

## Zaphod integration

The Zaphod WASM plugin should not read Git repositories or state files. A
native Zaphod adapter should consume Spacedock's versioned projection:

```text
spacedock events project --workflow-dir <dir> --after <cursor> --jsonl
spacedock events reduce  --workflow-dir <dir> --json
spacedock events watch   --workflow-dir <dir> --after <cursor> --jsonl
```

The adapter will combine durable snapshots with runtime observations, then
send tab-targeted rows through Zaphod's existing native-to-plugin boundary.
The plugin remains a read-only renderer and action surface.

The launcher must bind a tab to a workflow explicitly. A task worker should
receive stable IDs through its environment and pane metadata, including
`SPACEDOCK_WORKFLOW_ID`, `SPACEDOCK_TASK_ID`, `SPACEDOCK_STAGE_RUN_ID`, and
`SPACEDOCK_DISPATCH_ATTEMPT_ID`. Zaphod may display a title or emoji derived
from this binding, but those decorations remain non-authoritative.

The sidebar may render these projections:

- approved, pending dispatch;
- prepared;
- running;
- awaiting gate;
- feedback;
- failed;
- superseded;
- completed;
- merged;
- stale, unconfirmed, orphaned, or conflicted.

Each row should expose its source: durable state commit, product verification,
or runtime observation.

## Ownership

### Spacedock

- Define the event and receipt schemas.
- Project state commits and reduce events.
- Persist stable workflow, stage-run, dispatch-attempt, and gate-attempt
  identities.
- Expose project, reduce, watch, and reconciliation commands.
- Commit gate decisions and any selected external-side-effect receipts.

### Subspace and `m3`

- `m3` supplies stable task, stage, Spacedock gate-attempt, portable Briefing, actor,
  attempt sequence, and digest metadata on gate covers and Resolution records.
- The wrapper returns actual open/result observations with those IDs.
- Spacedock decides which observations become state receipts.

### Zaphod

- Consume the supported Spacedock stream through a native adapter.
- Correlate rows through stable IDs.
- Render durable state, observations, and discrepancies.
- Never write workflow state or claim dispatch success from package creation.

### Roborev and workers

- Carry the same task, stage-run, dispatch-attempt, product-head, and evidence
  IDs in reports and review records.
- Store evidence by reference and digest rather than copying review prose into
  events.

## Delivery plan

### Phase 1: Freeze the commit-derived contract

1. Add stable `workflow_id` to commissioned workflow metadata.
2. Specify `spacedock.state-event/v1` and deterministic event IDs.
3. Implement Git DAG projection and the pure reducer.
4. Persist stable Spacedock gate attempts, one current or frozen resolved Briefing
   binding per attempt, and exact adopted Resolutions in the entity's `gates`
   collection; keep full Briefing/log/lens/assessment history in Subspace.
5. Persist structured feedback-cycle route edges and bind rework stage runs to
   their rejected gate attempts.
6. Derive `approved_pending_dispatch` and `feedback_rework` without runtime
   instrumentation.

### Phase 2: Add dispatch identity and selected receipts

1. Mint `dispatch_attempt_id` before runtime spawn.
2. Pass all stable IDs through packages, environments, pane metadata, reports,
   and Roborev evidence.
3. Commit spawn success, spawn failure, completion signal, and supersession
   receipts.
4. Add restart reconciliation without automatic state repair.

### Phase 3: Build the Zaphod projection

1. Add the native event-stream adapter.
2. Bind one managed tab to one explicit workflow identity.
3. Render durable attempt state and stale-state diagnostics.
4. Add the managed-tab indicator and provenance view without making either an
   ownership input.

After Phase 1 freezes envelope and reducer fixtures, Phase 2 boundary
instrumentation and Phase 3 Zaphod work may proceed in parallel.

## Acceptance tests

1. **Determinism:** two projections of the same state ref produce byte-identical
   events, IDs, order, and reduced snapshots.
2. **Resume:** replay after a stored cursor emits no duplicate event.
3. **Merge:** concurrent state branches merged in either ordinary supported
   shape produce each parent event once and emit only genuine merge-resolution
   changes at the merge commit.
4. **Approval:** a committed approval with no consumed dispatch reduces to
   `approved_pending_dispatch`; changing the artifact digest makes it stale.
5. **No overclaim:** a package-build commit cannot reduce to spawned or running
   without a spawn receipt or live observation.
6. **Retry:** two dispatch attempts in one stage run retain separate IDs;
   superseding one cannot rewrite the other.
7. **Rewrite:** a non-descendant canonical ref causes a visible rebuild warning.
8. **Reconciliation:** missing and orphaned workers produce diagnostics without
   changing workflow state.
9. **Projection safety:** Zaphod renders the stream but cannot advance, resolve,
   dispatch, or merge a task.
10. **Privacy:** event fixtures contain identifiers, digests, and references but
    no prompt, transcript, or review body.
11. **Rejected rework:** a validation rejection routed to implementation
    projects the current stage as implementation plus durable
    `feedback_rework` context naming the validation gate and cycle; rebuilding
    after restart produces the same text and JSON, while an ordinary repeated
    implementation stage does not acquire that context.
12. **Physical authority:** recording approve or revise appends the exact binding
    Resolution to the entity `gates` collection without changing `status` or producing
    a dispatch receipt; a cold current-entity read enumerates every logical gate,
    attempt, current or frozen Briefing binding, selection, and latest application
    state, while Git replay reproduces prior pointer transitions.
13. **Snapshot evolution:** lens or reviewed-input changes advance one open attempt
    by replacing its current Briefing binding without a `revise` event. Current
    frontmatter embeds no revision list; Git reconstructs pointer history and Subspace
    retains full snapshots/logs/lenses/assessments. Only a Resolution for the exact
    current Briefing closes and freezes it; cross-Briefing `includes` fails.

## Open questions

1. Should the canonical state ref permit arbitrary merge commits, or should
   Spacedock serialize state commits through one writer? The projector above
   supports merges, but serialization would simplify ordering.
2. Which external receipts merit a state commit in the first release? Spawn
   success and failure are the minimum for a durable dispatch ledger.
3. Should `merge.verified` record only the product commit reference, or also
   the exact ancestry query and result digest?
4. How long should Zaphod retain runtime observations after their source
   disappears?
5. Should snapshots be cached beside the state checkout? Any cache must remain
   disposable and rebuildable from Git history.

## Recommendation

Implement Phase 1 in Spacedock before adding a Zaphod ledger reader. Treat the
committed entity `gates` collection as directly readable gate truth, the state commit
graph as transition history, and the event projector/reducer as replayable public interpretations. Add
runtime receipts only for external facts that must survive restart. Keep Zaphod a
read-only projection of that truth.
