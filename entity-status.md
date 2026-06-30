---
title: Defer the Status Viewer + Issue Filing reference out of the boot-resident core (entity-status)
status: ideation
group: cleanup
id: 84521x23qhnvy0xy6st4qwmh
sprint: 0240-lean-contract
started: 2026-06-30T16:20:13Z
---
The boot-resident FO core (`skills/first-officer/references/first-officer-shared-core.md`) loads ~675 tokens of phase-specific reference at greet that isn't needed until a specific phase: the **Status Viewer** section (~657 tok — `status --set` field docs, the canonical invocations, the Captain-Facing State Display rendering — needed only when answering an ad-hoc status question or mutating state) and **Issue Filing** (~18 tok — rare). Defer both into one lazily-loaded `entity-status` reference, loaded on first status-query / mutate / issue-file, using the same name-a-pointer / defer-the-body pattern the dispatch and merge modules already use. Sibling to z4's Dispatch/Merge pointer consolidation, which provides the registry index.

Measured: ~675 tok ≈ 9–10% of the ~7,170-token boot core (part of the ~32% total deferrable surface; FO Write Scope / ID Styles / Probe discipline are separate sibling tasks).

## Acceptance criteria
- **AC-1 (occupancy — the reason this exists)** — first-officer-shared-core.md drops by the deferred sections' size (~675 tok), measured net-NEGATIVE vs origin/main; the deferred content lives in a new lazily-loaded reference, off the boot path.
- **AC-2 (correctness — the greet must not depend on the deferred content)** — a boot that greets-and-stops (the common interactive path) completes WITHOUT loading the entity-status reference: verified by exercising the greet on a fixture and confirming it renders the state summary + any ready gate using only Startup + `«state.*»` + `status --boot`, never the Status-Viewer display rules. (The staff-review caveat, made into a gate.)
- **AC-3 (reachability)** — the FO loads the entity-status reference at its trigger (a status question, a `--set` mutation, an issue filing) via a boot-core pointer; a status-query / new-entity / mutate path still resolves all the moved guidance.
