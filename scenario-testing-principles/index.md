---
id: m8yy9wbsm1b832a5078fgyza
title: Scenario-testing principles — a scenario is a natural-language spec; codified test and LLM run are two executor implementations; cost/coverage key by scenario × {mode, runtime}
status: implementation
source: captain (2026-06-03, session 12) — "we should have scenario as natural language, like bdd. and codified test is a variant of an implementation. running through llm is another implementation." + "document this as scenario-testing principles first, the semantics not the syntax. 8y should have concluded with a few scenarios; we build upon it. 4n is tracking cost, which should be keyed by the scenario (with mode, runtime as variants)."
started: 2026-06-03T16:55:07Z
completed:
verdict:
score: 0.41
worktree: .worktrees/spacedock-ensign-scenario-testing-principles
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

## The failable proof (how this escapes "a document about itself")

The README bar: an AC whose only proof is reviewing the doc's own prose can never fail. The teeth here is that the doc is **bound to a code source-of-truth that already exists** — `8y`'s `sharedRuntimeScenarios()` table (`internal/ensigncycle/shared_scenarios_test.go`) is the single, tested registry of the seed scenario IDs (`gate-guardrail`, `rejection-flow`, `merge-hook-guardrail`), already pinned by `TestSharedRuntimeScenarioDefinitions` and held in host parity by `TestSharedScenarioRunnerCoverage`. The doc declares those same IDs in a machine-readable block; a **lock test compares the doc's declared IDs against `sharedRuntimeScenarios()` and reds on any divergence**. This is not a prose-presence check — it is an invariant over two real value sources (the spec's declared IDs and the Go table). If `8y` adds/renames a scenario and the doc isn't updated, or the doc names a scenario the code doesn't have, the test fails in either direction (proven below). The doc becomes the human-readable face of a code-backed truth, the way `8y`'s `TestSharedScenarioDocsContract` makes `docs/dev/README.md` the face of the runner contract — the exact pattern `docs/specs/state-behavior-extension.md` lacks (it is prose-only, bound to nothing; this entity must not repeat that).

Two consumer bindings carry the model into the entities that depend on it. They are named here as the contract; the consumers themselves are separate entities/lanes that cite this doc and add the matching test in their own lane:
- **`4n` cost-key binding** — `4n`'s `Record` (`internal/journeymetrics/types.go`) today keys cost by a flat `JourneyID` + `Host` + `Model`; its cycle-2 captain rejection mandates re-keying by the **shared scenario ID** with mode/runtime/host/model as variants (it had used host-prefixed `claude-gate-guardrail`, which these principles forbid). The binding: `4n`'s scenario key MUST be drawn from `sharedRuntimeScenarios()`, not a host-prefixed re-derivation. `4n` lands the test (a host-prefixed scenario key reds it); these principles define the key it adopts.
- **`p4` citation-gate binding** — `p4` cites the executor model to justify why a runtime-observable AC (only the LLM executor can decide it) needs a cited live run; `p4`'s gate is the enforcement of the "two executors of one scenario" model, not a re-derivation of it.

## Acceptance criteria

Each AC names an end-state property of the finished doc and a check **outside the doc** that can fail.

**AC-1 — The doc formalizes the scenario/executor/variant-axes semantics in a machine-readable form bound to the seed-scenario source-of-truth, and the binding fails on divergence.**
Verified by: a Go test `TestScenarioPrinciplesSeedLock` in `package ensigncycle` (reachable to `sharedRuntimeScenarios()`, which is `_test.go`-internal to that package — so the lock lives there, beside `TestSharedScenarioDocsContract`) reads `docs/specs/scenario-testing-principles.md`, extracts the seed scenario IDs the doc declares in a `<!-- seed-scenarios -->` block, and asserts that set `reflect.DeepEqual`s `sharedRuntimeScenarios()` names. RED when the doc drops, adds, or renames a scenario relative to the code table (both directions). This is an invariant over two real value sources, not a substring search. Spike below exercised it: matching → PASS; dropped scenario → FAIL; bogus extra scenario → FAIL.

**AC-2 — The doc's required semantic clauses are present at the claim's own level (the floor, not the load-bearing proof).**
Verified by: a Go presence check (same file/pattern as `TestSharedScenarioDocsContract`) that the doc carries the named terms — `scenario` defined as a natural-language behavioral spec graded on durable outcomes (entity state before→after, archive state, on-disk artifacts, durable user-facing output) and **never transcript phrasing**; the two executor kinds (codified — deterministic Go test, proves the modeled/consumer logic, cannot prove the producer; LLM — real agent run, proves the real producer, graded on durable outcomes); the variant axes `scenario × {mode, runtime}`; and the named seed instances. Legitimate as a presence check because the claim *is* the text — but it is the floor; AC-1's lock is what makes the doc fail on real drift.

**AC-3 — The doc states the cluster contract that its consumers key to: the `(scenario, mode, runtime)` cost/coverage tuple and the executor model the citation gate enforces.**
Verified by: the same presence check asserts the doc declares the `(scenario, mode, runtime)` tuple as the unit that is run, graded, and measured — the key `4n` adopts (forbidding host-prefixed scenario IDs like `claude-gate-guardrail`) — and the executor distinction `p4` cites. The *failable* edge of this AC is realized in the consumers' lanes: `4n`'s re-key test (host-prefixed key reds) and `p4`'s gate, both citing this doc. Within this entity, AC-3 is proof at the text's level; its teeth bite in `4n`/`p4`, which is why those bindings are named explicitly above rather than left as "a consumer somewhere honors this."

## Test plan

- **AC-1 (load-bearing lock, Go):** `TestScenarioPrinciplesSeedLock` in `internal/ensigncycle`, reading `../../docs/specs/scenario-testing-principles.md` via the proven `os.Getwd()` + relative-join path (`TestSharedScenarioDocsContract` already does this) and `reflect.DeepEqual`-comparing the doc's `<!-- seed-scenarios -->` IDs to `sharedRuntimeScenarios()`. Adversarial guard already run in the spike: drop a scenario from the doc → red; add a bogus one → red. Cost: trivial; no model, no network.
- **AC-2 / AC-3 (floor presence checks, Go):** assert the required semantic + contract clauses over the real doc text, modeled on `TestSharedScenarioDocsContract`'s specific-clause set (concrete phrases, not vague topic mentions). Cost: trivial.
- **Consumer bindings (other lanes, named not authored here):** `4n` adds the scenario-key test (a host-prefixed key reds); `p4` adds its citation gate. Both cite this doc. This entity does not author them — it defines the contract and lands AC-1's lock so the seed IDs cannot silently drift.
- **No spike needed beyond the one already run.** The design composes three already-proven mechanisms: (1) `8y`'s `sharedRuntimeScenarios()` table + its meta/parity/docs tests; (2) the in-package doc-reading path (`TestSharedScenarioDocsContract` reads a markdown file from `package ensigncycle` via `os.Getwd()`); (3) the captain-articulated semantic model. The one genuinely-unverified mechanism — *can a `package ensigncycle` test read a `docs/specs/*.md` doc, extract scenario IDs, and lock them to `sharedRuntimeScenarios()` such that drift reds it?* — was the riskiest unknown (a doc-bound-to-code lock is what separates this from a prose-only spec), so it was exercised first. See the Spike result below.

## Spike (riskiest unknown — RUN at ideation)

**The riskiest unknown for a *principles* entity is the failable binding itself**: a doc that only asserts prose cannot fail. The unexercised mechanism the whole design rests on is whether a `package ensigncycle` test can read a `docs/specs` markdown doc, extract the seed scenario IDs from a machine-readable block, and lock them against the unexported-but-in-package `sharedRuntimeScenarios()` so that divergence reds in both directions.

Exercised with a throwaway `TestSpikeScenarioPrinciplesSeedLock` in `internal/ensigncycle` reading a fixture `docs/specs/scenario-testing-principles.md` with a `<!-- seed-scenarios -->` block:

| Doc state vs. `sharedRuntimeScenarios()` | Expected | Result |
|---|---|---|
| doc lists all three seed IDs | PASS | `ok ... ensigncycle` |
| doc DROPS `merge-hook-guardrail` | FAIL | `doc IDs [gate-guardrail rejection-flow] != code IDs [...merge-hook-guardrail...]` |
| doc ADDS `bogus-not-in-code` | FAIL | `doc IDs [bogus-not-in-code ...] != code IDs [...]` |

The lock discriminates both directions; the spike artifacts were removed. This is the failable proof seed for AC-1's implementation test. Everything else composes already-proven behavior, so **no further spike is needed**.

## Cluster relationships and sequencing

This entity is FOUNDATIONAL: it is the single shared model the cluster cites instead of re-deriving.

- **`p4` (live-verification-gate)** slims to the **citation gate** (extends `2a`'s `ClassifyEntityACs`). It cites this doc's executor model — a runtime-observable AC is one only the LLM executor can decide — to justify why such an AC needs a cited live run. `p4` does not redefine scenario/executor; it enforces them.
- **The authoring primitive becomes its OWN entity** (deferred; out of scope here). It makes a scenario authorable outside the `ensigncycle` test package — generalizing `8y`'s `_test.go`-buried triple — and decides the SYNTAX (Gherkin vs runbook, file layout, Go API), which these principles explicitly leave open. That entity is where `sharedRuntimeScenarios()` likely graduates from `_test.go`-internal to an importable registry; until then, AC-1's lock lives in-package beside the table.
- **`4n` (journey-cost-ledger, codex peer's lane)** keys cost by the `(scenario, mode, runtime)` tuple, scenario drawn from the shared seed IDs — its cycle-2 captain rejection already mandates exactly this re-key. **Coordinate via shared state; do not author `4n`.** These principles define the key; `4n` adopts it and lands the host-prefixed-key-reds test in its own lane.

**Sequencing:** independent of the `at → 2a → n3` critical path; ideated in parallel. AC-1's lock has no dependency (`8y`'s table is merged on this root). `p4` sequences after `2a` lands (its own dependency, not this doc's). `4n` is in flight and re-keying now.

**No-spike determination:** recorded in the Spike section above. The design composes already-proven behavior (`8y`'s scenario table + meta/parity/docs tests; the in-package doc-reading path from `TestSharedScenarioDocsContract`; the captain-articulated model). The single genuinely-unverified mechanism — a doc-bound-to-code lock that reds on drift — was exercised first because it is what makes this a failable spec rather than prose-only. No further spike needed.

**Deliverable home:** `docs/specs/scenario-testing-principles.md` (semantics only), plus `TestScenarioPrinciplesSeedLock` + the presence checks in `internal/ensigncycle/`.

## Stage Report: ideation

- DONE: Formalize the SEMANTICS only (syntax explicitly deferred): scenario = a natural-language behavioral spec graded on DURABLE outcomes ... Executors = interchangeable implementations ... Variant axes ... Name 8y's three shipped scenarios as the seed instances.
  Proposed-approach points 1–4 carry the model; durable-outcomes-not-transcript, codified-vs-LLM (consumer vs producer fidelity), `scenario × {mode, runtime}`, and the three seed IDs are stated. Syntax (Gherkin/runbook/layout/API) is kept in Out of scope as a follow-on entity.
- DONE: Find the FAILABLE proof so this is NOT 'a document about itself' (README bar). Each AC must be decided by something OUTSIDE the doc that can fail ... Prefer a consumer-binding gate over a prose-presence check. Spell out which consumer makes it failable.
  AC-1 is a Go lock test (`TestScenarioPrinciplesSeedLock`) binding the doc's `<!-- seed-scenarios -->` IDs to `8y`'s `sharedRuntimeScenarios()` via `reflect.DeepEqual` — an invariant over two real value sources, not a substring search; spiked PASS/drop-FAIL/add-FAIL (table in Spike). AC-2/AC-3 are floor presence checks; their teeth bite in the named consumer bindings — `4n`'s scenario-key re-key (host-prefixed key reds) and `p4`'s citation gate.
- DONE: State the cluster relationships + sequencing: p4 slims to the citation gate (extends 2a), the authoring primitive becomes its OWN entity, and 4n keys cost by the (scenario, mode, runtime) tuple — all CITING these principles. Coordinate with the codex peer on 4n via shared state (do not author 4n). Record the no-spike determination or name what to spike.
  "Cluster relationships and sequencing" section states all three. `4n` named as codex peer's lane (coordinate, do not author); its cycle-2 rejection already mandates the shared-scenario re-key. No-spike determination recorded: composes proven `8y` behavior + the in-package doc-read path + the captain model; the one unverified mechanism (doc↔code lock) was exercised first.

### Summary

Refined the seed into a gate-ready ideation. The hard part — the failable proof — is solved by binding the doc to `8y`'s already-tested `sharedRuntimeScenarios()` table: a `package ensigncycle` lock test compares the doc's declared seed IDs to the Go table and reds on drift in either direction (spiked: PASS / dropped-FAIL / added-FAIL). This is what separates it from `docs/specs/state-behavior-extension.md` (prose-only, bound to nothing). Syntax is explicitly deferred to a follow-on authoring-primitive entity; `4n` (codex lane) and `p4` are named as the consumer bindings that carry the model's teeth. No further spike needed.

## Stage Report: implementation

- DONE: Write `docs/specs/scenario-testing-principles.md` capturing the SEMANTICS verbatim-in-spirit — scenario = natural-language behavioral spec graded on durable outcomes (never transcript phrasing); two interchangeable executors (codified = consumer logic, LLM = real producer); variant axes scenario × {mode, runtime}; the three seed instances; the consumer model (citation gate, authoring primitive, cost ledger keyed by the tuple).
  Doc at `docs/specs/scenario-testing-principles.md`, committed on worktree branch `spacedock-ensign/scenario-testing-principles` as `ffb829b8`. Includes a machine-readable `<!-- seed-scenarios -->` block whose IDs were exercised against `sharedRuntimeScenarios()` and matched (gate-guardrail, rejection-flow, merge-hook-guardrail).
- DONE: Match `docs/specs/` house style (read `state-behavior-extension.md`); markdown only, no ABOUTME comments; tight and reference-grade.
  Followed the read reference's structure: top `Status:` line, short purpose, sectioned semantics, fenced `text` blocks. ~69 lines, ~1.5 pages. Syntax kept explicitly out of scope per captain direction.
- DONE: Doc-only deliverable — no tests/ACs/validation gate; note doc path + no-ceremony design doc per captain direction.
  No-ceremony doc-write per captain-directed dispatch (doc-only entities banned, ceremony collapsed). The entity's AC-1 lock test (`TestScenarioPrinciplesSeedLock`) is intentionally NOT authored here; the `<!-- seed-scenarios -->` block is present so the p4/primitive lanes can bind a lock test to it later.

### Summary

Distilled the ideation body into a concise, reference-grade design doc at `docs/specs/scenario-testing-principles.md` (semantics only; syntax deferred). The doc carries the scenario/executor/variant-axes model, the three seed scenarios mirroring `sharedRuntimeScenarios()` (verified matching), and the consumer contract the citation gate / cost ledger / authoring primitive cite. Per captain's no-ceremony direction this is doc-only: no lock test or presence checks were authored, but the machine-readable seed block is in place so a future lane can bind to it. FO merges the doc directly — no PR, no validation gate.
