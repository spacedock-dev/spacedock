# Gate Resolution frontmatter contract

Status: proposed v1 contract for `3k` ideation review  
Date: 2026-07-18

## Recovered design lineage

This contract completes an existing design; it does not replace it with an event
ledger.

1. Closed Spacedock PR #474 (`iamcxa/status-apply-gate`, commits
   `685fe7bcda4a51b8e2c06da52e80c079f62ac8e0` through
   `5dee22831856db65db2acfefee0849c5f990f5d1`) established the original storage
   decision: an externally captured gate decision becomes entity frontmatter
   (`gate-id` and `gate-verdict`) through a binary-owned writer. Its apply operation
   also changed `status`, which is the coupling this task must remove.
2. PR #474 was closed, not merged, because Review & Gate v1 superseded its portable
   vocabulary: `Briefing` and `Resolution` replace the older gate-verdict packet, and
   `reject` is not a portable review decision. The current Review & Gate source is
   `spacedock-subspace` commit `bd17bdb23318f815d17a1d10ea2a6d39ab449520`,
   `docs/review-and-gate.md` blob `14f3eb91ec85bfcc08bb3330c21b94cc77f4529f`.
3. The Draft Ledger binding in Spacedock commit
   `61b9a66107ff3de155e4319c8d7681a6af9ba720` independently preserves the same
   ownership boundary: external systems own gate and Resolution identity; Spacedock
   owns entity state and application coordination; later application facts do not
   get folded into the provider binding.

The retained decision is therefore: **the entity's committed YAML frontmatter is the
authoritative current Spacedock gate record.** Git history is its durable history.
Projectors derive events from those commits; projected events never substitute for
the entity representation.

## Physical representation

An entity may carry one top-level `gate` mapping for its current or most recently
consumed gate attempt. Existing scalar frontmatter remains unchanged. The binary must
read and mutate this mapping through `yaml.Node`; the legacy scalar-only
`ParseFrontmatter` view may continue to expose `gate` as an empty scalar value to old
callers.

```yaml
gate:
  version: 1
  stage: ideation
  id: gate:docs-dev:3kd1x1gfxr8mdwzbmnwtjbw8:ideation
  attempt:
    id: briefing:3k-ideation-r1
    round: 1
    reviewed-digest: sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
    resolution:
      id: resolution:captain-3k-r1
      by: person:captain
      at: 2026-07-18T10:30:00Z
      decision: approve
      includes: []
    application:
      action: advance
      target-stage: implementation
      consumption:
        id: application:3k-ideation-r1
        state: pending
      blockers: []
```

This is a current-state record, not an embedded append-only history. Replacing the
current attempt requires a new state commit; the prior exact YAML remains in Git and
projects into immutable historical events.

### Field contract

| Field | Required | Meaning |
|---|---:|---|
| `gate.version` | yes | Integer `1`; unsupported versions fail closed. |
| `gate.stage` | yes | Exact workflow stage whose gate was reviewed. It is independent of current `status` after consumption. |
| `gate.id` | yes | Stable logical gate identity for this entity's occupancy at that gated stage. Re-review rounds retain it. An external gate surface may supply it; applying an existing decision never mints a replacement id. |
| `gate.attempt.id` | yes | Immutable Review & Gate `Briefing.id`; this is the gate-attempt identity. A changed question or reviewed revision set requires a new id. |
| `gate.attempt.round` | yes | Positive workflow review round, an external label as Review & Gate v1 requires. |
| `gate.attempt.reviewed-digest` | yes | `sha256:` digest of the immutable canonical Briefing described below. |
| `gate.attempt.supersedes` | no | Prior attempt id for the same `gate.id`; required when a new round or changed reviewed digest replaces an unconsumed or rejected attempt. |
| `resolution.id` | yes after resolution | Exact portable `Resolution.id`. |
| `resolution.by` | yes after resolution | Actor stamped by the review surface; authority is validated externally and is not self-asserted by this string. |
| `resolution.at` | yes after resolution | RFC 3339 time stamped by the review surface. |
| `resolution.decision` | yes after resolution | Review & Gate v1 `approve`, `revise`, or `hold`. There is no portable `reject`. |
| `resolution.reason` / `includes` | conditional | Preserved rationale fields. `revise` and `hold` require one per Review & Gate v1. |
| `application.action` | yes after resolution | Spacedock interpretation: `advance`, `feedback`, or `none`. Review tooling never executes it. |
| `application.target-stage` | for `advance`/`feedback` | Exact stage expected when the Resolution is consumed. |
| `consumption.id` | yes | Stable one-use application identity minted before any effect; retries and reconciliation retain it. |
| `consumption.state` | yes | `pending`, `prepared`, `consumed`, `ambiguous`, `not-applicable`, or `superseded`. |
| `consumption.dispatch-attempt-id` | for `prepared`/`ambiguous`; optional on `consumed` | Pre-effect identity used at the idempotent or queryable spawn boundary. |
| `consumption.consumed-at` | for `consumed` | Time the stage transition and durable effect receipt were recorded. The event envelope's source commit supplies the canonical commit identity; frontmatter cannot self-reference its own commit SHA. |
| `application.blockers[]` | yes, possibly empty | Committed dispatch prerequisites and their latest durable checks. |
| `application.execution-hold` | no | Workflow-owned pause after approval; distinct from a Review & Gate `hold` decision. |
| `application.feedback` | for `feedback` | Cycle and finding identity used to derive the durable rejection-to-rework route. |

Each blocker has this exact shape:

```yaml
- id: blocker:production-coordinator
  kind: entity-stage
  ref: production-coordinator
  expected-revision: state:4f92c1d
  expected-state: done
  state: unsatisfied
  observed-revision: state:17ae4c0
  checked-at: 2026-07-18T10:30:00Z
  failure-code:
```

`state` is `unsatisfied`, `satisfied`, `unknown`, or `failed`. `unknown` and
`failed` are never eligible. `expected-revision` pins the declaration being checked;
`observed-revision`, `checked-at`, and `failure-code` make the latest committed check
auditable. A scheduler re-check writes a new state commit rather than silently
changing an in-memory answer.

An execution hold has this exact shape:

```yaml
execution-hold:
  id: hold:captain:3k-r1
  state: active
  by: person:captain
  at: 2026-07-18T10:30:00Z
  reason: Approve the design, but do not dispatch until I release this hold.
```

`state` is `active` or `released`. Release preserves the same hold id and adds
`released-by` and `released-at` in a new commit. A later separate hold gets a new id.

For feedback application, the additional shape is:

```yaml
feedback:
  cycle: 1
  finding-ref: resolution:captain-cmd-validation-r1
  finding-digest: sha256:abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd
```

The finding digest covers the deterministic routed-feedback payload produced from the
Resolution reason and included Annotation ids; narrative `### Feedback Cycles` prose
is a human projection, not the route identity.

## Reviewed digest

`reviewed-digest` is SHA-256 over the RFC 8785 JSON Canonicalization Scheme bytes of
the exact Review & Gate v1 `Briefing`. The Briefing already contains its immutable
question, artifact ids and `sha256:` revisions, routing context, criteria, and
evidence. It is therefore the canonical manifest of what the actor reviewed; hashing
an entity file or Stage Report alone would omit gate-defining inputs.

The recorder must resolve and verify every Briefing artifact revision before accepting
the Resolution. A new artifact revision, question, or gate-defining context produces
a new Briefing id, new digest, and new attempt. An approval whose stored digest differs
from the current reconstructed Briefing is stale and cannot be prepared or consumed.

## Record and consume are separate commits

Persisting a binding Resolution never changes `status` and never spawns a worker.

1. **Record:** validate gate/stage/attempt/digest/authority, then commit the `gate`
   mapping with `consumption.state: pending`. Temporary review packages may be deleted
   only after this commit succeeds.
2. **Block or hold:** compute eligibility from the stored decision, digest, blockers,
   execution hold, and consumption state. Persist changed blocker observations or a
   hold release before acting.
3. **Prepare:** for an eligible action that needs a spawn, mint and commit
   `dispatch-attempt-id` with `state: prepared` before the external effect.
4. **Consume:** after the idempotent/queryable effect succeeds, atomically commit the
   expected `status` transition, `state: consumed`, `consumed-at`, and the durable
   receipt under the same consumption/application id. If the effect outcome cannot be
   resolved, commit or project `ambiguous`; never retry under a new identity.

For a feedback action, consumption also records the Feedback Cycle and binds its route
to the target `stage_run_id`. A cycle-3 escalation remains pending human action and has
no target stage run.

## Example 1: approve without dispatch

Before the captain decides, the entity remains at its gate and has no `gate` record:

```yaml
---
id: 3kd1x1gfxr8mdwzbmnwtjbw8
title: Persist gate approval while dispatch blockers remain
status: ideation
score: "0.80"
source: Captain design feedback, 2026-07-13.
worktree:
---
```

After the captain says “approve, but do not dispatch,” recording the Resolution changes
only frontmatter. `status` remains `ideation`; there is no dispatch attempt or worker
receipt:

```yaml
---
id: 3kd1x1gfxr8mdwzbmnwtjbw8
title: Persist gate approval while dispatch blockers remain
status: ideation
score: "0.80"
source: Captain design feedback, 2026-07-13.
worktree:
gate:
  version: 1
  stage: ideation
  id: gate:docs-dev:3kd1x1gfxr8mdwzbmnwtjbw8:ideation
  attempt:
    id: briefing:3k-ideation-r2
    round: 2
    reviewed-digest: sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
    supersedes: briefing:3k-ideation-r1
    resolution:
      id: resolution:captain-3k-r2
      by: person:captain
      at: 2026-07-18T10:30:00Z
      decision: approve
      reason:
      includes: []
    application:
      action: advance
      target-stage: implementation
      consumption:
        id: application:3k-ideation-r2
        state: pending
      blockers: []
      execution-hold:
        id: hold:captain:3k-r2
        state: active
        by: person:captain
        at: 2026-07-18T10:30:00Z
        reason: Approve the design, but do not dispatch until I release this hold.
---
```

Status projects `ideation (approved, execution hold)` from this record. Releasing the
hold does not consume the approval; it only makes the still-current digest eligible.

## Example 2: reject without dispatch

Review & Gate v1 deliberately has no `reject` decision. A captain-facing rejection
that requests rework is stored as portable `decision: revise` plus Spacedock-owned
`action: feedback`. Abandoning the task is a separate workflow action.

Before validation is rejected:

```yaml
---
id: cmd-cutover
title: Cut over CMD production coordinator
status: validation
score: "0.90"
source: CMD cutover
worktree: .worktrees/spacedock-ensign-cmd-cutover
---
```

After the captain rejects the gate, but before feedback routing or worker dispatch,
`status` remains `validation` and the application is visibly pending:

```yaml
---
id: cmd-cutover
title: Cut over CMD production coordinator
status: validation
score: "0.90"
source: CMD cutover
worktree: .worktrees/spacedock-ensign-cmd-cutover
gate:
  version: 1
  stage: validation
  id: gate:cmd:cmd-cutover:validation
  attempt:
    id: briefing:cmd-validation-r1
    round: 1
    reviewed-digest: sha256:23456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef01
    resolution:
      id: resolution:captain-cmd-validation-r1
      by: person:captain
      at: 2026-07-18T11:00:00Z
      decision: revise
      reason: The production coordinator is missing.
      includes: []
    application:
      action: feedback
      target-stage: implementation
      consumption:
        id: application:cmd-validation-r1
        state: pending
      blockers: []
      feedback:
        cycle: 1
        finding-ref: resolution:captain-cmd-validation-r1
        finding-digest: sha256:abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd
---
```

Status projects `validation (rejected; rework pending, cycle 1)`. Only the later
consume commit may change `status` to `implementation`; after it does, status projects
`implementation (rework: validation rejected, cycle 1)` from the same durable gate
attempt and route edge.

## Approval, Review & Gate hold, and execution hold

These are distinct states:

| State | Resolution | Application | Dispatchable |
|---|---|---|---:|
| Approved, dependency blocked | `approve` | blocker unsatisfied/unknown/failed | no |
| Approved, captain says do not dispatch | `approve` | active `execution-hold` | no |
| Approved, ready | `approve` | no active hold; all blockers satisfied; digest current; `pending` | yes |
| Review held | `hold` | `action: none`, `consumption: not-applicable` | no |
| Revision requested | `revise` | `action: feedback`, pending route | no direct forward dispatch |

The prior design modeled only the first row. A captain's “approve but do not dispatch”
is not a dependency and cannot be represented by Review & Gate `hold`, because that
would erase the approval. It therefore requires the first-class durable
`execution-hold` above.

## Supersession and history

- A Resolution is immutable. Changing its id, actor, time, decision, reason, or
  includes under the same attempt id is a conflict.
- Changed reviewed content creates a new attempt with a new Briefing id/digest and
  `supersedes` naming the prior attempt for the same logical gate. The prior
  application's unconsumed state reduces to `superseded`; it cannot dispatch.
- A second Resolution for the same immutable Briefing is not accepted after the
  authorized binding Resolution. Advisory Resolutions remain in the external review
  log and are referenced only through the binding Resolution's provenance.
- A later different gated stage gets its own `gate.id`. Replacing the entity's current
  `gate` mapping does not erase the earlier gate: the preceding state commit remains
  the historical authority and projects into the event stream.
- Git-DAG disagreement is a conflict until an explicit merge-resolution commit selects
  a complete `gate` mapping. Field-wise merging of two attempts is forbidden.

## Commit-derived events and projections

The projector reads the complete old and new `gate` YAML nodes at each reachable state
commit. It emits:

- `gate.attempt_recorded` when a new attempt appears;
- `gate.resolution_recorded` when its immutable `resolution` first appears;
- `gate.application_prepared`, `gate.application_consumed`, or
  `gate.application_ambiguous` from valid consumption transitions;
- `gate.attempt_superseded` from a new attempt's `supersedes` link;
- `feedback.cycle_recorded` and the rejection-to-rework route when a pending feedback
  application is consumed with the matching target stage transition.

The event envelope adds source commit, ordering, and normalized identity; it does not
invent missing gate fields. A status reducer may cache the projection, but deleting
the cache and replaying entity history must reproduce the same snapshot. If the
current entity frontmatter has no gate record, a projected “current gate decision” is
invalid even if a temporary Subspace log or stale cache contains one.

The existing `spacedock-state-commit-event-proposal.md` remains the broader event and
reconciliation design. This file is the physical state contract that proposal was
missing.

## Privacy and portability

The entity stores identities, digests, terse reason text, included Annotation ids,
blocker checks, and effect receipts. It never stores prompts, transcripts, temporary
package paths, pane/session ids, credentials, or private runtime observations. A
decision-log path may be retained only as a portable repository-relative evidence
reference; it is not authority for current workflow state.
