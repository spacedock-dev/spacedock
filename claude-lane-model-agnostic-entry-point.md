---
title: The Claude live lane selects a model by cadence, so a new model needs a workflow edit
status: ideation
source: "Captain, 2026-10-07: a manual run on claude-haiku-5-5 was requested and cannot be dispatched. The captain wants opus and sonnet to be a parameter, not a mode, and the underlying entry point to be model-agnostic."
id: f77r1rj2g7sgj3hkrht7n37f
gates:
    version: 1
    records:
        - id: gate:f77r1rj2g7sgj3hkrht7n37f:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:f77r1rj2g7sgj3hkrht7n37f-backlog-1
              briefing:
                id: briefing:f77r1rj2g7sgj3hkrht7n37f:backlog:attempt-1:revision-1
                digest: sha256:431932767c9eefcead1a0540c5d3118b8f86f71c1f6fa6e6d552029da7f51eae
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:f77r1rj2g7sgj3hkrht7n37f:backlog:1
                briefing: briefing:f77r1rj2g7sgj3hkrht7n37f:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-10-07T18:14:48.37703Z"
                decision: approve
                reason: 'Captain, 2026-10-07: "file a task to reorg how we assign model to the claude-code harness. opus/sonnet should be a param, not a mode... dispatch."'
              application:
                target-stage: ideation
                state: consumed
---

The Claude live lane has no model input. `claude_version` pins the CLI version, not the model. The model comes from `internal/release/live_models.txt`, and the `live_cadence` input selects between the entries. The lane's matrix maps one cadence to one model and one environment.

So a new model needs a workflow edit before anyone can run it. That happened today: `claude-haiku-5-5` landed, and no dispatch could exercise it.

## Proposed direction

Keep `internal/release/live_models.txt` as the single source. Make the model a parameter of the entry point, and keep the cadence mapping as a convenience over that parameter. Then any named model runs without a workflow change.

The CI environments stay as they are. The four environments exist because secrets differ, not because models differ.

## Out of scope

No new lane. No change to the environment approval model. No change to the Codex or Pi lanes beyond what the same parameterization requires.

## Evidence to gather during ideation

- the exact shape of the workflow matrix line, currently `claude-sonnet-5-5, max, CI-E2E`
- how `live_models.txt` reaches the workflow, and whether a model parameter can be passed without breaking the single-source rule
- whether the detached adversarial audit applies, since the CI and release machinery is one of the declared high-stakes surfaces
