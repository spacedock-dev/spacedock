---
title: Codex rejection-flow #141 reviewer-keepalive — diagnose real-gap vs assertion-false-positive, fix the right target
status: implementation
sprint: 0230-stable-finalization
score: 0.5
group: contract
issue:
id: 25q1qfeeae29j6wcxva82wk8
worktree: .worktrees/spacedock-ensign-codex-rejection-flow-keepalive-diagnosis
started: 2026-06-21T21:53:41Z
---

Surfaced by the 0230 flake-landscape investigation: the codex rejection-flow lane fails 6/15 deterministically (all 6 identical) — the codex FO fresh-dispatches the cycle-2 validator instead of reusing the kept-alive cycle-1 reviewer (`assertCodexReviewerReuse`, the #141 keepalive contract). The v0.23.0 tag fires only on a green Runtime Live E2E, and this does NOT re-run-clear, so it must be resolved (fix or scoped) before the cut.

## Problem — TWO candidate root causes; do NOT assume

The investigation graded this a real codex contract gap, but there is genuine tension with a known nuance, so the FIRST job is to distinguish (the M7/M8 discipline):

1. **Real contract gap:** the codex FO genuinely violates the #141 reviewer-keepalive contract — it should reuse the kept-alive cycle-1 validation reviewer for the cycle-2 re-review but fresh-dispatches instead. Fix → harden the codex-FO reuse prose (`## Feedback reviewer reuse` / `«addressable-worker»`).
2. **Assertion false-positive (M7-shape):** if codex's `«addressable-worker»` is legitimately ABSENT on the live codex surface, fresh-dispatch is the CONTRACT-CORRECT behavior — and `assertCodexReviewerReuse` should be taking its `codexAddressableWorkerAbsent → assertCodexFreshValidationWhenAddressableAbsent` branch but isn't (or the absent-detection isn't firing on the live transcript). Fix → the assertion / its absent-detection, not the FO. (M1's implementer flagged exactly this: "the FO judging «addressable-worker» ABSENT on the live Codex surface and fresh-dispatching.")

## Proposed approach

1. **Diagnose first (the deliverable's riskiest step).** Pull the 6 identical failing codex transcripts; grade what the codex FO actually did and WHY (did it observe a turn-starting reuse route / `followup_task` capability, or did it correctly determine no reuse route is exposed and fresh-dispatch?). Read `assertCodexReviewerReuse` + `codexAddressableWorkerAbsent` (`internal/ensigncycle/shared_reviewer_reuse_test.go`) and the codex-FO `## Feedback reviewer reuse` prose. Determine candidate 1 vs 2 with cited evidence — do not assume.
2. **Fix the right target.** Real gap → harden codex-FO reuse prose so the FO reuses when a reuse route IS exposed; keep contractlint green. False-positive → fix the assertion / absent-detection so a contract-correct fresh-dispatch (addressable-worker absent) PASSES while a genuine missed-reuse (route exposed, FO ignored it) still FAILS (non-tautological).
3. **Prove on the live codex lane** N≥3.

## Acceptance criteria

- **AC-1 (diagnosis)** — candidate 1 (real codex contract gap) vs candidate 2 (assertion false-positive) is determined with cited evidence: the failing codex transcripts, `assertCodexReviewerReuse` + `codexAddressableWorkerAbsent`, and the codex-FO reuse prose. The reuse-route-exposed-or-not question is answered from the transcript, not asserted.
- **AC-2 (fix + non-tautology)** — the correct target is fixed; the codex rejection-flow lane passes N≥3 after the change; a GENUINE missed-reuse (a reuse route is exposed and the FO ignores it) still FAILS, and a contract-correct fresh-dispatch (route absent) PASSES.
- **AC-3 (value/no-regression)** — any codex-FO edit keeps codex-FO ≤ its v0.22.0 baseline (6004 B; report the count); `go test ./internal/ensigncycle ./internal/contractlint` + `go build ./...` green.
- **AC-4 (e2e-gate de-risk)** — lands before the cut so the codex rejection-flow lane is not a coin-flip on the tag's green E2E (or, if the captain scopes it non-gating, that decision is recorded here instead).

## Test plan

Live codex rejection-flow N≥3 (rotated codex auth; an auth wall is not a scenario fail — surface to the FO). In-process grade of the captured codex `--json` streams through the assertion entry point. Offline contractlint/ensigncycle/build as the cheap gate.
