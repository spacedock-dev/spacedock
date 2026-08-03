---
id: 3616bd590vjvfaafngevnmsk
title: Wire TestLivePiSubagentEnsignSmoke into pi-live CI selector
status: backlog
source: "Two audits this session found TestLivePiSubagentEnsignSmoke (internal/ensigncycle/pi_live_runner_test.go:21) -- which carries assertPiEnsignBootContract, the AC-1 pi ensign-boot-contract grader landed in 66aa8a610 -- has never been in the pi-live CI job's -run selector (.github/workflows/runtime-live-e2e.yml:830, currently 'TestLivePiFrontDoorSmoke|TestLivePiRecordedGateLifecycle'), despite a debrief (docs/dev/.spacedock-state/_debriefs/2026-08-01-02-pi-kimi-k3.md:15) incorrectly claiming it runs green there. git history confirms this test's name has never appeared in any CI workflow file since its introduction (2026-06-03, ce192a7c2) -- orphaned from creation, not regressed. Captain directed: file tasks to reconcile the gaps found, with docs/runtime-live-ci.md as the SSOT to keep current."
started:
completed:
verdict:
score: 0.8
worktree:
issue:
sprint: live-test-truth
group: runtime-specific
sprint-readiness: ready
---

Add `TestLivePiSubagentEnsignSmoke` to the pi-live job's `-run` selector so the AC-1 ensign-boot-contract grader it carries actually gates CI, instead of being a landed-but-never-checked assertion.

## Proposed approach

Update `.github/workflows/runtime-live-e2e.yml`'s pi-live step's `-run` pattern to include `TestLivePiSubagentEnsignSmoke` alongside the existing `TestLivePiFrontDoorSmoke|TestLivePiRecordedGateLifecycle`. Confirm the test passes live in that CI context (not just locally) before landing -- this is the point of the change. Correct the inaccurate claim in `docs/dev/.spacedock-state/_debriefs/2026-08-01-02-pi-kimi-k3.md:15` if that file is still live/active (check first; do not edit an archived debrief without checking whether that's appropriate). Update `docs/runtime-live-ci.md` to record this test as pi-live-selected.

## Acceptance criteria

**AC-1 (VALUE)** -- `TestLivePiSubagentEnsignSmoke` executes and passes in an actual pi-live CI run (not just local `-tags live` invocation), proving `assertPiEnsignBootContract` gates real CI instead of sitting unreachable.
Verified by: a green pi-live CI job whose log shows this specific test name executed, on the PR that makes this change.

**AC-2** -- No regression to the two tests already selected on pi-live.
Verified by: `TestLivePiFrontDoorSmoke` and `TestLivePiRecordedGateLifecycle` (the latter already skips per its own TODO, unaffected) still behave identically.

**AC-3** -- `docs/runtime-live-ci.md` reflects this test as pi-live-selected going forward.

## Test plan

This requires an actual live CI run on the PR (pi-live lane) since the claim being proven is "this executes in CI," not just "the Go code compiles." No offline-only substitute proves AC-1.

## Out of scope

Any other pi-live selector changes, any change to `assertPiEnsignBootContract` itself (already landed and correct), the other 8 unwired live scenario tests.
