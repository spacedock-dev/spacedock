---
title: Early rejection-round record grades as never-recorded
status: backlog
source: "Run 31996696789 claude lane, rejection-flow: round recorded at stream line 287 with entries=2, rework spawned line 309; captain approved filing owners for tolerated residual modes at the 0.27 composite-green ruling, 2026-08-17"
id: zf7rymtke3b6xp7r0337hjj4
---

## Problem

Two coupled defects around one residual mode. The mode: the live claude FO invoked `gate record --round validation/1` immediately after the rejection — before routing the correction — so the immutable round room durably holds the reviewer's 2 entries instead of the complete 4 the contract's record-after-correction order produces. The red is correct in substance. The label is not: `claudeRecordedRejectionRound` accepts only a success line pinned to `entries=4` (`rejectionRoundSuccess`, shared_round_recording_test.go), so a successful early record is reported as "resolved launcher never invoked `gate record --round validation/1`" — the lying-label class this release spent itself killing. A second label defect in the same oracle: the workflow-README immutability boundary condition reports under `rejection-round-missing` ("README.md changed from its exact expected bytes", observed run 31991864922 attempt 2), which reads as a recording failure when the round WAS recorded.

## Proposed approach

Split honesty from strictness in the oracle. Match the recorder's success line generically (`entries=(\d+)`) and grade the count as its own condition with an honest code (`rejection-round-incomplete`, naming got-vs-want and that the record preceded the correction); keep `rejection-round-missing` for genuinely absent invocations. Move the README-immutability condition to its own code (`rejection-workflow-doc-mutated`). Both proven by falsifying edits: replay the run-31996696789 stream bytes through the recognizer (early record must grade incomplete, not missing) and a README-mutated fixture (must grade doc-mutated, not round-missing). Skill-side, the record-after-correction order is already explicit in step 6; the owner tracks the mode's recurrence rate with the metrics instrument rather than adding prose the model already had.

## Out of scope

- The round recorder binary (neutral; its output already carries the true entry count).
- Rejection-flow fixture text (the Cycle-line target pin shipped with the 0.27 stack).
