---
id: m8yy9wbsm1b832a5078fgyza
title: Scenario-testing principles — a scenario is a natural-language spec; codified test and LLM run are two executor implementations; cost/coverage key by scenario × {mode, runtime}
status: ideation
source: captain (2026-06-03, session 12) — "we should have scenario as natural language, like bdd. and codified test is a variant of an implementation. running through llm is another implementation." + "document this as scenario-testing principles first, the semantics not the syntax. 8y should have concluded with a few scenarios; we build upon it. 4n is tracking cost, which should be keyed by the scenario (with mode, runtime as variants)."
started: 2026-06-03T16:55:07Z
completed:
verdict:
score: 0.41
worktree:
issue:
---

Establish the **semantic foundation** for scenario-based testing in Spacedock: what a *scenario* is, what an *executor* is, and the *variant axes* a scenario runs under — so `8y`'s shared runtime scenarios, `p4`'s live-verification gate, the authoring primitive, and `4n`'s cost-ledger all build on one shared model instead of re-inventing it. **Semantics first; syntax (Gherkin vs runbook, file layout, API) is explicitly deferred.**

## Problem

`8y` shipped a host-neutral `sharedRuntimeScenarios()` table + per-host runner adapters + durable-outcome assertions — but it concluded as CI-internal Go plumbing, not as a *named, reusable scenario foundation*. As a result every downstream consumer re-derives the model:
- `p4` (live-verification-gate) needs "a runtime-observable claim is checked by a real agent run" but has no shared definition of scenario/executor to anchor on.
- `4n` (journey-cost-ledger, codex peer's lane) tracks cost per journey but isn't keyed by a scenario identity with variant dimensions, so costs aren't comparable across model/host.
- The recurring `yy` failure (offline proof passed, live failed) is the symptom of treating "offline test" and "live run" as *different kinds of thing* rather than two *executors of the same scenario*.

Without a documented semantic model, agents default to the cheap deterministic proxy and never reach for the faithful live executor (eval design is unintuitive), and the same scenario gets modeled twice with no shared identity.

## Proposed approach (semantics to document — ideation refines)

A principles document (likely `docs/specs/scenario-testing-principles.md`) defining:

1. **A scenario is a natural-language behavioral specification.** It names a behavior / user-or-agent journey and its expected **durable outcomes** (entity state before→after, archive state, on-disk artifacts, durable user-facing output obligations) — independent of *how* it is checked. Assertions are over durable outcomes, **never transcript phrasing** (non-deterministic).

2. **Executors are interchangeable implementations of a scenario's check**, at different fidelity:
   - **Codified executor** — a deterministic test (Go fixture/unit). Cheap, fast, always-on in CI; proves the *modeled* logic (the *consumer* side — e.g. a watcher over a recorded stream). Cannot prove the producer (a recording is frozen).
   - **LLM executor** — a real agent run (Claude/Codex). Expensive, non-deterministic; proves the real *producer* (the agent actually behaving), graded on the scenario's durable outcomes.
   - The model is open to further executors. The scenario is the spec; executors are pluggable.

3. **Variant axes: a scenario executes as `scenario × {mode, runtime}`.** `runtime` = host (claude / codex) — `8y`'s per-host runner adapters are this axis. `mode` = model / effort (e.g. sonnet / opus). The `(scenario, mode, runtime)` tuple is the unit that is run, graded, and **measured** — the key `4n` uses for cost, and the key a coverage matrix uses to demand parity.

4. **`8y`'s three scenarios are the seed instances**, named here as the first foundation: `gate-guardrail`, `rejection-flow`, `merge-hook-guardrail`. The principles formalize what `8y` should have concluded with.

5. **The model the consumers build on:**
   - **`p4` citation gate** — a runtime-observable AC (only the LLM executor can decide it) is satisfied only by a cited LLM-executor run; the codified executor alone does not satisfy it. This is what stops "default to the cheap proxy" (the `yy` hole) by construction.
   - **Authoring primitive** — makes a scenario authorable outside the test package (generalizes `8y`'s substrate), so writing a scenario is as cheap as a Go test.
   - **`4n` cost-ledger** — cost keyed by `(scenario, mode, runtime)`, so cost is attributable per scenario and comparable across variants.

6. **Why this dissolves the `yy` failure mode:** write the scenario once; run the codified executor every CI (proves logic) AND the LLM executor when producer-fidelity is required (proves the real agent), with the gate forcing the LLM run to be cited for producer-claims.

## Out of scope

- Concrete **syntax** — Gherkin vs freer runbook+assertions, file format, the primitive's Go API. (A follow-on entity; this is semantics only.)
- Re-implementing `8y` or building the primitive itself (separate entities that cite these principles).
- `4n`'s implementation (codex peer's lane) — these principles define the *key* (scenario × {mode, runtime}) it should adopt; coordinate, don't author.

## Acceptance criteria

{Ideation defines these as real gates. The challenge per the README: a principles doc must not be "a document about itself." Candidate teeth — the doc's semantics are *referenced/honored by a consumer that can fail*: e.g. a presence check that the doc names scenario/executor/variant-axes AND a check that `8y`'s seed scenarios + `4n`'s cost key are expressed in the doc's terms. Ideation must find the external, failable proof — likely "a consuming entity (p4 / 4n) is keyed to these principles," not the prose alone.}

## Test plan

{Ideation. Semantics-doc + the failable consumer-binding. No spike likely needed — this composes already-proven 8y behavior + the captain-articulated model — but record the determination.}

## Notes

- Reframes the cluster: this is FOUNDATIONAL; `p4` slims to the citation gate (extends `2a`), the authoring primitive becomes its own entity, and `4n` keys cost by the scenario tuple — all citing these principles.
- Sequencing: independent of the `at → 2a → n3` critical path; can ideate in parallel. `4n` is the codex peer's lane — these principles set the contract it adopts (coordinate via state).
