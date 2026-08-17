---
id: c9qz40bza726q58m3s8hsdq5
title: Prerelease releases ship a stable-stamped binary under the default asset name
status: validation
source: "Live cross-host report, CL 2026-08-17. A Linux box (linux/amd64) running `spacedock 0.27.0-pre7` reinstalled the Claude plugin as `spacedock@spacedock` on every launch, overwriting a hand-installed edge plugin. `go version -m $(command -v spacedock)` on that binary returned `devBranch=main`; the darwin/arm64 edge-cask binary at the identical version returned `devBranch=next`. Same tag, same `--version` output, opposite channel."
started: 2026-08-17T19:38:45Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-prerelease-ships-stable-stamped-default-artifact
issue:
mod-block: merge:pr-merge
pr: "#726"
---

A prerelease tag publishes both channel builds and gives the stable-stamped one the unsuffixed asset name, so installing edge from a `-pre` release silently yields a stable binary — and nothing in the binary's own output reveals which one you got.

## Problem

`.goreleaser.yaml` defines two builds over one shared goos/goarch matrix: `spacedock-stable` stamping `-X ...cli.devBranch=main` (line 49) and `spacedock-edge` stamping `devBranch=next` (line 63). Both stamp `cli.Version={{ .Version }}` identically. The archive templates (lines 72-85) name them `spacedock_{version}_{os}_{arch}.tar.gz` and `spacedock_{version}_{os}_{arch}_edge.tar.gz`.

So tag `v0.27.0-pre7` publishes both a stable-stamped and an edge-stamped linux/amd64 tarball into one prerelease, and the *edge* build is the one carrying a distinguishing suffix. A user installing edge reaches for the default-looking unsuffixed asset and receives `devBranch=main`.

That binary then drives the whole channel surface as stable: `channelMarketplace` returns `spacedock`, `channelMarketplaceSource` returns the bare repo, `channelPluginID` returns `spacedock@spacedock` (internal/cli/host_exec.go:236,246,265). The frontdoor re-ensures that id on every launch — correct behavior for a stable binary, which is why the symptom reads as a frontdoor bug and is not one.

The landing spot compounds it: marketplace `spacedock`'s `spacedock` entry pins git ref `stable` at v0.20.0, while the 0.27 binary's version gate requires plugin minor 0.27. The install therefore lands a plugin the binary itself then rejects.

Two things kept this undiagnosable. `internal/release/channel_agreement_guard_test.go` guards the stable/edge stamp pairing by parsing `.goreleaser.yaml` — it validates the CONFIG and never inspects a built artifact, so it passes while releases hand users the wrong binary. And no command prints the channel: `--version` reports version, OS, runtime, sandbox and contract, so two artifacts with opposite plugin behavior emit identical text. The stamp was recoverable only out-of-band via `go version -m`.

## Proposed approach

**Naming decision: on a real prerelease tag, no asset carries the default name.** Both channel builds stay, but the stable archive's `name_template` gains a conditional suffix:

```
name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if and .Prerelease (not .IsSnapshot) }}_stable{{ end }}"
```

A `-pre` release then publishes `spacedock_{v}_{os}_{arch}_stable.tar.gz` and `..._edge.tar.gz` and nothing under the bare name; a stable tag and a snapshot keep today's names byte-identical. Both prerelease assets self-describe their channel, so the incident's reach — "the default-looking asset" — ceases to exist.

Rejected directions:

- **Suppress the stable build on `-pre` tags** (the cleaner-sounding candidate): spiked and blocked. `builds[].skip` templates correctly, and the archive/checksum pipes tolerate the empty set, but the stable cask pipe hard-errors `no linux/macos archives found matching ... ids=[spacedock-stable]` before `skip_upload: auto` applies, and OSS goreleaser 2.16 casks expose no `disable` field (schema-verified). Every `-pre` cut would red mid-release.
- **Suffix both channels on every release:** invalidates the stable install surface (install.sh's constructed asset name, every published cask URL shape) to fix a prerelease-only trap, and forces a transition where old releases and new install.sh disagree about names.
- **Make install.sh/tap channel-explicit:** neither was involved. install.sh resolves `/releases/latest`, which GitHub never points at a prerelease, and the stable cask already skips `-pre` uploads. The trap is the hand-download from the release page, which only naming fixes.

Nothing consumes a prerelease's `_stable` tarball (the stable cask skips prerelease upload, install.sh cannot reach it, and the auto-pre0 catch-up run builds its own artifacts). It continues to exist only because the cask pipe requires its build; the suffix marks it honestly.

**Verification is descoped by captain ruling (ideation gate rejection, 2026-08-17).** Cycle 1 proposed three artifact-verification layers here — a `verify-stamp` build hook, a `verify-channel-stamps` dist sweep wired into two CI lanes, and a template-execution extension of the config guard: about 477 lines of verification for a five-line defect. CL rejected the kind, not just the size. The evidence is this incident itself: `internal/release/channel_agreement_guard_test.go` passed, the release pipeline was green, and a real user still downloaded a stable-stamped binary from a prerelease page and lost a day to it. A green suite coexisted with the defect; more tests of the same kind do not correct that. All three layers are cut and NOT replaced with smaller tests of the same kind. What survives as the guard is the one thing a human reads during actual usage — the `Channel:` line below — plus a hand-exercised check of the first real `-pre` cut (see Test plan).

**Legibility (AC-2):** `--version` gains one line in BOTH output shapes, directly after `OS:`:

```
Channel: stable (spacedock@spacedock)
Channel: edge (spacedock@spacedock-edge)
```

The channel word maps from the effective devBranch exactly as `channelMarketplace` does (main → stable, anything else → edge); the parenthetical is `channelPluginID(devBranch)` — the id the frontdoor re-ensures, so the incident's symptom (the `spacedock@spacedock` reinstall loop) and its cause read off one line. This line is NOT a test: it prints on every run, so a human reads it during actual usage, and its absence is why this defect cost an investigation instead of one command. It reports the EFFECTIVE value: the `--version` path calls `applyDevBranchOverride(env)` before `printVersion` — the same helper every install path already calls, so `SPACEDOCK_DEV_BRANCH` override semantics (including explicit-empty) are exact and the line renders what the binary will actually do; `printVersion` then reads the package `devBranch` var directly, a two-way branch mirroring `channelMarketplace` plus `channelPluginID(devBranch)`. Line 1 stays the bare version token; consumers audited: the FO version gate parses line 1 only, fo-install-gate greps `^Sandbox: ` (prefix-anchored), and the frozen `contract 3` token is position-independent below line 1.

Mechanism → value audit (the design has exactly two mechanisms; simplest alternative for each):

- Conditional `_stable` suffix (serves AC-1). Alternative: skip the stable build — blocked by the cask pipe (spiked). Alternative: docs-only warning — leaves the trap in place.
- `Channel:` line (serves AC-2). Alternative: keep `go version -m` as the channel oracle — precisely the out-of-band-only state that made the incident an investigation instead of one command.

No further spike is needed: the kept mechanisms rest on the already-spiked conditional suffix (Spike record 2-3) and a print statement over existing fixture-tested surfaces (`printVersion`, `applyDevBranchOverride`, `channelPluginID`).

## Spike record

Spiked 2026-08-17 with goreleaser 2.16.0: full `goreleaser release` runs in a scratch clone against throwaway local tags `v0.99.0-pre1` / `v0.99.0`, cask pushes aimed at a local bare repo via `repository.git.url`, `release.disable: true` plus an explicit cask `url.template` standing in for the GitHub API. Facts established:

1. `builds[].skip` accepts a template (schema: string|bool) and fires on a `-pre` tag (`skip is set id=spacedock-stable`), and the archive/checksum pipes tolerate the resulting empty stable set — but the stable cask pipe errors `no linux/macos archives found matching goos=[darwin linux] goarch=[amd64 arm64] ids=[spacedock-stable]` before `skip_upload: auto` is consulted, and casks have no `disable` field. The build-skip direction is dead on OSS goreleaser.
2. The conditional name suffix works end to end: the `-pre` tag produced only `_stable` + `_edge` archives (checksums.txt matching), both casks generated, and the stable cask publish skipped itself with `prerelease detected with 'auto' upload`. The stable tag produced today's unsuffixed + `_edge` names with cask URLs referencing the unsuffixed asset — the tap contract is byte-stable.
3. The `not .IsSnapshot` guard is load-bearing: without it a snapshot names the stable asset `..._stable.tar.gz` (snapshot versions carry a semver prerelease part) and install-e2e's end-anchored install.sh glob finds nothing.
4. `go version -m` / `debug/buildinfo` reads `-X ...cli.devBranch=...` from cross-compiled tarball binaries on a foreign OS (linux ELF read on darwin) — the read mechanism the cut verification layers used, identical to the incident's own forensic read.
5. `builds[].hooks.post` runs once per built artifact with `{{ .Path }}` resolved; a hook asserting the wrong stamp fails the whole release at the build stage (`post hook failed: exit status 1`) — the pre-publish abort direction proven, not assumed.

Items 4-5 established mechanisms for the verification layers cut at the ideation gate (captain ruling recorded in Proposed approach). The facts stand on record; nothing in the kept design consumes them.

## Migration for already-installed prerelease binaries

The naming change is forward-only: already-published `-pre` releases keep their unsuffixed stable assets, and an installed binary cannot change its stamp. For an operator holding one (the incident machine: linux/amd64, 0.27.0-pre7, devBranch=main):

- **Detect:** on existing binaries, `go version -m $(command -v spacedock)` showing `cli.devBranch=main` means stable; from this change on, `spacedock --version` prints the channel in-band. Symptom key: a prerelease-versioned binary that re-ensures `spacedock@spacedock` on every launch and overwrites an edge plugin is a stable-stamped binary, not a frontdoor bug.
- **Remedy:** replace the binary with the same release's `_edge` tarball or `brew install spacedock-dev/tap/spacedock@next`. No state migration is needed: the frontdoor re-ensures the channel plugin per launch, so plugin/marketplace state converges on the next launch once the binary's channel is right.
- **Optional retroactive closure (captain decision at the gate — destructive, published surface):** delete the unsuffixed stable assets from existing `-pre` releases (0.27.0-pre1..pre7 and earlier `-pre` cuts). Nothing sanctioned references them — the stable cask never published from a prerelease and install.sh cannot resolve one — and their checksums.txt lines dangle harmlessly (lookups are by filename). Until this is done, old prerelease pages still carry the trap.
- Stable-release installs are untouched: names, install.sh, and tap behavior stay byte-identical for `vX.Y.Z` tags.

## Out of scope

The duplicate edge routes in the marketplace (`spacedock-edge@spacedock` and `spacedock@spacedock-edge` resolve the same payload) — a published-surface decision, filed separately if CL wants it changed.

Whether `channelEntry`'s constant entry name (internal/cli/host_exec.go:223) makes an *edge* binary reinstall over a route-A install: plausible from reading, unverified, and not the cause of this report. Do not fold it in without first confirming the frontdoor's actual presence-check behavior.

The marketplace and re-pull work is task `2d`: `installArgvSequence`, the marketplace manifests, and `bump-calendar` are untouched here.

## Expected surface and tolerance

Estimate net LOC change: about +48/−8 across 5 files, most of it documentation (supersedes cycle 1's +560/11-file declaration, rejected at the gate; the captain-directed target is about 55 lines):

- `.goreleaser.yaml` +5 (stable archive name_template conditional plus comment) — the defect fix
- `internal/cli/cli.go` +15 (`applyDevBranchOverride(env)` on the `--version` path; Channel line in `printVersion`; comments)
- `internal/cli/version_session_test.go` +6 net (six exact-match `want` literals gain the Channel line; the two-lines/three-lines shape comments reworded in place) — forced maintenance, not new test surface
- `docs/releasing.md` +18/−3 (naming sentence in "What the Tag Push Does"; migration bullet in "Notes")
- `docs/site/reference/command-reference.md` +4/−2 (`--version` intro and three-line example)

The dispatch's "about 4 files" folds to 5 in practice because the migration note and the `--version` reference live in separate doc files; the line total holds. Tolerance: ±1 file, +25 LOC. Hard ceiling per captain direction: if the implementation passes 80 inserted lines, stop and bring it back to the gate with the reason stated.

## Declared semantic changes

- **Command grammar:** `--version` gains a third line `Channel: <stable|edge> (<plugin-id>)` in BOTH shapes; the outside-session shape goes two lines → three. Known consumers audited (FO version gate: line 1 only; fo-install-gate: `^Sandbox: ` prefix grep; `contract 3`: position-independent). No new subcommands.
- **Published-release surface:** on `-pre` tags only, the stable tarball's asset name gains `_stable`; checksums.txt keys follow. Stable-tag and snapshot asset names stay byte-identical to today.
- **Stored formats:** none. **Authority:** none. **Runtime install behavior:** none — frontdoor, marketplace mapping, install.sh, and tap flow untouched.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Verification is deliberately thin by captain ruling: the incident proved that a green suite of artifact-agreement tests coexists with users receiving the wrong binary, so the evidence of record for this task is the first real `-pre` cut, exercised by hand and recorded in the entity body at validation.

**AC-1 - A prerelease release offers no default-named asset: every asset a `-pre` release publishes names its channel, so reaching for "the" asset is impossible.**
This is the measuring AC: the unsuffixed-archive count over the first post-landing `-pre` release's published asset list must be ZERO — a count over the real published surface that moves the wrong way the moment the template regresses. Verified by hand at validation: read that release's page; only `_stable` and `_edge` assets exist. Fails if `.goreleaser.yaml` gives `spacedock-stable` the bare `{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}` name while prerelease publishing is enabled.

**AC-2 - A running binary reports the channel it drives, in-band, so identifying a wrong-channel binary is one command instead of an investigation.**
Verified by: (fixture, forced maintenance) the six exact-match literals in `version_session_test.go` now include the Channel line, so dropping the line, moving it off position 3, or breaking the edge rendering turns `go test ./...` red; (live, at validation) download the `_edge` asset from the AC-1 cut, run `spacedock --version`, and the Channel line reads `edge (spacedock@spacedock-edge)`. The line renders the EFFECTIVE devBranch (override applied via the same helper the install paths use), so it reports what the binary will actually do. Fails if the line is dropped or diverges from the plugin id the frontdoor ensures. This AC serves AC-1's end value: the release being correct is the goal; this is what lets an operator confirm it rather than trust it.

## Test plan

Deliberately thin, by captain ruling at the cycle-1 gate: the artifact verifier, both CI lane wirings, and the config-guard template extension are cut and NOT replaced with smaller tests of the same kind. The evidence against that kind is this incident itself — `channel_agreement_guard_test.go` passed, the release pipeline was green, and a user still downloaded a stable-stamped binary from a prerelease page.

- **Fixture (forced maintenance):** the six exact-match `want` literals in `version_session_test.go` gain the Channel line after `OS:`; `go test ./...` green. These pin the line's presence, its position, and the edge rendering of the unstamped test default (`devBranch=next` → `Channel: edge (spacedock@spacedock-edge)`). The stable-side rendering (main → `stable (spacedock@spacedock)`) is a two-way branch mirroring `channelMarketplace`, reviewed at the gate and confirmed live the first time `--version` runs on any stable binary; no dedicated fixture is added for it, per the no-new-test-surface ruling. Stated residual: a hardcoded-to-edge regression would pass the fixtures — accepted, because the artifact-verification kind that would catch it is what the captain cut.
- **Live (the proof of record, recorded in the entity body at validation):** the first `-pre` cut after this lands, exercised by hand:
  1. Read the release page. Only `_stable` and `_edge` assets exist (AC-1's zero count).
  2. Download the `_edge` asset and run `spacedock --version`. The Channel line reads edge (AC-2).
- Estimated cost: near zero — no new test files, no CI changes; the live proof rides an already-planned prerelease cut.

## Documentation diff

install.sh and the tap need NO change: install.sh resolves `/releases/latest` (never a prerelease), its local-dist glob is end-anchored on `_{arch}.tar.gz` (a suffixed asset can never be mistaken for the default), snapshots keep the unsuffixed name (spiked), and casks are regenerated per release from actual asset names with the stable cask never publishing from a `-pre` tag (spiked). The two docs that describe the changed surfaces:

**docs/releasing.md, "What the Tag Push Does" first bullet — before:**

```
- cross-builds darwin and linux (arm64 + amd64) tarballs plus `checksums.txt`,
  stamping `git describe --tags` into `internal/cli.Version`, for BOTH channels —
  a stable build (`cli.devBranch=main`) and an edge build (`cli.devBranch=next`);
```

**after:**

```
- cross-builds darwin and linux (arm64 + amd64) tarballs plus `checksums.txt`,
  stamping `git describe --tags` into `internal/cli.Version`, for BOTH channels —
  a stable build (`cli.devBranch=main`) and an edge build (`cli.devBranch=next`).
  The asset name carries the channel: the edge tarball always ends `_edge`; the
  stable tarball is unsuffixed on a stable tag but ends `_stable` on a `-pre`
  tag, so a prerelease never offers a default-named asset. Nothing consumes a
  prerelease's `_stable` tarball — it exists only because the cask pipe cannot
  tolerate a skipped build;
```

**docs/releasing.md, "Notes" — add bullet:**

```
- Prerelease releases cut before the `_stable` suffix landed shipped the
  stable-stamped binary under the unsuffixed default asset name. An operator
  who installed one holds a stable-channel binary at a prerelease version: it
  re-ensures the `spacedock@spacedock` plugin on every launch. Detection:
  `spacedock --version` prints a `Channel:` line on new binaries; on older
  ones `go version -m $(command -v spacedock)` shows `cli.devBranch=main`.
  Remedy: reinstall from the release's `_edge` tarball or
  `brew install spacedock-dev/tap/spacedock@next`.
```

**docs/site/reference/command-reference.md, `--version` section — before:**

```
`spacedock --version` reports the binary version and the host OS/arch, and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints two lines:

    spacedock 0.26.0
    OS: darwin/arm64
```

**after (the in-session example block and the closing "the output is the two lines shown above" sentence update in step; the intro paragraph's parenthetical gains "and release channel"):**

```
`spacedock --version` reports the binary version, the host OS/arch, and the release channel the binary drives (`stable` installs `spacedock@spacedock`; `edge` installs `spacedock@spacedock-edge`), and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints three lines:

    spacedock 0.26.0
    OS: darwin/arm64
    Channel: stable (spacedock@spacedock)
```

**Release-notes snippet for the first suffixed `-pre` cut:**

```
Prerelease asset names now carry their channel: `_edge` (install this one) and
`_stable` (built for pipeline reasons; not for prerelease use). If you
previously installed the unsuffixed tarball from a -pre release you have a
stable-channel binary — `spacedock --version` now prints the channel;
reinstall from the `_edge` tarball or
`brew install spacedock-dev/tap/spacedock@next`.
```

## Stage Report: ideation

- DONE: Decide the prerelease asset-naming shape and state what happens to operators who already installed the unsuffixed prerelease binary — the naming decision is worthless without that migration story.
  Decision: conditional `_stable` suffix on `-pre` tags only (stable/snapshot names byte-identical); migration recorded in "Migration for already-installed prerelease binaries" — detect via `go version -m` (old) / `Channel:` line (new), reinstall `_edge` or `spacedock@next` cask, optional captain-gated asset deletion from published `-pre` releases.
- DONE: Design the artifact-level stamp check as a check that can FAIL: it must read a built binary's embedded ldflags, not parse .goreleaser.yaml, because config-parsing is exactly the blind spot that let this incident through a passing guard.
  `verify-stamp` (per-binary `builds[].hooks.post`, pre-publish) + `verify-channel-stamps` (dist sweep: buildinfo read per archive plus native-OS execution of the `Channel:` line); the FAIL direction exercised live in the spike — a wrong-stamp hook aborted `goreleaser release` at the build stage (`post hook failed: exit status 1`).
- DONE: Produce the concrete doc diff for any install instructions the naming change invalidates (install.sh, tap, docs/releasing.md), since this changes a user-visible install surface.
  install.sh and tap proven NOT invalidated (latest-release resolution, end-anchored glob, snapshot names preserved — spiked); before/after diffs recorded for docs/releasing.md and docs/site/reference/command-reference.md plus a release-notes snippet, in "Documentation diff".

### Summary

Chose the prerelease naming shape by spiking the riskiest mechanisms first with goreleaser 2.16.0 in a scratch clone: the build-skip direction is dead on OSS goreleaser (stable cask pipe errors on zero artifacts before `skip_upload: auto`; no cask `disable` field), while the conditional `_stable` name suffix passed end to end on `-pre`, stable, and snapshot runs — including proof that the `not .IsSnapshot` guard is load-bearing for install-e2e. Verification moves onto artifacts in three layers (pre-publish build hooks, a whole-set dist sweep wired into install-e2e and release.yml, and a template-executing config guard), with `--version` gaining a `Channel:` line that doubles as the sweep's execution oracle. Expected surface recalibrated from the seed's +140/5-file guess to +560/11-file with tolerance; semantic changes declared (version-output grammar, prerelease asset names; no runtime install behavior change).

## Stage Report: ideation (cycle 2)

- DONE: Cut the verifier tool, both CI lane wirings, and the config-guard extension. Do not replace them with a smaller test of the same kind — the captain's objection is to the kind, not the size.
  All references to `verify-stamp`/`verify-channel-stamps`, the install-e2e/release.yml wirings, and the guard-test template extension removed from approach, ACs, test plan, surface, semantics, and the releasing.md doc diff; no replacement test added, and the ruling plus its evidence (a green `channel_agreement_guard_test.go` through the incident) is recorded in the body so the kind stays cut.
- DONE: Keep only the asset-name conditional, the Channel line in --version, and the migration doc note. Add no new mechanism.
  Mechanism audit now lists exactly two mechanisms (suffix → AC-1, Channel line → AC-2); the Channel line spec gained one concrete detail — `applyDevBranchOverride(env)` on the `--version` path reuses the existing install-path helper, so effective-value semantics (including explicit-empty override) are exact with zero new helpers; Migration section and remaining doc diffs kept, trimmed of the verifier sentence.
- DONE: Pin the proof to a hand-exercised prerelease cut recorded at validation, and re-declare the surface at about 55 lines.
  Proof of record is now the first post-landing `-pre` cut checked by hand (release page shows only `_stable`+`_edge`; downloaded `_edge` binary prints `Channel: edge`), recorded in the entity body at validation; surface re-declared at ~+48/−8 across 5 files (5 not 4: the migration note and the `--version` reference live in separate doc files), tolerance +25 LOC with the captain's 80-line hard ceiling named.

### Summary

Descoped the design per captain rejection: the three artifact-verification layers (~477 lines) are cut without replacement, with the ruling and its incident evidence recorded in the body — the surviving guard is the human-read `Channel:` line plus a hand-exercised check of the first real `-pre` cut. Spike facts stand untouched (items 4-5 annotated as supporting cut layers); ACs reduced to the measuring zero-unsuffixed-assets count and in-band channel legibility; the fixture change is literal maintenance only, with the hardcoded-to-edge residual stated for the gate rather than hidden. Body-only ideation change on this entity (46 insertions / 43 deletions including this report); no code touched.

## Stage Report: implementation

- DONE: `.goreleaser.yaml` stable archive `name_template` gains the conditional `_stable` suffix on a real `-pre` tag, byte-identical on a stable tag/snapshot.
  `.goreleaser.yaml` (+8/-1); confirmed live with `goreleaser check` (config valid) and a local `goreleaser release --snapshot --clean` run — snapshot names stayed unsuffixed (`spacedock_..._darwin_arm64.tar.gz`), matching the already-spiked `not .IsSnapshot` guard.
- DONE: `--version` prints a `Channel: <stable|edge> (<plugin-id>)` line directly after `OS:`, in both output shapes, reading the EFFECTIVE devBranch via `applyDevBranchOverride(env)` on the `--version` path.
  `internal/cli/cli.go` (+23/-6): `applyDevBranchOverride(env)` added before `printVersion` in the `versionFlag` branch; `printVersion` prints the Channel line via a two-way branch mirroring `channelMarketplace`. Live-exercised: built real stable+edge snapshot binaries via `goreleaser release --snapshot`, ran `--version` on each (`Channel: stable (spacedock@spacedock)` / `Channel: edge (spacedock@spacedock-edge)`), then confirmed `SPACEDOCK_DEV_BRANCH=next|main|""` overrides each binary's printed Channel to the override's value, including the explicit-empty case (renders `edge`, since only `main` maps to `stable`).
- DONE: Update the six exact-match `want` literals in `version_session_test.go` to include the Channel line; no new test file, CI step, or verifier tool added.
  `internal/cli/version_session_test.go` (+14/-6): all six `TestVersionSessionRender` cases gain `Channel: edge (spacedock@spacedock-edge)\n` after the `OS:` line (devBranch stays package-default `next` in this file, so every case renders `edge`); shape comments reworded two-lines→three-lines / TWO→THREE in place. `go test ./... ` and `go test ./... -race` both green — a dropped or repositioned Channel line, or a hardcoded-to-stable/edge regression on any of these six cases, turns the exact-match red.
- DONE: Update `docs/releasing.md` (naming sentence + migration bullet) and `docs/site/reference/command-reference.md` (`--version` section) to match the shipped surface.
  `docs/releasing.md` (+14/-1): "What the Tag Push Does" first bullet gains the asset-naming sentence; "Notes" gains the migration bullet verbatim per the entity's Documentation diff. `docs/site/reference/command-reference.md` (+5/-3): intro parenthetical, `--version` prose, both example blocks (outside- and in-session), and the closing "three lines" sentence all updated.

### Summary

Implemented both kept mechanisms exactly as specified: the conditional `_stable` suffix (the same template string already spiked and proven in ideation, re-validated here with `goreleaser check` plus a live `--snapshot` build showing unsuffixed snapshot names) and the `Channel:` line in `--version` reading the effective devBranch through `applyDevBranchOverride(env)`. Proved the Channel line live by building real stable and edge binaries via `goreleaser release --snapshot`, running `--version` on each, and confirming `SPACEDOCK_DEV_BRANCH` overrides (including explicit-empty) render what each binary will actually do. Total diff: +64/-17 across 5 files (.goreleaser.yaml, internal/cli/cli.go, internal/cli/version_session_test.go, docs/releasing.md, docs/site/reference/command-reference.md) — within the ±25 LOC tolerance of the +48/-8 estimate and well under the 80-line hard ceiling. No new test file, CI step, or verifier tool added, per the captain's ruling. `go test ./...` and `go test ./... -race` both green. Code committed to `spacedock-ensign/prerelease-ships-stable-stamped-default-artifact` at bc4f03d57.

## Stage Report: validation

- DONE: AC-1's evidence lane is a real `-pre` release page, and no such cut exists yet — establish what CAN be proven now (the template's three naming contexts, exercised) and state plainly that the published-page count is unproven until a cut happens. Do not rubber-stamp it and do not invent a substitute that cannot fail.
  Real `goreleaser release` (2.16.0, throwaway clone of bc4f03d57) over all three contexts: `-pre` tag `v0.99.0-pre1` → 8 archives, ZERO unsuffixed; stable tag `v0.99.0` → 4 unsuffixed + 4 `_edge`; `--snapshot` → 4 unsuffixed + 4 `_edge`; checksums.txt keys follow in every case. Falsified against the PRE-FIX config on the same tags: 4 unsuffixed on `-pre` (the defect reproduced), while stable-tag and snapshot names *and* checksums keys came out byte-identical pre-fix vs post-fix. The published-page count stays unproven — see AC-1.
- DONE: Attack the residual the ideation body already discloses: the fixture default is devBranch=next, so a Channel line hardcoded to edge would pass every fixture. Prove the stable rendering independently.
  Mutation matrix on a throwaway clone: deleting the Channel line, moving it below Runtime/Sandbox, hardcoding the word to `stable`, and inverting the map each turn `TestVersionSessionRender` RED; hardcoding word+id to `edge` stays GREEN — the disclosed residual is real, confirmed rather than assumed. Stable rendering proven with no fixture involved: the `devBranch=main` binary extracted from the `-pre` run's `_stable` tarball (stamp read independently via `go version -m`) prints `Channel: stable (spacedock@spacedock)`.
- DONE: Confirm the captain's cut held: no new test file, no CI step, no verifier tool was added, and the diff stays within about 64 inserted lines across 5 files.
  `git diff --numstat 04e936a12..bc4f03d57`: 64 insertions / 17 deletions across exactly the 5 declared files; `--diff-filter=A` empty (no file added); no `.github/` or `scripts/` path in the diff; `verify-stamp`/`verify-channel-stamps` match nowhere in the tree. Within the +25 LOC tolerance on the +48/−8 estimate, under the 80-line ceiling.

### Acceptance criteria

- **AC-1 — mechanism proven; the AC's own measure not yet measurable.** The zero-unsuffixed-archive count is proven over the artifacts the shipped `.goreleaser.yaml` really produces on a real `-pre` tag (0 of 8, four os/arch pairs, checksums keys following), and it moves the wrong way (0 → 4) the instant the template regresses to its pre-fix form. NOT proven: the count over a *published* release page — no `-pre` cut exists yet, and the one unexercised step is goreleaser's upload of `dist` to the GitHub release. Audited what else could put a default-named asset on a release: nothing in `release:` (no `extra_files`; source archives off by default), and the only other uploader, `.github/workflows/release.yml:114`, adds `journey-costs-v<ver>.json` — not an archive.
- **AC-2 — PASS.** Real released binaries, not fixtures: the `_stable` and `_edge` tarballs from the `-pre` run print `Channel: stable (spacedock@spacedock)` / `Channel: edge (spacedock@spacedock-edge)` at line 3 in BOTH shapes (outside-session and `CLAUDECODE=1`). Effective-value semantics exercised on both binaries across `SPACEDOCK_DEV_BRANCH` ∈ {`main`, `next`, `""`, `feature-x`, `MAIN`, `" main"`}: only exact `main` renders `stable`, matching `channelMarketplace`'s branch including case and whitespace. Divergence check against the frontdoor, live: with a stub `claude` on PATH in an isolated HOME, `install --host claude` issued `plugin install spacedock@spacedock` from the stable binary and `plugin install spacedock@spacedock-edge` from the edge one — the parenthetical is the id the binary actually ensures.

### Checks run

- `go test ./...` and `go test ./... -race` green in the worktree; `gofmt -l ./cmd ./internal` empty.
- Consumers, run verbatim against a real three-line binary: install-e2e's `case "$ver" in spacedock\ ?*)` matches; `grep -c '^Sandbox: '` = 1; `grep -c '^contract 3$'` = 1; line 1 is still the bare `spacedock <version>` token the FO gate parses.
- Cask pipe against the renamed archive (casks redirected to a local bare repo; builds+archives stanzas asserted byte-identical to the shipped config): both casks generate on both tags; on `-pre` the stable cask self-skips with `prerelease detected with 'auto' upload` while the edge cask — the one that does publish on a prerelease — still pins `_edge`; on a stable tag the stable cask pins the unsuffixed URLs, so the tap contract is unchanged.
- install.sh: from a `--snapshot` dist (install-e2e's real shape) it installs and the binary runs; from a prerelease dist it dies `no spacedock_*_darwin_arm64.tar.gz` rather than installing the wrong thing; its URL path still constructs the unsuffixed name, so it cannot reach a prerelease `_stable` asset.

### Deferred risks

- **A hardcoded-to-edge regression passes every fixture.** Released user/workflow: any operator reading the channel off `--version`. Harm: a stable binary would print `edge`. Authority: `value-ac[AC-2]`. Trigger: an edit dropping the `devBranch == "main"` branch — no automated guard exists. Disclosed in the Test plan and captain-accepted under the no-new-test-surface ruling; confirmed real by mutation here. Promotes to material if a Channel-line change ever ships without a live stable-binary run.
- **The `--version` path's `applyDevBranchOverride(env)` call has no guard at all.** Deleting it leaves the whole `go test ./...` suite green. Harm: `--version` would report the stamp instead of the effective channel under an explicit `SPACEDOCK_DEV_BRANCH`. Authority: `value-ac[AC-2]` ("renders the EFFECTIVE devBranch"). Trigger: refactoring the version branch; the override is a dev-only path and the shipped behavior is correct (exercised live above). Same accepted kind as the residual above, but NOT disclosed in the body — recording it. Promotes to material if the override becomes part of a documented operator procedure.
- **`release.prerelease: auto` and the archive template read the same semver signal today.** If `release.prerelease` were ever pinned `false` on a hyphenated tag, install.sh's `/releases/latest` path would build the unsuffixed name and 404. Fail-loud (`die "download failed"`), never the incident's silent wrong binary. Promotes to material on any change to `release.prerelease`.

### Needs decision

- **AC-1's stated evidence lane cannot close inside this task.** The AC names "the first post-landing `-pre` release's published asset list", and landing this change is what enables that cut. Either the task terminalizes carrying the published-page count as an explicit obligation on the next `-pre` cut, or validation blocks until a cut happens — a sequencing call the captain owns. Recommending the former: every step under this task's control is proven at the strongest fidelity available without publishing, and the one unexercised step (goreleaser uploading `dist`) is untouched by this change.

### Material Findings

- None.

### Recommendation

**PASSED**, with one carried obligation: on the first `-pre` cut after this lands, read the release page, confirm zero unsuffixed archives, and record that count against AC-1. No material finding; three deferred risks recorded above with promote conditions.

### Summary

Proved the naming mechanism by running real goreleaser 2.16.0 releases over all three contexts and falsifying each against the pre-fix config on the same tags — 4 unsuffixed archives before, 0 after on a `-pre` tag, byte-identical names and checksums keys on stable tags and snapshots. Proved AC-2 on real stable- and edge-stamped binaries rather than fixtures, including a live check that the Channel line's parenthetical equals the plugin id the binary actually passes to `claude plugin install`, and confirmed by mutation that the disclosed hardcoded-to-edge residual is genuinely unguarded — plus a second, undisclosed instance of the same kind (the `--version` override call). The captain's cut held exactly (64/17 across the 5 declared files, no new file, no CI step, no verifier tool); AC-1's published-page count is the one thing that stays open, by construction, until a cut happens.
