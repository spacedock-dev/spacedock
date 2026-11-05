# Local plugin development

For normal setup, use [Installation](../get-started/install.md).

## Local checkout

Use `--plugin-dir <checkout>` before `--` to select a local plugin checkout.

For Claude Code and Pi, the override applies to that session and does not use the Safehouse sandbox.

For Codex, it creates a persistent installation that replaces the existing Spacedock plugin. Run the command again after changes to the checkout.

Codex does not accept `--plugin-dir` after `--`.

## Alternate marketplace

Set `SPACEDOCK_MARKETPLACE_SOURCE` to select a local or alternate marketplace:

```bash
SPACEDOCK_MARKETPLACE_SOURCE=/path/to/local/marketplace spacedock install --host codex
```
