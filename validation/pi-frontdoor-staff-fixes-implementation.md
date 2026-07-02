# Pi frontdoor staff fixes implementation

Implemented staff-review fixbacks for `docs/dev/.spacedock-state/pi-runtime-frontdoor/index.md`.

## Captain decision

Pi auth readiness is accepted as not this fixback's problem. No Pi auth readiness behavior was changed.

## Changes

- Restored non-Pi `--plugin-dir` compatibility for setup commands:
  - `spacedock install --host claude/codex --plugin-dir ...` now reaches the existing non-Pi parser and rejects `--plugin-dir` instead of silently stripping it.
  - `spacedock doctor --host claude/codex --plugin-dir ...` does the same.
- Kept Pi `--plugin-dir` support for `install --host pi` and `doctor --host pi`.
- Updated `docs/runtime-support.md` stale wording to describe current `spacedock pi`, `spacedock install --host pi`, and `spacedock doctor --host pi` behavior.

## Changed files

- `internal/cli/pi.go`
- `internal/cli/pi_frontdoor_test.go`
- `docs/runtime-support.md`

## Commits

- Product commit: `b63c1cba242e5bdba31d8d2e9258b662449a7273` (`fix pi setup plugin-dir compatibility`)
- State commit: `349194e89a59d12b121b512942b1088d362cb688` (`implementation: pi runtime frontdoor staff fixes`)

## Validation

- PASS: `gofmt -w ./cmd ./internal`
- PASS: `go test ./internal/cli -run 'Test(PiInstallAcceptedAndDoesNotUsePluginCommands|PiInstallMissingSubagentsPrintsActionableInstructions|PiDoctorReportsMissingAndHealthyRuntime|NonPiSetupRejectsPluginDir)' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

## Residual risks

- Live Pi smoke was not rerun for this narrow compatibility/doc fixback.
- Pi auth readiness remains unchanged by design.
