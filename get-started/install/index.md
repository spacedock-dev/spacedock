---
title: "Install"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-06-20 20:38:16"
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
claude plugin marketplace add spacedock-dev/marketplace
claude plugin install spacedock@spacedock
```

## Sandboxing

See [supported sandboxes](../../reference/sandbox/).

## Troubleshooting

Run `spacedock doctor`.

## Next

Read about the [survey report](../survey/) to understand your usage pattern with coding agents, or start with [your first workflow](../first-workflow/).

## Sitemap

- [Survey your project](../survey/index.md)
- [Your first workflow](../first-workflow/index.md)
