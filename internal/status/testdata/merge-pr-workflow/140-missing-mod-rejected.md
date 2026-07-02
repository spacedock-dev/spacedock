---
id: "140"
title: Rejected entity blocked on a merge mod that no longer exists
status: implementation
verdict: rejected
mod-block: merge:ghost-merge
score: "0.40"
source: roadmap
---
# Rejected entity blocked on a merge mod that no longer exists

`mod-block: merge:ghost-merge` names a merge mod that is NOT registered under
`_mods/`, and `verdict: rejected` — the entity never merged, so the missing mod
file cannot have stranded a landed merge. `merge guard` must finalize as today;
the rejected-verdict escape takes priority over the missing-mod refusal.
