---
title: Live opus FO intermittently routes the cycle-2 rejection re-review to the impl worker (fix≠reviewer violation) — rejection-flow flake
status: backlog
sprint: 0230-stable-finalization
group: fo-reliability
id: 7hczkc0c6ezgwy1p627ejp6x
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
