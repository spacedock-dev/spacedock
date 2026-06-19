---
id: ga1e7heckpazwbmtxzrg2kms
title: Strip deferred-72 tier vocabulary leaked into czw's «gate.assemble-verdict» body (pre-cut ship-blocker)
status: validation
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
mod-block: merge:pr-merge
pr: "#403"
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

## Stage Report: validation

- DONE: grep the shipped FO contract cores for level-2-only / level-3 / judge / «fo.tier» / L3 → ZERO.
  `grep -rniE "level-2-only|level-3|«fo\.tier»|\bL3\b|\bjudge\b" skills/first-officer/references/*.md` exits 1 (zero matches). The named ship-blocker tokens are fully stripped.
- FAILED: Read the restructured «gate.assemble-verdict» body and confirm it reads mechanism-neutrally with NO dangling reference to an undefined decider/judge/level.
  Body bullets (lines 127, 130) are clean and mechanism-neutral. BUT the strip left line 124's heading — `## «gate.assemble-verdict»(slug, stage): assemble the gate review, route the verdict to the decider` — and line 116's intro `route the verdict:`. "the decider" appears EXACTLY ONCE in all of skills/ (line 124), entered in czw's restructure (c25fee26), and is DEFINED NOWHERE. It is the residual tail of the tier-delegation framing (route to a level-2-FO-or-level-3-judge decider); the body it summarizes now says "The FO renders its own `Recommend` line" — there is no decider to route to, the FO IS the decider. This is precisely the dangling-decider reference the checklist names. Base intent (2c9609f8) had no "decider" routing concept at all.
- DONE: Absence-guard is load-bearing, not a tautology; reproduce a mutation.
  BOTH mutations reproduced. (1) Planted `level-3 judge` into the real first-officer-shared-core.md → `TestFOContractCoresHaveNoDeferredTierToken` REDS (flagged the planted line); reverted → green. (2) Stubbed the shared `lineLeaksDeferredTierToken` scanner to never-flag → `TestDeferredTierTokenScannerDiscriminates` REDS on all three genuine leak shapes WHILE the absence check passes VACUOUSLY — proving the discriminator is the non-vacuity guard. `skillsRoot` resolves the real `skills/` dir, so the scan is over shipped files, not a fixture.
- DONE: Surgical scope — diff vs main is ONLY the gate body prose + one new test file; «fn» notation, routing oracle, layering_restore, boot-resident closure + ceremony, six behavioral clusters UNCHANGED.
  `git diff main --name-only` = exactly 2 files. Contract-core change is exactly the 2 body lines (124 heading NOT in diff — preserved as-is). TestProseFunctionNotationBindsToRouting/RoutingGuard/RoutingOracleDiscriminates, TestBootResidentDeferredLoadPointsResolve, TestHostNeutralCoresResolveAndCarryCeremony, layering_restore + PRView suites all PASS; their source files byte-identical to main.
- DONE: Full `go test ./...` green; go vet clean.
  All packages `ok`, no FAIL. `go vet ./...` exit 0. gofmt: this task's new test file is clean; `gofmt -l` flags two PRE-EXISTING files (internal/cli/prose_function_routing_test.go, internal/status/section_read.go) that are byte-identical to main and untouched by this branch — environmental, not a regression of this change.

### Summary

REJECTED. The named-token strip is mechanically complete (grep → ZERO) and the contractlint absence-guard is proven load-bearing by both required mutations. But the strip is INCOMPLETE on coherence: it removed the two body sentences that routed the verdict to a tier "decider" yet left the `«gate.assemble-verdict»` heading (line 124) reading "route the verdict to **the decider**" and the intro (line 116) "route the verdict" — naming an undefined "decider" that the now-mechanism-neutral body contradicts. This is the exact dangling-decider reference the checklist's second item rejects, and it diverges from base intent (2c9609f8 had no decider-routing concept). Bounce-back (surgical, one file): rephrase line 124's heading to drop "to the decider" — e.g. `assemble the gate review and render the verdict` — and align line 116 so the heading/intro promise matches the body (the FO renders its own `Recommend` line); no code change, no test change, the absence-guard already covers it.

## Stage Report: implementation (cycle 2)

- DONE: Drop the dangling "decider" — rephrase line 124's heading to drop "to the decider" and align line 116's intro so heading/intro match the body (the FO renders its own `Recommend` line).
  Both edits made in first-officer-shared-core.md (code commit 5127caa1): line 124 heading → "assemble the gate review and render the verdict"; line 116 intro → "render the verdict:". Diff this cycle is exactly those two lines. `grep -rn 'decider' skills/` now returns ZERO. czw's «fn» notation, the `→ prose` routing oracle, and the six behavioral clusters remain untouched.
- DONE: Re-run `go test ./...` green (no test change expected); commit, clean status, HEAD sha.
  `go test ./...` all packages `ok` (no FAIL); no test file changed (prose-only — "decider" is a one-off framing word, not a systematic token, so the absence-guard correctly does not cover it). Committed 5127caa1 on spacedock-ensign/strip-deferred-tier-vocabulary, working tree clean.

### Summary

Closed the last dangling reference the validation cycle flagged: "the decider" — the level-3-judge tier framing under a different word, defined nowhere, contradicting the now-mechanism-neutral body. Aligned the `«gate.assemble-verdict»` heading and the gated-stage intro to the body (assemble the gate review and render the verdict). `grep -rn decider skills/` → ZERO; full suite green; prose-only, no test change. Code HEAD: 5127caa1.
