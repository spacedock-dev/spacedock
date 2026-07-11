---
title: Bind Codex context budget to per-thread token telemetry
status: ideation
source: "Captain request 2026-07-11 plus live Codex 0.144.1 schema and rollout evidence gathered while reusing 88t workers."
started: 2026-07-11T05:50:29Z
completed:
verdict:
score: 0.9
worktree:
issue:
id: bct1zbqhbkatrwmqf6s6qd9v
---

## Problem

The Codex first-officer adapter declares `«context-budget»` ABSENT, so reuse-condition 0 always passes. That is no longer an honest capability boundary: Codex emits per-thread token-window data, but the Spacedock runtime has not bound a worker path to that data.

The risk is already visible in a live rollout. The implementation child `019f4f17-ccbe-7343-a432-bc2d445706ec` reached 307,584 / 353,400 active tokens (45,816 remaining, 13.0%) while it was repeatedly reused. Its validation sibling `019f4f2c-1776-7713-97b4-7a1b2e750688` was at 192,895 / 353,400. Both identify parent `019f499a-3e03-7893-ae75-b7ca10ac1ff6`. The current `list_agents` surface exposes task paths and status, not thread IDs or token usage.

This work is limited to the Codex budget binding: identity, telemetry, fail-safe reuse, its command surface, and tests. It does not add PR, mod, or generic roster behavior.

## Spike findings

The 2026-07-11 spike ran Codex CLI 0.144.1's generated app-server v2 schema and an isolated local app-server session.

- `thread/tokenUsage/updated` carries `threadId`, `turnId`, `tokenUsage.last`, `tokenUsage.total`, and `modelContextWindow`.
- A parent spawned `/root/alpha` and `/root/beta`. Parent `subAgentActivity` mapped alpha to `019f4fc2-547f-7720-8ed9-a6183fe51e21`; `thread/read(alpha)` reported the same parent ID and `source.subAgent.thread_spawn.agent_path: /root/alpha`. Token updates arrived on that exact child ID.
- Manual compaction on alpha kept `total.totalTokens` at 18,893 while `last.totalTokens` fell from 18,893 to 4,097, with a 353,400-token window. `last.totalTokens` is therefore the active post-compaction window; `total.totalTokens` is lifetime accounting and must never drive reuse.
- The completed alpha accepted a parent-routed follow-up, started a new child turn, and emitted a new update under the same child thread ID. A completed-but-addressable worker remains a valid telemetry subject.
- Rollout JSONL independently persists the required metadata: `session_meta` has `id`, `parent_thread_id`, and `source.subagent.thread_spawn.agent_path`; `token_count` has `last_token_usage`, `total_token_usage`, and `model_context_window`. The current worker's latest record is 188,603 / 353,400 active tokens while lifetime total is 5,372,887.

### Bridge gate

App-server is the currently exposed experimental integration surface, pinned to the probed 0.144.1 schema, and it can host the normal terminal UI through `codex --remote` over a local Unix socket. The implementation must first prove that a Spacedock observer can subscribe beside that remote UI and receive its child events. A direct two-proxy attempt ended with a broken pipe, so multiplexing is not yet evidence. An unknown/incompatible schema or failed two-client replay is unavailable; it must not create a binding or promote the capability to PRESENT. Do not substitute a rollout-tail-only binding.

## Proposed direction

Add a Codex telemetry bridge owned by `spacedock codex`. It starts app-server on a per-session Unix socket in a mode-0700 directory, launches the ordinary terminal UI with `codex --remote unix://…`, and keeps a metadata-only observer connected to the same server. The launcher creates an opaque per-session capability and passes it to the observer and all Spacedock probe calls; it is never inferred from a task path or reported in workflow state. The observer writes atomically to `$CODEX_HOME/spacedock/telemetry/v1/<capability>.json` with mode 0600; it never stores prompts, responses, reasoning, or tool output.

The runtime records one of three launch states. `legacy-absent` preserves the current behavior for a direct, non-Spacedock Codex session. `present` requires a successful listener, schema, observer, and UI handshake. `required-unavailable` begins as soon as `spacedock codex` elects the bridge and remains so after an observer/listener/schema failure, reconnect, or old bridged-session resume without fresh telemetry. Only `legacy-absent` satisfies the core's ABSENT condition; `required-unavailable` always fresh-dispatches.

The observer records a worker key `{session_capability, task_path, dispatch_epoch}` with its parent thread and child thread. It assigns `dispatch_epoch` from a confirmed parent-side spawn event, not from task text. It accepts a binding only when all available evidence agrees:

1. The parent item identifies the exact task path and exactly one child (`subAgentActivity.agentPath`/`agentThreadId`; `collabAgentToolCall.receiverThreadIds` is corroborating evidence when emitted).
2. The child thread reports that same `parentThreadId` and `source.subAgent.thread_spawn.agent_path`.
3. The mapping belongs to the current dispatch epoch. A fresh replacement creates a new epoch; late telemetry from a superseded thread is ignored.

For each bound child, retain the newest event in observer receive order: `{observer_generation, receive_seq, thread_id, turn_id, active_tokens, context_window, observed_at}`. `active_tokens` is `tokenUsage.last.totalTokens`; `context_window` is `tokenUsage.modelContextWindow`. Token values are not required to rise: a later compaction is allowed to lower active tokens. On reconnect, discard every cached snapshot and enter `required-unavailable` until the new observer receives a fresh authoritative binding and token event. The bridge rejects a nil/zero window, a missing `last.totalTokens`, duplicate candidates, a parent/path mismatch, or an event older than two minutes. It compares tokens with integers: `active_tokens * 100 <= context_window * 60`. The equality boundary remains reusable; a value above 60% is not.

Expose a new, Codex-specific command rather than changing Claude's byte-stable command:

```text
spacedock dispatch codex-context-budget --worker <task-path>
```

On success it emits stable JSON in this order: `worker`, `thread_id`, `turn_id`, `active_tokens`, `context_window`, `remaining_tokens`, `remaining_pct`, `usage_pct`, `threshold_pct`, `source`, `observed_at`, `reuse_ok`. For example, the compacted alpha snapshot is 4,097 active tokens, 349,303 remaining, and `reuse_ok: true`. The 307,584 / 353,400 snapshot is 87.0% used and emits `reuse_ok: false`.

On absent, stale, malformed, ambiguous, unsafe, or capability-mismatched telemetry, the command exits non-zero, writes no JSON result, and names the reason on stderr. It reads the inherited opaque session capability and rejects a missing or mismatched capability before looking up a worker. The Codex adapter treats every `required-unavailable` result as PRESENT-but-unavailable and fresh-dispatches. It must never relabel a failed bridge as ABSENT merely to satisfy reuse-condition 0.

Rollout JSONL is recovery and diagnosis only. Its reader may inspect only `session_meta` identity fields and `token_count` token fields under the canonical `$CODEX_HOME/sessions` root. It resolves the root, rejects symlinks/non-regular files and paths outside that root, validates the exact capability-bound parent and task path, uses the JSON event timestamp for the same two-minute freshness bound, and refuses every other record type. The observer applies the same allowlist before serialization; it never writes raw event JSON. `/agent` and `/status` remain manual diagnostics, not protocol inputs.

The Codex first-officer runtime changes from `ABSENT` to a positive `PRESENT` binding only after the bridge handshake succeeds. A bridged-session resume or bridge restart is unavailable until it receives fresh telemetry and therefore fresh-dispatches. A direct, non-Spacedock Codex session stays `legacy-absent` and retains the current host behavior.

## Acceptance criteria

- **AC-1:** Each spawned Codex worker has one deterministic `{session_capability, task_path, dispatch_epoch} -> child_thread_id` binding, and the binding's latest record has the same thread ID. A fixture table covers zero, one, and multiple candidates, wrong parent/path, two parents with the same task path, and replacement; a live two-child replay confirms distinct alpha/beta bindings through parent activity plus `thread/read`.
- **AC-2:** `codex-context-budget` returns active tokens, context window, remaining tokens, remaining percentage, usage percentage, source, timestamp, and `reuse_ok` for a running worker and a completed-but-addressable worker. Golden CLI tests cover field order and arithmetic; the live replay sends a follow-up to a completed child and observes a new turn under its original thread ID.
- **AC-3:** A bridge-required Codex capability produces no reusable result for missing, stale, malformed, ambiguous, mismatched, reconnect-stale, observer-startup-failed, or over-threshold telemetry. Table and CLI tests assert a non-zero, no-JSON result for unavailable cases; 60.0% is reusable and the next token is not. The adapter's decision test asserts fresh dispatch for every `required-unavailable` result while direct legacy sessions retain ABSENT behavior.
- **AC-4:** Compaction decisions use `last.totalTokens`, never lifetime `total.totalTokens`. A fixture reproduces 18,893 -> 4,097 while total remains 18,893 and the recorded 307,584 -> 21,504-style reset while lifetime total grows; both must flip only the active-window calculation.
- **AC-5:** Telemetry remains isolated across siblings, a completed worker's follow-up turn, a reconnect, and a replacement worker. Event-order tests use observer receive sequence, permit the lower post-compaction value, clear snapshots after reconnect, retain the current epoch only, and reject late superseded-thread updates; a live replay routes a follow-up to alpha while beta's snapshot remains unchanged.
- **AC-6:** Neither the observer nor the recovery reader persists prompt/response fields, and both accept only exact, fresh metadata/token-count records beneath a canonical sessions root. Temp-tree tests cover a valid record, stale timestamp, wrong parent/path, malformed JSON, symlink escape, and non-regular file; sentinel-bearing raw events prove the telemetry file, stderr, and reports omit the sentinel. Tests assert a mode-0700 directory, mode-0600 file, symlink-safe atomic write, and the two-minute diagnostic.
- **AC-7:** The feedback-cycle decision fresh-dispatches a 307,584 / 353,400-style worker and reuses a sub-threshold sibling. A deterministic decision test injects both snapshots; the live-gated bridge replay uses a low-cost threshold override to make one real child over budget and its compacted sibling under budget, then observes one new child thread and one follow-up on the retained thread.

## Test plan

1. Start with the two-client bridge spike: app-server listener, ordinary `codex --remote` terminal UI, observer client, two named children, and a compacted child. Save only normalized IDs, counters, and event timestamps. If the observer cannot see the UI-owned threads, stop: keep the feature `required-unavailable` and report the host-interface blocker rather than shipping a log-tail substitute.
2. Add focused Go table tests for capability-scoped identity reduction, observer-generation/receive-sequence selection, threshold arithmetic, compaction, stale/missing/ambiguous states, startup/restart failure, replacement isolation, and the metadata-only fallback reader. Include prompt/response sentinels in raw fixture events and assert that the state file, diagnostics, and reports never contain them.
3. Add golden command tests for `codex-context-budget` success, over-budget, and every unavailable failure. Keep the existing Claude `context-budget` goldens byte-identical.
4. Add a live-gated app-server replay. It must prove one retained thread receives a parent-routed follow-up while a separate over-budget worker is replaced, then verify child IDs, telemetry state, entity report, state-checkout git log, and clean status. A test-only threshold makes this replay cheap; production remains 60%.

Estimated cost: small parser/registry and command packages plus a medium Codex-launch integration. The bridge spike is the first and blocking proof; unit and golden tests are cheap, while the live replay remains opt-in.

## Documentation change

The implementation updates `docs/runtime-support.md` and the Codex first-officer runtime binding. The public-facing addition is:

```diff
+ Codex sessions launched by `spacedock codex` keep a local, metadata-only
+ context-budget record. `spacedock dispatch codex-context-budget --worker <task-path>`
+ reports a worker's active context window. A bridge startup, restart, or telemetry
+ failure causes a fresh dispatch; Spacedock never reads transcript content for this decision.
```

The runtime binding replaces `«context-budget»: ABSENT` with the bridge command and its failure rule. It does not mention raw app-server field names outside the Codex-specific binding.

## Stage Report: ideation

- DONE: AC-1 worker/task-path identity. Live spike: parent subAgentActivity mapped /root/alpha to child thread 019f4fc2-547f-7720-8ed9-a6183fe51e21; child thread/read confirmed its parent and task path.
- DONE: AC-4 active-window accounting. Compaction changed last.totalTokens 18,893 -> 4,097; lifetime total.totalTokens stayed 18,893 (modelContextWindow 353,400).
- SKIPPED: AC-2, AC-3, AC-5, AC-6, and AC-7. Intentionally deferred to implementation/validation: AC-2 golden CLI and routed follow-up (3-4); AC-3 threshold/adapter (2-3); AC-5 ordering/alpha-beta (2,4); AC-6 sentinel/mode/symlink (2); AC-7 snapshot/live threshold (4).

### Summary

Only AC-1 and AC-4 have live-spike evidence. The other criteria remain proof-planned; the two-client observer is the first implementation gate.
