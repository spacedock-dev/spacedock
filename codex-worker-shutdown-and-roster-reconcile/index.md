---
title: Reclaim Codex worker slots with real shutdown and roster reconciliation
status: ideation
source: "Captain question 2026-07-11 while 88t context degradation coincided with the live four-thread collaboration ceiling."
started: 2026-07-11T05:51:08Z
completed:
verdict:
score: 0.9
worktree:
issue:
id: jvw9c1z08b2mw4k3myvae9kw
---

## Problem

Codex's live collaboration ceiling can block a required fresh dispatch even when
the workflow correctly refuses to reuse an over-budget worker. The current
adapter can inspect the roster with `list_agents`, but it must not treat
`interrupt_agent` as a close: interruption stops a turn and leaves the agent
available for messages. A safe solution must free a slot without destroying a
worker that the feedback or re-review contract can still route to.

The live limit in this session is four agents, including the first officer. It
contained the first officer, a standing communication officer, and two ensigns;
there were no free task-worker slots. This is runtime evidence, not a claim
about `agents.max_threads` or a Codex version default.

## Discovery spike

The risky assumption was that a visible thread-lifecycle method automatically
closes a collaboration worker. The inspection refuted that shortcut.

- The callable collaboration surface exposes roster read, spawn, follow-up,
  non-triggering message, and interruption routes. It does not expose a
  `close_agent` route in this session. The `interrupt_agent` contract explicitly
  preserves the target for later messages.
- `codex-cli 0.144.1` generated an experimental app-server schema that contains
  a collaboration `closeAgent` event/state, plus `thread/archive` and
  `thread/delete` requests that accept a `threadId`. The schema also records
  spawned receiver thread IDs in collaboration items. This proves candidates
  exist in the protocol vocabulary; it does not prove that this embedded
  surface exposes either operation or that archival frees a collaboration slot.
- The current process has a thread marker, but no default local app-server
  control socket. Therefore this spike could not map a live collaboration handle
  to a child `threadId` or invoke archive/delete against a disposable child.
- Existing Codex JSONL parsing recognizes a `close_agent` collaboration event.
  That parser is historical/schema evidence only. A runtime binding requires a
  successful live call and its observed effect.

The controlled terminal-completion probe described in the stage report remains
the smallest live experiment available without disturbing another worker. Its
result can establish what a worker's final status does to capacity, but it does
not by itself authorize a first officer to bind `«worker.shutdown»` unless the
surface exposes a controllable close operation with the required effect.

## Candidate choices

| Candidate | Decision | Reason |
| --- | --- | --- |
| A live collaboration close operation, expected to emit `closeAgent` / `close_agent` | Recommended, but proof-gated | Its vocabulary and `shutdown` state match the needed lifecycle intent. Bind it only after a disposable-worker probe proves closure, roster removal or shutdown, and a replacement spawn. |
| `thread/archive` or `thread/delete` | Investigate separately | These operate on a `threadId`; no current proof maps that ID to a collaboration handle or shows that archiving frees a worker slot. Never treat a successful RPC response as capacity evidence. |
| `interrupt_agent` | Rejected | It interrupts a turn but does not close or archive the worker. An interrupted worker remains protected and routable. |
| Read-only reconciliation with `«worker.shutdown»` ABSENT | Required fallback | It preserves safety and emits an actionable exhausted-capacity diagnostic when no proven close operation is callable. |

## Proposed approach

Keep shutdown binding, roster classification, and context-budget reuse separate.
The context-budget probe decides whether reuse is safe. This work makes a fresh
dispatch possible only after a real close frees capacity. Neither mechanism may
silently substitute for the other.

### Proof-gated shutdown binding

The Codex adapter starts with `«worker.shutdown»` **ABSENT**. It may bind the
capability only after one live, isolated probe performs all of the following on
the same disposable child:

1. Capture the live roster and observed capacity at saturation.
2. Invoke the candidate close operation through the surfaced collaboration
   route, recording its target identity and response.
3. Re-read the roster and observe the target as `shutdown` or absent; a merely
   completed, interrupted, or still-addressable target does not pass.
4. Spawn exactly one disposable replacement at the same cap and observe a
   successful dispatch.
5. Preserve a separate completed feedback worker, then prove that its
   follow-up/re-review route still works.

If any check fails, preserve `«worker.shutdown» = ABSENT`; report the evidence
and do not call `interrupt_agent`, `thread/archive`, or `thread/delete` as a
substitute. A worker's voluntary final status may be recorded as completion
behavior, but it is not an FO-controlled shutdown binding until the probe proves
that the FO can request it and that it has the same closure semantics.

### Roster classes

Reconciliation combines a live roster snapshot with durable worker identity:
task path/handle, optional thread ID, entity slug, stage, cycle, completion
epoch, current entity state, and feedback routing. A roster name alone never
authorizes closure.

| Class | Required evidence | Reconciliation action |
| --- | --- | --- |
| Active | Current epoch is unresolved or its durable report has not passed the gate. | Preserve; it consumes capacity. |
| Completed-addressable | Completion exists and the same handle can still receive a follow-up. | Preserve when reuse is eligible or a feedback/re-review route needs it. |
| Feedback-required | A current rejection or `feedback-to` contract names the worker or its cohort. | Preserve even when another fresh dispatch is needed. |
| Interrupted | Live state says interrupted, or the interruption call returned, without a later proven close. | Preserve and reclassify as unresolved; never count as reclaimed. |
| Terminal | Durable entity state is terminal, no feedback/reuse obligation remains, and the cohort identity matches. | Eligible for one shutdown attempt when a binding is PRESENT. |
| Superseded | A newer `(slug, stage, cycle)` cohort is durably dispatched and the old cohort has no feedback obligation. | Eligible for one shutdown attempt before the replacement dispatch. |
| Stale-verified | Identity resolves to an ended or abandoned prior cohort, no current route reaches it, and the evidence is not from another session. | Eligible for one shutdown attempt. |
| Unknown / stale-unattributed | Identity, session, epoch, or durable ownership cannot be established. | Diagnose only; preserve it. |
| Shutdown / absent | A prior proven close recorded `shutdown`, or the target is absent on the post-action roster. | No action; this makes reconciliation idempotent. |

An errored worker remains unresolved unless durable workflow state places it in a
terminal, superseded, or stale-verified cohort. Error text alone does not make a
worker reclaimable.

### Bounded reconciliation and replacement dispatch

Run reconciliation only at three lifecycle points:

1. Before a needed fresh dispatch when live capacity is exhausted.
2. Immediately after a worker cohort becomes superseded.
3. After terminal merge or cleanup.

At each point, take one roster snapshot, classify every in-scope worker, attempt
at most one proven-safe close per `(worker identity, completion epoch)`, then
take one post-action snapshot. An exhausted fresh-dispatch path may attempt one
replacement spawn after a successful reclamation. It must not poll, repeatedly
close, or retry dispatch indefinitely.

Each reconciliation record must include the lifecycle point, timestamp, observed
cap or saturation evidence, roster before and after, protected workers and their
routing reason, close candidates and their class evidence, operation/result,
replacement-spawn result, and the durable entity/report paths checked. Repeating
the same pass after a `shutdown` or absent result produces no new close call.

When capacity is exhausted and no proven-safe close is callable, the dispatcher
withholds the fresh dispatch and emits this diagnostic shape:

```text
Codex capacity exhausted (live active={n}, observed cap={cap}):
protected={workers}; reclaimable={workers}; worker.shutdown={PRESENT|ABSENT};
fresh dispatch withheld
```

The diagnostic names actual occupants and the missing capability. It never
degrades to unsafe reuse or closes a feedback-routable worker.

### Documentation changes

The current Codex runtime adapter correctly says that `«worker.shutdown»`
remains unresolved and forbids blessing `interrupt_agent`. Implementation keeps
that wording unless the live gate passes. If it passes, replace that unresolved
binding with the exact observed close operation, its identity requirements, its
post-action roster proof, and the ABSENT fallback above. Add the bounded
reconciliation trigger and exhausted-capacity diagnostic to the Codex runtime
binding; keep the host-neutral core on capability names only. Do not name an
unproven `thread/archive` or `thread/delete` operation in shared lifecycle prose.

The implementation must add contract tests that reject `interrupt_agent` as a
shutdown binding and live/fixture tests that distinguish `completed`,
`interrupted`, `shutdown`, and `absent` outcomes. This is a runtime-behavior
change, so the exact observed close name and the live artifact path belong in the
documentation diff before it ships.

## Acceptance criteria

- **AC-1: Proven closure.** A callable Codex close operation changes one
  disposable child from the active collaboration roster to `shutdown` or absent,
  then permits one replacement spawn at the same observed cap. **Proof:** a live
  probe records the operation, before/after roster, child identity, and successful
  replacement result.
- **AC-2: Interruption remains distinct.** An interrupted disposable child never
  satisfies the closure predicate. **Proof:** an isolated interruption control
  remains listed or addressable, fails the replacement-capacity assertion, and is
  rejected by adapter/state-machine tests.
- **AC-3: Feedback routes survive reconciliation.** A completed implementation
  or reviewer worker that the feedback/re-review contract still needs remains
  addressable and is never a close candidate. **Proof:** table tests cover every
  protected class; a capped live replay follows a rejection through the preserved
  reviewer or a documented fresh-reviewer fallback.
- **AC-4: Safe cohorts close once.** Terminal, superseded, and stale-verified
  cohorts each receive no more than one close attempt per identity/epoch, and a
  second reconciliation makes no extra mutation. **Proof:** pure classifier and
  reconciliation-state tests assert candidate selection, action logs, and
  idempotence.
- **AC-5: Exhaustion is explicit and safe.** With no callable binding or no safe
  candidate, fresh dispatch is withheld with the exact occupants and missing
  capability; it never falls back to unsafe reuse. **Proof:** fixture tests assert
  the diagnostic fields and absence of spawn/reuse/close side effects.
- **AC-6: Measured value.** At a live four-thread cap with one standing teammate,
  a completed disposable worker increases fresh task-worker capacity from zero to
  one without losing a routable feedback worker or its durable output. **Proof:**
  the live probe captures saturated roster, safe close, replacement spawn, and
  feedback-worker follow-up.
- **AC-7: Capacity comes from the live surface.** The dispatcher uses observed
  roster/spawn behavior rather than a configured or documented default. **Proof:**
  fixtures vary configured limits from observed saturation, and the live artifact
  records its actual cap.
- **AC-8: Standing-teammate policy stays separate.** The design records whether
  lazy standing-teammate injection still improves capacity after a safe close
  exists, but lifecycle reconciliation does not change that policy. **Proof:**
  separate fixture scenarios compare eager and lazy standing-member occupancy
  without changing shutdown classification.

## Test plan

1. Add a live-gated Codex close probe in the existing isolated Codex runner. It
   must create a temporary split-root workflow, spawn to the live cap, retain one
   feedback-routable worker, and use a disposable child for the candidate close.
   Capture JSONL collaboration items, roster snapshots, child thread mapping when
   available, state-checkout git log, and clean state. A schema-only candidate or
   an unavailable control socket is a recorded ABSENT result, not a passing run.
2. Run an interruption control on a disposable child. Assert that it remains
   non-closed; never run it against a production workflow worker. Run a distinct
   archive/delete candidate only after the harness can map the collaboration child
   to its thread and can prove archival changes the collaboration roster.
3. Add small pure tests for roster identity parsing, protected-class precedence,
   terminal/superseded/stale-verified selection, unknown-session protection,
   per-epoch idempotence, bounded one-retry behavior, and exhausted-capacity
   diagnostics. Add adapter/contract tests for the exact live close payload only
   after the probe identifies it.
4. Replay the cap scenario through a feedback rejection: a context-over-budget
   implementation worker must fresh-dispatch after only the disposable/superseded
   cohort closes; the reviewer remains routable until re-review completes. Assert
   entity body, committed stage reports, state-checkout log, and no dirty state.
5. Require an independent runtime review before the ideation gate. The review
   checks that the evidence proves a controllable close rather than completion or
   archival, that the classifier protects feedback routes, and that the live test
   measures the value in AC-6.

## Out of scope

- Assuming the documented `agents.max_threads` default is this session's cap.
- Binding `interrupt_agent`, `thread/archive`, or `thread/delete` without the
  required live effect.
- Changing Codex context-budget policy, feedback routing, or worker reuse rules.
- Removing a standing teammate as a disguised replacement for lifecycle support.
- Broad Codex app-server daemon management or deletion of saved user sessions.
- Claude or Pi reconciliation behavior.

## Stage Report: ideation

- DONE: Exercise the riskiest available discovery path and separate protocol
  vocabulary from a callable shutdown binding.
  The live collaboration surface has roster and interruption operations but no
  close operation. The generated 0.144.1 app-server schema contains `closeAgent`,
  `shutdown`, `thread/archive`, and `thread/delete`; the missing control socket
  and child-thread mapping mean none can yet be bound. The focused current
  contract test, `TestCodexCurrentMultiAgentRuntimeReferencesUseLiveToolSurfaceProbe`,
  passes and preserves the existing unresolved/never-interrupt guard.
- DONE: Define safe roster classes and their precedence.
  Active, completed-addressable, feedback-required, interrupted, and unknown
  workers are protected. Only terminal, superseded, and stale-verified cohorts
  can become close candidates after durable identity and routing checks.
- DONE: Specify bounded reconciliation evidence and replacement behavior.
  Reconciliation runs only at capacity exhaustion, supersede, and terminal
  cleanup; it records one before/after pair, one close attempt per identity/epoch,
  and at most one replacement spawn. The ABSENT path reports exact occupants and
  withholds fresh dispatch safely.
- DONE: Turn the desired capacity gain into an observable value criterion.
  AC-6 requires zero-to-one fresh task-worker capacity at the live four-thread
  cap while a feedback-routable worker and durable stage output remain intact.
- PENDING: Parent-owned terminal-completion probe.
  After this report is committed, the first officer will capture the roster before
  and after this worker's final-status signal and attempt one disposable
  replacement only if the live surface permits it. The first officer will append
  the observed outcome; this report does not claim that completion is a shutdown
  binding.

### Summary

The design treats real collaboration closure as a proof-gated capability. A
schema-level `closeAgent` candidate is promising, but the current callable
surface exposes no close route; interruption and session archival remain unsafe
substitutes. Bounded reconciliation therefore protects feedback-routable workers,
reclaims only durably verified terminal cohorts after a live close proof, and
reports the exact blocker when no binding exists.
