---
title: Repair Pi live harness — common-journey parallelism, custom-provider model mirror, and slow-model timeout
status: backlog
sprint: test-behavior-completeness
group: common-evidence
source: "fo write-scope review of direct test edits"
id: pnc09c4hz4wbzyegh3pnpb1d
---

## Problem

The Pi live suite has three harness gaps that block running the common journeys with a custom slow-thinking provider/model such as `lunaroute/glm-5.2-vision-background:max`:

1. `shared_live_runner_test.go` calls `t.Parallel()` only when `SPACEDOCK_LIVE_RUNTIME == "claude"`. Pi journeys never parallelize, so `-parallel 2` is a no-op and all 17 common journeys run sequentially even when the operator requests 2-at-a-time.
2. `seedPiLiveAuth` copies only `auth.json` into the isolated Pi home. Custom providers (e.g. `lunaroute`) declare their models in `models.json`, not `auth.json`, so an auth-only isolated home cannot resolve `SPACEDOCK_PI_LIVE_CHILD_MODEL` for them — Pi reports `Model ... not found`.
3. The Pi per-run cap is a fixed `12*time.Minute` (common) / `10*time.Minute` (smoke). A slow `:max`-thinking model taking minutes per turn exceeds the cap on multi-dispatch journeys, producing a timeout `--- FAIL` instead of a graded result.

## Value

A Pi operator can run the full common-journey suite against a custom slow-thinking model with `-parallel 2` actually fanning out, the model resolves in the isolated home, and multi-dispatch journeys complete under a raised per-run cap. This unlocks running the 17 common journeys (and the 4 known Pi XFAIL journeys) on `lunaroute/glm-5.2-vision-background:max` with parallelism.

## Acceptance criteria

- AC-1 (parallelism): `shared_live_runner_test.go` calls `t.Parallel()` for `SPACEDOCK_LIVE_RUNTIME == "pi"` as well as `"claude"`, so `-parallel 2` fans out Pi common journeys two at a time. Codex stays sequential (unchanged).
- AC-2 (model mirror): `seedPiLiveAuth` copies `~/.pi/agent/models.json` into the isolated Pi home when present, so custom-provider models declared in `models.json` resolve. `auth.json` and `models.json` are both copied with `0o600`.
- AC-3 (slow-model timeout): a `piLiveRunTimeout(dflt)` helper reads `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` and returns `dflt` when unset or invalid. The common `run` and the smoke `runPiLiveCommand` use `piLiveRunTimeout(12*time.Minute)` / `piLiveRunTimeout(10*time.Minute)` instead of the bare constants; the timeout-fatal message names the env var and the default.
- AC-4: The doc `docs/runtime-live-ci.md` describes Pi common journeys as running at most two at a time, documents `SPACEDOCK_PI_LIVE_CHILD_MODEL` with the `provider/model:thinking` form, documents the `models.json` copy and the `Model ... not found` check, and documents `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` with the "make the outer go test -timeout longer than the per-run cap" guidance.
- AC-5: `gofmt`, `go vet -tags live ./internal/ensigncycle`, `go build -tags live ./internal/ensigncycle`, and the offline `PiLiveEnv|PiIntercom|TestPiLive` unit tests pass. The existing `TestPiLiveEnvDropsForeignRuntimeMarkers` and `TestPiLiveEnvScrubsAmbientPiSubagentMarkers` assertions still hold (the `models.json` copy is additive and does not change env scrubbing).

## Notes

- These three changes are independent but bundled because they are all Pi runner infrastructure enabling the same operator scenario (custom slow-thinking model, parallel common journeys).
- Do not touch `assertConflictOwnerHandoff` or any XFAIL binding — those are owned by the separate `fix-conflict-owner-handoff-xfail-grading` task and the runtime-specific repair entities.
- Pre-existing uncommitted edits to `pi_live_runner_test.go`, `pi_shared_live_runner_test.go`, and `docs/runtime-live-ci.md` in the working tree implement parts of AC-2/AC-3/AC-4; the worker should adopt or re-derive them from this spec rather than assume they are absent.
