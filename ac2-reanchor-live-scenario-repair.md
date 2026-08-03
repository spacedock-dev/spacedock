---
title: Repair the ac2 re-anchor live scenario so it can fail
status: backlog
source: "0260 preflight staff review, 2026-07-20 — re-filed from the archived ac2-reanchor-scenario-falsifiable, whose repair scope was dropped in a double archive merge."
id: 1azrdbz8bke5m0c3qbehye5c
sprint: live-test-truth
group: common-journey
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:1azrdbz8bke5m0c3qbehye5c:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:1azrdbz8bke5m0c3qbehye5c-backlog-1
              briefing:
                id: briefing:1azrdbz8bke5m0c3qbehye5c:backlog:attempt-1:revision-1
                digest: sha256:dfd9cd4fe892c9aff524770f758295873c5af5cf323561c40f537f834f03e4eb
                request-digest: sha256:a26381d34a975a709f05ca5d067131ac911ed5b052a9b0b7006a400016464107
                room-ref: ./ac2-reanchor-live-scenario-repair/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:1azrdbz8bke5m0c3qbehye5c:backlog:1
                briefing: briefing:1azrdbz8bke5m0c3qbehye5c:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T11:21:00.280158Z"
                decision: approve
                reason: Captain directed the First Officer to continue the next risk-first ideation wave.
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

`internal/livescenario/ac2_reanchor.go` stays green when the deliverable clause it exists to verify is deleted — a committed live check that cannot fail (verified pre-tag major in the archived entity). The archived owner was merged into ac2-design-proof-fixture, which was itself archived into the check-ordering task's lure catalog carrying only the reviewer-side fixture; the Go-scenario repair vanished from every sprint document while the file stayed on disk. Repair: drive the scenario both ways (correct and incorrect FO behavior) asserting divergent durable state, proven RED-first against the clause-removed deliverable — the archived entity's test plan, unchanged. Next train; not part of 0260 (captain decision, staff-review round: 0.26.0 ships with this known and recorded, not silently).
