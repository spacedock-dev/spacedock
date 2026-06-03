---
id: gnanx8t260dyax3x6s841bgc
title: Make the Codex no-wait subagent completion behavior retryable
status: ideation
source: "captain (2026-06-03) — dogfood session found Codex may queue subagent completion notifications without explicit wait_agent; document a repeatable idle-wake probe in case Codex changes behavior"
started: 2026-06-03T16:09:19Z
completed:
verdict:
score: "0.26"
worktree:
issue:
---

Codex first-officer operation depends on knowing when `wait_agent` is required and when a worker can safely run in the background. The current adapter says the Codex completion signal is the async final-status notification in the FO mailbox and `wait_agent` is only an optional accelerator. A 2026-06-03 dogfood session confirmed a useful but subtle behavior: a no-wait probe worker completed and its notification was queued before the captain's next message. It did not fully prove that the FO is automatically re-entered while the user stays idle, because queued notification delivery and actual model wake-up are observably different.

## Problem

The current evidence for Codex no-wait completion is anecdotal and easy to misread. If Codex changes how interactive sessions wake on background subagent completion, the FO runtime contract could become either too blocking or too optimistic:

- too blocking: the FO foreground-waits after every dispatch even when Codex could wake from the mailbox, slowing interactive workflow operation and making steering awkward;
- too optimistic: the FO assumes background completion will wake it, leaves critical-path workers unwatched, and does not continue dispatching until the captain sends another message.

The team needs a documented recipe that separates three cases:

1. `wait_agent` foreground wait returns or is interrupted.
2. No `wait_agent`, but later user/tool activity flushes a queued completion notification.
3. No `wait_agent` and no user/tool activity; Codex independently wakes the FO from the subagent notification.

## Proposed approach

Make the Codex runtime contract conservative about state transitions but explicit
about the three observable completion paths. The future implementation should
change instruction text and add small, checked artifacts rather than building a
live probe runner before the host exposes timestamped mailbox events.

1. Update `skills/first-officer/references/codex-first-officer-runtime.md`.
   Keep the existing mailbox model, but split `## Awaiting Completion` into
   named outcomes:

   - **Foreground wait:** the FO calls `wait_agent(handle)` because the next
     stage transition, gate result, or dispatch decision depends on the worker's
     final status. A timeout is normal and retryable with the same handle. A
     captain message or shell-out during the wait is operator activity, not idle
     wake evidence.
   - **Queued notification flushed by later activity:** the FO does not call
     `wait_agent`, ends the turn, and a later captain message, tool action, or
     shell-out causes Codex to deliver a worker final-status notification that
     had already been queued. This proves mailbox ordering, not autonomous FO
     wake-up.
   - **Autonomous idle FO wake-up:** the FO does not call `wait_agent`, performs
     no later activity, and Codex starts a new assistant turn from the worker
     final-status notification alone. This is the only observation that proves
     no-wait idle wake-up.

   The contract should say: non-critical background work may end the FO turn and
   rely on the mailbox, but a critical-path result still requires `wait_agent` or
   another explicit operator-visible waiting posture. The guidance must not say
   every Codex dispatch must immediately foreground-wait.

2. Add a rerunnable manual recipe at
   `docs/dev/codex-idle-notification-probe.md`. The recipe should have three
   sections named `Foreground wait comparison`, `No-wait idle probe`, and
   `Queued-notification flush check`. It should include:

   - the exact no-write worker prompt;
   - the minimum idle window, currently 90 seconds after the FO turn ends;
   - instructions to avoid captain messages, shell-outs, terminal jobs, and tool
     calls during the idle window;
   - a comparison run that calls `wait_agent(handle)`, records timeout or final
     status, then retries the same handle after a non-stopping interruption;
   - interpretation rules that classify the run as `foreground_wait`,
     `queued_flush`, `autonomous_idle_wake`, or `no_notification_observed`.

3. Add timestamped evidence under
   `docs/dev/_evidence/codex-idle-notification-probe/`. Use JSON so Go tests can
   validate it with the standard library. The first file can capture the
   2026-06-03 dogfood run as
   `docs/dev/_evidence/codex-idle-notification-probe/2026-06-03-dogfood.json`.
   A run record should include these fields:

   - `schema_version`: integer, starting at `1`;
   - `run_id`: stable short name, such as `2026-06-03-dogfood`;
   - `host`: `codex`;
   - `codex_cli_version`: exact output if available, otherwise `unknown`;
   - `session_started_at_utc`: RFC3339 timestamp or `null`;
   - `worker_handle`: Codex worker handle or task id;
   - `worker_prompt`: the exact probe prompt;
   - `worker_delay_seconds`: intended worker delay;
   - `spawned_at_utc`: RFC3339 timestamp or `null`;
   - `fo_turn_ended_at_utc`: RFC3339 timestamp or `null`;
   - `idle_window_seconds`: intended idle observation window;
   - `first_user_activity_at_utc`: RFC3339 timestamp or `null`;
   - `final_status_notification_at_utc`: RFC3339 timestamp or `null`;
   - `notification_delivered_before_user_message`: boolean or `null`;
   - `autonomous_wake_observed`: boolean;
   - `classification`: one of `foreground_wait`, `queued_flush`,
     `autonomous_idle_wake`, or `no_notification_observed`;
   - `interpretation`: one short sentence explaining what the run proves.

4. Add focused tests in `skills/integration/`, for example
   `codex_idle_notification_test.go`.

   - A runtime-text test should parse the `## Awaiting Completion` section of
     `skills/first-officer/references/codex-first-officer-runtime.md` and assert
     that the three named outcomes are present, that the critical-path
     `wait_agent` rule is present, and that blanket "wait after every dispatch"
     wording is absent.
   - A recipe-shape test should read
     `docs/dev/codex-idle-notification-probe.md` and assert that the three probe
     sections exist, that the idle window is at least 90 seconds, and that queued
     delivery is not classified as autonomous wake-up.
   - An evidence-schema test should load every
     `docs/dev/_evidence/codex-idle-notification-probe/*.json`, validate the
     required fields, parse non-null timestamp fields as RFC3339, and reject
     classifications outside the enum.

## Spike / risk order

The riskiest mechanism is Codex host behavior, not Markdown wording. Because the
current multi-agent surface does not expose a live, timestamped event stream that
a Go test can drive, the implementation should first add the manual recipe and
evidence schema, then capture one current run before relaxing or tightening the
runtime guidance. The existing dogfood note is only queued-delivery evidence; it
must not be treated as proof of autonomous idle wake-up.

## Acceptance criteria

Each AC names a property of the finished task and a check outside this task body that can fail.

**AC-1 - The Codex FO runtime contract distinguishes all three completion paths.**
Verified by: `go test ./skills/integration -run TestCodexIdleNotificationRuntimeContract` parses the `## Awaiting Completion` section of `skills/first-officer/references/codex-first-officer-runtime.md` and fails unless it names foreground `wait_agent`, queued notification flushed by later activity, and autonomous idle FO wake-up as separate outcomes.

**AC-2 - Critical-path guidance remains conservative without forcing foreground waits for every dispatch.**
Verified by: the same runtime-contract test fails unless the Codex adapter says a critical-path worker result requires `wait_agent` or explicit operator-visible waiting, and also fails on blanket wording such as "must call `wait_agent` after every dispatch".

**AC-3 - A rerunnable manual probe documents the foreground-wait comparison, no-wait idle probe, and queued-notification flush check.**
Verified by: `go test ./skills/integration -run TestCodexIdleNotificationRecipeShape` reads `docs/dev/codex-idle-notification-probe.md` and fails unless the three sections, minimum idle window, no-activity rule, and interpretation rules are present.

**AC-4 - Timestamped evidence captures the current Codex behavior in a machine-checkable format.**
Verified by: `go test ./skills/integration -run TestCodexIdleNotificationEvidenceSchema` parses `docs/dev/_evidence/codex-idle-notification-probe/*.json`, validates required fields and RFC3339 timestamps, and rejects unknown classifications.

**AC-5 - Existing queued-notification evidence is not promoted to autonomous idle-wake evidence.**
Verified by: the evidence-schema test fails if the 2026-06-03 dogfood evidence classifies the `Dalton` run as `autonomous_idle_wake` while `first_user_activity_at_utc` is non-null or `autonomous_wake_observed` is false.

## Test plan

- Add `skills/integration/codex_idle_notification_test.go` with three focused tests: runtime contract, recipe shape, and evidence schema. Cost: low.
- Run `go test ./skills/integration -run TestCodexIdleNotification` first to see the new tests fail for missing wording/artifacts, then implement the runtime text, recipe, and JSON evidence. Cost: low.
- Run the repo gates after the focused tests pass: `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. Cost: medium, dominated by the race test.
- Defer a live automated Codex idle-wake probe until the host exposes timestamped mailbox events or a stable automation API. Cost if added later: medium to high.

## Notes

- 2026-06-03 dogfood observation: `Dalton` (`019e8e34-ec78-7691-b56d-14528a0a9f71`) was spawned without `wait_agent` and returned `Done: idle-wake Codex notification probe completed.` The notification appeared in the delivered batch before the captain asked "and how long should we wait?", but that only proves queued notification ordering, not necessarily autonomous FO re-entry while the user remains idle.
- The same session found a foreground wait can be interrupted by captain input, and a shell-out like `! echo ...` is a steering path while waiting. That path should be documented as operator interaction, not mistaken for an idle-wake proof.

## Stage Report: ideation

- DONE: Separate the documented outcomes: foreground `wait_agent`, queued notification flushed by later activity, and autonomous idle FO wake-up.
  Evidence: `Proposed approach` now defines the three outcomes and binds them to runtime contract updates.
- DONE: Turn the manual recipe into an implementation-ready design with artifact/evidence paths, timestamp fields, and interpretation rules.
  Evidence: the task names `docs/dev/codex-idle-notification-probe.md`, JSON evidence under `docs/dev/_evidence/codex-idle-notification-probe/`, required fields, and classification rules.
- DONE: Tighten acceptance criteria so future Codex behavior drift is caught without requiring every dispatch to foreground-wait.
  Evidence: AC-1 through AC-5 define runtime, recipe, and evidence-schema tests, including an explicit ban on blanket `wait_agent` wording.

### Summary

Fleshed out the seed into an implementation-ready task body that separates foreground waiting, queued notification delivery, and autonomous idle wake-up. The design keeps Codex runtime guidance conservative for critical-path results while preserving background dispatch compatibility, and it gives future workers concrete tests and artifact paths to implement.
