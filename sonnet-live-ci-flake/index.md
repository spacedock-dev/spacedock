---
id: yyqez6npx8qb11b5v7fgwjtf
title: Sonnet live-CI fails with a repeatable shape — FO subprocess ends its turn after shutdown_request before the streamwatcher's expected dispatch-close fires
status: validation
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

## Stage Report: implementation

- DONE: Make terminal TeamDelete a bounded retry-to-success loop in first-officer-shared-core.md step 10 + the Claude FO runtime — do NOT end the turn on the first failed TeamDelete; retry (re-sending shutdown_request to active members) until success or a small cap.
  Commit `1367f7b8`: shared-core Merge-and-Cleanup step 10 now requires a bounded retry-to-success team teardown (re-issue cooperative shutdown to still-active members between attempts, cap reached → surface to captain); new Claude-runtime `## Terminal Team Teardown` section realizes it for `TeamDelete` (4-step procedure, ≈3–5 attempt cap), and the `## Awaiting Completion` TeamDelete ban is scoped to the pre-completion phase so the two phases no longer read as contradictory. Prose-oracle lint `skills/integration/terminal_teardown_retry_test.go` (2 tests) pins the retry clause in both files and pins that the pre-completion ban survives — both red before the edit, green after.
- DONE: Land the offline regression fixture (AC-2) in internal/ensigncycle — drip the archived sonnet stream-json (committed testdata) through the streamwatcher and assert expectExit trips the stepTimeout that kills the never-exiting proc.
  Commit `337685d8`: `internal/ensigncycle/testdata/sonnet_teamdelete_hang.stream.jsonl` is the verbatim 341-line captured sonnet stream from the failing CI run (PR #277, run 26865865357, artifact `runtime-live-e2e-claude-live-sonnet`, re-downloaded via gh and tee-extracted from `live-e2e-transcript.txt`). `TestSonnetTeamDeleteHangReplay` drips it one-line-per-poll through the real `streamWatcher` and asserts TeamCreate matches, the dispatch-close fires (closedCount=1 via the keyed task_notification), and `expectExit` trips a `stepTimeout("expect_exit")` that kills the never-exiting proc. Non-vacuity proven by mutating the closedCount assertion (failed reporting the real count) before reverting.

### Summary
Both deliverables landed and committed in the code worktree. The FO-contract fix (commit `1367f7b8`) adds the missing retry-to-success discipline to terminal team teardown in both the shared core (step 10) and the Claude runtime (new `## Terminal Team Teardown` section), and disambiguates the two TeamDelete phases (pre-completion ban vs. terminal retry) so they no longer contradict. The AC-2 regression fixture (commit `337685d8`) is a real captured sonnet stream replayed through the real streamwatcher offline, pinning the watcher's hang localization. AC-2's offline scope ("pin the diagnosis + guard the watcher," AC-1 as the fix oracle) was settled in the dispatch checklist/summary — no separate captain round-trip needed; the honest limit (a recorded stream cannot flip red→green on a producer-side FO fix) is stated in the test and commit. Full offline suite green (`go test ./...` 836 passed); live-tagged build vets clean. AC-1/AC-3 confirmation (live-e2e on sonnet+opus for #275/#277) is the validation-stage oracle and cannot run from this stage.

## Stage Report: validation

- DONE: Reproduce evidence for every AC: run TestSonnetTeamDeleteHangReplay (AC-2 offline fixture) and go test ./... (AC-3 opus+offline); confirm AC-1 is by-construction-pending the live sonnet run; flag any AC without reproduced evidence.
  AC-2: `TestSonnetTeamDeleteHangReplay` PASS (real watcher + real 341-line captured stream); adversarial mutation (assert proc exited cleanly) flipped it red, then reverted — non-vacuous. AC-3: `go test ./...` 10/10 packages green, `go vet -tags live ./internal/ensigncycle/` exit 0. AC-1: by-construction pending — its only true oracle is a live sonnet cycle (no recording can prove a producer-side FO fix); the FO triggers it post-validation. No AC lacks evidence.
- DONE: Confirm the FO-contract retry-to-success fix is sound and the pre-completion-ban vs terminal-retry disambiguation is not self-contradictory (read both edited contract sections in first-officer-shared-core.md + claude-first-officer-runtime.md).
  Sound: step 10 (shared-core) + new `## Terminal Team Teardown` (claude runtime) mandate "do NOT end the turn on the first active-member(s) failure; re-issue cooperative shutdown to still-named members and retry TeamDelete until success or a small cap" — exactly the opus-improvised path that passed. Not self-contradictory: the two phases are split by the completion-signal boundary, cross-reference each other explicitly (awaiting-completion L225 → terminal section; terminal L253 → awaiting-completion), and L244 states "the ban is pre-completion, the retry is terminal." No single situation triggers both rules.

### Summary
PASSED (offline + opus + fixture), with AC-1 by-construction-pending the live sonnet run the FO triggers next. AC-2's reframe HOLDS and is not rubber-stamped: the fixture's actual event ordering matches the diagnosed trace exactly — TeamCreate (L91) → Agent dispatch (L213) → premature shutdown_request (L252, before the completion signal) → task_notification completed (L270) → TeamDelete (L322) → active-member fail (L323) → 2nd shutdown_request (L338) → turn ends WITHOUT retrying TeamDelete (L340). That last "no-retry-after-fail" is the real defect; a recording of the buggy producer cannot demonstrate the producer-side fix, so AC-2 legitimately only pins the watcher's localization while AC-1 is the fix oracle. Change set is exactly scoped (5 files: 2 contract-prose, fixture+replay, prose-oracle lint) — no `streamwatch.go` logic change, no `am`/n3/2a coupling, per the out-of-scope list. High-stakes surface (FO contract + CI machinery): a detached adversarial read-only audit of the merge result is REQUIRED before merge per the dev README validation stage — my own AC-2 mutation was a spot-check, not that full audit.

## Feedback Cycles

### Cycle 1 — detached adversarial audit (2026-06-03) — MATERIAL

Validation PASSED, but the detached audit found two test-strength holes plus a fix-completeness gap that is load-bearing for the live AC-1 oracle:

- **MATERIAL — the contract oracles are substring-presence lints, not meaning oracles.** `TestTerminalTeardownRetriesToSuccess` checks only `strings.Contains` for tokens. An **inverted** mandate ("End the turn on the first failure; do NOT re-send shutdown_request; do NOT call TeamDelete again" — the exact #275 bug) passes GREEN in BOTH `first-officer-shared-core.md` and `claude-first-officer-runtime.md` because the grep tokens survive. The "locked by this test" claim is overstated.
  **Fix:** strengthen the oracle to detect inversion — assert the directional mandate AND the ABSENCE of the negating phrases (ban e.g. "do not retry", "end the turn on the first failure", "do NOT call TeamDelete again") within the terminal-teardown region, scoped to step 10 / the `## Terminal Team Teardown` step (not the whole `## Merge and Cleanup` region). The two inverting edits the audit applied must turn the oracle RED.

- **MATERIAL — the AC-2 replay step-3 is tautological.** `TestSonnetTeamDeleteHangReplay`'s `fakeProc` never exits, so `expectExit` ALWAYS trips `stepTimeout` regardless of stream content; truncating the captured stream to drop the real `TeamDelete`-failure (L323 "active member(s)") and the no-retry turn-end (L341) leaves the test GREEN. Step 3 is proven by the test's own never-exiting proc, not by the recording.
  **Fix:** make step 3 depend on the captured failure — assert the watcher OBSERVED the `TeamDelete`-failure line (L323) and/or drive proc-exit from the stream so removing the failure evidence changes the outcome. Truncating the stream to drop L271–341 must turn the test RED.

- **POLISH but load-bearing for AC-1 — the fix may not close the SONNET case.** (a) No inter-attempt settle/backoff between `TeamDelete` retries → an agent can fire all ~3–5 attempts in milliseconds and exhaust the cap before any member leaves the roster (exactly #275's "retried 3x but all raced before the registry settled, then stopped"). (b) Cap-exhaustion in non-interactive `claude -p` mode → "surface to captain" = emit text + end turn, leaving the team un-deleted (still hangs).
  **Fix:** prescribe an inter-attempt settle/wait in the contract (do NOT fire retries in parallel/instantly), and define the non-interactive cap-exhaustion behavior so the subprocess does not hang. This is what determines whether the live AC-1 sonnet run actually passes — do not skip it.

## Stage Report: implementation (cycle 2)

- DONE: Finding 1 (MATERIAL) — strengthen `TestTerminalTeardownRetriesToSuccess` so an inverted mandate turns it RED.
  Commit `c95f0809`: the oracle now asserts BOTH directions — positive mandate present AND negating phrases (`end the turn on the first`, `do not call teamdelete again`, `do not re-send`, `do not retry`, `stop at the first`) absent — scoped to step 10 alone (new `numberedStep` helper) and the `## Terminal Team Teardown` section, not the whole region. PROVEN: applying the audit's two inverting edits to both contract files turned it RED (inverted-phrase-present + positive-phrase-missing errors), then reverted.
- DONE: Finding 2 (MATERIAL) — make `TestSonnetTeamDeleteHangReplay` step-3 depend on the captured failure.
  Commit `c95f0809`: step 3 now records the watcher's teed transcript and asserts it OBSERVED the diagnosed beats in order after the dispatch-close — terminal `TeamDelete` tool_use, its `Cannot cleanup team with 1 active member(s)` failure, the `terminal_reason":"completed"` no-retry turn-end — plus exactly one terminal `TeamDelete` (a recovered/opus stream would show a retried second one). PROVEN: truncating the fixture to drop L271–341 (the failure evidence) turned it RED, then reverted. The never-exiting `fakeProc` is no longer the sole proof.
- DONE: Finding 3 (load-bearing for AC-1) — prescribe inter-attempt settle + non-interactive cap-exhaustion behavior.
  Commit `c95f0809`: both contract surfaces (shared-core step 10 + Claude `## Terminal Team Teardown`) now require a SETTLE between attempts (re-issue cooperative shutdown, then `Bash("sleep 2")` before the next `TeamDelete` — explicitly distinguished from the banned `## Awaiting Completion` idle-poll sleep) so retries don't millisecond-re-race the same registry update (#275). Cap-exhaustion now MUST drive the subprocess to exit: non-interactive `claude -p` has no operator to read a "surfaced" failure, so the FO keeps settling+retrying `TeamDelete` until success (the cap bounds the FAST retries, not the teardown obligation); only a `TeamDelete` success lets the FO end its turn.

### Summary
All three cycle-1 audit findings closed in commit `c95f0809`. Both test-strength holes are now proven non-vacuous by adversarial mutation: the contract oracle turns RED under the audit's inverting edits (it asserts negating-phrase absence, not just positive-token presence, scoped to the single teardown step), and the AC-2 replay turns RED when the fixture is truncated to drop the captured TeamDelete-failure (it asserts the watcher observed the diagnosed beats in order, not merely that a never-exiting proc trips expectExit). The fix-completeness gap is closed in both contract surfaces — an inter-attempt settle makes the retry effective rather than a millisecond re-race, and a non-interactive exit obligation keeps cap-exhaustion from degrading to an end-turn hang in `claude -p`. Full offline suite green (`go test ./...` 836 passed); `go vet -tags live` clean. AC-1 (the live sonnet+opus run on #275/#277) remains the behavioral oracle for the fix and is the validation/CI-trigger step — it cannot run from this stage.

### Cycle 2 — detached re-audit (2026-06-03) — MATERIAL

Re-validation independently reproduced the cycle-1 adversarial edits in a throwaway detached worktree (`/tmp/sonnet-flake-audit` at `c95f0809`). Findings 1 and 2 are CLOSED; finding-3's PROSE is present and correct, but its ORACLE coverage is NOT — checklist item 2 ("an oracle covers them") fails.

- **CLOSED — finding 1 (inversion-resistant contract oracle).** Independently re-applied the two inverting rewrites. Inverting shared-core step 10 → `TestTerminalTeardownRetriesToSuccess` RED (negating phrase `do not re-send` present + positive mandates `retry the team-teardown call to success` / `Do NOT end the turn` missing). Inverting the Claude `## Terminal Team Teardown` steps → RED (`end the turn on the first` + `do not re-send` present, `do NOT end the turn` missing). A bare substring lint would have passed both; the strengthened two-sided lint catches them.
- **CLOSED — finding 2 (non-tautological AC-2 replay).** Truncating the fixture to drop L271–341 (the dispatch-close onward, incl. the L323 active-member failure) → `TestSonnetTeamDeleteHangReplay` RED (missing `"name":"TeamDelete"` beat; `countOccurrences == 0`). Sharper proof: removing ONLY the L323 `Cannot cleanup team with 1 active member(s)` line while keeping the TeamDelete call AND the never-exiting `fakeProc` → still RED (the diagnosed-beat-in-order assertion fails). The never-exiting proc is no longer the sole proof.
- **MATERIAL — finding-3 is present in prose but NOT covered by any oracle.** The settle + non-interactive-exit clauses ARE in both contract surfaces (correct for the live AC-1 behavior — the live FO reads prose, not the oracle). But `TestTerminalTeardownRetriesToSuccess` does NOT pin them. Gutting BOTH contract files' teardown steps back to the exact #275 hang — `immediately call TeamDelete again` (the millisecond re-race) + cap-exhaustion → `surface to the captain and end the turn` (the non-interactive hang) — leaves the oracle GREEN. Cause: the shared-core required-phrase `"settle"` is satisfied by the trailing delegation sentence `"the settle interval, the cap value, and the non-interactive exit obligation are the adapter's"` (which NAMES the concepts while delegating them away), not by the behavioral mandate; the Claude-runtime required set (`TeamDelete`, `active member`, `do NOT end the turn`, `cap`) checks no settle / non-interactive-exit token at all. So a future edit can silently revert finding-3's behavior with the full suite staying green — an assertion that green-lights the exact regression the fix exists to prevent.
  **Fix:** extend `TestTerminalTeardownRetriesToSuccess` to pin finding-3's BEHAVIOR, not a delegated keyword: in shared-core step 10 require the actionable settle mandate (e.g. assert the settle phrase appears OUTSIDE the trailing `are the adapter's` delegation sentence, or require the explicit non-interactive-exit obligation `MUST keep settling and retrying ... until it succeeds`); in the Claude `## Terminal Team Teardown` steps require the inter-attempt settle (`sleep 2` / `wait for the settle before the next`) AND the non-interactive exit obligation (`keep the settle-then-... loop going ... until TeamDelete succeeds`). Both my gutting edits — `immediately call TeamDelete again` + `surface to the captain and end the turn` — must then turn the oracle RED.

Findings 1 and 2 hold and the live-AC-1 behavior fix (finding-3 prose) is present, so the live oracle is not blocked by this. The gate is held because checklist item 2 (an oracle covers finding-3) is unmet — a test-strength hole that green-lights a future revert of the settle/non-interactive-exit clauses. Routes back to implementation for the oracle extension above; no contract-prose change is needed.

## Stage Report: validation (cycle 2)

- DONE: Independently re-run the cycle-1 audit's adversarial edits in a throwaway worktree: (a) the two inverted-mandate rewrites, (b) the AC-2 fixture truncation dropping L271-341 — confirm the strengthened oracle and the de-tautologized replay now go RED on ALL of them.
  All RED as required. (a) shared-core step-10 inversion → oracle RED; Claude `## Terminal Team Teardown` inversion → oracle RED (each independently). (b) truncate-to-270-lines → replay RED; plus surgical removal of ONLY the L323 active-member failure line → replay RED (proves non-tautological). Baseline GREEN before each edit, reverted after; detached worktree `/tmp/sonnet-flake-audit` at `c95f0809`, left clean. The cycle-1 holes for findings 1 and 2 are closed.
- FAILED: Verify finding-3 is real, not just present: the inter-attempt settle clause + the non-interactive cap-exhaustion exit obligation are in BOTH contract surfaces and an oracle covers them; confirm go test ./... green and no AC regressed.
  The CLAUSES are in both surfaces (correct for live AC-1). But NO oracle covers them: gutting both files' teardown steps to `immediately call TeamDelete again` + cap-exhaustion `surface to the captain and end the turn` (the exact #275 hang) leaves `TestTerminalTeardownRetriesToSuccess` GREEN — the `"settle"` token is anchored to the trailing `are the adapter's` delegation sentence, and the Claude side checks no settle/non-interactive token. `go test ./...` 836 passed and `go vet -tags live` clean (no AC regressed), but checklist item 2's "an oracle covers them" is unmet. See Cycle-2 feedback finding for the fix.

### Summary
REJECTED. Independently reproduced all of cycle-1's adversarial edits: findings 1 (inversion-resistant oracle) and 2 (non-tautological replay) are genuinely CLOSED — both turn RED exactly where they must, and a sharper surgical edit (removing only the L323 active-member failure line) confirms the replay is no longer tautological. Finding-3's PROSE (inter-attempt settle + non-interactive exit obligation) is present and correct in both contract surfaces, so the LIVE AC-1 fix behavior is in place. But checklist item 2 fails: no oracle COVERS finding-3 — I gutted both files' teardown steps back to the #275 millisecond-re-race + non-interactive hang and the suite stayed GREEN, because the `"settle"` requirement is satisfied by a non-load-bearing `"...are the adapter's"` delegation sentence and the Claude side checks no settle/non-interactive token. That is a test-strength hole that green-lights the exact regression the fix prevents. AC-1 stays live-pending (not a blocker); AC-2 and AC-3 hold (`go test ./...` 836 passed, `go vet -tags live` clean). Routes back to implementation for the oracle extension; no contract-prose change needed.

### Cycle 2 — detached re-audit (in cycle-2 validation, 2026-06-03) — MATERIAL

Findings 1 & 2 from cycle 1 are CLOSED (the re-audit independently reproduced both inverted mandates → oracle RED, and the fixture truncation + L323 removal → replay RED). One NEW hole remains:

- **MATERIAL — finding-3's settle + non-interactive-exit prose has no load-bearing oracle.** The inter-attempt settle clause and the non-interactive cap-exhaustion exit clause ARE present in both contract surfaces (so the live AC-1 behaviour fix is in place), but `TestTerminalTeardownRetriesToSuccess` does NOT cover them: gutting both files' teardown steps back to the #275 "immediately retry" millisecond-re-race + "surface to captain and end the turn" non-interactive hang leaves the oracle GREEN. The "settle" token is satisfied by a non-load-bearing "...are the adapter's" delegation sentence; the Claude side checks no settle / non-interactive token. So an inversion of the exact regression the fix prevents passes green.
  **Fix (oracle extension only — do NOT change the contract prose; the prose is correct):** extend the oracle so it requires the LOAD-BEARING settle clause AND the non-interactive cap-exhaustion clause within the teardown step (scoped so a delegation sentence cannot satisfy it). While you're there, make the oracle comprehensive over EVERY load-bearing teardown clause (retry-to-success, settle-between-attempts, non-interactive cap-exhaustion exit) so prose and oracle can't drift again. Reproduce the audit's gut-edit — revert both teardown steps to immediate-retry + surface-and-end-turn — and confirm the strengthened oracle goes RED on each clause.

This is **feedback cycle 2** (next reject escalates to the captain). AC-1 remains live-pending (not a blocker). Re-run the targeted oracle + `go test ./...`, append a cycle-3 Stage Report, commit on your worktree branch, then SendMessage(to="team-lead", message="Done: ...").

## Stage Report: implementation (cycle 3)

- DONE: Extend `TestTerminalTeardownRetriesToSuccess` to cover finding-3's load-bearing settle + non-interactive-exit clauses, scoped so the delegation sentence can't satisfy it (ORACLE ONLY — no contract-prose change).
  Commit `c9cfb265`: new `stripDelegationTail` drops the trailing `…are the adapter's` sentence from shared-core step 10 before asserting (so a delegation-only mention no longer counts — it targets the LAST occurrence, since step 10 also says "roster and decomposition are the adapter's" early). shared-core step 10 now requires the behavioral settle mandate (`wait a short settle interval before the next teardown`) + the non-interactive exit obligation (`keep settling and retrying the teardown until it succeeds`); the Claude `## Terminal Team Teardown` now requires the inter-attempt settle (`wait for the settle before the next`, `sleep 2`) + the non-interactive loop obligation (`keep the settle-then-`, `until \`TeamDelete\` succeeds`). Gut-edit fingerprints (`immediately call teamdelete again`, `surface to the captain and end the turn`) added to the negating set. NB: dropped the over-broad `immediately re-fire` fingerprint — the CORRECT prose says "do NOT immediately re-fire", which would false-positive; the gut-edit's affirmative `immediately call TeamDelete again` is the real fingerprint.
- DONE: Make the oracle comprehensive over EVERY load-bearing teardown clause and PROVE the gut-edit goes RED on each.
  The oracle now pins all three clauses per surface (retry-to-success, settle-between-attempts, non-interactive cap-exhaustion exit). PROVEN by reproducing the audit's gut-edit (both files: settle → `immediately call TeamDelete again`; cap-exhaustion → `surface to the captain and end the turn`): oracle RED on `wait a short settle interval before the next teardown` + `keep settling and retrying the teardown until it succeeds` (shared-core) and on `wait for the settle before the next` + `sleep 2` + `keep the settle-then-` + `until \`TeamDelete\` succeeds` + the negating `surface to the captain and end the turn` (Claude). Cycle-1 inversion still turns it RED (regression-checked); correct prose stays GREEN. Restored both files (empty `git diff --stat skills/first-officer/` confirms no prose change).

### Summary
Cycle-2's one remaining hole is closed by an oracle extension only — the contract prose was already correct and is untouched (verified: the only changed file is `skills/integration/terminal_teardown_retry_test.go`). The delegation-sentence loophole is shut: `stripDelegationTail` removes the trailing `…are the adapter's` tail from shared-core step 10 so the load-bearing settle/non-interactive-exit phrases must appear in the BEHAVIORAL prose, not the delegation tail. The oracle is now comprehensive over all three teardown clauses on both surfaces, and the audit's gut-edit (immediate-retry + surface-and-end-turn) turns it RED on every clause — proven by reproducing the gut-edit, plus a regression check that the cycle-1 inversion still trips. Full offline suite green (`go test ./...` 836 passed); `go vet -tags live` clean. AC-1 (the live sonnet+opus run) remains the behavioral oracle for the fix and is the validation/CI-trigger step — it cannot run from this stage.
