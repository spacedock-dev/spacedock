---
title: Simplify the unreleased v1 gate-state schema
status: implementation
source: "Durable-decisions sprint implementation-shape audit, 2026-07-24."
score: "0.7"
id: jccbpvjv5bg1jn0jbmj2yf8s
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:jccbpvjv5bg1jn0jbmj2yf8s:ideation
    records:
        - id: gate:jccbpvjv5bg1jn0jbmj2yf8s:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:jccbpvjv5bg1jn0jbmj2yf8s-backlog-1
              briefing:
                id: briefing:jccbpvjv5bg1jn0jbmj2yf8s:backlog:attempt-1:revision-1
                digest: sha256:4d61c845361a3aaba15ec68a09a5090b03c963c532d61227de695e4473ae32c3
                digest-domain: canonical-bytes
                request-digest: sha256:7b480163e7d4760504a04d5d91b979e8b942b0de4ab99b690f1619145e3c4db3
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:jccbpvjv5bg1jn0jbmj2yf8s:backlog:1
                briefing: briefing:jccbpvjv5bg1jn0jbmj2yf8s:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:54:16.651475Z"
                decision: approve
                reason: Captain approved dispatching this durable-decisions ideation lane in parallel with wj, hq, and nth.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:jccbpvjv5bg1jn0jbmj2yf8s:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:jccbpvjv5bg1jn0jbmj2yf8s-ideation-1
              briefing:
                id: briefing:jccbpvjv5bg1jn0jbmj2yf8s:ideation:attempt-1:revision-1
                digest: sha256:3cf59dd9edae9ccd6acbf748297061107c7e8d80760740852bad494fdf6d9c30
                digest-domain: canonical-bytes
                request-digest: sha256:0d2eaeb477ff763f7ab39613d61800400946b586b9451ddac93e951630379021
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:jccbpvjv5bg1jn0jbmj2yf8s:ideation:1
                briefing: briefing:jccbpvjv5bg1jn0jbmj2yf8s:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-02T01:32:41.555714Z"
                decision: revise
                reason: 'Per Science Officer advisory and the Captain-delegated send-back rule: split the exact Acceptance criteria section from Test plan so the mandatory AC scan passes; tighten the 12-18-file/80-180-insertions and 140-300-deletions scope to the minimum consumer surface; retain the stale-approval adversarial fixture and full/race plan. Repairing the entity invalidates this Briefing, so require a fresh current-stage binding.'
              application:
                action: feedback
                target-stage: ideation
                state: superseded
            - id: gate-attempt:jccbpvjv5bg1jn0jbmj2yf8s-ideation-2
              briefing:
                id: briefing:jccbpvjv5bg1jn0jbmj2yf8s:ideation:attempt-2:revision-1
                digest: sha256:c6b47fae9ccadec75149e579e644e171e0d9c4caac755b25b6667f3e75c5e555
                digest-domain: canonical-bytes
                request-digest: sha256:e4e77cf1029b03939cbb1a8066f1f6a7f3f3f3d6d01913a8b7bd3e04e62ad32c
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:jccbpvjv5bg1jn0jbmj2yf8s:ideation:2
                briefing: briefing:jccbpvjv5bg1jn0jbmj2yf8s:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-02T08:06:14.933867Z"
                decision: approve
                reason: Captain authorized sprint advancement when Science Officer advisory approves; Science approved ideation attempt 2 with the six-production-file boundary, strict 24-file/180-plus/230-minus/net-negative guard, reducer consistency, and merge-order constraints.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
started: 2026-08-01T14:01:17Z
---

The unreleased v1 gate-state implementation still carries prototype compatibility and a mutable current-gate pointer that duplicates derivable state and has already projected a stale approval.

Ideation should define the smallest clean-v1 schema:

- Remove `raw-file-pin` support and its compatibility fixtures; canonical bytes are the sole shipped binding.
- Determine whether `digest-domain` is redundant once the v1 schema fixes canonical digest semantics.
- Exercise multi-stage re-entry, multiple historical attempts, changed-Briefing supersession, and same-stage replay to determine whether `gates.current.gate` can be derived from current stage plus immutable records.
- If derivation is sufficient, remove the stored pointer and all pointer-repair/rebind behavior. If one counterexample requires stored selection, record that minimal counterexample and retain only the smallest non-stale selector.

The test must reproduce the observed failure class: an older approved attempt must not make a newly rejected candidate appear approved-awaiting-merge. No prototype migration or compatibility path is required.

## Problem statement

Canonical v1 currently stores `gates.current.gate` even though workflow `status` already names the active stage and each record names its stage. The two mutable selectors can disagree. Worse, `Validate` allows two records to claim one stage, so a stale pointer can select an older approved record while a newer same-stage candidate is rejected and make readiness report an approval that the captain did not grant to the current candidate.

The schema also repeats a fixed digest rule in every Briefing as `digest-domain: canonical-bytes`. Since unreleased v1 has no accepted alternative and already rejects `raw-file-pin`, that field and its compatibility-only refusal fixture add representation without choice. The clean v1 value is one authoritative route from workflow stage to the last ordered attempt, with canonical digest semantics fixed by the schema.

## Proposed approach

1. Remove `Document.Current`, `Selection`, and serialized `gates.current`. Resolve the active record by the entity's authoritative `status`, requiring exactly one `records[*].stage` match. Validation rejects duplicate record stages; no match means no gate exists yet for that stage, which prepare may create. The last ordered attempt remains current.
2. Route readiness, summary, eligibility, consume, delivery/finalize/rework, prepare, and record through the same stage lookup. Writers stop selecting or repairing a pointer. Historical records and attempts remain immutable and ordered; re-entry selects the existing stage record and appends an attempt.
3. Remove `Briefing.DigestDomain` and serialized `digest-domain`. All v1 Briefing and review-round digests are unconditionally canonical bytes. Remove `raw-file-pin` refusal/compatibility fixtures and domain switches; malformed extra fields still fail through closed-schema decoding.
4. Keep gate IDs, attempt IDs, room references, application semantics, retained-authority checks, and canonical digest computation unchanged. No migration or dual reader is added because this schema is unreleased.

The record-by-stage helper already exists as `recordForStage` and rejects ambiguous matches. The implementation should make it the single selection primitive (or a small equivalent that distinguishes create-time absence from read-time invalidity), rather than introduce a new index or selector.

## Common scenarios and expected behavior

- First gate at a stage: no matching record exists; prepare creates one, and subsequent reads derive it from `status`.
- Same-stage replay: an open identical preparation is a no-op; a divergent open candidate still fails closed.
- Multiple attempts: only the last attempt in the uniquely matched stage record affects readiness; older approvals cannot override a later revise/hold/open attempt.
- Cross-stage progress and re-entry: historical records for other stages remain untouched; returning to a gate stage reuses its one record and appends the next attempt.
- Changed Briefing after a closed attempt: preparation supersedes that attempt's pending application and appends a successor; readiness follows the successor without pointer rebinding.
- Duplicate stage records: read/validate fails closed with an ambiguity error and leaves bytes unchanged.
- Terminal approval: only the last attempt of the status-matched record can produce `approved-awaiting-merge`; a newly rejected candidate produces `feedback-pending` even if older attempts or other records were approved.

## Acceptance criteria

- **AC-1 — Current-candidate authority:** For an entity at a gate stage, readiness is determined only by the last attempt of the sole record whose `stage` equals entity `status`. A table-driven Go test covers open, approve, revise, hold, consumed, superseded, multi-attempt, cross-stage history, and re-entry. It fails if any lookup uses an ID selector or a non-last attempt.
- **AC-2 — Observed stale-approval class is eliminated:** A behavior fixture with an older terminal-target approval and a newly rejected same-stage candidate reports `feedback-pending` and never `approved-awaiting-merge`. A companion duplicate-stage fixture fails closed and preserves entity bytes. This is the independent end-value measure: false approval projections are **0 of 2 adversarial cases**, versus **1 of 2** under the current pointer-selected baseline.
- **AC-3 — Stored v1 is minimal and closed:** Newly prepared/recorded entities contain neither `gates.current` nor any `digest-domain`; canonical digest association, stale-input detection, exact replay, and room-backed recording still pass. Golden/entity-byte assertions fail if either removed field is emitted or canonical digest behavior changes.
- **AC-4 — Historical behavior survives:** Multi-stage records, multiple attempts, changed-Briefing supersession, same-stage replay, nonterminal consume, terminal finalize, and terminal rework retain their existing observable status/application outcomes. Focused gate package and status behavior tests exercise resulting bytes and readiness, followed by `go test ./...` and `go test ./... -race`.
- **AC-5 — Documentation names one authority:** The schema and gate-resolution contract describe status-derived unique-stage selection and implicit canonical digest semantics, with no normative `gates.current`, `digest-domain`, or `raw-file-pin` support. A review of the rendered YAML example plus repository search over normative spec/schema files verifies the old fields are absent; historical roadmap evidence is intentionally preserved.

## Test plan

Test cost is medium: Go unit and behavior fixtures only; no provider or live workflow is required. Implement AC-2 first as a table-driven fixture in `internal/gates/gates_test.go`: the stale-pointer baseline must reproduce one false approval across the two adversarial cases, while status-derived unique-stage selection must produce zero. The newer rejected-attempt case must return `feedback-pending`; the duplicate-stage case must return an ambiguity error, leave the entity bytes unchanged, and fail if validation accepts both records.

Focused package tests cover the remaining criteria. `go test ./internal/gates -count=1` must exercise unique-stage lookup, last-attempt readiness, prepare/replay/supersession, consume, terminal finalize/rework, canonical digest association, exact replay, and closed-schema rejection. `go test ./internal/status ./internal/cli ./internal/ensigncycle -count=1` must prove boot/status projection and command-level lifecycle behavior still consume the shared reducer without reintroducing a selector. Golden and entity-byte assertions must fail if either removed field is emitted, if stale input is accepted, or if a rejected candidate projects approval.

Repository validation is `gofmt -w ./cmd ./internal`, then `go test ./...`, then `go test ./... -race`. The implementation report must record the focused commands and both full-suite commands with the behavior each would falsify; repository search over `internal`, `skills/fo-gate-lifecycle`, `docs/specs`, and `docs/schema` must show no normative or emitted `gates.current`, `digest-domain`, or `raw-file-pin`, while allowing explicitly labeled historical evidence outside those paths.

## Risk spike and alternatives

The riskiest mechanism was whether status plus immutable records can select correctly through same-stage history and re-entry. A throwaway `internal/gates` test constructed two `ideation` records: an older terminal-target approval selected by `gates.current`, and a newer rejected record. Current `Validate` accepted the duplicate-stage document and `CurrentStageReadiness` returned `approved-awaiting-merge` from the old record; `recordForStage(doc, "ideation")` rejected the ambiguity. `go test ./internal/gates -run TestSchemaIdeationSpikeReproducesStaleTerminalApproval -count=1 -v` passed on 2026-08-01, and the throwaway file was removed after the exercise. This reproduces the failure mechanism and establishes unique-stage validation plus status lookup as the test seed.

Retaining the pointer with repair/rebind logic is the simplest alternative considered, but cannot deliver AC-2: it preserves two mutable authorities and every writer/read path must synchronize them. Deriving by status without enforcing unique stages is also insufficient because the spike proves ambiguous records exist. A replacement record index is unnecessary because record counts are small and linear lookup already exists.

Removing only the `raw-file-pin` branch while retaining `digest-domain` is insufficient for AC-3: it leaves a mandatory field with one legal value and preserves domain-switch control flow. Fixing canonical semantics in v1 makes malformed alternative fields ordinary unknown-field failures.

## Expected surface, estimate, and semantic boundary

The direct field inventory fixes the minimum semantic surface. `gates.current` is produced by `internal/gates/prepare.go` and `internal/gates/operation.go`; it is validated and read for readiness in `internal/gates/model.go`, consumed for application/retained authority in `internal/gates/application.go`, and consumed for finalize/rework in `internal/gates/delivery.go`. `digest-domain` is produced by `internal/gates/prepare.go` and `internal/gates/round.go`; it is validated or compared in `internal/gates/model.go`, `internal/gates/operation.go`, `internal/gates/application.go`, and `internal/gates/round.go`. `internal/status/discover.go` and `internal/cli/cli.go` are transitive consumers through the unchanged gates APIs and are not expected to need production edits.

Expected implementation surface is **20–22 files**: those six `internal/gates` production files; the five paired gate tests plus `internal/gates/testdata/advisory-round/recorded-entity.md`; command/status/integration fixtures in `internal/cli/{gate,terminal_consume}_test.go`, `internal/status/{gates_coexist,boot_identify}_test.go`, and `internal/ensigncycle/{gate_assert_impl,gate_assert,shared_round_recording}_test.go`; and the three normative consumers `skills/fo-gate-lifecycle/SKILL.md`, `docs/specs/gate-resolution-frontmatter-contract.md`, and `docs/schema/entity.mdschema.yml`. Estimate: **90–150 inserted lines** for the adversarial table and revised assertions, **110–190 deleted lines** from selector/domain code and fixtures, with a net deletion of at least 20 lines.

Tolerance is **24 files, 180 insertions, and 230 deletions** to allow a discovered direct fixture/helper consumer or formatting split; the net delta must remain negative. Any production edit outside the six named gate files, any edit to the transitive `internal/status/discover.go` or `internal/cli/cli.go` consumers, any compatibility reader/migration, or exceeding a tolerance requires gate re-review.

Declared semantic changes: stored gate frontmatter removes `current` and Briefing `digest-domain`; duplicate stage records become invalid; runtime gate selection derives from entity status and false stale approvals disappear. Command grammar, stdout/stderr vocabulary except ambiguity diagnostics, Resolution/application authority, IDs, room layout, locking, atomicity, and workflow status transitions do not change. No compatibility reader, migration command, or historical-roadmap rewrite is in scope.

## Required documentation diff

Apply this wording in `docs/specs/gate-resolution-frontmatter-contract.md` and mirror the invariant in `docs/schema/entity.mdschema.yml`:

```diff
 gates:
   version: 1
-  current:
-    gate: gate:example:sample:validation
   records:
@@
             id: briefing:sample-validation-1a
             digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
-            digest-domain: canonical-bytes
             request-digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
@@
-`gates.current.gate` selects the logical gate eligible for later application.
+Entity `status` selects exactly one record by its `stage`; duplicate stage records fail closed.
+The last ordered attempt in that record is current and eligible for later application.
@@
-Every Briefing binding includes an id, SHA-256 digest, the `canonical-bytes` digest
-domain, and an exact file or room reference.
+Every Briefing binding includes an id, canonical SHA-256 digest, and an exact file
+or room reference. Version 1 digests are unconditionally RFC 8785/JCS canonical bytes.
```

## Stage Report: ideation

- DONE: Produce a concrete, bounded problem statement, proposed approach, common scenarios, and expected surface for simplifying the stable-v1 gate-state schema.
  The body fixes selection to unique status-matched records, removes fixed digest metadata, and bounds the expected surface at 12–18 files with explicit tolerance.
- DONE: Define independently falsifiable end-state acceptance criteria with a reproducible test plan, declared semantic scope, tolerance, and required documentation diff.
  AC-1–AC-5 name falsifiers; AC-2 measures false approvals at 0/2 versus the current 1/2 baseline; commands, format changes, limits, and exact doc wording are declared.
- DONE: Exercise or explicitly record the riskiest unverified mechanism and explain why the simplest alternative cannot deliver the value.
  A throwaway Go test reproduced stale approval through duplicate same-stage records and proved the existing stage lookup fails closed; pointer repair and field-only trimming alternatives remain insufficient.

### Summary

Ideation defines a clean unreleased v1 with no stored current-gate pointer and no one-value digest-domain field. The spike identified unique record-per-stage validation as the necessary condition for safe status derivation and seeded an adversarial regression test for the observed false-approval class.

## Stage Report: ideation (cycle 2)

- DONE: Repair the ideation contract with a distinct exact Acceptance criteria section and a separate Test plan section.
  The exact headings now expose AC-1 through AC-5 to `status --read jc --ac-scan`; the separate plan cites focused, full, and race commands and the failures they detect.
- DONE: Tighten the proposed schema surface to the minimum consumer inventory and honest file/LOC tolerance while preserving the stale-approval adversarial fixture.
  The inventory names each direct producer and consumer, bounds 20–22 expected files with a 24-file tolerance, and retains AC-2's 0-of-2 stale-approval measure plus AC-3's closed stored format.
- DONE: Rerun the mandatory AC scan and record the full/race validation plan and all required stage-report evidence.
  The corrected report explicitly evidences AC-1, AC-2, AC-3, AC-4, and AC-5; the independent scan is rerun after this report, and the test plan requires focused suites, `go test ./...`, and `go test ./... -race`.

### Summary

The ideation contract is now independently scannable and implementation-ready without changing its clean-v1 direction. The repaired package narrows the implementation boundary to observed producers, consumers, and fixtures while keeping the stale-approval adversary and full/race validation obligations intact.
