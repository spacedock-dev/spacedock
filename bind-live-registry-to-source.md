---
title: Bind the desired live-test registry to tests and fixture builders
status: ideation
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
gates:
    version: 1
    records:
        - id: gate:3w2rx3aw4vcympx84zt8mtv7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3w2rx3aw4vcympx84zt8mtv7-backlog-1
              briefing:
                id: briefing:3w2rx3aw4vcympx84zt8mtv7:backlog:attempt-1:revision-1
                digest: sha256:31091557195718f78a0f4425d97cdb4fad7653b621456e030b063c5ff7e09fb1
                request-digest: sha256:6dbc7d42c596779e0d531fdd9a0dd73238546e9f3b6279f9ab3882dede2112dc
                room-ref: ./bind-live-registry-to-source/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3w2rx3aw4vcympx84zt8mtv7:backlog:1
                briefing: briefing:3w2rx3aw4vcympx84zt8mtv7:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T10:33:36.145597Z"
                decision: approve
                reason: Captain approved the prepared Sol ideation cohort with make it so.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
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

