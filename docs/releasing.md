# Releasing Spacedock

Stable releases are cut from `main`. `next` remains a dev-only branch for source
builds and pre-stable publishing tests. Do not present `next` as the stable
marketplace source.

## What the Tag Push Does

`release.yml` runs one goreleaser job on `macos-latest` that:

- cross-builds darwin and linux (arm64 + amd64) tarballs plus `checksums.txt`,
  stamping `git describe --tags` into `internal/cli.Version`, for BOTH channels —
  a stable build (`cli.devBranch=main`) and an edge build (`cli.devBranch=next`);
- publishes the GitHub Release with those assets;
- bumps BOTH `spacedock-dev/homebrew-tap` casks (`spacedock` stable +
  `spacedock@next` edge) via `HOMEBREW_TAP_TOKEN`;
- stamps the plugin manifests' `version` on `main`.

The marketplace manifest no longer lives in the plugin branch. It is the
standalone `spacedock-dev/marketplace` repo, whose entries pin each channel by
ref — `spacedock` (stable, repointed to the released tag) and `spacedock-edge`
(edge, tracking `next`). Repointing the stable entry is a commit in that repo,
not a manifest stamp on `main`.

The tag triggers the release, but goreleaser publishes only after the `e2e-gate`
job confirms the tagged commit has a green Runtime Live E2E run (or a recorded
`SPACEDOCK_E2E_GATE_WAIVER`) — so the commit you tag must be one that has been
exercised green. A manual release-prep commit is still useful because it produces
a reviewable annotated-tag changelog and manifest diff before the tag is pushed.

## Cutting a Stable Release

1. Ensure all release content is merged to `main`. Choose the version `X.Y.Z`.

2. Create a release worktree off `main`:

   ```bash
   git worktree add .worktrees/release-X.Y.Z -b release/X.Y.Z origin/main
   ```

3. Bump the version stamps with the release tool, then commit. `stamp-version`
   writes the release `X.Y.Z` into the plugin manifests:

   ```bash
   go run ./cmd/spacedock-release stamp-version X.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json
   git commit -m "release: bump version to spacedock@X.Y.Z" -- .claude-plugin/plugin.json .codex-plugin/plugin.json
   ```

   The marketplace entry is no longer stamped here — repoint the `spacedock`
   stable entry's ref in the standalone `spacedock-dev/marketplace` repo.

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
