---
title: Bind the FO's own gate and triage decisions to the evidence bar in the runtime-neutral Working Principles
status: backlog
source: "ezf/hf merge incident (2026-06-16): the FO merged a Claude-adapter change (ezf, claude-fo-dispatch.md) on deterministic lanes while leaving claude-live unapproved, and labeled a live-CI red 'the known flake' without reading the actual failing test (it was TestLiveZeroDiscover, not TestLiveEnsignCycle). Root cause: the FO contract aims its proof discipline (code-gate-over-prose, observe-the-behavior, the AC cross-check) at the ensign's DELIVERABLE and the gate review, but does not bind the FO's OWN dispatcher decisions — what verification a change requires, whether a result is green, what a failure means. The dev-workflow README now encodes the dev-specific path->lane realization; the generic principle belongs in the runtime-neutral FO Working Principles so it is intrinsic to every FO on every host, not a per-workflow override."
sprint:
sprint-readiness: ready
issue:
id: z25gzd3s3p0a90c18t8tzs6r
---

## Problem
The FO operating contract (`skills/first-officer/references/first-officer-shared-core.md`, `## Working Principles`) enforces "prefer a code gate over a prose-only rule" and the gate AC cross-check on the ENSIGN's deliverable and the gate review. It does not bind the FO's OWN dispatcher decisions to the same evidence bar. Two real failures followed (ezf/hf, 2026-06-16): (1) the FO judged `claude-live` "unrelated" to a change in `claude-fo-dispatch.md` and merged without it; (2) the FO classified a live-CI red as "the known flake" by inheriting a handoff label instead of reading this run's failing test. Both are FO-side judgment calls the contract leaves unconstrained, so the discipline is self-exempting exactly where it matters.

## What's needed
A runtime-neutral Working Principle (host-agnostic, no `claude-live`/`codex-live` specifics) that the FO holds its own gate/merge/triage decisions to the bar it imposes on workers:
- Required verification follows from WHAT CHANGED, not the FO's sense of relevance; a relevant check is not waived as "unrelated" by intuition, and a flaky-but-relevant check is re-run to green, never skipped.
- A result is "green" only when the relevant check actually ran and passed — an unapproved/skipped/cancelled lane is not a pass.
- A failure is READ from this run's evidence (failing test, assertion, error); a prior session's or handoff's label is a hypothesis to confirm, not a verdict to apply.

The dev-workflow README's `Proof policy` then INHERITS this (its path->lane bullet becomes the dev-specific realization and can slim to reference the generic principle).

## Notes for ideation
This is itself a contract/scaffolding change, so it must go through implementation + validation and — per the principle it adds — the `claude-live`/`codex-live`/`pi-live` drives must be green before merge (the dogfood). The hard part is the AC: a behavioral FO principle resists a clean code gate (a contractlint "the phrase is present" check is the prose-grep tautology this very policy bans). Ideation must decide what genuinely proves it — a live FO scenario observing the right gate/triage decision, the structural absence of host-coupling, the binary-level merge check, or some composition — not a self-referential prose check.
