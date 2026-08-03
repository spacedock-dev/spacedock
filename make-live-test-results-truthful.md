---
title: Make live test results tell the truth
status: ideation
source: "Captain recarve of live-test-truth, 2026-08-03. Absorbs 1a and wp as design inputs."
score: 1.0
sprint: live-test-truth
group: truthful-results
sprint-readiness: ready
id: 3d2rqxrgvqky085mn170x3zp
gates:
    version: 1
    records:
        - id: gate:3d2rqxrgvqky085mn170x3zp:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3d2rqxrgvqky085mn170x3zp-backlog-1
              briefing:
                id: briefing:3d2rqxrgvqky085mn170x3zp:backlog:attempt-1:revision-1
                digest: sha256:a7eebcb6d8b2b5e90be67cd6c6ada2294f0a241bec98632fa93596e4c902b33b
                request-digest: sha256:6437632e4af360e5c2a901a16ea56e7362e805727387ffb8937e12284a63e91e
                room-ref: ./make-live-test-results-truthful/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3d2rqxrgvqky085mn170x3zp:backlog:1
                briefing: briefing:3d2rqxrgvqky085mn170x3zp:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T12:18:43.129773Z"
                decision: approve
                reason: Captain explicitly approved the outcome-shaped recarve and directed immediate redispatch.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:3d2rqxrgvqky085mn170x3zp:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3d2rqxrgvqky085mn170x3zp-ideation-1
              briefing:
                id: briefing:3d2rqxrgvqky085mn170x3zp:ideation:attempt-1:revision-1
                digest: sha256:f97998b1bff713e954f7193f927be3ed5b6d980435df74c398f21f2c93e9c884
                request-digest: sha256:42cad5cfdc0b6ea8c20d137ac5e21fcc2898d88924a0cfa5d3e8e315a1fbd127
                room-ref: ./make-live-test-results-truthful/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T14:22:46.909612Z"
                reason: Preflight staff review found cross-member ownership and sequencing defects, and the shared sprint index changed. Withdraw this stale binding before the authorized fold.
            - id: gate-attempt:3d2rqxrgvqky085mn170x3zp-ideation-2
              briefing:
                id: briefing:3d2rqxrgvqky085mn170x3zp:ideation:attempt-2:revision-1
                digest: sha256:e548fefa006f3370bfd6916b7ffeedb3b8fb10e8473e8a983d5b33e5a21596a2
                request-digest: sha256:77ed13df7454fc323fa438511104b0f1007a9bba45b122753f8db4943dd82a47
                room-ref: ./make-live-test-results-truthful/review/ideation/briefing-2
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T14:33:52.416943Z"
                reason: Preflight review closure changed the shared sprint index after attempt 2 was frozen. Replace it with a final package that binds the review artifact.
started: 2026-08-03T12:19:34Z
---

## Outcome

A bad first-officer decision produces a red result. A healthy progressing Codex run does not produce a false timeout.

The task owns truthful classification at both edges. Oracle storage and liveness timers are implementation steps, not separate delivery units.

## Problem

The live suite can report green after a wrong first-officer decision. It can also report red while a healthy Codex run progresses.

The AC2 scenario grades narration and unchanged state. The Codex runner buffers JSONL under one fixed 15-minute deadline.

These faults invert the result at opposite edges. One task must make both faults mechanically visible.

## Completed evidence

### Wrong-green evidence

The `ac2-reanchor-live-scenario-repair` report reproduced the existing false green at commit `2997da9a2`.

The scenario passed in 96.70 seconds after the operative contract clauses were removed. Its entity stayed unchanged at `status: ideation`.

The same spike proved a durable two-way oracle. The correct branch stored `revise`, `feedback`, and `rework`.

The forced incorrect branch stored `approve`, `advance`, and `accepted`. The oracle rejected that branch in 148.38 seconds.

This plan reuses that evidence. If source changes do not invalidate the spike, implementation must not repeat the negative live actuator.

### Codex activity-reset spike

The spike used the real `spacedock codex` front door at commit `17d3f5a71`. It used the host arguments of the shared runner.

The final harness kept `CODEX_HOME` outside the plugin checkout. It also started the launcher from a directory without `.safehouse`.

The run used a 30-second probe quiet budget and a 2-second fixed comparison gate. It produced these results:

- The process exited 0 after 8.72 seconds.
- Four valid JSONL events reset the quiet timer.
- The first event arrived after 1.61 seconds.
- The largest gap between events was 6.46 seconds.
- The final message was `activity-reset-spike-ok`.

Thus, the run continued beyond the fixed comparison gate and stopped cleanly. Stream activity can control Codex liveness end to end.

The 30-second value made the probe faster. It does not change the product quiet budget.

The existing controls also passed:

```text
TestDrainToExitResetsDeadlineOnActivity
TestDrainToExitKillsStalledStream
```

The first control runs longer than four quiet budgets under activity. The second control kills a silent process after one quiet budget.

## Proposed approach

### Visible result 1: wrong decisions produce red results

Repair the standalone AC2 scenario with the proved two-way gate fixture.

1. Replace the held ideation gate with a decision-bearing validation gate.
2. Add parked `rework` and `accepted` target stages.
3. Give the first officer authority for one gate decision.
4. Keep only the raw values: baseline 10,000 bytes, target 8,000 bytes, actual 10,200 bytes.
5. Remove the re-anchor rule and required result from the runbook.
6. Remove decision cues from fixture text.
7. Grade only the durable entity body.
8. Accept only `revise`, `feedback`, and `rework`.
9. Reject `approve`, `advance`, `accepted`, and unchanged state.
10. Preserve the fixture ID `ac-reanchor/means-pass-value-regressed`.

This gate oracle serves AC-1. Narration grading is the simplest alternative, but the completed spike proved that it false-greens.

Terminal merge state is another alternative. It adds unrelated behavior and is not necessary to separate the two decisions.

### Visible result 2: healthy Codex progress does not produce red results

Replace the Codex process deadline with the existing shared `streamWatcher`.

The runner will start one process and stream stdout into `codex-exec.jsonl`. It will apply `quietBudgetDefault`, which is 60 seconds.

Each complete JSONL line will reset this budget.

The runner will keep stderr in `codex-exec.stderr.txt`. It will not retry or start a second Codex process.

On stream silence, the watcher will kill the process. The error will name the last event and the artifact directory.

Move the existing process poller beside `streamWatcher`. Both Claude and Codex runners will use the same timer and process adapter.

This stream watcher serves AC-2 and AC-3. A longer fixed deadline is simpler, but any fixed deadline can false-red progress.

Polling Git state is another alternative. It adds fixture-specific knowledge and misses valid model progress before a durable write.

Keep the suite `-timeout 40m` value unchanged. It remains the backstop for a chatty runaway process or a broken test binary.

This suite backstop serves AC-4. A second scenario deadline is insufficient because it recreates the false-red fault.

## Expected surface

| File | Expected change | Estimate |
|---|---|---:|
| `internal/livescenario/ac2_reanchor.go` | Add the decision fixture and durable oracle | +55 / -50 |
| `internal/livescenario/ac2_reanchor_test.go` | Add correct, wrong, unchanged, and identity cases | +85 / -0 |
| `internal/ensigncycle/ac2_reanchor_live_test.go` | Describe and run the durable branch proof | +8 / -6 |
| `internal/ensigncycle/codex_single_run_test.go` | Stream JSONL, apply quiet liveness, and retain fault evidence | +115 / -45 |
| `internal/ensigncycle/codex_live_runner_test.go` | Select the shared quiet budget | +8 / -8 |
| `internal/ensigncycle/streamwatch_test.go` | Share the existing command poller | +45 / -0 |
| `internal/ensigncycle/live_test.go` | Remove the moved command poller | +0 / -45 |
| `internal/ensigncycle/live_budget_test.go` | Guard the Codex process boundary | +8 / -2 |
| `docs/runtime-live-ci.md` | Describe Codex quiet-progress behavior | +7 / -5 |

Expected total: nine files, 331 insertions, and 161 deletions. Tolerance: 70 insertions and 40 deletions.

No command grammar, stored format, or authority semantics change. Runtime behavior changes only for the live scenario fixtures and Codex liveness.

Failure text will add the last progress event and artifact directory. The suite-wide timeout stays unchanged.

## Acceptance criteria

**AC-1 (VALUE) - The AC2 scenario rejects the wrong durable decision.**

Tested by: the oracle accepts `revise/feedback/rework`. It rejects `approve/advance/accepted` and unchanged state.

**AC-2 (VALUE) - A progressing Codex run stays alive.**

Tested by: a deterministic process runs beyond four quiet budgets under JSONL activity. One focused live Codex journey also passes.

**AC-3 - A stalled Codex run fails within the 60-second quiet budget.**

Tested by: a helper emits one event and stalls. The failure names that event and the artifact directory.

**AC-4 - The suite timeout stays a runaway backstop.**

Tested by: the source guard rejects a fixed Codex deadline. The CI command retains `-timeout 40m`, and full tests pass.

## Test plan

Add the offline tests before implementation.

For AC-1, make the unchanged and wrong stored branches fail against the current oracle. Then implement the durable oracle.

Run:

```bash
go test ./internal/livescenario -run '^TestACReanchor' -count=1
```

Run the correct live AC2 branch once. Do not repeat the completed negative live actuator.

```bash
SPACEDOCK_LIVE_ARTIFACT_DIR="$artifact_dir" go test -tags live -count=1 -timeout 12m -run '^TestLiveReanchorGateRejectsMeansOnlyRegressed$' ./internal/ensigncycle -v
```

For AC-2, add `TestCodexProcessActivityResetsQuietBudget`. The helper must emit JSONL for more than four quiet budgets and exit 0.

For AC-3, add `TestCodexProcessQuietTimeoutPreservesFaultEvidence`. Remove the final event to make this test fail.

The stall test must assert the partial JSONL, stderr, process result, last event, and artifact path.

It must also assert the process kill and one invocation.

Run the focused process and watcher tests:

```bash
go test ./internal/ensigncycle -run 'TestCodexProcess|TestDrainToExit' -count=1
```

Before `ys` lands, run one real journey through the changed Codex process boundary.

`TestLiveCodexSharedScenarios` is pre-`ys` evidence for `3d`:

```bash
SPACEDOCK_LIVE_ARTIFACT_DIR="$artifact_dir" go test -tags live -count=1 -timeout 40m -run '^TestLiveCodexSharedScenarios/gate-guardrail$' ./internal/ensigncycle -v
```

`ys` owns the final canonical selector migration to `TestLiveSharedScenarios`. After `ys` lands, use its canonical selector for final sprint validation.

For AC-4, add `codex_single_run_test.go` to the live-budget source guard. Ban `context.WithTimeout` and `codexScenarioTimeout` on this path.

Make sure that `.github/workflows/runtime-live-e2e.yml` still uses the 40-minute suite backstop. Do not add a scenario timeout.

Run the required checks:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

## Documentation change

Change the Codex liveness paragraph in `docs/runtime-live-ci.md`.

Before:

> A fixed 15-minute wall-clock process limit is its only scenario-level liveness guard. JSONL activity does not extend the deadline.

After:

> The shared stream watcher applies a 60-second quiet budget to each Codex scenario. Each complete JSONL line resets the budget.
>
> On stream silence, the runner kills the process and reports the last event and artifact directory. The 40-minute suite timeout remains a runaway backstop.

`docs/site/contributing/architecture-notes.md` already states the desired shared-watcher behavior. No change is necessary in that file.

## Stage Report: ideation

- DONE: Reuse the completed AC2 false-green and durable-oracle evidence without repeating that spike.
  The plan reuses the 96.70-second false green and both durable decision branches from the completed report.
- DONE: Exercise the unproved Codex activity-reset path end to end before selecting a liveness design.
  A real front-door run exceeded a 2-second fixed gate and exited 0 after four JSONL resets in 8.72 seconds.
- DONE: Produce one plan that makes both wrong green results and wrong red results mechanically visible.
  AC-1 grades durable decisions. AC-2 and AC-3 grade progress and silence through one shared watcher.
- DONE: Use `$simple-english` in pragmatic mode for the complete plan.
  The plan uses short sentences, active voice, consistent terms, and conditions before commands.

### Summary

The plan combines the proved AC2 gate oracle with a proved Codex JSONL reset path. It makes wrong decisions red and keeps healthy Codex progress green.
