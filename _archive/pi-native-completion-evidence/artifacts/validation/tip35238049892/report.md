# Pi retained live triage — run 35238049892

## Stage Report: validation

- DONE: Establish exact native and durable outcomes for both failing Pi leaves, distinguish runtime failure from observer defect, and audit actual variant execution.
  Native parent/child chronology, plain state bundle and Go events are retained alongside this report.
- DONE: Propose materiality, ownership and disposition with the four workflow evidence fields, preserving existing successful controls and unchanged candidate.
  Three bounded findings below remain proposals for FO disposition; no repair or rerun was performed.
- DONE: Commit artifact-only report and selected evidence with source hashes and explicit acceptance limits.
  This path-scoped artifact commit owns only this directory; source-sha256.json identifies retained inputs.

## Summary

Pi job 105259444883, artifact 10505664184 remains red. Same-stage plain has a real incomplete report-repair outcome and an independent observer blind spot. Continuation completed both validators and prepared a replacement gate, but its observer combines different attempts. Neither failure demonstrates that workers were absent.

## Provenance and limits

Candidate is 0997b6b7590a5a89ba979194238e93f7f47d94e4, unchanged. CI checked out synthetic merge d5b3190b36b824f521d276b11d04ef6ce5b3f2f5, not a presumed main checkout. The retained step log says “Merge 0997b6b7590a5a89ba979194238e93f7f47d94e4 into 1b2fa6b7a98e0ffb2c76966d50364c151f95dfc8.” My local lookup could not resolve d5b3190. FO separately supplied GitHub commit API evidence: parents are that base and candidate, tree 2b5fd3fa31f735ee1f227310a30222200985d898. FO's local candidate tree lookup returned the identical tree. Tree equality is attributed to FO verification, not an API call performed by this validator.

Selected native records retain source path and one-based line number; thinking and signatures are omitted. The plain bundle was cloned into a disposable bare repository for read-only inspection. fsck and recorded Git reads exited 0. No continuation bundle was retained: its Git commit and gate evidence comes from actual tool results, not an independent object replay. No new tests, native/model calls, candidate edits, frontmatter updates, pushes or CI runs occurred.

## Actual execution

| Same-stage variant | Pi outcome |
| --- | --- |
| plain | FAIL, 450.39s |
| review-required | Unstarted; no run/pass/fail/skip event |
| separate-review-required | Unstarted; no run/pass/fail/skip event |
| round-required | Unstarted; no run/pass/fail/skip event |
| round-missing | Unstarted; no run/pass/fail/skip event |
| cycle-limit | Unstarted; no run/pass/fail/skip event |

Continuation failed in 604.93s. Coverage package failed in 2217.942s. The inventory retains all terminal events: 16 other top-level common journeys passed. Front-door smoke passed (its separate package took 158.220s). These controls do not turn the five unstarted same-stage variants into passes. Claude's separate report at workflow-owned-same-stage-revision/artifacts/validation/tip35238049892/report.md (66bf2a4b0) owns Claude 5/6 and Codex 6/6 actual same-stage outcomes, plus Claude missing-source recording, roadmap refusal and naming/recovery findings. They are not new #806 defects. Pi's smallest-mechanism case passed in this run.

## Finding 1 — plain correction committed, repair terminated with provider error

Proposed classification: material live outcome failure. Ownership: runtime/report-repair outcome, coordinated by FO; do not misfile an authentication configuration defect. Disposition: HOLD; retain unchanged strict scenario as incomplete until FO chooses the next action.

1. Supported workflow/user: captain-authorized same-stage correction, complete stage report, then next decision boundary.
2. Observable harm: corrected plan exists, but the report omits a dispatched checklist item and no replacement open gate was prepared.
3. Authority: contract[skills/ensign/references/ensign-shared-core.md#stage-report-protocol]
4. Observed trigger: initial worker completed with one of two checklist items accounted for; revival to repair its report ended with a native provider error before any repair commit.

Parent session 01a0aff0-7bf5-76e8-b01f-7eaf97428d76 dispatched workflow e50172d2-7e00-49f9-be6c-f78a814092c8 at line 88. The workflow completion at line 94 identifies child 5a2756d7-4e57-4aa7-a9a4-7b70a909d545 and its session artifact. Child lines 35/39 record successful commit e6495fe040584f41325f1e92aa9eac653b6b0a10, comparison success and clean staged state; line 40 stops normally. This is corroborated by bundle HEAD e6495fe: plan and frozen input are byte-identical. The entity still has only the prior revised gate attempt. Its cycle-2 report omits “Append a ## Stage Report: validation documenting the correction and validation evidence, then commit the entity state.”

The same child file is revived: session_info line 41 names run 27266d0c-c89a-4409-afeb-c459cc794e44; line 42 carries the report-only repair assignment. Line 47 has stopReason:error and errorMessage “Codex error: Unable to verify Daybreak Blue access. Please try again.” It is direct native error evidence, not an inference from final prose. The error does not identify credential root cause or prove a host-wide outage. The initial report omission and interrupted repair remain distinct facts. Gate/report acceptance cannot be granted merely because plan correction succeeded.

## Finding 2 — workflow wrapper is outside the current Pi observer schema

Proposed classification: material evidence-observer gap, with a fixture host-instruction mismatch requiring disposition. Ownership: #806 native observer, coordinated with same-stage fixture owner. Disposition: HOLD / Needs decision; no blanket acceptance repair justified by this incomplete journey.

1. Supported workflow/user: Pi-native named correction workers carrying the generated ensign assignment.
2. Observable harm: the grader reports “the run routed no workers at all” despite a retained child dispatch, successful initial completion and committed correction.
3. Authority: contract[skills/ensign/references/pi-ensign-runtime.md#runtime-implementation]
4. Observed trigger: parent line 88 uses subagent workflowScript with runs.run, rather than a top-level task argument; report repair later revives the same child session file.

Current piNativeCompletions and piRejectionRoutes recognize task-bearing subagent calls. The workflow call has no top-level task and its result mode is workflow. Its outer run ID differs from the child run ID. The workflow completion text names both, but this is not the existing single-run Session-file notification schema. The retained child also now contains two session_info names and two user assignments: passing the whole revived file to the current fresh single-run verifier would violate its identity/terminal assumptions. These are concrete schema differences; a prose “completed” match or treating the outer run as child identity would be unsafe.

Current rejectionHostRealization has a Codex arm and otherwise supplies Claude Agent/name/run_in_background/SendMessage wording, including for Pi. The actual Pi prompt contains that wording. The model chose a named runs.run wrapper. This is a fixture mismatch; causation of the wrapper choice is not proven.

Smallest proposal: first decide the intended Pi realization in the existing host adapter. If workflow wrappers are supported, require a correlated outer tool result and structured child receipt/inventory plus parent, run, assignment, epoch and terminal checks; do not accept truncated notification prose alone. Treat each revival epoch explicitly, preserving the error terminal. Do not add a parallel harness or weaken missing/error completion and wrong-owner/run negatives. Even perfect route extraction must leave this particular plain journey red because report repair and replacement gate are absent.

## Finding 3 — continuation observer crosses gate-attempt boundaries

Proposed classification: material observer chronology defect; separately retain initial report omission as a recovered runtime defect. Ownership: #806 lifecycle observer. Disposition: propose bounded correction for FO approval; current live result remains FAIL pending approved implementation and verification.

1. Supported workflow/user: complete independent validation and report repair before presenting the current open decision gate.
2. Observable harm: validation-worker-not-dispatched falsely describes two completed native validators; the predicate combines the latest completion with an earlier withdrawn gate attempt.
3. Authority: contract[skills/ensign/references/pi-ensign-runtime.md#runtime-implementation]
4. Observed trigger: first gate preparation precedes a report-repair dispatch; FO withdraws it, receives the repair completion, and prepares attempt 2.

Parent session 01a0affc-3251-7468-b97d-5a46f9d90fd3 lines 44/45 dispatch and correlate run ac0f4b70-1733-49ad-9e6d-f3fb3a5af03f. Its child identity is 01a0affd-ca87-705b-b027-eb776546a0c6, session_info names that run, starts 15:30:40.392Z and stops normally 15:32:32.311Z. Notify line 50 follows at 15:32:32.521Z. Actual git show and worktree results at lines 54–56 establish committed report 4a084486d33c801c7d480b74cf4e67c856dad3d0.

First gate prepare is line 81, observer stream index 99. The report lacked the second dispatched checklist item. Line 110 explicitly confirms attempt 1 withdrawn. Lines 117/118 dispatch repair run 26ab8c8d-5af7-4359-a934-7a8a38ca6562. Child 01a0b003-7bcd-70c3-bfa2-f970a9e638c4 has matching run identity, starts 15:36:53.454Z, and stops normally 15:38:08.227Z. Parent git show line 122 verifies repaired report with both DONE items at e5c47e7d33d16d176e600a0d20f8fb8f5b40d638. Notify line 123 follows at 15:38:13.557Z, observer index 141. Second prepare is line 126 (index 144). Line 135 returns attempt-2 briefing bound to that exact repair commit and digest sha256:2b646a4dc973209ddc75a52f3b4b85bafa16bba41c702f1f4770ccf6bd447fb8. No approval/consumption is observed.

The existing native async correction is exercised: completed=141 is the actual second notification. piNativeCompletion selects the latest dispatch; assertWorkerLifecycle retains the first nextSignal. Its resulting 141 >= 99 rejection is consistent with that implementation, but not with the final attempt's chronology. This is different from the prior missing async/sync completion support; no synchronous result is claimed for these two failing leaves.

Smallest proposal: use the existing correlated native dispatch/completion records to associate each gate boundary with its preceding completed worker, and verify the current gate's committed report after repair. Preserve the first-attempt chronology and explicit withdrawal; do not simply select the last gate or earliest completion globally. Keep missing/error completion, wrong identity/owner, premature gate, omitted withdrawal, missing report and uncommitted report negative controls. The recovered initial checklist omission must not disappear from the evidence. No new general observer framework is needed.

## Acceptance limits

This is completed triage, not completed live acceptance. Plain remains a real incomplete runtime outcome; continuation remains a recorded red pending FO disposition. Five same-stage Pi variants did not start. Existing independently validated #806/#801 controls and front-door smoke remain evidence only for their actual scope. Earlier deterministic evidence in final-tip-0997 still retains the original default-timeout failures and the exact resolver failure from the authorized serial normal/race runs; this investigation changes none of those results.
