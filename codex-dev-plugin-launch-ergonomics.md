---
title: Codex dev/local-plugin launch ergonomics — a `--plugin-dir`-equivalent for the local-marketplace dance
status: backlog
source: "Captain request 2026-06-30 (0240 Commander session). Testing a local/dev build with Codex is far more cumbersome than Claude. Claude: `spacedock claude --plugin-dir <checkout>` — one flag, loads the local plugin checkout directly, bypasses installed-plugin resolution. Codex has no equivalent: you hand-build a local marketplace (`.agents/plugins/marketplace.json` with `source: local` + a `plugins/spacedock` symlink to the checkout), export `SPACEDOCK_MARKETPLACE_SOURCE`, run `spacedock install --host codex`, AND get the channel-in-the-name right (a plain `go build` edge binary needs the marketplace named `spacedock-edge`; the entry stays `spacedock`), plus the `.codex-plugin/plugin.json` version-masquerade gotcha."
group: tooling
id: 4q01qqyx4g2z3rctts1400av
---

## Problem
Launching Codex against a local/dev plugin checkout is a multi-step, footgun-laden ritual, while Claude is a single flag.

- **Claude (the ergonomic target):** `spacedock claude --plugin-dir <checkout>` loads the local plugin directly, bypassing install resolution (`docs/site/get-started/install.md`; the `claude-live` CI lane uses `spacedock claude --plugin-dir "$GITHUB_WORKSPACE"`).
- **Codex (the pain):** no `--plugin-dir` path. The dev must (1) author a local marketplace dir — `.agents/plugins/marketplace.json` (`source: local`, `path: ./plugins/spacedock`) + a `plugins/spacedock` symlink to the checkout (the exact shape `internal/ensigncycle/codex_marketplace.go::writeCodexLocalMarketplace` already builds for CI); (2) `SPACEDOCK_MARKETPLACE_SOURCE=<dir> spacedock install --host codex`; (3) name the marketplace for the channel — a plain `go build` is `devBranch=next` -> edge -> the marketplace MUST be named `spacedock-edge` while the entry stays `spacedock` (must equal the plugin's `plugin.json` name); wrong name -> silently wrong channel; (4) live with the version-masquerade (a fresh `spacedock@spacedock` can report the prior stable version).

## Desired direction (for ideation to refine — not a committed design)
A one-command Codex dev-launch matching Claude's `--plugin-dir` ergonomics — e.g. `spacedock codex --plugin-dir <checkout>` (or `spacedock install --host codex --plugin-dir <checkout>`) that:
- auto-generates the local marketplace from the checkout (promote/reuse `writeCodexLocalMarketplace` into a user-facing helper, not just the CI harness),
- installs from it and launches Codex against the current checkout,
- picks the correct channel-name automatically from the binary's `devBranch` (no `spacedock-edge` footgun),
- surfaces the real version (or at least does not silently masquerade).

The CI building block already exists; the gap is a user-facing front door.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- A single documented command launches Codex against a local checkout's plugin with NO manual marketplace authoring and NO channel-name footgun.
- Verified behaviorally (not prose): a live or fixture-backed Codex install smoke that loads the local checkout — the `github.com` / `ref next` absence guard passes (the same guard the `codex-live` lane uses) and `codex plugin list` resolves the local path.

## Related
- Ergonomic target to match: Claude's `--plugin-dir`.
- Reusable building block: `internal/ensigncycle/codex_marketplace.go::writeCodexLocalMarketplace`.
- Adjacent backlog: `next-post-release-preversion-bump` (the version-masquerade half).
- Docs to update if shipped: `docs/site/get-started/install.md`, `docs/runtime-live-ci.md`.
