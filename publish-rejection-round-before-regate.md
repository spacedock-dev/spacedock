---
title: Publish the rejected validation round before correction and re-gating
status: implementation
source: "Replacement for archived rejected zbc; Runtime Live E2E rejection-flow evidence showed the FO claimed validation/1 was recorded without invoking gate record --round validation/1."
started: 2026-08-09T18:34:24Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
worktree: .worktrees/spacedock-ensign-publish-rejection-round-before-regate
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
        - id: gate:zhcb4bcz1qgcn7ajx2ctxpxk:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zhcb4bcz1qgcn7ajx2ctxpxk-ideation-1
              briefing:
                id: briefing:zhcb4bcz1qgcn7ajx2ctxpxk:ideation:attempt-1:revision-1
                digest: sha256:00429838e456b758e78fab8b2133e09bb8b80a909c4c243625003a297dbf695e
                request-digest: sha256:d95410918ea0207d251a276b57d3fb80d890834e29d2db77bdcf158afb9cc526
                room-ref: ./publish-rejection-round-before-regate/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zhcb4bcz1qgcn7ajx2ctxpxk:ideation:1
                briefing: briefing:zhcb4bcz1qgcn7ajx2ctxpxk:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T21:33:17.763816Z"
                decision: approve
                reason: Captain approved complete rejected-round publication before correction.
              application:
                target-stage: implementation
                state: consumed
---

Restore durable recorder publication for supported rejection-flow failures. Keep the complete rejected round visible before correction re-gating.

## Problem

The rejection flow must publish validation/1 after the first rejection. Operators
must see the complete rejected round before the corrected candidate enters
re-gating. The existing recorder supports this state. Live conduct can call it
at the wrong point.

Pi is the current supported recorder-publication case. Run `31016570689`, job
`92342373497`, and artifact `8935708302` show an exact recorder result with
`entries=2`. The worker log and Cycle line came later. No second recorder call
published the complete four-entry log. This evidence supports a typed semantic
`rejection-round-incomplete` at the durable assertion boundary.

Archived Sonnet and Opus results do not provide exact stable recorder evidence.
Sonnet passed one full run and had FAIL/PASS/PASS focused results. Opus had
FAIL/FAIL/FAIL focused results. Do not bind either target to this task until an
exact live artifact shows the recorder failure code and a repeat shows the same
code.

Codex does not reach the recorder boundary. Run `31032033236`, job
`92395174900`, and artifact `8941373026` stopped after one implementation report
and one REJECTED validation report. The separate task `Continue Codex rejection
after the first validation` (`dvddbpsf4tdt3yjw1yjyp14k`) owns correction,
validation/2, and the fresh final gate. It is the dependency seam for Codex
coverage. This task does not bind Codex.

The source currently names this task for Sonnet, Opus, Codex, and Pi. Recarve
that source. Admit Pi now. Admit Sonnet or Opus only after exact stable recorder
evidence. Use the Codex task for Codex continuation. Do not use archived `zbc`.

## Value

After the first rejection, operators can see one complete rejected round before
correction re-gating starts. The retained round contains reviewer and worker
entries. The next reviewer and ordinary gate then use the corrected candidate.

The end value is measurable for each admitted target. `TestLiveCommonRejectionFlow`
must leave two implementation reports, two validation reports, a retained
validation/1 round with four entries, and a fresh open validation gate. Its
command log must show one successful exact `gate record --round validation/1`
call. A missing call, an `entries=2` result, or one-cycle state must fail the
journey.

## Proposed approach

Use one recorder-reuse change for the admitted target. Keep the current neutral
`gate record --round` command and the current rejection routing. Move its
explicit invocation to the first point where the complete reviewer and worker
log exists.

1. Keep the target lookup, authorized package, budget probe, and worker reuse
   rules unchanged.
2. Wait for correction worker completion and its worker entries. Do not
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
source. The complete log cannot exist before the worker entries exist. This
boundary is the earliest valid publication point for a complete round. It keeps
the rejected round visible before correction re-gating.

### XFAIL-first dependency and target seam

Task `ts7gq0mr9s3chx2w4wppd1kt`, `Run known live behavior gaps as target-level XFAIL`,
must land before this product change. Run the real
`TestLiveCommonRejectionFlow` cell at the durable assertion boundary.

- Pi has current exact evidence for `rejection-round-incomplete`. Bind Pi as
  XFAIL after the strict runner is available.
- Sonnet has no stable recorder code. Run an exact artifact and a repeat before
  adding Sonnet to this task.
- Opus has no stable recorder code. Run an exact artifact and a repeat before
  adding Opus to this task.
- Codex has `rejection-flow-not-completed` evidence. The separate Codex task
  owns that code and the continuation journey. Do not bind it here.

The recorder codes have narrow meanings:

- `rejection-round-not-published` means that no successful exact recorder call
  was observed.
- `rejection-round-incomplete` means that the exact call succeeded but the
  retained round is not the complete four-entry round.

`rejection-flow-not-completed` does not belong to this recorder task. It means
that the journey stopped before the recorder boundary. Infrastructure,
authentication, timeout, fixture, and other semantic failures remain FAIL. Keep
TODO only when the journey cannot run. A target with unstable semantic outcomes
remains unbound until a new disposition exists. After the repair, XPASS must
remove each target XFAIL binding. PASS requires the unchanged journey.

## Out of scope

- Do not revive `zbc`'s `correction-round` field or Reference-binding schema.
- Do not add another recorder, flag, freshness protocol, or test harness.
- Do not weaken the durable rejection-flow assertion or accept a final-message claim.
- Do not change gate storage, round formats, worker reuse, or reviewer identity.
- Do not own Codex continuation. The separate Codex task owns its correction,
  validation/2, and final-gate value.
- Do not bind Sonnet or Opus without exact stable recorder evidence.

## Acceptance criteria

**AC-1 (VALUE) - An admitted target publishes one complete rejected validation round before correction re-gating.**

Verified by: the unchanged `TestLiveCommonRejectionFlow` command log contains
one successful exact `gate record --round validation/1` result, the retained
round validates with four entries, and the durable two-cycle assertion passes.
The no-invocation control, an `entries=2` result, and a final gate without the
retained round each fail.

**AC-2 (VALUE) - Recorder-failure coverage has honest target ownership.**

Verified by: Pi reports XFAIL only for `rejection-round-incomplete` after its
exact evidence. Sonnet and Opus remain unbound until an exact repeat confirms a
stable recorder code. Codex has no binding in this task. The separate task
`dvddbpsf4tdt3yjw1yjyp14k` owns `rejection-flow-not-completed` and the complete
Codex continuation.

**AC-3 - The change reuses the current rejection lifecycle and recorder.**

Verified by: the final diff changes the feedback-flow ordering, mirrored process
wording, and admitted target binding only. It adds no gate field, recorder
command, flag, freshness schema, parallel harness, or new authority source.

**AC-4 - Workflow-owned policy remains separate from recorder bytes.**

Verified by: the existing durable recorder test finds the Cycle line unchanged,
preserves lifecycle sentinels, and retains only the canonical Briefing and log
in the round room. A recorder edit that rewrites the Cycle line fails the test.

## Expected surface and semantic budget

Baseline: **4 existing files, about 12 gross insertions and 10 deletions, for a
net change of about 2 lines. Tolerance is ±1 file and ±12 net lines.**

- `skills/feedback-rejection-flow/SKILL.md`: about 7 insertions and 5
  deletions. Move the existing command to the complete-log checkpoint and make
  its success boundary explicit.
- `docs/dev/README.md`: about 2 insertions and 2 deletions. State that the
  recorder runs before reviewer re-run and next-gate preparation.
- `skills/commission/references/templates/development.md`: about 2 insertions
  and 2 deletions. Keep new development workflows aligned with the process doc.
- `internal/ensigncycle/shared_live_runner_test.go`: replace the current
  rejection-flow target binding with the Pi XFAIL binding after `ts7g` lands.
  Do not add Codex, Sonnet, or Opus bindings without their evidence gate.

The estimate excludes the target-level XFAIL runner and the separate Codex task. The
XFAIL task owns its result type, reconciliation, and metrics changes. The Codex
task owns continuation after the recorder boundary. This task owns only the
recorder publication order and its admitted target binding.

Declared semantic changes:

- **Command grammar:** none. Reuse the existing `gate record --round`
  invocation without a new flag or output requirement.
- **Stored formats:** none. Keep the existing review-round pointer and
  canonical two-file room.
- **Authority:** none. Keep the existing First Officer and neutral recorder
  writers and authority.
- **Runtime behavior:** publish one complete round after worker entries and
  before reviewer re-run and final gate preparation.
- **Coverage:** classify Pi as the current recorder-publication lane. Admit
  Sonnet and Opus only after exact stable evidence. Keep Codex in its own task.
- **Documentation:** state the same observable order and ownership seam.

## Test plan

First, keep the recorder spike as the offline gate. The exact command was:

```bash
go test ./internal/ensigncycle -run 'TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl|TestRejectionFlowRoundInvocationExtractors' -count=1 -v
```

It passed in 1.592 seconds. The durable test passed the four-entry recorder,
preserved status and lifecycle sentinels, and failed its inverted no-invocation
control. The extractor test passed wrong entity, round, suffix, and file
controls, and rejected missing or failed results. This proves the existing
recorder path before repair selection. No binary extension is indicated by this
spike.

After `ts7gq0mr9s3chx2w4wppd1kt` lands, build the current binary and run the Pi
focused cell:

```bash
SPACEDOCK_LIVE_RUNTIME=pi \
  go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle -v
```

Run the Sonnet and Opus cells only as evidence gates for later admission:

```bash
SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet \
  go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle -v
SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=claude-opus-4-8 \
  go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle -v
```

Record each target, candidate SHA, duration, run, job, artifact, result, and
semantic code. Pi must reproduce `rejection-round-incomplete`. Sonnet or Opus
can join only when the exact recorder command, durable code, and repeat match.
Do not run Codex as part of this task. The separate Codex task owns its live
continuation and `rejection-flow-not-completed` seam.

Implement the one skill ordering change, the two mirrored process wording
changes, and the admitted Pi binding. Re-run the focused offline recorder tests.
Then run the Pi XFAIL cell again. An XPASS removes the Pi binding. A PASS removes
the TODO or XFAIL state. A different FAIL stops the repair and requires a new
disposition.

Finally run the unchanged full and race suites, the live registry reconciliation,
and `gofmt -w ./cmd ./internal`. The final proof must include the complete
round, two-cycle durable state, fresh final gate, one successful recorder call,
and honest target coverage.

## Proposed documentation diff

```diff
--- a/skills/feedback-rejection-flow/SKILL.md
+++ b/skills/feedback-rejection-flow/SKILL.md
@@
-3. If the workflow declares a `### Feedback Cycles` correction-round projection, the First Officer appends its authorized line directly.
+3. Wait for complete correction worker entries. Then append the authorized `### Feedback Cycles` line.
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
--- a/internal/ensigncycle/shared_live_runner_test.go
+++ b/internal/ensigncycle/shared_live_runner_test.go
@@
-liveJourney(... []liveJourneyTODO{... Sonnet ..., ... Opus ..., ... Codex ..., ... Pi ...})
+liveJourney(... []liveJourneyGap{liveXFail("pi", "zhcb4bcz1qgcn7ajx2ctxpxk")})
```

## Stage Report: ideation

- DONE: Exercise the rejection-round recorder path before selecting the repair.
  `go test ./internal/ensigncycle -run 'TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl|TestRejectionFlowRoundInvocationExtractors' -count=1 -v` passed in 1.592 seconds. It proved the four-entry recorder, lifecycle preservation, exact invocation extraction, and no-invocation control.
- DONE: Fold staff finding M2 into a stable recorder-publication scope.
  Pi is the current admitted lane with run `31016570689`, job `92342373497`, artifact `8935708302`, and `entries=2` evidence for `rejection-round-incomplete`. Sonnet and Opus require an exact stable recorder result and a matching repeat. Codex does not reach the recorder boundary. Task `Continue Codex rejection after the first validation` (`dvddbpsf4tdt3yjw1yjyp14k`) owns its continuation and `rejection-flow-not-completed` seam.
- DONE: Define one recorder-reuse change and give gross and net line estimates.
  The plan moves the existing recorder call to the complete-log checkpoint before reviewer re-run and next-gate preparation. It estimates 4 existing files, about 12 gross insertions and 10 deletions, for a net change of about 2 lines. It records the archived Sonnet, Opus, and Codex evidence without binding those outcomes to this task.

### Summary

Ideation keeps the end-user value: the complete rejected round is visible before correction re-gating. It limits this task to stable recorder-publication failures, admits Pi from current evidence, gates Sonnet and Opus on exact repeatable evidence, and assigns Codex continuation to its separate owner.

## Stage Report: implementation

- DONE: Implement the approved complete rejected-round publication before reviewer rerun and next-gate preparation.
  Commit `267de20d6` changes `skills/feedback-rejection-flow/SKILL.md:17-22` and `skills/commission/references/templates/development.md:118`.
- DONE: Admit no new source binding now; preserve the existing recorder, round format, authority, and task-specific ownership.
  The change reuses `gate record --round`. It adds no field, command, flag, writer, or target binding.
- DONE: Add task-local focused proof for one complete four-entry validation/1 publication and the no-invocation control.
  `internal/contractlint/feedback_rejection_publication_smoke_test.go:9-45` fails if publication moves after reviewer or gate work.
  `TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl` fails without the four-entry round or observed invocation.
- DONE: Do not edit shared XFAIL bindings, registry reconciliation, sprint package files, or shared process documentation before ts lands.
  This checkpoint changes no shared binding, registry, sprint package, or `docs/dev/README.md` file.
- DONE: Commit the product behavior checkpoint and report exact files, lines, tests, and deferred post-ts live proof.
  Commit `267de20d6` contains three files. The Pi XFAIL binding and live proof remain deferred until `ts7gq0mr9s3chx2w4wppd1kt` lands.

### Summary

The rejection flow now publishes the complete round before reviewer rerun or gate preparation. Recorder failure or an incomplete log now holds the flow.

The focused contract and recorder tests pass. `go test ./...` and `go test ./... -race` pass after the required `gofmt` run.

## Stage Report: implementation (cycle 2)

- DONE: Rebase onto `origin/main` at or after `a8688cabf`, preserving the product repair.
  Rebased commit `1547a7bb5` contains the complete-round ordering, template update, and focused smoke test.
- DONE: Run the exact Claude Sonnet target and remove its `zh` XFAIL after proof.
  The bound run reported XPASS. The unbound rerun passed `TestLiveCommonRejectionFlow` in 513.65 seconds.
- DONE: Run the exact Claude Opus target and remove its `zh` XFAIL after proof.
  The bound run reported XPASS. The unbound rerun passed `TestLiveCommonRejectionFlow` in 635.20 seconds.
- FAILED: Run the exact Pi target and remove its `zh` XFAIL after proof.
  OpenRouter returned HTTP 402 before the journey started. The account limit was 22226 tokens, not the requested 128000 tokens.
- DONE: Preserve XFAIL ownership and update the exact reconciliation row.
  Commit `ea11da278` removes only the Sonnet and Opus `zh` rows at `internal/ensigncycle/shared_live_runner_test.go:125`.
  `internal/contractlint/live_registry_reconciliation_test.go:56` keeps the `dvd` Codex row and the blocked `zh` Pi row.
- DONE: Run the required checks and push the branch.
  The focused recorder tests, registry reconciliation, `go test ./...`, and `go test ./... -race` pass after `gofmt`.

### Summary

Sonnet and Opus now pass the repaired rejection flow without XFAIL bindings. Pi remains bound because an external credit limit blocked its exact run.

The branch is rebased on `a8688cabf` and pushed. No `dvd` or other task binding changed.
