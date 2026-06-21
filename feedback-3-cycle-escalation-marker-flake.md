---
title: feedback-3-cycle-escalation lane flaky — live FO escalates correctly on cycle 3 but omits the exact human-escalation marker
status: backlog
sprint: 0230-stable-finalization
score: 0.4
group: cleanup
issue:
id: v56dg1amaa33zgfrg6r0q3bn
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
