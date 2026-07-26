---
title: Walk the assembled durable-decision journey before the pre-release
status: backlog
source: "Captain request on 2026-07-26 to exercise the sprint Definition of Done end to end and surface cross-member seams before release."
started:
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: ph0zv6azcrhcxmg57wwnxah7
gates:
    version: 1
    current:
        gate: gate:ph0zv6azcrhcxmg57wwnxah7:backlog
    records:
        - id: gate:ph0zv6azcrhcxmg57wwnxah7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ph0zv6azcrhcxmg57wwnxah7-backlog-1
              briefing:
                id: briefing:ph0zv6azcrhcxmg57wwnxah7:backlog:attempt-1:revision-1
                digest: sha256:0782c65c06c7ee9378226b3a7ef88d92939a54c05d916fe3690cc7d99804278f
                digest-domain: canonical-bytes
                request-digest: sha256:77aabae5f9e5af378e377bc1eaefccde931c8932e5e6023a661f5eae4a22e438
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ph0zv6azcrhcxmg57wwnxah7:backlog:1
                briefing: briefing:ph0zv6azcrhcxmg57wwnxah7:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T10:49:10.054082Z"
                decision: approve
                reason: The sprint's end value is an operable durable-decision journey; repeated late integration seams make one real release-candidate walking skeleton necessary before the pre-release.
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

The sprint has strong component proofs but repeatedly discovered integration seams only when a First Officer tried to operate the assembled lifecycle. Before the pre-release, drive one real walking skeleton through the actual command and skill surfaces and preserve correctness plus agent-friction evidence.

This is a one-time release-candidate acceptance journey, not a new standing harness or a second implementation. It runs after `s4`, `rq`, `v21`, and `withdraw-stale-open-gate-attempt` land. Findings route to their owning ticket; this ticket does not absorb product fixes.

## Captain journey

Describe intent in ordinary language: review a real prepared gate, approve it for a later Commander, re-scope a different still-open review, and then have a cold Commander apply the approved decision. Run the same presentation boundary once through default chat and once through the real `/subspace:r gate <room>` route.

## Required walking skeleton

1. A Shaping First Officer prepares a real folder-form gate room whose canonical Briefing selects one Artifact and one Reference by `git-root://` full-commit locators and raw SHA-256 revisions. The room begins with exactly the two authoritative metadata files and no copied payloads.
2. Default chat presents one bound review. The Captain approves for later execution. `gate record` retains the exact Resolution and leaves the ticket at its gated stage with `approved-awaiting-advance`; the Shaping FO does not call the Commander's `gate consume` verb.
3. A separate prepared attempt is legitimately re-scoped before decision. The First Officer withdraws it truthfully, cold-boots, prepares attempt N+1, and presents the replacement without deleting or rewriting the original room.
4. A cold Commander discovers the pending approval from boot, consumes it exactly once, advances, dispatches the declared successor, and durably commits the effect. A repeat consume fails without another transition or dispatch authorization.
5. The real Subspace room-only entry materializes the pinned Artifact and Reference from moved main/state checkouts, displays the canonical question, summaries, and exact bytes, retains provider evidence below the room, records the binding Result, and leaves no `association.json` or copied durable source payloads.
6. Tampered Briefing/source digests, missing Git objects, stale attempt identity, repeated withdrawal, and repeated consume each fail before presentation or mutation with an actionable error. Catchable provider failures retain diagnostics but remove ephemeral resolved payloads.

## Acceptance criteria

**AC-1 (VALUE)** A Captain and two role-separated First Officers can complete the natural-language journey without authoring JSON/YAML, reconstructing room coordinates, hand-editing `gates:`, manually committing package paths, or reverting an unintended stage transition.

**AC-2 (CORRECTNESS)** The retained state proves exact Briefing/Result/inventory identity, truthful withdrawal, record-versus-apply separation, one-use authorization consumption, one successor dispatch, frozen historical rooms, and provider-neutral recorder validation.

**AC-3 (CHANNEL PARITY)** Default chat and real Subspace presentation reach the same recorder-valid decision boundary; the provider override changes presentation only, never gate authority or state semantics.

**AC-4 (DURABILITY)** One `state commit` per lifecycle boundary captures the ticket plus new/changed room evidence, excludes seeded sibling dirt, and survives a cold boot from a clean checkout.

**AC-5 (OPERABILITY)** The report includes a timestamped step table with actor intent, exact command/skill entry, pre-state, output, post-state, retained artifact/digest, elapsed time, and friction. Every friction is classified as address now, deferred with promotion condition, polish, or needs decision and routed to one owner.

## Proof policy

Use exact release-candidate Spacedock and pinned Subspace commits. Run the smallest failure mutants inline with the journey; do not create a parallel fake provider or new generic test framework. The final report must distinguish command correctness from operator role error, and check the end value instead of counting commands.
