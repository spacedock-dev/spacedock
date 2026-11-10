---
title: "Split-root state"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-10 16:50:35"
---

# Split-root state

Split-root state keeps workflow progress commits off your code branch. Workflow rules stay with the code; work records live in a separate Git checkout.

Ask your agent to configure this layout. On a fresh clone, restore the state checkout with:

```
spacedock state init
```

Conflicting edits to the same work item stop for your decision.

## Sitemap

- [Mods & standing teammates](../mods-and-standing-teammates/index.md)
- [Multiple workflows](../multi-workflow/index.md)
- [External trackers](../external-tracker/index.md)
- [Workflow refits](../refit/index.md)
- [Local plugin development](../plugin-development/index.md)
