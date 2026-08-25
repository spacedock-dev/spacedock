---
title: "Refit a workflow"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-25 04:42:39"
---

# Refit a workflow

When you upgrade Spacedock, run `/spacedock:refit path/to/workflow` to bring the workflow's generated files up to date while leaving your edits in place. Nothing is auto-replaced: you see a diff and decide, file by file, and if a schema change affects your entities, refit proposes the migration and waits for your approval. Git is the safety net; ask the agent about anything else.

## Sitemap

- [Mods & standing teammates](../mods-and-standing-teammates/index.md)
- [Multiple workflows](../multi-workflow/index.md)
- [Split-root state](../split-root-state/index.md)
- [Bridge an external tracker](../external-tracker/index.md)
