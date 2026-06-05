---
title: Pi readiness checks should report intercom supervisor-talkback prerequisites
status: validation
source: captain (2026-06-04) — follow-up from cq pi-intercom-supervisor-talkback spike after live passed evidence showed the behavior can work but current Pi doctor/install only check base runtime setup
score: "0.32"
started: 2026-06-04T00:00:00Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs
issue:
id: qj2m6f3v8n9c0p4r7t1x5kz2
mod-block: merge:pr-merge
---

# Pi readiness checks should report intercom supervisor-talkback prerequisites

Implement the product follow-up from `pi-intercom-runtime-capability-probe` (`cq`, Pi intercom supervisor-talkback spike): extend Spacedock's Pi readiness surfaces so operators can see whether the local Pi setup has the prerequisites for supervisor talkback.

This is a real implementation task, not another spike. The spike and live evidence should be used as input, but setup readiness checks must still be described honestly: `spacedock doctor --host pi` or `spacedock install --host pi` can verify/report prerequisites, but they do **not** by themselves prove that a live child subagent can send progress, block for a decision, receive the first officer reply, resume, and write durable evidence.

## Problem

Current product code in `internal/cli/pi.go` checks the Pi CLI, Pi auth, `pi-subagents` extension and skill paths, and local Spacedock first-officer/ensign skills. It does not report `pi-intercom` / supervisor-talkback readiness even though recent Pi runtime work now depends on child workers being able to contact the supervising session through `contact_supervisor`.

The cq spike established two important boundaries:

- Setup evidence, such as package path or bridge readiness, is necessary for supervisor talkback but is not capability proof.
- Live behavior proof requires the whole progress update -> decision request -> FO reply -> child resume -> durable marker chain.

Operators need the setup side surfaced in `spacedock install --host pi` and `spacedock doctor --host pi` so missing prerequisites fail clearly before a launch/runtime attempt. At the same time, docs and command output must not over-claim that doctor/install success equals live talkback proof.

## Spike input / evidence to reference

Use the existing cq entity as the source of record, not as implementation code:

- State entity: `.spacedock-state/pi-intercom-runtime-capability-probe/index.md`
- Product probe commit: `81a11d066fa369b1e70e25c9e7d46cc2693824b7` (`Add pi intercom runtime capability probe`)
- Validation evidence commit: `a50bb349dbef34d0d39638c8747478cb5b4ade6e` (`Record pi intercom validation not-run evidence`)
- Live passed evidence commit: `ff4e7362982fbfac8001e2f8d89875b5f13ae373` (`Record passed pi intercom live evidence`)
- Live run id: `4395f7ae-5ff7-41c2-9f80-7845f6b57439`
- Passed evidence file: `docs/dev/_evidence/pi-intercom-runtime-capability-probe/2026-06-05-passed-live-talkback.json`
- Durable marker from the live evidence: `/tmp/pi-intercom-runtime-capability-probe-Avf2Rx/pi-intercom-smoke-marker.txt` containing `PI-INTERCOM-SMOKE-APPROVED`

The implementation should read that evidence to understand the prerequisite names and honest wording, but should add product readiness checks and tests in the normal product tree.

## Proposed approach

1. **Model Pi intercom readiness separately from live capability.** Extend the Pi runtime check data model with explicit fields for the supervisor-talkback setup prerequisites, likely including:
   - `pi-intercom` command or package availability, if that is the supported operator-facing executable.
   - `subagents-doctor` availability or equivalent bridge-health command, if that is the stable way to inspect bridge/intercom setup.
   - Resolved package/path environment inputs such as `PI_SUBAGENTS_PACKAGE_ROOT`, `PI_CODING_AGENT_DIR`, and any Pi intercom package-path variable required by the supported substrate.
   - Any readable extension/skill path needed for the intercom bridge to be loaded in Pi child sessions.
2. **Report readiness in doctor/install output.** Extend `spacedock doctor --host pi` and `spacedock install --host pi` / `--check` so missing intercom setup is visible, actionable, and reflected in exit behavior where the command is meant to be a readiness gate.
3. **Keep the claim boundary in command text.** Output should say these checks verify setup prerequisites for supervisor talkback. It should also say live talkback proof requires the documented runtime probe and durable evidence; doctor/install success alone is not that proof.
4. **Preserve existing Pi checks.** Do not regress existing checks for Pi CLI/auth, `pi-subagents` extension/skill, and Spacedock skills. Compatibility-first behavior matters: add the new surface without changing unrelated host behavior.
5. **Update runtime docs if needed.** `docs/runtime-support.md` should name the new readiness checks and keep the setup-vs-live-proof distinction aligned with the cq probe.

## Acceptance criteria

**AC-1 - `spacedock doctor --host pi` reports supervisor-talkback setup readiness.**
Verified by: CLI tests that run the doctor path with fake filesystem/PATH fixtures and assert output includes explicit Pi intercom/supervisor-talkback prerequisite rows, actionable remedies for missing setup, and the resolved paths or env-derived package roots used by the check.

**AC-2 - Doctor exit behavior is failable for missing required intercom setup.**
Verified by: CLI tests covering all-present and missing-intercom fixtures. The all-present fixture exits 0; missing required supervisor-talkback setup exits non-zero and names the missing prerequisite. If any intercom check is intentionally advisory rather than required, the test must assert that wording and the chosen exit behavior explicitly.

**AC-3 - `spacedock install --host pi` / `--check` uses the same readiness model.**
Verified by: CLI tests showing install check output and exit behavior match doctor for Pi intercom prerequisites, while non-check install prints clear next steps without mutating unrelated global Pi state.

**AC-4 - Package/path environment variables are fixture-tested.**
Verified by: unit tests over the Pi config/check code that set env vars for the Pi subagents/intercom package roots and auth/session dirs, then assert the exact resolved paths checked and printed in the doctor/install report.

**AC-5 - Docs and command output do not over-claim live talkback proof.**
Verified by: docs/runtime-support or integration invariant tests asserting doctor/install wording says setup checks are necessary but insufficient, and that live supervisor talkback proof still requires the cq-style probe/evidence chain. Tests should fail on wording that equates bridge-active or doctor success with `pi-intercom-supervisor-talkback` passed behavior.

**AC-6 - Existing Pi runtime readiness remains covered.**
Verified by: existing and new CLI tests still cover Pi CLI/auth, `pi-subagents` extension/skill, and Spacedock first-officer/ensign skill checks, with no regression in non-Pi doctor/install behavior.

## Test plan

- **Focused CLI fixture tests:** Add or extend tests around `runDoctorWithPi`, `runInitWithPi`, `checkPiRuntime`, and `printPiDoctorReport` using fake `piRuntimeOps` to simulate present/missing `pi`, auth, subagents paths, intercom paths, and doctor/bridge command availability.
- **Exit-code matrix:** Cover healthy, missing auth, missing `pi-subagents`, missing `pi-intercom`, and missing bridge-health command cases for both `doctor --host pi` and `install --host pi --check`.
- **Environment resolution tests:** Use fixture env values for `PI_SUBAGENTS_PACKAGE_ROOT`, `PI_CODING_AGENT_DIR`, `SPACEDOCK_REPO_ROOT`, and the selected intercom package/path variable. Assert the exact paths in the resulting config and command output.
- **Docs/output invariant tests:** Parse `docs/runtime-support.md` and/or golden command output to require the setup-vs-live-proof warning and references to the Pi intercom supervisor-talkback probe.
- **Regression gates:** Run `go test ./... -count=1`; run `go test ./... -race` if the implementation changes shared runtime/checking code; run `gofmt -w ./cmd ./internal` after Go edits.
- **Manual/live check (optional for this implementation):** Do not require a new live probe run to pass this task. If a live run is performed, record it as evidence using the existing cq probe schema and do not replace CLI fixture tests with transcript claims.

## Residual risks

- The stable operator-facing way to locate or invoke `pi-intercom` / `subagents-doctor` may need confirmation in product code. If the implementation discovers there is no stable command/path contract, pause for a design decision rather than inventing one silently.
- Existing live passed evidence captured `pi-intercom-supervisor-talkback` behavior but did not capture clean Pi CLI/subagents/intercom versions for every package, so readiness output should tolerate unknown versions while still checking required paths/commands.
- Making intercom setup a hard doctor/install failure could affect users who only use the Pi front door without supervisor talkback. If product direction wants a warning-only mode, that should be explicit and test-covered.

## Out of scope

- Re-running the cq live supervisor-talkback smoke as a required gate for this task.
- Changing `pi-subagents` or `pi-intercom` internals.
- Adding PR/mod behavior or new workflow semantics.
- Treating setup checks, bridge-active output, or doctor success as live talkback proof.
- Editing the cq entity except to read it for context.

## Stage Report: implementation

- Product worktree: `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs`
- Product branch: `spacedock-ensign/pi-runtime-readiness-intercom-prereqs`
- Product commit: `8c6a2e8683a4552b8a6c627d632f7463e48f4c55` (`cli: report pi supervisor talkback prerequisites`)
- Changed product files:
  - `internal/cli/pi.go`
  - `internal/cli/pi_frontdoor_test.go`
  - `docs/runtime-support.md`
- DONE: Extended the Pi runtime readiness model to resolve/report `PI_INTERCOM_PACKAGE_ROOT`, `PI_CODING_AGENT_SESSION_DIR`, `pi-intercom`, and `subagents-doctor` alongside the existing Pi CLI/auth/pi-subagents/Spacedock skill checks.
- DONE: Made doctor and install check readiness fail when required supervisor-talkback prerequisites are missing, while non-check install remains idempotent/instructive.
- DONE: Added command output and docs wording that setup checks are necessary but insufficient; live proof still requires the cq-style progress -> decision -> supervisor reply -> child resume -> durable marker probe.
- DONE: Updated CLI tests for healthy/missing Pi doctor output, install output/exit behavior, environment path resolution, and runtime-support documentation invariants.
- Validation commands:
  - `gofmt -w internal/cli/pi.go internal/cli/pi_frontdoor_test.go` — passed
  - `go test ./internal/cli -count=1` — passed
  - `go test ./... -count=1` — passed
  - `go test ./... -race` — passed
- AC coverage:
  - AC-1/AC-2: `TestPiDoctorReportsMissingAndHealthyRuntime` covers explicit supervisor-talkback rows, remedies, all-present exit 0, and missing prerequisites non-zero.
  - AC-3: `TestPiInstallAcceptedAndDoesNotUsePluginCommands`, `TestPiInstallMissingSubagentsPrintsActionableInstructions`, and `TestPiInstallCheckFailsForMissingSupervisorTalkbackPrerequisites` cover install and install `--check` behavior without plugin mutation.
  - AC-4: `TestPiRuntimeConfigResolvesEnvPathsForSubagentsIntercomAuthAndSessions` and `TestPiRuntimeConfigDefaultsIntercomAndAuthPathsUnderHome` pin env/default path resolution.
  - AC-5: command-output assertions plus `TestRuntimeSupportDocsKeepPiDoctorVsLiveTalkbackBoundary` pin the setup-vs-live-proof boundary.
  - AC-6: existing Pi CLI/auth/pi-subagents/Spacedock skill checks remain in the doctor/install fixture matrix and non-Pi plugin-dir behavior tests still pass.
- Residual risks:
  - The implementation treats `pi-intercom`, `subagents-doctor`, and `PI_INTERCOM_PACKAGE_ROOT` as required readiness prerequisites per this task; if the Pi substrate later changes its stable package/command names, these rows will need adjustment.
  - The readiness checks verify availability/path setup only and do not run the live `pi-intercom-supervisor-talkback` probe.

### Summary

Implemented Pi doctor/install supervisor-talkback prerequisite readiness reporting with required `pi-intercom`, `subagents-doctor`, and intercom package-root checks, preserving the cq setup-vs-live-proof boundary in CLI output, tests, and runtime docs.
### Stage Report: validation

- Product worktree: `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs`
- Product branch: `spacedock-ensign/pi-runtime-readiness-intercom-prereqs`
- Product commit validated: `8c6a2e8683a4552b8a6c627d632f7463e48f4c55` (`cli: report pi supervisor talkback prerequisites`)
- Recommendation: **REJECT / needs fix before merge**

#### AC evidence

- AC-1: Mostly satisfied. `spacedock doctor --host pi` now prints a `Supervisor-talkback setup prerequisites` section with `pi-intercom command`, `subagents-doctor bridge-health command`, and `pi-intercom package root` rows, plus remedies and the resolved intercom package root/auth/session paths.
- AC-2: Satisfied for supervisor-talkback prerequisites. Missing `pi-intercom`, `subagents-doctor`, or `PI_INTERCOM_PACKAGE_ROOT` fixture state exits non-zero and names the missing prerequisite.
- AC-3: **Not satisfied.** `spacedock install --host pi --check` does not consistently use the same readiness model/exit behavior as doctor. With all Pi/subagents/intercom prerequisites present but auth missing, install `--check` returns 0 and prints `Pi runtime ready`, while doctor returns 1 and reports `MISSING Pi auth` for the same fixture. This violates the AC requirement that install check output and exit behavior match doctor for readiness gating, and regresses the existing Pi auth readiness surface.
- AC-4: Partially satisfied. Unit tests cover env/default resolution for `PI_SUBAGENTS_PACKAGE_ROOT`, `PI_INTERCOM_PACKAGE_ROOT`, `PI_CODING_AGENT_DIR`, and `PI_CODING_AGENT_SESSION_DIR`; doctor output includes resolved paths. The tests do not strongly assert exact printed env-derived intercom/auth/session paths, but behavior is present.
- AC-5: Satisfied. CLI output and `docs/runtime-support.md` keep the necessary-but-insufficient setup-vs-live-proof boundary and reference the cq-style `pi-intercom-supervisor-talkback` probe/evidence chain.
- AC-6: **Not satisfied due to AC-3 auth mismatch.** Doctor still covers Pi auth, but install `--check` can report ready despite missing auth because the early `piRuntimeLaunchReady` path ignores `authOK`.

#### External/failable evidence for rejection

Using fake `pi`, `pi-intercom`, and `subagents-doctor` commands, valid package/skill paths, and no auth file:

```text
$ spacedock install --host pi --check
exit=0
Pi runtime ready.
  pi-subagents: <tmp>/pi-subagents
  pi-intercom: <tmp>/pi-intercom
  Spacedock skills: <worktree>
NOTE: These checks verify necessary supervisor-talkback setup prerequisites only; they are insufficient to prove live child talkback.
NOTE: Live proof still requires the cq-style progress -> decision -> supervisor reply -> child resume -> durable marker probe for pi-intercom-supervisor-talkback.

$ spacedock doctor --host pi
exit=1
Pi runtime check
OK pi CLI: <tmp>/bin/pi
MISSING Pi auth: <tmp>/home/.pi/agent/auth.json
  remedy: run `pi` login/auth flow; live tests copy this file into an isolated PI_CODING_AGENT_DIR
...
```

Root cause: `runInitWithPi` checks `piRuntimeLaunchReady(check)` before honoring `--check`, and `piRuntimeLaunchReady` intentionally excludes `authOK`; `piDoctorHealthy` includes auth. Therefore install `--check` bypasses the doctor readiness gate whenever launch prerequisites are present but auth is absent.

#### Validation commands

- `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs status --short --branch` — confirmed branch `spacedock-ensign/pi-runtime-readiness-intercom-prereqs` at `8c6a2e8683a4552b8a6c627d632f7463e48f4c55`.
- `git show --stat --patch 8c6a2e8683a4552b8a6c627d632f7463e48f4c55 -- internal/cli/pi.go internal/cli/pi_frontdoor_test.go docs/runtime-support.md` — inspected implementation and tests.
- `go test ./internal/cli -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test ./... -race` — passed.
- Manual CLI fixture with fake commands/package roots and missing auth — install `--check` exited 0 while doctor exited 1, as shown above.

#### Residual risks / requested fix

- Fix `install --host pi --check` so it uses the same readiness/exit model as doctor, including Pi auth, rather than short-circuiting through `piRuntimeLaunchReady`.
- Add a regression test for the missing-auth install `--check` case; the current tests pass because the healthy install path omits auth and does not exercise check-mode parity with doctor.

## Stage Report: implementation fixback

- Product worktree: `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs`
- Product branch: `spacedock-ensign/pi-runtime-readiness-intercom-prereqs`
- Product commit: `546207eaa80d` (`cli: align pi install check with doctor readiness`)
- Changed product files:
  - `internal/cli/pi.go`
  - `internal/cli/pi_frontdoor_test.go`
- DONE: Moved `install --host pi --check` ahead of the non-check launch-readiness shortcut so check mode always prints the doctor report and returns `piDoctorExit`, including Pi auth.
- DONE: Preserved non-check install behavior: it still reports `Pi runtime ready` using launch readiness when launch prerequisites are present, without making auth a non-check install blocker.
- DONE: Added a regression test for the missing-auth parity case so install check now fails with `MISSING Pi auth` instead of printing ready.
- Validation commands:
  - `gofmt -w internal/cli/pi.go internal/cli/pi_frontdoor_test.go` — passed
  - `go test ./internal/cli -count=1` — passed
  - `go test ./... -count=1` — passed
  - `go test ./... -race` — passed
- AC coverage:
  - AC-3: fixed; `install --host pi --check` now uses the same doctor readiness/exit model, including Pi auth.
  - AC-6: fixed for the rejected auth mismatch; Pi auth readiness remains a failing check-mode gate.
- Residual risks:
  - No new live Pi probe was run; this fix is covered by fixture tests and preserves the setup-vs-live-proof boundary.

### Summary

Fixed the validation rejection by making Pi install check mode use doctor readiness before the launch-readiness shortcut, and added missing-auth regression coverage.

### Stage Report: validation rerun

- Product worktree: `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs`
- Product branch: `spacedock-ensign/pi-runtime-readiness-intercom-prereqs`
- Product commits validated:
  - `8c6a2e8683a4552b8a6c627d632f7463e48f4c55` (`cli: report pi supervisor talkback prerequisites`)
  - `546207eaa80d` (`cli: align pi install check with doctor readiness`)
- Recommendation: **PASS / mergeable**

#### AC evidence

- AC-1: Satisfied. `doctor --host pi` prints a `Supervisor-talkback setup prerequisites` section with explicit `pi-intercom command`, `subagents-doctor bridge-health command`, and `pi-intercom package root` rows, actionable remedies in missing fixtures, and resolved auth/session/package paths.
- AC-2: Satisfied. Fixture tests cover all-present doctor exit 0 and missing runtime/prerequisite exit non-zero; missing intercom rows name the missing prerequisites and remedies.
- AC-3: Satisfied after fixback. `install --host pi --check` now enters the same doctor report/exit path before the non-check launch-readiness shortcut. The prior rejection was falsified with an external fake-PATH/manual fixture: with Pi/subagents/intercom setup present but auth absent, both `install --host pi --check` and `doctor --host pi` exited 1 and printed `MISSING Pi auth`; `install --check` no longer printed `Pi runtime ready`.
- AC-4: Satisfied. `TestPiRuntimeConfigResolvesEnvPathsForSubagentsIntercomAuthAndSessions` and `TestPiRuntimeConfigDefaultsIntercomAndAuthPathsUnderHome` assert env/default resolution for `PI_SUBAGENTS_PACKAGE_ROOT`, `PI_INTERCOM_PACKAGE_ROOT`, `PI_CODING_AGENT_DIR`, `PI_CODING_AGENT_SESSION_DIR`, auth path, session dir, extension path, and skill paths. The doctor output also prints the resolved paths.
- AC-5: Satisfied. CLI output and `docs/runtime-support.md` state that doctor/install checks are necessary setup checks but insufficient to prove live supervisor talkback, and preserve the cq-style progress -> decision request -> supervisor reply -> child resume -> durable marker proof boundary.
- AC-6: Satisfied. Existing Pi CLI/auth/pi-subagents/Spacedock skill checks remain in the readiness model and fixture matrix; non-Pi install/doctor plugin-dir behavior tests still pass.

#### External/failable evidence

Manual fake-PATH fixture with executable `pi`, `pi-intercom`, and `subagents-doctor`, readable `pi-subagents` extension/skill paths, readable `PI_INTERCOM_PACKAGE_ROOT`, Spacedock skills from the product worktree, and no auth file:

```text
$ spacedock install --host pi --check
exit=1
Pi runtime check
OK pi CLI: <tmp>/bin/pi
MISSING Pi auth: <tmp>/home/.pi/agent/auth.json
OK pi-subagents extension: <tmp>/pi-subagents/src/extension/index.ts
OK pi-subagents skill: <tmp>/pi-subagents/skills/pi-subagents
INFO Pi auth/session dirs: auth=<tmp>/home/.pi/agent/auth.json session=<tmp>/home/.pi/agent/sessions
Supervisor-talkback setup prerequisites
OK pi-intercom command: <tmp>/bin/pi-intercom
OK subagents-doctor bridge-health command: <tmp>/bin/subagents-doctor
OK pi-intercom package root: <tmp>/pi-intercom
OK Spacedock first-officer skill: <worktree>/skills/first-officer
OK Spacedock ensign skill: <worktree>/skills/ensign
NOTE: These checks verify necessary supervisor-talkback setup prerequisites only; they are insufficient to prove live child talkback.
NOTE: Live proof still requires the cq-style progress -> decision -> supervisor reply -> child resume -> durable marker probe for pi-intercom-supervisor-talkback.

$ spacedock doctor --host pi
exit=1
Pi runtime check
OK pi CLI: <tmp>/bin/pi
MISSING Pi auth: <tmp>/home/.pi/agent/auth.json
...
OK pi-intercom command: <tmp>/bin/pi-intercom
OK subagents-doctor bridge-health command: <tmp>/bin/subagents-doctor
OK pi-intercom package root: <tmp>/pi-intercom
```

Manual healthy fixture with the same prerequisites plus `<tmp>/home/.pi/agent/auth.json` exited 0 for `doctor --host pi` and printed all Pi/subagents/intercom/Spacedock skill rows as `OK` while retaining the live-proof warning.

#### Validation commands

- `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs status --short --branch` — confirmed branch `spacedock-ensign/pi-runtime-readiness-intercom-prereqs` at `546207eaa80d` with no product edits from validation.
- `git show --stat --oneline 8c6a2e8683a4552b8a6c627d632f7463e48f4c55` — inspected product implementation commit (`internal/cli/pi.go`, `internal/cli/pi_frontdoor_test.go`, `docs/runtime-support.md`).
- `git show --stat --oneline 546207eaa80d` — inspected fixback commit (`internal/cli/pi.go`, `internal/cli/pi_frontdoor_test.go`).
- `go test ./internal/cli -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test ./... -race` — passed.
- Manual fake-PATH missing-auth parity check for `install --host pi --check` and `doctor --host pi` — both exited 1 with `MISSING Pi auth`; prior rejection no longer reproduces.
- Manual fake-PATH healthy doctor check — exited 0 with supervisor-talkback prerequisite rows and live-proof warning.

#### Residual risks

- No new live Pi supervisor-talkback probe was run; this validation only confirms setup readiness behavior and fixture-backed CLI tests, consistent with the task scope.
- The implementation hard-codes `pi-intercom`, `subagents-doctor`, and `PI_INTERCOM_PACKAGE_ROOT` as required setup contracts; future Pi substrate naming changes will require product/test updates.

### Summary

Validated the implementation and fixback. The prior install-check/auth mismatch is fixed, all acceptance criteria are satisfied with fixture and manual CLI evidence, and the setup-vs-live-proof boundary remains explicit.
