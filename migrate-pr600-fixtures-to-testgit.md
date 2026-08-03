---
id: 69cg49x2fzprarbg32pkbvgn
title: Migrate 7 PR #600 test fixtures to testgit.InitRepo (git-scaffold guard reds on main)
status: backlog
source: "main is currently broken: go test ./internal/contractlint/... -run TestNoHandRolledGitInitOutsideTestgit fails with 7 non-bare git init sites outside testgit.InitRepo, all in collapse-gate-approval-ceremony's (#600) own test files. shared-git-scaffold-helper's (#601) ideation explicitly anticipated this exact scenario ('Whichever branch lands second should migrate it to testgit.InitRepo... otherwise the guard will flag it on merge') -- #600 merged before the guard existed on main, #601 landed the guard second, and #600's fixtures were never migrated. Captain directed: fix, main is broken."
started:
completed:
verdict:
score: 1.0
worktree:
issue:
---

`main`'s `go test ./internal/contractlint/... -run TestNoHandRolledGitInitOutsideTestgit` fails: 7 non-bare `git init` scaffold sites outside `testgit.InitRepo`, all introduced by `collapse-gate-approval-ceremony` (#600):

- `internal/cli/gate_ceremony_count_test.go:71`
- `internal/cli/gate_ceremony_count_test.go:84`
- `internal/cli/gate_consume_sync_test.go:133`
- `internal/cli/gate_consume_sync_test.go:144`
- `internal/cli/gate_consume_sync_test.go:385`
- `internal/cli/gate_consume_sync_test.go:396`
- `internal/dispatch/build_stamp_test.go:30`

## Proposed approach

Mechanical substitution, same pattern as `shared-git-scaffold-helper`'s own migration of the other 38 sites: replace each hand-rolled `git init` + identity-config sequence with a single `testgit.InitRepo(t, dir, ...)` call. No behavior change -- this is a DRY/hermeticity fix identical in kind to the ones already done across the rest of the repo.

## Acceptance criteria

**AC-1 (VALUE)** — `go test ./internal/contractlint/... -run TestNoHandRolledGitInitOutsideTestgit` passes (0 offenders) on `main`.
Verified by: the guard test itself, red before the fix (7 sites), green after.

**AC-2** — No behavior change to any of the 7 affected tests beyond the init/identity call substitution.
Verified by: `go test ./internal/cli/... ./internal/dispatch/...` and the full `go test ./...` plus `-race` green, diff review confirms each hunk is scaffold-site-only.

## Test plan

Offline only. `go test ./internal/contractlint/...`, `go test ./internal/cli/... ./internal/dispatch/...`, then full `go test ./...` and `go test ./... -race`.

## Out of scope

Any other test-fixture cleanup, any production code change.
