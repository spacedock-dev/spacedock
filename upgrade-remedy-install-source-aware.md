---
id: vzsastkvv2r6dpjakw1vq6wx
title: Binary-upgrade prompt must be install-source-aware (brew formula, non-brew, sandbox)
status: implementation
source: "Captain report (CL) 2026-07-16 — the skill-upgrade version gate prompts `brew upgrade spacedock` regardless of how the binary was installed."
started: 2026-07-21T15:58:31Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-upgrade-remedy-install-source-aware
issue:
mod-block:
pr:
---

When the skill is upgraded and determines the binary needs upgrading, the remedy it emits hardcodes `brew upgrade spacedock` — even when that is not how the binary was installed. It must detect the real install source and runtime context and emit the correct upgrade instruction. Likely a follow-up gap in the shipped upgrade-hint work (`install-refresh-and-upgrade-hint`, `init-upgrade-and-contract-remedy`).

## Problem

> **Captain TRIM (2026-07-22).** The all-install-sources design (a `BrewStable`/`BrewEdge`/`NonBrew`/sandbox `InstallSource` matrix) was implemented and validated PASSED, then the captain trimmed the scope to ONLY the actual reported bug — the edge (`@next`) channel. The sections in THIS spec region (Problem → Expected surface) are revised to that @next-only scope. The historical ideation / implementation / validation reports below describe the superseded four-source design and are kept as the record.

The too-old-binary remedy is generated in ONE place: `internal/contract/contract.go` → `tooOldBinaryRemedy()` (called from `compareNamed`, reached by `spacedock doctor`, `spacedock init/install`'s post-install doctor, and the `spacedock claude|codex` front-door version gate via `ManifestVerdict`). It returns a fixed block that leads with `brew upgrade spacedock`:

```
  Upgrade via Homebrew: brew upgrade spacedock
  Or build from source: go build -o spacedock ./cmd/spacedock
  Or refresh the plugin instead: spacedock install
```

**The reported bug (captain-observed).** A user who installed via the `spacedock@next` cask ran the prompt, followed `brew upgrade spacedock`, and it did nothing — the stable cask isn't installed. The edge cask is a *separate* cask token (`spacedock@next`, `conflicts_with cask: spacedock`), so the correct command is `brew upgrade spacedock@next`. That is the only case this task fixes; every other install keeps today's unchanged block.

No skill/agent text and no user-facing doc restates the `brew upgrade spacedock` remedy, so the fix is contained to the one code remedy generator.

## Proposed approach

Minimal @next-only fix: thread a small `edgeCask bool` into the single remedy generator (mirroring how `tooOldPluginRemedy(host)` already threads a value), and detect the one case — "is the running binary the `spacedock@next` cask?" — in `internal/cli`. The `contract` package stays free of `os.Executable`/`exec`: it receives a bool and picks the Homebrew formula.

1. **`tooOldBinaryRemedy(edgeCask bool)` in `internal/contract`.** `edgeCask=true` swaps line 1's formula to `brew upgrade spacedock@next`; the other two lines are unchanged. `edgeCask=false` reproduces today's 3-line block byte-for-byte — the safe default and the reason the public `Compare` signature stays untouched (it passes `false`).

2. **Thread the bool through the remedy path only:** `compareNamed`, `ManifestVerdict`, and `RunDoctor` gain an `edgeCask bool` param; `Compare` keeps its signature and passes `false`. Callers computing it: `frontdoor.go` `gateHost` (×1) and `init.go`'s doctor / codex / manifest arms (×3), each calling `runningEdgeCask()`.

3. **Detection `runningEdgeCask()` / `isEdgeCaskPath(execPath)` in `internal/cli/edge_cask.go`.** Resolve the RUNNING binary via the existing `resolvedLauncherBin()` (`os.Executable` → `filepath.EvalSymlinks`), split the path, and report whether the segment immediately after a `Caskroom` segment is exactly `spacedock@next`. Any other path (the stable `spacedock` cask, a source checkout, an empty/unresolvable path) → `false` → the unchanged remedy. No `brew` subprocess (the resolved Caskroom token carries the cask name; spike-proven below).

**Why a bool, not the `SourceKind` enum?** The captain's reported bug is exactly one case (`@next`). A single bool is the smallest value that fixes it and keeps `contract` pure; the non-brew / stable-formula / sandbox arms of the earlier enum design are out of scope (below).

## Out of scope

- **Every install source except the `@next` edge cask.** Non-brew builds (source `go build`, proxy `go install`, downloaded archive), a separate remedy for the stable `spacedock` cask, and sandbox / brew-unreachable execution are NOT addressed — they keep today's unchanged 3-line block. The captain trimmed the four-source design that had covered these; if any recurs as a real report, it is a separate task.
- **The `SourceKind`/`InstallSource` enum, `HostOnly`/sandbox detection, and the per-source rendering/detection matrices.** Replaced by the single `edgeCask` bool.
- **Version-COMPARISON logic** (`ParseMajorMinor`, `Compare`, verdict classification) — untouched. Only the too-old-binary remedy's Homebrew formula line changes, and only for an @next install.
- **The other verdicts' messages** (compatible, too-old-plugin, malformed-version, no-plugin-found) — byte-identical.
- **Doc change.** The command-reference never described the too-old-binary remedy, so no doc edit is made (keeping the trim minimal).
- **Auto-upgrade / self-update.** The remedy tells the user what to run; it never runs it.

## Acceptance criteria

- **AC-1 (the captain's regression — an @next install names the @next formula, never the bare stable command).** When the running binary is the `spacedock@next` cask (`edgeCask=true`), the too-old-binary remedy names `brew upgrade spacedock@next` and contains NO bare stable `brew upgrade spacedock` command, word-boundary matched so `brew upgrade spacedock@next` does NOT satisfy the stable substring. (This is exactly the `@next`-user-got-`brew upgrade spacedock`-no-op bug.)
  - *Test:* `internal/contract/version_message_test.go` `TestTooOldBinaryRemedyEdgeChannel` — `tooOldBinaryRemedy(true)` asserts `brew upgrade spacedock@next` present AND regexp `` `brew upgrade spacedock(\s|$)` `` absent.
- **AC-2 (no collateral change).** For any non-@next install (`edgeCask=false`), the remedy is byte-for-byte the existing 3-line block, and the compatible / too-old-plugin / malformed / no-plugin messages are unchanged — so the untouched `Compare`-based verdict tests pass without edits.
  - *Test:* same `TestTooOldBinaryRemedyEdgeChannel` pins `tooOldBinaryRemedy(false)` byte-for-byte; existing `internal/contract/contract_test.go` and `TestTooOldBinaryRemedyLeadsWithBrew` (call-only reconcile) stay green.
- **AC-3 (detection grounded in the resolved path).** `isEdgeCaskPath` classifies a `…/Caskroom/spacedock@next/…` path as edge (true), and the stable cask / a source checkout / an empty path as not-edge (false).
  - *Test:* `internal/cli/edge_cask_test.go` `TestIsEdgeCaskPath` table over four representative paths.

## Test plan

Focused Go unit tests feeding inputs and asserting emitted bytes — no mocks.

1. **`internal/contract/version_message_test.go` — add `TestTooOldBinaryRemedyEdgeChannel`.** `tooOldBinaryRemedy(true)` → contains `brew upgrade spacedock@next`, the word-boundary regexp for the bare stable command does NOT match, and it keeps the `spacedock install` plugin-refresh line. `tooOldBinaryRemedy(false)` → byte-for-byte the pinned 3-line block.
2. **`internal/cli/edge_cask_test.go` (new) — `TestIsEdgeCaskPath`.** Table: edge Caskroom path → true; stable Caskroom path → false; source checkout → false; empty → false.
3. **Signature-ripple caller reconciles:** `internal/contract/doctor_test.go` (×2), `internal/cli/upgrade_from_stale_test.go` (×2), and `version_message_test.go` (×3) pass `false`; `frontdoor.go` (×1) and `init.go` (×3) pass `runningEdgeCask()`. `Compare` callers are UNCHANGED.
4. **Hygiene:** `gofmt -l` (clean), `go vet ./...`, `go test ./...`, `go test ./... -race`.

## Expected surface

| File | Kind | as-built |
|---|---|---|
| `internal/contract/contract.go` | prod — `tooOldBinaryRemedy(edgeCask)`, thread `compareNamed`/`Compare` | +18/-8 |
| `internal/contract/doctor.go` | prod — `ManifestVerdict`/`RunDoctor` gain `edgeCask` | +10/-6 |
| `internal/cli/edge_cask.go` | prod — `runningEdgeCask`, `isEdgeCaskPath` (new file) | +37 |
| `internal/cli/frontdoor.go` | prod — `gateHost` passes `runningEdgeCask()` | +1/-1 |
| `internal/cli/init.go` | prod — 3 doctor callers pass `runningEdgeCask()` | +3/-3 |
| `internal/contract/version_message_test.go` | test — reconcile + `TestTooOldBinaryRemedyEdgeChannel` | +33/-3 |
| `internal/cli/edge_cask_test.go` | test — `TestIsEdgeCaskPath` (new) | +28 |
| `internal/contract/doctor_test.go`, `internal/cli/upgrade_from_stale_test.go` | test — ripple `false` | +4/-4 |

As-built: 9 files, +134/-25 (net +109). No doc change (out of scope per the trim). Roughly ⅓ the surface of the superseded four-source build (+416/-32).

## Docs

No doc change: the command-reference never described the too-old-binary remedy, so the @next-only trim adds none. (The four-source design's proposed command-reference note is dropped with the rest of that scope.)

## Spike (riskiest mechanism — DONE)

**Question:** can we reliably detect which Homebrew cask owns the running binary vs a non-brew path, without shelling out to `brew`?

**Result — proven on this machine (edge box):**
- `command -v spacedock` → `/opt/homebrew/bin/spacedock`; `realpath`/`EvalSymlinks` → `/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock`. The segment after `Caskroom/` is `spacedock@next` — exactly the cask/formula name to emit. `brew --prefix` Caskroom listing confirms `spacedock@next` is the installed cask.
- The edge cask ships `binary "spacedock"` (`internal/release/testdata/spacedock@next.rb`), so Homebrew symlinks `<prefix>/bin/spacedock` → the Caskroom staged path; `EvalSymlinks` (already used by `resolvedLauncherBin()` at `internal/cli/frontdoor.go:125`) resolves it with no `brew` subprocess.
- Anchor on `os.Executable` of the RUNNING process, NOT `command -v`: on this box `SPACEDOCK_BIN=/Users/clkao/git/…/spacedock` is a source checkout build, distinct from the brew binary on PATH — so which binary "owns" the session depends on which one is executing. `resolvedLauncherBin()` already keys off `os.Executable`, so a checkout-launched session correctly classifies as `NonBrew` and a brew-launched one as `BrewEdge`.

**No spike needed for** the sandbox `HostOnly` branch (a plain `exec.LookPath("brew")` with an injectable stub, mirroring `safehouse.Available(lookPath)`) or the message threading (mirrors the proven `tooOldPluginRemedy(host)` host-threading already in `contract.go`).

## Stage Report: ideation

- DONE: Fill in Problem, Proposed approach, and Out of scope in the task body from the seeds
  Problem names the single generator `contract.go:tooOldBinaryRemedy()` and the 3 wrong cases; approach and out-of-scope written with rejected-alternative justification.
- DONE: Design install-source detection distinguishing brew stable, brew edge (@next), non-brew, and sandbox
  `detectInstallSource(execPath, brewLookPath, devBranch)` → `contract.InstallSource{Kind, HostOnly}`; brew-vs-nonbrew + formula from the resolved Caskroom token, sandbox via `brew` unresolvable.
- DONE: Specify the exact upgrade instruction emitted per source
  Per-source wording block in Proposed approach: `brew upgrade spacedock@next` (edge), `brew upgrade spacedock` (stable), source rebuild (non-brew), run-on-host (sandbox); all keep the `spacedock install` plugin-refresh line.
- DONE: Locate and name the remedy code path and any skill abort text; record before/after wording
  Path: `internal/contract/contract.go:tooOldBinaryRemedy()` via `compareNamed`←`ManifestVerdict`/`RunDoctor` (front-door gate + `init.go` doctor). Grep confirms NO skill/agent/site-doc restatement of the binary remedy. Before/after wording recorded.
- DONE: Write Acceptance criteria including a value-measuring AC and a NO-`brew upgrade spacedock`-for-@next/non-brew/sandbox AC
  AC-1 (matrix value, oracle external to message), AC-2 (word-boundary-matched no-bare-stable for edge/non-brew), AC-3 detection, AC-4 no-collateral, AC-5 docs.
- DONE: Write the Test plan (table tests across the source matrix; gofmt, go test ./..., go test ./... -race)
  Detection + rendering tables in `install_source_test.go`, per-source in `version_message_test.go`, skip-guarded live doctor smoke, hygiene commands listed.
- DONE: Declare the expected surface (files + LOC) and tolerance
  Table: ~9 files, prod ~117 / test ~178 / doc ~3 LOC; tolerance ±3 files, ±60 LOC (ripple capped by keeping `Compare` stable).
- DONE: Spike the riskiest unverified mechanism (which formula/tap owns the running binary vs a non-brew path)
  Proven on this edge box: `EvalSymlinks` → `/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock`; token after `Caskroom/` is the exact cask name; no `brew` subprocess; anchor on `os.Executable` (SPACEDOCK_BIN here is a distinct source checkout).
- DONE: Keep version-COMPARISON logic out of scope; propose the doc diff for user-visible remedy text
  Out of scope lists comparison logic + other verdict messages; doc diff proposed for `docs/site/reference/command-reference.md`.

### Summary

Filled the entity with a concrete, testable, source-aware remedy design. The remedy is generated in one place (`tooOldBinaryRemedy()`); the plan keeps `contract` pure by threading a detected `InstallSource` value (mirroring existing `host` threading) and keeps the public `Compare` signature stable to cap the caller ripple. The riskiest mechanism — reading the owning cask from the resolved Caskroom path without shelling out to `brew` — was spiked live and confirmed (`spacedock@next` token extracted from the real path). Key subtlety banked for implementation: `brew upgrade spacedock@next` contains the substring `brew upgrade spacedock`, so the "no bare stable" assertion must be word-boundary matched. No user-facing doc or skill text restates the binary remedy, so the doc change is a single command-reference addition.

## Gate: ideation — APPROVED (FO)

- **Verdict:** approved for implementation. Threads a detected `InstallSource` value into the single remedy generator (keeps `contract` pure, `Compare` signature stable to cap the caller ripple), spike proven live on the edge box (`Caskroom/spacedock@next` token extracted from the resolved path, no `brew` subprocess), and correctly catches the word-boundary subtlety. No captain decision required.
- **Validation carries the detached adversarial audit** — threads a param through `frontdoor.go gateHost` (the front-door launcher, a high-stakes surface). AC-1's oracle is the external install-source input, so the AC-provenance sub-trigger does not fire.
- **Implementation base:** worktree off `origin/main` (`ca136f83`), not local `main`.

### Gate note — implementation surface overage ACCEPTED (FO)

Implementation landed at +416/-32 (~+118 over the ±60 LOC tolerance). Accepted as an estimation miss, NOT scope growth: the 10 touched files are exactly the approved surface (no new files, no new behavior), all value ACs are met with falsifiable evidence, `Compare`'s signature stayed stable (its 15 cases pass unedited), and the overage is concentrated in AC-mandated test tables (`install_source_test.go` 198: the AC-1 four-source matrix + AC-3 six-row detection table + live smoke) plus a `currentInstallSource` DRY helper and repo-norm comment density. Trimming would cut AC-mandated coverage. The tolerance baseline is recalibrated to the as-built surface for any later correction round; proceeding to validation.

## Stage Report: implementation

- DONE: tooOldBinaryRemedy is install-source-aware via a threaded contract.InstallSource; emits the right instruction per source (brew stable / brew upgrade spacedock@next / non-brew rebuild / sandbox run-on-host) with NO bare `brew upgrade spacedock` for @next or non-brew (word-boundary matched) — AC-1/AC-2.
  `tooOldBinaryRemedy(src)` switches on `InstallSource` (contract.go); code commit `b323cced`. AC-1: cli `TestTooOldBinaryRemedyRendersPerSource` drives each source through `contract.RunDoctor` — reddens if any arm's emitted token stops matching its source input (the oracle is the source Kind, external to the message). AC-2: that test AND contract `TestTooOldBinaryRemedyPerSource` assert `regexp` `` `brew upgrade spacedock(\s|$)` `` does NOT match the edge/non-brew/host-only rows while `spacedock@next` IS present — fails the instant the remedy regresses to the bare stable command.
- DONE: detectInstallSource in internal/cli reads the resolved Caskroom token off os.Executable; Compare signature stays unchanged (ripple capped to ManifestVerdict/RunDoctor callers); SourceUnknown reproduces today's 3-line block byte-for-byte (AC-4).
  New `internal/cli/install_source.go` (`detectInstallSource`/`caskToken`/`currentInstallSource`), anchored on `resolvedLauncherBin()` (os.Executable→EvalSymlinks). AC-3: `TestDetectInstallSource` classifies 6 real-shaped paths (edge/stable Caskroom, brew-stripped→HostOnly, checkout→NonBrew, unknown-token→SourceUnknown, ""→SourceUnknown) — fails if token/HostOnly logic misreads any. `Compare` untouched → `contract_test.go` (15 Compare cases) passes UNEDITED. AC-4: `TestTooOldBinaryRemedyPerSource` byte-for-byte-equals the old 3-line block for `InstallSource{}` (fails on any drift); `TestTooOldBinaryRemedyLeadsWithBrew` keeps its assertions, only its call is reconciled.
- DONE (tolerance sub-claim breached): Detection + rendering table tests + the command-reference doc note; gofmt clean, go test ./... and -race green; within ~9 files / tolerance +-3 files, +-60 LOC.
  Doc note added to `docs/site/reference/command-reference.md` (AC-5). `gofmt -l` clean; `go vet ./...` clean; `go test ./...` and `go test ./... -race` both green (every package `ok`). Files: 10 (within ±3 of ~9). LOC OVER the ±60 tolerance: committed numstat vs merge-base is +416/-32 (net +384) vs the ~298 estimate — ~+118 on added lines. Concentrated in test tables (`install_source_test.go` 198 vs est 115: the AC-1 four-source RunDoctor matrix + AC-3 six-row detection table + skip-guarded live smoke, at the repo's dense-comment norm); prod +35 is the `currentInstallSource` DRY helper shared by the 4 call sites plus doc-comment updates. No AC scope grew; trimming is a cheap follow-up if the gate wants it, but it would cut AC-mandated coverage or fight the surrounding comment density.

### Summary

Threaded a detected `contract.InstallSource` through the single remedy generator so the too-old-binary message names the upgrade path that actually applies — the captain's `@next`-user regression (`brew upgrade spacedock` no-op) is fixed, and non-brew/sandbox users get a source-rebuild / run-on-host hint instead of an inapplicable brew line. Kept `contract` pure and `Compare`'s signature stable (zero-value `InstallSource` = SourceUnknown = today's block byte-for-byte), capping the caller ripple to the `ManifestVerdict`/`RunDoctor` path; detection reads the Caskroom token off the resolved running binary with no `brew` subprocess. The live cask smoke ran (not skipped) on this Caskroom box, grounding the value check. One deviation to flag at the gate: total LOC exceeds the ±60 tolerance (+416/-32 vs ~298), driven by AC-mandated test tables and style-consistent comments, not scope growth.

## Gate: validation — PASSED (recommendation)

## Stage Report: validation

- DONE: AC-1 (per-source remedy matrix, oracle = source Kind external to the message)
  `TestTooOldBinaryRemedyRendersPerSource` (cli) drives 4 sources through `contract.RunDoctor`; oracle is the source Kind input, moves each row independently. Green. Mutation A (force BrewEdge→`spacedock`) reds it (`must contain "brew upgrade spacedock@next"`) — not tautological.
- DONE: AC-2 (word-boundary: NO bare `brew upgrade spacedock` for @next/non-brew/host-only; `spacedock@next` present in edge row)
  Both cli render table and contract `TestTooOldBinaryRemedyPerSource` assert regexp `brew upgrade spacedock(\s|$)` absent on edge/non-brew/host-only, present on stable. Probed the regexp for holes: `brew upgrade spacedock@next` (and `…@next\n…`) does NOT match; bare stable at `\n` and EOF DOES match. No hole.
- DONE: AC-3 (detectInstallSource classifies the 6 real-shaped paths)
  `TestDetectInstallSource` 6 rows green; oracle = input path + brew stub. Mutation B (token off-by-one → returns "Caskroom") reds all 3 brew rows; Mutation C (exact-switch → HasPrefix) reds the `spacedock@beta` unknown-token row. Probes: `spacedock-foo`→SourceUnknown, trailing-slash→BrewEdge, substring-non-Caskroom→NonBrew, `Caskroom-fake`→NonBrew — no misclassification.
- DONE: AC-4 (SourceUnknown reproduces today's 3-line block byte-for-byte AND Compare's cases pass UNEDITED)
  `source-unknown` subtest pins the exact 3-line string; `contract_test.go` has an EMPTY diff vs merge-base and its Compare-based tests (11 named `TestCompare` subtests + `TestCompareMessageShape`/`TestCompareHostSubstitution`) pass unedited. `TestTooOldBinaryRemedyLeadsWithBrew` reconciled call-only, green. (Note: implementer's "15 Compare cases" label = actually 11 named subtests + 2 siblings; immaterial — the unedited-and-green invariant holds.)
- DONE: AC-5 (doc note)
  `docs/site/reference/command-reference.md:56` documents the binary-out-of-date remedy as install-source-aware (brew stable / `@next` / source rebuild / run-on-host sandbox). Prose, reviewed here.
- DONE: Detached adversarial audit on a THROWAWAY detached git worktree (never the impl worktree; torn down after)
  Threaded surface probed via `frontdoor.go gateHost` path (`currentInstallSource`). 3 claim-breaking mutations each reddened a guard (A: edge→bare-stable reds 5 guards incl. live smoke; B: token parse reds 3 detection rows; C: prefix-match reds unknown-token). No AC coverage faked/tautological.
- DONE: Semantic adversarial pass (os.Executable anchor, EvalSymlinks, missing/empty Caskroom, brew-stripped HostOnly)
  os.Executable-vs-PATH anchor proven LIVE: built the worktree binary (a non-Caskroom source build) and ran its REAL `doctor --plugin-manifest` — it emitted the NonBrew remedy (`Rebuild from source … or re-download`, NO brew line), NOT the SourceUnknown generic block (which carries `brew upgrade spacedock`). So `currentInstallSource()` threads a DETECTED (non-zero) source end-to-end through the real front-door/doctor composition, anchored on the running binary — not hardcoded `InstallSource{}`, not the PATH `@next` binary. EvalSymlinks handled by pre-existing `resolvedLauncherBin`; missing/empty/HostOnly all table-covered + probed.
- DONE: Test hygiene — gofmt/vet/build clean; `go test ./...` green; race on touched packages green
  `gofmt -l` (changed files) clean; `go vet ./...` clean; `go build ./...` clean; full `go test ./...` green; `go test -race ./internal/cli ./internal/contract` green. NOTE: an initial `go test ./...` hit `ENOSPC: no space left on device` (disk at 100%, 516Mi free) failing 4 tests (`TestUpgradeFromStaleMovesToGreen`, `TestFreshBoxInstallSucceeds`, `TestNewVerbMintsInDiscoveredWorkflow`, `TestNewVerbFolderForm`) — 3 untouched by this change, the 4th failed past its 1-line param edit at the `claude plugin install` step. Cleared the regenerable go-build cache (freed 2.4Gi); all 4 then PASS. Environmental, not a code regression.

### Summary

PASSED. All five value ACs verified with reproduced evidence and, where an automated guard exists, mutation-tested (each claim-breaking edit reds a specific test) — no AC proof is self-referential or tautological. The captain's regression is fixed at the value level: BrewEdge→`brew upgrade spacedock@next` (render table + live smoke on this real `@next` Caskroom box), non-brew→source rebuild with no brew line, sandbox→run-on-host, SourceUnknown→byte-for-byte the old block. The FO-flagged high-stakes surface (`gateHost`/`currentInstallSource` threading) is confirmed live: the built binary's real `doctor` path emits a detected NonBrew remedy, proving the source flows end-to-end rather than degrading to the zero-value generic block. The `Compare` signature stayed stable (contract_test.go empty diff). No material findings. Deferred/observational only: (1) the documented, out-of-scope known limitation — a sandbox that RELOCATES the binary off its Caskroom path AND strips brew degrades to NonBrew/SourceUnknown (trigger outside the promised in-place-sandbox case; HostOnly handles the common case); (2) BrewEdge cannot be driven through the *built* binary end-to-end on this box (the only Caskroom binary is the pre-fix installed one) — covered instead by the union of exercised links (real-@next-path detection + render(BrewEdge) + live currentInstallSource→runDoctor threading). The FO-accepted +118 LOC overage is not re-litigated; confirmed it is AC-mandated test-table volume, not faked coverage.

## Stage Report: implementation (cycle 2 — captain @next-only trim)

Supersedes the four-source implementation above. The captain trimmed the scope (2026-07-22) to ONLY the @next edge-channel fix; the branch was reset to clean `origin/main` (`dceede3a`; main advanced past the earlier `ca136f83` base) and the four-source code discarded. Items track the trim directive, not the superseded four-source checklist.

- DONE: @next-only remedy — when the running binary is the `spacedock@next` cask, `tooOldBinaryRemedy` emits `brew upgrade spacedock@next`; otherwise the 3-line block is byte-for-byte unchanged.
  Threaded a small `edgeCask bool` (NOT a `SourceKind` enum) through `compareNamed`/`ManifestVerdict`/`RunDoctor`; `Compare` passes `false`. Code commit `37a872e2`. AC-1: `TestTooOldBinaryRemedyEdgeChannel` asserts `tooOldBinaryRemedy(true)` contains `brew upgrade spacedock@next` AND regexp `` `brew upgrade spacedock(\s|$)` `` does NOT match — reds if the formula switch is dropped (edge falls back to bare stable). AC-2: same test pins `tooOldBinaryRemedy(false)` byte-for-byte equal to the old block — reds on any drift.
- DONE: Minimal detection — `isEdgeCaskPath` reads the resolved `os.Executable` Caskroom token for the @next case only; no enum, no non-brew/sandbox/HostOnly.
  New `internal/cli/edge_cask.go` (`runningEdgeCask`/`isEdgeCaskPath`), anchored on `resolvedLauncherBin()`. AC-3: `TestIsEdgeCaskPath` — edge Caskroom path→true; stable cask / source checkout / empty→false. Reds if the segment-after-`Caskroom` check misreads (substring match or off-by-one).
- DONE: Entity body revised to @next-only + trimmed title proposed.
  Problem / Proposed approach / Out of scope / Acceptance criteria / Test plan / Expected surface rewritten to the @next-only scope with a Captain-TRIM provenance banner; the four-source ideation/impl/validation reports kept as history. **Proposed title (frontmatter — for the FO to apply, since ensigns do not edit frontmatter):** `Binary-upgrade prompt must name spacedock@next for edge-cask installs`.
- DONE: TDD; gofmt/vet clean; `go test ./...` and `-race` green.
  `gofmt -l` clean, `go vet ./...` clean. `go test ./...` green (0 FAIL); `go test ./... -race` green (0 FAIL). NOTE: one earlier `-race` run hit a TRANSIENT `internal/contractlint` failure — its repo-root filesystem walk raced a concurrent goreleaser `dist/` snapshot build mutating the git-ignored `dist/`; environmental, unrelated to this change; the settled re-run is clean. (Latent robustness gap in that boundary-guard test — it walks git-ignored `dist/` — noted, not fixed here: unrelated package.) As-built surface: 9 files, +134/-25 (net +109) — above the ~30-60 target, driven by the 4-signature threading ripple + a new detection file + tests at repo comment density, NOT behavioral scope (one formula line). ~⅓ the surface of the discarded four-source build (+416/-32).

### Summary

Minimal @next-only fix per the captain's trim: a single `edgeCask bool` threaded into the one remedy generator emits `brew upgrade spacedock@next` when the running binary is the edge cask (resolved `Caskroom/spacedock@next/` token), leaving every other install's remedy byte-for-byte unchanged. No `SourceKind` enum, no non-brew/sandbox arms — those are now Out of scope. Before/after for the @next case: `brew upgrade spacedock` (a no-op for @next users — the bug) → `brew upgrade spacedock@next`. Entity spec sections revised to the trimmed scope; four-source history preserved; trimmed title proposed for the FO to apply. `go test ./...` and `-race` green. Surface is net +109 (above the ~30-60 target but ~⅓ of the discarded build) — the plumbing to keep `contract` pure, not behavioral scope.
