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

### Failure mode 1 (#55297) — residual-flake measurement (riskiest unknown)

The riskiest unknown is whether the path fix + `at`'s shipped `drivePrompt` anti-shutdown override leaves a **residual** flake from upstream #55297. This must be measured by genuine `-count=1` live runs (3–5) with the path fix applied — real model spend, bounded.

**Status: BLOCKED on auth.** The first probe run fast-failed in 2.9s with `apiKeySource:"none"` + `error_status:401 authentication_failed` (`total_cost_usd:0` — no budget burned). The `~/.claude/benchmark-token` (an `sk-ant-oat01` OAuth token) is expired/invalid and re-rotation needs Keychain access this ensign thread lacks. Per the assignment stop condition, escalated to the FO to re-rotate. The measurement (and therefore the final choice between re-promotion vs. documented flake-tolerance) completes once auth is fresh. The path-fix mechanism does not depend on this measurement — it is the deterministic half and is already validated.

## Acceptance criteria

- **AC-1 — Path divergence is closed (deterministic, shippable now).** `spacedock dispatch build` emits an absolute entity path in BOTH the dispatch-body entity-read line and the completion signal, for any `entity_path` input (relative or absolute), so the dispatched ensign resolves the same file regardless of its cwd.
  - **Verified by:** a new offline test in `internal/dispatch/` (the `filepath.Rel`-derived-relative-input pattern above) that fails before the absolutization and passes after — asserting both emitted lines are `filepath.IsAbs`. Plus the existing live stage-report-shape assertions (`liveStageReportHeading` / `doneMarker` / `### Summary`) showing no "two entity files" divergence across the measured live runs.
- **AC-2 — Full end-state assertions restored OR an honest documented flake-tolerance.** The demoted end-state assertions (`liveStageReportHeading`, `doneMarker`, `### Summary`, `checkboxBullet`, `status: done`, `verdict:`, `someCommitNamesOnly`) in `internal/ensigncycle/live_test.go` are re-promoted from non-fatal `t.Logf` to hard `t.Errorf`/`t.Fatalf` AND pass across **N = 5** consecutive genuine `-count=1` runs; OR — if the residual-#55297 measurement shows a flake the override cannot eliminate — a documented flake-tolerance is in place that can NEVER mask a Spacedock regression (see decision below).
  - **Verified by:** the re-promoted hard assertions in `live_test.go` (the demotion `t.Logf` lines with the `#55297` comment are gone, replaced by `t.Errorf`/`t.Fatalf`) and the recorded measurement of N consecutive greens; OR the flake-tolerance construct + the test/comment that documents the upstream tie and the regression-safety argument.

### Residual-#55297 decision (firm rule, pending the measurement)

If N=5 consecutive genuine greens are achievable with the path fix + the override, **re-promote to hard** (simplest, no tolerance) — preferred.

If a residual flake remains, the ONLY acceptable tolerance is a **bounded retry of the whole live cycle inside the test (cap = 3 attempts)** that is HONEST about tolerating a known upstream bug, with these non-negotiable safety properties so it can never mask a Spacedock regression:
- The retry triggers ONLY on the #55297 signature (FO premature teardown / entity left un-terminalized before any end-state assertion ran), identified by the stream signal (no terminal-status marker / un-verdicted entity), NOT on a generic assertion failure. An assertion that runs and fails on a PRESENT-but-WRONG end state (e.g. malformed stage report, wrong commit scope) is a real regression and must NOT be retried — it fails hard immediately.
- The retry is capped and the cap exhaustion is a hard FAIL (never a skip, never a silent pass).
- The tolerance carries a comment naming #55297 and this entity, and is removed when #55297 is fixed upstream (version-pin note: regression in 2.1.126, we run 2.1.161).

Version-pin guidance and "best-effort until upstream" are explicitly rejected as the primary strategy: they leave the assertions demoted (no coverage), which is the status quo this entity exists to end.

## Test plan

- **Offline (no model spend, runs in the default gate):** the `dispatch build` entity-path-absolute contract test above — `go test ./internal/dispatch -run <new test>`. Failing before the fix, green after. This is the AC-1 gate and the deterministic keystone.
- **Live (bounded model spend, gated `-tags live`):** after re-rotation, run 3–5 genuine `go test -tags live -count=1 -run TestLiveEnsignCycle ./internal/ensigncycle` with the path fix applied; record pass rate + failure modes. `-count=1` is mandatory (live tests are silently cacheable). If ≥5 consecutive greens → re-promote (AC-2 path A). If a residual #55297 flake → implement the bounded-retry tolerance (AC-2 path B).
- **Cost/complexity:** offline test is trivial (one fixture + two `IsAbs` asserts, copies an existing pattern). Live measurement is the cost driver but bounded to ~3–5 runs; each clean run is a few minutes of sonnet. The path fix itself is ~6 lines.
- **Sequencing (riskiest-first):** the path fix + offline test are deterministic and validated — they go first and ship the AC-1 keystone. The live residual-flake measurement gates only the AC-2 strategy choice (re-promote vs. tolerance), so it is the last gate, paid once auth is fresh.

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
