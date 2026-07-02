---
title: "codex-live reuse characterization: negation must bind to the concept, not co-occur in the message"
status: backlog
group: tooling
source: "PR #464 codex-live (2026-07-02, run 28568410914): TestLiveCodexSharedScenarios/rejection-flow failed 'addressable-worker was characterized ABSENT, but transcript contains turn-starting reuse tool send_input' — while the scenario's behavior was CORRECT (two-cycle rejection flow, kept-alive reviewer, clean tree). Root cause read from the artifact: codexNarrationNegatesReuseRoute (internal/ensigncycle/shared_reviewer_reuse_test.go:391) flags any single FO message containing a reuse concept AND any negation token; both tripping messages AFFIRM reuse with unrelated negations ('does not close or redispatch the worker', 'the host has no dedicated shutdown tool ... the two reusable workers'). Same over-broad live-heuristic family as the 2w wrong-root detector fixed in #462. Narration-dependent: can red any branch."
id: mq42796928asq686vxs74mpm
---

## Problem
The codex runner characterizes the addressable-worker capability from FO narration: one message containing a reuse-route concept ("follow-up", "followup", "addressable", "reuse", "reusable") plus any negation token ("no ", "not ", "cannot", "n't", "without ", "never ") anywhere in the message flips the harness to its ABSENT assertion branch. An affirmative-reuse sentence with an unrelated negation ("waiting for the reused worker; interruption does not close it") mischaracterizes the runtime, and the harness then flags the FO's legitimate reuse as a contradiction. The function's own doc comment claims message-scoping prevents exactly this; it does not.

## Desired direction (for ideation to refine)
The negation must bind to the reuse-route concept, not co-occur: candidates (ideation picks the simplest that kills both observed false positives while keeping the true-positive fixtures red) — a bounded window between negation and concept, explicit characterization phrases ("followup_task is not available", "no reuse route", "send_input is unavailable"), or keying on the FO's capability-probe output rather than free narration. The two real tripping messages from run 28568410914 become checked-in must-NOT-match fixtures; the existing true-positive shapes stay matched. Consider whether claude/pi runners share the pattern.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- Both captured PR #464 tripping messages replay through the characterization returning PRESENT-compatible (no ABSENT flip); genuine absence narrations still flip it (RED/GREEN fixture pair).
- The rejection-flow scenario passes on a transcript replaying the full #464 artifact (behavioral proof, not prose assertion).
- Harness-internal change only; no product surface. Offline fixtures are the gating proof; one codex-live run as field confirmation.
