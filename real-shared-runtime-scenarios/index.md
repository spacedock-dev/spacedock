---
id: 8y7yten220npj6g4kj4680p2
title: Extract real shared runtime scenarios for Claude and Codex live CI
status: ideation
source: "captain (2026-06-03) - follow-up from f4 codex-live-ci: make the shared scenario reuse real, not Codex-only"
started: 2026-06-03T07:09:59Z
completed:
verdict:
score: "0.68"
worktree:
issue:
mod-block:
---

f4 ports the old Python Claude/Codex journey overlap into the Codex live lane, but the Go implementation is still Codex-shaped: `codexSharedScenarios()` defines the scenario list and the current Claude live CI continues to run `TestLiveEnsignCycle` rather than consuming the same scenario table.

This task makes the sharing real. The scenario definitions should be host-neutral, with Claude and Codex runners implementing the same scenario IDs through their runtime-specific launch, auth, plugin, artifact, and transcript mechanisms.

## Problem

Runtime regressions should be caught once per user journey and then exercised by each supported host. Today the old shared journey intent is preserved for Codex, but Claude and Codex do not share a single Go scenario contract. That allows parity drift: a scenario can exist for Codex only, the Claude lane can keep testing a different full-cycle shape, and CI can appear broad while the actual common journey set is split.

## Proposed approach

Extract the f4 scenario list into a host-neutral table, such as `sharedRuntimeScenarios()`, that contains only runtime-neutral facts: scenario ID, old Python provenance, behavior intent, fixture/journey shape, timeout/cost class, and expected durable outcome.

Move host-specific details behind runner adapters:

- Codex runner: `codex exec --json`, isolated `CODEX_HOME`, local Codex marketplace/plugin install, Codex artifacts.
- Claude runner: `spacedock claude`, isolated Claude config/auth, team/session artifacts, current checkout plugin path.

Both runners should implement the same scenario IDs. Meta tests should fail if either runtime is missing a runner for a shared scenario. Assertions should prefer durable workflow state and output over transcript phrasing: frontmatter, archive/no-archive state, stage reports, exact fix markers, merge-hook refusal, and only durable user-facing final-message obligations such as a gate review and decision prompt.

## Out of scope

- Do not add a large new scenario matrix in this task.
- Do not add token/cost ledger serialization; that belongs to journey-cost-ledger work.
- Do not remove the existing Claude full-cycle smoke unless the shared scenario suite fully replaces its coverage and validation proves the replacement.

## Acceptance criteria

**AC-1 - Shared runtime scenarios are defined once in host-neutral code.**
Verified by: focused Go tests under `internal/ensigncycle` assert the shared scenario table includes at least `gate-guardrail`, `rejection-flow`, and `merge-hook-guardrail`, and that the table type does not encode Claude-only or Codex-only runner fields.

**AC-2 - Codex live tests consume the shared table.**
Verified by: `go test -tags live -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` iterates the host-neutral table and fails if any shared scenario lacks a Codex runner.

**AC-3 - Claude live tests consume the same shared table.**
Verified by: `go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v` or the implemented equivalent iterates the same host-neutral table and fails if any shared scenario lacks a Claude runner.

**AC-4 - CI proves both runtime lanes run the shared scenario suite.**
Verified by: `.github/workflows/runtime-live-e2e.yml` runs the shared Claude scenario suite in `claude-live` and the shared Codex scenario suite in `codex-live`, with environment-gated secrets and artifact upload preserved.

**AC-5 - Scenario assertions are behavior/state oriented, not transcript-shape tautologies.**
Verified by: unit or fixture tests exercise at least one negative case per shared scenario where a broken state transition, missing rejection route, or merge-hook bypass makes the assertion fail.

**AC-6 - Evergreen docs explain the shared-scenario contract.**
Verified by: `docs/dev/README.md` documents how to add a shared runtime scenario, what belongs in the host-neutral definition, what belongs in each runner, and the commands for local Claude/Codex live execution.

## Test plan

- Focused unit/meta tests for the shared scenario table and runner coverage.
- Focused negative assertion tests for gate hold, rejection routing, and merge-hook guard behavior.
- `go test ./internal/ensigncycle`
- `go test ./...`
- `go test ./... -race`
- Local live run for Codex when Codex auth is available.
- Local live run for Claude when Claude auth is available, or approved `CI-E2E`/`CI-E2E-OPUS` evidence if local auth is unavailable.
- Approved `CI-E2E-CODEX` evidence for the Codex lane before validation passes.
