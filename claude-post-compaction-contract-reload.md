---
id: cdbhzxcees9hmdqx2zk2w3p3
title: Claude post-compaction contract reload delivery
status: backlog
source: "Split from codex-post-compaction-contract-reload (host-neutral reframe); captain directed 2026-07-18."
started:
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
---

Claude-specific delivery binding for the host-neutral post-compaction reload rules defined in `codex-post-compaction-contract-reload`. On Claude, unlike Codex, the reminder can reach the model automatically.

## Problem

The parent task (`codex-post-compaction-contract-reload`) defines two host-neutral FO rules — suggest compaction only at a durable boundary, and reread the authoritative `spacedock:first-officer` contract + reconcile durable state after a compaction. Codex can only deliver the post-compaction reminder as a captain-facing UI warning (`systemMessage`; the live probe found the model never received it), so Codex falls back to a manual captain cue.

On Claude Code the delivery can be stronger: a `SessionStart` hook with the `compact` matcher writes stdout that IS injected into the post-compaction model context (Claude docs' "re-inject context after compaction" pattern). That would make the reload reminder reach the FO automatically rather than depending on a captain cue — the Claude-specific piece that the parent task deliberately scopes out.

## Proposed approach

Bundle a Claude `SessionStart(compact)` hook (matching manual + auto) that emits the reread-and-reconcile instruction to stdout so it enters the post-compaction model context. `PreCompact` optional. Requires a hook-shipping surface the plugin does not have yet (`.claude-plugin/plugin.json` has no `hooks` key today). Failure-open: hooks need trust and can be disabled — on absence, the host-neutral FO rule + manual captain cue (from the parent task) still applies, so nothing is blocked and no Spacedock state file/process is created.

## Out of scope

- The host-neutral FO rules themselves (owned by `codex-post-compaction-contract-reload` / the shared core).
- Codex and Pi delivery bindings.
- Any continuation controller / ledger / watchdog (rejected in the parent task).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A trusted bundled `SessionStart(compact)` hook injects the reread-and-reconcile instruction into the post-compaction MODEL context (not just the UI).**
Verified by: the live spike below — after compaction, the next model turn (no tools) reproduces the sentinel rather than answering NONE. Confirmed for both manual `/compact` and auto-compact.

**AC-2 - Harmless absence.**
Verified by: with the hook disabled/untrusted/absent, Claude compaction and the next turn continue with no Spacedock-created state file, background process, blocked stop, or workflow mutation; the manual captain cue still triggers the reload.

## Spike evidence

Riskiest mechanism: does Claude `SessionStart(compact)` `additionalContext`/stdout actually reach the post-compaction model (the exact thing the Codex probe found unsupported for `systemMessage`)? To be exercised on model `claude-sonnet-4-6` at 200k context, mirroring the Codex 0.144.4 probe methodology. {Spike worker fills this in.}

## Test plan

{Ideation refines. Likely: a hook fixture asserting the exact stdout JSON on manual+auto, an absence/disabled/failing matrix (AC-2), and the opt-in live sonnet-4-6/200k probe confirming model-context delivery (AC-1).}
