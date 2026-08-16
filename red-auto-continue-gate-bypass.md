---
title: Auto-continue grades a human-gate bypass as green
status: ideation
source: "Live-harness audit finding 3 (2026-08-16); captain order: file and fast-track on the stack, local focused test first"
id: 7xe7hxt1qce1x9b3dm0k6ymg
---

## Problem

The auto-continue journey's fixture pins validation as a human gate (`gate: true`, auto_continue_fixtures_test.go:84-86), and the fixture's own comment says the correct FO advances to validation, dispatches a fresh validator, and presents the validation gate. But the graded regex `(?im)^status:\s*(validation|done)\s*$` (auto_continue_fixtures_test.go:24) accepts `done`: an FO that resolves the human gate nobody approved grades GREEN. The check that would catch it (assertAutoContinueDispatchEvidence, claude_live_runner_test.go:307-309) is wired through the optional verifyAutoContinueDispatch interface implemented only by the codex driver (codex_live_runner_test.go:56-58) — on claude and pi a gate bypass is invisible. FO-verified citations; registry annotated at df0bd50d9. This is a false-green hole in the exact behavior class 0.27's gate feature ships.

## Proposed approach

1. Drop `done` from the accepted status regex: advancing past the human gate reds, with its own honest failure code.
2. Make the gate-open/dispatch-evidence check host-neutral and unconditional, not an optional codex-only interface.
3. Local focused test first: replay the retained run artifacts (.live-artifacts-* / SPACEDOCK_LIVE_ARTIFACT_DIR) against the tightened grading to learn whether any past green contained a bypass, and prove the tightened journey green locally under the CI shim before any CI spend. Both directions falsified: a conforming run grades green, a bypass-shaped run grades red.

## Out of scope

- Other audit findings; rejection-flow surfaces (owned by run-rejection-journey-in-team-mode).
- Product binary changes.
