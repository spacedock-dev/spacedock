---
id: c9qz40bza726q58m3s8hsdq5
title: Prerelease releases ship a stable-stamped binary under the default asset name
status: ideation
source: "Live cross-host report, CL 2026-08-17. A Linux box (linux/amd64) running `spacedock 0.27.0-pre7` reinstalled the Claude plugin as `spacedock@spacedock` on every launch, overwriting a hand-installed edge plugin. `go version -m $(command -v spacedock)` on that binary returned `devBranch=main`; the darwin/arm64 edge-cask binary at the identical version returned `devBranch=next`. Same tag, same `--version` output, opposite channel."
started: 2026-08-17T19:38:45Z
completed:
verdict:
score:
worktree:
issue:
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

**Verification moves to the artifact** (AC-2), three layers:

1. **Pre-publish, in-pipeline:** each `builds[].hooks.post` runs `go run ./cmd/spacedock-release verify-stamp {{ .Path }} <main|next>`. The checker reads the binary with `debug/buildinfo` (the `go version -m` surface) and exits non-zero unless the recorded `-ldflags` carry `cli.devBranch=<expected>`. A failing hook aborts `goreleaser release` at the build stage, before anything uploads (spiked — see Spike record 5).
2. **Whole-set sweep:** `spacedock-release verify-channel-stamps <dist>` maps every `*.tar.gz` name to its advertised channel (`_edge` → next, else → main), extracts the binary, and asserts the buildinfo stamp matches; for the native-OS binary it additionally EXECUTES `spacedock --version` and asserts the new `Channel:` line agrees — execution catches the silent-no-op class (an `-X` flag whose target symbol was renamed still appears in recorded ldflags but changes nothing; the config comments already document that hazard). It also enforces AC-1 on artifacts: a real-prerelease dist containing any unsuffixed archive fails. Fail-closed on zero archives. Wired into install-e2e.yml after its existing snapshot step (every PR, pre-merge) and into release.yml right after goreleaser (the real cut; post-publish but loud within minutes — OSS goreleaser has no pre-publish seam for whole-set checks; layer 1 covers pre-publish).
3. **Config guard:** extend `internal/release/channel_agreement_guard_test.go` to EXECUTE the stable archive's name-suffix template under three contexts (real prerelease → `_stable`; snapshot → unsuffixed; stable → unsuffixed), and pin the edge archive's unconditional `_edge` suffix plus the archive-ids ↔ name-template binding. Executing the template, not string-matching it, keeps the test on behavior.

**Legibility (AC-3):** `--version` gains one line in BOTH output shapes, directly after `OS:`:

```
Channel: stable (spacedock@spacedock)
Channel: edge (spacedock@spacedock-edge)
```

The channel word maps from the effective devBranch exactly as `channelMarketplace` does (main → stable, anything else → edge); the parenthetical is `channelPluginID(devBranch)` — the id the frontdoor re-ensures, so the incident's symptom (the `spacedock@spacedock` reinstall loop) and its cause read off one line. It reports the EFFECTIVE value: `printVersion` already receives `getenv`, so a `SPACEDOCK_DEV_BRANCH` override renders what the binary will actually do. Line 1 stays the bare version token; consumers audited: the FO version gate parses line 1 only, fo-install-gate greps `^Sandbox: ` (prefix-anchored), and the frozen `contract 3` token is position-independent below line 1.

Mechanism → value audit (simplest alternative for each):

- Conditional `_stable` suffix (serves AC-1). Alternative: skip the stable build — blocked by the cask pipe (spiked). Alternative: docs-only warning — leaves the trap in place.
- `verify-stamp` build hook (serves AC-2). Alternative: sweep only — on the real cut the sweep runs post-publish; the hook is the only OSS pre-publish enforcement point.
- `verify-channel-stamps` sweep (serves AC-2 and AC-1-on-artifacts). Alternative: hooks only — a hook sees one binary and never the archive name, so a crossed `archives.ids` mapping (the stable archive packaging the edge build) is invisible to it.
- `Channel:` line (serves AC-3; doubles as layer 2's execution oracle). Alternative: keep `go version -m` as the channel oracle — precisely the out-of-band-only state that made the incident an investigation.

## Spike record

Spiked 2026-08-17 with goreleaser 2.16.0: full `goreleaser release` runs in a scratch clone against throwaway local tags `v0.99.0-pre1` / `v0.99.0`, cask pushes aimed at a local bare repo via `repository.git.url`, `release.disable: true` plus an explicit cask `url.template` standing in for the GitHub API. Facts established:

1. `builds[].skip` accepts a template (schema: string|bool) and fires on a `-pre` tag (`skip is set id=spacedock-stable`), and the archive/checksum pipes tolerate the resulting empty stable set — but the stable cask pipe errors `no linux/macos archives found matching goos=[darwin linux] goarch=[amd64 arm64] ids=[spacedock-stable]` before `skip_upload: auto` is consulted, and casks have no `disable` field. The build-skip direction is dead on OSS goreleaser.
2. The conditional name suffix works end to end: the `-pre` tag produced only `_stable` + `_edge` archives (checksums.txt matching), both casks generated, and the stable cask publish skipped itself with `prerelease detected with 'auto' upload`. The stable tag produced today's unsuffixed + `_edge` names with cask URLs referencing the unsuffixed asset — the tap contract is byte-stable.
3. The `not .IsSnapshot` guard is load-bearing: without it a snapshot names the stable asset `..._stable.tar.gz` (snapshot versions carry a semver prerelease part) and install-e2e's end-anchored install.sh glob finds nothing.
4. `go version -m` / `debug/buildinfo` reads `-X ...cli.devBranch=...` from cross-compiled tarball binaries on a foreign OS (linux ELF read on darwin) — the AC-2 read mechanism, identical to the incident's own forensic read.
5. `builds[].hooks.post` runs once per built artifact with `{{ .Path }}` resolved; a hook asserting the wrong stamp fails the whole release at the build stage (`post hook failed: exit status 1`) — the pre-publish abort direction proven, not assumed.

## Migration for already-installed prerelease binaries

The naming change is forward-only: already-published `-pre` releases keep their unsuffixed stable assets, and an installed binary cannot change its stamp. For an operator holding one (the incident machine: linux/amd64, 0.27.0-pre7, devBranch=main):

- **Detect:** on existing binaries, `go version -m $(command -v spacedock)` showing `cli.devBranch=main` means stable; from this change on, `spacedock --version` prints the channel in-band. Symptom key: a prerelease-versioned binary that re-ensures `spacedock@spacedock` on every launch and overwrites an edge plugin is a stable-stamped binary, not a frontdoor bug.
- **Remedy:** replace the binary with the same release's `_edge` tarball or `brew install spacedock-dev/tap/spacedock@next`. No state migration is needed: the frontdoor re-ensures the channel plugin per launch, so plugin/marketplace state converges on the next launch once the binary's channel is right.
- **Optional retroactive closure (captain decision at the gate — destructive, published surface):** delete the unsuffixed stable assets from existing `-pre` releases (0.27.0-pre1..pre7 and earlier `-pre` cuts). Nothing sanctioned references them — the stable cask never published from a prerelease and install.sh cannot resolve one — and their checksums.txt lines dangle harmlessly (lookups are by filename). Until this is done, old prerelease pages still carry the trap.
- Stable-release installs are untouched: names, install.sh, and tap behavior stay byte-identical for `vX.Y.Z` tags.

## Out of scope

The duplicate edge routes in the marketplace (`spacedock-edge@spacedock` and `spacedock@spacedock-edge` resolve the same payload) — a published-surface decision, filed separately if CL wants it changed.

Whether `channelEntry`'s constant entry name (internal/cli/host_exec.go:223) makes an *edge* binary reinstall over a route-A install: plausible from reading, unverified, and not the cause of this report. Do not fold it in without first confirming the frontdoor's actual presence-check behavior.

## Expected surface and tolerance

Estimate net LOC change: +560 across 11 files (supersedes the +140/5-file seed guess, which did not count the verifier's unit tests, the version-fixture literal updates, or the two workflow wirings):

- `.goreleaser.yaml` +20 (name_template conditional, two post hooks, comments)
- `cmd/spacedock-release/verify_stamp.go` +170 (new: both subcommands)
- `cmd/spacedock-release/verify_stamp_test.go` +180 (new)
- `cmd/spacedock-release/main.go` +10 (routing)
- `internal/cli/cli.go` +15 (Channel line in printVersion)
- `internal/cli/version_session_test.go` +40 net (literal updates plus channel/override cases)
- `internal/release/channel_agreement_guard_test.go` +90 (template-execution + ids-binding guards)
- `.github/workflows/install-e2e.yml` +12 (sweep step)
- `.github/workflows/release.yml` +15 (sweep step)
- `docs/releasing.md` +25/−5
- `docs/site/reference/command-reference.md` +10/−4

Tolerance: ±2 files, +150/−200 LOC. Beyond that, back to the gate.

## Declared semantic changes

- **Command grammar:** `--version` gains a third line `Channel: <stable|edge> (<plugin-id>)` in BOTH shapes; the outside-session shape goes two lines → three. Known consumers audited (FO version gate: line 1 only; fo-install-gate: `^Sandbox: ` prefix grep; `contract 3`: position-independent). New `spacedock-release verify-stamp` / `verify-channel-stamps` subcommands (release tooling, not the user CLI).
- **Published-release surface:** on `-pre` tags only, the stable tarball's asset name gains `_stable`; checksums.txt keys follow. Stable-tag and snapshot asset names stay byte-identical to today.
- **Stored formats:** none. **Authority:** none. **Runtime install behavior:** none — frontdoor, marketplace mapping, install.sh, and tap flow untouched.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A prerelease release offers no default-named asset: every asset a `-pre` release publishes names its channel, so reaching for "the" asset is impossible and the edge asset is the only undecorated-looking choice a `-pre` page offers.**
This is the measuring AC: the unsuffixed-archive count over a real-prerelease artifact set must be ZERO — a count over the actually-built dist that moves the wrong way the moment the template regresses. Verified by: (config) a guard test that EXECUTES the stable archive's name-suffix template under three contexts — real prerelease → `_stable`, snapshot → unsuffixed, stable → unsuffixed — and pins the edge archive's unconditional `_edge` suffix plus the archive-ids ↔ name-template binding; (artifact) `verify-channel-stamps` fails any real-prerelease dist containing an unsuffixed archive — enforced on every PR (install-e2e) and on the real cut (release.yml); (live) the first `-pre` cut after landing lists only `_stable` + `_edge` assets on its release page. Fails if `.goreleaser.yaml` gives `spacedock-stable` the bare `{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}` name while prerelease publishing is enabled.

**AC-2 - The channel guard proves the stamp on a built ARTIFACT, not only in config text.**
Verified by: `spacedock-release verify-channel-stamps <dist>` — run over the real goreleaser snapshot in install-e2e (every PR) and over `dist/` in release.yml (every cut) — extracts each archive's binary and reads its embedded ldflags via `debug/buildinfo` (the `go version -m` surface), asserting `cli.devBranch` equals the channel the asset name advertises; the native-OS binary is additionally EXECUTED and its `--version` `Channel:` line must agree, which catches the silent-no-op `-X` class that recorded ldflags cannot. `verify-stamp` in `builds[].hooks.post` enforces the same read pre-publish inside the goreleaser pipeline. Unit tests drive both verifiers against really-built binaries stamped each way; the mismatch direction must exit non-zero. Fails if a stable-stamped binary is produced under an edge asset name or vice versa — the exact condition config-parsing cannot detect, and the reason the existing guard passed through this incident.

**AC-3 - A running binary reports the channel it drives, so a recurrence is detectable in-band.**
Verified by: fixture tests over `--version` output asserting the channel line tracks the effective devBranch — devBranch=main renders `Channel: stable (spacedock@spacedock)`, devBranch=next renders `Channel: edge (spacedock@spacedock-edge)`, and a `SPACEDOCK_DEV_BRANCH` override is reflected. Fails if the line is dropped, hardcoded to one channel, or diverges from the plugin id the frontdoor actually ensures. This AC serves AC-1's end value: the release being correct is the goal, and this is what lets an operator confirm it rather than trust it — it is also the execution oracle AC-2's sweep uses.

## Test plan

- **Unit (`go test ./...`, no goreleaser dependency):** `verify_stamp_test.go` builds `./cmd/spacedock` via `go build -ldflags "-X ...devBranch=<main|next>"` into `t.TempDir` (build cache keeps repeat cost low), then drives `verify-stamp` in both directions and `verify-channel-stamps` over synthesized dist dirs: a matched dir passes; a crossed name/stamp dir fails; an empty dir fails; a real-prerelease dir containing an unsuffixed asset fails. Each test's red condition is the named regression, not a tautology over its own fixtures.
- **Config guard (unit):** template execution over the three naming contexts; the edge suffix and ids binding pinned — the AC-1 config half.
- **Fixture/CLI:** `version_session_test.go` literal updates plus the two channel cases and the override case — the AC-3 half.
- **CI:** install-e2e's existing snapshot leg gains the sweep step, proving name↔stamp on real cross-compiled artifacts for both OSes on every PR; its existing install.sh legs keep proving the unsuffixed snapshot contract (the `not .IsSnapshot` guard's consumer).
- **Live:** the first `-pre` cut after landing is the end-value measurement — release page lists only `_stable`+`_edge`, release.yml's sweep step green over the published dist, and `spacedock --version` from the downloaded `_edge` tarball prints `Channel: edge (spacedock@spacedock-edge)`. Rides an already-planned prerelease cut; no dedicated live workflow.
- Estimated cost: moderate — the expensive goreleaser runs ride existing CI legs; new code is one verifier file plus test updates.

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
  tolerate a skipped build — and every binary's embedded stamp is verified
  against its asset name before publish (`spacedock-release verify-stamp` in
  the build hooks) and after (`verify-channel-stamps` over `dist/`);
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
