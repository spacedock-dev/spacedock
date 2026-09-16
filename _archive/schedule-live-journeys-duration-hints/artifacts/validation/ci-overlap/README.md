# CI overlap evidence — run 35058663190

Conclusion: actual common/substrate overlap, three occupied Go-test slots, and distinct retained runtime paths are demonstrated. **AC-2 is not fully closed:** 20 functions were invoked, but owner-handoff failed before a host launch, rejection-flow hit the quiet budget, and keep-moving is an existing semantic XFAIL. Both sequential break-glass variants executed successfully (bare is XPASS; team passes). This is evidence review only; the approved local report, gate, frontmatter and candidate are unchanged.

## Sources and provenance

- Run: https://github.com/spacedock-dev/spacedock/actions/runs/35058663190; Claude Sonnet artifact downloaded under `/tmp/spacedock-tip-ci/claude`.
- PR source head: `ac3c06a445fcfbc773639254ff7242965248ba78`.
- Executed checkout and embedded binary revision, both recorded in `spacedock/spacedock/live-artifacts/claude/claude-sonnet-5/candidate-binary-provenance.txt`: `c03362a699626fa7be44de86504102ce08861e33`, `vcs_modified=false`.
- FO independently confirmed GitHub PR merge parents `[95de6889a7b1fcb6b648584e19c5355d06d9c568, ac3c06a445fcfbc773639254ff7242965248ba78]`; compare reports `files:[]`, so executed merge has the source tip's tree. This API attribution is FO-provided, not inferred from artifact timestamps.
- `source-SHA256SUMS` records every downloaded source file's relative path and SHA-256, including candidate binary, JSON test events, streams, metrics and project transcripts. `evidence.json` retains extracted event line numbers, exact nanosecond Go timestamps, millisecond host timestamps, session/workflow/config paths, outcomes and minimal failure/residue excerpts. No full prompts or credentials are copied into this report.

## Coverage and occupied slots

A sweep over exactly one `run` and one terminal event for each `TestLiveScheduled/slot-N/<function>` yields peak **3**, ending at zero. Nested variants are not additional slots. There are 18 Go passes, 2 Go failures and **0 skips** among these 20 functions. A Go pass can contain XFAIL and is not automatically a clean durable assertion result.

All times below are UTC on 2026-09-16. Exact timestamps/line numbers are in `evidence.json`.

| Function | Slot | Start | Terminal | Go outcome |
|---|---:|---|---|---|
| TestLiveCommonAutoContinueAfterImplementation | slot-0 | 05:16:23.266 | 05:21:46.655 | pass |
| TestLiveCommonRejectionFlow | slot-2 | 05:16:23.267 | 05:20:21.874 | fail |
| TestLiveCommonFeedbackThreeCycleEscalation | slot-1 | 05:16:23.267 | 05:19:09.457 | pass |
| TestLiveCommonACValueReanchor | slot-1 | 05:19:09.457 | 05:22:12.877 | pass |
| TestLiveCommonKeepMovingPosture | slot-2 | 05:20:21.874 | 05:24:37.605 | pass |
| TestLiveCommonOwnedConflictOwnerHandoff | slot-0 | 05:21:46.655 | 05:21:46.934 | fail |
| TestLiveBreakGlassShimRecovery | slot-0 | 05:21:46.934 | 05:29:37.057 | pass |
| TestLiveCommonSmallestSufficientMechanism | slot-1 | 05:22:12.877 | 05:25:42.180 | pass |
| TestLiveCommonDefaultHeadlessGateStop | slot-2 | 05:24:37.605 | 05:28:06.282 | pass |
| TestLiveCommonRecordedGateLifecycle | slot-1 | 05:25:42.180 | 05:28:53.405 | pass |
| TestLiveCommonFullEnsignCycle | slot-2 | 05:28:06.282 | 05:30:20.236 | pass |
| TestLiveMergedTeamModeDispatch | slot-1 | 05:28:53.405 | 05:30:49.709 | pass |
| TestLiveCommonSelfEvidenceMergeTriage | slot-0 | 05:29:37.057 | 05:32:57.049 | pass |
| TestLiveBareReachable | slot-2 | 05:30:20.236 | 05:31:38.400 | pass |
| TestLiveCommonWithdrawnGateRecovery | slot-1 | 05:30:49.709 | 05:31:51.211 | pass |
| TestLiveCommonGateGuardrail | slot-2 | 05:31:38.400 | 05:33:15.860 | pass |
| TestLiveCommonFiling | slot-1 | 05:31:51.211 | 05:32:39.439 | pass |
| TestLiveCommonMergeHookGuardrail | slot-1 | 05:32:39.439 | 05:33:04.379 | pass |
| TestLiveCommonShallowBoot | slot-0 | 05:32:57.049 | 05:33:14.977 | pass |
| TestLiveCommonZeroDiscovery | slot-1 | 05:33:04.379 | 05:33:17.879 | pass |

Both break-glass variants remained within slot-0: selected-bare ran 05:21:47.172083473–05:23:55.032575202; selected-team began 05:23:55.032584279 and passed at 05:29:37.056647265. Bare emitted `XPASS ALERT ... observed=[]`; its metric is `xpass`, not XFAIL. Team metric is `passed`. No TODO skip occurred. Queue progress after rejection-flow and owner-handoff failures is directly visible in the subsequent slot events.

## Direct runtime overlap and path evidence

The following intervals use timestamped events in actual Claude project transcripts, not just Go test intervals. They bound observed session activity, not exact OS process start/exit or uninterrupted CPU activity. Retained init streams independently bind session IDs to launch cwd. Cross-session tool events are interleaved within these windows.

- Break-glass selected-bare and common AC-value-reanchor overlap 05:21:47.691–05:22:12.365 (24.674 seconds of intersecting observed session spans).
- Merged-team-mode and common full-ensign-cycle overlap 05:28:53.940–05:30:19.884 (85.944 seconds).
- Bare-reachable and common self-evidence-merge-triage overlap 05:30:20.962–05:31:38.073 (77.111 seconds).

Config prefix `C` is `/home/runner/work/_temp/spacedock-claude-config/claude-sonnet-5`. Each row below has its own `C/<scenario>/projects/<encoded-workflow>/<session>.jsonl`; complete source paths and SHA-256 are in `evidence.json` and the manifest.

| Scenario/config child | Observed first–last UTC | Session | Initial workflow root |
|---|---|---|---|
| ac-value-reanchor | 05:19:10.215–05:22:12.365 | `bf5fac0d-0c35-46ee-9c86-e3f80f2ed963` | `/tmp/TestLiveScheduledslot-1TestLiveCommonACValueReanchor2844895375/003` |
| bare-reachable | 05:30:20.962–05:31:38.073 | `4f4ed0a4-9dbf-49ba-9ada-b18ff7b86616` | `/tmp/TestLiveScheduledslot-2TestLiveBareReachable1041975429/003` |
| break-glass-shim-selected-bare | 05:21:47.691–05:23:54.719 | `a89e2acd-f0a2-42d8-b739-25254fa18f49` | `/tmp/TestLiveScheduledslot-0TestLiveBreakGlassShimRecoveryselected-bare120101119/001` |
| break-glass-shim-selected-team | 05:23:55.641–05:29:36.717 | `f4b63484-3f14-46ed-b6ad-f8a7aa4f4bb6` | `/tmp/TestLiveScheduledslot-0TestLiveBreakGlassShimRecoveryselected-team1140847924/001` |
| default-headless-gate-stop | 05:24:38.656–05:28:05.768 | `956e8395-c0d9-4cc6-a9fd-f9a0080d7139` | `/tmp/TestLiveScheduledslot-2TestLiveCommonDefaultHeadlessGateStop1237732248/003` |
| feedback-3-cycle-escalation | 05:16:24.370–05:19:09.014 | `0db71a6d-afc3-4157-81e9-49ea2da2bf47` | `/tmp/TestLiveScheduledslot-1TestLiveCommonFeedbackThreeCycleEscalation1182985962/003` |
| filing | 05:31:51.940–05:32:39.008 | `740c238b-28db-46a4-8155-6af51939e9fc` | `/tmp/TestLiveScheduledslot-1TestLiveCommonFiling3443954526/003` |
| full-ensign-cycle | 05:28:06.984–05:30:19.884 | `e0502a24-d0ac-4b56-ba48-13ca411565e9` | `/tmp/TestLiveScheduledslot-2TestLiveCommonFullEnsignCycle3965061912/003` |
| gate-guardrail | 05:31:39.459–05:33:15.550 | `84845e33-40f5-47ff-b445-1113f26b0dfb` | `/tmp/TestLiveScheduledslot-2TestLiveCommonGateGuardrail221434753/003` |
| keep-moving-posture | 05:20:22.852–05:24:36.978 | `83809ed4-ed58-4e57-98c8-e1ca5a5fb4e2` | `/tmp/TestLiveScheduledslot-2TestLiveCommonKeepMovingPosture1451290018/003` |
| merge-hook-guardrail | 05:32:40.206–05:33:03.984 | `7e6c8493-2b83-4317-a023-13f23cc2958d` | `/tmp/TestLiveScheduledslot-1TestLiveCommonMergeHookGuardrail483471058/003` |
| merged-team-mode | 05:28:53.940–05:31:06.426 | `6abd960a-e960-4902-a206-17a6ef6101ff` | `/tmp/TestLiveScheduledslot-1TestLiveMergedTeamModeDispatch761946429/002` |
| recorded-gate-lifecycle | 05:25:43.173–05:28:53.013 | `fb257e0d-fdce-49ab-856c-11bc8d430361` | `/tmp/TestLiveScheduledslot-1TestLiveCommonRecordedGateLifecycle2992277225/004` |
| rejection-flow | 05:16:24.409–05:19:21.506 | `3e347c0e-62c7-4014-a95e-100ace3c2c0f` | `/tmp/TestLiveScheduledslot-2TestLiveCommonRejectionFlow3795307379/003` |
| self-evidence-merge-triage | 05:29:37.769–05:32:56.623 | `199d98c1-3d08-40b7-825d-2e494aba3f3a` | `/tmp/TestLiveScheduledslot-0TestLiveCommonSelfEvidenceMergeTriage3438515223/003` |
| shallow-boot | 05:32:57.780–05:33:14.496 | `51358476-103e-4ce7-8198-174bfc703231` | `/tmp/TestLiveScheduledslot-0TestLiveCommonShallowBoot1269532478/003` |
| withdrawn-gate-recovery | 05:30:50.870–05:31:50.736 | `b8ae81a9-863b-42cf-9809-9f0cecac04b5` | `/tmp/TestLiveScheduledslot-1TestLiveCommonWithdrawnGateRecovery2625847151/003` |
| zero-discovery | 05:33:05.127–05:33:17.499 | `4917e6a9-8656-49e2-82a5-036a88820c91` | `/tmp/TestLiveScheduledslot-1TestLiveCommonZeroDiscovery410443155/003` |

Eighteen root session transcripts above cover distinct scheduled functions/variants with no config/session collision. All 22 retained init streams have distinct session IDs; the two smallest-mechanism sessions intentionally share their one function's workflow cwd and run within its slot. Auto-continue has two separate fixture roots. The downloaded archive lacks project transcripts for the two auto-continue and two smallest-mechanism sessions; their init streams provide cwd/session evidence, but this review does not invent an archived config-path proof for those four. Owner-handoff has no host transcript because its fixture failed before launch. No cross-test state mutation was observed; distinct paths do not constitute a blanket sandbox/security guarantee.

## Concrete failures and classifications

1. **Owner-handoff fixture:** `initial stamped dispatch produced incomplete owner tuple` shows `Entity:""` with stage `implementation`, worker `conflict-owner-implementation`, branch `conflict-owner` and worktree `.worktrees/spacedock-ensign-conflict-owner`. Go failure occurs after 0.28s, with no host stream. Naming owner is investigating per FO; this establishes missing durable proof, not a scheduler defect.
2. **Rejection-flow:** live Claude emitted real host events, then failed `made no stream progress within 1m0s (no-progress quiet budget)`. Diagnostic identifies a filesystem-wide `grep -rn ... /` search. This is a liveness/transport completion failure, not an auth skip or XFAIL. No evidence here attributes it to queue ordering or shared config.
3. **Keep-moving:** Go passes but emits `XFAIL claude-sonnet/keep-moving-posture ... observed=[keep-moving-violation]`; retained metrics confirm `xfail`. The registered assertion ran and did not cleanly pass. Existing policy is preserved rather than concealed.

## Review-finding disposition — residual merged host activity

Observation: merged test passes at 05:30:49.709577014 and its slot admits withdrawn-gate at 05:30:49.709588896. Merged session `6abd960a-e960-4902-a206-17a6ef6101ff` nevertheless emits fresh ToolSearch, SendMessage and Bash activity through 05:31:06.426. Its project transcript line 151 at 05:30:58.424 says its workflow directory was deleted and shell cwd recovered to `/tmp`; line 156 reports its old root is absent. It then summarizes completion. With self-evidence, bare-reachable and withdrawn-gate, four distinct FO session activity spans intersect from 05:30:50.870 to 05:31:06.426 (15.556s). This is **not** a fourth occupied Go-test slot or proof of four simultaneously busy CPUs.

- Released user/normal workflow: scheduled Claude CI, merged test completes terminal state while other queued journeys continue.
- Observable harm: residual host tools query its deleted workflow and fail; no retained evidence of writing another session's state or config collision.
- Authority: `none: AC-3 explicitly introduces no new cancellation or process-cleanup guarantee; AC-2 bounds occupied tests, not descendant process count. No value-AC breach from the observed residue is established.`
- Trigger evidence: exact post-pass records and deleted-root tool result retained in `evidence.json`; source SHA-256 in manifest. Candidate uses pre-existing `cmdPoller.kill()` direct-process kill; its pre-existing merged test does not establish descendant process reaping. No new reproduction or root-cause fix was attempted.
- Worker proposal: **Deferred risk**, ownership pre-existing live-harness cleanup outside scheduling scope, decline scheduling candidate fix. Trigger is observed; a process-count/reaping promise is outside current scope. Supported slot bound and distinct config/session paths still hold.
- FO authorization: **DECLINE candidate fix; accept Deferred risk** (addressable-worker message during this review). Promote if residue writes another journey's state, config/session collision occurs, or future scope promises process-count/reaping. No process supervisor or new cleanup contract is proposed.

## Per-AC conclusion and completion

- **AC-1:** existing deterministic admission-order validation remains applicable. This run supplies no new performance threshold and no defensible provider-speed comparison. Go start ordering can differ from mutex admission ordering; no wallclock benefit is claimed.
- **AC-2:** real common/substrate overlap and distinct retained actual config/workflow/session paths are demonstrated, with exactly three occupied tests and sequential break-glass variants. **Overall closure remains incomplete** because required all-function durable success is absent (two failures and the registered keep-moving XFAIL); the four missing archived config transcripts and residual-host boundary are explicit limits. Passing overlap evidence must not be discarded merely because the suite failed, nor may it turn the suite green.
- **AC-3:** all 20 functions were invoked exactly once, original assertions exposed the two failures and XFAIL, continued scheduling completed the remaining queue, and JSON/stream/metric/project evidence survived the failed run. Codex's 17-function execution is outside this Claude-only review; previous deterministic AC-3 conclusions remain unchanged.

## Stage Report: validation — CI evidence follow-up

- DONE: Verify actual execution coverage, common/substrate host overlap, three-slot bound and runtime isolation from retained tip CI evidence.
  Exact test events plus actual transcript timestamps/init identities establish the supported facts above; unproved host coverage and cleanup limitations remain explicit.
- DONE: Record a durable AC-2 conclusion with exact sources and limitations without changing approval or rerunning tests.
  This artifact alone records the incomplete closure, source digests, evidence tables and FO-declined deferred risk; approved entity report, gate and candidate are untouched.

### Summary

The full tip run proves live scheduling overlap and bounded occupied-test concurrency, including both successful break-glass variants, while preserving isolation evidence and failed-run artifacts. It does not close AC-2's complete required proof: two functions failed, one registered semantic XFAIL remains, and actual descendant-process cleanup is neither promised nor established.
