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
gates:
    version: 1
    records:
        - id: gate:cvrwyn3xynbwg4wkgc5khyfq:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:cvrwyn3xynbwg4wkgc5khyfq-backlog-1
              briefing:
                id: briefing:cvrwyn3xynbwg4wkgc5khyfq:backlog:attempt-1:revision-1
                digest: sha256:2544a2af2281e604befe554972fd7e5132abba4caacdc03bb0991794ec9bc943
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:cvrwyn3xynbwg4wkgc5khyfq:backlog:1
                briefing: briefing:cvrwyn3xynbwg4wkgc5khyfq:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-10T18:37:03.84826Z"
                decision: approve
                reason: Seed names the diagnosed cause with local evidence, the unproven journeys, and a bounded pin-bump-plus-drift-check outcome; ready for ideation.
                conn:
                    quote: dispatch pi drift one too and approve toward PR stack on top of the haiku one.
                    source: captain chat, 2026-11-10
              application:
                target-stage: ideation
                state: pending
---

The pi-live job of Runtime Live E2E has not passed since 2026-10-07. Worker launches die before a child session exists with a Pi package mismatch, so no commit on `main` can be greened on the Pi lane. This blocks the 0.28.0 release, which must green its exact release SHA.

## Problem

- Runs 42848114247 (`kent/pi-parity` @ `8d7307c99`) and 42849777496 (`spacedock-ensign/state-durability` @ `02536b6c3`), both 2026-11-10 and both now merged into `main` @ `5785dbed7`, fail the steps "Run live Pi common journeys" and "Run live Pi front-door smoke".
- Both logs carry, in the `recorded-gate-lifecycle` journey: ``Error: `@earendil-works/pi-codemode/declarations` does not export `MCP_TYPESCRIPT_PREAMBLE`.``
- The same error appears in failing runs 39608837699 (2026-10-21) and 41857395650 (2026-11-04).
- Last green pi-live: `main` @ `077fb3bc1` on 2026-10-06, and a PR branch on 2026-10-07. Every pi-live job since has failed or been cancelled.
- The Pi pins live in `internal/pilive/pilive.go` (`pi-coding-agent@1.0.2`, `pi-subagents@0.75.0`, `pi-intercom@0.16.0`), refreshed by `pi-live-lane-pin-refresh` (PRs #816, #830, #833). `pi-codemode` is not among them.

## Unknowns

- Whether `recorded-gate-lifecycle` passes on Linux CI. In the local macOS run it failed with `successor dispatch build attempts = 1 (0 ok)`: `dispatch build --stamp` exited 2 with `--entity-path` under `/var/folders/…` and `--workflow-dir` under `/private/var/folders/…`, which points at a macOS symlink spelling rather than Pi, but that is unconfirmed.
- Whether `rejection-flow` passes: the local run hit the 50-minute suite timeout with it unfinished at 25 minutes.
- Which `pi-codemode` release first dropped the export (1.1.0 was published 2026-10-07, the day the lane went red).

## Diagnosis

`internal/pilive/pilive.go` pins `@earendil-works/pi-coding-agent@1.0.2`, which depends on `@earendil-works/pi-codemode` `^1.0.2`. The install has no lockfile, so npm resolves `pi-codemode` 1.5.0, whose `./declarations` no longer exports `MCP_TYPESCRIPT_PREAMBLE`. Pi 1.0.2's bundled `dist/extensions/codemode/tool.js` imports that symbol, so the extension fails to load. Pi 1.5.0 does not import it.

Local evidence, 2026-11-10 (macOS, model `openai-codex/gpt-6-luna`):

- Pi 1.0.2 installed the way CI installs it: `TestLivePiFrontDoorSmoke` fails in 102 s with the export error.
- Pi 1.5.0 with the pinned substrates (`pi-subagents` 0.75.0, `pi-intercom` 0.16.0): the smoke passes in 263 s.
- Pi 1.5.0 with `pi-subagents` 0.84.0 and `pi-intercom` 0.18.1: 16 of 18 tests pass, including every journey that failed in CI (`smoke`, `smallest-sufficient-mechanism`, `full-ensign-cycle`, `filing`, `gate-guardrail`). The two that did not pass are the unknowns above.
- In that run three `pi` xfail bindings passed: `owned-conflict-owner-handoff`, `default-headless-gate-stop`, and `auto-continue-after-implementation` split-root.

Latest published versions: Pi 1.5.2, `pi-subagents` 0.87.0, `pi-intercom` 0.18.1.

## Outcome

The pi-live pin moves to a current Pi family that loads cleanly, and the install can no longer drift through an unpinned transitive dependency without a deterministic check failing. pi-live passes in CI on a tip containing the fix. Any journey that still fails for a different cause is named with its evidence rather than repaired here.

The change ships as a pull request stacked on the `claude-live-forbid-schedulewakeup` branch, which is itself stacked on `ci/go-126-workflows` (PR #849).
