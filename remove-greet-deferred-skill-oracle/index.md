---
id: zath5jk8c6txzq3rwn8a6m1g
title: Complete removal of the assertGreetInvokesNoDeferredFOSkill greet oracle
status: ideation
source: Follow-up from PR #512 (restore-fo-merge-write-lazy-loading-reset)
started: 2026-07-15T04:14:21Z
completed:
verdict:
score:
worktree:
issue:
---

Complete the removal of the brittle `assertGreetInvokesNoDeferredFOSkill` AC-2 greet oracle, plus its two dedicated unit-test controls and its dedicated fixture, that PR #512 started but left uncommitted in its worktree — while RETAINING the durable shallow-boot guards. The started-but-unmerged removal is preserved verbatim as `preserved-lazy-load-wip.patch` beside this entity (captured from `.worktrees/spacedock-ensign-restore-fo-merge-write-lazy-loading-reset`).

## Problem

{Ideation fills this in. Seed: the oracle keys on the Skill-argument string of a captured live stream — a per-skill behavioral check that is brittle and, after PR #512 restored lazy loading of the FO write/merge cores, redundant with the durable structural guards. The removal was drafted in the #512 worktree but never committed or merged, so `main` still carries the oracle and its controls.}

## Proposed approach

{Ideation fills this in. Seed: on a fresh branch off main, apply the preserved patch — (1) remove `assertGreetInvokesNoDeferredFOSkill` + the `deferredFOSkillNames` var + the now-unused `strings` import from `internal/ensigncycle/shallow_boot_measure_test.go`; (2) remove its two unit tests (`TestAssertGreetInvokesNoDeferredFOSkillOffline`, `TestAssertGreetInvokesNoDeferredFOSkillCatchesLaterDeltaInvoke`) from `shallow_boot_measure_unit_test.go`; (3) remove the call site and reword the AC-2 comment in `claude_live_runner_test.go`; (4) delete `testdata/greet-invokes-fo-status-viewer-skill.stream.jsonl`. RETAIN `assertNoTeamCreateBeforeGreet` and `assertShallowBootMeasured`.}

## Out of scope

{Ideation fills this in. Seed: no behavior change to the retained shallow-boot guards; no change to the lazy-loading behavior PR #512 shipped; no touching of `greetTurnIndex` (stays used by surviving functions).}

## Acceptance criteria

{Ideation fills this in.}

## Test plan

{Ideation fills this in. Seed: gofmt clean; `go build`/`go vet` the package; `go test ./...` and `go test ./... -race` green with the oracle, its two controls, and the fixture gone and the durable guards still present and passing.}
