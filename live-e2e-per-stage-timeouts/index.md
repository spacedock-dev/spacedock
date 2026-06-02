---
id: 3780rb1fcaw0dt8pytyk8vmx
title: Live-e2e per-stage timeouts — port the FOStreamWatcher pattern, ban the monolithic long timeout
status: backlog
source: "captain (2026-06-02): live-e2e flakiness dive — v1 regressed to a single 12m monolithic timeout; restore the Python harness's per-stage timeouts (no individual timeout > 60s, long timeout banned) and verify reliable e2e"
started:
completed:
verdict:
score: "0.36"
worktree:
issue:
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

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — Per-stage timeouts; no monolithic long timeout.** No individual timeout exceeds 60s; the 12m `liveTimeout` monolithic ctx is removed; each stage event has its own ≤60s budget.
Verified by: the test source carries no timeout > 60s; an injected hung/no-progress stage fails in ≤~60s localized to that stage (not at a 10m cap).

**AC-2 — Streaming transcript (diagnosable on hang).** The FO's stream-json streams to the test log as it arrives; a hang leaves a partial transcript naming the last stage reached.
Verified by: a simulated no-progress hang yields a non-empty transcript identifying the stalled stage (vs today's zero-output-on-hang).

**AC-3 — Reliable e2e established.** The live cycle drives FO→done reliably under the per-stage watcher, OR any failure is fast + localized + diagnosable (no silent 10m panic).
Verified by: multiple CI-E2E runs on this entity's PR (sonnet + opus) — consistent green, or a captured localized transcript if it flakes (feeding the runtime-follow-up decision).

## Notes
- `internal/ensigncycle` (live_test.go) + a Go stream-watcher helper (port of `~/git/spacedock/scripts/test_lib.py` FOStreamWatcher: `expect`/`expect_dispatch_close`/`expect_exit`, each with `timeout_s`). NOT the status lane.
- Live verification needs CI-E2E (the local box's `~/.claude/benchmark-token` is stale → 401), so the binding reliability proof is this entity's PR CI-E2E, as with the 0.19.3 fixes.
- Subsumes the 0.19.3-dive cheap fixes (timeout config, transcript streaming, mode-agnostic). The FO-runtime headless-team-await reliability question is DEFERRED pending the transcript this fix captures — do not file it speculatively.
