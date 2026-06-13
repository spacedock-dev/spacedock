---
id: ndpfbqqvezggnrydnvrxjmh2
title: Prefer `spacedock new` over manual --next-id for filing entities
status: backlog
source: captain (2026-06-13) — noticed FO still calls --next-id when filing seeds
started:
completed:
verdict:
score:
worktree:
issue:
---

Make `spacedock new` the contract-blessed atomic-create path for filing entities, and have `status --next-id` emit a hint pointing to it.

## Problem

`spacedock new [--folder] SLUG` (alias `status --new`) has existed since #242 (2026-05-31): it mints the id and atomically writes a valid stamped entity in one call. But the FO operating contract never adopted it — both `first-officer-shared-core.md` (Status/Dispatch sections) and `claude-first-officer-runtime.md` still teach the manual flow: call `status --next-id`, hand-assemble the frontmatter, `Write` the file. The result: every FO files entities the slow, drift-prone way (a `--next-id` candidate is not a reservation, so the fetched id can drift from what finally lands), and operators get no signal the atomic path exists.

## Proposed approach

Two coupled changes:
1. **Contract** (`skills/first-officer/references/first-officer-shared-core.md` + `claude-first-officer-runtime.md`): teach `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < body` as the blessed atomic-create path for seed filing. Keep `--next-id` documented only for its candidate-preview use. Note that for split-root state checkouts the FO still does the path-scoped commit + push after `new` (new writes, it does not commit).
2. **Binary** (`internal/status`): `status --next-id` emits a stderr hint that `spacedock new` files atomically (so an operator reading next-id output is pointed at the better path). Hint must not pollute the stdout id value that callers parse.

## Out of scope

Changing what `new` itself does (it already mints + atomic-writes correctly). Auto-committing from `new` (commit/push stays the caller's, per split-root concurrency-safety rules).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - `status --next-id` emits a use-`new` hint on stderr without altering the stdout id.**
Verified by: a Go test in `internal/status` asserting stdout is exactly the id (unchanged from current) and stderr contains the `new` hint; exit code unchanged.

**AC-2 - The FO contract teaches `spacedock new` as the atomic-create path and a live drive shows the FO filing via `new`, not the manual --next-id+Write flow.**
Verified by: a live first-officer drive (the ensigncycle shared-scenario harness or an equivalent recorded run) in which the FO files a new entity and the durable result shows an entity minted via `new`; the prose change alone does not satisfy this AC.

## Test plan

AC-1 is a cheap Go unit test (stdout/stderr separation). AC-2 needs a live drive that observes the FO's filing behavior — the expensive half; reuse the existing live ensigncycle scaffolding rather than building new. Contract prose edits are authoring work but are not an AC on their own.
