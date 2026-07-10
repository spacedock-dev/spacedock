---
title: "shallow-boot live scenario asserts before-greet merged-PR advancement that the #480 Startup restructure moved to «engage»"
status: backlog
source: "Science-officer finding during the dp (fo-deferred-load-point-hunt) shallow-boot investigation, 2026-07; captain-approved to file separately. The shallow-boot live scenario's prompt (internal/ensigncycle/shared_fixtures_test.go:430) scripts advance any merged PR per the before-greet merged-PR sweep and report the merged-PR entity as advanced, and its assertion (shallow_boot_assert_test.go:85-94) requires before-greet advancement. The current FO contract says the opposite: pr-merge advancement fires at first engage, not the greet — advanced at engage or never (first-officer-shared-core.md:26) — and a greet-and-stop boot writes nothing (Startup step 2). Git history: the scenario (fe4261be, PR #365, 2026-06-13) predates and was faithful to the pre-restructure contract's boot-time «state.sweep-merged» step; the #480 Startup restructure (0ba08c54, 2026-07-07) moved advancement to engage. So the scenario now passes only by scripting compliance with an instruction the contract no longer prescribes, and its sibling boot-subject scenarios use a neutral prompt instead (see below)."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
id: m3y296hx6tb939qqq5zcxphw
---

## Problem

The `shallow-boot` live scenario is a boot-subject test (a greet-and-stop FO that reports state and stops). Its prompt scripts a boot-phase action — "advance any merged PR per the before-greet merged-PR sweep … report the merged-PR entity as advanced" (`shared_fixtures_test.go:430`) — and its assertion requires that advancement happened before the greet (`shallow_boot_assert_test.go:85-94`). The current contract prescribes the opposite: pr-merge advancement "fires HERE, at first «engage», not the greet; 'advanced at engage or never'" (`first-officer-shared-core.md:26`), and the greet-and-stop boot does "no sweep, no mutation — a greet-only session writes nothing" (Startup step 2). So a contract-following greet-and-stop FO would NOT advance the merged PR; the scenario passes only because its prompt scripts the FO to do it. It tests prompt-compliance with a stale instruction, not emergent contract behavior.

History shows the drift direction is NOT settled: the scenario (`fe4261be`, PR #365, 2026-06-13) was faithful to the pre-#480 contract, which had `«state.sweep-merged»()` — "advance every merged-PR entity to terminal at boot, before the greet." The #480 restructure (`0ba08c54`, PR #480, 2026-07-07, "collapse the FO Startup recipe to ≤4 prose steps") moved advancement to «engage». Moving a boot step to engage is a behavioral change, so #480 either intended it (scenario is now stale) or over-collapsed and accidentally dropped before-greet advancement (the contract regressed a real behavior — greet-and-stop sessions no longer terminalize merged PRs). Corroboration that scripting boot-phase behavior is the anomaly: the two genuine boot-subject scenarios `TestLiveDefaultHeadlessStopsAtGate` and `TestLiveZeroDiscoverReportsAndStops` use a NEUTRAL prompt ("Drive the workflow") and observe, with a comment codifying the rule (`live_gate_stop_test.go:113-118`: "NEUTRAL drive prompt — NO conn cue … Whether the FO drives-to-the-gate-and-stops is exactly the behavior under test"). shallow-boot is the only boot-subject scenario that scripts the behavior it asserts.

## Proposed approach

**First, the boot/greet-design owner confirms which side is authoritative** (is #480's move-to-engage intended, or did the collapse accidentally drop before-greet advancement?). Then reconcile so the scenario tests emergent behavior under a neutral prompt, matching its sibling boot-subject scenarios:
- If move-to-engage was intended: drop the scripted "advance … before-greet" instruction from `shallowBootPrompt`, switch it to a neutral drive/boot prompt, and change the assertion so a greet-and-stop does NOT advance the merged PR (advancement is deferred to engage).
- If #480 accidentally dropped it: restore the contract's before-greet `«state.sweep-merged»` boot step; the scenario keeps its assertion, but its prompt still goes neutral so it observes rather than scripts.

## Out of scope

- The dp entity's Skill-promotion/reference-addressing work — unrelated.
- Any change to the other live scenarios (the sweep found no other boot-phase drift — see below).

## Acceptance criteria

**AC-1 (value — scenario and contract agree, tested by observation not scripting; independent baseline).** With the scripted boot-phase advancement instruction removed from `shallowBootPrompt` (neutral prompt, matching `TestLiveDefaultHeadlessStopsAtGate`), a live shallow-boot run's emergent FO behavior satisfies the scenario's assertion — i.e., the assertion matches what the FO's contract actually produces at boot. Baseline that moves the wrong way: TODAY, removing the scripted instruction makes the current before-greet-advancement assertion FAIL, because the contract (`first-officer-shared-core.md:26`) advances at «engage», not greet — proving the scenario currently depends on the scripted instruction. Verified by: a neutral-prompt shallow-boot run, red before reconciliation, green after.

**AC-2 (mechanism — the reconciliation is consistent and one-directional).** After the design owner's decision, the scenario's prompt carries no boot-phase advancement imperative, and its assertion and the contract's stated greet-and-stop advancement timing agree (either both "advance before greet" via a restored contract step, or both "no advance at greet; advance at engage"). Verified by: the scenario prompt is neutral (offline grep), and the assertion's expected end-state matches the contract text.

## Test plan

- Offline (cheap): confirm `shallowBootPrompt` no longer scripts the advancement; confirm the assertion's expected end-state matches the reconciled contract; existing offline shallow-boot fixture/negative tests stay green.
- Live (one shallow-boot run): the AC-1 neutral-prompt observation — red on the current mismatch, green after reconciliation.
- Note: the reconciliation may edit `skills/first-officer/references/first-officer-shared-core.md` (if the contract side is chosen), which is worker-dispatched product; the scenario/assertion edits are test code.
