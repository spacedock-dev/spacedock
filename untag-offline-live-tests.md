---
id: dve63qqe4dxsktfwbyp1327c
title: Untag 12 offline-nature tests wrongly gated //go:build live in internal/ensigncycle
status: backlog
source: "Two audits this session (live E2E suite composition; unwired-test archaeology) found 33 live-tagged test functions in internal/ensigncycle, only 12 reachable from any CI live lane's -run selector. Of the 21 unwired, 12 are not real live scenarios at all -- they are offline-nature unit/guard tests wrongly marked //go:build live, so they never run anywhere, ever, at zero live spend and zero flakiness risk if untagged: TestClaudeTODOModelScope, TestClaudeRejectionFlowTODOModelScope, TestClaudeSonnetGateGuardrailTODOModelScope, TestCleanupKeepMovingRootRetainsOnlyFailures, TestCodexLiveRunnerExecArgvEnablesMultiAgentV2, TestCodexLiveRunnerUsesSpacedockFrontDoorBeforeHostArgs, TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle, TestShallowBootFixtureContainsOnlyHeldGate, TestPiLiveSmokePromptRequiresExactStageReportHeading, TestPiLiveEnvDropsForeignRuntimeMarkers, TestPiLiveEnvScrubsAmbientPiSubagentMarkers, TestPiIntercomPackageRootDefaultsBesideSubagents. Captain directed: file tasks to reconcile the gaps found, with the testing doc (docs/runtime-live-ci.md) as the SSOT to keep current -- not a documentation-only entity, a real code change that includes updating that doc."
started:
completed:
verdict:
score: 0.7
worktree:
issue:
---

Remove the `//go:build live` constraint from 12 test functions in `internal/ensigncycle` that are offline-nature (no real model call, no live host launch) but are currently gated behind the `live` build tag, meaning they never run in any CI job -- not `go test ./...`, not any live lane. Untagging costs zero live spend and lets the existing offline `go test ./...` job start running them.

## Proposed approach

For each of the 12 named tests: confirm it genuinely makes no live model call or host launch (read the test body; if it constructs fixtures/assertions purely and calls no `liveDriver`/host-launch path, it qualifies), remove its file's `//go:build live` constraint (or move the test to a non-live-tagged file if the file also contains genuine live tests), and add it to `docs/runtime-live-ci.md`'s test inventory noting it now runs under the standard offline suite.

## Acceptance criteria

**AC-1 (VALUE)** -- All 12 named tests execute under `go test ./...` (no `-tags live`, no live credentials, no CI live-lane selector needed).
Verified by: `go test ./internal/ensigncycle/... -run '<name>'` (no live tag) passing for each, where today it reports "no tests to run" (build-excluded).

**AC-2** -- No behavior change to what any test actually asserts -- this is a build-constraint fix only, not a rewrite.
Verified by: diff review showing only the build-tag/file-placement change, no assertion edits; full `go test ./...` and `-race` green.

**AC-3** -- `docs/runtime-live-ci.md` is updated to list these 12 as offline tests (not part of the live-lane estate), so the doc stays the accurate source of truth for what's live vs. offline.
Verified by: the doc's test inventory reflects the post-change reality; a future reader does not need to re-run this session's audit to know these are offline.

## Test plan

Offline only. `go test ./internal/ensigncycle/...` and full `go test ./...` / `-race` before and after, confirming the 12 newly run and nothing else changes.

## Out of scope

Any of the 9 genuinely-live, currently-unwired scenario tests (separate entities/decisions), the ac2 re-anchor tautology (already filed separately), the TODO-skip liveness-resolution mechanism (separate, needs the state checkout).
