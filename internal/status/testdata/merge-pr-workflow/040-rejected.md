---
id: "040"
title: Rejected entity with no sentinel and no mod-block
status: implementation
verdict: rejected
score: "0.40"
source: roadmap
---
# Rejected entity with no sentinel and no mod-block

`verdict: rejected` with no `pr` sentinel and no `mod-block`. A rejected entity
never ran the merge ceremony — there is no PR to require, no merge to gate on —
so the merge-hook requirement is vacuous for it. Both `--set status=done` and
`--archive` must agree on this entity: under the default `merge: pr` policy each
surface exempts `verdict: rejected` and lets the entity through without `--force`.
