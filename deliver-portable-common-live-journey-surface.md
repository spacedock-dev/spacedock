---
title: Deliver one portable common live-journey surface
status: validation
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
worktree: .worktrees/spacedock-ensign-deliver-portable-common-live-journey-surface-recarve
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
- Cycle 8: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `2b0177dfe13da63dfd7cf30ceb6162a889e51890`; surface 42 files/+2286/-666, strict Pi evidence, six bindings, focused controls, full, and race are green. Exact-head hosted Codex run `31029501075`, job `92386614637`, artifact `8940307070` passed `full-ensign-cycle`, `gate-guardrail`, and `default-headless-gate-stop`, then failed runnable `withdrawn-gate-recovery`: the intended prepare rejected uncommitted `evidence/command.log`, Codex bound a throwaway attempt 2, withdrew it, and opened attempt 3, leaving three attempts instead of the required original-withdrawn plus one open successor. Existing product owner `live-agent-withdraw-redo-extra-attempt` (`47gnqfm1ft6f2hcahz98m2jv`) records this exact signature. Captain disposition: FIX representation only. Add only Codex `withdrawn-gate-recovery` → `TODO(47gnqfm1ft6f2hcahz98m2jv)` as the seventh target-scoped binding; add red-first exact count/owner/target and wrong-target controls; preserve every product, fixture, prompt, adapter, launch, Pi metric/evidence, desired-registry, and prior six-binding byte. Commit watched source first, then one documentation-only SHA binder; run focused/full/race and return to the same validator before any further hosted spend. No new repair filing, product repair, retry, controller, member, Sonnet, Pi, or Opus execution is authorized in the correction.
- Cycle 9: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `88488379684801f02949544f54e2248b5be4488d`; surface 42 files/+2306/-666 and exact seven bindings are green. Exact-head isolated Codex run `31032033236`, job `92395174900`, artifact `8941373026` passed `full-ensign-cycle`, `gate-guardrail`, `default-headless-gate-stop`, skipped the owned withdrawn gap, and passed `recorded-gate-lifecycle`, then failed runnable `rejection-flow`: durable state contains one implementation and one REJECTED validation only; the FO prepared an ordinary human gate and stopped awaiting a decision instead of loading the feedback-rejection flow, recording the correction round, routing rework, and running cycle-2 validation. Existing cross-runtime correction-round owner `bind-post-rework-briefing-at-rejection-regate` (`zbcj98qfwtax61vxdzrf615e`) already owns Pi `rejection-flow` and records the same prior Codex one-implementation-report failure. Captain disposition: FIX representation only. Add only Codex `rejection-flow` → `TODO(zbcj98qfwtax61vxdzrf615e)` as the eighth target-scoped binding; add red-first exact count/owner/target and Pi/Sonnet/Opus wrong-target controls; preserve every product, fixture, prompt, adapter, launch, Pi metric/evidence, pure desired registry, and prior seven-binding byte. Update the operator inventory in the watched core, then one SHA-only documentation binder; run focused/full/race and return to the same validator before any further hosted spend. No new filing, zbc repair, retry, controller, member, Sonnet, Pi, or Opus execution is authorized in the correction.
- Cycle 10: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `f390bf5d3ac37882db6fb1af9edb34a9c991c9ff`; surface 42 files/+2325/-666, pure 64-cell desired state, exact eight bindings, focused controls, full, and race are green. Exact-head isolated Codex run `31034783581`, job `92404363889`, artifact `8942417379` passed `full-ensign-cycle`, `gate-guardrail`, `default-headless-gate-stop`, `recorded-gate-lifecycle`, `feedback-3-cycle-escalation`, `merge-hook-guardrail`, `filing`, `shallow-boot`, and `zero-discovery`, skipped the two exact owned Codex gaps, then reached `auto-continue-after-implementation`. That result is harness-invalid, not a product failure: the FO's sole required boot call correctly reported no commissioned workflow because `autoContinueReadme()` lacks `commissioned-by: spacedock@...`; the common runner also stages that worktree-backed lifecycle through the no-Git builder and invokes only the single-root fixture despite the contract requiring both single-root and split-root variants. Captain disposition: FIX only these ys-owned fixture/reachability defects before any further spend. Add red-first offline controls proving both exact fixture READMEs are discoverable commissioned workflows, both staged roots have the Git/state shape their lifecycle needs, the common journey launches exactly both variants through one top-level identity, and artifact labels distinguish the fixture IDs. Then add the commissioning markers, stage the real Git-backed fixtures, and run both variants without changing the host adapters, prompts, assertions, product skill/binary behavior, desired registry, or eight TODO bindings. Update the operator evidence note, commit watched core first, then one SHA-only documentation binder; run focused/full/race and return to the same validator before retrying Codex. No product TODO, owner filing, product repair, controller, retry mechanism, new file, new member, Sonnet, Pi, or Opus execution is authorized in the correction.
- Cycle 11: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `37d91f673bda7ccdf9b7e63f6dc0134407895173`; surface 42 files/+2558/-671, pure 64-cell desired state, exact eight bindings, fixture discovery, focused controls, full, and race are green. Exact-head hosted attempts `31038408585` and `31038934635` both failed in the secret-free offline gate before any environment approval or model spend. The first lost a reachable parent while `TestDurableQuestionedRejectsUnrelatedEdit` traversed a throwaway repository; the second failed `TestRecordedGateLifecycleEnteredStageRecoveryMatrix/committed_heading_only` because TempDir cleanup found `.git/objects` still being written. CI runs Git 2.54.0, whose geometric automatic maintenance defaults to enabled and detached at a 100-loose-object threshold; both Git-heavy fixtures cross that boundary, while the untouched first test passes 50/50 normal and 10/10 race-enabled on local Git 2.50.1. Captain disposition: FIX only the shared throwaway-repository lifecycle. Add a red-first `testgit.InitRepo` invariant, then disable automatic maintenance in every initialized test repository using repo-local configuration compatible with current and older Git. Preserve product/runtime/fixture/prompt/adapter behavior, desired registry, all eight TODO bindings, and the watched-core/SHA-only binder shape; run focused/full/race, return to the same validator, then retry the offline gate before approving any runtime. No product TODO, owner filing, model spend, broader containment repair, controller, retry mechanism, new member, Sonnet, Pi, or Opus execution is authorized in the correction.
- Cycle 12: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `99f27c41c88bdd414240344495013e0c393a4c41`; surface 42 files/+2668/-681, pure 64-cell desired state, exact eight bindings, focused controls, full, and race are green. Exact-head isolated Codex run `31041561707` was external-capacity INCONCLUSIVE at `feedback-3-cycle-escalation`; its serial retry `31042421641`, job `92430021762`, artifact `8945373289` passed nine runnable journeys and skipped the two exact owned Codex gaps before `auto-continue-after-implementation--auto-continue/single-root` failed. The durable transcript proves Codex advanced `implementation -> validation`, spawned a fresh validator, the validator committed `## Stage Report: validation` in `.worktrees/spacedock-ensign-auto-continue-task/auto-continue-task.md`, and the FO read that report and prepared the validation gate. The common runner nevertheless graded the stale pipeline-root entity copy, whose frontmatter advanced but whose worktree-owned report correctly remained absent. Captain disposition: FIX only this ys-owned end-state resolver false negative. Add a red-first single-root overlay control that reproduces pipeline/worktree divergence and fails unless the active worktree entity copy supplies the grade; retain split-root resolution from its state checkout and archive fallback; make the smallest common-runner/test-helper change. Preserve product skills/binary behavior, fixtures, prompts, host adapters, launch shapes, desired registry, all eight TODO bindings, Pi metrics/evidence, and the existing 42-file/+2750 ceiling. Commit watched core first, then one SHA-only documentation binder; run focused/full/race and return to the same independent validator before any hosted retry. No product TODO, owner filing, product repair, new file, controller, retry mechanism, new member, Sonnet, Pi, or Opus execution is authorized in the correction.
- Cycle 13: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `0aa74035cc96d672c71480466ed4e1b1a3169213`; surface 42 files/+2737/-681, pure 64-cell desired state, exact eight bindings, focused controls, full, and race are green. Exact-head isolated Codex run `31045591048`, job `92440554439`, artifact `8946459512` passed the offline gate, then failed the first runnable journey, `full-ensign-cycle`, after 166.28s. The durable artifact proves a product failure rather than capacity or harness error: each fresh Codex child received only the helper's `Read /tmp/...` pointer; the generated file claimed to contain shared ensign discipline but carried no executable `ensign-shared-core` bootstrap; the child wrote generic `## Stage Report` checkbox sections twice, and the FO accepted and archived them without any anchored `## Stage Report: {stage}`. The new exact repair owner is `restore-codex-ensign-contract-bootstrap` (`nvz2ym82ydfn07jp04yfxg9r`); it preserves fresh-context isolation and pointer-only assignment ownership while repairing the missing Codex contract bootstrap, and is a backlog owner rather than a new sprint member. Captain disposition: FIX ys representation only. Add only Codex `full-ensign-cycle` -> `TODO(nvz2ym82ydfn07jp04yfxg9r)` as the ninth target-scoped binding; add red-first exact count/owner/target and Sonnet/Pi/Opus wrong-target controls; preserve every product, fixture, prompt, adapter, launcher, Pi metric/evidence, pure desired-registry, and prior eight-binding byte. Update the operator inventory in the watched core, then one SHA-only documentation binder; remain within 42 files/+2750, run focused/full/race, and return to the same independent validator before any further hosted spend. No product repair, new file, controller, retry mechanism, new member, Codex rerun, Sonnet, Pi, or Opus execution is authorized in the correction; Opus remains required and unverified after the representation gate reopens execution.
- Cycle 14: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5`; surface 42 files/+2742/-681, pure 64-cell desired state, exact nine bindings, focused controls, full, and race are green. Exact-head isolated Pi run `31048506420`, job `92450151141`, artifact `8947855457` passed the offline gate and seven runnable journeys, skipped only owned Pi `rejection-flow`, then failed `filing`. The durable root session proves a ys-owned observer false negative rather than a product failure: Pi emitted one `bash` tool call containing `cat ... | ${SPACEDOCK_BIN:-spacedock} new wire-the-thing`, its correlated `role: toolResult` has `isError: false` and `created: .../wire-the-thing.md id=001`, and the correct backlog entity landed, but the shared filing runner graded Pi stdout/stderr with the Claude `tool_use` parser instead of grading `sessionJSONL`. Captain disposition: FIX only this Pi filing observation seam. Add a red-first distilled control for the exact Pi assistant `toolCall` plus correlated successful tool-result shape; route Pi filing through its archived root-session JSONL; accept only the exact successful atomic-create observation and fail closed on missing, failed, duplicate, reused-ID, or manual `--next-id` evidence. Preserve Claude/Codex filing behavior, every product/fixture/prompt/launcher/runtime front door, native Pi metrics, the pure registry, all nine target-scoped bindings, and the 42-file/+2750 ceiling. Commit watched core first, then one SHA-only documentation binder; run focused/full/race and return to the same independent validator before resuming isolated Pi evidence. Do not file a product owner or TODO, repair product behavior, add a file/member/controller/retry, or run Codex/Sonnet/Opus/hosted/model/substrate/smoke execution during the correction; Opus remains required and unverified.
- Cycle 15: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `6408e2ad658f05d9c990b88a46cc8d002203c729`; surface 42 files/+2750/-684, pure 64-cell desired state, exact nine bindings, Pi filing controls, focused checks, full, and race are green. Exact-head isolated Pi run `31052117603`, job `92461666431`, artifact `8949169013` passed offline and eight runnable journeys including `filing`, skipped only owned Pi `rejection-flow`, then failed `shallow-boot` after 38.51s with `stream carried no assistant turns — nothing to check`. The durable/native evidence proves another ys-owned observer false negative: the host-neutral state assertion already passed; Pi stdout contains the exact required held-review greet; the root session ends in a native assistant text greet after only `read`/`bash` calls, with no `subagent`, `member_spawn`, or `delegate`; but the runner passes Pi stdout/stderr to both Claude-only turn checks (`assertNoTeamCreateBeforeGreet` and then `assertShallowBootMeasured`). Captain disposition: FIX only the Pi shallow-boot observation seam. Add a red-first distilled native-session positive plus causal missing-greet and pre-greet Pi dispatch/team-call negatives; route Pi through one native root-session structural/no-dispatch assertion while leaving both Claude stream checks and Claude boot-window metrics byte/behavior intact. Preserve the durable shallow-boot assertion, Codex path, native Pi whole-journey metrics, product/fixture/prompt/launcher/front door, pure registry, all nine bindings, and exactly 42 files/+2750; add no file. Commit watched core then one SHA-only binder, run focused/full/race, and return to the same validator before another isolated Pi run. Do not add an owner/TODO, repair product behavior, add a member/controller/retry, or run Codex/Sonnet/Opus/hosted/model/substrate/smoke during correction; Opus remains required and unverified.
- Cycle 16: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `37c7e2660caa6148a590b21671c2d5f3693b0c74`; surface 42 files/+2748/-708, pure 64-cell desired state, exact nine bindings, complete Pi shallow-boot controls, focused checks, full, and race are green. Exact-head isolated Pi run `31054684260`, job `92469677841`, artifact `8949983771` passed offline, `full-ensign-cycle`, and `gate-guardrail`, then failed runnable `default-headless-gate-stop` after 228.61s with `gate hold crossed its committed no-authority boundary`. The durable native sessions prove a product failure rather than capacity or harness error: the root advanced to validation and dispatched a validation worker; that worker prepared the correct single open validation gate but wrote a validation report and committed the state with raw `git commit` instead of the required `spacedock state commit`; the final gate remained open and unconsumed. This is the same direct-gate outcome owned for Sonnet by `headless-recorded-gate-stop-stage-coherence` (`26nk8qd48zknqnn4kc123sez`), whose current evidence records the forbidden validation-worker/report route. Captain disposition: FIX ys representation only. Add only Pi `default-headless-gate-stop` -> `TODO(26nk8qd48zknqnn4kc123sez)` as the tenth target-scoped binding; add red-first exact count/owner/target controls and preserve Sonnet/Codex/Opus discrimination. Preserve every product, fixture, prompt, adapter, launcher, runtime front door, Pi observer/metric path, pure desired registry, and prior nine bindings. Update the operator inventory in the watched core, then one SHA-only documentation binder; remain within 42 files/+2750 with no new file, run focused/full/race, and return to the same independent validator before another isolated Pi run. Do not repair 26n, add a member/controller/retry, or run Codex/Sonnet/Opus/hosted/model/substrate/smoke during correction; Opus remains required and unverified.
- Cycle 17: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `23e0dbd02b3c4aa8753776596eadf9980698a0b5`; surface 42 files/+2749/-708, pure 64-cell desired state, exact ten bindings, focused checks, full, and race are green. Exact-head isolated Pi run `31056949526`, job `92476577900`, artifact `8951115959` passed offline and eleven runnable journeys, skipped only the two exact owned Pi gaps, then failed runnable `smallest-sufficient-mechanism` after 186.96s. The durable native session proves a product failure rather than capacity or harness error: for both `ready-one` and `ready-two`, Pi changed `ready -> done`, committed `dispatch: <slug> entering done`, and built `--stage done`; the resulting ready-stage reports do not supply the required path-scoped `ready` entry carrying `started`. This exactly matches existing owner `repair-entered-stage-dispatch-and-post-gate-terminalization` (`9adv48yhye5s2vkhwd7ge52d`), which records the same initial-stage terminal dispatch defect for Sonnet and Codex. Captain disposition: FIX ys representation only. Add only Pi `smallest-sufficient-mechanism` -> `TODO(9adv48yhye5s2vkhwd7ge52d)` as the eleventh target-scoped binding; add red-first exact count/owner/target controls and preserve Sonnet/Codex/Opus discrimination. Preserve every product, fixture, prompt, adapter, launcher, runtime front door, Pi observer/metric path, pure desired registry, and prior ten bindings. Update the operator inventory in the watched core, then one SHA-only documentation binder; remain within 42 files/+2750 with no new file, run focused/full/race, and return to the same independent validator before resuming isolated Pi evidence. Do not repair 9a, add a member/controller/retry, or run Codex/Sonnet/Opus/hosted/model/substrate/smoke during correction; Opus remains required and unverified.
- Cycle 18: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `00fcde8a5121fe937a75d152cb7b498f1d16a10a`; surface 42 files/+2748/-708, pure 64-cell desired state, exact eleven bindings, focused checks, full, and race are green. Exact-head isolated Pi run `31060796822`, job `92488310460`, artifact `8952141495` passed offline and `full-ensign-cycle`, then failed runnable `gate-guardrail` after 101.98s with `prepared fixture request count = 0, want 1`. The durable native session proves a product failure rather than capacity or harness error: Pi selected the committed review, snapshot, and recorder contract plus untracked `evidence/command.log`; `gate prepare` correctly rejected that source with the generic exact-committed-file error, and Pi halted with no bound room or state mutation. This is the exact pathless rejected-source seam documented by existing owner `sonnet-gate-guardrail-no-authority` (`3zzpdw704df1g8pg1x9thzmw`), whose repair attributes the rejected source so the FO can correct its selection. Captain disposition: FIX ys representation only. Add only Pi `gate-guardrail` -> `TODO(3zzpdw704df1g8pg1x9thzmw)` as the twelfth target-scoped binding; add red-first exact count/owner/target controls and preserve Sonnet/Codex/Opus discrimination. Preserve every product, fixture, prompt, adapter, launcher, runtime front door, Pi observer/metric path, pure desired registry, and prior eleven bindings. Update the operator inventory in the watched core, then one SHA-only documentation binder; remain within 42 files/+2750 with no new file, run focused/full/race, and return to the same independent validator before resuming isolated Pi evidence. Do not repair 3z, add a member/controller/retry, or run Codex/Sonnet/Opus/hosted/model/substrate/smoke during correction; Opus remains required and unverified.
- Cycle 19: REJECTED by the Captain after fresh independent Sol/medium validation PASSED at `2240fabe5d975a1ba724848ac2fc91ef11eb2cc7`; surface 42 files/+2749/-708, pure 64-cell desired state, exact twelve bindings, focused checks, full, and race are green. Exact-head isolated Pi run `31062662883`, job `92493828626`, artifact `8953098453` passed ten runnable journeys, skipped the four exact owned Pi gaps, then failed runnable `keep-moving-posture` after 377.86s with `questioned must remain active and nonterminal`. The durable native session proves a product failure rather than capacity or harness error: Pi set `questioned` to ideation with `verdict=revise`, later advanced only `status=review`, and never cleared the verdict, so the final prose called it active while durable state remained terminal. This is the same keep-moving/post-gate terminalization seam already owned for Sonnet and Codex by `repair-entered-stage-dispatch-and-post-gate-terminalization` (`9adv48yhye5s2vkhwd7ge52d`). Captain disposition: FIX ys representation only. Add only Pi `keep-moving-posture` -> `TODO(9adv48yhye5s2vkhwd7ge52d)` as the thirteenth target-scoped binding; add red-first exact count/owner/target controls and preserve Sonnet/Codex/Opus discrimination. Preserve every product, fixture, prompt, adapter, launcher, runtime front door, Pi observer/metric path, pure desired registry, and prior twelve bindings. Update the operator inventory in the watched core, then one SHA-only documentation binder; remain within 42 files/+2750 with no new file, run focused/full/race, and return to the same independent validator before resuming isolated Pi evidence. Do not repair 9a, add a member/controller/retry, or run Codex/Sonnet/Opus/hosted/model/substrate/smoke during correction; Opus remains required and unverified.

- Cycle 23: REJECTED — fresh independent surface-recarve validation at `b1f49a395ccccf296b294617ed0f8b762c7c5e9d`; surface exactly 34 files/+1083/-1947 against the Captain hard bound of 34 files/<=+1083/>=-1883; AC unchanged. FO disposition: FIX exactly two Material, ys-owned representation defects. Route to implementation: make the exact bound builder/assertion values typed exercise parameters, call both directly, and reject dead references both statically and at runtime; move all four concrete Claude branches and the copied-plugin grade behind explicit driver capabilities; remove only stale scenario metadata that creates honest line headroom. Preserve the sixteen declarations and direct seven-argument calls, exact 34-file set and line bounds, every fixture/assertion and Claude evidence behavior, product bytes, Sonnet=4/Opus=0/Codex=5/Pi=5 bindings, three exact selectors, and the untouched rejected branch. Run focused dead-reference and wrapped-driver controls, live-tag compilation, formatting, isolated full/race, then return to the same validator. No live/model lane, PR, product repair, new file, map, switch, controller, mutation lab, or SHA guard is authorized before validation passes.

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

## Stage Report: validation (cycle 8)

- DONE: Independently inspect exact candidate 2b0177dfe/core ce6992d36, implementation report, merge base, clean origin, docs-only binder, and exact 42-file/+2750 ceiling.
  Candidate/origin are `2b0177dfe13da63dfd7cf30ceb6162a889e51890`, core is `ce6992d36cb3d84972d44049a4a45c9557c80bb2`, merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`; binder is docs-only and exact surface is 42 files/+2286/-666.
- DONE: Reproduce Cycle-7 mixed entries=2/entries=4 false positives on e952ed65c, then prove 2b0177dfe rejects same-block and split-block forms in both orders.
  Detached rejected-head tests observe all four false positives; the identical transcripts return false at Cycle 8 because canonical summary cardinality is two.
- DONE: Attack one-result content with entries=4 plus entries=2/3/5, repeated entries=2, repeated entries=4, duplicate text blocks, unrelated diagnostic text containing entry-like prose, malformed/noncanonical summaries, and one exact canonical entries=4 line; require exactly one canonical summary and exact complete success.
  Detached controls reject canonical counts 0/2/3/5/999999, mixed/repeated summaries, and split canonical text; one exact `entries=4` with unrelated or noncanonical diagnostics passes.
- DONE: Reprove every Cycle-7 invocation/result/order/reuse/tool-name/error/duplicate-call/duplicate-result/stdout/unrelated-tool adversary and the exact archived run 31016570689 entries=2 rejection.
  All prior correlation mutants fail closed, unrelated traffic remains irrelevant, and the exact 734,304-byte run `31016570689` root session remains false for its sole `entries=2` result.
- DONE: Reprove all Pi-native metric routes and runtime/host/provider/model/token/cost/duration/provenance/large-row/partial-attribution controls, with Claude/Codex isolation.
  Focused metric and 16-route controls pass; the archived session independently parses as `openai/gpt-5.4` with nonzero tokens/cost, while partial/mixed/Claude-attributed forms fail.
- DONE: Reprove pure 64-cell desired state and exactly six source bindings (Sonnet3/Codex2/Pi1/Opus0) plus every owner/target/global/duplicate/suppression/ledger mutation.
  TODO enumeration and reconciliation emit only the exact six cells, including Pi rejection-flow→zbc; all removal, retag, move, duplicate, ledger, and proven-pass suppression controls are discriminating.
- DONE: Reprove parity, first-failure, gate fixtures, launcher shim, workflow selectors, artifacts, current/stale SHA, formatting, diff, origin, history shape, and ceiling.
  Focused ensigncycle/release/contractlint, SHA, parity, fixture, selector, and shim tests pass; `gofmt -l`, `git diff --check`, two-commit history, exact origin, and 42/+2286 ceiling checks are clean.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race independently with explicit exit results; never mutate candidate.
  Read-only formatting verification is empty under the no-mutation hold; normal exits 0 (`ensigncycle` 221.329s), race exits 0 (`ensigncycle` 225.588s), and candidate bytes remain clean.
- DONE: Classify any finding before mutation and keep all product/containment repairs outside ys.
  No finding remains. The zbc/26n/9a product outcomes stay represented only by their six exact source bindings; no candidate or containment repair was authorized or performed.
- DONE: Send the First Officer a complete local checkpoint before any live/model/hosted action; do not launch or control CI.
  The checkpoint reported exact SHAs, full detached matrix, archived evidence, metrics, bindings, and no finding before full/race; no live, hosted, model, smoke, substrate, or runtime lane was launched or controlled.
- DONE: Recommend PASSED only when all result-content and prior ambiguity matrices, archived rejection, metrics, bindings, full/race/history/ceiling proof are green; append/push explicit Cycle-8 validation verdict.
  Recommend PASSED. AC-1 through AC-4 have reproducible local evidence, the complete adversarial matrix is green, the archived product failure remains truthful, and no Material or deferred finding remains.

### Summary

Cycle 8 closes the final Pi result-content ambiguity while preserving native metrics, the pure 64-cell desired table, and all six exact product bindings. Independent focused, detached, full, race, history, and ceiling proof is green, so validation recommends PASSED without any live spend or candidate mutation.

## Stage Report: implementation (cycle 9)

- DONE: Start from the clean approved Cycle-8 candidate and preserve its product behavior.
  Baseline candidate is `2b0177dfe13da63dfd7cf30ceb6162a889e51890`; its watched core is `ce6992d36cb3d84972d44049a4a45c9557c80bb2`.
- DONE: Represent the exact hosted Codex failure as the seventh source binding.
  `Codex + withdrawn-gate-recovery` now maps to existing repair `47gnqfm1ft6f2hcahz98m2jv`.
- DONE: Ground the binding in durable hosted evidence already supplied by the First Officer.
  Codex run `31029501075`, job `92386614637`, artifact `8940307070` passed full-ensign, gate-guardrail, and default-headless before withdrawn-gate recovery exhausted three attempts.
- DONE: Add the expected binding red-first to the live TODO evidence set.
  At baseline, `go test -tags live ./internal/ensigncycle -run '^TestSharedLiveTODOEvidenceSet$' -count=1` exited 1 because the Codex withdrawn-gate cell had no reason.
- DONE: Add reconciliation coverage and owner/target mutation controls red-first.
  At baseline, `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1` exited 1 because the seventh exact owner binding was absent.
- DONE: Make the smallest representation-only implementation.
  Shared scenario desired state adds one exact Codex withdrawn-gate case and reuses the existing `47gnqfm1ft6f2hcahz98m2jv` owner; no repair behavior was added.
- DONE: Make owner mutations discriminating.
  Reconciliation tests move only that cell to Pi, Claude Sonnet, and Claude Opus and require each wrong-owner form to fail.
- DONE: Keep the operator inventory truthful.
  `docs/runtime-live-ci.md` now says seven target-scoped `MISSING-EVIDENCE` results and names the Codex run, journey, repair, prior passes, and three-attempt failure.
- DONE: Preserve pure desired-state semantics.
  The matrix remains 64 desired cells; TODO bindings remain non-passing evidence, and unverified Pi/Opus targets remain required and runnable.
- DONE: Preserve all prior six source bindings and evidence boundaries.
  Sonnet remains 3, prior Codex remains 2, Pi rejection-flow remains 1, and the new Codex binding raises the exact total to 7.
- DONE: Re-run focused live, rejection, Pi metric, large-row, reconciliation, mutation, and SHA controls on final bytes.
  The final focused commands exited 0; ensigncycle completed in 0.478s and contractlint in 0.407s.
- DONE: Apply required formatting and repository checks.
  `gofmt -w ./cmd ./internal` made no final diff, and `git diff --check` is clean.
- DONE: Run the full repository suite on the final candidate.
  `go test ./...` exited 0; the slowest reported package was `internal/status` at 18.431s.
- DONE: Run the full race suite on the same final bytes.
  `go test ./... -race` exited 0; `internal/ensigncycle` completed in 119.825s.
- DONE: Preserve the exact two-commit watched-core/SHA-binder history.
  Core `af214fc4f3452241205f170ec21676bae49f66c8` is followed by docs-only binder `88488379684801f02949544f54e2248b5be4488d`.
- DONE: Verify the binder and remote state.
  The registry SHA names the exact core; the binder changes only `docs/runtime-live-ci.md`; candidate and origin both equal `88488379684801f02949544f54e2248b5be4488d` and the worktree is clean.
- DONE: Stay within the fixed review surface.
  The merge-base diff is exactly 42 files, 2,306 insertions, and 666 deletions, below the +2,750 insertion ceiling.
- DONE: Preserve the execution hold and ownership boundaries.
  No product fix, fixture, prompt, adapter, launcher, metric, live/model/hosted/substrate run, retry, approval, rejection, or external evidence spend occurred in implementation.

### Summary

Cycle 9 truthfully adds the supplied Codex withdrawn-gate failure as the seventh exact TODO source binding. The pushed representation-only pair preserves all product/runtime behavior and prior evidence, while focused, full, race, history, origin, and 42-file ceiling checks are green.

## Stage Report: validation (cycle 9)

- DONE: Independently audit exact candidate, core, merge base, history, origin, cleanliness, and surface before trusting the implementation report.
  Candidate/origin are `88488379684801f02949544f54e2248b5be4488d`, core is `af214fc4f3452241205f170ec21676bae49f66c8`, merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`; the binder is SHA-only and surface is 42 files/+2306/-666.
- DONE: Independently verify hosted run `31029501075`, job `92386614637`, artifact `8940307070`, and owner fit.
  Downloaded baseline-head artifacts prove passes for full-ensign-cycle (100.81s), gate-guardrail (35.26s), and default-headless (167.14s), then withdrawn recovery fails at 86.35s with withdrawn attempts 1/2 and open attempt 3; existing owner `47gnqfm1ft6f2hcahz98m2jv` describes that exact CAS-rejection/withdraw-redo signature.
- DONE: Replay red/green binding logic at baseline and candidate.
  A detached `2b0177dfe` control exits 1 because Codex withdrawn recovery has no TODO; the same exact-owner assertion and source extraction pass at current head.
- DONE: Prove the desired registry stays pure at 64 cells and source evidence contains exactly seven bindings.
  Enumeration/reconciliation report Sonnet=3, Codex=3, Pi=1, Opus=0, with only Codex withdrawn-gate-recovery→47gn added and no desired-state evidence ledger.
- DONE: Attack the new cell and all prior evidence mappings in both directions.
  Detached new-cell removal, exact wrong owner, Pi/Sonnet/Opus moves, and global scope fail; built-in duplicate/cardinality, ledger, unowned, prior-owner, other-target, and Codex proven-pass suppression mutations all remain discriminating.
- DONE: Verify operator documentation and watched-core/binder partitioning.
  Operator docs state seven bindings and accurately name run `31029501075`, the three prior passes, three-attempt outcome, Codex target, journey, and 47gn owner; core contains every narrative/source/test change and the second commit changes only the bound SHA.
- DONE: Prove no product/runtime surface or new owner was introduced.
  The Cycle-9 core touches only four existing docs/reconciliation/TODO-test files; no product, fixture, prompt, adapter, launcher, Pi metric/evidence, runtime behavior, new file, or owner entity changes.
- DONE: Recheck prior Pi strict evidence, native metrics, six bindings, parity, fixtures, shim, workflow, artifact, and SHA controls proportionately.
  Focused live-tagged offline, ensigncycle, release, reconciliation/mutation, and current/stale SHA commands exit 0; Pi incomplete/mixed evidence and partial metrics still fail while exact native forms pass.
- DONE: Run formatting, diff, full, race, origin, history, cleanliness, and ceiling checks on exact bytes.
  `gofmt -l ./cmd ./internal` is empty and `git diff --check` passes; `go test ./...` exits 0 (`ensigncycle` 132.741s), race exits 0 (`ensigncycle` 136.020s), candidate/origin match, and 42/+2306 stays below +2750.
- DONE: Classify findings before mutation and preserve execution/ownership holds.
  No Material, deferred, or polish finding remains; no candidate/product repair or local-live, model, hosted, substrate, smoke, Sonnet, Codex, Pi, or Opus action was launched or controlled.
- DONE: Record explicit verdict without transitioning state.
  Recommend PASSED. AC-1 through AC-4 remain reproducible, the seventh binding is exact and artifact-grounded, all mutations and repository checks are green, and the entity remains in validation for First Officer handling.

### Summary

Cycle 9 independently confirms that the prior hosted Codex failure is represented by exactly one new target-scoped 47gn TODO, raising the truthful source set to seven without changing product behavior. Artifact, detached mutation, focused, full, race, history, and ceiling proof is green, so validation recommends PASSED with no live spend or candidate mutation.

## Stage Report: implementation (cycle 10)

- DONE: Start from the clean pushed Cycle-9 candidate and preserve its watched core.
  Baseline candidate is `88488379684801f02949544f54e2248b5be4488d`; watched core is `af214fc4f3452241205f170ec21676bae49f66c8`.
- DONE: Represent the supplied Codex rejection-flow failure as binding eight.
  `Codex + rejection-flow` now maps to existing cross-runtime repair `zbcj98qfwtax61vxdzrf615e`; no owner was filed.
- DONE: Ground the binding in the exact durable hosted evidence.
  Run `31032033236`, job `92395174900`, artifact `8941373026` passed full-ensign, gate-guardrail, default-headless, and recorded-gate-lifecycle, skipped owned withdrawn, then failed rejection-flow.
- DONE: Preserve the durable failure distinction.
  The artifact has one implementation report and one REJECTED validation; the transcript prepared an ordinary gate and stopped instead of routing rework and running cycle-2 validation.
- DONE: Add the expected Codex binding red-first.
  Baseline `TestSharedLiveTODOEvidenceSet` exited 1 because Codex rejection-flow returned no reason instead of exact `TODO(zbcj98qfwtax61vxdzrf615e)`.
- DONE: Add reconciliation and mutation coverage red-first.
  Baseline `TestRuntimeLiveRegistryReconciliation` exited 1 with exact missing journey=rejection-flow, target=codex, owner=zbc diagnostic.
- DONE: Implement the smallest source representation.
  One Codex rejection-flow switch case reuses the existing rejection owner and leaves the desired journey selected.
- DONE: Preserve the Pi rejection-flow binding distinctly.
  The exact audited map contains separate Pi and Codex rejection-flow keys with the same zbc owner; neither replaces the other.
- DONE: Prove exact count, owner, target, and target-move behavior.
  TODO/reconciliation controls require eight bindings and move the new Codex cell to Pi, Sonnet, or Opus; each owner/target/cardinality mutation fails.
- DONE: Reprove prior source-boundary mutations.
  Removal, retag, wrong target, global scope, duplicate, desired-ledger, unowned, and proven-pass suppression controls remain discriminating.
- DONE: Preserve the pure 64-cell desired state and exact distribution.
  Final bindings are Sonnet=3, Codex=4, Pi=1, Opus=0; all prior seven mappings remain exact and unverified cells remain runnable.
- DONE: Update the operator inventory in watched core.
  Docs now state eight bindings and name run `31032033236`, prior passes/skip, ordinary-gate stop, missing rejection routing, and cycle-2 validation outcome.
- DONE: Reprove prior Pi rejection evidence and native metric boundaries.
  Focused round-recording/invocation and Pi identity/cost/large-row controls exit 0; incomplete or partially attributed forms remain rejected.
- DONE: Run required formatting and focused final-byte checks.
  `gofmt -w ./cmd ./internal`, `git diff --check`, live TODO, reconciliation/mutations/SHA, Pi rejection, and Pi metric commands exit 0.
- DONE: Run both repository-wide suites on exact candidate bytes.
  `go test ./...` exits 0 (`ensigncycle` 104.890s); `go test ./... -race` exits 0 (`ensigncycle` 116.595s).
- DONE: Preserve the exact watched-core/SHA-binder history and remote equality.
  Core `79d4055e8fae70c4e2f74c9edabc3f5f1ab45544` is followed by SHA-only docs binder `f390bf5d3ac37882db6fb1af9edb34a9c991c9ff`; origin equals clean head.
- DONE: Stay inside the fixed review ceiling without product/runtime expansion.
  Surface is exactly 42 files/+2325/-666; Cycle 10 core changes only four existing representation/test/docs files and adds no file, member, controller, retry, or product repair.
- DONE: Preserve the execution and state-transition holds.
  No local-live/model/hosted/substrate/smoke/Sonnet/Codex/Pi/Opus action or evidence spend occurred, and implementation did not transition entity state.

### Summary

Cycle 10 adds the supplied Codex rejection-flow failure as the eighth exact target-scoped TODO while preserving the distinct Pi binding to the same owner. The pushed representation-only pair is green under focused, mutation, SHA, full, race, history, origin, and ceiling proof with no live spend or product/runtime change.

## Stage Report: validation (cycle 10)

- DONE: Verify exact candidate provenance and immutable remote equality.
  Candidate/origin are clean at `f390bf5d3ac37882db6fb1af9edb34a9c991c9ff`; watched core is `79d4055e8fae70c4e2f74c9edabc3f5f1ab45544` and merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`.
- DONE: Independently verify the supplied hosted Codex evidence without launching a run.
  Run `31032033236`, job `92395174900`, artifact `8941373026`, and artifact `source-head.txt` bind baseline `88488379684801f02949544f54e2248b5be4488d`.
- DONE: Verify the exact hosted journey outcome and durable failure state.
  Full-ensign, gate-guardrail, default-headless, and recorded-gate-lifecycle passed; owned withdrawn recovery skipped; rejection-flow retained one implementation and one REJECTED validation, prepared an ordinary gate, and never ran rework or cycle-2 validation.
- DONE: Verify the existing repair owner fits the observed Codex failure.
  Archived owner `zbcj98qfwtax61vxdzrf615e` explicitly owns rejection-flow rework/revalidation and the post-rework gate boundary; no owner was created or filed.
- DONE: Reproduce the Cycle-10 red baseline and green candidate independently.
  A detached baseline test exits 1 because Codex rejection-flow is empty at `884883796`; the identical exact-owner assertion passes on `f390bf5d3` while the separate Pi binding remains present.
- DONE: Prove pure desired state and the exact observed evidence representation.
  Independent parsing proves 16 journeys × four targets = 64 pure cells and exactly eight bindings: Sonnet=3, Codex=4, Pi=1, Opus=0; only the new Codex rejection-flow cell maps to `zbcj98qfwtax61vxdzrf615e`.
- DONE: Perform the semantic adversarial and mutation pass.
  New-cell removal, wrong owner, global scope, extra cardinality, and moves to Pi/Sonnet/Opus all fail; repository controls also reject duplicate source, desired ledger, unowned/wrong target, and proven-pass suppression mutations.
- DONE: Reproduce focused AC, reconciliation, SHA, release, and prior Pi controls.
  Parity/runtime selection, promoted entrypoints, gate/auto-continue/zero-discovery/re-anchor, workflow/cost, exact four-entry Pi rejection recording, native metric identity/completeness/large-row, observer, and negative single-cycle controls all exit 0.
- DONE: Verify documentation, history, surface, and semantic boundaries.
  Docs state eight and accurately cite run `31032033236`; Cycle 10 changes only four existing representation/test/docs files, the binder changes only the SHA, and total surface is 42 files/+2325/-666.
- DONE: Verify formatting and exact-byte repository suites.
  Read-only `gofmt -l ./cmd ./internal` is empty and `git diff --check` passes; `go test ./...` exits 0 (`ensigncycle` 103.286s, wall 115.22s) and `go test ./... -race` exits 0 (`ensigncycle` 116.375s, wall 129.29s).
- DONE: Preserve candidate, runtime, and workflow authority boundaries.
  Validation changed no candidate/product/runtime/fixture/prompt/adapter/launcher/Pi metric/observer file, launched no local-live/model/hosted/substrate/smoke action, filed no owner, and did not transition entity state.
- DONE: Recommend an explicit validation verdict.
  PASSED — all latest Cycle-10 truth requirements have independent evidence, with no material findings, deferred risks, or polish findings.

### Summary

Cycle 10 PASSED independently at exact candidate `f390bf5d3ac37882db6fb1af9edb34a9c991c9ff`. The supplied hosted failure validly adds only Codex rejection-flow to the existing `zbcj98qfwtax61vxdzrf615e` owner while preserving Pi's distinct binding, pure 64-cell desired state, and every no-product/no-live boundary.

## Stage Report: implementation (cycle 11)

- DONE: Start from the exact clean Cycle-10 candidate and preserve its evidence registry.
  Baseline candidate is `f390bf5d3ac37882db6fb1af9edb34a9c991c9ff`; watched core is `79d4055e8fae70c4e2f74c9edabc3f5f1ab45544`.
- DONE: Classify the supplied Codex evidence before changing any product ownership.
  Run `31034783581`, job `92404363889`, artifact `8942417379`, source `f390bf5d3` is harness-invalid, not a product failure.
- DONE: Record the exact hosted outcome without overstating its staged surface.
  The run passed the prior unowned journeys, skipped owned withdrawn/rejection, then auto-continue stopped in 12.58s because its staged single-root README was not commissioned.
- DONE: Confirm the static reachability defects independently.
  Both exact README definitions lacked `commissioned-by`, while the common runner called one non-Git single-root setup and never staged split-root.
- DONE: Add discovery controls red-first against real status discovery.
  Baseline focused run exited 1 because `status.DiscoverWorkflowDir` rejected both exact READMEs; deleting either added marker makes its control fail again.
- DONE: Add Git/state-root baseline controls with removal cases.
  Single-root and split-root checks require exact roots, existing entity, one seed commit, and clean status; raw no-Git staging for either variant is rejected.
- DONE: Add the exact common-runner reachability control red-first.
  Baseline fake-driver run observed one single-root launch, no Git, no discovery, and artifact label `auto-continue-after-implementation` without the fixture ID.
- DONE: Make both exact workflow READMEs discoverable.
  Single-root and split-root definitions now carry `commissioned-by: spacedock@1`; marker-removal controls prove the discovery predicate is causal.
- DONE: Stage both variants with proven Git ordering.
  Single-root initializes its definition root; split-root initializes state checkout first and definition root second, leaving both committed and clean.
- DONE: Run exactly both variants under one canonical journey.
  One top-level auto-continue runner invokes single-root then split-root serially; canonical scenario/assertion identity remains unchanged.
- DONE: Keep artifact evidence distinct without adapter or metric changes.
  Launch-only labels contain each exact fixture ID, so host artifact directories differ while Pi still aggregates results into the canonical journey metric.
- DONE: Preserve the durable oracle and prompt semantics.
  The unchanged neutral runbook and `assertAutoContinue` pass for both fake-driver durable transitions; stopped-after-implementation negatives remain red.
- DONE: Preserve pure desired state and exact observed bindings.
  Reconciliation/mutations prove 64 cells and eight TODOs remain Sonnet=3, Codex=4, Pi=1, Opus=0; auto-continue gains no owner binding.
- DONE: Reprove prior fixture parity, SHA, Pi rejection, and metric boundaries.
  Final focused live-tagged/offline, reconciliation/mutation/current-stale SHA, Pi round, identity/cost, and large-row commands exit 0 (`ensigncycle` 0.903s).
- DONE: Run required formatting and repository-wide suites on final replacement bytes.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 99.422s) and race exits 0 (`ensigncycle` 117.785s).
- DONE: Preserve exact watched-core/SHA-binder history after documentation correction.
  Initial `1dc8a19e2`/`b025eb033` was replaced by core `1fa6d52254cf7acb009e474b688e4861bb6fbc20` and SHA-only binder `37d91f673bda7ccdf9b7e63f6dc0134407895173`.
- DONE: Stay within surface, authority, and execution holds.
  Origin equals clean head; surface is 42 files/+2558/-671 with five existing Cycle-11 files, no product/adapter/launcher/runtime/new file, no state transition, and no live/model/hosted/substrate/smoke action or spend.

### Summary

Cycle 11 repairs only auto-continue fixture reachability: both commissioned Git-backed shapes now run serially through one canonical journey with distinct artifacts and the same durable oracle. The harness-invalid run creates no product evidence or owner, and focused/full/race/history/origin/ceiling proof is green on the pushed replacement pair.

## Stage Report: validation (cycle 11)

- DONE: Verify exact candidate provenance, history, and remote equality.
  Candidate/origin are clean at `37d91f673bda7ccdf9b7e63f6dc0134407895173`; core is `1fa6d52254cf7acb009e474b688e4861bb6fbc20` and merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`.
- DONE: Independently verify the supplied hosted artifact without launching a run.
  Run `31034783581`, job `92404363889`, artifact `8942417379`, and `source-head.txt` bind exact source `f390bf5d3ac37882db6fb1af9edb34a9c991c9ff`.
- DONE: Verify the exact hosted outcome and root cause.
  Nine runnable journeys passed, exact owned withdrawn/rejection cells skipped, and auto-continue alone failed in 12.58s after its sole boot-identify in fixture cwd correctly found no commissioned workflow.
- DONE: Preserve the hosted-versus-static evidence classification.
  The archived transcript shows no post-probe mutation; the missing marker and absent split-root/Git staging make the run harness-invalid, with no product evidence, repair owner, or TODO binding.
- DONE: Reproduce red baseline and green candidate discovery behavior independently.
  Both exact README definitions fail real `status.DiscoverWorkflowDir` at `f390bf5d3`; both resolve to their exact fixture roots on `37d91f673`, and removing either commissioned marker fails again.
- DONE: Verify exact Git and state-root baselines.
  Single-root uses its definition root; split-root uses `.spacedock-state` for the entity and state Git root; every repository has exactly one clean seed commit, while raw no-Git variants fail.
- DONE: Verify one canonical journey launches exactly both variants serially.
  The common runner launches `auto-continue/single-root` then `auto-continue/split-root`; removal, duplication, and reverse-order source mutants each fail the exact sequence control.
- DONE: Verify distinct artifact identity and unchanged semantic oracle.
  Both artifact labels contain their exact fixture IDs and a wrong-label mutant fails; neutral prompt and durable `assertAutoContinue` hashes are byte-identical to baseline, including the stopped-after-implementation negative.
- DONE: Verify adapters, product surfaces, selectors, and evidence state remain unchanged.
  Cycle 11 changes exactly five existing fixture/test/docs files, no host adapter/product skill/binary/launcher/runtime/workflow selector/new file; desired state remains pure 64 and evidence remains eight bindings (Sonnet=3, Codex=4, Pi=1, Opus=0).
- DONE: Reprove focused parity, fixture accountability, SHA, binding mutation, and Pi boundaries.
  Runner/definition/runtime parity, annotations/use, current+stale SHA, exact TODO mutations, Pi four-entry rejection recording, native identity/cost/complete-large-row metrics, and release controls all exit 0.
- DONE: Verify formatting, surface ceiling, and exact-byte repository suites.
  Read-only `gofmt -l ./cmd ./internal` is empty and `git diff --check` passes; surface is 42 files/+2558/-671; full passes (`ensigncycle` 103.856s, wall 115.74s) and race passes (`ensigncycle` 115.016s, wall 127.82s).
- DONE: Preserve candidate, runtime, and workflow authority boundaries.
  Validation changed no candidate bytes, transitioned no state, launched no live/model/hosted/substrate/smoke action, filed no owner, and made no product fix; PASSED with no material, deferred-risk, or polish findings.

### Summary

Cycle 11 PASSED independently at exact candidate `37d91f673bda7ccdf9b7e63f6dc0134407895173`. The correction repairs only fixture reachability and two-shape coverage, while the frozen hosted failure remains correctly classified as harness-invalid and all product/evidence boundaries stay unchanged.

## Stage Report: implementation (cycle 12)

- DONE: Start from the exact clean Cycle-11 candidate and inspect durable failure evidence.
  Baseline is `37d91f673bda7ccdf9b7e63f6dc0134407895173`; run `31038408585` lost reachable parent `8354e03b…` in `TestDurableQuestionedRejectsUnrelatedEdit`.
- DONE: Identify the exact rerun cleanup failure.
  Read-only log for run `31038934635` names `TestRecordedGateLifecycleEnteredStageRecoveryMatrix/committed_heading_only` and writes continuing under `.spacedock-state/.git/objects` during TempDir removal.
- DONE: Confirm the shared root cause and supported Git behavior.
  Git-heavy throwaway repositories cross automatic-maintenance thresholds; Git 2.54 enables detached geometric maintenance, while local Git 2.50.1 exposes the same child-launch seam in TRACE2.
- DONE: Inspect reconciliation watch scope before changing history.
  `internal/testgit` alone is not watched; the ceiling-preserving annotation relocation changes watched `internal/ensigncycle`, so an exact core plus SHA-only binder remains required.
- DONE: Add a focused version-neutral RED test first.
  Baseline `TestInitRepoDisablesAutomaticMaintenanceLocally` exited 1 at missing local `maintenance.auto`; the paired plain-init repo proves the assertion is not ambient config.
- DONE: Implement the smallest repo-local maintenance disablement.
  `InitRepo` stores `maintenance.auto=false` to block modern automatic maintenance and `gc.auto=0` as the older-Git fallback, alongside the existing local identity.
- DONE: Add causal behavioral discrimination without sleeps or detached work.
  A synchronous positive control sets strategy=gc, auto=true, autoDetach=false, gc.auto=1 and emits maintenance in TRACE2; default `InitRepo` emits neither maintenance nor gc.
- DONE: Prove the TRACE2 supplement executes rather than skips.
  `TestInitRepoSuppressesTraceableAutomaticMaintenance` passes ten consecutive runs on Git 2.50.1 with its positive child observed each time.
- DONE: Preserve the fixed 42-file ceiling while applying the shared fix.
  Restored annotation-only `recorded_gate_lifecycle_test.go` and `shallow_boot_fixture_live_test.go` to merge-base bytes, freeing exactly two paths for `testgit.go` and `testgit_test.go`.
- DONE: Preserve fixture bytes and improve common annotation ownership.
  Exact IDs `recorded-gate/prepared` and `boot/held-gate` moved to common wrappers in already-changed `shared_fixtures_test.go`; Claude/Codex call them and Pi retains the Claude path.
- DONE: Reprove fixture inventory and SHA discrimination.
  Reconciliation, mutation, current/stale SHA, shallow-boot shape, gate overlap, and Git scaffold controls exit 0; fixture count and desired/evidence registries are unchanged.
- DONE: Stress both previously failing tests serially.
  Exact tests pass once and 50× before commit, then 25× on final bytes with `-parallel=1` (`46.068s`) and no lost object or cleanup race.
- DONE: Run the complete shared-helper and focused regression fanout.
  `go test ./internal/testgit` and focused fixture/Pi rejection/native metric/large-row controls exit 0 on the final candidate.
- DONE: Run required formatting and repository-wide suites.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 109.100s, `testgit` 2.327s) and race exits 0 (`ensigncycle` 121.338s, `testgit` 2.234s).
- DONE: Preserve exact history, remote equality, and surface limits.
  Core `ba0fa02f233cfc3efec4e3bac87af9e61aadc275` is bound by docs-only `99f27c41c88bdd414240344495013e0c393a4c41`; origin equals clean head at exactly 42 files/+2668/-681.
- DONE: Preserve product, runtime, registry, owner, and execution boundaries.
  No product/live fixture bytes/prompt/adapter/workflow/TODO/owner changed, no file was added, no state transition occurred, and no hosted/live/model/substrate/smoke action ran.

### Summary

Cycle 12 hardens every throwaway repository through the shared `InitRepo` helper, preventing modern detached maintenance and legacy auto-GC during test cleanup. Config, synchronous TRACE2, exact-failure stress, reconciliation, full/race, history, origin, and 42-file ceiling proof are green with no product or live execution change.

## Stage Report: validation (cycle 12)

- DONE: Verify exact candidate provenance, history, and remote equality.
  Candidate/origin are clean at `99f27c41c88bdd414240344495013e0c393a4c41`; core is `ba0fa02f233cfc3efec4e3bac87af9e61aadc275`, baseline is `37d91f673bda7ccdf9b7e63f6dc0134407895173`, and merge base is `507a4bc12e48a3e4a813219602c488f09c81a5d8`.
- DONE: Independently verify both hosted offline failures and execution holds.
  Run `31038408585`/job `92416384344` lost reachable parent `8354e03b…`; run `31038934635`/job `92418150514` hit `.git/objects: directory not empty`; all live jobs were explicitly skipped in both runs.
- DONE: Verify the Git automatic-maintenance causal classification.
  Git docs establish `maintenance.auto=true`, detached `maintenance.autoDetach=true`, and 100-object auto thresholds by default; both archived failures are consistent with background writes in Git-heavy throwaway repositories.
- DONE: Reproduce the focused red baseline and exact green configuration.
  A detached baseline test fails because both local keys are absent; exact candidate `InitRepo` persists `maintenance.auto=false` and `gc.auto=0`, while a plain-init control persists neither.
- DONE: Perform independent per-line mutation attacks.
  Removing `maintenance.auto=false` or `gc.auto=0` separately in a detached candidate makes `TestInitRepoDisablesAutomaticMaintenanceLocally` fail on that exact missing key; restored bytes pass.
- DONE: Prove the maintenance boundary behaviorally without cleanup races.
  Ten runs of the synchronous TRACE2 control (`strategy=gc`, auto=true, autoDetach=false, gc.auto=1) launch maintenance, while default `InitRepo` emits neither maintenance nor gc.
- SKIPPED: Launch an `autoDetach=true` background control.
  Deliberately spawning detached maintenance would recreate the prohibited cleanup race and risks leaving work; the synchronous child-launch control proves the same causal seam safely without sleeps.
- DONE: Stress both prior failures and the shared helper serially.
  The exact two failure sites pass 25× combined with `-parallel=1` (`ensigncycle` 45.904s, wall 46.10s); `internal/testgit` passes 25× (`9.024s`, wall 9.14s) with no object loss or cleanup race.
- DONE: Verify annotation relocation is behavior- and byte-neutral.
  Restored recorded-gate and shallow-boot files hash exactly to merge base and are absent from final diff; each ID occurs once on a common wrapper used by Claude/Codex, with Pi routing through Claude.
- DONE: Reprove fixture, registry, SHA, mutation, and Pi evidence boundaries.
  Fixture parity/use, current+stale SHA, exact binding mutations, Pi four-entry rejection recording, native identity/cost/complete-large-row metrics, pure 64 desired cells, and eight bindings (Sonnet=3, Codex=4, Pi=1, Opus=0) all pass.
- DONE: Verify formatting, binder shape, surface ceiling, and exact-byte suites.
  Read-only `gofmt -l ./cmd ./internal` is empty and `git diff --check` passes; binder changes only the SHA; surface is 42 files/+2668/-681; full passes (`ensigncycle` 104.090s, wall 116.46s) and race passes (`ensigncycle` 114.547s, wall 127.99s).
- DONE: Preserve candidate, product, runtime, owner, and workflow authority boundaries.
  Validation changed no candidate, transitioned no state, launched/approved/rejected/cancelled no hosted/live/model/substrate/smoke action, filed no owner, and repaired no product; PASSED with no material, deferred-risk, or polish findings.

### Summary

Cycle 12 PASSED independently at exact candidate `99f27c41c88bdd414240344495013e0c393a4c41`. The shared throwaway-repository fix is causally discriminating and eliminates both observed Git cleanup signatures while preserving every fixture, registry, runtime, product, and execution boundary.

## Stage Report: implementation (cycle 13)

- DONE: Start from the exact clean and pushed Cycle-12 candidate.
  Baseline was `99f27c41c88bdd414240344495013e0c393a4c41`; no prior candidate bytes were rewritten.
- DONE: Inspect the exact hosted Codex evidence without rerunning it.
  Run `31042421641`, job `92430021762`, artifact `8945373289` was read from `/tmp/ys-cycle12-codex-retry-31042421641`.
- DONE: Preserve the hosted result's evidence classification.
  Nine runnable journeys passed and exact owned gaps skipped; only `auto-continue-after-implementation` failed, so this cycle changes harness truth only and creates no product or owner binding.
- DONE: Identify the stale single-root grading source.
  The root entity was nonempty at validation with `worktree: .worktrees/spacedock-ensign-auto-continue-task` but lacked a validation report; the active worktree entity contained the committed report and PASSED recommendation.
- DONE: Add the focused single-root RED before implementation.
  `TestResolveAutoContinueEndState/single-root_active_worktree` returned the stale pipeline body instead of the active worktree report on the Cycle-12 baseline.
- DONE: Correct the split-root control after review and capture a second RED.
  The adversarial file now occupies the exact probed path, `stateRoot/.worktrees/spacedock-ensign-auto-continue-task/auto-continue-task.md`; the unguarded resolver returned `WRONG CODE WORKTREE COPY`.
- DONE: Encode root authority explicitly in the common runner.
  Setup records whether the fixture's state root differs from its workflow root; only single-root grading follows a safe relative stamped worktree path.
- DONE: Preserve split-root and archive semantics.
  Split-root grading always trusts the separate state checkout, while empty post-state still falls back to both flat and folder archive shapes.
- DONE: Reprove the focused harness contracts.
  Resolver, common-runner, discovery, Git baseline, negative, parity, rejection, TODO, deterministic Pi metrics, reconciliation, mutation, and current/stale SHA controls exit 0.
- DONE: Keep the registry binder isolated from the watched core.
  Core `7ef1636c5ff751ff2e04716e39eb075cdd8aaf42` is followed only by docs binder `0aa74035cc96d672c71480466ed4e1b1a3169213`.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 103.072s) and race exits 0 (`ensigncycle` 115.215s).
- DONE: Preserve the fixed surface and remote equality.
  Candidate/origin are `0aa74035cc96d672c71480466ed4e1b1a3169213` at exactly 42 files/+2737/-681, below the +2750 ceiling.
- DONE: Preserve fixtures, prompts, adapters, registry, product, and workflow authority.
  No file was added, no runtime fixture/prompt/adapter/workflow/TODO/product byte changed, no owner was filed, and no state transition was made.
- DONE: Record the one verification execution anomaly exactly.
  An overbroad local name filter selected the opt-in Pi front-door smoke; the provider rejected it at its pre-execution 402 credit check, after which all filters were narrowed to deterministic tests and no hosted rerun was launched.

### Summary

Cycle 13 makes common auto-continue grading follow the durable active entity in single-root workflows while explicitly retaining state-checkout authority for split-root workflows. Both causal REDs, focused contracts, full/race suites, SHA binding, remote equality, and the 42-file ceiling are proved at candidate `0aa74035cc96d672c71480466ed4e1b1a3169213` without a product repair or workflow transition.

## Stage Report: validation (cycle 13)

- DONE: Verify exact immutable provenance, history, binder, and remote equality.
  Candidate and origin are clean at `0aa74035cc96d672c71480466ed4e1b1a3169213`; watched core is `7ef1636c5ff751ff2e04716e39eb075cdd8aaf42`, baseline is `99f27c41c88bdd414240344495013e0c393a4c41`, and the binder changes only `docs/runtime-live-ci.md` from the prior core SHA to the watched SHA.
- DONE: Independently inspect the hosted result without executing a runtime.
  Run `31042421641`, Codex job `92430021762`, and artifact `8945373289` have the assigned baseline head; nine runnable journeys passed and exact owned Codex gaps `withdrawn-gate-recovery`/`47gnqfm1ft6f2hcahz98m2jv` and `rejection-flow`/`zbcj98qfwtax61vxdzrf615e` skipped before the single-root auto-continue false negative.
- DONE: Prove the durable single-root lifecycle outcome from the archive.
  At validation the root entity was nonempty but had no validation report; the stamped active worktree entity had `started`, a complete committed validation report, and `Recommendation: PASSED`; the FO read that exact copy, prepared an open validation briefing, and presented the gate.
- DONE: Reproduce baseline RED and restored candidate GREEN independently.
  A detached Cycle-12 baseline resolver test stayed on the stale pipeline body and failed; the exact candidate passes the resolver, common-runner, discovery, and Git-baseline focused controls.
- DONE: Attack both root-authority branches causally.
  Removing the single-root overlay makes `single-root_active_worktree` fail on the missing report; bypassing `splitRoot` makes `split-root_remains_in_state_checkout` return the exact adversarial `WRONG CODE WORKTREE COPY`; restored bytes pass both.
- DONE: Verify resolver safety and terminal layouts.
  Absolute and parent-relative worktree values fall back to the pipeline body; split-root remains state-checkout authoritative, and both flat and folder archive layouts pass.
- DONE: Preserve fixture, prompt, oracle, and launch semantics.
  Both fixture IDs remain discoverable with clean committed Git/state-root baselines and run serially through the common runner; `auto_continue_fixtures_test.go` and `shared_scenarios_test.go` hashes are identical to baseline.
- DONE: Reprove reconciliation, mutations, SHA, Pi, and evidence cardinality.
  Focused ensigncycle/release/contractlint controls pass, including current/stale SHA, registry mutations, deterministic Pi evidence/metrics, the pure 64-cell desired matrix, and exactly eight target-scoped bindings (Sonnet 3, Codex 4, Pi 1, Opus 0).
- DONE: Verify formatting, exact surface, and repository-wide suites.
  Required `gofmt -w ./cmd ./internal` preserved identical bytes; `git diff --check` passes; full exits 0 (`ensigncycle` 103.771s, wall 116.14s) and race exits 0 (`ensigncycle` 115.010s, wall 128.43s); surface is exactly 42 files/+2737/-681, below +2750.
- DONE: Preserve scope and authority boundaries.
  Cycle 13 changes only two existing harness test files plus the docs-only SHA binder; no product, fixture, prompt, adapter, launcher, workflow, TODO, owner, or new-file change exists, and validation made no candidate or state transition.
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; the reported overbroad local filter reached only Pi's provider pre-execution 402 credit rejection, so no claim is made beyond that evidence and no hosted retry occurred.
- DONE: Render the independent verdict.
  PASSED with no Material, Deferred-risk, or Polish findings.

### Summary

Cycle 13 PASSED independently at exact candidate `0aa74035cc96d672c71480466ed4e1b1a3169213`. The narrow common-runner resolver fixes the archive-proven single-root grading false negative, rejects unsafe stamped paths, retains split-root and archive authority, and preserves every product, fixture, prompt, registry, owner, execution, and surface boundary.

## Stage Report: implementation (cycle 14)

- DONE: Start from the exact clean and pushed Cycle-13 candidate.
  Baseline and origin were `0aa74035cc96d672c71480466ed4e1b1a3169213`; watched core was `7ef1636c5ff751ff2e04716e39eb075cdd8aaf42`.
- DONE: Represent the exact deterministic hosted product failure.
  Runtime Live E2E run `31045591048`, Codex job `92440554439`, artifact `8946459512` failed `full-ensign-cycle` because fresh workers lacked the executable ensign contract bootstrap.
- DONE: Bind only the newly filed repair owner.
  Codex `full-ensign-cycle` now returns `TODO(nvz2ym82ydfn07jp04yfxg9r)` for backlog owner `restore-codex-ensign-contract-bootstrap`.
- DONE: Add the smallest exact controls before implementation.
  The TODO table expected the exact Codex owner and reconciliation expected the exact target/journey pair; both REDs failed on the missing binding before the switch case existed.
- DONE: Prove exact cardinality and wrong-target behavior.
  `TestSharedLiveTODOEvidenceSet` passes Sonnet, Opus, Codex, and Pi subtests at exactly Sonnet=3, Codex=5, Pi=1, Opus=0; moving the binding to any non-Codex target would make that target non-runnable and the Codex count short.
- DONE: Preserve all prior evidence bindings and mutation controls.
  Reconciliation logs exactly nine bindings, while owner, target, global, duplicate, proven-pass, and all pre-existing wrong-target mutation controls still fail their injected candidates.
- DONE: Preserve the normative desired registry.
  The 16 journeys × four required targets remain the pure 64-cell desired matrix; only observed target-scoped missing evidence changed.
- DONE: Update the operator evidence inventory in watched history.
  The inventory states nine bindings and cites run `31045591048`; no desired-registry, fixture, workflow, or runtime selector text changed.
- DONE: Preserve the core-plus-binder history boundary.
  Watched core `84693a0f4feb02512bdb7994663c406dd32a0929` is followed only by SHA binder `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5`.
- DONE: Reprove current and stale SHA discrimination.
  Before binding, the guard named both changed watched paths; after binding, `TestRuntimeLiveRegistryReconciliationSHA` passes its current-base and stale-parent assertions.
- DONE: Run required formatting and repository-wide suites.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 105.770s) and race exits 0 (`ensigncycle` 116.979s).
- DONE: Preserve exact surface and remote equality.
  Candidate/origin are `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5` at exactly 42 files/+2742/-681, eight insertions below the +2750 ceiling.
- DONE: Preserve every prohibited product and execution boundary.
  No product skill/binary, dispatch builder, fixture, prompt, adapter, launcher, oracle, Pi evidence/metric, workflow, other TODO, file, member, controller, or retry changed; no live/hosted/model/substrate/smoke command ran and no state transition was made.

### Summary

Cycle 14 records the archive-proven Codex `full-ensign-cycle` product failure as the ninth target-scoped missing-evidence binding and leaves the pure 64-cell desired registry intact. Exact RED/GREEN controls, mutation and SHA discrimination, full/race suites, remote equality, scope audit, and the fixed surface ceiling are green at `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5`.

## Stage Report: validation (cycle 14)

- DONE: Verify exact immutable provenance, history, binder, and remote equality.
  Candidate/origin are clean at `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5`; core `84693a0f4feb02512bdb7994663c406dd32a0929` changes four existing representation/test/inventory files, and the binder changes only the reconciliation SHA to that core.
- DONE: Independently verify the archived hosted product failure.
  Run `31045591048`, Codex job `92440554439`, artifact `8946459512` has baseline head `0aa74035cc96d672c71480466ed4e1b1a3169213`; first runnable `full-ensign-cycle` failed at 166.28s after the FO archived two generic no-colon checkbox reports and the anchored-report oracle correctly rejected them.
- DONE: Verify the missing executable bootstrap seam.
  The generated fresh-child file only claimed to contain shared ensign discipline and carried stage/entity/checklist/fetch/completion content; it contained no executable `ensign-shared-core` bootstrap, matching the observed generic worker output.
- DONE: Verify the exact repair owner and sprint boundary.
  Active `restore-codex-ensign-contract-bootstrap.md` is `status: backlog`, ID `nvz2ym82ydfn07jp04yfxg9r`, and `sprint:` is empty, so it exactly owns the product seam without being a sprint member.
- DONE: Prove the pure desired matrix and exact target-scoped evidence set.
  Registry source has exactly 16 common journeys and four required targets (64 cells); reconciliation logs exactly nine missing-evidence bindings with Sonnet=3, Codex=5, Pi=1, Opus=0, and only Codex/full-ensign-cycle maps to `nvz2ym82ydfn07jp04yfxg9r`.
- DONE: Preserve other targets and all prior evidence bindings.
  Sonnet, Pi, and Opus full-ensign-cycle remain runnable; the Cycle-14 core diff only inserts the new constant/case/expectation/map row, leaving all prior eight target/owner bindings byte-preserved and every TODO classified as missing evidence rather than pass or exception.
- DONE: Exercise removal, target, owner, global, duplicate, and proven-pass mutations causally.
  Detached removal left Codex short; Sonnet relocation suppressed its pass and left Codex short; wrong owner reported exact expected/got IDs; global form left the Codex cell missing; duplicate form emitted `DUPLICATE MISSING-EVIDENCE`; suppressing Sonnet/Pi/Opus full-cycle passes failed all three runnable-target assertions.
- DONE: Reproduce focused registry, TODO, mutation, and SHA controls.
  `TestSharedLiveTODOEvidenceSet` and the three contractlint reconciliation/mutation/SHA tests exit 0; current core-to-head watched diff is empty, while stale `84693a0f4^` names both changed ensigncycle paths.
- DONE: Verify the operator inventory.
  `docs/runtime-live-ci.md` states nine rows, cites run `31045591048` for Codex full-ensign-cycle, preserves the prior eight descriptions, and explicitly retains unverified Pi/Opus plus missing-evidence-not-pass semantics.
- DONE: Verify required formatting, exact surface, and repository-wide suites.
  `gofmt -w ./cmd ./internal` preserved exact bytes; `git diff --check` passes; full exits 0 (`ensigncycle` 102.953s, wall 115.46s) and race exits 0 (`ensigncycle` 114.538s, wall 127.97s); surface is exactly 42 files/+2742/-681, eight insertions below +2750.
- DONE: Preserve prohibited scope and execution boundaries.
  No skill/binary/dispatch/fixture/prompt/adapter/launcher/oracle/Pi metric/workflow/new-file/member/controller/retry byte changed; validation mutated no candidate, transitioned no state, and ran no hosted/live/model/substrate/smoke command.
- DONE: Render the independent verdict.
  PASSED with no Material, Deferred-risk, or Polish findings.

### Summary

Cycle 14 PASSED independently at exact candidate `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5`. The ninth binding accurately withholds only the archive-disproved Codex full-cycle evidence, preserves the pure 64-cell desired matrix and all prior evidence semantics, and remains causally guarded within the exact surface and authority limits.

## Stage Report: implementation (cycle 15)

- DONE: Start from the exact clean and pushed Cycle-14 candidate.
  Baseline and origin were `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5`; no root-main update or unrelated state write occurred.
- DONE: Classify the exact Pi evidence without rerunning it.
  Run `31048506420`, job `92450151141`, artifact `8947855457` passed offline and seven runnable journeys, skipped owned Pi rejection-flow, then false-negatived on filing despite the correct durable file.
- DONE: Identify the observer root cause.
  Pi archived a native assistant `bash` toolCall and correlated successful `role:toolResult` in `sessionJSONL`, while `runClaudeFilingScenario` graded Pi stdout/stderr with Claude's `tool_use` parser.
- DONE: Add the focused Pi-native RED before implementation.
  `TestAssertPiFilingViaNew` initially failed to compile because the Pi correlator did not exist; its distilled exact archived-success call/result shape is now the positive control.
- DONE: Route only Pi through native root-session evidence.
  `runtime == "pi"` selects `sessionJSONL` plus `assertPiFilingViaNew`; Claude retains `result.stream` plus its original parser and Codex retains its separate filing runner.
- DONE: Require exact successful atomic-create evidence.
  The correlator accepts exactly one requested-slug `bash` atomic-create call, exactly one later same-ID `toolResult`, exact `toolName: bash`, and explicit `isError: false`; the durable entity `os.Stat` assertion remains separate.
- DONE: Add causal ambiguity and manual-flow negatives.
  Missing/failed result, duplicate invocation, duplicate result, unrelated call reusing the filing ID, and a distinct non-atomic `--next-id` call each fail closed; all six named subtests pass.
- DONE: Preserve existing host-specific filing and Pi correlation behavior.
  Claude synthetic/PR399 negatives, Codex launcher/manual-flow negatives, and the full Pi rejection-round ambiguity matrix pass after extracting the shared native correlator.
- DONE: Preserve registry and evidence boundaries.
  Reconciliation, mutation, current/stale SHA, pure 64-cell desired state, and all nine TODO bindings remain green at Sonnet=3, Codex=5, Pi=1, Opus=0; no owner or TODO was added.
- DONE: Stay within the fixed surface ceiling without losing executable proof.
  Branch-added explanatory comments were compacted from already-changed harness/helper files; every affected executable control remains and the final surface is exactly 42 files/+2750/-684.
- DONE: Preserve watched-core plus SHA-only binder history.
  Watched core `9c38ebca94ad23d65cc8648e67d0bce6cea4fe94` is followed only by binder `6408e2ad658f05d9c990b88a46cc8d002203c729` changing the recorded SHA.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 107.390s) and race exits 0 (`ensigncycle` 119.521s).
- DONE: Preserve scope, execution, and authority boundaries.
  No product/fixture/prompt/launcher/front-door, native Pi metric/evidence beyond this observer, file/member/controller/retry, or other entity changed; no hosted/live/model/substrate/smoke command ran and no workflow transition was made.

### Summary

Cycle 15 grades Pi filing from its durable root-session tool evidence instead of stdout/stderr, while leaving Claude and Codex filing paths intact. The exact archived success, fail-closed mutation set, prior correlation suites, registry/TODO boundaries, full/race verification, pushed two-commit history, and hard 42-file/+2750 ceiling are green at `6408e2ad658f05d9c990b88a46cc8d002203c729`.

## Stage Report: validation (cycle 15)

- DONE: Verify exact provenance, history, binder, cleanliness, and remote equality.
  Candidate/origin are clean at `6408e2ad658f05d9c990b88a46cc8d002203c729`; baseline is `d7cdd4e91f6fcfdc7308138434908b37fa6e54a5`, watched core is `9c38ebca94ad23d65cc8648e67d0bce6cea4fe94`, and the following commit changes only the documented reconciliation SHA.
- DONE: Independently classify the archived hosted result without rerunning it.
  Run `31048506420`, Pi job `92450151141`, artifact `8947855457` passed offline and seven runnable common journeys, skipped only owned Pi rejection-flow, then produced a filing harness false negative despite the correct durable entity.
- DONE: Prove the exact archived native filing evidence.
  Root session `2026-08-05T21-39-49-697Z_019fd3de-4f41-70f8-8900-0e1c1ed84a17.jsonl` contains exactly one assistant `bash` call of `${SPACEDOCK_BIN:-spacedock} new wire-the-thing` and exactly one later same-ID `bash` result with `isError:false` and `created: .../wire-the-thing.md id=001`; no bash `--next-id` call exists.
- DONE: Verify runtime-specific evidence routing independently.
  A detached live-tagged fake-driver audit accepts valid Pi `sessionJSONL` despite invalid Claude stream, accepts valid Claude stream despite invalid Pi session, and rejects a Pi result with the wrong tool name; Codex retains its separate runner.
- DONE: Attack every correlator invariant causally.
  Removing call-ID uniqueness makes `reused_call_ID` pass incorrectly; removing the manual-flow rejection makes `manual_next-id` pass incorrectly; loosening result cardinality makes `duplicate_result` pass incorrectly; removing explicit error checking makes `failed_result` pass incorrectly; each restored candidate control passes.
- DONE: Reprove the complete host-specific evidence suites.
  The Pi filing positive and missing/failed/duplicate/reused-ID/manual negatives, Pi rejection-round ambiguity matrix, Claude synthetic and PR399 controls, and Codex filing/manual-flow controls all pass.
- DONE: Reprove registry truth, ownership, and SHA discrimination.
  The desired inventory remains exactly 16 journeys × four required targets = 64 cells; exactly nine missing-evidence bindings remain Sonnet=3, Codex=5, Pi=1, Opus=0; current SHA and stale-parent checks plus all reconciliation mutations pass.
- DONE: Verify scope, comment reclamation, and fixed surface.
  Cycle 15 changes only existing observer/test/helper files plus the SHA binder; reclaimed text is comment-only outside the filing observer route, and no product, fixture, prompt, launcher, front-door, metric, workflow, owner, or TODO changed. The cumulative surface is exactly 42 files/+2750/-684.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` is idempotent and `git diff --check` passes; full exits 0 in 153.60s and race exits 0 in 158.94s.
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; hosted inspection was read-only and all new runtime-routing evidence came from deterministic fake-driver or archived-session controls.
- DONE: Render the independent verdict.
  PASSED with no Material, Deferred-risk, or Polish findings; validation changed no candidate and made no workflow state transition.

### Summary

Cycle 15 PASSED independently at exact candidate `6408e2ad658f05d9c990b88a46cc8d002203c729`. Pi filing is now graded from one unique, later, successful native bash call/result pair while durable state remains separately required; causal negatives, prior host paths, pure registry truth, the fixed surface ceiling, and full/race suites all remain green.

## Stage Report: implementation (cycle 16)

- DONE: Start from the exact clean and pushed Cycle-15 candidate.
  Baseline and origin were `6408e2ad658f05d9c990b88a46cc8d002203c729`; no root-main update or unrelated state write occurred.
- DONE: Classify the exact Pi evidence without rerunning it.
  Run `31052117603`, job `92461666431`, artifact `8949169013`, extracted at `/tmp/ys-cycle15-pi-31052117603`, passed offline and all eight runnable journeys, including filing, and skipped only the owned Pi rejection-flow journey.
- DONE: Identify the shallow-boot observer false negative.
  The run failed shallow boot at 38.51s with `stream carried no assistant turns — nothing to check`; the durable `assertShallowBoot` check passed, Pi stdout held the exact greeting, and the native root session showed that greeting after only read/bash activity with no `subagent`, `member_spawn`, or `delegate` call.
- DONE: Isolate the runtime mismatch.
  `runClaudeShallowBootScenario` passed Pi stdout/stderr to Claude-only `assertNoTeamCreateBeforeGreet` and `assertShallowBootMeasured`, even though the structural evidence lives in Pi `sessionJSONL`.
- DONE: Add focused RED evidence before routing the fix.
  The native-session test first failed to compile with `undefined: assertPiShallowBootSession`; the fake-routing test then failed Pi with the exact non-Claude-stream error before the runtime route changed.
- DONE: Grade Pi shallow boot from the whole native root session.
  `assertPiShallowBootSession` requires assistant text containing both exact held-gate and engagement-hint lines and scans the entire session for forbidden `subagent`, `member_spawn`, or `delegate` tool calls.
- DONE: Close ordering and short-circuit gaps causally.
  Missing greeting, each forbidden pre-greeting call, a same-turn greeting followed by dispatch, and dispatch after greeting all fail; the last two were added after adversarial review exposed the initial early-return weakness.
- DONE: Prove runtime-specific routing with the full scenario runner.
  A fake Pi result succeeds from valid `sessionJSONL` despite `non-Claude Pi diagnostics` in its stream, while a fake Claude result succeeds from its original fixture stream despite deliberately invalid Pi-session data.
- DONE: Preserve durable and host-specific boundaries.
  The durable entity assertion remains separate and first; Claude keeps its original stream assertions and boot-window metrics, Codex keeps its separate runner, and Pi keeps its existing whole-run metrics.
- DONE: Preserve registry and evidence truth.
  The desired registry remains exactly 16 common journeys × four required targets = 64 cells, with exactly nine missing-evidence TODO bindings and no owner change.
- DONE: Stay within the fixed surface ceiling without deleting executable controls.
  Existing branch-added test registries and mutation maps were compacted; every affected control still runs, and the cumulative surface is exactly 42 files/+2748/-708.
- DONE: Preserve watched-core plus SHA-only binder history.
  Watched core `b01171a2bf50b05750a85edac65125f9cbc49416` is followed only by binder `37c7e2660caa6148a590b21671c2d5f3693b0c74`, which changes the recorded reconciliation SHA.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 104.219s) and race exits 0 (`ensigncycle` 119.608s).
- DONE: Preserve scope, execution, and authority boundaries.
  No product, fixture, prompt, adapter, launcher, workflow, owner, TODO, file/member/controller/retry behavior, or other entity changed; no hosted/live/model/substrate/smoke command ran and no workflow transition was made.

### Summary

Cycle 16 grades Pi shallow boot from complete native root-session structure while leaving durable state and the Claude, Codex, and Pi metric paths intact. Exact archived evidence, causal ordering negatives, fake-driver runtime routing, pure 64-cell registry truth, full/race suites, remote equality, and the fixed surface ceiling are green at `37c7e2660caa6148a590b21671c2d5f3693b0c74`.

## Stage Report: validation (cycle 16)

- DONE: Verify exact provenance, history, binder, cleanliness, and remote equality.
  Candidate/origin are clean at `37c7e2660caa6148a590b21671c2d5f3693b0c74`; baseline is `6408e2ad658f05d9c990b88a46cc8d002203c729`, watched core is `b01171a2bf50b05750a85edac65125f9cbc49416`, and the binder changes only the documented reconciliation SHA.
- DONE: Independently classify the archived hosted result without rerunning it.
  Run `31052117603`, Pi job `92461666431`, artifact `8949169013` passed offline and eight runnable common journeys, skipped only owned Pi rejection-flow, then shallow-boot false-negatived because Pi stdout/stderr reached Claude-only stream parsers.
- DONE: Prove the exact native root-session evidence over the whole session.
  The assigned session has 23 records, nine assistant messages, one assistant text block containing both exact held-review and engage lines, only `read`/`bash` assistant tool calls, and zero `subagent`, `member_spawn`, or `delegate` calls.
- DONE: Prove strict native evidence and ordering behavior.
  Focused controls accept read/bash then greet and reject missing greet plus every forbidden pre-greet, same-turn, and post-greet call; detached additions also reject malformed, non-assistant, stdout, wrong-record-type, and split-text lookalikes.
- DONE: Prove runtime-specific routing causally.
  Valid Pi session with invalid Claude stream and valid Claude stream with invalid Pi session both pass; disabling the Pi branch reproduces `stream carried no assistant turns`, while an early return at greeting makes the post-greet dispatch negative fail.
- DONE: Preserve durable state, metrics, and host boundaries.
  `assertShallowBoot` remains first and its broken-end-state controls pass; Claude retains its original no-TeamCreate, measured-stream, and boot-window metric controls, Pi retains whole-run metrics, and Codex remains separate.
- DONE: Preserve prior observers and registry truth.
  Pi filing and rejection ambiguity suites pass; the desired inventory remains 16 journeys × four targets = 64 cells with exactly nine bindings at Sonnet=3, Codex=5, Pi=1, Opus=0 and no owner or TODO change.
- DONE: Verify scope, compacted controls, and fixed surface.
  Cycle 16 modifies three existing observer/test files plus the SHA binder, adds no file, and removes no test function; no product, fixture, prompt, adapter, launcher, workflow, owner, or TODO behavior changed. Surface is exactly 42 files/+2748/-708.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` is idempotent and `git diff --check` passes; full exits 0 in 153.46s and race exits 0 in 157.87s.
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; hosted inspection was read-only and all causal runtime checks used deterministic fake drivers or archived evidence.
- DONE: Render the independent verdict.
  PASSED with no Material, Deferred-risk, or Polish findings; validation changed no candidate and made no workflow state transition.

### Summary

Cycle 16 PASSED independently at exact candidate `37c7e2660caa6148a590b21671c2d5f3693b0c74`. The Pi shallow-boot observer now requires the complete native assistant greeting and rejects forbidden dispatch anywhere in the root session, while durable state, host-specific metrics, prior observers, registry truth, and the fixed surface remain intact.

## Stage Report: implementation (cycle 17)

- DONE: Start from the exact clean and pushed Cycle-16 candidate.
  Baseline and origin were `37c7e2660caa6148a590b21671c2d5f3693b0c74`; no root-main update or unrelated state write occurred.
- DONE: Classify the exact hosted evidence without rerunning it.
  Isolated Pi run `31054684260`, job `92469677841`, artifact `8949983771` at `/tmp/ys-cycle16-pi.wa29wg` passed offline, full-ensign-cycle in 121.89s, and gate-guardrail in 94.80s before default-headless-gate-stop failed at 228.61s; smoke was skipped and upload passed.
- DONE: Bind the failure to the existing product owner.
  The root advanced implementation to validation and dispatched a validator; the child prepared one correct open gate and report but used raw Git commit instead of `spacedock state commit`, crossing the no-authority boundary without decision, consume, or advance, which matches owner `26nk8qd48zknqnn4kc123sez`.
- DONE: Preserve product ownership and behavior.
  No new defect was filed and existing `headless-recorded-gate-stop-stage-coherence` remained unchanged; only target-scoped representation changed.
- DONE: Add exact failing controls before the source binding.
  The contract reconciliation first failed with `MISSING-EVIDENCE journey=default-headless-gate-stop target=pi owner=26nk8qd48zknqnn4kc123sez`, and the live runner test first reported an empty Pi reason for that exact target/journey.
- DONE: Add only the Pi default-headless source binding.
  `liveEvidenceTargetPi` plus `default-headless-gate-stop` now returns owner `26nk8qd48zknqnn4kc123sez`; every prior source binding is preserved.
- DONE: Preserve the pure desired matrix and exact observed counts.
  Controls prove 16 common journeys × four required targets = 64 cells and exactly ten missing-evidence bindings: Sonnet=3, Codex=5, Pi=2, Opus=0.
- DONE: Prove target and owner discrimination causally.
  Removal and retagging of either Sonnet or Pi default-headless fail; adding Codex fails its proven-pass control, adding Opus fails its unverified-target control, and the runner expectation independently keeps Codex/Opus runnable while retaining Sonnet.
- DONE: Update the watched operator inventory.
  `docs/runtime-live-ci.md` records the run, job, artifact, exact Pi route, existing owner, ten-result count, Codex passing evidence, and Opus runnable truth while retaining all prior nine descriptions.
- DONE: Stay inside the fixed surface without deleting executable proof.
  Repeated branch-added mutation registry bodies were factored/compacted; all controls still execute, no file was added, and cumulative surface is exactly 42 files/+2749/-708.
- DONE: Preserve watched-core plus SHA-only binder history.
  Core `59338039470002bb8dd1e6583edd2050ec625db1` is followed only by binder `23e0dbd02b3c4aa8753776596eadf9980698a0b5`, which changes the reconciliation SHA; the stale guard names both watched ensigncycle paths and the current guard passes.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 105.103s) and race exits 0 (`ensigncycle` 117.086s).
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; all new proof used deterministic source, mutation, registry, and archived-evidence controls.
- DONE: Preserve scope and authority boundaries.
  No runtime, product, fixture, prompt, adapter, launcher, front-door, observer, metric, workflow, owner, other TODO, file/member/controller/retry behavior, or other entity changed, and no workflow transition was made.

### Summary

Cycle 17 records the archive-proven Pi default-headless product failure against its existing owner while preserving the pure 64-cell desired matrix and all prior evidence truth. Exact RED/GREEN target and owner controls, current/stale SHA discrimination, full/race suites, pushed two-commit history, and the hard surface ceiling are green at `23e0dbd02b3c4aa8753776596eadf9980698a0b5`.

## Stage Report: validation (cycle 17)

- DONE: Verify exact provenance, history, binder, cleanliness, and remote equality.
  Candidate/origin are clean at `23e0dbd02b3c4aa8753776596eadf9980698a0b5`; baseline is `37c7e2660caa6148a590b21671c2d5f3693b0c74`, watched core is `59338039470002bb8dd1e6583edd2050ec625db1`, and the binder changes only the documented reconciliation SHA.
- DONE: Independently classify the archived hosted evidence without rerunning it.
  Run `31054684260`, Pi job `92469677841`, artifact `8949983771` passed offline, full-ensign, and gate-guardrail before default-headless failed at 228.61s with `gate hold crossed its committed no-authority boundary`.
- DONE: Prove the failure matches the existing owner rather than capacity or harness drift.
  Root-session evidence advances to validation and dispatches a child; the child prepares one valid open gate/report, then raw-Git commits state instead of using `spacedock state commit`; root validation leaves the gate open, unconsumed, and without successor advance, exactly matching owner `26nk8qd48zknqnn4kc123sez`.
- DONE: Verify the exact representation-only source change.
  The core adds only Pi/default-headless to the existing 26n switch and expectations, compacts existing mutation tables, and updates operator evidence; all prior nine target/owner bindings are byte-preserved.
- DONE: Prove pure desired registry and exact observed evidence.
  Focused runner/reconciliation controls retain 16 journeys × four targets = 64 desired cells and log exactly ten bindings at Sonnet=3, Codex=5, Pi=2, Opus=0.
- DONE: Prove target and owner discrimination causally.
  Detached source mutations make removal and retag of both Sonnet and Pi fail their exact owner checks; adding the same skip to proven Codex or runnable Opus fails those targets, while restored candidate tests pass all four targets.
- DONE: Preserve every compacted registry mutation control.
  The verbose mutation suite executes all prior structural, desired-registry, target/owner, global, duplicate, and proven-pass cases plus Pi and Opus additions; direct delete/set helpers cannot mask a no-op because unchanged evidence makes reconciliation return nil and the subtest fail.
- DONE: Verify operator text and scope boundaries.
  Inventory text names the exact run/job/artifact, child raw-Git seam, open/unconsumed gate, existing owner, Codex pass, and Opus runnable truth; no product, fixture, prompt, adapter, launcher, front-door, observer, metric, workflow, owner, or unrelated TODO behavior changed.
- DONE: Verify fixed surface and retained tests.
  Cycle 17 modifies four existing files, adds no file, removes no test function, and keeps the cumulative surface exactly 42 files/+2749/-708.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` is idempotent and `git diff --check` passes; full exits 0 in 154.48s and race exits 0 in 159.30s.
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; hosted metadata inspection was read-only and all causal proof used archived sessions or deterministic tests.
- DONE: Render the independent verdict.
  PASSED with no Material, Deferred-risk, or Polish findings; validation changed no candidate and made no workflow state transition.

### Summary

Cycle 17 PASSED independently at exact candidate `23e0dbd02b3c4aa8753776596eadf9980698a0b5`. The archive-proven Pi default-headless failure is accurately attached to the existing 26n owner, while target discrimination, pure registry truth, prior evidence, operator documentation, and the fixed surface remain intact.

## Stage Report: implementation (cycle 18)

- DONE: Start from the exact clean and pushed Cycle-17 candidate.
  Baseline and origin were `23e0dbd02b3c4aa8753776596eadf9980698a0b5`; no root-main update or unrelated state write occurred.
- DONE: Classify the exact hosted evidence without rerunning it.
  Pi run `31056949526`, job `92476577900`, artifact `8951115959` passed 11 runnable journeys before `smallest-sufficient-mechanism` failed; the archived run alone supplied product evidence.
- DONE: Bind the failure to the existing product owner.
  The run changed `ready-one` and `ready-two` directly from ready to done, committed entry into done, and built stage done, exactly matching existing owner `9adv48yhye5s2vkhwd7ge52d`; no owner was changed or filed.
- DONE: Add exact failing controls before the source binding.
  The live runner first reported an empty TODO reason for Pi/smallest-sufficient-mechanism, and reconciliation first emitted exact `MISSING-EVIDENCE` for target Pi and owner `9adv48yhye5s2vkhwd7ge52d`.
- DONE: Add only the Pi smallest-sufficient source binding.
  The existing multi-target smallest-sufficient case now includes `liveEvidenceTargetPi`; every prior ten target/owner binding remains intact.
- DONE: Preserve the pure desired matrix and exact observed counts.
  Controls prove 16 common journeys × four required targets = 64 cells and exactly eleven missing-evidence bindings: Sonnet=3, Codex=5, Pi=3, Opus=0.
- DONE: Prove target, owner, and runnable-Opus discrimination.
  Removing or retagging the Pi binding fails its exact check, adding the same owner to Opus fails the unverified-target control, and the four-target runner expectation preserves existing Sonnet/Codex bindings plus runnable Opus.
- DONE: Update the watched operator inventory.
  `docs/runtime-live-ci.md` records run `31056949526`, job `92476577900`, artifact `8951115959`, the exact two-entity ready-to-done behavior, existing repair owner, eleven-result count, and runnable Opus truth.
- DONE: Stay inside the fixed surface without deleting executable proof.
  Repeated branch-added registry string mutations now use one shared constructor; all mutation cases still execute, no file was added, and cumulative surface is exactly 42 files/+2748/-708.
- DONE: Preserve watched-core plus SHA-only binder history.
  Core `0e7f2379f9d63896f003a0f10d3a5beac1c31670` is followed only by binder `00fcde8a5121fe937a75d152cb7b498f1d16a10a`, which changes the reconciliation SHA; stale and current SHA controls pass their opposite expectations.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; under heavy host contention full exits 0 (`ensigncycle` 373.660s) and race exits 0 (`ensigncycle` 430.431s).
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; all new proof used archived evidence and deterministic source, mutation, registry, and SHA controls.
- DONE: Preserve scope and authority boundaries.
  No product, fixture, prompt, adapter, launcher, front-door, runtime observer/metric, workflow, owner, unrelated TODO, file/member/controller/retry behavior, or other entity changed, and no workflow transition was made.

### Summary

Cycle 18 records the archive-proven Pi smallest-sufficient-mechanism failure against its existing `9adv…` owner while preserving the pure 64-cell desired matrix and all prior evidence truth. Exact RED/GREEN owner and target controls, runnable Opus, SHA discrimination, full/race suites, pushed two-commit history, and the hard surface ceiling are green at `00fcde8a5121fe937a75d152cb7b498f1d16a10a`.

## Stage Report: validation (cycle 18)

- DONE: Verify exact provenance, history, binder, cleanliness, and remote equality.
  Candidate/origin are clean at `00fcde8a5121fe937a75d152cb7b498f1d16a10a`; baseline is `23e0dbd02b3c4aa8753776596eadf9980698a0b5`, watched core is `0e7f2379f9d63896f003a0f10d3a5beac1c31670`, and the binder changes only the documented reconciliation SHA.
- DONE: Independently classify the archived hosted evidence without rerunning it.
  Run `31056949526`, Pi job `92476577900`, artifact `8951115959` passed offline and 11 runnable journeys, skipped only the two prior owned Pi cases, then failed smallest-sufficient at 186.96s with a completed process and no timeout; smoke remained skipped.
- DONE: Prove the failure matches the existing owner rather than capacity or harness drift.
  Native root-session commands changed both `ready-one` and `ready-two` directly from ready to done, built and dispatched stage done, and left no dispatchable entities; child sessions wrote ready reports, confirming the missing path-scoped ready-stage dispatch entry exactly matches owner `9adv48yhye5s2vkhwd7ge52d`.
- DONE: Verify the representation-only source change and prior observers.
  The core adds only Pi/smallest-sufficient to the existing 9a source case and exact expectation, compacts equivalent registry mutations, and updates operator evidence; Pi rejection ambiguity and metric-attribution controls remain green.
- DONE: Prove pure desired registry and exact observed evidence.
  Focused runner/reconciliation controls retain 16 journeys × four targets = 64 desired cells and log exactly 11 bindings at Sonnet=3, Codex=5, Pi=3, Opus=0; current and stale SHA guards pass opposite expectations.
- DONE: Prove target and owner discrimination causally.
  Detached source mutations make Pi removal, wrong owner, wrong target, and Opus suppression fail their exact checks; restoring the exact candidate returns all four target subtests to green with a clean worktree.
- DONE: Verify operator text, scope, and fixed surface.
  Text names the exact run/job/artifact and two ready-to-done transitions; Cycle 18 changes four existing files, adds no file, removes no test function, changes no runtime/product/fixture/prompt/adapter/workflow/owner behavior, and retains exactly 42 files/+2748/-708 cumulatively.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` is idempotent and `git diff --check` passes; full exits 0 in 611.04s and race exits 0 in 618.97s.
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; hosted inspection was read-only and causal proof used archived native sessions or deterministic tests.
- DONE: Render the independent verdict.
  PASSED. Material: none. Deferred-risk: none. Polish: none. Validation changed no candidate and made no workflow state transition.

### Summary

Cycle 18 PASSED independently at exact candidate `00fcde8a5121fe937a75d152cb7b498f1d16a10a`. Pi's archive-proven smallest-sufficient failure is accurately attached to the existing `9adv…` owner, while target discrimination, runnable Opus, pure registry truth, prior observers, operator documentation, and the fixed surface remain intact.

## Stage Report: implementation (cycle 19)

- DONE: Start from the exact clean and pushed Cycle-18 candidate.
  Baseline and origin were `00fcde8a5121fe937a75d152cb7b498f1d16a10a`; no root-main update or unrelated state write occurred.
- DONE: Classify the exact hosted evidence without rerunning it.
  Pi run `31060796822`, job `92488310460`, artifact `8952141495` passed full-ensign-cycle, then gate-guardrail selected untracked `evidence/command.log` with committed review sources and halted.
- DONE: Bind the failure to the existing product owner.
  Generic exact-committed-file rejection was pathless and left no prepared room, matching existing owner `3zzpdw704df1g8pg1x9thzmw`; no owner was changed or filed.
- DONE: Add exact failing controls before the source binding.
  The runner first reported an empty Pi/gate-guardrail TODO reason, and reconciliation emitted exact missing evidence for target Pi and owner `3zzpdw704df1g8pg1x9thzmw`.
- DONE: Add only the Pi gate-guardrail source binding.
  `liveEvidenceTargetPi` plus `gate-guardrail` now returns owner `3zzpdw704df1g8pg1x9thzmw`; every prior eleven binding remains intact.
- DONE: Preserve the pure desired matrix and exact observed counts.
  Controls prove 16 journeys × four required targets = 64 cells and exactly twelve bindings: Sonnet=3, Codex=5, Pi=4, Opus=0.
- DONE: Prove target, owner, and runnable-Opus discrimination.
  Removal and owner retag of Pi/gate-guardrail fail directly; adding the binding to Opus fails its unverified-target control, while the four-target expectation preserves all prior bindings.
- DONE: Update the watched operator inventory.
  The inventory names run/job/artifact, untracked `command.log`, committed companions, pathless rejection, no-room result, existing owner, twelve-result count, and runnable Opus.
- DONE: Stay inside the fixed surface without deleting executable proof.
  Existing branch-added mutation constructors now cover counted and replace-all variants and repeated map bodies; every control executes, no file was added, and surface is exactly 42 files/+2749/-708.
- DONE: Preserve watched-core plus SHA-only binder history.
  Core `1cb4585ceb5ec77b29b113d87b95eedbfcb82dc2` is followed only by binder `2240fabe5d975a1ba724848ac2fc91ef11eb2cc7`, which changes the reconciliation SHA; stale and current controls pass their opposite expectations.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; full exits 0 (`ensigncycle` 180.994s) and race exits 0 (`ensigncycle` 206.493s).
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; proof used archived evidence and deterministic registry, mutation, runner, and SHA controls.
- DONE: Preserve scope and authority boundaries.
  No product, fixture, prompt, adapter, launcher, observer/metric, workflow, owner, unrelated TODO, file/member/controller/retry behavior, or other entity changed, and no transition was made.

### Summary

Cycle 19 records the archive-proven Pi gate-guardrail failure against its existing `3zz…` owner while preserving pure registry truth and runnable Opus. Exact RED/GREEN controls, SHA discrimination, full/race suites, pushed two-commit history, and the surface ceiling are green at `2240fabe5d975a1ba724848ac2fc91ef11eb2cc7`.

## Stage Report: validation (cycle 19)

- DONE: Verify exact provenance, history, binder, cleanliness, and remote equality.
  Candidate/origin are clean at `2240fabe5d975a1ba724848ac2fc91ef11eb2cc7`; baseline is `00fcde8a5121fe937a75d152cb7b498f1d16a10a`, watched core is `1cb4585ceb5ec77b29b113d87b95eedbfcb82dc2`, and the binder changes only the documented reconciliation SHA.
- DONE: Independently classify the archived hosted evidence without rerunning it.
  Run `31060796822`, Pi job `92488310460`, artifact `8952141495` passed offline and full-ensign-cycle, then gate-guardrail failed at 101.98s; the process completed without timeout, smoke was skipped, and artifact upload passed.
- DONE: Prove the exact existing 3z product seam rather than harness or capacity.
  Native root-session evidence selects three committed review sources plus untracked `evidence/command.log`; the real prepare call returns the generic pathless exact-committed-file rejection, creates no room, leaves readiness at needs-preparation, and halts, matching owner `3zzpdw704df1g8pg1x9thzmw`'s attributed-source repair.
- DONE: Verify the representation-only change and all prior bindings.
  The core adds only Pi/gate-guardrail to the source case and exact expectation, compacts equivalent mutation helpers, and updates operator evidence; prior observer/metric suites and all eleven prior target-owner bindings remain green.
- DONE: Prove pure desired registry and exact observed evidence.
  Focused runner/reconciliation controls retain 16 journeys × four targets = 64 desired cells and log exactly 12 bindings at Sonnet=3, Codex=5, Pi=4, Opus=0.
- DONE: Prove target, owner, and SHA discrimination causally.
  Detached source mutations make Pi removal, wrong owner, wrong target, and Opus suppression fail their exact checks; restoring the candidate returns all targets green. The current core has no watched delta, while the stale baseline names both changed ensigncycle paths.
- DONE: Verify scope, retained controls, and fixed surface.
  Cycle 19 changes four existing files, adds no file, removes no test function, changes no product/fixture/prompt/adapter/workflow/owner behavior, and retains exactly 42 files/+2749/-708 cumulatively.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` is idempotent and `git diff --check` passes; full exits 0 in 276.71s and race exits 0 in 281.56s.
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; hosted inspection was read-only and causal proof used archived native evidence or deterministic tests.
- DONE: Render the independent verdict.
  PASSED. Material: none. Deferred-risk: none. Polish: none. Validation changed no candidate and made no workflow state transition.

### Summary

Cycle 19 PASSED independently at exact candidate `2240fabe5d975a1ba724848ac2fc91ef11eb2cc7`. Pi's archived gate-guardrail failure is accurately attached to the existing `3zz…` path-attribution owner, while target discrimination, runnable Opus, pure registry truth, prior evidence, operator documentation, and the fixed surface remain intact.

## Stage Report: implementation (cycle 20)

- DONE: Start from the exact clean pushed Cycle-19 candidate.
  Baseline and origin were `2240fabe5d975a1ba724848ac2fc91ef11eb2cc7`; only two unpublished commits were rebuilt after explicit authorization.
- DONE: Classify the exact hosted evidence without rerunning it.
  Pi run `31062662883`, job `92493828626`, artifact `8953098453` passed ten runnable journeys, skipped four owned Pi gaps, then failed keep-moving-posture.
- DONE: Bind only Pi keep-moving to the existing owner.
  Pi set questioned `verdict=revise`, later set `status=review` without clearing the verdict, and left durable state terminal; this matches existing owner `9adv48yhye5s2vkhwd7ge52d`.
- DONE: Add RED-first exact controls.
  Runner and reconciliation first failed only Pi/keep-moving with exact owner `9adv…`; removal, retag, and runnable-Opus wrong-target controls pass after the single source binding.
- DONE: Preserve pure desired and observed truth.
  The desired matrix remains 16×4=64; exactly thirteen target-scoped bindings remain Sonnet=3, Codex=5, Pi=5, Opus=0, and TODO remains explicitly non-passing evidence.
- DONE: Preserve the durable Captain completion contract.
  The 692ae776f audit confirms every nonpassing desired cell has exact target-scoped TODO(owner), unbound targets remain runnable, and owners 26n/9a/nv/47g/zbc/3z remain unchanged; rm stays deferred.
- DONE: Resolve the registry/guide documentation finding.
  Captain status audit found `docs/runtime-live-ci-registry.md` still claimed reconciliation and SHA were future work although the guide ships both; classification was Material, ys-owned, disposition fix.
- DONE: Prove the documentation correction RED/GREEN.
  `TestRuntimeLiveRegistryGuideStateCurrent` first failed both stale phrases, then passed after only lines 7-8 and 349-352 were rewritten to present-state reconciliation/SHA/guard truth.
- DONE: Record the Captain-visible surface exception.
  No prior doc could be sacrificed: the principles file supplies six IDs required for the exact 16-journey seed, architecture notes preserve adapter truth, and the guide is normative. Captain authorized 43 files while retaining +2750 and no-new-file boundaries.
- DONE: Stay inside the authorized surface.
  Existing move mutations were factored through one helper with every case retained; final cumulative surface is exactly 43 files/+2744/-714 and adds no file.
- DONE: Preserve core plus SHA-only binder history.
  Core `b42835481c6944e239129a78bbef0eca240a8473` is followed only by binder `0f9e096e7f534f702c31fdcaa43aa1f416c2784f`, changing the guide SHA; stale/current controls pass opposite expectations.
- DONE: Run required verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean; focused coherence/registry/mutation/runner/SHA controls pass; full exits 0 (`ensigncycle` 303.903s) and race exits 0 (`ensigncycle` 282.025s).
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; proof used archived evidence and deterministic tests.
- DONE: Preserve scope and authority boundaries.
  No product, fixture, prompt, adapter, launcher, observer/metric, workflow, owner, unrelated TODO, member/controller/retry behavior, or transition changed.

### Summary

Cycle 20 binds Pi keep-moving to existing `9adv…`, fixes the Material registry/guide time-state contradiction, and preserves the Captain completion contract. The authorized 43/+2744/-714 candidate, exact RED/GREEN controls, full/race suites, and pushed two-commit history are green at `0f9e096e7f534f702c31fdcaa43aa1f416c2784f`.

## Stage Report: validation (cycle 20)

- DONE: Verify exact provenance, history, binder, cleanliness, and remote equality.
  Candidate/origin are clean at `0f9e096e7f534f702c31fdcaa43aa1f416c2784f`; baseline is `2240fabe5d975a1ba724848ac2fc91ef11eb2cc7`, watched core is `b42835481c6944e239129a78bbef0eca240a8473`, and the binder changes only the documented reconciliation SHA.
- DONE: Independently classify the archived hosted evidence without rerunning it.
  Run `31062662883`, Pi job `92493828626`, artifact `8953098453` passed offline and ten runnable journeys, skipped four owned Pi gaps, then failed keep-moving at 377.86s; the process completed without timeout, smoke was skipped, and upload passed.
- DONE: Prove the exact existing 9a product seam rather than harness or capacity.
  Native root-session state ends with `questioned` at `status=review`, empty completion/archive, and stale `verdict=revise`; the durable oracle treats any verdict as terminal and rejects that history, matching existing owner `9adv48yhye5s2vkhwd7ge52d`'s post-gate terminalization repair.
- DONE: Verify the single source binding and durable completion semantics.
  The core adds only Pi/keep-moving to the existing 9a case and exact expectation; all prior twelve bindings remain unchanged, the 16×4 desired matrix yields exactly thirteen bindings at Sonnet=3, Codex=5, Pi=5, Opus=0, and target-scoped TODO(owner) remains explicitly non-green evidence while unbound targets remain runnable.
- DONE: Verify the authorized documentation correction.
  Registry future-state claims are replaced by present-state reconciliation, SHA, and watched-guard truth without changing normative-guide ownership or mechanism scope; the coherence test is green and a detached restoration of the stale claim fails it exactly.
- DONE: Prove target, owner, Opus, and SHA discrimination causally.
  Detached source mutations make Pi removal, wrong owner, wrong target, and Opus suppression fail their exact checks; restoration returns all target tests green. Current watched paths are clean after core, while stale baseline names the registry and both ensigncycle paths.
- DONE: Audit the Captain-authorized surface exception and retained documentation.
  The principles seed remains 16/16 with code, architecture retains the shared adapter truth, and the operating guide remains normative; Cycle 20 adds no file or removes any test, stays below +2750, and lands exactly 43 files/+2744/-714 cumulatively.
- DONE: Run required formatting and repository-wide verification.
  `gofmt -w ./cmd ./internal` is idempotent and `git diff --check` passes; full exits 0 in 240.94s and race exits 0 in 245.48s.
- SKIPPED: Run hosted/live/model/substrate/smoke execution.
  Execution was prohibited; hosted inspection was read-only and causal proof used archived native evidence or deterministic tests.
- DONE: Render the independent verdict.
  PASSED. Material: none. Deferred-risk: none. Polish: none. Validation changed no candidate and made no workflow state transition.

### Summary

Cycle 20 PASSED independently at exact candidate `0f9e096e7f534f702c31fdcaa43aa1f416c2784f`. Pi keep-moving is accurately attached to the existing `9adv…` owner, the Material registry time-state defect is closed with causal proof, and the Captain-authorized 43-file surface preserves completion semantics, runnable Opus, documentation roles, and all prior evidence.

## Stage Report: validation (cycle 21)

- DONE: Revalidate the complete Pi target locally before any hosted execution.
  Exact candidate `0f9e096e7f534f702c31fdcaa43aa1f416c2784f` ran the shared 16-journey selector with isolated Pi 0.80.10, pinned `pi-subagents` 0.35.1 and `pi-intercom` 0.6.0, and subscription model `openai-codex/gpt-5.4`. The run exited 0 in 1531.85s: all 11 runnable journeys passed, five cells skipped under four exact owners, and final `ac-value-reanchor` passed after the new keep-moving binding.
- DONE: Run the Claude Sonnet target locally through the supported subscription path.
  The first package launch stopped before model execution because npm had left its documented native placeholder; running the package installer produced Claude Code 2.1.161. The exact selector then used `SPACEDOCK_LIVE_MODEL=sonnet`, resolved `claude-sonnet-4-6`, and retained artifacts under `/tmp/spacedock-cycle20-sonnet-local.fkxGfs`.
- DONE: Classify the runnable Sonnet evidence before changing candidate bytes.
  `full-ensign-cycle`, `gate-guardrail`, `withdrawn-gate-recovery`, and `recorded-gate-lifecycle` passed; `default-headless-gate-stop` skipped under exact owner `26nk8qd48zknqnn4kc123sez`. In `rejection-flow`, the validation ensign committed the REJECTED report and JSONL package, then `SendMessage(to="team-lead")` returned `No agent named 'team-lead' is currently addressable`; no further stream progress occurred for 60 seconds and the runner killed the process at 419.27s.
- DONE: Prove the Sonnet failure is product conduct with an existing owner, not auth, capacity, fixture, or oracle drift.
  The retained stream shows successful launch, delegation, durable validation mutation, and commit before the addressability failure. The FO never received completion, never recorded the complete correction round, never routed rework, and never ran cycle-2 validation. Archived owner `zbcj98qfwtax61vxdzrf615e` explicitly owns cross-runtime rejection-flow rework/revalidation and already records Sonnet and Opus triggers; it remains in `test-behavior-completeness`.
- DONE: Preserve the durable completion rule and route only representation truth.
  Sonnet `rejection-flow` is currently runnable without passing evidence, so completion is blocked. The correction is exactly one target-scoped `TODO(zbcj98qfwtax61vxdzrf615e)` binding plus count/target/owner/runnable-Opus mutation proof and operator inventory; TODO remains non-green. No product, prompt, fixture, adapter, retry, controller, owner, or other target binding may change.
- DONE: Preserve candidate and sprint boundaries during classification.
  Candidate and origin remain clean and equal at `0f9e096e7f534f702c31fdcaa43aa1f416c2784f`; no product repair, hosted CI, Opus run, Codex run, new owner, new member, or rm work occurred.
- FAILED: Recommend PASSED only if every runnable target cell has passing evidence or an exact reproduced owner binding.
  REJECTED. Sonnet `rejection-flow` has reproduced product-owned nonpassing evidence but no current target-scoped binding, so the ys representation is incomplete until the one-cell zbc correction is independently validated.

### Summary

Cycle 21 confirms the Cycle-20 Pi correction locally, then reproduces a Sonnet rejection-flow addressability/continuation failure after durable validation work. The failure matches existing cross-runtime owner `zbc…`; validation is REJECTED for one representation-only Sonnet binding, with product behavior and all unrelated cells held fixed.

## Stage Report: implementation (cycle 21)

- DONE: Start from the exact clean pushed Cycle-20 candidate.
  Baseline and origin were `0f9e096e7f534f702c31fdcaa43aa1f416c2784f`; watched core was `b42835481c6944e239129a78bbef0eca240a8473`.
- DONE: Bind only the reproduced Sonnet rejection-flow result.
  Local Claude Code 2.1.161 resolved `SPACEDOCK_LIVE_MODEL=sonnet` to `claude-sonnet-4-6`; full-ensign-cycle, gate-guardrail, withdrawn-gate-recovery, and recorded-gate-lifecycle passed, while owned default-headless skipped.
- DONE: Classify the exact rejection failure against the existing owner.
  The validation ensign durably committed its REJECTED report/JSONL, then `SendMessage(to:"team-lead")` found no addressable agent; the 60-second no-progress guard stopped the run. Existing repair `zbcj98qfwtax61vxdzrf615e` owns the missing rework/revalidation continuation.
- DONE: Qualify the local artifact correctly.
  `/tmp/spacedock-cycle20-sonnet-local.fkxGfs/artifacts/claude-shared-scenarios/rejection-flow` is ephemeral local evidence; the binding justification is the reproduced defect plus archived owner zbc's durable prior Sonnet/Opus triggers.
- DONE: Add RED-first exact controls and preserve target discrimination.
  Before the source binding, the runner failed the empty Sonnet/rejection TODO and reconciliation emitted exact `MISSING-EVIDENCE` with target `claude-sonnet` and owner `zbc…`; removal, retag, and unverified-Opus suppression mutations remain causal.
- DONE: Preserve pure desired and observed truth.
  The desired matrix remains 16×4=64; exactly fourteen target-scoped bindings are Sonnet=4, Codex=5, Pi=5, Opus=0. Opus remains runnable, and TODO remains explicitly non-green evidence.
- DONE: Preserve the durable Captain completion contract.
  The rule audited at `692ae776f` remains intact: every desired cell must pass or carry an exact target-scoped TODO(owner) backed by a reproduced defect; a missing owner or runnable path still blocks completion.
- DONE: Detect and classify external current-state oracle drift.
  The first full suite failed only because PR #629 / state archive commit `3f85cb081` moved `codex-launch-multi-agent-v2.md` under `_archive/`; focused default and race gate checks reproduced the same missing active path.
- DONE: Apply only the Captain-authorized archive reconciliation.
  Audit showed the candidate already counted 23 archived paths because of the prior gate-agent reconciliation, correcting the initial 22→23 assumption. Only the manifest path changed to `_archive/codex-launch-multi-agent-v2.md`, and only the archive predicate/diagnostic changed from 23 to 24.
- DONE: Prove the corrected integration oracle.
  `TestV1PilotManifestReadsAndValidates` passes and names the archived path; no other gate or product behavior changed. Classification: Material current-state integration-oracle drift, externally triggered, disposition fixed.
- DONE: Preserve core plus SHA-only binder history.
  Core `702d2b13ba1f9ac794e09435f4520bd20c3a015c` is followed only by binder `3ee17e6f3f324c9e973d217fb1f38cf0dcdf235b`, changing the guide SHA. The stale guard named exactly both changed ensigncycle paths before the binder; current and stale controls then passed opposite expectations.
- DONE: Stay inside the authorized surface.
  Final cumulative surface is exactly 43 files/+2750/-714, adds no file, and contains no product/runtime, fixture, prompt, adapter, launcher, observer, metric, workflow, owner, retry, controller, or unrelated TODO change.
- DONE: Run required verification and publish the candidate.
  `gofmt -w ./cmd ./internal` is idempotent, `git diff --check` is clean, focused reconciliation/mutation/coherence/TODO/manifest checks pass, full exits 0 in 33.23s, and race exits 0 in 184.20s. Candidate and origin are clean and equal at `3ee17e6f3f324c9e973d217fb1f38cf0dcdf235b`.
- SKIPPED: Bind the separate Sonnet-5 probe or run other live targets.
  Claude Code 2.1.220 resolving alias `sonnet` to `claude-sonnet-5` in a no-tools probe is not journey evidence and caused no scope change; no Opus, hosted, model, substrate, or smoke run was performed during implementation.
- DONE: Preserve workflow authority.
  This report records Cycle 21 only; implementation made no entity transition.

### Summary

Cycle 21 binds only Sonnet rejection-flow to existing `zbc…`, reconciles the externally archived PR #629 pilot path with the corrected 24-path oracle, and preserves the 692ae completion rule. The exact 43/+2750/-714 two-commit candidate is fully green and pushed at `3ee17e6f3f324c9e973d217fb1f38cf0dcdf235b` without a workflow transition.

## Stage Report: validation (cycle 21, independent)

- DONE: Verify exact immutable provenance and surface.
  Candidate and upstream are clean and equal at `3ee17e6f3f324c9e973d217fb1f38cf0dcdf235b`; baseline is `0f9e096e7f534f702c31fdcaa43aa1f416c2784f`, watched core is `702d2b13ba1f9ac794e09435f4520bd20c3a015c`, and binder `3ee17e6f3` changes only the guide's reconciliation SHA. The merge-base surface is exactly 43 existing files/+2750/-714 with no added file.
- DONE: Independently classify the retained local Sonnet artifact.
  `/tmp/spacedock-cycle20-sonnet-local.fkxGfs/artifacts/claude-shared-scenarios/rejection-flow` records `claude-sonnet-4-6`: launch, dispatch, durable REJECTED report/JSONL mutation, and commit all succeeded before `SendMessage(to:"team-lead")` returned no addressable agent and the 60-second quiet guard stopped the run at 419.27s. This is continuation/addressability product conduct, not authentication, capacity, fixture, or oracle failure.
- DONE: Verify the existing cross-runtime owner and exact one-cell binding.
  Archived owner `zbcj98qfwtax61vxdzrf615e` explicitly owns rejection-flow rework/revalidation and retains prior Sonnet and Opus triggers. The core adds only Sonnet/rejection-flow to that owner; all thirteen prior bindings remain intact, Opus remains runnable, and TODO remains non-green evidence.
- DONE: Verify pure desired and observed truth.
  The desired matrix remains 16 journeys x four targets = 64 cells. Reconciliation reports exactly fourteen target-scoped gaps: Sonnet=4, Codex=5, Pi=5, Opus=0.
- DONE: Prove causal discrimination.
  Reconciliation mutation controls pass for Sonnet removal, owner retag, wrong target, duplicate source binding, global TODO, missing/unowned TODO, and unverified Opus suppression, while the restored candidate is green. SHA controls reject stale watched state and accept the exact core.
- DONE: Verify current state integration without relaxing the oracle.
  State has `_archive/codex-launch-multi-agent-v2.md` and no active counterpart. Baseline already counted 23 archives because `gate-agent-ergonomics.md` was archived; the candidate changes only this manifest path and 23→24, retains the strict 31-path count, and the state-aware pilot test passes all 31 entries. No override, restore, schema relaxation, or unrelated state path is present.
- DONE: Audit scope and separate probe evidence.
  The six-file cycle diff changes only operator evidence, reconciliation/TODO tests and source binding, plus the pilot manifest/count. It changes no product, prompt, fixture, adapter, launcher, observer, metric, owner, retry, controller, or workflow mechanism. The Claude Code 2.1.220 Sonnet-5/high no-tools probe is correctly excluded from journey evidence.
- DONE: Run formatting, focused, and repository-wide verification at the exact candidate.
  `gofmt -l ./cmd ./internal` emitted nothing and `git diff --check` passed without changing bytes. Focused runner/reconciliation/mutation/SHA/pilot tests passed; full exited 0 in 242.95s and race exited 0 in 247.05s.
- DONE: Render the independent verdict.
  PASSED. Material: none. Deferred-risk: none. Polish: none. Validation changed no candidate bytes and made no workflow state transition.

### Summary

Cycle 21 PASSED independently at exact pushed candidate `3ee17e6f3f324c9e973d217fb1f38cf0dcdf235b`. The retained Sonnet failure is accurately bound to the existing cross-runtime `zbc…` owner, while runnable Opus, the pure 16x4 matrix, all prior evidence, strict pilot-state validation, fixed surface, and scope boundaries remain intact.

## Stage Report: validation (cycle 22 — surface recarve)

- FAILED: Prove that source reconciliation derives evidence truth rather than restating it.
  Material finding 1: `internal/contractlint/live_registry_reconciliation_test.go` is 795 lines and copies all fourteen source gaps into `auditedMissingEvidence`, then tests the extracted source against that copy. The duplicated current-gap oracle makes the result tautological: source and expectation can drift together while the proof remains green.
- FAILED: Hold the approved operator surface near the original estimate and remove competing identities.
  Material finding 2: candidate `3ee17e6f3f324c9e973d217fb1f38cf0dcdf235b`, core `702d2b13ba1f9ac794e09435f4520bd20c3a015c`, is exactly +2750/-714, net +2036. That is +1626 net over the original +410 estimate. Three standalone journeys remain as unreachable `runRetired…` wrappers, while one synthetic runner map forwards every identity to the same trampoline and three runtime-specific 16-case switches redispatch by journey name.
- DONE: Record the Captain's replacement disposition.
  REJECTED. Replace the evidence-only candidate from `origin/main`. Implement sixteen real `TestLiveCommon…` Go functions with adjacent stable-ID and fixture-binding metadata; each test must call an executable `liveJourney` directly. Runtime selection belongs inside the shared helper or fixture parameter, never in a central journey map, runner registry, or journey-name switch.
- DONE: Define the replacement operator surface.
  Correct the exact registry entrypoints, delete the competing Pi auto-continue surface, and select the same real tests in every runtime lane with common prefix `^TestLiveCommon` plus `-failfast`. Do not add a map, switch, copied oracle, mutation lab, reconciliation SHA self-guard, or speculative Pi observer/metrics layer.
- DONE: Preserve earned evidence during replacement.
  Banked live artifacts and the fourteen target-scoped gap facts remain evidence inputs. Recarving representation and routing must not erase passing runs, infer a pass from TODO, change repair ownership, or broaden a target gap.
- DONE: Preserve authority and mutation boundaries.
  This report changes only the entity record. It does not transition status, modify candidate/code, run a model lane, or supersede any product repair owner.

### Summary

Cycle 22 REJECTED `3ee17e6f3` as an evidence-only implementation whose reconciliation is tautological and whose net surface is 1,626 lines over the original estimate. The Captain approved a replacement from `origin/main`: sixteen directly executable `TestLiveCommon…` tests, adjacent source metadata, one fail-fast prefix selector across runtimes, and source-derived reconciliation without parallel maps, switches, copied gap oracles, mutation laboratories, SHA self-guards, or speculative Pi observation/metrics machinery.

## Stage Report: implementation (cycle 23 — portable common live journey recarve)

- DONE: Rebuild from the exact approved main baseline.
  Baseline was `8b5af99baa5c37fe7c969a904819041688420e22`; the replacement is the single clean pushed commit `b1f49a395ccccf296b294617ed0f8b762c7c5e9d` on `spacedock-ensign/deliver-portable-common-live-journey-surface-recarve`.
- DONE: Replace competing journey identities with sixteen real common tests.
  Exactly sixteen top-level exported `TestLiveCommon…` functions carry adjacent stable journey/fixture annotations and call `liveJourney` directly. Runtime selection occurs only inside the shared helper; no journey registry, runner map, or journey-name switch remains.
- DONE: Put every supported runtime on the same executable surface.
  Claude, Codex, and Pi transports normalize their real command evidence behind one driver seam. The workflow contains exactly three `-run '^TestLiveCommon' -failfast` lanes, and the retired Pi auto-continue and runtime-specific shared-coverage surfaces are removed.
- DONE: Preserve earned gaps without copying an expected-gap ledger.
  Compact AST reconciliation derives desired journeys and targets from the registry and observed bindings from source, checks exact builder ownership and builder/assertion reachability, rejects malformed or duplicate target-scoped TODO metadata, checks orphan fixtures and workflow selectors, and emits stable sorted derived totals: Sonnet=4, Opus=0, Codex=5, Pi=5.
- DONE: Preserve real host evidence and operator rationale.
  Claude keeps anti-shutdown timing, stream-watcher liveness, shallow-boot stream grading, and merged-host ground-truth rationale. Codex and Pi command extraction use native transcript shapes; filing and zero-discovery assertions consume normalized commands without a speculative Pi observer or metrics layer.
- DONE: Meet the Captain's exact replacement surface.
  The final baseline-relative change is exactly 34 files, +1083/-1947. It removes obsolete tables, maps, wrappers, duplicate orchestration, and superseded meta-tests; `shared_live_runner_unit_test.go` is absent and `codex_wait_matrix_live_test.go` is unchanged.
- DONE: Run focused formatting and contract verification.
  `gofmt -w ./cmd ./internal` and `git diff --check` are clean. `TestRuntimeLiveRegistryReconciliation` passes with the source-derived sorted totals, and `go test -tags=live ./internal/ensigncycle -run '^$'` compiles the live surface.
- DONE: Run the required repository-wide suites against uncontaminated state.
  In clean clone `/tmp/spacedock-live-final.wdbRyb/repo` at the exact candidate, `go test -p 1 ./...` exited 0 and `go test -p 1 ./... -race` exited 0. The shared development checkout's independently archived pilot files were therefore excluded from repository discovery rather than altering product or state truth.
- SKIPPED: Run hosted/model live journeys.
  Cycle 23 changes representation and routing only; it preserves the banked live evidence and exact owner bindings without claiming TODO as passing evidence or spending new model/substrate capacity.
- DONE: Preserve authority boundaries and publish the immutable candidate.
  No product behavior, fixture semantics, repair owner, status transition, or unrelated state changed. HTTPS push lacked workflow scope, so the already-authenticated SSH transport published the identical candidate; the remote branch resolves exactly to `b1f49a395ccccf296b294617ed0f8b762c7c5e9d`.

### Summary

Cycle 23 replaces the rejected evidence-only surface with sixteen directly executable, runtime-portable common journeys and a compact source-derived reconciler. The exact 34/+1083/-1947 candidate is focused-green, full/race-green in a clean clone, and pushed at `b1f49a395ccccf296b294617ed0f8b762c7c5e9d` without a workflow transition.

## Stage Report: validation (cycle 23 — independent surface recarve)

- DONE: Verify immutable provenance, remote equality, base, and hard surface bounds.
  The isolated clone's candidate and upstream are exactly `b1f49a395ccccf296b294617ed0f8b762c7c5e9d`; merge base is exactly `8b5af99baa5c37fe7c969a904819041688420e22`. The change is 34 files/+1083/-1947, `shared_live_runner_unit_test.go` is absent, and the rejected `3ee17e6f` worktree/branch was not touched.
- DONE: Verify real declarations, adjacent metadata, selectors, and source-derived gap truth.
  Source has exactly sixteen exported top-level `TestLiveCommon…` declarations and sixteen adjacent stable journey/fixture annotations. Registry entrypoints match. All three workflow lanes use exact `-run '^TestLiveCommon' -failfast`. Reconciliation derives and logs Sonnet=4, Opus=0, Codex=5, Pi=5; TODO remains target-scoped and non-green, with ownership unchanged.
- DONE: Prove representative drift controls in disposable detached checkouts.
  Removing one declaration failed at 16/15; adding an orphan fixture failed with its exact ID; changing only the Pi selector failed the exact three-lane check. None was committed or applied to the candidate.
- FAILED: Prove that each declared builder/assertion binding is executable rather than nominal.
  Material finding 1: `liveJourney` receives builder and assertion function values but only reflects their names; it executes only `exercise`. The reconciler's recursive `functionReaches` treats any identifier reference as reachability. In a detached control, replacing the required `someCommitNamesOnly(...)` durable assertion call with `_ = someCommitNamesOnly` still made `TestRuntimeLiveRegistryReconciliation` pass with all expected gap totals. The declared assertion can therefore become non-executable while the authoritative proof remains green, contrary to the Captain's direct executable binding boundary.
- FAILED: Confine runtime selection and behavior to the one approved selector/helper seam.
  Material finding 2: shared journey exercises retain four `runner.(claudeLiveRunner)` branches across recorded-gate, gate-stop, and withdrawn-recovery paths. The recorded-gate branch also copies and grades a Claude-only plugin body. Wrapping the same driver changes those exercise semantics, so runtime-specific selection remains outside `liveDriverForRuntime`; names and interface conformance do not make the exercise runtime-portable.
- DONE: Classify the retained scenario fields and scope.
  `sharedRuntimeScenario.oldPythonTest` and `.intent` survive only as substrate-test initialization baggage; common declarations populate only `name`. This is polish, not a central common table or independent Material defect. No central journey list/map/runner registry, per-runtime journey-name switch, copied gap oracle, committed mutation lab, SHA guard, speculative Pi observer/metrics layer, competing Pi auto-continue/shared wrapper, product behavior change, fixture-semantic deletion, or PR #629 value deletion was found.
- DONE: Run required verification in the isolated clone.
  `gofmt -l ./cmd ./internal` emitted nothing; `git diff --check`, focused reconciliation, live-tag compile, focused ensign checks, and release workflow suites passed. Exact `go test ./...` exited 0 in 156.11s; exact `go test ./... -race` exited 0 in 161.04s. No live/model lane ran.
- DONE: Render the independent verdict and preserve authority.
  REJECTED. Fix the two representation/mechanism findings without changing banked live evidence, target-gap facts, repair ownership, product behavior, or the approved bounds. Validation changed no candidate/code and made no workflow status transition.

### Summary

Cycle 23 is scope-, budget-, selector-, registry-, focused-, full-, and race-green at exact pushed candidate `b1f49a395ccccf296b294617ed0f8b762c7c5e9d`, but is independently REJECTED on two Material Captain-boundary failures: nominal builder/assertion bindings can pass without execution, and shared exercises still branch on the concrete Claude driver outside the sole runtime selector helper.

## Stage Report: implementation (cycle 24 — executable-binding correction)

- DONE: Correct only the two independently rejected Cycle-23 findings on the immutable recarve.
  The exact starting candidate was `b1f49a395ccccf296b294617ed0f8b762c7c5e9d`; the correction is the clean pushed successor `5db8cec3e17b128fa07d3fca4909022d80a98324` on the same recarve branch. The rejected `3ee17e6f` branch and worktree remain untouched.
- DONE: Make every declared builder and assertion an executed value.
  `liveJourney` now passes typed builder/assertion wrappers into the bound exercise, counts their real runtime invocations with `reflect.MakeFunc`, and fails when either count is zero. All sixteen exercises directly call the received builder and assertion parameters.
- DONE: Replace nominal identifier reachability with direct executable-binding reconciliation.
  The reconciler retains its original registry parsing, ensigncycle+livescenario scan, fixture ownership, TODO-owner validation, diagnostics, and derived-gap logic; only `functionReaches` is replaced by an AST check requiring direct `ast.CallExpr` uses of exercise parameters four and five. The disposable dead-reference control that stayed green in Cycle 23 is therefore rejected.
- DONE: Remove concrete Claude selection from shared exercises.
  All four `runner.(claudeLiveRunner)` branches are gone. Recorded-gate preparation and its Claude-only copied-skill grade live behind the `liveDriver.prepareRecordedGate` capability; Codex and Pi provide host-neutral no-op preparation. `liveDriverForRuntime` remains the sole Claude/Codex/Pi selector, with no journey-name dispatch.
- DONE: Preserve fixture and transcript behavior while staying inside the approved surface.
  The reused durable-field reader has the same first-unindented-field and path-rejection behavior as the deleted local loop. The derived split-root Pi README is byte-equivalent to the prior fixture, and generalized Pi transcript parsing retains exact read and nonempty bash/shell-command semantics, including rejection of empty tool names.
- DONE: Preserve banked evidence, ownership, and architecture boundaries.
  Exactly sixteen exported top-level `TestLiveCommon…` declarations retain adjacent stable annotations and direct seven-argument `liveJourney` calls. The workflow retains exactly three `-run '^TestLiveCommon' -failfast` selectors. Reconciliation derives Sonnet=4, Opus=0, Codex=5, and Pi=5; no product repair, behavior, fixture semantics, TODO owner, PR-629 manifest value, central table/map/list, runtime journey switch, mutation lab, SHA guard, speculative Pi observer/metrics, or competing wrapper was added.
- DONE: Meet the exact Captain-approved envelope.
  The cumulative base-relative surface at `5db8cec3e` is exactly 34 files, +1083/-2018; `shared_live_runner_unit_test.go` remains absent. `gofmt -w ./cmd ./internal` is idempotent and `git diff --check` is clean.
- DONE: Run focused and repository-wide verification.
  Reconciliation passes with the four exact sorted gap totals; live-tag package compilation passes; focused `internal/ensigncycle` and `internal/release` pass (`ensigncycle` 296.236s). In clean clone `/tmp/spacedock-cycle24-verify.i0zrRM/repo`, exact HEAD, blank status, and absent development state checkout were logged before each command: `go test ./...` passed in 377.75s and `go test ./... -race` passed in 188.09s.
- DONE: Classify discarded operator-precondition attempts without attributing them to the candidate.
  One command checked out the clean clone but omitted `cd`, so it ran in the contaminated source checkout and hit the known external state-root manifests; a second just-started run was interrupted because its preconditions were asserted silently rather than logged. Neither is candidate evidence; the two separately logged clean-clone suites above are authoritative.
- SKIPPED: Run hosted or local model/live lanes.
  Implementation spent no model or substrate capacity. Those lanes remain after independent validation of this representation correction.
- DONE: Publish the immutable correction without changing workflow state.
  SSH push succeeded and the remote recarve branch resolves exactly to `5db8cec3e17b128fa07d3fca4909022d80a98324`. Implementation made no status transition.

### Summary

Cycle 24 closes both independent Material findings: all sixteen declarations now execute their exact typed builder/assertion bindings, and shared exercises contain no concrete Claude-driver selection. The exact 34/+1083/-2018 candidate is focused-, full-, and race-green in uncontaminated state and pushed at `5db8cec3e17b128fa07d3fca4909022d80a98324`, with live/model execution held for fresh independent validation.

## Stage Report: validation (cycle 24 — independent correction rerun)

- DONE: Verify immutable provenance, clean isolation, and the exact approved envelope.
  Remote, local branch, HEAD, and upstream are exactly `5db8cec3e17b128fa07d3fca4909022d80a98324`; merge base is exactly `8b5af99baa5c37fe7c969a904819041688420e22`. The clean isolated clone `/tmp/spacedock-cycle24-independent.LOm7ei` had an empty status and no `docs/dev/.spacedock-state`. Surface is exactly 34 files/+1083/-2018, `shared_live_runner_unit_test.go` is absent, and rejected evidence candidate `3ee17e6f` was untouched.
- DONE: Revalidate declarations, metadata, executable calls, selectors, and gap truth.
  Source contains exactly sixteen exported top-level `TestLiveCommon…` declarations with sixteen adjacent stable journey/fixture annotations and direct seven-argument `liveJourney` calls. The registry names those exact entrypoints. All three workflow lanes use `-run '^TestLiveCommon' -failfast`. Source reconciliation derives sorted Sonnet=4, Opus=0, Codex=5, Pi=5; TODO remains target-scoped and non-green with byte-identical owner bindings from Cycle 23.
- DONE: Reproduce and close the Cycle-23 nominal-binding finding adversarially.
  In a detached disposable checkout, replacing the required full-cycle assertion call with `_ = assert` while retaining its nominal binding failed reconciliation: `runFullEnsignCycleJourney` no longer directly called its bound builder and assertion. Restored immutable source passed. The reconciler now requires direct `ast.CallExpr` use of exercise parameters four and five.
- DONE: Prove runtime execution enforcement independently of AST spelling.
  A disposable live-tag control wrapped an exact function with `countedLiveFunction`: pre-call wrapper/underlying counts were 0/0, one wrapped invocation produced 1/1, and a separate non-invocation control remained observably zero. This proves `liveJourney` can reject a bound builder or assertion that is not executed at runtime.
- DONE: Reproduce and close the Cycle-23 concrete-driver portability finding.
  Static inspection finds zero `runner.(claudeLiveRunner)` branches or other concrete Claude assertions. Recorded-gate copied-skill preparation/grade is behind `liveDriver.prepareRecordedGate`; Claude implements the copy/grade and Codex/Pi return the same driver with `noLiveGrade`. A disposable interface-wrapper control invoked that capability through `liveDriver`, incremented the underlying preparation count exactly once, and returned a usable driver/grade without journey-name redispatch. `liveDriverForRuntime` remains the only Claude/Codex/Pi selector.
- DONE: Verify fixture and parser equivalence rather than accepting refactor names.
  A disposable control proved the derived split-root Pi README byte-for-byte equal to the prior `b1f49a3` `piAutoContinueReadme`. The `durableField` resolver matched the prior worktree parser for missing, valid, dot, parent, absolute, and first-duplicate-wins inputs. Generalized transcript parsing preserved read path plus nonempty bash/shell command extraction and rejected empty tool names and empty commands.
- DONE: Revalidate architecture, scope, and retained value.
  `sharedRuntimeScenario` now contains only `name`. No central journey table/list/map/runner registry, per-runtime journey-name switch, copied gap oracle, committed mutation lab, reconciliation SHA guard, speculative Pi observer/metrics, competing Pi auto-continue/shared wrapper, product behavior change, fixture-semantic drift, TODO-owner drift, or PR #629 manifest-value deletion was found.
- DONE: Run formatting, focused, and repository-wide verification.
  `gofmt -l ./cmd ./internal` emitted nothing and `git diff --check` passed. Focused reconciliation, live-tag compile, ensign durable/rejection checks, release workflow checks, and all disposable controls passed. Exact `go test ./...` exited 0 in 242.82s; exact `go test ./... -race` exited 0 in 247.64s. No live/model lane ran.
- DONE: Render the independent verdict and preserve authority.
  PASSED. Material: none. Deferred-risk: none. Polish: none. Validation changed no candidate/code and made no workflow status transition.

### Summary

Cycle 24 PASSED independently at exact pushed candidate `5db8cec3e17b128fa07d3fca4909022d80a98324`. Both prior Material findings are closed causally: nominal dead references fail static reconciliation, missing runtime invocation is count-detectable, and shared exercises are portable through capabilities with runtime selection confined to one helper. Scope, evidence ownership, fixture behavior, exact selectors, hard bounds, full, and race remain green.

## Exact-Head Local Live Evidence: cycle 24 (checkpoint 1)

- Candidate and remote branch are clean, equal, and upstream-tracked at `5db8cec3e17b128fa07d3fca4909022d80a98324`; no PR or hosted CI exists.
- Mandatory Opus attempt 1 used Claude Code 2.1.220, `SPACEDOCK_LIVE_MODEL=claude-opus-4-8`, and exact `-run '^TestLiveCommon' -failfast`. `full-ensign-cycle`, `gate-guardrail`, and `default-headless-gate-stop` passed. `withdrawn-gate-recovery` then exhausted API retries with DNS `ENOTFOUND` before usable FO work. The lane stopped after 1119.88s; retained artifacts: `/tmp/spacedock-cycle24-opus-local.vnbAhx`.
- Immediate DNS/HTTPS checks resolved `api.anthropic.com` and connected normally, so attempt 2 reran the complete exact selector. `full-ensign-cycle`, `gate-guardrail`, `default-headless-gate-stop`, and `withdrawn-gate-recovery` passed. `recorded-gate-lifecycle` then ended after 479.17s with `API Error: Connection closed mid-response`; the complete selector stopped after 2001.86s. Retained artifacts: `/tmp/spacedock-cycle24-opus-rerun.myD6Xq`.
- The second failure's diagnostic also preserved one FO command that broadly searched `.spacedock-state`; because transport terminated before the journey's post-run copied-skill and durable-state grades, this is an ancillary observation, not a passing result or a final product verdict. A clean targeted rerun must disposition it before Opus completion.
- Neither stopped selector is exact-head lane completion. Pi `openai-codex/gpt-5.4`, Sonnet 5/high, and Codex `gpt-5.6-luna`/max have not started. All fourteen source TODO bindings remain unchanged and non-green.

## Exact-Head Local Live Evidence: cycle 24 (checkpoint 2)

- A direct subscription probe returned `OK` from `claude-opus-4-8` without API error, so `TestLiveCommonRecordedGateLifecycle` was rerun alone at the same exact candidate.
- The targeted run reached the gate presentation and called `gate record --room`. The command failed because `provider/result.json` did not exist; the FO then made no stream progress for 60 seconds and the harness killed it. This is a real nonpassing product/runtime disposition, not DNS, authentication, capacity, or an oracle-only failure. Retained artifacts: `/tmp/spacedock-cycle24-opus-recorded-gate.Af7c8Z`.
- Existing active owner `prove-or-cut-provider-backed-gate-closure` (`a732sahay8wzgqrd2yr0xxr7`) owns this exact boundary. Its approved implementation direction is the clean-cut branch: a provider run stopping before recording is a cut signal, not permission to add retry/fallback machinery, and the task owns removal of `gate record --room`/provider-only skill wiring while preserving chat closure.
- Disposition: do not modify ys product behavior, add an Opus TODO, or change any of the fourteen exact bindings. Opus remains non-green on the exact head pending the existing owner's clean-cut delivery (or an explicit captain decision). Continue the other three local exact-head lanes before any hosted CI.

## Exact-Head Local Live Evidence: cycle 24 (checkpoint 3)

- The Pi lane used Pi 0.82.1 with subscription model `openai-codex/gpt-5.4`, `pi-subagents` 0.35.1, `pi-intercom` 0.6.0, `SPACEDOCK_PI_LIVE_REQUIRED=1`, and exact `-run '^TestLiveCommon' -failfast` at candidate `5db8cec3e17b128fa07d3fca4909022d80a98324`.
- Ten runnable journeys passed before the final runnable cell: `full-ensign-cycle`, `withdrawn-gate-recovery`, `recorded-gate-lifecycle`, `feedback-3-cycle-escalation`, `merge-hook-guardrail`, `filing`, `shallow-boot`, `zero-discovery`, both single-root and split-root `auto-continue`, and `self-evidence-merge-triage`. The five exact Pi TODO cells remained skipped and non-green: `gate-guardrail`, `default-headless-gate-stop`, `rejection-flow`, `smallest-sufficient-mechanism`, and `keep-moving-posture`.
- `ac-value-reanchor` failed after 115.09s; the complete selector stopped after 1665.019s. The Pi process itself completed without transport, authentication, or timeout error. The FO correctly found AC-2 regressed (`10,200 > 8,000`), prepared and validated the gate, and presented REJECT, but did not record `revise` or apply the configured route despite having the scenario connection. Durable gate state therefore remained canonical `open`, and the external oracle correctly failed it.
- Retained artifacts: `/tmp/spacedock-cycle24-pi-local.eBEv6Y`. This is a real model/product-conduct nonpass, not a TODO pass or oracle-only failure. The prior archived `ac2-reanchor-live-scenario-repair` owned falsifiability of the durable oracle, not this fresh runtime failure to finish the authorized gate lifecycle; current-owner mapping and one targeted exact-head rerun remain pending.
- Disposition: do not weaken the durable oracle, modify ys product behavior, or change any of the fourteen exact bindings. Pi remains non-green until the targeted rerun and owner mapping produce a truthful disposition. Sonnet 5/high and Codex `gpt-5.6-luna`/max have not started; hosted CI remains prohibited.

## Exact-Head Local Live Evidence: cycle 24 (checkpoint 4 — Pi correction)

- The targeted Pi rerun retained at `/tmp/spacedock-cycle24-pi-reanchor.5L698j` disproved checkpoint 3's provisional classification of the shared error as a second model-conduct failure. Pi prepared and committed the bound review, recorded `decision: revise`, and routed the entity to `status: rework`; its final report accurately described that durable outcome.
- The common test still failed because `assertACReanchorScenario` constructed a fresh `AuthorACReanchorScenario()` after setup. That fresh assertion closure had an empty `entityPath`, so the wrapper—not the product state—emitted `validation gate is not canonical: open : no such file or directory`. The exact error therefore could not distinguish the first full-selector run's real present-only stop from the targeted rerun's successful gate route.
- The first full-selector trace remains a genuine instance of the behavior owned by active task `pi-delegated-gate-continuation-reliability` (`9w59t6m1qc46hccd54p04z2j`): correct bound presentation followed by an early stop before record/application under delegated conn. The targeted rerun proves the capability is present but not yet reliable; it does not justify a new source TODO or a weaker durable oracle.
- A focused red test now reproduces the candidate wiring bug before correction: the common assertion must receive the exact authored scenario whose setup populated the durable entity path. No source TODO binding or product behavior is being changed. Pi remains non-green until the test-only correction is pushed and the exact-head lane is rerun.

## Exact-Head Local Live Evidence: cycle 24 (checkpoint 5 — oracle wiring fixed)

- The focused regression was added before implementation. It failed to compile because `assertACReanchorScenario` accepted only state values and could not receive the built scenario. The correction passes the exact authored `livescenario.Scenario` into the assertion binding and delegates to that instance's setup-populated closure; the regression then passed and rejects an unchanged entity for the durable-state reason rather than an empty-path error.
- The correction changes only `internal/ensigncycle/shared_promoted_live_test.go` (+24/-4). It does not change product behavior, fixture semantics, target gaps, or any of the fourteen TODO bindings. Focused live-tag regression, `TestACReanchor`, source reconciliation, and live-tag compilation all pass after `gofmt -w ./cmd ./internal`.
- The new clean pushed and remote-equal candidate is `5732cc822493d12921fff01085a7a6d1df06b60b`. Clean clone `/tmp/spacedock-cycle24-oracle-fix-verify.aSGx1w/repo` has no development state checkout. Exact `go test ./...` passed in 310.57s and exact `go test ./... -race` passed in 369.57s.
- One preceding full-suite attempt in the shared worktree is not candidate evidence: all code-bearing packages passed, but `internal/gates` observed two concurrently archived pilot manifests through the external development state checkout. The clean-clone suites above exclude that state and are authoritative.
- The successful targeted Pi trace at the prior head remains retained diagnostic evidence, not exact-head completion. The corrected `5732cc822` journey must rerun before Pi receives a disposition. Opus remains bound to `a732sahay8wzgqrd2yr0xxr7`; Sonnet 5/high and Codex `gpt-5.6-luna`/max remain unstarted; hosted CI remains prohibited.

## Exact-Head Local Live Evidence: cycle 24 (checkpoint 6 — final validation candidate)

- Pi's focused `ac-value-reanchor` rerun passed at exact `5732cc822` in 171.76s, proving the authored-scenario assertion correction against canonical durable gate state. The subsequent full Pi selector passed through `zero-discovery`, including recorded-gate chat closure, then failed `auto-continue/single-root` because the fixture's README was not commissioned. Pi correctly stopped on `status --boot --identify --json`; retained artifacts: `/tmp/spacedock-cycle24-pi-full-fixed.A2RKgR`.
- Source inspection found the recarve had removed `commissioned-by: spacedock@1` from only `autoContinueReadme()` while the split-root derivative re-added it. A behavioral `status.DiscoverWorkflowDir` regression was added first and failed on the malformed fixture; restoring the one semantic line made it and existing auto-continue assertions pass. This is fixture-truth restoration, not product behavior or an owner change.
- Sonnet's superseded-head run used canonical `claude-sonnet-5` through an enforced `--effort high` shim. It passed full-cycle, gate-guardrail, withdrawn recovery, recorded gate, feedback escalation, merge guard, filing, shallow boot, zero discovery, and the single-root auto-continue variant. It was intentionally interrupted after 1248.14s because the candidate had become invalid; retained artifacts: `/tmp/spacedock-cycle24-sonnet-high.CcDCE9`. Those passes are diagnostic only, not final-head lane completion.
- Mandatory Opus `recorded-gate-lifecycle` remains a reproduced product nonpass owned by active task `prove-or-cut-provider-backed-gate-closure` (`a732sahay8wzgqrd2yr0xxr7`). Its PR #632 is not landed or integrated. Per the Commander completion rule, only target `claude-opus` on `recorded-gate-lifecycle` is now bound to `TODO(a732sahay8wzgqrd2yr0xxr7)`; the prior fourteen bindings are byte-unchanged and remain non-green. Reconciliation derives Sonnet=4, Opus=1, Codex=5, Pi=5, and an exact Opus selector proves the new cell skips rather than passes.
- The clean pushed and remote-equal final validation candidate is `546bc5220a3ddac07e14f28d820de538575abbdc`, cumulatively 35 files/+1116/-2018 (net -902) from exact base `8b5af99baa5c37fe7c969a904819041688420e22`. It adds no product behavior. Focused discoverability, auto-continue assertions, reconciliation, live-tag compilation, formatting, and diff checks pass.
- All model evidence before `546bc5220` is retained diagnostic evidence only. One independent validation of this combined final head must pass before restarting exact-head Pi, Sonnet 5/high, and Codex Luna/max. Hosted CI and this task's PR remain prohibited until those local dispositions exist.
