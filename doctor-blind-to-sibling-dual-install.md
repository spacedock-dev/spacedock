---
id: x0petxt7xvr459b6zh4vf4wj
title: doctor is blind to a co-installed enabled sibling spacedock plugin
status: validation
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

`spacedock doctor` resolves ONE manifest — the running binary's own channel via `hostOps.ResolveManifest(host)` — so a machine holding both `spacedock@spacedock` and `spacedock@spacedock-edge` enabled reports OK. The corrected behavior inventories providers on the launch critical path and reuses the existing one-shot install repair before a host can resolve the wrong first-officer skill.

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

Add a plugin inventory to `hostOps`. The production implementation runs
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

The conflict keeps the compatibility exit code. Doctor remains diagnostic-only and
never mutates the host.

If inventory fails, doctor prints this line and keeps the compatibility exit code:

```text
INCOMPLETE: doctor checked compatibility but did not read the claude plugin enablement state: {error}
```

This line prevents a silent `OK` when doctor cannot check the new condition.
`--plugin-manifest` remains a manifest-only override and does not inspect host state.

After a compatible manifest verdict, `resolveHealableGate` reads the same inventory.
With an enabled sibling, the default front door invokes its existing `Install` seam
once. That sequence already removes sibling providers and reinstalls the binary's
channel. The launcher then rechecks manifest compatibility and inventory exclusivity
before host launch. It never loops or adds a second repair implementation.

Under `--no-install`, the front door refuses before launch and prints
`spacedock install --host <host>` as the repair. If inventory cannot be read, the
critical path also refuses rather than launching without provider identity evidence.

This is the smallest sufficient mechanism for AC-1. `ResolveManifest` discards sibling
and enablement state, while the existing one-shot install path already owns safe
removal. A new compatibility verdict or a second uninstall sequence would add contract
surface without improving the guarantee. Lowering the skill version floor would hide
the wrong-provider fault.

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
- A new compatibility verdict.
- Automatic plugin removal from `doctor`.
- The existing install-sequence sibling cleanup.
- Local `spacedock@spacedock-local` and retired provider IDs.
- A disabled sibling, which cannot supply the runtime skill.
- Pi, which does not use the stable and edge plugin pair.
- Host inventory for `doctor --plugin-manifest`, which is an explicit manifest check.

## Expected surface and tolerance

Approved estimate remains +255 net LOC across 7 files, with net LOC +/-80 and file
count +/-2 tolerance. The release-package live proof must therefore fit at or below
+335 net LOC across 9 files; a worker cannot revise this baseline mid-stage.

Implemented files:

- `internal/cli/frontdoor.go`: add the inventory type and one-shot sibling heal/refusal.
- `internal/cli/host_exec.go`: read and normalize Claude and Codex plugin JSON.
- `internal/cli/init.go`: add the doctor-only conflict and incomplete reports.
- `internal/cli/frontdoor_test.go`: extend `fakeHost` and prove heal/refusal behavior.
- `internal/ensigncycle/live_test.go`: package the current plugin and stamp one shared stable binary.
- `internal/ensigncycle/claude_live_runner_test.go`: install that package before every common Claude journey.
- `internal/ensigncycle/codex_live_runner_test.go`: install that package before every common Codex journey.
- `internal/ensigncycle/team_capability_test.go`: keep common Codex launches on the ordinary front door.
- `docs/site/reference/command-reference.md`: correct the doctor contract.

Declared semantic changes:

- Command grammar: unchanged.
- Stored formats: unchanged.
- Authority: unchanged. Doctor remains read-only.
- Runtime behavior: doctor adds a read-only plugin-list query and remains diagnostic.
  Front doors inventory after compatibility; an enabled sibling is repaired through
  one install attempt, or refused under `--no-install`.
- Front-door behavior: a host launches only after the selected channel is compatible
  and no enabled sibling remains.
- Healthy output: a single-channel compatible install keeps the current one-line `OK`.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - An enabled sibling cannot reach runtime first-officer resolution.**

Verified by: deterministic launcher-seam tests provide a compatible selected manifest
plus an enabled sibling, require the existing `Install` path exactly once, then require
a clean re-gate, re-inventory, and launch. Every existing Claude/Codex common live
journey now installs the current release-stamped stable package and uses the ordinary
front door without bypass flags. A bad stable skill floor therefore fails before the
journey workload rather than relying on a marker-only special test.

**AC-2 - `--no-install` refuses a dual-enabled launch with the actionable repair.**

Verified by: both front-door command tests require exit `1`, no `Install`, no `Launch`,
and the exact host-correct `spacedock install --host <host>` command.

**AC-3 - Doctor reports the same conflict without repairing it.**

Verified by: Claude/Codex command fixtures require `CONFLICT`, checked/sibling identity,
exit `0`, and no install call. Live-schema parser fixtures preserve installed/enabled,
and a disabled sibling remains quiet.

**AC-4 - Healthy single-channel launches and the command reference remain correct.**

Verified by: existing front-door suites plus `go test ./...`, race, and strict docs.
The reference describes doctor diagnostics, launcher auto-repair, and manual repair.

## Test plan

1. Add pure parser, doctor, heal, and `--no-install` command tests.
2. Reuse `resolveHealableGate` and `Install`; re-gate and re-inventory once.
3. Package the current plugin with a binary stamped from its structured manifest
   version and `devBranch=main`; install it into each isolated common-live host home.
4. Remove `--plugin-dir` and `--skip-compat-check` from the common Claude/Codex paths.
5. Run formatting, focused tests, `go test ./...`, race, strict docs, and one existing
   common live journey per host. Do not add a bespoke dual-install journey.

Because `internal/cli/frontdoor.go` is in the launch path, Claude and Codex live CI are
required before merge. Pi live CI is not required because the Pi doctor and launch
paths do not use this inventory.

### Documentation diff

Replace the current sentence in `docs/site/reference/command-reference.md`:

```diff
-For what is installed for each host — plugin versions and enablement — use `spacedock doctor`.
+Use `spacedock doctor --host <host>` to compare this binary with the selected channel plugin.
+If another channel is enabled, doctor names both plugins and reports the load conflict.
+A normal launch repairs that conflict before starting the host; under `--no-install`, run the repair command manually.
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
The corrected implementation must expose that difference in doctor and remove it on
the launch critical path before runtime skill resolution.

## Stage Report: implementation

- DONE: Add a normalized Claude/Codex plugin inventory to doctor and the launcher.
  Doctor now emits exact non-mutating `CONFLICT` or `INCOMPLETE` reports after its
  compatibility result. An installed but disabled sibling keeps the healthy one-line
  `OK` result.
- DONE: Heal an enabled sibling through the existing one-shot install path.
  Both normal front doors now inventory after compatibility, invoke `Install` exactly
  once for a conflict, re-gate and re-inventory, and launch only after the enabled
  sibling is gone. `--no-install` and inventory failure refuse before launch with the
  host-correct repair.
- DONE: Put the critical bootstrap prerequisite in every existing common live journey.
  The shared live package reads the plugin version from its JSON manifest, stamps one
  stable-channel binary, installs the current plugin into each isolated Claude/Codex
  home, and launches through the ordinary front door without `--plugin-dir` or
  `--skip-compat-check`. No bespoke dual-install journey, marker workflow, transcript
  oracle, registry entry, remote tag clone, or prose grep remains.
- DONE: Verify AC-1 through AC-4.
  Focused launcher, doctor, inventory-failure, disabled-sibling, and real-schema tests
  pass. `go test ./...`, `go test ./... -race`, and a temporary-environment
  `mkdocs build --strict` pass. Existing Codex `TestLiveCommonShallowBoot` passes in
  33.09 seconds through the installed stable package. Claude reaches the same ordinary
  front door and reports exactly `spacedock@spacedock` v0.28.0-pre0 loaded from that
  package, then the local benchmark credential fails with an expired-token 401 before
  FO work; the package/bootstrap portion passed, but a green Claude model journey is
  externally blocked on credential renewal.
- DONE: Stay inside the approved implementation envelope.
  The final code candidate changes 9 files with 385 insertions and 83 deletions: +302
  net LOC versus +255, and +2 files versus 7. Production is +125 net LOC, tests are
  +177 net LOC, and documentation is 0 net LOC. This is inside the declared +/-80 LOC
  and +/-2 file tolerance. The deliverable is commit
  `62f0b103c25e1fd12e0af17c098311dfedda4400`.

### Summary

Doctor now exposes enabled sibling-channel conflicts without mutation, while both
front doors reuse the existing install sequence to restore exclusivity before launch.
Every common Claude/Codex live journey now exercises the current stable-stamped package
through the ordinary front door, so an incompatible stable skill floor fails at shared
bootstrap rather than in a special-purpose test.
