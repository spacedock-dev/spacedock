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
        - id: gate:b9pjkz3rv0svx9rt63kw8yg7:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:b9pjkz3rv0svx9rt63kw8yg7-ideation-1
              briefing:
                id: briefing:b9pjkz3rv0svx9rt63kw8yg7:ideation:attempt-1:revision-1
                digest: sha256:803ff34b274f373d3fdf25c3c04e356d2ebc25eb7048035bad8fe8e362aadf97
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:b9pjkz3rv0svx9rt63kw8yg7:ideation:1
                briefing: briefing:b9pjkz3rv0svx9rt63kw8yg7:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T21:48:05.669602Z"
                decision: approve
              application:
                target-stage: implementation
                state: pending
---

Make a flat entity's gate rooms durable through the supported command and let it record
advisory rounds without moving the entity or rewriting frozen gate history.

## Two defects that compose into a cul-de-sac

A flat entity (`<slug>.md`) may legitimately carry a sibling room directory (`<slug>/review/...`). The shipped precedent is `source-build-compatibility-identity`, whose backlog briefing lives exactly that way. Both defects below apply to that shape.

### 1. The observed `state commit` defect is already fixed in current source

Observed on every gate room bound during the 2026-07-27/28 session — eight occurrences
across `79` and `cn`, at backlog, ideation and validation gates. At that time,
`spacedock state commit <slug>` committed `<slug>.md` and pushed while
`<slug>/review/<stage>/briefing-N/` remained untracked.

The consequence is not cosmetic. The committed frontmatter carries a `room-ref` pointing at a room that does not exist on the state remote, so a peer or a later session cannot retrieve the Briefing that justified a recorded decision. That is the "digests no committed tree reproduces" failure class the durable-decisions sprint was created to eliminate, reproduced through the supported path.

Every occurrence was worked around by hand:

```
git add -- <slug>/ && git commit -m "…" -- <slug>/ && git push origin HEAD:spacedock-state/dev
```

Current `main` now includes `flatEntityCommitPaths`: explicit state commit and implicit
gate sync include `<slug>.md` plus the present or tracked `<slug>/` companion, including
tracked deletions. `TestStateCommitFlatIncludesExactCompanionDirectoryAndTrackedDeletions`
passes. This ticket adds a fresh-clone regression proof but no state-sync production
change.

### 2. Converting to folder form silently breaks every prior closed gate

The obvious repair is to convert the entity to folder form, which also unlocks advisory round recording (the recorder requires `<slug>/index.md` and refuses flat entities before locking or writing). It cannot be done safely today.

`room-ref` resolves relative to the **entity's own directory**. Proven against `prepare-provider-neutral-gate-room`, a folder-form entity whose stored `./review/backlog/briefing-1` maps to `<slug>/review/backlog/briefing-1`. For a flat entity the entity directory is the state root, so a correct stored ref reads `./<slug>/review/...`. After `git mv <slug>.md <slug>/index.md` that same ref resolves to `<slug>/<slug>/review/...` — one level too deep, and pointing at nothing.

Nothing catches it. `gate validate <entity>` exits 0 and reports gate, attempt, briefing, resolution and decision, because it reads frontmatter without dereferencing `room-ref`. Observed directly: after converting `79`, validate stayed green while the backlog gate's room pointer was broken.

There is no supported repair once broken. `gate record --briefing` will not rewrite a closed attempt, and hand-editing `gates:` is forbidden by the lifecycle skill. The conversion on `79` was reverted for exactly this reason.

### Why the two compose badly

Current source has repaired the durability half but still requires folder form for round
recording. A flat entity therefore remains stuck with hand-authored rounds, while moving
it would change the meaning of its frozen relative refs. Canonical `@review/...` refs
remove the form dependency for all new rooms without rewriting history.

## What a fix needs to decide

The selected design makes migration unnecessary. A ticket's **review home** is always
`<state-root>/<slug>/`, independent of whether its mutable entity is the flat
`<state-root>/<slug>.md` or folder-form `<state-root>/<slug>/index.md`. Gate Briefings
and advisory-round rooms live below that home at `review/<stage>/...`.

The selected stored namespace is `@review/<path>`. It always resolves from the ticket's
review home, so `@review/ideation/briefing-1` names
`<state-root>/<slug>/review/ideation/briefing-1` for both entity forms. The resolver is
a small gate-package helper, not a new entity identity subsystem. It accepts only the
reserved `@review/` prefix plus a nonempty normalized slash path with no backslash,
absolute segment, `.` segment, or `..` segment. The resolved path must stay below the
review root.

Every newly prepared gate room and newly recorded round stores `@review/...`. Existing
flat `./<slug>/review/...`, folder `./review/...`, and other legacy refs retain their
existing entity-directory-relative reader path. No command rewrites a retained binding,
and exact replay of a legacy open gate or folder round leaves its stored ref unchanged.
There is no migration or flat-to-folder move in this task. The existing declared-folder
preparation policy remains; round recording reuses its small grandfathering check when
the workflow declares folder form.

`gate record --round` removes only its unconditional folder-form refusal. It derives
the room with the same review-home function already used by prepare, reuses prepare's
room-ancestry check, and stores `@review/<stage>/round-<cycle>`. Folder and flat rounds
therefore share one physical and stored shape. `readRoundPointerData` accepts the new
canonical ref and the frozen legacy folder `./review/...` ref for exact replay.

`resolveRound` returns that derived review home with the round room instead of making
later validation reconstruct a root from `filepath.Dir(entityPath)`. The resulting
round location carries `reviewHome` into `loadValidateRound` and
`verifyRoundArtifacts` during publication and through `ValidateRoundFile` during
replay. A relative Artifact is valid only when its resolved regular path stays inside
`<state-root>/<slug>/`. A relative Artifact that reaches a sibling entity or any other
path outside that home refuses byte-clean. The mutable entity also refuses: folder
`<slug>/index.md` remains covered by the explicit mutable-entity check, while flat
`<slug>.md` is outside the shared home and fails containment.

Every gate reader routes `@review/...` through the new helper before its existing
Briefing, digest, and authority checks. All non-`@review/` values take the unchanged
legacy path. The existing `status --validate` missing-room diagnostic uses the same
helper so a valid new ref does not produce a false warning. Broader validation changes
are outside this task.

### Why this design

- `@review/...` serves AC-1 through AC-3. Keeping new refs entity-relative is
  insufficient because the same physical home needs two stored spellings and a manual
  entity move changes the base. Rewriting legacy refs is insufficient because closed
  bindings are frozen.
- The small dual-mode resolver serves AC-2 and AC-3. Replacing every ref with a new
  general path model is unnecessary; dispatch only the reserved prefix, then preserve
  the existing legacy path verbatim.
- Reusing the current flat commit unit serves AC-1. Reworking state sync is unnecessary
  because source and real-Git tests already include a flat tracked companion and its
  tracked deletions.

No new spike is needed. Current-code evidence already proves the physical and Git
mechanisms:
`TestPrepareCreatesOneFileRecorderRoomForFolderAndFlatEntities` proves the common
physical review home for both forms;
`TestStateCommitFlatIncludesExactCompanionDirectoryAndTrackedDeletions` proves the
flat commit unit and tracked deletion; `flatEntityCommitPaths` is already shared by
explicit state commit and implicit gate sync; and
`TestPrepareRejectsSymlinkedFlatCompanionWithoutChangingBytes` proves prepare's room
ancestry check can be reused by round. `TestRoundRequiresFolderFormWithoutCrossEntityCollision`
pins the one guard this task removes. `round.go` also supplies the direct defect
evidence: both publication and `ValidateRoundFile` currently pass
`filepath.Dir(entityPath)` as the Artifact trust root, which is the ticket home for
folder `index.md` but the whole state root for flat `<slug>.md`. These focused
baselines passed on 2026-08-27.

## Consumer inventory and boundaries

### Review-home and ref resolution

- `internal/gates/prepare.go` already owns `entitySlug`, `preparedRoomPath`, room
  ancestry, and legacy replay. It adds the reserved-prefix parser, canonical writer,
  and one resolver that maps `@review/...` through the common physical review home.
- The same resolver returns the unchanged entity-directory-relative path for every
  legacy ref. `prepare.go`, `io.go`, and `operation.go` replace their five direct joins
  with this helper; their Briefing, request, digest, and authority logic stays intact.
- `internal/status/validate.go:unresolvedRoomRefs` uses the exported resolver only to
  understand the new namespace. It keeps its current scope and warning semantics.

### State commit and publication

- `internal/cli/state_sync.go` already resolves a flat active commit unit as
  `<slug>.md` plus `<slug>/` when present or tracked. `commitEntityPathsScoped` already
  stages and commits only literal paths, including tracked deletions.
- `internal/cli/gate_ceremony.go` already consumes that seam for implicit gate sync.
  No production change is required in either file.
- `internal/cli/state_commit_test.go` extends the existing two-host proof so Host B
  retrieves an `@review`-bound flat round after the supported command alone.

### Advisory rounds

- `internal/cli/cli.go:newGateCommand` routes `gate record --round` and performs the
  immediate readback; its grammar and output stay unchanged.
- `internal/gates/operation.go:RecordSemanticSummary` removes the unconditional
  folder-only guard and reuses the existing declared-folder grandfathering check.
- `internal/gates/round.go`: `resolveRound` derives the review home and room;
  its result carries `reviewHome` as the Artifact trust root into
  `loadValidateRound`/`verifyRoundArtifacts` for `recordRoundLockedWith` publication
  and `ValidateRoundFile` replay. `recordRoundLockedWith` writes `@review/...`;
  `readRoundPointerData` accepts canonical and frozen legacy folder refs. The recorder
  reuses prepare's existing ancestry check before publication.
- `skills/feedback-rejection-flow/SKILL.md` already passes an entity operand and then
  calls `state commit`; its command text needs no change.

### Validation and historical path resolution

- Historical gate refs continue to resolve entity-directory-relative in prepare replay,
  `validateRetainedAuthorityExcept`, `Withdraw`, prepared-room classification, and
  `boundBriefingPath`. Frozen legacy fixtures exercise resolved bytes and digests.
- New gate and round refs use the same canonical writer. The parser rejects malformed
  values inside the reserved namespace; it does not reinterpret other legacy strings.
- Full archive/round-only validation, opaque-provider classification, and new global
  local-ref hardening remain follow-up work. This ticket only prevents the existing
  unresolved-room warning from misreading a valid `@review/...` gate ref.

### Observable semantic boundary

- **Command grammar:** unchanged. No new command or flag.
- **Stored formats:** the only addition is reserved `room-ref` syntax
  `@review/<normalized-path>` for newly written rooms. Existing keys, legacy refs, and
  `git-root://<main|state>/<commit>/<path>` Artifact identity remain unchanged.
- **Authority:** unchanged. Closed gate attempts stay frozen; round Resolutions remain
  advisory; only the First Officer mutates entity state.
- **Runtime behavior:** flat `gate record --round` changes from refusal to publication
  and exact replay through the shared review home. New gate and round writes use one
  canonical ref spelling for both forms. Legacy replay, folder physical layout, state
  commit output, validation severity, and ordinary status remain compatible. Round
  Artifact containment is consistently rooted at that ticket home, preventing a flat
  round from trusting a sibling entity or its mutable flat entity file.
- **Excluded semantics:** no entity move, migration command, automatic ref rewrite,
  room retention change, request-digest ordering change, broader Git staging, dispatch
  stamp change, general entity resolver, or validation expansion.

### Separable risks and existing guards

- Active flat-plus-folder collision is already an explicit `status --validate` error,
  and supported workflows do not operate in that shape. `state commit` preferring the
  folder path in an invalid collision is a separable state-identity hardening task.
- Gate prepare already rejects a symlinked flat companion byte-clean through
  `validatePreparedRoomAncestry`; round reuses that guard. Broader entity-leaf and
  dispatch-stamp symlink parity are separable hardening, not required to publish a
  valid flat round.
- `state commit` already includes a present or tracked flat companion and its tracked
  deletions through literal Git pathspecs. No state-sync repair remains in current
  source; only a fresh-clone regression proof belongs here.
- Archive coverage, round-only status validation, and opaque-provider classification
  are defects or compatibility work independent of resolving the reserved namespace.
  The dispatcher leaves every non-`@review/` ref on its current reader path, so this
  task neither fixes nor regresses those cases.
- Dispatch stamp already includes a present flat companion directory. Its broader
  tracked-deletion and invalid-shape parity do not block the supported explicit
  `state commit` outcome measured by this task.

## Expected surface

Estimate net LOC change: **+180, across 12 files**. Expected movement is approximately
270 insertions and 90 deletions. Tolerance is **net +/-70 LOC and +/-2 files**. The
estimate is test-heavy: five production files receive a small helper or call-site
replacement, five test files add compatibility and behavior rows, and two docs change.
No state-sync, dispatch, archive, or new-package production work is included. Crossing
either tolerance bound, narrowing an AC, or changing a semantic outside the declaration
above requires a design reset.

Expected files:

1. `internal/gates/prepare.go` — review-home helper, reserved ref parser/writer, and prepare readers.
2. `internal/gates/io.go` — route retained request-backed room reads through the helper.
3. `internal/gates/operation.go` — remove the folder-only round guard and route gate-room readers.
4. `internal/gates/round.go` — flat/folder physical room, canonical pointer, shared-home Artifact trust root, and frozen legacy replay.
5. `internal/status/validate.go` — resolve `@review/...` in the existing missing-room diagnostic and remove obsolete conversion advice.
6. `internal/gates/prepare_test.go` — canonical new gate refs and legacy open replay.
7. `internal/gates/prepare_shape_test.go` — reserved syntax and frozen legacy reader table.
8. `internal/gates/round_test.go` — flat/folder publication, Artifact containment, legacy folder replay, policy, ancestry, and divergence.
9. `internal/cli/state_commit_test.go` — two-host fresh-clone durability proof for a flat round room.
10. `internal/status/hybrid_flat_rooms_warn_test.go` — canonical-ref resolution and removal of conversion warnings.
11. `docs/specs/gate-resolution-frontmatter-contract.md` — normative namespace and direct flat-round contract.
12. `docs/site/reference/command-reference.md` — user-visible canonical refs, flat rounds, and supported commit boundary.

## Out of scope

- The folder-form `state commit` boundary, which is `rd`'s.
- Gate-room retention size, which is `9t`'s.
- The request-digest ordering trap, which is a separate defect in the same command family.
- General active-form collision handling and dispatch-stamp parity.
- Archive/round-only validation and opaque-provider ref classification.
- Entity-leaf and repository-wide symlink hardening beyond the reused room-ancestry guard.

## Acceptance criteria

1. **AC-1 — One supported state commit makes a flat round durable.** After Host A
   records a flat advisory round and runs one `spacedock state commit <slug>`, a fresh
   Host B clone contains `<slug>.md`, both immutable round files under
   `<slug>/review/...`, and a resolving `@review/...` pointer whose Briefing digest
   matches. Exact replay changes zero bytes, and dirty sibling state remains absent from
   the commit. The existing two-host harness measures remote paths, digest, commit
   scope, and replay delta; removing the companion from `flatEntityCommitPaths` makes
   it fail.
2. **AC-2 — New gate and round refs have one canonical meaning for both forms.** Every
   newly prepared gate room and recorded round stores normalized `@review/...`, and
   that ref resolves to `<state-root>/<slug>/review/...` for flat and folder entities.
   The table exercises both writers and both readers; restoring entity-relative output
   or resolving from the entity directory makes at least one form fail.
3. **AC-3 — Frozen historical refs retain their meaning without migration.** Existing flat
   `./<slug>/review/...` and folder `./review/...` gate bindings still resolve to the
   same canonical Briefing bytes, while pre-existing folder `review-round` pointers
   replay with their stored refs unchanged. Existing `git-root://` Artifact identities
   remain byte-for-byte unchanged. Historical fixtures drive the actual gate and round
   readers and compare stored refs, resolved bytes, and digests; routing legacy refs
   through the new namespace or rewriting a binding makes them fail.
4. **AC-4 — Flat and folder advisory rounds publish and replay identically.** A
   policy-valid flat entity and folder entity each publish exactly one immutable
   `review/<stage>/round-<cycle>` room, retain the supplied Briefing/log bytes, and
   preserve status, gates, and body bytes. For both forms, an in-home relative Artifact
   validates; a relative Artifact outside `<state-root>/<slug>/` and the mutable entity
   refuse byte-clean during publication and `ValidateRoundFile` replay. Exact replay is
   a whole-tree no-op and divergence fails byte-clean. The form table fails if the
   folder-only guard returns, the flat physical home or Artifact trust root is wrong,
   an invalid Artifact leaves residue, or replay changes any byte.
5. **AC-5 — The reserved namespace fails closed without broadening legacy semantics.**
   Empty, absolute, dot-segment, traversal, and backslash `@review` values are rejected;
   a valid value cannot escape the review root. Round publication reuses the existing
   room-ancestry guard, uses the derived review home rather than the entity directory as
   its Artifact trust root, and the existing declared-folder preparation policy also
   gates new flat rounds. All other ref strings take the unchanged legacy reader.
   Parser, Artifact, and policy tables name the rejected value and compare the
   entity/room tree before and after; accepting a malformed reserved ref, trusting an
   out-of-home or mutable-entity Artifact, or reclassifying a legacy ref makes them
   fail.

## Test plan

Use the existing `twoHostStateWorkflow` real-Git harness in
`internal/cli/state_commit_test.go`; do not add another clone framework. Host A records
one flat round, leaves dirty sibling state, and runs one supported `state commit`. A
fresh Host B clone integrates the state branch, resolves `@review/...`, checks both
retained files and the Briefing digest, and exactly replays the round. Assert the commit
contains only the flat entity and companion paths. Cost: medium, one focused extension
to the existing real-Git harness.

Add a gate-package table for `@review/...`: canonical gate and round refs for flat and
folder entities; empty suffix, absolute-like suffix, dot, traversal, and backslash
refusals; and physical resolution below the shared review root. The test uses literal
expected paths rather than the production writer as its oracle. Cost: small.

Extend the existing prepare and binding-shape fixtures with frozen legacy flat and
folder refs. Exercise prepare replay, retained-authority validation, withdraw/record
room reads, and canonical Briefing digest comparison. Assert the stored legacy ref is
unchanged. Cost: medium; it reuses current fixtures.

Replace the unconditional-flat round test with a flat/folder table. Assert canonical
pointer bytes, exact retained Briefing/log bytes, body/gates/status preservation,
whole-tree exact replay, divergent byte-clean refusal, declared-folder grandfathering,
and reuse of prepare's ancestry guard. Run these Artifact rows for both forms and for
both publication and `ValidateRoundFile` replay:

| Artifact row | Flat result | Folder result |
| --- | --- | --- |
| Relative regular file inside `<state-root>/<slug>/` | validates | validates |
| Relative path resolving outside `<state-root>/<slug>/` | refuses byte-clean | refuses byte-clean |
| Mutable entity (`<slug>.md` or `<slug>/index.md`) | refuses byte-clean | refuses byte-clean |

Use literal paths so the test cannot reproduce the production root-selection bug in
its oracle. Retain a frozen legacy folder pointer fixture. Cost: medium.

Update the existing status warning fixture only enough to prove that `@review/...`
resolves and that the obsolete blanket conversion warning is absent. Do not add the
separate archive/opaque/general-local-ref matrix to this ticket.

Run the focused tests first, then the repository-required verification:

```bash
go test ./internal/gates ./internal/status ./internal/cli
go test ./...
go test ./... -race
gofmt -w ./cmd ./internal
```

### Documentation diff proposed at ideation

- In `docs/specs/gate-resolution-frontmatter-contract.md`, replace the conversion
  prescription with the `@review/...` namespace. State its physical resolution for both
  forms and the frozen legacy fallback. Replace the folder-only round paragraph with
  direct flat/folder publication and exact replay. Specify the shared review home as
  the relative Artifact trust root for both forms, including out-of-home and mutable
  entity refusal. Preserve `git-root://` Artifact identity text and the existing
  supported state-commit unit.
- In `docs/site/reference/command-reference.md`, change the `status --validate` row
  so canonical refs do not produce a missing-room or blanket conversion warning. Change
  gate prepare and round rows to name `@review/...`, both physical forms, frozen legacy
  refs, and no migration. Keep the current state-commit description because it already
  documents the flat companion commit unit.

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

## Stage Report: ideation (cycle 3)

- DONE: Select one backward-compatible review-home and room-ref design that supports flat and folder tickets without migration or rewriting frozen gate bindings.
  Captain correction selects canonical `@review/...` for new rooms; all legacy refs remain frozen on their current reader path.
- DONE: Enumerate the exact state-commit, advisory-round, validation, and path-resolution consumers; declare semantic boundaries and expected net LOC/files with tolerance.
  Revised baseline is exactly net +180 LOC across 12 files, about 270 insertions/90 deletions, with +/-70 net LOC and +/-2 files tolerance.
- DONE: Define falsifiable historical-fixture and two-clone tests proving flat sibling rooms commit and replay, flat rounds record, existing refs retain meaning, and broken refs fail validation.
  Focused tables prove canonical writes, frozen legacy reads, flat/folder round replay, malformed reserved-ref refusal, and Host-B retrieval after state commit.
- DONE: AC-1 supported state-commit durability evidence.
  Existing flat commit-unit code is unchanged; fresh Host B must retrieve both round files and matching digest after one Host-A commit.
- DONE: AC-2 canonical namespace evidence.
  Literal expected paths fail if either entity form writes or resolves `@review/...` from the entity directory.
- DONE: AC-3 frozen compatibility evidence.
  Actual readers must return the same bytes/digests and preserve stored flat, folder, and legacy round refs on replay.
- DONE: AC-4 flat/folder round evidence.
  Whole-tree snapshots fail on folder-only refusal, wrong physical home, retained-byte change, replay mutation, or divergent-write residue.
- DONE: AC-5 reserved-prefix boundary evidence.
  Parser and policy tables fail on malformed acceptance, review-root escape, legacy reclassification, or mutation after refusal.

### Summary

Cycle 3 replaces the broad +420/16 architecture with the captain's canonical
`@review/...` namespace and the smallest consumer changes needed for flat rounds.
Already-shipped flat commit durability stays production-unchanged and receives a
fresh-clone regression proof; broader hardening is recorded as follow-up scope.

## Stage Report: ideation (cycle 4)

- DONE: Select one backward-compatible review-home and room-ref design that supports flat and folder tickets without migration or rewriting frozen gate bindings.
  Canonical `@review/...` and the frozen legacy reader remain unchanged; `resolveRound` now carries the already-derived shared review home through validation.
- DONE: Enumerate the exact state-commit, advisory-round, validation, and path-resolution consumers; declare semantic boundaries and expected net LOC/files with tolerance.
  The only new consumer obligation is in `internal/gates/round.go`; the estimate remains exactly net +180 LOC across 12 files, about 270 insertions/90 deletions, with +/-70 net LOC and +/-2 files tolerance.
- DONE: Define falsifiable historical-fixture and two-clone tests proving flat sibling rooms commit and replay, flat rounds record, existing refs retain meaning, and broken refs fail validation.
  The existing proof plan is retained and its flat/folder round table now runs Artifact-containment rows in publication and `ValidateRoundFile` replay.
- DONE: AC-1 supported state-commit durability evidence.
  The supported flat commit unit and fresh-clone proof are unchanged by this correction.
- DONE: AC-2 canonical namespace evidence.
  `@review/...` still resolves from the shared review home for both entity forms.
- DONE: AC-3 frozen compatibility evidence.
  Legacy flat, folder, and round refs remain on the unchanged reader path and are never rewritten.
- DONE: AC-4 round Artifact trust-root evidence.
  Flat and folder rows require an in-home regular Artifact to validate, while an out-of-home relative Artifact and the mutable entity refuse byte-clean in publication and replay.
- DONE: AC-5 closed-boundary evidence.
  Literal-path tests fail if `resolveRound` does not pass its derived review home into `loadValidateRound` and `verifyRoundArtifacts`, or if any rejected Artifact leaves residue.

### Summary

Cycle 4 corrects the one in-scope trust-root defect without broadening the ticket.
Round publication and replay use the resolved ticket review home instead of
`filepath.Dir(entityPath)`, so a flat round cannot bind a sibling entity or mutable
`<slug>.md`. Canonical refs, frozen legacy behavior, narrowed scope, and the +180/12
estimate with its existing tolerance remain unchanged.
