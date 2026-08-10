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

Codex can receive a valid implementation spawn envelope and continue without a
native worker handle. The First Officer can then change the status to validation.

The current prose rule is not deterministic. Identical fixture, prompt, adapter,
and envelope bytes produced two native-spawn passes and one empty-wait failure.

The binary records the stage start, but it does not record a Codex dispatch
handle. It also does not record a completion receipt for that handle. Therefore,
the status command cannot distinguish completion from a skipped spawn.

## End-user value

A Codex user gets one real implementation worker before validation starts. The
workflow binds the worker handle and its completion receipt to the implementation
stage. The existing default-headless journey passes without XFAIL.

## Constraints

- Keep the existing `spacedock dispatch build` command surface and stored workflow formats unless ideation proves a change is necessary.
- The binary owns the acknowledgment boundary. Do not add another prose guard or simulator.
- Status advance and `--advance` must fail closed when a spawn envelope lacks a matching native handle and completion receipt.
- Preserve Sonnet, Pi, gate authority, and existing worker ownership behavior.
- Remove only this task's Codex default-headless XFAIL after exact passing product evidence.

## Current boundaries

The current path has five separate boundaries:

1. `spacedock dispatch build --stamp` writes `started` and emits the Codex envelope.
2. `spawn_agent` returns the canonical Codex task path as its native handle.
3. The Codex mailbox sends a final notification from that same task path.
4. The worker writes and commits its implementation report.
5. `spacedock status --set` can change the stage after the report exists.

The binary owns boundaries 1 and 5. Codex owns boundaries 2 and 3. The report
owns boundary 4. No durable value currently joins all five boundaries.

The 98a cycle-4 evidence shows this gap. The exact journey used the same spawn
envelope in each run. Two runs contained a native spawn and matching completion.
One run contained an empty wait, then crossed into validation without a spawn.

## Spike result

The ideation spike used the live Codex agent surface. `spawn_agent` returned the
canonical task path
`/root/commander/spacedock_ensign_n28423efmj_ideation/n284_receipt_spike`.
The final mailbox notification used the same sender path. `list_agents` then
reported that path as completed with `N284_RECEIPT_SPIKE_COMPLETE`.

This exercise proves that the native surface exposes one stable opaque handle at
spawn and completion. It also proves the remaining gap: no current Spacedock
command consumes that handle. The selected design adds that command and makes
the existing status boundary consume its durable result.

## Proposed approach

Add one binary-owned Codex dispatch receipt. Keep the build envelope unchanged.

1. For a Codex `dispatch build --stamp`, write a pending `dispatch-receipt`
   record in the entity before the existing path-scoped commit and state sync.
2. Store the entity ID, stage, host, epoch, state, and native handle in the
   record. The first state is `pending`, and the handle is empty.
3. After `spawn_agent` returns, run `spacedock dispatch acknowledge` with
   `--spawned` and the returned task path. The binary stores the handle.
4. After the final mailbox notification, run the same command with `--completed`
   and the sender task path. The binary requires an exact handle match.
5. Refuse a status change away from the recorded stage until the receipt is
   complete and the current stage report is complete and committed.
6. Refuse `dispatch build --advance` when the preceding Codex receipt is absent,
   pending, spawned without completion, stale, or bound to another handle.
7. Commit and sync each receipt update with the existing entity commit unit.
   A sync failure returns a nonzero exit and keeps the local receipt durable.
8. Do not create a receipt for Claude, Pi, or a legacy entity without a Codex
   receipt. Their current behavior stays unchanged.

The stored-format change is necessary. Separate command invocations cannot share
process memory, and a stage report contains no native handle. The optional
`dispatch-receipt` map keeps the receipt with the entity and its existing sync
unit. The record is replaced only after its prior epoch is complete.

The simplest alternative was another imperative sentence after `dispatch build`.
The 98a evidence disproves that alternative. A component-only receipt API was
also rejected because status and `--advance` can still bypass it. Parsing a
private Codex rollout was rejected because a normal Spacedock process does not
own that host file or its lifecycle.

This mechanism serves AC-1. The status and `--advance` guards make the receipt
part of the end-to-end value path instead of an optional component.

## Released workflow

The normal user still runs Codex through Spacedock. The First Officer builds the
dispatch and calls the native spawn tool. It then acknowledges the returned
handle. After Codex sends completion, the First Officer acknowledges that same
handle. Only then can validation start.

The command grammar adds one subcommand:

```text
spacedock dispatch acknowledge --workflow-dir DIR --entity REF --stage STAGE --spawned HANDLE
spacedock dispatch acknowledge --workflow-dir DIR --entity REF --stage STAGE --completed HANDLE
```

The existing `dispatch build` output does not change. The existing status command
grammar does not change.

## Acceptance criteria

**AC-1 (VALUE) — Codex implementation dispatch completes before validation.**

The exact Codex `default-headless-gate-stop` journey records one native
implementation spawn. The receipt contains its matching handle and completion
before validation. The implementation report contains DONE, and the journey
passes without XFAIL.

Verified by: run the exact live target on the candidate after the focused
controls pass. The test fails if spawn is absent, the handles differ, completion
follows validation, the report lacks DONE, or the gate opens early.

**AC-2 — Missing acknowledgment fails closed.**

An envelope followed by an empty wait, a report read, a status change, or
`--advance` cannot enter validation without the matching handle and completion
receipt.

Verified by: command tests start from each receipt state. Each status or
`--advance` attempt must exit nonzero and leave entity bytes unchanged.

**AC-3 — The mechanism is end to end.**

The Codex adapter acknowledges the native handle, the binary stores it, and the
status path consumes it. The exact live journey passes through this path.

Verified by: a fixture drives build, spawned acknowledgment, completed
acknowledgment, report commit, and status change through the CLI. Removing any
step makes the fixture fail.

**AC-4 — A stale or different handle has no authority.**

A completion from another task path or dispatch epoch cannot complete the
current stage.

Verified by: focused tests send a different handle and replay the prior epoch.
Both commands must exit nonzero without changing the receipt.

**AC-5 — Other runtimes keep their current behavior.**

Claude, Pi, and entities without a Codex receipt do not acquire the new status
guard.

Verified by: control fixtures run the same transitions without a Codex receipt.
Their output, exit code, and entity state remain unchanged.

## Explicit non-goals

- The task does not make the binary call `spawn_agent`.
- The task does not parse Codex private session or rollout files.
- The task does not add a workflow driver or change scheduler priority.
- The task does not change Claude, Pi, gate authority, or worker ownership.
- The task does not add receipts for reused workers after `--advance`.
- The task does not change the dispatch envelope schema.
- The task does not remove the Pi default-headless XFAIL.

## Expected surface

The implementation changes exactly these 11 files:

- `internal/status/dispatch_receipt.go`: 150 insertions.
- `internal/status/dispatch_receipt_test.go`: 190 insertions.
- `internal/status/handlers.go`: 22 insertions.
- `internal/dispatch/dispatch.go`: 48 insertions and 2 deletions.
- `internal/dispatch/stamp.go`: 30 insertions and 3 deletions.
- `internal/dispatch/acknowledge_test.go`: 140 insertions.
- `skills/first-officer/references/codex-first-officer-runtime.md`: 8 insertions and 4 deletions.
- `skills/first-officer/references/fo-dispatch-core.md`: 5 insertions and 3 deletions.
- `docs/runtime-support.md`: 10 insertions and 2 deletions.
- `internal/ensigncycle/shared_live_runner_test.go`: 1 insertion and 1 deletion.
- `internal/contractlint/live_registry_reconciliation_test.go`: 1 insertion and 1 deletion.

The estimate is 605 insertions and 16 deletions. The total is 621 gross lines
and 589 net lines. The tolerance is 75 gross lines, with no additional files.

The command grammar adds `dispatch acknowledge`. The entity format adds the
optional `dispatch-receipt` map. Authority labels and the build envelope stay
unchanged. Runtime behavior changes only for a recorded Codex dispatch and its
immediate stage exit.

## Documentation diff

Current Codex runtime text says that a successful build is not a dispatch and
that the First Officer must call `spawn_agent`.

The implementation replaces that paragraph with this behavior:

> After `spawn_agent` returns, acknowledge its task path with `spacedock dispatch
> acknowledge --spawned`. After the matching final notification, acknowledge the
> same task path with `--completed`. Do not change status or use `--advance`
> before the completed receipt exists.

`docs/runtime-support.md` will add this user-visible statement:

> Codex dispatch stores a pending receipt with the entity. The native spawn and
> final notification must acknowledge the same task path before the next stage.

## Test plan

1. Add focused status tests for pending, spawned, completed, mismatched, and
   stale receipts. These tests serve AC-2 and AC-4.
2. Add dispatch command tests for valid transitions, idempotent retry, invalid
   order, sync failure, and `--advance` refusal. These tests serve AC-2 and AC-3.
3. Add an end-to-end CLI fixture. Drive build, both acknowledgments, report
   commit, and status change. This fixture serves AC-1 and AC-3.
4. Run Claude, Pi, and legacy controls. These controls serve AC-5.
5. Run the focused contract and registry checks. Then run `gofmt -w ./cmd
   ./internal`, `go test ./...`, and `go test ./... -race`.
6. Run the exact Codex live target with this task's XFAIL still present. The
   repaired target must report XPASS, which is a strict failure.
7. Remove only the Codex binding owned by `n28423efmj358m5av61z2fxx`. Update
   the reconciliation row, then rerun the exact target for a normal pass.

The focused fixture has low runtime cost. The full and race suites have normal
repository cost. The exact Codex live target has model and runtime cost, so it
runs last on the exact candidate.

## Stage Report: ideation

- DONE: Trace the current Codex dispatch envelope, native spawn handle, completion receipt, report, and status-advance boundaries.
  The boundary trace identifies five steps and shows that no durable value joins the binary and native Codex steps.
- DONE: Reproduce or inspect the exact 98a pass/failure evidence and identify the smallest binary-owned acknowledgment mechanism.
  The 98a evidence has two native-spawn passes and one empty-wait failure from identical envelopes. One entity receipt closes the gap.
- DONE: Spike the mechanism end to end through the real Codex path before selecting the design. Do not ship a component-only API.
  The native spike returned one task path at spawn, completion, and roster read. The design makes status consume that receipt.
- DONE: Define the released user workflow, end-user value, exact semantic boundary, and explicit out-of-scope items.
  The body defines the unchanged launch flow, the Codex-only stage guard, and seven explicit non-goals.
- DONE: Write acceptance criteria with exact falsifiers, including missing-handle fail-closed behavior and exact live XFAIL removal.
  Five criteria name command, fixture, control, stale-handle, and exact live falsifiers.
- DONE: Declare exact files, gross additions/deletions, net lines, tolerance, and any command or stored-format impact.
  The surface is 11 files, 621 gross lines, 589 net lines, and 75 gross lines of tolerance.
- DONE: Provide a test plan using focused controls first and exact Codex live proof last.
  The plan runs focused receipt controls, CLI proof, full checks, strict XPASS, binding removal, and normal live pass.
- DONE: Commit and push the complete ideation report in Simplified Technical English.
  This report uses short active sentences and consistent receipt terminology. The state commit and push provide durable proof.

### Summary

The selected design adds one durable Codex dispatch receipt and one acknowledgment
command. The native spike proves that spawn and completion expose the same task
path. Status and `--advance` consume the receipt, so a component-only API cannot
satisfy the task.
