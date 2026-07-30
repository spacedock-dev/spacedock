---
title: Require Codex worker spawn after approved gate advance
status: ideation
source: "KD validation cycle 4, Codex keep-moving-posture at b60d1c8"
started:
completed:
verdict:
score: 0.95
worktree:
issue:
milestone: 0.27.0
id: f02j6dbnd4jakwczahv1tg2h
gates:
    version: 1
    current:
        gate: gate:f02j6dbnd4jakwczahv1tg2h:backlog
    records:
        - id: gate:f02j6dbnd4jakwczahv1tg2h:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:f02j6dbnd4jakwczahv1tg2h-backlog-1
              briefing:
                id: briefing:f02j6dbnd4jakwczahv1tg2h:backlog:attempt-1:revision-1
                digest: sha256:7deff3c5b9e3616c2f7619faf23da6168456c71d8cd968d5c68a3a73c08306f7
                digest-domain: canonical-bytes
                request-digest: sha256:bd94ed0c5d9047c300db27b7014a1289f1bedd84f037fa6ecf938bd72bf53150
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f02j6dbnd4jakwczahv1tg2h:backlog:1
                briefing: briefing:f02j6dbnd4jakwczahv1tg2h:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T23:39:29.432938Z"
                decision: approve
                reason: Captain conn approves immediate ideation because the exact supported Codex run already proves a Material worker-authority breach; the task owns diagnosis and the smallest runtime/contract correction, while KD remains frozen.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

A supported headless Codex keep-moving journey advanced `approved-gate`, built its implementation dispatch artifact, never invoked `spawn_agent` / `worker.spawn`, then read a report and terminalized the task. This violates the shipped dispatch boundary and silently bypasses worker execution and write-scope authority.

The exact evidence is retained at `/tmp/spacedock-validation-cycle4-b60d1c8.qPdIXO/codex/codex-shared-scenarios/keep-moving-posture/codex-exec.jsonl`: item 17 advances and builds; items 26–27 read and finalize; no worker identity or spawn exists. The live harness reported: `the FO advanced the approved entity "approved-gate" but did not dispatch its next stage`.

Ideation must find the smallest fix that makes the supported Codex journey actually spawn and await a worker after gate approval. It must not repair the observer, parse transcript/provider dialects, or let a dispatch commit or generated artifact substitute for an observed worker spawn.

## Acceptance criteria

**AC-1 — Approved work executes under worker authority.** In the supported Codex keep-moving journey, gate approval advances the task and is followed by a real addressable worker spawn before any implementation report read, completion credit, or terminalization.

**AC-2 — Missing spawn remains a hard failure.** A control that advances and builds the dispatch artifact but omits `worker.spawn` is rejected and cannot reach terminal state.

**AC-3 — The correction stays runtime- and observer-coherent.** The implementation uses the existing Codex dispatch binding and durable workflow boundary without adding transcript grammar, provider-event parsing, or a second observer protocol.
