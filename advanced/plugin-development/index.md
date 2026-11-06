---
title: "Local plugin development"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-06 00:33:04"
---

# Local plugin development

For normal setup, use [Installation](../../get-started/install/).

## Local checkout

Use `--plugin-dir <checkout>` before `--` to select a local plugin checkout.

For Claude Code and Pi, the override applies to that session and does not use the Safehouse sandbox.

For Codex, it creates a persistent installation that replaces the existing Spacedock plugin. Run the command again after changes to the checkout.

Codex does not accept `--plugin-dir` after `--`.

## Alternate marketplace

Set `SPACEDOCK_MARKETPLACE_SOURCE` to select a local or alternate marketplace:

```
SPACEDOCK_MARKETPLACE_SOURCE=/path/to/local/marketplace spacedock install --host codex
```

## Sitemap

- [Mods & standing teammates](../mods-and-standing-teammates/index.md)
- [Multiple workflows](../multi-workflow/index.md)
- [Split-root state](../split-root-state/index.md)
- [External trackers](../external-tracker/index.md)
- [Workflow refits](../refit/index.md)
