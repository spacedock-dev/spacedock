---
id: 85z12f0ywkzy47akg9gwh6hm
title: "merge-guard arm phase has no keep-moving clause: armed reads as a stopping point"
status: done
source: "FO self-diagnosis, 2026-07-08 live session. After the captain approved three validation gates and said \"push it,\" the FO ran `spacedock merge guard <slug> --verdict passed` for each entity, which only ARMS the merge (sets mod-block=merge:pr-merge) — then stopped to read the pr-merge.md hook file instead of immediately constructing and presenting the PR draft in the same turn. Before finishing even one entity's draft, the FO got pulled into an unrelated task and the arm sat untouched. When the captain later asked \"what did you do when I said push it,\" the honest answer was: armed three entities, pushed nothing."
started: 2026-07-20T03:29:40Z
completed: 2026-07-20T17:26:28Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-merge-guard-arm-not-a-stopping-point
issue:
sprint: 0260-proportionality
group: verification
gates:
    version: 1
    current:
        gate: gate:docs-dev:85:validation
        attempt: gate-attempt:85-validation-2
    records:
        - id: gate:docs-dev:85:validation
          stage: validation
          current-attempt: gate-attempt:85-validation-2
          attempts:
            - id: gate-attempt:85-validation-2
              sequence: 2
              state: closed
              briefing:
                id: briefing:85-validation-2
                note: "Validation cycle 2 stage report (PASSED, zero material findings), read in full by the FO before resolving. The report IS the briefing. Cycle 1 REJECTED with three material findings; all three are resolved or dispositioned."
              resolution:
                type: Resolution
                id: resolution:fo-conn-85-validation-2
                briefing: briefing:85-validation-2
                by: agent:first-officer
                at: 2026-07-20T14:45:00Z
                decision: approve
                reason: "Approved by the FO under the captain's explicit conn grant (2026-07-20). Delegated approval, NOT a captain decision — recorded as agent:first-officer. Grounds: both shipped ACs verified against evidence the validator REPRODUCED. AC-1's blind A/B was re-run from a pre-registration written before any reader — 0/3 before-text vs 3/3 after-text surface-and-stop — and in doing so the validator found that the RECORDED probe had tested the parked 843-byte paragraph rather than the shipped clause (cmp separates them at char 435), so this is the first evidence that the text actually shipping moves the baseline. It replicates. AC-2 measured two independent ways agreeing exactly: net -24 bytes, ratchet green, full suite green uncached. The payload substitution survives its hardest test: the original AC-1 is byte-identical to a frozen ideation briefing committed nine hours BEFORE the substitution (an immutable independent source, not self-reference), and two blind readers given only the entity body each identified the shipped text, flagged the parked paragraph as unshippable, and answered NO to whether the entity delivered what it set out to — so the retreat reads as a retreat. Detached adversarial audit held on both misattribution and information-loss attacks. Zero material findings; three deferred risks and four polish items recorded with triggers and promote conditions."
              application:
                action: advance
                target-stage: done
                state: pending
              note: "MERGE PRECONDITION, unmet at approval: this diff touches skills/**/references/** (host-neutral FO contract), so offline plus every host lane is required. `offline` is RUN and GREEN — the cycle-1 red is fixed, and that red was caused by the parked clause's bytes. The three live lanes had NOT RUN at approval because the branch was unpushed and no PR existed; the validator stated that plainly rather than calling them green, and noted correctly that the captain's pi waiver covers a pi RED, not an UNRUN lane. FO note on its own conduct: the FO previously cited the cycle-2 probe's 3/3 result to the captain when recommending the payload substitution, without verifying which text that probe exercised. It had exercised the parked paragraph. The conclusion survived re-testing, but the evidence chain was broken and was caught only because the validation dispatch required re-running rather than accepting recorded numbers."
        - id: gate:docs-dev:85:ideation
          stage: ideation
          current-attempt: gate-attempt:85-ideation-2
          attempts:
            - id: gate-attempt:85-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:85-ideation-1
                digest: sha256:843de43172a64d695f0423dc81a357fd4b3af6c3c5653623e02480ffbe4983e7
              resolution:
                type: Resolution
                id: resolution:actor-1784520391265880000
                briefing: briefing:85-ideation-1
                by: person:reviewer
                at: 2026-07-20T04:06:31Z
                decision: approve
                reason: "based on our new principle, write a estimated change, so future stage can refer to it to judge deviation"
              application:
                action: advance
                target-stage: implementation
                state: superseded
              note: "Subspace advisory float, captain at the keyboard as person:reviewer; the resolution reason directs appending a declared expected surface to the body — applied post-approval as part of the approval's own terms, not drift. Superseded by attempt 2 (captain-approved staff-review vocabulary sweep); the approval itself stands."
            - id: gate-attempt:85-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:85-ideation-1
              state: closed
              briefing:
                id: briefing:85-ideation-2-chat
                digest: sha256:cff3a819597e625058429c606ac788046316d0e0e31f72941f210e493b2394ba
                note: "chat presentation; ADVISORY digest — it hashes the working file at recording time (sweep applied, this attempt's own record excluded), which no single committed tree reproduces because an entity cannot self-bind its gates record. For drift checking, diff the entity BODY against the state commit that introduced this attempt; do not re-hash the current file."
              resolution:
                type: Resolution
                id: resolution:captain-chat-85-ideation-2
                briefing: briefing:85-ideation-2-chat
                by: person:captain
                at: 2026-07-20T10:23:28Z
                decision: approve
                reason: "Staff-review sweep, captain-approved in chat: the banned coined vocabulary removed from the group field and four body spots (plain 'captain judgment' / 'check ordering' wording; group regrouped to verification). No design change."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "FO applied the sweep directly under the captain's edit-directly grant; codex finding 5 and fable delta finding 19."
pr: pr-merge:537
archived: 2026-07-20T17:26:28Z
---
# Frozen production frontmatter replay fixture.
