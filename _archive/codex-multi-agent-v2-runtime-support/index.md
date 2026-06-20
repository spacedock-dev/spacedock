---
title: Codex multi_agent_v2 runtime dispatch support
status: validation
source: "FO live dispatch observation (2026-06-20): while dispatching status-validate-determinism under Codex multi_agent_v2, the current Codex runtime adapter no longer matches the live tool surface. v2 exposes spawn_agent(task_name,message,fork_turns), list_agents, wait_agent(timeout_ms), send_message, followup_task, and interrupt_agent. The shipped adapter still names send_input, wait_agent(handle), and hyphenated dispatch-build names as direct handles."
score:
started: 2026-06-20T03:58:46Z
completed: 2026-06-20T09:13:35Z
verdict: PASSED
worktree:
issue:
sprint:
sprint-readiness:
id: f2r8cnyxj9pf24xrsf71szb0
mod-block:
pr: "#414"
archived: 2026-06-20T09:13:39Z
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

## Stage Report: implementation

- DONE: Update the existing Codex runtime references with standalone-safe multi_agent_v2 host-binding language, using `«fn»` anchors without requiring b2/shared-core changes to have landed first.
  Added Codex multi_agent_v2 bindings to the existing Codex FO and ensign runtime references, the feedback rejection flow, and the idle-notification probe. The implementation keeps `host: codex`, uses `«worker.spawn»`, `«worker-identity»`, `«addressable-worker»`, `«completion-signal»`, `«roster-reconcile»`, and leaves `«worker.shutdown»` unresolved until probed instead of adding a separate Codex v2 runtime file.
- DONE: Replace or version-gate stale Codex assumptions around `send_input`, `wait_agent(handle)`, task-name mapping, completion delivery, roster reads, and unresolved shutdown behavior.
  Replaced unversioned runtime prose with `spawn_agent(task_name,message,fork_turns)`, `send_message(target,message)`, `followup_task(target,message)`, global `wait_agent(timeout_ms)`, `list_agents(path_prefix?)`, queued/activity-driven no-wait delivery, and an explicit "do not bless `interrupt_agent`" shutdown statement. Legacy `send_input` remains only in legacy fixture/parser contexts.
- DONE: Add focused tests or contractlint/live-fixture updates that prove the Codex v2 wording and dispatch mapping expectations without weakening existing host-neutral helper guarantees.
  Added `internal/dispatch` tests and adapter code proving helper JSON maps to v2 spawn args while preserving prompt bytes and omitting unsupported helper fields; added contractlint coverage for v2 binding prose and absence of unversioned stale wording; updated Codex shared-scenario reuse parsing and fixtures to accept `followup_task` while preserving legacy `send_input` fixture coverage.

### Verification

- PASS: `go test ./...` -> 1637 passed in 17 packages.
- PASS: `go test ./... -race` -> 1637 passed in 17 packages.
- RAN: `gofmt -w ./cmd ./internal`, then scoped `gofmt` on touched Go files after reverting unrelated comment corruption from the whole-tree formatting pass.
- FAIL: `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` -> `TestLiveCodexSharedScenarios/rejection-flow` failed because the live FO emitted 2 validation `spawn_agent` calls, meaning it fresh-dispatched the cycle-2 validator instead of reusing the kept-alive cycle-1 reviewer. The same live run reported 5 passed, 2 failed in `internal/ensigncycle`.

### Summary

Implementation added standalone Codex multi_agent_v2 host bindings to the existing Codex runtime references, introduced a small dispatch-build-to-v2-spawn adapter with collision detection, and updated contract/parser tests for v2 `followup_task` reuse. Offline and race suites pass; the Codex live lane still exposes the live rejection-flow fresh-dispatch behavior and is recorded as a remaining runtime failure.

## Stage Report: implementation (feedback cycle 1)

- DONE: Root-cause why the v2 runtime bindings/tests still allow or induce fresh cycle-2 validation dispatch instead of `followup_task` to the existing validation worker.
  The live transcript showed the FO read the new Codex v2 reviewer-reuse requirement, then self-classified the host as lacking `followup_task` because the observed stream had only used spawn/wait so far. That caused a fallback to fresh cycle-2 validation. The root cause was instructional ambiguity: the runtime contract did not forbid inferring missing `followup_task` from prior transcript shape.
- DONE: Fix the implementation so multi_agent_v2 rejection-flow reuses the kept-alive validation reviewer via the correct addressable-worker binding, while preserving durable state assertions and legacy fixture compatibility where explicitly versioned.
  Tightened `codex-first-officer-runtime.md` and `feedback-rejection-flow/SKILL.md`: Codex multi_agent_v2 is not one-shot; feedback rejection must keep the first validation reviewer addressable, re-run it with `followup_task(target,message)`, and only fresh-spawn when the existing reviewer is not addressable or reuse conditions fail. The FO must attempt the v2 binding and report a concrete tool-surface blocker instead of silently fresh-spawning.
- DONE: Add/adjust focused tests so this failure is caught without relying only on the live lane.
  Extended `TestCodexMultiAgentV2RuntimeReferencesUseLiveHostBindings` to require the kept-reviewer `followup_task` rule, the fresh-spawn exception boundary, and the "do not infer absence" wording. Existing shared-scenario parser tests still preserve legacy `send_input` fixture compatibility while accepting v2 `followup_task`.

### Verification

- PASS: `go test ./internal/contractlint -run TestCodexMultiAgentV2RuntimeReferencesUseLiveHostBindings` -> 1 passed.
- PASS: `go test ./internal/ensigncycle -run TestAssertCodexReviewerReuse` -> 9 passed.
- RAN: `gofmt -w ./cmd ./internal`; it again introduced unrelated comment punctuation churn in `internal/cli/prose_function_routing_test.go` and `internal/status/section_read.go`, which was reverted before commit.
- PASS: `go test ./...` -> 1637 passed in 17 packages.
- PASS: `go test ./... -race` -> 1637 passed in 17 packages.
- PASS: `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-codex-live-f2r8cnyxj9-followup-attempt go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` -> 7 passed in 1 package.

### Summary

Feedback cycle 1 fixed the remaining live rejection-flow failure by making Codex multi_agent_v2 reviewer reuse explicit and non-optional for feedback re-review. The final local Codex live lane passed all seven shared scenarios.

## Stage Report: implementation (feedback cycle 2)

- DONE: Remove agent-facing `Codex multi_agent_v2` / v2 behavior guidance and align Codex runtime prose with probe-based `«fn»` host bindings.
  Replaced runtime and skill wording with live tool-surface probes: bind `«worker.spawn»` from `spawn_agent(task_name,message,fork_turns)`, `«addressable-worker»` from `send_message(target,message)` / `followup_task(target,message)`, foreground `«completion-signal»` from `wait_agent(timeout_ms)`, and `«roster-reconcile»` from `list_agents(path_prefix?)`. Kept `«worker.shutdown»` unresolved until a shutdown-specific probe proves a binding.
- DONE: Preserve concrete blocker behavior instead of version inference.
  Codex runtime instructions now say missing or shape-rejecting tools are concrete tool-surface blockers, and fallback is allowed only to behavior the live surface actually supports. They also retain the feedback-rejection rule to attempt the probed `followup_task` binding instead of inferring absence from older transcript labels.
- DONE: Update contractlint so version-taxonomy leakage fails.
  Renamed the Codex contractlint test around current runtime live tool-surface probes and banned `Codex multi_agent_v2`, `multi_agent_v2`, `Codex v2`, and `pre-v2` from the agent-facing Codex runtime/skill references.

### Verification

- PASS: `go test ./internal/contractlint -run TestCodexCurrentMultiAgentRuntimeReferencesUseLiveToolSurfaceProbe` -> 1 passed.
- PASS: `go test ./internal/ensigncycle -run TestAssertCodexReviewerReuse` -> 9 passed.
- RAN: `gofmt -w ./cmd ./internal`; unrelated formatting/comment churn in `internal/cli/prose_function_routing_test.go`, `internal/status/section_read.go`, and `internal/ensigncycle/haiku_loop_spike_live_test.go` was reverted before commit.
- PASS: `go test ./...` -> 1679 passed in 17 packages.
- PASS: `go test ./... -race` -> 1679 passed in 17 packages.

### Summary

Feedback cycle 2 removes version-based Codex behavior guidance from agent-facing runtime files. Codex instructions now bind capabilities from the live tool surface and report concrete tool-shape blockers instead of reasoning from a runtime version label.

## Stage Report: validation

- DONE: Verified AC-1, initial Codex v2 dispatch mapping.
  Evidence: `internal/dispatch/codex_v2_adapter.go` maps helper JSON `name` to a lowercase/underscore `task_name`, preserves `prompt` as `message`, omits unsupported helper fields from `ToolArgs`, and returns slug/stage identity. `TestCodexMultiAgentV2SpawnInputMapsBuildOutput` and `TestCodexMultiAgentV2SpawnInputRejectsSanitizedCollision` passed via `go test ./internal/dispatch` (243 passed).
- DONE: Verified AC-2, continuation and steering wording plus legacy fixture handling.
  Evidence: `skills/first-officer/references/codex-first-officer-runtime.md` and `skills/feedback-rejection-flow/SKILL.md` distinguish `send_message(target,message)` from `followup_task(target,message)` and require kept-reviewer re-review instead of fresh-spawning. `TestCodexMultiAgentV2RuntimeReferencesUseLiveHostBindings` passed, and `TestAssertCodexReviewerReuse` passed with the parser accepting v2 `followup_task` plus explicitly legacy `send_input`.
- DONE: Verified AC-3, waiting semantics.
  Evidence: Codex FO runtime and `docs/dev/codex-idle-notification-probe.md` describe global `wait_agent(timeout_ms)`, interruptible foreground waits, queued/activity-driven no-wait delivery, and attribution through mailbox/task-path/roster/durable state. Contractlint coverage passed in `TestCodexMultiAgentV2RuntimeReferencesUseLiveHostBindings`.
- DONE: Verified AC-4, terminal/supersede shutdown remains unresolved rather than over-blessed.
  Evidence: `skills/first-officer/references/codex-first-officer-runtime.md` keeps `«worker.shutdown»` unresolved and says not to bless `interrupt_agent`; `skills/ensign/references/codex-ensign-runtime.md` keeps cooperative shutdown separate from hard interruption. Contractlint coverage passed. No concrete shutdown binding is claimed.
- DONE: Verified AC-5, old Codex assumptions are either preserved as legacy or retired from v2 references.
  Evidence: `rg --files | rg '(^|/)codex[-_]?v2|codex.*multi_agent_v2|multi_agent_v2'` found only the dispatch adapter/test and contractlint test, not a separate Codex v2 runtime reference file. Legacy `send_input` remains in parser/fixture contexts while the v2 runtime references reject unversioned `send_input` and `wait_agent(handle)`.
- DONE: Reproduced the feedback-cycle fix against the live shared-scenario lane.
  Evidence: `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-codex-live-validation-f2r8cnyxj9-20260619235358 go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` passed all six subtests in 778.905s, including `rejection-flow` in 424.86s. The rejection-flow artifact emitted exactly one validation `spawn_agent` and routed cycle-2 validation to the original validation thread; final durable state contained `- Cycle 1: REJECTED`, `- Cycle 2: PASSED`, and `shared-rejection-fix: applied`.
- DONE: Ran required regression gates.
  Evidence: `gofmt -w ./cmd ./internal` completed; it again produced unrelated comment punctuation churn in `internal/cli/prose_function_routing_test.go` and `internal/status/section_read.go`, which was restored before this report. `go test ./...` passed with 1637 tests in 17 packages. `go test ./... -race` passed with 1637 tests in 17 packages.

### Findings

- Residual risk: the live rejection-flow artifact used Codex stream events labeled `send_input` for reuse, even though the v2 runtime references now require `followup_task(target,message)` for turn-triggering rework. The live run still proves the behavioral fix that failed earlier: the FO did not fresh-spawn the cycle-2 validator and did re-review on the kept validation worker. Direct live proof of a `followup_task` event remains outside this artifact.
- Scope note: the branch includes `docs/dev/_debriefs/2026-06-19-02.md`, which is not part of the AC surface for this runtime-support task. It did not affect tests or validation, but it should be reviewed before merge as unrelated branch content.

### Recommendation

PASSED. The acceptance criteria are satisfied by code/fixture tests, updated runtime references, preserved legacy parser handling, and a passing live shared-scenario lane. The only material caveat is that the live transcript proves kept-reviewer reuse under the current stream label rather than a direct `followup_task` event.

### Summary

Validation passed AC-1 through AC-5, reproduced the feedback-cycle fix in the live Codex shared-scenario lane, and confirmed the required Go, race, formatting, and focused checks. The recommended gate result is PASSED with the live `send_input` stream-label caveat recorded.

## Stage Report: implementation (feedback cycle 3)

- DONE: Remove runtime-host and runtime-tool leakage from the shared feedback rejection procedure.
  `skills/feedback-rejection-flow/SKILL.md` now routes only through `«addressable-worker»` and waits only through `«completion-signal»`. It no longer names Codex, Claude, `followup_task`, `send_message`, `SendMessage`, `list_agents`, or other host-specific tool names.
- DONE: Keep concrete Codex call shapes in the Codex runtime adapter binding section.
  `skills/first-officer/references/codex-first-officer-runtime.md` keeps the explicit bindings in `## Live Tool Surface Probe`; other Codex runtime and probe prose now refers to capabilities such as `«worker.spawn»`, `«addressable-worker»`, `«completion-signal»`, and `«roster-reconcile»`.
- DONE: Extend contractlint for the host-neutral boundary.
  Added checks that fail if feedback-rejection-flow contains runtime host/tool names, and checks that concrete Codex tool call shapes stay in the Codex live binding section while still proving the live Codex binding exists.

### Verification

- PASS: `go test ./internal/contractlint -run 'TestCodex|TestFeedback' -count=1` -> 4 passed in 1 package.
- PASS: `go test ./internal/ensigncycle -run TestAssertCodexReviewerReuse -count=1` -> 9 passed in 1 package.
- RAN: `gofmt -w ./cmd ./internal`; as before, it produced unrelated comment/alignment churn in `internal/cli/prose_function_routing_test.go`, `internal/status/section_read.go`, and `internal/ensigncycle/haiku_loop_spike_live_test.go`, which was reverted before commit.
- PASS: `go test ./...` -> 1681 passed in 17 packages.
- PASS: `go test ./... -race` -> 1681 passed in 17 packages.

### Summary

Feedback cycle 3 makes the rejection-flow procedure host-neutral and moves concrete runtime call shapes behind the Codex adapter binding boundary. Focused contract and reviewer-reuse checks pass, along with the full Go and race suites.

## Stage Report: implementation (feedback cycle 4)

- DONE: Remove Codex adapter negative host-contrast prose.
  Removed the Codex FO `## Team Creation` section, team registry / `TeamDelete` contrast, `team-lead` mailbox warning, and no-op/absent-team framing. The Codex ensign runtime now describes dispatch files, context-budget availability, and completion behavior positively.
- DONE: Remove shared merge cleanup procedure from the Codex adapter.
  Replaced `## Terminal Teardown (Merge-and-Cleanup step 10)` and the numbered teardown paragraph with `## Cleanup Capability Bindings`, which only states Codex binding status for `«worker.shutdown»`, `«roster-reconcile»`, `«addressable-worker»`, and durable workflow state authority.
- DONE: Extend contractlint for the cleaner adapter boundary.
  Added Codex FO/ensign negative-contrast guards banning step-number coupling, `Team Creation`, `TeamDelete`/`TeamCreate`, `team registry`, `team-lead`, `Backstop`, `Codex declares none`, `no-op`, and related ensign contrast wording while requiring positive capability binding text.

### Verification

- PASS: `go test ./internal/contractlint -run 'TestCodex|TestFeedback' -count=1` -> 6 passed in 1 package.
- PASS: `go test ./internal/ensigncycle -run TestAssertCodexReviewerReuse -count=1` -> 9 passed in 1 package.
- PASS: sweep `rg` for removed contrast terms in Codex FO/ensign runtime references -> no matches.
- RAN: `gofmt -w ./cmd ./internal`; it again produced unrelated comment/alignment churn in `internal/cli/prose_function_routing_test.go`, `internal/status/section_read.go`, and `internal/ensigncycle/haiku_loop_spike_live_test.go`, which was reverted before commit.
- PASS: `go test ./...` -> 1683 passed in 17 packages.
- PASS: `go test ./... -race` -> 1683 passed in 17 packages.

### Summary

Feedback cycle 4 leaves Codex runtime references as positive adapter bindings instead of negative comparisons or shared procedure restatements. Terminal cleanup semantics now live in shared-core capabilities; the Codex adapter only records binding status and available Codex evidence.
