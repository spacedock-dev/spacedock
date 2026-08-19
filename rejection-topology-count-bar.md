---
id: 12zytd0ksb47r3g61qw1s5nc
title: The rejection topology bar counts routing events instead of judging conduct
status: backlog
source: "Captain CL, 2026-08-19, after PR #736's live lane reddened on both hosts and neither red was attributable to the diff. Failing assertion: codex-live rejection-flow, observed=[rejection-worker-topology], 'the reuse branch owes 8 routing events, the run produced 10'. Retained evidence preserved at /tmp/9g-evidence (run 32270990171) and readable in each host's rejection-topology.tsv."
started:
completed:
verdict:
score:
worktree:
issue:
---

`rejection-worker-topology` grades the rejection journey by counting routing events against an exact number. On the reuse branch that number is 8. A First Officer that ran one extra review round produced 10 and graded RED.

The retained topology shows the run, and it is not flailing:

```
0  spawn  implementation      5  done   validation
1  done   implementation      6  reuse  implementation
2  spawn  validation          7  done   implementation
3  done   validation          8  reuse  validation
4  reuse  validation   <--    9  done   validation
```

Events 4 and 5 are the extra pair: the FO ran validation, re-ran validation, then routed the correction and re-reviewed. Re-running a reviewer before routing a rejection is defensible conduct, and an exact-count bar has no way to say so.

## Why this is the fourth instance, not an isolated red

Three graders in this release reddened compliant conduct for the same underlying reason — a pattern written against one observed shape, then treated as the definition of correct behavior:

- `filing-recognizer-newline-terminator` (merged #731): the create recognizer's terminator class omitted `\n`, so a real filing counted 0.
- `decide-dispatch-build-count-bar` (merged #732): `assertRecordedGateLifecycle` demanded exactly one dispatch build, so a corrective rebuild graded RED. That entity's finding was that the exact count was an assumption nobody had justified.
- `own-claude-early-rejection-round-record` (`zf7`, backlog): `rejectionRoundSuccess` pins `entries=4`, so a successful early record reports as "resolved launcher never invoked".

This entity is the same defect in the topology oracle. `6ht` is the closest precedent and its remedy is the model: it replaced a bare count with a rule that names which shapes are acceptable and why, and kept every genuinely bad shape red.

## What must NOT be lost

`6ht`'s validation proved the widened bar still reds flailing, waste, and final failure. Any change here owes the same proof. The bar exists because worker topology is real evidence of conduct: a journey that never spawns a validator, or that abandons a rejection, must stay red. Widening the bar until nothing fails is worse than the current false red, because a false red gets investigated and a permissive bar does not.

## Ideation must settle

1. **What the 8 is.** Derive where the number came from and whether any run ever justified it, exactly as `6ht` did for its 1/1. If it encodes a happy path rather than a contract, say so.
2. **Which shapes are acceptable.** An extra review round before routing a correction appears benign. Decide whether it is, and name the shapes that are not: no validator, no correction round, a rejection never routed, unbounded rounds. The claude branch in the same run took `fresh` with 8 events — so the bar must hold across both branches without special-casing.
3. **Whether the bar should count at all.** The alternative to a better count is an ordering rule: the journey must contain a spawn, a rejection, a routed correction, and a re-review, in order, with bounded repetition. Weigh that against a widened count on necessity.
4. **Relationship to `9g`.** `persist-codex-rollout-for-rejection-topology` was filed on the belief that this evidence dies with the isolated `CODEX_HOME`. It does not — `rejection-topology.tsv` is retained per host in the CI artifact. Recommend whether `9g` is superseded, and say plainly that its premise was wrong.

## Out of scope

`zf7`'s `rejection-round-missing` mislabel, which is the same journey's other red and has its own entity. The rejection journey's fixture text. Any change to what the rejection flow itself does.

## Value

A live lane is worth its cost only when a red means something. Four graders in one release have reddened correct behavior, and each one teaches every reader to discount the next red. This is the one still shipping.
