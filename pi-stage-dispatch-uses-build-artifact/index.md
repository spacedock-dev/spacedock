---
title: Pi stage dispatches should use dispatch build artifacts, not hand-rolled prompts
status: validation
source: captain (2026-06-04) — FO manually composed Pi subagent task prompts for fc/d2 instead of routing the canonical spacedock dispatch build artifact that carries entity slug/stage context
score: "0.30"
started: 2026-06-04T08:05:59Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact
issue:
id: z68h8vwxeetp011b1484c2jx
mod-block: merge:pr-merge
pr: "#300"
---

Spacedock stage dispatches should be driven by the canonical `spacedock dispatch build` artifact. That artifact is responsible for carrying the entity slug, entity path, workflow directory, target stage, stage definition, worktree path, checklist, and host/runtime constraints. During recent Pi work, the FO manually composed `subagent(...)` task prompts for `launcher-binary-path-passthrough` and `pi-stage-dispatch-fresh-context`, adding slug/stage details by hand. That was sufficient for the moment but bypassed the dispatch-builder contract and makes omissions likely.

This task was manually described because it is about formalizing the dispatch-build requirement itself. Future Pi stage dispatches should be routine FO work driven by the same assignment artifact as the other hosts, with only a small Pi wrapper around that artifact.

## Problem

Manual Pi subagent prompts can drift from the Spacedock stage contract. They may omit the entity slug, target stage, workflow directory, entity path, current stage definition, worktree layout, or completion checklist; they can also phrase stage-completion expectations differently from Claude/Codex dispatches. That weakens parity with the host-neutral dispatch builder and makes FO behavior harder to test.

The important separation is:

- `pi-stage-dispatch-fresh-context` controls the context boundary: a Pi stage worker starts from fresh context.
- This task controls the assignment source of truth: the Pi worker's assignment content comes from `spacedock dispatch build`, not from the FO hand-writing a new prompt.

## Proposed approach

Use the existing dispatch-builder host seam for Pi stage dispatches and make the FO runtime guidance plus tests prove that the FO wraps a canonical artifact instead of replacing it.

1. **Canonical assignment source.** For an initial Pi stage dispatch, the FO builds the assignment with `spacedock dispatch build --host pi` or an equivalent schema-version-2 request containing `host: "pi"`. The FO forwards the emitted dispatch file path or dispatch-file content as the worker assignment. The builder remains responsible for entity slug/name derivation, entity path, workflow directory, target stage, stage definition, worktree path, checklist, split-root state guidance, and completion-signal wording.
2. **Additive Pi wrapper.** The Pi `subagent(...)` wrapper may add transport/runtime fields that are not part of the canonical assignment: `context: "fresh"`, and optional human-facing `phase` / `label` values. Those values can make the Pi session legible, but they must not redefine the stage, slug, workflow, worktree, checklist, or completion contract that the builder emitted.
3. **No Pi acceptance contract.** Pi stage dispatches must not use a `subagent(... acceptance: ...)` contract. Spacedock's implementation-to-validation workflow owns acceptance through committed state/product changes and independent validation, so same-agent acceptance criteria in the Pi wrapper would duplicate or confuse the real gate.
4. **Preserve artifact fields through the wrapper.** Tests should cover that wrapping a builder result still leaves the worker able to see the builder-derived entity slug/name, entity path, workflow directory, target stage, worktree path when applicable, and checklist. The wrapper may display a label/phase, but those are labels around the same artifact, not replacement facts.
5. **Limited manual fallback.** Manual Pi prompt composition is a break-glass path only when `spacedock dispatch build` is unavailable or exits non-zero, or when a developer is explicitly debugging the dispatch builder itself. The fallback must state the builder failure/unavailability reason and should include the minimum fields from the canonical schema so the degraded dispatch is auditable.
6. **External/failable proof.** The final implementation should add tests that fail if the Pi FO guidance or helper/wrapper code regresses to hand-rolled prompts. Static text checks are acceptable only for text-level claims over runtime guidance; behavioral claims should use fixture or unit tests that parse real builder output and wrapper values.

## Current repo observations

- `skills/first-officer/references/pi-first-officer-runtime.md` already says to use `spacedock dispatch build` with `host: "pi"`, forward the emitted dispatch-file prompt, and avoid `subagent(... acceptance: ...)` for Spacedock stage dispatches.
- `internal/dispatch/build_pi_host_test.go` already verifies the Pi host shape omits Claude team syntax and preserves split-root entity paths in dispatch content.
- Existing coverage does not yet appear to prove the FO/Pi wrapper relationship end-to-end: that `context: "fresh"` and optional `phase` / `label` are additive around the builder artifact, while `acceptance` is absent and builder-derived slug/stage/workflow/worktree/checklist fields remain the source of truth.

No spike needed: the design composes already-present mechanisms (`dispatch build` host selection, Pi dispatch-file prompt emission, split-root path preservation tests, and Pi runtime guidance). The riskiest remaining work is test coverage for wrapper behavior rather than an unknown parser or runtime capability.

## Out of scope

- Changing the fresh-context boundary itself; that belongs to `pi-stage-dispatch-fresh-context`.
- Adding a `pi-subagents` acceptance contract for Spacedock stage work.
- Redesigning `spacedock dispatch build` output for all hosts.
- Adding PR/mod behavior or changing stage advancement semantics.
- Solving reusable Pi worker/session follow-up beyond ensuring this initial stage dispatch uses a canonical assignment artifact.

## Acceptance criteria

Each AC names a property of the finished implementation, not stage work, and cites proof outside this task body.

**AC-1 - Pi FO stage-dispatch guidance makes `spacedock dispatch build` the required assignment source.**
Verified by: a skill/integration invariant test over `skills/first-officer/references/pi-first-officer-runtime.md` (or the generated installed skill text) that fails unless the Dispatch section requires `spacedock dispatch build` with `host: "pi"` and requires forwarding the emitted dispatch file prompt/content instead of composing a replacement assignment.

**AC-2 - Pi subagent wrapper fields are additive to the builder artifact.**
Verified by: a unit or fixture test for the Pi dispatch/wrapper path that builds a real `host: "pi"` dispatch artifact, wraps it for `subagent(...)`, and asserts the wrapper includes `context: "fresh"`, may include `phase` / `label`, and still forwards the exact dispatch file path or content produced by the builder.

**AC-3 - Pi dispatch wrapping preserves builder-derived stage facts.**
Verified by: `internal/dispatch` or integration fixture tests that parse the emitted Pi dispatch artifact and wrapped worker input, asserting the entity slug/name, target stage, entity path, workflow directory, stage checklist, and worktree path when applicable are present from the builder output and are not replaced by wrapper-local strings.

**AC-4 - Pi stage dispatches do not use same-agent acceptance contracts.**
Verified by: a unit/integration test over the Pi wrapper or runtime guidance that fails if a Spacedock stage dispatch includes a `subagent` `acceptance` field while still allowing acceptance requirements to appear inside the canonical dispatch content/checklist.

**AC-5 - Manual prompt fallback is break-glass and auditable.**
Verified by: a guidance invariant test or command fixture that allows manual Pi prompt composition only when the helper is unavailable/non-zero or the task is explicitly debugging the builder, and requires the fallback path to report the helper failure/unavailability reason in observable output or recorded state.

**AC-6 - Pi dispatch behavior is covered by a failable runtime or fixture smoke.**
Verified by: a focused fixture or live Pi test that runs a stage dispatch from a temp workflow through the FO/Pi path and checks durable evidence outside the transcript: process exit, entity/state changes or lack of inappropriate mutation, state-checkout git log when applicable, and the stage report content.

## Test plan

- **Guidance invariant, low cost.** Extend `skills/integration` tests to parse the Pi first-officer runtime Dispatch section and assert the text-level contract: use `spacedock dispatch build` with `host: "pi"`, forward the emitted dispatch file prompt/content, do not hand-roll replacement assignments, and do not use `acceptance` for Spacedock stage dispatches. This proves the instruction text, not runtime behavior.
- **Builder fixture, low/medium cost.** Add or extend `internal/dispatch` Pi host tests using a split-root, folder-form entity and a worktree stage. Assert the emitted artifact carries the slug/name, entity path, workflow directory, stage, checklist, and worktree path. Reuse existing fixture style from `build_pi_host_test.go` and dispatch integration tests.
- **Wrapper fixture, medium cost.** Add the smallest test around the code or adapter that constructs the Pi subagent call. Feed it a real dispatch-build result and assert the wrapper adds `context: "fresh"`, optional `phase` / `label`, omits `acceptance`, and forwards the builder artifact unchanged. Include a negative case that would fail if the wrapper uses a hand-written prompt.
- **Fallback fixture, low/medium cost.** Exercise the builder-failure branch with a fake non-zero helper result. Assert the fallback is not used when the helper succeeds and that the fallback records the reason when it is used.
- **Pi live smoke, higher cost / only if runtime behavior is in scope for implementation.** Dispatch a Pi ensign against a temp split-root workflow and verify durable state evidence: exit code, entity state/report content, state checkout git log, and clean/expected git status. Do not rely on transcript phrasing as the proof.

## Residual risks and questions

- The exact location of Pi wrapper code may still be moving; implementation should choose the smallest existing seam rather than inventing a new architecture.
- If the current Pi runtime only has instruction text and no wrapper helper, AC-2/AC-4 may need to be proven by a fixture harness that models the FO-produced call from real builder output until a concrete adapter exists.
- Manual fallback should remain narrow; broadening it would reintroduce the hand-rolled prompt drift this task is meant to prevent.

## Stage Report: ideation

- Read `docs/dev/README.md` workflow rules and this entity fully.
- Reframed the task around assignment source of truth rather than fresh-context boundaries.
- Added the proposed canonical dispatch-build flow, additive Pi wrapper constraints, no-acceptance rule, break-glass fallback, acceptance criteria with external/failable proof, and a staged test plan.

## Stage Report: implementation

- DONE: Pi FO guidance requires canonical dispatch-builder assignment content.
  Evidence: product commit `613cb086` updates `skills/first-officer/references/pi-first-officer-runtime.md`; `skills/integration/skill_surface_test.go` now requires the Dispatch section to use `spacedock dispatch build` with `host: "pi"`, forward the emitted dispatch-file prompt/content, and avoid replacement hand-written assignments.
- DONE: Pi subagent wrapper fields are additive to the builder artifact.
  Evidence: `internal/piruntime/subagents.go` adds `SubagentStageDispatch`; `internal/piruntime/subagents_test.go` asserts the wrapper carries `context: "fresh"`, optional `phase`/`label`, no acceptance contract, and the exact builder assignment text.
- DONE: Pi dispatch artifacts expose the stage facts needed for naming/audit.
  Evidence: `internal/dispatch/build_pi_host_test.go` now builds a real `host: "pi"` artifact and asserts title/slug-stage path, target stage, entity path, workflow directory stage-definition command, worktree path, and checklist facts are present and survive wrapping.
- DONE: Same-agent acceptance remains forbidden for Spacedock stage dispatches.
  Evidence: wrapper tests assert no `acceptance` field is emitted; skill-surface tests retain the `subagent(... acceptance: ...)` ban and require prompt/entity-based acceptance proof.
- DONE: Manual prompt fallback is documented as break-glass/debug only.
  Evidence: skill-surface invariant requires fallback guidance to mention builder failure/unavailability/debug, observable reason reporting, and canonical schema facts.
- SKIPPED: Live Pi FO dispatch smoke.
  Rationale: this implementation slice proves dispatch-builder, wrapper-helper, and guidance-invariant seams. No full FO live Pi execution path was added.

Changed files:
- `internal/piruntime/subagents.go`
- `internal/piruntime/subagents_test.go`
- `internal/dispatch/build_pi_host_test.go`
- `skills/first-officer/references/pi-first-officer-runtime.md`
- `skills/integration/skill_surface_test.go`

Product commit:
- `613cb086` (`Ensure Pi dispatch wraps build artifacts`)

Validation commands reported by implementation worker:
- PASS: `gofmt -w ./cmd ./internal` (incidental unrelated `internal/status/*` gofmt churn reverted before commit)
- PASS: `go test ./internal/piruntime ./internal/dispatch ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

### Residual risks

- No full FO live Pi dispatch execution path was exercised in this slice; coverage is at the dispatch-builder, wrapper-helper, and guidance-invariant seams.

## Stage Report: validation

Recommendation: REJECTED

Validated product commit: `613cb0867561469c8bad1dc3cc43ed6e9f45072a` (`Ensure Pi dispatch wraps build artifacts`).

Commands run:
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact status --short --branch`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact rev-parse HEAD`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact show --stat --oneline --decorate --no-renames 613cb086`
- PASS: `go test ./internal/piruntime ./internal/dispatch ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

Acceptance criteria evidence:
- AC-1 PASSED: `skills/first-officer/references/pi-first-officer-runtime.md` Dispatch section requires `spacedock dispatch build` with `host: "pi"`, calls the build artifact the assignment source of truth, and says to forward the emitted dispatch file prompt/content without composing a replacement assignment. `skills/integration/skill_surface_test.go` includes `TestPiFirstOfficerRuntimeRequiresDispatchBuildArtifactForStages` locking those invariants.
- AC-2 PASSED: `internal/piruntime/subagents.go` adds `SubagentStageDispatch`; `internal/piruntime/subagents_test.go` asserts `context: "fresh"`, phase/label preservation, exact assignment preservation, and no serialized `acceptance` field. `internal/dispatch/build_pi_host_test.go` wraps a real Pi dispatch-build prompt and asserts the wrapper task equals the builder prompt.
- AC-3 PASSED: `internal/dispatch/build_pi_host_test.go` builds a real `host: "pi"` artifact from a temp split-root workflow and asserts builder-derived title/slug-stage path, target stage, entity path, workflow dir stage-definition command, worktree path, and checklist facts appear in the dispatch body; it also asserts the wrapper forwards the dispatch file path.
- AC-4 PASSED: `internal/piruntime/subagents_test.go` asserts the wrapper JSON omits `acceptance`; `skills/integration/skill_surface_test.go` includes `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages` over the Pi FO Dispatch guidance.
- AC-5 PASSED: `skills/first-officer/references/pi-first-officer-runtime.md` documents manual Pi prompt composition as break-glass only for builder unavailable/non-zero or explicit builder debugging, with required reason recording and canonical schema facts; `TestPiFirstOfficerRuntimeLimitsManualPromptFallback` locks that guidance.
- AC-6 REJECTED: no focused runtime/fixture smoke runs a stage dispatch through the FO/Pi path or checks durable evidence such as process exit plus state/report content and state-checkout git log. The new builder/wrapper tests are useful seam coverage, but they stop before an FO/Pi stage dispatch and do not verify state/report evidence outside the transcript. The implementation report also explicitly marked live Pi FO dispatch smoke as skipped.

Residual risks:
- The concrete FO/Pi dispatch execution path remains unproven by durable fixture/live evidence.
- Guidance tests are string-invariant checks; they do not prove a first officer actually invokes `SubagentStageDispatch` or handles builder fallback at runtime.

## Stage Report: implementation fixback

- DONE: Added focused AC-6 fixture smoke through the FO/Pi dispatch path.
  Evidence: product commit `2f46b7e6` adds `TestPiStageDispatchSmokeUsesBuildArtifactThroughWrapper` in `internal/dispatch/build_pi_host_test.go`. The test builds a real `host: "pi"` dispatch artifact, wraps it with `piruntime.SubagentStageDispatch`, runs a fixture Pi worker as a subprocess, and requires the worker to read the dispatch artifact rather than transcript text.
- DONE: The smoke checks durable evidence outside transcript.
  Evidence: the fixture worker appends a stage report to a split-root entity, commits that entity path in a temp state checkout with message `fixture Pi stage report`, and the parent test verifies process exit, entity report content, state-checkout git log, and clean path-scoped state status.
- DONE: AC-1..AC-5 remain covered by the existing guidance, wrapper, and builder tests from product commit `613cb086`, and the fixback keeps those tests green.

Changed files:
- `internal/dispatch/build_pi_host_test.go`

Product commits:
- `613cb086` (`Ensure Pi dispatch wraps build artifacts`)
- `2f46b7e6` (`Add Pi stage dispatch smoke fixture`)

Validation commands run by implementation fixback worker:
- PASS: `gofmt -w internal/dispatch/build_pi_host_test.go`
- PASS: `go test ./internal/dispatch -run 'TestPiStageDispatchSmoke|TestBuildPiHost' -count=1`
- PASS: `go test ./internal/piruntime ./internal/dispatch ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

AC coverage:
- AC-1: PASS, retained skill-surface invariant requiring `spacedock dispatch build` with `host: "pi"` as the assignment source.
- AC-2: PASS, retained wrapper helper/test and smoke wrapper use of `context: "fresh"`, phase, and label around the builder prompt.
- AC-3: PASS, retained builder fact assertions and added smoke worker assertions over entity path, workflow stage-def command, worktree path, stage, title, and checklist in the actual dispatch file.
- AC-4: PASS, retained no-acceptance assertions and added smoke assertion that wrapper JSON does not serialize `acceptance`.
- AC-5: PASS, retained break-glass fallback guidance invariant.
- AC-6: PASS, added failable fixture smoke with subprocess result plus durable state/report/git-log evidence.

Residual risks:
- The smoke uses a fixture Pi worker process rather than a live Pi runtime session; it proves the FO/Pi artifact-and-wrapper path with durable state evidence without depending on external Pi availability.

## Stage Report: validation z6

Recommendation: PASSED

Validated product commits:
- `613cb0867561469c8bad1dc3cc43ed6e9f45072a` (`Ensure Pi dispatch wraps build artifacts`)
- `2f46b7e6ce6dc740486c3b364989029efae5c116` (`Add Pi stage dispatch smoke fixture`)

Commands run:
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact status --short --branch`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact log --oneline -5`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact show --stat --oneline --decorate --no-renames 613cb086`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact show --stat --oneline --decorate --no-renames 2f46b7e6`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact branch --show-current`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact rev-parse HEAD`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact status --short`
- PASS: `go test ./internal/dispatch -run 'TestPiStageDispatchSmoke|TestBuildPiHost' -count=1`
- PASS: `go test ./internal/piruntime ./internal/dispatch ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`
- PASS: final product and state `git status --short` checks were clean before this validation report append.

Acceptance criteria evidence:
- AC-1 PASSED: `skills/first-officer/references/pi-first-officer-runtime.md` Dispatch guidance requires `spacedock dispatch build` with `host: "pi"`, the dispatch artifact as assignment source of truth, and forwarding emitted dispatch prompt/content without composing a replacement assignment. `skills/integration/skill_surface_test.go` locks these requirements in `TestPiFirstOfficerRuntimeRequiresDispatchBuildArtifactForStages`.
- AC-2 PASSED: `internal/piruntime.SubagentStageDispatch` adds `context: "fresh"` plus optional phase/label around the supplied assignment; `internal/piruntime/subagents_test.go` verifies exact assignment preservation and no serialized `acceptance`; `internal/dispatch/build_pi_host_test.go` wraps a real Pi build prompt and asserts it is forwarded unchanged.
- AC-3 PASSED: `TestBuildPiHostArtifactCarriesCanonicalStageFactsThroughPiWrapper` builds a real split-root `host: "pi"` artifact and asserts builder-derived title/slug-stage path, target stage, entity path, workflow stage-def command, worktree path, and checklist facts are present and that the wrapper forwards the dispatch file path.
- AC-4 PASSED: wrapper JSON is asserted not to contain `acceptance` in both piruntime tests and the AC-6 smoke fixture; skill-surface guidance test also forbids `subagent(... acceptance: ...)` for Spacedock stage dispatches.
- AC-5 PASSED: `TestPiFirstOfficerRuntimeLimitsManualPromptFallback` verifies guidance limits manual Pi prompt composition to builder unavailable/non-zero or explicit builder debugging and requires recording an auditable failure/unavailability reason with minimum canonical schema facts.
- AC-6 PASSED: product commit `2f46b7e6` adds `TestPiStageDispatchSmokeUsesBuildArtifactThroughWrapper`, a failable fixture smoke that builds a Pi dispatch artifact, wraps it, runs a fixture Pi worker subprocess, requires the worker to read the dispatch file, appends a durable stage report, commits the entity path in a temp state checkout, and verifies process exit, state/report content, state git log (`fixture Pi stage report`), and clean path-scoped state status.

Residual risks:
- AC-6 uses a fixture Pi worker subprocess, not a live external Pi runtime session; this is acceptable for the stated fixture/live option but does not prove external runtime availability.
- The product branch is ahead of and behind `origin/next` in the worktree status observed at validation start; validation was performed on the assigned commits at HEAD without rebasing or editing product files.

## Stage Report: implementation z6 CI fixback

- DONE: Fixed the PR #300 offline `internal/dispatch` failure caused by launcher-command rendering from PR #298.
  Evidence: product commit `f1aec201` updates `internal/dispatch/build_pi_host_test.go` so Pi dispatch tests no longer require the obsolete bare `spacedock dispatch show-stage-def ...` substring.
- DONE: The builder/wrapper tests now assert the current semantic stage-definition command survives launcher rendering.
  Evidence: added `assertRenderedStageDefCommand`, which requires the rendered `spacedock_launcher() { ... }` shim and the semantic `spacedock_launcher dispatch show-stage-def --workflow-dir ... --stage implementation` invocation in both the canonical artifact test and fixture smoke worker.
- DONE: Existing AC coverage is preserved: canonical build artifact source, additive fresh Pi wrapper, builder-derived stage facts, no same-agent acceptance, break-glass fallback guidance, and fixture smoke durable evidence all remain covered.

Changed files:
- `internal/dispatch/build_pi_host_test.go`

Product commits:
- `f1aec201` (`Fix Pi dispatch stage-def wrapper assertions`)

Validation commands run:
- PASS: `gofmt -w internal/dispatch/build_pi_host_test.go`
- PASS: `go test ./internal/dispatch -run 'TestPiStageDispatchSmoke|TestBuildPiHost' -count=1`
- PASS: `go test ./internal/piruntime ./internal/dispatch ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`

AC coverage:
- AC-1: PASS, retained skill-surface invariant requiring `spacedock dispatch build` with `host: "pi"` as assignment source.
- AC-2: PASS, retained wrapper helper/test; the fixture smoke continues to use the Pi wrapper around the builder prompt.
- AC-3: PASS, builder fact assertions now accept current launcher-rendered stage-def commands while still checking the semantic `dispatch show-stage-def` invocation.
- AC-4: PASS, no-acceptance assertions remain unchanged.
- AC-5: PASS, break-glass fallback guidance invariant remains unchanged.
- AC-6: PASS, fixture smoke now validates the current launcher-rendered fetch command and durable state/report/git-log evidence.

Residual risks:
- The smoke remains a fixture Pi worker subprocess rather than a live external Pi runtime session.

## Stage Report: validation z6 rerun

Recommendation: PASSED

Validated product commits:
- `b68bb89ebde8f00adec99b9de239aae3478fce4c` (`Ensure Pi dispatch wraps build artifacts`)
- `7b289f5ab6fa05a73b0ea53ffdaf24736d823c08` (`Add Pi stage dispatch smoke fixture`)
- `f1aec20183e700dbee7f00a89b54e7abc10c4e2c` (`Fix Pi dispatch stage-def wrapper assertions`)

Commands run:
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact status --short --branch`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact log --oneline -5`
- PASS: `git -C /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-stage-dispatch-uses-build-artifact show --stat --oneline --decorate --no-renames b68bb89e 7b289f5a f1aec201`
- PASS: `spacedock dispatch show-stage-def --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --stage validation`
- PASS: `go test ./internal/dispatch -run 'TestPiStageDispatchSmoke|TestBuildPiHost' -count=1`
- PASS: `go test ./internal/piruntime ./internal/dispatch ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`
- PASS: final product `git status --short` was clean before this report append.

Acceptance criteria evidence:
- AC-1 PASSED: Pi FO guidance requires `spacedock dispatch build` with `host: "pi"`, treats the build artifact as the assignment source of truth, and requires forwarding emitted prompt/content without composing a replacement. Locked by `TestPiFirstOfficerRuntimeRequiresDispatchBuildArtifactForStages`.
- AC-2 PASSED: `internal/piruntime.SubagentStageDispatch` adds only Pi transport fields (`context: "fresh"`, optional phase/label) around the supplied assignment, preserves the assignment exactly, and serializes no `acceptance` field; dispatch tests wrap a real builder prompt unchanged.
- AC-3 PASSED: `TestBuildPiHostArtifactCarriesCanonicalStageFactsThroughPiWrapper` verifies builder-derived title/slug-stage path, target stage, entity path, workflow directory, worktree path, checklist facts, and the launcher-rendered semantic `dispatch show-stage-def` invocation survive through the Pi wrapper.
- AC-4 PASSED: wrapper and smoke tests assert no serialized `acceptance`; skill-surface tests forbid `subagent(... acceptance: ...)` for Spacedock stage dispatches while preserving acceptance content inside canonical dispatch/checklists.
- AC-5 PASSED: guidance invariant test keeps manual Pi prompt composition limited to builder unavailable/non-zero or explicit builder debugging and requires an auditable failure/unavailability reason plus canonical schema facts.
- AC-6 PASSED: `TestPiStageDispatchSmokeUsesBuildArtifactThroughWrapper` builds a Pi dispatch artifact, wraps it, runs a fixture Pi worker subprocess, requires reading the dispatch artifact, appends a durable stage report, commits the entity path in a temp state checkout, and verifies process exit, state/report content, state git log, and clean path-scoped state status.

Residual risks:
- AC-6 is fixture-backed, not a live external Pi runtime session; it proves the artifact/wrapper path and durable state evidence without proving external Pi runtime availability.
- Guidance tests remain string-invariant checks for FO prose; runtime proof of an actual human FO using the helper is represented by the fixture seam, not a live orchestrator run.
