---
title: The Pi worker-lifecycle assert credits a completion surface this host does not provide
status: validation
score: 0.8
source: "Live lane run 37101046846, journey default-headless-gate-stop, 2026-10-03: the journey reported implementation-worker-not-dispatched while the launcher log proved both dispatches succeeded."
id: mk72bnt1b5hsp9sfv83979xs
started: 2026-10-04T04:35:01Z
worktree: .worktrees/spacedock-ensign-repair-pi-worker-lifecycle-observation
---

Two Pi journeys report a false cause. The assert cannot see how the First Officer
observes a worker completion on this host.

## Problem

`assertWorkerLifecycle` credits completion only from a `subagent` tool result that
carries `Run: <id>` and `State: complete`, or, since task `nta`, from a
`subagent_wait` result keyed on the spawned run id.

On this host neither surface appears. Run `37101046846`, journey
`default-headless-gate-stop`, produced these facts:

- The First Officer's root transcript holds 3 `bg_wait` calls and **zero**
  `subagent({action:"status"})` calls.
- The same transcript holds 3 `Background task completed` notices.
- `subagent_wait` appears once, only inside tool-description text. The wait tool is
  named `bg_wait` here.
- `command.log` shows `status --set` and `dispatch build --stamp` succeeding for the
  implementation stage and then the validation stage.
- Two child session files exist under the journey's `sessions/` directory.

The assert therefore computes `completed=-1` and reports
`implementation-worker-not-dispatched`, although the worker was dispatched and its
completion was observed. Task `nta` made the same class of repair once before, when
the assert credited only a `status` result while the First Officer used
`subagent_wait`.

## Value

The lane names the real fault. An operator who reads a red lane learns the actual
cause instead of a false one.

## Out of scope

The two interface faults observed in the same run: the `gate prepare --artifact`
argument, and the split-state read that resolved the stale entity file. Route those
to `vpf multi-artifact-gate-prepare` and to `qx state-ready-cwd-path-resolution-reset`
or `93 document-dispatch-entity-path-base`.

## Expected surface and tolerance

Estimate net LOC change: +40, across 2 files (`internal/ensigncycle/claude_runtime_helpers_test.go`
and a replay test beside it). Insertions ~+55, deletions ~-15. Tolerance: +/-20 net LOC,
+/-1 file.

Declared semantic change: which completion surfaces the assert credits. This task must
NOT change the spawn tolerance, the `completed < validation` ordering, or any other
runtime's grading.

## Acceptance criteria

**AC-1 - The assert credits a native completion notice and a `bg_wait` completion.**
Verified by: a deterministic test that feeds a captured transcript shape containing
those two surfaces and asserts `completed` is set for the spawned run id. Independent
baseline that can move the wrong way: the captured transcript from run
`37101046846` yields `completed=-1` before the change.
Falsifying edit: remove the new credit; the test must turn RED.

**AC-2 - A genuinely missing dispatch still fails.**
Verified by: the existing zero-spawn control. Falsifying edit: credit any transcript
regardless of spawn count; the control must turn RED.

**AC-3 - The ordering contract is preserved.**
Verified by: the existing `completed < validation` assertion. Falsifying edit: drop
the ordering check; the inverted-ordering control must turn RED.

**AC-4 (no-regression) -** `go test ./internal/ensigncycle/...` is green, and no other
runtime's grading changes.

## Test plan

Primary proof owner: `internal/ensigncycle/claude_runtime_helpers_test.go` (the shared
assert) plus one replay test. The captured root transcript is preserved under
`docs/dev/.spacedock-state/_evidence/pi-delegated-gate-continuation-reliability/retained-pi-recorded-gate/`
and must be copied to repository testdata before the replay test depends on it. Deterministic
tests only.

## Stage Report: implementation

- DONE: Fix the shared assert so it credits the completion surfaces this host supplies: a bg_wait result keyed on the spawned run id, and the native completion notice. Today it credits only a subagent result carrying Run/State, or a subagent_wait result, so it reports spawns=1 completed=-1 while the officer did dispatch and did observe the completion. Keep the rule that a genuinely missing dispatch still fails, and keep the completed-before-validation ordering. Add the deterministic test from the captured transcript shape, and keep the negative control that removes the new credit so the test turns red.
  Commit `2d5bd25f6`: `assertWorkerLifecycle` now credits a `bg_wait` result that names the spawned run id and reports `done`/`complete` (alongside `subagent_wait`), and a top-level `custom_message` `customType: subagent-notify` whose body is `Background task completed`. Both stay behind the existing spawn-count and `completed < validation` checks. Positive control `TestPiWorkerLifecycleObservationReplay` feeds the captured shape; removing either new credit makes it RED with `validation lifecycle incomplete: spawns=1 completed=-1 validation=5 report=<nil>` — the same completed=-1 run 37101046846 reported.
- DONE: Prove it: the focused deterministic tests, then a local live run of the affected journey auto-continue-after-implementation using SPACEDOCK_LIVE_RUNTIME=pi and a distinct SPACEDOCK_LIVE_ARTIFACT_DIR. Record the exact result. Do not start a CI lane run.
  `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/live-observer SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 45m -run '^TestLiveCommonAutoContinueAfterImplementation$' ./internal/ensigncycle -v` -> `--- PASS: TestLiveCommonAutoContinueAfterImplementation (766.12s)`. Both fixtures passed: single-root `error=<nil> timeout=false` 7m40.9s, split-root `error=<nil> timeout=false` 5m03.6s, model `openai-codex/gpt-5.6-luna:max`. Grade `pass`, observed codes `[]`; no XFAIL/XPASS line, because the binding is gone. No CI lane run started.
- DONE: When that journey passes locally with an exact passing pair, clear its XFAIL binding in the shared runner and record the evidence in the entity. Do not leave the binding as XPASS.
  Commit `a2166585a`: `TestLiveCommonAutoContinueAfterImplementation` gaps are `nil` (binding to `mk72bnt1b5hsp9sfv83979xs` removed with its stale evidence-defect comment). Rebased onto the peer tip `fb4428e9c`, where the binding lived; original base was main `4436ec14c`. Peer branch untouched.
- DONE: Keep a genuinely missing dispatch failing, and the completed-before-validation ordering.
  Same test includes the zero-spawn control (spawn retargeted to implementation -> spawns=0, RED) and the inverted-order control (completion moved after `gate prepare` -> `completed >= validation`, RED). Each turns RED only while its guard exists.
- DONE: no-regression `go test ./internal/ensigncycle/...` and no other runtime's grading changes.
  `go test ./internal/ensigncycle/... -count=1` green (286.7s); `go test -race ./internal/ensigncycle/ -run TestPiWorkerLifecycleObservationReplay` green; `go test ./internal/contractlint/` green; gofmt clean. The claude/codex branches are untouched.

### Summary

The shared worker-lifecycle assert now credits the two completion surfaces a Pi host supplies — a `bg_wait` result keyed on the spawned run id and a native `subagent-notify` completion notice — while preserving the spawn-count and `completed < validation` guards. A deterministic replay of the captured transcript shape plus its removed-credit, zero-spawn, and inverted-order controls proves each guard is load-bearing. After rebasing onto the peer stack tip, the auto-continue Pi XFAIL was cleared, and the journey passed locally with an exact two-variant passing pair (766.12s) and no XPASS.
