---
id: "050"
title: Rejected entity with an in-flight mod-block
status: implementation
verdict: rejected
mod-block: merge:local-merge
score: "0.30"
source: roadmap
---
# Rejected entity with an in-flight mod-block

`verdict: rejected` but a `mod-block` is set: the verdict escape relaxes only the
merge-hook pr-requirement, NOT the policy-independent mod-block-pending guard. A
live mod-block must still refuse `--archive` here — mirroring how `merge: local`
relaxes only the pr-requirement while the mod-block guard survives.
