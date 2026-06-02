---
id: 3780rb1fcaw0dt8pytyk8vmx
title: Live-e2e per-stage timeouts — port the FOStreamWatcher pattern, ban the monolithic long timeout
status: validation
source: "captain (2026-06-02): live-e2e flakiness dive — v1 regressed to a single 12m monolithic timeout; restore the Python harness's per-stage timeouts (no individual timeout > 60s, long timeout banned) and verify reliable e2e"
started: 2026-06-02T07:02:10Z
completed:
verdict:
score: "0.36"
worktree: .worktrees/spacedock-ensign-live-e2e-per-stage-timeouts
issue:
mod-block: merge:pr-merge
---

`TestLiveEnsignCycle` (internal/ensigncycle/live_test.go) drives the entire FO→done cycle through one headless `spacedock claude -p` process under a SINGLE `liveTimeout = 12*time.Minute` ctx + `cmd.CombinedOutput()`. This regressed from the upstream Python harness (`~/git/spacedock/scripts/test_lib.py`, `FOStreamWatcher` ~L1288), which bounds EACH stage/event with its own per-stage timeout (`w.expect(event, timeout_s=…)`; upstream `PER_STAGE_OVERALL_S` 60–120s, `SUBPROCESS_EXIT_BUDGET_S` 180s).

Consequences (from the 0.19.3 flakiness dive): a hung stage isn't caught until the monolithic cap — and because CI runs `go test` with no `-timeout` (.github/workflows/runtime-live-e2e.yml:148), Go's default 10m fires before the unreachable 12m, yielding a silent `panic` instead of a clean message; `CombinedOutput()` buffers the whole transcript and returns it only after the child exits, so a hang yields ZERO output (undiagnosable); and there is no localization to the stalled stage.

**Captain directive: no individual timeout > 60s; the long monolithic timeout is banned. Port the per-stage stream-watcher.**

## Scope
- Replace the monolithic `liveTimeout` ctx + `CombinedOutput()` with a STREAMING per-stage watcher (a Go port of `FOStreamWatcher`): consume the FO's stream-json as it arrives and `expect` each stage-progress event — TeamCreate, dispatch, per-ensign completion signal, terminalize, archive, process-exit — each bounded by its own per-stage timeout ≤60s. A stage that doesn't progress within its budget fails FAST, localized to that stage, with the streamed transcript intact.
- Subsumes the dive's cheap fixes: the per-stage bound replaces the monolithic timeout; streaming fixes the transcript-blindness; the watcher is dispatch-mode-agnostic (catches a hang under team OR bare mode) — so no team-vs-bare change is needed.
- Live stages legitimately take minutes, so the ≤60s budget is a NO-PROGRESS / quiet bound (≤60s with no new stream activity → hung), not a total-stage cap. Ideation fixes the exact semantics from the FOStreamWatcher reference.

## Verify reliable e2e (captain)
After the fix, ESTABLISH reliability: run the live cycle multiple times (this entity's PR CI-E2E, sonnet + opus) and confirm consistent green OR a fast + localized + diagnosable failure. If the FO genuinely flakes (doesn't drive to done) even with bounded detection, capture the streamed transcript and root-cause it — that may surface a separate FO-runtime headless-await follow-up (deferred until this fix produces that evidence).

## Design: Go port of FOStreamWatcher

### What the watcher is

A small, default-build-tag Go helper (`streamWatcher`) that consumes the FO's
stream-json JSONL **as it arrives** from `claude`'s stdout and bounds each
stage-progress step with its own ≤60s budget, replacing the monolithic
`liveTimeout` ctx + `cmd.CombinedOutput()`. It is the direct port of
`~/git/spacedock/scripts/test_lib.py` `FOStreamWatcher` (`expect` ~L1586,
`expect_dispatch_close` ~L1628, `expect_exit` ~L1705), specialised to the v1
live cycle. It compiles under DEFAULT tags so an offline unit test can exercise
it against synthetic logs with no model spend (the same split the package
already uses for `locateEntity`/`someCommitNamesOnly` in `liveassert_test.go` +
`liveassert_unit_test.go`); the `//go:build live` `live_test.go` wires it to the
real subprocess.

### Wiring change in live_test.go (the riskiest plumbing)

Today: `cmd.CombinedOutput()` under a single 12m ctx. New: launch the SAME argv
under a **plain `context.Background()`** (no deadline ctx — the per-step budgets
ARE the timeout discipline; AC-1 bans the monolithic ctx), take `cmd.StdoutPipe()`
with `cmd.Stderr = cmd.Stdout`-equivalent (`cmd.Stderr = w` where `w` is the
same pipe-tee target), `cmd.Start()`, and feed the pipe to the watcher. The
watcher tees every line it reads to the test log (`t.Log`) immediately, so a hang
leaves a partial transcript naming the last step reached (AC-2). Post-watch, the
existing `locateEntity` / frontmatter / `someCommitNamesOnly` end-state
assertions run unchanged against the completed-and-archived entity.

**Spike-confirmed mechanism (see Spike below):** `spacedock claude` is a
`syscall.Exec` front door (`internal/cli/frontdoor.go:138` `ops.Launch` →
`internal/cli/host_exec.go:204` `syscall.Exec`) that REPLACES the process with
`claude`, forwarding `fd.passthrough...` verbatim (`frontdoor.go:124`). So
`--output-format stream-json --verbose` reach the real `claude`, and `claude`
writes stream-json straight to the inherited stdout pipe. The event shapes are
therefore the SAME standard Claude Code stream-json the upstream Python watcher
already parses — there is no v1-specific event-shape risk.

### Steps watched, and what is NOT a stream event

The proven upstream sequence (`tests/test_dispatch_completion_signal.py:173-192`,
`tests/test_feedback_keepalive.py:143-181`) watches exactly THREE kinds of step,
then asserts end-state from the FILESYSTEM after exit:

1. `expect(isTeamCreate, label="TeamCreate")` — the FO's `TeamCreate` assistant
   `tool_use` (teams mode engaged). The live cycle runs teams mode, so this is
   the first progress beat.
2. `expectDispatchClose(...)` per ensign — open on an `Agent(subagent_type="spacedock:ensign")`
   assistant `tool_use`, close on the matching completion signal (the
   mode-agnostic close anchors the Python port already proved: bare-mode synchronous
   `Agent()` tool_result, teams-mode `task_notification status=completed`, or the
   headless inbox-poll `from: spacedock-ensign-…` + `text: Done:` Bash tool_result).
   The flat-entity fixture drives backlog→done, which is ONE ensign dispatch that
   writes `## Stage Report: done` — so the live cycle expects exactly one
   dispatch-close.
3. `expectExit(...)` — the FO subprocess exits.

**Decision — terminalize + archive are POST-EXIT filesystem state, not stream
events.** The dispatch checklist item 1 lists "terminalize → archive →
process-exit" among the events to `expect`; the proven upstream pattern does NOT
watch terminalize/archive on the stream — it watches dispatch-close + exit, then
asserts the archived entity + `status: done` + path-scoped commit from disk/git
(exactly v1's existing `locateEntity`/`frontmatterField`/`someCommitNamesOnly`).
Forcing the watcher to match a synthetic "terminalize"/"archive" stream event
would invent a brittle new event shape the FO does not deterministically emit on
the stream, and would re-introduce a guess-the-mechanism risk. So this fix
**keeps the three watched stream steps above** and **leaves terminalize+archive as
the unchanged post-exit FS assertions**. This is faithful to the FOStreamWatcher
reference and to v1's current end-state checks; it tightens the dispatch's intent
(catch a hang fast + localized) without modelling runtime internals.

### Budget semantics — NO-PROGRESS (quiet) bound, not total-stage

This is the exact semantic the dispatch asked ideation to nail down, and it is
where the Go port **intentionally differs** from a literal copy of upstream
`expect()`:

- Upstream `expect(timeout_s)` is a per-CALL **overall** budget: the deadline is
  `now + timeout_s` at the start of the call and does NOT reset as new stream
  lines arrive; upstream sets it to 120s (`PER_STAGE_OVERALL_S`) precisely
  because a live stage legitimately takes minutes of wallclock and it needs
  headroom.
- The captain bans any individual timeout > 60s. A 60s OVERALL per-stage budget
  would false-fail a healthy stage that simply takes >60s of real work. The
  reconciliation, stated in Scope: the ≤60s budget is a **NO-PROGRESS / quiet**
  bound — the deadline resets to `now + 60s` every time the watcher drains ANY
  new stream line (a heartbeat of FO activity), and trips only after 60s of
  stream SILENCE. A live stage that progresses (assistant turns, tool_use,
  tool_result all count as activity) never trips it no matter how long it runs;
  a genuinely hung stage emits no new lines and trips in ≤~60s, localized to the
  step's `label`.

  Concretely the Go `expect` loop, each poll: drain new lines → if any new line
  arrived, `deadline = now + quietBudget`; if predicate matched, return; if proc
  exited, do a final drain then fail with `StepFailure(label, exitCode, tail)`; if
  `now >= deadline`, fail with `StepTimeout(label, tail)`. This is the one-line
  delta from upstream (`reset deadline on any drained line`) that turns upstream's
  overall budget into the captain's quiet budget.

- Budget constants (all ≤60s per the directive):
  - `quietBudgetS = 60` — the per-step no-progress bound for `expect` and
    `expectDispatchClose`.
  - `exitBudgetS = 60` — `expectExit` waits up to 60s for the FO to exit AFTER
    the last watched step matched; on timeout it `proc.Kill()`s and fails
    `StepTimeout("expect_exit")` with the transcript tail. (Upstream uses 180s
    here; the directive caps it at 60 — acceptable because by `expect_exit` the
    contract-bearing work is already done and the FO only needs to wind down.)
  - `pollIntervalMs = 200` — matches upstream `POLL_INTERVAL_S = 0.2`.

  NO constant in the test exceeds 60s, and there is NO monolithic deadline ctx —
  AC-1 enforced both by the source-scan test and the injected-hang test.

### Error surface (port of StepTimeout / StepFailure)

Two failure shapes, each carrying the step `label` + a bounded transcript tail
(the last ~8000 bytes / ~20 lines, reusing the existing `tail` helper) so the
red is self-diagnosing:

- `StepTimeout{label}` — no progress within the quiet budget (or no exit within
  `exitBudgetS`). The killed/abandoned subprocess + the tail name the stalled step.
- `StepFailure{label, exitCode}` — the FO subprocess exited BEFORE the expected
  step matched (e.g. crashed mid-cycle). Carries the non-zero exit code + tail.

On timeout the watcher `proc.Kill()`s the subprocess so a hung `claude` does not
outlive the test. Both shapes `t.Fatalf` with the label + tail already streamed.

## Spike (riskiest unverified mechanism, run first)

**Question that would invalidate the rest of this work if it broke:** does the v1
`spacedock claude` front door actually emit parseable stream-json JSONL to a pipe
the test can read line-by-line (rather than buffering, munging, or swallowing it),
and are the event shapes the same ones the upstream Python watcher matches?

**Done — confirmed by reading the front-door path, no model spend needed.**
`spacedock claude` parses front-door flags then `ops.Launch(argv)` →
`syscall.Exec` (`internal/cli/frontdoor.go:97-142`, `internal/cli/host_exec.go:197-205`),
which REPLACES the spacedock process image with `claude` — spacedock does not sit
between `claude`'s stdout and the test's pipe, does not buffer, and forwards
`--output-format stream-json --verbose` verbatim via `fd.passthrough`. Therefore
the bytes on the test's `StdoutPipe()` are `claude`'s own stream-json, identical
in shape to the logs the upstream `FOStreamWatcher` parses (it parses the same
`claude` binary's stream-json). The line-buffering + partial-line-hold logic the
Go port needs is the proven upstream `_drain_entries` algorithm. Residual unknown:
the LIVE FO actually driving backlog→done reliably under a real model — that is
not a parser risk, it is the reliability question AC-3 establishes via CI-E2E.

## Acceptance criteria

**AC-1 — Per-stage no-progress timeouts; no monolithic long timeout.** The live
test carries NO individual timeout > 60s and NO monolithic deadline ctx: the 12m
`liveTimeout` const and its `context.WithTimeout` are gone, the subprocess runs
under a plain `context.Background()`, and each watched step is bounded by its own
≤60s quiet (no-progress) budget that resets on stream activity.
Verified by: (a) a Go test asserting every timeout literal/const reachable from
the live test is ≤60s and that no `context.WithTimeout`/`liveTimeout` remains in
`live_test.go` — parses real values, not a substring grep for spelling; (b) an
offline unit test feeding the `streamWatcher` a synthetic log that goes silent and
asserting `expect` returns a `StepTimeout` carrying the step's label within ~1×
the (test-shrunk) quiet budget — proving the no-progress trip is localized, not
total-stage; (c) an offline unit test feeding a log that keeps emitting unrelated
lines past the quiet budget and asserting `expect` does NOT trip (the reset-on-
activity semantic) until the predicate matches.

**AC-2 — Streaming transcript, diagnosable on hang.** The watcher tees each
stream line to the test log as it is drained, so a hang leaves a non-empty partial
transcript whose tail names the last step reached — versus today's
`CombinedOutput()` which returns ZERO bytes until the child exits.
Verified by: an offline unit test where the synthetic log advances through
`TeamCreate` then goes silent before the dispatch-close; the `StepTimeout` error
message and the streamed log both name the stalled step (`"dispatch close"`) and
include the `TeamCreate` line — asserting a non-empty, step-naming tail.

**AC-3 — Reliable e2e established (LIVE).** The live cycle drives FO→done under
the per-stage watcher, OR any failure is fast + localized + diagnosable (a labelled
`StepTimeout`/`StepFailure` with a partial transcript — never a silent Go 10m
default-timeout panic).
Verified by: multiple CI-E2E runs on THIS entity's PR across the existing matrix
(sonnet on CI-E2E + claude-opus-4-8 on CI-E2E-OPUS, `.github/workflows/runtime-live-e2e.yml`)
— the binding proof is consistent green across runs, OR, if the FO genuinely
flakes (does not drive to done) even with bounded detection, the captured labelled
transcript + step is attached and root-caused, feeding the DEFERRED FO-runtime
headless-await follow-up decision. (The local box's `~/.claude/benchmark-token` is
stale → 401, so local live runs cannot bind this; CI-E2E is the proof surface, as
with the 0.19.3 fixes.)

## Test plan

| Layer | What it proves | Cost | Fixture/CLI/live |
|---|---|---|---|
| Offline unit (`streamwatch_unit_test.go`, default tags) | `expect` returns matched entry; `StepTimeout` on quiet-silence carries label; reset-on-activity does NOT trip on a noisy-but-unmatched stream; `StepFailure` on early proc exit carries exit code + label; final-drain-before-poll ordering; partial-line hold across reads; dispatch-close open/close on the three mode anchors | seconds, no model | synthetic JSONL + a fake exit-poller (port of Python `_FakeProc` + `test_fo_stream_watcher.py`) |
| Offline static (`streamwatch_unit_test.go` or `live_budget_test.go`, default tags) | AC-1: no timeout literal/const reachable from the live path > 60s; no `liveTimeout`/`context.WithTimeout` in `live_test.go` | seconds, no model | parse the test source's real timeout values (a relationship over real values, not a spelling grep) |
| Live (`live_test.go`, `//go:build live`) | AC-3: the watched-step sequence (TeamCreate → one dispatch-close → exit) completes, then the unchanged end-state assertions pass; on a real hang it fails fast + labelled | ~350s/run, model spend, gated | CI-E2E PR run (sonnet + opus matrix) |

**Riskiest path first (mechanism check before comprehensive run):** the offline
unit tests for `expect`'s quiet-budget trip + reset-on-activity are written and
green BEFORE the live wiring is touched — they validate the one semantic that
differs from upstream (no-progress reset). The live CI-E2E run is the comprehensive
run that follows, per the "smallest end-to-end exercise of the riskiest path
first" discipline.

**Cost/complexity:** small. The watcher is ~150 lines (port of a proven Python
class), the live-test change is swapping `CombinedOutput()` for `StdoutPipe()` +
the watch sequence, and the end-state assertions are untouched. The offline unit
tests mirror the existing `liveassert_unit_test.go` split, so the package's
default-tags-cover-the-helpers convention is preserved.

## Scope confirmation & subsumption

- **In scope:** `internal/ensigncycle/live_test.go` (swap monolithic ctx +
  `CombinedOutput()` for the streaming per-step watcher) + a new
  default-build-tag `streamWatcher` helper and its offline unit test. The watcher
  is **dispatch-mode-agnostic** (the close anchors cover bare AND teams mode), so
  NO team-vs-bare change is needed.
- **Subsumes the 0.19.3-dive cheap fixes:** the per-step quiet bound replaces the
  monolithic 12m timeout (and the Go default-10m silent panic it caused);
  streaming the pipe to `t.Log` fixes the `CombinedOutput()` transcript-blindness;
  the mode-agnostic close anchor catches a hang under either mode. One mechanism
  subsumes all three.
- **NOT in scope (untouched):** the FO runtime (`internal/cli/frontdoor.go`,
  `internal/status/*`), the serialized internal/status lane, and the `agents/` /
  `references/` scaffolding. This is a TEST/HARNESS change only.
- **DEFERRED (do not file speculatively):** the FO-runtime headless-team-await
  reliability question. This fix is what PRODUCES the evidence — a captured,
  labelled, partial transcript of a genuine FO-drive flake — that would justify
  opening that follow-up. Until AC-3's CI-E2E runs surface such a transcript, the
  follow-up stays unfiled.

## Notes
- Reference: `~/git/spacedock/scripts/test_lib.py` `FOStreamWatcher` (class ~L1288,
  `expect` ~L1586, `expect_dispatch_close` ~L1628, `expect_exit` ~L1705); budget
  constants `PER_STAGE_OVERALL_S = 120` / `SUBPROCESS_EXIT_BUDGET_S = 180` in
  `tests/test_dispatch_completion_signal.py` (the v1 port caps both at ≤60s per the
  captain directive and reinterprets per-stage as no-progress); offline unit
  reference `tests/test_fo_stream_watcher.py`.
- Live verification needs CI-E2E (the local box's `~/.claude/benchmark-token` is
  stale → 401), so the binding reliability proof is this entity's PR CI-E2E.

## Stage Report: ideation

- DONE: Design the Go port of FOStreamWatcher (expect/expect_dispatch_close/expect_exit, each with per-stage timeout_s); name the events + decide per-stage budget semantics
  Design section added: watches TeamCreate → one dispatch-close → exit; budget reinterpreted as a ≤60s NO-PROGRESS quiet bound (deadline resets on any stream activity) — the one-line delta from upstream's overall budget; quietBudgetS=60, exitBudgetS=60, pollIntervalMs=200; StepTimeout/StepFailure carry label + transcript tail.
- DONE: AC-3 reliability-verification plan
  AC-3 binds reliability to multiple CI-E2E runs on this PR (sonnet on CI-E2E + opus on CI-E2E-OPUS via the existing matrix); consistent green OR a captured labelled transcript on a genuine FO flake, feeding the DEFERRED FO-runtime-await decision; local box can't bind (stale benchmark-token → 401).
- DONE: Confirm scope + subsumption
  Scope-confirmation section: live_test.go + a default-tag streamWatcher helper; mode-agnostic close anchors (no team-vs-bare change); subsumes the dive's three cheap fixes (monolithic timeout, CombinedOutput transcript-blindness, mode); FO runtime + status lane + scaffolding untouched; FO-runtime headless-await stays DEFERRED.

### Summary

Hardened the provisional spec into a complete ideation: a default-build-tag Go `streamWatcher` (port of the upstream Python `FOStreamWatcher`) replaces the 12m monolithic ctx + `CombinedOutput()` with a streaming, per-step ≤60s quiet-budget watch over the FO's stream-json. Two load-bearing design decisions resolved against the reference + the captain directive: (1) per-stage budget is a NO-PROGRESS bound (deadline resets on any drained line) not a total-stage cap, reconciling "≤60s" with stages that legitimately take minutes — the single intentional delta from upstream's 120s overall budget; (2) terminalize + archive are NOT stream events — upstream and v1 both assert them post-exit from filesystem/git, so the watcher watches only TeamCreate → dispatch-close → exit and the existing end-state assertions are left untouched (avoids inventing a brittle terminalize event shape). The riskiest mechanism (does `spacedock claude` emit parseable stream-json to a pipe?) was spiked by reading the `syscall.Exec` front door — confirmed it forwards `--output-format stream-json` verbatim and replaces the process, so the pipe carries `claude`'s own stream-json identical to what the Python watcher already parses; no model spend needed for the spike. Test plan layers offline unit (synthetic JSONL + fake exit-poller, no model) for the quiet-trip/reset-on-activity semantics + a static ≤60s/no-monolithic-ctx check, with CI-E2E as the live reliability proof.

## Stage Report: implementation

- DONE: RISKIEST-FIRST (TDD) — offline unit tests (default tags, synthetic JSONL + a fake exit-poller, port of `_FakeProc`) for the watcher's quiet-budget trip, reset-on-activity, and StepFailure-on-early-exit, GREEN before any live wiring
  `internal/ensigncycle/streamwatch_unit_test.go` (commit 00126297): written first, watched fail on undefined `streamWatcher`/`streamEntry`/`stepTimeout`, then GREEN. `TestExpectResetsDeadlineOnActivity` runs ~4× the (shrunk 150ms) quiet budget pushing unrelated lines and does NOT trip — proving no-progress reset, not a total-stage cap; `TestExpectQuietTimeoutCarriesLabel` trips at ~1× on silence carrying the step label; `TestExpectStepFailureOnEarlyExit` carries exit code + label. 14/14 pass, race-clean (`go test -race`).
- DONE: Implement the default-tag `streamWatcher` (port FOStreamWatcher: expect / expectDispatchClose / expectExit; quietBudgetS=60, exitBudgetS=60, pollIntervalMs=200; deadline resets on any drained line; StepTimeout/StepFailure carry label + transcript tail; proc.Kill() on timeout) AND swap live_test.go's `CombinedOutput()` + 12m `liveTimeout` ctx for `StdoutPipe()` + `context.Background()` + the watch sequence (TeamCreate → one dispatch-close → exit); LEAVE post-exit end-state assertions UNCHANGED
  `internal/ensigncycle/streamwatch.go` + `live_test.go` (commit 3e7a5621): `liveTimeout` const + `context.WithTimeout` removed; subprocess runs under plain `context.Background()`; stdout+stderr feed one `io.Pipe`, a background `cmd.Wait` closes the write-end so the scanner drains final lines, `cmdPoller` surfaces exit non-blockingly. Three mode-agnostic close anchors ported (bare Done tool_result, teams `task_notification` completed, headless inbox-poll `from:`+`Done:`). `locateEntity`/`frontmatterField`/`someCommitNamesOnly` end-state checks untouched (terminalize+archive stay post-exit FS assertions).
- DONE: AC-1 static guard + suite green — a default-tag test asserting no timeout literal/const reachable from the live path exceeds 60s AND no `liveTimeout`/`context.WithTimeout` remains in live_test.go (parse real values, not a spelling grep); `go test ./internal/ensigncycle/` + `go vet ./...` green; `go vet -tags live` compiles
  `internal/ensigncycle/live_budget_test.go` (commit 3e7a5621): parses the live-path AST and evaluates the REAL `time.Duration` each `N * time.Unit` literal denotes. PROVEN by injecting `120 * time.Second` + a `context.WithTimeout` call → both guards went RED (caught `2m0s`, the string "60" never appearing), then reverted. `go test ./internal/ensigncycle/` ok, `go test -race` ok, `go vet ./...` clean, `go vet -tags live ./internal/ensigncycle/` clean, full-repo `go test ./...` green.

### Summary

Replaced the live e2e test's monolithic 12m `liveTimeout` ctx + `CombinedOutput()` with a streaming `streamWatcher` (Go port of upstream `FOStreamWatcher`) that bounds each watched step — TeamCreate → the one backlog→done dispatch-close → FO exit — by its own ≤60s NO-PROGRESS quiet budget (deadline resets on any drained stream line), tees every line to `t.Log` so a hang leaves a partial step-naming transcript, and kills a hung subprocess on a labelled `stepTimeout`/`stepFailure`. The one intentional delta from upstream — quiet (no-progress) bound vs overall budget — is the riskiest semantic and was the first thing proven, TDD-first, by an offline unit test with synthetic JSONL + a fake exit-poller (race-clean). Terminalize+archive stay as the unchanged post-exit filesystem assertions, faithful to the ideation decision. The AC-1 guard parses real durations from the AST rather than grepping the spelling "60", and was confirmed to go red on injected violations before reverting. AC-3 (live reliability) is NOT establishable locally — the box's `~/.claude/benchmark-token` is stale → 401 — so the binding proof is CI-E2E on this PR (sonnet + opus matrix); validation establishes it. CODE committed on branch `spacedock-ensign/live-e2e-per-stage-timeouts` (commits 00126297, 3e7a5621).

## Stage Report: validation

- DONE: Reproduce AC-1: confirm no monolithic deadline ctx and no live-path timeout > 60s — re-run live_budget_test (the AST duration-eval guard) AND independently confirm it goes RED on an injected `120 * time.Second` / `context.WithTimeout` then reverts; confirm streamWatcher's quiet budget trips fast + localized on stream silence (the offline unit, shrunk budget).
  Independently injected `_ = 120 * time.Second` + `context.WithTimeout(...)` into live_test.go → BOTH guards went RED (TestNoTimeoutLiteralExceeds60s caught the real evaluated `2m0s`, NOT a "60" spelling grep; TestLiveTestHasNoMonolithicDeadlineCtx caught the WithTimeout); reverted → both GREEN, `git diff` clean. Only live-path durations are 60s/60s/200ms (streamwatch.go:37,40,42). Diff vs merge-base confirms `liveTimeout` const + `context.WithTimeout` + `CombinedOutput()` were REMOVED (not commented). TestExpectQuietTimeoutCarriesLabel trips at 0.15s = 1× the 150ms shrunk budget carrying label "dispatch close" (localized, not total-stage).
- DONE: Reproduce AC-2: the streaming/diagnosable-on-hang property — the offline unit where the synthetic log advances through TeamCreate then goes silent yields a NON-EMPTY partial transcript naming the stalled step (vs the old CombinedOutput zero-output-on-hang). Confirm lines are teed to t.Log as drained.
  TestExpectFailureCarriesTranscriptTail: StepTimeout names "dispatch close" AND its tail includes the prior `TeamCreate` line. Authored + ran a disposable probe (removed after) capturing the tee callback: it fired on EVERY drained line (2 lines captured BEFORE the hang tripped) — a non-empty step-naming partial transcript exists at hang time, the property CombinedOutput's zero-byte block-until-exit lacks. drainEntries tees each line before recording (streamwatch.go:127); live test passes `t.Log` as the tee (live_test.go:120).
- DONE: AC-3 wiring + scope + suite: confirm the live wiring is sound (StdoutPipe + plain context.Background() + the TeamCreate → one dispatch-close → exit watch; the 3 mode-agnostic close anchors; proc.Kill on timeout; terminalize/archive LEFT as post-exit FS assertions). Full `go test ./...` + `go test -race ./internal/ensigncycle/` + `go vet ./...` + `go vet -tags live ./internal/ensigncycle/` green; scope is internal/ensigncycle ONLY (no FO-runtime/status-lane touch). Note AC-3 LIVE reliability binds at the PR's CI-E2E.
  Wiring sound: subprocess launched via `exec.Command` (NO deadline ctx — no CommandContext/WithTimeout), stdout+stderr both into one `io.Pipe` (live_test.go:107-109 — realizes the spec's "same pipe-tee target" intent; the merged-stream case `cmd.StdoutPipe()` cannot cover); watch = TeamCreate→one dispatch-close→exit (lines 128/131/134); 3 mode anchors all pass (TestExpectDispatchCloseOnThreeAnchors 3/3); proc.Kill via cmdPoller.kill→Process.Kill (line 290) on expectExit timeout; terminalize/archive are post-exit FS assertions (locateEntity/frontmatterField/verdictSet/someCommitNamesOnly, lines 147-191). Suite: `go test ./...` exit 0 / 0 FAIL; `go test -race ./internal/ensigncycle/` exit 0 / 0 DATA RACE (41 pass); `go vet ./...` clean; `go vet -tags live ./internal/ensigncycle/` clean; live-tag package compiles. Scope: 4 files, ALL internal/ensigncycle/ (streamwatch.go, streamwatch_unit_test.go, live_budget_test.go new + live_test.go modified) — no frontdoor/status/agents/references touch. AC-3 LIVE reliability NOT bindable locally (stale benchmark-token → 401); CI workflow (.github/workflows/runtime-live-e2e.yml:148) runs `go test -tags live` with NO `-timeout` over the sonnet(CI-E2E)+opus(CI-E2E-OPUS) matrix — the per-step ≤60s quiet budgets now trip + Kill well before Go's default 10m, so the no-`-timeout` invocation no longer yields the old silent panic. The wiring I validated is exactly what CI runs.

### Summary

PASSED. All three ACs reproduced by exercising, not re-reading. AC-1's static guard is the load-bearing claim and I independently confirmed it: injecting `120 * time.Second` + `context.WithTimeout` turned both guards RED (catching the real evaluated `2m0s`, proving it parses durations not the string "60"), and reverting returned them GREEN — and the only live-path timeouts are 60s/60s/200ms. AC-2's streaming-diagnosability holds: the existing unit shows a hang leaves a tail naming the stalled step + the prior TeamCreate line, and a disposable tee probe confirmed every drained line is teed before the hang trips (non-empty partial transcript vs CombinedOutput's zero bytes). AC-3 wiring is sound (no deadline ctx, merged io.Pipe, TeamCreate→one dispatch-close→exit, 3 mode anchors, proc.Kill, post-exit FS assertions for terminalize/archive); full suite + race + both vet variants green; scope confined to internal/ensigncycle (4 files, no FO-runtime/status-lane). AC-3 LIVE reliability is correctly NOT bindable locally (stale benchmark-token → 401) and binds at this PR's CI-E2E — and the CI invocation (`-tags live`, no `-timeout`) is consistent with the fix: the per-step quiet budgets trip + Kill long before Go's default 10m, eliminating the old silent panic this entity targets.

## Stage Report: implementation (cycle 2 — validation touch-ups)

Validation PASSED + detached adversarial audit confirmed the fix correct (no blocking defect). Two audit-surfaced touch-ups applied before the PR:

- DONE: [BUG fix] Leaked child on an EARLY trip — a TeamCreate-stall or dispatch-close-stall did `t.Fatalf` WITHOUT killing the subprocess (only `expectExit` killed), orphaning a token-spending `claude`
  Added `defer poller.kill()` in `TestLiveEnsignCycle` right after the poller is created, so a hang at ANY step kills the child (no-op once exited). Offline guard `TestEarlyStepTripCanKillProc` asserts an early-step `stepTimeout` does NOT self-kill (killing is the caller's deferred job) and that the still-running fake proc remains killable — the unit proof of the `defer` contract (the live defer runs only under `//go:build live`). 15/15 offline tests pass, race-clean.
- DONE: [DOC known-gap, NO long timeout] The "stuck-but-emitting" (chatty) hang is not caught by the quiet budget (by design — the deadline resets on any drained line)
  Added explicit prose in `live_test.go` at the watch sequence: the quiet budget catches a SILENT hang (≤60s, localized) but a chatty hang falls to Go's BUILT-IN default test timeout — a conscious, documented trade-off (the captain banned long individual timeouts; chatty hangs are far rarer than silent ones). No explicit `-timeout` added.
- SKIPPED: AC-1 guard's exotic bypasses (`time.Duration(120)*time.Second` / raw-nanosecond literals)
  Out of scope per FO; committed code is clean (`60 * time.Second`). KNOWN MINOR GUARD-COVERAGE GAP: the AST guard evaluates `N * time.Unit` literals but does not catch a duration spelled via a `time.Duration(N)` conversion or a raw-nanosecond integer. Acceptable — the guard's job is to red-flag the obvious regression (a re-introduced large `N * time.Second`/`liveTimeout`/`context.WithTimeout`), which it does and is proven to do.

### Summary

Applied the two validation touch-ups on the code worktree branch: a real bug fix (`defer poller.kill()` so an early-step `t.Fatalf` no longer orphans a token-spending `claude`, guarded by a new offline unit test) and an explicit prose note documenting the conscious chatty-hang known-gap (no long timeout added — left to Go's built-in default per the captain's directive). Re-ran the offline suite (15/15 pass), `-race` (clean), `go vet -tags live` (clean), gofmt (clean), and the AC-1 guards (still green). The AC-1 guard's exotic-bypass gap is noted as a known minor limitation per FO scope. CODE committed on `spacedock-ensign/live-e2e-per-stage-timeouts`; NOT pushed (pr-mod pushes at the merge boundary).
