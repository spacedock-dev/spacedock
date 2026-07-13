---
title: Restore lazy loading for first-officer merge and write cores
status: ideation
source: Captain correction after fresh-session boot trace, 2026-07-13
started: 2026-07-13T15:56:42Z
completed:
verdict:
score:
worktree:
issue:
id: 1kevganrmr2csr539ktfjerh
---

Restore the intended first-officer loading boundary: boot reads the shared core and active runtime adapter, while the write core loads at the first FO-authored mutation and the merge core loads only at terminal or merge-mod recovery handling.

## Problem

The archived `dp` task (`dpwp415wfzj6yrcwbs0krrea`, `fo-deferred-load-point-hunt-vs-skill-addressing`) diagnosed a real Claude failure class in which the model hunted the filesystem for delayed references; its structural addressing direction fed PR #491. PR #495 (`f22360de`) was not `dp`: it explicitly superseded #491's root-addressing direction and bundled `6h`, `p4`, and the `m1y` small-core preload. That preload overcorrected by eagerly importing `fo-merge-core.md` and `fo-write-core.md` from the first-officer entry skill. Commit `6baeed70` moved merge eager and `1e4423e1` moved write eager; later contractlint tests canonized exactly three eager imports. This regressed the earlier split explicitly designed in `shared-merge-dispatch-contract`: an interactive greet that stops should not pay for mutation or terminal ceremony it never uses.

The failure was deterministic discovery, not laziness itself. Restoring vague skill discovery would recreate the original hunt. The delayed cues must name exact canonical reference paths and load them at behaviorally enforced triggers.

## Proposed approach

- Keep `first-officer-shared-core.md` and the selected host runtime adapter eager at boot.
- Remove `@references/fo-write-core.md` and `@references/fo-merge-core.md` from `skills/first-officer/SKILL.md`.
- Restore an exact-path deferred cue for `references/fo-write-core.md` before the first FO-authored file write or state mutation. The cue must run `write.classify` before mutation and must not rely on model skill discovery, a compatibility shim, or filesystem search.
- Restore an exact-path deferred cue for `references/fo-merge-core.md` at terminalization and when resuming a merge `mod-block`. Dispatch and non-terminal gates do not load it.
- Replace the eager-topology tests introduced by #495 with trigger/reachability tests that protect both prompt economy and deterministic resolution. Preserve the resident smallest-sufficient rule and the runtime adapter's eager boot load.

## Out of scope

- Deferring the active runtime adapter.
- Reintroducing the removed standalone write-core skill wrapper or the rejected callable-skill discovery path.
- Changing merge, mutation, gate, or dispatch behavior beyond reference load timing.
- Deferring `fo-dispatch-core`, status viewer, gate presentation, feedback routing, or recovery differently from their existing triggers.

## Acceptance criteria

**AC-1 (VALUE) - A fresh interactive first-officer boot that greets and stops loads the shared core and selected runtime adapter but does not read merge-core or write-core.**
Verified by: fresh Claude and Codex shallow-boot scenarios record the files read before the greeting; both terminal/write references are absent, the runtime adapter is present, the greeting and boot state remain correct, and prompt bytes decrease against the pre-change eager baseline.

**AC-2 - The first FO-authored mutation deterministically loads the canonical write-core before any write classification or mutation, without broad search or skill-discovery guessing.**
Verified by: a live/fixture mutation journey observes the exact reference read before `status --set` or file mutation and succeeds with no `find`, recursive search, missing-skill retry, or pre-load write. A planted removal or post-write reorder turns the scenario red.

**AC-3 - Merge-core loads at the first terminal or merge-mod recovery boundary and nowhere earlier.**
Verified by: one non-terminal dispatch journey reaches a gate without reading merge-core; terminal and `mod-block=merge:*` journeys read the exact canonical reference before merge guard/hook handling and preserve existing archive, teardown, and refusal behavior.

**AC-4 - Contract tests enforce one eager canonical reference plus two deterministic delayed reference cues instead of the regressed three-import topology.**
Verified by: replace `TestFirstOfficerEagerReferencesKeepDispatchCoreDeferred` and related byte/closure assumptions with a topology test requiring only `@references/first-officer-shared-core.md`, exact resolvable delayed merge/write paths at their triggers, and the selected runtime adapter instruction. Planted eager merge/write imports, missing cues, dangling paths, and a restored wrapper each fail.

**AC-5 - Loading-boundary repair preserves first-officer behavior across supported hosts.**
Verified by: focused contractlint and live shallow-boot/mutation/terminal scenarios pass for the affected hosts, followed by `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; detached adversarial audit proves a trigger-breaking edit is caught.

## Test plan

Start with RED topology and boot-read assertions against the current three-import state. Update the entry skill and shared-core deferred load points without changing the referenced bodies. Add read-order fixtures for first mutation and terminal/mod-block handling, including broad-search negatives reproducing #495's original failure class. Run focused contractlint and host live scenarios before the full/race gates. Because shipped contract scaffolding changes, finish with a detached audit on a throwaway checkout.
