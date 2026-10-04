---
title: Repair the Pi recorded-gate-lifecycle journey
status: implementation
source: "CI run 31770740214 (PR #685 pi-live, model openai/gpt-5.6-luna:max): TestLiveCommonRecordedGateLifecycle FAIL observed=[recorded-gate-lifecycle-violation], 'Blocked at the validation gate. The required committed reference is missing.'"
score: 0.85
sprint: pi-live-completeness
sprint-readiness: ready
group: pi-live-followup
id: gcmfwfjd9735b58sbzw7xsb8
started: 2026-10-04T04:35:04Z
worktree: .worktrees/spacedock-ensign-repair-pi-recorded-gate-lifecycle
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

## Stage Report: implementation

- DONE: Fix the officer instruction so the gate-prepare reference path is composed from the observed workflow root or the committed entity path, and is never reproduced from prose in the prompt.
  `skills/fo-gate-lifecycle/SKILL.md` Prepare step now reads "Compose the selected-source `<R>/<P>` (observed workflow root or committed entity path); never retype prose. Use absolute if cwd/state spellings differ." This restores the `<R>/<P>` composition that commit 6060c382a replaced with "Supply absolute, launch-cwd-relative, or state-relative judgment paths" (the wording that invited the hand-typed path). The 7700-byte component cap holds: file is 7695 B.
- DONE: Prove it: the focused checks, then a local live run of default-headless-gate-stop using SPACEDOCK_LIVE_RUNTIME=pi and a distinct SPACEDOCK_LIVE_ARTIFACT_DIR.
  After-fix run (base 05bdfa7da; branch later rebased onto the sibling's amended tip a2166585a — binding removed, skill fixed): `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/live-refpath SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 45m -run '^TestLiveCommonDefaultHeadlessGateStop$' ./internal/ensigncycle -v` -> `--- PASS: TestLiveCommonDefaultHeadlessGateStop (551.12s)`, no XFAIL and no XPASS alert (the binding is gone, so an XPASS is impossible), empty observed set. The run's `command.log` shows `gate prepare ... --artifact .../entity-snapshot.md --reference .../recorder-contract.md` exit=0, then `state commit` exit=0 (`state-head d70d84bf`).
- DONE: When that journey passes locally with an exact passing pair, clear its XFAIL binding and record the evidence in the entity.
  Removed `[]liveJourneyGap{liveXFail("pi","gcmfwfjd9735b58sbzw7xsb8")}` -> `nil` for `TestLiveCommonDefaultHeadlessGateStop` in `internal/ensigncycle/shared_live_runner_test.go` and dropped its now-obsolete comment; committed with the skill fix in bb3265618.
- DONE: Do not weaken the shared assertion, the binary's refusal, or the no-retry rule.
  Shared assertions (`assertGateHeld`, `assertRecordedGateHoldLog`, `assertImplementationWorkerLifecycle`) are byte-unchanged; `internal/gates/prepare.go` is untouched, so the fail-closed refusal of an unreadable/missing selected source stands; the skill's single-prepare/no-retry sentence is untouched.
- DONE: Focused offline checks.
  `gofmt -l` clean; `go vet -tags live ./internal/ensigncycle` clean; `go test -tags live -run '^$' ./internal/ensigncycle` compiles; `go test ./internal/contractlint/... ./internal/ensigncycle/... ./internal/gates/... ./internal/cli/... -run '...' -count=1` green (ensigncycle offline 308s, contractlint 1.3s, gates/cli gate-prepare green). Per the FO scope, the repository-wide `go test ./...`/`-race` was not run.
- SKIPPED: none.

### Summary

This is hardening against a transcription slip, not the repair of a reproducible fault. The before-fix local run (peer base fb4428e9c, binding present, skill unfixed) did NOT reproduce the missing-reference fault at all: `gate prepare` succeeded (exit=0, `state=open`, correct `--artifact`/`--reference` paths, committed) and the sole red was the separate observer defect `implementation-worker-not-dispatched` (owned by mk72bnt1b5hsp9sfv83979xs), producing `XFAIL pi/default-headless-gate-stop observed=[implementation-worker-not-dispatched]`. The reference-path failure recorded from CI kept the random suffix while dropping one word, which is only possible when the path was written from the prompt prose rather than read from the boot record; the instruction change makes that harder without claiming a reproduced fault. Because the headless journey passes only once the observer fix is present, my branch is stacked on the sibling tip a2166585a (peer fb4428e9c + `2d5bd25f6` bg_wait/native-completion credit + `a2166585a`); the branch was first based on the sibling's then-tip 05bdfa7da and rebased onto a2166585a when the sibling amended its replay test, and the observer fix in `claude_runtime_helpers_test.go` is byte-identical across the two, so the live PASS holds. My commit on the final base is bb3265618. Residual risk: the journey is stochastic — the after-fix run needed a clean artifact directory (a killed prior run in the same dir made the harness see two root sessions and fail before the journey); a future rerun must start from an empty `SPACEDOCK_LIVE_ARTIFACT_DIR`.
