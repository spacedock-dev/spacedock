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

---

## SPIKE RESULTS (ideation, 2026-06-03) — the riskiest unknown, exercised first

The riskiest unknown was: **how does a non-interactive `claude -p` FO subprocess cleanly EXIT when the team is un-deletable — does ending the turn with no further work terminate the process, or is an explicit signal needed?** This is exactly what yy got wrong, so it was paid first.

**Exercise.** Real `claude -p` invocations (claude 2.1.161) with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` (the CI env), `--permission-mode bypassPermissions --output-format stream-json`. Each spike dispatched a real general-purpose member running `sleep 600` (a member that will NOT settle within the window), then varied what the FO did at teardown. Exit code from `timeout` distinguishes clean self-exit (0) from hang (124 = SIGTERM-on-timeout).

| Spike | FO behaviour at teardown | Member | Result |
|---|---|---|---|
| plain `-p`, no team | end turn (text only) | — | **exit 0** — process exits on `end_turn` |
| empty team | TeamCreate, no members, end turn | — | **exit 0** (2 turns: harness re-woke once, then exited) |
| active member, kept retrying | dispatched member, retried TeamDelete across turns until it SUCCEEDED | sleep 600 | **exit 0** — exited only AFTER a TeamDelete success (member finally settled) |
| **active member, STOPPED at cap** | dispatched member, one/two failed TeamDelete, then ended turn with NO further calls (modelling "hard-exit at cap") | sleep 600 | **HUNG → `timeout` SIGTERM (124)** — ending the turn did NOT exit; harness kept the process alive |
| Bash self-kill | `kill -TERM $PPID; kill -TERM $$` from a Bash tool | sleep 600 | exit **143** (unclean SIGTERM) — rejected by the live test's `exitCode != 0` gate; orphans the team |

**The FO model stated the mechanism verbatim** in the kept-retrying run's final result message:
> "The process did **not** exit while the member was active. Instead, **the harness re-invoked me in non-interactive mode with a hard requirement: I cannot return a response until the team is shut down.** … My first `TeamDelete` failed … then [team] is fully shut down and cleaned up."

**Pinned mechanism (load-bearing, changes the design):** In non-interactive `claude -p` teams mode the harness enforces a hard invariant — **the process will NOT exit while the team has active members.** `end_turn` does NOT terminate the subprocess when a member is still in the roster; the harness re-invokes the FO *specifically to force teardown*. The **only** clean exit (code 0, which the live `expectExit`/`exitCode==0` gate requires) is a **successful `TeamDelete`** (empty roster). A Bash self-kill exits non-zero and orphans the team, so it is not viable.

### What the spike INVALIDATES (escalation — captain-approved A+C is not achievable as written)

The captain's approved direction A — "hard-exit at the cap (best-effort; residual team cleanup happens at session/process death), NOT keep going" — **rests on a runtime mechanism that does not exist.** There is no process death to fall back on: the harness keeps `claude -p` alive *until the team is empty*, so "stop retrying and end the turn" produces exactly the cap.jsonl HANG (124), which is the same failure class as #275/#277. "Residual cleanup at process death" cannot happen because the process does not die while members remain.

Consequently the seed's AC-1 ("the FO completes a bounded teardown and the process EXITS") and AC-3 ("the contract mandates process EXIT at the cap, not keep retrying") are **mutually unsatisfiable on this runtime**: the only thing that makes the process EXIT is a TeamDelete SUCCESS, which is precisely "keep retrying until it succeeds" (yy), not "hard-exit at the cap." A is the opposite of what the runtime requires.

### What the spike VALIDATES about yy, and the real (narrower) bug

yy's "retry-to-success" is the **correct and only** exit mechanism — the kept-retrying spike proves a TeamDelete success is reachable and exits cleanly once the member settles. yy is right in kind. The REAL bug the live sonnet cycle exposed is narrower and is about TIMING, not the retry/exit dichotomy:

1. **Members can take longer to settle than the FAST cap covers.** The member is mid-tool (e.g. a long Bash); it has not processed the cooperative `shutdown_request` yet. yy's cap is ~5 fast attempts × ~2s settle ≈ 10s. If settle takes longer, the cap's FAST phase exhausts and — per yy's own step 4 — the FO "MUST keep settling and retrying," i.e. it keeps emitting tool calls, which never reaches an empty `end_turn`. That is fine for exit *eventually*, BUT:
2. **The live streamwatcher's `exitBudgetDefault = 60s` is too tight for a slow-settling member.** The new sonnet failure (seed's runs 26891717026 / 26891718562) showed 6 TeamDelete attempts and still timed out — yy DID keep retrying, but the member had not settled by the 60s `expectExit` budget, so `expectExit` tripped and killed it. The hang the seed attributes to "reads cap-exhaustion as keep-retrying and never exits" is really "the member had not settled within the 60s exit budget, so the still-correct retry loop was killed before it could reach success."

So the fix is **not** A (hard-exit-at-cap is impossible) — it is to make yy's retry-to-success reliably *reach* success within the live budget:
- **(C, kept and central) Widen `expectExit`'s budget** to tolerate a slow member settle (the retry loop is emitting TeamDelete attempts the whole time, so this is a no-progress-style tolerance, not an unbounded wait). The streamwatcher already fires `expectExit` on ANY clean process exit and does NOT require TeamDelete success — that half of C is already satisfied (see `TestExpectExitWaitsThenKills/exits_cleanly`); the missing half is budget.
- **(A', replacing A) Strengthen the settle, not abandon the retry.** Between attempts, after the FAST cap is exhausted, escalate to a LONGER settle (the cooperative shutdown was acknowledged; the member just needs more wall-time to leave a long tool). The cap should bound the FAST cadence then switch to longer-settle retries — never "stop and exit," which cannot exit.

This is a substantive divergence from the captain-approved A+C and is being escalated to team-lead/captain before the body is treated as final. The spike, the mechanism, and the reframed direction are recorded here regardless of the decision (the spike is the deliverable).

### Reframed acceptance criteria (PENDING captain confirmation of the A-reframe above)

**AC-1 — The sonnet live cycle PASSES.** Verified by: a live-e2e sonnet run where the FO retries teardown to a TeamDelete SUCCESS and the process EXITS code 0 within the (widened) exit budget — `expectExit` fires, no timeout. (Runtime-observable, by-construction-pending-live.)

**AC-2 — An offline fixture pins retry-to-success-then-exit within budget.** Verified by: a Go test (extending `TestSonnetTeamDeleteHangReplay` or a sibling) dripping a stream where TeamDelete fails N times then SUCCEEDS, with the `fakeProc` flipping to exited on the success line, asserting `expectExit` fires (returns code 0, no `stepTimeout`) — RED on a fixture that never reaches a TeamDelete success (the pure-hang shape), GREEN once the success+exit beat is present. This pins the exit-on-success path the never-exiting `fakeProc` of the current replay cannot.

**AC-3 — The streamwatcher's exit budget tolerates a slow member settle.** Verified by: a unit test on `streamWatcher.expectExit` asserting it does NOT trip while the proc stays alive up to the new budget and DOES fire on exit just after — pinning the concrete `exitBudgetDefault` value, RED if the budget is reverted to a value too tight for the observed settle latency. (Note: this is a real CODE seam — `internal/ensigncycle/streamwatch.go:40` — unlike the prose-only teardown clause.)

**AC-4 — The contract mandates retry-to-success with a longer-settle escalation past the FAST cap, and forbids the impossible hard-exit-at-cap.** Verified by: extending `TestTerminalTeardownRetriesToSuccess` so the required set keeps "retry … to success" and ADDS a longer-settle-past-the-cap clause, and the NEGATING set ADDS the now-disproven "stop retrying and exit the process at the cap" / "hard-exit at the cap" phrasings (which the spike proved cannot exit). An inverted edit that re-introduces hard-exit-at-cap OR drops retry-to-success reds it. (Prose-only ceiling — the behavioural oracle is AC-1.)

### Reframed test plan

- Spike: DONE (above) — pinned the harness exit invariant; invalidated A, validated the C-plus-budget direction.
- AC-2/AC-3/AC-4: offline Go tests (`internal/ensigncycle/streamwatch_*_test.go` + `skills/integration/terminal_teardown_retry_test.go`). AC-3 is a genuine code gate on the budget constant; AC-2 extends the replay with a success+exit beat; AC-4 extends the directional-mandate lint.
- AC-1: live sonnet run (by-construction-pending-live). High-stakes (FO contract + CI machinery) → detached adversarial audit before merge, per the seed.
- Estimated cost: AC-2/3/4 are sub-second offline Go tests; AC-1 is one live-e2e sonnet cycle (~minutes, CI-only / Linux-bound).

> NOTE: this reframe diverges from the captain-approved A+C because the spike proved A's exit-at-cap mechanism does not exist on this runtime. Escalated to team-lead 2026-06-03. If the captain insists on a literal hard-exit, the only runtime-available "exit without TeamDelete success" is a non-zero Bash self-kill (143), which the live `exitCode==0` gate rejects and which orphans the team — so that path needs the live gate relaxed too, and is strictly worse. Recommendation: adopt the C-plus-budget reframe.

## Stage Report: ideation

- DONE: SPIKE FIRST (the riskiest unknown): pin the REAL mechanism by which a non-interactive `claude -p` FO subprocess cleanly EXITS when the team is un-deletable
  Ran 5 real `claude -p` teams-mode spikes (table in SPIKE RESULTS). Pinned invariant: the harness will NOT exit `claude -p` while the team has active members; `end_turn` does not terminate; ONLY a successful TeamDelete (empty roster) exits clean (code 0). cap.jsonl hung→124; kept-retrying→0 after TeamDelete success.
- FAILED: Land the A+C design with each AC backed by a CODE or TEST gate (the contract hard-exit-at-cap clause A + streamwatcher best-effort-teardown-then-exit C)
  The spike INVALIDATED contract A: "hard-exit at the cap, residual cleanup at process death" is impossible — there is no process death while members remain; ending the turn at the cap reproduces the #275/#277 hang (124). A and the seed's exit-EXITS AC-1 are mutually unsatisfiable on this runtime. Reframed to a C-plus-budget direction (retry-to-success is the only exit; widen `expectExit` budget + longer-settle-past-cap) and ESCALATED to team-lead, since this diverges from the captain-approved A+C.
- DONE: Test plan that separates offline ACs (AC-2/AC-3/AC-4 Go fixtures + oracle) from runtime-observable AC-1 (live sonnet), marking AC-1 by-construction-pending-live and naming the detached-audit trigger
  Reframed ACs/test plan recorded: AC-3 is now a genuine CODE gate on `streamwatch.go:40` exitBudgetDefault (not prose-only); AC-2 extends the replay with a TeamDelete-success+exit beat; AC-4 extends the directional lint to forbid the disproven hard-exit-at-cap. AC-1 marked pending-live; high-stakes detached audit named.

### Summary
Spiked the riskiest unknown first and it overturned the captain-approved A+C: real `claude -p` teams-mode runs prove the harness keeps the subprocess alive until the team is empty, so "hard-exit at the cap" cannot exit (it hangs exactly like #275/#277) and only a successful TeamDelete exits cleanly. The real, narrower bug is timing — a slow-settling member outlasts the 60s `expectExit` budget while yy's retry-to-success loop is still (correctly) running. Reframed the design to keep retry-to-success and widen the exit budget + escalate to a longer settle past the fast cap, with AC-3 now a concrete code gate on the budget constant. Because this diverges from the approved A+C, it is flagged in-body and escalated to team-lead before the body is final.
