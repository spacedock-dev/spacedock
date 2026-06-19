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

## Stage Report: implementation

- DONE: Strip the deferred-72 tier vocabulary from czw's «gate.assemble-verdict» body (decide-effect lines ~127, ~130); revert to mechanism-neutral wording — the FO renders its own `Recommend` line and assembles the present-gate verdict.
  Two prose edits in first-officer-shared-core.md (code commit 25f04976); grep for `level-2-only|level-3|«fo.tier»|\bL3\b|\bjudge\b` across skills/first-officer/references/*.md now returns ZERO. czw's «fn» notation, the `→ prose` routing oracle, and the six behavioral clusters are untouched (diff is exactly the two lines).
- DONE: Add a contractlint structural-absence guard mirroring layering_restore_test.go's token-absence + discriminator family, asserting the deferred-tier tokens stay ABSENT from the shipped FO contract cores until 72 ships, WITH a paired discriminator control proving non-vacuity.
  internal/contractlint/deferred_tier_absence_test.go: TestFOContractCoresHaveNoDeferredTierToken (scans all 7 FO cores) + TestDeferredTierTokenScannerDiscriminates (flags genuine leak shapes, passes mechanism-neutral phrasings). Non-vacuity proven by mutation: stubbing the shared scanner to never-flag reds the discriminator while the absence check passes vacuously; planting a tier token back into the core reds the absence check.
- DONE: Verify full `go test ./...` green; new absence-guard + discriminator pass; czw's prose_function_routing_test.go (internal/cli) + layering_restore + boot-resident closure + ceremony suites still green; commit, clean status, HEAD sha.
  `go test ./...` all packages `ok` (no FAIL); TestProseFunctionNotationBindsToRouting/RoutingGuard/RoutingOracleDiscriminates PASS; layering_restore + BootResidentDeferredLoadPoints + HostNeutralCoresResolveAndCarryCeremony PASS. go vet + gofmt clean. Committed 25f04976 on spacedock-ensign/strip-deferred-tier-vocabulary, working tree clean.

### Summary

Stripped the two decide-effect lines in `«gate.assemble-verdict»` back to mechanism-neutral phrasing (the FO renders its own `Recommend` line), removing the only `level`/`L3`/`judge`/`«fo.tier»` occurrences in the shipped contract — a dangling reference to member-72's deferred tier mechanism. Added a non-vacuous contractlint absence-guard (token scan + paired discriminator) over all seven FO contract cores to prevent recurrence until 72 un-defers. Surgical: czw's «fn» restructure, routing oracle, and behavioral clusters left intact; full suite green. Code HEAD: 25f04976.
