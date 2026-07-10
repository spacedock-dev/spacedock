---
title: Bound Claude live failure classification to the phase it actually observes
status: ideation
score: 0.7
source: Captain-approved recovery item 1 on 2026-07-10, after review of merged PR #490 showed classifyBootPreambleFailure scans the full transcript and can label later contract-lookup searches as boot-preamble failures before the scenario assertion is evaluated.
id: p4h6a5wcqe5ddkhnmrac1w9a
started: 2026-07-10T12:56:28Z
---

Separate true boot-orientation failures from later contract-lookup hunts, and preserve the scenario's primary failure evidence so the live harness reports the phase and cause it actually observed.

## Problem statement

PR #490 added `classifyBootPreambleFailure(stream, workflowRoot)` to the Claude shared-scenario runner. The classifier correctly detects wrong-root operations and broad filesystem hunts, but `claudeLiveRunner.run` passes the full transcript to it after the process exits or stalls. A matching search at any later point therefore becomes a "boot-preamble" failure and aborts the test before its scenario assertion runs.

The archived #490 cycle-2 evidence exposed the consequence: a run could boot successfully, perform scenario work, and then search for a deferred contract reference. The full-stream classifier discarded the phase distinction. Its later retry experiment even manufactured a false scenario failure by rerunning against state already mutated by the first attempt. The retry was reverted before merge, but the classifier still scans the full stream and still pre-empts the scenario's primary assertion.

This task fixes classification and reporting only. It does not permit contract hunts, retry a scenario, reset fixture state, or weaken the scenario assertions.

## Evidence probe

The risky question was whether the Claude stream exposes a stable boundary without reading model prose. It does. Saved sonnet streams contain a `Bash` `tool_use` for `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json --workflow-dir <fixture>` followed by a correlated, non-error `tool_result`. The keep-moving stream then loads `spacedock:fo-dispatch-core` and runs `find <plugin> -iname "claude-fo-dispatch.md"`; the feedback-escalation stream runs `find <fixture> ... -iname "README*"` after the same successful boot result. In both streams, the tool events order boot completion before the hunt without relying on assistant text.

Implementation tests will embed minimal, sanitized JSONL events from these shapes. They must not depend on the local `/tmp` artifacts.

## Approaches considered

1. **Recommended: split the transcript at a successful boot result.** Pair the `status --boot` tool-use ID with its non-error tool result. Run boot classifiers only on the prefix through that result; scan the suffix for a late-search diagnostic. This boundary describes an observed command completion, survives model wording changes, and treats a failed or missing boot result as unfinished orientation.
2. **Stop boot at the `status --boot` tool call.** This is simpler, but it moves concurrent tool uses and failed boot recovery into the later phase before orientation has succeeded. It would under-report genuine boot failures.
3. **Stop boot at the first scenario-specific mutation, dispatch, or skill load.** This recognizes progress, but every scenario would need its own marker and a contract hunt can occur before the first mutation. It couples the classifier to scenario details and leaves the disputed gap unresolved.

## Proposed approach

### Observable boundary

Parse the stream with the existing `streamEntry`, `toolUseBlocks`, `resultBlocks`, and `toolResultFailed` machinery. Track a `Bash` tool use that invokes `status --boot` for the expected fixture workflow. Boot orientation ends only when the stream carries a correlated, non-error tool result for that call.

The prefix includes the successful result. The suffix begins with the next stream event. A failed boot result, a missing result, or process death before the result leaves the whole transcript in the boot phase. Tool calls emitted in the same turn before the result also remain in the boot phase.

Run `detectWrongRootBoot` and the boot form of the broad-search detector on the prefix. Run a phase-neutral broad-search finding function on the suffix and format that finding as a **late contract-lookup diagnostic**, never as a boot-preamble failure. Keep detector input and output model-agnostic: tool name, command/path, tool-use ID, tool-result status, and event order only.

### Failure composition

An early boot finding remains fatal before scenario grading because the scenario never acquired a valid orientation. Preserve the existing wrong-root-before-broad-search priority and attach an underlying stall when present.

A late search still makes the test red, but it cannot pre-empt scenario grading. On a clean process exit, return the `liveResult`, register the late diagnostic for end-of-subtest reporting, and let the existing scenario assertion run first. If that assertion fails, its current message remains the primary failure and the late diagnostic follows as additional evidence. If the scenario assertion passes, the deferred late-search diagnostic becomes the sole failure.

When the runner cannot return a gradable result because of a stall or final-message extraction failure, that runner error is primary and the late-search diagnostic is appended. The formatter must never replace, relabel, or omit the primary error.

Expected developer-facing output changes from one misleading failure:

```text
FO broad-searched the filesystem at boot: ... a boot-preamble filesystem sweep, not the scenario's own assertion
```

to ordered evidence such as:

```text
<existing scenario assertion or runner failure>
Additional diagnostic: FO broad-searched the filesystem after successful boot orientation: ...
```

True prefix failures retain the existing `FO booted the wrong root` or `FO broad-searched the filesystem at boot` form.

### Implementation scope and #490 preservation

Change only the classifier and Claude-runner failure composition:

- `internal/ensigncycle/boot_preamble_classify_impl_test.go`: add the successful-result boundary and return separate early failure and late diagnostic.
- `internal/ensigncycle/boot_preamble_classify_test.go`: replace the ambiguous full-stream example with paired early/late controls built from correlated tool events.
- `internal/ensigncycle/broad_search_detect_impl_test.go`: expose a phase-neutral finding while retaining the existing boot wrapper and zero-discovery wording contract.
- `internal/ensigncycle/broad_search_detect_test.go`: keep its current detector matrix green and add formatting coverage only if the extraction requires it.
- `internal/ensigncycle/claude_live_runner_test.go`: classify by phase, preserve runner/scenario failure priority, and defer a suffix diagnostic until scenario grading finishes.
- `internal/ensigncycle/claude_live_failure_compose_test.go` (new): prove primary-first composition without launching a model.

Keep #490's fixture fixes intact: the absolute-root prompt signatures and commissioned markers in `shared_fixtures_test.go`; the corresponding root arguments in `claude_live_runner_test.go`, `codex_live_runner_test.go`, and `livescenario_adapter_live_test.go`; the gate-stop fixture relocation from `live_gate_stop_test.go`; and `boot_discovery_test.go`'s 10 discoverable-fixture plus markerless zero-discovery controls. Do not revive the reverted retry files or alter fixture content.

## Acceptance criteria

- **AC-1 (value — correct phase and primary evidence against the current classifier).** In a four-stream offline corpus containing two searches before boot completion and the same two searches after a correlated successful boot result, **2/2 early streams** classify as boot failures, **0/2 late streams** carry a boot-failure label, and **2/2 late streams** reach an injected scenario assertion before reporting their search diagnostic. The current full-stream classifier labels all four as boot failures and reaches **0/2** injected late-stream assertions, so either 1+ late stream mislabeled as boot or 1+ suppressed primary assertion fails this criterion. *Test:* `TestClassifyClaudeLiveFailurePhaseBoundary` and `TestClaudeLiveFailureCompositionKeepsPrimaryFirst` exercise sanitized JSONL from the keep-moving and feedback-escalation shapes.
- **AC-2 (boundary semantics).** Boot ends only at the correlated, non-error result of the expected fixture's `status --boot` call. A missing result, an error result, a result for another tool ID, and a search emitted before the successful result all remain in the boot prefix; a search after it enters the suffix. *Test:* table-driven `TestClaudeLiveBootBoundary` covers every case and asserts the exact prefix/suffix event counts.
- **AC-3 (failure priority).** Wrong-root remains ahead of broad-search within the boot prefix; an early boot failure remains primary over a stall; a stall or extraction error remains primary over a late diagnostic; and a scenario assertion remains primary when both it and a late diagnostic fail. A late diagnostic alone still reds the subtest. *Test:* `TestClassifyBootPreambleFailureWrongRootTakesPriority` plus table-driven `TestComposeClaudeLiveFailure` and `TestClaudeLiveFailureCompositionKeepsPrimaryFirst` assert both membership and output order.
- **AC-4 (#490 fixture non-regression).** All 10 intended shared fixtures remain discoverable, the zero-discovery fixture remains markerless and undiscoverable, all eight non-boot-subject prompts keep their absolute workflow root, and the reverted retry mechanism remains absent. *Test:* `TestSharedScenarioFixturesAreDiscoverable`, `TestZeroDiscoverFixtureStaysUndiscoverable`, the shared prompt-signature tests, and a diff review against merged PR #490's fixture files.

## Test plan

Implementation starts with the offline red controls. Use exact tool-use IDs and paired tool results; changing only the search's position around the successful boot result must flip early failure to late diagnostic. The composition tests use a fake reporter or pure formatter, cost milliseconds, and spend no model tokens.

Run focused classifier and composition tests:

```bash
go test ./internal/ensigncycle -run 'Test(ClaudeLiveBootBoundary|ClassifyClaudeLiveFailurePhaseBoundary|ClassifyBootPreambleFailure|ComposeClaudeLiveFailure|ClaudeLiveFailureCompositionKeepsPrimaryFirst)'
```

Run #490's fixture and detector controls:

```bash
go test ./internal/ensigncycle -run 'Test(SharedScenarioFixturesAreDiscoverable|ZeroDiscoverFixtureStaysUndiscoverable|DetectBroadSearchAtBoot|DetectWrongRootBoot)'
```

Then run the repository gates and live-tag compile check:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
go test -tags live -run '^$' ./internal/ensigncycle/...
```

No new live workflow run is required: the claim concerns deterministic post-hoc classification of an already captured stream, not model behavior. The saved live shapes established that the command/result boundary exists; the sanitized offline fixtures provide the durable proof. No product or docs-site change is needed because this alters only developer-facing live-test failure output; the exact before/after output contract appears above.

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
