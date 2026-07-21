---
id: bwr6j6edkmfx5sbz73cr2952
title: Feedback-cycle record convention and design-reset decision — declared estimate, per-round actuals, recorded reframe
status: done
source: "captain (2026-06-04) — forked from xa (feedback-guarantee-binary-gate) per the roadmap-the-decision + separate-build-task call. xa's ideation determined Candidate 1 (3-cycle escalation) is mechanizable via a dedicated cycle-record command (a spike disproved a --set status guard) and Candidate 2 (budget-probe) is not. This task SHIPS the Candidate-1 guard; xa closed as a roadmap decision."
score: "0.30"
started: 2026-07-20T03:29:33Z
completed: 2026-07-21T04:41:32Z
verdict: passed
worktree: .worktrees/spacedock-ensign-feedback-cycle-record-command
issue:
sprint: 0260-proportionality
group: reframe
gates:
    version: 1
    current:
        gate: gate:docs-dev:bw:ideation
        attempt: gate-attempt:bw-ideation-3
    records:
        - id: gate:docs-dev:bw:ideation
          stage: ideation
          current-attempt: gate-attempt:bw-ideation-3
          attempts:
            - id: gate-attempt:bw-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:bw-ideation-1
                digest: sha256:51b38d47a1b15d0ea7b34d4908d20d257ae3e39f5bc587839075a78466d10167
              resolution:
                type: Resolution
                id: resolution:actor-1784521963753201000
                briefing: briefing:bw-ideation-1
                by: person:reviewer
                at: 2026-07-20T04:32:43Z
                decision: revise
                reason: "Four annotations: (1) briefing packaging should use separate artifacts (FO-side, 3k experiment); (2) is the new record command needed at all — apply the cheapest-check ordering; (3) add a final landing-spot review AC (core vs dev-specific) and propose the dev README change; (4) roborev-shaped in-stage AC coverage."
              application:
                action: feedback
                target-stage: ideation
                state: consumed
              note: "Subspace advisory float; four captain annotations included by id in the resolution. Annotation 1 is FO-owned (briefing packaging), 2-4 routed to the worker; next attempt opens at re-presentation."
            - id: gate-attempt:bw-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:bw-ideation-1
              state: closed
              briefing:
                id: briefing:bw-ideation-2
                digest: sha256:520f84901aae2b1e9fc0c78eaaf974cbec3a8b8cbb82db3d8eb7fecb9a778a6d
              resolution:
                type: Resolution
                id: resolution:actor-1784523098406183000
                briefing: briefing:bw-ideation-2
                by: person:reviewer
                at: 2026-07-20T04:51:38Z
                decision: approve
              application:
                action: advance
                target-stage: implementation
                state: superseded
              note: "Subspace advisory float, no annotations; third presentation (prose-only convention) approved. Superseded by attempt 3 (captain-approved staff-review folds); the approval itself stands."
            - id: gate-attempt:bw-ideation-3
              sequence: 3
              previous-attempt: gate-attempt:bw-ideation-2
              state: closed
              briefing:
                id: briefing:bw-ideation-3-chat
                digest: sha256:837779a0b96ebddc7e695106109a6026034c28337686a73bba3d8d18f2ff8c6f
                note: "chat presentation; ADVISORY digest — it hashes the working file at recording time (body folds applied, this attempt's own record excluded), which no single committed tree reproduces because an entity cannot self-bind its gates record. For drift checking, diff the entity BODY against the state commit that introduced this attempt; do not re-hash the current file."
              resolution:
                type: Resolution
                id: resolution:captain-chat-bw-ideation-3
                briefing: briefing:bw-ideation-3-chat
                by: person:captain
                at: 2026-07-20T10:18:41Z
                decision: approve
                reason: "Staff-review folds, captain-approved in chat with the concrete strengthened wording shown: past-tolerance deviation requires a recorded reconfirm/re-scope/park/escalate decision before further repair dispatch (no automatic re-dispatch); tolerance default named as 2x unless the entity declares otherwise; the testlint arm struck from AC-1 (one-off script, no committed check); the skill frontmatter description line added to the declared surface; title corrected to the prose-only convention the body ships."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "FO applied the folds directly under the captain's edit-directly grant; codex staff review finding 1 and fable delta findings 8-10."
mod-block:
pr: pr-merge:541
archived: 2026-07-21T04:41:32Z
---
# Frozen production frontmatter replay fixture.
