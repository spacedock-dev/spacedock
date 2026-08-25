---
id: 7p4vxcb9edzrph0zsfn75cd9
title: claude install leaves the sibling edge plugin installed and enabled
status: backlog
source: "Captain CL work-machine report 2026-08-25: clean documented stable install on a machine with spacedock@spacedock-edge left both plugins installed and enabled; doctor OK; which plugin serves unpredictable"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:7p4vxcb9edzrph0zsfn75cd9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:7p4vxcb9edzrph0zsfn75cd9-backlog-1
              briefing:
                id: briefing:7p4vxcb9edzrph0zsfn75cd9:backlog:attempt-1:revision-1
                digest: sha256:4b88d6e7f37061f48b8d0475ee4f5a15476c1728bf602afc4d29856f86b71fb9
                request-digest: sha256:7652910a297d84b88917d9d48ad7e4009f83582811873821de42860027c6a649
                room-ref: ./claude-install-sibling-channel-cleanup/review/backlog/briefing-1
---

`spacedock install --host claude` never uninstalls the sibling channel's modern plugin id. A stable install on a machine that already has `spacedock@spacedock-edge` (or the reverse) leaves BOTH channel plugins installed and enabled. Both are entry `spacedock` shipping the same skill set, so which provider serves is unpredictable. Mirror the codex sequence's sibling-channel remove, and give doctor eyes on a co-installed enabled sibling.

## Problem

`installArgvSequence` (internal/cli/host_exec.go:406) uninstalls only the own-channel id and the retired route-A id (`spacedock-edge@spacedock`). The modern sibling id is never touched. `codexInstallArgvSequence` (host_exec.go:445) already removes the sibling channel first — justified by codex's global skill namespace — and claude has the same practical collision: both channel plugins are entry `spacedock` with identical skill names. Neither safety net catches the dual install: `spacedock doctor` resolves only the running binary's own channel manifest (channel-aware resolver) and exits 0; the binary gate passes because both plugin snapshots can declare the same required binary minor. {Ideation refines.}

## Proposed approach

{Ideation fills this in. Seeded: add a tolerated `plugin uninstall spacedock@<otherChannelMarketplace(devBranch)>` step to installArgvSequence, mirroring codexInstallArgvSequence's sibling remove; decide whether doctor should flag a co-installed enabled sibling spacedock plugin in the same task or as a routed follow-up.}

## Risk evidence

{Backlog: captain live repro on v0.27.0 — after the documented stable install steps on a machine with an existing edge plugin, `claude` listed spacedock@spacedock v0.27.0 AND spacedock@spacedock-edge v0.27.0-pre8, both enabled; workaround was `claude plugin uninstall spacedock@spacedock-edge`. Source check 2026-08-25: v0.27.0 and main installArgvSequence are byte-identical — no sibling uninstall step exists. This decides design should start.}

## Out of scope

CLAUDE_CODE_PLUGIN_PREFER_HTTPS documentation (plugin-install-https-keyless-machines / rb0). Product README edge-channel mention (product-readme-edge-channel-mention). The codex install sequence (already removes the sibling).

## Expected surface and tolerance

Estimate net LOC change: +40, across 4 files. {Backlog seed; ideation refines with tolerance.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - After `spacedock install --host claude` on a host with the sibling channel plugin installed and enabled, exactly one spacedock channel plugin remains installed and enabled — the selected channel's.**
Verified by: {ideation refines — seed: a sequence-shape test in internal/cli asserting the claude install steps include the sibling uninstall argv as an independent literal (mirror channel_selection_test.go's literal-id style, not derived from channelPluginID), plus the existing tolerance asymmetry preserved. Falsifying edit: drop the sibling step or misspell it back to the retired route-A id — the test must fail.}

**AC-2 - `spacedock doctor --host claude` surfaces a co-installed enabled sibling spacedock plugin instead of reporting OK.**
Verified by: {ideation refines or splits out — seed: doctor fixture with both channel manifests present and enabled; falsifying edit: doctor keeps resolving only the own-channel manifest and stays silent — the fixture must fail.}

## Test plan

{Ideation fills this in. Seeded: unit test over installArgvSequence literals; doctor dual-manifest fixture; decide whether the claude-live lane is required (host_exec.go feeds the launcher/install path — check the path-to-lane mapping in the Proof policy).}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
