---
id: 95bcs48mr3jemtsb2zq7zbtb
title: FO contract residual-prose audit + comm-officer polish (post-split)
status: ideation
source: 0.20.3 / 0203-fo-efficiency sprint (T3); captain 2026-06-13
started: 2026-06-13T18:09:33Z
completed:
verdict:
score:
worktree:
issue:
---

The T3 of the 0.20.3 (0203-fo-efficiency) sprint: after j9 Phase-1 splits the FO contract references into boot-resident vs deferred modules, sweep the slimmed refs for residual dead/redundant prose and run a comm-officer light-touch polish pass.

## Dependency

BLOCKED on j9 Phase-1 (the contract structural split). The audit cut-list does not exist until the structure is split — do NOT dispatch ideation until Phase-1 lands. Stays a backlog seed until then.

## Notes

Behavior-preserving CLEANUP: tightened/deduped prose in the already-split refs, NOT the structural split (j9 Phase 1) and NOT a behavioral contract change. At ideation, prove a real checkable change — the existing live FO scenarios stay green (behavior-preserving) plus a measurable size reduction of the refs — never a review of its own prose. If the split left no residual prose to cut, this collapses to a recorded roadmap decision, not a shipped task.
