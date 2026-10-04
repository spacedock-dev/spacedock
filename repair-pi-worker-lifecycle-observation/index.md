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
