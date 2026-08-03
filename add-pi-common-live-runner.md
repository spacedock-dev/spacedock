---
title: Add Pi as a first-class adapter for every common live journey
status: ideation
source: "Desired live-test registry requires common journeys on every supported runtime, 2026-08-03"
started: 2026-08-03T10:48:52Z
completed:
verdict:
score: 0.9
worktree:
issue:
sprint: live-test-truth
group: pi-common-runner
sprint-readiness: ready
id: tj41e4f404mz7ast3yh9enwc
gates:
    version: 1
    records:
        - id: gate:tj41e4f404mz7ast3yh9enwc:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:tj41e4f404mz7ast3yh9enwc-backlog-1
              briefing:
                id: briefing:tj41e4f404mz7ast3yh9enwc:backlog:attempt-1:revision-1
                digest: sha256:75c779adf9f75759ab31035e7433523ee1a64a5a351eb725e3bc838061a5a7ca
                request-digest: sha256:b2e61d14956245e69ed188a2d46d2d53dd2bc2e1c2af97fa94def22578bc1346
                room-ref: ./add-pi-common-live-runner/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:tj41e4f404mz7ast3yh9enwc:backlog:1
                briefing: briefing:tj41e4f404mz7ast3yh9enwc:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T10:34:56.368971Z"
                decision: approve
                reason: Captain approved the prepared Sol ideation cohort with make it so.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:tj41e4f404mz7ast3yh9enwc:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:tj41e4f404mz7ast3yh9enwc-ideation-1
              briefing:
                id: briefing:tj41e4f404mz7ast3yh9enwc:ideation:attempt-1:revision-1
                digest: sha256:ef75e23f510c2b349155973dda33caeab6684190f275bbd11f74492a1a8b92e6
                request-digest: sha256:2ed73cad1e0b03f316e6bbf5cf30c1811912f012caab5e7e1885218ae50811b3
                room-ref: ./add-pi-common-live-runner/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T12:16:42.029799Z"
                reason: Captain recarved the sprint into outcome-shaped delivery units. Preserve this report as design input; do not present or consume this component-shaped attempt.
---

## Problem

Pi is a supported runtime target. It has a coverage map, a front-door smoke, and one quarantined recorded-gate test.

Pi has no runner for the common live journey registry. As a result, no common journey has proof across all supported runtimes.

## Spike result

The spike used the existing `shallow-boot` fixture, prompt, observation builder, and host-neutral assertion. It used the real `spacedock pi` front door.

The first run reached Pi and OpenRouter. OpenRouter returned HTTP 402 because the `gpt-5.4` request reserved a 128K output limit.

The second run used `openrouter/openai/gpt-5.4-mini`. It passed in 28.21 seconds with these results:

- The root used 7 model requests.
- Usage was 18,304 input tokens, 1,635 output tokens, and 78,848 cache-read tokens.
- OpenRouter reported a cost of $0.0269991.
- The final message contained both required held-gate lines.
- The entity and Git state stayed byte-identical.
- Pi created no archive, worktree, or child-agent artifact.

The one-off spike test left no repository file. The diagnostic artifacts are under `/tmp/spacedock-pi-common-spike-20260803-mini/` on the spike host.

This result proves the front door, isolated Pi home, local package, session, final-message, and shared assertion path. The cost error was harness friction.

## Proposed approach

Add `piScenarioRunners()` beside the Claude and Codex runner maps. Each map entry will use an existing shared fixture, prompt, and assertion.

Add `TestLivePiSharedScenarios` as the only canonical Pi common-suite entry point. The suite will run every registry entry without a coverage substitute.

The Pi runner adapter will own these runtime facts:

- Package: load the current checkout and the pinned `pi-subagents` and `pi-intercom` packages from isolated Pi state.
- Child agent: grade Pi run IDs, child session records, model stamps, and completion records from `pi-subagents` artifacts.
- Session: create one root session directory and isolated Pi home for each journey.
- Artifact: retain root output, root session JSONL, child JSONL, process result, cost record, and durable workflow evidence.
- Launch: call `spacedock pi --plugin-dir <checkout> -- --print --model <model> --session-dir <dir>`.
- Liveness: apply one fixed 12-minute process deadline per journey. Do not retry a failed live journey.
- Observation: use Pi root assistant events for final text and Pi tool events for mechanism-specific traces.

The common registry, fixtures, prompts, and durable assertions will contain no Pi condition. This boundary serves AC-1 and AC-2.

The Pi trace reader is necessary for `filing` and `smallest-sufficient-mechanism`. Durable state alone cannot distinguish the required command and worker choices.

The simplest alternative was the current coverage map. It cannot exercise a journey, so it cannot serve AC-1 or AC-3.

## Complete journey order and price

The estimates use `openrouter/openai/gpt-5.4-mini`. On 2026-08-03, OpenRouter listed these rates in its [model API](https://openrouter.ai/api/v1/models):

- Input: $0.75 per million tokens.
- Output: $4.50 per million tokens.
- Cache read: $0.075 per million tokens.

The spike cost matches these rates. The estimates below are per-journey planning ceilings, not pass criteria.

| Order | Common journey | Expected Pi work | Price ceiling |
|---:|---|---|---:|
| 1 | `shallow-boot` | No mutation and no child | $0.03 |
| 2 | `merge-hook-guardrail` | Guard refusal and no mutation | $0.06 |
| 3 | `gate-guardrail` | Retained gate package and root presentation | $0.06 |
| 4 | `feedback-3-cycle-escalation` | Durable human handoff and no fourth cycle | $0.07 |
| 5 | `self-evidence-merge-triage` | Root evidence decision | $0.08 |
| 6 | `filing` | One atomic seed write | $0.08 |
| 7 | `recorded-gate-lifecycle` | Three authority barriers and one successor child | $0.15 |
| 8 | `smallest-sufficient-mechanism` | Direct edits plus two commissioned children | $0.25 |
| 9 | `rejection-flow` | Two review cycles and fresh Pi child runs | $0.35 |
| 10 | `keep-moving-posture` | Three concurrent children and one corrected hold | $0.45 |

The full estimate is $1.55. A 25 percent reserve gives a $1.94 planning ceiling for one complete run.

The runner will use these cost controls:

- Run the two non-mutating journeys first.
- Continue in the table order from low cost to high cost.
- Stop the suite after the first failed journey.
- Use no automatic live retry.
- Keep the parent and child model in the Pi adapter configuration.
- Record tokens, cost, duration, process result, and session paths for each journey.
- Keep the live CI environment approval before the suite starts.

## Journey bindings and proof

Every Pi binding uses the same shared fixture and host-neutral assertion as the existing hosts.

| Journey | Shared proof | Pi-only evidence |
|---|---|---|
| `shallow-boot` | `assertShallowBoot` | No Pi child metadata |
| `merge-hook-guardrail` | `assertMergeHookGuardHeld` | Root final assistant text |
| `gate-guardrail` | `assertGateHeld` and command-log proof | Ordered root assistant presentation |
| `feedback-3-cycle-escalation` | `assertThirdCycleEscalation` | Root session and durable Git state |
| `self-evidence-merge-triage` | `assertSelfEvidenceMergeTriage` | Root final assistant text |
| `filing` | Existing entity-path proof | Pi tool trace proves `spacedock new` |
| `recorded-gate-lifecycle` | `assertRecordedGateLifecycle` | Root event order and one child model stamp |
| `smallest-sufficient-mechanism` | `assertDurableSmallestMechanism` | Pi tool and child-run trace |
| `rejection-flow` | `assertRejectionFlow` | Distinct Pi child run IDs for fresh cycles |
| `keep-moving-posture` | `assertDurableKeepMoving` | Three child records and one held task |

## Expected surface

| File | Planned change | Estimated insertions |
|---|---|---:|
| `internal/ensigncycle/pi_shared_runner_live_test.go` | Add Pi setup, launch, timeout, sessions, artifacts, and cost records | 240 |
| `internal/ensigncycle/pi_shared_scenarios_live_test.go` | Add the runner map and ten shared journey bindings | 280 |
| `internal/ensigncycle/pi_session_trace_live_test.go` | Parse Pi root and child events for host-specific proof | 140 |
| `internal/ensigncycle/shared_coverage_meta_test.go` | Require a callable Pi runner in both map directions | 25 |
| `internal/ensigncycle/pi_shared_coverage_test.go` | Delete the coverage-map substitute | 0, delete about 81 |
| `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` | Move its useful proof into the shared runner and delete quarantine code | 0, delete about 104 |
| `.github/workflows/runtime-live-e2e.yml` | Run the complete Pi shared suite as required evidence | 10 |
| `docs/runtime-live-ci.md` | Document Pi parity, the command, costs, and artifacts | 35 |

Expected total: eight files, about 730 insertions, and about 220 deletions. Tolerance: two extra files or 30 percent more insertions.

No product package, command grammar, stored format, or authority rule changes. Runtime test behavior and CI evidence change.

## Test-first implementation plan

1. Change the parity helper to require a callable Pi runner. Prove that an injected missing Pi entry fails.
2. Add adapter tests for Pi arguments, isolated directories, artifact paths, session parsing, model stamps, and the 12-minute deadline.
3. Add the two non-mutating journey bindings. Run them before any higher-cost journey.
4. Add the remaining eight bindings in price order. Use the existing shared assertions unchanged.
5. Delete the coverage map and recorded-gate quarantine after all ten runners exist.
6. Change CI and the runtime-live document. Then run the complete proof set.

## Test plan

Offline proof:

```bash
go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions' ./internal/ensigncycle -v
go test ./...
go test ./... -race
gofmt -w ./cmd ./internal
```

The parity negative must fail when one Pi map entry is absent. It prevents a map or missing runner from satisfying completeness.

Cheap live proof:

```bash
SPACEDOCK_PI_LIVE_REQUIRED=1 go test -tags live -count=1 -timeout 15m -run '^TestLivePiSharedScenarios$/(shallow-boot|merge-hook-guardrail)$' ./internal/ensigncycle -v
```

This command proves the real front door, the shared fixtures, no mutation, and the Pi adapter artifacts. A changed entity or child record makes it fail.

Complete live proof:

```bash
SPACEDOCK_PI_LIVE_REQUIRED=1 go test -tags live -count=1 -timeout 40m -run '^TestLivePiSharedScenarios$' ./internal/ensigncycle -v
```

The log must show one real `PASS` for each of the ten registry IDs. A skip, quarantine, absent runner, or coverage entry does not count.

## Documentation diff

Replace this text in `docs/runtime-live-ci.md`:

> one host-neutral scenario table, per-host runner adapters (Claude and Codex today, with Pi tracked through an explicit live/codified/gap coverage map until its shared runners are live-safe)

with:

> one host-neutral scenario table and one live runner adapter for each supported host: Claude, Codex, and Pi

Replace the Pi runner table entry `pi_shared_coverage_test.go` with `pi_shared_runner_live_test.go` and `pi_shared_scenarios_live_test.go`.

Replace the quarantined two-test Pi command with the complete `TestLivePiSharedScenarios` command from the test plan.

## Acceptance criteria

**AC-1 (VALUE)** Pi invokes every common registry journey through the canonical shared entry point. Each journey uses the same fixture contract and host-neutral assertion.
Proven by: the required Pi live lane runs all ten registered journeys with no skip or coverage substitute.

**AC-2** Pi-specific launch, package, child-agent, session, and artifact behavior stays within the Pi adapter. Common scenarios contain no Pi branching.
Proven by: adapter tests and source inspection pass, followed by the complete real Pi run.

**AC-3** A missing or quarantined Pi runner cannot count as evidence. The required suite reports it as a failure or skip.
Proven by: the parity negative fails after removal of one injected Pi runner entry.

## Stage-specific test gates

- Ideation must price live cost and sequence the cheapest non-mutating journeys first without weakening the all-journey end state.
- Spike one existing shared fixture through Pi before designing bulk migration.
- Validation requires the exact Pi live suite and full offline/race tests.

## Stage Report: ideation

- DONE: Spike one existing shared fixture through Pi and record the real result or the exact blocker.
  `shallow-boot` passed through `spacedock pi` in 28.21 seconds with the unchanged shared assertion and no durable mutation.
- DONE: Price and sequence complete Pi coverage for every common registry journey without coverage-map substitutions.
  The ten-journey table totals $1.55, with a $1.94 reserve ceiling and the two non-mutating journeys first.
- DONE: Produce a complete plan with adapter boundaries, expected files, line estimates, acceptance checks, and live-cost controls.
  The plan defines eight files, about 730 insertions, all ten bindings, three AC proofs, and fail-fast live controls.

### Summary

The spike proved that Pi can run an existing shared fixture through the real front door. The plan replaces all Pi coverage substitutes with live runners and keeps host behavior inside the Pi adapter.
