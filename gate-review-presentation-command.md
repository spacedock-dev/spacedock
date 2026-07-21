---
title: One-command gate review presentation with atomic result retention
status: ideation
source: "Split from the gate-recorder task (3k), captain-approved 2026-07-21. The subspace-coupled presentation half; 3k cycles 11-12 are its banked design history."
id: xbatj4hxtxw9t83vvmfem27f
sprint: durable-decisions
group: recorder
started: 2026-07-21T01:43:36Z
---

One binary command presents a gate: it validates an explicit briefing package (gate summary, frozen design snapshot, frozen probe input/history), derives the canonical title, launches the Subspace TUI as a blocking child through the caller's terminal transport, and atomically validates and retains the review log, resolution, and diagnostics on success or failure — the ensign presenter stays addressable until TUI exit plus retention, and pane creation or timeout is never completion. Includes the provider id-mapping adapter (the provider resolution binds its own envelope briefing id; normalize to the attempt briefing id after digest validation) as SPECIFIED in 3k's gate-resolution-frontmatter-contract.md. Scope moved from 3k at the split: ACs 7 and 15, the presentation-side AC-8 mutants (early-completion, detached-worker, live-Reference append, controller/child/validation/retention), and the probes companion (gate-review-probes.md rides as provider-owned convention; probes.jsonl proved out in one dry run). Evidence base: the 0260 shaping float findings 1-13 (blank-float EOF defect, launcher repair, probe-first ritual, retention-deleted-on-failure incidents) in the shaping debrief. Cross-repo sequencing: depends on the subspace-tui briefing-package and result surfaces; the working-copy-skill ritual recorded in the debrief is the interim. Land after the recorder (3k).
