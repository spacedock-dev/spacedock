# Gate Resolution frontmatter contract

Status: implemented recorder contract
Date: 2026-07-22

## Outcome and ownership

The recorder makes a captain's decision durable without changing workflow status,
routing feedback, creating application state, or dispatching work. It owns the logical
gate, ordered attempts, immutable Briefing binding, and portable Resolution. The
application layer owns what an approval does and may add an opaque `application` subtree
to a closed attempt; the recorder preserves that known boundary and freezes it with the
closure.

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
            room-ref: ./review/validation/briefing-1
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
migration or compatibility rewrite. The one intentionally opaque exception is the known
`application` field, whose nested keys belong to the application layer.

Every Briefing binding includes an id, SHA-256 digest, explicit digest domain, and room
reference. The approved domains are:

- `canonical-bytes`: SHA-256 over RFC 8785/JCS canonical Briefing JSON bytes. New
  recorder binds always use this domain.
- `raw-file-pin`: an explicitly labelled raw-byte pin that may remain in a canonical v1
  record. It is never silently reinterpreted as a canonical digest.

## Recorder lifecycle

`spacedock gate record` derives lifecycle under the entity lock:

1. `--briefing` derives the logical gate from the entity's current workflow stage.
2. With no record for that stage, it opens the first attempt.
3. With an open last attempt, an identical binding is a no-op and a changed binding
   replaces that attempt's Briefing.
4. With a closed last attempt, it appends a successor. Existing closed attempts,
   Resolutions, and application data are frozen.
5. A Result or chat decision closes only the last open attempt for the current stage.

Cross-logical-gate re-entry is ordinary: workflow stage selects the target record even
when `gates.current.gate` names a different closed gate. The successful write selects the
target record but does not modify either record's earlier closures.

## Provider Result association

The provider form consumes exact `review-v1-result` bytes plus a retained
`spacedock-result-association` v1. The association binds:

- the raw Result digest and provider Briefing id;
- the authorized actor;
- the canonical Briefing id and canonical revision;
- the canonical artifact id/revision list; and
- a one-to-one presentation mapping from provider artifacts to canonical artifacts.

The association is not trusted to declare package completeness. The recorder resolves
the bound room's exact `briefing.json`, recomputes its frozen JCS digest, checks its id,
and derives the complete artifact inventory from those independently authenticated
bytes. The association's canonical list must equal that inventory and its presentation
mapping must cover every inventory item exactly once, including the Result's primary
artifact. Only after all checks and actor authorization pass does the recorder normalize
the provider Resolution's Briefing id to the canonical binding. Advisory results also
require an adoption note naming the authorizer. Artifact payloads may remain external
URI + SHA references; the recorder does not copy them merely to establish inventory.

## Write boundary and invariants

The recorder rebuilds only the canonical `gates:` subtree. Before atomic replacement it
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
spacedock gate record ENTITY --briefing FILE [--workflow-dir DIR]
spacedock gate record ENTITY --result FILE --association FILE --actor ID [--adoption-note TEXT] [--workflow-dir DIR]
spacedock gate record ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--directive TEXT] [--workflow-dir DIR]
spacedock gate validate ENTITY [--workflow-dir DIR]
```

Exactly one semantic source is required. The binary derives operation, ids, stage target,
and compare-and-swap state; callers cannot submit an operation envelope or candidate
identities. `gate validate` is read-only and reports the selected record's last attempt.

Delegated chat decisions require the quoted directive. A delegated First Officer
approval also requires a reason. The recorder constructs the portable Resolution and
records it under the identity that actually rendered the decision; it does not apply the
result.

## Explicitly outside v1

- Prototype-format compatibility, migration, and arbitrary unknown-field preservation
  inside `gates`.
- Provider launch, polling, result retention, presentation UI, and Subspace-specific
  behavior.
- Application lifecycle, workflow transition, dispatch, or effect receipts.
- A second schema version or provider operation envelope.

## Behavioral proof

The release tests must fail if any of these outcomes regress:

1. Open/rebind/close/successor behavior changes, or a successor mutates a frozen closure.
2. Cross-gate re-entry targets global selection instead of current workflow stage.
3. A prototype field or arbitrary unknown binary-owned field becomes readable or writable.
4. A stale, invalid, or lock-contended write changes the entity.
5. A canonical write changes bytes outside `gates` or alters opaque application data.
6. Removing canonical artifacts and matching presentation entries together is accepted.
7. Result identity is normalized before exact bytes, authority, bound Briefing digest,
   complete inventory, and full presentation mapping are verified.
