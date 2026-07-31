---
id: d2k1ww3ghw0y4f8d04e69s9c
title: "First-officer boot upgrade journey: check latest release and hint the user to upgrade"
status: backlog
source: "Split from fo-boot-install-hint-linux-direct-sandbox (GitHub issue spacedock-dev/spacedock#581) at the ideation subspace:r gate review 2026-07-31 — captain annotation: 'let's focus on install here. i want a separate upgrade journey where we can check latest release and hint the user for upgrade. file that separately.' Covers the #581 'binary present but wrong version' abort class."
started:
completed:
verdict:
score:
worktree:
issue: spacedock-dev/spacedock#581
---

When the FO boot version gate finds a binary present but wrong version (major.minor mismatch), today it aborts and points at `spacedock doctor`. This task builds a distinct upgrade journey: the FO checks the latest available release (from the version in the manifest / a release-query mechanism the previous design deferred), hints the user with the correct OS-aware upgrade command, and — where safe — offers to run the upgrade and re-check `--version`, converging to resume. Kept separate from the install journey (fo-boot-install-hint-linux-direct-sandbox) so the two journeys own their own convergence and test surfaces.

## Problem

{Ideation fills this in, seeded from the parent task's review: a wrong-minor binary aborts with no self-repair path; the operator copy-pastes a command and restarts. The upgrade path needs to know what "latest" is (currently only the manifest-skew comparison exists, not a newest-release query).}

## Proposed approach

{Ideation fills this in. Seeded: a release-check mechanism (e.g. GitHub releases query or the installed-version-vs-latest comparison), then an OS-aware upgrade hint (brew upgrade / curl|sh reinstall), then the optional upgrade-and-resume convergence the parent design carved out. The brew-upgrade convergence spike deferred from the parent task lives here.}

## Out of scope

{Ideation fills this in. Explicitly out: the binary-absent install journey (owned by fo-boot-install-hint-linux-direct-sandbox).}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded from the split; ideation refines and re-anchors to end value.

**AC-1 - When the boot gate finds a wrong-minor binary, the user receives a hint naming the latest available release and the correct OS-aware upgrade command, not a bare abort.**
Verified by: {ideation names the test — assert the wrong-minor path emits the latest-release identifier and the OS-appropriate upgrade command; name the concrete change that would make it fail.}

**AC-2 - The upgrade path, when the operator accepts the offer, converges: run upgrade, re-check `--version`, resume on compatible binary; one-attempt bound with hint-and-abort fallback.**
Verified by: {ideation names the test — behavior fixture driving the upgrade path with a captive upgrade; assert re-check and resume (or fallback on failure); the brew-upgrade convergence spike deferred from the parent task lands here.}

## Test plan

{Ideation fills this in, seeded from the parent review: the brew-upgrade convergence mechanism was asserted but never run (staff-review finding N1); a real brew-upgrade spike or an explicit reasoned-not-run caveat belongs here, plus the one-attempt fallback fixture.}
