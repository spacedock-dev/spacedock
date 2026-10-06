---
title: The Pi launcher demands an extension path that Pi itself already resolves from settings
status: implementation
score: 0.8
source: "FO review of the Pi stack, 2026-10-04: internal/cli/pi.go:645-648 resolves pi-subagents from the env var or HOME, while Pi resolves a registered npm entry from the agent directory."
id: 7qksxxpbdcqz96mh0basxxzb
worktree: .worktrees/spacedock-ensign-launcher-resolves-registered-pi-extensions
---

## Problem

The launcher does not use the same resolution as the host it launches.

`piRuntimeConfigFromEnv` in `internal/cli/pi.go` reads `PI_SUBAGENTS_PACKAGE_ROOT`,
and when that is unset it falls back to
`join(HOME, ".pi/agent/npm/node_modules/pi-subagents")`. Pi itself resolves the
`npm:pi-subagents` entry in the agent directory's `settings.json` to
`<agentDir>/npm/node_modules/pi-subagents`. Those are different roots whenever
`HOME` is not the operator's real home.

So an operator with an isolated or sandboxed home, or with
`PI_CODING_AGENT_DIR` set, gets the setup-incomplete block and its step 5 telling
them to export the variable, although Pi can already find the packages. The
doctor probes were fixed for this class in #817; the launch path was not, and the
hint is still printed from `internal/cli/pi.go:472`.

## Visible value

An operator whose home is not the default runs `spacedock pi` and both
extensions resolve from the registration Pi already reads, with neither
package-root variable exported.

## Out of scope

- Pi's own resolution rules.
- The isolated-home harness contract, which is folded into the Pi lane work.
- Installing the packages, and the plugin and skill paths.

## Expected surface and tolerance

`internal/cli/pi.go`: the launch-path resolution and its source label.
`internal/cli/pi_frontdoor_test.go`: the resolution tests. Estimate net +80,
across 2 files. Tolerance net +40 to +140, at most 3 files.

## Acceptance criteria

**AC-1 (VALUE) - Both roots resolve from the registration with no variable exported.**
Verified by: a test that registers both npm entries in a settings.json inside a
temporary agent directory, points `HOME` elsewhere, exports neither variable, and
requires both roots to resolve under that agent directory's npm path.
Falsifier: remove the registration read and that test must fail.

**AC-2 - An exported override still wins.**
Verified by: the existing override tests. Falsifier: ignore the variable and they
must fail.

**AC-3 - A missing install still fails closed.**
Verified by: the setup-incomplete path, with no registration and no variables.
Falsifier: return a path that does not exist and the test must fail.

**AC-4 (no-regression) -** The doctor's probe semantics are unchanged - the same lines (`pi-subagents extension`, `pi-subagents skill`, `pi-intercom package root`, `pi-intercom skill`) test the same conditions and keep their OK/MISSING meaning - and the `Pi runtime setup incomplete` message text is byte-identical. Because the doctor shares `piRuntimeConfigFromEnv`, it resolves the same registered roots as the launch path; only the no-override/no-registration (default-home) case is byte-identical to today. The offline suite passes.

## Test plan

Deterministic tests only. No live run. Every test must be able to fail for a
reason other than editing its own expectation.

## Stage Report: implementation

- DONE: Implement the registered-extension resolution in the Pi launch path, exactly as the entity's ideation stage report designs it.
  `resolveRegisteredPiRoot` (internal/cli/pi.go) resolves a settings.json `packages` entry to its root+source; `piRuntimeConfigFromEnv` applies the three cases (env override -> registration -> default home). Commit 55f3a48ff. Note: the entity has no `## Stage Report: ideation` section (headings: Problem/Visible value/Out of scope/Expected surface/Acceptance criteria/Test plan); the functions, three cases and source-label field were taken from the entity's Problem + Expected surface + dispatch brief, and no product decision was outstanding.
- DONE: Stay inside the stated tolerance: net +40 to +140 across at most 3 files.
  net +121 (148 insertions / 27 deletions) across 2 files: internal/cli/pi.go, internal/cli/pi_frontdoor_test.go (`git diff --numstat 894ff5dfd..HEAD`).
- DONE: State the falsifier for each acceptance criterion you claim, and run it.
  Each falsifier below was executed against a temporarily mutated pi.go, observed RED, then reverted (clean gofmt).
- DONE: AC-1 (VALUE) - Both roots resolve from the registration with no variable exported.
  Test `TestPiRuntimeConfigResolvesRegisteredRootsWithoutEnv`: registers `npm:pi-subagents`+`npm:pi-intercom` in a temp agent dir, HOME elsewhere, no vars -> both roots resolve under `<agentDir>/npm/node_modules`, sources `settings.json npm:*`. Falsifier (own test): neutering `resolveRegisteredPiRoot` fell both roots back under HOME and the test went RED.
- DONE: AC-2 - An exported override still wins.
  Pre-existing `TestPiRuntimeConfigResolvesEnvPathsForSubagentsIntercomAuthAndSessions` asserts `packageRootSource`/`intercomPackageSource` = the env var names. Falsifier (pre-existing test): forcing the env branch to be skipped made it resolve the default and go RED.
- DONE: AC-3 - A missing install still fails closed.
  Test `TestPiRuntimeConfigFailsClosedWithoutRegistrationOrOverride`: no vars, no settings.json -> default roots under HOME (asserted absent on disk) and `checkPiRuntime` reports extensionOK/intercomPackageOK false. Falsifier (own test): pointing the default at an existing dir made the test RED.
- DONE: AC-4 (no-regression) - doctor probe semantics unchanged; `Pi runtime setup incomplete` text byte-identical.
  Diff touches no probe `printPiCheck` line and no setup-message line; only additive `INFO ... (source: ...)` lines. Falsifier (pre-existing test): rewording the setup message to `Pi runtime incomplete.` made `TestPiInstallMissingSubagentsPrintsActionableInstructions` RED. Offline `go test ./internal/cli/` green (172.9s).
- DONE: Make the source label operator-visible.
  `piCheckResult` now carries `packageRootSource`/`intercomPackageSource`; `printPiDoctorReport` prints `INFO pi-subagents root: <root> (source: <source>)` (and the intercom twin) and the `Pi runtime ready` line prints `(source: ...)`. Test `TestPiDoctorPrintsResolutionSource`; falsifier (own test): dropping the INFO lines made it RED.
- DONE: Commit on the branch and push only that branch.
  55f3a48ff on `spacedock-ensign/launcher-resolves-registered-pi-extensions`; pushed with `git push origin HEAD:spacedock-ensign/...`. No push to main.

### Summary

The Pi launch path now uses the same package resolution as the host it launches: an exported `PI_SUBAGENTS_PACKAGE_ROOT`/`PI_INTERCOM_PACKAGE_ROOT` wins, else the package registered in the agent directory's `settings.json` (pi's own discovery), else the default home layout. The winning source label — previously populated but unprinted — is now carried into `piCheckResult` and shown in the doctor report and the ready line, so an operator whose HOME is not the default learns which root resolved. Net +121 across 2 files; no probe-line or setup-message bytes changed.
