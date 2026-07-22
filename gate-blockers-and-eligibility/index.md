---
title: Approved-pending blockers, execution holds, and dispatch eligibility for recorded gates
status: validation
source: "Split from the gate-recorder task (3k), captain-approved 2026-07-21. Carries 3k's original seed concern (captain design feedback 2026-07-13: an approval evaporated while its dispatch stayed blocked)."
id: h1y616vjh64wc961z5t1031d
gates:
    version: 1
    current:
        gate: gate:docs-dev:h1:ideation
        attempt: gate-attempt:h1-ideation-1
    records:
        - id: gate:docs-dev:h1:ideation
          stage: ideation
          current-attempt: gate-attempt:h1-ideation-1
          attempts:
            - id: gate-attempt:h1-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:docs-dev:h1:ideation:briefing-1
                digest: sha256:f98f7ac3f9b6933a83ec8d573204c44ae1e2ba598f63378bc71ac09e604dbc78
                room-ref: "./review/ideation/briefing-1"
                note: "Multi-artifact package: gate summary, frozen entity snapshot, frozen recorder-contract snapshot — each digest-pinned inside the briefing; the digest above binds briefing.json itself."
              resolution:
                type: Resolution
                id: resolution:captain-chat-h1-ideation-1
                briefing: briefing:docs-dev:h1:ideation:briefing-1
                by: person:captain
                at: 2026-07-21T12:52:44Z
                decision: approve
                reason: "Captain approve in chat after a left-open float and a requested justification of the application abstraction (what breaks without it: exactly-once consumption across sessions, staleness marking, the approve-consequence routing, and the structural advisory-cannot-advance guarantee — each anchored to a lived incident). The left-open envelope is retained in-room (float-result.json); the captain saw the floated content and rendered this resolution in chat, so the captain identity is correct under the recording-identity ruling."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "The ships-vs-declines split stands: the application record ships; the blocker-evaluator half stays a recorded decline with its promotion condition."
sprint: durable-decisions
group: recorder
started: 2026-07-21T01:43:36Z
worktree: .worktrees/spacedock-ensign-gate-blockers-and-eligibility
---

The application layer is what a recorded gate decision *does*: it turns a durable approval into exactly one workflow action, or holds it durably when it must not act yet. The recorder owns the resolution record (what the decision *is*); this layer owns the one-use `application` — its `pending → consumed` transition through the existing workflow transition and dispatch path, its `superseded` marking on reviewed-input drift, its `not-applicable` state under a portable `hold`, and the fail-closed eligibility read that decides whether a pending action may act now. The record's field shapes are authored in the settled contract (`durable-gate-approval-pending-blockers/gate-resolution-frontmatter-contract.md`); this task designs the behavior over the `application.*` sections that contract marks application-layer-owned.

## Problem

A durable approval is only half a decision. In the 0260 shaping dry run the First Officer hand-recorded five resolutions (two approve) and the approved entities sat at their current stage with an `advance/pending` application, awaiting a later cold-booted session to apply them. That un-consumed pending application *is* durable "approved but not yet dispatched": it survives restart, and a cold reader reconstructs the entity's true position from frontmatter alone. What that dry run never mechanized is the act of consuming it safely — advancing and dispatching **exactly once**, never twice, and never after the reviewed input has drifted.

Exactly-once is not free here. The existing dispatch path has no transactional "this stage was already dispatched" guard: a stage advance writes only `status`, dispatch is a separate step that writes nothing back to the entity, and the only de-facto "in progress" marker is the `worktree` field consulted in the read path — which is empty at a non-worktree stage and empty in the window between the advance and the worktree stamp. So a repeated consumption pass can re-advance and re-dispatch, and nothing downstream catches the duplicate. The application layer's `pending → consumed` transition is therefore the *sole* exactly-once authority; it must close its own window rather than lean on a downstream guard that does not exist.

The seed concern (captain, 2026-07-13: "an approval evaporated while its dispatch stayed blocked") names two failures — the approval was not durable, and the "must not act yet" state had nowhere to live. The durable pending application answers the first fully. The second splits into two very different things: an approval the First Officer is deliberately not consuming yet (covered by pending + the First Officer's judgment), versus a **machine-enforced** refusal to dispatch until a declared dependency is satisfied. The first is proven and shipped here; the second is the blocker/eligibility-computation half whose live need this ideation re-examines.

## The honest asymmetry (what ships vs. what is declined)

The two halves of this task have unequal evidence, and the split explicitly requires judging them separately.

- **The application record has a demonstrated live consumer.** The 0260 Commander consumed pending `advance` applications through the normal transition path. Credit that as live evidence: the one-use lifecycle, staleness, and application-state eligibility surfacing all serve a consumer that already exists. **Build it.**
- **Blockers, execution holds, and the blocker-satisfaction eligibility computation do not.** All eight recorded approvals in the dry run carried zero declared blockers and no execution hold. No consumer has ever queried blocker satisfaction. Judged separately, this half has no demonstrated live consumer, so it lands as a **recorded decline** (below), not a build — while the eligibility read still honors any blocker or hold that is present, fail-closed.

## Proposed approach

**One write surface — extend the recorder binary, never a second gates writer.** The application-state transitions (`pending → consumed`, `→ superseded`, `→ not-applicable`) are new verbs on the same binary that owns every `gates:` write; the eligibility read is a pure function surfaced through status. This is the responsibility-boundary rule "one write surface — the recorder binary, extended by this layer's verbs." No new schema field is introduced: the layer mutates only the `application.*` scalars the recorder already round-trips, so the recorder's eight-entity replay stays green by construction.

**Exactly-once = one atomic frontmatter commit.** Consuming a pending `advance` writes the effect (`status: <target-stage>`) and the consumption mark (`application.state: consumed`) in the *same* frontmatter node-round-trip and the *same* atomic write. Either both land or neither does, so there is no crash window for a re-pass to observe `pending` after the advance already happened. Dispatch follows as the existing separate step; if it fails, the entity is already `status: advanced, application: consumed` and the existing dispatch retry path re-dispatches without the application being re-consumed. Re-reading `consumed` never re-authorizes.
  - *Simplest alternative considered:* two commits — advance, then mark consumed. *Insufficient:* with no downstream dispatch idempotency (see the spike), a crash between the two commits leaves `pending` atop an advanced `status`, and the next pass double-applies with nothing to catch it. The single commit is what makes exactly-once real.

**Eligibility is a pure, fail-closed read of the record — it never queries another entity.** A pending application is eligible only when it is the current gate/attempt/Briefing at the current stage, its reviewed input is unchanged, its state is `pending`, its decision/action/target agree, no execution-hold is active, and its `blockers` list is empty. Any present blocker, any active hold, any non-pending state, or any missing/ambiguous field reads **ineligible**. Because the layer ships no blocker-satisfaction evaluator, a present blocker is by definition unqueryable and therefore never satisfied — declining the evaluator is the strictly *more* fail-closed choice, not a gap.
  - *Simplest alternative considered:* a blocker-satisfaction evaluator that queries a dependency entity's stage/revision. *Insufficient / declined:* no demonstrated consumer, and the pure read already covers the proven consumer (the Commander needs current + not-stale + pending), while being strictly more fail-closed.

**Closure shapes the application; drift supersedes it.** A recorded `approve` closure yields `advance/pending`; a portable `hold` yields `action: none, state: not-applicable` (the reviewer did not approve the material — never a pending advance); a `revise` yields a pending feedback application whose routing is owned by the feedback route, not here. Any digest-bound reviewed-input change before consumption marks the pending application `superseded` and requires a new attempt; closed attempts never reopen.

## Recorded decline (blocker-satisfaction evaluator + execution-hold authoring)

This layer **declines to build**, as a recorded ideation outcome the split authorized:

1. a blocker-satisfaction **evaluator** — the code that would query a dependency entity's stage/revision to flip a declared blocker from `unsatisfied` to `satisfied`; and
2. a captain-facing **command to set or release** an execution-hold.

**Justification.** No demonstrated live consumer: all eight recorded approvals in the 0260 dry run carried zero declared blockers and no execution hold; cross-task sequencing in this sprint is handled by convention (the durable-decisions sequencing section), not by querying a per-entity blocker record; and a pending application the First Officer has simply not consumed already provides durable "approved but not dispatched" without a machine-enforced refusal. Declining is the more fail-closed choice: with no evaluator, any blocker present renders the application permanently ineligible (never auto-consumed), which the eligibility read enforces today.

**What still ships from this half:** the eligibility read *honors* any `blockers` entry or active `execution-hold` that is present — fail-closed, ineligible — even though nothing in this layer authors one. There is deliberately no authoring path for a hold or a blocker: until the promotion condition below fires, neither can be created. The promotion is the only route to authoring — a live consumer appearing promotes the hold/blocker half to its own captain-approved design round that ships the authoring surface with it; no hand-edit or ad-hoc write is implied in the meantime.

**Promotion condition.** The first time a real gate must durably carry "approved but blocked on X landing" across a session boundary *and* have a machine — not the First Officer's judgment — refuse dispatch until X satisfies a queryable predicate, the evaluator graduates to a binding gate attempt naming the queried entity and predicate. Until then the decline stands and the fail-closed honoring is the safe floor.

## Spike: riskiest unverified mechanism

The riskiest unverified seam is the exactly-once co-write, because the existing dispatch path provides no idempotency to fall back on. Spiked before designing the rest, mirroring the real frontmatter writer (`internal/status/mutate.go`'s `yaml.Node` parse → mutate → single `yaml.Marshal` → atomic write):

- **Co-write:** flipping the top-level `status` scalar (advance) and the nested `application.state` scalar (`pending → consumed`) in one marshal lands both and round-trips the frozen closure (resolution id, `decision: approve`, briefing digest) field-exact. PASS.
- **Exactly-once:** re-reading the now-`consumed` record reports ineligible — no re-authorization. PASS.
- **Fail-closed:** a record whose current application carries an unsatisfied blocker and an active hold reads ineligible. PASS.
- **Downstream idempotency (read-level finding):** confirmed by reading the code that Spacedock has no transactional dispatch guard (advance writes only `status`; `dispatch build` writes nothing back; the only marker is the `worktree` field in the read path, empty at non-worktree stages and in the advance→dispatch window). This is *why* the co-write must be a single commit.

The throwaway (scratchpad `h1spike/`) seeds the implementation's first test: the atomic-co-write round-trip and the eligibility table.

## Acceptance criteria

Entity-level properties of the finished layer. Each names how it is tested.

**AC-A1 (application durability + surfacing).** A recorded approval with a pending advance application survives process restart and reports the recorded application state — `advance/pending`, rendered as `approved-pending` — byte-for-byte, without advancing `status` or dispatching. *Test:* restart the fixture; status text and JSON show `advance/pending` and the `approved-pending` condition alongside the recorder's `approve` surfacing.

**AC-A2 (VALUE — consumed exactly once).** Consuming a pending advance application advances `status` exactly once and flips the application to `consumed` in a single atomic frontmatter commit; repeated consumption passes produce zero further advances and zero further dispatches. *Baseline that can move the wrong way:* the count of advance+dispatch effects the existing transition/dispatch fake records across three consecutive eligibility passes equals **exactly 1** (a re-apply regression makes it 2 or 3). *Test:* drive the consume verb against the fake, assert effect-count == 1 after three passes, and assert `status` and `application.state` changed in one write.

**AC-A3 (staleness / supersession).** Any change to the reviewed artifact or other digest-bound gate input before consumption marks the pending application `superseded` and produces zero advance/dispatch effects until a replacement Resolution is recorded on a new attempt; closed attempts never reopen. *Test:* mutate the reviewed digest, run the eligibility read → `superseded`, assert zero effects from a consumption pass.

**AC-A4 (second-pending-application refusal — RED FIXTURE).** A current gate attempt binds at most one application. Creating a second pending application while one already exists (pending or consumed) on that attempt is refused, not appended; the record still holds exactly one application. The invariant holds **across attempts, not just within one**: after a supersede opens a new attempt, the gate carries at most one pending application over all its attempts — the superseded attempt's application must read `superseded`, never left `pending`. This is a banked incident, not a hypothetical: the recorder's own ideation gate record briefly carried pending advances on both attempt 7 and attempt 8 when attempt 8 superseded attempt 7 but attempt 7's application state was left live; it was corrected `pending → superseded` at preflight and noted in attempt 7's application record. *Test:* red fixtures — (a) call the create-application verb twice on the same closed attempt; the second is rejected and the record is unchanged; (b) drive a supersede on an attempt that carries a pending application, then run the eligibility read across the whole gate and assert exactly one pending application remains (the new attempt's), the old one having moved to `superseded`.

**AC-A5 (closure shapes application).** A portable `hold` closure yields `action: none, state: not-applicable` and never a pending advance; a `revise` closure yields a pending feedback application; an `approve` closure yields `advance/pending`. *Test:* three closure fixtures assert the resulting application shape.

**AC-B1 (fail-closed eligibility).** The eligibility computation is a pure function of the record that queries no external entity. The exact current pending advance — current pointers, current stage, unchanged reviewed input, `state: pending`, agreeing decision/action/target, no active hold, empty blockers — is the only eligible case. Missing, ambiguous, or unqueryable blocker state never reads as satisfied and never consumes an approval; a stale, unknown, held, superseded, consumed, wrong-stage, wrong-decision, or any-blocker-present record reads ineligible. *Test:* a table test enumerating the one eligible case and each fail-closed case.

## Behavioral test plan

Proof by exercising the binary and observing on-disk state and the transition/dispatch fake — not prose. Estimated cost medium; fixture and CLI/table tests, one live-workflow smoke only for the consume path.

1. **Exactly-once against the fake (AC-A2).** Drive the consume verb against the existing transition/dispatch fake; assert exactly one advance+dispatch across three passes and a single-commit co-write of `status` + `application.state`. Seeded by the spike's co-write check.
2. **Application durability and surfacing (AC-A1).** Record an approval with a pending advance; restart; assert status text/JSON reports `advance/pending`/`approved-pending` byte-for-byte with no `status` change and no worker.
3. **Eligibility table (AC-B1).** One eligible row (current pending advance, empty blockers, no hold) and fail-closed rows: stale digest, superseded, consumed, active hold, unsatisfied/unknown blocker present, wrong stage, wrong decision, missing field. Only the first admits an action.
4. **Staleness (AC-A3).** Mutate the reviewed digest before consumption; assert `superseded`, zero effects, and that a new attempt is required (closed attempt does not reopen).
5. **Second-pending refusal (AC-A4, red fixture).** Two create-application calls on one closed attempt; assert the second is refused and the record holds one application.
6. **Closure shapes (AC-A5).** `approve`/`revise`/`hold` closure fixtures assert `advance/pending`, feedback/`pending`, and `none`/`not-applicable` respectively.
7. **Replay stays green.** Replay the recorder's eight production entities (they carry `application` blocks) through the extended binary; assert byte-identical round-trip — the layer added no schema field.

## Expected surface + tolerance

Go product code extending the gate recorder binary (never a second `gates:` writer): application-state mutation on the nested `application` node (`pending → consumed` co-written with `status`; `→ superseded`; `→ not-applicable`), the pure fail-closed eligibility function, application-state surfacing in `internal/status`, and 1-2 `spacedock gate ...` verb entries in `cmd/`. **~250-400 production LOC** (smaller than the recorder — it rides the recorder's model and parser, adds no schema), roughly equal test LOC (the eligibility table, the exactly-once fake, the red fixtures, the eight-entity replay). **Tolerance 2×.** Hard self-check — any of these trips a reconfirm: a schema change that breaks the eight-entity replay, a second `gates:` writer instead of extending the recorder, a blocker-satisfaction evaluator or execution-hold authoring command (declined), or any subspace-tui coupling (the presentation command's surface).

## Documentation change proposal

This layer changes user-visible status output (a recorded application now surfaces `approved-pending`, `consumed`, `superseded`, and `stale`). The recorder task's proposed edit to `docs/site/concepts/gates-and-decisions.md` already describes this layer's behavior — "Eligible work then advances exactly once; unresolved blockers or an explicit approve-but-do-not-dispatch hold keep it at the gate without losing your approval." Coordinate rather than duplicate: that wording is accurate for the shipped behavior, given that a present blocker keeps the work ineligible (fail-closed) even though this layer does not auto-clear it.

Add to `docs/site/reference/frontmatter-contract.md`, after the recorder's `gates` paragraph:

```diff
+A recorded approval carries a one-use application. Status surfaces it as `approved-pending` until it is consumed, `consumed` once the approval has advanced and dispatched exactly once through the ordinary transition path, `superseded` if the reviewed input changed before it was consumed, and `not-applicable` for a reviewer hold. A present dependency blocker or an active execution hold keeps the application ineligible; the application is never consumed twice.
```

## Out of scope

- A blocker-satisfaction evaluator or any external-entity dependency query (recorded decline above).
- A captain-facing command to set or release an execution hold (recorded decline above).
- A second dispatch identity, receipt, or crash-recovery state machine — existing transition and dispatch state owns effect identity and recovery.
- The rejected-gate → rework route context projection (deferred with the recorder's AC-9).
- Any `gates:` write surface separate from the recorder binary.

## Stage Report: ideation

- DONE: The application layer designed against the recorder's settled contract sections this task owns — one-use lifecycle (pending → consumed exactly-once via the existing transition path; superseded on drift; not-applicable on hold), staleness, eligibility surfacing — with the 0260 Commander credited as the live consumer and the second-pending-application refusal as the red fixture.
  AC-A1..A5 + B1 designed against the h1-owned `application.*` sections; consumer credited in Problem/asymmetry; AC-A4 is the second-pending red fixture.
- DONE: Blockers, execution holds, and the eligibility computation judged separately and closed as a recorded decline against fail-closed honoring.
  Declined the blocker-satisfaction evaluator and hold-authoring command (no demonstrated consumer; zero blockers across all eight dry-run approvals); eligibility read honors any present blocker/hold fail-closed; promotion condition recorded.
- DONE: Expected surface + tolerance declared; design EXTENDS the recorder binary; riskiest unverified mechanism spiked first; eight-entity replay stays green.
  ~250-400 LOC, tolerance 2×, no second gates writer, no schema field; spike (scratchpad `h1spike/`) proved the atomic status+application.state co-write, exactly-once re-read, and fail-closed eligibility — all PASS.

### Summary

Split the layer by evidence: the application record (one-use consume-exactly-once, staleness, application-state eligibility surfacing) ships because the 0260 Commander is a demonstrated live consumer; the blocker-satisfaction evaluator and execution-hold authoring are a recorded decline because no gate has ever carried a blocker and the FO sequences by convention. The riskiest seam — exactly-once with no downstream dispatch idempotency — is closed by co-writing `status` and `application.state: consumed` in one atomic frontmatter commit, spiked and proven against the real writer's node round-trip. The design extends the recorder binary (one write surface), adds no schema field, and keeps the recorder's eight-entity replay green.

## Stage Report: ideation (cycle 2 — preflight folds)

- DONE: Preflight fold, second material finding's criterion half — cross-attempt pending-application invariant.
  Extended AC-A4: after a supersede, the gate holds at most one pending application across ALL attempts (the superseded attempt's must read `superseded`); cited the banked incident (the recorder's own ideation gate briefly held pending advances on attempts 7 and 8, corrected at preflight) and added a supersede-then-eligibility red-fixture test.
- DONE: Preflight fold, third decline — fallback reworded.
  Removed the sentence naming a non-existent recorder gates-write verb; the recorded promotion condition is now the only authoring route (a live consumer promotes the hold/blocker half to its own captain-approved design round), no hand-edit fallback implied.

### Feedback Cycles

- Cycle 1: CHANGES_REQUESTED — Roborev branch_final job 545; surface 14 files/971 added LOC (543 production, 404 test, 24 docs) vs estimate 250–400 production plus roughly equal tests and small docs (136% of production upper estimate; within 2×); AC unchanged
- Cycle 2: CHANGES_REQUESTED — Roborev branch_final job 551; surface 14 files/1,077 added LOC (605 production, 448 test, 24 docs) vs estimate 250–400 production plus roughly equal tests and small docs (151% of production upper estimate; within 2×); AC unchanged
- Cycle 3: CHANGES_REQUESTED — Roborev branch_final job 557; surface 14 files/1,143 added LOC (631 production, 488 test, 24 docs) vs estimate 250–400 production plus roughly equal tests and small docs (158% of production upper estimate; within 2×); AC unchanged
- Write-scope incident: the implementation worker directly authored the three FO-owned cycle lines before requesting First-Officer review. The First Officer halted the worker at the Cycle-3 boundary, independently reconciled each job id, reviewed ref, surface total, and unchanged-AC claim, and adopts the lines here; no worker-authored state write is treated as self-authorizing.

## Stage Report: implementation (FO closeout after Cycle 3)

- DONE: Implement the approved one-use application lifecycle, atomic consume, staleness/supersession, closure shapes, and pure fail-closed eligibility without adding the declined blocker evaluator, hold authoring, Subspace coupling, or a second gates writer.
  Commits `b1513f80`, `6c41403a`, `2cba20e0`, and `c7612661` form the clean implementation branch. The final commit is retained only for its independently useful literal-byte evidence, requested-only eligibility materialization, split-root definition-root correction, and removal of the test-local dispatch counter.
- DONE: Close Roborev jobs 545, 551, and 557 at the Cycle-3 human boundary, with the two unresolved product findings decided by the captain rather than another automatic correction pass.
  Legacy/pilot application compatibility is DECLINED/DEFERRED under the settled unreleased strict-v1 ruling; promote only if a future captain explicitly promises compatibility to a released external consumer. Approval at a terminal or single-stage gate with no successor is DECLINED/DEFERRED because h1's approved application is a one-use advance to a successor and no supported consumer exists; promote only when a real supported commissioned workflow gates a stage with no successor.
- DONE: Preserve and route the unintended post-limit review as audit evidence without treating it as a fourth correction cycle.
  Job 563 was initiated before the Cycle-3 stop and is not authorization for more implementation. Its compatibility and terminal findings inherit the captain's declines. Its medium claim that the same-gate supersession test can false-green is routed as a mandatory fresh-validation attack: assert the prior application state directly and prove a mutant that leaves it pending fails.
- DONE: Reconcile final surface and verification responsibility.
  Exact `fa240a76...c7612661` surface is 14 files, 1,146 additions and 52 deletions: 642 production additions, 480 test additions, and 24 documentation additions. Production is 160.5% of the 400-line upper estimate and remains below the approved 2× ceiling. Focused checks, formatting, and `git diff --check` were green at `c7612661`; full normal/race suites were green before the final evidence/performance commit and are explicitly validation-owned at the exact final head.

### Summary

Implementation closes at clean `c7612661` under the Cycle-3 captain decision. The supported successor-gate application path remains the only promised path; pilot compatibility and successorless terminal approvals are explicit deferred semantics with concrete promotion conditions. No fourth implementation round is authorized, and job 563's test-hole claim must be attacked independently before any validation recommendation.

## Stage Report: validation

- DONE: Independently validate exact h1 candidate c7612661 and full range fa240a76...c7612661 against the approved application-layer design, Cycle-3 captain decisions, current strict-v1 recorder contract, and every AC; do not rely on implementation narration.
  `git rev-parse HEAD` returned `c7612661cce857e90ae2073ac861f5b8b32b72c0`; diff/code/contract inspection and fresh behavioral drives found no material defect.
- DONE: Reproduce AC-A1 application durability/surfacing across restart: a canonical approve resolution yields advance/pending and approved-pending without changing status or dispatching before consumption.
  Detached fresh-process `TestH1AuditApprovalSurvivesFreshProcessesWithoutAdvance` passed; it would fail on status/worktree mutation, non-identical restart output, or loss of `approve`, `advance/pending`, `approved-pending`, or eligibility bytes.
- DONE: Reproduce AC-A2 value behavior through the supported successor-gate path: consume atomically co-writes status plus application.state=consumed, produces exactly one authorization consumption across repeated passes/crash-boundary re-entry, and leaves ordinary dispatch retry separate.
  Final-head consume/crash tests and detached fresh-process `TestH1AuditFreshProcessConsumeSpendsAuthorizationOnce` passed; a second process returned `condition=consumed consumed=false` and byte-identical state, while dispatch receipts/recovery remained absent and separate.
- DONE: Reproduce AC-A3 staleness/supersession and AC-A4 second-pending refusal across same and successor attempts. Mandatory job-563 attack: assert the prior application state directly, then introduce a mutant that leaves it pending; the named same-gate supersession test or an independent control must fail, never false-green via a no-op string replacement.
  Stale/same-/cross-attempt and second-close tests passed; the job-563 mutant made the original same-gate replacement test false-green (exit 0) but failed the fresh direct-state control and existing cross-gate direct-state control (both exit 1).
- DONE: Reproduce AC-A5 closure shapes for approve/revise/hold and AC-B1's one eligible row plus every fail-closed row, including stale, unknown, held, superseded, consumed, wrong-stage, wrong-decision, missing field, and any blocker present.
  `TestRecordClosureShapesApplication` and `TestEligibilityFailClosedTable` passed with exact action/target/state and condition/boolean assertions; changing any row's condition or eligibility falsifies them.
- DONE: Exercise the two captain-declined paths explicitly and prove zero mutation: legacy/pilot non-strict application shapes are rejected without compatibility inference or migration; approval at a terminal/single-stage gate with no successor is rejected without creating a Resolution/application. Supported gates with a real successor must remain green.
  Detached strict/pilot and single-/terminal-stage controls passed byte-identical refusal checks; compatibility promotes only on a future captain promise to a released external consumer, terminal semantics only on a real supported successorless workflow.
- DONE: Verify workflow ownership/discovery and split-root behavior: implicit and explicit operations use the owning workflow definition; requested gate-condition/gate-eligible status projection is correct; unrequested status paths avoid the added eligibility I/O.
  Owning-definition CLI/status tests and detached split-root/requested-only control passed; wrong-definition-root and forced-unrequested-I/O mutants each failed that control (exit 1).
- DONE: Audit exact surface and boundaries: 14 files, 1,146 additions/52 deletions with 642 production, 480 test, and 24 docs additions; under 2x; no blocker evaluator, hold authoring, Subspace coupling, second gates writer, compatibility layer, terminal semantics, dispatch receipt, or crash-recovery state machine.
  `git diff --numstat fa240a76...c7612661` independently totals exactly those values (production 160.5% of the 400-line upper estimate); boundary inspection found none of the prohibited surfaces.
- DONE: Verify canonical byte preservation and eight supported application shapes using independent literal expectations; exercise atomic co-write byte preservation for unrelated frontmatter/body content. Review Roborev jobs 545/551/557 and post-limit 563 dispositions plus the Feedback Cycles write-scope incident.
  Literal eight-shape replay and detached atomic-strip comparison passed; `roborev show 545/551/557/563` matched the adopted dispositions, job 563 remained audit-only, and no implementation-worker state write was treated as self-authorization.
- DONE: Run a detached adversarial audit for gates/status high-stakes mutations, then gofmt -w ./cmd ./internal, focused gates/status/CLI tests, go test ./..., go test ./... -race, git diff --check, and exact-head cleanliness checks.
  All commands exited 0 at exact head; detached pending-supersession, forced-eligibility-I/O, and wrong-split-root mutants each produced the expected non-zero control, and final `git status --short` was empty.

### Summary

PASSED: every AC has fresh executable evidence and no material finding remains at `c7612661`. Deferred risks are the three captain declines: blocker/hold authoring promotes only when a real cross-session machine-enforced predicate consumer appears; strict-v1 compatibility and successorless approvals promote only under their conditions above. The original same-gate assertion is tautological in isolation, but independent direct-state and existing cross-gate controls both killed the job-563 mutant, so it does not block this gate.
