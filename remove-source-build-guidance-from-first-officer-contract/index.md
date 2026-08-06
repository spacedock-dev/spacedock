---
id: xfde00anxs2khetvfdm5kwdn
title: Remove source-build guidance from the First Officer contract
status: backlog
source: "Captain correction, 2026-08-06: the shipped First Officer contract must not tell users or agents to build Spacedock from source"
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

Make every shipped First Officer path use an installed launcher or a supported installation method. The shared contract must not suggest or execute a Spacedock source build.

## Problem

The binary-absent startup text names `go build -o spacedock ./cmd/spacedock` as a fallback. The Claude reconciliation text also runs that command after local trunk drift. These source-tree assumptions do not belong in the shipped First Officer contract.

## Proposed approach

Remove source-build guidance from the binary-absent and install-gate outcomes. Keep the supported Linux and macOS installation methods. Remove source rebuilding from generic reconciliation. Preserve any Spacedock-development-only rebuild behavior in the development workflow only when that workflow still needs it.

## Out of scope

This task does not redesign installation, add an installer, change Cowork bootstrap task `8v`, or repair generic consumer-repository portability. It does not change runtime dispatch semantics.

## Acceptance criteria

**AC-1 (VALUE) - A binary-absent First Officer gives only supported installation guidance.**
Verified by: binary-absent replays for Linux, macOS, sandboxed, and unsupported systems contain no source-build instruction. Restoring the source-build fallback makes the replay fail.

**AC-2 (VALUE) - Generic reconciliation never builds Spacedock from the current workflow repository.**
Verified by: a local-trunk-drift replay updates only the declared repository state and invokes no Go build. Restoring the build action makes the command-event assertion fail.

**AC-3 - Spacedock development behavior remains explicit and local.**
Verified by: if the development workflow still requires a rebuild after its own trunk moves, its workflow-owned hook or process record proves that behavior without adding it to the shipped shared contract.

## Test plan

Extend the existing version-gate and reconciliation behavior fixtures. Assert command events and user-visible remediation output. Use one negative control that restores each source-build path. Run the applicable contract and integration tests.
