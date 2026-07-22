---
id: vn15pvn4nt9zce55m757f23a
title: Make state commit include folder-form entity artifacts
status: ideation
source: "Captain-directed follow-up from Roborev setup entity 00 dogfood, 2026-07-14"
started: 2026-07-14T13:28:29Z
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
sprint: durable-decisions
group: recorder
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
