---
title: Deliver one portable common live-journey surface
status: ideation
source: "Captain recarve of live-test-truth, 2026-08-03. Absorbs 3w, h3, tj, and r4 as design inputs."
score: 1.0
sprint: live-test-truth
group: portable-common-surface
sprint-readiness: ready
id: ys7ncwh9kr8w5h9hdkz5apat
gates:
    version: 1
    records:
        - id: gate:ys7ncwh9kr8w5h9hdkz5apat:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ys7ncwh9kr8w5h9hdkz5apat-backlog-1
              briefing:
                id: briefing:ys7ncwh9kr8w5h9hdkz5apat:backlog:attempt-1:revision-1
                digest: sha256:ff92cf6104cdc81f88ab51f6a59f33f3d1f42e72d7e72f712a2f50d5ae61fc47
                request-digest: sha256:f9af92f7f1823701ddcf301aadddf68012f7af8eb4656ea6a6ced3772ec14e48
                room-ref: ./deliver-portable-common-live-journey-surface/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ys7ncwh9kr8w5h9hdkz5apat:backlog:1
                briefing: briefing:ys7ncwh9kr8w5h9hdkz5apat:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T12:18:32.293007Z"
                decision: approve
                reason: Captain explicitly approved the outcome-shaped recarve and directed immediate redispatch.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

## Outcome

A test operator can select one common journey by one stable identity on Claude, Codex, or Pi. The registry makes missing journeys and orphan fixtures visible.

The task delivers the complete operator journey. Registry annotations, fixture bindings, runtime adapters, selectors, and reconciliation are implementation steps.

## Design inputs

- `bind-live-registry-to-source` (`3w`): annotation grammar, source join, and reconciliation guard spike.
- `converge-shared-live-suite-entrypoint` (`h3`): one selector passed on Claude and Codex.
- `add-pi-common-live-runner` (`tj`): the Pi shallow-boot spike passed and measured cost.
- `promote-standalone-common-live-journeys` (`r4`): six missing common journeys and fixture overlap.

Preserve these reports as evidence. Do not repeat their live spikes unless a source change invalidates them.

## Acceptance criteria

**AC-1 (VALUE) — One journey identity works on every supported runtime.**
Verified by: the same `TestLiveSharedScenarios/<journey-id>` selector runs unchanged on Claude, Codex, and Pi for one representative journey.

**AC-2 — All desired common journeys have canonical executable identities.**
Verified by: reconciliation reports no missing binding for the 16 registry journeys, except an explicit desired-state gap that the registry still names.

**AC-3 — Fixtures have accountable use.**
Verified by: each registered fixture ID resolves to its source builder and at least one journey. The reconciliation result names each orphan fixture.

**AC-4 — Runtime differences stay behind adapters.**
Verified by: common fixtures and assertions contain no host-specific branch. CI lanes select the same canonical entry point.

## Ideation requirements

Use the four design-input reports. Produce one implementation plan, one collision map, and one landing sequence. Reconcile overlaps instead of repeating component plans.

Use `$simple-english` in pragmatic mode for the complete plan.

