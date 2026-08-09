---
title: Restore live evidence whose original repair owners are complete
status: ideation
source: "Live-test-truth close reconciliation, 2026-08-09"
started: 2026-08-09T18:34:33Z
completed:
verdict:
score: 0.9
worktree:
issue:
pr:
mod-block:
sprint: test-behavior-completeness
group: common-evidence
sprint-readiness: ready
id: xp6c9qfe7y4wwp46enc3f85n
gates:
    version: 1
    records:
        - id: gate:xp6c9qfe7y4wwp46enc3f85n:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:xp6c9qfe7y4wwp46enc3f85n-backlog-1
              briefing:
                id: briefing:xp6c9qfe7y4wwp46enc3f85n:backlog:attempt-1:revision-1
                digest: sha256:df59dc0041b492b9dd552cc7d2e574b73036ebd23b4b924735f9243b2bbb91b9
                request-digest: sha256:497346c14775e33d34376b86d620b664986f817b81386f7169752ef061d2655b
                room-ref: ./restore-live-evidence-after-completed-repairs/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:xp6c9qfe7y4wwp46enc3f85n:backlog:1
                briefing: briefing:xp6c9qfe7y4wwp46enc3f85n:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:22.408925Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; the task is a late evidence-only capstone with exact target proof.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:xp6c9qfe7y4wwp46enc3f85n:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:xp6c9qfe7y4wwp46enc3f85n-ideation-1
              briefing:
                id: briefing:xp6c9qfe7y4wwp46enc3f85n:ideation:attempt-1:revision-1
                digest: sha256:0e0c2f5cc14a9b1f86efa9c40576c4772e31531fe80b8e1a5a89f5e897c11fcc
                request-digest: sha256:7d25c93ec3eb7f55727f6e47d83c464d685940d0ade98592f01a08ad581e9f6e
                room-ref: ./restore-live-evidence-after-completed-repairs/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:xp6c9qfe7y4wwp46enc3f85n:ideation:1
                briefing: briefing:xp6c9qfe7y4wwp46enc3f85n:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T21:33:28.842159Z"
                decision: approve
                reason: Captain approved the evidence-only capstone with no repair bytes.
              application:
                target-stage: implementation
                state: pending
---

## Problem

Seven ordinary xp6 target cells still use target-scoped `TODO` bindings. One
additional withdrawn-gate cell also needs an xp6-owned binding removal. The
source hides this evidence behind pre-execution skips.

The `liveTODO(...)` helper skips before the fixture builder and journey run. It
cannot report a repair, an expected semantic failure, or an infrastructure
failure. Strict XFAIL support from `ts7gq0mr9s3chx2w4wppd1kt` must define that
boundary before any binding edit.

The Codex `withdrawn-gate-recovery` rerun passed. Staff finding M5 keeps task
`47g` outside this sprint. xp6 removes that stale binding after an exact rerun.
Task `47g` remains the evidence origin, not the removal owner.

## Value

The common registry will show the result for each xp6 target cell. A passing
cell can remove its `TODO` or XPASS binding. A stable failure stays visible
with a named repair owner. An unexecutable cell keeps a `TODO` with its setup
reason.

The xp6 action baseline is eight target-scoped `TODO` bindings: seven ordinary
cells plus the withdrawn-gate removal. The expected result is four passing
cells, two stable semantic failures, and two unexecutable cells. No cell is
unclassified.

## Spike result

The riskiest mechanism was the live assertion boundary. A temporary isolated
checkout removed only target `TODO` rows. The real fixtures, hosts, and durable
assertions then ran at source commit `a929fcb60`.

The Codex `gate-guardrail` cell passed in 138.89 seconds. The Sonnet and Pi
`owned-conflict-owner-handoff` cells passed in 312.44 and 135.39 seconds.
These probes did not change the product checkout.

The Pi `gate-guardrail` probe failed with `state commit missing or before the
successful gate prepare`. The Pi `default-headless-gate-stop` probe failed with
`gated entity is not held at its open validation boundary`. These are semantic
failures, not skips.

The Opus lane skipped because the local benchmark token was empty and
`ANTHROPIC_API_KEY` was absent. The default Pi model stopped at an OpenRouter
credit error. A lower-cost Pi model enabled the two semantic probes.

The separate Codex `withdrawn-gate-recovery` rerun passed in 149.68 seconds.
This is an xp6-owned evidence binding removal. Task `47g` remains outside
sprint membership. The Codex `default-headless-gate-stop` failure remains the
known `98a` route.

The `/tmp` artifact paths are local probe locations. The result, command, and
failure text above are the durable evidence for this ideation record.

## Target classification

| Journey | Target | Classification | Binding or repair owner | Evidence |
| --- | --- | --- | --- | --- |
| `gate-guardrail` | Codex | passing | — | Exit 0; 138.89s; `codex-shared-scenarios/gate-guardrail` |
| `gate-guardrail` | Pi | stable semantic failure | `2e4fe65gy9vcr4xck6akzmdd` (`commit-pi-gate-prepare-before-presentation`) | `state commit missing or before the successful gate prepare`; lower-cost Pi model |
| `recorded-gate-lifecycle` | Claude Opus | unexecutable | xp6 `xp6c9qfe7y4wwp46enc3f85n` | Runner skipped because local Claude auth was unavailable |
| `default-headless-gate-stop` | Pi | stable semantic failure | `fh6rv0k6wr25zty0jjan4jp7` (`hold-pi-default-headless-validation-gate`) | `gated entity is not held at its open validation boundary`; lower-cost Pi model |
| `owned-conflict-owner-handoff` | Claude Sonnet | passing | — | Exit 0; 312.44s; `claude-shared-scenarios/owned-conflict-owner-handoff` |
| `owned-conflict-owner-handoff` | Claude Opus | unexecutable | xp6 `xp6c9qfe7y4wwp46enc3f85n` | Runner skipped because local Claude auth was unavailable |
| `owned-conflict-owner-handoff` | Pi | passing | — | Exit 0; 135.39s; lower-cost Pi model |
| `withdrawn-gate-recovery` | Codex | passing | xp6 `xp6c9qfe7y4wwp46enc3f85n` | Exit 0; 149.68s; exact rerun; stale TODO originated in `47g` |

The table has eight xp6 evidence cells: seven ordinary cells plus one
withdrawn-gate binding removal. Four pass, two have stable semantic failures,
two are unexecutable, and zero are unclassified. The `47g` entity stays outside
sprint membership.

## Repair routing

Staff finding M3 gives each Pi failure one owner and one mechanism.

1. `2e4fe65gy9vcr4xck6akzmdd` (`commit-pi-gate-prepare-before-presentation`)
   owns `gate-guardrail`. It starts from strict XFAIL code
   `gate-prepare-state-commit-missing`. It must commit and reread the prepared
   gate before the root session presents it.
2. `fh6rv0k6wr25zty0jjan4jp7` (`hold-pi-default-headless-validation-gate`)
   owns `default-headless-gate-stop`. It must identify the failed final-state
   clause before it chooses a repair. It must not assume the Sonnet or Codex
   worker-dispatch mechanism.

These tasks are separate. xp6 does not implement either repair or invent either
semantic code. The `98aa776adg66gn823a8gamdq` task remains the owner of the
adjacent Codex and Sonnet default-headless failure. It does not expand xp6.

The withdrawn-gate row needs no repair owner. xp6 owns its evidence-only TODO
removal. Task `47g` remains outside sprint membership and does not receive this
binding change.

## Required order

1. Land or expose strict XFAIL behavior from `ts7gq0mr9s3chx2w4wppd1kt` at the
   durable assertion boundary.
2. Run each xp6 target on its exact candidate and record the result artifact.
3. Remove a binding only after its exact target passes its durable assertion.
4. Keep each Pi failure with its named repair owner and its exact observed text.
   Do not combine the two mechanisms.
5. Keep each Opus target as `TODO` while authenticated execution is unavailable.
6. Stop on an unclassified result. Do not convert it to `TODO` or XFAIL.
7. Keep `47g` outside the sprint. After its passed withdrawn-gate rerun, xp6
   removes the stale binding as an evidence-only change. Do not assign this
   removal to `47g` or defer it outside xp6.

## Proposed approach

Use the existing common journey entry points, fixtures, exercises, and
assertions. Any later xp6 source edit is limited to passing binding removal in
`internal/ensigncycle/shared_live_runner_test.go`.

The four passing xp6 cells are Codex `gate-guardrail`, Sonnet
`owned-conflict-owner-handoff`, Pi `owned-conflict-owner-handoff`, and Codex
`withdrawn-gate-recovery`. Remove only their target bindings after exact-
candidate evidence.

Keep the two Pi failure records routed to their separate repair tasks. Keep both
Opus `TODO` bindings until an authenticated pre-release run can execute them.
Do not remove a TODO because a host lacks auth or credits.

The `47g` withdrawn-gate pass remains evidence from an outside-sprint entity.
Its stale binding is still an xp6 removal. No deferred repair enters xp6 scope.

No new scenario table, simulator, helper-only change, or component-only landing
is allowed. A product-repair landing is also outside xp6.

## Acceptance criteria

**AC-1 (VALUE) — Four passing xp6 evidence cells have no target-scoped TODO or XPASS binding.**

The Codex `gate-guardrail`, Sonnet `owned-conflict-owner-handoff`, and Pi
`owned-conflict-owner-handoff` probes pass their durable assertions. The Codex
`withdrawn-gate-recovery` probe also passes. Each binding removal is
evidence-only.

**AC-2 (VALUE) — The two Pi semantic failures have separate named repair owners.**

`gate-guardrail` routes to `2e4fe65gy9vcr4xck6akzmdd`; `default-headless-gate-stop`
routes to `fh6rv0k6wr25zty0jjan4jp7`. Each owner has a distinct mechanism. A
different code or infrastructure error fails the lane.

**AC-3 — Two Opus target cells retain honest TODO bindings while auth is unavailable.**

The `recorded-gate-lifecycle` and `owned-conflict-owner-handoff` Opus lanes
record the auth-missing setup reason. They do not claim semantic evidence.

**AC-4 (VALUE) — All eight xp6 evidence cells have one classification.**

The final xp6 state records four passing cells, two stable semantic failures,
two unexecutable Opus cells, and zero unclassified cells. Seven cells are
ordinary xp6 targets. The eighth is the `47g`-origin withdrawn-gate binding
removal, which xp6 owns while `47g` stays outside sprint membership.

**AC-5 — This task remains evidence-only.**

The correction changes only this workflow record. The eventual xp6 product
landing can remove passing bindings only. It adds no repair code, component-only
change, fixture change, durable assertion change, desired registry entry,
command grammar, stored format, authority rule, or metrics format.

## Test plan

Run the strict-XFAIL grade and metrics tests from `ts7gq0mr9s3chx2w4wppd1kt`
before any binding edit. A matching code must produce XFAIL. A pass must produce
XPASS. A different code, launch failure, timeout, or auth failure must remain
FAIL or TODO.

Run each focused common journey on its exact candidate. Use local subscription
auth before paid CI. A lower-cost Pi model can diagnose behavior. It cannot
justify removal from the default Pi lane.

Run the exact Codex `withdrawn-gate-recovery` candidate before removing its
stale TODO. The pass permits the binding removal. It does not move `47g` into
the sprint.

Run these deterministic checks after any binding edit:

- `go test ./internal/ensigncycle -run 'TestGateGuardrailNegativeBrokenStateTransition|TestAssertGateHeld'`
- `go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle'`
- `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$'`
- `SPACEDOCK_LIVE_STATE_DIR=docs/dev/.spacedock-state go test ./internal/contractlint -run '^TestRuntimeLiveTODOOwnersAreActive$'`

Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
No new fixture, simulator, runner test, or reconciliation test is needed.

## Expected surface and estimates

The product surface is one existing test file. The state entity and report are
workflow records, not product LOC. M5 keeps product net at zero.

| Surface | Gross additions | Gross deletions | Net | Purpose |
| --- | ---: | ---: | ---: | --- |
| `internal/ensigncycle/shared_live_runner_test.go` | ~3 | ~3 | 0 | Remove four passing xp6 TODO entries across three binding rows after exact evidence |
| Pi repair tasks | 0 | 0 | 0 | Product repairs stay with the two named owners |
| **Product total** | **~3** | **~3** | **0** | Binding-only evidence landing; no repair bytes |

The estimate has a tolerance of two changed lines. A component-only landing is
not acceptable. A product-repair landing is not acceptable. If no exact passing
binding remains after reconciliation, xp6 lands no product source change.

## Observable semantic boundary

- Command grammar: unchanged.
- Stored formats: unchanged.
- Authority: unchanged.
- Runtime behavior: no repair behavior changes in xp6.
- Evidence: use the existing journey result and artifact.
- Documentation: this state record carries the plan and report.

## Scope

- Reuse the current common journey entry points, fixtures, exercises, and assertions.
- Record the exact result for each of the eight xp6 evidence cells.
- Remove all four passing xp6 bindings after exact evidence.
- Name the two Pi repair owners and keep their mechanisms separate.
- Keep the two Opus TODOs only while authenticated execution is unavailable.
- Keep `47g` outside sprint membership while xp6 removes its passed cell binding.

## Out of scope

- Repairing either Pi failure.
- Combining the two Pi failures under one owner.
- Weakening a durable assertion.
- Adding a simulator or a new scenario table.
- Adding tests for reconciliation code.
- Changing the desired journey registry.
- Changing First Officer behavior, gate authority, or stored formats.
- Adding deferred task `47g` to sprint membership or assigning its binding removal to `47g`.
- Landing component-only or product-repair code.

## Stage Report: ideation

- DONE: Keep this task evidence-only and include the withdrawn-gate Codex cell from deferred task 47g.
  M5 keeps `47g` outside sprint membership, while xp6 removes its passed stale binding after the exact rerun. No repair enters xp6.
- DONE: Classify each target as passing, stable semantic failure, unexecutable, or unclassified.
  The corrected matrix has eight evidence cells: four passing, two stable semantic failures, two unexecutable Opus cells, and zero unclassified cells. Seven are ordinary xp6 cells; one is the `47g`-origin removal.
- DONE: Give gross and net line estimates with product-repair routing outside this capstone.
  The binding-only estimate is about three additions and three deletions, net zero. Pi repairs route to `2e4fe65gy9vcr4xck6akzmdd` and `fh6rv0k6wr25zty0jjan4jp7`; no repair bytes enter xp6.

### Summary

This corrected ideation package is evidence-only. It records eight evidence
cells, routes the two distinct Pi failures to their named owners, and keeps both
Opus TODOs until authenticated execution is available. xp6 removes the passed
`47g`-origin binding while that deferred task stays outside the sprint.

## Stage Report: ideation (cycle 2)

- DONE: Keep this task evidence-only and include the withdrawn-gate Codex cell from deferred task 47g.
  M5 keeps `47g` outside sprint membership. xp6 owns the passed cell's stale TODO removal after the exact rerun.
- DONE: Classify each target as passing, stable semantic failure, unexecutable, or unclassified.
  M3 is folded into the eight-cell matrix. Seven cells are ordinary xp6 targets. The eighth is the passed withdrawn-gate removal. It names `2e4fe65gy9vcr4xck6akzmdd` for the Pi gate-commit failure and `fh6rv0k6wr25zty0jjan4jp7` for the Pi headless-hold failure.
- DONE: Give gross and net line estimates with product-repair routing outside this capstone.
  The binding-only estimate is about three additions and three deletions, net zero. No component-only or product-repair landing is allowed.
