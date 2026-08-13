---
title: "Live CI captures API-error/retry logs so a stalled stream is diagnosable, not just dead"
status: backlog
group: tooling
source: "2026-07-02 session: PR #461's claude-live sonnet lane failed the filing scenario on a 60s no-progress kill whose transcript tail ends mid-thinking (stream died mid-generation, work already complete on disk); the same window produced two explicit 'API Error: Connection closed mid-response' failures in an interactive subagent. The rerun went green with zero code change. Root cause was NOT determinable from CI: nothing captured distinguishes provider-side API weather / runner networking from the one actionable variant — the claude CLI hanging instead of retrying after a dropped stream. Captain direction: 'in any case, we should have api error logs etc, if we are not currently capturing that.'"
id: r5y6qjr10k4m3gw9w5p3b2vj
sprint: pi-ux
---

## Problem
When a live-lane scenario dies of stream silence, the archived evidence (stream jsonl + step log) shows only the silence. API-level errors, retry attempts, and connection lifecycle events are not captured, so a stall cannot be classified: transient weather (re-run and move on) vs a CLI retry/hang defect (file upstream, work around). Every such failure costs a full lane re-run and still teaches nothing.

## Desired direction (for ideation to refine)
The live runners capture the host CLI's API-error/retry/debug output as a per-scenario CI artifact alongside the existing stream jsonl — whatever channel the claude CLI (and codex/pi equivalents, if cheap) exposes: debug/verbose flags, error log files, env-var-enabled logging. On a no-progress kill, the harness failure message points at the captured log. Ideation determines what the CLI actually exposes (read the docs/schema first — usage presence is not existence evidence), the artifact size/noise cost, and whether capture is always-on or armed only for the kill path.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- A live scenario killed for no-progress leaves an artifact that shows the API request lifecycle around the stall (attempts, errors, retries or their absence) — demonstrated by inspection on a real or induced stall.
- The capture adds negligible cost to green runs (measured artifact size / wall-clock delta), and the harness kill message names the artifact path.
- Related, separate candidate recorded (not this task): grade-on-outcome-state when the stream dies after the deliverable landed — the filing kill discarded a satisfiable pass.
