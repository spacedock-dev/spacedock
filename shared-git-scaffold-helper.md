---
id: zfmbm75wvmfj38h73wrtpqqy
title: Add a shared, persisting git-scaffold test helper to close the identity-config gap
status: ideation
source: "Found while dogfooding collapse-gate-approval-ceremony's CI run, 2026-08-02: TestStampCommitsInlineBeforeWorktreeCreation's fixture omitted git config user.name/user.email, passing locally (ambient git identity) but failing on the CI runner (none). A follow-up search found no shared, persisting git-scaffold-init helper anywhere in the repo -- roughly 15-20 fixture functions across internal/cli and internal/dispatch each hand-roll their own init+config sequence independently, including one file-local gitInit helper (internal/dispatch/parity_harness_test.go) that only sets identity via scoped -c flags on its own seed commit, not a persisted git config, so it doesn't even cover a later plain git commit in the same repo. Captain directed: file and fast-track."
started: 2026-08-02T14:31:58Z
completed:
verdict:
score: 0.5
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:zfmbm75wvmfj38h73wrtpqqy:backlog
    records:
        - id: gate:zfmbm75wvmfj38h73wrtpqqy:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zfmbm75wvmfj38h73wrtpqqy-backlog-1
              briefing:
                id: briefing:zfmbm75wvmfj38h73wrtpqqy:backlog:attempt-1:revision-1
                digest: sha256:ed3ebaef25d66a395a77d778b29879d7ff5ad90f40fc8b65d157c3df486059fd
                digest-domain: canonical-bytes
                request-digest: sha256:2ee4ca7169ea0dbbe53f8d8a436d26bdd6bb84a7608e059a04ce7681d557efe2
                room-ref: ./shared-git-scaffold-helper/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zfmbm75wvmfj38h73wrtpqqy:backlog:1
                briefing: briefing:zfmbm75wvmfj38h73wrtpqqy:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T14:31:53.236223Z"
                decision: approve
                reason: 'Captain directed in chat: file the git-scaffold-helper follow-up and dispatch it fast-track.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Add one shared, exported test helper that initializes a scratch git repo with a persisted `user.name`/`user.email` config in a single call, so no future fixture can reproduce the "works locally on ambient identity, fails on a clean CI runner" bug class, and migrate the existing hand-rolled fixtures to it.

## Problem

No shared git-scaffold-init helper exists anywhere in this repo. `internal/cli` and `internal/dispatch` each have ~15-20 test fixture functions that independently call `git init` + `git config user.name` + `git config user.email` (or, in the file-local `gitInit` helper in `internal/dispatch/parity_harness_test.go`, only scope identity via `-c` flags on one seed commit, not a persisted repo config). This is a real, already-observed defect class: `TestStampCommitsInlineBeforeWorktreeCreation`'s fixture omitted the config calls, passed on every local run (this and other dev environments already have an ambient git identity configured), and failed on the GitHub Actions runner (no ambient identity) -- caught only by an actual CI run, not by any local or code-review signal.

## Proposed approach

Add one small, exported helper (name and package location left to ideation -- candidates: a new `internal/cli/testlib` or `internal/dispatch/testlib` package, or a single shared location both packages can import) that takes `(t *testing.T, dir string)` and performs `git init` + a **persisted** `git config user.name`/`user.email` in one call, returning nothing or the same information the existing ad hoc callers already use. Migrate the existing ~15-20 hand-rolled fixtures to call it instead of repeating the sequence, including fixing the file-local `gitInit`/`repoGitInit` helpers in `internal/dispatch` to persist identity rather than scope it per-commit. Do not change any test's actual assertions or fixture *content* beyond this mechanical substitution -- this is a DRY/hermeticity fix, not a behavior change.

## Common scenarios and expected behavior

- A new fixture that uses the shared helper never has to remember the three-step init sequence, and cannot omit persisted identity by construction.
- Existing fixtures, once migrated, produce byte-identical test behavior -- this is a refactor, not a functional change to any test's assertions.
- The helper works identically regardless of the host machine's ambient git identity (or lack of it), which is exactly the gap that let this bug through undetected until a live CI run.

## Out of scope

Changing any test's assertions, coverage, or behavior beyond the init/config mechanism. Any production (non-test) code change. Consolidating other kinds of test-fixture duplication not related to git identity/init.

## Acceptance criteria

**AC-1 (VALUE) - No test fixture in `internal/cli` or `internal/dispatch` can omit a persisted git identity by construction.**
Verified by: every fixture that creates a scratch git repo calls the new shared helper (grep for raw `git(t, ..., "init"` / direct `exec.Command("git", ..., "init"...)` outside the helper itself returns none in these two packages), and a test asserting the helper's own repo has a resolvable `git config user.name`/`user.email` after one call.

**AC-2 - Existing tests are unaffected.**
Verified by: `go test ./internal/cli ./internal/dispatch` and the full `go test ./...` plus `-race` pass unchanged after migration, with no assertion or fixture-content changes beyond the init/config call substitution.

**AC-3 - The specific bug class cannot recur silently.**
Verified by: a new focused test that runs `TestStampCommitsInlineBeforeWorktreeCreation` (or an equivalent using the shared helper) with `HOME` unset and no global git config present (mirroring the CI runner's actual failure condition), confirming it passes without relying on ambient identity.

## Test plan

Offline only, no live workflow needed. `go test ./internal/cli ./internal/dispatch -count=1` after migration, then the full `go test ./...` and `go test ./... -race`. The AC-3 test is the one genuinely new piece of coverage; everything else is a mechanical substitution validated by the existing suite staying green.
