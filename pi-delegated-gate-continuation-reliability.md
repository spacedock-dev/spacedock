---
title: "Make Pi delegated gate approval continue reliably through successor dispatch"
status: backlog
source: "se0 exact-tip Pi recorded-gate proof on 2026-07-28: one run invented the Briefing digest before reading it; the single retry presented the correct bound gate but stopped before recording or consuming the delegated approval."
score: 0.9
sprint: pi-ux
group: gate
sprint-readiness: ready
issue:
id: 9w59t6m1qc46hccd54p04z2j
---

## Problem

The supported Pi recorded-gate journey is not reliable under an explicit
delegated conn. At exact tip `ce4ac943`, the first retained run completed the
gate lifecycle but presented an invented Briefing digest before reading the
canonical value. The one authorized retry presented the exact bound Briefing
ID and digest and an actionable approval question, then stopped immediately:
it did not commit the bound package, record or consume approval, advance,
dispatch the successor, wait, or verify the durable handoff.

This is agent conduct, not an oracle or fixture defect. The strict oracle
correctly rejected both runs. Prior Pi executions passed, so the problem is
reliability across the full presentation-to-application boundary rather than
absence of a mechanism.

`TestLivePiRecordedGateLifecycle` is temporarily TODO under this ticket. Its
skip does not prove the capability and must be removed when this ticket lands.
Keep deterministic gate tests, Pi front-door smoke, and the Claude/Codex
recorded-gate journeys active.

## Acceptance criteria

**AC-1 (VALUE)** Repeated clean Pi recorded-gate journeys under delegated conn
present the exact bound Briefing ID and digest, record and consume the approval,
dispatch exactly one successor on the requested model, wait, and verify the
durable successor report before yielding.

**AC-2** The remedy uses the canonical retained gate package and existing
provider-neutral lifecycle. It does not hardcode fixture values, weaken the
oracle, add retries, or introduce a Pi-only gate protocol.

**AC-3** The exact retained counterexamples remain red: wrong digest before
decision; correct presentation followed by an early stop; child/tool-only
presentation; decision before presentation; and missing successor evidence.

**AC-4** Remove the linked TODO and run the focused Pi recorded-gate journey
repeatedly plus the registered Pi live package at the exact candidate tip.

## Boundary

Coordinate with `w5bfnrvpcphw857nzz93340c`, which owns exact-digest
presentation reliability, and `gqsw81ghf48hr2n3jg6k7nx8`, which owns dispatch
after a gate has been consumed. This ticket owns the Pi agent's reliable
end-to-end conduct across both obligations under delegated conn. Prefer the
smallest authoritative contract, skill-wiring, or presentation seam; do not
teach the fixture transcript or add compatibility behavior.

## Evidence

- `live-artifacts/local-proof/final-ce4ac943-pi-recorded-gate-pinned/`:
  lifecycle completed, but the root review presented an invented digest and
  the transcript acknowledged the invention.
- `live-artifacts/local-proof/final-ce4ac943-pi-recorded-gate-pinned-retry/`:
  exact review presented, then the trace ended with zero child sessions and no
  approval application.
