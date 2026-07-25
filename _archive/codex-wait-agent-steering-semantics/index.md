---
title: Codex wait_agent steering semantics describe captain input as active-loop resumption
status: done
source: "Captain request 2026-07-23: replace misleading wait-interruption language and use the corrected behavior in-session"
started: 2026-07-23T14:43:01Z
completed: 2026-07-25T21:01:05Z
verdict: passed
score: 0.9
worktree: .worktrees/spacedock-ensign-codex-wait-agent-steering-semantics
issue:
id: 6gkz4z2qweheyj17ck5tythn
gates:
    version: 1
    current:
        gate: gate:docs-dev:6g:validation
    records:
        - id: gate:docs-dev:6g:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6g-ideation-1
              briefing:
                id: briefing:docs-dev:6g:ideation:attempt-1:revision-1
                digest: sha256:a91e6243eab3b12d756db99283f05e7d74aa48f037bc418dc0327060536fc768
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6g:ideation:1
                briefing: briefing:docs-dev:6g:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T14:53:56.796844Z"
                decision: approve
                reason: The ideation supplies a real steering trace, preserves durable completion authority, and bounds implementation to five Codex-only files with independent negative controls.
                adoption-note: file codex runtime issue; dispatch, and do not forget to asyncwait; use captain steering as active-loop resumption while workers continue unchanged.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6g-ideation-2
              briefing:
                id: briefing:docs-dev:6g:ideation:attempt-2:revision-1
                digest: sha256:85cca2e27e401b206fa7d0d375f3a331256c520bfd1af66f3c7416851b3eeaec
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6g:ideation:2
                briefing: briefing:docs-dev:6g:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-25T17:48:22.634251Z"
                decision: revise
                reason: The steering semantics and adversarial proof are sound, but the gate package still makes file and LOC estimates binding reset authority; revise them to advisory drift evidence and retain only semantic scope expansion as a reset trigger.
              application:
                action: feedback
                target-stage: ideation
                state: superseded
            - id: gate-attempt:6g-ideation-3
              briefing:
                id: briefing:docs-dev:6g:ideation:attempt-3:revision-1
                digest: sha256:d14f37ba0ac249ae495ebc96e3c90f3144cae9780546d0d645ae8c142d74b227
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6g:ideation:3
                briefing: briefing:docs-dev:6g:ideation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-07-25T17:55:05.089935Z"
                decision: approve
                reason: The corrected design preserves the proven Codex-only steering semantics and durable completion authority, supplies falsifiable lifecycle and messaging controls, and makes file/LOC projections advisory while semantic scope remains binding.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:docs-dev:6g:validation
          stage: validation
          attempts:
            - id: gate-attempt:6g-validation-1
              briefing:
                id: briefing:docs-dev:6g:validation:attempt-1:revision-1
                digest: sha256:12ef44e574a581df377d90ae4b3e321558a8ad8fb44c480af40705b9dae506ca
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6g:validation:1
                briefing: briefing:docs-dev:6g:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T20:08:14.744192Z"
                decision: approve
                reason: The real Codex steering drive demonstrates the intended active-loop resumption, all 36 adversarial mutants and repository gates pass, the stale-empty false pass is closed, and the proposed evidence-line schema was correctly declined as an acceptance narrowing.
              application:
                action: advance
                target-stage: done
                state: consumed
                blockers: []
mod-block:
pr: pr-merge:569
archived: 2026-07-25T21:01:05Z
---

## Problem

The Codex First Officer runtime currently imports the harness label “Wait interrupted by new input” into its operating language. That label does not describe the behavior: `wait_agent` is asynchronous monitoring, captain input resumes the FO's active loop, and unresolved workers continue unchanged. Repeating an interruption disclaimer before each wait epoch makes normal steering sound destructive and adds noise.

The defect appears in three Codex-only surfaces. The runtime binding says steered input “interrupts” `wait_agent`; the wait notes require a captain-facing interruption disclaimer at each epoch; and the idle-notification recipe calls the mechanism a foreground wait/non-terminal wait return. Those phrases pull a harness result label into the FO's lifecycle model even though the live trace shows no worker cancellation, closure, or redispatch.

## Current surface map

| Surface | Current clause or fixture | Disposition |
| --- | --- | --- |
| `skills/first-officer/references/codex-first-officer-runtime.md:21` | `«async-dispatch»` says steered captain input “interrupts it immediately.” | Replace with active-loop resumption and unchanged workers. |
| Same file, `## Codex wait notes` at lines 31-37 | Lead paragraph repeats “interrupts,” mandates the epoch-start disclaimer, and the `### Idle-monitoring interruption` subsection repeats it. | Replace the lifecycle wording, remove the mandatory cue, and quarantine the harness label as non-semantic. |
| Same file, lines 23 and 33 | A wait return is an idle observation; no wait result alone completes a worker. | Preserve verbatim in meaning. |
| `docs/dev/codex-idle-notification-probe.md:3-32,59-68` | Calls explicit `wait_agent` use “foreground wait” and operator input a non-terminal foreground-wait return. | Rename the comparison/classification to async idle monitoring and use the same active-loop wording. Preserve queued-flush/autonomous-wake distinctions. |
| `skills/integration/codex_idle_notification_test.go:13-18` | Allowed evidence enum contains `foreground_wait`. | Rename only that unused enum member to `async_idle_monitoring`; the checked-in evidence is `queued_flush`, so no evidence migration is needed. |
| `internal/ensigncycle/codex_dispatch_evidence_test.go:21-89` and regression fixtures | Wait is credited only between dispatch-build and a later durable report read. | Reuse its durable-evidence rule; do not weaken it or make the harness label completion evidence. |
| `internal/ensigncycle/shared_reviewer_reuse_test.go:244-367` | Correlates Codex worker identity through task/thread handles and rejects uncorrelated reuse. | Reuse this handle-correlation pattern in the new steering fixture; no production change. |

No shared-core contradiction was found. The shared core decides when work is ready or the FO is idle; only the Codex adapter names `wait_agent` and the harness label. Claude and Pi remain out of scope.

## Required behavior

The Codex runtime contract must express these semantics directly:

- Captain input resumes the FO's active loop while workers continue unchanged.
- When the FO becomes idle again, it resumes monitoring unresolved workers.
- The contract does not require a repetitive interruption disclaimer before every wait epoch.
- A wait return or captain message never becomes worker-completion evidence; durable reports and the final-status signal remain authoritative.

Scope is the Codex runtime contract and its behavioral proof. Do not change `wait_agent`, invent cancellation/restart state, or broaden this into a generic scheduler redesign.

## Proposed wording delta

Apply this exact semantic edit to the Codex adapter:

```diff
-`wait_agent` is asynchronous with respect to worker progress and captain interaction: it is idle monitoring only, does not stop the FO event loop, and steered captain input interrupts it immediately so the FO can resume active work.
+`wait_agent` is asynchronous monitoring: captain input resumes the FO's active loop while workers continue unchanged.
```

In `## Codex wait notes`, replace the interruption/cue cluster with:

> Captain input resumes the FO's active loop while workers continue unchanged. When the FO becomes idle again, it MUST resume monitoring unresolved workers. Do not preface each monitoring epoch with an interruption disclaimer.

Replace `### Idle-monitoring interruption` with:

> ### Harness wait-return label
>
> “Wait interrupted by new input” is a harness implementation label only. Do not repeat it as FO language or derive worker cancellation, closure, redispatch, or completion from it.

Keep the timeout/retry paragraph and the following rule intact in meaning: a wait result must be attributed by mailbox/task/durable state and no wait result alone completes a worker.

For `docs/dev/codex-idle-notification-probe.md`, change `Foreground wait comparison` / `foreground_wait` to `Async idle-monitoring comparison` / `async_idle_monitoring`, and replace the step-4 language with:

> If captain input resumes the FO's active loop, record the worker as unchanged and continue useful active-scope work. When the FO becomes idle again, resume monitoring the same unresolved worker. The harness return label is not worker completion, failure, closure, redispatch, or idle-wake evidence.

The simplest alternative was to delete only the word “interrupts.” It is insufficient because the mandatory disclaimer subsection and probe vocabulary would continue teaching the old lifecycle model. Moving the rule to shared core was also rejected: it would broaden a Codex harness correction into Claude/Pi semantics without evidence.

## Spike: real steering trace

No new runtime primitive is assumed. A real Codex 0.145.0 session trace from 2026-07-23 exercised the risky path:

1. `wait_agent(timeout_ms:300000)` began at `14:35:21.736Z`; the harness returned “Wait interrupted by new input” at `14:39:46.715Z`, immediately followed by captain input.
2. At `14:39:54.317Z`, `list_agents` still reported `/root/skill_wiring_commander/spacedock_ensign_r4xva464wf_implementation_correction` as `running`.
3. The FO performed active-scope work, including 20 roster/command/spawn calls in the reduced trace; there were zero `interrupt_agent` calls and zero redispatches of that task path.
4. Only after that active work did the FO call `wait_agent` again at `14:47:46.147Z`.

A `jq` reduction over the raw session produced: same worker running `true`, active-work calls `20`, later waits `1`, target cancellations `0`, target redispatches `0`. This is the cheapest observed discriminator for active-loop resumption. It does not by itself prove eventual durable completion, so the implementation fixture must append correlated final-status plus durable-report events and must fail when either is removed or stale.

## Behavioral proof approach

Add one Codex-only ordered-event fixture and assertion. It serves AC-1 through AC-3 by correlating one task path/completion epoch across: spawn → idle monitor → captain input → same worker still running → useful FO work → active scope empty → monitor again → matching final status → durable report read. The assertion must reject four planted variants: a target cancellation/second spawn, monitoring reinstalled while active work remains, a wrong task path or stale completion epoch credited as done, and an old captain-facing disclaimer before each of two epochs. Harness tool output may contain its implementation label; captain-facing FO messages may not repeat it.

The simplest alternative, a contractlint phrase check over the adapter, is insufficient: this repository retired those checks because they prove only that prose contains a phrase. A shared live-scenario entry is also unnecessary and would force Claude/Pi parity for a Codex-only harness concern. The ordered fixture is the low-cost gate; validation repeats the live drive and records the durable report.

## Acceptance criteria

**AC-1 (VALUE)** In a Codex drive where one worker remains unresolved across captain input, the FO handles the input, completes all captain-authorized active work, and later monitors the same task path/completion epoch without cancellation, closure, or redispatch. The measured baseline is one correlated worker before and after steering, at least one useful active-work event between the two monitoring calls, zero target lifecycle mutations, and eventual matching final-status plus durable report. Verified by the ordered fixture and a validation-time live replay.

**AC-2** Captain input has no completion effect: the worker remains unresolved until a matching final-status signal and durable report are observed, and monitoring resumes only after active work is exhausted. Verified by planted cancellation, redispatch, premature-monitoring, wrong-handle, stale-epoch, and wait-return-only fixtures; each must fail independently.

**AC-3** Two or more monitoring epochs produce zero repeated captain-facing uses of the harness label or old mandatory interruption disclaimer. Verified by separating harness tool output from FO `agent_message` output; a planted disclaimer before both epochs fails while ordinary progress/idle communication passes.

**AC-4** The shipped change is limited to the Codex adapter, Codex probe/evidence, and Codex behavioral tests; Claude, Pi, shared core, and runtime tool/API behavior remain unchanged. Verified by scoped diff review, `go test` focused/full/race gates, and live-tag compile or replay appropriate to the touched paths.

## Intended surface and reconciliation

The intended implementation surface remains the following advisory planning estimate:

- `skills/first-officer/references/codex-first-officer-runtime.md`: 6-10 changed lines, net negative or flat.
- `docs/dev/codex-idle-notification-probe.md`: 12-20 changed lines.
- `skills/integration/codex_idle_notification_test.go`: 1-4 changed lines for the classification rename.
- `docs/dev/_evidence/codex-wait-agent-steering-semantics/2026-07-23-dogfood.json`: 35-55 new lines containing only the reduced event slice, not the full private transcript.
- `internal/ensigncycle/codex_wait_agent_steering_test.go`: 135-195 new lines for ordered replay, handle/epoch correlation, output-channel separation, and the seven planted negative cases.

These five likely files and their roughly 190-285 changed LOC are drift evidence, not a cap, tolerance, or reset trigger. Implementation reports the actual surface against this estimate and explains additions, removals, helper extraction, or larger/smaller LOC; no arithmetic delta alone blocks the stage or forces a design reset.

The intended semantic surface is authoritative: the Codex adapter, Codex probe/evidence, and Codex behavioral proof needed for this correction. A design reset is required only if implementation changes shared-core, Claude, or Pi behavior; changes the `wait_agent` or runtime API; adds cancellation/restart state; broadens into scheduler redesign; weakens final-status plus durable-report completion authority; or creates another obligation outside this Codex-only correction. A different file count, fixture placement, helper split, or test harness within that semantic boundary requires reconciliation, not re-gating by arithmetic.

## Test plan

- First add the reduced evidence fixture and ordered replay assertion. Run the passing trace, then each planted variant independently; cost low, fixture-only, no model spend.
- Apply the exact adapter/probe wording. The behavioral fixture—not a markdown phrase grep—remains the semantic proof. Run `go test ./internal/ensigncycle ./skills/integration -run 'TestCodex.*(Steering|IdleNotification)' -count=1`.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- Compile live-tagged Codex tests. During validation, repeat the live steering drive with one delayed no-write worker, one captain input while monitoring, one useful active action, the same task path still unresolved, a later monitor call, matching final status, and a durable report read. Store only the reduced evidence record.

## Stage test gates

- Ideation identifies every current clause and behavioral fixture affected, proposes an exact wording delta, and spikes the cheapest scenario that distinguishes steering from worker interruption.
- Implementation reconciles its actual file/LOC surface against the advisory estimate, stays within the semantic boundary above, runs focused and full/race gates, and requests Roborev because this is shipped runtime scaffolding.
- Validation performs the required detached adversarial audit and a Codex behavioral drive capable of catching cancellation, redispatch, stale completion, and repetitive-disclaimer regressions.

## Stage Report: ideation

- DONE: Map the current Codex wait clauses and affected behavioral fixtures, then propose the exact minimal wording delta.
  The surface map covers the adapter binding/wait notes, probe recipe and enum, durable dispatch evidence, and handle-correlation fixture; the proposed wording removes interruption semantics while preserving no-wait-result completion.
- DONE: Demonstrate the cheapest falsifiable scenario where captain steering preserves one unresolved worker and monitoring resumes only when the FO is idle again.
  A reduced real Codex 0.145.0 trace observes the same task path still running, 20 active-work calls, zero cancellation/redispatch, and a later wait after active work.
- DONE: Produce acceptance evidence and an expected file/LOC surface that distinguish active-loop resumption from cancellation, redispatch, stale completion, and repetitive disclaimer behavior.
  AC-1 through AC-4 define correlated task/epoch and output-channel checks; the five-file 190-285 LOC surface includes seven independent planted negative cases.

### Summary

Ideation narrows the correction to Codex runtime wording, the Codex idle-notification vocabulary, and a behavioral replay based on the real session that prompted the task. The design treats captain input as active-loop resumption, preserves durable completion authority, and explicitly rejects cancellation, redispatch, premature monitoring, stale completion, and repetitive disclaimer variants.

## Stage Report: ideation (cycle 2)

- DONE: Replace the exact-file/LOC tolerance and automatic re-gate clauses with advisory surface reconciliation; only semantic expansion beyond the Codex-only correction may trigger a design reset.
  `Intended surface and reconciliation` keeps the five-file/190-285 LOC estimate as drift evidence, removes the ±60/cap authority, and enumerates the semantic reset triggers.
- DONE: Preserve the proven active-loop-resumption behavior, durable completion authority, planted negative controls, and exclusions for shared core, Claude, Pi, runtime API behavior, and scheduler redesign.
  Problem, wording, spike, AC-1 through AC-4, and the seven negative cases remain unchanged; the authoritative boundary restates every exclusion from the captain's correction.
- DONE: Append a corrected ideation report that states the chosen boundary clearly and gives the First Officer enough evidence to prepare a clean successor Briefing.
  This cycle records advisory numeric reconciliation versus binding semantic expansion and supersedes the prior gate package's arithmetic reset condition.

### Summary

The correction changes authority, not behavior: file and LOC estimates now support reconciliation without acting as caps or automatic reset triggers. A design reset is reserved for semantic expansion beyond the accepted Codex-only steering correction, including cross-runtime/API/scheduler changes or weakened durable completion authority.

### Feedback Cycles

- Cycle 1: CHANGES REQUESTED — Roborev job 2362; surface 5 files/519 LOC vs estimate 5 files/190–285 LOC (182% of upper estimate); AC unchanged
- Cycle 2: CHANGES REQUESTED — Roborev job 2378; surface 6 files/704 LOC vs estimate 5 files/190–285 LOC (247% of upper estimate); AC unchanged
- Cycle 3: CHANGES REQUESTED — Roborev job 2386; surface 6 files/704 LOC vs estimate 5 files/190–285 LOC (247% of upper estimate); AC unchanged
- Cycle 4: CHANGES REQUESTED — Roborev job 2392; surface 6 files/713 LOC vs estimate 5 files/190–285 LOC (250% of upper estimate); AC unchanged
- Cycle 5: CHANGES REQUESTED — Roborev job 2398; surface 6 files/768 LOC vs estimate 5 files/190–285 LOC (269% of upper estimate); AC unchanged
- Cycle 6: CHANGES REQUESTED — Roborev job 2404; surface 6 files/768 LOC vs estimate 5 files/190–285 LOC (269% of upper estimate); AC unchanged
- Cycle 7: CHANGES REQUESTED — Roborev job 2405; surface 6 files/865 LOC vs estimate 5 files/190–285 LOC (304% of upper estimate); AC unchanged
- Cycle 8: CHANGES REQUESTED — Roborev job 2409; surface 6 files/886 LOC vs estimate 5 files/190–285 LOC (311% of upper estimate); AC unchanged
- Cycle 9: REJECTED — independent validation; surface 6 files/886 LOC vs estimate 5 files/190–285 LOC (311%); AC unchanged; material stale-empty ordering false pass routed to implementation, while a new mandatory evidence-line schema was declined as an AC-narrowing report-format obligation.

## Stage Report: implementation

- DONE: Implement the Codex-only active-loop-resumption wording and probe/evidence vocabulary while preserving final-status plus durable-report completion authority and all cross-runtime/API exclusions.
  Commits `12c9a123` through `2e6e88c5` update only the Codex adapter, Codex probe/evidence, and Codex-focused test vocabulary; shared core, Claude, Pi, runtime APIs, scheduler behavior, and cancellation state are unchanged.
- DONE: Add the ordered correlated-worker behavioral proof and independent mutants for cancellation/redispatch, premature monitoring, wrong identity or stale epoch, wait-return completion, and repeated captain-facing disclaimer language.
  Commits `c1950744`, `134460f5`, `dbc58d04`, `88ad9121`, `6ff6361d`, and `9c903da6` make each listed mutant fail independently and also reject unattributable spawns, duplicate completion authorities, quoted/fenced report claims, and empty durable stage results.
- DONE: Declare and reconcile the actual semantic surface, run focused/live-tag/full/race verification, request final Roborev review, triage findings by materiality, and report exact acceptance evidence.
  The final surface is 6 files/886 changed LOC versus the advisory 5 files/190–285 LOC estimate: the sixth file is the retained durable report, and the additional LOC are explicit ordered-event validation and independent mutants added through review without semantic expansion.
  `go test ./internal/ensigncycle ./skills/integration -run 'TestCodex.*(Steering|IdleNotification)' -count=1`, the contractlint containment test, and the live-tag compile passed.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed after the final change.
  Roborev jobs 2362–2405 produced material correlation, ordering, uniqueness, containment, and durable-report fixes; job 2409 repeated the raw/private live-replay request, declined because the approved plan assigns that evidence to validation.

### Summary

Implementation now describes captain input as resuming the First Officer's active loop while unresolved workers continue unchanged, then reinstalls monitoring only after active work is exhausted. Completion remains gated by one matching final-status signal plus a parsed durable implementation report; the Codex-only semantic boundary is preserved.

## Stage Report: validation

- DONE: Run the required real Codex steering drive: one unresolved assignment across captain input, useful active work, resumed idle monitoring, matching final status, and a parsed durable implementation report; retain only a reduced evidence record.
  Live task `/root/spacedock_ensign_6gkz4z2qwe_validation/live_steering_sentinel`, epoch 1, remained running after routed captain input; active work verified candidate `9c903da6402c82fbc595143c64392733d761d243` and a clean worktree, a later 10-second wait monitored the same unresolved task, and the matching final status completed once. `spacedock status --read ... --json` parsed the durable implementation report; only these reduced facts were retained.
- DONE: Reproduce every planted fixture mutant and run focused, containment, live-tag, full, and race checks; verify repeated captain-facing harness/disclaimer language fails and Claude/Pi/shared-runtime behavior is untouched.
  All 35 independent planted mutant subtests rejected their target changes; focused steering/idle, contractlint containment, live-tag compile, `go test ./...`, and `go test ./... -race` passed. A repeated captain-facing disclaimer is one rejecting mutant, and the six-file diff contains only Codex adapter/probe/evidence/test paths.
- DONE: Perform a semantic adversarial audit of assignment/path/epoch attribution, event order, uniqueness, and durable completion authority; classify findings and recommend PASSED or REJECTED.
  A detached checkout at `9c903da6` reproduced assignment/path/epoch, order, uniqueness, and completion-authority controls, then exposed two material evidence false passes. Recommendation: REJECTED; both are narrow AC-1/AC-2 evidence defects, with no observed outcome defect or design reset.
- FAILED: AC-1/AC-2 material evidence defect — resumed-monitoring order can be falsely established after later active work.
  In the detached audit, inserting `active_work(count=1)` after `active_scope_empty` and before the second wait passed `assertCodexWaitAgentSteering`; lines 445-452 latch the first empty observation and never invalidate it. Reset that observation on later active work and require a renewed empty-scope event before monitoring.
- FAILED: AC-1/AC-2 material evidence defect — a durable stage report with result bullets but no evidence lines is credited as authoritative completion.
  In the detached audit, deleting both evidence lines while retaining the two required `DONE` bullets passed; lines 553-563 validate bullet presence only. Require a non-empty, unfenced evidence line for each required result, or validate the complete stage-report record atomically.
- DONE: AC-3 captain-facing language verification.
  The live drive used two monitoring calls with zero captain-facing harness-label/disclaimer messages, and the planted repeated-disclaimer mutant fails if either forbidden phrase is accepted.
- DONE: AC-4 containment verification.
  Scoped diff review and the contractlint/full/race gates show no Claude, Pi, shared-core, scheduler, or runtime-API change; `git diff --check` is clean.

### Summary

The live Codex drive exhibited the intended outcome: captain steering resumed active work without cancelling or redispatching the unresolved worker, and monitoring resumed afterward. The shipped tests reject every planted mutant and all required suites pass, but the detached audit found two supported false-pass traces in ordering and durable-report authority. Both material findings affect AC-1/AC-2 evidence and require a narrow correction before validation can recommend PASSED; there are no deferred risks or polish-only findings.

## Stage Report: implementation (cycle 2)

- DONE: Invalidate the latched active-scope-empty observation whenever later active work occurs, and require a renewed empty-scope event before the second monitoring call.
  Commit `9e3a5ca0` resets the observed empty boundary on every later positive active-work event; monitoring without a renewed empty event is rejected while runtime behavior remains unchanged.
- DONE: Add the exact active-work-after-empty mutant and prove it fails at the stale-empty ordering boundary without changing runtime behavior or the durable report schema.
  The new focused mutant first passed before the fix, then failed with `captain-authorized active work was not exhausted before monitoring resumed`; all prior steering mutants and the original trace still pass.
- DONE: Run focused, containment, live-tag, full, race, and one targeted final Roborev pass; update the implementation report with the narrow disposition.
  Focused steering/idle tests, contractlint containment, live-tag compile, `go test ./...`, and `go test ./... -race` passed; targeted Roborev job 2419 found no issues.
  The final surface is 6 files/896 changed LOC versus the advisory 5 files/190–285 LOC estimate, a 10-line net increase confined to the existing Codex behavioral proof; AC remains unchanged.
  Validation's mandatory evidence-line proposal remains declined as an AC-narrowing report-format schema; parsed durable stage-report authority is unchanged.

### Summary

The correction closes the stale-empty ordering false pass by making later active work revoke the prior idle observation. A fresh empty-scope event is now required before monitoring can resume, with no runtime, adapter, or durable-report schema change.

## Stage Report: validation (cycle 2)

- DONE: Reproduce the active-work-after-empty mutant at corrected commit `9e3a5ca0`; verify it fails until a renewed empty-scope event precedes the second wait, while the original reduced trace still passes.
  A detached three-state audit passed: the original trace succeeds, later active work without renewed empty scope fails at `captain-authorized active work was not exhausted`, and inserting a new empty-scope event before the second wait succeeds.
- DONE: Confirm all prior mutants, focused, containment, live-tag, full, and race checks remain green and that the runtime adapter plus durable report schema are byte-unchanged by this correction.
  The original trace and all 36 planted mutants passed their assertions; focused steering/idle, contractlint containment, live-tag compile, `go test ./...`, and `go test ./... -race` passed. Diff and function-slice comparisons show only `internal/ensigncycle/codex_wait_agent_steering_test.go` changed; the adapter, evidence fixture, report artifact, and report parser are byte-unchanged.
- DONE: Re-enter AC-1 through AC-4 judgment, preserve the evidence-line decline unless an accepted AC is actually violated, and report PASSED or REJECTED.
  PASSED: AC-1 is covered by the successful live drive plus the corrected order matrix; AC-2 by all lifecycle, identity, epoch, authority, and ordering mutants; AC-3 by the repeated-disclaimer mutant and zero live captain-facing uses; AC-4 by the one-file correction and containment/full/race gates. Evidence lines are not an accepted completion-schema requirement, so the declined proposal does not affect any AC.
- DONE: Reviewer findings and release scope.
  The detached correction audit is clean; there are no material findings, deferred risks, or polish-only findings.

### Summary

Cycle-2 validation confirms that later active work now invalidates a stale empty-scope observation and that monitoring becomes valid again only after a fresh empty boundary. All acceptance evidence and repository gates pass without runtime, adapter, or durable-report changes. Recommendation: PASSED.
