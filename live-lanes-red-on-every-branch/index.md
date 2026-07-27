---
id: se0v37bt7mhsrmhta1nyns0r
title: "Runtime Live E2E has failed on every branch for three days, so no merge since 2026-07-24 carries a live signal"
status: backlog
source: "Found attributing PR #571's four red live lanes, 2026-07-27. The change was exonerated by baseline comparison and local reproduction; the lanes were already red on every branch."
started:
completed:
verdict:
score: 0.85
worktree:
issue:
---

Restore a trustworthy live signal. The lanes that grade First Officer conduct end-to-end have been red continuously for three days across unrelated branches, so every merge in that window shipped without the verification those lanes exist to provide.

## The observation

Last green `Runtime Live E2E` run: **30097092217, 2026-07-24T13:29**. Every run since has failed — 13 of the last 14, the fourteenth having no conclusion — across at least five unrelated branches: `prepare-provider-neutral-gate-room`, `codex-wait-agent-steering-semantics`, `sync-merge-guard-archive-state`, `recorded-gate-lifecycle-*`, and `version-output-runtime-and-sandbox-state`. Offline, build and both install jobs are green throughout, so this is specific to the live lanes.

The consequence is not that one PR is blocked. It is that **no merge since 2026-07-24 has a live signal at all**. PR #565 merged with `pi-live` at FAILURE, and the merge commit `deac7f8a` carried a single check-run. A pre-cut review of that commit already recorded "the audited commit has no live-lane evidence" as a release prerequisite; this entity is the cause behind that symptom.

The sharpest cost: the `durable-decisions` sprint's deliverable *is* recorded gate conduct, and the lanes that grade gate conduct have been dark for the entire period it was being built.

## The signature says nondeterminism, not regression

The failing sub-assertion **rotates between runs** on the same code. Observed variants across runs and baselines:

- `live_gate_stop_test.go:28: open bound entity count for "state: open" is not 1`
- `live_gate_stop_test.go:28: gated entity is not held at its open validation boundary`
- `live_gate_stop_test.go:28: gate hold retried after a failed prepare`
- `claude_live_runner_test.go:130: recorded gate lifecycle graded FAIL: gate review omits its decision facts`
- `codex_live_runner_test.go:41: rejection trajectory left 1 implementation reports, want at least 2`
- `codex_live_runner_test.go:41: the FO advanced the approved entity but did not dispatch its next stage`

A deterministic regression hits the same assertion every run. Rotation across the same test functions is the signature of nondeterministic live-conduct grading. Note also that `claude_live_failure_diagnostic_impl_test.go` is diagnostic-only — its own ABOUTME says it "reports it only after another failure" and is "silent on success" — so its wrong-root and broad-search lines annotate whichever primary failure occurred and must not be chased as failures themselves.

## Recurring core, from a local reproduction

The failures cluster on First Officer gate-binding conduct rather than on any product output. In a local run the live FO bound `gate:recorded-gate-task:implementation` where the grader expected the `3k-validation-1` attempt; rejection trajectories stop at one implementation report where the grader wants two (original plus cycle-2 rework). Whether the graders drifted from current FO contract behaviour, or FO conduct regressed, or both, is the question ideation must answer — and the two are distinguishable, because the contract and the graders are separately versioned.

## It reproduces locally, which makes this cheap to work

`TestLiveDefaultHeadlessStopsAtGate` fails locally against real sessions on both `52e0c6a6` and its merge-base `50f8d1fb` with a byte-identical assertion (`gated entity is not held at its open validation boundary`), at roughly 200 seconds per run. No CI minutes are needed to iterate.

## The attribution technique worth keeping

The method that settled PR #571's attribution should be the standard response to a red live lane, because it distinguishes "my change broke it" from "it was already broken" without guessing:

1. Compare against the merge-base's own CI run, plus a second unrelated-branch run.
2. Run the failing test locally on both the branch commit and the merge-base, and compare the exact assertion text.
3. Grep the graders for every string the change touched, to test an output-shape hypothesis rather than assume one.

Applied to #571 that produced: same lanes, same test functions, byte-identical local assertion on both commits, and zero grader references to any changed string. Not attributable.

## Out of scope

- Fixing any individual PR. #571 was exonerated by this work and needs no change.
- Widening a feature branch to repair shared CI. The diagnosing worker correctly declined to.
- The nine-runtime-token and live-CI-diagnostics entities, which are adjacent but narrower.

## Acceptance criteria

Ideation fills these in. The end state is a green live signal whose greenness is trustworthy, with at least one criterion measuring against something that can move the wrong way — a count of consecutive green runs across branches, not the presence of a fix. Ideation must also decide whether the graders or FO conduct drifted, since repairing the wrong side would produce a green lane that grades nothing.

## Test plan

Ideation fills this in. The local reproduction above is the substrate; CI is confirmation, not discovery.
