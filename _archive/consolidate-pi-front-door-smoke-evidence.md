---
title: Consolidate Pi front-door and direct subagent smoke evidence
status: backlog
source: "Desired runtime-specific registry selects one lean Pi front-door substrate proof, 2026-08-03"
started:
completed:
verdict:
score: 0.85
worktree:
issue:
sprint:
group: runtime-specific
sprint-readiness: absorbed
id: 5p51c69b215nyyfm1x72ewzc
---

## Problem

TestLivePiFrontDoorSmoke is selected, while TestLivePiSubagentEnsignSmoke uses the same fixture but carries additional ensign boot-contract grading and is unselected. Wiring both would duplicate model spend rather than define one lean supported substrate proof.

## Acceptance criteria

**AC-1 (VALUE)** One selected Pi front-door smoke proves current-checkout installation, child-agent dispatch, durable output, and the ensign boot contract in a single run.
Verified by: the selected smoke fails when any of those four boundaries is broken and passes through the production Pi front door.

**AC-2** Duplicate direct-launch coverage is removed or reduced to deterministic helper coverage that does not spend another live run.
Verified by: the Pi selector and test inventory contain one runtime-specific live smoke for this fixture.

**AC-3** The common Pi journey adapter remains responsible for portable workflow semantics; this smoke stays limited to Pi substrate wiring.
Verified by: scope and assertion separation in the candidate.

## Stage-specific test gates

- Ideation compares the two current entry points and chooses the smallest merged oracle.
- Validation runs the resulting Pi smoke plus offline, full, and race suites.

