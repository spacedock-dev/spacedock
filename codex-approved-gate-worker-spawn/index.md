---
title: Require Codex worker spawn after approved gate advance
status: backlog
source: "KD validation cycle 4, Codex keep-moving-posture at b60d1c8"
started:
completed:
verdict:
score: 0.95
worktree:
issue:
milestone: 0.27.0
id: f02j6dbnd4jakwczahv1tg2h
---

A supported headless Codex keep-moving journey advanced `approved-gate`, built its implementation dispatch artifact, never invoked `spawn_agent` / `worker.spawn`, then read a report and terminalized the task. This violates the shipped dispatch boundary and silently bypasses worker execution and write-scope authority.

The exact evidence is retained at `/tmp/spacedock-validation-cycle4-b60d1c8.qPdIXO/codex/codex-shared-scenarios/keep-moving-posture/codex-exec.jsonl`: item 17 advances and builds; items 26–27 read and finalize; no worker identity or spawn exists. The live harness reported: `the FO advanced the approved entity "approved-gate" but did not dispatch its next stage`.

Ideation must find the smallest fix that makes the supported Codex journey actually spawn and await a worker after gate approval. It must not repair the observer, parse transcript/provider dialects, or let a dispatch commit or generated artifact substitute for an observed worker spawn.

## Acceptance criteria

**AC-1 — Approved work executes under worker authority.** In the supported Codex keep-moving journey, gate approval advances the task and is followed by a real addressable worker spawn before any implementation report read, completion credit, or terminalization.

**AC-2 — Missing spawn remains a hard failure.** A control that advances and builds the dispatch artifact but omits `worker.spawn` is rejected and cannot reach terminal state.

**AC-3 — The correction stays runtime- and observer-coherent.** The implementation uses the existing Codex dispatch binding and durable workflow boundary without adding transcript grammar, provider-event parsing, or a second observer protocol.
