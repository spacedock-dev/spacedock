---
id: gf0jvhj4y8vjhd6ww62vsb2q
title: Gate rooms make a flat entity a hybrid, and converting it breaks every later gate
status: ideation
source: "Issue #739 (captain, 2026-08-20): gate prepare fails after a flat-to-folder conversion because a retained room-ref resolves relative to the new entity folder and duplicates the slug."
started:
completed:
verdict:
score:
worktree:
issue: "#739"
gates:
    version: 1
    records:
        - id: gate:gf0jvhj4y8vjhd6ww62vsb2q:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:gf0jvhj4y8vjhd6ww62vsb2q-backlog-1
              briefing:
                id: briefing:gf0jvhj4y8vjhd6ww62vsb2q:backlog:attempt-1:revision-1
                digest: sha256:2539b1d74e3671d51fc174a7b96ef4a4368e4ad0807828c8a567b02bb113ec08
                request-digest: sha256:6ba250fdd78031eabba8a169306935a272bb9c6d734c72192e859adfd856c308
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:gf0jvhj4y8vjhd6ww62vsb2q:backlog:1
                briefing: briefing:gf0jvhj4y8vjhd6ww62vsb2q:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-20T17:23:13.463297Z"
                decision: approve
                reason: 'Captain: ''dispatch gf0j''. Shape the fix, weighing enforcement at gate prepare against a tolerant read, and account for the 10 existing hybrids.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:gf0jvhj4y8vjhd6ww62vsb2q:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:gf0jvhj4y8vjhd6ww62vsb2q-ideation-1
              briefing:
                id: briefing:gf0jvhj4y8vjhd6ww62vsb2q:ideation:attempt-1:revision-1
                digest: sha256:5d1488086867e2d3e03308bdf01c5b01bee2b84dbefbba97689a5dc5857b8df8
                request-digest: sha256:b0078af8601623d5a540a5af924d0b26b18eb6c1419be1df54e43dea8caf9780
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:gf0jvhj4y8vjhd6ww62vsb2q:ideation:1
                briefing: briefing:gf0jvhj4y8vjhd6ww62vsb2q:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-20T18:16:57.802436Z"
                decision: revise
                reason: 'Captain: send it back. The repository and first officer that actually hit issue #739 converted and repaired the entity by hand with no difficulty, which contradicts the premise the converter rests on. Reprice with grandfathering as the primary option: refuse minting a NEW hybrid, leave the existing closed set alone, and make it visible with the validator warning. The 9 hybrids resolve correctly today and break only on conversion, so not converting them is a real option the design never weighed. Target under 100 net LOC.'
---

Stop `gate prepare` from leaving an entity in a form whose retained references break on conversion, and unblock the entities already in that state.

## Problem

A retained `room-ref` is written relative to the entity file's directory. That directory changes meaning between the two entity forms, so the same room resolves to two different paths:

| form | directory | ref written |
|---|---|---|
| flat `<state>/<slug>.md` | `<state>` | `./<slug>/review/...` |
| folder `<state>/<slug>/index.md` | `<state>/<slug>` | `./review/...` |

Convert flat to folder and every retained ref resolves to `<state>/<slug>/<slug>/review/...`. Issue #739 reports the doubled slug in that path.

One old attempt blocks all later gates. `internal/gates/io.go:207` loops over every prior attempt that carries a `RequestDigest` and returns an error when any retained room is unreadable. A consumed historical attempt therefore makes every future gate unpreparable, which is what #739 observed with a green, captain-approved validation gate.

Five sites join a room-ref the same way: `internal/gates/io.go:207`, `internal/gates/application.go:240`, `internal/gates/operation.go:181`, `internal/gates/operation.go:560`, and the write side at `internal/gates/prepare.go:222` (`relativeRoomRef`, line 671).

The deeper fault is the form itself. `gate prepare` on a flat entity creates `<state>/<slug>/review/...`, so a `<slug>/` directory appears beside `<slug>.md`. Every flat entity that has ever prepared a gate is already this hybrid. In this checkout: 122 flat entities and 68 folder entities; 10 flat entities carry a `gates:` key; 9 of those are hybrids holding 33 retained room-refs across 64 files. The tenth, `repair-sonnet-live-flakes`, has an empty `gates:` key, no room, and nothing to migrate. The hybrid is one `git mv` away from folder form, and that move is exactly what breaks the refs.

Nothing warns. `status --validate` reports VALID on the hybrid, and still reports VALID after the conversion that breaks all 33 refs; the plain `status` table lists the entity normally. The only signal is `gate prepare` failing at the moment the first officer needs the gate — the green-then-blocked sequence #739 reported.

`gate record --round` already refuses this: `internal/gates/operation.go:108` returns "gate record --round requires folder-form entity <slug>/index.md because review artifacts accumulate beside the entity". `gate prepare` creates review artifacts beside the entity and enforces nothing. Two commands in one package disagree about the same rule.

## Spike results

Every claim below was exercised against a built binary, not read off the source. Scripts and fixtures are throwaway; they seed the implementation's first tests.

**#739 reproduces end to end (offline, real binary).** A flat entity at a `backlog` gate, prepared, approved, consumed, then converted with one `git mv`, then asked for its `ideation` gate:

```
Error: attempt gate-attempt:task-backlog-1 retained request.json: open
  <state>/task/task/review/backlog/briefing-1/request.json: no such file or directory
```

The doubled slug is `task/task`. Without the `git mv` the same script prepares the later gate successfully, so the conversion is the whole cause. One consumed historical attempt is enough: `io.go:207` loops every attempt carrying a `RequestDigest`, so the failure is permanent, not confined to the current gate.

**The same break happens on the real hybrids, and the migration repairs it.** Copying all 9 live hybrids into an isolated checkout and exercising the form-sensitive half of `validateRetainedAuthority`:

| state | refs resolving + digest-verifying |
|---|---|
| as they are today (flat) | 33 / 33 |
| after `git mv` alone | 0 / 33 — every one of the 9 broken |
| after `git mv` + ref rewrite | 33 / 33 |

All 64 retained room files were byte-for-byte identical (sha256) across the whole sequence, and a second rewrite pass matched nothing and changed nothing. `entity-session-claim-lease` carries the older single-file `briefing.json` room shape with no request digest; it breaks and repairs under the same one-line rule, so no per-shape special case is needed.

**A form-conflict `Error` locks the first officer out of their own listing.** `findEntityFormConflicts` runs unconditionally at `validate.go:158`, on the read path as well as `--validate`. With a conflict present, plain `spacedock status` prints the error and **exits 1**; `--next` and `--validate` print it and exit 0. So the new hybrid finding cannot be an error at that site. This is the failure the comment at `validate.go:152-156` warns about, now measured.

**Enforcement costs 36 tests.** A one-line folder-form guard in `Prepare` fails 16 tests in `internal/gates`, 19 in `internal/cli`, and 1 in `internal/status`. The same three packages are green at that commit apart from one pre-existing machine-local failure (`TestCodexResolveManifestAgainstInstalledHost`, a locally installed Codex plugin), which is not attributable. All 36 build a flat entity and prepare on it; they funnel through a handful of shared fixtures.

## Proposed approach

**Enforce folder form at `gate prepare`, and ship the conversion that makes enforcement safe.** These are not two options to pick between — either alone is wrong.

Enforcement alone is destructive: `define-fo-moving-target-conflict-ownership`, `merge-guard-requires-preceding-report`, and `preserve-pi-terminal-fields-on-nonterminal-advance` sit at a `validation` gate right now and work fine. Refusing flat form makes their next gate unpreparable, and the obvious remedy — a hand `git mv` — silently destroys all 33 retained refs, as measured above. Conversion alone is insufficient: it repairs today's 9 and leaves `gate prepare` free to mint the tenth.

Together they close the class rather than patch it. Once a room can only be created beside a folder-form entity, `preparedRoomPath` is a pure function of the entity's own directory, `relativeRoomRef` can only ever emit `./review/<stage>/briefing-N`, and the form-dependent ref does not exist to break. Both helpers have exactly one caller each, both inside `Prepare`, so the flat branch of `preparedRoomPath` is deleted rather than left dormant. This also ends the disagreement inside one package: `gate record --round` already refuses flat form for this exact reason (`operation.go:108`), while `gate prepare` creates the same class of artifact and enforces nothing.

Four mechanisms, each with the value AC it serves and the simpler thing it beats:

1. **A folder-form guard in `Prepare`** (serves AC-2). Mirrors `operation.go:108`'s wording and adds the fix: `gate prepare requires folder-form entity <slug>/index.md because review artifacts accumulate beside the entity; convert it with: spacedock convert <slug> --folder`. The same fix clause is added to the `--round` message so the two agree. Simpler alternative — leave `gate prepare` permissive and rely on the converter — is insufficient: it is exactly today's behavior, which mints a new hybrid on every flat entity's first gate.

2. **`spacedock convert <slug> --folder`** (serves AC-1, AC-3, AC-5). One command that moves `<slug>.md` to `<slug>/index.md` and rewrites `room-ref: ./<slug>/…` to `room-ref: ./…` in the same operation, verifies every ref that resolved before still resolves after, and rolls the move back if any does not. On an already-folder entity it repairs stale refs; a second run is a no-op. Simpler alternative — put the two git commands in the refusal message and let the first officer run them — is insufficient, and the spike shows why concretely: the `git mv` half succeeds, the rewrite half is forgettable, and the resulting broken state passes `status --validate` and lists normally in `status`. The step that must not be skipped is the one nothing checks, so the binary owns both halves atomically. Rejected alternative: have `gate prepare` convert the entity itself. It moves the entity path underneath the first officer's in-flight dispatch artifacts, and it breaks `gate prepare`'s own contract that a refusal changes no bytes and a success writes exactly two files.

3. **A warn-tier hybrid finding under `--validate` only** (serves AC-4). Emitted from `gateValidationDiagnostics`, which already walks entities with gate records and already emits `entityEvidenceLine("Warning", …)` scoped to active entities. Condition: the entity is flat and `<slug>/review/` exists. This is the only way to enumerate exposure in a checkout nobody has looked at, which is the risk that cannot be sized from here. It must not go in `findEntityFormConflicts`: measured above, an error there exits 1 on plain `status`. Simpler alternative — a `--list` mode on the converter — is insufficient for the enumeration case, because it has to be run on purpose by someone who already suspects the problem.

4. **Delete the flat branch of `preparedRoomPath`** (serves AC-2). Not a separate feature; it is what makes the guard structural instead of a policy check sitting in front of code that still knows how to do the wrong thing.

Rejected as the primary fix: **resolving legacy refs tolerantly at read time** (retry without the leading slug segment when a ref misses). It is much smaller — roughly 40 lines against five join sites and no test churn — and it unblocks any checkout without anyone running a migration. It is rejected because it makes the doubled-slug ref a permanently supported shape in five places, keeps `gate prepare` minting hybrids, and prices immutability wrong: a fallback path that guesses at a retained room's location is a compatibility branch on exactly the bytes whose value is that nobody may reinterpret them. The residual risk enforcement leaves is a user who hand-converts anyway; `spacedock convert` on an already-folder entity repairs that state, which is why mechanism 2 is specified as repairing rather than only moving.

## Out of scope

Do not change the room layout, the Briefing schema, digest computation, or the gate ceremony's ordering. Do not discard gate history. Do not create a replacement attempt to route around an unreadable room.

The filing-time default and its commission-time parameter belong to a separate entity. Until it lands, newly filed flat entities meet the refusal at their first gate and convert in one command; that friction is the accepted cost, not a defect to solve here. `docs/dev/README.md`'s File Naming guidance is filing guidance and is left alone. `state commit`'s handling of a flat file plus its `<slug>/` companion stays exactly as is — it is how an unmigrated hybrid remains committable.

## Expected surface and tolerance

Estimate net LOC change: +595, across 22 files. Insertions ~666, deletions ~71.

| area | files | net |
|---|---|---|
| `internal/gates/convert.go` (new) | 1 | +150 |
| `internal/gates/prepare.go`, `operation.go` | 2 | +5 |
| `internal/status/validate.go` | 1 | +22 |
| `internal/cli/cli.go` | 1 | +42 |
| new tests (converter, CLI verb, hybrid warning) | 3 | +355 |
| fixture churn in the 36 failing tests | 12 | +15 |
| `docs/site/reference/command-reference.md`, `docs/specs/gate-resolution-frontmatter-contract.md` | 2 | +8 |

Tolerance: net +595 ±150 LOC, 22 ±5 files. The two uncertain terms are the converter's rollback path and the fixture churn across 36 tests; the production guard and the validator warning are each under 25 lines and should not move. The converter lives in `internal/gates` rather than a new package because that is where `room-ref` knowledge and `entitySlug` already are, and `internal/status` already imports it.

Observable semantics this task may change:

- **Command grammar:** one new verb, `spacedock convert <slug> --folder`. `gate prepare` gains one refusal condition.
- **Runtime behavior:** `gate prepare` on a flat entity exits nonzero and writes nothing, where it previously succeeded and created a `<slug>/` companion. `status --validate` gains warn-tier lines.
- **Stored formats:** none. The gates schema, Briefing schema, digest computation, room layout, and the `room-ref` field's grammar are all unchanged. Migrated entities' `room-ref` *values* change; the field's meaning does not.
- **Authority:** none. No new actor, no change to who approves or consumes, no retained byte rewritten.
- **Deliberately unchanged, and asserted as such:** the plain-`status` read path's output and exit code, and `state commit`'s flat-plus-companion unit.

## Acceptance criteria

**AC-1 (value).** For a fixture carrying a consumed gate attempt whose entity is then converted, the number of retained room-refs that resolve and digest-verify is N/N, where the same fixture measures 0/N before the change. Baseline moves the wrong way if any ref stops resolving.
*Test:* Go test in `internal/gates` building a flat entity with a consumed `backlog` attempt, converting it, then counting resolving refs — plus the same count taken against the pre-change binary's behavior, encoded as the broken-input case (a folder-form entity with a `./<slug>/…` ref) which the converter repairs.

**AC-2 (value).** The number of hybrid entities `spacedock` can create is zero. After a `gate prepare` attempt on a flat entity, the process exits nonzero and the state tree is byte-identical and entry-identical to before, with no `<slug>/` directory created.
*Test:* CLI behavior test asserting exit code, the refusal text including `spacedock convert <slug> --folder`, and a before/after tree snapshot (the existing `prepareTreeSnapshot` helper, which compares entry types and a tree digest).

**AC-3 (value).** `spacedock convert <slug> --folder` leaves every retained room file byte-for-byte unchanged. Baseline is the sha256 of every file under `<slug>/review/` before the run; any differing digest fails.
*Test:* Go test hashing the whole `review/` subtree before and after, on a fixture with attempts across three stages including one non-request-backed `briefing.json` room.

**AC-4 (guardrail).** With hybrids present, plain `spacedock status` exits 0 and its output is unchanged, while `status --validate` reports one warning per hybrid and still exits 0.
*Test:* One test with 2 hybrids and 1 clean folder entity asserting: plain-`status` stdout equals the no-hybrid-present baseline and exit 0; `--validate` stdout contains exactly 2 warning lines naming the right slugs and exit 0. This fails if the check is added to `findEntityFormConflicts`, which is the demonstrated wrong home.

**AC-5 (property).** `spacedock convert` is idempotent and repairing: run on an already-folder entity holding stale `./<slug>/…` refs it fixes them; run again it changes no byte of the entity file. A conversion whose rewritten ref would not resolve leaves the entity at its original path with its original bytes.
*Test:* Go test running the converter twice and comparing entity bytes; plus a negative case where a room directory is deleted before converting, asserting nonzero exit, the entity still at `<slug>.md`, and byte-identical content.

**AC-6 (migration completed).** All 9 live hybrids in `docs/dev/.spacedock-state` are folder-form with resolving refs, and their retained room files are byte-identical to their pre-migration digests.
*Test:* Validation stage runs `spacedock convert` against each, records the before/after digest comparison and the resolving-ref count, and confirms `status --validate` reports zero hybrid warnings afterward. Precondition: convert only entities with no open gate attempt and no in-flight dispatch, one commit per entity, because the state checkout is shared.

## Test plan

Fixture and CLI-level Go tests, no live workflow test needed — the mechanism is filesystem and exit-code behavior, all of it reachable offline. Cost is low; the spike already built every fixture shape the tests need.

- `internal/gates/convert_test.go` (new, ~190 lines): AC-1, AC-3, AC-5. Table-driven over room shapes (request-backed, legacy `briefing.json`, multi-stage), plus the rollback case.
- `internal/cli/convert_test.go` (new, ~110 lines): the verb's grammar, usage errors, and the refusal-to-fix round trip — refuse a flat prepare, run the named command, prepare again successfully. This is AC-2's positive half and the one test that proves the two mechanisms compose.
- `internal/status`: the AC-4 warning test, asserting both the `--validate` findings and the plain-`status` invariance.
- Existing suites: 36 tests move to folder-form fixtures. 35 are one-token fixture changes. `TestPrepareRejectsSymlinkedFlatCompanionWithoutChangingBytes` is the exception — its subject is a symlinked flat `<slug>/` companion, which enforcement makes unreachable, so it is retargeted to a symlinked `<slug>/review` inside a folder-form entity. The `validatePreparedRoomAncestry` guard it protects still matters; only the symlink's location changes. Do not delete it.

Green means `go test ./...` and `go test ./... -race` pass with no new failures, allowing for the pre-existing machine-local `TestCodexResolveManifestAgainstInstalledHost`.

## Documentation changes

`docs/site/reference/command-reference.md`, the `gate prepare` row — replace:

> At an actionable current workflow stage (`gate: true` and nonterminal), derive and bind a recorder-ready room for folder or flat form.

with:

> For a folder-form entity (`<slug>/index.md`) at an actionable current workflow stage (`gate: true` and nonterminal), derive and bind a recorder-ready room. Flat entities are refused because review artifacts accumulate beside the entity; convert with `spacedock convert <slug> --folder`.

Same file, a new row after the `spacedock new` row:

> \| `spacedock convert <slug> --folder` \| Convert a flat entity (`<slug>.md`) to folder form (`<slug>/index.md`) in one move and rewrite its retained gate `room-ref` values from `./<slug>/review/…` to `./review/…` in the same operation. Retained Briefing and request bytes are never rewritten. Refuses without changing anything if a rewritten ref would not resolve. Idempotent: on an already-folder entity it repairs stale `./<slug>/` refs, and a second run is a no-op. \|

`docs/specs/gate-resolution-frontmatter-contract.md` — replace the "Folder and flat entities share the same companion-room layout" paragraph and its sentence "For folder form, the entity is `<slug>/index.md`; for flat form it is `<slug>.md` and `<slug>/` is its artifact companion." with:

> Preparation requires a folder-form entity at `<slug>/index.md`, so `room-ref` is always `./review/<stage>/briefing-<attempt>` — relative to the entity's own directory, and therefore invariant under any later move of that directory. Flat entities refuse before locking or writing, the same rule and reason as round recording below. Legacy flat entities carry `./<slug>/review/…` refs written before this rule; `spacedock convert <slug> --folder` moves the entity and rewrites those refs together, and never rewrites a retained room byte. State commit and archive operations continue to treat a flat Markdown file plus its `<slug>/` companion directory as one literal path-scoped unit, so an unmigrated entity stays committable.

## Stage Report: ideation

- DONE: Decide between enforcing folder form at gate prepare and resolving legacy refs tolerantly, and account for the 10 existing hybrids either way
  Enforcement plus a shipped converter; tolerant read rejected with its price stated (~40 LOC, 5 join sites, permanent). The population is 9 hybrids holding 33 refs, not 10 — `repair-sonnet-live-flakes` has an empty `gates:` key and nothing to migrate.
- DONE: Name the migration's exact shape if enforcement wins, preserving every retained Briefing and request digest byte-for-byte
  `spacedock convert <slug> --folder`: move plus `room-ref: ./<slug>/…` → `./…` in one operation, verify-then-rollback. Proven on all 9 real hybrids: 33/33 refs resolve after, 64 retained files byte-identical by sha256.
- DONE: Decide whether the validator should flag flat-plus-rooms as a form conflict, given it already rejects both-forms slugs
  Yes, but warn-tier in `gateValidationDiagnostics` under `--validate` only, NOT in `findEntityFormConflicts`. Measured: an error at that site makes plain `spacedock status` exit 1, locking the FO out of the listing.

### Summary

Reproduced #739 end to end with a built binary — a flat entity's consumed `backlog` attempt makes its later `ideation` gate unpreparable after one `git mv`, failing on `task/task/review/…`. Then exercised the candidate migration against copies of all 9 live hybrids: 33/33 refs resolve today, 0/33 after `git mv` alone, 33/33 after mv-plus-rewrite, with all 64 retained room files byte-for-byte unchanged and the rewrite idempotent. Two measurements changed the design rather than confirming it: a one-line folder-form guard in `Prepare` fails exactly 36 tests across three packages (baselined against a clean run, excluding one pre-existing machine-local Codex failure), and an error emitted from `findEntityFormConflicts` makes plain `status` exit 1 — which is why the new hybrid finding is warn-tier under `--validate` and why AC-4 asserts the read path's exit code directly.

The gate decision worth the captain's attention: enforcement and the converter are one indivisible change, because three hybrids sit at a `validation` gate right now and enforcement alone makes their next gate unpreparable while the obvious hand remedy silently destroys all 33 refs.
