---
id: 87j19afq4tj5te1hjvgd6rs4
title: pr-merge mod hardcodes base branch `next` (pre-flip); refit to `main`/config-driven
status: backlog
source: "0202 Commander drive (2026-06-13). The pr-merge mod (v0.12.1) opens PRs against `next`; the Commander overrode the base to `main` per the dispatch doc on every merge. Same post-flip stale-trunk class as dispatch reconcile."
group: cleanup
---

The `pr-merge` merge-hook mod (`docs/dev/_mods/pr-merge.md`, v0.12.1) opens code-branch PRs against base `next`, the pre-flip trunk. Post-flip the trunk is `main`, so every merge this drive required a manual base override.

## Problem

`_mods/pr-merge.md` describes the PR lifecycle against `origin`, base branch `next` (header + `gh pr create --base next`). The 2026-06-08 flip made `main` the trunk; the mod was written pre-flip (v0.12.1) and never refit. The 0202 Commander overrode `--base main` on all five member PRs. A fresh FO following the mod literally would open PRs against the deprecated `next` branch.

## Proposed approach (ideation to firm)

Refit the mod so the PR base is `main` (the post-flip trunk) — ideally sourced from a workflow/README config key rather than a literal, so a future trunk change doesn't require another mod edit. Reconcile with the sibling reconcile-helper stale-trunk seed (they share the root cause).

## Acceptance criteria (sketch)

**AC-1 (sketch) — the pr-merge mod opens PRs against the post-flip trunk (`main`), not `next`.**
Verified by: the mod text + any merge-hook behavior test asserting the base is `main`/config-sourced; if behavior-tested, the expected base derives from config, not the mod prose.

## Notes
Mod files are refit-or-worker territory, not FO direct edits. Sibling to the dispatch-reconcile stale-trunk seed.
