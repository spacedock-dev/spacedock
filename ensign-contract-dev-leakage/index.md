---
id: ep0ra3zjf4hhkhx5rrkwsxbb
title: Universal ensign contract has absorbed dev-workflow assumptions (TDD, code-only deliverables, "CODE only" worktree) — re-home dev policy out of the shared core
status: ideation
source: session-10 detached audit (Task waymqcmru, 2026-06-03) — overall verdict MATERIAL-PRESENT; validated the captain's grilling instinct ("why am I seeing dev-workflow specific ones like TDD in the universal contract?")
score: "0.19"
worktree:
started: 2026-06-03T07:09:59Z
completed:
verdict:
issue:
---

The shared ensign contract (`skills/ensign/references/ensign-shared-core.md`) is loaded by **every** dispatched ensign across all three commission templates (development, experiment, refinement). A detached adversarial audit confirmed it has absorbed dev-workflow-specific discipline — TDD authoring, a code-shaped closed deliverable list, and a "CODE only" worktree assumption — and asserts them as universal. The FO contract was already corrected to deliverable-shape-agnostic phrasing (`first-officer-shared-core.md` lines 342-353); the ensign contract is the **uncorrected twin**.

## Problem

`ensign-shared-core.md` lines 22-27 carry a `## Working Practices` block that duplicates dev-workflow policy verbatim into the portable contract. The project's own docs already tag this same policy as dev-only opt-in:
- `docs/dev/README.md` line ~110: "This is dev-workflow policy: an AC's proof here is code/command/state. A non-development workflow's AC proof may legitimately be a published artifact, a metric, or a human review."
- `skills/commission/references/templates/development.md` lines ~111-117: a "Recommended practices (opt-in)" carve-out stating "the universal first-officer contract does not impose them."

The shared ensign core strips that qualifier and presents the rules as universal ensign discipline, breaking at least two non-dev template classes: refinement entities whose deliverable **is** prose (PRD, outreach reply, content — `refinement.md` says they never touch the repo), and experiment entities whose deliverable is a hypothesis/analysis verdict against pre-registered success criteria.

**Material findings from the audit (Task waymqcmru):**
- **L2-01 (umbrella):** The entire `## Working Practices` section duplicates dev policy into the universal contract. Lift it out; re-home each principle.
- **L1-F1 / L2-02:** TDD bullet (line 24, "Write the failing test first… The test is what the gate judges") asserts test-first as the universal authoring contract — false for refinement (gate judges the artifact body) and experiment (gate judges evidence vs. pre-registered criteria) workflows.
- **L2-03:** "Every task produces a real, checkable change" + its escape clause under-serve prose-only deliverables (refinement/experiment bodies that are checkable via pre-fixed success criteria, not external files).
- **L2-05:** The Split-Root State Contract (line ~35) hard-codes "the worktree isolates **CODE only**" — excludes non-code work product (experiment evidence files, refinement attachments). Same fix needed in **both** `ensign-shared-core.md` and `first-officer-shared-core.md`.
- **L3-1:** Structural twin of the FO's already-corrected four-principle block — the ensign block is the uncorrected version.
- **L4-P2 / runtime adapters (Polish):** `claude-ensign-runtime.md` and `codex-ensign-runtime.md` enumerate "worktree path" as a top-level always-present assignment field, while the shared core already uses the substrate-neutral "workflow location" for the same slot.

## Proposed approach

Mirror the FO-contract correction (`first-officer-shared-core.md` 342-353) onto the ensign contract:

1. **Lift the `## Working Practices` block out of `ensign-shared-core.md`.** Re-home each principle to its rightful workflow:
   - TDD + the closed code-shaped deliverable list → `development.md` Recommended-practices / Adoption section (and `docs/dev/README.md`, which already carries the deliverable-proof policy).
   - `experiment.md` / `refinement.md` already encode their own analogues (pre-registered success criteria; draft/review/polish) — no new homes needed.
   - If the shared core should retain a universal restatement, phrase it discipline-neutrally: "fix the success criterion before you start, in a form a later check can read back" — covers TDD (dev), pre-registered criteria (experiment), acceptance criteria (refinement) without asserting any one tradition.
2. **Replace "CODE only" with "the work product only"** (or "the deliverable artifacts only") in both `ensign-shared-core.md` and `first-officer-shared-core.md` so the worktree-vs-state-checkout boundary is substrate-neutral.
3. **Rename "worktree path" → the shared core's neutral term** ("workflow location" / "workspace path") in both runtime adapters.
4. **Add lock-in oracles** (the proven Phase 0.A / qs cycle-3 pattern — `TestNoAuditTrailExposition`, `TestNoCrossFileRestatement`): a Go test asserting the universal shared core is free of the dev-only vocabulary and that the dev-specific homes carry it, to prevent re-inflation.

The shared ensign core should retain only universals: read the assignment, follow the stage definition's proof requirements, commit the entity to the state checkout, signal completion.

## Out of scope

- Rewriting the dev-workflow disciplines themselves — they are correct *for dev*; this is a re-homing, not a policy change.
- Changing the FO contract's already-corrected four-principle block (only its "CODE only" line, shared with the ensign core, is in scope — finding L2-05).
- The experiment/refinement template bodies beyond confirming they already carry the re-homed analogues.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action. Proof is presence/absence over the instruction files (legitimate per the README: "a presence check over instruction files proving they carry a required clause or stay free of a banned token is proof at the claim's own level") plus the lock-in oracle. Ideation refines.

**AC-1 — The universal ensign shared core no longer asserts dev-only authoring discipline.**
Verified by: a Go oracle test asserting `ensign-shared-core.md` is free of the dev-only test-first vocabulary ("failing test", "feature or bugfix", "the test is what the gate judges") in a universal-imperative position; the test fails if the Working Practices block is reintroduced.

**AC-2 — The dev disciplines survive in their dev-specific homes.**
Verified by: the same (or paired) oracle asserting `development.md` / `docs/dev/README.md` carry the re-homed TDD + deliverable-proof policy, so the lift relocated rather than deleted the guidance.

**AC-3 — The worktree-isolation boundary is substrate-neutral in both shared cores.**
Verified by: grep/oracle confirming neither `ensign-shared-core.md` nor `first-officer-shared-core.md` contains the "CODE only" phrasing, and both express the boundary as work-product-agnostic.

**AC-4 — Both runtime adapters use the shared core's neutral location vocabulary.**
Verified by: oracle/grep confirming `claude-ensign-runtime.md` and `codex-ensign-runtime.md` no longer enumerate "worktree path" as a universal assignment field where the shared core uses "workflow location".

**AC-5 — The contract still drives a real workflow correctly.**
Verified by: a live ensign dispatch on the swept ensign contract completing a stage cleanly (the strongest signal that load-bearing meaning was preserved — same bar Phase 0.A met by driving on both sonnet and opus).

## Test plan

- Go oracle tests for AC-1..AC-4 (the proven Phase 0.A negative-proof + cross-file-restatement pattern). Cost: low — text invariants over instruction files.
- One live ensign dispatch for AC-5. Cost: medium (one live cycle).
- High-stakes surface (shipped contract/scaffolding) → a detached adversarial audit is required before merge per the dev README's validation stage.

## Notes

- The full audit synthesis (6 material + 2 polish findings, MATERIAL-PRESENT) is Task `waymqcmru` from session 10. Findings are quoted above; re-run the audit if deeper detail is needed.
- This is the ensign-side completion of the contract-simplification arc shipped in 0.19.4 (Phase 0.A swept the FO + ensign cores for *length*; this sweeps the ensign core for *dev-leakage* — a different axis the captain's grilling surfaced).
- Sequence: 0.19.5. Independent of the sonnet-live-CI-flake entity, though both touch the ensign-contract surface.
