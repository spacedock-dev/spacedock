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
sprint: 019x-pre-flip-cleanups
group: release-hygiene
sprint-readiness: ready
mod-block:
pr: "#323"
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

## Stage Report: validation

- DONE: AC-1: Read .goreleaser.yaml homebrew_casks.caveats directly (not the stage report) and confirm it contains NO manual xattr/quarantine fallback text, NO multi-line safehouse block, and DOES briefly name installing safehouse + agentsview.
  Read worktree `.goreleaser.yaml:75-78` directly: caveats are a 3-line "Recommended companions" note naming safehouse (with https://agent-safehouse.dev) + agentsview — no xattr/quarantine text and no multi-line safehouse block (the only `xattr` is in `hooks.post.install` at :86, outside caveats).
- DONE: AC-2: confirm hooks.post.install still runs `xattr -dr com.apple.quarantine`, run `goreleaser check` (expect exit 0), run `go test ./...` (expect green), then emit a PASSED/REJECTED recommendation.
  `.goreleaser.yaml:84-87` `hooks.post.install` runs `system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/spacedock"]`; `goreleaser check` → "1 configuration file(s) validated", exit 0; `go test ./...` exit 0, 1141 passed, no FAIL lines; independent grep for caveats/quarantine/safehouse-asserting Go tests found none → nothing to update.

### Summary

PASSED. Independently confirmed against the actual config (not the report): the caveats are trimmed to a 3-line companions note for safehouse + agentsview with no xattr/quarantine fallback and no multi-line safehouse wall, while the `hooks.post.install` xattr auto-clear is preserved. `goreleaser check` exits 0 and `go test ./...` is green (1141 passed, exit 0). This is a cosmetic caveats-string edit, not a high-stakes surface, so no detached adversarial audit was warranted. Recommendation: PASSED.
