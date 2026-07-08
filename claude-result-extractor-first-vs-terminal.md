---
id: 6h08n9jrwa9g5kgm3b3fy8vr
title: "Claude final-message extractor returns the first result event, not the terminal one"
status: ideation
source: "Codex-session CI investigation of PR #483's claude-live (opus) red (2026-07-08), independently re-verified file:line by the FO. Confirmed unrelated to PR #483's own commit (43396704, scoped only to codex_liveenv.go/codex_liveenv_test.go) — this is pre-existing shared test-harness infrastructure, not a regression from that change."
started: 2026-07-08T04:10:04Z
completed:
verdict:
score:
worktree:
issue:
---

## Problem statement

`extractClaudeFinalMessage` (`internal/ensigncycle/claude_final_message_impl_test.go:55-89`) loops a Claude stream-json transcript and returns on the FIRST `result`-type event with `IsError == false` — the early `return result.Result, nil` at line 72, with no continuation to check for a later result event. The function's own doc comment (lines 44-51) says "(1) a **terminal** result/success event's `result` field ... is preferred," and *terminal* means last. The code implements first; the doc comment specifies last. They contradict each other, and the code is wrong.

When a transcript contains more than one non-error success `result` event, the extractor silently returns the wrong one. Observed live: PR #483's failing `claude-live (claude-opus-4-8, CI-E2E-OPUS)` job (`TestLiveClaudeSharedScenarios/keep-moving-posture`) produced two success `result` events — an early one saying only "All four ensigns are dispatched..." and a later one that actually names the corrected entity and presents the gate. The extractor returned the early one, so the assertion checking the later behavior failed. The prior green opus run on the same code also had two success results; it passed only because its first one happened to mention the asserted phrase. So the oracle's pass/fail depends on the incidental wording of an intermediate result, not the run's actual final behavior — a flake independent of real FO behavior.

Confirmed NOT caused by, and unrelated to, PR #483's own commit `43396704` (`git show --stat` confirms it touches only `codex_liveenv.go`/`codex_liveenv_test.go`). This is pre-existing shared test-harness infrastructure.

## Proposed approach (retain-last-success-result)

1. Change the extractor to retain the LAST non-error success `result` event seen while scanning the full stream, instead of returning on the first match. Track the last seen `result` string and a `haveResult` flag; `continue` past a success result rather than returning. After the loop, prefer the retained result over the assistant-text fallback.
2. Keep the loud-failure behavior for `is_error` / `api_error_status` result events UNCHANGED and short-circuiting: an `is_error` event returns `errClaudeLaunchFailed` immediately, regardless of position. A later event must never mask an earlier launch failure, and a launch failure that follows a success must still fail loudly — never treated as a fallback.
3. The existing precedence order is otherwise unchanged: result event(s) preferred over assistant text; last assistant text block as fallback when no result event exists; error when neither a result event nor assistant text is present.

## Out of scope

- The live-runner scenario assertions and their wording (`TestLiveClaudeSharedScenarios`) — this fix corrects the extractor the assertions consume, not the assertions.
- The Codex extractor / `--output-last-message` path — unaffected.
- The `streamEntry` watcher, which does not parse the result event.
- Any live-workflow test — this is pure parsing logic, fully covered offline.
- Documentation: none required. This is internal live-test-harness code with no user-visible surface (no CLI output, command surface, banner, or docs-site content changes), so no doc diff is proposed.

## Acceptance criteria

- **AC-1 (value against baseline).** Given a stream-json transcript with ≥2 non-error success `result` events, `extractClaudeFinalMessage` returns the LAST event's `result` field byte-for-byte. The baseline that moves the wrong way: current HEAD returns the FIRST event's `result`. Verified by an offline fixture subtest that FAILS red on current HEAD (returns the early "All four ensigns are dispatched..." text) and PASSES after the fix (returns the terminal text). The red-on-HEAD / green-after transition is the measured end-value.
- **AC-2 (loud failure preserved and unmaskable).** An `is_error` result event yields an error wrapping `errClaudeLaunchFailed` even when a non-error success `result` event precedes it in the same transcript; a 401 event's error still names the `api_error_status`. Verified by a fixture subtest (success-then-is_error/401) asserting `errors.Is(err, errClaudeLaunchFailed)`, plus the existing `surfaces_401_result_as_launch_failure` and `error_result_without_401_is_still_launch_failure` subtests staying green.
- **AC-3 (no regression on existing precedence).** Single-result, assistant-text fallback, and no-result-no-text behaviors are unchanged. Verified by the five existing `TestExtractClaudeFinalMessage` subtests staying green after the fix.

## Test plan

- **Offline fixture test only.** Add subtests to `TestExtractClaudeFinalMessage` in `internal/ensigncycle/claude_final_message_test.go` (beside the existing subtests):
  - `returns_last_of_multiple_success_results` — the AC-1 regression test; transcript with two success `result` events, assert the last one's `Result` is returned. This MUST fail red against current HEAD before the fix lands.
  - `success_then_error_still_fails_loud` — the AC-2 guard; a success `result` followed by an `is_error:true, api_error_status:401` event, assert `errors.Is(err, errClaudeLaunchFailed)`.
  - Keep all five existing subtests (AC-3).
- **Cost / kind.** Pure string parsing over synthetic transcripts; spends no live credential; sub-second. No CLI behavior fixture and no live-workflow test needed — the claim is parsing logic, not runtime behavior.
- **Detached adversarial audit (validation, required).** This touches the CI-and-release-machinery high-stakes surface named in the workflow Proof policy (`docs/dev/README.md:77`). Before merge, run a read-only audit on a throwaway checkout (never the implementation worktree) that constructs an adversarial edit the new tests should catch — e.g. revert the impl to first-return, or weaken the AC-1 assertion so a first-return would still pass — and confirm the test goes red. A test that stays green under a claim-breaking edit is a hole. Record material findings as a `### Feedback Cycles` entry; "refuted nothing material" is a valid recorded outcome.

## Spike result (riskiest mechanism first)

Riskiest unverified path: whether a fixture stream with multiple success `result` events actually reproduces the wrong-result return, and whether retain-last preserves the `is_error` loud-failure. Both were exercised as throwaway tests against the real `extractClaudeFinalMessage` (no code committed to main; tree restored after):

- **Red on HEAD (bug reproduced).** A throwaway `TestSpikeMultipleSuccessResultsReturnsLast` with two success results failed against current code:
  `final message = "All four ensigns are dispatched...", want the LAST success result "Gate review: names the corrected entity and presents the gate."` — HEAD returns the first, confirming the diagnosis.
- **Green after fix sketch.** Temporarily patching the impl to retain-last + short-circuit on `is_error`, then running the spikes and all five existing subtests: the multi-success spike passed, a `success_then_error` spike (success followed by 401 `is_error`) still failed loud via `errClaudeLaunchFailed`, and all five existing subtests stayed green. The impl patch and spike files were then reverted; the working tree is clean and no code lands in the ideation stage.

## Stage Report: ideation

- DONE: Problem statement, proposed approach (retain-last-success-result fix), and out-of-scope boundary written into the entity body, citing claude_final_message_impl_test.go:55-89 and the doc-comment/implementation contradiction
  Problem statement cites line 72's early return and doc-comment lines 44-51 ("terminal ... preferred") as the contradiction; proposed approach and out-of-scope sections written.
- DONE: Acceptance criteria (each with how it's verified) plus a test plan added, covering the multi-result regression test and unchanged is_error/401 loud-failure behavior
  AC-1 (last-of-multiple, red-on-HEAD/green-after), AC-2 (is_error unmaskable + 401 named), AC-3 (no regression on 5 existing subtests); test plan names offline fixture subtests and the detached adversarial audit for the CI/release high-stakes surface.
- DONE: Riskiest mechanism (a fixture stream with multiple success result events) spiked/written first so it fails red against current code before the fix lands
  Throwaway `TestSpikeMultipleSuccessResultsReturnsLast` failed red on HEAD (got "All four ensigns are dispatched...", want terminal text); a reverted impl patch proved retain-last green + success-then-401 still fails loud + 5 existing subtests green. Tree restored clean.

### Summary

Formalized the already-diagnosed extractor bug into problem statement / proposed approach / out-of-scope / acceptance criteria / test plan. The riskiest mechanism was spiked first: a two-success-result fixture fails red against current HEAD (returns the first result), and a reverted fix sketch confirmed retain-last returns the terminal result while the `is_error`/401 loud-failure stays short-circuiting and unmaskable. Pure offline parsing fix; no doc changes; validation must run the detached adversarial audit per the CI/release-machinery Proof policy.
