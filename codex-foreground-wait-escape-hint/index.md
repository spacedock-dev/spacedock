---
id: 5hn35sfb4aenhzjfrr15g9jp
title: Codex foreground wait should tell the captain Esc only returns control
status: implementation
source: "FO dogfood (2026-06-06) - Codex runtime text documents foreground wait semantics, but does not tell the operator before wait_agent that Esc/interrupt is safe and does not mean the worker failed or should be closed."
score: "0.24"
worktree: .worktrees/spacedock-ensign-codex-foreground-wait-escape-hint
issue:
started: 2026-06-06T04:16:46Z
---

Codex foreground waits are easy to misread from the operator seat. The runtime
contract says a `wait_agent` timeout is normal and retryable, and that a captain
message or shell-out during the wait is operator activity rather than idle-wake
evidence. It does not require the first officer to say, before entering a
foreground wait, that pressing Esc or otherwise interrupting the wait only
returns control to the captain. The interruption is not a worker failure signal
and must not cause the FO to close, redispatch, or mark the worker failed.

This matters during long validation and live-CI waits: without the hint, the
captain has to infer whether interrupting a wait is safe. The intended worker
lifecycle behavior is already clear enough operationally; this task only makes
the operator-facing Codex foreground-wait guidance explicit.

## Scope decisions

This is a narrow Codex foreground-wait wording task, not a broader runtime
semantics task.

- Archived `codex-idle-notification-probe` (`gn`) already owns the distinction
  among explicit foreground wait, queued-notification flush, autonomous idle
  wake-up, and no-notification observations. This task should preserve that
  taxonomy and only add that an operator interruption during the explicit
  foreground-wait comparison is a non-terminal return of control.
- Active `codex-followup-task-reuse` (`82`) owns reusable-worker continuation
  semantics for `followup_task`, stale completions, and validation re-review
  handle correlation. This task should not change reuse routing, `send_input`,
  `send_message`, `followup_task`, or feedback-rejection behavior.
- The hint belongs in Codex first-officer foreground-wait guidance, not in
  dispatch setup, every post-dispatch path, worker completion handling, or
  terminal teardown. The FO should still foreground-wait only after ready
  completions, gate decisions, state transitions, and newly dispatchable work
  are exhausted.

## Proposed approach

1. Add a focused integration test before changing the runtime text. The test can
   live beside the existing Codex idle-notification tests, for example
   `TestCodexForegroundWaitEscapeHint`, and should parse the `### Foreground
   wait` subsection of
   `skills/first-officer/references/codex-first-officer-runtime.md`.
2. Make the test require a narrow operator hint in that foreground-wait
   subsection:
   - the FO tells the captain before calling `wait_agent(handle)` that Esc or an
     operator interruption safely returns control;
   - the hint states the worker is not failed, closed, or redispatched by that
     interruption;
   - the retryable same-handle foreground-wait semantics remain intact.
3. Keep the existing blanket-wait guardrails. The implementation should not add
   wording that implies Codex must foreground-wait after every dispatch or that
   an Esc/interruption hint is needed when the FO ends the turn and relies on a
   mailbox notification.
4. Update `docs/dev/codex-idle-notification-probe.md` or its existing recipe
   test so `Foreground wait comparison` records an operator interruption as a
   non-terminal return of control. It should remain distinct from the
   `foreground_wait`, `queued_flush`, `autonomous_idle_wake`, and
   `no_notification_observed` classifications.
5. No spike is needed. The work composes already-proven mechanisms: the existing
   Codex `### Foreground wait` adapter section, the existing idle-notification
   recipe, and the existing `skills/integration` test helpers. No new Codex host
   primitive or live mailbox behavior is being asserted.

## Acceptance criteria

**AC-1 - Codex foreground-wait guidance includes the operator hint.**
Verified by `go test ./skills/integration -run TestCodexForegroundWaitEscapeHint`
parsing the `### Foreground wait` subsection of
`skills/first-officer/references/codex-first-officer-runtime.md` and failing
unless that subsection tells the captain that Esc or operator interruption
safely returns control and does not mark the worker failed, close it, or
redispatch it.

**AC-2 - The hint is tied to foreground wait, not every dispatch.**
Verified by the same focused test plus the existing blanket-wait guardrail in
`TestCodexIdleNotificationRuntimeContract`; the tests fail if the hint appears
as blanket dispatch guidance or introduces wording such as "wait after every
dispatch".

**AC-3 - The probe recipe remains semantically aligned.**
Verified by `go test ./skills/integration -run TestCodexIdleNotificationRecipeShape`
or a focused companion test reading `docs/dev/codex-idle-notification-probe.md`;
the test fails unless `Foreground wait comparison` records an operator
interruption as a non-terminal control return rather than a worker completion,
failure, closure, redispatch, or idle-wake classification.

**AC-4 - Existing Codex runtime and probe coverage stays green.**
Verified by `go test ./skills/integration -run 'TestCodexIdleNotification|TestCodexForegroundWaitEscapeHint'`
passing after the focused implementation.

## Test plan

- First run the new focused test and confirm it fails on the missing foreground
  wait hint. Cost: low; fixture-backed integration test only.
- Update only the Codex first-officer runtime adapter and the idle-notification
  probe recipe text needed for the hint. Cost: low; no code-path changes.
- Run `go test ./skills/integration -run 'TestCodexIdleNotification|TestCodexForegroundWaitEscapeHint'`
  to verify the hint, existing taxonomy, blanket-wait guardrails, and recipe
  shape together. Cost: low.
- Run the repo gates required by this project before validation claims:
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
  Cost: medium, dominated by the race test.

## Stage test gates

- Ideation should compare this against archived `codex-idle-notification-probe`
  (`gn`) and active `codex-followup-task-reuse` (`82`) to avoid duplicating
  their broader runtime concerns.
- Implementation should update only the Codex runtime adapter/probe text and
  focused integration tests.
- Validation should run the focused integration test plus the relevant
  `skills/integration` package.

## Stage Report: ideation

- DONE: compare the foreground-wait hint task against archived codex-idle-notification-probe and active codex-followup-task-reuse so the scope does not duplicate broader runtime concerns
  Evidence: `Scope decisions` limits the task to Codex foreground-wait operator wording, preserving `gn`'s idle-notification taxonomy and leaving `82`'s follow-up reuse semantics untouched.
- DONE: produce a precise implementation plan and test plan for adding the operator Esc/interruption hint only to Codex foreground-wait guidance, not blanket dispatch/wait behavior
  Evidence: `Proposed approach`, AC-1 through AC-4, and `Test plan` specify the adapter subsection, probe recipe update, focused tests, and blanket-wait guardrail.
- DONE: write an ideation stage report with the evidence and any scope decisions needed for implementation
  Evidence: this stage report records the comparison, implementation scope, and verification plan in the state entity.

### Summary

The ideation pass narrows the work to a Codex foreground-wait operator hint: before `wait_agent`, tell the captain Esc/interruption only returns control and does not fail, close, or redispatch the worker. It explicitly avoids changing the archived idle-wake taxonomy or the active follow-up reuse task, and gives implementation a focused test path plus the required repo gates.

## Stage Report: implementation

- DONE: Add the Codex foreground-wait operator hint in the source runtime guidance: before `wait_agent(handle)`, the FO tells the captain Esc/operator interruption only returns control and does not fail, close, or redispatch the worker; preserve same-handle retry semantics and the idle-notification taxonomy.
  Evidence: code commit `f305f1c5` updates `skills/first-officer/references/codex-first-officer-runtime.md` and preserves same-handle wording.
- DONE: Keep proof 4q-compliant: if a test reads first-officer runtime instruction text, put it in `internal/contractlint` as a quarantined structural check; do not add behavior-proof greps over `SKILL.md`/runtime prose in `skills/integration`.
  Evidence: code commit `f305f1c5` adds `internal/contractlint/codex_foreground_wait_shape_test.go`; `go test ./internal/contractlint -count=1` passed.
- DONE: Update the idle-notification probe recipe if needed so operator interruption is a non-terminal foreground-wait return, then run focused tests plus `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
  Evidence: code commit `f305f1c5` updates `docs/dev/codex-idle-notification-probe.md`; focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed.

### Summary

Implemented the Codex foreground-wait operator hint in the source runtime adapter and aligned the idle-notification probe recipe so Esc/operator interruption is a non-terminal return of control with same-handle retry. Added quarantined contractlint structural coverage for the runtime section and probe recipe, then ran the required focused and full repo gates.
