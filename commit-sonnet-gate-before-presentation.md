---
title: Commit Sonnet gate before presentation
status: ideation
score: "0.90"
source: "n28 exact Claude default-headless finding, 2026-08-10"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: kky8pg7wc8xgb985epwss092
gates:
    version: 1
    records:
        - id: gate:kky8pg7wc8xgb985epwss092:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kky8pg7wc8xgb985epwss092-backlog-1
              briefing:
                id: briefing:kky8pg7wc8xgb985epwss092:backlog:attempt-1:revision-1
                digest: sha256:a45f5479086379c8c2fa589df282f75f9fd2e98afa472f24c553480dae7f0398
                request-digest: sha256:2315b198901a80e371d222f7a35ec38f11019774d5b8ef5237989e6c0fe5cf8d
                room-ref: ./commit-sonnet-gate-before-presentation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kky8pg7wc8xgb985epwss092:backlog:1
                briefing: briefing:kky8pg7wc8xgb985epwss092:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T18:44:24.797313Z"
                decision: approve
                reason: Captain created this active Sonnet owner to repair the target-external gate lifecycle defect.
              application:
                target-stage: ideation
                state: consumed
started: 2026-08-10T18:44:40Z
---
## Problem

The exact local Sonnet default-headless journey prepared an open validation gate but did not commit the gate state. The First Officer then read and presented the gate. The target reported `gate-hold-violation`.

## Value

After `gate prepare`, the Sonnet First Officer commits the durable gate state. It stops at the same clean open validation gate and dispatches no successor.

## Scope

- Repair the Sonnet gate lifecycle at the prepare-and-bind boundary.
- Use the exact n28 Claude artifact as the baseline.
- Own the Sonnet `default-headless-gate-stop` binding after n28 transfers it.
- Do not change n28 acknowledgment mechanics, Codex behavior, Pi, or the zero-discovery repair.
- Use local Sonnet subscription authentication before required PR CI.

## Acceptance criteria

- AC-1: The command order is successful `gate prepare`, then `state commit`, then structured reads and presentation.
- AC-2: The exact local Sonnet default-headless target first reports bound XPASS-green and then passes normally after binding removal.
- AC-3: The final entity is a clean open validation gate with no terminal fields and no successor dispatch.
- AC-4: A focused negative control rejects a missing, failed, or late state commit.
- AC-5: Full, race, format, registry, active-owner, and required exact PR checks pass. Pi remains skipped.

## Baseline evidence

- Released user and workflow: local-subscription Sonnet default-headless gate stop without Captain authority.
- Observable harm: the prepared gate is not proven durable before presentation.
- Value authority: `skills/fo-gate-lifecycle/SKILL.md` requires prepare, commit, structured reads, and presentation in that order.
- Trigger: `/tmp/n284-happy-claude/claude-shared-scenarios/default-headless-gate-stop/command.log` has gate prepare, reads, and presentation, but no state commit. The target observed `gate-hold-violation`.

## Ideation requirements

- Name exact files and gross/net estimate before product edits.
- Preserve the existing gate grammar and authority.
- Define one focused falsifier for the missing commit.
