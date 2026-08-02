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
        gate: gate:nthcevf1snz7hm75gny3kd2e:ideation
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
        - id: gate:nthcevf1snz7hm75gny3kd2e:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:nthcevf1snz7hm75gny3kd2e-ideation-1
              briefing:
                id: briefing:nthcevf1snz7hm75gny3kd2e:ideation:attempt-1:revision-1
                digest: sha256:d7ad337095f1c498c77434c61e75c6fcefe5d4934c2afd723373d243ea290336
                digest-domain: canonical-bytes
                request-digest: sha256:c4ea8d1d3f35a04792d3ef1191cf42e17433969db2276446db4cb9b4c4c8afdf
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nthcevf1snz7hm75gny3kd2e:ideation:1
                briefing: briefing:nthcevf1snz7hm75gny3kd2e:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-01T16:08:30.024517Z"
                decision: revise
                reason: 'Captain directed send-back for Science Officer REVISE: AC-2 and AC-3 are planned only, with no concrete evidence; coordinate the design with jc and the internal-gates schema before re-gating.'
              application:
                action: feedback
                target-stage: ideation
                state: superseded
            - id: gate-attempt:nthcevf1snz7hm75gny3kd2e-ideation-2
              briefing:
                id: briefing:nthcevf1snz7hm75gny3kd2e:ideation:attempt-2:revision-1
                digest: sha256:38b43486d7e09452ef4f5eb27e367fd09822d4bcca8fee67e6a845ab1d23f7c2
                digest-domain: canonical-bytes
                request-digest: sha256:5ec7f938d894704f866a320c687dabe108f7dce7221768ec092201d0e8143da7
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:nthcevf1snz7hm75gny3kd2e:ideation:2
                briefing: briefing:nthcevf1snz7hm75gny3kd2e:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-02T10:02:00.481726Z"
                decision: hold
                reason: Science confirms NTH evidence is complete but approval would authorize premature implementation. Hold until Captain-authorized JC scope reset and durable selector/API; then re-present NTH and WJ in order.
              application:
                action: none
                state: not-applicable
            - id: gate-attempt:nthcevf1snz7hm75gny3kd2e-ideation-3
              briefing:
                id: briefing:nthcevf1snz7hm75gny3kd2e:ideation:attempt-3:revision-1
                digest: sha256:feac214c3d3d152dbf626f9dd34aa34c5653c1d1e34990cba174ebdc13392e96
                digest-domain: canonical-bytes
                request-digest: sha256:c15f2c195fb130e2b6eee9f3674b5bdae17895e789a205eee7f96216e37c6d05
                room-ref: ./review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:nthcevf1snz7hm75gny3kd2e:ideation:3
                briefing: briefing:nthcevf1snz7hm75gny3kd2e:ideation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-08-02T14:48:27.141744Z"
                decision: approve
                reason: 'Science advisory approves the cycle-2 ideation direction now that JC #599 is merged and the WJ/NTH seam is explicit; implementation must use a fresh current-head binding and preserve the serialized JC-to-NTH-to-WJ order.'
              application:
                action: advance
                target-stage: implementation
                state: pending
                blockers: []
---

Keep only application state that a supported command creates and a supported consumer spends. The unreleased v1 model currently carries blockers, execution holds, feedback payloads, and leaf fields that were designed speculatively but have no live producer or consumer.

## Problem

The unreleased canonical-v1 `application` object models six leaf keys (`action`, `target-stage`, `state`, `blockers`, `execution-hold`, and `feedback`) and three actions, but the shipped journey has only one value transaction: an approval authorizes one later stage transition or terminal delivery. The active-state audit found 16 gate-bearing tickets; 14 carry only the recorder-emitted `blockers: []`, and none carries a nonempty blocker, execution hold, or feedback payload. There is no supported command that creates or updates those three concepts and no consumer that can satisfy a blocker, release a hold, or spend feedback.

Strict `KnownFields` decoding turns every accepted leaf into a durable format commitment despite that absence. The recorder also emits `application: {action: none, state: not-applicable}` for `hold` and `application: {action: feedback, ...}` for `revise`, duplicating facts already owned by the durable Resolution and the First Officer's declared `feedback-to` route. Keeping these shapes makes the v1 contract larger without preserving authority or enabling behavior.

## Proposed approach

Make `application` an approval-only authority token with exactly two required fields:

```yaml
application:
  target-stage: implementation
  state: pending # later consumed or superseded
```

`gate record --decision approve` remains the only producer, so the enclosing approved Resolution makes `action: advance` redundant. `gate consume` consumes or supersedes a nonterminal token; terminal delivery consumes it and `merge guard --rework` supersedes it. `hold` and `revise` close the attempt with a Resolution and no application. Remove `Action`, `Blocker`, `ExecutionHold`, and application `Feedback` from the Go model, remove their eligibility and validation branches, stop requiring an empty blockers list, and constrain application validation to a nonblank target and `pending|consumed|superseded`.

Normalize tracked pilot entities and fixtures in the same unreleased-v1 cut. Afterward, strict decoding rejects `action`, `blockers`, `execution-hold`, and `feedback`; validation rejects `not-applicable`. Do not add schema versioning, dual decoding, migration commands, or unknown-field preservation.

Mechanism-to-value mapping: the model/validator deletion and one-time fixture/state normalization serve AC-1 and AC-2 by making the smaller format observable and exclusive. The simplest alternative, merely stopping new emission while retaining decode support, is insufficient: it leaves all six leaf keys in the accepted public contract and cannot make removed-shape negative tests fail closed. Retaining constant `action: advance` is likewise insufficient because the enclosing approval already fixes that operation and no consumer can observe an alternative. Omitting applications for non-approvals serves AC-1 and AC-3; retaining `none`/`feedback` wrappers is insufficient because no supported application consumer spends either wrapper and the Resolution already supplies the decision evidence.

## Streamlined common scenarios

- **Approve, nonterminal:** recording writes the Resolution plus one `<target>/pending` token. `gate consume` atomically changes status and token state to `consumed`; repeats have no effect.
- **Approve, stale input:** eligibility detects changed reviewed input and `gate consume` changes only `pending` to `superseded`.
- **Approve, terminal:** recording writes the same minimal pending token. `gate consume` leaves it pending and reports `approved-awaiting-merge`; terminal delivery is the sole consumer, while `merge guard --rework` supersedes it.
- **Hold:** recording writes only the `hold` Resolution and stops at the gate. Eligibility projects ineligible from the absent application.
- **Revise:** recording writes only the `revise` Resolution. The First Officer routes its reason/findings through the workflow stage's `feedback-to`; no application duplicates that route.
- **Old/prototype input:** any removed leaf or state fails strict read/validation without changing entity bytes.

## Out of scope

Do not redesign terminal delivery, current-gate selection, Resolution contents, review-round storage, feedback routing, provider presentation, or top-level status projections. Do not add blocker/hold authoring commands or retain prototype fields for compatibility.

## Declared semantic scope

- **Command grammar:** unchanged; no verbs or flags are added or removed.
- **Stored format:** intentionally breaking within unreleased v1. Approval applications lose `action`, `blockers`, `execution-hold`, and `feedback`; hold/revise attempts lose `application`; `not-applicable` becomes invalid.
- **Authority:** unchanged. Resolution remains the durable decision; only approval application state is spendable authority.
- **Runtime behavior:** approve/consume/terminal-delivery behavior is preserved. Hold still stops and revise still routes through First Officer workflow logic, but neither emits dead application state. Status/readiness must derive those outcomes from Resolution plus application absence.

## Expected surface and tolerance

Expected implementation surface: 8-12 product/test/doc files, about 80-160 insertions and 140-260 deletions (net negative), concentrated in `internal/gates/{model,operation,application}.go`, their focused tests, affected `internal/status`/`internal/cli` fixtures, `docs/specs/gate-resolution-frontmatter-contract.md`, `docs/site/reference/frontmatter-contract.md`, and normalized `docs/dev/.spacedock-state` entity files. The exercised pre-cut manifest contains 31 canonical-v1 pilot entities (16 active, 15 archived); those mechanical state normalizations are counted separately from the product/test/doc surface.

Tolerance: up to 14 product/test/doc files, up to 220 insertions, and up to 35 normalized state/fixture files are acceptable if every extra file is a direct occurrence of a removed shape. Any new command, schema version, compatibility reader, migration engine, feedback router, or net-positive application-model LOC is outside tolerance and requires a new gate.

## Acceptance criteria

**AC-1 - The canonical v1 application value surface is reduced from six leaf keys to exactly two, and only approvals carry it.**
Verified by: a table-driven recorder test closes approve, hold, and revise attempts and decodes resulting YAML nodes (not Go zero values). It asserts exact approval keys `{target-stage,state}` and no application node for hold/revise: a 67% leaf-key reduction. The test fails if any other key is emitted or either non-approval gains an application.

**AC-2 - The clean v1 schema is the only accepted stored format after the cut.**
Verified by: table-driven strict-read negatives independently inject `action`, `blockers`, `execution-hold`, and `feedback`, plus invalid `not-applicable`; each must return an error and leave fixture bytes unchanged. Positive fixtures cover only pending, consumed, and superseded approval tokens. A checked-in 31-path pilot manifest is then iterated by a Go test through `gates.Read` and `Validate` for all 16 active and 15 `_archive` entities after one-time normalization; any legacy shape, unknown field, omitted path, or decode failure fails the run. `spacedock gate validate` remains an active-entity operator check, not evidence for archive coverage.

**AC-3 - Removing non-approval applications and speculative leaves does not change approval spending, hold stopping, revise routing, or terminal-delivery authority.**
Verified by: real-CLI fixtures drive (a) approve then repeated consume, observing one status transition, `pending→consumed`, and derived `application=advance/consumed`; (b) stale approve, observing only `pending→superseded`; (c) hold and revise, observing unchanged status, preserved Resolution identity/reason, and absent application; and (d) terminal consume/finalize/rework, observing the existing pending/consumed/superseded transitions. Eligibility retains a derived `Action: "advance"`, so the existing CLI vocabulary `application=advance/<state>` does not require a stored `action`. A duplicate spend, lost Resolution, changed route, non-approval application, or changed CLI route fails its leg.

**AC-4 - The implementation stays within the approved deletion-oriented boundary.**
Verified by: `git diff --numstat` and changed-path classification against the implementation base. It fails if product/test/doc insertions exceed 220, product/test/doc files exceed 14, normalized files exceed 35, or any changed path introduces command grammar, versioning, compatibility, migration, or routing machinery not declared above.

## Test plan

1. Add the table-driven model/recorder tests first (about 40-60 LOC): exact emitted YAML keys for three decisions, three valid approval states, and one independent negative per removed key/state. These are package tests because they exercise strict decoding and canonical writes directly; changing any promised shape makes a named row fail.
2. Update the model, validator, producer, eligibility logic, and package fixtures (about 40-80 net insertions, with larger deletions). Run `go test ./internal/gates -count=1` and the focused status/CLI gate suites.
3. Drive repository-built CLI fixtures for approve/consume/stale/hold/revise and the existing terminal delivery round trip (about 30-60 test LOC). Assert exit code and resulting on-disk YAML/status, not output substrings alone.
4. Check in the exercised 31-path pilot manifest, normalize each listed active/archive entity once, and add a Go test that iterates every manifest path through `gates.Read`/`Validate`. Also run the repository-built `gate validate` over active fixtures, but do not cite it for archives: that command reads only the resolved active path and `status --validate` does not call `gates.Validate` for every gate-bearing file. No live workflow test is needed because the claim is deterministic decoder, CLI, and on-disk behavior.
5. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; record any environment-correlated baseline failure separately. Check expected-surface thresholds with `git diff --numstat` and path review.

## Documentation diff required

Apply these user-visible changes during implementation:

- In `docs/specs/gate-resolution-frontmatter-contract.md`, change “The application layer owns what an approval does through the typed `application` subtree” to “A closed approval carries the typed one-use `application` subtree; revise and hold are complete Resolutions with no application.” Remove `action: advance` and `blockers: []` from the canonical YAML and replace the application validation/invariants text with the two-field approval-only rule.
- In `docs/site/reference/frontmatter-contract.md`, replace “Status surfaces it as ... `not-applicable` for a reviewer hold. A present dependency blocker or an active execution hold keeps the application ineligible” with “Status surfaces an approval as `approved-pending`, `consumed`, or `superseded`. Hold and revise carry no application; their Resolution remains the durable decision.”
- In `docs/site/reference/command-reference.md`, remove “held or blocked approvals are refused” from `gate consume`; retain stale, exactly-once, and terminal-routing wording.
- In `docs/site/concepts/gates-and-decisions.md`, add after the three calls: “Only approve creates an application. Revise and hold are complete when their Resolution is recorded; workflow routing, not an application payload, handles feedback.”

## Risk spike

Exercised on 2026-08-01 with `go test ./internal/gates -run 'TestRecordClosureShapesApplication|TestEightCanonicalApplicationShapesReplayByteIdentical|TestEligibilityFailClosedTable' -count=1` (exit 0). The current implementation demonstrably emits applications for all three decisions, requires `action: advance` and `blockers: []` for approval eligibility, and round-trips blocker, execution-hold, feedback, none, and not-applicable shapes. This proves the risky mechanism is the coordinated strict-decoder/model cut: producer-only cleanup cannot achieve AC-1/AC-2 because retained typed fields remain accepted canonical state. The spike seeds the exact-shape and negative rows in steps 1-2.

## Cycle-2 evidence

The strict negative baseline is deliberately red-before-cut. On 2026-08-02, `TestEightCanonicalApplicationShapesReplayByteIdentical` passed all eight legacy rows, including `approval_held`, `portable_hold`, and both feedback states; `TestRecordClosureShapesApplication` also passed approve/revise/hold and proved the current recorder emits all three applications. Therefore each AC-2 negative is independent and meaningful: removing only emission leaves the decoder accepting the exercised old rows, while deleting any one negative row leaves that legacy leaf unproved. The implementation must invert this baseline with named `action`, `blockers`, `execution-hold`, `feedback`, and `not-applicable` rejection rows and assert unchanged input bytes after each failed read.

A throwaway Go harness (removed after the run) exercised `gates.Read` plus `Validate` over canonical entity paths in both roots rather than scanning prose or calling the active-only CLI. It found 16 readable active pilot entities with 44 applications (`action` 44, `blockers` 34) and 15 readable archived pilot entities with 67 applications (`action` 67, `blockers` 52). Those exact 31 paths are the normalization manifest; older pre-v1 archives and non-entity Markdown are outside this unreleased-v1 cut. The post-cut harness must report 31/31 valid, exactly `{target-stage,state}` on every remaining approval application, no application on hold/revise, and zero occurrences of the four removed keys or `not-applicable`.

The runtime baseline is also exercised, not inferred. On 2026-08-02 the focused gates, status, and real-CLI commands exited 0 for approve/repeated consume, stale supersession, crash-window non-reconsumption, terminal pending/finalize/rework, shared readiness, and CLI `application=advance/pending`. The implementation reuses those fixtures and changes their YAML assertions: authority transitions and CLI route remain identical while hold/revise applications disappear. This is concrete AC-3 evidence because changing a status transition, duplicating a spend, losing a Resolution, or dropping the derived action/display breaks an already-running leg.

## Coordination and merge seam

The three lanes meet at one neutral reducer boundary, not through compatibility fields. `jc` owns selection of the current attempt by authoritative entity status plus a unique stage record; its implementation is currently byte-for-byte held under `mod-block=scope-reset:captain`, so this task requires no new jc API and must not mutate that worktree. After a Captain-authorized reset, jc first lands the status-derived unique-stage lookup; nth then rebases and removes stored application policy fields while consuming that lookup, projecting `Eligibility.Action = "advance"` from an approved pending/consumed/superseded token, and preserving CLI `application=advance/<state>`. If jc's approved selector contract changes, nth returns to the gate rather than adding a selector or compatibility field.

`wj` confirmed this seam in its pushed Cycle-3 design at state commit `a0870d39c`: wj owns only correction-round policy/projection deletion and retains the neutral round producer/room; it does not alter route or application fields. The shared stage graph derives the approval successor; only approve stores `{target-stage,state}`; revise/hold remain Resolution-only; workflow `feedback-to` routing stays outside Application. Merge-order conflicts resolve as jc selector first, nth application reducer second, then wj's orthogonal round cleanup; no lane may reintroduce stored `action`, a policy-specific application shape, or a hidden compatibility field.

## Stage Report: ideation

- DONE: Produce a concrete, bounded problem statement, proposed approach, streamlined scenarios, and expected surface for minimizing the stable-v1 gate application schema.
  The body defines the approval-only two-field schema, six scenarios, explicit exclusions, and a bounded deletion-oriented file/LOC estimate.
- DONE: Define independently falsifiable end-state acceptance criteria with a reproducible test plan, declared semantic scope, tolerance, and required documentation diff.
  AC-1 through AC-4 name independent YAML/CLI/diff observations and their falsifying changes; the body declares all four semantic dimensions, thresholds, and exact documentation replacements.
- DONE: Exercise or explicitly record the riskiest unverified mechanism and explain why the simplest alternative cannot deliver the value.
  Focused gates tests exited 0 and proved current strict/canonical support for all dead shapes; the body explains why stopping emission alone leaves the accepted v1 surface unchanged.

### Summary

Ideation reduces application state to one approval-only `{target-stage,state}` token while preserving Resolution authority and existing consume/terminal behavior. It rejects compatibility machinery, requires a coordinated strict-format cut and one-time normalization, and supplies falsifiable package, CLI, state-validation, and surface-budget checks for implementation.

### Feedback Cycles

- Cycle 1: REVISE — Captain-directed send-back preserving the Science Officer finding: AC-2 and AC-3 are planned only, with no concrete evidence; coordinate the design with jc and the internal-gates schema before re-gating. Correction assignment: add concrete AC evidence and coordinate the reduced application schema with jc/internal-gates before fresh presentation.

## Stage Report: ideation (cycle 2)

- DONE: Add concrete falsifiable evidence for AC-1 through AC-4, including strict decode negatives, CLI behavior, normalization, and the current status of each scan.
  AC-1, AC-2, AC-3, and AC-4 are each cited here: focused gates/status/CLI suites exited 0; a throwaway `gates.Read`/`Validate` harness measured 16 active plus 15 archived pilot files and 111 applications; the body fixes post-cut negatives and deletion-budget falsifiers.
- DONE: Coordinate the approval-only schema with wj and jc through one neutral stage-route seam, with no hidden compatibility field or policy-specific application shape.
  Wj confirmed and pushed the same seam at state `a0870d39c`; jc remains byte-clean under Captain hold, and the body records the exact selector-to-reducer handoff and post-reset merge order while preserving derived Action/CLI output.
- DONE: Commit the Cycle-2 ideation report and rerun the authoritative checklist and AC scan before gate preparation.
  This report is the path-scoped state commit artifact; the authoritative `status --read nth --checklist` and `--ac-scan` results are rerun immediately before completion and cited in the commit handoff.

### Summary

Cycle 2 turns AC-2 and AC-3 from future claims into executable baselines and exact post-cut falsifiers, including a real archived-entity decoder path and a bounded 31-file normalization manifest. It also fixes the jc-to-nth-to-wj seam without compatibility storage: jc selects, nth reduces approval authority and derives display action, and wj removes only round policy.
