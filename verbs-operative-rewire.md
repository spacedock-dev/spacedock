---
title: Make the shipped state/merge verbs the operative contract path + bind «fn» bodies to their verbs (oracle)
source: '0221-layered-fo rework (2026-06-19): validated findings 1+2 — the contract does not operatively USE the shipped verbs. `### Split-Root State Sync` names an abstract "status tool" not `spacedock state commit`, and `«state.commit»`s effect restates the hand git sequence; `«merge.guard»`s prose claims it "invoke[s] the registered merge hook" / runs "as one call" / "default-merge[s]", all FALSE of the shipped re-entrant partial envelope (merge.go doc-comment: "It does NOT invoke the merge hook"). Bundles the contract rewire (A) with a routing oracle (B) written test-first.'
status: ideation
score: 0.75
sprint: 0221-layered-fo
group: foundation
id: 4asxw7kxvdzdtf87w9rjkxwx
started: 2026-06-19T22:53:22Z
---

Complete the vertical slice `rgq`/`mz`/`czw` shipped empty: the verbs exist and route, but the contract instructs the hand sequence and mis-describes the merge verb. Make the verb the operative path AND guard it mechanically. (Per the captain's call, A and B are ONE bundled slice, not two layer-tasks.)

## Proposed approach (firm up in ideation)
**A — rewire the contract prose (`skills/first-officer/references/`, dispatched worker — scaffolding guardrail, not an FO edit):**
1. `### Split-Root State Sync` "Preferred" bullet → name `spacedock state commit <slug>` as the path (keep hand `git -C add/commit` as the genuine degraded / no-origin fallback only).
2. `«state.commit»` effect bullet → delegate to the verb, not restate the hand git sequence.
3. `«merge.guard»` body (`fo-merge-core.md`) → describe the verb ACCURATELY as a re-entrant partial envelope (arm: set mod-block + emit "invoke the hook, then re-run"; the FO invokes the hook; finalize: clear mod-block + terminalize + archive). Drop the false "invoke[s] the registered merge hook", "default-merge", and "as one call" claims. Worktree-removal (step 9) + worker-teardown (step 10) stay the FO's, correctly.

**B — the «fn»→shipped-verb routing oracle (`internal/contractlint`), written TEST-FIRST:**
For each `→ shipped: spacedock <verb>` in an «fn» body, invoke the verb (argv) and assert it routes. RED on today's overstated `«merge.guard»` (or a mutation pointing an «fn» at a non-existent verb); GREEN after A. This REPLACES the deferred-tier absence-grep — a behavioral routing oracle, not a literal-string prose blocklist.

## Seed acceptance criteria (firm up in ideation)
- AC-1: An FO following `### Split-Root State Sync` / `«state.commit»` is directed to `spacedock state commit` as the primary path (hand-git is the named fallback only). Verified by the rewritten prose + the routing oracle GREEN on the `state commit` binding.
- AC-2: `«merge.guard»`s description matches the shipped verb's actual behavior (re-entrant arm→re-run→finalize; does NOT itself invoke the hook). Verified by the routing oracle GREEN + a diff removing the three false claims.
- AC-3: The routing oracle is non-vacuous — RED against the pre-rewrite overstated prose (or a non-existent-verb mutation), GREEN after. Verified by the RED-then-GREEN control, replacing `deferred_tier_absence_test.go`.

## Live-usage proof
The end-value proof that a live FO actually CALLS the verbs as its path rides `kt` (haiku-drive-validation / `kt-full`), sequenced after the premise spike. This slice ships the rewire + the mechanical oracle; the live drive is `kt`'s.
