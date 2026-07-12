---
id: eqn4ecmdy9d5a0meeqxwjwfa
title: Make survey work inside Claude Cowork
status: ideation
source: captain request
started: 2026-07-12T06:27:41Z
completed:
verdict:
score:
worktree:
issue:
---

Enable the `survey` skill to recognize Claude Cowork and guide the user through the network-access prerequisite for obtaining the Spacedock binary before continuing with Cowork-native session-history analysis.

## Problem

`survey` assumes a normal local Claude Code checkout with a Spacedock binary on `PATH`, local project history, and repo-shaped evidence. Claude Cowork has a different tool and filesystem environment. Without an explicit Cowork runtime check, `survey` can follow unavailable local-history paths or stop at a missing binary without telling the user how to enable the download.

The Cowork path must detect the environment positively. When the Spacedock binary is unavailable, it must pause and ask the user to enable network access so the correct release binary can be downloaded. It must not attempt to bypass the permission boundary or silently treat missing access as runtime impossibility.

## Proposed approach

Add a Claude Cowork environment probe at the start of `survey`. Route detected Cowork sessions to the host-native session inventory and transcript tools instead of local Claude Code history. Add an explicit user-facing approval prompt before any binary download: explain that network access must be enabled, why the binary is needed, and what will be downloaded. After approval, select the release asset from the detected OS and architecture, install it in a writable user location, verify `spacedock --version`, and resume the survey path.

Keep runtime proof fixture-backed where possible, with a focused live Cowork check for the actual environment and permission interaction if the host supports it.

## Out of scope

- General Claude Cowork workflow commissioning or git support.
- Changes to connected-folder deletion semantics.
- Persisting or publishing Cowork transcripts, session identifiers, account data, filesystem paths, credentials, or other personal/private information.
- Automatic network-policy changes or downloads without the user's explicit enablement.
- GitHub issue creation.

## Acceptance criteria

**AC-1 - Survey positively detects Claude Cowork and selects the Cowork runtime path.**
Verified by: a fixture-backed or live test that presents Cowork's host signals and proves `survey` chooses the Cowork adapter without probing local Claude Code project history.

**AC-2 - A missing Spacedock binary produces a clear network-access prompt before download.**
Verified by: a test that starts in detected Cowork with no usable binary and asserts the user is asked to enable network access, with no download attempted before affirmative permission.

**AC-3 - After permission, installation selects the correct release asset and verifies the binary.**
Verified by: architecture-table tests for supported Cowork Linux architectures plus an install-path test showing the binary lands in a writable user location and `spacedock --version` succeeds; denied or unavailable network access leaves the environment unchanged and reports the next action.

**AC-4 - Cowork survey uses host-native session evidence and clearly marks unavailable repo cross-checks.**
Verified by: a Cowork adapter test showing session inventory and transcript reads feed workstream clustering while git/PR/scaffold conclusions are labeled unverified when no repository evidence exists.

**AC-5 - Tests and diagnostics do not expose personal or private information.**
Verified by: sanitized fixtures and assertions showing logs, reports, and errors exclude real transcript content, session identifiers, account details, credentials, machine-specific paths, and tokens.

## Test plan

Add focused tests for Cowork detection, permission-denied and permission-granted download flows, OS/architecture asset selection, writable install-path verification, and Cowork session-tool routing. Use synthetic session metadata and transcript fixtures only. Add a bounded live Cowork smoke when the host exposes a stable test surface; record only exit status, redacted capability results, and durable non-sensitive state evidence.
