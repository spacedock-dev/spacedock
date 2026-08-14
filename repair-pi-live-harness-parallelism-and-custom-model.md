---
title: Repair Pi live harness — common-journey parallelism, custom-provider model mirror, and slow-model timeout
status: implementation
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
