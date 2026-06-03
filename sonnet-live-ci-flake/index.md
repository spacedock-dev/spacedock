---
id: yyqez6npx8qb11b5v7fgwjtf
title: Sonnet live-CI fails with a repeatable shape — FO subprocess ends its turn after shutdown_request before the streamwatcher's expected dispatch-close fires
status: implementation
source: session-10 — identical failure shape on PR #275 (n3 frontmatter-hash-quoting) and PR #277 (2a require-external-proof-guard); both offline-green + opus-green + sonnet-red. Two matching failures = mechanism issue, not random flake. Blocks merging n3 + 2a into 0.19.5.
score: "0.19"
worktree: .worktrees/spacedock-ensign-sonnet-live-ci-flake
started: 2026-06-03T07:09:59Z
completed:
verdict:
issue:
---

The live-e2e CI cycle passes on offline and on opus but **fails on sonnet with a repeatable shape** across two unrelated PRs (n3 #275, 2a #277). The cycle drives the entity to `done` correctly; the FO subprocess self-reports `terminal_reason: completed`, `is_error: false`. But the FO **never exits** — it gets stuck at terminal teardown when its `TeamDelete` fails ("active member(s)") and it does not retry to success. The test fails at `expectExit` ("FO subprocess did not exit within 1m0s"), not at the dispatch-close. Two identical failures on different PRs (plus an identical opus path that only passes by luck) confirm a mechanism issue, not a flake.

## Root cause (CONFIRMED from archived JSONL — both hypotheses below REFUTED)

Archived sonnet runs traced: #277 = CI run `26865865357` artifact `runtime-live-e2e-claude-live-sonnet`; #275 = run `26865122632` (sonnet job cancelled but artifact retained). Both fail **identically** at `expectExit`. The opus job on #277 (run `26865865357`) passes — its stream shows the *same* defect but recovers.

The observed event ordering the streamwatcher consumed (#277 sonnet transcript line numbers):

| Event | What happens |
|---|---|
| `L92` TeamCreate | `expect(isTeamCreate)` passes |
| `L214` Agent dispatch | dispatch opens |
| `L253` SendMessage `shutdown_request` | **PREMATURE** — sent ~16s after dispatch, BEFORE the ensign's `Done:` (violates `## Awaiting Completion`) |
| `L271` `system task_notification status=completed` (keyed on the Agent tool_use_id) | **dispatch-close fires here** → `expectDispatchClose` passes. It fires because the premature shutdown made the ensign terminate early. |
| `L323` TeamDelete | terminal teardown |
| `L324` TeamDelete **FAILS** | `"Cannot cleanup team with 1 active member(s): spacedock-ensign-make-it-work-done"` — races the ensign's shutdown settling in the team registry |
| `L339` SendMessage `shutdown_request` (2nd) | sent to the already-gone ensign; never answered |
| `L341` FO turn ends `terminal=completed` | **subprocess goes idle and never retries TeamDelete → never exits** |

`expectExit` waits 60s, never sees exit, kills the subprocess → FAIL.

**The opus/sonnet delta (decisive):** opus hits the *exact same* premature-shutdown + TeamDelete-fails-with-active-member sequence, but opus **retries `TeamDelete` a third time after the registry settles** (#277 opus L160 fail → L166 retry succeeds → subprocess exits → PASS). Sonnet, after `TeamDelete` fails, re-sends a `shutdown_request` and ends its turn **without ever retrying `TeamDelete`**, so the team is never torn down and the `claude -p` subprocess hangs forever. #275 sonnet retried 3× but all 3 raced *before* the registry settled, then stopped — same hang.

So the bug is a **two-part FO-contract defect**, and the failure point is `expectExit`:
1. **Trigger:** the FO emits a premature `shutdown_request` right after dispatch (before the completion signal), which the `## Awaiting Completion` clause is *supposed* to prevent but sonnet violates. This forces the teardown race to start early.
2. **Hang:** the terminal-teardown `TeamDelete` races the ensign's cooperative shutdown and fails with "active member(s)", and the FO contract (shared-core Merge-and-Cleanup **step 10** / `## Awaiting Completion`) has **no clause that deterministically retries `TeamDelete` until the team is gone**. Opus improvised the retry; sonnet did not.

### Hypotheses from the seed — both REFUTED by the trace
1. ~~Streamwatcher budget mismatch on sonnet's verbosity.~~ **Refuted.** The dispatch-close (`task_notification completed`) fires well within budget (`L271`); the failure is `expectExit`, not a budget trip on the close.
2. ~~Contract ambiguity about turn-ending after `shutdown_request`.~~ **Refuted as the cause.** Ending the turn after `shutdown_request` is correct per contract; the close IS observed. The real gap is the missing *retry-TeamDelete-to-success* loop in terminal teardown.

### Spike (riskiest mechanism) — DONE, mechanism validated
The riskiest unknown was "can the archived JSONL reproduce the hang offline through the real streamwatcher?" Validated: extracting the captured stream-json lines (the `live_test.go:127` tee lines) and **dripping them incrementally** into a `fakeLineSource` reproduces the exact behavior — `expect(isTeamCreate)` ok, `expectDispatchClose` fires (closedCount=1 via the task_notification anchor), `expectExit` hangs and kills the never-exiting proc. **Note from the spike:** an all-at-once push misrepresents the bug (it makes `expectDispatchClose` capture its baseline *after* the close was already counted, so the close looks like it never fired). The replay MUST drip lines one-per-poll to faithfully reproduce; the regression fixture below encodes that.

## Proposed approach

1. **FO contract fix (the real defect).** In the shared-core terminal teardown (`first-officer-shared-core.md` Merge-and-Cleanup step 10) and the Claude FO runtime, make terminal team teardown a **bounded retry-to-success loop**: after the cooperative `shutdown_request`, the FO MUST retry `TeamDelete` (re-sending `shutdown_request` to any named active member between attempts) until `TeamDelete` succeeds or a small attempt cap is hit, rather than ending the turn on the first "active member(s)" failure. Specific before/after wording to be drafted at implementation; the load-bearing change is "first-officer-shared-core.md:226 step 10 must not stop at a failed TeamDelete."
2. **Tighten the premature-`shutdown_request` guard.** The premature shutdown at `L253` (16s post-dispatch) is what starts the teardown race early. The `## Awaiting Completion` decision procedure already bans this; assess whether the sonnet violation needs a sharper clause or is subsumed by the step-10 retry fix (the retry fix alone may be sufficient — even with the premature shutdown, a retry-to-success teardown exits cleanly, as opus demonstrates).
3. **Add the offline regression fixture** (AC-2a): a Go test in `internal/ensigncycle` that drips the archived sonnet stream through the streamwatcher and asserts the diagnosis — TeamCreate matches, dispatch-close fires, and `expectExit` trips a `stepTimeout` that kills the never-exiting proc. This locks the watcher's localization and the captured evidence in the repo.
4. **Re-trigger live CI on #275 and #277** (AC-1, the only true oracle): both must pass sonnet + opus + offline after the contract fix.

## Out of scope

- Merging n3 / 2a themselves — that is the follow-on once this fix lands (tracked on those PRs).
- The `am` streamwatch rename (test-only file → `*_test.go`) — orthogonal; do not couple.
- Changing the streamwatcher's close/exit logic — the trace proves the watcher is **correct**; it accurately localized a real hang. No `streamwatch.go` logic change is in scope.
- Broader live-CI redesign; this is a targeted FO-contract fix for the terminal-teardown hang.

## Acceptance criteria

The bug is an FO-contract defect whose only true oracle is a live CI run (`looks-right-here ≠ runs-right-everywhere`). The streamwatcher is not buggy, so there is no streamwatcher red→green test to write — see AC-2's honest reframe.

**AC-1 — The sonnet live cycle reaches its bounded stop condition (the FO subprocess exits) on the failing surface.**
Verified by: a live-e2e CI run on sonnet (on n3 #275 and 2a #277, or an equivalent reproduction) passing — specifically `expectExit` returns 0 instead of timing out, because terminal `TeamDelete` now retries to success. The run is the oracle.

**AC-2 — An offline regression fixture replays the archived sonnet stream and pins the diagnosis.**
Verified by: a Go test in `internal/ensigncycle` that drips the captured sonnet stream-json (committed as a testdata fixture) through the streamwatcher and asserts: `expect(isTeamCreate)` matches, `expectDispatchClose` fires (the close anchor is fine), and `expectExit` trips a `stepTimeout` that kills the never-exiting proc. This is a real, checkable artifact (real captured stream + real watcher) that fails if the watcher's localization regresses.
- HONEST LIMIT: this fixture does NOT flip red→green on the FO fix — the fix changes a *future live FO stream*, not this recording, so a recorded stream cannot demonstrate the fix. The original AC-2 ("red before fix, green after" in the streamwatcher) rested on the false premise that the streamwatcher is buggy; it is not. The fix's red→green proof is AC-1 (the live run). The implementation stage should confirm with the captain that AC-2's offline scope is "pin the diagnosis + guard the watcher," with AC-1 as the fix oracle.

**AC-3 — Opus and offline cycles stay green (no regression).**
Verified by: the opus live run still passing + `go test ./...` still passing after the contract fix. (Opus already passes; the step-10 retry clause must not break its working teardown.)

## Test plan

- Diagnostic (DONE): archived JSONL from the failed sonnet runs (#275, #277) traced; root cause confirmed at `expectExit` (TeamDelete-fails-then-no-retry), both seed hypotheses refuted. Spike validated the offline incremental-drip replay reproduces the hang.
- Regression (AC-2): commit the captured sonnet stream as `internal/ensigncycle/testdata/` and a Go test that drips it through the watcher, asserting close-fires + exit-hangs-and-is-killed. Cheap, offline, in `go test ./...`.
- Confirmation (AC-1, AC-3): live-e2e CI on sonnet + opus. Cost: medium-high (live cycles). The contract fix is the only thing that flips sonnet green.
- High-stakes surface (CI/release machinery + FO contract) → detached adversarial audit required before merge per the dev README's validation stage.

## Notes

- Failure evidence: PR #277 sonnet run `26865865357` (opus on the same run PASSED), PR #275 run `26865122632` (sonnet cancelled, artifact retained). Both sonnet: offline ✅, opus ✅, sonnet ❌ at `expectExit`. Artifacts downloaded and traced this stage.
- Exact failure line (both): `live_test.go:159: live cycle failed waiting for FO exit: FO subprocess did not exit within 1m0s.` — NOT the dispatch-close.
- The fix surface is the FO contract prose (`skills/first-officer/references/first-officer-shared-core.md` step 10 + `claude-first-officer-runtime.md` `## Awaiting Completion`), NOT `streamwatch.go`. `TeamDelete` is a Claude-runtime tool with no spacedock-binary wrapper, so there is no binary-level guard to add — the retry-to-success discipline lives in contract prose, and its proof is the live run (AC-1).
- This is the gating fix that carries n3 + 2a back into mergeable shape for 0.19.5.
- Related context: `internal/ensigncycle/streamwatch.go` is the Go port of the upstream Python `FOStreamWatcher`; the `am` entity proposes renaming it to `*_test.go` but does not change its logic — keep the two entities decoupled. The `8y` shared-scenarios surface and this entity both touch `internal/ensigncycle` testdata — coordinate at implementation to avoid worktree conflicts.

## Stage Report: ideation

- DONE: Root-cause the flake from the archived sonnet JSONL (PRs #275/#277): confirm which hypothesis holds — streamwatcher bounded-stop budget vs FO-contract turn-ending-after-shutdown_request — with the event trace recorded in the body.
  Downloaded + traced both sonnet artifacts; BOTH seed hypotheses REFUTED. Real cause: failure is at `expectExit` (not dispatch-close) — terminal `TeamDelete` fails "active member(s)" and sonnet does not retry to success; opus passes only because it retries. Trace table + opus/sonnet delta recorded in the body's "Root cause" section.
- DONE: Pin the fix approach AND a regression test that reproduces the sonnet shutdown_request -> dispatch-close ordering offline (red before the fix, green after).
  Fix pinned: FO-contract step-10 terminal teardown must retry `TeamDelete` to success (body "Proposed approach"). Regression-test approach validated via a throwaway spike (incremental-drip replay of the captured stream through the real streamwatcher reproduces close-fires + exit-hangs). HONEST REFRAME recorded in AC-2: the streamwatcher is NOT buggy, so no streamwatcher red→green test exists; the offline fixture pins the diagnosis, and the fix's red→green oracle is the live run (AC-1). Flagged for captain confirmation at implementation.

### Summary
The seed framed this as a streamwatcher-budget or turn-ending-after-shutdown bug; the archived JSONL refutes both. The cycle completes correctly and the dispatch-close fires — the FO subprocess hangs at terminal teardown because its `TeamDelete` races the ensign's cooperative shutdown, fails with "active member(s)", and sonnet (unlike opus) never retries `TeamDelete` to success, so `claude -p` never exits and `expectExit` kills it after 60s. A premature post-dispatch `shutdown_request` (an `## Awaiting Completion` violation) starts the race early. Fix is FO-contract prose (step-10 retry-to-success); the streamwatcher is correct and untouched. AC-2 was reframed honestly: an offline fixture replay pins the diagnosis, but only the live run (AC-1) can prove the contract fix — flagged for captain confirmation since this changes AC-2's original "red→green offline" intent.
