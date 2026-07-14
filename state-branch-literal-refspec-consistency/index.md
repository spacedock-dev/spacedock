---
title: Make literal state-branch refspecs consistent
status: ideation
source: "Roborev job 1434 on e6j exact head 58b304ef; captain filing request 2026-07-15."
started: 2026-07-14T16:49:02Z
completed:
verdict:
score:
worktree:
issue:
milestone: 0.25.0
id: 0by8fscax8t88f5wnktggees
---

## Problem

`internal/status.StateBranch` documents that an explicit `state-branch:` override wins verbatim. Git branch creation therefore treats a legal configured value such as `refs/heads/foo` as the literal branch name `refs/heads/foo`, whose full local ref is `refs/heads/refs/heads/foo`. Existing origin-backed fetch, pull, and push calls pass the configured string directly, where Git instead interprets it as the qualified ordinary branch `foo`. Local creation/verification and remote synchronization can therefore address different histories.

Roborev exposed this while reviewing `state-ready-cwd-path-resolution-reset`, but the defect is shared by state ready and state commit. It must not be folded into that task's linked-worktree/no-origin recovery scope.

## Required outcome

- Preserve the documented verbatim meaning of `state-branch:` overrides across local branch creation, checkout verification, fetch, pull/rebase, and push.
- Make `state ready` and `state commit` address the same full local and remote branch for qualified-looking but legal literal names.
- Preserve behavior for ordinary configured names such as `spacedock-state/dev`.
- Fail closed before pulling, rebasing, or pushing a different history when branch/refspec resolution is invalid or inconsistent.
- Do not change linked-worktree path anchoring, stale-registration recovery, workflow discovery, or merge behavior.

## Mechanism/value trace

- Served value: every state command synchronizes the one branch the workflow actually configured, without silently selecting a different history.
- Simplest route: centralize the existing configured-literal-to-full-ref/refspec mapping and feed explicit refs to the existing Git commands.
- Rejected heavier route: no new state protocol, lock, journal, generation, publication layer, daemon, lease, or lifecycle supervisor.

## Acceptance criteria

- **AC-1 (VALUE):** In a real origin-backed fixture with `state-branch: refs/heads/foo`, local creation, `state ready`, pull/rebase, `state commit`, and push all address the literal branch whose full ref is `refs/heads/refs/heads/foo`; the ordinary remote branch `foo` remains untouched.
- **AC-2:** Local-only and origin-backed fixtures prove the same configured override resolves to the same history, and a mismatch fails before any wrong-history mutation.
- **AC-3:** Existing behavior for `spacedock-state/dev` remains unchanged across ready and commit paths.
- **AC-4:** Focused real-Git regressions, `go test ./...`, and `go test ./... -race` pass; the change reuses existing Git command flow and introduces no coordination or lifecycle subsystem.

## Relationship

This is a prerequisite for rebase and revalidation of `state-ready-cwd-path-resolution-reset`; that task remains frozen and unpushed until this shared consistency defect lands.
