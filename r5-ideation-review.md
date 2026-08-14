# IDEATION GATE — Live CI API-error/retry log capture (`r5`)

Recommendation: **APPROVE and dispatch implementation.**

## Selected approach

Always-on capture for the **claude lane only**, armed via the host flags the
runner already constructs (not a second mechanism). Codex and pi need no change
— the spike confirmed their existing stderr + JSONL/session artifacts already
cover API-error/retry output.

1. **Claude:** add `--debug api --debug-file <artifact-dir>/claude-api-debug.log`
   to the `--` host-flag block in `claudeLiveRunner.run()`. The debug file is a
   separate sink — it does NOT touch the stream-json pipe the watcher drains for
   liveness. `--debug api` filters to the API category, suppressing
   permission/MCP/plugin noise.
2. **Kill message:** append the debug-log path to the stall failure message so a
   no-progress kill points at the captured lifecycle.

## Risk evidence

- **Spike done (the unverified mechanism).** Live probe
  `claude -p "say hi" --debug api --debug-file /tmp/claude-api-debug-probe.log`
  confirmed the `api` category emits the request lifecycle to the debug file:
  `[API:timing] dispatching…`, `[API REQUEST] /v1/messages…`,
  `Stream started - received first chunk`, `[API:timing] first byte after 758ms`.
  On a dropped stream or retry, this channel is where the attempt/error/retry (or
  its absence) appears — the exact hole a mid-stream hang leaves in the
  stream-json transcript.
- **Codex/pi confirmed-absent, not assumed.** `codex exec --help` / `pi --help`
  probed: no debug/API-lifecycle flag; their runners already capture stderr +
  JSONL/session artifacts separately (codex `codex-exec.stderr.txt` +
  `codex-exec.jsonl`; pi `pi-stdout.txt` + `pi-stderr.txt` + session JSONL).
- **Cost risk:** `--debug api` file was ~50KB on a trivial run (startup noise the
  `api` filter doesn't suppress); long scenarios could reach low-MB. Acceptable
  for a diagnostic-only artifact; the test plan measures on the first real green
  run.

## Expected surface and tolerance

- `internal/ensigncycle/claude_live_runner_test.go` — ~5-10 lines (argv + kill
  message).
- An argv shape unit test asserting `--debug api --debug-file` present.
- **Estimate: ~+10 lines net, tolerance ±5.** No semantic changes — diagnostic
  capture only; no gating on the captured output; no stream-pipe change.

## Semantic changes

None. The debug file is a separate sink; the stream-json transcript the watcher
drains is unchanged. No CLI output, stored format, authority, or runtime
behavior change.

## Proposed proof per acceptance criterion

- **AC-1 (a stalled scenario leaves an API-lifecycle artifact):** an induced or
  real stall's `claude-api-debug.log` shows the request lifecycle around the stall
  (opportunistic; the shape test is the gate). Baseline: today only stream
  silence is archived.
- **AC-2 (negligible green-run cost):** measure artifact size / wall-clock delta
  on the first real green run; report.
- **AC-3 (harness kill message names the artifact):** the stall failure message
  includes the debug-log path.
- **Shape test (the gate):** argv asserts `--debug api --debug-file <path>`
  present in the claude runner's host flags.

## Decision ask

Approve to dispatch implementation (worktree) for this claude-only diagnostic
capture, or revise/hold with a concrete boundary.
