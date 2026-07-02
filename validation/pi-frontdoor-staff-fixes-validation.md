# Pi runtime frontdoor staff-review fixes validation

Recommendation: PASSED

Evidence:
- Inspected isolated worktree `.worktrees/spacedock-ensign-pi-frontdoor-polish-fix` on branch `spacedock-ensign/pi-frontdoor-polish-fix` at product commit `b63c1cba242e5bdba31d8d2e9258b662449a7273`.
- Verified changed files: `internal/cli/pi.go`, `internal/cli/pi_frontdoor_test.go`, and `docs/runtime-support.md`.
- Non-Pi install/doctor now reject `--plugin-dir` instead of silently ignoring it; Pi install/doctor still accept `--plugin-dir`.
- Runtime support docs now describe current `spacedock pi`, `install --host pi`, and `doctor --host pi` behavior.
- Auth readiness behavior was not changed; captain decision is documented in the implementation report.

Commands:
- PASS: `gofmt -w ./cmd ./internal`
- PASS: `go test ./internal/cli -run 'Test(PiInstallAcceptedAndDoesNotUsePluginCommands|PiInstallMissingSubagentsPrintsActionableInstructions|PiDoctorReportsMissingAndHealthyRuntime|NonPiSetupRejectsPluginDir)' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

Residual risks:
- No live Pi smoke was run, consistent with the narrow CLI compatibility/docs scope.
- Pi auth readiness remains intentionally unchanged and was not treated as a rejection reason.
