---
title: Simplify the unreleased v1 gate-state schema
status: backlog
source: "Durable-decisions sprint implementation-shape audit, 2026-07-24."
score: "0.7"
id: jccbpvjv5bg1jn0jbmj2yf8s
---

The unreleased v1 gate-state implementation still carries prototype compatibility and a mutable current-gate pointer that duplicates derivable state and has already projected a stale approval.

Ideation should define the smallest clean-v1 schema:

- Remove `raw-file-pin` support and its compatibility fixtures; canonical bytes are the sole shipped binding.
- Determine whether `digest-domain` is redundant once the v1 schema fixes canonical digest semantics.
- Exercise multi-stage re-entry, multiple historical attempts, changed-Briefing supersession, and same-stage replay to determine whether `gates.current.gate` can be derived from current stage plus immutable records.
- If derivation is sufficient, remove the stored pointer and all pointer-repair/rebind behavior. If one counterexample requires stored selection, record that minimal counterexample and retain only the smallest non-stale selector.

The test must reproduce the observed failure class: an older approved attempt must not make a newly rejected candidate appear approved-awaiting-merge. No prototype migration or compatibility path is required.
