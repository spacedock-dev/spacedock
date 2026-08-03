---
id: 69cg49x2fzprarbg32pkbvgn
title: Migrate 7 PR #600 test fixtures to testgit.InitRepo (git-scaffold guard reds on main)
status: validation
source: "main is currently broken: go test ./internal/contractlint/... -run TestNoHandRolledGitInitOutsideTestgit fails with 7 non-bare git init sites outside testgit.InitRepo, all in collapse-gate-approval-ceremony's (#600) own test files. shared-git-scaffold-helper's (#601) ideation explicitly anticipated this exact scenario ('Whichever branch lands second should migrate it to testgit.InitRepo... otherwise the guard will flag it on merge') -- #600 merged before the guard existed on main, #601 landed the guard second, and #600's fixtures were never migrated. Captain directed: fix, main is broken."
started: 2026-08-03T01:43:27Z
completed:
verdict:
score: 1.0
worktree: .worktrees/spacedock-ensign-migrate-pr600-fixtures-to-testgit
issue:
gates:
    version: 1
    records:
        - id: gate:69cg49x2fzprarbg32pkbvgn:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:69cg49x2fzprarbg32pkbvgn-backlog-1
              briefing:
                id: briefing:69cg49x2fzprarbg32pkbvgn:backlog:attempt-1:revision-1
                digest: sha256:3dcb2e496bec15fb1f6fff236ea46e95ff6f3c3c9d458b5bce7718eac1a91848
                request-digest: sha256:974a38278549ec3905b599efe3b1954fda0a14f1489e814fea511a36d1a78854
                room-ref: ./migrate-pr600-fixtures-to-testgit/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:69cg49x2fzprarbg32pkbvgn:backlog:1
                briefing: briefing:69cg49x2fzprarbg32pkbvgn:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T01:42:24.047511Z"
                decision: approve
                reason: 'Captain approved in chat: main is red on this, fast-track to implementation.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:69cg49x2fzprarbg32pkbvgn:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:69cg49x2fzprarbg32pkbvgn-ideation-1
              briefing:
                id: briefing:69cg49x2fzprarbg32pkbvgn:ideation:attempt-1:revision-1
                digest: sha256:e6f9816f47a12ef491c570d9ded8aa1368334d00259606695b1cfdc022c827cb
                request-digest: sha256:9cb31ff23a7a31fa4ee5559b81890b085e9fcb15206e21c7775bd354252c9668
                room-ref: ./migrate-pr600-fixtures-to-testgit/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:69cg49x2fzprarbg32pkbvgn:ideation:1
                briefing: briefing:69cg49x2fzprarbg32pkbvgn:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T01:43:17.007022Z"
                decision: approve
                reason: Captain directed fast-track straight to implementation; design fully pre-scoped in the filing.
              application:
                target-stage: implementation
                state: consumed
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

## Stage Report: implementation

- DONE: Migrate all 7 hand-rolled git init scaffold sites to testgit.InitRepo
  gate_ceremony_count_test.go:71,84; gate_consume_sync_test.go:133,144,385,396; dispatch/build_stamp_test.go:30 — commit 54d52d1c2 on spacedock-ensign/migrate-pr600-fixtures-to-testgit
- DONE: Run go test ./internal/contractlint/... -run TestNoHandRolledGitInitOutsideTestgit -v and confirm 0 offenders (currently 7)
  before: 7 offenders listed by name; after: PASS, 0 offenders
- DONE: Run go test ./internal/cli/... ./internal/dispatch/... and the full go test ./... plus go test ./... -race; confirm all green
  internal/cli and internal/dispatch both ok; full `go test ./...` ok across all 17 packages; full `go test ./... -race` ok across all 17 packages, no data races
- DONE: gofmt -w ./cmd ./internal; confirm clean
  `gofmt -l ./cmd ./internal` reported no files both before commit and after

### Summary

Read the already-migrated idiom from the 38-site migration (commit 97007928d) in the same packages — internal/cli/gate_test.go and internal/dispatch/{parity_harness_test.go,reconcile_test.go} — then applied the identical substitution: each 3-line `git init` + `config user.name` + `config user.email` (or file-local `runGitFatal` equivalent) collapsed to one `testgit.InitRepo(t, dir, "-q"[, extra init args])` call, plus the `internal/testgit` import. No assertions, fixture content, or identity values used elsewhere in these tests were touched. Guard test dropped from 7 offenders to 0; full suite and `-race` green with no other changes.
