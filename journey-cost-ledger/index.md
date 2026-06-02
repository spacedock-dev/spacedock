---
id: 4nef7frwwsrasfbntqfjd11c
title: Track tool calls and token costs for release journeys
status: ideation
source: "captain (2026-06-02) - track tool calls and token consumption for common user journeys; serialize per release; mark selected test cases for token/turn tracking"
started: 2026-06-02T16:53:20Z
completed:
verdict:
score: "0.29"
worktree:
issue:
---

Spacedock has behavior tests and live workflow tests, but it does not preserve a release-to-release view of agent cost. The missing artifact is a stable, redacted journey ledger: for selected user journeys, record tool calls, turns, token usage, duration, host/model metadata, and outcome, then publish the aggregate with each release.

## Problem

Pass/fail tests answer "does the journey work?" They do not answer "did this release make the same journey twice as expensive?" Today that cost signal is scattered across host transcripts, CI logs, and ad hoc observations. The result is that regressions in FO loop efficiency, dispatch shape, prompt size, or host adapter behavior are only visible when a human notices the session felt expensive.

The cost signal overlaps with tests, but it is not the same contract. Token counts and tool-call counts vary with host version, model, cache state, and transcript format. They should be measured from real journeys, serialized, and compared as trends. They should not become exact golden outputs in ordinary unit tests.

## Proposed approach

Add a journey metrics layer that attaches to selected tests and live journeys.

1. Tests opt in by declaring a stable journey ID and metadata, for example `fresh-install-claude`, `fresh-install-codex`, `fo-startup-status`, `one-entity-ideation`, `implementation-validation-cycle`, and `codex-runtime-dispatch`.
2. A small `internal/journeymetrics` package owns the schema, aggregation, and budget policy. Test helpers or live harnesses write one JSON record per marked journey into a known output directory.
3. Host-specific parsers extract cost from transcripts after the run. Start with Claude, reusing the existing usage-field knowledge in `internal/claudeteam/contextbudget.go`: assistant `usage.input_tokens`, `cache_creation_input_tokens`, and `cache_read_input_tokens`, plus tool-use and turn counts from transcript events. Add Codex only after its transcript schema is verified.
4. CI uploads raw per-run journey records as artifacts. The release process aggregates accepted records into `journey-costs-vX.Y.Z.json` and a short markdown summary showing deltas from the prior release.
5. Budgets are optional and coarse. By default a tracked journey records metrics and reports deltas without failing. A test may declare explicit ceilings such as `max_total_tokens` or `max_tool_calls` only for stable journeys where a budget violation should block.

The overlap with test cases is intentional: a tracked journey is still a normal test or live harness that proves behavior. The metrics layer observes that run and emits cost facts. Most unit tests stay untracked.

## Out of scope

- Storing raw transcripts in release assets. CI may upload transcripts as temporary artifacts for diagnosis, but release ledgers contain aggregate counts and metadata only.
- Exact token-count golden tests. The ledger compares trends and optional ceilings, not byte-for-byte cost totals.
- Replacing `spacedock status`, `dispatch context-budget`, or existing live E2E pass/fail logic.
- Codex transcript parsing before the actual Codex transcript schema is exercised and recorded.

## Acceptance criteria

Each AC names an end-state property and a check outside this task body that can fail.

**AC-1 - Selected tests and live harnesses can opt into journey cost tracking without changing their behavioral assertions.**
Verified by: Go tests for `internal/journeymetrics` register at least two fake journeys, one ordinary unit-style journey and one live-harness-style journey, and assert that both still report behavior pass/fail independently from metrics emission.

**AC-2 - The Claude journey parser extracts turns, tool calls, and token usage from transcript fixtures.**
Verified by: fixture-based Go tests over representative Claude JSONL transcript snippets assert total turns, total tool calls, tool-call counts by name, input tokens, cache creation tokens, cache read tokens, and total tokens. The fixtures include at least one assistant entry with usage and one tool-use event.

**AC-3 - Release journey metrics serialize to a stable, versioned ledger schema.**
Verified by: a golden-file Go test feeds multiple journey records to the aggregator and asserts the resulting `journey-costs.json` contains `schema_version`, release/version metadata, host/model metadata, outcome, duration, turns, tool calls, token totals, and per-tool counts.

**AC-4 - The release/live CI path preserves journey metrics as artifacts.**
Verified by: a test or workflow-lint check confirms `.github/workflows/runtime-live-e2e.yml` uploads the journey metrics output directory, and release documentation names the release asset `journey-costs-vX.Y.Z.json`.

**AC-5 - Cost budgets are explicit and coarse, not hidden golden values.**
Verified by: Go tests over the budget policy assert that a journey with no budget never fails from cost drift, a journey with `max_total_tokens` fails or reports the configured violation when exceeded, and exact token equality is not required anywhere.

## Test plan

- Unit tests for the schema, aggregator, and optional budget policy in `internal/journeymetrics`.
- Fixture tests for Claude transcript parsing, using small checked-in JSONL snippets rather than live host output.
- A workflow or release-doc check that the live E2E job uploads journey metric files and release docs describe the versioned asset.
- No live Codex parser test until a Codex transcript fixture is captured and its schema is understood.

## Stage Report: ideation

- DONE: The overlap with tests is explicit.
  Tracked journeys are normal behavior tests or live harnesses with extra metrics metadata. The metrics layer observes selected runs; most tests remain untracked.
- DONE: The release artifact shape is specified.
  The design serializes aggregate records as `journey-costs-vX.Y.Z.json` plus a markdown summary, with raw transcripts kept out of release assets.
- DONE: Cost budgets are trend-oriented by default.
  Exact token golden tests are out of scope. Optional ceilings are allowed only when a journey is stable enough for a coarse budget.

### Summary

Designed a release journey cost ledger that reuses selected tests as measurement points without making token counts brittle test oracles. The implementation starts with a small metrics package, Claude transcript fixtures, JSON aggregation, CI artifact upload, and optional explicit budgets.
