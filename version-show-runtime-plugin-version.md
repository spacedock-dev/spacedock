---
title: "Show the per-runtime plugin version in --version (drop enablement jargon)"
source: "captain (2026-06-13) — spacedock --version reports install/enablement posture but discards the installed plugin VERSION (which the runtime probe already returns) and labels an unreadable probe with the invented noun 'enablement'. Show the version pulled from the runtime, in plain words. Follows gj (#350, startup-sandbox-status)."
sprint: 0201-post-flip-release-model
group: ux-cleanup
sprint-readiness: ready
id: dag3bk4p0xe6tydc66k29ev3
status: implementation
worktree: .worktrees/spacedock-ensign-version-show-runtime-plugin-version
---

`spacedock --version`'s per-runtime block reports install + enablement posture but (1) discards the installed plugin VERSION — which the runtime probe already returns (`claude plugin list --json` carries a `version` field; `codex plugin list` prints a VERSION column) — and (2) labels an unreadable probe with the invented noun "enablement" ("enablement unknown"). The captain expects the plugin version pulled from the runtime, in plain words. Follows from gj (#350, startup-sandbox-status); fixes the version-display + jargon gap surfaced post-merge.

## Approach (locked)

Show the installed plugin version per runtime, version-forward, in plain words; delete the `enablement` type and all its jargon. Read the version ROBUSTLY from the resolved plugin manifest (the same source `doctor` reads via `ops.ResolveManifest`), NOT the fragile `plugin list` probe — so it renders even when `plugin list` errors (the captain's case). The enabled/disabled marker stays best-effort from the probe; if unreadable, omit it (bare version) rather than invent an "unknown" state.

Output format:

    claude: spacedock 0.20.0
    codex: spacedock 0.20.0 (disabled)
    pi: spacedock ready
    <host>: spacedock not installed     (host present, no plugin)
    <host>: not installed               (host binary absent)

The `Sandbox:` line is unchanged.

## Acceptance criteria

- **AC-1** — `spacedock --version` prints the version-forward per-runtime block; the noun "enablement" appears nowhere in cli output or the type vocabulary (grep-clean). Verified by command output + a tree grep.
- **AC-2** — the per-runtime version is sourced from the resolved manifest, so it still renders when `plugin list` errors. Verified by a fake-probe test where the enabled-probe fails but the manifest resolves and the version still shows (bare, no marker).
- **AC-3** — `version_runtime_test.go` rewritten to the new format; whole-repo `go test ./...` green; ZERO `.md` (yw rewrites the `## --version` doc to match).

## Test plan

TDD: rewrite `version_runtime_test.go` expectations FIRST (it locks the old "enablement"-style output), watch fail, then change `internal/cli/host_runtime.go` (runtimeStatus carries a version; probe extracts it from the resolved manifest; runtimeLine renders the new format) + the `cli.go` printVersion path. Keep the injectable runtimeProbe seam (no live host CLI in the test path). Confirm no caller depends on the deleted enablement constants.

gj-shaped: code + tests only; docs are yw's.
