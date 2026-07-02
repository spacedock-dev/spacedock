# Pi install --plugin-dir correction validation

Recommendation: PASSED

State commit: `ad3c677de2dcf967e18cf76544d2e55507bd7758` (`validation: pi install rejects plugin-dir`)

Validated isolated worktree: `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-install-no-plugin-dir`
Branch: `spacedock-ensign/pi-install-no-plugin-dir`
Product commit: `8f2c55a469d491a724bce4a9f44efa665716ce89`

Findings:
- `install --host pi --plugin-dir /checkout` rejects clearly with exit code 2 and `--plugin-dir is not supported` guidance.
- `install --host pi` still works as the readiness/check/instructions path and does not use Claude/Codex plugin install seams.
- `spacedock pi --plugin-dir /checkout` remains supported for launch skill checkout selection.
- `doctor --host pi --plugin-dir /checkout` remains supported for local skill checkout diagnostics.
- Install help/docs no longer claim Pi install accepts `--plugin-dir`; doctor docs/help may still show it.
- Pi auth readiness behavior is unchanged: launch readiness excludes auth, doctor health includes auth.

Validation commands:
- PASS: `gofmt -w ./cmd ./internal`
- PASS: `go test ./internal/cli -run 'TestPi(InstallRejectsPluginDir|InstallAcceptedAndDoesNotUsePluginCommands|InstallMissingSubagentsPrintsActionableInstructions|DoctorReportsMissingAndHealthyRuntime|FrontDoorLaunchesWithNativeResourcePaths)|TestNonPiSetupRejectsPluginDir' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

Residual risks:
- Validation was limited to CLI/docs as requested; no live Pi smoke was run.
- `spacedock install --host pi` still relies on existing working-directory or `SPACEDOCK_REPO_ROOT` checkout resolution by design.
