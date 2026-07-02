---
title: "Reuse-advance messages go through a dispatch-built file pointer instead of FO-hand-assembled verbatim stage sections"
status: backlog
group: tooling
source: "fable-token-trim-scout analysis 2026-07-02 (captain-ordered fresh-angle token review): the reuse path explicitly does NOT route through dispatch build (fo-dispatch-core.md:40) — the FO hand-assembles a SendMessage embedding the full README stage subsection verbatim (claude-fo-dispatch.md:40), costing a README section read PLUS a verbatim echo per advance, ~400-800 tok x ~10 advances in a 5-entity/3-stage session = 4-8k tok. The biggest recurring dispatch-path cost found; fresh dispatch already proves the file-pointer mechanism (~175-char prompt pointing at a written dispatch file)."
id: jhe1244c8cjdymfbnnpnsvsw
---

## Problem
Advancing a reused worker to its next stage is the only dispatch-shaped message the FO still assembles by hand: the contract requires copying the next stage's full `### Stage definition` subsection from the README verbatim into a SendMessage, plus the checklist and continuation instruction. Every advance pays the section read + echo twice over; initial dispatch already solved this shape with `dispatch build` writing the assignment to a file and emitting a short pointer prompt.

## Desired direction (for ideation to refine)
A `dispatch build --advance` mode (or `dispatch advance` verb) that assembles the advancement message — next stage name, stage definition, completion checklist, continue-on-entity instruction — writes it to a dispatch file, and emits the short pointer + SendMessage-ready fields; the FO's reuse-advance becomes a one-line pointer message. Contract prose (fo-dispatch-core.md reuse section, claude-fo-dispatch.md reuse-advance template) updates to route through it.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- A reuse advance sends O(pointer) tokens instead of O(stage-section) — measured against the current template shape.
- The advancement file carries the same content contract the current template names (stage definition verbatim, checklist, entity path, commit-before-signal instruction).
- Contract prose routes the reuse path through the helper; live ensign behavior on a pointer advance verified (touches skills/**, so claude-live gates the merge).
