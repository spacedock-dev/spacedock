---
title: Live harness 60 s quiet budget kills runs while a background worker has a silent turn
status: ideation
source: captain (2026-11-03), from the repair-sonnet-live-flakes ideation probes
score: 0.6
sprint:
group: live-lane
sprint-readiness: ready
id: w6kcbj4k8ev9t8n2hpw52dds
started: 2026-11-04T07:38:35Z
---

`drainToExit` in the live harness kills a run when no stream line arrives for 60 s (`quietBudgetDefault`). In the keep-moving probes of `repair-sonnet-live-flakes` (Sonnet 5.5, effort max, head `64274cef5`), 7 of 8 unchanged-byte runs were killed this way. Each kill came after the FO ended its turn while the background `questioned` worker was in a long thinking turn. In run b1 (budget raised to 180 s), silent gaps of 78 to 128 s occurred five times. These kills are infra failures, not grades, and the XFAIL binding does not absorb them, so they fail the lane.

The kill can depend on load: three other workers shared the machine, and earlier the same journey passed 3 of 3 locally. CI Sonnet 5.5 ran the journey in about 62 s.

Evidence: the `### Finding: quiet budget kills runs while a background worker is in a silent turn` section of `repair-sonnet-live-flakes.md`.

Acceptance criteria (ideation fleshes): a live run is not killed while a dispatched background worker is alive and its process is making progress, and a run whose processes are all stalled still fails fast; proven by a deterministic harness test with a fake stream that goes silent for longer than the budget while a child is alive (no kill) and while no child is alive (kill).
