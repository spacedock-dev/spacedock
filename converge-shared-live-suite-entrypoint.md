---
title: Converge common live journeys on one runtime-neutral test entry point
status: ideation
source: "Desired live-test registry common-journey semantics, 2026-08-03"
started: 2026-08-03T10:48:38Z
completed:
verdict:
score: 0.95
worktree:
issue:
sprint: live-test-truth
group: common-runner
sprint-readiness: ready
id: h3b5tgk77vx9qqdmbjtpsh98
gates:
    version: 1
    records:
        - id: gate:h3b5tgk77vx9qqdmbjtpsh98:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:h3b5tgk77vx9qqdmbjtpsh98-backlog-1
              briefing:
                id: briefing:h3b5tgk77vx9qqdmbjtpsh98:backlog:attempt-1:revision-1
                digest: sha256:3f30288713a92f34c70c13890ae1c396a4a4a70eab5c70e1cf2d4f9275bcd122
                request-digest: sha256:62c1eb7f852de5d5c7cf6a56772e5d61aa7c59d77ad8b1b66cd5b704f2ea9e1f
                room-ref: ./converge-shared-live-suite-entrypoint/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:h3b5tgk77vx9qqdmbjtpsh98:backlog:1
                briefing: briefing:h3b5tgk77vx9qqdmbjtpsh98:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T10:34:18.234548Z"
                decision: approve
                reason: Captain approved the prepared Sol ideation cohort with make it so.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:h3b5tgk77vx9qqdmbjtpsh98:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:h3b5tgk77vx9qqdmbjtpsh98-ideation-1
              briefing:
                id: briefing:h3b5tgk77vx9qqdmbjtpsh98:ideation:attempt-1:revision-1
                digest: sha256:775f2c29f39d9d6ca3113f4b7ca0929d56aad827ae534b015d116b32966d2919
                request-digest: sha256:07744725dc0976cfee2f7e78aaf7069613a9394287fd281cfdfd343fc8a747c3
                room-ref: ./converge-shared-live-suite-entrypoint/review/ideation/briefing-1
---

## Problem

The same common journeys currently enter through separate Claude and Codex top-level tests. This makes runtime selection part of journey identity and permits lane-specific scenario drift. The registry requires one TestLiveSharedScenarios/<journey-id> identity with transport differences behind adapters.

## Acceptance criteria

**AC-1 (VALUE)** A common journey has one test/subtest identity that is invoked unchanged by every implemented runtime lane.
Verified by: focused live invocations of the same named subtest through Claude and Codex that reach the same fixture and host-neutral assertion.

**AC-2** Runtime authentication, launch, output capture, and liveness remain adapter-owned and no common fixture or assertion is forked per host.
Verified by: unit coverage of adapter selection plus source inspection and focused live evidence.

**AC-3** Existing shared scenario parity remains complete while the separate TestLiveClaudeSharedScenarios and TestLiveCodexSharedScenarios entry points are retired or reduced to non-authoritative compatibility wrappers.
Verified by: the shared coverage tests and CI selector behavior.

## Stage-specific test gates

- Spike one shared journey end to end through the proposed unified entry point before migrating the suite.
- Validation runs focused Claude and Codex live evidence plus offline, full, and race suites.

## Ideation plan

### Target behavior

`TestLiveSharedScenarios` becomes the only authoritative entry point for common live journeys. Each runtime lane sets `SPACEDOCK_LIVE_RUNTIME` and invokes the same test selector.

Each journey keeps one subtest identity: `TestLiveSharedScenarios/<journey-id>`. A runtime policy can skip that identity only with an explicit non-applicability or tracked-defect reason.

The common runner owns fixture setup, the prompt, durable state capture, and the host-neutral assertion. The selected adapter owns authentication, launch, output extraction, liveness, and host-output parsing.

### Spike result

The spike used `shallow-boot` because it has the riskiest current split. Its durable oracle is common, but Claude also checks team state and stream measurements.

A throwaway detached worktree added one `TestLiveSharedScenarios/shallow-boot` path and two adapters. The common fixture and durable assertion contained no runtime branch.

Both live runs passed with the exact same selector:

```text
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
PASS: TestLiveSharedScenarios/shallow-boot (36.87s)

SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
PASS: TestLiveSharedScenarios/shallow-boot (28.31s)
```

The Claude adapter retained the no-team, pre-greet, and measured-stream checks. The Codex adapter retained its fixed process deadline and final-message extraction.

### Proposed approach

1. Add one common live entry point and one common journey runner map.
2. Select the adapter from an exact `SPACEDOCK_LIVE_RUNTIME` value of `claude` or `codex`.
3. Return one normalized result with the final message, observed output, artifact directory, duration, and adapter-private evidence.
4. Move each duplicated journey body into one common function.
5. Keep output-dialect parsing behind narrow adapter methods.
6. Keep runtime policy, model-specific skips, metrics, authentication, launch, output capture, and liveness inside each adapter.
7. Replace both CI selectors with `TestLiveSharedScenarios` and set the runtime value in each lane.

The common file must not switch on a runtime name or concrete adapter type. The adapter factory is the only runtime-selection branch.

The shared runner map contains all ten current journey IDs. Claude and Codex use the same map, so a lane cannot add a private common journey.

Runtime supplements remain only for real host differences. Claude keeps its `TeamCreate` and stream-token checks because Codex has neither Claude team storage nor the Claude stream dialect.

### Migration map

| Current surface | Final surface |
|---|---|
| `TestLiveClaudeSharedScenarios` | Removed. The Claude lane invokes `TestLiveSharedScenarios`. |
| `TestLiveCodexSharedScenarios` | Removed. The Codex lane invokes `TestLiveSharedScenarios`. |
| `claudeLiveScenarios()` and `codexLiveScenarios()` | One `sharedLiveScenarios()` list. |
| `claudeScenarioRunners()` and `codexScenarioRunners()` | One common journey runner map plus an adapter registry. |
| Ten `runClaude*Scenario` functions | Ten host-neutral `runShared*Scenario` functions. |
| Ten `runCodex*Scenario` functions | Removed after their common behavior moves to the shared functions. |
| Claude stream and Codex JSONL parsing | Adapter methods that return the same observation types. |
| Separate metrics calls | One adapter method that records a passing journey with the current runtime format. |

`recorded-gate-lifecycle` keeps its runtime setup behind adapter methods. This move removes the current concrete-type switch from common orchestration.

`filing`, `rejection-flow`, and `smallest-sufficient-mechanism` use their current runtime parsers. Each parser returns a shared command or worker observation to one common assertion path.

### Mechanism choices

The adapter registry serves AC-1. Two compatibility wrappers were the simpler alternative, but they preserve two authoritative test identities.

The common journey map serves AC-2 and AC-3. A top-level rename alone is insufficient because separate journey bodies can still drift.

Narrow observation methods serve AC-2. A universal transcript event model adds parser scope that these durable journeys do not need.

The explicit runtime environment value serves AC-1. Credential inference is insufficient because local authentication can exist for both hosts.

### Acceptance criteria

**AC-1 (VALUE)** The count of authoritative common-suite top-level identities decreases from two to one. The remaining identity is exactly `TestLiveSharedScenarios` in both live lanes.

Verified by: a Go test-list count, the workflow execution guard, and focused Claude and Codex runs of the same named subtest.

**AC-2** Every common journey uses one fixture and one host-neutral assertion path. Authentication, launch, output capture, liveness, and output-dialect parsing remain adapter-owned.

Verified by: focused fake-adapter tests, existing liveness tests, output-parser tests, and source review of the common runner file.

**AC-3** All ten current shared journey IDs remain reachable in both implemented runtime lanes. A runtime-specific skip keeps the common identity and states its reason.

Verified by: the revised shared coverage test, both full live lane selectors, and an unchanged `sharedRuntimeScenarios()` ID set.

**AC-4** Existing artifact directories, metrics records, host launch commands, and timeout behavior remain unchanged.

Verified by: adapter unit tests, workflow artifact guards, and focused live evidence from both hosts.

### Acceptance checks and test plan

The implementation starts with focused tests for runtime selection, one top-level identity, and common-map coverage. These tests use fake adapters and spend no model tokens.

The workflow guard must reject either lane that uses an old selector. It must also reject a missing or changed `SPACEDOCK_LIVE_RUNTIME` value.

The shared coverage test must compare the scenario table, the common runner map, and the adapter registry in both directions. A missing or extra entry must fail.

Existing negative fixtures continue to prove the durable assertions. Existing parser tests continue to prove Claude stream and Codex JSONL handling.

The focused offline tests have low cost and finish in seconds. The focused live proof uses two model launches and took about 66 seconds in the spike.

The full live lane cost stays unchanged because this task does not add a journey. Fixture tests, workflow tests, and live tests are all required.

Run these focused offline checks:

```bash
go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestSharedLiveRuntimeSelection|TestSharedLiveEntrypoint' ./internal/ensigncycle -v
go test ./internal/release -run 'TestWorkflowsPreserveAndPublishJourneyCosts|TestRuntimeLiveWorkflowGuardRejectsMissingSharedScenarioRun' -v
```

Run the same focused live journey through both adapters:

```bash
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 20m -run '^TestLiveSharedScenarios$/^shallow-boot$' ./internal/ensigncycle -v
```

Run the required full checks:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

CI runs the full common suite once per Claude matrix entry and once for Codex. Each job uses the same `TestLiveSharedScenarios` selector.

### Expected surface

| File | Estimated change |
|---|---:|
| `internal/ensigncycle/shared_live_runner_test.go` | new, about 330 insertions |
| `internal/ensigncycle/shared_live_runner_unit_test.go` | new, about 90 insertions |
| `internal/ensigncycle/claude_live_runner_test.go` | about 80 insertions and 380 deletions |
| `internal/ensigncycle/codex_live_runner_test.go` | about 70 insertions and 290 deletions |
| `internal/ensigncycle/shared_coverage_meta_test.go` | about 35 insertions and 30 deletions |
| `internal/ensigncycle/journey_metrics_live_test.go` | about 35 insertions and 20 deletions |
| `.github/workflows/runtime-live-e2e.yml` | about 4 insertions and 2 deletions |
| `internal/release/workflow_exec_guard_test.go` | about 25 insertions and 15 deletions |
| `internal/release/journey_workflow_test.go` | about 10 insertions and 6 deletions |
| `docs/runtime-live-ci.md` | about 35 insertions and 25 deletions |
| `docs/site/contributing/architecture-notes.md` | about 8 insertions and 6 deletions |

The estimate is about 720 insertions and 770 deletions across 11 files. The tolerance is 25 percent total and two adjacent test or documentation files.

No fixture or durable assertion file belongs in the expected surface. A change to a journey outcome requires a new ideation decision.

### Observable semantic boundary

This task changes test identity and CI selector semantics. It adds the `SPACEDOCK_LIVE_RUNTIME` test configuration value.

This task does not change command grammar, stored formats, workflow authority, host launch behavior, artifact formats, or liveness budgets.

### Compatibility disposition

The old top-level test names are removed without wrappers. Wrappers would keep duplicate authoritative identities and defeat AC-1.

The repository has no supported external API for Go test names. All in-repository selectors and documentation move in the same change.

CI step names, JSON detail filenames, artifact directories, and journey-metrics paths stay unchanged. Historical diagnostic fixtures can keep old temporary paths because they describe captured runs.

### Documentation diff

In `docs/runtime-live-ci.md`, replace the two local commands with:

Before:

```text
go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v
go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v
```

After:

```text
SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 40m -run TestLiveSharedScenarios ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run TestLiveSharedScenarios ./internal/ensigncycle -v
```

Replace the current runtime-only layer row with these two rows:

```text
| Suite entry point | `shared_live_runner_test.go` (`TestLiveSharedScenarios`) | Yes |
| Runtime adapter | `codex_live_runner_test.go`, `claude_live_runner_test.go`, `pi_shared_coverage_test.go` | No |
```

Replace the GitHub setup sentences with this text:

> The `claude-live` lane sets `SPACEDOCK_LIVE_RUNTIME=claude` and runs `TestLiveSharedScenarios`. The `codex-live` lane sets `SPACEDOCK_LIVE_RUNTIME=codex` and runs the same test.

In `docs/site/contributing/architecture-notes.md`, replace this sentence:

> A per-host runner adapter (Claude and Codex today, with Pi tracked through a live/codified/gap coverage map) turns each scenario into a real launch.

With this text:

> `TestLiveSharedScenarios` owns one subtest identity for each common journey. The lane selects a runtime adapter through `SPACEDOCK_LIVE_RUNTIME`, while fixtures and assertions remain host-neutral.

## Stage Report: ideation

- DONE: Spike one shared journey through one canonical TestLiveSharedScenarios subtest identity on two runtime adapters.
  `shallow-boot` passed as the same named subtest through Claude in 36.87s and Codex in 28.31s.
- DONE: Map the migration from separate Claude and Codex entry points without runtime branches in common fixtures or assertions.
  The migration map moves ten duplicated journey bodies into one common map and keeps runtime differences behind adapters.
- DONE: Produce a complete plan with expected files, line estimates, acceptance checks, selector changes, and compatibility disposition.
  The plan defines 11 expected files, a 25 percent tolerance, exact selectors, tests, semantics, and removal of old names.

### Summary

The spike proved one canonical `shallow-boot` identity across both current runtime adapters. The plan converges all common journeys without changing host launch or liveness behavior.
