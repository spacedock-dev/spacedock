---
title: Codex edge channel install — marketplace-name channels, programmatic install, source override
status: implementation
sprint: 0230-stable-finalization
score: ""
source: captain request after local Codex marketplace setup; extended after the v0.23.0-pre.2 edge dogfood exposed the Codex edge install is broken end-to-end
priority: high
id: z2tjv3570ahjxewv1c309rbc
started: 2026-06-21T06:05:14Z
worktree: .worktrees/spacedock-ensign-spacedock-marketplace-source-env
mod-block:
pr: "#424"
---

## Problem

The Codex edge channel does not install end-to-end — it was unit-tested but never run against a real `codex plugin add`. Dogfooding the v0.23.0-pre.2 edge build surfaced a cluster of defects, on top of the original dev-ergonomics gap (no marketplace-source override).

### Codex rejects the edge entry (the blocker)
Codex enforces that a marketplace entry's name equals the plugin's own `plugin.json` `name`. The channel is encoded in the ENTRY name (`channelEntry` → `spacedock` stable / `spacedock-edge` edge), but the product repo ships a manifest named `spacedock` on every ref. So `codex plugin add spacedock-edge@spacedock` always fails:

```
Error: plugin.json name `spacedock` does not match marketplace plugin name `spacedock-edge`
```

This fails whether the binary runs the command or the user pastes it by hand.

### `install --host codex` is prose-only on a fresh box
`runInit` (init.go) runs the install programmatically ONLY when a Codex plugin already resolves; with no plugin installed it prints the commands (`printCodexInstallProse`) and runs nothing — so it neither installs nor ensures the right marketplace. The claude arm always installs programmatically; Codex is inconsistent.

### The resolver hardcodes the stable id
`resolveCodexManifest` (host_exec.go) confirms the install by checking only for `spacedock@spacedock`, never the channel's id (`spacedock-edge@spacedock`). An edge binary therefore cannot recognize or refresh an installed edge plugin, even once the entry installs.

### No marketplace-source override (original scope)
`spacedock codex` and `spacedock install --host codex` install from the hardcoded `spacedock-dev/marketplace`. A developer can set up a local/alternate marketplace, but the wrapper has no override for the source on the auto-install/install path — awkward for dogfooding, which is the natural way to validate the channel fix below before it reaches the production marketplace.

## Approach

**Encode the channel in the marketplace NAME, not the entry name.** Edge resolves as `spacedock@spacedock-edge`: a marketplace named `spacedock-edge` whose single entry is `spacedock` (tracking next), so the entry name equals the manifest `spacedock` and Codex's name-match passes. The plugin keeps one identity (`/spacedock:` commands), and the scheme is host-agnostic (claude installs the same way). Stable stays `spacedock@spacedock`. This requires the standalone marketplace repo (see `marketplace-repo-decouple`, id `w6bhzvezybbrarkk56zemndd`) to expose a `spacedock-edge` marketplace name — a second `marketplace.json` source, not one marketplace carrying two entries.

With that shape in place: make the Codex fresh-install path programmatic and marketplace-ensuring (unify with the refresh path / the claude arm), make the resolver channel-aware, and add the `SPACEDOCK_MARKETPLACE_SOURCE` override (which then lets a developer point the install at a corrected local marketplace to dogfood the channel fix).

## Host parity

The binary installs to both claude and codex from the same `marketplaceSource`, but the defects split unevenly:

- **Name-match blocker** — Codex confirmed rejects the `spacedock-edge` entry (`plugin.json name 'spacedock' does not match …`). claude's enforcement is unverified (verifying it needs a real `claude plugin install`, which mutates the global claude plugin state). The marketplace-name fix removes the mismatch for BOTH hosts by construction, so it covers claude whether or not claude currently enforces the check.
- **Fresh-install prose-only** — Codex-only. The claude arm of `runInit` always installs programmatically; only the codex arm falls back to prose.
- **Stale-id resolver** — BOTH hosts. `resolveClaudeManifest` (host_exec.go:58) and `resolveCodexManifest` (host_exec.go:81) each hardcode `spacedock@spacedock`; neither resolves the channel's id.
- **Source override** — BOTH hosts. They share the hardcoded `marketplaceSource`; `SPACEDOCK_MARKETPLACE_SOURCE` must apply to both install paths.
- **Marketplace registration health** — claude's `spacedock` marketplace already points at the standalone `spacedock-dev/marketplace`; codex's was stale on the plugin repo (now re-pinned during diagnosis). The programmatic re-pin (AC-3) must be robust on codex regardless.

## Acceptance criteria

- **AC-1 (end-to-end value, REQUIRES one refinement)** — On a clean box the edge channel installs end-to-end with NO entry-name vs plugin.json name mismatch and the launcher then resolves/launches the edge plugin. The REQUIRED, blocking evidence is a real `codex plugin add` of the EDGE channel succeeding where it fails today: extend the existing live lane (runtime-live-e2e.yml:400-402 already runs `codex plugin marketplace add` + `codex plugin add` + `codex plugin list`, today only the STABLE `spacedock@spacedock` entry from a synthesized local marketplace) to exercise the edge entry-name path that currently fails. Codex failing today is the independent baseline that moves the wrong way; the fix must flip it to success. A unit/fixture assertion of the plugin id alone does NOT satisfy this AC.
- **AC-1b (claude half, refined out of AC-1's hard CI gate)** — The same marketplace-name fix makes the claude edge install name-match-safe by construction; claude must NOT regress. Claude enforcement is currently unverified and a real `claude plugin install` mutates GLOBAL claude plugin state, so it is confirmed on a throwaway checkout / out-of-band (recorded in the task), NOT wired as a blocking CI hard-gate. Claude coverage in CI is the channel-aware resolver test (AC-4) plus the construction guarantee, not a live global-state-mutating install.
- **AC-2 (channel naming)** — `channelEntry`/`channelPluginID` and both hosts' install sequences emit a name-match-safe id where the entry name equals manifest `spacedock`: stable `spacedock@spacedock` and edge `spacedock@spacedock-edge` both resolve on claude and codex. channel_selection_test.go's current expectations (`spacedock-edge@spacedock`) are updated to the new shape. The standalone marketplace repo exposes the `spacedock-edge` marketplace name.
- **AC-3 (programmatic fresh install - Codex only)** — `spacedock install --host codex` with no plugin installed RUNS the marketplace add + plugin add and ensures the marketplace source is (re-)pinned, instead of printing prose (init.go:67-72). A real re-pin failure surfaces rather than being masked by the tolerated `marketplace remove`. The claude arm already installs programmatically; no change there.
- **AC-4 (channel-aware resolver - both hosts)** — resolveCodexManifest AND resolveClaudeManifest (plus codexCacheManifest's cache path) resolve the channel's id, so an edge binary recognizes/refreshes an installed edge plugin on either host, not the hardcoded stable id.
- **AC-5 (source override)** — `spacedock codex`, the front-door auto-install, and `spacedock install --host codex` honor `SPACEDOCK_MARKETPLACE_SOURCE`, preserving the default `spacedock-dev/marketplace` when unset. Tests observe the actual source passed to the install seam (the fakeHost seam at frontdoor_test.go:42 records {host, source, devBranch}) for unset/default + override cases.
- **AC-6 (docs)** — Codex / front-door development guidance (docs/site/get-started/install.md and/or docs/site/contributing/agent-development.md) documents the channel scheme, the env override + local-marketplace dogfood use case, and the `--plugin-dir` caveat (bypasses installed-plugin resolution, does not solve launcher safehouse wrapping). Per docs/dev/README.md:108 the concrete before/after doc diff is recorded in the task body at ideation (see "AC-6 doc diff" below), applied at implementation.
- **AC-7 (validation)** — `go test ./internal/cli` green; a test exercises the Codex entry-name vs plugin.json name-match constraint (or a faithful fixture) so this integration gap cannot regress silently. This is the mechanism-shipped backstop that serves AC-1; it does NOT substitute for AC-1's real `codex plugin add`.

## Implementation note (channel_selection_test.go expectation flip)

`channel_selection_test.go` currently LOCKS the broken shape `spacedock-edge@spacedock` (~lines 49/85/127/193), GREEN against today's broken `channelPluginID` (host_exec.go:230-232). After the fix flips the id to `spacedock@spacedock-edge`, those expectations must update to the new shape as PART of AC-2 — it is NOT a regression. A Commander should not mistake the post-fix failing test for a break.

## AC-6 doc diff (before/after for docs/site/get-started/install.md)

Today install.md "## Skills" (lines 45-48) shows only the claude stable install:

```bash
claude plugin marketplace add spacedock-dev/marketplace
claude plugin install spacedock@spacedock
```

After (document the channel scheme + the source override; entry name equals manifest `spacedock` on both channels so the host name-match passes):

```bash
# Stable (default channel) — marketplace named `spacedock`, entry `spacedock`
claude plugin marketplace add spacedock-dev/marketplace
claude plugin install spacedock@spacedock

# Edge (tracks next) — marketplace named `spacedock-edge`, entry still `spacedock`
claude plugin install spacedock@spacedock-edge
```

Plus a short note that `SPACEDOCK_MARKETPLACE_SOURCE` overrides the default `spacedock-dev/marketplace` install source (for dogfooding a local/alternate marketplace), and the `--plugin-dir` caveat (it bypasses installed-plugin resolution but does not solve launcher safehouse wrapping). Final exact wording is firmed at implementation; this is the gate-reviewed before/after the ideation gate approves.

## Notes

- Root cause of the integration gap: `channel_selection_test.go` asserts `spacedock-edge@spacedock` but never ran a real `codex plugin add`, so Codex's name-match was never exercised — the "validate the riskiest path end-to-end" step that was skipped.
- Related: `marketplace-repo-decouple` (`w6bhzvezybbrarkk56zemndd`) — the standalone marketplace repo this builds on (done/REJECTED-superseded by `marketplace-repo-and-pinned-channels`, PR #352; the `spacedock-edge` marketplace-name prerequisite is SHIPPED, not a dependency).

## Stage Report: implementation

- DONE: A real `codex plugin add` of the edge channel succeeds on the live lane — the failing-today baseline is flipped to success (AC-1); record the before (codex rejects the spacedock-edge entry) and after (installs as spacedock@spacedock-edge from a distinct spacedock-edge marketplace name)
  Validated locally on codex-cli 0.141.0 (isolated CODEX_HOME): BEFORE `codex plugin add spacedock-edge@spacedock` → `Error: plugin.json name 'spacedock' does not match marketplace plugin name 'spacedock-edge'`; AFTER `codex plugin add spacedock@spacedock-edge` → `Added plugin 'spacedock' from marketplace 'spacedock-edge'`, installed at cache/spacedock-edge/spacedock/0.22.0, listed `installed, enabled`. Wired into runtime-live-e2e.yml as a baseline-fails → fix-succeeds step (commit 7be7e2fc).
- DONE: channel_selection_test.go expectation is flipped from spacedock-edge@spacedock to spacedock@spacedock-edge, both resolvers are channel-aware, and go test ./... passes
  Inverted the model (commit 7241aa67): channelEntry is always `spacedock`; new channelMarketplace(devBranch) carries the channel (main→spacedock, edge→spacedock-edge); channelPluginID = entry@marketplace. resolveClaudeManifest, resolveCodexManifest + codexCacheManifest are channel-aware via the package devBranch. Flipped channel_selection_test.go and the sibling tests that locked the old shape. `go test ./...` green (all ok).
- DONE: SPACEDOCK_MARKETPLACE_SOURCE override is wired and the codex fresh-install path is programmatic (not prose-only on a fresh box)
  AC-5 (commit 4fdcbddf): marketplaceSource is now a var; applyMarketplaceSourceOverride mirrors applyDevBranchOverride, applied in the claude/codex/install cobra commands; tests observe the source at the install seam for default + override cases. AC-3 (commit 7241aa67): runInit's codex arm drives the install seam (marketplace add + plugin add, re-pin) on a fresh box too; removed the dead printCodexInstallProse.

### Summary
Fixed the Codex edge install end-to-end by encoding the channel in the marketplace NAME (`spacedock-edge`) rather than the entry name, so the entry stays `spacedock` (== plugin.json name) and Codex's name-match passes. Made both resolvers channel-aware, made the codex fresh-install programmatic, added the `SPACEDOCK_MARKETPLACE_SOURCE` override, and updated docs. The riskiest path (AC-1) was proven live on codex-cli 0.141.0: the old `spacedock-edge@spacedock` shape fails with the exact blocker error, the new `spacedock@spacedock-edge` shape installs and lands at the resolver's cache path. AC-7 ships a hermetic real-codex backstop (skips when codex absent). AC-6 docs the scheme + override + --plugin-dir caveat. AC-1b (claude half) is covered by the AC-4 resolver test + the construction guarantee, not a global-state-mutating live install. `go test ./...` is green.

## Stage Report: validation

- DONE: Reproduce AC-1 independently on an isolated CODEX_HOME (or via the hermetic real-codex backstop if codex absent): `codex plugin add spacedock@spacedock-edge` installs (entry `spacedock` from marketplace `spacedock-edge`), and the old `spacedock-edge@spacedock` shape fails with the plugin.json-name-vs-marketplace-name mismatch error. Confirm the failing-today→success flip is real, not asserted.
  codex-cli 0.141.0 present. Reproduced by hand on a throwaway isolated CODEX_HOME (NOT the Go harness): OLD `spacedock-edge@spacedock` → exit 1, `plugin.json name 'spacedock' does not match marketplace plugin name 'spacedock-edge'`; NEW `spacedock@spacedock-edge` → exit 0, `Added plugin 'spacedock' from marketplace 'spacedock-edge'`, listed `installed, enabled`, cache at `cache/spacedock-edge/spacedock/0.0.0/.codex-plugin/plugin.json`. Also exercised the CI step's exact `.agents/plugins`+`source:local` layout on 0.141.0 — NEW shape installs/lists. AC-7 backstop (TestCodexEntryNameMustMatchPluginName) ran (not skipped), PASS 0.19s.
- DONE: Confirm the channel model and tests: channelEntry is always `spacedock`, channelMarketplace carries the channel (main→spacedock, edge→spacedock-edge), both resolveClaudeManifest/resolveCodexManifest are channel-aware, channel_selection_test.go + siblings flipped to the new shape, and `go test ./...` is green.
  Verified in host_exec.go (channelEntry/channelMarketplace/channelPluginID + both resolvers + codexCacheManifest, all keyed on devBranch). channel_selection_test.go and siblings assert `spacedock@spacedock-edge` and guard against the old `spacedock-edge@spacedock`; resolver/seam tests use independent constructed values (not substring-over-file). `go test ./...` all packages ok (cli 37.7s).
- DONE: Confirm SPACEDOCK_MARKETPLACE_SOURCE override + the programmatic codex fresh-install path are wired (printCodexInstallProse removed; install seam drives marketplace add + plugin add + re-pin), and the runtime-live-e2e.yml baseline-fails→fix-succeeds step exists.
  cli.go applies applyMarketplaceSourceOverride on claude/codex/install; init.go codex arm drives ops.Install on a fresh box (printCodexInstallProse gone); seam-observing tests cover default+override and install+front-door. runtime-live-e2e.yml:412+ has the "Codex edge channel install (name-match baseline → fix)" step with a baseline-regression guard.

### Summary
PASSED. All seven ACs verified with reproduced evidence; `go test ./...` green. AC-1's failing-today→success flip is REAL, reproduced by hand on an isolated CODEX_HOME (codex 0.141.0) outside the Go harness, and the live-CI step's exact marketplace layout was also exercised. Tests are behavioral (real codex CLI / seam observation / independent constructed values), not substring-over-instruction-file. Fixed one stale doc comment in channel_selection_test.go that named the old broken edge id (commit d5981619) — comment-only, behavior was already correct. Non-blocking notes: a pre-existing `go vet` finding (pi_frontdoor_test.go:701, also on main, out of scope); the pre-existing claude decoupling test still installs `spacedock-edge@spacedock` and passes, confirming claude does NOT enforce the entry-name match (consistent with AC-1b's "claude unverified" stance — the fix is name-match-safe by construction regardless).
