---
title: Pi stage dispatches should force fresh subagent context
status: validation
source: captain (2026-06-04) — FO mistakenly dispatched a Spacedock implementation worker with pi-subagents context=fork; stage workers should be fresh and independent
score: "0.29"
started: 2026-06-04T07:38:36Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-stage-dispatch-fresh-context
issue:
id: d2w8z614c0q1yssmyr33a38y
mod-block: merge:pr-merge
---

Spacedock's Pi first-officer runtime should make stage dispatch context explicit. During `launcher-binary-path-passthrough`, the FO dispatched an implementation worker with `context: "fork"` because the builtin worker agent defaults to forked context. That preserved parent context but conflicts with the intended Spacedock stage model: implementation and validation workers should receive the assignment/context from the dispatch prompt and entity state, not inherit the FO session transcript by default.

## Problem

Pi support currently documents that Spacedock stage dispatches may use the Pi-native `subagent(...)` tool, must not use `pi-subagents` acceptance contracts, and should prefer fresh redispatch for follow-up/retry cycles. The missing guard is the `context` parameter passed to `subagent(...)`. If the FO omits it, the worker agent default can select forked context. That creates a hidden dependency on the FO transcript and can contaminate stage work with prior reasoning, validation expectations, captain-side discussion, or unrelated delegation state that was never part of the durable assignment.

This weakens the Spacedock workflow boundary in three ways:

1. **Independence is accidental.** Stage workers are expected to be seeded by the dispatch prompt, workflow directory, entity file, and assigned stage checklist. A forked worker may succeed only because inherited transcript context filled gaps in the assignment.
2. **Validation can be biased.** A validation worker inheriting implementation discussion may test the implementer's assumptions rather than independently reproducing the acceptance criteria from the entity and resulting files.
3. **Feedback cycles can become stale.** Reusing a previous subagent context for a retry can make an old completion or old mental model appear current unless the assignment is explicitly a manual/debug resume with durable epoch evidence.

## Proposed approach

Implement a compatibility-first guard around the Pi first-officer stage-dispatch contract:

1. Update the Pi first-officer runtime guidance so any Spacedock stage dispatch through `pi-subagents` explicitly passes `context: "fresh"` to `subagent(...)`. The guidance should state that relying on the worker agent's default context is forbidden for Spacedock stages.
2. Keep the existing rule that Spacedock stage dispatches must not use `subagent(... acceptance: ...)`. Acceptance requirements remain part of the dispatch prompt/entity contract, and completion remains proven by product/state commits plus entity stage reports and independent validation.
3. Clarify follow-up and retry behavior: normal feedback/rework dispatch is a new fresh assignment cycle. Any non-fresh resume is an explicit manual/debug exception only, and it must be visibly marked as such and tied to durable entity-stage evidence such as worker metadata, stage, epoch, and current state.
4. Add tests at the smallest enforceable surface:
   - a skill-surface/instruction invariant test that parses `skills/first-officer/references/pi-first-officer-runtime.md` and fails unless the Dispatch section contains the required fresh-context invariant;
   - extend the existing acceptance-contract invariant test so the fresh-context rule and no-acceptance rule are checked together for Pi stage dispatches;
   - add or extend a follow-up/reuse invariant over the Follow-up section so normal retry cycles are fresh and manual/debug resume is the only documented exception.

The first implementation can be instruction-text plus invariant tests because the current Pi stage dispatch path is runtime-guided rather than a typed Go dispatcher. If a later implementation introduces a structured Pi dispatch builder or adapter, this task's tests should move downward into that typed API and fail on an actual `subagent(...)` call without `context: "fresh"`.

### Spike determination

No live-runtime spike is needed for ideation: the design relies on existing, already-shipped mechanisms in this repo:

- `skills/first-officer/references/pi-first-officer-runtime.md` already defines the Pi stage-dispatch contract;
- `skills/integration/skill_surface_test.go` already contains `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages`, a failing static invariant over the same runtime guidance;
- Pi subagent dispatch is already allowed by the current runtime contract, while Spacedock-owned acceptance via `subagent(... acceptance: ...)` is already forbidden.

The riskiest unverified behavior is whether the live Pi host honors `context: "fresh"` exactly as named. That should be covered by a later live Pi runtime scenario if/when live Pi dispatch tests are added; it is not necessary to establish the instruction-level guard because this task's implementation target is the Spacedock dispatch contract and its repo tests.

## Out of scope

- Changing `launcher-binary-path-passthrough` or redispatching any existing worker.
- Adding PR/mod behavior or changing non-Pi runtime dispatch semantics.
- Banning subagent dispatch itself; Pi subagent dispatch remains allowed.
- Reintroducing or depending on `pi-subagents` `acceptance` contracts for Spacedock stages.
- Building a full typed Pi subagent adapter if the current code path is still instruction-driven.
- Proving live Pi host behavior in CI unless the implementation already has a live Pi harness available at low cost.

## Acceptance criteria

Each AC names an end-state property of the finished deliverable and cites proof outside this task body.

**AC-1 - Pi/Spacedock stage dispatch instructions require explicit fresh subagent context.**
Verified by: a Go test in `skills/integration` that reads `skills/first-officer/references/pi-first-officer-runtime.md`, isolates the `## Dispatch` section, and fails unless Spacedock stage dispatch guidance for `pi-subagents` requires `context: "fresh"` and forbids relying on the worker agent's default context.

**AC-2 - Spacedock stage dispatches still avoid `pi-subagents` acceptance contracts.**
Verified by: `go test ./skills/integration -run TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages` or its extended successor, which fails unless the Pi Dispatch section forbids `subagent(... acceptance: ...)` and says acceptance requirements live in the task prompt/dispatch content with completion proven by entity stage reports, product/state commits, and independent validation.

**AC-3 - Normal Pi feedback/retry dispatches are fresh assignment cycles, not context resumes.**
Verified by: a Go test in `skills/integration` that reads the `## Follow-up and Reuse` section of `skills/first-officer/references/pi-first-officer-runtime.md` and fails unless normal follow-up/retry dispatch is documented as fresh by default, previous completions cannot satisfy a new epoch, and any non-fresh manual/debug resume is explicitly marked as an exception requiring durable metadata.

**AC-4 - The repo-level test suite catches regressions in the Pi stage-dispatch contract.**
Verified by: `go test ./...` after implementation, with the Pi skill-surface tests included in the normal suite and failing if the fresh-context or no-acceptance invariants are removed.

## Test plan

- **Focused invariant tests (low cost, fixture/static-contract level):** add or extend tests in `skills/integration/skill_surface_test.go` that parse sections of `skills/first-officer/references/pi-first-officer-runtime.md`. These tests prove the shipped instruction contract carries the enforceable invariants at the same abstraction level as the current implementation.
- **Baseline Go suite (medium cost):** run `go test ./...` to prove the new invariant tests are integrated into the repo's normal gate and no unrelated packages regress.
- **Race suite (medium/high cost, optional unless implementation touches concurrent Go code):** if the implementation only changes skill prose and invariant tests, `go test ./... -race` is low-risk but not necessary to prove this text-contract change. If a typed dispatcher or runtime adapter is changed, run it before validation.
- **Live Pi runtime smoke (defer unless harness exists):** when a durable live Pi harness is available, add a scenario that dispatches an implementation or validation stage with a sentinel string present only in the FO transcript and asserts the worker does not see it unless included in the dispatch prompt. This would prove host-level fresh context, but it is not required for the first instruction/test guard.

## Stage report

### Ideation — 2026-06-04

- Read `docs/dev/README.md` and this entity.
- Checked the current Pi runtime contract in `skills/first-officer/references/pi-first-officer-runtime.md`.
- Found existing test coverage for the no-`acceptance` Pi stage-dispatch invariant in `skills/integration/skill_surface_test.go`.
- Refined the problem, proposed approach, out-of-scope boundary, acceptance criteria, and test plan to align with current Pi contracts: subagent dispatch is allowed, `acceptance` contracts are forbidden, and feedback/retry defaults to fresh assignment cycles unless an explicit manual/debug resume is introduced with durable metadata.
- Noted that `/Users/clkao/git/spacedock-research/spacedock-v1/context.md` and `/Users/clkao/git/spacedock-research/spacedock-v1/plan.md` were requested but are not present in this checkout.

### Implementation — 2026-06-04

- Updated the Pi first-officer runtime dispatch contract so Spacedock stage dispatches through `pi-subagents` must call `subagent(...)` with explicit `context: "fresh"` and must not rely on the worker agent default context.
- Preserved the no-`subagent(... acceptance: ...)` invariant for Spacedock stage dispatches; acceptance remains in the task prompt/dispatch content and completion proof remains entity stage reports, product/state commits, and independent validation.
- Clarified Follow-up and Reuse guidance so normal follow-up/retry dispatches are fresh assignment cycles, previous completions cannot satisfy a new epoch, and non-fresh resume is only a marked manual/debug exception tied to durable metadata.
- Added skill-surface invariant tests for explicit fresh Pi stage context and fresh follow-up/retry behavior, alongside the existing no-acceptance-contract invariant.
- Product commit: `ea80f3a2` (`Require fresh Pi stage dispatch context`).
- Validation run from the product worktree:
  - `go test ./skills/integration`
  - `go test ./...`
  - `gofmt -w ./cmd ./internal` (temporarily reformatted two unrelated comments under `internal/status`; those unrelated changes were reverted before commit)
  - `go test ./... -race`
- AC coverage: AC-1 covered by `TestPiFirstOfficerRuntimeRequiresFreshSubagentContextForStages`; AC-2 covered by `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages`; AC-3 covered by `TestPiFirstOfficerRuntimeFollowupsAreFreshByDefault`; AC-4 covered by `go test ./...` including `skills/integration`.
- Residual risk: live Pi host behavior for the exact `context: "fresh"` argument remains unproven by this instruction-level guard and should be covered later by a live Pi harness if available.

### Validation — 2026-06-04 (d2 validation)

Recommendation: PASSED.

Validated product commit `ea80f3a2edb0fc2eeb07bb7eef7c8012c6e47324` in `.worktrees/spacedock-ensign-pi-stage-dispatch-fresh-context` on branch `spacedock-ensign/pi-stage-dispatch-fresh-context`.

Evidence reviewed:
- Product commit changes only `skills/first-officer/references/pi-first-officer-runtime.md` and `skills/integration/skill_surface_test.go`.
- `## Dispatch` now requires Spacedock Pi stage dispatches through `pi-subagents` to call `subagent(...)` with explicit `context: "fresh"` and forbids relying on the worker agent default context.
- `## Dispatch` still forbids `subagent(... acceptance: ...)` for Spacedock stage dispatches and keeps acceptance/completion proof in task prompt/dispatch content, entity stage reports, product/state commits, and independent validation.
- `## Follow-up and Reuse` now documents normal follow-up/retry dispatches as fresh assignment cycles, requires epoch increment, forbids previous completions satisfying the new epoch, and restricts non-fresh resume to an explicit manual/debug exception tied to durable metadata.
- Added invariant tests: `TestPiFirstOfficerRuntimeRequiresFreshSubagentContextForStages`, `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages`, and `TestPiFirstOfficerRuntimeFollowupsAreFreshByDefault`.

Commands run from the product worktree:
- `git rev-parse HEAD` → `ea80f3a2edb0fc2eeb07bb7eef7c8012c6e47324`.
- `git status --short --branch` → `## spacedock-ensign/pi-stage-dispatch-fresh-context...origin/next [ahead 1]`.
- `git show --stat --oneline --decorate --no-renames HEAD`.
- `git show --name-only --format=fuller --no-renames HEAD`.
- `go test ./skills/integration -run 'TestPiFirstOfficerRuntime(RequiresFreshSubagentContextForStages|ForbidsSubagentAcceptanceForStages|FollowupsAreFreshByDefault)' -count=1` → PASS.
- `go test ./skills/integration -count=1` → PASS.
- `go test ./... -count=1` → PASS.
- `go test ./... -race -count=1` → PASS.
- `git status --short` → clean.

AC evidence:
- AC-1: PASS — focused invariant test `TestPiFirstOfficerRuntimeRequiresFreshSubagentContextForStages` passed and the Dispatch section contains the required explicit `context: "fresh"`/no-default-context language.
- AC-2: PASS — focused invariant test `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages` passed and the Dispatch section keeps `subagent(... acceptance: ...)` forbidden for Spacedock stages.
- AC-3: PASS — focused invariant test `TestPiFirstOfficerRuntimeFollowupsAreFreshByDefault` passed and the Follow-up and Reuse section documents fresh retry cycles plus the manual/debug-only non-fresh exception with durable metadata.
- AC-4: PASS — `go test ./... -count=1` passed with the Pi skill-surface tests included in the normal suite.

Residual risks:
- Live Pi host behavior for the exact `context: "fresh"` argument remains unproven; this task intentionally validates the instruction-level guard and static invariants only.
- The requested root `context.md` and `plan.md` files were absent in the parent checkout, so validation used the dispatch assignment file and entity state as the available operating contract.
