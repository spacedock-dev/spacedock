---
title: Trim state.commit contract token load
status: ideation
source: "Captain follow-up after per-runtime contract token count (2026-06-20): shared-core grew +943 tokens since v0.22.0; biggest new body is «state.commit» at ~587 tokens. Determine and ship the leanest safe form."
started: 2026-06-20T18:36:31Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0221-layered-fo
sprint-readiness:
id: 6cccykszyvxz5mxhrf270fa2
---

The v0.22.0-to-current contract token comparison showed `skills/first-officer/references/first-officer-shared-core.md` grew by about 943 `o200k_base` tokens. The largest new section is `«state.commit»(slug)` at roughly 587 tokens. This task should determine whether that body can collapse now that the binary/state-verb path exists, and then implement a safe reduction or produce a documented no-cut decision if the mechanics remain load-bearing.

## Problem

`«state.commit»` currently carries detailed split-root commit mechanics in the boot-resident shared core. That may be the right temporary prose-function form, but it is now the largest shared-core token-growth driver. The contract should keep the state safety guarantees while avoiding duplicated or obsolete mechanics that the binary already owns.

## Proposed approach

Ideation should compare the current `«state.commit»` body against the shipped `spacedock state commit` behavior and related state-management text. Identify which parts are still unique FO obligations, which are binary-owned, and which are duplicated elsewhere. Prefer a minimal shared-core body that names the intent, required invariant, shipped command, and failure handling, with detailed mechanics moved to binary docs/tests or removed when already enforced.

Coordinate with `trim-dispatch-adapter-prose` (`ad`) only where the same contract-hygiene principles apply. This task is specifically about state/gate shared-core token load, starting with `«state.commit»`; it is not the dispatch-adapter trim.

## Out of scope

Do not redesign state branch topology, merge hooks, PR completion, or dispatch capability bindings. Do not weaken path-scoped commit, non-fast-forward retry, rebase-conflict halt, or split-root safety guarantees.

## Acceptance criteria

**AC-1 - The `«state.commit»` shared-core body is either measurably reduced or explicitly justified as still load-bearing.**
Verified by: an implementation report with before/after `o200k_base` token counts for `first-officer-shared-core.md` and the `«state.commit»` section, plus a diff or no-cut rationale.

**AC-2 - State safety guarantees remain enforced by behavior or tests, not only prose.**
Verified by: relevant Go tests or command-level fixtures covering path-scoped state commit, push/rebase handling where practical, and rebase-conflict halt/refusal behavior; if an existing test is the proof, cite it and run it.

**AC-3 - The resulting contract text follows runtime-support principles.**
Verified by: contractlint or targeted review showing the shared core names lifecycle/state capabilities and avoids unnecessary duplicated mechanics, mutable step-number coupling, and runtime-specific tool bindings.

## Test plan

Run focused tests for the state commit command and contract lint touched by the change, then `go test ./internal/contractlint ./internal/status ./internal/cli` at minimum. If the implementation changes Go state-commit behavior, run `go test ./...`.
