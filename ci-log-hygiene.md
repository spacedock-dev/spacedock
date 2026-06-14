---
id: ecn07f3hwp5wgs8xf14h59sj
title: CI log hygiene — live-runner stream jsonl belongs in the artifact, not stdout
status: implementation
source: "captain (2026-06-14) — observed while debugging #368's opus `gate-guardrail` no-progress failure: the live-runner jsonl-to-stdout dump (`internal/ensigncycle/claude_live_runner_test.go:365`) bloated the CI log (~143KB for ~80 lines on one failed step) and buried the actual failure line."
started: 2026-06-14T19:16:23Z
completed:
verdict:
score: "0.30"
worktree: .worktrees/spacedock-ensign-ci-log-hygiene
issue:
sprint: 0203-fo-efficiency
---

The live shared-scenario runner dumps the full host jsonl stream to CI stdout — `internal/ensigncycle/claude_live_runner_test.go:365` `t.Logf`s every parsed stream line — so the whole conversation lands in the test log. That bloats CI output and buries the actual test failure. The jsonl is ALREADY uploaded as a per-scenario artifact (`claude-stream.jsonl`), so the stdout dump is pure redundancy. Keep the stream in the artifact; keep stdout to clean test-framework output (FAIL lines, assertions) — or only dump the stream to stdout on failure.

## Problem

- `claude_live_runner_test.go:365` logs each parsed stream line to the test log → the entire conversation lands in CI stdout.
- Debugging #368's opus failure (2026-06-14) meant fighting ~143KB of jsonl to find one `claude_live_runner_test.go:101` "no-progress quiet budget" line.
- The jsonl is already captured to an artifact, so duplicating it on stdout adds log cost and reviewer friction with no benefit.
- Likely mirrored in the Codex / Pi runners — confirm at ideation.

## Decision — artifact-only no-op tee, scoped to the two shared-scenario runners

**Recommendation: artifact-only.** Replace the per-line `t.Log` tee in the shared-scenario runners with a named no-op sink. The stream is already persisted to the per-scenario artifact (`claude-stream.jsonl` for Claude; the Codex runner's equivalent), so the stdout dump is pure redundancy. Dump-on-failure is rejected as unnecessary complexity (YAGNI): the failure paths already carry independent, bounded diagnostics that the tee never provided.

### Why artifact-only beats dump-on-failure

The runner's failure messages already surface the stream WITHOUT relying on the tee:
- `claude_live_runner_test.go:386-391` interpolates `transcriptTail()` and the artifact dir on the wrong-root / stall paths.
- `claude_live_runner_test.go:400-401` interpolates `tail(stream, 4000)` on a launch-extract failure.
- The full stream is on disk at `streamPath` (written at `claude_live_runner_test.go:373`) before any assertion runs.

So on failure a reviewer gets: the bounded tail in the log line itself, plus the complete jsonl in the artifact. The per-line `t.Log` tee adds a third copy — the entire conversation, unbounded, interleaved into framework output — which is exactly the ~143KB bloat that buried the #368 `gate-guardrail` no-progress line. Dump-on-failure would re-introduce a whole-stream `t.Log` (just gated on failure), re-bloating the log on precisely the runs a reviewer is debugging, while the artifact + tail already cover that case. No new diagnostic value, added branching — rejected.

### Codex / Pi scope (confirmed by reading the runners)

- **Codex** (`codex_live_runner_test.go:320`) mirrors the exact same `func(line string) { t.Log(line) }` tee and writes its stream to an artifact the same way → **in scope**, same one-line fix.
- **Pi** (`pi_live_runner_test.go:93`) sets `cmd.Stdout = stdout` (an `os.Create`d file) — it never tees stream lines to `t.Log` → **no dump to fix; out of scope** because there is nothing to change.
- **`live_test.go:150`** (the live FO-cycle smoke, `//go:build live`) ALSO tees `t.Log(line)`, but it does NOT persist the stream to an artifact (it only uses `fullTranscript()` for wrong-root detection). Removing its tee would lose the stream entirely on failure. That is a different change (it first needs an artifact write) and is **out of scope** for this task — recorded as a follow-on seed below.

## Out of scope

- `live_test.go:150`'s `t.Log` tee — it has no artifact backing, so silencing it would drop all stream diagnostics on failure. Follow-on: add a `streamPath` artifact write to the live FO-cycle test, then silence its tee the same way. (Backlog seed.)
- The Pi runner — it writes its stream to a file, never to stdout; nothing to change.
- Changing the `streamWatcher.tee` mechanism or its `transcriptTail()` failure-message wiring — both stay; only the runner-supplied callback changes.

## The change (concrete)

In `internal/ensigncycle/claude_live_runner_test.go:365` and `internal/ensigncycle/codex_live_runner_test.go:320`, replace:

```go
watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) { t.Log(line) })
```

with a named, shared no-op sink, e.g.:

```go
watcher := newStreamWatcher(newPipeLineSource(pr), poller, discardStreamLine)
```

where `discardStreamLine` is a single package-level `func(string) {}` declared once (with a comment naming WHY: the full stream is persisted to the per-scenario artifact and failure messages already carry the bounded `transcriptTail()`, so teeing every line to `t.Log` is redundant and bloats CI output). The watcher still drains, parses, records into `transcript`, and tracks dispatches exactly as before — only the per-line stdout echo is dropped.

No user-visible CLI surface, banner, or docs-site behavior changes (this is internal test wiring), so no doc diff is owed.

## Acceptance criteria

- **AC-1 — runner stream sink emits no stdout.** An offline Go unit test drives `newStreamWatcher` with the SAME callback the runners pass (`discardStreamLine`) over the real captured fixture `internal/ensigncycle/testdata/sonnet_teamdelete_hang.stream.jsonl` (341 lines: 65 `"type":"assistant"`, 32 `"type":"user"`), draining the whole stream, and asserts the sink received ZERO lines while the watcher's `transcript` still holds all 341 lines (parse/record path intact).
  - **Verified by:** a new test in `internal/ensigncycle/streamwatch_regression_test.go` (or sibling, default build tag) that wires a recording sink alongside `discardStreamLine`, runs the watcher to exit, and fails if the recorder count != 0 or the transcript length != fixture line count. The expected values (0 emitted; 341 recorded) come from the fixture and the recorder — independent of any prose in this task. Inverting the runner back to `t.Log`/a teeing sink makes the recorder count non-zero and the test fails. Run: `go test ./internal/ensigncycle/ -run TestStreamSinkDiscardsLines`.
- **AC-2 — both shared-scenario runners use the silent sink, neither uses `t.Log` per line.** A guard test asserts the runner files no longer wire a per-line `t.Log` tee.
  - **Verified by:** since this is a wiring fact, the durable proof is AC-1's behavioral test bound to `discardStreamLine` PLUS a `go vet ./internal/ensigncycle/` + full-package `go test ./internal/ensigncycle/` green after the edit (compiles, the shared `discardStreamLine` symbol is referenced by both runners). The behavioral anchor is AC-1; AC-2 is the compile/reference fact, not a grep over instruction prose.
- **AC-3 — the artifact still contains the full stream.** The `os.WriteFile(streamPath, []byte(stream), …)` at `claude_live_runner_test.go:373` (and the Codex equivalent) is unchanged, so the per-scenario `claude-stream.jsonl` artifact still holds the complete jsonl.
  - **Verified by:** the diff touches only the `newStreamWatcher(...)` tee argument and adds the `discardStreamLine` declaration — the `streamPath` write lines are untouched (confirmable in the implementation diff / git log). The existing live-runner suite's artifact write is exercised by any live scenario run.

## Test plan

- **Primary (offline, free, deterministic):** AC-1's unit test in `internal/ensigncycle/`. Reuses the existing `drippingLineSource` + `transcriptRecorder` harness from `streamwatch_regression_test.go` and the existing 341-line fixture — no new fixture authoring. Cost: ~sub-second, no credentials, runs in normal CI. This is the riskiest-mechanism exercise AND the regression guard in one: it proves that the callback the runners hand `newStreamWatcher` emits nothing to stdout while the drain/parse path is unaffected.
- **Compile/reference:** `go vet` + `go test ./internal/ensigncycle/` green confirms both runners reference the shared `discardStreamLine` and the package builds (AC-2).
- **No live run required for the AC.** The bloat reproduces only under a live API run, but the BEHAVIOR being fixed — "the runner's stream sink writes to the test log" — is fully captured offline by feeding a recorded stream through the real watcher with the real callback. A live `gate-guardrail` scenario run would confirm the end-to-end log size drop, but it is not needed to make the AC checkable and is left as optional manual confirmation.
- **No spike needed:** the design only composes already-proven behavior — `newStreamWatcher`'s tee/transcript split (proven by `streamwatch_unit_test.go` and `streamwatch_regression_test.go`), the existing artifact `WriteFile`, and the existing failure-message `transcriptTail()`/`tail()` diagnostics. The one mechanism the fix's soundness rests on (the watcher still drains and records every line when the tee is a no-op) is exercised directly by AC-1's test, not assumed.

## Stage Report: ideation

- DONE: Decide and specify the fix for claude_live_runner_test.go:365's full-stream `t.Logf` dump: artifact-only vs dump-on-failure-only — recommend one with rationale
  Recommended artifact-only (named no-op `discardStreamLine` sink); dump-on-failure rejected — failure messages already carry `transcriptTail()`/`tail(stream,…)` and the artifact holds the full jsonl, so dump-on-failure re-bloats logs on exactly the debugged runs with no added value (Decision section).
- DONE: Confirm whether the Codex/Pi runners mirror the same stdout dump and scope them in this task or a follow-on.
  Read all three: Codex (codex_live_runner_test.go:320) mirrors the tee → IN scope; Pi (pi_live_runner_test.go:93) writes `cmd.Stdout=stdout` to a file, no `t.Log` tee → nothing to fix; live_test.go:150 tees but has no artifact backing → follow-on backlog seed (out of scope).
- DONE: AC bound to an external check (NOT a prose/regex match): captured test stdout contains NO raw stream jsonl while `claude-stream.jsonl` remains in the artifact.
  AC-1 is an offline Go unit test driving `newStreamWatcher` with the runner's `discardStreamLine` callback over the real 341-line `testdata/sonnet_teamdelete_hang.stream.jsonl` (65 assistant/32 user lines), asserting sink-emitted==0 AND transcript==341; expected values come from the fixture+recorder, independent of the task body, and inverting the tee back to `t.Log` fails it. On-failure path not chosen, so its fixture verification is N/A.

### Summary

Recommended artifact-only over dump-on-failure: the stream is already persisted to the per-scenario artifact and the failure messages already carry a bounded `transcriptTail()`/`tail()`, so the per-line `t.Log` tee is pure redundancy and the source of the ~143KB #368 bloat. Scoped the fix to the two shared-scenario runners (Claude + Codex), excluded Pi (no stdout tee) and live_test.go (no artifact backing — follow-on seed). The AC is an offline unit test over the existing recorded fixture using the runners' actual `discardStreamLine` callback — checkable in normal CI, no live API run, no spike needed (the only load-bearing mechanism, the watcher still draining/recording with a no-op tee, is exercised by the test itself).
