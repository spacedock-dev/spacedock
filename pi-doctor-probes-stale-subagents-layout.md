---
title: Pi doctor probes a pi-subagents source layout the package no longer ships
status: backlog
source: "Captain ran `spacedock doctor --host pi` against pi-subagents 0.74.0, 2026-10-01, and reported the two MISSING lines."
score:
started:
completed:
verdict:
worktree:
issue:
id: mc0ajnpb4wh5nhd2p4q0vx9n
gates:
    version: 1
    records:
        - id: gate:mc0ajnpb4wh5nhd2p4q0vx9n:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mc0ajnpb4wh5nhd2p4q0vx9n-backlog-1
              briefing:
                id: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:backlog:attempt-1:revision-1
                digest: sha256:1238b011f4548835cdc914d79e2b7ec094be63f774d7ed21e6cea30f2c69ce8f
                room-ref: '@review/backlog/briefing-1'
---

## Problem

`spacedock doctor --host pi` reports two prerequisites missing against a
correctly installed, up-to-date `pi-subagents`:

```
MISSING pi-subagents extension: .../pi-subagents/src/extension/index.ts
MISSING pi-subagents intercom bridge: .../pi-subagents/src/intercom/intercom-bridge.ts
```

Both probes test for TypeScript sources. The package ships compiled output. In
pi-subagents 0.74.0:

- `package.json` declares `"pi": { "extensions": ["./index.js"] }`
- the root `index.js` exists and imports `./src/extension/index.js`
- `src/extension/index.js` and `src/intercom/intercom-bridge.js` exist
- no `.ts` file exists under `src/`

The printed remedy cannot clear the line it appears under: `pi install
npm:pi-subagents` reports `up to date`.

Four call sites hardcode the source paths in `internal/cli/pi.go`, and the
comment at `:303` states the same assumption ("package discovery loads
`<pkg>/index.ts`").

## Risk evidence

The extension is live, not merely present: a session with pi-subagents loaded
carries the `subagent`, `subagent_supervisor`, and `bg_wait` tools. The defect is
in the probe, so the design question is narrow — how the probe should resolve a
package's entry instead of assuming one.

## Out of scope

The doctor's other pi checks, the `--version` floor, and the duplicate
registration warning. This task changes only how the two pi-subagents probes
resolve their target paths.

## Expected surface and tolerance

Estimate net LOC change: +60 to +90 across 2 files (`internal/cli/pi.go` and
`internal/cli/pi_frontdoor_test.go`). Ideation refines this.

## Acceptance criteria

**AC-1 - The two pi-subagents probes report OK for a package that declares a compiled entry.**
Verified by: a Go test that builds a fixture package from a declarative manifest —
`package.json` declaring `pi.extensions`, plus the declared entry file, and no
`.ts` file anywhere — then asserts both the extension and the intercom bridge
lines print `OK`. Falsifying edit: restore the hardcoded
`src/extension/index.ts` probe path, which must turn the test RED.

**AC-2 - The same probes report MISSING when the declared entry is absent.**
Verified by: the same test with the declared entry file removed, asserting both
lines print `MISSING`. Falsifying edit: make the probe accept any existing file
under the package root, which must turn the test RED.

**AC-3 - The probe expectation does not come from the probe implementation.**
Verified by: the fixture builds its expected paths from the manifest it wrote,
not from a stat map naming the paths the production code checks. Falsifying edit:
delete the manifest from the fixture and the test must fail, proving the
expectation is derived from the artifact rather than assumed by the test.

## Test plan

The existing owner is `internal/cli/pi_frontdoor_test.go`, which currently
fabricates a `statOK` map naming `src/extension/index.ts` and
`src/intercom/intercom-bridge.ts`. Those entries are why the defect went
unnoticed: the test supplies the exact path the production code stats, so it
passes by construction and cannot fail when a real package changes layout. The
primary proof owner must be reshaped to a manifest-driven fixture, which is AC-3.

Deterministic Go tests only. No live lane is needed: the claim is about path
resolution, not runtime behavior.
