---
title: "Show the per-runtime plugin version in --version (drop enablement jargon)"
source: "captain (2026-06-13) — spacedock --version reports install/enablement posture but discards the installed plugin VERSION (which the runtime probe already returns) and labels an unreadable probe with the invented noun 'enablement'. Show the version pulled from the runtime, in plain words. Follows gj (#350, startup-sandbox-status)."
sprint: 0201-post-flip-release-model
group: ux-cleanup
sprint-readiness: ready
id: dag3bk4p0xe6tydc66k29ev3
status: done
worktree: .worktrees/spacedock-ensign-version-show-runtime-plugin-version
completed: 2026-06-13T18:58:58Z
verdict: PASSED
pr: "#354"
archived: 2026-06-13T18:59:23Z
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

## Stage Report: validation

- DONE: AC-1 — version-forward output, NO "enablement".
  Built `/tmp/da-val` from the worktree; `--version` renders `claude: spacedock 0.19.9` / `codex: spacedock 0.20.0` / `pi: spacedock ready` (exit 0). `git grep enablement -- internal/cli` returns only the negative-assertion banned-string list in `version_runtime_test.go` (not the old vocabulary); the `enablement` type + all old constants/funcs are GONE tree-wide (`type enablement`, `claudeEnablement`, `probe*Enablement`, etc. → no Go matches). `TestVersionVocabularyHasNoEnablementJargon` passes AND fails when jargon is reintroduced (adversarial refutation #2 caught it).
- DONE: AC-2 — version sourced ROBUSTLY from the resolved manifest, not the plugin-list probe.
  `probeVersion` → `execHost{}.ResolveManifest` → `contract.ManifestVersion` (reads the resolved plugin.json, the source doctor reads), structurally decoupled from `probeClaudeMarker`. `TestVersionMarkerUnknownRendersBareVersion` exercises the marker-failed-but-version-resolves case, asserts the bare `claude: spacedock 0.20.0` (no marker, not "not installed", no invented "unknown") — and FAILS under an edit that suppresses the version on markerUnknown (adversarial refutation #1 caught it). Live binary renders real manifest versions, confirming end-to-end.
- DONE: AC-3 — tests rewritten + suite green + clean additive helper.
  `go test ./...` exit 0, 1317 passed / 16 packages, no FAIL (the implementer-noted ensigncycle flake did not appear). `gofmt -l` clean, `go vet` exit 0. `contract.ManifestVersion` is purely additive — full `internal/contract/` diff is the new helper only; `ManifestVerdict`/`RunDoctor`/`Compare`/`readManifest` untouched. `git diff --name-only origin/main...HEAD` = 0 `.md` (yw owns the `## --version` doc — scope note honored, not flagged).
- DONE: Detached adversarial audit (front-door `--version` surface).
  Read-only refutation on a throwaway `/tmp/da-audit` checkout (deliverable never mutated, checkout removed). Refutation #1 (suppress version on markerUnknown) and #2 (reinject jargon) were both CAUGHT by the suite. Refutation #3 (make `contract.ManifestVersion` return a hardcoded constant) was NOT caught — see Polish below.

### Summary

PASSED. All three ACs verified by exercising behavior (live binary `--version`, command grep, suite exit code) — not prose reads. The version-forward block renders the locked five states, the version is sourced robustly from `contract.ManifestVersion` independent of the best-effort enabled marker, the suite is green (1317/16), and zero `.md` ship. Adversarial audit refuted two of three regression edits.

Polish (non-blocking, recorded): `contract.ManifestVersion` has NO direct unit test — adversarial edit #3 (return a hardcoded `99.99.99` instead of reading the manifest) left both `internal/contract` and `internal/cli` green. The cli suite injects a fake `runtimeProbe` and never drives the real `probeVersion`→`ManifestVersion` path, and no contract test calls the new helper. The shipped behavior IS proven (the live binary rendered the real manifest versions, and `ManifestVersion` is a 2-line delegation to the already-tested `readManifest`), so this is a test-strength gap, not an AC failure: a future refactor of the one-line helper would not be caught by the suite. A tiny `internal/contract` table test reading a fixture manifest's version would close it.

## Stage Report: implementation (polish fold — ManifestVersion direct test)

- DONE: Add an `internal/contract` direct unit test for `ManifestVersion` that closes the seam adversarial edit #3 exposed.
  `internal/contract/manifest_version_test.go` (commit 55ea7c6f): reads `testdata/compatible.json` and asserts the EXACT declared version `0.12.1` (not merely non-empty), plus the missing-file (`errNoManifest`) and unparseable-JSON (parse error) cases.
- DONE: Write it so it FAILS against a hardcoded-return edit, watch it catch the edit, then confirm green against the real helper.
  Hardcoded `ManifestVersion` to return `"99.99.99"` → all three cases FAILED (`= "99.99.99", want "0.12.1"`; missing-file and unparseable both lost their error). Reverted (clean diff vs committed), suite green again.
- DONE: Constraints — ZERO .md, whole-repo green, commit on the worktree branch.
  `go test ./...` 1320 passed / 16 packages (the three new cases); `gofmt`/`vet` clean; diff is one new `.go` test file; no `.md`.

### Summary

Closed the prove-the-seam gap the validator's audit surfaced: `ManifestVersion` now has a direct table test asserting the real fixture version (refutation #3 — a hardcoded return — is now caught), plus the missing-file and unparseable cases. Committed `55ea7c6f` on `spacedock-ensign/version-show-runtime-plugin-version`; whole-repo `go test ./...` green (1320/16).
