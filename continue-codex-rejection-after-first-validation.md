---
title: Continue Codex rejection after the first validation
status: validation
source: "Staff review M2 for test-behavior-completeness, 2026-08-09"
started: 2026-08-09T20:36:16Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree: .worktrees/spacedock-ensign-continue-codex-rejection-after-first-validation
issue:
pr:
mod-block:
id: dvddbpsf4tdt3yjw1yjyp14k
gates:
    version: 1
    records:
        - id: gate:dvddbpsf4tdt3yjw1yjyp14k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:dvddbpsf4tdt3yjw1yjyp14k-backlog-1
              briefing:
                id: briefing:dvddbpsf4tdt3yjw1yjyp14k:backlog:attempt-1:revision-1
                digest: sha256:f8211760d0bf561280fa0671e01ccd07238ebcc2d984f118fe07f8b7a955ce67
                request-digest: sha256:aa4c58b93bc1e78c48359d493dcd20922bda0a810de46ab3d7dcc5308b93bd62
                room-ref: ./continue-codex-rejection-after-first-validation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:dvddbpsf4tdt3yjw1yjyp14k:backlog:1
                briefing: briefing:dvddbpsf4tdt3yjw1yjyp14k:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T20:35:44.798466Z"
                decision: approve
                reason: The Captain authorized shaping and requires end-user value; this task owns the complete Codex correction journey.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:dvddbpsf4tdt3yjw1yjyp14k:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:dvddbpsf4tdt3yjw1yjyp14k-ideation-1
              briefing:
                id: briefing:dvddbpsf4tdt3yjw1yjyp14k:ideation:attempt-1:revision-1
                digest: sha256:e7ea43c9d2cf658ba396a887006d1f6eb35c2553c608513a223228c5553a080d
                request-digest: sha256:d7130b947de519759c629aac7fa524b15d8d6b79131095e608fb7311f28f6ecf
                room-ref: ./continue-codex-rejection-after-first-validation/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-09T21:28:27.351902Z"
                reason: The ideation AC labels are not readable by the workflow AC extractor.
            - id: gate-attempt:dvddbpsf4tdt3yjw1yjyp14k-ideation-2
              briefing:
                id: briefing:dvddbpsf4tdt3yjw1yjyp14k:ideation:attempt-2:revision-1
                digest: sha256:d42af88b082ed9618b19409f52f7ae84861fdc41d901ab7d1e4bb5a2a4b2dfb8
                request-digest: sha256:cff06998837ee1db46e617fb191d8b13df6a0b01c2fa892e8fcd1aed997e27b2
                room-ref: ./continue-codex-rejection-after-first-validation/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:dvddbpsf4tdt3yjw1yjyp14k:ideation:2
                briefing: briefing:dvddbpsf4tdt3yjw1yjyp14k:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-09T21:33:20.653918Z"
                decision: approve
                reason: Captain approved the complete Codex correction journey and fresh final gate.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:dvddbpsf4tdt3yjw1yjyp14k:validation
          stage: validation
          attempts:
            - id: gate-attempt:dvddbpsf4tdt3yjw1yjyp14k-validation-1
              briefing:
                id: briefing:dvddbpsf4tdt3yjw1yjyp14k:validation:attempt-1:revision-1
                digest: sha256:38ae354a6e8f34d8f82c078252a2dffeff90656dc27322d3496e6e8aa67922e8
                request-digest: sha256:e2d4a6243cea4aba4bf58b3bc5c7a5fa833111e69bdaf4fa8201b62102ed6709
                room-ref: ./continue-codex-rejection-after-first-validation/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:dvddbpsf4tdt3yjw1yjyp14k:validation:1
                briefing: briefing:dvddbpsf4tdt3yjw1yjyp14k:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T17:52:00.369733Z"
                decision: approve
                reason: Captain conn approves dvd after independent PASSED validation of exact candidate 29a4dd5dc and the target-level XFAIL disposition.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:dvddbpsf4tdt3yjw1yjyp14k-validation-2
              briefing:
                id: briefing:dvddbpsf4tdt3yjw1yjyp14k:validation:attempt-2:revision-1
                digest: sha256:ca626f734a1f32129070102a853cc50fffffc0cc5f35018ba828529fca6e0f81
                request-digest: sha256:29a8c2ed9bf2e2a6feb2ea10012a4411dab48c29f44fd436e84ee20e28ff86cd
                room-ref: ./continue-codex-rejection-after-first-validation/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:dvddbpsf4tdt3yjw1yjyp14k:validation:2
                briefing: briefing:dvddbpsf4tdt3yjw1yjyp14k:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-10T18:36:18.337897Z"
                decision: approve
                reason: Captain conn approves dvd after independent PASSED cycle-2 validation of exact candidate 3d0c02c194; required PR lanes remain the merge boundary.
              application:
                target-stage: done
                state: superseded
---

## Problem

The Codex rejection-flow cell does not prove a complete correction journey.
The current declaration skips Codex with a `liveTODO` entry. The first live
spike removed only that skip in an isolated worktree. It ran this real command:

```bash
SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_LIVE_REQUIRED=1 \
SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-codex-rejection-live-artifacts-1d6Hgc \
go test -tags live ./internal/ensigncycle \
  -run '^TestLiveCommonRejectionFlow$' -v -count=1
```

Codex authenticated with the local login and ran for 434.41 seconds. The test
ended with this error:

```text
resolved launcher never invoked `gate record --round validation/1`
```

The error came from the shared runner calling `claudeRecordedRejectionRound`
against Codex JSONL. The Codex transcript contains two successful
`gate record --round` calls. `codexRecordedRejectionRound` finds both calls.
This proves that a Claude parser swap is needed for evidence, but it does not
complete the task.

The same transcript contains zero successful `gate prepare` calls. Its durable
fixture state has two implementation reports, two validation reports, the
exact marker `shared-rejection-fix: applied`, status `validation`, and review
round cycle 2. It has no fresh open gate. Codex ended with a nonterminal state
and said that it had stopped after recording validation/2.

This spike found two separate gaps:

1. The shared evidence path cannot grade the Codex command shape.
2. The Codex correction path records validation/2 but does not re-enter the
   normal gate lifecycle and prepare one fresh open gate.

An extractor-only repair would fix the first error while leaving the operator
without the required final gate. A recorder-only repair has the same defect.
The task value is the full Codex journey: rejected candidate, correction,
validation/2, and one fresh final gate.

## Visible value

After this task, a Codex operator can follow one rejected candidate through its
correction. The operator sees a second validation pass and one fresh open
validation gate for the corrected candidate. The gate remains unresolved. No
decision, application, terminal status, or successor dispatch occurs.

The independent baseline is the live spike above. It reached two validation
reports but prepared zero gates and ended in nonterminal `validation` state.
The target outcome is measurable on disk and in the Codex transcript:

- two implementation reports;
- two validation reports, with validation/1 rejected and validation/2 passed;
- one retained validation/1 advisory round with four entries and its canonical
  two-file review room;
- the exact correction marker;
- one successful `gate prepare` after validation/2;
- one fresh open validation gate bound to the prepared Briefing; and
- no gate resolution, application, verdict, completion, or successor dispatch.

The final-gate count moves from zero in the baseline to one in the candidate.
That count is the end-value measure. The report and round counts prove that the
gate is reached through the complete journey, not through a shortcut.

## Proposed approach

### Host-specific evidence

Keep the shared fixture and exercise. Route the rejection-round assertion
through the existing host shape. Claude uses `claudeRecordedRejectionRound`.
Codex uses `codexRecordedRejectionRound` and its observed command list. Keep
the existing Codex reviewer evidence helper as a secondary assertion when the
runtime emits a structured reviewer identity. Do not require a new identity
format for this task. Durable validation/2 state and the second reviewer run
remain the primary independence proof.

Add a small rejection-flow evidence capability to the existing live driver
surface. Select it from the driver that already owns the host transport. Do
not add a scenario registry or a second Codex journey. Keep Claude and Pi on
the current shared path.

The simplest alternative is to use the Claude extractor for every host. The
spike proves that this gives a false missing-invocation error for Codex. A
second alternative is a new Codex-only scenario. That duplicates the fixture,
prompt, metrics, and durable assertions. The existing driver capability is the
smallest mechanism that serves the evidence criterion.

### Complete Codex correction journey

After the Codex validation reviewer passes on the corrected candidate, the
Codex adapter and the feedback-rejection instructions must state the next
actions explicitly:

1. Read the completed validation/2 report and the workflow-owned Cycle line.
2. Re-enter the existing First Officer gate lifecycle.
3. Commit one gate-review artifact for the corrected candidate.
4. Invoke the existing neutral `gate prepare` path once.
5. Read the resulting durable gate state and present the fresh open gate.
6. Stop without consuming the gate or setting a decision.

The existing `gate record --round validation/2` call remains an advisory round
publication. It is not the final gate. The final gate must use the prepared
Briefing `briefing:rejection-task:validation:attempt-1:revision-1`, not the
validation/1 advisory Briefing. The existing `gates` authority and lifecycle
own the record. No new gate command, field, or authority is needed.

The simplest alternative is to leave “re-enter the normal gate flow” as the
only instruction. The live spike proves that this wording permits Codex to
stop after `gate record`. The explicit sequence is required for the visible
value. A new gate command or an automatic gate resolution is broader than the
value and is out of scope.

### Durable final-gate oracle

Keep the existing round oracle and add a strict final-gate assertion for this
journey. The strict assertion must reject an absent gate, an advisory round
used as a gate, a missing prepare command, a closed gate, a second attempt, a
decision, an application, or a terminal state. Preserve the current offline
round test that accepts a round-only state. The Codex live journey calls the
strict form after it calls the common round oracle.

The simplest alternative is to make `assertRejectionRoundGateBoundary` always
require a gate. That would change the existing round-only contract and would
break the control that proves `gate record --round` can retain advisory state
without preparing a gate. The optional strict form keeps that boundary while
measuring the Codex value.

### Non-product baseline and strict-XFAIL source binding

After the known `ts` target-XFAIL machinery and the preceding `zh` recorder
work land, make one non-product baseline commit. This commit contains only
the Codex host extractor, the strict final-gate oracle, its negative controls,
and the target binding. It contains no skill text, Codex runtime instruction,
or user-visible documentation change.

The baseline uses the Codex `rejection-flow` target binding:

```text
target: codex
owner: dvddbpsf4tdt3yjw1yjyp14k
```

Run the complete Codex rejection journey from that baseline commit. It must
not skip. The Codex host extractor must select
`codexRecordedRejectionRound`, and the strict final-gate oracle must reach the
durable missing-gate state. A typed semantic failure is target-level XFAIL.
A parser error, auth error, launch error, timeout, or malformed state is an
ordinary failure.

Do not change skills, Codex runtime instructions, or user-visible behavior
until the baseline has produced a target-level XFAIL. This order keeps
the baseline proof independent from the repair.

With the binding still present, the repaired exact candidate must produce an
XPASS and fail the lane. Remove the binding only after that XPASS. Run the same
candidate again and require a PASS with no semantic code. The candidate keeps
the Sonnet, Opus, and Pi ownership rows unchanged.

## Acceptance criteria

**AC-1 — The Codex baseline executes the rejection-flow journey and records a typed semantic XFAIL.**

The cell does not skip. The focused live command, the strict outcome record,
and the Codex artifact prove this. Any parser, process, auth, timeout, or
second semantic error fails the baseline instead of becoming XFAIL.

**AC-2 — The repaired exact candidate leaves the fixture in a complete correction state.**

The state contains two implementation reports, two validation reports,
validation/1 rejected, validation/2 passed, the exact fix marker, the unchanged
validation/1 four-entry round, and the unchanged canonical two-file room. The
live assertion reads the entity, reports, round summary, room bytes, and
candidate bytes. Negative controls remove the second report, change the
marker, change the round, or alter a preserved lifecycle sentinel, and each
control must fail.

**AC-3 (VALUE) — The repaired exact candidate leaves exactly one fresh open validation gate after validation/2.**

The gate has the prepared Briefing ID, one attempt, no resolution, no
application, and no terminal entity fields. The Codex transcript has one
successful `gate prepare` after the validation/2 record. The strict gate
oracle rejects an absent gate, an advisory Briefing, a closed gate, a duplicate
attempt, a decision, an application, or a prepare before validation/2.

**AC-4 — The Codex result proves the host-specific and independent parts of the journey.**

The assertion finds the Codex `command_execution` round record, finds the
second validation run, and rejects an implementation worker serving as its own
validator when structured reviewer data is present. The existing Codex
reuse-or-fresh helper supplies the reviewer evidence. The durable validation/2
report remains required even when the host does not expose a structured
identity.

**AC-5 — The strict-XFAIL sequence is complete.**

The repaired candidate is an XPASS while the binding remains. The same
candidate is a PASS after the binding is removed. The two results come from
the same focused target and preserve the Codex artifact and metric record.

## Expected surface and semantic budget

The expected surface is six existing files. The first three files form one
non-product baseline commit. The last three files change only after the live
baseline produces a target-level XFAIL.

- `internal/ensigncycle/claude_live_runner_test.go`: route rejection evidence
  by host and call the strict Codex final-gate oracle.
- `internal/ensigncycle/shared_round_recording_test.go`: add the strict oracle
  and its malformed, closed, duplicate, and pre-round controls.
- `internal/ensigncycle/shared_live_runner_test.go`: bind Codex to the strict
  XFAIL owner and code after the dependency lands.
- `skills/feedback-rejection-flow/SKILL.md`: state the post-review gate
  preparation and stop point after the baseline evidence.
- `skills/first-officer/references/codex-first-officer-runtime.md`: bind the
  explicit Codex correction-to-gate sequence after the baseline evidence.
- `docs/runtime-live-ci.md`: document the Codex journey result and the strict
  XFAIL sequence after the baseline evidence.

Estimate about 46 gross insertions and 20 deletions, for a net increase of
about 26 lines. Tolerance is plus or minus one file and plus or minus 14 net
lines. The extra budget records the live spike's two-boundary result. A
component-only extractor diff does not meet this estimate or the acceptance
criteria.

The semantic budget is:

- Command grammar: unchanged. Use the existing `gate record --round` and
  `gate prepare` commands.
- Stored formats: unchanged. Keep the existing round, review-room, and gate
  records. Do not change recorder bytes.
- Authority: unchanged. The First Officer prepares the gate. The Captain
  resolves it. No worker gains gate authority.
- Runtime behavior: Codex now continues from a passed validation/2 review to
  one fresh open gate and stops there. It does not decide or terminalize.
- Other runtimes: unchanged. Do not change the Pi recorder-order repair or the
  Sonnet and Opus rows.

## Spike record and proven mechanisms

The live spike used an isolated worktree and changed only its local test TODO.
It did not change the product checkout or the state entity. The artifact was:

```text
/tmp/spacedock-codex-rejection-live-artifacts-1d6Hgc/codex-shared-scenarios/rejection-flow/codex-exec.jsonl
```

The transcript proved two successful Codex `gate record` commands, for
validation/1 and validation/2, and zero successful `gate prepare` commands.
It also proved that `commandRecordsRejectionRound` and
`codexRecordedRejectionRound` match the real Codex shell command. The existing
offline test `TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl`
already drives `gates.Prepare` and proves the prepared open-gate bytes. That
mechanism needs no new format spike. The initial live run is diagnostic only:
the Claude extractor stopped the test before a semantic final-gate result, so
it is not the strict-XFAIL proof. The non-product baseline commit must close
that evidence gap. The unverified mechanism is the Codex post-review handoff.
The exact-candidate live run is the first implementation test for it.

## Test plan

The test has medium code complexity and high live-runtime cost. The observed
baseline took about seven minutes. Use an isolated artifact directory and the
documented 40-minute Codex timeout.

1. Rebase after the `ts` strict-XFAIL work and the preceding `zh` recorder
   work. Add the Codex host extractor, strict final-gate oracle, negative
   controls, and XFAIL binding in one non-product baseline commit. Do not
   change skills or Codex runtime instructions.
2. Run the exact Codex command from that baseline commit. Require the complete
   journey and a typed semantic XFAIL. A parser or infrastructure failure stops
   the repair. Save the baseline JSONL,
   final message, entity, reports, round room, and gate evidence.
3. Only after the target-level XFAIL, add the feedback skill and Codex runtime
   instruction changes. Add the documentation diff. Run the focused offline
   round, extractor, strict-oracle, and negative-control tests.
4. Run the exact Codex command with the binding retained. Require XPASS. Save
   the repaired JSONL, final message, entity, reports, round room, and gate
   room.
5. Remove only the Codex binding. Run the same command again. Require PASS and
   inspect the final durable state for the complete correction journey.
6. Run the normal focused and full Go tests, the race suite, registry
   reconciliation, and `gofmt -w ./cmd ./internal` in the implementation
   stage. Do not use a Claude or Pi live run as a substitute for the Codex
   proof.

## Proposed documentation diff

In `docs/runtime-live-ci.md`, add the following text after the Codex common
journey command. This records the user-visible stop state without changing the
CLI surface:

```diff
@@
 SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommon' -failfast ./internal/ensigncycle -v
+The `rejection-flow` Codex journey is admitted through a strict-XFAIL
+baseline with `rejection-flow-not-completed`. A repaired candidate must show
+the rejected first validation, correction, passed validation/2, and one fresh
+open validation gate. The gate remains unresolved. A parser, launch, auth,
+or timeout error is not an expected semantic outcome.
```

## Out of scope

- recorder bytes, round format, `gate record --round` grammar, or gate schema;
- the Pi recorder-order repair and its source binding;
- a new Codex scenario, worker protocol, reviewer identity format, or gate
  authority;
- automatic gate resolution, application, terminalization, or successor
  dispatch;
- changes to the `ts` strict-XFAIL classifier beyond the target binding;
- product implementation in this ideation stage.

## Stage Report: ideation

- DONE: Exercise the real Codex rejection-flow failure before selecting a
  repair. The isolated live command ran for 434.41 seconds. Its Codex JSONL
  contained successful validation/1 and validation/2 records but no
  `gate prepare`; the shared assertion used the Claude extractor.
- DONE: Define one complete correction journey with strict-XFAIL and exact-
  candidate proof. The plan binds Codex to
  `rejection-flow-not-completed`, requires correction and validation/2, then
  requires one fresh open gate. It defines XPASS with the binding and PASS
  after its removal.
- DONE: Give a visible-value statement and gross and net line estimates. The
  baseline has zero prepared gates. The candidate must have one fresh open
  gate, two implementation reports, two validation reports, and the retained
  four-entry round. The estimate is 46 gross insertions, 20 deletions, and a
  net increase of 26 lines across six files.

### Summary

The ideation record assigns one complete Codex correction journey. It repairs
both the host-specific evidence boundary and the missing final-gate handoff.
It preserves recorder bytes, round format, authority, and the Pi repair. The
strict-XFAIL and exact-candidate sequence measures the final open gate rather
than accepting a component-only helper.

## Stage Report: ideation (cycle 2)

- DONE: Correct the baseline order after staff re-review. The Codex extractor,
  strict final-gate oracle, negative controls, and target binding now land in
  one non-product baseline commit. The live baseline runs only after that
  commit. Skill text, Codex runtime instructions, and user-visible docs wait
  for the baseline result.
- DONE: Require the complete Codex journey to produce a target-level XFAIL.
  The baseline must reach a durable semantic failure through validation/2.
  A parser or infrastructure error fails the baseline.
- DONE: Preserve the complete correction value and the exact-candidate proof.
  After the target-level XFAIL, the repair produces XPASS with the binding and
  PASS after its removal. The final state still requires one fresh open gate,
  two implementation reports, two validation reports, and the retained
  four-entry advisory round.

### Summary

The material ordering error is resolved. Test-only Codex evidence establishes
the honest strict-XFAIL baseline before any skill or runtime behavior changes.
The task still owns the complete rejected-candidate correction journey and its
fresh final gate.

## Stage Report: ideation (cycle 3)

- DONE: Restore extractor-visible acceptance criteria without changing their
  content or semantic scope. The five criteria now use exact bold labels with
  em dashes. The final fresh-gate criterion is `AC-3 (VALUE)`.
- DONE: Run the DVD acceptance-criteria scan. The command returned:
  `ac=AC-1 line=222 unevidenced=false citations=1`,
  `ac=AC-2 line=228 unevidenced=false citations=1`,
  `ac=AC-3 line=238 unevidenced=false citations=2`,
  `ac=AC-4 line=246 unevidenced=false citations=1`, and
  `ac=AC-5 line=255 unevidenced=false citations=1`. All five ACs are now
  discovered by the workflow extractor with evidence references.
- DONE: Preserve the complete correction journey and the strict evidence order.
  The non-product baseline still precedes the skills and runtime changes, and
  the fresh open final gate remains the measured end-user value.

### Summary

The AC structure now matches the DVD extractor. The scan finds AC-1 through
AC-5, including the VALUE final-gate criterion. No acceptance-criteria or
semantic requirement changed.

## Implementation Work Plan

The exact code base is `e2f07a40e6bf45eeec8d133723b25cef18f5cb9a`.
The approved surface contains these six existing files:

- `internal/ensigncycle/claude_live_runner_test.go`
- `internal/ensigncycle/shared_round_recording_test.go`
- `internal/ensigncycle/shared_live_runner_test.go`
- `skills/feedback-rejection-flow/SKILL.md`
- `skills/first-officer/references/codex-first-officer-runtime.md`
- `docs/runtime-live-ci.md`

The estimate is 46 additions and 20 deletions. The net estimate is 26 lines.
The tolerance is one file and 14 net lines.

A Codex operator gets one complete correction journey. The journey ends at one
fresh unresolved validation gate after validation/2.

## Implementation Finding

The exact baseline SHA is `06e87a3cb7f5608afe6b12e518f10de3642ea6c8`.
The local Codex run produced two semantic codes. The expected code was
`rejection-flow-not-completed`. The additional code was
`rejection-reviewer-flow`.

- Released user and normal workflow: the subscription-backed Codex rejection
  journey reached validation/2.
- Observable harm: the extra code prevents the required single typed XFAIL.
- Affected value: `value-ac[AC-4]` requires structured identity only when Codex
  emits that identity.
- Trigger evidence: the transcript contains only `wait` collaboration events.
  It contains no structured worker spawn or follow-up event. The validation/2
  report is durable and PASSED.

The proposal is Material. This task owns the finding. The proposed disposition
is a test-only fix for the absent structured-identity case. Candidate bytes
remain unchanged until the First Officer gives a distinct authorization.

The First Officer authorized the fix. The correction changes only
`internal/ensigncycle/claude_live_runner_test.go`. The estimate is four added
lines and no deleted lines. The correction accepts only
`errReviewerIdentityUnsupported`. A detected wrong reviewer route remains red.

## Implementation Design Reset

The exact baseline changes these two files:

- `internal/ensigncycle/claude_live_runner_test.go`
- `internal/ensigncycle/shared_round_recording_test.go`

The current baseline has 87 additions and 24 deletions. The net increase is 63
lines.

The remaining repair changes these three files:

- `skills/feedback-rejection-flow/SKILL.md`
- `skills/first-officer/references/codex-first-officer-runtime.md`
- `docs/runtime-live-ci.md`

The remaining estimate is 20 additions and two deletions. The projected full
change has 107 additions and 26 deletions. The projected net increase is 81
lines across five existing files. The proposed revised tolerance is 81 plus or
minus 10 net lines and zero additional files.

The end-user value is one complete Codex correction journey. Codex records the
rejected first validation, applies the correction, passes validation/2, and
stops at one fresh unresolved validation gate.

The strict oracle must prove these independent facts:

- Codex records the validation/1 advisory round.
- Codex reaches a durable validation/2 report.
- Codex runs exactly one successful gate preparation after validation/2.
- The durable gate has one open attempt and the prepared Briefing.
- A decision, application, duplicate attempt, terminal state, or wrong reviewer
  route remains red.
- An absent structured reviewer identity does not create false evidence.

The current oracle is the smallest implementation that keeps these controls.
It reuses the existing gate reader, round oracle, Codex command parser, and
reviewer classifier. It does not add a fixture, command, format, or authority.

A smaller alternative can omit transcript ordering and accept the durable gate
alone. This alternative removes the successful-command parser and its controls.
It loses AC-3 value because it cannot prove that preparation occurred once and
after validation/2. It can accept a failed, duplicated, or early command path.

## Implementation Finding: Current Round Pointer

The repaired candidate is `274c6b1f2768eb43b79db891a23f0e53b94c12b6`.
The exact Codex run published validation/1 and validation/2. It then prepared
one fresh open gate. The oracle rejected this correct state because it asked
`ValidateRoundFile` to resolve validation/1 through the current validation/2
pointer.

The First Officer classified this finding as Material and owned by this task.
The authorized correction changes only
`internal/ensigncycle/shared_round_recording_test.go`. The estimate is 18
additions and three deletions, for a net increase of 15 lines. The projected
full net increase is 87 lines. This value is inside the approved 71-to-91-line
range and adds no file.

The correction validates the current validation/2 summary when that pointer
exists. It keeps baseline validation/1 behavior. It also keeps exact byte checks
for the validation/1 room. Focused controls reject a missing or malformed
validation/1 room and an invalid current validation/2 summary.

The end-user value does not change. Codex completes the correction, publishes
validation/2, prepares one fresh unresolved gate, and stops.

## Implementation Finding: Gate Artifact Location

The unbound run used candidate
`4c1bc36e44df8638a466e6ade42a72f6955a4341`. Codex completed both validation
rounds and prepared one open gate. It placed `gate-review.md` inside the
immutable validation/2 advisory room. The strict oracle correctly rejected the
extra room file.

The First Officer classified this finding as Material and owned by this task.
The authorized correction replaces one existing line in each product file:

- `skills/feedback-rejection-flow/SKILL.md`
- `skills/first-officer/references/codex-first-officer-runtime.md`

The estimate is two additions and two deletions. The net change is zero. The
full candidate stays at the approved 91-line net limit and adds no file.

The lines require the gate-review artifact to stay outside every
`review/STAGE/round-CYCLE` advisory room. The oracle continues to reject any
extra file in an immutable round room. The end-user value remains one fresh
open gate after a valid validation/2 publication.

## Stage Report: implementation

- DONE: Rebase from exact zh landing e2f07a40e6bf45eeec8d133723b25cef18f5cb9a and preserve Pi and all non-dvd bindings.
  The final rebase target is `4dc83c0f8`. Candidate `29a4dd5dc` keeps Pi owner `p17swb3375rt525fn7f8xt7e` and all non-dvd rows.
- DONE: Before product edits, update the task report with exact files, gross/net estimate, tolerance, and Codex operator value.
  The Implementation Design Reset records five planned files, +107/-26, net +81 plus or minus 10, and the complete Codex journey.
- DONE: Build the non-product Codex extractor, strict final-gate oracle, negative controls, and dvd Codex XFAIL baseline only.
  Commits `6ecd0abf7` and `ba024e403` add the extractor, strict oracle, absent-identity rule, and red reviewer controls.
- DONE: Run the exact local subscription-backed Codex rejection-flow baseline and retain typed target-level XFAIL artifacts; stop on parser, auth, launch, timeout, or a second semantic failure.
  `/tmp/spacedock-dvd-codex-baseline-local.3wLIwh` ran 395.28s and produced the two recorded typed codes. Its SHA-256 is `abede885b2bd575f61bff10dc9f05e850af3b9a3be89210479a0a699e496c050`.
- DONE: After valid baseline evidence, implement only the approved feedback and Codex runtime instructions plus documentation.
  Commits `cf24e4030` and `fb9a41507` change the approved skill, Codex runtime, and behavior documentation.
- DONE: Prove the complete two-validation correction state and exactly one fresh unresolved validation gate with all required negative controls.
  The oracle proves validation/1, validation/2, one prepare, one open attempt, the Briefing, correct rooms, and all recorded red controls.
- DONE: Run the exact repaired Codex target with binding retained and retain XPASS evidence, then remove only dvd Codex binding and retain normal PASS evidence.
  First XPASS `/tmp/spacedock-dvd-codex-repaired-xpass-2.Z4shfn` ran 574.99s with SHA-256 `c845f4eaee450eb00caa50028e20a95b4f6664c441ca96f40c0d8c680e769cd0`. Adjusted XPASS `/tmp/spacedock-dvd-codex-adjusted-xpass.4AbVtR` ran 554.31s with SHA-256 `ec484832ca840edccfe43418b54cff464ae1dc708d84557bbae808867f86d4e3`. Final PASS `/tmp/spacedock-dvd-codex-adjusted-pass.YvjTPg` ran 535.78s with SHA-256 `2c71f508fc2da50081cc8cc2db72907986a4bedfe3b255921d6fcf6df3b63b55`.
- DONE: Do not run Pi, Opus, or Sonnet. Preserve the Pi XFAIL and all unrelated ownership rows.
  No Pi, Opus, or Sonnet run occurred. The post-rebase registry and active-owner controls passed in 0.416s.
- DONE: Run focused tests, gofmt, registry, active-owner, go test ./..., and go test ./... -race.
  All named commands passed. `internal/ensigncycle` took 366.333s normally and 556.475s with the race detector.
- DONE: Keep within six existing files, plus/minus one file, and net +26 plus/minus 14 lines. Stop for Captain approval before exceeding tolerance or changing grammar, formats, authority, or other runtime behavior.
  The Captain approved the reset before mutation. The final seven-file change has 122 additions, 31 deletions, and net +91.
- DONE: Write a Simplified Technical English implementation report with exact commits, commands, timings, artifact paths, digests, checklist status, and AC evidence.
  This report gives commits, commands, times, four artifact paths and digests, checklist status, and falsifiable acceptance evidence.
- DONE: Commit and push the exact product candidate and durable report, then send a completion signal that starts with "We love you."
  Candidate `29a4dd5dc` is pushed. The path-scoped report commit and completion signal finish this item.

### Summary

Codex records validation/1, applies the correction, and records validation/2. It prepares one fresh open gate and stops without a successor.
The pushed product candidate stays unchanged at `29a4dd5dc440d8fbc34cb3b635396ee8040144c6`.

## Stage Report: validation

- DONE: Confirm exact candidate 29a4dd5dc440d8fbc34cb3b635396ee8040144c6, remote match, clean worktree, and base 4dc83c0f8.
  Local HEAD and the remote branch match `29a4dd5dc`. Base `4dc83c0f8` is an ancestor, and the worktree is clean.
- DONE: Verify the exact seven-file, +122/-31, +91 net surface stays within the approved five-file reset plus previously authorized oracle files and adds no unapproved file.
  The diff contains seven approved existing files and no new file.
- DONE: Inspect baseline XFAIL, adjusted bound XPASS, and final unbound Codex PASS artifacts and SHA-256 values.
  The JSONL digests are `abede885b2bd575f61bff10dc9f05e850af3b9a3be89210479a0a699e496c050`, `ec484832ca840edccfe43418b54cff464ae1dc708d84557bbae808867f86d4e3`, and `2c71f508fc2da50081cc8cc2db72907986a4bedfe3b255921d6fcf6df3b63b55`. Each process reached its terminal event.
- DONE: Verify validation/1 and validation/2 publications, exactly one fresh unresolved gate, and the gate-review artifact outside immutable advisory-round rooms.
  The final transcript has both successful publications, one successful prepare, one open attempt, and an external gate-review artifact.
- DONE: Verify no-reuse, wrong-reviewer, duplicate prepare, malformed round, invalid summary, terminal state, and extra-round-file negative controls.
  Focused controls passed. Temporary probes also rejected terminal fields and an extra validation/2 room file.
- DONE: Verify only dvd's Codex rejection-flow binding and mirrored row were removed; preserve Pi and every unrelated binding.
  The base diff removes only owner `dvddbpsf4tdt3yjw1yjyp14k`. Pi owner `p17swb3375rt525fn7f8xt7e` remains.
- DONE: Run focused reviewer, round, gate, registry, active-owner, binding, format, and diff checks independently.
  The focused Go checks passed. The format and diff checks produced no error.
- DONE: Inspect the implementation full/race evidence without duplicating owned full/race/live runs.
  The retained report records normal completion in 366.333s and race completion in 556.475s. No duplicate run occurred.
- DONE: Report every finding with released user, observable harm, value authority, and exact trigger evidence.
  One proposed evidence finding entered disposition. The First Officer declined it under the Captain ruling from 2026-08-10.
- DONE: Write and push a Simplified-English validation report with PASSED or REJECTED recommendation.
  PASSED: all value criteria have applicable evidence, and no Material finding remains.

### Finding disposition

- Released user and normal workflow: The Codex operator relies on the target-level XFAIL admission sequence.
- Observable harm: The validator proposed that two semantic codes left the baseline without valid admission evidence.
- Value authority: `value-ac[AC-1]` requires a typed target-level XFAIL and rejects infrastructure failures.
- Trigger evidence: Artifact `abede885…` ran source `06e87a3c` and produced the expected target-level XFAIL plus a reviewer-identity code.
- First Officer disposition: DECLINE. The Captain permits multiple semantic codes for dvd, and the artifact contains no infrastructure failure.

### Summary

The exact candidate preserves the complete Codex correction journey. It removes only the dvd Codex binding.
The retained final run stops at one fresh unresolved gate.
The validation recommendation is PASSED. No Material finding remains after the recorded First Officer disposition.

## Validation Finding: Launcher Variable

- Released user and workflow: The exact Codex `rejection-flow` uses candidate `29a4dd5dc`.
- Observable harm: The oracle reports `rejection-flow-not-completed` for a complete journey. This false result blocks required PR evidence.
- Value authority: `value-ac[AC-3]` requires validation/2 before exactly one fresh unresolved gate.
- Trigger evidence: Run `31416271663`, job `93546268920`, and artifact `9074101972` contain the failing transcript.

The transcript contains validation/1, validation/2, one gate prepare, an open
validation gate, and the terminal final message. The successful validation/2
command uses `launcher=${SPACEDOCK_BIN:-spacedock}; "$launcher" gate record ...`.
The `rejectionValidation2Command` expression does not accept `$launcher`.

The First Officer classified this finding as Material and owned by dvd.
The disposition is FIX. The correction changes only
`internal/ensigncycle/shared_round_recording_test.go` with a same-line
replacement. The planned net change is zero, so the seven-file net +91 cap
does not change. Product instructions and bindings stay unchanged.

The focused positive uses the exact `$launcher` command. Existing controls
continue to reject a wrong entity, round, file, or failed command. Cancelled
and superseded PR live evidence stays retained.

Commit `c334febb28d45705362d3485abbda6dc7993f47f` adds the `$launcher`
alternative and changes the existing positive to the exact command shape.
The focused test first failed with `stage/prepares = 0/1`. It then passed with
the retained negative controls. The correction changes two lines in and two
lines out of one existing file. The full candidate remains net +91.

## Stage Report: implementation (cycle 2)

- DONE: Record the exact corrected candidate head.
  The corrected candidate head is `3d0c02c1947824d30e3cac3f874ee179092ee50d`.
- DONE: Keep the correction in one existing file with a net-zero diff.
  `internal/ensigncycle/shared_round_recording_test.go` changes by +2/-2, so the correction is net zero.
- DONE: Prove the focused failure before the correction and the focused pass after it.
  The focused control first failed with `stage/prepares = 0/1`. The same control passed after the regex correction.
- DONE: Preserve the wrong entity, wrong round, wrong file, and failed-command controls.
  The focused extractor run passed with all four negative controls still active.
- DONE: Keep product instructions, bindings, and the approved net cap unchanged.
  The correction changes no product or binding file. The full seven-file candidate remains net +91.
- DONE: Retain the superseded PR live evidence.
  Run `31416271663` and artifact `9074101972` remain in the durable finding record.
- DONE: Push the exact corrected candidate to the remote branch.
  Local HEAD and the remote branch both identify `3d0c02c1947824d30e3cac3f874ee179092ee50d`.
- DONE: Bound the validation/2 evidence correction before mutation.
  `internal/ensigncycle/shared_round_recording_test.go` has an estimated +10/-10 change. Net +91 and the seven-file tolerance stay unchanged. Users get exact validation/2 evidence with canonical files.
- DONE: Require canonical validation/2 artifacts and preserve all focused controls.
  Commit `3d0c02c1947824d30e3cac3f874ee179092ee50d` is +10/-10. Wrong briefing, log, entity, round, and failed-command controls pass.

### Summary

The validator accepts the exact `$launcher` validation/2 command only with canonical briefing and log files. All negative controls remain active.
Product instructions and bindings do not change.
The candidate remains within the approved net +91 cap and retains the superseded PR evidence.

## Stage Report: validation (cycle 2)

- DONE: Confirm exact corrected candidate 3d0c02c1947824d30e3cac3f874ee179092ee50d, remote match, clean worktree, and base 4dc83c0f8.
  Local HEAD and the remote branch match `3d0c02c`. Base `4dc83c0f8` remains an ancestor, and the worktree is clean.
- DONE: Verify the cycle-2 correction changes only one existing oracle file by +10/-10 and keeps the full candidate seven files and +91 net.
  The bounded correction is net zero. The full candidate is +127/-36 across seven existing files.
- DONE: Inspect run 31416271663, job 93546268920, artifact 9074101972 and confirm the product journey completed validation/1, validation/2, one prepare, and one open gate.
  The retained transcript contains both successful publications before one successful prepare. Its final state is open and nonterminal.
- DONE: Verify the corrected oracle accepts the exact `$launcher` successful command and rejects wrong briefing, log, entity, round, and failed command controls.
  The focused oracle tests passed in 0.894s. Each mutation removes the required `stage=2, prepares=1` result.
- DONE: Verify product instructions, bindings, Pi ownership, and all unrelated bytes are unchanged from validated candidate 29a4dd5dc.
  The correction changes only `internal/ensigncycle/shared_round_recording_test.go`. Pi owner `p17swb3375rt525fn7f8xt7e` remains.
- DONE: Inspect focused pre-fix failure and post-fix pass evidence without repeating full, race, or live runs.
  The retained failure was `stage/prepares = 0/1`. The corrected focused tests passed. No full, race, or live run occurred.
- DONE: Report every finding with released user, observable harm, value authority, and exact trigger evidence.
  The wrong-file evidence defect entered disposition before rerun. The bounded correction closes it, and no new Material finding remains.
- DONE: Write and push a Simplified-English validation cycle-2 report with PASSED or REJECTED recommendation.
  PASSED: the corrected oracle matches the retained journey and rejects all required adjacent command variants.

### Finding disposition

- Released user and normal workflow: A Codex operator relies on the strict validation/2 command oracle for PR evidence.
- Observable harm: Candidate `c334febb` accepted a command that used a wrong briefing or log file.
- Value authority: `value-ac[AC-3]` requires a successful validation/2 publication before exactly one fresh gate.
- Trigger evidence: The changed regex lacked file constraints. Its named wrong-file control exercised the separate validation/1 recognizer.
- First Officer disposition: FIX. Commit `3d0c02c` requires both canonical files and adds five direct negative controls.

### Summary

The exact candidate matches the remote branch and preserves the complete Codex correction journey. The bounded oracle correction rejects every required adjacent variant.
The validation recommendation is PASSED. No Material finding remains.

## Stage Report: implementation (cycle 3)

- DONE: Read the exact Sonnet artifact and Captain parser disposition.
  Artifact `9076566633` shows validation/1 with `entries=2`, validation/2, one external gate-review artifact, and one open gate.
- DONE: Change only the validation/1 result count from entries=4 to entries=2.
  Commit `d02dd14096bc7a5a87f1b8e3f8ae53399ee9f01a` makes the exact count change in `internal/ensigncycle/shared_round_recording_test.go`.
- DONE: Add the exact Sonnet multiline positive and malformed-round negative in the same test file.
  The focused extractor test accepts the multiline command and three-line result. It rejects the same result after its round identifier becomes malformed.
- DONE: Keep all product instructions, bindings, and other oracle behavior unchanged.
  The correction changes one test file by +4/-4. The full candidate has seven existing files, +131/-40, and net +91.
- DONE: Run focused, full, race, format, registry, and owner checks for the changed test bytes.
  Focused and contract checks passed. Normal and race suites passed. `internal/ensigncycle` took 283.091s and 271.190s.
- DONE: Commit and push the corrected candidate and a Simplified-English implementation cycle report.
  Local and remote code heads match `d02dd14096bc7a5a87f1b8e3f8ae53399ee9f01a`. No live runtime test ran before fresh validation.

### Summary

The parser now accepts the exact Sonnet validation/1 result with two entries. The matching malformed result remains red.
The correction changes no product instruction or binding. The full candidate remains within the approved seven-file and net +91 limits.

## Stage Report: implementation (cycle 4)

- DONE: Rebase exact candidate d02dd14096bc7a5a87f1b8e3f8ae53399ee9f01a onto immutable origin/main 0bbe9d46c02328930253bfbe619f9827d6da5109.
  Rebased candidate `47d39b624edbdb9590b9fbabf7c8c9dd27534989` has `0bbe9d46c02328930253bfbe619f9827d6da5109` as its merge base.
- DONE: Stop on any conflict except the two known registry and binding files.
  Conflicts occurred only in the registry file. The binding file merged automatically during the authorized commits.
- DONE: Keep all current main Opus and Pi owners and current Sonnet rows.
  The registry and binding checks accept the current main owners. The dvd Codex row is absent from both files.
- DONE: Preserve every product and oracle byte from d02dd140.
  The five product and oracle blob identifiers match `d02dd1409`. The two binding files contain the authorized semantic union.
- DONE: Run focused parser, registry, owner, format, and diff checks only.
  Parser checks passed in 0.630s. Registry and owner checks passed in 0.340s. Format and diff checks passed.
- DONE: Do not repeat full, race, or live runs.
  No full, race, or live run occurred after the rebase. The prior durable results remain applicable to unchanged behavior bytes.
- DONE: Push the rebased exact head and update the Simplified-English report.
  Local and remote code heads match `47d39b624edbdb9590b9fbabf7c8c9dd27534989`. The seven-file candidate remains +131/-40 and net +91.

### Summary

The rebased candidate keeps current main ownership and removes only the dvd Codex binding. All product and oracle bytes remain unchanged.
The focused checks pass. The remote branch now identifies `47d39b624edbdb9590b9fbabf7c8c9dd27534989`.
