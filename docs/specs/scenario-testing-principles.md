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

A scenario produces one or more observations. Each observation is described by the scenario plus its execution axes:

```text
scenario × {executor, mode, runtime, model?}
```

- **executor** — the implementation of the check (`codified` / `llm`).
- **mode** — the evidence path used for the observation. Current modes are `codified` for deterministic Go/fixture checks and `llm-live` for a real model-backed agent run. Future modes may add replay or simulation, but a mode is never a model name.
- **runtime** — the host or runner surface (`claude` / `codex`). Codified observations still carry runtime: they prove the modeled consumer behavior for that host. Live LLM observations prove the real producer for that same host.
- **model** — the actual model used when the executor invokes one, such as `sonnet`, `opus`, or `gpt-5-codex`. Codified observations usually leave this empty.

The baseline coverage matrix for each shared scenario is:

```text
runtime {claude, codex} × mode {codified, llm-live}
```

The `(scenario, mode, runtime)` tuple is the primary variant row that is **run, graded, and measured**. When more than one model can run under the same mode/runtime pair, `model` is a separate observation dimension; do not fold it into `mode`. This is the key shape the cost ledger and coverage matrix use to compare hosts, evidence paths, and models without changing the scenario identity.

## Seed Scenarios

The first foundation is the three host-neutral runtime scenarios already shipped and held in host parity by the shared coverage tests. They are the named seed instances:

<!-- seed-scenarios -->
- `gate-guardrail` — the FO halts at a human gate and presents the review without self-approval, mutation, or archival.
- `rejection-flow` — the FO observes a rejected validation report and routes the concrete finding back through implementation.
- `merge-hook-guardrail` — the FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force.
<!-- /seed-scenarios -->

These IDs are the code-backed source of truth. They mirror the `sharedRuntimeScenarios()` table in `internal/ensigncycle`; the seed IDs declared above must equal that table. This block is machine-readable so a lock test can bind the doc to the code and red on drift in either direction — adding, dropping, or renaming a scenario on one side without the other. This is what makes the doc the human-readable face of a code-backed truth rather than prose bound to nothing.

## Prioritizing New Cross-Runtime Scenarios

A scenario belongs in the shared cross-runtime set only when the same user journey should hold for every supported host. Host-specific probes, such as a Codex-only idle-notification experiment, should live in a host-runtime test lane instead of `sharedRuntimeScenarios()`.

Prefer new scenarios that start from the first officer's normal event loop: boot the workflow, inspect startup state and dispatchable work, run the required action, and grade durable outcomes. This catches real producer behavior instead of only testing a helper in isolation.

Prioritize these next shared scenarios:

- `pr-lifecycle-from-boot` — boot a workflow, observe PR-pending and dispatchable state, let the PR lifecycle advance, and verify the entity reaches the correct durable state without bypassing merge hooks or archival rules.
- `dispatch-before-wait` — with ready workflow work available, the first officer finishes dispatch/advance/gate handling before waiting for idle worker completion.
- `split-root-bootstrap-resume` — from a fresh or stale split-root checkout, the first officer halts on missing state, initializes or pulls state, then resumes dispatch only after the entity directory is present.
- `live-metrics-artifact-capture` — live Claude and Codex scenario runs emit raw per-run journey metrics and archive them as CI artifacts before any aggregation or post-processing.
- `feedback-reuse-boundary` — after a rejected validation result, the first officer routes the concrete finding back to the implementation stage, chooses reuse or fresh dispatch by the runtime contract, and records the second-cycle durable outcome.

## The Model Consumers Build On

These principles are the contract the cluster keys to. The consumers are separate entities/lanes that cite this doc and land their own enforcing tests.

- **Citation gate** — a runtime-observable acceptance criterion is one only the LLM executor can decide. It is satisfied only by a cited LLM-executor run; the codified executor alone does not satisfy it. This stops "default to the cheap proxy" by construction — the gate enforces the two-executors-of-one-scenario model rather than re-deriving it.

- **Authoring primitive** — makes a scenario authorable outside the test package, so writing a scenario is as cheap as a Go test. This is the entity that decides the deferred syntax.

- **Cost ledger** — groups cost by scenario and variant axes, with the scenario ID drawn from the shared seed IDs. Host-prefixed scenario keys (e.g. `claude-gate-guardrail`) are forbidden: the scenario identity is host-neutral, `runtime` is a variant axis, and `model` is separate from `mode`. This makes cost attributable per scenario and comparable across hosts, modes, and models.

## Why This Dissolves the Offline/Live Split

Write the scenario once. Run the codified executor every CI to prove the modeled logic, and run the LLM executor when producer fidelity is required to prove the real agent — with the citation gate forcing the LLM run to be cited for any producer claim. The offline pass and the live pass are no longer different kinds of evidence about different things; they are two executors of one named scenario, and the gate names which one a given claim requires.
