---
title: "codex --plugin-dir version-masquerade advisory leaks the internal branch name into shipped CLI output"
status: backlog
source: "0240 pre-cut antipattern audit (2026-07-01, lens 1). internal/cli/codex_marketplace.go installCodexLocalPluginDir prints a version-masquerade advisory on every --plugin-dir codex/pi install ending '...see next-post-release-preversion-bump' — an internal roadmap/branch identifier an end user cannot act on. The advisory itself is legitimate (version stamping is deferred); the dangling internal-branch pointer is the issue. Non-blocking; captain fast-tracked."
group: tooling
id: acy7gdv88md7jgzsfea85zx6
---

## Problem
`installCodexLocalPluginDir` (internal/cli/codex_marketplace.go, ~L132-142) prints a version-masquerade advisory to stderr on every `--plugin-dir` codex/pi install, ending with "...see next-post-release-preversion-bump" — an internal roadmap/branch name an end user cannot act on.

## Desired fix
Reword the advisory to drop the internal branch identifier. Keep the legitimate meaning ("the reported version reflects the checkout's checked-in manifest, not necessarily its current HEAD") but end it with something user-actionable — a public docs pointer or nothing — never an internal branch/roadmap name.

## Rough acceptance sketch (ideation tightens into measured ACs + test)
- The shipped advisory string no longer contains "next-post-release-preversion-bump" (or any internal branch/roadmap identifier), verified by a test asserting its absence; the user-facing meaning is preserved.
- The AC-3 advisory presence/absence test (internal/cli/codex_plugin_dir_test.go) still passes (present on --plugin-dir, absent on plain install), updated to the reworded string.
- go build ./... + go test ./internal/cli/ green.
