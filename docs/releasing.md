# Releasing Spacedock

This document describes the post-main-flip release lane. Starting with
`v0.20.0`, stable releases are cut from `main`; `next` remains a dev-only branch
for source builds and pre-stable publishing tests. Do not present `next` as the
stable marketplace source after the flip.

## 0.20.0 Main Flip

The `0.20.0` release is the branch transition. It must not ship until the release
gate fixes, runtime-live checks, README/install reconciliation, and upgrade-path
confirmation have passed.

1. Record the current pre-v1 `origin/main` SHA, then archive it:

   ```bash
   git fetch origin main next
   preflip_main="$(git rev-parse origin/main)"
   git tag -a v0-archived "$preflip_main" -m "archive pre-v1 main before 0.20.0 flip"
   git push origin v0-archived
   ```

2. Prepare the release line from `origin/next`. The release-prep commit must
   stamp `0.20.0`, move stable marketplace references to `main`, and remove the
   released-binary install pin to `next`. Keep any dev-only `next` publish path
   separate.

   ```bash
   git worktree add .worktrees/release-0.20.0 -b release/0.20.0 origin/next
   ```

3. Replace `origin/main` with the prepared line using a guarded update. The
   lease must name the pre-flip `main` SHA recorded in step 1:

   ```bash
   git push origin --force-with-lease=main:"$preflip_main" release/0.20.0:main
   ```

   Do not delete `next`. After the flip, `main` is the stable branch and `next`
   is the dev-only branch.

4. Cut and push the annotated `v0.20.0` tag from the `main` release line:

   ```bash
   git fetch origin main
   git switch release/0.20.0
   git reset --hard origin/main
   git tag -a v0.20.0 -F <changelog-file>
   git push origin v0.20.0
   ```

Before pushing the tag, verify `.github/workflows/release.yml` and
`.goreleaser.yaml` no longer stamp or pin the stable release to `next`. Those
release-mechanics changes are part of the main-flip release work, not this docs
reconciliation.

## What the Tag Push Does

Once the main-flip release mechanics are in place, `release.yml` runs one
goreleaser job on `macos-latest` that:

- cross-builds the darwin arm64 and amd64 tarballs plus `checksums.txt`,
  stamping `git describe --tags` into `internal/cli.Version`;
- publishes the GitHub Release with those assets;
- bumps the `spacedock-dev/homebrew-tap` cask via `HOMEBREW_TAP_TOKEN`;
- stamps the plugin manifests' `version` on `main`;
- keeps the marketplace entry serving the stable plugin from `main`.

The tag push is therefore enough to publish the stable release once the release
config targets `main`. A manual release-prep commit is still useful because it
produces a reviewable annotated-tag changelog and manifest diff before the tag
is pushed.

## Cutting a Stable Release After 0.20.0

1. Ensure all release content is merged to `main`. Choose the version `X.Y.Z`.

2. Create a release worktree off `main`:

   ```bash
   git worktree add .worktrees/release-X.Y.Z -b release/X.Y.Z origin/main
   ```

3. Bump the version stamps with the release tool, then commit. `stamp-version`
   writes the release `X.Y.Z` into the plugin manifests; `bump-calendar` advances
   the marketplace entry's separate `0.0.YYYYMMDDNN` calendar key:

   ```bash
   go run ./cmd/spacedock-release stamp-version X.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json
   go run ./cmd/spacedock-release bump-calendar .claude-plugin/marketplace.json
   git commit -m "release: bump version to spacedock@X.Y.Z" -- .claude-plugin/plugin.json .codex-plugin/plugin.json .claude-plugin/marketplace.json
   ```

4. Write a changelog. Summarize the commits since the last tag into plain text:

   ```bash
   git log $(git describe --tags --abbrev=0 origin/main)..HEAD --oneline
   ```

   One sentence names the release theme, then user-value-led `- ` bullets name
   what upgrading gives users. Ignore workflow-state churn
   (dispatch/advance/archive/mod-block/pr/report commits).

5. Create the annotated tag locally:

   ```bash
   git tag -a vX.Y.Z -F <changelog-file>
   ```

6. Review:

   ```bash
   git show vX.Y.Z
   git diff origin/main..release/X.Y.Z
   ```

   To amend the changelog: `git tag -d vX.Y.Z` then re-tag.

7. Publish after confirmation:

   ```bash
   git push origin release/X.Y.Z:main
   git push origin vX.Y.Z
   ```

8. Clean up:

   ```bash
   git worktree remove .worktrees/release-X.Y.Z
   git branch -d release/X.Y.Z
   ```

## Dev-Only `next` Publishing

Keep `next` for development. Source builds may use
`go install github.com/spacedock-dev/spacedock/cmd/spacedock@next`, local
checkouts may use `--plugin-dir`, and the deliberate `next-publish` workflow may
bump the marketplace calendar key for dev testers.

Do not send stable users to `next`. If a command or manifest uses `@next`, it is
a dev-only path.

## Notes

- Do not stamp the version via a pull request; the release branch and annotated
  tag are the release mechanism.
- macOS binaries are adhoc-signed, not yet notarized; the Homebrew cask's
  postflight strips the `com.apple.quarantine` xattr as the interim Gatekeeper
  fix until Developer-ID notarization lands.
- This flow is the v1 adaptation of the original `scripts/release.sh` from the
  upstream spacedock plugin repo.
