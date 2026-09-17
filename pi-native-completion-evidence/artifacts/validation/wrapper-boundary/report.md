# Retained Pi workflow-wrapper boundary

Conclusion: the retained plain leaf does **not** contain a complete structured workflow-to-child receipt. It contains useful outer identity and child epoch boundaries, but these do not suffice to admit wrapper completion through the existing strict checks. Stop at this missing-evidence boundary. No adapter implementation or acceptance rerun is recommended from these bytes alone.

## Checklist

- DONE: Determine whether retained Pi workflow and revival evidence contains sufficient structured identities to use existing completion checks without weakening them.
  The outer inventory is explicitly incomplete and empty; the referenced receipt is absent from this retained leaf; the notification's embedded JSON cannot be decoded. Epoch boundaries are present but do not replace the missing join.
- DONE: Give the smallest concrete implementation proposal or a precise missing-evidence boundary, preserving the actual incomplete plain outcome and existing approved scope.
  Separate extraction, epoch and host-wording boundaries below. HOLD pending complete existing receipt evidence and distinct FO authorization; original plain remains FAIL.

## Summary

This refines Finding 2 from tip35238049892/report.md; it does not reopen accepted final-tip-7c63 validation or introduce a separate product defect. Two investigation items DONE, zero SKIPPED/FAILED does not mean wrapper acceptance passed. Candidate 7c63ad87f34a18c80e3e8b24f6d201e227f661e6 and its clean worktree are unchanged. No code, canonical entity, frontmatter, refs, tests, model, CI or push action occurred.

## Exact retained join inventory

All line numbers are one-based in the original files identified and hashed in boundary-check.json. selected-records.jsonl retains only relevant native records, without thinking/signatures. The parent is session 01a0aff0-7bf5-76e8-b01f-7eaf97428d76. The child session's actual ID is 01a0aff4-42bd-76f9-9ea6-f2ebf79cba92; its directory UUID is not treated as run identity.

| Boundary | Actual retained fields | Sufficiency |
| --- | --- | --- |
| Parent request to outer workflow | Parent88 toolCall ID; parent89 matching message.toolCallId, isError=false, details.mode=workflow, details.runId=e50172d2-7e00-49f9-be6c-f78a814092c8 | Sufficient outer dispatch/result join |
| Outer owner | Parent89 details.mission.ownerSessionId points to this exact parent; mission.cwd matches parent cwd | Existing piRetainedChild owner/path checks can preserve this authority |
| Workflow inventory | Parent89 details.workflowChildren.version=1, parentToolCallId matches, workflowRunId matches, inventoryComplete=false, workflowState=running, children=[] | Explicitly insufficient to establish any launched child or completed inventory |
| Workflow completion | Parent94 custom_message/subagent-notify has a receipt path and a Return text preview; no structured details/result block | Receipt contents are not present; preview is insufficient |
| Initial child epoch | Child4 session_info.name identifies worker/run5a2756d7-4e57-4aa7-a9a4-7b70a909d545/epoch1; child5 exact task; child40 assistant stop | Establishes child-local identity, assignment and successful stop, not structured ownership by the outer workflow |
| Repair request/result | Parent103 action=resume, id=5a2756d7..., exact message; parent104 correlated successful toolResult details.runId=27266d0c-c89a-4409-afeb-c459cc794e44, context=fresh, launchContractDigest and sourceLaunchContractDigest | Establishes old-run request to new-run result; no ownerSessionId or original structured launch digest exists here to complete independent ownership correlation |
| Repair child boundary | Child41 session_info.name names new run; its parentId equals child40.id; child42 revival assignment; child47 assistant stopReason=error | Explicit linked epoch boundary and failed terminal; cannot credit success |

Parent88's workflowScript contains a single literal `return runs.run(key, {agent, skill, context, cwd, task})`. The script is a JSON string, not an evaluated child launch record. A regex or JavaScript evaluator for arbitrary scripts would change the evidence boundary and is not a narrow substitute for a receipt. The embedded literal task matches the initial child assignment, but that alone is not a structured outer/inner run join.

Parent94 references `/tmp/pi-subagents-uid-1001/async-subagent-runs/e50172d2-7e00-49f9-be6c-f78a814092c8/workflow-receipt.json`. The retained leaf inventory contains stdout/stderr, model/process/duration, topology, state.bundle and two native session files; no receipt, events.jsonl, status.json or mission JSON was retained there. This is a claim about the available artifact, not about whether the runner originally persisted those files or whether some other archive might hold them.

The notification preview names the inner run, agent, requested/resolved fresh context and artifact path, but truncates inside `results[0].sessionName` before “Trace: 2 event(s).” A read-only standard JSON raw-decode probe fails at offset1019 with an invalid control character. Its trailing “Child runs … (completed)” text is not a complete structured receipt. Recovering selected keys from that broken preview would accept prose/partial data as the missing join. The probe merely checked these retained bytes; no replay/model/new harness was run.

## Epoch selection is a separate, partially supported change

The child file has one session header and two distinct session_info/user epochs. Initial session_info is child4, task5, successful terminal40. Revival session_info41 links directly to terminal40 by parentId, followed by task42 and error terminal47. The fresh epoch ends before revival begins; first stop15:21:30.736Z precedes workflow notification15:21:57.612Z, then revival boundary15:22:53.753Z. Thus selecting the earlier validated epoch need not erase the later failure. Conversely, selecting the last successful assistant anywhere would hide that failure and is unacceptable.

Current piVerifyNativeChild reads the entire file, requires exactly one session_info and task, and takes its final record as terminal. piSessionRecord currently omits parentId. A future authorized bounded selector could retain the original session header/cwd, identify exactly one requested run/epoch, and select that epoch's task/terminal up to the next linked boundary. A revival's start must be its own boundary timestamp rather than the old session header timestamp. Unknown/duplicate identities, unlinked boundaries, mismatched task, wrong owner/path, early/late terminal and error terminal must still reject. Same-role revival must not be counted as an independent reviewer.

This structural observation is not sufficient authorization or data to complete wrapper acceptance: the initial outer-to-inner receipt join is missing. The revival message's “Original run” and “Follow-up” text are corroborating assignment bytes, not an alternative authority source. Parent104 lacks a mission owner and its sourceLaunchContractDigest cannot be compared with an original child launch digest in the retained records.

## Smallest next step and conditional implementation surface

Proposed disposition: retain HOLD / Needs decision on evidence sufficiency. First recover the **already referenced existing workflow receipt**, if it survives in the completed run's stored artifacts. Required contents to inspect are its outer parent-call/workflow/owner identity and complete child inventory: exact child run, agent, resolved assignment/context/cwd, session locator, and successful/failed terminal outcome. This is an evidence request, not authorization for a new run, instrumentation, discovery framework or candidate change. If the receipt cannot be recovered, stop; no code estimate can substitute for missing authoritative schema.

Only if the complete receipt supports those joins should FO authorize a narrow adapter proposal. Exact existing functions/types affected would be:

- piSessionMessage/piSessionRecord and piToolCallArgs in pi_rejection_extractors_test.go: decode proven receipt/parentId/resume fields, without parsing arbitrary JavaScript.
- piNativeCompletions and piNativeDispatch in pi_native_completion_test.go: normalize a receipt-backed child with distinct outer and inner identities and actual timeline indexes; never equate workflow ID with child ID.
- piVerifyNativeChild in that file: select a correlated epoch before applying the existing task/cwd/path/terminal/time checks; piRetainedChild remains the ownership and filesystem guard.
- piRejectionRoutes in pi_rejection_extractors_test.go: consume normalized child dispatch/completion identities without merging reviewer identity or counting a revival as fresh independent review.
- Existing pi_auto_continue_double_dispatch_replay_test.go and/or pi_rejection_extractors_test_test.go: captured receipt/epoch positives plus wrong outer/child/owner, missing/ambiguous receipt, duplicated epoch, error and downstream gate/report negatives.

Planning estimate only, contingent on the actual missing receipt schema: 2 existing helper files, 1–2 existing test-owner files, compact projected receipt/parent/child fixtures; approximately80–140 helper lines and100–180 test lines. This is not an approved estimate or a reason to implement now. No new general parser/controller/event framework is proposed. Existing lifecycle/Git/report/gate owners stay in place. Receipt recovery would either substantiate or invalidate this estimate before a correction assignment.

## Fixture wording is independently owned

rejectionHostRealization in claude_live_runner_test.go has a Codex-specific arm and a Claude Agent/name/run_in_background/SendMessage fallback. Pi receives that fallback in the retained parent prompt. A runtime-appropriate Pi wording correction is a separate fixture-owner decision. It neither establishes receipt identity nor proves Pi will avoid workflows. Do not couple that wording change to a guarantee of direct task-bearing subagent calls or use it to close this gap without evidence.

## Finding 2 policy refinement

Classification remains Material evidence-observer gap with an unresolved evidence boundary; #806 observer owner coordinates with fixture owner, and FO owns the disposition. This is refinement of the held finding, not a new fix authorization.

1. Released user/normal workflow: authorized Pi same-stage correction dispatched through a native workflow wrapper and followed by report repair.
2. Observable harm: the existing route grader reports no workers despite child-local successful correction/commit evidence; the available archive cannot support a strict structured wrapper-to-child join.
3. Authority: value-ac[AC-1] Completion must credit only the exact dispatched worker at the correct boundary, without substituting partial notification text for identity evidence.
4. Trigger evidence: parent89 inventoryComplete=false/children=[], parent94 absent referenced receipt plus truncated Return JSON, and child4/41 distinct linked run epochs. These are observed retained artifacts, not a hypothetical runtime schema.

Even a complete future extraction must leave this original plain journey FAIL: report repair terminal47 is an error (“Codex error: Unable to verify Daybreak Blue access. Please try again.”), the committed report omits its second checklist item, and the retained state bundle contains no replacement open gate. No authentication root cause is inferred. The accepted 7c63 bounded deterministic result, all prior real CI failures and five unstarted Pi variants remain unchanged.
