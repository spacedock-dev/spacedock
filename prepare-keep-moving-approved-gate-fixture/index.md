---
title: Prepare the approved gate in the keep-moving live fixture
status: backlog
source: "Captain filing from PR #784 runtime CI attempt 2, Codex job 103174402234"
id: wvdv5c5mkkd23yhmr9p16czd
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
- Implementation: add the focused regression before the fixture change; run relevant deterministic tests, full Go tests and race tests, and formatting required by AGENTS.md.
- Validation: independently verify setup and refusal behavior, then run the targeted Codex scenario. Any remaining live failure must be diagnosed separately from the repaired setup.
