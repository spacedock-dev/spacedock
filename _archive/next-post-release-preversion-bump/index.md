---
title: Post-release `next` preversion bump for dev/plugin manifests
status: done
source: "Captain question 2026-06-21 after Codex local marketplace dogfood: should dev bump to 0.(next).0-pre.1 once stable releases? FO recommendation: yes, as a next-only post-release step."
started:
completed: 2026-08-25T05:35:42Z
verdict: PASSED
score: 0.25
worktree:
issue:
sprint: 0201-post-flip-release-model
group: release-model
related: "next-independent-release-line; stamp-then-tag-release-ritual; spacedock-marketplace-source-env"
id: grhae4mq53h8vv3zy2w6j62m
---

After a stable release `0.N.0`, make the dev/edge line visibly advance to `0.(N+1).0-pre.1` instead of leaving in-tree plugin manifests and local-marketplace installs showing the just-released or older stable version.

## Problem

The stable release flow stamps plugin manifests for the stable tag, but the dev line can remain version-incoherent afterward. During Codex local-marketplace dogfood, a fresh `spacedock@spacedock` install still displayed `0.22.0` because `.codex-plugin/plugin.json` and `.claude-plugin/plugin.json` carried that version even though the code being loaded was the in-tree/dev plugin.

That makes it hard to tell whether Codex is loading the latest in-tree plugin, a stale cache, or the stable channel. It also undermines the intended edge-channel shape: after stable `0.N.0`, development should identify as the next release candidate line, not as the previous stable.

## Proposed approach

Add an explicit post-release, `next`-only version bump ritual:

1. Cut stable `0.N.0` from the release branch/tag using the stable release ritual.
2. Repoint/refresh the stable marketplace entry to that release.
3. Move or refresh `next`.
4. On `next`, bump `.codex-plugin/plugin.json`, `.claude-plugin/plugin.json`, and the binary-reported development version to `0.(N+1).0-pre.1`.
5. Use subsequent prerelease bumps (`pre.2`, `pre.3`, etc.) when the edge/dev line needs a cache-visible refresh.

Before making this policy, prove the host accepts these prerelease version strings: `codex plugin add/list` must handle a manifest version like `0.24.0-pre.1` without rejecting or normalizing it incorrectly. Claude should be checked too if the same manifests are shared by both hosts.

This is related to, but narrower than, `next-independent-release-line`: it does not decide the whole edge cask/publish architecture. It pins the immediate dev-process invariant that the in-tree plugin and local marketplace stop masquerading as an old stable version after a release.

## Out of scope

- Reworking the full edge-channel marketplace architecture; that belongs to `spacedock-marketplace-source-env` and the marketplace-repo work.
- Changing the stable release tag ritual; that belongs to `stamp-then-tag-release-ritual`.
- Deciding whether `next` gets an independent binary/cask publish trigger; that belongs to `next-independent-release-line`.

## Acceptance criteria

**AC-1 - The release runbook includes a post-stable `next` preversion bump.**
Verified by: `docs/releasing.md` or its successor runbook names the sequence after stable `0.N.0`: bump dev/edge manifests and binary dev version on `next` to `0.(N+1).0-pre.1`, with a concrete example such as `0.23.0 -> 0.24.0-pre.1`.

**AC-2 - The version stamp path accepts prerelease versions for all dev-visible version sources.**
Verified by: a test or scripted smoke run stamps a prerelease semver string into `.codex-plugin/plugin.json`, `.claude-plugin/plugin.json`, and the binary/version source, then asserts the resulting JSON/version output contains exactly the prerelease string.

**AC-3 - Codex accepts the prerelease plugin manifest version.**
Verified by: a live or fixture-backed Codex install smoke using a local marketplace whose plugin manifest has a version like `0.24.0-pre.1`; `codex plugin add` succeeds and `codex plugin list --json` reports the prerelease version exactly.

**AC-4 - The local dev marketplace display is cache-visible and not stale-stable-looking after the bump.**
Verified by: after refreshing the local marketplace/install path, the installed `spacedock@spacedock` or edge equivalent reports the new prerelease version rather than the prior stable version.

## Test plan

Start with the riskiest path: a local Codex marketplace fixture or dogfood marketplace pointing at an in-tree plugin stamped to `0.(N+1).0-pre.1`, then run real `codex plugin add/list` if available. Add a small Go/unit or script-level test for the stamp command/version-source update so the prerelease string does not regress. Documentation is required but should not be the only proof.
