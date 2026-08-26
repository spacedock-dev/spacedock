---
id: pkwmaawe57bkssqcv3twdv9m
title: One archive verb owns the move, the reason, the commit, and the publish
status: backlog
source: "Recurring FO toil, three occurrences on 2026-08-25/26 (two in this repo's session, one in the email-triage run): state commit refuses archived scope as publish-only, so every archive is a hand-rolled sequence of --set archived, a body note, git mv of the file and its room dir, a path-scoped git commit, and a publish"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

The write scope grants the FO archive moves. The tooling does not. Each archive is five manual steps, and two of them (the room-dir move, the path-scoped commit) are easy to get wrong in a shared checkout.

## Problem

{Ideation fills this in. Seeded: `state commit <slug>` deliberately refuses archived scope ("publish-only"); the manual sequence is --set archived=TS, append the supersede reason, git mv the entity file AND its companion room dir, path-scoped git add/commit, then state commit to publish. The flat-entity-plus-room-dir pair is the sharp edge.}

## Proposed approach

{Ideation fills this in. Seeded: `spacedock archive <slug> --reason TEXT [--verdict PASSED|REJECTED]` — stamps archived, records the reason in the body, moves the entity and its companion dir to _archive/, commits path-scoped, and publishes. Supersede convention (empty verdict) is the default.}

## Risk evidence

{Backlog: three hand-rolled occurrences in two days decide design should start.}

## Out of scope

Un-archive. Changing the publish-only rule for ad-hoc archived edits.

## Expected surface and tolerance

Estimate: production +60 across 3 files; proof +80 across 2 files. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - One command archives a flat or folder entity with its rooms: the entity lands in _archive/, the reason is recorded, and the state branch is published — with no manual git step.**
Verified by: {ideation refines — seed: a fixture test archiving a flat entity with a room dir, asserting the final tree, the commit's pathspec, and the pushed ref. Falsifying edit: drop the room-dir move — the tree assertion reds.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
