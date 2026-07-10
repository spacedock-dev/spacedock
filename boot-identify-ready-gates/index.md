---
title: Expose ready-gate entities in boot identify JSON
status: ideation
score: 0.9
source: PR #493 local-live m3 investigation on 2026-07-10. The default Claude greeting omitted the already-gated Gate Check entity because status --boot --identify --json returned only dispatchable entities; dispatchAnalysis intentionally suppresses current gate stages, so the authoritative boot record made the shipped greet requirement impossible to satisfy.
id: 8n55etrw9wj10jfejdq5f1s8
---

## Problem

The interactive first-officer contract requires the greeting to name ready gates without rendering their full reviews. `status --boot --identify --json` exposes only `dispatchable`, which intentionally excludes entities whose current stage is a gate. In the real shallow-boot fixture, `gate-check.md` existed at `review`, but the boot record contained only `merged-pr`; the model therefore had no authoritative ready-gate entity to name.

## Proposed approach

Add a stable `ready_gates` array to the boot JSON. Populate it from active entities whose current workflow stage has `gate: true`, preserving deterministic entity ordering. Keep `dispatchable` semantics unchanged. Each ready-gate record must expose enough identity for a greeting to name it without reading entity bodies: `id`, `slug`, and `current` stage. Emit an empty array when none exist.

## Out of scope

- Rendering gate reviews during boot.
- Changing dispatchability, stage transitions, or gate ownership.
- Adding recursive entity reads or model-facing prompt cues.

## Acceptance criteria

- **AC-1:** Boot identify JSON lists every active current-gate entity in `ready_gates`, with stable `id`, `slug`, and `current` fields; the shallow-boot fixture exposes `gate-check` at `review`.
- **AC-2:** `dispatchable` remains byte-compatible and continues excluding current gate entities; terminal and unknown-stage entities do not enter `ready_gates`.
- **AC-3:** No-gate workflows emit `ready_gates: []`, and boot key-order/identify fixtures are updated deliberately.
- **AC-4:** The m3 default-startup live greeting can name `Gate Check` from the authoritative boot record without filesystem hunting or a custom prompt.

## Test plan

- Add focused boot JSON tests for one/multiple/zero ready gates, deterministic ordering, and dispatchable non-regression.
- Run affected status/CLI packages, full tests, race tests, and live-tag compile.
- Re-run m3's real Claude default-startup shallow-boot scenario on the coupled integration tree for AC-4.
