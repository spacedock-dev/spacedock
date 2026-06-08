---
id: v3yf6j3nkey6tvxvsa1y42sm
title: Ship Linux binaries + a Linux install path (release is darwin-only)
status: backlog
source: "captain (2026-06-08) - dispatch before 0.20.0. Mirrors GitHub spacedock-dev/spacedock#321 (public report). The release is darwin-only though the code already builds + tests on Linux in CI; the gap is release config + a Linux distribution channel."
started:
completed:
verdict:
score:
worktree:
issue: "spacedock-dev/spacedock#321"
---

Ship Linux binaries and give Linux users an install path, before the 0.20.0 cut. The code is already portable — the gap is release packaging + distribution.

## Problem

`.goreleaser.yaml` builds `goos: [darwin]` / `goarch: [arm64, amd64]` only, and `release.yml`'s goreleaser job runs on `macos-latest`. A `v*` release publishes ONLY macOS tarballs + a macOS Homebrew cask. There is NO Linux install path: the cask is macOS-only, and `go install …@latest` is a dev/toolchain fallback (unstamped version). So Linux users can't install spacedock.

The code already works on Linux: `runtime-live-e2e.yml`'s offline gate (`runs-on: ubuntu-latest`) runs `go build ./...` + `go test ./...`, and other jobs `go build -o ./spacedock` on ubuntu. So building a Linux binary is a trivial cross-compile (CI proves it). The work is release config + distribution, not portability.

## Proposed approach

1. **Build (one-liner):** add `linux` to goreleaser `builds.goos` (amd64 + arm64) and the matching `archives` (the existing `name_template` `{ProjectName}_{Version}_{Os}_{Arch}` already yields `spacedock_{version}_linux_{arch}.tar.gz`). Pure-Go cross-compile from the existing runner.
2. **Distribution — a `curl | sh` installer (recommended):** add `install.sh` that detects OS/arch, downloads the matching tarball from the latest GitHub Release (verifying against `checksums.txt`), and installs `spacedock` to `~/.local/bin` (or `/usr/local/bin`). Universal — works on Linux AND macOS; brew stays the mac-preferred path. No Linuxbrew requirement, no new package infra. Document the Linux/script path in `docs/install-journey.md`.
   - Alternatives considered (ideation can revisit): Homebrew formula on Linuxbrew (loses the cask's Gatekeeper postflight, uncommon); deb/rpm/AUR packages (heavier maintenance); `go install` (dev-only, unstamped).
3. **Linux runtime caveats:** the host CLIs (Claude Code / Codex / Pi) run on Linux. The macOS release steps (adhoc signing, `com.apple.quarantine` / Gatekeeper, notarization) do NOT apply to Linux. **safehouse** on Linux is a separate question — its sandbox is macOS-oriented; confirm Linux support or state the Linux sandbox story honestly rather than implying it works (do not overclaim, per the README proof-policy).

## Out of scope

- The two-channel brew tap / per-channel `devBranch` stamp / next-publish pipeline (separate flip release-mechanics).
- safehouse's own Linux sandbox implementation (a safehouse-side concern) — this task only states the Linux story honestly.

## Acceptance criteria

Ideation/implementation fills in. Sketch:

- A `v*` release publishes `linux/amd64` + `linux/arm64` tarballs alongside the darwin ones (verified by the release artifact list / a goreleaser `--snapshot` producing them).
- The `install.sh` installs a runnable `spacedock` on Linux AND macOS from a Release tarball (verified by a fixture/CI run that fetches + installs + `spacedock --version` works), checksum-verified.
- `docs/install-journey.md` documents the Linux install path.
- The safehouse-on-Linux story is stated accurately (not implied to work if it doesn't).

## Test plan

Ideation/implementation fills in. `goreleaser release --snapshot --clean` produces the linux tarballs (offline check); a CI/fixture run of `install.sh` against a Release tarball asserts `spacedock --version`; a content check that install-journey carries the Linux path.
