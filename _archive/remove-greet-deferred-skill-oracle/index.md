---
id: zath5jk8c6txzq3rwn8a6m1g
title: Complete removal of the assertGreetInvokesNoDeferredFOSkill greet oracle
status: done
source: Follow-up from PR #512 (restore-fo-merge-write-lazy-loading-reset)
started: 2026-07-15T04:14:21Z
completed: 2026-07-15T06:35:47Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-remove-greet-deferred-skill-oracle
issue:
mod-block:
pr: pr-merge:513
archived: 2026-07-15T06:35:47Z
---

Complete the removal of the brittle `assertGreetInvokesNoDeferredFOSkill` AC-2 greet oracle, plus its two dedicated unit-test controls and its dedicated fixture, that PR #512 started but left uncommitted in its worktree — while RETAINING the durable shallow-boot guards. The started-but-unmerged removal is preserved verbatim as `preserved-lazy-load-wip.patch` beside this entity (captured from `.worktrees/spacedock-ensign-restore-fo-merge-write-lazy-loading-reset`).

## Problem

`assertGreetInvokesNoDeferredFOSkill` is an AC-2-sibling behavioral oracle that parses a captured live claude-stream and fails if any pre-greet turn invokes `fo-status-viewer` / `fo-write-core` / `fo-dispatch-recovery` via a `Skill(skill=...)` call. It keys on the *skill-argument string* of the captured stream — a brittle per-skill contract: it must enumerate every deferred FO skill by name, must special-case `present-gate` as the one legitimately-pre-greet skill, and breaks whenever the FO skill set is renamed or re-partitioned. PR #512 restored lazy loading of the FO write/merge cores, and the durable structural guards now cover the greet-and-stop boot without keying on skill-argument strings: `assertNoTeamCreateBeforeGreet` catches an eager team, `assertShallowBootMeasured` catches a missing greet turn. The oracle is therefore redundant and a maintenance liability. #512 drafted its removal in its worktree but never committed or merged it, so `main` still carries the oracle, its two dedicated unit-test controls, and its dedicated negative fixture.

## Proposed approach

Removal only — no new code and no new mechanism. On a fresh branch off `main`, apply the preserved seed patch (`docs/dev/.spacedock-state/remove-greet-deferred-skill-oracle/preserved-lazy-load-wip.patch`), which makes four deletions, all contained to `internal/ensigncycle`:

1. `shallow_boot_measure_test.go` — delete `assertGreetInvokesNoDeferredFOSkill`, the `deferredFOSkillNames` var, and the now-unused `strings` import (its only use was `strings.Contains` inside the removed function).
2. `shallow_boot_measure_unit_test.go` — delete the two dedicated controls `TestAssertGreetInvokesNoDeferredFOSkillOffline` and `TestAssertGreetInvokesNoDeferredFOSkillCatchesLaterDeltaInvoke`.
3. `claude_live_runner_test.go` — delete the AC-2 call site in `runClaudeShallowBootScenario` and reword its doc comment to describe only the retained `assertNoTeamCreateBeforeGreet` signal.
4. `testdata/greet-invokes-fo-status-viewer-skill.stream.jsonl` — delete the negative fixture (referenced only by the removed control).

Containment confirmed by `grep -rn` across `internal/ensigncycle`: every reference to `assertGreetInvokesNoDeferredFOSkill`, `deferredFOSkillNames`, and the fixture lives inside a removed hunk. `greetTurnIndex` is NOT orphaned — after removal it keeps live callers in `assertNoTeamCreateBeforeGreet`, `assertShallowBootMeasuredTurns`, `BuildShallowBootWindowRecord`, and the surviving `TestParserExtractsTeamCallsFromRealHangCapture`. The `strings` import is removed because the deleted function was its sole user (an unused import would fail compilation). RETAINED unchanged: `assertNoTeamCreateBeforeGreet` and `assertShallowBootMeasured`.

Spike (riskiest path, exercised during ideation): applied the seed patch to a throwaway working tree → `git apply` clean, `gofmt -l` clean, `go vet ./internal/ensigncycle/` clean (proves no orphaned symbol and no unused import), fixture deleted, and the 7 retained offline oracle tests green (`TestAssertShallowBootMeasuredOffline`, both `TestAssertNoTeamCreateBeforeGreet*`, `TestParserExtractsTeamCallsFromRealHangCapture`, `TestBuildShallowBootWindowRecord`, `TestShallowBootMeasureSignalsAreIndependent`, `TestPreGreetPeakCacheCreationFindsSpikeNotOnFirstTurn`); then reverted, leaving the tree clean. The riskiest unverified path — "does the package still compile and do the retained guards still pass after the removal" — is proven. No new mechanism is introduced, so there is no enabling-mechanism justification to record. No user-visible surface changes (no CLI output, startup banner, command surface, or docs-site content), so no documentation diff is required.

## Out of scope

- No behavior change to the retained shallow-boot guards (`assertNoTeamCreateBeforeGreet`, `assertShallowBootMeasured`).
- No change to the lazy-loading behavior PR #512 shipped.
- No touching of `greetTurnIndex` (stays used by the surviving functions).
- No changes outside `internal/ensigncycle`.
- No documentation changes — no user-visible surface is affected.

## Acceptance criteria

- **AC-1 (value):** The cumulative test-infrastructure line delta across `internal/ensigncycle/` versus `origin/main` is NEGATIVE — removing the oracle, its two controls, and its fixture nets −101 lines (+2 / −103), with no compensating additions elsewhere in the package. `origin/main`'s line count is the independent baseline; a change that re-added or merely relabeled the oracle instead of removing it would drive the delta non-negative.
  - Verified by: `git diff --numstat origin/main -- internal/ensigncycle/ | awk '{a+=$1;d+=$2} END{print a-d}'` prints a negative integer (−101 for the seed patch).
- **AC-2:** `assertGreetInvokesNoDeferredFOSkill` and the `deferredFOSkillNames` var are absent from the package, and no unused `strings` import remains.
  - Verified by: `grep -rn 'assertGreetInvokesNoDeferredFOSkill\|deferredFOSkillNames' internal/ensigncycle/` prints nothing, and `go vet ./internal/ensigncycle/` is clean (an orphaned symbol or an unused import fails the build).
- **AC-3:** The two dedicated controls and the negative fixture no longer exist.
  - Verified by: `go test ./internal/ensigncycle/ -list 'TestAssertGreetInvokesNoDeferredFOSkill.*'` lists no test, and `test ! -e internal/ensigncycle/testdata/greet-invokes-fo-status-viewer-skill.stream.jsonl`.
- **AC-4:** The durable guards are retained and green. `assertNoTeamCreateBeforeGreet` and `assertShallowBootMeasured` remain present, `greetTurnIndex` remains referenced (not orphaned), and `runClaudeShallowBootScenario` retains only the `assertNoTeamCreateBeforeGreet` check with a comment describing that single signal.
  - Verified by: `go test ./internal/ensigncycle/ -run 'TestAssertShallowBootMeasuredOffline|TestAssertNoTeamCreateBeforeGreet|TestBuildShallowBootWindowRecord|TestShallowBootMeasureSignalsAreIndependent|TestParserExtractsTeamCallsFromRealHangCapture'` — 7 offline tests pass.
- **AC-5:** Package and repo build clean and the full suite is green with the removal in place.
  - Verified by: `gofmt -l internal/ensigncycle/` empty; `go test ./...` and `go test ./... -race` exit 0.

## Test plan

- **Focused internal/ensigncycle tests** (offline, no live API): the 7 retained oracle unit tests in `shallow_boot_measure_unit_test.go` — `TestAssertShallowBootMeasuredOffline`, `TestAssertNoTeamCreateBeforeGreetOffline`, `TestAssertNoTeamCreateBeforeGreetCatchesLaterDeltaTeamCreate`, `TestParserExtractsTeamCallsFromRealHangCapture`, `TestBuildShallowBootWindowRecord`, `TestShallowBootMeasureSignalsAreIndependent`, `TestPreGreetPeakCacheCreationFindsSpikeNotOnFirstTurn`. These prove the retained guards and the turn parser stay green after the oracle and its controls are gone.
- **gofmt / vet:** `gofmt -l internal/ensigncycle/` (empty) and `go vet ./internal/ensigncycle/` (clean — this is what catches the orphaned-symbol / unused-`strings`-import regressions).
- **Package-wide green suite:** `go test ./...` and `go test ./... -race`.
- **Implementation seed:** `docs/dev/.spacedock-state/remove-greet-deferred-skill-oracle/preserved-lazy-load-wip.patch` — the verbatim removal from the #512 worktree, confirmed to apply cleanly to `origin/main` (9aadd89e) and to be the implementation seed for this task.
- **Cost/complexity:** trivial — a 4-file deletion with no new code. Offline focused tests run in <1s; the full `-race` suite is the usual repo cost. No fixture authoring, CLI, or live-workflow test is needed (the removed fixture is deleted, not replaced).

## Stage Report: ideation

- DONE: Formalize acceptance criteria for the oracle removal: at least one value-measuring AC (e.g. a NEGATIVE test-infrastructure line delta from removing the oracle + its two controls + its fixture), each AC naming its test/command verification.
  AC-1 measures the net line delta across `internal/ensigncycle/` vs `origin/main` (NEGATIVE, −101 = +2/−103, verified by `git diff --numstat origin/main | awk`); AC-2..AC-5 each name a grep/`go vet`/`go test -list`/`go test` verification.
- DONE: Test plan names the focused internal/ensigncycle test(s) plus gofmt, `go test ./...`, and `go test ./... -race`; confirm the preserved patch is the implementation seed.
  Test plan lists the 7 retained offline oracle tests by name, `gofmt -l` + `go vet`, `go test ./...` and `-race`, and names the preserved patch as the seed (confirmed applies cleanly to origin/main 9aadd89e).
- DONE: Approach confirms the removal is contained to internal/ensigncycle, RETAINS assertNoTeamCreateBeforeGreet and assertShallowBootMeasured, and neither orphans greetTurnIndex nor leaves an unused `strings` import.
  Approach states containment (grep-confirmed), retention of both guards verbatim, and that `greetTurnIndex` keeps 4 live callers post-removal while `strings`' sole user (the deleted function) is removed with its import.

### Summary

Filled in the ideation body for a pure-removal task: delete the brittle skill-argument-keyed `assertGreetInvokesNoDeferredFOSkill` oracle, its two dedicated controls, and its negative fixture — all contained to `internal/ensigncycle` — while retaining the durable `assertNoTeamCreateBeforeGreet` and `assertShallowBootMeasured` guards. Ran the riskiest-path spike (applied the preserved seed patch to a throwaway tree: `git apply`/gofmt/`go vet` clean and 7 retained offline oracle tests green, then reverted), so the "still compiles + guards still pass after removal" mechanism is proven, not assumed. The value AC is a NEGATIVE −101-line delta vs origin/main; no new mechanism and no user-visible/doc surface is touched.

## Stage Report: implementation

- DONE: Apply the preserved seed patch (docs/dev/.spacedock-state/remove-greet-deferred-skill-oracle/preserved-lazy-load-wip.patch) on this worktree branch: delete assertGreetInvokesNoDeferredFOSkill + deferredFOSkillNames + the now-unused strings import (shallow_boot_measure_test.go), the two controls (shallow_boot_measure_unit_test.go), the AC-2 call site + reword the comment (claude_live_runner_test.go), and the fixture greet-invokes-fo-status-viewer-skill.stream.jsonl — retaining assertNoTeamCreateBeforeGreet and assertShallowBootMeasured. Contained to internal/ensigncycle; do not orphan greetTurnIndex.
  `git apply` clean; commit 536f96e6. AC-2 grep empty; fixture absent; `greetTurnIndex` keeps 4 live callers (shallow_boot_window_record.go + surviving guards); both retained guards present.
- DONE: gofmt -l clean; go vet ./internal/ensigncycle/ clean (catches orphaned symbol / unused import); go test ./... and go test ./... -race exit 0 with the oracle+controls+fixture gone and the durable guards passing.
  `gofmt -l internal/ensigncycle/` empty; `go vet ./internal/ensigncycle/` clean; `go test ./...` all ok; `go test ./... -race` RACE_EXIT=0; AC-4 retained-guard tests pass; AC-3 `-list` shows no removed control.
- DONE: Commit the deliverable on the branch; report the cumulative test-infra line delta vs origin/main (git diff --numstat origin/main -- internal/ensigncycle/ | awk '{a+=$1;d+=$2} END{print a-d}') — must be NEGATIVE (AC-1).
  Committed as 536f96e6 on spacedock-ensign/remove-greet-deferred-skill-oracle; AC-1 delta = -101 (NEGATIVE), matching the seed patch's +2/-103.

### Summary

Applied the preserved seed patch verbatim, removing the brittle skill-argument-keyed `assertGreetInvokesNoDeferredFOSkill` oracle, its `deferredFOSkillNames` var, the now-unused `strings` import, the two dedicated unit-test controls, and the negative fixture — all contained to `internal/ensigncycle`, with `runClaudeShallowBootScenario` retaining only the `assertNoTeamCreateBeforeGreet` check and a reworded comment. All five acceptance criteria verified by exercising: containment grep empty, `go vet`/gofmt clean, retained durable guards green, full `go test ./...` and `-race` both exit 0, and the AC-1 cumulative line delta vs origin/main is -101 (NEGATIVE). Committed as 536f96e6.

## Stage Report: validation

- DONE: Independently reproduce each AC's cited evidence (AC-1 line delta; AC-2 grep+vet; AC-3 controls+fixture absent; AC-4 retained guards green; AC-5 gofmt + go test ./... + -race).
  AC-1: `git diff --numstat origin/main -- internal/ensigncycle/ | awk` = -101 (NEGATIVE, +2/-103). AC-2: grep for both symbols empty, `go vet ./internal/ensigncycle/` exit 0, no `strings` import/use left. AC-3: `go test -list 'TestAssertGreetInvokesNoDeferredFOSkill.*'` lists no test; fixture absent; no dangling filename refs. AC-4: 6/6 cited-regex guard tests PASS (7th green in package suite), both guards live-called in runClaudeShallowBootScenario (L397/L403), greetTurnIndex has 3 live callers. AC-5: `gofmt -l` empty; `go test ./...` exit 0; `go test ./... -race` exit 0; uncached ensigncycle -race exit 0.
- DONE: Semantic adversarial pass — no neutered guard, no dangling reference.
  Both guard bodies byte-unchanged from origin/main (diff removed only the sibling fn/var/import). assertNoTeamCreateBeforeGreet has a passing RED-path control (TestAssertNoTeamCreateBeforeGreetCatchesLaterDeltaTeamCreate) — fires on bad input, not a stub. assertShallowBootMeasured still errors on parse-fail/zero-turns/no-greet (not `return nil`), live-called at L403. greetTurnIndex: 3 live + 1 test caller, not orphaned. Reworded comment matches the two retained calls; no reference to the removed oracle remains. Detached audit judged NOT warranted: pure test-infra DELETION removing an oracle (not weakening one), no production code, no new mechanism — low blast radius; retained RED path verified intact.
- DONE: PASSED/REJECTED recommendation with reproduced evidence per AC; classify findings on defect-kind and release-scope.
  Recommendation: PASSED. All five ACs have valid, independently-reproduced evidence; no outcome defect, no evidence defect, no material finding. Two non-blocking deferred observations below.

### Summary

Validated the pure-removal deliverable (commit 536f96e6, branch spacedock-ensign/remove-greet-deferred-skill-oracle, atop origin/main 9aadd89e). Reproduced all five ACs by exercising, not re-reading: net line delta -101 (NEGATIVE) vs the git baseline, grep-empty + `go vet` clean, no listed control + fixture absent, retained guards present/green/live-called, and full suite + `-race` exit 0. Adversarial pass confirmed neither retained guard was neutered, greetTurnIndex is not orphaned, and the reworded comment matches the single retained signal. Recommendation: PASSED.

Deferred (non-blocking):
- AC-4 prose says "7 offline tests pass" but its cited regex expands to 6 names; the 7th (TestPreGreetPeakCacheCreationFindsSpikeNotOnFirstTurn) passes in the package suite, just isn't matched by the regex. Evidence-annotation nit in the entity text, not a deliverable defect.
- assertShallowBootMeasured has no RED-path unit test in the retained set. Pre-existing, byte-unchanged by this task, explicitly out of scope. Promotes to material only if that guard is later changed without adding a negative test.
