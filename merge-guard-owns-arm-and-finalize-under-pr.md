---
title: merge guard owns the full merge:pr lifecycle — auto-arm + finalize-from-detected-merged (eliminate FO hand-rolling)
status: validation
sprint: 0221-layered-fo
group: binary-ux
id: xdcf177r7sqtkb9w3mdafp4w
worktree: .worktrees/spacedock-ensign-merge-guard-owns-arm-and-finalize-under-pr
started: 2026-06-20T06:38:47Z
mod-block: merge:pr-merge
pr: "#415"
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

## Stage Report: implementation

- DONE: AUTO-ARM — when entering terminal under merge:pr with an empty mod-block, set mod-block=merge:{hook} and signal `armed`
  `internal/status/merge.go` auto-arm now fires under BOTH policies (removed the `policy == mergeLocal` gate); `TestMergeGuardAutoArmsUnderBothPolicies` green. (Arm value is `merge:{registered-hook}`; in production the merge mod is `pr-merge` → `merge:pr-merge`.)
- DONE: FINALIZE — when the PR is detected MERGED, clear mod-block then terminalize then archive, EVEN from a non-armed state
  New `prIndicatesMerged()` predicate: finalize keys off a merge sentinel in `pr` (`pr-merge:{n}` / `local-merge:{sha}`), never raw pr-presence. `TestMergeGuardFinalizesFromMergedSentinelNonArmed` (stranded empty-mod-block case) green.
- DONE: the archive move is committed PATH-SCOPED by the verb (no rename mis-stage / no bare git add -A)
  New `commitArchiveMove()` stages only the entity's two rename paths and commits them; `TestMergeGuardFinalizeCommitsArchivePathScoped` proves a dirty sibling is NOT swept in.
- DONE: TDD — AC-1 auto-arm, AC-2 finalize-from-non-armed-MERGED, AC-3 path-scoped archive commit (failing tests first, then to green)
  All three written red first (confirmed failing), then implemented to green. Fixtures `080-pr-merged.md` (merged sentinel, non-armed) + `090-pr-open-unmerged.md` (open #42, must block).
- DONE: premature-finalize gate (Codex FO steering) — bare/open PR must block, not finalize
  `TestMergeGuardBlocksOnOpenPRNoModBlock` pins it: open `#42` + empty mod-block → `blocked`, NOT archived.
- DONE: `go build ./...` and the merge-package tests green
  Full `go build ./...` clean; `go test ./internal/status ./internal/cli` green; full `go test ./...` green.
- DONE: Update `skills/first-officer/references/fo-merge-core.md` to the one-verb-per-phase flow; keep mod-block guard semantics
  Rewrote Merge-and-Cleanup + the «merge.guard» entry to armed/blocked/finalized-on-sentinel; added the SPACEDOCK_BIN launcher-invariant note. Mod-block + merge-hook guard sections unchanged. `TestProseFunctionNotationBindsToRouting` + contractlint green.

### Summary
The verb now owns the full merge:pr lifecycle: auto-arm under both policies, finalize off a local merge SENTINEL (never raw pr-presence — fixing the premature-finalize bug where an open-PR entity was archived before its PR landed), and a path-scoped archive commit by the verb. Two superseded policy-gated tests were reframed (not deleted) to the new contract while preserving their guard-not-bypassed intent. CONSISTENCY HANDOFF: fo-merge-core.md now names the `pr-merge:{n}` sentinel as the FINALIZE key, but the fo-realm `docs/dev/_mods/pr-merge.md` MERGED-detection still finalizes directly rather than recording that sentinel + re-running `merge guard` — left for the FO (fo-realm file, outside this stage's stated deliverable; the verb accepts both sentinel forms today).

## Stage Report: implementation (cycle 2)

- DONE: prIndicatesMerged finalizes ONLY on a well-formed merge sentinel — a positive-integer pr-merge:{N} or a non-empty local-merge:{sha}; pr-merge:abc, pr-merge:0, and empty/quoted-empty sentinels must NOT finalize. Write the red test first (proving the current fail-open), then implement to green.
  Red first: `TestPRIndicatesMerged` + `TestMergeGuardDoesNotFinalizeOnMalformedSentinel` (fixture `100-pr-malformed-sentinel.md`, `pr: pr-merge:abc`) confirmed FAILING against the bare-HasPrefix code (it finalized+archived `abc`). Fix: `prIndicatesMerged` now `strconv.Atoi`-validates a positive pr-merge suffix and `isSHALike`-validates a non-empty hex local-merge suffix; all other forms (garbage/zero/empty/quoted-empty/bare-#N) return false. Both tests green; happy-path `pr-merge:42`/`local-merge:{sha}` still finalize. Code commit 98769858 (post-rebase).
- DONE: Branch rebased onto current origin/main (now 2 behind) and full `go test ./...` green.
  `git fetch origin && git rebase origin/main` — clean (the 2 ahead were the "fail closed on invalid validate roots" pair `c32473dd`/`3e2cb0a7`, no merge.go overlap). Full `go test ./...` green (all packages ok; status 27s, cli 43s).

### Summary
Closed the fail-OPEN hole the detached audit found: `prIndicatesMerged` (internal/status/merge.go) did a bare `strings.HasPrefix` with zero suffix validation, so any garbage sentinel drove an irreversible finalize+archive. The smallest reasonable fix validates the suffix after the prefix — positive integer for `pr-merge:`, non-empty hex token for `local-merge:` — failing CLOSED on everything else. Auto-arm and archive-commit logic untouched (out of scope); the folder-form archive path-awareness and mod↔core sentinel convergence remain separately backlogged.

### Feedback Cycles

#### Detached adversarial audit (validation, FO-run via parallel subagent workflows; mutation-probed, tests executed)

Two rounds, per the validation stage's high-stakes detached-audit-on-throwaway requirement.

- **Round 1 — pre-fix (985e9f82 → rebased 04a9362b), 3 lenses.** Verdict CONCERNS → GO_WITH_FIXES.
  - **MATERIAL (CLOSED in cycle 2):** `prIndicatesMerged` bare `strings.HasPrefix`, zero suffix validation → `pr-merge:abc` / `pr-merge:0` / quoted-empty drove a full irreversible finalize+archive (fail-OPEN). Adversarial edit drove garbage sentinels end-to-end through MergeGuard and confirmed the archive. Closed by commit 98769858 (suffix validation, red-first TDD).
  - note1 reframed tests: **SOUND** — mutation probe (forcing `--force` into emitSet) turned `TestMergeGuardArmedButNoMergeRefusesFinalize` RED, revealing a dual guard (set + archive); the never-bypass invariant is relocated, not removed.
  - note2 mod↔core sentinel: **REAL but BENIGN** — bare `#N` fail-safe blocks, paths mutually exclusive (resolveEntityPath misses _archive) → backlog.
  - CONCERN C: `commitArchiveMove` hardcodes flat `slug.md`, would exit-128 on folder-form entities → backlog (dogfood entity is flat-form, unaffected).
- **Round 2 — post-fix (throwaway checkout @ 98769858), 2 lenses.** Verdict SOUND → GO. No fail-CLOSED regression: `pr-merge:{N}` + `local-merge:{short-sha}` still finalize; all status+cli tests green; never-bypass guard + bare-#N block untouched. Residual benign fail-OPEN (`pr-merge:+42` / `0042` via Atoi sign/leading-zero leniency; unbounded/uppercase isSHALike) — still names a real positive PR, never emitted by the honest hook → backlog hardening.

Backlog from the audits (none cut-blocking): folder-form `commitArchiveMove` path-awareness; mod↔core sentinel convergence; canonical-sentinel hardening (reject sign/leading-zeros, bound SHA length); archive-commit robustness vs the non-deterministic validate pre-commit hook.

## Stage Report: implementation (cycle 3)

Closes three audit-backlog items (folder-form, mod↔core convergence, archive-commit robustness) before PR #415 merges. Captain-approved.

- DONE: commitArchiveMove handles folder-form (folder-form fixture + test proves finalize+archive+path-scoped commit); flat-form byte-identical.
  `commitArchiveMove` now resolves the rename pathspecs via `archiveMovePathspecs`, mirroring `runArchive`: folder-form renames the whole `{slug}/` dir, flat-form the `{slug}.md` file. `TestMergeGuardFinalizesFolderFormEntity` (fixture `110-pr-merged-folder/index.md`) was red first (`exit status 128` on the `git add`), now green; the flat-form `TestMergeGuardFinalizeCommitsArchivePathScoped` stays green (byte-identical path).
- DONE: pr-merge mod MERGED-detection (startup + idle) converges to: set pr-merge:{N} sentinel + re-run `merge guard --verdict passed`; end state identical.
  `docs/dev/_mods/pr-merge.md` (fo-realm) startup + idle MERGED paths rewritten to record `pr=pr-merge:{N}` then finalize through the verb (clear+terminalize+archive+commit), replacing the hand-rolled two-`--set`+`--archive`. CLOSED-without-merge escalation unchanged; noted pr now records the sentinel post-finalize (was bare `#N`). Version 0.12.3→0.12.4. Prose-only (no Go test); contractlint + full suite green.
- DONE: finalize is atomic: a failing archive commit rolls back to pre-finalize state (live, non-terminal, mod-block intact); test injects a failing commit and asserts unchanged.
  `finalize` snapshots the pre-finalize bytes + live location before mutating (`captureArchiveState`); on `commitArchiveMove` failure it reverses the rename and restores the original content (`rollbackArchive`). `TestMergeGuardFinalizeRollsBackOnCommitFailure` (real failing pre-commit hook injected into the staged repo) was red first (entity stranded in `_archive`), now green — entity byte-identical to pre-finalize. Rollback scoped to commit failure only, so the never-bypass guard-refusal behavior (`TestMergeGuardArmedButNoMergeRefusesFinalize`) is unchanged.
- DONE: rebase + full suite + push.
  Rebased clean onto origin/main @ `d4b7ac61` (#416 had landed). Full `go test ./...` GREEN (all packages; status 26s, cli 39s, dispatch 27s). Code branch force-with-lease pushed (rebase rewrote hashes) — PR #415 updated. Single code commit `52380e72` (pre-rebase) carries all three fixes.

### Summary
Three fixes folded into PR #415: (1) folder-form-aware archive commit (mirrors runArchive resolution; flat-form byte-identical); (2) atomic finalize that rolls back the rename + frontmatter on a commit failure rather than stranding the entity half-archived; (3) the pr-merge mod's MERGED-detection now finalizes through the `merge guard` verb (dogfooded == shipped path), end state identical with the pr field recording the sentinel. All TDD red-first. Residual canonical-sentinel hardening (Atoi sign/leading-zero leniency, unbounded SHA) remains backlogged — benign, never emitted by the honest hook.
