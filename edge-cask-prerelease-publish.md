---
title: Edge cask publishes on pre-releases — spacedock@next skip_upload auto→false
status: implementation
sprint: 0221-layered-fo
group: release-engineering
id: g0n7592mf99saddc5g3vp91h
worktree: .worktrees/spacedock-ensign-edge-cask-prerelease-publish
started: 2026-06-21T03:12:46Z
---

So the `spacedock@next` Homebrew cask (the edge tap) tracks pre-release tags. Today both casks carry `skip_upload: auto`, so a `-pre` tag skips BOTH — the edge tap stays stale. Let the EDGE cask publish on pre-releases while keeping the STABLE cask clean.

## Fix
`.goreleaser.yaml` — the EDGE cask (`name: spacedock@next`, `ids: [spacedock-edge]`, ~line 57): `skip_upload: auto` → `skip_upload: false` (publishes on every tag, including pre-releases). KEEP the STABLE cask (`name: spacedock`, `ids: [spacedock-stable]`, ~line 18) at `skip_upload: auto` (stable never gets a pre-release). Smallest change — one line on the edge cask only. Do NOT touch the release.yml guards or the `prerelease: auto` key.

## Acceptance criteria
- **AC-1** — on a PRE-RELEASE tag (`vX.Y.Z-pre.N`), the `spacedock@next` (edge) cask PUBLISHES while the stable `spacedock` cask SKIPS. Verified by goreleaser dry-run/source semantics on a `-pre` version (edge `skip_upload:false` → publish; stable `auto` + prerelease → skip), not key-presence.
- **AC-2** — on a FINAL tag (`vX.Y.Z`), BOTH casks still publish (no regression to the stable two-channel cut).
- **AC-3** — `goreleaser check` passes.

Release-critical: prove the OUTCOME (edge publishes / stable skips on a -pre; both on a final), not just the key value. No tag is pushed by this task — the cut is the FO's separate step.
