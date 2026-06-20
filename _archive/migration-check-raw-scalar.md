---
title: Migration-check compares RAW scalar frontmatter values (un-red main; fix the auto-typed-date false failure)
status: validation
score: 0.6
sprint: 0221-layered-fo
group: binary-ux
id: vkh9e4398r126dq17r0njz5p
worktree:
started: 2026-06-20T01:20:23Z
mod-block:
pr: "#410"
completed: 2026-06-20T02:45:10Z
verdict: PASSED
archived: 2026-06-20T02:45:10Z
---

Pre-existing `main`-red: `TestMigrationCheckFixturesParseConsistently` (`internal/status`) fails on `docs/dev/_debriefs/2026-06-19-01.md`'s unquoted `session-date: 2026-06-19`. The test's `direct` decode — `yaml.Unmarshal` into `map[string]any` (`migration_check_test.go:100`) — lets yaml.v3 auto-type the bare date to `time.Time`, which `scalarString` (`:246`) renders `"2026-06-19T00:00:00Z"`, disagreeing with `ParseFrontmatter`'s string `"2026-06-19"`. The product reader is CORRECT per the contract's "every frontmatter value is a string"; the TEST is over-strict (it compares the string reader against yaml's auto-TYPED view). This reds `offline` on EVERY PR (verbs #404 had to `--admin` past it).

## Fix (test-only, root-cause)
In `internal/status/migration_check_test.go`, decode `direct` into `map[string]yaml.Node` (not `map[string]any`) and compare `node.Value` (the RAW scalar text) for scalar nodes / `""` for nested (Mapping/Sequence kinds) — so a bare date stays `"2026-06-19"`, matching the reader + the contract's string rule. Do NOT edit the debrief or quote any dates; fix the root cause (the test's typed comparison). The test must STILL catch real parser divergences (a quoting/escaping/multiline case the line-reader would mishandle).

## Acceptance criteria
- **AC-1** — `go test ./internal/status -run TestMigrationCheckFixturesParseConsistently` PASSES, including against the unquoted `session-date` debrief. Verified by the green run on the current tree (the failing fixture is present).
- **AC-2** — the relaxation is NOT vacuous: the test still goes RED on a genuine reader-vs-yaml divergence (e.g. a deliberately mis-parsed quoted/escaped fixture, or a reader mutation). Verified by a negative control that makes it RED, then restored.
- **AC-3** — test-only change: `go test ./internal/status` and `go build ./...` GREEN; NO product-code change and NO debrief/frontmatter edits (`git diff --stat` touches only `migration_check_test.go`).

## Stage Report: implementation

- DONE: Apply the parser relax — decode `direct` into `map[string]yaml.Node` and compare `node.Value` (raw scalar) for scalar nodes / `""` for Mapping/Sequence kinds
  `migration_check_test.go:100` now decodes into `map[string]yaml.Node`; new `nodeScalarString` helper returns `n.Value` for `ScalarNode`, `""` otherwise — mirrors `parseFrontmatterContent` exactly. Commit 4220693f.
- DONE: Did NOT edit any debrief or quote any dates
  `git status --short` shows only `migration_check_test.go` modified; debrief `2026-06-19-01.md` left at the peer's quoted state, frontmatter.go restored byte-clean after the negative control.
- DONE: PROVE the root-cause fix handles UNQUOTED bare dates independent of the band-aid
  Added `TestMigrationCheckUnquotedDateStaysRawString` owning its own UNQUOTED `some-date: 2026-06-19` fixture; reader and relaxed direct-decode both yield raw `"2026-06-19"`, GREEN. Anchored with a guard asserting yaml.v3 still auto-types the bare date to non-string (so the relax stays load-bearing).
- DONE: AC-1 — `TestMigrationCheckFixturesParseConsistently` PASSES on current tree
  Verbose run: `--- PASS: TestMigrationCheckFixturesParseConsistently`, 68 frontmatters verified.
- DONE: AC-2 non-vacuity — test goes RED on a genuine reader-vs-yaml divergence, then restored
  Transient reader mutation (`v.Value + "X"`) reddened both the walk (`reader="taskX" direct="task"`) and the new subtest (`reader="2026-06-19X"`); frontmatter.go restored, re-run GREEN.
- DONE: AC-3 — `go test ./internal/status` and `go build ./...` GREEN; diff touches ONLY `migration_check_test.go`
  Full package suite `ok` (18.4s); `go build ./...` clean; `git diff --stat origin/main...HEAD` = `internal/status/migration_check_test.go | 85 ...` (1 file). No product code, no debrief edits.

### Summary

Replaced the over-strict `map[string]any` + `scalarString` comparison (which let yaml.v3 auto-type bare dates to `time.Time`, rendering `2026-06-19T00:00:00Z` against the reader's raw `2026-06-19`) with a `map[string]yaml.Node` decode comparing `node.Value`, exactly mirroring the product reader's "raw scalar text / nested → empty string" contract. A peer had band-aided the live debrief date to quoted on origin/main (335c59fa); I reproduced the original RED on an unquoted date, fixed the root cause in the test only, and proved robustness with a self-owned UNQUOTED fixture that survives any future re-quoting. Negative control confirms the relaxation is not vacuous. Committed test-only on `spacedock-ensign/migration-check-raw-scalar` (4220693f).

## Stage Report: validation

- DONE: Independently reproduce AC-1 — `go test ./internal/status -run TestMigrationCheckFixturesParseConsistently` PASSES on the rebased worktree
  Verbose run GREEN: `--- PASS: TestMigrationCheckFixturesParseConsistently (0.05s)`, 68 frontmatters verified. Worktree rebased on origin/main `f64259fb`; band-aid `335c59fa` confirmed in HEAD ancestry; debrief `2026-06-19-01.md` carries the quoted `session-date: "2026-06-19"`; self-owned UNQUOTED fixture `TestMigrationCheckUnquotedDateStaysRawString` present (commit `8023ff11`).
- DONE: Reproduce AC-2 non-vacuity — relaxed `node.Value` comparison still goes RED on a genuine reader-vs-yaml divergence, then restore
  Transient reader mutation (`fields[key] = v.Value + "X"` in `frontmatter.go:121`) reddened BOTH the live walk (`migration_check_test.go:115: reader="taskX" direct="task"`) and the date-specific subtest (`migration_check_test.go:176: reader: some-date="2026-06-19X"`). `frontmatter.go` restored byte-clean (`git status --short` empty, no diff vs HEAD) — relax proven load-bearing on the exact date case it protects.
- DONE: Confirm AC-3 test-only — `go build ./...` and full `go test ./internal/status` GREEN; diff touches ONLY `migration_check_test.go`
  `go build ./...` exit 0; full package suite `ok` (16.9s); `git diff --name-only origin/main...HEAD` = `internal/status/migration_check_test.go` (1 file, +57/-28). No product code, no debrief/frontmatter edits.

### Summary

PASSED. All three acceptance criteria independently reproduced on the rebased worktree (`8023ff11`, atop origin/main `f64259fb`). The root-cause fix — `map[string]yaml.Node` decode comparing `node.Value` via a `nodeScalarString` helper that mirrors `parseFrontmatterContent` exactly (scalar → raw text, nested → "") — un-reds the migration check without quoting any date. AC-2's negative control confirms the relax is non-vacuous (reds on both the live walk and the self-owned unquoted-date subtest, restored byte-clean). The diff is strictly test-only (1 file). The bidirectional key-coverage assertions (reader→direct and direct→reader) and the self-owned unquoted fixture make the relaxation robust to any future re-quoting of the live debrief.
