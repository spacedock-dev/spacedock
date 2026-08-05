---
id: w9h7bzax0qdmtmp8r5eepk0h
title: Commission runnable workflows for each supported runtime
status: backlog
source: Captain-authorized filing from the shipped-surface containment audit, 2026-08-06
started:
completed:
verdict:
score: 0.95
group: shipped-surface-containment
worktree:
issue:
pr:
mod-block:
---

Make Commission generate a runnable workflow for each declared runtime adapter without Claude-only paths or shared temporary state.

## Problem

Commission always emits `spacedock claude` guidance and requires Claude-only team tools and paths. Generated shell text does not quote project paths. A fixed `/tmp/orphan-birth` path makes concurrent runs collide. The development template also assumes `main` and carries a repository-local ruling.

## Proposed approach

Generate launch and dispatch guidance from the selected declared runtime. Use safe path transport and unique temporary resources. Keep the generated development template portable and use its configured trunk.

## Out of scope

This task does not add a runtime adapter or change a runtime's lifecycle semantics. It does not repair First Officer, Debrief, Survey, or `mods/pr-merge.md`. Existing host-detection tasks retain their narrower runtime boundaries.

## Acceptance criteria

**AC-1 (VALUE) - A commissioned workflow starts through each declared runtime adapter.**
Verified by: generated Claude, Codex, and Pi packages use their native launch contracts. A fixture-backed adapter check proves shape. A live smoke proves runtime integration when the task changes that runtime path.

**AC-2 - Generated workflows work in project paths that contain spaces.**
Verified by: a temporary project with a spaced path commissions and starts without shell argument changes. Removing safe path transport makes the run fail.

**AC-3 - Concurrent commission runs do not share scratch state.**
Verified by: two simultaneous fixture runs complete with distinct temporary resources. Restoring the fixed scratch path makes the check fail.

**AC-4 - The generated development template uses declared workflow configuration.**
Verified by: a non-`main` fixture uses its configured trunk and contains no local Captain ruling. Restoring either hard-coded value makes the check fail.

## Test plan

Exercise generation for every declared adapter. Start the generated workflow through the cheapest real front door for each changed runtime. Use a spaced path and two concurrent runs as focused negative controls.
