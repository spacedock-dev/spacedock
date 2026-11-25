---
id: 1g5zccpxcf00gpszsdneg33n
title: An ensign assignment form that does not push
status: backlog
source: "MetaFO friction list 2026-11-24, sent on the word of CL, /tmp/spacedock-frictions-2026-11-24.md, group B item 6"
group: B-paseo-host-backend
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

A worker can run in a project that forbids pushes, or in a blind run, and still finish its stage correctly.

## Problem

The split-root assignment tells an ensign to commit and then push the state branch (`internal/dispatch/build.go:925-940`, `skills/ensign/references/ensign-shared-core.md:40`). A carve-out exists only when the state checkout has no `origin` (`ensign-shared-core.md:42`). A project with an origin but a "no push" rule, or a blind run, has no form of the assignment that commits without pushing.

## Proposed approach

{Ideation fills this in.}

## Risk evidence

Backlog: the push reminder and the no-origin carve-out are confirmed in the source above. Seen in use on the Studfleet VM.

## Out of scope

Who later pushes the commits; the Paseo transport.

## Expected surface and tolerance

Estimate net LOC change: {ideation}, across {M} files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - With the no-push form selected, a dispatched ensign's assignment contains no push instruction and the stage completes with a local path-scoped commit.**
Verified by: a dispatch build test on a split-root fixture with an origin: the no-push assignment has no `git push`; without the form, it does. Removing the form's branch makes the first case fail.

## Test plan

{Ideation fills this in.}

### Feedback Cycles
