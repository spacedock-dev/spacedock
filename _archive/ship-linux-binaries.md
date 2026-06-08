---
id: v3yf6j3nkey6tvxvsa1y42sm
title: Ship Linux binaries + a Linux install path (release is darwin-only)
status: done
source: "captain (2026-06-08) - dispatch before 0.20.0. Mirrors GitHub spacedock-dev/spacedock#321 (public report). The release is darwin-only though the code already builds + tests on Linux in CI; the gap is release config + a Linux distribution channel."
started: 2026-06-08T18:17:44Z
completed: 2026-06-08T22:41:53Z
verdict: PASSED
score:
worktree:
issue: "spacedock-dev/spacedock#321"
sprint: 0199-pre-flip-mechanics
group: distribution
sprint-readiness: ready
mod-block:
pr: "#332"
archived: 2026-06-08T22:41:53Z
---

Ship Linux binaries and give Linux users an install path, before the 0.20.0 cut. The code is already portable — the gap is release packaging + distribution.

## Problem

`.goreleaser.yaml` builds `goos: [darwin]` / `goarch: [arm64, amd64]` only, and `release.yml`'s goreleaser job runs on `macos-latest`. A `v*` release publishes ONLY macOS tarballs + a macOS Homebrew cask. There is NO Linux install path: the cask is macOS-only, and `go install …@latest` is a dev/toolchain fallback (unstamped version). So Linux users can't install spacedock.

The code already works on Linux: `runtime-live-e2e.yml`'s offline gate (`runs-on: ubuntu-latest`) runs `go build ./...` + `go test ./...`, and other jobs `go build -o ./spacedock` on ubuntu. So building a Linux binary is a trivial cross-compile (CI proves it). The work is release config + distribution, not portability.

## Proposed approach

1. **Build (one-liner):** add `linux` to goreleaser `builds.goos` (amd64 + arm64) and the matching `archives` (the existing `name_template` `{ProjectName}_{Version}_{Os}_{Arch}` already yields `spacedock_{version}_linux_{arch}.tar.gz`). Pure-Go cross-compile from the existing runner.
2. **Distribution — a `curl | sh` installer (recommended):** add `install.sh` that detects OS/arch, downloads the matching tarball from the latest GitHub Release (verifying against `checksums.txt`), and installs `spacedock` to `~/.local/bin` (or `/usr/local/bin`). Universal — works on Linux AND macOS; brew stays the mac-preferred path. No Linuxbrew requirement, no new package infra. Document the Linux/script path in `docs/install-journey.md`.
   - Alternatives considered (ideation can revisit): Homebrew formula on Linuxbrew (loses the cask's Gatekeeper postflight, uncommon); deb/rpm/AUR packages (heavier maintenance); `go install` (dev-only, unstamped).
3. **Linux runtime caveats:** the host CLIs (Claude Code / Codex / Pi) run on Linux. The macOS release steps (adhoc signing, `com.apple.quarantine` / Gatekeeper, notarization) do NOT apply to Linux, so the cask's `xattr -dr com.apple.quarantine` post-install has no Linux analogue (and is not needed). The bare `spacedock` binary is statically linked (`CGO_ENABLED=0`), so it runs on Linux with no shared-library dependencies (the spike's `file` output confirms `statically linked` ELF).
4. **safehouse on Linux (stated honestly, not implied to work):** spacedock does NOT implement the sandbox — `internal/safehouse` only (a) detects a `.safehouse` profile in the workdir, (b) checks whether a `safehouse` binary is resolvable on `PATH`, and (c) wraps the inner command as `safehouse --trust-workdir-config … -- <inner>`. That seam is pure argv composition and platform-agnostic: it behaves identically on Linux. **What is NOT proven by this task** is whether the external `safehouse` binary itself runs on Linux — that is a safehouse-side concern (the install hint points users at the upstream safehouse repo). So the honest Linux story is: *spacedock's sandbox integration works on Linux exactly as on macOS, but sandboxing only actually happens if a Linux-capable `safehouse` binary is installed; when it is absent spacedock prints its existing install-hint and the user can launch unsandboxed.* The README's "Native sandbox integration — drop a profile and Spacedock runs the agent sandboxed" claim stays true on Linux only to the extent the `safehouse` binary is present; the docs must not imply spacedock ships or guarantees a Linux sandbox. This task ships no safehouse code — it only states this fact in `docs/install-journey.md` (and does not over-claim).

## Out of scope

- The two-channel brew tap / per-channel `devBranch` stamp / next-publish pipeline (separate flip release-mechanics).
- safehouse's own Linux sandbox implementation (a safehouse-side concern) — this task only states the Linux story honestly.
- deb/rpm/AUR native packages and a Linuxbrew formula (heavier maintenance; revisit later if demanded). The `curl | sh` installer covers the Linux need without new package infra.
- Notarization / Developer-ID signing (a macOS-only concern, already tracked separately).

## Riskiest-mechanism spike (done — result on the record)

The design's load-bearing unverified mechanism was: *does goreleaser actually cross-compile and archive `linux/amd64` + `linux/arm64` tarballs from the existing runner, given only a one-line `goos` addition?* Paid that bill first.

**Exercise:** copied `.goreleaser.yaml` to a throwaway config with `linux` added to `builds.goos` (nothing else changed) and ran `goreleaser release --snapshot --clean --skip=publish,homebrew` (goreleaser v2.16.0, matching the config's `version: 2`) on the macOS dev host.

**Result (exit 0):** all four targets built (`darwin_amd64`, `darwin_arm64`, `linux_amd64`, `linux_arm64`) and four `tar.gz` archives + `checksums.txt` were produced. The archive name template yielded the predicted `spacedock_{version}_linux_{arch}.tar.gz`. `file` on the extracted binaries reports `ELF 64-bit LSB executable, x86-64, statically linked` (amd64) and `ELF 64-bit LSB executable, ARM aarch64, statically linked` (arm64) — genuine, dependency-free Linux binaries, with the bare `spacedock` at the archive root (alongside the auto-included `README.md`). The native darwin binary from the same run ran and reported its stamped version (`spacedock 0.19.8-snapshot-… (contract 1)`), confirming the version-stamp path the install AC depends on. `checksums.txt` covered all four tarballs.

This de-risks the entire build half: the change is genuinely a one-line `goos` addition, pure-Go cross-compile, no new toolchain. The install half's only unproven external is the GitHub "resolve latest release" call, addressed in the test plan below (also probed live: the API returns `tag_name` + per-asset `browser_download_url`, unauthenticated, and the current `v0.19.7` release confirms ONLY darwin assets exist today — the exact gap this task closes).

## Acceptance criteria

Each criterion is an end-state property with an independent, fail-able check (not a re-read of this task).

1. **A `v*` release publishes `linux/amd64` + `linux/arm64` tarballs alongside the darwin ones.**
   *Verified by:* the goreleaser config's `builds.goos` includes `linux`, exercised by `goreleaser release --snapshot --clean --skip=publish,homebrew` whose `dist/` contains `spacedock_*_linux_amd64.tar.gz` and `spacedock_*_linux_arm64.tar.gz`, each an extractable tarball with a bare `spacedock` ELF at root, and `checksums.txt` listing both. A guard test in `internal/release` asserts the config parses and its build target set includes the two linux pairs (the expected value — the linux/{amd64,arm64} target set — is independent of the YAML it checks: a config that drops linux fails it).

2. **`install.sh` installs a runnable `spacedock` on Linux AND macOS from a Release tarball, checksum-verified.**
   *Verified by:* a CI job on a `[ubuntu-latest, macos-latest]` matrix that (a) runs `goreleaser release --snapshot --clean --skip=publish,homebrew` so the runner's native-OS tarball exists locally, (b) runs `install.sh` pointed at that local dist (via an injectable source override), and (c) asserts the installed `spacedock --version` exits 0 and prints the stamped version. The macOS leg runs the darwin binary; the ubuntu leg runs the linux binary — both real, no mocks. A tamper case (corrupt the tarball or its checksum line) asserts `install.sh` exits non-zero and installs nothing, proving the checksum gate is load-bearing.

3. **`install.sh` resolves and fetches the latest GitHub Release tarball for the detected OS/arch.**
   *Verified by:* a test that drives the installer's OS/arch detection (`uname -s`/`-m` → `{darwin,linux}` × `{amd64,arm64}`) and latest-version resolution against the GitHub API and asserts the constructed asset name + download URL match the live release's published asset for this host. (The production fetch path; the local-dist override in AC-2 isolates the OS-logic from network flakiness, this one proves the network path's URL construction.)

4. **`docs/install-journey.md` documents the Linux install path** (the `curl | sh` line) and states the safehouse-on-Linux story accurately.
   *Verified by:* the doc carries a runnable `curl … install.sh | sh` invocation for Linux and a safehouse note that does NOT claim spacedock ships a Linux sandbox (it states sandboxing requires a Linux-capable `safehouse` binary, matching `internal/safehouse`'s detect-and-wrap-only behavior). Checkable by a test/CI grep that the Linux install line is present AND that the doc does not assert an unqualified "sandboxed on Linux" claim. *(Note: a grep over a doc is a weak proof on its own — it asserts wording, not behavior. The load-bearing behavioral proofs are AC-1/AC-2/AC-3; AC-4 is documentation correctness, and the honesty of the safehouse wording is a human/reviewer judgment at the gate, not a machine assertion.)*

5. **The macOS release path is unregressed.**
   *Verified by:* the existing `internal/release` guard tests (`release_test.go`, `workflow_exec_guard_test.go`) stay green, and the snapshot still produces both darwin tarballs + the homebrew cask generation is unaffected (the `homebrew_casks` block is darwin-only by asset and is left untouched).

## Test plan

| What | How | Cost | Type |
| --- | --- | --- | --- |
| Linux tarballs are built (AC-1) | `goreleaser release --snapshot --clean --skip=publish,homebrew`; assert `dist/` has both linux tarballs + checksums; extract + `file` shows ELF. Already done as the spike — implementation re-runs it as the offline gate. | seconds (~7s observed) | offline goreleaser snapshot, run in CI |
| Config guard (AC-1) | Go test in `internal/release` that loads `.goreleaser.yaml` and asserts the build target set ⊇ {linux/amd64, linux/arm64, darwin/amd64, darwin/arm64}. Fails if a future edit drops linux. | ms | Go unit test |
| `install.sh` end-to-end on both OSes (AC-2) | CI matrix `[ubuntu-latest, macos-latest]`: snapshot → `install.sh --from <local-dist>` → `spacedock --version` exits 0. Real binaries, no mocks. | ~1–2 min/leg | CI job (new) |
| Checksum gate is load-bearing (AC-2) | Same harness, corrupt a tarball byte (or its checksum line); assert `install.sh` exits non-zero and the install dir is empty. | seconds | shell test in the CI job |
| OS/arch + latest-version resolution (AC-3) | Test driving detection + URL construction against the live GitHub API (unauthenticated); assert the asset name/URL matches the host's published asset. Skippable offline; runs in CI. | seconds | CI / integration test |
| install-journey carries the Linux path + honest safehouse note (AC-4) | grep/content check in CI that the `curl … | sh` Linux line is present and no unqualified "sandboxed on Linux" claim exists. | ms | content check |
| macOS path unregressed (AC-5) | `go test ./internal/release/...` stays green; snapshot still emits both darwin tarballs. | seconds | existing Go tests + snapshot |

**Implementation note on the `install.sh` source override (AC-2/AC-3 seam):** the installer needs one injectable input — *where to fetch the tarball from*. Production resolves the latest GitHub Release; the test passes a local snapshot `dist/` directory. This is the single seam that lets AC-2 run the full extract/checksum/install/`--version` path with a real native binary on each CI OS without first publishing a Linux release (chicken-and-egg), while AC-3 separately proves the production GitHub-resolution path. Keep the override a plain env var or flag (e.g. `SPACEDOCK_INSTALL_FROM=<dir|url>`), not a code branch that diverges behavior — the same extract/verify/install logic runs in both.

## Stage Report: ideation

- DONE: Exercise the riskiest mechanism FIRST: run goreleaser --snapshot (or equivalent) and record that it produces linux/amd64 + linux/arm64 tarballs alongside darwin — the artifact list is the evidence.
  Ran `goreleaser release --snapshot --clean --skip=publish,homebrew` (v2.16.0) against a throwaway config with `linux` added to `builds.goos`; exit 0, produced all four tarballs (`darwin_{amd64,arm64}`, `linux_{amd64,arm64}`) + `checksums.txt`; `file` confirms statically-linked ELF x86-64 / ARM aarch64 binaries with bare `spacedock` at root. Result recorded in the body's "Riskiest-mechanism spike" section.
- DONE: Design the curl|sh installer (OS/arch detect, fetch latest Release tarball, checksum-verify, install to ~/.local/bin) with AC + test plan proving a runnable spacedock on Linux AND macOS — not a doc about it.
  AC-2/AC-3 specify a `[ubuntu-latest, macos-latest]` CI matrix that snapshots, runs `install.sh` against the local native-OS tarball via a `SPACEDOCK_INSTALL_FROM` source override, and asserts `spacedock --version` exits 0 on each OS (real binaries, no mocks), plus a tamper case proving the checksum gate. GitHub latest-release resolution probed live (returns `tag_name` + per-asset `browser_download_url`, unauthenticated).
- DONE: State the safehouse-on-Linux story honestly (do not imply the sandbox works on Linux if it doesn't); keep it in scope as a stated fact.
  Read `internal/safehouse/safehouse.go`: spacedock only detects a `.safehouse` profile, checks the `safehouse` binary on PATH, and wraps argv — it ships no sandbox. Proposed-approach item 4 states this: the integration seam is platform-agnostic, but actual sandboxing on Linux requires a Linux-capable external `safehouse` binary; the docs must not imply spacedock guarantees a Linux sandbox.

### Summary

Fleshed out the entity with a concrete problem statement, proposed approach (one-line goreleaser `goos` add + a universal `curl | sh` installer), five entity-level ACs each with an independent fail-able check, and a test-plan table. Paid the riskiest bill first: a throwaway goreleaser `--snapshot` proved the linux tarballs build and are genuine static ELF binaries from the existing macOS runner — the build half is de-risked to a one-line change. Key honesty call recorded: spacedock's safehouse seam works identically on Linux, but spacedock ships no sandbox, so the Linux sandbox story is gated on an external Linux-capable `safehouse` binary and the docs must not over-claim. Skipped comm-officer polish: the document's value is its load-bearing technical precision (tables + AC blocks), not prose style.

## Stage Report: implementation

- DONE: AC-1 — `.goreleaser.yaml` `builds.goos` includes linux (amd64+arm64); `internal/release` guard asserts the parsed target-set ⊇ {linux,darwin}×{amd64,arm64}; offline snapshot produces both linux tarballs + checksums.
  Commit `e632462c`: added one `- linux` line to `builds.goos`; `goreleaser_guard_test.go` parses the goos×goarch cross-product (expected set written independent of the YAML) + an adversarial test proving a dropped-linux config reds it. Offline gate re-run: `goreleaser release --snapshot --clean --skip=publish,homebrew` exit 0 → 4 tarballs + checksums.txt; extracted linux/amd64 is `ELF 64-bit … statically linked`.
- DONE: AC-2 — `install.sh` installs a runnable `spacedock` on Linux AND macOS from a checksum-verified Release tarball; CI matrix + tamper case.
  Commit `8a2b625d`: `install.sh` (uname OS/arch detect, fail-closed sha256 gate, extract bare binary to ~/.local/bin) with one `SPACEDOCK_INSTALL_FROM` seam (same extract/verify/install path in prod + test, not a divergent branch). `.github/workflows/install-e2e.yml` `[ubuntu-latest, macos-latest]` snapshots → installs native tarball → asserts `spacedock --version` exits 0 + stamped, then a corrupted-tarball tamper case asserts non-zero + nothing installed. Verified locally on darwin: happy path exit 0 (`0.19.9-snapshot-… (contract 1)`); both tamper variants (corrupt tarball + corrupt checksum line) → exit 1, install dir never created.
- DONE: AC-3 — live GitHub latest-release OS/arch + URL construction test (the ONLY network path, isolated).
  Commit `8a2b625d`: `install_url_test.go` drives install.sh's production path via a `SPACEDOCK_PRINT_TARGET` inspection mode (resolves + prints, stops before fetch — no logic duplication, no divergent install branch) against the live API; asserts the constructed asset name + URL match the live release's published `browser_download_url` for this host. Ran live (0.29s round-trip) → PASS against `v0.19.8` darwin/arm64. Skips when the API is unreachable so flakiness never reds the suite.
- DONE: AC-4 — `docs/install-journey.md` carries the runnable Linux `curl … install.sh | sh` line + honest safehouse-on-Linux note; content check.
  Commit `631d3832`: new "Install on Linux (or macOS without Homebrew)" section with the `curl -fsSL …/install.sh | sh` line; safehouse note states spacedock ships no sandbox and a run is sandboxed only with a Linux-capable `safehouse` binary on PATH (no unqualified "sandboxed on Linux" claim). `install_doc_test.go` locks both halves; the over-claim guard is non-vacuous (adversarially trips on "runs sandboxed on Linux").
- DONE: AC-5 — macOS release path unregressed; existing `internal/release` guards green.
  `.goreleaser.yaml` functional diff is the single `+ - linux` line; `homebrew_casks`/`archives`/`checksum`/darwin goarch untouched. Snapshot still emits both darwin tarballs. Full `internal/release` suite green (57 tests). Pre-existing unrelated failure flagged below.

### Summary

Shipped Linux binaries + a universal install path: a one-line goreleaser `goos` add (de-risked in ideation) now publishes linux/{amd64,arm64} tarballs, and a new fail-closed `install.sh` installs a checksum-verified `spacedock` on both Linux and macOS, fronted by a `[ubuntu,macos]` CI matrix (happy path + tamper) and a live-API URL-construction test isolating the single network path. The darwin/homebrew path is untouched (single functional config line) and the whole `internal/release` suite is green. NOTE for validation: `go test ./...` surfaces ONE pre-existing failure in `internal/release`'s sibling package `internal/status` — `TestMigrationCheckFixturesParseConsistently` (a `session-date: 2026-06-08` parse inconsistency in `docs/roadmap/0198-pre-flip-hardening/debrief.md`); it fails identically on the base commit `83a95b86` with none of my changes present (verified in a detached base-commit worktree), is unrelated to this task, and touches no file I modified. The checksum gate was genuinely exercised under tamper (both tarball-corruption and checksum-line-corruption fail closed), per the high-stakes/adversarial-audit flag.

### Post-completion adjustment (FO polish)

Commit `a6bafcb3`: the AC-3 live test now `t.Skip`s — not `t.Fatal`s — when the install.sh subprocess's OWN live GitHub `curl` can't connect (the production path resolves the latest tag via curl; a transport failure previously hard-failed). The direct `fetchLatestRelease` already skipped on a net/http error; this closes the second network seam. A logic failure (malformed print output) still hard-fails, and the OS-logic + URL-construction assertions stay hard assertions — only the actual live fetch skips. Proven both ways: with a failing `curl` stub on PATH the test SKIPs (lane green); network-up it PASSes against the live `browser_download_url`.

## Stage Report: validation

- DONE: AC-1: re-run the offline gate `goreleaser release --snapshot --clean --skip=publish,homebrew` — confirm `dist/` has both `spacedock_*_linux_amd64.tar.gz` + `spacedock_*_linux_arm64.tar.gz` (extractable, bare static-ELF `spacedock` at root) + `checksums.txt` listing both; the `internal/release` config guard passes.
  Snapshot exit 0 (goreleaser v2.16.0): all 4 tarballs built; both linux tarballs extract to a bare `spacedock` — `file` reports `ELF 64-bit … x86-64, statically linked` (amd64) + `ELF 64-bit … ARM aarch64, statically linked` (arm64); `checksums.txt` lists all 4. Guard tests `TestGoreleaserBuildsLinuxAndDarwin` + `TestGoreleaserBuildGuardRejectsDroppedLinux` PASS (expected set written independent of YAML).
- DONE: AC-2: the CI job installs a runnable `spacedock --version` (exit 0, stamped) from a native snapshot tarball via the `SPACEDOCK_INSTALL_FROM` override AND the tamper case fails closed (non-zero exit, nothing installed).
  Exercised the darwin leg locally against the snapshot dist: happy path → `install.sh` exit 0, `spacedock --version` exit 0 = `0.19.9-snapshot-0a5d87a5 (contract 1)` (real binary). FOUR tamper variants all fail closed (exit 1, install dir never created): corrupt tarball byte, corrupt checksum line, stripped checksum line, empty checksums.txt. (ubuntu leg = same code path, CI-only.)
- DONE: AC-3: the live-API URL-construction test asserts the URL against the host's published asset and t.Skips (not t.Fatal) on network failure.
  `TestInstallScriptResolvesLiveReleaseAsset` PASS (0.59s) against live API: install.sh resolves tag `v0.19.8`, constructs `spacedock_0.19.8_darwin_arm64.tar.gz`, byte-matches the live published `browser_download_url`. Transport-failure handling proven via a failing `curl` stub → test SKIPs (lane `PASS ok`). Live release confirms ONLY darwin assets exist today — the gap this task closes.
- DONE: AC-4: `docs/install-journey.md` has the runnable Linux `curl … | sh` line + an honest safehouse-on-Linux note (no unqualified "sandboxed on Linux").
  Both doc tests PASS; both proven non-vacuous (injecting "sandboxed on Linux" reds the over-claim guard; removing the curl line reds the positive guard; both restored). Human-judgment read: the safehouse wording matches `internal/safehouse`'s detect-and-wrap-only behavior (`Present`/`Available`/`Wrap`) — names the required `safehouse binary` PATH dependency, makes no unqualified sandbox promise.
- DONE: AC-5: darwin tarballs + homebrew_casks unregressed; full `internal/release` suite green.
  Functional `.goreleaser.yaml` diff vs base `83a95b86` is the single `+ - linux` line; both darwin tarballs still emitted; `homebrew_casks` block + `com.apple.quarantine` postflight intact. Full `internal/release` suite green (40 PASS, package ok).
- DONE: Detached adversarial audit (v3 = release/CI machinery) on a SEPARATE throwaway checkout of the merge result.
  Audited `/tmp/v3-audit` (detached merge of v3 into main `b1bdcf19`, clean FF-less merge, no conflicts). Probe 1: dropping linux from `builds.goos` reds the guard. Probe 2: a CI-tamper-step replica PASSES on the real gate and FAILS/reds the lane when the gate is gutted (always-pass comparison) — the tamper case is genuinely load-bearing. Stripped/empty checksum cases also fail closed. REFUTED NOTHING MATERIAL. One polish note: the `[ -n "$expected" ]` guard is redundant with the comparison for the stripped-line case (defense-in-depth, clearer error) — keeping it is correct.

### Summary

VERDICT: PASSED. All five ACs verified with independent evidence outside the task body (snapshot dist + `file` output + exit codes + live-API match + non-vacuity probes), and the REQUIRED detached adversarial audit on the merge result refuted nothing material — the goreleaser config guard reds when linux is dropped, and the install.sh checksum gate fails closed under every tamper (corrupt tarball, corrupt/stripped checksum line, empty file) with the CI tamper step proven to red the lane when the gate is broken. The macOS/homebrew path is unregressed (single `+ - linux` functional config line). The one full-suite failure — `internal/status` `TestMigrationCheckFixturesParseConsistently` — is a pre-existing failure (reproduced identically on base commit `83a95b86` with none of v3's changes; v3 touches zero `internal/status` files), owned by f1, not a v3 regression.
