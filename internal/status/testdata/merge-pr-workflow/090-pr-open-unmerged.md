---
id: "090"
title: Entity with an open unmerged PR and no mod-block
status: implementation
score: "0.50"
source: roadmap
pr: "#42"
---
# Entity with an open unmerged PR and no mod-block

A bare PR reference (`pr: #42`) for an OPEN, NOT-yet-merged PR, with an EMPTY
mod-block. `merge guard` must NEVER finalize on `pr`-presence alone — a bare
open-PR reference means the PR is still in review. The verb must signal
`blocked`/`await-pr`, leave the entity at its non-terminal status, and NOT
archive. Only a merge sentinel (`pr-merge:{n}` / `local-merge:{sha}`) finalizes.
This pins the premature-finalize bug: archiving a task before its PR landed.
