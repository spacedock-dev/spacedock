---
title: Make the headless recorded-gate stop fixture stage-coherent
status: ideation
source: "PR #580 run 30591046287, Sonnet job 91033369022"
started: 2026-07-31T00:32:10Z
completed:
verdict:
score: 0.9
worktree:
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
                state: pending
                blockers: []
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
the initial fixture has exactly one dispatchable entity at implementation, no
selected gate, no completed validation report, and an implementation definition
whose next state is gated validation. Verified by a focused fixture test that
boots the real fixture through `spacedock status --boot --identify --json` and
asserts the resulting dispatchable/ready-gate state.

**AC-2 — Headless no-authority remains load-bearing.** A supported Sonnet run
dispatches implementation once, enters validation before preparation, binds and
commits exactly one validation Briefing, presents it, and stops open. After the
first successful prepare there is no decision, consume, withdrawal, status
repair, duplicate prepare, or successor dispatch. Verified by the durable entity,
gate room, state Git history, and provider-neutral command log; the final attempt
has no Resolution or Application.

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
