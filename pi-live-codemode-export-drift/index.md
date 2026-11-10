---
id: cvrwyn3xynbwg4wkgc5khyfq
title: pi-live fails because pi-codemode no longer exports MCP_TYPESCRIPT_PREAMBLE
status: backlog
source: "Captain directive, 2026-11-10: the pi-live lane is red on main; pin-refresh already merged, so the repair is a new task."
started:
completed:
verdict:
score: 0.9
worktree:
issue:
pr:
mod-block:
---

The pi-live job of Runtime Live E2E has not passed since 2026-10-07. Worker launches die before a child session exists with a Pi package mismatch, so no commit on `main` can be greened on the Pi lane. This blocks the 0.28.0 release, which must green its exact release SHA.

## Problem

- Runs 42848114247 (`kent/pi-parity` @ `8d7307c99`) and 42849777496 (`spacedock-ensign/state-durability` @ `02536b6c3`), both 2026-11-10 and both now merged into `main` @ `5785dbed7`, fail the steps "Run live Pi common journeys" and "Run live Pi front-door smoke".
- Both logs carry, in the `recorded-gate-lifecycle` journey: ``Error: `@earendil-works/pi-codemode/declarations` does not export `MCP_TYPESCRIPT_PREAMBLE`.``
- The same error appears in failing runs 39608837699 (2026-10-21) and 41857395650 (2026-11-04).
- Last green pi-live: `main` @ `077fb3bc1` on 2026-10-06, and a PR branch on 2026-10-07. Every pi-live job since has failed or been cancelled.
- The Pi pins live in `internal/pilive/pilive.go` (`pi-coding-agent@1.0.2`, `pi-subagents@0.75.0`, `pi-intercom@0.16.0`), refreshed by `pi-live-lane-pin-refresh` (PRs #816, #830, #833). `pi-codemode` is not among them.

## Unknowns

- Which package versions disagree, and whether a declared pin, an unpinned transitive dependency, or the runner install step lets them drift.
- Whether the export error explains every failing journey and the front-door smoke. Failing in both runs: `smallest-sufficient-mechanism`, `recorded-gate-lifecycle`, `TestLivePiFrontDoorSmoke`. Failing in one run only: `gate-guardrail`; `filing` and `TestLiveCommonFullEnsignCycle`.

## Outcome

pi-live passes on a tip containing the fix, the install can no longer drift the same way without a deterministic check failing, and any failing journey with a different cause is named with its evidence rather than repaired here.
