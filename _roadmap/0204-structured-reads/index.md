# 0204 — structured reads (0.20.4)

**Sprint:** the entities matching `sprint: 0204-structured-reads` — list current members with `spacedock status --workflow-dir docs/dev --where sprint=0204-structured-reads`. Membership and per-task state are the query, never enumerated or tracked in this doc.
**Theme:** make reading entities, reports, and the README cheap and structured.

## Goal (success criterion)

An FO or ensign reads only the part of a workflow file it needs — a section, the frontmatter, a CI summary — and mutating a file does not re-echo the whole thing. The backbone is a `spacedock status` read helper that returns a file's parsed frontmatter plus a section-heading→line-offset map, so a caller reads one section with `Read(offset, limit)` instead of the whole file. Around it the sprint trims the other recurring read/mutation sinks (the status listing's SOURCE render, the read-then-`set` staleness echo, whole-log CI reads) and makes the shipped commission templates defer to the operating contract. Proven by real behavior — helper output over fixtures, live section reads, token-flow measurements — never a prose-grep.

## Why

This session's dogfood surfaced the cost: 280-line entity bodies and stage reports read whole, a ~315-line README loaded at boot, 143KB CI logs — plus the Read-then-Bash staleness echoes. Reading is a recurring FO/ensign token sink. Two levers: a binary helper that returns a document's structure (read only the section you need), and a slimmer README/templates so the whole-file fallback is cheaper.

## Definition of Done

0.20.4 ships when, merged to `main`:
- **The `status` read helper** (the backbone) returns an entity/markdown file's parsed frontmatter + a section-heading→line-offset map, so a caller can read a target section with `Read(offset, limit)` instead of the whole file. Proven by Go tests over real helper output (correct FM; heading offsets that locate the right section bytes), plus a live exercise of an FO/ensign reading a section through it.
- **The other read/mutation-cost reductions the sprint scopes are resolved** — each member meets its acceptance criteria, or, where a spike proves a sink is not tool-fixable (e.g. the read-then-`set` echo if it turns out harness-inherent), it is recorded as a roadmap decision rather than forced into code.
- `v0.20.4` cut after the pre-cut antipattern audit is clean.

Membership (the specific reductions) is the query above, not enumerated here. The dev README's own slim-down is **FO-direct upkeep** (the FO owns and edits the workflow `README.md`), done alongside the helper so the slimmed README stays cheap to read. The commission-template **restructure** (lead-with-the-end + defer-to-contract) is now a sprint member; `ey` (the narrower proof-policy rule-port into the operating contract) stays separate and composes with it.

## Out of scope

Defined as the sprint forms (ideation / the Commander): whether the helper extends `spacedock status` vs a new `read` surface; output shape (text vs JSON, byte offsets); coverage beyond entities + the README. The `p2`/`vc` binary-simplification line is not pulled in unless explicitly added.
