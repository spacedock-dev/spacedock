---
session-date: "2026-06-08"
sequence: 1
first-commit: 7a732b31
last-commit: 41994596
duration: ~2h45m
sprint: 0198-pre-flip-hardening
deliverable: spacedock 0.19.8 (cut + released)
---

# Session Debrief — 2026-06-08 (sprint 0198-pre-flip-hardening)

Commander session: drove the 4 sprint members from dispatch through validation, merge, the pre-cut antipattern audit, and the **0.19.8 release cut**. All four shipped PASSED; 0.19.8 is tagged and released.

## Shipped
- **kb** `migration-check-prune-state-walk` — [#327](https://github.com/spacedock-dev/spacedock/pull/327). Migration-check test tree-prunes the gitignored `.spacedock-state` checkout (matching the 3 sibling prunes) + deletes orphaned survey-scaffold fixtures.
- **qa** `spacedock-binary-missing-install-journey` — [#328](https://github.com/spacedock-dev/spacedock/pull/328). doctor + launch-gate messages now show display versions, not "contract N" jargon; brew-led stale-binary remedy; FO binary-absent abort points at the `spacedock claude` payoff.
- **z9** `codex-plugin-auto-install` — [#329](https://github.com/spacedock-dev/spacedock/pull/329). `spacedock codex` auto-installs a missing plugin then launches (codex analog of #311), channel-tracked via shared `devBranch`. Clean detached adversarial audit.
- **vh** `survey-skill-correctness-pass` — [#330](https://github.com/spacedock-dev/spacedock/pull/330). Survey corrected to agentsview's real git-root-basename model (resolves preflight B2) + codex-presence count, SCAFFOLD state-the-fact, sandbox probe; consolidates 69/1p/4t. +1 feedback cycle.

## Filed (backlog)
- **47rx3x8a809wx35vx6rbqqhv** `survey-codex-and-sandbox-followups` — 0.19.9 follow-ups: surface Codex sessions in the survey *body*, the silent-0 sandbox caveat, and a source-readability detector. (Captain pulled into 0199.)

## Non-PR commits (workflow-only)
- `41994596` release: bump to spacedock@0.19.8 — annotated tag `v0.19.8`, re-cut after a lightweight-tag botch.
- `7dd35c62` roadmap(0198): pre-cut antipattern audit — SHIP-CLEAR (0 blockers, 4 record-for-next).
- vh feedback cycle 1 recorded in entity body; `fo-friction-log.md` entries (rtk stale git output; kept-alive ensign worktree race).

## Decisions
- **vh pulled into 0.19.8** (was 0199) — B2 resolved by implementing the corrected git-root model.
- **Merge-gated on the offline CI check**, not the live E2E lanes (the 4 changes don't touch the shared runtime scenarios).
- **Deferred to 0.19.9:** survey silent-0 caveat + Codex body-surfacing.
- **DoD#3 (qa live-drive) → post-release** on another machine running released 0.19.8; issues → 0.19.9.
- **Pre-cut antipattern audit run before the tag** (captain-directed gate) — SHIP-CLEAR.

## Issues — Workflow
- **Kept-alive-ensign worktree race (qa):** impl ensign committed late comm-officer polish after Done while validation ran in the same worktree. Fixed by **hold-before-advance** (confirm no pending polish before validation) — applied cleanly on vh.
- **vh feedback cycle 1 (stale-DB phantom-0):** codex-presence rendered a confident "0" from a stale persisted DB; §1 now mandates a fresh sync + a non-vacuous sync→codex-presence e2e test (the gap fixture-only AC-2 missed).

## Issues — Spacedock
- **rtk-filtered git output was stale**, hiding a just-landed commit / giving a bogus "files identical" — caught only by raw blob-SHA compare. Hazard for commit-identity checks in a live, sibling-mutated worktree.
- **Release-process gotcha (operator):** cut a lightweight tag without consulting `docs/releasing.md` → `release.yml`'s `git cat-file tag` failed; re-cut as an annotated tag with changelog. (Already addressed on next: the cut-the-release checklist now points the Commander at `docs/releasing.md`.)

## Observations
- **hold-before-advance** cleanly prevented the worktree race on vh — worth codifying in the FO/ensign contract.
- **Pre-cut audit** is the right gate — independent, before the tag, when an antipattern can still be caught.
- **B2 resolved by construction** — driving vh closed the captain-gated blocker with no separate decision.
- **Read the runbook before outward ops** — the lightweight-tag botch was avoidable.

## What's Next
- **DoD#3** — captain's post-release qa live-drive on another machine (version-bearing doctor/gate messages); issues → 0.19.9.
- **0.19.9** — `survey-codex-and-sandbox-followups` + the audit's 4 record-for-next findings (kb test walk-loop dup, deferred-not-archived survey supersessions, a pre-existing ensigncycle instruction-file read, a doctor exit-0-on-no-plugin note).
- **0199-pre-flip-mechanics** (parallel track) — k6 devBranch retarget, v3/th/jm/m1; the 0.20.0 flip.
