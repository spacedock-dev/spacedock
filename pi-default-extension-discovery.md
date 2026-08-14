---
title: Pi default extension discovery for an isolated home
status: ideation
source: "Pi-UX carve, 2026-08-13: runtime-support first-contact friction"
score: 0.8
sprint: pi-ux
sprint-readiness: ready
group: tooling
id: 3w1ncf1thj12aryvkf5gj1rd
gates:
    version: 1
    records:
        - id: gate:3w1ncf1thj12aryvkf5gj1rd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3w1ncf1thj12aryvkf5gj1rd-backlog-1
              briefing:
                id: briefing:3w1ncf1thj12aryvkf5gj1rd:backlog:attempt-1:revision-1
                digest: sha256:ea7fae5fe635ebedf7504254e21a6e10ac5b467da3049466a2ae2ffaf9f47856
                request-digest: sha256:5bd2c1bebf60fb5a6261e2d6ec9e2f2b54564577d606af9f5d87079b59d884fd
                room-ref: ./pi-default-extension-discovery/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3w1ncf1thj12aryvkf5gj1rd:backlog:1
                briefing: briefing:3w1ncf1thj12aryvkf5gj1rd:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-14T06:43:06.589777Z"
                decision: approve
                reason: Captain approved backlog gate; advance to ideation for extension discovery.
              application:
                target-stage: ideation
                state: consumed
---

## Problem

An isolated Pi home does not auto-discover the `pi-subagents` / `pi-intercom`
extensions. The live harness works around this by hand-wiring paths: `piSubagentsPackageRoot`
requires `PI_SUBAGENTS_PACKAGE_ROOT` or falls back to `~/.pi/agent/npm/node_modules/pi-subagents`,
and the smoke launches pi with an explicit `--extension .../src/extension/index.ts`.
The operator's normal Pi install knows where these extensions live; an isolated
home forgets and has to be told.

`docs/runtime-support.md:147` names this as first-contact friction that is
supposed to be harness work ("an extension not auto-discovered in a temp home...
is harness work"), and the "assume it works" operating prompt expects auth and
package paths to be ironed out without a real blocker. Today the ironing is
manual and per-harness; a better default probe would make the isolated home
discover the operator's installed extensions the same way a normal home does.

## Visible value

A Pi runner in an isolated home resolves `pi-subagents` and `pi-intercom`
without the operator exporting `PI_SUBAGENTS_PACKAGE_ROOT` or the harness
hard-coding the `--extension` path. Measured against baseline: before, an
isolated-home Pi run with no `PI_SUBAGENTS_PACKAGE_ROOT` exported fails to find
the subagents extension; after, the same run discovers it from the operator's
installed package location and proceeds.

## Out of scope

- Changing where Pi itself stores or loads extensions.
- The `models.json` / `auth.json` copy (owned by `repair-pi-live-harness-parallelism-and-custom-model`, pnc).
- The intercom supervisor-talkback capability (archived spike `pi-intercom-runtime-capability-probe`).
- A new runtime, fixture, result format, or CI lane.

## Acceptance criteria

**AC-1 (VALUE) — An isolated-home Pi run finds the subagents extension with no env var exported.**

Verified by: an isolated-home Pi run that does NOT export `PI_SUBAGENTS_PACKAGE_ROOT`
resolves the `pi-subagents` extension (and `pi-intercom`, sibling package) from
the operator's installed package location, instead of erroring that the package
extension was not found. The baseline is the current `piSubagentsPackageRoot` fatal
path when the env var is unset and the fallback path is absent.

**AC-2 — The discovery is read from a real installed location, not a hard-coded fallback.**

Verified by: the probe reads the operator's actual installed extension/package
location (e.g. the Pi home's npm node_modules or package manifest), not a
hard-coded absolute path; a machine with extensions installed elsewhere is
discovered correctly.

**AC-3 — The explicit env var override still wins.**

Verified by: when `PI_SUBAGENTS_PACKAGE_ROOT` IS exported, it takes precedence
over discovery (existing harness behavior preserved); the `--extension` wiring
in the smoke stays valid.

**AC-4 — Offline and front-door smoke pass.**

Verified by: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
`go build -tags live ./internal/ensigncycle`, and `TestLivePiFrontDoorSmoke`
pass with no `PI_SUBAGENTS_PACKAGE_ROOT` exported.

## Test plan

Use the offline `PiLiveEnv|PiIntercom|TestPiLive` unit tests first. Then one
`TestLivePiFrontDoorSmoke` run with the env var unset only when Pi work is
authorized. Preserve the explicit-override path.

## Notes

- General first-contact-friction reduction, not a journey repair.
- Coordinate with `repair-pi-live-harness-parallelism-and-custom-model` (pnc),
  which owns the `models.json`/`auth.json` copy into the isolated home; both
  touch `seedPiLiveAuth` / isolated-home setup.
