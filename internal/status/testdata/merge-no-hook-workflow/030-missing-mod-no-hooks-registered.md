---
id: "030"
title: Mod-block names a merge mod but NO merge hook is registered at all
status: implementation
mod-block: merge:ghost-merge
score: "0.40"
source: roadmap
---
# Mod-block names a merge mod but NO merge hook is registered at all

`mod-block: merge:ghost-merge`, and this workflow (unlike the sibling
merge-pr-workflow fixture) registers NO merge hook whatsoever — `_mods/` has no
`## Hook: merge` mod at all, not merely a mismatched one. This is the likeliest
real-world D5 trigger: the deleted mod file WAS the workflow's only merge hook.
`merge guard` must still refuse (not silently finalize via the len(mergeHooks)==0
default-Phase-C path) since no sentinel records a landed merge.
