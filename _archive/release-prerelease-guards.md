---
title: Pre-release tag guards — release.yml skips stable-promotion + goreleaser skips cask bump for v*-pre tags
status: done
sprint: 0221-layered-fo
group: release-engineering
id: evkr57zq9g6c2fbk9v8b40kr
worktree: .worktrees/spacedock-ensign-release-prerelease-guards
started: 2026-06-20T21:13:05Z
pr: local-merge:e8dfd687
verdict: passed
completed: 2026-06-21T03:07:04Z
archived: 2026-06-21T03:07:04Z
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

## Stage Report: implementation

- DONE: release.yml stamp/stable-advance step SKIPS on a hyphenated tag and STILL runs on a final vX.Y.Z (proven by condition evaluation on both ref shapes)
  Added `if: "!contains(github.ref, '-')"` to the single "Stamp plugin manifests" step (its run block holds BOTH the manifest stamp and `git push origin main:refs/heads/stable`, so one guard covers both). Confirmed no other step advances stable/stamps. Condition evaluated on `refs/tags/v0.23.0` → RUNS and `refs/tags/v0.23.0-pre` → SKIPS via GitHub Actions `contains` (case-insensitive substring; docs.github.com/actions expressions#contains). Worktree commit cc009c4d.
- DONE: .goreleaser.yaml casks get `skip_upload: auto` so a -pre tag does NOT bump the casks while goreleaser still marks the release a pre-release
  Added `skip_upload: auto` to both casks (`spacedock`, `spacedock@next`). Proven end-to-end against goreleaser's OWN publish pipe (mock client): pre→publish SKIPPED + no cask file; final→PUBLISHED + cask file created. Source is cask.go:147-148 (`skip_upload=="auto" && ctx.Semver.Prerelease != ""` ⇒ skip).
- DONE: `goreleaser check` passes and release.yml lints clean (actionlint)
  `goreleaser check` v2.16.0 → "1 configuration file(s) validated". `actionlint` (via `go run github.com/rhysd/actionlint/cmd/actionlint@latest`) on release.yml → exit 0, no findings (actionlint not on PATH; go-run fallback documented).
- DONE: [BEYOND ENUMERATED GUARDS] added `release.prerelease: auto` to close an AC-2 gap
  The assignment/entity assumed goreleaser's default already marks a -pre tag a pre-release. FALSE for v2.16: default is `false`/`''` (release.go:74-82 has only `auto`/`true` cases; upstream test TestDefaultPreRelease/"release" proves Release{}+Prerelease="rc1" ⇒ PreRelease=FALSE). Without the key a -pre tag would publish as a FULL/latest GitHub release — violating AC-2 + "WITHOUT promoting to stable". Added `prerelease: auto` (one line). Flagged to team-lead before finalizing.

### Final-vs-pre outcome table

| Behavior | final `v0.23.0` | pre `v0.23.0-pre` |
|---|---|---|
| release.yml stamp manifests + advance `refs/heads/stable` (`!contains(ref,'-')`) | RUNS (no `-`) | SKIPS (has `-`) |
| Homebrew casks bump (`skip_upload: auto`) | PUBLISHED (Semver.Prerelease="") | SKIPPED (Semver.Prerelease="pre") |
| GitHub release flag (`release.prerelease: auto`) | normal/latest (PreRelease=false) | pre-release (PreRelease=true) |
| e2e-gate (unchanged) | gates on green live E2E / waiver | gates on green live E2E / waiver |

### Proof commands (re-runnable)

- `goreleaser check` (in worktree) → config valid.
- `go run github.com/rhysd/actionlint/cmd/actionlint@latest .github/workflows/release.yml` → exit 0.
- goreleaser v2.16.0 source (cloned to verify): `go test ./internal/pipe/cask/ -run 'TestRunPipeNoUpload|TestSkipUploadAutoPrereleaseVsFinal'` and `go test ./internal/pipe/release/ -run 'TestPrereleaseAutoMarksHyphenTag|TestNoPrereleaseKeyDefaultsToFinal'` → all PASS (pre→skip/pre-release, final→publish/latest). These tests live in the throwaway goreleaser clone, not this repo.
- `go test ./internal/release/...` (this repo) → ok; the channel-agreement guard still matches the stamp step (only an `if:` + comment were added).

### Summary

Three minimal guards keep a `v0.23.0-pre` tag out of the stable channel while still shipping it as a GitHub pre-release: a `!contains(github.ref,'-')` gate on the release.yml stamp/stable-advance step, `skip_upload: auto` on both casks, and `release.prerelease: auto` so goreleaser actually flags the pre-release. The third went beyond the two enumerated guards because the assignment's stated goreleaser default was wrong (v2.16 defaults `prerelease` to `false`, not `auto`); without it AC-2 ("still marks the release a pre-release") would fail. Proof exercises goreleaser's actual publish/release pipes (mock client) plus condition evaluation on both ref shapes — outcome, not key-presence. No tags were pushed.
