---
title: Continue Codex rejection after the first validation
status: backlog
source: "Staff review M2 for test-behavior-completeness, 2026-08-09"
started:
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree:
issue:
pr:
mod-block:
id: dvddbpsf4tdt3yjw1yjyp14k
---

A Codex operator must receive a complete correction journey after the first rejected candidate.

The task must start from strict XFAIL code `rejection-flow-not-completed`. It must then reach correction, validation/2, and one fresh final gate.

A recorder helper or feedback-flow component does not satisfy this task. The complete Codex journey is the required value.

Do not change recorder bytes, round format, or the Pi recorder-order repair.
