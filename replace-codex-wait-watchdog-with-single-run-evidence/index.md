---
title: Replace Codex wait watchdog with single-run evidence
status: validation
source: "Test-infrastructure audit 2026-07-14."
started: 2026-07-15T06:16:10Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-replace-codex-wait-watchdog-with-single-run-evidence
issue:
milestone: 0.26.0
id: 15sdwn85ekhwkf0jxnh38h8j
mod-block: merge:pr-merge
pr: pr-merge:515
---

## Problem

The Codex rejection-flow harness owns a parallel wait protocol: it parses runtime events, tracks wait epochs, fingerprints Markdown and Git state, kills the process, and retries the entire journey once. Archived evidence includes a green result on attempt 2, so a failed first run can disappear behind a false green.

There are three different bounds in the current lane, and they must not be conflated:

- the suite-wide `go test -timeout 40m` backstop bounds the whole test binary;
- the shared 60-second stream-silence watchdog resets whenever stdout emits another line;
- the Codex foreground-wait watchdog parses `wait`, `wait_agent`, and `collab:wait`, tracks a wait epoch against Markdown/Git fingerprints, and may authorize a second rejection-flow journey.

The latter two make the test a second runtime controller. A chatty runaway can extend the silence deadline, while a wait event that the harness classifies as a typed stall can discard the first fixture and start again. Neither is needed to decide whether one supported `codex exec` produced the required durable workflow outcome.

## Required outcome

Make one supported Codex run authoritative. Bound the real process, preserve its artifacts, and grade its exit plus durable workflow state. Do not recreate Codex wait semantics or add another controller, recovery protocol, daemon, lease, or lifecycle layer.

## Scope and non-goals

This task changes only the Codex shared-scenario process boundary, the rejection-flow runner, their focused tests, and the Codex paragraph in `docs/runtime-live-ci.md`. Claude's stream watcher, Codex's producer-side first-officer wait contract, the host-neutral rejection-flow assertions, reviewer-reuse grading, prompts, fixtures, and scenario coverage stay intact. This task does not make runtime events a new success oracle and does not tune, replace, or model `wait_agent` semantics.

## Proposed approach

1. Give every Codex shared-scenario launch a fixed 15-minute wall-clock context. The deadline starts once, immediately before `codex exec`, and never resets for JSONL activity, wait events, or durable writes. `codex exec` is invoked with `exec.CommandContext` exactly once; the suite-wide 40-minute Go timeout remains only an outer emergency backstop.
2. Connect stdout and stderr directly to `codex-exec.jsonl` and `codex-exec.stderr.txt`. Use `cmd.Run()` as the supported process boundary, close the files after it returns, and record `started`, `exit_code`, `timed_out`, and duration in `codex-process-result.txt`. A timeout is distinct from an ordinary nonzero exit; both fail the scenario.
3. Flatten the artifact layout to `<artifact-root>/<scenario>/`; remove the attempt number from `codexLiveRunner.run`, delete `runCodexRejectionFlowWithRetry`, and call the runner once from `runCodexRejectionFlowScenario`. There is no `attempt-1` or `attempt-2` directory because there is no scenario retry concept.
4. After that single process ends, read the preserved JSONL for the existing producer assertions and read the workflow for the existing durable assertions. For rejection-flow diagnostics, preserve `rejection-task.after.md`, `git-head.txt`, `git-log.txt`, and `git-status.txt` before reporting a process or assertion failure. This is post-run evidence capture, not an in-run progress probe.
5. Preserve the current local-checkout plugin proof and add `_setup/source-head.txt`, populated from the checkout used to build/install the plugin. An exact-head live run compares this file with the pre-run `git rev-parse HEAD`, so a green artifact is bound to the implementation commit rather than merely to a local path.

The success decision is therefore linear:

`one codex exec` -> `fixed deadline or process exit` -> `artifacts closed and recorded` -> `exit code must be 0` -> `durable rejection and reviewer assertions must pass`.

## Spike: bounded single-run evidence

The riskiest mechanism was exercised during ideation with a throwaway Go test that used the proposed `context.WithTimeout` + `exec.CommandContext` + direct artifact-file shape. A fake Codex child incremented an invocation file, emitted one JSONL event and one stderr marker, appended a validation failure to an entity, committed that entity, and then either exited 23 or stalled past the hard deadline. The throwaway source was removed after the run; it seeds the implementation test.

Command: `go test ./.codex-single-run-spike -run TestBoundedSingleRunPreservesEvidence -count=1 -v`.

- nonzero case: `invocations=1 exit=23 deadline=false jsonl=present stderr=present entity=present git_head="fault: durable first-run state"`;
- stall case: `invocations=1 exit=-1 deadline=true jsonl=present stderr=present entity=present git_head="fault: durable first-run state"`;
- result: 2/2 passed in 2.638s on 2026-07-15.

This proves the process/evidence mechanism without parsing a runtime event or launching a second attempt. It does not substitute for the exact-head live producer proof in AC-4.

## Acceptance criteria

- **AC-1 (VALUE — first run is authoritative):** Each Codex rejection-flow test invocation starts exactly one `codex exec`; a nonzero exit or 15-minute expiry makes that invocation red, and no later run can replace its verdict.
  - **Test:** a helper-process test records its invocation count, exits 23 in one case and stalls in another, and asserts `count == 1`, the expected exit/timeout classification, and a failing scenario result. A mutation that calls the helper a second time or accepts its failure must make the test red.
- **AC-2 (supported boundary, no wait controller):** Codex scenario liveness is decided only by the fixed wall-clock process deadline; JSONL volume, `wait`/`wait_agent`/`collab:wait` events, and intermediate Markdown/Git changes neither extend nor shorten it. Scenario success requires process exit 0 followed by the existing durable-state and producer assertions.
  - **Test:** the timeout helper emits repeated wait-shaped JSONL and commits durable state but still expires at the original short test deadline. Focused tests and a code review confirm the old watchdog, progress probe/fingerprint, stream-watcher use in the Codex adapter, typed-stall classification, and retry entry points are absent; existing rejection-flow and reviewer-reuse assertion tests remain green.
- **AC-3 (diagnosable first-run failure):** After an ordinary nonzero exit or hard timeout, the sole scenario artifact directory retains JSONL, stderr, the process result with exact exit/timeout status, and rejection-flow entity/Git evidence from that same run.
  - **Test:** both helper fault cases assert nonempty `codex-exec.jsonl`, the exact stderr marker, `codex-process-result.txt`, the appended entity report, the committed Git head/log, and no `attempt-*` directories. Temporarily disabling each production artifact write/capture must make the focused test red.
- **AC-4 (exact-head live producer proof):** One real rejection-flow run against the implementation HEAD exits 0 on its only `codex exec`, records cycle 1 REJECTED and cycle 2 PASSED in durable state, satisfies the existing separate-worker/reviewer assertion, and binds its artifacts to that source HEAD.
  - **Test:** from a committed implementation with clean staged and unstaged diffs, run `SPACEDOCK_LIVE_ARTIFACT_DIR=<fresh-dir> go test -tags live -count=1 -timeout 20m -run '^TestLiveCodexSharedScenarios/rejection-flow$' ./internal/ensigncycle -v`; assert `_setup/source-head.txt` equals the pre-run `git rev-parse HEAD`, the flat `rejection-flow/` process result says exit 0 and not timed out, the preserved entity/Git artifacts satisfy the existing graders, and no `attempt-*` path exists.
- **AC-5 (VALUE — simpler harness):** Relative to audit baseline `9aadd89e5915f57587f5eea93d5b76982d8de9f1`, the Codex single-run change removes at least 400 net Go lines across `codex_collab_wait_watchdog_impl_test.go`, `codex_collab_wait_watchdog_test.go`, `codex_live_runner_test.go`, the new `codex_single_run_test.go`, and `shared_reviewer_reuse_table_test.go`; none of the removed controller/retry symbols remain, and no replacement poller, watcher, epoch, lease, daemon, or recovery loop is introduced.
  - **Test:** run `git diff --numstat 9aadd89e5915f57587f5eea93d5b76982d8de9f1..HEAD -- internal/ensigncycle/{codex_collab_wait_watchdog_impl_test.go,codex_collab_wait_watchdog_test.go,codex_live_runner_test.go,codex_single_run_test.go,shared_reviewer_reuse_table_test.go}` and require summed `deletions - additions >= 400`; run `rg` for `codexCollabWaitWatchdog|workflowStateFingerprint|runCodexRejectionFlowWithRetry|codexCollabWaitStallError|attempt-2` with no implementation hits, then review the new process helper as one `CommandContext.Run` path. Baseline sizes are 355 lines in `codex_collab_wait_watchdog_impl_test.go`, 377 in `codex_collab_wait_watchdog_test.go`, and 492 in `codex_live_runner_test.go`.

## Test plan

- **Offline process-boundary test (low cost, <5s):** promote the ideation spike into a focused Go helper-process test. Cover exit 23 and a hard stall, exact one-invocation count, direct JSONL/stderr retention, process-result fields, and post-process entity/Git evidence. This is the fault-injected proof for AC-1 through AC-3.
- **Offline regression tests (low cost):** keep the host-neutral rejection-flow and Codex reviewer-reuse tables, run `go test ./internal/ensigncycle -count=1`, then `go test ./...` and `go test ./... -race`. Apply two temporary mutations: retry after the exit-23 fault, and ignore the exit code; each must red the focused test.
- **Quantitative audit (low cost):** measure the named files against baseline `9aadd89e…` and record the net delta and removed-symbol sweep in the implementation report. This proves AC-5 can move the wrong way if replacement machinery grows.
- **Exact-head Codex run (live, medium cost, expected 9-15 minutes):** run only `rejection-flow` once with a fresh persistent artifact root, verify source HEAD, exit/result artifacts, flat layout, entity cycle state, Git log/status, and existing producer/reviewer assertions. Do not rerun a red invocation as evidence; diagnose that first run from its artifacts and fix the implementation or runtime contract before a new validation run.

No golden fixture is needed: the process-boundary fault is deterministic Go behavior, while the actual producer claim requires the one live Codex run. The full live Codex suite remains a later regression gate, not a substitute for the isolated exact-head proof.

## Documentation diff

Implementation changes the Codex paragraph in `docs/runtime-live-ci.md` from:

> Codex shared runs use two liveness guards: the shared 60s stream-silence guard, and a Codex foreground-wait watchdog for repeated `collab:wait` / `wait_agent` activity with no durable workflow-state progress or for no-progress silence after a foreground wait. The watchdog can fail or retry a typed stall, but scenario success still comes from process exit plus durable entity/git-state assertions.

to:

> Each Codex shared scenario launches one `codex exec`. A fixed 15-minute wall-clock process limit is its only scenario-level liveness guard; JSONL activity, `wait_agent` events, and durable writes do not extend the deadline, and the runner does not retry. The runner preserves JSONL, stderr, the process result, and post-run durable entity/Git evidence, then requires exit 0 and grades the existing workflow assertions. The suite-wide `-timeout 40m` remains a loose outer backstop.

## Mechanism/value trace

The fixed `CommandContext` limit serves AC-1/AC-2. The simpler alternative, relying only on `go test -timeout`, can terminate the entire test binary before one scenario closes its diagnostic files; retaining the resettable 60-second silence guard would preserve a test-owned runtime protocol and would not stop chatty runaway output.

Direct stdout/stderr files plus a small process-result record serve AC-3. Reusing the pipe/scanner is insufficient because it exists to feed the in-run controller being removed; reading JSONL once after exit is enough for the existing producer assertions.

The post-run entity/Git snapshot serves AC-3/AC-4. Leaving state only under `t.TempDir()` is insufficient because Go cleanup removes the evidence that explains a CI red. The snapshot never participates in liveness or retry.

The source-head artifact serves AC-4. The existing local-plugin path check proves checkout locality but does not bind a particular live artifact set to a commit SHA.

The line-delta and removed-symbol audit serves AC-5. A symbol-absence check alone could pass after replacing the same controller under new names; the independent net-line reduction can move the wrong way. If the single boundary plus existing graders cannot prove the value, stop for architecture review instead of extending the harness.

## Stage Report: ideation

- DONE: Define a concrete minimal single-run design that removes whole-scenario retry and Codex wait-state parsing while preserving one hard process bound and diagnostic artifacts.
  The proposed path is one 15-minute `CommandContext.Run`, direct artifacts, flat per-scenario layout, post-run grading, and no retry or in-run event/state controller.
- DONE: Exercise and record the riskiest path: a bounded Codex process preserves its exit, JSONL, stderr, and durable state without a test-owned runtime protocol or second attempt.
  Throwaway fault injection passed 2/2 in 2.638s: one exit-23 run and one hard-timeout run each preserved JSONL, stderr, entity, Git head, and invocation count 1.
- DONE: Pair the revised acceptance criteria with fault-injected, exact-head live, and quantitative code-reduction evidence sufficient for implementation and validation.
  AC-1 through AC-5 bind the helper fault matrix, exact-head live command/artifacts, and >=400 net-line reduction against `9aadd89e5915f57587f5eea93d5b76982d8de9f1`.

### Summary

Ideation replaces the Codex harness's silence/wait/retry stack with a single fixed process boundary and post-run evidence grading. The risky OS-process path is exercised, the real Codex proof is pinned to one exact-head run, and the design carries an independently measurable simplification target plus a concrete documentation diff.

## Stage Report: implementation

- DONE: Replace the retrying Codex wait controller with one fixed-deadline CommandContext.Run path; add no watcher, poller, retry, or recovery loop.
  Commit `d199fc98630408649e89f39a5db57a76a9049732` uses one 15-minute `CommandContext` and one `cmd.Run`; removed-symbol and controller review sweeps are clean.
- DONE: Preserve same-run process, JSONL, stderr, entity, Git, and source-HEAD evidence; make focused fault tests and existing rejection/reviewer assertions green.
  Exit-23 and hard-timeout faults preserve all artifacts with one invocation; retry and ignored-exit mutations red, and the exact-head live rejection flow passed once in 402.64s under `/tmp/spacedock-codex-single-run-d199fc98.yWE8r1`.
- DONE: Remove at least 400 net Go lines against the named baseline, apply the runtime-live documentation diff, and run the required focused, full, race, and formatting checks.
  The named-file audit reports 335 additions, 820 deletions, and 485 net lines removed; `go test ./...`, `go test ./... -race`, focused/live-tag tests, `gofmt`, and `git diff --check` pass.

### Summary

The Codex shared runner now launches one authoritative process under a fixed deadline, writes process streams directly to flat scenario artifacts, and grades durable state only after exit. The rejection-flow lane retains post-run entity/Git evidence and exact source-HEAD provenance while deleting the wait watchdog, state fingerprinting, typed-stall handling, and whole-scenario retry.

## Stage Report: validation

- FAILED: Independently reproduce AC-1 through AC-5 evidence, including the fault cases, mutation reds, exact-head live rejection flow, and the 485-line net reduction.
  AC-1, AC-2, AC-4, and AC-5 reproduced; AC-3's artifact-content proof false-greens under two independent empty-content mutations.
- DONE: Reproduce AC-1 first-run authority and fault classification.
  Exit-23 and hard-deadline cases passed 2/2 with one invocation; retry mutation red on `[1 1]`, and ignored-exit mutation red on an unexpectedly passing scenario.
- DONE: Reproduce AC-2's one fixed deadline without a wait controller.
  The stall emitted `wait_agent` JSONL every 20ms yet expired at the original 2s bound; review found one `CommandContext`, one `cmd.Run`, and no old controller/retry symbols.
- FAILED: Reproduce AC-3's diagnosable same-run fault evidence.
  Actual artifacts are correct, but writing empty `rejection-task.after.md` or empty `git-head.txt` independently left the focused test green; existence alone does not prove retained content.
- DONE: Reproduce AC-4's exact-head live producer proof.
  One fresh run passed in 322.68s under `/tmp/spacedock-codex-validation-d199fc98.ulLI8a`; source HEAD is `d199fc98`, exit is 0/not timed out, cycles are REJECTED then PASSED, Git status is clean, and layout is flat.
- DONE: Reproduce AC-5's simplification evidence.
  The named-file audit is 335 additions, 820 deletions, and 485 net lines removed; removed-symbol sweep is clean and the process path has no replacement poller, watcher, retry, or recovery loop.
- FAILED: Adversarially test whether retries, ignored exits, wait-shaped output, or missing artifacts could still produce a false green; verify one fixed deadline and one authoritative invocation.
  Retry, ignored exit, omitted process-result file, and omitted entity file all red, but zero-byte entity and Git-head artifacts false-green, so missing same-run content is not fully guarded.
- DONE: Confirm the full and race suites, documentation, removed-symbol sweep, clean worktree, and same-run artifact provenance; classify every finding as material, deferred, or polish and recommend PASSED or REJECTED.
  Package, `go test ./...`, and fresh `go test ./... -race -count=1` pass; docs and provenance match, removed symbols are absent, and the implementation worktree is clean.
- FAILED: Material evidence defect — AC-3 artifact-content boundary.
  Supported ordinary-exit/timeout diagnostics can regress to empty entity or Git-head snapshots without the focused test failing; narrowly compare snapshot bytes to the entity and `rev-parse HEAD` from the same helper run.
- DONE: Deferred and polish classification.
  No deferred risks found; polish only: `gofmt -l ./cmd ./internal` reports pre-existing, out-of-scope `internal/release/journeydelta.go`, unchanged from parent `00c2fa6f`.
- FAILED: Recommend PASSED or REJECTED.
  REJECTED because the material AC-3 evidence defect violates the promised mutation-sensitive proof even though the current production and live artifacts are correct.

### Summary

Validation confirms the single-run mechanism, exact-head live outcome, suites, provenance, documentation, and 485-line reduction. Recommendation is REJECTED on one narrow material evidence defect: required entity and Git-head snapshots can be emptied without making the focused fault test red; no outcome defect or deferred risk was found.

## Stage Report: validation (classification correction)

- DONE: Reassess the AC-3 finding using the release-scope classifier.
  The prior material classification and REJECTED recommendation are superseded: both empty-content cases are unreachable code mutations, not released-input or normal-CI triggers.
- DONE: Finding 1 — Released user and normal workflow.
  The only user is a CI maintainer diagnosing a Codex rejection-flow fault; the normal path reads the nonempty entity and writes those exact bytes at `codex_single_run_test.go:107-113`.
- DONE: Finding 1 — Observable harm.
  Current users lose nothing: the fresh live artifact is 1,600 bytes; only a hypothetical future edit that writes different bytes would remove the entity snapshot diagnostic.
- DONE: Finding 1 — Affected value AC or non-negotiable boundary.
  AC-3's diagnostic value would be weakened by that hypothetical edit, but current production satisfies AC-3 and no safety, security, data-integrity, or compatibility boundary fails.
- DONE: Finding 1 — Trigger evidence and classification.
  The green mutation replaced the write argument with `nil` while still returning the original entity in memory; no supported input causes that split, and read/write errors already fail capture, so this is polish-level assertion hardening.
- DONE: Finding 2 — Released user and normal workflow.
  The same CI maintainer consumes `git-head.txt`; production runs `git rev-parse HEAD` and writes its output at `codex_single_run_test.go:119-133`.
- DONE: Finding 2 — Observable harm.
  Current users lose nothing: the fresh live artifact is a 41-byte SHA line; only a hypothetical future edit that discards successful command output would empty it.
- DONE: Finding 2 — Affected value AC or non-negotiable boundary.
  AC-3's Git diagnostic would be weakened by that hypothetical edit, while AC-4 uses the separate nonempty source-HEAD artifact; no non-negotiable boundary currently fails.
- DONE: Finding 2 — Trigger evidence and classification.
  The green mutation assigned `out = nil` after successful `rev-parse`; a valid initialized fixture returns a SHA, and any command or artifact-write error fails capture, so this too is polish rather than material or deferred.
- DONE: Return the corrected verdict.
  PASSED: all value ACs have current behavioral evidence, with no material finding or real deferred trigger; exact artifact-byte comparisons remain optional polish.

### Summary

Plainly, the test can be made stronger against a future programmer deliberately disconnecting two writes from their sources, but CI cannot naturally reach either empty-file state through the implemented path. The completed four-field triage therefore corrects the release-scope classification to polish and the validation verdict to PASSED.
