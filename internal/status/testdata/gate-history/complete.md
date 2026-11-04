---
id: task
title: Gate history fixture
status: validation
gates:
  version: 1
  records:
    - id: gate:task:validation
      stage: validation
      attempts:
        - id: attempt:validation:1
          briefing:
            id: briefing:validation:1
            digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
            request-digest: sha256:2222222222222222222222222222222222222222222222222222222222222222
            room-ref: '@review/validation/briefing-1'
          provider-evidence:
            result-digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
          resolution:
            type: Resolution
            id: resolution:validation:1
            briefing: briefing:validation:1
            by: person:captain
            at: 2026-09-16T00:00:00.1200Z
            decision: approve
            reason: ''
            includes: []
          application:
            target-stage: done
            state: superseded
            legacy: {live: false, counts: [1, 2], note: '  old  '}
        - id: attempt:validation:2
          briefing:
            id: briefing:validation:2
            digest: sha256:2222222222222222222222222222222222222222222222222222222222222222
            room-ref: '@review/validation/briefing-2'
          resolution:
            type: Resolution
            id: resolution:validation:2
            briefing: briefing:validation:2
            by: agent:first-officer
            at: '2026-09-17T00:00:00Z'
            decision: approve
            reason: "  Résumé — 保留 <reason> & \"quote\"\nline\tend  "
            conn:
              quote: Proceed with the 3.
              source: Captain message
          application:
            target-stage: done
            state: superseded
        - id: attempt:validation:3
          briefing:
            id: briefing:validation:3
            digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
            room-ref: '@review/validation/briefing-3'
          resolution:
            type: Resolution
            id: resolution:validation:3
            briefing: briefing:validation:3
            by: person:captain
            at: '2026-09-18T00:00:00Z'
            decision: approve
            reason: Final approval
            includes: [annotation:first, annotation:last]
          application:
            target-stage: done
            state: pending
    - id: gate:task:ideation
      stage: ideation
      attempts:
        - id: attempt:ideation:1
          briefing:
            id: briefing:ideation:1
            digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
            request-digest: sha256:5555555555555555555555555555555555555555555555555555555555555555
            room-ref: '@review/ideation/briefing-1'
          withdrawal:
            by: agent:first-officer
            at: 2026-09-15T00:00:00.123400Z
            reason: '  旧資料  '
        - id: attempt:ideation:2
          briefing:
            id: briefing:ideation:2
            digest: sha256:5555555555555555555555555555555555555555555555555555555555555555
            room-ref: '@review/ideation/briefing-2'
---
## Body
