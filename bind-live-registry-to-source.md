---
title: Bind the desired live-test registry to tests and fixture builders
status: backlog
source: "Captain semantic lock for docs/runtime-live-ci-registry.md, 2026-08-03"
started:
completed:
verdict:
score: 0.95
worktree:
issue:
sprint: live-test-truth
group: registry
sprint-readiness: ready
id: 3w2rx3aw4vcympx84zt8mtv7
---

## Problem

The desired live-test registry now names common journeys, runtime-specific proofs, canonical entry points, and fixture identities, but source declarations do not yet carry stable journey and fixture bindings. A reader cannot reliably move from a registry promise to its test and builder, and unused fixtures remain hard to distinguish from deliberate helpers.

## Acceptance criteria

**AC-1 (VALUE)** Every registered journey and runtime-specific proof resolves to a concrete source binding, and every registered fixture resolves to a concrete builder or an explicitly unimplemented desired binding.
Verified by: a one-off inventory at the candidate tip that joins registry IDs to annotated test, scenario, and fixture declarations and reports no ambiguous bindings.

**AC-2** docs/runtime-live-ci.md names docs/runtime-live-ci-registry.md as the normative desired-state registry and records the required reconciliation SHA discipline without duplicating the registry entries.
Verified by: the rendered documentation and a path-scoped git diff from the recorded SHA.

**AC-3** Live-tagged tests and fixture builders outside the registry are either registered, moved to default-tagged coverage, explicitly classified as non-gating experiments, or removed.
Verified by: the same one-off inventory with each discovered symbol accounted for.

## Stage-specific test gates

- Ideation must choose the smallest annotation grammar and must not assume a permanent lint or generator without separate captain consent.
- Implementation updates source annotations and documentation together.
- Validation runs the one-off inventory and the repository-required Go tests.

