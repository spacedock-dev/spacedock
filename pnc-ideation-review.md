# IDEATION GATE — Pi live harness: parallelism + custom-provider model + slow-model timeout (`pnc`)

Recommendation: **APPROVE and dispatch implementation.**

## Selected approach

Three independent Pi-runner-harness changes, bundled because they enable the
same operator scenario (run the common journeys against a custom slow-thinking
model with real parallelism):

1. **Parallelism** — `shared_live_runner_test.go:44`: widen the `t.Parallel()`
   condition from `SPACEDOCK_LIVE_RUNTIME == "claude"` to `== "claude" || "pi"`.
   Codex stays sequential (no `t.Parallel()` for codex). Each Pi journey already
   isolates per-journey (`t.TempDir()` Pi-home, session-dir, clean-home, fixture
   root), so concurrency is safe.
2. **Custom-provider model mirror** — `seedPiLiveAuth` in
   `pi_live_runner_test.go`: copy `~/.pi/agent/models.json` into the isolated Pi
   home (0o600, alongside the existing `auth.json` copy) when present. Custom
   providers (e.g. `lunaroute`) declare models in `models.json`, not `auth.json`;
   an auth-only isolated home currently gets `Model ... not found`.
3. **Slow-model timeout** — new `piLiveRunTimeout(dflt)` helper reads
   `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` and returns `dflt` when unset/invalid;
   the common `run` (`pi_shared_live_runner_test.go`) and the smoke
   `runPiLiveCommand` (`pi_live_runner_test.go`) use it instead of the bare
   `12*time.Minute`/`10*time.Minute` constants. Fatal message names the env var
   + default.

## Risk evidence

- **No spike needed.** All three mechanisms are proven: `t.Parallel()` already
  fans out Claude journeys (registry proves `-parallel 2` works for Claude);
  `seedPiLiveAuth` already copies `auth.json` (the `models.json` copy is
  additive); `os.Getenv` parsing is standard. The `lunaroute/...:max` model
  resolved and ran the front-door smoke + `FullEnsignCycle` + `GateGuardrail` +
  `RecordedGateLifecycle` + the 4 XFAIL journeys locally this session once the
  `models.json` mirror and `piLiveRunTimeout` were in place (the edits are in
  `stash@{0}`, verified working before filing).
- **Env-scrub risk** — the additive `models.json` copy must not break
  `TestPiLiveEnvDropsForeignRuntimeMarkers` / `TestPiLiveEnvScrubsAmbientPiSubagentMarkers`.
  The copy is in `seedPiLiveAuth` (file copy), not `piLiveEnv` (env scrubbing),
  so the scrub tests are unaffected. Verified offline in the stashed edits.
- **Parallel-contention risk** — running 2 `:max`-thinking sessions in parallel
  roughly triples per-journey latency; observed this session (solo `keep-moving`
  ~335s; parallel `keep-moving` ~1079s). That's why AC-1 pairs the per-run cap
  raise with the parallelism; without the cap raise, parallel journeys time out.
  The outer `go test -timeout` must exceed the per-run cap × selected journeys.

## Expected surface and tolerance

- `internal/ensigncycle/shared_live_runner_test.go` — ~1 line (the `t.Parallel()`
  condition).
- `internal/ensigncycle/pi_live_runner_test.go` — `seedPiLiveAuth` +~7 lines
  (models.json copy) + `piLiveRunTimeout` helper +~10 lines + `runPiLiveCommand`
  call-site swap +~1 line.
- `internal/ensigncycle/pi_shared_live_runner_test.go` — `run` call-site swap
  +~1 line, message +~1 line.
- `docs/runtime-live-ci.md` — doc updates (~+20 lines).

**Total estimate: ~+40 lines net, tolerance ±10.** No semantic changes to
command grammar, stored formats, authority, or runtime behavior — test-harness
only (the `t.Parallel()` change is test concurrency, not product behavior).

## Semantic changes

None. Test-harness only. No CLI output, no stored format, no authority, no
runtime behavior change. The doc change describes existing-but-undocumented
behavior (Pi parallelism, the env var, the models.json copy).

## Proposed proof per acceptance criterion

- **AC-1 (VALUE — custom slow `:max` model runs common suite two-at-a-time to
  graded results):** a focused live run
  `SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_CHILD_MODEL='lunaroute/glm-5.2-vision-background:max' SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES=40 go test -run '^TestLiveCommon' -parallel 2`
  shows two `=== CONT` lines together (fan-out) and journeys complete to
  XFAIL/XPASS/PASS, not `Model ... not found` or a 12m timeout. Baseline: today
  the suite skips (`Model ... not found`) without the models.json mirror, or
  times out at 12m without the cap raise.
- **AC-2 (Pi `t.Parallel()`; Codex sequential):** `go test -v -run '^TestLiveCommon' -parallel 2`
  log shows `=== PAUSE`/`=== CONT` for Pi (parallel) and only sequential RUN
  lines for Codex. Offline: grep the widened condition.
- **AC-3 (models.json mirror):** `seedPiLiveAuth` copies `models.json` (0o600)
  when present; a focused run with a custom-provider model resolves (no
  `Model ... not found`). The `TestPiLiveEnv*` scrub tests still pass (the copy
  is not in `piLiveEnv`).
- **AC-4 (overridable cap):** `piLiveRunTimeout(12*time.Minute)` returns the env
  value when set, `dflt` when unset/invalid; the fatal message names the env var.
  Offline unit test for the helper.
- **AC-5 (doc describes surface):** `docs/runtime-live-ci.md` diff shows the Pi
  parallel + slow-model + models.json sections.
- **AC-6 (offline + env-scrub tests):** `gofmt`, `go vet -tags live`,
  `go build -tags live`, `go test -run 'PiLiveEnv|PiIntercom|TestPiLive' ./internal/ensigncycle`
  pass; `TestPiLiveEnvDropsForeignRuntimeMarkers` and
  `TestPiLiveEnvScrubsAmbientPiSubagentMarkers` hold.

## Decision ask

Approve to dispatch implementation (worktree) for this three-part Pi-harness
bundle, or revise/hold with a concrete boundary.
