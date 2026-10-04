---
title: The lane model ids are written down five times, and the Pi lane runs an older Luna
status: backlog
score: 0.8
source: "Captain, 2026-10-04: pin gpt-6-luna for CI, and move the Claude lane to sonnet-5.5. Local smoke on openai-codex/gpt-6-luna passed (root and child, 166s, grade artifact written)."
id: q67rfhn33tncbvf9xqbq960r
---

## Problem

Two problems with one cause.

Each lane model id is written down in several places. The Pi lane's
`openai-codex/gpt-5.6-luna:max` appears in `internal/ensigncycle/pi_liveenv.go`, in
three branches of the codex shim inside `.github/workflows/runtime-live-e2e.yml`, in
its step summary, in `docs/runtime-live-ci.md`, and in fixtures. The Claude lane's
`claude-sonnet-5` appears in the workflow matrix, in the harness's accepted-value
list, in the context-budget table, and in an artifact name.

And the Pi lane runs an older Luna. A local front-door smoke on
`openai-codex/gpt-6-luna:max` passed today: 5 reads, ensign skill at read 2, root and
child both `openai-codex/gpt-6-luna`, 166s, grade artifact written, no model error.

## Design

**One place for the lane models.** Put the ids in `internal/release` or
`internal/ensigncycle` as named constants, and print them for the workflow the way
`stamp-paths` prints the stamp targets: a small command that emits the model id per
lane and per auth path. The workflow matrix, the codex shim, the step summary, the
docs and the fixtures then consume that output instead of repeating literals. No
model id may remain in the workflow file.

**Pin the Pi lane to gpt-6-luna**, both auth paths: the OAuth spelling used with the
subscription, `openai-codex/gpt-6-luna:max`, and the API-key spelling beside it. Verify
the API-key spelling against the provider before pinning it; do not assume it from the
OAuth one.

**Claude moves only when the id is confirmed.** `claude-sonnet-5.5` is refused by the
installed CLI today: "There's an issue with the selected model (claude-sonnet-5.5). It
may not exist or you may not have access to it." Keep `claude-sonnet-5` until an
exact id is confirmed, and add the new value in the same single source when it is.

## Out of scope

- Recorded fixtures that name a past model on purpose. `auto_continue_replay_test.go`
  references real `claude-sonnet-5` streams and the CI artifact they came from; those
  are evidence, not configuration, and renaming them falsifies the record.
- Auth, quota, and thinking level.

## Expected surface and tolerance

`internal/ensigncycle/pi_liveenv.go`, the new print command, the workflow, the docs,
and the fixtures that assert the shim's argv. Estimate net +70, across 6 files.
Tolerance net +30 to +140, at most 8 files.

## Acceptance criteria

**AC-1 (VALUE) - One place holds each lane model.**
Verified by: the printed model list is what the workflow and the harness use, and a
test fails when they diverge. Falsifier: change the constant and require the workflow
site to follow.
**AC-2 - The Pi lane runs gpt-6-luna on both auth paths.**
Verified by: the harness constants and the shim's argv name it; a local front-door
smoke passes with root and child on it. Falsifier: leave the old literal and the shim
`--model` argument must differ from the constant.
**AC-3 - No model literal remains in the workflow.**
Verified by: searching the workflow for each id and requiring no match outside the
command that prints them. Falsifier: add one and the check must fail.
**AC-4 (no-regression) -** Real Pi discovery is unaffected, the Claude lane keeps its
current model until the id is confirmed, and the offline suite passes.

## Test plan

Behavioural tests only, with the model supplied independently. Every test must state
the reason it could fail. No test may assert a file's own text against a copy of it,
and no expected value may come from the thing it tests. The smoke is a local run with
the captain's credentials, never a CI lane run.
