---
id: at3jmf35cgt1ghgm53m91en8
title: status --boot ORPHANS misreports DIR_EXISTS/BRANCH_EXISTS for split-root code worktrees
status: backlog
source: spacedock-dev/spacedock#251 + FO boot (2026-06-01) — reproduced live this session: boot flagged the 7h code worktree as dir_exists/branch_exists=no
started:
completed:
verdict:
score: "0.28"
worktree:
issue: spacedock-dev/spacedock#251
---

Under split-root, `status --boot` ORPHANS cross-references each entity's `worktree:` against the filesystem and git state, but resolves the repo root via `FindGitRoot(roots.entityDir)` (`internal/status/handlers.go:194`) — which points at the `.spacedock-state` checkout, NOT the code repo where the worktree actually lives. So a present code worktree is reported `dir_exists=no` / `branch_exists=no` (a false orphan). Confirmed live this session: boot flagged the `7h release-notes-local-summary` code worktree as both-no.

Root cause is the same `entityDir`-vs-`definitionDir` class that `f2` (#252) fixed for `scanMods`: the orphan scan must resolve `FindGitRoot(definitionDir)` (the code repo root) for the worktree dir/branch existence checks, while entity state stays in the state checkout. This is the serialized `internal/status` lane.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — Split-root ORPHANS reports real worktree existence.** A present code worktree under a split-root workflow reports `dir_exists=yes`/`branch_exists=yes`; a removed one reports `no`. Single-root workflows are unaffected.
Verified by: a boot test over a split-root fixture with a live code worktree (the missing case that let this ship). Extend `internal/status/boot_orphan_abs_test.go` to the split-root code-root path; the new case must fail against the current `entityDir` resolution and pass after switching to `definitionDir`.
