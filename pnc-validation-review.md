# VALIDATION GATE — Pi live harness: parallelism + custom-provider model + slow-model timeout (`pnc`)

Recommendation: **APPROVE (PASSED).**

## What was verified

A fresh validator independently verified the worktree deliverable
(`.worktrees/spacedock-ensign-repair-pi-live-harness-parallelism-and-custom-model`,
commit `864a3ca97`, +43/-12 across 4 files) against the entity spec, without
editing it:

- **AC-2 (Pi `t.Parallel()`; Codex sequential):** `shared_live_runner_test.go:44-47`
  switches on `claude`/`pi` → `t.Parallel()`; only line 46 inside the switch; no
  `codex.*parallel` match (codex falls through).
- **AC-3 (models.json mirror):** `pi_live_runner_test.go:240-247` — after the
  `auth.json` write (0o600), reads `~/.pi/agent/models.json`; on nil err + non-empty,
  writes `piHome/models.json` (0o600).
- **AC-4 (overridable cap):** `piLiveRunTimeout(dflt)` at `:285-292` parses
  `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` (strconv.Atoi, >0), returns `dflt` otherwise;
  `run` uses `piLiveRunTimeout(12*time.Minute)` (`pi_shared_live_runner_test.go:62`);
  `runPiLiveCommand` uses `piLiveRunTimeout(10*time.Minute)` (`pi_live_runner_test.go:68`).
- **AC-5 (doc describes surface):** `docs/runtime-live-ci.md` covers Pi
  `-parallel 2`, `provider/model:thinking`, `models.json` mirror + `Model ... not
  found` check, `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES` guidance, common Pi command uses
  `-timeout 90m -parallel 2`.
- **AC-6 (offline + env-scrub):** `gofmt -l` empty; `go vet -tags live` exit 0;
  `go build -tags live` exit 0; `go test -run 'PiLiveEnv|PiIntercom|TestPiLive'` →
  `ok 0.619s`; `TestPiLiveEnvDropsForeignRuntimeMarkers` and
  `TestPiLiveEnvScrubsAmbientPiSubagentMarkers` both PASS.
- **Scope:** `git diff --stat` names only the 4 expected files; no journey body,
  XFAIL binding, `assertConflictOwnerHandoff`, or CI lane touched. Worktree clean
  (read-only validation).

## Reviewer findings

None. The deliverable is the three-part harness enablement matching the spec;
no material issue found.

## Risks (deferred)

- AC-1 (VALUE — a live custom-model `-parallel 2` run to graded results) is not
  exercised offline — the test plan is offline-first and a live Pi run requires
  authorized Pi work. The validator named the exact confirmation run for when Pi
  work is authorized:
  `SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_CHILD_MODEL='lunaroute/glm-5.2-vision-background:max' SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES=40 go test -tags live -run '^TestLiveCommon' -parallel 2 -timeout 120m ./internal/ensigncycle -v`
  (expect two `=== CONT` lines fan out + graded results, not `Model ... not found`
  or a 12m timeout). This is a residual evidence step, not a regression — the
  harness enablement is offline-proven; the live run confirms end-to-end.

## Decision

Approve to terminalize (merge the worktree branch
`spacedock-ensign/repair-pi-live-harness-parallelism-and-custom-model`). The
deliverable is the smallest sufficient mechanism, harness-only, with all
journeys/bindings/CI lanes preserved.
