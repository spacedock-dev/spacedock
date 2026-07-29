---
title: "Bind the post-rework Briefing when rejection flow re-enters the gate"
status: backlog
source: "se0 complete Opus live suite, 2026-07-28: after validation/1 rejection, successful advisory round recording, implementation rework, and cycle-2 PASS, the FO rebound the original round-1 Briefing as the ordinary final gate instead of the distinct post-rework Briefing."
score: 0.85
sprint: durable-decisions
group: gate
sprint-readiness: ready
issue:
id: zbcj98qfwtax61vxdzrf615e
---

## Problem

A supported feedback-rejection journey can correctly retain the validation/1 advisory review round, complete rework and a second validation, then bind the stale pre-rework Briefing as the ordinary final gate. In the retained complete Opus run, `gate record --round validation/1` succeeded, cycle 2 passed, and the FO selected `rejection-task/inputs/briefing.json` again. PR run `30412397240` reproduced the same defect on Sonnet: two validations completed, but the final gate bound `briefing:rejection-task:validation:round-1`, so `assertReviewRoundRecorded` correctly rejected the advisory round being retained as a gate attempt.

A prior focused Opus run on the same candidate correctly created and bound a distinct `gate-validation-round2/briefing.json`. The behavior is therefore nondeterministic contract conduct, not a recorder failure or an oracle false positive. After feedback changes the candidate, final approval must spend a freshly bound post-rework Briefing; an advisory review-round package cannot silently become that gate.

The Opus and Sonnet executions of the shared rejection-flow journey are temporarily TODO under this task. Keep Codex, Pi, recorded-gate, keep-moving, and deterministic round/gate coverage active. Re-enable both Claude cases when this task lands.

## Acceptance criteria

**AC-1 (VALUE) - A captain deciding after rejected feedback always reviews the post-rework candidate, never the stale rejected snapshot.**
Verified by: repeated clean Opus and Sonnet rejection-flow journeys record validation/1 as advisory, create a distinct post-rework canonical Briefing, and bind only that new Briefing to the ordinary final gate.

**AC-2 - Advisory round and binding gate roles remain mechanically distinct.**
Verified by: the retained round record preserves the original reviewer/worker disposition, while gate state identifies a separate Briefing ID/digest produced after rework and cycle-2 validation.

**AC-3 - The remedy is agent-ergonomic and provider-neutral.**
Verified by: the FO can select or prepare the post-rework Briefing through the normal gate lifecycle without reconstructing metadata, depending on Opus wording, or adding compatibility behavior.

**AC-4 - The quarantined Claude live journeys are restored.**
Verified by: remove the linked TODO for both models, run the focused Opus and Sonnet rejection-flow journeys repeatedly, then run both complete Claude shared suites at the exact candidate tip.

## Boundary

Do not weaken the semantic oracle and do not allow a review-round Briefing to double as a later binding gate. Determine whether the missing authority belongs in feedback-rejection routing, gate lifecycle preparation, or the shipped briefing-selection contract, then fix the smallest authoritative seam. V1 is unreleased; add no compatibility layer.
