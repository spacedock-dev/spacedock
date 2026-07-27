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

## The cause: PR #565 merged live scenarios that had never been green

Established 2026-07-27 by tracing the failing tests' provenance, after the captain observed that four lanes failing simultaneously is far likelier to share one cause than to be four independent flakes. It is not a pre-existing flake and not a gradual drift.

All three failing test files belong to `first-officer-gate-command-lifecycle` (6y, PR #565):

- `internal/ensigncycle/recorded_gate_lifecycle_test.go` — added by `c9633279`, 2026-07-23, titled **"WIP counterexample: FO recorded gate lifecycle"**.
- `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` — added by the same commit.
- `internal/ensigncycle/live_gate_stop_test.go` — pre-existing, but last substantially rewritten by `9577380d`, 2026-07-25, "feat: complete first-officer gate lifecycle".

The timeline is exact:

| when | what | live lanes |
| --- | --- | --- |
| 2026-07-23 | `c9633279` adds the recorded-gate live scenarios on 6y's branch | — |
| 2026-07-24 13:29 | last green run `30097092217`, on #564's branch, which does not carry them | GREEN |
| 2026-07-24 17:15 | first red run, on 6y's branch, where they now execute | RED |
| 2026-07-25 22:22 | 6y merges as `deac7f8a` with `pi-live` at FAILURE and three lanes never run | RED |
| since | every branch rebasing onto main inherits them | RED |

#564 is exonerated: `git diff 642ca090..cc51e518` is empty, so the merge commit's tree is identical to the tree that ran green.

**These scenarios have never passed.** They entered as a WIP counterexample — a test written to demonstrate a failure — and no run of them has been green on any branch. So the correct question is not "what regressed" but "what were they written to expect, and was that expectation ever met."

That reframes the fix. Repairing FO conduct to satisfy them assumes the scenarios encode correct expectations; they may instead encode a contract behaviour that was never implemented, or a grader written against an intended design rather than a shipped one. Ideation must establish which before changing either side.

## The rotation, correctly interpreted

The failing sub-assertion **rotates between runs** on the same code. Observed variants across runs and baselines:

- `live_gate_stop_test.go:28: open bound entity count for "state: open" is not 1`
- `live_gate_stop_test.go:28: gated entity is not held at its open validation boundary`
- `live_gate_stop_test.go:28: gate hold retried after a failed prepare`
- `claude_live_runner_test.go:130: recorded gate lifecycle graded FAIL: gate review omits its decision facts`
- `codex_live_runner_test.go:41: rejection trajectory left 1 implementation reports, want at least 2`
- `codex_live_runner_test.go:41: the FO advanced the approved entity but did not dispatch its next stage`

A deterministic regression hits the same assertion every run. Rotation across the same test functions is what a set of never-passing new scenarios looks like: several independent expectations are unmet, and which one trips first varies with live-session timing. It is not evidence of a pre-existing flake. Note also that `claude_live_failure_diagnostic_impl_test.go` is diagnostic-only — its own ABOUTME says it "reports it only after another failure" and is "silent on success" — so its wrong-root and broad-search lines annotate whichever primary failure occurred and must not be chased as failures themselves.

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

Ideation fills these in. The end state is a green live signal whose greenness is trustworthy, with at least one criterion measuring against something that can move the wrong way — a count of consecutive green runs across branches, not the presence of a fix. Because these scenarios have never passed, ideation must first establish what each was written to expect and whether that expectation was ever met, then decide per assertion whether the grader or the First Officer conduct is wrong. Repairing the wrong side yields a green lane that grades nothing, and assuming the scenarios are correct because they are committed is the specific error to avoid.

## Test plan

Ideation fills this in. The local reproduction above is the substrate; CI is confirmation, not discovery.
