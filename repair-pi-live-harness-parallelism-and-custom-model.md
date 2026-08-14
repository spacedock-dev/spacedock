---
title: Repair Pi live harness — common-journey parallelism, custom-provider model mirror, and slow-model timeout
status: validation
source: "FO write-scope review, 2026-08-13: direct test edits routed through filing"
score: 0.85
sprint: pi-ux
sprint-readiness: ready
group: pi-live-followup
id: pnc09c4hz4wbzyegh3pnpb1d
gates:
    version: 1
    records:
        - id: gate:pnc09c4hz4wbzyegh3pnpb1d:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:pnc09c4hz4wbzyegh3pnpb1d-backlog-1
              briefing:
                id: briefing:pnc09c4hz4wbzyegh3pnpb1d:backlog:attempt-1:revision-1
                digest: sha256:54bc1775cf5a638d29455dd68415d04c558715bb1f6009f0e6d93964c2e55b5f
                request-digest: sha256:f166ecca01e94cb9f9f8861a808d8c9c2ac27140ac9080f31b7ee45578430e2d
                room-ref: ./repair-pi-live-harness-parallelism-and-custom-model/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pnc09c4hz4wbzyegh3pnpb1d:backlog:1
                briefing: briefing:pnc09c4hz4wbzyegh3pnpb1d:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:40:14.7204Z"
                decision: approve
                reason: Captain approved backlog gate; advance to ideation for the Pi-harness enablement bundle.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:pnc09c4hz4wbzyegh3pnpb1d:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:pnc09c4hz4wbzyegh3pnpb1d-ideation-1
              briefing:
                id: briefing:pnc09c4hz4wbzyegh3pnpb1d:ideation:attempt-1:revision-1
                digest: sha256:0f88d0082b23330340b7234b5741b42d0828d42173d5212b5a26c6f0443814da
                request-digest: sha256:7d6beada4a235da00201933a3bdcc516ec29dc4345d4d00adfcd2c5286c63030
                room-ref: ./repair-pi-live-harness-parallelism-and-custom-model/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pnc09c4hz4wbzyegh3pnpb1d:ideation:1
                briefing: briefing:pnc09c4hz4wbzyegh3pnpb1d:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-14T01:52:17.337776Z"
                decision: approve
                reason: Captain approved ideation gate; dispatch implementation for the Pi-harness bundle.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:pnc09c4hz4wbzyegh3pnpb1d:validation
          stage: validation
          attempts:
            - id: gate-attempt:pnc09c4hz4wbzyegh3pnpb1d-validation-1
              briefing:
                id: briefing:pnc09c4hz4wbzyegh3pnpb1d:validation:attempt-1:revision-1
                digest: sha256:6e39275c667140ebb5ab71731a576c62df243a793414a5847377cfb8bcc25f91
                request-digest: sha256:4ae1a20e637685989b1a286e580704d7b5f4da9b7780cca978e33724aff7a746
                room-ref: ./repair-pi-live-harness-parallelism-and-custom-model/review/validation/briefing-1
started: 2026-08-14T01:55:02Z
worktree: .worktrees/spacedock-ensign-repair-pi-live-harness-parallelism-and-custom-model
---

## Problem

The Pi live harness has three gaps that block running the common journeys against
a custom slow-thinking provider/model such as `lunaroute/glm-5.2-vision-background:max`:

1. `shared_live_runner_test.go` calls `t.Parallel()` only when
   `SPACEDOCK_LIVE_RUNTIME == "claude"`. Pi journeys never parallelize, so
   `-parallel 2` is a no-op and all 17 common journeys run sequentially even when
   the operator requests 2-at-a-time.
2. `seedPiLiveAuth` copies only `auth.json` into the isolated Pi home. Custom
   providers (e.g. `lunaroute`) declare their models in `models.json`, not
   `auth.json`, so an auth-only isolated home cannot resolve
   `SPACEDOCK_PI_LIVE_CHILD_MODEL` for them — Pi reports `Model ... not found`.
3. The Pi per-run cap is a fixed `12*time.Minute` (common) / `10*time.Minute`
   (smoke). A slow `:max`-thinking model taking minutes per turn exceeds the cap
   on multi-dispatch journeys, producing a timeout `--- FAIL` instead of a graded
   result.

These are harness enablement gaps, not product repairs. They do not change any
journey's exercise or durable assertion; they let an operator run the existing
suite against a custom slow model with parallelism.

## Visible value

A Pi operator runs the full common-journey suite against a custom slow-thinking
model with `-parallel 2` actually fanning out two journeys at a time, the model
resolves in the isolated home, and multi-dispatch journeys complete under a raised
per-run cap. Measured against baseline: before this fix, the suite either skips
(`Model ... not found`) or runs all 17 sequentially under a 12m cap that a `:max`
multi-dispatch journey exceeds; after, two journeys run concurrently and a
raised `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` lets them complete to a graded result.

## Out of scope

- Any journey's exercise, fixture, or durable assertion. This task changes only
  the Pi runner harness and the shared `t.Parallel()` gate.
- The `assertConflictOwnerHandoff` assert or any XFAIL binding — owned by
  `fix-conflict-owner-handoff-xfail-grading` (rzr) and the runtime-specific
  repair entities.
- The `0a` CI workflow cadence (archived, done) — this task does not change
  GitHub Actions or the merge-evidence policy.
- A new registry, result artifact, host-switch table, or lifecycle model.

## Acceptance criteria

**AC-1 (VALUE) — A custom slow `:max` model runs the common suite two at a time to graded results.**

Verified by: `SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_CHILD_MODEL='lunaroute/glm-5.2-vision-background:max' SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES=40 go test ... -run '^TestLiveCommon' -parallel 2` fans out two journeys concurrently (two `=== CONT` lines together) and the journeys complete to XFAIL/XPASS/PASS graded results rather than `Model ... not found` or a 12m timeout `--- FAIL`. The baseline is the current sequential-only, 12m-capped, auth-only behavior.

**AC-2 — Pi common journeys call `t.Parallel()`; Codex stays sequential.**

Verified by: `shared_live_runner_test.go` calls `t.Parallel()` for
`SPACEDOCK_LIVE_RUNTIME == "pi"` as well as `"claude"`. Codex remains sequential
(no `t.Parallel()` for `"codex"`).

**AC-3 — The isolated Pi home mirrors `models.json` for custom providers.**

Verified by: `seedPiLiveAuth` copies `~/.pi/agent/models.json` into the isolated
Pi home when present (alongside `auth.json`, both `0o600`), so custom-provider
models declared in `models.json` resolve.

**AC-4 — The per-run cap is overridable for slow models.**

Verified by: a `piLiveRunTimeout(dflt)` helper reads
`SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` and returns `dflt` when unset or invalid;
the common `run` and the smoke `runPiLiveCommand` use it instead of the bare
constants; the timeout-fatal message names the env var and the default.

**AC-5 — The doc describes the Pi parallel and slow-model surface.**

Verified by: `docs/runtime-live-ci.md` describes Pi common journeys as running
at most two at a time, documents `SPACEDOCK_PI_LIVE_CHILD_MODEL` in
`provider/model:thinking` form, documents the `models.json` copy and the
`Model ... not found` check, and documents `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES`
with the "make the outer `go test -timeout` longer than the per-run cap" guidance.

**AC-6 — Offline checks and existing env-scrub assertions pass.**

Verified by: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
`go build -tags live ./internal/ensigncycle`, and the offline
`PiLiveEnv|PiIntercom|TestPiLive` unit tests pass. `TestPiLiveEnvDropsForeignRuntimeMarkers`
and `TestPiLiveEnvScrubsAmbientPiSubagentMarkers` still hold — the `models.json`
copy is additive and does not change env scrubbing.

## Test plan

Use the offline `PiLiveEnv|PiIntercom|TestPiLive` unit tests first. Then one
focused Pi live run with the custom slow model and `-parallel 2` only when Pi
work is authorized. Preserve all Sonnet and Codex behavior.

## Notes

- Three independent changes bundled because they are all Pi runner infrastructure
  enabling the same operator scenario.
- Pre-existing uncommitted edits to `pi_live_runner_test.go`,
  `pi_shared_live_runner_test.go`, and `docs/runtime-live-ci.md` in the working
  tree implement parts of AC-3/AC-4/AC-5; the worker should adopt or re-derive
  them from this spec rather than assume they are absent.
- Filed by the FO after a write-scope review caught direct test edits; the
  worker re-derives the edits from this spec.

## Stage Report: implementation

- DONE: Pi common journeys call `t.Parallel()`, codex stays sequential.
  `internal/ensigncycle/shared_live_runner_test.go:44-46` — `switch runtime := os.Getenv("SPACEDOCK_LIVE_RUNTIME"); runtime { case "claude", "pi": t.Parallel() }` replaces the `== "claude"` check; codex falls through with no `t.Parallel()`.
- DONE: `seedPiLiveAuth` mirrors `~/.pi/agent/models.json` into the isolated Pi home.
  `internal/ensigncycle/pi_live_runner_test.go:240-247` — after the `auth.json` write, reads `~/.pi/agent/models.json` and writes `piHome/models.json` (0o600) when present and non-empty, so custom-provider models resolve.
- DONE: `piLiveRunTimeout(dflt)` helper reads `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES`.
  `internal/ensigncycle/pi_live_runner_test.go:280-291` — `piLiveRunTimeout` parses the env var via `strconv.Atoi` (>0) and returns `dflt` when unset/invalid; `strconv` added to imports.
- DONE: `runPiLiveCommand` uses `piLiveRunTimeout` (default 10m) with env-naming fatal message.
  `internal/ensigncycle/pi_live_runner_test.go:68` uses `piLiveRunTimeout(10*time.Minute)`; `:88` fatal message names `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES, default 10m`.
- DONE: `pi_shared_live_runner_test.go` `run` uses `piLiveRunTimeout(12*time.Minute)` with env-naming fatal message.
  `internal/ensigncycle/pi_shared_live_runner_test.go:62` uses `piLiveRunTimeout(12*time.Minute)`; `:77` fatal message names `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES, default 12m`.
- DONE: `docs/runtime-live-ci.md` documents Pi `-parallel 2`, `provider/model:thinking`, `models.json` mirror + `Model ... not found` check, and `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` guidance.
  `docs/runtime-live-ci.md:62-66` describes Pi running at most two journeys at a time; `:97-100` documents `provider/model:thinking` and the `models.json` mirror + `Model ... not found` check; `:102` documents `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` with "make the outer `go test -timeout` longer than the per-run cap"; common Pi command now uses `-parallel 2` and `-timeout 90m`.
- DONE: Offline checks pass.
  `gofmt -w internal/ensigncycle` (no diff after apply), `go vet -tags live ./internal/ensigncycle` (no output), `go build -tags live ./internal/ensigncycle` (BUILD_OK), `go test -tags live -run 'PiLiveEnv|PiIntercom|TestPiLive' ./internal/ensigncycle` → `ok … 0.542s`; `TestPiLiveEnvDropsForeignRuntimeMarkers` and `TestPiLiveEnvScrubsAmbientPiSubagentMarkers` PASS.
- DONE: No journey exercise/fixture/durable assertion, XFAIL binding, `assertConflictOwnerHandoff`, or CI lane touched.
  `git diff --stat` for the code commit names only `shared_live_runner_test.go`, `pi_live_runner_test.go`, `pi_shared_live_runner_test.go`, and `docs/runtime-live-ci.md`; no CI workflow or journey body edits.

### Summary

Implemented the three Pi live-harness enablement changes from the spec: widened the shared `t.Parallel()` gate to include `pi` (codex stays sequential), mirrored `~/.pi/agent/models.json` into the isolated Pi home alongside `auth.json` so custom-provider models resolve, and added an overridable `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` per-run cap (via `piLiveRunTimeout`) used by both the smoke and common runners with env-var-naming timeout messages. Updated `docs/runtime-live-ci.md` to describe Pi `-parallel 2`, the `provider/model:thinking` child model form, the `models.json` mirror and `Model ... not found` check, and the timeout override guidance. Offline `gofmt`/`vet`/`build` and the `PiLiveEnv|PiIntercom|TestPiLive` unit tests (including the two env-scrub assertions) pass.

## Stage Report: validation

- DONE: `shared_live_runner_test.go` calls `t.Parallel()` for `SPACEDOCK_LIVE_RUNTIME == "pi"` as well as `"claude"`; codex stays sequential.
  Read `internal/ensigncycle/shared_live_runner_test.go:44-47`: `switch runtime := os.Getenv("SPACEDOCK_LIVE_RUNTIME"); runtime { case "claude", "pi": t.Parallel() }`; `grep -n Parallel` shows only line 46 inside that switch; `grep -in 'codex.*parallel'` returns no match (codex falls through with no `t.Parallel()`).
- DONE: `seedPiLiveAuth` copies `~/.pi/agent/models.json` into the isolated Pi home (0o600) when present.
  Read `internal/ensigncycle/pi_live_runner_test.go:240-247`: after the `auth.json` write (0o600), `os.ReadFile(filepath.Join(realHome, ".pi", "agent", "models.json"))`; on nil err and non-empty trimmed bytes, `os.WriteFile(filepath.Join(piHome, "models.json"), models, 0o600)`.
- DONE: A `piLiveRunTimeout(dflt)` helper reads `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` and returns `dflt` when unset/invalid; used by both `run` (12m) and `runPiLiveCommand` (10m).
  `pi_live_runner_test.go:285-292` defines `piLiveRunTimeout(dflt)` parsing `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` via `strconv.Atoi` (n>0) and returning `dflt` otherwise; `grep -rn piLiveRunTimeout` shows `pi_shared_live_runner_test.go:62` uses `piLiveRunTimeout(12*time.Minute)` and `pi_live_runner_test.go:68` uses `piLiveRunTimeout(10*time.Minute)`; `strconv` added to imports.
- DONE: `docs/runtime-live-ci.md` documents the Pi parallel + custom-model + timeout surface.
  `docs/runtime-live-ci.md:62-64` states Claude and Pi run at most two journeys at a time, Codex sequential; `:97-100` documents the `models.json` mirror (0o600) and `Model ... not found` check plus `provider/model:thinking`; `:102` documents `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` with "make the outer `go test -timeout` longer than the per-run cap"; common Pi command at `:118` uses `-timeout 90m -parallel 2`.
- DONE: No journey exercise/fixture/durable assertion, XFAIL binding, `assertConflictOwnerHandoff`, or CI lane touched; `git diff --stat` reports the exact files/lines changed.
  `git diff --stat HEAD~1 HEAD` names only `docs/runtime-live-ci.md` (+14/-6), `internal/ensigncycle/pi_live_runner_test.go` (+26/-2), `internal/ensigncycle/pi_shared_live_runner_test.go` (+2/-2), `internal/ensigncycle/shared_live_runner_test.go` (+2/-1); `git diff --name-only | grep -iE 'journey|xfail|assertConflictOwnerHandoff|\.github|workflow|fixture'` returns no matches.
- DONE: Offline checks re-run in the worktree pass.
  `gofmt -l internal/ensigncycle` → empty; `go vet -tags live ./internal/ensigncycle` → exit 0, no output; `go build -tags live ./internal/ensigncycle` → exit 0, no output; `go test -tags live -run 'PiLiveEnv|PiIntercom|TestPiLive' ./internal/ensigncycle -v` → `ok github.com/spacedock-dev/spacedock/internal/ensigncycle 0.619s` with `TestPiLiveEnvDropsForeignRuntimeMarkers` and `TestPiLiveEnvScrubsAmbientPiSubagentMarkers` both PASS.
- DONE: Deliverable not edited during validation.
  `git status --short` in the worktree is empty (no staged or unstaged changes); validation was read-only against the implementation commit `864a3ca97`.

### Summary

Recommend PASSED. The deliverable at commit `864a3ca97` satisfies all four checklist verifications and AC-2 through AC-6: the shared `t.Parallel()` gate now includes `pi` (codex sequential), `seedPiLiveAuth` mirrors `~/.pi/agent/models.json` (0o600), `piLiveRunTimeout(dflt)` overrides the per-run cap via `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` for both runners, and `docs/runtime-live-ci.md` documents the Pi parallel + custom-model + timeout surface. The offline `gofmt`/`vet`/`build` and the `PiLiveEnv|PiIntercom|TestPiLive` unit tests (including both env-scrub assertions) pass, and no journey body, XFAIL binding, `assertConflictOwnerHandoff`, or CI lane was touched. AC-1 (live custom-model `-parallel 2` run) is not exercised here per the test plan (offline-first; live Pi run only when Pi work is authorized).
