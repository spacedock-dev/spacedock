---
id: 5n4k6djrq8gtvd54zg9s6zhs
title: recorded-gate-lifecycle grader has two unrelated bugs that both produce "successor dispatch not observed after consume"
status: backlog
source: "Originally filed 2026-08-02 as an opus-specific, root-cause-unknown candidate flake (TestLiveClaudeSharedScenarios/recorded-gate-lifecycle, PR #600 run 30754109029 job 91513297850). Amended 2026-08-03 after a dedicated diagnostic (opus, high effort) root-caused it with a local repro: the same error string is produced by TWO distinct, deterministic one-line grader bugs in internal/ensigncycle/recorded_gate_lifecycle_test.go, not live-model nondeterminism. Confirmed recurring on a second, independent instance: codex-live hit the identical error string on PR #600's rerun (run 30754109029, job 91518827444) via the SECOND bug below, proving this is not opus-specific. Classification: test-defect, high confidence (opus reproduced locally; codex confirmed from live stdout). Captain directed: file/amend based on the classification workflow's findings."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
---

`assertRecordedGateLifecycle` (internal/ensigncycle/recorded_gate_lifecycle_test.go:121) collapses `!o.dispatch.ordered || !o.dispatch.committed` into one error string ("successor dispatch was not observed after consume"), so two unrelated bugs in how `ordered`/`committed` get computed both surface identically. Both are one-line fixes.

## Bug 1 — unanchored substring match lets a harmless `--help` probe poison the ordering check

`recordedGateLiveObservation` (recorded_gate_lifecycle_test.go:983-986) locates command-log phases with unanchored `strings.Index` over the whole log:

```go
prepareAt := strings.Index(log, "exit=0\tgate prepare ")
bindCommitAt := strings.Index(log, "exit=0\tstate commit recorded-gate-task")
decisionAt := strings.Index(log, "exit=0\tgate record ")
ordered = prepareAt >= 0 && bindCommitAt > prepareAt && decisionAt > bindCommitAt
```

`"exit=0\tgate record "` also matches the read-only line `exit=0\tgate record --help`. In the failing opus run the agent ran `gate prepare --help` then `gate record --help` (both harmless, read-only) before doing any real work. Measured on the actual run artifact (artifact 8835777822): `prepareAt=251` (the `--help` probe), `decisionAt=303` (the `gate record --help` probe, matched instead of the real decision-record command later in the log), `bindCommitAt=1648` -> `decisionAt < bindCommitAt` -> `ordered=false`. A passing rerun's log (artifact 8836227123) has only the `gate prepare --help` probe, no `gate record --help` probe, so `decisionAt` correctly lands on the real command and `ordered=true`. The line-983 `dispatch build` phase-detection already excludes `--help` (confirmed at ~line 1010); this same guard is simply missing for the `prepare`/`record` phases.

Reproduced locally with the real binary (scratch test at the diagnostic's scratchpad, `internal/ensigncycle/zz_repro_test.go`): the same command sequence without help probes grades PASS; inserting a `gate record --help` probe before the real `state commit` flips `ordered` to false and reproduces the exact failure message.

**Fix:** anchor the phase-detection to exclude `--help` invocations, the same way the existing `dispatch build` detection already does (or match complete lines rather than substrings).

## Bug 2 — hardcoded resolution-ID literal misses the room-backed close path

`recordedGateLiveObservation` (recorded_gate_lifecycle_test.go:1018) pickaxes a literal chat-path resolution ID:

```go
closeCommit := ... git log --reverse --format=%H -S"id: resolution:spacedock:recorded-gate-task:validation:1" ...
```

That exact literal is produced only by `chatResolutionID` (internal/gates/operation.go:758-763) -- the `--decision`/`--actor` close path. Codex's live run closed via the **room-backed** path instead: `gate record ... --room ... --consume`, which produces a provider-authored resolution ID (`resolution=resolution:captain-recorded-gate-task-validation-1`) written verbatim from the room's provider `result.json` by `recordRoomLocked` (internal/gates/operation.go ~line 400). The hardcoded pickaxe string never matches that ID, so `close == ""`, and `recordedGateCommittedBeforeDispatch` (line 1193) returns false on its very first guard clause -- independent of Bug 1, and independent of anything about `--consume`/`--stamp` timing (both the collapsed and classic ceremony shapes grade PASS correctly once this is accounted for).

Ironically, `recordedGateLiveObservation` already derives `resolutionID` correctly by regex from the entity's own post-state a few lines later (line 1023) and `assertRecordedGateLifecycle` uses that derived value elsewhere -- only this one commit-lookup pickaxe still hardcodes the literal.

**Fix:** use the regex-derived `resolutionID` (already computed at line 1023) for the pickaxe search instead of the hardcoded chat-path literal, or pickaxe on a room-agnostic substring common to both closing shapes.

## What's ruled out

- Not model nondeterminism -- both bugs are deterministic string-matching gaps triggered by legal, unmodeled agent behavior (a `--help` probe; closing via the room-backed path), not by anything random.
- Not the codex-live dispatch-build-checklist-race issue (different assertion entirely; filed separately as `codex-live-dispatch-build-checklist-race`).
- Not a YAML-nesting-depth regression from PR #599's schema simplification -- measured directly (16 spaces, unchanged).
- Not `-S`'s occurrence-count semantics, not `--stamp` timing/durability, not a real difference in behavior between the collapsed `--consume` and classic 2-step ceremony shapes -- all confirmed fine in the local repro; the two bugs above are the entire explanation for both recurrences.

## Acceptance criteria (rough, for ideation to firm up)

- **AC-1** — A `--help` probe (of `gate prepare`, `gate record`, or `gate consume`) issued at any point in a live agent's command log no longer perturbs `ordered`'s phase detection. Verified by a fixture/unit test reproducing Bug 1's exact repro (insert a `gate record --help` line before the real sequence; assert `ordered` stays true).
- **AC-2** — A room-backed close (`gate record --room ... --consume`) is correctly recognized as committed by `recordedGateCommittedBeforeDispatch`. Verified by a fixture/unit test reproducing Bug 2 (a room-backed close followed by dispatch; assert `committed` is true).
- **AC-3** — The classic chat-path close (`--decision`/`--actor`) continues to work unchanged; no regression on either existing passing shape.

## Out of scope

Any change to the actual gate/dispatch mechanics, `--stamp`, or `--consume` -- confirmed correct. This is a test/grader-only fix.
