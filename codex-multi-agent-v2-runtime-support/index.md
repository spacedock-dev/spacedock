---
title: Codex multi_agent_v2 runtime dispatch support
status: implementation
source: "FO live dispatch observation (2026-06-20): while dispatching status-validate-determinism under Codex multi_agent_v2, the current Codex runtime adapter no longer matches the live tool surface. v2 exposes spawn_agent(task_name,message,fork_turns), list_agents, wait_agent(timeout_ms), send_message, followup_task, and interrupt_agent. The shipped adapter still names send_input, wait_agent(handle), and hyphenated dispatch-build names as direct handles."
score:
started: 2026-06-20T03:58:46Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-codex-multi-agent-v2-runtime-support
issue:
sprint:
sprint-readiness:
id: f2r8cnyxj9pf24xrsf71szb0
---

# Codex multi_agent_v2 runtime dispatch support

## Problem

The Codex first-officer and ensign runtime references describe the older collab surface:

- initial dispatch uses `spawn_agent`, but `dispatch build` emits hyphenated `name` values such as `spacedock-ensign-status-validate-determinism-implementation`; live v2 `spawn_agent.task_name` rejects hyphens and accepts only lowercase letters, digits, and underscores
- feedback/reuse/shutdown text names `send_input`, but the live surface exposes `send_message` and `followup_task`
- foreground waiting text says `wait_agent(handle)`, but the live surface exposes global `wait_agent(timeout_ms)` and captain messages interrupt the wait without killing the worker
- the no-wait probe still did not prove autonomous idle wake-up; completion surfaced after later root activity, so the correct classification is queued/activity-driven delivery
- `list_agents` is now a usable roster-like read for active and completed task paths

The existing Codex runtime adapter needs forward-compatible multi_agent_v2 host bindings that capture these semantics without adding a separate Codex v2 runtime file or breaking existing v1/live-test assumptions before they are intentionally retired.

## Proposed approach

Update the existing Codex runtime references with an explicit **Codex multi_agent_v2** section rather than adding a separate runtime variant. Codex is still one host from the Spacedock workflow perspective; the incompatible piece is the live collab tool surface exposed inside newer Codex sessions. A separate `codex-v2` host would duplicate front-door, install, dispatch-build, and live-lane plumbing while making mixed-version support harder to reason about. The safer compatibility path is:

1. Keep `host: codex` in `spacedock dispatch build`.
2. Add a narrow adapter mapping from dispatch-build output to live v2 tool arguments.
3. Version-gate or retire stale v1 assumptions in contractlint and live assertions.
4. Preserve legacy transcript parsers only behind fixtures that explicitly name the old surface.

No spike is needed for the core v2 mapping: this task is seeded by a live Codex multi_agent_v2 dispatch observation with rejected hyphenated task names, successful underscore retry, `list_agents` roster reads, global wait interruption, and successful `send_message` delivery. The remaining risky behavior is shutdown semantics; implementation must prove the shutdown binding with a small live or fixture-backed probe before blessing any concrete tool.

## Host-binding language

Implement Codex multi_agent_v2 as host bindings on the existing `«fn»`-style dispatch contract, not as a parallel capability registry or a separate `codex-v2` runtime file:

- `«worker.spawn»` → Codex multi_agent_v2: call `spawn_agent(task_name,message,fork_turns)` for initial dispatch. Parse `spacedock dispatch build` JSON, pass emitted `prompt` byte-for-byte as `message`, and omit unsupported `description`, `subagent_type`, and `model` arguments. Use `fork_turns` only when the Codex adapter deliberately wants inherited context; file-pointer initial dispatch should normally use no inherited turn context.
- `«worker-identity»` → Codex multi_agent_v2: helper-emitted `name` remains the host-neutral worker identity source, but the live spawn argument must be sanitized to a v2-valid `task_name` by replacing non `[a-z0-9_]` runs with `_` and preventing collisions. The adapter must retain enough mapping from sanitized task path back to entity slug/stage/cycle to classify cohorts and teardown targets.
- `«addressable-worker»` → Codex multi_agent_v2: realized by the mailbox tools for an existing worker. Use `send_message(target,message)` for non-turn-triggering delivery such as context notes or preservation-aware shutdown text. Use `followup_task(target,message)` for reuse, feedback, rework, or any instruction that must trigger a worker turn. These are two concrete realizations of the same `«addressable-worker»` binding, not separate named capabilities. Legacy `send_input` belongs only to explicitly versioned pre-v2 fixtures or references.
- `«completion-signal»` → Codex multi_agent_v2: the ensign still completes by sending the concise final message observed as a final-status mailbox update. `wait_agent(timeout_ms)` is a global, interruptible foreground wait for mailbox updates, not a handle-scoped wait. Captain input or other operator activity cancels the wait call but does not fail, close, or redispatch the worker; if the worker is still unresolved after control returns, the FO reinstalls the wait when waiting is again the next useful idle action. No-wait completion remains classified as queued/activity-driven unless a clean idle probe proves autonomous wake-up.
- `«roster-reconcile»` → Codex multi_agent_v2: use `list_agents(path_prefix?)` as the roster read for active/completed Codex task paths. It can support stale-cohort classification, debugging, and terminal cleanup targeting, but durable workflow state remains authoritative for workflow advancement.
- `«worker.shutdown»` → Codex multi_agent_v2: unresolved until probed. Do not bless `interrupt_agent` as the shutdown binding yet. Implementation should first prove whether `interrupt_agent` terminates, pauses, or leaves a worker addressable, and should keep cooperative preservation text separate from any hard interruption.

## Compatibility risks

- Existing Codex live fixture parsers bind reuse to `send_input` events and thread IDs. A blind text replacement would make old fixtures fail for the wrong reason, so tests need either a versioned v1 parser path or fixtures marked `codex_multi_agent_v2`.
- `wait_agent(handle)` assumptions are embedded in prose and watchdog naming. v2's global wait can return an unrelated worker or an interruption, so entity attribution must come from mailbox content, task path, or durable state.
- Sanitized `task_name` values can collide if two helper names differ only by punctuation. The adapter should either include a deterministic suffix when sanitation changes the name or test that the known dispatch-build name format remains collision-free for `(worker_key, slug, stage, cycle)`.
- `description` and `model` are mandatory host-neutral helper outputs today, but unsupported by v2 `spawn_agent`. The Codex adapter must clearly say "helper emits them, v2 spawn omits them" so the host-neutral core does not weaken Claude/Pi forwarding guarantees.
- `send_message` may not wake an idle/completed worker for rework. Reuse and feedback must use the `followup_task` realization of `«addressable-worker»` unless a probe proves `send_message` starts the needed turn.
- `interrupt_agent` may be stronger than cooperative shutdown and could lose uncommitted worker state. `«worker.shutdown»` must remain unresolved until probed; implementation must not silently equate interrupt with cooperative shutdown.

## Acceptance criteria

**AC-1 - Initial Codex v2 dispatch maps helper output to the live spawn surface.**  
Verified by: a test or live fixture showing `dispatch build` output with a hyphenated `name` is converted to a valid v2 `task_name` while forwarding the helper `prompt` unchanged as `message`, omitting unsupported `description` and `model` fields.

**AC-2 - Codex v2 continuation and steering use the live message tools.**  
Verified by: adapter text and tests distinguish `send_message` context delivery from `followup_task` turn-triggering reuse/feedback, and remove or explicitly version-gate `send_input` from v2 instructions.

**AC-3 - Codex v2 waiting semantics match observed behavior.**  
Verified by: adapter text and live/fixture evidence state that `wait_agent(timeout_ms)` is global, captain input interrupts wait without failing the worker, and no-wait completion remains queued/activity-driven unless an autonomous wake probe proves otherwise.

**AC-4 - Terminal/supersede shutdown has a probed v2 binding.**  
Verified by: adapter text leaves `«worker.shutdown»` unresolved until a live or fixture-backed probe proves the concrete v2 behavior; implementation does not bless `interrupt_agent` before that proof and names residual risk after the binding is chosen.

**AC-5 - Old Codex runtime assumptions are preserved or deliberately retired.**  
Verified by: existing Codex live/shared tests either remain green under a versioned v1 path or are updated with v2-specific expectations; stale assertions around `send_input` and `wait_agent(handle)` are not left ambiguous.

## Evidence to seed ideation

- `spawn_agent` rejected helper name `spacedock-ensign-status-validate-determinism-ideation` with `agent_name must use only lowercase letters, digits, and underscores`.
- Retrying as `spacedock_ensign_status_validate_determinism_ideation` succeeded with the same helper prompt.
- `list_agents` showed `/root`, completed probe workers, and the running Spacedock ensign task.
- Captain messages interrupted `wait_agent`, but the worker remained running and could later complete.
- A no-wait window did not wake root autonomously; completion became visible after a later `list_agents` call and mailbox notification.
- `send_message` successfully delivered a non-interrupting note to the running worker, which then committed and finalized.

## Documentation delta to apply during implementation

Implementation should update the existing Codex adapter references in `skills/first-officer/references/codex-first-officer-runtime.md`, `skills/ensign/references/codex-ensign-runtime.md`, `skills/feedback-rejection-flow/SKILL.md`, and `docs/dev/codex-idle-notification-probe.md`. It should add Codex multi_agent_v2 bindings on the existing `«fn»` contract rather than creating a new Codex v2 runtime file or a parallel capability registry.

Replace unversioned Codex reuse wording like:

```text
Route feedback or continuation to an existing Codex worker with `send_input`.
```

with v2-aware wording like:

```text
On Codex multi_agent_v2, `«addressable-worker»` is realized by `send_message(target,message)` for non-triggering notes and `followup_task(target,message)` for reuse, feedback, or rework that must start a worker turn. Legacy `send_input` belongs only to explicitly versioned pre-v2 fixtures or references.
```

Replace foreground-wait wording like:

```text
The FO calls `wait_agent(handle)` as the next useful idle action.
```

with:

```text
On Codex multi_agent_v2, `«completion-signal»` may use global `wait_agent(timeout_ms)` only after workflow and gate work is exhausted. Captain input cancels the wait call but not the worker; if the worker remains unresolved, reinstall the wait when waiting is again the next useful idle action. A wait return must be attributed by mailbox content, task path, roster state, or durable workflow state, not by a handle argument.
```

Add a Codex v2 dispatch paragraph:

```text
For Codex multi_agent_v2, bind `«worker.spawn»` to `spawn_agent(task_name,message,fork_turns)`: sanitize emitted `name` through `«worker-identity»` to a lowercase digit/underscore `task_name`, pass emitted `prompt` unchanged as `message`, omit unsupported `description`, `subagent_type`, and `model`, and record the sanitized task path as the live worker handle.
```

## Test plan

- **Dispatch-build mapping tests:** add a Go fixture around the existing Codex adapter's `«worker.spawn»` / `«worker-identity»` mapping from helper JSON to v2 spawn call input. Assert that a helper name such as `spacedock-ensign-status-validate-determinism-implementation` becomes a valid underscore `task_name`, `prompt` is byte-identical as `message`, unsupported `description`/`subagent_type`/`model` are absent, and the reverse map can still recover slug/stage for cohort classification. Include a collision case or deterministic suffix test if sanitation can collapse distinct names.
- **Contractlint / prose assertions:** add fixture tests that fail on unversioned `send_input`, `wait_agent(handle)`, handle-scoped wait wording, or a new separate Codex v2 runtime file in Codex v2 references. Permit stale strings only in files or fixture blocks that explicitly mark legacy pre-v2 behavior. Assert the docs bind `send_message` and `followup_task` under `«addressable-worker»`, classify `wait_agent(timeout_ms)` under `«completion-signal»`, classify no-wait completion as queued/activity-driven unless proven otherwise, and name `list_agents` under `«roster-reconcile»`.
- **Live or durable fixture assertions:** update Codex shared-scenario parsing so v2 collaboration events recognize `spawn_agent`, `send_message`, `followup_task`, `wait_agent(timeout_ms)`, `list_agents`, and `interrupt_agent` where relevant. Reuse/rejection assertions should prove entity state, stage reports, and state-checkout commits, then tie those durable results to the correct v2 message tool instead of only counting old `send_input` events.
- **Shutdown probe:** before finalizing `«worker.shutdown»`, run a small live or fixture-backed probe documenting whether `interrupt_agent` terminates, pauses, or leaves a worker addressable. If live probing is not available in CI, record a fixture-backed adapter contract and require one manual live verification before release. Until this probe lands, implementation must not bless `interrupt_agent` as the shutdown binding.
- **Regression gates:** run `go test ./...`, `go test ./... -race`, and the Codex live lane (`go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`) for implementation that changes runtime behavior or live assertions.

## Stage Report: ideation

- DONE: Decide whether Codex multi_agent_v2 support should update the existing Codex runtime references or introduce a separate variant, with compatibility risks named.
  Chose existing `host: codex` references with a multi_agent_v2 adapter section; risks are documented under Compatibility risks.
- DONE: Map each observed v2 capability to concrete FO/ensign contract wording: spawn task-name sanitization, message/follow-up routing, interruptible global waits, roster reads, and shutdown.
  Host-binding language names the concrete v2 tools and the wording/semantics to apply in FO and ensign contracts.
- DONE: Produce an implementation-ready test plan covering dispatch-build mapping, contractlint/live assertions, and retirement or version-gating of stale send_input and wait_agent(handle) assumptions.
  Test plan covers mapping fixtures, contractlint, live/durable assertions, shutdown probing, and regression gates.

### Summary

Ideation narrowed Codex multi_agent_v2 support to an adapter update inside the existing Codex runtime contract, not a new runtime variant. The entity now records concrete tool mappings, compatibility risks, documentation deltas, and test coverage needed to retire or gate stale `send_input` and `wait_agent(handle)` assumptions.

## Stage Report: ideation (cycle 2)

- DONE: Align the Codex multi_agent_v2 design with the PR #409 capability-refactor review's `«fn»` host-binding direction.
  Replaced the parallel capability table with `«worker.spawn»`, `«addressable-worker»`, `«worker-identity»`, `«completion-signal»`, `«roster-reconcile»`, and `«worker.shutdown»` host-binding language.
- DONE: Preserve the conclusion that Codex multi_agent_v2 updates the existing Codex runtime adapter rather than creating a separate v2 runtime file.
  Proposed approach and documentation delta now explicitly say to add v2 bindings in the existing Codex adapter and reject a separate Codex v2 runtime file.
- DONE: Keep shutdown unresolved until probed and update implementation guidance accordingly.
  `«worker.shutdown»`, AC-4, and the shutdown probe now say not to bless `interrupt_agent` until live or fixture-backed evidence proves its behavior.

### Summary

The ideation body now follows the `«fn»` binding model from the cross-sprint review instead of adding a second capability registry. `send_message` and `followup_task` are framed as Codex v2 realizations of `«addressable-worker»`, `wait_agent(timeout_ms)` sits under `«completion-signal»`, task-name sanitation under `«worker-identity»`, and `list_agents` under `«roster-reconcile»`.
