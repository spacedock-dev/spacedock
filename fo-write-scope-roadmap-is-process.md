---
id: kbhtzfkdy1a4az56f383pmpe
title: The write classifier blocks docs/roadmap as product, but it is the FO's own planning surface
status: backlog
source: "0260 Commander drive, 2026-07-20. The FO needed to add a boot-order line to its own sprint's Commander package and could not: `docs/roadmap/**` is classified blocked-product, so it required an explicit captain path-class grant to write a file no worker builds and no test covers. The captain granted it and asked for this task."
started:
completed:
verdict:
score:
worktree:
issue:
---

`docs/roadmap/**` is classified `blocked-product` in the FO write classifier, but it is process the FO operates rather than product built under test.

## Problem

`skills/first-officer/references/fo-write-core.md`'s classifier lists `docs/roadmap/**` alongside `cmd/**`, `internal/**`, `**/*_test.go`, `skills/**`, `agents/**`, `plugin.json`, `.github/**`, `fixtures/**` and `docs/dev/_mods/**` — code, tests, shipped scaffolding, release machinery. The rule for that class is "do not write; route through a dispatched worker."

But `docs/roadmap/` contains only sprint folders and roadmap documents: sprint `index.md` strategy, `staff-review*.md` readiness analyses, `dispatch-sprint-execution.md` Commander packages, plus `README.md` and a couple of architecture notes. Nothing there is built by a worker under test. It is the surface the Shaping FO authors and the Commander reads — the FO's own planning artefacts.

The classifier already draws exactly this distinction for a sibling case, and draws it the other way. `docs/dev/README.md` is `allowed-process` with the stated rationale: "The FO may edit the workflow README it operates because that file defines process, not the product being built." And the off-limits list is explicit that the workflow README "is process the FO owns, not product, so it is NOT in this list." The roadmap has the same character and the opposite classification.

Observed cost, 2026-07-20: driving 0260, the FO identified that a post-compaction successor would have no pointer to the Commander debrief, and that the correct home for that pointer was the Commander package's own boot order — a file in `docs/roadmap/0260-proportionality/`. It could not write it. The options were an exact-target captain grant, dispatching a worker to insert one line, or putting the pointer somewhere less correct (the workflow README, which is durable process text, for a transient session condition). The captain granted the path class. Note the asymmetry that makes this more than an inconvenience: a Shaping FO authors these files as its primary deliverable, so either it operates under a standing grant every session, or the classifier is wrong.

## Proposed approach

Ideation fills this in. The question is whether the whole path class moves or only part of it.

Candidate: reclassify `docs/roadmap/**` as `allowed-process`, with the same rationale line the workflow README already carries, and add it to the "process the FO owns" carve-out in the off-limits paragraph.

Consider before deciding: whether any roadmap file is genuinely product (`bootstrap-roadmap.md` is referenced by the dev workflow's Testing Resources table as stage-specific required tests, which is a different role from sprint strategy); whether a sprint index that a validator reads as an authority should be FO-writable mid-drive without a gate; and whether the correct unit is `docs/roadmap/**` or only the sprint-folder shape `docs/roadmap/NNN-*/`.

The classifier block is delimited by `FO-WRITE-CLASSIFIER:START/END` markers and `internal/contractlint/fo_write_core_mutation_gate_test.go` binds it, so any change has an existing mechanical check — verify what that test asserts before editing, and expect it to red on an unconsidered edit.

## Out of scope

Widening FO write authority generally. Changing the `allowed-state` or `blocked-product` treatment of any other path. The override mechanism itself, which worked exactly as designed here: the FO stopped, the captain granted the class explicitly, and the FO quoted the grant before writing.

## Acceptance criteria

Ideation fills these in. The value AC must measure whether an FO can perform its own planning-surface work without a per-session grant, against the observed baseline that it currently cannot — and the mechanical check is `fo_write_core_mutation_gate_test.go`, which should red if the classifier is edited in a way that contradicts its assertions.
