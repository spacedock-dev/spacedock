---
title: Make the split-root state checkout operable from foreign path surfaces (Cowork sandbox)
status: backlog
source: live Cowork dogfood 2026-07-13 (survey-claude-cowork-runtime-detection run-1 session)
started:
completed:
verdict:
score:
worktree:
issue:
id: qwp7vn4dnt0sx96wzz1wfy8b
---

Make `spacedock state commit|ready|sweep` work when the repo is reached through a path surface different from where the state worktree was created — observed live: a Cowork sandbox mounts the repo at `/sessions/<name>/mnt/<folder>` while the linked worktree's `.git` pointer carries the host-absolute gitdir, so the verb and bare git both fail with `not a git repository`.

## Problem

`git worktree add` writes absolute paths into the checkout's `.git` pointer file and the `gitdir` back-pointer. The split-root state checkout lives INSIDE the repo tree (`docs/dev/.spacedock-state` → `.git/worktrees/<name>`), so any second path surface for the same repo (Cowork sandbox mount, container bind-mount, NFS re-mount) breaks every state verb even though all the actual git data is present and reachable relatively. Observed live in a Cowork session: `spacedock state commit` failed; a manual `GIT_DIR=.git/worktrees/-spacedock-state GIT_WORK_TREE=docs/dev/.spacedock-state` override committed fine.

## Proposed approach

Two lanes, smallest first:

1. **Verb-side gitdir re-derivation (compat, no repo mutation).** The state verbs already discover the project root. When the checkout's `.git` pointer names a gitdir path that does not exist, recompute `<project-root>/.git/worktrees/<basename>` and invoke git with explicit `GIT_DIR`/`GIT_WORK_TREE` (or `--git-dir/--work-tree`). Works on any git version (sandbox observed: 2.34) and heals both directions (host→sandbox, sandbox→host).
2. **Relative pointers at creation (structural).** Git ≥2.48 supports `worktree.useRelativePaths` / `git worktree add --relative-paths`; because the state checkout is inside the repo tree, relative pointers make the checkout portable across surfaces with no runtime logic. Adopt in `spacedock state init` when the executing git is new enough; keep lane 1 as the fallback for existing checkouts and old gits.

## Out of scope

- Credential/push handling from sandboxes (separate concern; the verb's push-failure surfacing already covers it).
- Non-state (code) worktrees under `.worktrees/` — same mechanism could apply later, but state verbs are the broken path today.

## Acceptance criteria

**AC-1 - `spacedock state commit` succeeds when the repo is addressed via a path different from the checkout's recorded absolute gitdir.**
Verified by: a fixture that clones a split-root repo, moves/bind-mount-simulates it to a second path (rename the dir), and shows `state commit` path-scoped-commits from the new path; failing today, green after lane 1.

**AC-2 - A state checkout created by a git ≥2.48 `state init` carries relative pointers and needs no runtime override.**
Verified by: version-gated fixture asserting the `.git` pointer content is relative when supported, with graceful skip on older git.

## Test plan

Fixture-level: moved-repo scenario for lane 1 (cheap, no sandbox needed); version-gated init assertion for lane 2. Live Cowork smoke optional once lane 1 ships — the survey task's run-2 session can double as it.
