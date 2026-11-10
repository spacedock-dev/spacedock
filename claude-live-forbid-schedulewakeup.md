---
id: ef4p1sbvkk9zxjxm0ez9npdr
title: Claude live harness forbids ScheduleWakeup so haiku runs exit when the work is done
status: backlog
source: "Captain directive, 2026-11-10: the fix is a test-harness injection on haiku that forbids the tool, not a first-officer instruction."
started:
completed:
verdict:
score: 0.6
worktree:
issue:
pr:
mod-block:
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

- That a pending wakeup is what holds `claude -p` open rests on the 3-of-3 versus 0-of-41 correlation; it has not been reproduced in isolation.
- Whether the three killed journeys pass their assertions once they exit cleanly.
- Whether the restriction applies to the haiku lane only or to every Claude live lane. Sonnet does not schedule wakeups, so an all-lane restriction would need no per-model branch.

## Outcome

The Claude live harness withholds or disallows `ScheduleWakeup` for haiku runs, a deterministic test pins the injected restriction in the built `claude` argument list, and a haiku live run shows no "no stream progress" kill that follows a `ScheduleWakeup` call. First-officer skill text is unchanged.
