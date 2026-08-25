# Releasing Spacedock

Stable releases are cut from `main`. `next` remains a dev-only branch for source
builds (`go install …@next`, `--plugin-dir`) — it is NOT a re-pull source for any
installer. The edge marketplace entry tracks `main` directly, the same branch
prereleases are tagged on, so `main`'s tip already carries the latest edge
content once a `-pre` tag lands on it. Do not present `next` as the edge
marketplace source.

## What the Tag Push Does

`release.yml` runs one goreleaser job on `macos-latest` that:

- cross-builds darwin and linux (arm64 + amd64) tarballs plus `checksums.txt`,
  stamping `git describe --tags` into `internal/cli.Version`, for BOTH channels —
  a stable build (`cli.devBranch=main`) and an edge build (`cli.devBranch=next`).
  The asset name carries the channel: the edge tarball always ends `_edge`; the
  stable tarball is unsuffixed on a stable tag but ends `_stable` on a `-pre`
  tag, so a prerelease never offers a default-named asset. Nothing consumes a
  prerelease's `_stable` tarball — it exists only because the cask pipe cannot
  tolerate a skipped build;
- publishes the GitHub Release with those assets;
- bumps BOTH `spacedock-dev/homebrew-tap` casks (`spacedock` stable +
  `spacedock@next` edge) via `HOMEBREW_TAP_TOKEN`;
- advances the stable channel ref to the tagged commit. See below.
- stops the cut before publication if the tag is older than the release that
  `stable` points at. See "The Stable Regression Gate" below.
- stamps the plugin manifests and the FO prose pin on `main`, but only on a
  latest-line stable tag. The stamp writes the `X.(Y+1).0-pre0` version that the
  same job auto-cuts. An old-line tag does not change `main`.
- on a latest-line stable tag, auto-cuts a `vX.(Y+1).0-pre0` prerelease tag so
  the edge binary catches up to the new release line within minutes (see
  "Advancing the Edge Line" below).

The marketplace manifest no longer lives in the plugin branch. It is the
standalone `spacedock-dev/marketplace` repo, where each channel is a branch with
its own root `marketplace.json`: the repo root (a marketplace named `spacedock`,
the stable entry pinned to `source.ref=stable`) and the `edge` branch (a
marketplace named `spacedock-edge`, entry tracking `main`). The channel lives in
the marketplace NAME, so a binary adds the channel's branch source —
`spacedock-dev/marketplace` for stable, `spacedock-dev/marketplace@edge` for edge —
and the marketplace it registers carries the matching name.

A hand-edit in the marketplace repo does not repoint the stable entry per
release. `stable` is a MOVING BRANCH in this repo. After the tag fires,
release.yml resolves the tagged commit and pushes it to `refs/heads/stable`. A
fresh `spacedock@spacedock` install resolves whatever `stable` points at. That
push is what publishes the release to the stable channel. No marketplace-repo
commit is necessary.

The `main` stamp is no longer part of the stable-tag path. The `edge-advance`
job stamps `main`, and only on a latest-line tag. The stamp writes the next
prerelease version. The tagged commit keeps the release version. The two values
are different by design, so the stamp always makes one commit on `main`.

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
   go run ./cmd/spacedock-release stamp-version X.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json skills/first-officer/references/first-officer-shared-core.md
   git commit -m "release: bump version to spacedock@X.Y.Z" -- .claude-plugin/plugin.json .codex-plugin/plugin.json skills/first-officer/references/first-officer-shared-core.md
   git push origin release/X.Y.Z:main
   ```

   The stamped commit's manifest version now equals `X.Y.Z`, and the FO
   shared-core's release-stamped minor literal now equals `X.Y` (D5). Guard
   both before tagging — a tag whose semver disagrees with the tagged commit's
   manifest, or whose major.minor disagrees with the tagged commit's
   prose-stamped minor, is the pre-stamp inversion the gate has historically
   caught (v0.20.0 tagged a commit whose `plugin.json` still read 0.19.9):

   ```bash
   go run ./cmd/spacedock-release manifest-tag-gate vX.Y.Z .claude-plugin/plugin.json .codex-plugin/plugin.json skills/first-officer/references/first-officer-shared-core.md
   ```

   The marketplace entry is not stamped or repointed here: release.yml advances
   the moving `stable` ref to this commit after the tag fires (see "What the Tag
   Push Does").

4. Green that exact commit. Capture the release SHA, then dispatch Runtime Live
   E2E on it and wait for a `conclusion: success` run — the `e2e-gate` matches a
   green run to the tagged SHA, and the workflow is `workflow_dispatch`-only, so
   nothing greens the commit unless you dispatch it. A lane that flakes can be
   re-run to green (`gh run rerun <run-id> --failed`); the re-run-to-green run
   satisfies the gate — no fresh dispatch needed:

   ```bash
   REL_SHA=$(git rev-parse HEAD)
   gh workflow run "Runtime Live E2E" --ref main
   gh run watch "$(gh run list --workflow 'Runtime Live E2E' --branch main --limit 1 --json databaseId --jq '.[0].databaseId')"
   ```

   `REL_SHA` is the stamped commit from step 3 (the worktree HEAD you just pushed);
   it is the SHA this run must go green on, and the SHA you tag in step 6. (For an
   emergency cut when the live matrix is unavailable, the gate accepts the auditable
   `SPACEDOCK_E2E_GATE_WAIVER` repo variable instead.)

   **Equivalent prior green.** The gate is mechanical — it matches a run to the
   exact SHA — but the evidence question is about trees, not SHAs. Before burning a
   fresh matrix, check what actually changed since the last greened commit:

   ```bash
   git diff --name-only <greened-sha> "$REL_SHA"
   ```

   When that diff contains only the stamp manifests (`.claude-plugin/plugin.json`,
   `.codex-plugin/plugin.json`, an unchanged-minor shared-core stamp) and files no
   live lane loads (this repo's own workflow docs under `docs/dev/`), the prior
   green already proved this tree: a fresh run re-buys the same evidence plus one
   roll of host stochasticity. In that case the cutter MAY satisfy the gate with
   `SPACEDOCK_E2E_GATE_WAIVER` instead of a fresh dispatch, citing the equivalent
   run id and the `git diff --name-only` output in the waiver's audit trail
   (captain ruling, 2026-08-15: the v0.27.0-pre5 delta was three files, zero
   live-lane bytes). Any product file in the diff — anything a lane builds, loads,
   or drives — voids the equivalence and a fresh dispatch is required.

5. Write a changelog. For a stable release, inventory everything from the prior
   stable tag through the candidate, including changes that first appeared in a
   prerelease. Inventory first; compress it into public notes second:

   ```bash
   PREV_STABLE="$(git tag --list --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n 1)"
   git log "$PREV_STABLE..$REL_SHA" --oneline
   ```

   Use commits and PRs as source evidence, but write for users: one theme
   sentence followed by 4–7 `- ` bullets that say what upgrading lets them do,
   avoid, or rely on. Omit PR numbers, task IDs, filenames, tests and fixture or
   oracle names, prose counts, CI wiring, and workflow-state churn. Finish with
   one stable-to-stable `Full changelog` comparison link when useful.

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

## Advancing the Edge Line

The edge marketplace entry resolves `main` directly — the same branch stable
AND prerelease tags are both cut from — so `main`'s tip already carries the
latest edge content once a tag lands on it: `claude plugin update` / `codex`
re-pull on a manifest-version difference at that source ref (no separate
tracking branch, no calendar-key re-pull mechanism to maintain). The only thing
a tag push still automates for the edge line is auto-cutting the NEXT
prerelease after a latest-line stable release, in a job that `needs: goreleaser`
(a sibling job, so a rare conflict here cannot unwind or block the release that
already published):

- **Stable (`vX.Y.Z`) tag, latest line:** an ANNOTATED `vX.(Y+1).0-pre0` tag is
  auto-created **on the greened release commit** and pushed over SSH with a
  dedicated write **deploy key** (`EDGE_RELEASE_DEPLOY_KEY`), scoped to this
  repo. That prerelease tag's own release run reuses the greened commit's
  e2e-gate pass, builds+publishes the `X.(Y+1)`-minor edge binary, and bumps the
  `spacedock@next` cask — so the edge binary's minor catches up to the skills'
  gate line within minutes instead of waiting for the next hand-cut prerelease.
  The auto-tag MUST be annotated with a non-empty body (the release-notes
  extraction step rejects a lightweight tag). The push MUST use a
  trigger-capable credential: a `GITHUB_TOKEN` push — and, as observed on the
  v0.25.0 and v0.26.0 cuts, the cross-repo tap PAT — does NOT create the pre0
  `release.yml` run, so the step pushes with the deploy key and then **verifies
  a run was created for the pre0 tag, failing `edge-advance` loudly if none
  appears** rather than leaving the edge binary silently behind. Re-running the
  job after such a failure is safe: an existing pre0 tag is never re-minted or
  moved (the step checks `refs/tags/` before tagging), the push of an
  already-landed tag is a no-op, and the run-verification poll executes again —
  so the re-run turns green once the pre0 run exists, and still fails loudly if
  the credential remains suppressed. Expect two GitHub releases per stable cut.
  After the job verifies the pre0 run, the same job stamps the manifests and the
  FO pin on `main` to the pre0 version. An edge install then passes the FO version
  gate with no human step. If a merge to `main` happens during the run, the push
  can fail. Run the job again, because the stamp gives the same result each time.
- **Old-line / patch (`vX.Y.1`, cut after a newer `X.(Y').Z` line already
  shipped):** the decision step compares this tag's own version against the
  highest OTHER existing bare stable tag in the repo and prints `skip`, logged
  as a `::notice::`; the auto-pre0 step is gated on that decision, so it does
  NOT run. Without this check, the patch would auto-cut a WRONG, lower
  `vX.(Y+1).0-pre0` tag whose release run would bump the `spacedock@next` cask
  DOWN to an older minor. The same decision also gates the `main` stamp. An
  old-line tag cannot move `main` to a lower version.

The `vX.(Y+1).0-pre0` release run does not recurse: its `-pre0` tag is itself a
prerelease (hyphenated), so both the decision step and the auto-pre0 step —
which carry `!contains(github.ref, '-')` — skip on that tag's own run. goreleaser
stamps the edge binary from the highest tag pointing at the commit (`git tag
--points-at HEAD --sort -version:refname`), which is the `-pre0` tag, so the
`X.(Y+1)` binary version is correct without an override; a
`GORELEASER_CURRENT_TAG` pin is set on the goreleaser step as belt-and-suspenders.

Both edge binary install paths track those prerelease tags: the `spacedock@next`
cask, bumped on every tag including prereleases, and
`SPACEDOCK_CHANNEL=edge … install.sh | sh`, which resolves the newest release
including prereleases and fetches that tag's `_edge` asset. The default stable
path resolves `/releases/latest`, which excludes prereleases — so a default
install on an edge machine lands a lower minor and aborts at first-officer boot.
Casks are macOS-only; on Linux the script is the only edge binary path.

## The Stable Regression Gate

A stable tag that is older than the release that `stable` points at can damage
two binary channels. goreleaser marks each bare `vX.Y.Z` tag as a full release, so
GitHub moves `/releases/latest` to it. goreleaser also bumps the stable Homebrew
cask, because `skip_upload: auto` skips only a prerelease. An old tag moves both
channels DOWN, and every job stays green. The `e2e-gate` job stops this before
goreleaser starts.

- The gate reads the version in the `.claude-plugin/plugin.json` file that
  `stable` points at.
- If the tag is older than that version, the gate fails the run. goreleaser does
  not start, because it needs the `e2e-gate` job.
- If the tag is the same version or newer, the gate lets the run continue. A
  re-run of a release that already reached `stable` is still possible.
- On the first stable release the `stable` ref does not exist. The gate writes a
  notice and lets the run continue.
- If the gate cannot read the `stable` ref, it fails the run. A read failure is
  not a permission to publish.

The gate does not make a patch line deliverable. A patch that is newer than the
version `stable` points at still publishes, and `stable` then leaves the history
of `main`. The next latest-line release cannot advance `stable` after that. Do not
cut a patch line until that work is complete.

## Dev-Only `next` Publishing

Keep `next` for development. Source builds may use
`go install github.com/spacedock-dev/spacedock/cmd/spacedock@next`, and local
checkouts may use `--plugin-dir`. `next` is not consulted by any installer
(the edge marketplace entry resolves `main` directly — see "Advancing the Edge
Line" above), so there is no re-pull mechanism to run for it; `next` is a
source-build convenience branch only. A `go install …@vX.Y.Z` proxy build
carries no ldflags, so it self-reports `X.Y.Z+dev` (the tagged manifest at the
module-proxy commit equals the tag) — gates correctly under minor-version
coupling; the `+dev` suffix on an otherwise-tagged build is a cosmetic oddity,
not a compatibility issue.

Do not send stable users to `next`. If a command or manifest uses `@next`, it is
a dev-only path.

## Notes

- Do not stamp the version via a pull request; the release branch and annotated
  tag are the release mechanism.
- macOS binaries are adhoc-signed, not yet notarized; the Homebrew cask's
  postflight strips the `com.apple.quarantine` xattr as the interim Gatekeeper
  fix until Developer-ID notarization lands.
- The `spacedock-dev/marketplace` repo's `edge` entry's `version` field is
  display metadata only — inert on both current hosts (`claude plugin update`
  keys on the plugin MANIFEST version at the source ref; codex re-clones on add
  and displays the manifest version), and a tag-vs-branch channel structure
  cleanup want. That repo is standalone and unreachable from here, so the
  cleanup is deferred to a marketplace-repo task and is NOT part of this flow.
- This flow is the v1 adaptation of the original `scripts/release.sh` from the
  upstream spacedock plugin repo.
- Prerelease releases cut before the `_stable` suffix landed shipped the
  stable-stamped binary under the unsuffixed default asset name. An operator
  who installed one holds a stable-channel binary at a prerelease version: it
  re-ensures the `spacedock@spacedock` plugin on every launch. Detection:
  `spacedock --version` prints a `Channel:` line on new binaries; on older
  ones `go version -m $(command -v spacedock)` shows `cli.devBranch=main`.
  Remedy: reinstall from the release's `_edge` tarball or
  `brew install spacedock-dev/tap/spacedock@next`.
