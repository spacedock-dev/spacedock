---
id: pt0mz5stt4c1ve7ynz24r3yv
title: docs site — address reviewer feedback (simplify, dedupe, fix rendering)
status: implementation
source: "PR #343 intake (captain 2026-06-12). Addresses reviewer (Karen) feedback on the docs site — headline asks: the docs are wordy / Spacedock feels complex, and apply the Recce doc-writing principle. Intaken directly at implementation for captain-steered interactive revision on the existing PR branch."
started: 2026-06-12T17:34:58Z
completed:
verdict:
score:
worktree: .worktrees/docs-site-feedback
pr: "#343"
issue:
---

Intake of **PR #343** (`docs-site-feedback` → `main`) for continued, captain-steered revision. The PR already does two things; this entity tracks refining it to done.

## Problem

Reviewer feedback on the docs site: the writing is wordy and Spacedock reads as complex, and the docs should follow the Recce doc-writing principle (structure + simplicity). Three issues were visible in the rendered site (e.g. a dark-on-dark header in slate mode; the Home page leading with the "multi-agent orchestrator" claim and ~7 terms up front).

## Approach (in-flight on the PR; revision is captain-steered)

PR #343 already:
1. Adds `docs/site/CLAUDE.md` — a standing authoring directive (adapted from Recce's `doc.md`) governing structure + simplicity for everything under `docs/site/`, auto-loaded when editing docs, deferring to `voice-and-tone.md` for voice, and excluded from the published build via `exclude_docs`.
2. Acts on each feedback item — fixes that needed no decision plus ones gated on a maintainer call (resolved with the maintainer).

Remaining work is **interactive**: the captain steers the doc voice/content revisions directly with the ensign. The ensign applies that steering on the `docs-site-feedback` branch (so #343 updates), keeps the rendered docs correct, and does not regress the build.

## Out of scope

- Opening a new PR or new branch — revisions land on the existing `docs-site-feedback` branch.
- The docs-deploy mechanism (the landing site owns deploy; the mkdocs Pages job is descoped) — this is content/voice, not deploy.

## Acceptance criteria

(Captain-steered; validation confirms. Each verified by the rendered docs / build, not a prose-grep.)

**AC-1 (sketch) — the reviewer feedback items are addressed in the rendered docs.** Verified by: the rendered site shows the simplified Home (problem-first, fewer terms) and the resolved rendering issues (e.g. header legible in slate mode) — observed in the built site / before-after, not the source prose alone.

**AC-2 (sketch) — the authoring directive exists and is excluded from the published build.** Verified by: `docs/site/CLAUDE.md` present and absent from the built `site/` output (`exclude_docs` honored) — observed in the build output.

**AC-3 (sketch) — the docs build stays green.** Verified by: `mkdocs build --strict` exits 0 after the revisions.

## Notes

Interactive intake — the captain steers revision in the ensign's pane directly. SSH push is currently down; the ensign commits locally on `docs-site-feedback`, and the FO pushes the updated PR branch (via the gh-HTTPS route or once the key is restored) when the captain is satisfied.

## Progress (session 1, 2026-06-12)

Phase 1 complete and committed on `docs-site-feedback` (19 commits ahead of main, working tree clean): structural restructure + tree-wide tone sweep + all of Karen's 2026-06-12 feedback. See the handoff in the next-session prompt for the full done-list and the Phase-2 plan.

Key protocol notes for the next session:
- `mkdocs build --strict` (AC-3) is **deliberately unenforced during iteration** per captain; build with plain `mkdocs build`, re-enable `--strict` at the end and fix the link punch-list (several inbound links now point at GitHub URLs / changed anchors after the Contributing/Advanced/Reference restructure).
- comm-officer polish is **best-effort and non-blocking** (2-min timeout, then proceed un-polished). Never block on it. `SendMessage` takes ONLY `to`/`summary`/`message` — adding `type`/`recipient`/`content` is what corrupted earlier sends.
- Phase 2 (per-page paragraph polish) NOT yet started on disk. comm-officer returned polish for `index.md` and `install.md` that is not yet applied — re-derive or re-request in the fresh session.
