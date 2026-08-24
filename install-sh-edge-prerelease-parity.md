---
id: tdng3g6fe5p4y08tj5dphywc
title: install.sh resolves releases/latest and silently breaks edge version parity
status: backlog
source: "Captain CL fresh-install VM experience report, 2026-08-24: fresh VM tracking edge would have gotten v0.26.0 from install.sh with no warning"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:tdng3g6fe5p4y08tj5dphywc:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:tdng3g6fe5p4y08tj5dphywc-backlog-1
              briefing:
                id: briefing:tdng3g6fe5p4y08tj5dphywc:backlog:attempt-1:revision-1
                digest: sha256:d89298a1f4bcf4e18e6fcb6b6dd823249563f6367f773312dc641a5cf342c939
                request-digest: sha256:c7ed6226a494766f117a58b14240301160dcfa1c0c3826598321c35b70a7932d
                room-ref: ./install-sh-edge-prerelease-parity/review/backlog/briefing-1
---

`install.sh:90` resolves `api.github.com/repos/$REPO/releases/latest`, and GitHub's `/latest` excludes prereleases. Verified 2026-08-24: it resolves `v0.26.0` while the edge-parity target is 0.27.0-pre8 — a fresh machine following the binary-install path silently lands a binary one minor behind the edge plugin, which the FO boot version gate (requires minor 0.27) then aborts on. The skew is silent at install time: nothing prints the resolved version or hints that prereleases are excluded. This defeated the exact parity claim a fresh-VM prototype run was testing.

## Problem

{Ideation fills this in. Seeded: edge-channel machines have no scripted binary-install path that lands the current prerelease; install.sh's stable-only resolution is implicit, not stated or printed.}

## Proposed approach

{Ideation fills this in. Options to weigh: a channel/version opt-in on install.sh (e.g. env var or flag that lists releases including prereleases), printing the resolved version before install so the skew is visible, and/or documenting install.sh as stable-only with a distinct documented edge binary path. Related but distinct: the FO upgrade-hint journey (d2k) handles the wrong-version abort after the fact; this task is about not installing the wrong version in the first place.}

## Out of scope

Homebrew tap channel behavior. The FO boot upgrade-hint journey (d2k). Marketplace README repair (separate task).

## Expected surface and tolerance

Estimate net LOC change: ~+40, across 2 files (install.sh, its test/fixture surface).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded; ideation refines.

**AC-1 - A fresh machine tracking edge has a scripted install path that lands the current edge-parity binary version.**
Verified by: running that path and checking `spacedock --version` reports the prerelease, on a machine (or clean prefix) without a prior binary; fails if the path resolves a stable tag while a newer prerelease is the edge target.

**AC-2 - The stable path's resolution is explicit: install.sh states/prints what it resolved before installing.**
Verified by: script output naming the resolved tag; fails if a version skew can occur with no install-time signal.

## Test plan

{Shell-level test of release resolution (mock or fixture the API response — internal/release/install_url_test.go already exercises the endpoint), plus one manual/live fresh-prefix run for AC-1.}
