---
title: Make Codex implementation dispatch mechanically acknowledged
status: ideation
source: "98a cycle-4 live evidence: identical spawn envelopes produced two native-spawn passes and one empty-wait failure"
started: 2026-08-10T09:23:06Z
completed:
verdict:
score: 0.95
worktree:
issue:
milestone: 0.27.0
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: n28423efmj358m5av61z2fxx
gates:
    version: 1
    records:
        - id: gate:n28423efmj358m5av61z2fxx:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:n28423efmj358m5av61z2fxx-backlog-1
              briefing:
                id: briefing:n28423efmj358m5av61z2fxx:backlog:attempt-1:revision-1
                digest: sha256:2d7ed1f6f562e4b034a06241033654f1ec8f4c7c2b24888d096387f1c4b6a782
                request-digest: sha256:5886795dd5d37b928946dd83b3364574ea2f37e56ddcb7a1cf77d79a1b0bd257
                room-ref: ./mechanically-acknowledge-codex-implementation-dispatch/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:n28423efmj358m5av61z2fxx:backlog:1
                briefing: briefing:n28423efmj358m5av61z2fxx:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T09:22:46.146101Z"
                decision: approve
                reason: Captain explicitly approved creating and shaping this smallest end-user mechanical prerequisite after 98a cycle-4 nondeterminism.
              application:
                target-stage: ideation
                state: consumed
---

## Problem

Claude or Codex can receive a valid fresh dispatch envelope without starting a
native worker. The First Officer can then move the entity to the next stage.

The current binary records stage entry, but it does not attest the native worker.
A report read, empty wait, or caller statement can look like completion.

The captain expanded this task to one host-neutral mechanism for Claude and
Codex. Pi keeps its current behavior.

## End-user value

A user gets one native worker for each fresh Claude or Codex envelope. The host
attests that worker before any later stage can start.

The same receipt also binds native completion to the exact entity, stage, and
dispatch epoch. Caller text cannot complete the receipt.

## Constraints

- The launcher and binary own the acknowledgment boundary.
- Do not add another prose guard.
- Do not deliver a component-only receipt API.
- A successful fresh envelope must create pending state.
- A second build for the same entity and stage must refuse pending state.
- Every stage mutation door must use one shared guard.
- Preserve gate authority and worker report ownership.
- Preserve Pi behavior and its current XFAIL.
- Remove an exact live XFAIL only after passing product evidence.

## Staff findings and disposition

| Finding | Materiality | Owner | Disposition and evidence |
|---|---|---|---|
| Caller acknowledgments are forgeable. | Material | This task | Accepted. The old command trusted caller text. The new design removes that command. |
| The grammar has no unique epoch. | Material | This task | Accepted. The new envelope has a random epoch token in a native field. |
| Pending receipts have no recovery authority. | Material | This task | Accepted. A confirmed operator can abandon one exact epoch. |
| Receipt writes race report writes. | Material | This task | Accepted. The spike used separate immutable files during a worker report commit. |
| A Codex build without `--stamp` bypasses receipts. | Material | This task | Accepted. Every fresh Claude or Codex build creates pending state. |
| Mutation doors, schema, help, and tests are absent. | Material | This task | Accepted. The exact surface below includes each item. |

## Current boundaries

The current path has five boundaries:

1. `dispatch build` emits a host envelope.
2. The First Officer calls the native spawn tool.
3. The host creates a worker identity.
4. The host emits a completion event.
5. A Spacedock command changes the entity stage.

The binary owns boundaries 1 and 5. The host owns boundaries 2 through 4.
No durable value currently joins these boundaries.

## Authorized second spike

The spike used `/tmp/n284-second-spike.9E5G54`. It contained a disposable
split-root state checkout and a throwaway observer.

The spike did not edit product files. It exercised native hosts and disposable
state only.

### Codex evidence

The Codex rollout emitted a `SubAgentActivity` start event 107 milliseconds
after the native spawn call.

The event contained thread ID
`019feb0b-e7a9-7d72-9bf5-6609043a8a25`. It also contained the canonical agent
path with epoch token `n284_e1_20260810a`.

The final mailbox event arrived 5.040 seconds after that start. Its sender path
matched the start event and mapped to the same thread ID.

The observer moved one receipt through `pending`, `spawned`, and `completed`.
An old completion did not complete a replacement epoch.

### Claude evidence

The live probe used Claude Code 2.1.226. The stream emitted `task_started` with
task ID `a6253614e277e84a5`.

The event also contained tool-use ID `toolu_01CK3Wnms2YJdZTjKWHrcBB2` and epoch
token `n284_claude_ce1`.

The later `task_notification` contained both IDs and status `completed`. The
durable Claude session transcript contained the same structured values.

The observer moved the Claude receipt from `pending` to `spawned`, then to
`completed`. A second same-stage build returned exit code 20 while pending.

### Recovery and concurrent-write evidence

The observer recorded a spawned Codex worker, then accepted an operator
abandonment with an operator name and reason.

Abandonment did not permit stage advance. It permitted one replacement epoch for
the same stage.

The observer wrote a receipt sidecar 160 times while a worker committed the
entity report. Commit `cb081be` changed only `entity.md`.

Observer commit `57996dc` then changed only the receipt files. The report and
receipt remained intact.

### Mutation-door evidence

The pending receipt refused these six prototype doors:

- direct `status --set` stage change;
- `gate consume` stage change;
- terminal status change;
- merge or archive finalization;
- `dispatch build --advance`;
- fresh Claude or Codex build without receipt creation.

The first five doors passed only after a matching completion. The no-stamp fresh
build always created pending state or refused.

## Selected approach

Add one launcher-owned envelope observer for Claude and Codex. The launcher
already stays resident as the host parent.

The observer reads the host's durable structured event stream. It does not trust
First Officer text or a caller acknowledgment command.

Use this sequence:

1. A fresh `dispatch build` creates a random epoch.
2. The binary adds the epoch token to a native envelope field.
3. The binary writes one immutable `pending.json` receipt event.
4. The launcher matches one host start event to the epoch token.
5. The launcher writes one immutable `spawned.json` event.
6. The launcher matches completion to the native worker and call IDs.
7. The launcher writes one immutable `completed.json` event.
8. Every stage mutation command requires the completed event.

For Claude, the token is in the `description` and worker `name`. For Codex, the
sanitized `task_name` contains the same token.

The build reserves suffix space before it caps the native name. Therefore, the
token cannot exceed either host's name limit.

The observer stores a source cursor after each event. On resume, it replays the
host stream from that cursor.

The launcher binds the observer to the child host session and its start time.
The observer ignores event files from other sessions.

Replay is idempotent. A completion must match the epoch, native call ID, and
native worker ID.

The observer fails closed on an unknown event schema. It never infers completion
from a report, roster, empty wait, or parent message.

This is the smallest host-attested end-to-end mechanism found. Both hosts already
persist the required identities while the launcher remains resident.

## Receipt format

Each active entity uses this directory:

```text
<entity-companion>/dispatch-receipts/<epoch>/
  pending.json
  spawned.json
  completed.json | abandoned.json
```

The entity companion is `<slug>/` for both flat and folder entities. The worker
continues to write the entity Markdown file.

Each event uses receipt schema version 1. Required common fields are:

```json
{
  "schema_version": 1,
  "event": "pending|spawned|completed|abandoned|superseded",
  "entity_id": "n28423efmj358m5av61z2fxx",
  "entity_path": "mechanically-acknowledge-codex-implementation-dispatch.md",
  "stage": "implementation",
  "host": "claude|codex",
  "epoch": "random-128-bit-token",
  "session_nonce": "launcher-session-token",
  "observed_at": "RFC3339 timestamp"
}
```

`spawned.json` also requires `native_call_id`, `native_worker_id`, and
`source_cursor`.

`completed.json` requires the same native IDs and a later `source_cursor`.
`abandoned.json` requires operator identity, reason, and confirmation time.

Receipt event files are immutable. The store uses a per-epoch lock and atomic
create, then commits only the new event path.

The state sync retries Git index-lock contention. A sync failure leaves the local
event durable and blocks stage advance.

## Replay and build rules

A fresh Claude or Codex build always creates pending state. `--stamp` still
controls the existing stage-entry stamp only.

The build locks the entity and stage, commits pending state, and then emits the
envelope. A pending-state failure prevents envelope output.

An output failure can leave pending state. Recovery then follows the normal
operator abandonment path.

Plain builds without `--stamp` cannot bypass receipt creation. `--validate-only`
and `--print-schema` do not create envelopes or receipts.

A build refuses when the same entity and stage has a pending, spawned, or
completed epoch. It emits no second envelope.

An abandoned epoch permits one replacement. The replacement pending event names
the old epoch in `supersedes`.

The store writes `superseded.json` for the old epoch in the same locked
transaction. A replayed completion for the old epoch has no authority.

Pi builds do not create these receipts. Existing entities without a Claude or
Codex fresh receipt keep their current behavior.

## Recovery authority

Only the host observer can write `spawned.json` or `completed.json`. No public
command accepts those states or native IDs.

The operator can run this recovery command:

```text
spacedock dispatch abandon --workflow-dir DIR --entity REF --stage STAGE \
  --epoch EPOCH --reason TEXT
```

The command requires a controlling terminal. It displays the exact receipt and
requires an explicit confirmation.

There is no `--yes` option. The event records the operating-system user,
launcher session, reason, and confirmation time.

Abandonment never permits stage advance. It permits only a replacement build for
the same stage.

On launcher resume, the observer first reconciles pending and spawned receipts.
It writes missing host events when the durable stream proves them.

If the stream is absent or corrupt, the receipt remains blocked. The operator
can then abandon the exact epoch.

## Shared mutation guard

One guard reduces the receipt event journal before any stage mutation.

The guard covers these product doors:

- `status --set` when `status` changes;
- terminal status changes, including `--force`;
- `status --archive`;
- `gate consume` when it applies `target-stage`;
- `merge guard` terminalization and archive;
- `dispatch build --advance`.

`--force` cannot bypass this guard. Abandonment also cannot satisfy it.

The guard checks the entity ID, entity path, stage, host, epoch, native IDs,
terminal event, and committed receipt state.

## Command and envelope impact

Remove the proposed `dispatch acknowledge` command. It must never ship.

Add `dispatch abandon` to command routing, root help, dispatch help, and command
tests.

Fresh Claude and Codex build output changes from schema version 2 to version 3.
It adds `dispatch_epoch` and `receipt_path`.

Fresh `description`, Claude `name`, and Codex `task_name` contain the bounded
epoch token. Advance envelopes keep their current shape.

The build command requires a live launcher observer for fresh Claude and Codex
envelopes. It fails before output when the observer session is absent.

## Acceptance criteria

**AC-1 (VALUE) — Each fresh envelope has one native worker.**

The exact Claude and Codex journeys create pending state, record one host start,
record matching completion, and then enter validation.

Falsifier: either journey advances without all receipt events, or creates more
than one native worker for the same epoch.

**AC-2 — Caller text has no acknowledgment authority.**

No command or API accepts caller-supplied spawned or completed state.

Falsifier: a caller can create either event without a matching host event.

**AC-3 — Same-stage replay fails closed.**

A second build refuses pending, spawned, or completed state. An old completion
cannot complete a replacement epoch.

Falsifier: a second envelope is emitted, or an old host event changes the new
receipt.

**AC-4 — Recovery has explicit authority.**

Only an operator-confirmed abandonment permits a replacement. Abandonment never
permits stage advance.

Falsifier: a noninteractive caller abandons a receipt, or an abandoned receipt
satisfies the mutation guard.

**AC-5 — Concurrent writes preserve both records.**

A worker report commit and observer receipt commit keep both paths intact during
the race suite.

Falsifier: either file is lost, overwritten, staged by the wrong writer, or left
under an index lock.

**AC-6 — Every mutation door consumes the receipt.**

Direct status, terminal status, archive, gate consume, merge, and `--advance`
all use the shared guard.

Falsifier: any door changes entity or gate bytes before matching completion.

**AC-7 — Exact live proof removes only owned XFAILs.**

The Claude and Codex `default-headless-gate-stop` cells first report strict
XPASS. The implementation then removes only this task's bindings.

Falsifier: a cell still XFAILs, another binding changes, or the normal rerun does
not pass.

**AC-8 — Pi and legacy behavior stays stable.**

Pi and entities without a new receipt keep their current command results.

Falsifier: a Pi envelope gains receipt fields, or a legacy fixture gains a new
guard.

## Explicit non-goals

- The binary does not call the host spawn tool.
- The observer does not judge worker report content.
- The task does not change scheduler priority.
- The task does not change gate approval authority.
- The task does not add Pi receipt support.
- The task does not remove the Pi live XFAIL.
- The task does not accept caller-supplied native IDs.

## Expected implementation surface

The implementation changes exactly these 34 files:

- `internal/dispatchreceipt/model.go`: 170 insertions.
- `internal/dispatchreceipt/store.go`: 220 insertions.
- `internal/dispatchreceipt/guard.go`: 130 insertions.
- `internal/dispatchreceipt/observer.go`: 210 insertions.
- `internal/dispatchreceipt/claude_events.go`: 135 insertions.
- `internal/dispatchreceipt/codex_events.go`: 135 insertions.
- `internal/dispatchreceipt/store_test.go`: 260 insertions.
- `internal/dispatchreceipt/observer_test.go`: 330 insertions.
- `internal/dispatchreceipt/guard_test.go`: 230 insertions.
- `docs/specs/dispatch-receipts.md`: 150 insertions.
- `internal/dispatch/build.go`: 85 insertions and 18 deletions.
- `internal/dispatch/codex_v2_adapter.go`: 20 insertions and 6 deletions.
- `internal/dispatch/dispatch.go`: 95 insertions and 8 deletions.
- `internal/dispatch/help_test.go`: 45 insertions and 10 deletions.
- `internal/dispatch/build_receipt_test.go`: 210 insertions.
- `internal/dispatch/advance_receipt_test.go`: 125 insertions.
- `internal/dispatch/abandon_test.go`: 155 insertions.
- `internal/cli/host_exec.go`: 95 insertions and 12 deletions.
- `internal/cli/host_launch_test.go`: 180 insertions.
- `internal/cli/cli.go`: 4 insertions and 1 deletion.
- `internal/cli/help.go`: 10 insertions and 3 deletions.
- `internal/cli/help_test.go`: 25 insertions and 5 deletions.
- `internal/status/handlers.go`: 28 insertions.
- `internal/status/mutate.go`: 18 insertions.
- `internal/status/merge.go`: 18 insertions.
- `internal/status/dispatch_receipt_guard_test.go`: 220 insertions.
- `internal/gates/application.go`: 18 insertions.
- `internal/gates/application_test.go`: 110 insertions.
- `skills/first-officer/references/fo-dispatch-core.md`: 12 insertions and 7 deletions.
- `skills/first-officer/references/claude-fo-dispatch.md`: 10 insertions and 5 deletions.
- `skills/first-officer/references/codex-first-officer-runtime.md`: 10 insertions and 5 deletions.
- `docs/runtime-support.md`: 35 insertions and 10 deletions.
- `internal/ensigncycle/shared_live_runner_test.go`: 1 insertion and 1 deletion.
- `internal/contractlint/live_registry_reconciliation_test.go`: 1 insertion and 1 deletion.

The estimate is 3,134 insertions and 92 deletions. The total is 3,226 gross
lines and 3,042 net lines.

The tolerance is 320 gross lines. Any new file or larger change requires a new
review before implementation continues.

The two one-line XFAIL edits occur only after strict XPASS evidence.

## Test plan

1. Add schema tests before store code.
2. Add store tests for immutable events, locks, atomic create, and state sync.
3. Add parser fixtures from the exact Claude and Codex spike events.
4. Add parser tests for unknown schemas, truncation, duplicate events, and stale IDs.
5. Add build tests for pending creation with and without `--stamp`.
6. Add second-build tests for pending, spawned, completed, and abandoned states.
7. Add operator tests with an injected terminal confirmer.
8. Add guard tests for every mutation door and `--force`.
9. Add a race test for observer events and a path-scoped worker report commit.
10. Add launcher tests for start, resume, cursor replay, crash, and signal exit.
11. Run Claude, Codex, Pi, and legacy controls.
12. Run `gofmt -w ./cmd ./internal`.
13. Run `go test ./...`.
14. Run `go test ./... -race`.
15. Rebind the Claude and Codex XFAIL rows to this task in one baseline commit.
16. Run both exact live cells and require strict XPASS.
17. Remove only the two bindings owned by this task.
18. Rerun both exact live cells and require normal PASS.

The exact live cells run last because they have model and runtime cost.

## Readiness

The revised design is ready for implementation review. All six Material
findings have task-owned mechanisms and exact falsifiers.

Readiness recommendation: **APPROVE**.

## Stage Report: ideation

- DONE: Accept and disposition all six Material findings.
  The table gives materiality, ownership, disposition, and evidence.
- DONE: Run the exact second spike without product edits.
  The disposable observer exercised Claude, Codex, replay, recovery, and races.
- DONE: Prove timely native host identities.
  Both hosts emitted stable worker IDs before completion.
- DONE: Replace the forgeable caller acknowledgment design.
  The launcher observer is the only spawn and completion writer.
- DONE: Define schema, commands, recovery authority, and mutation doors.
  The body defines immutable events and one shared guard.
- DONE: Preserve end-to-end and exact live proof.
  The acceptance criteria prohibit a component-only result.
- DONE: Define the exact implementation surface and test order.
  The surface has 34 files and focused tests before live tests.
- DONE: Write the revised report in Simplified Technical English.
  The report uses short active sentences and stable terms.
- DONE: Commit and push the revised ideation report.
  The state commit and push make the revised design durable.

### Summary

One launcher-owned observer now replaces caller acknowledgment. Claude and Codex
both provide stable native start and completion identities.

The receipt journal is separate from the worker report. Every stage mutation
door consumes the same completed receipt.
