---
id: 2ae8r33r18g0w0g21559yc57
title: Dev template ships the rigor scar tissue and refit propagates it to commissioned workflows
status: done
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.6"
sprint: 0260-proportionality
group: template
started: 2026-07-20T06:40:51Z
gates:
    version: 1
    current:
        gate: gate:docs-dev:2ae:ideation
        attempt: gate-attempt:2ae-ideation-2
    records:
        - id: gate:docs-dev:2ae:ideation
          stage: ideation
          current-attempt: gate-attempt:2ae-ideation-2
          attempts:
            - id: gate-attempt:2ae-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:single-file:b567c1211ed3a2257a92f1725c2e93bc
                digest: sha256:4760be51f28b83d92b5671119ab26916113c4b470d44603377a2a88cc2800448
                room-ref: review/ideation/briefing-1/gate-summary.md
                note: Subspace advisory float (single-file, working-copy skill launcher); the artifact is the gate summary with the frozen entity snapshot appended
              resolution:
                type: Resolution
                id: resolution:actor-1784539060005018000
                briefing: briefing:single-file:b567c1211ed3a2257a92f1725c2e93bc
                by: person:reviewer
                at: 2026-07-20T09:17:40Z
                decision: approve
                reason: "on validation gate, present the refitted delta on the workflow readme for human review"
              application:
                action: advance
                target-stage: implementation
                state: superseded
              note: "The resolution reason is a binding captain instruction for the VALIDATION gate: its presentation must include the refit diff against the workflow README for human review. Carry into the Commander package. bw's Feedback Cycles format stays deferred (surfaced in the briefing, no annotation overriding it). Superseded by attempt 2 (captain-approved staff-review folds); the approval and the validation-gate instruction both stand."
            - id: gate-attempt:2ae-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:2ae-ideation-1
              state: closed
              briefing:
                id: briefing:2ae-ideation-2-chat
                digest: sha256:1f71f711733bec3fe6d6d6a243c818767938cdc78388dab2cca5056ef32f3132
                note: "chat presentation; ADVISORY digest — it hashes the working file at recording time (body folds applied, this attempt's own record excluded), which no single committed tree reproduces because an entity cannot self-bind its gates record. For drift checking, diff the entity BODY against the state commit that introduced this attempt; do not re-hash the current file. Digest refreshed once in the same fold round after the closure pass caught a residual cross-reference (Piece 6's placement parenthetical)."
              resolution:
                type: Resolution
                id: resolution:captain-chat-2ae-ideation-2
                briefing: briefing:2ae-ideation-2-chat
                by: person:captain
                at: 2026-07-20T10:20:31Z
                decision: approve
                reason: "Staff-review folds, captain-approved in chat: Piece 4 restored word-for-word to 02av's approved block (the compressed sentences returned) and moved from the validation to the implementation stage-def mirroring 02av's dev placement; the commission-skill pattern sentence names implementation-with-review-rounds stages; AC-2 reworded to a live refit-skill dry-run (a dispatched agent drives Phase 3b; validator-performed regeneration proves nothing about the skill); Pieces 1-3 re-anchor against landed sibling text at implementation."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "FO applied the folds directly under the captain's edit-directly grant; fable delta findings 4-6 and the codex finding-4 wave-4 sharpening. The attempt-1 validation-gate instruction (refit README delta presented for human review) carries forward unchanged."
              carried-finding:
                from: "roborev branch review of sibling member az, 2026-07-20"
                routed-by: agent:first-officer
                record: "Declined on az as out-of-surface and ROUTED HERE because this entity owns the files. Finding, verified: skills/commission/SKILL.md and skills/commission/references/templates/development.md still present `grep` as first-class proof in their AC suggestion lists and omit the falsifying-change clause, so a normal `commission` run generates acceptance criteria that CONTRADICT the policy az landed in docs/dev/README.md (evidence must name the concrete change that would make it fail; a one-off grep is structural evidence only and cannot satisfy a behavioral AC). This is the template-gap half of the original anti-tautology scope, which the 0260 re-lock moved to the template group — so it is this entity's by design, not new scope. Fold it into the Verified-by fix already in this entity's declared surface rather than treating it as an addition. IMPORTANT: the reviewer also proposed adding a cross-file consistency test to keep the templates aligned; that was DECLINED OUTRIGHT and must not be implemented — a new committed check needs explicit captain approval and its own entity, and a test asserting two instruction files agree in wording is a prose-to-prose consistency check, the banned shape. Verify against az's LANDED text, not this note."
              coverage-gap:
                raised-by: agent:first-officer
                at: 2026-07-20T14:40:00Z
                question: "Captain asked whether az's docs/dev/README.md changes consolidate into the dev workflow TEMPLATE later. FO checked the Pieces against az's landed edits and against the template's current text; recording the answer here because two gaps are real and this entity is the only propagation path."
                covered-by-re-anchor: "Piece 1 re-anchors on docs/dev/README.md's AC-template Verified-by line, which az sharpened with 'name the concrete change that would make it fail' — so it picks that up automatically. Piece 2 carries docs/dev/README.md:76 verbatim, which now includes the captain's prose-grep ruling AND the honesty-of-evidence bounding clause — so it picks those up too. Both work ONLY if implementation honours the wave-4 re-anchor rule and quotes the LANDED line rather than this entity's ideation-time snapshot."
                gap-1: "az's Edit B standalone bullet — 'Evidence must be able to fail: each AC's cited evidence names the concrete change that would flip it; an author who cannot name what would make the evidence fail has not shown it can fail, and the criterion does not count' — is carried by NO Piece. Verified receiving surface: development.md has ZERO occurrences of a can-fail rule; its 'External-proof acceptance criteria' bullet (~line 93) requires evidence from outside the task body but never asks the author to name the falsifying change. Piece 1 fixes only the AC-template stub. So a workflow commissioned today inherits the weaker rule. DECIDE DELIBERATELY: port it into the external-proof bullet, or record a decline with grounds — do not let it fall through by omission."
                gap-2: "az's Edit D — the detached audit ALSO fires on AC provenance (an AC whose expected value derives from the same package's production functions or constants), scoped to that AC's adversarial-edit check — is carried by NO Piece. Verified receiving surface: development.md DOES already carry a detached-adversarial-audit bullet (~line 94), so there is somewhere for it to land; the gap is not structural. DECIDE DELIBERATELY: port, or decline on the grounds that the provenance trigger is a dev-repo-specific sharpening a fresh workflow does not need yet."
                do-not-propagate: "NOT everything az landed belongs in the template, and porting indiscriminately would be its own error. The required-CI-lane rule and its path-to-lane mapping are a dev-lane realization tied to this repo's specific lanes, and the validation stage-def's routine-change exemption qualification is likewise docs/dev-specific. Those stay put. The generic disciplines (evidence must be able to fail; the prose-grep honesty boundary; arguably the provenance trigger) are the propagation candidates. Judge each on whether a NEWLY COMMISSIONED workflow with no CI and no lanes would be served by it."
worktree: .worktrees/spacedock-ensign-template-rigor-propagation
pr: pr-merge:542
verdict: passed
completed: 2026-07-21T04:43:00Z
archived: 2026-07-21T04:43:00Z
---
# Frozen production frontmatter replay fixture.
