# Spacedock

**Spacedock is a multi-agent orchestrator where nothing ships without a decision.** It lives within your existing harness — Claude Code, Codex, or Pi.

[TODO] — the Home page: the pitch, what's different, and where to go next (Get started, Concepts, Running workflows, Contributing).

## For agents

Spacedock's docs are read by agents too — a user's first officer parsing these docs is itself an agent. The build emits an agent-readable surface: an `llms.txt` index of the docs at the site root (discoverable from each page's `<head>` via `<link rel="alternate" type="text/markdown">`) and the repo's [agent instructions](agents/index.md).
