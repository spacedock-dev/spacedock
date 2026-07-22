---
id: vn15pvn4nt9zce55m757f23a
title: Make state commit include folder-form entity artifacts
status: implementation
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
        gate: gate:state-commit-folder-entity-scope:ideation
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
