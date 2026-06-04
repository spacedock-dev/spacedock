---
id: t3w3s0q6a89me2kjkrgpz7nq
title: Extract Gate Presentation from first-officer-shared-core into a lazy spacedock-owned skill
status: backlog
source: "captain (2026-06-04) — token-efficiency decomposition of first-officer-shared-core.md (~9,730 tok, the largest boot-read file). Gate Presentation (the template + captain-facing assembly rules) is judgment/format prose needed only when presenting a gate, not in a session's first turns; defer it off the eager boot read via the zd lazy-skill pattern."
score: "0.33"
worktree:
started:
completed:
verdict:
issue:
---

`first-officer-shared-core.md` loads in full at every FO boot, but its `## Gate Presentation` block (the format template + the captain-facing assembly rules) is only needed at a gate. Lift it into a lazy spacedock-owned skill loaded via `Skill(skill=...)` at the gate-presentation point, the way zd (#291) lifted the team lifecycle to `using-claude-team`.

## Proposed approach (ideation firms)

Move the `## Gate Presentation` section (template + the captain-facing assembly rules: lede-first/decision-last, chosen-direction-required, cite-the-report, reviewer-tiers, etc.) into a new lazy skill; replace it in the FO core with a `Skill(skill=...)` invocation anchored at the gate-handling step. Judgment/format content → lazy skill (NOT a binary command — it is FO prose, not mechanical ceremony).

**Two load-bearing constraints (from the zd audit + the superpowers comparison):**
1. **Faithfulness** — the assembly rules are welded into how the FO adjudicates gates; a dropped clause is a correctness regression. zd-grade faithfulness audit (normalized diff of moved text vs pre-change).
2. **Load-trigger discoverability** — the `Skill()` invocation must sit at the gate-handling point in the always-on skeleton, or the FO won't load it when a gate arrives.

Keep it spacedock-owned (no external superpowers dependency).

## Acceptance criteria (seed)

- **AC-1 (seed):** The Gate Presentation block is ABSENT from `first-officer-shared-core.md` and PRESENT in the new skill; the FO core carries a `Skill()` invocation at the gate point — verified by an instruction-text oracle (the zd AC-1 pattern: fingerprint-absent-from-core / present-in-skill + invocation-present).
- **AC-2 (seed):** Faithfulness — moved text semantically complete (normalized diff; host-neutrality/portability oracles green).
- **AC-3 (seed):** No regression — a live FO drive that reaches a gate presents it correctly (closes via fresh session against the built plugin).

## Notes

Split from the umbrella analysis (see the binary-simplification roadmap, refreshed 2026-06-04, and sibling tasks). Template + faithfulness-audit pattern: zd `extract-team-orchestration-skill` (#291). Siblings: `feedback-rejection-flow-skill-extraction` (same lever), `pr-complete-binary-command` (Merge-and-Cleanup → binary, the ceremony counterpart).
