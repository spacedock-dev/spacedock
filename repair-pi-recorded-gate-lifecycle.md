---
title: Repair the Pi recorded-gate-lifecycle journey
status: backlog
source: "CI run 31770740214 (PR #685 pi-live, model openai/gpt-5.6-luna:max): TestLiveCommonRecordedGateLifecycle FAIL observed=[recorded-gate-lifecycle-violation], 'Blocked at the validation gate. The required committed reference is missing.'"
score: 0.85
sprint: pi-live-completeness
sprint-readiness: ready
group: pi-live-followup
id: gcmfwfjd9735b58sbzw7xsb8
---

## Problem

The Pi FO does not complete the `recorded-gate-lifecycle` journey cleanly. The
strict oracle rejects it with `observed=[recorded-gate-lifecycle-violation]`:
"Blocked at the validation gate. The required committed reference is missing."
The FO reaches the validation gate but the retained reference the lifecycle
requires is absent, so the recorded-gate lifecycle is not durably complete on Pi.

This is an ordinary lane FAIL on Pi: `recorded-gate-lifecycle` is XFAIL-bound only
for `claude-opus` (`66dpwxgvsxt7cbxhmgvt3qp4`). There is no `liveXFail("pi",…)`
or `liveTODO("pi",…)` binding, so on Pi the journey is expected to PASS. The CI
`pi-live` lane is correctly red on an unowned, real Pi conduct gap.

## Visible value

A Pi operator runs a delegated-authority recorded gate and the lifecycle binds,
records, commits, and consumes exactly once before successor dispatch, with the
required committed reference present at the validation gate. Measured against
baseline: before this fix, the Pi `recorded-gate-lifecycle` journey FAILs with
`recorded-gate-lifecycle-violation` (missing committed reference); after, the same
run completes the recorded gate lifecycle to a PASS.

## Evidence

- CI: `Runtime Live E2E` run 31770740214, `pi-live` job, model
  `openai/gpt-5.6-luna:max`; `TestLiveCommonRecordedGateLifecycle` FAIL,
  `observed=[recorded-gate-lifecycle-violation]`, "Blocked at the validation
  gate. The required committed reference is missing." Recorded in the pnc
  pi-live log (17 common journeys ran; this was 1 of 3 failures, the other two
  being the known `default-headless-gate-stop` gap (nta) and the
  `keep-moving-posture` XPASS).

## Out of scope

- The `claude-opus` XFAIL binding (`66d`). This task owns only the Pi conduct.
- Shared XFAIL policy, the assert, or the fixture.
- Sonnet, Codex, or any other runtime's behavior on this journey.
- A new runtime, fixture, result format, or CI lane.

## Acceptance criteria

**AC-1 (VALUE) — The exact Pi recorded-gate-lifecycle target passes.**

Verified by: the focused live Pi `TestLiveCommonRecordedGateLifecycle` target
exits successfully — the FO binds, records, commits, and consumes the delegated
authority exactly once before successor dispatch, with the required committed
reference present at the validation gate. Baseline: the current
`recorded-gate-lifecycle-violation` FAIL.

**AC-2 — The Pi binding stays honest.**

Verified by: no `liveXFail("pi",…)` is added to mask the gap; the journey reaches
PASS by completing the lifecycle, not by weakening the assertion. A temporary
`liveTODO("pi",…)` if needed names this active task as owner and is removed on PASS.

**AC-3 — Other runtimes and the shared assert are preserved.**

Verified by: the `claude-opus` XFAIL binding and the shared assert are unchanged;
Sonnet and Codex behavior on this journey is unaffected.

**AC-4 — Offline and required-lane checks pass.**

Verified by: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
`go build -tags live ./internal/ensigncycle`, `go test ./...`, and
`go test ./... -race` pass; the Pi live lane passes the focused target.

## Test plan

Use focused offline gate and terminalization controls first. Use one exact Pi
`recorded-gate-lifecycle` target sequence only when Pi work is authorized.
Preserve all Sonnet and Codex behavior and the shared assert.

## Notes

- Coordinate with `pi-delegated-gate-continuation-reliability` (9w, group: gate),
  which owns the Pi recorded-gate journey under delegated conn — the missing
  committed reference may share a root cause in the gate-record/dispatch seam.
- Filed from the pnc pi-live run (31770740214), which surfaced this gap because
  pnc's parallelism change let all 17 common journeys run (vs -failfast stopping
  at the first failure on non-parallel branches).
