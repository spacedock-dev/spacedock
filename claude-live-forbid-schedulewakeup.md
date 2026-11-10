---
id: ef4p1sbvkk9zxjxm0ez9npdr
title: Claude live harness forbids ScheduleWakeup so haiku runs exit when the work is done
status: ideation
source: "Captain directive, 2026-11-10: the fix is a test-harness injection on haiku that forbids the tool, not a first-officer instruction."
started:
completed:
verdict:
score: 0.6
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:ef4p1sbvkk9zxjxm0ez9npdr:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ef4p1sbvkk9zxjxm0ez9npdr-backlog-1
              briefing:
                id: briefing:ef4p1sbvkk9zxjxm0ez9npdr:backlog:attempt-1:revision-1
                digest: sha256:ecbd5ac3e78d627a49fcd0be234df9b0515c0f45e3a56dc8c7d4e50863ec533b
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:ef4p1sbvkk9zxjxm0ez9npdr:backlog:1
                briefing: briefing:ef4p1sbvkk9zxjxm0ez9npdr:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-10T18:34:22.537219Z"
                decision: approve
                reason: Seed states the problem with run evidence, a bounded harness-only scope, and the captain's chosen mechanism; ready for ideation.
                conn:
                    quote: for haiku - dispatch this with the disallowed tool change the remaining failing marked as xpass. create PR onto the current stack and we can run ci
                    source: captain chat, 2026-11-10
              application:
                target-stage: ideation
                state: consumed
---

On the claude-haiku live lane, journeys that complete their work are killed as hung. The haiku first officer calls the host tool `ScheduleWakeup` as a fallback heartbeat while waiting on background workers; after the work finishes and the final turn ends, the pending wakeup keeps the `claude -p` process alive and silent, and the harness kills it after its 60-second no-progress budget, before grading.

## Problem

- Run 42865474062 (`ci/go-126-workflows` @ `eef2d1302`, model `claude-haiku-5-5`), attempts 1 and 2, 2026-11-10. Three journeys were killed with "made no stream progress within 1m0s (no-progress quiet budget)": `owned-conflict-owner-handoff` and `smallest-sufficient-mechanism` in attempt 1, `rejection-flow` in attempt 2.
- Each killed stream contains `ScheduleWakeup` calls with a delay (60s once, 1200s otherwise), a final `end_turn` with no background tasks, then 60 seconds of silence. `ScheduleWakeup` scheduling calls appear in exactly these three streams and in none of the other 41 haiku streams across the two attempts.
- The same journeys in the other attempt never call `ScheduleWakeup` and exit within about a second of the last event.
- Three green sonnet runs (66 streams: 41838996310, 42848106064, 42849771310) contain no scheduling call; one stream calls `ScheduleWakeup` with `stop: true` only.
- No file under `skills/` mentions `ScheduleWakeup`; it comes from the Claude Code host tool list.
- The kill is `drainToExit` in `internal/ensigncycle/streamwatch_test.go`, invoked with `quietBudgetDefault` from `internal/ensigncycle/claude_live_runner_test.go`, which builds the `claude` argument list (`--permission-mode bypassPermissions`, stream-json, model).

## Unknowns

- That a pending wakeup is what holds `claude -p` open rests on a correlation: 4 kills across three CI attempts, each in a stream with a scheduling call, and no scheduling call in any stream that exited cleanly. It has not been reproduced in isolation.
- Whether the restriction applies to the haiku lane only or to every Claude live lane. Sonnet does not schedule wakeups, so an all-lane restriction would need no per-model branch.

## Evidence with the restriction applied

A local run on `main` @ `5785dbed7` with `--disallowedTools ScheduleWakeup` appended to the `claude` argument list for haiku (Claude Code 2.1.317, 2026-11-10): `ScheduleWakeup` was absent from the session tool list, none of 22 streams called it, and no journey was killed for stream silence. `owned-conflict-owner-handoff` and `smallest-sufficient-mechanism` passed.

Failures that remain on haiku once the kill is gone:

- `rejection-flow`: graded `rejection-worker-topology` — the reuse/implementation dispatch is followed by reuse/validation before its own completion. CI attempts 2 and 3 killed this journey before grading.
- `TestLiveBareReachable`: no bare-shaped `Agent()` call observed, in the one local run. It passed in all three CI attempts.
- `ac-value-reanchor`: failed in CI attempt 1 (the first officer recorded `revise` and stopped at `validation`), passed in attempts 2 and 3 and locally.

Haiku xfail bindings that passed in all four runs: `default-headless-gate-stop`, `withdrawn-gate-recovery`, dispatch recovery selected-bare. `keep-moving-posture` failed once and passed three times.

## Outcome

The Claude live harness disallows `ScheduleWakeup` for haiku runs, and a deterministic test pins the injected restriction in the built `claude` argument list. Haiku journeys that still fail with the restriction applied carry `claude-haiku` xfail bindings, so the haiku lane reports them as known gaps. A haiku live run in CI shows no "no stream progress" kill that follows a `ScheduleWakeup` call. First-officer skill text is unchanged.

The change ships as a pull request stacked on `ci/go-126-workflows` (PR #849).
