---
id: e87dyx4ytp64g8bh5a5fat1a
title: A durable record that a captain-only act is done
status: backlog
source: "MetaFO friction list 2026-11-24, sent on the word of CL, /tmp/spacedock-frictions-2026-11-24.md, group C item 11"
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

When only the captain can do an act (a post, a send), the workflow records that it was done, so a pane or a later first officer does not ask again.

## Problem

Gate decisions are only `approve`, `revise` or `hold` (`internal/gates/model.go:352-358`). An outward act that only the captain can do has no entity form, so the Ready Room pane keeps it as view state.

## Proposed approach

{Ideation fills this in. Ideation decides whether this is a gate decision, a frontmatter field or a stage.}

## Risk evidence

Backlog: the decision set is confirmed in the source above.

## Out of scope

Doing the act on the captain's behalf.

## Expected surface and tolerance

Estimate net LOC change: {ideation}, across {M} files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - After the captain records an act as done, a fresh read of the entity reports it as done, and the first officer does not present it again.**
Verified by: {ideation names a fixture}.

## Test plan

{Ideation fills this in.}

### Feedback Cycles
