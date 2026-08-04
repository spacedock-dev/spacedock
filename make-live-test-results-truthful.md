---
title: Make live test results tell the truth
status: validation
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
            - id: gate-attempt:3d2rqxrgvqky085mn170x3zp-ideation-3
              briefing:
                id: briefing:3d2rqxrgvqky085mn170x3zp:ideation:attempt-3:revision-1
                digest: sha256:300b09a445442cb7cb7d61b8a707479749a9bd9acd8110a0bb3216de1f9311e9
                request-digest: sha256:2e880d7988a60bb97bfc6b6b99d08809dcd17978300082a782846d08b6a123c7
                room-ref: ./make-live-test-results-truthful/review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:3d2rqxrgvqky085mn170x3zp:ideation:3
                briefing: briefing:3d2rqxrgvqky085mn170x3zp:ideation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-03T14:45:21.039996Z"
                decision: approve
                reason: Approved after staff review. Land truthful live-result oracles first.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:3d2rqxrgvqky085mn170x3zp:validation
          stage: validation
          attempts:
            - id: gate-attempt:3d2rqxrgvqky085mn170x3zp-validation-1
              briefing:
                id: briefing:3d2rqxrgvqky085mn170x3zp:validation:attempt-1:revision-1
                digest: sha256:2a92323df769b11785a7a522557ef05797ae25b3f44f412688c408f4bd1c94df
                request-digest: sha256:083a25c48521ad3b983c0f875f87c5750ee1890a05732c73f7ea43519c0677de
                room-ref: ./make-live-test-results-truthful/review/validation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-04T15:00:14.613395Z"
                reason: 'Validation attempt is stale: its sole repository blocker was the v1 manifest drift, which merged in PR #612; replace the attempt before fresh validation at candidate tip 4d2d45cd5.'
            - id: gate-attempt:3d2rqxrgvqky085mn170x3zp-validation-2
              briefing:
                id: briefing:3d2rqxrgvqky085mn170x3zp:validation:attempt-2:revision-1
                digest: sha256:8869f0652edb6b43aa8fd0927bf59d9f5a4cd230afb807d55e71828fe40539d4
                request-digest: sha256:d0e193fc453de044a59344378fd36b36678162b9da1ea10ec082d105044508ca
                room-ref: ./make-live-test-results-truthful/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:3d2rqxrgvqky085mn170x3zp:validation:2
                briefing: briefing:3d2rqxrgvqky085mn170x3zp:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-04T16:07:51.851039Z"
                decision: approve
                reason: 'Fresh independent evidence accepts the direction: wrong durable AC decisions now fail, meaningful complete Codex JSONL progress prevents false timeout, true silence still kills with retained evidence, and all required suites pass.'
              application:
                target-stage: done
                state: pending
started: 2026-08-03T12:19:34Z
worktree: .worktrees/spacedock-ensign-make-live-test-results-truthful
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

## Stage Report: implementation

- DONE: Make the AC re-anchor scenario grade the durable gate branch.
  The neutral fixture now records raw measurements at `validation`. The oracle accepts only `revise`, the configured `feedback` route, and parked status `rework`. Offline controls reject approve, advance, accepted, and unchanged states even when the final message narrates rejection.
- DONE: Replace the Codex fixed process deadline with shared stream liveness.
  One `codex exec` process streams complete JSONL lines through `streamWatcher`. Each line resets the 60-second quiet budget. Silence kills and reaps that process without retry.
- DONE: Preserve useful fault evidence and the suite backstop.
  Stall failures retain JSONL, stderr, the process result, the last event, and the artifact directory. Source guards ban the removed Codex deadline, and `.github/workflows/runtime-live-e2e.yml` still uses `-timeout 40m`.
- DONE: Prove both changed live paths.
  The AC2 positive branch passed in 254.46 seconds with durable `decision: revise` and status `rework`; artifacts are in `/tmp/spacedock-live-test-truth-ac2-retry.QV9MqO`. An earlier run in `/tmp/spacedock-live-test-truth-ac2.kkSVCQ` exposed a contradictory stop instruction after the model chose revise; the runbook was corrected, and no negative actuator was repeated. The pre-`ys` Codex `gate-guardrail` journey passed in 117.94 seconds with `gpt-5.6-luna` at maximum reasoning; `/tmp/spacedock-live-test-truth-codex.bz7ka2/codex-process-result.txt` records exit 0, no timeout, and the final event.
- DONE: Run formatting, focused controls, full tests, and race tests.
  `gofmt -w ./cmd ./internal`, live-tag compilation, the AC controls, the Codex activity/stall controls, and five focused race repetitions passed. A raw `go test ./...` exposed seven paths that the current shared state branch archived while `v1_pilot_manifest.txt` still names them. Against immutable state snapshot `73f41e2a2232ebb561710bce568641ec976d5f3d`, which contains all 31 declared fixtures, both `go test ./...` and `go test ./... -race` passed.
- DONE: Keep the approved implementation surface.
  Commit `fcd8f7881` changes the expected nine files by 400 insertions and 172 deletions, within the approved tolerance.

### Summary

Wrong AC decisions now fail on durable state, while healthy Codex JSONL progress resets one shared quiet budget. The implementation preserves one-process execution, actionable stall artifacts, and the 40-minute suite backstop.

### Feedback Cycles

- Cycle 1: REJECTED — detached validation/First Officer; surface 9 files/571 changed LOC vs estimate 9 files/492 changed LOC (116%); AC unchanged. Material task-owned AC-1 evidence defect: the partial oracle accepted malformed gate state and an older `revise` that masked a later `approve`. FO disposition: fix. Replace the partial parser and synthetic action with canonical complete-record validation plus the workflow-derived `feedback-to` route; retain malformed, duplicate-record, later-approve, wrong-route, wrong-decision, and unchanged-state negatives. The authorized correction remains within the approved nine-file +401/-201 tolerance. The separate seven-path current manifest drift is outside 3d and snapshot substitution does not count as canonical full-suite evidence.

## Stage Report: validation

- DONE: Reproduce every AC against independent evidence, including the durable correct/wrong/unchanged AC2 branches and the progress/stall Codex controls; reject any oracle that can pass while the visible result is wrong.
  AC-1's correct branch passes; malformed, duplicate-record, later-approve, wrong-route, wrong-decision, accepted-status, and unchanged controls fail through canonical `gates.Read` validation and workflow-derived `feedback-to` routing.
- DONE: Perform the mandated detached adversarial audit on a throwaway checkout, tracing complete JSONL identity, order, lifecycle, quiet-reset cardinality, kill/reap cleanup, one-process execution, retained fault artifacts, the 60-second quiet budget, and the unchanged 40-minute suite backstop.
  The audit observed 31 ordered complete events, one invocation, byte-identical retained JSONL, progress past four quiet budgets, and a 60.44-second kill/reap; CI still selects `-timeout 40m`.
- FAILED: Reproduce the standard repository checks and both changed live paths on the candidate, resolve the current-state manifest failure versus the immutable-snapshot workaround, verify the nine-file +400/-172 surface and pre-ys selector boundary, and classify every finding by defect kind and release scope.
  Both live paths pass and commit `17985b60f` is nine files, +401/-170, on pre-`ys` `TestLiveCodexSharedScenarios`; current-checkout `go test ./...` and `go test ./... -race` fail the same seven stale pilot-manifest paths, so the immutable snapshot is not accepted.
- DONE: AC-1 (VALUE) - The AC2 scenario rejects the wrong durable decision.
  Offline bidirectional controls and detached hostile records pass; the corrected positive live branch passed in 194.60s with artifacts under `/tmp/spacedock-live-truth-rereview-artifacts/ac2`.
- DONE: AC-2 (VALUE) - A progressing Codex run stays alive.
  The deterministic process survives 31 events and more than four quiet budgets; Luna/max `gate-guardrail` passed in 321.98s with exit 0 and no timeout.
- DONE: AC-3 - A stalled Codex run fails within the 60-second quiet budget.
  The scaled stall retains the last event, stderr, process result, artifact path, kill/reap, and one invocation; the production-budget audit failed at 60.44s as required.
- DONE: AC-4 - The suite timeout stays a runaway backstop.
  Source guards reject a fixed Codex deadline, `codex_single_run_test.go` has none, and `.github/workflows/runtime-live-e2e.yml` retains the 40-minute Codex suite timeout.
- FAILED: Material finding — current-state v1 pilot manifest drift (evidence defect; release-blocking; outside 3d ownership).
  Seven entries introduced by code commit `9ff2aa50c` now live under `_archive`; this common current-checkout trigger violates `contract[AGENTS.md#Expected Commands]` and keeps both mandatory suites red until the repository-level manifest/count owner updates the canonical current-state check.
- SKIPPED: Immutable state snapshot substitution.
  Snapshot `73f41e2a2` predates the archive moves and is not the current substrate selected by `TestV1PilotManifestReadsAndValidates`; its green result cannot satisfy the repository gate.
- SKIPPED: Deferred risk — a complete Codex JSONL line above 16 MiB can stop `bufio.Scanner` and deadlock reap.
  The trigger is unobserved and outside the current host evidence (largest live line: 35,627 bytes); promote to Material if a supported Codex lane can emit a line at or above 16 MiB or promises unbounded event size.
- FAILED: Recommendation — REJECTED.
  All candidate value ACs have valid current evidence and the cycle-1 task-owned finding is closed, but the independent material current-checkout repository blocker remains.

### Summary

The corrected candidate makes wrong durable decisions red and healthy JSONL progress green. Validation recommends REJECTED only because the canonical current checkout cannot pass the repository's required full and race suites; the stale manifest belongs outside 3d and must not be hidden by a historical snapshot.

## Stage Report: validation (cycle 2)

- DONE: Independently verify immutable candidate 03bd07e58ef7a74a404c880c73599c2f80de7609 is 0 behind current main db7f1e84aef5df2daf20fb02deac440df4ae1af1 and exactly the approved nine-file +401/-170 surface, with no planning/registry path; reproduce AC-1 through AC-4 and close the prior current-manifest blocker against the live split-root checkout.
  `rev-list` is 0/4 and merge-base is current main; the exact nine-path diff is +401/-170, while current-state `TestV1PilotManifestReadsAndValidates`, full, and race suites now pass and would red on stale declared fixture paths.
- DONE: Run the focused AC re-anchor, Codex process/watcher, live-tag compile, manifest, gofmt, full, and race checks; perform the mandated detached adversarial audit of complete gate-record identity/order and Codex quiet-reset/kill/reap evidence without modifying candidate bytes.
  All named checks passed; throwaway audit mutations red on later approve, mismatched briefing identity, malformed complete record, and duplicate stage, while 31 ordered events reset liveness and a production 60.31s stall killed/reaped one process to exit -1.
- DONE: Verify the already-recorded changed-live-path artifacts remain attributable to the unchanged feature commits, classify all findings, recommend PASSED or REJECTED, append a complete final validation report, and commit only make-live-test-results-truthful.md in split-root state; do not approve or consume a gate.
  Codex `source-head.txt` names `17985b60f`; its changed live paths are byte-unchanged in the candidate, and the AC2 artifact records canonical revise plus rework from the same post-correction run window. No gate action was taken.
- DONE: AC-1 (VALUE) - The AC2 scenario rejects the wrong durable decision.
  Before, a wrong durable decision could go green; now canonical revise/feedback/rework passes while approve, later approve, wrong target/route, malformed/duplicate records, and unchanged state red independently of narration.
- DONE: AC-2 (VALUE) - A progressing Codex run stays alive.
  Before, healthy progress could false-timeout; the exact audit retained 31 complete ordered JSONL events over more than four quiet budgets with one invocation and exit 0, and the unchanged pre-ys live artifact ran 5m21s without timeout.
- DONE: AC-3 - A stalled Codex run fails within the 60-second quiet budget.
  A fresh production-budget audit tripped at 60.31s, retained byte-identical partial JSONL/stderr/result evidence, named the last event and artifact directory, and proved kill plus reap with one invocation.
- DONE: AC-4 - The suite timeout stays a runaway backstop.
  The source guard and live-tag compile pass, the Codex runner uses `quietBudgetDefault`, and CI still selects the pre-ys Codex journey with `-timeout 40m` and no fixed Codex process deadline.
- DONE: Closed findings — cycle-1 partial oracle and current-state manifest drift.
  Both were Material evidence defects: `17985b60f` closes AC-1 at the complete-record boundary, and landed prerequisite `0b8ca50de` makes the current split-root manifest/full/race controls green; neither remains release-blocking.
- SKIPPED: Subscription-backed Pi/Codex/Claude actuators.
  Dispatch reserves those local subscription lanes to the Commander; validation reused attributable live artifacts and reran all deterministic/offline controls plus live-tag compilation.
- SKIPPED: Deferred risk — a complete Codex JSONL line above 16 MiB can stop the scanner.
  The exact trigger remains outside the supported/observed workflow (largest attributable live line: 35,627 bytes); promote to Material if a supported lane emits at least 16 MiB or the product promises unbounded event size.
- DONE: Recommendation — PASSED.
  All four ACs have direct falsifiable evidence, all mandatory current-checkout checks pass, and no Material finding remains.

### Summary

The standalone value is directly proved: wrong durable AC decisions now produce red results, and meaningful complete JSONL progress keeps a healthy Codex run alive. Silence still fails with retained fault evidence, while the 40-minute suite timeout remains only the runaway backstop. Final validation recommends PASSED; candidate topology is scope hygiene only.
