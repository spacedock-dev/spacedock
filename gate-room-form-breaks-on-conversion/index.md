---
id: gf0jvhj4y8vjhd6ww62vsb2q
title: Gate rooms make a flat entity a hybrid, and converting it breaks every later gate
status: validation
source: "Issue #739 (captain, 2026-08-20): gate prepare fails after a flat-to-folder conversion because a retained room-ref resolves relative to the new entity folder and duplicates the slug."
started: 2026-08-20T19:59:06Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-gate-room-form-breaks-on-conversion
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
            - id: gate-attempt:gf0jvhj4y8vjhd6ww62vsb2q-ideation-2
              briefing:
                id: briefing:gf0jvhj4y8vjhd6ww62vsb2q:ideation:attempt-2:revision-1
                digest: sha256:12f0bae5b1bd133d005061d80d50dbd6de7ba4a2156afc5a9836ae1cd89837e8
                request-digest: sha256:2528069ebf817b30c260f8544deffc9fc1584a9d1f8bac7748f4b8d211a89c68
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:gf0jvhj4y8vjhd6ww62vsb2q:ideation:2
                briefing: briefing:gf0jvhj4y8vjhd6ww62vsb2q:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-20T19:27:53.006647Z"
                decision: approve
                reason: 'Captain: ''approve both''. Grandfathering approved, converter dropped on the field evidence that the repository which hit #739 converted by hand without difficulty. Approval ratifies +283 net against the 100 ceiling: production is +66, the numbers are measured on a green tree, and the two available cuts to ~150 both delete the checks on this design''s most likely later mistakes.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:gf0jvhj4y8vjhd6ww62vsb2q:validation
          stage: validation
          attempts:
            - id: gate-attempt:gf0jvhj4y8vjhd6ww62vsb2q-validation-1
              briefing:
                id: briefing:gf0jvhj4y8vjhd6ww62vsb2q:validation:attempt-1:revision-1
                digest: sha256:b55953e69b3b4706f163d8f85aac9f68f128e0e7b409a6a6d83d1f68fd3ba28e
                request-digest: sha256:4f75ec56794b9aa8c3671a97541036b7760c94bfaf04631a535142414fdc9d2f
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:gf0jvhj4y8vjhd6ww62vsb2q:validation:1
                briefing: briefing:gf0jvhj4y8vjhd6ww62vsb2q:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-20T22:33:47.323393Z"
                decision: approve
                reason: 'Captain: ''push it and approve CI for the stack tip''. Validation PASSED with no material finding: all five ACs reproduced on a throwaway clone, every live number re-derived independently, and the armed-gate concern resolved — the three live prepares ran in clones and never touched live state. Two polish findings, both of the same class: an assertion that cannot fail, and a claim about an assertion that is not true. Surface 292 net against the ratified 283 plus-or-minus 40.'
              application:
                target-stage: done
                state: superseded
mod-block:
pr:
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

## Cycle 2 narrowing: the workflow declares the form

Cycle 1 shipped the refusal unconditional, and it refused the shipped default. `--folder` is opt-in on `spacedock new`, so a flat entity holding no rooms is what filing mints, and both live lanes reddened when the fixture workflow's first gate met the guard. The captain's rule: refuse only where the WORKFLOW declares the form.

A workflow declares it as `entity-form: folder` in README frontmatter — the sibling of `entity-type`, `entity-label`, and `id-style`. Preparation only READS the key. `qpa` (`workflow-declared-entity-form`) owns the commission-time parameter and the `spacedock new` default; nothing here writes the key, and no workflow in this repo declares it yet.

- Workflow declares folder form, entity is flat and holds no rooms → refuse, with the two-part remedy.
- Workflow declares nothing → allow either shape, no refusal. This is every workflow today, including every live fixture.
- An unreadable or unparseable README is not a declaration. `validatePreparedStage` reads the same file a few lines later and owns that error, so the guard introduces no new failure mode.

**Grandfathering is kept, and the reason changed.** Under the unconditional guard it protected the 9 hybrids from a rule they never opted into. Under the declaration it protects them from a rule a workflow opts into later — which is the concrete next step here, since `docs/dev` is both the most likely first workflow to declare `entity-form: folder` and the one holding all 9 hybrids, 3 of them at active validation gates. Without the branch, adding one README line arms a gate failure that fires at each hybrid's next gate: green now, blocked later, which is exactly #739's shape and the thing this task exists to remove. The branch costs four production lines and one test case now that the declaration removed all its fixture churn. `status --validate` still names every hybrid with the full conversion remedy, so the declaration is not silent about them — it just does not block them.

## Acceptance criteria

**AC-1 (value — prevention, where the workflow asked for it).** `spacedock` mints no new hybrid in a workflow that declares folder form, and refuses nothing in a workflow that does not. Where the README declares `entity-form: folder` and the entity is flat with no `<slug>/review/`, `gate prepare` exits nonzero, the state tree is entry- and digest-identical before and after, and no `<slug>/` directory appears. Where the README declares no form, the same flat entity prepares and binds its room.
*Test:* `TestPrepareRefusesFlatEntityOnlyWhereTheWorkflowDeclaresFolderForm`, branches `declared-folder-form-refuses-flat` and `undeclared-workflow-prepares-flat`. Falsified in both directions: deleting the guard reds the refusal branch, and deleting the declaration check reds the undeclared branch. Also reds if the key spelling or the declared value changes, since either makes a declaring workflow read as undeclared.

**AC-2 (value — no collateral).** Every grandfathered entity keeps working unchanged, including after its workflow declares folder form: a flat entity that already holds `<slug>/review/` prepares its next gate and binds a `./<slug>/review/<stage>/briefing-N` ref.
*Test:* the same test's `declared-folder-form-grandfathers-rooms` branch asserts the room path and the exact slug-prefixed ref. Fails if the guard refuses grandfathered entities — which is what would strand the three hybrids sitting at validation gates the day `docs/dev` declares a form — and fails if anyone "corrects" the ref shape, which would silently invalidate 34 retained refs.

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

**Cycle 2 supersedes the first and third bullets.** The guard test is `TestPrepareRefusesFlatEntityOnlyWhereTheWorkflowDeclaresFolderForm`, three branches, covering AC-1 and AC-2. No pre-existing test needs repair: gating the refusal on the declaration means a fixture workflow that declares nothing is never refused, so all 34 repairs reverted to their stronger pre-change assertions. The `--validate` bullet is unchanged. Cycle 2 also adds what cycle 1 lacked: a targeted live journey run locally against the built change, since `gate prepare` is on every live journey's path.

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

**Cycle 2 amends both replacements**, since the refusal is now conditional. The `gate prepare` row keeps "for folder or flat form" and adds the condition after it: where the workflow README declares `entity-form: folder`, a flat entity holding no rooms is refused instead, a workflow declaring no form accepts either shape, and a flat entity already holding rooms is grandfathered under the declaration. The contract's paragraph gains the same condition — a workflow states its form with `entity-form: folder` and preparation refuses only where that declaration is present — with the grandfathering, conversion remedy, and state-commit sentences unchanged. One doc is added: `docs/site/concepts/workflows-and-entities.md`, whose "Each work item is one file" paragraph is where a reader looks for how a workflow states its entity form. The `--validate` addition is unchanged.

### Feedback Cycles

- Cycle 1: REVISED — captain at the ideation gate; surface estimate 22 files/net +595 vs the captain's 100-net ceiling (595%); AC unchanged. The defect closes in +5 (the `Prepare` guard) plus +22 (the validator warning). The other +492, 83% of the estimate, is a new `spacedock convert` command and its tests. The design rejected doing nothing because a hand `git mv` destroys 33 retained refs and the rewrite half is forgettable — measured, and true as far as it goes. The captain reports the field outcome that contradicts its conclusion: the repository and first officer that actually hit issue #739 converted and repaired the entity by hand without difficulty. So the measurement argues the migration must be SAFE IF PERFORMED, not that it must be performed or shipped as a permanent verb. Reprice with grandfathering as the primary option: the 9 hybrids resolve correctly today and break only on conversion, so refusing to mint a NEW hybrid while leaving the existing closed set alone breaks no gate and needs no migration. FO note: the FO dispatched this without capping the surface, and the scope note recorded the captain's enforcement steer without the ceiling. That is on the FO, not the worker — the estimate is an honest reading of the design it was asked to shape.

### Feedback Cycles

- Cycle 1: REJECTED — captain, after the stack-tip CI run graded the delivered guard (narrow rule change, not a design reset); surface 13 files/net +292 vs estimate +283 ±40 (103%); AC narrowed: AC-1's refusal condition is rewritten. Material: the guard refuses `gate prepare` on ANY flat entity that holds no rooms, which is the shipped default shape — `--folder` is opt-in and no workflow declares a form — so it refuses the product's own default, not merely a fixture convenience. Both live lanes reddened at tip run 32425117427: the journey's own words are "Gate preparation blocked: the workflow requires folder-form `auto-continue-task/index.md`, but the task is stored as `auto-continue-task.md`". Reproduced locally against the built guard binary on a fresh flat entity before reading the CI log. Captain's rule: LOOSEN — refuse only when the WORKFLOW declares folder form. A workflow that declares nothing allows either shape, so enforcement arrives with the declaration rather than ahead of it, and `qpa` (`workflow-declared-entity-form`) owns the commission-time parameter and the `spacedock new` default. FO note, the cause of the miss: the ideation declared "no live workflow test needed" and the FO approved that twice, for a change to a command every live journey drives. The workflow README's live-run rule enumerates grader, fixture, runner, and contract text, and this change was none of those — it was product code the journeys exercise. The FO read the enumeration instead of its intent. That is on the FO, not the worker, and it is the second occurrence of this class this session.

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

## Stage Report: implementation

- DONE: Verify the inherited on-disk diff hunk by hunk against the entity's ACs, then commit what is green before further work
  All 13 hunks matched their AC; committed as `a64751c26` before any further work. One defect found and fixed in `2fcaa35e5`: in `gateValidationDiagnostics` the two comment paragraphs each sat above the other's code.
- DONE: Land the docs and prove by name that the three flat-with-rooms entities at validation gates still prepare
  Docs in `2fcaa35e5`. All three prepare at exit 0 — `define-fo-moving-target-conflict-ownership` → `validation/briefing-6`, `merge-guard-requires-preceding-report` → `validation/briefing-1`, `preserve-pi-terminal-fields-on-nonterminal-advance` → `validation/briefing-2`. No destructive variant.
- DONE: Full suite and race pass, stage report with executed falsifiers
  `go test ./... -race` green in every package; `gofmt -l ./cmd ./internal` clean. Only failure is pre-existing `TestCodexResolveManifestAgainstInstalledHost`, which fails identically on unmodified `main` (verified this session).

### Falsifiers executed

- **AC-1.** `TestPrepareRefusesNewFlatCompanionAndGrandfathersExistingRooms/new-companion-refused` asserts the refusal text, `prepareTreeSnapshot` equality (entry types + tree digest) and the companion's absence; it fails if the guard is removed or moved after the entity lock — the guard sits before `lockEntity`, after read-only input validation only. Also run against a **real** entity: stripping `merge-guard-requires-preceding-report`'s companion in a clone of both repos makes prepare exit 1 with the folder-form message, mints no companion, and leaves the whole state tree byte-identical with 0 dirty files.
- **AC-2.** The same test's `existing-rooms-grandfathered` branch asserts the room path and the exact slug-prefixed ref `./task/review/validation/briefing-1`; it fails if the guard refuses a grandfathered entity, and fails if anyone "corrects" the ref shape. Against the three live entities above, all pre-existing retained files were byte-identical by sha256 after preparing (14, 2 and 8 files respectively) — so the new attempt appended without touching a retained byte.
- **AC-3.** Read-only live diff of `--validate` stderr, guard binary vs the **unmodified `main` binary**: exactly 9 lines added, 0 removed, naming the 9 hybrid slugs; `--validate` exits 0 printing VALID. The fixture test asserts 2 hybrid warnings, both remedy tokens, and no finding for the folder-form entity; it fails if either half of the remedy drops or the finding fires on folder form.
- **AC-4.** `TestHybridFindingLeavesPlainStatusReadPathUnaffected` plus the live run: plain `status` exits 0 with 0 bytes of stderr and all 19 table rows. Fails if the finding is moved to `findEntityFormConflicts`, the placement cycle 1 measured as exiting 1.
- **AC-5.** `TestValidateReportsRetainedRoomThatNoLongerResolves` performs #739's hand conversion inside the fixture and asserts exactly one finding naming the ref with no false positive on the healthy siblings; the healthy-count assertion is what falsifies it in the over-reporting direction. Live baseline: 0 unresolved refs across the 34 held by the 9 hybrids.

### Surface

13 files, +312 / −20, **net +292** against the declared +283 ±40 LOC and 13 ±2 files — within tolerance. Docs came in at net +9 rather than the projected +6 because the frontmatter contract's existing sentence about state commit treating the flat file plus companion as one path-scoped unit was preserved verbatim, including tracked deletions and rollback, rather than replaced by a shorter restatement.

Population confirmed unchanged in kind but grown by one since the spike: the 9 hybrids now hold 34 retained room-refs, not 33 (`define-fo-moving-target-conflict-ownership` gained one). 122 flat and 68 folder entities; 10 flat entities carry a `gates:` key and the 10th is still empty, so it is correctly not flagged.

### Summary

Inherited the prior worker's uncommitted tree, verified every hunk against the five ACs before committing it, and found one real defect — reversed comment paragraphs in `gateValidationDiagnostics` — which is fixed. The named risk did not materialize: all three flat-with-rooms entities at validation gates prepare at exit 0 under the guard, proven against clones of both real repos with history so the retained rooms' `git-root` selected sources still resolve, with every pre-existing retained file byte-identical afterward.

Two things a reviewer should weigh. The strongest new evidence is the read-only live diff against the unmodified `main` binary: this change adds exactly 9 `--validate` lines and removes none, which is AC-3's "and nothing else" measured rather than asserted. Separately, the first non-race full run showed 3 timing failures in `internal/ensigncycle` (Codex 250ms no-progress budget) that are not attributable — the package is green in isolation with this change present (291s alone against 522s under full-suite load) and green under the `-race` full suite, so the failures need CPU contention, not this diff. `internal/ensigncycle` does import both changed packages, so this was established behaviorally rather than by import graph.

## Stage Report: validation

- DONE: Reproduce each AC's cited evidence, re-running the named falsifiers on a throwaway checkout and confirming each reds exactly its claimed cases
  All five ACs reproduced on a throwaway clone at `2fcaa35e5`, diffed against the stack parent `8f29ad577` (stricter isolation than "unmodified main" — it excludes #743/#738). Eleven mutations run; ten red exactly their claimed case, one does not (finding 1).
- DONE: Independently confirm the three live entities still prepare and that the nine hybrids' retained files were untouched, from this run's evidence
  Reproduced from room contents in a fresh clone of both repos at live state HEAD `3a5081022`: all three prepare at exit 0 binding `validation/briefing-6`, `briefing-1`, `briefing-2` — the exact attempts live state would next mint — each with a slug-prefixed `./<slug>/review/validation/briefing-N` ref. Of 66 pre-existing retained files across the 9 hybrids, path-keyed sha256 gives **0 changed, 0 vanished, 6 added** (only the three new rooms' `gate-briefing.json` + `request.json`).
- DONE: Judge the crash-inheritance chain: the diff was applied from a spike whose measuring session died, verified by a successor — decide whether any hunk still lacks an owner who checked it against its AC
  Twelve of 13 hunks have an owner who checked them against an AC by exercise. One did not and now does: `internal/gates/operation.go`'s `--round` remedy text is covered by no AC and asserted by no test (finding 2). I exercised it directly against both binaries.

### Falsifier mutation matrix

Each row mutates source in the throwaway checkout and re-runs the named test.

- **AC-1** guard deleted -> `new-companion-refused` FAILS. Guard moved below every mutation (before the final return) -> FAILS. Guard moved directly after `lockEntity` -> **still PASSES** (finding 1).
- **AC-2** grandfathering early-return deleted -> `existing-rooms-grandfathered` FAILS. `relativeRoomRef` "corrected" to strip the slug prefix -> FAILS while `new-companion-refused` stays green, so the branches are independent.
- **AC-3** literal `git mv` token removed -> remedy test FAILS; the `room-ref` rewrite clause removed -> FAILS; finding made to fire on folder form -> FAILS. Note the flat-form check alone is not falsifiable: the stat path `<dir>/<slug>/review` is itself form-dependent, so deleting the `index.md` condition changes no behavior. A genuine folder-form firing required editing both.
- **AC-4** hybrid finding relocated into `findEntityFormConflicts` (the placement cycle 1 measured as exiting 1) -> `TestHybridFindingLeavesPlainStatusReadPathUnaffected` FAILS. Claimed falsifier accurate.
- **AC-5** `unresolvedRoomRefs` forced to nil -> FAILS; forced to report every ref -> FAILS. Falsifiable in both the under- and over-reporting directions.

### Live evidence (read-only against the live workflow)

- **AC-3** `--validate` stderr diffed between the two binaries: exactly **9 lines added, 0 removed**, 9 distinct slugs, all 9 carrying both remedy halves; stdout byte-identical; both exit 0.
- **AC-4** plain `status`: both binaries exit 0 with **0 bytes of stderr** and byte-identical stdout.
- **AC-5** live baseline **0 unresolved refs**. Population recounted independently: 122 flat, 68 folder, 10 flat carry `gates:`, 9 are hybrids holding **34** refs — confirming the implementation's corrected 34 against the spike's 33.
- **AC-1** live: stripping `merge-guard-requires-preceding-report`'s companion in the clone makes prepare exit 1 with the folder-form message, mints no companion, leaves the tree entry-list identical, and leaves **no `.gates.lock` residue**.
- Placement read directly: everything above `prepare.go:95` is pure input validation, and the guard precedes `lockEntity` — itself the first byte-writing call (`<entity>.gates.lock`).

### Reviewer findings

1. **Polish (evidence defect) — AC-1's cited falsifier overstates by one clause.** AC-1 claims the test "Fails if the guard is removed, or if it is placed after any mutation **or after the entity lock**." The final clause is false: moving the guard directly below `lockEntity` leaves both branches green. It is also harmless — `lockEntity` creates `<entity>.gates.lock` and the deferred `unlock()` removes it, so a refusal there still leaves the tree entry- and digest-identical, which is what AC-1's value claim actually asserts. The mutation that would violate that claim (guard below room publication) does red the test. No value AC fails. Promote to material only if `lockEntity` ever leaves durable residue on its error path.
2. **Polish — one hunk carried no AC and no assertion.** `internal/gates/operation.go`'s new `gate record --round` remedy text is asserted by no test; the only `requires folder-form entity` assertion in the tree is `prepare_test.go:1132`, which asserts the *prepare* message. Its prior owner-check was "the suite is green", which cannot fail for an unasserted string. Exercised live: base prints the bare refusal, the guard binary appends the two-part remedy, using the same literal `<slug>` placeholder style as the pre-existing half of that message. Correct as written; no user-visible loss.

### Checks run

`go test ./... -race` on the throwaway checkout: **19 packages green**, including `internal/ensigncycle` (290s, the package the implementation flagged for non-race timing flakes). The only failure is `TestCodexResolveManifestAgainstInstalledHost` in `internal/cli`, which I confirmed fails identically at base commit `8f29ad577` without the change — machine-local, not attributable. **Zero data races.** `gofmt -l ./cmd ./internal` clean. Surface re-measured at 13 files, +312/-20, **net +292** against the declared +283 +/-40 and 13 +/-2 — within tolerance.

Fixture churn reviewed for strength loss rather than just greenness: `grandfatherFlatRooms` adds a tracked `<slug>/review/.gitkeep`, and the two assertions that previously read "no `<slug>/` companion exists" now read "no `<slug>/review/<stage>` room exists" — the necessary narrowing once the fixture must hold a companion. `TestPrepareRejectsSymlinkedFlatCompanionWithoutChangingBytes` is not short-circuited by the guard: it creates `external/review` so the guard grandfathers, and still asserts the error contains `symlink`.

### Recommendation

**PASSED.** All five ACs carry valid, independently reproduced evidence and no material finding remains. The two polish findings above are non-blocking and neither consumes a feedback cycle.

### Summary

Reproduced all five ACs on a throwaway clone rather than reading them off the report, and re-derived every live number independently: 9 hybrids, 34 refs, 0 unresolved, exactly 9 added `--validate` lines, plain `status` byte-identical. AC-2's byte-identity claim was re-verified from room contents — 66 retained files, 0 changed, 0 vanished, 6 added — after reproducing all three live prepares at their exact claimed attempt numbers.

The dispatch's armed-gate-state concern is resolved: those prepares never touched live state. Live `define-fo` still tops out at `validation/briefing-5`, `preserve-pi` at `validation/briefing-1` (superseded), `merge-guard` holds no validation room, and the state checkout is clean. The implementation ran them in clones, and the attempt numbers it reported are exactly what live state would next mint — which corroborates that those clones were faithful rather than contradicting the report.

On the crash-inheritance chain the concern was justified but narrow. Ten of eleven claimed falsifiers red exactly their stated case under mutation, so the AC-to-test bindings are real rather than asserted. Only two hunks rested on "tests pass" alone: the `--round` message, which no test asserts and which I exercised directly, and AC-1's lock-placement clause, which no test can red and which needs none.

## Stage Report: implementation (cycle 2)

- DONE: Gate the refusal on a workflow-declared form so a workflow that declares nothing allows either shape
  `8849dc880`. `refuseNewFlatCompanion` returns early unless the README declares `entity-form: folder`; `workflowDeclaresFolderForm` reads that one key through the same `frontmatterNode`/`mappingValue` path `entityIdentity` already uses for `id-style`, and treats an unreadable README as undeclared so the guard adds no new failure mode.
- DONE: Prove the live-fixture shape now prepares: a fresh flat entity in an undeclared workflow succeeds, and a flat entity in a folder-declaring workflow is refused
  Four cases run against the built binary on the live auto-continue fixture's own README bytes. Undeclared: exit 0, binds `./auto-continue-task/review/validation/briefing-1`. The cycle-1 binary on the identical fixture: exit 1 with the message the journeys quoted. Declared: exit 1, tree entry list unchanged, no companion minted. Declared with rooms already present: exit 0.
- DONE: Rerun the targeted live journey locally against the built change before signalling, and report its observed codes
  `TestLiveCommonAutoContinueAfterImplementation`, claude lane / sonnet, built from this branch: **PASS 393.04s, `go test` exit 0**. Codes and the one red I also observed are in `### Live evidence` below.

### Grandfathering: kept, with a changed reason

Reconsidered as instructed, and kept. Under the declaration it no longer protects the nine hybrids from a rule they never opted into; it protects them from one a workflow opts into later, which is the concrete next step here — `docs/dev` is both the likeliest first workflow to declare `entity-form: folder` and the one holding all nine hybrids, three at active validation gates. Without the branch, adding one README line arms a gate failure that fires at each hybrid's next gate: green now, blocked later, which is #739's shape. It costs four production lines and one test case, because gating on the declaration removed all of its fixture churn. Full reasoning in `## Cycle 2 narrowing`.

### Falsifiers executed

Five source mutations, each re-running `TestPrepareRefusesFlatEntityOnlyWhereTheWorkflowDeclaresFolderForm`. Every one reds exactly one branch and only that branch, so the three branches are independent rather than co-satisfied.

- **AC-1 refusal.** Guard deleted -> `declared-folder-form-refuses-flat` FAILS. Key read as `entity-shape` -> FAILS. Declared value compared against `flat` -> FAILS. The last two matter because either makes a declaring workflow read as undeclared, which is the cycle-1 defect inverted.
- **AC-1 permission.** Declaration check deleted, so the guard fires unconditionally -> `undeclared-workflow-prepares-flat` FAILS. This is the exact regression that reddened both live lanes; it is now a red offline test.
- **AC-2.** Grandfathering early-return deleted -> `declared-folder-form-grandfathers-rooms` FAILS.
- **AC-3/4/5 untouched.** The two `status --validate` findings are unchanged from cycle 1 and keep their validated falsifiers. Re-confirmed live read-only: 9 hybrid findings naming 9 slugs, 0 unresolved refs, `--validate` exit 0.

### Live evidence

- **`TestLiveCommonAutoContinueAfterImplementation` (the mandated journey) — PASS, 393.04s, exit 0.** Both fixture variants drove `gate prepare` on flat `auto-continue-task.md` and bound `briefing:auto-continue-task:validation:attempt-1:revision-1` with `state=open`. Occurrences of `requires folder-form entity` in either transcript: **0**.
- **`TestLiveCommonRejectionFlow` (extra evidence, the other flat-entity gate fixture) — one RED, not attributable, and resolved.** Run 1 on this branch FAILED 356.67s on grade `rejection-worker-topology`: "the fresh branch owes 8 routing events, the run produced 6", with `reuse` where `spawn` was owed. `gate prepare` had already succeeded in that same run — it bound `briefing:rejection-task:validation:attempt-1:revision-1` with 0 refusals — so the guard was not what failed; the FO under-drove the rejection chain by two events. Settled by two further runs rather than by argument: the same journey at stack parent `8f29ad577` with no change present PASSED 496.36s on the full 8-event chain, and a rerun on this branch PASSED 432.68s. This branch changes 0 files under `skills/`, `internal/dispatch/`, or `internal/ensigncycle/`, which is where the FO's fresh-versus-reuse routing is decided.

### Surface

8 files, +331 / -7, **net +324**, against a declared +283 ±40 LOC and 13 ±2 files. Both bounds are missed and the estimate predates the captain's narrowing, so it is stale rather than met: files land 3 under the low bound because the declaration removed the fixture churn (6 files, 34 repairs, all reverted to their stronger pre-change assertions), and net LOC lands 1 over the high bound because `prepare.go` grew from +21 to +53 for the declaration reader and the guard's reasoning, `prepare_test.go` from +76 to +90 for the third branch, and docs from +9 to +11 for the added concepts page. Production alone is +99/-1. Flagging the deviation rather than reconciling it: the shape of the change moved when the rule did.

### Checks run

`go test ./... -race`: **19 packages green, 0 data races**. `gofmt -l ./cmd ./internal` clean. The only failure is `TestCodexResolveManifestAgainstInstalledHost`, which I confirmed fails identically at stack parent `8f29ad577` without this change — machine-local (a locally installed Codex plugin), not attributable.

### Summary

The guard now fires only where a workflow declares `entity-form: folder` in its README, so the shape `spacedock new` mints by default is refused nowhere and enforcement arrives with the declaration instead of ahead of it. Preparation only reads the key; `qpa` owns writing it. The clearest evidence that cycle 1's defect is closed is the before/after on the live fixture's own README bytes: the cycle-1 binary exits 1 there with the message the journeys quoted, this one exits 0 and binds the room, and the mandated live journey passes end to end with zero refusals in either variant's transcript.

Two things a reviewer should weigh. First, the estimate is missed in both directions and I did not reconcile it — the narrowing deleted 34 fixture repairs and added a frontmatter reader, so 8 files/+324 is the honest measurement against a stale +283 ±40/13 ±2. Second, `gate record --round` still refuses every flat entity unconditionally at `operation.go:108`. That is pre-existing and untouched here except for cycle 1's remedy text, and no live journey drives `--round` on a flat entity today, but it does leave the two commands disagreeing about whose rule the form is: prepare now asks the workflow, `--round` still asserts it. Whether `--round` should take the same declaration is a scope decision I did not make.
