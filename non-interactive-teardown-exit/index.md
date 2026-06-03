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

**Hard proof (raw stream, not self-report).** The captain challenged the FO's self-report, so the critical case was re-run with the raw stream preserved at `/tmp/at-spike/cap-confirm.jsonl` (sha `93dc87c4b1ce56b4a1692003c6e91e38d8c870e4`, 25 lines). The FO emitted a genuine `end_turn` (final text, no tool calls) and the harness injected a fresh `system/init` turn — three times — and the process still never exited (`EXIT=124` at the 180s `timeout`). The `L17 result(end_turn) → L18 system/init` transition is the hard evidence that the no-exit is a HARNESS INVARIANT, not a model choice. (The earlier self-report quote — "I cannot return a response until the team is shut down" — is now merely corroborating colour.)

**Root-cause confirmation from the REAL CI artifact (not a synthetic spike).** The streamwatcher's flat-60s exit kill was triggered because the roster never empties: the ensign member APPROVED its shutdown and its session ENDED (14:37:30), but Claude Code never cleared it from `members[]`, so the FO saw `active member(s)` for ~55s against a DEAD member (still active at 14:38:25+) — upstream **#38116 / #57681**. This is why retry-to-success is unreachable and is the load-bearing fact behind the captain-approved direction below.

**Pinned mechanism (load-bearing).** Two facts, both empirically settled:

1. **The `claude -p` harness will NOT let the subprocess self-exit while the team's `members[]` is populated.** A genuine `end_turn` (final text, zero tool calls) does NOT terminate; the harness injects a fresh `system/init` turn and re-invokes the FO. Proven on the raw stream of `/tmp/at-spike/cap-confirm.jsonl` (sha `93dc87c4…`): the FO ended its turn three times (`stop_reason: end_turn`, each followed by a harness `system/init`), narrated that it was deliberately making no further tool calls, and the process still never exited — `timeout` killed it at 180s (`EXIT=124`). Reproduced on BOTH sonnet and opus. This is a harness invariant, not a model choice.

2. **The roster never empties, because the member is dead-but-listed.** From the REAL failed CI artifact (the ensign member's own session jsonl, runs `26891717026` / `26891718562`): the dispatched ensign finished its work, RECEIVED the FO's `shutdown_request`, APPROVED it ("my work is complete… I'll approve"), and its session TERMINATED at **14:37:30**. But the Claude Code harness **never removed it from the team `members[]`**. The FO kept getting `Cannot cleanup team with 1 active member(s)` for ~55s AFTER the member was already dead (FO still sees it active at **14:38:25+**), until the flat 60s `expectExit` killed it. This is upstream Claude Code bug **#38116 / #57681** (an approved-shutdown member is never cleared from `members[]`).

**What these jointly INVALIDATE:**
- yy's "retry `TeamDelete` to success" is **UNREACHABLE**: success never comes, because the dead member is never removed from `members[]`. The retry loop runs forever (bounded only by the 60s exit budget that then kills it).
- contract A's "hard-exit at the cap; residual cleanup at process death" is **IMPOSSIBLE**: the harness will not let `claude -p` self-exit while `members[]` is non-empty (fact 1), so there is no process death the FO can cause. (My earlier no-progress-reset reframe is also dropped — it only widens the wait for a success that, per fact 2, never arrives.)

### Approved direction (captain, 2026-06-03): grade correct best-effort teardown; the LAUNCHER owns the exit

The exit cannot come from the FO (impossible) or from a `TeamDelete` success (unreachable). So we stop demanding either. Two coordinated changes:

1. **(Test / streamwatcher) Grade the terminal teardown on CORRECT BOUNDED BEST-EFFORT BEHAVIOR, then the launcher kills the subprocess and the cycle PASSES.** The live cycle replaces its final `expectExit → exitCode==0` step (live_test.go:157-163) with: observe that the FO did the right sequence — sent `shutdown_request` → member approved → FO made BOUNDED `TeamDelete` attempts — then STOP the subprocess via the existing `defer poller.kill()` (live_test.go:125-133) and PASS. It MUST NOT require a clean self-exit, which the upstream bug makes impossible. This is a STRONGER assertion than `exitCode==0`: it checks the FO performed the correct teardown sequence, not that the process happened to exit. No-progress-reset is UNNECESSARY (we observe a FINITE bounded sequence, not an open-ended success-wait) — drop it; the ≤60s guard stays green because the new step still uses the ≤60s no-progress budget for each beat.
2. **(Contract — `first-officer-shared-core.md` + the Claude runtime adapter) Non-interactive terminal teardown is BOUNDED best-effort, then STOP and emit a terminal status; the PROCESS EXIT is the LAUNCHER's responsibility, NOT the FO's.** The FO: cooperative `shutdown_request` to the cohort → a BOUNDED set of `TeamDelete` attempts with the inter-attempt settle → then STOP retrying and emit a terminal status (it cannot delete a team the harness won't let it empty). The process death comes from the launcher (the live test's `kill()`, or a real automation's timeout) — the FO cannot cause it. After the bound, the harness will re-invoke the FO (proven); the FO just HOLDS (text only, no tool calls) until the launcher kills it. This **rehabilitates A's SPIRIT** — best-effort then process death — except the death is launcher-owned. REMOVE/correct yy's "keep retrying `TeamDelete` to success in non-interactive mode" (unreachable → the hang) and the disproven "hard-exit at the cap / residual cleanup at process death" (the FO cannot cause the death).

### Acceptance criteria (captain-approved direction)

**AC-1 — The sonnet live cycle PASSES via correct bounded best-effort teardown + launcher kill (runtime-observable oracle).**
Verified by: a live-e2e sonnet run where the FO performs the correct bounded teardown sequence (shutdown_request → member approval → bounded `TeamDelete` attempts) and the harness/launcher kills the subprocess — NO clean self-exit required, NO `exitCode==0` gate. `expectExit`-style flat-cap-then-fail is gone. By-construction-pending-live; MUST pass on a real sonnet run (the oracle yy's fix failed).

**AC-2 — An offline fixture reproduces the dead-but-listed shape and the watcher grades it as correct best-effort teardown → PASS.**
Verified by: a Go test (sibling to / extending `TestSonnetTeamDeleteHangReplay`) dripping a stream where the member APPROVES shutdown and its session ENDS, yet `TeamDelete` keeps failing `active member(s)`; the new grading step asserts the watcher OBSERVED the correct beats (shutdown_request sent + member approval + ≥1 bounded `TeamDelete` attempt) and PASSES (then the modelled `kill()` stops the proc). RED on a fixture where the FO did NOT do the correct teardown — never sent `shutdown_request`, OR never attempted `TeamDelete` — so the grade is behavioral, not "the proc was killed regardless." (The existing replay's `fakeProc`-never-exits tautology is replaced by a behavioral grade.)

**AC-3 — The contract mandates bounded best-effort teardown + launcher-owned exit, and FORBIDS the unreachable/impossible framings.**
Verified by: extending `TestTerminalTeardownRetriesToSuccess` so the required set asserts "bounded best-effort teardown" + a "terminal status, then hold" + "the launcher/process owner ends the subprocess" clause, and the NEGATING set ADDS "retry … to success" (now unreachable) and "hard-exit at the cap" / "residual cleanup at process death" / "the FO exits the process" (now impossible). An inverted edit re-introducing retry-to-success OR FO-self-exit reds it. (Prose-only ceiling — the behavioral oracle is AC-1.)

### Test plan

- Spike: DONE — the riskiest unknowns are empirically settled by (a) the raw `cap-confirm.jsonl` harness-reinvocation reproduction (sonnet+opus, `EXIT=124`) and (b) the real CI artifact's dead-but-listed timeline (member dies 14:37:30 / FO still sees it active 14:38:25+, upstream #38116/#57681). NO further live spike needed for the mechanism.
- AC-2: offline Go test in `internal/ensigncycle/` (the dead-but-listed fixture + the new behavioral grade); sub-second, no model spend.
- AC-3: contract oracle in `skills/integration/terminal_teardown_retry_test.go` (extend the directional-mandate lint); sub-second.
- AC-1: one live-e2e sonnet cycle (~minutes, CI-only / Linux-bound), by-construction-pending-live.
- **≤60s guard stays green.** The new grading step uses the existing ≤60s no-progress budget per beat and the `defer poller.kill()` for exit, so no timeout literal grows; `TestNoTimeoutLiteralExceeds60s` (live_budget_test.go) is unaffected. Confirm it stays green at validation.
- High-stakes (FO contract + CI machinery) → detached adversarial audit before merge.

> OPEN QUESTION (flagged, non-blocking; default assumption = live CI is the only consumer): are there real headless/cron `spacedock … -p` FO runs that gate on a CLEAN exit code? If so, they need a launcher-side timeout (a wrapper that bounds the run and kills it), NOT a contract change and NOT any reach into Claude's private team registry (captain rejected the config-rewrite host-hack). Default: only the live-e2e CI consumes the FO `-p` exit, and it already owns a `kill()`, so no extra launcher work is needed.

## Stage Report: ideation

- DONE: SPIKE FIRST (the riskiest unknown): pin the REAL mechanism by which a non-interactive `claude -p` FO subprocess cleanly EXITS when the team is un-deletable
  Ran 5 real `claude -p` teams-mode spikes (table in SPIKE RESULTS). Pinned invariant: the harness will NOT exit `claude -p` while the team has active members; `end_turn` does not terminate; ONLY a successful TeamDelete (empty roster) exits clean (code 0). cap.jsonl hung→124; kept-retrying→0 after TeamDelete success.
- FAILED: Land the A+C design with each AC backed by a CODE or TEST gate (the contract hard-exit-at-cap clause A + streamwatcher best-effort-teardown-then-exit C)
  The spike INVALIDATED contract A: "hard-exit at the cap, residual cleanup at process death" is impossible — there is no process death while members remain; ending the turn at the cap reproduces the #275/#277 hang (124). A and the seed's exit-EXITS AC-1 are mutually unsatisfiable on this runtime. Reframed to a C-plus-budget direction (retry-to-success is the only exit; widen `expectExit` budget + longer-settle-past-cap) and ESCALATED to team-lead, since this diverges from the captain-approved A+C.
- DONE: Test plan that separates offline ACs (AC-2/AC-3/AC-4 Go fixtures + oracle) from runtime-observable AC-1 (live sonnet), marking AC-1 by-construction-pending-live and naming the detached-audit trigger
  Reframed ACs/test plan recorded: AC-3 is now a genuine CODE gate on `streamwatch.go:40` exitBudgetDefault (not prose-only); AC-2 extends the replay with a TeamDelete-success+exit beat; AC-4 extends the directional lint to forbid the disproven hard-exit-at-cap. AC-1 marked pending-live; high-stakes detached audit named.

### Summary
Spiked the riskiest unknown first and it overturned the captain-approved A+C: real `claude -p` teams-mode runs prove the harness keeps the subprocess alive until the team is empty, so "hard-exit at the cap" cannot exit (it hangs exactly like #275/#277) and only a successful TeamDelete exits cleanly. The real, narrower bug is timing — a slow-settling member outlasts the 60s `expectExit` budget while yy's retry-to-success loop is still (correctly) running. Reframed the design to keep retry-to-success and widen the exit budget + escalate to a longer settle past the fast cap, with AC-3 now a concrete code gate on the budget constant. Because this diverges from the approved A+C, it is flagged in-body and escalated to team-lead before the body is final.

## Stage Report: ideation (constraint fold — ≤60s budget guard)

- DONE: Fold the team-lead ≤60s constraint into the design without a validation surprise
  `TestNoTimeoutLiteralExceeds60s` (live_budget_test.go:30) AST-bounds every timeout literal in streamwatch.go/live_test.go to ≤60s. Reframed C to achieve slow-settle tolerance via `expectExit`'s NO-PROGRESS reset-on-activity (the pattern `expect`/`expectDispatchClose` already use), NOT a larger literal — `exitBudgetDefault` stays `60 * time.Second`, guard stays GREEN, ≤60s invariant untouched. AC-3 split into gate (a) reset-on-activity behavior + gate (b) the ≤60s guard stays green. Interaction noted explicitly in the test plan.

### Summary
Team-lead supplied the ≤60s budget-guard constraint (`TestNoTimeoutLiteralExceeds60s` over streamwatch.go/live_test.go). Rather than widen a literal (which would red the guard), the design now makes `expectExit` reset its deadline on stream activity like its sibling steps — slow-settle tolerance with no literal >60s, so the captain's ≤60s invariant is preserved, not revised. AC-3 now carries a second gate asserting the ≤60s guard stays green, so a future implementer cannot "fix" the settle by bumping the literal without a RED.

## Stage Report: ideation (rework — captain-approved test-infra best-effort teardown direction)

- DONE: Rework the body around the captain-approved pivot (test-infra grades best-effort teardown; launcher owns the exit)
  Replaced the no-progress-reset reframe with the approved direction: the live cycle grades correct bounded best-effort teardown (shutdown_request → member approval → bounded TeamDelete attempts) then kills via the existing `defer poller.kill()` (live_test.go:125-133) and PASSES — no exitCode==0 gate. Contract: bounded best-effort then STOP + terminal status; process exit is the LAUNCHER's job. Both prior dead directions (yy retry-to-success / A self-exit) marked unreachable/impossible with the root cause.
- DONE: Root cause cited from the REAL CI artifact + the raw spike, no further live spike
  SPIKE RESULTS now cites (a) the raw cap-confirm.jsonl harness-reinvocation proof (end_turn→system/init, EXIT=124, sonnet+opus) and (b) the CI artifact dead-but-listed timeline (member dies 14:37:30 / FO sees it active 14:38:25+, upstream #38116/#57681). Mechanism is settled; no further live spike.
- DONE: Reworked ACs as real code/test gates + flagged the open question
  AC-1 live grades best-effort+launcher-kill (no self-exit); AC-2 offline dead-but-listed fixture graded behaviorally (RED if FO skipped shutdown_request or TeamDelete); AC-3 contract oracle forbids retry-to-success + FO-self-exit framings. ≤60s guard stays green (no literal grows). Open question (real headless/cron -p consumers needing a clean exit → launcher-side timeout, NOT a registry hack) flagged non-blocking, default = live CI is the only consumer.

### Summary
Reworked the ideation around the captain's approved pivot after the spikes + the real CI artifact settled the mechanism: the dispatched member approves shutdown and dies, but Claude Code never clears it from members[] (upstream #38116/#57681), so TeamDelete never succeeds (yy unreachable) AND the harness won't let claude -p self-exit while members[] is populated (A impossible — proven on sonnet+opus, EXIT=124, raw jsonl preserved). The fix is test-infra: the live cycle grades CORRECT bounded best-effort teardown (a stronger assertion than exitCode==0) then the launcher's kill() ends the subprocess; the contract is corrected to bounded-best-effort-then-stop with launcher-owned exit. Three ACs are real code/test gates; the ≤60s guard is untouched; the only open question (real -p consumers gating on clean exit) is flagged non-blocking with the launcher-timeout answer.
