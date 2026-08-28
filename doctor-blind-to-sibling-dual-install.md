---
id: x0petxt7xvr459b6zh4vf4wj
title: doctor is blind to a co-installed enabled sibling spacedock plugin
status: implementation
source: "Split out of claude-install-sibling-channel-cleanup ideation (2026-08-25): making doctor see a sibling reaches gateHost's fail-fast branch — a front-door product decision that task could not own"
started: 2026-08-28T06:48:16Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-doctor-blind-to-sibling-dual-install
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:x0petxt7xvr459b6zh4vf4wj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:x0petxt7xvr459b6zh4vf4wj-backlog-1
              briefing:
                id: briefing:x0petxt7xvr459b6zh4vf4wj:backlog:attempt-1:revision-1
                digest: sha256:6543c3f888c603117c5f6b46c3e2ba8add209f492fe8fecef71ce55349cbdf29
                room-ref: ./doctor-blind-to-sibling-dual-install/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:x0petxt7xvr459b6zh4vf4wj:backlog:1
                briefing: briefing:x0petxt7xvr459b6zh4vf4wj:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-28T06:47:47.308556Z"
                decision: approve
              application:
                target-stage: ideation
                state: consumed
        - id: gate:x0petxt7xvr459b6zh4vf4wj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:x0petxt7xvr459b6zh4vf4wj-ideation-1
              briefing:
                id: briefing:x0petxt7xvr459b6zh4vf4wj:ideation:attempt-1:revision-1
                digest: sha256:8691bd6bcfa6fe506664f012bd9b0006633b042f268b02c4cb77ae2974f08d93
                room-ref: ./doctor-blind-to-sibling-dual-install/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:x0petxt7xvr459b6zh4vf4wj:ideation:1
                briefing: briefing:x0petxt7xvr459b6zh4vf4wj:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-28T07:43:30.984809Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
---

`spacedock doctor` resolves ONE manifest — the running binary's own channel via `hostOps.ResolveManifest(host)` — so a machine holding both `spacedock@spacedock` and `spacedock@spacedock-edge` enabled reports OK. After the sibling-cleanup install fix, the affected population is machines that already hold a dual install AND never run `spacedock install` again: a Compatible own-channel manifest never triggers the launcher auto-heal, so the condition does not self-clear on launch.

## Problem

`spacedock doctor` checks one manifest. The binary channel selects that manifest:

- A stable binary checks `spacedock@spacedock`.
- An edge binary checks `spacedock@spacedock-edge`.

Claude and Codex resolve the `spacedock:first-officer` skill in a global skill
namespace. If both channel plugins are installed and enabled, the runtime can load a
different plugin from the one that doctor checked.

This fault produced a real bootstrap abort. A stable `0.27.1` binary checked a
compatible stable `0.27.x` manifest and doctor reported `OK`. Both runtimes then loaded the edge
`0.28.0-pre0` first-officer skill. That skill requires binary minor `0.28`, so Startup
stopped before `state.boot`.

The version floor is correct. Lowering it lets a `0.28` skill use commands that a
`0.27` binary does not own. The fault is the missing provider identity check between
doctor and the runtime loader.

The existing install sequences already remove sibling providers. However, a compatible
own-channel manifest does not start the launcher auto-heal. A machine with an old dual
install can therefore stay in this state until the operator runs install again.

## Proposed approach

Add a doctor-only plugin inventory to `hostOps`. The production implementation runs
`<host> plugin list --json` and normalizes these fields:

- plugin ID
- plugin version
- installed state
- enabled state

Claude and Codex use different JSON envelopes. Keep the normalized type inside
`internal/cli`. `internal/contract` continues to own only manifest compatibility.

After the current compatibility report, doctor finds the binary channel and its
sibling. If the sibling is installed and enabled, doctor prints this non-fatal report:

```text
CONFLICT: claude can load a different Spacedock plugin than doctor checked.
  checked: spacedock@spacedock 0.27.1 (installed, enabled)
  sibling: spacedock@spacedock-edge 0.28.0-pre0 (installed, enabled)
Run `spacedock install --host claude` to keep only the stable channel.
```

The selected plugin can be disabled. In that case, the `checked:` line reports
`(installed, disabled)`. An installed but disabled sibling does not produce the
conflict because it cannot supply the runtime skill.

The conflict keeps the compatibility exit code. A compatible dual install therefore
exits `0`. Doctor is a report, and `spacedock install --host <host>` is the repair.

If inventory fails, doctor prints this line and keeps the compatibility exit code:

```text
INCOMPLETE: doctor checked compatibility but did not read the claude plugin enablement state: {error}
```

This line prevents a silent `OK` when doctor cannot check the new condition.
`--plugin-manifest` remains a manifest-only override and does not inspect host state.

`gateHost` does not call the inventory method. A dual install remains recoverable, and
the front door continues to launch. A later `spacedock install --host <host>` removes
the sibling through the cleanup sequence that already shipped.

This is the smallest sufficient mechanism for AC-1. Reusing `ResolveManifest` is not
sufficient because that method discards sibling and enablement state. A concrete
`execHost` call outside the interface is not sufficient because tests could not supply
the failing state.

A new compatibility verdict is not necessary. It adds contract and `gateHost` paths
for a condition that does not change compatibility. The default `gateHost` branch can
also turn a new verdict into a launch refusal.

An automatic front-door repair is not necessary. The install command already repairs
the state, and a launch refusal would block access to a recoverable session. Lowering
the skill version floor is also not sufficient because it hides the wrong provider.

The doctor output serves AC-1. The inventory method supplies the missing evidence for
that output. The simpler single-manifest method cannot observe the end-value failure.

## Risk evidence

The riskiest path was exercised before this design. The proof used real host plugin
commands and real runtime skill loading on 2026-08-28.

**Claude, ambient dual install**

- `claude plugin list --json` reported `spacedock@spacedock 0.27.0-pre7+dev` and
  `spacedock@spacedock-edge 0.28.0-pre0`. Both entries were installed and enabled.
- A release-shaped stable binary reported `spacedock 0.27.1` and channel `stable`.
- Its `doctor --host claude` command exited `0`. It reported compatibility with the
  stable `0.27.0-pre7+dev` manifest.
- A bounded `claude --agent spacedock:first-officer` run loaded the skill from
  `/Users/clkao/.claude/plugins/cache/spacedock-edge/spacedock/0.28.0-pre0/skills/first-officer`.
- The loaded skill required binary minor `0.28`. It observed binary `0.27.1` and
  stopped at Startup step 1, before `state.boot`.

**Codex, isolated dual install**

- An isolated `CODEX_HOME` installed exact `v0.27.1` and current `0.28.0-pre0`
  plugin snapshots through real local Codex marketplaces.
- `codex plugin list --json` reported both channel IDs as `installed: true` and
  `enabled: true`, with versions `0.27.1` and `0.28.0-pre0`.
- Stable and edge doctor runs each exited `0` and reported only their own manifest as
  compatible.
- A bounded `codex exec` run selected
  `/private/tmp/spacedock-codex-dual-repro2/plugins/cache/spacedock-edge/spacedock/0.28.0-pre0/skills/first-officer/SKILL.md`.
- With `SPACEDOCK_BIN` set to the stable `0.27.1` binary, the skill rejected the
  binary because it requires minor `0.28`.
- The required `doctor --host codex` call then exited `0` and reported the stable
  `0.27.1` manifest as compatible. The runtime gate and doctor disagreed in one run.

The source explains the measurements. `execHost.ResolveManifest` selects only
`channelPluginID(devBranch)`. `ManifestVerdict` reads only that path. Neither function
reads a sibling ID or enabled state. Both runtime loaders register the duplicate
`spacedock:first-officer` name and selected the edge provider in these runs.

The Codex install sequence already enforces exclusivity when Spacedock runs it. Manual
host commands can still create the dual install, as the isolated proof shows. The new
doctor report therefore applies to both hosts.

## Out of scope

- Changes to the first-officer version floor.
- Changes to Claude or Codex skill-loader precedence.
- A new compatibility verdict or a front-door launch gate.
- Automatic plugin removal from `doctor` or the launch path.
- The existing install-sequence sibling cleanup.
- Local `spacedock@spacedock-local` and retired provider IDs.
- A disabled sibling, which cannot supply the runtime skill.
- Pi, which does not use the stable and edge plugin pair.
- Host inventory for `doctor --plugin-manifest`, which is an explicit manifest check.

## Expected surface and tolerance

Estimate net LOC change: +255, across 7 files. Expected insertions are 263 lines.
Expected deletions are 8 lines. Tolerance: net LOC +/-80 and file count +/-2.

Expected files:

- `internal/cli/frontdoor.go`: add the inventory capability and normalized entry type.
- `internal/cli/host_exec.go`: read and normalize Claude and Codex plugin JSON.
- `internal/cli/init.go`: add the doctor-only conflict and incomplete reports.
- `internal/cli/frontdoor_test.go`: extend `fakeHost` and prove non-gating behavior.
- `internal/cli/doctor_sibling_test.go`: add doctor output and exit-code behavior tests.
- `internal/cli/plugin_inventory_test.go`: add live-captured host JSON parser fixtures.
- `docs/site/reference/command-reference.md`: correct the doctor contract.

The two new test files can be combined with existing test files. This choice stays
inside the file-count tolerance.

Declared semantic changes:

- Command grammar: unchanged.
- Stored formats: unchanged.
- Authority: unchanged. Doctor remains read-only.
- Runtime behavior: default Claude and Codex doctor runs add a read-only plugin-list
  query. They print `CONFLICT` for an enabled sibling and `INCOMPLETE` on inventory
  failure. Both reports keep the current compatibility exit code.
- Front-door behavior: unchanged. `gateHost` continues to use only manifest
  compatibility.
- Healthy output: a single-channel compatible install keeps the current one-line `OK`.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - On a host with an enabled sibling channel, `spacedock doctor` cannot report only `OK`.**

Verified by: a table-driven command test for Claude and Codex. Each fixture has a
compatible selected manifest and an enabled sibling. The test requires the exact
`CONFLICT`, `checked:`, `sibling:`, and repair lines. It also requires exit `0`.
Removing the inventory call or either sibling entry makes this test fail.

**AC-2 - A valid stable `0.27.1` bundle still boots, and a dual install does not block the front door.**

Verified by: a front-door test with binary `0.27.1`, channel `main`, and a stable
`0.27.1` manifest. The single-channel fixture must reach `Launch` with no new output.
A second fixture supplies the dual inventory and must also reach `Launch` because
`gateHost` does not inspect it. Calling inventory from `gateHost` makes this test fail.
The bounded validation proof then loads the stable `0.27.1` skill and passes its
binary gate on Claude and Codex.

**AC-3 - Doctor reports host enablement from both supported JSON schemas and does not flag a disabled sibling.**

Verified by: parser and command tests with live-captured Claude and Codex JSON. The
fixtures contain both-enabled, selected-disabled, and sibling-disabled states.
Both-enabled and selected-disabled cases produce a conflict. The sibling-disabled case
keeps the exact one-line `OK` output. Ignoring `enabled` makes at least one case fail.

**AC-4 - The command reference states the selected-manifest limit and the non-fatal sibling report.**

Verified by: review of the exact diff below and `mkdocs build --strict`. Validation
also compares the text with the command fixture output from AC-1. This criterion serves
AC-1 by making the repair discoverable. Leaving the current overstatement is not
sufficient because it does not state the selected-manifest limit or the conflict.

## Test plan

1. Add failing command tests for AC-1 and AC-2 before the implementation.
   Cost: small, with no host process.
2. Add pure parser tests for the two live-captured JSON envelopes.
   Cost: small, with no host process.
3. Add the inventory method and the doctor report. Run the focused tests after each
   change.
4. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.
5. Run `mkdocs build --strict` for the command-reference change.
6. Repeat the bounded isolated host proof during validation.
   Install both plugin snapshots, run doctor, and require more than the `OK` line.
   Remove the sibling through `spacedock install`, then require the one-line `OK`.
7. Run a stable-only `0.27.1` bounded bootstrap on Claude and Codex.
   The runtime must load the stable skill and pass the binary gate.

The deterministic Go tests prove the product behavior. The bounded live runs prove the
host JSON and runtime-loader assumptions. No permanent live harness is necessary.

Because `internal/cli/frontdoor.go` is in the launch path, Claude and Codex live CI are
required before merge. Pi live CI is not required because the Pi doctor and launch
paths do not use this inventory.

### Documentation diff

Replace the current sentence in `docs/site/reference/command-reference.md`:

```diff
-For what is installed for each host — plugin versions and enablement — use `spacedock doctor`.
+Use `spacedock doctor --host <host>` to compare this binary with the selected channel plugin.
+If another channel is enabled, doctor names both plugins and reports the load conflict.
+This report does not block a launch. Run `spacedock install --host <host>` to keep one channel.
```

This is a reference page. The replacement states the visible contract and the repair.
It does not describe the parser or the host inventory.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: Reproduce the 0.27.1 bootstrap abort and trace the runtime-selected first-officer skill source versus the manifest and enablement state doctor reports on both Claude and Codex; identify the exact divergence rather than lowering the version floor.
  Claude and Codex both selected the enabled edge `0.28.0-pre0` skill while the stable `0.27.1` doctor checked only stable and reported `OK`.
- DONE: Choose the smallest release-safe behavior that prevents or loudly exposes a cross-channel dual-install mismatch without making a valid 0.27.1 bundle fail or making the front door refuse a merely recoverable state.
  The design adds a doctor-only `CONFLICT` report, keeps exit `0` for compatible bundles, and leaves `gateHost` unchanged.
- DONE: Define falsifiable ACs, fixtures or bounded live proof, exact files, LOC/file estimate and tolerance, including a valid 0.27.1 boot and a mismatched sibling-install case that cannot falsely report only OK.
  Four ACs bind exact output, front-door launch, both host JSON schemas, disabled-state behavior, documentation, and bounded valid-boot proof.

### Summary

The live proof found the same split-brain state on Claude and Codex. Doctor follows the
binary channel, while the runtime skill loader selected the enabled edge provider.
The proposed report exposes that difference without changing compatibility or launch
behavior.
