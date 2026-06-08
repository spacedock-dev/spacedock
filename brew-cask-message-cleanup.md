---
id: 78zrmrh8j3fv2b1t7bh77ykr
title: Clean up the Homebrew cask caveats — drop the xattr note, shrink the safehouse wall
status: implementation
source: "captain (2026-06-08, flip-line) - the brew cask install message is too verbose. The postflight hook already auto-clears com.apple.quarantine, so the manual xattr fallback is redundant; the safehouse block is a wall of text. Just briefly recommend installing safehouse + agentsview."
started: 2026-06-08T03:54:12Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-brew-cask-message-cleanup
issue:
---

Clean up `homebrew_casks.caveats` in `.goreleaser.yaml` (the message `brew` prints on install). Verified this session: the cask's `hooks.post.install` already runs `xattr -dr com.apple.quarantine` on the staged binary, so the quarantine is auto-cleared — the manual fallback is redundant.

## Change (fully specified — dispatched straight to implementation)

- DROP from the caveats: the notarization/quarantine explanation AND the manual `xattr -dr …` fallback line. The postflight auto-clears it ("no need for xattr" in the user-facing text). KEEP the `hooks.post.install` xattr hook itself (that's the auto-clear mechanism).
- REPLACE the multi-line safehouse block with a BRIEF recommend-install line covering BOTH optional dependencies: **safehouse** (sandboxes agent runs) and **agentsview** (powers `/spacedock:survey`). A line or two, not a wall; keep the safehouse link if it stays short.
- Leave homepage / description / license unchanged.

## Out of scope

The README/install-journey docs (separate — already in PR #322); the safehouse-on-Linux question; the two-channel tap; notarization itself (5w).

## Acceptance criteria

**AC-1 — caveats are trimmed and recommend the two optional deps.**
Verified by: the `.goreleaser.yaml` caveats contain NO manual `xattr`/quarantine fallback text and NO multi-line safehouse block, and DO briefly name installing safehouse + agentsview — a content check over the config file (the deliverable, with values that can fail).

**AC-2 — the auto-clear is preserved and the config is valid.**
Verified by: `hooks.post.install` still runs `xattr -dr com.apple.quarantine`; `goreleaser check` exits 0 (config parses); any Go test asserting the caveats content (e.g. in `internal/release`) is updated; `go test ./...` green.

## Test plan

`goreleaser check` (config valid); a content assertion that the caveats are trimmed + name safehouse/agentsview; update any caveats-asserting test; `go test ./...`. Cheap, offline.
