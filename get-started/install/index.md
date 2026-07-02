---
title: "Install"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-07-02 04:16:03"
---

# Install Spacedock

Spacedock works with a coding agent you already have: Claude Code, Codex, or Pi. Install one of those first.

```
brew tap spacedock-dev/homebrew-tap
brew install spacedock
```

`brew install` also pulls in `agentsview` (it powers `/spacedock:survey`). The optional sandbox, safehouse, is installed separately — see [Sandboxing](../../reference/sandbox/).

```
curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
```

Installs a checksum-verified binary to `~/.local/bin`.

## Launch

In a project you already have:

```
spacedock claude "/spacedock:survey"
```

Or launch directly:

```
spacedock claude "what can spacedock do for me in this project"
```

Replace `claude` with `codex` or `pi` for the respective coding agents.

## Skills

Spacedock installs the relevant skills on launch. To install them manually:

```
# Stable (default channel) — marketplace named `spacedock`, entry `spacedock`
claude plugin marketplace add spacedock-dev/marketplace
claude plugin install spacedock@spacedock

# Edge (tracks next) — marketplace named `spacedock-edge`, entry still `spacedock`
claude plugin marketplace add spacedock-dev/marketplace@edge
claude plugin install spacedock@spacedock-edge
```

The channel is the marketplace name; the entry name stays `spacedock` on both channels (it equals the plugin's own `name`, so the host's entry-name vs plugin-name check passes). Each channel adds its own marketplace source — the stable marketplace lives at the repo root (named `spacedock`), the edge one on the `@edge` branch (named `spacedock-edge`) — so the `@edge` ref is what registers the `spacedock-edge` marketplace the edge entry resolves from. Codex installs the same way with `codex plugin add`.

Set `SPACEDOCK_MARKETPLACE_SOURCE` to install from a local or alternate marketplace instead of the default `spacedock-dev/marketplace` — useful for dogfooding a marketplace change before it reaches the production marketplace:

```
SPACEDOCK_MARKETPLACE_SOURCE=/path/to/local/marketplace spacedock install --host codex
```

Launching with `--plugin-dir` loads a local plugin checkout directly. On Claude and Pi this is an ephemeral, install-free override — it bypasses installed-plugin resolution for that one launch and does not wrap the launch in the safehouse sandbox; use it for plugin development, not as an install substitute.

Codex has no such flag on its own CLI, so `spacedock codex --plugin-dir <checkout>` and `spacedock install --host codex --plugin-dir <checkout>` build a local marketplace from the checkout and install it under the binary's own channel (`spacedock` stable / `spacedock-edge` edge — matching whatever `spacedock codex` would otherwise install), then launch. This IS a persistent install, replacing whatever Codex plugin was previously configured, and it is a point-in-time snapshot: editing the checkout afterward has no effect until the command is re-run. The command prints an advisory that the reported version reflects the checkout's checked-in manifest, not necessarily its current HEAD.

## Sandboxing

See [supported sandboxes](../../reference/sandbox/).

## Troubleshooting

Run `spacedock doctor`.

## Next

Read about the [survey report](../survey/) to understand your usage pattern with coding agents, or start with [your first workflow](../first-workflow/).

## Sitemap

- [Survey your project](../survey/index.md)
- [Your first workflow](../first-workflow/index.md)
