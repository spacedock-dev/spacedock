---
id: qcwrzza9xkr5kfwmdmvbqhkv
title: Pin the installed Pi package to the launcher release
status: backlog
source: "Released v0.27.2 Pi bootstrap abort reported by the captain on 2026-08-28"
started:
completed:
verdict:
score: "1.0"
worktree:
issue:
pr:
mod-block:
---

A stable launcher can install skills from a newer release and then fail its first binary gate.

## Problem

`spacedock install --host pi` uses the unpinned source `git:github.com/spacedock-dev/spacedock`.
Pi updates that source from the repository default branch. A stable v0.27.2 launcher therefore installed the v0.28.0-pre1 first-officer skill from `main`. The skill correctly rejected the older binary before `status --boot`.

The prior sibling-provider fix covered Claude and Codex only. Its approved scope explicitly excluded Pi.

## Proposed approach

Bind the default Pi package source to the running launcher's release identity: a release-stamped binary pins the source to its own release ref and never floats across tags. The unstamped development-build source behavior must be stated explicitly in ideation — float to the default branch, track `next`, or refuse — because an unpinned default silently re-creates this skew for development builds. State the pinning tradeoff as intended semantics: a pinned source does not auto-update, so stable users receive package fixes by upgrading the launcher, matching the Claude/Codex marketplace pin. Preserve `--plugin-dir` as the explicit development override.

Make the normal `spacedock pi` front door own one repair attempt. If the package is missing, unpinned, or incompatible, install the binary's pinned source. Recheck the package before launch. Refuse the launch if the repair fails or the package still does not match.

Exercise the ordinary installed Pi front door. Do not use `--plugin-dir`, `SPACEDOCK_REPO_ROOT`, a prose marker, or a direct skill path as the proof.

## Risk evidence

Pi's installed package manager marks a Git source with `@ref` as pinned. It fetches and checks out that ref. An unpinned source updates from its default branch. The released v0.27.2 source is unpinned, while its embedded first-officer contract requires minor 0.27. The current `main` contract requires minor 0.28.

## Out of scope

- Lowering the first-officer binary floor.
- Moving or replacing the published v0.27.2 tag.
- Claude or Codex provider inventory.
- A second package manager or a new Pi discovery protocol.
- Prose-grep tests.

## Expected surface and tolerance

Estimate net LOC change: +100, across 4 files. Tolerance: +/-60 net LOC and +/-2 files.

## Acceptance criteria

**AC-1 (VALUE) - A released launcher repairs and boots the Pi skill suite from the same release line.**
Verified by: an ordinary `spacedock pi` run starts with a missing or unpinned package, performs one pinned install, and reaches `status --boot`. The run uses no development override. The binary and loaded first-officer contract have the same major.minor. Removing the repair or its pin makes the run fail before workflow work.

**AC-2 - The default Pi install source is release-pinned, and the development override remains local.**
Verified by: command behavior tests observe the source passed to Pi for a release-stamped binary and for `--plugin-dir`. Changing the release source to omit its ref or changing the override to use the release source fails the assertions.

**AC-3 - The release fix is safe for v0.27.3 and the v0.28 prerelease line.**
Verified by: table tests cover stable, prerelease, and unstamped development identities. The test owns literal expected sources instead of deriving them from the production source function.

**AC-4 - A failed or ineffective repair cannot launch Pi.**
Verified by: command behavior tests require one install attempt, one recheck, no launch, and an actionable error after an install failure or remaining mismatch.

## Test plan

Use focused command behavior tests first. Run one existing Pi live journey through the installed-package front door without a development override. Start that run with no valid Spacedock package. Then run formatting, the full suite, race, and a detached adversarial audit because this changes a front door.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
