---
title: "Reject gate recording when the Briefing and current workflow stage disagree"
status: implementation
source: "se0 live CI recovery, 2026-07-28: TestLiveDefaultHeadlessStopsAtGate reproduced a validation Briefing durably bound as an implementation gate; captain separated the feature defect from CI restoration."
score: 0.85
sprint: durable-decisions
group: recorder
sprint-readiness: ready
issue:
id: q3vpb8hes1b3k3f1jps1kvpk
gates:
    version: 1
    current:
        gate: gate:q3vpb8hes1b3k3f1jps1kvpk:ideation
    records:
        - id: gate:q3vpb8hes1b3k3f1jps1kvpk:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:q3vpb8hes1b3k3f1jps1kvpk-backlog-1
              briefing:
                id: briefing:q3vpb8hes1b3k3f1jps1kvpk:backlog:attempt-1:revision-1
                digest: sha256:d84af4e055a57f211a97dbaa3faa7ec5c8dbc2dbcb85b0134c3096bcbc215652
                digest-domain: canonical-bytes
                request-digest: sha256:67af66f54068626fea34e776301eb17b8f5e59bee743606828c6f3713ca394af
                room-ref: ./gate-record-stage-coherence-guard/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:q3vpb8hes1b3k3f1jps1kvpk:backlog:1
                briefing: briefing:q3vpb8hes1b3k3f1jps1kvpk:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T15:39:37.060215Z"
                decision: approve
                reason: 'Sprint conn: task is a bounded, reproduced gate-semantic defect whose ideation can proceed independently of other critical lanes.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:q3vpb8hes1b3k3f1jps1kvpk:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:q3vpb8hes1b3k3f1jps1kvpk-ideation-1
              briefing:
                id: briefing:q3vpb8hes1b3k3f1jps1kvpk:ideation:attempt-1:revision-1
                digest: sha256:a82475a2afe17b637bac82476aa968ae78f38c8ff12b1ec3d1b441b339e9b585
                digest-domain: canonical-bytes
                request-digest: sha256:dd7e3e5fbd9c128b8f2bad933e63be4c533dca0e88a7e30aa98cc769d18979af
                room-ref: ./gate-record-stage-coherence-guard/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:q3vpb8hes1b3k3f1jps1kvpk:ideation:1
                briefing: briefing:q3vpb8hes1b3k3f1jps1kvpk:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T15:51:49.129297Z"
                decision: approve
                reason: 'Sprint conn: corrected design closes the noncanonical and cross-stage authority bypass with falsifiable byte-clean guards, bounded six-file scope, and no compatibility path.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
---

## Problem

`gate record` could bind a stage-qualified validation Briefing while the task was
still in `implementation`. At the reproduced tip, `RecordSemantic` held the entity
lock, but `recordBriefingLocked` had no workflow-stage authority check.
`initialIDs` recognized a structured Briefing ID only when its stage matched current
status; on mismatch it silently fell back to `gate:<entity>:<current-stage>`. The
result was an implementation attempt whose retained Briefing was qualified for
validation.

The retained `se0` Sonnet journey demonstrated the supported trigger: a task resumed
at `implementation` while a later validation Briefing remained on disk. The First
Officer read the current-stage failure, selected the retained validation evidence
anyway, and `gate record` accepted it. The final durable state remained
`status: implementation`, `gate:...:implementation`, with a validation-qualified
Briefing. This is a genuine recorder authority defect, not a test-oracle defect.

The provider-neutral preparation work landed afterward at `dbaebc795`. `gate prepare`
now constructs and binds
`briefing:<entity>:<current-stage>:attempt-N:revision-1`; ordinary `gate record`
accepts only a chat decision or the already-bound room. That removes the original
arbitrary-bind entrance, but not the authority defect at the spending boundary:
an already malformed open attempt can still be closed. A throwaway package spike on
current `main` constructed an `implementation` record containing
`briefing:task:validation:attempt-1:revision-1`; this command succeeded and wrote an
approval:

```text
go test ./internal/gates \
  -run TestSpikeCurrentRecorderClosesCrossStageBriefing -count=1 -v
PASS
```

The spike file was removed. It proves the residual defect independently of the
quarantined live journey and identifies `currentStageAttempt`, shared by room and chat
closure, as the smallest current seam.

A cycle-2 throwaway changed only the retained ID to `briefing:legacy`;
`TestSpikeCurrentRecorderClosesUnqualifiedBriefing` also passed and wrote a hold
Resolution. That file was removed. Unparseable identity is therefore an exercised
authority bypass, not a hypothetical compatibility concern.

`TestLiveDefaultHeadlessStopsAtGate` remains temporarily TODO under this task. It is one
test definition executed by the Sonnet and Opus Claude CI matrix legs; no Codex, Pi,
or shared gate/rejection/keep-moving scenario is quarantined. Removing that TODO is CI
restoration after the deterministic product tests are green, not a substitute for
those tests and not permission to change the live oracle.

## Fail-closed contract

For ordinary `gate record --decision` and `gate record --room`, the binary derives the
authoritative stage from the locked entity's `status`, resolves that stage in the
workflow definition, selects that stage's logical gate and last attempt, then checks
the attempt's Briefing identity before reading a decision into durable state.

A Briefing ID is valid for ordinary close only when it matches the shipped canonical
v1 shape
`briefing:<entity-identity>:<stage>:attempt-<positive>:revision-<positive>`.
`<entity-identity>` may itself contain colons, so parsing anchors from the
`:attempt-N:revision-N` suffix and takes the immediately preceding component as the
stage. An unqualified or malformed ID exits nonzero with exactly:

```text
Briefing id briefing:legacy is not a canonical stage-qualified v1 identity
```

If a canonical ID's stage differs from current `status`, recording exits
nonzero with exactly:

```text
Briefing stage validation does not match current workflow stage implementation
```

If current status is undefined, not `gate: true`, or terminal, ordinary recording
exits nonzero. A defined non-gated or terminal stage uses the stable diagnostic:

```text
current workflow stage implementation is not an actionable gate:true stage
```

The checks run under the existing entity lock and before either close path constructs
a Resolution or invokes `writeDocument`. Every refusal leaves the complete entity
byte-identical and removes the lock. The CLI emits no success stdout and prefixes the
diagnostic with its existing `Error: ` convention.

There is no compatibility parser or migration path: the task adds no alternate
format, rewrites no retained record, and accepts no prototype field. Existing
malformed or unqualified state remains readable, but an ordinary close is a new
authority mutation and must prove canonical current-stage identity. Recovery uses the
supported current-stage prepare/re-gate route rather than spending the retained
malformed attempt.

`gate record --round` remains outside this check. Its explicit `STAGE/CYCLE` denotes
historical advisory evidence and may legitimately differ from current status.

## Legitimate re-gates

The guard must preserve all of these:

- A matching current-stage first attempt or successor attempt can close.
- Re-entry to a prior logical gate selects the record named by current `status`, even
  when `gates.current.gate` still names another closed stage; a matching Briefing for
  the re-entered stage can close without rewriting earlier records.
- After rejection, a later validation gate can close a new validation-qualified
  Briefing. Whether that Briefing is fresh relative to rework belongs to `zbc`, not
  this stage-equality check.

## Proposed approach and mechanism choices

1. Add one private validator/parser for the complete canonical Briefing shape and one
   `validateRecordStage` helper. Both room and chat paths call it through
   `currentStageAttempt` while the existing lock is held. This serves AC-1 and AC-2.
   The simpler CLI-only check is insufficient because package callers and one of the
   two semantic sources could bypass it; treating an unparsable ID as “no stage
   claim” is also insufficient because it recreates the authority bypass.
2. Extend the gates package's existing workflow-stage projection with `gate` and
   `terminal` booleans; do not import status internals or create a second workflow
   parser. This serves AC-2. Checking only for a matching gate record is insufficient
   because malformed durable state can contain a record at a non-gated stage.
3. Keep the model/storage schema, `gate prepare`, round recording, application,
   consume, and presentation/provider logic unchanged. This serves AC-3 by making the
   guard a close-time precondition rather than a new lifecycle.
4. Add package and command tests first, then remove the live TODO. This serves AC-1
   through AC-4. A prose-only precondition is insufficient: the recovered Sonnet run
   loaded and ignored that exact class of instruction.

## Dependencies and collisions

- `dbaebc795`/s4 is a landed prerequisite and zero-delta dependency. Its generated
  matching IDs supply the positive path; this task does not edit `gate prepare`.
- `0m6` owns withdrawn-attempt classification and replacement preparation. If it lands
  first, rebase the shared attempt-selection helper and retain its withdrawn refusal;
  do not add withdrawal fields, diagnostics, or successor rules here.
- `zbc` owns post-rework freshness. Stage equality cannot prove freshness: a stale
  validation Briefing still says `validation`. Keep both guards independently
  falsifiable and do not absorb correction-epoch or Briefing-selection behavior.
- `se0` owns the broader live-CI recovery. This task removes only its linked
  `TestLiveDefaultHeadlessStopsAtGate` TODO after deterministic proof; it does not
  green unrelated live failures, alter provider behavior, or change expected journey
  semantics.

## Acceptance criteria

**AC-1 (VALUE) - An ordinary gate decision closes only a canonical current-stage Briefing.**

Verified by: package and command fixtures begin at `implementation` with a
validation-qualified, unqualified, and malformed open attempt in turn; both chat and
room record forms exit nonzero with the exact mismatch or canonical-identity
diagnostic, no success stdout, byte-identical entity state, no Resolution, and no
lock residue. The test fails if malformed identity falls through as “no stage claim,”
the comparison is deleted or moved after write, or only one semantic source is
guarded.

**AC-2 - Ordinary gate recording requires an actionable current workflow stage.**

Verified by: a table exercises matching IDs at a non-gated stage and a terminal stage,
asserting the exact actionable-gate diagnostic and byte/lock cleanliness, while an
existing genuine gated-stage fixture still closes. The test fails if stage metadata
is ignored or terminal is treated as actionable.

**AC-3 - Legitimate re-gates remain expressible without storage or grammar changes.**

Verified by: existing prepare/record/round positives remain green and one focused
re-entry fixture closes the current-status record while preserving the formerly
selected record byte-for-byte. A schema snapshot and CLI help assertion remain
unchanged. The test fails if the guard uses `gates.current` as authority, rejects
canonical current-stage successors, touches rounds, or introduces another
command/field.

**AC-4 - Both Claude matrix legs actively prove default-headless stop-open behavior.**

Verified by: the linked TODO is removed only after AC-1 through AC-3 are green; focused
Sonnet and Opus `TestLiveDefaultHeadlessStopsAtGate` runs use the exact candidate tip,
then the affected Claude live lane passes. No expected final-state assertion changes.
The check fails if the skip remains, either model is omitted, or the journey closes,
consumes, advances, or dispatches past the open gate.

## Expected surface and observable semantics

Baseline estimate: 6 files, about 145 insertions and 8 deletions. Tolerance is at most
2 additional files or 60 additional inserted lines; exceeding either requires design
re-entry. Expected files:

- `internal/gates/operation.go` — parser, workflow flags, and shared close guard,
  approximately +40/-5.
- `internal/gates/gates_test.go` — package refusal/re-gate table, approximately +55.
- `internal/cli/gate_test.go` — exact output and byte/lock controls, approximately +35.
- `docs/specs/gate-resolution-frontmatter-contract.md` — normative close precondition,
  approximately +8/-1.
- `docs/site/reference/command-reference.md` — user-facing recovery sentence,
  approximately +4/-1.
- `internal/ensigncycle/live_gate_stop_test.go` — remove one TODO, -1.

Command grammar and stored formats do not change. Authority becomes stricter only for
ordinary gate closure. Runtime behavior gains three byte-clean exit-1 diagnostics:
stage mismatch, non-actionable stage, and noncanonical identity. Round recording,
valid closure, preparation, provider association, application, and reads retain
current behavior. No Subspace/provider semantics change.

## Documentation diff

Apply this normative addition after the recorder-close paragraph in
`docs/specs/gate-resolution-frontmatter-contract.md`:

```diff
+Before either ordinary close, the recorder resolves authoritative current status in
+the workflow taxonomy and requires a nonterminal `gate: true` stage. The bound
+Briefing must use the canonical v1 stage-qualified identity and name that same stage.
+Malformed identity, mismatch, or non-actionable stage fails before Resolution
+construction and leaves entity bytes unchanged.
```

Append this sentence to both `gate record` rows in
`docs/site/reference/command-reference.md`:

```diff
+The current workflow stage must be an actionable gate, and the bound Briefing must
+use the canonical v1 stage-qualified identity and name that stage; malformed or
+mismatched identity fails without mutation.
```

## Test plan

Add the smallest failing package and CLI cases before production changes. Reuse the
existing semantic-decision, room, cross-gate, and byte-clean helpers; do not create a
parallel gate harness.

1. Package tests: table chat/room mismatch, unqualified, malformed, non-gated,
   terminal, canonical matching, successor, and cross-gate re-entry. Cost: small,
   deterministic fixture expansion; no new golden or live harness.
2. CLI tests: assert exit 1, exact stderr, empty stdout, byte identity, absence of a
   Resolution, and removed `.gates.lock` for mismatch, unqualified/malformed identity,
   and non-actionable stage. Cost: one consolidated command table.
3. Regression suites: run focused `./internal/gates` and `./internal/cli`, then
   `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
4. Live restoration: remove the TODO and run focused Sonnet and Opus
   `TestLiveDefaultHeadlessStopsAtGate`, then the affected Claude live CI lane at the
   exact candidate SHA. This is required because the changed test file is live-loaded,
   but it is not evidence in place of the deterministic product guard.

No fixture or expectation may be changed merely to make `se0` green. A deterministic
test that still permits a cross-stage or noncanonical-identity Resolution is a product
failure; a deterministic green followed by unrelated live infrastructure red is CI
restoration work and routes back to `se0`.

## Stage Report: ideation

- DONE: Define the smallest fail-closed stage-coherence contract at gate-record time, including the legitimate re-gate cases that must remain expressible.
  AC-1, AC-2, and AC-3 bind both ordinary close sources to locked current status while preserving successors, cross-logical-gate re-entry, unqualified IDs, and advisory rounds; deleting either comparison makes the named negative fixtures accept a Resolution.
- DONE: Turn the recovered semantic work into falsifiable ACs and tests that distinguish a product defect from CI restoration or test expectation changes.
  AC-4 follows the throwaway current-main spike and deterministic byte/output/lock controls with unchanged-oracle Sonnet/Opus restoration, routing unrelated live red to se0.
- DONE: Declare expected files, LOC, and observable semantics, spike the riskiest mechanism first, and add no compatibility behavior for unreleased formats.
  The body declares 6 files, +125/-8 with tolerance, exact authority/runtime deltas, zero schema/grammar/provider change, and a removed throwaway spike; introducing migration, identity storage, or zbc/0m6 behavior breaches the baseline.

### Summary

Ideation narrows the recovered pre-s4 bind defect to the still-reproducible close-time
authority gap and defines one shared locked guard for chat and room decisions. It
records exact diagnostics, byte-clean proof, legitimate re-gates, current dependency
collisions, documentation wording, and a test-first path that keeps CI restoration
separate from the product fix.

## Stage Report: ideation (cycle 2)

- DONE: Define the smallest fail-closed stage-coherence contract at gate-record time, including the legitimate re-gate cases that must remain expressible.
  AC-1, AC-2, and AC-3 now require canonical v1 identity plus current-stage equality on both close sources, preserve only canonical current-stage successors/re-entry and rounds, and make unqualified or malformed closure a named refusal.
- DONE: Turn the recovered semantic work into falsifiable ACs and tests that distinguish a product defect from CI restoration or test expectation changes.
  AC-1 adds chat/room unqualified and malformed mutants with exact diagnostics and byte/Resolution/lock controls; AC-4 still restores the unchanged live oracle only after those deterministic guards pass.
- DONE: Declare expected files, LOC, and observable semantics, spike the riskiest mechanism first, and add no compatibility behavior for unreleased formats.
  The corrected baseline remains 6 files but rises to +145/-8 for the negative identity cases; reads remain tolerant, new closes fail, and no migration, alternate identity, round, 0m6, zbc, se0, or provider behavior enters scope.

### Summary

Cycle 2 closes the unqualified-ID bypass identified in material feedback: every
ordinary decision now requires the canonical v1 shape and a current-stage match.
Diagnostics, ACs, tests, estimate, documentation wording, and the appended report
move together while historical reads, advisory rounds, and legitimate canonical
re-gates remain intact.
