---
title: Run known live behavior gaps as strict XFAIL
status: implementation
source: "Captain decision after live-test-truth sprint close, 2026-08-09"
started: 2026-08-09T18:34:17Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-run-known-live-gaps-as-strict-xfail
issue:
pr:
mod-block:
sprint: test-behavior-completeness
group: common-evidence
sprint-readiness: ready
id: ts7gq0mr9s3chx2w4wppd1kt
gates:
    version: 1
    records:
        - id: gate:ts7gq0mr9s3chx2w4wppd1kt:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ts7gq0mr9s3chx2w4wppd1kt-backlog-1
              briefing:
                id: briefing:ts7gq0mr9s3chx2w4wppd1kt:backlog:attempt-1:revision-1
                digest: sha256:8dd057d5494ec56568e793d8cc503cc1dee58720ed236092096811e43ab4f9de
                request-digest: sha256:3ddbb008cf288eaa9aca4ede97e3d06d670684c6a5c4de66356884895bb4f4dd
                room-ref: ./run-known-live-gaps-as-strict-xfail/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ts7gq0mr9s3chx2w4wppd1kt:backlog:1
                briefing: briefing:ts7gq0mr9s3chx2w4wppd1kt:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:02.472302Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; the seed defines strict semantic outcomes and requires real Sonnet and Codex XFAIL cells.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:ts7gq0mr9s3chx2w4wppd1kt:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ts7gq0mr9s3chx2w4wppd1kt-ideation-1
              briefing:
                id: briefing:ts7gq0mr9s3chx2w4wppd1kt:ideation:attempt-1:revision-1
                digest: sha256:498f410c939dc1deab7eaf7e25eb83f57b0e0cc167bf6a8f5041adc0656d5fad
                request-digest: sha256:ec78e971d896eb5f13706bb5dd404c7dc6f3364989dec405e90909e7ba44f526
                room-ref: ./run-known-live-gaps-as-strict-xfail/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ts7gq0mr9s3chx2w4wppd1kt:ideation:1
                briefing: briefing:ts7gq0mr9s3chx2w4wppd1kt:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T21:33:07.558424Z"
                decision: approve
                reason: Captain approved the strict-XFAIL direction and Commander activation.
              application:
                target-stage: implementation
                state: consumed
---

## Problem

A `liveTODO(...)` stops before the fixture or journey runs. It records ownership,
but it gives no current behavior evidence. It also cannot report that a product
repair made the journey pass.

The current `default-headless-gate-stop` cells for Sonnet and Codex have a real
product failure. Their headless runs prepare the validation gate without first
dispatching the implementation worker. The live evidence labels this failure
as `implementation-worker-not-dispatched` for candidate
`7c8708c8537fc73761e56813ddbd6a498959ef19`.

Go has no native XFAIL result. A broad inversion of the test process can hide
authentication, launch, timeout, fixture, and unrelated assertion failures.

## Value

Known product gaps run in their required live lanes. CI records the expected
failure, rejects a different failure, and detects a repaired journey immediately.
The live artifact then measures execution instead of counting an early skip.

## Required semantics

A strict XFAIL applies only at the durable semantic assertion boundary.

- `PASS` means that the journey passed.
- `XFAIL` means that the journey ran and the named semantic failure occurred.
- `XPASS` means that an expected failure disappeared. `XPASS` fails the lane until
  its source binding is removed.
- `FAIL` means that infrastructure failed or a different semantic failure
  occurred.

The source binding owns the target, active task ID, and stable failure code. A
`liveTODO(target, owner)` binding keeps the current skip behavior. A
`liveXFail(target, owner, code)` binding runs the fixture and expects `code`.
The desired-state registry does not copy this actual-state information.

Keep TODO only when the journey cannot run. Examples include a missing adapter,
fixture, selector, or authentication path.

Infrastructure failures stay outside the grade result. Fixture setup, process
start, authentication, timeout, output parsing, and state reads must succeed
before semantic assertions start. Then every durable semantic assertion runs,
even after another assertion returns a semantic error. Metric emission remains a
hard failure after grading.

The runner collects the unique stable codes from all durable semantic
assertions. XFAIL is valid only when this set is exactly `{expected}`. An empty
set for an XFAIL binding is XPASS. The expected code plus any additional code,
or any different code, is FAIL. The runner does not stop at the first semantic
failure. Only typed semantic grade errors reach this classification.

## Acceptance criteria

**AC-1 (VALUE) — Sonnet and Codex execute the known headless gap.**

The Sonnet and Codex `default-headless-gate-stop` cells run the real fixture,
host, exercise, and every durable assertion. On the current failing candidate,
each cell reports `xfail` only when the unique code set is
`{implementation-worker-not-dispatched}` for owner
`98aa776adg66gn823a8gamdq`. Neither cell reports a skip.

Verified by: run both exact-candidate live commands and inspect the two existing
journey metric records plus the test result. The independent baseline is the two
current TODO skips.

**AC-2 (VALUE) — A repaired journey reports XPASS and fails the lane.**

When all durable assertions produce no semantic code for an XFAIL binding, the
live test fails with `XPASS`. XPASS means no semantic failure, not only that the
first assertion passed. The result names `default-headless-gate-stop`, target,
owner, and expected code. The source binding remains until a later change
removes it.

Verified by: focused classifier tests feed no semantic codes to each expected
binding and assert `xpass` plus a failing disposition. A later exact-candidate
run after task `98a` repairs the behavior and proves the live lane fails.

**AC-3 — Infrastructure and different semantic failures remain failures.**

Authentication, launch, timeout, fixture, state-read, and artifact failures keep
their existing hard-failure behavior. A semantic grade with a different code,
or with the expected code plus another code, also fails and never becomes
XFAIL. The classifier must inspect all durable assertion results before it
chooses this outcome.

Verified by: focused grade-matrix tests cover the matching code, no code, a
different code, and expected plus additional code. Existing runner failure tests
continue to exercise launch and stream failure paths. No test-only failure
inversion is added.

**AC-4 — Reconciliation reports execution state and checks ownership.**

Reconciliation distinguishes TODO from XFAIL. It validates target, active owner,
and expected code for each XFAIL. The mutable owner join checks both binding kinds.

Verified by: `TestRuntimeLiveRegistryReconciliation` parses one TODO and one
XFAIL source binding, then rejects mutations that change kind, owner, target, or
code. `TestRuntimeLiveTODOOwnersAreActive` runs with the state checkout and checks
both rows.

**AC-5 — Journey metrics include the strict outcome.**

The existing Claude and Codex journey records include `pass`, `xfail`, `xpass`,
or `fail`. An XFAIL record includes its owner and its sole observed failure
code. An XPASS record includes its owner and expected code as binding metadata,
but no observed semantic failure. A FAIL record includes every observed
semantic code. Existing metric consumers continue to read the same artifact
directory and schema shape.

Verified by: journeymetrics serialization tests build all four strict outcomes,
round-trip the JSON, and assert owner, expected-code, and observed-code fields.
Emitter tests use one Claude and one Codex fixture record. A missing outcome or
code set fails.

## Scope

- Add the smallest grade result needed at the semantic assertion boundary.
- Preserve immediate failures for infrastructure and setup errors.
- Convert the two evidenced `98a` cells in the same change.
- Keep the registry as desired state.
- Use the existing journey metrics artifact.

## Out of scope

- Inverting the exit status of the complete test process.
- Matching unstable CLI text or complete error strings.
- Converting an unexecuted TODO without stable failure evidence.
- Adding tests that only test test infrastructure.
- Repairing the product behavior owned by `98a`.

## Proposed approach

### 1. Bind the expected failure in source

Keep the existing `liveJourney` call shape and replace its target gap list with
one list of typed gap bindings. `liveTODO` keeps two arguments. `liveXFail` adds
the stable semantic code. The first landing changes only this call:

```go
[]liveJourneyGap{
    liveXFail("claude-sonnet", "98aa776adg66gn823a8gamdq", "implementation-worker-not-dispatched"),
    liveXFail("codex", "98aa776adg66gn823a8gamdq", "implementation-worker-not-dispatched"),
    liveTODO("pi", "xp6c9qfe7y4wwp46enc3f85n"),
}
```

The source parser must record the binding kind, target, owner, and code. It must
reject a missing code on `liveXFail`, an unexpected code on `liveTODO`, a bad
target, or an inactive owner. No copied gap table belongs in the registry.

### 2. Grade at the semantic boundary

Enrich the existing `gradedErr` with a stable code. The gate-hold oracle assigns
`implementation-worker-not-dispatched` only to its committed no-authority branch.
Other semantic branches use different codes. The classifier compares codes, not
the unstable prose after the code.

The shared scenario carries a small grade state. After setup, launch, timeout,
fixture, and state reads succeed, `runGateStopScenario` invokes every durable
assertion. It appends each typed semantic error to the state and does not return
after the first one. Metric errors remain hard test failures. The shared runner
resolves the unique observed code set after the exercise:

1. Record `XFAIL` only when the set is exactly `{expected}`.
2. Record `XPASS` and fail the test when an expected binding has an empty set.
3. Record `FAIL` and fail the test when the set has an additional or different
   code. This includes expected plus another code.
4. Record `PASS` when no expected binding exists and all grades pass.

The classifier is host-neutral. Claude and Codex continue to use their current
launchers, authentication, output, and liveness paths.

### 3. Extend the existing metrics record

Use the current journey metrics artifact. Strict live records use
`outcome.status` values `pass`, `xfail`, `xpass`, and `fail`. The outcome carries
the source `owner`, the `expected_code` when a strict binding exists, and the
unique observed `failure_codes` list. XFAIL has one observed code and that code
is the expected code. XPASS has an empty observed-code list, because it means
no semantic failure. FAIL preserves every observed code, including the expected
code when an additional code also occurs.

Add these optional fields to the existing outcome object. Do not create a second
artifact or a copied gap ledger. Preserve the current `passed` and `failed`
values for older non-strict records unless the metrics package already provides
a compatible normalization point.

### 4. Reconcile both gap states

`TestRuntimeLiveRegistryReconciliation` must distinguish TODO and XFAIL in its
derived execution report. `TestRuntimeLiveTODOOwnersAreActive` must join active
owners for both binding kinds. The desired registry remains unchanged because it
describes required journeys, not current evidence.

### Alternatives rejected

- Keep TODO and count the gap: this hides the current journey and cannot detect
  a repaired behavior.
- Invert every test failure: this hides authentication, launch, timeout, and
  fixture failures.
- Match the full assertion text: model and harness text can change without a
  semantic change.
- Copy actual gap rows into the desired registry: two owners can drift.
- Emit a second metrics artifact: release jobs would need another reader.

### Riskiest mechanism spike

The existing semantic boundary was exercised first:

```text
go test ./internal/ensigncycle -run 'TestAssertRecordedGateHoldLog|TestAssertGateHeld' -count=1 -v
```

The targeted tests passed. `assertRecordedGateHoldLog` already returns a typed
`*gradedErr` for the exact missing-worker branch. Its setup, commit, decision,
withdrawal, status, and successor branches remain separate semantic errors. The
runner calls `t.Fatalf` for launch and stall errors before this boundary. This
spike supports enriching the existing error with a code. It does not support a
test-process inversion.

The exact candidate and live evidence remain the independent baseline. Sonnet
ran for 414.17 seconds and Codex ran for 166.56 seconds before both reported the
missing-worker failure. The first landing must rerun both real cells and record
two XFAIL artifacts before task `98a` changes First Officer behavior.

## Test plan

Run the pure grade matrix before any live spend. It must invoke every durable
assertion after infrastructure succeeds and prove these outcomes: exactly the
expected code is XFAIL, no code is XPASS, a different code is FAIL, and the
expected code plus an additional code is FAIL. It must also prove that a later
assertion runs after an earlier assertion returns a semantic error. A normal
binding with no semantic code is PASS. Run the registry parser and owner join
tests with source mutations. Run the metrics round-trip tests against the
existing `Record` type. Do not add a second JSON format.

Then run the real Sonnet and Codex default-headless cells on the exact candidate:

```bash
SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet \
  go test -tags live -count=1 -timeout 40m \
  -run '^TestLiveCommonDefaultHeadlessGateStop$' ./internal/ensigncycle -v

SPACEDOCK_LIVE_RUNTIME=codex \
  go test -tags live -count=1 -timeout 40m \
  -run '^TestLiveCommonDefaultHeadlessGateStop$' ./internal/ensigncycle -v
```

The expected first-landing result is one executed XFAIL per lane with the sole
code `implementation-worker-not-dispatched`. A skipped test, an infrastructure
error classified as XFAIL, or a different or additional semantic code fails the
proof. After task `98a` repairs the product behavior, rerun both commands and
expect XPASS failure with no observed semantic code before removing each source
binding.

Estimated complexity is moderate. The work needs fixture-backed live tests for
the two real lanes, Go unit tests for classifier and metrics behavior, and a
registry parser test. It needs no new fixture, CLI command, host adapter, or
product behavior change.

## Expected surface

The implementation may touch only the following files. The estimates include
tests. The aggregate net limit is `+210` lines. Each file has a tolerance of
`±20%`. A larger change needs a new design review.

| File | Purpose | Gross lines | Net delta |
|---|---|---:|---:|
| `internal/ensigncycle/shared_live_runner_test.go` | Gap bindings and host-neutral grade disposition | 70 | +40 |
| `internal/ensigncycle/claude_runtime_helpers_test.go` | Stable semantic failure code at the existing oracle | 25 | +15 |
| `internal/ensigncycle/claude_live_runner_test.go` | Record semantic assertion results without catching infrastructure errors | 25 | +5 |
| `internal/ensigncycle/journey_metrics_live_test.go` | Emit strict outcomes through the existing artifact | 25 | +18 |
| `internal/journeymetrics/types.go` | Add outcome owner, expected-code, and observed-code fields | 15 | +10 |
| `internal/journeymetrics/record.go` | Preserve legacy mapping and add strict mapping | 20 | +12 |
| `internal/journeymetrics/tracking_test.go` | Round-trip strict outcomes and legacy compatibility | 55 | +35 |
| `internal/contractlint/live_registry_reconciliation_test.go` | Parse TODO/XFAIL and join both owner states | 55 | +40 |
| `internal/ensigncycle/live_grade_unit_test.go` | Exercise the pure disposition matrix | 45 | +35 |
| **Total estimate** |  | **335** | **+210** |

The total row is the hard aggregate limit. Shared helpers and compact tests must
keep the implementation at or below `+210` net lines. No First Officer or
launcher product package, skill, workflow, desired registry, or new metrics
artifact belongs in the surface.

## Observable semantic scope

- Command grammar: unchanged.
- Stored formats: the existing journey metrics `outcome` object gains optional
  strict status, owner, expected-code, and observed-code values. No second
  artifact is added.
- Authority: unchanged. The registry remains desired state. Source bindings own
  current target evidence and active task ownership.
- Runtime behavior: the two named live cells execute and classify expected
  semantic failures. XPASS and unexpected failures fail their lanes. The First
  Officer product behavior is unchanged.
- Documentation: update the runtime live guide with the binding and outcome
  semantics described below.

## Documentation diff

Apply these replacements to `docs/runtime-live-ci.md`.

Before:

```text
Each declaration has an adjacent `liveJourney(...)` call that binds its stable
journey ID, fixture builder, target-scoped TODO owner, runtime-neutral exercise,
and durable assertion. There is no scenario table or runtime runner registry.
```

After:

```text
Each declaration has an adjacent `liveJourney(...)` call that binds its stable
journey ID, fixture builder, target-scoped TODO or strict-XFAIL owner, runtime-
neutral exercise, and durable assertion. A TODO skips only when the journey
cannot run. An XFAIL runs the journey and names its expected semantic code.
There is no scenario table or runtime runner registry.
```

After the existing mutable-owner paragraph, add:

```text
Strict live records use `pass`, `xfail`, `xpass`, or `fail`. After infrastructure
succeeds, the grade runs every durable semantic assertion. XFAIL requires the
sole observed code to equal the expected code. An empty code set is XPASS and
fails the lane until the source binding is removed. An additional or different
code is FAIL. Infrastructure failures remain ordinary test failures.
```

## Stage Report: ideation

- DONE: Define strict XFAIL at the semantic assertion boundary and keep infrastructure failures normal.
  The spike passed `TestAssertRecordedGateHoldLog` and `TestAssertGateHeld`. The design tags only typed semantic errors and keeps runner failures hard.
- DONE: Convert the real Sonnet and Codex default-headless-gate-stop cells in the first landing.
  The first binding replaces only the Sonnet and Codex TODO rows with `implementation-worker-not-dispatched` XFAIL rows. Pi remains TODO.
- DONE: Define source ownership, metrics outcomes, XPASS failure, and net line limits.
  The body defines source-owned target/owner/code, strict metrics statuses, XPASS lane failure, reconciliation, documentation wording, and a `+210` net cap.

### Summary

The design runs the two known product gaps and classifies only the named semantic
failure. It preserves hard infrastructure failures, the desired registry, and
the existing metrics artifact. The first landing requires two exact-candidate
XFAIL records before the product repair task changes First Officer behavior.

## Stage Report: ideation (cycle 2)

- DONE: Fold staff finding M1 into the semantic boundary.
  After infrastructure succeeds, every durable assertion runs. The unique code
  set must equal `{expected}` for XFAIL. Expected plus another code is FAIL.
- DONE: Define XPASS as no semantic failure.
  An empty code set for an XFAIL binding is XPASS and fails the lane. A different
  code also fails. The grade matrix must prove these cases without short-circuit.
- DONE: Keep the first landing tied to the real cells.
  Only the Sonnet and Codex `default-headless-gate-stop` cells use the known
  `98a` XFAIL code. Pi remains TODO. Metrics preserve the complete observed code
  set in the existing artifact, within the `+210` net line limit.

### Summary

The ideation record now closes M1. It requires full durable semantic assertion
coverage after infrastructure success, strict code-set classification, and two
real first-landing cells. No product change is included.

## Stage Report: implementation

- DONE: Implement host-neutral strict grading after infrastructure succeeds: run every durable assertion and classify exact sole code as XFAIL, empty as failing XPASS, and different/additional codes as FAIL.
  Commit `886f2d6ae` adds the code-set grade. If a strict result changes,
  `TestGradeLiveStrictMatrix` fails. If a later assertion does not run,
  `TestGradeLiveRunsEveryAssertion` fails.
- FAILED: Convert only the real Sonnet and Codex `default-headless-gate-stop` cells to owner `98aa776adg66gn823a8gamdq` and sole code `implementation-worker-not-dispatched`, then obtain exact-candidate executed XFAIL evidence for both.
  Commit `886f2d6ae` converts both cells and keeps Pi as TODO. The exact Codex
  cell reported XFAIL in 385.77 seconds. Its metric SHA-256 is
  `edded162c79f2433b02b7479bff9784430b8a494e1ca92771ca16a1ff4c403bf`.
  Sonnet did not run because neither `~/.claude/benchmark-token` nor
  `ANTHROPIC_API_KEY` exists. The harness skipped the test before the fixture launch.
- DONE: Extend existing metrics, reconciliation, active-owner proof, and runtime-live documentation within the approved file surface and hard +210 net cap; run all required offline/race/registry checks.
  The code change is +182 net lines. `go test ./...` and `go test ./... -race`
  passed. The two contractlint checks passed. A malformed binding or inactive
  owner makes these checks fail.

### Summary

The commit implements strict grading, source bindings, metrics, reconciliation,
and documentation. The Codex cell proves the exact XFAIL result. The Sonnet
evidence remains blocked by missing local authentication, and no paid CI ran.
