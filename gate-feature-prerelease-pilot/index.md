---
title: Cut and field-test the current gate feature as a prerelease
status: backlog
score: 1.0
sprint: durable-decisions
source: "Captain direction on 2026-07-30 to release the current gate feature before sprint closure and use real installations to discover friction."
id: 0hympdejewzwkhe60ygqk15a
gates:
    version: 1
    current:
        gate: gate:0hympdejewzwkhe60ygqk15a:backlog
    records:
        - id: gate:0hympdejewzwkhe60ygqk15a:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0hympdejewzwkhe60ygqk15a-backlog-1
              briefing:
                id: briefing:0hympdejewzwkhe60ygqk15a:backlog:attempt-1:revision-1
                digest: sha256:8e001f16c030b80444488bc79a581b04da24ec4f9b9aee31bfa734a88525caa2
                digest-domain: canonical-bytes
                request-digest: sha256:519a9e44c1703eb2a751839e98755b9710a25e4b43781f37827045f8547ebcff
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0hympdejewzwkhe60ygqk15a:backlog:1
                briefing: briefing:0hympdejewzwkhe60ygqk15a:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T14:03:29.385661Z"
                decision: approve
                reason: Captain directed an instrumented prerelease before sprint closure; the bounded pilot preserves release integrity while deferring non-blocking friction.
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

## End value

Publish `v0.27.0-pre1` from the current releasable `main` gate surface, then exercise the installed artifact through real default-chat and provider-backed gate journeys so additional friction is observed from user and agent experience before the remaining sprint scope hardens the design.

## Boundary

This is an instrumented prerelease, not the durable-decisions sprint acceptance walk. It does not require every sprint member to be terminal, and it does not silently absorb defects found while piloting. Known limitations and newly observed friction are recorded with owner, severity, reproduction, and promotion condition; only failures that prevent publishing, installing, preparing a valid room, recording an authentic decision, or consuming it without corrupting authority block the prerelease.

## Acceptance criteria

**AC-1 — Published candidate.** An annotated `v0.27.0-pre1` names an exact `main` commit whose required format, full Go, race, install, and Runtime Live E2E gates passed, and the release assets and checksums install successfully in a clean environment.

**AC-2 — Default path.** A clean installed Spacedock drives one real chat gate through prepare, presentation, exact decision recording, validation, and one-use consumption without caller-authored JSON or direct state editing.

**AC-3 — Provider path.** The same installed Spacedock prepares a provider-neutral room consumed by the current Subspace gate entry; retained request, Briefing, Result, inventory, actor, digest, and final application are independently inspectable. A known provider rendering limitation may be retained as a filed finding when it does not compromise the recorded decision.

**AC-4 — Useful feedback.** The pilot report separates release blockers, sprint-critical integrity/operability defects, and deferrable convenience or generalization work. Each friction names the failing journey step, observable impact, evidence, owner, and next action; the release task itself fixes only release blockers.

## Explicit non-prerequisites

Completion of the final `ph` sprint walking skeleton, multi-Artifact preparation, generalized advisory-round recording, extra dashboards, compatibility with prototype formats, and resolution of every host-specific convenience issue are not prerequisites for this prerelease.
