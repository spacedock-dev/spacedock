---
title: Simplify the unreleased v1 gate-state schema
status: backlog
source: "Durable-decisions sprint implementation-shape audit, 2026-07-24."
score: "0.7"
id: jccbpvjv5bg1jn0jbmj2yf8s
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:jccbpvjv5bg1jn0jbmj2yf8s:backlog
    records:
        - id: gate:jccbpvjv5bg1jn0jbmj2yf8s:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:jccbpvjv5bg1jn0jbmj2yf8s-backlog-1
              briefing:
                id: briefing:jccbpvjv5bg1jn0jbmj2yf8s:backlog:attempt-1:revision-1
                digest: sha256:4d61c845361a3aaba15ec68a09a5090b03c963c532d61227de695e4473ae32c3
                digest-domain: canonical-bytes
                request-digest: sha256:7b480163e7d4760504a04d5d91b979e8b942b0de4ab99b690f1619145e3c4db3
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:jccbpvjv5bg1jn0jbmj2yf8s:backlog:1
                briefing: briefing:jccbpvjv5bg1jn0jbmj2yf8s:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:54:16.651475Z"
                decision: approve
                reason: Captain approved dispatching this durable-decisions ideation lane in parallel with wj, hq, and nth.
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

The unreleased v1 gate-state implementation still carries prototype compatibility and a mutable current-gate pointer that duplicates derivable state and has already projected a stale approval.

Ideation should define the smallest clean-v1 schema:

- Remove `raw-file-pin` support and its compatibility fixtures; canonical bytes are the sole shipped binding.
- Determine whether `digest-domain` is redundant once the v1 schema fixes canonical digest semantics.
- Exercise multi-stage re-entry, multiple historical attempts, changed-Briefing supersession, and same-stage replay to determine whether `gates.current.gate` can be derived from current stage plus immutable records.
- If derivation is sufficient, remove the stored pointer and all pointer-repair/rebind behavior. If one counterexample requires stored selection, record that minimal counterexample and retain only the smallest non-stale selector.

The test must reproduce the observed failure class: an older approved attempt must not make a newly rejected candidate appear approved-awaiting-merge. No prototype migration or compatibility path is required.
