---
title: "Release FreeBSD binaries"
source: "Captain request during the 0.28.0 release, 2026-11-10; Subspace added FreeBSD release targets on 2026-11-05"
id: h0xv2p2a7kts1k124hjdgb2v
status: backlog
gates:
    version: 1
    records:
        - id: gate:h0xv2p2a7kts1k124hjdgb2v:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:h0xv2p2a7kts1k124hjdgb2v-backlog-1
              briefing:
                id: briefing:h0xv2p2a7kts1k124hjdgb2v:backlog:attempt-1:revision-1
                digest: sha256:f7e4c48a9345b32bd99144a6d3cb586efcc0a304692709cbbe31aa4732c73505
                room-ref: '@review/backlog/briefing-1'
---

Spacedock releases ship darwin and linux binaries only, and `install.sh` rejects any other OS. A FreeBSD user has to build from source. Subspace, which the Spacedock installer also installs, now builds FreeBSD amd64 and arm64 binaries; Spacedock should match so a FreeBSD host can install the whole kit.

## What is known

- `CGO_ENABLED=0 GOOS=freebsd go build ./cmd/spacedock` and `GOOS=freebsd go vet ./...` pass for amd64 and arm64 on the 0.28.0 tree. The release build is already `CGO_ENABLED=0`, and OS-specific code is split on the `unix` build tag, which includes FreeBSD.
- `.goreleaser.yaml` lists `darwin` and `linux` in the `goos` anchor shared by the stable and edge builds.
- `install.sh` `detect_os` maps only `Darwin` and `Linux`.
- Subspace cross-builds FreeBSD in CI as a compile check (output discarded), adds `freebsd/amd64` and `freebsd/arm64` to its release targets, and maps `FreeBSD` in its installer. Its latest full release, v0.16.4, has no FreeBSD assets; v0.17.0-beta.1 is the first candidate to carry them.

## Outcome

A tagged Spacedock release publishes `freebsd_amd64` and `freebsd_arm64` tarballs for both channels with checksums, `install.sh` on FreeBSD installs the matching one, and CI fails a change that breaks the FreeBSD build.

## Unknowns for ideation

1. Whether the goreleaser Homebrew cask step tolerates FreeBSD archives in the same release, or needs them excluded from the cask.
2. What the combined installer does on FreeBSD when Subspace or Compass has no FreeBSD asset: skip that extension with a named message, or fail.
3. Whether `install.sh` uses anything the FreeBSD base system lacks (checksum tool, `install` flags, `curl` versus `fetch`).
4. Which hosts run on FreeBSD (Claude Code, Codex, Pi), and what the docs should promise; a binary without a supported host is of little use.
5. Whether the Safehouse sandbox front door, a macOS mechanism, fails clearly on FreeBSD.
6. How to prove the installer on FreeBSD: a FreeBSD VM job in CI, or a compile check plus installer unit tests with a stubbed `uname`.

## Out of scope

Windows and other BSDs. Notarization and packaging for FreeBSD ports or pkg.
