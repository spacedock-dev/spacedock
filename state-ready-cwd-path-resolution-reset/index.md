---
title: state ready resolves split state from linked-worktree cwd (clean reset)
status: validation
source: Clean reset from rejected e6j implementation and GitHub issue #484, captain direction 2026-07-14
started: 2026-07-14T14:12:39Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-state-ready-cwd-path-resolution-reset
issue: spacedock-dev/spacedock#484
milestone: 0.25.0
id: qx18j81e0b1pe84ck9jxrbkj
mod-block: merge:pr-merge
---

Make `spacedock state ready` recover the declared split-root state checkout from inside an agent worktree and in a local-only repository, without becoming a worktree/session manager.

## Problem

Issue #484 showed two user-visible failures while recovering a deleted split-root checkout. From `.worktrees/<worker>/<workflow>`, discovery anchors the state path below that linked worktree instead of the repository's main worktree. In a repository with no `origin`, recovery hard-fails on `git fetch origin` even when the state branch exists locally and ordinary `git worktree add` can restore it.

A stale registration may still block `git worktree add` after the directory was manually deleted. Automatic stale-registration repair is not required here: fail closed with precise, main-root-anchored remediation. The prior implementation branch was rejected for growing locks, outcome files, generations, publication protocols, quarantine, and OS-specific atomic machinery around this narrow command path. This task carries none of its stage reports or design.

## Required outcome

- `state ready`, when invoked from a linked worktree, targets the workflow's declared checkout under the repository's main worktree and never creates a nested state checkout.
- With no `origin`, it restores from the existing local state branch.
- With a stale registration, it makes no destructive repair and reports the exact registered/main-root path plus the smallest manual remediation.
- Do not expand behavior for `state commit`, `state new`, `state sweep`, dispatch, merge, or general workflow discovery unless a shared helper change is strictly necessary for the two outcomes above.

## Mechanism/value trace

- Served value: the issue-484 recovery command converges from a linked-worktree cwd and a local-only repo.
- Simplest route: resolve the existing Git main-worktree identity, then use the existing local branch/worktree-add path; exercise it with real Git fixtures.
- Rejected heavier route: no lock protocol, durable outcome journal, generation token, private publication scheme, admin-record quarantine, OS-specific no-replace primitive, daemon, lease, or lifecycle supervisor. Encountering a need for one is a design-reset trigger.

## Acceptance criteria

- **AC-1 (VALUE):** In the issue-484 real-Git fixture, `state ready` invoked under an agent worktree restores/uses the checkout only at the main-worktree declared path, exposes the state files there, and creates nothing under `.worktrees/`.
- **AC-2:** With no `origin` and an existing local state branch, `state ready` exits 0 and restores the checkout from that branch; an origin-backed repository retains its existing fetch-first behavior.
- **AC-3:** A stale registration fails closed without removing or rewriting any registration and emits precise remediation naming the canonical state path. A present checkout remains untouched.
- **AC-4:** Focused real-Git tests, `go test ./...`, and `go test ./... -race` pass. The implementation introduces no new coordination/lifecycle subsystem and remains a narrow change around existing Git commands.

## Implementation result

Code commit `0384a012` changes only `internal/cli/state_sync.go` and its real-Git
`state_ready_test.go` fixtures.

The exact mechanism/value trace is:

- **Value observed:** from a workflow discovered inside a linked agent worktree,
  `state ready` restores `first-task.md` at the canonical main-worktree state
  path and leaves the agent-worktree state path absent.
- **Smallest mechanism:** Git supplies the workflow-relative prefix and canonical
  main worktree; the existing branch plus `git worktree add` restores the checkout.
  Origin-backed recovery fetches first; no-origin recovery verifies the local ref.
- **Fail-closed boundary:** a missing path still present in `git worktree list
  --porcelain` exits non-zero, leaves the registry byte-identical, and prints the
  exact `git -C <main-root> worktree remove <state-path>` manual remediation.
- **RED to GREEN:** the three new fixtures initially failed 3/3 (nested-path
  recovery, mandatory `fetch origin`, and generic stale-registration failure),
  then passed 3/3; the complete `TestStateReady` family, full suite, and race suite
  passed after the change.
- **Rejected machinery:** no lock, journal, generation, publication protocol,
  automatic repair, daemon, lease, or lifecycle controller was introduced.

## Stage Report: implementation

- DONE: Make state ready use the canonical main-worktree state path from a linked-worktree cwd and restore from a local branch when origin is absent.
  Commit `0384a012`; real-Git linked-worktree and no-origin fixtures both pass and assert the resulting on-disk entity path.
- DONE: Fail closed with precise remediation for stale registration; do not add automatic repair, locks, outcome journals, generations, publication protocols, or lifecycle machinery.
  The stale fixture exits non-zero with an unchanged worktree registry and exact path-scoped `worktree remove` guidance; the diff adds no coordination subsystem.
- DONE: Prove the narrow command behavior with real-Git fixtures, then run focused, full, and race suites and record the exact mechanism/value trace.
  `TestStateReady*` passed (9.977s), `go test ./internal/cli` passed (67.292s), `go test ./...` passed, `go test ./... -race` passed, and the RED 3/3 baseline is recorded above.

### Summary

Implemented the clean reset as a narrow `state ready` recovery change backed by
real Git behavior, with canonical main-root placement, local-branch fallback,
and non-mutating stale-registration refusal. All required focused, full, and race
gates passed; code is committed at `0384a012` and ready for independent validation.

### Roborev request

- **Identity:** job `1426`, review `1412`, UUID
  `6448bbfa-5461-4264-8a14-c5ac84355a15`, patch
  `3db7a3fd7140070070ed7e99ceb866776a9e3367`; exact `git_ref`
  `0384a012faa25bae41ac56f7cc8b1347122942b2`, agent `codex`, reasoning
  `thorough`.
- **Status:** `done`, verdict `F`, zero retries; enqueued
  `2026-07-14T15:27:33Z` and finished `2026-07-14T23:30:46+08:00`.
- **HIGH — `internal/cli/state_sync.go:159`:** `branch` and `stateRel` still
  come from the linked-worktree README after the checkout path moves to the
  main worktree. If an agent branch changes `state-branch`, `state ready` can
  pull/rebase that branch into the existing main state checkout and corrupt its
  history. Re-resolve and validate the main-worktree README, fail closed on any
  disagreement, and verify the present checkout is on the expected branch
  before pulling.
- **MEDIUM — `internal/cli/state_sync.go:228`:** both worktree-list parsers use
  line-delimited porcelain, so Git C-quoted paths (including backslashes, tabs,
  newlines, or quoted non-ASCII byte sequences) can produce an invalid main root
  or hide a stale registration. Use `git worktree list --porcelain -z`, parse
  NUL-delimited records, and cover a path that triggers quoting.
- **MEDIUM — `internal/cli/state_sync.go:249`:** `registeredWorktree` converts a
  registry-query failure into “not registered,” after which recovery can fetch
  and attempt mutation while registry state is unknown. Return registration plus
  error and abort recovery when the query fails.

No code changed and no verification was rerun during this request; these concrete
findings are recorded for the feedback/rejection implementation cycle.

## Stage Report: validation

- DONE: Verify AC-1 and AC-2 with independent real-Git fixtures: linked-worktree cwd restores the canonical main-worktree state path, and a missing origin falls back to the local state branch.
  The exact-head binary recovered `state-ready-cwd-path-resolution-reset/index.md` only at the canonical main path in both origin-backed and no-origin runs (exit 0); forced origin failure exited 1 before checkout creation, proving fetch-first remained intact.
- DONE: Verify AC-3: stale linked-worktree registration fails closed with precise path-scoped remediation, leaves the registry unchanged, and adds no repair or coordination subsystem.
  The `state ready` stale run exited 1 and left the registry byte-identical; separately, the printed manual `git worktree remove` command exited 0 and removed only that registration, while a present-checkout run preserved registry, HEAD, and gitfile byte-for-byte.
- DONE: Reproduce the focused behavior checks plus applicable package, full, and race suites on exact implementation head 0384a012; map independent evidence to every AC and issue PASSED or a precisely classified REJECTED recommendation.
  At `0384a012faa25bae41ac56f7cc8b1347122942b2`: focused AC tests passed 4/4, all `TestStateReady*` tests passed 10/10, `go test ./internal/cli`, `go test ./...`, and `go test ./... -race` passed; `gofmt -w ./cmd ./internal` completed and the final tree is clean.

### Summary

PASSED. AC-1 through AC-4 are independently reproduced with real Git state,
process exits, exact paths, registry snapshots, and the required suites. The
implementation stays narrow: only existing `state_sync.go` and its real-Git test
file changed, with no repair, locking, publication, or lifecycle subsystem.

### Reviewer Findings

- CLEAN — Detached adversarial audit of the high-stakes stale-registration guard at exact head `0384a012faa25bae41ac56f7cc8b1347122942b2` used a detached throwaway checkout, never the implementation worktree.
  The one-line adversarial edit bypassed `registeredWorktree`; `go test ./internal/cli -run '^TestStateReadyStaleRegistrationFailsClosed$' -count=1 -v` exited 1 at `state_ready_test.go:351` because the raw `git worktree add` failure replaced the required non-mutating guard/remediation output.
  Restoring that one line returned the same test to PASS; the throwaway checkout was removed and the implementation worktree remained clean, so the audit found no test-strength hole or material feedback cycle.
