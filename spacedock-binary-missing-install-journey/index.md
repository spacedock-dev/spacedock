---
id: qak50v1d6pghavfc0ewg5hjd
title: Binary-missing journey — help the user install spacedock + show they can launch with `spacedock claude`
status: ideation
source: "captain (2026-06-05) — 0.19.6 readiness. When the spacedock binary is absent/non-executable, the user should get a helpful install-and-launch journey (a hint that helps them install the binary, then shows they can launch with `spacedock claude`), not a bare abort."
score: "0.32"
started: 2026-06-08T02:14:39Z
completed:
verdict:
worktree:
issue:
---

When `spacedock` is missing or non-executable, the user today gets a terse abort (the FO contract's startup contract-gate names an install hint, but that is FO-prose seen only mid-workflow). A first-time or returning user who runs into the missing-binary case should get a clear, friendly journey: detect the absence, help them install (the right channel for their platform), and show the payoff — that they can then launch with `spacedock claude`.

## Problem

The missing-binary path is an abort, not an on-ramp. The install hint lives in the FO contract (model-ingested prose) rather than in a user-facing surface the person actually hits when the binary isn't there. There is no guided "install → launch with `spacedock claude`" journey.

## Direction (for ideation)

- Decide WHERE the journey lives: the launcher front door / a `spacedock doctor`-style helper / the contract-gate abort message — whichever the user actually encounters when the binary is absent. (Note the FO startup contract-gate already emits install hints — `brew install spacedock-dev/homebrew-tap/spacedock` or `go build -o spacedock ./cmd/spacedock` — and is told NOT to run `spacedock doctor` when the binary is missing; reconcile with that.)
- The journey: detect absence → platform-appropriate install hint → confirm/verify → show the launch payoff (`spacedock claude` to start a Claude-hosted run). Keep it short and actionable.
- Relates to (do not duplicate): `44 bundle-asset-distribution` (zero-config plugin install), `5w notarize-macos-release` (Gatekeeper), `s0cq install-marketplace-ref-refresh`. This task is the USER-FACING missing-binary journey, distinct from the packaging/marketplace mechanics.

## Out of scope

The packaging itself (bundling, notarization, marketplace ref) — those are their own tasks. This is the detect-and-guide journey for a user who hits the missing/unusable binary.

## Acceptance criteria

**AC-1 — a missing/unusable binary produces a helpful install-and-launch journey, not a bare abort.**
Verified by: a behavior fixture/command that invokes the relevant surface with no usable `spacedock` on PATH and asserts the output (a) names the platform-appropriate install step and (b) shows the `spacedock claude` launch payoff — output/exit-code assertion, not prose review.

## Test plan

Behavior fixture driving the chosen surface with the binary absent (PATH-stripped / non-executable), asserting the emitted guidance. Cheap, offline. Ideation picks the surface + the exact copy (specific before/after wording).
