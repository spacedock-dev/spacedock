---
id: sbkske5gyadwfxs9vmd61zb3
title: Make the First Officer portable across supported consumer repositories
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

Make the shipped First Officer operate from a supported project that consumes Spacedock and does not contain the Spacedock source tree.

## Problem

Shared First Officer instructions assume this repository's Go layout, the literal `main` branch, and plugin-private source paths. Some recovery paths try to build Spacedock from the consumer project. These assumptions can misclassify files, select the wrong trunk, or fail when the project only has an installed launcher.

## Proposed approach

Use resolved workflow roots, the configured trunk, and the launcher selected at the version gate. Define the supported consumer layouts from existing product contracts. Exercise one end-to-end First Officer journey in a temporary repository that has no Spacedock source tree.

## Out of scope

This task does not replace `qx` state-ready path resolution, `z93` Claude launch-directory ownership, `kd` self-contained ensign dispatch, or `kbh` and `88c` roadmap classification. It does not change host adapter semantics, Commission, Debrief, Survey, or `mods/pr-merge.md`.

## Acceptance criteria

**AC-1 (VALUE) - A First Officer completes the declared journey from a supported consumer layout.**
Verified by: a temporary consumer repository uses a non-`main` trunk, split-root state, and the installed launcher. The run verifies durable workflow state. Removing any resolved-root or configured-trunk binding makes the run fail.

**AC-2 - Shared First Officer behavior does not require the Spacedock source tree.**
Verified by: the same consumer fixture contains no `cmd/spacedock` or `internal` source. A launcher or recovery path that tries a consumer-repository source build makes the run fail.

**AC-3 - File authority follows declared workflow roots instead of this repository's product layout.**
Verified by: behavior fixtures classify equivalent process and product paths in two different consumer layouts. Reintroducing a `docs/dev` or Go-layout assumption makes a fixture fail.

## Test plan

Start with the smallest end-to-end consumer fixture that exercises boot, root resolution, one state mutation, and launcher reuse. Add focused behavior tests. Run the applicable live runtime lane only when shared host behavior changes.
