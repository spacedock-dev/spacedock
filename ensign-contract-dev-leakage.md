---
title: Re-home dev-only discipline out of the universal ensign core (per-dispatch occupancy)
status: ideation
sprint: 0240-lean-contract
group: cleanup
id: scr2rx4589p7j6mpgh50hdct
started: 2026-06-30T15:51:56Z
---
The ensign shared core (`skills/ensign/references/ensign-shared-core.md`) loads on EVERY worker dispatch (per-dispatch tier — a token here recurs every spawn). It bakes dev-workflow-specific discipline — TDD framing, "code-only deliverables," the "CODE only" worktree rule — into the UNIVERSAL ensign contract, even though the ensign also runs non-dev workflows (ticket, experiment, survey). Re-home that dev-only prose into the dev-shape scaffolding (loaded only by dev workflows), leaving the universal core host-/stage-neutral. Saves per-dispatch tokens (the leakage rides every ensign spawn) AND is a contract-correctness fix (the universal core shouldn't assert dev assumptions).

Sourced from `_sprint-notes.md` (`ep0ra3z…`). Companion to `read-guidance-redundant-with-grep` (82k) — both trim the ensign core; coordinate edits to avoid collision.

## Acceptance criteria
- **AC-1 (correctness)** — the universal ensign core asserts no dev-workflow-specific discipline; a token/structural oracle confirms the dev-only markers (TDD, "CODE only" worktree, code-only deliverable) are ABSENT from the universal core and PRESENT in the dev-shape scaffolding, with a non-vacuity control.
- **AC-2 (per-dispatch occupancy)** — the universal ensign core's token footprint drops by the re-homed prose, measured net-NEGATIVE vs origin/main; the dev-shape scaffolding's gain is offset (prose moved, not duplicated).
- **AC-3 (no behavior loss for dev workflows)** — a dev-workflow ensign dispatch still receives the dev discipline via the dev-shape scaffolding; verified by a dispatch-build test on a dev fixture that the assembled assignment carries it.
