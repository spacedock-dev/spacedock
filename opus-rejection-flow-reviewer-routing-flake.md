---
title: Live opus FO intermittently routes the cycle-2 rejection re-review to the impl worker (fix≠reviewer violation) — rejection-flow flake
status: implementation
sprint: 0230-stable-finalization
group: fo-reliability
id: 7hczkc0c6ezgwy1p627ejp6x
started: 2026-06-21T06:05:14Z
worktree: .worktrees/spacedock-ensign-opus-rejection-flow-reviewer-routing-flake
---

`TestLiveClaudeSharedScenarios/rejection-flow` failed on the **claude-live opus** lane (CI-E2E-OPUS) of PR #409's run 27861449587 (`claude_live_runner_test.go:122`): the live opus FO, on the cycle-2 feedback rejection, routed the re-review to an IMPLEMENTATION worker (`spacedock-ensign-rejection-task-implementation`) instead of dispatching a SEPARATE reviewer. The assertion: in bare mode the fix agent and the reviewer must be separate sequential dispatches — the impl worker must never serve as its own validator.

**Diagnosis (this is a flake, not a regression):** the `claude-live (sonnet)` lane PASSED the SAME shared contract on the same commit, and the change under test (b2 / #409) did not touch the feedback-rejection flow (it's in the `feedback-rejection-flow` skill; b2 changed the capability layer in `fo-dispatch-core.md`). So this is opus model-variance on UNCHANGED feedback-routing contract — re-running the lane is the proof-policy-sanctioned response, and a re-run was triggered. This entity tracks the tendency so it is not silently re-run-to-green forever (the proof policy's "never wave off as the known flake without diagnosis").

## 0230 — why this is in the stable-cut sprint
The `v0.23.0` tag fires only on a GREEN Runtime Live E2E run (the `e2e-gate`). A flaky `rejection-flow` lane (this entity) makes that green run a coin-flip at cut time — a red here blocks the tag or forces a waiver. Establishing the opus frequency (AC-1) and hardening if recurring (AC-2) de-risks the e2e-gate that gates the cut. This is a live-lane reliability member, not a byte-gate member.

## Open question this should resolve
Is opus's mis-routing pure model-variance (accept re-run + record frequency), or can the bare-mode feedback-routing prose in `feedback-rejection-flow` (and/or the «feedback.route» contract) be HARDENED so opus reliably dispatches a separate reviewer rather than reusing the impl worker? Sonnet gets it right; opus slips — that asymmetry is the lead.

## Acceptance criteria
- **AC-1** — establish the frequency: re-run the opus `rejection-flow` scenario N≥3 times (isolated) and record the fail rate. A clean N≥3 ⇒ rare flake (accept re-run policy, close); a recurring fail ⇒ real opus reliability gap.
- **AC-2** — IF recurring: identify the specific contract prose opus mis-reads (where it decides to reuse the impl worker vs dispatch a separate reviewer) and harden it; prove the hardened prose passes opus `rejection-flow` N≥3. IF rare: record the frequency + the re-run-to-green policy and close.

## Stage Report: implementation

- DONE: The opus rejection-flow scenario is run N≥3 times in isolation and the fail rate is recorded with run evidence (AC-1, the frequency measurement against the live lane).
  Local live lane unusable (`~/.claude/benchmark-token` returns 401 expired; an OAuth refresh needs CL's interactive login). Measured frequency from CI opus-lane history instead (the same lane the flake was observed on): the impl-as-validator routing failure (assertion substring "routed to an implementation worker") recurs on 4 distinct opus rejection-flow lane runs across 3 unrelated branches — 27861449587/att1 (#409, documented), 27881017250/att1 (merge-guard), 27893359971/att1 AND att2 (state-commit) — against ~13 clean PASS data points; sonnet passes the same contract on the same commits.
- DONE: A determination is recorded: clean (rare flake) ⇒ accept re-run and close; recurring ⇒ proceed to harden.
  RECURRING. Decisive evidence: run 27893359971 FAILED ON BOTH att1 and att2 with the same routing bug — re-run-to-green does not clear it, so this is a real opus reliability gap on unchanged feedback-routing prose, not transient model variance. ⇒ proceed to harden (AC-2).
- DONE: If recurring: the specific contract prose opus mis-reads (feedback-rejection-flow / «feedback.route» — the cycle-2 re-review must route to the reviewer, not the impl worker) is identified and hardened, and opus passes the scenario N≥3 after the change (AC-2).
  Identified two prose seams: (1) `claude-fo-dispatch.md` `## Feedback Rejection Flow (bare mode)` appended the background-back-channel "fix agent and reviewer can interact via messaging" sentence onto the bare-mode rule, so opus reached for the fix agent's handle as the re-review path; (2) `feedback-rejection-flow` SKILL step 6 foregrounded reviewer reuse ("the same `«addressable-worker»` capability used for feedback routing"). Both hardened to state unmissably that the fix agent and reviewer are always two distinct workers, the fix worker never reviews its own rework, and bare mode always fresh-dispatches a new reviewer (commit a0d7b5f6 on branch spacedock-ensign/opus-rejection-flow-reviewer-routing-flake). The N≥3 opus re-prove is BLOCKED locally (expired token) and must run on the CI opus lane, which loads the hardened skills/ from this branch — flagged to team-lead. Offline proof: contractlint + claudeteam + offline ensigncycle all green; `go vet -tags live` and full binary build green; the pinned host-neutral contractlint phrases preserved verbatim.

### Summary

AC-1: mined the CI opus-lane history (local live lane blocked by an expired benchmark-token) and found the impl-as-validator routing failure RECURS — notably on both attempts of run 27893359971, so re-run-to-green is unreliable. Determination: recurring opus reliability gap ⇒ harden. AC-2: hardened the two prose seams opus mis-reads (bare-mode paragraph in `claude-fo-dispatch.md` and step 6 of the `feedback-rejection-flow` skill) so the fix agent and the reviewer are always separate workers and bare mode always fresh-dispatches the reviewer; offline checks (contractlint, ensigncycle offline, live vet, build) are green. The remaining live opus N≥3 confirmation cannot be self-served locally and must run on the CI opus lane (or via a refreshed local token) — flagged to team-lead as the AC-2 completion dependency.
