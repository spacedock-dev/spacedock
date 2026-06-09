---
session-date: 2026-06-08
sequence: 1
first-commit: 85404bdd
last-commit: 4082b90a
duration: ~3h
---

# Session Debrief — Sprint 0200-flip (Commander drive) — 2026-06-08 #1

Cold-boot Commander drive of the 0200-flip **pre-flip** work to merged-on-`next` — `nzb` + `k6d` (the packaged pair) plus `cmx` (joined mid-drive from the Shaping FO). All three shipped, a CL-requested test fix landed, and the mandated pre-cut antipattern audit ran SHIP-CLEAR. **No release cut** (the flip, `pj`, is held for FO + captain).

## Shipped
- **nzb** `gate-release-on-e2e` — [#337](https://github.com/spacedock-dev/spacedock/pull/337). Release-time `e2e-gate` job that goreleaser `needs:`, binding the cut to a `conclusion:success` live-e2e run whose `headSha` equals the tagged commit (audited `SPACEDOCK_E2E_GATE_WAIVER`).
- **cm** `frontdoor-launch-banner-ux` — [#338](https://github.com/spacedock-dev/spacedock/pull/338). Pre-launch banner rewritten to one consistent first-officer framing (no self-serve `spacedock status`); banner added to `runPi`; `Version` defaulted to the `dev` sentinel.
- **k6d** `two-channel-release-devbranch-stamp` — [#339](https://github.com/spacedock-dev/spacedock/pull/339). Per-channel `devBranch` stamp (stable→`main` / edge→`next`), two casks (`spacedock` + `spacedock@next`), single-target `release.yml` stamp swap, and a binary-pair channel-agreement guard.

## Filed (backlog)
- None newly filed this session. (`cmx` was *carved into* the sprint by the Shaping FO mid-drive — `39f746d6` — not newly filed.)

## Non-PR commits (workflow-only)
- [#340](https://github.com/spacedock-dev/spacedock/pull/340) `fix(dispatch): make launcher-command test hermetic to ambient SPACEDOCK_BIN` — test-only; CL-requested fix of a pre-existing env-dependent failure (not a sprint entity).
- `docs/roadmap/0200-flip/post-sprint-audit.md` — the pre-cut antipattern audit document, committed with this debrief.
- All entity state transitions (dispatch/advance/done/archive/pr/mod-block) are routine split-root churn in `.spacedock-state`, rolled up in the PRs above.

## Decisions
- **Took the conn** (captain delegated approve + merge + CI authority). Merged all four PRs to `next` on judgment **without burning the env-gated live lanes** — `next` is unprotected, and these surfaces (release.yml / goreleaser / internal/release / frontdoor / a test) don't touch the runtime FO/ensign code the live lanes exercise.
- **Sequential on the shared file:** `nzb` and `k6d` both edit `release.yml`, so `nzb` landed first, `cmx` ran in parallel (disjoint surface), then `k6d` branched off the updated `next`. No conflicts by construction.
- **Staff-review angle corrected twice:** per-entity → sprint-wide → the roadmap's prescribed **pre-cut antipattern audit** (staff-eng persona, ship-blocker-vs-defer), after the captain flagged it and pointed to `roadmap/README.md` + the `0198/post-sprint-audit.md` template.
- **R1 (stamp-to-`main` landmine):** the audit downgraded it from candidate ship-blocker to record-for-next-sprint (bounded impact) — but its fix note is "before any pre-flip cut." Fix-now-vs-defer pending the captain's release intent.

## Issues — Workflow
- **R1** — `release.yml:203-217` stamp step retargeted to `main` and **armed on `next` now**, while `AGENTS.md:28` still says cut-from-`next`; a pre-flip 0.19.x patch would push a wrong version-stamp commit (0.12.1 code labeled 0.19.x) to the legacy `origin/main` (runtime-reproduced). Bounded (the flip's archive+force-replace absorbs it; default installs still resolve from `next`). Fix-now-or-defer pending the captain's call.
- **R2** — after the stamp moved to `main`, no workflow updates `next`'s `plugin.json` `version`; the edge version-panel freezes at 0.19.9; the `release.yml:195` comment ("rides the unchanged next-publish.yml") is imprecise. Cosmetic.
- **R3** — archived entities have empty `issue:` frontmatter; PR provenance (#337/#338/#339) lives only in git/GitHub, not queryable Spacedock state.
- Self-caught: on `cmx` the mod-block was cleared before `pr=#338` was recorded; the merge-hook guard correctly refused terminalization; the merged PR was recorded truthfully and the entity completed — no data loss. Lesson applied on `k6d` (recorded `pr=#339` right after PR creation).

## Issues — Spacedock
- **Process-doc gap (not filed — proposed as direct edits):** `roadmap/README.md`'s pre-cut-audit checklist line is under-specified vs its own preflight line (no antipattern catalog, no `→ post-sprint-audit.md` output target, and the artifact is missing from "Where things live"), AND the Commander's `dispatch-sprint-execution.md` package omits the pre-cut-audit step entirely — which is why the audit was not run until the captain asked. Symmetry fix + a packaging step that copies the Drive-phase checklist into the dispatch package were proposed.
- **FO-runtime note (not filed):** early in boot, a `git worktree add` ran while the working directory had drifted into the `.spacedock-state` sub-repo (a linked worktree sharing the object store), creating the worktree under the wrong tree. Self-corrected by relocating with absolute paths. Worth an FO-runtime habit: always absolute-path worktree creation, and never leave CWD inside the state checkout.

## Observations
- The **sprint-wide audit earned its keep**: R1 is a cross-cutting landmine *no per-task detached audit could see* — it only emerges from `AGENTS.md` (cut-from-`next`) × the retargeted stamp × the divergent legacy `main`. Validates the prescribed sprint-wide angle over a per-entity re-review.
- The flip hand-off is **proven green-by-construction**: the auditor applied pj's one marketplace-ref flip (+ the paired `marketplace_manifest_test.go` edit) in a throwaway worktree and the tri-surface guard un-skipped to PASS, with an inverse honesty check (drifted devBranch reds it).
- Each high-stakes surface got its mandatory detached adversarial audit at validation (nzb 7 edits, cmx 5, k6d 4) — every weakening caught RED, zero material findings; the pre-cut audit then re-confirmed a sample of those claims independently.

## What's Next
- **`pj` — the 0.20.0 flip + marketplace ref flip + cut** — SHIP-CLEAR to proceed per its captain-gated runbook. Owned by FO + captain, not the Commander.
- **R1 decision** — guard the stamp now (`fix-now-on-next`: refuse/no-op the stamp when not flipping) vs accept + a "no pre-flip cuts from `next`" constraint that pj's flip work cleans up.
- **R2 / R3** → next-sprint backlog (the Shaping FO's Close step folds the recorded findings forward).
- **Doc edits** — tighten `roadmap/README.md` + the dispatch-package template so the next Commander boots with the pre-cut-audit step (catalog + output target + no-cut nuance).
- **Concurrent:** the Shaping FO is driving `wv` (mkdocs-material docs-site) in the shared `.spacedock-state` checkout — interleaved state commits are theirs, not this drive's.
