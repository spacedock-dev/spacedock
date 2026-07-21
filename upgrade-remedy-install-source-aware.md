---
id: vzsastkvv2r6dpjakw1vq6wx
title: Binary-upgrade prompt must be install-source-aware (brew formula, non-brew, sandbox)
status: ideation
source: "Captain report (CL) 2026-07-16 — the skill-upgrade version gate prompts `brew upgrade spacedock` regardless of how the binary was installed."
started: 2026-07-21T15:58:31Z
completed:
verdict:
score:
worktree:
issue:
---

When the skill is upgraded and determines the binary needs upgrading, the remedy it emits hardcodes `brew upgrade spacedock` — even when that is not how the binary was installed. It must detect the real install source and runtime context and emit the correct upgrade instruction. Likely a follow-up gap in the shipped upgrade-hint work (`install-refresh-and-upgrade-hint`, `init-upgrade-and-contract-remedy`).

## Problem

{Ideation fills this in. Seed: the wrong-version abort / `spacedock doctor` remedy assumes a stable Homebrew install and prints `brew upgrade spacedock`. That is wrong for at least three cases and leaves the user stuck on the old binary: (1) installed from the `spacedock@next` formula/tap — upgrading the stable `spacedock` formula does nothing and misleads; the remedy should target the `@next` formula; (2) non-brew install (source `go build`, a downloaded/notarized binary, `SPACEDOCK_BIN` pointing at a checkout build) — `brew upgrade` does not apply at all; (3) within-sandbox execution — brew may be unavailable or inappropriate in the sandbox context, so a brew remedy is unactionable there.}

## Proposed approach

{Ideation fills this in. Seed: make the remedy generator install-source-aware. Detect which Homebrew formula/tap owns the running binary (`spacedock` vs `spacedock@next`) vs a non-brew path (source build / `SPACEDOCK_BIN` / downloaded), plus whether execution is sandboxed, and emit the matching upgrade instruction: `brew upgrade <owning-formula>` for a brew install, the source-build/download-refresh instruction for non-brew, and a sandbox-appropriate path otherwise. Most likely lives in `spacedock doctor` / the version-gate remedy path (and any skill abort text that restates it).}

## Out of scope

{Ideation fills this in. Seed: no change to the version-COMPARISON logic (what counts as out-of-date) — only the remedy MESSAGE; no auto-upgrade; no new install method.}

## Acceptance criteria

{Ideation fills this in.}

## Test plan

{Ideation fills this in. Seed: table tests over the remedy generator across install sources — brew stable, brew `@next`, non-brew/source, and sandbox — asserting the emitted upgrade instruction matches the source; specifically assert NO `brew upgrade spacedock` (stable) is emitted for an `@next`, non-brew, or sandbox install. gofmt, go test ./..., go test ./... -race.}
