---
id: f1tqvtmr5xeqh6j3rg3qa6hk
title: migration-check test shares the production walk composition (extract a walk-step helper)
status: implementation
source: "0198 pre-cut antipattern audit, finding R1 (2026-06-08). kb's TestMigrationCheckPrunesStateTree (migration_check_test.go:213-227) re-implements the filepath.Walk callback instead of sharing the production walk composition (the info.IsDir()→SkipDir step). It DOES exercise the real shared predicate (isMigrationCheckPrunedDir), so this is not a hole — but if the walk composition ever diverges from the production path, the hermetic test would not catch it. Filed for 0.19.9 per captain."
started: 2026-06-08T20:43:10Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-migration-check-share-walk-helper
issue:
group: test-hygiene
sprint: 0199-pre-flip-mechanics
sprint-readiness: ready
---

Extract the migration-check walk-step composition into a shared helper that both the production check and the hermetic prune test call, so the prune composition is tested once instead of re-implemented in the test.

## Problem

`TestMigrationCheckPrunesStateTree` inlines a `filepath.Walk` callback (`migration_check_test.go:213-227`) that re-implements the production walk's `info.IsDir()` → `SkipDir` composition. The test exercises the real shared predicate (`isMigrationCheckPrunedDir`), which is the load-bearing prune logic, so the prune *decision* is genuinely tested. But the *walk composition* (how the predicate is applied during the walk) is duplicated, not shared — if production's composition diverges from the test's inlined copy, the hermetic test would stay green on a broken production walk.

## Proposed approach

Extract the walk-step (the `info.IsDir()`/`SkipDir` composition around `isMigrationCheckPrunedDir`) into a shared helper in the production package; have both the production migration-check and the test drive it. Net: the prune composition is defined and tested once.

## Out of scope

- The prune predicate itself (`isMigrationCheckPrunedDir`) — already shared and tested; this is only the walk-step composition around it.

## Acceptance criteria

Ideation/implementation fills in (small — likely fast-track). Sketch:

- The production migration-check and `TestMigrationCheckPrunesStateTree` both drive ONE shared walk-step helper (no inlined re-implementation in the test). Verified by: the test no longer contains its own `filepath.Walk` callback composition; the existing prune assertions (`.spacedock-state` subtree never visited AND the real entity IS visited) stay green against the shared helper.

## Test plan

The existing `migration_check_test.go` assertions carry over against the shared helper. `go test ./internal/status -run TestMigrationCheck` green. Small, mechanical; candidate for fast-track (no gate).
