---
title: Premise spike — can a Haiku FO hold the mechanical dispatch loop (no tier/L3)?
source: '0221-layered-fo rework (2026-06-19): the load-bearing premise of the Haiku-operable/layered FO is that a weak model can reliably DRIVE the mechanical dispatch loop (boot→dispatch→advance→terminalize), needing help only on judgment. Untested — opus itself deviates (5xs). This is the cheapest path that, if it fails, invalidates the tier bet (72, standing-L3, kt-full). Run it FIRST.'
status: backlog
score: 0.85
sprint: 0221-layered-fo
group: validation
id: mvctb79y19fvhbsepyagdd8f
---

Test the riskiest, cheapest path first: a live Haiku-model FO drives ONE entity through the mechanical loop end to end — boot → dispatch a worker → review the report → advance → terminalize/merge — with NO tier self-identification, NO standing `level-3-judge`, NO gate-verdict escalation (all of that is deferred-72 / `kt-full`). The question is binary: can a weak model hold the contract loop at all?

If NO → the "Haiku-operable" goal and the tier bet (72) need rethink; `kt-full` stays shelved. If YES → `kt-full` (Haiku + L3 escalation) becomes worth running as the cut proof.

The deliverable is checkable: a reusable live harness (build on `internal/ensigncycle/haiku_loop_spike_live_test.go`) **plus** the recorded, durable-state-graded result of running it. Grade on DURABLE STATE only — terminal `verdict:` frontmatter + the entity's committed status transitions — never transcript phrasing (README proof policy + 72's grading discipline).

## Seed acceptance criteria (firm up in ideation)
- AC-1: A live Haiku-model FO drives one entity from initial dispatch to terminal, and the durable workflow state (terminal `verdict:` + committed status transitions) shows the loop completed without a human or stronger-model FO taking over the mechanical steps. Graded by reading the archived entity's state, not the drive transcript.
- AC-2: The spike records each Haiku deviation class it observes (e.g. skipped `spacedock new`, broad-search after zero-discover, exit-before-terminalize) against the captured transcript — the finding is the deliverable whether the verdict is can-hold or cannot-hold.

## Out of scope (deferred-72 / kt-full)
Tier self-identification (`«fo.tier»`), the standing `level-3-judge` mod, gate-verdict escalation, the `### Gate Verdicts` durable line. This spike isolates the MECHANICAL loop; the judgment-escalation half is the deferred bet this spike gates.
