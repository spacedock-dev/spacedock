---
title: Make the headless recorded-gate stop fixture stage-coherent
status: validation
source: "PR #580 run 30591046287, Sonnet job 91033369022"
started: 2026-07-31T00:32:10Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-headless-recorded-gate-stop-stage-coherence
issue:
milestone: 0.27.0
id: 26nk8qd48zknqnn4kc123sez
gates:
    version: 1
    current:
        gate: gate:26nk8qd48zknqnn4kc123sez:ideation
    records:
        - id: gate:26nk8qd48zknqnn4kc123sez:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:26nk8qd48zknqnn4kc123sez-backlog-1
              briefing:
                id: briefing:26nk8qd48zknqnn4kc123sez:backlog:attempt-1:revision-1
                digest: sha256:fea869611abb6a21b3bdf569d264e8c7dbc6166b5869203beec12d8aec962afb
                digest-domain: canonical-bytes
                request-digest: sha256:c6dd2c6b17d18deb57e14686317e8a856fb17c96ae5f6072c601fd0beba9b649
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:26nk8qd48zknqnn4kc123sez:backlog:1
                briefing: briefing:26nk8qd48zknqnn4kc123sez:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-31T00:32:05.343329Z"
                decision: approve
                reason: 'Captain conn approves focused ideation because the required Sonnet lane exposes a no-authority breach from contradictory fixture state; the task must distinguish fixture ownership from a real withdrawal-contract defect before any PR #580 mutation.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:26nk8qd48zknqnn4kc123sez:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:26nk8qd48zknqnn4kc123sez-ideation-1
              briefing:
                id: briefing:26nk8qd48zknqnn4kc123sez:ideation:attempt-1:revision-1
                digest: sha256:9f7a7ba5dbde944edbd8684a2fd8e26f6305540a7161c2ea14c92258829cf3e8
                digest-domain: canonical-bytes
                request-digest: sha256:b7db54faa91d87e7ad21b7f5c00ad09a1b32223be38ac796f429f7c43426b544
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:26nk8qd48zknqnn4kc123sez:ideation:1
                briefing: briefing:26nk8qd48zknqnn4kc123sez:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-31T00:47:13.663192Z"
                decision: approve
                reason: 'Approved under sprint conn: the corrected ideation isolates a fixture-owned contradiction, preserves the no-authority boundary, limits changes to three test files, and requires one supported Sonnet live proof before validation.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
review-round:
    id: round:26nk8qd48zknqnn4kc123sez:implementation:10
    stage: implementation
    cycle: 10
    briefing:
        id: briefing:26nk8qd48zknqnn4kc123sez:implementation:round-10
        digest: sha256:7ece44670335a772324ed9622d7655664ae4d1679c4becd22e796e6c4bb757ec
        digest-domain: canonical-bytes
        room-ref: ./review/implementation/round-10
---

## Problem statement

PR #580 run 30591046287, Sonnet job 91033369022 failed
`TestLiveDefaultHeadlessStopsAtGate/default-headless-recorded-gate-stop`.
The fixture booted as dispatchable `implementation`, but its body still ended in
`## Stage Report: validation` and its selected entity snapshot said that validation
was complete. Its README alone gained an implementation stage definition.

The archived command/state timeline isolates the breach:

1. `status --boot` reported `current=implementation`, `next=validation`, and no
   ready gate, while `status --read` exposed the retained validation report.
2. Sonnet prepared and committed an implementation gate instead of dispatching
   implementation.
3. Both attempted `gate withdraw` calls exited 1. The first said implementation
   was not an actionable `gate:true` stage; the second said no logical validation
   gate existed. No withdrawal was recorded and `gate withdraw` did not change
   status.
4. Sonnet itself ran `status --set ... status=validation started`, committed it,
   then prepared and committed a second, validation gate.
5. `assertRecordedGateHoldLog` rejected the two successful prepares. No decision,
   consume, or successor dispatch occurred after the gate.

This is a fixture-owned contradiction, not evidence for a product authority
change. In the same job and product commit, the coherent `gate-guardrail` control
passed in 114.93s and the coherent withdrawn-attempt recovery passed in 101.74s;
each prepared and committed exactly once and stopped without decision, consume,
or successor dispatch. The contradictory default scenario alone failed in
329.06s.

## Proposed approach

Make the default scenario construct a distinct implementation-start fixture
instead of partially rewriting the validation-ready fixture:

- README: define `### implementation` with exact instructions to append one
  implementation Stage Report and return; retain validation as the next
  `gate:true` stage.
- Entity: start at `status: implementation`, with no selected gate and no
  validation Stage Report. The implementation worker creates the first report.
- Selected entity snapshot: remove the false claim that a validation Stage Report
  is already complete; describe the retained package without claiming completed
  stage work.
- After implementation completion, the normal FO transition enters validation.
  Only then may the run prepare and commit one validation Briefing and present it.

Keep `assertRecordedGateHoldLog` load-bearing. Strengthen its post-prepare boundary
to reject `gate withdraw` and `status --set` as well as decision, consume,
successor dispatch, duplicate prepare, and missing commit. Add a
default-scenario-only ordered command-log assertion requiring one implementation
dispatch before the first successful validation prepare. This uses the existing
provider-neutral logging shim; it does not inspect transcript, model, provider, or
shell grammar.

The fresh local Sonnet spike used the above entity change in a detached scratch
clone of PR head and passed
`TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle`. Its live launch was
blocked before any FO work by OAuth 429 (weekly limit), so it is not claimed as
live proof. Implementation must supply one fresh supported Sonnet live result.

## Alternatives considered

- Starting directly at validation is coherent and the same-job control passes,
  but it drops the default journey's value: proving headless dispatch from the
  initial implementation stage to the gate.
- Accepting the repair, duplicate prepare, or post-prepare status mutation weakens
  the no-authority boundary and is rejected.
- Adding a product guard around withdrawal or status transitions is unsupported:
  both withdrawals failed and the coherent product controls passed. If the fresh
  coherent Sonnet journey reproduces a repair, stop this task and commission that
  distinct product-contract correction.
- Transcript/model/provider parsing observes a symptom rather than durable
  authority and is out of scope.

## Expected surface and semantics

- `internal/ensigncycle/shared_fixtures_test.go`: implementation-start fixture
  helper plus an offline coherence test; about 35 insertions.
- `internal/ensigncycle/claude_live_runner_test.go:310-340`: select the coherent
  helper and apply the scenario-specific ordered log assertion; about 8
  insertions and 5 deletions.
- `internal/ensigncycle/livescenario_adapter_live_test.go:16-48`: preserve and
  strengthen the no-authority grader and mutants; about 15 insertions and 3
  deletions.

Tolerance: at most these 3 test files, 60 insertions, and 10 deletions. No
production Go, skill, schema, command reference, or site documentation changes.
Command grammar, stored formats, product authority, and product runtime behavior
must remain unchanged. The only changed observable is the live fixture journey
and its stricter test-only grading.

No documentation diff is required because no user-visible behavior changes.

## Acceptance criteria

**AC-1 — The live journey starts from one coherent workflow episode.** On disk,
the initial fixture is at `queued`, has no selected gate or completed stage
report, and defines implementation before gated validation. Verified by a
focused fixture test that boots through
`spacedock status --boot --identify --json` and requires exactly one
dispatchable row with `current=queued`,
`next=implementation`, and no ready gate.

**AC-2 — Headless no-authority remains load-bearing.** A supported Sonnet run
starts from the queued projection, dispatches implementation once, records its
Stage Report, enters validation, binds and commits exactly one validation
Briefing, presents it, and stops open. A successful status transition is allowed
only after implementation dispatch and before prepare; repair before dispatch or
after prepare is rejected, as are decision, consume, withdrawal, duplicate
prepare, or successor dispatch. Verified by the durable entity, gate room, state
Git history, and provider-neutral command log; the final attempt has no Resolution
or Application.

**AC-3 — Product semantics remain invariant.** The delivered diff is confined to
the three declared test files, adds no transcript/provider observer, and changes
no command grammar, stored format, product authority, or runtime implementation.
Verified by diff inspection plus `go test ./...`, `go test ./... -race`, and the
focused live-tag mutant tests.

## Test plan

1. Add the offline fixture-coherence test described by AC-1. Cost: small, fixture
   test only; it fails if status, report history, selected gate, or stage taxonomy
   drift apart.
2. Extend the command-log mutant table with post-prepare withdrawal and status
   repair, and add ordered initial-dispatch mutants. Cost: small, no model; each
   mutant names the authority crossing that must fail.
3. Run `go test -tags live -count=1 -run
   '^TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle$'
   ./internal/ensigncycle/`.
4. Run the required supported Sonnet proof:
   `SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 20m -run
   '^TestLiveDefaultHeadlessStopsAtGate/default-headless-recorded-gate-stop$'
   ./internal/ensigncycle/ -v`. Preserve the test output and artifact directory.
5. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and
   `go test ./... -race`.

## Stage Report: ideation

- DONE: Reproduce the Sonnet breach and isolate whether contradictory status/gate/Stage Report setup is the causal fixture defect.
  Run 30591046287 artifacts show implementation status plus validation report led to two prepares and manual status repair; same-job coherent controls passed.
- DONE: Spike one stage-coherent fixture and the existing no-authority assertions before proposing any product-contract change.
  Detached scratch fixture removed the stale report and the focused mutant test passed; the live launch reached no FO work because local OAuth returned 429.
- DONE: Define the smallest owned correction, AC evidence, exact file/LOC semantics, and required focused/live proof without transcript/provider parsing.
  The design limits work to three test files, provider-neutral durable/log controls, and one mandatory supported Sonnet proof.

### Summary

Ideation identifies a fixture-owned contradictory episode; the archived evidence
does not show a successful withdrawal or product-driven status change. The
correction makes the implementation-start fixture coherent, strengthens
post-prepare authority grading, and requires a fresh Sonnet live proof before
validation.

## Stage Report: implementation

- DONE: Make the default live fixture boot as one coherent implementation-start episode with no selected gate or completed validation report.
  Commit 854ffa62b; `TestImplementationStartRecordedGateFixtureIsStageCoherent` shells real `status --boot --identify --json` and fails on a wrong stage pair, ready gate, selected gate, stale report, or completed-stage snapshot claim.
- DONE: Preserve the strict no-authority grader in the fixture-only candidate.
  `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` fails on missing/late/duplicate implementation dispatch, missing commit, duplicate prepare, decision, consume, withdrawal, status repair, or successor dispatch.
- DONE: Keep all changes test-only and prove them with focused fixture/mutant tests, full/race/format.
  The committed diff is 52 insertions/10 deletions in the three authorized test files; focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed.
- DONE: Prove one implementation dispatch, one validation prepare/commit, then stop open in one fresh supported Sonnet live journey.
  PR #583 Runtime Live E2E run 30595653496, Sonnet job 91047460694 passed `TestLiveDefaultHeadlessStopsAtGate` in 303.92s against exact PR head candidate 854ffa62b; its durable entity and provider-neutral command-log graders require that ordered lifecycle and reject resolution, application, withdrawal, status repair, duplicate prepare, or successor dispatch.

### Summary

The fixture-only candidate coherently starts implementation and strengthens the
provider-neutral command-log boundary without changing product behavior. The
later exact-head Sonnet CI journey proves AC-2: implementation dispatched once,
validation prepared and committed once, and the selected attempt remained open
without a resolution, application, or successor dispatch.

### Feedback Cycles

- Cycle 1: REJECTED — validation/PR #583 Sonnet evidence; surface 3 test files/+52/-10 vs estimate 3 test files/+60/-10 (87% of insertion ceiling); AC unchanged. The supported run reached the open validation gate without authority crossing, but it skipped implementation and the grader credited a validation dispatch plus pre-prepare status repair as the required implementation dispatch. Route only this Material task-owned false green: require a successful implementation-stage dispatch before prepare and reject validation/status-repair substitution; do not add provider/transcript grammar or change product gate semantics.
- Cycle 2: REJECTED — validation/adversarial ordering trace; surface 3 test files/+59/-10 vs estimate 3 test files/+60/-10 (98% of insertion ceiling); AC unchanged. The Cycle-1 predicate rejects every successful pre-prepare status mutation, including the required post-implementation transition into validation, so the correct implementation→validation→open-gate lifecycle can never pass. Design reset: reaffirm AC-2 and the same three-test-file/provider-neutral boundary, reset the ceiling to +75/-12, and route only command ordering—reject status repair before implementation dispatch while permitting the normal validation transition after it.
- Cycle 3: REJECTED — validation/root-cause escalation/captain; surface 3 test files/+63/-10 vs reset ceiling 3 test files/+75/-12 (84% of insertion ceiling); AC revised, not narrowed. Captain ruling: the behavior is not Sonnet-specific and is not a product scheduler defect. Initial stages intentionally project their successor, so `current=implementation,next=validation` correctly dispatched validation; 26n's fixture and AC-2 contradicted that contract while its coherence test blessed the row. Re-scope the same three-test-file task to a `queued` initial stage whose successor is implementation, then require implementation report → validation transition → one open gate. Retain a line-local successful-status-set detector for the pre-dispatch and post-prepare authority boundaries, without transcript/provider parsing. Reset the ceiling to +90/-15; no product command, storage, authority, or runtime semantic changes.
- Cycle 4: REJECTED — validation/exact journey trace; surface 3 test files/+81/-10 vs reset ceiling 3 test files/+90/-15 (90% of insertion ceiling); AC unchanged. The corrected queued topology is coherent, but the grader rejects the normal successful queued→implementation status transition before implementation dispatch. FIX authorized only for that line-local false negative: allow the implementation transition, reject pre-dispatch validation repair, and preserve all post-prepare mutation rejection; no product, provider, transcript, or shell-parser expansion.
- Cycle 5: REJECTED — validation/evidence-retention preflight; surface 3 test files/+85/-10 vs reset ceiling 3 test files/+90/-15 (94% of insertion ceiling); AC unchanged. The exact queued journey and prior mutants pass, but the mandatory live runner deletes the durable entity, state Git, gate room, and command log after a green run. FIX authorized only to retain the exercised workflow tree beneath the existing scenario artifact directory and prove those existing artifacts remain gradeable; no generic artifact framework or transcript/provider substitution.
- Cycle 6: REJECTED — validation/failure-path retention; surface 3 test files/+90/-10 vs reset ceiling 3 test files/+90/-15 (100% of insertion ceiling); AC unchanged. One fresh Sonnet run failed after 399.79s, but the assertion fired before the success-only retention call and cleanup deleted the provider-neutral workflow evidence needed to classify it. FIX authorized only to register the existing retention action before model/assertion failure and prove the failing path preserves the same tree; remain at the insertion ceiling by refactoring/cutting or stop on exact overage. No rerun or transcript diagnosis before that proof.
- Cycle 7: REJECTED — validation/retained exact journey; surface 3 test files/+90/-11 vs reset ceiling 3 test files/+90/-15 (100% of insertion ceiling); AC unchanged. The failure-safe artifact from one Sonnet run proves a coherent implementation dispatch → validation dispatch → one prepare/commit → open unresolved attempt, but the oracle false-rejects every second dispatch and therefore mistakes the required validation reviewer for a forbidden successor. FIX authorized only to require the two ordered pre-prepare dispatches and reject substitution, duplicate/late dispatch, or any post-prepare successor; replace/cut the incorrect count logic to remain at the insertion ceiling.
- Cycle 8: REJECTED — validation/detached handoff audit; surface 3 test files/+90/-11 vs reset ceiling 3 test files/+90/-15 (100% of insertion ceiling); AC unchanged. Ordered dispatch checks pass, but detached missing/late-transition mutants show the oracle never requires the successful implementation→validation status transition between the two workers. FIX authorized only to require exactly that between-dispatch transition and reject missing, late, duplicate, or other repair while retaining all prior mutants; no live spend before the focused trace is green and no insertion overage.
- Cycle 9: REJECTED — validation/retention cleanup lifecycle; surface 3 test files/+90/-11 vs reset ceiling 3 test files/+90/-15 (100% of insertion ceiling); AC unchanged. The 21-mutant event table is clean and one Sonnet run reached the correct retained open-gate outcome, but success copied the workflow and registered cleanup copied it again over read-only Git objects, false-failing the lane. FIX authorized only for single-copy-per-outcome behavior: success retains once, failure cleanup retains once, and cleanup no-ops on an existing retained tree; no new artifact mechanism or insertion overage.

- Cycle 5: REJECTED — validation/exact queued journey trace; surface 3 test files/+85/-10 vs estimate 3 test files/+90/-15 (94%); AC unchanged

- Cycle 6: REJECTED — validation/evidence-retention preflight; surface 3 test files/+90/-10 vs estimate 3 test files/+90/-15 (100%); AC unchanged

- Cycle 7: REJECTED — validation/failure-path retention; surface 3 test files/+90/-11 vs estimate 3 test files/+90/-15 (100%); AC unchanged

- Cycle 8: REJECTED — validation/ordered stage dispatches; surface 3 test files/+90/-11 vs estimate 3 test files/+90/-15 (100%); AC unchanged

- Cycle 9: REJECTED — validation/status transition placement; surface 3 test files/+90/-11 vs estimate 3 test files/+90/-15 (100%); AC unchanged

- Cycle 10: REJECTED — validation/single-copy retention; surface 3 test files/+90/-11 vs estimate 3 test files/+90/-15 (100%); AC unchanged

## Stage Report: validation

- DONE: Verify AC-1 against the real fixture boot and on-disk starting state.
  `TestImplementationStartRecordedGateFixtureIsStageCoherent` passed from an isolated Go cache; changing current/next, adding a ready or selected gate, retaining the validation report, or restoring the completed-stage snapshot makes it fail.
- FAILED: Verify AC-2 from the supported Sonnet journey's durable and provider-neutral evidence.
  PR #583 run 30595653496/job 91047460694 passed at 303.92s, but archived command-log lines 15-22 show `status ... --set ... status=validation started` followed by `dispatch build ... --stage validation`; the final entity has only a validation report, so implementation was never dispatched.
- DONE: Verify the exact candidate bytes and AC-3's declared semantic surface.
  Candidate `854ffa62b` and CI's synthetic merge commit `6cbb0981f` share tree `c1d7f9235`; diff inspection is exactly 52 insertions/10 deletions in the three authorized test files, with no production, grammar, format, authority, or runtime implementation change.
- DONE: Reproduce the focused live-tag grader and perform the semantic adversarial pass.
  The focused fixture/mutant tests pass, but replacing the accepted pre-prepare dispatch stage with validation still satisfies `assertRecordedGateHoldLog`, proving the grader can pass while AC-2's observable lifecycle is wrong.
- DONE: Classify the finding by trigger, harm, value boundary, and evidence.
  Released workflow: supported PR #583 Sonnet lane; observable harm: false green and false implementation-dispatch claim; `value-ac[AC-2]:` requires implementation → validation → one open gate; trigger evidence: archived command-log lines 15-22 plus the validation-only final entity.
- DONE: Record First Officer disposition and the minimum owned remedy.
  Material; outcome and evidence defect; task-owned; FIX authorized for implementation, requiring a successful `dispatch build --stage implementation` before validation and rejecting any earlier status repair/validation dispatch that substitutes for it.
- SKIPPED: Re-run full and race suites after the material finding.
  The review-finding checkpoint keeps candidate bytes unchanged and forbids a reviewer rerun after the finding; prior implementation evidence and the exact-tree offline CI job are green, but they cannot cure AC-2.
- DONE: Produce a PASSED/REJECTED recommendation with deferred risks separated.
  REJECTED on the single Material AC-2 finding; no separate deferred-risk or polish findings.

### Summary

REJECTED. The fixture starts coherently, but the supported Sonnet proof skipped
implementation and dispatched validation after a status repair; the new grader
accepted that wrong lifecycle. The candidate remains exactly `854ffa62b`, and
the minimum correction is confined to making the ordered lifecycle proof require
implementation dispatch before any transition that could substitute validation.

## Stage Report: implementation (cycle 2)

- DONE: Add a red-first provider-neutral command-log case proving validation dispatch plus pre-prepare status repair cannot satisfy the required implementation dispatch.
  Before the fix, focused subtest `validation_substitution` failed because the exact archived substitution was accepted; commit ccd0f367b retains that mutant as `substitution`.
- DONE: Require exactly one successful `dispatch build --stage implementation` before the first successful validation prepare while preserving the open no-authority hold and all existing negative controls.
  `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` now rejects validation-stage dispatch, pre-prepare status repair, their combined substitution, and every prior mutant while accepting the coherent lifecycle.
- DONE: Keep the correction within the three declared test files and estimate tolerance; run focused, full, race, format, and diff gates without a live or Roborev rerun.
  Incremental surface is one test file, +8/-1; cumulative surface is three test files, +59/-10; focused, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and diff checks passed, with live and Roborev intentionally not rerun.

### Summary

Cycle 2 fixes the task-owned false green in the provider-neutral command-log
grader without changing product semantics or expanding the authorized surface.
The deterministic mutants now make implementation-stage dispatch and the absence
of pre-prepare status repair load-bearing.

## Stage Report: validation (cycle 2)

- DONE: Re-review exact candidate ccd0f367b and prove validation dispatch plus pre-prepare status repair no longer substitutes for the required implementation dispatch.
  The focused live-tag grader passed; its validation-swap, early-repair, and combined-substitution mutants all fail when either substitution is accepted.
- FAILED: Verify the coherent implementation→validation→open-gate lifecycle and every prior no-authority/stage-coherence mutant still behave as declared.
  Prior mutants and fixture coherence pass, but a detached trace with implementation dispatch, then normal `status ... --set ... status=validation started`, then validation prepare/commit/open hold is rejected by the new any-pre-prepare-status predicate.
- DONE: Re-cross-check all acceptance criteria, the cumulative three-test-file +59/-10 surface, and focused/full/race/format/diff evidence; recommend PASSED or REJECTED with candidate bytes frozen.
  AC-1 and AC-3 remain evidenced; cumulative diff is exactly three test files, +59/-10, focused and diff checks pass, and implementation reports full/race/format green, but AC-2 is not validly provable, so the recommendation is REJECTED at unchanged head `ccd0f367b`.
- DONE: Classify the cycle-2 finding and retain the First Officer disposition.
  Released workflow: supported headless implementation completion; harm: the correct lifecycle is false-red and cannot prove AC-2; `value-ac[AC-2]:` requires implementation → validation → open gate; trigger: the detached coherent trace fails solely on the successful post-dispatch status transition. Material; task-owned evidence defect; FIX authorized.
- DONE: Bound the minimum remedy to the existing provider-neutral ordered evidence.
  Reject a successful validation/status mutation only when it substitutes before the required implementation dispatch; permit the normal post-dispatch transition before validation prepare, with no product, transcript/provider, runtime, or compatibility expansion.
- SKIPPED: Re-run full, race, live-model, or Roborev suites after the Material finding.
  The review-finding checkpoint freezes candidate bytes and stops reviewer reruns; broader green evidence cannot cure the focused coherent-trace false negative.

### Summary

REJECTED. Cycle 2 rejects the original validation-dispatch substitution, but it
also rejects the intended implementation-complete transition into validation.
The correction remains narrowly task-owned: order the successful status mutation
relative to the required implementation dispatch instead of banning every
pre-prepare status mutation.

## Stage Report: implementation (cycle 3)

- DONE: Add a red-first coherent trace proving implementation dispatch followed by the normal validation status transition and validation prepare/commit/open hold must pass.
  Before the fix, the focused test rejected the new implementation-dispatch → validation-transition → prepare/commit trace; commit d5b821f56 retains it as a positive control.
- DONE: Reject validation/status repair only when it occurs before the required implementation dispatch; retain the Cycle-1 substitution and every earlier mutant as red.
  `assertRecordedGateHoldLog` now scans successful status mutation only before the implementation dispatch; the coherent trace passes while validation swap, early repair, combined substitution, and all prior mutants remain rejected.
- DONE: Stay within the reset three-test-file +75/-12 ceiling; run focused, full, race, format, and diff checks without live or Roborev reruns.
  Incremental surface is one test file, +5/-1; cumulative surface is three test files, +63/-10; focused, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and diff checks passed, with live and Roborev intentionally not rerun.

### Summary

Cycle 3 makes the command-log proof order-sensitive without weakening the
implementation-stage requirement or the open no-authority hold. The correct
post-implementation validation transition is accepted, while a status or
validation substitution before implementation dispatch remains rejected.

## Stage Report: validation (cycle 3)

- DONE: Re-review exact candidate d5b821f56 and prove the coherent implementation-dispatch → validation-transition → prepare/commit/open-hold trace passes.
  The focused live-tag grader passed and its positive control exercises implementation dispatch, a successful post-dispatch validation transition, then one prepare/commit/open hold.
- FAILED: Verify early status repair, validation-dispatch substitution, their combined form, and every prior no-authority/stage-coherence mutant remain red.
  Existing mutants remain red, but a detached exact-flag-order mutant shows successful post-prepare `status --workflow-dir … --set` is accepted; a paired coherent trace with a failed early set is also false-rejected by cross-line substring matching.
- DONE: Cross-check all acceptance criteria, cumulative three-test-file +63/-10 surface, and focused/full/race/format/diff evidence; recommend PASSED or REJECTED with candidate frozen.
  AC-1 and AC-3 remain evidenced, the cumulative diff is exactly three test files +63/-10 under the +75/-12 reset, and focused/diff pass with implementation-reported full/race/format green; AC-2's no-authority proof remains invalid, so the frozen candidate is REJECTED.
- DONE: Perform a semantic adversarial matrix over status identity, order, flag placement, and exit result.
  Coherent post-dispatch transition passes; successful early repair and validation substitution fail; failed early repair plus an unrelated successful status read false-rejects; successful post-prepare repair with the supported flag order false-passes.
- DONE: Classify the finding and retain the First Officer's Cycle-3 disposition.
  Released workflow: the supported PR #583 command shape and provider-neutral log; harm: false green after a post-prepare authority crossing plus false red on a harmless failed attempt; `value-ac[AC-2]:` requires implementation → validation → open hold with no post-prepare repair; trigger: two failing detached traces. Material; task-owned evidence defect; Needs decision at Cycle 3, candidate held unchanged.
- DONE: State the narrow captain decision recommended by the First Officer.
  Recommend FIX with one line-local successful-status-set detector reused for pre-dispatch and post-prepare windows; allow the coherent post-dispatch transition and reject successful repair outside it without transcript/provider or product expansion.
- SKIPPED: Re-run full, race, live-model, or Roborev suites after the Material finding.
  The review-finding checkpoint and Cycle-3 hold freeze candidate bytes pending captain ruling; broader reruns cannot cure the focused evidence defect.

### Summary

REJECTED and held for captain decision. Cycle 3 fixes ordering for its positive
trace, but status mutation detection still combines different log lines and
assumes one flag order, allowing a real post-prepare authority crossing to pass.
The recommended correction is confined to a line-local successful-status-set
predicate over the existing provider-neutral command log.

## Stage Report: implementation (cycle 4)

- DONE: Make the fixture topology mechanically coherent: one initial queued stage projects implementation, which is dispatched and reported before validation can be entered.
  Commit 623d53e6f; `TestQueuedImplementationRecordedGateFixtureIsStageCoherent` fails unless boot returns `current=queued,next=implementation` with no ready/selected gate or prior report, and the live runner fails if validation is reached without an implementation report.
- DONE: Make the provider-neutral authority grader line-local and ordering-correct: accept the normal post-implementation validation transition, reject successful status repair before implementation dispatch or after gate preparation, and retain every prior mutant.
  `successfulStatusSet` inspects one successful command-log line at a time; focused positives cover the normal transition and failed early set, while early repair, supported flag-order post-prepare repair, validation substitution, and every prior mutant remain rejected.
- DONE: Revise AC-1/AC-2 mechanism wording without narrowing their value; stay within three test files and cumulative +90/-15, then pass focused/full/race/format/diff checks without live or Roborev reruns.
  AC-1/AC-2 now specify queued → implementation → gated validation; incremental code is three files, +33/-15 and cumulative code is three files, +81/-10; focused, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and diff checks passed, with live and Roborev intentionally not rerun.

### Summary

Cycle 4 aligns the test fixture with the existing initial-stage successor
contract while preserving product scheduling unchanged. The provider-neutral
grader now distinguishes successful status mutations line by line and enforces
their allowed window around implementation dispatch and gate preparation.

## Stage Report: validation (cycle 4)

- FAILED: Reconstruct the corrected queued -> implementation -> gated-validation journey and independently verify every prior authority-crossing mutant still fails.
  Fixture coherence and prior mutants pass, but a detached exact trace with successful queued→implementation status transition before `dispatch build --stage implementation`, then validation transition and prepare/commit/open hold is rejected as a forbidden pre-dispatch mutation.
- SKIPPED: Run one fresh supported Sonnet live journey and preserve its durable state, gate-room, command-log, and open-unresolved gate evidence; do not substitute aggregate or cached evidence.
  The required cheap spot-check proved the exact supported journey cannot satisfy the grader, so validation stopped before spending a live run and claims no substitute evidence.
- DONE: Cross-check AC-1 through AC-3, the three-file +81/-10 surface, focused/full/race/format evidence, and recommend PASSED or REJECTED with the candidate frozen on any new finding.
  AC-1 and AC-3 remain evidenced; the cumulative diff is exactly three test files +81/-10 under the +90/-15 reset, focused/diff pass, and implementation reports full/race/format green, but AC-2 fails at frozen `623d53e6f`, so the recommendation is REJECTED.
- DONE: Classify the finding and retain the First Officer disposition.
  Released workflow: Captain-approved queued initial-stage successor projection; harm: the correct mandatory Sonnet journey is false-red; `value-ac[AC-2]:` queued must project and dispatch implementation before gated validation; trigger: the detached exact command-log trace fails. Material; task-owned evidence defect; FIX authorized.
- DONE: Bound the authorized correction to the exact provider-neutral ordering defect.
  Permit the line-local successful queued→implementation transition before implementation dispatch, reject successful pre-dispatch transition/repair to validation, and retain all post-prepare mutation rejection; no transcript/provider parser, product semantics, or unrelated cleanup.

### Summary

REJECTED. The queued topology is mechanically correct, but the authority grader
still forbids the normal successor transition that makes implementation
dispatchable. The required live run was intentionally not spent after this
focused false negative; candidate `623d53e6f` remains unchanged.

## Stage Report: implementation (cycle 5)

- DONE: Add the exact queued-status-transition → implementation-dispatch positive trace red-first and make it pass without weakening existing negative controls.
  Before the correction, the focused test rejected the supported successful `status=implementation started` transition; commit 4aa59dafe retains that trace as a positive control while early validation repair, validation substitution, combined substitution, post-prepare status mutation, and every prior mutant remain rejected.
- DONE: Restrict the line-local pre-dispatch predicate to allow only the normal successful transition to implementation.
  `successfulStatusSet` now accepts an explicit allowed token pair only on the same successful status-set line; the implementation-start call site supplies that allowance, while all other call sites continue rejecting every successful status mutation.
- DONE: Stay within the authorized surface and run focused, full, race, format, and diff checks without a live or Roborev rerun.
  Incremental code is one test file, +7/-3; cumulative code is three test files, +85/-10 under the +90/-15 ceiling. The focused live-tag grader, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` passed.
- DONE: Record the Material finding and worker FIX disposition through the binary-owned advisory-round surface.
  Implementation round 5 validates as `all-fixed` against immutable candidate manifest `candidate-4aa59dafe.txt`; its reviewer evidence binds the released workflow, observable harm, `value-ac[AC-2]:` boundary, and exact frozen-candidate trigger, and its worker resolution records one material finding fixed with none declined.

### Summary

Cycle 5 fixes only the queued-to-implementation false negative in the
provider-neutral journey trace. The supported transition now passes, every
pre-dispatch validation repair and post-prepare mutation remains rejected, and
the cumulative candidate stays within the authorized three-test-file ceiling.

## Stage Report: validation (cycle 5)

- DONE: Reconstruct queued -> implementation -> gated validation and verify the prior authority and stage-order mutants against exact candidate `4aa59dafe`.
  The focused live-tag fixture/coherence and command-log tests pass. The retained positives cover queued-to-implementation and implementation-to-validation transitions; validation substitution, early validation repair, combined substitution, post-prepare mutation, decision, consume, withdrawal, duplicate prepare, missing or late implementation dispatch, and successor dispatch remain rejected.
- SKIPPED: Spend the fresh supported Sonnet journey.
  The required preflight found that the current live runner cannot preserve the admissible proof after a green run, so validation stopped before model execution and claims no aggregate, cached, transcript, provider, or shell-derived substitute.
- FAILED: Preserve the fresh journey's durable state, gate room, state Git history, provider-neutral command log, and open unresolved gate evidence.
  `runClaudeGateGuardrailScenario` creates the exercised workflow and command log under `t.TempDir()`, but `claudeLiveRunner.run` copies only `claude-stream.jsonl` and `claude-final-message.txt` into the persistent scenario artifact directory. Test cleanup therefore deletes the required durable evidence. The archived PR #583 Sonnet artifact confirms that the default scenario directory contains only those two files.
- DONE: Cross-check AC-1 through AC-3, the authorized surface, and the implementation verification record.
  AC-1's queued topology and AC-3's three-test-file `+85/-10` surface remain evidenced; focused and diff checks pass, and the implementation report records full, race, and format checks green. AC-2 cannot be validly proven because its required durable evidence is not retained, so the frozen candidate is REJECTED.
- DONE: Classify the finding and retain the First Officer disposition.
  Released workflow: mandatory supported Sonnet default-headless gate-stop proof; harm: a green run destroys the provider-neutral evidence and forces forbidden transcript evidence or an unverifiable aggregate; `value-ac[AC-2]:` requires the durable entity, gate room, state Git history, and command log; trigger: the current artifact writer plus the prior archived artifact's two-file contents. Material; task-owned outcome/evidence defect; FIX authorized.
- DONE: Bound the correction without mutating validation's frozen candidate.
  Preserve `4aa59dafe`; after the gate assertions, copy the exercised workflow tree beneath the existing scenario artifact directory and prove the retained entity/state Git, gate room, and command log remain gradeable. No generic artifact framework, transcript/provider observer, product change, or unrelated cleanup; remain within three files and `+90/-15`.

### Summary

REJECTED. The corrected command-log grader now accepts the exact queued journey
and rejects the prior mutants, but the mandatory live lane discards AC-2's
durable proof on success. Sonnet was intentionally not spent after that preflight
failure, and candidate `4aa59dafe` remains unchanged.

## Stage Report: implementation (cycle 6)

- DONE: Retain the exercised gate workflow tree beneath the existing scenario artifact directory after successful assertions, using the existing copy mechanism rather than a new artifact framework.
  Commit 1df0092ee calls the existing `copyTree` only after the source entity and command-log assertions pass, preserving the complete exercised workflow under the scenario's `workflow/` artifact subtree.
- DONE: Add a focused proof that the retained entity/state Git, gate room, command log, and unresolved attempt are present and still gradeable; keep the exact queued journey and all prior mutants green/red as intended.
  `TestRetainRecordedGateWorkflowEvidence` failed red on the absent retention seam, then passed by reopening the copied state Git, resolving its gate room, re-grading its unresolved entity, and re-running its copied log oracle; the queued/coherent positives and every prior mutation control also pass.
- DONE: Remain within three files and cumulative +90/-15, run focused/full/race/format/diff checks, and record the advisory disposition before returning to the same validator.
  The final cumulative diff is exactly three test files, +90/-10; focused retention/journey and fixture tests, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` passed, and implementation round 6 validates `all-fixed` with one Material FIX and no declines.
- SKIPPED: Spend a supported Sonnet or Roborev run during implementation.
  The assignment requires deterministic retention proof before model spend and routes the committed candidate back to the same validator for the fresh supported journey.

### Summary

Cycle 6 preserves the already-graded headless gate workflow as durable scenario
evidence and proves the copied state remains independently gradeable. It adds no
artifact framework or product/runtime semantics and lands exactly on the
captain-approved insertion ceiling.

## Stage Report: validation (cycle 6)

- DONE: Reconstruct the corrected queued -> implementation -> gated-validation journey and independently verify every prior authority-crossing mutant still fails.
  Exact candidate `1df0092ee`; the focused queued fixture, command-log positive/mutant matrix, and retained-workflow regrade passed in 2.6s, and each would fail on the previously identified topology, ordering, authority, or missing-evidence mutations.
- FAILED: Run one fresh supported Sonnet live journey and preserve its durable state, gate-room, command-log, and open-unresolved gate evidence; do not substitute aggregate or cached evidence.
  The one authorized Sonnet run failed after 399.79s with the provider-neutral grader's generic no-authority-boundary result, but the assertion fired before retention; the persistent scenario artifact contains only stream/final files and no `workflow/` subtree. Validation did not inspect the forbidden transcript and will not rerun.
- DONE: Cross-check AC-1 through AC-3, the three-file +81/-10 surface, focused/full/race/format evidence, and recommend PASSED or REJECTED with the candidate frozen on any new finding.
  AC-1 remains green; AC-3 is exactly three test files and actual cumulative `+90/-10` at the captain-reset `+90/-15` ceiling, with implementation-recorded full/race/format and validation-focused/diff checks green. AC-2 cannot be established because its observed failing journey discarded the required durable evidence, so the recommendation is REJECTED at frozen `1df0092ee`.
- DONE: Classify the new finding and retain the First Officer disposition before candidate mutation or reviewer rerun.
  Released workflow: mandatory supported Sonnet gate-stop proof; harm: a failing assertion destroys the only admissible state/log evidence and prevents diagnosis without forbidden transcript substitution; `value-ac[AC-2]:` requires durable entity, gate room, state Git history, and command log; trigger: the exact 399.79s run plus absent artifact `workflow/`. Material; task-owned evidence defect; FIX authorized.
- DONE: Bound the authorized correction at the existing insertion ceiling.
  Register the existing retention mechanism before the model/assertion path so the same provider-neutral workflow survives success, assertion failure, and cleanup, with the cheapest focused failure-path proof. No second Sonnet now, generic framework, provider/transcript diagnosis, product change, or silent overage beyond `+90/-15`.

### Summary

REJECTED. The deterministic journey, mutants, and success-path retention proof
pass, but the single fresh Sonnet run exposed that assertion failures still
delete AC-2's durable evidence before it can be retained. Candidate `1df0092ee`
remains unchanged and the exact failed run is preserved without transcript use.

## Stage Report: implementation (cycle 7)

- DONE: Make the existing workflow-evidence retention failure-safe by registering it before the model/assertion path, so success, assertion failure, and cleanup retain the same provider-neutral tree.
  Commit 0691f25a9 registers the existing `copyTree` cleanup immediately after creating the known scenario artifact directory and before model launch; cleanup ordering preserves `workflow/` before the earlier `t.TempDir` removal, while the existing success-path regrade remains in place.
- DONE: Add the cheapest focused failing-path proof and preserve the existing retained-tree regrade, queued journey, and all authority-crossing mutants without introducing another artifact mechanism.
  `TestRetainRecordedGateWorkflowEvidence` failed red on the absent registrar, then simulated a failed reporter, ran registered cleanup, and used the retained tree as the source of a fresh complete regrade; the queued/coherent positives and every prior mutation control remain green/red.
- DONE: Stay at or below the existing three-file +90/-15 ceiling by relocating/refactoring/cutting, run focused/full/race/format/diff checks, and record the advisory disposition; stop with the exact minimum overage if clean compliance is impossible.
  Previously added local scaffolding was compacted without dropping assertions; the cumulative diff is three files, +90/-11. Focused failure/retained-tree, journey, and fixture tests, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` passed, and implementation round 7 validates `all-fixed`.
- SKIPPED: Spend a second supported Sonnet or inspect the failed run transcript.
  The authorization requires deterministic failure-path proof before another live run and forbids transcript/provider substitution; the same validator receives the committed provider-neutral candidate.

### Summary

Cycle 7 moves the existing workflow copy onto a pre-launch cleanup boundary, so
the exact durable evidence survives both successful and failing live assertions.
The correction remains test-only, uses no second artifact mechanism, and stays
within the captain-approved cumulative ceiling.

## Stage Report: validation (cycle 7)

- DONE: Reconstruct the corrected queued -> implementation -> gated-validation journey and independently verify every prior authority-crossing mutant still fails.
  Exact candidate `0691f25a9`; the focused queued fixture, prior command-log positive/mutant matrix, success-path regrade, and simulated failing-path cleanup/regrade passed in 4.1s and remain falsifiable by the previously retained topology, ordering, authority, and evidence-loss mutations.
- FAILED: Run one fresh supported Sonnet live journey and preserve its durable state, gate-room, command-log, and open-unresolved gate evidence; do not substitute aggregate or cached evidence.
  Exactly one fresh Sonnet run executed for 348.20s and failed the command-log oracle, while failure-safe cleanup preserved `workflow/`; inspection used only that subtree and found ordered state Git, both stage reports, one canonical open attempt, and the full provider-neutral log. No transcript/provider artifact was read and no rerun occurred.
- DONE: Cross-check AC-1 through AC-3, the three-file +81/-10 surface, focused/full/race/format evidence, and recommend PASSED or REJECTED with the candidate frozen on any new finding.
  AC-1 is green and AC-3 is three files at actual cumulative `+90/-11` within the captain-reset `+90/-15`, with focused/diff green and implementation-recorded full/race/format green. AC-2's supported journey is false-red because the oracle counts all successful dispatches as one, so the recommendation is REJECTED at frozen `0691f25a9`.
- DONE: Classify the finding from retained provider-neutral evidence and retain the First Officer disposition before mutation or rerun.
  Released workflow: queued implementation worker followed by gated validation worker; harm: the valid open-gate journey is rejected; `value-ac[AC-2]:` requires implementation then validation and forbids only a post-prepare successor; trigger: the retained log's ordered implementation and validation dispatches plus one prepare/commit and no later authority action. Material; task-owned evidence defect; FIX authorized.
- DONE: Bound the authorized ordered-dispatch correction.
  Require exactly one implementation dispatch followed by exactly one validation dispatch before prepare; reject validation before implementation, missing/duplicate/late dispatch, and every post-prepare successor or authority action. Seed the retained log shape deterministically; no transcript/provider parser, product change, generic grammar, unrelated cleanup, silent overage, or immediate live rerun.

### Summary

REJECTED. Failure-safe retention worked and exposed the remaining false red: the
oracle mistakes the required validation-worker dispatch for a forbidden
successor dispatch. The exact Sonnet evidence remains at
`/tmp/spacedock-sonnet-cycle7.zlt1Yn`, and candidate `0691f25a9` is unchanged.

## Stage Report: implementation (cycle 8)

- DONE: Turn the exact retained Sonnet command-log shape into the deterministic positive: one implementation dispatch, then one validation dispatch, then one prepare/commit and an open unresolved attempt.
  Commit e6bd2a387 seeds the retained queued→implementation status/dispatch→validation status/dispatch→prepare/commit order directly; it failed red under the frozen single-dispatch predicate and now passes with failure-safe retained-tree grading unchanged.
- DONE: Replace the single-total-dispatch predicate with ordered stage-specific checks that still reject validation substitution, missing/duplicate/late dispatch, and every post-prepare successor or authority action.
  The oracle requires exactly two successful pre-prepare dispatch lines, with implementation first and validation second; focused mutants remove either stage, reverse them, duplicate validation, substitute validation, move implementation late, or add post-prepare dispatch, status, decision, consume, withdrawal, and duplicate prepare actions, and every mutant remains rejected.
- DONE: Remain within three files and cumulative +90/-15 by replacing/cutting, run focused/full/race/format/diff checks, and record the advisory disposition; stop on exact overage and do not spend Sonnet in implementation.
  The incremental correction is two test files, +12/-12; cumulative is exactly three files, +90/-11. Focused ordered-dispatch/retention and fixture tests, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` passed, and implementation round 8 validates `all-fixed`.
- SKIPPED: Spend another supported Sonnet or inspect transcript/provider evidence.
  The retained provider-neutral log supplied the authorized deterministic shape, and the same validator receives the committed candidate for any fresh live proof.

### Summary

Cycle 8 makes the headless oracle stage-specific: implementation and validation
are both required before the gate, while any dispatch after preparation remains
an authority crossing. The correction changes no product behavior and stays at
the existing cumulative insertion ceiling.

## Stage Report: validation (cycle 8)

- FAILED: Reconstruct the corrected queued -> implementation -> gated-validation journey and independently verify every prior authority-crossing mutant still fails.
  Exact candidate `e6bd2a387`; the focused queued/ordered-dispatch/retention suite passed in 2.8s, but a detached exact-trace audit proved both missing and late implementation-to-validation status transitions are accepted while both dispatches and the open-gate tail remain intact.
- SKIPPED: Run one fresh supported Sonnet live journey and preserve its durable state, gate-room, command-log, and open-unresolved gate evidence; do not substitute aggregate or cached evidence.
  The required cheap preflight exposed a deterministic AC-2 false green, so validation stopped before model spend and claims no prior, aggregate, cached, transcript, provider, or shell-derived substitute.
- DONE: Cross-check AC-1 through AC-3, the three-file +81/-10 surface, focused/full/race/format evidence, and recommend PASSED or REJECTED with the candidate frozen on any new finding.
  AC-1 remains green and AC-3 is three files at actual cumulative `+90/-11` within the captain-reset `+90/-15`, with focused/diff green and implementation-recorded full/race/format green. AC-2's exact lifecycle can still false-green, so the recommendation is REJECTED at frozen `e6bd2a387`.
- DONE: Classify the detached finding and retain the First Officer disposition before candidate mutation or rerun.
  Released workflow: implementation-to-validation handoff before validation dispatch; harm: missing or late handoff is certified as coherent; `value-ac[AC-2]:` requires entering validation after implementation and before its worker dispatch; trigger: both detached mutants return nil from the oracle. Material; task-owned evidence defect; FIX authorized.
- DONE: Bound and preserve the authorized handoff correction.
  Require exactly one successful `status=validation` after implementation dispatch and before validation dispatch; reject missing, late, duplicate, or other repair while retaining every prior dispatch and post-prepare authority mutant. Mutants remain at `/tmp/spacedock-cycle8-audit.aV1FQF/worktree`; no live rerun, parser, product change, unrelated cleanup, or overage.

### Summary

REJECTED. Ordered dispatches are now stage-specific, but the oracle does not
require the state handoff between them and accepts both missing and late
transitions. Sonnet was not spent, and candidate `e6bd2a387` is unchanged.

## Stage Report: implementation (cycle 9)

- DONE: Make the exact supported trace require one successful `status=validation` transition after implementation dispatch and before validation dispatch.
  Commit `43a1b8db4` partitions the command log around the two ordered dispatches, requires one exact successful validation transition in the middle interval, and rejects every successful status set after validation dispatch or gate preparation.
- DONE: Retain the ordered two-dispatch, queued-transition, post-prepare authority, retention, and prior mutation controls; add the preserved missing/late-transition mutants as red-first evidence.
  The missing and late transition mutants both failed red against frozen candidate `e6bd2a387`; they now pass as rejection cases alongside a duplicate-transition mutant and every prior missing, reordered, duplicated, substituted, late-dispatch, alternate-repair, gate, status, and retention control.
- DONE: Remain within three files and cumulative +90/-15 by replacing/cutting, run focused/full/race/format/diff checks, and record the advisory disposition; stop on overage and do not spend Sonnet in implementation.
  The correction is one existing test file, +7/-7; cumulative scope remains exactly three files, +90/-11. Focused transition/retention and queued-fixture tests, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` passed, and implementation round 9 validates `all-fixed`.
- SKIPPED: Run a live provider journey or inspect provider, transcript, parser, or shell evidence.
  The authorized deterministic mutants reproduced the defect and prove the bounded correction without another model spend or a change outside the test oracle.

### Summary

Cycle 9 makes the validation handoff both mandatory and position-specific while
preserving the ordered dispatch journey and every established authority guard.
The correction remains test-only and at the existing cumulative insertion cap.

## Stage Report: validation (cycle 9)

- DONE: Reconstruct the corrected queued -> implementation -> gated-validation journey and independently verify every prior authority-crossing mutant still fails.
  Exact candidate `43a1b8db4`; the repository focused suite passed in 2.2s, and a detached 21-case event table accepted only the complete supported trace while rejecting every required-event removal/reorder/duplicate and forbidden post-prepare decision, consume, withdrawal, status, duplicate prepare, or successor insertion.
- FAILED: Run one fresh supported Sonnet live journey and preserve its durable state, gate-room, command-log, and open-unresolved gate evidence; do not substitute aggregate or cached evidence.
  Exactly one fresh Sonnet run reached the correct open-gate outcome in 411.13s and its direct retained-tree regrade succeeded, but cleanup copied onto the same retained Git tree and failed the test. Inspection used only `workflow/`; no transcript/provider artifact was read and no rerun occurred.
- DONE: Cross-check AC-1 through AC-3, the three-file +81/-10 surface, focused/full/race/format evidence, and recommend PASSED or REJECTED with the candidate frozen on any new finding.
  AC-1 and the AC-2 lifecycle are evidenced by retained state, but AC-2's supported lane false-fails during duplicate retention; AC-3 is three files at actual cumulative `+90/-11` within `+90/-15`, with focused/diff green and implementation-recorded full/race/format green. Recommendation: REJECTED at frozen `43a1b8db4`.
- DONE: Classify the new finding and retain the First Officer disposition before mutation or rerun.
  Released workflow: successful supported Sonnet proof retention; harm: correct open-gate evidence is followed by a cleanup false failure; `value-ac[AC-2]:` requires a green supported proof from retained durable evidence; trigger: the 411.13s run plus detached first-copy success/second-copy read-only Git-object permission failure. Material; task-owned evidence-retention defect; FIX authorized.
- DONE: Bound and preserve the authorized single-copy outcome rule.
  Successful assertion path retains once; failure cleanup retains once; cleanup no-ops after successful retention. Retain the full event matrix and cheapest success/failure cleanup proof; no live rerun, framework, parser, provider/transcript/shell evidence, product change, unrelated cleanup, or overage. Evidence remains at `/tmp/spacedock-sonnet-cycle9.hoJ8hq` and `/tmp/spacedock-cycle9-audit.2nKrdz/worktree`.

### Summary

REJECTED. The supported Sonnet journey itself is now stage-coherent and stops at
one unresolved gate, but the harness copies its retained Git tree twice and
fails after success. Candidate `43a1b8db4` remains unchanged.

## Stage Report: implementation (cycle 10)

- DONE: Make existing workflow retention single-copy per outcome: success retains once, failure cleanup retains once, and cleanup does not overwrite an already retained tree.
  Commit `92dadede5` leaves direct success retention and absent-destination failure copying unchanged, but makes registered cleanup stop when `artifactDir/workflow` is already a directory, preventing the second write over read-only Git objects.
- DONE: Prove both success and failure cleanup paths with the existing retention/regrade test while preserving the 21-event journey matrix and all prior authority mutants.
  The focused test now directly retains and regrades the success tree, registers a separate failure destination, runs both cleanups, and regrades the failure-retained tree. It failed red at the duplicate success cleanup under frozen `43a1b8db4`; after the guard, both outcomes and the complete journey mutation test pass.
- DONE: Remain within three files and cumulative +90/-15 by replacing/cutting, run focused/full/race/format/diff checks, and record the advisory disposition; stop on overage and do not spend Sonnet in implementation.
  The correction is two existing test files, +10/-10; cumulative scope remains exactly three files, +90/-11. Focused success/failure retention, 21-event journey, and queued-fixture tests, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` passed, and implementation round 10 validates `all-fixed`.
- SKIPPED: Run another live provider journey or add an artifact framework, provider/transcript parser, shell grammar, or product behavior.
  The retained failure and deterministic permission reproduction isolate the duplicate-copy lifecycle, and the authorized existing mechanism now owns exactly one copy per outcome.

### Summary

Cycle 10 makes retained workflow ownership idempotent at cleanup: success keeps
its first copy, while failure still receives one cleanup copy. The full journey
oracle remains unchanged and the correction stays at the cumulative line cap.
