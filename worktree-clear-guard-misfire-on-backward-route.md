---
title: Clearing worktree on a backward route trips the terminal merge guard
status: backlog
source: "Live FO session, 2026-07-21, routing an implementation entity back to ideation."
id: 2kzefmqxg44vg3d7kx2nt1fy
---

`status --set <slug> worktree=` on an implementation-stage entity is refused by the terminal merge-hook guard ("cannot advance to terminal — merge hook(s) [pr-merge] have not run") even when no status change to a terminal stage is requested — observed live while routing ensign-finding-triage-disposition backward (implementation -> ideation, captain-directed reframe) where the worktree pointer referenced a sibling's worktree and clearing it was bookkeeping, not merge ceremony. The guard conflates worktree-clearing with terminalization; a backward feedback-style route is a legitimate non-terminal worktree clear. --force worked as the sanctioned bypass and is recorded here as the reproduction. Fix: scope the guard to actual terminal status transitions (or to worktree-clear only when the requested/current status is terminal), with a behavior test driving the backward-route case that fails on the current binary.
