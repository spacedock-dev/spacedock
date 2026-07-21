---
title: Contractlint runtime-semantics retirement — codex and pi phrase checks become behavior tests
status: done
source: "Contractlint antipattern sweep, 2026-07-11: codex_multi_agent_v2_contract_test.go and Codex portions of runtime_binding_block_test.go assert runtime meaning from host-adapter prose."
score: 0.34
id: 8413fc05vpp8116k54x8br15
sprint: 0260-proportionality
group: contract-cleanups
started: 2026-07-20T05:04:09Z
gates:
    version: 1
    current:
        gate: gate:docs-dev:841:validation
        attempt: gate-attempt:841-validation-2
    records:
        - id: gate:docs-dev:841:validation
          stage: validation
          current-attempt: gate-attempt:841-validation-2
          attempts:
            - id: gate-attempt:841-validation-2
              sequence: 2
              state: closed
              briefing:
                id: briefing:841-validation-2
                note: "Validation cycle 2 stage report (PASSED), scoped to what changed since the cycle-1 REJECT. The report IS the briefing."
              resolution:
                type: Resolution
                id: resolution:fo-conn-841-validation-2
                briefing: briefing:841-validation-2
                by: agent:first-officer
                at: 2026-07-20T15:30:00Z
                decision: approve
                reason: "Approved by the FO under the captain's explicit conn grant (2026-07-20). Delegated approval, NOT a captain decision. Grounds: the cycle-1 material finding is closed by the right instrument. The validator re-derived the uncovered-token set with a DELIBERATELY DIFFERENT implementation — its own Go-string-literal lexer, its own delimiter rules, the token universe pulled mechanically from main's three retired files — and reproduced exactly the same NINE, with no tenth, including probing the argument names both enumerations discarded as noise. It independently confirmed both retractions (subagent/intercom are named-not-asserted) and both mechanism corrections as necessary rather than preference. It verified the production-caller basis of the declaration-not-wire accepts: ToolArgs() has zero production callers and all five pi constructors have zero non-test callers. Both HARD tolerances hold exactly — runtime-meaning literal-in-adapter-prose assertions 0 across 8 re-inventoried sites, committed test-function count 10 vs a 10 baseline. go test ./..., -race, and go vet all exit 0."
              application:
                action: advance
                target-stage: done
                state: pending
              note: "Closure history worth preserving: this member ran four review rounds and two validation cycles, and every round found something real. What closed it was NOT another round — it was the captain's ruling to replace the hand-audit with a mechanical enumeration, after three successive hand-audits each missed a different token. The script then corrected the implementer in BOTH directions (nine uncovered rather than six or seven; two prior 'no owner anywhere' claims retracted to 'named, not asserted'). The FO deliberately closed the hardening cycle with five findings deferred, two of them genuine false-green paths in code this entity introduced. The validator tested that closure empirically instead of accepting it and UPHELD it on stronger grounds than the FO had: deferral 2's delegate/message_dm swap does pass contractlint, but internal/piruntime/teams_test.go reds on it, so the regression cannot ship green anywhere. Required lanes: this diff touches only internal/contractlint/*_test.go — no skills/**/references/**, no dispatch or launch path, and none of the live lanes' own tests — so deterministic lanes suffice. Proven rather than assumed: go vet -tags live ./internal/... exits 0 on this branch, so live-tagged compilation is unaffected."
        - id: gate:docs-dev:841:ideation
          stage: ideation
          current-attempt: gate-attempt:841-ideation-1
          attempts:
            - id: gate-attempt:841-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:841-ideation-1-chat
                digest: sha256:7b46c619b5f83efe557122fb3a5e2e016f13c75500713be226851984bd1515c5
                note: chat presentation (per-check disposition + binding counts); digest is the entity content immediately before this record
              resolution:
                type: Resolution
                id: resolution:captain-chat-841-ideation-1
                briefing: briefing:841-ideation-1-chat
                by: person:captain
                at: 2026-07-20T05:40:00Z
                decision: approve
              application:
                action: advance
                target-stage: implementation
                state: consumed
worktree: .worktrees/spacedock-ensign-contractlint-codex-runtime-semantics-retirement
pr: pr-merge:539
verdict: passed
completed: 2026-07-20T15:33:33Z
archived: 2026-07-20T15:33:33Z
---
# Frozen production frontmatter replay fixture.
