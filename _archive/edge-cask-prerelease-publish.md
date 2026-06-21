---
title: Edge cask publishes on pre-releases — spacedock@next skip_upload auto→false
status: done
sprint: 0221-layered-fo
group: release-engineering
id: g0n7592mf99saddc5g3vp91h
worktree: .worktrees/spacedock-ensign-edge-cask-prerelease-publish
started: 2026-06-21T03:12:46Z
pr: local-merge:68c127fd
verdict: passed
completed: 2026-06-21T03:21:31Z
archived: 2026-06-21T03:21:31Z
---

So the `spacedock@next` Homebrew cask (the edge tap) tracks pre-release tags. Today both casks carry `skip_upload: auto`, so a `-pre` tag skips BOTH — the edge tap stays stale. Let the EDGE cask publish on pre-releases while keeping the STABLE cask clean.

## Fix
`.goreleaser.yaml` — the EDGE cask (`name: spacedock@next`, `ids: [spacedock-edge]`, ~line 57): `skip_upload: auto` → `skip_upload: false` (publishes on every tag, including pre-releases). KEEP the STABLE cask (`name: spacedock`, `ids: [spacedock-stable]`, ~line 18) at `skip_upload: auto` (stable never gets a pre-release). Smallest change — one line on the edge cask only. Do NOT touch the release.yml guards or the `prerelease: auto` key.

## Acceptance criteria
- **AC-1** — on a PRE-RELEASE tag (`vX.Y.Z-pre.N`), the `spacedock@next` (edge) cask PUBLISHES while the stable `spacedock` cask SKIPS. Verified by goreleaser dry-run/source semantics on a `-pre` version (edge `skip_upload:false` → publish; stable `auto` + prerelease → skip), not key-presence.
- **AC-2** — on a FINAL tag (`vX.Y.Z`), BOTH casks still publish (no regression to the stable two-channel cut).
- **AC-3** — `goreleaser check` passes.

Release-critical: prove the OUTCOME (edge publishes / stable skips on a -pre; both on a final), not just the key value. No tag is pushed by this task — the cut is the FO's separate step.

## Stage Report: implementation

- DONE: EDGE cask `skip_upload: auto` → `false`; STABLE cask STAYS `auto`
  `.goreleaser.yaml` edge cask (`name: spacedock@next`, `ids: [spacedock-edge]`, line 153/158) now `skip_upload: false`; stable cask (`name: spacedock`, `ids: [spacedock-stable]`, line 112/117) untouched at `auto`. Edge comment rewritten to describe the divergence. Worktree commit 68c127fd; only `.goreleaser.yaml` changed.
- DONE: AC-1 — on a `-pre` version, edge PUBLISHES and stable SKIPS (outcome, not key-presence)
  Exercised goreleaser v2.16.0's real cask publish pipe with a mock client on semver prerelease `pre.2` (shape of `v0.23.0-pre.2`): edge `false` → `CreatedFile==true`; stable `auto` → `CreatedFile==false`. Source: cask.go:143-148 (`"true"`→skip; `"auto" && Semver.Prerelease!=""`→skip; else publish, so `"false"` never skips).
- DONE: AC-2 — on a FINAL version, BOTH casks PUBLISH (no stable-cut regression)
  Same pipe, prerelease `""` (shape of `v0.23.0`): edge `false` → publishes; stable `auto` → publishes (`Semver.Prerelease==""` so the auto guard does not fire). Both `CreatedFile==true`.
- DONE: AC-3 — `goreleaser check` passes
  `goreleaser check` v2.16.0 in the worktree → "1 configuration file(s) validated", exit 0.

### Edge-vs-stable outcome table

| skip_upload | `v0.23.0-pre.2` (Prerelease="pre.2") | final `v0.23.0` (Prerelease="") |
|---|---|---|
| EDGE `spacedock@next` = `false` | PUBLISHES | PUBLISHES |
| STABLE `spacedock` = `auto` | SKIPS | PUBLISHES |

### Proof commands (re-runnable)

- `goreleaser check` (in worktree) → exit 0, config valid.
- goreleaser v2.16.0 cloned to /tmp/goreleaser-src; added `internal/pipe/cask/spacedock_edge_skipupload_test.go` and ran `go test ./internal/pipe/cask/ -run TestSpacedockEdgeSkipUpload -v` → 4/4 PASS (edge-pre publish, stable-pre skip, edge-final publish, stable-final publish). Test lives in the throwaway clone, not this repo — it exercises goreleaser's own pipe, which this repo does not vendor.
- `release.yml` `!contains(github.ref,'-')` guard and `release.prerelease: auto` key NOT touched (diff is `.goreleaser.yaml`-only, 5 insertions / 3 deletions).

### Summary

One-line behavioral change: the edge `spacedock@next` cask's `skip_upload` flips from `auto` to `false`, so it publishes on every tag including pre-releases while the stable `spacedock` cask keeps `auto` and never publishes a pre-release. Proven by exercising goreleaser v2.16.0's actual cask publish pipe against a mock client (edge publishes on a `-pre`; stable skips; both publish on a final) plus a passing `goreleaser check` — outcome, not key-presence. The accompanying comment was rewritten to describe the new divergence; no other config was touched, and no tag was pushed (the FO cuts that separately).
