---
id: x0petxt7xvr459b6zh4vf4wj
title: doctor is blind to a co-installed enabled sibling spacedock plugin
status: backlog
source: "Split out of claude-install-sibling-channel-cleanup ideation (2026-08-25): making doctor see a sibling reaches gateHost's fail-fast branch — a front-door product decision that task could not own"
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
        - id: gate:x0petxt7xvr459b6zh4vf4wj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:x0petxt7xvr459b6zh4vf4wj-backlog-1
              briefing:
                id: briefing:x0petxt7xvr459b6zh4vf4wj:backlog:attempt-1:revision-1
                digest: sha256:6543c3f888c603117c5f6b46c3e2ba8add209f492fe8fecef71ce55349cbdf29
                room-ref: ./doctor-blind-to-sibling-dual-install/review/backlog/briefing-1
---

`spacedock doctor` resolves ONE manifest — the running binary's own channel via `hostOps.ResolveManifest(host)` — so a machine holding both `spacedock@spacedock` and `spacedock@spacedock-edge` enabled reports OK. After the sibling-cleanup install fix, the affected population is machines that already hold a dual install AND never run `spacedock install` again: a Compatible own-channel manifest never triggers the launcher auto-heal, so the condition does not self-clear on launch.

## Problem

{Ideation fills this in. Seeded: doctor needs, at minimum, a new hostOps capability to enumerate sibling installs (a seam every fake host in internal/cli implements), a verdict-or-Hint decision, an exit-code decision, and a gateHost decision — gateHost's default branch fails fast, so a sibling condition reaching it as non-Compatible would make the FRONT DOOR refuse to launch on a dual-install machine. That product/compatibility call is this task's core question. Related doc overstatement found in the same ideation: docs/site/reference/command-reference.md:42 claims doctor reports per-host "plugin versions and enablement" — today it reads one manifest's version and no enablement at all.}

## Proposed approach

{Ideation fills this in. Candidate shapes: (a) report-only Hint on dual install, exit 0, front door unaffected; (b) new non-fatal verdict class doctor prints but gateHost treats as Compatible; (c) full gate. The gateHost interaction decides.}

## Risk evidence

{Backlog: claude-install-sibling-channel-cleanup's ideation spike measured the dual-install state live (claude plugin list --json showing both channels enabled, doctor OK on the same machine) and mapped the code path (internal/contract/doctor.go single-manifest ManifestVerdict; internal/cli hostOps.ResolveManifest). Decides design should start.}

## Out of scope

The install-sequence sibling cleanup (claude-install-sibling-channel-cleanup ships it). Codex doctor behavior unless ideation proves the same blindness there.

## Expected surface and tolerance

Estimate net LOC change: +80, across 5 files. {Backlog seed; ideation refines with tolerance.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - On a host with both spacedock channel plugins installed and enabled, `spacedock doctor` surfaces the dual install instead of reporting only OK — without making the front door refuse to launch.**
Verified by: {ideation refines — seed: doctor fixture with both channel manifests present/enabled asserting the dual-install report appears and the exit code matches the chosen shape; a gateHost test asserting launch behavior on the same fixture. Falsifying edit: revert doctor to single-manifest resolution — the fixture must fail.}

**AC-2 - The command-reference doctor description matches what doctor actually reports.**
Verified by: {ideation refines — seed: the corrected docs/site/reference/command-reference.md:42 wording shipped with the capability, reviewed at validation; docs strict build as the mechanical check.}

## Test plan

{Ideation fills this in. Seeded: dual-manifest doctor fixture; gateHost launch test; no live lane expected unless the front-door path changes.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
