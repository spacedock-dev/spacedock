---
id: fea266y405b95053aq86q5d8
title: Codex edge install — per-channel marketplace source (edge branch + channel-specific binary source)
status: implementation
source: handoff-codex-edge-0230 (z2 spacedock-marketplace-source-env follow-up)
sprint: 0230-stable-finalization
started: 2026-06-22T00:36:49Z
completed:
verdict:
score: 0.7
worktree: .worktrees/spacedock-ensign-codex-edge-channel-marketplace-source
issue:
---

The codex edge channel cannot install. The edge binary builds plugin id `spacedock@spacedock-edge`, but `codex plugin marketplace add spacedock-dev/marketplace` registers a marketplace NAMED `spacedock` (the repo's root `marketplace.json` `name`, which carries `spacedock-edge` only as an *entry*), so `codex plugin add spacedock@spacedock-edge` fails: `plugin spacedock was not found in marketplace spacedock-edge`. z2 shipped the binary half but validated its marketplace half only against a synthesized live-lane fixture — the real repo + real source were never wired, which masked this. Blocks pre.3 and v0.23.0.

## Problem

The binary's `channelMarketplace(edge)` resolves marketplace name `spacedock-edge`, but `internal/cli/init.go:22 marketplaceSource = "spacedock-dev/marketplace"` is a single var used for BOTH channels. The bare source's root `marketplace.json` registers under the name `spacedock`, so an edge install adds a marketplace named `spacedock`, then looks up `spacedock@spacedock-edge` and fails. Stable is unaffected (entry `spacedock` in marketplace `spacedock` → `spacedock@spacedock` matches).

## Proposed approach (captain pre-approved — handoff)

1. **edge branch in `spacedock-dev/marketplace`** (NOT a new repo). Its root `.claude-plugin/marketplace.json` is NAMED `spacedock-edge`, with a single entry NAMED `spacedock` whose `source` points at the product repo `spacedock.git` ref `next` (model the entry source on the current main `marketplace.json`'s `spacedock-edge` entry: url `https://github.com/spacedock-dev/spacedock.git`, ref `next`).
2. **VERIFY the mechanism FIRST (riskiest path — spike before any binary change):** `codex plugin marketplace add spacedock-dev/marketplace@edge` must register a DISTINCT marketplace named `spacedock-edge` (alongside the existing `spacedock`), and `codex plugin add spacedock@spacedock-edge` must SUCCEED (entry `spacedock` == the plugin's own `.codex-plugin/plugin.json` name → name-match passes). Record the transcript. If `owner/repo@ref` does NOT yield a distinct-named marketplace, STOP and report to the captain — the fallback is a separate `spacedock-dev/marketplace-edge` bare-source repo, which is the captain's call; do NOT create a new repo unilaterally.
3. **Binary: channel-specific marketplace source.** Make the edge install source `spacedock-dev/marketplace@edge` while stable stays `spacedock-dev/marketplace`; keep the `SPACEDOCK_MARKETPLACE_SOURCE` override working for both. TDD.
4. **Validate against the REAL marketplace** (an actual `codex plugin add` / `spacedock install --host codex` from the published edge), NOT a synthesized fixture — the fixture is exactly what masked this.

## Out of scope

- Cutting pre.3 / v0.23.0 (separate release step, gated on this fix).
- The Claude channel (already installs correctly).
- Re-pinning the edge entry version policy (tracked elsewhere); the edge entry tracks `next` HEAD.

## Acceptance criteria

**AC-1 — A real codex edge install succeeds end-to-end (end-value, measured against the broken baseline).**
Verified by: an actual codex command transcript captured in the validation stage report — `codex plugin marketplace add spacedock-dev/marketplace@edge` registers a marketplace named `spacedock-edge` distinct from `spacedock`, and `codex plugin add spacedock@spacedock-edge` exits 0 with `codex plugin list` showing `spacedock@spacedock-edge` installed. Baseline that moves the wrong way: today this errors `plugin spacedock was not found in marketplace spacedock-edge`.

**AC-2 — The binary resolves a channel-specific marketplace source.**
Verified by: a Go test asserting the edge devBranch resolves the install source `spacedock-dev/marketplace@edge`, stable resolves `spacedock-dev/marketplace`, and `SPACEDOCK_MARKETPLACE_SOURCE` overrides both. Code gate, not prose.

**AC-3 — `spacedock install --host codex` drives a real, succeeding edge install through the install seam.**
Verified by: a real `spacedock install --host codex` run against an edge binary, validated against the published edge marketplace (not a synthesized fixture), exit 0 with the edge plugin installed/refreshed.

## Test plan

- Spike first (step 2): run the manual codex commands against the new edge branch; record the transcript in the task body. Seeds the implementation's first test.
- Go unit test: channel-specific source resolution for edge vs stable + `SPACEDOCK_MARKETPLACE_SOURCE` override (AC-2).
- Real validation: an actual `codex plugin add spacedock@spacedock-edge` and `spacedock install --host codex` against the published edge (AC-1, AC-3).
- High-stakes surface (install/release machinery): detached adversarial audit at validation per the README proof policy.
