---
title: Fix the #442 launcher signal-forwarding data race — start forwardHostSignals after cmd.Start()
status: validation
sprint: 0230-stable-finalization
score: 0.75
source: "pre-tag spot-audit of #441/#442 (the Haiku-shipped parallel members), 2026-06-30. Confirmed real + reproduced under -race; verified non-ship-blocker (benign on shipped arches) but a known -race defect on the front door. Gates the v0.23.0 stable cut per the captain (ship this, then tag)."
id: 3p0ccnj99jbjf6h938fsgvk8
started: 2026-06-30T06:33:25Z
worktree: .worktrees/spacedock-ensign-launcher-signal-forward-race-fix
---

`internal/cli/host_exec.go:294` spawns the `forwardHostSignals` goroutine BEFORE `cmd.Start()` at `host_exec.go:296`. The goroutine reads `cmd.Process` (`host_launch_unix.go:43`) while `cmd.Start()` writes it — no happens-before edge, so a SIGTERM/SIGHUP delivered to the launcher during the fork/exec startup window is an unsynchronized read/write (UB per the Go memory model). Reproduced under `-race`. This is net-new in #442 (the prior `syscall.Exec` model had no resident goroutine).

## Problem

The race is benign on the shipped arches (darwin/arm64, linux/amd64): `cmd.Process` is a pointer and an aligned pointer load/store is atomic, so no crash/corruption — the worst case is a dropped SIGTERM/SIGHUP forward in the few-ms fork/exec window (rare, and a session manager/supervisor escalates to SIGKILL past its grace timeout). But it is a shipped data race on the high-stakes front-door teardown path, and it carries a TEST-STRENGTH HOLE: every signal test (`host_launch_test.go:313`, `host_launch_pty_test.go:93`) fires its signal only AFTER the stub logs `STARTED` — long after `Start()` returns — so `go test ./... -race` stays GREEN while the race ships. This is exactly the coverage gap the README proof policy's adversarial-audit clause targets.

## Proposed approach (audit-specified)

Keep `signal.Notify` before `cmd.Start()` (the buffered channel already queues any early signal), but start the forwarding goroutine AFTER `cmd.Start()` returns, so the goroutine-creation happens-before edge publishes `cmd.Process` safely. This also closes the latent dropped-early-signal logic gap (an early signal now drains from the buffered chan once the goroutine starts).

## Acceptance criteria

- **AC-1 (VALUE)** — the data race is gone: a `-race` regression test that runs the production launch ordering **in-process under the `testing` harness** and drives the forwarding goroutine's `cmd.Process` read with a SIGTERM goes RED on the current `host_exec.go` goroutine-before-Start ordering and GREEN after the fix. Verified by: `go test ./internal/cli/... -race` with the new in-process launch test — RED pre-fix, GREEN post-fix (the divergeable proof the existing tests lack). See **Test design (RED-first)** below for the determinism mechanism — the divergence is HB-guaranteed, not timing-stress, and was proven by spike (RED pre-fix / GREEN under the simulated post-fix ordering).
- **AC-2** — no behavior regression: an early signal (delivered before the goroutine starts) is still forwarded to the host (drained from the buffered chan), and signals delivered after startup forward as today; the existing TestLaunch* signal tests stay green.
- **AC-3** — no other launcher-startup goroutine reads a field written by `cmd.Start()` before a happens-before edge (sweep the spawn path for the same antipattern).

## Code verification (audit claims confirmed against live `internal/cli`)

- `host_exec.go:294` `stop := forwardHostSignals(cmd)` spawns the forwarding goroutine BEFORE `host_exec.go:296` `if err := cmd.Start()` writes `cmd.Process`. Confirmed.
- `signal.Notify(ch, …SIGTERM, SIGHUP)` is `host_launch_unix.go:38`, before the `go func()` at `:39` and before `Start` — the buffered (cap-16) chan already queues an early signal. Confirmed; the fix keeps this placement.
- The goroutine reads the `cmd.Start()`-written `cmd.Process` at `host_launch_unix.go:43` (`if cmd.Process != nil`) and `:44` (`cmd.Process.Signal(sig)`, which dereferences it). Confirmed — these are the racy reads (the spike's race traces point exactly here).

### AC-3 antipattern-sweep scope

Swept `internal/cli` (non-test) for any launcher-startup goroutine reading a `cmd.Start()`-written field before a happens-before edge. The **only** production goroutine in the package is the one at `host_launch_unix.go:39`. The only `cmd.Start()`-written field read off the spawning goroutine is `cmd.Process` at `:43`/`:44`, inside that one goroutine. The other `cmd.Process`/`ProcessState` reads (`host_exec.go:299` `cmd.Wait()`, `:300` `cmd.ProcessState`) run on the **main** goroutine, sequenced after `Start()` returns — no concurrency, no race. `host_launch_other.go` (non-unix) spawns no goroutine. So AC-3's scope is exactly this single goroutine; there is no second instance of the antipattern to fix.

## Test design (RED-first)

The AC-1 regression test drives the **real `execHost.Launch` in-process** (in the test goroutine), so the forwarding goroutine runs under the `testing` harness — Go's `testing`+race integration then fails the test on any detected race. Shape (extends the existing re-exec helper harness in `host_launch_test.go`):

1. Build the stub argv (`self`, logpath, `wait`) and the `…ROLE=stub` env, exactly like `roleCmd`.
2. Spawn a goroutine that `waitForLog(STARTED)` then `syscall.Kill(os.Getpid(), SIGTERM)`. Gating the self-signal on `STARTED` is the safety + liveness lever: STARTED is logged by the stub only after the launcher armed `signal.Notify` and `Start` returned, so the SIGTERM is provably caught by the launcher's handler (never the test's default disposition) and provably drives the forwarding goroutine's `cmd.Process` read.
3. Call `execHost{}.Launch(argv, env)`. Pre-fix: the goroutine reads `cmd.Process` with no HB edge to `Start`'s write → the race detector flags it → `testing` fails the test (RED). The forwarded SIGTERM also reaches the stub, which exits 143, so `Launch` returns cleanly (a secondary AC-2 signal that the forward still works).

**Determinism mechanism (HB-guaranteed, not timing-stress).** The forwarding goroutine is created (pre-fix) BEFORE `cmd.Start()`, so its vector clock never includes the `Start` write — *every* read it performs of `cmd.Process` is unsynchronized with the write, independent of when the signal lands. The test therefore does not need to "overlap the fork/exec window"; it only needs the goroutine to perform the read once, which the STARTED-gated SIGTERM guarantees. **Post-fix GREEN is guaranteed** by the Go memory model's go-statement edge: spawning the goroutine AFTER the `Start` write (on the same main goroutine) makes write → `go` → read a happens-before chain, so the read can never race the write — regardless of signal timing.

**Why the existing tests stay GREEN (corrected diagnosis).** The filing premise — "the existing tests fire after `STARTED`, so `-race` stays green" — misdiagnoses the cause. Timing is irrelevant: a late in-process read still races (proven). The existing `TestLaunch*` signal tests are green because they run the launcher in a **separate child process** whose race report is *swallowed*: (a) `roleCmd` never sets the child's `Stderr`, so the child's `WARNING: DATA RACE` goes to `/dev/null`; (b) `launcherRole` calls `os.Exit(code)`, bypassing the race runtime's exit-code override; (c) no `testing` harness runs in the child to consult `race.Errors()`. Proven by spike: capturing the child's stderr and running it with `GORACE=halt_on_error=1` while firing the SIGTERM *late* (exactly as `host_launch_test.go:313` does) made the child exit **66** and emit `WARNING: DATA RACE` at `host_launch_unix.go:43` vs `:39`. The fix to this TEST-STRENGTH HOLE is to exercise the ordering **in-process**, not to fire earlier.

**Cost/complexity.** One Go `-race` unit test, ~25 lines, fixture-only (re-exec helper; no live workflow, no network). Reuses `newStubLog`/`roleCmd`/`waitForLog`/`withoutEnv`. Runs in well under a second.

## Test plan

- New in-process Go `-race` test as above, proven RED-first on the current ordering and GREEN under the simulated post-fix ordering (both demonstrated in the ideation spike).
- `go test ./internal/cli/... -race` green; full `go test ./... -race` green.
- The existing `TestLaunchResidentParent` / `TestLaunchExitCodePropagation` / `TestLaunchPTYTerminalSignals` stay green (AC-2 no-regression).
- High-stakes front-door surface: detached adversarial audit at validation (per the README proof policy).

## Doc-diff determination

No doc diff needed. This is an internal concurrency fix to the resident launcher's signal-forwarding startup ordering. It changes no user-visible behavior — the same signals are forwarded with the same exit-code mapping and the same CLI surface, banners, and host integration; only the read/write of `cmd.Process` becomes correctly synchronized. Nothing the docs site describes changes.

## Spike record

The riskiest unverified mechanism here was "a `-race` test that reliably diverges (RED pre-fix, GREEN post-fix)." Exercised first, in-process in package `cli`:

- **RED pre-fix:** replicating the live `forwardHostSignals(cmd)` → `cmd.Start()` ordering (and, separately, driving the real `execHost.Launch` in-process) fails under `go test ./internal/cli/... -race` with `WARNING: DATA RACE` at `host_launch_unix.go:43`/`:44` (read) vs `:39` (goroutine carrying the `Start` write). Reliable across early-injection, signal-storm, and late-injection variants — timing-independent.
- **GREEN post-fix (simulated):** the same injection with the goroutine spawned AFTER `cmd.Start()` passes clean — confirming the go-statement HB edge closes the race and the divergence is real.
- **Swallow proof:** the separate-process path races identically but is hidden (child exit 66 + `DATA RACE` banner only when stderr is captured under `GORACE=halt_on_error=1`).

This throwaway exercise (since removed) seeds implementation's first test: the in-process Design-A test above is the concrete RED-first artifact implementation writes and turns GREEN. No further unverified mechanism remains.

## Stage Report: ideation

- DONE: Design the RED-first -race test concretely … State the determinism mechanism … do not hand-wave "it should race".
  Concrete test in **Test design (RED-first)** (drive real `execHost.Launch` in-process, STARTED-gated self-SIGTERM); determinism is HB-guaranteed, proven by spike — RED pre-fix, GREEN under simulated post-fix. NOTE: the checklist's "force the read and Start write to OVERLAP during the fork/exec window" framing was disproven by the spike — timing is irrelevant in-process; the real lever is running under the `testing` harness (a late read races just as well). Delivered the stronger, proven mechanism instead.
- DONE: Verify the audit-specified approach against the live internal/cli code … name the AC-3 antipattern-sweep scope.
  **Code verification** section confirms `host_exec.go:294` goroutine-before-`:296` Start, `signal.Notify` at `host_launch_unix.go:38` before Start, the `cmd.Process` read at `:43`/`:44`. **AC-3 sweep** finds exactly one production goroutine (`host_launch_unix.go:39`); no second antipattern instance.
- DONE: Record in the task body the doc-diff determination and the spike record per the proof policy.
  **Doc-diff determination** (internal race fix, no user-visible change → no doc diff) and **Spike record** (re-ran the reproduction: RED pre-fix / GREEN post-fix / swallow proof; seeds implementation's first test) both in the body.

### Summary

Verified every audit claim against live `internal/cli` and designed the divergent RED-first `-race` test concretely. Key correction: the filing premise that the existing tests are green because they fire "after STARTED" is wrong — they are green because the launcher runs in a separate child process whose race report is swallowed (child stderr→/dev/null, `os.Exit` bypass, no `testing` harness); a late in-process read races identically. The fix and the test were both validated by a throwaway in-process spike (since removed): pre-fix RED, simulated-post-fix GREEN, plus a `GORACE=halt_on_error=1` capture proving the existing path races (child exit 66). AC-3 sweep is clean: one launcher-startup goroutine, no other instance. No doc diff; no remaining unverified mechanism.

## Stage Report: implementation

- DONE: Write the in-process RED-first -race test per the approved "Test design (RED-first)" … CONFIRM it goes RED on the CURRENT pre-fix ordering BEFORE you apply the fix — capture that RED run.
  `TestLaunchSignalForwardNoDataRace` (host_launch_test.go) drives real `execHost.Launch` in-process with a STARTED-gated self-SIGTERM, reusing newStubLog/roleCmd-shape argv/waitForLog/withoutEnv. Captured RED pre-fix: `--- FAIL … race detected`, `WARNING: DATA RACE` read at host_launch_unix.go:43/:44 vs `cmd.Start()` write at host_exec.go:296, goroutine created at host_launch_unix.go:39 — exactly the AC-1 antipattern.
- DONE: Apply the approved fix in host_exec.go: start the forwardHostSignals goroutine AFTER cmd.Start() (keep signal.Notify before Start) … Confirm GREEN and that `go test ./internal/cli/... -race` plus full `go test ./... -race` are green.
  Commit e9c82046: `forwardHostSignals` now returns `(forward, stop)`; `signal.Notify` still arms before Start, the launcher calls `forward()` after `cmd.Start()` so the go-statement HB edge orders Start's `cmd.Process` write → goroutine read. Non-unix shim + caller updated. New test GREEN; `go test ./internal/cli/... -race` ok; full `go test ./... -race` all packages ok; `gofmt` clean, `go vet` clean, `GOOS=windows go build` ok.
- DONE: Confirm AC-2 no-regression: existing TestLaunch* signal tests stay green, and an early signal delivered before the goroutine starts is still forwarded (drained from the buffered chan once the goroutine starts).
  TestLaunchResidentParent / TestLaunchExitCodePropagation / TestLaunchPTYTerminalSignals all PASS under -race. Early-signal-forward preserved: `signal.Notify` armed before Start (host_launch_unix.go) into a cap-16 buffered chan, and `forward()`'s `for sig := range ch` drains any queued signal once the pump starts — the same range-drain-and-forward path the new test exercises end-to-end (forwarded SIGTERM → stub exit 143). The HB-not-timing design has no externally injectable "after Start, before forward()" window, so no flaky timing test was added; the buffered-chan queue + exercised forward path is the proof.

### Summary

Closed the #442 launcher signal-forwarding data race by deferring the `forwardHostSignals` pump goroutine until after `cmd.Start()` returns, so the go-statement publishes the started `cmd.Process` to the goroutine with a happens-before edge. `signal.Notify` stays before Start (early signals queue on the cap-16 buffered chan and now drain correctly once the pump runs, closing the latent dropped-early-signal gap). Proved RED-first: the new in-process `TestLaunchSignalForwardNoDataRace` fails under `-race` on the pre-fix ordering (DATA RACE at host_launch_unix.go:43 vs host_exec.go:296) and passes after the fix. Full `go test ./... -race`, the existing TestLaunch* signal tests, `go vet`, `gofmt`, and a windows cross-build are all green. AC-3 needed no further change — ideation's sweep already confirmed this was the only launcher-startup goroutine reading a Start-written field, and the fix synchronizes exactly that read.
