---
id: f5yp12pqbp646dhrz6ayks4z
title: Make Survey portable from an installed package
status: backlog
source: Captain-authorized filing from the shipped-surface containment audit, 2026-08-06
started:
completed:
verdict:
score: 0.90
group: shipped-surface-containment
worktree:
issue:
pr:
mod-block:
---

Make Survey find packaged resources and run safely on supported systems without access to the Spacedock source repository.

## Problem

Survey requires the external `timeout` command, which stock macOS does not provide. It resolves installed resources from the consumer repository. Repository paths enter SQLite dot-command text without safe transport.

## Proposed approach

Use packaged resource resolution and a portable bounded-process mechanism. Pass repository paths to SQLite as data. Exercise the installed package from a consumer repository that has no Spacedock source tree.

## Out of scope

This task does not replace `eqn` Cowork runtime detection or the existing Survey presentation and attribution tasks. It does not change Survey query definitions, which passed the audit.

## Acceptance criteria

**AC-1 (VALUE) - Survey completes from an installed package in a supported consumer repository.**
Verified by: a stock-macOS-compatible fixture has no Spacedock source tree and still resolves packaged resources. Requiring `timeout` or a source checkout makes the run fail.

**AC-2 - Survey handles repository paths as data.**
Verified by: paths with spaces and SQLite-sensitive characters select only the intended repository. Restoring direct dot-command interpolation makes the negative control fail.

**AC-3 - The portable timeout preserves bounded Survey execution.**
Verified by: a stalled child stops within the declared budget and leaves no active child process. Removing the bounded-process mechanism makes the check fail.

## Test plan

Run the installed Survey entry point from a temporary consumer repository. Use a stock-macOS command environment, adversarial path values, and a stalled-child fixture. Observe exit status, selected repository, output, and child cleanup.
