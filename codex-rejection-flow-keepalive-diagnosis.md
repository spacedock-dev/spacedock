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

## Stage Report: implementation

- DONE: DIAGNOSE FIRST — candidate 1 (real codex contract gap) vs candidate 2 (assertion false-positive)
  Candidate 2 (assertion absent-detection false-positive). Graded 14 captured live transcripts in-process through `assertCodexReviewerReuse`: 6 PASS, 8 FAIL. Of the 8, 6 are the identical deterministic failure; 2 are truncated/stalled (watchdog-killed, not contract end-states). For ALL 6 identical failures the codex FO PROBED the live surface, found only `spawn_agent`+`wait_agent` (NO reuse route), correctly declared `«addressable-worker»` ABSENT in natural language, and contract-correctly fresh-dispatched a separate cycle-2 reviewer. The reuse route IS non-deterministically exposed — 6 OTHER transcripts reuse via `send_input` (PRESENT branch). So fresh-dispatch when the route is absent is CONTRACT-CORRECT; the FO is not violating #141.
- DONE: Reuse-route-exposed-or-not answered from the transcript, not asserted
  Read from the FO's own `agent_message` narration: PRESENT runs say "the live surface has spawn_agent, wait_agent, send_input ... so the validation reviewer can be kept addressable and reused"; ABSENT runs say "spawn_agent and wait_agent, but no follow-up/send-message binding" / "the cycle-1 reviewer is not addressable on this host". The route's presence is non-deterministic across runs.
- DONE: Root cause located — `codexAddressableWorkerAbsent` matched a fixed list of 5 canned phrases
  The FO's real absence wording matched NONE of the 5, so the assertion fell through to its PRESENT branch and RED-flagged a contract-correct fresh-dispatch (the exact M7-shape false-positive the M1 implementer flagged).
- DONE: Fix the RIGHT target (the assertion, not the FO) + non-tautology
  Replaced the brittle phrase list with a structural detector reading only `agent_message` narration (never command output, which echoes the runtime doc): a message that negates a reuse-route concept (follow-up / followup_task / addressable / reuse) is the absence observation. A reuse tool actually present still overrides absent via `assertCodexFreshValidationWhenAddressableAbsent`. Commit 6e426526 (worktree branch). NON-TAUTOLOGICAL, pinned by existing + new table guards: genuine missed-reuse (route present via spawn, FO ignores it → `freshDispatch`) still FAILS; `absentButReuseTool` still FAILS; `absentOnlyOneValidation` still FAILS; contract-correct fresh-dispatch (route absent) PASSES across 10 verbatim live-FO wordings.
- DONE: codex rejection-flow passes N≥3 after the change; genuine missed-reuse still FAILS, contract-correct fresh-dispatch PASSES
  Live N=3 all PASS: run1 436s (PRESENT/reuse), run2 820s (attempt-2 ABSENT branch — 2 fresh validation spawns + narrated "this host has no addressable reviewer reuse route", graded PASS through the assertion), run3 612s (PRESENT/reuse). Both branches exercised live; the previously-false-positive ABSENT path is now green. All 6 real failing transcripts now grade PASS; the 6 reuse runs still PASS (no regression).
- DONE: Value + no-regression — codex-FO ≤ 6004 B; `go test ./internal/ensigncycle ./internal/contractlint` + `go build ./...` green
  No codex-FO edit was needed (the FO is already correct); codex-first-officer-runtime.md unchanged at 5991 B ≤ 6004 B. ensigncycle + contractlint suites green; `go build ./...` exit 0. Diff is 2 test-tree files only.

### Summary

The 6-of-15 codex rejection-flow #141 failure is candidate 2: an assertion absent-detection false-positive, NOT a codex FO contract gap. The live Codex tool surface non-deterministically exposes the `send_input` reuse route; when it is absent the FO correctly declares absence and fresh-dispatches a separate cycle-2 reviewer (contract-correct), but `codexAddressableWorkerAbsent` only recognized 5 canned phrases that never matched the FO's real wording, so the assertion red-flagged correct behavior. The fix broadens the detector to read the FO's narration structurally (negation + reuse-route concept), keeps every non-tautology guard, adds 10 verbatim live-FO absence wordings as offline cases, and edits no FO prose. Live N=3 all PASS covering both the reuse and fresh-dispatch branches; the ABSENT path that was the false-positive is now green live.
