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
                state: pending
---

Remove the public `gate record --round` and `gate validate --round` surface from the stable v1 cut. The implementation is coherent only for this repository's development review policy; retaining it would freeze that policy into the generic durable-decision contract.

## Problem statement

`internal/gates` currently recognizes `material` and `correct-but-disproportionate`, decides fixed versus declined, requires `actor:ensign`, parses the development workflow's files/LOC/estimate/percentage/AC cycle sentence, and mutates a `### Feedback Cycles` body section. Those are workflow judgments and projections, not binding gate primitives. The command therefore gives a generic package authority it must not have, and a prose-only rename leaves that jurisdiction error intact.

The bounded v1 problem is to remove this advisory-round operation while preserving every binding gate operation and its stored `gates:` format. Historical entity files may still contain `review-round` or `### Feedback Cycles`; v1 simply stops creating or validating those advisory records.

## Proposed approach

Delete the round-only CLI flags, routing, summary output, recorder/validator, model types and regexes, room publication helpers, policy classifier, fixtures, tests, and generic schema/spec claims. Keep shared Briefing, Annotation, Resolution, canonicalization, and atomic-write machinery when a binding operation still uses it. Add command tests that pin the reduced help and reject `--round`, `--log`, and `--feedback-cycle` as unknown flags while exercising all retained binding verbs.

Move ownership of correction-round retention back to the active workflow text: reviewers report findings; workers investigate and propose materiality, ownership, and disposition; the First Officer authorizes candidate mutation and retains the workflow's report/`### Feedback Cycles` entry directly. The binary neither parses nor certifies those entries. This serves AC-3. The simplest alternative—renaming the command or treating its payload as opaque while retaining the present implementation—is insufficient because the code would still enforce the ensign taxonomy and project workflow prose (AC-2).

Keep `workflow-neutral-advisory-round-recorder` as the only reintroduction path. It is out of scope here and may return only with a structural storage graph whose classifications and projections are owned by the active workflow.

## Streamlined scenario

During implementation review, a reviewer reports a finding. The worker records evidence and proposes materiality, task ownership, and disposition in the stage report; the First Officer separately authorizes any fix and, when another pass is needed, directly appends the workflow-defined cycle line. Binding stage approval still uses `gate prepare`, `gate record`, `gate validate`, `gate eligibility`, and `gate consume`. No v1 gate command receives, interprets, or stores the advisory review round.

## Semantic scope

- Command grammar changes: remove record flags `--round`, `--log`, and `--feedback-cycle`, remove validate flag `--round`, and remove round usage/output. These spellings become usage errors (exit 2).
- Stored formats: stop creating `review-round` frontmatter, `review/<stage>/round-<cycle>/`, and binary-authored `### Feedback Cycles` lines. Do not migrate, reject, or rewrite historical occurrences. The binding `gates:` graph and retained gate-room formats are byte-compatible.
- Authority changes: finding classification, disposition, correction ownership, cycle projection, and escalation remain workflow/First-Officer responsibilities; the generic binary has none.
- Runtime behavior: only advisory-round record/validate calls cease to work. Binding prepare, record (room and chat), validate, withdraw, eligibility, consume, and merge behavior are unchanged.

## Expected surface and tolerance

Expected implementation surface: 13-18 files, approximately 1,150-1,500 deleted lines and no more than 180 inserted lines (tests and replacement workflow wording). Primary code is `internal/cli/cli.go`, `internal/cli/gate_test.go`, `internal/gates/{round.go,round_test.go,review.go,model.go,io.go,operation.go}`, and `internal/gates/testdata/advisory-round/`. Contract/documentation is `docs/specs/gate-resolution-frontmatter-contract.md`, `docs/schema/entity.mdschema.yml`, `docs/site/reference/{command-reference.md,frontmatter-contract.md}`, `docs/dev/README.md`, `skills/feedback-rejection-flow/SKILL.md`, and `skills/commission/references/templates/development.md`; directly affected FO references may be changed only where they instruct binary-backed round retention.

Tolerance: up to 20 files, 1,700 deletions, and 250 insertions is allowed when compile errors or existing contract fixtures reveal another direct round dependency. Any change to the binding `gates:` schema, gate application semantics, status transitions, or runtime-host behavior is outside tolerance regardless of LOC and requires a design reset.

## Risk exercise

Exercised 2026-08-01 from the repository root. `go test ./internal/gates ./internal/cli -run 'TestGate|TestRound|TestRecord|TestValidate' -count=1` passed; `go run ./cmd/spacedock gate --help` exposed both advisory round forms alongside all retained binding verbs; and `go run ./cmd/spacedock gate validate missing --round implementation/1 --workflow-dir docs/dev` routed through round-capable parsing before failing entity lookup. Source call inventory shows `classifyCompletedRound` is called only by `round.go`, while `RecordSemanticSummary` has separate room/chat branches for binding decisions. This demonstrates the riskiest boundary—round policy can be cut without replacing binding recording—and seeds the implementation tests. A rename cannot produce the required absence because the same exercised parser/classifier/projector would remain reachable.

## Acceptance criteria

**AC-1 — The stable CLI has zero advisory-round grammar while every binding gate verb remains callable.**

Verified by a binary-level CLI test that asserts help contains prepare, both binding record forms, validate, eligibility, and consume; asserts no `--round`, `--log`, `--feedback-cycle`, or round summary; and invokes each removed spelling expecting exit 2/unknown flag plus byte-identical fixture state. The user-visible baseline is two accepted advisory command forms in current help and the target is zero. Any removed spelling accepted, any filesystem mutation, or any retained verb missing/failing its established fixture falsifies AC-1.

**AC-2 — Generic gate code has zero reachable development finding-policy or cycle-projection behavior.**

Verified by deleting round-only tests/fixtures, running focused `internal/gates` binding tests, and a repository source inventory over shipped Go code for `correct-but-disproportionate`, `actor:ensign`, the Material fixed/declined classifier, feedback-cycle grammar, `review-round`, and `spliceFeedbackCycle`; the measured ideation baseline is 31 matching production references and the expected count is zero. A surviving reachable parser/projector symbol or a binding test regression falsifies AC-2.

**AC-3 — The development workflow completes a correction pass without claiming that the gate binary records advisory rounds.**

Verified by existing skill smoke/contract tests plus a fixture-backed correction-flow scenario that reaches reviewer rerun and gate re-entry while no generated instruction contains `gate record --round`, `gate validate --round`, or a claim that the gate binary classifies findings. Missing worker/First-Officer authority or any removed invocation falsifies AC-3. No live host test is needed because host handoff semantics do not change.

**AC-4 — Binding durable decisions remain format- and behavior-compatible.**

Verified by the existing gate package and CLI fixtures for prepare, room/chat Resolution record, validate, withdraw, eligibility, consume, and terminal merge guard, followed by `go test ./...` and `go test ./... -race`. Any fixture-byte change to `gates:`/binding rooms or changed application/status result falsifies AC-4.

## Documentation diff

- `docs/specs/gate-resolution-frontmatter-contract.md`: remove “Round records and triage dispositions (advisory; owner: 02av)” and round-only summary claims; add: “The v1 gate contract records binding gate attempts and Resolutions only. Advisory review rounds, finding classification, and workflow report projection are outside this contract.”
- `docs/schema/entity.mdschema.yml`: remove the `review-round` writer/atomic-projection claim and remove `### Feedback Cycles` from binary-owned sections; do not forbid either key/heading in historical prose.
- `docs/site/reference/command-reference.md`: delete the `gate record ... --round` and `gate validate ... --round` table rows; retained binding rows remain unchanged.
- `docs/site/reference/frontmatter-contract.md`: delete the paragraph beginning “`review-round` is the current pointer” and the sentence “Advisory-round Briefings remain summary-free”; retain the binding Briefing and `gates` contract unchanged.
- `docs/dev/README.md` and the commissioned development template: replace “retain the advisory round with `${SPACEDOCK_BIN:-spacedock} gate record`” with “After FO consultation, retain the reviewer and worker evidence in the stage report. When another pass follows, the First Officer appends the workflow-defined Cycle line directly; no gate command records or validates this advisory report.”
- `skills/feedback-rejection-flow/SKILL.md` and directly affected FO references: retain the opaque authorized-package append and cycle-3 escalation, but remove any implication that the gate binary parses, stores, or validates the round.

## Test plan

Implementation begins with CLI tests for AC-1, then removes routing and round-only package code until they pass. Run focused gate package/CLI tests for AC-2 and AC-4, schema/spec contract checks, commission and feedback-rejection skill smoke tests for AC-3, then `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. Estimated complexity is medium deletion work: fixture and CLI tests are required; one fixture-backed workflow scenario is required; a live workflow/host run is not required because neither runtime dispatch nor host integration changes. Do not build a replacement recorder, compatibility alias, hidden legacy mode, migration, or schema version.

## Stage Report: ideation

- DONE: Produce a concrete, bounded problem statement, proposed approach, streamlined scenario, and expected surface for removing workflow-specific advisory rounds.
  The body bounds removal to advisory round code/commands/contracts, estimates 13-18 files and 1,150-1,500 deletions, and declares numeric and semantic tolerance.
- DONE: Define independently falsifiable end-state acceptance criteria with a reproducible test plan, declared semantic scope, tolerance, and required documentation diff.
  AC-1 through AC-4 name observable failure conditions, independent baselines, focused/full/race tests, on-disk evidence, and exact documentation changes.
- DONE: Exercise or explicitly record the riskiest unverified mechanism and explain why the simplest alternative cannot deliver the value.
  Focused binding and round tests exercised their separation; the body traces the call boundary and explains why rename/help-only approaches retain the jurisdiction error.

### Summary

The ideation defines a deletion-first v1 cut: remove advisory round grammar and development policy from generic gates while preserving binding decisions byte-for-byte. It records measured baselines, bounded surface and semantics, workflow-owned correction evidence, reproducible tests, and concrete documentation edits without designing a replacement recorder.
