---
title: The user documentation has no Pi path
status: implementation
score: 0.7
source: "FO review, 2026-10-04: docs/site/get-started covers Claude and Codex; nothing covers spacedock on Pi."
id: 7gq0w79d76gcx2v2qnda5js6
worktree: .worktrees/spacedock-ensign-pi-getting-started-in-user-docs
---

## Problem

The getting-started documentation does not mention Pi. A search of
`docs/site/get-started/` finds no `spacedock pi`, no `--host pi`, and no
`pi-coding-agent`. So a Pi user has no documented path: which host command to run,
which two extensions are required, how to install them, how to authenticate, and
how to check that the setup is complete.

## Visible value

A Pi user follows the getting-started documentation from install to a first
workflow, without guessing a command or an extension name.

## Out of scope

- Version numbers. The family stamp lives in one place in the code. The
  documentation says how to check versions and never what they are.
- Rewriting the Claude or Codex sections.
- Runtime-support design notes, which stay in `docs/runtime-support.md`.

## Expected surface and tolerance

`docs/site/get-started/`: a Pi section in the install page, or one page beside it,
plus any index link. Estimate net +60, across 2 files. Tolerance net +20 to +120,
at most 4 files.

## Acceptance criteria

**AC-1 (VALUE) - The Pi path is documented end to end.**
Verified by: the page names `spacedock pi`, the two launch-required packages
`pi-subagents` and `pi-intercom` with their install commands, the Pi auth step,
`spacedock install --host pi`, and `spacedock doctor --host pi`.
Falsifier: remove any one of those and a reader cannot complete the path.

**AC-2 - The documented commands and setup claims match the tools that ship.**
Verified by: running each documented Spacedock command's help (`spacedock pi`,
`spacedock install`, `spacedock doctor`) and comparing subcommands and flags;
comparing the Pi-side commands against the Pi CLI's own help; and reading the
launcher and doctor source for the extension-resolution and readiness paths.
Falsifier: document a Spacedock flag the binary rejects, a Pi subcommand the Pi
CLI rejects, or a resolution or auto-install behaviour the shipped launcher does
not perform.

**AC-3 - No version number and no model id appears.**
Verified by: searching the changed pages for the family version string and for
every exact lane model id owned by `internal/release/live_models.txt`, requiring
no match. Falsifier: add either and the check must fail.

**AC-4 (no-regression) -** The documentation build passes and the existing pages
keep their links.

## Test plan

The documentation build, plus the command help comparisons. No live run.
