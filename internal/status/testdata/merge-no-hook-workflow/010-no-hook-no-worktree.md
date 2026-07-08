---
id: "010"
title: No merge hook registered, no worktree recorded
status: implementation
score: "0.40"
source: roadmap
---
# No merge hook registered, no worktree recorded

No `_mods/` merge hook is registered for this workflow, no `pr` sentinel, no
`mod-block`, and no `worktree` frontmatter field. `merge guard --verdict passed`
reaches Phase C default and finalizes; the finalized line must name the manual
`--no-ff` merge onto trunk since nothing automated it, with no worktree-removal
clause (there is no worktree to remove).
