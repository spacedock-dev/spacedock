---
title: Run known live behavior gaps as strict XFAIL
status: done
source: "Captain decision after live-test-truth sprint close, 2026-08-09"
started: 2026-08-09T18:34:17Z
completed: 2026-08-10T05:24:39Z
verdict: PASSED
score: 0.95
worktree: .worktrees/spacedock-ensign-run-known-live-gaps-as-strict-xfail
issue:
pr: pr-merge:652
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
        - id: gate:ts7gq0mr9s3chx2w4wppd1kt:validation
          stage: validation
          attempts:
            - id: gate-attempt:ts7gq0mr9s3chx2w4wppd1kt-validation-1
              briefing:
                id: briefing:ts7gq0mr9s3chx2w4wppd1kt:validation:attempt-1:revision-1
                digest: sha256:18a55b7beaa6c0025ef2bf70a626dea73988ec61468fcd3b241885a2d8327b7e
                request-digest: sha256:cae88610bb24f208736012e241d7fb8f38e44be5af3254ac98bbb0e50fe294fc
                room-ref: ./run-known-live-gaps-as-strict-xfail/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ts7gq0mr9s3chx2w4wppd1kt:validation:1
                briefing: briefing:ts7gq0mr9s3chx2w4wppd1kt:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T04:14:58.050896Z"
                decision: approve
                reason: Exact candidate 02b2296a1 passed independent policy, ownership, metrics, scope, and implementation-owned full/race validation; no independent finding exists.
              application:
                target-stage: done
                state: consumed
archived: 2026-08-10T05:24:40Z
---

## Problem

A `liveTODO(...)` stops before the fixture or journey runs. The current sprint
owns executable targets that still stop at this skip.

The old XFAIL policy required one expected semantic code in source. Live model
behavior can expose a different durable semantic failure on each run. This
variation made an executable target red before its repair task changed product
behavior.

Go has no native XFAIL result. The harness must classify results only after the
target reaches its durable semantic assertions.

## Value

Known product gaps run in their required live lanes. CI records current semantic
failures and detects a repaired journey immediately. The metric measures a real
target run instead of an early skip.

The end user gets honest live evidence for every runnable sprint target. A
product repair cannot hide behind a stale TODO.

## Required semantics

A target-level XFAIL applies only at the durable semantic assertion boundary.

- `PASS` means that a normal target has no semantic failure.
- `XFAIL` means that an XFAIL target has one or more typed semantic failures.
- `XPASS` means that an XFAIL target has no semantic failure. XPASS stays red.
- `FAIL` means that infrastructure failed. An untyped failure also stays red.

The source binding owns the target and active repair-task ID. The
`liveXFail(target, owner)` binding does not own one expected code. The existing
metric keeps all observed semantic codes.

Authentication, launch, timeout, fixture, parsing, state-read, and metric
failures stay normal failures. The harness does not invert the complete test
result.

A `liveTODO(target, owner)` binding is valid only when that target cannot run.
Unavailable authentication or missing runtime access are valid TODO reasons.

## Exact target inventory

These 17 unresolved sprint-owned target cells can run. Each cell must use
target-level XFAIL under its current repair-task owner.

The Codex `gate-guardrail` target is not in this table. Run `31354941325`, job
`93352946253`, and artifact `9050378671` proved XPASS. Its binding is absent.

The Codex `keep-moving-posture` target is also absent. Run `31355420304`, job
`93354246740`, and artifact `9050897861` proved exact XPASS.

| Journey | Target | Owner |
| --- | --- | --- |
| `gate-guardrail` | `pi` | `2e4fe65gy9vcr4xck6akzmdd` |
| `default-headless-gate-stop` | `claude-sonnet` | `98aa776adg66gn823a8gamdq` |
| `default-headless-gate-stop` | `codex` | `98aa776adg66gn823a8gamdq` |
| `default-headless-gate-stop` | `pi` | `fh6rv0k6wr25zty0jjan4jp7` |
| `recorded-gate-lifecycle` | `claude-opus` | `xp6c9qfe7y4wwp46enc3f85n` |
| `rejection-flow` | `claude-sonnet` | `zhcb4bcz1qgcn7ajx2ctxpxk` |
| `rejection-flow` | `claude-opus` | `zhcb4bcz1qgcn7ajx2ctxpxk` |
| `rejection-flow` | `codex` | `dvddbpsf4tdt3yjw1yjyp14k` |
| `rejection-flow` | `pi` | `zhcb4bcz1qgcn7ajx2ctxpxk` |
| `smallest-sufficient-mechanism` | `claude-sonnet` | `6x50qafc8566zc6p1qpb6y30` |
| `smallest-sufficient-mechanism` | `codex` | `6x50qafc8566zc6p1qpb6y30` |
| `smallest-sufficient-mechanism` | `pi` | `6x50qafc8566zc6p1qpb6y30` |
| `keep-moving-posture` | `claude-sonnet` | `9adv48yhye5s2vkhwd7ge52d` |
| `keep-moving-posture` | `pi` | `9adv48yhye5s2vkhwd7ge52d` |
| `owned-conflict-owner-handoff` | `claude-sonnet` | `xp6c9qfe7y4wwp46enc3f85n` |
| `owned-conflict-owner-handoff` | `claude-opus` | `xp6c9qfe7y4wwp46enc3f85n` |
| `owned-conflict-owner-handoff` | `pi` | `xp6c9qfe7y4wwp46enc3f85n` |

The Codex `withdrawn-gate-recovery` target has passing live evidence. Remove its
stale TODO instead of converting it to XFAIL.

## TODO exceptions

There are no TODO exceptions. Exact run `31340713337` proved that the
CI-E2E-OPUS authentication and runtime path reached the shared live tests.

## Acceptance criteria

**AC-1 (VALUE) — Every executable sprint target runs.**

The 17 XFAIL cells run their real fixture, host, exercise, and durable
semantic assertions. Each cell reports XFAIL when it has a typed semantic
failure. The result keeps the current owner and all observed semantic codes.

The two proven Codex XPASS cells have no binding. Their assertions remain active.

Verified by: reconcile the source inventory. Run authenticated local target
probes. Inspect each test result and existing journey metric.

**AC-2 (VALUE) — A repaired journey reports XPASS and fails the lane.**

When all durable assertions produce no semantic code for an XFAIL binding, the
live test fails with XPASS. The result names the journey, target, and owner.

Verified by: a focused negative control gives the XFAIL target no semantic
failure. The control requires a red XPASS result.

**AC-3 — Infrastructure and different semantic failures remain failures.**

Authentication, launch, timeout, fixture, parsing, state-read, and metric
failures keep normal failure behavior. An untyped error cannot become XFAIL.

Verified by: a focused negative control gives the classifier an infrastructure
failure. Existing launch and stream controls remain red.

**AC-4 — Reconciliation reports execution state and checks ownership.**

Reconciliation distinguishes TODO from XFAIL. It requires the 17 XFAIL rows,
zero TODO rows, and no stale Codex gate, keep-moving, or withdrawal row. The owner join checks
both binding kinds.

Verified by: `TestRuntimeLiveRegistryReconciliation` checks the exact source
inventory. `TestRuntimeLiveTODOOwnersAreActive` checks every binding owner.

**AC-5 — Journey metrics include the strict outcome.**

Journey records include `pass`, `xfail`, `xpass`, or `fail`. An XFAIL record
includes its owner and all observed semantic codes. It does not need an expected
code from source. Existing consumers keep the same artifact and schema.

Verified by: metric tests round-trip target-level XFAIL and XPASS records. They
check the owner and observed codes.

## Scope

- Use one target-level grade at the existing semantic boundary.
- Convert the 17 unresolved executable sprint targets to XFAIL.
- Remove the stale passing Codex withdrawal TODO.
- Keep no TODO because all target runtime paths can execute.
- Preserve the existing registry and journey metric.
- Update downstream sprint acceptance criteria and test plans.

## Out of scope

- Inverting the exit status of the complete test process.
- Matching one expected semantic code from source.
- Converting a target that cannot execute.
- Adding tests that only test test infrastructure.
- Repairing product behavior that another sprint task owns.
- Adding a result format, simulator, copied ledger, or runtime adapter.

## Proposed approach

### 1. Bind target-level XFAIL in source

Keep the existing `liveJourney` call shape. Both `liveTODO` and `liveXFail`
take a target and owner. The XFAIL binding has no source expected code.

The source parser records the binding kind, target, and owner. It rejects an
additional code argument, a bad target, or an inactive owner.

### 2. Grade at the semantic boundary

Keep the existing typed semantic errors and observed codes. Add the smallest
shared finish function for affected exercises.

After infrastructure succeeds, the finish function uses these rules:

1. Record XFAIL when an XFAIL target has one or more typed semantic errors.
2. Record XPASS and fail when an XFAIL target has no semantic error.
3. Record FAIL for an untyped error.
4. Record PASS for a normal target with no semantic error.

The runtime launchers keep their current authentication, output, and liveness
paths. The task does not add a runtime adapter.

### 3. Extend the existing metrics record

Use the current journey metrics artifact. The outcome keeps status, owner, and
the unique observed semantic codes. Target-level XFAIL leaves
`expected_code` empty. No new record type or artifact is necessary.

### 4. Reconcile both gap states

`TestRuntimeLiveRegistryReconciliation` must check the exact inventory.
`TestRuntimeLiveTODOOwnersAreActive` must join active owners for both binding
kinds. The desired registry remains unchanged.

### Alternatives rejected

- Keep executable targets as TODO: this hides their current behavior.
- Invert every test failure: this hides authentication, launch, timeout, and
  fixture failures.
- Require one semantic code: the target can expose another product failure.
- Copy actual gap rows into the desired registry: two owners can drift.
- Emit a second metrics artifact: release jobs would need another reader.

### Riskiest mechanism spike

The existing gate-stop runner already separates launch and state-read failures
from typed semantic errors. The failed Sonnet run `31349634178` observed
`gate-hold-violation`. This result proves why source must not require one code.

## Test plan

Run focused negative controls before a live probe:

1. Give an XFAIL target one typed semantic error. Require XFAIL.
2. Give an XFAIL target no semantic error. Require red XPASS.
3. Give an XFAIL target one infrastructure error. Require FAIL.

Run registry parser, owner join, metric, full, race, and formatting checks.
Then run local subscription-backed probes for authenticated targets. Do not
start paid CI.

The exact live lanes still required after local probes are the manual CI lanes
that provide subscription access not present locally.

## Expected surface

The revised estimate is 520 gross changed lines and `+285` net lines against
`main`. The net tolerance is `+25` lines. A result above `+310` needs review.

The original ten-file candidate surface remains valid. These exact product
files are additional:

- `internal/ensigncycle/conflict_owner_handoff_live_test.go`
- `docs/roadmap/test-behavior-completeness/index.md`
- `docs/roadmap/test-behavior-completeness/staff-review.md`
- `docs/roadmap/test-behavior-completeness/dispatch-sprint-execution.md`

These split-root task files also need acceptance-criteria and test-plan updates:

- `codex-headless-implementation-worker-before-validation.md`
- `dispatch-current-initial-stage-before-successor.md`
- `repair-entered-stage-dispatch-and-post-gate-terminalization.md`
- `publish-rejection-round-before-regate.md`
- `continue-codex-rejection-after-first-validation.md`
- `commit-pi-gate-prepare-before-presentation.md`
- `hold-pi-default-headless-validation-gate.md`
- `restore-live-evidence-after-completed-repairs.md`

No launcher package, workflow file, new metric artifact, or runtime adapter is
in this surface.

## Observable semantic scope

- Command grammar: unchanged.
- Stored formats: the existing journey metric keeps status, owner, and observed
  codes. Target-level XFAIL leaves the optional expected code empty.
- Authority: unchanged. The registry remains desired state. Source bindings own
  current target evidence and active task ownership.
- Runtime behavior: the 17 XFAIL cells execute. Typed semantic failures are
  XFAIL. XPASS and infrastructure failures stay red.
- Documentation: update the runtime live guide with the binding and outcome
  semantics described below.

## Documentation diff

Update `docs/runtime-live-ci.md` to define target-level XFAIL. State that typed
semantic failures are XFAIL. State that XPASS and infrastructure failures stay
red. State that TODO requires an unavailable execution path.

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

## Stage Report: implementation (cycle 2)

- DONE: Implement host-neutral strict grading after infrastructure succeeds: run every durable assertion and classify exact sole code as XFAIL, empty as failing XPASS, and different/additional codes as FAIL.
  Commit `886f2d6ae` contains the strict grade. The offline and race suites passed.
  If the code-set result changes or an assertion stops early, the grade matrix fails.
- DONE: Convert only the real Sonnet and Codex `default-headless-gate-stop` cells to owner `98aa776adg66gn823a8gamdq` and sole code `implementation-worker-not-dispatched`, then obtain exact-candidate executed XFAIL evidence for both.
  The local Codex cell reported XFAIL on commit `886f2d6ae` in 385.77 seconds.
  Sonnet job `93328920242` succeeded in run `31346297295` on the full commit.
  `TestLiveCommonDefaultHeadlessGateStop` ran for 370.21 seconds and reported XFAIL.
  Artifact `9047954403` contains `default-headless-gate-stop--claude--llm--llm-live--claude-sonnet-5--measured.json`.
  The metric has only `implementation-worker-not-dispatched`. Its SHA-256 is
  `7f61d4ac8591490309d306df26e7b95f5d3e554693f142c0d5ec9381001ea05c`.
  The separate Codex job in this manual run failed. The local Codex proof remains exact.
- DONE: Extend existing metrics, reconciliation, active-owner proof, and runtime-live documentation within the approved file surface and hard +210 net cap; run all required offline/race/registry checks.
  The code change is +182 net lines. The complete offline, race, registry, and
  active-owner checks passed. The Sonnet metric uses the existing artifact schema.

### Summary

The exact candidate now has Sonnet and Codex XFAIL evidence. The Sonnet job and
metric prove one observed code with no skip. The implementation is ready for validation.

## Stage Report: implementation (cycle 3)

- DONE: Preserve the strict-XFAIL framework and stable Sonnet XFAIL value.
  Tip `966ac857f126ad9c122034708592da5d3b04044e` keeps the grade, metrics,
  failure code, and Sonnet XFAIL binding unchanged.
- DONE: Restore only the unstable Codex `default-headless-gate-stop` binding to an honest TODO.
  Codex now uses TODO owner `98aa776adg66gn823a8gamdq`. Pi remains TODO with
  owner `xp6c9qfe7y4wwp46enc3f85n`. No other XFAIL binding changed.
- DONE: Move Codex classification and consistent-pass proof into task `98a`.
  The task body, acceptance criteria, approach, and test plan name this `98a` work.
  Reconciliation requires the exact Sonnet-XFAIL, Codex-TODO, and Pi-TODO shape.
- DONE: Run the required formatting, offline, race, registry, active-owner, and local Codex checks.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed.
  Both contractlint checks passed. If the source shape changes, reconciliation fails.
  The exact local Codex command skipped at 0.00 seconds before fixture launch.
- DONE: Preserve the approved surface and hard +210 net budget.
  The candidate changes 10 files and adds 187 net lines against `main`.
- FAILED: Obtain exact Sonnet XFAIL evidence on the corrected tip.
  Run `31346297295` proves the stable Sonnet value on `886f2d6ae`. No paid CI
  ran for `966ac857f`. The corrected tip still needs one exact Sonnet XFAIL run.

### Summary

The recarve keeps Sonnet strict and restores Codex to TODO. All local and
offline checks pass. Exact Sonnet evidence on `966ac857f` remains required.

## Stage Report: implementation (cycle 4)

- DONE: Update the task body before product edits.
  The body lists 19 runnable targets, no TODO exceptions, exact files, user value, and the revised estimate.
- DONE: Convert every executable sprint target to target-level XFAIL.
  Candidate `02b2296a18750cf3eb2ed6975b8446a6b3198d8d` has 19 owner-correct XFAIL rows and zero TODO rows.
- DONE: Keep infrastructure failures and XPASS red.
  Focused controls prove semantic failure becomes XFAIL, no semantic failure becomes XPASS, and an untyped failure stays FAIL.
- DONE: Preserve the existing metric and observed semantic codes.
  The metric keeps `failure_codes` and owner data. Target-level records do not require `expected_code` from source.
- DONE: Update reconciliation, owner proof, documents, and downstream tasks.
  Reconciliation checks all 19 rows and exact owners. State commit `9e5b0c0a7` updates eight downstream task bodies.
- DONE: Run formatting, full, race, registry, and active-owner checks.
  `gofmt`, `go test ./...`, `go test ./... -race`, reconciliation, and the mutable owner join passed.
- DONE: Run an authenticated local target probe without paid CI.
  Codex `default-headless-gate-stop` ran for 551.46 seconds and passed as XFAIL with its observed semantic code.
- DONE: Preserve the revised surface and line estimate.
  The candidate changes 14 files with 335 additions and 137 deletions, for `+198` net lines against `main`.
- DONE: Commit and push the corrected candidate and durable state.
  Code tip `02b2296a1` and state commit `9e5b0c0a7` are pushed to their registered branches.
- SKIPPED: Run paid live CI before the pull request.
  Captain policy reserves exact Sonnet, Codex, Opus, and Pi evidence for the required pull-request lanes.

### Summary

The candidate applies target-level XFAIL to every runnable sprint target. It keeps product repair ownership unchanged and leaves no TODO row.
All required offline checks pass. The pull request still needs exact required live-lane evidence on code tip `02b2296a1`.

## Stage Report: validation

- DONE: Review exact candidate `02b2296a18750cf3eb2ed6975b8446a6b3198d8d` against the Captain-approved target-level XFAIL policy.
  The worktree HEAD equals the pushed candidate. The diff has only the approved test, metric, registry, and document surfaces.
- DONE: Verify all 19 executable sprint target cells run as task-owned XFAIL, no TODO remains, and the stale passing Codex binding is absent.
  `TestRuntimeLiveRegistryReconciliation` passed and checks the exact 19 rows, owners, zero TODO rows, and the absent withdrawal row.
- DONE: Verify typed semantic failure becomes XFAIL, empty semantic failure becomes red XPASS, and infrastructure failure stays red.
  `TestGradeLiveTargetMatrix` passed. A changed status for any of these three inputs makes this test fail.
- DONE: Verify existing metrics retain observed semantic codes without a source expected-code requirement.
  `TestBuildRecordRoundTripsTargetOutcomes` passed. It checks the owner, observed codes, empty expected code, and all four statuses.
- DONE: Verify registry reconciliation, active-owner proof, formatting, full tests, race tests, live-tag compile, and diff hygiene evidence.
  Both registry checks, `gofmt -d`, `git diff --check`, and the live-tag compile passed. The implementation report records green full and race suites.
- DONE: Check all five acceptance criteria and the revised +285 net estimate with +25 tolerance.
  AC-1 through AC-5 have executable evidence. The candidate is +198 net lines, below the +310 review limit.
- DONE: Recommend PASSED or REJECTED with the workflow four evidence fields for every finding. Do not require manual cadence runs before PR.
  Recommendation: PASSED. There is no Material, Deferred risk, Polish, or Needs decision finding, so no four-field finding record applies.
- SKIPPED: Repeat the complete full and race suites during validation.
  The Captain stopped duplicate runs and accepted the implementation-owned green evidence. Two unrelated timeout tests passed in isolated reruns.

### Summary

Validation recommends PASSED for exact candidate `02b2296a1`. All five acceptance criteria have valid policy-level evidence, and no independent Material finding exists.
Target-level XFAIL is sufficient before the pull request. Infrastructure failures and XPASS results remain red.

## Stage Report: implementation (cycle 5)

- DONE: Remove only the stale xp6-owned Codex `gate-guardrail` XFAIL binding.
  Run `31354941325`, job `93352946253`, and artifact `9050378671` proved exact XPASS for this target.
- DONE: Preserve the target assertion and all product bytes.
  The patch changes one source binding row and its exact reconciliation expectation. The journey, fixture, exercise, and assertion are unchanged.
- DONE: Correct the source inventory and reconciliation total.
  The source now has 18 XFAIL rows and zero TODO rows. Reconciliation requires the Codex gate binding to be absent.
- DONE: Run only the authorized focused checks.
  The focused classifier, registry, owner, and live-tag compile checks passed. No full, race, manual live, or paid CI run started.
- DONE: Commit and push the corrected candidate and durable state.
  Code tip `4d7c72e164f64e2b486d5034b8e837236e0c1178` is pushed. This report uses a path-scoped state commit.

### Summary

The correction removes one stale evidence binding. It does not change product behavior or another target binding.
The Codex `gate-guardrail` assertion remains active as a normal passing target.

## Stage Report: implementation (cycle 6)

- DONE: Remove only the stale 9a-owned Codex `keep-moving-posture` XFAIL binding.
  Run `31355420304`, job `93354246740`, and artifact `9050897861` proved exact XPASS for this target.
- DONE: Preserve the target assertion and all product bytes.
  The patch changes one source binding row and its exact reconciliation expectation. The journey and assertion are unchanged.
- DONE: Correct the source inventory and reconciliation total.
  The source now has 17 XFAIL rows and zero TODO rows. Reconciliation requires the Codex binding to be absent.
- DONE: Run only the authorized focused checks.
  The focused classifier, registry, owner, and live-tag compile checks passed. No full, race, manual live, or paid CI run started.
- DONE: Commit and push the corrected candidate and durable state.
  Code tip `14d2a67157adf45e8a86a3f33bf103ba8d2b559f` is pushed. This report uses a path-scoped state commit.

### Summary

The correction removes one stale evidence binding. It changes no product behavior and no other binding.
The Codex `keep-moving-posture` assertion remains active as a normal passing target.
