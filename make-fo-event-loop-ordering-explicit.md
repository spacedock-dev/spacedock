---
title: Make FO event-loop ordering and idle wait explicit
status: ideation
source: "Captain follow-up after the 2026-08-03 durable-decisions execution-gap diagnosis."
started: 2026-08-03T16:00:51Z
completed:
verdict:
score: 0.98
worktree:
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: ej9kwkvw94w6rh6n5ek7qrbf
gates:
    version: 1
    records:
        - id: gate:ej9kwkvw94w6rh6n5ek7qrbf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ej9kwkvw94w6rh6n5ek7qrbf-backlog-1
              briefing:
                id: briefing:ej9kwkvw94w6rh6n5ek7qrbf:backlog:attempt-1:revision-1
                digest: sha256:4dbaa852a23e16b0c2c732a3fac6956ef5f95db8c603f428b1636a5983d3064b
                request-digest: sha256:2716adf5a2f9014015e76146c012c4feb29a71001469f747748059142ecf26cb
                room-ref: ./make-fo-event-loop-ordering-explicit/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ej9kwkvw94w6rh6n5ek7qrbf:backlog:1
                briefing: briefing:ej9kwkvw94w6rh6n5ek7qrbf:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T15:55:01.310359Z"
                decision: approve
                reason: Captain conn approves ideation. SO concurs with parallel ideation for EJ and G3; preserve EJ-before-G3 implementation landing because both define the shared FO dispatch/merge contract.
              application:
                target-stage: ideation
                state: consumed
---

Make the First Officer event loop mechanically explicit so a dispatch-only empty result cannot hide merge recovery, ready gates, or a required idle/reconcile pass. The task preserves the existing state and runtime boundaries; it makes the ordering observable and testable.

## Problem

The current contract has the required branches, but an FO can run a filtered `status --next` query, see `dispatchable: []`, and report idle without processing `mod-block`/PR recovery, the separate `ready_gates` surface, the one-shot idle hook and reconcile retry, or the unresolved-worker wait rule. This loses the distinction between no dispatchable worker, a pending merge action, a gate awaiting Captain input, and a genuinely idle session.

## Proposed approach

Specify one event-loop trace and bind every runtime adapter to it. Each iteration drains worker messages when supported, processes every `mod-block` and PR action before querying dispatchables, handles `ready_gates` as gate or merge actions rather than worker dispatches, and then runs `status --next`. After the first empty result, run the idle hook exactly once, reconcile the roster, and query `status --next` once more. Install async wait only when the worker set is unresolved and there is no dispatchable entity, presentable gate, mod action, or other state work; a completed, errored, or absent worker set never qualifies. Preserve keep-moving for unrelated ready tasks and do not add a new state field or scheduler command.

## Out of scope

- Changing the commissioned workflow definition or stage taxonomy.
- Adding a new status command, daemon, watchdog, or generic resolver worker.
- Auto-resolving PR conflicts, force-pushing, or changing merge/approval authority.
- Rewriting task-specific findings or changing the Codex, Claude, or Pi transport APIs.

## Acceptance criteria

**AC-1 (VALUE) — A ready action cannot be hidden by a dispatch-only empty result.**
Verified by: a fixture-backed event-loop trace with one `mod-block`/PR row, one `awaiting-captain` or `approved-awaiting-merge` gate, and one genuinely dispatchable task; the recorded order shows mod/PR and gate handling before `status --next`, and the unrelated ready task still dispatches.

**AC-2 — Empty dispatch results trigger one idle/reconcile retry and a truthful stop.**
Verified by: a command-log fixture that returns an empty first `status --next`, releases a task from the idle hook or reconcile, and asserts one idle hook, one reconcile, one second `status --next`, and dispatch of the newly unblocked task; the unchanged-empty variant asserts the explicit no-dispatchable stop reason.

**AC-3 — Async wait is reserved for unresolved workers.**
Verified by: Codex runtime fixtures covering active unresolved, completed, errored, and absent worker sets; only the active unresolved case with no other work may emit `wait_agent(timeout_ms: 300000)`, and no other case emits a wait.

**AC-4 — Gate readiness is not conflated with worker dispatchability.**
Verified by: boot/status fixtures that expose `awaiting-captain`, `approved-awaiting-advance`, and `approved-awaiting-merge` rows, and assert presentation, consume/advance, and merge routing respectively without a worker spawn or a false idle report.

**AC-5 — One blocked task does not park unrelated work.**
Verified by: a multi-entity fixture with one conflict/mod-blocked task and two independent ready tasks; the trace shows the blocked task held with its evidence while both independent tasks reach their declared dispatch or gate action.

## Expected surface and semantic boundaries

Expected change is limited to the First Officer dispatch/event-loop contract and its behavior fixtures: `skills/first-officer/references/first-officer-shared-core.md`, `skills/first-officer/references/fo-dispatch-core.md`, the host runtime adapter notes needed to bind the wait predicate, and focused event-loop/recorded-workflow tests. Estimate 4–8 files and +120/-20 lines, with a 2x tolerance. Allowed semantic change: ordering and stop-reason observability only. Stored state, command grammar, gate authority, worker transport, and stage behavior remain unchanged.

## Test plan

Start with the smallest falsifying fixture: replay the current false-stop trace and assert that a filtered `status --next=[]` cannot terminate the iteration while a mod-block, ready gate, or unresolved worker remains. Add deterministic command-log tests for each branch and one runtime-backed Codex wait test. Run the focused status/dispatch/ensigncycle suites, then `go test ./...`, `go test ./... -race`, formatting, and the required host live lanes when contract/runtime adapter files change.

## Stage-specific test gates

- Ideation must replay the falsifying fixture and record the exact ordered trace before selecting a landing shape.
- Implementation must keep the trace assertions independent of contract prose and preserve the existing state and command bytes.
- Validation must cover the empty-first-query retry, gate/merge separation, unresolved-worker wait predicate, and keep-moving matrix, then run full, race, formatting, and relevant live evidence.

## Stage Report: ideation

### False-stop replay and ordered trace

The current checkout provides a direct falsifier for treating an empty dispatch query as idle (run 2026-08-04):

```text
status --where 'mod-block !=' --json --fields id,slug,mod-block
  26n/headless-recorded-gate-stop-stage-coherence  merge:pr-merge
  32/boot-identify-multi-workflow-llm-retry-friction merge:pr-merge
  93/document-dispatch-entity-path-base              merge:pr-merge
status --next --json --fields id,slug
  {"dispatchable":[]}
status --boot --identify --json
  {"dispatchable":[],"ready_gates":[],"pr_state":{"status":"local",...},"team_state":{"present":true,...}}
```

The first query exposes three pending PR/merge actions; the second query alone therefore cannot be a stop decision. The boot-identify status and recorded fixtures provide the complementary gate case: `internal/status/boot_identify_test.go` passes with `ready_gates` rows `mf` (`awaiting-captain`), `r4` (`awaiting-captain`), `wd` (`withdrawn-awaiting-prepare`), and `2n` (`approved-awaiting-merge`) alongside dispatchable `dispatch-me`. The terminal-gate test keeps `approved-awaiting-merge` in `ready_gates` through `gate consume` until merge delivery, proving that consume is not worker dispatch. The Codex wait fixture `TestCodexWaitAgentSteeringEvidence` also passes: after spawn and active-scope exhaustion it records `wait_agent(timeout_ms:300000)`, retries after captain steering, then verifies the durable report after `final_status`.

The implementation-gate trace must be explicit about both the existing baseline and the required decision order. For each iteration:

1. If a roster is addressable, reconcile it; then drain addressable inbound worker messages. If the runtime has no roster/message primitive, record that adapter boundary rather than inferring a clean roster.
2. Query `status --where "mod-block !="`; process every blocking mod/PR before dispatch. A merged PR records its sentinel and runs merge guard; an open PR remains held; a closed PR asks the Captain; an unavailable check is skipped without clearing the block. No new worker is dispatched by this branch.
3. Read the boot/status `ready_gates` projection and route each gate as state work: present `awaiting-captain`, prepare a successor for `withdrawn-awaiting-prepare`, consume/advance `approved-awaiting-advance` (with successor and observed spawn), and deliver/merge `approved-awaiting-merge` (never spawn). Any unresolved gate prevents an idle claim.
4. Run `status --next`; dispatch every returned row, requiring worker-spawn evidence (a successful dispatch build is not spawn evidence). Keep moving independent ready rows when one row is blocked.
5. Only after the first empty `status --next`, run the idle hook exactly once, reconcile when the adapter supports it, and run one second `status --next`. Dispatch newly released work; if the second result is still empty, report an explicit `no-dispatchable` stop only when there is no mod/PR action, gate action, or other state work.
6. In Codex, emit `wait_agent(timeout_ms:300000)` only when an active unresolved worker remains and the post-retry scope has no dispatchable, gate, mod, or state work. A completed, errored, or absent worker never qualifies. A timeout retries the same epoch; a completion notification starts durable report/state verification and another iteration. Pi has no `wait_agent`; Claude's reconcile and idle behavior follows its bound team identity.

This order makes the false-stop trace observable: the current status-only path is `mod-block? -> status --next=[] -> idle/stop`, while the required path is `message drain -> mod/PR -> ready_gates -> status --next -> (idle once + reconcile + status --next retry) -> explicit stop or guarded wait`.

### Acceptance-criteria and boundary reconciliation

The ACs and existing boundaries are consistent, with the following implementation fixtures required:

| AC | Fixture assertion | Existing boundary preserved |
| --- | --- | --- |
| AC-1 | Mixed mod/PR + gate + dispatchable command log proves mod and gate actions precede `status --next`, while the independent row dispatches. | `status` is read-only; PR merge guard and gate lifecycle own mutation. |
| AC-2 | First empty result emits one idle hook, one reconcile, one retry; released work dispatches; unchanged empty reports `no-dispatchable`. | `hooks.run` remains caller-owned and exact-once; no new scheduler command. |
| AC-3 | Active-unresolved/completed/errored/absent matrix permits wait only for active-unresolved with an otherwise empty scope. | Codex uses `list_agents`/`wait_agent`; completion attribution is mailbox/task-path plus durable state, never a handle alone. |
| AC-4 | `awaiting-captain`, `approved-awaiting-advance`, and `approved-awaiting-merge` route to presentation, consume/advance, and merge respectively, with no false worker spawn/idle. | `ready_gates` and gate bytes remain authoritative; terminal consume remains pending until merge delivery. |
| AC-5 | One blocked row does not prevent two independent ready rows from reaching dispatch/gate actions. | Keep-moving semantics and worker write scopes remain unchanged. |

Focused status, gate-lifecycle, and Codex wait tests already pass. Add a host-neutral command-log/event-loop fixture first, then adapter-specific assertions; do not invent a shipped `dispatch next-action` command. Run the focused suites, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`; run host live lanes only if runtime adapter notes change.

### Unresolved risk and recommendation

The separate unhandled-completion-ledger risk remains: a final-status/mailbox completion can arrive while the FO is handling another request, and no durable completion feed currently records that it was handled. Ordering the available drain and reconcile operations reduces the window but cannot recover a lost completion or make a reconcile without session identity authoritative; the live `dispatch reconcile --workflow-dir docs/dev` probe reported git/filesystem drift only because no team identity was resolved. Keep that ledger/reclaim work separate, and make wait fail-safe when roster identity is absent.

Recommendation for implementation-gate entry: **APPROVE ideation**. Land a contract-first numbered event-loop trace plus a small host-neutral command-log reducer/fixture and the Codex wait matrix. Make the mixed falsifying fixture red against the current `status --next=[]` stop before updating the contract and adapter notes. Preserve all current state fields, command bytes, transport APIs, gate authority, and stage behavior; the intended semantic change is ordering and stop-reason observability only.
