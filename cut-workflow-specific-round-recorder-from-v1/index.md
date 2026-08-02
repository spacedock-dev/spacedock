---
title: Cut workflow-specific advisory round recording from v1
status: ideation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: generic gate code embeds the development workflow's Material/fixed/declined taxonomy, LOC estimate grammar, ensign role, and Feedback Cycles projection."
started: 2026-08-01T14:00:47Z
completed:
verdict:
score: "0.95"
worktree:
issue:
pr:
sprint: durable-decisions
id: wjkhq0sktbbe3txx6jhnvcv2
gates:
    version: 1
    current:
        gate: gate:wjkhq0sktbbe3txx6jhnvcv2:ideation
    records:
        - id: gate:wjkhq0sktbbe3txx6jhnvcv2:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wjkhq0sktbbe3txx6jhnvcv2-backlog-1
              briefing:
                id: briefing:wjkhq0sktbbe3txx6jhnvcv2:backlog:attempt-1:revision-1
                digest: sha256:45e247730eceb6bb9fb0e68f806630bab3cdd2df47830efc5f4f9e40792351f8
                digest-domain: canonical-bytes
                request-digest: sha256:13584896ffe679ba3def35d84099ea4402ac80d53f7836c89be8fa535efc32af
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wjkhq0sktbbe3txx6jhnvcv2:backlog:1
                briefing: briefing:wjkhq0sktbbe3txx6jhnvcv2:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:50:02.905125Z"
                decision: approve
                reason: Captain approved dispatching this durable-decisions ideation lane in parallel with hq, nth, and jc.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:wjkhq0sktbbe3txx6jhnvcv2:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:wjkhq0sktbbe3txx6jhnvcv2-ideation-1
              briefing:
                id: briefing:wjkhq0sktbbe3txx6jhnvcv2:ideation:attempt-1:revision-1
                digest: sha256:36d1434d6fe326e0b9a94e81b131ed047ee1e5986a4d5b2d1bcbe69e49f82f33
                digest-domain: canonical-bytes
                request-digest: sha256:bb1d202e76847364f51a17ff726b6c7e6f71a6820298d600aae3c67e658c01fa
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wjkhq0sktbbe3txx6jhnvcv2:ideation:1
                briefing: briefing:wjkhq0sktbbe3txx6jhnvcv2:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-01T16:07:10.322447Z"
                decision: revise
                reason: 'Captain directed send-back for Science Officer REVISE: AC-2 and AC-3 lack concrete evidence, and the design must preserve the zbc correction-round producer; route for ideation correction.'
              application:
                action: feedback
                target-stage: ideation
                state: superseded
            - id: gate-attempt:wjkhq0sktbbe3txx6jhnvcv2-ideation-2
              briefing:
                id: briefing:wjkhq0sktbbe3txx6jhnvcv2:ideation:attempt-2:revision-1
                digest: sha256:66eb240047269d7861200ce4c5c0667fb3139a8b7c621c34edc49f63f59d3c67
                digest-domain: canonical-bytes
                request-digest: sha256:915c92a2caad3262388d2465d2564073eb7bea15d5c59b5801aa3a72d53d3b68
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:wjkhq0sktbbe3txx6jhnvcv2:ideation:2
                briefing: briefing:wjkhq0sktbbe3txx6jhnvcv2:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-02T09:57:43.891639Z"
                decision: hold
                reason: Science approves the bounded WJ design, but NTH cycle 2 evidence and the JC reset path are not yet durable. Hold WJ at ideation to avoid independent implementation on a stale shared seam.
              application:
                action: none
                state: not-applicable
            - id: gate-attempt:wjkhq0sktbbe3txx6jhnvcv2-ideation-3
              briefing:
                id: briefing:wjkhq0sktbbe3txx6jhnvcv2:ideation:attempt-3:revision-1
                digest: sha256:182bbf2881ea281e0b37e5df2a8a78cf5deb60f57691dc784b93a893199a5c9c
                digest-domain: canonical-bytes
                request-digest: sha256:ed5062b5037618d28123ba45285282d81120a3f70f5e6a3aa469e1fca1aa54b7
                room-ref: ./review/ideation/briefing-3
---

Cut the development-specific classification and Cycle projection from the public round surface while retaining its workflow-neutral correction-record producer for stable v1. Keeping policy would freeze this repository's review taxonomy into the generic durable-decision contract; deleting the structural producer would break zbc's post-rework freshness boundary.

## Problem statement

`internal/gates` currently recognizes `material` and `correct-but-disproportionate`, decides fixed versus declined, requires `actor:ensign`, parses the development workflow's files/LOC/estimate/percentage/AC cycle sentence, and mutates a `### Feedback Cycles` body section. Those are workflow judgments and projections, not binding gate primitives. The command therefore gives a generic package authority it must not have, and a prose-only rename leaves that jurisdiction error intact.

The bounded v1 problem is to remove the development workflow policy from this operation while preserving both every binding gate operation and the workflow-neutral correction-round producer consumed by `zbc` (`bind-post-rework-briefing-at-rejection-regate`). Historical entity files may contain `### Feedback Cycles`; the binary stops creating or validating that projection. It continues to retain canonical round bytes in an immutable room and advance the structural `review-round` pointer, without deciding what the findings mean.

## Proposed approach

Keep `gate record --round STAGE/CYCLE --briefing ... --log ...` and `gate validate --round STAGE/CYCLE` as a workflow-neutral producer/validator. They continue to verify canonical Briefing and review-log structure, atomically publish the two immutable room files, advance the `review-round` pointer, replay byte-cleanly, and leave binding `gates:` state unchanged. Delete only the development policy: `--feedback-cycle`, its fixed Cycle-line grammar and body splicer, the Material/fixed/declined/`actor:ensign` classifier, the `triage=` summary field, and their fixtures/tests. Add command tests that pin the reduced help, reject `--feedback-cycle` as an unknown flag with byte-identical state, accept structurally valid logs whose labels/actors do not follow development policy, and exercise retained record/validate plus every binding verb.

Move ownership of correction-round interpretation and projection back to the active workflow text: reviewers report findings; workers investigate and propose materiality, ownership, and disposition; the First Officer authorizes candidate mutation and writes the workflow's `### Feedback Cycles` entry directly. The neutral producer then retains the already-completed Briefing/log bytes and publishes the pointer that `zbc` consumes. The binary neither parses nor certifies workflow classifications or the Cycle line. This serves AC-3. The simplest alternative—removing only the help spelling while retaining classifier/projection internals—is insufficient because the same generic code would still enforce the ensign taxonomy and project workflow prose (AC-2). Removing the producer is also insufficient: it deletes the only on-disk correction identity that `zbc` uses to reject stale post-rework Briefings.

Do not build a second recorder or compatibility alias. The existing structural producer is the retained neutral boundary; all classifications and prose projections are owned by the active workflow.

## Streamlined scenario

During implementation review, a reviewer reports a finding. The worker records evidence and proposes materiality, task ownership, and disposition in the stage report; the First Officer separately authorizes any fix and directly appends the workflow-defined Cycle line. The First Officer invokes the neutral round producer to retain the canonical Briefing/log and advance `review-round`; the command does not receive or interpret the Cycle line. A later `gate prepare` may use that structural pointer as `zbc` specifies. Binding stage approval still uses `gate prepare`, binding `gate record`, `gate validate`, `gate eligibility`, and `gate consume`.

## Semantic scope

- Command grammar changes: remove only record flag `--feedback-cycle` and the `triage=` success field. Retain record flags `--round`, `--briefing`, and `--log`, plus validate flag `--round`; their structural identity fields remain stable.
- Stored formats: retain `review-round` frontmatter and `review/<stage>/round-<cycle>/` as the `zbc` producer boundary. Stop binary-authored `### Feedback Cycles` lines; do not migrate, reject, or rewrite historical prose. The binding `gates:` graph and retained gate-room formats are byte-compatible.
- Authority changes: finding classification, disposition, correction ownership, cycle projection, and escalation remain workflow/First-Officer responsibilities; the generic binary has none.
- Runtime behavior: neutral round record/validate remains callable and no longer rejects a structurally valid log for missing development-specific reviewer/ensign disposition patterns. It never mutates the body. Binding prepare, record (room and chat), validate, withdraw, eligibility, consume, and merge behavior are unchanged.

## Shared route seams

This task owns no stage-route policy. It preserves the workflow-neutral seam shared with `nth` and `jc`: `jc`'s status-derived unique-stage lookup selects the gate attempt before the reducer; the reducer derives the next route from the workflow stage graph; approval alone may produce the existing neutral Application/Eligibility projection (`action: advance`, with the existing pending/consumed state); and revise/hold remain Resolution-only while workflow feedback routing stays outside Application. `nth` may remove policy-specific application fields without changing that neutral reducer. Wj must not add a round-derived route, application field, feedback target, or status transition.

The separate `zbc` seam is structural rather than routing policy: wj retains `gate record --round`, `gate validate --round`, the `review-round` identity, and the existing immutable two-file room exactly so `zbc` can bind a later post-rework Briefing to the newest correction episode. Round publication continues to leave binding `gates:` and the neutral stage-route reducer untouched.

## Expected surface and tolerance

Expected implementation surface: exactly 16 files, approximately 300-520 deleted lines and no more than 220 inserted lines (replacement producer tests and workflow wording):

1. `internal/cli/cli.go`
2. `internal/cli/gate_test.go`
3. `internal/gates/round.go`
4. `internal/gates/round_test.go`
5. `internal/gates/review.go`
6. `internal/gates/model.go`
7. `internal/gates/io.go`
8. `internal/gates/operation.go`
9. `internal/ensigncycle/shared_round_recording_test.go`
10. `docs/specs/gate-resolution-frontmatter-contract.md`
11. `docs/schema/entity.mdschema.yml`
12. `docs/site/reference/command-reference.md`
13. `docs/site/reference/frontmatter-contract.md`
14. `docs/dev/README.md`
15. `skills/feedback-rejection-flow/SKILL.md`
16. `skills/commission/references/templates/development.md`

This is the candidate, not a lower range: all 16 files are expected to change. `internal/cli/pi_launch_test.go` contains a historical comment mentioning “feedback-cycle-2” but no recorder policy or command dependency, so it is deliberately outside the candidate. `zbc`, `nth`, and `jc` entity/worktree files are also outside this task.

Tolerance: the file cap is 16 (zero unlisted-file tolerance), with up to 650 deletions and 300 insertions across those files for compile fixes or existing contract fixtures. A required 17th file is evidence that this inventory is incomplete and requires a design reset before editing it. Deleting or changing the identity/immutability semantics of `review-round`, its two-file room, `gate record --round`, or `gate validate --round` is outside tolerance, as is any change to the neutral `jc` lookup/`nth` stage-route seam, binding `gates:`, gate application/status semantics, or runtime-host behavior; any such change requires a design reset regardless of LOC.

## Risk exercise

Exercised 2026-08-02 from the repository root. `go test ./internal/ensigncycle -run TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl -count=1` passed: its positive control creates `review-round: round:rejection-task:validation:1`, exactly two immutable room files, and no ordinary gate application, while its inverted no-invocation control fails. That is concrete evidence that `zbc` has a producer boundary worth preserving. A production-source inventory over `internal/gates` and `internal/cli` finds 28 policy references matching `correct-but-disproportionate|actor:ensign|feedbackCycleRE|FeedbackCyclePath|feedback-cycle|spliceFeedbackCycle|classifyCompletedRound|\.Triage|Triage:`; all are confined to the classifier, projection, input plumbing, or triage output, while `readRoundPointerData`, `publishRound`, and `ValidateRoundFile` supply the independent structural producer. This demonstrates the riskiest boundary: policy can reach zero without deleting the pointer/room producer. A help-only rename cannot produce the required absence because the exercised classifier/projector remains reachable; wholesale round deletion breaks the passing zbc oracle.

## Acceptance criteria

**AC-1 — The stable CLI retains a neutral correction-round producer with zero workflow-policy grammar while every binding gate verb remains callable.**

Verified by a binary-level CLI test that asserts help contains neutral round record/validate plus prepare, both binding record forms, eligibility, and consume; asserts no `--feedback-cycle` or `triage=`; records and validates structurally valid non-development-labeled round bytes; and invokes `--feedback-cycle` expecting exit 2/unknown flag plus byte-identical fixture state. Any policy flag/output surviving, neutral producer failure, filesystem mutation on rejection, or retained binding verb missing/failing its established fixture falsifies AC-1.

**AC-2 — Generic gate code has zero reachable development finding-policy or cycle-projection behavior.**

Verified by deleting classifier/projection tests while retaining producer tests, then running `rg -n 'correct-but-disproportionate|actor:ensign|feedbackCycleRE|FeedbackCyclePath|feedback-cycle|spliceFeedbackCycle|classifyCompletedRound|\.Triage|Triage:' internal/gates internal/cli --glob '*.go' --glob '!**/*_test.go'`. The measured 2026-08-02 baseline is 28 matching production lines and the target is zero; a new producer test supplies a structurally valid log with workflow-foreign actors/labels and must publish and validate unchanged. Any inventory hit, rejection based on development labels/actors, body projection, or binding test regression falsifies AC-2.

**AC-3 — The development workflow completes a correction pass through the retained neutral producer without delegating workflow policy to generic gate code.**

Verified by adapting `TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl`: the workflow fixture prewrites its authorized Cycle line, invokes `gate record --round validation/1` without `--feedback-cycle`, and must retain the exact `review-round` identity and two immutable room files while leaving the Cycle line and binding `gates:` bytes unchanged; its existing no-invocation control must still fail. Skill smoke/contract tests must show First-Officer ownership of classification/projection, reviewer rerun, and gate re-entry. Missing pointer/room production, any binary-written Cycle bytes, any generic label/actor judgment, or any instruction that passes the Cycle line to the binary falsifies AC-3. This is the producer `zbc` names as its authoritative correction-episode key; no live host test is needed because host handoff semantics do not change.

**AC-4 — Binding durable decisions remain format- and behavior-compatible.**

Verified by the existing gate package and CLI fixtures for prepare, room/chat Resolution record, validate, withdraw, eligibility, consume, and terminal merge guard, followed by `go test ./...` and `go test ./... -race`. Focused route fixtures must also retain status-derived unique-stage selection before reduction, stage-graph-derived advance, approve-only Application, revise/hold Resolution-only behavior, and the existing `Eligibility.Action=advance` plus `application=advance/<state>` projection. Any fixture-byte change to `gates:`/binding rooms, changed application/status result, round-derived route, or policy field added to the neutral Application falsifies AC-4.

## Documentation diff

- `docs/specs/gate-resolution-frontmatter-contract.md`: replace development triage/projection claims with: “The neutral round producer retains canonical Briefing/log bytes and advances `review-round`; it does not classify findings or write workflow body projections.” Retain the immutable room/pointer contract needed by `zbc`.
- `docs/schema/entity.mdschema.yml`: retain the `review-round` writer and room pointer; remove `### Feedback Cycles` from binary-owned projection claims and say its grammar/ownership is workflow-defined.
- `docs/site/reference/command-reference.md`: retain `gate record ... --round` and `gate validate ... --round`, remove `--feedback-cycle` and `triage=` wording, and describe structural retention only.
- `docs/site/reference/frontmatter-contract.md`: retain the paragraph defining `review-round` as the current structural pointer; replace triage/projection claims with the neutral producer boundary.
- `docs/dev/README.md` and the commissioned development template: replace “retain the advisory round with `${SPACEDOCK_BIN:-spacedock} gate record`” with “After FO consultation, the First Officer appends the workflow-defined Cycle line directly, then invokes the neutral round recorder with the canonical Briefing/log. The recorder retains bytes and advances `review-round`; it does not parse the Cycle line or classify findings.”
- `skills/feedback-rejection-flow/SKILL.md` and directly affected FO references: retain opaque authorized-package append, cycle-3 escalation, and invocation of the neutral producer; remove any implication that the binary classifies findings or owns the projection.

## Test plan

Implementation begins with CLI tests for AC-1 and the split positive/negative producer oracle for AC-3, then deletes classifier/projection code until the AC-2 inventory reaches zero. Run focused gate/CLI producer and binding tests for AC-2 and AC-4, including the existing neutral lookup/reducer/Application projections named above; run schema/spec contract checks and commission/feedback-rejection skill smoke tests for AC-3; then run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. Before implementation claims completion, compare `git diff --name-only` with the exact 16-file list and reset the design on any 17th file. Estimated complexity is medium deletion work: fixture and CLI tests are required; one fixture-backed workflow scenario is required; a live workflow/host run is not required because neither runtime dispatch nor host integration changes. Do not build a replacement recorder, compatibility alias, hidden legacy mode, migration, or schema version, and do not modify `zbc`, `nth`, or `jc` implementation code.

## Stage Report: ideation

- DONE: Produce a concrete, bounded problem statement, proposed approach, streamlined scenario, and expected surface for removing workflow-specific advisory rounds.
  The body bounds removal to advisory round code/commands/contracts, estimates 13-18 files and 1,150-1,500 deletions, and declares numeric and semantic tolerance.
- DONE: Define independently falsifiable end-state acceptance criteria with a reproducible test plan, declared semantic scope, tolerance, and required documentation diff.
  AC-1 through AC-4 name observable failure conditions, independent baselines, focused/full/race tests, on-disk evidence, and exact documentation changes.
- DONE: Exercise or explicitly record the riskiest unverified mechanism and explain why the simplest alternative cannot deliver the value.
  Focused binding and round tests exercised their separation; the body traces the call boundary and explains why rename/help-only approaches retain the jurisdiction error.

### Summary

The ideation defines a deletion-first v1 cut: remove advisory round grammar and development policy from generic gates while preserving binding decisions byte-for-byte. It records measured baselines, bounded surface and semantics, workflow-owned correction evidence, reproducible tests, and concrete documentation edits without designing a replacement recorder.

### Feedback Cycles

- Cycle 1: REVISE — Captain-directed send-back preserving the Science Officer finding: AC-2 and AC-3 lack concrete evidence, and the design must preserve the zbc correction-round producer. Correction assignment: add concrete AC evidence and retain the zbc producer while keeping advisory-round ownership in the workflow/First Officer.

## Stage Report: ideation (cycle 2)

- DONE: Add concrete evidence for AC-2 and AC-3 that can falsify the proposed v1 deletion boundary.
  AC-1 pins retained/removed grammar; AC-2 records the reproducible 28-line production baseline and zero target; AC-3 names the producer oracle and inverted control; AC-4 keeps binding fixtures byte-stable.
- DONE: Preserve the zbc correction-round producer while keeping generic gate code free of workflow policy.
  The design retains `gate record/validate --round`, `review-round`, and the immutable two-file room consumed by zbc, while deleting only the classifier, `--feedback-cycle` projection, and `triage=` output.
- DONE: Rerun the independent AC scan and record the refreshed bounded surface before re-gating ideation.
  Fresh `status --read ... --stage ideation --ac-scan --json` emitted `unevidenced:false` for AC-1 through AC-4; the baseline is 10-14 files, 300-520 deletions, at most 220 insertions, with tolerance of 16 files, 650 deletions, and 300 insertions.

### Summary

Cycle 2 narrows the deletion to development-specific policy and preserves the neutral correction-round producer required by zbc. AC-2 and AC-3 now carry measured baselines, executable controls, exact expected on-disk state, and explicit falsifiers; the mandatory independent scan cites both checklist evidence lines.

## Stage Report: ideation (cycle 3)

- DONE: Reconcile the exact 16-file wj surface with the expected-surface estimate and tolerance, and update the design so the estimate cannot understate the candidate.
  The design now enumerates all 16 expected files, makes 16 the hard cap with a reset at file 17, and records measured evidence `candidate_files=16`; deletion/insertion tolerance remains 650/300 only within those named files.
- DONE: Make the shared neutral stage-route seam explicit for wj, nth, and jc while preserving zbc's correction-round producer and the existing round room.
  Wj preserves jc's upstream unique-stage lookup, nth's stage-graph/approve-only Application boundary, Resolution-only revise/hold, and zbc's unchanged `review-round` plus immutable two-file room; any route or room semantic change is a design reset.
- DONE: Rerun the authoritative ideation checklist and AC scan, and record the revised bounded design and documentation diff.
  AC-1 evidence: the CLI oracle rejects policy grammar while retaining round/binding verbs. AC-2 evidence: the independent production inventory measures 28 policy lines with target zero. AC-3 evidence: the producer oracle passes with the exact pointer/two-file room and its no-invocation control fails. AC-4 evidence: existing binding/route fixtures remain the mandatory full/race-test compatibility baseline.

### Summary

Cycle 3 replaces the understated range with an exact, enumerated 16-file candidate and zero unlisted-file tolerance. It also separates the three shared boundaries cleanly: jc selects, nth reduces/routes, wj removes round policy while retaining zbc's structural producer; the documentation diff and executable AC falsifiers remain bounded to wj.
