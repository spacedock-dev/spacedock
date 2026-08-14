# BACKLOG GATE — Live CI API-error log capture (`r5`)

Recommendation: **APPROVE and dispatch ideation.**

## Capability and value

When a live-lane scenario dies of stream silence, the archived evidence (stream
jsonl + step log) shows only the silence. API-level errors, retry attempts, and
connection lifecycle events are not captured, so a stall cannot be classified:
transient weather (re-run and move on) vs a CLI retry/hang defect (file upstream).
Every such failure costs a full lane re-run and teaches nothing. This task
captures the host CLI's API-error/retry/debug output as a per-scenario artifact;
the harness kill message points at the captured log.

This is the tooling I wished for twice this session: the `gate-guardrail` flake
and the `owned-conflict` model-garbage runs both died mid-stream with no stack
trace (the temp Pi home was cleaned up), leaving the failure unclassifiable.

## Binding boundaries

- Diagnostic capture only. Does not change any journey's exercise, fixture, or
  assertion; does not gate on the captured output.
- Adds a per-scenario artifact (size/noise cost to characterize in ideation).
  Ideation determines what each CLI actually exposes (read the docs/schema first
  — usage presence is not existence evidence), whether capture is always-on or
  armed only for the kill path, and the artifact size cost.
- No new runtime, result format, or CI lane. The capture is additive evidence.

## Proof direction

Ideation determines the CLI surfaces (claude debug/verbose flags, codex/pi
equivalents if cheap), the artifact shape, and always-on vs kill-path-only
capture. Implementation proves a stalled scenario leaves an artifact showing the
API request lifecycle around the stall, with negligible cost on green runs, and
the harness kill message names the artifact path.

## Decision ask

Approve this diagnostic capture for ideation, or revise/hold with a concrete
boundary.
