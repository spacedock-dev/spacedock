---
id: "060"
title: Passed entity with no sentinel and no mod-block
status: implementation
verdict: passed
score: "0.40"
source: roadmap
---
# Passed entity with no sentinel and no mod-block

`verdict: passed` with no `pr` sentinel and no `mod-block`. This is exactly the
case the merge-hook guard exists for: an entity claiming an accepted outcome
without having run the merge ceremony. The `verdict: rejected` exemption must NOT
widen to other verdicts — `--archive` of this entity under the default `merge: pr`
policy must still refuse without `--force`. Pins the over-wide-exemption
complement of the rejected-verdict escape.
