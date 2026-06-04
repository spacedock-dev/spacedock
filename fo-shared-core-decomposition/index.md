---
id: mdyyxf7dyannrxpjz6a47dvh
title: Decompose first-officer-shared-core.md — defer event-specific ceremony off the always-loaded boot read
status: backlog
source: "captain (2026-06-04) — token-efficiency analysis: the FO boot read loads first-officer-shared-core.md (~9,730 tok) IN FULL every session, but most of it (Gate Presentation, Merge-and-Cleanup, Feedback Rejection Flow) is event-specific, not needed in a session's first turns. superpowers keeps each skill ≤~2.6k tok and lazy; our boot read is a ~6x-median monolith loaded eagerly. zd (#291) proved the lazy spacedock-owned skill decomposition mechanism."
score: "0.35"
worktree:
started:
completed:
verdict:
issue:
---

`first-officer-shared-core.md` is the single largest boot-read file (~9,730 tok) and loads unconditionally at every FO boot. The deciding insight from the superpowers comparison: the difference that matters isn't size, it's *when* it loads. The boot read should carry only what every session needs in its first turns; by that test, most of this file is event-specific ceremony that should be pulled when its event fires.

## Problem

Eager monolith. Gate Presentation, Merge-and-Cleanup (+ Ship-Local Ceremony), and Feedback Rejection Flow are each needed only at specific event-loop moments, yet all load at boot. Startup + the dispatch-loop skeleton + Write-Scope/policy are the genuine always-on core.

## Proposed approach (ideation firms) — route each block by content TYPE

The superpowers comparison says the two slimming levers sort cleanly by content:
- **Ceremony (mechanical) -> binary command, NOT a skill.** Merge-and-Cleanup / Ship-Local Ceremony should become the roadmap's #3 `spacedock pr complete` (and #1 `state sync` / #2 `dispatch advance`), which removes the prose outright. Track those via the binary-simplification roadmap; this entity coordinates with them rather than duplicating.
- **Judgment / format -> lazy spacedock-OWNED skill** (the zd pattern), loaded via `Skill(skill=...)` AT its trigger point in the always-on skeleton: Gate Presentation (template + captain-facing assembly rules); Feedback Rejection Flow.
- **Irreducible policy -> stays** always-on: Startup, Status Viewer, ID Styles, Single-Entity Mode, FO Write Scope, the dispatch-loop skeleton.

**Two load-bearing constraints (from the zd audit + the superpowers comparison):**
1. **Faithfulness.** These blocks are welded into control flow (the gate's AC-refusal, the teardown retry). A botched extraction drops a load-bearing clause or alters semantics — a correctness regression, not just a token win. Each extraction needs zd-grade faithfulness auditing (normalized diff of moved text vs pre-change).
2. **Load-trigger discoverability.** Each deferred block needs a `Skill()` invocation anchored at its use-point in the always-on skeleton (the way zd anchored using-claude-team at the Team-Creation step), or the FO won't load it when needed — a correctness gap. Consider a `using-superpowers`-style index in the skeleton naming what gets pulled when.

Do NOT take a runtime dependency on the external superpowers plugin — keep extracted blocks spacedock-owned (the zd rationale).

## Out of scope

- The binary-command moves themselves (#1/#2/#3) — owned by the roadmap; this entity hands Merge-and-Cleanup to #3 rather than skill-ifying it.
- The token-budget guard (`contract-token-budget-guard`) — the forcing function that should land first.
- The ensign contract (ep already did that side).

## Acceptance criteria (seed)

- **AC-1 (seed):** The FO boot read (always-on core) drops materially — target the ~3–4k tok always-on skeleton — verified by a measured before/after token count and the deferred blocks no longer appearing in the boot-read files.
- **AC-2 (seed):** Faithfulness — each moved block is semantically complete (no dropped clause/caveat), verified by a normalized diff vs the pre-change file and the existing host-neutrality/portability oracles staying green (zd-grade).
- **AC-3 (seed):** Each deferred block has a load-trigger anchored at its use-point in the always-on skeleton, verified by an instruction-text oracle (the block's fingerprint is ABSENT from the boot read and PRESENT in the deferred skill, and the skeleton carries the `Skill()` invocation at the right section — the zd AC-1 pattern).
- **AC-4 (seed):** No behavior regression — a live FO drive (gate -> merge) on the decomposed contract completes cleanly (the zd AC-3 / p4 live-verification bar; closes via a fresh FO session against the built plugin).

## Test plan (seed)

- Instruction-text oracles for AC-1/AC-3 (extend the proven `skill_text_test.go` / `using_claude_team_test.go` pattern). Faithfulness diff for AC-2. One live FO cycle for AC-4. High-stakes surface (the FO operating contract itself) -> detached adversarial audit before merge.

## Notes

- Sibling to `contract-token-budget-guard` (file the cap first) and the binary-simplification roadmap (refreshed 2026-06-04 — the refresh's headline reframe: install the budget, then route each block ceremony->binary / judgment->lazy-skill / policy->stays). zd (`extract-team-orchestration-skill`, #291) is the proof-of-concept and the faithfulness-audit template.
