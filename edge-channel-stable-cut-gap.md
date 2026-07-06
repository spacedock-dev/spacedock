---
title: "Edge channel survives the stable-cut window (no binary/skills minor skew between a stable tag and the first prerelease)"
status: backlog
source: "Captain report, 0250 Commander session 2026-07-07: edge tap installs the 0.24.0 binary while the spacedock-edge marketplace serves next-branch skills stamped 0.25.0-pre1. Verified: spacedock@next cask version 0.24.0 (updates only on tag push); origin/next .claude-plugin/plugin.json version 0.25.0-pre1; origin/next shared-core gate line 'require binary minor 0.25'. Result: every edge boot since the 0.24.0 stable cut (Jul 4) aborts at the FO version gate (binary too old, minor 0.24 < required 0.25). Broken by the release flow's own design — 'Advancing the Edge Line' bumps next to the post-release dev pre-version at the stable tag, but the edge BINARY only updates on the next tag push — not by any interim merge. Immediate remediation (separate from this task): push the first prerelease tag to realign."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
id: zr2rbsjsak7xx6tetr3n37hc
---

The release flow guarantees a broken edge channel from every stable cut until the first subsequent prerelease tag: the stable tag bumps next's manifests and contract gate line to the post-release pre-version while the spacedock@next cask keeps serving the stable-tag edge build, so the shipped version gate (same major.minor required) aborts every edge boot in the window. Direction options for ideation: (a) defer next's version bump so it rides the first prerelease tag instead of the stable tag; (b) have the stable tag's pipeline also cut a post-release pre-version edge build + cask bump in the same run, closing the window to minutes; (c) make the version gate tolerate skills exactly one minor ahead when the binary channel is edge — weakest, contract-side complexity. Acceptance sketch: value — after a stable cut, an edge-channel install boots green with zero manual steps (baseline: the 2026-07-04..07 window, every edge boot aborting); mechanism — the chosen pipeline change ships. High-stakes surface (CI and release machinery): detached adversarial audit + the release-flow dry-run treatment per docs/releasing.md.

## Root cause localization (FO forensics, 2026-07-07)

- `.github/workflows/release.yml:322-339`, job `edge-advance`, step "Reconcile the edge line past the stable release" (`if: "!contains(github.ref, '-')"`): computes `spacedock-release dev-preversion X.Y.Z` -> `X.(Y+1).0-pre1` and `stamp-version`s it into `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, AND `skills/first-officer/references/first-officer-shared-core.md` (the `FO_PROSE` var) — rewriting the boot version-gate line to a binary minor that does not exist anywhere yet.
- `release.yml:345-353`: the marketplace calendar bump runs on both paths, so on a stable tag it actively triggers every edge installer's re-pull of the incompatible skills.
- The binary half (goreleaser edge build + `spacedock@next` cask bump) only publishes the TAG's version (0.24.0) and does not move again until the next tag push — no 0.25-line binary exists in the same run.
- Latent interaction, not a single bad line: under the retired contract-integer gate, one-minor-ahead skills booted fine (contract 3 == 3). #468 (`511dae11`) replaced it with minor-version coupling BEFORE the 0.24.0 cut; v0.24.0 (Jul 4) was the first stable tag under the new gate, converting the deliberate stamp-ahead ("never masquerade as the shipped stable") into a channel-wide boot abort. docs/releasing.md's own 0.24.0-pre1 note records the drift-cousin of this class.
