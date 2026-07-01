---
title: Release process must advance the `next` edge line (reconcile + calendar bump) on every release, not just `stable`
status: ideation
source: "0.24.0-pre1 prerelease cut, 2026-07-01 (Commander). Cutting v0.24.0-pre1 advanced the binary + stamped main to 0.24.0-pre1, but left the `next` edge line stale — it had diverged 40 commits during the 0240 sprint, and release.yml's hyphenated-tag carve-out advances NEITHER `stable` (correctly skipped) NOR `next`. So `spacedock install --host codex` (the spacedock-edge marketplace serves the plugin from `next`) kept serving the old 0.23.0-pre plugin, and the binary(0.24.0-pre1)/plugin(0.23.0-pre) version-compat check hard-blocked `spacedock codex`. Manually reconciled next -> main@0.24.0-pre1 (merge favoring main) + bumped the marketplace calendar key (origin/next now 1bb3da06), which unblocked the edge install. The release process should do this automatically."
group: tooling
id: s20pdb1pzexwkbp5b4cz30av
sprint: 0240-lean-contract
started: 2026-07-01T02:38:47Z
---

## Problem
`release.yml` advances the `stable` ref only on a stable (non-hyphenated) `vX.Y.Z` tag; on a prerelease (hyphenated `-pre`) tag it advances NEITHER `stable` (correctly, per the carve-out) NOR the `next` edge line. So the edge channel (the `spacedock-edge` marketplace, which serves the plugin from `next`) does not advance with a prerelease cut — it drifts behind `main` until someone manually reconciles it. The 0240 sprint left `next` 40 commits behind main; the 0.24.0-pre1 cut then produced a binary(0.24.0-pre1)/edge-plugin(0.23.0-pre) skew that the strict version-compat check hard-blocked `spacedock codex`.

## Desired direction (ideation to refine)
The release process keeps the `next` edge line advanced on every release:
- On a PRERELEASE (`-pre`) tag: release.yml's existing hyphenated-tag carve-out ALSO advances `next` to the tagged commit (reconcile the edge line to the release SHA) AND bumps the marketplace calendar key (`spacedock-release bump-calendar .claude-plugin/marketplace.json`), so the edge marketplace serves the prerelease and codex/claude re-pull.
- On a STABLE (`vX.Y.Z`) tag: advance `next` to the post-release dev pre-version (couples with `next-post-release-preversion-bump`) so the edge line does not masquerade as the just-released stable.
- Document the edge-line advance in `docs/releasing.md` alongside the `stable`-advance step.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- Cutting a prerelease tag leaves `origin/next` reconciled to the tagged commit's content, with the plugin version == the tag and a bumped calendar key — verified by a release-machinery test/fixture (or a `spacedock-release` gate), NOT a manual step.
- The edge marketplace serves the prerelease version after the cut (a codex/fixture install smoke resolving the edge/local marketplace to the prerelease version).
- `docs/releasing.md` documents the edge-line advance for both prerelease and stable cuts.

## Related
- `next-post-release-preversion-bump` — the version-stamp masquerade half of the dev/edge line advance.
- `minor-version-compat-coupling` — why the skew hard-blocks; a laxer contract-compatible check softens the failure mode.
- The 0.24.0-pre1 cut (this session) that exposed it; the manual reconcile (`origin/next` @ 1bb3da06).
- `release.yml` prerelease carve-out (`if: "!contains(github.ref, '-')"`), `next-publish.yml` (calendar bump).
