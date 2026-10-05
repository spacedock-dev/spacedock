---
title: Repair the Pi recorded-gate-lifecycle journey
status: validation
source: "CI run 31770740214 (PR #685 pi-live, model openai/gpt-5.6-luna:max): TestLiveCommonRecordedGateLifecycle FAIL observed=[recorded-gate-lifecycle-violation], 'Blocked at the validation gate. The required committed reference is missing.'"
score: 0.85
sprint: pi-live-completeness
sprint-readiness: ready
group: pi-live-followup
id: gcmfwfjd9735b58sbzw7xsb8
started: 2026-10-04T04:35:04Z
worktree: .worktrees/spacedock-ensign-repair-pi-recorded-gate-lifecycle
gates:
    version: 1
    records:
        - id: gate:gcmfwfjd9735b58sbzw7xsb8:validation
          stage: validation
          attempts:
            - id: gate-attempt:gcmfwfjd9735b58sbzw7xsb8-validation-1
              briefing:
                id: briefing:gcmfwfjd9735b58sbzw7xsb8:validation:attempt-1:revision-1
                digest: sha256:945b09c94a62a4944c64f0c1a6d18f6bbd30d65e69c78a0b365189aa1a9f747c
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:gcmfwfjd9735b58sbzw7xsb8:validation:1
                briefing: briefing:gcmfwfjd9735b58sbzw7xsb8:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-10-05T17:22:55.89839Z"
                decision: approve
              application:
                target-stage: done
                state: pending
---

## Problem

The Pi FO does not complete the `recorded-gate-lifecycle` journey cleanly. The
strict oracle rejects it with `observed=[recorded-gate-lifecycle-violation]`:
"Blocked at the validation gate. The required committed reference is missing."
The FO reaches the validation gate but the retained reference the lifecycle
requires is absent, so the recorded-gate lifecycle is not durably complete on Pi.

This is an ordinary lane FAIL on Pi: `recorded-gate-lifecycle` is XFAIL-bound only
for `claude-opus` (`66dpwxgvsxt7cbxhmgvt3qp4`). There is no `liveXFail("pi",…)`
or `liveTODO("pi",…)` binding, so on Pi the journey is expected to PASS. The CI
`pi-live` lane is correctly red on an unowned, real Pi conduct gap.

## Visible value

A Pi operator runs a delegated-authority recorded gate and the lifecycle binds,
records, commits, and consumes exactly once before successor dispatch, with the
required committed reference present at the validation gate. Measured against
baseline: before this fix, the Pi `recorded-gate-lifecycle` journey FAILs with
`recorded-gate-lifecycle-violation` (missing committed reference); after, the same
run completes the recorded gate lifecycle to a PASS.

## Evidence

- CI: `Runtime Live E2E` run 31770740214, `pi-live` job, model
  `openai/gpt-5.6-luna:max`; `TestLiveCommonRecordedGateLifecycle` FAIL,
  `observed=[recorded-gate-lifecycle-violation]`, "Blocked at the validation
  gate. The required committed reference is missing." Recorded in the pnc
  pi-live log (17 common journeys ran; this was 1 of 3 failures, the other two
  being the known `default-headless-gate-stop` gap (nta) and the
  `keep-moving-posture` XPASS).

## Out of scope

- The `claude-opus` XFAIL binding (`66d`). This task owns only the Pi conduct.
- Shared XFAIL policy, the assert, or the fixture.
- Sonnet, Codex, or any other runtime's behavior on this journey.
- A new runtime, fixture, result format, or CI lane.

## Acceptance criteria

**AC-1 (VALUE) — The gate-prepare selected-source path follows the observed root.**

Verified by: the `fo-gate-lifecycle` Prepare step composes the selected-source path from the observed
root instead of the launch directory, so a gate room prepared from another working directory resolves
the right file; and the Pi `default-headless-gate-stop` binding is cleared, because that journey can
hold once the prepare step is corrected.

**Not delivered by this layer (recorded, not claimed):** the Pi delegated-authority journey
`recorded-gate-lifecycle`, where the first officer binds, records, commits and consumes delegated
authority exactly once before successor dispatch. This diff does not exercise that journey, and its
binding is already `nil`. So this layer claims the prepare-path fix and the headless binding clear,
and does not claim the delegated-authority journey. If that journey is still wanted, it needs its own
work.

**AC-2 — The Pi binding stays honest.**

Verified by: no `liveXFail("pi",…)` is added to mask the gap; the journey reaches
PASS by completing the lifecycle, not by weakening the assertion. A temporary
`liveTODO("pi",…)` if needed names this active task as owner and is removed on PASS.

**AC-3 — Other runtimes and the shared assert are preserved.**

Verified by: the `claude-opus` XFAIL binding and the shared assert are unchanged;
Sonnet and Codex behavior on this journey is unaffected.

**AC-4 — Offline and required-lane checks pass.**

Verified by: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
`go build -tags live ./internal/ensigncycle`, `go test ./...`, and
`go test ./... -race` pass; the Pi live lane passes the focused target.

## Test plan

Use focused offline gate and terminalization controls first. Use one exact Pi
`recorded-gate-lifecycle` target sequence only when Pi work is authorized.
Preserve all Sonnet and Codex behavior and the shared assert.

## Notes

- Coordinate with `pi-delegated-gate-continuation-reliability` (9w, group: gate),
  which owns the Pi recorded-gate journey under delegated conn — the missing
  committed reference may share a root cause in the gate-record/dispatch seam.
- Filed from the pnc pi-live run (31770740214), which surfaced this gap because
  pnc's parallelism change let all 17 common journeys run (vs -failfast stopping
  at the first failure on non-parallel branches).

## Stage Report: implementation

- DONE: Fix the officer instruction so the gate-prepare reference path is composed from the observed workflow root or the committed entity path, and is never reproduced from prose in the prompt.
  `skills/fo-gate-lifecycle/SKILL.md` Prepare step now reads "Compose the selected-source `<R>/<P>` (observed workflow root or committed entity path); never retype prose. Use absolute if cwd/state spellings differ." This restores the `<R>/<P>` composition that commit 6060c382a replaced with "Supply absolute, launch-cwd-relative, or state-relative judgment paths" (the wording that invited the hand-typed path). The 7700-byte component cap holds: file is 7695 B.
- DONE: Prove it: the focused checks, then a local live run of default-headless-gate-stop using SPACEDOCK_LIVE_RUNTIME=pi and a distinct SPACEDOCK_LIVE_ARTIFACT_DIR.
  After-fix run (base 05bdfa7da; branch later rebased onto the sibling's amended tip a2166585a — binding removed, skill fixed): `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/live-refpath SPACEDOCK_LIVE_RUNTIME=pi go test -tags live -count=1 -timeout 45m -run '^TestLiveCommonDefaultHeadlessGateStop$' ./internal/ensigncycle -v` -> `--- PASS: TestLiveCommonDefaultHeadlessGateStop (551.12s)`, no XFAIL and no XPASS alert (the binding is gone, so an XPASS is impossible), empty observed set. The run's `command.log` shows `gate prepare ... --artifact .../entity-snapshot.md --reference .../recorder-contract.md` exit=0, then `state commit` exit=0 (`state-head d70d84bf`).
- DONE: When that journey passes locally with an exact passing pair, clear its XFAIL binding and record the evidence in the entity.
  Removed `[]liveJourneyGap{liveXFail("pi","gcmfwfjd9735b58sbzw7xsb8")}` -> `nil` for `TestLiveCommonDefaultHeadlessGateStop` in `internal/ensigncycle/shared_live_runner_test.go` and dropped its now-obsolete comment; committed with the skill fix in bb3265618.
- DONE: Do not weaken the shared assertion, the binary's refusal, or the no-retry rule.
  Shared assertions (`assertGateHeld`, `assertRecordedGateHoldLog`, `assertImplementationWorkerLifecycle`) are byte-unchanged; `internal/gates/prepare.go` is untouched, so the fail-closed refusal of an unreadable/missing selected source stands; the skill's single-prepare/no-retry sentence is untouched.
- DONE: Focused offline checks.
  `gofmt -l` clean; `go vet -tags live ./internal/ensigncycle` clean; `go test -tags live -run '^$' ./internal/ensigncycle` compiles; `go test ./internal/contractlint/... ./internal/ensigncycle/... ./internal/gates/... ./internal/cli/... -run '...' -count=1` green (ensigncycle offline 308s, contractlint 1.3s, gates/cli gate-prepare green). Per the FO scope, the repository-wide `go test ./...`/`-race` was not run.
- SKIPPED: Repository-wide `go test ./...` and `go test ./... -race` (AC-4).
  Scoped out by the FO for this stage: offline proof was limited to the focused packages (contractlint, ensigncycle, gates, cli) and the focused live Pi lane supplied the journey proof, so the repo-wide suites were not run in this stage.

### Summary

This is hardening against a transcription slip, not the repair of a reproducible fault. The before-fix local run (peer base fb4428e9c, binding present, skill unfixed) did NOT reproduce the missing-reference fault at all: `gate prepare` succeeded (exit=0, `state=open`, correct `--artifact`/`--reference` paths, committed) and the sole red was the separate observer defect `implementation-worker-not-dispatched` (owned by mk72bnt1b5hsp9sfv83979xs), producing `XFAIL pi/default-headless-gate-stop observed=[implementation-worker-not-dispatched]`. The reference-path failure recorded from CI kept the random suffix while dropping one word, which is only possible when the path was written from the prompt prose rather than read from the boot record; the instruction change makes that harder without claiming a reproduced fault. Because the headless journey passes only once the observer fix is present, my branch is stacked on the sibling tip a2166585a (peer fb4428e9c + `2d5bd25f6` bg_wait/native-completion credit + `a2166585a`); the branch was first based on the sibling's then-tip 05bdfa7da and rebased onto a2166585a when the sibling amended its replay test, and the observer fix in `claude_runtime_helpers_test.go` is byte-identical across the two, so the live PASS holds. My commit on the final base is bb3265618. Residual risk: the journey is stochastic — the after-fix run needed a clean artifact directory (a killed prior run in the same dir made the harness see two root sessions and fail before the journey); a future rerun must start from an empty `SPACEDOCK_LIVE_ARTIFACT_DIR`.

## Stage Report: validation

- DONE: Re-run the focused checks.
  At bb3265618: changed-Go-file `gofmt -l` and `git diff --check` clean; `go vet -tags live ./internal/ensigncycle`, `go build -tags live ./internal/ensigncycle`, live-tag compile-only check, full contractlint package, focused gates/CLI prepare and ensigncycle gate/lifecycle/observer controls PASS. Commands/results: `/tmp/spacedock-validation-gcm-20261004/focused.log`; no candidate edits.
- DONE: Re-run the live Pi journey locally with the captain's credentials, a distinct and EMPTY SPACEDOCK_LIVE_ARTIFACT_DIR, only on a quiet machine.
  Started 2026-10-04T06:12:17Z at one-minute load 4.44 / 10 cores (5/15-minute loads 10.03/15.99, recorded rather than hidden). Newly created `/tmp/spacedock-gcm-validation-live-3zm6n9uf` had zero entries. Harness copied operator Pi auth into its isolated home; no API-key/OAuth override; model `openai-codex/gpt-5.6-luna:max`. Preflight: `/tmp/spacedock-validation-gcm-20261004/live-preflight.json`.
- DONE: Require a clean PASS with an empty observed set and no XFAIL and no XPASS alert.
  `SPACEDOCK_LIVE_RUNTIME=pi SPACEDOCK_PI_LIVE_REQUIRED=1 SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-gcm-validation-live-3zm6n9uf go test -tags live -count=1 -timeout 45m -run '^TestLiveCommonDefaultHeadlessGateStop$' ./internal/ensigncycle -v` PASS (456.79s), exit 0; no XFAIL/XPASS. Empty observed set follows from the unbound `gradeLive(false, ...)` PASS branch (the runner does not print observed on PASS), not from absence of error text alone.
  Artifact `pi-common/default-headless-gate-stop/command.log`: exactly one prepare begin/exit pair, exit=0, followed by successful state commit, head `9bbae0b79b1b23b905c8054fee3064c4492e04e4`; one root session. Selected artifact was committed entity `index.md`, reference workflow `README.md`. Candidate skill read appears at 06:14:47Z. Audit: `/tmp/spacedock-validation-gcm-20261004/live-audit.json`; raw run log alongside. This pair justifies removing the headless target's binding, not closing AC-1.
- DONE: Judge the skill wording on its own terms.
  It supplies a sensible exact-path recipe but remains a skippable officer instruction: old text already said keep paths/derive R and P; no executable construction/enforcement changed. It can avoid transcription if followed, cannot guarantee prevention, and this run does not measure a reduction in slip probability.
- DONE: Confirm the component cap held and report the measured byte value.
  `wc -c skills/fo-gate-lifecycle/SKILL.md`: 7695 bytes, cap 7700, headroom 5; base 7697. `TestFOInstructionComponentCaps` compares actual bytes against an independent fixed limit and would fail at 7701, not merely on changed phrasing.
- DONE: Judge whether "hardening against a transcription slip" is the right description.
  Yes: instruction-only hardening, not demonstrated fault repair. The recorded before-run already prepared successfully and failed only the sibling observer; no missing-reference reproduction or isolated before/after causal proof is supplied. The prior run's causal account remains reported evidence, not an independently reproduced fault. This rerun used README.md, not the prior recorder-contract.md reference.
- FAILED: AC-1 (VALUE) — The exact Pi recorded-gate-lifecycle target passes; judge the scope mismatch.
  No authorized live `TestLiveCommonRecordedGateLifecycle` evidence was produced here. That target exercises delegated record/consume/successor dispatch; the headless target intentionally forbids these. Shared missing-reference symptoms make the instruction relevant, but cannot broaden the entity's explicit journey-specific AC. F1 below: material evidence defect; captain scope decision required, not another automatic repair cycle.
- DONE: AC-2 — The Pi binding stays honest.
  Diff removes only the headless Pi XFAIL and adds no masking binding. RecordedGateLifecycle remains unbound. Unbound semantic failure controls grade FAIL; a semantic failure in the rerun would fail rather than XFAIL.
- DONE: AC-3 — Other runtimes and the shared assert are preserved (diff-level check).
  Only Prepare prose and the headless binding/comment changed; shared assertions and binary refusal are byte-unchanged versus a2166585a. RecordedGateLifecycle binding is already nil at both base and candidate, so the entity's claimed claude-opus XFAIL is stale. No cross-runtime live compatibility claim is made for the shared wording.
- SKIPPED: AC-4 repository-wide tests/race and CI lane.
  Explicit FO prohibition; focused vet/build/format/test checks and local Pi proof passed. Full AC-4 is not reproduced and no CI run was started.
- DONE: Audit for tautological tests in the captain's sense.
  No new regression test was added by bb3265618. Existing `TestGradeLiveRunsEveryAssertion` counts argument evaluation before gradeLive executes; `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` has a cleanup Replace whose needle is absent (no-op). Neither supplies behavioral proof. Other controls reject mutated gate identity/digest/state, missing commit, extra successful prepare, consume/dispatch, missing worker and inverted completion; these can fail on real behavior changes. F2 below records the excluded checks and the retry coverage gap.
- DONE: Semantic adversarial pass over the changed behavior.
  Existing exercised matrix: absolute/state-relative/launch-relative sources share immutable binding; missing/ambiguous/permission-denied paths refuse; symlink/directory/foreign-root inputs preserve bytes; open/resolved/applied/wrong-digest gates differ. Live independently verifies one prepare without retry. No changed hot path, allocation, parser or size algorithm warrants a new scaling test; no deliverable repairs were made.

### Summary

REJECTED as completion of this entity: the independent headless PASS/empty-observed pair supports that target's binding removal and the 7695-byte wording is acceptable as modest hardening, but AC-1 still has the wrong journey as evidence. Hold task acceptance for a captain-authorized scope reconciliation; do not automatically repair or silently redefine the value AC. No outcome defect was reproduced by this validation.

## Review-finding disposition

### F1 — Wrong-journey evidence for the value AC

- Observation: `shared_live_runner_test.go:131` tests stopping open, while line 145's RecordedGateLifecycle and entity AC-1 require delegated lifecycle completion; neither a shared symptom nor a headless PASS establishes that promised outcome.
- Released user/normal workflow: Pi operator delegates authority through the recorded-gate lifecycle. Observable harm: accepting this report would declare that workflow repaired without evidence it can record, consume and dispatch successfully.
- Authority: value-ac[AC-1] requires the exact Pi RecordedGateLifecycle target to PASS. Trigger evidence: entity AC-1/test plan versus bb3265618's two-file diff and this exact headless-only live command.
- Advisory classification: evidence defect; release scope Material; ownership task acceptance/scope (captain), not an observed binary regression. Recommend Needs decision/hold; related missing-reference hardening may ship only under reconciled scope, not as proof this entity is complete. FO authorization remains outstanding.

### F2 — Existing vacuous checks and retry-proof blind spot

- Observation: `live_grade_unit_test.go:31` counts pre-call argument evaluation; `claude_runtime_helpers_test.go:669` replaces absent cleanup text. The latter test's positive fixture also includes a failed prepare followed by success, and `assertRecordedGateHoldLog` counts successful prepares only, so PASS alone is not proof of no retry.
- Released user/normal workflow: maintainers interpreting these checks as lifecycle/no-retry proof. Observable harm: overclaiming coverage; no user outcome regression is established here because independent state/negative controls and the actual single-attempt command log were checked.
- Authority: none: no value AC failure is attributed to these redundant checks in this candidate; they are explicitly excluded from the credited proof. Trigger evidence: argument evaluation semantics, absent replacement needle, and the failed-then-success fixture at `claude_runtime_helpers_test.go:658-659`.
- Advisory classification: evidence defect; release scope Deferred risk, not a blocker for this narrow diff. Ownership existing harness tests, outside this verification-only assignment. Trigger outside the accepted proof is relying on these checks alone; promote to Material if a future acceptance relies on them or a supported live retry is mislabeled compliant. Recommend separate disposition, no candidate mutation; FO authorization remains outstanding.

## Stage Report: implementation (cycle 2)

- DONE: Rebase this branch onto the pi-live lane layer at 5ead9b85c (spacedock-ensign/pi-live-lane-pin-refresh). Keep this task's own change intact. Do not take the base side wholesale, and do not drop this task's commits.
  `git rebase 5ead9b85c` on `spacedock-ensign/repair-pi-recorded-gate-lifecycle`; 5ead9b85c is now an ancestor of the new tip `cf2b42cb7`. Own change intact as three commits: `3479b614b` (credit bg_wait + native completion), `2e61d3d31` (clear auto-continue Pi XFAIL), `cf2b42cb7` (gate-prepare observed-root anchor + clear default-headless Pi XFAIL). `git diff 5ead9b85c..HEAD --stat` = 4 files, +101/-29; no own hunk dropped.
- DONE: Resolve the conflicts in internal/ensigncycle/shared_live_runner_test.go, and in skills/first-officer/references/fo-install.md where this branch touches it. Report each conflict, the two sides, and the resolution you chose.
  Only `internal/ensigncycle/claude_runtime_helpers_test.go` conflicted (NOT either predicted file). Base side: the single-tool `subagent_wait` condition + old comment (commit 974ed2bed). Incoming side `2d5bd25f6`: the same condition widened to `(subagent_wait || bg_wait)` + explanatory comment, with the native `subagent-notify` credit cleanly added below. Resolution: kept the incoming (widened) side — it is a strict superset of base's behavior and retains `subagent_wait`. `shared_live_runner_test.go`: applied cleanly, no conflict (a2166585a/bb3265618 replayed without merge). `skills/first-officer/references/fo-install.md`: no conflict and no delta — this branch's `17858e911` was dropped as already-applied (base carries its equivalent `c6170ea8c`); file is byte-identical to base.
- DONE: Run this task's focused tests after the rebase. Do not run the repository-wide suite and do not start a CI run.
  `gofmt -l` clean on the three changed Go files; `go vet -tags live ./internal/ensigncycle` OK; `go test -tags live -run '^$' ./internal/ensigncycle` compiles; `go test ./internal/ensigncycle -run '^(TestPiWorkerLifecycleObservationReplay|TestImplementationLifecycleAndObserverNegativeControls|TestCodexNativeLifecycleUsesCorrelatedSessionHandle|TestCodexNativeLifecycleParentRolloutLookupFailsClosed|TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle|TestAssertGateHeld|TestAssertGateHeldAcceptsPreparedFixtureBinding|TestMissingPreparedGateExpectationIsSemanticallyGraded)$'` PASS (1.99s); `go test ./internal/contractlint -run '^TestFOInstructionComponentCaps$'` PASS (SKILL.md 7695 B, cap 7700). Falsification: removing the bg_wait/native credits turns `TestPiWorkerLifecycleObservationReplay` RED (`completed=-1`), restored -> PASS. No repo-wide suite, no CI run.
- DONE: Report the new tip commit, the files that conflicted, and the evidence for each resolution.
  New tip `cf2b42cb783f471922f87bd525bf642117d6e5fe`. Conflicted file: `internal/ensigncycle/claude_runtime_helpers_test.go` (resolution above). `shared_live_runner_test.go` and `skills/first-officer/references/fo-install.md` did not conflict; that is the evidence, not a claim of resolution.

### Summary

Rebased cleanly onto 5ead9b85c with one real conflict, resolved in favor of this task's own widened completion-credit logic; the branch's change survives intact and its focused offline tests pass. On the FO's second point: the `fo-gate-lifecycle` Prepare wording edit is NOT the fix for the `implementation-worker-not-dispatched` red. That red is the shared worker-lifecycle assert computing `completed=-1` because Pi completes via `bg_wait`/native `subagent-notify`; commit `3479b614b` credits those surfaces and is the fix (falsification above shows the red returns without it). The SKILL.md edit addresses the separate missing/mistyped reference path (`gate-not-held`/`gate-hold-violation`), not dispatch observation. The AC-1 `recorded-gate-lifecycle` vs `default-headless-gate-stop` mismatch remains an open captain scope decision and was not resolved here.

## Stage Report: implementation (cycle 3)

- DONE: Rebase this branch onto the observer fix at 9755d96d1, which is repair-pi-worker-lifecycle-observation rebased onto the pi-live lane layer.
  `git rebase 9755d96d1` from prior tip `cf2b42cb7` (merge-base `5ead9b85c`). Rebase skipped `3479b614b` and `2e61d3d31` as previously applied; new tip `e5907d00bfad2312a7cd531d80a4d9f72b920652`. `git merge-base HEAD 9755d96d1` = `9755d96d1`, so the observer fix is a direct ancestor.
- DONE: What must remain is only this branch's own commit, and say so if more remains.
  `git log --oneline 9755d96d1..HEAD` yields exactly one commit, `e5907d00b "Anchor the gate-prepare selected-source path to the observed root"`, carrying both the gate-prepare anchor and the default-headless Pi XFAIL clear. Nothing else remains. `git diff 9755d96d1..HEAD --stat` = 2 files, +2/-13 (`internal/ensigncycle/shared_live_runner_test.go` binding/comment cleared; `skills/fo-gate-lifecycle/SKILL.md` Prepare step).
- DONE: The commits that dropped as already-applied were patch-identical to the observer fix.
  `3479b614b` patch-id `a7379382bdf3f6f9922dc8c9f575adfb79a78a4c` == observer `67d413a39`; `2e61d3d31` patch-id `66089fb358f20ac885c8d808b485dac6945c5eff` == observer `c65bba735`. Both `git show <c> | git patch-id --stable` pairs identical, so the drops are the intended already-applied result.
- DONE: Then push the rewritten branch with a lease-guarded force.
  `git push --force-with-lease origin HEAD:spacedock-ensign/repair-pi-recorded-gate-lifecycle` -> `+ cf2b42cb7...e5907d00b ... (forced update)`, exit 0. Lease was safe: local tracking and `git ls-remote origin` both showed `cf2b42cb7` (remote unchanged since last push), so the lease matched. `git ls-remote origin spacedock-ensign/repair-pi-recorded-gate-lifecycle` now reports `e5907d00bfad2312a7cd531d80a4d9f72b920652`.
- DONE: Report the new tip, the commits that dropped as already-applied, and the commands you used to confirm it.
  New tip `e5907d00bfad2312a7cd531d80a4d9f72b920652`; dropped `3479b614b`, `2e61d3d31` (patch-id evidence above); commands: `git rebase 9755d96d1`, `git log --oneline 9755d96d1..HEAD`, `git diff 9755d96d1..HEAD --stat`, `git show <c> | git patch-id --stable`, `git ls-remote origin ...`, `git push --force-with-lease ...`.
- DONE: Focused tests after the rebase; no repository-wide suite and no CI run started.
  `gofmt -l internal/ensigncycle/shared_live_runner_test.go` clean; `git diff --check 9755d96d1..HEAD` clean; `go vet -tags live ./internal/ensigncycle` OK; `go test -tags live -run '^$' ./internal/ensigncycle` compiles; `go test ./internal/contractlint -run '^TestFOInstructionComponentCaps$' -count=1` PASS (SKILL.md 7695 B, cap 7700); `go test ./internal/ensigncycle -count=1 -run '^(TestPiWorkerLifecycleObservationReplay|TestImplementationLifecycleAndObserverNegativeControls|TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle|TestAssertGateHeld|TestAssertGateHeldAcceptsPreparedFixtureBinding|TestMissingPreparedGateExpectationIsSemanticallyGraded)$'` PASS (1.68s). No repo-wide `go test ./...`/`-race`, no CI run.

### Summary

Rebased `spacedock-ensign/repair-pi-recorded-gate-lifecycle` onto the observer fix `9755d96d1`; the two replayed worker-lifecycle commits dropped as patch-identical, leaving only this task's own commit `e5907d00b` (gate-prepare selected-source anchor + default-headless Pi XFAIL clear). Pushed with an authorized lease-guarded force; remote branch tip is now `e5907d00b`. Focused offline checks pass; the live Pi journey was not rerun here. Residual risk unchanged from prior cycles: the value AC-1 journey-specific live evidence and the captain scope reconciliation remain open.
