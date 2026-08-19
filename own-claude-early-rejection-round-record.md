---
title: Early rejection-round record grades as never-recorded
status: backlog
source: "Run 31996696789 claude lane, rejection-flow: round recorded at stream line 287 with entries=2, rework spawned line 309; captain approved filing owners for tolerated residual modes at the 0.27 composite-green ruling, 2026-08-17"
id: zf7rymtke3b6xp7r0337hjj4
gates:
    version: 1
    records:
        - id: gate:zf7rymtke3b6xp7r0337hjj4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zf7rymtke3b6xp7r0337hjj4-backlog-1
              briefing:
                id: briefing:zf7rymtke3b6xp7r0337hjj4:backlog:attempt-1:revision-1
                digest: sha256:67e727f265794634086d12c6cd19fc36f1ced497cb6cf325905a606cf4aa0ebc
                request-digest: sha256:a95473793cb5a3ac73d64382a04502e43a92229e3c2112dc29b626c91951f847
                room-ref: ./own-claude-early-rejection-round-record/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zf7rymtke3b6xp7r0337hjj4:backlog:1
                briefing: briefing:zf7rymtke3b6xp7r0337hjj4:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-19T18:16:02.74664Z"
                decision: approve
                reason: 'Captain approved in chat: ''if the fable opinion lands reasonable, dispatch those two on to the stack of 736 so we can run tip CI lane to verify.'' The fable review landed and confirmed the diagnosis while finding two real remedy gaps, both now recorded.'
              application:
                target-stage: ideation
                state: pending
---

## Problem

Two coupled defects around one residual mode. The mode: the live claude FO invoked `gate record --round validation/1` immediately after the rejection — before routing the correction — so the immutable round room durably holds the reviewer's 2 entries instead of the complete 4 the contract's record-after-correction order produces. The red is correct in substance. The label is not: `claudeRecordedRejectionRound` accepts only a success line pinned to `entries=4` (`rejectionRoundSuccess`, shared_round_recording_test.go), so a successful early record is reported as "resolved launcher never invoked `gate record --round validation/1`" — the lying-label class this release spent itself killing. A second label defect in the same oracle: the workflow-README immutability boundary condition reports under `rejection-round-missing` ("README.md changed from its exact expected bytes", observed run 31991864922 attempt 2), which reads as a recording failure when the round WAS recorded.

## Recurrence confirmed (2026-08-19, FO)

The mode recurred in PR #736's live lane, run 32270990171, claude-live. Same shape, and the retained stream settles it beyond the earlier inference:

    tool_use:    ${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 \
                   --briefing .../briefing.json --log .../briefing.review.jsonl --workflow-dir ...
    tool_result: is_error = false
                 round=round:rejection-task:validation:1 stage=validation cycle=1
                 briefing=briefing:rejection-task:validation:round-1 entries=2
                 entry=annotation:rejection-task:missing-marker  type=Annotation
                 entry=resolution:rejection-task:reviewer        type=Resolution decision=revise

`entries=2`, exactly as this entity predicted: the reviewer's two entries, recorded before the correction was routed. `rejectionRoundSuccess` pins `entries=4`, so the oracle graded it "resolved launcher never invoked" — a successful, exit-0, correctly-formed invocation reported as absent.

One diagnosis to strike from the record, so it does not mislead the next reader: the FO initially attributed this red to `invokesRejectionRoundRecorder` failing to parse the contract-mandated `${SPACEDOCK_BIN:-spacedock}` form. That is WRONG. `directRoundLauncher` (`shared_round_recording_test.go:17`) accepts the default-expansion form explicitly via `(?::-[^}]*)?`, and it is checked first. The launcher parses fine; the entry count is the whole cause. The FO read the fallback branch's inner regex and never tested the branch that returns before it.

Evidence preserved at `/tmp/9g-evidence` (both hosts' full artifacts, 30 MB) independent of CI artifact expiry.

Sibling reds in the same journey, same release: `rejection-topology-count-bar` (`12z`, filed today) covers the codex-side `rejection-worker-topology` exact-count bar. Both are the same family as the merged `filing-recognizer-newline-terminator` (#731) and `decide-dispatch-build-count-bar` (#732) — a pattern written against one observed shape, then treated as the definition of correct conduct.

## Proposed approach

Split honesty from strictness in the oracle. Match the recorder's success line generically (`entries=(\d+)`) and grade the count as its own condition with an honest code (`rejection-round-incomplete`, naming got-vs-want and that the record preceded the correction); keep `rejection-round-missing` for genuinely absent invocations. Move the README-immutability condition to its own code (`rejection-workflow-doc-mutated`). Both proven by falsifying edits: replay the run-31996696789 stream bytes through the recognizer (early record must grade incomplete, not missing) and a README-mutated fixture (must grade doc-mutated, not round-missing). Skill-side, the record-after-correction order is already explicit in step 6; the owner tracks the mode's recurrence rate with the metrics instrument rather than adding prose the model already had.

## Adversarial review of this remedy (2026-08-19, fable ensign, read-only)

Diagnosis CONFIRMED against the live stream: one `--round` call at line 286, result at 287 carrying `entries=2`, `is_error=false`; the rework spawn is at index 307, so the record demonstrably preceded the correction. Every link holds.

It also settled the question the FO could not: is an honest label a way to make a real defect easier to ignore? No. The conduct IS wrong — the round room is immutable, so an early record permanently truncates the durable record to 2 of 4 entries. Under this remedy the red STAYS red; only its name changes. The opposite risk is already proven: the lying label sent the FO chasing a launcher-regex ghost. And it matches the codebase's own pattern (`assertRejectionCycleLine`'s heading-drift diagnostic, `TestRejectionUnpreparedGateReportsItsOwnCode`): the grade does not soften, the diagnostic must say what happened.

**But the remedy as written SHIPS HALF-DONE. Two gaps, both verified by the FO:**

1. **The `entries=4` pin lives in TWO places, not one.** The stream regex at `shared_round_recording_test.go:18`, and a durable summary check at `:400-405` — `len(summary.Entries) != 4`, message `"retained round summary = %#v"`, still under `rejection-round-missing`. Genericize only the stream regex and the same run reds at :404 under the same lying code with a WORSE message: a raw struct dump. The count grade must move to the DURABLE site under the honest code. And this remedy's proposed proof — replay the stream bytes through the recognizer — exercises only the stream path, so it would GREEN the half-fix. A durable fixture is required: record a 2-entry log, assert the code is `rejection-round-incomplete`.
2. **The byte-exact block at `:425-442` covers FOUR files**, not one: `briefing.json`, `briefing.review.jsonl`, `candidate.txt`, `README.md` — all reporting "changed from its exact expected bytes" under `rejection-round-missing`. This remedy splits only README. `candidate.txt` carries the identical lying label. The split criterion should be **"is this file part of the round record"** — the two room files keep the round code, everything else does not — rather than "is it README".

**Related conduct hotspot, named by the review and owned by neither entity.** Both hosts deviated at the SAME point in the flow: the gap between rejection and routed correction. Claude recorded the round early; codex inserted an extra review round. If recurrence continues, the fix is mechanical rather than prose — for example a precondition the FO can check before `gate record --round` (the log tail contains a correction entry). This entity's "metrics over prose" call stands for now, at two occurrences.

## Out of scope

- The round recorder binary (neutral; its output already carries the true entry count).
- Rejection-flow fixture text (the Cycle-line target pin shipped with the 0.27 stack).
