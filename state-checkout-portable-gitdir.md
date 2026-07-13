---
title: Make the split-root state checkout operable from foreign path surfaces (Cowork sandbox)
status: ideation
source: live Cowork dogfood 2026-07-13 (survey-claude-cowork-runtime-detection run-1 session)
started: 2026-07-13T02:54:30Z
completed:
verdict:
score:
worktree:
issue:
id: qwp7vn4dnt0sx96wzz1wfy8b
---

Make every Spacedock-owned state operation work when the same repository tree is reached through a path surface different from where its state checkout was created. Observed live: a Cowork sandbox mounts the project at `/sessions/<name>/mnt/<folder>` while the linked checkout's `.git` file still points at the host-absolute administrative directory, so Git reports `not a git repository` even though the checkout and repository metadata are both present under the mounted project.

## Problem

By default, `git worktree add` records absolute paths in both the checkout's `.git` pointer file and the administrative directory's `gitdir` back-pointer. The split-root checkout lives inside the project (`docs/dev/.spacedock-state` backed by `.git/worktrees/<id>`), so a second path surface for the same tree — Cowork mount, bind mount, moved directory, or remounted disk — leaves all data reachable but makes those recorded paths stale. Observed live in Cowork: `spacedock state commit` failed, while an explicit `GIT_DIR=.git/worktrees/-spacedock-state GIT_WORK_TREE=docs/dev/.spacedock-state` commit succeeded on Git 2.34.

The persisted administrative name is not necessarily the checkout basename: Git may append a numeric suffix when names collide. Recovery therefore must preserve the administrative id parsed from the stale `.git` pointer and re-anchor it under the current repository common directory; it must not guess `<basename>` alone.

## Proposed approach

Two complementary lanes, with compatibility first:

1. **Verb-side gitdir re-derivation (ship first; no metadata rewrite).** Introduce one state-checkout Git runner used by every Git call made by `state commit`, `state ready`, and the existing-checkout branch of `state init`. On a linked checkout whose recorded gitdir exists, preserve today's normal `git -C` behavior. When that path is missing, parse the administrative id from the checkout's `.git` pointer, resolve the current repository common directory from the still-operable project checkout, validate `<common-dir>/worktrees/<id>` as the matching state-worktree administrative directory, and invoke Git with explicit `GIT_DIR` and `GIT_WORK_TREE`. Never rewrite either pointer as part of fallback. This works with old Git (the observed Cowork Git is 2.34), preserves collision-suffixed ids, and works in both host→mounted and mounted→host directions.
2. **Relative pointers at creation (structural, compatibility-gated).** Git 2.48 introduced `git worktree add --relative-paths` and `worktree.useRelativePaths`. Apply the same creation policy to both persisted state-checkout paths, `state new` and `state init`; do not fix only resumed checkouts. However, Git documents that this enables `extensions.relativeWorktrees` and makes the repository incompatible with older Git. Therefore it cannot be an unconditional version-of-the-creator check while Cowork may use Git 2.34. The implementation must make the compatibility policy explicit: keep absolute-pointer creation plus lane 1 wherever older consumers remain supported, or enable relative worktrees only when the project/runtime support floor is Git 2.48+. Prefer a capability probe over parsing a vendor-specific version string. Lane 1 remains required for existing checkouts regardless.

`state sweep` currently reads state files directly and uses Git only against the main project checkout. It should already survive this failure mode. Keep it in the portability fixture to prove the whole FO convergence sequence (`ready` then `sweep`, followed later by `commit`) works from the foreign path surface, rather than adding recovery logic where no state-checkout Git call exists.

## Out of scope

- Credential/push handling from sandboxes (separate concern; the verb's push-failure surfacing already covers it).
- Non-state (code) worktrees under `.worktrees/` — same mechanism could apply later, but state verbs are the broken path today.
- Repairing a project whose main checkout's own `.git` metadata is unreachable; compat recovery assumes the project checkout remains discoverable and only the nested state checkout has stale links.
- Making bare `git -C <state-checkout>` portable for legacy absolute-pointer checkouts. Lane 1 guarantees Spacedock's command surface; structural creation prevents the problem for checkouts whose Git compatibility policy permits it.

## Acceptance criteria

**AC-1 - The persisted checkout's administrative id is re-anchored safely instead of guessed from its basename.**
Verified by: focused tests with a collision-suffixed `.git/worktrees/<id>` entry. A missing recorded absolute gitdir resolves to the matching current common-dir entry; a missing, malformed, escaping, or mismatched candidate fails closed with a useful diagnostic rather than selecting another worktree.

**AC-2 - `spacedock state commit` completes through the compat runner from a foreign path surface without repairing repository metadata.**
Verified by: create a split-root fixture at path A, record both link files, move the entire project to path B so the recorded gitdir no longer exists, edit one task, and show `state commit` makes the path-scoped commit from B. The checkout `.git` file and administrative `gitdir` back-pointer remain byte-identical. Exercise local-only success and an origin-backed push path.

**AC-3 - `spacedock state ready` performs all existing-checkout Git operations through the same compat runner.**
Verified by: the moved-project fixture runs the real fetch/pull-ready path from B and observes the same ready/no-op, integration, and rebase-conflict exit semantics as an unmoved checkout. No fallback branch may silently classify an unreadable state checkout as `no origin`.

**AC-4 - The complete first-officer state sequence remains usable from the foreign path surface.**
Verified by: from B, run the public `state ready`, `state sweep`, and `state commit` commands in order. `sweep` produces its normal read-only result without unnecessary pointer recovery, and the other verbs use the compat runner. Then address the same project again at A (or an equivalent second mount) and repeat, proving the fallback is path-surface-derived rather than a one-way rewrite.

**AC-5 - Persisted state checkouts created by both `state new` and `state init` follow one explicit relative-worktree compatibility policy.**
Verified by: fixture-backed capability lanes for (a) a Git without `--relative-paths`, which retains absolute links and remains operable through lane 1 after a move, and (b) a Git with relative-worktree support under a declared Git ≥2.48 consumer floor, where both creation verbs produce relative `.git` and `gitdir` links and Git operates after a move without fallback. Assert that the relative lane enables `extensions.relativeWorktrees`; do not claim it remains readable by older Git.

**AC-6 - Existing unmoved repositories keep their current behavior and output.**
Verified by: existing `state init|new|ready|sweep|commit` fixtures remain byte-for-byte stable where output is contractual, plus focused tests showing the runner uses ordinary `git -C` when the recorded gitdir is valid. `go test ./...` and `go test ./... -race` pass.

## Test plan

Start with focused resolver/runner tests, including a collision-suffixed administrative id and fail-closed candidates. Add one fixture that creates the checkout at A and renames the whole repository to B; run the real `ready → sweep → commit` public-command sequence and assert durable state via process exits, entity body, state-checkout log, pointer bytes, and clean status. Cover both local-only and origin-backed paths plus the existing rebase-conflict contract. Test both `state new` and `state init` creation through injected old/new Git capabilities, including the `extensions.relativeWorktrees` compatibility consequence. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. A live Cowork smoke remains valuable after the fixture is green, using the same durable evidence rather than transcript wording.
