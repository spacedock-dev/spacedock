---
id: f5yp12pqbp646dhrz6ayks4z
title: Make Survey portable across installed hosts and evidence providers
status: backlog
source: Captain-authorized filing from the shipped-surface containment audit, 2026-08-06
started:
completed:
verdict:
score: 0.90
group: shipped-surface-containment
worktree:
issue: spacedock-dev/spacedock#613
pr:
mod-block:
---

Make Survey work from an installed package on supported systems. Agentsview is one evidence provider, not a requirement.

## Problem

Survey requires the external `timeout` command, which stock macOS does not provide. It resolves installed resources from the consumer repository. Repository paths enter SQLite dot-command text without safe transport.

Issue #613 reports that Survey ran under Windows and WSL without Agentsview. The report does not establish what evidence remained available or whether Survey named the missing coverage.

## Proposed approach

Use packaged resource resolution and a portable bounded-process mechanism. Pass repository paths to SQLite as data. Define the minimum useful Survey result from repository-native evidence when Agentsview is unavailable. Name evidence-provider coverage in the result.

Exercise the installed package from a consumer repository that has no Spacedock source tree. Include stock macOS and the Windows/WSL topology from issue #613.

## Out of scope

This task does not replace `eqn` Cowork runtime detection or the existing Survey presentation and attribution tasks. It does not change Survey query definitions, which passed the audit.

Issue #613's dispatch path-mapping failure remains outside this task. This task owns only Survey behavior across that host boundary.

## Acceptance criteria

**AC-1 (VALUE) - Survey completes from an installed package in a supported consumer repository.**
Verified by: stock-macOS and Windows/WSL fixtures have no Spacedock source tree and still resolve packaged resources. Requiring `timeout` or a source checkout makes a run fail.

**AC-2 (VALUE) - Survey produces a useful and honest result when Agentsview is unavailable.**
Verified by: a fixture without Agentsview reports the defined repository-native evidence and names unavailable evidence. Restoring an Agentsview hard requirement or silently omitting coverage makes the check fail.

**AC-3 - Survey handles repository paths as data.**
Verified by: paths with spaces and SQLite-sensitive characters select only the intended repository. Restoring direct dot-command interpolation makes the negative control fail.

**AC-4 - The portable timeout preserves bounded Survey execution.**
Verified by: a stalled child stops within the declared budget and leaves no active child process. Removing the bounded-process mechanism makes the check fail.

## Test plan

Run the installed Survey entry point from temporary consumer repositories. Cover stock macOS, Windows/WSL paths, Agentsview present, and Agentsview absent. Use adversarial path values and a stalled-child fixture. Observe exit status, evidence coverage, selected repository, output, and child cleanup.
