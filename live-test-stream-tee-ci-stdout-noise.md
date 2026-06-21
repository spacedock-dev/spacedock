---
title: Live ensign-cycle tests re-flood CI stdout on failure — inline t.Log stream tees bypass the discardStreamLine sink
status: backlog
score: ""
source: captain noticed CI stdout still mixed with raw session JSONL (Runtime Live E2E run 27892346952) and was surprised it had been "fixed twice"
priority: medium
id: nf119vts6ggcyc4hz669hp0g
---

## Problem

The "Run live ensign cycle" CI step dumped ~227KB of raw Claude session JSONL into stdout (e.g. `zero_discover_live_test.go:82: {"type":"system",...}` / `{"type":"assistant",...}`), burying the real assertion failure. This bloated-stdout-from-session-JSONL problem was addressed twice before — but for OTHER surfaces:

- **#372 `80db9653` (ci-log-hygiene)** replaced the per-line `func(line){ t.Log(line) }` stream tee with the quiet `discardStreamLine` sink in the SHARED-SCENARIO runners (`claude_live_runner_test.go`, `codex_live_runner_test.go`), guarded by `TestStreamSinkDiscardsLines`.
- **#389 `82db42e6` (ci-log-read-summary)** moved every live step to `gotestsum --jsonfile … --format pkgname` (clean step log + archived `-json` detail).

### Why it is back (a new uncovered surface, not a regression of the fixes)

- **#374 `5ead2a74` (lean-boot)** created `internal/ensigncycle/zero_discover_live_test.go` ~1h after #372, reintroducing the inline `newStreamWatcher(…, func(line string){ t.Log(line) })` tee that #372 had just removed elsewhere. It was never converted to `discardStreamLine`.
- The #372 guard `TestStreamSinkDiscardsLines` pins only the `discardStreamLine` SYMBOL, so it never constrains inline-`t.Log` callers.
- `gotestsum --format pkgname` suppresses per-test output ONLY for PASSING tests. On a FAILED test it prints the failed test's full accumulated output, including every `t.Log(line)`. So the firehose appears precisely when a test reds (run 27892346952 failed → dump).
- The #389 state doc (line 156) already KNOWINGLY flagged this residual gap: reverting the tee to inline `t.Log` stays green because it sits outside the discardStreamLine-symbol contract.

### Scope

6 loud inline-tee call sites vs 4 quiet `discardStreamLine` sites:

- Loud: `zero_discover_live_test.go:82`, `live_test.go:305`, `live_gate_stop_test.go:141`, `merged_team_mode_live_test.go:137`, `pty_team_mode_live_test.go:100` & `:196`.
- Quiet: `claude_live_runner_test.go:403`, `codex_live_runner_test.go:320`, `haiku_loop_spike…:292`, `pty_live_driver…:327`.
- CI step `.github/workflows/runtime-live-e2e.yml:231` runs three loud tests together under one step, so the dump is their combined accumulated stream.

## Acceptance criteria

- **AC-1 (value)** — When a live ensign-cycle test FAILS in CI, the step's stdout no longer carries the raw per-event session JSONL firehose; the failing assertion is readable without scrolling past the stream dump. Baseline: run 27892346952 dumped ~227KB on failure — the measure must flip that to a clean, bounded failure log. Measured on a real failing run (or a faithful reproduction), not just a symbol assertion.
- **AC-2 (route through the quiet sink)** — The 6 inline `func(line){ t.Log(line) }` tees passed to `newStreamWatcher` are replaced with `discardStreamLine` (or an equivalent gated sink). The full transcript stays drained and is surfaced on failure via `transcriptTail()` (plus the per-scenario `*-stream.jsonl` artifact where applicable) — debuggability preserved.
- **AC-3 (guard the gap)** — A guard prevents recurrence: a source-level check (analogous to the cilog guards) asserts NO live `newStreamWatcher` call site passes an inline `t.Log` tee. Pinning the `discardStreamLine` symbol alone is what let #374 slip past, so the guard must constrain the call sites, not just the sink symbol.
- **AC-4 (validation)** — `go test ./internal/ensigncycle/…` and the cilog/streamsink guard package green.

## Notes

- Distinct from `zero-discover-broad-search-hardening` (the FO broad-filesystem-sweep model flake that surfaced in the same run): that entity tracks the test FAILURE; this tracks the log NOISE. Quieting the tee helps CI regardless of the flake.
- Prior fixes: #372 `80db9653`, #389 `82db42e6`; reintroduction: #374 `5ead2a74`.
