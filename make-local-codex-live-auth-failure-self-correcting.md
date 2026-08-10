---
id: pma2c1j7wmg9yvf5n25tx2ff
title: Make local Codex live auth failure self-correcting
status: ideation
source: "Captain fast-track direction on 2026-08-09 after a targeted local gate-guardrail run incorrectly set the CI-only required flag and bypassed supported isolated OAuth."
started: 2026-08-10T03:00:30Z
completed:
verdict:
score: 0.95
sprint: durable-decisions
sprint-readiness: ready
group: runtime-live-ux
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:pma2c1j7wmg9yvf5n25tx2ff:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:pma2c1j7wmg9yvf5n25tx2ff-backlog-1
              briefing:
                id: briefing:pma2c1j7wmg9yvf5n25tx2ff:backlog:attempt-1:revision-1
                digest: sha256:027f0c349f06b44d4299c4eeb239938ef923bce635f16cc50e369c4f8b3706ee
                request-digest: sha256:043aac914ad272034a2af4c8e58768a17afe191c817495ea7a31dfe20b9d7b1f
                room-ref: ./make-local-codex-live-auth-failure-self-correcting/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pma2c1j7wmg9yvf5n25tx2ff:backlog:1
                briefing: briefing:pma2c1j7wmg9yvf5n25tx2ff:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T03:00:04.328183Z"
                decision: approve
                reason: The seed isolates the operator error, keeps authentication behavior unchanged, and defines deterministic proof with no live-runtime spend.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:pma2c1j7wmg9yvf5n25tx2ff:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:pma2c1j7wmg9yvf5n25tx2ff-ideation-1
              briefing:
                id: briefing:pma2c1j7wmg9yvf5n25tx2ff:ideation:attempt-1:revision-1
                digest: sha256:cd0596992278060b9c5fbec95bd550798ec5990717fb5d057ea80f86fcd0bd0f
                request-digest: sha256:4a6cc4f46d1f7bd6d31be0d68255171439531fa29697dbb9f27e0e1b333979ce
                room-ref: ./make-local-codex-live-auth-failure-self-correcting/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pma2c1j7wmg9yvf5n25tx2ff:ideation:1
                briefing: briefing:pma2c1j7wmg9yvf5n25tx2ff:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T03:17:52.296781Z"
                decision: approve
                reason: 'The design is the smallest sufficient correction: two files, no auth-policy change, no output comparison, and deterministic proof only.'
              application:
                target-stage: implementation
                state: pending
---

Make the existing local Codex OAuth path obvious when a developer accidentally uses the CI-only live-test flag.

## Problem

The local Codex live harness already copies `~/.codex/auth.json` into an isolated `CODEX_HOME`. A local command that sets `SPACEDOCK_CODEX_LIVE_REQUIRED=1` instead selects the protected-CI policy and fails with an API-key requirement. The current message does not tell the developer to remove the CI-only flag, so a targeted local run can be escalated unnecessarily to GitHub Actions.

## Proposed approach

Keep all authentication behavior unchanged. Change only the existing required-lane fatal message to:

> `OPENAI_API_KEY is required when SPACEDOCK_CODEX_LIVE_REQUIRED is set; unset SPACEDOCK_CODEX_LIVE_REQUIRED for a local run that uses isolated Codex OAuth`

The function's decision order and all four modes remain untouched. Per the captain's 2026-08-09 ruling, no test may search, contain, equal, snapshot, or otherwise compare help/error prose. The diagnostic change is therefore verified by direct review against the concrete wording approved at the ideation gate, while existing tests prove only the auth behavior they exercise.

Insert this single paragraph immediately after the existing local Codex command in `docs/runtime-live-ci.md`:

> Leave `SPACEDOCK_CODEX_LIVE_REQUIRED` unset for this local path. When no `OPENAI_API_KEY` is set, the harness copies `~/.codex/auth.json` into an isolated `CODEX_HOME`; if the variable is already set, run `unset SPACEDOCK_CODEX_LIVE_REQUIRED` first.

The simplest alternative was documentation alone. It is insufficient because the developer who has already set the flag sees the fatal test output, not necessarily the guide; the corrective diagnostic serves AC-1 at that failure point. A diagnostic alone would leave the documented local/protected-CI distinction implicit, so the one-paragraph clarification serves AC-2. No new runtime mechanism is introduced.

## Expected surface

- `internal/ensigncycle/codex_liveenv.go`: replace one diagnostic string; zero net lines.
- `docs/runtime-live-ci.md`: insert the paragraph above beside the existing local Codex command; at most 3 net inserted lines.

Exactly two files, no new files, at most 3 net inserted lines and 5 changed lines. The only observable semantic allowed to change is the fatal diagnostic text when `SPACEDOCK_CODEX_LIVE_REQUIRED` is non-empty and `OPENAI_API_KEY` is empty. Command grammar, stored formats, authentication authority, credential paths, auth-mode selection, decision order, and successful/skip runtime behavior must not change. Exceeding the file list, line budget, or semantic declaration requires returning to ideation.

## Out of scope

- No environment-variable rename.
- No authentication-policy or decision-order change.
- No new credential copying or storage behavior.
- No model-backed test, manual Runtime Live E2E dispatch, or CI workflow change.

## Acceptance criteria

**AC-1 - A mistaken local CI-only flag gives the corrective local command.**
Verified by: implementation review compares the single changed string literal directly with the concrete wording approved above and rejects a diff that omits the variable name, the `unset` command, or the isolated-OAuth destination. This is a Captain-authorized narrowing of proof: no automated command-output or prose-presence comparison is permitted.

**AC-2 - The documented local path cannot be confused with protected CI.**
Verified by: review of the concrete doc diff confirms that the targeted local command remains free of `SPACEDOCK_CODEX_LIVE_REQUIRED` and its adjacent paragraph gives both the isolated-auth behavior and corrective `unset` command. `TestSeedCodexLocalAuthCopiesOnlyAuthIntoIsolatedHome` independently observes the documented `auth.json` copy on disk and fails if the copy is absent or extra operator state appears.

**AC-3 - Authentication behavior is unchanged.**
Verified by: all cases in `TestDecideCodexLiveAuth` retain their existing modes, while `TestSeedCodexLocalAuthCopiesOnlyAuthIntoIsolatedHome` and `TestSeedCodexLiveHomeCopiesOnlyAuthAndMinimalConfig` observe that the isolated home still contains the copied `auth.json` and minimal config only. A changed mode, decision priority, credential path, or extra copied state makes these tests fail.

## Test plan

1. Focused deterministic proof: `go test ./internal/ensigncycle -run '^(TestDecideCodexLiveAuth|TestSeedCodexLocalAuthCopiesOnlyAuthIntoIsolatedHome|TestSeedCodexLiveHomeCopiesOnlyAuthAndMinimalConfig)$' -count=1`. This proves the preserved mode table and isolated-home contents without selecting a model-backed test. It intentionally does not inspect diagnostic output; direct implementation review verifies the approved wording.
2. Formatting and full proof: `gofmt -w ./cmd ./internal` followed by `go test ./...`. This catches package-wide regressions; the expected formatting diff is empty outside the two-file surface.
3. Race proof: `go test ./... -race`. This is the repository-required concurrency regression check and remains deterministic/offline.

No model-backed run, manual Runtime Live E2E run, GitHub workflow dispatch, or diagnostic-output comparison is needed or permitted for this task: the only runtime-output change is produced by the pure auth-decision function before any Codex process launches. No spike needed: the focused baseline already exercised the existing decision function and both on-disk isolation helpers successfully on 2026-08-09; the task adds no unverified parser, handoff, format, or runtime mechanism.

## Stage Report: ideation

- DONE: Prove the smallest sufficient change is diagnostic text plus one documentation clarification, with no authentication-policy change.
  The proposal fixes the error at its source and documents the local path while explicitly freezing the decision order, modes, credential paths, and authority.
- DONE: Name the exact files and a tight line budget, and preserve the existing isolated-OAuth behavior and tests.
  Expected surface is exactly two named files, at most 3 net inserted and 5 changed lines; the focused baseline passed all three named auth/isolation tests without adding a prose assertion.
- DONE: Define deterministic focused, full, and race proof with no model-backed run or GitHub workflow dispatch.
  The exact focused, full, and race commands passed offline; the plan excludes live hosts, model selection, manual E2E, workflow dispatch, and all diagnostic-output comparison.

### Summary

The design limits implementation to one corrective fatal message and one adjacent documentation paragraph. Under the captain's proof narrowing, wording is reviewed directly while existing mode and on-disk tests preserve isolated OAuth and protected-CI policy; deterministic full and race verification add no live-runtime work.
