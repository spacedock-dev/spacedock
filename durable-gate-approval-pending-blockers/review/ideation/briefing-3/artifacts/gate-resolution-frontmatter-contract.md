# Gate Resolution frontmatter contract

Status: proposed v1 contract for `3k` ideation review
Date: 2026-07-18

## Authoritative sources and storage boundary

The portable authority is the complete
`../spacedock-subspace/docs/review-and-gate.md` contract at
`spacedock-subspace` commit `bd17bdb23318f815d17a1d10ea2a6d39ab449520`, blob
`14f3eb91ec85bfcc08bb3330c21b94cc77f4529f`, especially §1 “Scope and
ownership,” §2 “Model,” §3 “Decisions,” §5 “Review entries,” §7 “Review log,”
and §8 “Versioning and serialization.” The entity tree below indexes selected
portable objects; it does not change Review & Gate v1.

The companion [`gate-review-probes.md`](gate-review-probes.md) defines the first-use
question, later-revision comparison, and room-local result flow. This contract encodes
only the optional opaque provider room reference and durable Spacedock gate binding. It
does not copy Briefings, logs, probes, results, citations, responder attribution, or
deltas into entity frontmatter.

Closed Spacedock PR #474 (`iamcxa/status-apply-gate`, commits
`685fe7bcda4a51b8e2c06da52e80c079f62ac8e0` through
`5dee22831856db65db2acfefee0849c5f990f5d1`) supplied the retained physical direction:
the workflow binding belongs in binary-owned entity frontmatter. `3k` removes PR #474's
coupling between recording a decision and changing `status`.

The corrected identity hierarchy is:

```text
logical Spacedock gate
  -> stable Spacedock adjudication attempts (directly retained in the entity)
       -> current immutable portable Briefing reference while open
       -> frozen resolved Briefing reference + exact binding Resolution when closed
       -> mutable Spacedock workflow application, consumed separately
```

A Spacedock gate attempt is not a Review & Gate `Briefing`. It is a longer-lived
adjudication session. While open, adding/revising a lens or changing design/evidence can
replace its current Briefing pointer with a new immutable Briefing without implying
portable `revise` or creating a new attempt. The current entity does **not** duplicate
that attempt's full Briefing revision list. State Git commits preserve its prior
pointer/digest values, while Subspace retains the full Briefings, their separate logs,
lenses, assessments, and presentable deltas.

A binding Resolution closes the attempt only when it references the exact current
Briefing. The close commit freezes that Briefing id/digest (and stable Subspace room
reference when used) as `resolved-briefing`, removes the mutable `current-briefing`
slot, preserves the exact Resolution, and creates the Spacedock application. Re-entry
to the logical gate after a closed result is the normal new-attempt boundary. Closed
attempts never reopen.

## Physical representation

```yaml
gates:
  version: 1
  current:
    gate: gate:docs-dev:3k:validation
    attempt: gate-attempt:3k-validation-2
  records:
    - id: gate:docs-dev:3k:ideation
      stage: ideation
      current-attempt: gate-attempt:3k-ideation-2
      attempts:
        - id: gate-attempt:3k-ideation-1
          sequence: 1
          state: closed
          resolved-briefing:
            id: briefing:3k-ideation-1a
            digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
            room-ref: subspace-room:3k-gate-design
          resolution:
            type: Resolution
            id: resolution:captain-3k-ideation-1a
            briefing: briefing:3k-ideation-1a
            by: person:captain
            at: 2026-07-16T09:00:00Z
            decision: revise
            reason: Clarify the dispatch blocker contract.
            includes: []
          application:
            id: application:3k-ideation-1
            action: feedback
            target-stage: backlog
            state: consumed
            consumed-at: 2026-07-16T09:02:00Z
            blockers: []
            feedback:
              cycle: 1
              finding-ref: resolution:captain-3k-ideation-1a
              finding-digest: sha256:2111111111111111111111111111111111111111111111111111111111111111
        - id: gate-attempt:3k-ideation-2
          sequence: 2
          previous-attempt: gate-attempt:3k-ideation-1
          state: closed
          resolved-briefing:
            id: briefing:3k-ideation-2b
            digest: sha256:3222222222222222222222222222222222222222222222222222222222222222
            room-ref: subspace-room:3k-gate-design
          resolution:
            type: Resolution
            id: resolution:captain-3k-ideation-2b
            briefing: briefing:3k-ideation-2b
            by: person:captain
            at: 2026-07-17T09:00:00Z
            decision: approve
          application:
            id: application:3k-ideation-2
            action: advance
            target-stage: implementation
            state: consumed
            dispatch-attempt-id: dispatch:3k-implementation-1
            consumed-at: 2026-07-17T09:03:00Z
            blockers: []
    - id: gate:docs-dev:3k:validation
      stage: validation
      current-attempt: gate-attempt:3k-validation-2
      attempts:
        - id: gate-attempt:3k-validation-1
          sequence: 1
          state: closed
          resolved-briefing:
            id: briefing:3k-validation-1a
            digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
            room-ref: subspace-room:3k-gate-design
          resolution:
            type: Resolution
            id: resolution:captain-3k-validation-1a
            briefing: briefing:3k-validation-1a
            by: person:captain
            at: 2026-07-18T08:00:00Z
            decision: revise
            reason: The production coordinator is missing.
            includes: []
          application:
            id: application:3k-validation-1
            action: feedback
            target-stage: implementation
            state: consumed
            consumed-at: 2026-07-18T08:02:00Z
            blockers: []
            feedback:
              cycle: 1
              finding-ref: resolution:captain-3k-validation-1a
              finding-digest: sha256:4333333333333333333333333333333333333333333333333333333333333333
        - id: gate-attempt:3k-validation-2
          sequence: 2
          previous-attempt: gate-attempt:3k-validation-1
          state: closed
          resolved-briefing:
            id: briefing:3k-validation-2b
            digest: sha256:5444444444444444444444444444444444444444444444444444444444444444
            room-ref: subspace-room:3k-gate-design
          resolution:
            type: Resolution
            id: resolution:captain-3k-validation-2b
            briefing: briefing:3k-validation-2b
            by: person:captain
            at: 2026-07-18T10:30:00Z
            decision: approve
          application:
            id: application:3k-validation-2
            action: advance
            target-stage: done
            state: pending
            blockers:
              - id: blocker:production-coordinator
                kind: entity-stage
                ref: production-coordinator
                expected-revision: state:4f92c1d
                expected-state: done
                state: unsatisfied
                observed-revision: state:17ae4c0
                checked-at: 2026-07-18T10:30:00Z
                failure-code:
            execution-hold:
              id: hold:captain:3k-validation-2
              state: active
              by: person:captain
              at: 2026-07-18T10:30:00Z
              reason: Approval is durable; do not consume it until I release this hold.
```

This current entity directly contains two logical gates and two distinct adjudication
attempts for each. Every closed attempt freezes only the Briefing it actually resolved.
For example, Git history may show ideation attempt 2 advance from Briefing `2a` to `2b`
while open; the current tree stores only resolved `2b`. The opaque provider reference
`subspace-room:3k-gate-design` identifies one review-room/spec-lineage across all four
attempts. The provider owns the full Briefings, logs, lens/evidence changes, assessment
re-evaluation, and presentable deltas.

An open attempt has this lean shape:

```yaml
- id: gate-attempt:3k-validation-3
  sequence: 3
  previous-attempt: gate-attempt:3k-validation-2
  state: open
  current-briefing:
    id: briefing:3k-validation-3c
    digest: sha256:7555555555555555555555555555555555555555555555555555555555555555
    room-ref: subspace-room:3k-gate-design
```

It has no `resolution` or `application`. Earlier current pointers `3a` and `3b` are
recoverable from state commits and fully materialized in the stable Subspace room, not
repeated in current frontmatter.

## Field and layer contract

The `gates` tree and all fields except the copied Resolution are Spacedock-owned. The
YAML `resolution` node is a semantic transcription of the exact tagged-JSON portable
object, including valid additive v1 fields; YAML syntax does not define another
portable format.

| Field | Required | Meaning |
|---|---:|---|
| `gates.version` | yes | Integer `1`; unsupported versions fail closed. |
| `gates.current` | when an attempt is selected at the current gated stage | Pair of logical gate and Spacedock attempt ids. |
| `records[].id` / `stage` | yes | Stable logical entity/stage gate and Spacedock workflow stage. |
| `records[].current-attempt` | yes | Selected adjudication attempt for this logical gate. |
| `attempts[].id` | yes | Stable Spacedock adjudication-session identity; not a `Briefing.id`. |
| `attempts[].sequence` / `previous-attempt` | yes / after first | Spacedock ordering and chain within the logical gate. |
| `attempts[].state` | yes | `open` or `closed`; closed attempts never reopen. |
| `current-briefing` | iff open | Exact current Briefing id, Spacedock JCS digest, and optional stable `room-ref`. |
| `resolved-briefing` | iff closed | Frozen id/digest and optional room reference for the Briefing named by the adopted Resolution. |
| `attempts[].resolution` | iff closed | Exact adopted binding portable Resolution; `briefing` equals `resolved-briefing.id`. |
| `attempts[].application` | iff closed | One-use Spacedock application created when the attempt closes. |
| `application.action` / `target-stage` | yes / for advance or feedback | `advance`, `feedback`, or `none`, and its Spacedock target. |
| `application.state` | yes | `pending`, `prepared`, `consumed`, `ambiguous`, `superseded`, or `not-applicable`. |
| `application.dispatch-attempt-id` | before external dispatch | Stable pre-effect identity; feedback-only stage routing may omit it. |
| `application.consumed-at` | iff consumed | Time the workflow application was consumed. |
| `application.blockers[]` | yes, possibly empty | Durable prerequisite declarations and latest checks. |
| `application.execution-hold` / `feedback` | conditional | Approve-without-dispatch hold or rejection-to-rework route. |

`room-ref` is optional. When present, it is an opaque provider reference to one review
room/spec lineage. It is not attempt identity, a temporary package path, a UI session
id, or a new portable Review & Gate field. Attempts in the same lineage reuse the same
value. Spacedock may resume the provider through this reference but need not interpret
room storage, Probes, ProbeResults, logs, lenses, assessments, or deltas. Those provider
objects are not portable Review & Gate `Context` and add no fields to this encoding.

The whole-Briefing JCS digest, stage, attempt hierarchy, pointer selection, application
state, and routing are **Spacedock-only**. Review & Gate v1 separately requires each
Artifact `rev` to hash its unnormalized raw bytes; it does not define a whole-Briefing
digest, portable round/stage, mutable review status, lens, or routing executor.

## Identity, lifecycle, and selection invariants

1. Logical-gate, attempt, Briefing, Resolution, and application ids are unique within
   an entity in their respective identity classes.
2. A logical gate id names one entity/stage gate. Its `current-attempt` references one
   directly retained attempt. `gates.current`, when present, references that same
   attempt.
3. Attempt ids are Spacedock identities stable across open-attempt Briefing pointer
   updates. The first attempt omits `previous-attempt`; each later attempt names the
   formerly current, already-closed attempt. The chain cannot fork implicitly.
4. An open attempt has exactly one `current-briefing` and no `resolved-briefing`,
   Resolution, or application. Replacing its Briefing pointer/digest is one state commit
   and never mutates the portable Briefing it formerly referenced.
5. Adding/revising a lens or changing question, artifact revisions, design, evidence,
   or decision opportunity creates a new immutable Briefing in Subspace and advances
   the same open attempt's pointer. Subspace re-evaluates affected assessments and shows
   the delta. Spacedock records no `revise` unless a portable Resolution says so.
6. A binding Resolution closes an attempt only if `resolution.briefing` equals the open
   attempt's current Briefing id under compare-and-swap. The same commit freezes the
   exact id/digest/room reference as `resolved-briefing`, removes `current-briefing`,
   changes `open` to `closed`, preserves the Resolution, and creates one application.
7. A closed attempt has exactly one `resolved-briefing`, Resolution, and application.
   It never gains another Briefing/Resolution or reopens. Gate re-entry creates a new
   chained attempt. After rework, the room may carry a previously answered question to
   the new attempt and re-run it against the new Briefing; entity state records only the
   new attempt and current Briefing binding.
8. `approve` may omit portable rationale. `revise` and `hold` require a nonblank reason
   or an included earlier Annotation in the same Briefing log. `revise` maps to a
   Spacedock feedback application; `hold` maps to `action: none`, `not-applicable`.
9. Only `gates.current` and its closed attempt's frozen Briefing/application may be
   eligible. An `advance` application requires the exact binding `approve`. A
   `feedback` application requires the exact binding `revise`, a target stage, and
   feedback cycle/finding context. A binding `hold` remains `not-applicable`.
10. A post-close change to reviewed input cannot update the frozen Briefing. It marks
    the pending application stale/superseded and requires a new attempt.

## Per-Briefing logs and cross-Briefing provenance

Each immutable Briefing owns one separate logical ordered portable review log in
Subspace. A log may contain many Annotations and advisory Resolutions, but at most one
binding Resolution: the first Resolution attributed to the authorized approver identity
supplied externally by workflow tooling. The reviewer app accepts no later portable
entry in that log.

Advancing an open attempt to a new Briefing does not copy or merge the prior log. Review
& Gate `includes` may reference only earlier entries in the same log; forward, cyclic,
unknown, and cross-Briefing references are invalid. A Resolution for Briefing `3c`
cannot include an Annotation from `3b`.

If prior feedback must remain visible, the new Briefing explicitly carries a portable
Context `Reference` to stable revision-addressed material or receives a new Annotation
in its own log. Subspace's room may present the historical delta; the entity does not
duplicate log entries, lens definitions, assessment results, or delta objects.

A late portable binding Resolution for a no-longer-current Briefing remains binding in
that Briefing's provider log; Spacedock must not relabel it advisory. It is not adopted
as the attempt-closing Resolution because the pointer compare-and-swap fails. The open
attempt remains on its current Briefing and creates no application.

## Recording, closing, and consuming are separate operations

1. **Open attempt:** append a Spacedock attempt with one `current-briefing` reference;
   do not change lifecycle `status`.
2. **Advance presentation:** while open, replace the current Briefing reference/digest
   in one commit. Do not create a decision, application, or new attempt. Subspace owns
   revision history and presentation; state Git retains pointer history.
3. **Close with Resolution:** validate external authority, same-Briefing log rules, and
   pointer equality; atomically freeze `resolved-briefing`, store the exact Resolution,
   close the attempt, and create a `pending` (`approve`/`revise`) or `not-applicable`
   (`hold`) application. Do not advance `status` or dispatch.
4. **Observe:** refresh blockers or release an execution hold in later commits. Unknown
   and failed prerequisites fail closed.
5. **Prepare:** before an eligible application spawns a worker, commit a stable
   `dispatch-attempt-id` and `prepared`.
6. **Consume:** atomically commit a state-only feedback transition. For an external
   effect, commit the expected `status` transition, `consumed`, time, and durable receipt
   after the idempotent/queryable effect succeeds. An unresolved external outcome
   becomes `ambiguous` and is not retried under a new identity.

Every application must pass the common safeguards before consumption: matching
gate/attempt pointers, an exact current frozen Briefing digest, `pending` state, all
blockers satisfied, no active execution hold, and the expected lifecycle stage. Open,
non-current, stale, superseded, consumed, ambiguous, and `not-applicable` applications
fail closed.

Decision-specific eligibility then applies. `action: advance` requires the exact current
binding `approve`. `action: feedback` requires the exact current binding `revise`, its
`target-stage`, and durable feedback cycle/finding reference/digest. `decision: hold`
maps only to `action: none`, `state: not-applicable` and can never become eligible.

Consuming feedback atomically records the target-stage transition and durable route
context. It need not carry `dispatch-attempt-id` when that commit does not spawn a
worker; any later worker dispatch uses its own prepared identity. The consumed feedback
examples above are therefore complete: they retain `target-stage`, cycle, finding
reference/digest, and `consumed-at` beside the exact binding `revise`.

## Concurrency invariants

- Every write compares the expected entity/state revision under the state-checkout
  mutation lock. A stale writer re-reads and re-evaluates.
- Two Briefing pointer updates from the same base conflict. A pointer update and an
  attempt close compete on `current-briefing`; no timestamp, list position, or Git
  parent implicitly wins.
- A Resolution for the losing Briefing may remain portable provider truth but cannot
  create a Spacedock application.
- Frozen Briefing references and Resolution nodes never field-merge. Merge resolution
  selects a complete pointer/attempt/application state and emits an explicit resolution
  event; until then, eligibility is false.
- Sequence order has no selection semantics. A canonical writer orders gates by id and
  attempts by sequence; explicit pointers determine current state.

## Ownership and Spacedock-only conn policy

Review & Gate owns immutable portable object shapes and one-Briefing log invariants.
Workflow tooling owns workflow position, externally authorized-approver identity, and
routing interpretation/execution. The Subspace reviewer app/room owns full Briefings,
resource verification, attribution, per-Briefing logs, room-local probes, exact-
Briefing results, citations, responder identity, lenses, assessments, deltas, concrete
persistence, reconciliation, interaction mode, and UI closure.

The entity owns only the lean durable workflow index: logical gates, distinct
Spacedock attempts, current or frozen Briefing references/digests, exact adopted
Resolution, blockers/holds/routes, and application state. Prompts, transcripts,
temporary paths, pane/session ids, credentials, and private observations do not enter
it.

An ensign may publish and present a gate and transport Subspace's authenticated
Resolution through a recorder that checks the exact current Briefing and externally
authorized identity. The recorder cannot grant captain authority or change workflow
state. Subspace authenticates and stamps the acting approver but does not choose the
authorized approver. The First Officer validates the committed Resolution and alone
consumes its application to transition, route, or dispatch.

Spacedock retains a stricter authoring rule for a First Officer exercising delegated
conn authority: its auto-approval must contain a nonblank `reason`. This is
**Spacedock-only**. Base Review & Gate permits reasonless `approve`; generic portable
parsing and entity schema remain permissive, and a captain's ordinary reasonless
approval is valid. Enforcement belongs to the FO authoring path; the resulting exact
Resolution is preserved without a new portable field.

## Commit-derived events and provider room integration

The projector diffs old/new `gates` trees and emits logical-attempt opened,
current-Briefing changed, attempt closed, Resolution adopted, application,
stale-decision, supersession, and feedback-route events. The current entity directly
enumerates every gate and distinct adjudication attempt plus each attempt's current or
resolved Briefing binding. Git history reconstructs prior open-attempt pointer/digest
values and intermediate application observations. A cache or temporary review package
is never authority.

The optional opaque room reference lets the provider own full presentation history for
one spec lineage. Spacedock preserves and reuses the value but does not interpret room
shape, default view, storage, navigation, Probe records, or conflict rendering. A room
or lens may not replace gate, attempt, Briefing, Resolution, or application identity in
the entity.
