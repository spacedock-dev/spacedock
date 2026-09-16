# Tip CI same-stage investigation

Run [35148797301](https://github.com/spacedock-dev/spacedock/actions/runs/35148797301), exact head d0a6f413af32fa69ac9da4af2c74c0086fcca812. This is a read-only investigation artifact, not a replacement canonical Stage Report. Approved entity body/frontmatter were preserved. No candidate edits, native/model runs, broad reruns or push.

## Results inventory

Both hosts ran all six variants exactly once with no skips. Archived Go result records are retained in each host's same-stage-detail.jsonl/outcomes.json.

| Variant | Claude Sonnet5 | Codex |
|---|---|---|
| plain | PASS252.49s | PASS164.35s |
| review-required | FAIL309.68s | PASS145.87s |
| separate-review-required | PASS506.53s | PASS157.57s |
| round-required | PASS233.66s | PASS186.71s |
| round-missing | PASS233.09s | PASS135.28s |
| cycle-limit | PASS253.79s | FAIL173.20s |

Additional requested current-tip positives: TestLiveSemanticNamesCodex PASS85.93s; TestLiveCommonOwnedConflictOwnerHandoff Claude PASS310.96s and Codex PASS110.60s. Each has run+pass records, not skip; retained additional-detail.jsonl provides exact evidence. Pi was not an artifact/lane in this run, so this report makes no Pi execution claim.

## V3a — Claude worker-name collision misclassified as self-review

Exact source: claudeRejectionRoutes stores Agent input.name as target; assertSameStageWorkers compares rounds[1].target == rounds[0].target for independent review. CI reports two spawn/done rounds with equal semantic name recorded-gate-task-validation, so rejects independence.

Native source proves two fresh identities: first Agent tool_use toolu_013tQCCE4kV67J8wWfYJttpV returns task a891653619c5fb829; its matching completed task_notification occurs at zero-based line166. Second Agent tool_use toolu_012HxAjBzxemi3YCjWsCFFCt returns task a6f8a511fd35edfb4; matching completion is line331. Their spawn positions are94/274; both input.name values are identical. The host distinguishes them by native task and tool identity; the grader discards that distinction. Full bounded scenario stream plus identity-observations.json retain exact correlations.

Declared policy requires a different reviewer after correction, with missing external selected/reviewer-source.txt prohibiting a fabricated verdict. Matching state.bundle HEAD ccc760935e6dacfcc8804307e3989c26a5db7723 contains correction98756ad and fresh reviewer hold report ccc7609. Exact corrected plan/frozen bytes match; required source remains absent. One original rejected attempt remains, with no application/new gate/round. Original briefing and frozen input bytes match their earliest commits. Thus actual independent review correctly held; the reported self-review violation is an evidence-grader false rejection, not an observed model violation.

Four evidence fields: supported user/workflow is the required Claude routine same-stage review control; observable harm is false failure of independent-review proof for distinct native workers; authority is value-ac[AC-2] declared independent review must be evaluated by actual worker independence; trigger is two fresh same-stage workers assigned the same semantic name in this exact CI run.

Recommendation: Material evidence defect, task-owned extractor/grader boundary, FIX proposed pending distinct FO authorization. Preserve semantic name/stage separately from native worker identity. Correlate Claude spawn and completion using native task/tool identity rather than display name. Falsifiers: captured two-fresh-IDs/same-name must pass identity independence; same native identity under changed labels must still fail self-review; wrong/missing completion owner must fail. Do not weaken required-review hold or invent a reviewer verdict.

## V3b — Codex reused correction turn miscounted as extra worker

Exact source: assertSameStageWorkers calls parseRejectionRounds then requires len(rounds)==1 when no review is required. parseRejectionRounds deliberately produces one record per completed spawn/reuse turn, not per worker identity. CI has spawn+done then reuse+done on the same task, so the predicate reports extra workers.

Native source: spawn call call_V8lU1VVXKqm3vlGDc2ij2X9Y returns /root/recorded_gate_task_validation at line112; attributed FINAL_ANSWER/Done closes it at136. followup_task targets the same worker at151 (call_pH4Ubx6bZQiWNEhnC8iNLwkL); attributed completion at163 closes the followup. There is exactly one successful spawn and one reused worker identity, not two workers. Followup message body is encrypted; its outcome is established by the exact same-owner completion and corresponding committed escalation, not guessed from ciphertext.

Declared policy has no reviewer/round obligation and requires durable cycle3 escalation. State.bundle HEAD794eb95f3f817a0383f83c4180cf0f9c24cb143d contains correction efba244, correction report94269a2 and followup escalation794eb95. All three original rejected attempts and frozen briefings are unchanged; no fourth attempt/application/round exists. Plan equals frozen bytes. The worker's first report explicitly left escalation pending; the completed followup records it and stops for captain direction. This is legitimate completion of the existing correction, not an extra review worker or fourth cycle.

Four evidence fields: supported user/workflow is the no-reviewer cycle-limit control with addressable worker reuse; observable harm is false failure of a completed one-worker correction/escalation; authority is value-ac[AC-1] self-feedback correction must complete without invented review machinery; trigger is a completed same-worker followup in this exact routine Codex run.

Recommendation: Material evidence defect, task-owned same-stage grader, FIX proposed pending FO authority. Keep fail-closed ordered completion pairing, but count native correction worker identities separately from completed turns; allow valid same-worker reuse to finish the correction/escalation. Falsifiers: this one-spawn+completed-reuse trace must pass; missing followup completion, extra fresh correction worker, self-review, or fourth gate must still fail. Do not relax cycle limit or gate/round requirements.

## Scope and completion

- DONE: Identify each assigned failure at its actual observation boundary using retained runtime and durable evidence.
  Both are grader false rejections with native identity and matching immutable/committed state evidence, detailed above. Ten other variant results pass; all twelve ran once.
- DONE: Propose smallest disposition with affected AC, task ownership, and falsifying check; do not mutate candidate or weaken assertions.
  V3a/V3b propose a bounded existing extractor/grader identity correction, preserving completion order, independent-review and gate/round guards. No workflow/skill change or new lifecycle framework is indicated. No fix is authorized by this investigation.

### Summary

Current tip CI exercised the promoted same-stage journey on both hosts. Its two failures confuse display names and completed turns with native worker identity; retained durable results satisfy the intended controls. Findings and focused falsifiers are ready for distinct FO disposition, without changing the approved entity or candidate.

FO disposition received: both findings are task-owned Material evidence defects in native identity extraction and worker cardinality. Preserve actual outcomes and missing/wrong-owner/self-review/extra-fresh-worker negatives. FO will route minimal correction to the existing implementation owner; reviewer has no candidate mutation authority.
