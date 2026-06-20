# Codex Idle Notification Probe

Use this recipe when checking whether Codex can resume a first-officer turn from
a background worker final-status notification without an explicit foreground
wait. The probe separates foreground waiting, queued notification delivery, and
autonomous idle wake-up.

Use this exact no-write worker prompt:

```text
You are a no-write Codex idle-wake probe. Do not read or write files, do not run tools, and do not modify state. Sleep for 30 seconds, then reply exactly: Done: idle-wake Codex notification probe completed.
```

Record `worker_delay_seconds` as `30`. The minimum idle window is 90 seconds
after the FO turn ends.

During a no-wait idle window, avoid captain messages, shell-outs, terminal jobs,
and tool calls. Any such activity can flush a queued mailbox notification and
must be treated as operator activity, not idle wake evidence.

## Foreground wait comparison

1. Dispatch a worker with the exact no-write prompt and record its handle.
2. When the live Codex tool surface exposes a foreground wait binding, use `«completion-signal»` foreground waiting only when there is no ready workflow work.
3. Record whether the call returns a timeout or a final status.
4. If captain input, Esc, or another operator interruption returns control,
   record it as a non-terminal foreground-wait return; do not classify it as a
   worker final status, failure, closure, redispatch, or idle wake-up. If the
   worker remains unresolved, reinstall the global wait when waiting is again
   the next useful idle action and record the next timeout or final status.
5. Classify a final status observed through this explicit wait path as
   `foreground_wait`.

Use `«roster-reconcile»` only to inspect active/completed task paths for attribution and debugging. Durable workflow state remains authoritative.

## No-wait idle probe

1. Dispatch a worker with the exact no-write prompt and record its handle.
2. Do not use foreground wait; end the FO turn.
3. Keep the session idle for the minimum idle window of 90 seconds after the FO
   turn ends.
4. During that idle window, avoid captain messages, shell-outs, terminal jobs,
   and tool calls.
5. If Codex starts a new assistant turn from the worker final-status
   notification alone, classify the run as `autonomous_idle_wake`.
6. If no notification appears during the idle window, classify the idle portion
   as `no_notification_observed` unless later activity reveals a queued
   notification.

## Queued-notification flush check

1. After a no-wait idle window without autonomous wake-up, perform one explicit
   activity: a captain message, tool action, or shell-out.
2. If Codex then delivers the worker final-status notification, record that the
   notification was delivered after later activity.
3. A queued notification flushed by later activity is `queued_flush`, not
   `autonomous_idle_wake`.

## Interpretation Rules

- `foreground_wait`: foreground wait returned
  a timeout or final-status mailbox update, and any operator interruption that
  returned control was recorded as a non-terminal foreground-wait return
  followed by reinstalling the global wait only when waiting was again the next
  useful idle action. Legacy fixtures may describe handle-scoped waiting, but
  they must be explicitly marked as legacy evidence.
- `queued_flush`: no foreground wait was used, but later captain, tool, or
  shell-out activity caused a queued worker final-status notification to appear.
  This remains queued/activity-driven delivery unless an autonomous wake probe
  proves otherwise.
- `autonomous_idle_wake`: no foreground wait and no later activity occurred, and
  Codex began a new assistant turn from the worker final-status notification
  alone.
- `no_notification_observed`: no foreground wait was used and no final-status
  notification appeared during the minimum idle window.

Store run records as JSON under
`docs/dev/_evidence/codex-idle-notification-probe/`. Use RFC3339 UTC timestamps
when known, and `null` when the observation did not capture an exact timestamp.

Codex shutdown probing is separate from this idle-notification recipe.
`«worker.shutdown»` remains unresolved until a live or fixture-backed
probe proves whether live interruption terminates, pauses, or leaves a worker
addressable. Do not bless live interruption as a shutdown binding from idle-wake
evidence alone.
