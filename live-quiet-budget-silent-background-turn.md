---
title: Live harness 60 s quiet budget kills runs while a background worker has a silent turn
status: ideation
source: captain (2026-11-03), from the repair-sonnet-live-flakes ideation probes
score: 0.6
sprint:
group: live-lane
sprint-readiness: ready
id: w6kcbj4k8ev9t8n2hpw52dds
started: 2026-11-04T07:38:35Z
---

`drainToExit` in the live harness kills a run when no stream line arrives for 60 s (`quietBudgetDefault`). In the keep-moving probes of `repair-sonnet-live-flakes` (Sonnet 5.5, effort max, head `64274cef5`), 7 of 8 unchanged-byte runs were killed this way. Each kill came after the FO ended its turn while the background `questioned` worker was in a long thinking turn. In run b1 (budget raised to 180 s), silent gaps of 78 to 128 s occurred five times. These kills are infra failures, not grades, and the XFAIL binding does not absorb them, so they fail the lane.

The kill can depend on load: three other workers shared the machine, and earlier the same journey passed 3 of 3 locally. CI Sonnet 5.5 ran the journey in about 62 s.

Evidence: the `### Finding: quiet budget kills runs while a background worker is in a silent turn` section of `repair-sonnet-live-flakes.md`.

Seed acceptance criteria: a live run is not killed while a dispatched background worker is alive and its process is making progress, and a run whose processes are all stalled still fails fast; proven by a deterministic harness test with a fake stream that goes silent for longer than the budget while a child is alive (no kill) and while no child is alive (kill).

## Problem

The exact rule that kills a run (at `64274cef5`): `claudeLiveRunner.run` (`internal/ensigncycle/claude_live_runner_test.go:731`) calls `watcher.drainToExit(quietBudgetDefault, …)`. `drainToExit` (`internal/ensigncycle/streamwatch_test.go:396`) sets `deadline = now + 60 s` and moves it forward only when `drainEntries` returns one or more new stream-json lines. When `time.Now()` passes the deadline and the `spacedock claude` process has not exited, it kills the process and returns `stepTimeout`. The runner then calls `t.Fatalf` before it grades anything. No other input changes the decision: the watcher already parses host task events, but `drainToExit` ignores them.

In a1 to a4 the FO had ended its turn, so the only producer of lines was the background `questioned` worker. A max-effort thinking turn writes no line until it ends, and b1 had five gaps of 78 to 128 s.

### Can the harness tell a live background worker from a stalled run?

It can tell that a worker is **alive**. It cannot see that a worker is **making progress** during a silent thinking turn. I checked each candidate signal:

1. Host task events in the stream: usable. The host writes `system/task_started` (`task_id`, `task_type`) when a worker starts. It writes `system/task_notification` (same `task_id`, `status` `completed`, `stopped` and so on) when the worker ends. All three recorded Claude streams in `internal/ensigncycle/testdata/` contain this pair: `claude_live_auto_continue_run31915540750_sonnet` and `haiku_loop_spike_runs/run-1` (`task_type: local_agent`), and `sonnet_teamdelete_hang` (`task_type: in_process_teammate`, team mode, with no `subagent_type`). b2 saw `task_notification status: stopped` for the stopped worker. Host 2.1.303 also writes a `background_tasks_changed` snapshot (b2 saw `tasks: []`), but none of the recorded fixtures has it, so it depends on the host version. I do not use it.
2. `system/task_progress`: no help. A read of the host binary (2.1.303) shows it fires once per worker message that ends in a `tool_use`. That is the same moment a stream line appears, so it is also silent during a thinking turn.
3. `--include-partial-messages`: rejected. In the host binary, headless `stream_event` output is main-thread only (`parent_tool_use_id:null`) and drops pings, so it shows nothing from a background worker. This is unverified live. It would also make the stream much larger.
4. Process tree: rejected. Claude subagents run inside the host process. A thinking turn has no child process. To sample the host's CPU time, the harness would need a threshold that is not deterministic.
5. Subagent transcript files under `CLAUDE_CONFIG_DIR`: rejected. The host writes them once per message, the same timing as the stream.

The host has bounds of its own (binary reading, not documented). A stream idle watchdog stops a stalled API stream after `max(config, 300 s)` ("Agent stalled (stream watchdog)"). b2 also saw the host stop a background worker 600 s after the FO's last turn.

So the harness has a deterministic signal for a live worker, and no signal for progress. The proposal therefore puts a limit on silence while a worker is registered. It does not claim to detect progress.

### Spike (throwaway clone of `64274cef5`, now removed)

- Patch: +24/-6 in `streamwatch_test.go`, with the design below. Fake-clock table test, 5 cases, all pass: worker alive with 150 s of silence then exit, no kill. b1 gap replay (128, 104, 100, 97 and 78 s) with the worker alive, no kill, clean exit at 8m33s fake time. No worker and 61 s of silence, killed at 60.2 s. Worker `stopped`, then silence, killed at 70.2 s. Worker alive but silent for 181 s, killed at 180.2 s.
- Falsified: I disabled the extension, which gives the base rule. The two worker-alive cases then fail at 1m0.2s. I closed the registry only on `completed`. The `stopped` case then fails.
- Fixture replay over the three recorded streams: the registry is non-empty between each `task_started` and its `task_notification`, and empty at the end. My first filter (`subagent_type != ""`) missed the team-mode `in_process_teammate` worker. That is why the filter is on `task_type`.
- `go vet -tags live ./internal/ensigncycle` passes, and `go test ./internal/ensigncycle` passes with the patch (902 s).

## Proposed approach

Change one function, and reuse the sanctioned 3-minute constant.

1. `streamEntry` gets `TaskID` (`task_id`) and `TaskType` (`task_type`).
2. `streamWatcher` gets `runningWorkers map[string]bool`. `updateDispatches` adds the `task_id` of each `task_started` whose `task_type` is not `local_bash`. It removes the `task_id` on any `task_notification`, whatever the status.
3. `drainToExit` measures silence from the last drained line. The limit is `budget`, or `quietBudgetDispatchClose` (3 min) when `runningWorkers` is non-empty and that value is larger. The function uses the watcher's `now`/`sleep` clock seam, which `drainCodexToTerminal` already uses. The kill message names the limit that applied and the number of running worker tasks.

Callers: the Claude shared runner gets 60 s with no worker and 180 s with a worker. The haiku spike already passes 180 s, so nothing changes for it. The Codex runner uses `drainCodexToTerminal` and is not affected. `expect`, `expectDispatchClose` and `expectCondition` stay unchanged, because no kill has been observed in them.

Alternatives considered (each new mechanism serves AC-1 and AC-2):
- Raise `quietBudgetDefault` to 180 s. Rejected: a run with no worker would then take 180 s to fail, which breaks AC-2.
- Reuse `openDispatches`. Rejected: it closes only on `completed` and Done anchors, so a stopped worker stays open. If its close rule changed, `closedCount` would change, and `expectDispatchClose` relies on it.
- Use the `background_tasks_changed` snapshot. Rejected: the recorded older-host fixtures do not have it.
- Count every task type, `local_bash` included. Rejected: a stalled run that has a silent background Bash would then wait 180 s. The worker types differ by mode (`local_agent`, `in_process_teammate`), so the rule excludes Bash instead of listing worker types.
- Add a progress probe (CPU, partial messages). Rejected: see the signal list above.

Semantic changes: only runtime behavior of the test harness, namely when the Claude live lane kills a silent run, and the text of the kill message. No command grammar, stored format, authority or product code changes.

Expected surface: Estimate net LOC change: +140, across 3 files (insertions about +148, deletions about -8). The files are `internal/ensigncycle/streamwatch_test.go` (about +28/-6), `internal/ensigncycle/streamwatch_unit_test.go` (about +118, new tests and a scripted fake) and `internal/ensigncycle/claude_live_runner_test.go` (a comment at `:666`/`:728` that names the worker extension, about +2/-2). Tolerance: net +90 to +200, 2 to 3 files.

Docs: no doc diff. `docs/site/contributing/architecture-notes.md:30` ("Liveness remains the transport's no-progress quiet budget") stays accurate.

## Acceptance criteria

1. AC-1 (value): a Claude shared-scenario run is not killed by the quiet budget while the host reports a running worker task and no silence lasts longer than 180 s. Test: in a fake-clock `drainToExit` test, a replay of b1's five observed gaps (128, 104, 100, 97 and 78 s) with a registered worker ends in a clean exit with no kill. With the base rule (no extension), the same replay is killed at 60 s. Falsify: remove the extension. The subtest then fails at 1m0.2s, as in the spike.
2. AC-2: a run with no registered worker still fails at the 60 s budget. A worker that has ended (any `task_notification` status, `stopped` included) no longer extends the budget. Test: subtests "no worker, 61 s silent" (killed by 60.2 s) and "worker stopped, then silent" (killed within 60.2 s of the notification). Falsify: close the registry only on `completed`. The stopped subtest then fails.
3. AC-3: the extension has a limit. A registered worker that stays silent for 181 s is killed by 180.2 s. Test: subtest "worker alive, silent past 3 min". The fake clock fails the test once fake time passes 1 h, so an unlimited wait fails the test and cannot hang it. Falsify: make the extension unlimited while a worker runs. The subtest then fails on the clock cap.
4. AC-4: the worker registry follows real host streams. Over the three recorded Claude fixtures (background `local_agent`, auto-continue `local_agent`, team-mode `in_process_teammate`), the registry is non-empty between each `task_started` and its `task_notification`, and empty at the end. Test: a fixture-replay unit test. Falsify: filter on `subagent_type`. The team-mode fixture then fails, as seen in the spike.
5. AC-5: no regression. `go test ./...`, `go test ./... -race`, `gofmt -l ./cmd ./internal` (no new entries) and `go vet -tags live ./internal/ensigncycle` pass. `TestLiveBudget*` passes, and its allowlist does not change.

## Test plan

- Primary proof owner: the existing `TestDrainToExit*` tests in `internal/ensigncycle/streamwatch_unit_test.go`. They must keep passing on the clock seam. The new tests go in the same file.
- New `TestDrainToExitWorkerSilence`: a table test with a scripted line source, a fake process and a fake clock. It covers AC-1 to AC-3, takes about 0 s and is deterministic. The falsifying edit for each case is named in its AC.
- New `TestRunningWorkersFollowRecordedStreams`: it replays the existing testdata streams (AC-4) and takes about 0.02 s.
- Cost: under 1 s of test time. A full `ensigncycle` package run takes about 15 min on this loaded machine.
- Live runs: none are needed to prove the ACs, because each AC is about watcher behavior on a given stream. Whether local keep-moving runs then get graded depends on the real gap distribution (b1's largest gap was 128 s, so the margin is 52 s). The retirement runs that `repair-sonnet-live-flakes` already budgets (about $12 to $15) will show this. An optional extra local probe costs about $3 and 12 to 28 min. It is not needed for the gate.

## Risk: the change could hide a real hang

- A registered worker that is really hung (stalled API stream, deadlocked tool), or an FO that hangs while a worker is registered: the harness kills the run after 180 s of silence instead of 60 s. The extra delay is at most 120 s, and only once per run, because the kill ends the run. Behind that are the host's own limits: the stream watchdog (300 s or more) and the observed 600 s stop of a background worker.
- A stale registry entry (a `task_started` with no `task_notification`): the cost is the same, 120 s at most. The extension changes the length of the silence limit only. It never waits on worker liveness without a limit.
- Stalled runs with no worker, or with only a background Bash, keep the 60 s limit.
- Triage stays clear: the kill message states the limit that applied and the number of running workers. "Hung while a worker was registered" is therefore separate from a plain stall.
- What stays open: a real worker silence longer than 180 s still kills the run. This is a residual flake, not a hidden hang. Raise the limit only if there is gap evidence.

## Stage Report: ideation

- DONE: From the keep-moving probe evidence in repair-sonnet-live-flakes.md and the harness code (drainToExit, quietBudgetDefault), name the exact rule that kills a run and say whether the harness can tell a live background worker from a stalled run (process tree, stream events, or host task events), with a deterministic spike that proves the signal.
  Rule: `drainToExit` (`streamwatch_test.go:396`) kills when 60 s pass with no new stream line and the process is alive. The Claude runner (`:731`) then fails before it grades. Host `task_started`/`task_notification` events show that a worker is alive (seen in all 3 recorded Claude fixtures, team mode included). No signal shows progress during a thinking turn: `task_progress`, partial messages, the process tree and transcript files were each checked and rejected. Fake-clock spike on a throwaway clone of `64274cef5`: 5 of 5 cases pass. Two mutations each make the cases they target fail. Fixture replay passes. The package passes.
- DONE: Propose the smallest harness change with a falsifiable AC: no kill while a dispatched worker is alive and making progress, a fast fail when every process is stalled; include an exact file list, net LOC estimate and tolerance, and say whether live runs are needed to prove it.
  In `drainToExit`, the limit becomes 180 s (`quietBudgetDispatchClose`) while the host registry has a non-Bash task open, and stays 60 s otherwise. ACs 1 to 5 each name a falsifying edit. Net +140 across 3 files (`streamwatch_test.go`, `streamwatch_unit_test.go`, `claude_live_runner_test.go`), tolerance +90 to +200 and 2 to 3 files. No live run is needed for the ACs. Grading under load comes from the retirement runs already budgeted.
- DONE: State the risk that the change hides a real hang, and how the proposal bounds it.
  A hang while a worker is registered, or a stale registry entry, costs at most 120 s more, once per run. Behind that are the host watchdog (300 s or more) and the 600 s background stop. Runs with no worker keep 60 s. The kill message names the limit that applied and the worker count. A real silence longer than 180 s remains a flake.

### Summary

The kill rule is pure stream silence. The harness cannot see progress inside a silent thinking turn, but the host's `task_started`/`task_notification` pair tells it reliably that a worker is alive. The proposal raises the silence limit from 60 s to the sanctioned 180 s only while such a worker is registered, so the limit still applies. A fake-clock spike replayed b1's gaps without a kill, still killed stalled runs with no worker at 60 s, and showed that a `subagent_type` filter would miss team-mode workers. The throwaway clone was removed, and no live runs were made.
