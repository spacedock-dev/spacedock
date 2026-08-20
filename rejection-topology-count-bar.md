---
id: 12zytd0ksb47r3g61qw1s5nc
title: The rejection topology bar counts routing events instead of judging conduct
status: validation
source: "Captain CL, 2026-08-19, after PR #736's live lane reddened on both hosts and neither red was attributable to the diff. Failing assertion: codex-live rejection-flow, observed=[rejection-worker-topology], 'the reuse branch owes 8 routing events, the run produced 10'. Retained evidence preserved at /tmp/9g-evidence (run 32270990171) and readable in each host's rejection-topology.tsv."
started: 2026-08-19T18:17:48Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-rejection-topology-count-bar
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
                state: consumed
        - id: gate:12zytd0ksb47r3g61qw1s5nc:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:12zytd0ksb47r3g61qw1s5nc-ideation-1
              briefing:
                id: briefing:12zytd0ksb47r3g61qw1s5nc:ideation:attempt-1:revision-1
                digest: sha256:e8e5cd7362cc21574db4b4fa4da934f144492af6cc18ba48e60f71e4d4521e45
                request-digest: sha256:784486ce093347adc5d8e1f08b38b69ec42a3347e2c17723cd52883772c03940
                room-ref: ./rejection-topology-count-bar/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:12zytd0ksb47r3g61qw1s5nc:ideation:1
                briefing: briefing:12zytd0ksb47r3g61qw1s5nc:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-19T19:24:29.390013Z"
                decision: revise
                reason: 'Captain rejected in chat: ''it''s not about asking permission. it''s about not having that in the fricking AC.'' Same structural defect as zf7: no acceptance criterion requires a targeted live run of the journey this grader governs. Its ACs verify against retained chains and a mutant corpus — which proves the grammar changed, not that the journey passes. Add an AC whose verification IS a targeted local live run of rejection-flow, naming the journey and the expected observed codes. The grammar design itself is accepted and stands.'
            - id: gate-attempt:12zytd0ksb47r3g61qw1s5nc-ideation-2
              briefing:
                id: briefing:12zytd0ksb47r3g61qw1s5nc:ideation:attempt-2:revision-1
                digest: sha256:93e8fd75396dcf75abbb55c4678358f529f19f2c9d20fe20980a54314d07d155
                request-digest: sha256:bc0c1b70b4e3d6d889cd50718e44f73d3c018adbafc9010dfa25e15903442a77
                room-ref: ./rejection-topology-count-bar/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:12zytd0ksb47r3g61qw1s5nc:ideation:2
                briefing: briefing:12zytd0ksb47r3g61qw1s5nc:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-19T22:57:16.085558Z"
                decision: approve
                reason: 'Captain approves onto the #736 stack: replace the exact 8-event count with a round-parsing grammar so the bar stops reddening correct topologies.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:12zytd0ksb47r3g61qw1s5nc:validation
          stage: validation
          attempts:
            - id: gate-attempt:12zytd0ksb47r3g61qw1s5nc-validation-1
              briefing:
                id: briefing:12zytd0ksb47r3g61qw1s5nc:validation:attempt-1:revision-1
                digest: sha256:00c97f8084e975badc13cb71ae1a1d51950f17fe86d1b63ca83c511d119e3c50
                request-digest: sha256:d0ea2a75218518ffc999df9beecf9290c0a0efa966ed08a69d2864d81b2cd966
                room-ref: ./rejection-topology-count-bar/review/validation/briefing-1
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

## Ideation (2026-08-19)

### What the 8 encodes, and what justified it

`rejectionChain` (claude_runtime_helpers_test.go:408) hardcodes four rounds × two events. Two clauses are folded into that 8, and they have different standing:

- **The chain and its order** — candidate, review, routed rework, re-review, ending in a reviewer verdict — IS the contract: `feedback-rejection-flow/SKILL.md` steps 1 and 4, plus reviewer independence ("never to validate its own fix work") and bounded cycles ("On cycle 3 run no reviewer: escalate to the human instead of a fourth round").
- **Exactly one round per segment** is nowhere in the contract. It came from the one captured rollout the grader was written against (commit 52b94f9aa, 2026-08-16: `testdata/rejection_topology/codex-reuse-chain.jsonl`) plus one synthetic fresh chain. A single live observation justified it; the next live run (32270990171) falsified it — 10 events, order held, identity held, every durable oracle green.

Same finding shape as `6ht`'s 1/1 (`recordedGateBuildAttemptsAcceptable`, recorded_gate_lifecycle_test.go:98): the exact count encoded a happy-path assumption no run had ever been asked to justify, and the first benign deviation reddened compliant conduct.

### Decision: an ordering rule over parsed rounds — review validated, not re-derived

The adversarial review's three decisive claims were checked against the evidence, not accepted on authority:

1. **Shape and identity held** — retained TSV rows 0–9: the required chain with one validation round repeated in place; every validation round reached the validator handle; the re-review handle differs from the fix producer's.
2. **Outcome was clean** — codex-exec.jsonl lines 62–63: `gate record` exit 0 with `entries=4` INCLUDING the `fixed-marker` annotation and the post-fix resolution, i.e. the round was recorded after the correction; single publication; `rejection-worker-topology` was the lane's only red.
3. **#732 transfers as pattern, not mechanism** — verified in code: `recordedGateBuildAttemptsAcceptable` case 2 accepts a second attempt only when the first failed or the command differed, so its literal transfer reds this run's identical successful re-review. The pattern (replace a count with a rule naming acceptable shapes) transfers; the mechanism does not, because a dispatch build is semantics-free while a review round is a semantic act whose harmlessness the durable oracles independently vouched for.

No refutation. Two refinements to the review's text:

1. Its falsifying control "three-plus validation rounds" needs a precise boundary: three validation rounds TOTAL is the retained codex run and must green; what reds is exceeding the one-repeat budget — and any single segment reaching three rounds necessarily does (three in a segment plus the other segment's required one is four total, two repeats).
2. Its self-review control needs a fresh-branch variant: on the reuse branch a retargeted re-review usually trips handle-consistency first, so the independent-review check must be proven on the collapsed shape (all validation rounds targeting the producer, internally consistent) and on the fresh branch, where handle comparison is otherwise meaningless.

A widened count was considered and rejected for the reasons the review names: it re-encodes today's two observed shapes as the contract, is position-fragile (the identity pairs {0,4}/{2,6} shift per shape), and cannot say which extra pair is acceptable or where.

### The grammar (normative)

Parse the event list into rounds, then grade the grammar. Fail closed: any unpairable event is a graded red naming its position.

- **Round:** a dispatch event (`spawn`/`reuse`) immediately followed by a `done` of the same stage and target.
- **Sequence:** `implementation, validation+, implementation (the rework), validation+` — exactly one rework; the journey ends in a completed validation round.
- **Repetition budget:** at most ONE repeated validation round total, `(vPre−1)+(vPost−1) ≤ 1`. The diagnostic declares it as a budget citing this entity, so a future legitimate second repeat produces a red naming the knob instead of masquerading as a contract violation.
- **Routing rule (both branches, no special-casing):** a stage's FIRST round is a spawn; every subsequent round of that stage follows the branch key — reuse: a followup reaching the same handle the stage's first round opened; fresh: a fresh spawn. Codex (reuse, 10 events) and Claude (fresh, 8 events) both green under this one rule; each branch's chain still reds under the other, and all four probe fail-safe cases keep their verdicts. A pre-rework confirm re-run routed to a live worker WITHOUT a passing probe still reds on the fresh branch — that is the contract's fail-safe, not a false red.
- **Independent review:** no validation round's target may equal any implementation round's target.

Every diagnostic names the violated invariant. In particular the "no re-review after rework" shape stops redding via the bare length message ("owes 8 routing events") and reds as the missing re-review it is.

### Spike (riskiest mechanism, exercised)

Throwaway Go program implementing exactly the grammar above, run against the retained run-32270990171 TSVs and a mutant corpus: **20/20**. Greens: retained codex chain (10), retained claude chain (8), captured fixture chain (8), post-rework repeat. Reds, each under an invariant-naming diagnostic: no validator; no rework; rework never re-reviewed; re-review before rework (unpairable); self-review via retarget; reviewer-collapsed (internally consistent reuse, all reviews reach the producer); fresh-branch self-review; single-worker chain; second rework; two repeats in one segment; two repeats split across segments; both cross-branch chains; both probe fail-safe violations; two-validators-up-front (unpairable). The corpus seeds the implementation's table test.

### Acceptance criteria

1. **(Value, measured against a baseline that can move the wrong way)** Replayed through the grader, the retained run-32270990171 chains grade GREEN on both branches — 2/2, from today's 1/2 — while the falsifying-control corpus (at minimum the seven checklist shapes, both self-review variants, both cross-branch chains, both probe fail-safe violations) stays at 100% RED under code `rejection-worker-topology`. A single control greening is the widened-bar failure and fails the table test.
   *Test:* offline table test; the 10-event chain enters as the captured fixture chain with the repeat pair inserted at the observed position, cited to the retained TSV.
2. The conforming set is exactly the grammar above, and every red names its violated invariant: the missing-re-review shape's diagnostic names the missing re-review and no diagnostic path emits a bare owed-event-count message.
   *Test:* per-case diagnostic substring assertions; the missing-re-review case asserts presence of the re-review phrasing and absence of "owes 8".
3. One branch-parameterized routing rule replaces the per-position chains: mutual exclusivity and all four probe fail-safe verdicts hold unchanged.
   *Test:* `TestRejectionTopologyBranchesAreMutuallyExclusive` and `TestRejectionTopologyProbeFailSafe` pass with assertions unmodified.
4. Independent review holds at round level on both branches, including the reviewer-collapsed shape and the fresh-branch variant.
   *Test:* `reviewerCollapsedRoutes` kept; new fresh-branch self-review case.
5. **(Scope fence)** Extraction and the digest are untouched: `TestRejectionTopologyExtractsCapturedCodexChain`'s extraction assertions and the `rejectionTopologyDigest` TSV format are byte-identical; the diff touches grading and its proof only.
   *Test:* extraction test unmodified; digest function untouched in the diff.
6. **(Live journey)** A targeted local live run of the **rejection-flow** journey against the built change carries NO `rejection-worker-topology` in its observed codes, and its retained `rejection-topology.tsv` parses into a chain the grammar accepts.
   *Verification — this AC IS the run, performed at implementation or validation (it needs the built change):* `SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 30m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle/` with `SPACEDOCK_BIN` resolving to an installed binary (a stale path fails in seconds and is not a journey result) and `SPACEDOCK_LIVE_ARTIFACT_DIR` set so the digest survives the run. Expected observed codes, not "passes": `rejection-worker-topology` ABSENT. Any red under that code fails this criterion, and the diagnostic plus the retained digest attribute it — the old bare-count message ("owes 8 routing events") means the grammar did not land; an invariant-naming message over a digest outside the grammar means live non-conforming conduct, a reportable journey finding. Codes owned by other oracles (`zf7`'s `rejection-round-incomplete` in particular) are separate outcomes: reportable, neither satisfying nor failing this AC.
   *Unavailability is named, never silent:* the codex lane refuses to run in this sandbox — it cannot create its isolated `CODEX_HOME` at either candidate path — so the REUSE branch, the branch that produced the falsifying 10-event chain, stays live-unproven locally. That gap is a stated finding in the implementation/validation report, covered by AC-1's deterministic replay of the retained reuse chain and corroborated by the stack-tip CI codex lane (expected reading there: no `rejection-worker-topology` code). A 429/session limit or connection drop is a provider error, never a grade: rerun, or report the quota as the named unavailability with this AC unmet.

### Test plan

Grammar proof is offline: `go test ./internal/ensigncycle -run RejectionTopology` (plus `-race`), grading stated route sequences — no new fixtures, no new mechanisms. Journey proof is AC-6's targeted local live run of rejection-flow on the claude lane, run at implementation or validation against the built change; offline replay proves the grader changed, only the journey proves the journey passes, and those are two different claims. The tip-CI run over the #736 stack corroborates the codex reuse branch that cannot run locally: expected reading is NO `rejection-worker-topology` code on either lane, with the claude lane RED under `zf7`'s `rejection-round-incomplete` — that entity's honest outcome, not this one's. A topology red anywhere names its invariant and is investigable on the digest alone.

New-mechanism accounting: the only new mechanism is the round parser, serving AC-1/AC-2. Simplest alternative considered: enumerate the three conforming chains per branch — (1,1), (2,1), (1,2) validation-round distributions — and flat-match each. Rejected: the identity pairs shift per template (the position-fragility this entity exists to remove), a red reports "matched none of three chains" instead of the violated invariant, and the budget lives implicitly in the enumeration. No other unverified mechanism: extraction, branch keying, and digest persistence are proven and untouched ("no spike needed" beyond the grammar spike recorded above).

### Expected surface

Net **+105 LOC** (~+180 insertions, ~−75 deletions) across **2 files**; tolerance ±50 net, +1 file. AC-6 adds no code surface: the targeted runner, artifact retention, and digest persistence already exist — it adds a run obligation, not a mechanism. Breakout — product: 0 (the oracle is test-harness code); grader (`internal/ensigncycle/claude_runtime_helpers_test.go`: `rejectionChain` + `assertRejectionWorkerTopology` replaced by parser + grammar): ~+45 net; proof (`internal/ensigncycle/shared_rejection_topology_table_test.go`: new green/red cases, updated diagnostic substrings): ~+60 net; docs: 0 — no doc describes this oracle (grep over `docs/` confirms), so no doc diff is owed. Declared semantic changes: the grading semantics of the `rejection-worker-topology` oracle (conforming set widens by exactly the one-repeat budget; every diagnostic reworded to name its invariant). No command grammar, stored-format (TSV digest unchanged), authority, or product runtime changes.

### 9g: partial supersession, not retirement

This entity's original item-4 claim was overstated, as the adversarial review found. Only the topology TSV digest survives CI; the raw parent rollout — the only artifact that could establish the extra round's INTENT (confirm-the-rejection vs briefly-misroute-then-correct) — still dies with the isolated `CODEX_HOME` under `t.Cleanup`. Recommendation: **rescope `9g`, do not retire it.** Superseded: shape attribution (its "the red cannot be diagnosed" premise — this ideation diagnosed the red from the retained digest alone). Not superseded: conduct forensics — persisting the parent rollout (or its followup-payload slice) into run artifacts remains 9g's live scope. Killing 9g on this entity's original premise would repeat, against 9g, the unexamined-premise failure this entity exists to correct.

## Stage Report: ideation

- DONE: Settle the remedy as an ordering rule over parsed rounds, or refute the review's case for it with evidence. Do NOT ship a widened count: it re-encodes today's observed shapes as the contract and cannot say which extra pair is acceptable or where.
  Settled FOR the ordering rule after verifying the review's shape/identity claims against the retained TSVs, its outcome claim against codex-exec.jsonl:62-63 (entries=4 incl. fixed-marker, exit 0), and its #732 caveat against `recordedGateBuildAttemptsAcceptable` in code; widened count rejected with reasons recorded.
- DONE: Derive what the 8 encodes and whether any run ever justified it, exactly as decide-dispatch-build-count-bar did for its 1/1.
  8 = contract chain (justified by `feedback-rejection-flow/SKILL.md` steps 1/4) × an exactly-one-round-per-segment clause justified only by the single captured rollout (commit 52b94f9aa) and falsified by the next live run.
- DONE: Prove the new bar still reds every bad shape: no validator, no rework, rework never re-reviewed, re-review before rework, self-review (re-review reaching the fix producer), an extra implementation round, three-plus validation rounds. A bar that reds nothing is worse than the false red it replaces.
  Spike exercised the grammar over the retained live chains plus a 16-red mutant corpus: 20/20, every red under an invariant-naming diagnostic; boundary made precise (three rounds total with one repeat greens — the retained run; any segment at three necessarily exceeds the budget and reds).
- DONE: State the routing rule that generalizes both branches without special-casing — codex took reuse with 10 events, claude took fresh with 8, and both must grade correctly under one rule.
  A stage's first round is a spawn; every subsequent round of that stage follows the branch key (reuse: followup to the handle the stage opened; fresh: fresh spawn) — spike shows both retained chains green under it and all cross-branch/probe-fail-safe verdicts unchanged.
- DONE: Correct the entity's 9g claim in your report: only the TSV digest survives; the raw rollout that could establish intent still dies with the isolated CODEX_HOME. Recommend partial supersession, not retirement.
  Recorded in "9g: partial supersession, not retirement": superseded for shape attribution, not for conduct forensics; recommend rescoping 9g to rollout persistence for intent.

### Summary

Settled the remedy as a round-parsing grammar — `impl, val+, rework, val+`, one-repeat budget, branch-keyed routing rule, round-level independent-review check — validating the adversarial review's conclusions against the retained evidence rather than re-deriving or accepting them. The riskiest mechanism was spiked against both hosts' retained run-32270990171 chains and a falsifying corpus (20/20); the corpus seeds the implementation's table test. Surface declared at net +105 LOC (±50) across 2 test-harness files, no product or doc changes; tip-CI reading recorded as honesty-of-message, not lane color.

## Stage Report: ideation (cycle 2)

- DONE: Settle the remedy as an ordering rule over parsed rounds, or refute the review's case for it with evidence. Do NOT ship a widened count: it re-encodes today's observed shapes as the contract and cannot say which extra pair is acceptable or where.
  Stands from cycle 1 per the gate ("your grammar design is accepted"); see "Decision: an ordering rule over parsed rounds".
- DONE: Derive what the 8 encodes and whether any run ever justified it, exactly as decide-dispatch-build-count-bar did for its 1/1.
  Stands from cycle 1; see "What the 8 encodes, and what justified it".
- DONE: Prove the new bar still reds every bad shape: no validator, no rework, rework never re-reviewed, re-review before rework, self-review (re-review reaching the fix producer), an extra implementation round, three-plus validation rounds. A bar that reds nothing is worse than the false red it replaces.
  Stands from cycle 1; spike 20/20 recorded in "Spike (riskiest mechanism, exercised)".
- DONE: State the routing rule that generalizes both branches without special-casing — codex took reuse with 10 events, claude took fresh with 8, and both must grade correctly under one rule.
  Stands from cycle 1; see the routing rule in "The grammar (normative)".
- DONE: Correct the entity's 9g claim in your report: only the TSV digest survives; the raw rollout that could establish intent still dies with the isolated CODEX_HOME. Recommend partial supersession, not retirement.
  Stands from cycle 1; see "9g: partial supersession, not retirement".

### Summary

Cycle 2 addresses the REVISE finding: the live run moved from a test-plan line into AC-6, whose verification IS the targeted local rejection-flow run (`SPACEDOCK_LIVE_RUNTIME=claude ... -run '^TestLiveCommonRejectionFlow$'`), naming the expected observed codes — `rejection-worker-topology` ABSENT, with red-attribution via diagnostic text plus retained digest — and naming the codex sandbox refusal (no isolated `CODEX_HOME`) as a stated unavailability whose reuse-branch gap AC-1's retained-chain replay covers deterministically and stack-tip CI corroborates. The test plan now distinguishes grader-changed (offline replay) from journey-passes (AC-6's run) as the README's Testing Resources section requires; no other cycle-1 content changed, and AC-6 adds no code surface.

## Stage Report: implementation

- DONE: The round-parsing grammar replaces the hardcoded four-rounds-by-two-events chain, and every diagnostic names the invariant it violated instead of a bare count
  `rejectionChain` and the flat position-compare in `assertRejectionWorkerTopology` are gone (1dbd72aa4); `parseRejectionRounds` folds the chain into rounds and `gradeRejectionRounds` grades sequence, budget, routing and independence. `TestRejectionTopologyRedsNonConformingShapes` asserts a distinct invariant phrase per case AND bans "routing events"/"owes 8" on EVERY red — reinstating the old length check fails both halves.
- DONE: The one-repeat budget widens the conforming set without accepting a genuinely non-conforming topology
  `TestRejectionTopologyGreensRetainedRepeatChains` greens both retained run-32270990171 chains (2/2, from 1/2) plus the same round spent after the rework; 13 red controls stay red — deleting the `repeats > rejectionRepeatBudget` guard greens the two two-repeat cases, and raising the budget to 2 greens them too.
- DONE: AC-6's targeted local live rejection-flow run against the built change, expecting rejection-worker-topology ABSENT, with the codex sandbox refusal NAMED rather than silent
  BOTH lanes ran live against the built binary and PASSED with `rejection-worker-topology` absent. Claude/sonnet, fresh branch, 689.13s, started 17:00:58: `/tmp/12z-ac6-evidence/ac6-live-artifacts2/claude-shared-scenarios/rejection-flow/rejection-topology.tsv`. Codex, reuse branch, 453.44s, started 16:19:29: `/tmp/12z-ac6-evidence/ac6-codex-artifacts/codex-shared-scenarios/rejection-flow/rejection-topology.tsv`. Run logs beside them as `ac6-live2.log` / `ac6-codex-live.log`; `SPACEDOCK_BIN` was the worktree-built `.worktrees/spacedock-ensign-rejection-topology-count-bar/spacedock` (0.27.0-pre8+dev). See the three findings below.

### Summary

The grammar landed as ideation settled it: rounds (a dispatch plus the completion closing it), the sequence `impl, val+, rework, val+`, a one-repeat TOTAL budget declared as a budget, one branch-keyed routing rule, and round-level reviewer independence. Fail-closed parsing turned two shapes ("re-review before rework", "two validators up front") into position-named unpairable reds. `reviewerCollapsedRoutes` now retargets both events of a round rather than only the dispatch: the half-mutated chain it used to build is a shape no extractor can emit, and the collapse this case is about is internally consistent, which is exactly why only the round-level identity check catches it.

FINDING 1 — the codex-sandbox unavailability in AC-6 is FALSE from a worktree root, and I did not inherit it. Probed first: `~/Library/Caches/spacedock-live-codex` is blocked (`Operation not permitted`), but the repo-adjacent candidate `.worktrees/.spacedock-live-codex/<worktree>` is writable, so the codex lane created its isolated `CODEX_HOME` and ran to completion. The ideation's "either candidate path" reading holds only from the MAIN repo root, whose parent `/Users/clkao/git/spacedock-research` is blocked. The consequence is good news: the REUSE branch — the branch that produced the falsifying 10-event chain — is now live-proven locally, not just replayed, so AC-6's stated coverage gap does not exist for this run.

FINDING 2 — SURFACE OVERRUN, disclosed. Actual net **+247 LOC** (311 insertions, 64 deletions) across **2 files**, against the declared net +105 ±50 — 235% of estimate, past the 2x line, and the file count is within tolerance. Causes, none of them a design change: (1) the grammar needs 13 distinct invariant-naming diagnostics, one per clause, each a 1-2 line `fmt.Errorf` — AC-2 forbids collapsing them into a shared message; (2) AC-1's minimum corpus added 7 red cases and 3 green cases to the table at ~6 lines each; (3) +72 of the +247 are comment lines at this file's prose density. No second harness, no mechanism beyond the round parser, no third file — so by this entity's own tolerance rule none of it requires a return to ideation. I refactored once to cut it (errors carry no `routes` param; the summary and graded code are attached once at the boundary; three named mutant helpers collapsed into `insertRound` call sites) which took it from +302 to +247. Further cutting would mean fewer diagnostics or fewer controls, i.e. trading AC-1/AC-2 for the number.

FINDING 3 — the claude lane took TWO attempts, and the first one is evidence FOR the deliverable. Attempt 1 (started 16:17:32, digest `/tmp/12z-ac6-evidence/ac6-live-artifacts/claude-shared-scenarios/rejection-flow/rejection-topology.tsv`) reddened with all six oracles red and `rejection-worker-topology` PRESENT: the journey stopped after three routing events. Two things follow. First, the red itself is AC-2 working on a REAL failure rather than a corpus mutant — it came back as "the spawn/validation dispatch at position 2 is never completed, so the run ends mid-round" with the route list attached, naming its invariant instead of the old bare "owes 8 routing events"; a red that reports honestly is the deliverable doing its job. Second, on cause I can only offer inference, and label it as such: my own session hit its quota limit at ~16:51, four minutes AFTER that run died, so the timing is suggestive but proves nothing. Supporting the same inference, also not proof — attempt 1's final message ("keeping the implementation worker alive since it's the feedback-to target. Waiting for the validation completion signal") is CONTRACT-CORRECT conduct: the FO dispatched validation and was waiting, which is what it is supposed to do, so the shape is consistent with the journey dying externally rather than the FO misbehaving. Honest statement: attempt 1 ended mid-journey for an undetermined reason and attempt 2 completed cleanly. No third run was spent, since the cause would not change AC-6's outcome.

Offline: `go test ./internal/ensigncycle -run RejectionTopology` green plain and `-race`. Tree-wide `go test ./...` and `-race` are green except `TestCodexResolveManifestAgainstInstalledHost` (internal/cli), a pre-existing machine-state failure — a stale `spacedock-local` 0.27.0-pre7 install under `~/.codex/plugins/cache` that `codex plugin list` does not report — in a package this diff does not touch. AC-5 holds: `TestRejectionTopologyExtractsCapturedCodexChain` and `rejectionTopologyDigest` are unmodified in the diff.

## Review-finding disposition

- **Reviewer (validation, 2026-08-19) — Deferred risk: three grammar clauses are individually deletable with the whole `RejectionTopology` suite still green.** Released user and normal workflow: a future edit to `rejectionRoundRouting` or `parseRejectionRounds`. Observable harm: deleting only the wrong-reuse-handle case, only the first-round-must-be-spawn case, or only the non-journey-stage guard leaves `go test ./internal/ensigncycle -run RejectionTopology` green — three live clauses no shipped test defends, so a regression in any of them ships silently. Affected authority: `none:` — no AC fails today. I drove all six corpus-unreached diagnostic paths on a throwaway checkout and every one reds as a graded `rejection-worker-topology` naming its invariant with no bare-count vocabulary, so AC-2's universal claim is TRUE; it is the proof, not the behavior, that is short. Trigger evidence: each single-clause deletion run on a throwaway worktree at `1dbd72aa4`, suite green all three times. The wrong-handle gap is inherited, not introduced — deleting the parent's `{0,4}/{2,6}` handle-correlation check at `8b8d009ea` is equally silent. Promotes to material if any of the three clauses is edited without a covering case. Narrow fix, test-only: add the six shapes (empty chain, `done` at an even position, non-journey stage, odd-length mid-round, journey opening on validation, reuse followup to an unopened handle) to `TestRejectionTopologyRedsNonConformingShapes`.
- **Reviewer (validation, 2026-08-19) — Polish: FINDING 1 names the wrong `CODEX_HOME` candidate as the enabler.** `codexLiveIsolatedHomeParentCandidates` tries three parents in order; candidate 1 is `{artifactRoot}/_codex-home` when the artifact root is outside `os.TempDir()` and outside the repo, and the repo-adjacent `.worktrees/.spacedock-live-codex/<worktree>` is candidate 3, reached only if 1 and 2 fail. The passing run used candidate 1: `codex-exec.jsonl` carries 14 `_codex-home` references and ZERO `spacedock-live-codex` references. The real enabler was `SPACEDOCK_LIVE_ARTIFACT_DIR` pointing at the session scratchpad, not the worktree-vs-main-root distinction the finding draws — from the main repo root the same artifact dir would also have worked. The repo-adjacent path IS writable (its empty per-worktree dir exists, created 16:19), so the probe's observation stands; only the causal attribution to this run is wrong. No AC is affected; the correction is one sentence of report prose.
- **Reviewer (validation, 2026-08-19) — Polish: "the REUSE branch ... is now live-proven locally, not just replayed" overstates what the codex lane proved.** The codex live digest is the conforming 8-event chain (`branch reuse`, no repeat). What ran live is the reuse BRANCH; the falsifying 10-event repeat SHAPE remains proven only by AC-1's deterministic replay. AC-6 asks only for the code's absence, which held, so nothing fails — but the sentence invites a reader to think the falsifying chain itself was re-run live.
- **Reviewer (validation, 2026-08-19) — Polish: two of FINDING 2's own numbers are off.** The grader has 17 distinct invariant-naming diagnostics, not 13. And cause 3's "at this file's prose density" does not hold: the added lines are 29% comment (89 of 311) against the two files' prior 11% and 20% and a 16% package-wide average; matching prior density would have saved roughly 42 lines. Neither changes FINDING 2's conclusion, and 17 understates rather than inflates the cause.

## Stage Report: validation

- DONE: Reproduce each of AC-1 through AC-6's cited evidence, including opening the three AC-6 digests and confirming the two PASS logs
  All six hold. AC-1: reproduced the baseline independently by running the 10-event chain through the PRE-change grader at `8b8d009ea` — RED with the exact bare-count message ("owes 8 routing events, the run produced 10") while the fresh chain greened, confirming 1/2 before and 2/2 after against a baseline that did move. AC-2: 13 shipped reds each name a distinct invariant and none carries the banned vocabulary. AC-3/AC-5: `TestRejectionTopologyBranchesAreMutuallyExclusive`, `TestRejectionTopologyProbeFailSafe`, `TestRejectionTopologyExtractsCapturedCodexChain` and `claudeRejectionStream` diff byte-identical against the parent, and `rejectionTopologyDigest` appears zero times in the diff. AC-4: both independence cases present and both flip green when only the independence check is removed. AC-6: opened all three digests — `ac6-live2.log` ends `--- PASS (689.13s)` over an 8-row `branch fresh` chain, `ac6-codex-live.log` ends `--- PASS (453.44s)` over an 8-row `branch reuse` chain, attempt 1's digest holds 3 rows. Re-graded both PASS digests from their RETAINED bytes through the shipped grammar rather than trusting the in-process grade: both GREEN, and attempt 1's reds as "the spawn/validation dispatch at position 2 is never completed" with no bare-count fallback. The built change was under test on both lanes — `0.27.0-pre8+dev` in both streams, and the codex stderr names the worktree install. No live run was re-executed.
- DONE: Confirm the grammar reds non-conforming shapes with distinct invariant-naming diagnostics and greens the retained repeat chains, by breaking the budget guard and watching the controls flip
  Seven single-clause mutations on a throwaway worktree, never the implementation worktree; every clause is load-bearing. Raising `rejectionRepeatBudget` to 2 greens exactly the two two-repeat controls. Reinstating the old bare-count length check fails NINE subtests — both retained-repeat greens plus seven reds whose invariant phrasing and banned-vocabulary assertions both fire — so AC-2's ban is falsifiable, not a spelling check. Dropping the independence check flips exactly the two self-review cases; dropping the routing rule flips the mutual-exclusivity test plus two probe fail-safe cases; dropping the post-rework clauses flips the missing-re-review and second-rework cases; loosening round pairing flips the two unpairable cases; dropping the `preRev == 0` clause flips the single-worker chain. I also probed the widening from both sides: one extra VALIDATION round greens at every legal insertion point and reds where the validation stage has not yet been spawned, while an extra IMPLEMENTATION round stays red at all four insertion points — the conforming set widens by exactly the one-repeat budget and nothing else.
- DONE: Weigh FINDING 2's 235% surface overrun against the entity's own return-to-ideation rule and state whether the causes hold
  Causes 1 and 2 hold; cause 3 is partially refuted; the conclusion survives. Diagnostics are load-bearing, not padding: the 17 map one-to-one onto ideation's normative grammar clauses, and the mutation battery shows every clause changes a verdict. Corpus growth is load-bearing: I checked the shipped cases against AC-1's own declared minimum and every named shape is present with none redundant — the two budget cases are required by ideation's cycle-1 total-vs-per-segment refinement. Cause 3 fails on its own terms: added lines run 29% comment against the files' prior 11%/20%. Decisively, net CODE alone excluding all comments and blanks is +163 (grader +78, table +85), which already exceeds the declared band's +155 ceiling — so the overrun is real growth, not a prose artifact, and it is roughly 195-205 net even after normalising prose. On the rule itself: the workflow makes the estimate the baseline a correction round calibrates against and requires a captain-visible design-reset decision past declared tolerance BEFORE another pass; it does not make an overrun an automatic return to ideation. I tested the three conditions that WOULD force one and all three are absent — one new mechanism (the round parser) exactly as declared, 2 files within the +1 tolerance, and the declared semantics changed exactly as declared. Tolerance is captain-owned, so this is ratification at the gate, not a defect.

### Summary

Recommending PASSED. All six acceptance criteria have valid, independently reproduced evidence: the AC-1 baseline was re-measured against the pre-change grader rather than inherited, both AC-6 lanes were re-graded from their retained digest bytes rather than from the in-process verdict, and the AC-3/AC-5 scope fence was diffed byte-for-byte against the parent. A seven-mutation battery on a throwaway checkout shows every grammar clause is load-bearing and that AC-2's bare-count ban genuinely fails when the old message returns. FINDING 2's surface overrun is real and past tolerance — larger than the report's own framing, since net code alone clears the declared ceiling before a single comment is counted — but it is a disclosed threshold deviation for the captain to ratify, not a design change: no second harness, no mechanism beyond the round parser, no third file. Four findings recorded, all Polish or deferred risk, none blocking; the deferred risk is a proof gap, not a behavior gap, since I drove all six corpus-unreached diagnostic paths and they behave as AC-2 claims. Tree-wide `go test ./...` on an unloaded run reds only `TestCodexResolveManifestAgainstInstalledHost`, which I reproduced byte-identically on main at `6cb8f8211` in a package this diff does not touch; an earlier loaded run added four codex-process reds that all pass in isolation in seconds and are load artifacts of a 250ms quiet budget, not defects. `-race` was run on `internal/ensigncycle` and is green; I did not run a tree-wide race pass.
