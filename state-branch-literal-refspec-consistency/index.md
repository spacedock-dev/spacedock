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

## Proposed approach

Treat the configured value `B` as a branch **name**, never as an already-qualified ref. One small resolver validates and returns `H = refs/heads/<B>` plus the diagnostic remote-tracking ref `T = refs/remotes/origin/<B>`. It does not normalize, strip, or special-case a leading `refs/heads/`, and it does not introduce a generalized ref model.

Feed those values into the existing Git sequence according to the argument's meaning:

- branch-name-taking creation and attachment remain `checkout --orphan B` and `worktree add <path> B`, followed by an exact `symbolic-ref HEAD == H` guard;
- local existence and checkout verification use `H`;
- remote detection selects the exact `H` line from `ls-remote --heads origin H`;
- a fresh fetch uses the explicit refspec `H:H`, a present-checkout fetch/pull uses the exact remote ref `H`, and push/re-push uses `H:H`;
- conflict diagnostics resolve `T`, avoiding revision-name shorthand.

This literal-ref resolver serves AC-1. The simpler alternative—continue passing `B` to each Git command—is insufficient because Git interprets the same bytes as a branch name in `checkout --orphan` but as the qualified ordinary branch in remote operations. Duplicating `"refs/heads/" + B` at each call is also insufficient because it invites another partial fix and cannot validate the mapping once before mutation.

Validate `H` with Git's `check-ref-format` immediately after configuration resolution, before any mutable file or Git command, and verify the checkout's full symbolic HEAD before commit, pull/rebase, or push. This fail-closed guard serves AC-2. The simpler alternative—trusting the state checkout path—is insufficient because a valid checkout can be attached to a different branch and would then receive the commit or rebase.

Keep the current command flow and exit conventions. Do not add locks, journals, generations, publication protocols, daemons, leases, recovery controllers, lifecycle supervisors, or automatic repair. This scope boundary serves AC-1 by correcting the addressed history without creating a second synchronization system.

## Git mechanism spike

A disposable real-Git fixture proved the riskiest ref/DWIM boundary before design completion. `checkout --orphan refs/heads/foo` produced symbolic HEAD `refs/heads/refs/heads/foo`; `fetch origin H:H`, `pull --rebase origin H`, and `push origin H:H` synchronized that literal branch while the ordinary remote `refs/heads/foo` SHA stayed unchanged. The pull also updated `refs/remotes/origin/refs/heads/foo` as required for peer diagnostics.

Two adversarial observations constrain the design. `ls-remote --heads origin refs/heads/foo` can return both the ordinary and literal refs, so the existing exact-line parser remains load-bearing. Conversely, `worktree add <path> H` creates a detached checkout, while `worktree add <path> B` attaches the literal branch (even with an ordinary `foo` collision); therefore the resolver must not blindly replace branch-name arguments with full refs. No further spike is needed: the remaining mechanism is the already-proven exact symbolic-HEAD comparison from e6j, reused here without its linked-worktree recovery behavior.

## Acceptance criteria

- **AC-1 (VALUE):** In a real origin-backed fixture with `state-branch: refs/heads/foo`, local creation, `state ready`, pull/rebase, `state commit`, and push all address the literal branch whose full ref is `refs/heads/refs/heads/foo`; the ordinary remote branch `foo` remains untouched.
  **Test:** Drive `state new`, fresh-clone `state init`, peer advancement plus `state ready`, and an entity edit plus `state commit`; compare exact local/remote ref SHAs and snapshot ordinary `foo` before and after.
- **AC-2:** Local-only and origin-backed fixtures prove the same configured override resolves to the same history, and a mismatch fails before any wrong-history mutation.
  **Test:** Assert both fixtures' symbolic HEAD and local full ref, then attach the state path to the ordinary branch and separately use an invalid configured ref; `ready` and `commit` must exit nonzero with local HEAD/index and both remote refs byte-for-byte unchanged.
- **AC-3:** Existing behavior for `spacedock-state/dev` remains unchanged across ready and commit paths.
  **Test:** Retain the existing ordinary-name real-Git ready/commit suite and add a resolver table case proving `spacedock-state/dev` maps to `refs/heads/spacedock-state/dev`.
- **AC-4:** Focused real-Git regressions, `go test ./...`, and `go test ./... -race` pass; the change reuses existing Git command flow and introduces no coordination or lifecycle subsystem.
  **Test:** The implementer runs focused literal/collision/mismatch tests, `gofmt -w ./cmd ./internal`, full, and race gates before requesting exact-head Roborev; diff inspection confirms only mapping, call-site, tests, and the promised doc change.

## Test plan

- **Focused real-Git CLI fixture (moderate, seconds):** extend `internal/cli` helpers to create a bare origin containing both ordinary `foo` and literal `refs/heads/foo`. Exercise the complete `new -> init -> ready -> commit` sequence and assert process exits, symbolic refs, exact SHAs, on-disk entity state, and untouched ordinary-branch state. This mechanism serves AC-1; argv-only unit tests are simpler but insufficient because Git's context-sensitive ref interpretation is the defect.
- **Fail-closed fixture (moderate, seconds):** snapshot refs, HEAD, index, worktree status, and remote refs around wrong-branch and invalid-ref runs. This serves AC-2; checking stderr alone is simpler but insufficient because mutation safety is the claim.
- **Compatibility and gates (low):** retain ordinary-name ready/commit coverage, run the focused package, `go test ./...`, and `go test ./... -race`, with formatting first. This serves AC-3/AC-4; no live runtime workflow is needed because the claim is CLI behavior against real Git, not host integration.
- **Detached adversarial audit (moderate):** after implementation is green, a validator uses a detached throwaway checkout at the exact implementation head. Reverting one remote call to raw `B` and separately bypassing the symbolic-HEAD guard must make the focused collision/mismatch tests fail; restoring each line must return them green. Only then request exact-head Roborev. This serves AC-1/AC-2; ordinary code review is insufficient for a high-stakes state-mutation guard path.

## Documentation change

Implementation applies this concrete clarification to `docs/specs/state-behavior-extension.md`:

````diff
@@
 The simplest extension is one top-level field:
 ```yaml
 state: .spacedock-state
+state-branch: spacedock-state/plans
 ```
+
+`state-branch` is optional; the default remains `spacedock-state/<workflow-directory>`.
+An explicit value is a literal Git branch name, not a pre-expanded ref. For example,
+`state-branch: refs/heads/foo` names the branch whose full ref is
+`refs/heads/refs/heads/foo`; state creation, ready, and commit all synchronize that
+same full ref.
````

## Relationship

This is a prerequisite for rebase and revalidation of `state-ready-cwd-path-resolution-reset`; that task remains frozen and unpushed until this shared consistency defect lands. This task does not change e6j's linked-worktree path anchoring, stale-registration detection/remediation, worktree-list parsing, local-only recovery, or linked/main README agreement. The eventual rebase may reuse this task's literal resolver and branch guard, but e6j owns recovery behavior and its own validation.

## Stage Report: ideation

- DONE: Define the smallest literal-ref mapping design, tie it to AC-1, and exclude lifecycle machinery.
  The design maps literal `B` once to validated `H`/`T`, preserves branch-name arguments, reuses the current Git flow, and explicitly excludes coordination/lifecycle subsystems.
- DONE: Specify real-Git proof for literal refs/heads/foo, local/origin parity, fail-closed mismatch, and ordinary branch behavior.
  The test plan drives real `new/init/ready/commit`, snapshots both colliding remote refs, covers local-only parity and wrong-branch/invalid-ref refusal, then runs focused/full/race gates before exact-head review.
- DONE: Keep e6j linked-worktree recovery separate and record the riskiest Git-mechanism spike or proven no-spike basis.
  The disposable spike proved exact remote refs/refspecs and the full-ref worktree-detach trap; the Relationship section leaves every e6j recovery behavior frozen and separate.

### Summary

Ideation now defines a minimal, behavior-first literal-ref mapping across the existing state Git commands, with an exact pre-mutation branch guard and no lifecycle machinery. The real-Git collision spike fixes the command-by-command boundary, and the implementation/validation plan requires durable ref evidence plus a detached adversarial audit before exact-head Roborev.
