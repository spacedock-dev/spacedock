---
title: "Install"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-25 16:33:06"
---

# Install Spacedock

Spacedock works with a coding agent you already have: Claude Code, Codex, or Pi. Install one of those first.

```
brew tap spacedock-dev/tap
brew install spacedock
# Edge channel instead (tracks prereleases; conflicts with the stable cask):
# brew install spacedock-dev/tap/spacedock@next
```

Match the cask to the plugin channel you install below — the edge plugin pins a binary minor that no stable release satisfies.

`brew install` also pulls in `agentsview` (it powers `/spacedock:survey`). The optional sandbox, safehouse, is installed separately — see [Sandboxing](../../reference/sandbox/).

```
curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
```

Installs a checksum-verified binary to `~/.local/bin`. This is the stable channel: it resolves the latest stable release and prints the tag it resolved before it installs anything.

To track edge, set `SPACEDOCK_CHANNEL=edge`. Edge resolves the newest release including prereleases, and installs the edge-stamped binary:

```
curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | SPACEDOCK_CHANNEL=edge sh
```

Match the channel to the skills you run. The edge plugin requires an edge-line binary minor, so a stable binary with edge skills aborts at first-officer boot. On macOS, `brew install spacedock-dev/tap/spacedock@next` is the Homebrew equivalent; casks are macOS-only, so on Linux this script is the only edge path.

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

# Edge (tracks main) — marketplace named `spacedock-edge`, entry still `spacedock`
claude plugin marketplace add spacedock-dev/marketplace@edge
claude plugin install spacedock@spacedock-edge
```

The channel is the marketplace name; the entry name stays `spacedock` on both channels (it equals the plugin's own `name`, so the host's entry-name vs plugin-name check passes). Each channel adds its own marketplace source — the stable marketplace lives at the repo root (named `spacedock`), the edge one on the `@edge` branch (named `spacedock-edge`) — so the `@edge` ref is what registers the `spacedock-edge` marketplace the edge entry resolves from. Codex installs the same way with `codex plugin add`.

Set `SPACEDOCK_MARKETPLACE_SOURCE` to install from a local or alternate marketplace instead of the default `spacedock-dev/marketplace` — useful for dogfooding a marketplace change before it reaches the production marketplace:

```
SPACEDOCK_MARKETPLACE_SOURCE=/path/to/local/marketplace spacedock install --host codex
```

When the resolved `spacedock` executable sits at the root of a checkout whose host manifest names `spacedock` and points to an existing skills directory, the launcher selects that adjacent checkout automatically. A source build at `./spacedock` therefore needs no `--plugin-dir` for either Claude or Codex. A release binary in a `bin/` directory, or a checkout with a missing or invalid host manifest, keeps the normal installed-plugin resolution and compatibility gate.

An explicit `--plugin-dir <checkout>` before `--` overrides adjacent discovery. On Claude and Pi this is an ephemeral, install-free override — it bypasses installed-plugin resolution for that one launch and does not wrap the launch in the safehouse sandbox; use it for plugin development, not as an install substitute. On Claude, a `--plugin-dir` after `--` remains a native additional session plugin: it does not replace Spacedock or suppress the compatibility gate.

Codex has no such flag on its own CLI, so `spacedock codex --plugin-dir <checkout>` and `spacedock install --host codex --plugin-dir <checkout>` build a local marketplace from the checkout and install it under the binary's own channel (`spacedock` stable / `spacedock-edge` edge — matching whatever `spacedock codex` would otherwise install), then launch. This IS a persistent install and Spacedock makes it exclusive across Codex channels: the selected channel replaces any existing stable or edge Spacedock Codex plugin so `$spacedock:*` skills resolve from the selected install. It is also a point-in-time snapshot: editing the checkout afterward has no effect until the command is re-run. The command prints an advisory that names the selected channel and notes that the reported version reflects the checkout's checked-in manifest, not necessarily its current HEAD. Codex does not accept a post-`--` `--plugin-dir`; install additional Codex plugins persistently with `codex plugin add`.

## Sandboxing

See [supported sandboxes](../../reference/sandbox/).

## Troubleshooting

Run `spacedock doctor`.

## Next

Read about the [survey report](../survey/) to understand your usage pattern with coding agents, or start with [your first workflow](../first-workflow/).

## Sitemap

- [Survey your project](../survey/index.md)
- [Your first workflow](../first-workflow/index.md)
