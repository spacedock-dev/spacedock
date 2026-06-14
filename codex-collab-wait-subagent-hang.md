---
id: czza18qnjzj75fznszxm3z0s
title: Codex rejection-flow collab:wait spawned-subagent stall (live-infra flakiness)
status: validation
source: "j9 validation (2026-06-14) — the live Codex rejection-flow hung at collab:wait on 3/5 runs across j9's validation cycles (a spawned-subagent stall in the Codex runner). Orthogonal to j9 (the validator retried past it to a clean run); flagged by the validator for its own triage."
started: 2026-06-14T02:02:43Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-codex-collab-wait-subagent-hang
issue:
sprint:
---

The Codex live `rejection-flow` shared-runtime scenario intermittently hangs at `collab:wait` — a spawned-subagent stall in the Codex runner, observed on 3 of 5 runs across j9's validation cycles. When it hangs the assertion is never reached (so it is NOT a scenario regression), but it makes the Codex live lane flaky and forces retries.

## Problem

The Codex runner spawns the reviewer subagent and waits on `collab:wait`; on roughly 60% of the flagged runs the wait stalls long enough that the rejection-flow assertions are never reached. A retry usually yields a clean run. This is live-infra flakiness, not evidence that the rejection-flow scenario is behaviorally wrong, but it degrades the Codex live lane and can waste CI time unless the wait is bounded by an oracle that fails or retries cleanly.

The fix must preserve the runtime-support proof policy: the runner may use Codex JSONL to detect host-specific producer events and wait boundaries, but a passing runtime claim still comes from live or fixture-backed durable state evidence: process exit, entity body, state git log, and clean state. Transcript phrasing or instruction-prose matches do not satisfy the claim.

## Failure boundary

Current evidence bounds the first implementation:

- Observed j9 shape: the validation record for `lazy-teamcreate-shallow-boot` says Codex `rejection-flow` runs hung at `collab:wait` and the no-progress 60s budget killed them. That means the typed stall must cover a silent/no-progress interval after a foreground wait, not only a chatty stream of repeated wait events.
- Already present but too generic: the Codex shared runner runs `codex exec --json` through `streamWatcher.drainToExit(quietBudgetDefault, ...)`. If the Codex process stops emitting JSONL for more than the 60s no-progress budget, the existing watcher kills the subprocess and writes artifacts, but it reports a generic stream stall rather than a Codex foreground-wait stall when the active or last meaningful event was `collab:wait`.
- Second possible shape: foreground mailbox wait progress. A `collab:wait` / `wait_agent` loop can keep the parent Codex process alive, or emit enough wait-loop JSONL to reset the stream-silence budget, while no worker final-status notification and no durable workflow-state change arrives. That is also a runner liveness-classification gap, not a rejection-flow assertion failure.
- Possible but not the first boundary: a real worker stall where the spawned reviewer never completes, a Codex CLI mailbox-delivery bug where the worker completed but the parent never receives final status, or host resource contention that occurs only when the full serial live suite runs. The watchdog should report enough artifact detail to distinguish these after the first bounded failure.

Chosen bounded-stall strategy: add a Codex-specific foreground-wait stall classifier in the live runner, layered beside the existing 60s stream watcher. It has two arms: repeated wait events for the same worker with no durable progress, and stream silence/no-progress while a foreground wait is active or when the last meaningful event was a wait. Repeated wait events are not progress by themselves. A wait is cleared only by external evidence that can fail: a worker final-status/completion event, process exit, a durable entity-body/state-git change, or a non-wait scenario step that advances the runner toward an assertion. On either watchdog arm, the runner kills or observes the killed `codex exec`, writes the JSONL/stderr/final-message artifacts, and returns the same typed collab-wait stall error. The rejection-flow scenario may retry at most once on that typed watchdog error, in a fresh workflow/artifact attempt; it must not retry assertion failures, auth/setup failures, or post-exit durable-state failures.

## Proposed approach

1. Keep the shared `streamWatcher` as the process-level stream-silence guard. Do not weaken `quietBudgetDefault`, add a long scenario timeout, or replace durable assertions with transcript-only checks. Extend the Codex runner's handling of `streamWatcher` no-progress errors so a timeout whose active/last meaningful Codex collab event was a foreground wait is classified as the typed collab-wait stall, not as an untyped generic hang.
2. Add a small Codex wait classifier for `collab_tool_call` items whose tool is the Codex foreground wait surface (`wait`, `wait_agent`, `collab:wait`, or the actual emitted spelling from the first fixture/live artifact). The classifier should key by receiver thread id / handle when present and track both repeated wait epochs and the active-last-wait state used by the silence arm.
3. Add a durable-progress probe to the Codex scenario runner. For `rejection-flow`, the cheapest useful probe is the entity file hash plus state/git observation where available; if the fixture is single-root, the entity body itself is enough for the fixture-backed tests. The probe is used only to decide whether a wait loop is making externally visible progress, not to pass the scenario.
4. Refactor `codexLiveRunner.run` so watchdog stalls can be returned as typed errors instead of immediate `t.Fatalf`, then let the scenario wrapper perform one retry only for that error class. Artifacts should be written under attempt-specific directories so attempt 1 remains inspectable when attempt 2 passes.
5. Keep the final pass condition exactly where it belongs: `assertRejectionFlow` over durable entity state plus `assertCodexReviewerReuse` over the host-specific producer signal. The watchdog only prevents indefinite waiting; it does not make the scenario pass.

## Out of scope

- Changing the j9 deliverable. The validator retried past the hang and j9's Codex AC-3 passed on a clean run.
- Changing Codex CLI behavior, mailbox semantics, or the Codex FO runtime contract beyond the runner's bounded handling of foreground waits.
- Retrying every live failure. Only the typed collab-wait watchdog stall is retryable, and only once.
- Adding PR/mod behavior or changing shared scenario assertions.
- Treating transcript text, instruction prose, or a static grep over skills/docs as proof that runtime behavior works.

## Acceptance criteria

### AC-1: Codex foreground waits are bounded by a falsifiable typed stall

The Codex shared runner reports the same typed collab-wait stall when either foreground-wait shape occurs with no durable progress: repeated wait events for the same worker, or stream silence/no-progress while a foreground wait is active or when the last meaningful event was a wait. The error names the scenario, wait handle/thread if available, the stall arm, and artifact directory.

Verified by:

- Fixture A: synthetic Codex JSONL emits repeated foreground wait events for one receiver thread while the entity hash/state probe stays unchanged. The test must assert the fake proc is killed, the error is the typed collab-wait stall, and the diagnostic names the scenario, handle/thread, and repeated-wait arm.
- Fixture B: synthetic Codex JSONL emits a foreground wait event, then no further stream progress until the quiet budget expires. The test must assert the same typed collab-wait stall, not a generic `streamWatcher` timeout, and the diagnostic names the scenario, handle/thread if known, and silent-after-wait arm. This fixture covers the j9-recorded no-progress 60s kill shape.
- Positive controls in the same test file where a durable entity change, a worker final-status/completion event, or a non-wait scenario-progress event occurs before the budget and the watchdog does not trip.

### AC-2: A bounded retry is narrow and cannot hide scenario failures

The `rejection-flow` Codex runner retries at most once after the typed collab-wait watchdog error, writes attempt-specific artifacts, and never retries assertion failures, auth/setup failures, process start failures, or durable-state failures after process exit.

Verified by:

- A fixture-backed runner test with a fake Codex command/proc sequence: first attempt returns the typed watchdog stall, second attempt exits cleanly, and the wrapper records exactly two attempts.
- Negative controls where a non-watchdog error and a rejection-flow assertion failure each stop after one attempt.

### AC-3: Rejection-flow success remains proven by durable state plus the real Codex reuse signal

A passing Codex `rejection-flow` still requires the two-cycle durable entity outcome and a real `send_input` to the cycle-1 validation reviewer thread. The watchdog cannot turn a transcript-only wait/reuse story into a pass.

Verified by:

- Existing offline producer-signal tests continue to pass: `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse|TestFeedbackReflow' -count=1`.
- A new or extended negative fixture where the JSONL contains wait/reuse-looking events but the entity lacks the second-cycle durable state must fail the rejection-flow assertion.
- A targeted live run when Codex auth is available: `go test -tags live -count=1 -timeout 40m -run 'TestLiveCodexSharedScenarios/rejection-flow' ./internal/ensigncycle -v`. A valid failure is the typed watchdog stall with artifacts, not a Go timeout or silent hang; a valid pass must still satisfy the durable assertions.

### AC-4: Developer docs describe both liveness guards accurately

The runtime live-CI docs distinguish the existing stream-silence guard from the new Codex foreground-wait watchdog and state that runtime success is still proven by durable workflow state. The docs proof must be load-bearing: implementation adds or extends a docs-contract test that fails unless `docs/dev/README.md` contains the concrete watchdog clauses.

Verified by:

- A new focused docs test, or an explicit extension of `TestSharedScenarioDocsContract`, requires all four clauses: `Codex foreground-wait watchdog`, `collab:wait` / `wait_agent`, durable workflow-state progress, and typed stall behavior.
- Default-tag docs/definition guard: `go test ./internal/ensigncycle -run 'TestSharedScenarioDocsContract|TestCodexForegroundWaitWatchdogDocsContract|TestSharedRuntimeScenarioDefinitions' -count=1` after the docs update, using the actual focused docs test name if it differs.
- Live-tag no-secret parity guard, matching `docs/dev/README.md`: `go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPiSharedScenarioCoverage' ./internal/ensigncycle -v`.
- Full baseline gates from AGENTS before merge: `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

## Cheapest spot-check before live runs

Ideation did the cheap, non-live spot-check first:

```bash
go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse|TestFeedbackReflow' -count=1
```

Result: passed (`12 passed in 1 packages`). This confirms the existing deterministic rejection-flow seams still compile and that the Codex reviewer-reuse oracle already requires a real `send_input` to the bound validation reviewer thread. It does not prove the live wait is fixed; it only makes the implementation starting point explicit. The first implementation test should therefore be the deterministic collab-wait watchdog fixture, followed by the targeted live Codex `rejection-flow` run, and only then the full shared Codex suite.

## Test plan

1. Add the watchdog unit tests first. Cost: milliseconds, no auth. These are the cheapest reproduction of the failure class and should fail before implementation. Include both stall fixtures: repeated foreground wait events with unchanged durable state, and silent/no-progress after a foreground wait.
2. Add the retry-wrapper fixture tests. Cost: milliseconds, no auth. These prove the one-retry policy and prevent broad retry masking.
3. Run focused offline seams: `go test ./internal/ensigncycle -run 'TestCodexCollabWaitWatchdog|TestRunCodexRejectionFlowRetry|TestAssertCodexReviewerReuse|TestFeedbackReflow' -count=1`.
4. Run the default-tag docs and definition guards: `go test ./internal/ensigncycle -run 'TestSharedScenarioDocsContract|TestCodexForegroundWaitWatchdogDocsContract|TestSharedRuntimeScenarioDefinitions' -count=1` after adding the focused docs assertion.
5. Run the live-tag no-secret parity guards: `go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPiSharedScenarioCoverage' ./internal/ensigncycle -v`.
6. Run the baseline repo gates: `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`.
7. With Codex auth, run the targeted live scenario before paying for the full suite: `go test -tags live -count=1 -timeout 40m -run 'TestLiveCodexSharedScenarios/rejection-flow' ./internal/ensigncycle -v`. Inspect artifacts if it trips the watchdog.
8. Only after the targeted live run is clean or cleanly bounded, run `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`.

## Documentation impact

Implementation should update the Codex local-live paragraph in `docs/dev/README.md`. Proposed diff:

```diff
@@
-Run the Codex shared suite locally (`npm install -g @openai/codex` then `codex login`, or set `OPENAI_API_KEY`). Local runs may authenticate either through an existing Codex login at `~/.codex/auth.json` or through `OPENAI_API_KEY`. The test copies only `auth.json` into a temporary `CODEX_HOME` for the local subscription path; it does not copy local plugin state or the rest of the operator's Codex config. CI does not use local subscription auth.
+Run the Codex shared suite locally (`npm install -g @openai/codex` then `codex login`, or set `OPENAI_API_KEY`). Local runs may authenticate either through an existing Codex login at `~/.codex/auth.json` or through `OPENAI_API_KEY`. The test copies only `auth.json` into a temporary `CODEX_HOME` for the local subscription path; it does not copy local plugin state or the rest of the operator's Codex config. CI does not use local subscription auth.
+
+Codex shared runs use two liveness guards: the shared 60s stream-silence guard, and a Codex foreground-wait watchdog for repeated `collab:wait` / `wait_agent` activity with no durable workflow-state progress or for no-progress silence after a foreground wait. The watchdog can fail or retry a typed stall, but scenario success still comes from process exit plus durable entity/git-state assertions.
```

## Stage Report: ideation

- DONE: Define the failure boundary: distinguish whether the `collab:wait` stall is likely Codex CLI, subagent/mailbox, runner quiet-budget, or resource contention, and record the chosen bounded-stall strategy.

- DONE: Rewrite the seed acceptance criterion into end-state `AC-N` properties with external proof that can fail: live Codex evidence or a deterministic watchdog/retry oracle.

- DONE: Produce an implementation-ready test plan that starts with the cheapest reproduction or spot-check before any expensive live Codex run.

### Summary

Fleshed out the ideation body for the Codex `collab:wait` live-runner flake. The design keeps the entity in ideation, chooses a Codex-specific wait watchdog plus one narrow retry, preserves durable-state runtime proof, and starts implementation with deterministic fixture tests before any targeted live Codex run.

## Stage Report: ideation (review repair)

- DONE: Repaired the failure boundary and AC-1 to cover both typed stall arms: repeated foreground wait events with no durable progress and stream silence/no-progress while a foreground wait is active or was the last meaningful event. Added explicit fixture expectations for both, including the j9-recorded no-progress-after-`collab:wait` shape.

- DONE: Corrected no-secret guard commands so live-tagged parity tests run with `-tags live`, split from default-tag docs/definition checks.

- DONE: Made AC-4 load-bearing by requiring a docs-contract assertion for the concrete Codex foreground-wait watchdog wording, `collab:wait` / `wait_agent`, durable workflow-state progress, and typed stall behavior; split the proposed README diff into a short separate paragraph.

### Summary

Review repairs applied to the ideation body only. The design now covers both observed stall shapes, uses executable guard commands that match build tags, and requires a focused docs-contract test for the new watchdog wording.

## Stage Report: implementation

- DONE: Implement a typed Codex foreground-wait stall path covering both repeated wait events and silent/no-progress after a wait, with deterministic fixtures and positive controls.
  Code commit `5431f64c` adds the Codex foreground-wait watchdog, repeated-wait and silent-after-wait fixtures, durable-progress/completion/non-wait positive controls, and typed diagnostics with scenario, handle, arm, durable-progress flag, and artifact directory.
- DONE: Keep retry behavior narrow: at most one retry only for the typed stall, with attempt-specific artifacts and negative controls for non-watchdog/assertion failures.
  `TestRunCodexRejectionFlowRetry*` passed in the focused suite (`23 passed`), and the targeted live run wrote attempt-specific artifacts under `/tmp/spacedock-codex-live-rejection-flow-after-durable/codex-shared-scenarios/rejection-flow/{attempt-1,attempt-2}`.
- DONE: Preserve durable-state proof and docs coverage: rejection-flow success still depends on durable entity state plus real Codex reuse signal, and docs-contract tests guard the new watchdog wording.
  `assertRejectionFlow` plus `assertCodexReviewerReuse` still gate success; docs guard passed (`3 passed`), live parity passed (`3 passed`), `go test ./...` and `go test ./... -race` each passed (`1260 passed in 16 packages`), and targeted live `TestLiveCodexSharedScenarios/rejection-flow` passed on attempt 2 with `- Cycle 2: PASSED`.

### Summary

Implemented the Codex-specific foreground-wait watchdog beside the shared stream-silence guard, plus a one-retry wrapper for the `rejection-flow` live scenario only. The fix keeps runtime success tied to durable entity/git state and the real Codex `send_input` reuse signal; the live retry artifacts show attempt 1 bounded at an active foreground wait and attempt 2 completing the two-cycle rejection flow.

## Stage Report: validation

- FAILED: Reproduce AC-1 and AC-2 with deterministic tests: both foreground-wait stall arms, positive controls, one retry only for typed stall, and no retry for non-watchdog/assertion failures.
  Focused suite passed (`go test ./internal/ensigncycle -run 'TestCodexCollabWaitWatchdog|TestRunCodexRejectionFlowRetry|TestAssertCodexReviewerReuse|TestFeedbackReflow' -count=1 -v`: 23 passed), but AC-1 is not satisfied: `internal/ensigncycle/codex_collab_wait_watchdog_test.go:105` expects a typed stall after durable progress before the budget, and `internal/ensigncycle/codex_collab_wait_watchdog_impl_test.go:111` returns a typed stall when `probeChanged()` is true instead of clearing the wait.
- DONE: Verify AC-3 remains durable-state based: rejection-flow success still requires durable two-cycle state plus the real Codex `send_input` reuse signal, including the negative fixture.
  Focused suite passed (`23 passed`), including reviewer-reuse tests and `TestRejectionFlowRejectsWaitReuseTranscriptWithoutDurableSecondCycle`.
- DONE: Verify AC-4 and baseline gates: docs-contract wording guard, live-tag no-spend parity guard, `gofmt`, `go test ./...`, `go test ./... -race`.
  Docs/definition guard passed (`3 passed`), live-tag no-spend parity guard passed (`3 passed`), `gofmt -w ./cmd ./internal` left the code worktree clean, `go test ./...` passed (`1260 passed in 16 packages`), and `go test ./... -race` passed (`1260 passed in 16 packages`).
- SKIPPED: targeted live Codex rejection-flow when auth is available.
  Not rerun during validation because the deterministic AC-1 contradiction is already material and should be repaired before spending on another live Codex run.

### Summary

Recommendation: REJECTED. AC-1 requires durable entity progress before the budget to be a positive control that clears the foreground-wait watchdog, but the implementation's own test and `silenceStall()` path encode the opposite behavior. Expected correction: make durable progress clear the active wait/no typed stall for the silence-after-wait path, and replace the contradictory test with a no-stall positive control for durable progress before the budget.

### Feedback Cycles

- Cycle 1 (2026-06-14T04:13:02Z): validation rejected AC-1. Route back to implementation to make durable progress before the silence budget clear the active foreground wait instead of returning a typed stall, and replace the contradictory durable-progress test with a no-stall positive control.

## Stage Report: implementation (cycle 2)

- DONE: Fix AC-1 durable-progress semantics for the silence-after-wait path: durable progress before budget must clear the active wait and avoid a typed stall.
  Code commit `dfdbf353` makes durable workflow-state progress clear the Codex active wait and reset the shared quiet deadline instead of returning a typed silence-after-wait stall.
- DONE: Replace or rewrite the contradictory deterministic test so durable progress before the budget is a positive no-stall control.
  `TestCodexCollabWaitWatchdogDurableProgressBeforeBudgetClearsSilentWait` now fails red on the old behavior and passes after the fix with no fake proc kill.
- DONE: Re-run the focused watchdog/retry/reuse suite and any affected baseline checks, then write a new implementation stage report with DONE/SKIPPED/FAILED for these items.
  `gofmt -w ./cmd ./internal`, focused watchdog/retry/reuse suite (`23 passed`), `go test ./...` (`1260 passed in 16 packages`), and `go test ./... -race` (`1260 passed in 16 packages`) all passed.

### Summary

Repaired the feedback-cycle AC-1 issue by treating durable state changes during a silent Codex foreground wait as liveness evidence that clears the wait and extends the quiet budget. The deterministic durable-progress test is now a no-stall positive control, while the typed stall path remains for repeated waits and silent waits with no durable progress.
