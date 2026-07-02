---
title: "codex-live CI false reds: reuse characterization and shallow-boot gh probe shape"
status: validation
group: tooling
source: "PR #464 codex-live (2026-07-02, run 28568410914): TestLiveCodexSharedScenarios/rejection-flow failed 'addressable-worker was characterized ABSENT, but transcript contains turn-starting reuse tool send_input' — while the scenario's behavior was CORRECT (two-cycle rejection flow, kept-alive reviewer, clean tree). Root cause read from the artifact: codexNarrationNegatesReuseRoute (internal/ensigncycle/shared_reviewer_reuse_test.go:391) flags any single FO message containing a reuse concept AND any negation token; both tripping messages AFFIRM reuse with unrelated negations ('does not close or redispatch the worker', 'the host has no dedicated shutdown tool ... the two reusable workers'). Same over-broad live-heuristic family as the 2w wrong-root detector fixed in #462. Narration-dependent: can red any branch."
id: mq42796928asq686vxs74mpm
started: 2026-07-02T07:41:51Z
worktree: .worktrees/spacedock-ensign-codex-reuse-characterization-negation-scope
---

## Problem
Two unrelated codex-live failures share the same user-visible symptom: the live lane goes red even though the run made the behavior Spacedock wanted to prove. Keep them in one task because both are harness/bootstrap characterization bugs, both can red any unrelated branch, and both should be closed with offline fixtures before spending another codex-live confirmation run.

First, the Codex rejection-flow runner characterizes the addressable-worker capability from FO narration. `assertCodexReviewerReuse` first asks `codexAddressableWorkerAbsent`; when that returns true it runs the fresh-dispatch oracle and treats any `followup_task`/`send_input` reuse tool as a contradiction. The current detector (`codexNarrationNegatesReuseRoute`) returns true when one FO `agent_message` contains any reuse-route concept ("follow-up", "followup", "addressable", "reuse", "reusable") and any negation token ("no ", "not ", "cannot", "n't", "without ", "never ") anywhere in the same message. PR #464 and PR #465 both produced affirmative reuse narrations with unrelated negations, so the harness chose the ABSENT path and then failed on the real reuse tool it should have accepted.

Observed false-positive messages to preserve as must-NOT-match fixtures:

- PR #464: "does not close or redispatch the worker" in an otherwise affirmative kept-alive reviewer message.
- PR #464: "the host has no dedicated shutdown tool ... the two reusable workers" in a message about reusable workers, not absent reviewer reuse.
- PR #465: "The validation reviewer is being reused for re-review only; the implementation worker that applied the fix is not doing validation."

Second, PR #465's codex-live `shallow-boot` run booted with PR #42 reported as `MERGED`, but `state sweep --workflow-dir .` reported `merge state UNKNOWN -- gh unavailable; sweep skipped, not empty.` The fixture `gh` stub emits plain `MERGED`, matching `status --boot`'s `gh pr view ... --json state --jq .state` call. `state sweep` uses `dispatch.GhRunnerExec`, which calls `gh pr view ... --json state` and JSON-unmarshals `{"state":"MERGED"}`. A local reproduction showed the split exactly: plain stub => boot sees `MERGED` while sweep sees `UNKNOWN`; JSON stub => sweep finds the merged PR. The harness needs one PR-state probe contract, or at minimum one fixture stub that supports both contracts, so the before-greet startup archive path is tested instead of bypassed by manual frontmatter edits.

## Proposed approach
Use one implementation scope: "codex-live harness false-red hardening." Do not split shallow-boot out unless implementation discovers the PR-state probe parity change is larger than a small helper/fixture fix. The two changes are independent in code, but they share the same validation goal: offline replay of PR #464/#465 failures must fail before the fix and pass after it, then one codex-live run confirms the live lane.

Recommended reuse fix: replace broad concept-plus-negation matching with an explicit "absence declaration" detector. It should match the existing live true-absence wordings in `TestAssertCodexReviewerReuse` ("no follow-up/send binding", "not addressable", "no followup_task or equivalent turn-starting reuse route", "reuse is not supported", "cannot address the kept-alive reviewer") while rejecting affirmative reuse sentences that contain unrelated negations. A bounded-window regex is acceptable only if the PR #464/#465 false-positive phrases are checked in as negatives and all current live-absence positives remain green. Do not key the whole branch on generic words like "no", "not", or "without" unless they bind to an absent route/tool/capability phrase.

Recommended PR-state fix: make boot and sweep consume the same effective `gh pr view` contract in tests and production. The smallest durable option is to share a helper or teach `GhRunnerExec` to request the same `--jq .state` plain state that `status --boot` already uses, because then the existing shallow-boot stub models one real command shape. If changing `GhRunnerExec` is risky for other callers, the fallback is a dual-shape stub plus a parity test proving boot and sweep agree under both plain and JSON outputs. The task is not complete if it only changes the live fixture without an offline parity test; that would leave the next probe-shape drift unpinned.

Coordination update: checked `s0`/`s05`, the active task `error-path-guidance-in-binary-output` on PR #465. Its branch `spacedock-ensign/error-path-guidance-in-binary-output` already contains the shallow-boot probe-shape fix in commit `868fd10f` ("Feedback cycle 2: fix ghRunnerExec to parse the --jq .state shape, not a JSON envelope"). In that worktree, `GhRunnerExec` calls `gh pr view PR --json state --jq .state`, returns trimmed uppercase text, and `internal/dispatch/sweep_test.go` includes `TestGhRunnerExecParsesJqExtractedState` plus `TestSweepWithRealGhStubDoesNotReportUnknown`. The same branch still has the old `codexNarrationNegatesReuseRoute` broad concept-plus-negation detector, so it does not cover the reuse-characterization false red. If `s0` lands first, `mq` should not re-implement the PR-state fix; it should verify/reuse those tests and focus new code on the Codex reuse narration classifier. If `s0` is not yet on the target branch, `mq` should either depend on/rebase over `s0` or explicitly import the same small `GhRunnerExec` parity patch after coordination.

Riskiest path spike: before implementation chooses the mechanism, add or sketch a tiny local harness that runs the same workflow with a stub `gh` that emits plain `MERGED`, then exercises both `status --boot --json` and `state sweep --json`. Record the current failing observation (boot entry state `MERGED`, sweep `swept=[]` or UNKNOWN reason) and make that the first red test. No external GitHub access is needed.

Documentation impact: none expected. This is harness/internal command-surface behavior, not a user-facing CLI output change. If implementation changes `state sweep` JSON/text output or documented `gh` expectations, update this task before coding the docs.

## Acceptance criteria
- **AC-1**: Codex reuse characterization no longer false-ABSENTs on captured affirmative reuse narration with unrelated negation.
  Test: table fixtures in `internal/ensigncycle/shared_reviewer_reuse_table_test.go` replay the two PR #464 tripping messages and the PR #465 "being reused ... is not doing validation" message in transcripts that include one validation `spawn_agent` plus a turn-starting `send_input`/`followup_task` to that validation thread. All three return success from `assertCodexReviewerReuse`; on the current code at least one returns `Codex addressable-worker was characterized ABSENT`.

- **AC-2**: Genuine Codex absence declarations still choose the ABSENT oracle and require fresh validation dispatch.
  Test: the current live-absence wording fixtures remain present and green, including at least one "no follow-up/send binding", one "not addressable", one "not supported", and one "cannot address" phrasing. A transcript with an absence declaration plus only one validation spawn still fails, and an absence declaration plus a real reuse tool still fails, so the fix does not weaken the contradiction checks.

- **AC-3**: The net false-red count for the reduced PR #464/#465 reuse corpus is zero while true-red controls stay nonzero.
  Test: a focused offline test reports three captured false-positive transcripts accepted and at least two negative controls rejected: fresh cycle-2 validation spawn masquerading as reuse, and uncorrelated `send_input` to a non-validation thread. This is the measured value for the task: the false-red corpus moves from >0 failures to 0 without moving the true-red controls to 0.

- **AC-4**: `status --boot` and `state sweep` agree on the same stubbed merged PR state.
  Test: an offline parity test builds a split-root shallow-boot-style workflow with a PR-bearing non-terminal entity and one `gh` stub. With the stub reporting `MERGED`, `status --boot --json` exposes `pr_state.status=="ok"` and an entry state of `MERGED`; `state sweep --json` reports exactly one swept entity for the same PR. The test must fail on the observed split where boot reads plain `MERGED` but sweep expects JSON. If `s0`/commit `868fd10f` is already merged into the target branch, this AC is satisfied by verifying its two `internal/dispatch/sweep_test.go` regressions rather than adding duplicate tests.

- **AC-5**: The shallow-boot durable archive assertion covers the PR-state parity path.
  Test: the codex/host-neutral shallow-boot fixture uses the parity-safe stub and `assertShallowBoot` observes the merged-PR entity terminalized, `verdict: PASSED`, `mod-block:` cleared, and archived before the greet. The gate entity remains unchanged and unarchived, with no worktree created.

- **AC-6**: One codex-live confirmation run passes both affected scenarios after the offline fixtures are green.
  Test: run the codex-live shared scenario lane, or a targeted equivalent that includes `rejection-flow` and `shallow-boot`, and attach the artifact/run ID. The confirmation must show the validation reviewer reuse is accepted and the shallow-boot merged PR is archived before greet; transcript prose alone is not sufficient.

## Test plan
1. Add focused red fixtures for `codexNarrationNegatesReuseRoute`/`assertCodexReviewerReuse` before changing the detector. Keep the fixtures as complete JSONL snippets, not prose-only string searches, so they exercise the branch that failed live.
2. Update the detector with the smallest explicit-absence matcher that preserves all current true-absence rows. Run `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse'`.
3. Add a boot/sweep parity test near `internal/status/live_prstate_pin_test.go`, `internal/dispatch/sweep_test.go`, or a small cross-package CLI-level test if needed. Prefer a CLI-level test because the failure is between two command surfaces, not inside only one package.
4. Check whether `s0`/PR #465's commit `868fd10f` is already present. If present, run its `GhRunnerExec`/sweep regression tests and treat AC-4/AC-5 as dependency verification; if absent, coordinate with `s0` before copying the same `--jq .state` parity patch into `mq`. Run the focused parity test and `go test ./internal/status ./internal/dispatch ./internal/cli`.
5. Run `go test ./internal/ensigncycle` for shared assertion coverage, then `go test ./...` as the implementation gate.
6. Run one codex-live confirmation for `rejection-flow` and `shallow-boot` after offline tests pass. Preserve durable evidence: process exit, artifact path/run ID, entity/archive state for shallow-boot, and the reuse assertion result for rejection-flow.

## Stage Report: ideation

- DONE: firm one coherent task scope for both codex-live false reds, or explicitly justify splitting one out.
  Scoped both failures as one codex-live harness/bootstrap false-red hardening task; split only if PR-state parity proves larger than a small helper/fixture fix.
- DONE: write behavior-first ACs and tests that replay/reduce PR #464 and PR #465 failures and fail on the observed false-reds.
  Added AC-1 through AC-6 with JSONL reuse fixtures, boot/sweep parity proof, shallow-boot durable archive proof, and one codex-live confirmation.
- DONE: exercise or specify the riskiest boot/sweep gh parity path before choosing the implementation mechanism.
  Specified the first red test: one plain-MERGED `gh` stub drives both `status --boot --json` and `state sweep --json`, reproducing boot MERGED plus sweep empty/UNKNOWN before the mechanism is chosen.

### Summary

The task body now treats PR #464's reuse-characterization failure and PR #465's reuse plus shallow-boot failures as one harness-stability scope. The proposed approach favors offline fixtures first, pins the exact false-positive narrations and PR-state probe-shape split, and leaves one codex-live run as field confirmation rather than primary proof.

## Stage Report: ideation (cycle 2)

- DONE: Check whether task `s0` or its branch/artifacts already contain a relevant fix.
  `s0` resolves to `s05` / `error-path-guidance-in-binary-output`; branch `spacedock-ensign/error-path-guidance-in-binary-output` is clean at `868fd10f` and contains the `GhRunnerExec --jq .state` fix plus two sweep regressions.
- DONE: Keep the original `mq` ideation scope while accounting for `s0`.
  Added the coordination update: shallow-boot parity is dependency/verification work if `s0` lands first; the remaining new implementation is the Codex reuse narration classifier, which `s0` still leaves unchanged.
- DONE: Preserve behavior-first ACs without proposing duplicate implementation work.
  AC-4 and Test plan step 4 now direct implementers to verify/reuse `s0`'s regressions when present, or coordinate before importing the same patch.

### Summary

The amendment narrows duplicate-work risk without dropping the PR #465 shallow-boot requirement from the task. `mq` remains the combined false-red task, but the recommended path is now reuse-characterization code plus integration verification of `s0`'s shallow-boot `GhRunnerExec` fix when that branch is available.

## Stage Report: ideation (format repair)

- DONE: Update acceptance criteria to scanner-recognized bold `**AC-N**` form.
  Rewrote AC-1 through AC-6 prefixes only; behavior and test text are unchanged.
- DONE: Preserve the intended task scope.
  No problem statement, approach, AC semantics, or test-plan content was changed beyond AC marker formatting.
- DONE: Commit the state change path-scoped.
  State commit records this format repair for the `mq` entity only.

### Summary

Repaired the acceptance-criteria marker format so scanners can recognize AC-1 through AC-6. The repair is intentionally formatting-only and leaves the combined false-red scope intact.

## Stage Report: implementation

- DONE: fix Codex reuse characterization so captured affirmative reuse narration with unrelated negation no longer enters the ABSENT branch, while true absence controls still fail/pass as intended.
  Commit 063a8b9a adds PR #464/#465 JSONL reuse rows; `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse'` went RED on false-absence failures, then GREEN with 26 passed.
- DONE: account for `s0`/commit `868fd10f` shallow-boot gh probe parity work by reusing/verifying it if present, or coordinating before importing the same minimal fix.
  Verified `868fd10f` exists on `spacedock-ensign/error-path-guidance-in-binary-output` but is not an ancestor of this branch; ported the minimal `GhRunnerExec --jq .state` parity patch and two regressions in 063a8b9a.
- DONE: add/update focused offline tests for the PR #464/#465 reduced failures and run the relevant Go test gates before writing the stage report.
  Focused red/green covered Codex reuse and dispatch gh parity; broader gates passed: `go test ./internal/ensigncycle` 252, `go test ./internal/dispatch ./internal/status ./internal/cli` 1291, `go test ./...` and `go test ./... -race` 1906.

### Summary

Implemented the Codex reuse narration classifier as explicit absence-declaration matching, so affirmative reuse messages with unrelated negation no longer force the ABSENT oracle while genuine absence rows still exercise the fresh-dispatch contradiction checks. Also aligned `state sweep`'s production gh probe with boot's `--jq .state` shape by porting the known `s0` parity fix because that commit was not reachable from this worktree branch.

## Stage Report: validation

- DONE: inspected implementation commit `063a8b9a` for scope and regression risk.
  The diff is limited to `internal/ensigncycle/shared_reviewer_reuse_test.go`, `internal/ensigncycle/shared_reviewer_reuse_table_test.go`, `internal/dispatch/reconcile.go`, and `internal/dispatch/sweep_test.go`. The reuse change replaces broad concept-plus-negation matching with explicit absence-declaration regexes; the dispatch change makes `GhRunnerExec` call the same `gh pr view PR --json state --jq .state` shape used by boot.
- DONE: verified AC-1 through AC-3 for Codex reuse characterization.
  `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse'` passed with 26 cases. The table includes the two PR #464 affirmative-reuse narrations, the PR #465 "being reused ... is not doing validation" narration, live true-absence wordings, absence contradiction controls, a fresh-cycle-2 validation spawn control, and an uncorrelated `send_input` control.
  AC-2 Verified by: the same `go test ./internal/ensigncycle -run 'TestAssertCodexReviewerReuse'` pass covers true-absence rows (`no follow-up/send binding`, `not addressable`, `not supported`, `cannot address`) plus the absence contradiction controls in `internal/ensigncycle/shared_reviewer_reuse_table_test.go`.
- DONE: verified AC-4 and AC-5 for gh probe parity and shallow-boot archive coverage.
  `go test ./internal/dispatch -run 'TestGhRunnerExecParsesJqExtractedState|TestSweepWithRealGhStubDoesNotReportUnknown'` passed with 2 cases. `internal/status/boot.go` and `internal/dispatch/reconcile.go` now use the same `--jq .state` command shape, `internal/ensigncycle/shallow_boot_fixture_live_test.go` uses the bare `MERGED` stub, and `go test ./internal/ensigncycle -run 'TestShallowBootNegativeBrokenEndStates'` passed, preserving the durable archive assertion that requires the merged entity to be archived, terminal, `verdict: PASSED`, and `mod-block:` cleared while the gate entity stays unchanged.
- DONE: ran broader gates for shared-risk coverage.
  `go test ./internal/ensigncycle` passed with 252 cases; `go test ./internal/dispatch ./internal/status ./internal/cli` passed with 1291 cases; `go test ./...` passed with 1906 cases; `go test ./... -race` passed with 1906 cases. `git diff --check` was clean. A read-only `gofmt -l ./cmd ./internal` reported pre-existing formatting drift in `internal/cli/pi_frontdoor_test.go`, which is outside this task's diff; validation did not rewrite unrelated code.
- GAP: AC-6 live Codex confirmation is no longer auth-blocked, but is still not cleanly satisfied by a single combined green run.
  The existing isolated Codex OAuth path was usable: `~/.codex/auth.json` was present and the harness seeded isolated `CODEX_HOME` successfully. A preserved single-scenario run, `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/mq-codex-live-20260702T095355Z go test -tags live -count=1 -timeout 40m -run "TestLiveCodexSharedScenarios/rejection-flow$" ./internal/ensigncycle -v`, passed (`PASS ok github.com/spacedock-dev/spacedock/internal/ensigncycle 504.981s`) with artifacts at `/tmp/mq-codex-live-20260702T095355Z/codex-shared-scenarios/rejection-flow/attempt-1/codex-exec.jsonl` and `codex-final-message.txt`; the final message says cycle 2 used a separate fresh reviewer because this isolated Codex surface had no addressable reviewer reuse route. A preserved combined run, `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/mq-codex-live-both-20260702T100236Z go test -tags live -count=1 -timeout 40m -run "TestLiveCodexSharedScenarios/(rejection-flow|shallow-boot)$" ./internal/ensigncycle -v`, failed overall: `rejection-flow` failed after a no-progress stall while `shallow-boot` passed in 119.59s, with artifacts under `/tmp/mq-codex-live-both-20260702T100236Z/codex-shared-scenarios/rejection-flow/attempt-1`, `attempt-2`, and `/tmp/mq-codex-live-both-20260702T100236Z/codex-shared-scenarios/shallow-boot/codex-exec.jsonl`.
- Recommendation: NOT READY for final approval unless the captain waives AC-6 or routes the remaining live-lane instability / isolated-config addressable-reuse proof through the multi_agent_v2 launcher task.

### Summary

Validation found commit `063a8b9a` scoped to the requested false-red fixes and reproduced the focused and broad offline evidence for AC-1 through AC-5. The auth blocker for AC-6 is gone, and live artifacts now exist, but the combined `rejection-flow` plus `shallow-boot` confirmation is not green and the isolated Codex config still does not prove v2 addressable reviewer reuse.
