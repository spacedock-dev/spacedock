---
title: "Bridge an external tracker"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-28 20:28:40"
---

# Bridge an external tracker

Your backlog may already live somewhere else: Linear, GitHub Issues, another ticket ledger. The bridge is asking, not configuring: tell the first officer to intake an external item ("intake GitHub PR #134") and it files an entity carrying the reference. The external system keeps owning intake, discussion, and assignment; Spacedock stays the execution workflow.

Two frontmatter fields carry the reference:

```
issue: ENG-123
source: linear
```

`issue` is the human-facing external reference; `source` records where the entity came from. Keep the tracker out of Spacedock's stage semantics: sync through entity creation, state changes, and stage reports, not tracker-specific stage rules.

## Sitemap

- [Mods & standing teammates](../mods-and-standing-teammates/index.md)
- [Multiple workflows](../multi-workflow/index.md)
- [Split-root state](../split-root-state/index.md)
- [Refit a workflow](../refit/index.md)
