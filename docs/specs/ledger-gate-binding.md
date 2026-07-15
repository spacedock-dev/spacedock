# Ledger gate binding contract

> **Draft for CL and Jared ratification.** This document and
> [`gate-binding.v1.schema.json`](../schema/gate-binding.v1.schema.json) propose the
> Spacedock-owned boundary. They do not claim a shipped writer, projector, watcher,
> command, or endpoint.

Spacedock needs one stable way to bind a portable Ledger gate to workflow state
without taking ownership of the gate or Ledger facts. `spacedock.gate-binding.v1`
does that. It identifies the Spacedock target before a Resolution or application
exists. Later Ledger facts link through the containing gate; they do not expand the
binding payload.

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

1. Ledger opens the gate with the Spacedock binding, then records an authorized
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

A contract change must update the spec, schema, relevant valid and invalid fixtures,
and `internal/contract/gate_binding_test.go` in the same pull request. Contributors
must state whether the change is additive within v1 or requires a new version. A Helm
wire-shape change also updates the vendored schema bytes, manifest digest, and all
conforming Spacedock fixtures in the same pull request. Review-ready provenance
requires `provenance_status: complete`, the exact 40-character Helm source commit,
and the 40-character Git blob ID plus raw SHA-256 for every vendored schema.

Review must confirm that Helm-owned attestation fields were not copied into the
Spacedock schema, the orphan writer did not become public, unknown fields still round
trip, unsupported provider pins fail closed, and no document claims the Draft runtime
exists.
