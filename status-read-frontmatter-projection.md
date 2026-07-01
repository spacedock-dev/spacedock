---
title: "status --read --json --frontmatter — boot-lean projection (frontmatter + stages only, drop the headings index)"
status: backlog
source: "FO boot-occupancy analysis 2026-07-01 (Pi FO Stage B report + Claude Commander session). status --read README --json returns the full headings index (18 entries for docs/dev/README.md), but the FO boot (Startup step 4) consumes only the frontmatter + the stages array (+ entity-label/plural + mission + id-style). The headings index is carried into boot context and unused at greet — avoidable Stage-B occupancy. Stage A (the FO contract load) is the bigger prize (nt/deferrals target it); this is the smaller, separable status-tooling half."
group: tooling
id: fkejdqg2aw5t777jajpqhz34
---

## Problem
FO startup reads `status --read {workflow_dir}/README.md --json` for the stage taxonomy (Startup step 4). Today that returns the frontmatter, the full `stages` array, AND a complete headings index (every `##`/`###` — 18 entries for `docs/dev/README.md`). Boot consumes only frontmatter + stages array + entity-label/mission + id-style; the headings index is unused at greet yet occupies boot context on every launch.

## Desired direction (for ideation to refine)
A boot-lean projection — e.g. a `--frontmatter` flag (or a `--boot`-shaped narrowing) on `status --read --json` returning only the frontmatter + the `stages` array (the fields Startup step 4 actually consumes), dropping the headings index. The full-headings read stays available for the phases that need it (a dispatch copies a stage subsection via `show-stage-def`; the merge ceremony reads `merge:` policy).

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- A MEASURED before/after: the byte/token size of the boot README read drops by the headings-index payload, verified on `docs/dev/README.md` (18 headings) — a real number, not a prose claim.
- The FO boot still resolves the full stage taxonomy (names/ordering + initial/terminal/gate/worktree/feedback-to/agent flags) + entity-label/plural + mission + id-style from the projection.
- A golden/JSON test pins the projected shape; the full-headings `--read --json` output stays byte-identical (negative control).

## Related
- Sibling boot-ergonomics item: `state` subcommand workflow auto-discovery.
- The 0240 lean-contract sprint (Stage A / FO-contract occupancy is the bigger, separate lever).
