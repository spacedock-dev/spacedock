---
id: atwf2w6p68t9q1mda790dcfc
title: Non-interactive FO must hard-exit after a bounded teardown attempt — yy's retry-to-success loops forever in claude -p and still flakes sonnet live-CI
status: ideation
source: live AC-1 confirmation of yy (PR #282) FAILED — both n3 #275 and 2a #277 sonnet cycles failed with yy's fix in the checkout; the sonnet FO looped 6 TeamDelete attempts in the settle-then-retry loop and never exited claude -p. Captain chose A+C (2026-06-03).
score: "0.42"
worktree:
started: 2026-06-03T15:28:23Z
completed:
verdict:
issue:
---

yy (PR #282) fixed the original "end the turn on the first `TeamDelete` failure" hang by mandating retry-to-success + an inter-attempt settle + a non-interactive cap-exhaustion exit. The live sonnet confirmation **failed**: the FO now does the settle-then-retry loop (6 `TeamDelete` attempts observed), but in non-interactive `claude -p` it reads the cap-exhaustion clause as *"keep retrying"* and **never cleanly exits the process** → the streamwatcher's `expectExit` times out → FAIL. yy traded a no-retry hang for a retry-loop hang.

## Problem

This is the audit's cycle-1 **polish finding #3** manifesting (flagged then as "load-bearing for AC-1"). Two coupled root causes:

1. **Contract: the non-interactive cap-exhaustion EXIT is ambiguous.** "surface to captain and end the turn" does NOT terminate the `claude -p` process when the team can't be deleted (members stuck). The FO keeps thinking/acting → never finishes → never exits. From the failure JSONL, the FO's own thinking: *"the non-interactive exit obligation requires me to keep retrying."*
2. **Test/streamwatcher: requires `TeamDelete` success before the FO exits.** The live cycle's `expectExit` only fires on a clean process exit, which never comes if teardown can't complete. The streamwatcher over-constrains: it should accept a **best-effort-teardown-then-exit**.

**Evidence:** live runs `26891717026` (n3) + `26891718562` (2a), artifact `runtime-live-e2e-claude-live-sonnet` — the FO session JSONL shows 6 `TeamDelete` attempts + the "keep retrying" thinking, no process exit. opus passes (it happens to exit); sonnet loops.

yy's fix is kept as-is — retry-to-success genuinely helps a *real* session (members do terminate, so the retry recovers). This entity closes the live-CI / non-interactive exit path it does not cover.

## Proposed approach (A + C)

1. **(A — contract) Hard-exit at the cap.** In non-interactive mode, terminal teardown is: cooperative shutdown + `TeamDelete` up to a small cap WITH the inter-attempt settle, and if still failing at the cap, the FO **EXITS the process** (best-effort; residual team cleanup happens at session/process death) — NOT "surface and keep going." Make the cap-exhaustion EXIT unambiguous: stop retrying, terminate. An inverted "keep retrying past the cap" must be a contract violation the oracle catches.
2. **(C — streamwatcher/test) Accept best-effort-teardown-then-exit.** The live cycle / streamwatcher's `expectExit` should fire when the FO exits after the bounded teardown attempt, not require `TeamDelete` success. Widen the `expectExit` budget for the settle latency if needed. Update the yy AC-2 replay fixture/oracle to reflect the bounded-exit contract.

**Riskiest unknown — exercise first:** how does the FO subprocess cleanly EXIT `claude -p` from within the contract when there's an un-deletable team? Does ending the turn with no further work actually exit, or is an explicit signal needed? This is exactly what yy's fix got wrong; the spike must pin the real exit mechanism before designing the prose.

## Out of scope

- Reverting yy (kept as-is — it helps real sessions; this is the non-interactive exit complement).
- The opus path (already exits cleanly).

## Acceptance criteria

**AC-1 — The sonnet live cycle PASSES (the real behavioural oracle).**
Verified by: a live-e2e sonnet run (on n3 #275 / 2a #277 or a dedicated run) where the FO completes a bounded teardown and the process EXITS — `expectExit` fires, no timeout. This is the oracle yy's fix failed; this time it must pass.

**AC-2 — An offline fixture pins bounded-teardown-then-exit (no infinite loop).**
Verified by: a Go test dripping a stream where `TeamDelete` never succeeds, asserting the FO makes ≤cap attempts with settle and then the modelled process exits (the streamwatcher's `expectExit` fires), RED on the pre-fix loop-forever behaviour.

**AC-3 — The contract unambiguously prescribes the hard-exit at the cap.**
Verified by: an oracle requiring the cap-exhaustion clause to mandate process EXIT (not "keep retrying"); an inverted "retry past the cap forever" edit reds it.

**AC-4 — The streamwatcher accepts best-effort-teardown-then-exit.**
Verified by: a test asserting `expectExit` fires on a best-effort-teardown-then-exit stream (no `TeamDelete` success required), and that the budget tolerates the settle latency.

## Test plan

- Spike the claude -p exit mechanism FIRST (the riskiest unknown).
- Go tests for AC-2/AC-3/AC-4 (offline fixtures + oracle).
- The live sonnet run for AC-1 (by-construction-pending-live, but it MUST pass this time — that is the whole point).
- High-stakes (FO contract + CI machinery) → detached adversarial audit before merge.

## Notes

- This unblocks n3 #275 + 2a #277 (still blocked on the sonnet flake) once it lands and the live sonnet passes.
- The live-confirmation discipline (captain's plan A) caught this post-merge — offline + 3 audit cycles all passed, but the live cycle revealed the fix changed the hang mode rather than closing it. Strongest possible argument for the live oracle on runtime-observable ACs.
- Pairs with `yy` (sonnet-live-ci-flake, shipped #282) as its non-interactive-exit completion.
