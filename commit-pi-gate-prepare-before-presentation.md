---
title: Commit the Pi gate before presentation
status: backlog
source: "Staff review M3 for test-behavior-completeness, 2026-08-09"
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
id: 2e4fe65gy9vcr4xck6akzmdd
---

A Pi operator must receive a gate review that is bound to committed state.

The task must start from strict XFAIL code `gate-prepare-state-commit-missing`. The repaired journey must prepare, commit, reread, and present in order.

A command helper or prose-only change does not satisfy this task. The real Pi `gate-guardrail` cell must pass.

Reuse the current gate room and commands. Do not change gate storage or command grammar.
