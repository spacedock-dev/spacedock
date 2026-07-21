# Ledger gate binding contract

> **Draft for CL and Jared ratification.** This document and
> [`gate-binding.v1.schema.json`](../schema/gate-binding.v1.schema.json) propose the
> Spacedock-owned boundary. They do not claim a shipped writer, projector, watcher,
> ledger-bound adapter, provider attempt-state schema, command, or endpoint.

Spacedock needs one stable way to bind a portable Ledger gate to workflow state
without taking ownership of the gate or Ledger facts. `spacedock.gate-binding.v1`
does that for a ledger-bound gate attempt. It identifies the Spacedock target before
a Resolution or application exists. Later Ledger facts link through the containing
gate; they do not expand the binding payload. A standalone Spacedock gate attempt is
outside this binding contract and remains valid without Helm Ledger.

## Ownership

| Fact or responsibility | Owner | Contract consequence |
|---|---|---|
| Portable Review & Gate semantics | Subspace R&G | Spacedock does not redefine the Resolution object. |
| `gat_`, `resolution_id`, and `application_id` | Helm Ledger | Spacedock treats these as immutable external identities. |
| `spacedock.gate-binding.v1` | Spacedock | Helm may pin and preserve the schema but does not redefine it. |
| Workflow and entity state | Spacedock | Only the Spacedock-owned writer mutates and commits this state. |
| Apply coordination for an SD-bound gate | Spacedock | The coordinator owns every post-mutation uncertainty until Ledger acceptance is known. |
| Application facts and the apply fold | Helm Ledger | Ledger accepts idempotent attestations and owns `pending_apply -> applied`. |
| Missed-attestation and observation recovery | Reconciliation lane | Recovery can append evidence but cannot invent a Resolution or change authority. |

There is no synchronous transaction across Git and Ledger. The boundary makes each
uncertain window observable and assigns one recovery owner.

## Authority selection

A standalone-only Spacedock deployment remains outside this binding contract. A
future implementation that supports ledger-bound attempts or both modes must select
exactly one decision authority when an attempt opens:

- A standalone attempt explicitly persists standalone authority in the
  provider-owned attempt record before presenting the Briefing or accepting a
  Resolution. It has no Helm Ledger gate or Spacedock binding. Spacedock may persist
  the authoritative R&G Resolution and its local application state. Those are
  Spacedock-local facts, not Helm facts.
- A ledger-bound attempt explicitly persists ledger-bound authority before opening
  a Helm Ledger gate with `spacedock.gate-binding.v1`. Helm owns `gate_id`,
  `resolution_id`, `application_id`, the authoritative Resolution, and the
  application fold. Spacedock persists their references plus workflow-owned state.

The selection is immutable for that attempt and is never inferred from the absence
of a readable Ledger gate or binding. Missing selection, an incomplete ledger-bound
opening, Ledger unavailability, or projection lag all refuse closed. They never
downgrade an attempt to standalone authority. Spacedock must also refuse before
mutation if a ledger-bound attempt presents a Spacedock-local authoritative
Resolution or locally minted Ledger identity. Converting a standalone attempt into
a ledger-bound attempt requires a new attempt or a separately ratified import
contract; it is not an in-place fallback.

A ledger-bound workflow may omit the Resolution body entirely. If it caches a body
for display, the cache must be explicitly non-authoritative, bound to the Ledger
Resolution by digest, and immutable for that `resolution_id`. A digest mismatch
discards or quarantines that cache record; a refetch creates a new cache record
rather than replacing its contents in place. Ledger remains the source of truth.
Provider-owned probe and room history neither selects the authority mode nor becomes
a Ledger fact.

## Ledger-bound adapter preconditions

These requirements constrain a future ledger-bound adapter. They do not define or
ship a provider attempt-state wire format in v1. A dual-mode implementation considers
a standalone attempt ready only after its explicit authority selection is durable. A
ledger-bound attempt uses this opening sequence:

1. Spacedock durably selects ledger-bound authority and records the producer-owned
   `openGate` idempotency key plus the exact command body bytes and their stable
   digest.
2. The opener creates the Helm Ledger gate with `spacedock.gate-binding.v1` using
   that key and body.
3. Spacedock durably records the returned gate and binding reference. A response-lost
   retry reuses the same key and byte-identical body, so Ledger returns the same
   `gate_id`.
4. Only then may the attempt present the Briefing or accept a Resolution or
   application action.

If gate creation or the durable reference write has an unknown outcome, the attempt
stays blocked in opening. The opener reconciles through the durable idempotency key
before continuing or explicitly abandons the attempt. A retry cannot change the body
or create a standalone authority path for the same attempt.

## Binding identity

The minimal provider pin is the existing Helm opaque binding slot:

```json
{
  "ns": "spacedock.gate-binding.v1",
  "entity_ref": "docs/ship-flow/m4.13"
}
```

`ns` selects the Spacedock-owned vocabulary. `entity_ref` identifies the provider
target. Optional `stage`, `target_stage`, `workflow_ref`, `provider_instance_id`, and
`expected_revision` fields can tighten transition checks, routing, and optimistic
concurrency without changing the real two-field Helm slot.

The containing Ledger gate attaches this provider pin to `gate_id`; the binding never
replaces or rekeys the gate. `resolution_id` and `application_id` do not exist when the
gate is opened, so they are forbidden as required binding fields. They appear only in
later Ledger-owned Resolution and application facts. A later authorized application
uses a new `application_id` and may name `supersedes_application_id` in the Helm
application contract.

## Application contracts

The boundary uses these Helm-owned wire contracts after Resolution:

| Reference | Meaning | Fold effect |
|---|---|---|
| `helm.application.committed.v1` | The Spacedock-local canonical state commit and its before/after evidence | Ledger acceptance is the sole `pending_apply -> applied` input. |
| `helm.application.observed.v1` | Later observation of the application on a named ref or remote | Audit and divergence evidence only. |
| `helm.application.view.v1` | Stable read model for response-lost recovery | No new fact; reports the current application state. |

Spacedock consumes these contracts and does not copy their fields into its binding schema.
Committed and observed facts use the same `application_id` but remain separate.
Missing observation does not undo an applied gate. Observation arriving first does
not apply a gate.

For a Git target, the committed fact names `tree_or_blob_kind` (`tree` or `blob`), the
full 40/64 lowercase-hex `tree_or_blob_oid`, and `tree_or_blob_digest`. The digest is
`sha256:` plus lowercase SHA-256 over the exact raw Git object payload bytes named by
that kind and OID. It excludes the loose-object header `<kind> <size>\0` and forbids
pretty-printed/text-normalized output such as `git cat-file -p`. The empty-blob golden
vector uses OID `e69de29bb2d1d6434b8b29ae775ad8c2e48c5391` and payload digest
`sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

## Apply lifecycle

This lifecycle applies only to a ledger-bound attempt:

1. After the gate-opening sequence completes, Ledger records an authorized
   Resolution and exposes a stable `application_id` for required apply.
2. The Spacedock coordinator verifies the provider pin, target state, and required
   Ledger capabilities before mutation.
3. Spacedock writes its workflow state and commits durable references to
   `application_id`, `resolution_id`, and `gate_id`.
4. The same coordinator records `helm.application.committed.v1`; idempotency is scoped
   by `(application_id, idempotency_key)`, and response-loss replay reuses that tuple
   with the byte-identical fact/body digest.
5. Ledger acceptance changes the gate from `pending_apply` to `applied`.
6. A publisher or watcher may later record `helm.application.observed.v1`.

The coordinator must refuse before mutation when application commit or application
recovery is unsupported. An observation capability may be absent without blocking
the committed apply.

A Git hook may wake reconciliation. Correctness must remain the same when hooks are
missing, bypassed, duplicated, or replayed. A hook never mints a Resolution, chooses
authority, or makes Ledger availability part of Git commit success.

## Failure and recovery

| Window | Reported state | Recovery owner and action |
|---|---|---|
| Gate creation or binding persistence has an unknown outcome | Attempt remains in opening; no Resolution is accepted | The opener replays the durable `openGate` idempotency key and byte-identical body to recover the same `gate_id`, then reconciles the provider attempt or abandons it. It never falls back to standalone authority. |
| Binding validation fails before mutation | `binding_refused` | Spacedock refuses closed and returns the typed Ledger or provider validation error. |
| Provider version or digest differs from the pinned contract | `helm.provider_binding.unsupported_version.v1` | Spacedock refuses before mutation and waits for an explicitly supported pin. |
| Spacedock exits before commit | Ledger remains `pending_apply` | Spacedock cleans or reconciles its own worktree, then retries by `application_id`. |
| Commit succeeds and Ledger definitively refuses the attestation | `commit_succeeded_attestation_pending` | Spacedock detects the existing commit and replays the same request after the refusal is addressed. |
| Commit succeeds and the Ledger response is lost | `commit_succeeded_attestation_outcome_unknown` | Spacedock reads `helm.application.view.v1` or replays the same key and body. It does not report a false failure. |
| Observation arrives before committed attestation | `observed_without_attestation` | Spacedock or trusted reconciliation validates durable metadata and supplies the missing committed fact. The gate stays pending until acceptance. |
| Ledger accepts committed attestation but no observation arrives | `applied`, publication evidence absent | No apply recovery is needed. Observation may arrive later. |
| The same idempotency key carries a changed body | `helm.application.idempotency_conflict.v1` | Spacedock refuses the changed replay and reads the current application view. |

The Helm boundary supplies these stable error codes to the coordinator:

- `helm.application.not_found.v1`
- `helm.application.attestation_invalid.v1`
- `helm.application.idempotency_conflict.v1`
- `helm.application.observation_unlinked.v1`
- `helm.provider_binding.unsupported_version.v1`
- `helm.projection.cursor_invalidated.v1`
- `helm.projection.rewrite_quarantined.v1`

Unknown error details are additive and must survive forwarding. An unknown error code
is never interpreted as success.

Helm pins the supported provider namespace and schema digest outside the binding
payload. A different `ns` version or a digest mismatch returns
`helm.provider_binding.unsupported_version.v1` before mutation. Spacedock must not
guess fields, fall back to a nearby version, or fetch live upstream canon. The same
typed refusal covers version and digest mismatch so consumers have one fail-closed
recovery branch.

## Projection rewrites

The Ledger projection that exposes a binding records `projection_epoch` and
`exposed_binding_digest` outside the binding payload. Event cursors are valid only
inside their epoch. An epoch mismatch returns
`helm.projection.cursor_invalidated.v1` with restart guidance; it must not look like
a successful empty page.

If reprojection selects a different binding, the old exposed binding stays frozen.
The projection enters `rewrite_quarantined`, records the exact
`candidate_binding_digest`, and waits for explicit reconciliation. Both binding
digests are `"sha256:"` followed by the lowercase hexadecimal SHA-256 of the
UTF-8 bytes produced by the RFC 8785 JSON Canonicalization Scheme (JCS) for their
corresponding binding objects. Producers and consumers must use a compliant JCS
implementation; ordinary JSON serialization, including Go `encoding/json` or
JavaScript `JSON.stringify`, is not a substitute. A binding value that cannot be
represented as RFC 8785 I-JSON fails closed before any mutation, with no fallback
digest path. Contract vectors cover UTF-16 property ordering, string escaping, and
numeric representation. Reprojection cannot silently mint a new application
obligation or redirect an existing one.

An applied application never enters rewrite_quarantined. Its committed receipt remains
append-only proof; later binding/source changes append supersession or reconciliation
evidence and, if needed, create a new application obligation. Quarantine is only for an
uncommitted projection whose candidate differs from the frozen exposed binding.

## Orphan observations

An observation without `application_id` is not a binding and cannot apply a gate.
Helm may expose readable `helm.application.orphan_observed.v1` evidence keyed by
immutable source identity. Its writer remains a typed, reconciliation-internal
operation for trusted watcher or backfill adapters.

Spacedock must not advertise the orphan writer as a CLI, HTTP endpoint, or public
capability. The public capability set remains application commit, linked observation,
application recovery, and projection epoch support.

## Operation without Spacedock

When Spacedock is absent, this binding is absent. Helm keeps the existing Ledger path:

- a gate with `apply_requirement: none` ends at `decided`;
- an action gate can bind to its session executor and use the provider-neutral Helm
  attestation contract;
- Helm does not create synthetic Spacedock workflow state or a fake provider pin.

This keeps Spacedock optional without creating a second decision or application
model.

## Operation without Helm Ledger

When Helm Ledger is absent, the Spacedock binding is absent. Spacedock may run the
standalone authority mode described above only after explicitly persisting that
selection. Ledger unavailability or an unreadable binding is not standalone
selection. A standalone attempt may include an authoritative R&G Resolution and
local application state in its workflow record. It must not mint `gat_`,
`resolution_id`, or `application_id` values that appear to be Helm Ledger facts.

Adding Helm later does not retroactively convert the attempt. A new ledger-bound
attempt may reference the same provider target and Briefing lineage, but it receives
new Ledger-owned identities and authority.

## Compatibility rules

The schema follows JSON Schema 2020-12 and accepts the existing minimal Helm binding
slot. It permits additive fields so a consumer can preserve extensions it does not
understand. Consumers must validate known required fields and constants, retain
unknown fields on read/write round trips, and ignore unknown optional fields unless
another negotiated capability assigns them meaning.

`entity_ref` and every known optional string routing field must contain at least one
non-whitespace character. A present whitespace-only field fails validation; omission
remains valid for optional fields.

An additive optional field is compatible within v1. Removing or renaming a required
field, changing a constant, changing authority, or changing a fold effect requires a
new schema version. New versions must not reuse `spacedock.gate-binding.v1`.

The standalone versus ledger-bound clarification changes neither the v1 schema
document nor binding payload bytes. It defines preconditions for a future
ledger-bound adapter; the versioned provider attempt-state representation remains a
separate ratification surface. In a dual-mode implementation, the explicit
provider-owned attempt selection chooses the mode. A valid binding on the containing
Ledger gate confirms the ledger-bound path, but binding absence alone never selects
standalone authority. A hybrid record is invalid rather than a third mode.

The executable examples live in
`internal/contract/testdata/ledger-boundary/`. They include the existing minimal
binding, an extension-preservation case, invalid provider target cases, and a
boundary fixture that links later Resolution, application, committed, and observed
facts without putting those fields inside the binding. Additional fixtures prove
same-key replay, changed-body conflict, response-lost application reads, epoch cursor
invalidation, rewrite freeze, and typed provider-pin refusal.

`internal/contract/testdata/ledger-boundary/helm-schemas/` vendors the exact Helm
application and rewrite schemas used by those executable examples. Its manifest pins
the source repository, source path, Git revision and blob identity, and raw SHA-256
for each of the eight schemas. Contract tests validate the complete request fact,
receipt, observation, application view, source-superseded evidence, generic provider
binding, and rewrite evidence through those schemas without a live upstream read.
They recursively walk every Helm `$ref` and reject an incomplete vendored schema
closure, including a reference reachable only through an optional property.

The manifest uses `provenance_status: complete` and pins Helm source revision
`e852f31e0ae6c5f06c44b51b5c7a82d0edc7da7a`. Every schema entry records its exact
source path, Git blob ID, and raw SHA-256 so consumers can independently reproduce
the vendored bytes from the immutable source commit.

## Consumer checklist

A consumer of `spacedock.gate-binding.v1` must:

- determine whether the gate attempt is standalone or ledger-bound before mutation
  by reading its explicit durable selection, never by treating an unavailable or
  absent Ledger binding as standalone;
- keep that authority selection fixed for the attempt and block while a ledger-bound
  opening outcome is unknown;
- persist the producer-owned `openGate` idempotency key, exact body bytes, and stable
  body digest before calling Ledger, then reuse the same key and byte-identical body
  for response-lost recovery;
- refuse a ledger-bound record that also claims a Spacedock-local authoritative
  Resolution or locally minted Ledger identity;
- verify required `ns` and `entity_ref`, plus any supplied optional target fields,
  before using the provider pin;
- read `gate_id`, `resolution_id`, and `application_id` from their Ledger-owned facts,
  never by requiring them inside the binding;
- validate target state and required capabilities before mutation;
- persist `application_id` in durable mutation evidence;
- recompute Git tree/blob evidence from the declared kind, full object OID, and raw
  payload bytes without a loose-object header or text conversion;
- treat committed acceptance as the only apply fold and observation as audit only;
- recover a lost response by application read or replay of the same
  `(application_id, idempotency_key)` tuple and byte-identical body;
- preserve unknown additive fields;
- stop on projection invalidation or rewrite quarantine.

## Contribution checklist

A wire contract change must update the spec, schema, relevant valid and invalid
fixtures, and `internal/contract/gate_binding_test.go` in the same pull request.
Contributors must state whether the change is additive within v1 or requires a new
version. A prose-only scope or consumer clarification leaves the schema document
unchanged unless all raw-digest consumers update their pins in the same change. A
Helm wire-shape change also updates the vendored schema bytes, manifest digest, and
all conforming Spacedock fixtures in the same pull request. Review-ready provenance
requires `provenance_status: complete`, the exact 40-character Helm source commit,
and the 40-character Git blob ID plus raw SHA-256 for every vendored schema.

Review must confirm that Helm-owned attestation fields were not copied into the
Spacedock schema, the orphan writer did not become public, unknown fields still round
trip, unsupported provider pins fail closed, and no document claims the Draft runtime
exists.
