# Architecture notes

Notes on how the Spacedock codebase and its design contracts are organized, for contributors orienting in the repo.

## Project shape

The repo builds the Go-based `spacedock` launcher and the project-side skill integration:

- `cmd/spacedock/` — the process entry point only.
- `internal/cli/` — command routing, usage text, exit-code behavior.
- `internal/status/` — status implementation and the compatibility runner.
- `docs/specs/` — design contracts.
- `docs/dev/README.md` — the development workflow definition.
- `docs/dev/.spacedock-state/` — the development workflow's entities, in a separate state checkout.
- `skills/` — the skill-facing command integration.

The baseline gate for every change is `go test ./...`. Prefer small Go packages with clear boundaries; use the standard library unless a dependency removes real complexity; keep command output stable and test it with fixtures.

## Design contracts

Two specs under `docs/specs/` carry the load-bearing design semantics:

- **State behavior extension** — the split-root storage profile (definition vs. runtime state) and the external-tracker bridge principles. See [External-tracker bridge](../advanced/external-tracker.md) and [Multi-workflow & split-root state](../advanced/split-root-state.md).
- **Scenario-testing principles** — the semantic model for scenario-based testing: what a scenario is, what an executor is (codified vs. LLM), and the variant axes a scenario runs under. It dissolves the offline/live split by treating both as two executors of one named scenario.

## Runtime live CI

The live lanes prove runtime behavior, not text shape. One host-neutral scenario table, per-host runner adapters implementing the same scenario IDs, and a parity guard that fails if a scenario exists for one host only. Static grep over workflow YAML or skill prose is never a substitute for launching the real host front door, observing its output, and checking the resulting workflow state. For the full discipline, see [Proof policy](proof-policy.md).
