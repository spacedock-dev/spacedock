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

## Stage Report: ideation

Design round, no files written. Recorded here by the first officer; the round was instructed not to write
files, so its report was returned inline and its substance is preserved below.

- DONE: Chosen the page and its position.
  Evidence: one `## Pi` section in `docs/site/get-started/install.md`, between `## Launch` and `## Skills`,
  plus one line in `## Troubleshooting`; not a new page, because `docs/site/AGENTS.md` says to split pages by
  cadence rather than topic. `mkdocs.yml` is unchanged.
- DONE: Specified the reader's path in launch order.
  Evidence: `pi install npm:pi-subagents`, `pi install npm:pi-intercom`, then `spacedock install --host pi`
  (launch-required: `piRuntimeLaunchReady`, `internal/cli/pi.go:758`), then Pi auth, then
  `spacedock pi "/spacedock:survey"`.
- DONE: Mapped each doctor line to a remedy.
  Evidence: a seven-row table drawn from `internal/cli/pi.go:833-856`: pi CLI, Pi auth, the pi-subagents
  extension/skill, the intercom bridge, the pi-intercom package root/skill, the Spacedock package lines, and
  the pi version floor.
- DONE: Stated the prohibitions, each verified against source.
  Evidence: no version numbers (do not paste doctor output, which prints them); no model ids, which are owned
  by `internal/release/live_models.txt`; do not promise env-var-free resolution until the sibling launcher fix
  lands; do not claim auto-install parity with Claude and Codex; do not call Pi tier-1, since
  `docs/site/reference/command-reference.md:51` says Codex and Pi are experimental; do not present
  `--plugin-dir` as an install substitute; do not document unregistered Pi flags; keep runtime internals off a
  getting-started page.
- DONE: Identified the criteria needing rewording, with exact text.
  Evidence: AC-2, because half the path's commands belong to the Pi CLI and cannot be checked by
  `spacedock --help`; AC-3, because it banned version stamps but not model ids. Both replacements were
  applied, and AC-1 was tightened to name the two packages, `spacedock install --host pi` and
  `spacedock doctor --host pi`.

## Stage Report: implementation

- DONE: Load the ensign discipline before the assignment.
  Evidence: read `skills/ensign/SKILL.md` and `references/ensign-shared-core.md` + `pi-ensign-runtime.md`; no Pi env marker was set, so the Pi runtime adapter applies.
- DONE: Write the Pi getting-started path exactly as the ideation report designs it.
  Evidence: one `## Pi` section in `docs/site/get-started/install.md` between `## Launch` and `## Skills` (install.md:62), one `## Troubleshooting` line (install.md:170), no new page, `mkdocs.yml` untouched. Commit 5e8893294.
- DONE: Name every AC-1 element end to end.
  Evidence: the section names `spacedock pi`, `pi install npm:pi-subagents`, `pi install npm:pi-intercom`, the Pi auth step (`pi auth` + the `auth.json` doctor check), `spacedock install --host pi`, and `spacedock doctor --host pi`.
- DONE: Match documented commands and setup claims to the tools that ship.
  Evidence: built `./cmd/spacedock` and ran `pi --help`, `install --help`, `doctor --help`, `install --host pi --check` (exit 0), `doctor --host pi` (exit 0); compared `pi install npm:pi-subagents`/`npm:pi-intercom` against `pi install --help` and `pi auth` against `pi auth --help`; read `resolvePiSubagentsEntries`, `piRuntimeConfigFromEnv`, and `piRuntimeLaunchReady` in `internal/cli/pi.go`.
- DONE: State the registration-resolved resolution and keep the package-root escape hatch.
  Evidence: the section says `spacedock pi` resolves both packages from the entries `pi install` writes to Pi's settings (the sibling launcher fix's post-landing behaviour) and names `PI_SUBAGENTS_PACKAGE_ROOT` / `PI_INTERCOM_PACKAGE_ROOT` as the outside-the-agent-dir override.
- DONE: Obey the "must not claim" list.
  Evidence: no version number or doctor output pasted; no model id from `internal/release/live_models.txt`; no env-var-free promise beyond the registration read; no auto-install/tier-1 claim (the section says "Pi is an experimental host, alongside Codex", matching `reference/command-reference.md`); `--plugin-dir` not presented as an install substitute; no unregistered Pi flags; no runtime internals.
- DONE: Satisfy AC-3 with a search of the changed page.
  Evidence: `grep -nE "0\.28|0\.20\.2|claude-sonnet-5-5|claude-opus-5-5|gpt-6-luna" docs/site/get-started/install.md` → no match. Falsifier: adding `0.28` or any lane model id makes it match.
- DONE: Run `mkdocs build --strict` per AC-4.
  Evidence: installed `docs/requirements.txt` into `/tmp/mkdocs-venv` and ran `/tmp/mkdocs-venv/bin/mkdocs build --strict` → exit 0, "Documentation built". Falsifier: a broken internal link or nav entry makes the strict build exit non-zero.
- DONE: Confirm no Go regression and commit + push on the assigned branch.
  Evidence: `go test ./...` exit 0; worktree commit 5e8893294 pushed to `origin/spacedock-ensign/pi-getting-started-in-user-docs` only (never `main`).

### Summary

Added the missing Pi path to `docs/site/get-started/install.md`: a `## Pi` section
covering the two required Pi packages, Pi auth, `spacedock install --host pi`,
`spacedock pi`, the registration-based package resolution with the package-root
variables as the escape hatch, and `spacedock doctor --host pi`, plus one
troubleshooting line. No new page and no other file changed. Verified against the
binary's and Pi CLI's own help, kept clear of every item on the must-not-claim
list, and passed the strict docs build and the Go suite.
