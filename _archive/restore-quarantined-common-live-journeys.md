---
title: Restore quarantined common journeys after their product defects land
status: backlog
source: "Desired registry exposes required common cases hidden by tracked TODO skips, 2026-08-03"
started:
completed:
verdict:
score: 0.85
worktree:
issue:
sprint:
group: common-evidence
sprint-readiness: superseded
id: rmz8te1nt2a2akrkjfhaa37p
archived: 2026-08-09T15:00:00Z
---

## Problem

Required common journeys are selected through suite entry points but currently skip for one or more hosts: gate-guardrail, recorded-gate-lifecycle, rejection-flow, smallest-sufficient-mechanism, and keep-moving-posture. Their underlying product defects are owned by existing tasks, including work outside this sprint; this task owns only exact-tip restoration after those dependencies land.

This broad restoration task is superseded. Exact target-scoped owners in
`test-behavior-completeness` now own each remaining product repair and TODO
removal. The archived body remains as historical design context.

## Acceptance criteria

**AC-1 (VALUE)** Every named journey produces real live evidence on each supported runtime target rather than a green suite containing a tracked skip.
Verified by: exact-tip focused runs for every formerly quarantined journey and runtime target, followed by the complete common suites.

**AC-2** TODO removal follows demonstrated product repair and preserves every durable oracle and negative control.
Verified by: pre-removal reproduction, post-repair focused live evidence, and unchanged or strengthened assertions.

**AC-3** This task does not duplicate or reshape the product fixes owned by the referenced upstream tasks.
Verified by: a diff limited to quarantine removal, runner integration, and registry/runtime-live documentation unless the captain explicitly re-scopes it.

## Stage-specific test gates

- Remain deferred until every upstream product dependency is terminal or explicitly transferred.
- Ideation records the exact dependency-to-skip map without investigating another sprint member beyond its public dependency state.
- Validation requires clean, non-skipped live results rather than selector presence.
