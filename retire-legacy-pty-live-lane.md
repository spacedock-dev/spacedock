---
title: Retire the unreachable Claude pty team-mode live lane
status: backlog
source: "Desired registry keeps only the current merged-agent Claude substrate proof, 2026-08-03"
started:
completed:
verdict:
score: 0.85
worktree:
issue:
sprint: live-test-truth
group: runtime-specific
sprint-readiness: ready
id: b91c2hx2148wwy2451h3v9cr
---

## Problem

The Claude workflow installs a current merged host but still installs tmux and selects two legacy native-team pty tests that skip on that host. The desired registry retains the current merged-agent proof and excludes the obsolete transport.

## Acceptance criteria

**AC-1 (VALUE)** The Claude live lane spends no setup or execution time on a host regime it cannot exercise, while current merged-agent dispatch remains required evidence.
Verified by: the workflow no longer installs tmux or selects the two legacy pty tests, and TestLiveMergedTeamModeDispatch runs green.

**AC-2** Legacy pty-only tests, helpers, artifacts, and documentation are removed without deleting shared lifecycle fixtures still used by current tests.
Verified by: targeted reference inventory, workflow execution, and full/race tests.

**AC-3** docs/runtime-live-ci.md and the desired registry describe only the supported Claude live substrate.
Verified by: rendered documentation and the registry source bindings.

## Stage-specific test gates

- Ideation checks whether any pty helper remains used outside the two obsolete tests.
- Validation includes a current Claude merged-agent live run and workflow syntax validation.

