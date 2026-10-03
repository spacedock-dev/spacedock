---
title: Pi doctor probes a pi-subagents source layout the package no longer ships
status: ideation
source: "Captain ran `spacedock doctor --host pi` against pi-subagents 0.74.0, 2026-10-01, and reported the two MISSING lines."
score:
started: 2026-10-03T02:42:26Z
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
              resolution:
                type: Resolution
                id: resolution:spacedock:mc0ajnpb4wh5nhd2p4q0vx9n:backlog:1
                briefing: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-10-03T02:42:14.267902Z"
                decision: approve
                reason: 'Captain directed dispatch of the seed: the defect, scope, and the ACs pinning the probe''s proof owner are stated.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:mc0ajnpb4wh5nhd2p4q0vx9n:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:mc0ajnpb4wh5nhd2p4q0vx9n-ideation-1
              briefing:
                id: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:ideation:attempt-1:revision-1
                digest: sha256:c8cb45aa0f2ab5d611cf05fcf867023d63b56105599a2cd42f42d2e901a685d3
                room-ref: '@review/ideation/briefing-1'
              withdrawal:
                by: agent:first-officer
                at: "2026-10-03T02:58:08.665278Z"
                reason: Bound artifact's AC-1 identifier wraps its bold span across two lines (lines 176-177), so the shipped status --read --ac-scan reader omits AC-1 entirely; proven by unwrapping it in a copy, where the scan returns AC-1, AC-2, AC-3. The gate AC cross-check cannot see the primary value AC. Repairing before presentation.
---

## Problem

`spacedock doctor --host pi` reports two prerequisites missing against a
correctly installed, up-to-date `pi-subagents`:

```
MISSING pi-subagents extension: .../pi-subagents/src/extension/index.ts
MISSING pi-subagents intercom bridge: .../pi-subagents/src/intercom/intercom-bridge.ts
```

Both probes test for TypeScript sources. The package ships compiled output.
The problem statement recorded pi-subagents 0.74.0 declaring
`"pi": { "extensions": ["./index.js"] }`; the installed 0.75.0 has the same
layout and pins the exact manifest resolution this task uses. Inspected on
2026-10-03 at `~/.pi/agent/npm/node_modules/pi-subagents` (version 0.75.0):

- `pi.extensions` = `["./index.js"]`; the root `index.js` exists.
- `src/extension/index.js` and `src/intercom/intercom-bridge.js` exist.
- No plain `.ts` source exists under `src/` (the only `.ts` files are `.d.ts`
  declarations); neither `src/extension/index.ts` nor
  `src/intercom/intercom-bridge.ts` exists.
- The intercom bridge is **not** declared in `pi.extensions`. It is reachable
  through the package `exports` map: `exports["./intercom-bridge"].default` =
  `./src/api/intercom-bridge.js`, which re-exports
  `./src/intercom/intercom-bridge.js`.

The printed remedy cannot clear the line it appears under: `pi install
npm:pi-subagents` reports `up to date`.

Four call sites hardcode the source paths in `internal/cli/pi.go`:
`piRuntimeConfigFromEnv` (`:617`), `checkPiRuntime` (`:644`, `:646`), and
`printPiDoctorReport` (`:773`, `:777`). The comment at `:303` repeats the
assumption ("package discovery loads `<pkg>/index.ts`").

## Risk evidence

The extension is live, not merely present: a session with pi-subagents loaded
carries the `subagent`, `subagent_supervisor`, and `bg_wait` tools. The defect
is in the probe, so the design question is narrow — how the probe should resolve
a package's entry instead of assuming one.

Independent baseline that can move the wrong way: with the current binary both
lines print `MISSING` on a machine where the package is correctly installed and
loadable. The fix must move those two lines to `OK` against that baseline.

## Approach

Resolve each probe target from the package's own `package.json` — the same
artifact pi reads — instead of hardcoding a source path:

- **Extension**: the first entry of `pi.extensions`, resolved relative to the
  package root. This is what pi's own package discovery loads.
- **Intercom bridge**: the package's `exports["./intercom-bridge"].default`
  subpath, resolved relative to the package root. This is the package's declared
  bridge entry.

One resolver, `resolvePiSubagentsEntries(packageRoot)`, reads
`<packageRoot>/package.json` once and returns both absolute paths (empty when
the manifest is missing or unparseable, so the probes fail closed).
`checkPiRuntime` stats the resolved paths; `piCheckResult` carries them so
`printPiDoctorReport` prints the resolved path; and `runPi`'s
unregistered-package fallback passes the resolved extension path to
`--extension` (it previously passed the nonexistent `.ts` path).

**Design decision (A), approved by the first officer and recorded here.** The
seed's AC-1/AC-2 originally implied both probes derive from one `pi.extensions`
entry. The real package contradicts that: pi-subagents declares only
`./index.js` in `pi.extensions` and does not declare the bridge there at all.
The bridge is therefore resolved from `exports["./intercom-bridge"]`. The
rejected alternative — a source-extension fallback at
`src/intercom/intercom-bridge.{ts,js,mjs}` — was refused because it
reintroduces the exact layout assumption this task removes and fails AC-3 (the
expected path would come from the probe's own guess, not the artifact).

Mechanism justification:

- The resolver is the end-value mechanism for AC-1/AC-2, not an enabling
  mechanism proved for its own internals.
- Simplest alternative considered: keep the hardcoded path and add a `.js`
  sibling candidate. Insufficient — it still assumes `src/extension/` and
  `src/intercom/`, so it re-breaks on the next package reorganization, which is
  precisely the defect.
- The `exports["./intercom-bridge"]` lookup is the minimal manifest-driven way
  to name the bridge; the package publishes no other declaration for it.

**No spike needed.** The only mechanism is standard `encoding/json` decoding of
the package's `package.json` plus `ops.Stat` on the resolved paths — both
already proven in this file (`readPackagePiSkills`, `isWorkspaceRoot`) and in
the live 0.75.0 artifact inspected above. No runtime handoff or on-disk format
is unverified.

## Out of scope

The doctor's other pi checks, the `--version` floor, the duplicate-registration
warning, and the CI lane's own `src/...index.ts` presence assertions
(`internal/release/workflow_exec_guard_test.go`, `.github/workflows/*`) — a
separate surface. This task changes only how the two pi-subagents probes resolve
their target paths.

## Expected surface and tolerance

Code surface: **2 files, net +60 to +90 LOC** (tolerance ±25%), reported as
insertions **~+115** and deletions **~-40** (the deletions cover the retired
hardcoded-path lines and the reshaped stat maps).

- `internal/cli/pi.go` — the `resolvePiSubagentsEntries` helper plus wiring at
  the four call sites, and the `:303` comment.
- `internal/cli/pi_frontdoor_test.go` — manifest-driven fixture plus reshaped
  stat maps.

Documentation surface (proposed diff below): **2 files, net 0** (four one-line
path edits).

Observable semantics this task **may** change:

- The two doctor lines' verdict for a correctly installed package (`MISSING` →
  `OK`) and the path each line prints — the resolved manifest entry
  (e.g. `.../pi-subagents/index.js` and
  `.../pi-subagents/src/api/intercom-bridge.js`).
- The `--extension` target `runPi` passes in its unregistered-package fallback
  (now the resolved `pi.extensions` entry).

Observable semantics this task **must not** change:

- The report's line labels (`pi-subagents extension`, `pi-subagents intercom
  bridge`), the `OK`/`MISSING` grammar, remedy text, exit codes, and every other
  doctor line.
- `pi.skills` discovery, the `--version` floor, the Spacedock package checks,
  and the front-door launch shapes.

## Acceptance criteria

**AC-1 - The two pi-subagents probes report OK for a package that declares a compiled entry.**
Verified by: a Go test that writes a fixture package from a declarative
manifest — a `package.json` declaring `pi.extensions` and
`exports["./intercom-bridge"]`, plus those declared entry files and no plain
`.ts` anywhere — then asserts both the extension and the intercom bridge lines
print `OK` and name the manifest-resolved paths. Independent baseline that can
move the wrong way: on today's binary the same fixture prints `MISSING` for both
lines.
Falsifying edit: restore the hardcoded `src/extension/index.ts` probe path —
the compiled-entry fixture must turn the test RED.

**AC-2 - Each probe reports MISSING when its own declared target is absent.**
Restated from the seed's single-file phrasing under approved decision (A):
pi-subagents declares the extension in `pi.extensions` but the bridge in
`exports`, so one shared entry cannot drive both lines. The restatement is
reported at the ideation gate and is proposed, not pre-approved.
Verified by: the same manifest-driven fixture with the declared `pi.extensions`
entry file removed → the extension line prints `MISSING` (the bridge line stays
`OK`); and with the `exports["./intercom-bridge"]` target removed (or its
`exports` entry deleted) → the bridge line prints `MISSING`. Each absence turns
its own line `MISSING`.
Falsifying edit: make a probe accept any existing file under the package root —
removing the declared entry must turn the test RED.

**AC-3 - The probe expectation does not come from the probe implementation.**
Verified by: the fixture reads back the `package.json` it wrote and derives the
expected paths from that manifest, rather than naming the paths the production
code checks in a `statOK` map.
Falsifying edit: delete the manifest write from the fixture (or the written
`package.json`) — the test must fail because the expectation cannot be derived
from the artifact.

## Test plan

Primary proof owner: `internal/cli/pi_frontdoor_test.go`. Its current
`statOKForPiResources` and `writePiSubagentsFixtures` fabricate
`src/extension/index.ts` and `src/intercom/intercom-bridge.ts` — the exact paths
the production code stats — so the tests pass by construction and cannot fail
when a real package changes layout. Reshape:

- `writePiSubagentsManifest(t, pkg, manifest)` writes `<pkg>/package.json` from
  a small struct (`pi.extensions`, `exports`) and creates the declared files.
- `readPiSubagentsManifest(t, pkg)` reads the written `package.json` back and
  returns the resolved expected paths — the expectation derives from the
  artifact, which is AC-3.
- `fakePiRuntimeOps.Stat` keys become the manifest-resolved paths
  (`<pkg>/index.js`, `<pkg>/src/api/intercom-bridge.js`), not the `.ts` paths.

Deterministic Go tests only — the claim is path resolution, not runtime
behavior; no live lane is needed.

Per-check falsifying edits (each turns its own check RED):

- AC-1 extension: hardcode `src/extension/index.ts` again.
- AC-1 bridge: drop the `exports` lookup, restoring
  `src/intercom/intercom-bridge.ts`.
- AC-2: accept any file under the package root.
- AC-3: remove the manifest write, or stop deriving expectations from
  `package.json`.

Cost: low — one new resolver plus fixture helpers, no runtime changes.

## Proposed documentation diff

The manual `--extension` examples name the stale path this task retires. Same
edit in two docs (four lines):

`docs/runtime-support.md:192`

```diff
---extension ~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+--extension ~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

`docs/runtime-support.md:247`

```diff
-~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

`docs/site/contributing/adding-a-runtime.md:113`

```diff
---extension ~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+--extension ~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

`docs/site/contributing/adding-a-runtime.md:194`

```diff
-~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

## Stage Report: ideation

- DONE: Task body states the problem, the chosen approach, and criteria under the exact `## Acceptance criteria` heading, with AC-1 through AC-3 each paired to a named falsifying edit that turns it RED.
  Task body has `## Problem`, `## Approach`, and `## Acceptance criteria` (scanner-exact heading); AC-1 falsifier = re-hardcode `src/extension/index.ts`, AC-2 = accept any file under the package root, AC-3 = remove the manifest write.
- DONE: Expected surface declares net LOC change and file count with tolerance, reports insertions and deletions separately, and names the observable semantics the task may change (doctor per-probe OK/MISSING lines) versus those it must not.
  `## Expected surface and tolerance`: 2 code files net +60..+90 (±25%), insertions ~+115 / deletions ~-40; may-change = the two doctor lines' verdict + printed path and the `--extension` fallback target; must-not-change = labels, OK/MISSING grammar, remedies, exit codes, other lines.
- DONE: Test plan names `internal/cli/pi_frontdoor_test.go` as primary proof owner, reshaped to a manifest-driven fixture whose expected paths derive from the written `package.json` rather than a stat map, and records the spike result (or "no spike needed" with the proven mechanisms).
  `## Test plan` names `internal/cli/pi_frontdoor_test.go` and the `writePiSubagentsManifest`/`readPiSubagentsManifest` reshape; `## Approach` records "No spike needed" with the proven `encoding/json` + `ops.Stat` mechanisms.

### Summary

Authored the ideation baseline for the two stale pi-subagents doctor probes.
Inspected the real installed package (pi-subagents 0.75.0) and found the seed's
AC-1/AC-2 implied one shared `pi.extensions` entry, but the bridge is declared
only via the package `exports` map; escalated and got decision (A) approved
(extension from `pi.extensions`, bridge from `exports["./intercom-bridge"]`).
Output: a fleshed task body with the resolver approach, a restated per-probe
AC-2 (reported at the gate as proposed, not pre-approved), and a manifest-driven
test plan.

### Same-stage repair: AC-1 bold-span line wrap

Repaired the one defect from the gate withdrawal: AC-1's bold identifier
wrapped across two lines (176-177 at authoring), so the line-based `--ac-scan`
reader could not see AC-1. Joined the AC-1 bold span onto a single line; no
other paragraph was reflowed.

Evidence — `spacedock status --read docs/dev/.spacedock-state/pi-doctor-probes-stale-subagents-layout.md --ac-scan --json --workflow-dir docs/dev`:

```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"180","unevidenced":"false"},{"id":"AC-2","line":"191","unevidenced":"false"},{"id":"AC-3","line":"204","unevidenced":"false"}]}
```

Checklist count summary:

- DONE (1/3): AC-1's bold identifier closes on a single line; `--ac-scan` now
  lists AC-1, AC-2, and AC-3 (previously AC-2 and AC-3 only).
- DONE (2/3): every other byte of the acceptance-criteria section, including the
  `## Acceptance criteria` heading and each AC's 'Verified by' and falsifying
  edit, is unchanged apart from that line wrap.
- DONE (3/3): the repair is committed path-scoped in the state checkout with no
  other file touched and no staged residue.
