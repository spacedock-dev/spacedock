# 0204 — structured reads (0.20.4)

**Sprint:** the entities matching `sprint: 0204-structured-reads` — list current members with `spacedock status --workflow-dir docs/dev --where sprint=0204-structured-reads`. Membership and per-task state are the query, never enumerated or tracked in this doc.
**Theme:** make reading entities, reports, and the README cheap and structured.

## Goal (success criterion)

An FO or ensign reads an entity's frontmatter + one target section (e.g. the latest `## Stage Report`) without loading the whole file — via a `spacedock status` read helper that returns the FM plus a section-heading→offset map — and the dev README is slim enough that the whole-file fallback stays cheap. Proven by the helper's real output over fixtures (correct FM + heading offsets that actually locate the sections) and the README/template line-count + behavior checks, never a prose-grep.

## Why

This session's dogfood surfaced the cost: 280-line entity bodies and stage reports read whole, a ~315-line README loaded at boot, 143KB CI logs — plus the Read-then-Bash staleness echoes. Reading is a recurring FO/ensign token sink. Two levers: a binary helper that returns a document's structure (read only the section you need), and a slimmer README/templates so the whole-file fallback is cheaper.

## Definition of Done

0.20.4 ships when, merged to `main`:
- **The `status` read helper** returns an entity/markdown file's parsed frontmatter + a section-heading→line-offset map, so a caller can read a target section with `Read(offset, limit)` instead of the whole file. Proven by Go tests over real helper output (correct FM; heading offsets that locate the right section bytes), plus a live exercise of an FO/ensign reading a section through it.
- `v0.20.4` cut after the pre-cut antipattern audit is clean.

The dev README's own slim-down is **FO-direct upkeep** (the FO owns and edits the workflow `README.md`), not a tracked sprint member — done alongside the helper so the slimmed README stays cheap to read. The commission-template slim (shipped scaffolding) rides with `ey`.

## Out of scope

Defined as the sprint forms (ideation / the Commander): whether the helper extends `spacedock status` vs a new `read` surface; output shape (text vs JSON, byte offsets); coverage beyond entities + the README. The `p2`/`vc` binary-simplification line is not pulled in unless explicitly added.
