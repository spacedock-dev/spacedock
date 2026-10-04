---
title: The Pi worker-lifecycle assert credits a completion surface this host does not provide
status: implementation
score: 0.8
source: "Live lane run 37101046846, journey default-headless-gate-stop, 2026-10-03: the journey reported implementation-worker-not-dispatched while the launcher log proved both dispatches succeeded."
id: mk72bnt1b5hsp9sfv83979xs
started: 2026-10-04T04:35:01Z
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
