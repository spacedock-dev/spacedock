---
id: e6e81q69r89wtzjpm8reqgkk
title: README "see a real workflow" example link (post-flip)
status: backlog
source: "captain + nb (readme-main-flip-reconciliation) reconciliation 2026-06-07 — PR #315 linked docs/plans/ example/debrief paths in the README, but those paths are flip-fragile (docs/plans is the old workflow dir; paths change at the main flip). Add the link AFTER the flip so it's stable."
started:
completed:
verdict:
score:
worktree:
issue:
---

Add a single "see a real workflow" example link to the README so a newcomer can look at a concrete, real spacedock workflow in the repo.

## Problem

PR #315 added README links to example/debrief paths under `docs/plans/` — but `docs/plans/` is the pre-flip workflow dir and those paths shift at the `0.20.0` main flip, so a link added now goes stale. The reader-first README (nb) deliberately left this out to avoid shipping a flip-fragile path.

## Proposed approach

After the main flip, add ONE example link to a real workflow (a stable post-flip path) — enough for a newcomer to see a concrete workflow without a tour. Pick the path once the post-flip docs layout is settled.

## Out of scope

- Doing this pre-flip (the path would be flip-fragile).
- A full examples gallery — one good link.

## Acceptance criteria

Ideation/implementation fills in. Sketch: the README links one real workflow at a post-flip-stable path; a link-check (or the docs-site build, once `mkdocs-material-docs-site` lands) resolves it.

## Test plan

Link resolves (link-check / mkdocs `--strict` once the site exists).
