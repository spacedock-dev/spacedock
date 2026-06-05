---
id: 95517mg61x8c8fmqz4htchp7
title: Behavioral drive for the pre-completion-TeamDelete ban (no-drive claim surfaced by hwk)
status: backlog
source: "hwk (tautological-test-remediation) implementation 2026-06-05 — TestAwaitingCompletionStillBansPreCompletionTeamDelete asserts the FO must NOT emit TeamDelete before a worker's completion signal (the premature-teardown bug). It is a no-drive BEHAVIORAL claim: exercised only IMPLICITLY by every live team scenario (a premature teardown breaks the run), with no dedicated mutation-controlled assertion. hwk demoted it honestly and named this task as the owed oracle. Distinct from the terminal-teardown HANG, which #285 + TestSonnetTeamDeleteHangReplay already cover."
score: "0.28"
started:
completed:
verdict:
worktree:
issue:
---

The `using-claude-team` `## Awaiting Completion` discipline bans the FO from emitting `TeamDelete` (or a `shutdown_request`) before a dispatched worker's completion signal — the premature-teardown bug. hwk's standing sweep flagged `TestAwaitingCompletionStillBansPreCompletionTeamDelete` as a presence check over the contract prose; hwk demoted it (markNonAC) and named THIS task as the owed behavioral oracle. This is the only no-drive behavioral claim beyond the halt/sync/journey cluster (which ev3e owns).

## Problem

The ban is a BEHAVIORAL claim ("the FO does not tear down before completion"). Today it has:
- the (now-demoted) presence check over the contract — proves nothing;
- IMPLICIT coverage — every live team scenario would break if the FO prematurely tore down — but no DEDICATED, mutation-controlled assertion that reds specifically on a premature `TeamDelete`/`shutdown_request`.

Per the proof-policy, the behavior needs a drive that runs it and reds on the violation.

## Proposed approach (ideation formalizes)

A behavioral drive that places the FO mid-dispatch (a worker spawned, completion signal NOT yet observed) and asserts the FO HOLDS — does NOT emit `TeamDelete`/`shutdown_request` before the completion signal — graded on durable/observed state (no teardown event before the completion). Plus an offline negative that reds on a premature-teardown end-state.

Host scope: `TeamDelete` is Claude-team-specific (Codex/Pi have no equivalent), so this is likely a Claude-specific live drive rather than a host-parity shared scenario — assess in ideation (a shared cross-host scenario may not fit; a Claude-only live test or a host-neutral "premature-teardown" abstraction). The riskiest unknown: grading the ABSENCE of a premature action (similar to ev3e's halt — grade that nothing tore down before the completion signal), so spike the grading first.

## Acceptance criteria (ideation formalizes; per the proof-policy the proof is RUNNING the behavior)

**AC-1 (behavioral) — a drive reds when the FO emits TeamDelete/shutdown_request before a worker's completion signal, and passes when it holds.** Verified by: a cited live run (or the closest faithful drive) + an offline negative, mutation-controlled.
**AC-2 (gap closure) — the hwk demotion of TestAwaitingCompletionStillBansPreCompletionTeamDelete names this drive as its owed oracle, and the drive exists.**

## Test plan

- Spike: grade-the-absence-of-a-premature-action (offline), like ev3e's halt grading. Then the drive + offline negative (live-gated if a real team drive is needed). High-stakes (team-harness machinery) → detached audit.
- Pairs with hwk (the demotion names this) + ev3e (sibling no-drive-behavioral-claim drive, split-root scope) + the proof-policy (`f8b257cf`).

## Notes

Provenance: hwk implementation surfaced this as the only no-drive behavioral claim beyond halt/sync/journey; the captain's standing instruction was to file a follow-up rather than bloat hwk. Sibling pattern: ev3e (the halt/sync/journey live drive).
