---
title: Add Pi to runtime-live CI and shared runtime scenarios
status: implementation
source: captain (2026-06-05) — PRs that affect Pi runtime support need a Pi live lane using OPENAI_API_KEY; local/CI coverage should run most shared scenarios through Pi, either LLM-live or codified where a scenario is not yet live-safe
score: "0.33"
started: 2026-06-05T00:00:00Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-live-ci-runtime-scenarios
issue:
id: z8wm1jdhzmq4rcm0rzxw7bf5
---

# Add Pi to runtime-live CI and shared runtime scenarios

Add Pi as a first-class runtime-live CI lane so PRs that change Pi runtime support can be verified by the same GitHub workflow discipline as Claude and Codex. The lane should use `OPENAI_API_KEY` for Pi's OpenAI-backed runtime path, run against the current checkout, and exercise as much of the shared runtime scenario suite as the current Pi substrate can honestly support.

## Problem

The current `.github/workflows/runtime-live-e2e.yml` has an offline gate plus Claude and Codex live lanes. Pi runtime work is now significant enough that PRs can be green without any live Pi verification. Recent Pi work proved individual mechanisms locally — Pi subagent dispatch, fresh-context stage dispatch, dispatch-build artifact wrapping, and child-to-FO intercom talkback — but the CI workflow does not yet burn a Pi environment to prove the current checkout under the shared runtime scenarios.

This leaves PRs such as Pi readiness/doctor changes under-verified: offline tests can prove CLI fixtures, but not that `spacedock pi` can launch a real Pi FO and drive workflow state in CI. The runtime-support principle says runtime claims require live or fixture-backed durable state evidence, not transcript phrasing or setup checks alone.

## Proposed approach

1. **Add a Pi live CI lane.** Extend `.github/workflows/runtime-live-e2e.yml` with a `pi-live` job gated behind a Pi/OpenAI environment, likely `CI-E2E-PI`, using `OPENAI_API_KEY`. Keep the existing offline-gate dependency so Pi spend only happens after secret-free tests pass.
2. **Install and launch current-checkout Pi resources.** In CI, install Pi and the required local Pi extensions/skills (`pi-subagents`, `pi-intercom` where needed), configure authentication from `OPENAI_API_KEY`, and ensure `spacedock pi` or the live test harness loads the current checkout's Spacedock skills rather than remote `next`.
3. **Run Pi live coverage for shared scenarios where supported.** Add a Pi runner for shared runtime scenarios where the current Pi substrate can produce durable workflow state. Prefer running the existing shared scenarios through a real Pi LLM lane. If a scenario is not yet live-safe on Pi, add an honest codified/fixture-backed Pi variant or explicitly mark the gap with a follow-up task rather than pretending coverage.
4. **Preserve host-neutral scenario semantics.** Keep scenario definitions host-neutral. Pi-specific details belong in runner setup/adapter code, not in the shared scenario table.
5. **Upload useful Pi artifacts.** Store Pi session logs, live artifacts, journey metrics, and final output under `live-artifacts/pi/` so failures are debuggable like Claude/Codex.

## Acceptance criteria

**AC-1 - Runtime Live E2E has a Pi live lane.**
Verified by: `.github/workflows/runtime-live-e2e.yml` defining a `pi-live` job that depends on the offline gate, uses `OPENAI_API_KEY`, is environment-gated for Pi live spend, and uploads Pi-specific artifacts.

**AC-2 - Pi live setup uses the current checkout and required Pi substrates.**
Verified by: workflow commands and/or tests showing Pi installs/loads local Spacedock first-officer/ensign skills plus required Pi substrate resources (`pi-subagents` and, when supervisor talkback is required, `pi-intercom`) from deterministic paths, not a remote `next` install.

**AC-3 - Shared runtime scenarios include Pi coverage for supported scenarios.**
Verified by: Go tests or live-tagged runner coverage checks that fail if supported shared scenarios lack a Pi runner, while explicitly documenting any scenario that remains unsupported or codified-only on Pi.

**AC-4 - Pi live/codified scenario assertions grade durable workflow state.**
Verified by: tests asserting Pi scenario outcomes through entity frontmatter/body, state checkout git log, artifacts, or other durable state; transcript wording alone must not satisfy a Pi scenario.

**AC-5 - Missing Pi live prerequisites fail loudly after approval.**
Verified by: CI/workflow or live-test setup checks that emit clear errors when `OPENAI_API_KEY`, Pi CLI, Pi auth/config, `pi-subagents`, or required intercom prerequisites are missing in the live lane.

**AC-6 - Existing Claude/Codex live lanes are not regressed.**
Verified by: workflow structure and focused tests preserving the existing offline, Claude, and Codex jobs, secrets, environments, and artifact paths.

## Test plan

- Inspect existing Runtime Live E2E workflow and shared scenario runner structure before editing.
- Add/extend workflow or host-neutrality tests that parse `.github/workflows/runtime-live-e2e.yml` and fail if `pi-live` is missing required dependency, environment, secret, setup, run, or artifact steps.
- Add/extend `internal/ensigncycle` coverage tests so Pi runner coverage is explicit for supported shared scenarios, with honest skip/gap metadata for unsupported scenarios.
- Run focused tests for workflow/runner coverage, then `go test ./... -count=1`.
- If local Pi prerequisites are available, run the smallest local Pi live/codified scenario smoke before claiming live readiness; otherwise record the exact missing prerequisite and leave live execution to CI approval.

## Residual risks

- The exact CI install/auth path for Pi with `OPENAI_API_KEY` may need iteration; do not treat first-contact setup friction as proof the lane is impossible.
- Running the full shared scenario suite through Pi may expose real runtime gaps. Unsupported scenarios should be tracked honestly rather than hidden.
- Live Pi spend should remain environment-gated like Claude/Codex.

## Stage Report: implementation

### Summary
Implemented a Pi Runtime Live E2E lane and codified current Pi shared-scenario coverage status. Product commit: `05ba04af ci: add pi runtime live lane`.

### Checklist
- DONE: Added `pi-live` to `.github/workflows/runtime-live-e2e.yml` behind the `offline` gate and `CI-E2E-PI`, using `OPENAI_API_KEY`, `SPACEDOCK_PI_LIVE_REQUIRED=1`, and Pi-specific artifact paths.
- DONE: Added CI setup for `pi-coding-agent`, `pi-subagents`, `pi-intercom`, current-checkout `SPACEDOCK_REPO_ROOT`, and a `spacedock doctor --host pi --plugin-dir "$GITHUB_WORKSPACE"` verification step.
- DONE: Added Pi shared-scenario coverage metadata and parity checks that require every shared scenario to have an explicit Pi `live`/`codified`/`gap` entry; the existing scenarios are honestly marked `gap` until live-safe Pi shared runners exist.
- DONE: Kept durable Pi assertions in `TestLivePiFrontDoorSmoke` / smoke helpers: process success, entity body marker/stage report, state git log commit, and clean state checkout path.
- DONE: Taught Pi doctor/live auth checks to accept `OPENAI_API_KEY` as CI auth, while still copying local Pi OAuth auth when available.
- DONE: Preserved Claude/Codex lanes and extended workflow guard tests so Pi additions do not replace their shared scenario runs or artifact uploads.

### Validation
- `go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestPiSharedScenarioCoverage|TestSharedRuntimeScenarioDefinitions|TestPiLiveSmokePromptRequiresExactStageReportHeading' ./internal/ensigncycle -v` — PASS.
- `go test ./internal/cli ./internal/release ./internal/ensigncycle -count=1` — PASS.
- `go test ./... -count=1` — PASS.
- `go test -tags live -count=1 -run 'TestSharedScenarioRunnerCoverage|TestPiSharedScenarioCoverage|TestLivePiFrontDoorSmoke' ./internal/ensigncycle -v` — PASS, including local live Pi front-door smoke.
- `go test ./... -race` — PASS.

### Open risks
- Pi shared scenarios are not yet full live LLM runners for `gate-guardrail`, `rejection-flow`, or `merge-hook-guardrail`; the implementation records them as explicit gaps rather than pretending coverage.
- First CI execution may still expose package-install or auth-environment friction for `pi-coding-agent` / `pi-subagents`; the lane is environment-gated and fails loudly after approval.
