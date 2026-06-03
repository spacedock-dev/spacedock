---
id: atwf2w6p68t9q1mda790dcfc
title: "Non-interactive FO teardown — grade bounded best-effort + launcher-owned exit (dead-but-listed member: yy retry-to-success unreachable, FO self-exit impossible)"
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

The live-e2e sonnet cycle hangs at terminal teardown because of a verified upstream Claude Code defect, **#38116 / #57681** (an approved-shutdown member is never cleared from the team `members[]`). The mechanism, settled by the real failed CI artifact (the ensign member's own session jsonl, runs `26891717026` n3 + `26891718562` 2a, artifact `runtime-live-e2e-claude-live-sonnet`):

1. The dispatched ensign finished its work, RECEIVED the FO's cooperative `shutdown_request`, APPROVED it, and its session TERMINATED at **14:37:30**.
2. But the harness **never removed it from the team `members[]`**, so the FO's `TeamDelete` kept failing `Cannot cleanup team with 1 active member(s)` for ~55s against an already-DEAD member (FO still sees it active at **14:38:25+**).
3. The live cycle's flat-60s `expectExit` then killed the never-exiting subprocess → FAIL.

Two consequences, both empirically confirmed (see `## SPIKE RESULTS`):

- **yy's "retry `TeamDelete` to success" is UNREACHABLE** — success never comes, because the dead member is never removed from `members[]`. The retry loop runs until the 60s exit budget kills it.
- **The `claude -p` harness will NOT let the subprocess self-exit while `members[]` is populated** — a genuine `end_turn` does not terminate; the harness re-invokes the FO. Proven on sonnet AND opus (`EXIT=124`, raw stream preserved). So contract A's "hard-exit at the cap; residual cleanup at process death" is IMPOSSIBLE: there is no process death the FO can cause.

The fix therefore cannot demand a `TeamDelete` success or an FO self-exit. It grades the FO's CORRECT bounded best-effort teardown and makes the process exit the launcher's responsibility.

## Proposed approach (captain-approved, 2026-06-03): grade best-effort teardown; the LAUNCHER owns the exit

The exit cannot come from the FO (impossible) or from a `TeamDelete` success (unreachable). So we stop demanding either. Two coordinated changes:

1. **(Test / streamwatcher) Grade terminal teardown on CORRECT BOUNDED BEST-EFFORT BEHAVIOR, then the launcher kills the subprocess and the cycle PASSES.** The live cycle replaces its final `expectExit → exitCode==0` step (live_test.go:157-163) with: observe that the FO did the right sequence — sent `shutdown_request` → member approved → FO made BOUNDED `TeamDelete` attempts — then STOP the subprocess via the existing `defer poller.kill()` (live_test.go:125-133) and PASS. It MUST NOT require a clean self-exit, which the upstream bug makes impossible. This is a STRONGER assertion than `exitCode==0`: it checks the FO performed the correct teardown sequence, not that the process happened to exit. The grading step still uses the existing ≤60s no-progress budget per beat, so no timeout literal grows.
2. **(Contract — `first-officer-shared-core.md` + the Claude runtime adapter) Non-interactive terminal teardown is BOUNDED best-effort, then STOP and emit a terminal status; the PROCESS EXIT is the LAUNCHER's responsibility, NOT the FO's.** The FO: cooperative `shutdown_request` to the cohort → a BOUNDED set of `TeamDelete` attempts with the inter-attempt settle → then STOP retrying and emit a terminal status (it cannot delete a team the harness won't let it empty). The process death comes from the launcher (the live test's `kill()`, or a real automation's timeout) — the FO cannot cause it. After the bound, the harness will re-invoke the FO (proven); the FO just HOLDS (text only, no tool calls) until the launcher kills it. This **rehabilitates A's SPIRIT** — best-effort then process death — except the death is launcher-owned. REMOVE/correct yy's "keep retrying `TeamDelete` to success in non-interactive mode" (unreachable → the hang) and the disproven "hard-exit at the cap / residual cleanup at process death" (the FO cannot cause the death).

## Out of scope

- Reaching into Claude's private team registry to force-clear the dead member (captain rejected the config-rewrite host-hack).
- The opus path (it hits the SAME harness invariant; this fixes the grading, which is model-neutral).

## Acceptance criteria

**AC-1 — The sonnet live cycle PASSES via correct bounded best-effort teardown + launcher kill (runtime-observable oracle).**
Verified by: a live-e2e sonnet run where the FO performs the correct bounded teardown sequence (shutdown_request → member approval → bounded `TeamDelete` attempts) and the harness/launcher kills the subprocess — NO clean self-exit required, NO `exitCode==0` gate. The `expectExit`-style flat-cap-then-fail is gone. By-construction-pending-live; MUST pass on a real sonnet run (the oracle yy's fix failed).

**AC-2 — An offline fixture reproduces the dead-but-listed shape and the watcher grades it as correct best-effort teardown → PASS.**
Verified by: a Go test (sibling to / extending `TestSonnetTeamDeleteHangReplay` in `internal/ensigncycle/`) dripping a stream where the member APPROVES shutdown and its session ENDS, yet `TeamDelete` keeps failing `active member(s)`; the new grading step asserts the watcher OBSERVED the correct beats (shutdown_request sent + member approval + ≥1 bounded `TeamDelete` attempt) and PASSES (then the modelled `kill()` stops the proc). RED on a fixture where the FO did NOT do the correct teardown — never sent `shutdown_request`, OR never attempted `TeamDelete` — so the grade is behavioral, not "the proc was killed regardless." (Replaces the existing replay's `fakeProc`-never-exits tautology with a behavioral grade.)

**AC-3 — The contract mandates bounded best-effort teardown + launcher-owned exit, and FORBIDS the unreachable/impossible framings.**
Verified by: extending `TestTerminalTeardownRetriesToSuccess` (`skills/integration/terminal_teardown_retry_test.go`) so the required set asserts "bounded best-effort teardown" + a "terminal status, then hold" + a "the launcher/process owner ends the subprocess" clause, and the NEGATING set ADDS "retry … to success" (now unreachable) and "hard-exit at the cap" / "residual cleanup at process death" / "the FO exits the process" (now impossible). An inverted edit re-introducing retry-to-success OR FO-self-exit reds it. (Prose-only ceiling — the behavioral oracle is AC-1.)

## Test plan

- Spike: DONE (see `## SPIKE RESULTS`) — the riskiest unknowns are empirically settled by (a) the raw `cap-confirm.jsonl` harness-reinvocation reproduction (sonnet+opus, `EXIT=124`) and (b) the real CI artifact's dead-but-listed timeline (member dies 14:37:30 / FO still sees it active 14:38:25+, upstream #38116/#57681). NO further live spike needed for the mechanism.
- AC-2: offline Go test in `internal/ensigncycle/` (the dead-but-listed fixture + the new behavioral grade); sub-second, no model spend.
- AC-3: contract oracle in `skills/integration/terminal_teardown_retry_test.go` (extend the directional-mandate lint); sub-second.
- AC-1: one live-e2e sonnet cycle (~minutes, CI-only / Linux-bound), by-construction-pending-live.
- **≤60s guard stays green.** The new grading step uses the existing ≤60s no-progress budget per beat and the `defer poller.kill()` for exit, so no timeout literal grows; `TestNoTimeoutLiteralExceeds60s` (live_budget_test.go) is unaffected. Confirm it stays green at validation.
- **OPEN QUESTION (flagged, non-blocking; default = live CI is the only consumer):** are there real headless/cron `spacedock … -p` FO runs that gate on a CLEAN exit code? If so they need a launcher-side timeout (a wrapper that bounds the run and kills it), NOT a contract change and NOT any reach into Claude's private registry. Default: only the live-e2e CI consumes the FO `-p` exit, and it already owns a `kill()`, so no extra launcher work is needed.
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

**What these jointly INVALIDATE** (the canonical current design lives in `## Proposed approach` / `## Acceptance criteria` up top; this is the evidence behind it):
- yy's "retry `TeamDelete` to success" is **UNREACHABLE**: success never comes, because the dead member is never removed from `members[]`. The retry loop runs forever (bounded only by the 60s exit budget that then kills it).
- contract A's "hard-exit at the cap; residual cleanup at process death" is **IMPOSSIBLE**: the harness will not let `claude -p` self-exit while `members[]` is non-empty (fact 1), so there is no process death the FO can cause.

(The superseded directions these displace — the original A+C and the no-progress-reset reframe — are quarantined in `## Superseded directions (provenance)` near the bottom so exactly one current AC set exists.)

## Superseded directions (provenance)

These directions were explored and ruled out by the spikes + the CI artifact. Kept as the audit trail; they are NOT the current design (see `## Proposed approach` / `## Acceptance criteria` at the top).

1. **Original seed A+C — "hard-exit at the cap" (contract A) + "accept best-effort-teardown-then-exit" (test C).** RULED OUT: A is impossible — the `claude -p` harness will not let the subprocess self-exit while `members[]` is populated (proven, `EXIT=124` on sonnet+opus), so there is no process death the FO can cause and "residual cleanup at process death" never happens.
2. **No-progress-reset reframe — make `expectExit` reset its deadline on stream activity so yy's retry-to-success reaches success within budget.** RULED OUT: it only widens the wait for a `TeamDelete` success that NEVER arrives (the dead member is never cleared from `members[]`, upstream #38116/#57681). It treated the bug as slow-settle timing; the CI artifact proved it is a never-settle registry defect. The ≤60s-guard analysis it produced still holds and carried forward (the current grading step likewise grows no literal).

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

## Stage Report: ideation (consolidation — one coherent current AC set)

- DONE: Make the canonical sections SINGULAR and CURRENT; quarantine superseded material as labeled provenance
  Replaced the dead `## Problem` (contract-ambiguous framing) with the verified dead-but-listed root cause (#38116/#57681); replaced the dead `## Proposed approach (A+C)` with the captain-approved best-effort-teardown+launcher-exit direction; DELETED the dead AC-1..AC-4 ("process EXITS"/"hard-exit at the cap") so exactly ONE `**AC-N` set (the 3 approved ACs) exists; replaced the old `## Test plan`. Moved the original A+C and the no-progress-reset reframe into `## Superseded directions (provenance)`. SPIKE RESULTS + Stage Reports kept as the audit trail. Frontmatter title left for the FO to fix.
- DONE: Verify structure for the gate AC cross-check + 2a `ClassifyEntityACs` extractor
  Confirmed exactly 3 `**AC-N` blocks (grep), one each of Problem/Proposed approach/Acceptance criteria/Test plan; the remaining "hard-exit at the cap" mentions are all in the FORBIDS/negating context, not prescriptive. Baselines green (TestTerminalTeardownRetriesToSuccess + TestNoTimeoutLiteralExceeds60s).

### Summary
Consolidation pass before staff review: I had appended the approved direction rather than replacing the superseded seed, leaving two contradictory `## Acceptance criteria` sets and dead canonical sections that would break the gate's AC cross-check and 2a's ClassifyEntityACs extractor. Rewrote Problem/Proposed approach/Acceptance criteria/Test plan to be singular and current (the captain-approved best-effort-teardown + launcher-owned-exit direction grounded in the #38116/#57681 dead-but-listed root cause), and quarantined the original A+C and the no-progress-reset reframe into a labeled `## Superseded directions (provenance)` section. The file now reads top-to-bottom as one coherent current ideation with exactly one AC set; nothing is lost — the dead directions and full spike/CI evidence remain as audit trail.
