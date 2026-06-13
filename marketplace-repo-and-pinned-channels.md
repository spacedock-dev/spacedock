---
id: gpvg343kshg02e34kq7hrk9f
title: Separate marketplace repo + tag-pinned stable / edge channels (Model B core)
status: ideation
source: "Post-flip release-model decision (roadmap 0201, captain 2026-06-09): adopt Model B — stable channel serves a pinned release tag, edge serves next HEAD — grounded by the plugin-distribution research (wcdgsgd88). This is the structural core."
started: 2026-06-13T04:08:48Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness:
---

Move the marketplace manifest out of the plugin branches into its own repo, and serve two pinned channels: **stable** = a release tag, **edge** = `next` HEAD. This is the structural core of the post-flip release model (Model B). See `docs/roadmap/0201-post-flip-release-model/index.md`.

## Problem

Today `marketplace.json` lives *inside* the plugin branch, carrying a per-branch `source.ref` (`main` for stable, `next` for edge). Two consequences:
1. **`main` and `next` permanently differ on `source.ref`**, so `next → main` is not a clean fast-forward — every release has to re-settle the ref. Confirmed: `origin/main` marketplace `ref: main`, `origin/next` `ref: next`.
2. **The plugin is served from a branch HEAD** — per the Claude Code docs, a git source with no effective version pin means *"every new commit is treated as a new version"* (https://code.claude.com/docs/en/plugins-reference). So a fresh stable install pulls `main` HEAD regardless of any release boundary. No stable end-user channel in any mature ecosystem does this (npm/VS Code/Homebrew/Chrome all serve a pinned artifact; Anthropic's own curated `claude-plugins-official` sha-pins 125/125 entries).

## Proposed approach (ideation firms)

- Stand up a standalone marketplace repo holding `marketplace.json`; plugin branches carry **no** manifest → the main/next `source.ref` divergence disappears and `next → main` becomes a clean fast-forward.
- Two entries (the docs' documented channel convention — two refs of the same source, each resolving to a **distinct** version, https://code.claude.com/docs/en/plugin-marketplaces): `spacedock` (stable) pinned `ref: v0.X.Y`; `spacedock-edge` on `ref: next` (HEAD).
- Point `spacedock install` at the new marketplace; map the two-channel `devBranch` stamp (from `k6d`) onto the stable/edge entries so a stable binary resolves the stable entry and an edge binary the edge entry.
- **Seamless migration:** the next binary release carries the new marketplace source; `install` / the front door auto-repoints existing users (remove old marketplace + add new — the sequence `install` already runs). This **depends on** the install-refresh fix (`tes` / install-refresh-and-upgrade-hint) actually re-pulling reliably.

## Riskiest mechanism — exercise first (per the spike rule)

Before committing to the rest: clean-install from a **tag-pinned** stable entry, make a **no-op commit on `main`**, run a plugin update, and confirm **stable does NOT update while edge (next HEAD) does**. That single end-to-end exercise proves the whole decoupling — pinning a tag actually freezes the served bytes, and the two channels resolve distinct versions (the docs' "each channel must resolve to a different version or the update is skipped" caveat). Record the result in the task body.

## Out of scope

- The **release ritual ordering** (stamp-then-tag) — separate task `stamp-then-tag-release-ritual`.
- The **install-refresh correctness + upgrade hint** — existing task `tes` (this task depends on its refresh fix for migration, but does not own it).
- Contract-compatibility semantics — unchanged.

## Acceptance criteria

(Ideation firms. Each verified outside this task body — git state, on-disk installed manifest, or command output; never a prose-grep.)

**AC-1 (sketch) — a tag-pinned stable channel is decoupled from `main` HEAD.** Verified by: the decoupling exercise above — after installing from the stable entry and pushing a no-op `main` commit, the installed plugin manifest on the stable channel is unchanged while the edge channel advances — observed in on-disk installed-plugin state.

**AC-2 (sketch) — plugin branches no longer carry a marketplace manifest, and `next → main` is a clean fast-forward.** Verified by: absence of `marketplace.json` on the plugin branches + a fast-forward `next → main` in a controlled check — git state.

**AC-3 (sketch) — `spacedock install` resolves the intended channel from the new marketplace repo on both hosts.** Verified by: an install smoke against the new marketplace where a stable binary lands the stable entry and an edge binary the edge entry — resolved/installed manifest source, not the command's own claim alone.

## Test plan

(Ideation/implementation firms.) The decoupling spike first (cheapest end-to-end proof). Then: a git-level check for branch-manifest absence + clean ff; an install smoke per host resolving the right channel. Live-host install smoke only where the host resolution itself is the claim.
