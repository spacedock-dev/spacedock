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

Treat the configured value `B` as a branch **name**, never as an already-qualified ref. One small resolver returns `H = refs/heads/<B>` plus the diagnostic remote-tracking ref `T = refs/remotes/origin/<B>`. It does not normalize, strip, or special-case a leading `refs/heads/`, and it does not introduce a generalized ref model.

Before any mutating filesystem or Git operation, run both read-only Git validations as direct argv: `git check-ref-format --branch B` proves `B` is safe for branch-name-taking commands, and `git check-ref-format H` proves the constructed full ref is valid. Both must pass. This dual validation serves AC-2. Validating only `H` is simpler but insufficient: on the spiked Git, `refs/heads/-foo` is a valid full ref while `-foo` is not a valid branch name and is option-like. Validating only `B` is likewise insufficient because remote operations consume `H`, not `B`.

Feed those values into the existing Git sequence according to the argument's meaning:

- only after both validations, branch-name-taking creation and attachment remain `checkout --orphan B` and `worktree add <path> B`, followed by an exact `symbolic-ref HEAD == H` guard;
- local existence and checkout verification use `H`;
- remote detection selects the exact `H` line from `ls-remote --heads origin H`;
- a fresh fetch uses the explicit refspec `H:H`, a present-checkout fetch/pull uses the exact remote ref `H`, and push/re-push uses `H:H`;
- conflict diagnostics resolve `T`, avoiding revision-name shorthand.

This literal-ref resolver serves AC-1. The simpler alternative—continue passing `B` to each Git command—is insufficient because Git interprets the same bytes as a branch name in `checkout --orphan` but as the qualified ordinary branch in remote operations. Duplicating `"refs/heads/" + B` at each call is also insufficient because it invites another partial fix and cannot validate the mapping once before mutation.

Verify the checkout's full symbolic HEAD before staging/commit, pull/rebase, or push. This fail-closed guard serves AC-2. The simpler alternative—trusting the state checkout path—is insufficient because a valid checkout can be attached to a different branch and would then receive the commit or rebase.

## Consumer inventory and generated guidance

- `internal/cli/state.go` consumes the validated mapping for `state init`/`state new`: raw `B` only for `checkout --orphan` and `worktree add`, `H` for local/remote lookup and pull, and `H:H` for fetch/push.
- `internal/cli/state_sync.go` consumes it for `state ready`/`state commit`: exact-HEAD guards and peer diagnostics use `H`/`T`, pull uses `H`, and push/re-push uses `H:H`.
- `internal/dispatch/build.go` is an active mutation consumer today because `stateCommitGuidance` renders raw `git push origin B` and `git pull --rebase origin B`. Replace the manual add/commit plus raw sync recipe with one launcher-bound command: `"${SPACEDOCK_BIN:-spacedock}" state commit <slug> --workflow-dir <workflow-dir> -m "..."`. It must be the sole commit/sync action after the worker writes the entity, not a command appended after a manual commit; `state commit` already owns path-scoped add/commit, index-lock retry, origin/no-origin behavior, non-fast-forward rebase, and conflict halt, while this task adds the dual validation and literal mapping inside that executable path.
- `internal/dispatch/reconcile.go` calls `StateBranch` only to display the configured name in `state sweep` JSON. It never supplies the value to Git, so it preserves literal `B` and needs no refspec behavior.

Routing generated guidance through `state commit` serves AC-1 and removes the raw-branch escape hatch. Rendering explicit `H:H`/`H` commands in the prose is the simpler local edit but is insufficient because it duplicates the resolver and guard contract in generated text and can drift again. Appending `state commit` after the current manual commit is also insufficient: its clean-entity no-op returns before remote sync, so the binary must replace the manual commit/sync recipe rather than follow it.

Keep the current command flow and exit conventions. Do not add locks, journals, generations, publication protocols, daemons, leases, recovery controllers, lifecycle supervisors, or automatic repair. This scope boundary serves AC-1 by correcting the addressed history without creating a second synchronization system.

## Git mechanism spike

A disposable real-Git fixture on `git version 2.39.5 (Apple Git-154)` proved the riskiest ref/DWIM boundary before design completion. Compact evidence: `checkout --orphan B` plus `fetch origin H:H` created/attached `H` at `05b36503559d`; after a peer advanced `H` to `17ee0fe9f890`, `pull --rebase origin H` produced local `17ee0fe9f890`; `push origin H:H` advanced the literal remote to `1a784fbbcefa`, while ordinary `refs/heads/foo` stayed `afe25ec7857d` before and after. The pull also updated `T` as required for peer diagnostics.

Three adversarial observations constrain the design. `ls-remote --heads origin refs/heads/foo` returned two lines, so the existing exact-line parser remains load-bearing. `worktree add <path> H` produced detached HEAD (`symbolic-ref` exit 128), while `worktree add <path> B` attached the literal branch even with an ordinary `foo` collision. Finally, `git check-ref-format refs/heads/-foo` exited 0 while `git check-ref-format --branch -foo` exited 128, proving both validations are necessary. No further spike is needed: the remaining mechanism is the already-proven exact symbolic-HEAD comparison from e6j, reused here without its linked-worktree recovery behavior.

## Acceptance criteria

- **AC-1 (VALUE):** In a real origin-backed fixture with `state-branch: refs/heads/foo`, local creation, `state ready`, pull/rebase, `state commit`, and push all address the literal branch whose full ref is `refs/heads/refs/heads/foo`; the ordinary remote branch `foo` remains untouched.
  **Test:** Drive `state new`, fresh-clone `state init`, peer advancement plus `state ready`, and an entity edit plus `state commit`; compare exact local/remote ref SHAs and snapshot ordinary `foo` before and after. Build a dispatch for the same workflow and assert its emitted sole commit/sync action is launcher-bound `spacedock state commit`, with no direct `git push origin` or `git pull --rebase origin` recipe.
- **AC-2:** Local-only and origin-backed fixtures prove the same configured override resolves to the same history; invalid branch-name/full-ref mappings and a checkout mismatch fail before any wrong-history or local setup mutation.
  **Test:** Assert both fixtures' symbolic HEAD and local full ref, then attach the state path to the ordinary branch and run `ready`/`commit`; local HEAD/index and both remote refs remain byte-for-byte unchanged. With `state-branch: -foo`, prove full-ref validation alone would pass, then assert `state new` exits nonzero while `.gitignore` bytes/existence, the state path, NUL-delimited worktree registry, and all refs remain identical; invalid `ready`/`commit` likewise make no mutation.
- **AC-3:** Existing behavior for `spacedock-state/dev` remains unchanged across ready and commit paths.
  **Test:** Retain the existing ordinary-name real-Git ready/commit suite, add a resolver table case proving `spacedock-state/dev` maps to `refs/heads/spacedock-state/dev`, and prove ordinary/no-origin generated guidance uses the same launcher-bound `state commit` command while preserving local-only behavior at execution.
- **AC-4:** Focused real-Git regressions, `go test ./...`, and `go test ./... -race` pass; the change reuses existing Git command flow and introduces no coordination or lifecycle subsystem.
  **Test:** The implementer runs focused literal/collision/mismatch/new/guidance tests, `gofmt -w ./cmd ./internal`, full, and race gates before requesting exact-head Roborev; diff inspection confirms only mapping, existing consumers, tests, and the promised doc change.

## Test plan

- **Focused real-Git CLI fixture (moderate, seconds):** extend `internal/cli` helpers to create a bare origin containing both ordinary `foo` and literal `refs/heads/foo`. Exercise the complete `new -> init -> ready -> commit` sequence and assert process exits, symbolic refs, exact SHAs, on-disk entity state, and untouched ordinary-branch state. This mechanism serves AC-1; argv-only unit tests are simpler but insufficient because Git's context-sensitive ref interpretation is the defect.
- **Fail-closed fixtures (moderate, seconds):** snapshot refs, HEAD, index, worktree status, and remote refs around wrong-branch `ready`/`commit`. In a separate `-foo` `state new` fixture, snapshot `.gitignore` bytes or nonexistence, state-path nonexistence, `git worktree list --porcelain -z`, and `show-ref` before/after. This serves AC-2; checking stderr or only ready/commit is simpler but insufficient because `state new` writes `.gitignore` before orphan birth today.
- **Generated-guidance fixture (low):** run `dispatch build` for literal, ordinary, and no-origin split-root workflows; inspect the emitted assignment as the user-visible output. It must contain the exact launcher-bound `state commit <slug> --workflow-dir <dir>` action and no direct remote Git sync recipe; CLI fixtures separately exercise that command's remote and local-only outcomes. This serves AC-1/AC-3; rendering explicit refspec prose is simpler but insufficient because it duplicates the executable safety contract.
- **Compatibility and gates (low):** retain ordinary-name ready/commit coverage, run focused CLI and dispatch packages, `go test ./...`, and `go test ./... -race`, with formatting first. This serves AC-3/AC-4; no live runtime workflow is needed because the claim is CLI/generated-command behavior against real Git, not host integration.
- **Detached adversarial audit (moderate):** after implementation is green, a validator uses a detached throwaway checkout at the exact implementation head. Reverting generated guidance to a raw-`B` sync recipe, bypassing branch-name validation before `state new`, and bypassing the symbolic-HEAD guard must each make its focused guidance/new/mismatch test fail; restoring each line must return it green. Only then request exact-head Roborev. This serves AC-1/AC-2; ordinary code review is insufficient for a high-stakes state-mutation and guard path.

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
+Values invalid either as a Git branch name or as the constructed full ref are
+rejected before `.gitignore`, worktree, commit, or ref mutation. Generated split-root
+worker guidance delegates commit and synchronization to `spacedock state commit`
+instead of constructing refspecs in prose.
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

## Stage Report: ideation (cycle 2)

- DONE: Inventory and design literal-safe behavior for the active `stateCommitGuidance` consumer.
  The consumer inventory replaces manual commit/raw-`B` sync prose with sole-action, launcher-bound `spacedock state commit` guidance and adds literal/ordinary/no-origin output regressions.
- DONE: Validate both branch-name semantics and constructed full-ref semantics before mutation, with an option-like-name regression.
  Dual read-only checks cover `B` and `H`; the recorded `-foo` spike proves `H` alone passes while branch-name validation exits 128.
- DONE: Add invalid-branch `state new` proof before its `.gitignore` write boundary.
  AC-2 now snapshots `.gitignore`, state-path absence, NUL worktree registrations, and all refs before/after the failing option-like configuration.
- DONE: Record Git version and compact command/SHA evidence for the mechanism spike.
  The spike section records Git 2.39.5, literal birth/peer/final SHAs, ordinary-branch before/after parity, exact ls-remote collision count, and detached-HEAD exit evidence.
- DONE: Preserve the accepted mapping and e6j scope boundary.
  `B` to `H`/`T`, `H:H` fetch/push, pull on `H`, validated raw `B` attachment, exact remote selection, ordinary compatibility, and strict e6j separation remain unchanged.

### Summary

Cycle 2 closes the staff-review gaps at generated guidance, dual ref validation, and `state new`'s pre-orphan `.gitignore` mutation boundary. The corrected plan keeps synchronization in the existing binary, expands durable real-Git and generated-output proof, and adds no recovery or lifecycle scope.
