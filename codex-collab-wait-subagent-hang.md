---
id: czza18qnjzj75fznszxm3z0s
title: Codex rejection-flow collab:wait spawned-subagent stall (live-infra flakiness)
status: ideation
source: "j9 validation (2026-06-14) — the live Codex rejection-flow hung at collab:wait on 3/5 runs across j9's validation cycles (a spawned-subagent stall in the Codex runner). Orthogonal to j9 (the validator retried past it to a clean run); flagged by the validator for its own triage."
started: 2026-06-14T02:02:43Z
completed:
verdict:
score:
worktree:
issue:
sprint:
---

The Codex live `rejection-flow` shared-runtime scenario intermittently hangs at `collab:wait` — a spawned-subagent stall in the Codex runner, observed on 3 of 5 runs across j9's validation cycles. When it hangs the assertion is never reached (so it is NOT a scenario regression), but it makes the Codex live lane flaky and forces retries.

## Problem

The Codex runner spawns the reviewer subagent and waits on `collab:wait`; on ~60% of runs the wait stalls indefinitely with no assertion reached. A retry usually yields a clean run. Live-infra flakiness, not a behavioral defect — but it degrades the Codex live lane's reliability and burns retries (and would be a silent CI hang without a watchdog).

## Proposed approach (seed)

Triage the `collab:wait` stall: Codex CLI / pi-subagents issue, a timeout/quiet-budget gap in the runner, or resource contention under the serial suite. Add a bounded retry or a per-stage quiet-budget watchdog so a stall is detected and retried/failed cleanly rather than hanging. Ideation determines root cause + fix.

## Out of scope

The j9 deliverable — the validator retried past the hang and j9's Codex AC-3 passed on a clean run.

## Acceptance criteria (seed — ideation defines external proof)

A Codex rejection-flow run that no longer stalls indefinitely at `collab:wait`, or whose stall is detected and bounded-retried within a quiet budget — proven by a live Codex run, or a deterministic reproduction of the stall plus the watchdog catching it (an external timeout/exit-code oracle, not a prose claim).
