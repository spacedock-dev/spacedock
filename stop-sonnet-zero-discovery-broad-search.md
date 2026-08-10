---
title: Stop Sonnet zero-discovery broad search
status: ideation
score: "0.90"
source: "PR #663 Sonnet zero-discovery failure, 2026-08-10"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: 3rns0vh3svq49w43cfr0wdqd
gates:
    version: 1
    records:
        - id: gate:3rns0vh3svq49w43cfr0wdqd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3rns0vh3svq49w43cfr0wdqd-backlog-1
              briefing:
                id: briefing:3rns0vh3svq49w43cfr0wdqd:backlog:attempt-1:revision-1
                digest: sha256:065ca5a81f4f0358e6bf789fba0888899021fb8a1236a79d56a9cb8504ae8ab2
                request-digest: sha256:6478d6c2a4d6f48380b4bbed9541397c34fef4b40ff9afabb286da733bf8f069
                room-ref: ./stop-sonnet-zero-discovery-broad-search/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3rns0vh3svq49w43cfr0wdqd:backlog:1
                briefing: briefing:3rns0vh3svq49w43cfr0wdqd:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T18:38:26.692525Z"
                decision: approve
                reason: Captain directed immediate ideation for the exact Sonnet zero-discovery product repair.
              application:
                target-stage: ideation
                state: consumed
---
## Problem

The exact Sonnet zero-discovery journey still broad-searches the filesystem at boot. In PR #663 run 31415550991, job 93543929461, artifact 9074747236, the First Officer ran `find` across the installed First Officer references. The journey failed in 24.73 seconds.

## Value

A Sonnet CLI user can start Spacedock with no workflow and receive the declared local identification result. The First Officer does not search the project, plugin, skill, or reference filesystem.

## Scope

- Repair only the Sonnet zero-discovery boot behavior.
- Keep the change in the smallest product instruction or runtime surface plus focused proof.
- Do not change n28, Pi, Codex, shared XFAIL policy, or unrelated live journeys.
- Do not add a permanent XFAIL.
- Use local Sonnet subscription authentication before required PR CI.

## Acceptance criteria

- AC-1: The exact Sonnet `TestLiveCommonZeroDiscovery` target passes normally and retains artifacts.
- AC-2: Boot uses only the declared local identification path. No `find`, recursive `grep`, or equivalent filesystem or reference sweep occurs.
- AC-3: A focused negative control fails for the exact broad-search command from artifact 9074747236.
- AC-4: Existing full, race, format, registry, and active-owner checks pass.
- AC-5: Required exact PR lanes pass before merge. Pi remains skipped.

## Baseline evidence

- Released user and workflow: Sonnet CLI zero-discovery boot.
- Observable harm: the First Officer broad-searches reference files instead of stopping as declared.
- Value authority: the zero-discovery live journey and Captain direction for truthful Sonnet evidence.
- Trigger: run 31415550991, job 93543929461, artifact 9074747236, command `find /tmp/spacedock-live-plugin-3439009114/skills/first-officer/references -iname "*claude*"`.

## Ideation requirements

- Name exact files and gross/net estimate before product edits.
- Identify the smallest behavior boundary and one falsifying control.
- Keep the normal local Sonnet failing baseline, repaired normal PASS, validation, PR, and merge flow.
