---
id: e6aaveste2tm0nsyqt407k55
title: spacedock status read helper — entity/markdown FM + section-heading offsets for surgical reads
status: backlog
source: "captain (2026-06-14) — 0.20.4 backbone. Reading whole entity bodies / stage reports / the README is a recurring FO + ensign token sink (this session: 280-line bodies, ~315-line README, 143KB CI logs, the Read-then-status--set staleness echo). A `spacedock status` helper that returns FM + a section-heading→offset map lets callers read the one section they need (Read offset/limit) instead of the whole file. Helps the README work (rzp) and other report-reading areas."
started:
completed:
verdict:
score: "0.40"
worktree:
issue:
sprint: 0204-structured-reads
---

A `spacedock status` helper that, given an entity (or any markdown file), returns its parsed frontmatter PLUS a map of section headings with line offsets — so a caller reads the one section it needs (`Read(file, offset, limit)`) instead of loading the whole file. The general capability behind 0.20.4's reading-cost reduction.

## Problem

- The FO and ensigns read whole files to reach one part: a 280-line entity to get the latest `## Stage Report`, a ~315-line README to check the `## Sprints` section, a stage report buried mid-file. Tokens scale with the whole file, not the part needed.
- On Claude Code, a `Read` followed by a `status --set` mutation of the same file triggers the staleness echo (the whole file re-emitted as cache-write tokens) — another whole-file tax.
- There is no structured way to ask "what sections does this doc have, and where do they start" — callers read the whole file or guess offsets.

## Proposed approach (seed — ideation designs the interface + output shape)

A read helper (likely under `spacedock status`; the exact surface is ideation's call) that takes an entity ref or a markdown path and emits:
- the parsed frontmatter, and
- an ordered list of section headings, each with: heading text, level, start line-offset, and the section's line range (to the next heading) — so a caller can `Read(offset, limit)` exactly that section.

Use cases to keep in view: read an entity's latest stage report by heading; read the dev README's `## Sprints` section without the other ~300 lines (helps rzp); read any report-shaped markdown structurally.

## Out of scope

{Ideation fills — e.g. `status` subcommand vs a new `read` command; JSON vs text output; whether to also emit byte offsets; non-markdown formats.}

## Relation

- **The dev README** — the FO slims it directly (FO realm over the workflow `README.md`; not a tracked task). This helper makes the slimmed README (and entity reports) cheap to read structurally; the two compose.

## Acceptance criteria

{Ideation fills. PROOF must be EXTERNAL: Go tests over the helper's REAL output — given a fixture markdown (FM + several sections), the helper returns the correct parsed FM and a heading→offset map whose offsets, fed to a real section read, return exactly that section's text. The independent oracle is the fixture's known structure (offsets locate the right bytes), never a string/regex match over instruction prose. Plus a live exercise of an FO/ensign reading a target section through the helper.}

## Test plan

{Ideation fills.}
