---
title: Defer FO Write Scope + ID Styles out of the boot-resident core (sibling to entity-status)
status: backlog
sprint: 0240-lean-contract
group: cleanup
id: k408d2ydgj7s3s9yg7csyw81
---
Sibling to `entity-status`. The boot-resident FO core (`skills/first-officer/references/first-officer-shared-core.md`) loads ~1,027 tokens of write/filing-phase reference at greet that isn't needed until the FO writes or files: **FO Write Scope** (~696 tok — what the FO may write on main + the `spacedock new` atomic-create procedure, needed at write-time) and **ID Styles** (~331 tok — sd-b32/sequential/slug minting detail, needed at new-entity filing). Defer both into a lazily-loaded reference, loaded on first write / `--set` / `spacedock new`, using the same pointer/defer pattern as the dispatch/merge modules and entity-status.

Measured: ~1,027 tok ≈ 14% of the ~7,170-token boot core. After this + entity-status, the remaining deferrable surface (Probe & Ideation Discipline ~347, Single-Entity ~107, Mod Hook Convention ~147) is a smaller third sibling if wanted. Sequencing: coordinate with z4 #6 (the registry block) and entity-status — all three edit first-officer-shared-core.md.

## Acceptance criteria
- **AC-1 (occupancy)** — first-officer-shared-core.md drops by ~1,027 tok (Write Scope + ID Styles), measured net-NEGATIVE vs origin/main; the content lives in a lazily-loaded reference, off the boot path.
- **AC-2 (greet/boot independence)** — a boot that greets-and-stops completes WITHOUT loading the deferred reference; a boot that DISPATCHES loads it before the first frontmatter write or `spacedock new`. Verified on fixtures: greet (no load) vs first-dispatch (load-before-write). Confirm the boot's own `«state.sweep-merged»` write path does not depend on the deferred reference.
- **AC-3 (reachability — no guard lost)** — every Write-Scope rule + ID-Style minting path resolves when the FO writes/files; the FO loads the reference at its trigger via a boot-core pointer. FO Write Scope is load-bearing (it guards what the FO may touch on main) — confirm no write-permission guard is lost in the move.
