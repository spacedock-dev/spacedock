---
title: Codex edge channel install — marketplace-name channels, programmatic install, source override
status: backlog
score: ""
source: captain request after local Codex marketplace setup; extended after the v0.23.0-pre.2 edge dogfood exposed the Codex edge install is broken end-to-end
priority: high
id: z2tjv3570ahjxewv1c309rbc
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

## Acceptance criteria

- **AC-1 (end-to-end value)** — On a clean box, installing the edge channel on Codex SUCCEEDS against a real `codex plugin add`, with no `plugin.json name … does not match` error, and `spacedock codex` then resolves and launches the edge plugin. The same exercise fails today — that is the independent baseline that moves the wrong way, and the fix must flip it to success. A unit assertion of the plugin id alone does NOT satisfy this AC.
- **AC-2 (channel naming)** — `channelEntry` / `channelPluginID` and the Codex + claude install sequences emit a name-match-safe id (entry name equals the manifest `spacedock`); edge resolves via the `spacedock-edge` marketplace name, which the standalone marketplace repo exposes. Stable (`spacedock@spacedock`) and the claude arm are unaffected (verified).
- **AC-3 (programmatic fresh install)** — `spacedock install --host codex` with no plugin installed RUNS the marketplace add + plugin add and ensures the marketplace source is (re-)pinned, rather than printing prose. The tolerated `marketplace remove` no longer masks a failed re-pin — a real re-pin failure surfaces.
- **AC-4 (channel-aware resolver)** — `resolveCodexManifest` resolves the channel's id (an edge binary recognizes/refreshes an installed `spacedock@spacedock-edge`), not the hardcoded stable id.
- **AC-5 (source override)** — `spacedock codex` and `spacedock install --host codex` honor `SPACEDOCK_MARKETPLACE_SOURCE` on the install/auto-install path, preserving the default `spacedock-dev/marketplace` when unset. Tests observe the actual source passed to the install seam (unset/default + override cases).
- **AC-6 (docs)** — Codex / front-door development guidance documents the channel scheme, the env override + local-marketplace use case, and the `--plugin-dir` caveat (it bypasses installed-plugin resolution but does not solve launcher safehouse wrapping).
- **AC-7 (validation)** — `go test ./internal/cli` green; a test exercises the Codex entry-name vs `plugin.json` name-match constraint (or a faithful fixture) so this integration gap cannot regress silently.

## Notes

- Root cause of the integration gap: `channel_selection_test.go` asserts `spacedock-edge@spacedock` but never ran a real `codex plugin add`, so Codex's name-match was never exercised — the "validate the riskiest path end-to-end" step that was skipped.
- Related: `marketplace-repo-decouple` (`w6bhzvezybbrarkk56zemndd`) — the standalone marketplace repo this builds on.
