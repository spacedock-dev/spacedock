---
title: "Bind the post-rework Briefing when rejection flow re-enters the gate"
status: ideation
source: "se0 complete Opus live suite, 2026-07-28: after validation/1 rejection, successful advisory round recording, implementation rework, and cycle-2 PASS, the FO rebound the original round-1 Briefing as the ordinary final gate instead of the distinct post-rework Briefing."
score: 0.85
sprint: durable-decisions
group: gate
sprint-readiness: ready
issue:
id: zbcj98qfwtax61vxdzrf615e
gates:
    version: 1
    current:
        gate: gate:zbcj98qfwtax61vxdzrf615e:backlog
    records:
        - id: gate:zbcj98qfwtax61vxdzrf615e:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zbcj98qfwtax61vxdzrf615e-backlog-1
              briefing:
                id: briefing:zbcj98qfwtax61vxdzrf615e:backlog:attempt-1:revision-1
                digest: sha256:a8668228e65695fdea30226ee877edb1031da0356a36cca5b245d644c3434802
                digest-domain: canonical-bytes
                request-digest: sha256:70a3922bccbd3031f2b3e4b7f5921d1b081db5861350376d92d8f6be23b6cc35
                room-ref: ./bind-post-rework-briefing-at-rejection-regate/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zbcj98qfwtax61vxdzrf615e:backlog:1
                briefing: briefing:zbcj98qfwtax61vxdzrf615e:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T15:39:42.727381Z"
                decision: approve
                reason: 'Sprint conn: task protects decision integrity after rework, reuses existing machinery, and can be ideated independently of other critical lanes.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

## Problem

A supported feedback-rejection journey can correctly retain the validation/1 advisory review round, complete rework and a second validation, then bind the stale pre-rework Briefing as the ordinary final gate. In the retained complete Opus run, `gate record --round validation/1` succeeded, cycle 2 passed, and the FO selected `rejection-task/inputs/briefing.json` again. PR run `30412397240` reproduced the same defect on Sonnet: two validations completed, but the final gate bound `briefing:rejection-task:validation:round-1`, so `assertReviewRoundRecorded` correctly rejected the advisory round being retained as a gate attempt.

A prior focused Opus run on the same candidate correctly created and bound a distinct `gate-validation-round2/briefing.json`. The behavior is therefore nondeterministic contract conduct, not a recorder failure or an oracle false positive. After feedback changes the candidate, final approval must spend a freshly bound post-rework Briefing; an advisory review-round package cannot silently become that gate.

The Opus and Sonnet executions of the shared rejection-flow journey are temporarily TODO under this task. Keep Codex, Pi, recorded-gate, keep-moving, and deterministic round/gate coverage active. Re-enable both Claude cases when this task lands.

## Acceptance criteria

**AC-1 (VALUE) - Final approval is always spent on the reworked and revalidated candidate, never the rejected snapshot.**
Verified by: repeated clean Opus and Sonnet rejection-flow journeys record validation/1 as advisory, create a distinct post-rework canonical Briefing, and prove the ordinary final approval resolves the newest post-rework episode rather than any stale candidate, even if another path satisfies the literal Briefing-bind sequence.

**AC-2 - Advisory round and binding gate roles remain mechanically distinct.**
Verified by: the retained round record preserves the original reviewer/worker disposition, while gate state identifies a separate Briefing ID/digest produced after rework and cycle-2 validation.

**AC-3 - The remedy is agent-ergonomic and provider-neutral.**
Verified by: the normal gate lifecycle prepares and binds the post-rework Briefing without the FO reconstructing metadata, depending on model wording, or adding compatibility behavior. Gate preparation or binding mechanically computes or verifies freshness against the newest recorded post-rework episode and refuses a stale bind with an actionable diagnostic.

**AC-4 - The quarantined Claude live journeys are restored.**
Verified by: remove the linked TODO for both models, run the focused Opus and Sonnet rejection-flow journeys repeatedly, then run both complete Claude shared suites at the exact candidate tip.

## Boundary

Do not weaken the semantic oracle and do not allow a review-round Briefing to double as a later binding gate. The observed nondeterminism is evidence against a prose-only correction: identical contract text admitted both correct and stale binding, so ideation must price the prose-only variant and reject it explicitly unless a falsifiable exercise disproves this finding.

Ideation must define **the newest post-rework episode** once and choose its authoritative key: the correction-cycle record, the rework commit, or a new epoch marker. That definition must serve all consumers that need it, including Feedback Cycles legibility, Briefing freshness, and gate binding. Prefer a mechanical guard in `gate prepare` or the binding seam that computes or asserts that the selected Briefing post-dates the latest recorded rework and fails closed on a stale candidate.

Determine whether the missing authority belongs in feedback-rejection routing, gate lifecycle preparation, or the shipped briefing-selection contract, then fix the smallest authoritative seam. V1 is unreleased; add no compatibility layer.
