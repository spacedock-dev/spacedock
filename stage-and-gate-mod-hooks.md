---
id: 1q57y1hf3ys5wvhaxgenqkc3
title: Mod hooks at stage completion and before a gate
status: backlog
source: "MetaFO friction list 2026-11-24, sent on the word of CL, /tmp/spacedock-frictions-2026-11-24.md, group C item 10"
group: C-extension-points
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

A workflow can add a deterministic check (for example Taskless before validation, or a spec validation before ideation) as a mod, not as prose in its README.

## Problem

The first officer contract supports exactly three mod hook points, `startup`, `idle` and `merge` (`skills/first-officer/references/first-officer-shared-core.md:120-123`). None fires at stage completion, before a gate, or at a rejection. The Go scanner records any `## Hook:` name (`internal/status/mutate.go:516-560`), but only the three run. Today each such check is a README line, as the Taskless spike proposes (`.worktrees/taskless-spike/spike/taskless-report.md`, branch `spike/taskless-ste-rules`).

## Proposed approach

{Ideation fills this in. The contract treats a new standing check as a last resort with captain approval; ideation must show why a README line in the stage definition is not enough.}

## Risk evidence

Backlog: the three hook points are confirmed in the source above.

## Out of scope

Any specific check (Taskless, OpenSpec).

## Expected surface and tolerance

Estimate net LOC change: {ideation}, across {M} files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A mod registered at the new point runs once when a worker's stage report passes the completion gate, before the first officer prepares that stage's gate.**
Verified by: {ideation names a fixture or live journey with an observable side effect of the hook}.

## Test plan

{Ideation fills this in.}

### Feedback Cycles
