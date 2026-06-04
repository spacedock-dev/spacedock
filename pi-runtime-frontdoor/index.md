---
title: Pi runtime front door and install UX
status: validation
score: "0.44"
source: captain (2026-06-03) — after AC-2 proved Pi subagent dispatch, expose the working runtime contract as user-facing Spacedock commands
id: mfa9gb6kpcc14b5ca0gphv2m
started: 2026-06-03T00:00:00Z
worktree: .worktrees/spacedock-ensign-pi-install-no-plugin-dir
---
# Pi runtime front door and install UX

Expose the proven Pi runtime path as Spacedock UX. The previous `pi-runtime-support` task proved that a Pi parent can dispatch a worker through `pi-subagents` from an isolated Pi home with copied OAuth credentials. This task turns that mechanism into command behavior: `spacedock pi`, `spacedock install --host pi`, and `spacedock doctor --host pi`.

## Problem

Pi runtime support currently works only through a live test harness that shells `pi` directly with explicit extension and skill paths. Users cannot discover or launch it from the Spacedock command surface, and `spacedock install --host pi` currently fails with `unknown host "pi"`.

## Scope

Implement the smallest compatibility-first UX over the proven mechanism:

- Add a `spacedock pi` launch command that starts Pi as the Spacedock first officer using Pi-native resources.
- Add `spacedock install --host pi` behavior that verifies or wires the required Pi substrate without pretending Pi is a Claude/Codex plugin marketplace.
- Add `spacedock doctor --host pi` checks for the Pi CLI, auth file, `pi-subagents` package/extension, and local Spacedock skill paths.
- Preserve current Claude/Codex behavior and output.
- Keep PR/mod behavior out of scope.

## Proven mechanism to reuse

Use `docs/runtime-support.md` as the implementation guide. The known-good Pi bringup shape is:

- host CLI: `pi`
- substrate package: `~/.pi/agent/npm/node_modules/pi-subagents`
- extension: `pi-subagents/src/extension/index.ts`
- skill: `pi-subagents/skills/pi-subagents`
- Spacedock skills: `<checkout>/skills/first-officer` and `<checkout>/skills/ensign`
- launch mode: Pi parent uses `subagent(...)`; no Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete`

## Acceptance criteria

**AC-1 - `spacedock pi` is a registered launch command that uses Pi-native resources.**
Verified by: focused CLI tests that assert help lists `pi`, the launch argv begins with `pi`, includes the local or installed `pi-subagents` extension/skill paths plus Spacedock first-officer/ensign skill paths, and excludes Claude/Codex-only flags/tool names.

**AC-2 - `spacedock install --host pi` is accepted and idempotent.**
Verified by: focused CLI tests for `install --host pi` that do not call Claude/Codex plugin commands and either confirm the existing `pi-subagents` package path or print exact next-step instructions when it is missing.

**AC-3 - `spacedock doctor --host pi` reports actionable Pi health.**
Verified by: fixture/injected-host tests covering missing `pi`, missing `auth.json`, missing `pi-subagents`, and healthy local resources with stable output and exit codes.

**AC-4 - The live Pi smoke can launch through `spacedock pi` or an equivalent Spacedock-owned wrapper.**
Verified by: updating or adding a live-gated test, ideally `go test -tags live -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v -count=1`, that keeps isolated `PI_CODING_AGENT_DIR`/session dirs and asserts durable split-root entity/git-log outcomes.

**AC-5 - Existing Claude and Codex front-door/install/doctor behavior remains stable.**
Verified by: `go test ./... -count=1` and focused existing CLI tests for Claude/Codex front doors/install/doctor.

## Test plan

- Add TDD-first CLI tests for command registration/help and Pi launch argv shape.
- Add TDD-first install/doctor tests with injected host/filesystem seams before implementation.
- Reuse the previous live smoke fixture where possible, but route through the new Spacedock-owned Pi command or wrapper.
- Run `gofmt -w ./cmd ./internal`, `go test ./... -count=1`, `go test ./... -race -count=1`, and the live Pi smoke when credentials are available.

## Stage Report: implementation

- DONE: Add `spacedock pi` launch routing with Pi-native resources.
  Evidence: code commit `255bb0da` registers `pi` in the CLI command tree/help and adds `internal/cli/pi.go`. `TestPiCommandRegisteredInTopLevelHelp` and `TestPiFrontDoorLaunchesWithNativeResourcePaths` assert the argv begins with `pi`, loads the `pi-subagents` extension/skill plus local Spacedock first-officer/ensign skills, forwards Pi flags after `--`, and excludes Claude/Codex-only runtime tokens.
- DONE: Add accepted, idempotent `spacedock install --host pi` behavior.
  Evidence: `TestPiInstallAcceptedAndDoesNotUsePluginCommands` asserts the Pi install path does not call the Claude/Codex plugin install seam and reports ready when Pi resources are present. `TestPiInstallMissingSubagentsPrintsActionableInstructions` asserts the missing-substrate path exits as an instructive setup report with `pi install npm:pi-subagents` and `PI_SUBAGENTS_PACKAGE_ROOT` guidance.
- DONE: Add stable `spacedock doctor --host pi` health checks.
  Evidence: `TestPiDoctorReportsMissingAndHealthyRuntime` covers missing and healthy Pi runtime states for the Pi CLI, auth file, `pi-subagents` extension/skill, and local Spacedock skills.
- DONE: Route live Pi proof through the Spacedock-owned front door.
  Evidence: `TestLivePiFrontDoorSmoke` shells the current-checkout `spacedock pi` command with isolated `PI_CODING_AGENT_DIR`, temp session dir, copied OAuth auth, explicit `PI_SUBAGENTS_PACKAGE_ROOT`, and a split-root workflow fixture. It passed alongside the existing raw Pi smoke.
- DONE: Preserve Claude/Codex behavior.
  Evidence: fresh `go test ./... -count=1` and `go test ./... -race -count=1` passed, including existing CLI/front-door/install/doctor tests.

### Summary

Implemented the Pi runtime front-door UX in code commit `255bb0da`. Spacedock now exposes `spacedock pi`, accepts `install --host pi`, and reports Pi runtime health via `doctor --host pi` without reusing Claude/Codex plugin commands or team-tool semantics. Verification passed with fresh baseline, race, and live Pi smoke commands, including `TestLivePiFrontDoorSmoke` through the new Spacedock-owned wrapper.

## Stage Report: validation

- AC-1: PASSED. Commit `255bb0da` registers `spacedock pi` in the cobra command tree and top-level help. `internal/cli/pi.go` builds argv beginning with `pi`, `--extension <pi-subagents>/src/extension/index.ts`, `--skill <pi-subagents>/skills/pi-subagents`, and local Spacedock first-officer/ensign skill paths; `TestPiFrontDoorLaunchesWithNativeResourcePaths` asserts Claude/Codex-only tokens are absent.
- AC-2: PASSED. `install --host pi` is accepted via `runInitWithPi`; focused tests show it does not call the Claude/Codex plugin install seam, reports ready when resources exist, and prints `pi install npm:pi-subagents` plus `PI_SUBAGENTS_PACKAGE_ROOT` guidance when the substrate is missing.
- AC-3: PASSED. `doctor --host pi` reports Pi CLI, Pi auth, pi-subagents extension/skill, and Spacedock skill health with actionable remedies; focused tests cover missing and healthy states with expected exit codes.
- AC-4: NOT PASSED for the required live validation command. `TestLivePiFrontDoorSmoke` itself passed through `spacedock pi`, but the requested combined live command failed because `TestLivePiSubagentEnsignSmoke` produced `## implementation report` instead of the asserted `## Stage Report: implementation`. No credential/package path was missing; Pi auth and `PI_SUBAGENTS_PACKAGE_ROOT` were available.
- AC-5: PASSED. Baseline `go test ./... -count=1` and `go test ./... -race -count=1` passed from `.worktrees/spacedock-ensign-pi-runtime-support`, covering existing Claude/Codex front-door/install/doctor behavior.

Validation commands:
- PASS: `gofmt -w ./cmd ./internal` (no tracked product diffs after formatting).
- PASS: `go test ./... -count=1`.
- PASS: `go test ./... -race -count=1`.
- FAIL: `go test -tags live -run 'TestLivePi(SubagentEnsignSmoke|FrontDoorSmoke)' ./internal/ensigncycle -v -count=1` (`TestLivePiFrontDoorSmoke` passed; `TestLivePiSubagentEnsignSmoke` failed on the expected stage-report heading).

Recommendation: REJECTED

## Stage Report: implementation (validation fix)

- DONE: Read the rejected validation report for commit `d7348a12` and addressed only the AC-4 live heading failure. The failure was the raw Pi live smoke prompt asking for a generic implementation stage report while `assertPiLiveSmokeResult` requires `## Stage Report: implementation`.
- DONE: Added focused TDD coverage in isolated worktree `.worktrees/spacedock-ensign-pi-runtime-frontdoor-fixback`: `TestPiLiveSmokePromptRequiresExactStageReportHeading` first failed against the old prompt, then passed after requiring the exact heading.
- DONE: Updated `piLiveSmokePrompt` to instruct the Pi ensign to append an implementation report with the exact heading `## Stage Report: implementation`, matching the existing live assertion.
- DONE: Committed product fix `3bcaf487b8b57e810439ac9a39e7f345c8a08c34` (`test: require exact Pi live smoke report heading`) on branch `spacedock-ensign/pi-runtime-frontdoor-fixback`.

### Commands run

- FAIL before fix: `go test -tags live -run TestPiLiveSmokePromptRequiresExactStageReportHeading ./internal/ensigncycle -count=1` failed because the prompt did not contain `exact heading '## Stage Report: implementation'`.
- PASS after fix: `gofmt -w ./cmd ./internal`.
- PASS after fix: `go test -tags live -run TestPiLiveSmokePromptRequiresExactStageReportHeading ./internal/ensigncycle -count=1`.
- PASS after fix: `go test ./... -count=1`.
- PASS after fix: `go test ./... -race -count=1`.
- PASS after fix: `go test -tags live -run 'TestLivePi(SubagentEnsignSmoke|FrontDoorSmoke)' ./internal/ensigncycle -v -count=1`; both `TestLivePiSubagentEnsignSmoke` and `TestLivePiFrontDoorSmoke` passed with available Pi auth and `pi-subagents` package path.

### Summary

The validation rejection has been routed back through implementation with a narrow, test-backed fix. The live prompt now explicitly requires the exact stage-report heading already asserted by the smoke test. This does not mark final independent validation passed; it is ready for an independent validation rerun.

### Residual risks

- Live Pi smoke still depends on local Pi credentials and `pi-subagents` availability in the validation environment.

## Stage Report: validation (fix rerun)

- AC-1: PASSED. Re-inspected the isolated fixback worktree on branch `spacedock-ensign/pi-runtime-frontdoor-fixback` at product commit `3bcaf487b8b57e810439ac9a39e7f345c8a08c34`. `spacedock pi` remains registered in help/routing, and `internal/cli/pi.go` launches argv beginning with `pi` plus Pi-native `--extension <pi-subagents>/src/extension/index.ts`, `--skill <pi-subagents>/skills/pi-subagents`, and local Spacedock first-officer/ensign skill paths. Focused tests still assert Claude/Codex-only tokens are absent.
- AC-2: PASSED. `install --host pi` remains accepted and idempotent via the Pi runtime path. Focused tests verify the ready case, the missing-substrate instruction path (`pi install npm:pi-subagents` and `PI_SUBAGENTS_PACKAGE_ROOT`), and that the Claude/Codex plugin install seam is not called.
- AC-3: PASSED. `doctor --host pi` still reports actionable Pi runtime health for the Pi CLI, auth file, pi-subagents extension/skill, and Spacedock skill paths, with missing and healthy fixture coverage and stable exit-code behavior.
- AC-4: PASSED. Inspected product fix commit `3bcaf487`; it adds `TestPiLiveSmokePromptRequiresExactStageReportHeading` and changes the raw Pi live smoke prompt to require the exact heading `## Stage Report: implementation`, matching the existing smoke assertion. Fresh combined live command passed for both `TestLivePiSubagentEnsignSmoke` and `TestLivePiFrontDoorSmoke` with available Pi auth and pi-subagents package path.
- AC-5: PASSED. Fresh baseline and race suites passed from the isolated fixback worktree, covering existing Claude/Codex front-door/install/doctor behavior.

Validation commands:
- PASS: `gofmt -w ./cmd ./internal` (no tracked product diffs after formatting).
- PASS: `go test ./... -count=1`.
- PASS: `go test ./... -race -count=1`.
- PASS: `go test -tags live -run 'TestLivePi(SubagentEnsignSmoke|FrontDoorSmoke)' ./internal/ensigncycle -v -count=1` (`TestLivePiSubagentEnsignSmoke` passed in 217.73s; `TestLivePiFrontDoorSmoke` passed in 78.10s).

Recommendation: PASSED

### Residual risks

- The live Pi smoke remains dependent on local Pi credentials and a local `pi-subagents` package being available in future validation environments.

## Stage Report: implementation (pi-subagents acceptance contract fix)

- DONE: Updated the Pi first-officer runtime contract to forbid using the `subagent(... acceptance: ...)` contract for Spacedock stage dispatches through `pi-subagents` while preserving Pi dispatch through `subagent(...)`.
- DONE: Documented the required rationale: Spacedock puts acceptance requirements in the task prompt/dispatch content and gates work through entity stage reports, product/state commits, and independent validation rather than same-agent acceptance finalization.
- DONE: Added a focused instruction-text invariant test that scopes the rule to the Pi first-officer Dispatch section and checks the banned shape plus the Spacedock independent-validation rationale.

Changed files:
- `skills/first-officer/references/pi-first-officer-runtime.md`
- `skills/integration/skill_surface_test.go`

Product commit:
- `7119e74d83b8d7077a8824d03399fa3928171b5e` (`fix pi subagent acceptance stage contract`)

Tests run:
- PASS: `go test ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- Not needed: `gofmt -w ./cmd ./internal` because no `cmd/` or `internal/` Go files changed. `gofmt` was run on `skills/integration/skill_surface_test.go`.

Residual risks:
- This is an instruction-surface invariant, so independent validation should still review the contract text for intended semantics.

## Stage Report: validation (pi-subagents acceptance contract fix)

- VERIFIED: Inspected isolated implementation worktree `.worktrees/spacedock-ensign-pi-runtime-no-subagent-acceptance` on branch `spacedock-ensign/pi-runtime-no-subagent-acceptance` at product commit `7119e74d83b8d7077a8824d03399fa3928171b5e` (`fix pi subagent acceptance stage contract`). The product commit changes only `skills/first-officer/references/pi-first-officer-runtime.md` and `skills/integration/skill_surface_test.go`.
- PASSED: Pi runtime still allows Pi-native dispatch through `subagent(...)`. The Dispatch section says a Pi first officer may dispatch via `subagent(...)` when available, with the emitted dispatch prompt/content, workflow directory, entity path, completion checklist, and Spacedock ensign/Pi adapter loading requirements.
- PASSED: The Pi runtime now explicitly forbids `subagent(... acceptance: ...)` for Spacedock stage dispatches through `pi-subagents`, while instructing first officers to put acceptance requirements in the task prompt/dispatch content instead.
- PASSED: The rationale is present and matches the Spacedock model: implementation-to-validation gating is via entity stage reports, product/state commits, and independent validation rather than same-agent acceptance finalization by the child that did the work.
- PASSED: The invariant test `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages` scopes its check to the `## Dispatch` section and requires the banned acceptance-contract shape plus the task-prompt/dispatch-content and independent-validation rationale phrases.

Validation commands:
- PASS: `go test ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `gofmt -l skills/integration/skill_surface_test.go` produced no output.

Recommendation: PASSED

### Residual risks

- This validation covers the instruction and invariant-test surface only. It did not run expensive live Pi smokes, consistent with the task scope.

## Stage Report: implementation (staff-review fixes)

- Captain decision: Pi auth readiness is accepted as not this fixback's problem; no Pi auth readiness behavior was changed.
- DONE: Restored non-Pi `--plugin-dir` setup compatibility by routing `install`/`doctor` for Claude/Codex through the existing non-Pi parsers without stripping `--plugin-dir`, so those commands reject the unsupported flag instead of silently ignoring it.
- DONE: Kept Pi `--plugin-dir` setup support for `install --host pi` and `doctor --host pi` intact.
- DONE: Updated stale `docs/runtime-support.md` wording to describe current `spacedock pi`, `spacedock install --host pi`, and `spacedock doctor --host pi` behavior.

Changed files:
- `internal/cli/pi.go`
- `internal/cli/pi_frontdoor_test.go`
- `docs/runtime-support.md`

Product commit:
- `b63c1cba242e5bdba31d8d2e9258b662449a7273` (`fix pi setup plugin-dir compatibility`)

Tests run:
- PASS: `gofmt -w ./cmd ./internal`
- PASS: `go test ./internal/cli -run 'Test(PiInstallAcceptedAndDoesNotUsePluginCommands|PiInstallMissingSubagentsPrintsActionableInstructions|PiDoctorReportsMissingAndHealthyRuntime|NonPiSetupRejectsPluginDir)' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

Residual risks:
- Staff-review fixes were fixture/unit validated only; live Pi smoke was not rerun for this narrow compatibility/doc fixback.
- Pi auth readiness behavior intentionally remains unchanged per captain decision.

## Stage Report: validation (staff-review fixes)

- VERIFIED: Read the latest `## Stage Report: implementation (staff-review fixes)` and inspected isolated worktree `.worktrees/spacedock-ensign-pi-frontdoor-polish-fix` on branch `spacedock-ensign/pi-frontdoor-polish-fix` at product commit `b63c1cba242e5bdba31d8d2e9258b662449a7273` (`fix pi setup plugin-dir compatibility`). Changed product files are limited to `internal/cli/pi.go`, `internal/cli/pi_frontdoor_test.go`, and `docs/runtime-support.md`.
- PASSED: Non-Pi `--plugin-dir` compatibility fix is present. `runInitWithPi` and `runDoctorWithPi` now forward the original args to Claude/Codex setup paths instead of stripping `--plugin-dir`, so the existing non-Pi parsers reject the unsupported flag. `TestNonPiSetupRejectsPluginDir` covers install/doctor for Claude and Codex, expects exit code 2 and `unknown argument "--plugin-dir"`, and verifies the plugin install seam is not called.
- PASSED: Pi `--plugin-dir` compatibility remains intact. The Pi setup parser still accepts `--plugin-dir` for `--host pi`, and focused install/doctor tests invoke Pi setup with `--plugin-dir` and pass through the Pi runtime config path.
- PASSED: `docs/runtime-support.md` no longer says Pi install is future work. It now describes `spacedock pi` as the front door, `spacedock install --host pi` as an idempotent readiness check/setup guide, and `spacedock doctor --host pi` as reporting Pi CLI, auth, `pi-subagents`, and Spacedock skill health. It also documents that Pi setup commands accept `--plugin-dir <spacedock checkout>`.
- VERIFIED: Captain decision about auth is documented in the implementation report. The product diff does not change auth readiness behavior; unchanged Pi doctor/auth wording remains limited to reporting/checking the existing auth file.

Validation commands:
- PASS: `gofmt -w ./cmd ./internal` (no tracked product diffs after formatting).
- PASS: `go test ./internal/cli -run 'Test(PiInstallAcceptedAndDoesNotUsePluginCommands|PiInstallMissingSubagentsPrintsActionableInstructions|PiDoctorReportsMissingAndHealthyRuntime|NonPiSetupRejectsPluginDir)' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

Recommendation: PASSED

### Residual risks

- Validation was limited to CLI compatibility/docs as requested; no live Pi smoke was run.
- Pi auth readiness remains unchanged by captain decision and was not evaluated as a blocker for this fixback.

## Stage Report: implementation (install plugin-dir correction)

### Captain decision

- Do not support `spacedock install --host pi --plugin-dir ...` for this slice.
- Preserve `spacedock pi --plugin-dir ...` launch support.
- Preserve `spacedock doctor --host pi --plugin-dir ...` diagnostics support.
- Do not change Pi auth readiness behavior.

### Changed files

- `internal/cli/pi.go`: rejects `--plugin-dir` during `install` parsing with a clear usage error, keeps doctor parsing unchanged, and removes plugin-dir from install next-step prose.
- `internal/cli/pi_frontdoor_test.go`: adds focused rejection coverage for `install --host pi --plugin-dir /checkout`; updates install readiness/instruction tests to use `SPACEDOCK_REPO_ROOT`; keeps doctor and Pi launch plugin-dir coverage.
- `internal/cli/cli.go`: removes the install help example that showed `--plugin-dir` for Pi install.
- `docs/runtime-support.md`: documents that Pi install does not accept `--plugin-dir`, while Pi doctor still accepts it for local checkout diagnostics.

### Product commit

- `8f2c55a469d491a724bce4a9f44efa665716ce89` (`pi install rejects plugin-dir`)

### Tests

- PASS: `gofmt -w ./cmd ./internal`
- PASS: `go test ./internal/cli -run 'TestPi(InstallRejectsPluginDir|InstallAcceptedAndDoesNotUsePluginCommands|InstallMissingSubagentsPrintsActionableInstructions|DoctorReportsMissingAndHealthyRuntime|FrontDoorLaunchesWithNativeResourcePaths)|TestNonPiSetupRejectsPluginDir' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

### Residual risks

- `spacedock install --host pi` still resolves the checkout through the existing cwd/`SPACEDOCK_REPO_ROOT` path logic; local dev checkout path improvements remain intentionally deferred.
- Pi auth readiness behavior was not changed; launch readiness still excludes auth while doctor health includes it.

## Stage Report: validation (install plugin-dir correction)

- VERIFIED: Read the latest `## Stage Report: implementation (install plugin-dir correction)` and inspected isolated worktree `.worktrees/spacedock-ensign-pi-install-no-plugin-dir` on branch `spacedock-ensign/pi-install-no-plugin-dir` at product commit `8f2c55a469d491a724bce4a9f44efa665716ce89` (`pi install rejects plugin-dir`). Changed product files are `internal/cli/pi.go`, `internal/cli/pi_frontdoor_test.go`, `internal/cli/cli.go`, and `docs/runtime-support.md`.
- PASSED: `spacedock install --host pi --plugin-dir /checkout` now rejects before runtime checks with exit code 2 and a clear message: `spacedock install: --plugin-dir is not supported; use SPACEDOCK_REPO_ROOT or run from the Spacedock checkout`. Focused coverage is `TestPiInstallRejectsPluginDir`.
- PASSED: `spacedock install --host pi` still works as the Pi readiness/check/instructions path. The ready fixture reports `Pi runtime ready` without calling Claude/Codex plugin install seams, and the missing-substrate fixture prints actionable `pi install npm:pi-subagents` and `PI_SUBAGENTS_PACKAGE_ROOT` guidance. The install next-step prose no longer suggests `--plugin-dir`.
- PASSED: `spacedock pi --plugin-dir /checkout` remains supported. `parsePiFrontDoorArgs` still accepts repeatable `--plugin-dir`, `runPi` resolves the local skill checkout from it, and `TestPiFrontDoorLaunchesWithNativeResourcePaths` verifies Pi-native extension/skill argv with local Spacedock first-officer/ensign paths.
- PASSED: `spacedock doctor --host pi --plugin-dir /checkout` remains intentionally supported for local checkout diagnostics. `parsePiSetupArgs` only rejects `--plugin-dir` for `install`, and `TestPiDoctorReportsMissingAndHealthyRuntime` still exercises doctor with `--plugin-dir` in missing and healthy cases.
- PASSED: Docs/help no longer claim Pi install accepts `--plugin-dir`. `internal/cli/cli.go` install examples now show `spacedock install --host pi` without `--plugin-dir`, while doctor help still shows `spacedock doctor --host pi --plugin-dir ./checkout`; `docs/runtime-support.md` explicitly says Pi install does not accept `--plugin-dir` and documents `SPACEDOCK_REPO_ROOT`/working-directory resolution instead.
- VERIFIED: Pi auth readiness behavior was not changed. `piRuntimeLaunchReady` still excludes auth, `piDoctorHealthy` still requires auth, and the auth check/report text remains limited to doctor health.

Validation commands:
- PASS: `gofmt -w ./cmd ./internal` (no tracked product diffs after formatting).
- PASS: `go test ./internal/cli -run 'TestPi(InstallRejectsPluginDir|InstallAcceptedAndDoesNotUsePluginCommands|InstallMissingSubagentsPrintsActionableInstructions|DoctorReportsMissingAndHealthyRuntime|FrontDoorLaunchesWithNativeResourcePaths)|TestNonPiSetupRejectsPluginDir' -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

Recommendation: PASSED

### Residual risks

- Validation was limited to CLI/docs as requested; no live Pi smoke was run.
- `spacedock install --host pi` still relies on the existing working-directory or `SPACEDOCK_REPO_ROOT` checkout resolution path by design.
