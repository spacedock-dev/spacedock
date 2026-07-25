---
title: Scrub foreign runtime markers from live Claude journeys
status: validation
source: "Repeated 6y live evidence on 2026-07-25: nested Claude inherited CODEX_THREAD_ID from the Codex captain host and first failed dispatch host detection before recovering with explicit --host claude"
started: 2026-07-25T16:30:35Z
completed:
verdict:
score: 0.7
worktree: .worktrees/spacedock-ensign-live-claude-runner-scrub-foreign-runtime-markers
issue:
sprint: durable-decisions
id: v3vt8gp2yffmn62r8p95gkph
gates:
    version: 1
    current:
        gate: gate:docs-dev:v3:ideation
    records:
        - id: gate:docs-dev:v3:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:v3-backlog-1
              briefing:
                id: briefing:docs-dev:v3:backlog:attempt-1:revision-1
                digest: sha256:2cda6cc97a6f5db6ef432781f25aae29f9781316c73999451d134ac549d26dc6
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:v3:backlog:1
                briefing: briefing:docs-dev:v3:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T08:11:16.055672Z"
                decision: approve
                reason: Two retained 6y Claude journeys prove the live harness leaks the Codex runtime marker; the narrow ideation preserves production ambiguity refusal and is required for a truthful host journey.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:docs-dev:v3:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:v3-ideation-1
              briefing:
                id: briefing:docs-dev:v3:ideation:attempt-1:revision-1
                digest: sha256:39164702899552b7c9c24b4e4067d0776eead9834ad958da907f50da8f2654d2
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:v3:ideation:1
                briefing: briefing:docs-dev:v3:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T08:24:47.475905Z"
                decision: approve
                reason: The design isolates the proven harness leak at existing per-host environment seams, preserves production ambiguity refusal, and supplies falsifiable builder plus live evidence within a declared six-file surface; implementation must wait for landed 6y.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
review-round:
    id: round:v3vt8gp2yffmn62r8p95gkph:implementation:3
    stage: implementation
    cycle: 3
    briefing:
        id: briefing:v3vt8gp2yffmn62r8p95gkph:implementation:round-3
        digest: sha256:cea9432b6c15da35c94e8737d24d7fa4d7da6e63d3269583f687e99449531e49
        digest-domain: canonical-bytes
        room-ref: ./review/implementation/round-3
---

## Problem

The cross-host live harness currently simulates target-host authentication and state
but can retain the captain host's runtime identity. `isolatedClaudeEnv` removes the
parent `CLAUDECODE` before the Claude front door, yet retains `CODEX_THREAD_ID`,
`PI_CODING_AGENT`, and `PI_CODING_AGENT_DIR`. Claude then supplies its own
`CLAUDECODE`, so a Claude session launched by a Codex captain contains two runtime
families and ordinary `spacedock dispatch build` correctly refuses the ambiguity.

This is a harness defect, not a production resolver defect. The two retained 6y
Claude journeys both failed an initial flag-free dispatch build with
`multiple runtime markers are set (CODEX_THREAD_ID, CLAUDECODE)`, then recovered
with `--host claude`. The cycle-23 stream records the refusal and recovery at
`review/implementation/cycle-23-counterexample/evidence/claude/recorded-lifecycle-stream.jsonl`
lines 176-182; the backlog gate records the independent repeat.

## Spike

A temporary `internal/ensigncycle` test set all four detector keys in the parent
and exercised the existing environment builders under `-tags live`. The observed
foreign-key counts were Claude 3 (`CODEX_THREAD_ID`, `PI_CODING_AGENT`,
`PI_CODING_AGENT_DIR`), Codex 2 (both Pi keys), and Pi 2 (`CODEX_THREAD_ID`,
`CLAUDECODE`). The spike file was removed after the run; this seeds the permanent
table tests below.

`go test ./internal/dispatch -run
'TestBuildHostResolutionFromFlagJSONAndEnv/(ambiguous-runtime|ambiguous-pi-runtime)$'
-count=1 -v` passes now. It proves production still refuses genuinely mixed runtime
families and gives the implementation an unchanged control.

No further mechanism spike is needed. Existing retained Claude streams prove that
the Claude front door supplies a fresh `CLAUDECODE` after the outer builder removes
the parent's marker, and the three builders already expose a deterministic
environment slice before process launch.

## Proposed approach

Treat the exact production detector keys in `internal/dispatch/build.go` as the
runtime-family boundary: `CODEX_THREAD_ID`, `CLAUDECODE`, `PI_CODING_AGENT`, and
`PI_CODING_AGENT_DIR`. Do not invent a broader prefix scrub.

At each existing per-host environment builder, remove only foreign detector keys
at the last environment-construction seam:

- Claude: extend both auth branches of `isolatedClaudeEnv` to remove Codex and both
  Pi detector keys, retaining the existing pre-front-door `CLAUDECODE` removal.
- Codex: extend `codexLiveEnv` to remove Claude and both Pi detector keys while
  retaining the Codex marker supplied by a Codex parent and preserving
  `CODEX_HOME`.
- Pi: build from the existing filtered-environment seam, remove Claude/Codex
  detector keys and inherited values for Pi state that the builder replaces, then
  append the target `PI_CODING_AGENT_DIR` and session/package state exactly once.

This is smaller and more falsifiable than changing `resolveBuildHost` to prefer one
marker: preference would weaken the production safety check and could silently
misroute real nested tooling. Always passing `--host` in live prompts was also
considered and rejected because it masks the harness defect rather than measuring
a truthful top-level-host journey. A process-global scrub was rejected because
credentials, HOME, PATH, and target state are already owned at the per-host builder
boundary.

After 6y lands, change its recorded-gate test prompt from:

> Run `dispatch build` once successfully, using `--host` for the current root
> runtime on that sole attempt when inherited runtime markers are ambiguous, then
> dispatch that exact artifact without rebuilding it.

to:

> Run `dispatch build` exactly once without `--host`, then dispatch that exact
> artifact without rebuilding it.

The Claude recorded-gate runner will additionally reject either the exact mixed
marker error in the stream or a `dispatch build` command containing
`--host claude`. This turns the former recovery into a falsifiable live proof.
Cycle-32 of 6y does not edit this prompt or live/env surface; it only changes
chronological selectors elsewhere in `recorded_gate_lifecycle_test.go`, so the
implementation should rebase onto landed 6y and preserve those changes.

## Expected surface

Baseline is `main` after PR #565 / 6y lands.

| File | Planned change | Expected changed LOC |
| --- | --- | ---: |
| `internal/ensigncycle/liveenv_test.go` | Claude foreign-marker scrub in both auth branches | 8 |
| `internal/ensigncycle/liveenv_decision_test.go` | `TestIsolatedClaudeEnvDropsForeignRuntimeMarkers` | 18 |
| `internal/ensigncycle/codex_liveenv_test.go` | Codex scrub plus `TestCodexLiveEnvDropsForeignRuntimeMarkers` | 22 |
| `internal/ensigncycle/pi_live_runner_test.go` | Pi scrub plus `TestPiLiveEnvDropsForeignRuntimeMarkers` | 28 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | remove explicit-host recovery from `recordedGatePrompt` | 2 |
| `internal/ensigncycle/claude_live_runner_test.go` | recorded-gate negative oracle for ambiguity and `--host claude` | 10 |
| **Total** | six test-harness files only | **88 changed LOC** |

Tolerance is ±25 changed LOC total, with no file outside
`internal/ensigncycle/`. Moving the Pi builder into a default-tag helper file is
allowed only if required to keep its deterministic marker test in `go test ./...`;
that may add one file but must remain within the same total tolerance. Any change
to `internal/dispatch/`, skills, production CLI output, or runtime docs is
out-of-scope and requires a gate reset.

## Acceptance criteria

**AC-1** With every detector key seeded in the parent, each per-host child
environment contains zero foreign runtime-family detector keys (current baseline:
Claude 3, Codex 2, Pi 2); a Claude child launched from Codex reaches the session
with Claude identity and flag-free dispatch host derivation exits 0 as Claude.
*Verified by:* the three named environment-builder tests plus the Claude
recorded-gate live scenario. Any retained foreign key, wrong derived host, or
non-zero flag-free build fails.

**AC-2** Runtime-marker isolation preserves each builder's required credentials,
PATH, isolated HOME, and target-host state (`CLAUDE_CONFIG_DIR`, `CODEX_HOME`, or
Pi agent/session/package roots), with each replaced key appearing exactly once.
*Verified by:* the new builder tests assert marker absence and the relevant
preserved values in the same returned environment slice; existing auth/HOME/PATH
tests remain green.

**AC-3** The recorded-gate Claude journey performs exactly one successful
`dispatch build`, without any explicit `--host`, mixed-marker ambiguity, or
harness-only retry, and still produces the successor's durable marker and commit.
*Verified by:* `TestLiveClaudeSharedScenarios/recorded-gate-lifecycle`, whose
command-log and stream assertions fail on a second build, non-zero build,
`--host claude`, the exact ambiguity error, or missing durable successor effect.

**AC-4** Production mixed-marker detection continues to reject Codex+Claude and
Codex+Pi inputs with the existing error and explicit-host guidance.
*Verified by:* the unchanged
`TestBuildHostResolutionFromFlagJSONAndEnv/ambiguous-runtime` and
`/ambiguous-pi-runtime` tests.

## Test plan

1. Run the three new environment tests under default tags; run the Pi test with
   `-tags live` if its helper remains live-tagged. These are cheap, offline tests
   against explicit parent environment values and exact child slices.
2. Run existing focused auth, HOME, PATH, and production host-resolution tests.
   They fail on collateral loss or a weakened ambiguity guard.
3. Run `go test ./...`, `go test ./... -race`, and
   `go test -tags live ./internal/ensigncycle -run
   'Test(IsolatedClaudeEnv|CodexLiveEnv|PiLiveEnv)' -count=1`.
4. Run one credentialed live proof:
   `go test -tags live ./internal/ensigncycle -run
   'TestLiveClaudeSharedScenarios/recorded-gate-lifecycle$' -count=1 -v`.
   Retain the command log and stream; success requires one zero-exit flag-free
   build and the committed successor marker.

Estimated implementation/test complexity is small (one working session plus one
Claude live journey). No fixture update or product documentation change is
required: all behavior changes are confined to test harness process environments,
and the production error text remains deliberately unchanged.

## Stage Report: ideation

- DONE: Prove the exact foreign-marker leak and select the smallest per-host scrub boundary using existing environment seams.
  A live-tagged temporary builder test measured foreign-key baselines of Claude 3, Codex 2, and Pi 2; the design scrubs only detector keys at each existing builder.
- DONE: Name the exact files, tests, and LOC while preserving production mixed-marker refusal unchanged.
  The expected surface names six harness files and 88 changed LOC ±25; focused unchanged production refusal tests pass.
- DONE: Define a falsifiable live proof that removes ambiguity and explicit-host recovery from the 6y Claude journey.
  The revised recorded-gate prompt forbids `--host`, and its live oracle fails on ambiguity, retry, explicit host, or missing committed successor effect.

### Summary

The ideation confines the fix to per-host live environment construction and its
behavioral tests, with no production resolver change. It converts the repeated 6y
recovery into a one-attempt, flag-free Claude journey and records exact scope,
controls, and evidence.

### Feedback Cycles

- Cycle 1: REJECTED — Roborev job 2307; surface 6/112 vs estimate 88 (127%); AC unchanged
- Cycle 2: REJECTED — validation adversarial pass; surface 6/111 vs estimate 88 (126%); AC unchanged; adversarial edit: make the sole dispatch-build exit non-zero while retaining the successor effect

## Stage Report: implementation

- DONE: Each live host child environment removes every foreign runtime-family detector while preserving credentials, PATH, isolated HOME, and the target host's state exactly once.
  `TestIsolatedClaudeEnvDropsForeignRuntimeMarkers`, `TestCodexLiveEnvDropsForeignRuntimeMarkers`, and `TestPiLiveEnvDropsForeignRuntimeMarkers` seed all detector families and fail on a retained foreign key, duplicate replacement, lost credential, or wrong HOME/PATH/host state.
- DONE: The Claude recorded-gate journey performs one flag-free dispatch build with no ambiguity recovery and still commits the successor effect.
  `TestLiveClaudeSharedScenarios/recorded-gate-lifecycle` passed in 256.26s; retained `command.log` records one attempt/one success after consume, and the oracle fails on another attempt, either `--host` form, mixed-marker ambiguity, or missing committed successor ancestry.
- DONE: Production mixed-marker refusal remains unchanged, and focused, full, race, live-tag, credentialed Claude, and Roborev evidence are retained.
  Commit `16aa2ec3`; `go test ./...`, `go test ./... -race`, the focused live-tag builder suite, and both unchanged ambiguity controls passed; Roborev job 2318 returned no issues after cycle-1 triage was recorded in advisory rounds 1 and 2.

### Summary

The six-file, 111-LOC harness-only change scrubs foreign detector keys at each existing host environment seam and leaves production host resolution untouched. The recorded-gate proof now observes every build attempt, forbids both explicit-host flag forms, and retains its stream and command log under `/tmp/spacedock-v3vt8gp2yf-claude-live-reviewed/claude-shared-scenarios/recorded-gate-lifecycle`.

## Stage Report: validation

- DONE: The exact candidate removes every foreign runtime-family detector from each live host child while preserving target-host credentials, PATH, HOME, and state exactly once.
  The three builder tests seed all four detector keys and fail on a foreign key, duplicate replacement, or changed credential/PATH/HOME/state value; default- and live-tag runs passed.
- DONE: The Claude recorded-gate journey performs one successful flag-free dispatch build with no explicit host, ambiguity recovery, or retry and still commits the successor effect.
  Retained `command.log` has one attempt, one matching `exit=0`, no explicit host, and dispatch after consume; the successful stream records committed marker `RECORDED-GATE-SUCCESSOR-DISPATCHED`.
- DONE: Production mixed-marker refusal remains unchanged and its existing Codex+Claude and Codex+Pi controls pass.
  Both unchanged `TestBuildHostResolutionFromFlagJSONAndEnv` ambiguity controls passed and would fail if refusal, error identity, or explicit-host guidance weakened.
- FAILED: The six-file harness-only diff, Roborev rounds, focused/full/race/live-tag tests, and retained credentialed Claude evidence support every acceptance criterion without a new product or instruction obligation.
  AC-3's live oracle does not verify dispatch-build exit status: one failed attempt followed by supported break-glass dispatch can retain one attempt and one durable effect while grading PASS.
- DONE: AC-1 — foreign runtime-family detectors are absent from each target child and flag-free Claude derivation succeeds.
  Exact-slice builder assertions passed; the retained Claude command is flag-free and exits 0, so changing a scrubbed key or derived host would fail the reproduced evidence.
- DONE: AC-2 — credentials, PATH, isolated HOME, and target-host state are preserved exactly once.
  Both Claude auth branches plus Codex and Pi exact-cardinality tests passed alongside the existing auth/config/PATH controls.
- FAILED: AC-3 — exactly one successful flag-free build with no ambiguity or retry and a committed successor effect.
  Outcome evidence is healthy, but `recordedGateLiveObservation` counts `begin` lines without requiring an `exit=0` line; this is a material evidence defect on the promised one-successful-build boundary.
- DONE: AC-4 — production mixed-marker detection still refuses Codex+Claude and Codex+Pi.
  The candidate changes no production file, and both unchanged command-level refusal controls passed.
- SKIPPED: Re-run the credentialed Claude journey.
  The retained stream/log is candidate-specific: the candidate-added artifact copy exists, records one successful attempt, has no exact ambiguity or host flag, and ends in a successful Claude result.
- SKIPPED: Promote the Pi live-tag CI-selection concern to material.
  Deferred risk: only a future unrun Pi builder regression triggers it; the required current live-tag proof passes, and it becomes material if CI coverage is promised or a supported Pi journey exhibits contamination.
- FAILED: Recommendation — PASSED/REJECTED.
  REJECTED; return to implementation for a narrow AC-3 evidence fix that pairs the sole attempt with exactly one successful exit and adds a non-zero-exit mutant.

### Summary

Candidate `16aa2ec3` is scope-correct and its reproduced focused, full, race,
live-tag, production-control, and retained live evidence all pass. Validation
rejects one material evidence defect: the recorded-gate oracle can qualify a
failed build followed by break-glass successor dispatch, contrary to AC-3.

## Stage Report: implementation (cycle 2)

- DONE: Pair the sole recorded-gate dispatch-build attempt with exactly one corresponding exit=0 before the successor effect can qualify; a failed build must never pass AC-3.
  `recordedGateLiveObservation` now counts a success only when an `exit=0` command exactly matches the recorded `begin`; `assertRecordedGateLifecycle` requires attempts/successes `1/1`, so changing the exit status or command bytes fails qualification.
- DONE: Add an adversarial mutant with one nonzero dispatch-build exit and a retained successor effect, and prove the observer rejects it.
  `TestRecordedGateLifecycleRealCLIReplay/failed-build` changes only the valid log's dispatch-build exit to `exit=1` after the successor effect is committed; it failed red before the fix and passes only because the corrected observer rejects that mutant.
- DONE: Keep the correction harness-only, preserve all existing host-scrubbing and ambiguity controls, rerun focused/full/race/live-tag verification, and request final Roborev review.
  Correction commit `6ea4cdc8` changes only `recorded_gate_lifecycle_test.go`; focused live-tag builders, unchanged production ambiguity controls, `go test ./...`, and `go test ./... -race` passed, and Roborev job 2342 returned no issues in advisory round 3.

### Summary

The correction is a one-file, 26-LOC AC-3 oracle repair atop the existing six-file harness candidate, with no production or workflow behavior change. The aggregate candidate is six files/133 changed LOC; its successful retained Claude journey remains valid, while the new deterministic mutant closes the failed-build evidence hole.
