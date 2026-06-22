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
- stamps the plugin manifests' `version` on `main`, then advances the stable
  channel ref (see below).

The marketplace manifest no longer lives in the plugin branch. It is the
standalone `spacedock-dev/marketplace` repo, where each channel is a branch with
its own root `marketplace.json`: the repo root (a marketplace named `spacedock`,
the stable entry pinned to `source.ref=stable`) and the `edge` branch (a
marketplace named `spacedock-edge`, entry tracking `next`). The channel lives in
the marketplace NAME, so a binary adds the channel's branch source —
`spacedock-dev/marketplace` for stable, `spacedock-dev/marketplace@edge` for edge —
and the marketplace it registers carries the matching name.

The stable entry is NOT repointed per release by a hand-edit in the marketplace
repo. `stable` is a MOVING BRANCH in THIS repo: after the manifest stamp,
release.yml resolves the tagged commit as `RELEASE_COMMIT` and runs
`git push origin "$RELEASE_COMMIT:refs/heads/stable"` (release.yml's "Stamp plugin
manifests" step), advancing `stable` to that exact tagged commit. A fresh
`spacedock@spacedock` install resolves whatever `stable` points at, so that push
is what publishes the release to the stable channel — no marketplace-repo commit.

The post-tag manifest stamp is idempotent: when the tagged commit ALREADY carries
the release version in `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json`
(it does, under the stamp-then-tag ordering below), the stamp step finds no diff
and commits nothing; it still advances `stable` to that commit.

The tag triggers the release, but goreleaser publishes only after the `e2e-gate`
job confirms the tagged commit has a green Runtime Live E2E run (or a recorded
`SPACEDOCK_E2E_GATE_WAIVER`) — so the commit you tag must be one that has been
exercised green. The gate binds the green run to the EXACT tagged commit SHA
(`run.HeadSha == tagged-commit`, `internal/release/e2egate.go`), and Runtime Live
E2E is `workflow_dispatch`-only, so it never runs automatically on a fresh commit.
This is why the cut below stamps and pushes the release commit to `main` FIRST,
greens THAT commit, and only then tags it — a tag on any commit that has no green
run for its own SHA is blocked by the gate.

## Cutting a Stable Release

The cutter tags a commit that ALREADY has both a matching manifest stamp and a
green Runtime Live E2E run for its exact SHA. Stamp and push the release commit to
`main`, green that commit, then tag it.

1. Ensure all release content is merged to `main`. Choose the version `X.Y.Z`.

2. Create a release worktree off `main`:

   ```bash
   git worktree add .worktrees/release-X.Y.Z -b release/X.Y.Z origin/main
   ```

3. Stamp the plugin manifests to `X.Y.Z`, commit, and push to `main`. This is the
   commit the tag will name, so it must land on `main` and be greened BEFORE the
   tag — a fresh post-tag commit would have no green Runtime Live E2E run for its
   SHA and the `e2e-gate` would block the cut.

   ```bash
   go run ./cmd/spacedock-release stamp-version X.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json
   git commit -m "release: bump version to spacedock@X.Y.Z" -- .claude-plugin/plugin.json .codex-plugin/plugin.json
   git push origin release/X.Y.Z:main
   ```

   The stamped commit's manifest version now equals `X.Y.Z`. Guard that before
   tagging — a tag whose semver disagrees with the tagged commit's manifest is the
   pre-stamp inversion the gate has historically caught (v0.20.0 tagged a commit
   whose `plugin.json` still read 0.19.9):

   ```bash
   go run ./cmd/spacedock-release manifest-tag-gate vX.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json
   ```

   The marketplace entry is not stamped or repointed here: release.yml advances
   the moving `stable` ref to this commit after the tag fires (see "What the Tag
   Push Does").

4. Green that exact commit. Capture the release SHA, then dispatch Runtime Live
   E2E on it and wait for a `conclusion: success` run — the `e2e-gate` matches a
   green run to the tagged SHA, and the workflow is `workflow_dispatch`-only, so
   nothing greens the commit unless you dispatch it:

   ```bash
   REL_SHA=$(git rev-parse HEAD)
   gh workflow run "Runtime Live E2E" --ref main
   gh run watch "$(gh run list --workflow 'Runtime Live E2E' --branch main --limit 1 --json databaseId --jq '.[0].databaseId')"
   ```

   `REL_SHA` is the stamped commit from step 3 (the worktree HEAD you just pushed);
   it is the SHA this run must go green on, and the SHA you tag in step 6. (For an
   emergency cut when the live matrix is unavailable, the gate accepts the auditable
   `SPACEDOCK_E2E_GATE_WAIVER` repo variable instead.)

5. Write a changelog. Summarize the commits since the last tag into plain text:

   ```bash
   git log $(git describe --tags --abbrev=0 origin/main)..HEAD --oneline
   ```

   One sentence names the release theme, then user-value-led `- ` bullets name
   what upgrading gives users. Ignore workflow-state churn
   (dispatch/advance/archive/mod-block/pr/report commits).

6. Create the annotated tag on the greened commit:

   ```bash
   git tag -a vX.Y.Z -F <changelog-file> "$REL_SHA"
   ```

   `$REL_SHA` is the stamped, greened commit captured in step 4; tag THAT SHA so
   the `e2e-gate` finds its matching green run. (Tagging the captured SHA, not
   `git rev-parse origin/main`, keeps the target stable even if a stray `git fetch`
   moves `origin/main` between steps.)

7. Review:

   ```bash
   git show vX.Y.Z
   ```

   To amend the changelog: `git tag -d vX.Y.Z` then re-tag the same commit.

8. Publish after confirmation:

   ```bash
   git push origin vX.Y.Z
   ```

   The release commit is already on `main` from step 3; pushing the tag fires
   `release.yml`, which gates on the green run for this commit and then publishes.

9. Clean up:

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
- The `spacedock-dev/marketplace` repo's `next` entry carries a dead `version`
  field and a tag-vs-branch channel structure that want cleanup. That repo is
  standalone and unreachable from here, so the cleanup is deferred to a
  marketplace-repo task and is NOT part of this flow.
- This flow is the v1 adaptation of the original `scripts/release.sh` from the
  upstream spacedock plugin repo.
