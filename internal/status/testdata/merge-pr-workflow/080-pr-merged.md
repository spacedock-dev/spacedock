---
id: "080"
title: Entity whose PR merged but state not yet terminalized (non-armed)
status: implementation
score: "0.50"
source: roadmap
pr: pr-merge:99
---
# Entity whose PR merged but state not yet terminalized (non-armed)

The PR landed and the FO recorded the merge as a `pr-merge:{number}` sentinel —
the local "this PR already merged" signal the startup/idle pr-merge hook writes on
MERGED detection. The `mod-block` is EMPTY (a re-validation bounce cleared it, the
stranded-non-armed case). `merge guard` must FINALIZE this entity (terminalize +
archive) from the non-armed state, keying off the merge sentinel — NOT signal
blocked. A bare open-PR reference (`#42`) is the contrast: that still blocks.
