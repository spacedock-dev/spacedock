---
id: kbr8yxzknmk8y5451rypz3h4
title: migration-check test should prune the .spacedock-state tree, not name-match _debriefs (+ drop orphaned survey scaffolds fixtures)
status: done
source: "post-sprint antipattern audit of 019x-pre-flip-cleanups (2026-06-08, staff-eng persona) — qy's migration-check fix (PR #326) silences _debriefs by name, but the walk still descends into the gitignored .spacedock-state checkout and validates ~107 machine-local entity files. The SAME sprint prunes that tree wholesale 3x elsewhere."
score: "0.22"
started: 2026-06-08T17:26:14Z
completed: 2026-06-08T18:01:12Z
verdict: PASSED
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: test-hygiene
sprint-readiness: ready
mod-block:
pr: "#327"
archived: 2026-06-08T18:01:12Z
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

## Stage Report: implementation

- DONE: migration-check walk PRUNES the `.spacedock-state` tree (tree-prune matching the three sibling prunes, NOT the `_debriefs` name-match) — proven by a POSITIVE assertion that the walk skips the state tree, and the `checked == 0` Fatal guard stays non-vacuous after the prune.
  `isMigrationCheckPrunedDir` returns `.spacedock-state` (matching `handlers.go` discoverIgnoreDirs + `boundary_guard_test.go`); `_debriefs` lives only inside the state tree, so the tree-prune subsumes the old name-match. Commit `96f22f5f`.
- DONE: POSITIVE assertion that the walk skips the state tree.
  Two proofs — (1) live-corpus arm in `TestMigrationCheckFixturesParseConsistently` asserts no `visited` path crosses a `.spacedock-state` segment; (2) hermetic `TestMigrationCheckPrunesStateTree` drives the same predicate over a planted temp tree. Mutation control (predicate → no-op) with the real 159-file state tree copied in flagged every descended entity; both went red, proving non-vacuous.
- DONE: `checked == 0` Fatal guard stays non-vacuous after the prune.
  With `.spacedock-state` (159 `.md`) present, the walk verified 57 frontmatters (NOT 57+159) — the tree was skipped and real entities still validated. Hermetic test also asserts the real entity IS visited.
- DONE: Orphaned `skills/integration/testdata/survey/scaffolds/` fixtures deleted — `grep -rln scaffolds --include=*.go` finds zero consumers AND the fixture tree is gone.
  `git rm` removed 6 tracked files; `grep -rln 'survey/scaffolds'` over the repo (and `--include='*.go'`) returns zero matches; sibling `testdata/survey/fixture-sessions.sql` (the live consumer in `survey_queries_test.go`) preserved.
- DONE: `go test ./...` from the repo root stays green with `.spacedock-state` present (DoD#2).
  Ran full `go test ./...` from the worktree root after copying the live 159-file state tree into the (gitignored) worktree path — all packages `ok`, including `contractlint` (its own `.spacedock-state` prune) and `status`. Temp copy then removed; worktree tree clean.
- SKIPPED: Optional `teamNamePattern` simplification in `internal/dispatch/build.go`.
  Marked Polish / lowest-priority and "include only if cheap and clearly safe." Reusing `namePattern` would TIGHTEN a path-safety validator (reject single-char team names `teamNamePattern` currently admits) — a semantic narrowing of a defense-in-depth guard, not a safe no-op, and out of scope for this test-hygiene entity. Left untouched.

### Summary

Replaced the migration-check `_debriefs` name-match with a wholesale `.spacedock-state` tree-prune (`isMigrationCheckPrunedDir`), matching the three established sibling prunes. Proved the prune two ways: a live-corpus positive assertion (active on a dev machine) and a hermetic CI-stable test, both shown non-vacuous via a mutation control run against the real 159-file state tree. Deleted the 6 orphaned survey-scaffold fixtures (zero consumers). Full `go test ./...` is green both with and without `.spacedock-state` present. Skipped the optional `teamNamePattern` polish as a semantic change out of scope. Code committed to the worktree branch as `96f22f5f`.

## Stage Report: validation

- DONE: Reproduce DoD#2 independently: `go test ./...` from the repo root is green WITH `.spacedock-state` materialized (the ~159-file state tree present), and report the pass ratio. Confirm the suite is also green without it (CI shape).
  WITH state tree (copied live 162 `.md` into worktree `docs/dev/.spacedock-state`, gitignored via `**/.spacedock-state/`): `go test -count=1 ./...` all 16 packages `ok` (cmd/spacedock has no tests). Without it (CI shape): 1151 tests pass, 16 packages `ok`. Both green; reproduction tree removed after.
- DONE: Verify the prune is a NON-VACUOUS tree-prune: independently confirm the mutation-control claim — neutralizing `isMigrationCheckPrunedDir` (so the walk descends into `.spacedock-state`) turns the positive assertion RED, and that `checked > 0` after the prune (the Fatal guard still fires on real entities, not silently disabled).
  Neutralized predicate to `return false` with state tree present → BOTH arms RED: live-corpus arm emitted 162 "descended into the pruned .spacedock-state tree" hits (one per machine-local entity); hermetic `TestMigrationCheckPrunesStateTree` failed at line 235 (visited the planted poison `.spacedock-state/index.md`). Reverted (diff vs HEAD empty) → green, `migration-check verified 57 frontmatters` (57 real entities, NOT 57+162 — guard non-vacuous on real entities).
- DONE: Verify the orphaned fixtures are gone: `grep -rln scaffolds --include=*.go` finds zero consumers, the `skills/integration/testdata/survey/scaffolds/` tree no longer exists, and the live `testdata/survey/fixture-sessions.sql` consumer (used by `survey_queries_test.go`) is preserved.
  `grep -rln scaffolds --include='*.go'` → exit 1, zero matches; `git ls-files | grep -c scaffolds` → 0 tracked. Scaffolds tree GONE on disk; the 6 files are deleted in commit `96f22f5f`. `fixture-sessions.sql` preserved (11.4K) and actively loaded by `survey_queries_test.go:33`.

### Summary

Validated all three completion-checklist items with independent evidence on worktree HEAD `96f22f5f`. Full `go test ./...` is green both WITH the 162-file state tree materialized and without it (CI shape, 1151 tests). The prune is proven non-vacuous: a mutation control (neutralizing `isMigrationCheckPrunedDir`) turns both the live-corpus assertion and the hermetic test RED for the right reasons, while the un-mutated run still verifies 57 real frontmatters (`checked > 0`). Orphaned survey-scaffold fixtures are gone with zero Go consumers, and the live `fixture-sessions.sql` consumer is preserved. The change is test-hygiene only (no production status mutation/guard, launcher, contract, or CI/release code touched), so the detached adversarial audit is not triggered — the mutation control already supplies the equivalent refutation. Recommendation: PASSED.
