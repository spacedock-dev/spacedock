---
title: Polish the v1 gate command and documentation surface
status: backlog
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: top-level help omits gate prepare and the prose still describes prototype lifecycle details that the semantic cuts will remove."
started:
completed:
verdict:
score: "0.8"
worktree:
issue:
pr:
sprint: durable-decisions
id: f6cvn0s87ywbs158yy0b5q7k
gates:
    version: 1
    records:
        - id: gate:f6cvn0s87ywbs158yy0b5q7k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:f6cvn0s87ywbs158yy0b5q7k-backlog-1
              briefing:
                id: briefing:f6cvn0s87ywbs158yy0b5q7k:backlog:attempt-1:revision-1
                digest: sha256:2372b3af80ff6dcfbc792369d495bde4f3beb4b74238a334109cb1c72a2e46c9
                request-digest: sha256:557bd043331c2db7efd24534289d48d21a537bd6cae11efe2db628b7e3f738f6
                room-ref: ./review/backlog/briefing-1
---

Make the stable help, command reference, specification, and First Officer instructions describe the final minimal gate lifecycle after the semantic cuts land. This ticket bundles only small text, usage, and golden-output corrections.

## Problem

Top-level `spacedock --help` currently lists `gate record | validate | eligibility | consume` but omits `prepare`, even though `spacedock gate --help` exposes it. Several docs and skill lines still describe prototype application or advisory-round behavior. These discrepancies make the shortest supported journey harder to discover.

## Proposed approach

Update top-level help and existing help goldens. Then reconcile the command reference, gate concept/spec text, and First Officer lifecycle prose with the landed behavior from `1w6`, the actionable-stage guard, the application-schema cut, the gate-state simplification, the advisory-round cut, and the provider prove-or-cut decision. Keep one concise common-scenarios section; do not create a second lifecycle specification.

## Streamlined common scenarios

- Nonterminal approve: prepare, commit, present, record, commit, consume atomically into the next stage, commit, dispatch.
- Terminal approve: prepare, commit, present, record, commit, consume to `approved-awaiting-merge` without spending, then let the merge guard atomically spend and terminalize after delivery proof.
- Rework after delivery failure: merge guard supersedes the pending terminal application and routes to the declared feedback stage; a later gate uses a fresh attempt.
- Changed review before decision: withdraw the open prepared attempt, commit, and prepare a successor. Do not fabricate a hold.
- Hold or revise: retain the Resolution and follow the workflow route; do not manufacture unused application metadata.

## Out of scope

No new gate logic, schema fields, compatibility layer, live lane, or provider integration belongs here. If a doc correction requires behavior not already landed, route it to the owning semantic ticket.

## Acceptance criteria

**AC-1 - The command surface makes every retained v1 gate verb discoverable and no removed verb discoverable.**
Verified by: existing CLI help golden tests at the top-level and `gate` level.

**AC-2 - One common-scenarios section accurately describes the final executable lifecycle.**
Verified by: mapping each documented command sequence to an existing focused behavior test; any prose step with no executable counterpart fails review.

**AC-3 - The cleanup contains no behavior or storage-format change.**
Verified by: diff classification and unchanged non-golden behavior tests; a production semantic change sends the work back to its owning ticket.

## Test plan

Run focused help/golden tests and the existing gate lifecycle tests referenced by the docs. The full Go suite is sufficient; no new live exercise is justified for text-only reconciliation.
