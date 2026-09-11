---
title: "A taskless spacedock pi launch never boots the FO: no first request, silent startup"
status: backlog
source: "Captain observation 2026-09-11 in the wrapped dev-override TUI session: silence at startup was the evidence the FO never booted; manual /first-officer produced the greeting and workflow stats. Diagnosis trail: post-s98 the frontdoor no longer appends any launch task (AC-3 removed the inert $spacedock:first-officer sentence), so a taskless launch sends no first model request; the extension bootstrap (session_start-armed, context-hook injected) is request-time-only and once-per-turn (agent_end disarms), so it never fires. Pre-s98 behavior greeted at startup because the always-appended prompt forced a first request and the old ungated extension injected."
id: n315frdw60950kjde4cx017q
gates:
    version: 1
    records:
        - id: gate:n315frdw60950kjde4cx017q:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:n315frdw60950kjde4cx017q-backlog-1
              briefing:
                id: briefing:n315frdw60950kjde4cx017q:backlog:attempt-1:revision-1
                digest: sha256:0bde44c6681bd0f95ed27aa526533c8a4ad99776979c01b4888ed8dae11d7d39
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:n315frdw60950kjde4cx017q:backlog:1
                briefing: briefing:n315frdw60950kjde4cx017q:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-11T15:48:17.855233Z"
                decision: approve
                reason: 'Captain approved in chat 2026-09-11 at the seed presentation: fast-track, ideation steered toward mechanism (a) — extension delivers the bootstrap as a real first-turn message via the documented pi.sendUserMessage()/expandPromptTemplates upgrade path; frontdoor stays argv-and-env-only'
              application:
                target-stage: ideation
                state: pending
---
Problem: a fresh `spacedock pi` launch with no operator task sits silent — the FO never presents the boot summary or workflow stats until the operator types something (and a weak model may never self-boot even then; glm-5.3-flash misreported its own context twice today). The pre-s98 launch greeted at startup; s98's correct argv fix removed the trigger without replacing it.

Deliverable: a taskless launch boots and greets (binary gate, boot identify, session summary, stop for input) with no operator input. Two candidate mechanisms to weigh at ideation: (a) the extension delivers the bootstrap as a REAL first-turn message via pi.sendUserMessage()/expandPromptTemplates — the upgrade path the s98 body already documents — which also makes the bootstrap persist instead of being once-per-turn and request-time-only (the observability trap hit today: multi-turn probes made a working injection look absent); (b) the frontdoor appends a neutral default launch task when no operator task is given (no $spacedock syntax — AC-3 intact), restoring the startup request; the bootstrap stays ephemeral.

Acceptance criteria and test plan to be fleshed out at ideation. Value AC direction: a fresh taskless launch (wrapped and unwrapped, installed and dev-override) presents the interactive greet without operator input, measured against today's silent-startup baseline; plain `pi` sessions still receive nothing (single-owner gate preserved).
