---
id: c9qz40bza726q58m3s8hsdq5
title: Prerelease releases ship a stable-stamped binary under the default asset name
status: backlog
source: "Live cross-host report, CL 2026-08-17. A Linux box (linux/amd64) running `spacedock 0.27.0-pre7` reinstalled the Claude plugin as `spacedock@spacedock` on every launch, overwriting a hand-installed edge plugin. `go version -m $(command -v spacedock)` on that binary returned `devBranch=main`; the darwin/arm64 edge-cask binary at the identical version returned `devBranch=next`. Same tag, same `--version` output, opposite channel."
started:
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

{Ideation fills this in. Candidate directions for the asset trap: suppress the stable build on prerelease tags; give both channels an explicit suffix so neither owns the default name; or make install.sh and the tap channel-explicit and fail loudly on mismatch. Legibility is part of the same value, not a separate task: whatever ships should also make the channel visible in the binary's own output so a future mismatch is a one-line check rather than an investigation.}

## Out of scope

The duplicate edge routes in the marketplace (`spacedock-edge@spacedock` and `spacedock@spacedock-edge` resolve the same payload) — a published-surface decision, filed separately if CL wants it changed.

Whether `channelEntry`'s constant entry name (internal/cli/host_exec.go:223) makes an *edge* binary reinstall over a route-A install: plausible from reading, unverified, and not the cause of this report. Do not fold it in without first confirming the frontdoor's actual presence-check behavior.

## Expected surface and tolerance

Estimate net LOC change: +140, across 5 files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - Installing the default-named asset from a prerelease release yields a binary whose channel matches that release's channel.**
Verified by: a release-config test asserting that for a `-pre` version the stable build is either not published or not assigned the unsuffixed `name_template`. Fails if `.goreleaser.yaml` gives `spacedock-stable` the bare `{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}` name while prerelease publishing is enabled.

**AC-2 - The channel guard proves the stamp on a built ARTIFACT, not only in config text.**
Verified by: a test that builds or downloads an archive and reads its embedded ldflags via `go version -m`, asserting `devBranch` equals the channel the asset name advertises. Fails if a stable-stamped binary is produced under an edge asset name or vice versa — the exact condition config-parsing cannot detect, and the reason the existing guard passed through this incident.

**AC-3 - A running binary reports the channel it drives, so a recurrence is detectable in-band.**
Verified by: a fixture test over `--version` output asserting the channel value tracks the build's `devBranch` — a devBranch=main build renders stable, a devBranch=next build renders edge. Fails if the field is dropped or hardcoded to one channel. This AC serves AC-1's end value: the release being correct is the goal, and this is what lets an operator confirm it rather than trust it.
