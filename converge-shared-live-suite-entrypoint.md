---
title: Converge common live journeys on one runtime-neutral test entry point
status: backlog
source: "Desired live-test registry common-journey semantics, 2026-08-03"
started:
completed:
verdict:
score: 0.95
worktree:
issue:
sprint: live-test-truth
group: common-runner
sprint-readiness: ready
id: h3b5tgk77vx9qqdmbjtpsh98
---

## Problem

The same common journeys currently enter through separate Claude and Codex top-level tests. This makes runtime selection part of journey identity and permits lane-specific scenario drift. The registry requires one TestLiveSharedScenarios/<journey-id> identity with transport differences behind adapters.

## Acceptance criteria

**AC-1 (VALUE)** A common journey has one test/subtest identity that is invoked unchanged by every implemented runtime lane.
Verified by: focused live invocations of the same named subtest through Claude and Codex that reach the same fixture and host-neutral assertion.

**AC-2** Runtime authentication, launch, output capture, and liveness remain adapter-owned and no common fixture or assertion is forked per host.
Verified by: unit coverage of adapter selection plus source inspection and focused live evidence.

**AC-3** Existing shared scenario parity remains complete while the separate TestLiveClaudeSharedScenarios and TestLiveCodexSharedScenarios entry points are retired or reduced to non-authoritative compatibility wrappers.
Verified by: the shared coverage tests and CI selector behavior.

## Stage-specific test gates

- Spike one shared journey end to end through the proposed unified entry point before migrating the suite.
- Validation runs focused Claude and Codex live evidence plus offline, full, and race suites.

