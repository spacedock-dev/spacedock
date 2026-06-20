---
title: Pre-release tag guards — release.yml skips stable-promotion + goreleaser skips cask bump for v*-pre tags
status: implementation
sprint: 0221-layered-fo
group: release-engineering
id: evkr57zq9g6c2fbk9v8b40kr
---

Enable a single coherent pre-release tag (`v0.23.0-pre`, openai/codex style) that ships binary + content together as a GitHub pre-release WITHOUT promoting to stable. The release machinery currently treats every `v*` tag as a full stable cut: a `-pre` tag would (wrongly) stamp main's manifests, advance the `stable` ref, and bump both Homebrew casks. Add two guards so a hyphenated (pre-release) tag stays out of the stable channel; a normal `vX.Y.Z` tag is unaffected.

## Fix
1. **`.github/workflows/release.yml`** — the post-goreleaser "Stamp plugin manifests" + stable-ref-advance steps (~lines 199-224) run unconditionally. Add `if: "!contains(github.ref, '-')"` so they SKIP for a hyphenated pre-release tag (a final `vX.Y.Z` has no hyphen and still runs).
2. **`.goreleaser.yaml`** — the Homebrew cask block(s) have no `skip_upload`, so a `-pre` tag would bump the stable `spacedock` cask. Add `skip_upload: auto` so cask publishing is skipped for pre-releases (goreleaser's `release.prerelease: auto` already marks a `-pre` tag as a GitHub pre-release; testers install from the release assets or `go install ...@v0.23.0-pre`).

Make the SMALLEST change. Do NOT alter the e2e-gate (a `-pre` tag still gates on a green live E2E or the waiver — that is intended). Do NOT change the final-release path behavior.

## Acceptance criteria
- **AC-1** — a FINAL `vX.Y.Z` tag (no hyphen) still stamps the manifests, advances `stable`, and publishes the casks (no regression to stable releases). Verified by evaluating the `if` condition against a no-hyphen ref (GitHub Actions `contains` semantics — an independent source) AND a goreleaser dry-run/snapshot showing the cask publish path for a non-prerelease version.
- **AC-2** — a PRE-RELEASE `vX.Y.Z-pre` tag SKIPS the manifest-stamp + stable-ref-advance steps and does NOT bump the casks, while goreleaser still marks the GitHub release as a pre-release. Verified by evaluating the `if` against a hyphenated ref (must be false) AND goreleaser `skip_upload: auto` + `prerelease: auto` behavior on a `-pre` version (dry-run/snapshot evidence, not just that the YAML key is present).
- **AC-3** — `goreleaser check` passes (config valid) and `release.yml` lints clean (actionlint or equivalent); reported in the stage report.

This is release-critical machinery — proof must demonstrate the OUTCOME (a `-pre` tag avoids stable promotion; a final tag still promotes), not merely that the guard text is present. Use the strongest achievable proof short of pushing a real tag (goreleaser check/snapshot + condition evaluation); a live tag push is NOT part of this task.
