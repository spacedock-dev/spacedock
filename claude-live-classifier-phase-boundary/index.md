---
title: Append Claude live detector evidence only after a primary failure
status: ideation
score: 0.7
source: Captain-approved recovery item 1 on 2026-07-10, after review of merged PR #490 showed classifyBootPreambleFailure scans the full transcript and can label later contract-lookup searches as boot-preamble failures before the scenario assertion is evaluated.
id: p4h6a5wcqe5ddkhnmrac1w9a
started: 2026-07-10T12:56:28Z
---

# Append Claude live detector evidence only after a primary failure

Preserve the Claude live harness's real runner or scenario failure, and add wrong-root or broad-search observations only as secondary context when that primary failure already exists.

## End value

A passing Claude shared scenario stays passing even when the full transcript contains a wrong-root or filesystem-hunt detector match. A failing scenario reports its existing assertion or runner error byte-for-byte as the primary evidence, followed by an explicitly secondary detector observation; the observation can never create, replace, relabel, or pre-empt the failure.

## Problem statement

Merged PR #490 calls `classifyBootPreambleFailure(stream, workflowRoot)` in `claudeLiveRunner.run` immediately after transcript capture and before the existing stall check, final-message extraction, and every scenario assertion. The classifier scans the entire stream and calls `t.Fatalf` on any wrong-root or broad-search match. It therefore gives detector evidence independent pass/fail authority and can suppress the failure the scenario actually owns.

That authority is the bug. The detector cannot reliably infer a phase from a full transcript, and this task will not add a phase boundary. Wrong-root and filesystem-hunt matches are still useful when a run has already failed, but they are only diagnostic observations. They must remain invisible when the runner and scenario pass.

The fix is reporting-only. It does not permit a contract hunt, weaken any scenario assertion, alter fixture state, retry a scenario, or restore #490's reverted retry experiment.

## Current call-site finding

The merged #490 surface has one relevant classifier call: `claudeLiveRunner.run` invokes it after writing `claude-stream.jsonl`. Runner failures occur next through the existing stall and final-message-extraction `t.Fatalf` calls; scenario failures occur later in the nine `runClaude*Scenario` functions. All share the same `*testing.T`, so one cleanup registered by the runner can observe the final subtest outcome without changing those failure sites.

No spike needed: the design relies only on Go's documented `testing.T.Cleanup` and `testing.T.Failed` behavior, the existing full-transcript detectors, and the existing single-`*testing.T` runner/scenario call chain. The implementation's fake-reporter tests exercise the cleanup decision and output order without a model run.

## Approaches considered

1. **Recommended: register detector evidence for failure-only cleanup.** After transcript capture, select the existing wrong-root or broad-search observation and register a cleanup. The cleanup logs `Additional diagnostic: ...` only when `t.Failed()` is already true. It has no `Error`, `Fail`, or `Fatal` capability, so the observation cannot create a red; the existing failure remains untouched and appears first.
2. **Return a diagnostic in `liveResult` and compose at every assertion.** This can preserve ordering, but it requires editing every Claude scenario failure site and separately handling runner failures that return no result. The broader surface creates omission risk with no additional end value.
3. **Split the transcript into boot and post-boot phases.** Rejected by the captain. A boundary would add parsing and classification semantics while still granting some detector matches independent failure authority. The desired contract needs neither.
4. **Log detector matches immediately.** This is a small change but violates the silence requirement: passing scenarios would emit misleading warnings despite having no primary failure.

## Proposed approach

### Primary-error-first registration

Rename and reduce `classifyBootPreambleFailure` to `detectClaudeLiveFailureDiagnostic`. It may retain the current deterministic evidence priority—wrong-root observation before broad-search observation—but it returns context only and never decides pass/fail.

Add a narrow diagnostic reporter helper with only these capabilities:

```go
type claudeLiveDiagnosticReporter interface {
    Cleanup(func())
    Failed() bool
    Logf(string, ...any)
}
```

`claudeLiveRunner.run` detects the observation after persisting the transcript and registers it before the existing stall and extraction checks. The cleanup runs after the runner and scenario finish:

- if `Failed()` is false, it emits nothing;
- if `Failed()` is true, it logs `Additional diagnostic: <existing detector evidence>`;
- it never calls a failure-producing method and never wraps or edits the primary error.

This covers both failure classes with one hook. A stall or extraction `t.Fatalf` unwinds through the registered cleanup. A returned `liveResult` keeps the cleanup registered while the existing scenario assertions run. The first failure line and its exact text remain owned by the existing call site; the cleanup's later log is visibly secondary context.

### Classifier decision

Delete the classification concept and rename the helper rather than add a phase boundary. The helper's two-detector body can remain structurally identical because selecting the most applicable observation is useful; its contract and caller change from fatal classification to failure-only diagnostic registration. Renaming is the smallest code delta that removes misleading pass/fail semantics while retaining offline detector-priority coverage.

### Output contract

Passing run with a detector match:

```text
<no detector output; test remains passing>
```

Failing run with a detector match:

```text
<exact existing runner or scenario failure>
Additional diagnostic: <wrong-root or broad-search detector evidence>
```

The diagnostic text describes evidence found in the full transcript; it does not assert a phase, change the test status, or become the scenario's failure label.

### Implementation scope and #490 preservation

The exact proposed product-test files are:

- rename `internal/ensigncycle/boot_preamble_classify_impl_test.go` to `internal/ensigncycle/claude_live_failure_diagnostic_impl_test.go`; rename/reduce the classifier and add the failure-only reporter registration helper;
- rename `internal/ensigncycle/boot_preamble_classify_test.go` to `internal/ensigncycle/claude_live_failure_diagnostic_test.go`; replace fatal-classification expectations with fake-reporter controls for pass silence, primary preservation, secondary ordering, and detector priority;
- update `internal/ensigncycle/claude_live_runner_test.go`; remove the classifier-owned `t.Fatalf` branch and register the observation before the existing runner failure checks.

Do not change `broad_search_detect_impl_test.go`, `wrong_root_detect_impl_test.go`, any scenario assertion, or any fixture. Preserve #490's absolute-root prompt signatures and commissioned markers in `shared_fixtures_test.go`; the root arguments in the Claude, Codex, and live-scenario adapters; the gate-stop fixture relocation; and `boot_discovery_test.go`'s 10 discoverable fixtures plus markerless zero-discovery control. Do not add retry files or a second launch attempt.

## Acceptance criteria

- **AC-1 (value — detector evidence has zero independent pass/fail authority).** In an offline three-case matrix using the same broad-search observation, **1/1 passing case remains passing with zero diagnostic lines**, while **2/2 failing cases preserve the exact injected primary error** (one runner error and one scenario assertion) and append exactly one `Additional diagnostic:` line afterward. Against merged #490, the passing case is made red and **0/2 injected primaries** survive because classification pre-empts them. *Test:* table-driven `TestClaudeLiveFailureDiagnosticIsSecondaryOnly` uses a fake reporter and asserts final failed state, exact primary bytes, log count, and order.
- **AC-2 (primary failure sites remain authoritative).** The existing stall, extraction, and scenario assertion messages are unchanged; with a detector observation present, each remains the first reported failure and the diagnostic follows only during cleanup. *Test:* `TestRegisterClaudeLiveFailureDiagnostic` covers pass, runner-fail, and scenario-fail reporter timelines, while the live-tag compile check proves the real `*testing.T` call site satisfies the narrow interface.
- **AC-3 (diagnostic selection remains deterministic without classification).** When both detector observations apply, wrong-root evidence is selected before broad-search evidence; when neither applies, no cleanup is registered and no output is possible. Neither result can fail a test. *Test:* renamed offline controls `TestDetectClaudeLiveFailureDiagnosticWrongRootTakesPriority`, `TestDetectClaudeLiveFailureDiagnosticBroadSearch`, and `TestDetectClaudeLiveFailureDiagnosticCleanStream` assert the selected observation and registration count.
- **AC-4 (#490 fixture and no-retry non-regression).** All 10 intended shared fixtures remain discoverable, the zero-discovery fixture remains markerless and undiscoverable, all eight non-boot-subject prompts retain their absolute workflow root, and the runner still launches exactly once. *Test:* existing discovery and prompt-signature tests stay green; focused source review confirms no fixture diff and no retry path, and the full/race gates exercise the unchanged suite.

## Test plan

Implementation begins with deterministic red controls in the renamed diagnostic test. A fake reporter records cleanup registration, failed state, and log order; it gives the test direct proof that a detector match cannot fail a pass and cannot mutate an injected primary error. The tests cost milliseconds, use synthetic stream lines already supported by the detector fixtures, and spend no model tokens.

Run the focused diagnostic and detector controls:

```bash
go test ./internal/ensigncycle -run 'Test(ClaudeLiveFailureDiagnostic|RegisterClaudeLiveFailureDiagnostic|DetectClaudeLiveFailureDiagnostic|DetectBroadSearchAtBoot|DetectWrongRootBoot)'
```

Run #490's discovery and shared-fixture controls:

```bash
go test ./internal/ensigncycle -run 'Test(SharedScenarioFixturesAreDiscoverable|ZeroDiscoverFixtureStaysUndiscoverable|SharedPrompt|ClaudePrompt|CodexPrompt)'
```

Then run the required repository gates and live-tag compile check:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
go test -tags live -run '^$' ./internal/ensigncycle/...
```

No live workflow run is required. The behavior under change is deterministic Go test reporting after a captured stream exists, not model behavior. No documentation-site change is needed because this affects developer-facing live-test diagnostics only; the exact output contract is recorded above.

### Feedback Cycles

- Cycle 1 (ideation gate): REVISE. The captain questioned whether a classifier is needed and accepted the smaller direction: retain wrong-root and filesystem-hunt extraction only as diagnostic context for an already-existing runner or scenario failure. Do not invent a phase boundary, do not make a late search independently red, do not replace or pre-empt the primary failure, and remain silent when the scenario passes. Rework the proposal, acceptance criteria, tests, title, and scope around that primary-error-first diagnostic contract while preserving #490's fixture fixes.

## Stage Report: ideation

- DONE: Define an observable phase boundary that distinguishes true boot orientation from later contract lookup without relying on model wording.
  A correlated successful `status --boot` tool result ends the prefix; event order, IDs, and result status define the boundary.
- DONE: Specify failure composition that preserves the scenario's primary assertion while attaching any hunt diagnostic, with red controls for early boot and late-search streams.
  AC-1 through AC-3 require early fatal classification, deferred late diagnostics, and primary-first output across four paired streams.
- DONE: Keep the merged #490 fixture fixes intact, name exact files and focused tests, and produce a complete ideation Stage Report with measurable end-value criteria.
  The scope names six classifier/runner files, preserves six fixture surfaces, and measures 2/2 early, 0/2 mislabeled late, and 2/2 primary assertions reached.

### Summary

The design bounds boot classification at the first correlated successful fixture boot result and treats later hunts as deferred diagnostics. It preserves #490's fixture hardening, keeps genuine early failures fatal, and proves the corrected failure order offline with paired red controls.

## Stage Report: ideation (cycle 2)

- DONE: Re-open merged PR #490's current call sites and determine the smallest diagnostic-only design. Treat the captain's decision as binding: no phase boundary and no classifier that independently changes pass/fail; wrong-root and broad-search observations may only be appended to an already-existing runner or scenario failure.
  AC-2/AC-3: the current single classifier call is replaced by one failure-only cleanup registration on the shared `*testing.T`; the reporter surface cannot fail a test and introduces no phase parser.
- DONE: Rewrite the entity's title, end value, problem statement, considered approaches, proposed approach, implementation scope, acceptance criteria, and test plan around primary-error-first composition. Require silence when the scenario passes, preservation of the exact primary error when it fails, and secondary diagnostics that cannot replace, relabel, or pre-empt it. Explicitly decide whether `classifyBootPreambleFailure` should be deleted, reduced, or renamed based on the smallest code delta.
  AC-1/AC-2: the body title and all design sections now require pass silence and byte-preserved primary errors; `classifyBootPreambleFailure` is renamed/reduced to evidence selection, with YAML title mutation left to the FO-owned status surface.
- DONE: Preserve all merged #490 fixture/discovery controls and do not reintroduce retry behavior. Add deterministic red controls proving: passing scenario + observed hunt remains passing; failing scenario/runner remains primary; applicable detector evidence is appended only as diagnostic context. Append a complete ideation cycle-2 Stage Report with exact proposed files and measurable evidence.
  AC-1/AC-4: scope is limited to two renamed diagnostic files and the Claude runner, AC-1 measures 1/1 silent pass and 2/2 exact primary failures against merged #490's 0/2 baseline, and all existing fixture/discovery controls remain required.

### Summary

Cycle 2 removes phase classification and gives detector evidence no independent authority over test status. A narrow failure-only cleanup keeps passing scenarios silent, preserves runner and scenario failures exactly, and appends one secondary diagnostic only after an existing red while leaving #490's fixture hardening and single-attempt behavior untouched.
