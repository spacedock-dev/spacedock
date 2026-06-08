---
id: kbr8yxzknmk8y5451rypz3h4
title: migration-check test should prune the .spacedock-state tree, not name-match _debriefs (+ drop orphaned survey scaffolds fixtures)
status: implementation
source: "post-sprint antipattern audit of 019x-pre-flip-cleanups (2026-06-08, staff-eng persona) — qy's migration-check fix (PR #326) silences _debriefs by name, but the walk still descends into the gitignored .spacedock-state checkout and validates ~107 machine-local entity files. The SAME sprint prunes that tree wholesale 3x elsewhere."
score: "0.22"
started: 2026-06-08T17:26:14Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-migration-check-prune-state-walk
issue:
sprint: 0198-pre-flip-hardening
group: test-hygiene
sprint-readiness: ready
---

Two test-hygiene cleanups surfaced by the 019x post-sprint audit. Bounded (the suite is green); a consistency fix, not a correctness bug.

## Problem

- **migration-check walks into the state checkout.** `internal/status/migration_check_test.go` roots its walk at `docs/` and, on a dev machine where the split-root state checkout is materialized, descends into `docs/dev/.spacedock-state/` — validating ~107 gitignored, machine-local entity files the migration check was never meant to govern (plus the 22 `_debriefs` that qy's `_debriefs`-SkipDir now silences). On CI (fresh clone) the walk never enters the tree, so the skip is moot. The fix name-matches one subdir instead of pruning the tree.
- **Inconsistent with the established pattern.** The codebase prunes the whole `.spacedock-state` tree elsewhere — `internal/status/handlers.go:442` (ignore set), `internal/status/external_proof_test.go:140` (SkipDir), and **qy's own contractlint fix in the same sprint** (`internal/contractlint/boundary_guard_test.go:88`). migration-check should match.
- **Orphaned survey scaffolds fixtures.** Commit `2a87a80f` (xn, PR #324) deleted `skills/integration/survey_scaffold_test.go` but left `skills/integration/testdata/survey/scaffolds/` (6 tracked files, 0 consumers). xn's own validation report recommended pruning them. Dead testdata.

## Proposed approach

- migration-check: prune `docs/dev/.spacedock-state` in the walk (replace the `_debriefs` name-match with a tree-prune matching the three siblings).
- Delete the orphaned `skills/integration/testdata/survey/scaffolds/` fixtures.
- (Optional, Polish, audit-flagged: simplify `internal/dispatch/build.go`'s `teamNamePattern` back to reusing `namePattern` — it diverges only on a single-char case the harness never produces. Path-safe + tested either way; lowest priority.)

## Acceptance criteria (sketch — ideation/impl firms)

- `go test ./...` from the repo root stays green with `.spacedock-state` present, AND the migration-check no longer descends into the state checkout — verified by an assertion that the walk skips the tree (not just `_debriefs`).
- `grep -rln scaffolds --include=*.go` finds zero consumers AND the fixture tree is gone.

## Notes

Provenance: the 019x post-sprint two-persona audit (its one Material finding + a Polish cleanup). Same test-hygiene class as qy. Candidate for a 0.19.8 / pre-flip sprint.
