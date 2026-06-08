---
id: f1tqvtmr5xeqh6j3rg3qa6hk
title: migration-check test shares the production walk composition (extract a walk-step helper)
status: validation
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
mod-block:
pr: "#331"
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

## Stage Report: implementation

- DONE: Extract the migration-check walk-step composition (info.IsDir() → filepath.SkipDir wrapping isMigrationCheckPrunedDir) into ONE shared helper in the production internal/status package; both the production migration-check AND TestMigrationCheckPrunesStateTree drive that single helper; the test no longer inlines its own filepath.Walk callback composition.
  New production file internal/status/migration_check.go defines migrationCheckWalkDir + the moved isMigrationCheckPrunedDir; both walk callbacks (migration_check_test.go:69 consistency walk, :199 prune walk) call it. Commit ec44dc05 on branch spacedock-ensign/migration-check-share-walk-helper.
- DONE: The existing prune assertions stay green against the shared helper (.spacedock-state subtree never visited AND the real entity IS visited).
  go test ./internal/status -run TestMigrationCheckPrunesStateTree → 1 passed. Full package: 447 passed, 1 skipped, 1 pre-existing unrelated FAIL.

### Summary
Both walk callbacks duplicated the same info.IsDir()→isMigrationCheckPrunedDir→filepath.SkipDir step, and the predicate lived only in the test file. Moved the predicate to a new production file internal/status/migration_check.go and added migrationCheckWalkDir wrapping it; both the live consistency walk and the hermetic prune test now drive the one shared walk-step, so the prune composition is defined and tested once. The pre-existing RED TestMigrationCheckFixturesParseConsistently (docs/roadmap/0198-pre-flip-hardening/debrief.md unquoted session-date scalar, migration_check_test.go:115) is the documented separate item owned elsewhere — not touched, not introduced by this change (the walk still visits 58 frontmatters and the .spacedock-state prune assertion does not fire).

## Stage Report: implementation (cycle 2)

Scope addition (captain-approved, relayed by team-lead): also FIX the pre-existing migration-check RED (TestMigrationCheckFixturesParseConsistently), folding the prune widening into the shared walk-step. The earlier "do not fix" instruction was reversed.

- DONE: Primary deliverable unchanged — the shared walk-step refactor (migrationCheckWalkDir + isMigrationCheckPrunedDir in production internal/status/migration_check.go; both walk callbacks drive it; no inlined composition in the test).
  Commit ec44dc05; helper signature widened to take path in 8e72b5f2.
- DONE: Root-cause fix for the RED — widen the shared prune predicate so the docs/roadmap strategy tree (non-entity by design) is pruned wholesale. Path-aware (a `roadmap` dir directly under a `docs` dir) because `roadmap` is a generic basename that must be qualified by its parent, unlike the reserved `.spacedock-state` name. NOT a date-quote band-aid.
  isMigrationCheckPrunedDir now also skips a `roadmap` dir directly under `docs`. The debrief's bare-YAML session-date scalar (decodes as time.Time directly but string through the reader) is the same non-entity shape as the already-pruned .spacedock-state _debriefs. Commit 8e72b5f2.
- DONE: NEW AC — go test ./internal/status FULLY GREEN; BOTH TestMigrationCheckPrunesStateTree AND TestMigrationCheckFixturesParseConsistently pass.
  448 passed, 0 failed, 0 skipped (was 447/1 fail/1 skip). Full repo: go test ./... → 1172 passed, 16 packages.
- DONE: TDD — confirmed the RED first (migration_check_test.go:115 session-date mismatch), then greened it; extended the hermetic prune test to plant a docs/roadmap debrief and assert it is pruned. Proved non-vacuous: reverting the roadmap branch reds TestMigrationCheckPrunesStateTree at the new docs/roadmap assertion.

### Summary
The pre-existing RED was a scope-creep of the migration-check walk into docs/roadmap — the strategy layer, which holds session debriefs (same session-date non-entity shape as the pruned .spacedock-state _debriefs) but no entity frontmatter. Fixed at root by widening the SHARED prune predicate (the exact composition this task centralized) to skip docs/roadmap wholesale — path-aware (a `roadmap` dir under `docs`) because the basename is generic, rather than quoting the one date. The hermetic prune test now covers the new branch (planted roadmap debrief, asserted pruned, non-vacuity confirmed). go test ./internal/status fully green; full repo green.

## Stage Report: implementation (pre-merge polish + rebase)

- DONE: Fix the validator's polish finding (no-false-comments rule) — drop the false precedent citations from the migration_check.go comment AND the cycle-2 summary above. discoverIgnoreDirs prunes .spacedock-state but NOT roadmap; boundary_guard_test.go lives in internal/contractlint, not this package. The docs/roadmap prune is the only one of its kind here, now stated self-containedly (both trees hold non-entity frontmatter the migration check does not govern). Code-comment commit ffca15e5 (pre-rebase SHA).
- DONE: Rebase the worktree branch onto origin/next — clean, no conflicts. origin/next advanced by exactly one file (docs/roadmap/0199-pre-flip-mechanics/staff-review.md), which my prune already skips; zero overlap with my two migration-check files. Branch is now linear on origin/next 84b2005b; my 3 commits reapplied as 89375a42 / d9d4b31a / 6545a420.
- DONE: Re-run go test ./... on the rebased branch — 1172 passed, 16 packages, 0 fail. Both TestMigrationCheckPrunesStateTree and TestMigrationCheckFixturesParseConsistently green; gofmt + go vet clean. Code branch NOT pushed (team-lead opens the PR at the merge ceremony).

## Stage Report: validation

- DONE: Shared walk-step refactor — confirm the production migration-check AND TestMigrationCheckPrunesStateTree both drive the ONE shared helper (migrationCheckWalkDir + isMigrationCheckPrunedDir in internal/status/migration_check.go) — no inlined filepath.Walk composition remains in the test. go test ./internal/status -run TestMigrationCheckPrunesStateTree green.
  Both walk callbacks (migration_check_test.go:69 live consistency check = the "production migration-check"; :214 hermetic prune test) call migrationCheckWalkDir; grep over the test shows ZERO inlined info.IsDir()/SkipDir composition. Adversarial edit 4 (disabled the prune in the production helper) reds BOTH tests → they genuinely drive the one shared production helper, not a copy. `go test ./internal/status -run TestMigrationCheck` → 2 passed (57 frontmatters verified).
- DONE: RED fix (cut-critical) — go test ./internal/status FULLY GREEN; BOTH TestMigrationCheckPrunesStateTree AND TestMigrationCheckFixturesParseConsistently pass. Confirm go test ./... fully green on the merge result.
  Merge result (origin/next 7bcfffa5 + branch, clean no-ff merge touching only the 2 migration-check files): `go test ./...` → 15 ok, 0 FAIL. Non-vacuity proven: on origin/next WITHOUT the branch, TestMigrationCheckFixturesParseConsistently is RED (docs/roadmap/.../debrief.md session-date reader="2026-06-08" vs direct="2026-06-08T00:00:00Z") — the branch is the only thing un-REDing the suite for the tag.
- DONE: Detached adversarial audit (f1 = a status GUARD path; risk = OVER-pruning) on a separate throwaway checkout — confirm the widened prune did NOT blind the guard; construct an entity-frontmatter inconsistency the migration-check should catch and confirm the widened prune still catches it.
  Separate /tmp/sd-audit checkout (never the impl worktree), 4 adversarial edits, restored clean after each. Edit1: colon-space poison entity in walked docs/dev/ → guard RED (caught). Edit2: poison under docs/dev/roadmap/ (parent≠docs) → guard RED (caught) → prune is path-precise. Edit3: poison directly under docs/roadmap/ → guard PASS (intended skip) → prune fires, non-vacuous. Surgical measurement: neutering only the roadmap branch shifts checked 57→58, the single newly-visited file is debrief.md (NON-ENTITY, no id:/status:); a probe over every fenced .md under docs/roadmap found exactly one file, non-entity. The prune skips ZERO real entities. No MATERIAL finding.

### Summary
VERDICT: PASSED. The refactor is confirmed (one shared migrationCheckWalkDir drives both the live consistency check and the hermetic prune test — proven by the broken-helper edit reding both, not a copy), `go test ./...` is fully green on the merge result (15 ok / 0 FAIL, non-vacuous against an origin/next RED baseline), and the detached adversarial audit found NO over-pruning (the widened docs/roadmap prune skips exactly one non-entity debrief and still catches real entity inconsistencies in every non-pruned location, path-precisely). One Polish (non-blocking): the migration_check.go comment AND the cycle-2 implementation summary cite `boundary_guard_test.go` (no such file exists) and `discoverIgnoreDirs` (prunes .spacedock-state but NOT roadmap) as precedents for the roadmap-under-docs prune; neither establishes that precedent — this is the only such prune in the codebase. Comment-accuracy only; no behavioral or test-strength hole, does not block the cut.
