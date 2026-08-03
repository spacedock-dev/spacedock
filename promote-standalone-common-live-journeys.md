---
title: Promote standalone common behaviors into the shared live journey suite
status: ideation
source: "Desired live-test registry inventory, 2026-08-03"
started: 2026-08-03T11:21:39Z
completed:
verdict:
score: 0.9
worktree:
issue:
sprint: live-test-truth
group: common-runner
sprint-readiness: ready
id: r4qk46605sjcphj44cvkcsk4
gates:
    version: 1
    records:
        - id: gate:r4qk46605sjcphj44cvkcsk4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:r4qk46605sjcphj44cvkcsk4-backlog-1
              briefing:
                id: briefing:r4qk46605sjcphj44cvkcsk4:backlog:attempt-1:revision-1
                digest: sha256:88c06435c9d215a9101c54108777e7553ffbfd61eed6a8cc041d02829236d274
                request-digest: sha256:8f3ec87bf5943cdeed8428bfcc78f626e9f0b108612fe2037f11f506de0f30b6
                room-ref: ./promote-standalone-common-live-journeys/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:r4qk46605sjcphj44cvkcsk4:backlog:1
                briefing: briefing:r4qk46605sjcphj44cvkcsk4:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T11:20:28.996942Z"
                decision: approve
                reason: Captain directed the First Officer to continue the next risk-first ideation wave.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:r4qk46605sjcphj44cvkcsk4:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:r4qk46605sjcphj44cvkcsk4-ideation-1
              briefing:
                id: briefing:r4qk46605sjcphj44cvkcsk4:ideation:attempt-1:revision-1
                digest: sha256:c54e46c73cff782789ebca53b54860837d826ae3bd89e1546ef26e6e58cb9b37
                request-digest: sha256:6be3137e2d5255c6a21f319f0fdd8970250216f7a35e43d2716bc48f86c7da5f
                room-ref: ./promote-standalone-common-live-journeys/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T12:16:46.846839Z"
                reason: Captain recarved the sprint into outcome-shaped delivery units. Preserve this report as design input; do not present or consume this component-shaped attempt.
---

## Problem

The desired registry contains 16 common journeys. The current shared scenario table contains 10 journeys.

Six desired journeys remain outside the shared suite:

- `full-ensign-cycle`
- `default-headless-gate-stop`
- `withdrawn-gate-recovery`
- `zero-discovery`
- `auto-continue-after-implementation`
- `ac-value-reanchor`

The seed omitted `ac-value-reanchor`. This plan includes it because the registry lists it as common.

The six journeys use standalone tests or Claude-specific runner paths. `auto-continue-after-implementation` also has a separate Pi test.

These paths hide missing host implementations. They also let a standalone test use a different fixture or grader from the canonical shared identity.

The current default-headless test has the riskiest overlap. It calls the Claude gate-guardrail runner and changes setup through a scenario-name condition.

The final source must expose one canonical identity for each journey. Runtime launch and output parsing must remain inside adapters.

## Inputs and dependencies

This plan consumes the completed `bind-live-registry-to-source` design. New records and builders use its adjacent source annotations and stable IDs.

This plan also consumes the completed `converge-shared-live-suite-entrypoint` design. Implementation starts after that task creates `TestLiveSharedScenarios` and its adapter boundary.

The `ac2-reanchor-live-scenario-repair` task owns the falsifiability repair for `internal/livescenario/ac2_reanchor.go`. This task promotes that repaired scenario without weakening it.

Implementation order is:

1. Land the canonical Claude and Codex shared entry point.
2. Land the AC re-anchor falsifiability repair.
3. Promote all six journeys in this task.
4. Add the Pi common adapter against the resulting 16-entry shared map.

The Pi task must consume all 16 entries. Its earlier 10-entry estimate is not the final shared-suite contract.

## Spike result

The spike exercised the closest assertion overlap in a detached throwaway worktree. It used the real gate fixture and `assertGateHeld`.

The spike added two small composite graders:

- The gate-guardrail grader required a `validation` start and a held `validation` result.
- The default-headless grader required an `implementation` start and a held `validation` result.

Both baselines passed. Each grader rejected the starting state of the other journey.

```text
=== RUN   TestGateJourneyOverlapSpike
--- PASS: TestGateJourneyOverlapSpike (0.00s)
PASS
```

This result proves one shared gate-stop primitive is sufficient. Separate start-state assertions preserve the distinct journey outcomes.

The spike worktree was removed. It left no product change.

## Target behavior

`sharedRuntimeScenarios()` contains all 16 registry IDs. Each added record carries its adjacent `live-journey` annotation.

`TestLiveSharedScenarios/<journey-id>` is the only authoritative live identity for each promoted journey. Old standalone top-level tests do not remain as wrappers.

Each shared runner selects a semantic fixture by stable fixture ID. Concrete function names remain source details.

Fixture variants run inside their journey subtest. They do not create a second journey identity.

Claude and Codex run the six new journeys through their existing adapters. The later Pi adapter receives the same runner map.

## Reconciliation map

| Registry journey | Current path | Final fixture contract | Final host-neutral grade |
|---|---|---|---|
| `full-ensign-cycle` | `TestLiveEnsignCycle` | `realistic-lifecycle` starts one discoverable backlog member in a three-stage workflow. | The entity reaches `done`, remains locatable, has a valid stage report, and has a path-scoped commit. |
| `default-headless-gate-stop` | First case in `TestLiveDefaultHeadlessStopsAtGate` | `recorded-gate/pre-gate` starts at `implementation` before a human validation gate. | The member advances to `validation`, binds one review, presents it, and performs no authority action after prepare. |
| `withdrawn-gate-recovery` | Second case in `TestLiveDefaultHeadlessStopsAtGate` | `recorded-gate/withdrawn` starts with one immutable withdrawn attempt and no open successor. | One successor attempt opens and commits. The first room stays byte-identical. No decision or dispatch occurs. |
| `zero-discovery` | `TestLiveZeroDiscoverReportsAndStops` | `boot/no-workflow` is a Git repository with no commissioned workflow marker. | Discovery returns zero. The run reports that result, creates no team, performs no broad search, and changes no state. |
| `auto-continue-after-implementation` | Claude and Pi standalone tests | Both registry variants start with a completed implementation report and a fresh validation stage. | Both variants use `assertAutoContinue`. A stop after implementation fails both variants. |
| `ac-value-reanchor` | `TestLiveReanchorGateRejectsMeansOnlyRegressed` | `ac-reanchor/means-pass-value-regressed` has a passing mechanism criterion and a regressed value criterion. | The repaired scenario produces divergent durable outcomes for correct and incorrect gate decisions. |

## Fixture and variant contracts

### Gate fixtures

Keep `recorded-gate/held` unchanged for `gate-guardrail`. Extract `recorded-gate/pre-gate` from the current inline mutation.

The pre-gate builder writes the stage graph and entity directly. It does not select behavior through a journey-name condition.

Extract `recorded-gate/withdrawn` from the current runner setup. Its builder returns these saved facts:

- The entity path and state root.
- The first attempt ID and retained room path.
- The original `gate-briefing.json` and `request.json` bytes.
- The expected next attempt number.

One shared `assertGateStopBoundary` primitive grades the open retained attempt. Journey-specific graders add different state deltas.

`assertGateGuardrailJourney` requires `validation` before and after the run. `assertDefaultHeadlessGateStop` requires `implementation` before and `validation` after.

The default-headless command grade requires pre-gate work before prepare. It forbids decision, consume, dispatch, or archive after prepare.

The withdrawn grader requires one new attempt. It also requires unchanged bytes for the withdrawn room and no authority action.

### Full-cycle fixture

Extract a named `realistic-lifecycle` builder from `startRealisticLifecycleDrive`. The builder owns setup only.

Move process launch and liveness into the adapter. Move terminal state checks into one host-neutral full-cycle grader.

The grader receives the final entity, its resolved path, and normalized Git facts. It does not read Claude stream events.

### Zero-discovery fixture

Extract a named `boot/no-workflow` builder from the inline `.gitkeep` setup. Keep the default-tag undiscoverable fixture test.

The adapter returns normalized boot evidence. This evidence contains discovery count, team creation count, broad-search commands, exit state, and observed report text.

The common grader requires zero discovery, zero team creation, zero broad-search commands, a clean exit, and no durable mutation.

### Auto-continue variants

Keep `auto-continue/single-root` in `auto_continue_fixtures_test.go`. Move the split-root builder from the Pi live file into the same default-tag fixture file.

Both builders return one common fixture result. The result contains workflow root, state root, entity path, and archive candidates.

The journey runner executes both variants inside `auto-continue-after-implementation`. It passes each resolved entity body to the unchanged `assertAutoContinue` grader.

The runner does not create nested journey IDs. Artifact and failure labels include the stable fixture ID.

This task removes both runtime-specific auto-continue wrappers after it extracts the fixture contracts. The later Pi task adds only the common adapter binding.

### AC re-anchor fixture

Use the repaired `AuthorACReanchorScenario` as the common scenario. Change its scenario name to the registry ID during its owning repair task.

The shared runner invokes that scenario through the selected adapter. This task does not replace its divergent durable-state grader with transcript wording.

## Proposed approach

1. Add the six annotated records to `sharedRuntimeScenarios()` in registry order.
2. Add one shared runner binding for each new record.
3. Extract the inline setup into builders with the registry fixture annotations.
4. Extract durable assertions from the standalone tests.
5. Add the gate journey composites proven by the spike.
6. Add both auto-continue variants to one journey binding.
7. Delete obsolete standalone top-level tests and scenario-name setup conditions.
8. Remove old standalone names from workflow selectors and guards.
9. Run registry reconciliation and record the new SHA as the bind contract requires.
10. Run the focused, live, full, race, and format checks.

The common runner map serves the value criterion. Compatibility wrappers are insufficient because they preserve competing authoritative identities.

The fixture-variant list serves the auto-continue value criterion. One topology is insufficient because the original failure occurred in split-root state.

The shared gate-stop primitive serves the distinct gate criteria. Two complete gate graders duplicate the same authority-boundary logic.

The normalized boot evidence serves the zero-discovery value criterion. A universal transcript event model adds unrelated parser scope.

## Acceptance criteria

**AC-1 (VALUE)** Registry-to-runner parity is 16 of 16, and the count of standalone common-journey top-level tests decreases to zero.
Proven by: the registry reconciliation, the shared map parity test, and a live-tag test-list count.

**AC-2** `gate-guardrail` and `default-headless-gate-stop` preserve distinct start-state and transition outcomes through one shared gate-stop primitive.
Proven by: the cross-qualification negative test and focused live runs of both canonical selectors.

**AC-3** `auto-continue-after-implementation` runs both stable fixture variants under one journey identity and one `assertAutoContinue` grader.
Proven by: a variant coverage test and the stop-after-implementation negative control for each variant.

**AC-4** All six promoted journeys use host-neutral fixtures and graders. Authentication, launch, liveness, and output-dialect parsing remain adapter-owned.
Proven by: fake-adapter tests, adapter boundary tests, and focused Claude and Codex live runs.

**AC-5** `ac-value-reanchor` enters the common suite only after its separate repair proves that incorrect gate behavior fails.
Proven by: the divergent durable-state test from the repair and the canonical shared live selector.

**AC-6** The shared suite keeps exact journey IDs, stable fixture IDs, artifact ownership, metrics ownership, and existing runtime transport behavior.
Proven by: reconciliation, metrics tests, workflow guards, and adapter command tests.

## Acceptance checks and test plan

Add default-tag tests before the runner migration. These tests do not spend model credentials.

- Make the parity test compare all 16 registry IDs with records and runner bindings in both directions. Removing one binding must fail.
- Make the test-list guard require zero retired standalone names. Restoring any old wrapper must fail.
- Run the gate overlap test from the spike. Swapping either start fixture must fail the affected journey.
- Mutate the default-headless result to omit the stage transition. Its grader must fail while the held-gate baseline remains valid.
- Mutate the withdrawn room or add a decision command. The withdrawn grader must fail.
- Run existing full-cycle broken-output controls. A malformed report, nonterminal status, or missing path-scoped commit must fail.
- Feed a recorded broad search and a team-creation event into the zero-discovery evidence. Each mutation must fail.
- Run both auto-continue variants against a stopped implementation state. Both calls to `assertAutoContinue` must fail.
- Run the repaired AC re-anchor divergent control. Incorrect gate behavior must fail before promotion.
- Run fake Claude and Codex adapters against each fixture contract. A runtime branch in a common fixture or grader fails source review.

Run the focused offline checks:

```bash
go test -tags live ./internal/ensigncycle -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPromotedCommonJourneyEntrypoints' -count=1
go test ./internal/ensigncycle -run 'TestGateJourneyOverlap|TestAutoContinue|TestZeroDiscover|TestEnsignCycleGoesRed|TestACReanchor' -count=1
go test ./internal/release -run 'TestRuntimeLiveWorkflow|TestWorkflowsPreserveAndPublishJourneyCosts' -count=1
```

Run each promoted journey through both current adapters:

```bash
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 40m -run '^TestLiveSharedScenarios$/^(full-ensign-cycle|default-headless-gate-stop|withdrawn-gate-recovery|zero-discovery|auto-continue-after-implementation|ac-value-reanchor)$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveSharedScenarios$/^(full-ensign-cycle|default-headless-gate-stop|withdrawn-gate-recovery|zero-discovery|auto-continue-after-implementation|ac-value-reanchor)$' ./internal/ensigncycle -v
```

The auto-continue selector runs both fixture variants in one journey subtest. A missing variant fails before launch.

After the Pi adapter lands, run the same six canonical selectors with its required runtime configuration. No Pi-only journey name counts as evidence.

Run the repository checks:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

## Selectors and compatibility

The canonical selector for one promoted journey is:

```text
^TestLiveSharedScenarios$/^<journey-id>$
```

The Claude and Codex lanes select the full `TestLiveSharedScenarios` suite. They do not list standalone common tests beside it.

Delete these authoritative names without wrappers:

- `TestLiveEnsignCycle`
- `TestLiveDefaultHeadlessStopsAtGate`
- `TestLiveZeroDiscoverReportsAndStops`
- `TestLiveAutoContinueAfterImplementation`
- `TestLivePiAutoContinueAfterImplementation`
- `TestLiveReanchorGateRejectsMeansOnlyRegressed`

Go test names are not a supported external API. Historical artifacts can retain old names because they describe past runs.

## Expected surface

The baseline is 16 files. The estimate is about 690 insertions and 730 deletions.

The tolerance is 19 files and 870 insertions. The tolerance includes adjacent fixture or workflow-guard tests.

| File | Estimated change | Purpose |
|---|---:|---|
| `internal/ensigncycle/shared_scenarios_test.go` | 55 insertions | Add six annotated registry records. |
| `internal/ensigncycle/shared_live_runner_test.go` | 230 insertions | Add six common runners and the auto-continue variant loop. |
| `internal/ensigncycle/shared_live_runner_unit_test.go` | 70 insertions | Prove bindings, variants, and adapter-neutral inputs. |
| `internal/ensigncycle/shared_fixtures_test.go` | 85 insertions | Add named full-cycle, pre-gate, withdrawn, and zero-discovery builders. |
| `internal/ensigncycle/gate_assert_impl_test.go` | 55 insertions | Add shared gate-stop and distinct journey delta graders. |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | 75 insertions | Add cross-qualification and journey negative controls. |
| `internal/ensigncycle/auto_continue_fixtures_test.go` | 60 insertions | Add the split-root variant and common fixture result. |
| `internal/ensigncycle/live_test.go` | 210 deletions | Remove the full-cycle top-level test and Claude-only drive. |
| `internal/ensigncycle/live_gate_stop_test.go` | delete 51 lines | Remove the combined standalone gate test. |
| `internal/ensigncycle/zero_discover_live_test.go` | delete 115 lines | Remove the Claude-only zero-discovery test. |
| `internal/ensigncycle/auto_continue_live_test.go` | delete 89 lines | Remove the Claude-only auto-continue test. |
| `internal/ensigncycle/auto_continue_pi_live_test.go` | delete 194 lines | Remove the Pi-only wrapper and move its reusable fixture. |
| `internal/ensigncycle/ac2_reanchor_live_test.go` | delete 30 lines | Remove the Claude-only AC re-anchor wrapper. |
| `.github/workflows/runtime-live-e2e.yml` | 8 deletions, 4 insertions | Remove standalone selectors and keep the shared selector. |
| `internal/release/workflow_exec_guard_test.go` | 20 insertions, 10 deletions | Reject retired selectors and require the shared selector. |
| `docs/runtime-live-ci.md` | 35 insertions, 20 deletions | Document 16 canonical journeys, reconciliation, and local commands. |

The separate AC re-anchor repair owns changes to `internal/livescenario/ac2_reanchor.go`. The Pi adapter task owns Pi transport files.

The source-binding task can change exact insertion counts before this implementation starts. The percentage and file-count tolerance absorb annotation movement only.

## Observable semantic boundary

- Command grammar: no change.
- Stored formats: no change.
- Workflow authority: no change.
- Product runtime behavior: no change.
- Live test identity: six journeys move to canonical shared subtests.
- Fixture behavior: inline setup becomes stable, annotated fixture builders.
- CI behavior: each implemented common adapter runs six additional registry journeys.
- Artifact and metrics behavior: records use canonical journey IDs and stable fixture labels.
- Runtime transport: no change. Adapters retain authentication, launch, parsing, liveness, and host-only evidence.

## Documentation diff

In `docs/runtime-live-ci.md`, replace the standalone Claude description:

> Runs `TestLiveEnsignCycle` (the full-cycle smoke) and `TestLiveSharedScenarios` (the shared suite).

With:

> Sets `SPACEDOCK_LIVE_RUNTIME=claude` and runs `TestLiveSharedScenarios`. Its 16 canonical subtests include full-cycle, gate-transition, boot, and auto-continue behavior.

Add this local selector example:

```text
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 40m -run '^TestLiveSharedScenarios$/^default-headless-gate-stop$' ./internal/ensigncycle -v
```

Update the reconciliation result from six missing journey bindings to 16 bound journey IDs. Record the new reconciliation SHA by the approved two-commit procedure.

Do not change `docs/runtime-live-ci-registry.md`. Its IDs, outcomes, and fixture IDs are already the desired contract.

## Stage-specific test gates

- Ideation identifies overlaps with gate-guardrail and avoids duplicate assertions.
- Implementation changes no runtime-specific transport behavior.
- Validation runs promoted journeys through available live adapters plus full and race suites.

## Stage Report: ideation

- DONE: Reconcile every standalone common behavior with the desired registry and the canonical shared-suite plan.
  The map accounts for six missing journeys, including the seed-omitted `ac-value-reanchor`, their fixtures, graders, dependencies, and retired identities.
- DONE: Spike the closest assertion overlap and prove that promoted journeys keep distinct observable outcomes without duplicate grading.
  `TestGateJourneyOverlapSpike` passed both baselines and rejected both cross-qualified starts through one shared `assertGateHeld` primitive.
- DONE: Produce a complete plan with fixture contracts, variant handling, expected files, line estimates, selectors, and acceptance checks.
  The plan defines six fixture contracts, two auto-continue variants, 16 files, exact selectors, negative controls, documentation, and semantic limits.

### Summary

The plan promotes all six missing registry journeys into the one canonical shared suite. It preserves stable IDs and runtime-neutral grading while adapters retain transport behavior.
