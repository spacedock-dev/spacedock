---
title: Make Codex implementation dispatch mechanically acknowledged
status: backlog
source: "98a cycle-4 live evidence: identical spawn envelopes produced two native-spawn passes and one empty-wait failure"
started:
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
                state: pending
---

## Problem

Codex can receive a valid implementation spawn envelope and continue without a native worker handle. Later reports and empty waits can look like completion, so validation can start without the required implementation lifecycle.

The current prose-only dispatch contract is not deterministic. Identical fixture, prompt, adapter, and envelope bytes have produced both passing native spawns and a missing-spawn failure.

## End-user value

A Codex user gets one real implementation worker, its bound handle, and its completion receipt before validation starts. The completed workflow, not a component API, must pass the existing default-headless journey without XFAIL.

## Constraints

- Keep the existing `spacedock dispatch build` command surface and stored workflow formats unless ideation proves a change is necessary.
- The binary owns the acknowledgment boundary. Do not add another prose guard or simulator.
- Status advance and `--advance` must fail closed when a spawn envelope lacks a matching native handle and completion receipt.
- Preserve Sonnet, Pi, gate authority, and existing worker ownership behavior.
- Remove only this task's Codex default-headless XFAIL after exact passing product evidence.

## Seed acceptance criteria

**AC-1 (VALUE) — Codex implementation dispatch completes before validation.**

The exact Codex `default-headless-gate-stop` journey records one native implementation spawn, its matching handle and completion, and a DONE implementation report before validation. The journey passes without XFAIL.

**AC-2 — Missing acknowledgment fails closed.**

An envelope followed by an empty wait, report read, status advance, or `--advance` without the matching handle and completion receipt is refused before validation.

**AC-3 — The mechanism is end to end.**

The delivered change includes the supported Codex workflow path and exact live proof. A component-only acknowledgment API or isolated unit test does not satisfy the task.

## Ideation assignment

Trace the current envelope, native handle, completion, report, and status boundaries. Spike the smallest binary-owned acknowledgment through the real Codex path before selecting a design. Declare exact files, gross/net lines, tolerance, semantic changes, falsifiers, and the XFAIL-first removal order.
