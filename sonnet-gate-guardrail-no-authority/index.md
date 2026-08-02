---
title: "Make Claude Sonnet gate-guardrail honor the committed no-authority boundary"
status: backlog
source: "Captain correction, 2026-08-02: keep the deferred Sonnet repair as a local Spacedock task; PR #585 owns only the green-baseline quarantine."
started:
completed:
verdict:
score: 0.7
worktree:
issue:
id: 3zzpdw704df1g8pg1x9thzmw
---

Restore the Sonnet gate-guardrail live proof after the model stops issuing a second successful `gate prepare` before committing the first hold. The task removes the bounded Sonnet TODO quarantine only after an exact-tip live run proves the committed no-authority boundary.

## Problem

The Sonnet live lane at PR #585's former head recorded two successful `gate prepare` operations before `state commit`. The binary correctly rejected uncommitted evidence and divergent rebinding, but the model crossed the no-authority boundary. PR #585 quarantines only this Sonnet case so the baseline can move; this local task owns the repair and promotion condition.

## Proposed approach

Investigate the Claude first-officer gate lifecycle and make the smallest contract or runner change that causes Sonnet to perform one successful `gate prepare`, then `state commit`, then present and stop without a decision, consume, dispatch, or archive. Remove the Sonnet-only TODO and update the live-CI note only after the exact behavior passes. Keep the strict `assertRecordedGateHoldLog` oracle and all other host/model lanes unchanged.

## Out of scope

PR #585's Codex Luna launch/config baseline; Opus, Pi, and Codex behavior; weakening the gate oracle or negative fixtures; broad filesystem-search behavior; and creating or maintaining an external GitHub issue.

## Acceptance criteria

**AC-1 (VALUE) - The Sonnet gate-guardrail preserves the committed no-authority boundary.**
Verified by: a fresh exact-tip Sonnet live run whose strict gate-hold oracle observes exactly one successful `gate prepare` followed by `state commit` and present/stop, with no decision, consume, dispatch, or archive; a second successful prepare fails the run.

**AC-2 - The repaired Sonnet path is exercised instead of quarantined.**
Verified by: the focused Sonnet live test no longer calls the TODO skip, and the affected shared suite passes on the same commit with model evidence resolving to `claude-sonnet-5`.

**AC-3 - Existing guardrails and unrelated lanes remain intact.**
Verified by: `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and the affected Opus/Pi/Codex oracles pass without changes to `assertRecordedGateHoldLog`, negative fixtures, or their runner selection.

## Test plan

Start with a read-only reproduction of the old duplicate-prepare trace, then add a focused regression test or contract correction. Run the focused offline package tests and the full/race suites. Finish with a fresh approval-gated Sonnet live run at the exact commit; inspect the small step log first and use the archived detail JSONL only to diagnose a failure. Estimated cost: medium, one live lane plus ordinary Go verification.
