---
id: 5n4k6djrq8gtvd54zg9s6zhs
title: Claude-opus live agent's recorded-gate-lifecycle successor dispatch not observed as committed-before-dispatch
status: backlog
source: "Found on PR #600 (collapse-gate-approval-ceremony) claude-live (claude-opus-4-8, CI-E2E-OPUS) CI, 2026-08-02, run 30754109029, job 91513297850: TestLiveClaudeSharedScenarios/recorded-gate-lifecycle failed 'successor dispatch was not observed after consume' (recorded_gate_lifecycle_test.go:135, assertRecordedGateLifecycle's recordedGateCommittedBeforeDispatch check). Pulled the actual Claude transcript (runtime-live-e2e-claude-live-claude-opus-4-8 artifact): the agent issued the classic 3-step ceremony (gate record --decision approve, then gate consume, then dispatch build) in the correct order, and its final message narrates a coherent successful completion (marker recorded, durable commit, stage report present) -- no obvious agent-side ordering mistake like the sibling codex-live failure on the same run (dispatch build --checklist-file raced ahead of the file write, filed separately as codex-live-dispatch-build-checklist-race). Checked whether this is another PR #599 (gate-schema-simplification) casualty: built the binary from this branch and ran a real prepare/record/consume cycle in a scratch fixture to measure the actual YAML nesting depth of 'state: consumed' under the new schema -- still 16 spaces, unchanged, so the exact-string git log -S pickaxe search in recordedGateLiveObservation (recorded_gate_lifecycle_test.go ~L1019) should still match. Root cause NOT yet identified. Captain directed: treat as a candidate flake alongside the codex-live issue, do not block the merge, rerun to confirm green, file for diagnosis."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
---

`recordedGateCommittedBeforeDispatch` (internal/ensigncycle/recorded_gate_lifecycle_test.go:1193) verifies, via `git log -S<exact-string>` pickaxe searches and `git merge-base --is-ancestor` ancestry checks against the state checkout's git history, that the record-commit precedes the consume-commit precedes the dispatch-head commit. On this run it returned false for the Claude-opus live scenario even though: (a) `dispatch.builds == 1` and `successfulBuilds == 1` (the build-count check earlier in `assertRecordedGateLifecycle` passed), (b) the agent's own transcript shows the 3 commands issued in the correct order, and (c) the final message describes a complete, coherent success.

## What's ruled out so far

- Not the codex-live issue (dispatch build racing ahead of a checklist-file write) -- different assertion, different mechanism, and the transcript shows correct command order.
- Not a YAML-nesting-depth regression from PR #599's schema simplification -- measured directly (16 spaces, unchanged) via a real prepare/record/consume cycle built from this branch's binary.

## Open questions for diagnosis

- Is `close`/`consumed`/`dispatchHead` commit resolution (the `git log --reverse --format=%H -S...` pickaxe searches at recorded_gate_lifecycle_test.go ~L1017-1019) actually finding the right commits for this run, or resolving empty/wrong due to some other git pickaxe quirk (e.g. `-S` only detects a *change in occurrence count*, not any touching diff -- worth checking whether the fixture's repo history could trip that)?
- Is this reproducible on a clean rerun, or a one-off? (Rerun in isolation before assuming either way.)
- If reproducible: is the defect in the test's git-ancestry detection mechanism, or a genuine ordering/timing issue in how the live Claude harness records the consume commit relative to dispatch?

## Out of scope (for now)

Any code change to the assertion or to gate/dispatch mechanics -- this entity is for diagnosis first.
