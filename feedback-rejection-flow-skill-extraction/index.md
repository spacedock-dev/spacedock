---
id: a9nte184whfmz8ajzn4n51yr
title: Extract Feedback Rejection Flow from first-officer-shared-core into a lazy spacedock-owned skill
status: backlog
source: "captain (2026-06-04) — token-efficiency decomposition of first-officer-shared-core.md. The Feedback Rejection Flow is needed only when a gate rejects (or a feedback stage recommends REJECTED), not in a session's first turns; defer it off the eager boot read via the zd lazy-skill pattern."
score: "0.31"
worktree:
started:
completed:
verdict:
issue:
---

`first-officer-shared-core.md`'s `## Feedback Rejection Flow` (cycle tracking, route-to-feedback-to-target, 3-cycle escalation, budget-probe/reuse, re-run-reviewer) only fires on a rejection. Lift it into a lazy spacedock-owned skill loaded via `Skill(skill=...)` at the rejection-handling point (zd pattern).

## Proposed approach (ideation firms)

Move the `## Feedback Rejection Flow` section into a new lazy skill; replace it in the FO core with a `Skill()` invocation anchored where the event loop detects a REJECTED gate / feedback recommendation. Judgment/process prose → lazy skill.

**Two load-bearing constraints:**
1. **Faithfulness** — the flow is welded into the gate/event-loop control flow (cycle counting, the reuse-vs-fresh + budget-probe decision, the bare-mode sequential variant in the runtime adapter); a dropped clause mis-routes a rejection. zd-grade faithfulness audit.
2. **Load-trigger discoverability** — the `Skill()` invocation must sit at the rejection-detection point in the always-on skeleton.

Keep spacedock-owned. Note the runtime adapters also carry a bare-mode feedback variant — decide at ideation whether that moves too or stays as an adapter seam.

## Acceptance criteria (seed)

- **AC-1 (seed):** The Feedback Rejection Flow block is ABSENT from `first-officer-shared-core.md` and PRESENT in the new skill; FO core carries the `Skill()` invocation at the rejection point — instruction-text oracle (zd AC-1 pattern).
- **AC-2 (seed):** Faithfulness — moved text semantically complete (normalized diff; oracles green); the runtime-adapter bare-mode variant stays consistent.
- **AC-3 (seed):** No regression — a live feedback cycle (reject → route → re-review) routes correctly.

## Notes

Split from the umbrella analysis (binary-simplification roadmap, refreshed 2026-06-04). Template: zd `extract-team-orchestration-skill` (#291). Siblings: `gate-presentation-skill-extraction`, `pr-complete-binary-command`.
