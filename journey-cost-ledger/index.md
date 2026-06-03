---
id: 4nef7frwwsrasfbntqfjd11c
title: Track tool calls and token costs for release journeys
status: validation
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
