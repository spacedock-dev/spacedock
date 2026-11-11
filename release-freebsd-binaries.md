---
title: "Release FreeBSD binaries"
source: "Captain request during the 0.28.0 release, 2026-11-10; Subspace added FreeBSD release targets on 2026-11-05"
id: h0xv2p2a7kts1k124hjdgb2v
status: ideation
gates:
    version: 1
    records:
        - id: gate:h0xv2p2a7kts1k124hjdgb2v:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:h0xv2p2a7kts1k124hjdgb2v-backlog-1
              briefing:
                id: briefing:h0xv2p2a7kts1k124hjdgb2v:backlog:attempt-1:revision-1
                digest: sha256:f7e4c48a9345b32bd99144a6d3cb586efcc0a304692709cbbe31aa4732c73505
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:h0xv2p2a7kts1k124hjdgb2v:backlog:1
                briefing: briefing:h0xv2p2a7kts1k124hjdgb2v:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-11T05:43:35.152357Z"
                decision: approve
                reason: Captain asked for this task to be filed and dispatched; the cross-build already passes and the open questions are named for ideation.
                conn:
                    quote: file and dispatch freebsd build as well
                    source: captain chat, 2026-11-10
              application:
                target-stage: ideation
                state: consumed
        - id: gate:h0xv2p2a7kts1k124hjdgb2v:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:h0xv2p2a7kts1k124hjdgb2v-ideation-1
              briefing:
                id: briefing:h0xv2p2a7kts1k124hjdgb2v:ideation:attempt-1:revision-1
                digest: sha256:3d2b544b23bc9801f7ec996dc2caecb0ea226e37b8232aa040b7e3be5ea8b742
                room-ref: '@review/ideation/briefing-1'
started: 2026-11-11T05:43:41Z
---

Spacedock releases ship darwin and linux binaries only, and `install.sh` rejects any other OS. A FreeBSD user has to build from source. Subspace, which the Spacedock installer also installs, now builds FreeBSD amd64 and arm64 binaries; Spacedock should match so a FreeBSD host can install the whole kit.

## What is known

- `CGO_ENABLED=0 GOOS=freebsd go build ./cmd/spacedock` and `GOOS=freebsd go vet ./...` pass for amd64 and arm64 on the 0.28.0 tree. The release build is already `CGO_ENABLED=0`, and OS-specific code is split on the `unix` build tag, which includes FreeBSD.
- `.goreleaser.yaml` lists `darwin` and `linux` in the `goos` anchor shared by the stable and edge builds.
- `install.sh` `detect_os` maps only `Darwin` and `Linux`.
- Subspace cross-builds FreeBSD in CI as a compile check (output discarded), adds `freebsd/amd64` and `freebsd/arm64` to its release targets, and maps `FreeBSD` in its installer. Its latest full release, v0.16.4, has no FreeBSD assets; v0.17.0-beta.1 is the first candidate to carry them.

## Outcome

A tagged Spacedock release publishes `freebsd_amd64` and `freebsd_arm64` tarballs for both channels with checksums, `install.sh` on FreeBSD installs the matching one, and CI fails a change that breaks the FreeBSD build.

## Unknowns for ideation

1. Whether the goreleaser Homebrew cask step tolerates FreeBSD archives in the same release, or needs them excluded from the cask.
2. What the combined installer does on FreeBSD when Subspace or Compass has no FreeBSD asset: skip that extension with a named message, or fail.
3. Whether `install.sh` uses anything the FreeBSD base system lacks (checksum tool, `install` flags, `curl` versus `fetch`).
4. Which hosts run on FreeBSD (Claude Code, Codex, Pi), and what the docs should promise; a binary without a supported host is of little use.
5. Whether the Safehouse sandbox front door, a macOS mechanism, fails clearly on FreeBSD.
6. How to prove the installer on FreeBSD: a FreeBSD VM job in CI, or a compile check plus installer unit tests with a stubbed `uname`.

## Out of scope

Windows and other BSDs. Notarization and packaging for FreeBSD ports or pkg.

## Answers to the unknowns

Labels: SETTLED (code, source or probe), SETTLED BY SUBSPACE (Subspace's hand runs on FreeBSD 14.5 and 15.1 VMs, 2026-11-04, against an installer of the same shape), OPEN (needs a FreeBSD machine).

1. **Cask: SETTLED, no exclusion needed.** goreleaser's cask pipe filters archives with `artifact.ByGooses("darwin", "linux")` before it reads them, and routes only `darwin` and `linux` into the cask (`internal/pipe/cask/cask.go` lines 218 and 443-446 at v2.16.0, the local version; lines 257 and 503-505 on the latest, v2.20.1, which `~> v2` in CI can pick up). FreeBSD archives never reach the cask, so neither the `ErrMultipleArchivesSameOS` count nor the cask's URLs change. `goreleaser check` accepts the config with `freebsd` added to the `goos` anchor (v2.16.0, local). The full snapshot render did NOT run: the link step failed with `no space left on device` (this host had 172 MiB free). The rendered-cask assertion in the test plan is the first thing implementation runs.
2. **Extension without a FreeBSD asset: SETTLED, it already skips; the message is wrong.** `install_extra` resolves the tag, then `curl -f` gets 404, prints `subspace: download failed — skipped`, and `install_extras` adds `subspace failed to install; spacedock works without it`. Only a failed tag lookup sets `LOOKUP_FAILED`, so the run exits 0 (pinned by `TestInstallExitsZeroWhenAnExtensionDownloadFailsAfterItsTagResolved`). Live today: Subspace v0.16.4 (latest stable) and Compass v2.9.1 have no FreeBSD archive; Subspace v0.17.0-beta.1 has both. A FreeBSD stable install would therefore print "download failed" twice on every run and every `spacedock update`, which reads as a network fault. Design: name the absence and do not call it a failure.
3. **Base tools: SETTLED BY SUBSPACE, plus reading this script.** FreeBSD base lacks `curl`, `bash`, `git`, `jq`, `tmux`, `perl` and `shasum`; it has `fetch`, `tar`, `readlink`, `mktemp`, `awk`, `install`, and `sha256sum` in `/sbin`. Spacedock's `install.sh` needs `curl` (`need curl` stops with `required command not found: curl`), and prefers `sha256sum`, so it hashes with base tools. `git` is used only by the `newest_tag` fallback when GitHub does not answer; without it that fallback finds nothing and the existing rate-limit message prints. No Spacedock skill or Go code calls `shasum`, `bash`, `jq` or `tmux` (grep of `skills/`, `internal/`, `cmd/`). Review, which the installer also installs, needs `bash jq tmux` at run time (Subspace README). Not run on FreeBSD: the lines this script has beyond Subspace's (`ln -sf`, `basename`, `sort`, `cut`, `sed -E`) are POSIX and in base.
4. **Hosts: OPEN for Spacedock; facts known.** FreeBSD ports, read 2026-11-10: `misc/codex` 0.167.0 builds Codex natively from Rust. `misc/claude-code` 2.1.322 repackages Anthropic's Linux npm package to run on FreeBSD's node with a third-party shim for a private Bun API; it is not an Anthropic-published build. Subspace's hand run on 2026-11-04 instead saw `pkg install claude-code` require the Linux compatibility layer; the two disagree and only a FreeBSD machine settles it. Pi has no port. Nobody has run any host with the Spacedock plugin on FreeBSD. `spacedock install --host claude` drives `claude plugin marketplace add`, which probably needs `git` (not in base, not in the port's run dependencies); unverified.
5. **Safehouse: SETTLED, fails clearly.** The wrap path is OS-independent: a `.safehouse` profile or `--safehouse` flag calls `exec.LookPath("safehouse")`, and on a miss prints the pinned hint and returns 1 before launching anything. Probe on darwin with `PATH=/usr/bin:/bin spacedock claude --safehouse --skip-compat-check`: exit 1, no host started. Safehouse has no FreeBSD build, so FreeBSD always takes this path, exactly as Linux does today. Two cross-OS warts in the hint are listed under Findings; neither is FreeBSD-specific.
6. **Proof: compile check plus installer tests with a stubbed `uname`.** A FreeBSD VM lane is rejected for the reasons the captain ruled on for Subspace on 2026-11-05 ("yes cross compiling only"): it is a standing check that runs nothing the release ships beyond what a compile check covers, and it would need a third-party action or a QEMU runner. No FreeBSD VM can run on this host (disk). A one-off hand run by a person after the first release is a captain decision below.

## Reused from Subspace

Reused: cross-compile only, no FreeBSD runtime in CI; FreeBSD kept out of the Homebrew package; `FreeBSD) echo freebsd ;;` in `detect_os` with the unknown-OS refusal kept; `pkg install curl` documented instead of `fetch` support in the installer; a README claim that says what CI does not do; an agent session on FreeBSD as a person's non-blocking test.

Not reused:
1. **The dedicated cross-build CI step.** Subspace has a custom release pipeline, so it needed a separate `go build` loop. Spacedock's release is goreleaser, and `install-e2e.yml` already runs `goreleaser release --snapshot` on every pull request to `main`, on ubuntu and macos, building every target in `.goreleaser.yaml`. Adding `freebsd` to the `goos` anchor makes that existing step the FreeBSD compile check, with the release's own flags, ldflags and both channel builds. A second step would build the same thing twice.
2. **A Formula-exclusion code change.** Subspace's Formula is hand-generated; goreleaser's cask pipe excludes FreeBSD by itself (unknown 1). Only an assertion is added.
3. **The validator `shasum` fix.** Spacedock ships no script that calls `shasum`.
4. **"Tested by hand" in the README.** Subspace ran VMs; nothing of Spacedock has run on FreeBSD, so the docs must not say tested.

## Proposed approach

One slice.

1. **Release targets, both channels.** Add `- freebsd` to the `&channel-goos` anchor in `.goreleaser.yaml`. Both `spacedock-stable` and `spacedock-edge` read it, so each tag publishes `spacedock_<v>_freebsd_{amd64,arm64}.tar.gz` (with `_stable` on a `-pre` tag) and `…_freebsd_{amd64,arm64}_edge.tar.gz`, each with a `checksums.txt` line: 8 archives become 12. The header comments of `.goreleaser.yaml` and `release.yml` name FreeBSD (`TestReleaseHeaderNamesEveryBuildOS` requires the `release.yml` header to name every built OS, so it fails first).
2. **Casks.** No config change. Both casks keep exactly their four darwin and linux URLs; an assertion pins that.
3. **`install.sh` OS detection.** Add `FreeBSD) echo freebsd ;;` to `detect_os`. Unknown OSes still stop with `unsupported OS <name>` before anything is created. Comments and the usage line name FreeBSD.
4. **`install.sh` when an extension has no build for the host.** Separate the asset download from the checksum download, and classify the asset's `curl -f` exit: 22 (HTTP error, in practice 404 on a resolved tag) prints `<name> <tag> has no <os>_<arch> archive — skipped` and returns 0, so `install_extras` adds no "failed to install" line; any other exit keeps today's `download failed — skipped` and return 1. Exit status of the run is unchanged (0). A 5xx also exits 22 and would be mislabelled as absent; accepted, since the run's outcome is the same and the next run retries.
5. **Too-old-binary remedy.** `tooOldBinaryRemedy` in `internal/contract/contract.go` leads with `brew upgrade` on every OS except Linux, so a FreeBSD user would be told to run Homebrew. Change `if goos == "linux"` to `if goos != "darwin"`, so every non-macOS host gets the `curl … | sh` lead.
6. **FO install reference.** `skills/first-officer/references/fo-install.md` sends any OS other than Linux and macOS to "hint go build, ABORT". Before: `- **Linux, stable:**`, `- **Linux, edge:**`, `On Linux edge the assignment binds to \`sh\``. After: `- **Linux or FreeBSD, stable:**`, `- **Linux or FreeBSD, edge:**`, `On Linux or FreeBSD edge the assignment binds to \`sh\``. Commands unchanged.
7. **CI compile check.** No new step. `install-e2e.yml`'s snapshot builds FreeBSD from step 1 on. Its header comment gains one sentence saying so, so nobody narrows the snapshot to `--single-target` without seeing what it removes. On a tag, goreleaser builds every target before publishing, so a FreeBSD break also stops the release. `main` has no branch protection (GitHub: "Branch not protected"), so on a pull request the red check is advice, as every check is today.

Mechanisms and the value they serve: the exit-22 classification serves AC-4; the simplest alternative, leaving the message, prints a false network failure twice on every FreeBSD stable install and update, and a static per-OS skip list goes stale when Compass ships FreeBSD. The rendered-cask assertion serves AC-2; relying on goreleaser's source alone is not enough because CI floats `~> v2`. The CI check adds no mechanism.

## Documentation diff

`docs/site/get-started/install.md`:

```diff
-=== "Linux or macOS without Homebrew"
+=== "Linux, FreeBSD, or macOS without Homebrew"
 
     <div class="sd-agent-prompt sd-command-prompt">
       <code>curl -fsSL …/install.sh | SPACEDOCK_CHANNEL=edge sh&#10;export PATH=…</code>
       …
     </div>
+
+    FreeBSD's base system has no `curl`. Install it first, as root: `pkg install curl git`.
+
+    Spacedock publishes FreeBSD binaries for amd64 and arm64. They are cross-compiled
+    from the same source as the Linux and macOS binaries, and CI checks on every change
+    that they build. CI does not run them on FreeBSD, and no agent host has been tested
+    with Spacedock on FreeBSD. Review publishes FreeBSD builds from 0.17.0; Compass does
+    not. When an extension's release has no FreeBSD build, the installer skips it and
+    says so.
```

`README.md`, Install section, after the Homebrew block:

```diff
+On Linux and FreeBSD, install with the script in the
+[install guide](docs/site/get-started/install.md).
```

`docs/releasing.md`:

```diff
-- cross-builds darwin and linux (arm64 + amd64) tarballs plus `checksums.txt`,
+- cross-builds darwin, linux and freebsd (arm64 + amd64) tarballs plus `checksums.txt`,
@@
-Casks are macOS-only; on Linux the script is the only edge binary path.
+Casks are macOS-only; on Linux and FreeBSD the script is the only edge binary path.
```

`install.sh` comments: `# Usage (Linux or macOS):` becomes `# Usage (Linux, FreeBSD or macOS):`; `Detects OS (darwin|linux)` becomes `(darwin|freebsd|linux)`; `the only Linux install path` becomes `the only Linux and FreeBSD install path`.

## Captain decisions

1. **What the docs promise about hosts on FreeBSD.** The paragraph above promises the binary and the installer only, and says plainly that no host has been tested. It does not name Claude Code or Codex as working on FreeBSD. Recommendation: this wording. The alternative, naming `misc/codex` and `misc/claude-code` as routes, would commit the docs to hosts nobody has run with the plugin, and the Claude Code port is a third-party repackaging whose run mode is already in dispute (unknown 4).
2. **A person's hand run on FreeBSD after the first release.** Recommendation: non-blocking, like Subspace's human test: on a FreeBSD 15.x VM, `pkg install curl git`, the documented `curl | sh`, `spacedock --version`, `spacedock doctor`, then `pkg install codex` and one `spacedock codex` session. If it passes, the docs may say "tested by hand" with the version and date; if it fails, the failure becomes a task.

## Spike

1. `goreleaser check -f <config with freebsd in the goos anchor>`: 1 configuration file validated (goreleaser 2.16.0, local).
2. goreleaser cask filter read at v2.16.0 and on the current default branch (v2.20.1): darwin and linux only, both versions.
3. `goreleaser release --snapshot --skip=publish --clean` in a scratch clone with that config: failed at link with `no space left on device`; scratch removed. Not a config result. The cross-build itself passes per the seed (`CGO_ENABLED=0 GOOS=freebsd go build`, amd64 and arm64, 0.28.0 tree).
4. Safehouse probe as in unknown 5.

Riskiest unverified mechanism remaining: the real snapshot render with FreeBSD archives present. It is the first failing-first run in implementation and needs about 300 MB free.

## Acceptance criteria

1. **AC-1: A tagged release publishes 12 archives instead of 8, adding `freebsd_amd64` and `freebsd_arm64` for both channels, each with a `checksums.txt` line.** Baseline: v0.28.0 has 8 archives and no FreeBSD. Before merge: `TestGoreleaserBuildsLinuxAndDarwin` (renamed to cover FreeBSD) expects the 6 os/arch pairs; dropping `- freebsd` from the anchor fails it. The install-e2e snapshot log lists 12 archives. After the first tag: `gh release view <tag> --json assets` lists 12 archives, 4 of them FreeBSD, each named in `checksums.txt`.
2. **AC-2: Both rendered casks still offer exactly four URLs, none FreeBSD.** `TestRenderedCaskDependsOnAgentsview`'s `assertRenderedCask` gains: no `freebsd` in the body and exactly 4 `url "` lines. Falsifier: feeding it a rendered body with an added FreeBSD `url` line fails it. After the first tag: the homebrew-tap cask commit changes only version and sha256 lines.
3. **AC-3: With `uname -s` reporting FreeBSD, `install.sh` installs the `freebsd_<arch>` archive's binary; with an unknown OS it exits 1 naming the OS and creates no install directory.** New `internal/release/install_freebsd_test.go` with a `uname` shim and a dist holding darwin, linux and freebsd archives with distinct contents. Falsifiers: deleting the `FreeBSD)` line; mapping FreeBSD to `linux` (installs the linux decoy).
4. **AC-4: When an extension's resolved release has no archive for the host, `install.sh` names the extension, tag and platform, installs the rest, exits 0, and prints no "failed to install" line for it; a connection failure still prints `download failed`.** New test in the same file, with the existing curl stub returning exit 22 for a missing file. The existing download-failure test switches to a non-22 exit. Falsifier: reverting step 4 of the approach fails the new test.
5. **AC-5: A change that builds for darwin and linux but not FreeBSD turns the install-e2e check red on its pull request.** One-off CI run on a throwaway draft PR adding a `//go:build freebsd` file that references an undefined symbol: the snapshot step fails on both legs while `go test ./...` stays green; the PR is closed after. On a tag, goreleaser stops before publishing (read from goreleaser's build-then-publish order, not run).
6. **AC-6: The too-old-binary remedy on FreeBSD leads with the `curl … | sh` install command, not `brew upgrade`.** `version_message_test.go` gains a `freebsd` case; reverting the condition to `goos == "linux"` fails it.
7. **AC-7: The install guide, README, `docs/releasing.md` and `fo-install.md` carry the wording above, and the docs promise no host on FreeBSD.** Read at the PR tip; `docs.yml`'s strict build passes.
8. **AC-8: darwin and linux behaviour is unchanged.** `go test ./...` and `go test ./... -race` pass; install-e2e's existing install, edge and tamper legs pass; existing asset-name tests pass unchanged.

## Test plan

Failing-first order, all deterministic unless marked:

1. `goreleaser_guard_test.go`: add the two FreeBSD pairs to the expected set. Fails until the anchor gains `freebsd`. Then `TestReleaseHeaderNamesEveryBuildOS` fails until the `release.yml` header names FreeBSD. Cost: under 1 s.
2. `cask_dependencies_test.go`: the no-FreeBSD and four-URL assertions. Runs the real goreleaser snapshot when goreleaser is on PATH (implementation must run it locally; about a minute, about 300 MB disk). Its falsifier is a doctored body, since goreleaser itself drops FreeBSD; this pins goreleaser's behaviour across `~> v2` upgrades.
3. `install_freebsd_test.go`: the FreeBSD install and unknown-OS cases fail on `unsupported OS FreeBSD` before step 3 of the approach; the missing-extension case fails on `download failed` before step 4. Cost: under 1 s each. Primary proof owner for the installer remains `internal/release`'s install tests; these reuse their harness (`extrasRun`, `lookupRun`'s curl stub).
4. `version_message_test.go`: the `freebsd` case fails before step 5.
5. One-off live CI: the AC-5 throwaway PR.
6. After the first tag, one-off: AC-1 and AC-2 release checks, plus `tar -xzOf` of one FreeBSD archive piped to `file -`, expecting `ELF 64-bit … FreeBSD`.

`fo-install.md` and the docs are instruction and prose changes with no behavioural test: their effect on FreeBSD needs a FreeBSD host, and the gate reads the before/after above.

## Expected surface

Estimate net LOC change: +122, across 15 files (about 142 insertions, 20 deletions).

| File | Insertions | Deletions |
|---|---|---|
| `.goreleaser.yaml` | 2 | 1 |
| `.github/workflows/release.yml` | 1 | 1 |
| `.github/workflows/install-e2e.yml` (comment) | 2 | 0 |
| `install.sh` | 10 | 5 |
| `internal/contract/contract.go` | 3 | 3 |
| `skills/first-officer/references/fo-install.md` | 3 | 3 |
| `internal/release/goreleaser_guard_test.go` | 4 | 2 |
| `internal/release/cask_dependencies_test.go` | 12 | 0 |
| `internal/release/install_freebsd_test.go` (new) | 80 | 0 |
| `internal/release/install_extras_test.go` (stub) | 3 | 1 |
| `internal/release/install_lookup_test.go` | 2 | 1 |
| `internal/contract/version_message_test.go` | 8 | 0 |
| `docs/site/get-started/install.md` | 8 | 1 |
| `README.md` | 2 | 0 |
| `docs/releasing.md` | 2 | 2 |

Tolerance: report above +175 net or 17 files. Halt on touch: any new CI job, step or third-party action; any non-test Go file other than `contract.go`; the cask `homebrew_casks` block.

Observable semantics this task may change: the published asset set (8 to 12 archives per tag); `install.sh` accepts FreeBSD; `install.sh`'s message and failure classification for an extension with no build for the host (exit status unchanged); the too-old-binary remedy text on non-macOS, non-Linux hosts; the FO's install routing on FreeBSD. Unchanged: command grammar, stored formats, darwin and linux asset names, cask contents, channel routing, who can publish.

## Findings outside this task

1. The safehouse hint says "this directory has a .safehouse profile" when the launch was selected by `--safehouse` with no profile present (observed in the probe), and names a Homebrew command on every OS.
2. On a host where the installer skips Review or Compass, `spacedock install` still installs their plugins, whose skills will then not find the program. `doctor` reports the program as "not installed" outside Homebrew. Not investigated further.
3. This host's data volume had 139-234 MiB free during this stage; a goreleaser snapshot failed with ENOSPC. Other agents on the host share this.

## Stage Report: ideation

- DONE: Answer each of the six unknowns in the entity with evidence from the code, the goreleaser config and docs, or a local probe; say plainly which could not be settled without a FreeBSD machine
  "Answers to the unknowns": 1, 2, 3 and 5 settled (goreleaser source at v2.16.0 and v2.20.1, `goreleaser check`, install.sh reading, live release assets, safehouse probe exit 1); 4 is OPEN, and 3 rests partly on Subspace's VM runs.
- DONE: Read how Subspace shipped FreeBSD (source repo checkout at ~/git/spacedock-research/spacedock-subspace, archived task release-freebsd-binaries, PR 94) and reuse its decisions where they fit, naming any that do not
  "Reused from Subspace" lists six reused decisions and four not reused, with reasons (goreleaser snapshot already compiles every target; cask pipe self-excludes; no `shasum`; no hand run here).
- DONE: Design the change: release targets for both channels, the cask handling, install.sh OS detection and its behaviour when an extension has no FreeBSD asset, the CI compile check, and the doc changes
  "Proposed approach" steps 1-7 and "Documentation diff"; adds two gaps the seed did not name: the brew-only too-old-binary remedy and `fo-install.md` aborting on FreeBSD.
- DONE: Decide what the docs promise about hosts on FreeBSD and flag it as a captain decision if it commits to supporting a host that is unverified there
  "Captain decisions" 1: binary and installer only, no host promised; 2: non-blocking hand run after the first release.
- DONE: State numbered acceptance criteria, a test plan with a failing-first test per change, expected surface (files and net LOC with tolerance), and how a release would prove the FreeBSD assets
  AC-1 to AC-8 (AC-1 measures 8 to 12 archives against v0.28.0); test plan items 1-6; +122 net across 15 files, report above +175 or 17; release proof in AC-1, AC-2 and test plan item 6.

### Summary

The design adds `freebsd` to the shared `goos` anchor, which gives both channels FreeBSD archives. It reuses install-e2e's existing goreleaser snapshot as the CI compile check, so no CI step is added. goreleaser's cask pipe drops non-darwin, non-linux archives in its source, so the casks need no change, only an assertion. The real snapshot render did not run because the host disk was full (ENOSPC at link, 139-234 MiB free), so that render is the first thing implementation runs. Two captain decisions are open: what the docs promise about hosts, and a non-blocking hand run after release. The private journal server was unavailable, so insights are recorded here only.
