---
id: n1a2q2f3dvxm9ypbrvs6ryeg
title: "Deterministic full multi-agent live cycle — re-promote the scoped TestLiveEnsignCycle end-state assertions"
status: ideation
source: "captain (2026-06-04) — at #285 scoped TestLiveEnsignCycle's live gate to the teardown MARKER grade (deterministic, at's deliverable) and DEMOTED the full multi-agent end-state assertions (verdict / status:done / stage-report-shape / path-scoped-commit) to non-fatal logging, because the real FO+ensign lifecycle is non-deterministic. This entity tracks making that lifecycle deterministic enough to re-promote those assertions to hard."
score: "0.35"
started: 2026-06-04T06:30:55Z
completed:
verdict:
issue:
worktree:
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

## Acceptance criteria (seed — to be firmed at ideation)

- **AC-1 (seed):** The dispatched ensign's stage report lands in the FO-tracked entity file deterministically — verified by repeated genuine `-count=1` live runs showing the stage-report-shape assertions pass (no "two entity files" divergence), and an offline test pinning the dispatch entity-path contract.
- **AC-2 (seed):** The full TestLiveEnsignCycle end-state assertions are re-promoted to hard and pass across N consecutive genuine `-count=1` runs (N to be set at ideation), OR — if blocked on #55297 — an honest, documented flake-tolerance is in place that does not mask a Spacedock regression.

## Test plan (seed)

- Reproduce the path divergence deterministically (it appeared ~1/3); add an offline test for the dispatch entity-path contract; fix; confirm with genuine `-count=1` live runs.
- The riskiest unknown is whether #55297 leaves a residual flake even after the path fix — exercise that first (several genuine `-count=1` runs with the path fix) before committing to re-promoting the assertions.

## Notes / provenance

- The full root-cause trail (8+ live runs, the override discovery, the cache false-green catch, the two failure modes) is in `at` (`non-interactive-teardown-exit`) feedback cycle 3 + its stage reports, and in `/tmp/fo-live-run-*.log` (session-local).
- Do NOT re-litigate `at`'s shipped mechanism (teardown relax + verdict gate + override + scoping). This entity is purely about restoring full-lifecycle determinism + re-promoting the assertions.
