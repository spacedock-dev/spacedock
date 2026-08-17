# Releasing Spacedock

Stable releases are cut from `main`. `next` remains a dev-only branch for source
builds and pre-stable publishing tests. Do not present `next` as the stable
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
- stamps the plugin manifests' `version` on `main`, then advances the stable
  channel ref (see below).
- advances the `next` edge line to match the release — reconciled on a
  prerelease tag, reconciled plus bumped to the post-release dev pre-version on
  a stable tag (see "Advancing the Edge Line" below).

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

## Advancing the Edge Line (`next`)

Every tag push also advances `next` — the branch the `spacedock-edge`
marketplace resolves — in a job that `needs: goreleaser` (a sibling job, so a
rare conflict here cannot unwind or block the release that already published):

- **Prerelease (`-pre`) tag:** `next` is reconciled to the tagged commit's
  content — `git merge -X theirs "$RELEASE_COMMIT"`, favoring the release over
  whatever `next` had drifted to — then the marketplace calendar key is bumped
  (`spacedock-release bump-calendar`) so `claude plugin update` / `codex`
  re-pull. This is the automated form of the manual reconcile the 0.24.0-pre1
  cut required (`next` had drifted 40 commits behind `main`, hard-blocking
  `spacedock codex` on a binary/plugin version-compat check).
- **Stable (`vX.Y.Z`) tag, latest line:** `next` is reconciled the same way,
  stamped PAST the release to `X.(Y+1).0-pre1`
  (`spacedock-release dev-preversion X.Y.Z`) so the edge line never masquerades
  as the stable version it just shipped, and the calendar key is bumped — then an
  ANNOTATED `vX.(Y+1).0-pre0` tag is auto-created **on the greened release commit**
  and pushed over SSH with a dedicated write **deploy key**
  (`EDGE_RELEASE_DEPLOY_KEY`), scoped to this repo. That prerelease tag's own
  release run reuses the greened commit's e2e-gate pass, builds+publishes the
  `X.(Y+1)`-minor edge binary, and bumps the `spacedock@next` cask — so the edge
  binary's minor catches up to the skills' gate line within minutes instead of
  waiting for the next hand-cut prerelease. The auto-tag MUST be annotated with a
  non-empty body (the release-notes extraction step rejects a lightweight tag).
  The push MUST use a trigger-capable credential: a `GITHUB_TOKEN` push — and, as
  observed on the v0.25.0 and v0.26.0 cuts, the cross-repo tap PAT — does NOT
  create the pre0 `release.yml` run, so the step pushes with the deploy key and
  then **verifies a run was created for the pre0 tag, failing `edge-advance`
  loudly if none appears** rather than leaving the edge binary silently behind.
  Expect two GitHub releases per stable cut.
- **Old-line / patch (`vX.Y.1`, or any tag whose target edge version is not
  strictly greater than `next`'s current manifest version):** the whole
  `edge-advance` job SKIPS (`spacedock-release edge-advance-decision` prints
  `skip`, logged as a `::notice::`; every downstream step is gated on the
  decision). `next`'s tip — content, manifests, gate line, and the marketplace
  calendar key — is left untouched, so no edge installer re-pulls. The patch
  updates only the stable cask; its fix reaches edge through the normal
  `main`→`next` flow, never through a `-X theirs` reconcile that would clobber
  `next`'s newer `(Y+1)`-line content or rewind its manifest/gate line. The
  auto-pre0 step is stable-latest-line-only and decision-gated, so a patch never
  attempts a colliding pre0 tag.

The `vX.(Y+1).0-pre0` release run does not recurse: its `-pre0` tag routes to the
prerelease path, whose `edge-advance-decision` skips (`pre0 < next`'s `pre1`), and
the auto-pre0 step runs only on the stable path — so it neither re-tags nor
rewinds `next`. goreleaser stamps the edge binary from the highest tag pointing at
the commit (`git tag --points-at HEAD --sort -version:refname`), which is the
`-pre0` tag, so the `X.(Y+1)` binary version is correct without an override; a
`GORELEASER_CURRENT_TAG` pin is set on the goreleaser step as belt-and-suspenders.

The reconcile is a merge, never a reset or force-push: the previous `next` tip
is always a first-parent ancestor of the new commit, so `git push origin
<sha>:next` is a plain fast-forward. A real conflict (two sides changing the
same file in incompatible, non-superseded ways) fails the step loudly instead
of guessing — the same manual reconciliation this replaces remains the escape
hatch.

## Dev-Only `next` Publishing

Keep `next` for development. Source builds may use
`go install github.com/spacedock-dev/spacedock/cmd/spacedock@next`, local
checkouts may use `--plugin-dir`, and the deliberate `next-publish` workflow may
bump the marketplace calendar key for dev testers. A `go install …@vX.Y.Z`
proxy build carries no ldflags, so it self-reports `X.Y.Z+dev` (the tagged
manifest at the module-proxy commit equals the tag) — gates correctly under
minor-version coupling; the `+dev` suffix on an otherwise-tagged build is a
cosmetic oddity, not a compatibility issue.

Every release tag now advances `next` and bumps its calendar key automatically
(see "Advancing the Edge Line" above); `next-publish` stays for an out-of-band
re-pull between releases (e.g. a `next`-only fix that isn't worth a full cut).

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
- Prerelease releases cut before the `_stable` suffix landed shipped the
  stable-stamped binary under the unsuffixed default asset name. An operator
  who installed one holds a stable-channel binary at a prerelease version: it
  re-ensures the `spacedock@spacedock` plugin on every launch. Detection:
  `spacedock --version` prints a `Channel:` line on new binaries; on older
  ones `go version -m $(command -v spacedock)` shows `cli.devBranch=main`.
  Remedy: reinstall from the release's `_edge` tarball or
  `brew install spacedock-dev/tap/spacedock@next`.
