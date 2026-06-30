---
title: Fix the #442 launcher signal-forwarding data race — start forwardHostSignals after cmd.Start()
status: ideation
sprint: 0230-stable-finalization
score: 0.75
source: "pre-tag spot-audit of #441/#442 (the Haiku-shipped parallel members), 2026-06-30. Confirmed real + reproduced under -race; verified non-ship-blocker (benign on shipped arches) but a known -race defect on the front door. Gates the v0.23.0 stable cut per the captain (ship this, then tag)."
id: 3p0ccnj99jbjf6h938fsgvk8
started: 2026-06-30T06:33:25Z
---

`internal/cli/host_exec.go:294` spawns the `forwardHostSignals` goroutine BEFORE `cmd.Start()` at `host_exec.go:296`. The goroutine reads `cmd.Process` (`host_launch_unix.go:43`) while `cmd.Start()` writes it — no happens-before edge, so a SIGTERM/SIGHUP delivered to the launcher during the fork/exec startup window is an unsynchronized read/write (UB per the Go memory model). Reproduced under `-race`. This is net-new in #442 (the prior `syscall.Exec` model had no resident goroutine).

## Problem

The race is benign on the shipped arches (darwin/arm64, linux/amd64): `cmd.Process` is a pointer and an aligned pointer load/store is atomic, so no crash/corruption — the worst case is a dropped SIGTERM/SIGHUP forward in the few-ms fork/exec window (rare, and a session manager/supervisor escalates to SIGKILL past its grace timeout). But it is a shipped data race on the high-stakes front-door teardown path, and it carries a TEST-STRENGTH HOLE: every signal test (`host_launch_test.go:313`, `host_launch_pty_test.go:93`) fires its signal only AFTER the stub logs `STARTED` — long after `Start()` returns — so `go test ./... -race` stays GREEN while the race ships. This is exactly the coverage gap the README proof policy's adversarial-audit clause targets.

## Proposed approach (audit-specified)

Keep `signal.Notify` before `cmd.Start()` (the buffered channel already queues any early signal), but start the forwarding goroutine AFTER `cmd.Start()` returns, so the goroutine-creation happens-before edge publishes `cmd.Process` safely. This also closes the latent dropped-early-signal logic gap (an early signal now drains from the buffered chan once the goroutine starts).

## Acceptance criteria

- **AC-1 (VALUE)** — the data race is gone: a `-race` regression test that fires SIGTERM/SIGHUP into the launcher DURING the fork/exec startup window (not after `STARTED`) goes RED on the current `host_exec.go` goroutine-before-Start ordering and GREEN after the fix. Verified by: `go test ./internal/cli/... -race` with the new overlapping-signal test — RED pre-fix, GREEN post-fix (the divergeable proof the existing tests lack).
- **AC-2** — no behavior regression: an early signal (delivered before the goroutine starts) is still forwarded to the host (drained from the buffered chan), and signals delivered after startup forward as today; the existing TestLaunch* signal tests stay green.
- **AC-3** — no other launcher-startup goroutine reads a field written by `cmd.Start()` before a happens-before edge (sweep the spawn path for the same antipattern).

## Test plan

- New Go `-race` test exercising the overlapping window (fire the signal concurrently with / immediately around `cmd.Start()`), proven RED-first on the current ordering.
- `go test ./... -race` green.
- High-stakes front-door surface: detached adversarial audit at validation (per the README proof policy).
