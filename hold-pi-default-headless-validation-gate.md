---
title: Hold the Pi default headless validation gate
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
id: fh6rv0k6wr25zty0jjan4jp7
---

A headless Pi run must stop at the first open validation gate with durable evidence.

Ideation must identify the exact failed final-state clause before it selects a repair. Do not assume the Sonnet or Codex mechanism.

A fixture or assertion change does not satisfy this task. The unchanged Pi `default-headless-gate-stop` journey must pass.

The task depends on strict XFAIL and on the `98a` worker-dispatch repair.
