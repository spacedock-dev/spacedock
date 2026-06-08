---
id: 78zrmrh8j3fv2b1t7bh77ykr
title: Clean up the Homebrew cask caveats — drop the xattr note, shrink the safehouse wall
status: validation
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

## Stage Report: implementation

- DONE: In .goreleaser.yaml homebrew_casks.caveats: DROP the notarization/quarantine explanation + the manual `xattr -dr …` fallback line (KEEP the hook), and REPLACE the multi-line safehouse block with a BRIEF recommend-install line for BOTH safehouse and agentsview. Leave homepage/description/license unchanged.
  Worktree commit caca3d51; caveats block (.goreleaser.yaml:75-78) now a 3-line companions note; homepage/description/license untouched.
- DONE: Keep hooks.post.install (the xattr auto-clear) intact.
  hooks.post.install (.goreleaser.yaml:79-87) still runs `xattr -dr com.apple.quarantine #{staged_path}/spacedock`; verified present post-edit.
- DONE: run `goreleaser check` to confirm the config parses.
  `goreleaser check` → "1 configuration file(s) validated", exit 0.
- DONE: find and update any Go test asserting the caveats content (e.g. in internal/release).
  grep over `*.go` for caveats/.goreleaser/agent-safehouse.dev found no test asserting caveats content — nothing to update (SKIPPED rationale embedded: no such test exists).
- DONE: `gofmt` + `go test ./...` green.
  `gofmt -l .` clean (only YAML changed); `go test ./...` → 1141 passed in 16 packages.
- DONE: commit on the worktree branch.
  Branch spacedock-ensign/brew-cask-message-cleanup, commit caca3d51.

### Summary

Trimmed `homebrew_casks.caveats` in `.goreleaser.yaml` from a 12-line block to a 3-line "recommended companions" note that briefly names safehouse (sandboxes agent runs, with its link) and agentsview (powers /spacedock:survey). Dropped the redundant notarization/quarantine explanation and the manual `xattr` fallback, since the `hooks.post.install` hook already auto-clears `com.apple.quarantine` — that hook is left intact. No Go test asserts caveats content, so none needed updating. `goreleaser check` exits 0; `go test ./...` is green (1141 passed).
