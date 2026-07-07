---
title: "Riskiest-path-first binds the FO's own infrastructure changes, not just dispatched ideation"
status: backlog
source: "Principle-derived contract review, 0250 Commander session 2026-07-07 (captain-requested filing). Evidence: the edge-channel stable-cut gap (task zr) — the release flow's dev-preversion stamp composed with #468's minor-coupling gate was never exercised end-to-end before the first stable cut under the new gate (v0.24.0), hard-breaking every edge boot for days. The Probe and Ideation Discipline's spike rule binds dispatched-ideation designs only; nothing requires the smallest end-to-end exercise of the riskiest path for FO/operator-side changes to release machinery, CI lanes, or channels."
started:
completed:
verdict:
score: 0.3
worktree:
issue:
id: cyqk2gnw8mnjatap255259kx
---

Extend the spike discipline's reach: when a change touches the four high-stakes surfaces' infrastructure half (release machinery, CI lanes, channel plumbing), the riskiest composed path gets a minimal end-to-end exercise BEFORE the first real cut/run that depends on it — the integration-level kin of the existing ideation spike rule, applied to the process's own plumbing. Direction: a clause in the Proof policy (docs/dev/README.md, FO-owned process doc) plus, where cheap, a dry-run mode in the release flow (docs/releasing.md already models manual reconciliation escapes). Acceptance sketch: value — the zr incident class is non-reproducible under the rule (a stable-cut dry-run surfaces the skills-ahead-of-binary skew before any tag fires; baseline: 2026-07-04 shipped it live); mechanism — the clause + the dry-run check land together. Coordinate with zr (the concrete gap fix) — this is the rule, zr is the instance.
