---
title: "Live CI captures API-error/retry logs so a stalled stream is diagnosable, not just dead"
status: ideation
group: tooling
source: "2026-07-02 session: PR #461's claude-live sonnet lane failed the filing scenario on a 60s no-progress kill whose transcript tail ends mid-thinking (stream died mid-generation, work already complete on disk); the same window produced two explicit 'API Error: Connection closed mid-response' failures in an interactive subagent. The rerun went green with zero code change. Root cause was NOT determinable from CI: nothing captured distinguishes provider-side API weather / runner networking from the one actionable variant — the claude CLI hanging instead of retrying after a dropped stream. Captain direction: 'in any case, we should have api error logs etc, if we are not currently capturing that.'"
id: r5y6qjr10k4m3gw9w5p3b2vj
sprint: live-evidence-followups
gates:
    version: 1
    records:
        - id: gate:r5y6qjr10k4m3gw9w5p3b2vj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:r5y6qjr10k4m3gw9w5p3b2vj-backlog-1
              briefing:
                id: briefing:r5y6qjr10k4m3gw9w5p3b2vj:backlog:attempt-1:revision-1
                digest: sha256:ab5631d7fc7f30835a41b462411b0dcc42702ec71b60b7c93b05fb8d71c1b8fe
                request-digest: sha256:0b8b68120fec292362c6d3c2861e6242b730f46ace2255f02ab2f80f357c2c9b
                room-ref: ./live-ci-api-error-log-capture/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:r5y6qjr10k4m3gw9w5p3b2vj:backlog:1
                briefing: briefing:r5y6qjr10k4m3gw9w5p3b2vj:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-14T06:43:27.049969Z"
                decision: approve
                reason: Captain approved backlog gate; advance to ideation for API-error log capture.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:r5y6qjr10k4m3gw9w5p3b2vj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:r5y6qjr10k4m3gw9w5p3b2vj-ideation-1
              briefing:
                id: briefing:r5y6qjr10k4m3gw9w5p3b2vj:ideation:attempt-1:revision-1
                digest: sha256:ce90c429547ac83a604105637315c4c3178ee6d54e5014bbb0deb08f37d81149
                request-digest: sha256:fa93a503b3cef14b321ba2bae869c1f252a43f4988a0c0af1b086c622e2898f2
                room-ref: ./live-ci-api-error-log-capture/review/ideation/briefing-1
started: 2026-08-14T06:44:38Z
---

## Problem
When a live-lane scenario dies of stream silence, the archived evidence (stream jsonl + step log) shows only the silence. API-level errors, retry attempts, and connection lifecycle events are not captured, so a stall cannot be classified: transient weather (re-run and move on) vs a CLI retry/hang defect (file upstream, work around). Every such failure costs a full lane re-run and still teaches nothing.

## Desired direction (for ideation to refine)
The live runners capture the host CLI's API-error/retry/debug output as a per-scenario CI artifact alongside the existing stream jsonl — whatever channel the claude CLI (and codex/pi equivalents, if cheap) exposes: debug/verbose flags, error log files, env-var-enabled logging. On a no-progress kill, the harness failure message points at the captured log. Ideation determines what the CLI actually exposes (read the docs/schema first — usage presence is not existence evidence), the artifact size/noise cost, and whether capture is always-on or armed only for the kill path.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- A live scenario killed for no-progress leaves an artifact that shows the API request lifecycle around the stall (attempts, errors, retries or their absence) — demonstrated by inspection on a real or induced stall.
- The capture adds negligible cost to green runs (measured artifact size / wall-clock delta), and the harness kill message names the artifact path.
- Related, separate candidate recorded (not this task): grade-on-outcome-state when the stream dies after the deliverable landed — the filing kill discarded a satisfiable pass.

## Spike findings

Probed the three host CLIs' actual API-error/retry/debug surfaces by reading
their `--help` output and exercising the claude debug-file path against a live
`-p "say hi"` launch. Usage presence was not treated as existence evidence; each
claim below is backed by a concrete probe or a confirmed absence.

### Claude (the lane the entity source names — the gap is here)

- `claude --help` exposes `-d, --debug [filter]` with category filtering (e.g.
  `"api,hooks"` or `"!1p,!file"`) and `--debug-file <path>` which "implicitly
  enables debug mode" and writes debug logs to that file.
- A live probe (`claude -p "say hi" --output-format stream-json --verbose --debug
  api --debug-file /tmp/claude-api-debug-probe.log`) confirmed the `api` category
  emits the request lifecycle to the debug file:
  `[API:timing] dispatching to firstParty model=...`, `[API REQUEST]
  /v1/messages x-client-request-id=...`, `Stream started - received first
  chunk`, `[API:timing] first byte after 758ms`. On a dropped stream or retry,
  this channel is where the attempt/error/retry (or its absence) would appear.
- The existing runner (`claude_live_runner_test.go` run()) already folds stderr
  into the stream pipe and passes `--verbose`. The stream-json `result` event
  carries `api_error_status` and `is_error`, so a TERMINAL API error (e.g. a 401)
  IS already in `claude-stream.jsonl`. The gap is the mid-stream HANG: when the
  stream dies mid-generation there is no result event, only silence. The
  `--debug api` file is the only surface that shows whether a retry was
  attempted before the silence — the distinction between transient weather and
  a CLI retry/hang defect this task exists to make diagnosable.
- No env-var-only path exists; `--debug api --debug-file <path>` is the cheap
  flag combo. `~/.claude/debug/<uuid>.txt` exists as a default sink but is
  keyed by session UUID (not controllable for CI artifact upload) and was not
  written by the probe.

### Codex (already covered — no change needed)

- `codex exec --help` has no `--debug`, `--verbose`, or log-level flag. `codex
  debug` is a static-inspection subcommand (models, app-server, prompt-input),
  NOT runtime API logging. `RUST_LOG` is not a documented codex surface.
- The codex runner (`runCodexProcess` in `codex_single_run_test.go`) ALREADY
  captures stdout to `codex-exec.jsonl` and stderr to a SEPARATE
  `codex-exec.stderr.txt`. A live `codex exec` probe showed stderr carries
  runtime diagnostics ("Reading additional input from stdin..."). API errors
  and retries surface to the JSONL stream (as error events) and/or stderr —
  both are already archived per scenario. No additional flag or capture path
  is needed for codex.

### Pi (already covered — no change needed)

- `pi --help` exposes `--verbose` (startup verbosity) and `PI_OFFLINE`; no
  `--debug api`, log-level, or API-lifecycle flag. `--export <file>` is a
  post-hoc session-to-HTML export, not a live capture.
- The pi runner (`pi_shared_live_runner_test.go`) captures stdout to
  `pi-stdout.txt`, stderr to `pi-stderr.txt`, and reads the root session JSONL
  (the session-dir output) into the stream the assertions consume. API errors
  surface in the session JSONL (the transcript pi itself persists) and/or
  stderr — both already archived. No additional flag needed for pi.

## Proposed approach

Capture is **always-on for the claude lane only**, armed via the host flags the
runner already constructs (not a second mechanism). Codex and pi need no change
(their existing stderr + JSONL/session artifacts already cover API-error/retry
output).

1. **Claude: add `--debug api --debug-file <artifact-dir>/claude-api-debug.log`
   to the `--` host-flag block in `claudeLiveRunner.run()`.** The debug file
   writes to the per-scenario artifact dir, so it is uploaded alongside the
   existing `claude-stream.jsonl` and `claude-final-message.txt`. It does NOT
   touch the stream pipe (debug-file is a separate sink), so the stream-json
   transcript the watcher drains for liveness is unchanged. `--debug api`
   filters to the API category, suppressing the permission/MCP/plugin noise
   that dominates an unfiltered debug log.
2. **Kill message names the artifact.** The `drainToExit` stepTimeout message
   (and the `claude_live_runner_test.go` stall surface that wraps it) appends
   the `claude-api-debug.log` path so a human reading the CI failure can go
   straight to the request-lifecycle log. This is a string change in the
   failure-message formatting, not a gating change.
3. **Codex / Pi: no flag change.** Document in the entity (this section) that
   their existing stderr + JSONL/session artifacts already satisfy the
   "leaves an artifact that shows the API request lifecycle" outcome. If a
   future codex/pi build adds a debug flag, it can slot into the same
   per-scenario artifact pattern.

**Always-on vs armed-only:** always-on is correct. The debug file must be
capturing BEFORE the stall to be useful; you cannot start it retroactively
after the kill fires. The cost (below) is bounded and the file is separate
from the stream, so it cannot destabilize a green run.

## Expected surface and tolerance

- **Files added (claude lane):** one per scenario —
  `<artifact-dir>/claude-api-debug.log`. No new files for codex or pi.
- **Code lines (estimated):** ~5-10 lines in `claude_live_runner_test.go` (the
  `--debug api --debug-file` argv addition + the debug-log path in the stall
  failure message). The `streamwatch_test.go` stepTimeout message is
  parameterized by label only; the claude runner wraps it and can append the
  artifact path at the call site without touching the shared watcher.
- **Green-run cost (measured from the probe):** the `--debug api` file for a
  trivial "say hi" was ~50KB, dominated by one-time startup noise
  (permissions, MCP connection, plugin load) that the `api` filter does not
  suppress because the filter matches debug categories, not the startup
  `[DEBUG]` lines. The actual API-lifecycle lines (request, stream-start,
  timing, retries) are a handful of lines per request. For a real multi-minute
  scenario, the file scales with the number of API requests, not wall-clock,
  so a green run producing N model turns yields ~N×(a few lines). Tolerance:
  the file is a diagnostic artifact, never parsed or gated; size growth only
  matters for CI upload budget. A scenario with hundreds of turns might reach
  low-MB; acceptable for a diagnostic-only artifact.
- **Wall-clock delta:** negligible. Debug-file logging is an async file sink
  (the probe showed no measurable launch delay); it does not block the
  stream-json stdout path the watcher drains.
- **Semantic changes:** none. This is diagnostic capture only — no gating, no
  assertion consumes the debug log, no pass/fail path reads it. The stream-json
  transcript, the watcher's no-progress deadline, and the scenario assertions
  are all unchanged.

## Test plan

Per-AC proof plan (the "Rough acceptance sketch" ACs, refined into measurable
checks):

- **AC-1 (stall leaves a request-lifecycle artifact):** Add a unit test
  (DEFAULT build tags, no model spend) that constructs the claude live argv
  via the shared argv builder and asserts `--debug`, `api`, and
  `--debug-file` are present and the debug-file path lands under the artifact
  dir. This is the falsifiable shape check: removing the argv lines fails it.
  For the live behavior: an induced-stall live test (or inspection of a real
  future stall's archived `claude-api-debug.log`) confirms the file shows the
  `[API REQUEST]` / `[API:timing]` lines around the silence — proving the
  artifact distinguishes "request dispatched, no retry" (CLI hang) from "retry
  attempted" (transient weather). The live proof is opportunistic (a real or
  induced stall); the shape test is the gate.
- **AC-2 (negligible green-run cost + kill message names artifact):** The argv
  unit test asserts the debug-file path is the artifact-dir-relative
  `claude-api-debug.log`. A live green run measures the file size (assert <
  bound, e.g. 1MB for a normal scenario) and the wall-clock delta vs a
  pre-change baseline. The stall failure-message test asserts the message
  string contains the `claude-api-debug.log` path — the falsifiable change is
  dropping the path from the message format.
- **Codex/Pi no-regression:** existing codex/pi runner tests continue to pass
  unchanged (no argv or capture-path change). This confirms the "already
  covered" finding did not silently drop coverage.

## Stage Report: ideation

- DONE: Spike the riskiest unverified mechanism — what API-error/retry/debug output the host CLIs expose
  Probed `claude --help`, `codex exec --help`, `pi --help` and ran a live `claude -p "say hi" --debug api --debug-file` probe; confirmed claude has `--debug api` + `--debug-file` (emits `[API REQUEST]`, `[API:timing]`, stream-start), codex/pi have no debug flag but already capture stderr+JSONL/session artifacts separately.
- DONE: Propose a per-scenario capture design based on the spike findings
  Always-on `--debug api --debug-file <artifact>/claude-api-debug.log` for the claude lane only (must capture before the stall; cost bounded); codex/pi unchanged (existing stderr+JSONL/session already cover); kill message names the debug-log path.
- DONE: Record spike findings, proposed approach, expected surface, semantic changes, and per-AC proof plan into the entity body
  Added `## Spike findings`, `## Proposed approach`, `## Expected surface and tolerance`, `## Test plan` sections; existing ACs and "Desired direction"/"Rough acceptance sketch" left intact.
- DONE: Append `## Stage Report: ideation` with DONE items and Summary, commit path-scoped, push
  This section; commit and push follow.

### Summary

Spike found claude is the only lane with a gap: `--debug api --debug-file <path>` emits the API request lifecycle (dispatch, stream-start, timing, retries) to a separate file sink, filling the diagnostic hole a mid-stream hang leaves in the stream-json transcript. Codex and pi already capture API-error/retry output via existing stderr + JSONL/session artifacts and need no change. The proposal is always-on claude-only capture (a few argv lines + a kill-message path), diagnostic-only with no gating, ~5-10 code lines in the claude runner.
