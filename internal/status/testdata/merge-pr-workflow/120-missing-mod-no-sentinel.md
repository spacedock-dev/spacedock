---
id: "120"
title: Entity blocked on a merge mod that no longer exists, no sentinel
status: implementation
verdict: passed
mod-block: merge:ghost-merge
score: "0.40"
source: roadmap
---
# Entity blocked on a merge mod that no longer exists, no sentinel

`mod-block: merge:ghost-merge` names a merge mod that is NOT registered under
`_mods/` (the file was deleted mid-ceremony). No `pr` sentinel records a landed
merge, and `verdict` is not `rejected`. `merge guard` must refuse rather than
silently finalize — clearing the block here would archive the entity without the
hook ever having run.
