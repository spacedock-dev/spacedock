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
gates:
    version: 1
    records:
        - id: gate:12zytd0ksb47r3g61qw1s5nc:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:12zytd0ksb47r3g61qw1s5nc-backlog-1
              briefing:
                id: briefing:12zytd0ksb47r3g61qw1s5nc:backlog:attempt-1:revision-1
                digest: sha256:8740c7ebfa0bb6cf4a9c93a05bc7cb8b5afa7b1a10fce93e97a58dd1b0cf159d
                request-digest: sha256:8857c6c8494c10a0e340a88a64471a79e9cbcd5bdda52eee180209f04d9998d9
                room-ref: ./rejection-topology-count-bar/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:12zytd0ksb47r3g61qw1s5nc:backlog:1
                briefing: briefing:12zytd0ksb47r3g61qw1s5nc:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-19T18:16:30.467253Z"
                decision: approve
                reason: 'Captain approved in chat: ''if the fable opinion lands reasonable, dispatch those two on to the stack of 736 so we can run tip CI lane to verify.'' The review confirmed the diagnosis, supplied its missing argument, settled the remedy shape, and corrected the 9g supersession claim.'
              application:
                target-stage: ideation
                state: pending
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

## Adversarial review (2026-08-19, fable ensign, read-only)

Diagnosis CONFIRMED, and it supplied the argument this entity asserted without one.

**Why the extra pair is benign — on evidence, not assertion.** Order held: the sequence is the required chain with one validation round repeated in place. Identity held: all three validation routes reach the same live validator handle, and the re-review target differs from the fix producer, so it is not self-review. And decisively, **every other oracle on the codex lane was green** — the round record ran AFTER the correction with the complete `entries=4` (codex-exec.jsonl line 62), single publication, one round room, gate prepared, all four byte-exact files intact. The extra validation round left ZERO durable footprint. `observed=[rejection-worker-topology]` was the lane's only red. That triad — shape, identity, outcome — is the argument, and it is also the design principle for the replacement bar, because intent is not observable.

**The remedy is an ordering rule. Item 3 must not be read as a live choice.** Framing "widened count vs ordering rule" as alternatives invites someone to pick the count, and the count is the original sin again: it encodes today's two observed shapes as the contract, is position-fragile, and cannot say WHICH extra pair is acceptable or where.

**And the #732 precedent transfers as a PATTERN, not a mechanism.** `recordedGateBuildAttemptsAcceptable` (recorded_gate_lifecycle_test.go:98-117) accepts a corrective rebuild and reds an identical successful rebuild as waste. Transferred literally, that REDS this codex run — the extra validation round is an apparently-identical successful re-run of a stage that already succeeded, exactly #732's waste case. The cases differ in kind: a dispatch build is semantics-free, so an identical rerun can only be waste; a review round is a semantic act with legitimate re-run reasons, and its harmlessness is independently checkable through the durable oracles that all passed. Anyone who says "do what #732 did" gets the wrong answer here.

**The shape to build.** Parse the event list into rounds (a route event plus its `done`), then grade a grammar: implementation round; one-or-more validation rounds; exactly one rework implementation round; one-or-more re-review validation rounds; last round is a completed validation. Routing rule per round, which generalizes both branches without special-casing: the FIRST round of a stage is a spawn; every subsequent round of that stage follows the branch (reuse to reuse on the same handle, fresh to spawn). Bounded repetition in #732's spirit: at most one extra validation round total.

Falsifying controls that must stay red: no validator; no rework; rework never re-reviewed; re-review before rework; self-review (re-review reaches the fix producer); an extra IMPLEMENTATION round (in this one-cycle fixture a second rework means the fix did not hold); three-plus validation rounds. Note the existing table test's "no re-review after rework" case currently reds only via the length check — the grammar replaces that with "must end in a completed re-review", a strictly better diagnostic.

Framing correction: this is not count-versus-order. The grader already grades order. It is flat-event-matching versus round-level parsing.

**Correction to item 4 — 9g is only PARTIALLY superseded.** This entity's claim that "the evidence does not die" is overstated. Only the topology TSV digest survives in the artifact. The raw parent rollout — the only thing that could establish the extra round's INTENT — still dies with the isolated `CODEX_HOME`, and the reviewer hit exactly that wall: it can prove shape, identity and outcome, but cannot distinguish "re-ran the reviewer to confirm the rejection" from "briefly mis-routed the correction, then corrected". So 9g is superseded for shape attribution and NOT for conduct forensics. Ideation must not kill 9g on the premise this entity originally stated — which would repeat, against 9g, precisely the unexamined-premise failure this entity exists to complain about.

## Out of scope

`zf7`'s `rejection-round-missing` mislabel, which is the same journey's other red and has its own entity. The rejection journey's fixture text. Any change to what the rejection flow itself does.

## Value

A live lane is worth its cost only when a red means something. Four graders in one release have reddened correct behavior, and each one teaches every reader to discount the next red. This is the one still shipping.
