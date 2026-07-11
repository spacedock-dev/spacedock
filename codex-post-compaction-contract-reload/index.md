---
title: Codex post-compaction contract reload
status: ideation
source: Absorbed from task njr36mfyhbafy8zx9ydks8ep in another workflow; canonical handoff /tmp/first-officer-compaction-rehydration.md; captain directed repo-local absorption 2026-07-11
started: 2026-07-11T04:15:29Z
completed:
verdict:
score: 0.95
worktree:
issue:
id: c60nzb396vgf0f8a9v0sggwm
---

Codex can continue after context compaction with a summarized or stale understanding of the first-officer contract. The session must re-anchor itself to the current checkout contract before it mutates state, dispatches, presents a gate, or crosses a merge boundary.

## Problem

A launch-time bootstrap proves only initial loading. After compaction, the FO may retain an incomplete summary or an older contract version and continue without re-reading the current `skills/first-officer/SKILL.md`, its eager imports, and the Codex runtime adapter.

## Acceptance criteria

- **AC-1:** A forced-compaction Codex live scenario observes the FO re-read the current first-officer entry contract and Codex runtime adapter before its first post-compaction workflow action.
- **AC-2:** The post-compaction reload uses the current checkout/plugin source selected by the launcher, not an installed stale sibling.
- **AC-3:** Reloading does not replay the greet, duplicate dispatch, repeat a completed mutation, or reset durable workflow state.
- **AC-4:** The enforcement survives another compaction and fails closed when the current contract cannot be resolved.

## Test plan

Ideation must first prove which Codex surface can detect or reliably signal compaction. Use durable tool-call and workflow-state evidence in a real Codex session; transcript wording alone is insufficient.
