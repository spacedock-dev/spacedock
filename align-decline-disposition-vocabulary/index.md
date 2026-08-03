---
title: Accept all declared non-material finding classes
status: ideation
source: "Captain-approved local task from debrief issue 2, 2026-08-02"
score: 0.7
worktree:
issue:
id: shra0x0r2bf7ka0q1m4ft79a
gates:
    version: 1
    records:
        - id: gate:shra0x0r2bf7ka0q1m4ft79a:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:shra0x0r2bf7ka0q1m4ft79a-backlog-1
              briefing:
                id: briefing:shra0x0r2bf7ka0q1m4ft79a:backlog:attempt-1:revision-1
                digest: sha256:f49fe5e9642d6d82e859daaa33ac43e37ebf1482336fb24b7904c304d67dc8f0
                request-digest: sha256:ed9980aecefd08f4b0cc43f731466b3fecd8925fb3afba6d8f6a315795043b0a
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:shra0x0r2bf7ka0q1m4ft79a:backlog:1
                briefing: briefing:shra0x0r2bf7ka0q1m4ft79a:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T01:56:25.061748Z"
                decision: approve
                reason: Captain approves the corrected opaque workflow-owned finding-class direction; enter ideation to rewrite the task body, ACs, and test plan.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:shra0x0r2bf7ka0q1m4ft79a:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:shra0x0r2bf7ka0q1m4ft79a-ideation-1
              briefing:
                id: briefing:shra0x0r2bf7ka0q1m4ft79a:ideation:attempt-1:revision-1
                digest: sha256:e5f97bfce5a2299aa57619f4ae2cd821a4525030a67bacd6ae285ba6a042a53b
                request-digest: sha256:248cf154f72b2e140024d082e8001703eba6e83b55c9c672832df7a3b5704200
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:shra0x0r2bf7ka0q1m4ft79a:ideation:1
                briefing: briefing:shra0x0r2bf7ka0q1m4ft79a:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T03:32:43.98711Z"
                decision: approve
                reason: 'Captain authorizes approval after Science Officer agreement: classes remain opaque and workflow-owned, no generic allowlist is added, WJ owns parser-policy removal, and sh supplies only residual proof.'
              application:
                target-stage: implementation
                state: pending
started: 2026-08-03T02:00:11Z
---

Generic correction-round recording must preserve workflow-owned finding content without
predefining or restricting its class vocabulary.

## Problem

The generic correction-round recorder currently recognizes the development workflow's
`correct-but-disproportionate` decline grammar. The active workflow instead declares
Material, Deferred risk, Polish, and Needs decision, and another workflow may declare a
different taxonomy. A generic allowlist therefore rejects structurally valid workflow
content and gives the launcher policy authority that belongs to the workflow.

## End value

Any workflow can retain an arbitrary well-formed finding class and disposition in a
correction-round Annotation without teaching generic gate code that vocabulary. The
retained log remains byte-exact and subject to the recorder's workflow-neutral integrity
checks.

## Proposed approach

Land this task after WJ (`cut-workflow-specific-round-recorder-from-v1`). WJ owns deletion
of `declineDispositionRE`, `dispositionKind`, `classifyCompletedRound`, actor/taxonomy
judgments, and the development-specific projection path. This task does not repeat or
partially reimplement that deletion.

Add one package-level round-recorder regression. Starting from the existing structurally
valid fixture, give an Annotation a deliberately workflow-foreign body such as
`class: orbital-debt; disposition: carry-until-reentry; basis: workflow-local policy`.
Record the round, validate it through the retained round API, and compare the retained
`briefing.review.jsonl` with the supplied bytes. Then attempt a divergent replay whose
only change is the opaque class/disposition body; it must fail and leave the retained
room and entity byte-identical.

The arbitrary token is test-owned input, not a production constant. The test therefore
fails if generic code later restores an allowlist for the development classes or parses
workflow body semantics.

## Explicit no-allowlist decision

Generic gate code has no finding-class or disposition allowlist. It does not enumerate
`deferred-risk`, `polish`, `correct-but-disproportionate`, the test's `orbital-debt`, or
any other workflow label. Adding declared classes to the existing regular expression,
moving the list to another generic table, or accepting an extensible registry in the
round command is explicitly rejected: each alternative still makes generic code the
taxonomy owner.

## Semantic boundaries

- Command grammar: unchanged by this task; WJ's neutral `gate record --round` and
  `gate validate --round` surfaces are the prerequisite baseline.
- Stored formats: unchanged. `review-round` and the canonical two-file room remain, and
  the exact supplied Briefing and JSONL bytes remain immutable.
- Authority: workflow policy alone defines class names, materiality, required evidence,
  and disposition meaning. Generic code owns only record structure and integrity.
- Runtime validation retained: valid JSONL shape, unique and ordered identities,
  attribution, same-Briefing association, backward `includes` references, canonical
  Briefing/artifact binding, derived room shape, pointer consistency, and immutable
  replay behavior.
- Runtime validation forbidden: interpreting Annotation body fields, class labels,
  disposition labels, rationale vocabulary, promotion rules, or workflow actor roles.
- No workflow, skill, schema, command-output, host adapter, or binding-gate behavior is
  changed by this task. WJ owns its already-approved contract and documentation edits.

## Expected surface and tolerance

Expected implementation surface: exactly one file, `internal/gates/round_test.go`, with
35-70 inserted lines and at most 10 deleted lines, reusing the existing advisory-round
fixture and byte-tree helpers.

Tolerance: still exactly one file, at most 100 insertions and 20 deletions for fixture
adaptation. Any production-code change, documentation or skill edit, second file, or
change to command grammar, stored format, authority, or runtime behavior requires a
design reset. WJ's implementation must be merged first; an implementation worker must
not make this regression pass by deleting or modifying the parser itself.

## Mechanism and alternatives

The dedicated record/validate/divergent-replay regression serves AC-1 by exercising the
supported producer and observing disk bytes. Merely adding `deferred-risk` and `polish`
to the old regex is simpler locally but cannot deliver the end value because the next
workflow-defined class fails. Relying only on WJ's deletion inventory is also
insufficient: zero named development tokens does not prove that an arbitrary body can
traverse record, persistence, and validation unchanged.

No spike needed. Existing tests already prove the neutral producer, exact replay,
divergent replay refusal, and byte-clean failure helpers; WJ's approved design retains
those mechanisms. This task combines those proven paths with workflow-foreign body data
instead of introducing a mechanism.

## Out of scope

Parser deletion and generic-policy cleanup (WJ); changes to the development taxonomy;
classification or authorization policy; Feedback Cycles projection; round storage;
binding gate behavior; CLI output; migrations; registries; and compatibility modes are
out of scope.

## Acceptance criteria

**AC-1 - A structurally valid correction round with an arbitrary workflow-owned finding class records and validates without generic taxonomy knowledge, while preserving the exact supplied log bytes.**

Verified by: the new round test supplies `orbital-debt` (a label absent from production
vocabulary), records and validates successfully, and byte-compares the retained JSONL
with the independent input. Replacing opaque acceptance with any class allowlist makes
the positive control fail.

**AC-2 - Opaque workflow content does not weaken structural identity, association, reference, or immutability checks.**

Verified by: the same regression changes only the already-retained opaque body and
attempts a replay; the operation must fail and the entity plus retained room must remain
byte-identical. Existing round tests continue to reject malformed JSONL, duplicate or
forward identities, cross-Briefing records, invalid `includes`, bad room shape, and
changed retained bytes. Removing one of those guards or permitting the divergent replay
falsifies AC-2.

## Test plan

After WJ is merged, add the focused fixture-backed regression in
`internal/gates/round_test.go`. Run its positive record/validate control and its
divergent-replay negative control together; the observable proof is exit/error behavior
plus resulting on-disk bytes, not source inspection. Then run `gofmt -w ./cmd ./internal`,
`go test ./...`, and `go test ./... -race`. Estimated complexity is low: one Go test
file, no new fixture file, no CLI fixture, and no live workflow/host run because command
and runtime-host semantics do not change.

## Stage Report: ideation

- DONE: Flesh out the opaque workflow-owned finding-class design and remove the generic allowlist target.
  The design explicitly rejects all generic allowlists/registries and assigns taxonomy semantics solely to workflow policy, with WJ owning parser deletion.
- DONE: Rewrite the entity acceptance criteria and test plan around structural acceptance of arbitrary well-formed classes.
  AC-1 uses a production-foreign class and exact disk bytes; AC-2 couples it to a divergent-replay refusal and the existing structural failure suite.
- DONE: Record the expected surface, tolerance, semantic boundaries, and WJ coordination seam.
  The follow-on is capped at one test file and 100 insertions/20 deletions, requires WJ first, and forbids production, format, command, authority, and host changes.

### Summary

Reframed the task from expanding a generic class parser to pinning workflow ownership
with a falsifiable opaque-body regression. The bounded follow-on adds only a round test;
WJ remains the sole owner of parser deletion and shared contract updates.
