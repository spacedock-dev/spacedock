---
title: Pi readiness checks should report intercom supervisor-talkback prerequisites
status: implementation
source: captain (2026-06-04) — follow-up from cq pi-intercom-supervisor-talkback spike after live passed evidence showed the behavior can work but current Pi doctor/install only check base runtime setup
score: "0.32"
started: 2026-06-04T00:00:00Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-runtime-readiness-intercom-prereqs
issue:
id: qj2m6f3v8n9c0p4r7t1x5kz2
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
