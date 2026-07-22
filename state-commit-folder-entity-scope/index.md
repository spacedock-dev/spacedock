---
id: vn15pvn4nt9zce55m757f23a
title: Make state commit include folder-form entity artifacts
status: validation
source: "Captain-directed follow-up from Roborev setup entity 00 dogfood, 2026-07-14"
started: 2026-07-14T13:28:29Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-state-commit-folder-entity-scope
issue:
milestone: 0.26.0
sprint: durable-decisions
group: recorder
gates:
    version: 1
    current:
        gate: gate:state-commit-folder-entity-scope:validation
    records:
        - id: gate:state-commit-folder-entity-scope:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:state-commit-folder-entity-scope-ideation-1
              briefing:
                id: briefing:docs-dev:vn:ideation:canonical-v1:revision-1
                digest: sha256:8191dc40c67c3854119c189622f6a15006e8a21fa9ac7ce2ce6f0618a66f496d
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-v1
              resolution:
                type: Resolution
                id: resolution:spacedock:state-commit-folder-entity-scope:ideation:1
                briefing: briefing:docs-dev:vn:ideation:canonical-v1:revision-1
                by: person:captain
                at: "2026-07-22T14:19:30.731531Z"
                decision: approve
                reason: 'Captain explicitly confirmed in chat after the recovered Subspace advisory approval: vn is approved.'
        - id: gate:state-commit-folder-entity-scope:validation
          stage: validation
          attempts:
            - id: gate-attempt:state-commit-folder-entity-scope-validation-1
              briefing:
                id: briefing:docs-dev:vn:validation:canonical-v1:revision-1
                digest: sha256:af40be6a0218d6334de45aaec07c9e7d37777aa849b31c049cea93d44d93d31b
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:state-commit-folder-entity-scope:validation:1
                briefing: briefing:docs-dev:vn:validation:canonical-v1:revision-1
                by: agent:first-officer
                at: "2026-07-22T15:55:30.121749Z"
                decision: approve
                reason: All six acceptance criteria reproduced; Roborev material findings fixed; full, race, focused, format, and clean checks passed; live xb and vn gate-room commits proved exact folder scoping.
                adoption-note: 'Captain: you have the conn toward the sprint goal; authorized to approve gates, PR, relevant CI lanes, and merge; use your judgement.'
---

`spacedock state commit <slug>` commits a flat entity correctly but treats a folder-form entity as only `<slug>/index.md`. Reports, evidence, and artifacts stored beside the index remain dirty even though the state command reports that the entity was committed and pushed.

## Observed dogfood failure

Amending folder-form entity `roborev-workflow-setup-skill` changed its `index.md` and two canonical files below `artifacts/`. Running:

```bash
spacedock state commit roborev-workflow-setup-skill
```

created and pushed state commit `d55074bb` containing only `roborev-workflow-setup-skill/index.md`. Both associated artifact files remained modified in the state checkout.

The implementation explains the outcome: folder-form resolution returns `<slug>/index.md`, and `commitEntityPathScoped` uses that single file as its Git pathspec. This conflicts with the storage contract, which makes reports and artifacts below a folder-form entity part of the workflow record. Archive finalization already scopes folder-form entities to the whole `<slug>/` directory.

A manual cleanup also exposed a second boundary: `state commit` accepted nested pseudo-slugs naming each artifact without a `.md` suffix and committed them. A command whose operand is documented as an entity slug must not accept path-bearing aliases for arbitrary nested Markdown files.

The failure reproduced again while dogfooding the canonical v1 gate recorder on 2026-07-22. The First Officer added a new retained package at `durable-gate-approval-pending-blockers/review/validation/briefing-v1/` and ran:

```bash
spacedock state commit --workflow-dir docs/dev durable-gate-approval-pending-blockers
```

The command reported `Nothing to commit ... state checkout already up to date` while both new package files remained untracked. The exact-path fallback committed and pushed only those files as state commit `d8e4180c`; the worktree-built recorder then bound the package successfully, and the ordinary entity-only mutation committed as `2c616b7e`. This confirms both halves of the boundary on a second folder-form entity: `index.md` mutations work, but additional reports and artifacts are omitted and can trigger a false clean no-op.

## Proposed approach

- Resolve one entity commit unit from the canonical top-level slug:
  - flat form commits exactly `<slug>.md`;
  - folder form commits exactly `<slug>/`, including tracked changes and new reports or artifacts below it.
- Preserve the concurrency boundary. Never stage the state checkout wholesale; dirty sibling entities and unrelated state files remain untouched.
- Treat concurrent changes anywhere within one folder-form entity as changes to the same entity. Preserve the existing reject, rebase, conflict-HALT, and no-force behavior.
- Validate the operand as a canonical entity slug before path resolution. Reject path separators, traversal, and nested pseudo-slugs without changing the index, worktree, or remote.
- Keep flat-form output, exit codes, JSON results, no-origin behavior, and clean no-op behavior compatible.

## Approved implementation boundary

Expected surface: `internal/cli/state_sync.go` (~20-60 production LOC), existing real-Git/CLI state-commit tests plus focused new cases (~160-300 test LOC), and the state-command reference (~5-20 documentation lines), with 2× tolerance. No new package, dependency, command, persisted schema, or background process is expected.

No new spike is needed. The riskiest mechanism is already proven twice by live failure and by source inspection: folder-form resolution returns `<slug>/index.md`, and the existing path-scoped Git helper stages only that resolved path. Archive handling already demonstrates the repository's directory-scoped entity boundary. Implementation should change the resolved commit unit, not introduce another lifecycle.

Mechanism-to-value check: a folder-scoped Git pathspec is the smallest mechanism that makes one folder entity durable while preserving sibling isolation. Checkout-wide `git add -A` cannot protect sibling dirt; enumerating known artifact types would miss future valid room files; an artifact manifest or registry would duplicate the filesystem boundary and add a second protocol.

User-facing documentation change: the state command reference must say that a flat entity commits only `<slug>.md`, while a folder entity commits every changed/new/deleted non-ignored path below `<slug>/`; unrelated paths remain untouched. It must also say the operand is one canonical top-level entity slug, not a nested path.

## Acceptance criteria

**AC-1 - A folder-form state commit durably includes the entity index and all changed or newly created reports and artifacts below that entity.**
Verified by: a real-Git fixture changes `<slug>/index.md`, modifies a tracked report, and creates an untracked artifact. One `spacedock state commit <slug>` must exit 0, create and push one entity-scoped commit containing all three paths, and leave no dirty path below `<slug>/`.

**AC-2 - A folder-form state commit never sweeps dirty sibling entities or unrelated state-checkout paths.**
Verified by: the AC-1 fixture also dirties a flat sibling, a folder sibling, and a top-level untracked file. The committed name set must equal only paths below the target folder; all sibling dirt must remain present and unstaged.

**AC-3 - Nested artifact dirt prevents a false clean no-op even when `index.md` is unchanged.**
Verified by: change only `<slug>/artifacts/evidence.md`, run the command, and assert a new pushed commit contains that path. Re-running without further changes must return the existing clean no-op result.

**AC-4 - Flat-form behavior remains compatible.**
Verified by: existing state-commit real-Git tests plus an exact committed-path assertion for `<slug>.md`; output/JSON, no-origin, retry-on-reject, and conflict-HALT tests remain green.

**AC-5 - Concurrency remains entity-scoped for folder form.**
Verified by: two-host real-Git fixtures show disjoint entities rebase and push successfully, while concurrent edits within the same folder-form entity halt on conflict, name the conflicting path, abort the rebase cleanly, and never force-push or discard either writer.

**AC-6 - The command rejects noncanonical and path-bearing operands without side effects.**
Verified by: table-driven CLI and real-Git cases for slash and platform-separator nesting, `.`, `..`, traversal, absolute paths, and the dogfood pseudo-slug `roborev-workflow-setup-skill/artifacts/roborev-setup-skill/SKILL`. Each must exit nonzero with a clear invalid-slug diagnostic and leave HEAD, index, worktree bytes, and remote unchanged.

## Test plan

1. Add the AC-1/AC-3 real-Git regression first and prove current code leaves the nested files dirty or returns a false no-op.
2. Change folder-form resolution to return the folder as the commit pathspec while retaining the single-file pathspec for flat form.
3. Extend the sibling-dirt and two-host concurrency fixtures across both entity forms.
4. Add operand-validation cases before filesystem resolution and prove the nested pseudo-slug workaround is rejected.
5. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Out of scope

- Committing multiple independent entities in one invocation.
- Sweeping the state checkout with `git add -A`.
- Changing folder-form discovery, archive layout, or merge-guard finalization semantics.
- Moving product deliverables into the state checkout.

## Stage Report: ideation

- DONE: Capture the live folder-form artifact omission and the accepted nested pseudo-slug workaround as durable reproduction evidence.
- DONE: Bound the intended commit unit so folder artifacts become durable without weakening sibling-entity isolation or multi-writer conflict safety.
- DONE: Define regression coverage for folder inclusion, false no-op, flat compatibility, sibling dirt, concurrency, and operand validation.

### Summary

The task is ready for staff review as a 0.26 correctness fix. The key design is directory-scoped Git pathspecs for folder-form entities, not checkout-wide staging; canonical slug validation closes the workaround that exposed the gap.

## Stage Report: implementation

- DONE: Change `state commit <slug>` so flat-form entities stage exactly `<slug>.md` while folder-form entities stage exactly `<slug>/`, including tracked modifications/deletions and new non-ignored reports/artifacts.
  Code commit `e0822912`; `TestStateCommitFolderIncludesWholeEntity` fails if the folder pathspec omits the index, tracked report edit/deletion, or new artifact.
- DONE: Preserve strict entity scoping: never stage dirty sibling entities or unrelated top-level state paths, and keep existing local-only, push, reject/rebase, conflict-HALT, no-force, JSON/text, and clean no-op behavior.
  The same real-Git test leaves flat/folder/top-level sibling dirt untouched; the full CLI suite would fail on output, retry, HALT, local-only, or no-op regressions.
- DONE: Reuse the existing canonical entity discovery/slug authority where possible; reject path-bearing, traversal, absolute, dot, and nested pseudo-slug operands before side effects without minting a conflicting slug grammar.
  `status.EntitySlug` identifies resolved commit units; `TestStateCommitRejectsNoncanonicalSlugWithoutSideEffects` proves HEAD, index, worktree bytes, and origin remain unchanged for every rejected operand.
- DONE: Add the real-Git regression first: index + tracked report + untracked artifact commit together; sibling dirt remains; artifact-only change is not a false no-op; clean rerun is a no-op.
  The regression was run RED before production edits (only `folder-task/index.md` committed), then GREEN; it also pins a tracked deletion in the same entity unit.
- DONE: Extend two-host real-Git coverage so disjoint entities rebase/push and conflicting paths anywhere within the same folder entity halt cleanly without force or discard.
  `TestStateCommitFolderMultiWriterHappyPath` fails if either host's artifact misses origin; `TestStateCommitFolderConflictHalts` fails unless both the peer and local nested edits survive the aborted rebase boundary.
- DONE: Preserve flat-form compatibility with exact committed-path assertions and existing state-commit output/error behavior.
  `TestStateCommitIsPathScoped` now requires the committed path set to equal only `first-task.md`; all pre-existing state-commit tests pass unchanged.
- DONE: Update the state command reference with the approved flat/folder commit-unit and top-level-slug wording; add no new command, package, dependency, schema, registry, or background lifecycle.
  `docs/site/reference/command-reference.md` documents the canonical operand and exact commit units; commit `e0822912` adds no new protocol surface.
- DONE: Run gofmt -w ./cmd ./internal, focused state-commit tests, go test ./..., and go test ./... -race; record exact evidence, actual surface versus estimate, and commit the clean worktree.
  All commands passed; actual surface was 36 production changed lines vs 20-60 estimated, 218 test lines vs 160-300, and 9 documentation lines vs 5-20.

### Summary

Commit `e0822912` makes a folder-form entity one durable Git pathspec while retaining the single-file boundary for flat entities and rejecting noncanonical path aliases before resolution. Real-Git coverage now proves complete folder artifacts, sibling isolation, artifact-only commits, clean no-ops, flat compatibility, and both disjoint and conflicting two-host behavior; the full normal and race suites pass.

### Feedback Cycles

- Cycle 1: REJECTED — Roborev `branch_final` job 534; surface 3 files/263 LOC vs estimate 3 files/185-380 LOC (93% of 283-line midpoint); AC unchanged
- Cycle 2: REJECTED — Roborev `branch_final` job 537; surface 3 files/337 LOC vs estimate 3 files/185-380 LOC (119% of 283-line midpoint); AC unchanged
- Cycle 3: REJECTED → DECLINED/DEFERRED — Roborev `branch_final` job 540; surface 3 files/429 LOC vs estimate 3 files/185-380 LOC (152% of 283-line midpoint); AC unchanged — stale comment fixed in `d4d39ed6`; Commander retained exact current-form scope and deferred conversion until a supported entity-form conversion workflow exists

## Stage Report: implementation (cycle 2)

- DONE: Change `state commit <slug>` so flat-form entities stage exactly `<slug>.md` while folder-form entities stage exactly `<slug>/`, including tracked modifications/deletions and new non-ignored reports/artifacts.
  Commits `e0822912` and `f58dfe20`; real-Git tests fail if folder contents, index deletion, complete-folder deletion, or flat deletion are omitted.
- DONE: Preserve strict entity scoping: never stage dirty sibling entities or unrelated top-level state paths, and keep existing local-only, push, reject/rebase, conflict-HALT, no-force, JSON/text, and clean no-op behavior.
  Commit `030b9c3c` makes every scoped Git operation literal; the explicit `:(glob)` regression fails if matching tracked/untracked siblings reach the local commit or remote.
- DONE: Reuse the existing canonical entity discovery/slug authority where possible; reject path-bearing, traversal, absolute, dot, and nested pseudo-slug operands before side effects without minting a conflicting slug grammar.
  Canonical on-disk resolution retains exact current form, with tracked-path fallback only for deletion; invalid-operand tests prove HEAD, index, bytes, and origin unchanged.
- DONE: Add the real-Git regression first: index + tracked report + untracked artifact commit together; sibling dirt remains; artifact-only change is not a false no-op; clean rerun is a no-op.
  `TestStateCommitFolderIncludesWholeEntity` was RED before production edits and is GREEN with exact committed-path and residual-dirt assertions.
- DONE: Extend two-host real-Git coverage so disjoint entities rebase/push and conflicting paths anywhere within the same folder entity halt cleanly without force or discard.
  Folder disjoint-writer and nested-conflict tests prove linear rebase/push, exit-3 abort, named nested path, and preservation of both writers.
- DONE: Preserve flat-form compatibility with exact committed-path assertions and existing state-commit output/error behavior.
  Flat update and deletion tests require exactly `first-task.md`; the complete existing state-commit suite passes.
- DONE: Update the state command reference with the approved flat/folder commit-unit and top-level-slug wording; add no new command, package, dependency, schema, registry, or background lifecycle.
  The final four-commit branch changes only the assigned three files and retains the folder pathspec as the smallest mechanism.
- DONE: Run gofmt -w ./cmd ./internal, focused state-commit tests, go test ./..., and go test ./... -race; record exact evidence, actual surface versus estimate, and commit the clean worktree.
  All commands passed after the last semantic change; final surface is production 62 lines vs 20-60, tests 363 vs 160-300, docs 9 vs 5-20, within the approved 2x tolerance; code worktree is clean at `d4d39ed6`.

### Roborev disposition

- Job `534` (`correctness F`, `product P`, synthesis `F`): one medium tracked-deletion finding, fixed by `f58dfe20` with index-only, complete-folder, and flat-deletion regressions.
- Job `537` (`correctness F`, `product P`, synthesis `F`): one high literal-pathspec isolation finding, fixed by `030b9c3c` with explicit Git-magic sibling and nonexistent-alias regressions.
- Job `540` (`correctness F`, `product P`, synthesis `F`): low stale-comment finding fixed by `d4d39ed6`; medium flat↔folder conversion request declined/deferred by Commander because dual-form staging conflicts with the approved exact current-form unit. Promote when a supported entity-form conversion workflow exists.

### Summary

The final implementation makes folder artifacts and tracked deletions durable without weakening exact entity isolation, including for canonical filenames containing Git pathspec syntax. All material in-scope Roborev findings are fixed; the sole semantic expansion is explicitly deferred under the Commander ruling, with its promotion trigger recorded above.

## Stage Report: validation

- DONE: Reproduce that one folder-form state commit includes the index plus tracked deletions/changes and new room artifacts, never false-no-ops on artifact-only changes, and leaves sibling/top-level dirt untouched.
  `TestStateCommitFolderIncludesWholeEntity` passed against real Git; its exact four-path assertion and residual porcelain checks fail on any omitted target path, false no-op, or swept sibling/top-level path.
- DONE: Attack exact scoping with Git pathspec-magic names, invalid/path-bearing operands, flat-form exact-file compatibility, and two-host disjoint/conflicting folder writes without force or discarded bytes.
  Focused real-Git tests passed; their assertions compare exact committed names, HEAD/index/worktree/origin bytes, linear history, exit 3, clean rebase abort, and both hosts' surviving content.
- DONE: Verify the final 62-production/363-test/9-doc surface, Roborev 534/537 fixes and 540 deferred classification, plus focused/full/race/format/cleanliness evidence against the committed branch.
  `07811d27..d4d39ed6` is exactly 62 production, 363 test, and 9 documentation changed lines in the approved three files; gofmt and `git diff --check` left clean HEAD `d4d39ed6`.
- DONE: AC-1 - A folder-form state commit durably includes the entity index and all changed or newly created reports and artifacts below that entity.
  `TestStateCommitFolderIncludesWholeEntity` committed and pushed exactly index, tracked report change, tracked deletion, and new artifact, then observed no target-folder dirt.
- DONE: AC-2 - A folder-form state commit never sweeps dirty sibling entities or unrelated state-checkout paths.
  The same test retained flat sibling, folder sibling, and top-level dirt; `TestStateCommitTreatsSlugAsLiteralGitPathspec` also kept matching tracked/untracked siblings off the commit and origin.
- DONE: AC-3 - Nested artifact dirt prevents a false clean no-op even when `index.md` is unchanged.
  The artifact-only phase advanced HEAD with exactly `folder-task/artifacts/evidence.md`; an unchanged rerun returned the established JSON `no-op` result.
- DONE: AC-4 - Flat-form behavior remains compatible.
  `TestStateCommitIsPathScoped` and `TestStateCommitFlatDeletion` required exactly `first-task.md`; JSON/text HALT, no-origin, retry/rebase, and clean-no-op coverage passed in focused/full suites.
- DONE: AC-5 - Concurrency remains entity-scoped for folder form.
  `TestStateCommitFolderMultiWriterHappyPath` proved disjoint artifacts reach linear origin history; `TestStateCommitFolderConflictHalts` proved named nested conflict, clean abort, no force, and preservation of both writers' bytes.
- DONE: AC-6 - The command rejects noncanonical and path-bearing operands without side effects.
  Eight real-Git cases, including both separators, dot/traversal/absolute forms, and the dogfood pseudo-slug, preserved HEAD, index, worktree bytes/status, and remote while returning the invalid-slug diagnostic.
- DONE: Roborev material findings and deferred risk classification.
  Job 534 is closed by index-only/whole-folder/flat deletion tests; job 537 is closed by literal-pathspec and nonexistent-alias tests. Job 540 remains deferred only for simultaneous flat↔folder conversion, which is unsupported and canonically diagnosed as a conflict; promote when a supported conversion workflow exists, and reject any dual-form staging expansion before then.
- DONE: Focused, full, race, format, and cleanliness verification.
  Focused state-commit matrix and canonical dual-form validator passed; `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed with no worktree dirt.
- DONE: PASSED recommendation.
  All six value ACs have executable Git/state evidence, no material finding remains, and the only deferred risk has an unsupported trigger plus a concrete promotion condition.

### Summary

Fresh validation reproduced every promised folder, flat, operand, deletion, and two-host boundary against real Git at `d4d39ed6`. Recommendation: PASSED; job 540's form-conversion request remains a deferred risk until the product supports conversion, and no hidden dual-form compatibility expansion is accepted.
