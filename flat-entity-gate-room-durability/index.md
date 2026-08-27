---
id: b9pjkz3rv0svx9rt63kw8yg7
title: "A flat entity's gate room is neither committed by `state commit` nor migratable to folder form"
status: ideation
source: "Observed eight times driving gates on 79 and cn, 2026-07-27/28. Every gate room bound on a flat entity needed a manual path-scoped commit after `state commit` reported success; the obvious repair — converting to folder form — silently breaks prior closed gates."
started: 2026-08-27T17:15:11Z
completed:
verdict:
score: 0.7
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:b9pjkz3rv0svx9rt63kw8yg7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:b9pjkz3rv0svx9rt63kw8yg7-backlog-1
              briefing:
                id: briefing:b9pjkz3rv0svx9rt63kw8yg7:backlog:attempt-1:revision-1
                digest: sha256:5d63d8641afab33a2f57b0f4054d12583c8e1b6fd4e3d043f94ce0c6edfee97b
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:b9pjkz3rv0svx9rt63kw8yg7:backlog:1
                briefing: briefing:b9pjkz3rv0svx9rt63kw8yg7:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T17:14:04.981857Z"
                decision: approve
                reason: 'Captain approved the bound Subspace backlog review: support flat tickets directly, preserve historical room-ref meaning, and make sibling review rooms durable and round-recordable without migration.'
              application:
                target-stage: ideation
                state: consumed
---

Make a gate room on a flat entity durable through the supported command, and give a flat entity a migration path to folder form that does not corrupt its gate history.

## Two defects that compose into a cul-de-sac

A flat entity (`<slug>.md`) may legitimately carry a sibling room directory (`<slug>/review/...`). The shipped precedent is `source-build-compatibility-identity`, whose backlog briefing lives exactly that way. Both defects below apply to that shape.

### 1. `state commit` reports success having committed only the index

Observed on every gate room bound during the 2026-07-27/28 session — eight occurrences across `79` and `cn`, at backlog, ideation and validation gates. `spacedock state commit <slug>` commits `<slug>.md` and pushes, printing `Committed and pushed …`, while `<slug>/review/<stage>/briefing-N/` remains untracked.

The consequence is not cosmetic. The committed frontmatter carries a `room-ref` pointing at a room that does not exist on the state remote, so a peer or a later session cannot retrieve the Briefing that justified a recorded decision. That is the "digests no committed tree reproduces" failure class the durable-decisions sprint was created to eliminate, reproduced through the supported path.

Every occurrence was worked around by hand:

```
git add -- <slug>/ && git commit -m "…" -- <slug>/ && git push origin HEAD:spacedock-state/dev
```

`sync-merge-guard-archive-state` (`rd`) owns the folder-form version of this and explicitly scopes its work to folder-form entities. The flat-plus-sibling-room shape is uncovered by it.

### 2. Converting to folder form silently breaks every prior closed gate

The obvious repair is to convert the entity to folder form, which also unlocks advisory round recording (the recorder requires `<slug>/index.md` and refuses flat entities before locking or writing). It cannot be done safely today.

`room-ref` resolves relative to the **entity's own directory**. Proven against `prepare-provider-neutral-gate-room`, a folder-form entity whose stored `./review/backlog/briefing-1` maps to `<slug>/review/backlog/briefing-1`. For a flat entity the entity directory is the state root, so a correct stored ref reads `./<slug>/review/...`. After `git mv <slug>.md <slug>/index.md` that same ref resolves to `<slug>/<slug>/review/...` — one level too deep, and pointing at nothing.

Nothing catches it. `gate validate <entity>` exits 0 and reports gate, attempt, briefing, resolution and decision, because it reads frontmatter without dereferencing `room-ref`. Observed directly: after converting `79`, validate stayed green while the backlog gate's room pointer was broken.

There is no supported repair once broken. `gate record --briefing` will not rewrite a closed attempt, and hand-editing `gates:` is forbidden by the lifecycle skill. The conversion on `79` was reverted for exactly this reason.

### Why the two compose badly

Round recording requires folder form. A flat entity therefore cannot machine-record advisory rounds — and the migration that would fix that corrupts its closed gate records. An entity filed flat is permanently stuck with hand-authored rounds, and the only escape damages its decision history.

## What a fix needs to decide

The selected design makes migration unnecessary. A ticket's **review home** is always
`<state-root>/<slug>/`, independent of whether its mutable entity is the flat
`<state-root>/<slug>.md` or folder-form `<state-root>/<slug>/index.md`. Gate Briefings
and advisory-round rooms live below that home at `review/<stage>/...`.

One new standard-library-only package, `internal/entityhome`, owns that identity. Its
resolver receives the state root, slug, and an explicit operation mode. It returns the
unique entity form, canonical entity path, review home, and literal commit paths. It
refuses an active flat-plus-folder collision rather than preferring either form. An
existing entity leaf must be a non-symlink regular file; an existing review home and
each existing descendant on the requested room path must be non-symlink directories.
An absent flat companion is allowed only in the explicit controlled-creation mode used
by prepare and round publication. A tracked deletion is allowed only when the Git
caller proves the exact missing canonical path; it never enters creation mode. These
separate modes preserve tracked deletion without making an arbitrary missing entity
valid.

Every mutating consumer resolves before it locks, stages, commits, or creates a room.
Resolution rejects absolute or escaping operands, a non-directory home, symlinked
entity/home/room components, and a home outside the state root. The resolver does not
follow a symlink and then bless its external target as the boundary. Prepare, round,
`state commit`, and `dispatch build --stamp` use the same result, so one slug cannot
acquire different ownership under different commands.

Stored `room-ref` keeps its existing meaning: it is relative to the entity file's own
directory. Therefore the same physical room is stored as `./<slug>/review/...` for a
flat entity and `./review/...` for a folder entity. Existing bindings need no rewrite,
the frozen-record rule remains intact, and a flat ticket can stay flat for its whole
lifecycle. There is no flat-to-folder migration verb in this task; a manual move that
changes the base directory remains unsupported. A workflow's `entity-form: folder`
filing rule also remains authoritative: this design supports grandfathered and
otherwise-valid flat tickets, but does not make a declared-folder workflow accept a
newly misfiled flat ticket.

That filing rule is one shared gate-package policy check, not a preparation-only
special case. A declared-folder workflow refuses a flat entity whose valid review home
does not yet contain a real `review/` directory. It allows a grandfathered flat entity
whose canonical home already contains that real directory. An undeclared workflow
allows either valid form. Prepare and round call the same check before mutation, so a
round cannot create the first companion that preparation refuses.

`gate record --round` must derive both its room and pointer from that review home. It
must stop assuming `filepath.Dir(entityPath)` is the review home or that the pointer is
always `./review/...`. For a flat entity it publishes
`<slug>/review/<stage>/round-<cycle>` and stores
`./<slug>/review/<stage>/round-<cycle>`; folder behavior and bytes stay unchanged.
Artifact containment is rooted at the review home, not the state root, so a flat round
cannot cite a sibling entity. The mutable flat `<slug>.md` is outside that home and is
also explicitly forbidden as a round artifact.

Explicit `status --validate` becomes the shipped integrity check over retained evidence,
independent of gate-record presence. It inspects `gates` and `review-round` separately
for both active and archived entities. A normalized local `./...` ref resolves relative
to the historical entity-directory base and must remain inside that base. This preserves
legacy flat local-file refs outside the newer review home without permitting traversal.
Absolute, traversal, backslash, non-normalized, escaping, symlinked, missing, and
wrong-type local refs fail validation with the entity path and offending ref. A local
gate binding must resolve through its historical file-or-room reader; a round ref must
name its exact immutable two-file room. An opaque provider ref remains non-filesystem
evidence and is never dereferenced. Ordinary status stays nonblocking. Existing gate
mutation paths keep their stronger Briefing and digest checks.

### Why this design

- The review-home derivation serves AC-1 and AC-2. The simpler alternative, requiring
  folder form and moving every flat entity, is insufficient because it either strands
  existing flat tickets or rewrites frozen bindings.
- Shape-specific relative refs serve AC-2 and AC-3. A new state-root-relative ref
  scheme is insufficient because it changes the meaning of every existing folder ref;
  a single fixed `./review/...` ref is wrong for flat entities under the historical
  resolver.
- Validation dereference serves AC-4. Keeping the current warning is insufficient
  because exit 0 lets automation accept absent evidence; blindly promoting the old
  join is also insufficient because it would reject opaque provider history and miss
  round-only or archived state.
- Literal flat commit units serve AC-1 and AC-5. A top-level `git add -A` is
  insufficient because it can attribute a concurrent sibling writer's state to the
  wrong ticket.
- The shared resolver and form-policy check serve AC-2 and AC-5. Keeping four local
  derivations is insufficient because current state commit prefers folder form at an
  active collision while companion discovery follows symlinks; a prepare-only policy
  also leaves round publication as a bypass.

No separate throwaway spike is needed because the risky filesystem mechanisms are
already exercised in the repository:
`TestPrepareCreatesOneFileRecorderRoomForFolderAndFlatEntities` proves the common
physical review home and both historical ref spellings;
`TestStateCommitFlatIncludesExactCompanionDirectoryAndTrackedDeletions` proves the
flat commit unit and tracked deletion; and
`TestPrepareRejectsSymlinkedFlatCompanionWithoutChangingBytes` proves the required
byte-clean `Lstat` boundary behavior. `TestRoundRequiresFolderFormWithoutCrossEntityCollision`
pins the refusal that must become the policy matrix. These focused baselines passed on
2026-08-27. The implementation's first test must be the shared resolver matrix; the
first production change must make that independent table green before any consumer is
rewired.

## Consumer inventory and boundaries

### Canonical identity and form policy

- New `internal/entityhome` owns active-form collision detection, regular-file entity
  requirements, review-home derivation, symlink/non-directory refusal, state-root
  containment, controlled flat-home creation, and the Git-proven tracked-deletion mode.
- `internal/gates/prepare.go` retains the README parser but replaces
  `refuseNewFlatCompanion` with one shared form-policy check. Prepare and round receive
  the same allow/refuse result after canonical resolution and before mutation.
- `internal/status/validate.go` keeps the workflow-wide flat/folder diagnostic, while
  the shared resolver owns whether a specific operation has one safe active identity.

### State commit and publication

- `internal/cli/state_sync.go`: `runStateCommit`, `commitInlineEntity`,
  `syncActiveEntity`, `resolveEntityCommitTarget`, `flatEntityCommitPaths`, and
  `commitEntityPathsScoped` consume the canonical resolver. Existing entities require
  regular non-symlink leaves and unambiguous form. Git-proven tracked deletions retain
  their exact literal path plus a tracked companion without authorizing creation.
- `internal/cli/gate_ceremony.go` consumes `syncActiveEntity` for implicit split-root
  record/consume publication and therefore inherits the same flat unit.
- `internal/dispatch/stamp.go`: `commitAndPublishEntity` and its duplicate
  `entityCommitPaths` are the `dispatch build --stamp` consumer. They must consume the
  canonical resolver and match state-commit collision, symlink, non-directory, and
  literal-scope behavior rather than keeping a second `os.Stat` derivation.
- Archived publication and `internal/status/merge.go` archive moves remain out of
  scope; they do not publish a newly recorded active flat round.

### Advisory rounds

- `internal/cli/cli.go:newGateCommand` routes `gate record --round` and performs the
  immediate readback; its grammar and output stay unchanged.
- `internal/gates/operation.go:RecordSemantic` selects the round recorder and applies
  the shared form-policy preflight before the round lock or write. Authority and status
  behavior stay unchanged.
- `internal/gates/round.go`: `resolveRound` derives the review home and room;
  `recordRoundLockedWith` derives the shape-correct ref; `readRoundPointerData` and
  `ValidateRoundFile` validate that ref against the entity path; `loadValidateRound`
  and `verifyRoundArtifacts` use the canonical review home as the containment boundary.
- `skills/feedback-rejection-flow/SKILL.md` already passes an entity operand and then
  calls `state commit`; its command text needs no change.

### Validation and historical path resolution

- `internal/gates/io.go` and `internal/gates/operation.go` expose one secure historical
  local-ref reader used by gate authority checks and the explicit diagnostic. It
  distinguishes normalized local refs from opaque provider refs and enforces the
  canonical filesystem boundary without changing stored bytes.
- `internal/status/validate.go:gateValidationDiagnostics` becomes a retained-evidence
  pass that does not return on `ErrNoGateRecord`. It validates gates and
  `review-round` independently, runs for active and archived entities, keeps unknown
  application warnings active-only, and emits invalid local refs as errors.
- Historical gate refs continue to resolve entity-directory-relative at the exact
  readers in `internal/gates/prepare.go` (replay and prepared-room classification),
  `internal/gates/io.go:validateRetainedAuthorityExcept`, and
  `internal/gates/operation.go` (`Withdraw`, `boundBriefingPath`). These consumers use
  the shared safe reader but do not reinterpret `./<slug>/review/...` or
  `./review/...`.
- New round refs use `relativeRoomRef(entityPath, room)`, the same writer rule as gate
  preparation. No fallback search, existence-dependent rebasing, ref normalization,
  or closed-attempt rewrite is allowed.

### Observable semantic boundary

- **Command grammar:** unchanged. No new command or flag.
- **Stored formats:** unchanged keys and schemas. Flat round `room-ref` uses the
  already-supported `./<slug>/review/...` spelling; folder bytes remain
  `./review/...`. Existing `git-root://<main|state>/<commit>/<path>` Artifact identity
  remains unchanged.
- **Authority:** unchanged. Closed gate attempts stay frozen; round Resolutions remain
  advisory; only the First Officer mutates entity state.
- **Runtime behavior:** flat `gate record --round` changes from refusal to publication
  and exact replay when the shared form policy allows it. Newly misfiled declared-folder
  flats, active collisions, unsafe entity/home shapes, and invalid local refs fail
  byte-clean. Explicit validation changes invalid retained local evidence from an
  exit-0 warning or omission to an exit-nonzero error. Plain reads, opaque refs,
  grandfathered flats, undeclared flats, folder rounds, and commit output remain
  compatible.
- **Excluded semantics:** no entity move, migration command, automatic ref rewrite,
  room retention change, request-digest ordering change, or broader Git staging.

## Expected surface

Estimate net LOC change: **+420, across 16 files**. Expected movement is approximately
580 insertions and 160 deletions. Tolerance is **net +/-140 LOC and +/-2 files**.
No ideation baseline has yet been approved, so this is the candidate baseline rather
than a correction against +155/8. Crossing either tolerance bound, narrowing an AC, or
changing a semantic outside the declaration above requires a design reset before
another correction round.

Expected files:

1. `internal/entityhome/home.go` — canonical form, review-home, and commit-unit resolver.
2. `internal/entityhome/home_test.go` — collision, file type, symlink, containment, creation, and tracked-deletion matrix.
3. `internal/gates/prepare.go` — shared form policy and canonical home/ref derivation.
4. `internal/gates/io.go` — secure retained local-ref reads for authority and diagnostics.
5. `internal/gates/operation.go` — round policy preflight and shared gate-room reader.
6. `internal/gates/round.go` — flat room derivation, pointer validation, and artifact boundary.
7. `internal/gates/prepare_test.go` — prepare half of the form-policy and safe-home matrix.
8. `internal/gates/round_test.go` — round half, flat/folder replay, policy, collision, and containment.
9. `internal/cli/state_sync.go` — canonical active commit units with tracked-deletion mode.
10. `internal/cli/state_commit_test.go` — combined-room two-clone proof and commit adversaries.
11. `internal/dispatch/stamp.go` — canonical resolver parity for stamped state publication.
12. `internal/dispatch/build_stamp_test.go` — stamp collision/symlink/non-directory parity.
13. `internal/status/validate.go` — independent active/archive gate/round validation.
14. `internal/status/hybrid_flat_rooms_warn_test.go` — local/opaque/history validation and removal of conversion warnings.
15. `docs/specs/gate-resolution-frontmatter-contract.md` — normative identity, policy, ref, and validation contract.
16. `docs/site/reference/command-reference.md` — user-visible flat rounds and failing explicit validation.

## Out of scope

- The folder-form `state commit` boundary, which is `rd`'s.
- Gate-room retention size, which is `9t`'s.
- The request-digest ordering trap, which is a separate defect in the same command family.

## Acceptance criteria

1. **AC-1 — A supported flat commit is self-contained on the state remote.** After
   Host A prepares one gate room and records one advisory round for the same flat
   entity, one supported `spacedock state commit <slug>` publishes both. A fresh Host B
   clone contains the flat Markdown and both complete sibling rooms. Both refs resolve,
   both Briefing digests match, gate-prepare replay and round replay change zero bytes,
   and the publishing commit contains no dirty sibling. The two-clone CLI test measures
   the two retained rooms, two digests, one commit, and zero replay delta; removing the
   companion from the commit unit or publishing only one room makes it fail.
2. **AC-2 — Form policy and advisory rounds agree for every supported entity form.**
   Prepare and round both refuse a newly misfiled flat entity in a declared-folder
   workflow before lock or mutation. Both allow a grandfathered declared-folder flat,
   an undeclared valid flat, and a valid folder entity. Each allowed form publishes one
   immutable `review/<stage>/round-<cycle>` below its canonical review home, stores the
   shape-correct relative ref, and exact replay changes zero bytes. A shared four-row
   table drives both commands; omitting either policy call, restoring the unconditional
   flat refusal, or keeping a fixed `./review` ref makes a row fail.
3. **AC-3 — Historical refs retain their meaning without migration.** Existing flat
   `./<slug>/review/...` and folder `./review/...` gate bindings still resolve to the
   same canonical Briefing bytes, while pre-existing folder `review-round` pointers
   replay unchanged. Opaque provider refs retain their non-filesystem meaning, and
   `git-root://` Artifact identities remain byte-for-byte unchanged. Historical
   fixtures exercise the actual gate and round readers and compare refs, resolved
   bytes, and digests; changing either base, dereferencing opaque refs, or rewriting
   Artifact identity makes a fixture fail.
4. **AC-4 — Explicit validation covers every retained local evidence scope.** Gates and
   `review-round` validate independently for active and archived entities, including a
   round-only entity with no `gates` record. Valid normalized `./...` refs pass. Missing,
   absolute, traversal, backslash, non-normalized, escaping, symlinked, and wrong-type
   local refs make `status --validate` exit nonzero with entity path and offending ref.
   Opaque provider refs pass without filesystem access, and plain `status` stays
   readable. A scope-by-evidence CLI matrix asserts exit codes and on-disk targets;
   retaining the early return, active-only skip, old warning, or blind join makes a row
   fail.
5. **AC-5 — Canonical entity ownership is unambiguous, safe, and path-scoped.** Prepare,
   round, `state commit`, and dispatch stamp resolve the same entity form and review
   home. They refuse active flat-plus-folder collision; symlinked or wrong-type entity,
   home, or room components; and boundary escape before mutation or commit. Controlled
   creation can create only an absent valid flat companion. Git-proven tracked flat and
   folder deletions still commit their literal units. Dirty siblings remain untouched,
   and flat round Artifacts cannot escape `<state-root>/<slug>/` or resolve to mutable
   `<slug>.md`. Resolver, real-Git, round, and stamp-parity adversarial tests compare
   HEAD, index, worktree, external targets, and origin before/after; restoring
   `os.Stat` discovery, form preference, unscoped staging, or state/stamp divergence
   makes at least one test fail.

## Test plan

Use the existing `twoHostStateWorkflow` real-Git harness in
`internal/cli/state_commit_test.go`; do not add another clone framework. Host A prepares
a gate room and records an advisory round on the same flat entity, then invokes one
`state commit`. A fresh Host B clone integrates the state branch, opens and digest-checks
both retained Briefings, exactly replays gate preparation from the committed selected
Git objects, and exactly replays the round from its retained two-file room. Assert one
publishing commit contains `<slug>.md` and both `<slug>/review/...` trees, while dirty
sibling state remains untracked. Cost: high, one combined two-clone CLI journey.

Write `internal/entityhome` tests before its production resolver. The table covers flat,
folder, active collision, symlinked entity, symlinked home, home-as-file, index-as-dir,
outside root, controlled absent companion, and Git-proven flat/folder deletion. Each
refusal compares the entity, index, worktree, external sentinel, HEAD, and origin as
applicable. Then drive the same collision/symlink/non-directory rows through state
commit and dispatch stamp; both must return the same refusal class and leave the same
state unchanged. Cost: high, unit plus real-Git parity fixtures.

Replace `TestRoundRequiresFolderFormWithoutCrossEntityCollision` with a shared policy
table used by prepare and round: declared-folder new flat refuses; grandfathered flat,
undeclared flat, and folder succeed. For successful rows assert canonical review home,
exact stored ref, whole-tree replay no-op, divergent replay refusal, distinct homes for
two flat slugs, workflow-stage validation, and Artifact containment. Retain the
existing folder fixture as a historical byte baseline. Cost: medium.

Expand the historical validation fixture into active/archive by gate/round-only rows.
Include valid flat gate, folder gate, folder round, flat round, no-gates round, opaque
provider ref, missing room, absolute ref, traversal, backslash, non-normalized ref,
symlink escape, and wrong file type. Drive the native CLI rather than searching prose:
valid local history and opaque refs exit 0; every invalid local row exits 1 with stable
entity/ref evidence; ordinary status remains exit 0. Cost: high, fixture-backed CLI
behavior tests.

Retain focused regression coverage for prepared-room authority and historical gate
readers in `prepare.go`, `io.go`, and `operation.go`. These tests must exercise resolved
bytes and digests, not only returned path strings. No committed prose-grep test is part
of the proof.

Run the focused tests first, then the repository-required verification:

```bash
go test ./internal/gates ./internal/status ./internal/cli
go test ./...
go test ./... -race
gofmt -w ./cmd ./internal
```

### Documentation diff proposed at ideation

- In `docs/specs/gate-resolution-frontmatter-contract.md`, replace the conversion
  prescription and all conversion-warning language with the canonical review-home and
  form-policy rules. State that direct flat support needs no move or frozen-ref rewrite.
  Add the active-collision, regular-file, symlink/non-directory, controlled-creation,
  tracked-deletion, and shared state/stamp boundaries. Replace "Round recording requires
  a folder-form entity" with the flat/folder derived-room and shape-specific ref table.
  Define local versus opaque refs and the active/archive, gate/round-only failing
  explicit validation contract. Preserve `git-root://` Artifact identity text exactly.
- In `docs/site/reference/command-reference.md`, change the `status --validate` row
  from conversion warnings to failing invalid-local-ref validation across active and
  archived gate/round evidence, with opaque refs excluded from dereference. Change the
  round row to support policy-valid flat and folder entities below their canonical
  homes; name the two ref spellings and declared-folder grandfathering rule. Update the
  state-commit and stamped-dispatch descriptions to name collision/symlink refusal and
  tracked-deletion preservation. Remove every instruction to convert a flat entity.

## Stage Report: ideation

- DONE: Select one backward-compatible review-home and room-ref design that supports flat and folder tickets without migration or rewriting frozen gate bindings.
  Selected the form-independent `<state-root>/<slug>` review home with historical entity-directory-relative refs and no migration or rewrite.
- DONE: Enumerate the exact state-commit, advisory-round, validation, and path-resolution consumers; declare semantic boundaries and expected net LOC/files with tolerance.
  The consumer inventory names every writer/reader seam; expected surface is net +155 LOC across 8 files, tolerance +/-75 LOC and +/-2 files.
- DONE: Define falsifiable historical-fixture and two-clone tests proving flat sibling rooms commit and replay, flat rounds record, existing refs retain meaning, and broken refs fail validation.
  AC-1 through AC-5 name the exercised outcome and the change that makes each proof fail; the test plan reuses the existing two-clone harness.

### Summary

Ideation chooses direct flat-ticket support instead of migration: one physical review
home, unchanged historical relative-ref semantics, and shape-aware advisory rounds.
It scopes implementation to round derivation and failing explicit validation, backed by
fresh-clone durability, historical fixtures, replay, and adversarial containment tests.

## Stage Report: ideation (cycle 2)

- DONE: Select one backward-compatible review-home and room-ref design that supports flat and folder tickets without migration or rewriting frozen gate bindings.
  Retained direct flat support and historical entity-relative refs; added canonical identity resolution and shared declared-form policy before every mutation.
- DONE: Enumerate the exact state-commit, advisory-round, validation, and path-resolution consumers; declare semantic boundaries and expected net LOC/files with tolerance.
  Revised baseline is exactly net +420 LOC across 16 files, about 580 insertions/160 deletions, with +/-140 net LOC and +/-2 files tolerance.
- DONE: Define falsifiable historical-fixture and two-clone tests proving flat sibling rooms commit and replay, flat rounds record, existing refs retain meaning, and broken refs fail validation.
  The combined two-clone journey publishes one prepared gate and one round with one commit; Host B verifies both digests and zero-byte replays.
- DONE: AC-1 fresh-clone durability evidence.
  Fails if either room is absent, either digest differs, replay changes bytes, more than one publication is needed, or sibling dirt enters the commit.
- DONE: AC-2 shared form-policy and round evidence.
  Fails if prepare and round disagree on newly misfiled, grandfathered, undeclared-flat, or folder rows, or write a shape-wrong ref.
- DONE: AC-3 historical compatibility evidence.
  Fails if either historical local ref changes target, opaque refs are dereferenced, closed bytes change, or git-root Artifact identity is rewritten.
- DONE: AC-4 retained-validation evidence.
  Fails if gates or round-only history escapes validation in active/archive scope, an invalid local ref exits zero, or an opaque ref touches the filesystem.
- DONE: AC-5 canonical ownership evidence.
  Fails if any consumer prefers an active form, follows a symlink, accepts a wrong type or escape, loses tracked deletion, changes an external target, or stages a sibling.

### Summary

Cycle 2 incorporates all four authorized independent-review findings without migration
or frozen-ref rewrites. The revised design shares identity and form policy across every
writer, validates all retained evidence scopes, and provides adversarial parity proof.
