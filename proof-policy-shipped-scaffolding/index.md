---
id: eykbxewyrxg1vyy9fhzpxee7
title: Port the proof-policy (no tautological string-match over LLM-ingested files) to the shipped dev-shape scaffolding
status: backlog
source: "captain (2026-06-04) — after the dev-README proof-policy fix (f8b257cf, closing the text-claim 'proven at its own level' loophole), port the same principle to the shipped scaffolding so every new dev workflow + the FO/ensign contracts carry it. AI-engineer propagation list: skills/commission/references/templates/development.md, first-officer-shared-core.md, ensign-shared-core.md."
score: "0.33"
started:
completed:
verdict:
worktree:
issue:
---

The dev-workflow README now states the proof policy bluntly (commit `f8b257cf`): a string/substring/regex match over any instruction file the model reads NEVER satisfies a behavioral acceptance criterion; the one test is "does the expected value come from a source OTHER than the file under test?" — if not, it is a tautology that proves nothing. That fix lives only in THIS workflow's README. Port the same principle to the **shipped dev-shape scaffolding** so it governs every newly commissioned workflow and the universal FO/ensign contracts, not just this one.

## Problem

The carve-out to DELETE ("text claim proven at its own level") existed only in `docs/dev/README.md` and is now gone. But three shipped surfaces carry the *parallel* policy language and lack the explicit independent-source test — so a future worker/validator/commissioned-workflow can still rationalize a tautological presence check as proof:

- **`skills/commission/references/templates/development.md`** (~L135, the live-scenario-for-runtime-claims paragraph) — the **dev-shape workflow template** every new dev workflow is generated from. Without the rule here, each new workflow inherits the gap. **Primary target (captain-named).**
- **`skills/first-officer/references/first-officer-shared-core.md`** (~L298, "Prefer a code gate over a prose-only rule") — the FO contract; already says "wording-present is not behavior" but not the independent-source test the FO adjudicates gates by.
- **`skills/ensign/references/ensign-shared-core.md`** (~L25, "Prove by exercising, not by re-reading") — the ensign contract; closest existing match, lacks the explicit test.

None reproduces the deleted carve-out, so this is an **affirmative add** (the independent-source test + the "string-match over an ingested file never satisfies a behavioral AC" line), not a removal.

## Proposed approach (ideation formalizes; AI-engineer draft exists)

Add, in each surface's existing voice, the load-bearing test: **"Does the expected value come from a source OTHER than the file under test? If not → tautology → proves nothing. For a behavioral claim, run the behavior and observe durable state/output."** Keep it consistent across all three so the FO (gate adjudication), the ensign (self-check), and a commissioned workflow (its own README) read one rule. The AI-engineer proposal (this session) drafted the dev-README edits and the propagation list; reuse its framing.

## The proof wrinkle (READ THIS — it is the whole point)

**This task cannot prove itself with a presence oracle** — that is exactly the banned thing it ports. Per the new policy:
- The **text half** (the clause is present in each file) is real authoring work but is NOT an acceptance criterion on its own.
- The **behavioral half** — that an FO/validator/ensign actually *applies* the rule (rejects a tautological-only AC, demands a behavioral proof) — must be a **live drive** that runs the behavior and observes the durable outcome. Ideation must design that behavioral scenario (e.g. a live scenario seeding an entity whose only AC is a presence oracle, asserting the validator/FO REJECTS it; graded on durable state — the rejection, the routed feedback — never transcript). This is the same shape as gq's escalation scenario.

If the behavioral scenario proves infeasible to drive, that is a finding for the gate, not a license to fall back on a presence check.

## Acceptance criteria (preliminary — ideation formalizes, and per the new policy the text-presence facts below are NOT standalone ACs)

**AC-1 (behavioral) — A validator/FO presented with an entity whose only proof is a presence oracle over an instruction file REJECTS it and demands a behavioral proof.**
Verified by: a live scenario (`internal/ensigncycle` shared scenario or `internal/livescenario`) that seeds such an entity, drives a real reviewer, and asserts the durable end-state is a rejection/routed-feedback — with an offline negative proving the grader reds on a wrongly-accepted end-state. (The behavioral proof the policy itself demands.)

**AC-2 (consistency invariant) — the three shipped surfaces carry the rule.**
Verified by: NOT a standalone AC — it is the text half. If a check is wanted, it must bind to an independent source (e.g. a doc↔code or cross-file consistency invariant that can diverge), never a bare presence assert; and it never substitutes for AC-1.

## Test plan

- **Behavioral (the real proof):** the AC-1 live rejection scenario + its offline negative. Cost: live-gated; ideation sizes it.
- **High-stakes → detached adversarial audit before merge:** this edits shipped contract/scaffolding (a named high-stakes surface) AND it is the proof-policy itself — the audit must try to defeat AC-1's grader (a transcript that narrates rejection while the durable state accepted).
- **Open edge (carry the captain's pending call):** the doc↔code-sync exception (a sync invariant that binds to code — keep it as a consistency check but bar it from satisfying a behavioral AC, vs. the blunter "zero exceptions for model-ingested files"). The dev-README edit (f8b257cf) left this implicit; the ideation should state the chosen rule explicitly here too.

## Notes

Provenance: the wm detached audit (this session) empirically proved a negation-wrap defeats a presence oracle; the captain ruled presence-over-ingested-file checks tautological and banned as behavioral proof. Sibling: the dev-README fix `f8b257cf` (the THIS-workflow half, already shipped). The AI-engineer proposal (this session) carries the drafted wording + the full propagation list.
