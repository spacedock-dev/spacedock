---
id: e3g7s1jtr05fp6n1w0z89w4v
title: Per-stage model declaration must support divergent per-host values (and Codex's effort axis), not just ignore-elsewhere
status: backlog
source: Captain conversation 2026-07-03, following dispatch-model-space-fable-sonnet5-1m (archived, id wcex4yjx4mvecybxjb43gwtw, PR #463) — that task's "host overlay" design deliberately rejected per-host keys ("no workflow needs it yet") in favor of validate-on-claude / ignore-with-note-and-null-elsewhere. A concrete need surfaced immediately after: declaring `fable` for the `ideation` stage on Claude while wanting a different, specific model (e.g. a Codex model id) plus a reasoning-effort tier (e.g. "xhigh") when the same workflow runs under Codex. The current mechanism cannot express this — Codex/Pi dispatch unconditionally gets `model: null`, with zero per-host substitution capability.
started:
completed:
verdict:
score:
worktree:
issue:
---

The shipped per-stage `model:` mechanism (`internal/dispatch/build.go`'s `runBuildFields`) only routes a declared model on Claude; on Codex and Pi the declared value is discarded entirely (an "ignored-with-note" stderr line, `model` always emitted as `null`). There is no way today for a workflow README to say "this stage runs fable on Claude, but a specific Codex model at a specific reasoning effort on Codex" — the mechanism was scoped to graceful degradation, not actual cross-host model routing, despite that being the stated motivation for the original task.

## Problem

{Ideation to flesh out. Seed framing: a workflow author wants one stage (e.g. `ideation`) to run under a specific, DIFFERENT model per host — not just "a model on Claude, defaulted elsewhere" — and Codex additionally exposes a reasoning-effort dimension (e.g. "xhigh") that has no representation anywhere in the current schema. The original task's own text anticipated this need and explicitly deferred it.}

## Proposed approach

{Explicitly NOT decided here — ideation's job. Two candidate shapes surfaced in conversation, neither chosen:
1. Per-host suffix keys alongside the existing flat `model:` (e.g. `model-codex:`), minimal parser surface but no home for Codex's `effort` axis without a second field.
2. A structured per-host block (e.g. `model: {claude: fable, codex: {model: ..., effort: xhigh}}`) — one field, extensible to whatever dimensions a host adds later, bigger schema/validation lift.
Ideation should evaluate these, invent a better one if warranted, and settle whether Codex's `effort` value space needs its own validation/enum the way Claude's `model:` does today.}

## Out of scope

{Ideation to determine — likely candidates: validating that a given Codex model id/effort combination is actually acceptable to the Codex CLI (that's Codex's own concern, not this task's); Pi's model space (no positive per-stage model space exists there yet per fo-dispatch-core.md's `pi-dispatch-model-stamping` note).}

## Acceptance criteria

{Ideation to flesh out with concrete Verified-by evidence. Seed AC: a workflow README can declare, for one stage, a Claude-space model AND a distinct Codex model+effort pair, and `dispatch build` emits the correct host-specific value for each — measured against today's baseline where the Codex/Pi branch always emits `model: null` regardless of what's declared.}

## Test plan

{Ideation to fill in.}
