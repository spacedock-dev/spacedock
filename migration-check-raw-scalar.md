---
title: Migration-check compares RAW scalar frontmatter values (un-red main; fix the auto-typed-date false failure)
status: backlog
score: 0.6
sprint: 0221-layered-fo
group: binary-ux
id: vkh9e4398r126dq17r0njz5p
---

Pre-existing `main`-red: `TestMigrationCheckFixturesParseConsistently` (`internal/status`) fails on `docs/dev/_debriefs/2026-06-19-01.md`'s unquoted `session-date: 2026-06-19`. The test's `direct` decode — `yaml.Unmarshal` into `map[string]any` (`migration_check_test.go:100`) — lets yaml.v3 auto-type the bare date to `time.Time`, which `scalarString` (`:246`) renders `"2026-06-19T00:00:00Z"`, disagreeing with `ParseFrontmatter`'s string `"2026-06-19"`. The product reader is CORRECT per the contract's "every frontmatter value is a string"; the TEST is over-strict (it compares the string reader against yaml's auto-TYPED view). This reds `offline` on EVERY PR (verbs #404 had to `--admin` past it).

## Fix (test-only, root-cause)
In `internal/status/migration_check_test.go`, decode `direct` into `map[string]yaml.Node` (not `map[string]any`) and compare `node.Value` (the RAW scalar text) for scalar nodes / `""` for nested (Mapping/Sequence kinds) — so a bare date stays `"2026-06-19"`, matching the reader + the contract's string rule. Do NOT edit the debrief or quote any dates; fix the root cause (the test's typed comparison). The test must STILL catch real parser divergences (a quoting/escaping/multiline case the line-reader would mishandle).

## Acceptance criteria
- **AC-1** — `go test ./internal/status -run TestMigrationCheckFixturesParseConsistently` PASSES, including against the unquoted `session-date` debrief. Verified by the green run on the current tree (the failing fixture is present).
- **AC-2** — the relaxation is NOT vacuous: the test still goes RED on a genuine reader-vs-yaml divergence (e.g. a deliberately mis-parsed quoted/escaped fixture, or a reader mutation). Verified by a negative control that makes it RED, then restored.
- **AC-3** — test-only change: `go test ./internal/status` and `go build ./...` GREEN; NO product-code change and NO debrief/frontmatter edits (`git diff --stat` touches only `migration_check_test.go`).
