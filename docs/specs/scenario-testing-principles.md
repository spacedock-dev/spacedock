# Scenario-Testing Principles

Status: design foundation (semantics).

This document defines the semantic model for scenario-based testing in Spacedock: what a *scenario* is, what an *executor* is, and the *variant axes* a scenario runs under. It is the single shared model that downstream consumers — the citation gate, the cost ledger, and the authoring primitive — cite instead of re-deriving.

**Semantics only.** Concrete syntax — Gherkin vs. freer runbook+assertions, file format, the authoring primitive's Go API — is explicitly out of scope and deferred to a follow-on entity.

## Scenario

A **scenario** is a natural-language behavioral specification, in the spirit of BDD. It names a behavior — a user-or-agent journey — and its expected **durable outcomes**, independent of *how* the outcomes are checked.

Durable outcomes are facts that survive the run:

- entity state before → after,
- archive state,
- on-disk artifacts,
- durable user-facing output obligations.

Assertions are over durable outcomes only — **never transcript phrasing**. Transcript wording is non-deterministic and is not a behavioral fact; grading on it makes a scenario brittle and meaningless across executors.

## Executors

An **executor** is an implementation of a scenario's check. Executors are interchangeable: the scenario is the spec, executors are pluggable implementations of it at different fidelity. Two are defined here; the model is open to further executors.

- **Codified executor** — a deterministic test (Go fixture/unit). Cheap, fast, always-on in CI. It proves the *modeled* logic — the **consumer** side (e.g. a watcher over a recorded stream). It cannot prove the producer: a recording is frozen, so the agent that produced it is not under test.

- **LLM executor** — a real agent run (Claude / Codex). Expensive and non-deterministic. It proves the real **producer** — the agent actually behaving — graded on the scenario's durable outcomes.

The two executors check the *same* scenario at different fidelity. Treating "offline test" and "live run" as different *kinds of thing* — rather than two executors of one scenario — is the root of the recurring failure mode where an offline proof passes while the live run fails.

## Variant Axes

A scenario executes as a tuple:

```text
scenario × {mode, runtime}
```

- **runtime** — the host (claude / codex). Per-host runner adapters implement this axis.
- **mode** — the model / effort (e.g. sonnet / opus).

The `(scenario, mode, runtime)` tuple is the unit that is **run, graded, and measured**. It is the key the cost ledger uses, and the key a coverage matrix uses to demand parity across hosts and models.

## Seed Scenarios

The first foundation is the three host-neutral runtime scenarios already shipped and held in host parity by the shared coverage tests. They are the named seed instances:

<!-- seed-scenarios -->
- `gate-guardrail` — the FO halts at a human gate and presents the review without self-approval, mutation, or archival.
- `rejection-flow` — the FO observes a rejected validation report and routes the concrete finding back through implementation.
- `merge-hook-guardrail` — the FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force.
<!-- /seed-scenarios -->

These IDs are the code-backed source of truth. They mirror the `sharedRuntimeScenarios()` table in `internal/ensigncycle`; the seed IDs declared above must equal that table. This block is machine-readable so a lock test can bind the doc to the code and red on drift in either direction — adding, dropping, or renaming a scenario on one side without the other. This is what makes the doc the human-readable face of a code-backed truth rather than prose bound to nothing.

## The Model Consumers Build On

These principles are the contract the cluster keys to. The consumers are separate entities/lanes that cite this doc and land their own enforcing tests.

- **Citation gate** — a runtime-observable acceptance criterion is one only the LLM executor can decide. It is satisfied only by a cited LLM-executor run; the codified executor alone does not satisfy it. This stops "default to the cheap proxy" by construction — the gate enforces the two-executors-of-one-scenario model rather than re-deriving it.

- **Authoring primitive** — makes a scenario authorable outside the test package, so writing a scenario is as cheap as a Go test. This is the entity that decides the deferred syntax.

- **Cost ledger** — keys cost by the `(scenario, mode, runtime)` tuple, with the scenario ID drawn from the shared seed IDs. Host-prefixed scenario keys (e.g. `claude-gate-guardrail`) are forbidden: the scenario identity is host-neutral, and `runtime` is a variant axis, not part of the scenario name. This makes cost attributable per scenario and comparable across variants.

## Why This Dissolves the Offline/Live Split

Write the scenario once. Run the codified executor every CI to prove the modeled logic, and run the LLM executor when producer fidelity is required to prove the real agent — with the citation gate forcing the LLM run to be cited for any producer claim. The offline pass and the live pass are no longer different kinds of evidence about different things; they are two executors of one named scenario, and the gate names which one a given claim requires.
