---
title: Make review-finding disposition workflow-owned and separate in-stage review from rejection routing
status: backlog
source: "Captain boundary audit, 2026-07-30: an end-of-implementation Roborev round never enters feedback-rejection-flow, and finding classification must happen when findings arrive rather than only after rejection."
started:
completed:
verdict:
score: "0.90"
worktree:
issue:
pr:
sprint: durable-decisions
id: rhx820qrkn6vxpday10nch36
gates:
    version: 1
    current:
        gate: gate:rhx820qrkn6vxpday10nch36:backlog
    records:
        - id: gate:rhx820qrkn6vxpday10nch36:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rhx820qrkn6vxpday10nch36-backlog-1
              briefing:
                id: briefing:rhx820qrkn6vxpday10nch36:backlog:attempt-1:revision-1
                digest: sha256:f5c4251b4617310097bca5e1d14a03f5d177df80b0e452751b510e9b8eae30de
                digest-domain: canonical-bytes
                request-digest: sha256:50ec6a9d980e9ca88e45a5c22e5f7e394b6c0d6a5e3044a15a4dbf617515fcc7
                room-ref: ./review/backlog/briefing-1
---

Make the active workflow the authority for review-finding policy. Reviewers report findings; workers retain evidence and propose materiality, ownership, and disposition; the First Officer authorizes the disposition before finding-driven product changes, commits, or review reruns; validators recommend gate verdicts; and the captain decides changes to scope, accepted value, thresholds, or acceptance criteria. Rejection routing carries an existing disposition to `feedback-to`; it neither defines nor repeats the policy.

The shared rejection flow should route an already-decided correction and must not define Material, Deferred risk, Polish, Needs decision, correct-but-disproportionate, value-AC, ideation-estimate, or 2× policy. The development workflow and reusable development template should own Roborev classification, task ownership routing, expected-surface calibration, and advisory round recording. Shared gate preparation and presentation should preserve workflow-declared finding categories rather than force development tiers.

## Problem

The shared `feedback-rejection-flow` currently carries development-specific finding classification even though it triggers only after a feedback gate rejects. Findings also arrive during implementation, validation, consequential quick work, detached audits, and routed feedback. Classification at rejection is therefore globally over-scoped and too late.

The current wording also blurs two separate questions: whether a finding matters and whether the active task owns its remedy. A worker may gather read-only evidence and prepare a rebuttal before consultation, but a Material finding outside the task's approved semantic boundary is `Needs decision`, not permission to expand the task. In-stage round Resolutions remain advisory; the First Officer's disposition authorizes or stops current work, while only an actual gate Resolution binds scope or status.

## Proposed approach

Ideation should define one workflow-owned policy applied wherever findings arrive:

- reviewers report findings and trigger evidence without deciding task ownership;
- workers preserve evidence and propose materiality, ownership, and disposition, consulting the First Officer before any finding-driven product edit, commit, or review rerun;
- validators use the authorized disposition when recommending `PASSED` or `REJECTED`, without becoming a second binding classifier;
- the First Officer authorizes the disposition for the current work and routes `Needs decision` findings without silently expanding scope;
- the captain alone changes scope, accepted value, thresholds, or acceptance criteria;
- cross-stage rejection routing forwards the existing disposition and revise decision without reclassification.

Generic skills and commission scaffolding should refer to the active workflow's declared review policy without naming development classes. The development README and reusable template should carry the full Roborev taxonomy, ownership rule, expected-surface and semantic boundary, tolerance, and AC-drift behavior. Preserve the existing four-field release-scope evidence test. Its third field—the affected value AC or non-negotiable boundary—must cite the authority for that boundary, such as an acceptance criterion, captain ruling, or named contract; do not invent a fifth field.

## Out of scope

Do not redesign the advisory-round recorder's storage validation in this task; that is owned by the sibling workflow-neutral recorder task. Do not redesign the generic gate criteria-source interface or dispatch stage identity.
