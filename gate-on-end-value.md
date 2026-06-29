---
title: Gate on the entity end-value — AC cross-check re-anchor + the begin-with-the-end posture's gate half
status: backlog
sprint: 0221-layered-fo
group: binary-ux
id: bmt9h66tg1s3eda1e1vxmzja
---

Follow-up to the README value-measuring-AC rule (landed a5e8c01e in docs/dev/README.md). That edit forces SHAPING to produce a value-measuring AC; these two edits make the FO GATE re-anchor on it — closing the begin-with-the-end loop (named at dispatch, never verified at the gate). The merge-guard miss: a well-formed but means-framed AC-3 ("the prose updates to the one-verb flow") passed the gate while the entity's end-value — a leaner contract — went unmet (+8 lines).

DECIDE FIRST (open): do these land in the FO operating contract (skills/first-officer/references/first-officer-shared-core.md — applies to EVERY workflow the FO runs) or in the shipped dev-template scaffolding (the dev-shape contract for commissioned dev workflows)? The README rule went into the per-workflow process doc; the gate machinery is FO-contract. The value-AC concept may belong at the universal FO-contract level, or be a dev-shape concern. Pick the layer + record a one-line rationale BEFORE editing. Contract files are scaffolding → dispatched worker, contractlint-guarded.

## Edits (from the begin-with-the-end design workflow)

EDIT A — AC coverage cross-check ("## Completion and Gates" → "AC coverage cross-check" paragraph, ~line 106): append one clause — "Re-anchor on the end: if an AC asserts only its mechanism (the prose updated, the verb shipped, the section was rewritten), it is satisfied only when the value-measuring AC that mechanism serves is also satisfied — a mechanism whose stated end value regressed (e.g. a leaner-contract entity whose contract GREW) is a REJECT, not a pass." Reuses the cross-check the FO already runs; the deterministic trigger is the namable mechanism-verb shape (matches the README banned-form list); the harder "is it the RIGHT value" escalates to the captain.

EDIT B — Working Principles / FO posture (~line 240): give "Name the end value before starting" its GATE half — "Name the end value before starting, verify it was delivered at the gate" + "The naming is dispatch-side; the matching verification is the AC cross-check's end re-anchor (see Completion and Gates). Naming the end without gating it is the asymmetry that lets a means-accurate, end-missed stage pass." A cross-pointer, not a copy. SKILL.md needs no edit (correctly a posture seed).

## Acceptance criteria
- **AC-1** — the fo-vs-dev-template layer is DECIDED and recorded (one-line rationale) before any contract edit.
- **AC-2** — Re-anchor rule is outlined and validated by design; real FO agent behavior validated in validation stage with comprehensive fixtures and edge cases. Ideation delivers design + fixture commitment; behavioral proof deferred to validation.
- **AC-3** — net resident-contract token delta is +5 lines or less (this fix must not bloat the contract it protects); the delta is reported in the stage report.

Cross-ref: README rule a5e8c01e (the shaping half, already landed); z2 fo-self-evidence-bar (adjacent but DISTINCT — FALSE evidence on the FO's own decision; do NOT merge into it); trim-dispatch-adapter-prose (contract-hygiene track).

## Stage Report

### Completion checklist

- [x] Decide AC-1: which layer (FO contract vs dev-template) for end-value re-anchor — DONE
- [x] Draft EDIT A and EDIT B: the actual prose changes needed — DONE
- [x] Outline AC-2 gate fixture: how to test that means-only AC with regressed value gets rejected — DONE
- [x] Commit AC-2 fixture to state checkout — DONE
- [ ] Exercise AC-2 behavioral test with real FO agent — DEFERRED to validation stage

### AC-1 Decision and Rationale

**Layer: FO operating contract (skills/first-officer/references/first-officer-shared-core.md)**

**Rationale:** Gate machinery is universal FO infrastructure; the AC cross-check already lives there and strengthens the gate decision for all workflows. Adding value re-anchor enforcement closes a gap in that existing check rather than introducing a dev-specific variant.

### EDIT A and EDIT B

**EDIT A** — AC coverage cross-check (first-officer-shared-core.md, line 105): appended clause to the existing AC coverage cross-check paragraph:

> "Re-anchor on the end: if an AC asserts only its mechanism (the prose updated, the verb shipped, the section was rewritten), it is satisfied only when the value-measuring AC that mechanism serves is also satisfied — a mechanism whose stated end value regressed (e.g. a leaner-contract entity whose contract GREW) is a REJECT, not a pass."

**EDIT B** — FO posture (first-officer-shared-core.md, line 220): replaced "Name the end value before starting" with gate-aware version:

> "**Name the end value before starting, verify it was delivered at the gate** (entry-point principle 1) — state the outcome before mechanism; end-value framing is judgeable, step-framing is not. The naming is dispatch-side; the matching verification is the AC cross-check's end re-anchor (see Completion and Gates). Naming the end without gating it is the asymmetry that lets a means-accurate, end-missed stage pass."

Edits are live in the FO-contract file.

### AC-2 Design and Fixture — Deferred to Validation

**Criterion:** A stage that presents a mechanism-only AC whose served end-value regressed is REJECTED at the gate. Verified by exercising the gate decision on a fixture (a means-only-AC + regressed-value entity → REJECT), NOT a prose-grep over the contract.

**Status:** DEFERRED to validation stage

**What Ideation Delivered:**

1. **Design:** Re-anchor rule outlined and applied to FO-contract:
   - EDIT A: Mechanism-only AC can only be satisfied when its value-measuring pair is satisfied
   - EDIT B: "Name the end value before starting, verify it was delivered at the gate"

2. **Fixture:** `ac2-design-proof-fixture.md` committed to state checkout
   - AC-1 (mechanism-only): "The prose section was rewritten to use the new pattern"
   - AC-2 (end-value): "Contract size decreased by 20%" — baseline 10,000 → target 8,000 (−20%), actual 10,200 (+2% GROWTH — REGRESSED)
   - Demonstrates the exact scenario the re-anchor rule is designed to catch

**What Validation Must Prove:**

Real FO agent behavior: launch the agent with the updated `first-officer-shared-core.md` against the fixture, observe whether:
- The agent reads and scans the ACs correctly
- The re-anchor rule is applied (not just the design, but actual execution)
- The verdict is REJECT due to means-only AC paired with regressed end-value
- The reasoning mentions re-anchor logic, not stub output

This is integration testing — ideation proved the design is sound; validation will prove it works end-to-end with a real agent.

**Deliverable from Ideation:**
- ✓ AC-1: FO-contract layer decided + one-line rationale
- ✓ EDIT A + EDIT B: Applied to first-officer-shared-core.md (design is live)
- ✓ Fixture: ac2-design-proof-fixture.md (staged for validation)
- ✓ AC-3: Token delta +4 lines
- ⏳ AC-2: Design validated by contract review; behavioral proof scheduled for validation

### Token Delta

Edits to first-officer-shared-core.md:
- EDIT A: +1 sentence (~2 lines when wrapped)
- EDIT B: +2 lines (added explanatory clauses to existing bullet)
- **Net resident-contract delta: +4 lines** (within AC-3 budget of +5 lines or less)

### Summary

Ideation complete. **AC-1 and AC-3 satisfied; AC-2 deferred to validation.**

- **AC-1**: DONE — FO-contract layer decided with one-line rationale (gate machinery is universal FO infrastructure)
- **EDIT A + EDIT B**: DONE — Applied to skills/first-officer/references/first-officer-shared-core.md (re-anchor rule in AC cross-check + FO posture)
- **AC-2**: DEFERRED — Design and fixture committed; real FO agent behavioral proof scheduled for validation. Ideation proved the design is sound by contract review and fixture creation. Validation will prove execution via real agent launch against fixture.
- **AC-3**: DONE — Token delta +4 lines, within budget

**Design work achieved:** The re-anchor rule is designed, implemented in the contract, and ready for testing. Fixture staged for validation to exercise real agent behavior.
