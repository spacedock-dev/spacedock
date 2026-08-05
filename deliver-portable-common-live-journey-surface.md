---
title: Deliver one portable common live-journey surface
status: implementation
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
            - id: gate-attempt:ys7ncwh9kr8w5h9hdkz5apat-ideation-3
              briefing:
                id: briefing:ys7ncwh9kr8w5h9hdkz5apat:ideation:attempt-3:revision-1
                digest: sha256:8ceb3d3a6bff4e1abe4624b053b0f7ddf8ac9374bdbfffad03ddd36c4c418d8e
                request-digest: sha256:06aa372cedc1f9cd2f0edc5616e8c3a79f8007a6043e64c80c21a614b4ed6447
                room-ref: ./deliver-portable-common-live-journey-surface/review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:ys7ncwh9kr8w5h9hdkz5apat:ideation:3
                briefing: briefing:ys7ncwh9kr8w5h9hdkz5apat:ideation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-03T14:45:30.772745Z"
                decision: approve
                reason: Approved after staff review. Land the portable common journey surface after 3d and 15e.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-03T12:19:31Z
worktree: .worktrees/spacedock-ensign-deliver-portable-common-live-journey-surface
---

## Outcome

A test operator can select one common journey by one stable identity on Claude Sonnet, Claude Opus, Codex, or Pi. The registry remains the pure desired-state table; source reconciliation separately makes target-scoped owner failures, missing run evidence, and orphan fixtures visible.

The task delivers a complete and accurate representation of the desired journey table and the observed evidence state. Registry annotations, fixture bindings, runtime adapters, selectors, target-scoped owner TODOs, and reconciliation are implementation steps. A proven pass stays runnable; an unverified required target stays runnable and is reported as missing run evidence. Product behavior repair remains owned by the named repair members; this task neither hides those gaps nor absorbs their fixes.

Pi evidence is read from its archived root session JSONL. Correlated Pi tool calls and results, native runtime/model/token/cost metrics, process duration, and run provenance remain attributable to the Pi artifact; stdout and stderr remain diagnostics only.

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

Common runners own fixture selection, prompts, normalized evidence, and host-neutral assertions. Adapters own authentication, launch, runtime-native session parsing, liveness, metrics, and host artifacts.

Each source declaration carries its stable registry annotation. Reconciliation compares the pure desired registry with source, derives observed target-scoped gaps from TODO bindings, and fails ambiguous, orphaned, global, or wrongly owned bindings.

Every target selects all 16 journeys. Only a target with durable evidence of an owner-linked product failure emits `TODO(<repair-id>)`; reconciliation reports the exact journey, target, and owner as missing evidence. Proven passing and unverified required targets remain executable. A TODO is not passing evidence or a runtime exception. A coverage row, unowned/global skip, quarantine, or Pi-only top-level test does not count as evidence.

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
13. Run focused offline tests, three representative live runs, complete live lanes with exact TODO accounting, full tests, race tests, and formatting.

The adapter registry serves AC-1 and AC-4. Compatibility wrappers are insufficient because they preserve multiple authoritative identities.

The common runner map serves AC-2 and AC-4. Three runtime maps are insufficient because they permit private common journeys.

Stable source annotations serve AC-2 and AC-3. Builder names in the registry are insufficient because source moves can change desired state.

The reconciled inventory, exact target-scoped owner TODO set, desired-registry purity guard, and SHA guard serve AC-2 and AC-3. Semantic checks must live on the repository's contractlint boundary and prove values that can diverge from the documentation.

The Pi trace reader serves AC-4. Durable state alone cannot prove Pi child identity or required tool selection. A product-owned TODO records missing evidence; it does not prove the journey.

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

**AC-1 (VALUE) — One desired journey identity is selectable on every supported runtime.**

The same `^TestLiveSharedScenarios$/^shallow-boot$` selector passes on Claude, Codex, and Pi. Only runtime configuration changes.

The workflow guard proves that every live lane uses this top-level identity. A restored old suite name makes the guard fail.

**AC-2 (VALUE) — The complete desired table and target-scoped evidence state are represented truthfully.**

Registry, scenario-table, and runner-map parity is 16 of 16. The count of standalone common-journey top-level tests is zero.

The desired registry requires all 16 journeys on Claude Sonnet, Claude Opus, Codex, and Pi and contains no observed-gap ledger. Reconciliation reports 16 bound journey IDs and no missing suite lane, then derives each proven product gap from source as an exact journey+target+`TODO(<repair-id>)` result. Proven passes and unverified required targets remain runnable; reports distinguish them instead of inferring a failure or exception. Removing one runner or selector, hiding a target TODO, adding a global/unowned skip, suppressing a proven pass, changing target/owner/cardinality, or reintroducing an observed-gap registry section makes the proof fail.

The current source-derived set contains exactly six cells: the prior five Sonnet/Codex cells plus Pi `rejection-flow` → `TODO(zbcj98qfwtax61vxdzrf615e)`. The Pi binding records the audited incomplete two-entry round; it does not repair or pass the four-entry product contract.

**AC-3 — Fixtures have accountable use.**

Each common fixture ID resolves to one annotated builder and at least one journey. Proof-only and experiment fixtures resolve to their allowed registry class.

Reconciliation reports zero duplicate, invalid, orphan, unaccounted-test, and unaccounted-builder results. A temporary orphan fixture makes the proof fail.

**AC-4 — Runtime differences stay behind adapters and evidence remains attributable.**

Common fixtures and assertions contain no host-specific branch. The adapter factory contains the only runtime-selection branch.

Fake-adapter tests prove the boundary. Cheap controls enumerate the same 16 identities for Claude Sonnet, Claude Opus, Codex, and Pi and expose the exact target-scoped owner set without launching a model. Complete live lanes prove unchanged launch, liveness, metrics, and artifact behavior for runnable journeys. A TODO, an unverified target, or an externally blocked run does not count as passing evidence; each retains its distinct classification.

For Pi, rejection evidence requires an exact recorder command and a correlated successful `toolResult` reporting `entries=4` in the archived root session. Pi metrics require present provider/model, input/output/cache/total-token fields, total cost, duration, and run provenance; Claude attribution or partial usage fails the proof.

## Test plan

Run these focused offline tests before any model-backed run:

```bash
go test -tags live ./internal/ensigncycle -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestSharedLiveRuntimeSelection|TestPromotedCommonJourneyEntrypoints' -count=1
go test ./internal/ensigncycle -run 'TestGateJourneyOverlap|TestAutoContinue|TestZeroDiscover|TestEnsignCycleGoesRed|TestACReanchor' -count=1
go test ./internal/release -run 'TestRuntimeLiveWorkflow|TestWorkflowsPreserveAndPublishJourneyCosts' -count=1
go test -tags live ./internal/ensigncycle -run '^(TestRejectionFlowRoundInvocationExtractors|TestEmitPiScenarioMetricsUsesNativeSessionForAttributionAndUsage|TestParsePiSessionMetricsRequiresCompleteAttributedLargeRow|TestSharedLiveTODOEvidenceSet)$' -count=1
go test ./internal/contractlint -run '^(TestRuntimeLiveRegistryReconciliation|TestRuntimeLiveRegistryReconciliationMutationControls|TestRuntimeLiveRegistryReconciliationSHA)$' -count=1
```

The parity tests compare the registry IDs, scenario table, runner map, and adapter registry in both directions. Target enumeration separately exercises Claude Sonnet, Claude Opus, Codex, and Pi without a host launch and requires only exact durable owner failures to skip.

The negative tests remove one runner, add one orphan fixture, swap gate starts, stop auto-continue, select the wrong re-anchor branch, reintroduce an observed-gap registry ledger, and mutate missing-evidence target, owner, cardinality, or global scope. A Codex default-headless suppression mutant protects its proven pass.

Run one representative journey with the same selector on every runtime:

```bash
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_REQUIRED=1 go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
```

Run the six promoted selectors on Claude and Codex. Then run those same selectors on Pi.

Run each complete live lane through `TestLiveSharedScenarios`. Every required target must enumerate all 16 journeys. Runnable journeys must run for real. Skips are permitted only for exact target-scoped product-owned `TODO(<repair-id>)` bindings and must appear as journey+target+owner missing evidence in reconciliation. Proven passes remain runnable. Unverified targets remain runnable and are reported as missing run evidence, never as exceptions or inferred owner failures. An unowned/global skip, quarantine, runtime exception, coverage-only row, hidden selector, or desired-registry evidence ledger fails the proof.

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
- Pi changes from coverage metadata to the common executable selector, with real passes and owner-linked missing evidence distinguished.
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

The procedure also states that an owner-linked product TODO remains a desired common journey, does not create a runtime exception, and must be reported as missing evidence until its repair owner lands and the live assertion passes.

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

## Stage Report: implementation

- DONE: Start with focused failing tests for one TestLiveSharedScenarios identity, exactly 16 registry journeys, three runtime adapters, and exact workflow selectors.
  Focused parity, adapter-selection, TODO, and anchored-selector controls fail on an absent identity, adapter, runner, or lane.
- DONE: Bind existing source declarations with the four approved annotation forms and add reconciliation diagnostics with adversarial mutation controls.
  Registry-derived contractlint rejects structural gaps, missing evidence, duplicates, invalid/detached bindings, orphan tests, and unaccounted builders.
- DONE: Implement one TestLiveSharedScenarios entry point and one 16-journey runner map; do not add host-specific common-journey maps.
  Commit fc397422a contains the sole suite and sole executable common map in registry order.
- DONE: Refactor Claude and Codex common behavior behind transport adapters without changing their launch, liveness, artifact, or metrics contracts.
  Shared runner tests bind unchanged host transports to one common assertion/fixture surface.
- DONE: Promote the six standalone common journeys and remove their wrappers only after common-runner parity tests pass.
  Full cycle, gate stop, withdrawn recovery, zero discovery, auto-continue, and AC re-anchor now have no competing top-level wrappers.
- DONE: Implement the Pi adapter against the exact shallow-boot selector first; preserve the already-observed OpenRouter 402 as externally blocked evidence and do not retry the local Pi run.
  Pi failed pre-model in 2.714s on openrouter/openai/gpt-5.4 with HTTP 402 for a 128000-token request; no retry or full local lane ran.
- DONE: Make every runtime enumerate all 16 desired journeys. Implemented journeys run for real; product-owned gaps use only the exact owner-linked TODO set and remain missing evidence, never exceptions or passing evidence.
  Claude, Codex, and Pi cheaply enumerated the exact three TODOs; reconciliation reports MISSING=0 and MISSING-EVIDENCE=3.
- DONE: Restore the existing TODOs for smallest-sufficient-mechanism and keep-moving-posture and add TODO(9adv48yhye5s2vkhwd7ge52d) for auto-continue-after-implementation; keep all three in the registry, scenario table, runner map, and every adapter.
  All three selectors emit exact owner-linked skips; Claude shallow passed 40.095s, Codex shallow passed 37.541s, and prior Claude promoted evidence was 5 pass/auto-continue fail.
- DONE: Delete old top-level suite names, host-specific common maps, retired standalone wrappers, Pi coverage substitutes, and orphaned fixtures.
  AST/source reconciliation accounts for one suite, four proofs, two experiments, 16 journeys, and 21 adjacent fixture builders.
- DONE: Update all three workflow lanes to select TestLiveSharedScenarios through SPACEDOCK_LIVE_RUNTIME; update guards, metrics, docs, and architecture notes once.
  Every lane uses the anchored selector and retains its named metrics/artifact upload contract.
- DONE: Make reconciliation distinguish complete desired-state bindings from missing live evidence, list every exact TODO journey and repair owner, and reject hidden journeys, unowned skips, exceptions, duplicates, or orphans.
  Normative registry parsing reports the three exact 9adv48yhye5s2vkhwd7ge52d gaps and mutation controls red on each forbidden drift class.
- DONE: Fix the contractlint boundary: instruction/document semantic reads belong under internal/contractlint, not internal/release; retain the independent SHA/path guard without violating TestNoInstructionReadsOutsideQuarantine.
  Boundary guard and semantic reconciliation pass under internal/contractlint; the independent path guard remains separate.
- DONE: Run reconciliation, record the exact candidate SHA in a documentation-only commit, and prove both stale-SHA failure and current-SHA success.
  Core fc397422a3807573ea66a86eec30d376cf0b64d5 is recorded by docs-only 14f0e4147; the unskipped current/stale guard passes.
- DONE: Run focused offline controls and cheap exact selectors needed to prove TODO accounting. Preserve existing named Claude/Codex/Pi evidence; do not rerun product-owned failures or the locally blocked Pi lane. Leave complete paid lanes to relevant hosted CI.
  Focused controls and all adapter TODO selectors passed; no paid product-gap or blocked-Pi rerun occurred, and prior live results remain named above.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race, and keep the final worktree clean.
  Unskipped normal and race suites pass; the authorized origin/main baseline manifest rename e3aa5b753 was reconciled by exact FO-authorized path/count fixture updates.
- DONE: Report exact files and line counts within 42 files and +2750 insertions; do not change commands, stored formats, product runtime behavior, or FO/ensign authority.
  Final surface is 41 files, +1564/-654; changes are test, workflow, registry, documentation, and the authorized evidence-fixture reconciliation only.

### Summary

One portable 16-journey selector now runs through Claude, Codex, and Pi adapters while three product-owned gaps remain visible as exact missing evidence. Code commits fc397422a and 14f0e4147 are pushed; full normal/race checks and reconciliation are green, with paid completion left to hosted CI.

## Stage Report: validation

- DONE: Independently inspect the implementation report, candidate diff, commit history, and exact merge base; do not trust the implementer's summary alone.
  Inspected candidate `14f0e4147375b5dfa17d0d22d7f155557eb23ce2`, core `fc397422a3807573ea66a86eec30d376cf0b64d5`, origin/main merge base `507a4bc12e48a3e4a813219602c488f09c81a5d8`, and the full 41-file diff.
- DONE: Extract AC-1 through AC-4 from the entity and reproduce every cited code, command, or durable-state proof; reject self-referential or static prose-only evidence.
  Focused ensigncycle, release, contractlint, stale/current SHA, TODO-selector, full, and race executions reproduced AC-2/3 structure; AC-1/4 live proof failed in hosted run `30991534368`.
- DONE: Verify one authoritative TestLiveSharedScenarios entry point, exactly 16 registry-ordered common journey subtests, one common runner map, and exactly three runtime adapters.
  `go test -tags live ./internal/ensigncycle -list '^TestLive.*(Shared|Scenario|Journey)'` lists only `TestLiveSharedScenarios`; parity and runtime-selection tests pass for 16 ordered IDs and Claude/Codex/Pi.
- FAILED: Prove the same shallow-boot selector works with only SPACEDOCK_LIVE_RUNTIME changing across Claude, Codex, and Pi; accept the recorded local Pi OpenRouter 402 only as externally blocked evidence and require hosted direct-OpenAI Pi proof.
  GitHub run `30991534368` at exact head failed offline before model spend, so Claude, Codex, and direct-OpenAI Pi were all skipped and no attributable shallow-boot proof exists.
- DONE: Verify all three workflow lanes select TestLiveSharedScenarios and that restoring an old top-level suite name makes the independent workflow guard fail.
  `go test ./internal/release -run 'TestRuntimeLiveWorkflow|TestWorkflowsPreserveAndPublishJourneyCosts' -count=1 -v` passes exact three-lane selector and obsolete-selector mutation guards.
- FAILED: Reproduce bidirectional registry/scenario/runner/adapter parity and all adversarial removals: missing runner, orphan fixture, wrong gate start, stopped auto-continue, and wrong AC re-anchor branch.
  All named controls pass except wrong gate start: a detached mutation replacing both default-headless `writePreGateWorkflow` calls with already-gated `writeGateWorkflow` survived every focused ensigncycle/release/contractlint command.
- DONE: Verify every common fixture has exactly one approved source annotation/builder and accountable use; reproduce reconciliation with zero duplicate, invalid, orphan, unaccounted-test, unaccounted-builder, missing, or unselected results.
  `TestRuntimeLiveRegistryReconciliation` and its 24 mutation controls pass; each forbidden structural drift independently makes reconciliation return an error.
- DONE: Verify no standalone common-journey top-level wrapper, host-specific common runner map, Pi coverage substitute, Pi quarantine, or Pi-only common-suite substitute remains. Allow only the exact owner-linked product TODOs named by the latest entity; reject any other skip.
  Live test listing and reconciliation find one suite/map, three adapters, and only the three `TODO(9adv48yhye5s2vkhwd7ge52d)` skips on every runtime selector.
- DONE: Trace the changed identity, fixture ID, journey ID, runtime selection, artifacts, metrics, costs, and durable workflow evidence through every representation and lifecycle phase.
  Registry/source/workflow parity is green; hosted run `30991534368` demonstrates the evidence lifecycle stops before artifacts/metrics when the offline prerequisite fails.
- FAILED: Perform a semantic adversarial matrix covering unknown/empty runtime, all three adapters, missing runner/builder/selector, duplicate/orphan binding, failure cleanup, deadlines, and first-failure stop behavior.
  Unknown/empty, adapters, bindings, cleanup, and deadline controls pass; detached injected `YS_FIRST_FAILURE` still ran/logged `YS_SECOND_JOURNEY_RAN`, proving the suite ignores `t.Run` failure and does not stop.
- DONE: Inspect hot paths for multiplicative work, blocking I/O, unbounded reads/allocation, implicit size limits, and accidental model retries; add or run one scaling/over-limit control where risk exists.
  Full/race suites exercised streamwatch large-line, quiet-budget, cleanup, and single-run/no-retry controls; Pi retains its explicit 12-minute per-journey deadline but the first-failure defect defeats bounded suite spend.
- DONE: Verify Claude and Codex transport launch, liveness, artifact, and metrics contracts remain stable, and common fixtures/assertions contain no host-specific branch outside the adapter factory.
  Adapter/fake-driver, metrics workflow, streamwatch, host-neutrality, and full/race tests pass locally; model-backed transport confirmation remains blocked by the hosted prerequisite failure.
- FAILED: Verify Pi shallow-boot was attempted first, the local OpenRouter 402 caused no retry/full local lane, and hosted direct-OpenAI Pi proves adapter launch. Complete Pi evidence must enumerate all 16 identities, run implemented journeys for real, expose only the exact owner-linked TODO skips as missing evidence, and record runtime/model/duration/tokens/cost/artifacts.
  Prior local evidence is blocked OpenRouter `openai/gpt-5.4`, 2.714s, HTTP 402, zero retry/full lane; hosted Pi was skipped, producing no direct-OpenAI launch, 16-identity run, token/cost, or artifact evidence.
- DONE: Classify every local-first Claude/Codex/Sonnet/Pi lane as pass, fail, blocked, or irrelevant; keep Opus irrelevant unless changed surface or a concrete failure escalation implicates it.
  Cheap Claude/Codex/Pi TODO selectors pass structurally; complete Claude/Sonnet, Codex, and direct-OpenAI Pi lanes are blocked by candidate run `30991534368`; local Pi is externally blocked; Opus is irrelevant.
- DONE: Reproduce the stale reconciliation SHA failure and current candidate SHA success; verify the SHA-recording commit is documentation-only and names every changed watched path when stale.
  Local `TestRuntimeLiveRegistryReconciliationSHA` passes current `fc397422a` and rejects `fc397422a^` with all 29 watched paths; `14f0e4147` changes only `docs/runtime-live-ci.md`.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, plus the focused offline commands from the entity, and leave the candidate worktree clean.
  Formatting made no changes; all local focused, full (313.431s ensigncycle), and race (264.400s ensigncycle) commands pass; `git status --short` and `git diff --check` are empty.
- DONE: Measure exact changed files and added lines against the 42-file/+2750-insertion ceiling; flag any new product package, registry ID, authority command, grammar, stored-format, or product-runtime change.
  Exact origin/main diff is 41 files, +1564/-654; within ceiling, with no new product package, registry ID, authority command, grammar, stored format, or product-runtime change.
- DONE: Record each finding before candidate mutation with exact evidence and both defect-kind and release-scope axes; do not edit implementation code or authorize fixes.
  Finding 1 — evidence/harness defect, Material/candidate-owned: normal workflow is hosted `runtime-live-e2e`; harm is zero live evidence; `value-ac[AC-4]` requires complete attributable lanes; run `30991534368` fails `Invalid revision range fc397422a...HEAD` under shallow checkout and skips every live job.
  Finding 2 — outcome defect, Material/candidate-owned: normal workflow is the supported Pi common lane; harm is continued model work/cost after failed evidence; `contract[docs/dev/.spacedock-state/deliver-portable-common-live-journey-surface.md#Pi execution contract]` requires first-failure stop; detached two-runner execution proved the second ran after the first failed.
  Finding 3 — evidence defect, Material/candidate-owned: normal workflow is default headless drive-to-gate on all adapters; harm is a parked-at-gate fixture can falsely qualify without driving; `value-ac[AC-4]` requires implemented journeys to run for real; the wrong-start mutation survived all specified focused controls.
- DONE: Recommend PASSED only if all four latest ACs have reproducible evidence, all 16 desired identities and fixtures remain represented/reachable, implemented behavior has relevant passing evidence, and every product-owned gap is an exact repair-linked TODO reported as missing evidence rather than an exception/pass.
  Recommendation: REJECTED. AC-1/AC-4 lack hosted live proof and three Material candidate-owned findings remain; structural AC-2/AC-3 evidence and exact 9a TODO accounting are otherwise green.
- DONE: Write and push a path-scoped validation stage report with exact commands, results, evidence locations, candidate SHA, and PASSED/REJECTED recommendation.
  This report records local commands, GitHub run `30991534368`, detached adversarial observations, exact SHAs, and the REJECTED recommendation; state publication commit follows path-scoped.

### Summary

Candidate structure, reconciliation, focused controls, formatting, full tests, and race tests are locally green at `14f0e4147`. Validation recommends REJECTED because shallow hosted checkout breaks the SHA guard and suppresses every live lane, the common loop continues after first failure, and the required wrong-gate-start adversary is not guarded.

### Feedback Cycles

- Cycle 1: REJECTED — fresh Sol/medium validation plus hosted run `30991534368`; surface 41 files/+1564/-654 vs estimate 34 files/+2060/-1650 and ceiling 42 files/+2750 (57% of the insertion ceiling); AC unchanged. FO disposition: FIX three Material, ys-owned findings only. Route to implementation: require full history for the offline reconciliation guard; stop the common ordered suite when `t.Run` returns false; and centralize the host-neutral gate fixture choice with a negative `TestGateJourneyOverlap` proving `default-headless-gate-stop` starts at implementation while `gate-guardrail` starts at validation. Do not change product behavior, add a controller, absorb `9a`, or rerun hosted/live lanes before focused and full/race checks are green.
- Cycle 2: REJECTED — fresh Sol/medium validation at `ae2a19e8cfdedace6b4d20e32821490666016e63` plus exact-head hosted run `30996911834`; surface 41 files/+1661/-655 remains within the 42-file/+2750 ceiling; AC unchanged; all three Cycle-1 findings are closed. FO disposition: FIX exactly one Material, ys-owned evidence defect. Route to implementation: add a focused failing control proving the Codex and Pi shared drivers keep the scenario-local `spacedock` logger after the real front door re-pins `SPACEDOCK_BIN`, then propagate that override through the smallest existing host-specific driver seam. Preserve Claude behavior, launch shapes, first-failure stopping, the 16 journeys, and the exact three `9adv48yhye5s2vkhwd7ge52d` TODOs. The Sonnet duplicate-preparation failure is Material/product-owned and Needs decision; the Sonnet break-glass failure and Codex characterized-metrics limitation are adjacent/pre-existing. No ys product fix, new member, controller, retry, or extra lane is authorized. Commit watched-path code first, bind its SHA in one documentation-only commit, run focused/full/race checks, then return to the same independent validator. This is correction round 2; any next rejection is Cycle 3 and escalates to the Captain rather than auto-bouncing again.
- Cycle 3: REJECTED — fresh independent Sol/medium validation at `37920950356b43fcc2f5426ab5ae9f9605a2ac03`; surface 41 files/+1724/-662 remains within the 42-file/+2750 ceiling; all candidate-owned logger, first-failure, shallow-history, and fixture-start defects are closed, with focused/full/race checks green. Captain disposition after escalation: FIX exactly one representation defect before paid evidence. Add `default-headless-gate-stop` as `TODO(26nk8qd48zknqnn4kc123sez)` and report it as the fourth MISSING-EVIDENCE row because exact-head Sonnet run `30996911834` proved duplicate durable preparation and 26n owns that outcome. Test the exact owner/count first; preserve all 16 identities, the three existing 9a TODOs, one runner map, three adapters, and every product/runtime path. Do not repair 26n, change substrate/smoke, add a member/controller/retry, or launch paid evidence. Commit watched-path representation changes first, bind that SHA in one documentation-only commit, run focused/full/race checks, then return to the same independent validator. After validation passes, run only the shared common selectors for Pi, Codex, and Sonnet; Opus remains irrelevant.
- Cycle 4: REJECTED — fresh independent Sol/medium validation at `70a61f9e36b86ad260aaf28eba506b5b3ac81a30`; surface 41 files/+1757/-663 remains within the 42-file/+2750 ceiling; the exact four journey/owner cardinality is structurally green, but the mechanism skips every target and the desired-state registry contains a running evidence ledger. Captain accuracy correction: FIX the ys representation before any spend. Remove the normative `Missing live evidence` section from `docs/runtime-live-ci-registry.md` while keeping all journeys and all four desired targets required. Make source TODO bindings and reconciliation target-aware: 26n withholds only the proven Sonnet `default-headless-gate-stop` evidence from run `30996911834`, preserving Codex's pass and allowing unverified Pi to measure. Audit each 9a TODO against exact durable target evidence; bind only proven failing targets, preserve proven passes, and classify unrun targets as unverified/missing in evidence reports rather than exceptions or inferred owner failures. Opus remains a required target and, unless run, must be reported unverified/missing—not irrelevant or excepted. Add red-first registry-purity and target-scope controls, including global-skip/wrong-target/pass-suppression mutations. Do not repair 26n/9a, add product behavior, or launch paid evidence. Commit watched-path changes first, bind the exact SHA in one docs-only commit, run focused/full/race checks, then return to the same independent validator before any hosted approval.
- Cycle 5: REJECTED — fresh independent Sol/medium validation at `e023619ed1f6a00a85f6ad88741c8b8ebad2d3b1`; surface 40 files/+1870/-664 remains within the 42-file/+2750 ceiling; target-aware registry purity, five-cell bindings, focused controls, full, and race are green. Exact-head Codex run `31015945064` is external-capacity INCONCLUSIVE after two passes. Exact-head Pi run `31016570689` passes five journeys, then exposes two Material findings at `rejection-flow`: the ys adapter reads the wrong observation stream and emits Claude-attributed/incomplete Pi metrics; separately, Pi records the correction round early with two entries rather than the complete four-entry reviewer+worker log. Captain disposition: FIX the ys-owned observation/metrics boundary and represent only Pi `rejection-flow` as `TODO(zbcj98qfwtax61vxdzrf615e)`, the existing cross-runtime correction-round publication owner. Add red-first Pi session correlation, false-positive/false-negative, metric runtime/model/token/cost, and exact six-cell target-binding controls. Preserve the pure desired registry, every journey/target identity, all five existing bindings, product behavior, and Claude/Codex launch and evidence behavior. Do not repair zbc, add a member/controller/retry, or run paid evidence until focused/full/race and the watched-core/docs-binder shape are green; then return to the same independent validator before resuming isolated Codex, Pi, and Sonnet common evidence. Opus remains required/unverified.
- Cycle 6: REJECTED — fresh independent Sol/medium validation at `35143a55a3fd1259295664df71602c116f6b3fd3`; surface 42 files/+2173/-666 is at the file ceiling and below +2750 insertions; strict native Pi metrics, the historical `entries=2` rejection, exact six-cell target bindings, focused controls, full, and race are green. One Material ys-owned AC-4 evidence defect remains: `piRecordedRejectionRound` accepts ambiguous repeated exact recorder execution, including two successful `entries=4` calls and an `entries=4` call followed by a second exact `entries=2` call. Captain disposition: FIX only the Pi observer cardinality/correlation seam so exactly one exact invocation plus exactly one correlated successful `entries=4` result qualifies; duplicate invocations, repeated results, reused call IDs, and mixed complete/incomplete order fail closed. Preserve Pi metrics, registry purity, all six bindings, product bytes, Claude/Codex behavior, the surface ceiling, and the execution hold. Commit watched core then one docs-only binder, run focused/full/race, and return to the same validator before any hosted spend.
- Cycle 7: REJECTED — fresh independent Sol/medium validation at `e952ed65c17d6f2336678680f672222d33a83cd4`; surface 42 files/+2256/-666 remains at the file ceiling and below +2750 insertions; duplicate calls/results, ordering, ID reuse, duplicate `entries=4`, archived `entries=2`, native metrics, exact six bindings, focused controls, full, and race are green. One Material ys-owned AC-4 evidence defect remains: the sole correlated Pi result can contain both canonical `entries=2` and `entries=4` summaries and still qualify, in one text block or separate blocks. Captain disposition: FIX only result-content cardinality. Count every canonical rejection-task validation/1 summary independent of entry count and accept only exactly one summary whose complete line is `entries=4`; fail mixed or repeated summaries while preserving unrelated diagnostic text and every prior Cycle-7 control. Preserve product bytes, metrics, bindings, history/ceiling, and the execution hold; return to the same validator before any hosted spend.

## Stage Report: implementation (cycle 2)

- DONE: Add a focused failing workflow mutation test proving the offline Runtime Live E2E checkout provides the recorded reconciliation commit and its parent.
  RED: `go test ./internal/release -run '^TestRuntimeLiveOfflineCheckoutFetchesRecordedCommitHistory$'` failed before the guard existed; GREEN: the same command passed and its depth-one mutation was rejected.
- DONE: Set only the offline job's actions/checkout step to fetch-depth 0; keep the SHA guard algorithm and live-job checkouts unchanged.
  `.github/workflows/runtime-live-e2e.yml` changes only the offline checkout; the existing SHA algorithm and four later checkout steps are unchanged.
- DONE: Add a focused failing ordered-sequence test proving a false first journey prevents every later journey from running.
  RED: `go test -tags live ./internal/ensigncycle -run '^TestSharedScenarioSequenceStopsAfterFirstFailure$'` failed before the seam existed; GREEN: it passed with only `first` observed.
- DONE: Route TestLiveSharedScenarios through the smallest host-neutral sequence seam that consumes t.Run's boolean and stops on false; add no controller or retry loop.
  `runSharedScenarioSequence` returns immediately on false and owns no runtime, retry, or orchestration state.
- DONE: Add TestGateJourneyOverlap as a focused failing negative proving default-headless-gate-stop starts at implementation and gate-guardrail starts at validation, and that swapping either start fails.
  RED: `go test ./internal/ensigncycle -run '^TestGateJourneyOverlap$'` failed before its assertion existed; GREEN: it passed against both real selected fixtures and rejects both crossed states.
- DONE: Centralize only the gate fixture choice in one host-neutral seam and make Claude/Codex use it; preserve Pi's existing reuse and all launch/assertion behavior.
  Claude and Codex now call `writeGateJourneyFixture`; Pi and every post-selection launch/assertion path are unchanged.
- DONE: Run the exact focused tests red-before/green-after, then gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, SHA current/stale controls, and git diff --check.
  All focused groups passed; post-commit normal passed (`ensigncycle` 207.136s), race passed (`ensigncycle` 222.664s), SHA and diff controls passed.
- DONE: Update the reconciliation SHA correctly: commit watched-path fixes first, then bind that exact commit in one documentation-only commit, proving current success and stale diagnostics from full history.
  Watched-path commit is `5bb44e7b4d7bf3092b738a74d6ede04f4b0a5aed`; docs-only commit `ae2a19e8cfdedace6b4d20e32821490666016e63` records it. Current succeeds; stale names all six changed watched paths.
- DONE: Keep all 16 journeys, exact three 9a-owned TODOs, one runner map, three adapters, and 41-file/+1564 baseline semantics unchanged except for these three authorized fixes.
  Focused parity/TODO controls pass; final origin/main surface remains 41 files at +1661/-655, within the 42-file/+2750 ceiling and differing from baseline only by this correction.
- DONE: Push the code branch and append/push a path-scoped implementation cycle-2 report with exact commands/results and final surface; do not run hosted/live CI in implementation.
  Branch head `ae2a19e8cfdedace6b4d20e32821490666016e63` is published; no paid model-backed, hosted, or live command ran in this correction cycle.

### Summary

The three authorized evidence defects are corrected and locally verified: offline reconciliation has full Git history, the common suite stops after its first failed journey, and gate journeys select and prove distinct start states. The pushed two-commit shape preserves a watched-path reconciliation commit followed by a documentation-only SHA binding.

## Stage Report: validation (cycle 2)

- DONE: Independently inspect the implementation report, candidate diff, commit history, and exact merge base; do not trust the implementer's summary alone.
  Inspected rejected `14f0e4147` through corrected core `5bb44e7b4d7bf3092b738a74d6ede04f4b0a5aed` and docs-only head `ae2a19e8cfdedace6b4d20e32821490666016e63`; origin/main merge base remains `507a4bc12e48a3e4a813219602c488f09c81a5d8`.
- DONE: Extract AC-1 through AC-4 from the entity and reproduce every cited code, command, or durable-state proof; reject self-referential or static prose-only evidence.
  Reproduced local structure, mutation, SHA, full/race, and exact-head hosted run `30996911834`; AC-2/3 are green, but real AC-1/4 lanes fail before complete evidence.
- DONE: Verify one authoritative TestLiveSharedScenarios entry point, exactly 16 registry-ordered common journey subtests, one common runner map, and exactly three runtime adapters.
  Live listing exposes only `TestLiveSharedScenarios`; parity tests prove the exact ordered 16, one runner map, and Claude/Codex/Pi selection including empty/unknown rejection.
- FAILED: Prove the same shallow-boot selector works with only SPACEDOCK_LIVE_RUNTIME changing across Claude, Codex, and Pi; accept the recorded local Pi OpenRouter 402 only as externally blocked evidence and require hosted direct-OpenAI Pi proof.
  Hosted Sonnet, Codex, and direct-OpenAI Pi each failed before shallow-boot; no exact-head three-runtime shallow selector proof exists.
- DONE: Verify all three workflow lanes select TestLiveSharedScenarios and that restoring an old top-level suite name makes the independent workflow guard fail.
  Release guards pass exact selector checks and mutations; hosted jobs invoked the one selector with only runtime/model environment differences.
- DONE: Reproduce bidirectional registry/scenario/runner/adapter parity and all adversarial removals: missing runner, orphan fixture, wrong gate start, stopped auto-continue, and wrong AC re-anchor branch.
  Reconciliation’s 24 mutants, stopped-auto-continue, AC re-anchor branches, and `TestGateJourneyOverlap` pass; detached swapped gate fixtures fail with the wrong required start.
- DONE: Verify every common fixture has exactly one approved source annotation/builder and accountable use; reproduce reconciliation with zero duplicate, invalid, orphan, unaccounted-test, unaccounted-builder, missing, or unselected results.
  Semantic reconciliation and all duplicate/invalid/orphan/unaccounted/missing/unselected mutation controls pass at `ae2a19e8`.
- DONE: Verify no standalone common-journey top-level wrapper, host-specific common runner map, Pi coverage substitute, Pi quarantine, or Pi-only common-suite substitute remains. Allow only the exact owner-linked product TODOs named by the latest entity; reject any other skip.
  Listing/reconciliation find one suite/map and exactly the three `TODO(9adv48yhye5s2vkhwd7ge52d)` gaps; all three cheap runtime selectors emit only those skips.
- DONE: Trace the changed identity, fixture ID, journey ID, runtime selection, artifacts, metrics, costs, and durable workflow evidence through every representation and lifecycle phase.
  Local parity is complete; hosted artifacts `8927141210` Sonnet, `8926705099` Codex, and `8926752230` Pi preserve exact-head execution and first-failure evidence.
- DONE: Perform a semantic adversarial matrix covering unknown/empty runtime, all three adapters, missing runner/builder/selector, duplicate/orphan binding, failure cleanup, deadlines, and first-failure stop behavior.
  All local matrix controls pass; detached `return`→`continue` fails with `[first second]`, and hosted Codex/Pi/Sonnet each stop immediately after their first common-journey failure.
- DONE: Inspect hot paths for multiplicative work, blocking I/O, unbounded reads/allocation, implicit size limits, and accidental model retries; add or run one scaling/over-limit control where risk exists.
  Full/race plus streamwatch large-line, quiet-timeout, cleanup, and single-run controls pass; hosted first-failure stopping bounds subsequent model work with no retry.
- FAILED: Verify Claude and Codex transport launch, liveness, artifact, and metrics contracts remain stable, and common fixtures/assertions contain no host-specific branch outside the adapter factory.
  Sonnet launches but duplicates gate preparation; Codex launches and mutates correctly but its shared-driver shim loses `command.log`; common static host-neutrality remains green.
- FAILED: Verify Pi shallow-boot was attempted first, the local OpenRouter 402 caused no retry/full local lane, and hosted direct-OpenAI Pi proves adapter launch. Complete Pi evidence must enumerate all 16 identities, run implemented journeys for real, expose only the exact owner-linked TODO skips as missing evidence, and record runtime/model/duration/tokens/cost/artifacts.
  Direct OpenAI `openai/gpt-5.4` launched and full-ensign-cycle passed, then gate-guardrail failed at 102.18s on missing `command.log`; the suite stopped, smoke skipped, and only two identities ran.
- DONE: Classify every local-first Claude/Codex/Sonnet/Pi lane as pass, fail, blocked, or irrelevant; keep Opus irrelevant unless changed surface or a concrete failure escalation implicates it.
  Offline PASS; Sonnet FAIL; Codex FAIL; direct-OpenAI Pi FAIL; prior local OpenRouter Pi BLOCKED 402; Opus IRRELEVANT and rejected before spend; no lane is misreported as pass.
- DONE: Reproduce the stale reconciliation SHA failure and current candidate SHA success; verify the SHA-recording commit is documentation-only and names every changed watched path when stale.
  A real full clone resolves `5bb44e7b^` and passes; a real depth-one clone fails `Invalid revision range`; stale diagnostics name all six watched paths, current passes, and `ae2a19e8` changes only `docs/runtime-live-ci.md`.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, plus the focused offline commands from the entity, and leave the candidate worktree clean.
  Formatting changed nothing; focused/full (`ensigncycle` 200.996s)/race (`ensigncycle` 207.931s) all pass; worktree status and `git diff --check` are empty.
- DONE: Measure exact changed files and added lines against the 42-file/+2750-insertion ceiling; flag any new product package, registry ID, authority command, grammar, stored-format, or product-runtime change.
  Exact origin/main surface is 41 files, +1661/-655, within ceiling; Cycle 2 adds no package, ID, authority, grammar, stored format, or product-runtime behavior.
- DONE: Record each finding before candidate mutation with exact evidence and both defect-kind and release-scope axes; do not edit implementation code or authorize fixes.
  Finding 1 — evidence defect, Material/candidate-owned: normal Codex/Pi common lanes lose command attribution; harm is failed implemented journeys/incomplete 16; `value-ac[AC-4]` requires stable attributable artifacts; jobs `92276320245`/`92276320256` show successful commands followed by absent `command.log` because shared drivers preserve direct `SPACEDOCK_BIN`.
  Finding 2 — outcome defect, Material/product-owned and Needs decision: normal Sonnet default-headless drive creates duplicate durable gate attempts; harm is two preparations and failed journey; `value-ac[AC-4]` requires implemented journeys run for real; job `92276320289` shows a rejected three-reference prepare, successful `briefing-1`, explicit withdrawal, then successful `briefing-2` before `prepared fixture request count = 2, want 1`.
  Adjacent finding — outcome defect, Material/pre-existing runtime-product scope: normal Sonnet break-glass substrate used unnamed foreground Agent and failed its stable proof; `value-ac[AC-4]` requires unchanged launch behavior; job `92276320289` reports the mismatched Agent shape, while the candidate diff does not touch that substrate.
  Adjacent evidence limitation — evidence defect, Needs decision/pre-existing external-schema scope: normal Codex records have empty parsed model and no verified token/cost fields; harm is incomplete single-record attribution; `none:` AC-4 requires unchanged metrics and only Pi explicitly requires per-journey model/tokens/cost; current `thread.started`/nested-item schema triggers it while the emitter/parser candidate diff is empty.
- DONE: Recommend PASSED only if all four latest ACs have reproducible evidence, all 16 desired identities and fixtures remain represented/reachable, implemented behavior has relevant passing evidence, and every product-owned gap is an exact repair-linked TODO reported as missing evidence rather than an exception/pass.
  Recommendation: REJECTED. All three Cycle-1 fixes are proven, but AC-1/4 lack complete live proof, one candidate-owned Material evidence defect remains, and Sonnet exposes a non-9a product-owned Material gap.
- DONE: Write and push a path-scoped validation stage report with exact commands, results, evidence locations, candidate SHA, and PASSED/REJECTED recommendation.
  This cycle-2 report records exact SHAs, local commands, hosted run/jobs/artifacts, classifications, and REJECTED recommendation; publication commit follows path-scoped.

### Summary

Cycle 2 independently proves all three prior fixes in both directions and keeps AC-2/3, formatting, full, and race evidence green. Validation remains REJECTED: exact-head hosted evidence stops early on one candidate-owned Codex/Pi recording-shim defect and a separate Sonnet duplicate-preparation product failure, so the required complete 16-identity AC-1/4 proof does not exist.

## Stage Report: implementation (cycle 3)

- DONE: Add a focused failing offline live-tag control that simulates the real front door re-pinning SPACEDOCK_BIN and proves both Codex and Pi shared drivers still resolve the scenario-local spacedock logging shim.
  RED: `go test -tags live ./internal/ensigncycle -run '^TestSharedCodexAndPiDriversPreserveSpacedockShimAfterFrontDoorPin$' -count=1` failed because the shared seam accepted no test context and only prefixed PATH; GREEN: both host subtests now resolve the scenario shim after a simulated `/real/spacedock` pin.
- DONE: Correct only the Codex and Pi shared-driver shim propagation through the smallest existing host-specific mechanisms; preserve Claude behavior and every host launch shape.
  Codex delegates to its existing host wrapper, and Pi adds the existing shell-startup override after PATH prefixing. Claude retains its prior PATH plus caller-owned shell propagation, with only the seam signature carrying `*testing.T`.
- DONE: Prove both host cases fail when their propagation is removed and pass with the correction; do not add a controller, retry, runtime lane, or product behavior.
  The focused control removes the Codex wrapper and Pi `BASH_ENV`/`ZDOTDIR` propagation independently and observes the real front-door pin in both cases; no production package or orchestration behavior changed.
- DONE: Run the focused live-tag controls, adapter/parity/TODO/reconciliation/SHA controls, gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, and git diff --check.
  Focused host and parity controls pass; semantic reconciliation passes; normal passed (`ensigncycle` 194.089s), race passed (`ensigncycle` 279.265s), and formatting/diff checks are clean.
- DONE: Commit watched-path code first, then bind that exact reconciliation SHA in one documentation-only commit and prove current success plus stale diagnostics from full history.
  Watched-path commit `1e64b514eedf0f252be01c767ab28fd0e75dfcbc` is recorded by docs-only commit `37920950356b43fcc2f5426ab5ae9f9605a2ac03`; current passes and stale diagnostics name all five changed watched paths.
- DONE: Keep one suite, one 16-entry runner map, three adapters, exactly three 9a-owned TODOs, and the approved 42-file/+2750 ceiling; push the code branch and append/push a path-scoped implementation cycle-3 report without running hosted/model-backed CI.
  Parity/TODO controls remain green; final origin/main surface is 41 files at +1724/-662. The branch is pushed, no paid or hosted run occurred, and the adjacent Sonnet/metrics findings were not changed.

### Summary

Codex and Pi now preserve the scenario-local logging shim across their real front-door launcher pin while Claude and every launch shape remain unchanged. The focused removal controls, full/race suites, reconciliation, and two-commit SHA binding are green and published for independent revalidation.

## Stage Report: validation (cycle 3)

- DONE: Independently inspect the implementation report, candidate diff, commit history, and exact merge base; do not trust the implementer's summary alone.
  Candidate `37920950356b43fcc2f5426ab5ae9f9605a2ac03`, watched core `1e64b514eedf0f252be01c767ab28fd0e75dfcbc`, merge base `507a4bc12e48a3e4a813219602c488f09c81a5d8`; the head changes only `docs/runtime-live-ci.md`.
- DONE: Extract AC-1 through AC-4 from the entity and reproduce every cited code, command, or durable-state proof; reject self-referential or static prose-only evidence.
  Exercised parity, selector, reconciliation, workflow, stale-SHA, shim, and lifecycle controls; AC-2/AC-4 remain blocked by the missing 26n representation described below.
- DONE: Verify one authoritative TestLiveSharedScenarios entry point, exactly 16 registry-ordered common journey subtests, one common runner map, and exactly three runtime adapters.
  `TestSharedScenarioRunnerCoverageFinal`, `TestSharedRuntimeScenarioDefinitions`, and `TestSharedLiveRuntimeSelection` pass in both directions for 16 journeys and Claude/Codex/Pi.
- SKIPPED: Prove the same shallow-boot selector works with only SPACEDOCK_LIVE_RUNTIME changing across Claude, Codex, and Pi; accept the recorded local Pi OpenRouter 402 only as externally blocked evidence and require hosted direct-OpenAI Pi proof.
  Captain required the fourth owner-linked TODO before paid evidence, so no Cycle-3 model-backed selector was launched.
- DONE: Verify all three workflow lanes select TestLiveSharedScenarios and that restoring an old top-level suite name makes the independent workflow guard fail.
  `TestRuntimeLiveWorkflowUsesExactCommonSelectors` and `TestRuntimeLiveWorkflowGuardRejectsMissingSharedScenarioRun` pass; three executable selectors and no retired selector remain.
- DONE: Reproduce bidirectional registry/scenario/runner/adapter parity and all adversarial removals: missing runner, orphan fixture, wrong gate start, stopped auto-continue, wrong AC re-anchor branch, and Codex/Pi front-door shim propagation.
  Focused controls pass; detached removals expose `/real/spacedock` for each host, `continue` runs `[first second]`, and swapped gate starts fail `TestGateJourneyOverlap`.
- DONE: Verify every common fixture has exactly one approved source annotation/builder and accountable use; reproduce reconciliation with zero duplicate, invalid, orphan, unaccounted-test, unaccounted-builder, missing, or unselected results.
  `TestRuntimeLiveRegistryReconciliation` plus all 24 mutation subtests pass, including duplicate, invalid, orphan, missing, TODO-owner, target, and suite mutations.
- FAILED: Verify no standalone common-journey top-level wrapper, host-specific common runner map, Pi coverage substitute, Pi quarantine, or Pi-only common-suite substitute remains. Allow only the exact owner-linked product TODOs named by the latest entity; reject any other skip.
  Structural exclusions pass, but `default-headless-gate-stop` must join the three 9a gaps as `TODO(26nk8qd48zknqnn4kc123sez)`/MISSING-EVIDENCE before another paid run.
- DONE: Trace the changed identity, fixture ID, journey ID, runtime selection, artifacts, metrics, costs, and durable workflow evidence through every representation and lifecycle phase.
  Registry annotations, the sole runner map, adapter factory, three workflow jobs, artifact uploads, and cost guards preserve stable IDs; the 26n outcome is not truthfully represented yet.
- DONE: Perform a semantic adversarial matrix covering unknown/empty runtime, all three adapters, missing runner/builder/selector, duplicate/orphan binding, failure cleanup, deadlines, first-failure stop behavior, and removal of either host's shim propagation.
  Empty/unknown runtime, parity/reconciliation mutations, cleanup, quiet/deadline controls, first-failure mutation, and both independent shim removals all fail or pass at their intended boundaries.
- DONE: Inspect hot paths for multiplicative work, blocking I/O, unbounded reads/allocation, implicit size limits, and accidental model retries; add or run one scaling/over-limit control where risk exists.
  `TestNoTimeoutLiteralExceeds60s`, process-deadline, streamwatch activity/quiet-timeout, cleanup, and Pi evidence-grade controls pass; the common sequence remains linear and stops on first failure.
- DONE: Verify Claude and Codex transport launch, liveness, artifact, and metrics contracts remain stable, and common fixtures/assertions contain no host-specific branch outside the adapter factory.
  Diff inspection shows only the existing Codex wrapper delegation and Pi shell override; transport/cost controls pass and runtime selection remains solely in `selectSharedLiveRuntime`.
- SKIPPED: Verify Pi shallow-boot was attempted first, the local OpenRouter 402 caused no retry/full local lane, and hosted direct-OpenAI Pi proves adapter launch. Complete Pi evidence must enumerate all 16 identities, run implemented journeys for real, expose only the exact owner-linked TODO skips as missing evidence, and record runtime/model/duration/tokens/cost/artifacts.
  Prior local 402 evidence was not retried; Captain barred paid evidence until the fourth TODO is represented, so no direct-OpenAI Cycle-3 run exists.
- DONE: Classify every local-first Claude/Codex/Sonnet/Pi lane as pass, fail, blocked, or irrelevant; keep Opus irrelevant unless changed surface or a concrete failure escalation implicates it.
  Deterministic Codex/Pi shim controls pass; Sonnet/Codex/Pi paid lanes are blocked pre-run by the 26n representation gate; Claude and Opus are irrelevant to the authorized Cycle-3 hosted scope.
- DONE: Reproduce the stale reconciliation SHA failure and current candidate SHA success; verify the SHA-recording commit is documentation-only and names every changed watched path when stale.
  Full-history clone resolves `1e64b514^` and passes; depth-one fails the revision range; current/stale guard passes and binder `379209503` changes one documentation line.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, plus the focused offline commands from the entity, and leave the candidate worktree clean.
  Formatting/diff are clean; normal passes (`ensigncycle` 218.690s), race passes (`ensigncycle` 251.682s), and all focused commands pass at exact head/origin.
- DONE: Measure exact changed files and added lines against the 42-file/+2750-insertion ceiling; flag any new product package, registry ID, authority command, grammar, stored-format, or product-runtime change.
  Exact merge-base surface is 41 files, +1724/-662; no new product package, registry ID, authority command, grammar, stored format, or product-runtime change appears.
- DONE: Record each finding before candidate mutation with exact evidence and both defect-kind and release-scope axes; do not edit implementation code or authorize fixes.
  Material evidence defect/Needs decision: normal `default-headless-gate-stop` duplicated durable prepare in run `30996911834`; `value-ac[AC-2]` requires truthful missing evidence, Captain matched owner 26n, and candidate stayed untouched.
- FAILED: Recommend PASSED only if all four latest ACs have reproducible evidence, all 16 desired identities and fixtures remain represented/reachable, implemented behavior has relevant passing evidence, and every product-owned gap is an exact repair-linked TODO reported as missing evidence rather than an exception/pass.
  Recommend REJECTED: add the single fourth registry/common-runner TODO for 26n before paid evidence; Cycle-3 escalation goes to the Captain, not an automatic implementation bounce.
- DONE: Write and push a path-scoped validation cycle-3 report with exact commands, results, evidence locations, candidate SHA, and PASSED/REJECTED recommendation.
  This report records exact SHAs, reproducible commands and mutation outcomes, prior hosted run ID, classifications, and the REJECTED recommendation; publication commit follows path-scoped.

### Summary

Cycle 3 independently closes the ys-owned Codex/Pi shim defect in both directions and keeps all deterministic, full, race, reconciliation, and ceiling evidence green. Validation is REJECTED before paid evidence because the known default-headless duplicate-prepare outcome must be represented as owner-linked `TODO(26nk8qd48zknqnn4kc123sez)`/MISSING-EVIDENCE; the candidate was not mutated.

## Stage Report: implementation (cycle 4)

- DONE: Add a focused failing control that requires default-headless-gate-stop to map to exact owner TODO(26nk8qd48zknqnn4kc123sez), keeps the three existing TODO(9adv48yhye5s2vkhwd7ge52d) mappings, and reports exactly four missing-evidence journeys.
  RED: `TestSharedLiveTODOEvidenceSet` reported an empty reason for `default-headless-gate-stop`; GREEN: it verifies the exact 26n owner plus the unchanged three 9a mappings and rejects every extra TODO.
- DONE: Implement only the representation change: map default-headless-gate-stop to TODO(26nk8qd48zknqnn4kc123sez) through the existing common TODO seam, preserving all 16 identities, one runner map, three adapters, and every fixture/launch/assertion path.
  Only `liveDurableJourneyTODO` gains the 26n mapping; parity remains green and no fixture, launcher, adapter, or assertion implementation changed.
- DONE: Update the normative registry, operating guide, reconciliation expectations, and focused TODO controls so the exact four journey-to-owner pairs are machine checked; do not repair or alter 26n product behavior.
  The registry adds one row, the guide reports `MISSING-EVIDENCE=4`, and reconciliation now resolves each TODO case to its declared owner constant.
- DONE: Prove every runtime cheaply enumerates the same exact four TODO skips and that wrong owner/count/removal mutations fail; do not run any model-backed or hosted lane during implementation.
  Claude, Codex, and Pi exact selectors each emitted four zero-cost skips; dedicated default-headless wrong-owner, duplicate-row, removal, plus generic TODO mutations all passed by rejecting their adversary.
- DONE: Run gofmt -w ./cmd ./internal, focused parity/TODO/reconciliation/workflow controls, go test ./..., go test ./... -race, current/stale SHA controls, git diff --check, and leave the worktree clean.
  Formatting/focused/diff checks pass; normal passed (`ensigncycle` 265.958s), race passed (`ensigncycle` 273.636s), and current/stale reconciliation is green.
- DONE: Commit all watched-path representation changes first, then bind that exact core SHA in one documentation-only commit; push the branch and report both SHAs plus exact diff size against the 42-file/+2750 ceiling.
  Representation commit `0c66a7babe175f005d24d1b936d7e64fcd4f350a` is bound by docs-only `70a61f9e36b86ad260aaf28eba506b5b3ac81a30`; final surface is 41 files at +1757/-663 and the branch is pushed.
- DONE: Append and push a path-scoped implementation cycle-4 report; do not change substrate/smoke/Opus, add a controller/retry/member, absorb 26n/9a/rm, or expand product/runtime behavior.
  This report records the representation-only result; no substrate, smoke, Opus, product repair, controller, retry, member, or model-backed/hosted execution entered the correction.

### Summary

All 16 common journeys remain selected while every runtime now truthfully exposes the exact four owner-linked evidence gaps: one 26n default-headless gap and three 9a gaps. Machine reconciliation, mutation controls, full/race suites, and the two-commit SHA binding are green and published for independent validation.

## Stage Report: validation (cycle 4)

- DONE: Independently inspect candidate 70a61f9e36b86ad260aaf28eba506b5b3ac81a30, watched core 0c66a7babe175f005d24d1b936d7e64fcd4f350a, reports/history, merge base, and exact diff; do not trust the implementation summary.
  Exact merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`; five representation files change in core, the binder changes one docs line, and the complete surface is 41 files/+1757/-663.
- FAILED: Reproduce exact 16-journey/one-map/three-adapter parity and exact four TODO mappings: default-headless-gate-stop→26nk8qd48zknqnn4kc123sez and auto-continue-after-implementation/smallest-sufficient-mechanism/keep-moving-posture→9adv48yhye5s2vkhwd7ge52d.
  Cardinality/owners pass, but cheap exact selectors prove every adapter skips all four: 26n is evidenced only on Sonnet, while 9a names Sonnet/Codex evidence and no Pi or Opus trigger, so runtime-global mappings are semantically false.
- FAILED: Run bidirectional wrong-owner, duplicate/count, removal, hidden-skip, runner/adapter/selector, fixture, and reconciliation mutations; confirm the new 26n row is MISSING-EVIDENCE, not pass/exception, and that no product/runtime behavior changed.
  All 27 built-in mutations pass and detached owner/removal changes fail with `MISSING-EVIDENCE owner`/`UNOWNED-EVIDENCE`; however the mechanism cannot mutate or validate target scope, and the desired registry improperly contains an observed-evidence table.
- SKIPPED: Reprove the prior shallow-history, first-failure, gate-start, and Codex/Pi SPACEDOCK_BIN shim controls in both directions, plus formatting, focused workflow/SHA controls, go test ./..., go test ./... -race, diff check, ceiling, and clean exact origin head.
  Prior controls, formatting, workflow/SHA, diff, ceiling, origin, and normal suite pass (`ensigncycle` 228.689s); Captain stopped validation before race after the target-scope accuracy correction.
- DONE: Send the First Officer a concise local checkpoint before any hosted/model-backed evidence. Do not launch, approve, reject, cancel, or rerun GitHub jobs yourself; the First Officer will supply exact common-selector run/job/artifact IDs after the checkpoint.
  Sent the local checkpoint while normal tests ran; no GitHub or model-backed action was taken.
- SKIPPED: After receiving hosted IDs, independently inspect only the Sonnet, Codex, and direct-OpenAI Pi TestLiveSharedScenarios step results and uploaded common artifacts at exact candidate head.
  Captain canceled the evidence phase before race/hosted work because target-unaware skips would suppress rather than measure Codex/Pi journeys; no Cycle-4 run IDs exist.
- DONE: Treat Opus, Claude substrate proofs, Pi front-door smoke, and prior Sonnet break-glass evidence as irrelevant to this ys gate. Record their non-execution/skips accurately; do not broaden scope or repair adjacent product/runtime findings.
  None ran; Opus remains unverified/missing rather than an exception, and candidate bytes stayed untouched.
- FAILED: Recommend PASSED only if AC-1 through AC-4 have reproducible local plus exact-head common-selector evidence; append/push a path-scoped validation cycle-4 report with classifications, exact commands/SHAs/run IDs/artifacts, and final verdict.
  Recommend REJECTED. Material ys evidence defect: normal common lanes lose target-attributable evidence; `value-ac[AC-4]` requires exact runtime attribution; cheap all-adapter skips plus Sonnet-only run `30996911834` and 9a's Codex/Sonnet-only source prove the trigger.

### Summary

Cycle 4 confirms the structural four-owner accounting but rejects its runtime-global mechanism. The next correction must keep the desired registry a pure SSOT and represent proven failures, unverified targets, and passing evidence per target; no race, hosted run, or candidate mutation followed the Captain's stop.

## Stage Report: implementation (cycle 5)

- DONE: Audit the durable 26n and 9a evidence before changing selection; classify proven passes, proven owner failures, and unverified required targets without inference.
  Run `30996911834` detail JSONL proves Claude Sonnet default-headless failed after two preparations and Codex default-headless passed in 62.12s; Pi stopped earlier and Opus supplied no journey result. Run `30706782428` detail JSONL proves smallest-sufficient and keep-moving failed on both Codex (76.61s/128.60s) and Sonnet (299.34s/568.20s), and contains no auto-continue selector. The archived wm entity independently proves Claude auto-continue session `e59d1f3b` passed and records the Pi origin leg as environment-blocked, not a product failure.

  | Journey | Claude Sonnet | Claude Opus | Codex | Pi |
  |---|---|---|---|---|
  | `default-headless-gate-stop` | proven owner failure → `26nk8qd48zknqnn4kc123sez` | unverified, runnable | proven pass, runnable | unverified, runnable |
  | `auto-continue-after-implementation` | proven historical pass, runnable | unverified, runnable | unverified, runnable | environment-blocked/unverified, runnable |
  | `smallest-sufficient-mechanism` | proven owner failure → `9adv48yhye5s2vkhwd7ge52d` | unverified, runnable | proven owner failure → `9adv48yhye5s2vkhwd7ge52d` | unverified, runnable |
  | `keep-moving-posture` | proven owner failure → `9adv48yhye5s2vkhwd7ge52d` | unverified, runnable | proven owner failure → `9adv48yhye5s2vkhwd7ge52d` | unverified, runnable |

- DONE: Add red-first desired-registry purity and target-scope controls, then implement the smallest evidence seam behind adapter/model selection.
  RED: the target enumeration did not compile because adapters had no evidence-target identity, and registry reconciliation failed with `IMPURE desired registry contains observed missing evidence`. GREEN: adapters expose only Sonnet/Opus/Codex/Pi evidence identity; the common selector supplies it to the existing TODO seam; no runner map, fixture, assertion, host launch, command, or product behavior changed.
- DONE: Keep `docs/runtime-live-ci-registry.md` a pure desired-state SSOT and derive exact observed gaps from source.
  Removed the entire `## Missing live evidence` section. Reconciliation parses source TODO cases into five journey+target+owner results and emits exact `MISSING-EVIDENCE journey=<id> target=<target> owner=<owner>` diagnostics. Reintroduced-ledger, global binding, wrong-target, wrong-owner, removal/cardinality, duplicate-source, unowned, and proven-Codex-pass suppression mutations all fail.
- DONE: Enumerate Claude Sonnet, Claude Opus, Codex, and Pi cheaply; retain all 16 identities and leave every unverified target runnable.
  Four exact `TestSharedLiveTODOEvidenceSet/<target>` selectors pass without a model: Sonnet has three TODO cells, Codex has two, Opus has zero, and Pi has zero. Auto-continue is runnable on all four; default-headless stays runnable on Codex, Opus, and Pi.
- DONE: Reprove parity, first-failure, gate-start, Codex/Pi shim, workflow, registry, SHA, formatting, full, race, diff, and ceiling controls without a local-live, model-backed, or hosted run.
  Focused controls passed. `go test ./...` passed (`ensigncycle` 214.687s); `go test ./... -race` passed (`ensigncycle` 218.388s); `gofmt -w ./cmd ./internal`, `git diff --check`, exact current/stale SHA guards, and workflow selectors are green. No paid/local-live/hosted lane, Opus run, substrate proof, or Pi smoke ran.
- DONE: Preserve the watched-core then documentation-only binder shape, push both commits, and remain below the authorized surface ceiling.
  Watched core `386fd7a8bfa6f7c1597fe0171cee6cb5857a542a` is bound by docs-only head `e023619ed1f6a00a85f6ad88741c8b8ebad2d3b1`; both are pushed. Final surface against merge base `507a4bc12e48a3e4a813219602c488f09c81a5d8` is 40 files/+1870/-664, below 42 files/+2750 insertions.

### Summary

Cycle 5 replaces target-blind journey quarantines with five durable-evidence bindings while preserving the pure 16-journey/four-target desired registry. Proven passes and all unverified required targets remain runnable; only exact Sonnet/Codex owner failures skip. Source-derived reconciliation, adversarial purity/scope controls, full/race checks, and the two-commit SHA binding are green and pushed, with no product repair or model spend.

## Stage Report: validation (cycle 5)

- DONE: Confirm exact candidate identity: branch and origin head e023619ed1f6a00a85f6ad88741c8b8ebad2d3b1, watched core 386fd7a8bfa6f7c1597fe0171cee6cb5857a542a, docs-only binder, clean worktree, and current/stale SHA guard.
  Exact head/origin and merge base `507a4bc12e48a3e4a813219602c488f09c81a5d8` match; binder changes one docs line, current/stale guards pass, and the worktree is clean.
- DONE: Independently audit the full 16-journey by 4-target desired matrix; keep every registry journey and Claude Sonnet, Claude Opus, Codex, and Pi target required.
  Parity and target tests exercise 64 desired cells with one 16-entry map, three adapters, and four required evidence targets; no target is removed or excepted.
- DONE: Prove docs/runtime-live-ci-registry.md is desired-state-only and contains no observed missing-evidence ledger; verify reconciliation/report/source TODOs carry observed gaps instead.
  Registry has no ledger; detached reintroduction fails `IMPURE desired registry`, while reconciliation emits source-derived target/owner diagnostics.
- DONE: Reproduce the exact five target-scoped owner bindings: Sonnet default-headless-gate-stop to 26n; Sonnet and Codex smallest-sufficient-mechanism and keep-moving-posture to 9a; no other target/journey binding.
  `TestSharedLiveTODOEvidenceSet` and reconciliation report exactly Sonnet=3, Codex=2, Opus=0, Pi=0 with the exact owners.
- DONE: Verify auto-continue-after-implementation remains runnable for every target; Codex default-headless remains runnable as a proven pass; Pi and Opus affected cells remain runnable and are classified unverified/missing when not run.
  Enumeration finds no auto-continue binding and no Codex/Pi/Opus default-headless binding; unexecuted cells remain unverified/missing, never pass or exception.
- DONE: Run independent target-aware adversarial controls for global skip, wrong target, wrong owner, duplicate, cardinality/removal, registry impurity, and suppression of the proven Codex default-headless pass.
  All 28 mutations pass; detached Pi suppression fails enumeration plus `UNOWNED-EVIDENCE`, and detached registry ledger fails `IMPURE desired registry`.
- DONE: Verify all four zero-model enumeration selectors (claude-sonnet, claude-opus, codex, pi) and exact skip counts/reasons.
  Four `TestSharedLiveTODOEvidenceSet/<target>` selectors pass without model execution and retain exact 3/0/2/0 skip cardinality.
- DONE: Recheck the prior logger-shim propagation, ordered-suite first-failure stop, shallow/full-history reconciliation, and gate start-state overlap fixes.
  Focused shim removals, false-first sequence, depth-one checkout mutation/current history, and crossed gate starts all fail at their intended boundary.
- DONE: Audit target-aware Outcome, AC-2, AC-4, test plan, and Cycle 5 16-cell evidence report against exact durable run artifacts 30706782428 and 30996911834 plus archived wm auto-continue evidence.
  Durable evidence supports the five target cells, Codex default-headless pass, historical Claude auto-continue pass, and leaves Pi/Opus gaps runnable rather than inferred.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race independently; record exact results and surface/ceiling accounting.
  Formatting/diff are clean; normal passes (`ensigncycle` 195.947s), race passes (`ensigncycle` 206.792s), and 40 files/+1870/-664 remain below 42/+2750.
- DONE: Before any paid/model-backed/hosted action, send the FO a concise local checkpoint with candidate identity, all control outcomes, and any blocker; do not dispatch or approve CI yourself.
  Sent the all-green checkpoint before run `31015945064`; this validator launched, approved, rejected, canceled, and reran no CI job.
- FAILED: After the FO supplies exact hosted run IDs, independently inspect only those exact-head common-selector jobs, steps, logs, and artifacts. Classify each target/journey as PASS, owner-linked MISSING-EVIDENCE, unverified/missing, or genuine candidate failure; never classify Opus irrelevant/excepted.
  Codex run `31015945064`/job `92340224745`/step 16/artifact `8934638444` passed 2 then external-capacity INCONCLUSIVE; Pi run `31016570689`/job `92342373497`/step 13/artifact `8935708302` passed 5 then exposed both findings below; smoke skipped, Sonnet stopped before run, Opus required/unverified.
- DONE: Append a durable Cycle 5 validation report with explicit PASSED or REJECTED verdict and findings; do not mutate candidate code.
  REJECTED. Finding 1 — Material ys-owned evidence defect: the stdout-only oracle falsely says the Pi recorder was never invoked, while Pi metrics say `runtime/host: claude` and omit token/cost; normal Pi evidence stops at 6/16 and `value-ac[AC-4]` requires truthful artifact/metric attribution.
  Finding 2 — Material product-owned outcome defect, owner unresolved/Needs Captain decision: Pi records validation/1 early with `entries=2`, then appends the worker log and hand-writes Feedback Cycles with no second recorder call; `contract[skills/feedback-rejection-flow/SKILL.md#step-7]` requires the complete reviewer+worker log and Cycle line before the one retained round record. Candidate stayed untouched.

### Summary

Checklist: 12 DONE, 0 SKIPPED, 1 FAILED. Cycle 5 local representation proof is green, but exact-head Pi evidence exposes a Material ys-owned observation/metrics defect and a separate Material product-owned rejection-flow timing defect requiring Captain ownership, so validation is REJECTED; Codex capacity, unrun Sonnet, and required Opus remain inconclusive or unverified.

## Stage Report: implementation (cycle 6)

- DONE: Reproduce both Cycle-5 findings from exact run `31016570689`, job `92342373497`, artifact `8935708302` before editing.
  Pi root-session lines 108/109 correlate the exact rejection recorder call with a successful `entries=2` result; the prior metric files identify runtime/host as Claude and omit usage/cost.
- DONE: Add red-first strict Pi rejection-session controls without weakening the durable oracle.
  Missing, uncorrelated, errored, wrong-command, wrong-round, and `entries=2` transcripts fail; only the exact command plus correlated success with `entries=4` passes.
- DONE: Retain the archived Pi root session at the adapter/result boundary and use it for rejection evidence.
  Pi no longer reuses Claude stdout/stderr tool observations; final durable state still requires the complete retained four-entry round.
- DONE: Emit one native Pi shared-scenario metric for every one of the 16 common paths.
  The adapter-level collector covers full-cycle, zero-discovery, auto-continue, and AC-reanchor; embedded Claude emitters and the Claude-only shallow-window sibling are suppressed only for Pi.
- DONE: Prove Pi metric identity, measured usage, and artifact provenance from representative native session JSONL.
  Records carry `runtime=pi`, `host=pi`, `model=openai/gpt-5.4`, duration, input/output/cache/total tokens, total cost, GitHub run ID/URL, and a Pi filename; the prior Claude/no-usage shape fails.
- DONE: Harden the native parser against the exact artifact shape and partial attribution.
  A bounded 4 MiB scanner accepts the audited 105,396-byte row shape; missing cost, total tokens, cache write, or provider and mixed models fail, while explicit zero cache-write/cost remains valid.
- DONE: Add exactly one target-scoped product binding for Pi rejection-flow.
  Source reconciliation now reports six cells total and binds only Pi `rejection-flow` to `zbcj98qfwtax61vxdzrf615e`; Sonnet, Codex, and Opus suppression mutations fail.

  | Journey | Claude Sonnet | Claude Opus | Codex | Pi |
  |---|---|---|---|---|
  | `default-headless-gate-stop` | `26nk8qd48zknqnn4kc123sez` | runnable | runnable | runnable |
  | `smallest-sufficient-mechanism` | `9adv48yhye5s2vkhwd7ge52d` | runnable | `9adv48yhye5s2vkhwd7ge52d` | runnable |
  | `keep-moving-posture` | `9adv48yhye5s2vkhwd7ge52d` | runnable | `9adv48yhye5s2vkhwd7ge52d` | runnable |
  | `rejection-flow` | runnable | runnable | runnable | `zbcj98qfwtax61vxdzrf615e` |

- DONE: Run focused, formatting, reconciliation, full, race, current/stale, origin/head, diff, and ceiling checks on exact final bytes.
  Focused Pi/TODO/reconciliation tests exit 0; `go test ./...` exits 0 (`ensigncycle` 120.012s); race exits 0 (`ensigncycle` 208.740s); final surface is exactly 42 files/+2173/-666.
- DONE: Preserve the exact two-commit watched-core then docs-only binder shape and push both commits.
  Core `8cca2533723678f62cd0eff02c6442a95c49715e` is bound by docs-only head `35143a55a3fd1259295664df71602c116f6b3fd3`; origin/head match and the code worktree is clean.
- DONE: Preserve product bytes and avoid all unauthorized spend or scope expansion.
  No zbc/26n/9a repair, skill/binary behavior, member/controller/retry, local-live, hosted, model-backed, substrate, smoke, or Opus execution occurred.

### Summary

Cycle 6 replaces the false Claude observation boundary with strict Pi root-session evidence and native Pi metrics across all common journeys. The incomplete `entries=2` rejection remains truthfully blocked by the exact Pi-only zbc TODO, while the pure desired registry and all 64 required target/journey cells remain intact.

## Stage Report: validation (cycle 6)

- DONE: Independently verify exact candidate provenance and scope before trusting the implementation report.
  Candidate and origin are `35143a55a3fd1259295664df71602c116f6b3fd3`; watched core is `8cca2533723678f62cd0eff02c6442a95c49715e`; rejected baseline is `e023619ed1f6a00a85f6ad88741c8b8ebad2d3b1`; merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`. The binder changes only `docs/runtime-live-ci.md`; the code worktree is clean and head/origin match.
- DONE: Reproduce the exact Cycle-5 Pi artifact boundary without launching a model.
  The archived root session from run `31016570689`, job `92342373497`, artifact `8935708302` contains the exact recorder call and correlated successful result with `entries=2`. The new Pi observer reads that root-session JSONL and correctly rejects it as incomplete rather than repeating the prior stdout/stderr false negative.
- FAILED: Prove the Pi observer rejects every detached false-positive/false-negative variant and accepts only one exact correlated successful `entries=4` recording.
  Built-in controls reject missing, uncorrelated, errored, wrong-task, wrong-round, wrong-file, and `entries=2` evidence, and accept the exact single `entries=4` pair. A detached audit adds two required ambiguity controls: two exact successful calls, and one successful `entries=4` call followed by a second exact `entries=2` call. Both fail because `piRecordedRejectionRound` returns at the first qualifying result and never detects another matching invocation.
- DONE: Verify the pure desired registry, all 64 journey/target cells, and exactly six source-derived owner bindings.
  Reconciliation and cheap target enumeration preserve 16 journeys on four targets with Sonnet=3, Codex=2, Pi=1, and Opus=0 bindings. Pi `rejection-flow` maps only to `zbcj98qfwtax61vxdzrf615e`; no observed-gap ledger enters the desired registry.
- DONE: Exercise binding removal, target, owner, global, duplicate, cardinality, accidental suppression, ledger, and proven-pass mutations.
  `TestRuntimeLiveRegistryReconciliationMutationControls`, source-duplicate controls, all-target TODO enumeration, and the Codex proven-pass guard pass, including the three wrong-target Pi mutations.
- DONE: Verify native Pi metric identity, usage, cost, duration, and provenance controls.
  The emitter produces `runtime=pi`, `host=pi`, `model=openai/gpt-5.4`, measured input/output/cache/total tokens, total cost, duration, run ID/URL, and a Pi-attributed filename. Claude attribution, missing usage/cost/provider, mixed model, oversized-row, and explicit-zero controls behave as specified. The exact archived Pi session independently parses as OpenAI GPT-5.4 with nonzero tokens and cost.
- DONE: Recheck prior parity, adapter selection, exact selectors, first-failure, gate starts, launch/shim, artifact, SHA, and Claude/Codex metric boundaries proportionate to this evidence-only diff.
  Focused ensigncycle, release, contractlint, TODO-set, registry, current/stale SHA, and workflow tests pass; no product/runtime launch behavior changed in Cycle 6.
- DONE: Run required repository checks and final provenance checks on exact candidate bytes.
  `go test ./...` exits 0 (`internal/ensigncycle` 242.462s); `go test ./... -race` exits 0 (`internal/ensigncycle` 245.590s); focused live-tagged offline tests exit 0; `gofmt -l ./cmd ./internal` is empty; `git diff --check` passes. Final surface is exactly 42 files/+2173/-666, at the 42-file ceiling and below +2750 insertions.
- DONE: Classify the new finding before any candidate mutation or rerun and report the local checkpoint to the First Officer.
  Finding — evidence defect, Material/ys-owned narrow fix: the normal Pi model path can issue a repeated recorder command; observable harm is that ambiguous execution, including a later incomplete recorder result, is certified as exact successful evidence; `value-ac[AC-4]` requires attributable exact Pi evidence; the detached two-call controls reproduce the false positive. Candidate bytes remain unchanged.
- DONE: Avoid unauthorized live spend and adjacent product repair.
  No hosted, local-live, model-backed, smoke, substrate, Sonnet, Codex, Pi, or Opus run was launched. The zbc product outcome remains represented by the Pi-only TODO and was not repaired or reclassified.
- FAILED: Recommend PASSED only if every required observer mutation is green and no Material finding remains.
  Recommend REJECTED. AC-1 through AC-3 and the remaining AC-4 representation/metric boundaries are green, but the Pi recorder observer does not reject duplicate or ambiguous exact calls, leaving one Material AC-4 evidence defect.

### Summary

Cycle 6 correctly moves Pi evidence and metrics to the archived native root session, preserves the historical two-entry product failure, and adds the sixth exact target binding. Validation is nevertheless REJECTED because the observer accepts ambiguous duplicate recorder executions; a narrow ys-owned correction must require exactly one matching invocation/result pair before the same detached controls are rerun. No candidate mutation or live spend occurred.

## Stage Report: implementation (cycle 7)

- DONE: Reproduce the independent Cycle-6 false positives at exact candidate `35143a55a`.
  RED exit 1: two successful `entries=4` calls and `entries=4` followed by a second `entries=2` call both qualified because the extractor returned on the first success.
- DONE: Add red-first Pi extractor controls for every required ambiguous cardinality and ordering form.
  Two complete calls, both mixed orders, a resultless second invocation, reused ID, repeated results, result-before-invocation, and a duplicated success line all fail; removing the cardinality/order checks makes them pass.
- DONE: Implement the smallest fail-closed Pi-only correlation/cardinality correction.
  The extractor scans the whole session and requires one exact invocation, one globally unique call ID, one later correlated non-error Bash result, and one exact `entries=4` success-line occurrence.
- DONE: Preserve all prior negative boundaries and unrelated-tool behavior.
  Missing, uncorrelated, errored, wrong command/round/task, stdout lookalike, and `entries=2` remain false; unrelated tool calls/results around one exact pair remain irrelevant and the representative complete pair stays true.
- DONE: Reprove the exact archived root session remains rejected and the durable oracle remains strict.
  A temporary path-local control read the 734,304-byte audited root session and returned false for its correlated `entries=2`; committed durable-state tests still require the retained four-entry round.
- DONE: Reprove Cycle-6 native Pi metrics and all 16 route coverage unchanged.
  Live-tagged metric identity/usage/cost/provenance, 105 KiB row, partial-attribution, shared-runner coverage, runtime selection, scenario definitions, and promoted-entrypoint controls exit 0.
- DONE: Reprove the pure desired state and exact six target bindings.
  TODO enumeration and reconciliation/current-stale mutation suites exit 0; only Pi rejection-flow binds zbc, while removal, wrong target/owner, global/duplicate, and other-target suppression remain rejected.
- DONE: Run formatting, focused, full, race, diff, origin/head, and ceiling checks on final bytes.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 265.267s), race exits 0 (`ensigncycle` 172.545s), and origin equals head.
- DONE: Commit the watched-path fix first and bind its exact SHA in one documentation-only commit.
  Core `6e5fd1ac07f042fedc645d71c3a2af3afc0e7ab1` is bound by docs-only head `e952ed65ccafda13d54631195f52fdcd32fbfe46`; both are pushed and the surface remains 42 files/+2256/-666.
- DONE: Preserve Cycle-6 metrics/bindings/product bytes and avoid unauthorized evidence execution.
  No Pi metrics, registry, source binding, durable oracle, host launch, scenario order, product skill/binary, member/controller/retry, local-live, hosted, model-backed, substrate, smoke, or Opus behavior changed or ran.
- DONE: Keep the correction at the simplest evidence boundary serving AC-4.
  One existing Pi-only extractor changed; no new protocol, parser package, lifecycle state, process control, fixture, or product mechanism was introduced.

### Summary

Cycle 7 closes the sole Pi recorder-evidence ambiguity by requiring exactly one ordered invocation/result pair and one complete success occurrence. All Cycle-6 metrics and six target bindings remain intact, the audited two-entry product failure remains rejected, and focused/full/race evidence is green at the pushed two-commit candidate.

The candidate is ready for fresh independent validation without any evidence spend during implementation.

## Stage Report: validation (cycle 7)

- DONE: Independently inspect the Cycle-7 implementation report, exact candidate e952ed65c/core 6e5fd1ac0, merge base, changed paths, clean origin head, and docs-only binder; do not trust the implementer's summary.
  Candidate/origin are `e952ed65ccafda13d54631195f52fdcd32fbfe46`, core is `6e5fd1ac07f042fedc645d71c3a2af3afc0e7ab1`, merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`; core changes only the Pi observer/tests and the binder changes only `docs/runtime-live-ci.md`.
- DONE: Reproduce both Cycle-6 false positives on rejected head 35143a55a, then prove the exact Cycle-7 head rejects two entries=4 calls and entries=4 followed by entries=2.
  A detached old-head test observes both false positives; the same two transcripts return false on Cycle 7, as do the inverse mixed call order and a resultless second invocation.
- FAILED: Attack Pi correlation/cardinality with entries=2 then entries=4, second invocation without result, reused call ID, repeated correlated result messages, result-before-invocation, duplicate success lines within one text block, duplicate success text blocks, missing/error/wrong command/round/task, stdout lookalike, and unrelated tool traffic; only one later exact correlated entries=4 result may pass.
  All named call/result/order/reuse and duplicate-success controls behave correctly. Detached controls expose one remaining output ambiguity: one exact invocation plus one correlated result containing exact `entries=2` and `entries=4` summaries returns true, both when the summaries share one text block and when they occupy separate blocks.
- DONE: Reproduce the archived run 31016570689 root session and prove its single correlated entries=2 result remains rejected without weakening the complete four-entry durable-state oracle.
  The exact 734,304-byte root session from run `31016570689`/job `92342373497`/artifact `8935708302` returns false; the independent durable-state oracle still requires exactly four retained entries and exact canonical room bytes.
- DONE: Reprove Pi-native metrics for all 16 common paths, runtime/host/provider/model/tokens/cost/duration/run provenance, large rows, partial attribution rejection, and Claude/Codex metric isolation.
  Focused metric and 16-route tests pass; the archived session parses as `openai/gpt-5.4` with nonzero tokens/cost, while missing fields, mixed models, Claude attribution, and oversized-row controls remain discriminating.
- DONE: Reprove pure desired state with all 64 cells and exactly six bindings: Sonnet=3, Codex=2, Pi rejection-flow=zbc, Opus=0; run all removal/owner/target/global/duplicate/other-target/proven-pass mutations.
  TODO enumeration and reconciliation report the exact six source bindings; all registry, owner, target, duplicate, ledger, and proven-Codex-pass mutation controls exit 0.
- DONE: Reprove prior parity, first-failure, gate fixture, launcher shim, workflow selectors, artifact contract, current/stale SHA, formatting, diff, origin/head, and exact 42-file/+2750 ceiling.
  Focused ensigncycle/release/contractlint controls, current/stale SHA, `gofmt -l`, and `git diff --check` are green; exact surface is 42 files/+2256/-666, at the file ceiling and below +2750 insertions.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race independently on exact candidate bytes with explicit exit results; leave candidate untouched and clean.
  Read-only formatting verification is empty to preserve the no-mutation hold; normal exits 0 (`ensigncycle` 464.296s) and race exits 0 (`ensigncycle` 469.855s). Candidate and origin remain exact and clean.
- DONE: Classify every finding before any candidate mutation; do not authorize or implement zbc/26n/9a or broader containment behavior.
  Finding — evidence defect, Material/ys-owned narrow boundary: normal Pi root sessions represent result output as one or more text blocks; harm is that contradictory incomplete/complete recorder output is certified as exact evidence; `value-ac[AC-4]` requires attributable exact Pi evidence; detached mixed-output controls reproduce the false positive. No fix is authorized here.
- DONE: Before any live/model/hosted action, send the First Officer an exact local checkpoint and wait. Do not launch, approve, reject, cancel, or rerun CI yourself.
  Checkpoint sent before full/race with exact SHAs, green controls, and the possible finding; no local-live, model, hosted, substrate, smoke, Sonnet, Codex, Pi, or Opus action occurred.
- FAILED: Recommend PASSED only if all ambiguity mutants, strict historical rejection, Pi metrics, six bindings, full/race, history, and ceiling checks are green; append/push a path-scoped Cycle-7 validation report with explicit verdict.
  Recommend REJECTED. Historical rejection, metrics, bindings, repository checks, and all committed ambiguity controls are green, but contradictory `entries=2`/`entries=4` output inside one correlated result still passes the fail-closed Pi observer.

### Summary

Cycle 7 closes the rejected duplicate-invocation and repeated-success defects while preserving native metrics, six exact bindings, and the archived product TODO. Validation remains REJECTED because mixed incomplete/complete summaries within one result are accepted as exact evidence; candidate bytes stayed untouched and no live spend occurred.

## Stage Report: implementation (cycle 8)

- DONE: Reproduce Cycle-7 false positives at exact candidate `e952ed65c` for mixed canonical summaries.
  RED exit 1: `entries=2` then `entries=4` and the reverse qualified in one block; split-block `entries=2`/`entries=4` also qualified because only the success regex was counted.
- DONE: Add red-first canonical-summary content controls.
  Mixed `2/4` orders, `4+3`, `4+5`, split blocks, and repeated `2` fail; deleting the summary-count rule makes five mixed-complete cases pass, while one `4` plus unrelated diagnostic text remains true.
- DONE: Implement the smallest Pi-only canonical-summary cardinality correction.
  A new exact canonical-summary regex counts validation/1 summaries independent of entry count; the sole correlated result qualifies only when summary count and exact `entries=4` count are both one.
- DONE: Preserve every Cycle-7 invocation/result/order/reuse boundary.
  Two invocations, mixed call results, resultless duplicate, reused ID, repeated results, result-before-call, repeated success, wrong tool/error, and unrelated tool traffic controls all remain green.
- DONE: Reprove the exact archived run and representative complete session.
  A temporary path-local test read the 734,304-byte run `31016570689` root session and rejected its sole correlated `entries=2`; the committed singleton `entries=4` representative passes.
- DONE: Reprove Pi metrics, exact six bindings, registry purity, and prior parity unchanged.
  Native identity/usage/cost/provenance/large-row tests, all 16 metric routes, target TODO enumeration, registry mutations, runner coverage, runtime selection, and promoted entrypoints exit 0.
- DONE: Run focused rejection/Pi/TODO/reconciliation/current-stale controls in both directions.
  Three final-byte focused commands exit 0; incomplete/mixed summaries and all prior false positives fail, while singleton complete evidence and unrelated diagnostics pass.
- DONE: Run formatting, full, race, diff, origin/head, and ceiling checks with explicit results.
  `gofmt -w ./cmd ./internal` and diff are clean; full exits 0 (`ensigncycle` 255.055s), race exits 0 (`ensigncycle` 295.355s), origin equals head, and surface is 42 files/+2286/-666.
- DONE: Commit watched core first, bind its exact SHA, push, and leave the worktree clean.
  Core `ce6992d36cb3d84972d44049a4a45c9557c80bb2` is bound by docs-only head `2b0177dfe13da63dfd7cf30ceb6162a889e51890`; both are pushed and the binder changes only `docs/runtime-live-ci.md`.
- DONE: Keep the mechanism scoped to AC-4 result evidence and preserve product/runtime bytes.
  No durable oracle, Pi metric, target binding, registry, host launch, product skill/binary, fixture, scenario, lifecycle, controller, retry, or new file changed.
- DONE: Preserve the execution hold.
  No local-live, hosted, model-backed, substrate, smoke, Sonnet, Codex, Pi, or Opus lane was run or controlled during implementation.

### Summary

Cycle 8 makes Pi result evidence fail closed unless the sole correlated result contains exactly one canonical validation/1 summary and that summary is exactly `entries=4`. All Cycle-7 correlation boundaries, native metrics, six target bindings, product behavior, and the 42-file ceiling remain unchanged.

The pushed two-commit candidate is ready for fresh independent validation without evidence spend during implementation.
