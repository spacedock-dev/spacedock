---
title: Headless FO ends its turn while a background worker is open, and the host stops the worker
status: backlog
source: captain (2026-11-03), from the repair-sonnet-live-flakes ideation probe b2
score: 0.7
sprint:
group: live-lane
sprint-readiness: ready
id: 5779fk6mzntvv5j3x6hw8rdn
---

In a headless (`-p`) drive with background workers, the FO can end its turn while a worker is still running. In probe b2 of `repair-sonnet-live-flakes` (Sonnet 5.5, head `64274cef5`, keep-moving journey), the host stopped the background `questioned` ideation worker 600 s after the FO's last turn. The worker had not committed its report, so the journey graded `keep-moving-violation`.

The contract already covers one case: a single-entity `-p` run uses bare (blocking) dispatch, because "the blocking call prevents premature session termination". Multi-entity headless drives with background workers have no equal rule.

This task owns AC-1 and AC-3 of `repair-sonnet-live-flakes` (captain ruling 2026-11-03): an unbound live keep-moving PASS on Sonnet in 3 of 3 runs, and the retirement of the `060xp69y61yhrww23g3wvwqy` keep-moving binding and its registry row. Ownership of that binding moves here when `repair-sonnet-live-flakes` closes.

Evidence: the `### Finding: background worker stopped before it reported` section of `repair-sonnet-live-flakes.md`. The 600 s figure comes from one observation and is not confirmed in host documentation.

Acceptance criteria (ideation fleshes): a headless multi-entity drive does not end the FO turn while a dispatched worker has not reported; proven by a live keep-moving run where a slow background worker reports after the FO's last dispatch, plus the retirement evidence above. Options to weigh: a headless wait rule for the FO, and a stage-size rule for ensigns in ideation on fixture entities.
