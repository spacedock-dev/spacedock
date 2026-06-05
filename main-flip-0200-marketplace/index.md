---
id: pj6k0h83g6taszkmt92qvsc0
title: Main-flip milestone — cut 0.20.0 on main + flip the marketplace once 0.19.6 is tested
status: backlog
source: "captain (2026-06-05) — '0.19.6: flip main readiness — once tested i want to make 0.20.0 on main and flip the marketplace.' The capstone of the 0.19.6 line."
score: "0.38"
started:
completed:
verdict:
worktree:
issue:
---

The 0.19.6 capstone: once the line is tested/green, cut **0.20.0 on main** and **flip the marketplace** so users get the stable release off main instead of the `next` development branch. This is the "flip main readiness" milestone.

## Hard prerequisites (gate this milestone)

- **`release-gate-job-separation-fix` (bqqr) MUST land first** — without it the cut fails the Runtime-Live-E2E gate (like 0.19.5).
- **The 0.19.6 line ALL GREEN** — captain's "tested" gate = all lanes/scenarios green on the branch (gq merged, the pi-line settled, all live + offline green). No partial-green cut.

## Direction (captain-clarified 2026-06-05)

- **The flip = `next` becomes `origin/main`.** Make `main` the release branch carrying what `next` holds; the marketplace serves from `main` (the stable channel) instead of `next`. This is the "flip main readiness."
- **Version = 0.20.0, the FIRST actual release.** NO contract change — the contract stays at 1; 0.20.0 is the first real (non-dev) release, not a contract bump.
- **Bundling is OUT — figure it out later.** The `44 bundle-asset-distribution` (plugin-into-binary, --plugin-dir) path is NOT part of this milestone; the captain explicitly deferred it. This milestone is the branch/marketplace flip + the cut, not the packaging mechanism.
- Note the related stale-ref bug `s0cq install-marketplace-ref-refresh` — a clean marketplace flip likely needs that fix so the new `main` ref actually replaces the old `next` pin (don't let the flip no-op on a stale ref).

## Out of scope

The release-gate machinery fix itself (bqqr). The packaging/notarization tasks (44, 5w) unless they're prerequisites the captain names.

## Acceptance criteria

(To firm up at ideation once the specifics are clarified.)
**AC-1 — 0.20.0 is cut on main and the marketplace serves it.**
Verified by: the 0.20.0 tag/release exists on main (the cut succeeded through the fixed release-gate), and the marketplace install path resolves 0.20.0 from main (an install/resolve exercise confirming the flipped ref, not prose) — exact mechanics per the clarified specifics.

## Test plan

Per the clarified scope. At minimum: the cut runs green through the fixed `release.yml` (depends on bqqr), and a post-flip install resolves 0.20.0 from main. This is an outward-facing release — captain-gated at each outward step.
