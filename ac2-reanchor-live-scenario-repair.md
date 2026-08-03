---
title: Repair the ac2 re-anchor live scenario so it can fail
status: backlog
source: "0260 preflight staff review, 2026-07-20 — re-filed from the archived ac2-reanchor-scenario-falsifiable, whose repair scope was dropped in a double archive merge."
id: 1azrdbz8bke5m0c3qbehye5c
sprint: live-test-truth
group: common-journey
sprint-readiness: ready
---

`internal/livescenario/ac2_reanchor.go` stays green when the deliverable clause it exists to verify is deleted — a committed live check that cannot fail (verified pre-tag major in the archived entity). The archived owner was merged into ac2-design-proof-fixture, which was itself archived into the check-ordering task's lure catalog carrying only the reviewer-side fixture; the Go-scenario repair vanished from every sprint document while the file stayed on disk. Repair: drive the scenario both ways (correct and incorrect FO behavior) asserting divergent durable state, proven RED-first against the clause-removed deliverable — the archived entity's test plan, unchanged. Next train; not part of 0260 (captain decision, staff-review round: 0.26.0 ships with this known and recorded, not silently).
