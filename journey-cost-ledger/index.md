---
id: 4nef7frwwsrasfbntqfjd11c
title: Track tool calls and token costs for release journeys
status: implementation
source: "captain (2026-06-02) - track tool calls and token consumption for common user journeys; serialize per release; mark selected test cases for token/turn tracking"
started: 2026-06-02T16:53:20Z
completed:
verdict:
score: "0.29"
worktree: .worktrees/spacedock-ensign-journey-cost-ledger
issue:
---

Spacedock has behavior tests and live workflow tests, but it does not preserve a release-to-release view of agent cost. The missing artifact is a stable, redacted journey ledger: for selected user journeys, record tool calls, turns, token usage, duration, host/model metadata, and outcome, then publish the aggregate with each release.

## Problem

Pass/fail tests answer "does the journey work?" They do not answer "did this release make the same journey twice as expensive?" Today that cost signal is scattered across host transcripts, CI logs, and ad hoc observations. The result is that regressions in FO loop efficiency, dispatch shape, prompt size, or host adapter behavior are only visible when a human notices the session felt expensive.

The cost signal overlaps with tests, but it is not the same contract. Token counts and tool-call counts vary with host version, model, cache state, and transcript format. They should be measured from real journeys, serialized, and compared as trends. They should not become exact golden outputs in ordinary unit tests.

## Proposed approach

Add a journey metrics layer that attaches to selected tests and live journeys.

1. Tests opt in by declaring a stable journey ID and metadata, for example `fresh-install-claude`, `fresh-install-codex`, `fo-startup-status`, `one-entity-ideation`, `implementation-validation-cycle`, and `codex-runtime-dispatch`. Codex journey IDs are allowed in v0 only after the runner has captured a representative `codex-exec.jsonl` fixture and classified the record as `metrics_state: characterized`; they do not get budget ceilings or trend claims until a parser AC promotes them to measured records.
2. A small `internal/journeymetrics` package owns the schema, aggregation, and budget policy. Test helpers or live harnesses write one JSON record per marked journey into a known output directory.
3. Host-specific parsers extract cost from transcripts after the run. Start with Claude, reusing the existing usage-field knowledge in `internal/claudeteam/contextbudget.go` and the live `internal/ensigncycle/testdata/sonnet_teamdelete_hang.stream.jsonl` shape. Claude assistant rows can be split across multiple JSONL rows with the same `message.id`; count the message-level `usage` object once per unique `message.id`, count turns by unique assistant message ID, and count tool calls by unique `tool_use.id` with per-tool names from the content block. For stream-json runs that end with a terminal `result` row, `result.usage` is the run-level token aggregate and `modelUsage` is the per-model cost/token aggregate. Do not add assistant-row usage to terminal `result.usage`; keep assistant-dedup totals only as parser diagnostics or as a fallback when the terminal result row is absent.
4. The existing Codex live runner already writes `codex-exec.jsonl` under `live-artifacts/codex/...`. The implementation first captures a small checked-in characterization fixture from that output and asserts the fields Spacedock can rely on before any Codex token parser is considered in scope.
5. CI uploads raw per-run journey records as artifacts. A release-tool command, for example `spacedock-release journey-costs X.Y.Z --metrics-dir <dir> --out <path>`, validates accepted records and writes `journey-costs-vX.Y.Z.json`. The release job must run that command, fail if the file is absent or empty, and publish the generated file as a release asset; documentation alone is not proof.
6. Budgets are optional and coarse. By default a tracked journey records metrics and reports deltas without failing. A test may declare explicit ceilings such as `max_total_tokens` or `max_tool_calls` only for stable journeys where a budget violation should block.

The overlap with test cases is intentional: a tracked journey is still a normal test or live harness that proves behavior. The metrics layer observes that run and emits cost facts. Most unit tests stay untracked.

## Out of scope

- Storing raw transcripts in release assets. CI may upload transcripts as temporary artifacts for diagnosis, but release ledgers contain aggregate counts and metadata only.
- Exact token-count golden tests. The ledger compares trends and optional ceilings, not byte-for-byte cost totals.
- Replacing `spacedock status`, `dispatch context-budget`, or existing live E2E pass/fail logic.
- Codex token-cost parsing before a captured `codex-exec.jsonl` fixture proves the schema fields to consume.

## Acceptance criteria

Each AC names an end-state property and a check outside this task body that can fail.

**AC-1 - Selected tests and live harnesses can opt into journey cost tracking without changing their behavioral assertions.**
Verified by: Go tests for `internal/journeymetrics` register at least two fake journeys, one ordinary unit-style journey and one live-harness-style journey, and assert that both still report behavior pass/fail independently from metrics emission.

**AC-2 - The Claude journey parser extracts turns, tool calls, terminal usage, and per-model usage without double-counting split assistant rows.**
Verified by: fixture-based Go tests over representative Claude JSONL transcript snippets assert total turns, total tool calls, tool-call counts by name, input tokens, output tokens, cache creation tokens, cache read tokens, total tokens, total cost, and per-model `modelUsage`. The fixtures include split assistant rows that share the same `message.id` but carry repeated `message.usage`; the expected totals count that `usage` once for assistant diagnostics and use terminal `result.usage` as the run aggregate when present. A second fixture without a terminal result asserts the assistant-dedup fallback path.

**AC-3 - Codex v0 journey coverage is characterization-first, not an unverified cost claim.**
Verified by: a checked-in `codex-exec.jsonl` fixture captured from `TestLiveCodexSharedScenarios` or the existing runner shape is parsed by a characterization test that records the stable event kinds and fields available for journey attribution. The test asserts that Codex journey IDs serialize with `metrics_state: characterized` until a later parser AC supplies token totals; budget policy rejects Codex token ceilings while that state is active.

**AC-4 - Release journey metrics serialize to a stable, versioned ledger schema.**
Verified by: a golden-file Go test feeds multiple journey records to the aggregator and asserts the resulting `journey-costs.json` contains `schema_version`, release/version metadata, host/model metadata, outcome, duration, turns, tool calls, token totals, and per-tool counts.

**AC-5 - The release/live CI path preserves journey metrics as artifacts and produces the release ledger through a checked release-tool path.**
Verified by: tests for `cmd/spacedock-release journey-costs` feed fixture journey records and assert it writes `journey-costs-vX.Y.Z.json`, rejects empty accepted input, and fails on an output filename that does not match the release version. A workflow check parses `.github/workflows/runtime-live-e2e.yml` and `.github/workflows/release.yml` to assert the live job uploads the metrics directory, the release job runs the journey-cost builder before publishing, and the generated JSON path is included in the release assets.

**AC-6 - Cost budgets are explicit and coarse, not hidden golden values.**
Verified by: Go tests over the budget policy assert that a journey with no budget never fails from cost drift, a journey with `max_total_tokens` fails or reports the configured violation when exceeded, and exact token equality is not required anywhere.

## Test plan

- Unit tests for the schema, aggregator, and optional budget policy in `internal/journeymetrics`.
- Fixture tests for Claude transcript parsing, using small checked-in JSONL snippets. Include a split-row same-`message.id` case, a terminal `result.usage`/`modelUsage` case, and a no-terminal-result fallback case.
- A Codex characterization fixture test over captured `codex-exec.jsonl` output. This verifies the v0 schema surface and marks Codex records characterized, not measured, until a parser is added.
- Release-tool tests for `cmd/spacedock-release journey-costs`, plus a workflow check that the live E2E job uploads metrics and the release job builds and publishes `journey-costs-vX.Y.Z.json`.
- No Codex token parser or Codex budget gate until the characterization fixture proves the stable fields to consume.

### Feedback Cycles

- Cycle 1 (2026-06-03): validation rejected AC-5. Detached audit replaced the release builder with commented command text plus `{}` output, and `TestWorkflowsPreserveAndPublishJourneyCosts` still passed. Route back to implementation to make the workflow guard parse executable release steps or otherwise prove the release job actually runs `spacedock-release journey-costs`, not merely contains matching substrings in comments.
- Cycle 2 (2026-06-03): captain review rejected the sample ledger shape. Refer to `m8` (`scenario-testing-principles`): a scenario is the natural-language spec, and codified tests plus LLM runs are executor implementations of that same spec; cost/coverage key by `scenario × {mode, runtime}`. The ledger keyed logical journeys as host-prefixed IDs such as `claude-gate-guardrail` and `codex-gate-guardrail`, but `8y` already established `gate-guardrail`, `rejection-flow`, and `merge-hook-guardrail` as host-neutral seed scenarios. Route back to implementation so cost tracking is keyed by the shared scenario/spec ID, for example `gate-guardrail`, with mode, runtime/executor, host, and model represented as variants or observations under that scenario. Preserve Claude `measured` versus Codex `characterized` states, but do not make those states separate logical journeys.
- Cycle 3 (2026-06-03): captain review rejected the `mode` axis usage. `sonnet`, `opus`, and `gpt-5-codex` are models, not modes. Route back to implementation for a narrow fix: define `mode` as execution/evidence mode (for example `llm-live`, `codified`, and possibly future `replay`), keep runtime as host (`claude`/`codex`), keep model in `model`, and update golden/tests so model names never populate `mode`.
- Cycle 3 addendum (2026-06-03): captain clarified scope priority: first ensure per-run metric records are archived as CI artifacts from the live scenario jobs. Aggregation and release/post-processing can remain a tested consumer if already present, but it should not be the core acceptance point for this slice. The immediate proof is that CI uploads the raw journey-metrics JSON artifacts for the shared scenario runs.

## Stage Report: ideation

- DONE: The overlap with tests is explicit.
  Tracked journeys are normal behavior tests or live harnesses with extra metrics metadata. The metrics layer observes selected runs; most tests remain untracked.
- DONE: The release artifact shape is specified.
  The design serializes aggregate records as `journey-costs-vX.Y.Z.json` plus a markdown summary, with raw transcripts kept out of release assets.
- DONE: Cost budgets are trend-oriented by default.
  Exact token golden tests are out of scope. Optional ceilings are allowed only when a journey is stable enough for a coarse budget.

### Summary

Designed a release journey cost ledger that reuses selected tests as measurement points without making token counts brittle test oracles. The implementation starts with a small metrics package, Claude transcript fixtures, JSON aggregation, CI artifact upload, and optional explicit budgets.

## Stage Report: ideation (cycle 2)

- DONE: Specify Claude transcript parsing so split assistant rows with the same message id do not double-count `usage`, and terminal `result.usage` / `modelUsage` handling is explicit.
  Proposed approach item 3 and AC-2 now pin same-`message.id` dedupe, terminal `result.usage` as the run aggregate, and `modelUsage` as the per-model aggregate.
- DONE: Resolve the Codex v0 scope: either add a captured `codex-exec.jsonl` characterization AC or remove Codex journey IDs from the release-ledger claim.
  Proposed approach item 4 and AC-3 add the captured Codex characterization fixture path; Codex records stay `metrics_state: characterized` until token parsing is proven.
- DONE: Strengthen release proof so `journey-costs-vX.Y.Z.json` is produced by a release-tool/job check, not only uploaded or mentioned in docs.
  Proposed approach item 5 and AC-5 require `spacedock-release journey-costs` plus a release-workflow check that builds and publishes the generated JSON asset.

### Summary

Reworked the ideation gate findings without changing the frontmatter or touching sibling entities. The spec now separates Claude measured parsing, Codex v0 characterization, and release-ledger production through a concrete release-tool path.

## Stage Report: implementation

- DONE: Implement `internal/journeymetrics` opt-in records, aggregation schema, and coarse budget policy with tests for AC-1, AC-4, and AC-6.
  Commit 2b2e97d4 adds opt-in tracking, budget evaluation, versioned ledger aggregation, and focused tests including the golden ledger.
- DONE: Implement Claude transcript parsing fixtures and Codex characterization fixtures so AC-2 and AC-3 distinguish measured versus characterized journeys.
  Commit 2b2e97d4 adds Claude split-row/result fixtures and a Codex `codex-exec.jsonl` characterization fixture.
- DONE: Wire the release/live CI artifact path through checked release-tool/workflow tests so AC-5 proves `journey-costs-vX.Y.Z.json` is produced and publishable.
  Commit 2b2e97d4 adds `spacedock-release journey-costs`, live metrics artifact env/upload paths, release asset upload, and workflow guards.

### Summary

Implemented the release journey cost ledger as a checked Go metrics package, release-tool command, host fixture parsers, and live/release workflow path. Verified with `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and a live-tag compile check for the live runner hooks.

## Stage Report: validation

- DONE: Reproduce the implementation evidence for AC-1 through AC-6, including Go tests, release-tool behavior, workflow guards, and Codex characterization state.
  Focused AC command passed 10 tests across `internal/journeymetrics`, `cmd/spacedock-release`, and `internal/release`; live-tag compile check passed 1 test.
- DONE: Run the required high-stakes detached adversarial audit for CI/release machinery and record any material test-strength findings.
  Material: a detached audit replaced the release builder with commented command text plus `{}` output, and `TestWorkflowsPreserveAndPublishJourneyCosts` still passed.
- DONE: Write a validation report with a PASSED or REJECTED recommendation and evidence citations for every AC.
  Recommendation: REJECTED. AC-1/2/3/4/6 passed by focused tests; AC-5 is rejected because the workflow guard can green-light a release job that never invokes the builder.

### Summary

Recommendation: REJECTED. AC-1 passed via `TestTrackingOptInDoesNotOwnBehaviorOutcome`; AC-2 passed via the two Claude parser fixture tests; AC-3 passed via `TestCodexFixtureIsCharacterizedNotMeasured`; AC-4 passed via `TestAggregateLedgerMatchesGolden`; AC-6 passed via `TestBudgetPolicyIsExplicitAndCoarse`.

AC-5 nominal release-tool checks passed via `TestJourneyCostsCommand*` and the workflow guard test, but the required detached audit found a material test-strength hole: substring checks are satisfied by comments, so they do not prove the release workflow actually runs `spacedock-release journey-costs`. Full gates: `go test ./...` passed 886 tests in 13 packages; `go test ./... -race` passed 886 tests in 13 packages. Formatting note: `gofmt -l ./cmd ./internal` reports `internal/hostneutrality/prose_inflator_locks_test.go`, `internal/status/fence_helpers_test.go`, and `internal/status/mutate.go`, none touched by commit `2b2e97d4`; validation left code unchanged.

## Stage Report: implementation (cycle 2)

- DONE: Update the implementation so the release/workflow guard parses executable release workflow steps or otherwise proves the release job actually runs the journey-cost builder and publishes `journey-costs-vX.Y.Z.json`, not just comments or inert text.
  Commit 9ade930d adds an executable-step guard that parses `run: |` blocks, ignores comments, joins shell continuations, and matches command prefixes.
- DONE: Add or strengthen tests so the adversarial edit described above fails.
  Commit 9ade930d adds tests for commented-out builder plus fake `{}` output and commented-out publish command.
- DONE: Run focused tests needed for AC-5 and broader tests as appropriate.
  `go test ./internal/release ./cmd/spacedock-release` passed 20 tests; `go test ./...` and `go test ./... -race` each passed 888 tests.

### Summary

Fixed the AC-5 guard so matching substrings in comments no longer count as release proof. The release workflow test now verifies executable builder and publish commands for `journey-costs-v${RELEASE_VERSION}.json` and keeps the prior release-tool command checks.

## Stage Report: validation (cycle 2)

- DONE: Reproduce AC-5 after the cycle-2 fix: release workflow guard must ignore comments/inert text and prove executable builder and publish commands.
  Focused command passed 12 tests across `internal/release`, `cmd/spacedock-release`, and `internal/journeymetrics`; AC-5 guard tests include commented builder and commented publish rejection.
- DONE: Rerun the detached adversarial audit by commenting out or faking the journey-cost builder/publish path and confirm the tests fail.
  Throwaway audit clone failed `TestWorkflowsPreserveAndPublishJourneyCosts` on commented builder/fake `{}` with `no executable journey-cost builder command`, and failed publish-only edit with `no executable journey-cost release upload command`.
- DONE: Re-check the full AC set and write a validation report with PASSED or REJECTED recommendation and evidence citations.
  Recommendation: PASSED. AC-1/6 via tracking and budget tests; AC-2 via Claude parser fixtures; AC-3 via Codex characterization fixture; AC-4 via ledger golden; AC-5 via release command/workflow tests plus detached audit.

### Summary

Recommendation: PASSED. `go test ./...` and `go test ./... -race` each passed 888 tests across 13 packages; the focused AC command passed 12 tests across 3 packages, and detached adversarial edits now fail for executable-builder and executable-publish proof.

Formatting note: `gofmt -l ./cmd ./internal` still reports `internal/hostneutrality/prose_inflator_locks_test.go`, `internal/status/fence_helpers_test.go`, and `internal/status/mutate.go`; validation left code unchanged rather than rewriting pre-existing unrelated formatting drift.

## Stage Report: implementation (cycle 3)

- DONE: Re-key the ledger around a host-neutral scenario/spec ID, with executor/run variants carrying mode, runtime, host, model, and metrics_state.
  Commit 80e41ace changes records to `scenario_id`, groups release ledgers under `scenarios[].observations[]`, and carries mode/runtime/executor/host/model per observation.
- DONE: Update golden/release-tool tests so an 8y seed scenario sample is one logical scenario with Claude and Codex observations, not two host-prefixed journeys.
  Commit 80e41ace updates the golden ledger and release-tool test so `gate-guardrail` has one scenario and two Claude/Codex observations.
- DONE: Preserve measured versus characterized behavior, including Codex no-budget policy, and rerun focused plus full Go gates.
  Focused tests passed 27 tests; live-tag compile path passed 5 tests; `go test ./...` and `go test ./... -race` each passed 889 tests across 13 packages.

### Summary

Reworked the cost ledger around host-neutral scenario identity instead of host-prefixed journey IDs. The release artifact now represents Claude measured and Codex characterized runs as observations of the same seed scenario, while the Codex characterized token-budget rejection remains covered.

## Stage Report: validation (cycle 3)

- DONE: Verify the ledger schema and golden output use one host-neutral 8y seed scenario with Claude and Codex observations, not host-prefixed logical journey IDs.
  `TestAggregateLedgerMatchesGolden`, `TestJourneyCostsCommandWritesVersionedLedger`, and golden JSON inspection confirm one `gate-guardrail` scenario with Claude/Codex observations and no `claude-gate-guardrail` or `codex-gate-guardrail` logical IDs.
- DONE: Re-check measured versus characterized behavior, especially Codex characterized token-budget rejection and absence of Codex token-cost claims.
  `TestCodexFixtureIsCharacterizedNotMeasured` passed; a throwaway audit failed as expected when Codex characterized records claimed `Tokens.Total=1`, while the golden keeps Claude measured token/cost fields only.
- DONE: Rerun focused scenario-key tests plus full repo gates, and write a PASSED or REJECTED validation report with AC coverage.
  Recommendation: PASSED. Focused metrics/scenario, live-tag parity, release workflow, `go test ./...`, and `go test ./... -race` commands all passed.

### Summary

Recommendation: PASSED. AC-1 via `TestTrackingOptInDoesNotOwnBehaviorOutcome`; AC-2 via the two Claude parser fixture tests; AC-3 via `TestCodexFixtureIsCharacterizedNotMeasured`; AC-4 via `TestAggregateLedgerMatchesGolden`; AC-5 via `TestJourneyCostsCommand*` plus `TestWorkflowsPreserveAndPublishJourneyCosts`; AC-6 via `TestBudgetPolicyIsExplicitAndCoarse`.

The focused scenario-key run passed 11 tests across `internal/journeymetrics`, `cmd/spacedock-release`, and `internal/ensigncycle`; live-tag parity added 1 no-spend coverage test, and the release workflow guard passed directly. A throwaway audit under the assigned worktree failed as expected when the ledger grouped by variant keys (`scenario_count: 2`) and when Codex characterized records claimed tokens. Full gates passed across the repo. Formatting note: `gofmt -l ./cmd ./internal` still reports `internal/hostneutrality/prose_inflator_locks_test.go`, `internal/status/fence_helpers_test.go`, and `internal/status/mutate.go`; validation left code unchanged per validation-only scope.

## Stage Report: implementation (cycle 4)

- DONE: Replace model-name values in mode with execution/evidence mode values such as llm-live or codified, while preserving model in model.
  Commit e46306e8 adds `ModeLLMLive`/`ModeCodified`, removes model fallback into `mode`, and keeps Codex characterization filling only `model`.
- DONE: Update golden ledger and tests so Claude/Codex observations keep one scenario_id and runtime host, with model names absent from mode.
  Commit e46306e8 updates the golden ledger, journey metrics tests, release-tool samples, and live metric hooks to use `mode: llm-live` with `runtime: claude/codex`.
- DONE: Rerun focused journey-metrics/release tests plus full Go gates and write the implementation report.
  Focused packages passed 107 tests; live-tag no-run compile passed; `go test ./...` and `go test ./... -race` each passed 890 tests after formatting.

### Summary

Fixed the journey cost ledger axes so `mode` describes evidence/execution mode, `runtime` remains the host, and model names stay in `model`. The schema helpers no longer synthesize `mode` from `model`, preventing Codex characterization and aggregation paths from reintroducing the rejected shape.

## Stage Report: implementation (cycle 5)

- DONE: Make Runtime Live E2E archive raw per-run journey metric JSON records from the live shared scenario jobs as the primary proof for this slice.
  Commit 16cee94a adds a structural runtime-live workflow guard requiring active Claude/Codex shared-scenario runs and upload-artifact steps that include `live-artifacts/journey-metrics/**`.
- DONE: Keep aggregation and release/post-processing as tested consumers rather than the core acceptance point.
  Commit 16cee94a preserves the existing release workflow consumer checks while moving `TestWorkflowsPreserveAndPublishJourneyCosts` to call the live raw-metrics artifact guard first.
- DONE: Preserve the mode/model correction: mode is execution/evidence mode, while model carries sonnet, opus, gpt-5-codex, and related model names.
  Commit e46306e8 remains in the branch; the final scan only finds Codex characterization filling `model`, not `mode`.
- DONE: Rerun focused journey-metrics/release tests plus full Go gates and write the implementation report.
  Focused packages passed 110 tests; live-tag no-run compile passed; `go test ./...` and `go test ./... -race` each passed 893 tests after `gofmt -w ./cmd ./internal`.

### Summary

Adjusted the implementation evidence around the captain addendum: Runtime Live E2E now has a parsed workflow guard proving raw journey metric records are uploaded from both live shared scenario lanes. The release ledger builder remains covered as a downstream consumer, and the mode/model axis fix is preserved.
