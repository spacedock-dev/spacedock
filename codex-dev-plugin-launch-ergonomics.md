---
title: Codex dev/local-plugin launch ergonomics — a `--plugin-dir`-equivalent for the local-marketplace dance
status: validation
source: "Captain request 2026-06-30 (0240 Commander session). Testing a local/dev build with Codex is far more cumbersome than Claude. Claude: `spacedock claude --plugin-dir <checkout>` — one flag, loads the local plugin checkout directly, bypasses installed-plugin resolution. Codex has no equivalent: you hand-build a local marketplace (`.agents/plugins/marketplace.json` with `source: local` + a `plugins/spacedock` symlink to the checkout), export `SPACEDOCK_MARKETPLACE_SOURCE`, run `spacedock install --host codex`, AND get the channel-in-the-name right (a plain `go build` edge binary needs the marketplace named `spacedock-edge`; the entry stays `spacedock`), plus the `.codex-plugin/plugin.json` version-masquerade gotcha."
group: tooling
id: 4q01qqyx4g2z3rctts1400av
sprint: 0240-lean-contract
started: 2026-07-01T02:38:48Z
worktree: .worktrees/spacedock-ensign-codex-dev-plugin-launch-ergonomics
---

## Problem
Launching Codex against a local/dev plugin checkout is a multi-step, footgun-laden ritual, while Claude is a single flag — and today's `--plugin-dir` codex path is not merely absent, it is **actively broken**: `runCodex` (`internal/cli/frontdoor.go:483`) relaxes the contract gate on `hasPluginDir` exactly like `runClaude` does (`frontdoor.go:313`), then forwards `--plugin-dir <dir>` straight into the real `codex` argv. Codex's CLI has no such flag and hard-rejects it (Spike A, below): `error: unexpected argument '--plugin-dir' found`. The flag is silently accepted by spacedock's own parser (`bindFrontDoorFlags`, `frontdoor.go:743`) and only fails downstream, at the host boundary, with a confusing message that never mentions spacedock.

- **Claude (the ergonomic target):** `spacedock claude --plugin-dir <checkout>` loads the local plugin directly via Claude's own native `--plugin-dir` flag, bypassing install resolution entirely (ephemeral — no on-disk install state changes).
- **Codex (the pain):** no working `--plugin-dir` path. The dev must (1) author a local marketplace dir — `.agents/plugins/marketplace.json` (`source: local`, `path: ./plugins/spacedock`) + a `plugins/spacedock` symlink to the checkout (the exact shape `internal/ensigncycle/codex_marketplace.go::writeCodexLocalMarketplace` already builds, CI-only); (2) `SPACEDOCK_MARKETPLACE_SOURCE=<dir> spacedock install --host codex`; (3) name the marketplace for the channel — a plain `go build` is `devBranch=next` → edge → the marketplace MUST be named `spacedock-edge` while the entry stays `spacedock`; wrong name → silently wrong channel (`ResolveManifest` returns `""`, no error — Spike E); (4) live with the version-masquerade (the installed manifest can report a version behind the checkout's real HEAD).

## Design

### Command surface
Two entry points sharing one implementation, mirroring pi's existing ephemeral/persistent split (`internal/cli/pi.go`: `runPi` vs `runInitWithPi`):

1. **`spacedock install --host codex --plugin-dir <checkout>`** — the persistent primitive. `--host`/`--plugin-dir` are already parsed host-generically for every host by `parsePiSetupArgs` (used by both `install` and `doctor`); today the `host != "pi"` branch in `runInitWithPi` (`pi.go:333-339`) unconditionally rejects a non-empty `pluginDir` with `"--plugin-dir is not supported"`. The design carves out `host == "codex"` from that rejection and routes it to the new shared helper instead of `runInit`.
2. **`spacedock codex --plugin-dir <checkout>`** — launch sugar. `runCodex` (`frontdoor.go:483`) gets a new step, inserted right after `safehouse.TranslateFlags` and before the existing gate-check block: if `fd.passthrough` carries `--plugin-dir <dir>`, strip it from the passthrough (codex's own CLI must never see it — that's Spike A's bug) and call the same shared helper to install it, before falling through to the **existing, unmodified** gate-check/launch logic. Because the flag is now stripped, `hasPluginDir(fd.passthrough)` reads false afterward and the normal `gateHost` check runs against the plugin the helper just installed — which is correct: unlike Claude's ephemeral bypass, a codex `--plugin-dir` install is a real on-disk install, so gate-checking it is meaningful, not redundant.

Both surfaces call one new unexported helper in `internal/cli` (co-located with `runCodex`/`execHost`, no new package):

```go
// installCodexLocalPluginDir builds a throwaway local marketplace from checkout
// (via WriteCodexLocalMarketplace, named for the binary's own channel) and
// installs it through the same Install() sequence a normal `spacedock codex`
// install uses, so the resulting plugin ID and cache layout are indistinguishable
// from a marketplace install on the same channel.
func installCodexLocalPluginDir(ops hostOps, checkout string, stderr io.Writer) error {
    marketplaceRoot, err := os.MkdirTemp("", "spacedock-codex-plugin-dir-")
    if err != nil {
        return fmt.Errorf("create local marketplace dir: %w", err)
    }
    defer os.RemoveAll(marketplaceRoot) // Spike F: codex COPIES into its cache at `plugin add`; the source need not outlive the call.

    install, err := WriteCodexLocalMarketplace(marketplaceRoot, checkout, channelMarketplace(devBranch))
    if err != nil {
        return fmt.Errorf("build local marketplace: %w", err)
    }
    if _, err := ops.Install("codex", install.marketplaceRoot, devBranch); err != nil {
        return fmt.Errorf("install from local marketplace: %w", err)
    }
    fmt.Fprintf(stderr, "Installed codex plugin from %s (version-masquerade advisory: the reported version reflects %s's checked-in .codex-plugin/plugin.json, not necessarily its current HEAD — see next-post-release-preversion-bump)\n", checkout, checkout)
    return nil
}
```

### Code relocation: `writeCodexLocalMarketplace` → `internal/cli`
Today's helper (`internal/ensigncycle/codex_marketplace.go:16`) hardcodes `Name: "spacedock"` — fine for a CI-only harness that always tests the stable-named marketplace, wrong for a user-facing command that must name the marketplace after whatever channel the binary itself resolves. The design:
- Move it into `internal/cli` as exported `WriteCodexLocalMarketplace(marketplaceRoot, repoRoot, marketplaceName string) (CodexMarketplaceInstall, error)` — same body, `Name: "spacedock"` → `Name: marketplaceName`. Spike B proved this exact signature change compiles and keeps all existing tests green.
- `internal/ensigncycle`'s two call sites (`codex_marketplace_test.go`, `codex_live_runner_test.go`) become `cli.WriteCodexLocalMarketplace(dir, repo, "spacedock")` — the literal string preserves today's exact CI behavior (both existing scenarios pin the stable name), zero behavior change, confirmed by Spike B's `go build && go vet -tags live && go test` pass.
- Import direction: `internal/ensigncycle` (test/CI-harness package) importing `internal/cli` (production package) is the correct direction — a grep confirmed neither package imports the other today, so this is a new edge, and it points the way test-harness → production code should, not backwards.
- Explicitly out of scope: migrating `.github/workflows/runtime-live-e2e.yml`'s hand-rolled shell equivalent (lines 330-509, including the standalone edge-channel step at 412-498) onto the Go helper. That YAML never called the Go function to begin with; leaving it alone avoids scope creep into CI-workflow changes this task doesn't need.

### Auto channel-name resolution
No new logic — reuse `channelMarketplace(devBranch)` (`host_exec.go:235`), the same function `channelMarketplaceSource`/`installArgvSequence` already use for the normal marketplace install. The footgun today isn't a missing mechanism, it's that `writeCodexLocalMarketplace` never consulted it. Spikes D and E prove the fix and the baseline it replaces.

### Version-masquerade — advisory, not the full fix
`next-post-release-preversion-bump` (backlog) owns re-architecting version stamping; its "Out of scope" section rules out this task's real fix living there either. This design's scope is a **printed advisory**, not a stamping change: `installCodexLocalPluginDir` prints the note above on every `--plugin-dir` install (both command surfaces, since both call the shared helper), naming the checkout and pointing at the backlog item. This satisfies the entity's own bracketed fallback ("surfaces the real version (**or at least does not silently masquerade**)").

### Snapshot semantics (Spike F finding — a design constraint, not a gap)
Verified live: `codex plugin add` **copies** the marketplace's plugin directory into `<CODEX_HOME>/plugins/cache/<marketplace>/<entry>/<version>/` as real files, not a symlink chain. This means (a) `installCodexLocalPluginDir`'s marketplace root can be a throwaway temp dir removed immediately after `Install` returns — no persistent on-disk marketplace to manage or clean up later — but also (b) a `--plugin-dir` codex install is a **point-in-time snapshot**, unlike Claude/Pi's live ephemeral override: editing the checkout after installing does NOT change what codex runs until `spacedock codex --plugin-dir` (or `spacedock install --host codex --plugin-dir`) is re-run. The docs diff below calls this out explicitly so it isn't a silent surprise.

## Acceptance criteria

**AC-1 — one command, no manual marketplace authoring, and no regression to today's broken passthrough.**
`spacedock codex --plugin-dir <checkout>` launches Codex against the local checkout's plugin in one command, with no operator-authored marketplace file and no `--plugin-dir` token forwarded into the real `codex` argv.
*Verified by:* a `hostOps`-fake-backed unit test in `internal/cli` (e.g. `TestRunCodexPluginDirInstallsThenLaunchesWithoutTheFlag`) asserting (a) the fake `Install` is called exactly once with `host="codex"` and a source directory whose `.agents/plugins/marketplace.json` `name` equals `channelMarketplace(devBranch)` and whose `plugins/spacedock` resolves (via symlink) to the given checkout, and (b) the argv passed to the fake `Launch` contains no `--plugin-dir` element anywhere. (b) is a direct regression guard against Spike A's reproduced baseline (today's argv DOES contain it, and the real `codex` binary rejects it) — a re-introduction of the passthrough bug flips this assertion, not just "looks wrong" prose.

**AC-2 — the channel-name footgun is closed, measured against the reproduced baseline.**
A `--plugin-dir` codex install always names its marketplace via `channelMarketplace(devBranch)`, so an edge-devBranch build's install resolves through the real `ResolveManifest`, where the old hardcoded name silently did not.
*Verified by:* a new non-`live`-tagged test in `internal/cli` (style of `codex_resolve_test.go`, skips if `codex` not on `PATH`, no auth required — Spike C), e.g. `TestInstallCodexLocalPluginDirResolvesOnEdgeChannel`: pins `devBranch = "next"`, points `CODEX_HOME` at a fresh temp dir, calls the real `installCodexLocalPluginDir` against a throwaway checkout fixture, then asserts `execHost{}.ResolveManifest("codex")` returns a non-empty path. This is the independent baseline pairing the checklist requires: Spike E already proved, with the real production `ResolveManifest` and the OLD hardcoded `"spacedock"` name under the same `devBranch="next"`, that resolution returns `""` — a regression back to a hardcoded name reintroduces that exact empty-resolve failure, so this AC can genuinely move the wrong way, not just start at zero.

**AC-3 — the version-masquerade is surfaced, not silent.**
Every `--plugin-dir` codex install prints an explicit advisory that the reported version reflects the checkout's checked-in manifest, not necessarily its current HEAD; a normal (non-`--plugin-dir`) install prints no such advisory.
*Verified by:* a stderr-capturing unit test asserting the advisory substring is present after `installCodexLocalPluginDir` runs and absent from the plain `runInit`/`runCodex` non-`--plugin-dir` install path (a presence/absence pair, not a presence-only check, so the test can't pass by printing the advisory unconditionally).

## Test plan
- All three ACs above land as fast, non-`live`-tagged `internal/cli` tests (AC-1/AC-3 fully faked via `hostOps`; AC-2 shells the real `codex` binary but needs no auth, matching `codex_resolve_test.go`'s existing style — skip, don't fail, when `codex` is absent from `PATH`).
- `internal/ensigncycle`'s two existing tests (`TestWriteCodexLocalMarketplacePointsAtCurrentCheckout`, and the `codex-live` scenario runner) get their call sites updated to the relocated `cli.WriteCodexLocalMarketplace(dir, repo, "spacedock")` signature and must stay green unchanged — this is the CI-non-breakage half of the checklist, proven feasible by Spike B.
- No CI YAML changes proposed by this design (see "Code relocation," out-of-scope bullet).

## Spike results
All spikes used the real, unmodified production code paths and the real `codex` CLI (no reimplementation), and were fully reverted/cleaned up after each — `git status --short` returns clean on `internal/cli` and `internal/ensigncycle` throughout, confirmed as of this writeup.

- **Spike A (today's failure mode, live):** Built the binary, ran `spacedock codex --skip-contract-check --plugin-dir <checkout> -- --help` from a directory with no `.safehouse` profile (the repo root's own `.safehouse` profile masks the bug behind an unrelated safehouse-unavailable early exit — first run from the repo root gave a misleading result; re-ran from `/tmp` to reach the real `codex` exec). Result: `error: unexpected argument '--plugin-dir' found` — confirms the Problem statement's "actively broken," not merely absent.
- **Spike B (relocation safety):** Edited `writeCodexLocalMarketplace` to take a `marketplaceName string` param (`Name: "spacedock"` → `Name: marketplaceName`) and updated both call sites to pass `"spacedock"` explicitly. `go build ./...`, `go vet -tags live ./...`, and `go test` on the affected packages all passed. Reverted via `git checkout --`. Confirms the signature change is a safe, zero-behavior-change refactor.
- **Spike C (unauthenticated marketplace ops):** `codex plugin marketplace add <dir>` / `plugin add <entry>@<marketplace>` / `plugin list` all succeeded against an isolated `CODEX_HOME` with no auth configured. Confirms AC-2's test can be a plain (non-`live`-tagged) test, not gated behind the `live` build tag.
- **Spike D (edge-channel fix, live):** Built a local marketplace named `spacedock-edge` (the `channelMarketplace("next")` value), installed it into an isolated `CODEX_HOME`, then called the real `execHost{}.ResolveManifest("codex")` with `devBranch` pinned to `"next"`. Result: resolved a non-empty manifest path. Confirms the fix half of AC-2.
- **Spike E (edge-channel baseline, live):** Same setup but with the marketplace named `"spacedock"` (today's hardcoded default) under the same `devBranch="next"`. Result: `ResolveManifest` returned `""` (no error) — reproduces the exact "wrong name → silently wrong channel" footgun the entity describes, using real production code. Confirms the baseline half of AC-2.
- **Spike F (cache semantics, live):** Added a local marketplace + plugin under an isolated `CODEX_HOME`, then inspected `<CODEX_HOME>/plugins/cache/<marketplace>/spacedock/<version>/` on disk: it is a real, fully-copied directory tree (not a symlink chain back to the marketplace source). Confirms `installCodexLocalPluginDir` may safely use a throwaway temp marketplace dir removed immediately after install, and surfaced the snapshot-semantics design constraint documented above (not discovered any other way — the existing CI harness never exercises the plugin after its `t.TempDir()` marketplace is cleaned up).

## Docs diff (proposed — this stage does not apply it)

`docs/site/get-started/install.md`, replacing lines 71-73:

```diff
-Launching with `--plugin-dir` loads a local plugin checkout directly and bypasses
-installed-plugin resolution, but it does not wrap the launch in the safehouse
-sandbox — use it for plugin development, not as an install substitute.
+Launching with `--plugin-dir` loads a local plugin checkout directly. On Claude
+and Pi this is an ephemeral, install-free override — it bypasses installed-plugin
+resolution for that one launch and does not wrap the launch in the safehouse
+sandbox; use it for plugin development, not as an install substitute.
+
+Codex has no such flag on its own CLI, so `spacedock codex --plugin-dir
+<checkout>` and `spacedock install --host codex --plugin-dir <checkout>` build a
+local marketplace from the checkout and install it under the binary's own
+channel (`spacedock` stable / `spacedock-edge` edge — matching whatever
+`spacedock codex` would otherwise install), then launch. This IS a persistent
+install, replacing whatever Codex plugin was previously configured, and it is a
+point-in-time snapshot: editing the checkout afterward has no effect until the
+command is re-run. The command prints an advisory that the reported version
+reflects the checkout's checked-in manifest, not necessarily its current HEAD.
```

`docs/runtime-live-ci.md` (optional pointer, near lines 94-101 describing the codex-live CI marketplace mechanics): add one sentence noting that `internal/ensigncycle`'s local-marketplace helper now calls the same `internal/cli.WriteCodexLocalMarketplace` that backs `spacedock codex --plugin-dir` / `spacedock install --host codex --plugin-dir`, so CI and the user-facing command share one implementation (not a parallel one CI must keep in sync by hand).

## Related
- Ergonomic target to match: Claude's `--plugin-dir`.
- Reusable building block, relocating in this design: `internal/ensigncycle/codex_marketplace.go::writeCodexLocalMarketplace` → `internal/cli.WriteCodexLocalMarketplace`.
- Existing prior art for the command-surface split: `internal/cli/pi.go`'s `runPi` (ephemeral) vs `runInitWithPi` (persistent) `--plugin-dir` handling.
- Adjacent backlog, NOT resolved by this design: `next-post-release-preversion-bump` (the full version-masquerade fix; this design ships an advisory only).
- Docs to update if shipped: `docs/site/get-started/install.md` (diff above), optionally `docs/runtime-live-ci.md` (pointer above).

## Stage Report: ideation

- DONE: Flesh the command-surface design (spacedock codex --plugin-dir or install --host codex --plugin-dir) wired to writeCodexLocalMarketplace, auto channel-name resolution, and the version-masquerade fix
  See "## Design" — persistent primitive (`install --host codex --plugin-dir`, extending `pi.go:333-339`) plus launch sugar (`runCodex`), both sharing one new `installCodexLocalPluginDir` helper built on a relocated `cli.WriteCodexLocalMarketplace`.
- DONE: Tighten the rough acceptance sketch into measured ACs each naming its test (behavioral live-or-fixture Codex install smoke, not prose)
  See "## Acceptance criteria" — AC-1 (argv/install-call assertions + Spike A regression guard), AC-2 (real `ResolveManifest` against the Spike E baseline), AC-3 (advisory presence/absence pair).
- DONE: Spike the riskiest mechanism first: promoting writeCodexLocalMarketplace from a CI-only helper into a user-facing command path without breaking the CI harness's own usage
  Six real, reverted spikes (A-F) recorded under "## Spike results": A confirmed today's `--plugin-dir` is actively broken (`error: unexpected argument '--plugin-dir' found`); B proved the signature promotion is a safe zero-behavior-change refactor (build/vet/test green, then `git checkout --` reverted); C-E proved the channel-name fix and reproduced its baseline against the real `codex` CLI and the real `ResolveManifest`; F discovered `codex plugin add` copies into its cache (no persistent marketplace dir needed), which shaped the helper's temp-dir design and the snapshot-semantics doc note. `git status --short` confirmed clean after each.

### Summary

Reframed the problem from "no `--plugin-dir` path" to "an actively broken one" (Spike A), then designed a single shared `installCodexLocalPluginDir` helper feeding both the persistent `spacedock install --host codex --plugin-dir` primitive and the `spacedock codex --plugin-dir` launch sugar, reusing the existing `channelMarketplace`/`Install` machinery so channel-naming needs no new logic. All three ACs are measured against real evidence: two against reproduced before/after baselines (Spikes A and D/E), one as a presence/absence pair. A proposed `docs/site/get-started/install.md` diff documents the resulting persistent, point-in-time-snapshot semantics (Spike F) that distinguish codex's `--plugin-dir` from Claude/Pi's ephemeral one.

## Stage Report: implementation

- DONE: Relocate writeCodexLocalMarketplace -> exported internal/cli.WriteCodexLocalMarketplace(marketplaceRoot, repoRoot, marketplaceName) (Name hardcode -> param) and update the two internal/ensigncycle call sites to pass "spacedock" verbatim (zero CI behavior change)
  New `internal/cli/codex_marketplace.go` (exported type `CodexMarketplaceInstall` with exported fields, since the ensigncycle tests read them cross-package); deleted `internal/ensigncycle/codex_marketplace.go`; both call sites pass "spacedock". `go vet -tags live ./internal/ensigncycle/` + `go test ./internal/ensigncycle/` green. Commit 1e9098be.
- DONE: Add the installCodexLocalPluginDir helper + wire both surfaces (install --host codex --plugin-dir carve-out in pi.go's host!=pi rejection; runCodex --plugin-dir strip-from-passthrough + install in frontdoor.go before the gate check) and print the version-masquerade advisory on every --plugin-dir install
  Helper in codex_marketplace.go; `takeCodexPluginDir` + runCodex insertion in frontdoor.go; `host == "codex"` carve-out in pi.go `runInitWithPi`. Advisory printed by the shared helper (both surfaces). E2E-verified with the built edge binary + real codex 0.141.0: install surface + launch surface both work, zero Spike A `--plugin-dir` rejections.
- DONE: Land the 3 ACs + apply the docs/site/get-started/install.md diff; go build ./... and go test ./internal/cli/ ./internal/ensigncycle/ green
  AC-1/AC-3 faked in `codex_plugin_dir_test.go`; AC-2 (`TestInstallCodexLocalPluginDirResolvesOnEdgeChannel`) RAN against real codex (0.16s, not skipped) and passed. Docs diff applied to install.md (+ optional runtime-live-ci.md pointer). `go build ./...`, `go vet ./...`, `go test ./...` all green (offline gate clean).

### Summary

Two design assumptions were falsified by running the real codex 0.141.0 and corrected. (1) **Spike F was wrong**: codex records the marketplace source by path and re-loads it on every later `codex plugin` command, so a temp marketplace `os.RemoveAll`'d right after install hard-fails every subsequent codex invocation ("marketplace root does not contain a supported manifest"); removing the marketplace registration instead *uninstalls* the plugin. Fix: the local marketplace is rooted persistently under `codexHome()` (one stable dir per channel, re-pointed each call), which also auto-isolates in tests via `CODEX_HOME`. (2) Per the design, codex `--plugin-dir` now **installs-then-gates** rather than relaxing the gate (unlike claude/pi's ephemeral bypass), so three existing tests that encoded the old forward-and-relax/rejection behavior were updated: `TestPluginDirRelaxesGate` (codex sub-case → install-then-gate), `TestStrayPromptGuardNegatives` (codex `--plugin-dir` → consumed, not forwarded), and `TestNonPiSetupRejectsPluginDir` (dropped the codex-install rejection case). One guard note: the AC-1 fake originally read `.agents/plugins/marketplace.json` to assert the marketplace name, but the repo boundary guard false-positives on `.agents` (substring of the `agents` skill segment); the manifest-name assertion was dropped from the faked test since AC-2 already proves the channel-name property behaviorally against real codex.

## Stage Report: validation

- DONE: AC-1 (argv carries NO --plugin-dir + Install called once host=codex + marketplace name == channelMarketplace(devBranch) + plugins/spacedock resolves to the checkout; the Spike A regression guard)
  `TestRunCodexPluginDirInstallsThenLaunchesWithoutTheFlag` PASS; adversarial audit on a throwaway worktree: re-forwarding the flag (Spike A regression) flips the assertion to FAIL ("launch argv forwards a --plugin-dir token"), reverted → PASS. NOTE: the AC-1 body's "Verified by" names a marketplace-name assertion the fake does NOT make (dropped over the `.agents` boundary-guard false-positive); that property is instead proven behaviorally in AC-2 — coverage is complete but AC-1's prose is stale vs the test.
- DONE: AC-2 (real ResolveManifest non-empty on edge channel vs the Spike E hardcoded-"spacedock" empty baseline)
  `TestInstallCodexLocalPluginDirResolvesOnEdgeChannel` RAN against real codex 0.141.0 (0.16s, not skipped) → PASS. Adversarial audit: regressing the marketplace name to hardcoded "spacedock" under devBranch="next" makes the real `codex plugin add spacedock@spacedock-edge` FAIL at the install boundary (channel name is load-bearing), reverted → PASS. Genuinely discriminates.
- DONE: AC-3 (advisory presence on --plugin-dir + absence on the plain install — a real pair)
  `TestCodexPluginDirAdvisoryPresenceAndAbsence` both sub-cases PASS; `TestInstallCodexPluginDirInstallsViaSharedHelper` confirms the persistent-primitive surface also prints it. Pair structure means an unconditional print would break the absence case.
- DONE: go build ./... and go test ./internal/cli/ ./internal/ensigncycle/ green from the repo root; two ensigncycle call sites pass "spacedock" verbatim; internal/ensigncycle/codex_marketplace.go deleted
  `go build ./...` green; `go vet ./internal/cli/ ./internal/ensigncycle/` + `go vet -tags live` green; both test packages PASS. Both call sites (`codex_marketplace_test.go:17`, `codex_live_runner_test.go:111`) pass literal "spacedock"; old ensigncycle file confirmed absent (moved to internal/cli).
- DONE: Confirm docs/site/get-started/install.md diff applied + docs/runtime-live-ci.md pointer
  install.md documents both `--plugin-dir` surfaces, persistent point-in-time-snapshot semantics, and the version-masquerade advisory (matches proposed diff); runtime-live-ci.md adds the shared-`WriteCodexLocalMarketplace`-builder pointer.

### Summary

PASSED. All three ACs reproduced (not asserted) against real codex 0.141.0, and each was adversarially audited on a throwaway worktree to confirm it can genuinely move the wrong way — AC-1 flips when the Spike A passthrough is re-forwarded, AC-2 flips when the channel name regresses to the Spike E baseline. Build/vet/test green, ensigncycle call sites unchanged in behavior, old helper deleted, both docs surfaces applied. One non-blocking note for the FO/captain: AC-1's "Verified by" prose still describes a marketplace-name assertion that was intentionally relocated to AC-2's behavioral proof (disclosed in the implementation report); the property is fully covered, but the AC-1 body should be reconciled so it doesn't over-specify a check the test no longer makes.
