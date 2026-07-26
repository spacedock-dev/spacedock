---
title: Make ensign dispatch self-contained across launcher drift
status: backlog
sprint: durable-decisions
source: "Real source-build sprint dogfood, 2026-07-26"
score: "1.0"
id: kd7877nnbd19d528xnpwwaj4
---

## Outcome

A dispatched ensign receives the exact workflow stage contract selected by the First
Officer's launcher without depending on inherited `SPACEDOCK_BIN` or a different
`spacedock` on `PATH`. Workflow-control mutations remain First Officer-owned; running
the worktree binary as the product under test remains separate.

## Observed failure

The First Officer invoked `/tmp/spacedock-s4-73810d1e` (`0.27.0-pre1`), while the
Codex process environment exposed
`SPACEDOCK_BIN=/Users/clkao/git/spacedock-research/spacedock-v1/spacedock`
(`0.26.0+dev`) and PATH exposed Homebrew `0.27.0-pre0`. `spawn_agent` has no
environment argument. `dispatch build` nevertheless emitted a mandatory ensign fetch
command that re-resolved `${SPACEDOCK_BIN:-spacedock} dispatch show-stage-def`.
Following the generated assignment therefore read the stage through a different
binary than the First Officer used. The command happened to exist, masking the drift.

`dispatch build` already reads and validates the requested stage subsection before it
writes the assignment, so the second worker-side binary call is redundant.

## Initial acceptance criteria

**AC-1 (VALUE)** A two-binary fixture invokes `dispatch build` through binary A while
the worker environment and PATH point to binary B; the ensign assignment still carries
the exact stage subsection selected by A and requires no worker-side
`dispatch show-stage-def`.

**AC-2** The generated dispatch file is self-contained for the stage contract on both
fresh dispatch and reuse advance. The ensign reads the dispatch file and entity, then
can perform ordinary stage work without a workflow-control Spacedock binary.

**AC-3** Every remaining intentional ensign-side Spacedock helper is inventoried and
either moved to the First Officer or mechanically bound to the exact executing
launcher by the dispatch builder. No caller-authored launcher flag, global environment
mutation, plugin-private wrapper, or compatibility layer is added.

## Boundary

Prefer embedding the already-resolved stage subsection over introducing environment
transport. Do not fold this into s4 gate-room behavior or v21 source-build version
identity. The durable-decisions walking skeleton and pre-release wait for this
launcher-consistent dispatch path.
