---
id: qyc6g8bmvcdsj7bdz7sjwgbn
title: internal/status test-hygiene — go test from root fails on debrief frontmatter, plus gofmt-dirty files
status: backlog
source: "captain (2026-06-04) — surfaced by the xa ideation ensign and verified this session. `go test ./...` from the project root fails because TestMigrationCheckFixturesParseConsistently scans the .spacedock-state debrief fixtures; two internal/status files are also gofmt-dirty on next."
score: "0.26"
started:
completed:
verdict:
worktree:
issue:
---

`go test ./internal/status/` (and thus `go test ./...`) FAILS when run from the project root, and two `internal/status` files are gofmt-dirty on `next`. Neither is caught by CI (the offline gate runs a fresh checkout with no `.spacedock-state`; worktree ensigns also lack it), so the failures only bite a developer or a non-worktree ensign running tests from the repo root — which is the common dev loop.

## The failures (verified 2026-06-04)

- **`TestMigrationCheckFixturesParseConsistently` fails from root.** It scans `docs/dev/.spacedock-state/_debriefs/*.md` and the debrief `session-date` field parses inconsistently between two paths — `migration_check_test.go:113`: `key "session-date" reader="2026-06-03" direct="2026-06-03T00:00:00Z"`. A bare YAML date (`session-date: 2026-06-03`) is read as a string by one parser and a full timestamp by the other. The test is checking parser consistency across fixtures + live state, and the debrief frontmatter (which is NOT entity frontmatter) trips it.
- **Two gofmt-dirty files on `next`:** `internal/status/external_proof.go` and `internal/status/no_yaml_silent_drop_test.go` (a curly-quote-in-comment quirk per prior validation notes). `gofmt -l internal/status/` lists both on a clean tree.

## Direction (for ideation to flesh out)

- The migration-check test should likely **exclude `_debriefs/` (and other non-entity dirs)** from its fixture scan — debriefs are session records with their own frontmatter shape (`session-date`, `sequence`, `first-commit`), not entity fixtures the migration check governs. Alternatively, make the date-field parse consistent for bare YAML dates. Pick the smaller, correct fix.
- `gofmt -w` the two dirty files (and consider why CI's offline gate doesn't gofmt-check — a `gofmt -l` guard in the offline lane would prevent recurrence).
- AC: `go test ./...` passes from the project root with `.spacedock-state` present, and `gofmt -l ./...` is empty — both checkable by command.

## Notes

Cleanup/hygiene task (the change is the fix plus the now-passing command). Surfaced by the `feedback-guarantee-binary-gate` (xa) ideation ensign's out-of-scope flag; verified by the FO. Provenance: session 2026-06-04 #8.
