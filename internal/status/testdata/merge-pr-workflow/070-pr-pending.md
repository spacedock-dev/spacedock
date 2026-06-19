---
id: "070"
title: Entity with an open PR and an in-flight mod-block
status: implementation
mod-block: merge:local-merge
score: "0.50"
source: roadmap
pr: "#42"
---
# Entity with an open PR and an in-flight mod-block

The merge hook ran and BLOCKED on an open PR (`pr: #42`), leaving the `mod-block`
in flight. This is the Phase B blocked-by-state-delta case: `merge guard` must
detect the set `pr`, signal `blocked`/`await-pr`, leave the `mod-block` intact,
and NOT terminalize or archive. The entity stays at its non-terminal status until
the PR lands and the FO re-runs the ceremony.
