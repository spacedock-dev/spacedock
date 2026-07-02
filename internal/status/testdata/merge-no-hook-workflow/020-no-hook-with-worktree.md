---
id: "020"
title: No merge hook registered, worktree recorded
status: implementation
worktree: .worktrees/020-no-hook-with-worktree
score: "0.40"
source: roadmap
---
# No merge hook registered, worktree recorded

Same as the sibling no-hook entity, but a `worktree` frontmatter field is
recorded. `merge guard --verdict passed` finalizes; the finalized line must
carry BOTH the worktree-removal next-step clause and the no-merge-hook manual
`--no-ff` clause — the two conditions are independent.
