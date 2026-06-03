---
id: cbtjetk8e17y777zfynw1sgc
title: Staff-review & detached-audit reports escape the conciseness standard — they're ad-hoc conditional dispatches, not formalized stages
status: backlog
source: captain (2026-06-03, session 12) — "the staff review result is very verbose. maybe because this is not a regular workflow stage but our conditional dispatch? the reporting should be subject to the same conciseness standard."
started:
completed:
verdict:
score: 0.24
worktree:
issue:
---

The staff review (requested at a complex ideation gate) and the detached adversarial audit (required at validation for high-stakes surfaces) are dispatched via a HAND-ASSEMBLED reviewer prompt — outside `spacedock dispatch build`. So they inherit none of the conciseness discipline the formalized surfaces impose: the gate-presentation format (15-25 lines, lede-first, cite-don't-paste, no format-pedantry) and the stage-report protocol. The result is review reports that re-narrate everything verified-correct, paste evidence inline, and run many screenfuls — high signal but low signal-to-noise, captain-observed in session 12.

## Problem

Two formalized report surfaces already carry a conciseness standard:
- Gate presentation (`first-officer-shared-core.md` `## Gate Presentation` + `### Captain-facing assembly rules`): target 15-25 lines, verdict once, cite the report don't paste it, render gists not full text.
- The stage-report protocol (ensign contract): structured DONE/SKIPPED/FAILED accounting.

The review/audit report has NO such standard. The detached audit's output IS partly specified (`docs/dev/README.md` validation stage: two tiers Material / Polish) — but only the TIERS, not the per-finding conciseness or a length budget. The staff review (README `## ideation` "Staff review") has no report format at all. Because the FO assembles these prompts ad-hoc, each review's verbosity depends on how the FO worded that one prompt.

## Proposed approach

{Ideation refines. Direction: define a concise review/audit report format as a sibling of the gate-presentation format, and make the FO emit it on every review/audit dispatch.}
- A `## Review & Audit Report Format` convention (likely in `first-officer-shared-core.md`, beside `## Gate Presentation`): verdict-first (`APPROVE | REWORK — {one-line reason}`); Material/Polish tiers where each finding is ONE line = claim + `file:line` evidence + the ask; NO re-narration of what's correct unless a finding hinges on it; cite `file:line`, don't paste; a line budget that scales with finding count.
- Make it load-bearing rather than prose-only: the cleanest lever is a `dispatch build`-style helper (or a `dispatch review`) that emits the reviewer/auditor prompt WITH the report-format constraint baked in — so the FO can't forget it and the format is enforced at assembly, not by FO discipline. (Mirrors how `dispatch build` already standardizes the ensign assignment.)

## Out of scope

- Changing WHAT a review/audit checks (the Material/Polish tiers, the refute-the-validation discipline) — only the report's shape/length.
- The gate-presentation format itself (already correct; this mirrors it).

## Acceptance criteria

{Ideation defines. The teeth: the standard must be enforced by something that can fail — a helper that emits the constrained prompt (proven by a test over its output), or a presence check that the contract defines the format. Avoid a prose-only "the contract says be concise" with no gate.}

## Test plan

{Ideation. Likely: a contract presence check + (if the helper route) a Go test over the emitted review-prompt. Low cost.}
