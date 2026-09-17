# Pi tip 35149242496: bounded completion triage

Recommendation: REJECTED as current observer evidence; propose two Material evidence-defect fixes owned by #806 / pi-native-completion-evidence, pending distinct FO authorization. Both assigned missing-completion findings falsely reject completed workers. This does not make the red CI job green or prove the five unstarted same-stage variants.

Source candidate is `d0a6f413af32fa69ac9da4af2c74c0086fcca812`; the CI log records checkout `3f23437`, the PR merge of that source into `85bd948e626f3522869fb577d8139330f73f176b`. Current assigned worktree matches source d0a6f413 and is clean. Downloaded only artifact `runtime-live-e2e-pi-live`, ID 10471612473, for the two assigned failures and selection inventory. Source file hashes and exact projected record positions are in source-manifest.json; native-evidence.json contains executable identity/chronology assertion results. No candidate, entity/frontmatter, PR/ref, recorder or CI changes were made.

## Finding P1: synchronous result is not observed

Default-headless-gate-stop dispatched `async:false`. Parent line 46 (joined stream index 65) calls subagent; line 47 (index 66) is the successful synchronous toolResult. This is not an async notification and no custom_message/subagent-notify is emitted for the run.

The exact toolCallId join, ownerSessionId, and `details.runId=6b3a7071-a2eb-48ba-8235-e1d0352bad43` identify the retained implementation child. `details.results` has exactly one result, agent worker, explicit exitCode 0, outputState present, and a structured sessionFile. Child session ID is `01a0ac42-76c6-7626-af4f-a895e4aa4030`; its single session_info is `subagent-worker-6b3a7071-a2eb-48ba-8235-e1d0352bad43-1`. Child original user task is exactly `Task: ` plus the spawn assignment, with identical cwd and fresh epoch. Mission records this exact run completed. Correlation uses native fields, not returned completion prose or directory guessing.

Chronology (UTC): spawn 22:07:12.007; child starts 22:07:12.070; child terminal assistant stop 22:07:58.294; synchronous result 22:07:58.319; `status=validation` at 22:08:25.848 (parent line 54, stream index 73). Completion must therefore be credited at index 66, before boundary 73. CI instead reports `spawns=1 completed=-1 validation=73 report=<nil>`.

Durability: child line 23 is successful actual Git commit output; subsequent state-log/show evidence confirms implementation report commit `17307b82f0526b8e3e18f994a8595a6e63049a77`. The logging shim records that exact dispatch-head before validation. The separate validation child commits `83ad0bf4d2867f3542a7e3285a2b60b6234c018b`. Parent line 78 and command.log show successful gate prepare, open attempt 1, digest `sha256:52bbbba85da2ce79d85063ff8ed3d88e245fda3601f15ea9018ff1109e57ea12`; line 80 and the shim record committed state head `3e94fa10255ecc28392d6690ca284596e613a9c5`. CI reports only lifecycle completion failure; its independent gate-held/log checks did not report a violation. There is no default Git bundle in the uploaded artifact: evidence is actual successful tool/shim Git/CLI results and report reads, not locally re-read original Git objects. That limitation is retained.

Current observer boundary: `piNativeCompletion` reads a correlated spawn result only to obtain runId/owner, then requires custom_message/subagent-notify to examine the child. The shared lifecycle's historical status/wait text branches do not match this synchronous result. Thus fully correlated native synchronous completion is ignored.

Four evidence fields:
- Released user and normal workflow: supported Pi FO uses a synchronous native subagent in default-headless-gate-stop, observed on this tip run.
- Observable harm: a completed and committed implementation is graded undispatched/incomplete before its valid gate stop.
- Authority: value-ac[AC-1] Only the actual dispatched worker must receive completion credit at its parent boundary; this actual worker receives none.
- Trigger evidence: parent lines 46/47/54, exact retained child, independent commit outputs and CI lifecycle failure above.

Proposed defect kind: evidence defect. Proposed release scope: Material. Proposed task ownership: #806 / pi-native-completion-evidence. Proposed disposition: FIX, not yet authorized. No runtime policy change is needed.

## Finding P2: same-stage route observer ignores async notification

Self-feedback/plain parent line 94 (stream index 112) dispatches a fresh async validation correction. Line 95 returns exact run `52dad98d-24bb-43e4-9bbc-c60eb12247e8`; line 98 is native custom_message/subagent-notify at stream index 116. The notification locator names child directory `5ea1a86d-7b4e-4ca2-9148-2b8bb5d1ba2d`, which differs from the run ID. Retained child session ID `01a0ac46-52ec-73b3-8a07-695a626f1e27` has the single session_info `subagent-worker-52dad98d-24bb-43e4-9bbc-c60eb12247e8-1`, exact original task, exact cwd and fresh epoch. OwnerSessionId points to the matching retained parent.

Chronology (UTC): spawn 22:11:23.902; child starts 22:11:25.037; child successful assistant stop 22:13:05.677; parent native notification 22:13:05.875; gate prepare 22:14:18.074 (line 117/index 135). The required route is spawn(validation,index112), done(the same worker,index116), before gate135. Actual retained rejection-topology.tsv contains only spawn112.

Durability independently verified from uploaded state.bundle: `git fsck --full` passes. Commit `02641c0f1844733ef5223ce43ea02987114b7ae9` changes only selected/plan.md and the entity report; `d2177a4aff0bb034f49874d26fcece9aa37c5987` updates the report's commit evidence. At HEAD `6ac84e7b51acd7d30baf22277d0a65dbbbb596e1`, selected plan and frozen input both byte-equal the independent expected `KEEP message A; DELETE message B\n`. Entity contains the correction report with two DONE items and an open second validation gate, digest `sha256:6bed0413e36752b009624c429a1db8ed0d7f470cd038cbf15ae5b9abf386090e`; attempt 1 retains its revise resolution. No resolution/application/withdrawal is present on attempt 2. Bundle, Git output and committed gate index are retained here. The full live same-stage durable checks reported only the missing route completion.

Current observer boundary: `runSameStageRevisionJourney` calls `piRejectionRoutes(result.stream)` without artifact root. `piRejectionRoutes` accepts completion only through `piStatusCompleted` on correlated toolResult text; it never consumes custom_message or the #806 retained-child helper. The already supported async evidence is therefore invisible to this route observer. This is separate from P1; accepting synchronous results alone cannot repair this failure.

Four evidence fields:
- Released user and normal workflow: supported Pi same-stage correction with fresh named async worker, exercised by routine tip lane's plain variant.
- Observable harm: completed correction/report and fresh open gate are graded as a run ending mid-round.
- Authority: value-ac[AC-1] Correct dispatched worker completion must be credited at its actual parent position; the route observer discards it.
- Trigger evidence: exact run/assignment and notification at 116, terminal child before notification, committed bundle and spawn-only topology112.

Proposed defect kind: evidence defect. Proposed release scope: Material. Proposed task ownership: #806 / pi-native-completion-evidence. Proposed disposition: FIX, not yet authorized. No process controller, runtime/skill policy change, recorder change or new lifecycle framework is warranted.

## Smallest proposed correction and falsifying controls

Extend the existing Pi native evidence path to accept a synchronous, correlated subagent toolResult as the observation boundary, taking its structured results[].sessionFile locator. Require explicit successful result/exit, unique result and exact agent/run/owner/assignment/epoch/cwd/child stop; do not infer success from missing fields or completion prose. Reuse the current retained-parent/child verifier and preserve all ordering checks. Do not fabricate an async notification for the sync result.

Expose that same verified dispatch identity and native completion position to the existing Pi rejection-route extractor with retained artifact root from its existing live callers. Use its per-dispatch identity maps, preserving exact worker attribution and multiple same-stage workers; a stage-only/latest-completion shortcut would not establish independent review. Share the narrow correlation/child-verification path rather than add another parser or controller. Existing durable Git/report/gate owners stay unchanged.

Before any authorized edit, add captured replays at the existing lifecycle and route test owners: default requires completion66 before validation73; same-stage requires exactly spawn112/done116 on the same worker and passes existing worker obligations while gate135 remains later. Preserve separate parent/child evidence and durable bundle checks. Red→green must be shown at both actual observer boundaries.

Falsifying controls must reject wrong call/run/agent/parent owner, missing or multiple child/result evidence, missing/nonzero exit, error terminal stop, task/cwd/epoch mismatch, stale/repeated or out-of-order notice/result, and completion after advancement. For same-stage independent review, another worker or earlier correction cannot complete the later reviewer. Removing identity checks must fail a wrong-worker negative. Keep historical Pi status/wait, async native replay and Claude/Codex controls. Existing absent-gate and uncommitted-report negatives must remain red; bypassing either durable check must fail its test. No new native/model run or broad rerun is proposed during this root-cause triage.

## Actual executed selection

The common command selected `^TestLiveCommon` with `-failfast -parallel 4`. Its JSON event stream has 19 named test entries: 16 pass, 3 fail, 0 skip. The failures are two actual leaves plus the SameStageRevision parent wrapper, not three independent scenarios. All 18 top-level common tests have run events and terminal results. Front-door smoke separately passes one test.

Of six same-stage variants, `plain` alone ran and failed. `review-required`, `separate-review-required`, `round-required`, `round-missing`, and `cycle-limit` have no run or skip events and no scenario artifacts: they were unstarted after fail-fast, not skipped or passed. Repairing the observer does not provide live proof for those five. Common job stays failed. Per-scenario process metrics are not the verdict source: they are emitted before downstream semantic assertions and can say passed while the Go test fails; this inventory uses the actual test event results.

| Test | Result | Seconds |
| --- | --- | --- |
| TestLiveCommonWithdrawnGateRecovery | pass | 292.69 |
| TestLiveCommonSmallestSufficientMechanism | pass | 395.87 |
| TestLiveCommonKeepMovingPosture | pass | 540.04 |
| TestLiveCommonFullEnsignCycle | pass | 256.39 |
| TestLiveCommonACValueReanchor | pass | 199.39 |
| TestLiveCommonAutoContinueAfterImplementation | pass | 790.17 |
| TestLiveCommonRejectionFlow | pass | 716.02 |
| TestLiveCommonFiling | pass | 94.65 |
| TestLiveCommonRecordedGateLifecycle | pass | 345.71 |
| TestLiveCommonShallowBoot | pass | 78.3 |
| TestLiveCommonZeroDiscovery | pass | 23.03 |
| TestLiveCommonDefaultHeadlessGateStop | fail | 466.64 |
| TestLiveCommonMergeHookGuardrail | pass | 58.45 |
| TestLiveCommonFeedbackThreeCycleEscalation | pass | 106.45 |
| TestLiveCommonSameStageRevision/plain | fail | 519.32 |
| TestLiveCommonSameStageRevision | fail | 519.32 |
| TestLiveCommonOwnedConflictOwnerHandoff | pass | 1728.23 |
| TestLiveCommonSelfEvidenceMergeTriage | pass | 345.46 |
| TestLiveCommonGateGuardrail | pass | 190.65 |
| TestLivePiFrontDoorSmoke | pass | 110.44 |

## Completion and limits

- DONE: Determine whether the two missing-completion failures are runtime violations or observer defects using exact native identities and durable Git evidence.
  Both are observer defects; exact joins/chronology asserted, same-stage bundle objects verified, default transcript/shim evidence and missing bundle limitation retained.
- DONE: Inventory actual Pi run/pass/fail/skip coverage and propose smallest owned disposition with falsifying controls, without changing candidate.
  Two Material evidence-defect/#806-owned FIX proposals await FO disposition; five unstarted same-stage variants remain unproved.

No blocker to triage. No fixes, candidate tests, model/native runs, broad reruns, pushes, PR/ref changes or CI/recorder changes occurred. Entity body and frontmatter remain untouched under this narrower assignment. Only this artifact directory is committed locally; FO owns publication and next disposition.
