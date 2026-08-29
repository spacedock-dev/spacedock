---
title: "Split-root state"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-29 00:32:28"
---

# Split-root state

A busy workflow generates constant small state commits. Split-root keeps them off your code branch. The README and stage declarations stay on your main branch; the mutable entity state (frontmatter updates, stage reports, archive moves) lives in a separate state checkout. Your code history stays clean, and several agents or operators can drive the same workflow at once.

The README [opts in with one `state:` field](../../concepts/workflows-and-entities/#keep-workflow-state-off-your-code-branch); the split is transparent, and you read the workflow exactly as you would any other. On a fresh clone the state checkout is absent; run `spacedock state init` to restore it. The shipped [`docs/dev` workflow](https://github.com/spacedock-dev/spacedock/tree/main/docs/dev) runs split-root, a live example.

## Concurrent writers

The agents follow the commit and sync discipline that keeps concurrent writers from clobbering each other; it is theirs to follow, not yours. The one case that reaches you: conflicting edits to the same entity halt for your call rather than auto-resolving.

## Sitemap

- [Mods & standing teammates](../mods-and-standing-teammates/index.md)
- [Multiple workflows](../multi-workflow/index.md)
- [Bridge an external tracker](../external-tracker/index.md)
- [Refit a workflow](../refit/index.md)
