---
title: Publish the rejected validation round before correction and re-gating
status: ideation
source: "Replacement for archived rejected zbc; Runtime Live E2E rejection-flow evidence showed the FO claimed validation/1 was recorded without invoking gate record --round validation/1."
started: 2026-08-09T18:34:24Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
worktree:
issue:
pr:
mod-block:
id: zhcb4bcz1qgcn7ajx2ctxpxk
gates:
    version: 1
    records:
        - id: gate:zhcb4bcz1qgcn7ajx2ctxpxk:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zhcb4bcz1qgcn7ajx2ctxpxk-backlog-1
              briefing:
                id: briefing:zhcb4bcz1qgcn7ajx2ctxpxk:backlog:attempt-1:revision-1
                digest: sha256:ca2ad92ff6a096a252a9d8afa04ceac8c0f42479830ee44dfa2410561c419c8d
                request-digest: sha256:aa3328285343a91f97a61f069c02112a1146eb4bdb5d7f472a65dff4eb020fbc
                room-ref: ./publish-rejection-round-before-regate/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zhcb4bcz1qgcn7ajx2ctxpxk:backlog:1
                briefing: briefing:zhcb4bcz1qgcn7ajx2ctxpxk:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:12.199234Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; the seed names the recorder reuse boundary and exact live proof.
              application:
                target-stage: ideation
                state: consumed
---

Restore the common rejection journey on Sonnet, Opus, Codex, and Pi by making the First Officer publish the rejected validation round before correction and re-gating.

## Problem

The `rejection-flow` journey needs durable evidence of validation/1 before the
First Officer re-runs the reviewer and prepares the next gate. The existing
recorder works, but live conduct does not always call it at that boundary.

The archived validation record shows the defect. At `f2bd91044`, the Codex
focused run failed at 390.47 seconds. The Sonnet full run passed once. The Opus
full run failed at 610.16 seconds. Repeated focused runs were FAIL/PASS/PASS
for Sonnet and FAIL/FAIL/FAIL for Opus. The report names the missing
`gate record --round validation/1` call as the common harm.

Pi gives a more precise failure. Run `31016570689`, job `92342373497`, and
artifact `8935708302` show one successful recorder result with `entries=2`.
The worker log and Cycle line were added later. No second recorder call
published the complete four-entry log.

Codex gives a lifecycle failure. Run `31032033236`, job `92395174900`, and
artifact `8941373026` reached one implementation report and one REJECTED
validation report. The First Officer prepared an ordinary gate and stopped.
It did not route rework, publish the round, or run validation/2.

The source keeps four target-scoped TODOs for Sonnet, Opus, Codex, and Pi. The
TODOs name this task, not archived `zbc`. Strict XFAIL evidence must classify
each target before a product change.

## Value

After a rejection, operators can inspect one complete correction history. The
retained round contains the reviewer and worker entries. The next reviewer and
ordinary gate then use the corrected candidate.

The end value is measurable. `TestLiveCommonRejectionFlow` must leave two
implementation reports, two validation reports, a retained validation/1 round
with four entries, and a fresh open validation gate. Its command log must show
a successful `gate record --round validation/1` call. A missing call, an
`entries=2` round, or one-cycle state must fail the journey.

## Proposed approach

Use one recorder-reuse change. Keep the current neutral `gate record --round`
command and the current rejection routing. Move its explicit invocation to the
first point where the complete reviewer and worker log exists.

1. Keep the target lookup, authorized package, budget probe, and worker reuse
   rules unchanged.
2. Wait for the correction worker's completion and its worker entries. Do not
   append the workflow-owned Cycle line before those entries exist.
3. Append the authorized Cycle line. Then invoke the existing command once:

   ```text
   ${SPACEDOCK_BIN:-spacedock} gate record --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl
   ```

4. Require a successful recorder result with the complete round summary before
   re-running the reviewer or preparing the next gate. If the command fails,
   produces no result, or retains an incomplete log, hold the flow. Do not claim
   that the round was recorded.
5. Re-run the reviewer and re-enter the normal gate flow after publication.

This reuses the existing recorder, immutable room, pointer, command log, and
durable assertion. It adds no recorder, field, freshness rule, or authority
source. The complete log cannot be published before the worker entries exist;
the proposed boundary is therefore the earliest valid publication point and
still precedes reviewer re-run and final gate preparation.

### XFAIL-first dependency

Task `ts7gq0mr9s3chx2w4wppd1kt` must land before this product change. Run the
real `TestLiveCommonRejectionFlow` cell for all four targets. Bind an XFAIL only
when the journey runs and returns one stable semantic code at the durable
assertion boundary.

- `rejection-round-not-published` means that no successful exact recorder call
  was observed.
- `rejection-round-incomplete` means that the exact call succeeded but the
  retained round is not the complete four-entry round.
- `rejection-flow-not-completed` means that the run stopped before the
  rejection flow reached its recorder boundary. Treat this as a different
  failure, not as recorder XFAIL evidence.

Infrastructure, authentication, timeout, fixture, and other semantic failures
remain FAIL. A target with unstable semantic outcomes remains FAIL until a new
disposition exists. Keep TODO only when the journey cannot run. After the repair,
XPASS must remove that target's XFAIL binding; the target has passing evidence
only when the unchanged journey passes.

## Out of scope

- Do not revive `zbc`'s `correction-round` field or Reference-binding schema.
- Do not add another recorder, flag, freshness protocol, or test harness.
- Do not weaken the durable rejection-flow assertion or accept a final-message claim.
- Do not change gate storage, round formats, worker reuse, or reviewer identity.
- Do not create a separate Opus task. Opus is the pre-release lane.

## Acceptance criteria

**AC-1 (VALUE) - A supported First Officer publishes one complete rejected validation round before reviewer re-run and final gate preparation.**

Verified by: the unchanged `TestLiveCommonRejectionFlow` command log contains
one successful exact `gate record --round validation/1` result, the retained
round validates with four entries, and the durable two-cycle assertion passes.
The no-invocation control, an `entries=2` result, and a final gate without the
retained round each fail.

**AC-2 (VALUE) - Each executable target has honest coverage state.**

Verified by: strict XFAIL runs cover Claude Sonnet, Claude Opus, Codex, and Pi.
Each target reports XFAIL only for its bound stable semantic code, XPASS when
the repair removes that code, or PASS after the unchanged journey passes. A
different failure remains FAIL. No target entry names archived `zbc`.

**AC-3 - The change reuses the current rejection lifecycle and recorder.**

Verified by: the final diff changes the feedback-flow ordering and its mirrored
process wording only. It adds no gate field, recorder command, flag, freshness
schema, parallel harness, or new authority source.

**AC-4 - Workflow-owned policy remains separate from recorder bytes.**

Verified by: the existing durable recorder test finds the Cycle line unchanged,
preserves lifecycle sentinels, and retains only the canonical Briefing and log
in the round room. A recorder edit that rewrites the Cycle line fails the test.

## Expected surface and semantic budget

Baseline: **4 existing files, about 11 gross insertions and 10 deletions,
for a net change of about 1 line; tolerance is ±1 file and ±12 net lines.**

- `skills/feedback-rejection-flow/SKILL.md`: about 7 insertions and 5
  deletions. Move the existing command to the complete-log checkpoint and make
  its success boundary explicit.
- `docs/dev/README.md`: about 2 insertions and 2 deletions. State that the
  recorder runs before reviewer re-run and next-gate preparation.
- `skills/commission/references/templates/development.md`: about 2 insertions
  and 2 deletions. Keep new development workflows aligned with the process doc.
- `internal/ensigncycle/shared_live_runner_test.go`: transient XFAIL bindings
  are owned by `ts7g`; remove each repaired binding on XPASS. Count at most 4
  binding-line replacements in gross work and 4 removals in the final diff.

The estimate excludes the strict XFAIL runner itself. That dependency owns its
result type, reconciliation, and metrics changes. This task only supplies the
four target bindings and removes them after XPASS.

Declared semantic changes:

- **Command grammar:** none. The existing `gate record --round` invocation is
  reused without a new flag or output requirement.
- **Stored formats:** none. The existing review-round pointer and canonical
  two-file room remain unchanged.
- **Authority:** none. The existing First Officer and neutral recorder retain
  their writers and authority.
- **Runtime behavior:** the feedback flow must publish one complete round after
  worker entries and before reviewer re-run and final gate preparation.
- **Documentation:** process wording states the same observable order.

## Test plan

First, keep the recorder spike as the offline gate. The exact command was:

```bash
go test ./internal/ensigncycle -run 'TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl|TestRejectionFlowRoundInvocationExtractors' -count=1 -v
```

It passed in 1.592 seconds. The durable test passed the four-entry recorder,
preserved status and lifecycle sentinels, and failed its inverted no-invocation
control. The extractor test passed wrong entity, round, suffix, and file
controls, and rejected missing or failed results. This proves the existing
recorder path before any repair selection. No binary extension is indicated by
this spike.

After `ts7gq0mr9s3chx2w4wppd1kt` lands, build the current binary and run the
focused cell with each supported target:

```bash
SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet \
  go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=claude-opus-4-8 \
  go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=codex \
  go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=pi \
  go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle -v
```

Record the target, candidate SHA, duration, result, and stable semantic code
from each real artifact. Keep TODO for an unavailable target. Convert only a
stable semantic failure to XFAIL. Do not classify a skipped, timed-out, or
unrelated failure as XFAIL.

Implement the one skill ordering change and the two mirrored process wording
changes. Re-run the focused offline recorder tests. Then run each XFAIL cell
again. An XPASS removes its binding. A PASS removes its TODO. A different FAIL
stops the repair and requires a new disposition.

Finally run the unchanged full and race suites, the live registry reconciliation,
and `gofmt -w ./cmd ./internal`. The final proof must include the complete
round, two-cycle durable state, fresh final gate, one successful recorder call,
and clean target coverage state.

## Proposed documentation diff

```diff
--- a/skills/feedback-rejection-flow/SKILL.md
+++ b/skills/feedback-rejection-flow/SKILL.md
@@
-3. If the workflow declares a `### Feedback Cycles` correction-round projection, the First Officer appends its authorized line directly.
+3. Wait for the correction worker's complete entries. Then append the authorized `### Feedback Cycles` line.
@@
-7. Re-enter the normal gate flow with the updated result. After the reviewer log and any workflow-owned Cycle line are complete, invoke the neutral `${SPACEDOCK_BIN:-spacedock} gate record --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl`.
+6. Invoke the existing `${SPACEDOCK_BIN:-spacedock} gate record --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl` once. Require its complete round result before reviewer re-run or gate preparation.
+7. Re-run the reviewer and re-enter the normal gate flow after publication.
--- a/docs/dev/README.md
+++ b/docs/dev/README.md
@@
-After reviewer and worker entries and FO consultation, the First Officer appends the workflow-defined Cycle line directly from the authorized package, then invokes `${SPACEDOCK_BIN:-spacedock} gate record --round` with the canonical Briefing/log.
+After reviewer and worker entries and FO consultation, the First Officer appends the workflow-defined Cycle line directly from the authorized package, then invokes `${SPACEDOCK_BIN:-spacedock} gate record --round` with the canonical Briefing/log before reviewer re-run or next-gate preparation.
--- a/skills/commission/references/templates/development.md
+++ b/skills/commission/references/templates/development.md
@@
-After reviewer and worker entries and FO consultation, the First Officer appends the workflow-defined Cycle line directly from the authorized package, then invokes `${SPACEDOCK_BIN:-spacedock} gate record --round` with the canonical Briefing/log.
+After reviewer and worker entries and FO consultation, the First Officer appends the workflow-defined Cycle line directly from the authorized package, then invokes `${SPACEDOCK_BIN:-spacedock} gate record --round` with the canonical Briefing/log before reviewer re-run or next-gate preparation.
```

## Stage Report: ideation

- DONE: Exercise the rejection-round recorder path before selecting the repair.
  `go test ./internal/ensigncycle -run 'TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl|TestRejectionFlowRoundInvocationExtractors' -count=1 -v` passed in 1.592 seconds. It proved the four-entry recorder, lifecycle preservation, exact invocation extraction, and no-invocation control.
- DONE: Define one recorder-reuse change with XFAIL-first dependencies for all executable targets.
  The plan moves the existing `gate record --round` call to the complete-log checkpoint before reviewer re-run and next-gate preparation. It requires `ts7gq0mr9s3chx2w4wppd1kt` first and defines target-scoped stable codes for Sonnet, Opus, Codex, and Pi.
- DONE: Give gross and net line estimates with exact runtime evidence.
  The plan estimates 4 existing files, about 11 gross insertions and 10 deletions, for a net change of about 1 line. It records archived Codex, Sonnet, Opus, and Pi run evidence with run, job, artifact, duration, and failure shape.

### Summary

Ideation selected one small lifecycle-order change that reuses the existing recorder and keeps gate and round formats unchanged. The recorder spike passed before repair selection, and strict XFAIL classification now gates every supported runtime target before product work.
