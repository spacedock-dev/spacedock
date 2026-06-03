---
id: p48fjxe9e2rz23afaqmzj1pg
title: Encode the live-verification requirement — gate runtime-observable ACs on a live run + promote 8y's shared scenarios to a first-class live-scenario primitive
status: backlog
source: captain (2026-06-03, session 11) — after yy's fix passed offline + 3 detached-audit cycles but FAILED its live run, the captain noted the live-verification rule is prose-only (yy slipped) AND that agents don't intuitively reach for "descriptive runbook → run a real agent → grade durable outcomes" live tests
score: "0.40"
worktree:
started:
completed:
verdict:
issue:
---

The requirement that a **runtime-observable AC** (one whose truth can only be decided by running the real producer) be confirmed by a **live run** is encoded only as advisory prose — the README validation stage and the merge-on-proof bar in handoffs. Nothing refuses to mark such an AC PASSED until a live run is on the record. `yy` (sonnet-live-ci-flake, #282) was marked `pending-live-run` and merged with its live AC unconfirmed; the post-merge live retrigger then revealed the fix doesn't actually close the flake. This violates the workflow's own Working Principle ("prefer a code gate over a prose-only rule; a prose-only rule has a ceiling of wording-present, which is not behavior").

## Problem

Two coupled gaps:

1. **The live-verification rule is prose-only and unenforced.** A runtime-observable AC can reach terminal/PASSED on offline proof alone (unit tests, a recorded-stream fixture, a contract-text check). But a recording proves the *consumer* (the streamwatcher localizes a hang in a frozen stream), never the *producer* (does the real model, reading the contract, actually exit?) — the fix produces a *different, future* stream you can only get by running it live. yy passed every offline check + 3 adversarial audit cycles and was still wrong.

2. **The live-scenario pattern is not reachable/intuitive for agents.** Agents default to deterministic offline proxies because "write a runbook → run a real agent → grade the durable outcomes" is **eval design**, not software testing — it requires defining observable success criteria for non-deterministic behaviour, a skill the contract/templates neither surface nor make cheap. `8y` built the substrate (host-neutral `sharedRuntimeScenarios()` + per-host runner adapters + durable state/output assertions) but it is CI-internal, not a primitive the FO/ensign reaches for. So agents pick the proxy and don't think to run the thing.

## Proposed approach (three layers, code-first)

1. **AC classification + a real gate.** At ideation, classify each AC as *offline-checkable* (a unit test / fixture / text check decides it) or *runtime-observable* (only a live run decides). Extend the existing AC-evidence cross-check (`status --validate` / a `gate verify`) so the terminal gate **refuses** a runtime-observable AC marked PASSED unless it cites a live-run artifact (a CI run id or the archived session JSONL). This turns `pending-live-run` from a hopeful label into an enforced state.
2. **Promote 8y's scenarios to a first-class live-scenario primitive.** Generalise the host-neutral table + runner adapters + durable-outcome assertions into a reusable primitive the FO/ensign authors as easily as a Go test: a scenario = {descriptive runbook/prompt, setup/fixture, expected durable outcomes (state + output, NOT transcript phrasing)}, runnable by a real Claude/Codex agent and graded on those outcomes.
3. **Surface the pattern in the dev template.** Add to "choose proof at the claim's altitude": *if the claim is what an agent/model DOES at runtime, the proof is a scripted live scenario graded on durable state/output — a recording proves the watcher not the producer; a contract-text check proves the words not the behaviour.*

**Riskiest unknown — spike first:** the gate's classification + citation mechanism. How does the binary KNOW an AC is runtime-observable (an explicit tag? a `Verified by: live …` convention?) and how does it verify a live-run-artifact citation is real (a run id that exists / a present JSONL)? Pin this before designing the prose, since a sloppy classifier either over-gates (blocks offline-checkable ACs) or under-gates (lets runtime ACs through — the status quo).

## Out of scope

- Re-running yy's fix — that's `at` (non-interactive-teardown-exit); this entity is the *meta*-fix that makes such slips gate-caught.
- A general eval framework beyond what spacedock workflows need.

## Acceptance criteria

**AC-1 — A runtime-observable AC cannot reach terminal/PASSED without a cited live-run artifact.**
Verified by: a test driving the gate on a fixture entity carrying a runtime-observable AC — RED (refused) when no live-run artifact is cited, GREEN when one is; an offline-checkable AC is unaffected.

**AC-2 — The live-scenario primitive is authorable + runnable independent of 8y's CI-internal plumbing.**
Verified by: a scenario authored as {runbook, setup, durable-outcome assertions}, run by a real agent and graded on outcomes, with at least one negative case (a broken durable outcome reds the grade).

**AC-3 — The dev template surfaces the live-scenario-for-runtime-claims guidance.**
Verified by: a presence check over the dev template confirming the proof-altitude guidance names the live-scenario pattern and the recording-proves-the-watcher / text-check-proves-the-words distinction.

**AC-4 — The yy slip would now be gate-caught.**
Verified by: a fixture reproducing yy's shape (a runtime-observable AC marked PASSED with only offline proof) is REFUSED by the gate from AC-1.

## Test plan

- Spike the AC-classification + citation-check mechanism FIRST (the riskiest unknown).
- Go tests for AC-1/AC-4 (gate behaviour over fixture entities) and AC-3 (template presence check).
- A live-scenario exercise for AC-2 (the primitive grading a real agent run + a negative case).
- High-stakes (the gate machinery + the contract/template) → detached adversarial audit before merge.

## Notes

- Motivated by yy's live-run failure (session 11): offline + 3 audit cycles passed, the live run failed — the strongest possible case for the gate. `8y` built the substrate; `at` will be the first entity that *needs* this primitive (its live AC is a scripted scenario).
- Pairs with the binary-simplification `gate verify` candidate (the AC-extractor that doesn't yet exist in `validate.go`) — this is the runtime-AC-aware extension of it.
