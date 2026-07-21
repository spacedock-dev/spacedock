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
---

When the skill is upgraded and determines the binary needs upgrading, the remedy it emits hardcodes `brew upgrade spacedock` — even when that is not how the binary was installed. It must detect the real install source and runtime context and emit the correct upgrade instruction. Likely a follow-up gap in the shipped upgrade-hint work (`install-refresh-and-upgrade-hint`, `init-upgrade-and-contract-remedy`).

## Problem

The too-old-binary remedy is generated in ONE place: `internal/contract/contract.go` → `tooOldBinaryRemedy()` (called only from `compareNamed`, reached by `spacedock doctor`, `spacedock init/install`'s post-install doctor, and the `spacedock claude|codex` front-door version gate via `ManifestVerdict`). It returns a fixed block that leads with `brew upgrade spacedock` regardless of how the running binary was installed:

```
  Upgrade via Homebrew: brew upgrade spacedock
  Or build from source: go build -o spacedock ./cmd/spacedock
  Or refresh the plugin instead: spacedock install
```

That stable-`spacedock` command is wrong for the source the binary actually came from, and leaves the user stuck on the old binary:

1. **Edge / `@next` channel (captain-observed).** A user who installed via the `spacedock@next` cask ran the prompt, followed `brew upgrade spacedock`, and it did nothing — the stable cask isn't installed. The edge cask is a *separate* cask token (`spacedock@next`, `conflicts_with cask: spacedock`), so the correct command is `brew upgrade spacedock@next`.
2. **Non-brew install.** A source `go build`, a `go install …@next` proxy build, a downloaded release archive, or a `SPACEDOCK_BIN` pointing at a checkout build — `brew upgrade` does not apply at all; the user needs a rebuild / re-download.
3. **Sandbox / brew-unreachable execution.** When the doctor/gate runs where Homebrew is not on PATH (a safehouse sandbox, a minimal CI/container env), any `brew upgrade …` line is unactionable in place — the upgrade must be run on the host.

No skill/agent text and no user-facing doc restates the `brew upgrade spacedock` remedy (grep of `skills/`, `agents/`, `docs/site` finds only the too-old-*plugin* `spacedock install` remedy, and `install.md`'s `brew install` setup line). So the fix is contained to the code remedy generator plus one small doc addition — there is no scattered skill abort text to chase.

## Proposed approach

Make the remedy generator install-source-aware by passing it a detected install source, mirroring how `tooOldPluginRemedy(host)` already threads the detected host.

**Detection lives in `internal/cli` (new `install_source.go`); remedy selection lives in `internal/contract` (pure data-in).** The `contract` package stays free of `os.Executable`/`exec`/`safehouse` — it receives a small value and switches on it.

1. **New value in `internal/contract`:**
   ```go
   type SourceKind int
   const ( SourceUnknown SourceKind = iota; BrewStable; BrewEdge; NonBrew )
   type InstallSource struct { Kind SourceKind; HostOnly bool } // HostOnly: brew unreachable (sandbox/minimal)
   ```
   Zero value `{SourceUnknown, false}` reproduces today's 3-line block — the safe fallback and the reason the public `Compare` (17 callers) stays untouched.

2. **New detector `detectInstallSource(execPath string, brewLookPath func(string)(string,error), devBranch string) contract.InstallSource` in `internal/cli`.** All inputs injectable for tests. Algorithm:
   - Resolve the RUNNING binary via the existing `resolvedLauncherBin()` (`os.Executable` → `filepath.EvalSymlinks`). Anchor on the running process, not `command -v spacedock` (they can differ — see Spike).
   - `token, isCask := caskToken(execPath)`: split the resolved path on `/`, and if a `Caskroom` segment is present, `token` is the next segment. `token == "spacedock@next"` → `BrewEdge`; `token == "spacedock"` → `BrewStable`; other/absent-token cask path → `SourceUnknown`.
   - For a brew kind, `HostOnly = brewLookPath("brew") errored` (brew not resolvable here → run-on-host variant; the sandbox case).
   - Non-empty resolved path with NO `Caskroom` segment → `NonBrew`.
   - Empty/unresolvable path → `SourceUnknown` (generic block).

3. **Thread the source through the remedy path only:** `compareNamed` gains a `src InstallSource` param; `Compare` keeps its signature and passes `SourceUnknown`; `ManifestVerdict(manifestPath, host, binaryVersion, src)` and `RunDoctor(…, src, …)` gain the param. `tooOldBinaryRemedy(src)` switches on it. Callers computing the source: `frontdoor.go` `gateHost` (×1) and `init.go` `runDoctor`/codex/manifest arms (×3), each calling `detectInstallSource(...)`.

4. **Per-source remedy wording** (every arm keeps the binary-vs-plugin distinction line `Or refresh the plugin instead: spacedock install`, which the existing test pins):
   - `SourceUnknown` → unchanged 3-line block (brew upgrade spacedock / build from source / refresh plugin).
   - `BrewStable` → `  Upgrade via Homebrew: brew upgrade spacedock` + refresh line.
   - `BrewEdge` → `  Upgrade via Homebrew: brew upgrade spacedock@next` + refresh line.
   - `BrewStable`/`BrewEdge` with `HostOnly` → `  Homebrew isn't reachable here (e.g. a sandbox). Upgrade on your host, then relaunch: brew upgrade <formula>` + refresh line (formula per channel).
   - `NonBrew` → `  Rebuild from source: go build -o spacedock ./cmd/spacedock (or re-download the latest release)` + refresh line. No `brew` line.

**Why not simpler alternatives?**
- *Use the `devBranch` build stamp (`next`/`main`) alone to pick the formula.* Rejected as the sole signal: `devBranch` defaults to `next` for any plain `go build`, so it can't distinguish an edge cask from a source build — it would emit `brew upgrade spacedock@next` for a source checkout. The Caskroom-token path signal is what proves a brew install; `devBranch` is not needed once the path is read.
- *Shell out to `brew list --cask`/`brew --prefix` to find the owning cask.* Rejected: slower, adds a subprocess dependency, and fails inside a sandbox where `brew` is stripped — the exact case we must handle. The resolved Caskroom path already carries the token with no subprocess.
- *Overwrite the remedy string in the CLI after calling the verdict.* Rejected: duplicates message assembly and splits remedy ownership across two packages; threading a value mirrors the existing `host` threading and keeps assembly in `contract`.

## Out of scope

- **Version-COMPARISON logic.** What counts as out-of-date (`ParseMajorMinor`, `Compare`, the verdict classification) is untouched. Only the too-old-binary remedy MESSAGE changes.
- **The other verdicts' messages.** Compatible, too-old-plugin, malformed-version, no-plugin-found are byte-identical to before.
- **Auto-upgrade / self-update.** The remedy tells the user what to run; it never runs it.
- **New install methods** or changes to the cask/goreleaser channel model.
- **A new `safehouse` in-process "am I inside a sandbox" marker.** None exists today, and adding one to safehouse (a separate project) is out of scope. The sandbox case is detected by its observable consequence — `brew` unresolvable on PATH — via `HostOnly`, not a marker. Known limitation (documented, not fixed here): a sandbox that *relocates* the binary off its Caskroom path AND strips `brew` degrades to `NonBrew`/`SourceUnknown`; the common in-place sandbox (same path, brew stripped) is handled by `HostOnly`.
- **`SPACEDOCK_DEV_BRANCH` / channel overrides changing the emitted formula.** The formula is read from the Caskroom token (the real cask), not the overridable `devBranch` stamp.

## Acceptance criteria

- **AC-1 (value — remedy matches the actual source across the matrix).** For each install source in {`BrewStable`, `BrewEdge`, `NonBrew`, brew+`HostOnly` (sandbox)}, the too-old-binary remedy emitted through the doctor/gate path names the instruction that actually applies: brew-stable → `brew upgrade spacedock`; brew-edge → `brew upgrade spacedock@next`; non-brew → a source rebuild / re-download instruction with NO `brew` line; sandbox → a run-on-host instruction. Verified by a table test whose oracle is the install-source INPUT, which lives outside the message and moves each row independently — a remedy generator that ignored its input would fail ≥3 rows.
  - *Test:* `internal/cli/install_source_test.go` rendering table (source → emitted remedy tokens).
- **AC-2 (the captain's regression — no bare stable command where it's wrong).** For a `BrewEdge` install and for a `NonBrew` install, the emitted remedy contains NO bare stable `brew upgrade spacedock` command, word-boundary matched so that `brew upgrade spacedock@next` does NOT satisfy the stable substring. (This is exactly the `@next`-user-got-`brew upgrade spacedock` bug.)
  - *Test:* same table asserts `regexp \`brew upgrade spacedock(\s|$)\`` does NOT match the edge/non-brew rows, while `brew upgrade spacedock@next` IS present in the edge row.
- **AC-3 (detection grounded in the real path).** `detectInstallSource` classifies the resolved running-binary path: a `…/Caskroom/<token>/…` path → the brew kind named by `<token>` (`spacedock`→stable, `spacedock@next`→edge); a resolved non-Caskroom path → non-brew; a brew-owned path with `brew` unresolvable → the `HostOnly` variant. Anchored on the spike's observed real path (below).
  - *Test:* `internal/cli/install_source_test.go` detection table over `(execPath, brewLookPath stub, devBranch)` inputs.
- **AC-4 (no collateral change).** The `SourceUnknown` fallback reproduces today's 3-line remedy byte-for-byte, and the compatible / too-old-plugin / malformed / no-plugin messages are unchanged — so the untouched `Compare`-based verdict tests pass without edits.
  - *Test:* existing `internal/contract/contract_test.go` + `TestTooOldBinaryRemedyLeadsWithBrew` (drives the `SourceUnknown` arm) stay green unmodified.
- **AC-5 (docs).** The command-reference doctor section documents the binary-out-of-date remedy as install-source-aware, alongside the existing plugin-out-of-date remedy.
  - *Test:* the doc diff below is applied; reviewed at the gate (prose, no automated check).

## Test plan

Primary proof is Go unit/table tests exercising the generator and detector by feeding inputs and asserting the emitted bytes — no mocks, oracle external to the message.

1. **`internal/cli/install_source_test.go` (new, ~110 LOC).**
   - *Detection table* over `detectInstallSource(execPath, brewLookPath, devBranch)`:
     - `/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock`, brew-found → `{BrewEdge, HostOnly:false}`.
     - `/usr/local/Caskroom/spacedock/0.25.0/spacedock`, brew-found → `{BrewStable, false}`.
     - same Caskroom path, brew-NOT-found (stub errors) → `{BrewEdge, HostOnly:true}`.
     - `/Users/x/git/spacedock/spacedock` (checkout build) → `{NonBrew, false}`.
     - `""` (unresolvable) → `{SourceUnknown, false}`.
   - *Rendering table*: for each `contract.InstallSource`, drive a too-old fixture through `contract.RunDoctor(fixture, "claude", oldBinaryVersion, src, &out, &err)` and assert exact/forbidden tokens per AC-1/AC-2, using the word-boundary regexp for the stable command.
2. **`internal/contract/version_message_test.go` (reconcile + extend, ~+55 LOC).** `TestTooOldBinaryRemedyLeadsWithBrew` keeps driving the `SourceUnknown` arm (still `brew upgrade spacedock`, still no `@next`, its existing `@next`-forbid assertion stays valid for that arm). ADD `TestTooOldBinaryRemedyPerSource` covering `BrewEdge` (`brew upgrade spacedock@next`, no bare stable), `NonBrew` (no `brew`), and `HostOnly` (run-on-host, brew not the in-place action).
3. **Signature-ripple caller updates (~+8 LOC total):** `internal/contract/doctor_test.go` (×2) and `internal/cli/upgrade_from_stale_test.go` (×2) pass `contract.SourceUnknown` (verdict-only tests); `frontdoor.go` (×1) and `init.go` (×3) pass `detectInstallSource(...)`. `Compare` callers (`contract_test.go` ×15, `release/stamp_then_gate_test.go`, `skills/integration/contract_skew_test.go`) are UNCHANGED.
4. **Live smoke (skip-guarded, machine-grounded value check).** Resolve the `spacedock` on PATH; if it resolves under a `Caskroom/` segment, run `spacedock doctor --plugin-manifest <too-old fixture>` and assert the emitted remedy names that Caskroom token's formula (on a `@next` box: `brew upgrade spacedock@next`). Skips when the PATH binary isn't a Caskroom install (portable per "no hidden machine dependencies").
5. **Hygiene:** `gofmt -l` (clean), `go vet ./...`, `go test ./...`, `go test ./... -race`.

## Expected surface

| File | Kind | ~LOC |
|---|---|---|
| `internal/contract/contract.go` | prod — `InstallSource`/`SourceKind`, `tooOldBinaryRemedy(src)`, thread `compareNamed` | +45 |
| `internal/contract/doctor.go` | prod — `ManifestVerdict`/`RunDoctor` gain `src` | +5 |
| `internal/cli/install_source.go` | prod — `detectInstallSource`, `caskToken` (new file) | +55 |
| `internal/cli/frontdoor.go` | prod — `gateHost` computes + passes source | +4 |
| `internal/cli/init.go` | prod — 3 doctor callers pass source | +8 |
| `internal/cli/install_source_test.go` | test — detection + rendering tables + live smoke (new) | +115 |
| `internal/contract/version_message_test.go` | test — reconcile + per-source | +55 |
| `internal/contract/doctor_test.go`, `internal/cli/upgrade_from_stale_test.go` | test — ripple `SourceUnknown` | +8 |
| `docs/site/reference/command-reference.md` | doc — binary-out-of-date remedy note | +3 |

Estimate: ~9 files; prod ~117 LOC, test ~178 LOC, doc ~3. **Tolerance: ±3 files, ±60 LOC.** The dominant variable is the signature-ripple caller count; keeping `Compare` stable caps it (only `ManifestVerdict`/`RunDoctor` callers move).

## Doc diff

`docs/site/reference/command-reference.md`, the `doctor`/`init` paragraph (currently ~line 56). Append after the existing plugin-out-of-date sentence:

> When `doctor` reports the **binary** is out of date, it prints the upgrade path that matches how the binary was installed — `brew upgrade spacedock` for a stable Homebrew install, `brew upgrade spacedock@next` for the edge (`@next`) channel, a source rebuild for a non-Homebrew build, and a run-on-host hint when Homebrew isn't reachable (e.g. inside a sandbox).

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
