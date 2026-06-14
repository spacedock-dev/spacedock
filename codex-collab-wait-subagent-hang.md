---
id: czza18qnjzj75fznszxm3z0s
title: Codex rejection-flow collab:wait spawned-subagent stall (live-infra flakiness)
status: ideation
source: "j9 validation (2026-06-14) — the live Codex rejection-flow hung at collab:wait on 3/5 runs across j9's validation cycles (a spawned-subagent stall in the Codex runner). Orthogonal to j9 (the validator retried past it to a clean run); flagged by the validator for its own triage."
started: 2026-06-14T02:02:43Z
completed:
verdict:
score:
worktree:
issue:
sprint:
---

The Codex live `rejection-flow` shared-runtime scenario intermittently hangs at `collab:wait` — a spawned-subagent stall in the Codex runner, observed on 3 of 5 runs across j9's validation cycles. When it hangs the assertion is never reached (so it is NOT a scenario regression), but it makes the Codex live lane flaky and forces retries.

## Problem

The Codex runner spawns the reviewer subagent and waits on `collab:wait`; on roughly 60% of the flagged runs the wait stalls long enough that the rejection-flow assertions are never reached. A retry usually yields a clean run. This is live-infra flakiness, not evidence that the rejection-flow scenario is behaviorally wrong, but it degrades the Codex live lane and can waste CI time unless the wait is bounded by an oracle that fails or retries cleanly.

The fix must preserve the runtime-support proof policy: the runner may use Codex JSONL to detect host-specific producer events and wait boundaries, but a passing runtime claim still comes from live or fixture-backed durable state evidence: process exit, entity body, state git log, and clean state. Transcript phrasing or instruction-prose matches do not satisfy the claim.

## Failure boundary

Current evidence bounds the first implementation:

- Already covered: total stream silence. The Codex shared runner runs `codex exec --json` through the shared `streamWatcher.drainToExit(quietBudgetDefault, ...)`. If the Codex process stops emitting JSONL for more than the 60s no-progress budget, the existing watcher kills the subprocess and writes artifacts.
- Suspected gap: foreground mailbox wait progress. A `collab:wait` / `wait_agent` loop can keep the parent Codex process alive, or emit enough wait-loop JSONL to reset the stream-silence budget, while no worker final-status notification and no durable workflow-state change arrives. That is a runner liveness-classification gap, not a rejection-flow assertion failure.
- Possible but not the first boundary: a real worker stall where the spawned reviewer never completes, a Codex CLI mailbox-delivery bug where the worker completed but the parent never receives final status, or host resource contention that occurs only when the full serial live suite runs. The watchdog should report enough artifact detail to distinguish these after the first bounded failure.

Chosen bounded-stall strategy: add a Codex-specific collab-wait watchdog in the live runner, layered beside the existing 60s stream-silence watcher. Repeated wait events for the same worker handle are not progress by themselves. They reset the bounded wait only when paired with external evidence that can fail: a worker final-status/completion event, process exit, a durable entity-body/state-git change, or a non-wait scenario step that advances the runner toward an assertion. On watchdog trip, the runner kills `codex exec`, writes the JSONL/stderr/final-message artifacts, and returns a typed stall error. The rejection-flow scenario may retry at most once on that typed watchdog error, in a fresh workflow/artifact attempt; it must not retry assertion failures, auth/setup failures, or post-exit durable-state failures.

## Proposed approach

1. Keep the shared `streamWatcher` as the process-level stream-silence guard. Do not weaken `quietBudgetDefault`, add a long scenario timeout, or replace durable assertions with transcript-only checks.
2. Add a small Codex wait classifier for `collab_tool_call` items whose tool is the Codex foreground wait surface (`wait`, `wait_agent`, or the actual emitted spelling from the first fixture/live artifact). The classifier should key by receiver thread id / handle when present and track consecutive wait epochs without durable progress.
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

### AC-1: Codex collab waits are bounded by a falsifiable watchdog

The Codex shared runner detects repeated foreground wait activity for the same worker with no durable progress, kills the subprocess within the configured wait budget, and reports a typed watchdog error naming the scenario, wait handle/thread if available, and artifact directory.

Verified by:

- A deterministic unit/fixture test that feeds synthetic Codex JSONL wait events to the watchdog with a fake clock/proc and an unchanged entity file. It must assert the fake proc was killed, the error is the typed collab-wait stall, and the diagnostic names the scenario/handle.
- A positive control in the same test file where a durable entity change or final-status/completion event occurs before the budget and the watchdog does not trip.

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

The runtime live-CI docs distinguish the existing stream-silence guard from the new Codex foreground-wait watchdog and state that runtime success is still proven by durable workflow state.

Verified by:

- `go test ./internal/ensigncycle -run 'TestSharedScenarioDocsContract|TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions' -count=1` after the docs update.
- Full baseline gates from AGENTS before merge: `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

## Cheapest spot-check before live runs

Ideation did the cheap, non-live spot-check first:

```bash
go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse|TestFeedbackReflow' -count=1
```

Result: passed (`12 passed in 1 packages`). This confirms the existing deterministic rejection-flow seams still compile and that the Codex reviewer-reuse oracle already requires a real `send_input` to the bound validation reviewer thread. It does not prove the live wait is fixed; it only makes the implementation starting point explicit. The first implementation test should therefore be the deterministic collab-wait watchdog fixture, followed by the targeted live Codex `rejection-flow` run, and only then the full shared Codex suite.

## Test plan

1. Add the watchdog unit tests first. Cost: milliseconds, no auth. These are the cheapest reproduction of the failure class and should fail before implementation.
2. Add the retry-wrapper fixture tests. Cost: milliseconds, no auth. These prove the one-retry policy and prevent broad retry masking.
3. Run focused offline seams: `go test ./internal/ensigncycle -run 'TestCodexCollabWaitWatchdog|TestRunCodexRejectionFlowRetry|TestAssertCodexReviewerReuse|TestFeedbackReflow' -count=1`.
4. Run the no-secret shared scenario contract/docs guards: `go test ./internal/ensigncycle -run 'TestSharedScenarioDocsContract|TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions' -count=1`.
5. Run the baseline repo gates: `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`.
6. With Codex auth, run the targeted live scenario before paying for the full suite: `go test -tags live -count=1 -timeout 40m -run 'TestLiveCodexSharedScenarios/rejection-flow' ./internal/ensigncycle -v`. Inspect artifacts if it trips the watchdog.
7. Only after the targeted live run is clean or cleanly bounded, run `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`.

## Documentation impact

Implementation should update the Codex local-live paragraph in `docs/dev/README.md`. Proposed diff:

```diff
@@
-Run the Codex shared suite locally (`npm install -g @openai/codex` then `codex login`, or set `OPENAI_API_KEY`). Local runs may authenticate either through an existing Codex login at `~/.codex/auth.json` or through `OPENAI_API_KEY`. The test copies only `auth.json` into a temporary `CODEX_HOME` for the local subscription path; it does not copy local plugin state or the rest of the operator's Codex config. CI does not use local subscription auth.
+Run the Codex shared suite locally (`npm install -g @openai/codex` then `codex login`, or set `OPENAI_API_KEY`). Local runs may authenticate either through an existing Codex login at `~/.codex/auth.json` or through `OPENAI_API_KEY`. The test copies only `auth.json` into a temporary `CODEX_HOME` for the local subscription path; it does not copy local plugin state or the rest of the operator's Codex config. CI does not use local subscription auth. Codex shared runs use two liveness guards: the shared 60s stream-silence guard, and a Codex foreground-wait watchdog for repeated `collab:wait` / `wait_agent` activity with no durable workflow-state progress. The watchdog can fail or retry a typed stall, but scenario success still comes from process exit plus durable entity/git-state assertions.
```

## Stage Report: ideation

- DONE: Define the failure boundary: distinguish whether the `collab:wait` stall is likely Codex CLI, subagent/mailbox, runner quiet-budget, or resource contention, and record the chosen bounded-stall strategy.

- DONE: Rewrite the seed acceptance criterion into end-state `AC-N` properties with external proof that can fail: live Codex evidence or a deterministic watchdog/retry oracle.

- DONE: Produce an implementation-ready test plan that starts with the cheapest reproduction or spot-check before any expensive live Codex run.

### Summary

Fleshed out the ideation body for the Codex `collab:wait` live-runner flake. The design keeps the entity in ideation, chooses a Codex-specific wait watchdog plus one narrow retry, preserves durable-state runtime proof, and starts implementation with deterministic fixture tests before any targeted live Codex run.
