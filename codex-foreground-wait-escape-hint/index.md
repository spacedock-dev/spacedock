---
id: 5hn35sfb4aenhzjfrr15g9jp
title: Codex foreground wait should tell the captain Esc only returns control
status: backlog
source: "FO dogfood (2026-06-06) - Codex runtime text documents foreground wait semantics, but does not tell the operator before wait_agent that Esc/interrupt is safe and does not mean the worker failed or should be closed."
score: "0.24"
worktree: ""
issue:
---

Codex foreground waits are easy to misread from the operator seat. The runtime
contract says a `wait_agent` timeout is normal and retryable, and that a captain
message or shell-out during the wait is operator activity rather than idle-wake
evidence. It does not require the first officer to say, before entering a
foreground wait, that pressing Esc or interrupting the wait only returns control
to the captain. It is not a worker failure signal and it should not cause the FO
to close or redispatch the worker.

This matters during long validation and live-CI waits: without the hint, the
captain has to infer whether interrupting a wait is safe. The intended behavior
is already clear enough operationally, but the Codex runtime adapter should make
it explicit so future FOs do not treat Esc as a terminal worker event.

## Acceptance criteria

**AC-1 - Codex foreground-wait instructions include the operator hint.**
Verified by a focused integration test over
`skills/first-officer/references/codex-first-officer-runtime.md` that requires
the `## Awaiting Completion` foreground-wait guidance to state that Esc or
operator interruption safely returns control, does not mark the worker failed,
and must not trigger worker closure or redispatch.

**AC-2 - The hint is tied to foreground wait, not every dispatch.**
Verified by the same test ensuring the instruction is located in or directly
under the foreground-wait guidance and does not introduce blanket wait-after-
dispatch wording.

**AC-3 - The probe recipe remains semantically aligned.**
Verified by updating `docs/dev/codex-idle-notification-probe.md` or its
integration test so the foreground-wait comparison explicitly records that an
operator interruption is a non-terminal control return, not a completion or
failure classification.

## Stage test gates

- Ideation should compare this against archived `codex-idle-notification-probe`
  (`gn`) and active `codex-followup-task-reuse` (`82`) to avoid duplicating
  their broader runtime concerns.
- Implementation should update only the Codex runtime adapter/probe text and
  focused integration tests.
- Validation should run the focused integration test plus the relevant
  `skills/integration` package.
