---
id: "150"
title: Entity finalizing from a merged sentinel with a worktree recorded
status: implementation
worktree: .worktrees/150-worktree-finalize
score: "0.50"
source: roadmap
pr: pr-merge:88
---
# Entity finalizing from a merged sentinel with a worktree recorded

The PR landed (`pr: pr-merge:88`) and a `worktree` frontmatter field is still
recorded (the FO has not yet removed it). `merge guard` finalizes; the finalized
line must carry the worktree-removal/branch-cleanup/teardown next-step clause.
Since a merge hook IS registered (`local-merge`) and a sentinel IS recorded, the
no-merge-hook manual-merge clause must NOT appear.
