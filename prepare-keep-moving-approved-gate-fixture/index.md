---
title: Prepare the approved gate in the keep-moving live fixture
status: implementation
source: "Captain filing from PR #784 runtime CI attempt 2, Codex job 103174402234"
id: wvdv5c5mkkd23yhmr9p16czd
gates:
    version: 1
    records:
        - id: gate:wvdv5c5mkkd23yhmr9p16czd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wvdv5c5mkkd23yhmr9p16czd-backlog-1
              briefing:
                id: briefing:wvdv5c5mkkd23yhmr9p16czd:backlog:attempt-1:revision-1
                digest: sha256:4267b265607a4a9ccf99397d80bb759eafbd76382dfeb268c97d9cce9ad97e4c
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:wvdv5c5mkkd23yhmr9p16czd:backlog:1
                briefing: briefing:wvdv5c5mkkd23yhmr9p16czd:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T04:01:30.288629Z"
                decision: approve
                reason: Captain requested dispatch of the Codex live failure tasks, targeted local verification, and stacked PRs.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:wvdv5c5mkkd23yhmr9p16czd:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:wvdv5c5mkkd23yhmr9p16czd-ideation-1
              briefing:
                id: briefing:wvdv5c5mkkd23yhmr9p16czd:ideation:attempt-1:revision-1
                digest: sha256:10a5618ff09104122df3c972d5d2a5df112cc6128e553c98aa7df024b8b36bb1
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:wvdv5c5mkkd23yhmr9p16czd:ideation:1
                briefing: briefing:wvdv5c5mkkd23yhmr9p16czd:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T04:07:04.691315Z"
                decision: approve
                reason: Minimal real gate setup is proven; bounded fixture correction retains authority and dispatch proof. Proceed under captain scope through targeted verification and PR stack.
                conn:
                    quote: dispatch codex live test failure tasks, verify locally for targeted failure, and open PR as stack, then trigger codex ci on stack tip.
                    source: active thread goal supplied by captain
              application:
                target-stage: implementation
                state: consumed
started: 2026-09-15T04:01:46Z
worktree: .worktrees/spacedock-ensign-prepare-keep-moving-approved-gate-fixture
---

## Problem statement

The host-neutral keep-moving fixture says the captain just approved `approved-gate` at review, but seeds no gate record or bound briefing. A conforming `gate record --decision approve --consume` therefore fails before the continuation behavior under test can begin. The test grades the missing implementation dispatch as a keep-moving violation.

This task belongs to the development workflow because its deliverable is a corrected executable live-test fixture and regression proof in `internal/ensigncycle`. Keep the runtime's refusal of unprepared gates intact. This is separate from the existing Codex smallest-mechanism grader repair and from Pi or Sonnet posture defects.

## Evidence

- Run: https://github.com/spacedock-dev/spacedock/actions/runs/34568401666/attempts/2 ; Codex job `103174402234`; candidate `09f123d33b1d7a5b8d3d2dbc85e536739ba7ea1b`.
- Artifact `10188347210`, created 2026-09-11T07:03:37Z; member `live-artifacts/codex/codex-shared-scenarios/keep-moving-posture/codex-exec.jsonl`. Local archive `/tmp/pr784-codex-attempt2.zip` is only a temporary convenience.
- Test `TestLiveCommonKeepMovingPosture` failed after 158.91s with `durable keep-moving journeys = 2/3: map[approved-gate:missing path-scoped dispatch entry with stage and started]`.
- Parent trace shows `gate record approved-gate --decision approve --actor person:captain --reason "Captain approved review in the current request." --consume --workflow-dir <fixture-root>` exiting 1 with `Error: entity has no gates record`.
- The immediately preceding file read proves `approved-gate.md` contains id/title/status=review and no gates record. Codex continues the two independent tasks to archive and reports the approval blocker. This is not an auth error or a provider timeout.
- Source: `writeKeepMovingWorkflow`, `keepMovingApprovedEntity`, and `keepMovingPrompt` in `internal/ensigncycle/shared_fixtures_test.go`; grading in `assertDurableKeepMoving` in `internal/ensigncycle/shared_keep_moving_durable_test.go`.

## Scope and proposed direction

Construct the real prepared review attempt and committed briefing through the supported gate API before presenting the captain's approval to the FO. Keep approval consumption and successor dispatch as observed agent actions, not fixture-precompleted work. Confirm the exact preparation requirements before choosing the smallest fixture setup. Do not weaken approval guards, force the transition, or classify a missing precondition as an agent continuation failure.

The observed setup defect is established; whether additional agent failures remain after repairing it needs the targeted live run. This filing does not claim PR #784 is cleared for merge.

## Acceptance criteria

**AC-1 — The fixture's stated approval can be consumed through the supported gate lifecycle.**
Verified by: a deterministic setup exercise prepares the actual gate, records and consumes approval successfully, and observes implementation as the successor; omit the prepared record and the same consume must fail.

**AC-2 — The live scenario still measures continuation and independent work.**
Verified by: the targeted Codex keep-moving run observes the approved task's implementation dispatch, worker report, and terminal journey, alongside the independent tasks and nonterminal questioned task. Existing durable journey checks remain active; suppressing the approved task's dispatch must still fail the grader.

**AC-3 — Approval authority is not bypassed by fixture preparation.**
Verified by: inspect and exercise the initialized fixture before the captain decision; it has a prepared review attempt but no consumed approval or implementation dispatch. A missing approval cannot advance it. Do not use forced status changes as preparation.

## Stage-specific test gates

- Ideation: reproduce the seed/consume mismatch with the installed candidate and demonstrate the smallest valid gate setup in a throwaway fixture. Name the existing primary proof owners and expected file/LOC scope.
- Implementation: add the focused regression before the fixture change; run the relevant deterministic tests and format changed Go files. Per captain-approved stack scheduling, full Go/race suites and repository formatting run once at the combined stack tip.
- Validation: independently verify setup and refusal behavior, then run the targeted Codex scenario. Any remaining live failure must be diagnosed separately from the repaired setup.

## Implementation plan and expected surface

Use the existing single-root `writeKeepMovingWorkflow` fixture. Add one short review Markdown artifact before `gitInit`; after initialization run the real `gate prepare approved-gate --question "Advance to implementation?" --artifact <root>/gate-review.md --summary "Ready to proceed to implementation." --workflow-dir <root>`, then `state commit approved-gate --workflow-dir <root>`. Reuse `buildRecordedGateBinary` and `mustRecordedGate`; no new framework or split-root fixture is needed. The agent still owns recording the captain decision, consuming it, dispatching implementation, and completing the durable journey.

This preparation serves AC-1 and AC-3. Bare status seeding is insufficient because it has no gate record; hand-authored gate YAML would duplicate digest/binding machinery. The supported prepare and state-commit commands are the smallest exercised route to a real committed review package.

Estimate net LOC change: +100, across 3 files (100 insertions, 0 deletions), tolerance +50 net LOC and no additional files. Expected files: `internal/ensigncycle/shared_fixtures_test.go` (about 10 lines), new `internal/ensigncycle/shared_keep_moving_fixture_test.go` (about 80), and `internal/ensigncycle/shared_keep_moving_durable_test.go` (about 10 for an approved-task missing-dispatch control). Changes affect only test fixture initial state and regression coverage; CLI grammar, stored formats, approval authority, production runtime, skills, and durable grader semantics stay outside scope. No user-visible documentation diff is needed.

### Spike evidence and primary proof

At `origin/main` 2a7b8719843e40b79545f0bb4def6609cdd9ebbf, the installed pre3 binary reproduced unprepared `gate record --decision approve --actor person:captain --consume` exiting 1 with `entity has no gates record`. The throwaway `TestKeepMovingIdeationSpike` then proved prepare plus state commit retains `approved-gate/review/review/briefing-1/index.json` in HEAD; the prepared entity remains at review with no started, resolution, or application. Missing-decision consume exits 1, reports ineligible/consumed=false, and leaves entity bytes unchanged. Recording captain approval with consume exits 0 and writes status=implementation/application.state=consumed.

The spike also established that single-root `gate prepare` alone leaves its room uncommitted; `state commit` is required for the committed-fixture requirement. Retained local evidence: `/tmp/keep-moving-ideation-wvdv/internal/ensigncycle/keep_moving_spike_test.go` and `/tmp/keep-moving-gate-spike-wvdv4` (throwaway paths, not dependencies of shipped tests). Final exercised command: `go test ./internal/ensigncycle -run '^TestKeepMovingIdeationSpike$' -count=1 -v` (1.446 seconds, pass).

Primary proof owners are the existing real command helpers in `recorded_gate_lifecycle_test.go`, `TestDurableTaskJourneys`, and `TestLiveCommonKeepMovingPosture` with `assertDurableKeepMoving`. Add the focused `TestKeepMovingPreparedGate` before the fixture edit: verify committed open preparation, unchanged missing-decision refusal, successful captain approval consumption, and the unprepared refusal control. Removing preparation, preconsuming approval, or omitting the state commit must each fail it. Use on-disk parsed state, command exit codes, Git objects, and byte equality; no prose-grep proof.

AC-2 stays with the existing durable journey proof and one missing-approved-dispatch case in its existing test table: suppressing that dispatch must fail while preserving the other two journeys. The live test must still observe all three dispatch/report/terminal journeys and the questioned task's nonterminal correction. Deterministic checks cost seconds plus one binary build; the single Codex live run costs minutes and is scheduled serially by the FO after implementation.

Targeted deterministic command: `go test ./internal/ensigncycle -run '^(TestKeepMovingPreparedGate|TestDurableTaskJourneys|TestDurableKeepMoving.*|TestDurableQuestioned.*|TestRetainedAtomicWorkerJourney)$' -count=1`.

Exact targeted Codex command, from the prepared stack checkout: `SPACEDOCK_BIN=/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -tags live ./internal/ensigncycle -run '^TestLiveCommonKeepMovingPosture$' -count=1 -v -timeout 15m`. This requires Codex on PATH and the existing supported local auth or CI auth input; the harness handles isolation. The candidate changes only Go test fixtures, so the specified pre3 executable is sufficient. FO runs the full normal/race/CI checks once at the combined four-fix stack tip.

## Stage Report: ideation

- DONE: Prove the smallest real prepared-gate fixture setup and preserve the missing-approval refusal.
  Installed pre3 spike passed: unprepared refusal, real prepare + state commit, committed briefing, unchanged missing-decision refusal, then approval consumption to implementation; dropping state commit failed the Git-object assertion.
- DONE: Record a concise implementation plan, exact targeted Codex command, primary tests, and bounded surface for the captain-approved stack.
  Plan above bounds fixture/regression work to +100 net LOC across 3 files (tolerance +50); existing durable checks and serialized targeted Codex run own continuation proof.

### Summary

The fixture needs a real prepared review room and its state commit before the prompt grants approval. The approval decision and continuation remain agent actions; no production authority or grader relaxation is proposed. Live execution and combined-stack full suites remain scheduled after implementation.
