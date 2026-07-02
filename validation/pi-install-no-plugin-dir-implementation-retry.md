# Pi install plugin-dir correction — implementation retry

Implemented removal of `--plugin-dir` support from `spacedock install --host pi` in the isolated worktree:

`/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-install-no-plugin-dir`

## Captain decision applied

- `spacedock install --host pi --plugin-dir ...` now rejects with a clear error.
- `spacedock install --host pi` remains an idempotent readiness/check/instructions path using existing cwd / `SPACEDOCK_REPO_ROOT` / defaults resolution.
- `spacedock pi --plugin-dir ...` remains supported.
- `spacedock doctor --host pi --plugin-dir ...` remains supported for local skill checkout diagnostics.
- Pi auth readiness behavior was not changed.

## Changed files

Product/docs commit: `8f2c55a469d491a724bce4a9f44efa665716ce89` (`pi install rejects plugin-dir`)

- `internal/cli/pi.go`
- `internal/cli/pi_frontdoor_test.go`
- `internal/cli/cli.go`
- `docs/runtime-support.md`

State commit: `05feb180650cd2b68eedca4284d7f7437dcac0ca` (`implementation: pi install rejects plugin-dir`)

- `docs/dev/.spacedock-state/pi-runtime-frontdoor/index.md`

## Validation

- PASS: `gofmt -w ./cmd ./internal`
- PASS: `go test ./internal/cli -run 'TestPi(InstallRejectsPluginDir|InstallAcceptedAndDoesNotUsePluginCommands|InstallMissingSubagentsPrintsActionableInstructions|DoctorReportsMissingAndHealthyRuntime|FrontDoorLaunchesWithNativeResourcePaths)|TestNonPiSetupRejectsPluginDir' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

## Residual risks

- The local dev checkout install path remains the existing cwd / `SPACEDOCK_REPO_ROOT` behavior; further path-polish was intentionally deferred by captain decision.
- Auth behavior is unchanged: launch readiness still excludes auth, while doctor health includes auth.
