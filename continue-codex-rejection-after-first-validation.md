---
title: Continue Codex rejection after the first validation
status: ideation
source: "Staff review M2 for test-behavior-completeness, 2026-08-09"
started: 2026-08-09T20:36:16Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree:
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

After the known `ts` strict-XFAIL machinery and the preceding `zh` recorder
work land, make one non-product baseline commit. This commit contains only
the Codex host extractor, the strict final-gate oracle, its negative controls,
and the target binding. It contains no skill text, Codex runtime instruction,
or user-visible documentation change.

The baseline commit replaces only the Codex `rejection-flow` TODO with the
target-scoped expected semantic code:

```text
target: codex
owner: dvddbpsf4tdt3yjw1yjyp14k
code: rejection-flow-not-completed
```

Run the complete Codex rejection journey from that baseline commit. It must
not skip. The Codex host extractor must select
`codexRecordedRejectionRound`, and the strict final-gate oracle must reach the
durable missing-gate state. The strict grade must report exactly one semantic
code: `rejection-flow-not-completed`. A parser error, auth error, launch
error, timeout, malformed state, or any extra semantic error is an ordinary
failure. The baseline evidence is not valid until the live run reaches this
sole-code result.

Do not change skills, Codex runtime instructions, or user-visible behavior
until the baseline has produced the sole expected XFAIL code. This order keeps
the baseline proof independent from the repair.

With the binding still present, the repaired exact candidate must produce an
XPASS and fail the lane. Remove the binding only after that XPASS. Run the same
candidate again and require a PASS with no semantic code. The candidate keeps
the Sonnet, Opus, and Pi ownership rows unchanged.

## Acceptance criteria

**AC-1 — The Codex baseline executes the rejection-flow journey and records exactly `rejection-flow-not-completed` as its strict-XFAIL semantic outcome.**

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
baseline produces the sole expected XFAIL code.

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
   journey and the sole semantic code `rejection-flow-not-completed`. A parser
   failure or any additional code stops the repair. Save the baseline JSONL,
   final message, entity, reports, round room, and gate evidence.
3. Only after the sole-code XFAIL, add the feedback skill and Codex runtime
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
- DONE: Require the complete Codex journey to produce the sole expected
  XFAIL. The baseline must reach the durable missing-gate state through
  validation/2 and report only `rejection-flow-not-completed`. A parser error,
  infrastructure error, or extra semantic code fails the baseline.
- DONE: Preserve the complete correction value and the exact-candidate proof.
  After the sole-code XFAIL, the repair produces XPASS with the binding and
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
