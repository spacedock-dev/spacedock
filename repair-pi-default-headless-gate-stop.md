---
title: Repair the Pi implementation-worker-not-dispatched conduct (default-headless-gate-stop + auto-continue-after-implementation)
status: backlog
source: "CI run 31747645316 (PR #682, pi-live job) + 2026-08-13-02 Pi debrief; reproduced locally on lunaroute/glm-5.2-vision-background:max"
score: 0.85
sprint: pi-live-completeness
sprint-readiness: ready
group: pi-live-followup
id: ntarrp8jp5h34g6528d66kbe
---

## Problem

The Pi FO does not dispatch the preceding-stage implementation worker before
stopping at the first human gate in the `default-headless-gate-stop` journey. A
headless launch without decision authority is required to dispatch and complete
the preceding-stage worker, then present the first human gate and stop open.
On Pi, the FO reaches the `validation` gate, presents a `Recorded Gate Task —
validation` review recommending approve, and stops — without ever dispatching the
implementation worker. The durable assertion fails with
`observed=[implementation-worker-not-dispatched]`.

This is an ordinary lane FAIL on Pi: `default-headless-gate-stop` is XFAIL-bound
only for `claude-sonnet` (owner `kk`, `commit-sonnet-gate-before-presentation`).
There is no `liveXFail("pi",…)` or `liveTODO("pi",…)` binding, so on Pi the
journey is expected to PASS. The CI `pi-live` lane is correctly red on an
unowned, real Pi conduct gap.

## Visible value

A Pi operator runs a headless launch without the conn and the FO dispatches and
completes the implementation worker, then presents the first human gate and stops
open — the contract the scenario encodes. Measured against baseline: before this
fix, the Pi `default-headless-gate-stop` journey FAILs with
`implementation-worker-not-dispatched`; after, the same run dispatches the worker,
completes it, presents the gate, and stops open (PASS). Two independent models
reproduce the same failure, so the baseline is stable across the Pi platform.

## Evidence (two-model reproduction)

- CI: `Runtime Live E2E` run 31747645316, `pi-live` job, step "Run live Pi common
  journeys" failed; `TestLiveCommonDefaultHeadlessGateStop` (551.71s),
  `observed=[implementation-worker-not-dispatched]`. Model
  `openai-codex/gpt-5.6-luna:max`. Front-door smoke passed. Recorded in
  `_debriefs/2026-08-13-02-pi-gpt-5-6-luna.md`.
- Local (same session): `lunaroute/glm-5.2-vision-background:max`,
  `TestLiveCommonDefaultHeadlessGateStop` (363.78s),
  `observed=[implementation-worker-not-dispatched]`. Identical final message
  shape (presents validation gate, recommends approve, stops without dispatching
  the implementation worker).

Same journey, same `observed` code, same final-message shape, two different
models on the Pi platform. Strong evidence this is a Pi-platform FO conduct gap,
not a model-quality issue.

## Out of scope

- The `claude-sonnet` XFAIL binding (`kk`). This task owns only the Pi conduct.
- Shared XFAIL policy, the assert, or the fixture. `assertGateHeld` is unchanged.
- Sonnet, Codex, or any other runtime's behavior on this journey.
- A new runtime, fixture, result format, or CI lane.

## Acceptance criteria

**AC-1 (VALUE) — The exact Pi default-headless-gate-stop target passes.**

Verified by: the focused live Pi `TestLiveCommonDefaultHeadlessGateStop` target
exits successfully — the FO dispatches and completes the preceding-stage
implementation worker, then presents the first human gate and stops open. The
durable assertion `assertGateHeld` passes with no `implementation-worker-not-dispatched`
code. Baseline: the current two-model-reproduced FAIL.

**AC-2 — The Pi binding stays honest.**

Verified by: no `liveXFail("pi",…)` is added to mask the gap; the journey reaches
PASS by the FO dispatching the worker before stopping, not by weakening the
assertion. (If a temporary `liveTODO("pi",…)` is needed to keep the lane green
during repair, it names this active task as owner and is removed on PASS.)

**AC-3 — Other runtimes and the shared assert are preserved.**

Verified by: the `claude-sonnet` XFAIL binding and `assertGateHeld` are
unchanged; Sonnet and Codex behavior on this journey is unaffected.

**AC-4 — Offline and required-lane checks pass.**

Verified by: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
`go build -tags live ./internal/ensigncycle`, `go test ./...`, and
`go test ./... -race` pass; the Pi live lane passes the focused target.

## Test plan

Use focused offline gate and terminalization controls first. Use one exact Pi
`default-headless-gate-stop` target sequence only when Pi work is authorized.
Preserve all Sonnet and Codex behavior and the shared assert.

## Notes

- Coordinate with `commit-sonnet-gate-before-presentation` (kk), which owns the
  claude-sonnet XFAIL binding on the same journey; the Pi conduct gap may share a
  root cause in the gate-presentation/dispatch seam, but this task owns only the
  Pi result.
- Filed because the 2026-08-13-02 Pi debrief recorded the CI failure but filed
  nothing ("None newly filed in this session. The live journey failures overlap
  parallel repair work."). This entity gives the gap an owner.
