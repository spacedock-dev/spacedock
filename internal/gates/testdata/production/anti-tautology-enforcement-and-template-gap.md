---
id: azh879wdzm72ysxg16hbg39q
title: Falsifiable-test estate — contract evidence rule and the automated-gate decision, resolved once
status: done
source: "Two related findings from this session's audit work, not yet designed: (1) docs/dev's own Proof policy ('no prose-grep over instruction files', the detached adversarial audit, the AC template's 'Verified by: ... something outside this task body ... that can fail' clause) only catches tautology reactively, at high-stakes-surface merge time via human/reviewer judgment -- there is no standing, automatic check against the mirror-assertion / no-op-assertion patterns that four real tests in this repo turned out to have. (2) The commission-skill templates that scaffold NEW workflows do not carry equivalent discipline for GO TEST CODE tautology (as opposed to instruction-file/prose tautology, which is entity ey / proof-policy-shipped-scaffolding's separate, pre-existing scope -- see OVERLAP NOTE below): skills/commission/references/templates/development.md's AC template stub and skills/commission/SKILL.md's base AC template both lack docs/dev's own 'outside this task body'/'that can fail' clauses. Reference: github.com/kenn-io/middleman skills/testing-without-tautologies/SKILL.md. A design workflow + an independent fable review ran this session and both converged (mechanism choice, scope, sequencing all hold per fable's independent check) on: NO automated hard AST gate for mirror-assertions (an internal/testlint check for the OTHER pattern, assertion-free tests, is its own sibling entity: testlint-assertion-free-gate) -- instead extend the existing detached-adversarial-audit trigger to fire on AC PROVENANCE ('any AC whose expected value is derived from the same package's production functions or constants'), not the originally-drafted broader 'equality/byte-identity check' wording fable found would over-fire on nearly every unit test in this repo. Concrete diffs exist for docs/dev/README.md's Proof policy + pr-merge gate rule, development.md, and SKILL.md. CORRECTED accounting (fable caught this): internal/status/boot_probe_parity_test.go is NOT a stampID mirror-assertion instance -- it mirrors a different production CONSTANT (teamStateNeutralHint), not a function call; the confirmed mirror-assertion-via-shared-function count is 2 (native_new_test.go, zz_independent_parity_test.go), already covered under the separate tautological-test-fixes entity, don't double-count here. OVERLAP NOTE, UNRESOLVED -- flagging for the captain, not silently resolving: entity ey (proof-policy-shipped-scaffolding, filed 2026-06-04, pre-existing) targets the SAME file (skills/commission/references/templates/development.md) for a related-but-distinct concern -- porting the INSTRUCTION-FILE/prose tautology test (not code-test tautology) to shipped scaffolding, plus first-officer-shared-core.md and ensign-shared-core.md, with a heavier behavioral AC (a live scenario proving a validator REJECTS a presence-only proof). This entity's development.md diff and ey's development.md target could collide if ideated independently without coordination. Captain has not yet said how to reconcile (fold together, sequence, or keep fully separate with a coordination note in each) -- do not dispatch this entity's ideation until that's decided."
started: 2026-07-20T04:53:02Z
completed: 2026-07-20T15:52:09Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-anti-tautology-enforcement-and-template-gap
issue:
sprint: 0260-proportionality
group: test-cleanups
gates:
    version: 1
    current:
        gate: gate:docs-dev:az:validation
        attempt: gate-attempt:az-validation-1
    records:
        - id: gate:docs-dev:az:validation
          stage: validation
          current-attempt: gate-attempt:az-validation-1
          attempts:
            - id: gate-attempt:az-validation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:az-validation-1
                note: "Validation stage report (PASSED, zero material findings), read in full by the FO before resolving. The report IS the briefing."
              resolution:
                type: Resolution
                id: resolution:fo-conn-az-validation-1
                briefing: briefing:az-validation-1
                by: agent:first-officer
                at: 2026-07-20T14:25:00Z
                decision: approve
                reason: "Approved by the FO under the captain's explicit conn grant (2026-07-20). Delegated approval, NOT a captain decision — recorded as agent:first-officer so the trail does not overstate its authority. Grounds: all four ACs verified against evidence the validator REPRODUCED rather than cited — the mirror/control mutation exercise rebuilt from scratch and extended to a 3-mutation matrix, AC-4's git state recomputed against a merge base that had moved again, the Edit C span diffed against an approved text in a different file (two values that can diverge). The captain's Edit C ruling is byte-identical and the shipped boundary is the honesty framing the captain corrected to, not the FO's earlier over-restrictive formulation: legitimate existence/absence evidence preserved, committed greps still banned, and no relabelling escape because the clause keys on what the evidence can establish rather than on how the claim is labelled. Edit D reads as a conditioned widening (provenance fires wherever such an AC appears, covering that AC alone), inside the captain's approved scope, not an always-on audit. Zero material findings; four deferred risks and four polish items recorded with triggers and promote conditions."
              application:
                action: advance
                target-stage: done
                state: pending
              note: "MERGE PRECONDITION, not a deliverable defect: this diff touches skills/ensign/references/ensign-shared-core.md, which skills/ensign/SKILL.md loads unconditionally for EVERY host, making it the host-neutral ensign contract rather than one adapter. Under the dev README's required-lane rule every host lane is REQUIRED green before merge — claude-live, codex-live, pi-live — with the captain's pi-red waiver covering pi ONLY. At approval time the lanes were UNRUN, not green. The validator explicitly refused to discharge them by analogy with sibling member ht, which merged on deterministic lanes alone after PROVING (go vet -tags live) that its test-only diff could not affect live-tagged compilation; that proof does not transfer to a change in the shipped ensign contract the lanes actually load. The FO accepts that correction."
        - id: gate:docs-dev:az:ideation
          stage: ideation
          current-attempt: gate-attempt:az-ideation-1
          attempts:
            - id: gate-attempt:az-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:docs-dev:az:ideation:attempt-1:revision-4
                digest: sha256:610fcfab5250d0d23eb7ed01f10eb702657b05b7be0692b194819715dff5bdc4
                room-ref: "./review/ideation/briefing-1"
              resolution:
                type: Resolution
                id: resolution:captain-chat-az-ideation-1
                briefing: briefing:docs-dev:az:ideation:attempt-1:revision-4
                by: person:captain
                at: 2026-07-20T05:40:00Z
                decision: approve
                reason: "Approved in chat after the specifics review (landing files + line counts). Edit D (the audit-trigger widening) is NOT covered by this approval — it awaits its own explicit yes/no per the new-enforcement consent rule."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "Captain hold via float 2026-07-20 (resolution:actor-1784524247673759000): the FO gate summary was session-jargon dense and unreadable without full context. Attempt open; plain-language rewrite via the comm-officer before re-presentation."
              edit-d-resolution:
                decision: approved
                by: person:captain
                at: 2026-07-20T12:21:26Z
                record: "Edit D (the existing detached audit ALSO fires whenever a test's expected answer comes from the same code being tested) — EXPLICIT captain yes, given to the Commander in chat when the contradiction was surfaced. This satisfies the attempt resolution's standing condition that Edit D 'awaits its own explicit yes/no per the new-enforcement consent rule'. Scope of the yes: a trigger widening of the EXISTING detached audit, contract prose only — no new tool, test, gate, lint, or CI lane, so AC-4's zero-new-enforcement criterion still holds. az ships Edits A-D across the two named instruction files. Supersedes the earlier reading that derived consent from the blanket 'agree with all recommendations' staff-review ruling plus a flag-if-not; a consent-gated edit needs a direct answer, which this record now carries."
pr: pr-merge:536
archived: 2026-07-20T15:52:09Z
---
# Frozen production frontmatter replay fixture.
