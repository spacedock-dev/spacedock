---
id: 9ezthwqxnff2tfchvg5z65sm
title: Make Debrief report the correct consumer project safely
status: backlog
source: Captain-authorized filing from the shipped-surface containment audit, 2026-08-06
started:
completed:
verdict:
score: 0.95
group: shipped-surface-containment
worktree:
issue: spacedock-dev/spacedock#533
pr:
mod-block:
---

Make Debrief use the consumer project's repository and resolved workflow state without unsafe shell interpolation.

## Problem

Debrief can derive links from the Spacedock plugin repository and ignore the resolved split-root state directory. It assumes fixed stages, `main`, numeric IDs, and GitHub squash text. Arbitrary issue and report text can enter shell arguments.

## Proposed approach

Resolve repository, state, trunk, stages, and ID style from authoritative project data. Transport untrusted text without shell evaluation. Generate links only from the resolved consumer repository.

## Out of scope

This task does not change the default filename owned by `8x`. It does not change split-root commit behavior or generic GitHub client behavior. It incorporates the same failure boundaries reported in issues #533 and #577 instead of creating duplicate repair tasks.

## Acceptance criteria

**AC-1 (VALUE) - A Debrief artifact identifies the correct consumer project and durable workflow state.**
Verified by: a temporary consumer repository with split-root state, custom stages, and SD-B32 IDs produces links and state references for that repository. Restoring plugin-repository or flat-state assumptions makes the check fail.

**AC-2 - Debrief supports the declared trunk and merge history forms.**
Verified by: fixtures with a non-`main` trunk and supported merge histories report the correct commit and PR data. Restoring a literal trunk or squash-only assumption makes a fixture fail.

**AC-3 - User-controlled text cannot alter a Debrief shell command.**
Verified by: shell-sensitive issue and report text remains exact data and causes no extra command. Reintroducing inline shell interpolation makes the negative control fail.

## Test plan

Use temporary consumer repositories with split-root state and custom workflow vocabulary. Add adversarial shell-sensitive text. Exercise the real Debrief entry point and compare the resulting artifact, repository links, process exit, and durable files.
