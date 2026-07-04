---
id: k74gt0qv3j4b86knvy2rhsta
title: "Split Startup step 8 — interactive greet names ready gates; assembly moves to a per-entity Engage turn"
status: backlog
source: "FO session 2026-07-04: interactive boot ran ~5 minutes. Startup step 8's interactive path (first-officer-shared-core.md:31,37) directs the greet to render a full present-gate review (stage-report read + AC cross-check + any re-verification) for every ready gate:true stage before stopping, so greet cost scales with the number of ready gates/orphans instead of staying a flat manifest. Drafted by the session's science-officer teammate from concrete friction in that boot (inline full-suite re-verification, a red-herring unmerged-branch chase, two missed --stage flags, an unfiltered ~50-row status dump)."
started:
completed:
verdict:
score: 0.35
worktree:
issue:
---

Boot and Engage are fused into one turn: the interactive greet both composes a cheap discovery manifest AND eagerly assembles/renders every ready gate review (and could pull in orphan investigation), so a boot with K ready gates costs O(K) full gate-review reads instead of O(1). Proposed direction: split the interactive path into a manifest-only Boot stop (name dispatchable entities, ready-gate entities, and orphans — no bodies read) and a separate, captain-triggered per-entity Engage step that does the gate assembly/AC cross-check/dispatch. Headless stays untouched (it must author the full gate review at its one stop, since no Engage turn follows). Needs ideation to resolve: whether a single ready gate still auto-renders vs. always name-only, and whether a live drive can deterministically count pre-greet tool calls as the value proof.
