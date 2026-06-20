---
title: Codex multi_agent_v2 runtime dispatch support
status: ideation
source: "FO live dispatch observation (2026-06-20): while dispatching status-validate-determinism under Codex multi_agent_v2, the current Codex runtime adapter no longer matches the live tool surface. v2 exposes spawn_agent(task_name,message,fork_turns), list_agents, wait_agent(timeout_ms), send_message, followup_task, and interrupt_agent. The shipped adapter still names send_input, wait_agent(handle), and hyphenated dispatch-build names as direct handles."
score:
started: 2026-06-20T03:58:46Z
completed:
verdict:
worktree:
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

The runtime contract needs a forward-compatible v2 variant or adapter section that captures these semantics without breaking existing v1/live-test assumptions before they are intentionally retired.

## Proposed approach

Update the existing Codex runtime references with an explicit **Codex multi_agent_v2** section rather than adding a separate runtime variant. Codex is still one host from the Spacedock workflow perspective; the incompatible piece is the live collab tool surface exposed inside newer Codex sessions. A separate `codex-v2` host would duplicate front-door, install, dispatch-build, and live-lane plumbing while making mixed-version support harder to reason about. The safer compatibility path is:

1. Keep `host: codex` in `spacedock dispatch build`.
2. Add a narrow adapter mapping from dispatch-build output to live v2 tool arguments.
3. Version-gate or retire stale v1 assumptions in contractlint and live assertions.
4. Preserve legacy transcript parsers only behind fixtures that explicitly name the old surface.

No spike is needed for the core v2 mapping: this task is seeded by a live Codex multi_agent_v2 dispatch observation with rejected hyphenated task names, successful underscore retry, `list_agents` roster reads, global wait interruption, and successful `send_message` delivery. The remaining risky behavior is shutdown semantics; implementation should prove it with a small live or fixture-backed probe before choosing whether `interrupt_agent`, `send_message`, or `followup_task` is the canonical cooperative-shutdown call.

## Capability mapping

| Observed v2 capability | FO / ensign contract wording to ship |
| --- | --- |
| `spawn_agent(task_name,message,fork_turns)` | Initial Codex dispatch parses `spacedock dispatch build` JSON, sanitizes emitted `name` to a v2 `task_name` by replacing non `[a-z0-9_]` runs with `_`, forwards emitted `prompt` unchanged as `message`, and omits unsupported `description`, `subagent_type`, and `model` arguments. Use `fork_turns` only when a Codex reuse policy explicitly needs context; initial file-pointer dispatch should normally use `fork_turns:"none"` or omit it if the default is accepted by the adapter. |
| Task-name validation | The helper may continue to emit hyphenated names for host-neutral worker identity and worktree naming. The Codex adapter owns v2 spawn-name sanitation and must retain a reverse map from sanitized task path to entity slug/stage when classifying cohorts. |
| `send_message(target,message)` | Use for non-interrupting context delivery to a running worker: status notes, cooperative shutdown notice when no new turn is required, or lightweight information that can be queued until the worker reaches a message boundary. Do not treat a `send_message` call itself as a completion signal. |
| `followup_task(target,message)` | Use for reuse, feedback rework, and any instruction that must trigger a new worker turn after a worker is idle or completed-but-reusable. This replaces v2 references to `send_input` for advancement/re-review. The FO must still wait for the worker's later final-status notification before advancing entity state. |
| `wait_agent(timeout_ms)` | Foreground wait is global, not handle-scoped. The FO may wait only after ready workflow/gate work is exhausted. A captain message or other operator activity can interrupt the wait and return control without failing or closing any worker; retrying is legal. Because the wait is global, the FO must inspect the returned mailbox update or roster state before attributing completion to an entity. |
| `list_agents(path_prefix?)` | Use as the Codex roster read. It can derive active/completed task paths for terminal teardown, stale-cohort classification, and debugging, but durable state remains authoritative for workflow advancement. |
| `interrupt_agent(target)` | Candidate shutdown/supersede primitive. Implementation should probe whether interruption ends a worker, merely stops the current turn, or leaves an addressable worker. Until proven, cooperative shutdown should be phrased as best-effort: send shutdown text first when preservation matters, then interrupt only for explicit supersede/terminal cleanup if live evidence supports it. |
| Final-status mailbox notification | Completion remains the ensign's concise final message observed through the FO mailbox. No-wait completion is classified as queued/activity-driven unless a clean idle probe proves autonomous wake-up. |

## Compatibility risks

- Existing Codex live fixture parsers bind reuse to `send_input` events and thread IDs. A blind text replacement would make old fixtures fail for the wrong reason, so tests need either a versioned v1 parser path or fixtures marked `codex_multi_agent_v2`.
- `wait_agent(handle)` assumptions are embedded in prose and watchdog naming. v2's global wait can return an unrelated worker or an interruption, so entity attribution must come from mailbox content, task path, or durable state.
- Sanitized `task_name` values can collide if two helper names differ only by punctuation. The adapter should either include a deterministic suffix when sanitation changes the name or test that the known dispatch-build name format remains collision-free for `(worker_key, slug, stage, cycle)`.
- `description` and `model` are mandatory host-neutral helper outputs today, but unsupported by v2 `spawn_agent`. The Codex adapter must clearly say "helper emits them, v2 spawn omits them" so the host-neutral core does not weaken Claude/Pi forwarding guarantees.
- `send_message` may not wake an idle/completed worker for rework. Reuse and feedback must use `followup_task` unless a probe proves `send_message` starts the needed turn.
- `interrupt_agent` may be stronger than cooperative shutdown and could lose uncommitted worker state. Shutdown wording must keep the worker's preservation obligation and use interruption only at terminal/supersede boundaries with evidence.

## Acceptance criteria

**AC-1 - Initial Codex v2 dispatch maps helper output to the live spawn surface.**  
Verified by: a test or live fixture showing `dispatch build` output with a hyphenated `name` is converted to a valid v2 `task_name` while forwarding the helper `prompt` unchanged as `message`, omitting unsupported `description` and `model` fields.

**AC-2 - Codex v2 continuation and steering use the live message tools.**  
Verified by: adapter text and tests distinguish `send_message` context delivery from `followup_task` turn-triggering reuse/feedback, and remove or explicitly version-gate `send_input` from v2 instructions.

**AC-3 - Codex v2 waiting semantics match observed behavior.**  
Verified by: adapter text and live/fixture evidence state that `wait_agent(timeout_ms)` is global, captain input interrupts wait without failing the worker, and no-wait completion remains queued/activity-driven unless an autonomous wake probe proves otherwise.

**AC-4 - Terminal/supersede shutdown has a concrete v2 call.**  
Verified by: adapter text maps cooperative shutdown to the live v2 surface, likely `interrupt_agent` or a `send_message`/`followup_task` shutdown protocol depending on what probes prove, and names residual risk.

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

Implementation should update `skills/first-officer/references/codex-first-officer-runtime.md`, `skills/ensign/references/codex-ensign-runtime.md`, `skills/feedback-rejection-flow/SKILL.md`, and `docs/dev/codex-idle-notification-probe.md`.

Replace unversioned Codex reuse wording like:

```text
Route feedback or continuation to an existing Codex worker with `send_input`.
```

with v2-aware wording like:

```text
On Codex multi_agent_v2, route non-triggering notes with `send_message(target,message)` and route reuse, feedback, or rework that must start a worker turn with `followup_task(target,message)`. Legacy `send_input` belongs only to explicitly versioned pre-v2 fixtures or references.
```

Replace foreground-wait wording like:

```text
The FO calls `wait_agent(handle)` as the next useful idle action.
```

with:

```text
On Codex multi_agent_v2, the FO calls global `wait_agent(timeout_ms)` only after workflow and gate work is exhausted. A wait return must be attributed by mailbox content, task path, roster state, or durable workflow state, not by a handle argument.
```

Add a Codex v2 dispatch paragraph:

```text
For Codex multi_agent_v2, map `dispatch build` output as follows: sanitize emitted `name` to a lowercase digit/underscore `task_name`, pass emitted `prompt` unchanged as `message`, omit unsupported `description`, `subagent_type`, and `model`, and record the sanitized task path as the live worker handle.
```

## Test plan

- **Dispatch-build mapping tests:** add a Go fixture around the Codex adapter mapping from helper JSON to v2 spawn call input. Assert that a helper name such as `spacedock-ensign-status-validate-determinism-implementation` becomes a valid underscore `task_name`, `prompt` is byte-identical as `message`, unsupported `description`/`subagent_type`/`model` are absent, and the reverse map can still recover slug/stage for cohort classification. Include a collision case or deterministic suffix test if sanitation can collapse distinct names.
- **Contractlint / prose assertions:** add fixture tests that fail on unversioned `send_input`, `wait_agent(handle)`, or handle-scoped wait wording in Codex v2 references. Permit those strings only in files or fixture blocks that explicitly mark legacy pre-v2 behavior. Assert the docs distinguish `send_message` from `followup_task`, classify no-wait completion as queued/activity-driven unless proven otherwise, and name `list_agents` as the roster read.
- **Live or durable fixture assertions:** update Codex shared-scenario parsing so v2 collaboration events recognize `spawn_agent`, `send_message`, `followup_task`, `wait_agent(timeout_ms)`, `list_agents`, and `interrupt_agent` where relevant. Reuse/rejection assertions should prove entity state, stage reports, and state-checkout commits, then tie those durable results to the correct v2 message tool instead of only counting old `send_input` events.
- **Shutdown probe:** before finalizing terminal/supersede wording, run a small live or fixture-backed probe documenting whether `interrupt_agent` terminates, pauses, or leaves a worker addressable. If live probing is not available in CI, record a fixture-backed adapter contract and require one manual live verification before release.
- **Regression gates:** run `go test ./...`, `go test ./... -race`, and the Codex live lane (`go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`) for implementation that changes runtime behavior or live assertions.

## Stage Report: ideation

- DONE: Decide whether Codex multi_agent_v2 support should update the existing Codex runtime references or introduce a separate variant, with compatibility risks named.
  Chose existing `host: codex` references with a multi_agent_v2 adapter section; risks are documented under Compatibility risks.
- DONE: Map each observed v2 capability to concrete FO/ensign contract wording: spawn task-name sanitization, message/follow-up routing, interruptible global waits, roster reads, and shutdown.
  Capability mapping table names the concrete v2 tools and the wording/semantics to apply in FO and ensign contracts.
- DONE: Produce an implementation-ready test plan covering dispatch-build mapping, contractlint/live assertions, and retirement or version-gating of stale send_input and wait_agent(handle) assumptions.
  Test plan covers mapping fixtures, contractlint, live/durable assertions, shutdown probing, and regression gates.

### Summary

Ideation narrowed Codex multi_agent_v2 support to an adapter update inside the existing Codex runtime contract, not a new runtime variant. The entity now records concrete tool mappings, compatibility risks, documentation deltas, and test coverage needed to retire or gate stale `send_input` and `wait_agent(handle)` assumptions.
