---
title: Refresh the pi-live lane pins and substrate assertions for the Pi 1.0 family
status: ideation
score: 0.75
source: "Captain directive, 2026-10-03: Pi 1.0 shipped; update the CI pin and the relevant Pi extensions."
id: mh698y3ht6ydmr6ethaw9hg9
gates:
    version: 1
    records:
        - id: gate:mh698y3ht6ydmr6ethaw9hg9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mh698y3ht6ydmr6ethaw9hg9-backlog-1
              briefing:
                id: briefing:mh698y3ht6ydmr6ethaw9hg9:backlog:attempt-1:revision-1
                digest: sha256:6041994faf9dbfafbc683fc9a5b47729d9b34825fab9707546d87835726e2b53
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:mh698y3ht6ydmr6ethaw9hg9:backlog:1
                briefing: briefing:mh698y3ht6ydmr6ethaw9hg9:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T04:04:23.251551Z"
                decision: approve
                reason: Covers validating the published Pi family and deriving substrate assertions from the installed package.
                conn:
                    quote: i already said dispatch to ideation, but don't present the ideation gate until staff review finishes
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: ideation
                state: consumed
sprint: pi-ux
group: tooling
sprint-readiness: ready
---

The pi-live Runtime Live E2E lane pins an obsolete Pi family, so the shipped launcher
is validated against a runtime nobody installs. Refresh the pins and replace the two
hardcoded source-layout assertions that the pin bump breaks.

## Problem

`.github/workflows/runtime-live-e2e.yml` pins the Pi family at `PI_CODING_AGENT_VERSION=0.85.1`,
`PI_SUBAGENTS_VERSION=0.67.0`, `PI_INTERCOM_VERSION=0.13.0`, each with a sha512 integrity.
The current published releases are pi-coding-agent `1.0.0`, pi-subagents `0.75.0`, and
pi-intercom `0.16.0`. The lane therefore proves compatibility with a superseded Pi.

The lane also asserts the substrate ships TypeScript sources:

```
test -f "$pi_npm_root/node_modules/pi-subagents/src/extension/index.ts"
test -f "$pi_npm_root/node_modules/pi-subagents/src/intercom/intercom-bridge.ts"
```

Those paths no longer exist. pi-subagents `0.75.0` declares `pi.extensions` = `["./index.js"]`
and ships compiled output; no `.ts` file exists under `src/`. So the pin bump alone turns the
lane red, and the assertion encodes the same stale-layout assumption as the doctor probes
tracked in `pi-doctor-probes-stale-subagents-layout` (id `mc`).

## Risk evidence (measured locally, 2026-10-03)

Against pi `1.0.0` + pi-subagents `0.75.0` + pi-intercom `0.16.0`:

- `pi --version` prints `1.0.0`, and `spacedock doctor --host pi` reports
  `OK pi version: 1.0.0 (floor 0.83.0)` — the existing `piVersionFloor = "0.83.0"` accepts the
  major bump, so no floor change is forced by the bump itself.
- Every other pi doctor check is OK on that family: Pi CLI, Pi auth, the pi-subagents skill,
  the pi-intercom package root and skill, the Spacedock extension.
- The only failures are the two `.ts` probe lines, which are `mc`'s defect, not this task's.

So the mechanism this task rests on is exercised: pi 1.0.0 drives the existing launcher.
No spike beyond this recorded probe is needed.

## Expected surface and tolerance

Estimate net LOC change: +45, across 3 files (`.github/workflows/runtime-live-e2e.yml`,
a new or extended guard under `internal/release/`, and `docs/runtime-live-ci.md`).
Insertions ~+60, deletions ~-15. Tolerance: ±20 net LOC, ±1 file.

Declared semantic changes: the Pi family the pi-live lane installs, and the substrate
layout assertions the lane enforces. This task must NOT change the launcher's version-floor
behavior, the live scenario set, or the claude-live and codex-live pins.

## Acceptance criteria

**AC-1 — the pi-live lane installs the current published Pi family.**
Verified by: a Go structural guard under `internal/release/` that binds the workflow's
recorded pin variables to the current published versions, plus a green pi-live Runtime Live
E2E run on that pin set. Falsifying edit: revert `PI_CODING_AGENT_VERSION` to `0.85.1`; the
guard must turn RED.

**AC-2 — the substrate assertions are manifest-derived, not hardcoded source paths.**
Verified by: the lane resolves each asserted artifact from the installed package's own
`package.json` (the `pi.extensions` entry, and the `exports` entry for the intercom bridge)
instead of naming `src/**.ts`. Falsifying edit: restore the hardcoded
`src/extension/index.ts` assertion; the lane must fail against the currently pinned package.

**AC-3 — the integrity pins can fail.**
Verified by: the lane's `verified_pack` step rejects a corrupted integrity. Falsifying edit:
corrupt `PI_SUBAGENTS_INTEGRITY`; the step must fail rather than install.

**AC-4 — no stale pin text survives in the lane or its documentation.**
Verified by: the workflow comments and `docs/runtime-live-ci.md` name the pinned family
consistent with the variables. `docs/runtime-live-ci.md` currently records
`pi-coding-agent@0.80.10`, `pi-subagents@0.35.1`, and `pi-intercom@0.6.0`. Falsifying edit:
restore a stale version string in the doc; the check for text-variable agreement must fail.

**AC-5 (no-regression) —** `go test ./internal/release/...` and `go build ./...` are green,
and no other lane's action majors or pins change.

## Out of scope

The doctor's probe path resolution in `internal/cli/pi.go`, owned by
`pi-doctor-probes-stale-subagents-layout` (id `mc`). The harness node-runtime bump tracked by
`live-e2e-node-runtime-bump` (id `3g8`). The claude-live and codex-live CLI pins. The
`piVersionFloor` value, unless ideation shows the pinned family requires it.

## Test plan

Primary proof owner: a Go structural guard under `internal/release/`, following the existing
`node24_actions_guard_test.go` and `workflow_exec_guard_test.go` shape — per the dev proof
policy, `.github/**` release lanes are proven by Go unit tests for decision logic and
structural checks for YAML wiring, not by replay harnesses. The lane's own `verified_pack`
step is the integrity owner. A live pi-live run is required, because the claim is runtime
compatibility with the pinned family.

Deterministic Go tests plus one live lane run.
