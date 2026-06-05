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
- **The 0.19.6 line tested/green** — in particular gq (feedback-nonhappy-live-coverage) merged, and the pi-line settled. ("once tested" — define the gate: which lanes/scenarios green on `next` constitute "tested".)

## Direction (for ideation — CLARIFY the specifics with the captain)

Open questions to settle before/at ideation (flagged to the captain at filing):
- **What exactly does "flip the marketplace" mean?** Repoint the marketplace source/ref from the `next` branch to `main` (so `spacedock claude`/`codex` install pulls stable 0.20.0)? A version pin? Both Claude + Codex marketplaces? (Relates to `s0cq install-marketplace-ref-refresh` — the stale-ref bug — and `44 bundle-asset-distribution`.)
- **Is main currently a release branch, or does this establish main-as-release?** ("flip main readiness" implies main becomes the stable channel.)
- **Version:** 0.20.0 (a minor bump from the 0.19.x line) — confirm the bump rationale (contract change? feature line?).
- **The "tested" gate** — the explicit green-criteria on `next` that authorize the cut.

## Out of scope

The release-gate machinery fix itself (bqqr). The packaging/notarization tasks (44, 5w) unless they're prerequisites the captain names.

## Acceptance criteria

(To firm up at ideation once the specifics are clarified.)
**AC-1 — 0.20.0 is cut on main and the marketplace serves it.**
Verified by: the 0.20.0 tag/release exists on main (the cut succeeded through the fixed release-gate), and the marketplace install path resolves 0.20.0 from main (an install/resolve exercise confirming the flipped ref, not prose) — exact mechanics per the clarified specifics.

## Test plan

Per the clarified scope. At minimum: the cut runs green through the fixed `release.yml` (depends on bqqr), and a post-flip install resolves 0.20.0 from main. This is an outward-facing release — captain-gated at each outward step.
