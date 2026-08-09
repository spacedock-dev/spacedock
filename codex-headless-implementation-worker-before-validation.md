---
title: Make headless Sonnet and Codex spawn implementation before validation
status: ideation
source: "PR #583 run 31320596435, Codex job 93262943132 and Sonnet job 93262943118, 2026-08-09"
started: 2026-08-09T18:34:21Z
completed:
verdict:
score: 0.95
worktree:
issue:
pr:
mod-block:
sprint: test-behavior-completeness
group: common-evidence
sprint-readiness: ready
id: 98aa776adg66gn823a8gamdq
gates:
    version: 1
    records:
        - id: gate:98aa776adg66gn823a8gamdq:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:98aa776adg66gn823a8gamdq-backlog-1
              briefing:
                id: briefing:98aa776adg66gn823a8gamdq:backlog:attempt-1:revision-1
                digest: sha256:b0d7c918d0addf124f3fa3b605e3ed074b6824a05a5acb69734299696260b51f
                request-digest: sha256:54e4b661e7ba07b41d0a51b198e3a987479f991259beed74bf82eb05ff0f3610
                room-ref: ./codex-headless-implementation-worker-before-validation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:98aa776adg66gn823a8gamdq:backlog:1
                briefing: briefing:98aa776adg66gn823a8gamdq:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:07.295318Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; the seed names the exact live failure, repair boundary, and XFAIL-first order.
              application:
                target-stage: ideation
                state: consumed
---

## Problem

The Sonnet and Codex default-headless-gate-stop journeys built an implementation
dispatch package, but the durable oracle still found no worker dispatch before
validation.

Both runs moved the entity from implementation to validation and then prepared
the validation gate without accepted worker evidence. The oracle rejected both
runs with: "gate hold crossed its committed no-authority boundary:
implementation was not dispatched before validation".

The failures occurred on candidate
7c8708c8537fc73761e56813ddbd6a498959ef19. Codex ran for 166.56 seconds.
Sonnet ran for 414.17 seconds. The Sonnet transcript contains bare Agent calls
and a subagent transcript, but the durable command-log grade still rejected the
boundary. This shows that textual activity is not enough; the spawn and
completion must be observable through the bound worker mechanism.

## Value

A headless First Officer must complete the required worker before it presents
the next gate. A dispatch package, a status change, or narration is not a
worker dispatch.

## Proposed approach

Wait for task ts7gq0mr9s3chx2w4wppd1kt to commit strict XFAIL evidence for both
runtime lanes. Then add one host-neutral guard to the shared First Officer
dispatch contract:

1. After dispatch build exits successfully, call the bound worker.spawn in the
   same turn.
2. Forward every field emitted by dispatch build without changes.
3. Record the worker handle before a status change, report read, gate prepare,
   or wait.
4. Wait for the completion signal before advancing to validation.
5. Treat dispatch build, status mutation, narration, a self-authored report,
   and an empty wait as no worker evidence.

This guard applies to the initial worker only. The existing advance path keeps
its reuse-handle rule. The shipped Claude and Codex adapters already map the
shared envelope to their native Agent and spawn_agent calls, so they need no
adapter change. The binary cannot call a host-native worker API.

## Required order

Task ts7gq0mr9s3chx2w4wppd1kt must land before this product repair starts. It
must run the Sonnet and Codex cells as strict XFAIL evidence.

Do not change First Officer behavior before both cells report the expected
implementation-worker-not-dispatched failure code. After the repair, XPASS
must force removal of each repaired XFAIL binding.

## Acceptance criteria

**AC-1 (VALUE) — Sonnet and Codex run the implementation worker before validation.**

The default-headless-gate-stop journey observes one implementation worker spawn.
It also observes worker completion before the validation transition.

Verified by: strict-XFAIL runs first, then passing Sonnet and Codex runs of
TestLiveCommonDefaultHeadlessGateStop with the current fixture and durable
oracle.

**AC-2 — A status change cannot replace the worker dispatch.**

The First Officer must not credit implementation from dispatch build, narration,
a direct status change, a self-authored report, or an empty wait.

Verified by: the provider-neutral command log and durable entity reject each
missing-worker form before gate prepare.

**AC-3 — Each runtime TODO is removed only after passing live evidence.**

The journey keeps each runtime TODO, or its strict-XFAIL binding, until the exact
candidate passes that runtime lane.

Verified by: registry reconciliation and the passing live artifact use the same
commit.

## Test plan

- First, run task ts7gq0mr9s3chx2w4wppd1kt. The Sonnet and Codex cells must
  execute the real fixture and report XFAIL with the named semantic code.
- After the guard, run focused shared dispatch tests. Check that the emitted
  worker envelopes and adapter mapping stay unchanged.
- Run the exact live selector
  TestLiveCommonDefaultHeadlessGateStop on the repair commit for Sonnet and
  Codex. Require one worker spawn, completion before validation, an open first
  gate, and no post-prepare authority crossing.
- Remove each XFAIL binding only after its exact runtime lane passes. Reconcile
  the registry with the same commit, then run go test ./... and go test ./... -race.
- Keep the live runs as the acceptance proof. The observed run cost was 166.56
  seconds for Codex and 414.17 seconds for Sonnet; no simulator is proposed.

## Expected surface

The product repair is expected to touch two files:

- skills/first-officer/references/fo-dispatch-core.md: five inserted logical
  lines for the initial spawn and completion guard.
- docs/runtime-live-ci-registry.md: three inserted and two deleted lines to
  state that the worker completes before the first gate.

Estimate: 8 gross insertions, 2 gross deletions, and 6 net lines. Tolerance is
plus or minus 2 lines, with no additional files. This ideation stage changes
only the state entity, not those product files.

## Semantic changes

The initial implementation stage must have a bound worker spawn and a
completion signal before validation can start. Command grammar, stored state
formats, authority labels, the reuse-handle advance path, and Pi behavior stay
unchanged.

## Documentation diff

Current registry outcome: "A headless launch without decision authority
advances from the preceding stage to the first human gate, presents it, and
stops open."

Proposed outcome: "A headless launch without decision authority dispatches and
completes the preceding stage worker, then advances to the first human gate,
presents it, and stops open."

## Scope

- Correct the shared First Officer behavior that lets Sonnet and Codex skip the
  implementation worker.
- Preserve the coherent 26n fixture and its no-authority oracle.
- Use each shipped runtime worker-dispatch mechanism.
- Keep product behavior portable. Do not add a test-only host switch.

## Out of scope

- Weakening the live oracle.
- Restoring the contradictory fixture.
- Treating dispatch build as worker-spawn evidence.
- Adding a simulator or a test for test infrastructure.
- Repairing Pi behavior. Pi has no new failure evidence from this run.

## Evidence

Run 31320596435, jobs 93262943132 and 93262943118, records this sequence:

1. The First Officer changed `queued` to `implementation`.
2. The First Officer created the implementation dispatch package.
3. The durable oracle did not accept an implementation worker dispatch.
4. The First Officer changed `implementation` to `validation`.
5. The First Officer prepared and committed the validation gate.

The Codex artifact has dispatch build output and an empty native wait, followed
by a validation transition. The Sonnet artifact has bare Agent calls and a
subagent transcript, but the same durable oracle rejection. The run proves a
product behavior gap and identifies the missing observable spawn/completion
boundary. A local replay was not authenticated, so it was not used as proof.

## Stage Report: ideation

- DONE: Exercise the current Sonnet and Codex failure path before selecting the repair.
  GitHub run 31320596435 failed both lanes with the exact no-authority error.
  Artifact inspection found Codex dispatch-build output without an accepted
  worker boundary; Sonnet had Agent transcript activity but the same durable
  oracle rejection.
- DONE: Define the smallest portable worker-spawn repair after committed strict XFAIL evidence.
  The proposed repair is one shared-core initial-dispatch guard. It preserves
  adapter mappings, the advance reuse handle, the fixture, and the oracle.
  Product work waits for task ts7gq0mr9s3chx2w4wppd1kt strict XFAIL evidence.
- DONE: Give gross and net line estimates with exact-candidate live proof.
  The repair estimate is 8 insertions, 2 deletions, and 6 net lines, with plus
  or minus 2 lines tolerance and no extra files. The evidence candidate is
  7c8708c8537fc73761e56813ddbd6a498959ef19; Codex took 166.56 seconds and
  Sonnet took 414.17 seconds.

### Summary

Ideation records an authenticated failure on the exact candidate and defines a
small, portable shared-core guard. The next stage must land strict XFAIL
evidence first, then prove spawn and completion ordering in both live lanes.
