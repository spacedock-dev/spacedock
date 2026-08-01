---
title: Minimize the unreleased v1 gate application schema
status: ideation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: production emits empty blockers but has no producer or demonstrated consumer for blockers, execution holds, or feedback payloads."
started: 2026-08-01T14:01:07Z
completed:
verdict:
score: "0.9"
worktree:
issue:
pr:
sprint: durable-decisions
id: nthcevf1snz7hm75gny3kd2e
gates:
    version: 1
    current:
        gate: gate:nthcevf1snz7hm75gny3kd2e:backlog
    records:
        - id: gate:nthcevf1snz7hm75gny3kd2e:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:nthcevf1snz7hm75gny3kd2e-backlog-1
              briefing:
                id: briefing:nthcevf1snz7hm75gny3kd2e:backlog:attempt-1:revision-1
                digest: sha256:a969f6243a82f97def080d4c51f69732932877a9a96392b3df5ef48b28a7ba8f
                digest-domain: canonical-bytes
                request-digest: sha256:1c1c9a1b2ce8ff70509a08507db6e35cbea319b2411678da6e2a7d9aa3af160d
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nthcevf1snz7hm75gny3kd2e:backlog:1
                briefing: briefing:nthcevf1snz7hm75gny3kd2e:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:53:46.825149Z"
                decision: approve
                reason: Captain approved dispatching this durable-decisions ideation lane in parallel with wj, hq, and jc.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Keep only application state that a supported command creates and a supported consumer spends. The unreleased v1 model currently carries blockers, execution holds, feedback payloads, and leaf fields that were designed speculatively but have no live producer or consumer.

## Problem

The unreleased canonical-v1 `application` object models six leaf keys (`action`, `target-stage`, `state`, `blockers`, `execution-hold`, and `feedback`) and three actions, but the shipped journey has only one value transaction: an approval authorizes one later stage transition or terminal delivery. The active-state audit found 16 gate-bearing tickets; 14 carry only the recorder-emitted `blockers: []`, and none carries a nonempty blocker, execution hold, or feedback payload. There is no supported command that creates or updates those three concepts and no consumer that can satisfy a blocker, release a hold, or spend feedback.

Strict `KnownFields` decoding turns every accepted leaf into a durable format commitment despite that absence. The recorder also emits `application: {action: none, state: not-applicable}` for `hold` and `application: {action: feedback, ...}` for `revise`, duplicating facts already owned by the durable Resolution and the First Officer's declared `feedback-to` route. Keeping these shapes makes the v1 contract larger without preserving authority or enabling behavior.

## Proposed approach

Make `application` an approval-only authority token with exactly three required fields:

```yaml
application:
  action: advance
  target-stage: implementation
  state: pending # later consumed or superseded
```

`gate record --decision approve` remains the only producer. `gate consume` consumes or supersedes a nonterminal token; terminal delivery consumes it and `merge guard --rework` supersedes it. `hold` and `revise` close the attempt with a Resolution and no application. Remove `Blocker`, `ExecutionHold`, and application `Feedback` from the Go model, remove their eligibility and validation branches, stop requiring an empty blockers list, and constrain application validation to `action: advance`, a nonblank target, and `pending|consumed|superseded`.

Normalize tracked pilot entities and fixtures in the same unreleased-v1 cut. Afterward, strict decoding rejects `blockers`, `execution-hold`, and `feedback`; validation rejects `none`, `feedback`, and `not-applicable`. Do not add schema versioning, dual decoding, migration commands, or unknown-field preservation.

Mechanism-to-value mapping: the model/validator deletion and one-time fixture/state normalization serve AC-1 and AC-2 by making the smaller format observable and exclusive. The simplest alternative, merely stopping new emission while retaining decode support, is insufficient: it leaves all six leaf keys in the accepted public contract and cannot make removed-shape negative tests fail closed. Omitting applications for non-approvals serves AC-1 and AC-3; retaining `none`/`feedback` wrappers is insufficient because no supported application consumer spends either wrapper and the Resolution already supplies the decision evidence.

## Streamlined common scenarios

- **Approve, nonterminal:** recording writes the Resolution plus one `advance/<target>/pending` token. `gate consume` atomically changes status and token state to `consumed`; repeats have no effect.
- **Approve, stale input:** eligibility detects changed reviewed input and `gate consume` changes only `pending` to `superseded`.
- **Approve, terminal:** recording writes the same minimal pending token. `gate consume` leaves it pending and reports `approved-awaiting-merge`; terminal delivery is the sole consumer, while `merge guard --rework` supersedes it.
- **Hold:** recording writes only the `hold` Resolution and stops at the gate. Eligibility projects ineligible from the absent application.
- **Revise:** recording writes only the `revise` Resolution. The First Officer routes its reason/findings through the workflow stage's `feedback-to`; no application duplicates that route.
- **Old/prototype input:** any removed leaf, action, or state fails strict read/validation without changing entity bytes.

## Out of scope

Do not redesign terminal delivery, current-gate selection, Resolution contents, review-round storage, feedback routing, provider presentation, or top-level status projections. Do not add blocker/hold authoring commands or retain prototype fields for compatibility.

## Declared semantic scope

- **Command grammar:** unchanged; no verbs or flags are added or removed.
- **Stored format:** intentionally breaking within unreleased v1. Approval applications lose `blockers`, `execution-hold`, and `feedback`; hold/revise attempts lose `application`; `not-applicable`, `none`, and application `feedback` become invalid.
- **Authority:** unchanged. Resolution remains the durable decision; only approval application state is spendable authority.
- **Runtime behavior:** approve/consume/terminal-delivery behavior is preserved. Hold still stops and revise still routes through First Officer workflow logic, but neither emits dead application state. Status/readiness must derive those outcomes from Resolution plus application absence.

## Expected surface and tolerance

Expected implementation surface: 8-12 tracked files, about 80-160 insertions and 140-260 deletions (net negative), concentrated in `internal/gates/{model,operation,application}.go`, their focused tests, affected `internal/status`/`internal/cli` fixtures, `docs/specs/gate-resolution-frontmatter-contract.md`, `docs/site/reference/frontmatter-contract.md`, and normalized `docs/dev/.spacedock-state` entity files. Mechanical fixture/state normalization may raise the file count to 30 without raising semantic scope.

Tolerance: up to 14 product/test/doc files, up to 220 insertions, and up to 35 normalized state/fixture files are acceptable if every extra file is a direct occurrence of a removed shape. Any new command, schema version, compatibility reader, migration engine, feedback router, or net-positive application-model LOC is outside tolerance and requires a new gate.

## Acceptance criteria

**AC-1 - The canonical v1 application value surface is reduced from six leaf keys and three actions to three leaf keys and one action, and only approvals carry it.**
Verified by: a table-driven recorder test closes approve, hold, and revise attempts and decodes resulting YAML nodes (not Go zero values). It asserts exact approval keys `{action,target-stage,state}`, action `advance`, and no application node for hold/revise. The test fails if any removed key/action is emitted or either non-approval gains an application.

**AC-2 - The clean v1 schema is the only accepted stored format after the cut.**
Verified by: table-driven strict-read negatives independently inject `blockers`, `execution-hold`, and `feedback`, plus invalid `none`, application `feedback`, and `not-applicable`; each must return nonzero/error and leave bytes unchanged. Positive fixtures cover only pending, consumed, and superseded approval tokens. The repository-built CLI then validates every tracked active and archived pilot entity after one-time normalization; any legacy shape or unknown field fails the run.

**AC-3 - Removing non-approval applications and speculative leaves does not change approval spending, hold stopping, revise routing, or terminal-delivery authority.**
Verified by: real-CLI fixtures drive (a) approve then repeated consume, observing one status transition and `pending→consumed`; (b) stale approve, observing only `pending→superseded`; (c) hold and revise, observing unchanged status, preserved Resolution identity/reason, and absent application; and (d) terminal consume/finalize/rework, observing the existing pending/consumed/superseded transitions. A duplicate spend, lost Resolution, changed route, or non-approval application fails its leg.

**AC-4 - The implementation stays within the approved deletion-oriented boundary.**
Verified by: `git diff --numstat` and changed-path classification against the implementation base. It fails if product/test/doc insertions exceed 220, product/test/doc files exceed 14, normalized files exceed 35, or any changed path introduces command grammar, versioning, compatibility, migration, or routing machinery not declared above.

## Test plan

1. Add the table-driven model/recorder tests first (about 40-60 LOC): exact emitted YAML keys for three decisions, three valid approval states, and one independent negative per removed key/action/state. These are package tests because they exercise strict decoding and canonical writes directly; changing any promised shape makes a named row fail.
2. Update the model, validator, producer, eligibility logic, and package fixtures (about 40-80 net insertions, with larger deletions). Run `go test ./internal/gates -count=1` and the focused status/CLI gate suites.
3. Drive repository-built CLI fixtures for approve/consume/stale/hold/revise and the existing terminal delivery round trip (about 30-60 test LOC). Assert exit code and resulting on-disk YAML/status, not output substrings alone.
4. Normalize every tracked active/archive occurrence once, then run the repository's entity-validation command over both roots. No live workflow test is needed: the claim is deterministic CLI and on-disk behavior, already expressible by fixtures.
5. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; record any environment-correlated baseline failure separately. Check expected-surface thresholds with `git diff --numstat` and path review.

## Documentation diff required

Apply these user-visible changes during implementation:

- In `docs/specs/gate-resolution-frontmatter-contract.md`, change “The application layer owns what an approval does through the typed `application` subtree” to “A closed approval carries the typed one-use `application` subtree; revise and hold are complete Resolutions with no application.” Remove `blockers: []` from the canonical YAML and replace the application validation/invariants text with the three-field approval-only rule.
- In `docs/site/reference/frontmatter-contract.md`, replace “Status surfaces it as ... `not-applicable` for a reviewer hold. A present dependency blocker or an active execution hold keeps the application ineligible” with “Status surfaces an approval as `approved-pending`, `consumed`, or `superseded`. Hold and revise carry no application; their Resolution remains the durable decision.”
- In `docs/site/reference/command-reference.md`, remove “held or blocked approvals are refused” from `gate consume`; retain stale, exactly-once, and terminal-routing wording.
- In `docs/site/concepts/gates-and-decisions.md`, add after the three calls: “Only approve creates an application. Revise and hold are complete when their Resolution is recorded; workflow routing, not an application payload, handles feedback.”

## Risk spike

Exercised on 2026-08-01 with `go test ./internal/gates -run 'TestRecordClosureShapesApplication|TestEightCanonicalApplicationShapesReplayByteIdentical|TestEligibilityFailClosedTable' -count=1` (exit 0). The current implementation demonstrably emits applications for all three decisions, requires `blockers: []` for eligibility, and round-trips blocker, execution-hold, feedback, none, and not-applicable shapes. This proves the risky mechanism is the coordinated strict-decoder/model cut: producer-only cleanup cannot achieve AC-1/AC-2 because retained typed fields remain accepted canonical state. The spike seeds the exact-shape and negative rows in steps 1-2.

## Stage Report: ideation

- DONE: Produce a concrete, bounded problem statement, proposed approach, streamlined scenarios, and expected surface for minimizing the stable-v1 gate application schema.
  The body defines the approval-only three-field schema, six scenarios, explicit exclusions, and a bounded deletion-oriented file/LOC estimate.
- DONE: Define independently falsifiable end-state acceptance criteria with a reproducible test plan, declared semantic scope, tolerance, and required documentation diff.
  AC-1 through AC-4 name independent YAML/CLI/diff observations and their falsifying changes; the body declares all four semantic dimensions, thresholds, and exact documentation replacements.
- DONE: Exercise or explicitly record the riskiest unverified mechanism and explain why the simplest alternative cannot deliver the value.
  Focused gates tests exited 0 and proved current strict/canonical support for all dead shapes; the body explains why stopping emission alone leaves the accepted v1 surface unchanged.

### Summary

Ideation reduces application state to one approval-only `advance/target-stage/state` token while preserving Resolution authority and existing consume/terminal behavior. It rejects compatibility machinery, requires a coordinated strict-format cut and one-time normalization, and supplies falsifiable package, CLI, state-validation, and surface-budget checks for implementation.
