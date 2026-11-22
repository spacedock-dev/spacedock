---
id: x9mfvnkqtbwgm4tba7dy41tw
title: Paseo worker record, FO restart adoption, and dispatch recovery
status: backlog
source: captain (2026-11-21)
started:
completed:
verdict:
score: 0.70
sprint: paseo-backend
worktree:
issue:
pr:
mod-block:
sprint-readiness: ready
---

Persist the Paseo worker handle in the task so a restarted or compacted first officer adopts a live worker instead of double-dispatching, and handle wait timeouts, lost create receipts, and dead workers without wedging unrelated tasks.

## Problem

With `paseo-dispatch-backend` alone, an FO that restarts, compacts, or crashes mid-stage has no durable record of the in-flight Paseo agent ID, risking a duplicate spawn on the same worktree/branch. Similarly, `paseo wait --timeout` exits 0 with `status=timeout` during normal progress, and an ended worker without a stage report or a lost create receipt needs explicit single-retry-then-hold handling.

## Proposed approach

1. **Worker record in task body/state:** Record host, server ID, agent ID, stage, worktree, and requested vs. returned provider/mode/model when `paseo run --no-wait` succeeds.
2. **FO restart adoption:** On engage/resume, inspect the recorded agent ID (and reconcile `paseo ls --json` by `spacedock.worker` / `spacedock.epoch` labels if the create receipt was lost). If running or idle with a committed report, adopt and verify instead of spawning a second worker.
3. **Fault isolation and retry:** Treat `paseo wait` timeout as live progress (no retry spent); on confirmed ended worker with no report, allow one fresh `-retry` dispatch, then hold only that task on second failure. Always archive verified or superseded workers.
4. **Live adoption proof:** Add `TestLivePaseoAdoptionSmoke` in `internal/ensigncycle/` (`-tags live`, using `codex/gpt-6-luna` by default, `SPACEDOCK_PASEO_BIN` overridable): start a worker, stop the first FO mid-stage, boot a fresh FO on the same workflow and isolated daemon, and assert total spawned workers == 1, durable Git commits verified, and worker archived.

## Risk evidence

Depends on `paseo-dispatch-backend`. The earlier probe (`_roadmap/paseo-backend/probe-report.md`) already verified `paseo ls` label lookup (`spacedock.worker`, `spacedock.epoch`), `paseo wait --timeout` returning exit 0 + `status=timeout` during real progress, and `paseo archive` closing verified workers.

## Out of scope

- Remote workers (`serverId`) and cross-host adoption.
- Fake ACP stub servers.
- Multi-FO concurrency leases (tracked separately by `79` / `entity-session-claim-lease`).

## Expected surface and tolerance

Estimate net LOC change: +200 to +350, across 4 to 7 files (`paseo-fo-dispatch.md`, any small Go helper/parser tests if needed, and `TestLivePaseoAdoptionSmoke`).

## Acceptance criteria

**AC-1 (VALUE) - A fresh FO adopts an in-flight Paseo worker after restart with zero duplicate dispatch.**
Verified by: `go test -tags live -count=1 -run '^TestLivePaseoAdoptionSmoke$' ./internal/ensigncycle -v` stops the first FO while a Paseo worker is running, boots a second FO on the same workflow and isolated daemon, and asserts `paseo ls --all` contains exactly 1 worker for the task, the task advances from the adopted worker's committed stage report, and the worker is archived. Spawning a second worker or failing to advance fails the test.

**AC-2 - Wait timeout, lost create receipt, and dead-worker recovery follow the single-retry-then-hold rule.**
Verified by: unit/contract tests on the worker-record reconciliation/recovery logic plus `TestLivePaseoAdoptionSmoke` confirm a wait timeout spends zero retries, an existing labelled worker prevents duplicate spawn, and a second failed attempt holds only the affected entity. Spending a retry on timeout or double-spawning fails the test.

## Test plan

- `go test ./...` and `go test ./... -race`.
- Targeted local live run: `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoAdoptionSmoke$' ./internal/ensigncycle -v`.

### Feedback Cycles
