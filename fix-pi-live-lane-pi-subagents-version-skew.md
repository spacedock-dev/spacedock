---
id: 2686bggef0qz5hsrft2aks0t
title: Restore pi-live lane green by resolving the pi-subagents/pi-coding-agent version skew
status: implementation
source: c6 validation cycle-7 live-lane triage, 2026-07-20
started: 2026-07-20T13:56:42Z
completed:
verdict:
score:
worktree:
issue:
---

The pi-live CI lane is deterministically red for every commit, including main, due to an upstream extension version skew. Until it is fixed, no PR can satisfy the every-host-lane merge requirement on pi-live.

## Problem

`TestLivePiFrontDoorSmoke` fails at Pi extension load: `pi-subagents@0.35.1` (published 2026-07-18) requires the `@earendil-works/pi-ai` `/compat` export path, which `pi-coding-agent@0.74.2`'s vendored copy lacks. Evidence: identical failure on main's run 29691271679 (2026-07-19, pre-dating unrelated PRs' commits) and on PR #531's run 29709011051; last green pi-live run (2026-07-17) resolved `pi-subagents@0.34.0`. npm still serves 0.35.1, so reruns red identically — this is a deterministic break, not a flake, and re-run-to-green policy cannot apply.

## Proposed approach

{Ideation fills this in. Candidates: pin `pi-subagents@0.34.0` in the pi-live lane setup until pi-coding-agent vendors the `/compat` export; or bump pi-coding-agent when a release with the export ships; or both (pin now, unpin on bump).}

## Out of scope

{Ideation fills this in. Likely: any change to Pi runtime behavior itself — this is CI lane infrastructure.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - The pi-live lane is green on main.**
Verified by: a main-branch pi-live run passing after the change — an independent external result that can move the wrong way.

**AC-2 - The version constraint is explicit and evergreen, not incidental.**
Verified by: the lane setup declares the pinned/bumped versions with a comment stating the compatibility requirement, so the next skew is diagnosable from the file.

## Test plan

{Ideation fills this in. Expected: one pi-live lane run on a branch + one on main post-merge; no offline test surface.}
