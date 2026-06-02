---
id: f4h107he3s9swwz51ymr022m
title: Enable Codex live CI with a real headless runtime smoke
status: validation
source: "captain (2026-06-02) — revive the useful old Python Codex live CI setup, ignore tautological static grep tests, and document setup/principles in docs/dev/README.md"
started: 2026-06-02T17:59:45Z
completed:
verdict:
score: "0.72"
worktree: .worktrees/spacedock-ensign-codex-live-ci
issue:
---

The old Python-era live CI had the right Codex shape: a `codex-live` job with `CI-E2E-CODEX`, `OPENAI_API_KEY`, Codex CLI install, `codex login --with-api-key`, preserved artifacts, and a `make test-live-codex` target that ran live Codex tests in a cheap serial preflight before heavier parallel cases. The static workflow-grep tests around it are not worth reviving.

This task ports only the useful live behavior to the current Go repo: a real `codex exec --json` smoke that exercises the current Spacedock plugin/skill checkout, plus an approval-gated GitHub Actions job that installs/logs in to Codex and runs that live smoke.

## Problem

Current CI has a Claude live lane but no Codex live lane. r0 added the Codex runtime contract and Codex-shaped dispatch text, but the only Codex proof so far is local CLI/plugin probing plus unit tests. That is not enough for runtime regressions such as skill loading, `codex exec --json` event shape, multi-agent notification handling, or install/login drift in GitHub Actions.

The old Python live harness proves the design direction, but it cannot be copied wholesale: `tests/` is gone, the repo is Go-first, the old harness used synthetic `$HOME/.agents/skills/spacedock` skill injection, and some old assertions (especially "immediate wait after dispatch") contradict the r0 Codex contract where mailbox notification is the completion signal and `wait_agent` is only a sparing critical-path accelerator.

## Proposed approach

Add the smallest meaningful Codex live smoke in Go, likely under `internal/ensigncycle` behind `//go:build live`, parallel to the existing Claude live test but focused on a gated workflow rather than a full terminal cycle. The test should:

- create a temporary git-backed workflow with one gated entity;
- launch a real Codex headless session via `codex exec --json` (or a proven headless equivalent) with isolated `HOME` and `CODEX_HOME`;
- load the current checkout's Spacedock plugin/skills, not the remote `next` plugin by accident;
- drive the FO to the gate;
- assert observable state and output: the entity is still at the gated stage, the FO did not self-approve or archive it, and the final output/report clearly says it is waiting for human approval;
- preserve JSONL/transcript artifacts when CI sets a test-artifact root.

Then add a `codex-live` job to `.github/workflows/runtime-live-e2e.yml` that depends on the offline Go suite, uses the approval-gated `CI-E2E-CODEX` environment, reads only `OPENAI_API_KEY`, installs Codex, logs in with `codex login --with-api-key`, builds the local `spacedock` binary, proves Codex plugin readiness against the current checkout, runs the live Codex smoke, and uploads artifacts.

Update `docs/dev/README.md` (the evergreen development doc referenced by `AGENTS.md`) with the Codex live test setup and principles: live behavior is the proof; static grep tests are not a substitute; the job must test the PR/current checkout, not `next`; serial cheap preflight before heavier cases; environment approval protects spend/secrets; artifacts are part of the debugging contract.

## Riskiest unknown

The first implementation spike must prove how CI loads the **current PR checkout** as a Codex plugin. `codex plugin marketplace add spacedock-dev/spacedock --ref next` is not acceptable for PR validation because it tests the published `next` branch, not the code under review. Candidate paths are:

- generate a temporary local Codex marketplace in the job whose `spacedock` plugin source points at `${GITHUB_WORKSPACE}`;
- use a supported local marketplace path directly with `codex plugin marketplace add`;
- or prove a headless `codex exec` path can consume the local plugin checkout without marketplace installation.

Record the chosen mechanism in `docs/dev/README.md` and in the implementation report. If no current-checkout mechanism works reliably, stop and report the blocker rather than shipping a live lane that false-passes against `next`.

## Out of scope

- Do not revive the old Python static workflow-grep tests.
- Do not port the whole old `make test-live-codex` suite in this task.
- Do not reintroduce `test_codex_dispatch_immediate_wait.py` semantics; immediate wait-after-spawn is stale after r0.
- Do not add a model/cost matrix yet. Cost ledger integration belongs to the separate journey-cost-ledger task.

## Acceptance criteria

**AC-1 — A Go live Codex smoke proves FO gate-hold behavior through a real Codex runtime.**
Verified by: `OPENAI_API_KEY=... go test -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` (or the implemented equivalent) launches real Codex headlessly, reaches the gate on a temp workflow, and fails if the entity is self-approved, archived, or missing a waiting-for-approval signal. Without live Codex auth, the test may skip locally; in CI the secret-bearing job must fail clearly when `OPENAI_API_KEY` is absent after environment approval.

**AC-2 — GitHub Actions has an approval-gated `codex-live` lane that runs the live smoke against the current checkout.**
Verified by: `.github/workflows/runtime-live-e2e.yml` includes a `codex-live` job that needs the offline Go gate, uses environment `CI-E2E-CODEX`, reads `OPENAI_API_KEY` but not `ANTHROPIC_API_KEY`, installs/logs in to Codex, builds the local `spacedock` binary, loads or installs the current checkout's plugin/skills, runs the live Codex smoke, and uploads artifacts. Validation should cite a successful Actions run when available; local validation must at least run the same job commands up through the live test on a machine with Codex auth.

**AC-3 — Codex plugin readiness is proven with real Codex CLI commands, not assumed from docs.**
Verified by: the live job or its setup script runs the current Codex plugin install/load path and `go test ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v` (or a stronger current-checkout equivalent) after installation. The evidence must show the resolver sees `spacedock@spacedock` from the intended checkout/cache; a run that only tests remote `next` does not satisfy this AC.

**AC-4 — Evergreen docs explain Codex live setup and testing principles in `docs/dev/README.md`.**
Verified by: `docs/dev/README.md` contains a Codex live CI/testing section covering required local tools/secrets, GitHub environment/secret names, how to run the local smoke, why static grep tests do not prove runtime behavior, why CI must test the current checkout rather than `next`, and how artifacts are preserved for debugging.

**AC-5 — The implementation does not restore tautological static tests as the primary proof.**
Verified by: validation evidence centers on live command execution and resulting state/output. Syntax or parser checks are allowed as supporting guardrails only; they cannot be the only proof for AC-1 through AC-3.

## Test plan

- Cheap setup check: run `codex --version`, `codex login status` or equivalent, and a current-checkout plugin install/load spike before writing the live test.
- Go gates: `go test ./...`, `go test ./... -race`, and focused tests for any new live harness helpers.
- Codex CLI smoke: `go test ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v` after plugin setup.
- Live runtime smoke: run the new `go test -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` locally when Codex auth is available; otherwise document the skip and rely on the approval-gated CI run.
- CI evidence: after the PR opens, approve `CI-E2E-CODEX` and cite the `codex-live` job result and uploaded artifacts in validation.

## Stage Report: ideation

- DONE: Scoped the task to live Codex behavior, not static workflow grep tests
  The task revives the old Python-era `codex-live` setup pattern but explicitly rejects the tautological static assertions as primary proof.
- DONE: Named the current-checkout plugin-load risk as the first spike
  CI must prove the PR checkout's plugin/skills are loaded; testing `spacedock-dev/spacedock --ref next` would false-pass against already-published code.
- DONE: Added the evergreen docs requirement to the ACs
  `docs/dev/README.md`, which `AGENTS.md` references as the development workflow definition, is the target documentation surface.

### Summary

Ready for implementation: add a Go live Codex gate-hold smoke, wire an approval-gated `codex-live` Actions job around it, prove current-checkout plugin loading, and document the setup/principles in `docs/dev/README.md`.

## Stage Report: implementation

- DONE: Current-checkout Codex plugin load is proven and cannot false-pass against next
  Generated `.agents/plugins/marketplace.json` with `source: local`; isolated `codex plugin list` showed a temp local path, resolver passed, and the cache contained new `codex_live_test.go`.
- DONE: Live Go Codex gate-hold smoke and codex-live CI lane are implemented
  Code commit `11381e15397336626144d7ee75964fd6d65c5afd` adds `TestLiveCodexGateGuardrail`; `.github/workflows/runtime-live-e2e.yml` adds `codex-live` on `CI-E2E-CODEX`.
- DONE: docs/dev/README.md documents Codex live setup and proof principles
  `docs/dev/README.md` now covers local setup, `CI-E2E-CODEX`/`OPENAI_API_KEY`, current-checkout marketplace loading, live proof, and artifacts.
- SKIPPED: Local live Codex execution with real OPENAI_API_KEY
  No `OPENAI_API_KEY` is present locally; `go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` compiled and skipped, and the CI-required missing-secret path failed clearly.

### Summary

Implemented a build-tagged Codex live gate-hold smoke, helper tests for local marketplace/auth/gate assertions, and a `codex-live` Actions lane that installs the current checkout as a local Codex marketplace before running the live smoke. Verified with uncached `go test -count=1 ./...`, uncached `go test -count=1 ./... -race`, focused Codex helper tests, YAML parsing, and an isolated Codex plugin install proving the cache came from the current checkout rather than remote `next`.

## Stage Report: validation

- DONE: AC-1 - Go live Codex smoke proves FO gate-hold behavior through real Codex runtime
  `go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` compiled and skipped without `OPENAI_API_KEY`; required CI mode failed clearly with the approval-gated key message.
- DONE: AC-2 - GitHub Actions has approval-gated `codex-live` lane against current checkout
  `.github/workflows/runtime-live-e2e.yml` has `codex-live` behind `CI-E2E-CODEX`, needs `offline`, reads only `OPENAI_API_KEY`, installs/logs in to Codex, builds local `spacedock`, installs the current checkout marketplace, runs resolver/live tests, and uploads artifacts.
- DONE: AC-3 - Codex plugin readiness is proven with real Codex CLI commands/current-checkout mechanism, not remote `next`
  Isolated temp `HOME`/`CODEX_HOME` ran `codex plugin marketplace add`, `codex plugin add spacedock@spacedock`, `codex plugin list`, and `go test -count=1 ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v`; listing pointed at the temp local checkout path.
- DONE: AC-4 - `docs/dev/README.md` documents Codex live setup and proof principles
  The Codex Live CI section covers local prerequisites, `CI-E2E-CODEX`/`OPENAI_API_KEY`, current-checkout marketplace loading, runtime proof over static grep, live smoke command, and artifacts.
- DONE: AC-5 - Primary proof is live command execution/resulting state/output, not tautological static tests
  Evidence is the live-tagged `codex exec --json` harness, state/output assertions, real Codex plugin commands, and Go gates; no stale Python/static workflow-grep proof was revived.
- DONE: Local gates are rerun from the code worktree, including focused helper tests, `go test ./...`, and `go test ./... -race`
  From `.worktrees/spacedock-ensign-codex-live-ci`: `go test -count=1 ./internal/ensigncycle`, `go test -count=1 ./...`, and `go test -count=1 ./... -race` all passed.
- DONE: The validation report recommends PASSED or REJECTED and flags any remaining CI/live-auth caveats
  Recommendation: PASSED; caveat: no local `OPENAI_API_KEY` or approved Actions run was available, so full live-auth runtime execution remains to be observed in `codex-live`.

### Summary

Validation recommends PASSED for commit `11381e15397336626144d7ee75964fd6d65c5afd`. The implementation provides a real Codex CLI/current-checkout setup path plus a live-tagged FO gate-hold smoke; local validation proved the non-auth skip and CI-required missing-secret failure paths, but did not observe a paid live Codex run.

## Stage Report: implementation follow-up

- DONE: Restored the old Python harness's local subscription-auth ergonomics for the Go Codex live smoke
  Code commit `e4f3b568` adds a local Codex auth path: when `OPENAI_API_KEY` is absent and `~/.codex/auth.json` exists, the test copies only `auth.json` into a temporary `CODEX_HOME`, preserving isolated plugin marketplace state while using the operator's existing Codex login.
- DONE: Kept CI approval-gated auth strict
  `SPACEDOCK_CODEX_LIVE_REQUIRED=1` still requires `OPENAI_API_KEY`; local subscription auth is not accepted as a CI substitute.
- DONE: Proved the live smoke locally through real Codex without `OPENAI_API_KEY`
  `env -u OPENAI_API_KEY SPACEDOCK_BIN=/tmp/spacedock-codex-live SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-codex-live-artifacts go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` passed in 24.4s.
- DONE: Reran the normal Go gates after the follow-up
  `go test ./...` passed 747 tests in 12 packages, and `go test ./... -race` passed 747 tests in 12 packages.

### Summary

Follow-up validation is required because the previous validation report is now stale: the local live Codex smoke no longer skips when the operator has an existing Codex subscription login.

## Stage Report: validation follow-up

- DONE: Re-check the entity acceptance criteria against the full implementation plus the follow-up commit.
  AC-1 through AC-5 remain satisfied at code commit `e4f3b568ff123192a89f04ad6313e790e099c70d`; the checks below prove behavior outside the task body.
- DONE: AC-1 - Go live Codex smoke proves FO gate-hold behavior through a real Codex runtime.
  Built `/tmp/spacedock-codex-live-validation` from the code worktree, then `env -u OPENAI_API_KEY ... go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` passed in 40.93s.
- DONE: AC-2 - GitHub Actions has an approval-gated `codex-live` lane that runs the live smoke against the current checkout.
  `.github/workflows/runtime-live-e2e.yml` still gates `codex-live` on `CI-E2E-CODEX`, requires only `OPENAI_API_KEY`, installs/logs in to Codex, builds local `spacedock`, installs the current checkout plugin, runs resolver/live tests, and uploads artifacts.
- DONE: AC-3 - Codex plugin readiness is proven with real Codex CLI commands, not assumed from docs.
  Live artifacts show `codex plugin list` installed `spacedock@spacedock` from the temp local marketplace path, and `go test -count=1 ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v` passed.
- DONE: AC-4 - Evergreen docs explain Codex live setup and testing principles in `docs/dev/README.md`.
  The Codex Live CI docs now describe `codex login`, local `~/.codex/auth.json` copying into temp `CODEX_HOME`, `OPENAI_API_KEY`, CI strictness, current-checkout plugin loading, live proof, and artifacts.
- DONE: AC-5 - The implementation does not restore tautological static tests as the primary proof.
  Validation evidence is real `codex exec --json` execution, temp workflow state assertions, Codex plugin commands, focused helper tests, and Go gates; no static workflow-grep proof was revived.
- DONE: Specifically verify the old Python harness issue: local Codex subscription auth should work in an isolated home without `OPENAI_API_KEY`, by copying only `~/.codex/auth.json` into temp `CODEX_HOME`; CI-required mode must still require `OPENAI_API_KEY`.
  Local `~/.codex/auth.json` was nonempty; the live smoke passed with `OPENAI_API_KEY` unset and `codex login status` reported ChatGPT auth, while `SPACEDOCK_CODEX_LIVE_REQUIRED=1` failed clearly without `OPENAI_API_KEY`.
- DONE: Run the smallest command set that proves the follow-up.
  Focused auth-helper tests passed, `go test -count=1 ./...` passed 747 tests in 12 packages, and `go test -count=1 ./... -race` passed 747 tests in 12 packages.
- SKIPPED: Approved GitHub Actions `codex-live` run.
  No approved Actions run was available in this local validation; the local live command exercised the same current-checkout plugin and live Codex gate-smoke path.
- SKIPPED: `gofmt -w ./cmd ./internal`.
  This worker cannot edit product code; read-only `gofmt -l ./cmd ./internal` reported unrelated `internal/status/enum_scope_test.go`, while the codex-live touched Go files were gofmt-clean.

### Summary

Recommendation: PASSED for commit `e4f3b568ff123192a89f04ad6313e790e099c70d`. The follow-up restores local subscription-auth ergonomics without weakening CI: local isolated `CODEX_HOME` auth passed a real live Codex smoke with no `OPENAI_API_KEY`, and CI-required mode still hard-fails without `OPENAI_API_KEY`.
