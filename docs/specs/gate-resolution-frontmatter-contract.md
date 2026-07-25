# Gate Resolution frontmatter contract

Status: implemented recorder and application contract
Date: 2026-07-22

## Outcome and ownership

The recorder makes a captain's decision durable before any workflow status change or
dispatch. It owns the logical gate, ordered attempts, immutable Briefing binding, and
portable Resolution. The application layer owns what an approval does through the typed
`application` subtree on a closed attempt; the same recorder binary owns its guarded
writes.

Presentation remains an overridable channel of the present-gate skill, not a recorder
verb. Chat and provider channels both hand semantic decision input to the recorder.
Provider transport, retention, and UI stay outside this binary.

## Canonical v1 schema

The binary accepts and emits one canonical `gates:` shape:

```yaml
gates:
  version: 1
  current:
    gate: gate:example:sample:validation
  records:
    - id: gate:example:sample:validation
      stage: validation
      attempts:
        - id: gate-attempt:sample-validation-1
          briefing:
            id: briefing:sample-validation-1a
            digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
            digest-domain: canonical-bytes
            request-digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
            room-ref: ./review/validation/briefing-1
          provider-evidence:
            result-digest: sha256:5555555555555555555555555555555555555555555555555555555555555555
            presented-inventory-digest: sha256:6666666666666666666666666666666666666666666666666666666666666666
          resolution:
            type: Resolution
            id: resolution:captain-sample-validation-1a
            briefing: briefing:sample-validation-1a
            by: person:captain
            at: 2026-07-22T09:00:00Z
            decision: approve
          application:
            action: advance
            target-stage: done
            state: pending
            blockers: []
```

`records` and `attempts` are ordered. The last attempt in a record is current.
Resolution absence means open; Resolution presence means closed. `gates.current.gate`
selects the logical gate eligible for later application. These facts remove any need for
separate attempt pointers, sequence numbers, lineage pointers, or explicit lifecycle
state.

The binary-owned model is closed: unsupported fields inside `gates` fail validation.
In particular, the pilot-only `gates.current.attempt`, `current-attempt`, `sequence`,
`previous-attempt`, and explicit attempt `state` encodings are rejected. There is no
migration or compatibility rewrite. The `application` field is the typed one-use
lifecycle boundary owned by the application layer on the same canonical-v1 writer surface.

Every Briefing binding includes an id, SHA-256 digest, explicit digest domain, and room
reference. The approved domains are:

- `canonical-bytes`: SHA-256 over RFC 8785/JCS canonical Briefing JSON bytes. New
  recorder binds always use this domain.
- `raw-file-pin`: an explicitly labelled raw-byte pin that may remain in a canonical v1
  record. It is never silently reinterpreted as a canonical digest.

A prepared provider room also binds `request-digest`, the JCS digest of its
`request.json`. Chat-only attempts may omit it. Changing request authority after the
attempt binds therefore fails before Result validation or entity mutation.

## Recorder lifecycle

`spacedock gate record` derives lifecycle under the entity lock:

1. `--briefing` requires the retained package manifest basename `briefing.json`, then
   derives the logical gate from the entity's current workflow stage. Any other basename
   is rejected before locking or mutation.
2. With no record for that stage, it opens the first attempt.
3. With an open last attempt, an identical binding is a no-op and a changed binding
   replaces that attempt's Briefing.
4. With a closed last attempt, it supersedes any pending application and appends a
   successor. Existing Briefings and Resolutions remain frozen.
5. A Result or chat decision closes only the last open attempt for the current stage and
   derives its `advance/pending`, `feedback/pending`, or `none/not-applicable` application.

Cross-logical-gate re-entry is ordinary: workflow stage selects the target record even
when `gates.current.gate` names a different closed gate. The successful write selects the
target record but does not modify either record's earlier closures.

## Provider Result association

The provider form consumes one prepared gate room. Its frozen `request.json` binds the
logical gate, attempt, canonical Briefing id and digest, and captain actor/approver authority.
Delegated chat decisions record `agent:first-officer` with a nonblank evidence reason and no directive or `adoption-note`;
the attempt's `request-digest` rejects post-binding request changes.
The fixed provider outputs are `provider/result.json` and
`provider/presented-inventory.json`; callers supply neither path nor provider argv.

The recorder resolves the room's exact `briefing.json`, recomputes its JCS digest, and
derives the canonical inventory from every Artifact and recursively reached Reference.
It derives a private `spacedock-result-association` v1 by matching each presented id and
revision to that inventory and binding the raw Result digest. The mapping must cover the
whole inventory exactly once, including the Result's primary Artifact.

A direct binding Result uses Review v1's minimal envelope: authority comes from nested
`Resolution.by`, and redundant `status`, `binding`, `actor`, `approver`, or
`resolutionId` fields are absent. The recorder requires `Resolution.by` to equal the
request authority. Advisory output remains retained evidence; no adoption note can
promote it into a binding Resolution. Artifact payloads may remain external URI and
SHA references.

On provider close, the recorder stores only the raw-byte digests of
`provider/result.json` and `provider/presented-inventory.json` as `provider-evidence`.
They are part of the frozen attempt. `gate validate` recomputes both from the fixed room
files and fails if either is missing or changed. Chat-closed and open attempts carry no
provider evidence; the derived association remains ephemeral.

## Round records and triage dispositions (advisory; owner: 02av)

A correction round reuses the recorder vocabulary without becoming a gate: its reviewed
snapshot is an immutable, digest-bound Briefing; reviewer findings are same-Briefing
Annotations; and the reviewer verdict is an advisory Resolution. The worker's triage is
a separate advisory Resolution on that Briefing. Round Resolutions carry no application,
do not select a logical gate, and cannot advance workflow status.

For every correct-but-disproportionate finding, the triage Resolution `includes` a
worker-authored Annotation that itself includes the reviewer's finding Annotation. Its
body records the class, why the finding is not material to an entity value AC or
non-negotiable boundary, and the condition that promotes it to material. Material
findings are fixed; needs-decision findings are escalated. A finding neither fixed nor
represented by the triage Resolution is not triaged.

```yaml
- type: Annotation
  id: annotation:decline-symlink-prototype
  briefing: briefing:02av-implementation-round-1
  by: actor:ensign
  includes: [annotation:finding-symlink-prototype]
  body: >
    class: correct-but-disproportionate; why-not-material: no value AC breaks and the
    crafted-symlink trigger is outside the supported flow; promotes-when: a released
    user reaches it through an operator-selected repository.
- type: Resolution
  id: resolution:ensign-02av-implementation-round-1
  briefing: briefing:02av-implementation-round-1
  by: actor:ensign
  decision: revise
  reason: "triage: 0 material fixed; 1 declined"
  includes: [annotation:decline-symlink-prototype]
```

No findings means no triage Resolution. An all-declines round instead has a real triage
Resolution recording zero fixes and including every decline Annotation; those states
must never project alike. Once reviewer and authorized worker triage entries are
complete, `spacedock gate record <entity> --round STAGE/CYCLE --briefing
PATH/briefing.json --log PATH/briefing.review.jsonl --feedback-cycle FILE` publishes
the canonical two-file room at `review/<stage>/round-<cycle>`, then atomically writes
the exact `review-round` pointer and Feedback Cycles projection. A complete no-findings
log omits both worker triage and `--feedback-cycle`.

Round recording requires a folder-form entity at `<slug>/index.md`, so its accumulating
`review/` artifacts are scoped beside that entity. Flat entities refuse before locking
or writing; the recorder does not alter the approved derived room path to compensate.
`STAGE` must name a stage in the workflow definition, but need not equal current
`status`: explicit historical backfill remains supported. Decline bodies must use the
exact structured class/rationale/promotion fields above with substantive values. A
projected Feedback Cycles line must match the complete documented grammar, its cycle
must equal `CYCLE`, and its verdict must agree with the reviewer Resolution.

The room is immutable: exact whole-room replay is a whole-tree no-op; any different
Briefing, log, room shape, pointer, or projection fails closed. Findings-bearing
reviewer-only logs are incomplete and never persist. New-room publication rolls back
if the full-entity compare-and-swap or atomic pointer/projection replacement fails.
`spacedock gate validate <entity> --round STAGE/CYCLE` reads the ordered log through
the pointer and reports every Resolution as advisory.

Narrowing a value AC to make a finding pass is not a round disposition. It opens a real
gate attempt whose binding Resolution is captain-owned; the correction loop cannot
self-approve that design reset.

## Write boundary and invariants

Ordinary recorder operations rebuild only the canonical `gates:` subtree. Consumption
co-writes `status` and `application.state: consumed` in one atomic replacement. Before replacement it
validates the rebuilt full entity and compares the locked source subtree with the one it
read, so stale or invalid writes fail without replacing the file. All frontmatter fields
outside `gates` and the Markdown body are preserved byte-for-byte. The per-entity lock
rejects concurrent recorder writers; there is no retry, lease, daemon, or recovery
protocol.

The model enforces unique gate, attempt, Briefing, and Resolution ids; a resolvable
current logical gate; non-empty attempt histories; exact Resolution-to-Briefing binding;
and portable `approve`, `revise`, or `hold` decisions. `revise` and `hold` require a
reason or an included same-Briefing Annotation. Open attempts cannot carry application
data.

## Command surface

```text
spacedock gate record ENTITY --briefing PATH/briefing.json [--workflow-dir DIR]
spacedock gate record ENTITY --room PATH [--workflow-dir DIR]
spacedock gate record ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--workflow-dir DIR]
spacedock gate validate ENTITY [--workflow-dir DIR]
spacedock gate eligibility ENTITY [--workflow-dir DIR]
spacedock gate consume ENTITY [--workflow-dir DIR]
```

Exactly one semantic source is required. Supported chat actor IDs are `person:captain` and `agent:first-officer`. The binary derives operation, ids, stage target,
and compare-and-swap state; callers cannot submit an operation envelope or candidate
identities. `gate validate` is read-only and reports the selected record's last attempt.

New delegated chat resolutions use `by: agent:first-officer`, require a nonblank
evidence reason, and omit `adoption-note`. Readers continue to accept historical
`adoption-note` values, but the chat recorder neither emits nor treats them as
authority. The recorder constructs the portable Resolution under the asserted identity
that rendered the decision; it does not authenticate chat or apply the result.

## Explicitly outside v1

- Prototype-format compatibility, migration, and arbitrary unknown-field preservation
  inside `gates`.
- Provider launch, polling, result retention, presentation UI, and Subspace-specific
  behavior.
- Blocker-satisfaction evaluation, execution-hold authoring, dispatch identities, or effect receipts.
- A second schema version or provider operation envelope.

## Behavioral proof

The release tests must fail if any of these outcomes regress:

1. Open/rebind/close/successor behavior changes, or a successor mutates a frozen closure.
2. Cross-gate re-entry targets global selection instead of current workflow stage.
3. A prototype field or arbitrary unknown binary-owned field becomes readable or writable.
4. A stale, invalid, or lock-contended write changes the entity.
5. A canonical write changes bytes outside `gates` or alters opaque application data.
6. Removing canonical artifacts and matching presentation entries together is accepted.
7. A room Result closes before exact bytes, request authority, bound Briefing digest,
   complete Artifact/Reference inventory, and full presentation mapping are verified.
