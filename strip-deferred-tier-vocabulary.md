---
id: ga1e7heckpazwbmtxzrg2kms
title: Strip deferred-72 tier vocabulary leaked into czw's «gate.assemble-verdict» body (pre-cut ship-blocker)
status: implementation
source: '0221 pre-cut antipattern audit (2026-06-19): czw (prose-function restructure) leaked the DEFERRED 72 member''s tier vocabulary into the shipped contract — first-officer-shared-core.md:127,130 instruct a level-2-only FO to escalate the verdict to a level-3 judge, the ONLY occurrence of level/L3/judge in the contract, with nothing defining the concept (72 deferred). Ship-blocker for the v0.22.1 tag. Audit also recommended a contractlint absence-guard to prevent recurrence.'
started: 2026-06-19T21:31:06Z
completed:
verdict:
score: 0.5
worktree: .worktrees/spacedock-ensign-strip-deferred-tier-vocabulary
issue:
sprint: 0221-layered-fo
sprint-readiness: ready
group: cleanup
---

The prose-function restructure smuggled the deferred 72 member's tier vocabulary into the most safety-critical contract decision (the gate verdict). Strip it back to mechanism-neutral phrasing and add a structural guard so deferred-member vocabulary cannot leak until its member ships.

## Problem

czw re-expressed `«gate.assemble-verdict»`'s verdict-decision body as "a level-2-only FO escalates the decision to a level-3 judge; a capable FO renders its own Recommend line" (first-officer-shared-core.md ~127) and "routes to L3 or the FO's own present-gate Recommend line" (~130). This is 72 (fo-tier-delegation)'s tier/level-3-judge mechanism — DEFERRED. Nothing in the shipped contract defines an FO "level" or a "level-3 judge". A booting FO reading this is instructed to escalate to a non-existent judge: a dangling reference. The base (2c9609f8) phrasing was mechanism-neutral (the FO renders its own Recommend line).

## Proposed approach (the pre-cut audit IS the design)

1. Strip the level-2-only / level-3 judge / L3 vocabulary from the `«gate.assemble-verdict»` body; revert to the base mechanism-neutral phrasing (the FO renders its own `Recommend` line / present-gate verdict). Two prose edits, no code, no change to czw's «fn» notation / routing oracle / the six behavioral clusters.
2. Add a contractlint structural-absence guard (mirroring layering_restore_test.go's token-absence family) asserting the deferred-tier vocabulary stays ABSENT from the shipped contract cores until 72 ships, with a discriminator control.

## Out of scope

- Re-introducing the tier mechanism — that rides 72 (fo-tier-delegation) when it un-defers, WITH its «fo.tier» source. czw's leak was in the right place, just premature.
