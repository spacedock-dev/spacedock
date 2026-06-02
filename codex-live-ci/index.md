---
id: f4h107he3s9swwz51ymr022m
title: Enable Codex live CI with a real headless runtime smoke
status: ideation
source: "captain (2026-06-02) — revive the useful old Python Codex live CI setup, ignore tautological static grep tests, and document setup/principles in docs/dev/README.md"
started:
completed:
verdict:
score: "0.72"
worktree:
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
