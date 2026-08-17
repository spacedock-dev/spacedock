---
id: 2dfpswjwbxez1km6439wfrsn
title: Collapse the duplicate marketplace routes so the channel lives in the entry name
status: ideation
source: "Captain directive, CL 2026-08-17: \"we should simplify the marketplace, there should just be one way\", then the chosen shape: \"i think spacedock-edge@marketplace is better, because we now have other things like subspace in the same marketplace.\" Raised from the live cross-host channel investigation; verified against the published marketplace on both refs and against a route-A install cached on the captain's workstation."
started: 2026-08-17T19:38:45Z
completed:
verdict:
score:
worktree:
issue:
---

The edge plugin is reachable by two ids resolving identical content. Collapse to one, with the channel carried by the marketplace ENTRY name and a single shared marketplace hosting the whole plugin family.

## Problem

`spacedock-dev/marketplace` publishes a `marketplace.json` on two refs, and both expose an edge route:

| ref | marketplace `name` | entries |
| --- | --- | --- |
| default | `spacedock` | `spacedock` (→ spacedock.git ref `stable`), `spacedock-edge` (→ ref `next`), `subspace`, `cargento` |
| `edge` | `spacedock-edge` | `spacedock` (→ ref `next`) |

`spacedock-edge@spacedock` and `spacedock@spacedock-edge` deliver identical content from `spacedock.git` ref `next`. Verified live: the captain's workstation holds installs made both ways (`cache/spacedock/spacedock-edge/` and `cache/spacedock-edge/spacedock/`), and the reporting Linux host a third combination.

Two costs. Operators have no single correct answer to "how do I install edge", and which route they picked determines whether a later command recognizes the install. The binary can only express `spacedock@<marketplace>` because `channelEntry` returns a hardcoded `"spacedock"` (internal/cli/host_exec.go:223), so the route-A id is one it cannot construct or recognize.

The decisive argument against encoding the channel in the MARKETPLACE is catalog scope: the marketplace now carries `subspace` and `cargento`, while the edge-branch marketplace carries only `spacedock`. A per-channel marketplace either duplicates the whole catalog per channel or strands every non-spacedock plugin for edge users, and lets the edge catalog drift from stable silently because catalog content is versioned by git branch.

## Proposed approach

Captain-directed shape: retire the `edge`-branch marketplace. Keep ONE marketplace (`spacedock`) hosting the family, and carry the channel in the entry name — `spacedock` for stable, `spacedock-edge` for edge. In the binary this inverts the current pair: `channelMarketplace` becomes the constant, `channelEntry` becomes channel-varying.

{Ideation refines migration and confirms the risk below.}

## Out of scope

The prerelease asset-naming defect and the artifact-level channel guard — separate entity. Changing how `devBranch` selects a channel.

## Riskiest unverified mechanism — spike before designing

`internal/cli/host_exec.go:218-222` asserts the host "rejects a marketplace entry whose name differs from the manifest name (codex confirmed, claude by construction)". The claude half is already FALSIFIED: the cached route-A install at `~/.claude/plugins/cache/spacedock/spacedock-edge/0.19.9/.claude-plugin/plugin.json` declares `"name": "spacedock"` under an entry named `spacedock-edge`, and it installed and ran.

Unresolved, and gating: (a) does codex enforce the name match, and (b) does the skill/agent NAMESPACE derive from the manifest name or the entry name? If the entry name leaks into namespacing, this shape renames every skill to `spacedock-edge:*` and breaks all references — which would sink the approach. Exercise both against a real host before designing the migration; record the result in this body.

## Expected surface and tolerance

Estimate net LOC change: +40, across 3 files, plus marketplace-repo changes outside this repository. Semantics changed: the published plugin id for the edge channel, and the marketplace source the binary adds.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - Each channel resolves through exactly one published plugin id.**
Verified by: fetching `marketplace.json` from every published ref and asserting no two (ref, entry) pairs resolve to the same `spacedock.git` ref. Fails while both `spacedock-edge@spacedock` and `spacedock@spacedock-edge` point at `next` — the current state.

**AC-2 - The id the binary constructs is the id the marketplace publishes, on every channel and every supported host.**
Verified by: a test asserting `channelPluginID` output for each devBranch value exists as a real (marketplace, entry) pair in the published marketplace data, plus an install exercised on both claude and codex. Fails if the surviving route is one a supported host cannot resolve.

**AC-3 - The skill and agent namespace is unchanged by the migration.**
Verified by: installing the surviving edge id on a real host and asserting the first-officer skill still resolves as `spacedock:first-officer`. Fails if the entry rename leaks into namespacing — the outcome that would invalidate this approach rather than merely complicate it.

**AC-4 - An operator already installed via the retired route is not silently left broken.**
Verified by: the documented migration path exercised against a host holding the retired id, showing it either continues to resolve or produces an actionable message naming the replacement. Fails if a retired route degrades to plugin-not-found with no remedy text.
