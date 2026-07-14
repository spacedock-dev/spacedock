---
title: state ready resolves split state from linked-worktree cwd (clean reset)
status: implementation
source: Clean reset from rejected e6j implementation and GitHub issue #484, captain direction 2026-07-14
started: 2026-07-14T14:12:39Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-state-ready-cwd-path-resolution-reset
issue: spacedock-dev/spacedock#484
milestone: 0.25.0
id: qx18j81e0b1pe84ck9jxrbkj
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
