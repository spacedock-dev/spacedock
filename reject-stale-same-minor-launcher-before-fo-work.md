---
title: Reject a stale same-minor launcher before First Officer work
status: backlog
source: "Weekly friction audit and Captain direction, 2026-08-09: an installed 0.27 binary passed the minor-version gate but lacked the merged approval surface. Detect the missing command capability before workflow effects."
started:
completed:
verdict:
score: 0.85
sprint: durable-decisions
sprint-readiness: ready
group: launcher-contract
worktree:
issue:
pr:
mod-block:
id: 5f6m3jwhbrbneak5j8eeyh5r
---

Stop a First Officer before workflow work when the selected launcher has the correct minor version but lacks a required command surface.

## Problem

The startup gate accepts any `0.27` launcher. During the approval migration, the installed `0.27` binary passed this gate but lacked the current gate help and behavior.

The First Officer then used an obsolete approval path and discovered the mismatch after state work started.

## Required outcome

Define the smallest command-capability check that distinguishes a compatible `0.27` launcher from a stale `0.27` launcher.

Run this check after launcher selection and before boot or state discovery. If the check fails, stop and name the selected launcher, the missing capability, and the normal upgrade remedy.

Do not tell a user to build Spacedock. Do not refer to a consumer repository or source-build workflow.

Do not add network version lookup, repository identity, commit identity, or a second launcher-resolution path unless ideation proves that command capability is insufficient.

## Acceptance criteria

**AC-1 (VALUE) - A stale same-minor launcher cannot start workflow work.**
Verified by: a fixture launcher reports `0.27` but lacks one required gate capability. Startup stops before `status --boot` and names the missing capability.

**AC-2 - A compatible installed launcher continues normally.**
Verified by: a fixture with the same version and required capability reaches the existing single boot call without an additional launcher selection.

**AC-3 - The remedy describes installation, not development.**
Verified by: failure output names the supported upgrade command for the host and contains no source-build or consumer-repository instruction.

## Test plan

Use fixture launchers with the same version and different command capabilities. Observe command order, output, and exit status.

Run the existing First Officer live boot journey after the focused checks. Do not add a static prose-presence test.

