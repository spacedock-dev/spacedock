---
id: n1a2q2f3dvxm9ypbrvs6ryeg
title: "Deterministic full multi-agent live cycle — re-promote the scoped TestLiveEnsignCycle end-state assertions"
status: implementation
source: "captain (2026-06-04) — at #285 scoped TestLiveEnsignCycle's live gate to the teardown MARKER grade (deterministic, at's deliverable) and DEMOTED the full multi-agent end-state assertions (verdict / status:done / stage-report-shape / path-scoped-commit) to non-fatal logging, because the real FO+ensign lifecycle is non-deterministic. This entity tracks making that lifecycle deterministic enough to re-promote those assertions to hard."
score: "0.35"
started: 2026-06-04T06:30:55Z
completed:
verdict:
issue:
worktree: .worktrees/spacedock-ensign-live-cycle-end-state-determinism
---

`at` (#285) unblocked the path to TestLiveEnsignCycle's full-lifecycle end-state assertions (by fixing the teardown hang), and that exposed that a single `claude -p` FO + a dispatched ensign cannot deterministically complete a full multi-stage lifecycle to a verdicted, archived terminal state. `at` shipped with the live gate SCOPED (teardown marker grade hard; end-state non-fatal/logged) so it wasn't blocked on an inherently-flaky assertion. This entity is the follow-up: make the full cycle deterministic, then re-promote the demoted assertions.

## Problem

TestLiveEnsignCycle drives a REAL FO (`claude -p` + Agent Teams) that dispatches a REAL ensign through `backlog → implementation → done`. Each step is a non-deterministic model decision (dispatch, drive-to-terminal, set a verdict, write the stage report to the right file, terminalize, teardown). Across genuine `-count=1` runs the cycle is flaky (observed ~2/3 PASS), from at least TWO independent failure modes:

1. **Upstream claude-code #55297 (premature shutdown).** In `claude -p` with an active team, the harness injects a per-turn `<system-reminder>` ("you cannot return a response until your team is shut down… shut down before your final response"; regression in 2.1.126, we run 2.1.161). The FO reads it as "shut down now" and tears down BEFORE finishing the workflow → entity never reaches a verdicted terminal. **Mitigation already shipped in `at`:** a generic anti-early-shutdown clause in the live `-p` launch (`drivePrompt`) — it raised the drive-to-terminal rate substantially but does NOT eliminate the flake (you cannot fully out-prose a per-turn harness reminder). Upstream: https://github.com/anthropics/claude-code/issues/55297 (also related: #38116 / #57681 — an approved-shutdown member is never cleared from `members[]`, the basis for `at`'s launcher-owned-exit design).
2. **Ensign stage-report path divergence.** In a failing genuine run the FO DID drive to a verdict (`status=implementation started` → `status=done` → `completed verdict=done`), but the test still red'd on `:219/:222/:225` (stage-report heading / `- DONE:` marker / `### Summary`) because the dispatched ensign wrote its stage report to a DIFFERENT `make-it-work.md` than the FO tracked (FO thinking: "I see two different entity files…"). Likely a dispatch/path-handling inconsistency in how the live cycle hands the ensign the entity path — possibly fixable to determinism (unlike #55297).

The teardown MARKER grade (`expectTerminalTeardownGrade`) and the TeamCreate→dispatch-close steps passed reliably across all runs — that is `at`'s deliverable and stays the hard assertion.

3. **#55297 also hits the SHARED scenarios (not just the full cycle).** CI run `26926169461` (PR #285 @ `f6eb6bfc`) red'd `claude-live (sonnet)` on `TestLiveClaudeSharedScenarios/gate-guardrail` — `claude_live_runner_test.go:49: spacedock claude did not finish within 1m0s` (the FO HUNG in #55297's reminder loop). The shared scenarios launch via `claude_live_runner_test.go`'s `run` with the bare canonical prompt and NO override. **Captain directive: apply the #55297 `-p` override to EVERY `-p`+team live launch** (the shared runner + the full cycle). Even with the override, `gate-guardrail`'s 60s timeout may need lengthening or retry-tolerance (same residual non-determinism) — covered by this entity.

## What `at` left in place (the scoped state to build on)

- `internal/ensigncycle/live_test.go`: teardown marker grade is HARD (`t.Fatalf`); the end-state assertions (`liveStageReportHeading`, `doneMarker`, `### Summary`, `checkboxBullet`, `status: done`, `verdict:`, `someCommitNamesOnly`) are NON-FATAL logging, with a comment referencing #55297 + this entity.
- The 3-stage `readmeRealisticLifecycle` fixture, the canonical neutral prompt (`Use $spacedock:first-officer for this whole run.`), the `drivePrompt` #55297 override, and `-count=1` on the live invocation (CI + README) all shipped with `at`.
- The `status --set` verdict gate (refuse a verdict-less terminal finalize) shipped and is proven offline + audit — independent of this live flakiness.

## Proposed scope (for ideation)

1. Fix the deterministic-izable flake (the ensign stage-report path divergence): make the live dispatch hand the ensign an unambiguous absolute entity path and confirm the ensign writes its stage report there. This is the most likely path to determinism that does NOT depend on the upstream fix.
2. Track claude-code #55297 upstream; decide the residual strategy if the override + the path fix still leave a flake (e.g., a retry-with-cap inside the test that is HONEST about being a flake-tolerance for a known upstream bug — NOT flake-reliance on a real regression; or version-pin guidance; or keep best-effort until #55297 is fixed).
3. Once the full cycle is deterministic (measured: N consecutive genuine `-count=1` greens), RE-PROMOTE the demoted end-state assertions from non-fatal logging back to hard (`t.Errorf`/`t.Fatalf`), restoring the full-lifecycle coverage.

## Design (ideation)

### Root cause of failure mode 2 (path divergence) — pinned

`spacedock dispatch build` (`internal/dispatch/build.go:236`) reads `entity_path` from the request and echoes it **verbatim** into the two places the ensign acts on it:
- the entity-read line — `Read the entity file at <entity_path> for the current spec.` (`build.go:496`)
- the completion signal — `Report written to <entity_path>.` (`build.go:539` via `completionSignalBlock`)

It does **not** absolutize `entity_path`, even though it explicitly absolutizes `workflow_dir` a few lines later (`build.go:291`, with a comment about cwd-independence for the worktree worker). Asymmetry confirmed by reading the code.

The FO supplies `entity_path` itself: `spacedock status --next --json` / `--where … --json` emit `{id, slug, current, next, worktree}` with **no path field** (verified by running both against a staged fixture). The FO contract (`claude-first-officer-runtime.md:78`, `--entity-path {entity_file_path}`) never says absolute-vs-relative, and the FO runs with cwd = workflow root. So the FO naturally passes a **relative** path (e.g. `make-it-work.md`). `build`'s `isFile(entityPath)` readability check (`build.go:281`) passes because build runs in the FO's cwd — masking the defect. The dispatched ensign is a separate Agent whose cwd is **not** guaranteed to be the workflow root (in Claude teams the agent thread's cwd resets between bash calls and is not pinned to root). It resolves the relative path against its own cwd → a different / nonexistent `make-it-work.md` → its stage report lands in the wrong file. This is the "two make-it-work.md" / "two entity files" divergence `at` observed (~1/3 of runs).

### Fix (deterministic, NO upstream dependency)

Absolutize `entity_path` against the process cwd inside `runBuild`, immediately after the Rule 10 readability check, exactly mirroring the existing `workflow_dir` absolutization:

```go
if abs, err := filepath.Abs(entityPath); err == nil {
    entityPath = abs
}
```

This is the single chokepoint every dispatch flows through, so the fix covers all hosts and every FO regardless of the cwd it passed — strictly better than adding prose guidance to the FO contract (which a model can ignore). Absolutizing AFTER the readability error message keeps that diagnostic showing the FO's original spelling.

**Spike — offline-validated (`go build` + `dispatch build` with a relative `entity_path`):** before the fix the body emits `Read the entity file at make-it-work.md` / `Report written to make-it-work.md` (relative); after the fix it emits `…/tmp/sdtest/make-it-work.md` (absolute) for both lines. The mechanism works.

### Offline contract test (pins AC-1)

Mirror the proven precedent `TestStateCommitPathAbsoluteFromRelativeWorkflowDir` (`internal/dispatch/build_statecommit_test.go:19`), which pins the same absolutization machinery for `workflow_dir`. Because `t.Chdir` is unavailable (module targets go1.22), derive a **relative** `entity_path` spelling via `filepath.Rel(cwd, entityPath)`, pass it to `dispatch build`, and assert the emitted dispatch body's entity-read line AND completion-signal line both carry an **absolute** path (`filepath.IsAbs`), not the relative spelling. Belongs in `internal/dispatch/` beside the existing entity-path-contract tests (`build_hazards_test.go`'s space-path test, `build_statecommit_test.go`). This is a real failing test before the fix (the body carries the relative spelling) and green after — TDD-ready for implementation.

### Residual-flake measurement (riskiest unknown) — DONE; found a SECOND deterministic bug, not #55297

After re-rotation, ran 4 genuine `-count=1` live cycles with the path fix applied (sonnet, `-count=1`, fresh `t.TempDir` per run):

| Run | Result | End-state non-fatal misses | #55297 fired? | Stale-dispatch-file leak? |
|-----|--------|---------------------------|---------------|---------------------------|
| A | PASS, archived, all 7 end-state checks green | 0 | no (verdicted terminal) | no (first run, clean `/tmp/spacedock-dispatch`) |
| B | PASS (hard), **3 stage-report-shape misses** | heading + `- DONE:` + `### Summary` | no (verdicted terminal) | **YES** — ensign read a prior run's tempdir |
| C | PASS, archived, all 7 green | 0 | no | no (ran after `rm /tmp/spacedock-dispatch/…make-it-work-*.md`) |
| D | PASS, archived, all 7 green | 0 | no | no (dirty dir, but won the race) |

**Headline: #55297 did NOT fire in any of the 4 runs** — every cycle reached a verdicted, archived terminal (`status: done` + `verdict:` set). `at`'s `drivePrompt` override is holding for the full cycle. The flake we measured (1/4) is a **SECOND, independent, deterministic-izable bug**, not the upstream shutdown bug.

**Second root cause — dispatch-file path collision (global, deterministic name).** `dispatch build` writes the dispatch body to a hardcoded global path `/tmp/spacedock-dispatch/{derived_name}.md` (`build.go:23` `dispatchFileDir = "/tmp/spacedock-dispatch"`, `build.go:549`), and `derived_name = {worker_key}-{slug}-{stage}` is identical across every run of the same fixture (`spacedock-ensign-make-it-work-implementation`). So the dispatch files PERSIST and COLLIDE across runs. In run B the FO's own thinking captured it verbatim: *"the ensign reported writing to `…/TestLiveEnsignCycle112126198/…/make-it-work.md` (note: `112126198` vs `1287896541`)… There are two separate test instances… The ensign worked on the wrong entity file."* The ensign Read a STALE dispatch-file pointer (from a prior run's tempdir `112126198`, already cleaned up) instead of run B's path (`1287896541`), did its work correctly, but to the dead tempdir → run B's tracked entity archived without the report. Confirmed by isolation: clearing `/tmp/spacedock-dispatch` before run C made it green; the failure is a RACE (whether the FO's `dispatch build` overwrite lands before the ensign reads), so it is probabilistic when stale same-named files exist (runs A/C/D won the race, B lost it).

This bug is independent of the entity-path-absolutization fix (run B's dispatch body DID carry an absolute path — just the wrong, stale one). It is the SAME observable symptom (`at`'s "two make-it-work.md" / "two entity files") with a DIFFERENT root cause. Both are deterministic-izable.

### Second fix — make the dispatch-file path collision-free

The dispatch file must not collide across concurrent FOs or back-to-back runs sharing an entity slug+stage. Two viable approaches (decide at implementation):
- **(preferred) Make the dispatch-file path unique per dispatch** — include a per-dispatch nonce / the team_name / a content hash in the filename (e.g. `/tmp/spacedock-dispatch/{team_name}-{derived_name}.md` or `…/{derived_name}-{shorthash}.md`). The FO already forwards `dispatch_file_path` verbatim from build output into the Agent prompt, so a unique path flows through with zero FO-contract change. This closes BOTH the cross-run test flake AND the real-world hazard of two concurrent FOs on one machine working same-named entities.
- **(insufficient alone) Test-only `TMPDIR` isolation** — pointing `dispatchFileDir` at a per-run temp dir would fix the live TEST but not the real product hazard (concurrent FOs). Use only as a belt-and-braces ADDITION, not the primary fix.

This second fix is the new keystone for AC-1's "no two-entity-files divergence" — the absolutization fix alone does not deliver it (run B proves that). Both fixes ship together.

### Failure mode 1 (#55297) — measured clean across this sample

Across these 4 runs, #55297 did not fire; the override held. This does NOT prove #55297 is gone (it is an upstream timing bug and a larger sample could still catch it), so AC-2 keeps the regression-safe bounded-retry tolerance as the documented fallback — but the EVIDENCE says the dispatch-file collision, not #55297, was the live test's actual residual flake.

## Acceptance criteria

- **AC-1 — The "two entity files" divergence is closed (deterministic, shippable now), via BOTH fixes.**
  - **(1a) Absolute entity path.** `spacedock dispatch build` emits an absolute entity path in BOTH the dispatch-body entity-read line and the completion signal, for any `entity_path` input (relative or absolute), so the dispatched ensign resolves the same file regardless of its cwd.
  - **(1b) Collision-free dispatch-file path.** `spacedock dispatch build` writes the dispatch body to a path that does NOT collide across concurrent FOs or back-to-back dispatches sharing an entity slug+stage, so an ensign never reads a stale prior dispatch's pointer.
  - **Verified by:** (1a) a new offline test in `internal/dispatch/` (the `filepath.Rel`-derived-relative-input pattern above) that fails before the absolutization and passes after — asserting both emitted lines are `filepath.IsAbs`. (1b) an offline test asserting two `dispatch build` invocations for the same slug+stage but different entity tempdirs produce DISTINCT `dispatch_file_path` values (and the first file's content is not clobbered/aliased by the second). Plus the live stage-report-shape assertions (`liveStageReportHeading` / `doneMarker` / `### Summary`) green across the measured runs (the 1/4 miss in this sample was bug 1b, not #55297).
- **AC-2 — Full end-state assertions restored OR an honest documented flake-tolerance.** The demoted end-state assertions (`liveStageReportHeading`, `doneMarker`, `### Summary`, `checkboxBullet`, `status: done`, `verdict:`, `someCommitNamesOnly`) in `internal/ensigncycle/live_test.go` are re-promoted from non-fatal `t.Logf` to hard `t.Errorf`/`t.Fatalf` AND pass across **N = 5** consecutive genuine `-count=1` runs WITH both AC-1 fixes applied; OR — if a residual flake the fixes cannot eliminate remains (#55297 or other) — a documented flake-tolerance is in place that can NEVER mask a Spacedock regression (see decision below).
  - **Verified by:** the re-promoted hard assertions in `live_test.go` (the demotion `t.Logf` lines with the `#55297` comment are gone, replaced by `t.Errorf`/`t.Fatalf`) and the recorded measurement of N consecutive greens; OR the flake-tolerance construct + the test/comment that documents the upstream tie and the regression-safety argument.

### Re-promotion decision (firm rule, informed by the measurement)

The measurement shows the live residual flake was the dispatch-file collision (bug 1b), and #55297 did not fire across 4 runs. So the expected path is: ship BOTH AC-1 fixes, then **re-promote to hard** once N=5 consecutive genuine greens hold (simplest, no tolerance) — preferred.

If a residual flake STILL remains after both fixes (e.g. #55297 surfaces in a larger sample), the ONLY acceptable tolerance is a **bounded retry of the whole live cycle inside the test (cap = 3 attempts)** that is HONEST about tolerating a known upstream bug, with these non-negotiable safety properties so it can never mask a Spacedock regression:
- The retry triggers ONLY on the #55297 signature (FO premature teardown / entity left un-terminalized before any end-state assertion ran), identified by the stream signal (no terminal-status marker / un-verdicted entity), NOT on a generic assertion failure. An assertion that runs and fails on a PRESENT-but-WRONG end state (e.g. malformed stage report, wrong commit scope) is a real regression and must NOT be retried — it fails hard immediately. NOTE: bug 1b is NOT a valid retry trigger — it is a Spacedock bug, fixed by AC-1(1b), and retrying it would mask the regression. The retry covers upstream #55297 ONLY.
- The retry is capped and the cap exhaustion is a hard FAIL (never a skip, never a silent pass).
- The tolerance carries a comment naming #55297 and this entity, and is removed when #55297 is fixed upstream (version-pin note: regression in 2.1.126, we run 2.1.161).

Version-pin guidance and "best-effort until upstream" are explicitly rejected as the primary strategy: they leave the assertions demoted (no coverage), which is the status quo this entity exists to end.

## Test plan

- **Offline (no model spend, runs in the default gate):** two `internal/dispatch/` tests — (1a) entity-path-absolute contract (`filepath.Rel`-derived relative input → both emitted lines `filepath.IsAbs`); (1b) dispatch-file path uniqueness (two builds, same slug+stage, different entity tempdirs → distinct `dispatch_file_path`, no clobber). Both failing before their fix, green after. These are the AC-1 gates and the deterministic keystones.
- **Live (bounded model spend, gated `-tags live`):** measurement DONE this stage — 4 genuine `-count=1` runs with the absolutization fix: A/C/D green, B flaked on bug 1b (stale dispatch-file collision), #55297 never fired. At implementation, after BOTH fixes land, re-run for N=5 consecutive greens to gate re-promotion. `-count=1` is mandatory (live tests are silently cacheable). If ≥5 consecutive greens → re-promote (AC-2 path A). If a residual #55297 flake surfaces → bounded-retry tolerance (AC-2 path B).
- **Cost/complexity:** both offline tests are trivial (fixtures + a couple asserts, copy existing patterns). The absolutization fix is ~6 lines; the dispatch-file-path fix is a small filename change at `build.go:549` + its docstring. Live re-run at implementation is the cost driver, bounded to ~5 runs (~25–35 min sonnet total).
- **Sequencing (riskiest-first, validated):** the riskiest unknown (does the override leave a residual flake?) is now ANSWERED by the measurement — the residual was bug 1b, not #55297. Both offline fixes are deterministic and go first; the live N=5 re-run is the last gate at implementation.

## Notes / provenance

- The full root-cause trail (8+ live runs, the override discovery, the cache false-green catch, the two failure modes) is in `at` (`non-interactive-teardown-exit`) feedback cycle 3 + its stage reports, and in `/tmp/fo-live-run-*.log` (session-local).
- Do NOT re-litigate `at`'s shipped mechanism (teardown relax + verdict gate + override + scoping). This entity is purely about restoring full-lifecycle determinism + re-promoting the assertions.

## Stage Report: ideation

- DONE: Pin the deterministic-izable flake (the ensign stage-report PATH DIVERGENCE)
  Root cause pinned to `internal/dispatch/build.go:236` — `entity_path` is echoed verbatim into the body's entity-read line (`:496`) + completion signal (`:539`) but never absolutized, while `workflow_dir` IS (`:291`); FO derives a relative path (`status --next/--where` emit no path, verified) and the ensign's differing cwd resolves a different file. Fix designed (6-line `filepath.Abs` mirror) + offline test designed (`filepath.Rel`-derived relative input, mirroring `build_statecommit_test.go:19`).
- DONE: Exercise the riskiest unknown FIRST (residual flake after path fix + #55297 override)
  Spike build with the path fix applied confirmed OFFLINE that a relative `entity_path` now emits an absolute path in both body lines. The LIVE residual-flake measurement is BLOCKED: first probe run fast-failed in 2.9s with `apiKeySource:"none"` + 401 `authentication_failed` (`total_cost_usd:0`). Escalated to FO to re-rotate `~/.claude/benchmark-token` per the stop condition (I lack Keychain access).
- DONE: Decide the residual-#55297 strategy + re-promotion plan; firm AC-1/AC-2
  AC-1 firmed (absolute entity path in both emitted lines; gated by the offline contract test). AC-2 firmed with N=5 consecutive genuine `-count=1` greens to re-promote; the only acceptable residual tolerance is a bounded (cap=3) whole-cycle retry that triggers ONLY on the #55297 signature (never on a present-but-wrong end state), fails hard on cap exhaustion, and is documented + removable on upstream fix. Version-pin / best-effort rejected as primary.

### Summary

Root-caused failure mode 2 (path divergence) to `dispatch build` echoing a relative `entity_path` verbatim instead of absolutizing it like `workflow_dir`, and offline-validated the fix via a spike build. Designed the deterministic fix + offline contract test, firmed AC-1 (absolute path, offline-gated) and AC-2 (re-promote at N=5 greens, or a regression-safe bounded retry tolerance). The live residual-#55297 measurement — the one remaining riskiest unknown — is blocked on an expired benchmark-token (escalated to FO); the deterministic half of the work does not depend on it.

## Stage Report: ideation (cycle 2 — live measurement after token re-rotation)

- DONE: Pin the deterministic-izable flake (the ensign stage-report PATH DIVERGENCE)
  Found TWO independent deterministic-izable root causes for the same "two entity files" symptom: (a) `dispatch build` echoes a relative `entity_path` verbatim (`build.go:236`) instead of absolutizing it like `workflow_dir` (`build.go:291`); (b) `dispatch build` writes to a GLOBAL collision-prone path `/tmp/spacedock-dispatch/{derived_name}.md` (`build.go:23/:549`) that is identical across runs of the same slug+stage, so an ensign can Read a STALE prior-run pointer. Both fixes designed; both offline tests designed.
- DONE: Exercise the riskiest unknown FIRST (residual flake after path fix + #55297 override)
  Ran 4 genuine `-count=1` live cycles (sonnet) with the absolutization fix applied: A/C/D fully green (all 7 end-state checks), B flaked with 3 stage-report-shape misses. Isolated B's cause to bug (b) via the FO's own transcript ("two separate test instances… ensign worked on the wrong entity file", reading prior tempdir `112126198` vs B's `1287896541`) and a control (clearing `/tmp/spacedock-dispatch` → run C green). KEY: #55297 fired in ZERO of 4 runs — every cycle reached a verdicted archived terminal; `at`'s override is holding for the full cycle. The live residual flake is bug (b), not #55297.
- DONE: Decide the residual strategy + re-promotion plan; firm AC-1/AC-2
  AC-1 expanded to require BOTH fixes (absolute path 1a + collision-free dispatch-file path 1b) with an offline test each. AC-2 re-promotion at N=5 consecutive greens WITH both fixes; bounded-retry tolerance retained for #55297 ONLY (explicitly NOT for bug 1b, which is a Spacedock bug to fix). Preferred path: ship both fixes, re-promote.

### Summary

Token re-rotated; ran the bounded live measurement and it overturned the prior framing: #55297 did not fire across 4 runs (the override holds), and the residual live flake is a SECOND deterministic-izable Spacedock bug — `dispatch build`'s global collision-prone dispatch-file path `/tmp/spacedock-dispatch/{name}.md` lets an ensign read a stale prior-run entity pointer (proven by the FO's own transcript + a clear-dir control). The plan now ships TWO offline-gated fixes (absolute entity path + collision-free dispatch-file path); re-promotion at N=5 greens is the expected path, with the #55297 bounded-retry kept only as a documented fallback. Spike worktree + test dispatch files cleaned up.

## Stage Report: implementation

- DONE: AC-1(1a) — absolutize `entity_path` in `dispatch build` so both emitted lines are absolute
  `internal/dispatch/build.go` absolutizes `entityPath` (filepath.Abs mirror of the `:291` workflow_dir block) after the Rule-10 readability check; offline contract test `TestEntityPathAbsoluteFromRelativeInput` (filepath.Rel-derived relative input → both the entity-read line AND the completion signal assert `filepath.IsAbs`) was RED before, GREEN after. Committed da579260 on `spacedock-ensign/live-cycle-end-state-determinism`.
- DONE: AC-2 — measure N=5 and re-promote per the FIRM rule
  Measured pass rate: **5/5 consecutive genuine `-count=1` TestLiveEnsignCycle runs GREEN with the end-state assertions HARD** (sonnet; SPACEDOCK_BIN/SPACEDOCK_REPO_ROOT exported; `-count=1`; fresh `t.TempDir` per run; archived terminal + all 7 end-state checks held, zero `t.Errorf` lines). Plus 2 prior confirmation runs (run 1 hard-teardown-gate PASS; run 2 non-fatal with ZERO non-fatal misses = end-state held) = 7 clean runs total. **#55297 did NOT fire in any run** (every cycle reached a verdicted archived terminal). Per the FIRM decision rule, 5/5 GREEN → re-promote (PREFERRED); NO bounded-retry tolerance installed (no residual flake to tolerate). `internal/ensigncycle/live_test.go` end-state checks restored from non-fatal `t.Logf` to hard `t.Errorf`, `#55297` demotion comments dropped, checks gate on a PRESENT-and-CORRECT end state (a present-but-wrong end state fails immediately, never retried). `go vet -tags live` clean. Committed 7c3a1e61.
- DONE: Full offline `go test ./...` green; `go vet` clean
  1003 offline tests pass across 12 packages; `go vet ./...` and `go vet -tags live ./...` both clean (with the committed 1a + re-promotion; the stashed 1b changes excluded — see FAILED below).
- FAILED: AC-1(1b) — collision-free dispatch-file path (BLOCKED on a parity-contract decision; escalated)
  The current entity spec (ideation cycle 2, written AFTER my dispatch was packaged) expanded AC-1 to require a SECOND fix — a collision-free `/tmp/spacedock-dispatch/{derivedName}.md` so an ensign never reads a stale prior dispatch's pointer. My dispatch checklist did NOT include 1b. I implemented it TDD (RED→GREEN): `TestDispatchFilePathCollisionFree` (two builds, same slug+stage, different entity tempdirs → distinct `dispatch_file_path`, no clobber) + a content-hash suffix on the dispatch filename (`{derivedName}-{sha256(body)[:12]}.md`, keeping `derivedName` as the readable team-member name). BLOCKER: the hash diverges the frozen GOLDEN parity contract — `dispatch_file_path` and the pointer `prompt` are frozen in `internal/dispatch/testdata/golden/*.txt`, and `normPaths` only normalizes root/home, not the hash; worse, the body embeds the per-run absolute tempdir so the hash is non-deterministic across runs and cannot be frozen even after `-update`. This breaks 5 golden parity tests (TestBuildParityCrossProduct, TestBuildAbsoluteWorktreeParity, TestBuildSchemaVersionNumericParity, TestBuildModsParity, TestBuildNonAsciiTitleParity). Resolving it is an architectural decision on the golden parity contract (regenerate goldens + normalize the hash token, OR use a team_name-keyed filename instead of a content hash, OR defer 1b) on a HIGH-STAKES file. The 1b work is preserved in worktree `git stash@{0}` ("wip-1b-dispatch-file-collision-fix"); the committed tree is green. Escalated to the FO for the parity-contract decision rather than unilaterally re-architecting it.

### Summary

Shipped AC-1(1a) (absolute entity_path + offline contract test) and AC-2 (re-promoted the demoted end-state assertions to hard after measuring 5/5 consecutive genuine `-count=1` live greens — #55297 never fired, no tolerance needed, the PREFERRED path). The measurement ran against the SAME persisted global dispatch file across all 7 runs and won the race every time, consistent with the prior 3/4 win rate — so 5/5 green does NOT prove the dispatch-file collision (bug 1b) is absent, only that it is winnable. AC-1(1b) — required by the entity spec but NOT in my dispatch checklist — is implemented and TDD-validated locally but BLOCKED: the content-hash filename breaks the frozen-golden parity contract (the hash is non-deterministic across runs and is not normalized), which is an architectural decision on a high-stakes file. Escalated that decision to the FO; the 1b change is stashed and the committed tree is green.

## Stage Report: implementation (cycle 2 — AC-1(1b) team_name-keyed fix + post-1b N=5)

- DONE: AC-1(1b) — collision-free dispatch-file path
  Resolved the escalation per the FO's PREFERRED direction (team_name keying, not the content hash). `internal/dispatch/build.go` now writes the dispatch body to `/tmp/spacedock-dispatch/{teamName}-{derivedName}.md` when a team name is present, falling back to `{derivedName}.md` in bare mode (bare dispatches block until the subagent completes and cannot alias within a session). Each real FO TeamCreate yields a unique team name (`{project}-{dir}-{YYYYMMDD-HHMM}-{shortuuid}`), so back-to-back runs and concurrent FOs never alias — confirmed live: run produced `002-002-20260604-0736-f9506ff9-spacedock-ensign-make-it-work-implementation.md` instead of the bare colliding name. `derivedName` stays the readable team-member name. Offline test `TestDispatchFilePathCollisionFreeAcrossTeams` (two teams, same slug+stage → distinct `dispatch_file_path`, no clobber): RED before, GREEN after. Goldens regenerated (deterministic `fixture-team-` prefix; diff is ONLY the prefix on the 2 dispatch-path lines across 18 goldens, audited); `build_errors_test.go`'s write-failed fixture updated to the team-keyed path. team_name keying chosen over a content hash precisely because the hash embeds the per-run absolute tempdir → non-deterministic → un-freezable in goldens. Committed 9c6b7ae7.
- DONE: AC-2 (cycle 2) — N=5 consecutive genuine greens WITH BOTH fixes
  Rebuilt the binary with both 1a+1b, then ran **5/5 consecutive genuine `-count=1` TestLiveEnsignCycle runs GREEN** (sonnet; SPACEDOCK_BIN/SPACEDOCK_REPO_ROOT exported; `-count=1`; fresh `t.TempDir` per run; `/tmp/spacedock-dispatch` deliberately NOT cleared between runs — the 1b fix makes stale files unreadable regardless). All 5 archived the terminal entity with all 7 end-state checks HARD-green (zero `t.Errorf`/`live cycle failed` lines); durations 249–342s. **#55297 fired in 0/5** (every cycle reached a verdicted archived terminal; `at`'s drivePrompt override holds). Per the FIRM rule, 5/5 GREEN with both fixes → re-promote (PREFERRED); NO bounded-retry tolerance installed. The #55297 bounded-retry remains a documented design fallback only and does NOT cover bug 1b (a Spacedock bug, now fixed). The re-promoted hard assertions (committed 7c3a1e61) stand unchanged.
- DONE: Full offline `go test ./...` green; `go vet` clean
  1004 offline tests pass across 12 packages (the +1 over cycle 1 is the new 1b contract test); `go vet ./...` and `go vet -tags live ./...` both clean. Three commits on `spacedock-ensign/live-cycle-end-state-determinism`: da579260 (1a), 7c3a1e61 (re-promotion), 9c6b7ae7 (1b). Tree clean (only the untracked build binary).

### Summary (cycle 2)

Closed the AC-1(1b) escalation by keying the dispatch-file path on team_name (the FO's preferred approach), which is collision-free for concurrent FOs / back-to-back runs AND deterministic in the goldens — unlike the content hash, which embeds per-run absolute paths and can't be frozen. Both AC-1 fixes (1a absolute entity_path + 1b team-keyed dispatch file) and the re-promoted hard end-state assertions are committed; the post-both-fixes confirmation is 5/5 consecutive genuine `-count=1` live greens with #55297 firing 0/5, satisfying AC-2's PREFERRED re-promote path with no tolerance. Full offline `go test ./...` (1004) + `go vet` (offline and live tag) all green.
