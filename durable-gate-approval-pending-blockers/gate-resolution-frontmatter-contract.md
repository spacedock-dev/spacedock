# Gate Resolution frontmatter contract

Status: proposed v1 contract for `3k` ideation review
Date: 2026-07-18

## Recovered requirement and design lineage

The original `3k` filing asked how “superseding approvals and multiple review rounds
retain audit history while exposing one current decision.” A one-slot `gate` snapshot
does not answer that question: it requires Git replay to discover earlier logical gates
and attempts, so the entity does not directly contain its durable gate record.

Two earlier decisions still hold:

1. Closed Spacedock PR #474 (`iamcxa/status-apply-gate`, commits
   `685fe7bcda4a51b8e2c06da52e80c079f62ac8e0` through
   `5dee22831856db65db2acfefee0849c5f990f5d1`) put the externally captured decision in
   binary-owned entity frontmatter. Its apply operation also changed `status`; `3k`
   must split those operations.
2. Review & Gate v1, at `spacedock-subspace` commit
   `bd17bdb23318f815d17a1d10ea2a6d39ab449520` and
   `docs/review-and-gate.md` blob
   `14f3eb91ec85bfcc08bb3330c21b94cc77f4529f`, supplies immutable `Briefing` and
   `Resolution` identities. Portable decisions are `approve`, `revise`, and `hold`;
   workflow rejection for rework is `revise` plus a Spacedock feedback application.

The authoritative source for this proposal is the complete
`../spacedock-subspace/docs/review-and-gate.md` contract at that commit, especially
§1 “Scope and ownership,” §2 “Model,” §3 “Decisions,” §5 “Review entries,” §7
“Review log,” and §8 “Versioning and serialization.” The `gates` tree below is a
Spacedock workflow index over those portable objects; it is not a change to Review &
Gate v1.

The smallest representation that satisfies the original requirement is one plural
top-level `gates` mapping. It contains a collection of logical gates, every binding
resolved attempt for each gate, and explicit current-selection pointers. Git commits
remain the audit trail of mutations and the source for projected events, but replay is
not required merely to enumerate prior gates or attempts.

## Review & Gate v1 boundary and audit corrections

The source audit fixes these boundaries:

- One immutable portable `Briefing` is one decision opportunity. A different question,
  artifact revision set, or decision opportunity receives a new Briefing id. A
  Spacedock “attempt” is an index entry keyed by that exact id; it is not a second
  portable object type.
- One logical ordered portable review log belongs to one Briefing. It may contain many
  advisory Resolutions from other actors, but at most one binding Resolution: the first
  Resolution attributed to the externally authorized approver. The entity copies only
  that binding Resolution; it does not collapse advisory entries into it.
- Workflow tooling owns the authorized-approver identity and routing interpretation.
  The Subspace reviewer app receives that authority externally, stamps attribution,
  maintains entry order, and identifies the first matching Resolution. The prior draft
  imprecisely assigned “binding-authority determination” to Subspace alone.
- `approve` needs no portable rationale. `revise` and `hold` require a nonblank `reason`
  or an `includes` reference to at least one earlier Annotation in the same Briefing.
  An included advisory Resolution alone is not the required rationale witness.
- Round, workflow stage, supersession, current selection, blockers, mutable application
  status, and routing execution are not portable Review & Gate fields. Every such field
  below is Spacedock-owned metadata outside the copied Resolution.
- Tagged JSON defines Review & Gate's portable semantics. The YAML `resolution` node
  below is a semantic transcription of the exact JSON object for entity storage; YAML
  syntax, ordering, and scalar presentation do not create a different portable format.

## Physical representation

```yaml
gates:
  version: 1
  current:
    gate: gate:docs-dev:3k:validation
    attempt: briefing:3k-validation-r2
  records:
    - id: gate:docs-dev:3k:ideation
      stage: ideation
      current-attempt: briefing:3k-ideation-r2
      attempts:
        - id: briefing:3k-ideation-r1
          round: 1
          reviewed-digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
          resolution:
            type: Resolution
            id: resolution:captain-3k-ideation-r1
            briefing: briefing:3k-ideation-r1
            by: person:captain
            at: 2026-07-16T09:00:00Z
            decision: revise
            reason: Clarify the dispatch blocker contract.
            includes: []
        - id: briefing:3k-ideation-r2
          round: 2
          reviewed-digest: sha256:2222222222222222222222222222222222222222222222222222222222222222
          supersedes: briefing:3k-ideation-r1
          resolution:
            type: Resolution
            id: resolution:captain-3k-ideation-r2
            briefing: briefing:3k-ideation-r2
            by: person:captain
            at: 2026-07-17T09:00:00Z
            decision: approve
      applications:
        - attempt: briefing:3k-ideation-r1
          id: application:3k-ideation-r1
          action: feedback
          target-stage: backlog
          state: consumed
          consumed-at: 2026-07-16T09:02:00Z
          blockers: []
          feedback:
            cycle: 1
            finding-ref: resolution:captain-3k-ideation-r1
            finding-digest: sha256:2111111111111111111111111111111111111111111111111111111111111111
        - attempt: briefing:3k-ideation-r2
          id: application:3k-ideation-r2
          action: advance
          target-stage: implementation
          state: consumed
          dispatch-attempt-id: dispatch:3k-implementation-r1
          consumed-at: 2026-07-17T09:03:00Z
          blockers: []
    - id: gate:docs-dev:3k:validation
      stage: validation
      current-attempt: briefing:3k-validation-r2
      attempts:
        - id: briefing:3k-validation-r1
          round: 1
          reviewed-digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
          resolution:
            type: Resolution
            id: resolution:captain-3k-validation-r1
            briefing: briefing:3k-validation-r1
            by: person:captain
            at: 2026-07-18T08:00:00Z
            decision: revise
            reason: The production coordinator is missing.
            includes: []
        - id: briefing:3k-validation-r2
          round: 2
          reviewed-digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
          supersedes: briefing:3k-validation-r1
          resolution:
            type: Resolution
            id: resolution:captain-3k-validation-r2
            briefing: briefing:3k-validation-r2
            by: person:captain
            at: 2026-07-18T10:30:00Z
            decision: approve
      applications:
        - attempt: briefing:3k-validation-r1
          id: application:3k-validation-r1
          action: feedback
          target-stage: implementation
          state: consumed
          consumed-at: 2026-07-18T08:02:00Z
          blockers: []
          feedback:
            cycle: 1
            finding-ref: resolution:captain-3k-validation-r1
            finding-digest: sha256:4333333333333333333333333333333333333333333333333333333333333333
        - attempt: briefing:3k-validation-r2
          id: application:3k-validation-r2
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
            id: hold:captain:3k-validation-r2
            state: active
            by: person:captain
            at: 2026-07-18T10:30:00Z
            reason: Approval is durable; do not consume it until I release this hold.
```

This example directly contains two logical gates and two immutable resolved attempts
for each gate. `gates.current` selects validation round 2 for workflow evaluation;
`current-attempt` retains the selected attempt within each logical gate, including the
already-consumed ideation gate. Historical attempts and their latest application state
are available from a single entity read.

Existing scalar frontmatter remains unchanged. The binary must read and mutate this
mapping through `yaml.Node`; the legacy scalar-only `ParseFrontmatter` view may expose
`gates` as an empty scalar to old callers.

## Field contract

The layer boundary is structural: `gates`, logical-gate identity, stage, round, digest,
supersession, selection, and every `applications[]` field are Spacedock-owned.
`attempts[].id` copies the portable `Briefing.id`; `attempts[].resolution` is the only
copied portable object.

| Field | Required | Meaning |
|---|---:|---|
| `gates.version` | yes | Integer `1`; unsupported versions fail closed. |
| `gates.current` | when a resolved attempt is selected at the entity's current gated stage | Pair of `gate` and `attempt` ids. It is scheduler selection, not history ordering. |
| `records[].id` | yes | Stable logical-gate identity for this entity and gated stage. Re-review rounds retain it. |
| `records[].stage` | yes | Exact workflow stage reviewed by this logical gate. It does not change after the entity advances. |
| `records[].current-attempt` | yes | Attempt selected for this logical gate. It must reference an entry in the same record. |
| `attempts[].id` | yes | Exact immutable Review & Gate `Briefing.id`; unique across the entity. |
| `attempts[].round` | yes | Positive workflow review round. It is an external label, not portable Review & Gate vocabulary. |
| `attempts[].reviewed-digest` | yes | `sha256:` digest of the exact canonical Briefing. |
| `attempts[].supersedes` | after the first attempt | Immediately preceding attempt id in the same logical gate. |
| `attempts[].resolution` | yes | Exact binding portable `Resolution`, semantically preserved without normalization or field loss. |
| `applications[].attempt` | yes | Attempt id in the same logical-gate record; exactly one application exists per attempt. |
| `applications[].id` | yes | Stable one-use workflow application identity, unique across the entity. |
| `applications[].action` | yes | Spacedock interpretation: `advance`, `feedback`, or `none`; Review & Gate does not execute it. |
| `applications[].target-stage` | for `advance` or `feedback` | Exact target stage for consumption. |
| `applications[].state` | yes | `pending`, `prepared`, `consumed`, `ambiguous`, `superseded`, or `not-applicable`. |
| `applications[].dispatch-attempt-id` | for `prepared` or `ambiguous`; optional on `consumed` | Stable pre-effect identity at an idempotent or queryable spawn boundary. |
| `applications[].consumed-at` | for `consumed` | RFC 3339 time of the atomic workflow transition and durable effect receipt. |
| `applications[].blockers[]` | yes, possibly empty | Declared dispatch prerequisites and their latest durable checks. |
| `applications[].execution-hold` | no | Workflow-owned pause after approval, distinct from a Review & Gate `hold` decision. |
| `applications[].feedback` | for `feedback` | Cycle and finding identity for durable rejection-to-rework routing. |

The `resolution` node is not a hand-picked summary. It preserves the submitted
portable object, including `type`, `id`, `briefing`, `by`, `at`, `decision`, `reason`
when present, `includes`, and any valid additive version-1 fields. Its `briefing` must
equal the containing attempt id. Review & Gate attribution and binding authority are
validated before recording; the `by` string does not self-assert authority.

The portable review log may contain annotations and advisory Resolutions before the
binding entry. Those entries remain in their one-Briefing log. The recorder accepts
only the first Resolution attributed to the externally supplied authorized approver as
`attempts[].resolution`, rejects a second claimed binding Resolution for that Briefing,
and never derives binding from an entry's contents. An `approve` object may omit both
`reason` and `includes`, as the two approving examples do. A `revise` or `hold` without
a nonblank reason is accepted only when `includes` names at least one earlier Annotation
in that same ordered log.

`reviewed-digest` is SHA-256 over RFC 8785 JSON Canonicalization Scheme bytes of the
exact Review & Gate `Briefing`. The Briefing carries its immutable question, artifact
ids and revisions, routing context, criteria, and evidence. The Subspace reviewer app
resolves and verifies each Artifact revision before opening, as Review & Gate §4
requires. Separately, the Spacedock recorder verifies the exact Briefing identity and
digest it is binding, and the consumer compares that digest with the reconstructed
workflow decision opportunity. Changed gate-defining input requires a new Briefing id,
digest, and attempt.

That whole-Briefing JCS digest is a **Spacedock-only binding constraint**, not a Review &
Gate v1 field or canonicalization rule. Review & Gate separately requires each Artifact
`rev` to hash its unnormalized raw bytes and requires the reviewer app to resolve and
verify artifacts before opening. Spacedock stores the additional Briefing digest so it
can compare the workflow's reconstructed decision opportunity before consumption; it
must not present JCS as portable Review & Gate behavior.

## Identity, selection, and supersession invariants

1. Gate ids, attempt ids, Resolution ids, and application ids are each unique within
   an entity. `resolution.briefing` equals its attempt id.
2. A logical gate id names one entity/stage gate across review rounds. A different
   gated stage gets a different gate id. Renaming or redefining a stage requires an
   explicit migration; it never silently reuses an id.
3. Every admitted `attempts[]` element is wholly immutable: id, round, digest,
   supersedes link, and exact Resolution. Any later change under the same id is a
   conflict. Its separately keyed `applications[]` element is mutable only through the
   allowed workflow state transitions below.
4. The first attempt omits `supersedes`. Every later attempt names the gate's former
   `current-attempt`; insertion and pointer update occur in one commit. The links form
   one chain. A fork has no implicit winner and fails closed.
5. `records[].current-attempt` points to exactly one attempt in that record.
   `gates.current`, when present, points to an existing record and that record's
   `current-attempt`. Only this pair may be eligible for consumption at the current
   gated stage.
6. Each attempt has exactly one application in the same gate record, and each
   application names exactly one attempt. Applications cannot move between attempts.
7. `gates.current` is cleared when its application consumes a transition away from the
   gated stage. A later gate Resolution selects its own pair. Historical gate records
   and their per-gate current pointers remain in place.
8. A non-current or superseded attempt is never dispatchable, even if its Resolution
   says `approve` and its blockers later appear satisfied. A consumed application is
   one-use and cannot become pending again.
9. A captain-facing rejection for rework is stored as portable `decision: revise`
   with `application.action: feedback`. A portable `decision: hold` has action `none`
   and application state `not-applicable`; an approved execution hold remains
   `decision: approve`, action `advance`, state `pending`, with an active
   `execution-hold`.

`records[].stage`, `attempts[].round`, `attempts[].reviewed-digest`, and
`attempts[].supersedes` are all Spacedock index fields. They do not enter the copied
portable Resolution and do not imply that Review & Gate has portable stage, round,
history, or mutable review status. A changed question, artifact revision set, or
decision opportunity always creates a new Briefing/attempt id even if Spacedock's
logical gate and round label stay otherwise related.

## Application transitions and separate operations

Recording a binding Resolution and consuming its workflow consequence are separate
state commits.

1. **Record:** under an expected entity/state revision, append the complete immutable
   attempt, create its application, and update both selection pointers. `approve` and
   `revise` begin `pending`; `hold` begins `not-applicable`. This commit does not change
   `status`, prepare a dispatch, or spawn a worker.
2. **Observe:** blocker refreshes and execution-hold release update only the selected
   attempt's application in later commits. `unknown` and `failed` blockers fail closed.
3. **Prepare:** if the selected application needs an external effect and is eligible,
   commit a stable `dispatch-attempt-id` and move `pending` to `prepared` before the
   effect.
4. **Consume:** after the idempotent or queryable effect succeeds, atomically commit the
   expected `status` transition, `consumed`, `consumed-at`, and its durable receipt
   under the same application/dispatch identity. An unresolved effect becomes
   `ambiguous`; it is not retried under a new identity.
5. **Supersede:** recording a new attempt atomically moves the prior unconsumed
   application to `superseded`, appends the new attempt, and advances the per-gate and
   top-level pointers. The prior attempt and latest application state remain present.

Eligibility is the conjunction of matching current pointers, current reconstructed
Briefing digest, binding `approve`, `pending` application, all blockers satisfied, no
active execution hold, and the expected current stage. No single field substitutes for
that check.

## Concurrency invariants

- Every write is compare-and-swap against the expected entity/state revision and is
  committed under the state-checkout mutation lock. A stale writer re-reads and
  re-evaluates; it does not overwrite the collection.
- Concurrent writers may be reconciled as a set union only when they append disjoint
  logical-gate records or disjoint non-current historical data and leave every existing
  immutable node and selection/application field byte-semantically unchanged.
- Two attempts that both supersede the same current attempt are a fork. Two writers
  that change `gates.current`, one `current-attempt`, blocker observations, execution
  hold, or application state from the same base conflict. No timestamp, list order, or
  Git parent implicitly wins.
- Sequence order has no selection semantics. A canonical writer orders gates by id,
  attempts by their supersession chain, and applications beside the corresponding
  attempt; pointers and ids, never YAML position, determine meaning.
- Merge resolution selects a complete attempt/application state and records an explicit
  merge-resolution event. Field-wise merging of Resolution or application nodes is
  forbidden. Until resolution, status reports conflict and dispatch eligibility is
  false.

## Portable contract, Subspace app, and entity ownership

The entity deliberately does not become a second Review & Gate database.

The Review & Gate contract owns the immutable Briefing, portable object shapes, the
logical one-Briefing review log, and their invariants. Workflow tooling—including
Spacedock—owns workflow position, the externally authorized-approver identity, and
routing interpretation/execution.

The Subspace reviewer app owns resource resolution and verification, cached display
bytes, authentication and attribution stamping, drafts/edits, selector placement,
reconciliation, interaction mode, concrete review-log persistence, and UI closure. It
matches entries to the authorized identity supplied by workflow tooling; it does not
mint its own workflow authority. Whether a future persistent room or lens is the
durable home and navigation surface for app-owned objects is not settled here.

Spacedock's entity record owns only the workflow binding needed after the review
invocation ends: logical gate/stage identity, Briefing id and canonical digest, the
exact binding Resolution, selection and supersession links, blockers/hold, feedback
route, and one-use application state. Temporary package paths, prompts, transcripts,
pane/session ids, credentials, and private runtime observations never enter it.

If a binding Resolution relies on `includes` rather than a typed reason, the ids remain
exact in the entity while the referenced entries remain in the provider-owned review
log. Loss of that log makes rationale evidence unavailable; it does not authorize the
entity to synthesize annotation bodies.

## Spacedock-only policy constraints

The multi-gate index, whole-Briefing digest, application state machine, and
approve-but-do-not-dispatch execution hold are all Spacedock constraints layered around
portable Review & Gate v1.

Spacedock retains a stricter authoring rule: a First Officer exercising delegated conn
authority must give an explicit nonblank `reason` for an auto-approval. This is also
**Spacedock-only**. Base Review & Gate v1 permits `approve` with no rationale. The
conn-made Resolution uses the already-valid optional `reason` member; generic Review &
Gate validation must still accept an authorized reasonless `approve`, and a captain's
ordinary reasonless approval remains valid. Enforcement belongs to the First Officer's
Spacedock authoring path when it exercises the conn, not to the generic Review & Gate
parser or entity schema. The entity preserves the resulting reason; no new portable
field is invented.

## Commit-derived events and projections

The projector diffs complete old/new `gates` trees and emits attempt, Resolution,
selection, application, supersession, and feedback-route events. The current entity
tree is the directly readable durable collection; Git history adds who/when/order,
intermediate application observations, and explicit merge resolution. A projection
cache or temporary review package is never authority.

Cold reconstruction from the current entity alone must enumerate all logical gates,
all admitted attempts, their exact binding Resolutions, and latest application states.
Replaying Git additionally reconstructs the transition/event history. These are
different guarantees; the one-slot model incorrectly provided only the second.

## Lens and persistent-room integration questions

This contract intentionally does not define unresolved lens semantics. Integration
must later answer:

- whether a lens defaults to the selected attempt, one logical gate's chain, or the
  full entity gate collection;
- whether a persistent room stores the canonical Briefing/log or only navigates to a
  Subspace-owned store, and what availability guarantee an entity `includes` reference
  receives;
- whether room/lens identities become optional evidence references in an attempt
  without becoming gate, attempt, or application identity;
- how lens-visible concurrent branches present a conflict before an explicit workflow
  selection is committed.

None of those choices changes the entity's minimum durable workflow binding or permits
a lens projection to replace it.
