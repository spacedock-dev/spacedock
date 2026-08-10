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

## Ideation HOLD report: supported-hooks spike

Date: 2026-08-10

Readiness recommendation: **HOLD**.

This report records the authorized supported-hooks spike. It does not revise
the product design, acceptance criteria, expected surface, registry, or XFAIL
bindings.

### Durable spike artifact

The complete evidence report is:

```text
/tmp/n284-hooks-spike.9atOqU/EVIDENCE.md
sha256:57c348f1974c369df9291106cc6997cbe11e7a06c5c9b8adfb80eb30646a723b
```

The main evidence artifacts are:

```text
sha256:7d6b8cc069d2f9f862fd941c835b95217276cee7120d6062ac7127e767df7220  receipt-hook.sh
sha256:dd123005cf742ed0c66b07b61d9e3b4e5d60aa18f1aa73012ccc8318754b1d56  codex-v2-hooks.jsonl
sha256:9ba3641972295e6bc161325175b9b9258b046928c4c976b2b65ce58a703b2bb1  claude-v2-hooks.jsonl
sha256:eca9af9843cd1bb10081ae3f45f39b55dffbe68bc77efb342ae8e2ec259451a0  codex-parallel-v2-hooks.jsonl
sha256:336c4d52a4676f3938edc7f8f29b4f470249396a98e3629d5c4af9491676480e  codex-reuse-hooks.jsonl
sha256:2379efd23cdea9a23c5bd15fdfec2d95ad43bc47ef43eade2a8db0733ffb859b  recovery-fixed.log
sha256:07c7de09ab6af80ee30bc8806246063d6217d53ff5f04a879c13a58b489d3cfa  codex-disabled.stream.jsonl
sha256:3f0b425503ec9b1e5f9e1ec581819cd653d5cc087c5b58f856bff0b29ee111c2  claude-disabled.stream.jsonl
```

All paths are below `/tmp/n284-hooks-spike.9atOqU`.

### Passed cells

- Codex CLI created and completed one schema-v2 receipt.
- Claude CLI created and completed one schema-v2 receipt.
- A duplicate same-generation prepare returned exit code 128.
- Codex feedback generation g2 reused the g1 native agent ID.
- A consumed captain resolution abandoned one exact receipt.
- Resolution reuse failed, and abandonment did not satisfy completion.
- Private Git refs survived 240 concurrent compare updates.
- A production-style `git add -A` commit changed only `entity-a.md`.
- Codex with hooks disabled left its receipt pending.
- Claude with hooks disabled left its receipt pending.
- Both disabled-hook cells refused a second same-stage prepare.

The Codex CLI version was `0.147.0`. The Claude Code version was `2.1.226`.

### Failed parallel correlation

The Codex parallel cell emitted two `PreToolUse` events before any
`SubagentStart` event.

Each pre-use event had a task name and tool-use ID. Each start event had an
agent ID, but it had no parent correlation ID.

The probe assigned the starts by arrival order. This assignment swapped the
two receipts.

Entity n recorded agent `019feb30-33d6-7f40-98e3-0b43f75f7f85`. That agent
stopped with `V2_MARKER_O`.

Entity o recorded agent `019feb30-27b6-7d40-a67e-3639385ba896`. That agent
stopped with `V2_MARKER_N`.

The parent command returned zero. The attestation result was still incorrect.

### Unexecutable host surfaces

The local Codex Desktop application was not installed. The environment could
not execute the available VS Code launcher.

The Codex cloud CLI had no authorized disposable cloud project. It also had no
controllable web session.

Claude Remote Control could not start because feature-flag evaluation was
disabled.

Therefore, the IDE, Desktop, and web cells were unexecutable. This spike gives
no product evidence for those surfaces.

### Evidence-record defect

The race generator command was not saved. The final ref value, commit, Git
integrity check, and worktree result remain available.

This missing command is an evidence-record defect. A later spike must save the
complete race command before it makes a concurrency claim.

### Required host prerequisite

Codex must add one stable parent correlation value to supported
`SubagentStart` events.

The preferred value is the parent `tool_use_id`. A stable assignment ID is also
sufficient.

Codex must include the same value in `SubagentStop`. CLI and in-app hosts must
use the same contract.

After that host change, rerun the disposable spike. Then run supported in-app
cells before any product design revision.

No product implementation can start from this report.

## Independent staff-review HOLD

Date: 2026-08-10

Staff recommendation: **HOLD**.

This review does not revise the product design. It records the remaining proof
gaps from the supported-hooks spike.

### Open Material findings

| Finding | Status | Materiality | Review result |
|---|---|---|---|
| Parallel Codex correlation | OPEN | Material | Two starts had no parent correlation ID. The probe assigned the wrong native IDs to both receipts. |
| Recovery authority | OPEN | Material | The probe script created its own approval. The abandon action did not recompute and verify the stored digest. |

The recovery cell proved atomic consumption only. It did not prove independent
captain authority.

### Unproved cells

- Parallel Claude correlation remains unproved.
- Exact feedback reuse generation remains unproved.
- IDE, Desktop, and web host surfaces remain unproved.
- Direct `status --set` remains unproved.
- Terminal status mutation remains unproved.
- Archive mutation remains unproved.
- Gate consumption remains unproved.
- Merge finalization remains unproved.
- `dispatch build --advance` remains unproved.

The feedback cell reused one native agent. It did not bind the follow-up call to
one exact pending generation from supported host data.

### Closed commit-race finding

Private Git refs close the companion-file commit race. The production-style
`git add -A` commit included only the worker report.

The race generator command was not saved. This remains an evidence-record
defect, but it does not reopen the private-ref commit boundary.

### Authorization

No product implementation is authorized. Resolve the two open Material
findings and prove all listed cells before design revision or implementation.

## Captain rescope for staff review

Date: 2026-08-10

Review state: **HOLD — staff review requested**.

This section supersedes the earlier approach, criteria, surface estimate, and
readiness recommendation. It does not authorize product implementation.

### Narrow product claim

The supported hosts are Claude CLI and Codex CLI only.

One workflow entity and stage can have one in-flight native worker. The host
must attest the start and stop of that worker.

The mechanism denies ambiguous concurrent starts. It does not correlate a
general parallel spawn batch.

These surfaces are out of scope:

- Claude Desktop and Codex Desktop;
- IDE extensions;
- web and cloud sessions;
- general parallel worker acknowledgment;
- cryptographic proof of the local operating-system user;
- Pi acknowledgment;
- a new acknowledgment for `dispatch build --advance` worker reuse.

The in-app surfaces were unexecutable in the authorized environment. This task
makes no compatibility claim for them.

### Trust boundary

Supported host hooks are the attestation source. The mechanism does not parse a
Claude transcript or Codex rollout file.

The local host process, launcher process, and Spacedock binary are trusted. The
threat model does not resist a malicious local process with state-repository
write access.

No public command accepts a caller-supplied native start, stop, or worker ID.
Caller text, reports, final messages, and empty waits have no receipt authority.

### Single-flight mechanism

The CLI launcher creates a random launcher session ID. It also creates a private
runtime directory with mode `0700`.

The launcher exports both values to its Claude or Codex child. A fresh build
refuses when these launcher values are absent.

A fresh Claude or Codex `dispatch build` uses this sequence:

1. Resolve the entity ID, entity path, stage, host, and state Git repository.
2. Create a random 128-bit assignment generation and a random 128-bit epoch.
3. Create one schema-v2 pending receipt through a zero-old compare update.
4. Write one mode-`0600` runtime registration for the epoch.
5. Put the bounded epoch token in Claude `description` or Codex `task_name`.
6. Emit the existing schema-v2 build envelope with one receipt object.

The receipt ref is:

```text
refs/spacedock/dispatch/v2/<entity-id>/<stage>/<assignment-generation>
```

Each receipt update writes a new JSON blob. The blob contains the prior blob ID.
The ref moves with one compare-and-swap update.

The receipt states are:

```text
pending -> armed -> spawned -> completed -> advanced
                    \-> abandoned
pending|armed ------> abandoned
```

`abandoned` never satisfies a stage guard. `advanced` records that one completed
stage instance left its stage.

### Hook behavior

The plugin registers supported `PreToolUse`, `SubagentStart`, and
`SubagentStop` hooks for each CLI host.

`PreToolUse` extracts the bounded epoch from the native assignment field. It
then resolves the exact runtime registration and pending receipt.

The hook acquires one launcher-session `unbound-start` file with atomic create.
It records the host session and tool-use ID, then moves the receipt to `armed`.

If an unbound start already exists, `PreToolUse` returns the host's supported
deny result. The second native tool call must not execute.

This denial applies across entities in one launcher session. It is deliberately
stricter than the entity-and-stage rule.

`SubagentStart` can bind only the single armed receipt. It records the native
worker ID, moves the receipt to `spawned`, and removes the unbound-start file.

`SubagentStop` finds the spawned receipt by host session and native worker ID.
It moves only that receipt to `completed`.

An absent, duplicate, stale, or ambiguous event changes no authoritative state.
The affected receipt stays blocked or moves to an explicit error state.

The future spike must prove that both hosts deny the second tool call. A hook
warning or non-blocking error is not sufficient.

### Schema-v2 receipt

Each receipt blob has these required fields:

```json
{
  "schema_version": 2,
  "state": "pending|armed|spawned|completed|abandoned|advanced|error",
  "entity_id": "stable workflow entity ID",
  "entity_path": "state-root-relative path",
  "stage": "implementation",
  "host": "claude|codex",
  "assignment_generation": "random-128-bit value",
  "epoch": "random-128-bit value",
  "launcher_session_id": "random launcher value",
  "host_session_id": "supported hook session ID",
  "tool_use_id": "supported PreToolUse ID",
  "native_worker_id": "supported SubagentStart ID",
  "previous_blob": "Git object ID or empty",
  "observed_at": "RFC3339 timestamp"
}
```

Fields that are not known in the current state are empty. A state transition
must add every field that becomes required for that state.

The existing build request and response remain schema version 2. A fresh Claude
or Codex response adds this object:

```json
{
  "dispatch_receipt": {
    "schema_version": 2,
    "assignment_generation": "...",
    "epoch": "...",
    "ref": "refs/spacedock/dispatch/v2/..."
  }
}
```

Pi, validation-only, and schema-print output do not add this object.

### Build replay and reuse

A fresh build refuses an active `pending`, `armed`, `spawned`, or `completed`
receipt for the same entity and stage. It emits no second envelope.

A successful stage change moves the matching completed receipt to `advanced`.
A later feedback re-entry can then create a new assignment generation.

`dispatch build --advance` does not start a native worker. It stays outside the
new native-spawn acknowledgment claim.

The advance path must name the completed source receipt and existing native
worker handle. It refuses a pending, abandoned, missing, or ambiguous source.

The task does not claim a new host-attested completion for a follow-up message.
That separate reuse-generation claim needs a future task.

### Every mutation surface

One shared receipt guard covers all lifecycle mutation doors. `--force` does
not bypass this guard.

The exact guarded surfaces are:

1. `status --set` when `status` changes to another stage.
2. `status --set` when `completed` or `verdict` finalizes lifecycle state.
3. `status --archive <entity>`.
4. Standalone `gate consume <entity>` when it changes stage.
5. `gate record --consume` through the same consume path.
6. `merge guard --verdict passed|rejected` when it finalizes or archives.
7. `merge guard --rework` when it routes to `feedback-to`.
8. `dispatch build --advance` before it emits a reuse envelope.

Non-lifecycle field updates, gate preparation, gate presentation, and gate
recording without `--consume` do not consume a receipt.

A fresh `dispatch build` is also a receipt mutation. It uses the zero-old
compare update before it emits an envelope.

`dispatch build --stamp` reserves the receipt before stamp writes. If stamping
fails, the receipt stays blocked and uses recovery.

Each guarded command uses this order:

1. Read the active receipt and current entity bytes.
2. Require `completed` for the entity, stage, host, and generation.
3. Perform the existing entity or gate mutation.
4. Move the completed receipt to `advanced` after a successful stage exit.
5. Fail closed if receipt advancement cannot become durable.

No command changes entity or gate bytes after a receipt refusal.

### Recovery surfaces and authority

The public recovery surfaces are:

```text
spacedock dispatch receipt show --workflow-dir DIR --entity REF --stage STAGE
spacedock dispatch receipt recovery-request --workflow-dir DIR --entity REF \
  --stage STAGE --generation GEN --reason TEXT
spacedock dispatch receipt abandon --workflow-dir DIR --entity REF \
  --stage STAGE --generation GEN --reason TEXT --resolution-file FILE
spacedock dispatch receipt reconcile --workflow-dir DIR --entity REF --stage STAGE
```

`show` is read-only. `recovery-request` is also read-only and prints canonical
JSON plus its SHA-256 digest.

The captain resolution is a separate retained record. The abandon command does
not create or modify that approval.

The resolution must name the request digest, entity ID, stage, generation,
epoch, action, reason, decision, and captain authority.

`abandon` rebuilds the canonical request from current state. It recomputes the
digest and compares every bound field.

One Git ref transaction moves the receipt to `abandoned` and the resolution to
`consumed`. Resolution reuse or digest drift changes neither ref.

`reconcile` can apply only supported hook evidence already registered for the
same launcher session and native worker ID. It cannot accept caller IDs.

These failures stay blocked until reconciliation or authorized abandonment:

- envelope output fails after pending creation;
- hooks are absent or disabled;
- a hook exits, times out, or returns malformed output;
- the host CLI or launcher crashes;
- the state-ref push fails;
- an event is stale, duplicated, or ambiguous.

Abandonment permits a replacement generation. It never satisfies a mutation
guard and never advances the entity.

### Private-ref durability and concurrency

Receipt and resolution refs live outside the Git index. Therefore, a worker's
production-style `git add -A` cannot stage them.

Every ref write uses an expected old object ID. A compare failure preserves the
winning value and returns a blocking error.

Split-root publication pushes the exact private ref with a compare lease. A
push failure leaves local evidence and blocks lifecycle mutation.

The concurrency proof must save its full generator command before execution.
The prior 240-update result is useful, but its unsaved command is not sufficient.

### Acceptance criteria

**AC-1 — One fresh CLI envelope gets one native worker.**

Claude CLI and Codex CLI each create one pending receipt. Supported hooks move
that receipt through armed, spawned, and completed.

Falsifier: a caller statement changes state, a native ID is missing, or more
than one native worker starts for the epoch.

**AC-2 — Same entity and stage is single-flight.**

A second fresh build refuses pending, armed, spawned, and completed states. It
emits no envelope and changes no receipt.

Falsifier: any active state permits a second envelope.

**AC-3 — Ambiguous concurrency fails closed.**

For each CLI host, two spawn calls in one parallel batch cause one supported
PreToolUse denial. Only one native worker starts.

Falsifier: both tool calls execute, receipt identities swap, or denial is only a
warning.

**AC-4 — Completion uses supported host identity.**

The completing stop must match host session and native worker ID from the start.

Falsifier: a report, transcript, rollout, final message, or stale stop completes
the receipt.

**AC-5 — Every lifecycle mutation door uses one guard.**

All eight listed doors refuse pending, spawned, abandoned, missing, and error
states before entity or gate bytes change.

Falsifier: one door mutates bytes or `--force` bypasses the receipt.

**AC-6 — Recovery has separate authority.**

Abandonment requires a separate unconsumed captain resolution. The binary
recomputes the request digest and consumes both records atomically.

Falsifier: the abandon path creates its approval, accepts digest drift, reuses a
resolution, or advances the entity.

**AC-7 — Private refs survive production Git writes.**

A saved race command performs ref compare updates while a worker runs
`git add -A` and commits its report.

Falsifier: the report commit includes receipt data, either value is lost, or Git
integrity fails.

**AC-8 — Disabled or broken hooks fail closed.**

Both CLI hosts complete a native child with hooks disabled. The receipt stays
pending, every lifecycle door refuses, and a second build refuses.

Falsifier: disabled hooks permit stage change or replacement dispatch.

**AC-9 — Scope exclusions stay explicit.**

No test or documentation claims in-app or general parallel acknowledgment. Pi
and legacy entities keep their current results.

Falsifier: an in-app or Pi claim appears, or a legacy fixture gains a receipt.

**AC-10 — Both default-headless bindings leave only after exact live proof.**

The Claude Sonnet and Codex CLI cells first produce strict XPASS with their
bindings present. After binding removal, both cells pass normally.

Falsifier: either cell skips, XFAILs, uses private host files, or loses the open
human gate.

### Exact default-headless proof

Run the Claude CLI cell with this command:

```sh
SPACEDOCK_LIVE_RUNTIME=claude \
SPACEDOCK_LIVE_MODEL=sonnet \
SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/n284-live-claude \
go test -tags=live ./internal/ensigncycle \
  -run '^TestLiveCommonDefaultHeadlessGateStop$' -count=1 -v
```

Run the Codex CLI cell with this command:

```sh
SPACEDOCK_LIVE_RUNTIME=codex \
SPACEDOCK_CODEX_LIVE_REQUIRED=1 \
SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/n284-live-codex \
go test -tags=live ./internal/ensigncycle \
  -run '^TestLiveCommonDefaultHeadlessGateStop$' -count=1 -v
```

Each artifact directory must prove all these facts:

- one fresh build envelope contains one schema-v2 receipt object;
- one supported PreToolUse event names the same epoch and tool-use ID;
- one supported start records one native worker ID;
- one supported stop records the same native worker ID;
- the receipt journal ends at `completed` with no error or unbound lock;
- the implementation report is complete, committed, and owned by that worker;
- the entity enters validation only after completion;
- the validation gate is prepared, committed, presented, and still open;
- the entity is not terminal and is not archived;
- the command log contains no second fresh build;
- the harness reads receipt artifacts, not Claude transcripts or Codex rollouts.

With both XFAIL bindings present, each command must fail only as strict XPASS.
The log must show `XPASS ALERT` with no semantic failure code.

Then remove only the Claude Sonnet and Codex entries from
`TestLiveCommonDefaultHeadlessGateStop`. Remove their matching expected entries
from `live_registry_reconciliation_test.go` in the same binding-only commit.

Rerun both commands. Each must report normal PASS with the same durable receipt
and open-gate evidence.

Finally, run:

```sh
go test ./internal/contractlint -run LiveRegistryReconciliation -count=1
go test ./...
go test ./... -race
```

The Pi binding and every unrelated TODO or XFAIL entry must remain unchanged.

### Reconciliation of prior staff findings

| Prior finding | Narrow disposition |
|---|---|
| Private Claude transcripts and Codex rollouts are unsupported. | In scope and closed by design. Only supported hooks can attest start and stop. Live proof remains required. |
| The resident observer excludes in-app hosts. | Removed from the claim. In-app hosts are out of scope and were unexecutable. |
| Companion receipt files enter `git add -A`. | Closed by private refs. The race must be rerun with its command saved. |
| Completed same-stage state blocks feedback re-entry. | Addressed by `completed -> advanced` on stage exit. A later stage instance gets a new generation. |
| Exact reuse generation is unproved. | Removed from native-spawn acknowledgment. `--advance` starts no native worker and must name the completed source receipt. |
| TTY confirmation does not prove recovery authority. | TTY confirmation is removed. A separate retained captain resolution is required, and abandon recomputes its digest. |
| Schema and build-to-observer IPC were undefined. | Addressed by schema-v2 private refs and a mode-`0600` launcher-session registration. Proof remains required. |
| Parallel Codex correlation is unsafe. | General parallel acknowledgment is out of scope. The second unbound start must receive a supported denial. |
| Parallel Claude correlation is unproved. | General parallel acknowledgment is out of scope. The same negative denial proof is required for Claude CLI. |
| IDE, Desktop, and web cells are unproved. | They remain unexecutable and out of scope. |
| Six mutation doors were unproved. | The rescope names eight exact doors. Every door is an acceptance-test requirement. |

### Exact expected implementation surface

The estimate includes exactly 40 files:

| File | Insertions | Deletions | Purpose |
|---|---:|---:|---|
| `.claude-plugin/plugin.json` | 1 | 0 | Register the Claude hook file. |
| `claude-hooks.json` | 58 | 0 | Register supported Claude CLI events. |
| `hooks.json` | 55 | 2 | Register supported Codex CLI events. |
| `internal/dispatchack/model.go` | 150 | 0 | Define schema-v2 receipt and resolution models. |
| `internal/dispatchack/store.go` | 280 | 0 | Own private refs, blobs, compare updates, and pushes. |
| `internal/dispatchack/hook.go` | 260 | 0 | Process supported hook input and single-flight denial. |
| `internal/dispatchack/guard.go` | 165 | 0 | Reduce receipts and guard lifecycle mutation. |
| `internal/dispatchack/recovery.go` | 185 | 0 | Build requests, verify resolutions, and reconcile. |
| `internal/dispatchack/model_test.go` | 180 | 0 | Test schema and canonical digests. |
| `internal/dispatchack/store_test.go` | 310 | 0 | Test compare updates, ref pushes, and races. |
| `internal/dispatchack/hook_test.go` | 340 | 0 | Test both hosts, denial, replay, and native IDs. |
| `internal/dispatchack/guard_test.go` | 245 | 0 | Test allowed and refused receipt states. |
| `internal/dispatchack/recovery_test.go` | 260 | 0 | Test separate authority and digest recomputation. |
| `internal/dispatch/build.go` | 78 | 14 | Reserve and emit fresh receipt data. |
| `internal/dispatch/dispatch.go` | 120 | 8 | Route hook, show, request, abandon, and reconcile. |
| `internal/dispatch/build_ack_test.go` | 260 | 0 | Test fresh, stamped, no-stamp, Pi, and replay builds. |
| `internal/dispatch/dispatch_ack_command_test.go` | 285 | 0 | Test every receipt command and help result. |
| `internal/cli/frontdoor.go` | 28 | 4 | Create and export the CLI launcher session. |
| `internal/cli/host_exec.go` | 18 | 2 | Preserve the private runtime environment. |
| `internal/cli/host_launch_test.go` | 125 | 0 | Test CLI session setup, cleanup, and crash behavior. |
| `internal/status/handlers.go` | 32 | 2 | Guard status changes and archive. |
| `internal/status/merge.go` | 24 | 2 | Guard finalize and rework routes. |
| `internal/status/dispatch_ack_guard_test.go` | 260 | 0 | Test status, archive, force, merge, and byte-clean refusal. |
| `internal/gates/application.go` | 20 | 2 | Guard stage-changing consumption. |
| `internal/gates/application_test.go` | 105 | 0 | Test pending and completed consume states. |
| `internal/cli/gate_ceremony.go` | 14 | 2 | Keep standalone and record-consume on one guard path. |
| `internal/cli/terminal_consume_test.go` | 155 | 0 | Test terminal, consume, and merge command doors. |
| `internal/cli/help.go` | 28 | 4 | Document receipt and recovery commands. |
| `internal/cli/help_test.go` | 75 | 0 | Pin root and dispatch help. |
| `skills/first-officer/references/fo-dispatch-core.md` | 24 | 8 | Use the single-flight CLI workflow. |
| `skills/first-officer/references/claude-fo-dispatch.md` | 12 | 5 | State Claude CLI hook and denial behavior. |
| `skills/first-officer/references/codex-first-officer-runtime.md` | 12 | 5 | State Codex CLI hook and denial behavior. |
| `docs/specs/dispatch-acknowledgment.md` | 220 | 0 | Specify trust, refs, states, guards, and recovery. |
| `docs/runtime-support.md` | 35 | 8 | Limit the support claim to both CLIs. |
| `internal/ensigncycle/shared_live_runner_test.go` | 10 | 2 | Read receipts, then remove two runtime bindings. |
| `internal/ensigncycle/claude_live_runner_test.go` | 80 | 8 | Retain supported Claude hook artifacts. |
| `internal/ensigncycle/codex_live_runner_test.go` | 75 | 32 | Replace rollout reads with receipt evidence. |
| `internal/contractlint/live_registry_reconciliation_test.go` | 1 | 1 | Remove the two mirrored expected bindings. |
| `skills/integration/plugin_manifest_test.go` | 45 | 0 | Pin both hook registrations. |
| `internal/ensigncycle/shared_assertions_impl_test.go` | 55 | 10 | Assert receipt order and the open gate. |

The estimate is 4,685 insertions and 121 deletions. It is 4,806 gross lines
and 4,564 net lines.

The gross tolerance is 480 lines. A new file, a new host, or a general parallel
claim requires another staff review before implementation.

### Test order

1. Add model and canonical-digest tests.
2. Add private-ref compare and saved-command race tests.
3. Add Claude and Codex hook fixture tests.
4. Prove supported PreToolUse denial for both hosts in disposable CLI runs.
5. Add recovery authority and digest-drift tests.
6. Add fresh build, no-stamp, stamp-failure, replay, Pi, and legacy tests.
7. Add all eight mutation-door tests.
8. Add launcher crash, missing hooks, disabled hooks, and malformed hook tests.
9. Run both default-headless cells with bindings and require strict XPASS.
10. Remove only the two runtime bindings and their mirrored expected entries.
11. Rerun both cells and require normal PASS.
12. Run formatting, the full suite, and the race suite.

### Staff review request

Please review these exact questions:

1. Does supported PreToolUse denial establish safe single-flight behavior for
   each CLI host?
2. Is the local-host trust boundary sufficiently narrow and explicit?
3. Does the separate captain resolution close the recovery-authority finding?
4. Do the eight listed doors cover every lifecycle mutation path?
5. Is `completed -> advanced` sufficient for later feedback re-entry?
6. Does the exact XPASS-to-PASS sequence justify both binding removals?

No product implementation is authorized before this review accepts the
rescope and all Material findings.

## Smallest CLI happy-path rescope

Date: 2026-08-10

Review state: **HOLD — staff review requested**.

This section supersedes every earlier implementation approach, acceptance
criterion, file estimate, and readiness statement.

No product implementation is authorized.

### Exact claim

The task supports one fresh Claude CLI or Codex CLI worker for one workflow
entity and stage.

`dispatch build` creates one binary-owned pending envelope. Supported host hooks
consume it when the native worker starts.

While the envelope is pending, these existing paths fail closed:

- `dispatch build --advance` for that entity and stage;
- `status --set` when it changes that entity from the current stage.

After acknowledgment, the existing complete-stage-report guard still controls
the stage transition. This task does not add a completion receipt.

### Explicit exclusions

The task excludes these items:

- recovery commands and recovery files;
- abandonment and supersession;
- IDE, Desktop, web, and cloud hosts;
- general parallel acknowledgment;
- a new public schema family;
- a broad mutation inventory;
- Pi acknowledgment;
- host-attested follow-up completion;
- cryptographic local-user identity.

The in-app surfaces remain unexecutable. General parallel acknowledgment
remains unsupported.

### Small mechanism

The private active ref is:

```text
refs/spacedock/dispatch-ack/<entity-id>/<stage>
```

The active blob contains only these internal fields:

```json
{
  "state": "pending|armed|consumed",
  "entity_id": "stable entity ID",
  "entity_path": "state-root-relative path",
  "stage": "implementation",
  "host": "claude|codex",
  "epoch": "random-128-bit value",
  "host_session_id": "supported hook session ID",
  "tool_use_id": "supported PreToolUse ID",
  "native_worker_id": "supported SubagentStart ID"
}
```

This internal blob does not define a new public schema family. The existing
build envelope stays at schema version 2.

A fresh build adds only `dispatch_ack_epoch` and `dispatch_ack_ref` to its
existing output. Pi and `--advance` output do not add these fields.

The happy path is:

1. A fresh build creates the active ref with a zero-old compare update.
2. The build puts the bounded epoch in Claude `description` or Codex `task_name`.
3. `PreToolUse` matches that epoch and changes `pending` to `armed`.
4. `SubagentStart` records the native worker ID and changes `armed` to `consumed`.
5. The hook writes one consumed audit ref for live proof.
6. The existing worker writes and commits the complete stage report.
7. The existing status transition guard validates that report.
8. The successful stage transition removes the old active ref.

The consumed audit ref is:

```text
refs/spacedock/dispatch-ack-audit/<entity-id>/<stage>/<epoch>
```

The binary writes one mode-`0700` temporary directory for hook lookup. The
pending Git ref remains the authoritative block.

No public command accepts a native worker ID. Caller text, a report, a final
message, or an empty wait cannot consume the envelope.

### Single-flight and replay

Fresh build refuses an existing `pending`, `armed`, or `consumed` active ref. It
emits no second envelope.

The same envelope cannot arm twice. A stale hook event changes no ref.

This claim covers one native spawn at a time. If the hook sees more than one
armed candidate, it changes no ref and returns a blocking error.

The task does not claim that two different parallel spawns can be correlated.

### Exact fail-closed behavior

`dispatch build --advance` reads the active ref before it emits output. It
refuses only `pending` and `armed` states.

`status --set` reads the active ref before a stage change. It also refuses only
`pending` and `armed` states.

A `consumed` ref does not replace the existing complete-report guard. A stage
change still fails until that guard passes.

Non-stage status fields, gate commands, archive, merge, and recovery are outside
this rescope.

### Exact negative proof

The focused tests must prove these cases for both hosts:

1. A fresh build creates one pending ref and one epoch-bearing envelope.
2. A second fresh build returns nonzero and emits no envelope.
3. `dispatch build --advance` returns nonzero while pending or armed.
4. A stage-changing `status --set` returns nonzero and changes no entity bytes.
5. Disabled hooks leave the ref pending and keep both refusals active.
6. A malformed or stale hook event changes no ref.
7. `PreToolUse` cannot consume the ref; it can only arm it.
8. `SubagentStart` consumes only the single armed ref.
9. The consumed record contains the supported native worker ID.
10. A report or caller final message does not change the pending ref.

The negative host cells use local subscription access only. They must not use
paid CI.

### Exact default-headless proof

Use these existing exact targets:

```sh
SPACEDOCK_LIVE_RUNTIME=claude \
SPACEDOCK_LIVE_MODEL=sonnet \
SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/n284-happy-claude \
go test -tags=live ./internal/ensigncycle \
  -run '^TestLiveCommonDefaultHeadlessGateStop$' -count=1 -v
```

```sh
SPACEDOCK_LIVE_RUNTIME=codex \
SPACEDOCK_CODEX_LIVE_REQUIRED=1 \
SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/n284-happy-codex \
go test -tags=live ./internal/ensigncycle \
  -run '^TestLiveCommonDefaultHeadlessGateStop$' -count=1 -v
```

For each host, the bound run must prove these facts:

- one fresh implementation envelope created one pending ref;
- supported `PreToolUse` armed that exact epoch;
- supported `SubagentStart` consumed it with one native worker ID;
- no second implementation envelope was emitted;
- the implementation report became complete and committed;
- the entity did not enter validation before acknowledgment;
- the existing report guard passed before validation;
- the validation gate was prepared, committed, presented, and left open;
- the entity was not archived or terminalized;
- no Claude transcript or Codex rollout supplied acknowledgment evidence.

The current XPASS policy is lane-green. With each binding present, the exact
target must log `XPASS ALERT`, report no semantic failure code, and exit zero.

Then make one binding-only commit. Remove only these two runtime bindings:

- `claude-sonnet` for `default-headless-gate-stop`;
- `codex` for `default-headless-gate-stop`.

Remove the same two expected entries from the registry reconciliation test. Do
not change any other TODO, XFAIL, or runtime binding.

Rerun both exact commands. Each target must report normal PASS with the same
pending-to-consumed and open-gate evidence.

Finally, run:

```sh
go test ./internal/contractlint -run LiveRegistryReconciliation -count=1
go test ./...
go test ./... -race
```

If either host lacks a supported blocking `PreToolUse` or stable
`SubagentStart` worker ID, this design is **HOLD**. That missing hook primitive
would prevent removal of that host's binding.

The completed local spike already observed stable single-worker start IDs for
Claude CLI and Codex CLI. Implementation proof is still required.

### Prior finding reconciliation

| Prior finding | Smallest-scope disposition |
|---|---|
| Caller acknowledgment is forgeable. | Removed. Only supported `SubagentStart` can consume the pending ref. |
| Private transcript and rollout formats are unsupported. | Removed. The live proof uses the consumed audit ref. |
| Parallel Codex correlation swapped receipts. | General parallel acknowledgment is excluded. Ambiguous armed state blocks without mutation. |
| Parallel Claude correlation is unproved. | General parallel acknowledgment is excluded. |
| Recovery authority is unproved. | Recovery and abandonment are excluded. Pending state needs manual product support in a later task. |
| Companion files race with `git add -A`. | The active and audit records use private Git refs, not worktree files. |
| Exact reuse generation is unproved. | Follow-up acknowledgment is excluded. `--advance` only refuses an unacknowledged fresh spawn. |
| In-app surfaces are unproved. | They remain unexecutable and excluded. |
| Broad mutation doors are unproved. | The claim covers only `--advance` and stage-changing `status --set`. |
| A new schema and observer IPC were too large. | No new public schema exists. One small temporary lookup joins build and hooks. |

### Exact file and line estimate

The implementation target is exactly ten files:

| File | Insertions | Deletions | Purpose |
|---|---:|---:|---|
| `.claude-plugin/plugin.json` | 1 | 0 | Reuse the shared hook file in Claude CLI. |
| `hooks.json` | 36 | 2 | Register PreToolUse and SubagentStart for both CLI hosts. |
| `internal/dispatchack/ack.go` | 145 | 0 | Own active refs, audit refs, hook parsing, and guards. |
| `internal/dispatchack/ack_test.go` | 130 | 0 | Test both hosts, replay, disabled hooks, advance, and status. |
| `internal/dispatch/build.go` | 28 | 6 | Create pending state and guard advance. |
| `internal/dispatch/dispatch.go` | 35 | 4 | Route the internal hook command and print errors. |
| `internal/status/handlers.go` | 14 | 2 | Refuse a stage change while pending or armed. |
| `internal/ensigncycle/claude_live_runner_test.go` | 30 | 22 | Read audit refs instead of private host files. |
| `internal/ensigncycle/shared_live_runner_test.go` | 2 | 2 | Remove only the two runtime bindings after XPASS. |
| `internal/contractlint/live_registry_reconciliation_test.go` | 1 | 1 | Remove the two matching expected entries. |

The estimate is 422 insertions and 39 deletions. It is 461 gross lines and 383
net lines.

The hard limit is ten files and 500 gross lines. Any extra file or line above
that limit requires captain approval and another staff review.

### Staff review request

Please review these exact questions:

1. Can each supported CLI deny a duplicate `PreToolUse` and expose one stable
   `SubagentStart` worker ID?
2. Does pending-only acknowledgment close both default-headless failures when
   combined with the existing complete-report guard?
3. Are `--advance` and stage-changing `status --set` the complete happy-path
   mutation boundary for this narrow claim?
4. Does the XPASS-green sequence prove each binding before removal?
5. Can the ten-file, 461-gross limit hold without a component-only result?

No product bytes can change before staff accepts this rescope.

## Staff review of CLI single-flight rescope

Date: 2026-08-10

We love you.

Recommendation: **HOLD**.

The CLI-only rescope removes the unsafe parallel-correlation claim. The design
still has three open Material findings and one incomplete proof contract.

### 1. Supported PreToolUse denial

Answer: **Yes, with a required correction and live proof.**

Both supported hosts can deny a tool before execution. Claude supports a
`PreToolUse` deny result or exit code 2. Codex supports the same deny result.

The official Codex hook contract also states that `spawn_agent` matches
`Agent`. Thus, the selected event is a supported CLI boundary.

The design denies only when an `unbound-start` file exists (lines 751-755). It
does not explicitly deny a missing, stale, duplicate, or non-pending epoch.

Those events only change no receipt at lines 766-767. A second native worker
can still start if the hook returns success.

The hook must deny every assignment that cannot bind one exact pending receipt.
The live cells must prove this rule for both hosts.

The supported denial closes the prior assignment-swap defect only under that
rule. One armed receipt then exists when `SubagentStart` arrives.

This task does not claim general parallel correlation. It also makes no claim
for an in-app host (lines 681-692).

References:

- <https://learn.chatgpt.com/docs/hooks>
- <https://code.claude.com/docs/en/hooks>

### 2. Local-host trust boundary

Answer: **Yes, for a cooperative local CLI host.**

Lines 694-704 identify the trusted host, launcher, and binary. They also exclude
a malicious local process that can write the state repository.

This boundary is explicit. It is not a security boundary against a malicious
model, shell process, or operating-system user.

The product documentation must keep that limitation. It must not describe the
receipt as cryptographic proof or user authentication.

### 3. Separate recovery authority

Answer: **No. The recovery-authority finding remains Material.**

The new flow correctly separates request creation from abandonment. It also
requires digest recomputation and atomic consumption (lines 879-892).

However, `abandon` accepts `--resolution-file FILE` at lines 870-876. The design
does not identify an authority key, protected store, or trusted record creator.

A worker can run `recovery-request`, create the stated file, and name itself as
captain authority. The binary has no defined fact that rejects this forgery.

The text later calls the resolution a Git ref. It does not define how the file
maps to that retained ref.

The design must name a separate authority source and its verification rule. The
abandon command must consume that exact retained record.

### 4. Eight mutation doors

Answer: **No. The list misses one current lifecycle mutation.**

The eight doors cover the public stage, terminal, archive, gate, merge, and
advance commands at lines 836-845.

However, `status --set worktree=` also removes lifecycle ownership. Current
code classifies this update as terminal at `internal/status/handlers.go:133-145`.

The list guards `completed` and `verdict`, but it does not guard a standalone
worktree clear. `--force` makes this omission important.

Add this mutation to the shared guard or state why it cannot affect an active
receipt. Add a byte-clean refusal test for its forced and unforced forms.

The post-mutation order also needs a failure rule. Lines 856-864 mutate entity
bytes before the receipt moves to `advanced`.

If the ref update fails, entity bytes already changed. The design must specify
one atomic transaction or a tested rollback that restores all entity and gate
bytes.

### 5. Completed-to-advanced feedback re-entry

Answer: **Yes, if the stage-exit transaction is atomic.**

Lines 816-820 block a second active generation. A successful exit moves the
completed receipt to `advanced`.

A later feedback return can then create a new generation for the same entity
and stage. This state model closes the permanent same-stage block.

The guard must select the exact completed generation. The stage mutation and
the `advanced` transition must succeed or fail together.

### 6. Exact XPASS-to-PASS proof

Answer: **No. The stated strict-XPASS sequence does not match the harness.**

The two bindings exist at `internal/ensigncycle/shared_live_runner_test.go:108-111`.
Their mirrored entries exist at
`internal/contractlint/live_registry_reconciliation_test.go:51-54`.

The harness classifies a clean expected failure as XPASS at
`internal/ensigncycle/claude_runtime_helpers_test.go:330-356`.

However, `liveGradeFailsLane` fails only `fail` at lines 326-328. The runner
only logs `XPASS ALERT` at `internal/ensigncycle/claude_live_runner_test.go:216-228`.

Therefore, the command does not fail as strict XPASS. Lines 1036-1037 require
an outcome that the current harness cannot produce.

Define one exact proof rule. Either make XPASS fail this lane, or require a
successful command with an exact XPASS record and empty semantic codes.

Then remove only the two bindings and mirrored entries. Both reruns must pass
with the same durable receipt and open-gate evidence.

### Surface verification

The table has exactly 40 files. Its arithmetic is correct:

- 4,685 insertions;
- 121 deletions;
- 4,806 gross lines;
- 4,564 net lines;
- 480 lines of gross tolerance.

The surface is large for a CLI-only feature. The estimate is acceptable only as
a ceiling after the Material design gaps close.

No product implementation is authorized by this review.

## Staff review request: ten-file happy path

The preceding review evaluates the superseded 40-file proposal. It remains as
durable review history.

The current proposal is `Smallest CLI happy-path rescope`. It has ten files and
461 gross lines.

Please review that rescope against its five exact questions. No product
implementation is authorized before staff accepts it.

## Second correction: consumed receipt completion gate

Date: 2026-08-10

Review state: **HOLD — staff review requested**.

This correction applies to `Smallest CLI happy-path rescope`. It supersedes any
conflicting completion or line-estimate text in that section.

No product implementation is authorized.

### Material finding and correction

An active `consumed` receipt proves native worker start. It does not prove
worker completion.

Before a stage-changing `status --set`, the handler must inspect the active
receipt for the current entity and stage.

The corrected rules are:

- `pending` refuses the stage change before any entity write.
- `armed` refuses the stage change before any entity write.
- `consumed` always calls the existing complete-and-committed stage-report
  predicate.
- A false predicate result refuses the stage change and preserves entity bytes.
- A true predicate result permits the existing stage transition.
- No active receipt preserves current legacy behavior.

The handler calls the predicate even when `worktree` is empty. It does not use
worktree presence as a completion shortcut.

The predicate must check the exact entity path and current stage. It requires a
complete latest stage report and a clean tracked entity path in local `HEAD`.

`--force` cannot bypass this consumed-receipt predicate.

### Stamped and unstamped value

Every fresh Claude CLI or Codex CLI build creates the pending envelope. This
rule applies with and without `--stamp`.

`--stamp` keeps its current stage-entry work. It is not an acknowledgment
prerequisite and does not define the supported user value.

An unstamped fresh build gets the same pending, armed, consumed, report, and
stage-transition rules.

### Real audit-ref oracle

The implementation writes these immutable audit refs during the real state
transitions:

```text
refs/spacedock/dispatch-ack-audit/<entity-id>/<stage>/<epoch>/pending
refs/spacedock/dispatch-ack-audit/<entity-id>/<stage>/<epoch>/armed
refs/spacedock/dispatch-ack-audit/<entity-id>/<stage>/<epoch>/consumed
```

Each audit ref points to the actual active-ref blob from that transition. The
active compare update and matching audit-ref create use one `git update-ref`
transaction.

The live oracle reads these Git refs directly. It does not construct an audit
record after the run.

For each host, the oracle must prove this exact order:

1. A real fresh build writes `pending` with epoch E.
2. Supported `PreToolUse` writes `armed` with the same epoch E.
3. Supported `SubagentStart` writes `consumed` with epoch E.
4. The consumed blob contains one non-empty native worker ID.
5. The implementation report becomes complete and committed.
6. A later stage-changing status command succeeds.
7. The validation gate is prepared, committed, presented, and left open.

The pending, armed, and consumed refs must contain the same entity ID, entity
path, stage, host, and epoch.

The armed ref must contain one supported tool-use ID. The consumed ref must
contain one supported native worker ID.

The status command time must be later than the consumed audit transition. The
validation gate commit must be later than the complete implementation report
commit.

The oracle reads no Claude transcript, Codex rollout, final message, or private
host activity file to establish acknowledgment or completion.

### Corrected negative proof

Add these focused cases to the existing ten-file test plan:

1. `consumed` plus an empty `worktree` and no report refuses byte-clean.
2. `consumed` plus an empty `worktree` and an incomplete report refuses.
3. `consumed` plus an empty `worktree` and an uncommitted report refuses.
4. `consumed` plus an empty `worktree` and a complete committed report passes.
5. Each refusal also fails with `--force` and preserves entity bytes.
6. Stamped and unstamped fresh builds create equivalent pending refs.
7. The three real audit refs share one epoch and one ordered transition chain.

### Corrected XPASS-green proof

Keep both default-headless bindings for the first live runs. Each exact target
must exit zero and log `XPASS ALERT` with no semantic failure code.

The lane-green XPASS artifact must contain the real three-ref chain, one worker
ID, a complete committed implementation report, a later stage change, and the
open validation gate.

Then remove only the Claude Sonnet and Codex bindings. Remove only their two
mirrored registry expectations in the same binding-only commit.

Each normal rerun must exit zero and report PASS. It must prove the same audit
chain, report, stage order, and open gate.

### Recalculated ten-file estimate

The exact file set remains unchanged:

| File | Insertions | Deletions |
|---|---:|---:|
| `.claude-plugin/plugin.json` | 1 | 0 |
| `hooks.json` | 36 | 2 |
| `internal/dispatchack/ack.go` | 148 | 0 |
| `internal/dispatchack/ack_test.go` | 138 | 0 |
| `internal/dispatch/build.go` | 30 | 6 |
| `internal/dispatch/dispatch.go` | 35 | 4 |
| `internal/status/handlers.go` | 24 | 4 |
| `internal/ensigncycle/claude_live_runner_test.go` | 34 | 24 |
| `internal/ensigncycle/shared_live_runner_test.go` | 2 | 2 |
| `internal/contractlint/live_registry_reconciliation_test.go` | 1 | 1 |

The corrected estimate is 449 insertions and 43 deletions. It is 492 gross
lines and 406 net lines.

The exact set remains ten files. The 500-gross hard cap remains in force.

### Staff review request

Please confirm these three points:

1. Does the consumed-state predicate close the empty-worktree completion gap?
2. Does the real three-ref oracle prove the exact host acknowledgment order?
3. Can the corrected implementation stay within ten files and 492 gross lines?

No product bytes or live CI can run before staff accepts this correction.
