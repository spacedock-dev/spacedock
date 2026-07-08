---
id: "130"
title: Entity blocked on a missing merge mod but carrying a landed-merge sentinel
status: implementation
verdict: passed
mod-block: merge:ghost-merge
pr: pr-merge:77
score: "0.40"
source: roadmap
---
# Entity blocked on a missing merge mod but carrying a landed-merge sentinel

`mod-block: merge:ghost-merge` names a merge mod that is NOT registered under
`_mods/`, but `pr: pr-merge:77` is a well-formed merge sentinel — the merge
genuinely landed before the mod file was deleted. `merge guard` must finalize as
today: a recorded sentinel honestly proves the ceremony ran, so the missing mod
file is not a stuck state.
