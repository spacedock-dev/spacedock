---
title: Make headless Sonnet and Codex spawn implementation before validation
status: implementation
source: "PR #583 run 31320596435, Codex job 93262943132 and Sonnet job 93262943118, 2026-08-09"
started: 2026-08-09T18:34:21Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-codex-headless-implementation-worker-before-validation
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
        - id: gate:98aa776adg66gn823a8gamdq:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:98aa776adg66gn823a8gamdq-ideation-1
              briefing:
                id: briefing:98aa776adg66gn823a8gamdq:ideation:attempt-1:revision-1
                digest: sha256:2aba5b6df9421dd67cb8935d68bbbeaa3143e6833b91e7cbe1e9cfe2152fc04e
                request-digest: sha256:28e6d0584bc432e39011dab29a42972c17f0d05b21b3ff8ba78c7e8fbd844c3b
                room-ref: ./codex-headless-implementation-worker-before-validation/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:98aa776adg66gn823a8gamdq:ideation:1
                briefing: briefing:98aa776adg66gn823a8gamdq:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T21:33:10.060763Z"
                decision: approve
                reason: Captain approved the real worker-before-validation direction.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:98aa776adg66gn823a8gamdq:validation
          stage: validation
          attempts:
            - id: gate-attempt:98aa776adg66gn823a8gamdq-validation-1
              briefing:
                id: briefing:98aa776adg66gn823a8gamdq:validation:attempt-1:revision-1
                digest: sha256:990469b68099b67e2e0b079e4d50be7ea2d03a9bdb915957ac65b18e4ffffee2
                request-digest: sha256:e169af07ba1c664957e035be20095299c4294fbced09f2ed3a294d6bbd95e014
                room-ref: ./codex-headless-implementation-worker-before-validation/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:98aa776adg66gn823a8gamdq:validation:1
                briefing: briefing:98aa776adg66gn823a8gamdq:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T06:38:16.194318Z"
                decision: revise
                reason: Two Material AC-1 findings are confirmed and 98a-owned. Fix the native spawn/completion/report oracle and isolate the observer log under the Captain-approved two-file, 130-gross correction cap.
            - id: gate-attempt:98aa776adg66gn823a8gamdq-validation-2
              briefing:
                id: briefing:98aa776adg66gn823a8gamdq:validation:attempt-2:revision-1
                digest: sha256:0c136007c67438fd4f66b18d71247157dec9b9e7de6de0bd044cc372ca64aedd
                request-digest: sha256:04aa47939257d1690935c10d327599eb84d179ef15d82631fd001344d250fea7
                room-ref: ./codex-headless-implementation-worker-before-validation/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:98aa776adg66gn823a8gamdq:validation:2
                briefing: briefing:98aa776adg66gn823a8gamdq:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-10T07:17:01.629546Z"
                decision: approve
                reason: Corrected candidate 79b539410 satisfies AC-1 through AC-3, resolves both Material cycle-1 findings, and stays within the Captain-approved two-file and 217-gross caps.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:98aa776adg66gn823a8gamdq-validation-3
              briefing:
                id: briefing:98aa776adg66gn823a8gamdq:validation:attempt-3:revision-1
                digest: sha256:3372d5aad9060fea4dfd1d9886d00f457b32028313370e77a8e1f9a684a3adb8
                request-digest: sha256:74e2de026a129464c6011dda5c36790632f7b13dbb437a0cb81052eda6c66a43
                room-ref: ./codex-headless-implementation-worker-before-validation/review/validation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:98aa776adg66gn823a8gamdq:validation:3
                briefing: briefing:98aa776adg66gn823a8gamdq:validation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-08-10T07:27:47.0654Z"
                decision: approve
                reason: Exact rebased candidate 6e281f31c preserves 2e4 and 98a reconciliation, all focused and required offline checks pass, and no Material finding remains.
              application:
                target-stage: done
                state: superseded
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

Wait for task ts7gq0mr9s3chx2w4wppd1kt to commit target-level XFAIL bindings for both
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
must run the Sonnet and Codex cells as target-level XFAIL evidence.

Do not change First Officer behavior before both cells report target-level
XFAIL. After the repair, XPASS
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
- DONE: Define the smallest portable worker-spawn repair after committed target-level XFAIL evidence.
  The proposed repair is one shared-core initial-dispatch guard. It preserves
  adapter mappings, the advance reuse handle, the fixture, and the oracle.
  Product work waits for task ts7gq0mr9s3chx2w4wppd1kt target-level XFAIL evidence.
- DONE: Give gross and net line estimates with exact-candidate live proof.
  The repair estimate is 8 insertions, 2 deletions, and 6 net lines, with plus
  or minus 2 lines tolerance and no extra files. The evidence candidate is
  7c8708c8537fc73761e56813ddbd6a498959ef19; Codex took 166.56 seconds and
  Sonnet took 414.17 seconds.

### Summary

Ideation records an authenticated failure on the exact candidate and defines a
small, portable shared-core guard. The next stage must land target-level XFAIL
evidence first, then prove spawn and completion ordering in both live lanes.
The entity parser and gofmt check passed. The normal and race Go suites were
attempted, but internal/cli stalled and both runs were interrupted after about
five minutes; no product files changed.

## Stage Report: implementation

- DONE: Implement the approved shared First Officer worker-spawn and completion guard for the initial implementation stage.
  Commit 29d11fed3 adds the five rules at `skills/first-officer/references/fo-dispatch-core.md:21`.
- DONE: Add task-local focused proof for real spawn evidence, completion before validation, and rejection of build/status/narration/empty-wait substitutes.
  `TestInitialWorkerSpawnGuardPrecedesCompletionAndValidation` fails if a rule is absent or out of order.
- DONE: Do not edit shared XFAIL bindings, registry reconciliation, sprint package files, or shared runtime documentation before ts lands.
  The checkpoint changes only the task-owned dispatch core and its focused contract test.
- DONE: Keep runtime adapters, command grammar, stored formats, and product ownership unchanged.
  The diff changes no adapter, command implementation, stored format, or ownership rule.
- DONE: Commit the product behavior checkpoint and report exact files, lines, tests, and deferred post-ts live proof.
  Commit 29d11fed3 changes the dispatch core at line 21 and the test at line 20.
- DONE: Run the task-local normal and race proof.
  `go test ./internal/contractlint -count=1` and its race form pass. Removing or reordering a guard makes the focused test fail.
- FAILED: Complete the repository-wide normal and race suites.
  Both `go test ./...` commands stalled in `internal/cli`. The runs stopped after 334.291s and 254.890s.
- SKIPPED: Run the Sonnet and Codex live acceptance proof.
  The post-ts stage owns the exact-candidate live runs, XFAIL removal, and registry reconciliation.

### Summary

Commit 29d11fed3 requires a real initial spawn and verified completion before validation. The focused normal and race tests pass.
The shared XFAIL bindings, registry, sprint package, runtime documentation, adapters, command grammar, and stored formats remain unchanged.

## Stage Report: implementation (cycle 2)

- DONE: Rebase the product work onto `origin/main` at or after `a8688cabf`.
  The branch rebased without a conflict and preserved the original product checkpoint as `e407abc61`.
- DONE: Implement the approved shared First Officer worker-spawn and completion guard for the initial implementation stage.
  Commits `1ab889f1e`, `a371ff490`, `f9935a20a`, and `b2c1ad9ba` complete the guard at `fo-dispatch-core.md:17`.
- DONE: Add task-local focused proof for real spawn evidence, completion before validation, and rejection of build/status/narration/empty-wait substitutes.
  The three focused tests start at `initial_worker_spawn_guard_test.go:28`. Removing a pinned boundary makes a test fail.
- DONE: Run each exact owned live target on the same committed candidate.
  Codex passed in 228.03s. Sonnet passed in 408.55s on commit `a00bd2c97`.
- DONE: Remove only the 98a-owned XFAIL bindings after passing product evidence.
  `shared_live_runner_test.go:110` retains only the Pi binding owned by `fh6rv0k6wr25zty0jjan4jp7`.
- DONE: Update the exact reconciliation rows and runtime registry outcome.
  The gap row is at `live_registry_reconciliation_test.go:53`. The outcome is at `runtime-live-ci-registry.md:115`.
- DONE: Keep runtime adapters, command grammar, stored formats, and product ownership unchanged.
  The candidate changes no runtime adapter, command implementation, stored format, sprint package, or other task binding.
- DONE: Run the required checks.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` completed successfully.
- DONE: Commit and push the completed product repair.
  Commit `a00bd2c97` and its five predecessor commits are on the remote task branch.

### Summary

The repair now requires a real initial worker, completion evidence, started stage transitions, and direct gate entry after completion.
Both owned live targets pass without XFAIL bindings. The normal and race suites also pass.
The Pi binding and all bindings owned by other tasks remain unchanged.

## Stage Report: validation

- DONE: Inspect exact candidate a00bd2c97 and the implementation report. Do not change candidate bytes.
  `git rev-parse HEAD` returned `a00bd2c97`. The code worktree stayed clean.
- FAILED: Verify each acceptance criterion with independent evidence, including exact Sonnet and Codex target passes.
  Sonnet passed in 337.11s. Codex failed in 281.74s because gate preparation created no request.
- FAILED: Confirm worker spawn, completion evidence, stage-start stamping, and direct gate entry remain ordered across normal paths.
  Codex spawned and completed the worker, then `gate prepare` rejected the changed `evidence/command.log` reference.
- DONE: Inspect removal of only the 98a-owned XFAIL bindings and exact reconciliation rows.
  The diff removes only the Sonnet and Codex `98aa776adg66gn823a8gamdq` bindings. The reconciliation check passed.
- FAILED: Perform the semantic adversarial pass on missing spawn, premature validation, empty wait, terminal state, and failure cleanup.
  The accepted command-log baseline contains no spawn, completion signal, or implementation report. The gate-state mutants failed as expected.
- DONE: Use existing implementation full and race results. Run only independent focused falsifiers.
  I did not repeat the full or race suites. Four focused contract and four focused gate checks passed.
- DONE: Classify findings with four evidence fields and write a PASSED or REJECTED report in Simplified Technical English.
  The recommendation is REJECTED because two Material findings affect AC-1.

### Material findings

1. Evidence defect: the live oracle can pass without a worker spawn or completion.
   - Released user and normal workflow: A headless Sonnet or Codex run advances from implementation to the validation gate.
   - Observable harm: A skipped implementation worker can receive passing release evidence.
   - Affected authority: `value-ac[AC-1]` requires one implementation-worker spawn and completion before validation.
   - Trigger evidence: `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` passes a log with status, build, validation, and prepare only.
     `runGateStopScenario` checks gate state and this command log. It does not check the host spawn, completion signal, or implementation report.

2. Outcome defect: the exact Codex target does not pass on the candidate.
   - Released user and normal workflow: A headless Codex run must complete implementation and present the first validation gate.
   - Observable harm: The run stops at validation without a prepared gate request.
   - Affected authority: `value-ac[AC-1]` requires a passing Codex `TestLiveCommonDefaultHeadlessGateStop` run on the candidate.
   - Trigger evidence: The exact target failed in 281.74s on `a00bd2c97`.
     `gate prepare` rejected `evidence/command.log` because the command changed that selected source before preparation.

### Deferred risks

None.

### Summary

The exact Sonnet target passed, but the exact Codex target failed. The independent adversarial control also proves that the oracle does not establish worker spawn or completion. I recommend REJECTED.

### Feedback Cycles

- Cycle 1: REJECTED — validation; surface 5 files/96 gross LOC vs estimate 2 files/10 gross LOC (+860%); AC unchanged
- Cycle 2: PASSED — validation re-review; surface 7 files/217 gross LOC vs estimate 2 files/10 gross LOC (+2070%); AC unchanged under Captain-approved reset

## Captain-approved correction package

The Captain approved both Material fixes for implementation cycle 2.
The end-user value remains unchanged: Sonnet and Codex must run one implementation worker to completion before validation.

The correction can change only these two files:

- `internal/ensigncycle/claude_runtime_helpers_test.go`
- `internal/ensigncycle/claude_live_runner_test.go`

Finding 1 is estimated at `+92/-1`, or 93 gross lines and 91 net lines.
Finding 2 is estimated at `+26/-2`, or 28 gross lines and 24 net lines.
The combined estimate is `+118/-3`, or 121 gross lines and 115 net lines.

The correction cap is 130 gross changed lines across exactly these two additional files.
The projected full candidate cap is seven files and 217 gross changed lines.
The package adds no production, contract, registry, documentation, or XFAIL file.

Captain approval authorizes the two rejected-finding fixes within these limits.

## Stage Report: implementation (cycle 3)

- DONE: Require one native implementation-worker spawn and its matching completion before validation.
  The live oracle reads Claude events and the Codex parent rollout. A command-only baseline now fails.
- DONE: Require one parsed implementation report with a DONE result.
  The oracle uses the fence-safe section parser and rejects a missing or incomplete report.
- DONE: Move the observer log outside the workflow.
  The runner passes an external path to the shim and copies only the finished log into the live artifacts.
- DONE: Add the observer isolation control.
  The focused control rejects `evidence/command.log` inside the workflow and accepts an external temporary path.
- DONE: Keep both correction limits.
  The correction changes two approved files and 121 gross lines. The full candidate has seven files and 217 gross lines.
- DONE: Run the focused falsifiers.
  The normal and live-tag focused tests pass. The command-only and in-root controls fail the corrected assertions as required.
- DONE: Run the exact Codex and Sonnet targets.
  Codex passed in 229.76s. Sonnet passed in 262.46s after one model-variance retry.
- DONE: Run the required checks.
  `gofmt`, `go test ./...`, and `go test ./... -race` completed successfully on the final bytes.
- DONE: Commit and push the correction.
  Commit `79b539410` is on the remote task branch.

### Summary

The live target now proves one native implementation worker, matching completion, and a DONE report before validation.
The harness observer no longer changes the workflow that Codex inspects.

## Stage Report: validation (cycle 2)

- DONE: Re-review exact corrected candidate 79b539410. Preserve candidate bytes.
  `git rev-parse HEAD` returned `79b539410`. The code worktree stayed clean.
- DONE: Reproduce both cycle-1 negative controls: missing worker lifecycle and in-workflow observer path.
  The command-only lifecycle failed the oracle. The in-workflow observer path also failed its boundary check.
- DONE: Verify native Claude and Codex spawn, matching completion, DONE implementation report, and ordering before validation.
  Each exact target passed the native lifecycle oracle with one spawn, its matching completion, one DONE report, and later validation.
- DONE: Verify exact Codex and Sonnet passing evidence and the isolated finished observer artifact.
  Codex passed in 220.21s. Sonnet passed in 297.53s. Each artifact contains its finished external `command.log`.
- DONE: Confirm the correction stays within two test files, 121 gross lines, and the approved 217-gross full cap.
  The correction is 120 insertions and one deletion across two approved files. The full candidate is 217 gross lines across seven files.
- DONE: Re-check all ACs, classify any new finding, and write the cycle-2 PASSED or REJECTED report in Simplified Technical English.
  AC-1 and AC-2 have independent native lifecycle evidence. AC-3 has a passing reconciliation check on the same candidate.

### Acceptance criteria

- AC-1: PASSED.
  Both exact runtime targets observed one native implementation worker and its matching completion before validation.
- AC-2: PASSED.
  The command-only control cannot replace native lifecycle evidence. The report parser also requires one implementation report with a DONE result.
- AC-3: PASSED.
  The reconciliation check confirms that only the owned Sonnet and Codex bindings remain removed.

### Findings

No Material finding, deferred risk, or polish finding remains.

### Summary

The corrected candidate resolves both cycle-1 findings. The exact Codex and Sonnet targets pass with native lifecycle and isolated observer evidence. I recommend PASSED.

## Stage Report: implementation merge rework

- DONE: Rebase the exact candidate `79b539410` onto landed main `f2428ddee6b5b52d1478be25c345e47aae43aa97`.
  The rebased candidate is `6e281f31c624eab611c8a3ba6f05adf0f17e4229`.
- DONE: Resolve only `internal/contractlint/live_registry_reconciliation_test.go`.
  The result keeps the landed `gate-guardrail` removal and the 98a Sonnet and Codex removals.
- DONE: Verify the clean merge of `shared_live_runner_test.go`.
  The gate-guardrail inventory is empty. The default-headless gate inventory contains only Pi issue `fh6rv0k6wr25zty0jjan4jp7`.
- DONE: Run the focused reconciliation and lifecycle falsifiers.
  Both focused checks passed. The live-tag compile and focused run also passed.
- DONE: Run the required checks.
  `gofmt`, `go test ./...`, and `go test ./... -race` completed successfully.
- DONE: Preserve the approved product surface.
  The candidate still changes exactly seven files and 217 gross lines. Product behavior did not change during the rebase.
- DONE: Push the rebased candidate.
  The remote task branch now points to `6e281f31c624eab611c8a3ba6f05adf0f17e4229`.

### Summary

The merge rework preserves the landed 2e4 inventory and the approved 98a repair. All required checks pass on the rebased candidate.

## Stage Report: validation merge recheck

- DONE: Inspect exact rebased candidate 6e281f31c624eab611c8a3ba6f05adf0f17e4229. Preserve bytes.
  `git rev-parse HEAD` returned the exact candidate. The code worktree stayed clean.
- DONE: Verify the manual reconciliation retains landed 2e4 removals and only 98a-owned Sonnet/Codex removals.
  `gate-guardrail` has no gap. `default-headless-gate-stop` retains only the Pi binding owned by `fh6rv0k6wr25zty0jjan4jp7`.
- DONE: Verify shared runner inventory, lifecycle falsifiers, live-tag compile, and implementation full/race evidence.
  Reconciliation, both falsifiers, and the live-tag compile passed. The implementation report records passing full and race suites.
- DONE: Confirm the seven-file, 217-gross bounded surface and unchanged product semantics.
  The diff against rebased main has seven files and 217 gross lines. The rebase adds only landed upstream behavior.
- DONE: Recommend PASSED or REJECTED for the exact rebased merge candidate in Simplified Technical English.
  I recommend PASSED. No Material finding, deferred risk, or polish finding remains.

### Summary

The merge reconciliation preserves the landed 2e4 inventory and the approved 98a repair. The focused checks pass on the exact rebased candidate.
The exact live runs were not repeated. The required PR lanes will run on this head.

## Captain decision: validation and PR cycle 3

The Captain approved the 98a AutoContinue finding as a Material, task-owned FIX.
The end-user value stays unchanged: completed implementation must dispatch fresh validation before its gate.

The correction replaces one existing guard line in `skills/first-officer/references/fo-dispatch-core.md`.
The line must dispatch the stage under normal freshness rules before gate lifecycle.
It must preserve the implementation-worker completion guard and every XFAIL binding.

The correction estimate is `+1/-1`, or two gross lines and zero net lines.
The full candidate must remain at seven files and 217 gross lines.
This approval authorizes validation and PR cycle 3 within these limits.

## Stage Report: implementation (PR cycle 3)

- DONE: Dispatch fresh validation before its gate.
  The guard stamps the stage, applies normal freshness rules, and enters gate lifecycle only after worker completion and report acceptance.
- DONE: Update only the matching contract expectation.
  The focused contract check fails if the guard omits the `started` transition or retains the obsolete no-dispatch rule.
- DONE: Preserve the worker-completion guard and external bindings.
  No worker-lifecycle line, XFAIL binding, registry row, runtime adapter, or stored format changed.
- DONE: Run the focused negative controls.
  AutoContinue rejects status-only advancement. The lifecycle controls reject missing worker evidence and an in-workflow observer.
- DONE: Use the completed exact Sonnet evidence.
  AutoContinue passed in 653.04s. Default-headless passed in 306.88s with the native implementation lifecycle oracle.
- DONE: Run the required offline checks.
  `gofmt`, `go test ./...`, and the restored `go test ./... -race` run completed successfully.
- DONE: Keep the approved correction limits.
  The correction is two line replacements. The full candidate remains seven files and 217 gross lines.
- DONE: Commit and push PR #657 head.
  Commit `61d31a349ecbccf5f68fe11539575ad749966243` is on the remote task branch.

### Summary

Fresh validation now runs before its gate. The exact Sonnet proofs and all required offline checks pass.
