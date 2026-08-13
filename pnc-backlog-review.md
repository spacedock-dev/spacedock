# BACKLOG GATE — Pi live harness: parallelism + custom-provider model + slow-model timeout (`pnc`)

Recommendation: **APPROVE and dispatch ideation.**

## Capability and value

Three Pi runner harness gaps block running the common journeys against a custom
slow-thinking model (`lunaroute/glm-5.2-vision-background:max`):
1. `t.Parallel()` is Claude-only in `shared_live_runner_test.go`, so `-parallel 2`
   is a no-op for Pi (17 journeys run sequentially).
2. `seedPiLiveAuth` copies only `auth.json`; custom providers declare models in
   `models.json`, so an isolated home can't resolve the model (`Model ... not found`).
3. The per-run cap is a fixed `12m`/`10m`; a slow `:max` model exceeds it on
   multi-dispatch journeys, producing a timeout `--- FAIL` instead of a graded result.

This task makes all three overridable/parallel so an operator runs the suite
against a custom slow model with real parallelism to graded results.

## Binding boundaries

- Harness enablement only. No journey's exercise, fixture, or durable assertion
  changes. The shared `t.Parallel()` gate widens to `pi`; Codex stays sequential.
- Does not touch `assertConflictOwnerHandoff` or any XFAIL binding (owned by
  `rzr` and the runtime-specific repair entities).
- Does not change the `0a` CI workflow cadence (archived, done) or merge policy.
- No new registry, result artifact, host-switch table, or lifecycle model.

## Proof direction

Ideation confirms the three seams: (1) widen the `t.Parallel()` condition to
`"claude" || "pi"`; (2) `seedPiLiveAuth` mirrors `models.json` (0o600, alongside
`auth.json`); (3) a `piLiveRunTimeout(dflt)` helper reads
`SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` and the `run`/`runPiLiveCommand` use it.
Implementation proves offline `PiLiveEnv|PiIntercom|TestPiLive` unit tests pass
(including the env-scrub markers, which the additive `models.json` copy must not
break) and `docs/runtime-live-ci.md` describes the Pi parallel + slow-model
surface. A focused live run with `-parallel 2` confirms fan-out, authorized separately.

## Decision ask

Approve this Pi-harness enablement bundle for ideation, or revise/hold with a
concrete boundary.
