---
id: 2686bggef0qz5hsrft2aks0t
title: Restore pi-live lane green by resolving the pi-subagents/pi-coding-agent version skew
status: implementation
source: c6 validation cycle-7 live-lane triage, 2026-07-20
started: 2026-07-20T13:56:42Z
completed: 2026-07-20T14:19:37Z
verdict: implemented
score:
worktree: .worktrees/spacedock-ensign-fix-pi-live-lane-pi-subagents-version-skew
issue:
---

The pi-live CI lane is deterministically red for every commit, including main, due to an upstream extension version skew. Until it is fixed, no PR can satisfy the every-host-lane merge requirement on pi-live.

## Problem

`TestLivePiFrontDoorSmoke` fails at Pi extension load: `pi-subagents@0.35.1` (published 2026-07-18) requires the `@earendil-works/pi-ai` `/compat` export path, which `pi-coding-agent@0.74.2`'s vendored copy lacks. Evidence: identical failure on main's run 29691271679 (2026-07-19, pre-dating unrelated PRs' commits) and on PR #531's run 29709011051; last green pi-live run (2026-07-17) resolved `pi-subagents@0.34.0`. npm still serves 0.35.1, so reruns red identically — this is a deterministic break, not a flake, and re-run-to-green policy cannot apply.

## Proposed approach

Bump only the pi-live job to Node 24, install exact pinned Pi substrate tarballs (`@earendil-works/pi-coding-agent@0.80.10`, `pi-subagents@0.35.1`, `pi-intercom@0.6.0`) after checking their recorded npm `dist.integrity` sha512 values, and add a fast compatibility guard for the pi-ai `/compat` export required by pi-subagents 0.35.x.

## Out of scope

Pi runtime behavior changes and unrelated live lanes' Node versions.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - The pi-live lane is green on main.**
Verified by: a main-branch pi-live run passing after the change — an independent external result that can move the wrong way.

**AC-2 - The version constraint is explicit and evergreen, not incidental.**
Verified by: the lane setup declares the pinned/bumped versions with a comment stating the compatibility requirement, so the next skew is diagnosable from the file.

## Test plan

- `go test ./...`
- `go test ./... -race`
- Local npm smoke of the pinned tarball/integrity flow and pi-ai `/compat` guard using temporary install prefixes.
- Primary live evidence remains the PR pi-live run, then a main pi-live run after merge.

### Summary

Implemented the pi-live workflow fix: Node 24 for the Pi lane, exact pinned Pi substrate versions with recorded sha512 integrity checks before install, installed tarballs instead of bare latest selectors, and a deterministic compatibility guard that fails fast if pi-subagents cannot use the pi-ai `/compat` export. Updated release workflow guard tests to enforce the new pinned/integrity-based install contract.
