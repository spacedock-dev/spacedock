---
id: z7sfm93ccddg7x2tycp1smwy
title: Prefer the cheapest check that can fail — replaces "code gate over prose rule", with new-check consent, fan-out surfacing, and the no-minting authoring rule
status: done
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.75"
sprint: 0260-proportionality
group: verification
started: 2026-07-20T03:29:38Z
gates:
    version: 1
    current:
        gate: gate:docs-dev:z7:validation
        attempt: gate-attempt:z7-validation-1
    records:
        - id: gate:docs-dev:z7:validation
          stage: validation
          current-attempt: gate-attempt:z7-validation-1
          attempts:
            - id: gate-attempt:z7-validation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:z7-validation-1
                note: "Validation stage report (PASSED, zero material findings). The report IS the briefing; its lure-scenario matrix is this entity's only behavioural evidence, by captain ruling that the catalog's home is a report and not a committed suite."
              resolution:
                type: Resolution
                id: resolution:fo-conn-z7-validation-1
                briefing: briefing:z7-validation-1
                by: agent:first-officer
                at: 2026-07-20T15:45:00Z
                decision: approve
                reason: "Approved by the FO under the captain's explicit conn grant (2026-07-20). Delegated approval, NOT a captain decision. Grounds: the behavioural claims rest on 30 LIVE DRIVES rather than assertion — six lure scenarios x branch/main x Claude/codex plus a commissioned-check control. Three of six discriminate. Scenario 6 (fan-out authoring, whose fixture is this sprint's own 110-agent incident) discriminates under BOTH runtimes: branch arms declare worker count and tolerance before launch and refuse the harness's per-finding two-verifier guidance, while both main arms reproduce the incident shape. Scenario 3 (minting) discriminates under both. The entity funds itself at -81 bytes with the ratchet green and sibling headroom growing 403 -> 484, independently re-measured at both refs and matching the implementer exactly. Trim discipline verified: six largest cuts traced to named owners with zero orphans, one now loading EARLIER than the text it left. The detached adversarial audit found no route by which the ordering justifies building machinery, and confirmed the second-verifier rule does not license skipping the mandated audit — four skip readings tried, none survived."
              application:
                action: advance
                target-stage: done
                state: pending
              note: "HONEST LIMITS, recorded because they qualify the approval rather than decorate it. (1) AC-1's negative control discriminates under codex only — codex/main dispatches the PTY harness build outright while codex/branch holds, but Claude/main ALSO refuses by a different route, so under Claude the clause produces the specified stop FORM without changing the OUTCOME. (2) Scenarios 2, 4 and 5 pass on all four arms and test rules pre-existing to this diff; they prove non-regression and are not evidence FOR this clause. (3) The audit's sharpest result is against the entity: inverting the ordering to put 'build a new standing lint' first, and deleting the consent stop, fan-out checkpoint and second-verifier rule outright, leaves go test ./... fully green — the change is prose with zero mechanical coverage, by design and captain ruling, which is why the drives are the only evidence it has. (4) One contaminated cell was detected, discarded and re-run: a run with filesystem reads enabled recognised the planted fixture and quoted the catalog back; a marker scan found no second instance and a residual symmetric confound (Claude Code injecting recent commit subjects) is disclosed. Five deferred risks recorded with triggers and promote conditions. MERGE PRECONDITION: the diff touches two host-neutral contract files plus the Claude adapter, so claude-live (both legs), codex-live and pi-live are all REQUIRED; at approval time NONE had run, because the live workflow is PR-gated and the branch was pre-PR. The captain's waiver covers a pi lane RED, pi only — it covers neither an UNRUN lane nor claude-live nor codex-live."
        - id: gate:docs-dev:z7:ideation
          stage: ideation
          current-attempt: gate-attempt:z7-ideation-5
          attempts:
            - id: gate-attempt:z7-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:z7-ideation-1
                digest: sha256:f4fefb0b7be0d57c933fe27d71aa90235e6b01faef2301f916db22a8f0065b47
              resolution:
                type: Resolution
                id: resolution:actor-1784520766469687000
                briefing: briefing:z7-ideation-1
                by: person:reviewer
                at: 2026-07-20T04:12:46Z
                decision: revise
                reason: "Annotation on the clause headline: do not invent new terminology — the named ladder/rungs/climb concept is itself minted vocabulary; express the ordering rule in plain language."
              application:
                action: feedback
                target-stage: ideation
                state: consumed
              note: "Subspace advisory float; captain annotation included (annotation:captain-1784520762311801000, TextQuoteSelector on the clause opening). In-stage revision routed to the live worker; next attempt opens at re-presentation."
            - id: gate-attempt:z7-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:z7-ideation-1
              state: closed
              briefing:
                id: briefing:z7-ideation-2
                digest: sha256:baf726be8962c9e4238117c1cb80ae66577f12277aeb7da997277a7ad8b455d4
              resolution:
                type: Resolution
                id: resolution:actor-1784522359233707000
                briefing: briefing:z7-ideation-2
                by: person:reviewer
                at: 2026-07-20T04:39:19Z
                decision: revise
                reason: "Five annotations: drop the jargon term in the falsifiable-exercise step; format the ordering as bullets; the last-resort build requires explicit approval and usually its own entity, and the wording may still be too dev-specific for the shared core; justify or restyle the inline deferred-file reference; carry the prose-grep nuance — allowed as one-off validation evidence, banned as a committed test."
              application:
                action: feedback
                target-stage: ideation
                state: consumed
            - id: gate-attempt:z7-ideation-3
              sequence: 3
              previous-attempt: gate-attempt:z7-ideation-2
              state: closed
              briefing:
                id: briefing:z7-ideation-3
                digest: sha256:da21cb40ecfe9aa1e1c32bbdcd70d49217d4af696251914e50bbfd2e8fa5fcf3
              resolution:
                type: Resolution
                id: resolution:actor-1784523600509859000
                briefing: briefing:z7-ideation-3
                by: person:reviewer
                at: 2026-07-20T05:00:00Z
                decision: revise
                reason: "Annotation on the ACs: a few scenarios need real testing — fixtures that lure the expensive action and observe it — but NOT as a committed test suite; what are the options; and coverage must include gpt-5.6-sol. Budget re-baseline question rolls forward undecided."
              application:
                action: feedback
                target-stage: ideation
                state: consumed
            - id: gate-attempt:z7-ideation-4
              sequence: 4
              previous-attempt: gate-attempt:z7-ideation-3
              state: closed
              briefing:
                id: briefing:z7-ideation-4-chat
                digest: sha256:e30d6bddfd4bdb525c40bf852f498caa6876346bc2486deb1c2fbb2ff462798d
                note: chat presentation; digest is the entity content immediately before this record was written
              resolution:
                type: Resolution
                id: resolution:captain-chat-z7-ideation-4
                briefing: briefing:z7-ideation-4-chat
                by: person:captain
                at: 2026-07-20T05:30:00Z
                decision: approve
                reason: "Approved in chat: the cycle-5 design including the lure-scenario catalog and option recommendation; boot-resident budget re-baselined to +1055 — the authoring rule stays boot-resident."
              application:
                action: advance
                target-stage: implementation
                state: superseded
              note: "Superseded by attempt 5 (captain-approved staff-review folds); the approval itself stands."
            - id: gate-attempt:z7-ideation-5
              sequence: 5
              previous-attempt: gate-attempt:z7-ideation-4
              state: closed
              briefing:
                id: briefing:z7-ideation-5-chat
                digest: sha256:ade2baa0213f223b2a6953ba436a2441274b78ce9e7130cd6f3be41dde0e0088
                note: "chat presentation; ADVISORY digest — it hashes the working file at recording time (body folds applied, this attempt's own record excluded), which no single committed tree reproduces because an entity cannot self-bind its gates record. For drift checking, diff the entity BODY against the state commit that introduced this attempt; do not re-hash the current file."
              resolution:
                type: Resolution
                id: resolution:captain-chat-z7-ideation-5
                briefing: briefing:z7-ideation-5-chat
                by: person:captain
                at: 2026-07-20T10:22:38Z
                decision: approve
                reason: "Staff-review folds, captain-approved in chat: the fan-out checkpoint also binds the authoring moment (a scripted fan-out declares expected agent count, tolerance, and economic reasonableness before launch), with the Claude runtime pre-Workflow declaration line as the delivery-at-the-trigger binding (third contract file — reconfirm trigger fired and reconfirmed); the docs/dev/README.md:74 edit widens to the whole sentence including the binary-or-test-only satisfier; lure scenario six (fan-out-authoring lure, live fixture: the 110-agents-queued incident) joins the catalog."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "FO applied the folds directly under the captain's edit-directly grant; fable delta finding 7, the z7-failure analysis the captain requested, and the captain's pre-workflow directive."
              scope-amendment:
                decision: "each member pays its own way"
                by: person:captain
                at: 2026-07-20T13:05:00Z
                record: "Captain ruling, sprint-wide, recorded BEFORE this entity's implementation dispatch because it materially amends the approved surface. A committed ratchet (TestFOFunctionPromptSurfaceShrinks, internal/contractlint/fo_function_reference_invariant_test.go) caps the total bytes of 13 first-officer contract files at 122634; the measured total on the rebased main is 122231, leaving 402 bytes of headroom. This entity's approved surface adds roughly +3000 bytes to three of those measured files, so it cannot land as approved. The ruling: every ratcheted member funds its own additions with an offsetting trim of at least equal size, taken from the files it already touches wherever possible. Alternatives considered and rejected: raising the baseline (editing a check so the change passes is the anti-pattern this sprint removes), and a shared offsetting trim of the dead legacy skill (creates an unmodelled multi-member seam on one file). The approved DESIGN is unchanged — the amendment adds a funding obligation, it does not alter the clause content the captain approved. If the entity cannot fund itself without weakening the contract, that is a recorded finding and a captain decision, not a licence to overrun. Note the sprint index's stated mitigation (prefer lazy-loaded references over boot-resident lines) does NOT satisfy this check: the measured set includes the deferred cores fo-dispatch-core.md and claude-fo-dispatch.md, not only boot-resident files."
worktree: .worktrees/spacedock-ensign-falsifiability-ladder
pr: pr-merge:540
verdict: passed
completed: 2026-07-20T17:00:53Z
archived: 2026-07-20T17:00:53Z
---
# Frozen production frontmatter replay fixture.
