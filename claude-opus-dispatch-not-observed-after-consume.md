---
id: 5n4k6djrq8gtvd54zg9s6zhs
title: recorded-gate-lifecycle grader has two unrelated bugs that both produce "successor dispatch not observed after consume"
status: validation
source: "Originally filed 2026-08-02 as an opus-specific, root-cause-unknown candidate flake (TestLiveClaudeSharedScenarios/recorded-gate-lifecycle, PR #600 run 30754109029 job 91513297850). Amended 2026-08-03 after a dedicated diagnostic (opus, high effort) root-caused it with a local repro: the same error string is produced by TWO distinct, deterministic one-line grader bugs in internal/ensigncycle/recorded_gate_lifecycle_test.go, not live-model nondeterminism. Confirmed recurring on a second, independent instance: codex-live hit the identical error string on PR #600's rerun (run 30754109029, job 91518827444) via the SECOND bug below, proving this is not opus-specific. Classification: test-defect, high confidence (opus reproduced locally; codex confirmed from live stdout). Captain directed: file/amend based on the classification workflow's findings."
started: 2026-08-03T05:32:35Z
completed:
verdict:
score: 0.5
worktree: .worktrees/spacedock-ensign-claude-opus-dispatch-not-observed-after-consume
issue:
sprint: durable-decisions
gates:
    version: 1
    records:
        - id: gate:5n4k6djrq8gtvd54zg9s6zhs:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:5n4k6djrq8gtvd54zg9s6zhs-backlog-1
              briefing:
                id: briefing:5n4k6djrq8gtvd54zg9s6zhs:backlog:attempt-1:revision-1
                digest: sha256:8fc216f76f0d60cc8db0e45c265ec26a52da62664706961590568fedabd8dba6
                request-digest: sha256:e726acf13dbe439061d23043e0004e21b5f503a2a101b2bf931ec9158c9403ed
                room-ref: ./claude-opus-dispatch-not-observed-after-consume/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5n4k6djrq8gtvd54zg9s6zhs:backlog:1
                briefing: briefing:5n4k6djrq8gtvd54zg9s6zhs:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T05:20:31.164764Z"
                decision: approve
                reason: 'Captain approved the bounded test-only grader correction: exclude help probes from phase matching and derive the room-backed Resolution ID; preserve gate and dispatch behavior, then rerun exact-head Codex CI.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:5n4k6djrq8gtvd54zg9s6zhs:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:5n4k6djrq8gtvd54zg9s6zhs-ideation-1
              briefing:
                id: briefing:5n4k6djrq8gtvd54zg9s6zhs:ideation:attempt-1:revision-1
                digest: sha256:c25d90024d75d419ab353a4223a6cf894f50e4c290e330645dcfac59af04b323
                request-digest: sha256:cf56fe16fe0a4a5f00777da9a6ba3b906450dbe419ff5fc35b07bd7621efbf9c
                room-ref: ./claude-opus-dispatch-not-observed-after-consume/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5n4k6djrq8gtvd54zg9s6zhs:ideation:1
                briefing: briefing:5n4k6djrq8gtvd54zg9s6zhs:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T05:30:21.448425Z"
                decision: approve
                reason: 'Captain approved the bounded test-only grader correction: exclude help probes from ordered phase matching and resolve room-backed IDs without changing dispatch mechanics.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:5n4k6djrq8gtvd54zg9s6zhs:validation
          stage: validation
          attempts:
            - id: gate-attempt:5n4k6djrq8gtvd54zg9s6zhs-validation-1
              briefing:
                id: briefing:5n4k6djrq8gtvd54zg9s6zhs:validation:attempt-1:revision-1
                digest: sha256:c8d6a31ea31100f7d92a0b05f49f3d2da0951c425edc261a7516fad9756cd942
                request-digest: sha256:b87d46ce3eeb580c7c9bafaf785c3eaf1302ca43e8aed13285a9c54a403cd61d
                room-ref: ./claude-opus-dispatch-not-observed-after-consume/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5n4k6djrq8gtvd54zg9s6zhs:validation:1
                briefing: briefing:5n4k6djrq8gtvd54zg9s6zhs:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T06:46:01.674641Z"
                decision: approve
                reason: Captain approved the validated test-only grader correction; all acceptance criteria and adversarial controls pass, and the normal terminal merge ceremony is authorized.
              application:
                action: advance
                target-stage: done
                state: superseded
                blockers: []
            - id: gate-attempt:5n4k6djrq8gtvd54zg9s6zhs-validation-2
              briefing:
                id: briefing:5n4k6djrq8gtvd54zg9s6zhs:validation:attempt-2:revision-1
                digest: sha256:87ccb2889902d4c84f5ae2c14332c949cad15f0a9f2905734304f8928f2ee92e
                request-digest: sha256:94e31f44d5a7c87c6ded23c0a5a3fba2824cf938d5eb92bc71e3d8eb979a0317
                room-ref: ./claude-opus-dispatch-not-observed-after-consume/review/validation/briefing-2
mod-block:
---

`assertRecordedGateLifecycle` (internal/ensigncycle/recorded_gate_lifecycle_test.go:121) collapses `!o.dispatch.ordered || !o.dispatch.committed` into one error string ("successor dispatch was not observed after consume"), so two unrelated bugs in how `ordered`/`committed` get computed both surface identically. Both are one-line fixes.

## Bug 1 — unanchored substring match lets a harmless `--help` probe poison the ordering check

`recordedGateLiveObservation` (recorded_gate_lifecycle_test.go:983-986) locates command-log phases with unanchored `strings.Index` over the whole log:

```go
prepareAt := strings.Index(log, "exit=0\tgate prepare ")
bindCommitAt := strings.Index(log, "exit=0\tstate commit recorded-gate-task")
decisionAt := strings.Index(log, "exit=0\tgate record ")
ordered = prepareAt >= 0 && bindCommitAt > prepareAt && decisionAt > bindCommitAt
```

`"exit=0\tgate record "` also matches the read-only line `exit=0\tgate record --help`. In the failing opus run the agent ran `gate prepare --help` then `gate record --help` (both harmless, read-only) before doing any real work. Measured on the actual run artifact (artifact 8835777822): `prepareAt=251` (the `--help` probe), `decisionAt=303` (the `gate record --help` probe, matched instead of the real decision-record command later in the log), `bindCommitAt=1648` -> `decisionAt < bindCommitAt` -> `ordered=false`. A passing rerun's log (artifact 8836227123) has only the `gate prepare --help` probe, no `gate record --help` probe, so `decisionAt` correctly lands on the real command and `ordered=true`. The line-983 `dispatch build` phase-detection already excludes `--help` (confirmed at ~line 1010); this same guard is simply missing for the `prepare`/`record` phases.

Reproduced locally with the real binary (scratch test at the diagnostic's scratchpad, `internal/ensigncycle/zz_repro_test.go`): the same command sequence without help probes grades PASS; inserting a `gate record --help` probe before the real `state commit` flips `ordered` to false and reproduces the exact failure message.

**Fix:** anchor the phase-detection to exclude `--help` invocations, the same way the existing `dispatch build` detection already does (or match complete lines rather than substrings).

## Bug 2 — hardcoded resolution-ID literal misses the room-backed close path

`recordedGateLiveObservation` (recorded_gate_lifecycle_test.go:1018) pickaxes a literal chat-path resolution ID:

```go
closeCommit := ... git log --reverse --format=%H -S"id: resolution:spacedock:recorded-gate-task:validation:1" ...
```

That exact literal is produced only by `chatResolutionID` (internal/gates/operation.go:758-763) -- the `--decision`/`--actor` close path. Codex's live run closed via the **room-backed** path instead: `gate record ... --room ... --consume`, which produces a provider-authored resolution ID (`resolution=resolution:captain-recorded-gate-task-validation-1`) written verbatim from the room's provider `result.json` by `recordRoomLocked` (internal/gates/operation.go ~line 400). The hardcoded pickaxe string never matches that ID, so `close == ""`, and `recordedGateCommittedBeforeDispatch` (line 1193) returns false on its very first guard clause -- independent of Bug 1, and independent of anything about `--consume`/`--stamp` timing (both the collapsed and classic ceremony shapes grade PASS correctly once this is accounted for).

Ironically, `recordedGateLiveObservation` already derives `resolutionID` correctly by regex from the entity's own post-state a few lines later (line 1023) and `assertRecordedGateLifecycle` uses that derived value elsewhere -- only this one commit-lookup pickaxe still hardcodes the literal.

**Fix:** use the regex-derived `resolutionID` (already computed at line 1023) for the pickaxe search instead of the hardcoded chat-path literal, or pickaxe on a room-agnostic substring common to both closing shapes.

## What's ruled out

- Not model nondeterminism -- both bugs are deterministic string-matching gaps triggered by legal, unmodeled agent behavior (a `--help` probe; closing via the room-backed path), not by anything random.
- Not the codex-live dispatch-build-checklist-race issue (different assertion entirely; filed separately as `codex-live-dispatch-build-checklist-race`).
- Not a YAML-nesting-depth regression from PR #599's schema simplification -- measured directly (16 spaces, unchanged).
- Not `-S`'s occurrence-count semantics, not `--stamp` timing/durability, not a real difference in behavior between the collapsed `--consume` and classic 2-step ceremony shapes -- all confirmed fine in the local repro; the two bugs above are the entire explanation for both recurrences.

## Acceptance criteria (rough, for ideation to firm up)

- **AC-1** — A `--help` probe (of `gate prepare`, `gate record`, or `gate consume`) issued at any point in a live agent's command log no longer perturbs `ordered`'s phase detection. Verified by a fixture/unit test reproducing Bug 1's exact repro (insert a `gate record --help` line before the real sequence; assert `ordered` stays true).
- **AC-2** — A room-backed close (`gate record --room ... --consume`) is correctly recognized as committed by `recordedGateCommittedBeforeDispatch`. Verified by a fixture/unit test reproducing Bug 2 (a room-backed close followed by dispatch; assert `committed` is true).
- **AC-3** — The classic chat-path close (`--decision`/`--actor`) continues to work unchanged; no regression on either existing passing shape.

## Out of scope

Any change to the actual gate/dispatch mechanics, `--stamp`, or `--consume` -- confirmed correct. This is a test/grader-only fix.

## Stage Report: implementation

- DONE: Make recorded-gate phase matching line-anchored and ignore exact `--help` probes for prepare, record, and consume, so read-only discovery cannot poison ordered/consumed detection.
  Code commit `48f001276` adds `recordedGatePhaseAt`/`recordedGateHelpProbe` and a focused phase-order regression, plus an end-to-end replay with all three phase help probes.
- DONE: Derive the close pickaxe from the post-state Resolution ID instead of the chat-only literal, allowing provider-authored room-backed Resolution IDs while retaining classic chat close behavior.
  `TestRecordedGateCommittedBeforeDispatchResolutionPaths` exercises both room-backed and classic closes through consumed state and dispatch ancestry.
- DONE: Preserve gate and dispatch mechanics; changes are confined to `internal/ensigncycle/recorded_gate_lifecycle_test.go` and add no product command behavior.
  Evidence: `git show --name-only --format= 48f001276` lists only `internal/ensigncycle/recorded_gate_lifecycle_test.go`; no files under `internal/gates` or dispatch implementation paths changed.
- DONE: Run focused regressions and the complete ensigncycle package suite. Focused rows passed; `go test ./internal/ensigncycle -count=1` passed (105.234s).
  Evidence: `go test ./internal/ensigncycle -race -count=1` also exited 0 (`ok github.com/spacedock-dev/spacedock/internal/ensigncycle`, 123.877s).

### Summary

The recorded-gate grader now ignores harmless help probes when detecting ordered lifecycle phases and recognizes the actual Resolution ID written by either the room-backed or classic close path. Focused regressions cover all three help probes and both close mechanisms; the complete ensigncycle test package is green. Gate and dispatch implementation code remains unchanged.

## Stage Report: validation

- DONE: Reproduce AC-1 across `gate prepare`, `gate record`, and `gate consume` help probes while preserving ordered phase detection.
  `TestRecordedGateLifecyclePhaseDetectionIgnoresHelpProbes` passed with all three exact `--help` lines interleaved in the lifecycle log; it confirms the real prepare/commit/record offsets remain ordered and consume help is not treated as a phase. `TestRecordedGateLifecycleRealCLIReplay` also passed its end-to-end replay after inserting prepare, record, and consume help probes.
- DONE: Reproduce AC-2 for room-backed and classic close paths using the actual post-state Resolution ID.
  `TestRecordedGateCommittedBeforeDispatchResolutionPaths` passed both `classic-chat` and `room-backed` subtests. Each performs prepare, close, consume, a successor commit, and dispatch ancestry inspection; both observations reported `ordered=true` and `committed=true`.
- DONE: Confirm AC-3 and run focused, full, and race ensigncycle tests at candidate `48f001276` atop `7ece33938`.
  Focused command `go test ./internal/ensigncycle -run 'TestRecordedGateLifecyclePhaseDetectionIgnoresHelpProbes|TestRecordedGateCommittedBeforeDispatchResolutionPaths|TestRecordedGateLifecycleRealCLIReplay' -count=1 -v` passed. Full `go test ./internal/ensigncycle -count=1` passed in 279.254s; `go test ./internal/ensigncycle -race -count=1` passed in 282.522s. No files were modified during validation.
- DONE: Perform the required semantic adversarial pass over adjacent lifecycle variants.
  The focused control run passed `TestRecordedGateLifecycleMissingEventControls` (each omitted phase rejected), `TestRecordedGateLifecycleCollapsedConsumeTrace` (fast-path and classic traces both accepted while incomplete close is rejected), terminal/no-successor and entered-stage recovery matrices, plus both Resolution-ID paths. The phase helper is line-scoped and token-based; it ignores only successful lines containing an exact `--help` field, while preserving out-of-order/missing-phase rejection. `git diff --check` is clean and the changed path is test-only, with no new I/O or allocation hot path.

### Summary

Validation passes all three acceptance criteria. Help probes for prepare, record, and consume no longer poison ordered detection; room-backed and classic closes both resolve committed detection through the actual Resolution ID; and focused, full, and race ensigncycle suites are green on the exact candidate.

## Stage Report: implementation (actor correction)

- DONE: Derive the expected approval actor from the first successful close command while retaining `agent:first-officer` as the default for synthetic observations. Room-backed `gate record --room` resolves to `person:captain`; chat `--actor` values are validated and matched exactly at the YAML line boundary.
  Code commit `206a4737c` updates only `internal/ensigncycle/recorded_gate_lifecycle_test.go`; the existing help-probe and Resolution-ID fixes remain intact.
- DONE: Add the reviewed room-backed consumed lifecycle regression and preserve classic actor-swap protection.
  `TestRecordedGateCommittedBeforeDispatchResolutionPaths/room-backed` now executes `gate record --room ... --consume`, then uses the shared full lifecycle assertion (with provider reason optional), exact actor checks, and `gates.Read` durable actor verification. `TestRecordedGateLifecycleProvenanceMutants` retains the classic actor-swap mutant and adds a forged actor-suffix mutant.
- DONE: Run focused, full, race, and exact live Codex evidence for the corrected candidate.
  Focused actor/room/classic rows passed; `go test ./internal/ensigncycle -count=1` passed in 186.391s; `go test ./internal/ensigncycle -race -count=1` passed in 191.068s; exact `go test -tags live ./internal/ensigncycle -run '^TestLiveCodexSharedScenarios$/recorded-gate-lifecycle$' -count=1 -timeout=20m` passed in 250.915s.

### Summary

The grader now derives and exact-matches the expected approval actor from the successful close mode, so room-backed captain approvals no longer fail on a hardcoded First Officer identity. The room-backed `--consume` path is covered through full lifecycle, durable-state, and dispatch ancestry checks; classic actor and forged-actor negatives remain load-bearing. All local, race, and exact live Codex checks passed, with no gate or dispatch mechanics changed.

## Stage Report: implementation (cycle 2)

- DONE: Apply the authorized actor-oracle correction. `206a4737c` derives the expected approval actor from the successful close mode: room-backed close uses `person:captain`, an explicit chat `--actor` is validated and matched exactly, and synthetic observations default to `agent:first-officer`. The existing help-probe exclusions and post-state Resolution-ID pickaxe remain intact.
- DONE: Preserve the room-backed consumed lifecycle as a full assertion and retain actor-swap protection. The room-backed `gate record --room ... --consume` regression uses the shared full lifecycle assertion, including durable state and dispatch ancestry; the classic actor-swap negative remains load-bearing (with forged actor-suffix coverage retained).
- DONE: Preserve test-only scope. `206a4737c` changes only `internal/ensigncycle/recorded_gate_lifecycle_test.go`; gate and dispatch mechanics are unchanged.
- DONE: Evidence is green for the corrected candidate. Focused actor/room/classic regressions passed; `go test ./internal/ensigncycle -count=1` passed in 186.391s; `go test ./internal/ensigncycle -race -count=1` passed in 191.068s; exact `go test -tags live ./internal/ensigncycle -run '^TestLiveCodexSharedScenarios$/recorded-gate-lifecycle$' -count=1 -timeout=20m` passed in 250.915s.
- CHECKLIST OUTCOME: 3 DONE, 0 FAILED, 0 SKIPPED.

## Stage Report: implementation

- DONE: Preserve the corrected actor oracle. Commit `206a4737c` derives `person:captain` for a successful room-backed close, validates and exact-matches an explicit chat `--actor`, and defaults synthetic observations to `agent:first-officer`; the help-probe and post-state Resolution-ID fixes remain intact.
- DONE: Preserve room-backed lifecycle coverage and actor-swap protection. The room-backed `gate record --room ... --consume` regression uses the shared full lifecycle assertion, including durable state and dispatch ancestry; the classic actor-swap negative and forged actor-suffix control remain load-bearing.
- DONE: Preserve test-only scope and the corrected-candidate evidence. The code commit changes only `internal/ensigncycle/recorded_gate_lifecycle_test.go`, with no gate or dispatch mechanics changed. Focused actor/room/classic regressions passed; `go test ./internal/ensigncycle -count=1` passed in 186.391s; `go test ./internal/ensigncycle -race -count=1` passed in 191.068s; exact `go test -tags live ./internal/ensigncycle -run '^TestLiveCodexSharedScenarios$/recorded-gate-lifecycle$' -count=1 -timeout=20m` passed in 250.915s.
- CHECKLIST OUTCOME: 3 DONE, 0 FAILED, 0 SKIPPED.

### Summary

The corrected candidate derives approval identity from the close mode, fully verifies the room-backed consumed lifecycle, and retains actor-provenance negatives. The implementation is test-only, and focused, full, race, and exact live Codex evidence is green.

## Stage Report: implementation

- DONE: Preserve the actor oracle by close mode.
  Evidence: code commit `206a4737c` derives `person:captain` for a successful room-backed close, validates and exact-matches explicit chat `--actor`, defaults synthetic observations to `agent:first-officer`, and retains the help-probe and post-state Resolution-ID fixes.
- DONE: Preserve the room-backed consumed lifecycle and actor-swap controls.
  Evidence: `TestRecordedGateCommittedBeforeDispatchResolutionPaths/room-backed` runs `gate record --room ... --consume` through the shared full lifecycle assertion, including durable state and dispatch ancestry; `TestRecordedGateLifecycleProvenanceMutants` retains classic actor-swap and forged actor-suffix negatives.
- DONE: Preserve test-only scope and the corrected-candidate validation evidence.
  Evidence: `git show --name-only --format= 206a4737c` lists only `internal/ensigncycle/recorded_gate_lifecycle_test.go`; focused actor/room/classic regressions passed; `go test ./internal/ensigncycle -count=1` passed in 186.391s; `go test ./internal/ensigncycle -race -count=1` passed in 191.068s; exact `go test -tags live ./internal/ensigncycle -run '^TestLiveCodexSharedScenarios$/recorded-gate-lifecycle$' -count=1 -timeout=20m` passed in 250.915s.
- CHECKLIST OUTCOME: 3 DONE, 0 FAILED, 0 SKIPPED.

### Summary

The final implementation report preserves the close-mode actor oracle, full room-backed consumed lifecycle proof, actor-swap controls, test-only boundary, and focused/full/race/live evidence.

## Stage Report: validation (fresh candidate 206a4737c)

- DONE: Verify the corrected actor oracle and room-backed close path. `go test ./internal/ensigncycle -run 'TestRecordedGateLifecyclePhaseDetectionIgnoresHelpProbes|TestRecordedGateExpectedActorFromSuccessfulClose|TestRecordedGateCommittedBeforeDispatchResolutionPaths|TestRecordedGateLifecycleRealCLIReplay|TestRecordedGateLifecycleProvenanceMutants|TestRecordedGateLifecycleMissingEventControls|TestRecordedGateLifecycleCollapsedConsumeTrace' -count=1 -v` passed. This covers prepare/record/consume `--help` probes, room-backed `gate record --room ... --consume`, classic explicit `--actor`, synthetic default, exact actor-swap and forged-suffix negatives, missing-event rejection, and collapsed/classic traces.
- DONE: Reproduce the adjacent semantic adversarial matrix. `go test ./internal/ensigncycle -run 'TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume|TestRecordedGateLifecycleTerminalConsumeHasNoDispatchableSuccessor|TestRecordedGateLifecycleEnteredStageRecoveryMatrix|TestRecordedGateLifecycleAC5RefusalMatrix|TestRecordedGateLifecycleAC7ResumeMatrix' -count=1 -v` passed all 5 matrix tests and 33 subtests, including terminal/no-successor, entered-stage recovery, refusal, and resume controls.
- DONE: Run complete package and race evidence at exact candidate `206a4737cc073df243607117c8e19914f62905e6`. `go test ./internal/ensigncycle -count=1` passed (`ok`, 176.585s); `go test ./internal/ensigncycle -race -count=1` passed (`ok`, 197.796s).
- DONE: Run exact live Codex recorded-gate lifecycle evidence. `go test -tags live ./internal/ensigncycle -run '^TestLiveCodexSharedScenarios$/recorded-gate-lifecycle$' -count=1 -timeout=20m` passed (`ok`, 332.517s), with no residual grader failure.
- DONE: Detached audit and scope check. `git diff --check 7ece33938..206a4737c` is clean; candidate parent is `48f001276`; `git diff --name-status` reports only `M internal/ensigncycle/recorded_gate_lifecycle_test.go`; worktree is clean and `gofmt -d internal/ensigncycle/recorded_gate_lifecycle_test.go` is empty. No product gate or dispatch mechanics changed.
- CHECKLIST OUTCOME: 5 DONE, 0 FAILED, 0 SKIPPED.

### Summary

Validation PASSED for `206a4737c`. The actor expectation is derived from the successful close mode (captain for room-backed, exact explicit chat actor, First Officer default for synthetic observations), room-backed consumed lifecycle and classic actor controls pass, all adversarial matrices reject invalid variants, and focused/full/race/exact-live evidence is green. The change remains test-only.
