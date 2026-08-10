---
id: pma2c1j7wmg9yvf5n25tx2ff
title: Make local Codex live auth failure self-correcting
status: backlog
source: "Captain fast-track direction on 2026-08-09 after a targeted local gate-guardrail run incorrectly set the CI-only required flag and bypassed supported isolated OAuth."
started:
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
---

Make the existing local Codex OAuth path obvious when a developer accidentally uses the CI-only live-test flag.

## Problem

The local Codex live harness already copies `~/.codex/auth.json` into an isolated `CODEX_HOME`. A local command that sets `SPACEDOCK_CODEX_LIVE_REQUIRED=1` instead selects the protected-CI policy and fails with an API-key requirement. The current message does not tell the developer to remove the CI-only flag, so a targeted local run can be escalated unnecessarily to GitHub Actions.

## Proposed approach

Keep all authentication behavior unchanged. Improve the existing required-lane error so it names the CI-only variable and tells a local developer to unset it to use isolated OAuth. Add a focused unit assertion for the message and one short clarification beside the existing local command in `docs/runtime-live-ci.md`.

## Out of scope

- No environment-variable rename.
- No authentication-policy or decision-order change.
- No new credential copying or storage behavior.
- No model-backed test, manual Runtime Live E2E dispatch, or CI workflow change.

## Acceptance criteria

**AC-1 - A mistaken local CI-only flag gives the corrective local command.**
Verified by: the existing Codex auth-decision unit test requires the fatal message to name `SPACEDOCK_CODEX_LIVE_REQUIRED` and instruct the operator to unset it to use local Codex OAuth; removing either instruction makes the test fail.

**AC-2 - The documented local path cannot be confused with protected CI.**
Verified by: `docs/runtime-live-ci.md` keeps the targeted local command free of `SPACEDOCK_CODEX_LIVE_REQUIRED` and directly states that the harness copies `~/.codex/auth.json` into an isolated `CODEX_HOME`; the existing local-auth isolation tests continue to pass.

**AC-3 - Authentication behavior is unchanged.**
Verified by: `TestDecideCodexLiveAuth`, `TestSeedCodexLocalAuthCopiesOnlyAuthIntoIsolatedHome`, and `TestSeedCodexLiveHomeCopiesOnlyAuthAndMinimalConfig` pass without a decision-order or credential-path change.

## Test plan

Run the focused Codex auth-decision and isolated-home tests first, then `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. No live model test is required because this task changes only a deterministic diagnostic and documentation.
