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

Stored `room-ref` keeps its existing meaning: it is relative to the entity file's own
directory. Therefore the same physical room is stored as `./<slug>/review/...` for a
flat entity and `./review/...` for a folder entity. Existing bindings need no rewrite,
the frozen-record rule remains intact, and a flat ticket can stay flat for its whole
lifecycle. There is no flat-to-folder migration verb in this task; a manual move that
changes the base directory remains unsupported. A workflow's `entity-form: folder`
filing rule also remains authoritative: this design supports grandfathered and
otherwise-valid flat tickets, but does not make a declared-folder workflow accept a
newly misfiled flat ticket.

`gate record --round` must derive both its room and pointer from that review home. It
must stop assuming `filepath.Dir(entityPath)` is the review home or that the pointer is
always `./review/...`. For a flat entity it publishes
`<slug>/review/<stage>/round-<cycle>` and stores
`./<slug>/review/<stage>/round-<cycle>`; folder behavior and bytes stay unchanged.
Artifact containment is rooted at the review home, not the state root, so a flat round
cannot cite a sibling entity. The mutable flat `<slug>.md` is outside that home and is
also explicitly forbidden as a round artifact.

Explicit `status --validate` becomes the shipped integrity check. A retained local gate
or round `room-ref` (`./...`) that does not resolve is an error and makes validation
exit nonzero, not a warning beside `VALID`. Opaque historical provider refs keep their
non-filesystem meaning and are not dereferenced. Ordinary status reads remain
unaffected. Existing gate mutation paths keep their stronger Briefing/digest checks.

### Why this design

- The review-home derivation serves AC-1 and AC-2. The simpler alternative, requiring
  folder form and moving every flat entity, is insufficient because it either strands
  existing flat tickets or rewrites frozen bindings.
- Shape-specific relative refs serve AC-2 and AC-3. A new state-root-relative ref
  scheme is insufficient because it changes the meaning of every existing folder ref;
  a single fixed `./review/...` ref is wrong for flat entities under the historical
  resolver.
- Validation dereference serves AC-4. Keeping the current warning is insufficient
  because exit 0 lets automation publish or accept an entity whose recorded decision
  room is absent.
- Literal flat commit units serve AC-1 and AC-5. A top-level `git add -A` is
  insufficient because it can attribute a concurrent sibling writer's state to the
  wrong ticket.

No new spike is needed. The exercised baseline is already in the repository:
`TestPrepareCreatesOneFileRecorderRoomForFolderAndFlatEntities` proves the common
physical review home and both historical ref spellings;
`TestStateCommitFlatIncludesExactCompanionDirectoryAndTrackedDeletions` proves the
flat commit unit; and `TestRoundRequiresFolderFormWithoutCrossEntityCollision` pins the
one deliberate refusal this task removes. All three passed together on 2026-08-27.

## Consumer inventory and boundaries

### State commit and publication

- `internal/cli/state_sync.go`: `runStateCommit`, `commitInlineEntity`,
  `syncActiveEntity`, `resolveEntityCommitTarget`, `flatEntityCommitPaths`, and
  `commitEntityPathsScoped` own the supported active commit unit. The flat unit is
  `<slug>.md` plus the literal `<slug>/` companion when present or tracked.
- `internal/cli/gate_ceremony.go` consumes `syncActiveEntity` for implicit split-root
  record/consume publication and therefore inherits the same flat unit.
- `internal/dispatch/stamp.go`: `commitAndPublishEntity` and its duplicate
  `entityCommitPaths` are the `dispatch build --stamp` consumer. It is audited for
  parity but needs no behavior change for a present review home.
- Archived publication and `internal/status/merge.go` archive moves remain out of
  scope; they do not publish a newly recorded active flat round.

### Advisory rounds

- `internal/cli/cli.go:newGateCommand` routes `gate record --round` and performs the
  immediate readback; its grammar and output stay unchanged.
- `internal/gates/operation.go:RecordSemantic` selects the round recorder; no authority
  or status behavior changes.
- `internal/gates/round.go`: `resolveRound` derives the review home and room;
  `recordRoundLockedWith` derives the shape-correct ref; `readRoundPointerData` and
  `ValidateRoundFile` validate that ref against the entity path; `loadValidateRound`
  and `verifyRoundArtifacts` use the review home as the containment boundary.
- `skills/feedback-rejection-flow/SKILL.md` already passes an entity operand and then
  calls `state commit`; its command text needs no change.

### Validation and historical path resolution

- `internal/status/validate.go:gateValidationDiagnostics` and
  `unresolvedRoomRefs` are the explicit validation consumer. The unresolved result
  moves from warnings to errors and includes the current `review-round` pointer.
- Historical gate refs continue to resolve entity-directory-relative at the exact
  readers in `internal/gates/prepare.go` (replay and prepared-room classification),
  `internal/gates/io.go:validateRetainedAuthorityExcept`, and
  `internal/gates/operation.go` (`Withdraw`, `boundBriefingPath`). These call sites are
  regression-tested, not reinterpreted.
- New round refs use `relativeRoomRef(entityPath, room)`, the same writer rule as gate
  preparation. No fallback search, existence-dependent rebasing, ref normalization,
  or closed-attempt rewrite is allowed.

### Observable semantic boundary

- **Command grammar:** unchanged. No new command or flag.
- **Stored formats:** unchanged keys and schemas. Flat round `room-ref` uses the
  already-supported `./<slug>/review/...` spelling; folder bytes remain
  `./review/...`.
- **Authority:** unchanged. Closed gate attempts stay frozen; round Resolutions remain
  advisory; only the First Officer mutates entity state.
- **Runtime behavior:** flat `gate record --round` changes from refusal to publication
  and exact replay; explicit validation changes an unresolved retained room from
  exit-0 warning to exit-nonzero error. Plain reads, folder rounds, and commit output
  remain byte-compatible.
- **Excluded semantics:** no entity move, migration command, automatic ref rewrite,
  room retention change, request-digest ordering change, or broader Git staging.

## Expected surface

Estimate net LOC change: **+155, across 8 files**. Expected gross movement is about
215 insertions and 60 deletions. Tolerance is **net +/-75 LOC and +/-2 files**; crossing
either bound, or changing a semantic outside the declaration above, requires a design
reset before another correction round.

Expected files:

1. `internal/gates/prepare.go` — reuse/expose the entity review-home and relative-ref derivation.
2. `internal/gates/round.go` — flat room derivation, pointer validation, and artifact boundary.
3. `internal/gates/round_test.go` — flat/folder record, replay, collision, and containment fixtures.
4. `internal/status/validate.go` — fail explicit validation on broken gate/round refs.
5. `internal/status/hybrid_flat_rooms_warn_test.go` — historical refs and exit-code expectations.
6. `internal/cli/state_commit_test.go` — two-clone flat room durability proof.
7. `docs/specs/gate-resolution-frontmatter-contract.md` — normative review-home/ref contract.
8. `docs/site/reference/command-reference.md` — user-visible round and validation behavior.

## Out of scope

- The folder-form `state commit` boundary, which is `rd`'s.
- Gate-room retention size, which is `9t`'s.
- The request-digest ordering trap, which is a separate defect in the same command family.

## Acceptance criteria

1. **AC-1 — A supported flat commit is self-contained on the state remote.** After
   `spacedock state commit <slug>`, a fresh second clone contains the flat Markdown,
   every file in its sibling review room, and a `room-ref` that resolves to Briefing
   bytes whose digest matches the frozen binding. The two-clone CLI test measures this
   on-disk state and fails if the companion path is removed from the commit unit.
2. **AC-2 — Flat and folder tickets can record and replay advisory rounds.** Each form
   publishes exactly one immutable `review/<stage>/round-<cycle>` below its review home,
   stores the shape-correct relative ref, and exact replay changes zero bytes. A table
   fixture drives both forms; restoring the folder-only guard or fixed `./review` ref
   makes the flat row fail.
3. **AC-3 — Historical refs retain their meaning without migration.** Existing flat
   `./<slug>/review/...` and folder `./review/...` gate bindings still resolve to the
   same canonical Briefing bytes, while pre-existing folder `review-round` pointers
   replay unchanged. Historical fixtures exercise the actual readers and compare the
   resolved bytes/digests; changing the base for either form makes one fixture fail.
4. **AC-4 — Broken retained local refs fail the shipped validation surface.** Missing
   gate rooms and missing current round rooms make explicit `status --validate` exit
   nonzero with entity path and offending `./...` ref, while plain `status` remains
   readable and opaque provider refs remain accepted. CLI fixtures delete each room
   type and assert exit code and evidence; leaving the old warn-only path makes the
   test fail.
5. **AC-5 — Flat durability stays path-scoped and entity-confined.** Publishing a flat
   ticket never commits dirty sibling entities, and a flat round cannot resolve an
   artifact outside `<state-root>/<slug>/` or back to mutable `<slug>.md`. The real-Git
   test leaves sibling dirt behind, and adversarial round fixtures use sibling and
   mutable-entity paths that must fail byte-clean.

## Test plan

Use the existing `twoHostStateWorkflow` real-Git harness in
`internal/cli/state_commit_test.go`; do not add another clone framework. Host A prepares
or records a flat room and invokes only `state commit`. Host B integrates the state
branch, opens the frozen binding, validates its digest, and exactly replays the round.
Assert the originating commit contains only `<slug>.md` and `<slug>/review/...`, while
a dirty sibling remains untracked. Cost: medium, one additional two-clone CLI test.

Convert `TestRoundRequiresFolderFormWithoutCrossEntityCollision` into a flat/folder
table fixture. Assert distinct review homes for two flat slugs, exact stored refs,
whole-tree replay no-op, divergent replay refusal, workflow-stage validation, and
artifact containment. Retain the existing folder fixture as a historical byte baseline.
Cost: medium, package-level filesystem behavior tests.

Extend the historical validation fixture with valid flat gate, valid folder gate,
valid folder round, missing flat gate room, and missing round room rows. Drive the
native CLI rather than searching prose: valid history exits 0, either deletion exits 1,
and ordinary status remains exit 0. Cost: small, fixture-backed CLI behavior tests.

Run the focused tests first, then the repository-required verification:

```bash
go test ./internal/gates ./internal/status ./internal/cli
go test ./...
go test ./... -race
gofmt -w ./cmd ./internal
```

### Documentation diff proposed at ideation

- In `docs/specs/gate-resolution-frontmatter-contract.md`, replace the conversion
  prescription ("Only the folder-form ref is invariant ... converting one requires")
  with the review-home rule and state that flat entities remain supported without a
  move or frozen-ref rewrite. Replace "Round recording requires a folder-form entity"
  with the flat/folder derived-room and shape-specific ref table. Change unresolved
  room validation from a report/warning to a failing explicit check.
- In `docs/site/reference/command-reference.md`, change the `status --validate` row
  from "warns ... when a retained room-ref no longer resolves" to "fails explicit
  validation when a retained gate or round room-ref does not resolve." Change the
  round row from "For a folder-form entity ... Flat entities are refused" to "For a
  flat or folder entity, publish below its derived review home; flat refs are
  `./<slug>/review/...`, folder refs are `./review/...`." Keep the gate-prepare filing
  rule for workflows that explicitly declare `entity-form: folder`.

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
