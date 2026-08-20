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

## Cycle 2 measurements: grandfathering, built and run

Cycle 1 estimated. This cycle built the whole change in a throwaway worktree and ran `go test ./...`. Every number below is a `git diff --numstat` reading of a green tree, not an estimate.

**The change, whole and passing:** 11 files, +298 / −15, **net +283**.

| part | net |
|---|---|
| production — `Prepare` guard (+21), `--round` message (+1/−1), validator (+45) | **+66** |
| fixture churn on pre-existing tests, 6 files | **+16** |
| new proof — the guard test and the validator tests | **+201** |

**Grandfathering does not reduce the fixture churn.** 34 tests fail under grandfathering against 36 under full enforcement — near-identical, because every existing test creates the FIRST room beside a flat entity, which both variants refuse. What grandfathering changes is the *repair*: a 6-line helper that gives a flat fixture the companion it now needs, called at 6 sites, plus one assertion and one seed fixture moved to folder form. That is +16 net. Full enforcement forces the fixtures to folder form instead, which rewrites entity paths through the assertions — and a mechanical attempt at that took the CLI package from 19 failures to 28 before being reverted.

Two pre-existing tests needed more than a fixture line, both because their subject is the flat companion itself: `TestPrepareRejectsSymlinkedFlatCompanionWithoutChangingBytes` now removes the fixture's companion before symlinking it, and `TestPrepareRollsBackPublishedRoomAfterBindingWriteFailure` moves to folder form with its two branches re-aimed at the `review/` root, because in flat form the companion can no longer be the thing that gets created and rolled back.

**On the live checkout, read-only:** the validator flags exactly the 9 hybrids, names all 9 slugs, reports 0 findings against the 68 folder entities, finds 0 unresolved refs across all 33 live room-refs, and both `status --validate` and plain `status` exit 0.

## Proposed approach

**Guard the mint, grandfather the 9, and make a hand conversion checkable. No new command.**

`gate prepare` refuses to create the FIRST prepared room beside a flat entity, and grandfathers an entity that already holds one. The predicate is the thing being prevented: the entity is flat and `<state>/<slug>/review/` does not exist. `<slug>/review/` beside a flat file can only have come from `gate prepare` — `gate record --round` has always refused flat form — so it reads exactly as "has this entity already got rooms".

That is the whole prevention, and it costs nothing anyone is using. The 9 hybrids resolve correctly today precisely because their `./<slug>/review/…` refs are right for the form they are in; three of them sit at a `validation` gate and their next `gate prepare` keeps working untouched. Nothing is migrated, no retained byte is rewritten, and the closed set of 9 shrinks on its own as those entities archive.

Then two warn-tier findings under `--validate`, both in `gateValidationDiagnostics`, which already walks entities with gate records and already emits active-scope warnings:

1. **A flat entity holding gate rooms**, carrying the whole remedy rather than just the finding: `git mv <slug>.md <slug>/index.md` AND rewrite every `room-ref: ./<slug>/` to `room-ref: ./` in the same commit. The remedy is the payload. A reader who learns only that the entity is a hybrid still hand-converts and destroys every retained room; the field first officer who hit #739 succeeded because they did both halves.
2. **A retained room-ref that no longer resolves.** This is the one thing the first finding cannot do: it fires only while the entity is flat, so it vanishes after the move — including a botched move. Without this check, an operator who follows the remedy has no way to confirm it landed, and cycle 1 measured that `status --validate` reports VALID on the broken state. This turns #739's end state from a hard failure at the next gate into a validate-time line naming the exact ref. It is +31 production LOC against the ~+300 the converter cost, and it is the only mechanism here that addresses the residual rather than the population.

**The converter is dropped.** The checklist asked for a case grandfathering plus the warning does not cover, that is not the one the field already handled by hand. There is no such case. Every route to a broken ref runs through a conversion, conversion is now optional for all 9, the remedy is stated at the place an operator looks for state health, and the result is checkable. A permanent verb with its own tests, aimed at a closed set of 9 that needs no migration, is a mechanism in search of a transaction.

**Retained from cycle 1:** the tolerant read stays rejected, for the reason given — a fallback that guesses a retained room's location is a compatibility branch on exactly the bytes whose value is that nobody may reinterpret them.

**One accepted item is not available under grandfathering, and this is a consequence, not a reopening.** Deleting the flat branch of `preparedRoomPath` requires that flat-form room creation never happen again. Grandfathering keeps it happening for the 9, so that branch and `relativeRoomRef`'s form-dependent output are both live code and must stay. Cycle 1's structural argument — that enforcement makes the fragile ref impossible by construction — is the thing being traded away for the 100-LOC ceiling and for not touching 9 working entities. The trade looks right; it should be made with open eyes.

## Out of scope

Do not change the room layout, the Briefing schema, digest computation, or the gate ceremony's ordering. Do not discard gate history. Do not create a replacement attempt to route around an unreadable room. Do not migrate the 9: they work, and nothing here requires them to move.

The filing-time default and its commission-time parameter belong to a separate entity. Until it lands, a newly filed flat entity meets the refusal at its first gate; the message names both ways out (file it as `<slug>/index.md`, or `git mv` it). `docs/dev/README.md`'s File Naming guidance is filing guidance and is left alone. `state commit`'s flat-plus-companion unit stays exactly as is — it is how a grandfathered hybrid remains committable.

## Expected surface and tolerance

Estimate net LOC change: **+283, across 11 files** — measured on a green tree, not projected. Insertions 298, deletions 15.

| file | net |
|---|---|
| `internal/gates/prepare.go` (guard + its comment) | +21 |
| `internal/status/validate.go` (both findings + helper) | +45 |
| `internal/gates/operation.go` (`--round` remedy text) | 0 |
| `internal/gates/prepare_test.go` (new guard test + 3 fixture repairs) | +76 |
| `internal/status/hybrid_flat_rooms_warn_test.go` (new) | +125 |
| `internal/cli/` × 5 (fixture grandfathering helper and calls) | +16 |
| `internal/status/gate_readiness_initial_seed_test.go` (seed → folder form) | 0 |
| `docs/site/reference/command-reference.md`, `docs/specs/gate-resolution-frontmatter-contract.md` | not yet written; +6 |

Tolerance: net +283 ±40 LOC, 13 ±2 files. The uncertainty is only in the doc wording and in review-driven edits to the two new tests; the production code and the fixture churn are exact.

**Answering the ceiling directly: no, this does not land under 100 net LOC. It lands at +283.** The mechanisms do — the guard plus the hybrid warning are +35 net production, and all three production changes together are +66. The other +217 is fixture churn the guard forces (+16, unavoidable in any variant) and the tests that prove the guard and the findings (+201). Under a hard 100-net ceiling counting tests, the only way through is to prove less. Two cuts are available, in the order I would make them:

- Drop the unresolved-ref check: −31 production, −~40 test. Lands at ~+212. Cost: the residual in the next section goes back to unfixed.
- Drop the grandfathered branch of the guard test and the read-path invariance test: ~−60. Lands at ~+150. Cost: nothing then fails if the guard stops grandfathering — which would strand the three hybrids at validation gates — or if the finding is later moved to `findEntityFormConflicts`, the placement measured in cycle 1 to exit 1 on plain `status`.

I do not recommend either. Both delete the checks on the two mistakes this design is most likely to suffer later.

Observable semantics this task may change:

- **Command grammar:** no new verb, no new flag. `gate prepare` gains one refusal condition; two existing refusal messages gain remedy text.
- **Runtime behavior:** `gate prepare` on a flat entity holding no rooms exits nonzero and writes nothing, where it previously succeeded and created the companion. A flat entity that already holds rooms behaves exactly as today. `status --validate` gains two warn-tier findings.
- **Stored formats:** none, and nothing is migrated. No `room-ref` value changes anywhere.
- **Authority:** none.
- **Deliberately unchanged, and asserted:** the 9 hybrids' refs and their next gate; the plain-`status` read path's output and exit code; `preparedRoomPath`'s flat branch and `relativeRoomRef`; `state commit`'s flat-plus-companion unit.

## What grandfathering leaves unfixed

Stated plainly, because the answer to "is it survivable" depends on it being stated.

1. **The 9 stay conversion-fragile forever.** Their refs remain slug-prefixed. Anyone who moves one to folder form without the rewrite destroys all its retained rooms — measured: 0 of 33 refs survive a bare `git mv`.
2. **The fragile mechanism stays.** `relativeRoomRef` remains form-dependent and `preparedRoomPath` keeps its flat branch, so the class is prevented by a guard rather than made impossible. A future change that removes or bypasses the guard reopens it.
3. **`gate record --round` still refuses every flat entity, grandfathered or not.** Pre-existing, unchanged, and not currently live: this workflow records correction rounds as a `### Feedback Cycles` body section, and every `review/*/round-N` directory in the checkout belongs to a folder-form entity. A grandfathered hybrid that ever needs a real `--round` must convert first.

**The failure a reader of the warning still walks into** is not ignorance — the remedy is in the line. It is that they cannot tell whether they got it right. The rewrite is a hand edit to YAML that nothing verifies; cycle 1 measured `status --validate` reporting VALID on a conversion whose refs were all broken. Worse, the hybrid warning fires only while the entity is flat, so it disappears after the move whether the move was correct or not. They find out at the next gate, which is exactly #739.

That is the one thing the second finding exists for, and with it the residual is survivable: the population is closed at 9, every one works today, conversion is never required, and if someone converts anyway the mistake surfaces as a validate-time line instead of a blocked ceremony.

## Acceptance criteria

**AC-1 (value — prevention).** The number of new hybrids `spacedock` can create is zero: on a flat entity with no `<slug>/review/`, `gate prepare` exits nonzero, the state tree is entry- and digest-identical before and after, and no `<slug>/` directory appears.
*Test:* `TestPrepareRefusesNewFlatCompanionAndGrandfathersExistingRooms/new-companion-refused` asserts the refusal text, `prepareTreeSnapshot` equality, and the companion's absence. Fails if the guard is removed, or if it is placed after any mutation or after the entity lock.

**AC-2 (value — no collateral).** Every grandfathered entity keeps working unchanged: a flat entity that already holds `<slug>/review/` prepares its next gate and binds a `./<slug>/review/<stage>/briefing-N` ref.
*Test:* the same test's `existing-rooms-grandfathered` branch asserts the room path and the exact slug-prefixed ref. Fails if the guard refuses grandfathered entities — which would strand the three hybrids sitting at validation gates — and fails if anyone "corrects" the ref shape, which would silently invalidate 33 retained refs.

**AC-3 (value — visibility).** Every grandfathered hybrid is reported, and nothing else is: one warning per flat-entity-with-rooms, zero for folder-form entities. Against the live checkout that is 9 warnings naming the 9 slugs and 0 findings across 68 folder entities, exit 0.
*Test:* fixture test over 2 hybrids and 1 clean folder entity asserting the count, both remedy tokens (the `git mv` and the ref rewrite), and no finding for the folder entity; plus the read-only live run recorded above. Fails if the finding drops either half of the remedy, or fires on folder form.

**AC-4 (guardrail — the measured risk).** The plain `status` read path is unaffected with hybrids present: exit 0, no new stderr line, every entity still listed.
*Test:* `TestHybridFindingLeavesPlainStatusReadPathUnaffected`. This fails if the finding is added to `findEntityFormConflicts`, the placement cycle 1 measured as exiting 1 on plain `status` and locking the first officer out of the listing that shows the broken entity.

**AC-5 (value — closes the residual).** A conversion that omits the ref rewrite is a validate-time finding, not a gate failure: for an entity whose retained ref no longer resolves, `status --validate` names that exact ref and exits 0; on a healthy checkout the count is 0.
*Test:* `TestValidateReportsRetainedRoomThatNoLongerResolves` performs #739's hand conversion inside the fixture and asserts exactly one finding naming the ref, with no false positive on the healthy siblings. The healthy-count assertion is what makes it falsifiable in the over-reporting direction — it already caught an incomplete room directory in the fixture during the spike. Live baseline: 0 findings across 33 refs.

## Test plan

Written and green. `go test ./...` passes on the whole change; the only failure in the tree is the pre-existing machine-local `TestCodexResolveManifestAgainstInstalledHost` (a locally installed Codex plugin), which fails identically at the base commit. `gofmt -l ./cmd ./internal` is clean. Implementation should re-run with `-race`, which this spike did not.

- `internal/gates/prepare_test.go` — `TestPrepareRefusesNewFlatCompanionAndGrandfathersExistingRooms`, two branches, covering AC-1 and AC-2.
- `internal/status/hybrid_flat_rooms_warn_test.go` — three tests covering AC-3, AC-4, AC-5 over one 3-entity fixture.
- 34 pre-existing tests repaired: 31 by the fixture-grandfathering helper, 3 by the targeted edits described above.

## Documentation changes

`docs/site/reference/command-reference.md`, the `gate prepare` row — replace:

> At an actionable current workflow stage (`gate: true` and nonterminal), derive and bind a recorder-ready room for folder or flat form.

with:

> At an actionable current workflow stage (`gate: true` and nonterminal), derive and bind a recorder-ready room. A flat entity that holds no rooms yet is refused, because the room would land in a `<slug>/` companion and bind a ref that breaks if the entity later becomes `<slug>/index.md`; file it as `<slug>/index.md` instead. A flat entity that already holds rooms is grandfathered and prepares as before.

Same file, in the `spacedock status` row's `--validate` description, add: `--validate` also warns when a flat entity holds gate rooms, and when a retained room-ref no longer resolves.

`docs/specs/gate-resolution-frontmatter-contract.md` — replace "Folder and flat entities share the same companion-room layout:" and its following sentence "For folder form, the entity is `<slug>/index.md`; for flat form it is `<slug>.md` and `<slug>/` is its artifact companion." with:

> The room layout is the same for both entity forms:
>
> ```text
> <state-root>/<slug>/review/<stage>/briefing-<attempt>/
> ```
>
> but `room-ref` is written relative to the entity file's own directory, so folder form binds `./review/…` while flat form binds `./<slug>/review/…`. Only the folder-form ref is invariant under a later move of the entity. Preparation therefore refuses to create the first room beside a flat `<slug>.md`: for folder form the entity is `<slug>/index.md`, and for flat form `<slug>/` would be a companion whose refs break on conversion. Flat entities that already hold rooms are grandfathered, and their slug-prefixed refs stay correct while they stay flat; converting one requires `git mv <slug>.md <slug>/index.md` and rewriting every `room-ref: ./<slug>/` to `room-ref: ./` in the same commit, and `status --validate` reports both the grandfathered shape and any ref that stops resolving. State commit and archive operations continue to treat a flat Markdown file plus its `<slug>/` companion directory as one literal path-scoped unit.

### Feedback Cycles

- Cycle 1: REVISED — captain at the ideation gate; surface estimate 22 files/net +595 vs the captain's 100-net ceiling (595%); AC unchanged. The defect closes in +5 (the `Prepare` guard) plus +22 (the validator warning). The other +492, 83% of the estimate, is a new `spacedock convert` command and its tests. The design rejected doing nothing because a hand `git mv` destroys 33 retained refs and the rewrite half is forgettable — measured, and true as far as it goes. The captain reports the field outcome that contradicts its conclusion: the repository and first officer that actually hit issue #739 converted and repaired the entity by hand without difficulty. So the measurement argues the migration must be SAFE IF PERFORMED, not that it must be performed or shipped as a permanent verb. Reprice with grandfathering as the primary option: the 9 hybrids resolve correctly today and break only on conversion, so refusing to mint a NEW hybrid while leaving the existing closed set alone breaks no gate and needs no migration. FO note: the FO dispatched this without capping the surface, and the scope note recorded the captain's enforcement steer without the ceiling. That is on the FO, not the worker — the estimate is an honest reading of the design it was asked to shape.

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

## Stage Report: ideation (cycle 2)

- DONE: Price grandfathering as the primary option and say plainly whether it closes the defect under 100 net LOC
  No — it lands at net +283 (11 files, measured on a green tree, not estimated). The mechanisms are +66 net production; the other +217 is unavoidable fixture churn (+16) and the tests that prove the change (+201). Two named cuts reach ~+212 and ~+150; I recommend neither and said why.
- DONE: State what grandfathering leaves unfixed and whether the validator warning makes that survivable
  Three items in `## What grandfathering leaves unfixed`. The failure a warning-reader still walks into is that a correct-looking hand conversion is unverifiable — the hybrid warning disappears after the move, botched or not, and cycle 1 measured `--validate` reporting VALID on all-33-refs-broken. Survivable only with the second finding, which reports an unresolved ref for +31 production LOC.
- DONE: Keep or drop the converter on evidence, not on the risk measurement alone
  Dropped. No case exists that grandfathering plus the remedy-carrying warning does not cover and that is not the hand conversion the field already performed. Every route to a broken ref runs through a conversion that is now optional for all 9 and checkable afterward.

### Summary

Built the whole grandfathering change and ran it rather than estimating it: `Prepare` refuses the first room beside a flat entity and grandfathers the 9, plus two warn-tier `--validate` findings — the hybrid shape carrying the full conversion remedy, and an unresolved retained ref. `go test ./...` is green except the pre-existing machine-local Codex failure, and against the live checkout the validator flags exactly the 9 hybrids, 0 of 68 folder entities, and 0 unresolved refs across all 33 live room-refs, with both `--validate` and plain `status` at exit 0.

Two measurements are worth the gate's attention. First, grandfathering does not reduce the fixture churn — 34 tests fail against full enforcement's 36, because every existing test creates the FIRST room beside a flat entity; what it changes is that the repair is a 6-line helper (+16 net) instead of rewriting entity paths through the assertions. Second, one item the FO marked accepted is not available under grandfathering: deleting the flat branch of `preparedRoomPath` requires flat-form room creation to stop entirely, so that branch and `relativeRoomRef`'s form-dependent output must stay. Cycle 1's structural argument is what the 100-LOC ceiling buys, and the trade should be made knowingly. Cycle 1's own overestimate was ~2x on total, from estimating test volume rather than writing it — the +16 fixture-churn figure was the one part that estimated accurately.
