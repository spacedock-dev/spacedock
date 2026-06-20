---
title: merge guard owns the full merge:pr lifecycle — auto-arm + finalize-from-detected-merged (eliminate FO hand-rolling)
status: backlog
sprint: 0221-layered-fo
group: binary-ux
id: xdcf177r7sqtkb9w3mdafp4w
---

This session hand-rolled the terminal merge ceremony 5× (vk, mvc, ga, launcher-flag, b2): set `mod-block=merge:pr-merge`, commit, push the code branch, open/update the PR, set `pr:`, wait for merge, then clear `mod-block`, terminalize `verdict=PASSED worktree=`, archive, remove worktree+branch. `spacedock merge guard` was NOT usable end-to-end because of two documented gaps:

1. **No auto-arm under `merge: pr`.** The FO must hand-set `mod-block=merge:pr-merge` before invoking the hook; `merge guard` does not arm it.
2. **No finalize from a non-armed / detected-merged state.** `state sweep` detects a MERGED PR but won't finalize it unless armed; a re-validation bounce that clears `mod-block` strands the entity (observed live this session on a peer entity `status-validate-determinism`, left carrying `mod-block: merge:pr-merge`).

The hand-rolling is toil and error-prone (a newline-in-variable staging bug failed the archive commit twice; caught + fixed both, but the verb should own this).

## Fix
Make `spacedock merge guard <slug>` own the full `merge: pr` lifecycle so the FO invokes ONE verb per phase: (a) on entering terminal with no in-flight block, AUTO-ARM (`mod-block=merge:{hook}`) and signal `armed`; (b) on re-run, if the PR is detected MERGED, FINALIZE (clear `mod-block` standalone, terminalize `verdict=PASSED worktree=`, archive) even from a non-armed state; (c) path-scoped archive commit handled by the verb (no FO-side `git add` glob that can mis-stage a rename).

## Acceptance criteria
- **AC-1** — `spacedock merge guard <slug>` on a `merge: pr` entity entering terminal with an empty mod-block AUTO-ARMS (sets `mod-block=merge:pr-merge`) and reports `armed`. Verified by a test driving the verb on a fixture entity.
- **AC-2** — `spacedock merge guard <slug>` finalizes a detected-MERGED PR (clear mod-block + terminalize + archive) even when `mod-block` is empty (the stranded-non-armed case). Verified by a test seeding a merged-PR + non-armed fixture and asserting terminal+archived after one invocation.
- **AC-3** — the archive move is committed by the verb path-scoped (no rename mis-stage); `go test` for the merge package green; the FO contract's merge ceremony prose updates to the one-verb-per-phase flow.
