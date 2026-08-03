---
title: Deliver one portable common live-journey surface
status: ideation
source: "Captain recarve of live-test-truth, 2026-08-03. Absorbs 3w, h3, tj, and r4 as design inputs."
score: 1.0
sprint: live-test-truth
group: portable-common-surface
sprint-readiness: ready
id: ys7ncwh9kr8w5h9hdkz5apat
gates:
    version: 1
    records:
        - id: gate:ys7ncwh9kr8w5h9hdkz5apat:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ys7ncwh9kr8w5h9hdkz5apat-backlog-1
              briefing:
                id: briefing:ys7ncwh9kr8w5h9hdkz5apat:backlog:attempt-1:revision-1
                digest: sha256:ff92cf6104cdc81f88ab51f6a59f33f3d1f42e72d7e72f712a2f50d5ae61fc47
                request-digest: sha256:f9af92f7f1823701ddcf301aadddf68012f7af8eb4656ea6a6ced3772ec14e48
                room-ref: ./deliver-portable-common-live-journey-surface/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ys7ncwh9kr8w5h9hdkz5apat:backlog:1
                briefing: briefing:ys7ncwh9kr8w5h9hdkz5apat:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T12:18:32.293007Z"
                decision: approve
                reason: Captain explicitly approved the outcome-shaped recarve and directed immediate redispatch.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:ys7ncwh9kr8w5h9hdkz5apat:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ys7ncwh9kr8w5h9hdkz5apat-ideation-1
              briefing:
                id: briefing:ys7ncwh9kr8w5h9hdkz5apat:ideation:attempt-1:revision-1
                digest: sha256:035db0434c5ab60a4a04e6ee1151d00d7e09b6d49f9f1e3e8a7d34ac24fbfa16
                request-digest: sha256:dace16b8c6dea83386db8d0da677e97ea55a8e5c87ea975c3863ec1f32318a7d
                room-ref: ./deliver-portable-common-live-journey-surface/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T14:22:43.42636Z"
                reason: Preflight staff review found cross-member ownership and sequencing defects, and the shared sprint index changed. Withdraw this stale binding before the authorized fold.
            - id: gate-attempt:ys7ncwh9kr8w5h9hdkz5apat-ideation-2
              briefing:
                id: briefing:ys7ncwh9kr8w5h9hdkz5apat:ideation:attempt-2:revision-1
                digest: sha256:077a66a2949eaceae818bb2efbaf376e0d6d79683ec79413b4db4c516d30512c
                request-digest: sha256:0a0c021d5e5516a7b11bcb971078dd0a4601da10862ceed9337c581e2e9a735d
                room-ref: ./deliver-portable-common-live-journey-surface/review/ideation/briefing-2
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T14:33:49.183562Z"
                reason: Preflight review closure changed the shared sprint index after attempt 2 was frozen. Replace it with a final package that binds the review artifact.
started: 2026-08-03T12:19:31Z
---

## Outcome

A test operator can select one common journey by one stable identity on Claude, Codex, or Pi. The registry makes missing journeys and orphan fixtures visible.

The task delivers the complete operator journey. Registry annotations, fixture bindings, runtime adapters, selectors, and reconciliation are implementation steps.

## Problem

The desired registry contains 16 common journeys. The current shared scenario table contains 10 journeys.

Claude and Codex use separate top-level suite names and separate journey maps. Pi has only a coverage map and two standalone live paths.

Six desired journeys still use standalone tests. Source declarations also lack stable registry bindings.

These splits make runtime selection part of journey identity. They also hide missing runners, unused fixtures, and lane-specific behavior drift.

## Design evidence

- `bind-live-registry-to-source` (`3w`) proved one adjacent annotation join and the path-scoped reconciliation guard.
- `converge-shared-live-suite-entrypoint` (`h3`) ran `shallow-boot` through one selector on Claude and Codex.
- `add-pi-common-live-runner` (`tj`) ran the same fixture and assertion through the real Pi front door.
- `promote-standalone-common-live-journeys` (`r4`) mapped all six missing journeys and proved the closest gate-assertion overlap.
- `ac2-reanchor-live-scenario-repair` proved a durable correct branch and a durable wrong branch for `ac-value-reanchor`.

No new ideation spike is necessary. The combined plan joins two proved seams: the common selector and the Pi adapter boundary.

Implementation starts with the exact Pi selector as a focused test. A failure stops migration before the full Pi suite incurs cost.

## Dependencies

`3d` lands before this task. It owns the AC2 durable-oracle repair and unit tests.

`15e` lands after `3d` and before this task. This task lands last and consumes the final behavior from both tasks.

## Target behavior

`TestLiveSharedScenarios` is the only authoritative top-level test for common journeys. It contains exactly 16 registry subtests.

Every lane sets `SPACEDOCK_LIVE_RUNTIME` to `claude`, `codex`, or `pi`. Every lane invokes the same test selector.

One common runner map owns every journey. One adapter registry owns runtime selection.

Common runners own fixture selection, prompts, normalized evidence, and host-neutral assertions. Adapters own authentication, launch, output parsing, liveness, and host artifacts.

Each source declaration carries its stable registry annotation. Reconciliation reports desired gaps and fails ambiguous or orphaned bindings.

The Pi lane runs all 16 journeys. A coverage row, quarantine, skip, or Pi-only top-level test does not count as evidence.

## Complete approach

1. Add focused tests for one suite identity, adapter selection, map parity, and exact workflow selectors.
2. `3d` repairs the `ac-value-reanchor` durable oracle before this task promotes its common runner.
3. Add the four approved annotation forms to declarations that already exist.
4. Create `TestLiveSharedScenarios` and the common 10-journey runner map.
5. Move Claude and Codex transport behavior into adapters without changing their launch contracts.
6. Add the six missing journey records and their common runners in registry order.
7. Extract named fixture builders and host-neutral assertions from the retired standalone tests.
8. Add `pi` to the same adapter registry and implement its session, child, trace, and artifact behavior.
9. Delete parallel suite names, host runner maps, standalone common tests, Pi coverage substitutes, and Pi quarantine code.
10. Update workflow selectors, workflow guards, metrics, documentation, and source annotations once against the final shape.
11. Run semantic reconciliation and record the candidate commit SHA.
12. Commit the SHA as a documentation-only change and run its path guard.
13. Run focused offline tests, three representative live runs, complete live lanes, full tests, race tests, and formatting.

The adapter registry serves AC-1 and AC-4. Compatibility wrappers are insufficient because they preserve multiple authoritative identities.

The common runner map serves AC-2 and AC-4. Three runtime maps are insufficient because they permit private common journeys.

Stable source annotations serve AC-2 and AC-3. Builder names in the registry are insufficient because source moves can change desired state.

The manual inventory and SHA guard serve AC-2 and AC-3. A permanent parser adds scope without proving runtime behavior.

The Pi trace reader serves AC-4. Durable state alone cannot prove Pi child identity or required tool selection.

## Collision map

| Shared surface | Colliding evidence | One final disposition |
|---|---|---|
| `shared_scenarios_test.go` | `3w` adds bindings. `r4` adds six journeys. `tj` needs Pi parity. | Define 16 annotated records once. Compare one runner map and three adapters in both directions. |
| `shared_live_runner_test.go` | `h3` creates the common runner. `r4` adds six runners. | Build one 16-journey map. Do not create host-specific journey maps. |
| Claude and Codex runner files | `h3` moves common behavior. `3w` adds suite bindings. | Keep only transport adapters and host-specific evidence. Put the suite annotation on the common test. |
| Pi runner files | `tj` proposed a Pi suite and map. `r4` expands the final journey count. | Add one Pi adapter for the common 16-journey map. Do not add `TestLivePiSharedScenarios`. |
| Fixtures and assertions | `3w` binds builders. `r4` extracts six contracts. | Extract each builder once, add its annotation there, and reuse existing assertions. |
| `ac-value-reanchor` | `3d` changes the oracle. `r4` promotes the result. | Depend on the repair from `3d`. Remove only the standalone wrapper during promotion. |
| Workflow and release guards | All four inputs change selectors or coverage. | Edit these files after all runners exist. Require the exact common selector in every lane. |
| `docs/runtime-live-ci.md` | All four inputs add overlapping instructions. | Apply one final documentation change with three commands, 16 journeys, costs, artifacts, reconciliation, and one SHA. |

## Journey and fixture reconciliation

The final shared table contains these 16 IDs in registry order:

1. `full-ensign-cycle`
2. `gate-guardrail`
3. `default-headless-gate-stop`
4. `withdrawn-gate-recovery`
5. `recorded-gate-lifecycle`
6. `rejection-flow`
7. `feedback-3-cycle-escalation`
8. `merge-hook-guardrail`
9. `filing`
10. `shallow-boot`
11. `zero-discovery`
12. `auto-continue-after-implementation`
13. `self-evidence-merge-triage`
14. `smallest-sufficient-mechanism`
15. `keep-moving-posture`
16. `ac-value-reanchor`

Keep both `auto-continue` fixtures under one journey. Artifact labels include the fixture ID, but no nested journey identity exists.

Use one gate-stop assertion primitive for `gate-guardrail` and `default-headless-gate-stop`. Add start-state assertions to prevent cross-qualification.

Preserve withdrawn-room bytes for `withdrawn-gate-recovery`. Reject any decision, consume, or dispatch after the successor review opens.

Grade `ac-value-reanchor` only from durable decision state. Prompt wording and final-message text do not count as proof.

The reconciliation inventory uses the diagnostic classes from `3w`. Desired gaps are `MISSING` or `UNSELECTED`.

Duplicate, invalid, orphan, and unaccounted results fail reconciliation. Each nonzero result lists its exact IDs or paths.

## Pi execution contract

The Pi adapter loads the current checkout through `spacedock pi`. It uses isolated Pi state and pinned local packages.

The adapter records root output, session JSONL, child JSONL, model stamps, process status, cost, duration, and durable workflow evidence.

Each journey has one 12-minute process deadline. The suite stops after its first failed journey and performs no automatic live retry.

The 10-journey price table totals $1.58. The `shallow-boot` spike measured $0.0269991 against its $0.03 ceiling.

The six promoted journeys add these planning ceilings:

| Journey | Ceiling |
|---|---:|
| `full-ensign-cycle` | $0.25 |
| `default-headless-gate-stop` | $0.08 |
| `withdrawn-gate-recovery` | $0.10 |
| `zero-discovery` | $0.03 |
| `auto-continue-after-implementation` | $0.27 |
| `ac-value-reanchor` | $0.15 |

The six ceilings total $0.88. The complete estimate is $2.46.

A 25 percent reserve gives a $3.08 approval ceiling for one full Pi run.

These values are planning limits, not pass criteria. The lane records actual tokens, cost, duration, and model for each journey.

Run non-mutating journeys before mutating journeys. Within each group, run lower-cost journeys first.

## Landing sequence

Use one branch and one final workflow change. Do not merge partially authoritative selectors between steps.

1. Land focused failing tests for the final suite identity and promotion boundaries.
2. Land source annotation grammar and bindings for declarations that already exist.
3. Land the canonical Claude and Codex entry point with the first 10 common runners.
4. Land the six promoted journeys, fixture variants, and retired standalone tests.
5. Land the Pi adapter against the complete 16-journey common map.
6. Land final workflow selectors, release guards, documentation, and reconciliation results.
7. Record the candidate SHA in a final documentation-only commit.

Do not enable the SHA guard before watched source changes are complete. Otherwise, each intermediate commit creates a known stale result.

## Acceptance criteria

**AC-1 (VALUE) — One journey identity works on every supported runtime.**

The same `^TestLiveSharedScenarios$/^shallow-boot$` selector passes on Claude, Codex, and Pi. Only runtime configuration changes.

The workflow guard proves that every live lane uses this top-level identity. A restored old suite name makes the guard fail.

**AC-2 — All desired common journeys have canonical executable identities.**

Registry, scenario-table, and runner-map parity is 16 of 16. The count of standalone common-journey top-level tests is zero.

Reconciliation reports 16 bound journey IDs and no missing suite lane. Removing one runner or selector makes the proof fail.

**AC-3 — Fixtures have accountable use.**

Each common fixture ID resolves to one annotated builder and at least one journey. Proof-only and experiment fixtures resolve to their allowed registry class.

Reconciliation reports zero duplicate, invalid, orphan, unaccounted-test, and unaccounted-builder results. A temporary orphan fixture makes the proof fail.

**AC-4 — Runtime differences stay behind adapters.**

Common fixtures and assertions contain no host-specific branch. The adapter factory contains the only runtime-selection branch.

Fake-adapter tests prove the boundary. The complete live lanes prove unchanged launch, liveness, metrics, and artifact behavior.

## Test plan

Run these focused offline tests before any model-backed run:

```bash
go test -tags live ./internal/ensigncycle -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestSharedLiveRuntimeSelection|TestPromotedCommonJourneyEntrypoints' -count=1
go test ./internal/ensigncycle -run 'TestGateJourneyOverlap|TestAutoContinue|TestZeroDiscover|TestEnsignCycleGoesRed|TestACReanchor' -count=1
go test ./internal/release -run 'TestRuntimeLiveWorkflow|TestWorkflowsPreserveAndPublishJourneyCosts' -count=1
```

The parity tests compare the registry IDs, scenario table, runner map, and adapter registry in both directions.

The negative tests remove one runner, add one orphan fixture, swap gate starts, stop auto-continue, and select the wrong re-anchor branch.

Run one representative journey with the same selector on every runtime:

```bash
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_REQUIRED=1 go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
```

Run the six promoted selectors on Claude and Codex. Then run those same selectors on Pi.

Run each complete live lane through `TestLiveSharedScenarios`. The Pi log must contain 16 real passes and no skip.

Run the required repository checks:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

Run the SHA guard with one stale base and the recorded base. The stale base must name each changed watched path.

## Expected surface

The baseline is 34 files, about 2,060 insertions, and about 1,650 deletions. The tolerance is 42 files and 2,750 insertions.

| Surface | Expected change | Purpose |
|---|---:|---|
| `internal/ensigncycle/shared_live_runner_test.go` | new, about 560 insertions | One suite, one runner map, 16 runners, and adapter interface. |
| `internal/ensigncycle/shared_live_runner_unit_test.go` | new, about 180 insertions | Selection, parity, fixture variants, and fake-adapter negatives. |
| `internal/ensigncycle/shared_scenarios_test.go` | about 70 insertions | Six records and 16 journey annotations. |
| Claude and Codex live runner files | about 150 insertions and 670 deletions | Keep transport adapters and remove duplicate journey maps. |
| Pi shared runner and trace files | two new files, about 500 insertions | Add Pi launch, sessions, traces, costs, and artifacts. |
| Shared fixture and assertion files | about 320 insertions | Extract six contracts, gate composites, variants, and annotations. |
| Six standalone common live tests | delete about 690 lines | Remove competing top-level identities. For AC2, remove only its standalone wrapper. |
| Pi coverage and recorded-gate quarantine files | delete about 185 lines | Remove non-executable substitutes. |
| `shared_coverage_meta_test.go` and `journey_metrics_live_test.go` | about 70 insertions and 45 deletions | Require 16 runners, three adapters, and stable metrics. |
| Seven proof, experiment, and existing fixture files | about 25 insertions | Add proof, experiment, and fixture annotations. |
| `.github/workflows/runtime-live-e2e.yml` | about 30 insertions and 20 deletions | Use one selector, configure three adapters, and add the SHA guard. |
| Two `internal/release` workflow tests | about 55 insertions and 30 deletions | Reject old selectors and prove exact lane commands. |
| `docs/runtime-live-ci.md` | about 90 insertions and 30 deletions | Add final commands, 16 journeys, Pi costs, reconciliation, and SHA. |
| `docs/site/contributing/architecture-notes.md` | about 10 insertions and 6 deletions | Describe the common suite and adapter boundary. |

Adjacent fixture or parser test files can enter the tolerance. New product packages, registry IDs, or authority commands require a new gate decision.

## Observable semantic boundary

- Command grammar does not change.
- Stored workflow formats do not change.
- First-officer and ensign authority do not change.
- Product runtime behavior does not change.
- Live-test identity changes from host-specific names to one common name.
- The scenario count changes from 10 implemented records to all 16 desired records.
- Pi changes from coverage metadata to required executable evidence.
- CI gains exact runtime configuration and a path-scoped stale-reconciliation guard.
- Artifact and metrics records keep stable journey and fixture IDs.

## Documentation change

Replace the two host-specific local commands and the quarantined Pi command with these commands:

```text
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 40m -run TestLiveSharedScenarios ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run TestLiveSharedScenarios ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_REQUIRED=1 go test -tags live -count=1 -timeout 40m -run TestLiveSharedScenarios ./internal/ensigncycle -v
```

Replace the current per-host runner description with this text:

> `TestLiveSharedScenarios` owns one subtest identity for each common journey. Each lane selects a runtime adapter through `SPACEDOCK_LIVE_RUNTIME`.

Add this normative statement:

> `docs/runtime-live-ci-registry.md` is the normative desired-state registry. Source annotations bind its stable IDs to current declarations.

Add the reconciliation procedure, diagnostic classes, current inventory, watched paths, and this record:

```text
Registry reconciliation SHA: `<40-character candidate SHA>`
```

The procedure states that source or selector changes require reconciliation in the same pull request.

## Stage Report: ideation

- DONE: Compose the four completed reports into one outcome plan without repeating valid live spikes.
  The plan joins the proved Claude, Codex, Pi, registry, and gate seams. It records that no new ideation spike is necessary.
- DONE: Resolve file overlap and define one landing sequence for the registry, canonical entry point, journey promotion, and Pi runtime.
  The collision map removes the proposed Pi-only suite and orders seven dependency-safe landing steps.
- DONE: Produce one complete plan whose acceptance evidence is visible through the operator-facing test surface.
  The plan defines four end-state criteria, one stable selector, negative controls, complete live lanes, and repository checks.

### Summary

The plan delivers one 16-journey suite through Claude, Codex, and Pi adapters. It binds registry IDs to source and makes drift visible through one operator-facing test surface.

The combined Pi estimate is $2.46, with a $3.08 approval ceiling for one complete run.
