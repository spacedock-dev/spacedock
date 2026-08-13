---
title: Acknowledge unstamped same-stage owner handoff
status: backlog
score: "0.90"
source: "n28 no-stamp validation finding, 2026-08-10"
sprint: live-evidence-followups
sprint-readiness: defer
group: common-evidence
id: egsrea0tppbaphb61kc9wj5s
---
## Problem

A fresh same-stage owner-handoff build without `--stamp` can emit a worker envelope without an acknowledgment receipt.

## Value

A same-stage owner handoff cannot spawn an unacknowledged worker. The fresh worker uses one binary-owned acknowledgment chain before workflow progress.

## Scope

- Own only the no-stamp same-stage conflict-handoff path.
- Preserve n28 stamped commissioned dispatch behavior.
- Do not change Pi, in-app hosts, general parallelism, recovery, or unrelated dispatch paths.
- Do not dispatch this task during the current Sonnet and Codex priority.

## Acceptance criteria

- AC-1: A fresh no-stamp same-stage handoff creates one pending envelope before native spawn.
- AC-2: Supported host hooks arm and consume the exact envelope for one native worker.
- AC-3: Missing or ambiguous acknowledgment blocks progress and a second fresh build.
- AC-4: Exact local host proof, full, race, registry, owner, and required PR checks pass before merge.

## Baseline evidence

- Released workflow: supported fresh Claude or Codex CLI same-stage owner handoff without `--stamp`.
- Observable harm: the build emits null acknowledgment fields and creates zero refs.
- Value authority: an owner handoff cannot spawn an unacknowledged worker.
- Trigger: n28 disposable no-stamp Codex build on candidate `042f926db` exited 0 with no acknowledgment refs.
