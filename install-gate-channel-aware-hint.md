---
id: p1tgy61tbhxj9apvpswqbhcy
title: The skill-first install hint is channel-blind - an edge plugin tells the user to install the stable binary its own version gate then rejects
status: ideation
source: "Captain CL in chat, 2026-08-24, reviewing install-sh-edge-prerelease-parity (#756): 'does it work for the skill-first journey that tells the user to run it?' - it does not; follow-up scoped to the skill hint path #756 could not reach"
started: 2026-08-24T19:20:28Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:p1tgy61tbhxj9apvpswqbhcy:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:p1tgy61tbhxj9apvpswqbhcy-backlog-1
              briefing:
                id: briefing:p1tgy61tbhxj9apvpswqbhcy:backlog:attempt-1:revision-1
                digest: sha256:20c0bb050069025cdf5f2b5627e1e28a9b2898b5ba07ead3e4dfc696107d6424
                request-digest: sha256:e389ddd2fbb0e491876aca632fed54018124271aefbb8b98726f2b7afae32f95
                room-ref: ./install-gate-channel-aware-hint/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:p1tgy61tbhxj9apvpswqbhcy:backlog:1
                briefing: briefing:p1tgy61tbhxj9apvpswqbhcy:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-24T19:20:03.982043Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''let''s dispatch p1t on to the channel install script PR stack'' - accepts the seed into ideation with stacked delivery on the #756 branch'
              application:
                target-stage: ideation
                state: consumed
---

install.sh now takes SPACEDOCK_CHANNEL (#756), but the skill-first journey never uses it. The FO binary-absent gate (first-officer-shared-core.md:10) hints the channel-less `curl ... | sh` on Linux and plain brew on macOS; fo-install-gate.md has zero channel awareness. On an edge-channel plugin install the journey self-defeats: the skills pin "binary minor 0.27" (prerelease-only until 0.27 goes stable), the hinted command installs stable v0.26.0, and the same gate then aborts on the binary it just told the user to install - the fresh-VM trap of 2026-08-24, surviving in the path that tells humans what to run.

## Problem

{Ideation fills this in. Seeded: the hint must know the channel with NO binary present. Candidate signals: the host's plugin-install record naming the marketplace (spacedock-edge vs spacedock); a channel stamp shipped in the skill bytes per channel ref; or the minor-pin-as-proxy (a pin satisfiable by no stable release implies edge). Also: macOS has no edge cask, so the brew hint is wrong for edge on that OS - the script path with SPACEDOCK_CHANNEL=edge is the only edge remedy on both platforms.}

## Proposed approach

{Ideation fills this in. Constraint worth honoring: the shared core line is boot-resident contract prose - keep the hint short; channel detection detail belongs in fo-install-gate.md, which already loads at the binary-absent trigger outside sandbox. Coordinate with fo-boot-upgrade-hint-latest-release (d2k), which owns the binary-present-wrong-version journey; this task owns binary-absent.}

## Out of scope

install.sh behavior (shipped in #756). The wrong-version upgrade journey (d2k). An edge Homebrew cask (its own release-machinery decision).

## Expected surface and tolerance

Estimate net LOC change: ~+25, across 2 files (first-officer-shared-core.md install line, fo-install-gate.md). Ideation refines with tolerance.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded; ideation refines.

**AC-1 (value) - A binary-absent boot on an edge-channel plugin ends with an edge-parity binary, not an abort loop.**
Verified by: a live or fixture-backed binary-absent journey on an edge plugin install where the followed hint lands a binary whose --version reports Channel: edge and a minor satisfying the skill pin; fails if the hinted command installs a stable binary the version gate then rejects.

**AC-2 - The hint the user sees names the channel-correct command for their OS.**
Verified by: gate output on an edge install carrying SPACEDOCK_CHANNEL=edge in the curl form (both OSes when brew has no edge cask); stable installs keep today's hints byte-identical. Fails if an edge user is sent to plain brew or a channel-less curl.

## Test plan

{Ideation fills this in. Seeded: the install-gate sentinel behavior journey (xw) is building the binary-absent test harness - coordinate rather than duplicate; a captive-PATH fixture with a stubbed plugin-install record may cover AC-1 offline.}
