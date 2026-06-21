---
title: feedback-3-cycle-escalation lane flaky — live FO escalates correctly on cycle 3 but omits the exact human-escalation marker
status: validation
sprint: 0230-stable-finalization
score: 0.4
group: cleanup
issue:
id: v56dg1amaa33zgfrg6r0q3bn
worktree: .worktrees/spacedock-ensign-feedback-3-cycle-escalation-marker-flake
started: 2026-06-21T21:14:12Z
mod-block: merge:pr-merge
pr: pr-merge:427
---

The `feedback-3-cycle-escalation` shared scenario intermittently fails the live e2e-gate (observed on SONNET, CI-E2E, in the M5 and M6 PR runs 27913790926 / 27913761338): `claude_live_runner_test.go:122` — "escalation entity did not record the human-escalation marker in the `### Feedback Cycles` section on the 3rd cycle." The captured transcript shows the FO ESCALATES correctly on cycle 3 (it stops re-routing and hands to the human) but writes prose ("escalated to the human") instead of the fixture's exact marker string `feedback-escalation: human-review-required`. This is a second e2e-gate reliability risk alongside M7's rejection-flow false-positive — the v0.23.0 tag fires only on a green Runtime Live E2E.

## Problem

The fixture greps for an exact marker string. The live FO performs the correct BEHAVIOR (escalate on cycle 3, stop the re-route loop) but does not deterministically emit that exact token. Two candidate root causes, to be distinguished by evidence (mirror M7's discipline — do not assume):
1. **Over-strict fixture / test false-positive** (like M7): the assertion demands a literal string the contract never pins, so a correct escalation that words it differently fails. If so, the fix is the assertion/fixture (assert the behavior, not the spelling), keeping a genuine non-escalation RED.
2. **Real contract/skill gap**: the `feedback-rejection-flow` skill / escalation contract does not unambiguously instruct the FO to write the exact `feedback-escalation: human-review-required` marker on cycle 3, so the model omits it. If so, harden the prose so the marker is emitted deterministically, and code-gate it.

## Proposed approach

1. **Establish frequency (N≥3).** Re-run the `feedback-3-cycle-escalation` scenario isolated on sonnet (and opus) N≥3; record the fail rate and grade each transcript by hand for the ACTUAL behavior (did the FO escalate correctly?) vs the marker string. Distinguish flake-rate from a deterministic gap.
2. **Root-cause to candidate 1 or 2** with transcript + fixture + contract evidence, the way M7 did (download/grade the real failing stream; read `claude_live_runner_test.go` line 122 and the escalation fixture; read the `feedback-rejection-flow` skill's escalation step).
3. **Fix the right target.** If over-strict fixture → make the assertion behavior-aware (accept a correct escalation regardless of exact wording; keep genuine non-escalation RED). If contract gap → harden the skill/contract prose to pin the marker AND add a code-gate (contractlint or the live assertion) so it cannot silently regress. Do NOT touch M1/M2 contract surfaces beyond what the escalation path requires.

## Acceptance criteria

- **AC-1 (frequency)** — the scenario is run N≥3 isolated on sonnet (and opus where relevant); the fail rate is recorded with per-run transcript grading of the actual escalation behavior, and a determination (flake vs deterministic gap) is on the record.
- **AC-2 (root cause)** — the failure is attributed to candidate 1 (over-strict fixture) or candidate 2 (contract gap) with cited evidence (the failing transcript, the line-122 assertion, the escalation fixture, the feedback-rejection-flow skill text), not asserted.
- **AC-3 (fix + proof)** — the correct target is fixed and the scenario passes N≥3 after the change on the affected model(s); a genuine non-escalation (FO fails to escalate on cycle 3) still FAILS the assertion (non-tautological); offline `go test ./internal/ensigncycle ./internal/contractlint` + `go build ./...` green.
- **AC-4 (e2e-gate de-risk)** — the change lands before the v0.23.0 cut so the tag's required green Runtime Live E2E is not a coin-flip on this scenario.

## Test plan

Live `feedback-3-cycle-escalation` runs (rotated token; surface token-stale to the FO, never count as a scenario fail). In-process grade of captured streams through the live-runner assertion entry point (the M7-proven stronger-than-fixture method). Offline contractlint/ensigncycle/build as the cheap gate. No new eval framework (YAGNI).

## Stage Report: implementation

- DONE: Establish frequency N≥3; record fail rate; grade EACH transcript by hand for ACTUAL behavior vs the exact marker; distinguish stochastic flake from deterministic gap (AC-1).
  5 transcripts graded (1 CI + 4 fresh local sonnet). In EVERY one the FO escalated CORRECTLY on cycle 3 — recorded the third `### Feedback Cycles` entry, parked the entity at `validation`, dispatched no fourth round, did not terminalize. The exact-marker emission is STOCHASTIC: CI run 27913790926 (the documented failure, sonnet) wrote prose "Escalated to human" in the body and NEVER read the README body that pins the token (the marker string is absent from the entire captured stream), so it red on the old exact-token check; local runs 1-4 all happened to read the body and emit the marker. Determination: a stochastic flake on the exact-token assertion, NOT a deterministic behavioral gap (candidate 1).
- DONE: Root-cause to candidate 1 (over-strict fixture/test false-positive) OR candidate 2 (real contract/skill gap) with cited evidence (AC-2).
  Candidate 1 — test false-positive. Evidence: (a) the failing transcript (CI run 27913790926 sonnet artifact 7778928478) — FO escalated correctly, wrote `- Cycle 3: REJECTED … Escalated to human.`, exact marker absent from the whole stream; (b) the assertion `internal/ensigncycle/shared_assertions_impl_test.go` (the line-122 entry point `claude_live_runner_test.go:258`→`assertThirdCycleEscalation`) demanded `strings.Contains(section, "feedback-escalation: human-review-required")`; (c) the fixture pins that token ONLY in `escalationReadme()`'s `### Feedback Cycles` BODY prose; (d) the FO contract `skills/first-officer/references/first-officer-shared-core.md:16` tells the FO to DEFER the README body and read only the structural index, and `feedback-rejection-flow/SKILL.md:17` says only "On cycle 3, escalate to the human" — it never pins a marker. So the token is not a contract obligation; emitting it is non-deterministic.
- DONE: Fix the correct target and prove it: scenario passes N≥3 after the change; a genuine non-escalation STILL FAILS (non-tautological); offline `go test ./internal/ensigncycle ./internal/contractlint` + `go build ./...` green (AC-3).
  Fixed the assertion (worktree commit 09abdd7d, TEST-ONLY — no contract/skill/fixture edits): `assertThirdCycleEscalation` now accepts the offered marker OR the FO's own escalate-to-human wording (`escalationToHuman` regex), scoped to the `### Feedback Cycles` section. Non-tautological: added a negative case (three cycles recorded, NO handoff → still RED on the handoff check) and the CI-observed-body regression test; the four existing failure modes (4th auto-bounce, stalled-at-cycle-2, terminalized, stray report) stay RED on their own checks. Proof: 4/4 fresh sonnet live runs PASS (96s/106s/84s/77s); offline `go test ./internal/ensigncycle ./internal/contractlint` green, `go build ./...` green, `go vet -tags live ./internal/ensigncycle/` green.

### Summary

The flake is a test false-positive (candidate 1), not an opus/sonnet behavior gap. The exact escalation marker lives only in deferred README body prose the FO contract tells the FO not to read, so a contract-faithful FO that escalates correctly may word the handoff in its own prose — which the old exact-token check red-flagged (CI run 27913790926). The fix makes `assertThirdCycleEscalation` grade the escalation-to-human BEHAVIOR (offered marker OR the FO's own wording, section-scoped) while keeping every genuine failure mode RED, proven by a new no-handoff negative case + the CI-observed-body regression and 4/4 green live sonnet runs. The change is test-only with zero entanglement with the M1/M2 contract surfaces.

## Stage Report: validation

- DONE: Verify the assertion fix grades BEHAVIOR not the exact token: assertThirdCycleEscalation accepts the offered marker OR the FO's own escalate-to-human wording (section-scoped); the CI-observed failing body now PASSES; root cause sound.
  `assertThirdCycleEscalation` (shared_assertions_impl_test.go:122) now `!Contains(escalationMarker) && !escalationToHuman.MatchString` over the `### Feedback Cycles` section. `TestThirdCycleEscalationAcceptsCIObservedHandoff` (the run-27913790926 prose body, exact marker absent) PASSES. Root cause confirmed by reading fixtures: marker pinned ONLY in escalationReadme() body prose (shared_fixtures_test.go:184); `first-officer-shared-core.md:4` DEFERS the README body; `feedback-rejection-flow/SKILL.md:3` never pins a marker; grep of skills/agents/docs finds the token nowhere.
- DONE: Verify non-tautological / discriminating: the NEW no-handoff negative case STAYS RED, and the four existing failure modes (4th auto-bounce, stalled-at-cycle-2, terminalized, stray report) each STAY RED on their own checks.
  Mutation testing (4 independent edits, each reverted): disabling the handoff check → caught at shared_scenarios_negative_test.go:197 (no-handoff) + shared_assertions_test.go:81; disabling cycle-count → caught at shared_assertions_test.go:74; disabling report-count → caught at :233 (stray report); disabling terminal → caught at :252 (terminalized). regex probe: CI-prose/natural-wording/offered-marker/"handed off to human" MATCH; no-handoff & stalled-cycle-2 do NOT. No blanket-pass.
- DONE: Confirm TEST-ONLY (zero M1/M2 entanglement), offline `go test ./internal/ensigncycle ./internal/contractlint` + `go build ./...` green, and a live sonnet run passes.
  `git diff aa49853c HEAD` = only shared_assertions_impl_test.go + shared_scenarios_negative_test.go (no production/skill/fixture/agents edits). Offline gate green uncached (ensigncycle 6.7s, contractlint 0.45s); `go build ./...` exit 0; `go vet -tags live ./internal/ensigncycle` exit 0. Live sonnet 3/3 PASS (85s/92s/69s) via in-process assertThirdCycleEscalation against the real durable entity; run-4 transcript shows the FO escalated correctly (Cycle 3 recorded, parked at validation, no 4th round, marker + own wording both emitted).

### Summary

PASSED. The fix is a correct test-only behavior-aware oracle, fully discriminating (every check independently guarded by an isolating negative, proven by mutation), and the CI-observed flake body now passes. Live sonnet is green 3/3 against the real end-state. One initial live run failed on the macOS `/var`→`/private/var` symlink artifact in `detectWrongRootBoot` (claude_live_runner_test.go:422 passes the raw unresolved `t.TempDir()` root, unlike the resolved-root team-mode callers) — a pre-existing local-only harness gap that CI's Linux env does not hit and that the M8 diff never touches; eliminated by a canonical TMPDIR. Not an M8 regression and not a scenario fail.
