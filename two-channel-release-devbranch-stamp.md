---
id: k6d5xtg9hrxjcajrqyxnfah4
title: Two-channel release (stable→main / edge spacedock@next→next) + per-channel devBranch stamp + next-publish
status: ideation
source: "FO OWED, carried from the 2026-06-08-01 + 2026-06-08-02 debriefs (captain-nodded to file 2026-06-08). z9 (codex-plugin-auto-install) + #311 (Claude auto-install) install the plugin from the shared devBranch; the 0.20.0 flip needs each released channel's binary to install ITS OWN channel's plugin. Flip release-mechanics — a prerequisite of pj (main-flip-0200-marketplace), not a 0198 task."
started: 2026-06-08T22:48:35Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0200-flip
group: flip-mechanics
sprint-readiness:
---

Make each released spacedock channel auto-install the plugin from its own channel: the stable binary installs the `main` plugin, the edge binary installs the `next` plugin. Today there is one channel and `devBranch` is hardcoded `next`, so after the flip a stable binary would still install the `next` plugin.

## Problem

The front door auto-installs the plugin from a single shared `devBranch` (`frontdoor.go:49`, today `"next"`). Both `z9` (Codex auto-install) and the existing `#311` (Claude auto-install) consume this same `devBranch`. Two coupled gaps surface at the 0.20.0 flip:

1. **Channel ↔ plugin-source binding.** When the flip publishes a stable release on `main`, the stable binary must install the `main` plugin — but `devBranch` is still `next`, so it would install the edge plugin. The channel-tracking note in the 0198 sprint records this as the qa/z9 ↔ flip dependency: `z9` is correct *as long as it uses `devBranch`*; the retarget is what makes it install the right channel.
2. **Two distribution channels.** There is one brew artifact and one release lane today. The flip wants a stable channel (binary on `main`, installs `main` plugin) AND an edge channel (`spacedock@next`, installs `next` plugin), each stamped with its own `devBranch`, plus a `next`-publish step so edge keeps flowing after stable goes to `main`.

## Proposed approach

Ideation fills this in. Sketch:

1. **Per-channel `devBranch` stamp.** Stamp `devBranch` at build/release time per channel (stable→`main`, edge→`next`) instead of a hardcoded constant, so the auto-install path (`frontdoor.go` / `host_exec.go`) resolves the plugin source from the binary's own channel.
2. **Two-channel brew.** A stable cask/formula and an edge `spacedock@next`, each shipping its channel's binary.
3. **next-publish.** A release step that keeps publishing the edge channel from `next` after the stable line moves to `main`.

## Out of scope

- The 0.20.0 flip itself (`pj`) — this is its release-mechanics prerequisite, sequenced before the tag.
- `z9` / `#311` plugin-install behavior — they are correct against `devBranch`; this task only retargets/per-channels the value they read.
- Linux distribution (`v3`).

## Acceptance criteria

Ideation/implementation fills in. Sketch:

- A stable release stamps `devBranch=main` so the stable binary auto-installs the `main` plugin; an edge release stamps `devBranch=next` (verified by the built binary's resolved plugin source — a launch against a fresh HOME installs the channel-correct plugin, observed in workflow/plugin state, not by grepping the source constant).
- `brew install` of the stable channel and `spacedock@next` of the edge channel each deliver their channel's binary (verified by `--version` / install fixture per channel).
- The `next`-publish step keeps the edge channel current after the stable line moves to `main` (verified by the release pipeline producing both channels from a `--snapshot` or a release dry-run).

## Test plan

Ideation/implementation fills in. Per-channel stamp is the riskiest mechanism — exercise it first: build a binary with each `devBranch` value and confirm the auto-install path resolves the channel-correct plugin against a fresh HOME (live front-door smoke), not by string-matching the constant. goreleaser `--snapshot` for the two-channel artifacts; a brew/install fixture per channel.
