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

The Codex first-officer adapter currently declares `«context-budget»` ABSENT, so reuse-condition 0 is automatically satisfied. That is no longer an honest capability ceiling.

Codex app-server v2 emits `thread/tokenUsage/updated` with `threadId`, `turnId`, `last`, `total`, and `modelContextWindow`. Subagent threads expose `parentThreadId`; collaboration items expose `receiverThreadIds`. Local rollout JSONL also persists `session_meta` and `token_count` records.

Current live evidence demonstrates the risk:

- 88t implementation worker thread `019f4f17-ccbe-7343-a432-bc2d445706ec`: last usage 307,584 / 353,400 tokens, leaving 45,816 (about 13%) while it was repeatedly reused.
- 88t validation worker thread `019f4f2c-1776-7713-97b4-7a1b2e750688`: 192,895 / 353,400 tokens.
- Both child threads share parent `019f499a-3e03-7893-ae75-b7ca10ac1ff6`.
- The current `list_agents` surface exposes task paths/status but not thread IDs or token telemetry, so the missing piece is identity binding and access—not measurement.

## Proposed direction

Bind Codex `«context-budget»` prospectively from app-server notifications:

1. Record the child thread ID returned by `collabAgentToolCall.receiverThreadIds` / `ThreadStartedNotification.parentThreadId` alongside the existing Spacedock worker identity.
2. Maintain the latest `ThreadTokenUsageUpdatedNotification` per child thread.
3. Compute remaining active context only after a spike proves which field tracks the post-compaction active window; do not assume cumulative `total` means current context.
4. Return PRESENT + available when current telemetry maps to the worker; PRESENT + unavailable must force fresh dispatch under the existing fail-safe reuse rule.
5. Use rollout JSONL as a diagnostic/recovery fallback only if it can map the exact task/thread without reading prompt content or accepting stale records.
6. Keep `/agent` + `/status` as a manual operator diagnostic, not the protocol binding.

## Acceptance criteria

- **AC-1:** A spawned Codex worker's Spacedock identity maps deterministically to one app-server child thread ID and its latest token usage/window record.
- **AC-2:** The budget probe returns exact usage, model context window, remaining tokens, and remaining percentage for a running or completed-but-addressable worker.
- **AC-3:** Reuse fails safe when telemetry is absent, stale, ambiguous, or over the configured threshold; it never treats PRESENT-but-unavailable as satisfied.
- **AC-4:** A compaction spike proves which notification field represents active post-compaction context and prevents cumulative lifetime usage from producing a false over-budget decision.
- **AC-5:** Multiple sibling workers, reused turns, completed workers, and thread replacement cannot cross-attribute telemetry.
- **AC-6:** The fallback reads only metadata and token-count records, does not expose prompt/response content, rejects unsafe ownership/symlink paths, and names its freshness bound.
- **AC-7:** A live replay demonstrates that the 307,584/353,400-style worker is fresh-dispatched rather than reused while a comfortably under-budget sibling remains eligible.

## Test plan

First run an app-server v2 spike: spawn two named subagents, capture the collab receiver IDs, consume their `thread/tokenUsage/updated` events, compact one thread, and compare notification values with that thread's interactive `/status`. This must establish identity and active-window semantics before implementation.

Then add table tests for threshold boundaries, missing/stale/ambiguous telemetry, multiple siblings, compaction, completed-addressable reuse, and rollout fallback privacy/freshness. Finish with a live feedback-cycle replay proving over-budget fresh dispatch and under-budget reuse.
