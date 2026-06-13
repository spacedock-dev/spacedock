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

## Stage Report: implementation

- DONE: Fix the per-runtime block to show the installed plugin VERSION per runtime and DROP the "enablement" type + jargon (plain words only); CODE + tests, ZERO .md.
  Real binary `--version` renders: `claude: spacedock 0.19.9` / `codex: spacedock 0.20.0` / `pi: spacedock ready` (commit 3c76472a).
- DONE: Build EXACTLY the locked output format (version-forward per-runtime line states).
  `runtimeLine` renders all five states; `TestVersionPerRuntimeBlock` + `TestVersionHostAbsentAndNoPlugin` assert the exact whole lines.
- DONE: VERSION SOURCE — read from the RESOLVED plugin manifest (the source doctor reads), not the fragile `plugin list` probe; marker stays best-effort and omits when unreadable.
  `probeVersion` → `execHost{}.ResolveManifest` → new `contract.ManifestVersion`; marker is a separate probe. `TestVersionMarkerUnknownRendersBareVersion` proves a failed marker probe still renders the bare version.
- DONE: DELETE the `enablement` type and its constants and the retired phrasings; the word "enablement" must not appear in output or type vocabulary.
  Grep-clean: no "enablement" in cli non-test code, no old constants/funcs anywhere in internal/. `TestVersionVocabularyHasNoEnablementJargon` asserts output is jargon-free.
- DONE: pi rendered in plain terms WITHOUT "enablement" and without inventing a version.
  Chose `pi: spacedock ready` / `pi: spacedock not installed`, sourced from the existing `piRuntimeLaunchReady` (skills + extension), via new `probePiReady`.
- DONE: TDD — rewrite `version_runtime_test.go` FIRST, watch it FAIL, then change the code; keep the injectable runtimeProbe seam.
  Red confirmed (undefined `enabledMarker`/`markerEnabled`/`claudeMarker`, struct-field errors), then green. Fake probe still pins state — no live host CLI in the test path.
- DONE: Confirm no other caller depends on the deleted enablement constants.
  Verified via grep: `claudeEnablement`/`probe*Enablement`/`codexEntryEnabled`/`enablement*` constants gone, no references remain.
- DONE: Whole-repo `go test ./...` green; zero .md.
  1317 passed in 16 packages (stable across two runs); diff touches only 5 .go files. `go vet` + `gofmt` clean.

### Summary

Reworked `internal/cli/host_runtime.go`: `runtimeStatus` now carries the resolved plugin `version`, a best-effort `enabledMarker` (markerUnknown/markerEnabled/markerDisabled), and pi `ready`; the `enablement` type and constants are deleted. The version is read robustly from the resolved manifest via a new exported `contract.ManifestVersion` (the same source `doctor` reads), decoupled from the best-effort enabled-marker probe — so the version renders even when `plugin list` errors (AC-2), and an unread marker omits rather than invents an "unknown" state. pi renders `spacedock ready` (no marketplace version). Verified on the real binary and by the rewritten tests. One pre-existing flake in `internal/ensigncycle` (`TestSonnetTeamDeleteHangReplay`, a stream-watch replay timing test) appeared once under the full parallel `./...` run but passed in isolation 3x, on the clean baseline, and on two subsequent full runs — unrelated to this cli/contract change.
