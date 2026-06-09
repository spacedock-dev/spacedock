---
id: 6a36f26kab5nfw8wg8gcy3sw
title: Frontdoor / launch-stream-shape changes must gate on a live e2e run (cmx process gap)
status: backlog
source: FO (2026-06-09) 0.20.0 flip — cmx's launch banner broke the claude live journey-metrics parse and slipped through because cmx's validation was offline-only.
started:
completed:
verdict:
score:
worktree:
issue:
---

`cmx` (frontdoor-launch-banner-ux) shipped a stderr launch banner that the claude live runner folds into the stream-json pipe, breaking `journeymetrics.ParseClaudeJSONL` (first non-JSON line). It was caught only when the flip's live e2e ran — AFTER cmx merged — because cmx's validation was offline banner-render tests; the live e2e that catches runtime-stream-shape regressions was never run on it. (The parser was hardened to skip non-JSON lines during the flip — 6accd320 — so the symptom is fixed; this is the PROCESS guard.)

Establish that changes to the front door / launch path / host stream shape (`internal/cli/frontdoor.go`, `host_exec.go`, banner, install argv) gate on a live e2e run — or a targeted live front-door smoke — before merge, so this regression class is caught pre-merge rather than at the next release.
