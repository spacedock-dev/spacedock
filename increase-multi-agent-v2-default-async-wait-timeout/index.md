---
title: Increase multi_agent_v2 default async wait timeout
status: validation
score: 0.5
source: "Captain filing request 2026-07-11."
completed:
verdict:
worktree: .worktrees/spacedock-ensign-increase-multi-agent-v2-default-async-wait-timeout
issue:
id: 95we0fhydgx5rbay5fw3qy4q
started: 2026-07-11T11:51:34Z
mod-block: merge:pr-merge
pr: "#500"
---

`multi_agent_v2`’s short default wait causes repeated timeout churn while subagents perform normal verification or interactive work. Because the wait already returns immediately when an agent responds or new user input arrives, increase the default timeout to five minutes. This reduces unnecessary polling and tool traffic without delaying updates or reducing responsiveness.

## Problem

The current live Codex collaboration surface accepts an omitted `timeout_ms`, but
that omission produced a timeout after about 34 seconds of wall-clock time in
the ideation probe (the call itself has orchestration overhead; this establishes
the practical current default as roughly 30 seconds). A normal ensign that is
verifying, running a focused test, or waiting for an interactive clarification
therefore causes repeated timeout/re-wait churn even though it is healthy.

That churn is not needed for responsiveness. The live `wait_agent` contract
states that a mailbox update and new user input end the wait early. In the probe,
a `timeout_ms: 60000` wait returned `timed_out: false` after about 21 seconds
when a parent mailbox message arrived. The completion/user-input behavior is
therefore distinct from the timeout ceiling. A sibling agent's final status did
not wake this worker's five-minute wait, which is useful negative evidence that
the wait observes the caller's collaboration scope rather than acting as a
global completion bus.

## Evidence and ownership boundary

- **Live surface:** `wait_agent(timeout_ms?)` currently accepts 10,000 through
  3,600,000 milliseconds. Its surfaced behavior says it returns early for
  mailbox updates and steered user input. The no-argument probe ran from
  `2026-07-11T12:01:40.3Z` to `2026-07-11T12:02:14.3Z` and timed out; the
  60-second mailbox probe ran from `12:08:03.3Z` to `12:08:24.3Z` and completed
  early.
- **Actual owner:** the omitted-timeout default is owned by the Codex runtime,
  not this Go repository or the Spacedock plugin. Spacedock currently owns only
  its first-officer instructions and the `spacedock codex` argv it launches.
- **Supported Codex configuration:** the current Codex schema and both local
  executables expose
  `features.multi_agent_v2.default_wait_timeout_ms`, with a 0--3,600,000 ms
  range. On `codex-cli 0.144.1` and the bundled `0.144.0-alpha.4`,
  `--strict-config -c 'features.multi_agent_v2.default_wait_timeout_ms=300000'`
  passed configuration parsing (then correctly stopped only because stdin was
  not a terminal); an unknown nested key and `3600001` each failed validation.
  The initially plausible root key `multi_agent_v2.default_wait_timeout_ms` is
  invalid. Do not use it.
- **Repository boundary:** a user-level `~/.codex/config.toml` changes the
  operator's Codex default. A trusted project `.codex/config.toml` changes only
  that project. Neither makes an installed Spacedock plugin change an already
  open Codex session or every repository where the plugin is invoked. The
  current `runCodex` front door also emits no timeout override.

## Chosen direction

Bind the Spacedock Codex first officer's *foreground wait* explicitly to five
minutes: `wait_agent(timeout_ms: 300000)`. This is the smallest durable
in-repository behavior because it applies wherever the shipped first-officer
skill is used--including an existing Codex session that did not start through
`spacedock codex`--and it uses the already-probed live tool call shape.

Keep the live-surface probe generic (`wait_agent(timeout_ms)`) so a missing or
incompatible tool remains a concrete runtime blocker. Bind the exact
five-minute value in the adjacent `«async-dispatch»` runtime binding; the
generic Codex wait notes remain lifecycle-only. This is a Spacedock per-call
policy, not a claim that Codex's global default changed.

Implementation changes only these product surfaces:

1. `skills/first-officer/references/codex-first-officer-runtime.md` -- bind the
   foreground wait in the adjacent `«async-dispatch»` runtime binding to
   `timeout_ms: 300000`, while preserving timeout and interruption lifecycle
   rules in the generic wait notes.
2. A focused Codex live scenario under `internal/ensigncycle` -- prove a
   delayed owned worker crosses the old short default without a timeout/re-wait
   and receives a non-timeout foreground-wait return after the delayed hold.

Do **not** add a global setting, modify the operator's `~/.codex/config.toml`,
add a repository `.codex/config.toml`, or add `-c` to `runCodex` in this
change.

## Alternatives considered

1. **Session configuration in `runCodex`:** inject
   `-c features.multi_agent_v2.default_wait_timeout_ms=300000`. This can change
   the default for sessions launched by `spacedock codex`, and the exact current
   key has been parser-proven. It does not cover a skill invoked in an existing
   session, changes a broader host-session default than the Spacedock wait, and
   makes the front door fail on older Codex builds that do not support the
   nested feature configuration. Defer unless a later task explicitly wants a
   launcher-session policy with version/override handling.
2. **Check in `.codex/config.toml`:** applies only when the Spacedock source
   checkout itself is the trusted workspace. It would not travel as active
   configuration with the cached plugin, so it is unsuitable for the requested
   operational behavior.
3. **Ask operators to edit global config:** this changes state outside the
   repository, is per-user rather than shipped behavior, and is out of scope.
4. **Leave the timeout omitted:** preserves the current approximately
   30-second churn and fails the operational goal.

## Documentation impact

The shipped Codex adapter is the user-facing contract surface. No CLI output,
command, or site documentation changes because no launcher or configuration
surface changes. The implementation binds the exact per-call policy in the
adjacent async-dispatch runtime binding; the generic `## Codex wait notes`
section stays lifecycle-only:

```diff
@@ ## Runtime implementation
- `«async-dispatch»`: Codex dispatch is async; the spawn call returns a handle and the FO event loop remains available for mailbox notifications, gate work, and captain interaction.
+ `«async-dispatch»`: Codex dispatch is async; the spawn call returns a handle and the FO event loop remains available for mailbox notifications, gate work, and captain interaction. When async multi-agent dispatch is live, foreground waiting MUST call `wait_agent(timeout_ms: 300000)` as the named long-task default.
```

## Acceptance criteria

**AC-1 (value): A healthy owned worker that takes longer than the current short
wait avoids timeout churn.** In an isolated Codex live run, a delayed owned
worker taking at least 45 seconds causes one foreground wait to remain active
through the delayed hold and receive a non-timeout return; no timeout/re-wait
occurs before that first return. The measured hold exceeds the approximately
30-second baseline, so an omitted/default wait can fail this criterion.
Verified by: a focused `//go:build live` Codex scenario that captures the
collaboration stream plus durable workflow state, with an assertion that fails
on a timeout or repeated wait before the first return.

**AC-2 (value): Five-minute waiting observes an early post-hold mailbox
return.** The same delayed-worker run calls `wait_agent(timeout_ms: 300000)`,
returns `timed_out:false` after the delayed-hold end marker rather than waiting
for the five-minute ceiling, and then verifies the committed stage report
separately. The end-of-hold marker is a whole-second value emitted before the
report and commit; it is not a worker-final-status timestamp, and the live
surface may deliver a post-hold mailbox return before the child FINAL_ANSWER.
Verified by: the live scenario's timestamped wait trace, non-timeout output,
delayed-hold marker, and separate stage-report/git-state assertion. It makes
no marker-to-final-completion causality or latency-SLO claim.

**AC-3 (scope and safety): Spacedock's Codex wait policy is explicit without
changing host-wide configuration or interruption semantics.** The shipped
adapter names the per-call 300,000-ms duration, keeps a timeout normal and
retryable, and preserves the explicit rule that an operator interruption does
not fail, close, or redispatch the worker. The implementation changes neither
the global Codex config nor `spacedock codex` argv.
Verified by: the adjacent async-dispatch adapter binding, scoped-diff review
showing no added configuration or front-door override, and the live AC-1/AC-2
run exercising the actual timeout behavior rather than treating instruction
text as behavioral proof.

## Test plan

1. Add the live delayed-worker control first. Use the existing isolated
   `CODEX_HOME` Codex runner so no developer-local config can mask the baseline.
   Delay the owned no-write worker for 45 seconds--longer than the observed
   short default but below the existing 60-second stream-silence guard. Record
   the foreground wait begin/return, timeout/retry count, process exit,
   delayed-hold markers, entity body, and state-checkout git log. The markers
   are end-of-hold values only: do not impose a marker-relative completion
   latency budget. Estimated cost: one Codex run, roughly 1--2 minutes.
2. Do not use contractlint or any instruction-text assertion as AC-1/AC-2
   proof. Keep the live-surface probe and generic wait notes lifecycle-only;
   the adjacent async-dispatch binding owns the exact per-call duration.
3. Run the required offline gates: `go test ./...`, `go test ./... -race`, and
   `gofmt -w ./cmd ./internal`. Because the changed adapter is loaded by the
   Codex live lane, run the focused live control and then
   `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` before validation/merge.
4. If a later task elects the rejected launcher-config alternative, add separate
   parser/argv tests for the exact `features.multi_agent_v2...` key, an
   unsupported-version failure path, user-supplied override precedence, fresh
   vs resume behavior, and a fresh-process live proof. Those tests are not part
   of this per-call change.

## Out of scope

- Changing Codex's product-wide default or asking users to edit global config.
- Adding a user-facing Spacedock timeout flag, profile, or configuration file.
- Changing Claude or Pi wait behavior.
- Treating a sibling/cross-tree completion as the foreground wait's completion
  signal.
- Replacing the existing timeout/interruption lifecycle guarantees.

## Stage Report: ideation

- DONE: Identify the real configuration owner and supported surface without
  editing global settings. The actual key is
  `features.multi_agent_v2.default_wait_timeout_ms`; both local Codex binaries
  accepted 300,000 under strict configuration parsing, while wrong-key and
  out-of-range controls failed.
- DONE: Establish current wait and early-return behavior with the smallest safe
  live probe. An omitted wait timed out after about 34 seconds; a 60-second
  wait returned early on a mailbox message after about 21 seconds. The surfaced
  tool contract also explicitly guarantees early return on steered user input.
- DONE: Choose the smallest compatible repository behavior. Bind a five-minute
  value explicitly in the Codex first-officer adapter, rather than claiming a
  global default changed or injecting a config override that only applies to
  front-door launches.
- DONE: Specify a live delayed-worker control that measures the operational
  value against the observed short baseline and proves a post-hold non-timeout
  return without asserting worker-final-status causality.
- DONE: Trace each acceptance criterion to independent evidence so the
  implementation does not substitute adapter wording or an invented
  marker-relative latency budget for the live timeout behavior.
  - **AC-1:** the isolated 45-second owned-worker run must capture one
    foreground wait spanning the delayed hold and fail if it records a timeout
    or a second wait before its first non-timeout return. The hold is the
    moving baseline: it exceeds the observed approximately 30-second omitted
    wait.
  - **AC-2:** in that same run, record `wait_started`, the delayed-hold end
    marker, and `wait_returned`. The assertion requires non-timeout output and
    a return after the hold marker; report and state-checkout commit verification
    happen separately. The whole-second hold marker precedes those writes and
    cannot prove final-status timing.
  - **AC-3:** the adjacent async-dispatch binding owns the fixed per-call
    shape; scoped-diff review proves no `-c`, configuration, or front-door
    change. Neither text inspection nor a generic contractlint check is used
    as proof of AC-1 or AC-2 behavior.
  The observed 21-second return from a 60-second wait after a parent mailbox
  message is supporting behavior only: it demonstrates that an incoming
  mailbox event wakes the live wait early and is consistent with the surfaced
  user-input rule, but it does not establish child-final-answer causality.

### Summary

Codex owns the omitted wait default, but it exposes a validated session config
surface at `features.multi_agent_v2.default_wait_timeout_ms`. Spacedock should
not alter that global/session default for this task: the robust, compatibility-
first behavior is an explicit `wait_agent(timeout_ms: 300000)` in the shipped
Codex first-officer adapter's async-dispatch binding. A 45-second owned-worker
live control proves that the new five-minute ceiling removes approximately
30-second timeout churn and produces a non-timeout post-hold return, without
claiming a marker-to-final-completion latency SLO.

## Stage Report: implementation

- DONE: Bound the Codex `«async-dispatch»` foreground wait to the named
  long-task per-call policy `wait_agent(timeout_ms: 300000)`. The generic wait
  notes remain lifecycle-only; no global/project Codex configuration or
  `spacedock codex` argv changed.
- DONE: Added the live-only `TestLiveCodexForegroundWaitUsesFiveMinutePerCallPolicy`
  control. It creates one 45-second delayed owned worker, preserves raw Codex
  JSON/rollout records and timing evidence, verifies its report and
  path-scoped state commit, and fails on timeout or re-wait churn before the
  first non-timeout return.
- DONE: Final live control passed after the capability-based adapter wording:
  the worker held for 45 seconds; the first wait used `timeout_ms=300000`, ran
  for 1m25.706s, and returned `timed_out=false` after the delayed-hold marker.
  Evidence is retained at
  `/tmp/spacedock-95w-live-proof-final/codex-shared-scenarios/foreground-wait-timeout/`.
- DONE: Verification passed: `go test ./...`; `go test ./... -race`; focused
  live control; and `go test -tags live -count=1 -timeout 40m -run
  '^TestLiveCodexSharedScenarios$' ./internal/ensigncycle -v` (all nine shared
  scenarios passed in 1103.83s). The new Go test is gofmt-clean.
- DONE: No contractlint change is retained. The live delayed-worker control is
  the behavioral proof; the adapter's async-dispatch binding is documented as
  the policy mechanism, not a text-proof oracle.
- DONE: Code commit `b733ab5c` (`Codex: use long-task foreground waits`).

### Summary

Spacedock now instructs live async multi-agent Codex foreground waits to use a
five-minute per-call ceiling while preserving existing interruption and timeout
lifecycle rules. A delayed owned-worker live scenario proves that the wait
does not churn at the old short horizon and receives a non-timeout post-hold
mailbox return; durable report/commit verification is separate.

### Feedback Cycles

**Cycle 1 — timing-proxy design correction (2026-07-11).** A fresh independent
focused live run went RED at 30.011s against the former 30-second upper bound.
The bound compared `time.Unix(date +%s, 0)` from the end-of-hold marker with
the foreground wait return. That whole-second marker is emitted before report
append and commit; retained runs measured 20.324s and 29.282s, and a direct
trace showed report append 5.567s, commit 17.169s, and child FINAL_ANSWER
28.062s after the marker. The correction removes only this invalid latency SLO,
keeps the explicit 300,000-ms policy and all non-timeout/no-churn assertions,
and restates AC-2 as an early post-hold mailbox return rather than
worker-final-completion causality.

**Cycle 2 — quiet wait epochs (captain correction, 2026-07-12).** The
five-minute wait removes short-timeout churn, but the Codex FO remains too loud
when an ordinary timeout causes it to reinstall the same asynchronous wait.
Treat consecutive waits over the same unresolved worker set and unchanged
workflow state as one monitoring epoch: announce the epoch once when monitoring
starts, then reinstall `wait_agent` silently while that worker set and state
remain unchanged. Speak again only on worker completion, blocker, workflow
state transition, or a periodic heartbeat the captain explicitly requested.
Preserve the existing interruption explanation and `wait_agent` requirement;
an operator interruption still returns control without failing, closing, or
redispatching the worker, and the next idle action reinstalls the wait. This is
a Codex runtime-adapter simplification only. Do not edit the shared core or
weaken the existing timeout/interruption lifecycle rules. Implementation must
update the adapter and focused proof so repeated ordinary timeouts do not emit
repetitive “continuing with …” narration within one unchanged wait epoch.

**Cycle 3 — ordering-invariant escalation (2026-07-12).** Fresh validation
rejected the cycle-2 proof even though the retained run observed the intended
ordering: the assertion receives neither hold marker, so it can pass when a
timeout/re-wait occurs after the worker hold or a workflow-state change. The
missing executable invariant is concrete: require the first wait to overlap the
hold, its ordinary timeout to precede `holdFinished`, and the silent re-wait to
begin before that finish marker; add a fast negative boundary control. This is
the third feedback cycle, so it is escalated to the captain rather than
auto-routed to another implementation pass.

**Cycle 4 — captain scope correction (2026-07-12).** The five-minute per-call
wait policy is intentionally simple and can be verified through normal use and
iterated if it proves wrong. Remove the task-specific heavyweight Codex live
timeout/quiet-epoch test stack rather than carrying roughly a thousand lines of
slow scenario machinery. Retain the explicit policy and its scope boundaries;
do not replace the removed live proof with Contractlint, instruction-file
substring checks, or another prose-only surrogate. Keep only the smallest
meaningful ordinary verification the existing code supports, if any, and record
the captain-directed actual-use/iterate posture in the task's updated plan and
report.

## Stage Report: implementation (cycle 1)

- DONE: Removed only the invalid marker-relative 30-second upper bound.
  The live test now labels `date +%s` values as delayed-hold endpoints, retains
  the explicit 300,000-ms call, 35-second crossing, `timed_out:false`,
  post-hold return, single-owned-worker, no-re-wait, and durable report/commit
  checks, and does not invent a replacement latency SLO.
- DONE: Corrected the design, ACs, test plan, and prior implementation report:
  the exact policy lives in the adjacent async-dispatch binding, generic wait
  notes remain lifecycle-only, and AC-2 claims only an early post-hold mailbox
  return rather than child-final-answer causality.
- DONE: The corrected focused live control passed in 278.36s. Its retained
  evidence records a 45-second hold and a 1m34.656s foreground wait using
  `timeout_ms=300000` with `timed_out=false`; its 32.224-second
  marker-to-return gap is intentionally not treated as a completion SLO.
  Evidence: `/tmp/spacedock-95w-correction-proof/codex-shared-scenarios/foreground-wait-timeout/`.
- DONE: `go test ./...` and `go test ./... -race` passed; the repository-wide
  gofmt pass left the corrected live test clean (an unrelated pre-existing
  alignment change was restored).
- DONE: Code commit `069643f9` (`Codex: clarify foreground wait timing proof`).
- FAILED: The required full shared Codex live suite ran for 1183.01s with 8 of
  9 scenarios passing. Only `filing` failed, independently of this correction:
  the FO correctly used `BIN="${SPACEDOCK_BIN:-spacedock}"` and `"$BIN" new
  wire-the-thing`, but the existing filing parser accepts only an unquoted
  launcher capture. No filing behavior or parser was changed in this task.

### Summary

Cycle 1 removes a flaky, semantically invalid timing proxy without rolling back
the five-minute per-call policy. Focused live and offline/race verification are
green; the full shared lane exposed an unrelated quoted-launcher filing parser
gap that remains outside this task's scope.

## Stage Report: implementation (cycle 2)

- DONE: Repaired the legitimate shared-lane filing assertion gap exposed in
  cycle 1. The parser now recognizes only the observed, contract-blessed
  double-quoted launcher capture
  `BIN="${SPACEDOCK_BIN:-spacedock}"` followed by `"$BIN" new`, alongside the
  original unquoted form. It still ties the call to the exact captured variable.
- DONE: Added an exact JSON-valid positive for the quoted Codex
  capture-and-call idiom. Added RED negative controls for `$BINN new` and
  `${BINN} new`: before the boundary repair, both falsely matched a prior `BIN`
  capture because the call pattern consumed the suffix. The matcher now requires
  a non-variable, non-newline boundary after the exact `$BIN` or `${BIN}`
  expansion. The observed quoted positive remains accepted.
- DONE: Code commit `30e435e4` (`Codex: accept quoted filing launcher capture`)
  changes only `internal/ensigncycle/shared_filing_test.go` and
  `internal/ensigncycle/shared_filing_negative_test.go`; it changes no wait
  policy, adapter, configuration, or contractlint surface.
- DONE: Verification passed: focused
  `TestAssertCodexFilingViaNew` (latest 0.376s), `go test ./...`, and
  `go test ./... -race`. `gofmt -d` on the two changed Go files and
  `git diff --check` were clean. The focused five-minute foreground-wait live
  control also passed in this worktree (276.89s), retaining evidence at
  `/tmp/spacedock-95w-filing-repair-proof/codex-shared-scenarios/foreground-wait-timeout/`.
- DONE: The repaired `filing` scenario passed in the required shared live run
  (83.90s), proving the real quoted `"$BIN" new wire-the-thing` producer stream
  is accepted. The full retained run is
  `/tmp/spacedock-95w-filing-boundary-shared-live/codex-shared-scenarios/`.
- BLOCKED: That full shared live run exited nonzero after 559.22s for external
  Codex capacity, not an assertion or product behavior failure. It passed
  `gate-guardrail`, `rejection-flow`, `feedback-3-cycle-escalation`,
  `merge-hook-guardrail`, and `filing`; `shallow-boot`,
  `self-evidence-merge-triage`, `smallest-sufficient-mechanism`, and
  `keep-moving-posture` each emitted only `thread.started`, the standard
  `multi_agent_v2` warning, `turn.started`, then
  `You've hit your usage limit. Visit https://chatgpt.com/codex/settings/usage
  to purchase more credits or try again at Jul 12th, 2026 12:27 AM.`, followed
  by `turn.failed`. None reached a tool call, skill execution, or final-message
  write, so each runner then reported the missing `codex-final-message.txt`.
  The test-only filing matcher cannot causally affect that pre-agent capacity
  failure.
- BLOCKED: An isolated fresh `shallow-boot` rerun reproduced the same missing
  final-message result and retained the same usage-limit stream at
  `/tmp/spacedock-95w-filing-boundary-isolated-shallow-boot/codex-shared-scenarios/shallow-boot/`.
  A `self-evidence-merge-triage` rerun had already started when the capacity
  diagnosis was confirmed; it was intentionally interrupted (exit `signal:
  interrupt`, 76.924s) to avoid blind retries before the stated reset. It is not
  pass/fail evidence.
- TODO: Leave this entity in `implementation`. Rerun the required shared Codex
  lane after capacity is available; do not unhold or claim the full live gate
  green from the partial result.

### Summary

Cycle 2 preserves the five-minute wait implementation and fixes the exact
quoted-launcher parser hole that the prior live filing stream demonstrated. The
repair is unit-, offline-, race-, focused-live-, and filing-live proven. The
remaining full-lane result is explicitly capacity-blocked, with retained
pre-agent usage-limit evidence, so the workflow remains in implementation.

## Stage Report: implementation (cycle 3)

- DONE: The required full shared Codex live suite passed after the quota reset:
  `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-95w-postreset-shared-live-20260712 go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`.
  It exited 0; `TestLiveCodexSharedScenarios` passed in 1192.11s and the package
  completed in 1192.362s.
- DONE: All nine scenarios passed: `gate-guardrail` (81.69s),
  `rejection-flow` (269.70s), `feedback-3-cycle-escalation` (73.41s),
  `merge-hook-guardrail` (45.46s), `filing` (57.74s), `shallow-boot`
  (132.07s), `self-evidence-merge-triage` (71.45s),
  `smallest-sufficient-mechanism` (220.57s), and `keep-moving-posture`
  (239.45s).
- DONE: Retained artifacts, including nine final-message files, are at
  `/tmp/spacedock-95w-postreset-shared-live-20260712/codex-shared-scenarios/`.
- NOTE: This verification rerun changed no code. The entity remains in
  `implementation` by direction; this report does not change its status.

## Stage Report: validation

- DONE: Independently inspect the quoted-launcher matcher change: observed quoted positive, exact-variable prefix negatives, and retained manual-filing negatives.
  The retained filing stream invokes `launcher="${SPACEDOCK_BIN:-spacedock}"` then `"$launcher" new wire-the-thing`; fresh normal and race runs of `TestAssertCodexFilingViaNew` pass its quoted positive, `$BINN`/`${BINN}` negative, and `--next-id` manual-flow negatives.
- DONE: Verify the foreground wait policy remains capability-scoped and the focused wait/live evidence proves tool-call behavior rather than instruction-file text.
  The adapter changes only the `«async-dispatch»` per-call policy; its isolated rollout records one spawn and a first `wait_agent(timeout_ms:300000)` lasting 1m19.08s with `timed_out:false`, after the 45s hold and before any later wait.
- DONE: Inspect the retained post-reset nine-scenario live artifacts and implementation test results; run the smallest fresh focused checks needed, identify any material regression, and recommend PASSED or REJECTED.
  All nine retained scenario streams have a final message and `turn.completed`; the recorded full lane exited 0. Fresh `go test ./... -count=1`, `go test ./... -race -count=1`, live-tag compile, targeted gofmt, and `git diff --check` found no regression.
- DONE: **AC-1 (value):** A healthy 45-second owned worker avoided short-horizon timeout/re-wait churn.
  `codex-foreground-wait-rollout.jsonl` shows the first non-timeout foreground wait crossing the 35-second control and no later wait before its return.
- DONE: **AC-2 (value):** The five-minute call returned after the hold marker without a completion-latency claim.
  The evidence records a 45s hold, a 1m19.08s first wait, `timeout_ms=300000`, `timed_out=false`, and separately recorded path-scoped state commits.
- DONE: **AC-3 (scope and safety):** The policy is explicit and scoped without a global config or launcher-argv change.
  The branch diff from `f22360de` contains only the adapter, focused live control, and filing assertion repair; the live rollout supplies behavior proof rather than text matching.

### Summary

Independent validation found the quoted-launcher repair correctly bounded and the wait policy behaviorally supported by retained live tool-call evidence. The clean implementation branch passes the fresh offline/race gates; no material regression was found.

Recommendation: PASSED.

## Stage Report: implementation (cycle 4)

- DONE: Rebasing 95w onto current `origin/main` retained both the current
  named-capability/runtime-evidence structure and the five-minute per-call wait
  policy, without `-X`, skip, force push, or a PR. Replayed commits:
  `0e2c0d1e`, `0a3ba677`, and `c5565289`.
- DONE: Preserved main's stronger segmented launcher matcher and restored 95w's
  exact-variable `$BINN` / `${BINN}` negative coverage. Focused
  `TestAssertCodexFilingViaNew` passed.
- DONE: Implemented captain feedback cycle 2 in the Codex adapter only:
  unchanged-worker/state waits form one monitoring epoch; the FO announces once
  at epoch start and silently reinstalls an ordinary timed-out wait.
- DONE: Added `TestLiveCodexForegroundWaitRetriesQuietlyWithinEpoch` plus
  fast positive/negative trace controls. The control raises only its test-only
  stream-silence watchdog so a real five-minute wait is observable.
- DONE: Fresh quiet-epoch live proof passed in 495.19s. Its first
  `wait_agent(timeout_ms:300000)` returned `timed_out:true` at 03:30:16Z; the
  second began at 03:30:18Z, and the next FO narration occurred only after that
  re-wait completed. Evidence:
  `/tmp/spacedock-95w-quiet-epoch-green-20260712/codex-shared-scenarios/foreground-wait-quiet-epoch/`.
- DONE: Fresh original short-hold proof still passed in 219.46s, with a 45s
  hold and non-timeout foreground wait. Evidence:
  `/tmp/spacedock-95w-cycle2-short-focus-20260712/codex-shared-scenarios/foreground-wait-timeout/`.
- DONE: Verification passed sequentially: `go test ./...` (2,154),
  `go test ./... -race` (2,154), targeted live trace/lifecycle checks, and
  `gofmt -w ./cmd ./internal`; the known unrelated
  `internal/release/journeydelta.go` alignment drift was restored.
- SKIPPED: Completing the broad shared Codex lane begun before feedback.
  It was captain-interrupted at gate-guardrail after 45.361s so it is not
  evidence; the changed behavior has fresh focused live proof instead.
- DONE: Code commit `3043aff` (`Codex: quiet foreground wait epochs`).

### Summary

The rebase now carries the five-minute policy on current main, and the Codex
adapter explicitly treats repeated ordinary timeouts as one quiet monitoring
epoch. The live baseline was already quiet on this host; the new focused control
makes that behavior executable and guards it through a real timeout/re-wait.

## Stage Report: validation (cycle 2)

- DONE: Independently reviewed the rebased cycle-2 diff. It retains main's
  named capability structure, keeps the explicit `wait_agent(timeout_ms:
  300000)` policy scoped to the Codex adapter, and adds the quiet-monitoring
  rule without a launcher/configuration change or shared-core edit.
- DONE: Ran the fresh focused live quiet-epoch control with a launcher built
  from this worktree: `SPACEDOCK_BIN=/tmp/spacedock-95w-cycle2-validation-bin
  SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-95w-cycle2-validation-rerun-20260712
  go test -tags live -count=1 -timeout 25m -run
  '^TestLiveCodexForegroundWaitRetriesQuietlyWithinEpoch$'
  ./internal/ensigncycle -v`. It passed in 485.26s. The retained rollout
  observed a `300000`-ms wait from 04:21:44Z to a `timed_out:true` return at
  04:26:44Z, a second `300000`-ms wait beginning at 04:26:46Z, and no
  captain-visible FO narration between those two wait calls.
- DONE: The fresh run's retained artifact independently shows the intended
  concrete ordering: the worker's six-minute hold spans
  `CODEX_FOREGROUND_WAIT_STARTED_EPOCH=1783830120` through
  `CODEX_FOREGROUND_WAIT_FINISHED_EPOCH=1783830480`; the first wait starts
  before that interval, times out while it is still open, and the silent
  re-wait begins before the finish marker. Artifact:
  `/tmp/spacedock-95w-cycle2-validation-rerun-20260712/codex-shared-scenarios/foreground-wait-quiet-epoch/`.
- DONE: Offline and race gates pass on the clean rebased branch: `go test
  ./... -count=1` and `go test ./... -race -count=1` each passed 2,154 tests
  in 17 packages. Live-tag quiet-trace unit and race controls passed; `gofmt
  -d` on every changed Go file and `git diff --check origin/main...HEAD` were
  clean.
- FAILED: The cycle-2 focused test does not make its observed wait/hold
  ordering a failing invariant. `runCodexQuietEpochScenario` reads
  `holdStarted` and `holdFinished` only to check the six-minute duration, then
  calls `assertCodexQuietEpochRollout(result.rollout)` without either marker.
  That assertion can pass if the first timed-out wait or its silent re-wait is
  no longer over the same unresolved worker/state--for example, if either
  begins after the hold/report boundary. The fresh artifact happens to have
  the right order, but a control that can pass after that regression does not
  prove the required unchanged monitoring epoch.
- TODO: Return to implementation. Thread the hold markers into the rollout
  assertion and make it fail unless the first wait overlaps the hold, its
  ordinary timeout occurs before `holdFinished`, and the silent second wait
  begins before `holdFinished` (with a fast negative boundary control). This
  keeps the test tied to an unresolved worker and unchanged workflow state;
  then rerun the focused live control and applicable gates.

### Summary

The concrete behavior was observed twice and the branch is clean, but the
cycle-2 regression control does not enforce the crucial temporal relation that
makes the re-wait a quiet continuation of the same monitoring epoch.

Recommendation: REJECTED.

## Stage Report: implementation (cycle 5)

- DONE: Threaded `holdStarted` and `holdFinished` from the committed worker
  report into the quiet-epoch rollout assertion. It now requires the first wait
  to overlap the unresolved hold, its ordinary timeout to return before hold
  finish, and the silent re-wait to start before that finish marker.
- DONE: Added fast full-assertion boundary controls for a missing overlap, a
  timeout at the hold boundary, and a re-wait at the hold boundary. Each must
  fail; the inside-hold positive remains green.
- DONE: The focused live quiet-epoch control passed in 516.16s with the timing
  invariant active. The hold markers are 1783836755 through 1783837115; the
  first 300,000-ms wait timed out at 06:17:21Z and the silent re-wait began at
  06:17:23Z, both before hold finish. Evidence:
  `/tmp/spacedock-95w-cycle3-quiet-epoch-live-20260712/codex-shared-scenarios/foreground-wait-quiet-epoch/`.
- DONE: Fast live-tag controls, including the three boundary negatives and
  narration controls, passed 7/7 both normally and with `-race`.
- DONE: `go test ./...` and `go test ./... -race` each passed 2,154 tests in
  17 packages. `gofmt -w ./cmd ./internal` was run; only the known unrelated
  `internal/release/journeydelta.go` alignment drift was restored.
- DONE: No global Codex configuration, launcher argv, or shared-core behavior
  changed. Code commit `b305da7` (`Codex: bind quiet waits to active hold`).

### Summary

Cycle 3 turns the previously observed ordering into a failing invariant, so a
quiet re-wait can pass only while the same worker is still held and unresolved.
The feedback correction is scoped to the focused Codex proof and preserves the
five-minute per-call policy and interruption lifecycle.

## Stage Report: validation (cycle 3)

- DONE: Independently verify the quiet-epoch assertion consumes hold markers and fails on missing overlap, timeout-at-finish, and re-wait-at-finish boundaries.
  Fresh normal and `-race` focused controls passed the inside-hold positive, all three named boundary negatives, and both narration controls.
- DONE: Run a fresh real timeout/re-wait live proof that demonstrates overlap, timeout before hold completion, silent re-wait before completion, and no intervening FO narration.
  Fresh rollout: first 300000-ms wait 06:48:18.609Z–06:53:18.624Z timed out during the 06:48:36Z–06:54:36Z hold; silent re-wait began 06:53:20.392Z and returned after the hold, with no `agent_message` between the two wait events.
- DONE: Re-run applicable standard, race, and live-tag controls; assess every acceptance criterion and return PASSED or REJECTED without editing product code.
  `TMPDIR=/tmp go test ./... -count=1`, its `-race` counterpart, the focused live-tag normal/race controls, and the fresh 45-second live policy control all passed; validation made no product edit.
- DONE: **AC-1 (value):** A healthy owned worker longer than the old short wait avoids timeout churn.
  The fresh 45-second hold was inside its first 300000-ms wait, which lasted 1m31.21s and returned `timed_out:false`; its only later wait began after that first return.
- DONE: **AC-2 (value):** Five-minute waiting observes an early post-hold mailbox return.
  The same fresh rollout records the hold ending at 15:04:53+08:00 and the first non-timeout wait returning at 15:05:27.913+08:00, with the report's path-scoped state commit separately recorded.
- DONE: **AC-3 (scope and safety):** Spacedock's Codex wait policy is explicit without changing host-wide configuration or interruption semantics.
  `origin/main...HEAD` changes only the adapter and focused tests; no config or front-door file changed, while the live runs exercised explicit 300000-ms calls rather than treating adapter text as behavioral proof.

### Summary

Cycle 3 closes the prior material validation gap: the current assertion consumes
the worker hold markers and rejects each late/non-overlapping boundary. Fresh
live evidence independently reproduces both the long-wait policy and the quiet
timeout/re-wait ordering. Recommendation: PASSED.
