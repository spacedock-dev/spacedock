# Gate Resolution frontmatter contract

Status: proposed v1 contract for `3k` ideation review
Date: 2026-07-18

## Authoritative sources and corrected hierarchy

The portable authority is the complete
`../spacedock-subspace/docs/review-and-gate.md` contract at
`spacedock-subspace` commit `bd17bdb23318f815d17a1d10ea2a6d39ab449520`, blob
`14f3eb91ec85bfcc08bb3330c21b94cc77f4529f`, especially §1 “Scope and
ownership,” §2 “Model,” §3 “Decisions,” §5 “Review entries,” §7 “Review log,”
and §8 “Versioning and serialization.” The entity tree below indexes those portable
objects; it does not change Review & Gate v1.

Closed Spacedock PR #474 (`iamcxa/status-apply-gate`, commits
`685fe7bcda4a51b8e2c06da52e80c079f62ac8e0` through
`5dee22831856db65db2acfefee0849c5f990f5d1`) supplied the retained physical direction:
the workflow binding belongs in binary-owned entity frontmatter. `3k` removes PR #474's
coupling between recording a decision and changing `status`.

The corrected identity hierarchy is:

```text
logical Spacedock gate
  -> stable Spacedock adjudication attempts
       -> immutable portable Briefing snapshots
            -> one separate portable ordered log per Briefing
       -> exact adopted binding Resolution for the selected Briefing (closes attempt)
       -> mutable Spacedock workflow application (recorded, then consumed separately)
```

A gate attempt is **not** a Review & Gate `Briefing`. It is a stable Spacedock
adjudication session. While open it may advance its `current-briefing` pointer through
multiple immutable Briefings as the presentation gains a lens or as design/evidence is
revised. That does not itself imply portable `revise`, close the attempt, or create a
new workflow attempt. A binding Resolution closes the attempt only when it references
the attempt's exact current Briefing under the same compare-and-swap write.

Re-entry to a logical gate after a closed binding Resolution is the normal new-attempt
boundary. Closed attempts never reopen. If approved content changes while its workflow
application is still pending, the closed approval becomes stale and a new attempt is
required; a new Briefing cannot be appended to the closed attempt.

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
          current-briefing: briefing:3k-ideation-1a
          briefings:
            - id: briefing:3k-ideation-1a
              digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
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
          current-briefing: briefing:3k-ideation-2b
          briefings:
            - id: briefing:3k-ideation-2a
              digest: sha256:2222222222222222222222222222222222222222222222222222222222222222
            - id: briefing:3k-ideation-2b
              digest: sha256:3222222222222222222222222222222222222222222222222222222222222222
              previous-briefing: briefing:3k-ideation-2a
              change:
                cause: lens-added
                delta-ref: artifact:3k-ideation-2b-delta
                delta-digest: sha256:4222222222222222222222222222222222222222222222222222222222222222
                affected-assessments:
                  - id: assessment:dispatch-safety
                    state: re-evaluated
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
          current-briefing: briefing:3k-validation-1a
          briefings:
            - id: briefing:3k-validation-1a
              digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
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
          current-briefing: briefing:3k-validation-2b
          briefings:
            - id: briefing:3k-validation-2a
              digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
            - id: briefing:3k-validation-2b
              digest: sha256:5444444444444444444444444444444444444444444444444444444444444444
              previous-briefing: briefing:3k-validation-2a
              change:
                cause: evidence-revised
                delta-ref: artifact:3k-validation-2b-delta
                delta-digest: sha256:6444444444444444444444444444444444444444444444444444444444444444
                affected-assessments:
                  - id: assessment:production-coordinator
                    state: re-evaluated
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

The example directly contains two logical gates and two adjudication attempts for each.
Ideation attempt 2 and validation attempt 2 each advanced through two immutable
Briefings while open; their binding Resolutions reference `2b`, the selected snapshot
that closed the attempt. The earlier `2a` identities and change/delta records remain
directly readable. Approval recording and application consumption are still visibly
separate: validation attempt 2 is closed and approved while its application is pending.

An open attempt uses the same shape with `state: open`, one or more `briefings`, and a
valid `current-briefing`; it has no `resolution` and no `application`.

## Field and layer contract

The top-level `gates` tree and all fields except the copied Resolution are
Spacedock-owned. Briefing ids originate in portable Review & Gate; full immutable
Briefing objects and their logs remain provider-owned. The YAML `resolution` node is a
semantic transcription of the exact tagged-JSON portable object, including valid
additive v1 fields; YAML syntax does not define another portable format.

| Field | Required | Meaning |
|---|---:|---|
| `gates.version` | yes | Integer `1`; unsupported versions fail closed. |
| `gates.current` | when an attempt is selected at the current gated stage | Pair of logical gate and Spacedock attempt ids. |
| `records[].id` / `stage` | yes | Stable logical entity/stage gate and its Spacedock workflow stage. |
| `records[].current-attempt` | yes | Selected adjudication attempt for this logical gate. |
| `attempts[].id` | yes | Stable Spacedock adjudication-session identity; not a `Briefing.id`. |
| `attempts[].sequence` / `previous-attempt` | yes / after first | Spacedock ordering and chain within the logical gate. |
| `attempts[].state` | yes | `open` or `closed`; closed attempts never reopen. |
| `attempts[].current-briefing` | yes | Selected immutable Briefing snapshot for this attempt. |
| `attempts[].briefings[].id` | yes | Exact portable `Briefing.id`, unique across the entity. |
| `attempts[].briefings[].digest` | yes | Spacedock JCS digest of the exact portable Briefing. |
| `previous-briefing` | after first snapshot | Prior Briefing in this attempt; this is Spacedock provenance, not a portable field. |
| `change` | after first snapshot | Spacedock cause, delta reference/digest, and affected-assessment re-evaluation summary. |
| `attempts[].resolution` | iff closed | Exact adopted binding portable Resolution for `current-briefing`. |
| `attempts[].application` | iff closed | One-use Spacedock application created when the attempt closes. |
| `application.action` / `target-stage` | yes / for advance or feedback | `advance`, `feedback`, or `none`, and its Spacedock target. |
| `application.state` | yes | `pending`, `prepared`, `consumed`, `ambiguous`, `superseded`, or `not-applicable`. |
| `application.dispatch-attempt-id` / `consumed-at` | by state | Pre-effect identity and consumption time. |
| `application.blockers[]` | yes, possibly empty | Durable prerequisite declarations and latest checks. |
| `application.execution-hold` / `feedback` | conditional | Approve-without-dispatch hold or rejection-to-rework route. |

The whole-Briefing JCS digest, round/sequence, stage, attempt hierarchy, snapshot
selection, change/delta metadata, application state, and routing are **Spacedock-only**.
Review & Gate v1 separately requires each Artifact `rev` to hash its unnormalized raw
bytes; it does not define a whole-Briefing digest, portable round/stage, mutable review
status, supersession, lens, or routing executor.

The `change` node records why Spacedock selected a new immutable decision snapshot and
where the presentable delta evidence lives. It does not define lens internals. A changed
question, artifact revision set, or decision opportunity receives a new Briefing id as
Review & Gate §2 requires. Adding/revising a lens or revising design/evidence therefore
creates a new immutable Briefing, re-evaluates affected assessments, and advances the
pointer within the still-open attempt. It does not synthesize `revise`.

## Identity, lifecycle, and selection invariants

1. Logical-gate, attempt, Briefing, Resolution, and application ids are unique within
   an entity in their respective identity classes.
2. A logical gate id names one entity/stage gate. Its `current-attempt` references one
   attempt in that record. `gates.current`, when present, references that same attempt.
3. Attempt ids are Spacedock identities stable across Briefing changes. The first
   attempt omits `previous-attempt`; each later attempt names the formerly current,
   already-closed attempt. The chain cannot fork implicitly.
4. Every `briefings[]` entry is immutable. Its id, digest, previous link, and change
   metadata never change after insertion. Each later snapshot names the formerly
   current Briefing, and insertion plus pointer advance is one commit.
5. While an attempt is open, presentation augmentation or revised design/evidence may
   append a Briefing and advance `current-briefing`. It records no Resolution or
   application and does not change `status` or dispatch.
6. A binding Resolution closes an attempt only if `resolution.briefing` equals the
   attempt's `current-briefing` at compare-and-swap time. The close commit preserves the
   exact Resolution, changes `open` to `closed`, and creates exactly one application;
   it still does not change `status` or dispatch.
7. `approve` may omit portable rationale. `revise` and `hold` require a nonblank reason
   or an included earlier Annotation in the same Briefing log. Spacedock presents a
   workflow rejection as `revise` plus `application.action: feedback`.
8. `hold` closes the attempt with `action: none` and `state: not-applicable`.
   Reconsideration is a new attempt. `approve` with an active Spacedock execution hold
   remains approved with a pending, ineligible application.
9. A closed attempt never gains another Briefing or Resolution. Gate re-entry after its
   binding result creates a new attempt. Changed reviewed input while an approval waits
   makes its application stale/superseded and also requires a new attempt.
10. Only `gates.current` plus its attempt's `current-briefing` and application may be
    eligible. Non-current, open, stale, superseded, consumed, ambiguous, or
    not-applicable records never dispatch.

## Separate portable logs and cross-Briefing provenance

Each immutable Briefing owns one separate logical ordered portable review log. A log may
contain many Annotations and advisory Resolutions, but at most one binding Resolution:
the first Resolution attributed to the authorized approver identity supplied externally
by workflow tooling. The reviewer app accepts no later portable entry in that log.

Advancing an open Spacedock attempt to a new Briefing does not copy or merge the prior
log. Prior Annotations and advisory Resolutions remain attached to their prior Briefing.
Review & Gate `includes` may reference only earlier entries in the same log; forward,
cyclic, unknown, and cross-Briefing references are invalid. In particular, a Resolution
for Briefing `2b` cannot include an Annotation from `2a`.

If prior feedback must remain visible, the new Briefing explicitly carries a portable
Context `Reference` to a stable, revision-addressed artifact or the reviewer submits a
new Annotation in the new log. Spacedock may show the stored delta reference beside the
new Briefing. Neither path silently carries entry identity or attribution across logs.

A late portable binding Resolution for a no-longer-current Briefing remains binding in
that Briefing's provider log; Spacedock must not relabel it advisory. It is not adopted
as the attempt-closing Resolution because the current-pointer compare-and-swap fails.
The attempt remains open on its current Briefing and status reports the stale decision
observation; no application is created.

## Recording, closing, and consuming are separate operations

1. **Open attempt:** create a Spacedock attempt and its first immutable Briefing pointer
   without changing lifecycle `status`.
2. **Advance presentation:** while open, append a new immutable Briefing plus delta and
   affected-assessment evidence, then atomically advance `current-briefing`. Do not
   create a decision or application.
3. **Close with Resolution:** validate external authority, same-Briefing log rules, and
   current pointer; then atomically store the exact binding Resolution, close the
   attempt, and create a `pending` (`approve`/`revise`) or `not-applicable` (`hold`)
   application. Do not advance `status` or dispatch.
4. **Observe:** refresh blockers or release an execution hold in later commits. Unknown
   and failed prerequisites fail closed.
5. **Prepare:** for an eligible action with an external effect, commit a stable
   `dispatch-attempt-id` and `prepared` before the effect.
6. **Consume:** after the idempotent/queryable effect succeeds, atomically commit the
   expected `status` transition, `consumed`, time, and durable receipt. An unresolved
   outcome becomes `ambiguous` and is not retried under a new identity.

Eligibility requires matching gate/attempt/current-Briefing pointers, a closed binding
`approve`, current reconstructed Briefing digest, pending application, all blockers
satisfied, no active execution hold, and expected lifecycle stage.

## Concurrency invariants

- Every write compares the expected entity/state revision under the state-checkout
  mutation lock. A stale writer re-reads and re-evaluates.
- Two new Briefings that both name the same current Briefing form a snapshot fork; two
  closes against the same open attempt conflict. No timestamp, list position, or Git
  parent implicitly wins.
- Snapshot advancement and attempt closing compete on `current-briefing`. A Resolution
  for the losing snapshot may remain portable provider truth but cannot create a
  Spacedock application.
- Existing Briefing and Resolution nodes never field-merge. Merge resolution selects a
  complete pointer/attempt/application state and emits an explicit resolution event;
  until then, eligibility is false.
- Sequence order has no selection semantics. A canonical writer orders gates by id,
  attempts by sequence, and Briefings by their previous chain; pointers determine
  current state.

## Ownership and Spacedock-only conn policy

Review & Gate owns immutable portable object shapes and one-Briefing log invariants.
Workflow tooling owns workflow position, externally authorized-approver identity, and
routing interpretation/execution. The Subspace reviewer app owns resource resolution
and verification, attribution stamping, concrete log persistence, drafts/edits,
reconciliation, interaction mode, and UI closure.

The entity owns only the durable workflow index and adopted binding: logical gates,
Spacedock attempts, Briefing ids/digests and selection, change/delta evidence, exact
adopted Resolution, blockers/holds/routes, and application state. Full Briefings,
per-Briefing logs, annotations, advisory Resolutions, prompts, transcripts, temporary
paths, pane/session ids, credentials, and private observations do not enter it.

Spacedock retains a stricter authoring rule for a First Officer exercising delegated
conn authority: its auto-approval must contain a nonblank `reason`. This is
**Spacedock-only**. Base Review & Gate permits reasonless `approve`; generic portable
parsing and entity schema remain permissive, and a captain's ordinary reasonless
approval is valid. Enforcement belongs to the FO authoring path; the exact resulting
Resolution is then preserved without a new portable field.

## Commit-derived events and unresolved lens integration

The projector diffs complete old/new `gates` trees and emits logical-attempt opened,
Briefing selected, attempt closed, Resolution adopted, application, stale-decision,
supersession, and feedback-route events. The current entity directly enumerates every
gate, attempt, Briefing identity, adopted Resolution, and latest application state; Git
history adds transition order and intermediate observations. A cache or temporary
review package is never authority.

This contract records that a lens change can produce a new Briefing and delta, but does
not define unresolved lens/persistent-room semantics: lens object shape, default view,
room storage, navigation, and conflict rendering remain open. A future lens may point
to provider objects or enrich portable Briefing Context; it may not replace gate,
attempt, Briefing, Resolution, or application identity in the entity.
