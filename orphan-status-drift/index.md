---
id: at3jmf35cgt1ghgm53m91en8
title: status --boot ORPHANS misreports DIR_EXISTS/BRANCH_EXISTS for split-root code worktrees
status: ideation
source: spacedock-dev/spacedock#251 + FO boot (2026-06-01) — reproduced live this session: boot flagged the 7h code worktree as dir_exists/branch_exists=no
started: 2026-06-02T02:18:54Z
completed:
verdict:
score: "0.28"
worktree:
issue: spacedock-dev/spacedock#251
---

## Problem

Under split-root, `status --boot` ORPHANS cross-references each entity's `worktree:` against the filesystem and git state, but resolves the repo root via `FindGitRoot(roots.entityDir)` (`internal/status/handlers.go:229`). Under split-root, `roots.entityDir` is the `.spacedock-state` checkout, which is itself a git worktree of the code repo — so it carries a `.git` *pointer file*. `FindGitRoot` walks up from `startDir` and stops at the FIRST dir holding a `.git` entry (dir OR regular file, `path.go:83`), so `FindGitRoot(entityDir)` returns `.spacedock-state` itself, NOT the code repo root where the worktree actually lives.

`scanOrphans` then computes `dirPath := PyJoin(gitRoot, wt)` for a relative `worktree:` value, yielding `.spacedock-state/.worktrees/<name>` — a path that does not exist. DIR_EXISTS is `no`. The BRANCH_EXISTS check (`worktreePaths[realpathOf(dirPath)]`) compares that same wrong path against the real `git worktree list` entries, so it is `no` too. A present code worktree is reported `dir_exists=no` / `branch_exists=no` — a false orphan. Confirmed live this session: boot flagged the `release-notes-local-summary` code worktree (present ~7h) as both-no. Reproduced in a clean fixture against today's binary (see Spike).

The relative `worktree:` value is load-bearing for the bug: the existing parity test `boot_orphan_abs_test.go` uses an *absolute* worktree path, which `PyJoin` resolves by discarding `gitRoot` entirely (os.path.join absolute-component reset), so the wrong-root never bites. That is precisely why this shipped untested — the only boot orphan test could not exercise the broken path.

Root cause is the same `entityDir`-vs-`definitionDir` class that `f2` (#252) fixed for `scanMods` (`scanMods(entityDir)` → `scanMods(definitionDir)`, commit `0f1a185d`). The orphan scan must resolve the worktree dir/branch existence against the code repo root — `FindGitRoot(definitionDir)` — while entity state stays in the state checkout. `definitionDir` (e.g. `docs/dev`) walks up to the real code repo root; `entityDir` (the state worktree) does not.

## Approach

In `runRead` (`handlers.go:229`), change the boot-path git-root source:

```
-	gitRoot := FindGitRoot(roots.entityDir)
+	gitRoot := FindGitRoot(roots.definitionDir)
```

This `gitRoot` is computed once and consumed ONLY by the boot path (`gatherBoot`/`printBoot` at lines 300/307), which threads it solely into `scanOrphans` (`boot.go:172`). It is not used for any entity-state I/O: entities are loaded from `roots.entityDir` independently (`scanEntitiesActive`, line 287), and `gatherBoot` already takes `definitionDir`/`entityDir` separately for their own roles. So the change is scoped exactly to the ORPHANS worktree-existence resolution.

In single-root, `definitionDir == entityDir`, so the change is byte-identical — the same clean-cutover property the #252 `scanMods` fix relied on. No `--archived`/`--next`/`--table`/`--validate` path reads `gitRoot`, so they are unaffected.

`internal/status` is a single serialized writer lane; no other status work runs an implementation worktree concurrently, so this is a self-contained one-site change with no cross-stage coordination.

## Spike (riskiest unknown — exercised first)

Riskiest unverified mechanism: does a split-root fixture with a real git worktree reproduce the misreport against today's `entityDir` resolution AND flip to correct after switching to `definitionDir`? If the fixture topology didn't reproduce the wrong-root, the whole test plan would be invalid. Exercised before committing to the plan:

- Built a code repo with its own `.git`; a `docs/dev/README.md` declaring `state: .spacedock-state`; `.spacedock-state` materialized as a real `git worktree add` (so it carries a `.git` pointer file, matching the live topology); an entity with relative `worktree: .worktrees/feature-wt`; and a real code worktree at `<coderoot>/.worktrees/feature-wt`.
- `FindGitRoot(entityDir)` = `.spacedock-state` (stops at the pointer file); `FindGitRoot(definitionDir=docs/dev)` = the code repo root. Confirmed.
- **Today's binary** `--boot` over the fixture: `ORPHANS … DIR_EXISTS=no BRANCH_EXISTS=no` (the live misreport, reproduced).
- **Patched binary** (the one-line `entityDir`→`definitionDir` change) `--boot`: `DIR_EXISTS=yes BRANCH_EXISTS=yes`.

Both the dir-existence and branch-existence halves flip together, end-to-end through the real binary. The mechanism is proven; the spike's fixture topology seeds the implementation's first test. Spike binary edit reverted; working tree clean.

## Acceptance criteria

**AC-1 — Split-root ORPHANS reports real worktree existence against the code repo root.** Under a split-root workflow whose state checkout is a git worktree (own `.git` pointer), an entity with a relative `worktree:` pointing at a present code worktree reports `dir_exists=yes` / `branch_exists=yes`; when that worktree directory is removed it reports `dir_exists=no` (and `branch_exists=no`). The resolution targets `FindGitRoot(definitionDir)`, not `FindGitRoot(entityDir)`.
Verified by: a new boot test (test plan T-1) over a split-root fixture with a live code worktree — the missing case that let this ship. The test must FAIL against today's `FindGitRoot(roots.entityDir)` (asserting `no/no`) and PASS after switching to `FindGitRoot(roots.definitionDir)` (asserting `yes/yes`), and assert `no` for the removed-worktree variant.

**AC-2 — Single-root ORPHANS is unchanged.** A single-root workflow (`definitionDir == entityDir`, no `state:` field) with a relative `worktree:` pointing at a present code worktree still reports `dir_exists=yes`, and a non-existent one reports `no` — byte-identical to today.
Verified by: the single-root negative test (test plan T-2). It passes both before and after the fix (the change is a no-op when the two roots coincide), confirming no regression. The existing `TestBootAbsoluteWorktreeDirExists` (absolute-path parity) continues to pass.

**AC-3 — Change is scoped to the ORPHANS worktree-existence root resolution.** Only the boot-path `gitRoot` source in `handlers.go` changes; entity-state resolution (`scanEntitiesActive`, the `entityDir` argument threaded into `gatherBoot`/`validateWorkflow`/`--set`/`--archive`) still targets the state checkout, unchanged.
Verified by: the full `go test ./internal/status` suite stays green (no split-root state-I/O test regresses), and the diff touches exactly the one `handlers.go:229` line in non-test code.

## Test plan

All tests are Go unit tests driving the native runner in-process (`runNative`), matching `boot_orphan_abs_test.go`. Low cost (sub-second each); no live workflow or CLI-binary fixtures needed. The git-worktree fixture is the only new machinery.

- **T-1 (AC-1) — split-root live-worktree boot, present and removed.** Extend `internal/status/boot_orphan_abs_test.go`. Add a helper that builds the spike topology: `git init` a code repo, write `docs/dev/README.md` with `state: .spacedock-state` and a worktree-bearing stage, `git worktree add` the `.spacedock-state` checkout (so it carries a `.git` pointer), write an entity (`feature/index.md`) with relative `worktree: .worktrees/feature-wt`, and `git worktree add` the real code worktree at `<coderoot>/.worktrees/feature-wt`. Run `--boot --workflow-dir docs/dev`; assert the ORPHANS row is `DIR_EXISTS=yes BRANCH_EXISTS=yes` (use a `branch_exists`-aware reader extending the existing `orphanDirExists` cell parser, or read both cells via `strings.Fields`). Then remove the code worktree dir (and `git worktree prune`) and re-run; assert `DIR_EXISTS=no`. This case is `no/no` against today's `entityDir` resolution (the reproduced misreport) and `yes/yes` after the fix — it fails first, passes after. Native-only (split-root is an intentional native divergence; no `runOracle`, consistent with the other split-root tests).
- **T-2 (AC-2) — single-root negative, no regression.** A single-root workflow (no `state:`) with a relative `worktree:` at a present code worktree reports `yes`, and a non-existent worktree reports `no`. This asserts the byte-identical-when-roots-coincide property: it is green both before and after the fix. The existing `TestBootAbsoluteWorktreeDirExists` is the absolute-path companion and is left intact.
- **T-3 (AC-3) — suite stays green.** `go test ./internal/status` (the serialized status lane) passes in full after the fix, confirming no split-root state-I/O path regressed.

Authoring order (TDD): write T-1 first, run it, watch it fail with `no/no` against the current `entityDir` resolution (the right reason), then make the one-line `handlers.go` change, re-run, watch T-1 go `yes/yes`. T-2/T-3 guard the no-regression boundary.

## Stage Report: ideation

- DONE: AC pins the fix to resolving the worktree dir/branch existence checks against FindGitRoot(definitionDir) — the code repo root — NOT FindGitRoot(entityDir) (the state checkout); same entityDir-vs-definitionDir class f2/#252 fixed for scanMods. The boot test must FAIL against today's entityDir resolution and PASS after the switch.
  AC-1 pins the resolution to `FindGitRoot(definitionDir)` and requires T-1 to assert `no/no` (fail) on today's `entityDir` resolution and `yes/yes` after the switch; Approach section cites the parallel #252 commit `0f1a185d`.
- DONE: Test plan names the boot fixture: extend internal/status/boot_orphan_abs_test.go to a split-root workflow with a LIVE code worktree (present → dir_exists/branch_exists=yes; removed → no), plus a single-root negative confirming no regression. This is the missing case that let the misreport ship.
  T-1 extends `boot_orphan_abs_test.go` with a git-worktree split-root fixture (present→yes/yes, removed→no); T-2 is the single-root negative; the Problem section explains the absolute-path existing test could not exercise the broken path.
- DONE: Scope boundary stated: ONLY the ORPHANS worktree-existence root resolution in handlers.go changes; entity-state resolution still targets the state checkout. internal/status serialized lane — no other status work runs an impl worktree concurrently.
  AC-3 + Approach pin the change to the single `handlers.go:229` line consumed only by the boot path's `scanOrphans`; entity-state I/O (`scanEntitiesActive`, the `entityDir` arg) is unchanged; serialized-lane note recorded.

### Summary

`status --boot` ORPHANS misreports present split-root code worktrees as `dir_exists/branch_exists=no` because `gitRoot := FindGitRoot(roots.entityDir)` (handlers.go:229) resolves the state-checkout worktree (which carries its own `.git` pointer) instead of the code repo root — the same `entityDir`-vs-`definitionDir` class #252 fixed for `scanMods`. The fix is a one-line `entityDir`→`definitionDir` switch scoped to the boot path's `gitRoot`, byte-identical in single-root. I exercised the riskiest unknown first: a clean fixture with a real `git worktree` state checkout reproduces `no/no` against today's binary and flips to `yes/yes` against a patched binary, end-to-end — so the named extension to `boot_orphan_abs_test.go` will fail-then-pass as required, and the absolute-path existing test (which can't exercise the broken path) is why this shipped untested.
