---
id: tdng3g6fe5p4y08tj5dphywc
title: install.sh resolves releases/latest and silently breaks edge version parity
status: implementation
source: "Captain CL fresh-install VM experience report, 2026-08-24: fresh VM tracking edge would have gotten v0.26.0 from install.sh with no warning"
started: 2026-08-24T15:09:44Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:tdng3g6fe5p4y08tj5dphywc:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:tdng3g6fe5p4y08tj5dphywc-backlog-1
              briefing:
                id: briefing:tdng3g6fe5p4y08tj5dphywc:backlog:attempt-1:revision-1
                digest: sha256:d89298a1f4bcf4e18e6fcb6b6dd823249563f6367f773312dc641a5cf342c939
                request-digest: sha256:c7ed6226a494766f117a58b14240301160dcfa1c0c3826598321c35b70a7932d
                room-ref: ./install-sh-edge-prerelease-parity/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:tdng3g6fe5p4y08tj5dphywc:backlog:1
                briefing: briefing:tdng3g6fe5p4y08tj5dphywc:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-24T15:07:39.462852Z"
                decision: approve
                reason: 'Captain CL in chat, 2026-08-24: ''dispatch td'' - accepts the seed direction into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:tdng3g6fe5p4y08tj5dphywc:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:tdng3g6fe5p4y08tj5dphywc-ideation-1
              briefing:
                id: briefing:tdng3g6fe5p4y08tj5dphywc:ideation:attempt-1:revision-1
                digest: sha256:0541b9d8900726000b867f26ad061d756f9f9c7e19e4769ab03b5a3ae82b9101
                request-digest: sha256:74728ccf32c5eaabc527990b126a9876cd7374ee976755a2f9345731d4a87837
                room-ref: ./install-sh-edge-prerelease-parity/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-24T15:33:11.979924Z"
                reason: 'Captain scope amendment 2026-08-24: fold the docs/site/get-started/install.md:50 edge channel comment correction (formerly 765 scope) into this task''s doc diff; body will change, bound snapshot stale'
            - id: gate-attempt:tdng3g6fe5p4y08tj5dphywc-ideation-2
              briefing:
                id: briefing:tdng3g6fe5p4y08tj5dphywc:ideation:attempt-2:revision-1
                digest: sha256:21d4a15f84bf8203246499afcfc814e8ac0f2879e1a920c55ade27e0b1aca92d
                request-digest: sha256:8fe62dafdf31898791699cb874a0293872fadcee989ff0e32c3072f16da23a62
                room-ref: ./install-sh-edge-prerelease-parity/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:tdng3g6fe5p4y08tj5dphywc:ideation:2
                briefing: briefing:tdng3g6fe5p4y08tj5dphywc:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-24T15:40:44.28296Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''approve'' at re-presented attempt-2 (digest 21d4a15f) - accepts the SPACEDOCK_CHANNEL design with the folded install.md amendment'
              application:
                target-stage: implementation
                state: consumed
---

`install.sh:90` resolves `api.github.com/repos/$REPO/releases/latest`, and GitHub's `/latest` excludes prereleases. Verified 2026-08-24: it resolves `v0.26.0` while the edge-parity target is 0.27.0-pre8 — a fresh machine following the binary-install path silently lands a binary one minor behind the edge plugin, which the FO boot version gate (requires minor 0.27) then aborts on. The skew is silent at install time: nothing prints the resolved version or hints that prereleases are excluded. This defeated the exact parity claim a fresh-VM prototype run was testing.

## Problem

`install.sh` has no notion of a release channel, while every other install surface does. The binary it lands is always the stable one, and it never says so.

Two independent gaps produce that outcome, and the second was not visible from the seed:

1. **Tag resolution is stable-only.** `install.sh:90` resolves `api.github.com/repos/$REPO/releases/latest`, which excludes prereleases by definition. It resolves `v0.26.0` while the edge line is at `v0.27.0-pre8`.
2. **Asset naming has no channel, and on a prerelease the channel-less name does not exist.** `asset_name()` (`install.sh:100`) always builds `spacedock_<ver>_<os>_<arch>.tar.gz`. Every release publishes a *pair* of per-arch archives — the edge-stamped binary always carries an `_edge` suffix, and on a prerelease tag the stable one carries `_stable` (`.goreleaser.yaml:83`, deliberately, so no prerelease publishes a bare-looking asset). So the unsuffixed name only exists on stable tags.

The consequence of (2) is the load-bearing correction to the seeded design: **fixing tag resolution alone does not work.** With the newest prerelease resolved and today's `asset_name()`, the installer requests `spacedock_0.27.0-pre8_darwin_arm64.tar.gz`, which 404s (spike leg F). An edge install path needs the channel threaded through *naming* as well as *resolution*.

It also means `install.sh` cannot install an edge-stamped binary at **any** tag today, not merely at prereleases — the `_edge` asset is never the name it builds. Homebrew has a working edge variant (`spacedock@next`, `.goreleaser.yaml:176`, `skip_upload: false` so it bumps on every prerelease), but Homebrew casks are macOS-only. **A Linux machine tracking edge has no scripted binary path at all.**

The skew is silent throughout. Nothing prints the resolved tag, so a fresh machine installs a 0.26 binary, installs edge skills, and only discovers the mismatch when the first-officer boot gate — which requires binary minor 0.27 (`skills/first-officer/references/first-officer-shared-core.md:9`) — aborts. That defeated the exact parity claim a fresh-VM prototype run was testing.

## Spike: how an edge machine resolves the newest prerelease

Run in ideation against the live repo and the live GitHub API, unauthenticated, on darwin/arm64. A throwaway patched copy of `install.sh` was the probe; nothing was committed to the code tree. Legs are lettered as run.

**Resolution mechanism (the riskiest claim).** `GET /repos/$REPO/releases?per_page=1` returns the newest release *including* prereleases: `v0.27.0-pre8`. `/releases/latest` returns `v0.26.0` for the same repo at the same moment (legs A, B). Both are exactly **one** unauthenticated request, so the mechanism is cost-neutral: measured `x-ratelimit-limit: 60`, `x-ratelimit-resource: core`, consuming 1 per run either way (leg A). No token, no auth, no rate-limit regression. Ordering is `created_at` descending and `tag_name` appears only at release level, never inside an asset object, so the existing `grep '"tag_name"' | head -n 1 | sed` parse works unchanged on the array response (legs A, C). Drafts do not appear in an unauthenticated list, so there is no draft hazard on the anonymous path.

**Asset naming (the finding that redirects the design).** `v0.27.0-pre8` publishes eight archives — `{darwin,linux}_{amd64,arm64}` × `{_edge,_stable}` — plus `checksums.txt` with a line for each, in the `<sha256>␣␣<filename>` shape install.sh's `awk '$2 == f'` already parses (legs D, G). The name install.sh builds today returns **HTTP 404** on that tag; the `_edge` name returns **200** (leg F). `v0.26.0` by contrast publishes the unsuffixed archive *and* an `_edge` one, so the unsuffixed name is the stable channel's name on a stable tag (leg E). Edge naming is therefore uniform (`_edge` always) while stable naming is conditional — and because the stable path resolves `/releases/latest`, which is never a prerelease, **the stable path never encounters the `_stable` variant and needs no naming change at all.**

**End-to-end proof.** With channel-aware resolution *and* naming, a clean empty prefix takes the edge path, passes the existing checksum gate untouched, and the installed binary self-reports (leg 1):

```
install.sh: resolved edge channel to v0.27.0-pre8
spacedock 0.27.0-pre8
Channel: edge (spacedock@spacedock-edge)
```

Minor 0.27 — the boot gate the old path was aborting on now passes. The same run on the default path lands `spacedock 0.26.0` from the unsuffixed asset, with no `Channel:` line at all (leg 2). The documented pipe shape works identically with the script on stdin: `cat install.sh | SPACEDOCK_CHANNEL=edge sh` (leg O).

**Failure and escape paths.** An unreachable or 403 API is already a clean fail-closed abort: exit 1, `could not resolve the latest release tag`, install dir never created (leg H) — and that message is the string `install_url_test.go:111` already skips on, so the live-test pattern carries over. A tag pin needs **no new mechanism**: the existing `SPACEDOCK_INSTALL_FROM=<release download base>` + `SPACEDOCK_INSTALL_VERSION` composes with the channel to install a chosen prerelease (leg I, installed `0.27.0-pre7` edge). That is the escape hatch if list ordering ever resolves the wrong tag.

**Two regressions checked, both clean.** Stable `SPACEDOCK_PRINT_TARGET=1` stdout is byte-identical to the shipped script's for the same host (leg L), and the new resolved-tag line goes to **stderr** and contains no `=`, so it cannot pollute the `key=value` parser in `install_url_test.go:117` (legs K, M). In a local dist directory the two channels' globs are disjoint: `spacedock_*_${os}_${arch}.tar.gz` does not match an `_edge` file, so the edge channel needs its own glob and the existing stable glob is untouched (leg, fixture dir).

**One defect found in the probe itself, which the implementation must not copy.** The probe used `case "$CHANNEL" in edge) … ;; *) …stable… ;; esac`. A typo'd value silently installs stable: `SPACEDOCK_CHANNEL=egde` resolved `v0.26.0` and would have installed it (leg P). That is precisely the silent-skew class this task exists to close, so the channel value must be **validated and rejected**, not defaulted through. Hence AC-6.

**Residual risk, accepted and declared.** Edge resolution trusts the API's documented `created_at`-descending order rather than sorting semver in POSIX shell. `created_at` is the tag's commit date, so the realistic bad case — a patch cut later from an older commit — still sorts below the newer prerelease. A pathological case (a prerelease tagged on an old commit) could mis-resolve. Not worth a prerelease-aware shell semver sorter (~25 lines of fragile `sh`): the mitigation is that the resolved tag is now **printed**, converting a mis-resolution from silent to visible, with the zero-new-code pin above as the override. Turning silence into visibility is this task's whole point; a wrong-but-announced tag is a different and lesser failure than a wrong-and-silent one.

## Proposed approach

One new environment variable, `SPACEDOCK_CHANNEL`, accepting `stable` (default) and `edge`, threaded through the three places the channel actually changes behavior, plus one line of disclosure.

1. **Validate the channel first.** Unrecognized value → `die` naming the accepted values. No silent default (spike leg P).
2. **Tag resolution** (`resolve_latest_tag`): edge → `/releases?per_page=1`; stable → today's `/releases/latest`, unchanged.
3. **Asset naming** (`asset_name`): edge → `spacedock_<ver>_<os>_<arch>_edge.tar.gz`; stable → today's unsuffixed name, unchanged.
4. **Local-dist glob**: edge → `spacedock_*_${os}_${arch}_edge.tar.gz`; stable → today's glob, unchanged. This keeps `SPACEDOCK_INSTALL_FROM=<dir>` usable for both channels and is what makes the channel testable offline against a goreleaser `--snapshot` dist, which publishes both archives.
5. **Disclose the resolved tag**: one stderr line, `install.sh: resolved <channel> channel to <tag>`, on the API-resolution path only, before any download. Not on the `SPACEDOCK_INSTALL_FROM` path, where the caller supplied the version and there is nothing to disclose.

The checksum gate, extract/install path, install dir, and PATH note are untouched; both channels flow through the same verify/install block, as both source paths do today.

**Why each mechanism, and the simplest alternative rejected.** Every item above serves AC-1 (an edge machine lands an edge-parity binary) except item 5, which serves AC-2 (the skew is never silent).

- *Documentation only — declare `install.sh` stable-only and send edge users to Homebrew.* Zero code, and it was the seed's third option. Insufficient: casks are macOS-only, so this leaves Linux edge machines with no scripted path, which is AC-1's core case. It also cannot satisfy AC-2, since the stable path stays silent.
- *A tag pin (`SPACEDOCK_INSTALL_TAG=v0.27.0-pre8`) instead of a channel.* Insufficient as the primary path: it requires a human to look up the newest prerelease, goes stale at every `-pre` cut, and cannot be scripted, while the FO abort names the required *minor*, not a tag. It is also redundant — the spike proved the pin already works through existing variables (leg I), so it is the documented escape hatch at zero new cost, not the mechanism.
- *Resolve by sorting all releases by semver in shell.* Rejected above as disproportionate; see residual risk.
- *Gate the resolved-tag line behind a verbose flag.* Rejected: silence is the defect. Unconditional costs one line and no one has to know to ask.

`SPACEDOCK_CHANNEL` is unused anywhere in the repo today, so the name is free. It is deliberately distinct from `SPACEDOCK_DEV_BRANCH`, which overrides the *binary's own* stamped `devBranch` at runtime (`internal/cli/cli.go:808`); this variable selects which *artifact to fetch* and is read only by `install.sh`.

One adjacent one-line correction, declared rather than folded in silently: `install.sh:6`'s usage comment advertises the `next` raw URL, while `internal/contract/contract.go:239` and `first-officer-shared-core.md:10` both use `main`, and `CLAUDE.md:25` states `next` "is not a re-pull source for any installer." The header is the outlier and the same class of defect (an install path naming the wrong ref).

Related but out of scope: the FO upgrade-hint journey (`fo-boot-upgrade-hint-latest-release`, d2k) handles the wrong-version abort *after* it happens. This task is about not installing the wrong version in the first place. The two meet at a shared need for a newest-release query, which d2k's body notes does not yet exist; this task builds it in shell for the installer only and does not add a Go one.

## Out of scope

Homebrew tap channel behavior. The FO boot upgrade-hint journey (d2k). Marketplace README repair (separate task).

Added during ideation, as alternatives a reviewer may reasonably raise:

- **Changing goreleaser's asset naming** so prereleases publish an unsuffixed archive. This would make the seeded one-line fix work, but the `_stable` suffix is deliberate: `.goreleaser.yaml:83` records that it exists so no prerelease publishes a default-looking asset, "closing the trap where that asset silently held the stable-stamped binary." Reopening that trap to save a `case` statement in `install.sh` trades a visible 404 for a silent wrong-binary install. The release artifacts stay as they are; the installer learns about channels.
- **A Go newest-release query.** d2k needs one for the FO upgrade hint. This task needs it in shell, before any binary exists on the machine, so a Go helper cannot serve the installer. Building one here would be speculative work for a task that has not reached ideation.
- **A prerelease-aware semver sorter in POSIX shell.** See residual risk in the spike section.

## Expected surface and tolerance

Estimate net LOC change: **+155**, across **5 files**. Insertions ~+161, deletions ~-6 (the replaced `resolve_latest_tag` / `asset_name` / glob lines, the `next` URL in the header, and the `(tracks next)` comment line).

| File | Net | What |
| --- | --- | --- |
| `install.sh` | +30 | channel variable + validation, three channel branches, one disclosure line, header comment for the new grammar, `next`→`main` URL fix |
| `internal/release/install_channel_test.go` (new) | +85 | offline fixture tests: channel asset selection both ways, edge tamper leg, bogus-channel abort, stable golden; plus the live edge tag-resolution test |
| `.github/workflows/install-e2e.yml` | +18 | an edge leg beside the existing stable leg, on both runners |
| `docs/site/get-started/install.md` | +14 | the edge binary path, plus the `(tracks next)`→`(tracks main)` comment fix (doc diff items 1 and 4) |
| `docs/releasing.md` | +8 | edge binary install paths in "Advancing the Edge Line" (doc diff below) |

The captain-ordered amendment (doc diff item 4) is a one-line replacement in a file already in the table, so it is **net +0** — one insertion against one deletion, no new file. The +155 total and the ±40 LOC / ±1 file tolerance are unchanged, and the file count stays at 5.

**Tolerance: ±40 net LOC and ±1 file.** The extra file allowance is for the live edge resolution test landing in the existing `internal/release/install_url_test.go` instead of the new file, if the shared `runInstallPrintTarget`/`fetchLatestRelease` helpers are cleaner to extend in place than to export.

This revises the seeded estimate (~+40 across 2 files) upward, and the reason is on the record above: the seed assumed the fix was tag resolution alone. Spike leg F showed that alone 404s, so the channel has to reach `asset_name()` and the local-dist glob too, and the offline test surface that pins channel selection is most of the added lines. The revision is a correction to a pre-spike guess, not scope growth — no capability was added beyond what AC-1 requires.

## Semantic changes

Declared, since files and lines do not catch these:

- **New command grammar.** `SPACEDOCK_CHANNEL` env var, values `stable` | `edge`, default `stable`. No flag is added: the surface is `curl … | sh`, where an env prefix composes and a `sh -s --` flag does not (spike leg O). No existing variable changes meaning; `SPACEDOCK_INSTALL_FROM`, `SPACEDOCK_INSTALL_VERSION`, `SPACEDOCK_INSTALL_DIR`, and `SPACEDOCK_PRINT_TARGET` keep their current semantics and compose with the new one.
- **New output.** One stderr line on the API-resolution path: `install.sh: resolved <channel> channel to <tag>`. stderr, not stdout, and `=`-free, so `SPACEDOCK_PRINT_TARGET` stdout stays machine-parseable (spike legs K, L, M). Existing output lines are unchanged.
- **New abort condition.** An unrecognized `SPACEDOCK_CHANNEL` exits non-zero and installs nothing. This is the only newly *rejected* input; every input accepted today is still accepted.
- **Which artifact a given invocation fetches.** On `SPACEDOCK_CHANNEL=edge` the resolved tag may be a prerelease and the asset carries `_edge`. **Unchanged when the variable is unset or `stable`** — same endpoint, same asset name, same URL, byte-identical print-target stdout (AC-4 locks this).
- **Not changed:** the checksum gate and its fail-closed behavior, tarball extraction, install directory and permissions, the PATH advisory, `--version` output, the plugin/marketplace channel logic in `internal/cli/host_exec.go`, and every Go package. This task touches no Go production code — only shell, a test file, CI, and docs.
- **Authority and stored formats:** unchanged. No new credential, no token, no auth; the added request is the same unauthenticated public endpoint class already in use. Nothing new is written to disk beyond the binary the script already installs.

## Doc diff

Concrete before/after for the user-visible changes. Implementation applies these.

### 1. `docs/site/get-started/install.md` — the "Binary (macOS / Linux)" tab (lines 17-23)

Before:

```markdown
=== "Binary (macOS / Linux)"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
    ```

    Installs a checksum-verified binary to `~/.local/bin`.
```

After:

```markdown
=== "Binary (macOS / Linux)"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
    ```

    Installs a checksum-verified binary to `~/.local/bin`. This is the stable
    channel: it resolves the latest stable release and prints the tag it
    resolved before it installs anything.

    To track edge, set `SPACEDOCK_CHANNEL=edge`. Edge resolves the newest
    release including prereleases, and installs the edge-stamped binary:

    ```bash
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | SPACEDOCK_CHANNEL=edge sh
    ```

    Match the channel to the skills you run. The edge plugin requires an
    edge-line binary minor, so a stable binary with edge skills aborts at
    first-officer boot. On macOS,
    `brew install spacedock-dev/tap/spacedock@next` is the Homebrew
    equivalent; casks are macOS-only, so on Linux this script is the only
    edge path.
```

The `Skills` section below already documents both channels for the plugin, which is the asymmetry this closes: the page tells you how to get the edge *plugin* but not the edge *binary*. No change needed there.

### 2. `docs/releasing.md` — end of "Advancing the Edge Line" (after line 235)

Before: the section ends at the `GORELEASER_CURRENT_TAG` paragraph.

After: append one paragraph.

```markdown
Both edge binary install paths track those prerelease tags: the `spacedock@next`
cask, bumped on every tag including prereleases, and
`SPACEDOCK_CHANNEL=edge … install.sh | sh`, which resolves the newest release
including prereleases and fetches that tag's `_edge` asset. The default stable
path resolves `/releases/latest`, which excludes prereleases — so a default
install on an edge machine lands a lower minor and aborts at first-officer boot.
Casks are macOS-only; on Linux the script is the only edge binary path.
```

### 3. `install.sh` header comment (lines 5-6 and the behavior block at 10-19)

Line 6, `next` → `main`, matching `internal/contract/contract.go:239` and `CLAUDE.md:25`:

```diff
-#   curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/next/install.sh | sh
+#   curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
```

And the new grammar in the behavior block:

```diff
+# SPACEDOCK_CHANNEL selects the release channel: `stable` (default) resolves the
+#   latest stable release and the unsuffixed asset; `edge` resolves the newest
+#   release INCLUDING prereleases and the `_edge` asset. Any other value aborts —
+#   a typo must not silently install the other channel. The resolved tag is
+#   printed to stderr before any download.
```

### 4. `docs/site/get-started/install.md` — the `## Skills` edge snippet comment (line 50)

Folded in by captain order from a sibling task. Same defect class as the `install.sh:6` header fix above: an install surface naming a ref that is not the one resolved.

```diff
-# Edge (tracks next) — marketplace named `spacedock-edge`, entry still `spacedock`
+# Edge (tracks main) — marketplace named `spacedock-edge`, entry still `spacedock`
```

The claim is stale twice over, verified independently during this amendment (2026-08-24): the live edge-branch manifest at `spacedock-dev/marketplace@edge` declares marketplace `spacedock-edge` holding entry `spacedock` with `"ref": "main"`, and `git ls-remote --heads …/spacedock.git next` returns nothing — **no `next` branch exists on the repo at all**. So the comment names a resolution target that is both wrong and nonexistent. This matches `CLAUDE.md:25` and `docs/releasing.md:192` ("The edge marketplace entry resolves `main` directly"), and the marketplace repo's own README was already repaired to match (`spacedock-dev/marketplace` main `403c46f`, edge `d390cad`).

The surrounding prose in that section is already correct — it explains the channel as the marketplace *name* and never repeats the `next` claim — so this one comment line is the whole fix. Note the `@edge` in the adjacent `claude plugin marketplace add spacedock-dev/marketplace@edge` is correct and must NOT be touched: that is the marketplace repo's branch, which is where the edge entry lives, and is unrelated to the ref that entry resolves.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (value) — A fresh prefix taking the edge install path holds a binary whose major.minor equals the minor the first-officer boot gate requires, stamped to the edge channel.**
Verified by: install into an empty prefix with `SPACEDOCK_CHANNEL=edge`, then run `<prefix>/spacedock --version` and assert line 1's major.minor equals `release.ProseMinor("skills/first-officer/references/first-officer-shared-core.md")`, and that the output carries `Channel: edge`. The baseline is independent and can move the wrong way twice over: it is read from the repo's own prose stamp rather than hardcoded, so it tracks the release line, and the measured value is the *installed artifact*, not the script's intent. Measured pre-change at 0.26.0 with no `Channel:` line (spike leg 2) against a required 0.27 — the gap this closes. Fails if the edge path lands a stable-stamped binary, a lower minor, or nothing.

**AC-2 (value) — No install can resolve a version without saying which one: the API-resolution path names the channel and tag it resolved, on both channels, before any download.**
Verified by: run with `SPACEDOCK_PRINT_TARGET=1` on each channel and assert stderr carries `resolved stable channel to <tag>` / `resolved edge channel to <tag>` with a non-empty `v`-prefixed tag, emitted before the download step. Fails if a version skew can occur with no install-time signal — the silence that made the original defect invisible.

**AC-3 — The edge channel resolves the newest release including prereleases, and constructs a URL the release actually publishes.**
Verified by: a live print-target run with `SPACEDOCK_CHANNEL=edge`, asserting the resolved tag equals the first entry of `/releases?per_page=1` and the constructed tarball URL matches that release's published `browser_download_url` byte-for-byte. Skips on an unreachable API, per the existing pattern at `install_url_test.go:111`. Fails if edge resolves the stable tag, or builds an unsuffixed asset name — which 404s on a prerelease (spike leg F), the failure the byte-for-byte match catches that a name-shape assertion would not.

**AC-4 — The stable path is byte-identical to today's.**
Verified by: `SPACEDOCK_PRINT_TARGET=1` stdout on the default path compared against a golden captured from the pre-change script for the same host (already IDENTICAL in spike leg L; the test locks it), plus the local-dist stable glob still selecting the unsuffixed asset from a dist directory containing both. Fails if the channel work perturbs the stable asset name, URL, endpoint, or glob — i.e. if this task makes existing stable installs a variable.

**AC-5 — The edge channel is behind the same fail-closed checksum gate; a tampered edge tarball installs nothing.**
Verified by: the existing `buildInstallFixture` harness extended with an `_edge` asset, tampered after `checksums.txt` is computed; assert non-zero exit and an empty install dir. Fails if the edge branch bypasses, forks, or duplicates the verify block rather than flowing through it.

**AC-6 — An unrecognized channel value aborts instead of silently installing the other channel.**
Verified by: `SPACEDOCK_CHANNEL=egde` exits non-zero, names the accepted values, and creates no install dir. Fails if a typo'd channel resolves stable — reproduced in the ideation probe (spike leg P) and the same silent-skew class this task closes, so a defaulting `case *)` must not survive.

## Test plan

Three tiers. The channel *selection* is proven offline and deterministically; only tag *resolution* needs the live API, and that leg skips rather than reds on network failure — the split the existing tests already use.

**Offline, deterministic — `internal/release/install_channel_test.go` (new), Go, low cost.** Extends `buildInstallFixture` (`install_checksum_gate_test.go:33`) to write both an unsuffixed and an `_edge` tarball plus a `checksums.txt` covering both, then drives `install.sh` against it via `SPACEDOCK_INSTALL_FROM`:
- edge channel selects the `_edge` asset and the installed binary prints the edge fixture's marker (AC-1's mechanism, AC-4's converse)
- default/`stable` selects the unsuffixed asset from that same both-assets directory (AC-4)
- tampered `_edge` tarball → non-zero exit, empty install dir (AC-5)
- `SPACEDOCK_CHANNEL=egde` → non-zero exit, accepted values named, empty install dir (AC-6)
- stable print-target stdout matches a checked-in golden (AC-4)
No network, no goreleaser, runs under `go test ./...` and `-race`. This is where the channel branches are actually pinned, because a fixture can hold both assets at once and force the selection to be a choice.

**Live, skip-on-unreachable — same file or `install_url_test.go`, Go, low cost.** Mirrors `TestInstallScriptResolvesLiveReleaseAsset`, reusing `runInstallPrintTarget` and `fetchLatestRelease`: edge print-target resolves the newest listed release and the constructed URL matches its published `browser_download_url` (AC-3). One unauthenticated request, measured cost-neutral (spike leg A). Skips on the existing `could not resolve the latest release tag` signal, which is also what a 403 rate-limit produces (spike leg H), so a throttled runner skips rather than reds.

**CI, both runners — `.github/workflows/install-e2e.yml`, low cost.** The workflow already builds a goreleaser `--snapshot` dist and runs the stable leg on `ubuntu-latest` and `macos-latest`. A snapshot dist carries both channels' archives, and its version carries the current minor, so an added edge leg asserts the full AC-1 property offline: `SPACEDOCK_CHANNEL=edge SPACEDOCK_INSTALL_FROM=$PWD/dist sh ./install.sh`, then `--version` reports minor equal to the prose stamp and `Channel: edge`. This is what makes AC-1 a repeatable gate rather than a one-time observation, and it covers Linux — the platform with no other edge path.

**Live fresh-prefix run, manual — no new code.** Already executed in ideation on darwin/arm64 (spike legs 1, 2, O): clean prefix, real API, real download, `0.27.0-pre8` + `Channel: edge`. Implementation re-runs it on the final script as the AC-1 live observation and records the `--version` output. The CI leg above is the deterministic twin; this leg is the one that proves the *live* resolution and the download actually compose, which no fixture can.

No live workflow (first-officer boot) test is needed: AC-1 asserts the installed binary's minor against the very prose stamp the boot gate reads, so it measures the gate's own comparison without paying for a workflow run.

## Stage Report: ideation

- DONE: Chosen mechanism decided with the riskiest claim exercised first: how an edge machine resolves the newest prerelease tag (GitHub releases list vs /latest, auth-free rate limits included) - spike output recorded in the task body, or an auditable "no spike needed" naming the proven mechanisms it relies on.
  Spiked live and unauthenticated against the real repo; 12 lettered legs recorded under "Spike: how an edge machine resolves the newest prerelease". `/releases?per_page=1` -> v0.27.0-pre8 vs `/latest` -> v0.26.0, both exactly 1 request against a measured `x-ratelimit-limit: 60`, so cost-neutral.
- DONE: Task body carries a value-measuring AC (a fresh-machine edge path lands the current prerelease, observed via spacedock --version) plus install-time visibility of the resolved tag, a net-LOC expected surface with tolerance, and declared semantic changes (install.sh flag/env grammar and output).
  AC-1 measures the installed binary's major.minor against `release.ProseMinor(first-officer-shared-core.md)` — the boot gate's own comparison, read from the repo rather than hardcoded — plus `Channel: edge`; AC-2 covers disclosure. Surface table gives +155 net / 5 files, tolerance +/-40 LOC and +/-1 file, with the upward revision from the seeded ~+40 explained. "Semantic changes" declares the `SPACEDOCK_CHANNEL` grammar, the new stderr line, the new abort, and an explicit not-changed list.
- DONE: A concrete doc diff for the user-visible changes (install docs / releasing docs before-after wording) is recorded in the task body per the ideation doc-diff rule.
  "Doc diff" section: full before/after for the `docs/site/get-started/install.md` binary tab, an append paragraph for `docs/releasing.md` "Advancing the Edge Line", and a unified diff for the `install.sh` header (including the `next`->`main` raw-URL correction).

### Summary

The spike overturned the seeded design. Fixing tag resolution alone does not work: every release publishes a *pair* of per-arch archives, the edge one always suffixed `_edge` and the stable one suffixed `_stable` on prerelease tags, so the name `install.sh` builds today returns HTTP 404 on the newest prerelease (leg F) while `_edge` returns 200. The channel therefore has to reach `asset_name()` and the local-dist glob, not just the endpoint — which is also why the estimate moved from ~+40/2 files to +155/5 files. A corollary worth the gate's attention: `install.sh` cannot install an edge-stamped binary at *any* tag today, and since Homebrew casks are macOS-only, a Linux edge machine has no scripted binary path at all.

Chosen mechanism is one env var, `SPACEDOCK_CHANNEL` (stable default | edge), threaded through resolution, asset naming, and the dist glob, plus one stderr line disclosing the resolved tag. Proven end to end on a clean prefix: `spacedock 0.27.0-pre8` / `Channel: edge`, against 0.26.0 with no channel line on the old path. Two regressions were checked and are clean — stable print-target stdout is byte-identical, and the new line is stderr and `=`-free so it cannot pollute the existing parser. The probe also exposed a defect in itself: a `case *)` default silently installs stable on a typo'd channel value (leg P), which is the exact silent-skew class this task closes, so AC-6 requires validation instead. Residual risk is declared, not eliminated: edge resolution trusts the API's created_at-descending order rather than sorting semver in shell, mitigated by the now-printed tag and an escape-hatch pin that the spike proved already works through existing variables at zero new cost. Ideation touched no code — `git status` over `install.sh`, `internal/`, `docs/site/`, and `.github/` is clean; the throwaway probe lived in a scratchpad.

### Amendment: captain-ordered doc-diff addition (2026-08-24, pre-gate)

Folded in from a sibling task by captain order, before the ideation gate decision (the gate was withdrawn and will be re-prepared).

- DONE: add the `docs/site/get-started/install.md:50` `(tracks next)` -> `(tracks main)` correction to the doc diff.
  Recorded as doc-diff item 4 with before/after. Both halves of the staleness claim were verified independently rather than relayed: the live `spacedock-dev/marketplace@edge` manifest declares marketplace `spacedock-edge` holding entry `spacedock` with `"ref": "main"`, and `git ls-remote --heads .../spacedock.git next` returns nothing, so the named branch does not exist. Item 4 also flags that the adjacent `marketplace@edge` ref is correct and must not be swept up in the fix.
- DONE: reconcile the surface table.
  A one-line replacement in a file already listed, so net +0: the `install.md` row and the +155 / 5-file total and the +/-40 LOC, +/-1 file tolerance all stand. Only the gross split moved, +161/-6, now noted in the estimate line.
- DONE: confirm ACs and test plan need no change.
  Confirmed, none changed. AC-1 through AC-6 are properties of `install.sh` behavior; this is prose in a docs page that no AC asserts and no test reads. `internal/contractlint/workflow_trunk_test.go` lints stray `next` refs in *workflows*, not docs, so nothing existing covers or breaks on it either. Adding a lint for one comment line would be disproportionate; the ideation doc-diff rule already routes it to implementation.

Scope note for the gate: this is a fourth doc hunk riding the existing doc-diff machinery, not a design change. The chosen mechanism, the spike evidence, and the ACs are untouched.
