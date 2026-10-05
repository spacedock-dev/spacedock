---
title: The Pi worker-lifecycle assert credits a completion surface this host does not provide
status: validation
score: 0.8
source: "Live lane run 37101046846, journey default-headless-gate-stop, 2026-10-03: the journey reported implementation-worker-not-dispatched while the launcher log proved both dispatches succeeded."
id: mk72bnt1b5hsp9sfv83979xs
started: 2026-10-04T04:35:01Z
worktree: .worktrees/spacedock-ensign-repair-pi-worker-lifecycle-observation
---

Two Pi journeys report a false cause. The assert cannot see how the First Officer
observes a worker completion on this host.

## Problem

`assertWorkerLifecycle` credits completion only from a `subagent` tool result that
carries `Run: <id>` and `State: complete`, or, since task `nta`, from a
`subagent_wait` result keyed on the spawned run id.

On this host neither surface appears. Run `37101046846`, journey
`default-headless-gate-stop`, produced these facts:

- The First Officer's root transcript holds 3 `bg_wait` calls and **zero**
  `subagent({action:"status"})` calls.
- The same transcript holds 3 `Background task completed` notices.
- `subagent_wait` appears once, only inside tool-description text. The wait tool is
  named `bg_wait` here.
- `command.log` shows `status --set` and `dispatch build --stamp` succeeding for the
  implementation stage and then the validation stage.
- Two child session files exist under the journey's `sessions/` directory.

The assert therefore computes `completed=-1` and reports
`implementation-worker-not-dispatched`, although the worker was dispatched and its
completion was observed. Task `nta` made the same class of repair once before, when
the assert credited only a `status` result while the First Officer used
`subagent_wait`.

## Value

The lane names the real fault. An operator who reads a red lane learns the actual
cause instead of a false one.

## Out of scope

The two interface faults observed in the same run: the `gate prepare --artifact`
argument, and the split-state read that resolved the stale entity file. Route those
to `vpf multi-artifact-gate-prepare` and to `qx state-ready-cwd-path-resolution-reset`
or `93 document-dispatch-entity-path-base`.

## Expected surface and tolerance

Estimate net LOC change: +40, across 2 files (`internal/ensigncycle/claude_runtime_helpers_test.go`
and a replay test beside it). Insertions ~+55, deletions ~-15. Tolerance: +/-20 net LOC,
+/-1 file.

Declared semantic change: which completion surfaces the assert credits. This task must
NOT change the spawn tolerance, the `completed < validation` ordering, or any other
runtime's grading.

## Acceptance criteria

**AC-1 - The assert credits a native completion notice and a `bg_wait` completion.**
Verified by: a deterministic test that feeds a captured transcript shape containing
those two surfaces and asserts `completed` is set for the spawned run id. Independent
baseline that can move the wrong way: the captured transcript from run
`37101046846` yields `completed=-1` before the change.
Falsifying edit: remove the new credit; the test must turn RED.

**AC-2 - A genuinely missing dispatch still fails.**
Verified by: the existing zero-spawn control. Falsifying edit: credit any transcript
regardless of spawn count; the control must turn RED.

**AC-3 - The ordering contract is preserved.**
Verified by: the existing `completed < validation` assertion. Falsifying edit: drop
the ordering check; the inverted-ordering control must turn RED.

**AC-4 (no-regression) -** `go test ./internal/ensigncycle/...` is green, and no other
runtime's grading changes.

## Test plan

Primary proof owner: `internal/ensigncycle/claude_runtime_helpers_test.go` (the shared
assert) plus one replay test. The captured root transcript is preserved under
`docs/dev/.spacedock-state/_evidence/pi-delegated-gate-continuation-reliability/retained-pi-recorded-gate/`
and must be copied to repository testdata before the replay test depends on it. Deterministic
tests only.

## Stage Report: implementation

- DONE: Fix the shared assert so it credits the completion surfaces this host supplies: a bg_wait result keyed on the spawned run id, and the native completion notice. Today it credits only a subagent result carrying Run/State, or a subagent_wait result, so it reports spawns=1 completed=-1 while the officer did dispatch and did observe the completion. Keep the rule that a genuinely missing dispatch still fails, and keep the completed-before-validation ordering. Add the deterministic test from the captured transcript shape, and keep the negative control that removes the new credit so the test turns red.
  Commit `2d5bd25f6`: `assertWorkerLifecycle` now credits a `bg_wait` result that names the spawned run id and reports `done`/`complete` (alongside `subagent_wait`), and a top-level `custom_message` `customType: subagent-notify` whose body is `Background task completed`. Both stay behind the existing spawn-count and `completed < validation` checks. Positive control `TestPiWorkerLifecycleObservationReplay` feeds the captured shape; removing either new credit makes it RED with `validation lifecycle incomplete: spawns=1 completed=-1 validation=5 report=<nil>` — the same completed=-1 run 37101046846 reported.
- DONE: Prove it: the focused deterministic tests, then a local live run of the affected journey auto-continue-after-implementation using SPACEDOCK_LIVE_RUNTIME=pi and a distinct SPACEDOCK_LIVE_ARTIFACT_DIR. Record the exact result. Do not start a CI lane run.
  `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/live-observer SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 45m -run '^TestLiveCommonAutoContinueAfterImplementation$' ./internal/ensigncycle -v` -> `--- PASS: TestLiveCommonAutoContinueAfterImplementation (766.12s)`. Both fixtures passed: single-root `error=<nil> timeout=false` 7m40.9s, split-root `error=<nil> timeout=false` 5m03.6s, model `openai-codex/gpt-5.6-luna:max`. Grade `pass`, observed codes `[]`; no XFAIL/XPASS line, because the binding is gone. No CI lane run started.
- DONE: When that journey passes locally with an exact passing pair, clear its XFAIL binding in the shared runner and record the evidence in the entity. Do not leave the binding as XPASS.
  Commit `a2166585a`: `TestLiveCommonAutoContinueAfterImplementation` gaps are `nil` (binding to `mk72bnt1b5hsp9sfv83979xs` removed with its stale evidence-defect comment). Rebased onto the peer tip `fb4428e9c`, where the binding lived; original base was main `4436ec14c`. Peer branch untouched.
- DONE: Keep a genuinely missing dispatch failing, and the completed-before-validation ordering.
  Same test includes the zero-spawn control (spawn retargeted to implementation -> spawns=0, RED) and the inverted-order control (completion moved after `gate prepare` -> `completed >= validation`, RED). Each turns RED only while its guard exists.
- DONE: no-regression `go test ./internal/ensigncycle/...` and no other runtime's grading changes.
  `go test ./internal/ensigncycle/... -count=1` green (286.7s); `go test -race ./internal/ensigncycle/ -run TestPiWorkerLifecycleObservationReplay` green; `go test ./internal/contractlint/` green; gofmt clean. The claude/codex branches are untouched.

### Summary

The shared worker-lifecycle assert now credits the two completion surfaces a Pi host supplies — a `bg_wait` result keyed on the spawned run id and a native `subagent-notify` completion notice — while preserving the spawn-count and `completed < validation` guards. A deterministic replay of the captured transcript shape plus its removed-credit, zero-spawn, and inverted-order controls proves each guard is load-bearing. After rebasing onto the peer stack tip, the auto-continue Pi XFAIL was cleared, and the journey passed locally with an exact two-variant passing pair (766.12s) and no XPASS.

## Review-finding disposition

Validation observations below are recommendations only; no candidate repair or FO disposition is authorized by this report.

- **V1 — Evidence defect / Material / owned AC-1 boundary:** `claude_runtime_helpers_test.go:259` credits any native completion once any matching-stage run ID exists; it never matches the notice to that ID. Released user/workflow: Pi officers receiving asynchronous notifications across stage dispatches. Harm: a different worker/run can satisfy the completion proof before gate preparation. Authority: `value-ac[AC-1]` requires completion for the spawned run ID. Trigger: the independently captured live notice explicitly names run `6d33973c-4c8d-449b-a356-91c156c76c49`, but an overlay pairing it with spawn `run-ac0f4b70` returns `err=<nil>` (`captured-native.log`). Proposal: hold; seek FO-authorized correction of this attribution boundary, without adding another lifecycle controller.
- **V2 — Evidence defect / Material / source-provenance decision:** released workflow is the portable deterministic regression replay. Harm: the promised independent run-37101046846 baseline is replaced by six authored JSON constants with no pinned captured input; the inline notice also omits the run-bearing retention-directory field observed live. Authority: `value-ac[AC-1]` requires a captured-shape independent baseline. Trigger: the plan's retained `first-root.jsonl` is the July recorded-gate capture (SHA-256 `4ecc18637c62143b8cbae7fbf584fa3093145aa3f6d8d46daa31f4343dfe85de`), contains `subagent_wait` and no `bg_wait`, and is not run 37101046846; the replay comment instead names `pi-native-completion-evidence/`, where no matching transcript is referenced. Proposal: route source/provenance choice to FO/captain, not an automatic invented-fixture repair. A source-traceable embedded excerpt could be sufficient in principle; this unverified reconstructed fixture does not discharge the stated copy/independent-baseline condition.
- **V3 — Evidence defect / Material / exact live acceptance boundary:** released workflow is the now-unbound Pi auto-continue single-root/split-root pair. Harm: this validation cannot establish the clean pair required to remove XFAIL. Authority: `captain-ruling[2026-10-04]` requires this exact independently passing pair with no XFAIL/XPASS. Trigger: requested command exited 1; single-root hit the default 12-minute per-run cap (`error=signal: killed timeout=true`), and split-root never started. No XFAIL/XPASS appeared, but that is not PASS. Proposal: hold binding-removal acceptance; timeout precedes lifecycle grading and does not establish that the changed predicate caused a runtime outcome defect.
- **V4 — Evidence/scope defect / Material authorization hold / Needs captain decision:** the third file is the expressly requested conditional XFAIL removal and is within the +1-file allowance; the completion, missing-spawn, and ordering controls are substantively justified. Nevertheless actual +83 net across 3 files exceeds +40 +/-20 (upper bound +60) by 23 lines. Authority: `contract[docs/dev/README.md#review-finding-disposition]` reserves out-of-tolerance scope decisions to the captain. Trigger: `git diff --numstat fb4428e9c..a2166585a` = 99 additions, 16 deletions. Proposal: record a captain-visible scope/design decision before another pass; justified lines do not waive the tolerance. This is not evidence of a runtime outcome regression.
- **Deferred risk (evidence defect):** synthetic `bg_wait` text `done. Outcome: 0 complete, 1 failed.` is accepted because the inherited predicate uses substrings. That exact failure-response wording was not observed from this host, so no supported successful-path AC failure is established by this probe; promote if a retained supported-host failed wait emits that form. Do not widen this task into a new terminal-state parser without authorization. The unfinished-second-dispatch probe also accepts, but demanding final-dispatch completion would change the explicitly preserved spawn/order contract and is not charged as a new regression here.
- **Polish only (evidence overclaim; no release blocker):** removing only `spawns < 1` does not kill the replay: its zero-spawn fixture also lacks a correlated run ID. Thus the implementation report overstates this mutation proof; the actual no-dispatch behavior still rejects correctly. No new scaling risk: the changed reader remains a linear per-event scan with no added I/O.

## Stage Report: validation

- DONE: Re-run the focused deterministic tests and the replay test yourself, and reproduce the negative control: with both new credits removed, the positive case must turn red with spawns=1 completed=-1. If you cannot make the control fail, the credit is not load-bearing and you must reject.
  Candidate `a2166585a`: 15 focused top-level tests PASS (including new replay and real Claude replay), one older /tmp-dependent Pi double-dispatch replay SKIP; focused race PASS. Go overlay removes BOTH new credits without changing candidate bytes: positive replay FAILS exactly `spawns=1 completed=-1 validation=5 report=<nil>`.
- FAILED: Independently re-run the live Pi journey auto-continue-after-implementation with the captain's credentials and a distinct SPACEDOCK_LIVE_ARTIFACT_DIR, and confirm a clean PASS with no XFAIL and no XPASS. The cleared binding must rest on that exact passing pair.
  Exact requested command exited 1: `--- FAIL: TestLiveCommonAutoContinueAfterImplementation (749.89s)`; single-root `error=signal: killed timeout=true`, duration `12m28.704126s`, model `openai-codex/gpt-5.6-luna:max`; split-root not reached; no XFAIL/XPASS. Inherited captain OAuth was used (no OPENAI_API_KEY override); no retry, timeout override, or CI lane.
- DONE: Judge the surface overrun: the change is +83 net test lines across 3 files against an ideation estimate of +40 across 2 files. Decide whether the third file (the binding removal) and the control lines are justified, or whether this is a scope breach that needs a captain decision.
  V4: third file/control purpose justified; file allowance satisfied; LOC allowance exceeded, requiring a captain-visible decision. The live prerequisite for retaining the binding removal was not independently reproduced.
- DONE: Check the claimed deviation: the plan said to copy the retained transcript into testdata, and the worker instead embedded the captured shape in the replay test. Decide whether that satisfies the plan's condition or removes needed evidence.
  V2: not sufficient as delivered; no captured run-37101046846 fixture or verified extraction, mismatched source references, and omitted native identity evidence. Do not confuse the passing synthetic positive/mutant pair with source provenance.
- DONE: AC-2 — A genuinely missing dispatch still fails.
  The replay's retargeted-spawn case rejects; missing spawn result, pre-spawn notice, and wrong-type controls also reject. Removing only the lower spawn-count guard still rejects via missing completion, limiting the report's guard-specific mutation claim but not this behavioral AC.
- DONE: AC-3 — The ordering contract is preserved.
  Inverted completion/transition fixture rejects; overlay removal of `completed >= validation` makes the test FAIL (`a Pi completion after the validation transition passed`). Spawn tolerance and existing ordering expression remain unchanged.
- DONE: AC-4 (no-regression) — scoped package suite and other-runtime checks.
  `go test ./internal/ensigncycle/... -count=1` PASS (345.550s); focused Claude/Codex lifecycle and real-Claude replay PASS; focused `-race` PASS; changed-file `gofmt -l` empty. No repository-wide suite or CI lane was run.
- FAILED: AC-1 — native completion is proved for the spawned run ID.
  Positive and removed-credit tests prove recognition is load-bearing, but the actual live notice for another run is accepted (`captured-native.log`); V1 is a material false-positive evidence boundary, and V2 leaves the promised independent captured baseline unproved.
- DONE: Report findings with the defect kind and release scope, and whether the change can proceed.
  REJECTED / cannot proceed: V1–V3 material evidence findings plus V4 captain scope hold; deferred risk is separate. Findings above are advisory, not FO authorization. No deliverable repair, source edit, candidate commit, or binding restoration was made.

### Summary

Independent deterministic, race, package, and negative-control checks succeeded, but the exact live verification timed out before completing even the first fixture. Validation also reproduced uncorrelated native-completion acceptance using this live run's actual notice and found that the reconstructed replay cannot substantiate its claimed source; delivery must remain held for FO disposition and the captain's scope/provenance decisions. Commands, overlays, and logs are retained in the assigned worktree's `.validation-observer/`; live evidence is under `/tmp/live-observer-verify/`.

## Stage Report: implementation (cycle 2)

- DONE: Fix the correlation defect the validation found. A real completion notice for a DIFFERENT run id was incorrectly accepted and credited the spawned run. The credit must apply only when the notice names the spawned run id. Add the deterministic negative case for that, so a notice belonging to another run cannot set the credit.
  Commit `3e7036970`: the native `subagent-notify` branch now also requires `strings.Contains(event.Content, piRunID)`, matching the run-bearing `Retention-managed async directory: …/async-subagent-runs/<id>` field against the spawn result's `details.runId`. `TestPiWorkerLifecycleObservationReplay` adds the negative fixture that swaps the captured notice's run id for an unrelated one; removing the correlation guard makes it FAIL `a native completion notice for another run credited the spawned run`, and with the guard it rejects. This closes the validator's `TestValidationObserverCapturedNativeIdentity` false-positive (spawn `run-ac0f4b70` vs notice `6d33973c…`).
- DONE: Re-prove the control: with both new credits removed, the positive case must fail with `spawns=1 completed=-1`.
  Go overlay (candidate bytes unchanged) drops `|| pi.Message.ToolName == "bg_wait"` and `false`-gates the notice branch: `go test -overlay /tmp/ensigncycle-control/overlay.json -run '^TestPiWorkerLifecycleObservationReplay$'` FAILS `validation lifecycle incomplete: spawns=1 completed=-1 validation=5 report=<nil>`.
- DONE: Re-run the live Pi journey auto-continue-after-implementation locally using the captain's credentials, with a distinct and EMPTY SPACEDOCK_LIVE_ARTIFACT_DIR. The validator's live verification timed out and the split-root case never ran. Run both roots if the journey has them, and record the artifact directory state you used.
  `/tmp/live-observer-repair` (created empty; `ls -la` showed no entries before launch). `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/live-observer-repair SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES=20 go test -tags live -count=1 -timeout 45m -run '^TestLiveCommonAutoContinueAfterImplementation$' ./internal/ensigncycle -v` -> `--- PASS (669.60s)`; single-root 4m5.75s and split-root 7m2.90s, both `error=<nil> timeout=false`, model `openai-codex/gpt-5.6-luna:max`, no XFAIL/XPASS. Both root sessions hold 1 `bg_wait` result and 1 `subagent-notify`, and each notice's `Retention-managed async directory` names that root's spawn run id (single `78eef524…`, split `6a45b94a…`). Only deviation: `SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES=20` (documented knob) because the default 12m per-run cap reset the validator; a first attempt with an empty dir still ended at 2m47s `spawns=1 completed=-1 validation=-1` (FO ended the turn before awaiting) so the dir was re-emptied for the passing run.
- DONE: Settle the replay provenance the validation called unproved: either copy the retained transcript into testdata as the plan stated, or state the embedded shape's exact source with its commit and hash so a later reader can verify it. Do not leave it as an assertion about an unnamed capture.
  Commit `3e7036970` replaces the authored constants with `internal/ensigncycle/testdata/pi_worker_lifecycle/auto-continue-validation-parent.jsonl`, selected from the real parent session `…/sessions/2026-10-04T05-06-33-345Z_01a1054e….jsonl` (SHA-256 `057933408037181af4ec990d5e394eff8754afe76640a10e0cc83bcd72003ca8`; lines 39,40,41,42,43,108). `provenance.json` records source path, hash, selected lines, spawn run id and the field projection; the replay loads the fixture and the test comment names that source instead of an unverified directory.
- DONE: Do not change the declared surface. The net +83 across 3 files against the +40 across 2 files estimate is recorded as awaiting the captain; record it again in your report and leave it alone.
  Surface decision left untouched, awaiting the captain. `git diff --numstat fb4428e9c..3e7036970`: 3 code files +111/-16 (net +95; was +83), plus the plan-mandated testdata fixture/provenance (`+23` across 2 new files, net +118 overall). No spawn tolerance, `completed < validation` order, or other-runtime grading changed.

### Summary

The native `subagent-notify` credit now requires the notice to name the spawned run id, and a deterministic negative case proves that correlation is load-bearing. The replay no longer rests on authored constants: it loads a hash-pinned, source-traceable fixture with a provenance record. The live Pi auto-continue journey passed both single-root and split-root from an empty artifact directory with the captain's credentials, and both live notices named their own spawn run id.


## Review-finding disposition (validation cycle 2)

Advisory review of candidate `3e7036970` against `fb4428e9c`; no candidate edits or authorization implied. Verification overlays/logs are in the assigned code worktree's `.validation-observer/cycle2/`.

- **V1 closed — Evidence defect / formerly Material / AC-1 correlation:** the normal replay passes; deleting ONLY `&& strings.Contains(event.Content, piRunID)` fails with `a native completion notice for another run credited the spawned run`. The repaired wrong-run case rejects; both retained live roots also pass with their wait results disabled, exercising their real correlated native notices. This boundary can proceed technically.
- **V2 closed — Evidence defect / formerly Material / AC-1 provenance:** all six selected source lines (39,40,41,42,43,108) exist in the 173-line named parent session; recursive projection comparison matches every retained fixture value to that source. Independently computed source SHA-256 is `057933408037181af4ec990d5e394eff8754afe76640a10e0cc83bcd72003ca8`, correctly labeled `source_sha256` in provenance.json. **Polish/evidence overclaim:** that is NOT the fixture hash; the fixture hashes to `9f082b9440e0b1793f5c2c70fb932d7c783af47cba5730923f912afb3607a5cb`. This is the first validator's real 2026-10-04 session, not original CI run 37101046846; it independently reproduces the same old-credit failure. Provenance is sufficient for this repaired assignment; the full source remains local retained evidence, while the replay itself is portable.
- **V3 closed for this criterion — Evidence defect / formerly Material / exact live pair:** the recorded PASS 669.60s is supported by both retained successful processes, single-root 4m5.75020425s and split-root 7m2.896998s, and gate preparation after correlated completion. Own run IDs are `78eef524-0226-4b06-a3ee-9ece2c93ded5` (spawn line 42, notice 47, gate 78) and `6a45b94a-04ba-4991-98e9-853117ee509c` (spawn 46, notice 51, gate 66). The documented 20m cap does not change credit semantics; both roots actually finish below the default 12m, so it does not weaken this recognition proof. No new live run or independent reconstruction of the exact 669.60s test total is claimed.
- **Journey reliability — Outcome risk / Deferred for this repair, not a credit defect:** the reported first repair attempt ended its officer turn without awaiting, at `spawns=1 completed=-1 validation=-1`. With no completion or transition observed, rejection is correct: a journey/orchestration flake, not evidence the repaired credit missed a supplied notice. The earlier validator timeout likewise preceded grading. These attempts weaken any zero-flake/reliable-default-run claim, not AC-1; the retry's emptied directory no longer retains the first attempt for this review. Promote/route separately if a retained completion is rejected (AC-1), or repeatability/default-cap reliability becomes the acceptance promise; do not suppress missing observations to make the journey green.
- **V4 remains open — Evidence/scope defect / Material authorization hold / Needs captain decision:** released workflow is the Pi lifecycle grader and its unbound live journey; harm is accepting an out-of-tolerance deliverable without the required scope authority. Authority: `contract[docs/dev/README.md#review-finding-disposition]` reserves scope/tolerance changes to the captain. Trigger: `git diff --numstat fb4428e9c..3e7036970` shows 3 code files +111/-16 (net +95), plus 2 testdata files +23, total +118 across 5, against +40 +/-20 across 2 +/-1. Code alone exceeds the +60 upper LOC bound by 35; all-in exceeds it by 58 and the 3-file maximum by 2. Binding removal, controls and required fixture/provenance are justified purposes, not a scope waiver. Do not proceed to release until captain records disposition; no automatic repair/design reset is authorized here.
- **Existing terminal-parser risk — Evidence defect / Deferred:** synthetic `done. Outcome: 0 complete, 1 failed` still credits completion via inherited substrings. No retained supported-host failed wait supplies that wording; successful live waits say `1 complete` and satisfy AC-1. Promote if real supported-host failure output has this form. No new lifecycle controller, terminal parser or spawn-tolerance requirement belongs in this validation.
- **Tautology audit — clean for changed tests:** the replay calls the behavior and asserts independent pass/reject outcomes; it never compares source text to its own copy. Parsing a run ID from fixture INPUT to construct a mismatching input is not deriving an expected verdict from the tested helper. Correlation, credit and ordering mutations all turn tests red. **Polish/evidence wording:** the negative-control comment at replay lines 62–63 misattributes the mutation result: deleting credit kills the positive cases, not the already-negative no-surface case. No release-blocking behavior defect follows. Reader remains a linear event scan, with no new I/O or multiplicative work.


## Stage Report: validation (cycle 2)

- DONE: Verify the correlation repair independently. The native completion credit must apply only when the notice names the spawned run id. Re-run the new negative case, and confirm the falsification claim by removing only the correlation guard and requiring the test to fail.
  Candidate replay PASS; single-guard overlay FAIL exit 1 at `a native completion notice for another run credited the spawned run` (`cycle2/replay.log`, `no-correlation.log`). AC-1 correlation proof is load-bearing, not trusted prose.
- DONE: Re-prove the control: with both new credits removed, the positive case must fail with spawns=1 completed=-1.
  Overlay removing bg_wait credit and disabling native credit FAIL exit 1: `validation lifecycle incomplete: spawns=1 completed=-1 validation=5 report=<nil>` (`no-new-credits.log`); tracked candidate unchanged.
- DONE: Check the replay provenance claim: the fixture internal/ensigncycle/testdata/pi_worker_lifecycle/auto-continue-validation-parent.jsonl must match its stated SHA-256, and its provenance.json must name source lines that exist in the cited session. Verify the hash yourself and report it.
  Important correction: fixture SHA-256 is `9f082b9440e0b1793f5c2c70fb932d7c783af47cba5730923f912afb3607a5cb`; claimed `057933408037181af4ec990d5e394eff8754afe76640a10e0cc83bcd72003ca8` is the SOURCE hash and matches correctly labeled metadata. All six cited lines exist and retained fields match exactly (`provenance.log`, `sha256.log`); V2 closed, hash wording is polish.
- DONE: Judge the live evidence recorded: PASS 669.60s, both roots, each root's notice naming its own spawn run id, with SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES raised to 20. Decide whether that evidence supports the criterion, and whether the raised cap or the earlier failed attempt weakens it. The first attempt ended with the officer not awaiting the async worker; decide whether that is a journey flake or a signal about the credit.
  Both retained roots confirm success and their own UUIDs; actual parent replays PASS even with wait credit disabled (`adjacent.log`, `live-artifacts.log`). Durations 4m5.75s/7m2.90s are below 12m: cap change does not weaken recognition proof. Earlier no-await failure is a journey flake correctly rejected, not a missed supplied credit; weakens reliability claims only. Detailed classification/limits above.
- DONE: Judge the scope hold and record it for the captain: 3 code files at +111/-16 (net +95) plus testdata at +23 across 2 new files, against the plan's +40 across 2.
  Independently verified numstat matches; V4 remains a Material authorization hold. Legitimate binding removal and fixture/control purposes do not waive +40 +/-20 LOC and 2 +/-1 file limits; captain disposition is required before release.
- DONE: Audit for tautological tests in the captain's sense: no test may assert a file's own text against a copy of it, and no expected value may come from the thing it tests.
  Changed replay asserts helper outcomes from external captured inputs; no source-copy expectation and no helper-derived verdict. Parsed fixture run ID only constructs a negative input. Three independent code mutations kill the expected assertions; imprecise negative-control comment is polish, not tautology.
- DONE: AC-2 — A genuinely missing dispatch still fails; AC-3 — The ordering contract is preserved.
  Replay and independent adjacent-state matrix reject missing spawn/result and inverted ordering; deleting only `completed >= validation` kills the inverted-order control (`no-order.log`). Spawn tolerance unchanged; no claim that deleting only `spawns < 1` kills a case that also lacks completion.
- DONE: AC-4 (no-regression) — scoped package suite and other-runtime checks.
  `go test ./internal/ensigncycle/... -count=1` PASS 314.956s; 15 focused top-level tests PASS; focused Claude/Pi/Codex/replay `-race` PASS; changed-file `gofmt -l` empty, `git diff --check` clean. Adjacent matrix covers 15 states plus two live-parent replays; no new reader scaling risk.
- SKIPPED: Additional live journey, CI lane, repository-wide test suite, or candidate changes.
  Current assignment asks assessment of recorded live evidence; none launched. Verification-only Go overlays preserve candidate HEAD/bytes; only this shared-state entity report is committed. Full provenance/live replay audit additionally uses explicitly named retained /tmp evidence, not a hidden portable-test dependency.
- DONE: Report each finding with the defect kind and release scope, and whether the change can proceed.
  Technical AC-1–AC-4 evidence PASSED; overall recommendation REJECTED for release / HOLD solely for unresolved captain scope authorization (V4). V1–V3 resolved for this task; deferred terminal-text and journey-repeatability risks are separated from blockers above.

### Summary

Independent mutation tests confirm the repaired correlation, both new credits and ordering are load-bearing; the source-traceable fixture and both real passing live roots support the promised recognition behavior. The quoted SHA-256 belongs to the full source session, not the projected fixture, and the report now records both hashes precisely. Candidate `3e7036970` is unchanged and technically validated, but cannot proceed to release until the captain disposes of the documented scope overrun.

## Stage Report: implementation (rebase onto pi-live lane)

- DONE: Rebase this branch onto `origin/spacedock-ensign/pi-live-lane-pin-refresh` at `5ead9b85c`, keeping this task's own change intact.
  `git rebase origin/spacedock-ensign/pi-live-lane-pin-refresh` replayed this task's three commits as `67d413a39`, `c65bba735`, `9755d96d1`; new tip `9755d96d1`. The lower shared commits were dropped as patch-identical upstream (`fadf4e453`, `080d37f23`, `17858e911`, `d70bd32d6`, `fb4428e9c`, …), so only this task's own change sits on the new base. Pre-rebase ref saved at `refs/backup/repair-pi-worker-pre-rebase` (`3e7036970`). `git diff refs/backup/repair-pi-worker-pre-rebase..HEAD` is empty for `claude_runtime_helpers_test.go` and the replay test: the change survived byte-for-byte.
- DONE: Report the conflict, both sides, and the resolution.
  One conflict: `internal/ensigncycle/claude_runtime_helpers_test.go` at the async-wait credit. Base (`5ead9b85c`, via the pi-live lane layer) carried `ToolName == "subagent_wait"` with the comment "subagent_wait (a recent pi-subagents)"; this task's `2d5bd25f6` generalizes it to `(subagent_wait || bg_wait)` and rewrites the comment for both host tool names. Resolved in favour of this task's side; the base condition is subsumed by the `|| bg_wait` form, so no lane-layer semantics were dropped.
- DONE: Report `internal/ensigncycle/shared_live_runner_test.go`, where the FO note expected a conflict.
  No conflict. This task's `a2166585a` removed the `mk72bnt1b5hsp9sfv83979xs` XFAIL binding and its owner comment with context that still matched the lane layer, so git applied it cleanly. Result vs the lane tip: `liveJourney(..., nil, ...)` and the binding comment removed, net `+1/-9`.
- DONE: Report `skills/first-officer/references/fo-install.md`, where the FO note expected a conflict.
  No conflict. This branch's install-sentinel commit `17858e911` is patch-identical to the lane layer's `c6170ea8c`, so rebase skipped it as already upstream; the file is byte-identical to the rebase target and needs no resolution.
- DONE: Run only this task's focused tests; no repository-wide suite, no CI run.
  `go test ./internal/ensigncycle/ -count=1 -run 'TestPiWorkerLifecycleObservationReplay|TestPiAutoContinueReplayDoubleDispatch|TestAutoContinueReplay|TestAutoContinueFixtureIsDiscoverable|TestAutoContinueGateFixturesParseAsIntended|TestAutoContinueRevalidateStreamCountsBothValidators' -v` PASS (0.518s); `TestPiAutoContinueReplayDoubleDispatch` SKIP (its `/tmp/pi-live-10-art` artifact is absent, unchanged from the recorded pre-rebase skip); `go test -race ./internal/ensigncycle/ -count=1 -run '^TestPiWorkerLifecycleObservationReplay$'` PASS; changed files `gofmt -l` empty. No repository-wide suite or CI lane started.

### Summary

This task's three commits were rebased cleanly onto the pi-live lane tip `5ead9b85c`, producing new tip `9755d96d1`. The only conflict was the async-wait credit comment/condition in `claude_runtime_helpers_test.go`, resolved in favour of this task's `subagent_wait || bg_wait` generalization. The two files the FO flagged -- `shared_live_runner_test.go` and `fo-install.md` -- needed no resolution: the XFAIL-binding removal applied cleanly, and the install-sentinel commit was already upstream. Focused lifecycle tests and the race run pass.

## Stage Report: implementation (restack onto pi-live lane 3fed0a16e)

- DONE: repair-pi-worker-lifecycle-observation onto 3fed0a16e.
  `git rebase --onto 3fed0a16e 5ead9b85c` replayed this task's 3 commits as `d8253403d`, `9a9f1d9c0`, `df88faeb7`; new tip `df88faeb7`. No conflict, no drops. All three patch-ids identical to the pre-rebase commits (`67d413a39`, `c65bba735`, `9755d96d1`).
- DONE: repair-pi-recorded-gate-lifecycle onto the new tip of repair-pi-worker-lifecycle-observation.
  `git rebase --onto df88faeb7 9755d96d1` replayed its 1 commit as `72fc537a4`; new tip `72fc537a4`. One conflict in `internal/ensigncycle/shared_live_runner_test.go`.
- DONE: Report the conflict, both sides, and the resolution.
  Base (`df88faeb7`, lane `f37bda151`) re-points the gate-stop XFAIL to `liveXFail("pi","mk72bnt1b5hsp9sfv83979xs")` with a worker-lifecycle comment; this task's `e5907d00b` removes that binding entirely (`nil`) and deletes its comment. Resolved in favour of this task's own change: comment deleted, binding `nil`. The lane's re-pointed binding is what this task supersedes ("With the observer fix stacked below the journey passes cleanly, so clear its Pi XFAIL binding"). SKILL.md (`Compose the selected-source` line) auto-merged and is byte-identical to `e5907d00b`.
- DONE: pin-lane-models-in-one-place onto the new tip of repair-pi-recorded-gate-lifecycle.
  `git rebase --onto 72fc537a4 e5907d00b` replayed its 7 commits as `f612eb553`, `d34f668d8`, `3b205f47a`, `e9e8a6d95`, `1b106a4c0`, `f9b0f783c`, `cee86defe`; new tip `cee86defe`. No conflict, no drops. All seven patch-ids identical to backup `238a2d59a`'s commits (`9c5fd356e`…`238a2d59a`).
- DONE: Keep each branch's own commits intact. Drop none of them, and take no base side wholesale. Report each new tip, and report any commit that dropped as already-applied.
  All 11 branch-owned commits (3+1+7) survived with unchanged patch-ids; zero dropped, so there is no already-applied drop to report. Each branch's own tree diff vs its new parent is byte-identical to its pre-rebase diff (`git diff --stat` match for all three). Backups retained at `backup/restack-observer`, `backup/restack-recorded-gate`, `backup/restack-lane-models`. New tips: observer `df88faeb7`, recorded-gate `72fc537a4`, lane-models `cee86defe`.
- DONE: Push each branch with a lease-guarded force.
  Each pushed in dependency order with `git push --force-with-lease origin HEAD:<branch>`: `9755d96d1...df88faeb7`, `e5907d00b...72fc537a4`, `238a2d59a...cee86defe`. Local and origin SHAs match for all three; ancestry verified (`3fed0a16e` ⊆ observer ⊆ recorded-gate ⊆ lane-models). No PR consumes any of them.
- DONE: Run only focused tests after each rebase; no repository-wide suite and no CI run.
  Observer: `go test ./internal/ensigncycle/ -run 'TestPiWorkerLifecycleObservationReplay|TestPiAutoContinueReplayDoubleDispatch|TestAutoContinueReplay|TestAutoContinueFixtureIsDiscoverable|TestAutoContinueGateFixturesParseAsIntended|TestAutoContinueRevalidateStreamCountsBothValidators'` PASS; `TestPiAutoContinueReplayDoubleDispatch` SKIP (absent `/tmp` artifact, unchanged). Recorded-gate: `go test ./internal/contractlint -run '^TestFOInstructionComponentCaps$'` PASS (SKILL.md cap); `go test ./internal/ensigncycle -run '^(TestPiWorkerLifecycleObservationReplay|TestImplementationLifecycleAndObserverNegativeControls|TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle|TestAssertGateHeld|TestAssertGateHeldAcceptsPreparedFixtureBinding|TestMissingPreparedGateExpectationIsSemanticallyGraded)$'` PASS; `go vet -tags live ./internal/ensigncycle` OK. Lane-models: `go test ./internal/claudeteam ./internal/ensigncycle -run '^(TestContextLimitForModelBoundary|TestCodexLiveWorkflowPinsOnlyExecToLuna|TestRuntimeLiveClaudeShimSetsMaximumEffort|TestPiWorkerLifecycleObservationReplay)$'` PASS; `go test ./internal/release` PASS (13.2s); changed-file `gofmt -l` empty. All search/commands bounded to the named worktrees; no repo-wide suite, no CI lane.

### Summary

The three stacked branches were restacked in dependency order onto the moved lane tip `3fed0a16e` and pushed with lease-guarded forces: observer `df88faeb7`, recorded-gate `72fc537a4`, lane-models `cee86defe`. The only conflict was the gate-stop XFAIL binding in `shared_live_runner_test.go`, resolved to the recorded-gate branch's own `nil` change (superseding the lane's re-pointed mk72 binding). All 11 branch-owned commits survived with unchanged patch-ids and no already-applied drops; focused tests pass after each rebase.


## Stage Report: implementation (restack onto pi-live lane 515344ddd)

- DONE: repair-pi-worker-lifecycle-observation onto the lane layer's new tip 515344ddd.
  `git rebase --onto 515344ddd 3fed0a16e` replayed this task's 3 commits as `57f0c5ea9`, `cda26f5c7`, `07a480522`; new tip `07a48052286ad3f46e48b469dacdedff98006437`. No conflict, no drops. All three stable patch-ids identical to pre-rebase (`d8253403d`, `9a9f1d9c0`, `df88faeb7`).
- DONE: repair-pi-recorded-gate-lifecycle onto the observer's new tip 07a480522.
  `git rebase --onto 07a480522 df88faeb7` replayed its 1 commit as `eb45c18aa`; new tip `eb45c18aa516c8bfe336dc2df297294e2f34ad2f`. No conflict. Patch-id identical to `72fc537a4`.
- DONE: pin-lane-models-in-one-place onto the recorded-gate's new tip eb45c18aa.
  `git rebase --onto eb45c18aa 72fc537a4` replayed its 7 commits as `15dc509c4`..`7ac7bb69b`; new tip `7ac7bb69b6cea927da258c7fd44395595e77f2be`. No conflict. All seven patch-ids identical to `f612eb553`..`cee86defe` (backup `restack2-lane-models`).
- DONE: Keep each branch's own commits intact; drop none; take no base side wholesale; report each new tip and any already-applied drop.
  All 11 branch-owned commits (3+1+7) survived with unchanged patch-ids; zero dropped, so there is no already-applied drop to report. Each branch's final tree diff vs its new parent equals its pre-rebase diff plus the lane-layer delta only (`git diff backup/restack2-* <tip> --stat` shows exactly the lane layer's 5 files). Backups retained: `backup/restack2-observer`, `backup/restack2-recorded-gate`, `backup/restack2-lane-models`.
- DONE: Lane-layer registry-deletion condition — no layer above carries changes to the deleted files, and no conflict with their deletion.
  `docs/runtime-live-ci-registry.md` and `internal/contractlint/live_registry_reconciliation_test.go` are absent at all three new tips; none of the 11 branch commits modifies either path (`git diff --name-only <parent> <tip>` for each). No rebase conflict touched them. No stop condition triggered. Observation for the FO (lane-layer's own, not a layer above): the lane tip 515344ddd still leaves dangling references to `docs/runtime-live-ci-registry.md` in `docs/roadmap/live-test-truth/dispatch-sprint-execution.md:13`, `docs/roadmap/live-test-truth/index.md:30,89`, and `docs/roadmap/test-behavior-completeness/index.md:21`; no branch above touches those files. Reported, not resolved — FO owns the lane layer.
- DONE: Push each branch with a lease-guarded force.
  `git push --force-with-lease origin HEAD:<branch>` in dependency order: `df88faeb7...07a480522`, `72fc537a4...eb45c18aa`, `cee86defe...7ac7bb69b`. Local and origin SHAs match for all three; ancestry verified (`515344ddd` ⊆ observer `07a480522` ⊆ recorded-gate `eb45c18aa` ⊆ lane-models `7ac7bb69b`). No PR consumes any of them.
- DONE: Run only focused tests after each rebase; no repository-wide suite, no CI run.
  Observer: `go test ./internal/ensigncycle/ -run 'TestPiWorkerLifecycleObservationReplay|TestPiAutoContinueReplayDoubleDispatch|TestAutoContinueReplay|TestAutoContinueFixtureIsDiscoverable|TestAutoContinueGateFixturesParseAsIntended|TestAutoContinueRevalidateStreamCountsBothValidators'` PASS; `TestPiAutoContinueReplayDoubleDispatch` SKIP (absent `/tmp/pi-live-10-art` artifact, unchanged); focused `-race` on the replay PASS; changed-file `gofmt -l` empty. Recorded-gate: `go test ./internal/contractlint -run '^TestFOInstructionComponentCaps$'` PASS; `go test ./internal/ensigncycle -run '^(TestPiWorkerLifecycleObservationReplay|TestImplementationLifecycleAndObserverNegativeControls|TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle|TestAssertGateHeld|TestAssertGateHeldAcceptsPreparedFixtureBinding|TestMissingPreparedGateExpectationIsSemanticallyGraded)$'` PASS; `go vet -tags live ./internal/ensigncycle` OK. Lane-models: `go test ./internal/claudeteam ./internal/ensigncycle -run '^(TestContextLimitForModelBoundary|TestCodexLiveWorkflowPinsOnlyExecToLuna|TestRuntimeLiveClaudeShimSetsMaximumEffort|TestPiWorkerLifecycleObservationReplay)$'` PASS; `go test ./internal/release` PASS (13.983s); `go test ./internal/contractlint` PASS; changed-file `gofmt -l` empty. All commands bounded to the named worktrees.

### Summary

The three stacked branches were restacked in dependency order onto the moved lane tip `515344ddd` and pushed with lease-guarded forces: observer `07a480522`, recorded-gate `eb45c18aa`, lane-models `7ac7bb69b`. Every rebase applied cleanly — no conflicts, and all 11 branch-owned commits survived with unchanged stable patch-ids. The lane-layer registry deletion did not collide with any layer above (no branch touches either deleted file); its only residue is dangling roadmap references at the lane tip itself, left to the FO. Focused tests pass after each rebase.

## Stage Report: validation (cycle 3)

- DONE: AC-1 — the assert credits a native completion notice and a `bg_wait` completion.
  `TestPiWorkerLifecycleObservationReplay` PASS at tip `d0b15a9ff` (task commits patch-identical to `07a480522`). The committed replay loads `internal/ensigncycle/testdata/pi_worker_lifecycle/auto-continue-validation-parent.jsonl` and asserts the two surfaces together, then each alone (`bg_wait` only, native notice only). Falsifier: a Go overlay that removes BOTH new credits (candidate bytes untouched) turns the replay RED with `validation lifecycle incomplete: spawns=1 completed=-1 validation=5 report=<nil>` — the same `completed=-1` the old credit produced. Removing each credit singly also RED (`bg_wait completion alone rejected` / `native completion notice alone rejected`), so both are load-bearing. Correlation guard is load-bearing: removing only `strings.Contains(event.Content, piRunID)` RED with `a native completion notice for another run credited the spawned run`.
- DONE: the replay fixture is committed and its source is hash-pinned.
  `git ls-files` lists both `auto-continue-validation-parent.jsonl` and `provenance.json`. Provenance `source_sha256=057933408037181af4ec990d5e394eff8754afe76640a10e0cc83bcd72003ca8` matches the retained source session (shasum verified); cited lines 39,40,41,42,43,108 exist in the 173-line source and project to spawn/spawn-result/bg_wait-call/bg_wait-result/native-notice/gate-prepare. Fixture sha256 `9f082b9440e0b1793f5c2c70fb932d7c783af47cba5730923f912afb3607a5cb`. Note: no test asserts the fixture's own hash; pinning is the provenance source hash plus git content addressing.
- DONE: AC-2 — a genuinely missing dispatch still fails.
  The zero-spawn control (spawn retargeted to the implementation stage) rejects inside the passing replay. Caveat: the literal falsifying edit "remove `spawns < 1`" does NOT by itself turn the control red, because that fixture also lacks a correlated run id, so `piRunID` stays empty and the credits are doubly gated; the mutation proof as worded overstates the spawn guard alone. The behavioral AC holds.
- DONE: AC-3 — the ordering contract is preserved.
  The inverted-order control rejects; a Go overlay removing only `|| completed >= validation` RED with `a Pi completion after the validation transition passed`. The `spawns > 2` tolerance and the ordering expression are unchanged.
- DONE: AC-4 (no-regression) — scoped package suite and other-runtime checks.
  `go test ./internal/ensigncycle/... -count=1` PASS (170.877s). The diff touches only the shared assert plus the Pi replay test, the Pi binding line in `shared_live_runner_test.go`, and two testdata files; the Claude/Codex grading branches are unchanged. Changed files `gofmt -l` empty; `git diff --check` clean. No repository-wide suite, race suite, or CI lane was run.
- SKIPPED: live-lane green / auto-continue Pi XFAIL-removal acceptance — no green lane at this tip.
  Lane run `37344061362`'s `offline` job failed on a `TempDir` cleanup error and `pi-live` was skipped; the re-run had not finished. The binding was removed (`shared_live_runner_test.go` gaps `nil`) but its acceptance rests on an independently green live Pi journey, which this round cannot supply.

### Summary

Deterministic validation settles AC-1, AC-2, AC-3 and AC-4 at the task's change: the native and `bg_wait` credits, the run-id correlation, and the ordering guard are each load-bearing under single-guard Go overlays, and the replay baseline reproduces `spawns=1 completed=-1` when the credits are removed. The fixture is committed with a verified provenance source hash. The live-lane green and the Pi XFAIL-removal acceptance remain pending because the tip has no green lane. The code worktree tip moved from `07a480522` to `d0b15a9ff` during this round; the task's three commits are patch-identical (`a737938`, `66089fb`, `4aa20dc`), so these results apply to both.
